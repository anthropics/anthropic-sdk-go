package bedrock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/auth/bearer"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/telemetry"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

var DefaultEndpoints = map[string]bool{
	"/v1/messages": true,
}

func NewStaticBearerTokenProvider(token string) *bearer.StaticTokenProvider {
	return &bearer.StaticTokenProvider{
		Token: bearer.Token{
			Value:     token,
			CanExpire: false,
		},
	}
}

type eventstreamChunk struct {
	Bytes string `json:"bytes"`
	P     string `json:"p"`
}

// sseTranslatingBody converts an AWS binary EventStream response body into SSE
// wire format (text/event-stream), so that middleware outside the Bedrock
// adaptation and stream consumers observe the same response shape on Bedrock
// as on the first-party API.
type sseTranslatingBody struct {
	eventstream.Decoder

	rc  io.ReadCloser
	buf bytes.Buffer
	err error
}

func (b *sseTranslatingBody) Read(p []byte) (int, error) {
	// Buffered SSE bytes must drain before a translation error surfaces, so
	// events decoded ahead of a malformed frame still reach the consumer.
	for b.buf.Len() == 0 {
		if b.err != nil {
			return 0, b.err
		}

		msg, err := b.Decoder.Decode(b.rc, nil)
		if err != nil {
			b.err = err
			continue
		}
		b.translate(msg)
	}
	return b.buf.Read(p)
}

func (b *sseTranslatingBody) Close() error {
	return b.rc.Close()
}

func (b *sseTranslatingBody) translate(msg eventstream.Message) {
	messageType := msg.Headers.Get(eventstreamapi.MessageTypeHeader)
	if messageType == nil {
		b.err = fmt.Errorf("%s event header not present", eventstreamapi.MessageTypeHeader)
		return
	}

	switch messageType.String() {
	case eventstreamapi.EventMessageType:
		eventType := msg.Headers.Get(eventstreamapi.EventTypeHeader)
		if eventType == nil {
			b.err = fmt.Errorf("%s event header not present", eventstreamapi.EventTypeHeader)
			return
		}

		if eventType.String() == "chunk" {
			chunk := eventstreamChunk{}
			err := json.Unmarshal(msg.Payload, &chunk)
			if err != nil {
				b.err = err
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(chunk.Bytes)
			if err != nil {
				b.err = err
				return
			}
			b.emit(gjson.GetBytes(decoded, "type").String(), decoded)
		}

	case eventstreamapi.ExceptionMessageType:
		// See https://github.com/aws/aws-sdk-go-v2/blob/885de40869f9bcee29ad11d60967aa0f1b571d46/service/iotsitewise/deserializers.go#L15511C1-L15567C2
		exceptionType := msg.Headers.Get(eventstreamapi.ExceptionTypeHeader)
		if exceptionType == nil {
			b.err = fmt.Errorf("%s event header not present", eventstreamapi.ExceptionTypeHeader)
			return
		}

		// See https://github.com/aws/aws-sdk-go-v2/blob/885de40869f9bcee29ad11d60967aa0f1b571d46/aws/protocol/restjson/decoder_util.go#L15-L48k
		var errInfo struct {
			Code    string
			Type    string `json:"__type"`
			Message string
		}
		err := json.Unmarshal(msg.Payload, &errInfo)
		if err != nil && err != io.EOF {
			b.err = fmt.Errorf("received exception %s: parsing exception payload failed: %w", exceptionType.String(), err)
			return
		}

		errorCode := "UnknownError"
		errorMessage := errorCode
		if ev := exceptionType.String(); len(ev) > 0 {
			errorCode = ev
		} else if len(errInfo.Code) > 0 {
			errorCode = errInfo.Code
		} else if len(errInfo.Type) > 0 {
			errorCode = errInfo.Type
		}

		if len(errInfo.Message) > 0 {
			errorMessage = errInfo.Message
		}
		b.emitError(errorCode, errorMessage)

	case eventstreamapi.ErrorMessageType:
		errorCode := "UnknownError"
		errorMessage := errorCode
		if header := msg.Headers.Get(eventstreamapi.ErrorCodeHeader); header != nil {
			errorCode = header.String()
		}
		if header := msg.Headers.Get(eventstreamapi.ErrorMessageHeader); header != nil {
			errorMessage = header.String()
		}
		b.emitError(errorCode, errorMessage)
	}
}

// emitError emits an SSE error event shaped like the first-party API's, so
// stream consumers surface Bedrock stream errors as API errors.
func (b *sseTranslatingBody) emitError(errorType, message string) {
	data, _ := sjson.SetBytes([]byte(`{"type":"error","error":{}}`), "error.type", errorType)
	data, _ = sjson.SetBytes(data, "error.message", message)
	b.emit("error", data)
}

func (b *sseTranslatingBody) emit(eventType string, data []byte) {
	b.buf.WriteString("event: ")
	b.buf.WriteString(eventType)
	b.buf.WriteByte('\n')
	// The SSE format carries one "data:" line per payload line; consumers
	// rejoin multi-line data with '\n'. API event JSON contains no raw
	// newlines, but split defensively to keep the framing valid either way.
	for _, line := range bytes.Split(data, []byte("\n")) {
		b.buf.WriteString("data: ")
		b.buf.Write(line)
		b.buf.WriteByte('\n')
	}
	b.buf.WriteByte('\n')
}

// WithLoadDefaultConfig returns a request option which loads the default config for Amazon and registers
// middleware that intercepts request to the Messages API so that this SDK can be used with Amazon Bedrock.
//
// If you already have an [aws.Config], it is recommended that you instead call [WithConfig] directly.
//
// Register any [option.WithMiddleware] before this option so your middleware
// observes Anthropic-shaped requests and responses; see [WithConfig].
//
// Like [WithConfig], the returned option includes
// [option.WithoutEnvironmentDefaults], so first-party credential sources
// never apply to a client constructed with it.
func WithLoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) option.RequestOption {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic(err)
	}
	return WithConfig(cfg)
}

// withPlatform names Amazon Bedrock as the gen_ai.provider.name of the
// OpenTelemetry spans, and this package in their instrumentation scope.
var withPlatform = requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
	rc.OpenTelemetry.ProviderName, rc.OpenTelemetry.Platform = telemetry.ProviderBedrock, "bedrock"
	return nil
})

// WithConfig returns a request option that uses the provided config and registers middleware to
// intercept requests to the Messages API, enabling this SDK to work with Amazon Bedrock.
//
// Authentication is resolved once, in this order:
//  1. If the AWS_BEARER_TOKEN_BEDROCK environment variable is set, its value is sent as a bearer token.
//  2. If cfg.AuthSchemePreference lists "httpBearerAuth" first, cfg.BearerAuthTokenProvider is used.
//  3. If cfg.Credentials is set, requests are signed with AWS SigV4.
//  4. Otherwise cfg.BearerAuthTokenProvider is used, and one of the two must be set.
//
// [config.LoadDefaultConfig] sets cfg.BearerAuthTokenProvider from the SSO token cache for SSO
// profiles; step 3 keeps that token from being sent to Bedrock.
//
// The Bedrock adaptation (URL and body rewriting, request signing, and
// normalization of streaming responses to SSE) should run closest to the wire.
// Middleware runs in registration order, so register [option.WithMiddleware]
// before this option:
//
//	client := anthropic.NewClient(
//		option.WithMiddleware(loggingMiddleware),
//		bedrock.WithConfig(cfg),
//	)
//
// Ordered this way, your middleware observes Anthropic-shaped requests
// (POST /v1/messages with model and stream in the body, no AWS signature)
// and SSE-formatted streaming responses — identical to the first-party API.
// Note that mutating the request after the Bedrock middleware has signed it
// invalidates the SigV4 signature, so body- or header-mutating middleware
// must be registered before this option.
//
// Amazon Bedrock does not serve the User Profiles API, so a request to it
// returns an error instead of being sent. A user profile ID on a Messages
// request is sent as usual.
//
// The returned option is an [option.Join] of [option.WithoutEnvironmentDefaults]
// and the Bedrock configuration, so a client constructed with it skips the
// SDK's first-party credential autoload entirely: environment credentials
// and profiles from the shared config store (ANTHROPIC_API_KEY,
// ANTHROPIC_PROFILE, ANTHROPIC_CONFIG_DIR) never apply to Bedrock requests.
func WithConfig(cfg aws.Config) option.RequestOption {
	var credentialErr error

	cfg = resolveAuth(cfg)
	if cfg.BearerAuthTokenProvider == nil && cfg.Credentials == nil {
		credentialErr = fmt.Errorf("expected AWS credentials to be set")
	}

	signer := v4.NewSigner()
	middleware := bedrockMiddleware(signer, cfg)

	return option.Join(option.WithoutEnvironmentDefaults(), withPlatform, requestconfig.WithoutMessageThreads(), requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if credentialErr != nil {
			return credentialErr
		}
		return rc.Apply(
			option.WithBaseURL(fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", cfg.Region)),
			option.WithMiddleware(middleware),
		)
	}))
}

// resolveAuth returns cfg with BearerAuthTokenProvider set only when requests
// should carry a bearer token; a nil provider means SigV4 with cfg.Credentials.
func resolveAuth(cfg aws.Config) aws.Config {
	prefersBearer := len(cfg.AuthSchemePreference) > 0 && cfg.AuthSchemePreference[0] == "httpBearerAuth"

	token := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")

	switch {
	case token != "":
		cfg.BearerAuthTokenProvider = NewStaticBearerTokenProvider(token)
	case prefersBearer && cfg.BearerAuthTokenProvider != nil:
	case cfg.Credentials != nil:
		cfg.BearerAuthTokenProvider = nil
	}
	return cfg
}

func bedrockMiddleware(signer *v4.Signer, cfg aws.Config) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
		if strings.Contains(r.URL.Path, "/v1/user_profiles") {
			return nil, errors.New("bedrock: Amazon Bedrock does not serve the User Profiles API; manage user profiles with a first-party client")
		}

		// Bedrock takes the API version from the body.
		version := r.Header.Get("anthropic-version")
		r.Header.Del("anthropic-version")

		var body []byte
		if r.Body != nil {
			body, err = io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			r.Body.Close()

			// Raw (non-JSON) bodies must reach the wire byte-for-byte; only a
			// JSON body carries fields to inject or strip.
			isJSON := gjson.ValidBytes(body)

			if isJSON && version != "" && !gjson.GetBytes(body, "anthropic_version").Exists() {
				body, _ = sjson.SetBytes(body, "anthropic_version", version)
			}

			// pull the betas off of the header (if set) and put them in the body
			if betaHeader := r.Header.Values("anthropic-beta"); isJSON && len(betaHeader) > 0 {
				r.Header.Del("anthropic-beta")
				var betas []string
				for _, value := range betaHeader {
					for _, beta := range strings.Split(value, ",") {
						if beta = strings.TrimSpace(beta); beta != "" {
							betas = append(betas, beta)
						}
					}
				}
				if len(betas) > 0 {
					body, err = sjson.SetBytes(body, "anthropic_beta", betas)
					if err != nil {
						return nil, err
					}
				}
			}

			if isJSON && r.Method == http.MethodPost && DefaultEndpoints[r.URL.Path] {
				model := gjson.GetBytes(body, "model").String()
				stream := gjson.GetBytes(body, "stream").Bool()

				body, _ = sjson.DeleteBytes(body, "model")
				body, _ = sjson.DeleteBytes(body, "stream")

				var method string
				if stream {
					method = "invoke-with-response-stream"
				} else {
					method = "invoke"
				}

				r.URL.Path = fmt.Sprintf("/model/%s/%s", model, method)
				r.URL.RawPath = fmt.Sprintf("/model/%s/%s", url.QueryEscape(model), method)
			}

			reader := bytes.NewReader(body)
			r.Body = io.NopCloser(reader)
			r.GetBody = func() (io.ReadCloser, error) {
				_, err := reader.Seek(0, 0)
				return io.NopCloser(reader), err
			}
			r.ContentLength = int64(len(body))
		}

		// WithConfig resolves the auth mode: a provider here means bearer.
		if cfg.BearerAuthTokenProvider != nil {
			token, err := cfg.BearerAuthTokenProvider.RetrieveBearerToken(r.Context())
			if err != nil {
				return nil, err
			}
			r.Header.Set("Authorization", "Bearer "+token.Value)
		} else {
			ctx := r.Context()
			credentials, err := cfg.Credentials.Retrieve(ctx)
			if err != nil {
				return nil, err
			}

			hash := sha256.Sum256(body)
			err = awsauth.SignHTTP(ctx, signer, credentials, r, hex.EncodeToString(hash[:]), "bedrock", cfg.Region, time.Now())
			if err != nil {
				return nil, err
			}
		}

		res, err = next(r)
		if err != nil || res == nil {
			return res, err
		}

		// Normalize streaming responses to the SSE format the first-party API
		// uses, so layers above this middleware never see AWS EventStream.
		// Error responses stay untranslated: the SDK's error path reads the
		// raw body to build an error carrying the status and request ID.
		if mediaType, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type")); res.StatusCode < 400 && mediaType == "application/vnd.amazon.eventstream" {
			res.Body = &sseTranslatingBody{rc: res.Body}
			res.Header.Set("Content-Type", "text/event-stream")
			res.Header.Del("Content-Length")
			res.ContentLength = -1
		}

		return res, nil
	}
}

package bedrock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/auth/bearer"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const DefaultVersion = "bedrock-2023-05-31"

var DefaultEndpoints = map[string]bool{
	"/v1/complete": true,
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
	// events decoded ahead of a mid-stream exception still reach the consumer.
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
		b.err = fmt.Errorf("received exception %s: %s", errorCode, errorMessage)

	case eventstreamapi.ErrorMessageType:
		errorCode := "UnknownError"
		errorMessage := errorCode
		if header := msg.Headers.Get(eventstreamapi.ErrorCodeHeader); header != nil {
			errorCode = header.String()
		}
		if header := msg.Headers.Get(eventstreamapi.ErrorMessageHeader); header != nil {
			errorMessage = header.String()
		}
		b.err = fmt.Errorf("received error %s: %s", errorCode, errorMessage)
	}
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
func WithLoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) option.RequestOption {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic(err)
	}
	return WithConfig(cfg)
}

// WithConfig returns a request option that uses the provided config and registers middleware to
// intercept requests to the Messages API, enabling this SDK to work with Amazon Bedrock.
//
// Authentication is determined as follows: if the AWS_BEARER_TOKEN_BEDROCK environment variable is
// set, it is used for bearer token authentication. Otherwise, if cfg.BearerAuthTokenProvider is set,
// it is used. If neither is available, cfg.Credentials is used for AWS SigV4 signing and must be set.
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
func WithConfig(cfg aws.Config) option.RequestOption {
	var credentialErr error

	if cfg.BearerAuthTokenProvider == nil {
		if token := os.Getenv("AWS_BEARER_TOKEN_BEDROCK"); token != "" {
			cfg.BearerAuthTokenProvider = NewStaticBearerTokenProvider(token)
		}
	}
	if cfg.BearerAuthTokenProvider == nil && cfg.Credentials == nil {
		credentialErr = fmt.Errorf("expected AWS credentials to be set")
	}

	signer := v4.NewSigner()
	middleware := bedrockMiddleware(signer, cfg)

	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if credentialErr != nil {
			return credentialErr
		}
		opts := []option.RequestOption{
			option.WithBaseURL(fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", cfg.Region)),
			option.WithMiddleware(middleware),
		}
		if cfg.HTTPClient != nil {
			opts = append(opts, option.WithHTTPClient(cfg.HTTPClient))
		}
		return rc.Apply(opts...)
	})
}

func bedrockMiddleware(signer *v4.Signer, cfg aws.Config) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
		var body []byte
		if r.Body != nil {
			body, err = io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			r.Body.Close()

			if !gjson.GetBytes(body, "anthropic_version").Exists() {
				body, _ = sjson.SetBytes(body, "anthropic_version", DefaultVersion)
			}

			// pull the betas off of the header (if set) and put them in the body
			if betaHeader := r.Header.Values("anthropic-beta"); len(betaHeader) > 0 {
				r.Header.Del("anthropic-beta")
				body, err = sjson.SetBytes(body, "anthropic_beta", betaHeader)
				if err != nil {
					return nil, err
				}
			}

			if r.Method == http.MethodPost && DefaultEndpoints[r.URL.Path] {
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

		// Use bearer token authentication if configured, otherwise fall back to SigV4
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
			err = signer.SignHTTP(ctx, credentials, r, hex.EncodeToString(hash[:]), "bedrock", cfg.Region, time.Now())
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

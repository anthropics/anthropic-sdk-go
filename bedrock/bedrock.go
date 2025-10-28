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
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

const DefaultVersion = "bedrock-2023-05-31"

var DefaultEndpoints = map[string]bool{
	"/v1/complete": true,
	"/v1/messages": true,
}

// Context key for per-request extraBody
type contextKey string

const extraBodyKey contextKey = "bedrock_extra_body"

type eventstreamChunk struct {
	Bytes string `json:"bytes"`
	P     string `json:"p"`
}

type eventstreamDecoder struct {
	eventstream.Decoder

	rc  io.ReadCloser
	evt ssestream.Event
	err error
}

func (e *eventstreamDecoder) Close() error {
	return e.rc.Close()
}

func (e *eventstreamDecoder) Err() error {
	return e.err
}

func (e *eventstreamDecoder) Next() bool {
	if e.err != nil {
		return false
	}

	msg, err := e.Decoder.Decode(e.rc, nil)
	if err != nil {
		e.err = err
		return false
	}

	messageType := msg.Headers.Get(eventstreamapi.MessageTypeHeader)
	if messageType == nil {
		e.err = fmt.Errorf("%s event header not present", eventstreamapi.MessageTypeHeader)
		return false
	}

	switch messageType.String() {
	case eventstreamapi.EventMessageType:
		eventType := msg.Headers.Get(eventstreamapi.EventTypeHeader)
		if eventType == nil {
			e.err = fmt.Errorf("%s event header not present", eventstreamapi.EventTypeHeader)
			return false
		}

		if eventType.String() == "chunk" {
			chunk := eventstreamChunk{}
			err = json.Unmarshal(msg.Payload, &chunk)
			if err != nil {
				e.err = err
				return false
			}
			decoded, err := base64.StdEncoding.DecodeString(chunk.Bytes)
			if err != nil {
				e.err = err
				return false
			}
			e.evt = ssestream.Event{
				Type: gjson.GetBytes(decoded, "type").String(),
				Data: decoded,
			}
		}

	case eventstreamapi.ExceptionMessageType:
		// See https://github.com/aws/aws-sdk-go-v2/blob/885de40869f9bcee29ad11d60967aa0f1b571d46/service/iotsitewise/deserializers.go#L15511C1-L15567C2
		exceptionType := msg.Headers.Get(eventstreamapi.ExceptionTypeHeader)
		if exceptionType == nil {
			e.err = fmt.Errorf("%s event header not present", eventstreamapi.ExceptionTypeHeader)
			return false
		}

		// See https://github.com/aws/aws-sdk-go-v2/blob/885de40869f9bcee29ad11d60967aa0f1b571d46/aws/protocol/restjson/decoder_util.go#L15-L48k
		var errInfo struct {
			Code    string
			Type    string `json:"__type"`
			Message string
		}
		err = json.Unmarshal(msg.Payload, &errInfo)
		if err != nil && err != io.EOF {
			e.err = fmt.Errorf("received exception %s: parsing exception payload failed: %w", exceptionType.String(), err)
			return false
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
		e.err = fmt.Errorf("received exception %s: %s", errorCode, errorMessage)
		return false

	case eventstreamapi.ErrorMessageType:
		errorCode := "UnknownError"
		errorMessage := errorCode
		if header := msg.Headers.Get(eventstreamapi.ErrorCodeHeader); header != nil {
			errorCode = header.String()
		}
		if header := msg.Headers.Get(eventstreamapi.ErrorMessageHeader); header != nil {
			errorMessage = header.String()
		}
		e.err = fmt.Errorf("received error %s: %s", errorCode, errorMessage)
		return false
	}

	return true
}

func (e *eventstreamDecoder) Event() ssestream.Event {
	return e.evt
}

var (
	_ ssestream.Decoder = &eventstreamDecoder{}
)

func init() {
	ssestream.RegisterDecoder("application/vnd.amazon.eventstream", func(rc io.ReadCloser) ssestream.Decoder {
		return &eventstreamDecoder{rc: rc}
	})
}

// WithLoadDefaultConfig returns a request option which loads the default config for Amazon and registers
// middleware that intercepts request to the Messages API so that this SDK can be used with Amazon Bedrock.
//
// If you already have an [aws.Config], it is recommended that you instead call [WithConfig] directly.
func WithLoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) option.RequestOption {
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		panic(err)
	}
	return WithConfig(cfg)
}

// WithConfig returns a request option which uses the provided config and registers middleware that
// intercepts request to the Messages API so that this SDK can be used with Amazon Bedrock.
func WithConfig(cfg aws.Config) option.RequestOption {
	signer := v4.NewSigner()
	middleware := bedrockMiddlewareWithExtra(signer, cfg, nil)

	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		return rc.Apply(
			option.WithBaseURL(fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com", cfg.Region)),
			option.WithMiddleware(middleware),
		)
	})
}

// WithExtraBody returns a request option that adds extraBody fields to this specific request.
// The extraBody map will be merged into the request JSON before AWS signing.
//
// This is the recommended way to pass extraBody parameters like context_management.
// It provides more flexibility than WithConfigAndExtraBody since you can vary the
// extraBody per request.
//
// Example:
//
//	extraBody := map[string]any{
//	    "context_management": map[string]any{
//	        "edits": []map[string]any{{
//	            "type": "clear_tool_uses_20250919",
//	            "trigger": map[string]any{"type": "input_tokens", "value": 30000},
//	        }},
//	    },
//	}
//	response, err := client.Messages.New(ctx, params, bedrock.WithExtraBody(extraBody))
func WithExtraBody(extraBody map[string]any) option.RequestOption {
	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		if rc.Context == nil {
			rc.Context = context.Background()
		}
		rc.Context = context.WithValue(rc.Context, extraBodyKey, extraBody)
		// CRITICAL: Must also update the Request's context so middleware can access it
		if rc.Request != nil {
			rc.Request = rc.Request.WithContext(rc.Context)
		}
		return nil
	})
}

func bedrockMiddleware(signer *v4.Signer, cfg aws.Config) option.Middleware {
	return bedrockMiddlewareWithExtra(signer, cfg, nil)
}

func bedrockMiddlewareWithExtra(signer *v4.Signer, cfg aws.Config, staticExtraBody map[string]any) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
		// First check for per-request extraBody from context
		var extraBody map[string]any
		if r.Context() != nil {
			if val := r.Context().Value(extraBodyKey); val != nil {
				if eb, ok := val.(map[string]any); ok {
					extraBody = eb
				}
			}
		}

		// Fall back to static extraBody if no per-request one
		if extraBody == nil {
			extraBody = staticExtraBody
		}

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

			// Convert anthropic-beta header or beta query param to anthropic_beta field in body (similar to Python SDK)
			// The Go SDK uses MessagesStreamBeta which adds ?beta=true, but Bedrock needs this in the body
			betaHeader := r.Header.Get("anthropic-beta")
			betaQuery := r.URL.Query().Get("beta")

			if betaHeader != "" {
				// Split comma-separated betas into array
				betas := []string{}
				for _, b := range bytes.Split([]byte(betaHeader), []byte(",")) {
					betas = append(betas, string(bytes.TrimSpace(b)))
				}
				if len(betas) > 0 {
					betasJSON, _ := json.Marshal(betas)
					body, _ = sjson.SetRawBytes(body, "anthropic_beta", betasJSON)
				}
			} else if betaQuery == "true" {
				// Beta query param detected - infer which betas from extraBody
				// If context_management is present, we need context-management-2025-06-27
				inferred := []string{}
				if extraBody != nil {
					if _, hasContextMgmt := extraBody["context_management"]; hasContextMgmt {
						inferred = append(inferred, "context-management-2025-06-27")
					}
				}
				if len(inferred) > 0 {
					betasJSON, _ := json.Marshal(inferred)
					body, _ = sjson.SetRawBytes(body, "anthropic_beta", betasJSON)
				}
			}

			if r.Method == http.MethodPost && DefaultEndpoints[r.URL.Path] {
				model := gjson.GetBytes(body, "model").String()
				stream := gjson.GetBytes(body, "stream").Bool()

				body, _ = sjson.DeleteBytes(body, "model")
				body, _ = sjson.DeleteBytes(body, "stream")

				// WORKAROUND: Merge extraBody fields into request JSON
				// This allows sending fields like context_management that Bedrock
				// rejects when sent as standard SDK parameters (similar to Python's extra_body)
				for key, value := range extraBody {
					valueJSON, err := json.Marshal(value)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal extra_body field %s: %w", key, err)
					}
					body, err = sjson.SetRawBytes(body, key, valueJSON)
					if err != nil {
						return nil, fmt.Errorf("failed to set extra_body field %s: %w", key, err)
					}
				}

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

		return next(r)
	}
}

package anthropic_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"
)

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

// unreadBody is a response body that fails the test when it is read.
type unreadBody struct {
	t *testing.T
}

func (b unreadBody) Read([]byte) (int, error) {
	b.t.Error("the body of a response that is retried must not be read")
	return 0, io.EOF
}

func (b unreadBody) Close() error {
	return nil
}

func TestUserAgentHeader(t *testing.T) {
	var userAgent string
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					userAgent = req.Header.Get("User-Agent")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)
	_, _ = client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if userAgent != fmt.Sprintf("Anthropic/Go %s", internal.PackageVersion) {
		t.Errorf("Expected User-Agent to be correct, but got: %#v", userAgent)
	}
}

func TestRetryAfter(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}

	attempts := len(retryCountHeaders)
	if attempts != 3 {
		t.Errorf("Expected %d attempts, got %d", 3, attempts)
	}

	expectedRetryCountHeaders := []string{"0", "1", "2"}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestDeleteRetryCountHeader(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
		option.WithHeaderDel("X-Stainless-Retry-Count"),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}

	expectedRetryCountHeaders := []string{"", "", ""}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestOverwriteRetryCountHeader(t *testing.T) {
	retryCountHeaders := make([]string, 0)
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					retryCountHeaders = append(retryCountHeaders, req.Header.Get("X-Stainless-Retry-Count"))
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
		option.WithHeader("X-Stainless-Retry-Count", "42"),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}

	expectedRetryCountHeaders := []string{"42", "42", "42"}
	if !reflect.DeepEqual(retryCountHeaders, expectedRetryCountHeaders) {
		t.Errorf("Expected %v retry count headers, got %v", expectedRetryCountHeaders, retryCountHeaders)
	}
}

func TestRetryAfterMs(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After-Ms"): []string{"100"},
						},
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}
	if want := 3; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestNonJSONErrorBody(t *testing.T) {
	for _, body := range []string{"length limit exceeded", "too large", ""} {
		client := anthropic.NewClient(
			option.WithAPIKey("my-anthropic-api-key"),
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusRequestEntityTooLarge,
							Header: http.Header{
								http.CanonicalHeaderKey("Content-Type"): []string{"text/plain"},
							},
							Body: io.NopCloser(strings.NewReader(body)),
						}, nil
					},
				},
			}),
		)
		_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		})
		var apiErr *anthropic.Error
		if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusRequestEntityTooLarge {
			t.Fatalf("Expected a 413 API error, got %v", err)
		}
		if raw := apiErr.RawJSON(); raw != body {
			t.Errorf("Expected raw body %q, got %q", body, raw)
		}
		if want := strings.TrimSpace("413 Request Entity Too Large " + body); !strings.HasSuffix(err.Error(), want) {
			t.Errorf("Expected error %q to end with %q", err.Error(), want)
		}
	}
}

func TestErrorDumpRedactsSensitiveHeaders(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Header: http.Header{
							http.CanonicalHeaderKey("Set-Cookie"): []string{"session=dump-secret-cookie"},
						},
						Body: io.NopCloser(strings.NewReader("dump-body")),
					}, nil
				},
			},
		}),
		option.WithHeader("Authorization", "Bearer dump-secret-token"),
		option.WithHeader("X-Api-Key", "dump-secret-key"),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	var apiErr *anthropic.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected an API error, got %v", err)
	}
	for _, dump := range []string{string(apiErr.DumpRequest(true)), string(apiErr.DumpResponse(true))} {
		if strings.Contains(dump, "dump-secret") || !strings.Contains(dump, "***") {
			t.Errorf("Expected sensitive headers to be redacted, got:\n%s", dump)
		}
	}
	if dump := string(apiErr.DumpResponse(true)); !strings.Contains(dump, "dump-body") {
		t.Errorf("Expected a second response dump to include the body, got:\n%s", dump)
	}
	if got := apiErr.Request.Header.Get("X-Api-Key"); got != "dump-secret-key" {
		t.Errorf("Expected the request's X-Api-Key header to be unchanged, got %q", got)
	}
}

func TestContextCancel(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.Messages.New(cancelCtx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("Expected there to be a cancel error")
	}
}

func TestContextCancelDelay(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	_, err := client.Messages.New(cancelCtx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{
					Text: "x",
				},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaudeSonnet5_5,
	})
	if err == nil {
		t.Error("expected there to be a cancel error")
	}
}

func TestRequestTimeoutRetried(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	_, err := client.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithRequestTimeout((10 * time.Millisecond)),
		option.WithMaxRetries(1),
	)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded after exhausting retries, got: %v", err)
	}
	if want := 2; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestRequestTimeoutNotRetriedAfterContextDone(t *testing.T) {
	attempts := 0
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					cancel()
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	_, err := client.Messages.New(
		cancelCtx,
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithRequestTimeout(time.Second),
	)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	if want := 1; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestResponseBodyReadErrorRetried(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(iotest.ErrReader(io.ErrUnexpectedEOF)),
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithMaxRetries(1),
	)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("expected io.ErrUnexpectedEOF after exhausting retries, got: %v", err)
	}
	if want := 2; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestRetriedResponseBodyNotRead(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					if attempts == 1 {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Header: http.Header{
								http.CanonicalHeaderKey("Retry-After-Ms"): []string{"1"},
							},
							Body: unreadBody{t},
						}, nil
					}
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       http.NoBody,
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithMaxRetries(1),
	)
	if err == nil {
		t.Error("expected the error of the last response")
	}
	if want := 2; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestResponseBodyReadErrorKeepsResponse(t *testing.T) {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(iotest.ErrReader(io.ErrUnexpectedEOF)),
					}, nil
				},
			},
		}),
	)
	var res *http.Response
	_, err := client.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithMaxRetries(1),
		option.WithResponseInto(&res),
	)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("expected io.ErrUnexpectedEOF, got: %v", err)
	}
	if res == nil || res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected the response that was received, got: %v", res)
	}
}

func TestResponseBodyReadErrorNotRetriedWhenToldNotTo(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Header: http.Header{
							http.CanonicalHeaderKey("X-Should-Retry"): []string{"false"},
						},
						Body: io.NopCloser(iotest.ErrReader(io.ErrUnexpectedEOF)),
					}, nil
				},
			},
		}),
	)
	_, err := client.Messages.New(
		context.Background(),
		anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		},
		option.WithMaxRetries(2),
	)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("expected io.ErrUnexpectedEOF, got: %v", err)
	}
	if want := 1; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

// doerFunc is an HTTP client that is not an *http.Client, which a request keeps in a
// field of its own.
type doerFunc func(req *http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

// The request for the next page of a list is a clone of the request for the page before it.
func TestClonedRequestKeepsCustomHTTPClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request did not go through the custom HTTP client: %s", r.URL)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	requests := 0
	cfg, err := requestconfig.NewRequestConfig(
		context.Background(),
		http.MethodGet,
		"cards",
		nil,
		nil,
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(0),
		option.WithHTTPClient(doerFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: req}, nil
		})),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, cfg := range []*requestconfig.RequestConfig{cfg, cfg.Clone(context.Background())} {
		if err := cfg.Execute(); err != nil {
			t.Fatal(err)
		}
	}
	if requests != 2 {
		t.Errorf("the custom HTTP client got %d requests, want 2", requests)
	}
}

// fakeTracerProvider makes spans that keep their status and whether they ended. The
// method of the provider, of its tracer or of its spans that panicIn names panics.
type fakeTracerProvider struct {
	embedded.TracerProvider
	panicIn string
	// sampling is empty for spans that record, "drop" for spans that do not and
	// have a valid context, as a sampler that drops them makes them, and "invalid"
	// for spans that do not record and have no valid context.
	sampling string
	spans    []*fakeSpan
}

func (p *fakeTracerProvider) fail(method string) {
	if p.panicIn == method {
		panic(method + " failed")
	}
}

func (p *fakeTracerProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	p.fail("Tracer")
	return fakeTracer{provider: p}
}

type fakeTracer struct {
	embedded.Tracer
	provider *fakeTracerProvider
}

func (t fakeTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	t.provider.fail("Start")
	s := &fakeSpan{provider: t.provider}
	t.provider.spans = append(t.provider.spans, s)
	return trace.ContextWithSpan(ctx, s), s
}

type fakeSpan struct {
	noop.Span
	provider *fakeTracerProvider
	status   codes.Code
	ended    bool
}

func (s *fakeSpan) IsRecording() bool { return !s.ended && s.provider.sampling == "" }

func (s *fakeSpan) SpanContext() trace.SpanContext {
	config := trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{1}, TraceFlags: trace.FlagsSampled}
	switch s.provider.sampling {
	case "drop":
		config.TraceFlags = 0
	case "invalid":
		config = trace.SpanContextConfig{}
	}
	return trace.NewSpanContext(config)
}

func (s *fakeSpan) SetAttributes(...attribute.KeyValue) { s.provider.fail("SetAttributes") }

func (s *fakeSpan) SetStatus(code codes.Code, _ string) {
	s.provider.fail("SetStatus")
	s.status = code
}

func (s *fakeSpan) End(...trace.SpanEndOption) { s.ended = true }

func (s *fakeSpan) TracerProvider() trace.TracerProvider { return s.provider }

// fakePropagator counts how often the SDK has it write a trace context: traceparent
// for a span, and baggage with or without one.
type fakePropagator struct {
	injects int
	panics  bool
}

func (p *fakePropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	p.injects++
	if p.panics {
		panic("Inject failed")
	}
	propagation.TraceContext{}.Inject(ctx, carrier)
	carrier.Set("baggage", "tenant=acme")
}

func (p *fakePropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}

func (p *fakePropagator) Fields() []string { return []string{"traceparent", "baggage"} }

func usePropagator(t *testing.T, p *fakePropagator) {
	previous := otel.GetTextMapPropagator()
	t.Cleanup(func() { otel.SetTextMapPropagator(previous) })
	otel.SetTextMapPropagator(p)
}

// tracedAPI serves a card and the events of a card, after it has failed as many requests as
// failures says, and keeps the trace-context headers of each request.
type tracedAPI struct {
	failures     int
	traceHeaders [][]string
}

func (a *tracedAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.traceHeaders = append(a.traceHeaders, []string{r.Header.Get("traceparent"), r.Header.Get("baggage")})
	if a.failures > 0 {
		a.failures--
		w.Header().Set("Retry-After-Ms", "1")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if r.URL.Path == "/cards/card_1/events" {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: card\ndata: {\"message\":\"a\"}\n\nevent: card\ndata: {\"message\":\"b\"}\n\n")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"token\":\"card_1\"}")
}

func newTracedAPIClient(t *testing.T, a *tracedAPI, opts ...option.RequestOption) anthropic.Client {
	t.Setenv("ANTHROPIC_OPEN_TELEMETRY", "")
	server := httptest.NewServer(a)
	t.Cleanup(server.Close)
	return anthropic.NewClient(append([]option.RequestOption{
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(1),
	}, opts...)...)
}

// getCard gets the card the way an SDK method does, which names the span of its call.
func getCard(client anthropic.Client, opts ...option.RequestOption) error {
	var card struct {
		Token string `json:"token"`
	}
	err := client.Get(context.Background(), "cards/card_1", nil, &card, append(opts, requestconfig.WithSpan("cards.retrieve", "/cards/{card_id}"))...)
	if err == nil && card.Token != "card_1" {
		err = fmt.Errorf("the call returned the card %q", card.Token)
	}
	return err
}

func TestTraceContextPropagation(t *testing.T) {
	const withSpan = "00-01000000000000000000000000000000-0100000000000000-01"
	tests := []struct {
		name string
		// env is the value of the variable, which is as good as unset when it is empty.
		env        string
		opts       []option.RequestOption
		noProvider bool
		want       bool
	}{
		{name: "on by default", want: true},
		{name: "off by the option", opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false)}},
		{name: "off by the variable in any case", env: "FALSE"},
		{name: "only false turns it off", env: "0", want: true},
		{name: "the option wins over the variable", env: "false", opts: []option.RequestOption{option.WithOpenTelemetryPropagation(true)}, want: true},
		{name: "on without a tracer provider", noProvider: true, want: true},
		{name: "off without a tracer provider", noProvider: true, opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false)}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_OPEN_TELEMETRY_PROPAGATION", test.env)
			headers := &fakePropagator{}
			usePropagator(t, headers)
			spans := &fakeTracerProvider{}
			if !test.noProvider {
				test.opts = append(test.opts, option.WithTracerProvider(spans))
			}
			// The first attempt fails, so that the call makes a second one.
			server := &tracedAPI{failures: 1}
			client := newTracedAPIClient(t, server, test.opts...)

			if err := getCard(client); err != nil {
				t.Fatal(err)
			}

			want := []string{"", ""}
			if test.want {
				want[1] = "tenant=acme"
				if !test.noProvider {
					want[0] = withSpan
				}
			}
			if !slices.EqualFunc(server.traceHeaders, [][]string{want, want}, slices.Equal[[]string]) {
				t.Errorf("the attempts had the trace-context headers %q, want %q twice", server.traceHeaders, want)
			}
			if !test.want && headers.injects != 0 {
				t.Errorf("the SDK called the propagator %d times, want 0", headers.injects)
			}
			if test.noProvider {
				return
			}
			if len(spans.spans) != 1 || !spans.spans[0].ended || spans.spans[0].status != codes.Ok {
				t.Errorf("the call has the spans %+v, want one that ended with status Ok", spans.spans)
			}
		})
	}
}

// spanTransport keeps the span that the context of each request holds, which is the
// one that the application's HTTP instrumentation parents its spans to.
type spanTransport struct{ spans []trace.Span }

func (t *spanTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.spans = append(t.spans, trace.SpanFromContext(req.Context()))
	return http.DefaultTransport.RoundTrip(req)
}

// With the propagation off, the context of a span that a sampler dropped still
// reaches the application's HTTP instrumentation, in the context of each attempt.
func TestTraceContextPropagationOffKeepsDroppedSpan(t *testing.T) {
	tests := []struct {
		name        string
		sampling    string
		propagation bool
		// wantSpan is whether the context of each attempt holds the span of the call.
		wantSpan        bool
		wantTraceparent string
	}{
		{name: "a dropped span is in the context and not in the headers", sampling: "drop", wantSpan: true},
		{name: "a span without a valid context is in neither", sampling: "invalid"},
		{name: "with the propagation on a dropped span is in both", sampling: "drop", propagation: true, wantSpan: true, wantTraceparent: "00-01000000000000000000000000000000-0100000000000000-00"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			headers := &fakePropagator{}
			usePropagator(t, headers)
			spans := &fakeTracerProvider{sampling: test.sampling}
			seen := &spanTransport{}
			// The first attempt fails, so that the call makes a second one.
			server := &tracedAPI{failures: 1}
			client := newTracedAPIClient(t, server,
				option.WithTracerProvider(spans),
				option.WithOpenTelemetryPropagation(test.propagation),
				option.WithHTTPClient(&http.Client{Transport: seen}),
			)

			if err := getCard(client); err != nil {
				t.Fatal(err)
			}

			if len(spans.spans) != 1 || len(seen.spans) != 2 || len(server.traceHeaders) != 2 {
				t.Fatalf("the call started %d spans and made %d attempts, want 1 and 2", len(spans.spans), len(seen.spans))
			}
			for attempt, inContext := range seen.spans {
				if got := inContext == trace.Span(spans.spans[0]); got != test.wantSpan {
					t.Errorf("attempt %d: the context holds the span of the call: %t, want %t", attempt, got, test.wantSpan)
				}
				if got := inContext.SpanContext(); test.wantSpan && (!got.IsValid() || got.IsSampled()) {
					t.Errorf("attempt %d: the context holds the span context %+v, want a valid one that is not sampled", attempt, got)
				}
				if got := server.traceHeaders[attempt][0]; got != test.wantTraceparent {
					t.Errorf("attempt %d: traceparent is %q, want %q", attempt, got, test.wantTraceparent)
				}
			}
			if !test.propagation && headers.injects != 0 {
				t.Errorf("the SDK called the propagator %d times, want 0", headers.injects)
			}
		})
	}
}

func TestTelemetryPanicsStayOutOfTheCall(t *testing.T) {
	handled := 0
	previous := otel.GetErrorHandler()
	t.Cleanup(func() { otel.SetErrorHandler(previous) })
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) { handled++ }))
	// check runs call, which makes one API call that has spans as their provider p
	// makes them, and fails the test unless the panic reached the error handler
	// only, and every span that started ended.
	check := func(t *testing.T, p *fakeTracerProvider, wantSpans int, call func() error) {
		t.Helper()
		before := handled
		if err := call(); err != nil {
			t.Errorf("the call failed: %v", err)
		}
		if handled == before {
			t.Error("the OpenTelemetry error handler did not get the panic")
		}
		if len(p.spans) != wantSpans {
			t.Errorf("the call started %d spans, want %d", len(p.spans), wantSpans)
		}
		for _, s := range p.spans {
			if !s.ended {
				t.Error("the span of the call did not end")
			}
		}
	}

	// What the application gives the SDK: a tracer provider and a propagator.
	for method, wantSpans := range map[string]int{"Tracer": 0, "Start": 0, "SetAttributes": 1, "SetStatus": 1, "Inject": 1} {
		t.Run(method+" of the application", func(t *testing.T) {
			t.Setenv("ANTHROPIC_OPEN_TELEMETRY_PROPAGATION", "")
			usePropagator(t, &fakePropagator{panics: method == "Inject"})
			p := &fakeTracerProvider{panicIn: method}
			client := newTracedAPIClient(t, &tracedAPI{}, option.WithTracerProvider(p))
			check(t, p, wantSpans, func() error { return getCard(client) })
		})
	}

	// What an SDK method records of its own.
	failed := func(hook string) { panic(hook + " failed") }
	recorders := map[string]requestconfig.SpanRecorder{
		"StartAttributes": {StartAttributes: func(http.Header) []attribute.KeyValue { failed("StartAttributes"); return nil }},
		"RecordRequest":   {RecordRequest: func(*tracing.APICall) { failed("RecordRequest") }},
		"RecordResult":    {RecordResult: func(*tracing.APICall, error) { failed("RecordResult") }},
	}
	for hook, recorder := range recorders {
		t.Run(hook+" of the method", func(t *testing.T) {
			p := &fakeTracerProvider{}
			client := newTracedAPIClient(t, &tracedAPI{}, option.WithTracerProvider(p))
			check(t, p, 1, func() error { return getCard(client, requestconfig.WithSpanRecorder("cards.retrieve", recorder)) })
			if len(p.spans) == 1 && p.spans[0].status != codes.Ok {
				t.Errorf("the span has the status %v, want Ok", p.spans[0].status)
			}
		})
	}
	observers := map[string]func(*tracing.APICall) tracing.StreamObserver{
		"StreamObserver": func(*tracing.APICall) tracing.StreamObserver {
			failed("StreamObserver")
			return tracing.StreamObserver{}
		},
		"OnEvent": func(*tracing.APICall) tracing.StreamObserver {
			return tracing.StreamObserver{OnEvent: func(ssestream.Event) { failed("OnEvent") }}
		},
		"OnDecoded": func(*tracing.APICall) tracing.StreamObserver {
			return tracing.StreamObserver{OnDecoded: func(ssestream.Event, any) { failed("OnDecoded") }}
		},
		// The stream ends without its last event, so it is incomplete.
		"OnIncomplete": func(*tracing.APICall) tracing.StreamObserver {
			return tracing.StreamObserver{LastEvent: "never sent", OnIncomplete: func() { failed("OnIncomplete") }}
		},
	}
	for hook, observer := range observers {
		t.Run(hook+" of the method", func(t *testing.T) {
			p := &fakeTracerProvider{}
			client := newTracedAPIClient(t, &tracedAPI{}, option.WithTracerProvider(p))
			recorder := requestconfig.SpanRecorder{StreamObserver: observer}
			check(t, p, 1, func() error {
				// A streaming method of the SDK makes its stream from the raw response.
				var raw *http.Response
				err := client.Get(context.Background(), "cards/card_1/events", nil, &raw, requestconfig.WithSpanRecorder("cards.stream", recorder))
				stream := ssestream.NewStream[struct {
					Message string `json:"message"`
				}](ssestream.NewDecoder(raw), err)
				var messages []string
				for stream.Next() {
					messages = append(messages, stream.Current().Message)
				}
				if err := stream.Err(); err != nil {
					return err
				}
				if !slices.Equal(messages, []string{"a", "b"}) {
					return fmt.Errorf("the stream had the messages %q", messages)
				}
				return nil
			})
			wantStatus := codes.Unset
			if hook == "OnIncomplete" {
				wantStatus = codes.Error
			}
			if len(p.spans) == 1 && p.spans[0].status != wantStatus {
				t.Errorf("the span has the status %v, want %v", p.spans[0].status, wantStatus)
			}
		})
	}
}

func TestContextDeadline(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client := anthropic.NewClient(
			option.WithAPIKey("my-anthropic-api-key"),
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						<-req.Context().Done()
						return nil, req.Context().Err()
					},
				},
			}),
		)
		_, err := client.Messages.New(deadlineCtx, anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeSonnet5_5,
		})
		if err == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

func TestContextDeadlineStreaming(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client := anthropic.NewClient(
			option.WithAPIKey("my-anthropic-api-key"),
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Status:     "200 OK",
							Body: io.NopCloser(
								io.Reader(readerFunc(func([]byte) (int, error) {
									<-req.Context().Done()
									return 0, req.Context().Err()
								})),
							),
						}, nil
					},
				},
			}),
		)
		stream := client.Messages.NewStreaming(deadlineCtx, anthropic.MessageNewParams{
			MaxTokens: 1024,
			Messages: []anthropic.MessageParam{{
				Content: []anthropic.ContentBlockParamUnion{{
					OfText: &anthropic.TextBlockParam{
						Text: "x",
					},
				}},
				Role: anthropic.MessageParamRoleUser,
			}},
			Model: anthropic.ModelClaudeOpus5_5,
		})
		for stream.Next() {
			_ = stream.Current()
		}
		if stream.Err() == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

func TestContextDeadlineStreamingWithRequestTimeout(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})
	deadline := time.Now().Add(100 * time.Millisecond)

	go func() {
		client := anthropic.NewClient(
			option.WithAPIKey("my-anthropic-api-key"),
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Status:     "200 OK",
							Body: io.NopCloser(
								io.Reader(readerFunc(func([]byte) (int, error) {
									<-req.Context().Done()
									return 0, req.Context().Err()
								})),
							),
						}, nil
					},
				},
			}),
		)
		stream := client.Messages.NewStreaming(
			context.Background(),
			anthropic.MessageNewParams{
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{{
					Content: []anthropic.ContentBlockParamUnion{{
						OfText: &anthropic.TextBlockParam{
							Text: "x",
						},
					}},
					Role: anthropic.MessageParamRoleUser,
				}},
				Model: anthropic.ModelClaudeOpus5_5,
			},
			option.WithRequestTimeout((100 * time.Millisecond)),
		)
		for stream.Next() {
			_ = stream.Current()
		}
		if stream.Err() == nil {
			t.Error("expected there to be a deadline error")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}

type readerFunc func([]byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }
func (f readerFunc) Close() error               { return nil }

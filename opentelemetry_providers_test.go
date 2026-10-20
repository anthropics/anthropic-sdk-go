package anthropic_test

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicaws "github.com/anthropics/anthropic-sdk-go/aws"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
	"github.com/anthropics/anthropic-sdk-go/foundry"
	"github.com/anthropics/anthropic-sdk-go/googlecloud"
	"github.com/anthropics/anthropic-sdk-go/internal"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/vertex"
)

func TestOpenTelemetryBedrockProviderName(t *testing.T) {
	client, recorder := newTracedClient(t, toolRunnerTransport(),
		bedrock.WithConfig(aws.Config{
			Region:      "us-east-1",
			Credentials: credentials.NewStaticCredentialsProvider("id", "secret", ""),
		}))

	// The tool runner's own spans name the client's platform, as its requests do.
	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
	if _, err := runner.RunToCompletion(context.Background()); err != nil {
		t.Fatal(err)
	}
	kinds := map[trace.SpanKind]int{}
	for _, s := range recorder.Spans() {
		kinds[s.SpanKind()]++
		wantAttrs(t, s, map[string]any{"gen_ai.provider.name": "aws.bedrock"})
		if s.SpanKind() == trace.SpanKindClient {
			wantAttrs(t, s, map[string]any{
				"server.address": "bedrock-runtime.us-east-1.amazonaws.com",
				"server.port":    int64(443),
				// The route is that of the SDK's method, and the URL the one that Bedrock gets.
				"http.request.method": "POST",
				"http.route":          "/v1/messages",
				"url.scheme":          "https",
				"url.path":            "/model/claude-requested/invoke",
				"url.query":           nil,
				"url.full":            "https://bedrock-runtime.us-east-1.amazonaws.com/model/claude-requested/invoke",
			})
		} else {
			// Only a request has a server.
			wantAttrs(t, s, map[string]any{"server.address": nil, "server.port": nil})
		}
	}
	if kinds[trace.SpanKindInternal] != 2 || kinds[trace.SpanKindClient] != 2 {
		t.Errorf("got spans by kind %v, want two INTERNAL and two CLIENT", kinds)
	}
	// The platform's host is the client's base URL, so each span starts with it.
	starts := recorder.sampler.started["anthropic.messages.create"]
	for _, attrs := range starts {
		if attrs["server.address"] != "bedrock-runtime.us-east-1.amazonaws.com" || attrs["server.port"] != int64(443) {
			t.Errorf("a span started with server %v:%v, want bedrock-runtime.us-east-1.amazonaws.com:443", attrs["server.address"], attrs["server.port"])
		}
		if attrs["http.request.method"] != "POST" || attrs["http.route"] != "/v1/messages" {
			t.Errorf("a span started with %v %v, want POST /v1/messages", attrs["http.request.method"], attrs["http.route"])
		}
	}
	if len(starts) != 2 {
		t.Errorf("%d spans started, want 2", len(starts))
	}
}

func TestOpenTelemetryVertexProviderName(t *testing.T) {
	creds := &google.Credentials{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"})}
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON),
		vertex.WithCredentials(context.Background(), "us-east5", "project", creds))

	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"gen_ai.provider.name": "gcp.vertex_ai",
		"gen_ai.request.model": "claude-requested",
		"server.address":       "us-east5-aiplatform.googleapis.com",
		"server.port":          int64(443),
		// The route is that of the SDK's method, and the URL the one that Vertex AI gets.
		"http.request.method": "POST",
		"http.route":          "/v1/messages",
		"url.scheme":          "https",
		"url.path":            "/v1/projects/project/locations/us-east5/publishers/anthropic/models/claude-requested:rawPredict",
		"url.query":           nil,
		"url.full":            "https://us-east5-aiplatform.googleapis.com/v1/projects/project/locations/us-east5/publishers/anthropic/models/claude-requested:rawPredict",
	})
	started := recorder.sampler.last("anthropic.messages.create")
	if started["server.address"] != "us-east5-aiplatform.googleapis.com" || started["server.port"] != int64(443) {
		t.Errorf("the span started with server %v:%v, want us-east5-aiplatform.googleapis.com:443", started["server.address"], started["server.port"])
	}
	if started["http.request.method"] != "POST" || started["http.route"] != "/v1/messages" {
		t.Errorf("the span started with %v %v, want POST /v1/messages", started["http.request.method"], started["http.route"])
	}
}

// TestOpenTelemetryTraceHeadersAreSentUnsigned checks every client that signs
// its requests for AWS. An HTTP client instrumentation may replace the trace
// context headers with those of its own span after the request is signed, and a
// signature that covered them would then no longer match.
func TestOpenTelemetryTraceHeadersAreSentUnsigned(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	ctx := tracedContext(t)

	for name, send := range map[string]func(t *testing.T, opts ...option.RequestOption) error{
		"bedrock": func(t *testing.T, opts ...option.RequestOption) error {
			client := anthropic.NewClient(append(opts, option.WithAPIKey("key"), bedrock.WithConfig(aws.Config{
				Region:      "us-east-1",
				Credentials: credentials.NewStaticCredentialsProvider("id", "secret", ""),
			}))...)
			_, err := client.Messages.New(ctx, messageParams())
			return err
		},
		"aws": func(t *testing.T, opts ...option.RequestOption) error {
			client, err := anthropicaws.NewClient(context.Background(), anthropicaws.ClientConfig{
				AWSAccessKey: "id", AWSSecretAccessKey: "secret", AWSRegion: "us-east-1", WorkspaceID: "wrkspc_1",
			}, opts...)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.Messages.New(ctx, messageParams())
			return err
		},
		"bedrock mantle": func(t *testing.T, opts ...option.RequestOption) error {
			client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{
				AWSAccessKey: "id", AWSSecretAccessKey: "secret", AWSRegion: "us-east-1",
			}, opts...)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.Messages.New(ctx, messageParams())
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			var sent http.Header
			transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				sent = req.Header.Clone()
				return respond(200, "application/json", messageJSON).fn(req)
			}}
			recorder := newSpanRecorder()
			if err := send(t, option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMaxRetries(0), option.WithTracerProvider(recorder)); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)

			authorization := sent.Get("Authorization")
			_, signed, found := strings.Cut(authorization, "SignedHeaders=")
			if !strings.HasPrefix(authorization, "AWS4-HMAC-SHA256 ") || !found {
				t.Fatalf("Authorization = %q, want a signature that lists its signed headers", authorization)
			}
			signed, _, _ = strings.Cut(signed, ",")
			signedHeaders := strings.Split(signed, ";")
			if !slices.Contains(signedHeaders, "host") {
				t.Errorf("SignedHeaders = %q, want the host among them", signed)
			}
			for _, header := range []string{"traceparent", "tracestate", "baggage"} {
				if sent.Get(header) == "" {
					t.Errorf("the request has no %s header", header)
				}
				if slices.Contains(signedHeaders, header) {
					t.Errorf("SignedHeaders = %q: the signature covers %s", signed, header)
				}
			}
			if got, want := sent.Get("traceparent"), "00-"+span.SpanContext().TraceID().String()+"-"+span.SpanContext().SpanID().String()+"-01"; got != want {
				t.Errorf("traceparent = %q, want the span's %q", got, want)
			}
			if got := sent.Get("tracestate"); got != "vendor=value" {
				t.Errorf("tracestate = %q, want vendor=value", got)
			}
			if got := sent.Get("baggage"); got != "tenant=acme" {
				t.Errorf("baggage = %q, want tenant=acme", got)
			}

			if err := send(t, option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMaxRetries(0), option.WithTracerProvider(recorder), option.WithOpenTelemetryPropagation(false)); err != nil {
				t.Fatal(err)
			}
			wantTraceContextHeaders(t, sent, false)
			_, signed, _ = strings.Cut(sent.Get("Authorization"), "SignedHeaders=")
			signed, _, _ = strings.Cut(signed, ",")
			if !strings.HasPrefix(sent.Get("Authorization"), "AWS4-HMAC-SHA256 ") || !slices.Contains(strings.Split(signed, ";"), "host") {
				t.Errorf("Authorization = %q without propagation, want a signature that covers the host", sent.Get("Authorization"))
			}
			if got := len(recorder.Spans()); got != 2 {
				t.Errorf("got %d ended spans, want one for each call", got)
			}
		})
	}
}

// vendorPropagator sends the span's trace ID in a header of its own.
type vendorPropagator struct{}

func (vendorPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	carrier.Set("X-Vendor-Trace", trace.SpanContextFromContext(ctx).TraceID().String())
}
func (vendorPropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}
func (vendorPropagator) Fields() []string { return []string{"X-Vendor-Trace"} }

func TestOpenTelemetryPropagatorFieldsAreSentUnsigned(t *testing.T) {
	setPropagator(t, vendorPropagator{})
	var sent http.Header
	transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		sent = req.Header.Clone()
		return respond(200, "application/json", messageJSON).fn(req)
	}}
	client, recorder := newTracedClient(t, transport, bedrock.WithConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider("id", "secret", ""),
	}))
	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	if got, want := sent.Get("X-Vendor-Trace"), onlySpan(t, recorder).SpanContext().TraceID().String(); got != want {
		t.Errorf("X-Vendor-Trace = %q, want the span's trace ID %q", got, want)
	}
	_, signed, _ := strings.Cut(sent.Get("Authorization"), "SignedHeaders=")
	signed, _, _ = strings.Cut(signed, ",")
	if signed == "" || slices.Contains(strings.Split(signed, ";"), "x-vendor-trace") {
		t.Errorf("SignedHeaders = %q, want a signature that does not cover x-vendor-trace", signed)
	}
}

func TestOpenTelemetryMantleProviderName(t *testing.T) {
	recorder := newSpanRecorder()
	client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{APIKey: "key", AWSRegion: "us-east-1"},
		option.WithHTTPClient(&http.Client{Transport: respond(200, "application/json", messageJSON)}),
		option.WithTracerProvider(recorder))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	// Mantle serves the paths of the API, so the URL has the path of the route.
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"gen_ai.provider.name": "aws.bedrock",
		"http.route":           "/v1/messages",
		"url.path":             "/anthropic/v1/messages",
		"url.full":             "https://bedrock-mantle.us-east-1.api.aws/anthropic/v1/messages",
	})
}

// The client of each platform names its package in the instrumentation scope of
// its spans. The route is that of the SDK's method on every platform, and the
// URL is the one that the platform gets.
func TestOpenTelemetryPlatformClientsHaveAScopeAndAURLOfTheirOwn(t *testing.T) {
	ctx := context.Background()
	awsConfig := aws.Config{Region: "us-east-1", Credentials: credentials.NewStaticCredentialsProvider("id", "secret", "")}
	cases := map[string]struct {
		scope, url string
		send       func(t *testing.T, opts ...option.RequestOption) error
	}{
		"first party": {
			scope: "com.anthropic.sdk.go",
			url:   "https://api.anthropic.com/v1/messages",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client := anthropic.NewClient(option.WithoutEnvironmentDefaults(), option.WithAPIKey("key"))
				_, err := client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"bedrock": {
			scope: "com.anthropic.sdk.go.bedrock",
			url:   "https://bedrock-runtime.us-east-1.amazonaws.com/model/claude-requested/invoke",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client := anthropic.NewClient(bedrock.WithConfig(awsConfig))
				_, err := client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"bedrock mantle": {
			scope: "com.anthropic.sdk.go.bedrock",
			url:   "https://bedrock-mantle.us-east-1.api.aws/anthropic/v1/messages",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client, err := bedrock.NewMantleClient(ctx, bedrock.MantleClientConfig{APIKey: "key", AWSRegion: "us-east-1"})
				if err != nil {
					t.Fatal(err)
				}
				_, err = client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"aws": {
			scope: "com.anthropic.sdk.go.aws",
			url:   "https://aws-external-anthropic.us-east-1.api.aws/v1/messages",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client, err := anthropicaws.NewClient(ctx, anthropicaws.ClientConfig{APIKey: "key", AWSRegion: "us-east-1", WorkspaceID: "wrkspc_1"})
				if err != nil {
					t.Fatal(err)
				}
				_, err = client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"vertex": {
			scope: "com.anthropic.sdk.go.vertex",
			url:   "https://us-east5-aiplatform.googleapis.com/v1/projects/project/locations/us-east5/publishers/anthropic/models/claude-requested:rawPredict",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				creds := &google.Credentials{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"})}
				client := anthropic.NewClient(vertex.WithCredentials(ctx, "us-east5", "project", creds))
				_, err := client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"foundry": {
			scope: "com.anthropic.sdk.go.foundry",
			url:   "https://resource.services.ai.azure.com/anthropic/v1/messages",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client, err := foundry.NewClient(foundry.ClientConfig{APIKey: "key", Resource: "resource"})
				if err != nil {
					t.Fatal(err)
				}
				_, err = client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
		"googlecloud": {
			scope: "com.anthropic.sdk.go.googlecloud",
			url:   "https://claude.googleapis.com/v1/projects/project/locations/global/workspaces/wrkspc_1/invoke/v1/messages",
			send: func(t *testing.T, opts ...option.RequestOption) error {
				client, err := googlecloud.NewClient(ctx, googlecloud.ClientConfig{
					Project: "project", WorkspaceID: "wrkspc_1", TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "token"}),
				})
				if err != nil {
					t.Fatal(err)
				}
				_, err = client.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			var sent string
			transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				sent = req.URL.String()
				return respond(200, "application/json", messageJSON).fn(req)
			}}
			recorder := newSpanRecorder()
			if err := test.send(t, option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMaxRetries(0), option.WithTracerProvider(recorder)); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			if scope := span.InstrumentationScope(); scope.Name != test.scope || scope.Version != internal.PackageVersion {
				t.Errorf("scope %s %s, want %s %s", scope.Name, scope.Version, test.scope, internal.PackageVersion)
			}
			if sent != test.url {
				t.Errorf("the request went to %s, want %s", sent, test.url)
			}
			wantAttrs(t, span, map[string]any{"http.route": "/v1/messages", "url.full": test.url})
		})
	}
}

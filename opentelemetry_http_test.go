package anthropic_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// urlAttrs are the attributes that a span has from the URL of its first request.
var urlAttrs = []string{"url.scheme", "url.path", "url.query", "url.full"}

// absent returns want with each of keys absent.
func absent(want map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		want[key] = nil
	}
	return want
}

func TestOpenTelemetryRequestSpanHasTheMethodTheRouteAndTheURL(t *testing.T) {
	var sent string
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		sent = req.URL.String()
		return respond(200, "application/json", `{"id":"msgbatch_1","type":"message_batch"}`).fn(req)
	}}, option.WithBaseURL("https://user:secret@api.example.com:8443/proxy/"))

	if _, err := client.Beta.Messages.Batches.Get(context.Background(), "msgbatch_1", anthropic.BetaMessageBatchGetParams{}); err != nil {
		t.Fatal(err)
	}
	if want := "https://user:secret@api.example.com:8443/proxy/v1/messages/batches/msgbatch_1"; sent != want {
		t.Fatalf("the request went to %s, want %s", sent, want)
	}
	// The route is the endpoint's, whatever the base URL, and the URL is the one
	// that the request went to, without what the endpoint does not define.
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"http.request.method":        "GET",
		"http.route":                 "/v1/messages/batches/{message_batch_id}",
		"url.scheme":                 "https",
		"url.path":                   "/proxy/v1/messages/batches/msgbatch_1",
		"url.query":                  nil,
		"url.full":                   "https://api.example.com:8443/proxy/v1/messages/batches/msgbatch_1",
		"server.address":             "api.example.com",
		"server.port":                int64(8443),
		"http.request.resend_count":  nil,
		"network.peer.address":       nil,
		"network.peer.port":          nil,
		"network.protocol.version":   nil,
		"anthropic.message.batch.id": "msgbatch_1",
	})
}

func TestOpenTelemetryServerPortIsThatOfTheSchemeWhenTheBaseURLHasNone(t *testing.T) {
	cases := map[string]struct {
		port int64
		full string
	}{
		"http://api.example.com/":      {80, "http://api.example.com/v1/models/claude"},
		"https://api.example.com/":     {443, "https://api.example.com/v1/models/claude"},
		"https://api.example.com:443/": {443, "https://api.example.com:443/v1/models/claude"},
		"http://api.example.com:8080/": {8080, "http://api.example.com:8080/v1/models/claude"},
	}
	for baseURL, want := range cases {
		t.Run(baseURL, func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", `{"id":"claude","type":"model"}`), option.WithBaseURL(baseURL))
			if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err != nil {
				t.Fatal(err)
			}
			// url.full has a port only when the URL of the request gives one.
			wantAttrs(t, onlySpan(t, recorder), map[string]any{
				"server.address": "api.example.com",
				"server.port":    want.port,
				"url.full":       want.full,
			})
		})
	}
}

func TestOpenTelemetryURLQueryHasOnlyTheParametersThatTheEndpointDefines(t *testing.T) {
	ctx := context.Background()
	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	cases := map[string]struct {
		call func(client anthropic.Client, opts ...option.RequestOption) error
		opts []option.RequestOption
		// sent is the query string of the request, and query the url.query of the
		// span, or nil when the span has none.
		sent  string
		query any
	}{
		"an array, and a parameter that the application added": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Beta.Files.List(ctx, anthropic.BetaFileListParams{IDs: []string{"file_b", "file a"}, Limit: anthropic.Int(5), ScopeID: anthropic.String("sesn_1")}, opts...)
				return err
			},
			opts:  []option.RequestOption{option.WithQuery("trace", "on")},
			sent:  "ids%5B%5D=file_b&ids%5B%5D=file+a&limit=5&scope_id=sesn_1&trace=on",
			query: "ids%5B%5D=file_b&ids%5B%5D=file+a&limit=5&scope_id=sesn_1",
		},
		// scope_id is a parameter of the beta endpoint only.
		"a parameter that only the beta endpoint defines": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Files.List(ctx, anthropic.FileListParams{Limit: anthropic.Int(5)}, opts...)
				return err
			},
			opts:  []option.RequestOption{option.WithQuery("scope_id", "sesn_1")},
			sent:  "limit=5&scope_id=sesn_1",
			query: "limit=5",
		},
		// created_at[gte] is a name that the endpoint defines, and created_at is
		// none, so no other created_at[...] stays.
		"a name with brackets of its own": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Beta.Sessions.List(ctx, anthropic.BetaSessionListParams{CreatedAtGte: anthropic.Time(created), Statuses: []string{"idle"}}, opts...)
				return err
			},
			opts:  []option.RequestOption{option.WithQuery("created_at[eq]", "now")},
			sent:  "created_at%5Beq%5D=now&created_at%5Bgte%5D=2026-01-02T03%3A04%3A05Z&statuses%5B%5D=idle",
			query: "created_at%5Bgte%5D=2026-01-02T03%3A04%3A05Z&statuses%5B%5D=idle",
		},
		"a value of the application's for a parameter that the endpoint defines": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Models.List(ctx, anthropic.ModelListParams{Limit: anthropic.Int(5)}, opts...)
				return err
			},
			opts:  []option.RequestOption{option.WithQuery("limit", "9")},
			sent:  "limit=9",
			query: "limit=9",
		},
		"the lifecycle stages of the models list": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Models.List(ctx, anthropic.ModelListParams{Lifecycle: []string{"retired"}, Limit: anthropic.Int(5)}, opts...)
				return err
			},
			sent:  "lifecycle%5B%5D=retired&limit=5",
			query: "lifecycle%5B%5D=retired&limit=5",
		},
		"the lifecycle stages of the beta models list": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Beta.Models.List(ctx, anthropic.BetaModelListParams{Lifecycle: []string{"retired"}, Limit: anthropic.Int(5)}, opts...)
				return err
			},
			sent:  "lifecycle%5B%5D=retired&limit=5",
			query: "lifecycle%5B%5D=retired&limit=5",
		},
		"the default workspace of the workspaces list": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Organization.Workspaces.List(ctx, anthropic.OrganizationWorkspaceListParams{IncludeDefault: anthropic.Bool(true)}, opts...)
				return err
			},
			sent:  "include_default=true",
			query: "include_default=true",
		},
		"the default workspace of the beta workspaces list": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Beta.Organization.Workspaces.List(ctx, anthropic.BetaOrganizationWorkspaceListParams{IncludeDefault: anthropic.Bool(true)}, opts...)
				return err
			},
			sent:  "include_default=true",
			query: "include_default=true",
		},
		"the statuses of the session threads list": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Beta.Sessions.Threads.List(ctx, "sesn_1", anthropic.BetaSessionThreadListParams{Statuses: []anthropic.BetaManagedAgentsSessionThreadStatus{anthropic.BetaManagedAgentsSessionThreadStatusIdle}}, opts...)
				return err
			},
			sent:  "statuses%5B%5D=idle",
			query: "statuses%5B%5D=idle",
		},
		"nothing but a parameter that the application added": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Messages.New(ctx, messageParams(), opts...)
				return err
			},
			opts: []option.RequestOption{option.WithQuery("trace", "on")},
			sent: "trace=on",
		},
		"no query string": {
			call: func(c anthropic.Client, opts ...option.RequestOption) error {
				_, err := c.Messages.New(ctx, messageParams(), opts...)
				return err
			},
		},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			var sent, path string
			client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				sent, path = req.URL.RawQuery, req.URL.Path
				return respond(200, "application/json", `{"data":[],"has_more":false}`).fn(req)
			}})
			if err := test.call(client, test.opts...); err != nil {
				t.Fatal(err)
			}
			if sent != test.sent {
				t.Fatalf("the request had the query %q, want %q", sent, test.sent)
			}
			full := "http://localhost:4010" + path
			if test.query != nil {
				full += "?" + test.query.(string)
			}
			wantAttrs(t, onlySpan(t, recorder), map[string]any{"url.query": test.query, "url.full": full})
		})
	}
}

func TestOpenTelemetryURLQueryIsContentWhenAParameterOfItIs(t *testing.T) {
	// A parameter that filters like a search, or whose value comes from stored
	// content, makes url.query, and the url.full that has it, content.
	ctx := context.Background()
	cases := map[string]struct {
		call func(client anthropic.Client) error
		// query is the url.query of the span, and content whether it is content.
		query   string
		content bool
	}{
		"email": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Users.List(ctx, anthropic.BetaOrganizationUserListParams{Email: anthropic.String("a@example.com"), Limit: anthropic.Int(5)})
			return err
		}, "email=a%40example.com&limit=5", true},
		"path_prefix": {func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.Memories.List(ctx, "memstore_1", anthropic.BetaMemoryStoreMemoryListParams{PathPrefix: anthropic.String("/notes/"), Limit: anthropic.Int(5)})
			return err
		}, "limit=5&path_prefix=%2Fnotes%2F", true},
		"expected_content_sha256": {func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.Memories.Delete(ctx, "mem_1", anthropic.BetaMemoryStoreMemoryDeleteParams{MemoryStoreID: "memstore_1", ExpectedContentSha256: anthropic.String("abc123")})
			return err
		}, "expected_content_sha256=abc123", true},
		"limit": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Users.List(ctx, anthropic.BetaOrganizationUserListParams{Limit: anthropic.Int(5)})
			return err
		}, "limit=5", false},
		// filter matches names on the skills, plugins and connectors lists, so it
		// is content there. On the other lists it takes IDs or fixed values.
		"filter of analytics skills": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Skills.List(ctx, anthropic.BetaOrganizationAnalyticsSkillListParams{Filter: []string{"skill_name:Quarterly"}})
			return err
		}, "filter%5B%5D=skill_name%3AQuarterly", true},
		"filter of analytics plugins": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Plugins.List(ctx, anthropic.BetaOrganizationAnalyticsPluginListParams{Filter: []string{"plugin_name:Quarterly"}})
			return err
		}, "filter%5B%5D=plugin_name%3AQuarterly", true},
		"filter of analytics connectors": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Connectors.List(ctx, anthropic.BetaOrganizationAnalyticsConnectorListParams{Filter: []string{"connector_name:Quarterly"}})
			return err
		}, "filter%5B%5D=connector_name%3AQuarterly", true},
		"filter of analytics users": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Users.List(ctx, anthropic.BetaOrganizationAnalyticsUserListParams{Filter: []string{"user_id:user_1"}})
			return err
		}, "filter%5B%5D=user_id%3Auser_1", false},
		// The page cursor of these two lists encodes a value of the last row: an
		// email address, or a memory's path.
		"page of analytics users": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Users.List(ctx, anthropic.BetaOrganizationAnalyticsUserListParams{Page: anthropic.String("page_abc"), Limit: anthropic.Int(5)})
			return err
		}, "limit=5&page=page_abc", true},
		"page of memories": {func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.Memories.List(ctx, "memstore_1", anthropic.BetaMemoryStoreMemoryListParams{Page: anthropic.String("page_abc"), Limit: anthropic.Int(5)})
			return err
		}, "limit=5&page=page_abc", true},
		"page of analytics skills": {func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Analytics.Skills.List(ctx, anthropic.BetaOrganizationAnalyticsSkillListParams{Page: anthropic.String("page_abc"), Limit: anthropic.Int(5)})
			return err
		}, "limit=5&page=page_abc", false},
	}
	modes := []option.OpenTelemetryContentMode{option.OpenTelemetryContentModeMetadataOnly, option.OpenTelemetryContentModeContent}
	for name, test := range cases {
		for _, mode := range modes {
			t.Run(name+" "+string(mode), func(t *testing.T) {
				var path string
				client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
					path = req.URL.Path
					return respond(200, "application/json", `{"data":[],"has_more":false}`).fn(req)
				}}, option.WithOpenTelemetryTracesContentMode(mode))
				if err := test.call(client); err != nil {
					t.Fatal(err)
				}
				var query, full any = test.query, "http://localhost:4010" + path + "?" + test.query
				if test.content && mode == option.OpenTelemetryContentModeMetadataOnly {
					query, full = nil, nil
				}
				// The rest of the URL is metadata whatever the query has.
				wantAttrs(t, onlySpan(t, recorder), map[string]any{"url.query": query, "url.full": full, "url.scheme": "http", "url.path": path})
			})
		}
	}
}

func TestOpenTelemetryOrganizationMethodHasASpan(t *testing.T) {
	modes := []option.OpenTelemetryContentMode{option.OpenTelemetryContentModeMetadataOnly, option.OpenTelemetryContentModeContent}
	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", `{"data":[],"has_more":false}`), option.WithOpenTelemetryTracesContentMode(mode))
			if _, err := client.Organization.Users.List(context.Background(), anthropic.OrganizationUserListParams{Email: anthropic.String("a@example.com"), Limit: anthropic.Int(5)}); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			if span.Name() != "anthropic.organization.users.list" {
				t.Errorf("span %q, want anthropic.organization.users.list", span.Name())
			}
			// email filters like a search, so the query that has it is content.
			var query, full any
			if mode == option.OpenTelemetryContentModeContent {
				query, full = "email=a%40example.com&limit=5", "http://localhost:4010/v1/organizations/users?email=a%40example.com&limit=5"
			}
			wantAttrs(t, span, map[string]any{
				"http.request.method": "GET",
				"http.route":          "/v1/organizations/users",
				"url.path":            "/v1/organizations/users",
				"url.query":           query,
				"url.full":            full,
			})
		})
	}
}

func TestOpenTelemetryURLIsThatOfTheFirstRequest(t *testing.T) {
	overloaded := respond(529, "application/json", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`, "retry-after-ms", "1")
	cases := map[string]struct {
		failures int
		resent   any
	}{
		"one request":   {0, nil},
		"one retry":     {1, int64(1)},
		"two retries":   {2, int64(2)},
		"three retries": {3, int64(3)},
	}
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			var paths []string
			requests := 0
			// Each request after the first goes elsewhere.
			elsewhere := func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				requests++
				if requests > 1 {
					req.URL.Scheme, req.URL.Host, req.URL.Path, req.URL.RawQuery = "https", "later.example.com", "/later", "limit=2"
				}
				return next(req)
			}
			client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				paths = append(paths, req.URL.Path)
				if len(paths) <= test.failures {
					return overloaded.fn(req)
				}
				return respond(200, "application/json", `{"data":[],"has_more":false}`).fn(req)
			}}, option.WithMaxRetries(3), option.WithMiddleware(elsewhere))

			if _, err := client.Models.List(context.Background(), anthropic.ModelListParams{Limit: anthropic.Int(1)}); err != nil {
				t.Fatal(err)
			}
			last := "/v1/models"
			if test.failures > 0 {
				last = "/later"
			}
			if len(paths) != test.failures+1 || paths[len(paths)-1] != last {
				t.Fatalf("the requests went to %v, want %d of them, the later ones to /later", paths, test.failures+1)
			}
			wantAttrs(t, onlySpan(t, recorder), map[string]any{
				"url.scheme":                "http",
				"url.path":                  "/v1/models",
				"url.query":                 "limit=1",
				"url.full":                  "http://localhost:4010/v1/models?limit=1",
				"http.request.resend_count": test.resent,
				"http.response.status_code": int64(200),
			})
		})
	}
}

func TestOpenTelemetryResendCountOfACallThatFails(t *testing.T) {
	client, recorder := newTracedClient(t, respond(529, "application/json",
		`{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`, "retry-after-ms", "1"), option.WithMaxRetries(2))
	if _, err := client.Messages.New(context.Background(), messageParams()); err == nil {
		t.Fatal("the call did not fail")
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{"http.request.resend_count": int64(2), "error.type": "overloaded_error"})
}

func TestOpenTelemetryResendCountOfAStream(t *testing.T) {
	requests := 0
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		requests++
		if requests == 1 {
			return respond(529, "application/json", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`, "retry-after-ms", "1").fn(req)
		}
		return respond(200, "text/event-stream", streamBody).fn(req)
	}}, option.WithMaxRetries(1))
	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{"http.request.resend_count": int64(1), "url.path": "/v1/messages"})
}

func TestOpenTelemetryURLIsAbsentWhenNoRequestIsSent(t *testing.T) {
	failed := errors.New("no credentials")
	fails := func(*http.Request, option.MiddlewareNext) (*http.Response, error) { return nil, failed }
	client, recorder := newTracedClient(t, &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		t.Error("a request was sent")
		return nil, failed
	}}, option.WithMaxRetries(1), option.WithMiddleware(fails))

	if _, err := client.Models.List(context.Background(), anthropic.ModelListParams{Limit: anthropic.Int(1)}); !errors.Is(err, failed) {
		t.Fatalf("err = %v, want the failure of the middleware", err)
	}
	// The span started with what the method knows, and no request gave it a URL.
	wantAttrs(t, onlySpan(t, recorder), absent(map[string]any{
		"http.request.method": "GET",
		"http.route":          "/v1/models",
	}, append([]string{"http.request.resend_count"}, urlAttrs...)...))
}

func TestOpenTelemetryURLIsNotThatOfARedirect(t *testing.T) {
	transport := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/v1/models" {
			return respond(307, "text/plain", "", "Location", "http://elsewhere.example.com/moved?limit=3").fn(req)
		}
		return respond(200, "application/json", `{"data":[],"has_more":false}`).fn(req)
	}}
	client, recorder := newTracedClient(t, transport)
	if _, err := client.Models.List(context.Background(), anthropic.ModelListParams{Limit: anthropic.Int(1)}); err != nil {
		t.Fatal(err)
	}
	// The HTTP client followed the redirect: the SDK sent one request.
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"url.path":                  "/v1/models",
		"url.full":                  "http://localhost:4010/v1/models?limit=1",
		"http.request.resend_count": nil,
		"http.response.status_code": int64(200),
	})
}

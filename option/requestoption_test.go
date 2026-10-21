package option_test

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/auth"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

func TestIdentityTokenFile_ReadsAndTrims(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := os.WriteFile(path, []byte("  the-jwt\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	tok, err := option.IdentityTokenFile(path)(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "the-jwt" {
		t.Errorf("got %q, want %q", tok, "the-jwt")
	}
}

func TestIdentityTokenFile_RereadsOnEachCall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := os.WriteFile(path, []byte("v1"), 0o600); err != nil {
		t.Fatal(err)
	}
	fn := option.IdentityTokenFile(path)

	got1, err := fn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got1 != "v1" {
		t.Errorf("first call got %q, want %q", got1, "v1")
	}

	// Simulate token rotation.
	if err := os.WriteFile(path, []byte("v2"), 0o600); err != nil {
		t.Fatal(err)
	}
	got2, err := fn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got2 != "v2" {
		t.Errorf("after rotation got %q, want %q (function did not re-read)", got2, "v2")
	}
}

func TestIdentityTokenFile_Errors(t *testing.T) {
	ctx := context.Background()

	t.Run("empty path", func(t *testing.T) {
		if _, err := option.IdentityTokenFile("")(ctx); err == nil || !strings.Contains(err.Error(), "path is empty") {
			t.Errorf("got %v, want error containing 'path is empty'", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does-not-exist")
		if _, err := option.IdentityTokenFile(missing)(ctx); err == nil || !strings.Contains(err.Error(), "read") {
			t.Errorf("got %v, want read error", err)
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		empty := filepath.Join(t.TempDir(), "empty")
		if err := os.WriteFile(empty, []byte("   \n\t\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := option.IdentityTokenFile(empty)(ctx); err == nil || !strings.Contains(err.Error(), "is empty") {
			t.Errorf("got %v, want 'is empty' error", err)
		}
	})
}

// TestJoin_AppliesInOrder verifies that a joined option is equivalent to
// passing its members individually at the same position: members apply in
// order (a later WithMaxRetries wins) and the first error stops the join.
func TestJoin_AppliesInOrder(t *testing.T) {
	cfg := &requestconfig.RequestConfig{}
	joined := option.Join(
		option.WithMaxRetries(1),
		option.WithBaseURL("https://joined.example/"),
		option.WithMaxRetries(7),
	)
	if err := joined.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxRetries != 7 {
		t.Errorf("got MaxRetries %d, want the later member (7) to win", cfg.MaxRetries)
	}
	if got := cfg.BaseURL.String(); got != "https://joined.example/" {
		t.Errorf("got base URL %q, want %q", got, "https://joined.example/")
	}

	boom := errors.New("boom")
	applied := false
	failing := option.Join(
		requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return boom }),
		requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { applied = true; return nil }),
	)
	if err := failing.Apply(&requestconfig.RequestConfig{}); !errors.Is(err, boom) {
		t.Fatalf("got error %v, want the member's error", err)
	}
	if applied {
		t.Error("members after a failing member must not be applied")
	}
}

// TestHasWithoutEnvironmentDefaults_SeesThroughJoin pins the contract that
// platform option constructors rely on: a marker returned inside a Join
// (bedrock/vertex bundle it with their own options) must still be detected
// by anthropic.NewClient, however deeply the joins nest — while a Join with
// no marker must not be mistaken for one.
func TestHasWithoutEnvironmentDefaults_SeesThroughJoin(t *testing.T) {
	noop := option.WithHeader("x-noop", "1")
	testCases := []struct {
		name string
		opts []option.RequestOption
		want bool
	}{
		{
			name: "no options",
			opts: nil,
			want: false,
		},
		{
			name: "unrelated options only",
			opts: []option.RequestOption{noop, option.WithAPIKey("k")},
			want: false,
		},
		{
			name: "marker at the top level",
			opts: []option.RequestOption{noop, option.WithoutEnvironmentDefaults()},
			want: true,
		},
		{
			name: "marker nested one join deep",
			opts: []option.RequestOption{noop, option.Join(noop, option.WithoutEnvironmentDefaults(), noop)},
			want: true,
		},
		{
			name: "marker nested two joins deep",
			opts: []option.RequestOption{option.Join(noop, option.Join(noop, option.WithoutEnvironmentDefaults()))},
			want: true,
		},
		{
			name: "joins without a marker",
			opts: []option.RequestOption{option.Join(noop, option.Join(noop, option.WithAPIKey("k"))), option.Join()},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := option.HasWithoutEnvironmentDefaults(tc.opts); got != tc.want {
				t.Errorf("HasWithoutEnvironmentDefaults() = %v, want %v", got, tc.want)
			}
		})
	}
}

// apiServer answers token exchanges with "federated-token" and keeps the headers
// of the last API request, which stay nil if there was none.
type apiServer struct {
	*httptest.Server
	sent http.Header
}

func newAPIServer(t *testing.T) *apiServer {
	t.Helper()
	s := &apiServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			_, _ = w.Write([]byte(`{"access_token":"federated-token"}`))
			return
		}
		s.sent = r.Header.Clone()
		_, _ = w.Write([]byte(`{"data":[],"has_more":false}`))
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *apiServer) newClient(opts ...option.RequestOption) anthropic.Client {
	return anthropic.NewClient(append([]option.RequestOption{option.WithBaseURL(s.URL)}, opts...)...)
}

func listModels(client anthropic.Client, opts ...option.RequestOption) error {
	_, err := client.Models.List(context.Background(), anthropic.ModelListParams{}, opts...)
	return err
}

// sentHeaders builds a client from clientOpts, sends one request with reqOpts,
// and returns the headers that reached the server.
func sentHeaders(t *testing.T, clientOpts []option.RequestOption, reqOpts ...option.RequestOption) http.Header {
	t.Helper()
	server := newAPIServer(t)
	if err := listModels(server.newClient(clientOpts...), reqOpts...); err != nil {
		t.Fatal(err)
	}
	return server.sent
}

func assertAuthorization(t *testing.T, sent http.Header, want string) {
	t.Helper()
	if got := sent.Values("Authorization"); !slices.Equal(got, []string{want}) {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

// clearCredentialEnv unsets everything a client reads from the environment and
// points the default credential chain at an empty config directory.
func clearCredentialEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"ANTHROPIC_BASE_URL",
		"ANTHROPIC_API_KEY",
		"ANTHROPIC_AUTH_TOKEN",
		"ANTHROPIC_PROFILE",
		"ANTHROPIC_FEDERATION_RULE_ID",
		"ANTHROPIC_ORGANIZATION_ID",
		"ANTHROPIC_IDENTITY_TOKEN",
		"ANTHROPIC_IDENTITY_TOKEN_FILE",
		"ANTHROPIC_WEBHOOK_SIGNING_KEY",
		"ANTHROPIC_CUSTOM_HEADERS",
	} {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
}

// captureLog returns what the standard logger writes from now until the test
// ends, with the warnings that are given once per process ready to be given again.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var logged bytes.Buffer
	output, flags := log.Writer(), log.Flags()
	log.SetOutput(&logged)
	log.SetFlags(0)
	auth.ResetWarnOnceForTest()
	t.Cleanup(func() {
		log.SetOutput(output)
		log.SetFlags(flags)
		auth.ResetWarnOnceForTest()
	})
	return &logged
}

// signedWebhook returns a webhook payload and the headers that sign it with key.
func signedWebhook(t *testing.T, key string) ([]byte, http.Header) {
	t.Helper()
	payload := []byte(`{"id":"whe_fake","created_at":"2026-03-15T10:00:00Z","data":{"id":"sesn_fake","organization_id":"org_fake","type":"session.status_idled","workspace_id":"wrkspc_fake"},"type":"event"}`)
	wh, err := standardwebhooks.NewWebhook(key)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	signature, err := wh.Sign("1", now, payload)
	if err != nil {
		t.Fatal(err)
	}
	return payload, http.Header{
		"Webhook-Id":        {"1"},
		"Webhook-Timestamp": {strconv.FormatInt(now.Unix(), 10)},
		"Webhook-Signature": {signature},
	}
}

// tokenConfig returns a config that authenticates with token.
func tokenConfig(t *testing.T, token string) *config.Config {
	t.Helper()
	path := filepath.Join(t.TempDir(), "credentials.json")
	if err := config.WriteCredentials(path, config.Credentials{AccessToken: token}); err != nil {
		t.Fatal(err)
	}
	return &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: path,
			UserOAuth:       &config.UserOAuth{},
		},
	}
}

func saveProfile(t *testing.T, name, token string) {
	t.Helper()
	if err := config.SaveProfile(config.DefaultDir(), name, tokenConfig(t, token)); err != nil {
		t.Fatal(err)
	}
}

type credentialOption struct {
	name     string
	option   option.RequestOption
	wantAuth string
}

// credentialOptions returns one of every option that supplies a credential. Call
// it after [clearCredentialEnv], because it saves the profile WithProfile loads.
func credentialOptions(t *testing.T) []credentialOption {
	t.Helper()
	saveProfile(t, "work", "profile-token")
	return []credentialOption{
		{"WithAPIKey", option.WithAPIKey("key"), "Bearer key"},
		{"WithAuthToken", option.WithAuthToken("token"), "Bearer token"},
		{"WithConfig", option.WithConfig(tokenConfig(t, "config-token")), "Bearer config-token"},
		{"WithConfigQuiet", option.WithConfigQuiet(tokenConfig(t, "quiet-config-token")), "Bearer quiet-config-token"},
		{"WithProfile", option.WithProfile("work"), "Bearer profile-token"},
		{
			"WithFederationTokenProvider",
			option.WithFederationTokenProvider(
				func(context.Context) (string, error) { return "jwt", nil },
				option.FederationOptions{FederationRuleID: "fdrl_fake", OrganizationID: "fake-org"},
			),
			"Bearer federated-token",
		},
	}
}

// eachCredentialOptionPair runs f for every ordered pair of different options.
func eachCredentialOptionPair(t *testing.T, f func(t *testing.T, first, second credentialOption)) {
	clearCredentialEnv(t)
	opts := credentialOptions(t)
	for _, first := range opts {
		for _, second := range opts {
			if first.name == second.name {
				continue
			}
			t.Run(first.name+" then "+second.name, func(t *testing.T) { f(t, first, second) })
		}
	}
}

func TestAPIKeyHeader(t *testing.T) {
	setByHand := []option.RequestOption{
		option.WithHeader("X-Api-Key", "key"),
		option.WithHeaderDel("Authorization"),
	}
	tests := []struct {
		name       string
		env        map[string]string
		clientOpts []option.RequestOption
		reqOpts    []option.RequestOption
		wantAuth   string
		wantAPIKey string
	}{
		{
			name:       "WithAPIKey",
			clientOpts: []option.RequestOption{option.WithAPIKey("key")},
			wantAuth:   "Bearer key",
		},
		{
			name:     "ANTHROPIC_API_KEY",
			env:      map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			wantAuth: "Bearer env-key",
		},
		{
			name:     "both environment variables",
			env:      map[string]string{"ANTHROPIC_API_KEY": "env-key", "ANTHROPIC_AUTH_TOKEN": "env-token"},
			wantAuth: "Bearer env-key",
		},
		{
			name:       "WithAuthToken beats ANTHROPIC_API_KEY",
			env:        map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			clientOpts: []option.RequestOption{option.WithAuthToken("token")},
			wantAuth:   "Bearer token",
		},
		{
			name:       "WithAPIKey beats ANTHROPIC_AUTH_TOKEN",
			env:        map[string]string{"ANTHROPIC_AUTH_TOKEN": "env-token"},
			clientOpts: []option.RequestOption{option.WithAPIKey("key")},
			wantAuth:   "Bearer key",
		},
		{
			name:       "request WithAuthToken replaces the client's key",
			clientOpts: []option.RequestOption{option.WithAPIKey("key")},
			reqOpts:    []option.RequestOption{option.WithAuthToken("token")},
			wantAuth:   "Bearer token",
		},
		{
			name:       "the same option twice",
			clientOpts: []option.RequestOption{option.WithAPIKey("key"), option.WithAPIKey("other-key")},
			wantAuth:   "Bearer other-key",
		},
		{
			name:       "request WithConfig replaces the client's empty key",
			clientOpts: []option.RequestOption{option.WithAPIKey("")},
			reqOpts:    []option.RequestOption{option.WithConfig(tokenConfig(t, "config-token"))},
			wantAuth:   "Bearer config-token",
		},
		{
			name:       "Authorization set by hand",
			clientOpts: []option.RequestOption{option.WithAPIKey("key")},
			reqOpts: []option.RequestOption{
				option.WithHeader("Authorization", "Bearer by-hand"),
				option.WithConfig(tokenConfig(t, "config-token")),
			},
			wantAuth: "Bearer by-hand",
		},
		{
			name: "Authorization set by hand to the value the replaced key had",
			clientOpts: []option.RequestOption{
				option.WithAPIKey("key"),
				option.WithHeader("Authorization", "Bearer key"),
			},
			reqOpts:  []option.RequestOption{option.WithConfig(tokenConfig(t, "config-token"))},
			wantAuth: "Bearer key",
		},
		{
			name: "Authorization set by hand before the replaced key",
			clientOpts: []option.RequestOption{
				option.WithHeader("Authorization", "Bearer by-hand"),
				option.WithAPIKey("key"),
			},
			reqOpts:  []option.RequestOption{option.WithConfig(tokenConfig(t, "config-token"))},
			wantAuth: "Bearer by-hand",
		},
		{
			name:       "Authorization deleted",
			clientOpts: []option.RequestOption{option.WithAPIKey("key"), option.WithHeaderDel("Authorization")},
		},
		{
			name:       "X-Api-Key set by hand",
			clientOpts: setByHand,
			wantAPIKey: "key",
		},
		{
			name:       "X-Api-Key set by hand with ANTHROPIC_API_KEY",
			env:        map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			clientOpts: setByHand,
			wantAPIKey: "key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCredentialEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got := sentHeaders(t, tt.clientOpts, tt.reqOpts...)
			if auth := got.Values("Authorization"); len(auth) > 1 || got.Get("Authorization") != tt.wantAuth {
				t.Errorf("Authorization = %q, want %q", auth, tt.wantAuth)
			}
			if key := got.Get("X-Api-Key"); key != tt.wantAPIKey {
				t.Errorf("X-Api-Key = %q, want %q", key, tt.wantAPIKey)
			}
		})
	}
}

func TestTwoCredentialOptionsFail(t *testing.T) {
	eachCredentialOptionPair(t, func(t *testing.T, first, second credentialOption) {
		pair := []option.RequestOption{first.option, second.option}
		tests := []struct {
			name       string
			clientOpts []option.RequestOption
			reqOpts    []option.RequestOption
		}{
			{name: "client", clientOpts: pair},
			{name: "client, joined", clientOpts: []option.RequestOption{option.Join(pair...)}},
			{name: "request", clientOpts: []option.RequestOption{option.WithAPIKey("client-key")}, reqOpts: pair},
		}
		want := "option: pass at most one of WithAPIKey, WithAuthToken, WithConfig, WithConfigQuiet, WithProfile and " +
			"WithFederationTokenProvider; got " + first.name + " and " + second.name
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				server := newAPIServer(t)
				err := listModels(server.newClient(tt.clientOpts...), tt.reqOpts...)
				if err == nil || err.Error() != want {
					t.Errorf("error = %v, want %q", err, want)
				}
				if server.sent != nil {
					t.Errorf("a request was sent with Authorization %q", server.sent.Values("Authorization"))
				}
			})
		}
	})
}

func TestCredentialOptionIgnoresEnvironment(t *testing.T) {
	envs := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "key and token",
			env:  map[string]string{"ANTHROPIC_API_KEY": "env-key", "ANTHROPIC_AUTH_TOKEN": "env-token"},
		},
		{
			name: "profile that does not exist",
			env:  map[string]string{"ANTHROPIC_PROFILE": "missing"},
		},
	}
	clearCredentialEnv(t)
	for _, opt := range credentialOptions(t) {
		for _, tt := range envs {
			t.Run(opt.name+" with "+tt.name, func(t *testing.T) {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
				assertAuthorization(t, sentHeaders(t, []option.RequestOption{opt.option}), opt.wantAuth)
			})
		}
	}
}

func TestEnvironmentCredentialOrder(t *testing.T) {
	clearCredentialEnv(t)
	saveProfile(t, "default", "active-profile-token")
	saveProfile(t, "named", "named-profile-token")

	// Each source is unset once it has been seen to win over the ones after it.
	sources := []struct {
		env      map[string]string
		wantAuth string
	}{
		{map[string]string{"ANTHROPIC_API_KEY": "env-key"}, "Bearer env-key"},
		{map[string]string{"ANTHROPIC_AUTH_TOKEN": "env-token"}, "Bearer env-token"},
		{map[string]string{"ANTHROPIC_PROFILE": "named"}, "Bearer named-profile-token"},
		{
			map[string]string{
				"ANTHROPIC_FEDERATION_RULE_ID": "fdrl_fake",
				"ANTHROPIC_ORGANIZATION_ID":    "fake-org",
				"ANTHROPIC_IDENTITY_TOKEN":     "jwt",
			},
			"Bearer federated-token",
		},
		{nil, "Bearer active-profile-token"},
	}
	for _, source := range sources {
		for k, v := range source.env {
			t.Setenv(k, v)
		}
	}
	for _, source := range sources {
		assertAuthorization(t, sentHeaders(t, nil), source.wantAuth)
		for k := range source.env {
			os.Unsetenv(k)
		}
	}
}

func TestEnvironmentAPIKeyAndAuthTokenWarning(t *testing.T) {
	const warning = "Anthropic: ANTHROPIC_API_KEY and ANTHROPIC_AUTH_TOKEN are both set. Using ANTHROPIC_API_KEY.\n"
	both := map[string]string{"ANTHROPIC_API_KEY": "env-key", "ANTHROPIC_AUTH_TOKEN": "env-token"}
	type test struct {
		name        string
		env         map[string]string
		clientOpts  []option.RequestOption
		wantAuth    string
		wantWarning string
	}
	tests := []test{
		{name: "both set", env: both, wantAuth: "Bearer env-key", wantWarning: warning},
		{
			name:     "only the key",
			env:      map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			wantAuth: "Bearer env-key",
		},
		{
			name:     "only the token",
			env:      map[string]string{"ANTHROPIC_AUTH_TOKEN": "env-token"},
			wantAuth: "Bearer env-token",
		},
		{
			name:     "empty token",
			env:      map[string]string{"ANTHROPIC_API_KEY": "env-key", "ANTHROPIC_AUTH_TOKEN": ""},
			wantAuth: "Bearer env-key",
		},
		{
			name:     "empty key",
			env:      map[string]string{"ANTHROPIC_API_KEY": "", "ANTHROPIC_AUTH_TOKEN": "env-token"},
			wantAuth: "Bearer env-token",
		},
		{
			name: "both set with WithoutEnvironmentDefaults",
			env:  both,
			clientOpts: []option.RequestOption{
				option.WithoutEnvironmentDefaults(),
				option.WithHeader("Authorization", "Bearer by-hand"),
			},
			wantAuth: "Bearer by-hand",
		},
	}
	clearCredentialEnv(t)
	for _, opt := range credentialOptions(t) {
		tests = append(tests, test{
			name:       "both set with " + opt.name,
			env:        both,
			clientOpts: []option.RequestOption{opt.option},
			wantAuth:   opt.wantAuth,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			logged := captureLog(t)
			for range 2 {
				assertAuthorization(t, sentHeaders(t, tt.clientOpts), tt.wantAuth)
				if got := logged.String(); got != tt.wantWarning {
					t.Errorf("logged %q, want %q", got, tt.wantWarning)
				}
			}
		})
	}
}

func TestCustomHeadersAndWebhookKeyFromEnvironment(t *testing.T) {
	const envWebhookKey, codeWebhookKey = "whsec_ZW52LWtleQ==", "whsec_Y29kZS1rZXk="
	tests := []struct {
		name           string
		env            map[string]string
		clientOpts     []option.RequestOption
		wantAuth       string
		wantHeader     string
		wantWebhookKey string
	}{
		{
			name:     "ANTHROPIC_API_KEY",
			env:      map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			wantAuth: "Bearer env-key",
		},
		{
			name:     "ANTHROPIC_AUTH_TOKEN",
			env:      map[string]string{"ANTHROPIC_AUTH_TOKEN": "env-token"},
			wantAuth: "Bearer env-token",
		},
		{
			name:     "ANTHROPIC_PROFILE",
			env:      map[string]string{"ANTHROPIC_PROFILE": "named"},
			wantAuth: "Bearer named-profile-token",
		},
		{
			name: "federation",
			env: map[string]string{
				"ANTHROPIC_FEDERATION_RULE_ID": "fdrl_fake",
				"ANTHROPIC_ORGANIZATION_ID":    "fake-org",
				"ANTHROPIC_IDENTITY_TOKEN":     "jwt",
			},
			wantAuth: "Bearer federated-token",
		},
		{
			name:     "active profile",
			wantAuth: "Bearer active-profile-token",
		},
		{
			name:       "WithAPIKey with ANTHROPIC_API_KEY",
			env:        map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			clientOpts: []option.RequestOption{option.WithAPIKey("key")},
			wantAuth:   "Bearer key",
		},
		{
			name:       "WithAuthToken with the active profile",
			clientOpts: []option.RequestOption{option.WithAuthToken("token")},
			wantAuth:   "Bearer token",
		},
		{
			name:       "WithHeader",
			env:        map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			clientOpts: []option.RequestOption{option.WithHeader("X-Custom", "from-code")},
			wantAuth:   "Bearer env-key",
			wantHeader: "from-code",
		},
		{
			name:           "WithWebhookKey",
			env:            map[string]string{"ANTHROPIC_API_KEY": "env-key"},
			clientOpts:     []option.RequestOption{option.WithWebhookKey(codeWebhookKey)},
			wantAuth:       "Bearer env-key",
			wantWebhookKey: codeWebhookKey,
		},
	}
	clearCredentialEnv(t)
	saveProfile(t, "default", "active-profile-token")
	saveProfile(t, "named", "named-profile-token")
	t.Setenv("ANTHROPIC_CUSTOM_HEADERS", "X-Custom: from-env")
	t.Setenv("ANTHROPIC_WEBHOOK_SIGNING_KEY", envWebhookKey)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			server := newAPIServer(t)
			client := server.newClient(tt.clientOpts...)
			if err := listModels(client); err != nil {
				t.Fatal(err)
			}
			assertAuthorization(t, server.sent, tt.wantAuth)
			wantHeader := cmp.Or(tt.wantHeader, "from-env")
			if got := server.sent.Values("X-Custom"); !slices.Equal(got, []string{wantHeader}) {
				t.Errorf("X-Custom = %q, want %q", got, wantHeader)
			}

			payload, headers := signedWebhook(t, cmp.Or(tt.wantWebhookKey, envWebhookKey))
			if _, err := client.Beta.Webhooks.Unwrap(payload, headers); err != nil {
				t.Errorf("Unwrap: %v", err)
			}
		})
	}
}

func TestRequestCredentialOptionReplacesTheClients(t *testing.T) {
	eachCredentialOptionPair(t, func(t *testing.T, first, second credentialOption) {
		assertAuthorization(t, sentHeaders(t, []option.RequestOption{first.option}, second.option), second.wantAuth)
	})
}

// A service that was not built by a client gets its own options and a method's
// as one list, so a second credential option may be the method's replacement.
func TestServiceWithoutClientUsesTheLastCredentialOption(t *testing.T) {
	eachCredentialOptionPair(t, func(t *testing.T, first, second credentialOption) {
		server := newAPIServer(t)
		service := anthropic.NewModelService(option.WithBaseURL(server.URL), first.option)
		if _, err := service.List(context.Background(), anthropic.ModelListParams{}, second.option); err != nil {
			t.Fatal(err)
		}
		assertAuthorization(t, server.sent, second.wantAuth)
	})
}

func TestMiddlewareAroundACredentialOption(t *testing.T) {
	tests := []struct {
		name       string
		reqOpts    []option.RequestOption
		wantBefore string
		wantAfter  string
	}{
		{name: "runs in the order it was added", wantBefore: "", wantAfter: "Bearer config-token"},
		{
			name:       "stays when the credential is replaced",
			reqOpts:    []option.RequestOption{option.WithAuthToken("token")},
			wantBefore: "Bearer token",
			wantAfter:  "Bearer token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCredentialEnv(t)
			saw := map[string]string{}
			record := func(name string) option.RequestOption {
				return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
					saw[name] = req.Header.Get("Authorization")
					return next(req)
				})
			}
			sentHeaders(t, []option.RequestOption{
				record("before"),
				option.WithConfig(tokenConfig(t, "config-token")),
				record("after"),
			}, tt.reqOpts...)
			if want := map[string]string{"before": tt.wantBefore, "after": tt.wantAfter}; !maps.Equal(saw, want) {
				t.Errorf("Authorization seen by each middleware = %v, want %v", saw, want)
			}
		})
	}
}

func TestConfigBaseURLAndWorkspace(t *testing.T) {
	tests := []struct {
		name          string
		clientOpts    []option.RequestOption
		reqOpt        func(other *apiServer) option.RequestOption
		wantOther     bool
		wantAuth      string
		wantWorkspace string
	}{
		{
			name: "replaced by those of another config",
			reqOpt: func(other *apiServer) option.RequestOption {
				cfg := tokenConfig(t, "other-token")
				cfg.BaseURL = other.URL
				cfg.WorkspaceID = "wrkspc_other"
				return option.WithConfig(cfg)
			},
			wantOther:     true,
			wantAuth:      "Bearer other-token",
			wantWorkspace: "wrkspc_other",
		},
		{
			name: "removed by a config that has none",
			reqOpt: func(*apiServer) option.RequestOption {
				return option.WithConfig(tokenConfig(t, "other-token"))
			},
			wantOther: true,
			wantAuth:  "Bearer other-token",
		},
		{
			name:       "workspace set by hand is not removed",
			clientOpts: []option.RequestOption{option.WithHeader("anthropic-workspace-id", "wrkspc_by_hand")},
			reqOpt: func(*apiServer) option.RequestOption {
				return option.WithConfig(tokenConfig(t, "other-token"))
			},
			wantOther:     true,
			wantAuth:      "Bearer other-token",
			wantWorkspace: "wrkspc_by_hand",
		},
		{
			name:          "kept by any other credential option",
			reqOpt:        func(*apiServer) option.RequestOption { return option.WithAuthToken("token") },
			wantAuth:      "Bearer token",
			wantWorkspace: "wrkspc_client",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCredentialEnv(t)
			server, other := newAPIServer(t), newAPIServer(t)
			cfg := tokenConfig(t, "client-token")
			cfg.BaseURL = server.URL
			cfg.WorkspaceID = "wrkspc_client"

			clientOpts := append([]option.RequestOption{requestconfig.WithDefaultBaseURL(other.URL), option.WithConfig(cfg)}, tt.clientOpts...)
			client := anthropic.NewClient(clientOpts...)
			if err := listModels(client, tt.reqOpt(other)); err != nil {
				t.Fatal(err)
			}
			want, unwanted := server, other
			if tt.wantOther {
				want, unwanted = other, server
			}
			if unwanted.sent != nil {
				t.Errorf("Authorization %q went to the wrong server", unwanted.sent.Values("Authorization"))
			}
			assertAuthorization(t, want.sent, tt.wantAuth)
			if got := want.sent.Get("anthropic-workspace-id"); got != tt.wantWorkspace {
				t.Errorf("anthropic-workspace-id = %q, want %q", got, tt.wantWorkspace)
			}
		})
	}
}

func TestRequestCredentialOptionReplacesTheEnvironments(t *testing.T) {
	envs := []struct {
		name string
		env  map[string]string
	}{
		{
			name: "key and token",
			env:  map[string]string{"ANTHROPIC_API_KEY": "env-key", "ANTHROPIC_AUTH_TOKEN": "env-token"},
		},
		{
			name: "profile that does not exist",
			env:  map[string]string{"ANTHROPIC_PROFILE": "missing"},
		},
		{
			name: "federation",
			env: map[string]string{
				"ANTHROPIC_FEDERATION_RULE_ID": "fdrl_fake",
				"ANTHROPIC_ORGANIZATION_ID":    "fake-org",
				"ANTHROPIC_IDENTITY_TOKEN":     "jwt",
			},
		},
		{name: "nothing"},
	}
	clearCredentialEnv(t)
	for _, opt := range credentialOptions(t) {
		for _, tt := range envs {
			t.Run(opt.name+" with "+tt.name, func(t *testing.T) {
				for k, v := range tt.env {
					t.Setenv(k, v)
				}
				assertAuthorization(t, sentHeaders(t, nil, opt.option), opt.wantAuth)
			})
		}
	}
}

func TestClientBuiltFromAnothersOptions(t *testing.T) {
	t.Run("keeps the credential", func(t *testing.T) {
		clearCredentialEnv(t)
		for _, opt := range credentialOptions(t) {
			t.Run(opt.name, func(t *testing.T) {
				server := newAPIServer(t)
				parent := server.newClient(opt.option)
				t.Setenv("ANTHROPIC_API_KEY", "env-key")
				if err := listModels(anthropic.NewClient(parent.Options...)); err != nil {
					t.Fatal(err)
				}
				assertAuthorization(t, server.sent, opt.wantAuth)
			})
		}
	})

	t.Run("keeps the credential from the environment", func(t *testing.T) {
		clearCredentialEnv(t)
		t.Setenv("ANTHROPIC_API_KEY", "env-key")
		t.Setenv("ANTHROPIC_AUTH_TOKEN", "env-token")
		server := newAPIServer(t)
		if err := listModels(anthropic.NewClient(server.newClient().Options...)); err != nil {
			t.Fatal(err)
		}
		assertAuthorization(t, server.sent, "Bearer env-key")
	})

	t.Run("replaces the credential", func(t *testing.T) {
		eachCredentialOptionPair(t, func(t *testing.T, first, second credentialOption) {
			server := newAPIServer(t)
			parent := server.newClient(first.option)
			if err := listModels(anthropic.NewClient(append(slices.Clone(parent.Options), second.option)...)); err != nil {
				t.Fatal(err)
			}
			assertAuthorization(t, server.sent, second.wantAuth)
		})
	})

	t.Run("replaces the credential from the environment", func(t *testing.T) {
		clearCredentialEnv(t)
		opts := credentialOptions(t)
		t.Setenv("ANTHROPIC_API_KEY", "env-key")
		t.Setenv("ANTHROPIC_AUTH_TOKEN", "env-token")
		for _, opt := range opts {
			t.Run(opt.name, func(t *testing.T) {
				server := newAPIServer(t)
				parent := server.newClient()
				if err := listModels(anthropic.NewClient(append(slices.Clone(parent.Options), opt.option)...)); err != nil {
					t.Fatal(err)
				}
				assertAuthorization(t, server.sent, opt.wantAuth)
			})
		}
	})
}

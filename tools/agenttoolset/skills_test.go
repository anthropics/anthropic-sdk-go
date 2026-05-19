package agenttoolset

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// SetupSkills must apply the request options it is given (the environment key,
// for self-hosted callers) to its API calls. The session lookup and skill
// endpoints are environment-scoped: if the per-call options are dropped the
// request falls back to the client's default credentials and fails. This
// guards the regression where SetupSkills ignored its opts and skills were
// silently never downloaded under ANTHROPIC_ENVIRONMENT_KEY.
func TestSetupSkills_AppliesRequestOptions(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		// No skills -> SetupSkills does only the session lookup and returns.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent":{"skills":[]}}`))
	}))
	defer srv.Close()

	client := anthropic.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("client-default-key"),
		option.WithMaxRetries(0),
	)
	env := &AgentToolContext{Workdir: t.TempDir()}

	if err := env.SetupSkills(context.Background(), client, "sess_x",
		option.WithAuthToken("env-key-xyz")); err != nil {
		t.Fatalf("SetupSkills returned error: %v", err)
	}

	if want := "Bearer env-key-xyz"; gotAuth != want {
		t.Errorf("session lookup Authorization = %q, want %q (per-call options were dropped)", gotAuth, want)
	}
}

func TestNumericVersionHelpers(t *testing.T) {
	for _, s := range []string{"0", "1", "1759178010641129"} {
		if !isNumericString(s) {
			t.Errorf("isNumericString(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"", "latest", "12a", "1.0", "v1"} {
		if isNumericString(s) {
			t.Errorf("isNumericString(%q) = true, want false", s)
		}
	}
	cases := []struct {
		a, b string
		want bool
	}{
		{"2", "10", false},   // 2 < 10 (length wins)
		{"10", "2", true},    // 10 > 2
		{"100", "99", true},  // length wins
		{"300", "200", true}, // same length, lexical
		{"200", "300", false},
		{"1759178010641130", "1759178010641129", true},
	}
	for _, c := range cases {
		if got := numericGreater(c.a, c.b); got != c.want {
			t.Errorf("numericGreater(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

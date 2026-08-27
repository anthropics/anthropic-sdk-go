package agenttoolset

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// SetupSkills must apply the request options it is given (the environment key,
// for self-hosted callers) to its API calls. The skill endpoints are
// environment-scoped: if the per-call options are dropped the request falls
// back to the client's default credentials and fails. This guards the
// regression where SetupSkills ignored its opts and skills were silently never
// downloaded under ANTHROPIC_ENVIRONMENT_KEY.
//
// The caller fetches the session and hands the snapshot in — SetupSkills never
// looks it up — so the skill version get and download are its only requests.
func TestSetupSkills_AppliesRequestOptions(t *testing.T) {
	var mu sync.Mutex
	var skillAuth []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record := func() {
			mu.Lock()
			defer mu.Unlock()
			skillAuth = append(skillAuth, r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/v1/sessions/sess_x":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"agent":{"skills":[{"skill_id":"skill_1","version":"1"}]}}`))
		case "/v1/skills/skill_1/versions/1":
			record()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"1","name":"demo"}`))
		case "/v1/skills/skill_1/versions/1/content":
			record()
			_, _ = w.Write(zipBytes(t, map[string]string{"SKILL.md": "hello"}))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := anthropic.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("client-default-key"),
		option.WithMaxRetries(0),
	)
	session, err := client.Beta.Sessions.Get(context.Background(), "sess_x", anthropic.BetaSessionGetParams{})
	if err != nil {
		t.Fatalf("session lookup returned error: %v", err)
	}

	workdir := t.TempDir()
	env := &AgentToolContext{Workdir: workdir}
	if err := env.SetupSkillsFromSession(context.Background(), client, session,
		option.WithAuthToken("env-key-xyz")); err != nil {
		t.Fatalf("SetupSkills returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	want := []string{"Bearer env-key-xyz", "Bearer env-key-xyz"}
	if !slices.Equal(skillAuth, want) {
		t.Errorf("skill request Authorization = %q, want %q (per-call options were dropped)", skillAuth, want)
	}
	// A silently skipped skill would leave the same empty workdir as a dropped
	// option, so assert the download actually landed.
	if _, err := os.Stat(filepath.Join(workdir, "skills", "demo", "SKILL.md")); err != nil {
		t.Errorf("skill was not extracted: %v", err)
	}
}

// Cleanup must remove only the skill directories the setup itself downloaded.
// The skills dir is inside the caller's workdir; skills the caller placed
// there by hand are its content, not ours to delete.
func TestCleanup_RemovesOnlyDownloadedSkills(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/skills/skill_1/versions/1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"1","name":"demo"}`))
		case "/v1/skills/skill_1/versions/1/content":
			_, _ = w.Write(zipBytes(t, map[string]string{"SKILL.md": "hello"}))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := anthropic.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))

	workdir := t.TempDir()
	preSeeded := filepath.Join(workdir, "skills", "my-local-skill")
	if err := os.MkdirAll(preSeeded, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preSeeded, "file.md"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	var session anthropic.BetaManagedAgentsSession
	if err := session.UnmarshalJSON([]byte(`{"agent":{"skills":[{"skill_id":"skill_1","version":"1"}]}}`)); err != nil {
		t.Fatal(err)
	}
	env := &AgentToolContext{Workdir: workdir}
	if err := env.SetupSkillsFromSession(context.Background(), client, &session); err != nil {
		t.Fatalf("SetupSkillsFromSession returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workdir, "skills", "demo", "SKILL.md")); err != nil {
		t.Fatalf("skill was not extracted: %v", err)
	}

	if err := env.Cleanup(); err != nil {
		t.Fatalf("Cleanup returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(preSeeded, "file.md")); err != nil {
		t.Errorf("Cleanup deleted the caller's own skill: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workdir, "skills", "demo")); !os.IsNotExist(err) {
		t.Errorf("Cleanup left the downloaded skill behind (stat err = %v)", err)
	}
}

// A session with no skills downloads nothing, so Cleanup has nothing to remove
// — caller-placed content under {Workdir}/skills is untouched.
func TestCleanup_NoDownloadsIsANoOp(t *testing.T) {
	workdir := t.TempDir()
	preSeeded := filepath.Join(workdir, "skills", "my-local-skill")
	if err := os.MkdirAll(preSeeded, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preSeeded, "file.md"), []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}

	client := anthropic.NewClient(option.WithAPIKey("k"), option.WithMaxRetries(0))
	var session anthropic.BetaManagedAgentsSession
	if err := session.UnmarshalJSON([]byte(`{"agent":{"skills":[]}}`)); err != nil {
		t.Fatal(err)
	}
	env := &AgentToolContext{Workdir: workdir}
	if err := env.SetupSkillsFromSession(context.Background(), client, &session); err != nil {
		t.Fatalf("SetupSkillsFromSession returned error: %v", err)
	}
	if err := env.Cleanup(); err != nil {
		t.Fatalf("Cleanup returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(preSeeded, "file.md")); err != nil {
		t.Errorf("Cleanup deleted the caller's own skill: %v", err)
	}
}

// The configured skill version may be the alias "latest", which only the
// version retrieve endpoint resolves: SetupSkills retrieves by the configured
// version and downloads by the concrete id that came back.
func TestSetupSkills_DownloadsByResolvedVersionID(t *testing.T) {
	var mu sync.Mutex
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/v1/skills/skill_1/versions/latest":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"skillver_123","name":"demo"}`))
		case "/v1/skills/skill_1/versions/skillver_123/content":
			_, _ = w.Write(zipBytes(t, map[string]string{"SKILL.md": "hello"}))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	client := anthropic.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))

	var session anthropic.BetaManagedAgentsSession
	if err := session.UnmarshalJSON([]byte(`{"agent":{"skills":[{"skill_id":"skill_1","version":"latest"}]}}`)); err != nil {
		t.Fatal(err)
	}
	workdir := t.TempDir()
	env := &AgentToolContext{Workdir: workdir}
	if err := env.SetupSkillsFromSession(context.Background(), client, &session); err != nil {
		t.Fatalf("SetupSkillsFromSession returned error: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	want := []string{"/v1/skills/skill_1/versions/latest", "/v1/skills/skill_1/versions/skillver_123/content"}
	if !slices.Equal(paths, want) {
		t.Errorf("requests = %q, want %q", paths, want)
	}
	if _, err := os.Stat(filepath.Join(workdir, "skills", "demo", "SKILL.md")); err != nil {
		t.Errorf("skill was not extracted: %v", err)
	}
}

// SetupSkills is the back-compat entry point for callers that hold only a
// session id: it fetches the session, applying the same per-call options as the
// skill requests (the session endpoint is environment-scoped too), and then
// does exactly what SetupSkills does.
func TestSetupSkills_FetchesTheSessionWithTheSameOptions(t *testing.T) {
	var mu sync.Mutex
	var auth []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		auth = append(auth, r.URL.Path+" "+r.Header.Get("Authorization"))
		mu.Unlock()
		switch r.URL.Path {
		case "/v1/sessions/sess_x":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"agent":{"skills":[{"skill_id":"skill_1","version":"1"}]}}`))
		case "/v1/skills/skill_1/versions/1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"1","name":"demo"}`))
		case "/v1/skills/skill_1/versions/1/content":
			_, _ = w.Write(zipBytes(t, map[string]string{"SKILL.md": "hello"}))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := anthropic.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("client-default-key"),
		option.WithMaxRetries(0),
	)
	workdir := t.TempDir()
	env := &AgentToolContext{Workdir: workdir}
	if err := env.SetupSkills(context.Background(), client, "sess_x",
		option.WithAuthToken("env-key-xyz")); err != nil {
		t.Fatalf("SetupSkills returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	// The extra session fetch this entry point costs, then the same two skill
	// requests SetupSkills makes — all three carrying the caller's credential.
	want := []string{
		"/v1/sessions/sess_x Bearer env-key-xyz",
		"/v1/skills/skill_1/versions/1 Bearer env-key-xyz",
		"/v1/skills/skill_1/versions/1/content Bearer env-key-xyz",
	}
	if !slices.Equal(auth, want) {
		t.Errorf("requests = %q, want %q", auth, want)
	}
	if _, err := os.Stat(filepath.Join(workdir, "skills", "demo", "SKILL.md")); err != nil {
		t.Errorf("skill was not extracted: %v", err)
	}
}

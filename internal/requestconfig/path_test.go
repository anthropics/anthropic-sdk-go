package requestconfig

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format string
		params []string
		want   string
	}{
		{
			name:   "ordinary beta path unchanged",
			format: "v1/vaults/%s/credentials/%s?beta=true",
			params: []string{"vault_A", "cred_x"},
			want:   "v1/vaults/vault_A/credentials/cred_x?beta=true",
		},
		{
			name:   "slash in id stays in one segment",
			format: "v1/vaults/%s/credentials/%s?beta=true",
			params: []string{"vault_A", "foo/bar"},
			want:   "v1/vaults/vault_A/credentials/foo%2Fbar?beta=true",
		},
		{
			name:   "dot-dot stays in the id segment",
			format: "v1/vaults/%s/credentials/%s?beta=true",
			params: []string{"vault_A", "../../vault_B/credentials/cred_x"},
			want:   "v1/vaults/vault_A/credentials/..%2F..%2Fvault_B%2Fcredentials%2Fcred_x?beta=true",
		},
		{
			name:   "question mark does not become query",
			format: "v1/vaults/%s/credentials/%s?beta=true",
			params: []string{"vault_A", "cred?injected=1"},
			want:   "v1/vaults/vault_A/credentials/cred%3Finjected=1?beta=true",
		},
		{
			name:   "hash does not drop suffix",
			format: "v1/files/%s/content",
			params: []string{"foo#bar"},
			want:   "v1/files/foo%23bar/content",
		},
		{
			name:   "bare dot is encoded",
			format: "v1/files/%s",
			params: []string{"."},
			want:   "v1/files/%2E",
		},
		{
			name:   "bare dot-dot is encoded",
			format: "v1/files/%s",
			params: []string{".."},
			want:   "v1/files/%2E%2E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, FormatPath(tt.format, tt.params...))
		})
	}
}

func TestExecuteNewRequestPathParameters(t *testing.T) {
	t.Parallel()

	type captured struct {
		uri      string
		escaped  string
		rawQuery string
		fragment string
	}

	run := func(t *testing.T, path string) (captured, error) {
		t.Helper()
		var got captured
		var hits int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits++
			got = captured{
				uri:      r.URL.RequestURI(),
				escaped:  r.URL.EscapedPath(),
				rawQuery: r.URL.RawQuery,
				fragment: r.URL.Fragment,
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{}`)
		}))
		t.Cleanup(server.Close)

		err := ExecuteNewRequest(context.Background(), http.MethodGet, path, nil, nil, WithDefaultBaseURL(server.URL+"/"))
		require.Equal(t, 1, hits, "expected a single request, got %d (err=%v)", hits, err)
		return got, err
	}

	t.Run("slash in id is percent-encoded", func(t *testing.T) {
		t.Parallel()
		path := FormatPath("v1/vaults/%s/credentials/%s?beta=true", "vault_A", "foo/bar")
		got, err := run(t, path)
		require.NoError(t, err)
		assert.Contains(t, got.escaped, "foo%2Fbar")
		assert.NotContains(t, strings.TrimPrefix(got.escaped, "/v1/vaults/vault_A/credentials/"), "/")
	})

	t.Run("dot-dot does not escape intended prefix", func(t *testing.T) {
		t.Parallel()
		path := FormatPath("v1/vaults/%s/credentials/%s?beta=true", "vault_A", "../../vault_B/credentials/cred_x")
		got, err := run(t, path)
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(got.escaped, "/v1/vaults/vault_A/credentials/"), "path %q escaped the vault_A prefix", got.escaped)
		assert.Contains(t, got.escaped, "..%2F..%2Fvault_B")
		assert.NotContains(t, got.escaped, "/vault_B/")
	})

	t.Run("question mark does not inject query", func(t *testing.T) {
		t.Parallel()
		path := FormatPath("v1/vaults/%s/credentials/%s?beta=true", "vault_A", "cred?injected=1")
		got, err := run(t, path)
		require.NoError(t, err)
		assert.NotContains(t, got.rawQuery, "injected")
		q, err := url.ParseQuery(got.rawQuery)
		require.NoError(t, err)
		assert.Equal(t, "true", q.Get("beta"))
		assert.Contains(t, got.uri, "cred%3Finjected=1")
	})

	t.Run("hash does not drop suffix", func(t *testing.T) {
		t.Parallel()
		path := FormatPath("v1/files/%s/content", "foo#bar")
		got, err := run(t, path)
		require.NoError(t, err)
		assert.Empty(t, got.fragment)
		assert.True(t, strings.HasSuffix(got.escaped, "/content"), "path %q dropped /content suffix", got.escaped)
		assert.Contains(t, got.uri, "foo%23bar")
	})
}

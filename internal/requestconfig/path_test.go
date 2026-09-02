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

func TestEscapeRequestPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr string
	}{
		{
			name: "ordinary beta path unchanged",
			in:   "v1/vaults/vault_A/credentials/cred_x?beta=true",
			want: "v1/vaults/vault_A/credentials/cred_x?beta=true",
		},
		{
			name: "ordinary ga path unchanged",
			in:   "v1/files/file_id/content",
			want: "v1/files/file_id/content",
		},
		{
			name: "static beta query preserved",
			in:   "v1/messages?beta=true",
			want: "v1/messages?beta=true",
		},
		{
			name:    "dot-dot segments rejected",
			in:      "v1/vaults/vault_A/credentials/../../vault_B/credentials/cred_x?beta=true",
			wantErr: "invalid path segment",
		},
		{
			name:    "dot segment rejected",
			in:      "v1/files/./content",
			wantErr: "invalid path segment",
		},
		{
			name:    "percent-encoded dot-dot rejected",
			in:      "v1/files/%2e%2e/content",
			wantErr: "invalid path segment",
		},
		{
			name: "question mark in id does not become query",
			in:   "v1/vaults/vault_A/credentials/cred?injected=1?beta=true",
			want: "v1/vaults/vault_A/credentials/cred%3Finjected=1?beta=true",
		},
		{
			name: "question mark in ga id does not become query",
			in:   "v1/files/cred?injected=1",
			want: "v1/files/cred%3Finjected=1",
		},
		{
			name: "hash in id does not drop suffix",
			in:   "v1/files/foo#bar/content",
			want: "v1/files/foo%23bar/content",
		},
		{
			name: "hash in id keeps beta suffix",
			in:   "v1/files/foo#bar/content?beta=true",
			want: "v1/files/foo%23bar/content?beta=true",
		},
		{
			name: "injected query cannot override beta flag",
			in:   "v1/vaults/vault_A/credentials/cred?beta=false&tail=?beta=true",
			want: "v1/vaults/vault_A/credentials/cred%3Fbeta=false&tail=?beta=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := escapeRequestPath(tt.in)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExecuteNewRequestPathParameters(t *testing.T) {
	t.Parallel()

	type captured struct {
		uri      string
		path     string
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
				path:     r.URL.Path,
				rawQuery: r.URL.RawQuery,
				fragment: r.URL.Fragment,
			}
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{}`)
		}))
		t.Cleanup(server.Close)

		err := ExecuteNewRequest(context.Background(), http.MethodGet, path, nil, nil, WithDefaultBaseURL(server.URL+"/"))
		if err != nil && hits == 0 {
			return captured{}, err
		}
		require.Equal(t, 1, hits, "expected a single request, got %d (err=%v)", hits, err)
		return got, err
	}

	t.Run("dot-dot does not escape intended prefix", func(t *testing.T) {
		t.Parallel()
		got, err := run(t, "v1/vaults/vault_A/credentials/../../vault_B/credentials/cred_x?beta=true")
		if err == nil {
			assert.True(t, strings.HasPrefix(got.path, "/v1/vaults/vault_A/credentials/"), "path %q escaped the vault_A prefix", got.path)
			assert.NotContains(t, got.path, "/vault_B/")
			return
		}
		assert.Contains(t, err.Error(), "invalid path segment")
		assert.Empty(t, got.path)
	})

	t.Run("question mark does not inject query", func(t *testing.T) {
		t.Parallel()
		got, err := run(t, "v1/vaults/vault_A/credentials/cred?injected=1?beta=true")
		require.NoError(t, err)
		assert.NotContains(t, got.rawQuery, "injected")
		q, err := url.ParseQuery(got.rawQuery)
		require.NoError(t, err)
		assert.Equal(t, "true", q.Get("beta"))
		assert.Contains(t, got.uri, "cred%3Finjected=1")
		assert.True(t, strings.HasPrefix(got.path, "/v1/vaults/vault_A/credentials/"))
	})

	t.Run("hash does not drop suffix", func(t *testing.T) {
		t.Parallel()
		got, err := run(t, "v1/files/foo#bar/content")
		require.NoError(t, err)
		assert.Empty(t, got.fragment)
		assert.True(t, strings.HasSuffix(got.path, "/content"), "path %q dropped /content suffix", got.path)
		assert.Contains(t, got.uri, "foo%23bar")
	})
}

//go:build !windows

// The fixtures below are extensionless #!/bin/sh scripts and the vectors are
// the POSIX-runnable ones; Windows resolution (PATHEXT, implicit ".") is
// exec.LookPath's own behaviour and is covered by the standard library's tests.

package agenttoolset

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// lookPathFixture lays out the directories the lookPath tests search:
//
//   - bin:   an absolute directory holding a real executable rg — the hit a
//     correct lookup may return;
//   - plant: the process working directory for the test (t.Chdir), holding
//     an executable rg and rel/sub/rg — files a lookup must never return;
//   - dirs:  an absolute directory whose "rg" is a directory, not a file;
//   - noexec: an absolute directory whose rg is a regular 0o644 file.
type lookPathFixture struct {
	bin, binRg, plant, dirs, noexec string
}

func newLookPathFixture(t *testing.T) lookPathFixture {
	t.Helper()
	f := lookPathFixture{bin: t.TempDir(), plant: t.TempDir(), dirs: t.TempDir(), noexec: t.TempDir()}
	f.binRg = writeExecutable(t, f.bin, "rg", "exit 0")
	writeExecutable(t, f.plant, "rg", "exit 0")
	writeExecutable(t, f.plant, filepath.Join("rel", "sub", "rg"), "exit 0")
	require.NoError(t, os.Mkdir(filepath.Join(f.dirs, "rg"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(f.noexec, "rg"), []byte("#!/bin/sh\nexit 0\n"), 0o644))
	t.Chdir(f.plant)
	return f
}

// TestLookPath pins the executable-lookup invariant documented on lookPath.
// The V*/P* labels are the test vectors shared with the other Anthropic SDKs'
// resolvers so the suites can be compared side by side.
func TestLookPath(t *testing.T) {
	f := newLookPathFixture(t)
	backslashName := writeExecutable(t, f.bin, `a\b`, "exit 0")

	tests := []struct {
		description string
		path        []string // PATH entries, joined with the OS list separator
		name        string
		// want is the only path lookPath may return for this row; "" means it
		// must report not-found.
		want string
		// notFoundOK marks rows where the sibling SDKs' resolvers skip the
		// working-directory-relative entry and return want, whereas Go's
		// exec.LookPath stops at that entry with exec.ErrDot, which lookPath
		// reports as not-found. Both outcomes are safe; the planted file is
		// what must never come back.
		notFoundOK bool
	}{
		{
			description: "V1: a bare name is found in an absolute PATH directory",
			path:        []string{f.bin},
			name:        "rg",
			want:        f.binRg,
		},
		{
			description: "V2: a leading empty PATH entry never selects the rg planted in the working directory",
			path:        []string{"", f.bin},
			name:        "rg",
			want:        f.binRg,
			notFoundOK:  true,
		},
		{
			description: "V3: a leading dot PATH entry never selects the rg planted in the working directory",
			path:        []string{".", f.bin},
			name:        "rg",
			want:        f.binRg,
			notFoundOK:  true,
		},
		{
			description: "V4: a relative PATH entry never selects the rg reachable through it from the working directory",
			path:        []string{filepath.Join("rel", "sub"), f.bin},
			name:        "rg",
			want:        f.binRg,
			notFoundOK:  true,
		},
		{
			description: "V5: a PATH made only of dot, empty and relative entries finds nothing even though ./rg exists",
			path:        []string{".", "", filepath.Join("rel", "sub")},
			name:        "rg",
			want:        "",
		},
		{
			description: "V6: a regular file without an execute bit is not an executable",
			path:        []string{f.noexec},
			name:        "rg",
			want:        "",
		},
		{
			description: "V7: a directory named rg is skipped and the search continues to the next entry",
			path:        []string{f.dirs, f.bin},
			name:        "rg",
			want:        f.binRg,
		},
		{
			description: "V9: an absolute name is returned as given without consulting PATH",
			path:        nil,
			name:        f.binRg,
			want:        f.binRg,
		},
		{
			description: "V9: an absolute name that does not exist is not found",
			path:        []string{f.bin},
			name:        filepath.Join(f.bin, "missing"),
			want:        "",
		},
		{
			description: "V12: an empty PATH finds nothing rather than falling back to the working directory or a default path",
			path:        nil,
			name:        "rg",
			want:        "",
		},
		{
			description: "P1: a backslash is not a path separator on POSIX, so the name is searched for on PATH like any bare name",
			path:        []string{f.bin},
			name:        `a\b`,
			want:        backslashName,
		},
		{
			description: "G4: an explicit relative name is refused rather than resolved against the working directory (make it absolute first)",
			path:        []string{f.bin},
			name:        "./rg",
			want:        "",
		},
		{
			description: "an absolute entry listed before a dot entry wins, so PATH order alone decides between two legitimate installs",
			path:        []string{f.bin, "."},
			name:        "rg",
			want:        f.binRg,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			t.Setenv("PATH", strings.Join(tc.path, string(os.PathListSeparator)))

			got, err := lookPath(tc.name)
			if got == "" {
				require.Error(t, err, "lookPath returned neither a path nor an error")
				require.True(t, tc.want == "" || tc.notFoundOK, "lookPath reported not-found (%v); want %q", err, tc.want)
				return
			}
			require.NoError(t, err, "lookPath returned a path together with an error")
			require.Equal(t, tc.want, got, "lookPath must never return the planted file or a relative spelling")
			// V10: whatever comes back is absolute and already clean — it is
			// handed to exec.Command as-is.
			require.True(t, filepath.IsAbs(got), "result %q is not absolute", got)
			require.Equal(t, filepath.Clean(got), got, "result %q is not normalised", got)
		})
	}

	t.Run("V11: degenerate names are not found and do not panic", func(t *testing.T) {
		t.Setenv("PATH", f.bin)
		for _, name := range []string{"", ".", ".."} {
			got, err := lookPath(name)
			require.Error(t, err, "name %q", name)
			require.Empty(t, got, "name %q", name)
		}
	})
}

// TestLookPathWithExecerrdotDisabled covers the one configuration in which
// exec.LookPath itself hands back a working-directory-relative hit with no
// error: GODEBUG=execerrdot=0, which restores the pre-Go 1.19 behaviour and
// can be switched on by an environment variable, a go.mod godebug line or a
// //go:debug directive in the embedding program. lookPath must refuse the hit
// regardless, because exec.Command would run it.
func TestLookPathWithExecerrdotDisabled(t *testing.T) {
	f := newLookPathFixture(t)
	t.Setenv("GODEBUG", "execerrdot=0")
	t.Setenv("PATH", "."+string(os.PathListSeparator)+f.bin)

	if raw, err := exec.LookPath("rg"); err != nil || filepath.IsAbs(raw) {
		t.Skipf("this Go release no longer honours GODEBUG=execerrdot=0 (exec.LookPath returned %q, %v); nothing extra to guard", raw, err)
	}

	got, err := lookPath("rg")
	require.Error(t, err)
	require.True(t, errors.Is(err, exec.ErrDot), "want an error matching exec.ErrDot, got %v", err)
	require.Empty(t, got, "lookPath must not return the working-directory hit exec.LookPath let through")

	// An absolute entry that comes first still resolves as usual.
	t.Setenv("PATH", f.bin+string(os.PathListSeparator)+".")
	got, err = lookPath("rg")
	require.NoError(t, err)
	require.Equal(t, f.binRg, got)
}

package agenttoolset

import (
	"os/exec"
	"path/filepath"
)

// lookPath is how this package turns a bare helper-program name ("rg") into
// the path it hands to exec.Command. It wraps [exec.LookPath] and exists so
// the following invariant is explicit, documented and tested here rather than
// implied by standard-library defaults:
//
// Safe executable resolution. This package never launches a helper program by
// bare name. Every program it spawns is either an absolute path it wrote
// itself (/bin/bash) or a bare name resolved by exec.LookPath via lookPath,
// which searches only the absolute entries of PATH and never the current
// working directory — neither implicitly, as Windows CreateProcess and cmd.exe
// do, nor via a ".", empty or relative PATH entry. The absolute path it
// returns is what is handed to the OS. This holds on every platform, so a file
// planted in a directory the process merely works in (a cloned repository, an
// extracted archive) is never selected as a helper binary.
//
// How the guarantees shared with the other Anthropic SDKs map onto Go:
//
//   - G1, absolute argv[0]: lookPath returns an absolute path or an error,
//     nothing in between. exec.LookPath hands back an explicit relative name
//     ("./rg") as given, and under GODEBUG=execerrdot=0 (an environment
//     variable, a go.mod godebug line or a //go:debug directive in the
//     program embedding this package) it hands back relative PATH hits with
//     no error at all; lookPath reports both as [exec.ErrDot] instead.
//   - G2, the working directory is never a search location: since Go 1.19
//     exec.LookPath and exec.Command refuse a match that is only reachable
//     relative to the current directory — Windows' implicit ".", or a "." or
//     relative PATH entry on any platform (an empty entry means "." on Unix
//     and is skipped on Windows) — and report exec.ErrDot, but still return
//     the offending path alongside it. lookPath drops that path.
//     exec.LookPath stops at such a match rather than skipping it, so
//     PATH=".:/usr/bin" plus a planted ./rg yields not-found here (grep then
//     uses its built-in walker) where the sibling SDKs' resolvers skip the
//     entry and go on to find /usr/bin/rg. Either way the planted file is
//     never returned.
//   - G3, Windows runs native executables only: delegated to exec.LookPath,
//     which appends each PATHEXT extension to an extensionless name and never
//     matches the extensionless file itself. Unlike the Python and TypeScript
//     SDKs, which accept only .exe and .com, the default PATHEXT also admits
//     .bat and .cmd from a legitimate PATH directory. We deliberately do not
//     reimplement LookPath to narrow that; the only name resolved today is
//     "rg", with arguments this package controls.
//   - G4, explicit paths are the caller's decision: a name containing a path
//     separator is not searched for. An absolute one is returned as given if
//     it names an executable regular file; a relative one is refused (G1) —
//     make it absolute first if that is really what is meant.
//   - G5, one implementation: every process this package starts goes through
//     lookPath or an absolute path literal (see CONTRIBUTING.md, "Spawning
//     external programs"). Do not call exec.LookPath directly, do not treat
//     exec.ErrDot as success, do not hand-roll a PATH walk, do not build
//     "./"+name.
//   - D1, NoDefaultCurrentDirectoryInExePath: nothing to set. exec.LookPath
//     already honours it on Windows, and refuses working-directory matches
//     whether or not it is set.
//
// Mirrors claude-code's safeExecutableResolver; keep in sync with the exec
// modules of claude-agent-sdk-python, anthropic-sdk-python and
// anthropic-sdk-typescript.
func lookPath(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		// Includes exec.ErrDot, which arrives together with the relative path
		// it complains about. Drop the path so no caller can run it anyway.
		return "", err
	}
	if !filepath.IsAbs(path) {
		// An explicit relative name, or a relative PATH hit let through by
		// GODEBUG=execerrdot=0. exec.Command would resolve either against the
		// working directory, so refuse it like any other dot hit.
		return "", &exec.Error{Name: name, Err: exec.ErrDot}
	}
	return path, nil
}

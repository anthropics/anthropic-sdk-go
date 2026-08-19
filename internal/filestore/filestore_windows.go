//go:build windows

package filestore

import (
	"errors"
	"os"
)

// FileStore's guarantees depend on O_NOFOLLOW semantics that Windows does not
// provide, so rather than degrading to a best-effort implementation, Open
// refuses. This is a deployment condition, not caller input, so it is a plain
// error rather than a *FileStoreError.

func platformCheck() error {
	return errors.New("filestore: requires O_NOFOLLOW support; unsupported on this platform")
}

// replaceViaTemp and openRegularFile are unreachable — Open never hands out a
// store here — but must exist for the package to compile.

func replaceViaTemp(string, []byte, bool) error { return platformCheck() }

func openRegularFile(string, string) (*os.File, error) { return nil, platformCheck() }

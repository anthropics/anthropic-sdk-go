//go:build windows

package filestore

import (
	"strings"
	"testing"
)

func TestFileStoreOpenRefusesUnsupportedPlatform(t *testing.T) {
	store, err := Open(t.TempDir())
	if err == nil || store != nil {
		t.Fatalf("Open = %v, %v; want unsupported-platform error", store, err)
	}
	if !strings.Contains(err.Error(), "unsupported on this platform") {
		t.Fatalf("err = %v", err)
	}
}

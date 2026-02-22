package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// CompareGolden compares the given string with the golden file at path.
// If UPDATE_GOLDEN=1, it overwrites the file with the provided contents.
func CompareGolden(t *testing.T, goldenRelPath string, got []byte) {
	t.Helper()
	goldenPath := filepath.FromSlash(goldenRelPath)

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("creating golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("writing golden: %v", err)
		}
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden: %v", err)
	}
	if string(want) != string(got) {
		t.Fatalf("golden mismatch for %s\nwant:\n%s\n----\ngot:\n%s", goldenRelPath, string(want), string(got))
	}
}

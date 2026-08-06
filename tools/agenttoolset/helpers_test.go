package agenttoolset

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/require"
)

// writeExecutable writes a shell script named name into dir with mode 0o755
// and returns its path — a real executable, so a lookup that wrongly selects
// it (or a spawn that wrongly runs it) behaves exactly as a planted binary
// would.
func writeExecutable(t *testing.T, dir, name, script string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
	require.NoError(t, os.WriteFile(p, []byte("#!/bin/sh\n"+script+"\n"), 0o755))
	return p
}

// runTool executes a BetaTool the way a session/Messages tool runner would and
// flattens the outcome to (text, isError) — the shape the tests assert on.
func runTool(t *testing.T, tool anthropic.BetaTool, raw json.RawMessage) (string, bool) {
	t.Helper()
	out, err := tool.Execute(context.Background(), raw)
	if err != nil {
		return err.Error(), true
	}
	var sb strings.Builder
	for _, b := range out {
		if b.OfText != nil {
			sb.WriteString(b.OfText.Text)
		}
	}
	return sb.String(), false
}

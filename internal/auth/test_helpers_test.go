package auth

import (
	"bytes"
	"sync"
	"testing"
)

// syncBuffer is a goroutine-safe bytes.Buffer for capturing log output.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// resetWarnOnceForTest clears the warnOnce dedupe set so a test can observe
// the next warning, and restores the prior state when the test ends.
func resetWarnOnceForTest(t *testing.T) {
	t.Helper()
	warnOnceMu.Lock()
	prev := warnOnceSet
	warnOnceSet = map[string]bool{}
	warnOnceMu.Unlock()
	t.Cleanup(func() {
		warnOnceMu.Lock()
		warnOnceSet = prev
		warnOnceMu.Unlock()
	})
}

package jsonl

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Close must not panic when the stream was constructed with an error and has no reader.
func TestCloseNilPanic(t *testing.T) {
	stream := NewStream[map[string]any](nil, io.ErrUnexpectedEOF)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Close panicked: %v", r)
		}
	}()
	if err := stream.Close(); err != nil {
		t.Fatalf("Close returned unexpected error: %v", err)
	}
}

// A JSONL line larger than the default 64 KiB scanner buffer must yield one record.
func TestLargeLine(t *testing.T) {
	big := strings.Repeat("x", bufio.MaxScanTokenSize+10)
	body := `{"custom_id":"a","payload":"` + big + `"}` + "\n"
	res := &http.Response{Body: io.NopCloser(strings.NewReader(body))}

	stream := NewStream[map[string]any](res, nil)
	n := 0
	for stream.Next() {
		n++
	}
	if n != 1 {
		t.Fatalf("expected 1 record, got %d", n)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}

// A read error from the body must surface via Err rather than look like clean EOF.
func TestScannerErrorSurfaced(t *testing.T) {
	res := &http.Response{Body: io.NopCloser(&errReader{fail: io.ErrClosedPipe})}
	stream := NewStream[map[string]any](res, nil)
	for stream.Next() {
	}
	if err := stream.Err(); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("Err() = %v, want error matching io.ErrClosedPipe", err)
	}
}

type errReader struct{ fail error }

func (r *errReader) Read(p []byte) (int, error) { return 0, r.fail }

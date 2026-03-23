package jsonl_test

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/packages/jsonl"
)

type testItem struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func TestStream_SmallLines(t *testing.T) {
	body := `{"id":1,"data":"hello"}` + "\n" + `{"id":2,"data":"world"}` + "\n"
	res := &http.Response{
		Body: io.NopCloser(strings.NewReader(body)),
	}

	stream := jsonl.NewStream[testItem](res, nil)
	defer stream.Close()

	count := 0
	for stream.Next() {
		count++
		item := stream.Current()
		if item.ID != count {
			t.Errorf("expected id %d, got %d", count, item.ID)
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 items, got %d", count)
	}
}

func TestStream_LargeLineExceeding64KB(t *testing.T) {
	// Create a JSON line that exceeds the default bufio.Scanner limit (64 KB).
	// Before the fix, this caused the scanner to silently stop with
	// "bufio.Scanner: token too long", making the stream appear empty.
	largeData := strings.Repeat("x", bufio.MaxScanTokenSize+1024) // ~65 KB + 1 KB
	line := fmt.Sprintf(`{"id":1,"data":"%s"}`, largeData)
	body := line + "\n"

	res := &http.Response{
		Body: io.NopCloser(strings.NewReader(body)),
	}

	stream := jsonl.NewStream[testItem](res, nil)
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected stream.Next() to return true for large line, got false; err: %v", stream.Err())
	}
	item := stream.Current()
	if item.ID != 1 {
		t.Fatalf("expected id 1, got %d", item.ID)
	}
	if len(item.Data) != len(largeData) {
		t.Fatalf("expected data length %d, got %d", len(largeData), len(item.Data))
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStream_NilResponse(t *testing.T) {
	stream := jsonl.NewStream[testItem](nil, nil)
	if stream.Next() {
		t.Fatal("expected Next() to return false for nil response")
	}
	if err := stream.Err(); err == nil {
		t.Fatal("expected error for nil response")
	}
}

func TestStream_Error(t *testing.T) {
	stream := jsonl.NewStream[testItem](nil, fmt.Errorf("request failed"))
	if stream.Next() {
		t.Fatal("expected Next() to return false when error is set")
	}
	if err := stream.Err(); err == nil || err.Error() != "request failed" {
		t.Fatalf("expected 'request failed' error, got: %v", err)
	}
}

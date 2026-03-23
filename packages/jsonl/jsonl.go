// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package jsonl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Stream[T any] struct {
	rc  io.ReadCloser
	scn *bufio.Scanner
	cur T
	err error
}

func NewStream[T any](res *http.Response, err error) *Stream[T] {
	if err != nil {
		return &Stream[T]{err: err}
	}

	if res == nil || res.Body == nil {
		return &Stream[T]{err: fmt.Errorf("No streaming response body")}
	}

	scn := bufio.NewScanner(res.Body)
	// Increase the scanner buffer to match the SSE stream decoder (4 MB).
	// The default 64 KB limit silently truncates MessageBatch result lines
	// that exceed it — scn.Scan() returns false and scn.Err() returns
	// "bufio.Scanner: token too long", causing the stream to appear empty.
	scn.Buffer(nil, bufio.MaxScanTokenSize<<9)
	return &Stream[T]{
		rc:  res.Body,
		scn: scn,
		err: err,
	}
}

func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return false
	}

	if !s.scn.Scan() {
		return false
	}

	line := s.scn.Bytes()
	var nxt T
	s.err = json.Unmarshal(line, &nxt)
	s.cur = nxt
	return s.err == nil
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Close() error {
	return s.rc.Close()
}

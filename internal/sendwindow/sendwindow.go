// Package sendwindow carries the tool-result send retry window from
// lib/environments.EnvironmentWorker (which learns the work-item lease TTL from
// each heartbeat) to anthropic.SessionToolRunner (which bounds its send retries
// by it) without putting the plumbing on either package's public API.
package sendwindow

import (
	"context"
	"sync/atomic"
	"time"
)

// Window is a duration that one goroutine updates while another reads it. The
// zero value, and a nil *Window, report 0, meaning "use the default".
type Window struct{ nanos atomic.Int64 }

// Set stores d.
func (w *Window) Set(d time.Duration) { w.nanos.Store(int64(d)) }

// Get returns the last Set value, or 0.
func (w *Window) Get() time.Duration {
	if w == nil {
		return 0
	}
	return time.Duration(w.nanos.Load())
}

type ctxKey struct{}

// NewContext returns a child of ctx carrying w. A SessionToolRunner constructed
// on the returned context re-reads w before every send retry.
func NewContext(ctx context.Context, w *Window) context.Context {
	return context.WithValue(ctx, ctxKey{}, w)
}

// FromContext returns the Window carried by ctx, or nil.
func FromContext(ctx context.Context) *Window {
	w, _ := ctx.Value(ctxKey{}).(*Window)
	return w
}

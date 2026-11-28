package sendwindow_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/sendwindow"
)

func TestWindowZeroValueAndNilReportZero(t *testing.T) {
	var zero sendwindow.Window
	if got := zero.Get(); got != 0 {
		t.Fatalf("zero value Get() = %v, want 0", got)
	}

	var nilWindow *sendwindow.Window
	if got := nilWindow.Get(); got != 0 {
		t.Fatalf("nil Get() = %v, want 0", got)
	}
}

func TestWindowSetReplacesTheLastValue(t *testing.T) {
	var w sendwindow.Window
	w.Set(30 * time.Second)
	if got := w.Get(); got != 30*time.Second {
		t.Fatalf("Get() = %v, want 30s", got)
	}

	w.Set(5 * time.Second)
	if got := w.Get(); got != 5*time.Second {
		t.Fatalf("Get() after second Set = %v, want 5s", got)
	}

	w.Set(0)
	if got := w.Get(); got != 0 {
		t.Fatalf("Get() after Set(0) = %v, want 0", got)
	}
}

func TestContextCarriesTheSameWindow(t *testing.T) {
	w := &sendwindow.Window{}
	ctx := sendwindow.NewContext(context.Background(), w)

	got := sendwindow.FromContext(ctx)
	if got != w {
		t.Fatalf("FromContext returned %p, want the window stored by NewContext (%p)", got, w)
	}

	// A caller that holds the context sees updates made after it was built.
	w.Set(12 * time.Second)
	if d := sendwindow.FromContext(ctx).Get(); d != 12*time.Second {
		t.Fatalf("Get() through the context = %v, want 12s", d)
	}

	// A child context still finds the window.
	child, cancel := context.WithCancel(ctx)
	defer cancel()
	if sendwindow.FromContext(child) != w {
		t.Fatal("FromContext on a child context did not return the window")
	}
}

func TestFromContextWithoutWindowIsNil(t *testing.T) {
	if w := sendwindow.FromContext(context.Background()); w != nil {
		t.Fatalf("FromContext(Background) = %v, want nil", w)
	}
	if d := sendwindow.FromContext(context.Background()).Get(); d != 0 {
		t.Fatalf("Get() on the missing window = %v, want 0", d)
	}
}

func TestWindowConcurrentSetAndGet(t *testing.T) {
	var w sendwindow.Window
	var wg sync.WaitGroup
	for i := 1; i <= 8; i++ {
		wg.Add(2)
		go func(d time.Duration) {
			defer wg.Done()
			w.Set(d * time.Second)
		}(time.Duration(i))
		go func() {
			defer wg.Done()
			if got := w.Get(); got < 0 || got > 8*time.Second {
				t.Errorf("Get() = %v, want a value some Set stored", got)
			}
		}()
	}
	wg.Wait()
}

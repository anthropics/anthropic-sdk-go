package telemetry

import (
	"slices"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// TestMessageDeltaFieldsCoverTheResponseTable asserts no behavior. It compares
// messageDeltaFields, in stream.go, with the generated table of what the span of
// a Message reads from its JSON, and fails when one changes without the other.
func TestMessageDeltaFieldsCoverTheResponseTable(t *testing.T) {
	// The paths of the table whose value message_start already has.
	knownAtMessageStart := []string{
		"id",
		"model",
		"usage.cache_creation_input_tokens",
		"usage.cache_read_input_tokens",
		"usage.cache_creation.ephemeral_5m_input_tokens",
		"usage.cache_creation.ephemeral_1h_input_tokens",
		"usage.service_tier",
	}
	within := func(path, field string) bool { return path == field || strings.HasPrefix(path, field+".") }

	paths := tracing.SpanResponsePaths("anthropic.messages.create")
	if len(paths) == 0 {
		t.Fatal("the span anthropic.messages.create reads nothing from its response: this test and messageDeltaFields in stream.go need its new name")
	}
	for _, path := range paths {
		ends := slices.ContainsFunc(messageDeltaFields, func(field string) bool { return within(path, field) })
		switch known := slices.Contains(knownAtMessageStart, path); {
		case ends && known:
			t.Errorf("%q is in messageDeltaFields in stream.go and in knownAtMessageStart in this test: keep one", path)
		case !ends && !known:
			t.Errorf("the span reads %q from a Message: add it to messageDeltaFields in stream.go when only message_delta ends its value, and to knownAtMessageStart in this test when message_start has it", path)
		}
	}
	for _, field := range messageDeltaFields {
		if !slices.ContainsFunc(paths, func(path string) bool { return within(path, field) }) {
			t.Errorf("the span no longer reads %q from a Message: remove it from messageDeltaFields in stream.go", field)
		}
	}
	for _, path := range knownAtMessageStart {
		if !slices.Contains(paths, path) {
			t.Errorf("the span no longer reads %q from a Message: remove it from knownAtMessageStart in this test", path)
		}
	}
}

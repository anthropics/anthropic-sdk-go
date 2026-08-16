package compat

import (
	"reflect"
	"sort"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// TestAccumulateExhaustive fails when the generated code gains a stream-event or
// delta variant that Accumulate does not handle, so no new variant is dropped
// silently. On failure, handle the named variant in Accumulate, then add it
// to the list below.
func TestAccumulateExhaustive(t *testing.T) {
	events := []string{
		"content_block_delta", "content_block_start", "content_block_stop",
		"message_delta", "message_start", "message_stop",
	}
	deltas := []string{
		"citations_delta", "input_json_delta", "signature_delta",
		"text_delta", "thinking_delta",
	}
	for _, tt := range []struct {
		name    string
		union   any
		handled []string
	}{
		{"events", anthropic.MessageStreamEventUnion{}, events},
		{"deltas", anthropic.RawContentBlockDeltaUnion{}, deltas},
		{"beta events", anthropic.BetaRawMessageStreamEventUnion{}, events},
		{"beta deltas", anthropic.BetaRawContentBlockDeltaUnion{}, append([]string{"compaction_delta"}, deltas...)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			generated := variants(t, tt.union)
			if extra := diff(generated, tt.handled); len(extra) > 0 {
				t.Errorf("new %s variants %q: handle them in Accumulate, then list them here", tt.name, extra)
			}
			if stale := diff(tt.handled, generated); len(stale) > 0 {
				t.Errorf("listed %s variants no longer generated %q: remove them", tt.name, stale)
			}
		})
	}
}

// variants lists the wire discriminant of every variant of a union: the types
// implementing the interface AsAny returns, read from each one's Type field.
func variants(t *testing.T, union any) []string {
	t.Helper()
	typ := reflect.TypeOf(union)
	asAny, ok := typ.MethodByName("AsAny")
	if !ok {
		t.Fatalf("%s has no AsAny", typ.Name())
	}
	variant := asAny.Type.Out(0)
	var out []string
	for i := range typ.NumMethod() {
		ret := typ.Method(i).Type
		if ret.NumOut() != 1 || ret.Out(0) == variant || !ret.Out(0).Implements(variant) {
			continue
		}
		field, ok := ret.Out(0).FieldByName("Type")
		if !ok {
			t.Fatalf("%s: variant %s has no Type field", typ.Name(), ret.Out(0).Name())
		}
		out = append(out, field.Tag.Get("default"))
	}
	return out
}

// diff returns the members of a that b lacks, sorted.
func diff(a, b []string) []string {
	lacks := map[string]bool{}
	for _, s := range b {
		lacks[s] = true
	}
	var out []string
	for _, s := range a {
		if !lacks[s] {
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// TestAccumulateDeltaFieldsExhaustive fails when a message_delta type gains a
// field that Accumulate does not merge into the message. Every value the final
// message_delta carries has to reach the accumulated Message: an ignored one
// leaves the streamed Message differing from the non-streaming one. On failure,
// merge the named field in Accumulate, then add it to the list below.
func TestAccumulateDeltaFieldsExhaustive(t *testing.T) {
	usage := []string{
		"CacheCreationInputTokens", "CacheReadInputTokens", "InputTokens",
		"OutputTokens", "OutputTokensDetails", "ServerToolUse",
	}
	delta := []string{"Container", "StopDetails", "StopReason", "StopSequence"}
	event := []string{"Delta", "Type", "Usage"}
	for _, tt := range []struct {
		name    string
		value   any
		handled []string
	}{
		{"delta usage", anthropic.MessageDeltaUsage{}, usage},
		{"beta delta usage", anthropic.BetaMessageDeltaUsage{}, append([]string{"FallbackCredit", "Iterations"}, usage...)},
		{"delta", anthropic.MessageDeltaEventDelta{}, delta},
		{"beta delta", anthropic.BetaRawMessageDeltaEventDelta{}, delta},
		{"delta event", anthropic.MessageDeltaEvent{}, event},
		{"beta delta event", anthropic.BetaRawMessageDeltaEvent{}, append([]string{"ContextManagement"}, event...)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			generated := fields(tt.value)
			if extra := diff(generated, tt.handled); len(extra) > 0 {
				t.Errorf("new %s fields %q: merge them in Accumulate, then list them here", tt.name, extra)
			}
			if stale := diff(tt.handled, generated); len(stale) > 0 {
				t.Errorf("listed %s fields no longer generated %q: remove them", tt.name, stale)
			}
		})
	}
}

// fields lists the exported fields of a struct, minus the JSON metadata one.
func fields(value any) []string {
	typ := reflect.TypeOf(value)
	var out []string
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() || field.Name == "JSON" {
			continue
		}
		out = append(out, field.Name)
	}
	return out
}

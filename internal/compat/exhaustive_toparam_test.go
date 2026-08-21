package compat

import (
	"bytes"
	"encoding/json"
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// TestToParamExhaustive checks that ToParam keeps all the data of every content
// block the API can return. Blocks are synthesized in synthetic_test.go. A
// response field with no request counterpart is skipped; when codegen adds one,
// convert it in ToParam if a request can carry it, else add it to skipped below.
func TestToParamExhaustive(t *testing.T) {
	skipped := []string{
		"CitationCharLocation.file_id", "CitationContentBlockLocation.file_id", "CitationPageLocation.file_id",
	}
	betaSkipped := []string{
		"BetaCitationCharLocation.file_id", "BetaCitationContentBlockLocation.file_id", "BetaCitationPageLocation.file_id",
	}
	for _, tt := range []struct {
		name     string
		response reflect.Type
		request  reflect.Type
		skipped  []string
		toParam  func([]byte) (any, error)
	}{
		{
			"content blocks",
			reflect.TypeFor[anthropic.ContentBlockUnion](),
			reflect.TypeFor[anthropic.ContentBlockParamUnion](),
			skipped,
			func(wire []byte) (any, error) {
				var block anthropic.ContentBlockUnion
				if err := json.Unmarshal(wire, &block); err != nil {
					return nil, err
				}
				return block.ToParam(), nil
			},
		},
		{
			"beta content blocks",
			reflect.TypeFor[anthropic.BetaContentBlockUnion](),
			reflect.TypeFor[anthropic.BetaContentBlockParamUnion](),
			betaSkipped,
			func(wire []byte) (any, error) {
				var block anthropic.BetaContentBlockUnion
				if err := json.Unmarshal(wire, &block); err != nil {
					return nil, err
				}
				return block.ToParam(), nil
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			synth := newSynthesizer(t, tt.request)
			for _, v := range unionVariants(tt.response) {
				for _, s := range synth.responseShapes(v.Type) {
					t.Run(v.Type.Name()+"/"+s.Name, func(t *testing.T) {
						wire, err := json.Marshal(s.Value)
						if err != nil {
							t.Fatalf("Failed to marshal shape: %v", err)
						}
						param, err := tt.toParam(wire)
						if err != nil {
							t.Fatalf("Failed to unmarshal %s: %v", wire, err)
						}
						want := comparableJSON(t, wire)
						if got := requestJSON(t, param); !bytes.Equal(got, want) {
							t.Errorf("ToParam changed the block\n want: %s\n  got: %s", want, got)
						}
					})
				}
			}
			got := slices.Collect(maps.Keys(synth.skipped))
			if extra := diff(got, tt.skipped); len(extra) > 0 {
				t.Errorf("new response fields with no request counterpart %q: list them here if no request can carry them", extra)
			}
			if stale := diff(tt.skipped, got); len(stale) > 0 {
				t.Errorf("listed fields no longer skipped %q: remove them", stale)
			}
		})
	}
}

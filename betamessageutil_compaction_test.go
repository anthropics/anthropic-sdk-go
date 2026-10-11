package anthropic_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/tidwall/gjson"
)

func accumulateCompactionEvent(t *testing.T, message *anthropic.BetaMessage, raw string) {
	t.Helper()
	var event anthropic.BetaRawMessageStreamEventUnion
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatal(err)
	}
	if err := message.Accumulate(event); err != nil {
		t.Fatal(err)
	}
}

func TestBetaCompactionDeltaKeepsPresenceAndClearsValues(t *testing.T) {
	const start = `{"type":"compaction","content":"previous summary","encrypted_content":"previous metadata","signature":"signature","tool_changes":[]}`
	cases := []struct{ name, delta, content, encrypted string }{
		{"replace both", `"content":"new summary","encrypted_content":"new metadata"`, `"new summary"`, `"new metadata"`},
		{"empty both", `"content":"","encrypted_content":""`, `""`, `""`},
		{"null both", `"content":null,"encrypted_content":null`, `null`, `null`},
		{"content only", `"content":"new summary"`, `"new summary"`, `"previous metadata"`},
		{"metadata only", `"encrypted_content":"new metadata"`, `"previous summary"`, `"new metadata"`},
		{"null content", `"content":null`, `null`, `"previous metadata"`},
		{"empty metadata", `"encrypted_content":""`, `"previous summary"`, `""`},
	}
	for _, test := range cases {
		for _, blockStop := range []bool{false, true} {
			name := test.name
			if blockStop {
				name += "/block-stop"
			} else {
				name += "/message-stop-only"
			}
			t.Run(name, func(t *testing.T) {
				var message anthropic.BetaMessage
				accumulateCompactionEvent(t, &message, `{"type":"message_start","message":{"id":"msg_test","type":"message","role":"assistant","model":"test","content":[],"usage":{"input_tokens":1,"output_tokens":0}}}`)
				accumulateCompactionEvent(t, &message, `{"type":"content_block_start","index":0,"content_block":`+start+`}`)
				accumulateCompactionEvent(t, &message, `{"type":"content_block_delta","index":0,"delta":{"type":"compaction_delta",`+test.delta+`}}`)
				if blockStop {
					accumulateCompactionEvent(t, &message, `{"type":"content_block_stop","index":0}`)
				}
				accumulateCompactionEvent(t, &message, `{"type":"message_stop"}`)
				block := message.Content[0]
				for field, want := range map[string]string{"content": test.content, "encrypted_content": test.encrypted} {
					if got := gjson.Get(block.RawJSON(), field).Raw; got != want {
						t.Errorf("block %s=%s, want %s", field, got, want)
					}
					if got := gjson.Get(message.RawJSON(), "content.0."+field).Raw; got != want {
						t.Errorf("message %s=%s, want %s", field, got, want)
					}
				}
				replay, err := json.Marshal(message.ToParam())
				if err != nil {
					t.Fatal(err)
				}
				for field, want := range map[string]string{"content": test.content, "encrypted_content": test.encrypted, "signature": `"signature"`, "tool_changes": `[]`} {
					if got := gjson.GetBytes(replay, "content.0."+field).Raw; got != want {
						t.Errorf("replayed %s=%s, want %s in %s", field, got, want, replay)
					}
				}
				if block.Content.OfString != gjson.Parse(test.content).String() || block.EncryptedContent != gjson.Parse(test.encrypted).String() {
					t.Errorf("typed fields disagree with deltas: %+v", block)
				}
				variant := block.AsCompaction()
				if variant.Content != block.Content.OfString || variant.EncryptedContent != block.EncryptedContent {
					t.Fatal("variant and accumulated union disagree")
				}
			})
		}
	}
}

func TestBetaCompactionDeltaDoesNotResurrectClearedMetadata(t *testing.T) {
	for _, cleared := range []string{`null`, `""`} {
		t.Run(cleared, func(t *testing.T) {
			var message anthropic.BetaMessage
			accumulateCompactionEvent(t, &message, `{"type":"message_start","message":{"content":[]}}`)
			accumulateCompactionEvent(t, &message, `{"type":"content_block_start","index":0,"content_block":{"type":"compaction","content":"initial","encrypted_content":"old"}}`)
			accumulateCompactionEvent(t, &message, `{"type":"content_block_delta","index":0,"delta":{"type":"compaction_delta","content":`+cleared+`,"encrypted_content":`+cleared+`}}`)
			accumulateCompactionEvent(t, &message, `{"type":"content_block_delta","index":0,"delta":{"type":"compaction_delta","content":"replacement"}}`)
			accumulateCompactionEvent(t, &message, `{"type":"message_stop"}`)
			replay, err := json.Marshal(message.ToParam())
			if err != nil {
				t.Fatal(err)
			}
			if gjson.GetBytes(replay, "content.0.content").String() != "replacement" || gjson.GetBytes(replay, "content.0.encrypted_content").Raw != cleared {
				t.Fatalf("cleared metadata resurrected: %s", replay)
			}
		})
	}
}

func TestBetaCompactionWithoutDeltasRetainsWireFields(t *testing.T) {
	for _, body := range []string{`"content":null,"encrypted_content":null`, `"content":"","encrypted_content":""`, `"content":"summary","encrypted_content":"opaque"`} {
		t.Run(body, func(t *testing.T) {
			raw := `{"type":"compaction",` + body + `,"future_field":{"value":9007199254740993}}`
			var message anthropic.BetaMessage
			accumulateCompactionEvent(t, &message, `{"type":"message_start","message":{"content":[]}}`)
			accumulateCompactionEvent(t, &message, `{"type":"content_block_start","index":0,"content_block":`+raw+`}`)
			accumulateCompactionEvent(t, &message, `{"type":"message_stop"}`)
			actual := message.Content[0].RawJSON()
			if !strings.Contains(actual, `9007199254740993`) {
				t.Fatal("unknown metadata was lost")
			}
			for _, field := range []string{"content", "encrypted_content", "future_field"} {
				if gjson.Get(actual, field).Raw != gjson.Get(raw, field).Raw {
					t.Fatalf("untouched field %s changed: %s", field, actual)
				}
			}
		})
	}
}

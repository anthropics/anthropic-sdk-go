package anthropic

import (
	"testing"
)

func TestThinkingAndSignaturePreserved(t *testing.T) {
	// Create a simple Message with Content containing Thinking and Signature
	msg := &Message{
		Role: MessageRoleAssistant,
		Content: []ContentBlock{
			{
				Type:      ContentBlockTypeText,
				Text:      "Hello world",
				Thinking:  "This is a thinking block",
				Signature: "This is a signature",
			},
		},
	}

	// Convert to MessageParam
	param := msg.ToParam()

	// Check if Thinking and Signature fields are preserved
	contentBlockParam := param.Content.Value[0].(ContentBlockParam)

	if contentBlockParam.Thinking.Value != "This is a thinking block" {
		t.Errorf("Expected Thinking value to be 'This is a thinking block', got %q", contentBlockParam.Thinking.Value)
	}

	if contentBlockParam.Signature.Value != "This is a signature" {
		t.Errorf("Expected Signature value to be 'This is a signature', got %q", contentBlockParam.Signature.Value)
	}
}

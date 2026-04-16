package anthropic

import (
	"fmt"
	"os"
)

// modelsToWarnWithThinkingEnabled lists models for which `thinking.type=enabled`
// (i.e. budget_tokens-based extended thinking) is deprecated in favor of
// `thinking.type=adaptive`.
var modelsToWarnWithThinkingEnabled = map[string]bool{
	"claude-opus-4-6":       true,
	"claude-mythos-preview": true,
}

// warnIfThinkingEnabled prints a deprecation warning to stderr when a request
// uses `thinking.type=enabled` with a model that supports adaptive thinking.
// This matches the runtime warning emitted by the other Anthropic SDKs.
func warnIfThinkingEnabled(model Model, thinkingEnabled bool) {
	if !thinkingEnabled {
		return
	}
	if !modelsToWarnWithThinkingEnabled[string(model)] {
		return
	}
	fmt.Fprintf(
		os.Stderr,
		"Warning: Using Claude with %s and 'thinking.type=enabled' is deprecated. Use 'thinking.type=adaptive' instead which results in better model performance in our testing: https://platform.claude.com/docs/en/build-with-claude/adaptive-thinking\n",
		model,
	)
}

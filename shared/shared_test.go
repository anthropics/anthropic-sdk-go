// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared_test

import (
	"github.com/anthropics/anthropic-sdk-go"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

func TestCalculateNonStreamingTimeout(t *testing.T) {
	// Store original model token limits to restore after test
	originalModelTokenLimits := make(map[string]int)
	for k, v := range constant.ModelNonStreamingTokens {
		originalModelTokenLimits[k] = v
	}
	defer func() {
		// Restore original model token limits
		constant.ModelNonStreamingTokens = originalModelTokenLimits
	}()

	// Set up a test model for consistent testing
	constant.ModelNonStreamingTokens = map[string]int{
		"test-model": 8192,
	}

	tests := []struct {
		name          string
		maxTokens     int
		model         string
		expectTimeout time.Duration
		expectError   bool
	}{
		{
			name:          "small token count returns expected timeout",
			maxTokens:     1000,
			model:         "any-model",
			expectTimeout: time.Duration(float64(60*60) * float64(1000) / 128000.0 * float64(time.Second)),
			expectError:   false,
		},
		{
			name:          "large token count above default time limit throws error",
			maxTokens:     100000,
			model:         "any-model",
			expectTimeout: 0,
			expectError:   true,
		},
		{
			name:          "token count above model specific limit throws error",
			maxTokens:     9000,
			model:         "test-model",
			expectTimeout: 0,
			expectError:   true,
		},
		{
			name:          "token count below model specific limit is ok",
			maxTokens:     8000,
			model:         "test-model",
			expectTimeout: time.Duration(float64(60*60) * float64(8000) / 128000.0 * float64(time.Second)),
			expectError:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			timeout, err := anthropic.CalculateNonStreamingTimeout(tc.maxTokens, tc.model)

			if tc.expectError && err == nil {
				t.Error("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}

			if timeout != tc.expectTimeout {
				t.Errorf("Expected timeout %v but got %v", tc.expectTimeout, timeout)
			}
		})
	}
}

// Test specific model limits
func TestModelLimits(t *testing.T) {
	// Verify the model limits are defined for opus-4 models
	if _, exists := constant.ModelNonStreamingTokens["claude-opus-4-20250514"]; !exists {
		t.Error("Expected model limit for claude-opus-4-20250514 but not found")
	}

	if _, exists := constant.ModelNonStreamingTokens["anthropic.claude-opus-4-20250514-v1:0"]; !exists {
		t.Error("Expected model limit for anthropic.claude-opus-4-20250514-v1:0 but not found")
	}

	if _, exists := constant.ModelNonStreamingTokens["claude-opus-4@20250514"]; !exists {
		t.Error("Expected model limit for claude-opus-4@20250514 but not found")
	}
}

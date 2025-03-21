// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestUsage(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.ContentBlockParamUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText), CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{Type: anthropic.F(anthropic.CacheControlEphemeralTypeEphemeral)}), Citations: anthropic.F([]anthropic.TextCitationParamUnion{anthropic.CitationCharLocationParam{CitedText: anthropic.F("cited_text"), DocumentIndex: anthropic.F(int64(0)), DocumentTitle: anthropic.F("x"), EndCharIndex: anthropic.F(int64(0)), StartCharIndex: anthropic.F(int64(0)), Type: anthropic.F(anthropic.CitationCharLocationParamTypeCharLocation)}})}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude3_7SonnetLatest),
	})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v\n", message.Content)
}

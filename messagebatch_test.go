// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestMessageBatchNew(t *testing.T) {
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
	_, err := client.Messages.Batches.New(context.TODO(), anthropic.MessageBatchNewParams{
		Requests: anthropic.F([]anthropic.MessageBatchNewParamsRequest{{
			CustomID: anthropic.F("my-custom-id-1"),
			Params: anthropic.F(anthropic.MessageBatchNewParamsRequestsParams{
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.MessageParam{{
					Content: anthropic.F([]anthropic.ContentBlockParamUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText), CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{Type: anthropic.F(anthropic.CacheControlEphemeralTypeEphemeral)})}}),
					Role:    anthropic.F(anthropic.MessageParamRoleUser),
				}}),
				Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
				Metadata: anthropic.F(anthropic.MetadataParam{
					UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				}),
				StopSequences: anthropic.F([]string{"string"}),
				Stream:        anthropic.F(true),
				System:        anthropic.F([]anthropic.TextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.TextBlockParamTypeText), CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{Type: anthropic.F(anthropic.CacheControlEphemeralTypeEphemeral)})}}),
				Temperature:   anthropic.F(1.000000),
				ToolChoice: anthropic.F[anthropic.ToolChoiceUnionParam](anthropic.ToolChoiceAutoParam{
					Type:                   anthropic.F(anthropic.ToolChoiceAutoTypeAuto),
					DisableParallelToolUse: anthropic.F(true),
				}),
				Tools: anthropic.F([]anthropic.ToolParam{{
					InputSchema: anthropic.F(anthropic.ToolInputSchemaParam{
						Type: anthropic.F(anthropic.ToolInputSchemaTypeObject),
						Properties: anthropic.F[any](map[string]interface{}{
							"location": map[string]interface{}{
								"description": "The city and state, e.g. San Francisco, CA",
								"type":        "string",
							},
							"unit": map[string]interface{}{
								"description": "Unit for the output - one of (celsius, fahrenheit)",
								"type":        "string",
							},
						}),
					}),
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{
						Type: anthropic.F(anthropic.CacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
				}}),
				TopK: anthropic.F(int64(5)),
				TopP: anthropic.F(0.700000),
			}),
		}}),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchGet(t *testing.T) {
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
	_, err := client.Messages.Batches.Get(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchListWithOptionalParams(t *testing.T) {
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
	_, err := client.Messages.Batches.List(context.TODO(), anthropic.MessageBatchListParams{
		AfterID:  anthropic.F("after_id"),
		BeforeID: anthropic.F("before_id"),
		Limit:    anthropic.F(int64(1)),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchCancel(t *testing.T) {
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
	_, err := client.Messages.Batches.Cancel(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("abc"))
	}))
	defer server.Close()
	baseURL := server.URL
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	resp, err := client.Messages.Batches.Results(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if !bytes.Equal(b, []byte("abc")) {
		t.Fatalf("return value not %s: %s", "abc", b)
	}
}

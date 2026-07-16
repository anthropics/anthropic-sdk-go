package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/googlecloud"
)

// Claude Platform on Google Cloud: the full first-party Anthropic API served
// through the Google Cloud gateway. Authentication uses Google Application
// Default Credentials (run `gcloud auth application-default login`, or rely on
// the ambient service-account credentials in a GCP environment).
//
// Configure via environment, or set the fields directly:
//
//	ANTHROPIC_GOOGLE_CLOUD_PROJECT       GCP consumer project (or inferred from ADC)
//	ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID  Anthropic workspace ID for the org
//	ANTHROPIC_GOOGLE_CLOUD_LOCATION      override the GCP location (optional; defaults to "global")
//	ANTHROPIC_GOOGLE_CLOUD_BASE_URL      override the gateway base URL (optional)
func main() {
	// NewClient's ctx is used for credential discovery during construction
	// only; it is not retained for the client's lifetime.
	client, err := googlecloud.NewClient(context.Background(), googlecloud.ClientConfig{
		Project:     os.Getenv("ANTHROPIC_GOOGLE_CLOUD_PROJECT"),
		WorkspaceID: os.Getenv("ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID"),
	})
	if err != nil {
		panic(err)
	}

	content := "Write me a haiku about Go."
	fmt.Println("[user]: " + content)

	// Use a per-request context with a timeout for individual API calls.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model: "claude-sonnet-4-5",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("[assistant]: " + message.Content[0].Text)
}

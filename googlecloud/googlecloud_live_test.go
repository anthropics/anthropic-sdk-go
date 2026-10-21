package googlecloud_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/googlecloud"
)

// Live integration test for the Google Cloud client. Skipped unless ANTHROPIC_LIVE=1.
//
// Authentication uses Application Default Credentials (cloud-platform scope), e.g.
// after `gcloud auth application-default login`. Required env vars:
//
//	ANTHROPIC_GOOGLE_CLOUD_PROJECT       GCP consumer project (or inferred from ADC)
//	ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID  Anthropic workspace ID for the org
//	ANTHROPIC_GOOGLE_CLOUD_LOCATION      GCP location (test override; default us-central1)
//
// Run: ANTHROPIC_LIVE=1 go test ./googlecloud/... -run TestLiveGoogleCloud -v
func TestLiveGoogleCloud(t *testing.T) {
	if os.Getenv("ANTHROPIC_LIVE") != "1" {
		t.Skip("set ANTHROPIC_LIVE=1 to run live integration tests")
	}
	for _, name := range []string{"ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID"} {
		if os.Getenv(name) == "" {
			t.Fatalf("required env var %s is not set", name)
		}
	}

	location := os.Getenv("ANTHROPIC_GOOGLE_CLOUD_LOCATION")
	if location == "" {
		location = "us-central1"
	}
	model := os.Getenv("ANTHROPIC_LIVE_MODEL")
	if model == "" {
		model = "claude-haiku-4-5"
	}

	ctx := context.Background()
	client, err := googlecloud.NewClient(ctx, googlecloud.ClientConfig{
		Location: location,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	t.Run("messages", func(t *testing.T) {
		msg, err := client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     anthropic.Model(model),
			MaxTokens: 64,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock("Say hello in exactly three words.")),
			},
		})
		if err != nil {
			t.Fatalf("Messages.New failed: %v", err)
		}
		if len(msg.Content) == 0 {
			t.Fatal("expected non-empty response content")
		}
		t.Logf("response: %s", msg.Content[0].Text)
	})

	t.Run("streaming", func(t *testing.T) {
		stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
			Model:     anthropic.Model(model),
			MaxTokens: 64,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock("Count to three.")),
			},
		})
		defer stream.Close()

		var acc anthropic.Message
		var sawStart, sawStop bool
		for stream.Next() {
			ev := stream.Current()
			if err := acc.Accumulate(ev); err != nil {
				t.Fatalf("Accumulate failed: %v", err)
			}
			switch ev.AsAny().(type) {
			case anthropic.MessageStartEvent:
				sawStart = true
			case anthropic.MessageStopEvent:
				sawStop = true
			}
		}
		if err := stream.Err(); err != nil {
			t.Fatalf("stream error: %v", err)
		}
		if !sawStart || !sawStop {
			t.Fatalf("incomplete event sequence: start=%v stop=%v", sawStart, sawStop)
		}
		if len(acc.Content) == 0 || acc.Content[0].Text == "" {
			t.Fatal("expected accumulated text content")
		}
		t.Logf("streamed: %s", acc.Content[0].Text)
	})

	t.Run("error", func(t *testing.T) {
		_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
			Model:     "no-such-model",
			MaxTokens: 1,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
			},
		})
		if err == nil {
			t.Fatal("expected error for unknown model")
		}
		var apierr *anthropic.Error
		if !errors.As(err, &apierr) {
			t.Fatalf("error is not *anthropic.Error: %T %v", err, err)
		}
		if apierr.StatusCode < 400 || apierr.StatusCode >= 500 {
			t.Errorf("status = %d, want 4xx", apierr.StatusCode)
		}
		t.Logf("error: status=%d %v", apierr.StatusCode, apierr)
	})

	t.Run("organization", func(t *testing.T) {
		workspaceID := os.Getenv("ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID")
		workspace, err := client.Organization.Workspaces.Get(ctx, workspaceID)
		if err != nil {
			t.Fatalf("Organization.Workspaces.Get failed: %v", err)
		}
		if workspace.ID != workspaceID {
			t.Errorf("workspace ID = %q, want %q", workspace.ID, workspaceID)
		}
	})

	t.Run("managed_agents", func(t *testing.T) {
		agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
			Model: anthropic.BetaManagedAgentsModelConfigParams{ID: model},
			Name:  "sdk-live-test",
		})
		if err != nil {
			t.Fatalf("Beta.Agents.New failed: %v", err)
		}
		t.Cleanup(func() {
			archived, err := client.Beta.Agents.Archive(ctx, agent.ID, anthropic.BetaAgentArchiveParams{})
			if err != nil {
				t.Fatalf("Beta.Agents.Archive failed: %v", err)
			}
			if archived.ArchivedAt.IsZero() {
				t.Error("expected archived agent to have archived_at set")
			}
		})

		retrieved, err := client.Beta.Agents.Get(ctx, agent.ID, anthropic.BetaAgentGetParams{})
		if err != nil {
			t.Fatalf("Beta.Agents.Get failed: %v", err)
		}
		if retrieved.ID != agent.ID {
			t.Errorf("agent ID = %q, want %q", retrieved.ID, agent.ID)
		}

		page, err := client.Beta.Agents.List(ctx, anthropic.BetaAgentListParams{
			Limit: anthropic.Int(1),
		})
		if err != nil {
			t.Fatalf("Beta.Agents.List failed: %v", err)
		}
		if len(page.Data) == 0 {
			t.Error("expected at least one agent")
		}
	})
}

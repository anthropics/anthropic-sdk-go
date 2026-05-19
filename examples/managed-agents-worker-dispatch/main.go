// Self-hosted runner, "worker-dispatch" flavor: this process was handed ONE
// already-claimed work item by an upstream poller/orchestrator — e.g. an
// `ant worker poll --on-work <this binary>` loop, or your own dispatcher that
// spawns a fresh sandbox per work item. It does NOT create an agent or session,
// and it does NOT poll for work: something else did that and claimed the item.
//
// EnvironmentWorker.HandleItem with an empty HandleItemOptions reads the claimed
// item's identity from the environment variables the upstream poller sets:
//
//	ANTHROPIC_WORK_ID         - the claimed work item to serve
//	ANTHROPIC_ENVIRONMENT_ID  - the self-hosted environment it belongs to
//	ANTHROPIC_SESSION_ID      - the session to run tools for
//	ANTHROPIC_ENVIRONMENT_KEY - the environment key (the runner's single credential)
//
// It then sets up the workdir + downloads the session agent's skills, runs the
// session's tools while heartbeating the work-item lease, force-stops the item
// on exit, and returns — one item, then this process exits.
//
// Also required:
//
//	ANTHROPIC_API_KEY - your API key (read by the SDK client)
//
// Security model: the worker executes bash and file operations directly on the
// host. This is the "sandbox process" shape — the upstream orchestrator is
// expected to have spawned it inside a container or other isolation boundary.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/environments"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
)

// currentTimeTool is a custom anthropic.BetaTool that returns the local time.
// It demonstrates extending the default tool list alongside
// agent_toolset_20260401. The worker executes it whenever the session emits a
// matching agent.custom_tool_use event.
type currentTimeTool struct{}

func (currentTimeTool) Name() string        { return "current_time" }
func (currentTimeTool) Description() string { return "Get the current local time on the worker host." }
func (currentTimeTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{Properties: map[string]any{}}
}
func (currentTimeTool) Execute(context.Context, json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	return []anthropic.BetaToolResultBlockParamContentUnion{
		{OfText: &anthropic.BetaTextBlockParam{Text: time.Now().Format(time.RFC3339)}},
	}, nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	client := anthropic.NewClient()
	ctx := context.Background()

	// Base directory for the per-session agenttoolset.AgentToolContext. An
	// orchestrator typically points this at the sandbox's scratch space.
	workdir := os.Getenv("ANTHROPIC_WORKDIR")
	if workdir == "" {
		workdir = "."
	}

	// Build the worker with a tools factory — the standard agent_toolset_20260401
	// set plus our one custom tool — and nothing else. No EnvironmentID /
	// EnvironmentKey here: HandleItem resolves the work item (and the environment
	// key) from the ANTHROPIC_* env vars the upstream poller set.
	worker := environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{
		Workdir: workdir,
		Logger:  logger,
		ToolsFunc: func(env *agenttoolset.AgentToolContext) []anthropic.BetaTool {
			return append(agenttoolset.BetaAgentToolset20260401(env), currentTimeTool{})
		},
	})

	// Service the single claimed item to completion: set up the workdir +
	// download the session agent's skills, run the local tools against the
	// session's agent.tool_use / agent.custom_tool_use events while heartbeating
	// the lease, then force-stop the work item. HandleItem with an empty
	// HandleItemOptions reads ANTHROPIC_WORK_ID / ANTHROPIC_ENVIRONMENT_ID /
	// ANTHROPIC_SESSION_ID / ANTHROPIC_ENVIRONMENT_KEY from the environment.
	if err := worker.HandleItem(ctx, environments.HandleItemOptions{}); err != nil {
		logger.Error("handle item failed", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("work item handled")
}

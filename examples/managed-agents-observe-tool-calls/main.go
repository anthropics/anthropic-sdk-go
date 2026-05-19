// Self-hosted runner, "observe every tool call" flavor: the low-level
// client.Beta.Sessions.Events.NewToolRunner path. The SessionToolRunner it
// returns dispatches a session's agent.tool_use / agent.custom_tool_use events
// to your local tools, posts each result back, and yields one
// anthropic.DispatchedToolCall per completed call — so you can watch every
// dispatch (name, input, error flag, whether the result posted). Unlike
// environments.EnvironmentWorker, it does NOT poll for work and does NOT manage
// a work-item lease.
//
// Two scenarios, two functions in this file:
//
//   - main() — PRIMARY. A session you created and drive yourself: no work
//     queue, no lease. The runner just dispatches tools against the session's
//     events, so it works the same whether or not the session's environment is
//     self-hosted. Reach for this when you want per-call visibility on a
//     session you own.
//
//   - observeAsSelfHostedWorker() — SECONDARY (not called by default). If you
//     ARE a self-hosted worker but want per-call visibility, you have to
//     compose the pieces EnvironmentWorker would otherwise compose for you: the
//     work poller, the per-session agent tool context, AND your own lease
//     heartbeat running in parallel with the runner loop. Reach for
//     EnvironmentWorker instead unless you specifically need to see each call.
//
// Required environment variables:
//
//	ANTHROPIC_API_KEY         - your API key (read by the SDK client)
//	ANTHROPIC_ENVIRONMENT_ID  - the self-hosted environment the session runs in
//	ANTHROPIC_ENVIRONMENT_KEY - the environment's key, the runner's single
//	                            credential (only the secondary scenario needs it)
//
// Security model: the tools execute bash and file operations directly on the
// host. Run inside a container or other isolation boundary you control.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/environments"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
)

// currentTimeTool is a custom anthropic.BetaTool that returns the local time.
// It demonstrates extending the default tool list alongside
// agent_toolset_20260401. The runner executes it whenever the session emits a
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

// workdir is the base directory for the per-session agenttoolset.AgentToolContext
// — the directory the file tools confine to and where SetupSkills downloads the
// session agent's skills.
func workdir() string {
	if w := os.Getenv("ANTHROPIC_WORKDIR"); w != "" {
		return w
	}
	return "."
}

// ===== PRIMARY: observe a session you created and drive yourself =====

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	client := anthropic.NewClient()
	ctx := context.Background()

	envID := mustEnv("ANTHROPIC_ENVIRONMENT_ID")

	// 1. Create an agent that exposes both the default toolset and our custom
	//    tool, then a session against the self-hosted environment.
	agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name:   "observe-tool-calls-example",
		Model:  anthropic.BetaManagedAgentsModelConfigParams{ID: "claude-haiku-4-5"},
		System: param.NewOpt("You are running in a sandbox. Use the available tools to answer."),
		Tools: []anthropic.BetaAgentNewParamsToolUnion{
			{OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
				Type: "agent_toolset_20260401",
			}},
			{OfCustom: &anthropic.BetaManagedAgentsCustomToolParams{
				Type:        "custom",
				Name:        "current_time",
				Description: "Get the current local time on the worker host.",
				InputSchema: anthropic.BetaManagedAgentsCustomToolInputSchemaParam{
					Type:       "object",
					Properties: map[string]any{},
				},
			}},
		},
	})
	if err != nil {
		fatal(logger, "create agent", err)
	}
	logger.Info("created agent", slog.String("agent_id", agent.ID))
	defer func() {
		// Run cleanup with a fresh-but-deadlined context so a cancelled parent
		// ctx doesn't skip Archive.
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if _, err := client.Beta.Agents.Archive(cleanup, agent.ID, anthropic.BetaAgentArchiveParams{}); err != nil {
			logger.Warn("archive agent failed", slog.Any("error", err))
		}
	}()

	session, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{
		Agent:         anthropic.BetaSessionNewParamsAgentUnion{OfString: param.NewOpt(agent.ID)},
		EnvironmentID: envID,
		Title:         param.NewOpt("observe-tool-calls-example"),
	})
	if err != nil {
		fatal(logger, "create session", err)
	}
	logger.Info("created session", slog.String("session_id", session.ID))
	defer func() {
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if _, err := client.Beta.Sessions.Delete(cleanup, session.ID, anthropic.BetaSessionDeleteParams{}); err != nil {
			logger.Warn("delete session failed", slog.Any("error", err))
		}
	}()

	// 2. Build the per-session agent tool context: the workdir the file tools
	//    confine to, plus the skills SetupSkills downloads into {workdir}/skills/.
	//    Cleanup removes them again.
	env := &agenttoolset.AgentToolContext{Workdir: workdir()}
	if err := env.SetupSkills(ctx, client, session.ID); err != nil {
		logger.Warn("skill setup failed", slog.Any("error", err))
	}
	defer func() {
		if err := env.Cleanup(); err != nil {
			logger.Warn("skill cleanup failed", slog.Any("error", err))
		}
	}()

	tools := append(agenttoolset.BetaAgentToolset20260401(env), currentTimeTool{})
	defer agenttoolset.CloseAll(tools)

	// TODO(codegen): the autogenerated
	// BetaManagedAgentsEventParamsOfUserMessage helper omits the required
	// `type` field — server validators reject the payload with
	// "events[0].type: Field required". Workaround: construct the union
	// directly with Type set. Remove once the generator emits the assignment.
	_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
		Events: []anthropic.BetaManagedAgentsEventParamsUnion{
			{OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
				Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
				Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{
					{OfText: &anthropic.BetaManagedAgentsTextBlockParam{
						Type: anthropic.BetaManagedAgentsTextBlockTypeText,
						Text: "What is the current time? Also run pwd to show me the working directory.",
					}},
				},
			}},
		},
	})
	if err != nil {
		fatal(logger, "send user message", err)
	}

	// 3. Iterate the SessionToolRunner: it attaches to the session, runs each
	//    tool call locally, posts the result back, and yields one
	//    DispatchedToolCall per completed call. The runner stops on its own once
	//    the session goes idle (MaxIdle after an end_turn); the ctx timeout is
	//    just a hard cap for the demo. It does NOT touch any work-item lease.
	//    No RequestOptions: a session you own is driven with the SDK client's
	//    own ANTHROPIC_API_KEY auth, not an environment key.
	runCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	maxIdle := 10 * time.Second
	runner := client.Beta.Sessions.Events.NewToolRunner(runCtx, session.ID, anthropic.SessionToolRunnerOptions{
		Tools:   tools,
		MaxIdle: &maxIdle,
		Logger:  logger,
	})
	defer runner.Close()

	fmt.Println("\n--- tool calls ---")
	for runner.Next() {
		printCall(runner.Current())
	}
	if err := runner.Err(); err != nil &&
		!errors.Is(err, anthropic.ErrSessionTerminated) &&
		!errors.Is(err, anthropic.ErrIdleTimeout) {
		fatal(logger, "tool runner", err)
	}
}

// ===== SECONDARY: observe each call while ALSO being a self-hosted worker =====

// observeAsSelfHostedWorker is NOT called by main. It is the shape you reach for
// only if you are a self-hosted worker AND you want per-call visibility.
//
// IMPORTANT: the SessionToolRunner does NOT manage the work-item lease —
// environments.EnvironmentWorker is what normally does. EnvironmentWorker polls
// for work, runs the equivalent of the runner loop, AND heartbeats the lease
// (force-stopping on exit), all composed together. Drop down to the
// SessionToolRunner for per-call visibility and you give up that lease
// management — so you have to roll it back yourself: the heartbeat goroutine
// below runs in parallel with the runner loop for exactly that reason. It is a
// SIMPLIFIED shape (fixed interval, minimal error handling);
// EnvironmentWorker's internal heartbeat loop is the careful reference — it
// adapts the interval to the server's ttl_seconds and tolerates transient
// failures with backoff. Rolling your own heartbeat is the cost of getting
// per-call visibility AND lease management together.
func observeAsSelfHostedWorker(ctx context.Context, client anthropic.Client, logger *slog.Logger) {
	envID := mustEnv("ANTHROPIC_ENVIRONMENT_ID")
	environmentKey := mustEnv("ANTHROPIC_ENVIRONMENT_KEY")

	// Every per-session call (the runner's event stream/list/send, the lease
	// heartbeat, the force-stop) authenticates with the environment key.
	// option.WithAuthToken alone only ADDS an Authorization header — the
	// parent client's WithAPIKey middleware still sets X-Api-Key, so without
	// the explicit WithHeaderDel both creds would land on the wire and the
	// server would 401 the per-session calls. (lib/environments wraps the
	// same pair in helperReqOpts/bearerReqOpts for the in-tree helpers.)
	envKeyOpts := []option.RequestOption{
		option.WithHeaderDel("X-Api-Key"),
		option.WithAuthToken(environmentKey),
	}

	// The work poller claims items and yields them; it acks each one and posts
	// work.Stop after the loop body moves on (or Close runs).
	poller := environments.NewWorkPoller(ctx, client, environments.WorkPollerOptions{
		EnvironmentID:  envID,
		EnvironmentKey: environmentKey,
		Logger:         logger,
	})
	defer poller.Close()

	for poller.Next() {
		work := poller.Current()
		if work == nil {
			continue
		}
		sessionID := work.Data.ID
		log := logger.With(slog.String("work_id", work.ID), slog.String("session_id", sessionID))
		log.Info("claimed work")

		// Per-session agent tool context + skills. The session lookup and skill
		// download are environment-scoped, so they need the environment key
		// (envKeyOpts) just like the heartbeat and the runner below.
		env := &agenttoolset.AgentToolContext{Workdir: workdir()}
		if err := env.SetupSkills(ctx, client, sessionID, envKeyOpts...); err != nil {
			log.Warn("skill setup failed", slog.Any("error", err))
		}
		tools := append(agenttoolset.BetaAgentToolset20260401(env), currentTimeTool{})

		// Per-session context: the heartbeat goroutine cancels it when the
		// control plane reports the work is stopping/stopped or the lease was
		// not extended, which in turn stops the runner loop below.
		sessCtx, sessCancel := context.WithCancel(ctx)

		hbDone := make(chan struct{})
		go func() {
			defer close(hbDone)
			heartbeatLease(sessCtx, sessCancel, client, work, envKeyOpts, log)
		}()

		runner := client.Beta.Sessions.Events.NewToolRunner(sessCtx, sessionID, anthropic.SessionToolRunnerOptions{
			Tools:          tools,
			Logger:         log,
			RequestOptions: envKeyOpts,
		})
		for runner.Next() {
			printCall(runner.Current())
		}
		if err := runner.Err(); err != nil &&
			!errors.Is(err, anthropic.ErrSessionTerminated) &&
			!errors.Is(err, anthropic.ErrIdleTimeout) {
			log.Warn("tool runner exited with error", slog.Any("error", err))
		}
		_ = runner.Close()

		// Stop the heartbeat goroutine and wait for it before cleaning up.
		sessCancel()
		<-hbDone

		agenttoolset.CloseAll(tools)
		if err := env.Cleanup(); err != nil {
			log.Warn("skill cleanup failed", slog.Any("error", err))
		}

		// No explicit work.Stop here: the SessionToolRunner does not manage the
		// work item, but the WorkPoller does — it posts work.Stop for this item
		// when the loop moves on to the next poller.Next (or when the deferred
		// poller.Close runs for the final item).
	}
	if err := poller.Err(); err != nil {
		logger.Error("work poller failed", slog.Any("error", err))
	}
}

// heartbeatLease is a SIMPLIFIED lease heartbeat — see the comment block on
// observeAsSelfHostedWorker. It beats on a fixed interval; the first beat uses
// the "NO_HEARTBEAT" sentinel, each later one echoes the server's previous
// last_heartbeat. It calls cancel (which stops the runner loop) as soon as the
// control plane reports the work is stopping/stopped or the lease was not
// extended — or on any heartbeat error.
func heartbeatLease(ctx context.Context, cancel context.CancelFunc, client anthropic.Client, work *anthropic.BetaSelfHostedWork, reqOpts []option.RequestOption, logger *slog.Logger) {
	const interval = 30 * time.Second
	expectedLastHeartbeat := "NO_HEARTBEAT"

	// beat returns false when the worker should stop heartbeating (and the
	// caller should cancel the session).
	beat := func() bool {
		resp, err := client.Beta.Environments.Work.Heartbeat(ctx, work.ID, anthropic.BetaEnvironmentWorkHeartbeatParams{
			EnvironmentID:         work.EnvironmentID,
			ExpectedLastHeartbeat: param.NewOpt(expectedLastHeartbeat),
		}, reqOpts...)
		if err != nil {
			// A cancelled ctx is normal teardown, not a heartbeat failure.
			if ctx.Err() == nil {
				logger.Warn("heartbeat failed", slog.Any("error", err))
			}
			return false
		}
		expectedLastHeartbeat = resp.LastHeartbeat
		switch resp.State {
		case anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopping,
			anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopped:
			logger.Info("heartbeat reports shutdown", slog.String("state", string(resp.State)))
			return false
		}
		if !resp.LeaseExtended {
			logger.Warn("lease not extended; shutting down")
			return false
		}
		return true
	}

	for {
		if !beat() {
			cancel()
			return
		}
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
	}
}

// printCall logs one observed tool call: name, input, error flag, and whether
// the result posted back to the session.
func printCall(call anthropic.DispatchedToolCall) {
	status := "ok"
	if call.IsError {
		status = "error"
	}
	posted := ""
	if !call.Posted {
		posted = " [result post failed]"
	}
	fmt.Printf("tool %s(%s) -> %s%s\n", call.Name, truncate(callInput(call), 120), status, posted)
}

// callInput renders the raw tool input from the triggering event —
// CustomToolUse.Input for a custom tool call, ToolUse.Input otherwise.
func callInput(call anthropic.DispatchedToolCall) string {
	input := call.ToolUse.Input
	if call.Custom {
		input = call.CustomToolUse.Input
	}
	b, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf("<input: %v>", err)
	}
	return string(b)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Fprintf(os.Stderr, "missing required env var %s\n", k)
		os.Exit(1)
	}
	return v
}

func fatal(l *slog.Logger, op string, err error) {
	l.Error(op+" failed", slog.Any("error", err))
	os.Exit(1)
}

// Example demonstrating MCPBetaTools with the GitHub hosted MCP server.
//
// Prerequisites:
//   - GITHUB_TOKEN: a GitHub Personal Access Token with repo read access
//   - ANTHROPIC_API_KEY
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const githubMCPURL = "https://api.githubcopilot.com/mcp/"

// bearerTransport injects an Authorization header on every request.
type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(r)
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is required")
		os.Exit(1)
	}

	ctx := context.Background()

	// Connect to the GitHub hosted MCP server.
	mcpClient := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "anthropic-sdk-go-example", Version: "1.0.0"}, nil)
	session, err := mcpClient.Connect(ctx, &mcpsdk.StreamableClientTransport{
		Endpoint: githubMCPURL,
		HTTPClient: &http.Client{
			Transport: &bearerTransport{token: token, base: http.DefaultTransport},
		},
	}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer session.Close()

	// List tools and convert directly to BetaTools — no adapter code required.
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list tools: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Available tools (%d):\n", len(toolsResult.Tools))
	for _, t := range toolsResult.Tools {
		fmt.Printf("  - %s\n", t.Name)
	}
	fmt.Println()
	betaTools, err := mcp.NewBetaTools(toolsResult.Tools, session)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tools: %v\n", err)
		os.Exit(1)
	}

	client := anthropic.NewClient()
	question := "List the 5 most recently opened issues in the github/github-mcp-server repository. For each, include the issue number, title, and who opened it."
	fmt.Printf("[user]: %s\n\n", question)

	runner := client.Beta.Messages.NewToolRunner(betaTools, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_20250514,
			MaxTokens: 4096,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(question)),
			},
		},
		MaxIterations: 10,
	})

	finalMessage, err := runner.RunToCompletion(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, block := range finalMessage.Content {
		if tb, ok := block.AsAny().(anthropic.BetaTextBlock); ok {
			fmt.Println("[assistant]:", tb.Text)
		}
	}
}

package mcp_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func structuredResultTool(t *testing.T, result *mcpsdk.CallToolResult) anthropic.BetaTool {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	definition := &mcpsdk.Tool{Name: "result", InputSchema: map[string]any{"type": "object"}}
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "structured-result-test", Version: "1"}, nil)
	server.AddTool(definition, func(context.Context, *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) { return result, nil })
	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "structured-client", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	tool, err := mcp.NewBetaTool(definition, session)
	if err != nil {
		t.Fatal(err)
	}
	return tool
}

func TestMCPStructuredOnlyErrorsRetainPayload(t *testing.T) {
	for _, isError := range []bool{false, true} {
		for name, payload := range map[string]map[string]any{
			"empty":   {},
			"details": {"message": "try again", "retry": false, "details": map[string]any{"code": float64(0), "reason": nil}},
		} {
			prefix := "success/"
			if isError {
				prefix = "error/"
			}
			t.Run(prefix+name, func(t *testing.T) {
				tool := structuredResultTool(t, &mcpsdk.CallToolResult{IsError: isError, StructuredContent: payload})
				blocks, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
				var encoded string
				if isError {
					if err == nil {
						t.Fatal("expected tool error")
					}
					if blocks != nil {
						t.Fatalf("error also returned success blocks: %#v", blocks)
					}
					encoded = err.Error()
				} else {
					if err != nil {
						t.Fatal(err)
					}
					if len(blocks) != 1 || blocks[0].OfText == nil {
						t.Fatalf("unexpected success blocks: %#v", blocks)
					}
					encoded = blocks[0].OfText.Text
				}
				var actual map[string]any
				if err := json.Unmarshal([]byte(encoded), &actual); err != nil {
					t.Fatalf("lost structured payload: %q (%v)", encoded, err)
				}
				if !reflect.DeepEqual(actual, payload) {
					t.Fatalf("got %#v, want %#v", actual, payload)
				}
			})
		}
	}
}

func TestMCPStructuredErrorKeepsExistingContentPrecedence(t *testing.T) {
	for name, tc := range map[string]struct {
		content []mcpsdk.Content
		want    string
	}{
		"text":       {[]mcpsdk.Content{&mcpsdk.TextContent{Text: "first"}, &mcpsdk.TextContent{Text: "second"}}, "first\nsecond"},
		"empty text": {[]mcpsdk.Content{&mcpsdk.TextContent{Text: ""}}, "tool returned an error"},
	} {
		t.Run(name, func(t *testing.T) {
			tool := structuredResultTool(t, &mcpsdk.CallToolResult{IsError: true, Content: tc.content, StructuredContent: map[string]any{"ignored": true}})
			_, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
			if err == nil || err.Error() != tc.want {
				t.Fatalf("got %v, want %q", err, tc.want)
			}
		})
	}
}

func TestMCPEmptyResultStillDistinguishesSuccessAndError(t *testing.T) {
	for _, isError := range []bool{false, true} {
		tool := structuredResultTool(t, &mcpsdk.CallToolResult{IsError: isError})
		blocks, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
		if isError {
			if err == nil || err.Error() != "tool returned an error" {
				t.Fatalf("unexpected error: %v", err)
			}
		} else if err != nil {
			t.Fatal(err)
		}
		if blocks != nil {
			t.Fatalf("empty result became non-nil: %#v", blocks)
		}
	}
}

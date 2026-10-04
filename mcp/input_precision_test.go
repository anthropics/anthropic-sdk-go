package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func precisionTool(t *testing.T) (anthropic.BetaTool, <-chan json.RawMessage, *atomic.Int64) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	captured := make(chan json.RawMessage, 1)
	calls := new(atomic.Int64)
	definition := &mcpsdk.Tool{Name: "capture", InputSchema: map[string]any{"type": "object"}}
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "precision-test", Version: "1"}, nil)
	server.AddTool(definition, func(_ context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		calls.Add(1)
		captured <- append(json.RawMessage(nil), req.Params.Arguments...)
		return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "ok"}}}, nil
	})
	clientTransport, serverTransport := mcpsdk.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "precision-client", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	tool, err := mcp.NewBetaTool(definition, session)
	if err != nil {
		t.Fatal(err)
	}
	return tool, captured, calls
}

func TestMCPToolArgumentsPreserveNumbers(t *testing.T) {
	cases := map[string]string{
		"large integer":    `{"value":9007199254740993}`,
		"signed minimum":   `{"value":-9223372036854775808}`,
		"signed maximum":   `{"value":9223372036854775807}`,
		"unsigned maximum": `{"value":18446744073709551615}`,
		"precise decimal":  `{"value":0.12345678901234567890123456789}`,
		"large exponent":   `{"value":1e400}`,
		"small exponent":   `{"value":1e-400}`,
		"nested":           `{"value":{"items":[9007199254740993,-9007199254740993,0.1234567890123456789]}}`,
		"controls":         `{"value":[1,0,-0,1.5,"9007199254740993",true,null]}`,
		"empty object":     `{}`,
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			tool, captured, calls := precisionTool(t)
			original := json.RawMessage(input)
			before := append([]byte(nil), original...)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			output, err := tool.Execute(ctx, original)
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			var actual json.RawMessage
			select {
			case actual = <-captured:
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			// Compare numeric tokens without decoding them through a floating-point value.
			var expectedObject, actualObject map[string]json.RawMessage
			if err := json.Unmarshal(original, &expectedObject); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(actual, &actualObject); err != nil {
				t.Fatal(err)
			}
			if len(expectedObject) != len(actualObject) {
				t.Fatalf("wrong property count: %s", actual)
			}
			for key, want := range expectedObject {
				var expected, got bytes.Buffer
				if err := json.Compact(&expected, want); err != nil {
					t.Fatal(err)
				}
				if err := json.Compact(&got, actualObject[key]); err != nil {
					t.Fatal(err)
				}
				if expected.String() != got.String() {
					t.Errorf("%s: received %s; want %s", key, got.String(), expected.String())
				}
			}
			if !bytes.Equal(original, before) {
				t.Fatal("caller input changed")
			}
			if calls.Load() != 1 {
				t.Errorf("calls=%d, want 1", calls.Load())
			}
			if len(output) != 1 || output[0].OfText == nil || output[0].OfText.Text != "ok" {
				t.Fatalf("unexpected result: %#v", output)
			}
		})
	}
}

func TestMCPToolRejectsInvalidInputBeforeDispatch(t *testing.T) {
	for _, input := range []string{`{"value":`, `{"value":1} {}`, `[]`, `"text"`, `42`} {
		t.Run(input, func(t *testing.T) {
			tool, _, calls := precisionTool(t)
			_, err := tool.Execute(context.Background(), json.RawMessage(input))
			if err == nil || !strings.Contains(err.Error(), "failed to unmarshal input") {
				t.Fatalf("expected input error, got %v", err)
			}
			if calls.Load() != 0 {
				t.Fatalf("invalid input dispatched %d times", calls.Load())
			}
		})
	}
}

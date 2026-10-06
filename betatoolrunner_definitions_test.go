package anthropic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestBetaToolRunner_PreservesConfiguredDefinitions(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		for _, runnable := range []bool{false, true} {
			for _, configured := range []bool{false, true} {
				t.Run(fmt.Sprintf("runnable=%v/configured=%v", runnable, configured), func(t *testing.T) {
					var definitions []BetaToolUnionParam
					if configured {
						definitions = []BetaToolUnionParam{
							{OfWebSearchTool20250305: &BetaWebSearchTool20250305Param{MaxUses: Int(3)}},
							{OfTool: &BetaToolParam{Name: "external", InputSchema: BetaToolInputSchemaParam{Properties: map[string]any{}}}},
						}
					}
					// Spare capacity detects accidental appends into the caller's backing array.
					backing := make([]BetaToolUnionParam, len(definitions)+2)
					copy(backing, definitions)
					backing[len(definitions)] = BetaToolUnionParam{OfTool: &BetaToolParam{Name: "sentinel"}}
					params := pauseTurnParams(5)
					if configured {
						params.Tools = backing[:len(definitions)]
					}
					original, _ := json.Marshal(params)
					storage, _ := json.Marshal(backing)
					weather := &stubBetaTool{name: "weather"}
					var tools []BetaTool
					script := []string{pausedTurnJSON, endTurnJSON}
					if runnable {
						tools = []BetaTool{weather}
						script = []string{toolCallTurnJSON("toolu_weather", "weather"), endTurnJSON}
					}
					server, bodies := scriptedMessagesServer(t, script...)
					runner := newTurnRunner(newTestToolRunnerClient(server), stream, tools, params)
					expected := append([]BetaToolUnionParam{}, definitions...)
					if runnable {
						expected = append(expected, betaToolDefinition(weather))
					}
					want, _ := json.Marshal(expected)
					got, _ := json.Marshal(runner.Params.Tools)
					requireJSONEqual(t, string(got), string(want))
					runner.run(t, nil)
					sent := bodies()
					requireRequestCount(t, sent, 2)
					for _, body := range sent {
						value := gjson.GetBytes(body, "tools")
						if len(expected) > 0 {
							requireJSONEqual(t, value.Raw, string(want))
						} else if len(value.Array()) != 0 {
							t.Fatalf("unexpected tools: %s", value.Raw)
						}
					}
					if runnable && weather.runs.Load() != 1 {
						t.Fatalf("local tool runs=%d", weather.runs.Load())
					}
					if !runnable && weather.runs.Load() != 0 {
						t.Fatal("plain definitions must not become local tools")
					}
					after, _ := json.Marshal(params)
					afterStorage, _ := json.Marshal(backing)
					if string(after) != string(original) || string(afterStorage) != string(storage) {
						t.Fatal("runner changed caller-owned parameters or tool slice storage")
					}
					if configured {
						runner.Params.Tools[0] = BetaToolUnionParam{}
						if !reflect.DeepEqual(backing[0], definitions[0]) {
							t.Fatal("runner tool slice aliases caller storage")
						}
					}
				})
			}
		}
	})
}

func TestBetaToolRunner_ConfiguredDefinitionsSurviveCompaction(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := pauseTurnParams(5)
		params.Tools = []BetaToolUnionParam{{OfWebSearchTool20250305: &BetaWebSearchTool20250305Param{MaxUses: Int(2)}}}
		server, bodies := scriptedMessagesServer(t, compactionResponseJSON, endTurnJSON)
		runner := newTurnRunner(newTestToolRunnerClient(server), stream, nil, params)
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		runner.run(t, nil)
		expected, _ := json.Marshal(params.Tools)
		requireRequestCount(t, bodies(), 2)
		for _, body := range bodies() {
			requireJSONEqual(t, gjson.GetBytes(body, "tools").Raw, string(expected))
		}
		if len(runner.toolMap) != 0 {
			t.Fatal("server tools entered local dispatch map")
		}
	})
}

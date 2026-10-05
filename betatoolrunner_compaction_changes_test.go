package anthropic

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"
)

func compactedToolChanges(t *testing.T, changes string) BetaMessageParam {
	t.Helper()
	var message BetaMessageParam
	data := `{"role":"assistant","content":[{"type":"compaction","content":"summary","signature":"sig","tool_changes":` + changes + `}]}`
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		t.Fatal(err)
	}
	return message
}

const compactedWeatherRemoval = `[{"type":"tool_removal","tool":{"type":"tool_reference","name":"weather"}}]`
const compactedWeatherAddition = `[{"type":"tool_addition","tool":{"type":"tool_reference","name":"weather"}}]`

func TestCompactedToolChangesAffectDispatch(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		for _, tc := range []struct {
			name     string
			history  []BetaMessageParam
			wantRuns int32
		}{
			{"compacted removal", []BetaMessageParam{compactedToolChanges(t, compactedWeatherRemoval)}, 0},
			{"compacted readdition", []BetaMessageParam{systemToolChange(NewBetaToolRemovalBlock(weatherRef())), compactedToolChanges(t, compactedWeatherAddition)}, 1},
			{"later removal wins", []BetaMessageParam{compactedToolChanges(t, compactedWeatherAddition), systemToolChange(NewBetaToolRemovalBlock(weatherRef()))}, 0},
			{"later addition wins", []BetaMessageParam{compactedToolChanges(t, compactedWeatherRemoval), systemToolChange(NewBetaToolAdditionBlock(weatherRef()))}, 1},
		} {
			t.Run(tc.name, func(t *testing.T) {
				server, requests := scriptedMessagesServer(t, toolUseTurnJSON, endTurnJSON)
				weather := &stubBetaTool{name: "weather"}
				params := pauseTurnParams(10)
				params.Messages = append(tc.history, params.Messages...)
				original, err := json.Marshal(params.Messages)
				if err != nil {
					t.Fatal(err)
				}
				runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{weather}, params)
				runner.run(t, nil)
				if got := weather.runs.Load(); got != tc.wantRuns {
					t.Errorf("tool ran %d times, want %d", got, tc.wantRuns)
				}
				bodies := requests()
				requireRequestCount(t, bodies, 2)
				result := gjson.GetBytes(bodies[1], "messages").Array()
				last := result[len(result)-1].Get("content.0")
				if last.Get("is_error").Bool() != (tc.wantRuns == 0) {
					t.Errorf("wrong tool-result status: %s", last.Raw)
				}
				after, _ := json.Marshal(params.Messages)
				if string(original) != string(after) {
					t.Fatal("the supplied history was modified")
				}
			})
		}
	})
}

func TestCompactedToolChangesRetainOrderAndIgnoreUnrelatedBlocks(t *testing.T) {
	for _, tc := range []struct {
		name, changes string
		want          bool
	}{
		{"remove then add", `[{"type":"tool_removal","tool":{"type":"tool_reference","name":"weather"}},{"type":"tool_addition","tool":{"type":"tool_reference","name":"weather"}}]`, true},
		{"add then remove", `[{"type":"tool_addition","tool":{"type":"tool_reference","name":"weather"}},{"type":"tool_removal","tool":{"type":"tool_reference","name":"weather"}}]`, false},
		{"by value addition", `[{"type":"tool_removal","tool":{"type":"tool_reference","name":"weather"}},{"type":"tool_addition","tool":{"type":"tool_definition","definition":{"name":"weather","input_schema":{"type":"object"}}}}]`, true},
		{"empty changes", `[]`, true},
		{"future change", `[{"type":"future_change","tool":{"type":"tool_reference","name":"weather"}}]`, true},
		{"server tool removal", `[{"type":"tool_removal","tool":{"type":"mcp_tool_reference","server_name":"remote","name":"weather"}}]`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			message := compactedToolChanges(t, tc.changes)
			client := newTestToolRunnerClient(messagesServer(t))
			params := pauseTurnParams(10)
			params.Messages = []BetaMessageParam{message}
			runner := client.Beta.Messages.NewToolRunner([]BetaTool{&stubBetaTool{name: "weather"}}, params)
			_, got := runner.availableToolNames()["weather"]
			if got != tc.want {
				t.Errorf("available=%v, want %v", got, tc.want)
			}
		})
	}
	user := compactedToolChanges(t, compactedWeatherRemoval)
	user.Role = BetaMessageParamRoleUser
	client := newTestToolRunnerClient(messagesServer(t))
	params := pauseTurnParams(10)
	params.Messages = []BetaMessageParam{user}
	runner := client.Beta.Messages.NewToolRunner([]BetaTool{&stubBetaTool{name: "weather"}}, params)
	if _, ok := runner.availableToolNames()["weather"]; !ok {
		t.Fatal("a user block changed tool availability")
	}
}

func TestCompactedToolChangesSurviveAnotherCompaction(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		server, requests := scriptedMessagesServer(t, compactionResponseJSON, toolUseTurnJSON, endTurnJSON)
		weather := &stubBetaTool{name: "weather"}
		params := pauseTurnParams(10)
		params.Messages = append([]BetaMessageParam{compactedToolChanges(t, compactedWeatherRemoval)}, params.Messages...)
		runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{weather}, params)
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		runner.run(t, nil)
		requireRequestCount(t, requests(), 3)
		if got := weather.runs.Load(); got != 0 {
			t.Errorf("removed tool executed %d times after compaction", got)
		}
	})
}

func TestCompactedToolChangesInReturnedMessagesAffectTheNextDispatch(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		compactedTurn := `{"id":"msg_compacted","type":"message","role":"assistant","model":"m","content":[{"type":"compaction","content":"summary","signature":"sig","tool_changes":[{"type":"tool_removal","tool":{"type":"tool_reference","name":"weather"}}]},{"type":"tool_use","id":"toolu_1","name":"weather","input":{}}],"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":1}}`
		server, requests := scriptedMessagesServer(t, compactedTurn, endTurnJSON)
		weather := &stubBetaTool{name: "weather"}
		runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{weather}, pauseTurnParams(10))
		runner.run(t, nil)
		bodies := requests()
		requireRequestCount(t, bodies, 2)
		if got := weather.runs.Load(); got != 0 {
			t.Errorf("tool removed in model response ran %d times", got)
		}
		messages := gjson.GetBytes(bodies[1], "messages").Array()
		if !messages[len(messages)-1].Get("content.0.is_error").Bool() {
			t.Fatal("removed tool was not answered as unavailable")
		}
		requireJSONEqual(t, messages[len(messages)-2].Get("content.0.tool_changes").Raw, compactedWeatherRemoval)
	})
}

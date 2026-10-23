package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func threadRunnerParams() BetaToolRunnerParams {
	return BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1024,
		System:    []BetaTextBlockParam{{Text: "Be brief."}},
		Messages:  []BetaMessageParam{userMessage("Weather?")},
	}}
}

// weatherCall is a reply msg_<id> that calls the weather tool with call id.
func weatherCall(id string) threadResponse {
	return threadResponse{body: toolCallTurnJSON(id, "weather")}
}

func weatherCallJSON(id string) string {
	return `{"role":"assistant","content":` + gjson.Get(toolCallTurnJSON(id, "weather"), "content").Raw + `}`
}

func weatherResultJSON(id string) string {
	return toolResultMessageJSON(toolResultBlockJSON(id, "ok from weather"))
}

func summaryJSON() string {
	return `{"role":"assistant","content":` + gjson.Get(compactionResponseJSON, "content").Raw + `}`
}

// newThreadRunner returns a runner with the weather tool, which counts its
// runs, on a client of a server that answers with script.
func newThreadRunner(t *testing.T, stream bool, params BetaToolRunnerParams, script ...threadResponse) (*turnRunner, *threadServer, *stubBetaTool) {
	t.Helper()
	server := newThreadServer(t, script...)
	weather := &stubBetaTool{name: "weather"}
	return newTurnRunner(server.client(), stream, []BetaTool{weather}, params), server, weather
}

// mustNext takes the next turn, which must not fail.
func (r *turnRunner) mustNext(t *testing.T) *BetaMessage {
	t.Helper()
	message, err := r.next(context.Background())
	require.NoError(t, err)
	return message
}

// onReply returns a callback for run that calls do with the reply id.
func onReply(id string, do func(message *BetaMessage)) func(*BetaMessage) {
	return func(message *BetaMessage) {
		if message.ID == id {
			do(message)
		}
	}
}

func TestBetaToolRunner_Thread_FirstTurnCreatesAndLaterTurnsContinueWithTheToolResultsAlone(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), weatherCall("b"), reply("msg_c", "Sunny"))

		runner.run(t, nil)

		requests := server.requireRequests(t, 3)
		requireCreate(t, requests[0], messagesOf(userJSON("Weather?")))
		require.NotEmpty(t, requests[0].field("system"))
		require.NotEmpty(t, requests[0].field("tools"))
		requireContinue(t, requests[1], "msg_a", messagesOf(weatherResultJSON("a")))
		requireContinue(t, requests[2], "msg_b", messagesOf(weatherResultJSON("b")))
		for _, request := range requests {
			require.Equal(t, []string{threadsBeta}, request.betas())
			require.Equal(t, stream, gjson.GetBytes(request.body, "stream").Bool())
		}
	})
}

func TestBetaToolRunner_Thread_KeepsTheWholeHistoryInMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, _, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		want := messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), assistantJSON("Sunny"))
		requireMessages(t, runner.Messages(), want)
		requireMessages(t, runner.Params.Messages, want)
	})
}

// runnerKinds are the runners NewToolRunner returns: on a thread, and the two
// that send the whole history.
var runnerKinds = map[string]func(server *threadServer, params BetaToolRunnerParams) (Client, BetaToolRunnerParams){
	"on a thread": func(server *threadServer, params BetaToolRunnerParams) (Client, BetaToolRunnerParams) {
		return server.client(), params
	},
	"with DisableThread": func(server *threadServer, params BetaToolRunnerParams) (Client, BetaToolRunnerParams) {
		return server.client(), unthreaded(params)
	},
	"on a client without message threads": func(server *threadServer, params BetaToolRunnerParams) (Client, BetaToolRunnerParams) {
		return server.client(requestconfig.WithoutMessageThreads()), params
	},
}

// threadValues are the values a caller can give the Thread field.
var threadValues = map[string]BetaMessageThreadParamUnion{
	"create":   {OfCreate: &BetaThreadCreateParams{}},
	"continue": BetaMessageThreadParamOfContinue("msg_x"),
}

func TestBetaToolRunner_AThreadValueFromTheCallerIsRefused(t *testing.T) {
	for kind, of := range runnerKinds {
		for name, value := range threadValues {
			t.Run(kind+"/"+name, func(t *testing.T) {
				forEachRunner(t, func(t *testing.T, stream bool) {
					server := newThreadServer(t, reply("msg_a", "Sunny"))
					client, params := of(server, threadRunnerParams())
					params.Thread = value
					runner := newTurnRunner(client, stream, []BetaTool{&stubBetaTool{name: "weather"}}, params)

					_, err := runner.next(context.Background())

					requireLocalError(t, err)
					require.ErrorContains(t, err, "Params.Thread cannot be set on a tool runner")
					require.Empty(t, server.requests())
					require.Zero(t, runner.IterationCount())

					runner.Params.Thread = BetaMessageThreadParamUnion{}
					require.Equal(t, "msg_a", runner.mustNext(t).ID)
				})
			})
		}
	}
}

func TestBetaToolRunner_AThreadValueAssignedDuringTheRunIsRefusedBeforeAToolRuns(t *testing.T) {
	for kind, of := range runnerKinds {
		for name, value := range threadValues {
			t.Run(kind+"/"+name, func(t *testing.T) {
				forEachRunner(t, func(t *testing.T, stream bool) {
					server := newThreadServer(t, weatherCall("a"))
					client, params := of(server, threadRunnerParams())
					weather := &stubBetaTool{name: "weather"}
					runner := newTurnRunner(client, stream, []BetaTool{weather}, params)
					runner.mustNext(t)
					runner.Params.Thread = value

					_, err := runner.next(context.Background())

					requireLocalError(t, err)
					require.ErrorContains(t, err, "Params.Thread cannot be set on a tool runner")
					require.Len(t, server.requests(), 1)
					require.Zero(t, weather.runs.Load())
					require.Equal(t, 1, runner.IterationCount())
				})
			})
		}
	}
}

func TestBetaToolRunner_DisableThread_SendsTheWholeHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, unthreaded(threadRunnerParams()), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		requests := server.requireRequests(t, 2)
		requireUnthreaded(t, requests[0], messagesOf(userJSON("Weather?")))
		requireUnthreaded(t, requests[1], messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
		require.NotEmpty(t, requests[1].field("system"))
		require.NotEmpty(t, requests[1].field("tools"))
		require.Equal(t, []string{"BetaToolRunner"}, requests[1].header.Values("x-stainless-helper"))
	})
}

func TestBetaToolRunner_OnAClientWithoutMessageThreadsSendsTheWholeHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		server := newThreadServer(t, weatherCall("a"), reply("msg_b", "Sunny"))
		client := server.client(requestconfig.WithoutMessageThreads())
		runner := newTurnRunner(client, stream, []BetaTool{&stubBetaTool{name: "weather"}}, threadRunnerParams())

		runner.run(t, nil)

		requests := server.requireRequests(t, 2)
		requireUnthreaded(t, requests[0], messagesOf(userJSON("Weather?")))
		requireUnthreaded(t, requests[1], messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
		require.Equal(t, []string{"BetaToolRunner"}, requests[1].header.Values("x-stainless-helper"))
	})
}

func TestBetaToolRunner_Thread_KeepsTheHelperHeaderAndTheCallersBetas(t *testing.T) {
	params := threadRunnerParams()
	params.Betas = []AnthropicBeta{"files-api-2025-04-14"}
	server := newThreadServer(t, weatherCall("a"), reply("msg_b", "Sunny"))
	client := server.client()
	runner := client.Beta.Messages.NewToolRunner([]BetaTool{&stubBetaTool{name: "weather"}}, params, option.WithHeader("anthropic-beta", "mine-2026-01-01"))

	_, err := runner.RunToCompletion(context.Background())

	require.NoError(t, err)
	for _, request := range server.requireRequests(t, 2) {
		require.Equal(t, []string{"files-api-2025-04-14", threadsBeta, "mine-2026-01-01"}, request.betas())
		require.Equal(t, []string{"BetaToolRunner, message-thread"}, request.header.Values("x-stainless-helper"))
	}
}

func TestBetaToolRunner_Thread_StopsAtMaxIterations(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := threadRunnerParams()
		params.MaxIterations = 2
		runner, server, weather := newThreadRunner(t, stream, params, weatherCall("a"), weatherCall("b"))

		ids := runner.run(t, nil)

		require.Equal(t, []string{"msg_a", "msg_b"}, ids)
		require.Len(t, server.requests(), 2)
		require.True(t, runner.IsCompleted())
		require.EqualValues(t, 1, weather.runs.Load())
	})
}

func TestBetaToolRunner_Thread_ResumesAPausedTurnWithNoMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), pausedReply("msg_a", "Searching"), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", `[]`)
	})
}

func TestBetaToolRunner_Thread_SendsTheContainerOfTheReplyOnLaterTurns(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		call := strings.Replace(toolCallTurnJSON("a", "weather"), `"stop_sequence":null`, `"stop_sequence":null,"container":{"id":"container_a","expires_at":"2026-08-19T00:00:00Z"}`, 1)
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), threadResponse{body: call}, reply("msg_b", "Sunny"))

		runner.run(t, nil)

		requests := server.requireRequests(t, 2)
		require.Empty(t, requests[0].field("container"))
		requireContinue(t, requests[1], "msg_a", messagesOf(weatherResultJSON("a")))
		requireJSONEqual(t, requests[1].field("container"), `"container_a"`)
	})
}

func TestBetaToolRunner_DisableThread_SendsTheContainerOfTheReplyUnlessTheCallerNamedOne(t *testing.T) {
	for name, test := range map[string]struct {
		given BetaMessageNewParamsContainerUnion
		want  string
	}{
		"none named": {want: `"container_a"`},
		"one named":  {given: BetaMessageNewParamsContainerUnion{OfString: String("container_mine")}, want: `"container_mine"`},
	} {
		t.Run(name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				call := strings.Replace(toolCallTurnJSON("a", "weather"), `"stop_sequence":null`, `"stop_sequence":null,"container":{"id":"container_a","expires_at":"2026-08-19T00:00:00Z"}`, 1)
				params := unthreaded(threadRunnerParams())
				params.Container = test.given
				runner, server, _ := newThreadRunner(t, stream, params, threadResponse{body: call}, reply("msg_b", "Sunny"))

				runner.run(t, nil)

				requireJSONEqual(t, server.requireRequests(t, 2)[1].field("container"), test.want)
			})
		})
	}
}

func TestBetaToolRunner_Thread_ARejectedContinueIsReplayedUnderCreate(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, weather := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), threadRejection(http.StatusNotFound, "thread_not_found"), reply("msg_b", "Sunny"))

		ids := runner.run(t, nil)

		require.Equal(t, []string{"msg_a", "msg_b"}, ids)
		requests := server.requireRequests(t, 3)
		requireContinue(t, requests[1], "msg_a", messagesOf(weatherResultJSON("a")))
		requireCreate(t, requests[2], messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
		require.NotEmpty(t, requests[2].field("tools"))
		require.EqualValues(t, 1, weather.runs.Load())
	})
}

func TestBetaToolRunner_Thread_AReplayedThreadThatIsNotFoundEitherGoesOnWithoutAThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, weather := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), notFound(), weatherCall("b"), notFound(), weatherCall("c"), reply("msg_d", "Sunny"))

		ids := runner.run(t, nil)

		require.Equal(t, []string{"msg_a", "msg_b", "msg_c", "msg_d"}, ids)
		require.EqualValues(t, 3, weather.runs.Load())
		requests := server.requireRequests(t, 6)
		history := []string{userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), weatherCallJSON("b"), weatherResultJSON("b")}
		requireCreate(t, requests[2], messagesOf(history[:3]...))
		requireContinue(t, requests[3], "msg_b", messagesOf(weatherResultJSON("b")))
		requireUnthreaded(t, requests[4], messagesOf(history...))
		requireUnthreaded(t, requests[5], messagesOf(append(history, weatherCallJSON("c"), weatherResultJSON("c"))...))
	})
}

func TestBetaToolRunner_Thread_AnUnsupportedRequestGoesOnWithoutAThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), unsupportedRequest(`["inference_hooks"]`), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		requests := server.requireRequests(t, 3)
		requireCreate(t, requests[0], messagesOf(userJSON("Weather?")))
		requireUnthreaded(t, requests[1], messagesOf(userJSON("Weather?")))
		requireUnthreaded(t, requests[2], messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
	})
}

// What the caller changes through the runner.

func TestBetaToolRunner_Thread_AppendedMessagesFollowTheToolResultsOnTheContinue(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			runner.AppendMessages(userMessage("In Celsius."))
		}))

		requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(weatherResultJSON("a"), userJSON("In Celsius.")))
	})
}

func TestBetaToolRunner_DisableThread_MessagesAppendedAfterAPausedReplyAreSentWithIt(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, unthreaded(threadRunnerParams()), pausedReply("msg_a", "Searching"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			runner.AppendMessages(userMessage("In Celsius."))
		}))

		requireUnthreaded(t, server.requireRequests(t, 2)[1], messagesOf(userJSON("Weather?"), assistantJSON("Searching"), userJSON("In Celsius.")))
	})
}

func TestBetaToolRunner_Thread_AnExtendedHistoryKeepsTheThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(message *BetaMessage) {
			runner.Params.Messages = []BetaMessageParam{userMessage("Weather?"), message.ToParam(), userMessage("In Celsius.")}
		}))

		requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(weatherResultJSON("a"), userJSON("In Celsius.")))
	})
}

func TestBetaToolRunner_Thread_ARewrittenHistoryStartsANewThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(message *BetaMessage) {
			runner.Params.Messages = []BetaMessageParam{userMessage("Weather in Paris?"), message.ToParam()}
		}))

		request := server.requireRequests(t, 2)[1]
		requireCreate(t, request, messagesOf(userJSON("Weather in Paris?"), weatherCallJSON("a"), weatherResultJSON("a")))
		require.NotEmpty(t, request.field("tools"))
	})
}

// sentMessageEdits change, in place, the question or the tool result of a history
// that begins with "Weather?", the call a and its result.
var sentMessageEdits = map[string]struct {
	edit             func(messages []BetaMessageParam)
	question, result string
}{
	"the content of a message": {
		edit: func(messages []BetaMessageParam) {
			messages[0].OfUser.Content = userMessage("Weather in Paris?").OfUser.Content
		},
		question: userJSON("Weather in Paris?"),
		result:   weatherResultJSON("a"),
	},
	"the text of a block": {
		edit: func(messages []BetaMessageParam) {
			messages[0].OfUser.Content[0].OfText.Text = "Weather in Paris?"
		},
		question: userJSON("Weather in Paris?"),
		result:   weatherResultJSON("a"),
	},
	"the cache breakpoint of a block": {
		edit: func(messages []BetaMessageParam) {
			messages[0].OfUser.Content[0].OfText.CacheControl = NewBetaCacheControlEphemeralParam()
		},
		question: markedUserJSON("Weather?"),
		result:   weatherResultJSON("a"),
	},
	"the content of a tool result": {
		edit: func(messages []BetaMessageParam) {
			messages[2].OfUser.Content[0].OfToolResult.Content = []BetaToolResultBlockParamContentUnion{{OfText: &BetaTextBlockParam{Text: "trimmed"}}}
		},
		question: userJSON("Weather?"),
		result:   toolResultMessageJSON(toolResultBlockJSON("a", "trimmed")),
	},
}

func TestBetaToolRunner_ASentMessageChangedInPlaceIsSentAsItIsNow(t *testing.T) {
	for kind, of := range runnerKinds {
		for name, change := range sentMessageEdits {
			t.Run(kind+"/"+name, func(t *testing.T) {
				forEachRunner(t, func(t *testing.T, stream bool) {
					server := newThreadServer(t, weatherCall("a"), weatherCall("b"), weatherCall("c"), reply("msg_d", "Sunny"))
					client, params := of(server, threadRunnerParams())
					weather := &stubBetaTool{name: "weather"}
					runner := newTurnRunner(client, stream, []BetaTool{weather}, params)

					ids := runner.run(t, onReply("msg_b", func(*BetaMessage) { change.edit(runner.Params.Messages) }))

					edited := []string{change.question, weatherCallJSON("a"), change.result, weatherCallJSON("b"), weatherResultJSON("b")}
					whole := append(slices.Clone(edited), weatherCallJSON("c"), weatherResultJSON("c"))
					requests := server.requireRequests(t, 4)
					if kind == "on a thread" {
						requireCreate(t, requests[2], messagesOf(edited...))
						requireContinue(t, requests[3], "msg_c", messagesOf(weatherResultJSON("c")))
					} else {
						requireUnthreaded(t, requests[2], messagesOf(edited...))
						requireUnthreaded(t, requests[3], messagesOf(whole...))
					}
					require.Equal(t, []string{"msg_a", "msg_b", "msg_c", "msg_d"}, ids)
					require.EqualValues(t, 3, weather.runs.Load())
					requireMessages(t, runner.Messages(), messagesOf(append(whole, assistantJSON("Sunny"))...))
				})
			})
		}
	}
}

func TestBetaToolRunnerStreaming_AMessageChangedInPlaceWhileItsReplyStreamsIsSentAsItIsNow(t *testing.T) {
	for kind, of := range runnerKinds {
		t.Run(kind, func(t *testing.T) {
			server := newThreadServer(t, weatherCall("a"), reply("msg_b", "Sunny"))
			client, params := of(server, threadRunnerParams())
			runner := client.Beta.Messages.NewToolRunnerStreaming([]BetaTool{&stubBetaTool{name: "weather"}}, params)

			for events, err := range runner.AllStreaming(context.Background()) {
				require.NoError(t, err)
				for _, err := range events {
					require.NoError(t, err)
					runner.Params.Messages[0].OfUser.Content[0].OfText.Text = "Weather in Paris?"
				}
			}

			edited := messagesOf(userJSON("Weather in Paris?"), weatherCallJSON("a"), weatherResultJSON("a"))
			request := server.requireRequests(t, 2)[1]
			if kind == "on a thread" {
				requireCreate(t, request, edited)
			} else {
				requireUnthreaded(t, request, edited)
			}
		})
	}
}

func TestBetaToolRunner_AMessageTheCallerChangesAfterCreatingTheRunnerStaysAsItWasGiven(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := threadRunnerParams()
		runner, server, _ := newThreadRunner(t, stream, params, reply("msg_a", "Sunny"))

		params.Messages[0].OfUser.Content = userMessage("Changed").OfUser.Content
		runner.run(t, nil)

		requireCreate(t, server.requireRequests(t, 1)[0], messagesOf(userJSON("Weather?")))
	})
}

func TestBetaToolRunner_AppendMessages_AMessageTheCallerChangesAfterwardsStaysAsItWasGiven(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			own := userMessage("In Celsius.")
			runner.AppendMessages(own)
			own.OfUser.Content = userMessage("Changed").OfUser.Content
		}))

		requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(weatherResultJSON("a"), userJSON("In Celsius.")))
	})
}

func TestBetaToolRunner_Thread_AHistoryRewrittenAfterAPausedReplyStartsANewThreadThatEndsWithTheReply(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), pausedReply("msg_a", "Searching"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(message *BetaMessage) {
			runner.Params.Messages = []BetaMessageParam{userMessage("Weather in Paris?"), message.ToParam()}
		}))

		requireCreate(t, server.requireRequests(t, 2)[1], messagesOf(userJSON("Weather in Paris?"), assistantJSON("Searching")))
	})
}

func TestBetaToolRunner_Thread_AssignedParamsReplaceTheFields(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := threadRunnerParams()
		params.StopSequences = []string{"END"}
		runner, server, _ := newThreadRunner(t, stream, params, weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			runner.Params.BetaMessageNewParams = BetaMessageNewParams{
				Model:     runner.Params.Model,
				MaxTokens: 12,
				System:    runner.Params.System,
				Tools:     runner.Params.Tools,
				Messages:  runner.Params.Messages,
			}
		}))

		requests := server.requireRequests(t, 2)
		requireJSONEqual(t, requests[0].field("stop_sequences"), `["END"]`)
		requireContinue(t, requests[1], "msg_a", messagesOf(weatherResultJSON("a")))
		require.Empty(t, requests[1].field("stop_sequences"))
		requireJSONEqual(t, requests[1].field("max_tokens"), `12`)
	})
}

// historyChanges are what the caller does to a history that holds "Weather?" and
// the call a.
var historyChanges = map[string]func(runner *turnRunner){
	"left as it is": func(*turnRunner) {},
	"changed in place": func(runner *turnRunner) {
		runner.Params.Messages[0].OfUser.Content[0].OfText.Text = "Weather in Paris?"
	},
	"replaced": func(runner *turnRunner) {
		runner.Params.Messages = []BetaMessageParam{userMessage("Weather in Paris?"), runner.Params.Messages[1]}
	},
}

func TestBetaToolRunner_Thread_AChangedSystemIsRefusedBeforeAToolRuns(t *testing.T) {
	for name, changeHistory := range historyChanges {
		t.Run("with the history "+name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				runner, server, weather := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"))
				runner.mustNext(t)
				runner.Params.System = []BetaTextBlockParam{{Text: "Be verbose."}}
				changeHistory(runner)

				_, err := runner.next(context.Background())

				requireLocalError(t, err)
				require.ErrorContains(t, err, "System cannot be changed on a message thread")
				require.Len(t, server.requests(), 1)
				require.Zero(t, weather.runs.Load())
			})
		})
	}
}

func TestBetaToolRunner_DisableThread_AChangedSystemIsSent(t *testing.T) {
	for name, changeHistory := range historyChanges {
		t.Run("with the history "+name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				runner, server, _ := newThreadRunner(t, stream, unthreaded(threadRunnerParams()), weatherCall("a"), reply("msg_b", "Sunny"))

				runner.run(t, onReply("msg_a", func(*BetaMessage) {
					runner.Params.System = []BetaTextBlockParam{{Text: "Be verbose."}}
					changeHistory(runner)
				}))

				requireJSONEqual(t, server.requireRequests(t, 2)[1].field("system.0.text"), `"Be verbose."`)
			})
		})
	}
}

// toolChangeScript calls weather, then lookup, then weather again, then ends.
func toolChangeScript() []threadResponse {
	return []threadResponse{weatherCall("a"), {body: toolCallTurnJSON("b", "lookup")}, weatherCall("c"), reply("msg_d", "Sunny")}
}

// changeTools returns a callback for run that adds lookup on the first reply of
// toolChangeScript and removes weather on the second.
func changeTools(runner *turnRunner, lookup, weather BetaTool) func(*BetaMessage) {
	return func(message *BetaMessage) {
		switch message.ID {
		case "msg_a":
			runner.AddTools(lookup)
		case "msg_b":
			runner.RemoveTools(weather)
		}
	}
}

func TestBetaToolRunner_Thread_AddedAndRemovedToolsTravelInTheContinues(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, weather := newThreadRunner(t, stream, threadRunnerParams(), toolChangeScript()...)
		lookup := &stubBetaTool{name: "lookup"}

		ids := runner.run(t, changeTools(runner, lookup, weather))

		require.Equal(t, []string{"msg_a", "msg_b", "msg_c", "msg_d"}, ids)
		require.NoError(t, runner.Err())
		requests := server.requireRequests(t, 4)
		requireCreate(t, requests[0], messagesOf(userJSON("Weather?")))
		requireJSONEqual(t, requests[0].field("tools.#.name"), `["weather"]`)
		requireContinue(t, requests[1], "msg_a", messagesOf(weatherResultJSON("a"), toolChangeMessageJSON(stubAdditionJSON("lookup"))))
		requireContinue(t, requests[2], "msg_b", messagesOf(toolResultMessageJSON(toolResultBlockJSON("b", "ok from lookup")), toolChangeMessageJSON(toolRemovalJSON("weather"))))
		requireContinue(t, requests[3], "msg_c", messagesOf(toolResultMessageJSON(toolNotFoundBlockJSON("c", "weather"))))
		require.EqualValues(t, 1, lookup.runs.Load())
		require.EqualValues(t, 1, weather.runs.Load())
	})
}

// Compaction.

func TestBetaToolRunner_Thread_ACompactionCarriesTheWholeHistoryUnderANewThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "Sunny"))

		runner.run(t, compactOnToolUse(runner))

		requests := server.requireRequests(t, 3)
		requireCreate(t, requests[1], messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
		requireJSONEqual(t, requests[1].field("compaction"), `{"type":"summarize"}`)
		requireCreate(t, requests[2], messagesOf(summaryJSON()))
	})
}

func TestBetaToolRunner_Thread_ACompactionRejectedAsUnsupportedCostsOneRequestAndStopsTheThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, server, _ := newThreadRunner(t, stream, threadRunnerParams(), weatherCall("a"), unsupportedRequest(`["context_management"]`), threadResponse{body: compactionResponseJSON}, weatherCall("b"), reply("msg_c", "Sunny"))

		ids := runner.run(t, onReply("msg_a", func(*BetaMessage) {
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		}))

		require.Equal(t, []string{"msg_a", "msg_compaction", "msg_b", "msg_c"}, ids)
		requests := server.requireRequests(t, 5)
		history := messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"))
		requireCreate(t, requests[1], history)
		requireUnthreaded(t, requests[2], history)
		for _, request := range requests[1:3] {
			requireJSONEqual(t, request.field("compaction"), `{"type":"summarize"}`)
		}
		requireUnthreaded(t, requests[3], messagesOf(summaryJSON()))
		requireUnthreaded(t, requests[4], messagesOf(summaryJSON(), weatherCallJSON("b"), weatherResultJSON("b")))
	})
}

// Streams.

func TestBetaToolRunnerStreaming_Thread_AStreamThatEndsBeforeTheReplyStartsEndsTheRun(t *testing.T) {
	server := newThreadServer(t, weatherCall("a"), threadResponse{stream: true})
	client := server.client()
	runner := client.Beta.Messages.NewToolRunnerStreaming([]BetaTool{&stubBetaTool{name: "weather"}}, threadRunnerParams())

	for events, err := range runner.AllStreaming(context.Background()) {
		require.NoError(t, err)
		for _, err := range events {
			require.NoError(t, err)
		}
	}

	require.NoError(t, runner.Err())
	require.Len(t, server.requests(), 2)
}

// Runner over a thread.

func weatherThreadParams() BetaMessageThreadCreateParams {
	params := threadParams()
	params.Tools = []BetaToolUnionParam{betaToolDefinition(&stubBetaTool{name: "weather"})}
	return params
}

// newRunnerOverThread returns a thread that was created with the weather tool,
// has sent one turn and has a question staged, and a runner with the tool over
// it.
func newRunnerOverThread(t *testing.T, stream bool, params BetaMessageThreadCreateParams, script ...threadResponse) (*BetaMessageThread, *turnRunner, *threadServer) {
	t.Helper()
	params.Tools = weatherThreadParams().Tools
	thread, server := newThread(t, params, append([]threadResponse{reply("msg_0", "Hi")}, script...)...)
	send(t, thread, turn("Hello"))
	thread.Append(userMessage("Weather?"))
	return thread, newTurnRunnerOver(thread, stream, &stubBetaTool{name: "weather"}), server
}

func newTurnRunnerOver(thread *BetaMessageThread, stream bool, tools ...BetaTool) *turnRunner {
	if !stream {
		runner := thread.NewToolRunner(tools)
		return &turnRunner{&runner.betaToolRunnerBase, runner.NextMessage}
	}
	runner := thread.NewToolRunnerStreaming(tools)
	return &turnRunner{&runner.betaToolRunnerBase, func(ctx context.Context) (*BetaMessage, error) {
		var message BetaMessage
		events := 0
		for event, err := range runner.NextStreaming(ctx) {
			if err != nil {
				return nil, err
			}
			if err := message.Accumulate(event); err != nil {
				return nil, err
			}
			events++
		}
		if events == 0 {
			return nil, nil
		}
		return &message, nil
	}}
}

func TestBetaMessageThread_NewToolRunner_SendsNothing(t *testing.T) {
	thread, _, server := newRunnerOverThread(t, false, threadParams())

	require.Len(t, server.requests(), 1)
	require.Equal(t, "msg_0", thread.PreviousMessageID())
}

func TestBetaMessageThread_NewToolRunner_ContinuesTheThreadWhichHoldsEveryTurnAfterwards(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), weatherCall("a"), reply("msg_b", "Sunny"), reply("msg_c", "Welcome"))

		runner.run(t, nil)
		send(t, thread, turn("Thanks"))

		requests := server.requireRequests(t, 4)
		requireContinue(t, requests[1], "msg_0", messagesOf(userJSON("Weather?")))
		requireContinue(t, requests[2], "msg_a", messagesOf(weatherResultJSON("a")))
		requireContinue(t, requests[3], "msg_b", messagesOf(userJSON("Thanks")))
		requireMessages(t, thread.Messages(), messagesOf(
			userJSON("Hello"), assistantJSON("Hi"),
			userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), assistantJSON("Sunny"),
			userJSON("Thanks"), assistantJSON("Welcome"),
		))
	})
}

func TestBetaMessageThread_NewToolRunner_ARefusedTurnEndsTheRunAndLeavesItsToolResultsStaged(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), weatherCall("a"), refusedReply("msg_b"), reply("msg_c", "Welcome"))

		ids := runner.run(t, nil)

		require.Equal(t, []string{"msg_a", "msg_b"}, ids)
		require.True(t, runner.IsCompleted())
		require.Equal(t, BetaStopReasonRefusal, runner.LastMessage().StopReason)
		history := []string{userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")}
		requireMessages(t, runner.Messages(), messagesOf(history...))
		requireMessages(t, thread.Messages(), messagesOf(history...))

		send(t, thread, turn("Thanks"))

		requests := server.requireRequests(t, 4)
		requireContinue(t, requests[2], "msg_a", messagesOf(weatherResultJSON("a")))
		requireCreate(t, requests[3], messagesOf(append(history, userJSON("Thanks"))...))
		require.NotEmpty(t, requests[3].field("tools"))
	})
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedDuringARefusedTurnStaysStagedWithoutTheReply(t *testing.T) {
	for name, appendTo := range map[string]func(thread *BetaMessageThread, runner *turnRunner){
		"to the runner": func(_ *BetaMessageThread, runner *turnRunner) { runner.AppendMessages(userMessage("Briefly.")) },
		"to the thread": func(thread *BetaMessageThread, _ *turnRunner) { thread.Append(userMessage("Briefly.")) },
	} {
		t.Run(name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, runner, server := newRunnerOverThread(t, stream, threadParams(), weatherCall("a"), refusedReply("msg_b"))
				thread.Options = append(thread.Options, duringRequest(server, 2, func() { appendTo(thread, runner) }))

				runner.run(t, nil)

				want := messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), userJSON("Briefly."))
				requireMessages(t, runner.Messages(), want)
				requireMessages(t, thread.Messages(), want)
			})
		})
	}
}

func TestBetaToolRunner_Thread_GoesOnUnderANewThreadAfterAReplyWrittenByAFallbackModel(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := threadRunnerParams()
		params.Fallbacks, params.Betas = fallbackParams().Fallbacks, fallbackParams().Betas
		runner, server, weather := newThreadRunner(t, stream, params, writtenByFallback(weatherCall("a")), weatherCall("b"), reply("msg_c", "Sunny"))

		ids := runner.run(t, nil)

		require.Equal(t, []string{"msg_a", "msg_b", "msg_c"}, ids)
		require.EqualValues(t, 2, weather.runs.Load())
		require.Equal(t, params.Model, runner.Params.Model)
		require.Equal(t, params.Fallbacks, runner.Params.Fallbacks)
		requests := server.requireRequests(t, 3)
		call := `{"role":"assistant","content":[` + fallbackBlockJSON + `,` + gjson.Get(weatherCallJSON("a"), "content.0").Raw + `]}`
		requireCreate(t, requests[1], messagesOf(userJSON("Weather?"), call, weatherResultJSON("a")))
		require.NotEmpty(t, requests[1].field("tools"))
		requireContinue(t, requests[2], "msg_b", messagesOf(weatherResultJSON("b")))
		for _, request := range requests {
			requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
			requireJSONEqual(t, request.field("fallbacks"), `[{"model":"`+fallbackModel+`"}]`)
		}
	})
}

func TestBetaToolRunner_DisableThread_AReplyWrittenByAFallbackModelLeavesTheModelAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := unthreaded(threadRunnerParams())
		params.Fallbacks, params.Betas = fallbackParams().Fallbacks, fallbackParams().Betas
		runner, server, _ := newThreadRunner(t, stream, params, writtenByFallback(weatherCall("a")), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		request := server.requireRequests(t, 2)[1]
		requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
		require.NotEmpty(t, request.field("fallbacks"))
	})
}

func TestBetaToolRunner_DisableThread_KeepsARefusedTurnInTheHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		runner, _, _ := newThreadRunner(t, stream, unthreaded(threadRunnerParams()), weatherCall("a"), refusedReply("msg_b"))

		runner.run(t, nil)

		requireMessages(t, runner.Messages(), messagesOf(userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), assistantJSON("I can't help with that.")))
	})
}

func TestBetaMessageThread_NewToolRunner_AddedAndRemovedToolsTravelInTheContinuesAndStayInTheHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), append(toolChangeScript(), reply("msg_e", "Welcome"))...)
		lookup := &stubBetaTool{name: "lookup"}

		runner.run(t, changeTools(runner, lookup, &stubBetaTool{name: "weather"}))
		send(t, thread, turn("Thanks"))

		addition, removal := toolChangeMessageJSON(stubAdditionJSON("lookup")), toolChangeMessageJSON(toolRemovalJSON("weather"))
		lookupResult := toolResultMessageJSON(toolResultBlockJSON("b", "ok from lookup"))
		weatherNotFound := toolResultMessageJSON(toolNotFoundBlockJSON("c", "weather"))
		requests := server.requireRequests(t, 6)
		requireContinue(t, requests[1], "msg_0", messagesOf(userJSON("Weather?")))
		requireContinue(t, requests[2], "msg_a", messagesOf(weatherResultJSON("a"), addition))
		requireContinue(t, requests[3], "msg_b", messagesOf(lookupResult, removal))
		requireContinue(t, requests[4], "msg_c", messagesOf(weatherNotFound))
		requireContinue(t, requests[5], "msg_d", messagesOf(userJSON("Thanks")))
		require.EqualValues(t, 1, lookup.runs.Load())
		requireMessages(t, thread.Messages(), messagesOf(
			userJSON("Hello"), assistantJSON("Hi"),
			userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), addition,
			`{"role":"assistant","content":`+gjson.Get(toolCallTurnJSON("b", "lookup"), "content").Raw+`}`, lookupResult, removal,
			weatherCallJSON("c"), weatherNotFound, assistantJSON("Sunny"),
			userJSON("Thanks"), assistantJSON("Welcome"),
		))
	})
}

func TestBetaMessageThread_NewToolRunner_StoppedAtMaxIterationsLeavesTheThreadWithTheReplyAndItsContainer(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		call := strings.Replace(toolCallTurnJSON("a", "weather"), `"stop_sequence":null`, `"stop_sequence":null,"container":{"id":"container_a","expires_at":"2026-08-19T00:00:00Z"}`, 1)
		thread, server := newThread(t, weatherThreadParams(), threadResponse{body: call}, reply("msg_b", "Sunny"))
		thread.Append(userMessage("Weather?"))
		weather := &stubBetaTool{name: "weather"}
		runner := newTurnRunnerOver(thread, stream, weather)
		runner.Params.MaxIterations = 1

		runner.run(t, nil)
		server.requireRequests(t, 1)
		require.Zero(t, weather.runs.Load())
		send(t, thread, turn("Thanks"))

		request := server.requireRequests(t, 2)[1]
		requireContinue(t, request, "msg_a", messagesOf(userJSON("Thanks")))
		requireJSONEqual(t, request.field("container"), `"container_a"`)
	})
}

func TestBetaMessageThread_NewToolRunner_ContinuesFromTheReplyToATurnThatWasSentAgainWithoutItsPartialReply(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, weatherThreadParams(), textCutOff("msg_0", "par"), weatherCall("a"), reply("msg_b", "Sunny"))
		readStream(thread, turn("Weather?"))

		newTurnRunnerOver(thread, stream, &stubBetaTool{name: "weather"}).run(t, nil)

		requests := server.requireRequests(t, 3)
		requireCreate(t, requests[1], messagesOf(userJSON("Weather?")))
		requireContinue(t, requests[2], "msg_a", messagesOf(weatherResultJSON("a")))
	})
}

func TestBetaMessageThread_NewToolRunner_AddsTheDefinitionsThatTheThreadLacks(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, weatherThreadParams(), reply("msg_a", "Sunny"))
		thread.Append(userMessage("Weather?"))

		newTurnRunnerOver(thread, stream, &stubBetaTool{name: "weather"}, &stubBetaTool{name: "time"}).run(t, nil)

		request := server.requireRequests(t, 1)[0]
		requireJSONEqual(t, request.field("tools.#"), `2`)
		requireJSONEqual(t, request.field("tools.0.name"), `"weather"`)
		requireJSONEqual(t, request.field("tools.1.name"), `"time"`)
	})
}

func TestBetaMessageThread_NewToolRunner_AToolThatTheThreadWasNotCreatedWithIsRefused(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, weatherThreadParams(), reply("msg_0", "Hi"))
		send(t, thread, turn("Hello"))
		thread.Append(userMessage("Time?"))

		_, err := newTurnRunnerOver(thread, stream, &stubBetaTool{name: "time"}).next(context.Background())

		requireLocalError(t, err)
		require.ErrorContains(t, err, "Tools cannot be changed on a message thread")
		require.Len(t, server.requests(), 1)
	})
}

func TestBetaMessageThread_NewToolRunner_OverAForkRunsTheToolsAndLeavesTheOriginalAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, weatherThreadParams(), reply("msg_0", "Hi"), weatherCall("a"), reply("msg_b", "Sunny"))
		send(t, thread, turn("Hello"))
		fork := thread.Fork()
		fork.Append(userMessage("Weather?"))
		weather := &stubBetaTool{name: "weather"}

		newTurnRunnerOver(fork, stream, weather).run(t, nil)

		request := server.requireRequests(t, 3)[1]
		requireCreate(t, request, messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?")))
		requireJSONEqual(t, request.field("tools.#"), `1`)
		require.EqualValues(t, 1, weather.runs.Load())
		require.Equal(t, "msg_0", thread.PreviousMessageID())
		requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Hi")))
	})
}

func TestBetaMessageThread_NewToolRunner_AnUnsupportedRequestIsReturnedOnADefaultThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		_, runner, server := newRunnerOverThread(t, stream, threadParams(), unsupportedRequest(`["inference_hooks"]`))

		_, err := runner.next(context.Background())

		requireStatus(t, err, http.StatusBadRequest)
		require.Len(t, server.requests(), 2)
	})
}

func TestBetaMessageThread_NewToolRunner_AnUnsupportedRequestGoesOnWithoutAThreadOnAContinueUnthreadedThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		_, runner, server := newRunnerOverThread(t, stream, continueUnthreadedParams(), unsupportedRequest(`["inference_hooks"]`), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, nil)

		requests := server.requireRequests(t, 4)
		requireUnthreaded(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?")))
		requireUnthreaded(t, requests[3], messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")))
	})
}

func TestBetaMessageThread_NewToolRunner_ACompactionResponseChangedInPlaceStartsANewThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), reply("msg_a", "Sunny"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "Welcome"), reply("msg_c", "Bye"))
		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		}))
		send(t, thread, turn("Thanks"))

		thread.Params.Messages[0].OfAssistant.Content = assistantMessage("Summary").OfAssistant.Content
		send(t, thread, turn("Bye"))

		requireCreate(t, server.requireRequests(t, 5)[4], messagesOf(assistantJSON("Summary"), userJSON("Thanks"), assistantJSON("Welcome"), userJSON("Bye")))
	})
}

// The run ends with the compaction response: nothing has to step the runner again.
func TestBetaMessageThread_NewToolRunner_ACompactionAfterTheFinalReplyLeavesTheThreadWithTheSummaryHoweverTheRunIsDriven(t *testing.T) {
	ctx := context.Background()
	tools := []BetaTool{&stubBetaTool{name: "weather"}}
	compactAfterTheFinalReply := func(runner *betaToolRunnerBase) {
		if runner.LastMessage().ID == "msg_a" {
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		}
	}

	for name, run := range map[string]func(t *testing.T, thread *BetaMessageThread){
		"RunToCompletion": func(t *testing.T, thread *BetaMessageThread) {
			runner := thread.NewToolRunner(tools)
			_, err := runner.NextMessage(ctx)
			require.NoError(t, err)
			compactAfterTheFinalReply(&runner.betaToolRunnerBase)
			_, err = runner.RunToCompletion(ctx)
			require.NoError(t, err)
		},
		"All": func(t *testing.T, thread *BetaMessageThread) {
			runner := thread.NewToolRunner(tools)
			for _, err := range runner.All(ctx) {
				require.NoError(t, err)
				compactAfterTheFinalReply(&runner.betaToolRunnerBase)
			}
		},
		"NextMessage until IsCompleted": func(t *testing.T, thread *BetaMessageThread) {
			runner := thread.NewToolRunner(tools)
			for !runner.IsCompleted() {
				_, err := runner.NextMessage(ctx)
				require.NoError(t, err)
				compactAfterTheFinalReply(&runner.betaToolRunnerBase)
			}
		},
		"AllStreaming": func(t *testing.T, thread *BetaMessageThread) {
			runner := thread.NewToolRunnerStreaming(tools)
			for events, err := range runner.AllStreaming(ctx) {
				require.NoError(t, err)
				for _, err := range events {
					require.NoError(t, err)
				}
				compactAfterTheFinalReply(&runner.betaToolRunnerBase)
			}
		},
		"NextStreaming until IsCompleted": func(t *testing.T, thread *BetaMessageThread) {
			runner := thread.NewToolRunnerStreaming(tools)
			for !runner.IsCompleted() {
				for _, err := range runner.NextStreaming(ctx) {
					require.NoError(t, err)
				}
				compactAfterTheFinalReply(&runner.betaToolRunnerBase)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			thread, _, server := newRunnerOverThread(t, false, threadParams(), reply("msg_a", "Sunny"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "Welcome"))

			run(t, thread)
			send(t, thread, turn("Thanks"))

			request := server.requireRequests(t, 4)[3]
			requireCreate(t, request, messagesOf(summaryJSON(), userJSON("Thanks")))
			require.Empty(t, request.field("compaction"))
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AFailedCompactionAfterTheFinalReplyLeavesTheThreadAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), reply("msg_a", "Sunny"), plainError(http.StatusBadRequest, "rejected"), reply("msg_b", "Welcome"))

		_, err := runner.next(context.Background())
		require.NoError(t, err)
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		_, err = runner.next(context.Background())
		requireStatus(t, err, http.StatusBadRequest)
		require.True(t, runner.IsCompleted())
		send(t, thread, turn("Thanks"))

		requests := server.requireRequests(t, 4)
		require.NotEmpty(t, requests[2].field("compaction"))
		requireCreate(t, requests[3], messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), assistantJSON("Sunny"), userJSON("Thanks")))
		require.Empty(t, requests[3].field("compaction"))
	})
}

func TestBetaMessageThread_NewToolRunnerStreaming_ACompactionThatTheCallerStopsReadingLeavesNoCompactionFieldsOnTheThread(t *testing.T) {
	thread, _, server := newRunnerOverThread(t, true, threadParams(), reply("msg_a", "Sunny"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "Welcome"))
	runner := thread.NewToolRunnerStreaming([]BetaTool{&stubBetaTool{name: "weather"}})

	for _, err := range runner.NextStreaming(context.Background()) {
		require.NoError(t, err)
	}
	runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
	for _, err := range runner.NextStreaming(context.Background()) {
		require.NoError(t, err)
		break
	}
	require.True(t, runner.IsCompleted())
	send(t, thread, turn("Thanks"))

	requests := server.requireRequests(t, 4)
	require.NotEmpty(t, requests[2].field("compaction"))
	requireCreate(t, requests[3], messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), assistantJSON("Sunny"), userJSON("Thanks")))
	require.Empty(t, requests[3].field("compaction"))
}

// What the caller appends to the thread during the run.

// duringRequest returns an option that calls do on the caller's goroutine while
// request n of server, counted from 0, is in progress.
func duringRequest(server *threadServer, n int, do func()) option.RequestOption {
	return option.WithMiddleware(func(request *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if len(server.requests()) == n {
			do()
		}
		return next(request)
	})
}

// compactionWindows are the times at which runCompacting lets the caller act.
var compactionWindows = []string{"while the request is in progress", "while handling the response"}

// runCompacting drives runner to the end. The runner compacts after the reply
// msg_a, with request 2 of server, and act is called in window.
func runCompacting(t *testing.T, window string, thread *BetaMessageThread, runner *turnRunner, server *threadServer, act func()) {
	t.Helper()
	if window == "while the request is in progress" {
		thread.Options = append(thread.Options, duringRequest(server, 2, act))
	}
	runner.run(t, func(message *BetaMessage) {
		switch {
		case message.ID == "msg_a":
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		case message.ID == "msg_compaction" && window == "while handling the response":
			act()
		}
	})
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadDuringACompactionIsSentAfterIt(t *testing.T) {
	for name, compaction := range map[string]struct {
		response string
		history  []string
	}{
		"with a summary":    {compactionResponseJSON, []string{summaryJSON()}},
		"without a summary": {nullSummaryResponseJSON, []string{userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a")}},
	} {
		for _, window := range compactionWindows {
			t.Run(name+"/"+window, func(t *testing.T) {
				forEachRunner(t, func(t *testing.T, stream bool) {
					captureStderr(t)
					thread, runner, server := newRunnerOverThread(t, stream, threadParams(), weatherCall("a"), threadResponse{body: compaction.response}, reply("msg_b", "Sunny"))

					runCompacting(t, window, thread, runner, server, func() { thread.Append(userMessage("In Celsius.")) })

					sent := slices.Concat(compaction.history, []string{userJSON("In Celsius.")})
					requireCreate(t, server.requireRequests(t, 4)[3], messagesOf(sent...))
					requireMessages(t, thread.Messages(), messagesOf(append(sent, assistantJSON("Sunny"))...))
				})
			})
		}
	}
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadDuringACompactionAfterTheFinalReplyStaysStaged(t *testing.T) {
	for _, window := range compactionWindows {
		t.Run(window, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, runner, server := newRunnerOverThread(t, stream, threadParams(), reply("msg_a", "Sunny"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "20 degrees"))

				runCompacting(t, window, thread, runner, server, func() { thread.Append(userMessage("In Celsius.")) })

				require.True(t, runner.IsCompleted())
				server.requireRequests(t, 3)
				requireMessages(t, thread.Messages(), messagesOf(summaryJSON(), userJSON("In Celsius.")))
				require.Equal(t, runner.LastMessage(), thread.LastMessage())
				send(t, thread, noMessages)
				requireCreate(t, server.requireRequests(t, 4)[3], messagesOf(summaryJSON(), userJSON("In Celsius.")))
			})
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadDuringACompactionThatFailsStaysStaged(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), reply("msg_a", "Sunny"), plainError(http.StatusBadRequest, "rejected"))
		thread.Options = append(thread.Options, duringRequest(server, 2, func() {
			thread.Append(userMessage("In Celsius."))
		}))
		runner.mustNext(t)
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})

		_, err := runner.next(context.Background())

		requireStatus(t, err, http.StatusBadRequest)
		requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), assistantJSON("Sunny"), userJSON("In Celsius.")))
	})
}

func TestBetaMessageThread_NewToolRunner_TheCompactionResponseAppendedToTheThreadAgainIsKeptOnce(t *testing.T) {
	for name, test := range map[string]struct {
		firstReply threadResponse
		endsTheRun bool
	}{
		"before the next turn":  {weatherCall("a"), false},
		"after the final reply": {reply("msg_a", "Sunny"), true},
	} {
		t.Run(name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, runner, server := newRunnerOverThread(t, stream, threadParams(), test.firstReply, threadResponse{body: compactionResponseJSON}, reply("msg_b", "Sunny"))

				runner.run(t, func(message *BetaMessage) {
					if message.ID == "msg_a" {
						runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
					}
					thread.Append(message.ToParam())
				})
				if test.endsTheRun {
					send(t, thread, noMessages)
				}

				requireCreate(t, server.requireRequests(t, 4)[3], messagesOf(summaryJSON()))
				requireMessages(t, thread.Messages(), messagesOf(summaryJSON(), assistantJSON("Sunny")))
			})
		})
	}
}

func TestBetaMessageThread_NewToolRunnerStreaming_AMessageAppendedToTheThreadWhileACompactionStreamsIsSentAfterIt(t *testing.T) {
	for _, eventType := range []string{"message_start", "message_stop"} {
		t.Run(eventType, func(t *testing.T) {
			thread, _, server := newRunnerOverThread(t, true, threadParams(), weatherCall("a"), threadResponse{body: compactionResponseJSON}, reply("msg_b", "Sunny"))
			runner := thread.NewToolRunnerStreaming([]BetaTool{&stubBetaTool{name: "weather"}})

			for events, err := range runner.AllStreaming(context.Background()) {
				require.NoError(t, err)
				for event, err := range events {
					require.NoError(t, err)
					if event.Type == eventType && len(server.requests()) == 3 {
						thread.Append(userMessage("In Celsius."))
					}
				}
				if runner.LastMessage().ID == "msg_a" {
					runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
				}
			}

			requireCreate(t, server.requireRequests(t, 4)[3], messagesOf(summaryJSON(), userJSON("In Celsius.")))
			requireMessages(t, thread.Messages(), messagesOf(summaryJSON(), userJSON("In Celsius."), assistantJSON("Sunny")))
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadFollowsTheToolResultsOnTheContinue(t *testing.T) {
	// The runner is given "Briefly." while the caller handles the reply.
	for window, added := range map[string][]string{
		"while the request is in progress": {userJSON("In Celsius."), userJSON("Briefly.")},
		"while handling the reply":         {userJSON("Briefly."), userJSON("In Celsius.")},
		"in a tool":                        {userJSON("Briefly."), userJSON("In Celsius.")},
	} {
		t.Run(window, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, server := newThread(t, weatherThreadParams(), weatherCall("a"), reply("msg_b", "Sunny"))
				thread.Append(userMessage("Weather?"))
				appendToThread := func(in string) {
					if in == window {
						thread.Append(userMessage("In Celsius."))
					}
				}
				thread.Options = append(thread.Options, duringRequest(server, 0, func() { appendToThread("while the request is in progress") }))
				runner := newTurnRunnerOver(thread, stream, &stubBetaTool{name: "weather", run: func(context.Context, json.RawMessage) (string, bool) {
					appendToThread("in a tool")
					return "ok from weather", false
				}})

				runner.run(t, onReply("msg_a", func(*BetaMessage) {
					runner.AppendMessages(userMessage("Briefly."))
					appendToThread("while handling the reply")
				}))

				sent := slices.Concat([]string{weatherResultJSON("a")}, added)
				requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(sent...))
				want := messagesOf(slices.Concat([]string{userJSON("Weather?"), weatherCallJSON("a")}, sent, []string{assistantJSON("Sunny")})...)
				requireMessages(t, thread.Messages(), want)
				requireMessages(t, runner.Messages(), want)
			})
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadAroundAPausedReplyIsSentWithTheRequestThatResumesIt(t *testing.T) {
	for _, window := range []string{"while the request is in progress", "while handling the reply"} {
		t.Run(window, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, runner, server := newRunnerOverThread(t, stream, threadParams(), pausedReply("msg_a", "Searching"), reply("msg_b", "Sunny"))
				appendToThread := func(in string) {
					if in == window {
						thread.Append(userMessage("In Celsius."))
					}
				}
				thread.Options = append(thread.Options, duringRequest(server, 1, func() { appendToThread("while the request is in progress") }))

				runner.run(t, onReply("msg_a", func(*BetaMessage) { appendToThread("while handling the reply") }))

				requireContinue(t, server.requireRequests(t, 3)[2], "msg_a", messagesOf(userJSON("In Celsius.")))
				want := messagesOf(
					userJSON("Hello"), assistantJSON("Hi"),
					userJSON("Weather?"), assistantJSON("Searching"), userJSON("In Celsius."), assistantJSON("Sunny"),
				)
				requireMessages(t, thread.Messages(), want)
				requireMessages(t, runner.Messages(), want)
			})
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AMessageAppendedToTheThreadAfterAReplyThatStoppedForCompactionIsSentWithTheNextRequest(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), threadResponse{body: compactedTurnJSON}, reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_compacted", func(*BetaMessage) {
			thread.Append(userMessage("In Celsius."))
		}))

		requireContinue(t, server.requireRequests(t, 3)[2], "msg_compacted", messagesOf(userJSON("In Celsius.")))
	})
}

func TestBetaMessageThread_NewToolRunner_GoesOnAfterAFailedRequestThatLeftOutThePartialReplyOfTheThread(t *testing.T) {
	for _, window := range []string{"nothing appended", "appended before the run", "appended while the request is in progress", "appended after the error"} {
		t.Run(window, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, server := newThread(t, weatherThreadParams(), textCutOff("msg_0", "par"), plainError(http.StatusBadRequest, "rejected"), reply("msg_a", "Sunny"))
				readStream(thread, turn("Weather?"))
				if window == "appended before the run" {
					thread.Append(userMessage("In Celsius."))
				}
				runner := newTurnRunnerOver(thread, stream, &stubBetaTool{name: "weather"})
				sent := []string{userJSON("Weather?"), assistantJSON("par")}
				if window != "nothing appended" {
					sent = append(sent, userJSON("In Celsius."))
				}
				if window == "appended while the request is in progress" {
					thread.Options = append(thread.Options, duringRequest(server, 1, func() { thread.Append(userMessage("In Celsius.")) }))
				}
				_, err := runner.next(context.Background())
				requireStatus(t, err, http.StatusBadRequest)
				if window == "appended after the error" {
					thread.Append(userMessage("In Celsius."))
				}

				require.Equal(t, "msg_a", runner.mustNext(t).ID)

				requireCreate(t, server.requireRequests(t, 3)[2], messagesOf(sent...))
				requireMessages(t, thread.Messages(), messagesOf(append(sent, assistantJSON("Sunny"))...))
			})
		})
	}
}

func TestBetaMessageThread_NewToolRunner_AHistoryShortenedOnTheThreadIsReplacedWithTheRunners(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, runner, server := newRunnerOverThread(t, stream, threadParams(), weatherCall("a"), reply("msg_b", "Sunny"))

		runner.run(t, onReply("msg_a", func(*BetaMessage) {
			thread.Params.Messages = thread.Params.Messages[:1]
		}))

		requireContinue(t, server.requireRequests(t, 3)[2], "msg_a", messagesOf(weatherResultJSON("a")))
		requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Hi"), userJSON("Weather?"), weatherCallJSON("a"), weatherResultJSON("a"), assistantJSON("Sunny")))
	})
}

func TestBetaMessageThread_NewToolRunner_KeepsTheHelperHeaderAndItsOptions(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Sunny"))
	thread.Append(userMessage("Weather?"))
	runner := thread.NewToolRunner(nil, option.WithHeader("x-run", "1"))

	_, err := runner.RunToCompletion(context.Background())

	require.NoError(t, err)
	request := server.requireRequests(t, 1)[0]
	require.Equal(t, []string{"BetaToolRunner, message-thread"}, request.header.Values("x-stainless-helper"))
	require.Equal(t, "1", request.header.Get("x-run"))
	require.Equal(t, []string{threadsBeta}, request.betas())
}

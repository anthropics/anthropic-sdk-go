package anthropic

import "strings"

// BetaManagedAgentsEventAccumulator folds event_start / event_delta preview
// events into per-event-id agent.message snapshots. The zero value is ready
// to use.
//
//	var previews anthropic.BetaManagedAgentsEventAccumulator
//	for stream.Next() {
//		event := stream.Current()
//		previews.Accumulate(event)
//
//		if event.Type == "event_delta" {
//			fmt.Print(previews.AgentMessageText(event.EventID))
//		}
//	}
type BetaManagedAgentsEventAccumulator struct {
	AgentMessages map[string]BetaManagedAgentsAgentMessageEvent
}

func (acc *BetaManagedAgentsEventAccumulator) Accumulate(event BetaManagedAgentsStreamSessionEventsUnion) {
	if acc == nil {
		return
	}
	if acc.AgentMessages == nil {
		acc.AgentMessages = map[string]BetaManagedAgentsAgentMessageEvent{}
	}

	switch event.Type {
	case "event_start":
		if event.Event.Type == "agent.message" {
			acc.AgentMessages[event.Event.ID] = BetaManagedAgentsAgentMessageEvent{
				ID:   event.Event.ID,
				Type: BetaManagedAgentsAgentMessageEventTypeAgentMessage,
			}
		}

	case "event_delta":
		msg, ok := acc.AgentMessages[event.EventID]
		if !ok {
			return
		}
		idx := int(event.Delta.Index)
		if idx < 0 || idx > len(msg.Content) {
			return
		}
		if idx == len(msg.Content) {
			msg.Content = append(msg.Content, BetaManagedAgentsTextBlock{
				Type: event.Delta.Content.Type,
				Text: event.Delta.Content.Text,
			})
		} else {
			msg.Content[idx].Text += event.Delta.Content.Text
		}
		acc.AgentMessages[event.EventID] = msg

	case "agent.message":
		acc.AgentMessages[event.ID] = event.AsAgentMessage()

	case "span.model_request_end":
		acc.AgentMessages = map[string]BetaManagedAgentsAgentMessageEvent{}
	}
}

func (acc *BetaManagedAgentsEventAccumulator) AgentMessageText(eventID string) string {
	if acc == nil {
		return ""
	}
	var b strings.Builder
	for _, block := range acc.AgentMessages[eventID].Content {
		b.WriteString(block.Text)
	}
	return b.String()
}

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
	// AgentMessages holds one snapshot per event id. Reconciled canonical
	// agent.message events persist across model requests; unreconciled previews
	// are dropped at span.model_request_end because the server will never
	// complete them. The buffered event stream is the authoritative transcript.
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
		start := event.AsEventStart()
		if start.Event.Type == "agent.message" {
			preview := start.Event.AsAgentMessage()
			acc.AgentMessages[preview.ID] = BetaManagedAgentsAgentMessageEvent{
				ID:   preview.ID,
				Type: BetaManagedAgentsAgentMessageEventTypeAgentMessage,
			}
		}

	case "event_delta":
		delta := event.AsEventDelta()
		msg, ok := acc.AgentMessages[delta.EventID]
		if !ok || msg.JSON.ProcessedAt.Valid() {
			return
		}
		idx := int(delta.Delta.Index)
		if idx < 0 || idx > len(msg.Content) {
			return
		}
		if idx == len(msg.Content) {
			var block BetaManagedAgentsAgentMessageEventContentUnion
			_ = block.UnmarshalJSON([]byte(delta.Delta.Content.RawJSON()))
			msg.Content = append(msg.Content, block)
		} else {
			msg.Content[idx].Text += delta.Delta.Content.Text
		}
		acc.AgentMessages[delta.EventID] = msg

	case "agent.message":
		msg := event.AsAgentMessage()
		acc.AgentMessages[msg.ID] = msg

	case "span.model_request_end":
		for id, msg := range acc.AgentMessages {
			if !msg.JSON.ProcessedAt.Valid() {
				delete(acc.AgentMessages, id)
			}
		}
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

package anthropic

import (
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
)

// Accumulate builds up the Message incrementally from a MessageStreamEvent. The Message then can be used as
// any other Message, except with the caveat that the Message.JSON field which normally can be used to inspect
// the JSON sent over the network may not be populated fully.
//
//	message := anthropic.Message{}
//	for stream.Next() {
//		event := stream.Current()
//		message.Accumulate(event)
//	}
func (acc *BetaMessage) Accumulate(event BetaRawMessageStreamEventUnion) error {
	if acc == nil {
		return fmt.Errorf("accumulate: cannot accumlate into nil Message")
	}

	switch event := event.AsAny().(type) {
	case BetaRawMessageStartEvent:
		*acc = event.Message
	case BetaRawMessageDeltaEvent:
		acc.StopReason = event.Delta.StopReason
		acc.StopSequence = event.Delta.StopSequence
		acc.Usage.OutputTokens = event.Usage.OutputTokens
		acc.ContextManagement = event.ContextManagement
	case BetaRawMessageStopEvent:
		accJson, err := json.Marshal(acc)
		if err != nil {
			return fmt.Errorf("error converting content block to JSON: %w", err)
		}
		acc.JSON.raw = string(accJson)
	case BetaRawContentBlockStartEvent:
		acc.Content = append(acc.Content, BetaContentBlockUnion{})
		err := acc.Content[len(acc.Content)-1].UnmarshalJSON([]byte(event.ContentBlock.RawJSON()))
		if err != nil {
			return err
		}
	case BetaRawContentBlockDeltaEvent:
		if len(acc.Content) == 0 {
			return fmt.Errorf("received event of type %s but there was no content block", event.Type)
		}
		cb := &acc.Content[len(acc.Content)-1]
		switch delta := event.Delta.AsAny().(type) {
		case BetaTextDelta:
			cb.Text += delta.Text
		case BetaInputJSONDelta:
			if len(delta.PartialJSON) != 0 {
				if string(cb.Input) == "{}" {
					cb.Input = []byte(delta.PartialJSON)
				} else {
					cb.Input = append(cb.Input, []byte(delta.PartialJSON)...)
				}
			}
		case BetaThinkingDelta:
			cb.Thinking += delta.Thinking
		case BetaSignatureDelta:
			cb.Signature += delta.Signature
		case BetaCitationsDelta:
			citation := BetaTextCitationUnion{}
			err := citation.UnmarshalJSON([]byte(delta.Citation.RawJSON()))
			if err != nil {
				return fmt.Errorf("could not unmarshal citation delta into citation type: %w", err)
			}
			cb.Citations = append(cb.Citations, citation)
		}
	case BetaRawContentBlockStopEvent:
		if len(acc.Content) == 0 {
			return fmt.Errorf("received event of type %s but there was no content block", event.Type)
		}
		contentBlock := &acc.Content[len(acc.Content)-1]
		cbJson, err := json.Marshal(contentBlock)
		if err != nil {
			return fmt.Errorf("error converting content block to JSON: %w", err)
		}
		contentBlock.JSON.raw = string(cbJson)
	}

	return nil
}

// Param converters

func (r BetaContentBlockUnion) ToParam() BetaContentBlockParamUnion {
	switch variant := r.AsAny().(type) {
	case BetaTextBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfText: &p}
	case BetaToolUseBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfToolUse: &p}
	case BetaThinkingBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfThinking: &p}
	case BetaRedactedThinkingBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfRedactedThinking: &p}
	}
	return BetaContentBlockParamUnion{}
}

func (r BetaMessage) ToParam() BetaMessageParam {
	var p BetaMessageParam
	p.Role = BetaMessageParamRole(r.Role)
	p.Content = make([]BetaContentBlockParamUnion, len(r.Content))
	for i, c := range r.Content {
		contentParams := c.ToParam()
		p.Content[i] = contentParams
	}
	return p
}

func (r BetaRedactedThinkingBlock) ToParam() BetaRedactedThinkingBlockParam {
	var p BetaRedactedThinkingBlockParam
	p.Type = r.Type
	p.Data = r.Data
	return p
}

func (r BetaTextBlock) ToParam() BetaTextBlockParam {
	var p BetaTextBlockParam
	p.Type = r.Type
	p.Text = r.Text

	// Distinguish between a nil and zero length slice, since some compatible
	// APIs may not require citations.
	if r.Citations != nil {
		p.Citations = make([]BetaTextCitationParamUnion, len(r.Citations))
	}

	for i, citation := range r.Citations {
		switch citationVariant := citation.AsAny().(type) {
		case BetaCitationCharLocation:
			var citationParam BetaCitationCharLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndCharIndex = citationVariant.EndCharIndex
			citationParam.StartCharIndex = citationVariant.StartCharIndex
			p.Citations[i] = BetaTextCitationParamUnion{OfCharLocation: &citationParam}
		case BetaCitationPageLocation:
			var citationParam BetaCitationPageLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndPageNumber = citationVariant.EndPageNumber
			citationParam.StartPageNumber = citationVariant.StartPageNumber
			p.Citations[i] = BetaTextCitationParamUnion{OfPageLocation: &citationParam}
		case BetaCitationContentBlockLocation:
			var citationParam BetaCitationContentBlockLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			p.Citations[i] = BetaTextCitationParamUnion{OfContentBlockLocation: &citationParam}
		}
	}
	return p
}

func (r BetaThinkingBlock) ToParam() BetaThinkingBlockParam {
	var p BetaThinkingBlockParam
	p.Type = r.Type
	p.Signature = r.Signature
	p.Thinking = r.Thinking
	return p
}

func (r BetaToolUseBlock) ToParam() BetaToolUseBlockParam {
	var p BetaToolUseBlockParam
	p.Type = r.Type
	p.ID = r.ID
	p.Input = r.Input
	p.Name = r.Name
	return p
}

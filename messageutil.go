package anthropic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// CalculateNonStreamingTimeout calculates the appropriate timeout for a non-streaming request
// based on the maximum number of tokens and the model's non-streaming token limit
func CalculateNonStreamingTimeout(maxTokens int, model Model, opts []option.RequestOption) (time.Duration, error) {
	preCfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return 0, fmt.Errorf("error applying request options: %w", err)
	}
	// if the user has set a specific request timeout, use that
	if preCfg.RequestTimeout != 0 {
		return preCfg.RequestTimeout, nil
	}

	maximumTime := time.Hour // 1 hour
	defaultTime := 10 * time.Minute

	expectedTime := time.Duration(float64(maximumTime) * float64(maxTokens) / 128000.0)

	// If the model has a non-streaming token limit and max_tokens exceeds it,
	// or if the expected time exceeds default time, require streaming
	maxNonStreamingTokens, hasLimit := constant.ModelNonStreamingTokens[string(model)]
	if expectedTime > defaultTime || (hasLimit && maxTokens > maxNonStreamingTokens) {
		return 0, fmt.Errorf("streaming is required for operations that may take longer than 10 minutes")
	}

	return defaultTime, nil
}

// Accumulate builds up the Message incrementally from a MessageStreamEvent. The Message then can be used as
// any other Message, except with the caveat that the Message.JSON field which normally can be used to inspect
// the JSON sent over the network may not be populated fully.
//
//	message := anthropic.Message{}
//	for stream.Next() {
//		event := stream.Current()
//		message.Accumulate(event)
//	}
func (acc *Message) Accumulate(event MessageStreamEventUnion) error {
	if acc == nil {
		return fmt.Errorf("accumulate: cannot accumlate into nil Message")
	}

	switch event := event.AsAny().(type) {
	case MessageStartEvent:
		*acc = event.Message
	case MessageDeltaEvent:
		acc.StopReason = event.Delta.StopReason
		acc.StopSequence = event.Delta.StopSequence
		acc.Usage.OutputTokens = event.Usage.OutputTokens
	case MessageStopEvent:
		accJson, err := json.Marshal(acc)
		if err != nil {
			return fmt.Errorf("error converting content block to JSON: %w", err)
		}
		acc.JSON.raw = string(accJson)
	case ContentBlockStartEvent:
		acc.Content = append(acc.Content, ContentBlockUnion{})
		err := acc.Content[len(acc.Content)-1].UnmarshalJSON([]byte(event.ContentBlock.RawJSON()))
		if err != nil {
			return err
		}
	case ContentBlockDeltaEvent:
		if len(acc.Content) == 0 {
			return fmt.Errorf("received event of type %s but there was no content block", event.Type)
		}
		cb := &acc.Content[len(acc.Content)-1]
		switch delta := event.Delta.AsAny().(type) {
		case TextDelta:
			cb.Text += delta.Text
		case InputJSONDelta:
			if len(delta.PartialJSON) != 0 {
				if string(cb.Input) == "{}" {
					cb.Input = []byte(delta.PartialJSON)
				} else {
					cb.Input = append(cb.Input, []byte(delta.PartialJSON)...)
				}
			}
		case ThinkingDelta:
			cb.Thinking += delta.Thinking
		case SignatureDelta:
			cb.Signature += delta.Signature
		case CitationsDelta:
			citation := TextCitationUnion{}
			err := citation.UnmarshalJSON([]byte(delta.Citation.RawJSON()))
			if err != nil {
				return fmt.Errorf("could not unmarshal citation delta into citation type: %w", err)
			}
			cb.Citations = append(cb.Citations, citation)
		}
	case ContentBlockStopEvent:
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

// ToParam converters

func (r ContentBlockUnion) ToParam() ContentBlockParamUnion {
	switch variant := r.AsAny().(type) {
	case TextBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfText: &p}
	case ToolUseBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfToolUse: &p}
	case ThinkingBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfThinking: &p}
	case RedactedThinkingBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfRedactedThinking: &p}
	}
	return ContentBlockParamUnion{}
}

func (r Message) ToParam() MessageParam {
	var p MessageParam
	p.Role = MessageParamRole(r.Role)
	p.Content = make([]ContentBlockParamUnion, len(r.Content))
	for i, c := range r.Content {
		p.Content[i] = c.ToParam()
	}
	return p
}

func (r RedactedThinkingBlock) ToParam() RedactedThinkingBlockParam {
	var p RedactedThinkingBlockParam
	p.Type = r.Type
	p.Data = r.Data
	return p
}

func (r ToolUseBlock) ToParam() ToolUseBlockParam {
	var toolUse ToolUseBlockParam
	toolUse.Type = r.Type
	toolUse.ID = r.ID
	toolUse.Input = r.Input
	toolUse.Name = r.Name
	return toolUse
}

func (r TextBlock) ToParam() TextBlockParam {
	var p TextBlockParam
	p.Type = r.Type
	p.Text = r.Text

	// Distinguish between a nil and zero length slice, since some compatible
	// APIs may not require citations.
	if r.Citations != nil {
		p.Citations = make([]TextCitationParamUnion, len(r.Citations))
	}

	for i, citation := range r.Citations {
		switch citationVariant := citation.AsAny().(type) {
		case CitationCharLocation:
			var citationParam CitationCharLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndCharIndex = citationVariant.EndCharIndex
			citationParam.StartCharIndex = citationVariant.StartCharIndex
			p.Citations[i] = TextCitationParamUnion{OfCharLocation: &citationParam}
		case CitationPageLocation:
			var citationParam CitationPageLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndPageNumber = citationVariant.EndPageNumber
			citationParam.StartPageNumber = citationVariant.StartPageNumber
			p.Citations[i] = TextCitationParamUnion{OfPageLocation: &citationParam}
		case CitationContentBlockLocation:
			var citationParam CitationContentBlockLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			p.Citations[i] = TextCitationParamUnion{OfContentBlockLocation: &citationParam}
		}
	}
	return p
}

func (r ThinkingBlock) ToParam() ThinkingBlockParam {
	var p ThinkingBlockParam
	p.Type = r.Type
	p.Signature = r.Signature
	p.Thinking = r.Thinking
	return p
}

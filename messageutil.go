package anthropic

import (
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
	case MessageStopEvent, ContentBlockStopEvent:
		break
	}

	return nil
}

// ToParam converters

func (r Message) ToParam() MessageParam {
	var p MessageParam
	p.Role = MessageParamRole(r.Role)
	p.Content = make([]ContentBlockParamUnion, len(r.Content))
	for i, c := range r.Content {
		p.Content[i] = c.ToParam()
	}
	return p
}

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
	case ServerToolUseBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfServerToolUse: &p}
	case WebSearchToolResultBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfWebSearchToolResult: &p}
	default:
		panic(fmt.Sprintf("unexpected anthropic.anyContentBlock: %#v", variant))
	}
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
		case CitationsSearchResultLocation:
			var citationParam CitationSearchResultLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.CitedText = citationVariant.CitedText
			citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)

			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			citationParam.Source = citationVariant.Source
			p.Citations[i] = TextCitationParamUnion{OfSearchResultLocation: &citationParam}
		case CitationsWebSearchResultLocation:
			var citationParam CitationWebSearchResultLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.CitedText = citationVariant.CitedText
			citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
			p.Citations[i] = TextCitationParamUnion{OfWebSearchResultLocation: &citationParam}
		default:
			panic(fmt.Sprintf("unexpected anthropic.anyTextCitation: %#v", citationVariant))
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

func (r ServerToolUseBlock) ToParam() ServerToolUseBlockParam {
	var p ServerToolUseBlockParam
	p.Type = r.Type
	p.ID = r.ID
	p.Input = r.Input
	return p
}

func (r WebSearchToolResultBlock) ToParam() WebSearchToolResultBlockParam {
	var p WebSearchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	p.Content = r.Content.ToParam()
	return p
}

func (r WebSearchResultBlock) ToParam() WebSearchResultBlockParam {
	var p WebSearchResultBlockParam
	p.Type = r.Type
	p.EncryptedContent = r.EncryptedContent
	p.Title = r.Title
	p.URL = r.URL
	p.PageAge = paramutil.ToOpt(r.PageAge, r.JSON.PageAge)
	return p
}

func (r WebSearchToolResultBlockContentUnion) ToParam() WebSearchToolResultBlockParamContentUnion {
	var p WebSearchToolResultBlockParamContentUnion

	if len(r.OfWebSearchResultBlockArray) > 0 {
		for _, block := range r.OfWebSearchResultBlockArray {
			p.OfWebSearchToolResultBlockItem = append(p.OfWebSearchToolResultBlockItem, block.ToParam())
		}
		return p
	}

	p.OfRequestWebSearchToolResultError = &WebSearchToolRequestErrorParam{
		ErrorCode: WebSearchToolRequestErrorErrorCode(r.ErrorCode),
	}
	return p
}

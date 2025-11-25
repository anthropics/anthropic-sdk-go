package anthropic

import (
	"fmt"

	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
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
	case BetaRawContentBlockStopEvent, BetaRawMessageStopEvent:
		break
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
	case BetaWebSearchToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfWebSearchToolResult: &p}
	case BetaBashCodeExecutionToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfBashCodeExecutionToolResult: &p}
	case BetaCodeExecutionToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfCodeExecutionToolResult: &p}
	case BetaContainerUploadBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfContainerUpload: &p}
	case BetaMCPToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfMCPToolResult: &p}
	case BetaMCPToolUseBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfMCPToolUse: &p}
	case BetaServerToolUseBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfServerToolUse: &p}
	case BetaTextEditorCodeExecutionToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfTextEditorCodeExecutionToolResult: &p}
	case BetaWebFetchToolResultBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfWebFetchToolResult: &p}
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
		case BetaCitationSearchResultLocation:
			var citationParam BetaCitationSearchResultLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.CitedText = citationVariant.CitedText
			citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			citationParam.Source = citationVariant.Source
			p.Citations[i] = BetaTextCitationParamUnion{OfSearchResultLocation: &citationParam}
		case BetaCitationsWebSearchResultLocation:
			var citationParam BetaCitationWebSearchResultLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.CitedText = citationVariant.CitedText
			citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
			p.Citations[i] = BetaTextCitationParamUnion{OfWebSearchResultLocation: &citationParam}
		default:
			panic(fmt.Sprintf("unexpected anthropic.anyBetaTextCitation: %#v", citationVariant))
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

func (r BetaWebSearchResultBlock) ToParam() BetaWebSearchResultBlockParam {
	var p BetaWebSearchResultBlockParam
	p.Type = r.Type
	p.EncryptedContent = r.EncryptedContent
	p.Title = r.Title
	p.URL = r.URL
	p.PageAge = paramutil.ToOpt(r.PageAge, r.JSON.PageAge)
	return p
}

func (r BetaWebSearchToolResultBlock) ToParam() BetaWebSearchToolResultBlockParam {
	var p BetaWebSearchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID

	if len(r.Content.OfBetaWebSearchResultBlockArray) > 0 {
		for _, block := range r.Content.OfBetaWebSearchResultBlockArray {
			p.Content.OfResultBlock = append(p.Content.OfResultBlock, block.ToParam())
		}
	} else {
		p.Content.OfError = &BetaWebSearchToolRequestErrorParam{
			Type:      r.Content.Type,
			ErrorCode: r.Content.ErrorCode,
		}
	}
	return p
}

func (r BetaWebFetchToolResultBlock) ToParam() BetaWebFetchToolResultBlockParam {
	var p BetaWebFetchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	return p
}

func (r BetaMCPToolUseBlock) ToParam() BetaMCPToolUseBlockParam {
	var p BetaMCPToolUseBlockParam
	p.Type = r.Type
	p.ID = r.ID
	p.Input = r.Input
	p.Name = r.Name
	p.ServerName = r.ServerName
	return p
}

func (r BetaContainerUploadBlock) ToParam() BetaContainerUploadBlockParam {
	var p BetaContainerUploadBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r BetaServerToolUseBlock) ToParam() BetaServerToolUseBlockParam {
	var p BetaServerToolUseBlockParam
	p.Type = r.Type
	p.ID = r.ID
	p.Input = r.Input
	p.Name = BetaServerToolUseBlockParamName(r.Name)
	return p
}

func (r BetaTextEditorCodeExecutionToolResultBlock) ToParam() BetaTextEditorCodeExecutionToolResultBlockParam {
	var p BetaTextEditorCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	p.Content = param.Override[BetaTextEditorCodeExecutionToolResultBlockParamContentUnion](r.Content.RawJSON())
	return p
}

func (r BetaMCPToolResultBlock) ToParam() BetaRequestMCPToolResultBlockParam {
	var p BetaRequestMCPToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.OfString.Valid() {
		p.Content.OfString = paramutil.ToOpt(r.Content.OfString, r.Content.JSON.OfString)
	} else {
		for _, block := range r.Content.OfBetaMCPToolResultBlockContent {
			p.Content.OfBetaMCPToolResultBlockContent = append(p.Content.OfBetaMCPToolResultBlockContent, block.ToParam())
		}
	}
	return p
}

func (r BetaBashCodeExecutionToolResultBlock) ToParam() BetaBashCodeExecutionToolResultBlockParam {
	var p BetaBashCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID

	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestBashCodeExecutionToolResultError = &BetaBashCodeExecutionToolResultErrorParam{
			ErrorCode: BetaBashCodeExecutionToolResultErrorParamErrorCode(r.Content.ErrorCode),
		}
	} else {
		requestBashContentResult := &BetaBashCodeExecutionResultBlockParam{
			ReturnCode: r.Content.ReturnCode,
			Stderr:     r.Content.Stderr,
			Stdout:     r.Content.Stdout,
		}

		for _, block := range r.Content.Content {
			requestBashContentResult.Content = append(requestBashContentResult.Content, block.ToParam())
		}

		p.Content.OfRequestBashCodeExecutionResultBlock = requestBashContentResult
	}

	r.Content.AsResponseBashCodeExecutionResultBlock()
	p.Content = param.Override[BetaBashCodeExecutionToolResultBlockParamContentUnion](r.Content.RawJSON())
	return p
}

func (r BetaBashCodeExecutionOutputBlock) ToParam() BetaBashCodeExecutionOutputBlockParam {
	var p BetaBashCodeExecutionOutputBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r BetaCodeExecutionToolResultBlock) ToParam() BetaCodeExecutionToolResultBlockParam {
	var p BetaCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	p.Content = param.Override[BetaCodeExecutionToolResultBlockParamContentUnion](r.Content.RawJSON())
	return p
}
func (r BetaCodeExecutionOutputBlock) ToParam() BetaCodeExecutionOutputBlockParam {
	var p BetaCodeExecutionOutputBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

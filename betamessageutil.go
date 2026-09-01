package anthropic

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// Accumulate builds up the Message incrementally from a MessageStreamEvent. The Message then can be used as
// any other Message, including the Message.JSON field, which holds the wire JSON with the streamed deltas
// applied and is complete once message_stop arrives.
//
//	message := anthropic.Message{}
//	for stream.Next() {
//		event := stream.Current()
//		message.Accumulate(event)
//	}
func (acc *BetaMessage) Accumulate(event BetaRawMessageStreamEventUnion) error {
	if acc == nil {
		return fmt.Errorf("accumulate: cannot accumulate into nil Message")
	}

	switch event.Type {
	case "message_start":
		*acc = event.Message
	case "message_delta":
		// stop_reason, stop_sequence and stop_details are always sent and null is
		// their final value when the turn carries no such detail.
		acc.StopReason = event.Delta.StopReason
		acc.StopSequence = event.Delta.StopSequence
		acc.StopDetails = event.Delta.StopDetails
		if event.Delta.JSON.Container.Valid() {
			acc.Container = event.Delta.Container
		}
		// Every usage count here is a cumulative whole-message total, so it
		// overwrites rather than adds; the ones that do not apply are omitted, and
		// message_start keeps the last word on those.
		acc.Usage.OutputTokens = event.Usage.OutputTokens
		if event.Usage.JSON.InputTokens.Valid() {
			acc.Usage.InputTokens = event.Usage.InputTokens
		}
		if event.Usage.JSON.CacheCreationInputTokens.Valid() {
			acc.Usage.CacheCreationInputTokens = event.Usage.CacheCreationInputTokens
		}
		if event.Usage.JSON.CacheReadInputTokens.Valid() {
			acc.Usage.CacheReadInputTokens = event.Usage.CacheReadInputTokens
		}
		if event.Usage.JSON.ServerToolUse.Valid() {
			acc.Usage.ServerToolUse = event.Usage.ServerToolUse
		}
		if event.Usage.JSON.OutputTokensDetails.Valid() {
			acc.Usage.OutputTokensDetails = event.Usage.OutputTokensDetails
		}
		if event.Usage.JSON.FallbackCredit.Valid() {
			acc.Usage.FallbackCredit = event.Usage.FallbackCredit
		}
		if event.Usage.JSON.Iterations.Valid() {
			acc.Usage.Iterations = event.Usage.Iterations
		}
		if event.JSON.ContextManagement.Valid() {
			acc.ContextManagement = event.ContextManagement
			acc.JSON.raw, _ = sjson.SetRaw(acc.JSON.raw, "context_management", event.ContextManagement.RawJSON())
		}
		// Sent only when a mid-stream model fallback served the response; it then
		// replaces the list from message_start, even when empty.
		if event.JSON.InputTransformations.Valid() {
			acc.InputTransformations = event.InputTransformations
			acc.JSON.InputTransformations = event.JSON.InputTransformations
			acc.JSON.raw, _ = sjson.SetRaw(acc.JSON.raw, "input_transformations", event.JSON.InputTransformations.Raw())
		}
		acc.JSON.raw = mergeRaw(acc.JSON.raw, "", event.Delta.JSON.raw)
		acc.JSON.raw = mergeRaw(acc.JSON.raw, "usage", event.Usage.RawJSON())
	case "content_block_start":
		// Content blocks start in index order with no gaps: a start event always
		// addresses the slot right after the previous block, even when deltas and
		// stops for still-open blocks interleave after it.
		if event.Index != int64(len(acc.Content)) {
			return fmt.Errorf("received event of type %s for content block at index %d, expected index %d", event.Type, event.Index, len(acc.Content))
		}
		acc.Content = append(acc.Content, BetaContentBlockUnion{})
		err := acc.Content[event.Index].UnmarshalJSON([]byte(event.ContentBlock.RawJSON()))
		if err != nil {
			return err
		}
		// The final hop's fallback block names the model that served the response;
		// non-streaming responses already report that model, so relabel the
		// accumulated snapshot to match. Last block wins on multi-hop chains.
		if fallback, ok := acc.Content[event.Index].AsAny().(BetaFallbackBlock); ok {
			acc.Model = fallback.To.Model
			acc.JSON.raw, _ = sjson.SetRaw(acc.JSON.raw, "model", jsonString(string(fallback.To.Model)))
		}
	case "content_block_delta":
		if err := checkContentBlockIndex(event.Type, event.Index, len(acc.Content)); err != nil {
			return err
		}
		cb := &acc.Content[event.Index]
		switch event.Delta.Type {
		case "text_delta":
			cb.Text += event.Delta.Text
		case "input_json_delta":
			if len(event.Delta.PartialJSON) != 0 {
				if string(cb.Input) == "{}" {
					cb.Input = []byte(event.Delta.PartialJSON)
				} else {
					cb.Input = append(cb.Input, event.Delta.PartialJSON...)
				}
			}
		case "thinking_delta":
			cb.Thinking += event.Delta.Thinking
		case "signature_delta":
			cb.Signature += event.Delta.Signature
		case "citations_delta":
			citation := BetaTextCitationUnion{}
			err := citation.UnmarshalJSON([]byte(event.Delta.Citation.RawJSON()))
			if err != nil {
				return fmt.Errorf("could not unmarshal citation delta into citation type: %w", err)
			}
			cb.Citations = append(cb.Citations, citation)
		case "compaction_delta":
			cb.Content.OfString = event.Delta.Content
			cb.EncryptedContent = event.Delta.EncryptedContent
		}
	case "message_stop":
		// A block whose stop event never arrived has not been refreshed yet.
		for i := range acc.Content {
			refreshBetaContentBlockRaw(&acc.Content[i])
		}
		acc.JSON.raw, _ = sjson.SetRaw(acc.JSON.raw, "content", rawArray(acc.Content))
	case "content_block_stop":
		if err := checkContentBlockIndex(event.Type, event.Index, len(acc.Content)); err != nil {
			return err
		}
		refreshBetaContentBlockRaw(&acc.Content[event.Index])
	}

	return nil
}

// refreshBetaContentBlockRaw overlays the delta-mutated fields onto the block's
// wire JSON, leaving a block that received no deltas byte for byte.
func refreshBetaContentBlockRaw(cb *BetaContentBlockUnion) {
	raw := cb.JSON.raw
	if cb.Text != "" {
		raw, _ = sjson.SetRaw(raw, "text", jsonString(cb.Text))
	}
	if cb.Thinking != "" {
		raw, _ = sjson.SetRaw(raw, "thinking", jsonString(cb.Thinking))
	}
	if cb.Signature != "" {
		raw, _ = sjson.SetRaw(raw, "signature", jsonString(cb.Signature))
	}
	if json.Valid(cb.Input) {
		raw, _ = sjson.SetRaw(raw, "input", string(cb.Input))
	} else if len(cb.Input) > 0 {
		// A cut-off tool call left non-JSON input; empty it so the block marshals.
		cb.Input = json.RawMessage(`{}`)
	}
	if len(cb.Citations) > 0 {
		raw, _ = sjson.SetRaw(raw, "citations", rawArray(cb.Citations))
	}
	if cb.Content.OfString != "" {
		raw, _ = sjson.SetRaw(raw, "content", jsonString(cb.Content.OfString))
	}
	if cb.EncryptedContent != "" {
		raw, _ = sjson.SetRaw(raw, "encrypted_content", jsonString(cb.EncryptedContent))
	}
	cb.JSON.raw = raw
}

// ParseOutput finds the first text content block in the message and unmarshals it
// into dest. This is useful for streaming workflows where you accumulate the message
// first and then parse the structured output.
//
//	var msg anthropic.BetaMessage
//	for stream.Next() { msg.Accumulate(stream.Current()) }
//	msg.ParseOutput(&myStruct)
func (r *BetaMessage) ParseOutput(dest any) error {
	return parseOutputContent(r, dest)
}

// Param converters

func (r BetaContentBlockUnion) ToParam() BetaContentBlockParamUnion {
	return r.AsAny().toParamUnion()
}

func (variant BetaTextBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfText: &p}
}

func (variant BetaToolUseBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfToolUse: &p}
}

func (variant BetaThinkingBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfThinking: &p}
}

func (variant BetaRedactedThinkingBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfRedactedThinking: &p}
}

func (variant BetaWebSearchToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfWebSearchToolResult: &p}
}

func (variant BetaBashCodeExecutionToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfBashCodeExecutionToolResult: &p}
}

func (variant BetaCodeExecutionToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfCodeExecutionToolResult: &p}
}

func (variant BetaContainerUploadBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfContainerUpload: &p}
}

func (variant BetaAdvisorToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfAdvisorToolResult: &p}
}

func (r BetaAdvisorToolResultBlock) ToParam() BetaAdvisorToolResultBlockParam {
	var p BetaAdvisorToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	switch {
	case r.Content.JSON.ErrorCode.Valid():
		p.Content.OfRequestAdvisorToolResultError = &BetaAdvisorToolResultErrorParam{
			ErrorCode: BetaAdvisorToolResultErrorParamErrorCode(r.Content.ErrorCode),
		}
	case r.Content.JSON.EncryptedContent.Valid():
		p.Content.OfRequestAdvisorRedactedResultBlock = &BetaAdvisorRedactedResultBlockParam{
			EncryptedContent: r.Content.EncryptedContent,
			StopReason:       paramutil.ToOpt(r.Content.StopReason, r.Content.JSON.StopReason),
		}
	default:
		p.Content.OfRequestAdvisorResultBlock = &BetaAdvisorResultBlockParam{
			Text:       r.Content.Text,
			StopReason: paramutil.ToOpt(r.Content.StopReason, r.Content.JSON.StopReason),
		}
	}
	return p
}

func (variant BetaMCPToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfMCPToolResult: &p}
}

func (variant BetaMCPToolUseBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfMCPToolUse: &p}
}

func (variant BetaServerToolUseBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfServerToolUse: &p}
}

func (variant BetaTextEditorCodeExecutionToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfTextEditorCodeExecutionToolResult: &p}
}

func (variant BetaWebFetchToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfWebFetchToolResult: &p}
}

func (variant BetaToolSearchToolResultBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfToolSearchToolResult: &p}
}

func (variant BetaCompactionBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfCompaction: &p}
}

func (variant BetaFallbackBlock) toParamUnion() BetaContentBlockParamUnion {
	p := variant.ToParam()
	return BetaContentBlockParamUnion{OfFallback: &p}
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
		p.Citations = make([]BetaTextCitationParamUnion, 0, len(r.Citations))
	}
	for _, citation := range r.Citations {
		p.Citations = append(p.Citations, citation.AsAny().toParamUnion())
	}
	return p
}

func (r BetaCitationCharLocation) toParamUnion() BetaTextCitationParamUnion {
	var citationParam BetaCitationCharLocationParam
	citationParam.Type = r.Type
	citationParam.DocumentTitle = paramutil.ToOpt(r.DocumentTitle, r.JSON.DocumentTitle)
	citationParam.CitedText = r.CitedText
	citationParam.DocumentIndex = r.DocumentIndex
	citationParam.EndCharIndex = r.EndCharIndex
	citationParam.StartCharIndex = r.StartCharIndex
	return BetaTextCitationParamUnion{OfCharLocation: &citationParam}
}

func (citationVariant BetaCitationPageLocation) toParamUnion() BetaTextCitationParamUnion {
	var citationParam BetaCitationPageLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
	citationParam.CitedText = citationVariant.CitedText
	citationParam.DocumentIndex = citationVariant.DocumentIndex
	citationParam.EndPageNumber = citationVariant.EndPageNumber
	citationParam.StartPageNumber = citationVariant.StartPageNumber
	return BetaTextCitationParamUnion{OfPageLocation: &citationParam}
}

func (citationVariant BetaCitationContentBlockLocation) toParamUnion() BetaTextCitationParamUnion {
	var citationParam BetaCitationContentBlockLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
	citationParam.CitedText = citationVariant.CitedText
	citationParam.DocumentIndex = citationVariant.DocumentIndex
	citationParam.EndBlockIndex = citationVariant.EndBlockIndex
	citationParam.StartBlockIndex = citationVariant.StartBlockIndex
	return BetaTextCitationParamUnion{OfContentBlockLocation: &citationParam}
}

func (citationVariant BetaCitationsWebSearchResultLocation) toParamUnion() BetaTextCitationParamUnion {
	var citationParam BetaCitationWebSearchResultLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.CitedText = citationVariant.CitedText
	citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
	citationParam.EncryptedIndex = citationVariant.EncryptedIndex
	citationParam.URL = citationVariant.URL
	return BetaTextCitationParamUnion{OfWebSearchResultLocation: &citationParam}
}

func (citationVariant BetaCitationSearchResultLocation) toParamUnion() BetaTextCitationParamUnion {
	var citationParam BetaCitationSearchResultLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.CitedText = citationVariant.CitedText
	citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
	citationParam.EndBlockIndex = citationVariant.EndBlockIndex
	citationParam.SearchResultIndex = citationVariant.SearchResultIndex
	citationParam.Source = citationVariant.Source
	citationParam.StartBlockIndex = citationVariant.StartBlockIndex
	return BetaTextCitationParamUnion{OfSearchResultLocation: &citationParam}
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
	p.ToolsetName = paramutil.ToOpt(r.ToolsetName, r.JSON.ToolsetName)
	if r.JSON.Caller.Valid() {
		p.Caller = r.Caller.toParam()
	}
	return p
}

// toParam converts a caller. The other three caller unions share this layout
// and convert through it.
func (r BetaToolUseBlockCallerUnion) toParam() BetaToolUseBlockParamCallerUnion {
	var p BetaToolUseBlockParamCallerUnion
	switch v := r.AsAny().(type) {
	case BetaDirectCaller:
		c := v.ToParam()
		p.OfDirect = &c
	case BetaServerToolCaller:
		c := v.ToParam()
		p.OfCodeExecution20250825 = &c
	case BetaServerToolCaller20260120:
		c := v.ToParam()
		p.OfCodeExecution20260120 = &c
	default:
		p = param.Override[BetaToolUseBlockParamCallerUnion](json.RawMessage(r.RawJSON()))
	}
	return p
}

func (r BetaServerToolUseBlockCallerUnion) toParam() BetaServerToolUseBlockParamCallerUnion {
	return BetaServerToolUseBlockParamCallerUnion(BetaToolUseBlockCallerUnion(r).toParam())
}

func (r BetaWebSearchToolResultBlockCallerUnion) toParam() BetaWebSearchToolResultBlockParamCallerUnion {
	return BetaWebSearchToolResultBlockParamCallerUnion(BetaToolUseBlockCallerUnion(r).toParam())
}

func (r BetaWebFetchToolResultBlockCallerUnion) toParam() BetaWebFetchToolResultBlockParamCallerUnion {
	return BetaWebFetchToolResultBlockParamCallerUnion(BetaToolUseBlockCallerUnion(r).toParam())
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

	if r.Content.JSON.OfBetaWebSearchResultBlockArray.Valid() {
		// content is required, so send [] rather than omitting it.
		p.Content.OfResultBlock = make([]BetaWebSearchResultBlockParam, 0, len(r.Content.OfBetaWebSearchResultBlockArray))
		for _, block := range r.Content.OfBetaWebSearchResultBlockArray {
			p.Content.OfResultBlock = append(p.Content.OfResultBlock, block.ToParam())
		}
	} else {
		p.Content.OfError = &BetaWebSearchToolRequestErrorParam{
			Type:      r.Content.Type,
			ErrorCode: r.Content.ErrorCode,
		}
	}
	if r.JSON.Caller.Valid() {
		p.Caller = r.Caller.toParam()
	}
	return p
}

func (r BetaWebFetchToolResultBlock) ToParam() BetaWebFetchToolResultBlockParam {
	var p BetaWebFetchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestWebFetchToolResultError = &BetaWebFetchToolResultErrorBlockParam{
			ErrorCode: r.Content.ErrorCode,
		}
	} else {
		p.Content = param.Override[BetaWebFetchToolResultBlockParamContentUnion](json.RawMessage(r.Content.RawJSON()))
	}
	if r.JSON.Caller.Valid() {
		p.Caller = r.Caller.toParam()
	}
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
	if r.JSON.Caller.Valid() {
		p.Caller = r.Caller.toParam()
	}
	return p
}

func (r BetaTextEditorCodeExecutionToolResultBlock) ToParam() BetaTextEditorCodeExecutionToolResultBlockParam {
	var p BetaTextEditorCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestTextEditorCodeExecutionToolResultError = &BetaTextEditorCodeExecutionToolResultErrorParam{
			ErrorCode:    BetaTextEditorCodeExecutionToolResultErrorParamErrorCode(r.Content.ErrorCode),
			ErrorMessage: paramutil.ToOpt(r.Content.ErrorMessage, r.Content.JSON.ErrorMessage),
		}
	} else {
		p.Content = param.Override[BetaTextEditorCodeExecutionToolResultBlockParamContentUnion](json.RawMessage(r.Content.RawJSON()))
	}
	return p
}

func (r BetaMCPToolResultBlock) ToParam() BetaRequestMCPToolResultBlockParam {
	var p BetaRequestMCPToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	p.IsError = paramutil.ToOpt(r.IsError, r.JSON.IsError)
	if r.Content.JSON.OfString.Valid() {
		p.Content.OfString = paramutil.ToOpt(r.Content.OfString, r.Content.JSON.OfString)
	} else if r.Content.JSON.OfBetaMCPToolResultBlockContent.Valid() {
		// Send [] for an empty result rather than omitting content.
		p.Content.OfBetaMCPToolResultBlockContent = make([]BetaTextBlockParam, 0, len(r.Content.OfBetaMCPToolResultBlockContent))
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
			// content is required, so send [] rather than omitting it.
			Content:    make([]BetaBashCodeExecutionOutputBlockParam, 0, len(r.Content.Content)),
			ReturnCode: r.Content.ReturnCode,
			Stderr:     r.Content.Stderr,
			Stdout:     r.Content.Stdout,
		}
		for _, block := range r.Content.Content {
			requestBashContentResult.Content = append(requestBashContentResult.Content, block.ToParam())
		}
		p.Content.OfRequestBashCodeExecutionResultBlock = requestBashContentResult
	}

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
	// content is required, so send [] rather than omitting it.
	files := make([]BetaCodeExecutionOutputBlockParam, 0, len(r.Content.Content))
	for _, block := range r.Content.Content {
		files = append(files, block.ToParam())
	}
	switch {
	case r.Content.JSON.ErrorCode.Valid():
		p.Content.OfError = &BetaCodeExecutionToolResultErrorParam{
			ErrorCode: r.Content.ErrorCode,
		}
	case r.Content.JSON.EncryptedStdout.Valid():
		p.Content.OfRequestEncryptedCodeExecutionResultBlock = &BetaEncryptedCodeExecutionResultBlockParam{
			Content:         files,
			EncryptedStdout: r.Content.EncryptedStdout,
			ReturnCode:      r.Content.ReturnCode,
			Stderr:          r.Content.Stderr,
		}
	default:
		p.Content.OfResultBlock = &BetaCodeExecutionResultBlockParam{
			Content:    files,
			ReturnCode: r.Content.ReturnCode,
			Stderr:     r.Content.Stderr,
			Stdout:     r.Content.Stdout,
		}
	}
	return p
}

func (r BetaCodeExecutionOutputBlock) ToParam() BetaCodeExecutionOutputBlockParam {
	var p BetaCodeExecutionOutputBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r BetaToolSearchToolResultBlock) ToParam() BetaToolSearchToolResultBlockParam {
	var p BetaToolSearchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestToolSearchToolResultError = &BetaToolSearchToolResultErrorParam{
			ErrorCode:    BetaToolSearchToolResultErrorParamErrorCode(r.Content.ErrorCode),
			ErrorMessage: paramutil.ToOpt(r.Content.ErrorMessage, r.Content.JSON.ErrorMessage),
		}
	} else {
		p.Content.OfRequestToolSearchToolSearchResultBlock = &BetaToolSearchToolSearchResultBlockParam{
			// tool_references is required, so send [] rather than omitting it.
			ToolReferences: make([]BetaToolReferenceBlockParam, 0, len(r.Content.ToolReferences)),
		}
		for _, block := range r.Content.ToolReferences {
			p.Content.OfRequestToolSearchToolSearchResultBlock.ToolReferences = append(
				p.Content.OfRequestToolSearchToolSearchResultBlock.ToolReferences,
				block.ToParam(),
			)
		}
	}
	return p
}

func (r BetaToolReferenceBlock) ToParam() BetaToolReferenceBlockParam {
	var p BetaToolReferenceBlockParam
	p.Type = r.Type
	p.ToolName = r.ToolName
	return p
}

func (r BetaCompactionBlock) ToParam() BetaCompactionBlockParam {
	var p BetaCompactionBlockParam
	p.Type = r.Type
	// A failed compaction has null content, which must stay null rather than
	// become "".
	p.Content = paramutil.ToOpt(r.Content, r.JSON.Content)
	p.EncryptedContent = paramutil.ToOpt(r.EncryptedContent, r.JSON.EncryptedContent)
	return p
}

func (r BetaFallbackBlock) ToParam() BetaFallbackBlockParam {
	var p BetaFallbackBlockParam
	p.Type = r.Type
	p.From = BetaFallbackInfoParam{Model: r.From.Model}
	p.To = BetaFallbackInfoParam{Model: r.To.Model}
	if r.JSON.Trigger.Valid() {
		p.Trigger = json.RawMessage(r.Trigger.RawJSON())
	}
	return p
}

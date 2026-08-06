package anthropic

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
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
func (acc *Message) Accumulate(event MessageStreamEventUnion) error {
	if acc == nil {
		return fmt.Errorf("accumulate: cannot accumulate into nil Message")
	}

	switch event.Type {
	case "message_start":
		*acc = event.Message
	case "message_delta":
		acc.StopReason = event.Delta.StopReason
		acc.StopSequence = event.Delta.StopSequence
		if event.Delta.JSON.StopDetails.Valid() {
			acc.StopDetails = event.Delta.StopDetails
		}
		if event.Delta.JSON.Container.Valid() {
			acc.Container = event.Delta.Container
		}
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
		acc.JSON.raw = mergeRaw(acc.JSON.raw, "", event.Delta.JSON.raw)
		acc.JSON.raw = mergeRaw(acc.JSON.raw, "usage", event.Usage.RawJSON())
	case "content_block_start":
		// Content blocks start in index order with no gaps: a start event always
		// addresses the slot right after the previous block, even when deltas and
		// stops for still-open blocks interleave after it.
		if event.Index != int64(len(acc.Content)) {
			return fmt.Errorf("received event of type %s for content block at index %d, expected index %d", event.Type, event.Index, len(acc.Content))
		}
		acc.Content = append(acc.Content, ContentBlockUnion{})
		err := acc.Content[event.Index].UnmarshalJSON([]byte(event.ContentBlock.RawJSON()))
		if err != nil {
			return err
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
			citation := TextCitationUnion{}
			err := citation.UnmarshalJSON([]byte(event.Delta.Citation.RawJSON()))
			if err != nil {
				return fmt.Errorf("could not unmarshal citation delta into citation type: %w", err)
			}
			cb.Citations = append(cb.Citations, citation)
		}
	case "message_stop":
		// A block whose stop event never arrived has not been refreshed yet.
		for i := range acc.Content {
			refreshContentBlockRaw(&acc.Content[i])
		}
		acc.JSON.raw, _ = sjson.SetRaw(acc.JSON.raw, "content", rawArray(acc.Content))
	case "content_block_stop":
		if err := checkContentBlockIndex(event.Type, event.Index, len(acc.Content)); err != nil {
			return err
		}
		refreshContentBlockRaw(&acc.Content[event.Index])
	}

	return nil
}

// refreshContentBlockRaw overlays the delta-mutated fields onto the block's
// wire JSON, leaving a block that received no deltas byte for byte.
func refreshContentBlockRaw(cb *ContentBlockUnion) {
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
	cb.JSON.raw = raw
}

// jsonString encodes s as the server does, leaving <, > and & literal, which
// sjson.Set does inconsistently.
func jsonString(s string) string {
	var b strings.Builder
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.Encode(s)
	return strings.TrimSuffix(b.String(), "\n")
}

// mergeRaw merges the wire object src into dst at path, skipping the nulls
// message_delta sends for values it is not updating.
func mergeRaw(dst, path, src string) string {
	if path != "" {
		path += "."
	}
	gjson.Parse(src).ForEach(func(key, value gjson.Result) bool {
		if value.Type == gjson.Null {
			return true
		}
		dst, _ = sjson.SetRaw(dst, path+apijson.EscapeSJSONKey(key.String()), value.Raw)
		return true
	})
	return dst
}

// rawArray renders the wire JSON of each item as an array, skipping any item
// with none, such as a block whose start event carried no content_block.
func rawArray[T interface{ RawJSON() string }](items []T) string {
	raws := make([]string, 0, len(items))
	for _, item := range items {
		if raw := item.RawJSON(); raw != "" {
			raws = append(raws, raw)
		}
	}
	return "[" + strings.Join(raws, ",") + "]"
}

// checkContentBlockIndex reports an error if a stream event's index does not
// address one of the numBlocks content blocks accumulated so far. Delta and
// stop events may interleave across open content blocks, so they address
// blocks by index rather than applying to the most recently started block.
func checkContentBlockIndex(eventType string, index int64, numBlocks int) error {
	if index < 0 || index >= int64(numBlocks) {
		return fmt.Errorf("received event of type %s for content block at index %d but there are only %d content blocks", eventType, index, numBlocks)
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
	return r.AsAny().toParamUnion()
}

func (variant TextBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfText: &p}
}

func (variant ToolUseBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfToolUse: &p}
}

func (variant WebSearchToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfWebSearchToolResult: &p}
}

func (variant ServerToolUseBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfServerToolUse: &p}
}

func (variant ThinkingBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfThinking: &p}
}

func (variant RedactedThinkingBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfRedactedThinking: &p}
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

func (citationVariant CitationCharLocation) toParamUnion() TextCitationParamUnion {
	var citationParam CitationCharLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
	citationParam.CitedText = citationVariant.CitedText
	citationParam.DocumentIndex = citationVariant.DocumentIndex
	citationParam.EndCharIndex = citationVariant.EndCharIndex
	citationParam.StartCharIndex = citationVariant.StartCharIndex
	return TextCitationParamUnion{OfCharLocation: &citationParam}
}

func (citationVariant CitationPageLocation) toParamUnion() TextCitationParamUnion {
	var citationParam CitationPageLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
	citationParam.CitedText = citationVariant.CitedText
	citationParam.DocumentIndex = citationVariant.DocumentIndex
	citationParam.EndPageNumber = citationVariant.EndPageNumber
	citationParam.StartPageNumber = citationVariant.StartPageNumber
	return TextCitationParamUnion{OfPageLocation: &citationParam}
}

func (citationVariant CitationContentBlockLocation) toParamUnion() TextCitationParamUnion {
	var citationParam CitationContentBlockLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.DocumentTitle = paramutil.ToOpt(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
	citationParam.CitedText = citationVariant.CitedText
	citationParam.DocumentIndex = citationVariant.DocumentIndex
	citationParam.EndBlockIndex = citationVariant.EndBlockIndex
	citationParam.StartBlockIndex = citationVariant.StartBlockIndex
	return TextCitationParamUnion{OfContentBlockLocation: &citationParam}
}

func (citationVariant CitationsSearchResultLocation) toParamUnion() TextCitationParamUnion {
	var citationParam CitationSearchResultLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.CitedText = citationVariant.CitedText
	citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
	citationParam.EndBlockIndex = citationVariant.EndBlockIndex
	citationParam.SearchResultIndex = citationVariant.SearchResultIndex
	citationParam.Source = citationVariant.Source
	citationParam.StartBlockIndex = citationVariant.StartBlockIndex
	return TextCitationParamUnion{OfSearchResultLocation: &citationParam}
}

func (citationVariant CitationsWebSearchResultLocation) toParamUnion() TextCitationParamUnion {
	var citationParam CitationWebSearchResultLocationParam
	citationParam.Type = citationVariant.Type
	citationParam.CitedText = citationVariant.CitedText
	citationParam.Title = paramutil.ToOpt(citationVariant.Title, citationVariant.JSON.Title)
	citationParam.EncryptedIndex = citationVariant.EncryptedIndex
	citationParam.URL = citationVariant.URL
	return TextCitationParamUnion{OfWebSearchResultLocation: &citationParam}
}

func (r TextBlock) ToParam() TextBlockParam {
	var p TextBlockParam
	p.Type = r.Type
	p.Text = r.Text

	// Distinguish between a nil and zero length slice, since some compatible
	// APIs may not require citations.
	if r.Citations != nil {
		p.Citations = make([]TextCitationParamUnion, 0, len(r.Citations))
	}

	for _, citation := range r.Citations {
		p.Citations = append(p.Citations, citation.AsAny().toParamUnion())
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
	p.Name = ServerToolUseBlockParamName(r.Name)
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
		ErrorCode: WebSearchToolResultErrorCode(r.ErrorCode),
	}
	return p
}

func (variant WebFetchToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfWebFetchToolResult: &p}
}

func (variant CodeExecutionToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfCodeExecutionToolResult: &p}
}

func (variant BashCodeExecutionToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfBashCodeExecutionToolResult: &p}
}

func (variant TextEditorCodeExecutionToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfTextEditorCodeExecutionToolResult: &p}
}

func (variant ToolSearchToolResultBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfToolSearchToolResult: &p}
}

func (variant ContainerUploadBlock) toParamUnion() ContentBlockParamUnion {
	p := variant.ToParam()
	return ContentBlockParamUnion{OfContainerUpload: &p}
}

func (r WebFetchToolResultBlock) ToParam() WebFetchToolResultBlockParam {
	var p WebFetchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	return p
}

func (r ContainerUploadBlock) ToParam() ContainerUploadBlockParam {
	var p ContainerUploadBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r BashCodeExecutionToolResultBlock) ToParam() BashCodeExecutionToolResultBlockParam {
	var p BashCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID

	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestBashCodeExecutionToolResultError = &BashCodeExecutionToolResultErrorParam{
			ErrorCode: BashCodeExecutionToolResultErrorCode(r.Content.ErrorCode),
		}
	} else {
		requestBashContentResult := &BashCodeExecutionResultBlockParam{
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

func (r BashCodeExecutionOutputBlock) ToParam() BashCodeExecutionOutputBlockParam {
	var p BashCodeExecutionOutputBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r CodeExecutionToolResultBlock) ToParam() CodeExecutionToolResultBlockParam {
	var p CodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestCodeExecutionToolResultError = &CodeExecutionToolResultErrorParam{
			ErrorCode: r.Content.ErrorCode,
		}
	} else {
		p.Content.OfRequestCodeExecutionResultBlock = &CodeExecutionResultBlockParam{
			ReturnCode: r.Content.ReturnCode,
			Stderr:     r.Content.Stderr,
			Stdout:     r.Content.Stdout,
		}
		for _, block := range r.Content.Content {
			p.Content.OfRequestCodeExecutionResultBlock.Content = append(p.Content.OfRequestCodeExecutionResultBlock.Content, block.ToParam())
		}
	}
	return p
}

func (r CodeExecutionOutputBlock) ToParam() CodeExecutionOutputBlockParam {
	var p CodeExecutionOutputBlockParam
	p.Type = r.Type
	p.FileID = r.FileID
	return p
}

func (r TextEditorCodeExecutionToolResultBlock) ToParam() TextEditorCodeExecutionToolResultBlockParam {
	var p TextEditorCodeExecutionToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestTextEditorCodeExecutionToolResultError = &TextEditorCodeExecutionToolResultErrorParam{
			ErrorCode:    TextEditorCodeExecutionToolResultErrorCode(r.Content.ErrorCode),
			ErrorMessage: paramutil.ToOpt(r.Content.ErrorMessage, r.Content.JSON.ErrorMessage),
		}
	} else {
		p.Content = param.Override[TextEditorCodeExecutionToolResultBlockParamContentUnion](r.Content.RawJSON())
	}
	return p
}

func (r ToolSearchToolResultBlock) ToParam() ToolSearchToolResultBlockParam {
	var p ToolSearchToolResultBlockParam
	p.Type = r.Type
	p.ToolUseID = r.ToolUseID
	if r.Content.JSON.ErrorCode.Valid() {
		p.Content.OfRequestToolSearchToolResultError = &ToolSearchToolResultErrorParam{
			ErrorCode: ToolSearchToolResultErrorCode(r.Content.ErrorCode),
		}
	} else {
		p.Content.OfRequestToolSearchToolSearchResultBlock = &ToolSearchToolSearchResultBlockParam{}
		for _, block := range r.Content.ToolReferences {
			p.Content.OfRequestToolSearchToolSearchResultBlock.ToolReferences = append(
				p.Content.OfRequestToolSearchToolSearchResultBlock.ToolReferences,
				block.ToParam(),
			)
		}
	}
	return p
}

func (r ToolReferenceBlock) ToParam() ToolReferenceBlockParam {
	var p ToolReferenceBlockParam
	p.Type = r.Type
	p.ToolName = r.ToolName
	return p
}

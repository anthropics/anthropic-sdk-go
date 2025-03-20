// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/resp"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// CompletionService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCompletionService] method instead.
type CompletionService struct {
	Options []option.RequestOption
}

// NewCompletionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewCompletionService(opts ...option.RequestOption) (r CompletionService) {
	r = CompletionService{}
	r.Options = opts
	return
}

// [Legacy] Create a Text Completion.
//
// The Text Completions API is a legacy API. We recommend using the
// [Messages API](https://docs.anthropic.com/en/api/messages) going forward.
//
// Future models and features will not be compatible with Text Completions. See our
// [migration guide](https://docs.anthropic.com/en/api/migrating-from-text-completions-to-messages)
// for guidance in migrating from Text Completions to Messages.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *CompletionService) New(ctx context.Context, body CompletionNewParams, opts ...option.RequestOption) (res *Completion, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/complete"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// [Legacy] Create a Text Completion.
//
// The Text Completions API is a legacy API. We recommend using the
// [Messages API](https://docs.anthropic.com/en/api/messages) going forward.
//
// Future models and features will not be compatible with Text Completions. See our
// [migration guide](https://docs.anthropic.com/en/api/migrating-from-text-completions-to-messages)
// for guidance in migrating from Text Completions to Messages.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *CompletionService) NewStreaming(ctx context.Context, body CompletionNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[Completion]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/complete"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[Completion](ssestream.NewDecoder(raw), err)
}

type Completion struct {
	// Unique object identifier.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id,required"`
	// The resulting completion up to and excluding the stop sequences.
	Completion string `json:"completion,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,required"`
	// The reason that we stopped.
	//
	// This may be one the following values:
	//
	//   - `"stop_sequence"`: we reached a stop sequence — either provided by you via the
	//     `stop_sequences` parameter, or a stop sequence built into the model
	//   - `"max_tokens"`: we exceeded `max_tokens_to_sample` or the model's maximum
	StopReason string `json:"stop_reason,required"`
	// Object type.
	//
	// For Text Completions, this is always `"completion"`.
	Type constant.Completion `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		Completion  resp.Field
		Model       resp.Field
		StopReason  resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Completion) RawJSON() string { return r.JSON.raw }
func (r *Completion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type CompletionNewParams struct {
	// The maximum number of tokens to generate before stopping.
	//
	// Note that our models may stop _before_ reaching this maximum. This parameter
	// only specifies the absolute maximum number of tokens to generate.
	MaxTokensToSample int64 `json:"max_tokens_to_sample,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,omitzero,required"`
	// The prompt that you want Claude to complete.
	//
	// For proper response generation you will need to format your prompt using
	// alternating `\n\nHuman:` and `\n\nAssistant:` conversational turns. For example:
	//
	// ```
	// "\n\nHuman: {userQuestion}\n\nAssistant:"
	// ```
	//
	// See [prompt validation](https://docs.anthropic.com/en/api/prompt-validation) and
	// our guide to
	// [prompt design](https://docs.anthropic.com/en/docs/intro-to-prompting) for more
	// details.
	Prompt string `json:"prompt,required"`
	// Amount of randomness injected into the response.
	//
	// Defaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative and
	// generative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully
	// deterministic.
	Temperature param.Opt[float64] `json:"temperature,omitzero"`
	// Only sample from the top K options for each subsequent token.
	//
	// Used to remove "long tail" low probability responses.
	// [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopK param.Opt[int64] `json:"top_k,omitzero"`
	// Use nucleus sampling.
	//
	// In nucleus sampling, we compute the cumulative distribution over all the options
	// for each subsequent token in decreasing probability order and cut it off once it
	// reaches a particular probability specified by `top_p`. You should either alter
	// `temperature` or `top_p`, but not both.
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopP param.Opt[float64] `json:"top_p,omitzero"`
	// An object describing metadata about the request.
	Metadata MetadataParam `json:"metadata,omitzero"`
	// Sequences that will cause the model to stop generating.
	//
	// Our models stop on `"\n\nHuman:"`, and may include additional built-in stop
	// sequences in the future. By providing the stop_sequences parameter, you may
	// include additional strings that will cause the model to stop generating.
	StopSequences []string `json:"stop_sequences,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CompletionNewParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r CompletionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow CompletionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

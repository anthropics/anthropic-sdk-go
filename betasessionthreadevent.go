// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// BetaSessionThreadEventService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSessionThreadEventService] method instead.
type BetaSessionThreadEventService struct {
	Options []option.RequestOption
}

// NewBetaSessionThreadEventService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaSessionThreadEventService(opts ...option.RequestOption) (r BetaSessionThreadEventService) {
	r = BetaSessionThreadEventService{}
	r.Options = opts
	return
}

// List Session Thread Events
func (r *BetaSessionThreadEventService) List(ctx context.Context, threadID string, params BetaSessionThreadEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsSessionEventUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/threads/%s/events?beta=true", params.SessionID, threadID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Session Thread Events
func (r *BetaSessionThreadEventService) ListAutoPaging(ctx context.Context, threadID string, params BetaSessionThreadEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsSessionEventUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, threadID, params, opts...))
}

// Stream Session Thread Events
func (r *BetaSessionThreadEventService) StreamEvents(ctx context.Context, threadID string, params BetaSessionThreadEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[BetaManagedAgentsStreamSessionThreadEventsUnion]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return ssestream.NewStream[BetaManagedAgentsStreamSessionThreadEventsUnion](nil, err)
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return ssestream.NewStream[BetaManagedAgentsStreamSessionThreadEventsUnion](nil, err)
	}
	path := fmt.Sprintf("v1/sessions/%s/threads/%s/stream?beta=true", params.SessionID, threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[BetaManagedAgentsStreamSessionThreadEventsUnion](ssestream.NewDecoder(raw), err)
}

type BetaSessionThreadEventListParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSessionThreadEventListParams]'s query parameters as
// `url.Values`.
func (r BetaSessionThreadEventListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaSessionThreadEventStreamParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

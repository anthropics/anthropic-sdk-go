package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeTagBundleService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeTagBundleService] method instead.
type BetaOrganizationClaudeTagBundleService struct {
	Options []option.RequestOption
	Targets BetaOrganizationClaudeTagBundleTargetService
}

// NewBetaOrganizationClaudeTagBundleService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationClaudeTagBundleService(opts ...option.RequestOption) (r BetaOrganizationClaudeTagBundleService) {
	r = BetaOrganizationClaudeTagBundleService{}
	r.Options = opts
	r.Targets = NewBetaOrganizationClaudeTagBundleTargetService(opts...)
	return
}

// Get one Bundle. What the Bundle grants (connectors, repositories and
// credentials) is never returned.
func (r *BetaOrganizationClaudeTagBundleService) Get(ctx context.Context, bundleID string, query BetaOrganizationClaudeTagBundleGetParams, opts ...option.RequestOption) (res *BetaClaudeTagBundle, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "experimental-claude-tag-2026-10-21")}, opts...)
	if bundleID == "" {
		err = errors.New("missing required bundle_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/bundles/%s", url.PathEscape(bundleID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.retrieve", "/v1/organizations/claude_tag/bundles/{bundle_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.id", bundleID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the organization's Bundles, newest first by `created_at`. A Bundle that
// Claude Tag manages for itself is not listed.
func (r *BetaOrganizationClaudeTagBundleService) List(ctx context.Context, params BetaOrganizationClaudeTagBundleListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaClaudeTagBundle], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "experimental-claude-tag-2026-10-21"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/claude_tag/bundles"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.list", "/v1/organizations/claude_tag/bundles", "limit", "page"))
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

// List the organization's Bundles, newest first by `created_at`. A Bundle that
// Claude Tag manages for itself is not listed.
func (r *BetaOrganizationClaudeTagBundleService) ListAutoPaging(ctx context.Context, params BetaOrganizationClaudeTagBundleListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaClaudeTagBundle] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// A named set of access that Claude can use in a Slack channel or a Slack
// workspace: connectors, repositories and credentials. Created and edited in
// claude.ai.
type BetaClaudeTagBundle struct {
	// The Bundle's identifier. Starts with `capp_`.
	ID string `json:"id" api:"required"`
	// Time the Bundle was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What the Bundle is for. Null when it has none. Read-only.
	Description string `json:"description" api:"required"`
	// The Bundle's name. Always set. Read-only.
	DisplayName string          `json:"display_name" api:"required"`
	Type        constant.Bundle `json:"type" default:"bundle"`
	// Time of the last change to the Bundle.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		DisplayName respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagBundle) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagBundle) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagBundleGetParams struct {
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeTagBundleListParams struct {
	// Number of items to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagBundleListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagBundleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

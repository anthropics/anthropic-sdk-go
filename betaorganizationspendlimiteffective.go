package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// BetaOrganizationSpendLimitEffectiveService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationSpendLimitEffectiveService] method instead.
type BetaOrganizationSpendLimitEffectiveService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationSpendLimitEffectiveService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationSpendLimitEffectiveService(opts ...option.RequestOption) (r BetaOrganizationSpendLimitEffectiveService) {
	r = BetaOrganizationSpendLimitEffectiveService{}
	r.Options = opts
	return
}

// List each member's effective spend limit and period-to-date spend.
//
// Returns one row per (member, period) the member resolves a spend limit for, with
// the `source` scope the spend limit was inherited from. Paginates by member, so a
// member's periods never split across pages.
func (r *BetaOrganizationSpendLimitEffectiveService) List(ctx context.Context, query BetaOrganizationSpendLimitEffectiveListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaSpendSummary], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/spend_limits/effective?beta=true"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
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

// List each member's effective spend limit and period-to-date spend.
//
// Returns one row per (member, period) the member resolves a spend limit for, with
// the `source` scope the spend limit was inherited from. Paginates by member, so a
// member's periods never split across pages.
func (r *BetaOrganizationSpendLimitEffectiveService) ListAutoPaging(ctx context.Context, query BetaOrganizationSpendLimitEffectiveListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaSpendSummary] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationSpendLimitEffectiveListParams struct {
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Maximum number of members per page. A member's period rows never split across
	// pages, so a page may carry more rows than this. Defaults to `20`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Restrict the report to these limit periods. Omit to return one row per period
	// each member resolves a spend limit for.
	//
	// Any of "daily", "monthly", "weekly".
	Period []string `query:"period,omitzero" json:"-"`
	// Restrict the report to these members, by tagged user ID (`user_...`). At most
	// 100 entries.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationSpendLimitEffectiveListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationSpendLimitEffectiveListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

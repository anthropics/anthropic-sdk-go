package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// BetaOrganizationAnalyticsArtifactService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsArtifactService] method instead.
type BetaOrganizationAnalyticsArtifactService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsArtifactService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsArtifactService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsArtifactService) {
	r = BetaOrganizationAnalyticsArtifactService{}
	r.Options = opts
	return
}

// Get artifact-creation activity for a given day, broken out by MIME type.
//
// Returns the full (`artifact_type`, `is_shared`) cube for the organization;
// `next_page` is null except for grouped queries, which paginate. The cube can be
// broken out per product, per member, or per RBAC group via `group_by[]`, and
// scoped via `filter[]`. Requires an API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsArtifactService) List(ctx context.Context, query BetaOrganizationAnalyticsArtifactListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaAnalyticsArtifactActivity], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/artifacts?beta=true"
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

// Get artifact-creation activity for a given day, broken out by MIME type.
//
// Returns the full (`artifact_type`, `is_shared`) cube for the organization;
// `next_page` is null except for grouped queries, which paginate. The cube can be
// broken out per product, per member, or per RBAC group via `group_by[]`, and
// scoped via `filter[]`. Requires an API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsArtifactService) ListAutoPaging(ctx context.Context, query BetaOrganizationAnalyticsArtifactListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaAnalyticsArtifactActivity] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationAnalyticsArtifactListParams struct {
	// UTC date in YYYY-MM-DD format. The day to get artifact activity for. Data is
	// typically available with a 1-day lag (varies by query; the error for a
	// too-recent date names the latest available day) and may be revised by a few
	// percent over the following days. No earlier than 2026-01-01.
	Date time.Time `query:"date" api:"required" format:"date" json:"-"`
	// Maximum rows to return (1-1000, default 100). The ungrouped artifact-type cube
	// is finite and returned in full; `limit` is the page size only when `group_by[]`
	// multiplies the cube.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page` field. Only valid with
	// `group_by[]` — the ungrouped cube is never paginated.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filters as `dimension:value`, e.g. `filter[]=rbac_group_id:{id}`. Repeat the
	// param for OR within a dimension and across dimensions for AND. Supported
	// dimensions on this endpoint: `artifact_type`, `is_shared`, `product`,
	// `rbac_group_id`, `user_id`. Value forms: `artifact_type` is a canonical artifact
	// MIME type (e.g. `text/markdown`) or `other`; `is_shared` is `true` or `false`;
	// `product` is `chat`, `claude_code`, or `cowork` (the surfaces that create
	// artifacts); `rbac_group_id` takes the tagged id (`rbac_group_...`, as emitted in
	// responses and by the spend-limits API) or a bare group UUID, and matches users
	// who held the group at any point during each covered UTC day (time-of-usage
	// attribution); `user_id` takes a tagged user id (`user_...`), as emitted in
	// responses. An unsupported dimension returns 400. At most 100 entries.
	Filter []string `query:"filter,omitzero" json:"-"`
	// Dimensions to break results out by: `product`, `user_id` and/or `rbac_group_id`.
	// The ungrouped artifact-type cube is finite and returned in full; grouped queries
	// multiply the cube and paginate via `next_page`. `product` takes the values
	// `chat`, `claude_code`, or `cowork` (the surfaces that create artifacts).
	// `rbac_group_id` attributes a user to every group they held at any point during
	// the requested UTC day, so grouped rows are not an exclusive partition. At most
	// 100 entries.
	//
	// Any of "product", "rbac_group_id", "user_id".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsArtifactListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsArtifactListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

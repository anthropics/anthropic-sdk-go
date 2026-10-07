package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationSpendLimitService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationSpendLimitService] method instead.
type BetaOrganizationSpendLimitService struct {
	Options          []option.RequestOption
	Effective        BetaOrganizationSpendLimitEffectiveService
	IncreaseRequests BetaOrganizationSpendLimitIncreaseRequestService
}

// NewBetaOrganizationSpendLimitService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationSpendLimitService(opts ...option.RequestOption) (r BetaOrganizationSpendLimitService) {
	r = BetaOrganizationSpendLimitService{}
	r.Options = opts
	r.Effective = NewBetaOrganizationSpendLimitEffectiveService(opts...)
	r.IncreaseRequests = NewBetaOrganizationSpendLimitIncreaseRequestService(opts...)
	return
}

// Retrieve a spend limit by ID.
func (r *BetaOrganizationSpendLimitService) Get(ctx context.Context, spendLimitID string, opts ...option.RequestOption) (res *BetaSpendLimit, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitID == "" {
		err = errors.New("missing required spend_limit_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limits/%s?beta=true", url.PathEscape(spendLimitID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the organization's spend limits.
//
// A Claude Console organization's limits come in an order that is stable across
// pages. A Claude Enterprise organization's are grouped by scope type, in the
// order `organization`, `seat_tier`, `rbac_group`, `organization_service`, `user`;
// within a type they come in a fixed order that is not creation order. Listing
// Claude Console limits is in an early access preview. To request access, contact
// your Anthropic account team.
func (r *BetaOrganizationSpendLimitService) List(ctx context.Context, params BetaOrganizationSpendLimitListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaSpendLimit], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "spend-limit-reads-2026-09-26"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/spend_limits?beta=true"
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

// List the organization's spend limits.
//
// A Claude Console organization's limits come in an order that is stable across
// pages. A Claude Enterprise organization's are grouped by scope type, in the
// order `organization`, `seat_tier`, `rbac_group`, `organization_service`, `user`;
// within a type they come in a fixed order that is not creation order. Listing
// Claude Console limits is in an early access preview. To request access, contact
// your Anthropic account team.
func (r *BetaOrganizationSpendLimitService) ListAutoPaging(ctx context.Context, params BetaOrganizationSpendLimitListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaSpendLimit] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete a spend limit.
//
// For a Claude Enterprise organization, this deletes a per-user override, and the
// member falls back to any inherited spend limit at that period. Its seat-tier,
// group, and organization-level rows cannot be deleted via this endpoint. A Claude
// Console organization deletes its organization and workspace limits. Deleting
// them through the API is in an early access preview.
func (r *BetaOrganizationSpendLimitService) Delete(ctx context.Context, spendLimitID string, opts ...option.RequestOption) (res *BetaOrganizationSpendLimitDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitID == "" {
		err = errors.New("missing required spend_limit_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limits/%s?beta=true", url.PathEscape(spendLimitID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Set a spend limit.
//
// Upsert keyed on (scope, period): setting a limit that already exists overwrites
// it in place. A Claude Enterprise organization sets `user` limits. Its seat-tier,
// group, and organization-level defaults are configured in claude.ai. A Claude
// Console organization sets `organization` and `workspace` limits, which are
// monthly and always carry an amount. Setting those limits is in an early access
// preview. To request access, contact your Anthropic account team.
func (r *BetaOrganizationSpendLimitService) Set(ctx context.Context, params BetaOrganizationSpendLimitSetParams, opts ...option.RequestOption) (res *BetaSpendLimit, err error) {
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
	path := "v1/organizations/spend_limits?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// A configured spend limit: a cap on metered spend for one scope and period.
type BetaSpendLimit struct {
	// Unique tagged ID of the spend limit (`spl_...`).
	ID string `json:"id" api:"required"`
	// Limit amount as a non-negative integer decimal string in the minor unit of
	// `currency` (cents for USD): "50000" is $500.00. `null` means no numeric cap is
	// configured at this scope — see the effective report for whether a limit applies.
	Amount string `json:"amount" api:"required"`
	// RFC 3339 datetime at which the spend limit was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ISO 4217 code of the organization's billing currency; the unit for `amount`.
	Currency string `json:"currency" api:"required"`
	// Read-only. `false` when extra usage is switched off for this organization
	// (`organization` limit) or for this member (`user` limit); `amount` is kept and
	// applies again when it's switched back on. Always `true` for other limits.
	IsEnabled bool `json:"is_enabled" api:"required"`
	// Length of the window the limit resets over. `amount` caps spend within each
	// period.
	//
	// Any of "daily", "monthly", "weekly".
	Period BetaSpendLimitPeriod `json:"period" api:"required"`
	// What the limit applies to. A tagged union on `type`; each variant carries the
	// identifier for its scope.
	Scope BetaSpendLimitScopeUnion `json:"scope" api:"required"`
	// Object type. Always `spend_limit`.
	Type constant.SpendLimit `json:"type" default:"spend_limit"`
	// RFC 3339 datetime at which the spend limit was last modified.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Amount      respjson.Field
		CreatedAt   respjson.Field
		Currency    respjson.Field
		IsEnabled   respjson.Field
		Period      respjson.Field
		Scope       respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimit) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendLimitScopeUnion contains all possible properties and values from
// [BetaSpendLimitUserScope], [BetaSpendLimitSeatTierScope],
// [BetaSpendLimitRBACGroupScope], [BetaSpendLimitOrganizationServiceScope],
// [BetaSpendLimitOrganizationScope], [BetaSpendLimitWorkspaceScope].
//
// Use the [BetaSpendLimitScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendLimitScopeUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [BetaSpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaSpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [BetaSpendLimitWorkspaceScope].
	WorkspaceID string `json:"workspace_id"`
	JSON        struct {
		Type        respjson.Field
		UserID      respjson.Field
		SeatTier    respjson.Field
		RBACGroupID respjson.Field
		Service     respjson.Field
		WorkspaceID respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaSpendLimitScope is implemented by each variant of
// [BetaSpendLimitScopeUnion] to add type safety for the return type of
// [BetaSpendLimitScopeUnion.AsAny]
type anyBetaSpendLimitScope interface {
	implBetaSpendLimitScopeUnion()
}

func (BetaSpendLimitUserScope) implBetaSpendLimitScopeUnion()                {}
func (BetaSpendLimitSeatTierScope) implBetaSpendLimitScopeUnion()            {}
func (BetaSpendLimitRBACGroupScope) implBetaSpendLimitScopeUnion()           {}
func (BetaSpendLimitOrganizationServiceScope) implBetaSpendLimitScopeUnion() {}
func (BetaSpendLimitOrganizationScope) implBetaSpendLimitScopeUnion()        {}
func (BetaSpendLimitWorkspaceScope) implBetaSpendLimitScopeUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendLimitScopeUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserScope:
//	case anthropic.BetaSpendLimitSeatTierScope:
//	case anthropic.BetaSpendLimitRBACGroupScope:
//	case anthropic.BetaSpendLimitOrganizationServiceScope:
//	case anthropic.BetaSpendLimitOrganizationScope:
//	case anthropic.BetaSpendLimitWorkspaceScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendLimitScopeUnion) AsAny() anyBetaSpendLimitScope {
	switch u.Type {
	case "user":
		return u.AsUser()
	case "seat_tier":
		return u.AsSeatTier()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_service":
		return u.AsOrganizationService()
	case "organization":
		return u.AsOrganization()
	case "workspace":
		return u.AsWorkspace()
	}
	return nil
}

func (u BetaSpendLimitScopeUnion) AsUser() (v BetaSpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitScopeUnion) AsSeatTier() (v BetaSpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitScopeUnion) AsRBACGroup() (v BetaSpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitScopeUnion) AsOrganizationService() (v BetaSpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitScopeUnion) AsOrganization() (v BetaSpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitScopeUnion) AsWorkspace() (v BetaSpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendLimitScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendLimitScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSpendLimitOrganizationScope struct {
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitOrganizationScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitOrganizationScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaSpendLimitOrganizationScope to a
// BetaSpendLimitOrganizationScopeParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaSpendLimitOrganizationScopeParam.Overrides()
func (r BetaSpendLimitOrganizationScope) ToParam() BetaSpendLimitOrganizationScopeParam {
	return param.Override[BetaSpendLimitOrganizationScopeParam](json.RawMessage(r.RawJSON()))
}

func NewBetaSpendLimitOrganizationScopeParam() BetaSpendLimitOrganizationScopeParam {
	return BetaSpendLimitOrganizationScopeParam{
		Type: "organization",
	}
}

// This struct has a constant value, construct it with
// [NewBetaSpendLimitOrganizationScopeParam].
type BetaSpendLimitOrganizationScopeParam struct {
	Type constant.Organization `json:"type" default:"organization"`
	paramObj
}

func (r BetaSpendLimitOrganizationScopeParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaSpendLimitOrganizationScopeParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSpendLimitOrganizationScopeParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSpendLimitOrganizationServiceScope struct {
	Service string                       `json:"service" api:"required"`
	Type    constant.OrganizationService `json:"type" default:"organization_service"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Service     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitOrganizationServiceScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitOrganizationServiceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSpendLimitPeriod string

const (
	BetaSpendLimitPeriodDaily   BetaSpendLimitPeriod = "daily"
	BetaSpendLimitPeriodMonthly BetaSpendLimitPeriod = "monthly"
	BetaSpendLimitPeriodWeekly  BetaSpendLimitPeriod = "weekly"
)

type BetaSpendLimitRBACGroupScope struct {
	RBACGroupID string             `json:"rbac_group_id" api:"required"`
	Type        constant.RBACGroup `json:"type" default:"rbac_group"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RBACGroupID respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitRBACGroupScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitRBACGroupScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A scoped Admin API key acting on behalf of the organization.
type BetaSpendLimitScopedAPIKeyActor struct {
	ScopedAPIKeyID string                     `json:"scoped_api_key_id" api:"required"`
	Type           constant.ScopedAPIKeyActor `json:"type" default:"scoped_api_key_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ScopedAPIKeyID respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitScopedAPIKeyActor) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitScopedAPIKeyActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSpendLimitSeatTierScope struct {
	SeatTier string            `json:"seat_tier" api:"required"`
	Type     constant.SeatTier `json:"type" default:"seat_tier"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SeatTier    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitSeatTierScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitSeatTierScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user within the organization. `name` and `email_address` are null when the
// underlying account is unavailable or has been deleted; `deleted` is true only
// for deleted accounts.
type BetaSpendLimitUserActor struct {
	// True only when the underlying account has been deleted.
	Deleted bool `json:"deleted" api:"required"`
	// The user's email address. Null when the account is unavailable or has been
	// deleted.
	EmailAddress string `json:"email_address" api:"required"`
	// The user's current display name. Null when the account is unavailable, has been
	// deleted, or has no name set.
	Name string `json:"name" api:"required"`
	// Actor type. Always `user_actor`.
	Type constant.UserActor `json:"type" default:"user_actor"`
	// Tagged ID of the user.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted      respjson.Field
		EmailAddress respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Scope selecting a single member of the organization.
type BetaSpendLimitUserScope struct {
	// Scope type. Always `user` for this scope.
	Type constant.User `json:"type" default:"user"`
	// Tagged ID of the member the spend limit applies to.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitUserScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitUserScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaSpendLimitUserScope to a BetaSpendLimitUserScopeParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaSpendLimitUserScopeParam.Overrides()
func (r BetaSpendLimitUserScope) ToParam() BetaSpendLimitUserScopeParam {
	return param.Override[BetaSpendLimitUserScopeParam](json.RawMessage(r.RawJSON()))
}

// Scope selecting a single member of the organization.
//
// The properties Type, UserID are required.
type BetaSpendLimitUserScopeParam struct {
	// Tagged ID of the member the spend limit applies to.
	UserID string `json:"user_id" api:"required"`
	// Scope type. Always `user` for this scope.
	//
	// This field can be elided, and will marshal its zero value as "user".
	Type constant.User `json:"type" default:"user"`
	paramObj
}

func (r BetaSpendLimitUserScopeParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaSpendLimitUserScopeParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSpendLimitUserScopeParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Scope selecting one workspace of a Claude Console organization.
type BetaSpendLimitWorkspaceScope struct {
	// Scope type. Always `workspace` for this scope.
	Type constant.Workspace `json:"type" default:"workspace"`
	// Tagged ID of the workspace the spend limit applies to.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitWorkspaceScope) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitWorkspaceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaSpendLimitWorkspaceScope to a
// BetaSpendLimitWorkspaceScopeParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaSpendLimitWorkspaceScopeParam.Overrides()
func (r BetaSpendLimitWorkspaceScope) ToParam() BetaSpendLimitWorkspaceScopeParam {
	return param.Override[BetaSpendLimitWorkspaceScopeParam](json.RawMessage(r.RawJSON()))
}

// Scope selecting one workspace of a Claude Console organization.
//
// The properties Type, WorkspaceID are required.
type BetaSpendLimitWorkspaceScopeParam struct {
	// Tagged ID of the workspace the spend limit applies to.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Scope type. Always `workspace` for this scope.
	//
	// This field can be elided, and will marshal its zero value as "workspace".
	Type constant.Workspace `json:"type" default:"workspace"`
	paramObj
}

func (r BetaSpendLimitWorkspaceScopeParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaSpendLimitWorkspaceScopeParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSpendLimitWorkspaceScopeParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-member effective-limit report row (`GET /spend_limits/effective`).
type BetaSpendSummary struct {
	Actor BetaSpendSummaryActorUnion `json:"actor" api:"required"`
	// Effective limit amount as a non-negative integer decimal string in the minor
	// unit of `currency` (cents for USD). `null` means no limit applies for this row's
	// `period` — each period resolves independently, so another period may still cap
	// this member.
	Amount string `json:"amount" api:"required"`
	// ISO 4217 code of the organization's billing currency; the unit for `amount` and
	// `period_to_date_spend`.
	Currency string `json:"currency" api:"required"`
	// Period this row's effective limit and spend are reported for.
	//
	// Any of "daily", "monthly", "weekly".
	Period BetaSpendLimitPeriod `json:"period" api:"required"`
	// The member's spend so far in the current period, as a non-negative decimal
	// string in the minor unit of `currency` (cents for USD). May carry fractional
	// minor units up to three decimal places (e.g. `"12050.5"`) — metered usage is not
	// rounded to whole cents. Reads as `"0"` when the spend reading is temporarily
	// unavailable.
	PeriodToDateSpend string                      `json:"period_to_date_spend" api:"required"`
	Scope             BetaSpendSummaryScopeUnion  `json:"scope" api:"required"`
	Source            BetaSpendSummarySourceUnion `json:"source" api:"required"`
	SpendLimitID      string                      `json:"spend_limit_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Actor             respjson.Field
		Amount            respjson.Field
		Currency          respjson.Field
		Period            respjson.Field
		PeriodToDateSpend respjson.Field
		Scope             respjson.Field
		Source            respjson.Field
		SpendLimitID      respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendSummary) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendSummaryActorUnion contains all possible properties and values from
// [BetaSpendLimitUserActor], [BetaSpendLimitScopedAPIKeyActor].
//
// Use the [BetaSpendSummaryActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendSummaryActorUnion struct {
	// This field is from variant [BetaSpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaSpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaSpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitScopedAPIKeyActor].
	ScopedAPIKeyID string `json:"scoped_api_key_id"`
	JSON           struct {
		Deleted        respjson.Field
		EmailAddress   respjson.Field
		Name           respjson.Field
		Type           respjson.Field
		UserID         respjson.Field
		ScopedAPIKeyID respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaSpendSummaryActor is implemented by each variant of
// [BetaSpendSummaryActorUnion] to add type safety for the return type of
// [BetaSpendSummaryActorUnion.AsAny]
type anyBetaSpendSummaryActor interface {
	implBetaSpendSummaryActorUnion()
}

func (BetaSpendLimitUserActor) implBetaSpendSummaryActorUnion()         {}
func (BetaSpendLimitScopedAPIKeyActor) implBetaSpendSummaryActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendSummaryActorUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserActor:
//	case anthropic.BetaSpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendSummaryActorUnion) AsAny() anyBetaSpendSummaryActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u BetaSpendSummaryActorUnion) AsUserActor() (v BetaSpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryActorUnion) AsScopedAPIKeyActor() (v BetaSpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendSummaryActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendSummaryActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendSummaryScopeUnion contains all possible properties and values from
// [BetaSpendLimitUserScope], [BetaSpendLimitSeatTierScope],
// [BetaSpendLimitRBACGroupScope], [BetaSpendLimitOrganizationServiceScope],
// [BetaSpendLimitOrganizationScope], [BetaSpendLimitWorkspaceScope].
//
// Use the [BetaSpendSummaryScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendSummaryScopeUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [BetaSpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaSpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [BetaSpendLimitWorkspaceScope].
	WorkspaceID string `json:"workspace_id"`
	JSON        struct {
		Type        respjson.Field
		UserID      respjson.Field
		SeatTier    respjson.Field
		RBACGroupID respjson.Field
		Service     respjson.Field
		WorkspaceID respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaSpendSummaryScope is implemented by each variant of
// [BetaSpendSummaryScopeUnion] to add type safety for the return type of
// [BetaSpendSummaryScopeUnion.AsAny]
type anyBetaSpendSummaryScope interface {
	implBetaSpendSummaryScopeUnion()
}

func (BetaSpendLimitUserScope) implBetaSpendSummaryScopeUnion()                {}
func (BetaSpendLimitSeatTierScope) implBetaSpendSummaryScopeUnion()            {}
func (BetaSpendLimitRBACGroupScope) implBetaSpendSummaryScopeUnion()           {}
func (BetaSpendLimitOrganizationServiceScope) implBetaSpendSummaryScopeUnion() {}
func (BetaSpendLimitOrganizationScope) implBetaSpendSummaryScopeUnion()        {}
func (BetaSpendLimitWorkspaceScope) implBetaSpendSummaryScopeUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendSummaryScopeUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserScope:
//	case anthropic.BetaSpendLimitSeatTierScope:
//	case anthropic.BetaSpendLimitRBACGroupScope:
//	case anthropic.BetaSpendLimitOrganizationServiceScope:
//	case anthropic.BetaSpendLimitOrganizationScope:
//	case anthropic.BetaSpendLimitWorkspaceScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendSummaryScopeUnion) AsAny() anyBetaSpendSummaryScope {
	switch u.Type {
	case "user":
		return u.AsUser()
	case "seat_tier":
		return u.AsSeatTier()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_service":
		return u.AsOrganizationService()
	case "organization":
		return u.AsOrganization()
	case "workspace":
		return u.AsWorkspace()
	}
	return nil
}

func (u BetaSpendSummaryScopeUnion) AsUser() (v BetaSpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryScopeUnion) AsSeatTier() (v BetaSpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryScopeUnion) AsRBACGroup() (v BetaSpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryScopeUnion) AsOrganizationService() (v BetaSpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryScopeUnion) AsOrganization() (v BetaSpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummaryScopeUnion) AsWorkspace() (v BetaSpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendSummaryScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendSummaryScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendSummarySourceUnion contains all possible properties and values from
// [BetaSpendLimitUserScope], [BetaSpendLimitSeatTierScope],
// [BetaSpendLimitRBACGroupScope], [BetaSpendLimitOrganizationServiceScope],
// [BetaSpendLimitOrganizationScope], [BetaSpendLimitWorkspaceScope].
//
// Use the [BetaSpendSummarySourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendSummarySourceUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [BetaSpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaSpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [BetaSpendLimitWorkspaceScope].
	WorkspaceID string `json:"workspace_id"`
	JSON        struct {
		Type        respjson.Field
		UserID      respjson.Field
		SeatTier    respjson.Field
		RBACGroupID respjson.Field
		Service     respjson.Field
		WorkspaceID respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaSpendSummarySource is implemented by each variant of
// [BetaSpendSummarySourceUnion] to add type safety for the return type of
// [BetaSpendSummarySourceUnion.AsAny]
type anyBetaSpendSummarySource interface {
	implBetaSpendSummarySourceUnion()
}

func (BetaSpendLimitUserScope) implBetaSpendSummarySourceUnion()                {}
func (BetaSpendLimitSeatTierScope) implBetaSpendSummarySourceUnion()            {}
func (BetaSpendLimitRBACGroupScope) implBetaSpendSummarySourceUnion()           {}
func (BetaSpendLimitOrganizationServiceScope) implBetaSpendSummarySourceUnion() {}
func (BetaSpendLimitOrganizationScope) implBetaSpendSummarySourceUnion()        {}
func (BetaSpendLimitWorkspaceScope) implBetaSpendSummarySourceUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendSummarySourceUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserScope:
//	case anthropic.BetaSpendLimitSeatTierScope:
//	case anthropic.BetaSpendLimitRBACGroupScope:
//	case anthropic.BetaSpendLimitOrganizationServiceScope:
//	case anthropic.BetaSpendLimitOrganizationScope:
//	case anthropic.BetaSpendLimitWorkspaceScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendSummarySourceUnion) AsAny() anyBetaSpendSummarySource {
	switch u.Type {
	case "user":
		return u.AsUser()
	case "seat_tier":
		return u.AsSeatTier()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_service":
		return u.AsOrganizationService()
	case "organization":
		return u.AsOrganization()
	case "workspace":
		return u.AsWorkspace()
	}
	return nil
}

func (u BetaSpendSummarySourceUnion) AsUser() (v BetaSpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummarySourceUnion) AsSeatTier() (v BetaSpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummarySourceUnion) AsRBACGroup() (v BetaSpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummarySourceUnion) AsOrganizationService() (v BetaSpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummarySourceUnion) AsOrganization() (v BetaSpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendSummarySourceUnion) AsWorkspace() (v BetaSpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendSummarySourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendSummarySourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationSpendLimitDeleteResponse struct {
	ID   string                     `json:"id" api:"required"`
	Type constant.SpendLimitDeleted `json:"type" default:"spend_limit_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationSpendLimitDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationSpendLimitDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationSpendLimitListParams struct {
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Maximum number of limits per page. Defaults to `20`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Return only limits with these scope types. A Claude Console organization has
	// `organization` and `workspace` limits; a Claude Enterprise organization has
	// `organization`, `seat_tier`, `rbac_group`, `organization_service` and `user`
	// limits. Omit for all.
	//
	// Any of "organization", "organization_service", "rbac_group", "seat_tier",
	// "user", "workspace".
	ScopeType []string `query:"scope_type,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `spend-limit-reads-2026-09-26` in
	// this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationSpendLimitListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationSpendLimitListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationSpendLimitSetParams struct {
	// Limit amount as a non-negative integer decimal string in the minor unit of the
	// organization's billing currency (cents for USD): "50000" is $500.00. `null` sets
	// an explicit no-limit override for this scope and `period` only — each period
	// resolves independently, so caps for other periods still apply.
	Amount param.Opt[string] `json:"amount,omitzero" api:"required"`
	// What the limit applies to. Claude Enterprise organizations set `user` limits.
	// Claude Console organizations set `organization` and `workspace` limits. Any
	// other combination returns 400. Setting `organization` and `workspace` limits
	// through the API is in an early access preview. To request access, contact your
	// Anthropic account team.
	Scope BetaOrganizationSpendLimitSetParamsScopeUnion `json:"scope,omitzero" api:"required"`
	// Any of "daily", "monthly", "weekly".
	Period BetaSpendLimitPeriod `json:"period,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationSpendLimitSetParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationSpendLimitSetParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationSpendLimitSetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationSpendLimitSetParamsScopeUnion struct {
	OfUser         *BetaSpendLimitUserScopeParam         `json:",omitzero,inline"`
	OfOrganization *BetaSpendLimitOrganizationScopeParam `json:",omitzero,inline"`
	OfWorkspace    *BetaSpendLimitWorkspaceScopeParam    `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationSpendLimitSetParamsScopeUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUser, u.OfOrganization, u.OfWorkspace)
}
func (u *BetaOrganizationSpendLimitSetParamsScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationSpendLimitSetParamsScopeUnion) asAny() any {
	if !param.IsOmitted(u.OfUser) {
		return u.OfUser
	} else if !param.IsOmitted(u.OfOrganization) {
		return u.OfOrganization
	} else if !param.IsOmitted(u.OfWorkspace) {
		return u.OfWorkspace
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationSpendLimitSetParamsScopeUnion) GetUserID() *string {
	if vt := u.OfUser; vt != nil {
		return &vt.UserID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationSpendLimitSetParamsScopeUnion) GetWorkspaceID() *string {
	if vt := u.OfWorkspace; vt != nil {
		return &vt.WorkspaceID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationSpendLimitSetParamsScopeUnion) GetType() *string {
	if vt := u.OfUser; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfOrganization; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWorkspace; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationSpendLimitSetParamsScopeUnion](
		"type",
		apijson.Discriminator[BetaSpendLimitUserScopeParam]("user"),
		apijson.Discriminator[BetaSpendLimitOrganizationScopeParam]("organization"),
		apijson.Discriminator[BetaSpendLimitWorkspaceScopeParam]("workspace"),
	)
}

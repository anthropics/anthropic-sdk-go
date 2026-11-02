package anthropic

import (
	"encoding/json"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationSpendLimitService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationSpendLimitService] method instead.
type OrganizationSpendLimitService struct {
	Options          []option.RequestOption
	Effective        OrganizationSpendLimitEffectiveService
	IncreaseRequests OrganizationSpendLimitIncreaseRequestService
}

// NewOrganizationSpendLimitService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationSpendLimitService(opts ...option.RequestOption) (r OrganizationSpendLimitService) {
	r = OrganizationSpendLimitService{}
	r.Options = opts
	r.Effective = NewOrganizationSpendLimitEffectiveService(opts...)
	r.IncreaseRequests = NewOrganizationSpendLimitIncreaseRequestService(opts...)
	return
}

// A configured spend limit: a cap on metered spend for one scope and period.
type SpendLimit struct {
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
	// (`organization` limit); `amount` is kept and applies again when it's switched
	// back on. Always `true` for other limits.
	IsEnabled bool `json:"is_enabled" api:"required"`
	// Length of the window the limit resets over. `amount` caps spend within each
	// period.
	//
	// Any of "daily", "monthly", "weekly".
	Period SpendLimitPeriod `json:"period" api:"required"`
	// What the limit applies to. A tagged union on `type`; each variant carries the
	// identifier for its scope.
	Scope SpendLimitScopeUnion `json:"scope" api:"required"`
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
func (r SpendLimit) RawJSON() string { return r.JSON.raw }
func (r *SpendLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendLimitScopeUnion contains all possible properties and values from
// [SpendLimitUserScope], [SpendLimitSeatTierScope], [SpendLimitRBACGroupScope],
// [SpendLimitOrganizationServiceScope], [SpendLimitOrganizationScope],
// [SpendLimitWorkspaceScope], [SpendLimitUsageTierScope].
//
// Use the [SpendLimitScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendLimitScopeUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace", "usage_tier".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [SpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [SpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [SpendLimitWorkspaceScope].
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

// anySpendLimitScope is implemented by each variant of [SpendLimitScopeUnion] to
// add type safety for the return type of [SpendLimitScopeUnion.AsAny]
type anySpendLimitScope interface {
	implSpendLimitScopeUnion()
}

func (SpendLimitUserScope) implSpendLimitScopeUnion()                {}
func (SpendLimitSeatTierScope) implSpendLimitScopeUnion()            {}
func (SpendLimitRBACGroupScope) implSpendLimitScopeUnion()           {}
func (SpendLimitOrganizationServiceScope) implSpendLimitScopeUnion() {}
func (SpendLimitOrganizationScope) implSpendLimitScopeUnion()        {}
func (SpendLimitWorkspaceScope) implSpendLimitScopeUnion()           {}
func (SpendLimitUsageTierScope) implSpendLimitScopeUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendLimitScopeUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserScope:
//	case anthropic.SpendLimitSeatTierScope:
//	case anthropic.SpendLimitRBACGroupScope:
//	case anthropic.SpendLimitOrganizationServiceScope:
//	case anthropic.SpendLimitOrganizationScope:
//	case anthropic.SpendLimitWorkspaceScope:
//	case anthropic.SpendLimitUsageTierScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendLimitScopeUnion) AsAny() anySpendLimitScope {
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
	case "usage_tier":
		return u.AsUsageTier()
	}
	return nil
}

func (u SpendLimitScopeUnion) AsUser() (v SpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsSeatTier() (v SpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsRBACGroup() (v SpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsOrganizationService() (v SpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsOrganization() (v SpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsWorkspace() (v SpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitScopeUnion) AsUsageTier() (v SpendLimitUsageTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendLimitScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendLimitScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SpendLimitOrganizationScope struct {
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SpendLimitOrganizationScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitOrganizationScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SpendLimitOrganizationServiceScope struct {
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
func (r SpendLimitOrganizationServiceScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitOrganizationServiceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SpendLimitPeriod string

const (
	SpendLimitPeriodDaily   SpendLimitPeriod = "daily"
	SpendLimitPeriodMonthly SpendLimitPeriod = "monthly"
	SpendLimitPeriodWeekly  SpendLimitPeriod = "weekly"
)

type SpendLimitRBACGroupScope struct {
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
func (r SpendLimitRBACGroupScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitRBACGroupScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A scoped Admin API key acting on behalf of the organization.
type SpendLimitScopedAPIKeyActor struct {
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
func (r SpendLimitScopedAPIKeyActor) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitScopedAPIKeyActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SpendLimitSeatTierScope struct {
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
func (r SpendLimitSeatTierScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitSeatTierScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Scope of the monthly spend cap Anthropic sets for a Claude Console organization
// from its usage tier, or a higher cap Anthropic granted the organization.
// Read-only; its `created_at` and `updated_at` are the organization's creation
// time, so `amount` can change while `updated_at` stays the same.
type SpendLimitUsageTierScope struct {
	// Scope type. Always `usage_tier` for this scope.
	Type constant.UsageTier `json:"type" default:"usage_tier"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SpendLimitUsageTierScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitUsageTierScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user within the organization. `name` and `email_address` are null when the
// underlying account is unavailable or has been deleted; `deleted` is true only
// for deleted accounts.
type SpendLimitUserActor struct {
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
func (r SpendLimitUserActor) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Scope selecting a single member of the organization.
type SpendLimitUserScope struct {
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
func (r SpendLimitUserScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitUserScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Scope selecting one workspace of a Claude Console organization.
type SpendLimitWorkspaceScope struct {
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
func (r SpendLimitWorkspaceScope) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitWorkspaceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-member effective-limit report row (`GET /spend_limits/effective`).
type SpendSummary struct {
	Actor SpendSummaryActorUnion `json:"actor" api:"required"`
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
	Period SpendLimitPeriod `json:"period" api:"required"`
	// The member's spend so far in the current period, as a non-negative decimal
	// string in the minor unit of `currency` (cents for USD). May carry fractional
	// minor units up to three decimal places (e.g. `"12050.5"`) — metered usage is not
	// rounded to whole cents. Reads as `"0"` when the spend reading is temporarily
	// unavailable.
	PeriodToDateSpend string                  `json:"period_to_date_spend" api:"required"`
	Scope             SpendSummaryScopeUnion  `json:"scope" api:"required"`
	Source            SpendSummarySourceUnion `json:"source" api:"required"`
	SpendLimitID      string                  `json:"spend_limit_id" api:"required"`
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
func (r SpendSummary) RawJSON() string { return r.JSON.raw }
func (r *SpendSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendSummaryActorUnion contains all possible properties and values from
// [SpendLimitUserActor], [SpendLimitScopedAPIKeyActor].
//
// Use the [SpendSummaryActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendSummaryActorUnion struct {
	// This field is from variant [SpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [SpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [SpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitScopedAPIKeyActor].
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

// anySpendSummaryActor is implemented by each variant of [SpendSummaryActorUnion]
// to add type safety for the return type of [SpendSummaryActorUnion.AsAny]
type anySpendSummaryActor interface {
	implSpendSummaryActorUnion()
}

func (SpendLimitUserActor) implSpendSummaryActorUnion()         {}
func (SpendLimitScopedAPIKeyActor) implSpendSummaryActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendSummaryActorUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserActor:
//	case anthropic.SpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendSummaryActorUnion) AsAny() anySpendSummaryActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u SpendSummaryActorUnion) AsUserActor() (v SpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryActorUnion) AsScopedAPIKeyActor() (v SpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendSummaryActorUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendSummaryActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendSummaryScopeUnion contains all possible properties and values from
// [SpendLimitUserScope], [SpendLimitSeatTierScope], [SpendLimitRBACGroupScope],
// [SpendLimitOrganizationServiceScope], [SpendLimitOrganizationScope],
// [SpendLimitWorkspaceScope], [SpendLimitUsageTierScope].
//
// Use the [SpendSummaryScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendSummaryScopeUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace", "usage_tier".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [SpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [SpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [SpendLimitWorkspaceScope].
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

// anySpendSummaryScope is implemented by each variant of [SpendSummaryScopeUnion]
// to add type safety for the return type of [SpendSummaryScopeUnion.AsAny]
type anySpendSummaryScope interface {
	implSpendSummaryScopeUnion()
}

func (SpendLimitUserScope) implSpendSummaryScopeUnion()                {}
func (SpendLimitSeatTierScope) implSpendSummaryScopeUnion()            {}
func (SpendLimitRBACGroupScope) implSpendSummaryScopeUnion()           {}
func (SpendLimitOrganizationServiceScope) implSpendSummaryScopeUnion() {}
func (SpendLimitOrganizationScope) implSpendSummaryScopeUnion()        {}
func (SpendLimitWorkspaceScope) implSpendSummaryScopeUnion()           {}
func (SpendLimitUsageTierScope) implSpendSummaryScopeUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendSummaryScopeUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserScope:
//	case anthropic.SpendLimitSeatTierScope:
//	case anthropic.SpendLimitRBACGroupScope:
//	case anthropic.SpendLimitOrganizationServiceScope:
//	case anthropic.SpendLimitOrganizationScope:
//	case anthropic.SpendLimitWorkspaceScope:
//	case anthropic.SpendLimitUsageTierScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendSummaryScopeUnion) AsAny() anySpendSummaryScope {
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
	case "usage_tier":
		return u.AsUsageTier()
	}
	return nil
}

func (u SpendSummaryScopeUnion) AsUser() (v SpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsSeatTier() (v SpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsRBACGroup() (v SpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsOrganizationService() (v SpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsOrganization() (v SpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsWorkspace() (v SpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummaryScopeUnion) AsUsageTier() (v SpendLimitUsageTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendSummaryScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendSummaryScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendSummarySourceUnion contains all possible properties and values from
// [SpendLimitUserScope], [SpendLimitSeatTierScope], [SpendLimitRBACGroupScope],
// [SpendLimitOrganizationServiceScope], [SpendLimitOrganizationScope],
// [SpendLimitWorkspaceScope], [SpendLimitUsageTierScope].
//
// Use the [SpendSummarySourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendSummarySourceUnion struct {
	// Any of "user", "seat_tier", "rbac_group", "organization_service",
	// "organization", "workspace", "usage_tier".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserScope].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitSeatTierScope].
	SeatTier string `json:"seat_tier"`
	// This field is from variant [SpendLimitRBACGroupScope].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [SpendLimitOrganizationServiceScope].
	Service string `json:"service"`
	// This field is from variant [SpendLimitWorkspaceScope].
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

// anySpendSummarySource is implemented by each variant of
// [SpendSummarySourceUnion] to add type safety for the return type of
// [SpendSummarySourceUnion.AsAny]
type anySpendSummarySource interface {
	implSpendSummarySourceUnion()
}

func (SpendLimitUserScope) implSpendSummarySourceUnion()                {}
func (SpendLimitSeatTierScope) implSpendSummarySourceUnion()            {}
func (SpendLimitRBACGroupScope) implSpendSummarySourceUnion()           {}
func (SpendLimitOrganizationServiceScope) implSpendSummarySourceUnion() {}
func (SpendLimitOrganizationScope) implSpendSummarySourceUnion()        {}
func (SpendLimitWorkspaceScope) implSpendSummarySourceUnion()           {}
func (SpendLimitUsageTierScope) implSpendSummarySourceUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendSummarySourceUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserScope:
//	case anthropic.SpendLimitSeatTierScope:
//	case anthropic.SpendLimitRBACGroupScope:
//	case anthropic.SpendLimitOrganizationServiceScope:
//	case anthropic.SpendLimitOrganizationScope:
//	case anthropic.SpendLimitWorkspaceScope:
//	case anthropic.SpendLimitUsageTierScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendSummarySourceUnion) AsAny() anySpendSummarySource {
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
	case "usage_tier":
		return u.AsUsageTier()
	}
	return nil
}

func (u SpendSummarySourceUnion) AsUser() (v SpendLimitUserScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsSeatTier() (v SpendLimitSeatTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsRBACGroup() (v SpendLimitRBACGroupScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsOrganizationService() (v SpendLimitOrganizationServiceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsOrganization() (v SpendLimitOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsWorkspace() (v SpendLimitWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendSummarySourceUnion) AsUsageTier() (v SpendLimitUsageTierScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendSummarySourceUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendSummarySourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

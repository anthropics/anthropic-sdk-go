package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
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

// UserProfileService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewUserProfileService] method instead.
type UserProfileService struct {
	Options    []option.RequestOption
	ActionURLs UserProfileActionURLService
}

// NewUserProfileService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewUserProfileService(opts ...option.RequestOption) (r UserProfileService) {
	r = UserProfileService{}
	r.Options = opts
	r.ActionURLs = NewUserProfileActionURLService(opts...)
	return
}

// Create User Profile
func (r *UserProfileService) New(ctx context.Context, params UserProfileNewParams, opts ...option.RequestOption) (res *UserProfile, err error) {
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	path := "v1/user_profiles"
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.create", "/v1/user_profiles"), requestconfig.WithSpanCreatedID("anthropic.user_profile.id", &res, func(created *UserProfile) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get User Profile
func (r *UserProfileService) Get(ctx context.Context, userProfileID string, query UserProfileGetParams, opts ...option.RequestOption) (res *UserProfile, err error) {
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s", url.PathEscape(userProfileID))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.retrieve", "/v1/user_profiles/{user_profile_id}"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update User Profile
func (r *UserProfileService) Update(ctx context.Context, userProfileID string, params UserProfileUpdateParams, opts ...option.RequestOption) (res *UserProfile, err error) {
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s", url.PathEscape(userProfileID))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.update", "/v1/user_profiles/{user_profile_id}"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List User Profiles
func (r *UserProfileService) List(ctx context.Context, params UserProfileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[UserProfile], err error) {
	var raw *http.Response
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/user_profiles"
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.list", "/v1/user_profiles", "enforcement_status", "limit", "order", "order_by", "page"))
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

// List User Profiles
func (r *UserProfileService) ListAutoPaging(ctx context.Context, params UserProfileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[UserProfile] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete User Profile
func (r *UserProfileService) Delete(ctx context.Context, userProfileID string, body UserProfileDeleteParams, opts ...option.RequestOption) (res *DeletedUserProfile, err error) {
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s", url.PathEscape(userProfileID))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.delete", "/v1/user_profiles/{user_profile_id}"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedUserProfile struct {
	// Identifier of the deleted user profile, prefixed `uprof_`.
	ID string `json:"id" api:"required"`
	// Object type. Always `user_profile_deleted`.
	Type constant.UserProfileDeleted `json:"type" default:"user_profile_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DeletedUserProfile) RawJSON() string { return r.JSON.raw }
func (r *DeletedUserProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A record of an entity that the platform serves through the API, such as an
// end-user of the platform's product or a company that the platform resells Claude
// access to.
//
// A Messages, Message Batches or token counting request can send a profile's `id`
// in the `anthropic-user-profile-id` header to attribute the request to that
// entity.
type UserProfile struct {
	// Unique identifier for this user profile, prefixed `uprof_`.
	ID string `json:"id" api:"required"`
	// How the platform uses the API for this entity: `application` (default) or
	// `passthrough`. Absent under the legacy `user-profiles-2026-03-24` beta header.
	//
	// Any of "application", "passthrough".
	AccessType UserProfileAccessType `json:"access_type" api:"required"`
	// When this user profile was created, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The enforcement actions open on this profile; empty when none is open. Absent
	// under the `user-profiles-2026-09-04` and earlier beta headers.
	EnforcementActions []UserProfileEnforcementActionUnion `json:"enforcement_actions" api:"required"`
	// Anthropic's enforcement status of this profile, from its open enforcement
	// actions and any identity check under review, by the first rule that applies:
	// `banned` while any open action is a `ban`; `pending` while a restriction that an
	// identity check can clear is open, the end user has submitted that check,
	// Anthropic's review of it has not ended, and no open action's `reason` is
	// `account_compromise`; `restricted` while any action is open; `unrestricted`
	// while none is. Absent under the `user-profiles-2026-09-04` and earlier beta
	// headers.
	//
	// Any of "unrestricted", "pending", "restricted", "banned".
	EnforcementStatus UserProfileEnforcementStatus `json:"enforcement_status" api:"required"`
	// Details about the entity this profile represents, as the platform states them;
	// not verified by Anthropic. Every field is present, `null` until the platform
	// supplies a value. Absent under the legacy `user-profiles-2026-03-24` and
	// `user-profiles-2026-08-18` beta headers, which serve `reference_id` as the
	// top-level `external_id` and, under `user-profiles-2026-08-18`, `onboarded_at` as
	// `external_user_onboarded_at`.
	ExternalUserDetails UserProfileExternalUserDetails `json:"external_user_details" api:"required"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Real-world name of the entity this profile represents (company or individual).
	// For a company the platform resells Claude access to (`access_type`
	// `passthrough`) this is that company's name. `null` until the platform supplies
	// one.
	Name string `json:"name" api:"required"`
	// The profile's standing in each Anthropic program, such as the Cyber Verification
	// Program. Present when the request sends no user-profiles beta header; absent
	// under the `user-profiles-2026-09-04` and earlier beta headers, which serve
	// `trust_grants` instead.
	Programs UserProfilePrograms `json:"programs" api:"required"`
	// Object type. Always `user_profile`.
	Type constant.UserProfile `json:"type" default:"user_profile"`
	// When this user profile was last modified, in RFC 3339 format. It normally also
	// updates shortly after Anthropic gives the profile an approval or withdraws one.
	// Other changes to the statuses that Anthropic reports do not always update it, so
	// read the status itself rather than watching this timestamp.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		AccessType          respjson.Field
		CreatedAt           respjson.Field
		EnforcementActions  respjson.Field
		EnforcementStatus   respjson.Field
		ExternalUserDetails respjson.Field
		Metadata            respjson.Field
		Name                respjson.Field
		Programs            respjson.Field
		Type                respjson.Field
		UpdatedAt           respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfile) RawJSON() string { return r.JSON.raw }
func (r *UserProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How the platform uses the API on behalf of the entity this profile represents.
// `application`: the platform sells a product that uses the API behind the scenes,
// and the profile represents an individual end-user of that product.
// `passthrough`: the platform resells raw inference, and the profile identifies
// the resold-to company.
type UserProfileAccessType string

const (
	// The user profile represents an individual end-user of a product that the
	// platform builds on the API. New profiles get this value by default.
	UserProfileAccessTypeApplication UserProfileAccessType = "application"
	// The user profile represents a company that the platform resells Claude access
	// to.
	UserProfileAccessTypePassthrough UserProfileAccessType = "passthrough"
)

// A user profile's standing in the Cyber Verification Program, which relaxes cyber
// safeguards for requests that name the profile.
type UserProfileCyberProgram struct {
	// What a request that names the profile must allow to get the relaxed safeguards
	// of every `active` tier. `null` when no tier is `active`, and wherever `tiers` is
	// `{}`.
	Requirements UserProfileCyberRequirements `json:"requirements" api:"required"`
	// `active` when at least one tier is `active`, or when the profile has the
	// program's access to earlier models. That access is not a tier: a profile that
	// has it alone reads `active` with every tier `inactive`, or with `tiers` `{}` on
	// a response that does not report tiers. When nothing is `active`, `pending` means
	// that an application or enrollment is under review, or that what Anthropic
	// approved does not yet apply to this profile, and `rejected` means that the
	// application or enrollment was declined.
	//
	// Any of "active", "pending", "rejected".
	Status UserProfileProgramStatus `json:"status" api:"required"`
	// The program's tiers, keyed by tier: `a` (Specialized Access), `b` (Red Team
	// Access) and `c` (Defense Access). A response reports either every tier or none.
	// When it reports them, as retrieving the profile normally does, every tier has a
	// key, whether or not it applies to the profile. When it does not, as in the list
	// of user profiles, `tiers` is `{}`, which does not mean that the tiers are
	// inactive. No tier implies another, so read the tier you need. Anthropic may add
	// tiers; ignore a key you do not recognize.
	Tiers map[string]UserProfileCyberTier `json:"tiers" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Requirements respjson.Field
		Status       respjson.Field
		Tiers        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileCyberProgram) RawJSON() string { return r.JSON.raw }
func (r *UserProfileCyberProgram) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What a request that names the profile must allow to get the relaxed safeguards
// of the profile's `active` tiers. A field reads `true` when at least one of those
// tiers has the requirement and Anthropic has not waived it on that tier for the
// profile. A request that does not meet a requirement is still served, without the
// relaxed safeguards of the tiers that have it. The values differ between profiles
// and can change, so read them rather than assuming them.
type UserProfileCyberRequirements struct {
	// Whether Anthropic must be able to retain the request's inference data.
	DataRetention bool `json:"data_retention" api:"required"`
	// Whether Anthropic personnel must be allowed to review the data retained from the
	// request.
	HumanReview bool `json:"human_review" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DataRetention respjson.Field
		HumanReview   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileCyberRequirements) RawJSON() string { return r.JSON.raw }
func (r *UserProfileCyberRequirements) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One tier of the Cyber Verification Program.
type UserProfileCyberTier struct {
	// `active` when the tier applies to the profile, otherwise `inactive`. The
	// response does not say why a tier is `inactive`: it can be in review, declined,
	// withdrawn, never applied for, or approved but not applied to this profile.
	//
	// Any of "active", "inactive".
	Status UserProfileCyberTierStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileCyberTier) RawJSON() string { return r.JSON.raw }
func (r *UserProfileCyberTier) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `active` when the tier applies to the profile, otherwise `inactive`. The
// response does not say why a tier is `inactive`: it can be in review, declined,
// withdrawn, never applied for, or approved but not applied to this profile.
type UserProfileCyberTierStatus string

const (
	UserProfileCyberTierStatusActive   UserProfileCyberTierStatus = "active"
	UserProfileCyberTierStatusInactive UserProfileCyberTierStatus = "inactive"
)

// One feature of Enterprise Frontier Safeguards.
type UserProfileEfsFeature struct {
	// `active` when Anthropic has approved the entity the profile represents to use
	// the feature, otherwise `inactive`. `active` does not mean that the feature is
	// set up or switched on, which is a separate step.
	//
	// Any of "active", "inactive".
	Status UserProfileEfsFeatureStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileEfsFeature) RawJSON() string { return r.JSON.raw }
func (r *UserProfileEfsFeature) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `active` when Anthropic has approved the entity the profile represents to use
// the feature, otherwise `inactive`. `active` does not mean that the feature is
// set up or switched on, which is a separate step.
type UserProfileEfsFeatureStatus string

const (
	UserProfileEfsFeatureStatusActive   UserProfileEfsFeatureStatus = "active"
	UserProfileEfsFeatureStatusInactive UserProfileEfsFeatureStatus = "inactive"
)

// A user profile's standing in Enterprise Frontier Safeguards, a program under
// which Anthropic approves the entity the profile represents and may approve any
// of the program's features.
type UserProfileEfsProgram struct {
	// The program's features, keyed by feature: `data_retention_exemption` (use of
	// models that normally require data retention, without the data being retained),
	// `cmek` (customer-managed encryption keys) and `storage` (customer-owned
	// storage). Every feature that the caller's organization may read has a key,
	// whether or not it is approved; a missing key does not mean that the feature is
	// inactive. Anthropic may add features; ignore a key you do not recognize.
	Features map[string]UserProfileEfsFeature `json:"features" api:"required"`
	// Anthropic's decision on the entity the profile represents: `pending` while it is
	// awaited, `active` once Anthropic has approved the entity for the program, or
	// `rejected` once Anthropic has declined or has removed the entity from the
	// program. `active` does not mean that any one feature is approved: read the
	// feature you need in `features`.
	//
	// Any of "active", "pending", "rejected".
	Status UserProfileProgramStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Features    respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileEfsProgram) RawJSON() string { return r.JSON.raw }
func (r *UserProfileEfsProgram) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UserProfileEnforcementActionUnion contains all possible properties and values
// from [UserProfileEnforcementActionRestriction],
// [UserProfileEnforcementActionBan].
//
// Use the [UserProfileEnforcementActionUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type UserProfileEnforcementActionUnion struct {
	ID string `json:"id"`
	// This field is from variant [UserProfileEnforcementActionRestriction].
	Action    UserProfileEnforcementActionNextStep `json:"action"`
	ExpiresAt time.Time                            `json:"expires_at"`
	OpenedAt  time.Time                            `json:"opened_at"`
	// This field is from variant [UserProfileEnforcementActionRestriction].
	Reason UserProfileEnforcementActionReason `json:"reason"`
	// Any of "restriction", "ban".
	Type string `json:"type"`
	JSON struct {
		ID        respjson.Field
		Action    respjson.Field
		ExpiresAt respjson.Field
		OpenedAt  respjson.Field
		Reason    respjson.Field
		Type      respjson.Field
		raw       string
	} `json:"-"`
}

// anyUserProfileEnforcementAction is implemented by each variant of
// [UserProfileEnforcementActionUnion] to add type safety for the return type of
// [UserProfileEnforcementActionUnion.AsAny]
type anyUserProfileEnforcementAction interface {
	implUserProfileEnforcementActionUnion()
}

func (UserProfileEnforcementActionRestriction) implUserProfileEnforcementActionUnion() {}
func (UserProfileEnforcementActionBan) implUserProfileEnforcementActionUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := UserProfileEnforcementActionUnion.AsAny().(type) {
//	case anthropic.UserProfileEnforcementActionRestriction:
//	case anthropic.UserProfileEnforcementActionBan:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u UserProfileEnforcementActionUnion) AsAny() anyUserProfileEnforcementAction {
	switch u.Type {
	case "restriction":
		return u.AsRestriction()
	case "ban":
		return u.AsBan()
	}
	return nil
}

func (u UserProfileEnforcementActionUnion) AsRestriction() (v UserProfileEnforcementActionRestriction) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u UserProfileEnforcementActionUnion) AsBan() (v UserProfileEnforcementActionBan) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u UserProfileEnforcementActionUnion) RawJSON() string { return u.JSON.raw }

func (r *UserProfileEnforcementActionUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic has stopped serving the end user where the ban applies and does not
// plan to serve them there again.
type UserProfileEnforcementActionBan struct {
	// The action's ID (`enac_...`), stable while the action is open.
	ID string `json:"id" api:"required"`
	// Always `null`: the end user cannot clear a ban.
	//
	// Any of "verify".
	Action UserProfileEnforcementActionNextStep `json:"action" api:"required"`
	// Always `null`: a ban has no end date.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// When Anthropic opened the action, in RFC 3339 format.
	OpenedAt time.Time `json:"opened_at" api:"required" format:"date-time"`
	// Why Anthropic placed the limit: `unsupported_region`, `fraud_or_abuse`,
	// `high_harm`, `usage_policy_violation`, `account_compromise` or `unspecified`.
	//
	// Any of "unsupported_region", "fraud_or_abuse", "high_harm",
	// "usage_policy_violation", "account_compromise", "unspecified".
	Reason UserProfileEnforcementActionReason `json:"reason" api:"required"`
	Type   constant.Ban                       `json:"type" default:"ban"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Action      respjson.Field
		ExpiresAt   respjson.Field
		OpenedAt    respjson.Field
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileEnforcementActionBan) RawJSON() string { return r.JSON.raw }
func (r *UserProfileEnforcementActionBan) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The next step on this enforcement action: `verify` means the end user can clear
// the action by passing an identity check. More values may be added; read a value
// you do not recognize as `null`.
type UserProfileEnforcementActionNextStep string

const (
	UserProfileEnforcementActionNextStepVerify UserProfileEnforcementActionNextStep = "verify"
)

// Why Anthropic placed this enforcement action: `unsupported_region`,
// `fraud_or_abuse`, `high_harm`, `usage_policy_violation`, `account_compromise` or
// `unspecified`. More values may be added; read a value you do not recognize as
// `unspecified`.
type UserProfileEnforcementActionReason string

const (
	UserProfileEnforcementActionReasonUnsupportedRegion    UserProfileEnforcementActionReason = "unsupported_region"
	UserProfileEnforcementActionReasonFraudOrAbuse         UserProfileEnforcementActionReason = "fraud_or_abuse"
	UserProfileEnforcementActionReasonHighHarm             UserProfileEnforcementActionReason = "high_harm"
	UserProfileEnforcementActionReasonUsagePolicyViolation UserProfileEnforcementActionReason = "usage_policy_violation"
	UserProfileEnforcementActionReasonAccountCompromise    UserProfileEnforcementActionReason = "account_compromise"
	UserProfileEnforcementActionReasonUnspecified          UserProfileEnforcementActionReason = "unspecified"
)

// Anthropic has limited how it serves the end user until the action is resolved.
// The end user may still be served in part.
type UserProfileEnforcementActionRestriction struct {
	// The action's ID (`enac_...`), stable while the action is open.
	ID string `json:"id" api:"required"`
	// The next step on the action: `verify` or `null`. `verify`: the end user can
	// clear it now by passing an identity check. `null`: there is nothing the end user
	// can do now to clear it, including while Anthropic reviews an identity check the
	// end user submitted.
	//
	// Any of "verify".
	Action UserProfileEnforcementActionNextStep `json:"action" api:"required"`
	// When the action ends on its own, in RFC 3339 format; it stops showing once this
	// time has passed. `null` when it has no end date, and always `null` when `reason`
	// is `account_compromise`.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// When Anthropic opened the action, in RFC 3339 format.
	OpenedAt time.Time `json:"opened_at" api:"required" format:"date-time"`
	// Why Anthropic placed the limit: `unsupported_region`, `fraud_or_abuse`,
	// `high_harm`, `usage_policy_violation`, `account_compromise` or `unspecified`.
	//
	// Any of "unsupported_region", "fraud_or_abuse", "high_harm",
	// "usage_policy_violation", "account_compromise", "unspecified".
	Reason UserProfileEnforcementActionReason `json:"reason" api:"required"`
	Type   constant.Restriction               `json:"type" default:"restriction"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Action      respjson.Field
		ExpiresAt   respjson.Field
		OpenedAt    respjson.Field
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileEnforcementActionRestriction) RawJSON() string { return r.JSON.raw }
func (r *UserProfileEnforcementActionRestriction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic's enforcement status of this profile: `unrestricted`, `pending`,
// `restricted` or `banned`. More values may be added; read a value you do not
// recognize as `restricted`.
type UserProfileEnforcementStatus string

const (
	UserProfileEnforcementStatusUnrestricted UserProfileEnforcementStatus = "unrestricted"
	UserProfileEnforcementStatusPending      UserProfileEnforcementStatus = "pending"
	UserProfileEnforcementStatusRestricted   UserProfileEnforcementStatus = "restricted"
	UserProfileEnforcementStatusBanned       UserProfileEnforcementStatus = "banned"
)

// The status of the entity's account on the platform, as the platform states it:
// `active`; `suspended`, when the platform has restricted the account and may
// restore it; or `blocked`, when the platform has barred it. It records the
// platform's decision only. The statuses that Anthropic itself reports for the
// profile, such as its status in a program, are Anthropic's decisions and do not
// change with it.
type UserProfileExternalUserAccountStatus string

const (
	// The platform has neither restricted nor barred the account of the entity that
	// the user profile represents.
	UserProfileExternalUserAccountStatusActive UserProfileExternalUserAccountStatus = "active"
	// The platform has restricted the account of the entity that the user profile
	// represents and may restore it.
	UserProfileExternalUserAccountStatusSuspended UserProfileExternalUserAccountStatus = "suspended"
	// The platform has barred the account of the entity that the user profile
	// represents.
	UserProfileExternalUserAccountStatusBlocked UserProfileExternalUserAccountStatus = "blocked"
)

// Details about the entity this profile represents, as the platform states them.
// Anthropic does not verify them. Every field is present, `null` until the
// platform supplies a value.
type UserProfileExternalUserDetails struct {
	// The status of the entity's account on the platform: `active`, `suspended` or
	// `blocked`. `null` until the platform supplies one.
	//
	// Any of "active", "suspended", "blocked".
	AccountStatus UserProfileExternalUserAccountStatus `json:"account_status" api:"required"`
	// The country the platform associates with the entity, as an ISO 3166-1 alpha-2
	// code. `null` until the platform supplies one.
	Country string `json:"country" api:"required"`
	// The platform-computed hash of the entity's email address. `null` until the
	// platform supplies one.
	EmailHash string `json:"email_hash" api:"required"`
	// What kind of entity the profile represents: `individual`, `business`,
	// `non_profit` or `government`. `null` until the platform supplies one.
	//
	// Any of "individual", "business", "non_profit", "government".
	EntityType UserProfileExternalUserEntityType `json:"entity_type" api:"required"`
	// The platform-computed hash of the entity's name. `null` until the platform
	// supplies one.
	NameHash string `json:"name_hash" api:"required"`
	// When the entity opened its account with the platform, as stated by the platform,
	// in RFC 3339 format (UTC). `null` until the platform supplies one.
	OnboardedAt time.Time `json:"onboarded_at" api:"required" format:"date-time"`
	// The platform's own reference for the entity. `null` until the platform supplies
	// one.
	ReferenceID string `json:"reference_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccountStatus respjson.Field
		Country       respjson.Field
		EmailHash     respjson.Field
		EntityType    respjson.Field
		NameHash      respjson.Field
		OnboardedAt   respjson.Field
		ReferenceID   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileExternalUserDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileExternalUserDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UserProfileExternalUserDetailsParams struct {
	// The country of the entity (not of the platform), as the platform determines it:
	// an ISO 3166-1 alpha-2 code in upper case, for example `US`. Only the form, two
	// uppercase ASCII letters, is checked.
	Country param.Opt[string] `json:"country,omitzero"`
	// A hash of the entity's email address, computed by the platform. Anthropic treats
	// it as an opaque string and does not prescribe the hash function. 1 to 255
	// characters.
	EmailHash param.Opt[string] `json:"email_hash,omitzero"`
	// A hash of the entity's name, computed by the platform. Anthropic treats it as an
	// opaque string and does not prescribe the hash function. 1 to 255 characters.
	NameHash param.Opt[string] `json:"name_hash,omitzero"`
	// The platform's own reference for the entity, for example the key of the
	// end-user's row in the platform's database. Not interpreted by Anthropic and not
	// enforced unique. 1 to 255 characters.
	ReferenceID param.Opt[string] `json:"reference_id,omitzero"`
	// When the entity opened its account with the platform, in RFC 3339 format: for an
	// `application` profile, when the end-user signed up; for a `passthrough` profile,
	// when the company became the platform's customer. Must be a complete timestamp no
	// more than 1 minute in the future.
	OnboardedAt param.Opt[time.Time] `json:"onboarded_at,omitzero" format:"date-time"`
	// The status of the entity's account on the platform: `active`, `suspended` or
	// `blocked`.
	//
	// Any of "active", "suspended", "blocked".
	AccountStatus UserProfileExternalUserAccountStatus `json:"account_status,omitzero"`
	// What kind of entity the profile represents: `individual`, `business`,
	// `non_profit` or `government`.
	//
	// Any of "individual", "business", "non_profit", "government".
	EntityType UserProfileExternalUserEntityType `json:"entity_type,omitzero"`
	paramObj
}

func (r UserProfileExternalUserDetailsParams) MarshalJSON() (data []byte, err error) {
	type shadow UserProfileExternalUserDetailsParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UserProfileExternalUserDetailsParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What kind of entity the profile represents, as the platform states it:
// `individual`, `business`, `non_profit` or `government`.
type UserProfileExternalUserEntityType string

const (
	UserProfileExternalUserEntityTypeIndividual UserProfileExternalUserEntityType = "individual"
	UserProfileExternalUserEntityTypeBusiness   UserProfileExternalUserEntityType = "business"
	UserProfileExternalUserEntityTypeNonProfit  UserProfileExternalUserEntityType = "non_profit"
	UserProfileExternalUserEntityTypeGovernment UserProfileExternalUserEntityType = "government"
)

// The direction that a list of user profiles is sorted in.
type UserProfileListOrder string

const (
	// Oldest first when `order_by` is `created_at`, or names in ascending order when
	// `order_by` is `name`.
	UserProfileListOrderAsc UserProfileListOrder = "asc"
	// Newest first when `order_by` is `created_at`, or names in descending order when
	// `order_by` is `name`. This is the default.
	UserProfileListOrderDesc UserProfileListOrder = "desc"
)

// Sort field for listing user profiles: `created_at` (default) or `name`
// (case-insensitive; profiles without a name sort last).
type UserProfileListOrderBy string

const (
	// Sort by when each user profile was created. This is the default.
	UserProfileListOrderByCreatedAt UserProfileListOrderBy = "created_at"
	// Sort by `name`, ignoring the case of ASCII letters. Profiles without a name come
	// last in either direction.
	UserProfileListOrderByName UserProfileListOrderBy = "name"
)

// A user profile's status in a program. Each program's `status` field says what
// the values mean there.
type UserProfileProgramStatus string

const (
	// The profile is in the program.
	UserProfileProgramStatusActive UserProfileProgramStatus = "active"
	// The profile is not in the program yet, and something is still awaited.
	UserProfileProgramStatusPending UserProfileProgramStatus = "pending"
	// The profile is not in the program and nothing is awaited: Anthropic has
	// declined, or has withdrawn an earlier approval.
	UserProfileProgramStatusRejected UserProfileProgramStatus = "rejected"
)

// A user profile's standing in each program. A program is a set of capabilities or
// exemptions that Anthropic approves for the entity the profile represents. A
// program reads `null` when the profile has no status in it, when the caller's
// organization is not entitled to read programs, or when the response does not
// report that program. Retrieving the profile reports every program. Anthropic may
// add programs; ignore a key you do not recognize.
type UserProfilePrograms struct {
	// The profile's standing in the Cyber Verification Program. `null` when the
	// profile has no status in it. Only retrieving the profile gives the full answer:
	// a list or update response can show `null`, or a status other than `active`, for
	// a profile that reads `active` when retrieved.
	Cyber UserProfileCyberProgram `json:"cyber" api:"required"`
	// The profile's standing in Enterprise Frontier Safeguards. `null` unless the
	// platform's attestation for the program is in place on this profile, or the
	// entity it represents has applied to the program. A list of user profiles does
	// not report this program: there it reads `null` for every profile, so retrieve
	// the profile to read it.
	Efs UserProfileEfsProgram `json:"efs" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Cyber       respjson.Field
		Efs         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfilePrograms) RawJSON() string { return r.JSON.raw }
func (r *UserProfilePrograms) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UserProfileNewParams struct {
	// Optional for all profiles. Real-world name of the entity this profile represents
	// (company or individual); for a company the platform resells Claude access to
	// (`access_type` `passthrough`), that company's name where known. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// How the platform uses the API for this entity. `application` (default): the
	// profile represents an individual end-user of the platform's product.
	// `passthrough`: the profile identifies a company the platform resells Claude
	// access to.
	//
	// Any of "application", "passthrough".
	AccessType UserProfileAccessType `json:"access_type,omitzero"`
	// Details about the entity this profile represents, as the platform states them.
	// Every field is optional. Not accepted under the legacy
	// `user-profiles-2026-03-24` and `user-profiles-2026-08-18` beta headers, which
	// take the top-level `external_id` (and, under `user-profiles-2026-08-18`,
	// `external_user_onboarded_at`) instead.
	ExternalUserDetails UserProfileExternalUserDetailsParams `json:"external_user_details,omitzero"`
	// Free-form key-value data to attach to this user profile. Maximum 16 keys, with
	// keys up to 64 characters and values up to 512 characters. Values must be
	// non-empty strings.
	Metadata map[string]string `json:"metadata,omitzero"`
	paramObj
}

func (r UserProfileNewParams) MarshalJSON() (data []byte, err error) {
	type shadow UserProfileNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UserProfileNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UserProfileGetParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	paramObj
}

type UserProfileUpdateParams struct {
	// If present, replaces the stored name. Omit to leave unchanged. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// If present, replaces the stored access type. Omit to leave unchanged.
	//
	// Any of "application", "passthrough".
	AccessType UserProfileAccessType `json:"access_type,omitzero"`
	// Details about the entity this profile represents, as the platform states them.
	// Each field sent replaces the stored value; omit a field to leave it unchanged.
	// Once set, a value cannot be cleared and `null` is rejected. Not accepted under
	// the legacy `user-profiles-2026-03-24` and `user-profiles-2026-08-18` beta
	// headers.
	ExternalUserDetails UserProfileExternalUserDetailsParams `json:"external_user_details,omitzero"`
	// Key-value pairs to merge into the stored metadata. Keys provided overwrite
	// existing values. To remove a key, set its value to an empty string. Keys not
	// provided are left unchanged. Maximum 16 keys, with keys up to 64 characters and
	// values up to 512 characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	paramObj
}

func (r UserProfileUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow UserProfileUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UserProfileUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UserProfileListParams struct {
	// The maximum number of user profiles to return, from 1 to 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The cursor for the page to return, taken from `next_page` in a previous
	// response.
	//
	// Leave it out to get the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Return only the user profiles whose `enforcement_status` is one of these values:
	// `unrestricted`, `pending`, `restricted` or `banned`. Repeat the parameter to
	// match any of several. With it set, a page may hold fewer than `limit` profiles,
	// or none, and still carry a `next_page`: the service examines a bounded number of
	// profiles per request, and only a null `next_page` means the listing is complete.
	// Rejected with `invalid_request_error` under the `user-profiles-2026-09-04` and
	// earlier beta headers. Fails with `overloaded_error` when a profile's enforcement
	// actions cannot be read for the moment.
	EnforcementStatus []UserProfileEnforcementStatus `query:"enforcement_status,omitzero" json:"-"`
	// The sort direction, applied to the field that `order_by` selects. Defaults to
	// `desc`.
	//
	// Any of "asc", "desc".
	Order UserProfileListOrder `query:"order,omitzero" json:"-"`
	// The field to sort user profiles by, in the direction that `order` sets. Defaults
	// to `created_at`.
	//
	// Any of "created_at", "name".
	OrderBy UserProfileListOrderBy `query:"order_by,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [UserProfileListParams]'s query parameters as `url.Values`.
func (r UserProfileListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type UserProfileDeleteParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	paramObj
}

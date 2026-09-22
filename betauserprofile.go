package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
)

// BetaUserProfileService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaUserProfileService] method instead.
type BetaUserProfileService struct {
	Options []option.RequestOption
}

// NewBetaUserProfileService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaUserProfileService(opts ...option.RequestOption) (r BetaUserProfileService) {
	r = BetaUserProfileService{}
	r.Options = opts
	return
}

// Create User Profile
func (r *BetaUserProfileService) New(ctx context.Context, params BetaUserProfileNewParams, opts ...option.RequestOption) (res *BetaUserProfile, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18")}, opts...)
	path := "v1/user_profiles?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get User Profile
func (r *BetaUserProfileService) Get(ctx context.Context, userProfileID string, query BetaUserProfileGetParams, opts ...option.RequestOption) (res *BetaUserProfile, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s?beta=true", userProfileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update User Profile
func (r *BetaUserProfileService) Update(ctx context.Context, userProfileID string, params BetaUserProfileUpdateParams, opts ...option.RequestOption) (res *BetaUserProfile, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s?beta=true", userProfileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List User Profiles
func (r *BetaUserProfileService) List(ctx context.Context, params BetaUserProfileListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaUserProfile], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/user_profiles?beta=true"
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
func (r *BetaUserProfileService) ListAutoPaging(ctx context.Context, params BetaUserProfileListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaUserProfile] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Create Enrollment URL
func (r *BetaUserProfileService) NewEnrollmentURL(ctx context.Context, userProfileID string, body BetaUserProfileNewEnrollmentURLParams, opts ...option.RequestOption) (res *BetaUserProfileEnrollmentURL, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s/enrollment_url?beta=true", userProfileID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// A record of an entity that the platform serves through the API, such as an
// end-user of the platform's product or a company that the platform resells Claude
// access to.
//
// A Messages, Message Batches or token counting request can send a profile's `id`
// in the `anthropic-user-profile-id` header to attribute the request to that
// entity.
type BetaUserProfile struct {
	// Unique identifier for this user profile, prefixed `uprof_`.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Trust grants for this profile, keyed by grant name. Key omitted when no grant is
	// active or in flight.
	TrustGrants map[string]BetaUserProfileTrustGrant `json:"trust_grants" api:"required"`
	// Object type. Always `user_profile`.
	//
	// Any of "user_profile".
	Type BetaUserProfileType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// How the platform uses the API on behalf of the entity this profile represents.
	// `application`: the platform sells a product that uses the API behind the scenes,
	// and the profile represents an individual end-user of that product.
	// `passthrough`: the platform resells raw inference, and the profile identifies
	// the resold-to company.
	//
	// Any of "application", "passthrough".
	AccessType BetaUserProfileAccessType `json:"access_type"`
	// Platform's own identifier for this user. Not enforced unique. Present under the
	// `user-profiles-2026-03-24` and `user-profiles-2026-08-18` beta headers; under
	// `user-profiles-2026-09-04` the value is `external_user_details.reference_id`.
	ExternalID string `json:"external_id" api:"nullable"`
	// Details about the entity this profile represents, as the platform states them.
	// Anthropic does not verify them. Every field is present, `null` until the
	// platform supplies a value.
	ExternalUserDetails BetaUserProfileExternalUserDetails `json:"external_user_details"`
	// A timestamp in RFC 3339 format
	ExternalUserOnboardedAt time.Time `json:"external_user_onboarded_at" api:"nullable" format:"date-time"`
	// Real-world name of the entity this profile represents (company or individual).
	// For a company the platform resells Claude access to (`access_type`
	// `passthrough`) this is that company's name.
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                      respjson.Field
		CreatedAt               respjson.Field
		Metadata                respjson.Field
		TrustGrants             respjson.Field
		Type                    respjson.Field
		UpdatedAt               respjson.Field
		AccessType              respjson.Field
		ExternalID              respjson.Field
		ExternalUserDetails     respjson.Field
		ExternalUserOnboardedAt respjson.Field
		Name                    respjson.Field
		ExtraFields             map[string]respjson.Field
		raw                     string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUserProfile) RawJSON() string { return r.JSON.raw }
func (r *BetaUserProfile) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Object type. Always `user_profile`.
type BetaUserProfileType string

const (
	BetaUserProfileTypeUserProfile BetaUserProfileType = "user_profile"
)

// How the platform uses the API on behalf of the entity this profile represents.
// `application`: the platform sells a product that uses the API behind the scenes,
// and the profile represents an individual end-user of that product.
// `passthrough`: the platform resells raw inference, and the profile identifies
// the resold-to company.
type BetaUserProfileAccessType string

const (
	// The user profile represents an individual end-user of a product that the
	// platform builds on the API. New profiles get this value by default.
	BetaUserProfileAccessTypeApplication BetaUserProfileAccessType = "application"
	// The user profile represents a company that the platform resells Claude access
	// to.
	BetaUserProfileAccessTypePassthrough BetaUserProfileAccessType = "passthrough"
)

// A URL to give to the entity that a user profile represents, so that the entity
// can enroll for a trust grant.
type BetaUserProfileEnrollmentURL struct {
	// A timestamp in RFC 3339 format
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Object type. Always `enrollment_url`.
	//
	// Any of "enrollment_url".
	Type BetaUserProfileEnrollmentURLType `json:"type" api:"required"`
	// Enrollment URL to send to the end user. Valid until `expires_at`.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt   respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUserProfileEnrollmentURL) RawJSON() string { return r.JSON.raw }
func (r *BetaUserProfileEnrollmentURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Object type. Always `enrollment_url`.
type BetaUserProfileEnrollmentURLType string

const (
	BetaUserProfileEnrollmentURLTypeEnrollmentURL BetaUserProfileEnrollmentURLType = "enrollment_url"
)

// Details about the entity this profile represents, as the platform states them.
// Anthropic does not verify them. Every field is present, `null` until the
// platform supplies a value.
type BetaUserProfileExternalUserDetails struct {
	// The status of the entity's account on the platform, as the platform states it:
	// `active`; `suspended`, when the platform has restricted the account and may
	// restore it; or `blocked`, when the platform has barred it. It records the
	// platform's decision only; the statuses in `trust_grants` are Anthropic's and do
	// not follow it.
	//
	// Any of "active", "suspended", "blocked".
	AccountStatus BetaUserProfileExternalUserDetailsAccountStatus `json:"account_status" api:"required"`
	// The country the platform associates with the entity, as an ISO 3166-1 alpha-2
	// code. `null` until the platform supplies one.
	Country string `json:"country" api:"required"`
	// The platform-computed hash of the entity's email address. `null` until the
	// platform supplies one.
	EmailHash string `json:"email_hash" api:"required"`
	// What kind of entity the profile represents, as the platform states it:
	// `individual`, `business`, `non_profit` or `government`.
	//
	// Any of "individual", "business", "non_profit", "government".
	EntityType BetaUserProfileExternalUserDetailsEntityType `json:"entity_type" api:"required"`
	// The platform-computed hash of the entity's name. `null` until the platform
	// supplies one.
	NameHash string `json:"name_hash" api:"required"`
	// A timestamp in RFC 3339 format
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
func (r BetaUserProfileExternalUserDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaUserProfileExternalUserDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The status of the entity's account on the platform, as the platform states it:
// `active`; `suspended`, when the platform has restricted the account and may
// restore it; or `blocked`, when the platform has barred it. It records the
// platform's decision only; the statuses in `trust_grants` are Anthropic's and do
// not follow it.
type BetaUserProfileExternalUserDetailsAccountStatus string

const (
	// The platform has neither restricted nor barred the account of the entity that
	// the user profile represents.
	BetaUserProfileExternalUserDetailsAccountStatusActive BetaUserProfileExternalUserDetailsAccountStatus = "active"
	// The platform has restricted the account of the entity that the user profile
	// represents and may restore it.
	BetaUserProfileExternalUserDetailsAccountStatusSuspended BetaUserProfileExternalUserDetailsAccountStatus = "suspended"
	// The platform has barred the account of the entity that the user profile
	// represents.
	BetaUserProfileExternalUserDetailsAccountStatusBlocked BetaUserProfileExternalUserDetailsAccountStatus = "blocked"
)

// What kind of entity the profile represents, as the platform states it:
// `individual`, `business`, `non_profit` or `government`.
type BetaUserProfileExternalUserDetailsEntityType string

const (
	BetaUserProfileExternalUserDetailsEntityTypeIndividual BetaUserProfileExternalUserDetailsEntityType = "individual"
	BetaUserProfileExternalUserDetailsEntityTypeBusiness   BetaUserProfileExternalUserDetailsEntityType = "business"
	BetaUserProfileExternalUserDetailsEntityTypeNonProfit  BetaUserProfileExternalUserDetailsEntityType = "non_profit"
	BetaUserProfileExternalUserDetailsEntityTypeGovernment BetaUserProfileExternalUserDetailsEntityType = "government"
)

type BetaUserProfileExternalUserDetailsParams struct {
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
	// A timestamp in RFC 3339 format
	OnboardedAt param.Opt[time.Time] `json:"onboarded_at,omitzero" format:"date-time"`
	// The status of the entity's account on the platform, as the platform states it:
	// `active`; `suspended`, when the platform has restricted the account and may
	// restore it; or `blocked`, when the platform has barred it. It records the
	// platform's decision only; the statuses in `trust_grants` are Anthropic's and do
	// not follow it.
	//
	// Any of "active", "suspended", "blocked".
	AccountStatus BetaUserProfileExternalUserDetailsParamsAccountStatus `json:"account_status,omitzero"`
	// What kind of entity the profile represents, as the platform states it:
	// `individual`, `business`, `non_profit` or `government`.
	//
	// Any of "individual", "business", "non_profit", "government".
	EntityType BetaUserProfileExternalUserDetailsParamsEntityType `json:"entity_type,omitzero"`
	paramObj
}

func (r BetaUserProfileExternalUserDetailsParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileExternalUserDetailsParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileExternalUserDetailsParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The status of the entity's account on the platform, as the platform states it:
// `active`; `suspended`, when the platform has restricted the account and may
// restore it; or `blocked`, when the platform has barred it. It records the
// platform's decision only; the statuses in `trust_grants` are Anthropic's and do
// not follow it.
type BetaUserProfileExternalUserDetailsParamsAccountStatus string

const (
	// The platform has neither restricted nor barred the account of the entity that
	// the user profile represents.
	BetaUserProfileExternalUserDetailsParamsAccountStatusActive BetaUserProfileExternalUserDetailsParamsAccountStatus = "active"
	// The platform has restricted the account of the entity that the user profile
	// represents and may restore it.
	BetaUserProfileExternalUserDetailsParamsAccountStatusSuspended BetaUserProfileExternalUserDetailsParamsAccountStatus = "suspended"
	// The platform has barred the account of the entity that the user profile
	// represents.
	BetaUserProfileExternalUserDetailsParamsAccountStatusBlocked BetaUserProfileExternalUserDetailsParamsAccountStatus = "blocked"
)

// What kind of entity the profile represents, as the platform states it:
// `individual`, `business`, `non_profit` or `government`.
type BetaUserProfileExternalUserDetailsParamsEntityType string

const (
	BetaUserProfileExternalUserDetailsParamsEntityTypeIndividual BetaUserProfileExternalUserDetailsParamsEntityType = "individual"
	BetaUserProfileExternalUserDetailsParamsEntityTypeBusiness   BetaUserProfileExternalUserDetailsParamsEntityType = "business"
	BetaUserProfileExternalUserDetailsParamsEntityTypeNonProfit  BetaUserProfileExternalUserDetailsParamsEntityType = "non_profit"
	BetaUserProfileExternalUserDetailsParamsEntityTypeGovernment BetaUserProfileExternalUserDetailsParamsEntityType = "government"
)

// The status of one trust grant on a user profile, listed in the profile's
// `trust_grants` map under the grant's name.
type BetaUserProfileTrustGrant struct {
	// Status of the trust grant.
	//
	// Any of "active", "pending", "rejected".
	Status BetaUserProfileTrustGrantStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUserProfileTrustGrant) RawJSON() string { return r.JSON.raw }
func (r *BetaUserProfileTrustGrant) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the trust grant.
type BetaUserProfileTrustGrantStatus string

const (
	BetaUserProfileTrustGrantStatusActive   BetaUserProfileTrustGrantStatus = "active"
	BetaUserProfileTrustGrantStatusPending  BetaUserProfileTrustGrantStatus = "pending"
	BetaUserProfileTrustGrantStatusRejected BetaUserProfileTrustGrantStatus = "rejected"
)

type BetaUserProfileNewParams struct {
	// Platform's own identifier for this user. Not enforced unique. Maximum 255
	// characters. Accepted under the `user-profiles-2026-03-24` and
	// `user-profiles-2026-08-18` beta headers; under `user-profiles-2026-09-04` send
	// `external_user_details.reference_id` instead.
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// Optional for all profiles. Real-world name of the entity this profile represents
	// (company or individual); for a company the platform resells Claude access to
	// (`access_type` `passthrough`), that company's name where known. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// A timestamp in RFC 3339 format
	ExternalUserOnboardedAt param.Opt[time.Time] `json:"external_user_onboarded_at,omitzero" format:"date-time"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// How the platform uses the API on behalf of the entity this profile represents.
	// `application`: the platform sells a product that uses the API behind the scenes,
	// and the profile represents an individual end-user of that product.
	// `passthrough`: the platform resells raw inference, and the profile identifies
	// the resold-to company.
	//
	// Any of "application", "passthrough".
	AccessType BetaUserProfileNewParamsAccessType `json:"access_type,omitzero"`
	// Details about the entity this profile represents, as the platform states them.
	// Every field is optional. Accepted under the `user-profiles-2026-09-04` beta
	// header only.
	ExternalUserDetails BetaUserProfileExternalUserDetailsParams `json:"external_user_details,omitzero"`
	// Free-form key-value data to attach to this user profile. Maximum 16 keys, with
	// keys up to 64 characters and values up to 512 characters. Values must be
	// non-empty strings.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaUserProfileNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How the platform uses the API on behalf of the entity this profile represents.
// `application`: the platform sells a product that uses the API behind the scenes,
// and the profile represents an individual end-user of that product.
// `passthrough`: the platform resells raw inference, and the profile identifies
// the resold-to company.
type BetaUserProfileNewParamsAccessType string

const (
	// The user profile represents an individual end-user of a product that the
	// platform builds on the API. New profiles get this value by default.
	BetaUserProfileNewParamsAccessTypeApplication BetaUserProfileNewParamsAccessType = "application"
	// The user profile represents a company that the platform resells Claude access
	// to.
	BetaUserProfileNewParamsAccessTypePassthrough BetaUserProfileNewParamsAccessType = "passthrough"
)

type BetaUserProfileGetParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaUserProfileUpdateParams struct {
	// If present, replaces the stored external_id. Omit to leave unchanged. Maximum
	// 255 characters. Accepted under the `user-profiles-2026-03-24` and
	// `user-profiles-2026-08-18` beta headers; under `user-profiles-2026-09-04` send
	// `external_user_details.reference_id` instead.
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// If present, replaces the stored name. Omit to leave unchanged. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// A timestamp in RFC 3339 format
	ExternalUserOnboardedAt param.Opt[time.Time] `json:"external_user_onboarded_at,omitzero" format:"date-time"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// How the platform uses the API on behalf of the entity this profile represents.
	// `application`: the platform sells a product that uses the API behind the scenes,
	// and the profile represents an individual end-user of that product.
	// `passthrough`: the platform resells raw inference, and the profile identifies
	// the resold-to company.
	//
	// Any of "application", "passthrough".
	AccessType BetaUserProfileUpdateParamsAccessType `json:"access_type,omitzero"`
	// Details about the entity this profile represents, as the platform states them.
	// Each field sent replaces the stored value; omit a field to leave it unchanged.
	// Once set, a value cannot be cleared and `null` is rejected. Accepted under the
	// `user-profiles-2026-09-04` beta header only.
	ExternalUserDetails BetaUserProfileExternalUserDetailsParams `json:"external_user_details,omitzero"`
	// Key-value pairs to merge into the stored metadata. Keys provided overwrite
	// existing values. To remove a key, set its value to an empty string. Keys not
	// provided are left unchanged. Maximum 16 keys, with keys up to 64 characters and
	// values up to 512 characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaUserProfileUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How the platform uses the API on behalf of the entity this profile represents.
// `application`: the platform sells a product that uses the API behind the scenes,
// and the profile represents an individual end-user of that product.
// `passthrough`: the platform resells raw inference, and the profile identifies
// the resold-to company.
type BetaUserProfileUpdateParamsAccessType string

const (
	// The user profile represents an individual end-user of a product that the
	// platform builds on the API. New profiles get this value by default.
	BetaUserProfileUpdateParamsAccessTypeApplication BetaUserProfileUpdateParamsAccessType = "application"
	// The user profile represents a company that the platform resells Claude access
	// to.
	BetaUserProfileUpdateParamsAccessTypePassthrough BetaUserProfileUpdateParamsAccessType = "passthrough"
)

type BetaUserProfileListParams struct {
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
	// The sort direction, applied to the field that `order_by` selects. Defaults to
	// `desc`.
	//
	// Any of "asc", "desc".
	Order BetaUserProfileListParamsOrder `query:"order,omitzero" json:"-"`
	// The field to sort user profiles by, in the direction that `order` sets. Defaults
	// to `created_at`.
	//
	// Any of "created_at", "name".
	OrderBy BetaUserProfileListParamsOrderBy `query:"order_by,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaUserProfileListParams]'s query parameters as
// `url.Values`.
func (r BetaUserProfileListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// The sort direction, applied to the field that `order_by` selects. Defaults to
// `desc`.
type BetaUserProfileListParamsOrder string

const (
	// Oldest first when `order_by` is `created_at`, or names in ascending order when
	// `order_by` is `name`.
	BetaUserProfileListParamsOrderAsc BetaUserProfileListParamsOrder = "asc"
	// Newest first when `order_by` is `created_at`, or names in descending order when
	// `order_by` is `name`. This is the default.
	BetaUserProfileListParamsOrderDesc BetaUserProfileListParamsOrder = "desc"
)

// The field to sort user profiles by, in the direction that `order` sets. Defaults
// to `created_at`.
type BetaUserProfileListParamsOrderBy string

const (
	// Sort by when each user profile was created. This is the default.
	BetaUserProfileListParamsOrderByCreatedAt BetaUserProfileListParamsOrderBy = "created_at"
	// Sort by `name`, ignoring the case of ASCII letters. Profiles without a name come
	// last in either direction.
	BetaUserProfileListParamsOrderByName BetaUserProfileListParamsOrderBy = "name"
)

type BetaUserProfileNewEnrollmentURLParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

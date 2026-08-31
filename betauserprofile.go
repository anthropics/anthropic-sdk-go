// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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
	// Platform's own identifier for this user. Not enforced unique.
	ExternalID string `json:"external_id" api:"nullable"`
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
	BetaUserProfileAccessTypeApplication BetaUserProfileAccessType = "application"
	BetaUserProfileAccessTypePassthrough BetaUserProfileAccessType = "passthrough"
)

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
	// characters.
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// Optional for all profiles. Real-world name of the entity this profile represents
	// (company or individual); for a company the platform resells Claude access to
	// (`access_type` `passthrough`), that company's name where known. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// A timestamp in RFC 3339 format
	ExternalUserOnboardedAt param.Opt[time.Time] `json:"external_user_onboarded_at,omitzero" format:"date-time"`
	// How the platform uses the API on behalf of the entity this profile represents.
	// `application`: the platform sells a product that uses the API behind the scenes,
	// and the profile represents an individual end-user of that product.
	// `passthrough`: the platform resells raw inference, and the profile identifies
	// the resold-to company.
	//
	// Any of "application", "passthrough".
	AccessType BetaUserProfileNewParamsAccessType `json:"access_type,omitzero"`
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
	BetaUserProfileNewParamsAccessTypeApplication BetaUserProfileNewParamsAccessType = "application"
	BetaUserProfileNewParamsAccessTypePassthrough BetaUserProfileNewParamsAccessType = "passthrough"
)

type BetaUserProfileGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaUserProfileUpdateParams struct {
	// If present, replaces the stored external_id. Omit to leave unchanged. Maximum
	// 255 characters.
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// If present, replaces the stored name. Omit to leave unchanged. Maximum 255
	// characters.
	Name param.Opt[string] `json:"name,omitzero"`
	// A timestamp in RFC 3339 format
	ExternalUserOnboardedAt param.Opt[time.Time] `json:"external_user_onboarded_at,omitzero" format:"date-time"`
	// How the platform uses the API on behalf of the entity this profile represents.
	// `application`: the platform sells a product that uses the API behind the scenes,
	// and the profile represents an individual end-user of that product.
	// `passthrough`: the platform resells raw inference, and the profile identifies
	// the resold-to company.
	//
	// Any of "application", "passthrough".
	AccessType BetaUserProfileUpdateParamsAccessType `json:"access_type,omitzero"`
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
	BetaUserProfileUpdateParamsAccessTypeApplication BetaUserProfileUpdateParamsAccessType = "application"
	BetaUserProfileUpdateParamsAccessTypePassthrough BetaUserProfileUpdateParamsAccessType = "passthrough"
)

type BetaUserProfileListParams struct {
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Query parameter for order
	//
	// Any of "asc", "desc".
	Order BetaUserProfileListParamsOrder `query:"order,omitzero" json:"-"`
	// Query parameter for order_by
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

// Query parameter for order
type BetaUserProfileListParamsOrder string

const (
	BetaUserProfileListParamsOrderAsc  BetaUserProfileListParamsOrder = "asc"
	BetaUserProfileListParamsOrderDesc BetaUserProfileListParamsOrder = "desc"
)

// Query parameter for order_by
type BetaUserProfileListParamsOrderBy string

const (
	BetaUserProfileListParamsOrderByCreatedAt BetaUserProfileListParamsOrderBy = "created_at"
	BetaUserProfileListParamsOrderByName      BetaUserProfileListParamsOrderBy = "name"
)

type BetaUserProfileNewEnrollmentURLParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

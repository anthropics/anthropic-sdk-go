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
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaUserProfileActionURLService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaUserProfileActionURLService] method instead.
type BetaUserProfileActionURLService struct {
	Options []option.RequestOption
}

// NewBetaUserProfileActionURLService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaUserProfileActionURLService(opts ...option.RequestOption) (r BetaUserProfileActionURLService) {
	r = BetaUserProfileActionURLService{}
	r.Options = opts
	return
}

// Create Action URL
func (r *BetaUserProfileActionURLService) New(ctx context.Context, userProfileID string, params BetaUserProfileActionURLNewParams, opts ...option.RequestOption) (res *BetaUserProfileActionURL, err error) {
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
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-08-18")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s/action_urls?beta=true", url.PathEscape(userProfileID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.action_urls.create", "/v1/user_profiles/{user_profile_id}/action_urls"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaUserProfileActionUnionParam struct {
	OfEnrollEndUser *BetaUserProfileEnrollEndUserActionParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaUserProfileActionUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfEnrollEndUser)
}
func (u *BetaUserProfileActionUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaUserProfileActionUnionParam) asAny() any {
	if !param.IsOmitted(u.OfEnrollEndUser) {
		return u.OfEnrollEndUser
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaUserProfileActionUnionParam) GetType() *string {
	if vt := u.OfEnrollEndUser; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaUserProfileActionUnionParam](
		"type",
		apijson.Discriminator[BetaUserProfileEnrollEndUserActionParam]("enroll_end_user"),
	)
}

// A URL to give to the person who takes the requested action.
type BetaUserProfileActionURL struct {
	// When the URL stops working, in RFC 3339 format.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Object type. Always `user_profile_action_url`.
	Type constant.UserProfileActionURL `json:"type" default:"user_profile_action_url"`
	// The URL to give to the person who takes the action. Anyone who holds it can open
	// it until `expires_at`.
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
func (r BetaUserProfileActionURL) RawJSON() string { return r.JSON.raw }
func (r *BetaUserProfileActionURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaUserProfileEnrollEndUserActionParam() BetaUserProfileEnrollEndUserActionParam {
	return BetaUserProfileEnrollEndUserActionParam{
		Type: "enroll_end_user",
	}
}

// The end-user that the profile represents enrolls in a program.
//
// This struct has a constant value, construct it with
// [NewBetaUserProfileEnrollEndUserActionParam].
type BetaUserProfileEnrollEndUserActionParam struct {
	Type constant.EnrollEndUser `json:"type" default:"enroll_end_user"`
	paramObj
}

func (r BetaUserProfileEnrollEndUserActionParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileEnrollEndUserActionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileEnrollEndUserActionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaUserProfileActionURLNewParams struct {
	// The action that the person who opens the URL takes.
	Action BetaUserProfileActionUnionParam `json:"action,omitzero" api:"required"`
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

func (r BetaUserProfileActionURLNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileActionURLNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileActionURLNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

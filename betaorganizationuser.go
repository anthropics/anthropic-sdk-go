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

// BetaOrganizationUserService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationUserService] method instead.
type BetaOrganizationUserService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationUserService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationUserService(opts ...option.RequestOption) (r BetaOrganizationUserService) {
	r = BetaOrganizationUserService{}
	r.Options = opts
	return
}

// Retrieve a member of the organization by user ID.
func (r *BetaOrganizationUserService) Get(ctx context.Context, userID string, opts ...option.RequestOption) (res *BetaOrganizationUser, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", requestconfig.EscapePathSegment(userID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.users.retrieve", "/v1/organizations/users/{user_id}"), requestconfig.WithSpanAttribute("anthropic.organization.user.id", userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a member's organization role.
func (r *BetaOrganizationUserService) Update(ctx context.Context, userID string, body BetaOrganizationUserUpdateParams, opts ...option.RequestOption) (res *BetaOrganizationUser, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", requestconfig.EscapePathSegment(userID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.users.update", "/v1/organizations/users/{user_id}"), requestconfig.WithSpanAttribute("anthropic.organization.user.id", userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List the organization's members.
func (r *BetaOrganizationUserService) List(ctx context.Context, query BetaOrganizationUserListParams, opts ...option.RequestOption) (res *pagination.Page[BetaOrganizationUser], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/users"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.users.list", "/v1/organizations/users", "after_id", "before_id", "email", "limit", "roles"))
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

// List the organization's members.
func (r *BetaOrganizationUserService) ListAutoPaging(ctx context.Context, query BetaOrganizationUserListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaOrganizationUser] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

// Remove a member from the organization.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeCannotRemoveLastMember] (HTTP 400)
//   - [ErrorCodeConcurrentMembershipChange] (HTTP 409)
func (r *BetaOrganizationUserService) Remove(ctx context.Context, userID string, opts ...option.RequestOption) (res *BetaOrganizationUserRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", requestconfig.EscapePathSegment(userID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.users.remove", "/v1/organizations/users/{user_id}"), requestconfig.WithSpanAttribute("anthropic.organization.user.id", userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Sign a member out of Claude everywhere, without removing them from the
// organization. This ends their sign-in sessions and removes their trusted devices
// wherever they use Claude, including other organizations and individual plans. It
// also ends the OAuth tokens issued to them for this organization, so apps such as
// Claude Code must sign in again.
//
// A token can take up to an hour to stop working, though it almost always stops at
// once. The member keeps their role and can sign in again right away. To keep them
// out, remove them from the organization or, if your organization uses single
// sign-on, block them in your identity provider. Signing out does not end API
// keys, OAuth tokens that were not issued for this organization, or the member's
// Claude Code cloud sessions and routines.
//
// You cannot sign out the organization's primary owner or the person your API key
// or token acts for, and a request to sign out another owner, a member with owner
// permissions, or a membership admin may be refused. A refused request returns a
// `400` or `403` error that says why. If the request returns `503`, send it again
// to finish; anything already ended stays ended.
//
// Every request must include the beta header
// `anthropic-beta: admin-user-sign-out-2026-10-08`. A request without it returns
// `400`. Sign-out is in beta. Anthropic turns it on for each organization and has
// turned it on for every Claude Enterprise organization; until yours is on,
// requests also return `400`. In a Claude Enterprise organization, an API key
// needs the `write:members` scope; a key without it gets `403`.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaCannotSignOutPrimaryOwner] (HTTP 400)
//   - [ErrorCodeBetaCannotSignOutSelf] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEnrolledInBeta] (HTTP 400)
//   - [ErrorCodeBetaCannotSignOutManagedMember] (HTTP 403)
//   - [ErrorCodeBetaCannotSignOutMembershipAdmin] (HTTP 403)
//   - [ErrorCodeBetaCannotSignOutOwner] (HTTP 403)
//   - [ErrorCodeBetaCannotSignOutPrivilegedMember] (HTTP 403)
func (r *BetaOrganizationUserService) SignOut(ctx context.Context, userID string, body BetaOrganizationUserSignOutParams, opts ...option.RequestOption) (res *BetaOrganizationUser, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "admin-user-sign-out-2026-10-08")}, opts...)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s/sign_out", requestconfig.EscapePathSegment(userID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.users.sign_out", "/v1/organizations/users/{user_id}/sign_out"), requestconfig.WithSpanAttribute("anthropic.organization.user.id", userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type BetaOrganizationUser struct {
	// ID of the User.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the User joined the Organization.
	AddedAt time.Time `json:"added_at" api:"required" format:"date-time"`
	// Email of the User.
	Email string `json:"email" api:"required"`
	// Name of the User.
	Name string `json:"name" api:"required"`
	// Organization role of the User.
	//
	// Any of "admin", "billing", "claude_code_user", "developer", "managed",
	// "membership_admin", "owner", "primary_owner", "user".
	Role BetaOrganizationRole `json:"role" api:"required"`
	// Object type.
	//
	// For Users, this is always `"user"`.
	Type constant.User `json:"type" default:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AddedAt     respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationUser) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationUserRemoveResponse struct {
	// ID of the User.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For Users, this is always `"user_deleted"`.
	Type constant.UserDeleted `json:"type" default:"user_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationUserRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationUserRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationUserUpdateParams struct {
	// New role for the User.
	//
	// The accepted values depend on the organization type. Console and API
	// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
	// `admin` cannot be assigned through the API. Claude Enterprise organizations
	// accept `user` and `managed`; changing a User's role to `managed` requires the
	// `write:rbac_groups` scope in addition to `write:members`.
	//
	// Any of "billing", "claude_code_user", "developer", "managed", "user".
	Role BetaOrganizationUserUpdateParamsRole `json:"role,omitzero" api:"required"`
	paramObj
}

func (r BetaOrganizationUserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationUserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationUserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// New role for the User.
//
// The accepted values depend on the organization type. Console and API
// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
// `admin` cannot be assigned through the API. Claude Enterprise organizations
// accept `user` and `managed`; changing a User's role to `managed` requires the
// `write:rbac_groups` scope in addition to `write:members`.
type BetaOrganizationUserUpdateParamsRole string

const (
	BetaOrganizationUserUpdateParamsRoleBilling        BetaOrganizationUserUpdateParamsRole = "billing"
	BetaOrganizationUserUpdateParamsRoleClaudeCodeUser BetaOrganizationUserUpdateParamsRole = "claude_code_user"
	BetaOrganizationUserUpdateParamsRoleDeveloper      BetaOrganizationUserUpdateParamsRole = "developer"
	BetaOrganizationUserUpdateParamsRoleManaged        BetaOrganizationUserUpdateParamsRole = "managed"
	BetaOrganizationUserUpdateParamsRoleUser           BetaOrganizationUserUpdateParamsRole = "user"
)

type BetaOrganizationUserListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Filter by user email.
	Email param.Opt[string] `query:"email,omitzero" format:"email" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter to items whose `role` equals one of the supplied values. Repeatable;
	// values are OR'ed together.
	//
	// Accepted values depend on the organization type: Console and API organizations
	// accept `user`, `developer`, `billing`, `admin`, and `claude_code_user`; Claude
	// Enterprise organizations accept `user`, `owner`, `primary_owner`,
	// `membership_admin`, and `managed`.
	Roles []string `query:"roles,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationUserListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationUserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationUserSignOutParams struct {
	// This endpoint is in beta: requests must send `admin-user-sign-out-2026-10-08` in
	// this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

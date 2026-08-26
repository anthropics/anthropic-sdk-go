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
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationInviteService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationInviteService] method instead.
type BetaOrganizationInviteService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationInviteService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationInviteService(opts ...option.RequestOption) (r BetaOrganizationInviteService) {
	r = BetaOrganizationInviteService{}
	r.Options = opts
	return
}

// Invite a user to join the organization by email.
//
// On plans that draw members from a finite pool of purchased seats, the invite
// automatically consumes a seat from the lowest tier with availability; there is
// no seat-tier parameter. When no seat is free the request fails with a 400 error
// rather than purchasing a seat.
func (r *BetaOrganizationInviteService) New(ctx context.Context, body BetaOrganizationInviteNewParams, opts ...option.RequestOption) (res *BetaOrganizationInvite, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/invites?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an invite by ID.
func (r *BetaOrganizationInviteService) Get(ctx context.Context, inviteID string, opts ...option.RequestOption) (res *BetaOrganizationInvite, err error) {
	opts = slices.Concat(r.Options, opts)
	if inviteID == "" {
		err = errors.New("missing required invite_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/invites/%s?beta=true", inviteID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the organization's invites.
func (r *BetaOrganizationInviteService) List(ctx context.Context, query BetaOrganizationInviteListParams, opts ...option.RequestOption) (res *pagination.Page[BetaOrganizationInvite], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/invites?beta=true"
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

// List the organization's invites.
func (r *BetaOrganizationInviteService) ListAutoPaging(ctx context.Context, query BetaOrganizationInviteListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaOrganizationInvite] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a pending invite.
func (r *BetaOrganizationInviteService) Delete(ctx context.Context, inviteID string, opts ...option.RequestOption) (res *BetaOrganizationInviteDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if inviteID == "" {
		err = errors.New("missing required invite_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/invites/%s?beta=true", inviteID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaOrganizationInvite struct {
	// ID of the Invite.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the Invite was accepted, or null.
	AcceptedAt time.Time `json:"accepted_at" api:"required" format:"date-time"`
	// Email of the User being invited.
	Email string `json:"email" api:"required"`
	// RFC 3339 datetime string indicating when the Invite expires.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// RFC 3339 datetime string indicating when the Invite was created.
	InvitedAt time.Time `json:"invited_at" api:"required" format:"date-time"`
	// RBAC group IDs recorded on the Invite (Claude Enterprise organizations), to be
	// assigned to the User when the Invite is accepted. `[]` when none.
	RBACGroupIDs []string `json:"rbac_group_ids" api:"required"`
	// Organization role of the User.
	//
	// Any of "admin", "billing", "claude_code_user", "developer", "managed",
	// "membership_admin", "owner", "primary_owner", "user".
	Role BetaOrganizationRole `json:"role" api:"required"`
	// Status of the Invite.
	//
	// Any of "accepted", "deleted", "expired", "pending".
	Status BetaOrganizationInviteStatus `json:"status" api:"required"`
	// Object type.
	//
	// For Invites, this is always `"invite"`.
	Type constant.Invite `json:"type" default:"invite"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		AcceptedAt   respjson.Field
		Email        respjson.Field
		ExpiresAt    respjson.Field
		InvitedAt    respjson.Field
		RBACGroupIDs respjson.Field
		Role         respjson.Field
		Status       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationInvite) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationInvite) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the Invite.
type BetaOrganizationInviteStatus string

const (
	BetaOrganizationInviteStatusAccepted BetaOrganizationInviteStatus = "accepted"
	BetaOrganizationInviteStatusDeleted  BetaOrganizationInviteStatus = "deleted"
	BetaOrganizationInviteStatusExpired  BetaOrganizationInviteStatus = "expired"
	BetaOrganizationInviteStatusPending  BetaOrganizationInviteStatus = "pending"
)

type BetaOrganizationInviteDeleteResponse struct {
	// ID of the Invite.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For Invites, this is always `"invite_deleted"`.
	Type constant.InviteDeleted `json:"type" default:"invite_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationInviteDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationInviteDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationInviteNewParams struct {
	// Email of the User.
	Email string `json:"email" api:"required" format:"email"`
	// Role for the invited User.
	//
	// The accepted values depend on the organization type. Console and API
	// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
	// `admin` cannot be assigned through the API. Claude Enterprise organizations
	// accept `user` and `managed`.
	//
	// Any of "billing", "claude_code_user", "developer", "managed", "user".
	Role BetaOrganizationInviteNewParamsRole `json:"role,omitzero" api:"required"`
	// RBAC group IDs to assign to the User when the Invite is accepted. A non-empty
	// array is accepted only for a Claude Enterprise organization with RBAC groups,
	// and requires the key to carry the `write:rbac_groups` scope.
	RBACGroupIDs []string `json:"rbac_group_ids,omitzero"`
	paramObj
}

func (r BetaOrganizationInviteNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationInviteNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationInviteNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Role for the invited User.
//
// The accepted values depend on the organization type. Console and API
// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
// `admin` cannot be assigned through the API. Claude Enterprise organizations
// accept `user` and `managed`.
type BetaOrganizationInviteNewParamsRole string

const (
	BetaOrganizationInviteNewParamsRoleBilling        BetaOrganizationInviteNewParamsRole = "billing"
	BetaOrganizationInviteNewParamsRoleClaudeCodeUser BetaOrganizationInviteNewParamsRole = "claude_code_user"
	BetaOrganizationInviteNewParamsRoleDeveloper      BetaOrganizationInviteNewParamsRole = "developer"
	BetaOrganizationInviteNewParamsRoleManaged        BetaOrganizationInviteNewParamsRole = "managed"
	BetaOrganizationInviteNewParamsRoleUser           BetaOrganizationInviteNewParamsRole = "user"
)

type BetaOrganizationInviteListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Filter by the email address the Invite was sent to. Matches the same way as the
	// Users list's `email` filter (normalized, case-insensitive).
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
	// Filter by Invite status. Repeatable; values are OR'ed together. Omit to return
	// `pending`, `accepted`, and `expired` Invites alike.
	//
	// Any of "accepted", "expired", "pending".
	Statuses []string `query:"statuses,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationInviteListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationInviteListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

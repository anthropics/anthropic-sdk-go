// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"net/http"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationService] method instead.
type BetaOrganizationService struct {
	Options            []option.RequestOption
	APIKeys            BetaOrganizationAPIKeyService
	ExternalKeys       BetaOrganizationExternalKeyService
	Federation         BetaOrganizationFederationService
	Invites            BetaOrganizationInviteService
	ServiceAccounts    BetaOrganizationServiceAccountService
	Users              BetaOrganizationUserService
	Workspaces         BetaOrganizationWorkspaceService
	RateLimits         BetaOrganizationRateLimitService
	ComplianceSettings BetaOrganizationComplianceSettingService
}

// NewBetaOrganizationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationService(opts ...option.RequestOption) (r BetaOrganizationService) {
	r = BetaOrganizationService{}
	r.Options = opts
	r.APIKeys = NewBetaOrganizationAPIKeyService(opts...)
	r.ExternalKeys = NewBetaOrganizationExternalKeyService(opts...)
	r.Federation = NewBetaOrganizationFederationService(opts...)
	r.Invites = NewBetaOrganizationInviteService(opts...)
	r.ServiceAccounts = NewBetaOrganizationServiceAccountService(opts...)
	r.Users = NewBetaOrganizationUserService(opts...)
	r.Workspaces = NewBetaOrganizationWorkspaceService(opts...)
	r.RateLimits = NewBetaOrganizationRateLimitService(opts...)
	r.ComplianceSettings = NewBetaOrganizationComplianceSettingService(opts...)
	return
}

// Retrieve information about the organization associated with the authenticated
// API key.
func (r *BetaOrganizationService) Get(ctx context.Context, opts ...option.RequestOption) (res *BetaOrganization, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/me?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type BetaOrganization struct {
	// ID of the Organization.
	ID string `json:"id" api:"required" format:"uuid"`
	// Name of the Organization.
	Name string `json:"name" api:"required"`
	// Object type.
	//
	// For Organizations, this is always `"organization"`.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganization) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRole string

const (
	BetaOrganizationRoleAdmin           BetaOrganizationRole = "admin"
	BetaOrganizationRoleBilling         BetaOrganizationRole = "billing"
	BetaOrganizationRoleClaudeCodeUser  BetaOrganizationRole = "claude_code_user"
	BetaOrganizationRoleDeveloper       BetaOrganizationRole = "developer"
	BetaOrganizationRoleManaged         BetaOrganizationRole = "managed"
	BetaOrganizationRoleMembershipAdmin BetaOrganizationRole = "membership_admin"
	BetaOrganizationRoleOwner           BetaOrganizationRole = "owner"
	BetaOrganizationRolePrimaryOwner    BetaOrganizationRole = "primary_owner"
	BetaOrganizationRoleUser            BetaOrganizationRole = "user"
)

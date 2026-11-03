package anthropic

import (
	"context"
	"net/http"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationService] method instead.
type OrganizationService struct {
	Options            []option.RequestOption
	APIKeys            OrganizationAPIKeyService
	ExternalKeys       OrganizationExternalKeyService
	Federation         OrganizationFederationService
	Invites            OrganizationInviteService
	ServiceAccounts    OrganizationServiceAccountService
	Users              OrganizationUserService
	Workspaces         OrganizationWorkspaceService
	RateLimits         OrganizationRateLimitService
	ComplianceSettings OrganizationComplianceSettingService
	Analytics          OrganizationAnalyticsService
	SpendLimits        OrganizationSpendLimitService
	RBACGroups         OrganizationRBACGroupService
	RBACRoles          OrganizationRBACRoleService
}

// NewOrganizationService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewOrganizationService(opts ...option.RequestOption) (r OrganizationService) {
	r = OrganizationService{}
	r.Options = opts
	r.APIKeys = NewOrganizationAPIKeyService(opts...)
	r.ExternalKeys = NewOrganizationExternalKeyService(opts...)
	r.Federation = NewOrganizationFederationService(opts...)
	r.Invites = NewOrganizationInviteService(opts...)
	r.ServiceAccounts = NewOrganizationServiceAccountService(opts...)
	r.Users = NewOrganizationUserService(opts...)
	r.Workspaces = NewOrganizationWorkspaceService(opts...)
	r.RateLimits = NewOrganizationRateLimitService(opts...)
	r.ComplianceSettings = NewOrganizationComplianceSettingService(opts...)
	r.Analytics = NewOrganizationAnalyticsService(opts...)
	r.SpendLimits = NewOrganizationSpendLimitService(opts...)
	r.RBACGroups = NewOrganizationRBACGroupService(opts...)
	r.RBACRoles = NewOrganizationRBACRoleService(opts...)
	return
}

// Retrieve information about the organization associated with the authenticated
// API key.
func (r *OrganizationService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrganizationInfo, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/me"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.retrieve", "/v1/organizations/me"))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Another change to the organization's members is in progress. Retry the request
// once it completes.
type ConcurrentMembershipChangeErrorDetails struct {
	ErrorCode constant.ConcurrentMembershipChange `json:"error_code" default:"concurrent_membership_change"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConcurrentMembershipChangeErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ConcurrentMembershipChangeErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationInfo struct {
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
func (r OrganizationInfo) RawJSON() string { return r.JSON.raw }
func (r *OrganizationInfo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRole string

const (
	OrganizationRoleAdmin           OrganizationRole = "admin"
	OrganizationRoleBilling         OrganizationRole = "billing"
	OrganizationRoleClaudeCodeUser  OrganizationRole = "claude_code_user"
	OrganizationRoleDeveloper       OrganizationRole = "developer"
	OrganizationRoleManaged         OrganizationRole = "managed"
	OrganizationRoleMembershipAdmin OrganizationRole = "membership_admin"
	OrganizationRoleOwner           OrganizationRole = "owner"
	OrganizationRolePrimaryOwner    OrganizationRole = "primary_owner"
	OrganizationRoleUser            OrganizationRole = "user"
)

// The endpoint does not support the organization's type. Each endpoint's
// description says which organization types it supports.
type OrganizationTypeNotSupportedErrorDetails struct {
	ErrorCode constant.OrganizationTypeNotSupported `json:"error_code" default:"organization_type_not_supported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationTypeNotSupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *OrganizationTypeNotSupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A role attached to an RBAC Group in the request grants owner access, so this API
// does not add members to that group (directly or through an invite), remove its
// members, or delete it. Make the change in the organization's admin settings in
// Claude.
type RBACGroupGrantsOwnerAccessErrorDetails struct {
	ErrorCode constant.RBACGroupGrantsOwnerAccess `json:"error_code" default:"rbac_group_grants_owner_access"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACGroupGrantsOwnerAccessErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RBACGroupGrantsOwnerAccessErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API could not finish checking whether an RBAC Group in the request grants
// owner access (the check reached a fixed limit), so it made no change. Retrying
// the same request will not help. Make the change in the organization's admin
// settings in Claude.
type RBACGroupOwnerAccessUnconfirmedErrorDetails struct {
	ErrorCode constant.RBACGroupOwnerAccessUnconfirmed `json:"error_code" default:"rbac_group_owner_access_unconfirmed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACGroupOwnerAccessUnconfirmedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RBACGroupOwnerAccessUnconfirmedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The request depends on a user profile that does not exist in the organization or
// is no longer active. The request either names this user profile by ID, or names
// a resource, such as an External Key, that is bound to it.
type UserProfileNotFoundErrorDetails struct {
	ErrorCode constant.UserProfileNotFound `json:"error_code" default:"user_profile_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

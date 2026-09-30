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
	return
}

// Retrieve information about the organization associated with the authenticated
// API key.
func (r *OrganizationService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrganizationInfo, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/me"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
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

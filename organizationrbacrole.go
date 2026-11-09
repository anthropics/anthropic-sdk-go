package anthropic

import (
	"context"
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

// OrganizationRBACRoleService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationRBACRoleService] method instead.
type OrganizationRBACRoleService struct {
	Options     []option.RequestOption
	Permissions OrganizationRBACRolePermissionService
}

// NewOrganizationRBACRoleService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationRBACRoleService(opts ...option.RequestOption) (r OrganizationRBACRoleService) {
	r = OrganizationRBACRoleService{}
	r.Options = opts
	r.Permissions = NewOrganizationRBACRolePermissionService(opts...)
	return
}

// Retrieve an RBAC Role by ID.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACRoleService) Get(ctx context.Context, rbacRoleID string, opts ...option.RequestOption) (res *RBACRole, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacRoleID == "" {
		err = errors.New("missing required rbac_role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_roles/%s", requestconfig.EscapePathSegment(rbacRoleID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_roles.retrieve", "/v1/organizations/rbac_roles/{rbac_role_id}"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_role.id", rbacRoleID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List RBAC Roles in the organization.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACRoleService) List(ctx context.Context, query OrganizationRBACRoleListParams, opts ...option.RequestOption) (res *pagination.PageCursor[RBACRole], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rbac_roles"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_roles.list", "/v1/organizations/rbac_roles", "limit", "page"))
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

// List RBAC Roles in the organization.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACRoleService) ListAutoPaging(ctx context.Context, query OrganizationRBACRoleListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[RBACRole] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type RBACRole struct {
	// ID of the RBAC Role.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the RBAC Role was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Name of the RBAC Role. For a role created by Anthropic, this is Anthropic's name
	// for the access that the role grants. Anthropic may change this name, so to keep
	// a lasting reference to a role, store its `id`.
	DisplayName string `json:"display_name" api:"required"`
	// Object type.
	//
	// For RBAC Roles, this is always `"rbac_role"`.
	Type constant.RBACRole `json:"type" default:"rbac_role"`
	// RFC 3339 datetime string indicating when the RBAC Role was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		DisplayName respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACRole) RawJSON() string { return r.JSON.raw }
func (r *RBACRole) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRBACRoleListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationRBACRoleListParams]'s query parameters as
// `url.Values`.
func (r OrganizationRBACRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

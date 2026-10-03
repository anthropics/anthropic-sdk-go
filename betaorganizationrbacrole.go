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

// BetaOrganizationRBACRoleService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationRBACRoleService] method instead.
type BetaOrganizationRBACRoleService struct {
	Options     []option.RequestOption
	Permissions BetaOrganizationRBACRolePermissionService
}

// NewBetaOrganizationRBACRoleService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationRBACRoleService(opts ...option.RequestOption) (r BetaOrganizationRBACRoleService) {
	r = BetaOrganizationRBACRoleService{}
	r.Options = opts
	r.Permissions = NewBetaOrganizationRBACRolePermissionService(opts...)
	return
}

// Retrieve an RBAC Role by ID.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACRoleService) Get(ctx context.Context, rbacRoleID string, opts ...option.RequestOption) (res *BetaRBACRole, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacRoleID == "" {
		err = errors.New("missing required rbac_role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_roles/%s?beta=true", rbacRoleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List RBAC Roles in the organization.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACRoleService) List(ctx context.Context, query BetaOrganizationRBACRoleListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaRBACRole], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rbac_roles?beta=true"
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
func (r *BetaOrganizationRBACRoleService) ListAutoPaging(ctx context.Context, query BetaOrganizationRBACRoleListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaRBACRole] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaRBACRole struct {
	// ID of the RBAC Role.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the RBAC Role was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Name of the RBAC Role. For a role created by Anthropic, this name can differ
	// from the label claude.ai shows, and Anthropic may change the name. To keep a
	// lasting reference to a role, store its `id`.
	DisplayName string `json:"display_name" api:"required"`
	// Deprecated: use `display_name` instead. Name of the RBAC Role; always the same
	// value as `display_name`.
	//
	// Deprecated: Use `display_name` instead; `name` always has the same value.
	Name string `json:"name" api:"required"`
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
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACRole) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACRole) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACRoleListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationRBACRoleListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationRBACRoleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

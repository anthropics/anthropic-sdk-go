package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationRBACRolePermissionService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationRBACRolePermissionService] method instead.
type OrganizationRBACRolePermissionService struct {
	Options []option.RequestOption
}

// NewOrganizationRBACRolePermissionService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationRBACRolePermissionService(opts ...option.RequestOption) (r OrganizationRBACRolePermissionService) {
	r = OrganizationRBACRolePermissionService{}
	r.Options = opts
	return
}

// List the permissions an RBAC Role grants.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACRolePermissionService) List(ctx context.Context, rbacRoleID string, query OrganizationRBACRolePermissionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[RBACRolePermission], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if rbacRoleID == "" {
		err = errors.New("missing required rbac_role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_roles/%s/permissions", requestconfig.EscapePathSegment(rbacRoleID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_roles.permissions.list", "/v1/organizations/rbac_roles/{rbac_role_id}/permissions", "limit", "page"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_role.id", rbacRoleID))
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

// List the permissions an RBAC Role grants.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACRolePermissionService) ListAutoPaging(ctx context.Context, rbacRoleID string, query OrganizationRBACRolePermissionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[RBACRolePermission] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, rbacRoleID, query, opts...))
}

type RBACAllConnectorsPermissionResource struct {
	// Kind of resource the permission applies to.
	Type constant.AllConnectors `json:"type" default:"all_connectors"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACAllConnectorsPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *RBACAllConnectorsPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACConnectorPermissionResource struct {
	// ID of the connector the permission applies to.
	ConnectorID string `json:"connector_id" api:"required"`
	// Kind of resource the permission applies to.
	Type constant.Connector `json:"type" default:"connector"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorID respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACConnectorPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *RBACConnectorPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACConnectorScopePermissionResource struct {
	// ID of the connector the permission applies to.
	ConnectorID string `json:"connector_id" api:"required"`
	// OAuth scope the permission names — the role may receive this scope when tokens
	// are minted for the connector.
	//
	// Subject to the same encoding rule as `tool_name`: a scope containing characters
	// outside `[a-zA-Z0-9_-]` (or colliding with a reserved form) appears
	// server-encoded in a stable `{prefix}_{32-hex}` form. OAuth scopes routinely
	// contain `:` and `/`, so most appear encoded.
	Scope string `json:"scope" api:"required"`
	// Kind of resource the permission applies to.
	Type constant.ConnectorScope `json:"type" default:"connector_scope"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorID respjson.Field
		Scope       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACConnectorScopePermissionResource) RawJSON() string { return r.JSON.raw }
func (r *RBACConnectorScopePermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACConnectorToolPermissionResource struct {
	// ID of the connector the permission applies to.
	ConnectorID string `json:"connector_id" api:"required"`
	// Published name of the connector tool the permission applies to.
	//
	// When the published name contains characters outside `[a-zA-Z0-9_-]` (or collides
	// with a reserved form), it is server-encoded into a stable `{prefix}_{32-hex}`
	// form — a shortened readable prefix of the name plus a hash — from which the
	// published name is not recoverable.
	ToolName string `json:"tool_name" api:"required"`
	// Kind of resource the permission applies to.
	Type constant.ConnectorTool `json:"type" default:"connector_tool"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorID respjson.Field
		ToolName    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACConnectorToolPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *RBACConnectorToolPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACOrganizationPermissionResource struct {
	// UUID of the organization the permission applies to.
	OrganizationID string `json:"organization_id" api:"required"`
	// Kind of resource the permission applies to.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		OrganizationID respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACOrganizationPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *RBACOrganizationPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACRolePermission struct {
	// Action the permission grants on the resource.
	//
	// The vocabulary follows the resource: an `organization` grant carries a
	// product-feature entitlement (for example `chat`), an admin-panel permission
	// entitlement (`permission_*`), or a blanket capability-access mode —
	// `capability_access_all` grants every product-feature entitlement except
	// `inference_hooks_inspection`, which a role holds only through its own row, and
	// `capability_access_all_ga` grants the subset of those for stable features as it
	// stands at permission-check time; neither mode grants model-access entitlements.
	// A consumer enumerating a role's per-feature grants should treat a blanket row as
	// granting every product-feature entitlement it covers, or it will under-report
	// the role's effective access. A `connector_tool` grant carries a tool-access
	// action (`use` or `always_allow`); a `connector_scope` grant carries the scope
	// action `grant` (the role may receive the named OAuth scope when tokens are
	// minted for the connector); `connector` and `all_connectors` grants carry a
	// tool-access action, the scope action, or an authentication-method action
	// (`interactive` or `managed`).
	Action string `json:"action" api:"required"`
	// What the permission applies to.
	//
	// A tagged union: `type` names the kind of resource and determines which
	// identifier fields are present.
	Resource RBACRolePermissionResourceUnion `json:"resource" api:"required"`
	// Object type.
	//
	// For RBAC Role Permissions, this is always `"rbac_role_permission"`.
	Type constant.RBACRolePermission `json:"type" default:"rbac_role_permission"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Action      respjson.Field
		Resource    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACRolePermission) RawJSON() string { return r.JSON.raw }
func (r *RBACRolePermission) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// RBACRolePermissionResourceUnion contains all possible properties and values from
// [RBACOrganizationPermissionResource], [RBACConnectorToolPermissionResource],
// [RBACConnectorScopePermissionResource], [RBACConnectorPermissionResource],
// [RBACAllConnectorsPermissionResource].
//
// Use the [RBACRolePermissionResourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type RBACRolePermissionResourceUnion struct {
	// This field is from variant [RBACOrganizationPermissionResource].
	OrganizationID string `json:"organization_id"`
	// Any of "organization", "connector_tool", "connector_scope", "connector",
	// "all_connectors".
	Type        string `json:"type"`
	ConnectorID string `json:"connector_id"`
	// This field is from variant [RBACConnectorToolPermissionResource].
	ToolName string `json:"tool_name"`
	// This field is from variant [RBACConnectorScopePermissionResource].
	Scope string `json:"scope"`
	JSON  struct {
		OrganizationID respjson.Field
		Type           respjson.Field
		ConnectorID    respjson.Field
		ToolName       respjson.Field
		Scope          respjson.Field
		raw            string
	} `json:"-"`
}

// anyRBACRolePermissionResource is implemented by each variant of
// [RBACRolePermissionResourceUnion] to add type safety for the return type of
// [RBACRolePermissionResourceUnion.AsAny]
type anyRBACRolePermissionResource interface {
	implRBACRolePermissionResourceUnion()
}

func (RBACOrganizationPermissionResource) implRBACRolePermissionResourceUnion()   {}
func (RBACConnectorToolPermissionResource) implRBACRolePermissionResourceUnion()  {}
func (RBACConnectorScopePermissionResource) implRBACRolePermissionResourceUnion() {}
func (RBACConnectorPermissionResource) implRBACRolePermissionResourceUnion()      {}
func (RBACAllConnectorsPermissionResource) implRBACRolePermissionResourceUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := RBACRolePermissionResourceUnion.AsAny().(type) {
//	case anthropic.RBACOrganizationPermissionResource:
//	case anthropic.RBACConnectorToolPermissionResource:
//	case anthropic.RBACConnectorScopePermissionResource:
//	case anthropic.RBACConnectorPermissionResource:
//	case anthropic.RBACAllConnectorsPermissionResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u RBACRolePermissionResourceUnion) AsAny() anyRBACRolePermissionResource {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "connector_tool":
		return u.AsConnectorTool()
	case "connector_scope":
		return u.AsConnectorScope()
	case "connector":
		return u.AsConnector()
	case "all_connectors":
		return u.AsAllConnectors()
	}
	return nil
}

func (u RBACRolePermissionResourceUnion) AsOrganization() (v RBACOrganizationPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RBACRolePermissionResourceUnion) AsConnectorTool() (v RBACConnectorToolPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RBACRolePermissionResourceUnion) AsConnectorScope() (v RBACConnectorScopePermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RBACRolePermissionResourceUnion) AsConnector() (v RBACConnectorPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u RBACRolePermissionResourceUnion) AsAllConnectors() (v RBACAllConnectorsPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u RBACRolePermissionResourceUnion) RawJSON() string { return u.JSON.raw }

func (r *RBACRolePermissionResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRBACRolePermissionListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationRBACRolePermissionListParams]'s query
// parameters as `url.Values`.
func (r OrganizationRBACRolePermissionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

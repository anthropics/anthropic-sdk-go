package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationRBACRolePermissionService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationRBACRolePermissionService] method instead.
type BetaOrganizationRBACRolePermissionService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationRBACRolePermissionService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationRBACRolePermissionService(opts ...option.RequestOption) (r BetaOrganizationRBACRolePermissionService) {
	r = BetaOrganizationRBACRolePermissionService{}
	r.Options = opts
	return
}

// List the permissions an RBAC Role grants.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACRolePermissionService) List(ctx context.Context, rbacRoleID string, query BetaOrganizationRBACRolePermissionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaRBACRolePermission], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if rbacRoleID == "" {
		err = errors.New("missing required rbac_role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_roles/%s/permissions?beta=true", url.PathEscape(rbacRoleID))
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
func (r *BetaOrganizationRBACRolePermissionService) ListAutoPaging(ctx context.Context, rbacRoleID string, query BetaOrganizationRBACRolePermissionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaRBACRolePermission] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, rbacRoleID, query, opts...))
}

type BetaRBACAllConnectorsPermissionResource struct {
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
func (r BetaRBACAllConnectorsPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACAllConnectorsPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACConnectorPermissionResource struct {
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
func (r BetaRBACConnectorPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACConnectorPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACConnectorScopePermissionResource struct {
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
func (r BetaRBACConnectorScopePermissionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACConnectorScopePermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACConnectorToolPermissionResource struct {
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
func (r BetaRBACConnectorToolPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACConnectorToolPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACOrganizationPermissionResource struct {
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
func (r BetaRBACOrganizationPermissionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACOrganizationPermissionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACRolePermission struct {
	// Action the permission grants on the resource.
	//
	// The vocabulary follows the resource: an `organization` grant carries a
	// product-feature entitlement (for example `chat`), an admin-panel permission
	// entitlement (`permission_*`), or a blanket capability-access mode —
	// `capability_access_all` grants every product-feature entitlement, and
	// `capability_access_all_ga` grants the generally-available subset as it stands at
	// permission-check time; neither mode grants model-access entitlements. A consumer
	// enumerating a role's per-feature grants should treat a blanket row as granting
	// every product-feature entitlement it covers, or it will under-report the role's
	// effective access. A `connector_tool` grant carries a tool-access action (`use`
	// or `always_allow`); a `connector_scope` grant carries the scope action `grant`
	// (the role may receive the named OAuth scope when tokens are minted for the
	// connector); `connector` and `all_connectors` grants carry a tool-access action,
	// the scope action, or an authentication-method action (`interactive` or
	// `managed`).
	Action string `json:"action" api:"required"`
	// What the permission applies to.
	//
	// A tagged union: `type` names the kind of resource and determines which
	// identifier fields are present.
	Resource BetaRBACRolePermissionResourceUnion `json:"resource" api:"required"`
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
func (r BetaRBACRolePermission) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACRolePermission) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaRBACRolePermissionResourceUnion contains all possible properties and values
// from [BetaRBACOrganizationPermissionResource],
// [BetaRBACConnectorToolPermissionResource],
// [BetaRBACConnectorScopePermissionResource],
// [BetaRBACConnectorPermissionResource],
// [BetaRBACAllConnectorsPermissionResource].
//
// Use the [BetaRBACRolePermissionResourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaRBACRolePermissionResourceUnion struct {
	// This field is from variant [BetaRBACOrganizationPermissionResource].
	OrganizationID string `json:"organization_id"`
	// Any of "organization", "connector_tool", "connector_scope", "connector",
	// "all_connectors".
	Type        string `json:"type"`
	ConnectorID string `json:"connector_id"`
	// This field is from variant [BetaRBACConnectorToolPermissionResource].
	ToolName string `json:"tool_name"`
	// This field is from variant [BetaRBACConnectorScopePermissionResource].
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

// anyBetaRBACRolePermissionResource is implemented by each variant of
// [BetaRBACRolePermissionResourceUnion] to add type safety for the return type of
// [BetaRBACRolePermissionResourceUnion.AsAny]
type anyBetaRBACRolePermissionResource interface {
	implBetaRBACRolePermissionResourceUnion()
}

func (BetaRBACOrganizationPermissionResource) implBetaRBACRolePermissionResourceUnion()   {}
func (BetaRBACConnectorToolPermissionResource) implBetaRBACRolePermissionResourceUnion()  {}
func (BetaRBACConnectorScopePermissionResource) implBetaRBACRolePermissionResourceUnion() {}
func (BetaRBACConnectorPermissionResource) implBetaRBACRolePermissionResourceUnion()      {}
func (BetaRBACAllConnectorsPermissionResource) implBetaRBACRolePermissionResourceUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaRBACRolePermissionResourceUnion.AsAny().(type) {
//	case anthropic.BetaRBACOrganizationPermissionResource:
//	case anthropic.BetaRBACConnectorToolPermissionResource:
//	case anthropic.BetaRBACConnectorScopePermissionResource:
//	case anthropic.BetaRBACConnectorPermissionResource:
//	case anthropic.BetaRBACAllConnectorsPermissionResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaRBACRolePermissionResourceUnion) AsAny() anyBetaRBACRolePermissionResource {
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

func (u BetaRBACRolePermissionResourceUnion) AsOrganization() (v BetaRBACOrganizationPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRBACRolePermissionResourceUnion) AsConnectorTool() (v BetaRBACConnectorToolPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRBACRolePermissionResourceUnion) AsConnectorScope() (v BetaRBACConnectorScopePermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRBACRolePermissionResourceUnion) AsConnector() (v BetaRBACConnectorPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRBACRolePermissionResourceUnion) AsAllConnectors() (v BetaRBACAllConnectorsPermissionResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaRBACRolePermissionResourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaRBACRolePermissionResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACRolePermissionListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationRBACRolePermissionListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationRBACRolePermissionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

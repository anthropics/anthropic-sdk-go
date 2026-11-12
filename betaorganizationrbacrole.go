package anthropic

import (
	"context"
	"encoding/json"
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
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeOrganizationTypeNotSupported] (HTTP 400)
//   - [ErrorCodeParentOrganizationRequired] (HTTP 400)
//   - [ErrorCodeRBACRoleNotFound] (HTTP 404)
func (r *BetaOrganizationRBACRoleService) Get(ctx context.Context, rbacRoleID string, opts ...option.RequestOption) (res *BetaRBACRole, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacRoleID == "" {
		err = errors.New("missing required rbac_role_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_roles/%s", requestconfig.EscapePathSegment(rbacRoleID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_roles.retrieve", "/v1/organizations/rbac_roles/{rbac_role_id}"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_role.id", rbacRoleID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List RBAC Roles in the organization.
//
// The RBAC Roles API is available to Claude Enterprise organizations only.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeOrganizationTypeNotSupported] (HTTP 400)
//   - [ErrorCodeParentOrganizationRequired] (HTTP 400)
func (r *BetaOrganizationRBACRoleService) List(ctx context.Context, query BetaOrganizationRBACRoleListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaRBACRole], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rbac_roles"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
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
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeOrganizationTypeNotSupported] (HTTP 400)
//   - [ErrorCodeParentOrganizationRequired] (HTTP 400)
func (r *BetaOrganizationRBACRoleService) ListAutoPaging(ctx context.Context, query BetaOrganizationRBACRoleListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaRBACRole] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaRBACRole struct {
	// ID of the RBAC Role.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the RBAC Role was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Name of the RBAC Role. For a built-in role, this is Anthropic's name for the
	// access that the role grants. Anthropic may change this name, so to keep a
	// lasting reference to a role, store its `id`.
	DisplayName string `json:"display_name" api:"required"`
	// Whether the RBAC Role is built-in or custom. More values of `type` may be added,
	// so do not treat every role that is not `built_in` as `custom`.
	Source BetaRBACRoleSourceUnion `json:"source" api:"required"`
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
		Source      respjson.Field
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

// BetaRBACRoleSourceUnion contains all possible properties and values from
// [BetaRBACRoleSourceBuiltIn], [BetaRBACRoleSourceCustom].
//
// Use the [BetaRBACRoleSourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaRBACRoleSourceUnion struct {
	// Any of "built_in", "custom".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaRBACRoleSource is implemented by each variant of
// [BetaRBACRoleSourceUnion] to add type safety for the return type of
// [BetaRBACRoleSourceUnion.AsAny]
type anyBetaRBACRoleSource interface {
	implBetaRBACRoleSourceUnion()
}

func (BetaRBACRoleSourceBuiltIn) implBetaRBACRoleSourceUnion() {}
func (BetaRBACRoleSourceCustom) implBetaRBACRoleSourceUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaRBACRoleSourceUnion.AsAny().(type) {
//	case anthropic.BetaRBACRoleSourceBuiltIn:
//	case anthropic.BetaRBACRoleSourceCustom:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaRBACRoleSourceUnion) AsAny() anyBetaRBACRoleSource {
	switch u.Type {
	case "built_in":
		return u.AsBuiltIn()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u BetaRBACRoleSourceUnion) AsBuiltIn() (v BetaRBACRoleSourceBuiltIn) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRBACRoleSourceUnion) AsCustom() (v BetaRBACRoleSourceCustom) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaRBACRoleSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaRBACRoleSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACRoleSourceBuiltIn struct {
	// Anthropic provisions the RBAC Role, and the organization's administrators cannot
	// change or delete it.
	Type constant.BuiltIn `json:"type" default:"built_in"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACRoleSourceBuiltIn) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACRoleSourceBuiltIn) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRBACRoleSourceCustom struct {
	// The organization's administrators can change and delete the RBAC Role.
	Type constant.Custom `json:"type" default:"custom"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACRoleSourceCustom) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACRoleSourceCustom) UnmarshalJSON(data []byte) error {
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

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

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// OrganizationAPIKeyService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationAPIKeyService] method instead.
type OrganizationAPIKeyService struct {
	Options []option.RequestOption
}

// NewOrganizationAPIKeyService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationAPIKeyService(opts ...option.RequestOption) (r OrganizationAPIKeyService) {
	r = OrganizationAPIKeyService{}
	r.Options = opts
	return
}

// Get API Key
func (r *OrganizationAPIKeyService) Get(ctx context.Context, apiKeyID string, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/api_keys/%s", apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update API Key
func (r *OrganizationAPIKeyService) Update(ctx context.Context, apiKeyID string, body OrganizationAPIKeyUpdateParams, opts ...option.RequestOption) (res *APIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/api_keys/%s", apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List API Keys
func (r *OrganizationAPIKeyService) List(ctx context.Context, query OrganizationAPIKeyListParams, opts ...option.RequestOption) (res *pagination.Page[APIKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/api_keys"
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

// List API Keys
func (r *OrganizationAPIKeyService) ListAutoPaging(ctx context.Context, query OrganizationAPIKeyListParams, opts ...option.RequestOption) *pagination.PageAutoPager[APIKey] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

type APIKey struct {
	// ID of the API key.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the API Key was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The ID and type of the actor that created the API key, or `null` when the
	// creator is not recorded (legacy, workload-identity-federated, or system-created
	// keys).
	CreatedBy APIKeyCreatedBy `json:"created_by" api:"required"`
	// RFC 3339 datetime string indicating when the API Key expires, or `null` if it
	// never expires.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Name of the API key.
	Name string `json:"name" api:"required"`
	// Partially redacted hint for the API key.
	PartialKeyHint string `json:"partial_key_hint" api:"required"`
	// The principal the API key acts as (a User or a Service Account), or `null` if
	// the API key is not bound to a principal.
	Principal APIKeyPrincipalUnion `json:"principal" api:"required"`
	// Where the API key belongs: its Workspace
	// (`{"type": "workspace", "workspace_id": "wrkspc_..."}`, with the Workspace's
	// real ID even when it is the organization's default Workspace), or the
	// organization (`{"type": "organization"}`) for a principal-bound API key that has
	// no Workspace.
	Scope APIKeyScopeUnion `json:"scope" api:"required"`
	// Status of the API key.
	//
	// Any of "active", "archived", "expired", "inactive".
	Status APIKeyStatus `json:"status" api:"required"`
	// Object type.
	//
	// For API Keys, this is always `"api_key"`.
	Type constant.APIKey `json:"type" default:"api_key"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		CreatedBy      respjson.Field
		ExpiresAt      respjson.Field
		Name           respjson.Field
		PartialKeyHint respjson.Field
		Principal      respjson.Field
		Scope          respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKey) RawJSON() string { return r.JSON.raw }
func (r *APIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// APIKeyPrincipalUnion contains all possible properties and values from
// [APIKeyUserActor], [APIKeyServiceAccountActor].
//
// Use the [APIKeyPrincipalUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type APIKeyPrincipalUnion struct {
	// Any of "user_actor", "service_account_actor".
	Type string `json:"type"`
	// This field is from variant [APIKeyUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [APIKeyServiceAccountActor].
	ServiceAccountID string `json:"service_account_id"`
	JSON             struct {
		Type             respjson.Field
		UserID           respjson.Field
		ServiceAccountID respjson.Field
		raw              string
	} `json:"-"`
}

// anyAPIKeyPrincipal is implemented by each variant of [APIKeyPrincipalUnion] to
// add type safety for the return type of [APIKeyPrincipalUnion.AsAny]
type anyAPIKeyPrincipal interface {
	implAPIKeyPrincipalUnion()
}

func (APIKeyUserActor) implAPIKeyPrincipalUnion()           {}
func (APIKeyServiceAccountActor) implAPIKeyPrincipalUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := APIKeyPrincipalUnion.AsAny().(type) {
//	case anthropic.APIKeyUserActor:
//	case anthropic.APIKeyServiceAccountActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u APIKeyPrincipalUnion) AsAny() anyAPIKeyPrincipal {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "service_account_actor":
		return u.AsServiceAccountActor()
	}
	return nil
}

func (u APIKeyPrincipalUnion) AsUserActor() (v APIKeyUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u APIKeyPrincipalUnion) AsServiceAccountActor() (v APIKeyServiceAccountActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u APIKeyPrincipalUnion) RawJSON() string { return u.JSON.raw }

func (r *APIKeyPrincipalUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// APIKeyScopeUnion contains all possible properties and values from
// [APIKeyOrganizationScope], [APIKeyWorkspaceScope].
//
// Use the [APIKeyScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type APIKeyScopeUnion struct {
	// Any of "organization", "workspace".
	Type string `json:"type"`
	// This field is from variant [APIKeyWorkspaceScope].
	WorkspaceID string `json:"workspace_id"`
	JSON        struct {
		Type        respjson.Field
		WorkspaceID respjson.Field
		raw         string
	} `json:"-"`
}

// anyAPIKeyScope is implemented by each variant of [APIKeyScopeUnion] to add type
// safety for the return type of [APIKeyScopeUnion.AsAny]
type anyAPIKeyScope interface {
	implAPIKeyScopeUnion()
}

func (APIKeyOrganizationScope) implAPIKeyScopeUnion() {}
func (APIKeyWorkspaceScope) implAPIKeyScopeUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := APIKeyScopeUnion.AsAny().(type) {
//	case anthropic.APIKeyOrganizationScope:
//	case anthropic.APIKeyWorkspaceScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u APIKeyScopeUnion) AsAny() anyAPIKeyScope {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "workspace":
		return u.AsWorkspace()
	}
	return nil
}

func (u APIKeyScopeUnion) AsOrganization() (v APIKeyOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u APIKeyScopeUnion) AsWorkspace() (v APIKeyWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u APIKeyScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *APIKeyScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the API key.
type APIKeyStatus string

const (
	APIKeyStatusActive   APIKeyStatus = "active"
	APIKeyStatusArchived APIKeyStatus = "archived"
	APIKeyStatusExpired  APIKeyStatus = "expired"
	APIKeyStatusInactive APIKeyStatus = "inactive"
)

type APIKeyCreatedBy struct {
	// ID of the actor that created the object.
	ID string `json:"id" api:"required"`
	// Type of the actor that created the object.
	//
	// Any of "service_account", "user".
	Type APIKeyCreatedByType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyCreatedBy) RawJSON() string { return r.JSON.raw }
func (r *APIKeyCreatedBy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of the actor that created the object.
type APIKeyCreatedByType string

const (
	APIKeyCreatedByTypeServiceAccount APIKeyCreatedByType = "service_account"
	APIKeyCreatedByTypeUser           APIKeyCreatedByType = "user"
)

type APIKeyOrganizationScope struct {
	// Scope type. Always `"organization"`: the API key has no Workspace. Only a
	// principal-bound API key can have this scope.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyOrganizationScope) RawJSON() string { return r.JSON.raw }
func (r *APIKeyOrganizationScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyServiceAccountActor struct {
	// ID of the Service Account the API key acts as.
	ServiceAccountID string `json:"service_account_id" api:"required"`
	// Principal type. Always `"service_account_actor"` for a Service Account.
	Type constant.ServiceAccountActor `json:"type" default:"service_account_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyServiceAccountActor) RawJSON() string { return r.JSON.raw }
func (r *APIKeyServiceAccountActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyUserActor struct {
	// Principal type. Always `"user_actor"` for a User.
	Type constant.UserActor `json:"type" default:"user_actor"`
	// ID of the User the API key acts as.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyUserActor) RawJSON() string { return r.JSON.raw }
func (r *APIKeyUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type APIKeyWorkspaceScope struct {
	// Scope type. Always `"workspace"`: the API key belongs to one Workspace.
	Type constant.Workspace `json:"type" default:"workspace"`
	// ID of the Workspace the API key belongs to. Unlike the deprecated top-level
	// `workspace_id`, this is the Workspace's real ID even for the organization's
	// default Workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIKeyWorkspaceScope) RawJSON() string { return r.JSON.raw }
func (r *APIKeyWorkspaceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationAPIKeyUpdateParams struct {
	// Name of the API key.
	Name param.Opt[string] `json:"name,omitzero"`
	// Status of the API key.
	//
	// Any of "active", "archived", "inactive".
	Status OrganizationAPIKeyUpdateParamsStatus `json:"status,omitzero"`
	paramObj
}

func (r OrganizationAPIKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationAPIKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationAPIKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the API key.
type OrganizationAPIKeyUpdateParamsStatus string

const (
	OrganizationAPIKeyUpdateParamsStatusActive   OrganizationAPIKeyUpdateParamsStatus = "active"
	OrganizationAPIKeyUpdateParamsStatusArchived OrganizationAPIKeyUpdateParamsStatus = "archived"
	OrganizationAPIKeyUpdateParamsStatusInactive OrganizationAPIKeyUpdateParamsStatus = "inactive"
)

type OrganizationAPIKeyListParams struct {
	// Filter by the ID of the User who created the object.
	CreatedByUserID param.Opt[string] `query:"created_by_user_id,omitzero" json:"-"`
	// Filter by Workspace ID.
	WorkspaceID param.Opt[string] `query:"workspace_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter by API key status.
	//
	// Any of "active", "archived", "expired", "inactive".
	Status OrganizationAPIKeyListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationAPIKeyListParams]'s query parameters as
// `url.Values`.
func (r OrganizationAPIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by API key status.
type OrganizationAPIKeyListParamsStatus string

const (
	OrganizationAPIKeyListParamsStatusActive   OrganizationAPIKeyListParamsStatus = "active"
	OrganizationAPIKeyListParamsStatusArchived OrganizationAPIKeyListParamsStatus = "archived"
	OrganizationAPIKeyListParamsStatusExpired  OrganizationAPIKeyListParamsStatus = "expired"
	OrganizationAPIKeyListParamsStatusInactive OrganizationAPIKeyListParamsStatus = "inactive"
)

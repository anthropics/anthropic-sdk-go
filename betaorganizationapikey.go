// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

// BetaOrganizationAPIKeyService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAPIKeyService] method instead.
type BetaOrganizationAPIKeyService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAPIKeyService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationAPIKeyService(opts ...option.RequestOption) (r BetaOrganizationAPIKeyService) {
	r = BetaOrganizationAPIKeyService{}
	r.Options = opts
	return
}

// Get API Key
func (r *BetaOrganizationAPIKeyService) Get(ctx context.Context, apiKeyID string, opts ...option.RequestOption) (res *BetaAPIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/api_keys/%s?beta=true", apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update API Key
func (r *BetaOrganizationAPIKeyService) Update(ctx context.Context, apiKeyID string, body BetaOrganizationAPIKeyUpdateParams, opts ...option.RequestOption) (res *BetaAPIKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if apiKeyID == "" {
		err = errors.New("missing required api_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/api_keys/%s?beta=true", apiKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List API Keys
func (r *BetaOrganizationAPIKeyService) List(ctx context.Context, query BetaOrganizationAPIKeyListParams, opts ...option.RequestOption) (res *pagination.Page[BetaAPIKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/api_keys?beta=true"
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
func (r *BetaOrganizationAPIKeyService) ListAutoPaging(ctx context.Context, query BetaOrganizationAPIKeyListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaAPIKey] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

type BetaAPIKey struct {
	// ID of the API key.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the API Key was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The ID and type of the actor that created the API key, or `null` when the
	// creator is not recorded (legacy, workload-identity-federated, or system-created
	// keys).
	CreatedBy BetaAPIKeyCreatedBy `json:"created_by" api:"required"`
	// RFC 3339 datetime string indicating when the API Key expires, or `null` if it
	// never expires.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Name of the API key.
	Name string `json:"name" api:"required"`
	// Partially redacted hint for the API key.
	PartialKeyHint string `json:"partial_key_hint" api:"required"`
	// The principal the API key acts as (a User or a Service Account), or `null` if
	// the API key is not bound to a principal.
	Principal BetaAPIKeyPrincipalUnion `json:"principal" api:"required"`
	// Where the API key belongs: its Workspace
	// (`{"type": "workspace", "workspace_id": "wrkspc_..."}`, with the Workspace's
	// real ID even when it is the organization's default Workspace), or the
	// organization (`{"type": "organization"}`) for a principal-bound API key that has
	// no Workspace.
	Scope BetaAPIKeyScopeUnion `json:"scope" api:"required"`
	// Status of the API key.
	//
	// Any of "active", "archived", "expired", "inactive".
	Status BetaAPIKeyStatus `json:"status" api:"required"`
	// Object type.
	//
	// For API Keys, this is always `"api_key"`.
	Type constant.APIKey `json:"type" default:"api_key"`
	// Deprecated: use `scope` instead. ID of the Workspace associated with the API
	// key, or `null` if the API key belongs to the default Workspace. Also `null` for
	// a principal-bound API key that has no Workspace; `scope` tells the two apart.
	//
	// Deprecated: Use `scope` instead. `workspace_id` is `null` both for an API key in
	// the default Workspace and for a principal-bound API key that has no Workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
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
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAPIKey) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaAPIKeyPrincipalUnion contains all possible properties and values from
// [BetaAPIKeyUserActor], [BetaAPIKeyServiceAccountActor].
//
// Use the [BetaAPIKeyPrincipalUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaAPIKeyPrincipalUnion struct {
	// Any of "user_actor", "service_account_actor".
	Type string `json:"type"`
	// This field is from variant [BetaAPIKeyUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaAPIKeyServiceAccountActor].
	ServiceAccountID string `json:"service_account_id"`
	JSON             struct {
		Type             respjson.Field
		UserID           respjson.Field
		ServiceAccountID respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaAPIKeyPrincipal is implemented by each variant of
// [BetaAPIKeyPrincipalUnion] to add type safety for the return type of
// [BetaAPIKeyPrincipalUnion.AsAny]
type anyBetaAPIKeyPrincipal interface {
	implBetaAPIKeyPrincipalUnion()
}

func (BetaAPIKeyUserActor) implBetaAPIKeyPrincipalUnion()           {}
func (BetaAPIKeyServiceAccountActor) implBetaAPIKeyPrincipalUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaAPIKeyPrincipalUnion.AsAny().(type) {
//	case anthropic.BetaAPIKeyUserActor:
//	case anthropic.BetaAPIKeyServiceAccountActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaAPIKeyPrincipalUnion) AsAny() anyBetaAPIKeyPrincipal {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "service_account_actor":
		return u.AsServiceAccountActor()
	}
	return nil
}

func (u BetaAPIKeyPrincipalUnion) AsUserActor() (v BetaAPIKeyUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaAPIKeyPrincipalUnion) AsServiceAccountActor() (v BetaAPIKeyServiceAccountActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaAPIKeyPrincipalUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaAPIKeyPrincipalUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaAPIKeyScopeUnion contains all possible properties and values from
// [BetaAPIKeyOrganizationScope], [BetaAPIKeyWorkspaceScope].
//
// Use the [BetaAPIKeyScopeUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaAPIKeyScopeUnion struct {
	// Any of "organization", "workspace".
	Type string `json:"type"`
	// This field is from variant [BetaAPIKeyWorkspaceScope].
	WorkspaceID string `json:"workspace_id"`
	JSON        struct {
		Type        respjson.Field
		WorkspaceID respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaAPIKeyScope is implemented by each variant of [BetaAPIKeyScopeUnion] to
// add type safety for the return type of [BetaAPIKeyScopeUnion.AsAny]
type anyBetaAPIKeyScope interface {
	implBetaAPIKeyScopeUnion()
}

func (BetaAPIKeyOrganizationScope) implBetaAPIKeyScopeUnion() {}
func (BetaAPIKeyWorkspaceScope) implBetaAPIKeyScopeUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaAPIKeyScopeUnion.AsAny().(type) {
//	case anthropic.BetaAPIKeyOrganizationScope:
//	case anthropic.BetaAPIKeyWorkspaceScope:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaAPIKeyScopeUnion) AsAny() anyBetaAPIKeyScope {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "workspace":
		return u.AsWorkspace()
	}
	return nil
}

func (u BetaAPIKeyScopeUnion) AsOrganization() (v BetaAPIKeyOrganizationScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaAPIKeyScopeUnion) AsWorkspace() (v BetaAPIKeyWorkspaceScope) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaAPIKeyScopeUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaAPIKeyScopeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the API key.
type BetaAPIKeyStatus string

const (
	BetaAPIKeyStatusActive   BetaAPIKeyStatus = "active"
	BetaAPIKeyStatusArchived BetaAPIKeyStatus = "archived"
	BetaAPIKeyStatusExpired  BetaAPIKeyStatus = "expired"
	BetaAPIKeyStatusInactive BetaAPIKeyStatus = "inactive"
)

type BetaAPIKeyCreatedBy struct {
	// ID of the actor that created the object.
	ID string `json:"id" api:"required"`
	// Type of the actor that created the object.
	//
	// Any of "service_account", "user".
	Type BetaAPIKeyCreatedByType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAPIKeyCreatedBy) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyCreatedBy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of the actor that created the object.
type BetaAPIKeyCreatedByType string

const (
	BetaAPIKeyCreatedByTypeServiceAccount BetaAPIKeyCreatedByType = "service_account"
	BetaAPIKeyCreatedByTypeUser           BetaAPIKeyCreatedByType = "user"
)

type BetaAPIKeyOrganizationScope struct {
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
func (r BetaAPIKeyOrganizationScope) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyOrganizationScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAPIKeyServiceAccountActor struct {
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
func (r BetaAPIKeyServiceAccountActor) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyServiceAccountActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAPIKeyUserActor struct {
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
func (r BetaAPIKeyUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAPIKeyWorkspaceScope struct {
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
func (r BetaAPIKeyWorkspaceScope) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyWorkspaceScope) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationAPIKeyUpdateParams struct {
	// Name of the API key.
	Name param.Opt[string] `json:"name,omitzero"`
	// Status of the API key.
	//
	// Any of "active", "archived", "inactive".
	Status BetaOrganizationAPIKeyUpdateParamsStatus `json:"status,omitzero"`
	paramObj
}

func (r BetaOrganizationAPIKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationAPIKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationAPIKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of the API key.
type BetaOrganizationAPIKeyUpdateParamsStatus string

const (
	BetaOrganizationAPIKeyUpdateParamsStatusActive   BetaOrganizationAPIKeyUpdateParamsStatus = "active"
	BetaOrganizationAPIKeyUpdateParamsStatusArchived BetaOrganizationAPIKeyUpdateParamsStatus = "archived"
	BetaOrganizationAPIKeyUpdateParamsStatusInactive BetaOrganizationAPIKeyUpdateParamsStatus = "inactive"
)

type BetaOrganizationAPIKeyListParams struct {
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
	Status BetaOrganizationAPIKeyListParamsStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAPIKeyListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationAPIKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by API key status.
type BetaOrganizationAPIKeyListParamsStatus string

const (
	BetaOrganizationAPIKeyListParamsStatusActive   BetaOrganizationAPIKeyListParamsStatus = "active"
	BetaOrganizationAPIKeyListParamsStatusArchived BetaOrganizationAPIKeyListParamsStatus = "archived"
	BetaOrganizationAPIKeyListParamsStatusExpired  BetaOrganizationAPIKeyListParamsStatus = "expired"
	BetaOrganizationAPIKeyListParamsStatusInactive BetaOrganizationAPIKeyListParamsStatus = "inactive"
)

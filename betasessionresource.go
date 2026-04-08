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
	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
)

// BetaSessionResourceService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSessionResourceService] method instead.
type BetaSessionResourceService struct {
	Options []option.RequestOption
}

// NewBetaSessionResourceService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaSessionResourceService(opts ...option.RequestOption) (r BetaSessionResourceService) {
	r = BetaSessionResourceService{}
	r.Options = opts
	return
}

// Get Session Resource
func (r *BetaSessionResourceService) Get(ctx context.Context, resourceID string, params BetaSessionResourceGetParams, opts ...option.RequestOption) (res *BetaSessionResourceGetResponseUnion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/resources/%s?beta=true", params.SessionID, resourceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Session Resource
func (r *BetaSessionResourceService) Update(ctx context.Context, resourceID string, params BetaSessionResourceUpdateParams, opts ...option.RequestOption) (res *BetaSessionResourceUpdateResponseUnion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/resources/%s?beta=true", params.SessionID, resourceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Session Resources
func (r *BetaSessionResourceService) List(ctx context.Context, sessionID string, params BetaSessionResourceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsSessionResourceUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/resources?beta=true", sessionID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
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

// List Session Resources
func (r *BetaSessionResourceService) ListAutoPaging(ctx context.Context, sessionID string, params BetaSessionResourceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsSessionResourceUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Delete Session Resource
func (r *BetaSessionResourceService) Delete(ctx context.Context, resourceID string, params BetaSessionResourceDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeleteSessionResource, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/resources/%s?beta=true", params.SessionID, resourceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Add Session Resource
func (r *BetaSessionResourceService) Add(ctx context.Context, sessionID string, params BetaSessionResourceAddParams, opts ...option.RequestOption) (res *BetaManagedAgentsFileResource, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/resources?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Confirmation of resource deletion.
type BetaManagedAgentsDeleteSessionResource struct {
	ID string `json:"id" api:"required"`
	// Any of "session_resource_deleted".
	Type BetaManagedAgentsDeleteSessionResourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeleteSessionResource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeleteSessionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeleteSessionResourceType string

const (
	BetaManagedAgentsDeleteSessionResourceTypeSessionResourceDeleted BetaManagedAgentsDeleteSessionResourceType = "session_resource_deleted"
)

type BetaManagedAgentsFileResource struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	FileID    string    `json:"file_id" api:"required"`
	MountPath string    `json:"mount_path" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileResourceType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		FileID      respjson.Field
		MountPath   respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileResource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsFileResourceType string

const (
	BetaManagedAgentsFileResourceTypeFile BetaManagedAgentsFileResourceType = "file"
)

type BetaManagedAgentsGitHubRepositoryResource struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	MountPath string    `json:"mount_path" api:"required"`
	// Any of "github_repository".
	Type BetaManagedAgentsGitHubRepositoryResourceType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time                                              `json:"updated_at" api:"required" format:"date-time"`
	URL       string                                                 `json:"url" api:"required"`
	Checkout  BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		MountPath   respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		URL         respjson.Field
		Checkout    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsGitHubRepositoryResource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsGitHubRepositoryResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsGitHubRepositoryResourceType string

const (
	BetaManagedAgentsGitHubRepositoryResourceTypeGitHubRepository BetaManagedAgentsGitHubRepositoryResourceType = "github_repository"
)

// BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion contains all possible
// properties and values from [BetaManagedAgentsBranchCheckout],
// [BetaManagedAgentsCommitCheckout].
//
// Use the [BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion struct {
	// This field is from variant [BetaManagedAgentsBranchCheckout].
	Name string `json:"name"`
	// Any of "branch", "commit".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsCommitCheckout].
	Sha  string `json:"sha"`
	JSON struct {
		Name respjson.Field
		Type respjson.Field
		Sha  respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsGitHubRepositoryResourceCheckout is implemented by each
// variant of [BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion] to add type
// safety for the return type of
// [BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny]
type anyBetaManagedAgentsGitHubRepositoryResourceCheckout interface {
	implBetaManagedAgentsGitHubRepositoryResourceCheckoutUnion()
}

func (BetaManagedAgentsBranchCheckout) implBetaManagedAgentsGitHubRepositoryResourceCheckoutUnion() {}
func (BetaManagedAgentsCommitCheckout) implBetaManagedAgentsGitHubRepositoryResourceCheckoutUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsBranchCheckout:
//	case anthropic.BetaManagedAgentsCommitCheckout:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsAny() anyBetaManagedAgentsGitHubRepositoryResourceCheckout {
	switch u.Type {
	case "branch":
		return u.AsBranch()
	case "commit":
		return u.AsCommit()
	}
	return nil
}

func (u BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsBranch() (v BetaManagedAgentsBranchCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsCommit() (v BetaManagedAgentsCommitCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionResourceUnion contains all possible properties and
// values from [BetaManagedAgentsGitHubRepositoryResource],
// [BetaManagedAgentsFileResource].
//
// Use the [BetaManagedAgentsSessionResourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionResourceUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	Checkout BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [BetaManagedAgentsFileResource].
	FileID string `json:"file_id"`
	JSON   struct {
		ID        respjson.Field
		CreatedAt respjson.Field
		MountPath respjson.Field
		Type      respjson.Field
		UpdatedAt respjson.Field
		URL       respjson.Field
		Checkout  respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionResource is implemented by each variant of
// [BetaManagedAgentsSessionResourceUnion] to add type safety for the return type
// of [BetaManagedAgentsSessionResourceUnion.AsAny]
type anyBetaManagedAgentsSessionResource interface {
	implBetaManagedAgentsSessionResourceUnion()
}

func (BetaManagedAgentsGitHubRepositoryResource) implBetaManagedAgentsSessionResourceUnion() {}
func (BetaManagedAgentsFileResource) implBetaManagedAgentsSessionResourceUnion()             {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionResourceUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsGitHubRepositoryResource:
//	case anthropic.BetaManagedAgentsFileResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionResourceUnion) AsAny() anyBetaManagedAgentsSessionResource {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u BetaManagedAgentsSessionResourceUnion) AsGitHubRepository() (v BetaManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionResourceUnion) AsFile() (v BetaManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionResourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSessionResourceGetResponseUnion contains all possible properties and values
// from [BetaManagedAgentsGitHubRepositoryResource],
// [BetaManagedAgentsFileResource].
//
// Use the [BetaSessionResourceGetResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSessionResourceGetResponseUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	Checkout BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [BetaManagedAgentsFileResource].
	FileID string `json:"file_id"`
	JSON   struct {
		ID        respjson.Field
		CreatedAt respjson.Field
		MountPath respjson.Field
		Type      respjson.Field
		UpdatedAt respjson.Field
		URL       respjson.Field
		Checkout  respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaSessionResourceGetResponse is implemented by each variant of
// [BetaSessionResourceGetResponseUnion] to add type safety for the return type of
// [BetaSessionResourceGetResponseUnion.AsAny]
type anyBetaSessionResourceGetResponse interface {
	implBetaSessionResourceGetResponseUnion()
}

func (BetaManagedAgentsGitHubRepositoryResource) implBetaSessionResourceGetResponseUnion() {}
func (BetaManagedAgentsFileResource) implBetaSessionResourceGetResponseUnion()             {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSessionResourceGetResponseUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsGitHubRepositoryResource:
//	case anthropic.BetaManagedAgentsFileResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSessionResourceGetResponseUnion) AsAny() anyBetaSessionResourceGetResponse {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u BetaSessionResourceGetResponseUnion) AsGitHubRepository() (v BetaManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSessionResourceGetResponseUnion) AsFile() (v BetaManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSessionResourceGetResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSessionResourceGetResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSessionResourceUpdateResponseUnion contains all possible properties and
// values from [BetaManagedAgentsGitHubRepositoryResource],
// [BetaManagedAgentsFileResource].
//
// Use the [BetaSessionResourceUpdateResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSessionResourceUpdateResponseUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResource].
	Checkout BetaManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [BetaManagedAgentsFileResource].
	FileID string `json:"file_id"`
	JSON   struct {
		ID        respjson.Field
		CreatedAt respjson.Field
		MountPath respjson.Field
		Type      respjson.Field
		UpdatedAt respjson.Field
		URL       respjson.Field
		Checkout  respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaSessionResourceUpdateResponse is implemented by each variant of
// [BetaSessionResourceUpdateResponseUnion] to add type safety for the return type
// of [BetaSessionResourceUpdateResponseUnion.AsAny]
type anyBetaSessionResourceUpdateResponse interface {
	implBetaSessionResourceUpdateResponseUnion()
}

func (BetaManagedAgentsGitHubRepositoryResource) implBetaSessionResourceUpdateResponseUnion() {}
func (BetaManagedAgentsFileResource) implBetaSessionResourceUpdateResponseUnion()             {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSessionResourceUpdateResponseUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsGitHubRepositoryResource:
//	case anthropic.BetaManagedAgentsFileResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSessionResourceUpdateResponseUnion) AsAny() anyBetaSessionResourceUpdateResponse {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u BetaSessionResourceUpdateResponseUnion) AsGitHubRepository() (v BetaManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSessionResourceUpdateResponseUnion) AsFile() (v BetaManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSessionResourceUpdateResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSessionResourceUpdateResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionResourceGetParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSessionResourceUpdateParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// New authorization token for the resource. Currently only `github_repository`
	// resources support token rotation.
	AuthorizationToken string `json:"authorization_token" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSessionResourceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaSessionResourceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSessionResourceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionResourceListParams struct {
	// Maximum number of resources to return per page (max 1000). If omitted, returns
	// all resources.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's next_page field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSessionResourceListParams]'s query parameters as
// `url.Values`.
func (r BetaSessionResourceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaSessionResourceDeleteParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSessionResourceAddParams struct {
	// Mount a file uploaded via the Files API into the session.
	BetaManagedAgentsFileResourceParams BetaManagedAgentsFileResourceParams
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSessionResourceAddParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.BetaManagedAgentsFileResourceParams)
}
func (r *BetaSessionResourceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

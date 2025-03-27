// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/param"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
)

// ModelService contains methods and other services that help with interacting with
// the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewModelService] method instead.
type ModelService struct {
	Options []option.RequestOption
}

// NewModelService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewModelService(opts ...option.RequestOption) (r *ModelService) {
	r = &ModelService{}
	r.Options = opts
	return
}

// Get a specific model.
//
// The Models API response can be used to determine information about a specific
// model or resolve a model alias to a model ID.
func (r *ModelService) Get(ctx context.Context, modelID string, opts ...option.RequestOption) (res *ModelInfo, err error) {
	opts = append(r.Options[:], opts...)
	if modelID == "" {
		err = errors.New("missing required model_id parameter")
		return
	}
	path := fmt.Sprintf("v1/models/%s", modelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *ModelService) List(ctx context.Context, query ModelListParams, opts ...option.RequestOption) (res *pagination.Page[ModelInfo], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/models"
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

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *ModelService) ListAutoPaging(ctx context.Context, query ModelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ModelInfo] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

type ModelInfo struct {
	// Unique model identifier.
	ID string `json:"id,required"`
	// RFC 3339 datetime string representing the time at which the model was released.
	// May be set to an epoch value if the release date is unknown.
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// A human-readable name for the model.
	DisplayName string `json:"display_name,required"`
	// Object type.
	//
	// For Models, this is always `"model"`.
	Type ModelInfoType `json:"type,required"`
	JSON modelInfoJSON `json:"-"`
}

// modelInfoJSON contains the JSON metadata for the struct [ModelInfo]
type modelInfoJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	DisplayName apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ModelInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r modelInfoJSON) RawJSON() string {
	return r.raw
}

// Object type.
//
// For Models, this is always `"model"`.
type ModelInfoType string

const (
	ModelInfoTypeModel ModelInfoType = "model"
)

func (r ModelInfoType) IsKnown() bool {
	switch r {
	case ModelInfoTypeModel:
		return true
	}
	return false
}

type ModelListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Field[string] `query:"after_id"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Field[string] `query:"before_id"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Field[int64] `query:"limit"`
}

// URLQuery serializes [ModelListParams]'s query parameters as `url.Values`.
func (r ModelListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

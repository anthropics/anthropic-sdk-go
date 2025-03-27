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

// BetaModelService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaModelService] method instead.
type BetaModelService struct {
	Options []option.RequestOption
}

// NewBetaModelService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaModelService(opts ...option.RequestOption) (r *BetaModelService) {
	r = &BetaModelService{}
	r.Options = opts
	return
}

// Get a specific model.
//
// The Models API response can be used to determine information about a specific
// model or resolve a model alias to a model ID.
func (r *BetaModelService) Get(ctx context.Context, modelID string, opts ...option.RequestOption) (res *BetaModelInfo, err error) {
	opts = append(r.Options[:], opts...)
	if modelID == "" {
		err = errors.New("missing required model_id parameter")
		return
	}
	path := fmt.Sprintf("v1/models/%s?beta=true", modelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *BetaModelService) List(ctx context.Context, query BetaModelListParams, opts ...option.RequestOption) (res *pagination.Page[BetaModelInfo], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/models?beta=true"
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
func (r *BetaModelService) ListAutoPaging(ctx context.Context, query BetaModelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaModelInfo] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

type BetaModelInfo struct {
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
	Type BetaModelInfoType `json:"type,required"`
	JSON betaModelInfoJSON `json:"-"`
}

// betaModelInfoJSON contains the JSON metadata for the struct [BetaModelInfo]
type betaModelInfoJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	DisplayName apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaModelInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaModelInfoJSON) RawJSON() string {
	return r.raw
}

// Object type.
//
// For Models, this is always `"model"`.
type BetaModelInfoType string

const (
	BetaModelInfoTypeModel BetaModelInfoType = "model"
)

func (r BetaModelInfoType) IsKnown() bool {
	switch r {
	case BetaModelInfoTypeModel:
		return true
	}
	return false
}

type BetaModelListParams struct {
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

// URLQuery serializes [BetaModelListParams]'s query parameters as `url.Values`.
func (r BetaModelListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

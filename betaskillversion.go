// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apiform"
	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaSkillVersionService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSkillVersionService] method instead.
type BetaSkillVersionService struct {
	Options []option.RequestOption
}

// NewBetaSkillVersionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaSkillVersionService(opts ...option.RequestOption) (r BetaSkillVersionService) {
	r = BetaSkillVersionService{}
	r.Options = opts
	return
}

// Create Skill Version
func (r *BetaSkillVersionService) New(ctx context.Context, skillID string, params BetaSkillVersionNewParams, opts ...option.RequestOption) (res *BetaSkillVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s/versions?beta=true", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Skill Version
func (r *BetaSkillVersionService) Get(ctx context.Context, version string, params BetaSkillVersionGetParams, opts ...option.RequestOption) (res *BetaSkillVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if params.SkillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s/versions/%s?beta=true", params.SkillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Skill Versions
func (r *BetaSkillVersionService) List(ctx context.Context, skillID string, params BetaSkillVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaSkillVersion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s/versions?beta=true", skillID)
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

// List Skill Versions
func (r *BetaSkillVersionService) ListAutoPaging(ctx context.Context, skillID string, params BetaSkillVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaSkillVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, skillID, params, opts...))
}

// Delete Skill Version
func (r *BetaSkillVersionService) Delete(ctx context.Context, version string, params BetaSkillVersionDeleteParams, opts ...option.RequestOption) (res *BetaDeletedSkillVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if params.SkillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s/versions/%s?beta=true", params.SkillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Download a skill version's content as a zip archive.
func (r *BetaSkillVersionService) Download(ctx context.Context, version string, params BetaSkillVersionDownloadParams, opts ...option.RequestOption) (res *http.Response, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/binary")}, opts...)
	if params.SkillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s/versions/%s/content?beta=true", params.SkillID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type BetaDeletedSkillVersion struct {
	// Unique identifier for this Skill Version. The id addresses the version in paths
	// and pins it in references.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For Skill Versions, this is always `"skill_version_deleted"`.
	Type constant.SkillVersionDeleted `json:"type" default:"skill_version_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedSkillVersion) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedSkillVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSkillVersion struct {
	// Unique identifier for this Skill Version. The id addresses the version in paths
	// and pins it in references.
	ID string `json:"id" api:"required"`
	// ISO 8601 timestamp of when the skill was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Description of the skill version.
	//
	// This is extracted from the SKILL.md file in the skill upload.
	Description string `json:"description" api:"required"`
	// The Skill's immutable kebab-case slug, set at creation from the first upload's
	// SKILL.md frontmatter `name` (or its enclosing directory). Every later upload
	// must resolve to the same value. Also the top-level directory of the Skill's
	// mounted files and the base name of a downloaded archive.
	Name string `json:"name" api:"required"`
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	SkillID string `json:"skill_id" api:"required"`
	// Object type.
	//
	// For Skill Versions, this is always `"skill_version"`.
	Type constant.SkillVersion `json:"type" default:"skill_version"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		Name        respjson.Field
		SkillID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSkillVersion) RawJSON() string { return r.JSON.raw }
func (r *BetaSkillVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSkillVersionNewParams struct {
	// Files to upload for the skill.
	//
	// All files must be in the same top-level directory and must include a SKILL.md
	// file at the root of that directory.
	Files []io.Reader `json:"files,omitzero" api:"required" format:"binary"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSkillVersionNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err == nil {
		err = apiform.WriteExtras(writer, r.ExtraFields())
	}
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

type BetaSkillVersionGetParams struct {
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	SkillID string `path:"skill_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSkillVersionListParams struct {
	// Number of results to return per page.
	//
	// Ranges from `1` to `1000`. Defaults to `20`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSkillVersionListParams]'s query parameters as
// `url.Values`.
func (r BetaSkillVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaSkillVersionDeleteParams struct {
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	SkillID string `path:"skill_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSkillVersionDownloadParams struct {
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	SkillID string `path:"skill_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

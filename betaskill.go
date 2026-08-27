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

// BetaSkillService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSkillService] method instead.
type BetaSkillService struct {
	Options  []option.RequestOption
	Versions BetaSkillVersionService
}

// NewBetaSkillService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaSkillService(opts ...option.RequestOption) (r BetaSkillService) {
	r = BetaSkillService{}
	r.Options = opts
	r.Versions = NewBetaSkillVersionService(opts...)
	return
}

// Create Skill
func (r *BetaSkillService) New(ctx context.Context, params BetaSkillNewParams, opts ...option.RequestOption) (res *BetaSkill, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	path := "v1/skills?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Skill
func (r *BetaSkillService) Get(ctx context.Context, skillID string, query BetaSkillGetParams, opts ...option.RequestOption) (res *BetaSkill, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s?beta=true", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Skills
func (r *BetaSkillService) List(ctx context.Context, params BetaSkillListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaSkill], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/skills?beta=true"
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

// List Skills
func (r *BetaSkillService) ListAutoPaging(ctx context.Context, params BetaSkillListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaSkill] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete Skill
func (r *BetaSkillService) Delete(ctx context.Context, skillID string, body BetaSkillDeleteParams, opts ...option.RequestOption) (res *BetaDeletedSkill, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if skillID == "" {
		err = errors.New("missing required skill_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/skills/%s?beta=true", skillID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaDeletedSkill struct {
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For Skills, this is always `"skill_deleted"`.
	Type constant.SkillDeleted `json:"type" default:"skill_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedSkill) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSkill struct {
	// Unique identifier for the skill.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id" api:"required"`
	// ISO 8601 timestamp of when the skill was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable, single-line label for the Skill. Maximum 255 characters. Always
	// set: derived from the SKILL.md frontmatter `name` when omitted at creation. Not
	// unique.
	DisplayName string `json:"display_name" api:"required"`
	// ID of the newest Skill Version — what `latest` references resolve to. Always
	// set: a Skill holds at least one version.
	LatestVersionID string `json:"latest_version_id" api:"required"`
	// Where the Skill comes from.
	//
	// Possible values:
	//
	// - `"custom"`: authored by the platform user; private to their workspace
	// - `"anthropic"`: published by Anthropic; shared and read-only
	// - `"anthropic_example"`: Anthropic-published sample Skill
	// - `"plugin"`: resolved from an installed plugin
	Source BetaSkillSource `json:"source" api:"required"`
	// Object type.
	//
	// For Skills, this is always `"skill"`.
	Type constant.Skill `json:"type" default:"skill"`
	// ISO 8601 timestamp of when the skill was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		CreatedAt       respjson.Field
		DisplayName     respjson.Field
		LatestVersionID respjson.Field
		Source          respjson.Field
		Type            respjson.Field
		UpdatedAt       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSkill) RawJSON() string { return r.JSON.raw }
func (r *BetaSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSkillSource struct {
	// Where the Skill comes from.
	//
	// Possible values:
	//
	// - `"custom"`: authored by the platform user; private to their workspace
	// - `"anthropic"`: published by Anthropic; shared and read-only
	// - `"anthropic_example"`: Anthropic-published sample Skill
	// - `"plugin"`: resolved from an installed plugin
	//
	// Any of "custom", "anthropic", "anthropic_example", "plugin".
	Type BetaSkillSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSkillSource) RawJSON() string { return r.JSON.raw }
func (r *BetaSkillSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where the Skill comes from.
//
// Possible values:
//
// - `"custom"`: authored by the platform user; private to their workspace
// - `"anthropic"`: published by Anthropic; shared and read-only
// - `"anthropic_example"`: Anthropic-published sample Skill
// - `"plugin"`: resolved from an installed plugin
type BetaSkillSourceType string

const (
	BetaSkillSourceTypeCustom           BetaSkillSourceType = "custom"
	BetaSkillSourceTypeAnthropic        BetaSkillSourceType = "anthropic"
	BetaSkillSourceTypeAnthropicExample BetaSkillSourceType = "anthropic_example"
	BetaSkillSourceTypePlugin           BetaSkillSourceType = "plugin"
)

type BetaSkillNewParams struct {
	// Files to upload for the skill.
	//
	// All files must be in the same top-level directory and must include a SKILL.md
	// file at the root of that directory.
	Files []io.Reader `json:"files,omitzero" api:"required" format:"binary"`
	// Human-readable, single-line label for the Skill. Maximum 255 characters. Always
	// set: derived from the SKILL.md frontmatter `name` when omitted at creation. Not
	// unique.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSkillNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type BetaSkillGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSkillListParams struct {
	// Pagination token for fetching a specific page of results.
	//
	// Pass the value from a previous response's `next_page` field to get the next page
	// of results.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter skills by source.
	//
	// If provided, only skills from the specified source will be returned:
	//
	// - `"custom"`: only return user-created skills
	// - `"anthropic"`: only return Anthropic-created skills
	Source param.Opt[string] `query:"source,omitzero" json:"-"`
	// Number of results to return per page.
	//
	// Ranges from `1` to `1000`. Defaults to `20`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSkillListParams]'s query parameters as `url.Values`.
func (r BetaSkillListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaSkillDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

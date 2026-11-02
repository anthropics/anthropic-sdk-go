package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

// BetaOrganizationTelemetryDestinationWorkspaceService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationTelemetryDestinationWorkspaceService] method instead.
type BetaOrganizationTelemetryDestinationWorkspaceService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationTelemetryDestinationWorkspaceService generates a new service
// that applies the given options to each request. These options are applied after
// the parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationTelemetryDestinationWorkspaceService(opts ...option.RequestOption) (r BetaOrganizationTelemetryDestinationWorkspaceService) {
	r = BetaOrganizationTelemetryDestinationWorkspaceService{}
	r.Options = opts
	return
}

// List the workspaces on a destination's workspace list, newest entry first.
//
// A destination whose `workspaces.type` is `all` doesn't have a list, so its page
// is empty. The list includes archived workspaces.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-workspaces-for-a-destination)
// to learn more about choosing the workspaces of a destination.
func (r *BetaOrganizationTelemetryDestinationWorkspaceService) List(ctx context.Context, destinationID string, params BetaOrganizationTelemetryDestinationWorkspaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaTelemetryDestinationWorkspace], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "telemetry-destinations-2026-08-11"), option.WithResponseInto(&raw)}, opts...)
	if destinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s/workspaces", url.PathEscape(destinationID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.workspaces.list", "/v1/organizations/telemetry_destinations/{destination_id}/workspaces", "limit", "page"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", destinationID))
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

// List the workspaces on a destination's workspace list, newest entry first.
//
// A destination whose `workspaces.type` is `all` doesn't have a list, so its page
// is empty. The list includes archived workspaces.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-workspaces-for-a-destination)
// to learn more about choosing the workspaces of a destination.
func (r *BetaOrganizationTelemetryDestinationWorkspaceService) ListAutoPaging(ctx context.Context, destinationID string, params BetaOrganizationTelemetryDestinationWorkspaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaTelemetryDestinationWorkspace] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, destinationID, params, opts...))
}

// Add one workspace to a destination's workspace list.
//
// Under `include` the workspace's spans start going to the destination. Under
// `exclude` they stop. A list holds at most 1,000 workspaces.
//
// Adding a workspace that is already on the list changes nothing and returns the
// existing entry. An archived workspace that isn't on the list yet can't be added.
// A destination whose `workspaces.type` is `all` doesn't have a list, so the API
// returns 400. To change which workspaces such a destination covers, update its
// `workspaces`.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-workspaces-for-a-destination)
// to learn more about choosing the workspaces of a destination.
func (r *BetaOrganizationTelemetryDestinationWorkspaceService) Add(ctx context.Context, workspaceID string, params BetaOrganizationTelemetryDestinationWorkspaceAddParams, opts ...option.RequestOption) (res *BetaTelemetryDestinationWorkspace, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "telemetry-destinations-2026-08-11")}, opts...)
	if params.DestinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s/workspaces/%s", url.PathEscape(params.DestinationID), url.PathEscape(workspaceID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.workspaces.add", "/v1/organizations/telemetry_destinations/{destination_id}/workspaces/{workspace_id}"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", params.DestinationID), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.workspace.id", workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Remove one workspace from a destination's workspace list.
//
// Under `include` the workspace's spans stop going to the destination. Under
// `exclude` they start. An archived workspace can be removed too.
//
// The API returns 404 when the workspace isn't on the list. A destination whose
// `workspaces.type` is `all` doesn't have a list, so the API returns 400.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-workspaces-for-a-destination)
// to learn more about choosing the workspaces of a destination.
func (r *BetaOrganizationTelemetryDestinationWorkspaceService) Remove(ctx context.Context, workspaceID string, params BetaOrganizationTelemetryDestinationWorkspaceRemoveParams, opts ...option.RequestOption) (res *BetaDeletedTelemetryDestinationWorkspace, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "telemetry-destinations-2026-08-11")}, opts...)
	if params.DestinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s/workspaces/%s", url.PathEscape(params.DestinationID), url.PathEscape(workspaceID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.workspaces.remove", "/v1/organizations/telemetry_destinations/{destination_id}/workspaces/{workspace_id}"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", params.DestinationID), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.workspace.id", workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Confirmation that a workspace was removed from a destination's list.
type BetaDeletedTelemetryDestinationWorkspace struct {
	// ID of the telemetry destination (`tdst_...`).
	TelemetryDestinationID string                                        `json:"telemetry_destination_id" api:"required"`
	Type                   constant.TelemetryDestinationWorkspaceDeleted `json:"type" default:"telemetry_destination_workspace_deleted"`
	// ID of the workspace that was removed (`wrkspc_...`).
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		TelemetryDestinationID respjson.Field
		Type                   respjson.Field
		WorkspaceID            respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedTelemetryDestinationWorkspace) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedTelemetryDestinationWorkspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A workspace on a telemetry destination's workspace list.
//
// Under `include`, the destination covers the workspace. Under `exclude`, the
// destination covers every workspace except this one. The entry doesn't have an
// `id`. Its `telemetry_destination_id` and `workspace_id` together identify it.
type BetaTelemetryDestinationWorkspace struct {
	// Time the workspace was added to the list.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the telemetry destination (`tdst_...`).
	TelemetryDestinationID string                                 `json:"telemetry_destination_id" api:"required"`
	Type                   constant.TelemetryDestinationWorkspace `json:"type" default:"telemetry_destination_workspace"`
	// ID of the workspace (`wrkspc_...`).
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt              respjson.Field
		TelemetryDestinationID respjson.Field
		Type                   respjson.Field
		WorkspaceID            respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationWorkspace) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationWorkspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationTelemetryDestinationWorkspaceListParams struct {
	// Number of entries to return per page, between 1 and 1,000. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`. Leave it out to get the
	// first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11`
	// in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationTelemetryDestinationWorkspaceListParams]'s
// query parameters as `url.Values`.
func (r BetaOrganizationTelemetryDestinationWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationTelemetryDestinationWorkspaceAddParams struct {
	// ID of the destination (`tdst_...`).
	DestinationID string `path:"destination_id" api:"required" json:"-"`
	// This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11`
	// in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationTelemetryDestinationWorkspaceRemoveParams struct {
	// ID of the destination (`tdst_...`).
	DestinationID string `path:"destination_id" api:"required" json:"-"`
	// This endpoint is in beta: requests must send `telemetry-destinations-2026-08-11`
	// in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

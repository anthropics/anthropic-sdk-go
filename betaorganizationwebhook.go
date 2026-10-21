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

// BetaOrganizationWebhookService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationWebhookService] method instead.
type BetaOrganizationWebhookService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationWebhookService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationWebhookService(opts ...option.RequestOption) (r BetaOrganizationWebhookService) {
	r = BetaOrganizationWebhookService{}
	r.Options = opts
	return
}

// Create an organization webhook endpoint.
//
// An organization webhook endpoint belongs to the organization rather than to one
// of its workspaces. An organization can have at most 32 organization webhook
// endpoints, enabled or disabled.
//
// The response includes the endpoint's signing secret. No read returns it, so
// store it when you receive it.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) New(ctx context.Context, params BetaOrganizationWebhookNewParams, opts ...option.RequestOption) (res *BetaNewOrganizationWebhookEndpoint, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07")}, opts...)
	path := "v1/organizations/webhooks?beta=true"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.create", "/v1/organizations/webhooks"), requestconfig.WithSpanCreatedID("anthropic.organization.webhook.id", &res, func(created *BetaNewOrganizationWebhookEndpoint) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve an organization webhook endpoint by ID.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) Get(ctx context.Context, webhookEndpointID string, query BetaOrganizationWebhookGetParams, opts ...option.RequestOption) (res *BetaOrganizationWebhookEndpoint, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07")}, opts...)
	if webhookEndpointID == "" {
		err = errors.New("missing required webhook_endpoint_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/webhooks/%s?beta=true", url.PathEscape(webhookEndpointID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.retrieve", "/v1/organizations/webhooks/{webhook_endpoint_id}"), requestconfig.WithSpanAttribute("anthropic.organization.webhook.id", webhookEndpointID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an organization webhook endpoint. Only the fields you send are changed.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) Update(ctx context.Context, webhookEndpointID string, params BetaOrganizationWebhookUpdateParams, opts ...option.RequestOption) (res *BetaOrganizationWebhookEndpoint, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07")}, opts...)
	if webhookEndpointID == "" {
		err = errors.New("missing required webhook_endpoint_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/webhooks/%s?beta=true", url.PathEscape(webhookEndpointID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.update", "/v1/organizations/webhooks/{webhook_endpoint_id}"), requestconfig.WithSpanAttribute("anthropic.organization.webhook.id", webhookEndpointID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the organization's webhook endpoints.
//
// An organization webhook endpoint belongs to the organization rather than to one
// of its workspaces. The list does not include webhook endpoints that belong to a
// workspace.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) List(ctx context.Context, params BetaOrganizationWebhookListParams, opts ...option.RequestOption) (res *pagination.Page[BetaOrganizationWebhookEndpoint], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/webhooks?beta=true"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.list", "/v1/organizations/webhooks", "after_id", "before_id", "limit", "order"))
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

// List the organization's webhook endpoints.
//
// An organization webhook endpoint belongs to the organization rather than to one
// of its workspaces. The list does not include webhook endpoints that belong to a
// workspace.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) ListAutoPaging(ctx context.Context, params BetaOrganizationWebhookListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaOrganizationWebhookEndpoint] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// Delete an organization webhook endpoint.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) Delete(ctx context.Context, webhookEndpointID string, body BetaOrganizationWebhookDeleteParams, opts ...option.RequestOption) (res *BetaDeletedOrganizationWebhookEndpoint, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07")}, opts...)
	if webhookEndpointID == "" {
		err = errors.New("missing required webhook_endpoint_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/webhooks/%s?beta=true", url.PathEscape(webhookEndpointID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.delete", "/v1/organizations/webhooks/{webhook_endpoint_id}"), requestconfig.WithSpanAttribute("anthropic.organization.webhook.id", webhookEndpointID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Replace an organization webhook endpoint's signing secret. The old secret stops
// working immediately. The new secret is returned only in this response.
//
// Every request must include the beta header
// `anthropic-beta: webhooks-2026-10-07`. A request without it returns `404`,
// exactly as if this API did not exist. Anthropic turns organization webhook
// endpoints on for each organization; until yours is on, requests also return
// `404`. Organization webhook endpoints are in beta and are not available to
// organizations with HIPAA readiness enabled.
func (r *BetaOrganizationWebhookService) RegenerateSigningSecret(ctx context.Context, webhookEndpointID string, body BetaOrganizationWebhookRegenerateSigningSecretParams, opts ...option.RequestOption) (res *BetaOrganizationWebhookEndpointSigningSecret, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "webhooks-2026-10-07")}, opts...)
	if webhookEndpointID == "" {
		err = errors.New("missing required webhook_endpoint_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/webhooks/%s/regenerate_signing_secret?beta=true", url.PathEscape(webhookEndpointID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.webhooks.regenerate_signing_secret", "/v1/organizations/webhooks/{webhook_endpoint_id}/regenerate_signing_secret"), requestconfig.WithSpanAttribute("anthropic.organization.webhook.id", webhookEndpointID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// A deleted organization webhook endpoint.
type BetaDeletedOrganizationWebhookEndpoint struct {
	// ID of the deleted webhook endpoint.
	ID string `json:"id" api:"required"`
	// Always `webhook_endpoint_deleted`.
	Type constant.WebhookEndpointDeleted `json:"type" default:"webhook_endpoint_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedOrganizationWebhookEndpoint) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedOrganizationWebhookEndpoint) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A new organization webhook endpoint, with its signing secret.
type BetaNewOrganizationWebhookEndpoint struct {
	// ID of the webhook endpoint.
	ID string `json:"id" api:"required"`
	// When the endpoint was created, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Notes about the endpoint. An empty string when none were given.
	Description string `json:"description" api:"required"`
	// Why Anthropic disabled the endpoint automatically. An empty string when the
	// endpoint is enabled, or when you disabled it.
	DisabledReason string `json:"disabled_reason" api:"required"`
	// The event types delivered to the endpoint.
	//
	// Any of "user_profile.enforcement_action_opened",
	// "user_profile.enforcement_action_resolved".
	EnabledEvents []string `json:"enabled_events" api:"required"`
	// A short label for the endpoint. An empty string when none was given.
	Name string `json:"name" api:"required"`
	// The secret Anthropic signs deliveries to the endpoint with (prefix `whsec_`). No
	// other response returns it, so store it now. To replace it, regenerate the
	// signing secret.
	SigningSecret string `json:"signing_secret" api:"required"`
	// Whether deliveries are sent to the endpoint.
	//
	// Any of "disabled", "enabled".
	State BetaOrganizationWebhookEndpointState `json:"state" api:"required"`
	// When the endpoint was last changed, in RFC 3339 format.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The HTTPS URL that deliveries are sent to.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		Description    respjson.Field
		DisabledReason respjson.Field
		EnabledEvents  respjson.Field
		Name           respjson.Field
		SigningSecret  respjson.Field
		State          respjson.Field
		UpdatedAt      respjson.Field
		URL            respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaNewOrganizationWebhookEndpoint) RawJSON() string { return r.JSON.raw }
func (r *BetaNewOrganizationWebhookEndpoint) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A webhook endpoint that belongs to the organization rather than to one of its
// workspaces.
type BetaOrganizationWebhookEndpoint struct {
	// ID of the webhook endpoint.
	ID string `json:"id" api:"required"`
	// When the endpoint was created, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Notes about the endpoint. An empty string when none were given.
	Description string `json:"description" api:"required"`
	// Why Anthropic disabled the endpoint automatically. An empty string when the
	// endpoint is enabled, or when you disabled it.
	DisabledReason string `json:"disabled_reason" api:"required"`
	// The event types delivered to the endpoint.
	//
	// Any of "user_profile.enforcement_action_opened",
	// "user_profile.enforcement_action_resolved".
	EnabledEvents []string `json:"enabled_events" api:"required"`
	// A short label for the endpoint. An empty string when none was given.
	Name string `json:"name" api:"required"`
	// Whether deliveries are sent to the endpoint.
	//
	// Any of "disabled", "enabled".
	State BetaOrganizationWebhookEndpointState `json:"state" api:"required"`
	// When the endpoint was last changed, in RFC 3339 format.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The HTTPS URL that deliveries are sent to.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		Description    respjson.Field
		DisabledReason respjson.Field
		EnabledEvents  respjson.Field
		Name           respjson.Field
		State          respjson.Field
		UpdatedAt      respjson.Field
		URL            respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationWebhookEndpoint) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationWebhookEndpoint) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An organization webhook endpoint's new signing secret.
type BetaOrganizationWebhookEndpointSigningSecret struct {
	// The new signing secret. Deliveries are signed with it from now on, and no other
	// response returns it.
	SigningSecret string `json:"signing_secret" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SigningSecret respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationWebhookEndpointSigningSecret) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationWebhookEndpointSigningSecret) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Whether deliveries are sent to an organization webhook endpoint.
type BetaOrganizationWebhookEndpointState string

const (
	BetaOrganizationWebhookEndpointStateDisabled BetaOrganizationWebhookEndpointState = "disabled"
	BetaOrganizationWebhookEndpointStateEnabled  BetaOrganizationWebhookEndpointState = "enabled"
)

type BetaOrganizationWebhookNewParams struct {
	// The event types to deliver to the endpoint. Must list at least one.
	//
	// Any of "user_profile.enforcement_action_opened",
	// "user_profile.enforcement_action_resolved".
	EnabledEvents []string `json:"enabled_events,omitzero" api:"required"`
	// The HTTPS URL that deliveries are sent to.
	URL string `json:"url" api:"required"`
	// Notes about the endpoint.
	Description param.Opt[string] `json:"description,omitzero"`
	// A short label for the endpoint.
	Name param.Opt[string] `json:"name,omitzero"`
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationWebhookNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationWebhookNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationWebhookNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationWebhookGetParams struct {
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationWebhookUpdateParams struct {
	// Notes about the endpoint.
	Description param.Opt[string] `json:"description,omitzero"`
	// A short label for the endpoint.
	Name param.Opt[string] `json:"name,omitzero"`
	// The HTTPS URL that deliveries are sent to.
	URL param.Opt[string] `json:"url,omitzero"`
	// The event types to deliver to the endpoint. Replaces the whole set; an empty
	// list leaves it unchanged.
	//
	// Any of "user_profile.enforcement_action_opened",
	// "user_profile.enforcement_action_resolved".
	EnabledEvents []string `json:"enabled_events,omitzero"`
	// `disabled` stops deliveries. `enabled` resumes them and clears
	// `disabled_reason`. Events that happen while the endpoint is disabled are never
	// delivered.
	//
	// Any of "disabled", "enabled".
	State BetaOrganizationWebhookUpdateParamsState `json:"state,omitzero"`
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationWebhookUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationWebhookUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationWebhookUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `disabled` stops deliveries. `enabled` resumes them and clears
// `disabled_reason`. Events that happen while the endpoint is disabled are never
// delivered.
type BetaOrganizationWebhookUpdateParamsState string

const (
	BetaOrganizationWebhookUpdateParamsStateDisabled BetaOrganizationWebhookUpdateParamsState = "disabled"
	BetaOrganizationWebhookUpdateParamsStateEnabled  BetaOrganizationWebhookUpdateParamsState = "enabled"
)

type BetaOrganizationWebhookListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of a webhook endpoint to use as a cursor. Returns the first `limit` endpoints
	// that come before it in list order, not the `limit` endpoints just before it.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Number of items to return per page, up to `100`. The default, `0`, means `20`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort order by `created_at`. Defaults to `desc`, newest first.
	//
	// Any of "asc", "desc".
	Order BetaOrganizationWebhookListParamsOrder `query:"order,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationWebhookListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationWebhookListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort order by `created_at`. Defaults to `desc`, newest first.
type BetaOrganizationWebhookListParamsOrder string

const (
	BetaOrganizationWebhookListParamsOrderAsc  BetaOrganizationWebhookListParamsOrder = "asc"
	BetaOrganizationWebhookListParamsOrderDesc BetaOrganizationWebhookListParamsOrder = "desc"
)

type BetaOrganizationWebhookDeleteParams struct {
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationWebhookRegenerateSigningSecretParams struct {
	// This endpoint is in beta: requests must send `webhooks-2026-10-07` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

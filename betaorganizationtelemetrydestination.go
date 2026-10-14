package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationTelemetryDestinationService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationTelemetryDestinationService] method instead.
type BetaOrganizationTelemetryDestinationService struct {
	Options    []option.RequestOption
	Workspaces BetaOrganizationTelemetryDestinationWorkspaceService
}

// NewBetaOrganizationTelemetryDestinationService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationTelemetryDestinationService(opts ...option.RequestOption) (r BetaOrganizationTelemetryDestinationService) {
	r = BetaOrganizationTelemetryDestinationService{}
	r.Options = opts
	r.Workspaces = NewBetaOrganizationTelemetryDestinationWorkspaceService(opts...)
	return
}

// Register an OpenTelemetry destination for the organization.
//
// The destination starts in state `enabled`. An organization can have at most 10
// destinations, including disabled ones.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#add-a-destination)
// to learn more about adding a destination.
func (r *BetaOrganizationTelemetryDestinationService) New(ctx context.Context, params BetaOrganizationTelemetryDestinationNewParams, opts ...option.RequestOption) (res *BetaTelemetryDestination, err error) {
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
	path := "v1/organizations/telemetry_destinations?beta=true"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.create", "/v1/organizations/telemetry_destinations"), requestconfig.WithSpanCreatedID("anthropic.organization.telemetry_destination.id", &res, func(created *BetaTelemetryDestination) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Read a destination's configuration, state, and most recent delivery results.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#delivery-and-health)
// to learn how the API delivers spans and when it disables a destination.
func (r *BetaOrganizationTelemetryDestinationService) Get(ctx context.Context, destinationID string, query BetaOrganizationTelemetryDestinationGetParams, opts ...option.RequestOption) (res *BetaTelemetryDestination, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "telemetry-destinations-2026-08-11")}, opts...)
	if destinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s?beta=true", url.PathEscape(destinationID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.retrieve", "/v1/organizations/telemetry_destinations/{destination_id}"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", destinationID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Change a telemetry destination's settings, its state, or which workspaces it
// covers.
//
// Send at least one field. Omitted fields keep their stored values, except the
// registered headers, which a change of the endpoint's hostname or port removes.
// Only `display_name`, `traces`, `traces.endpoint_url`, `traces.content_mode`,
// `traces.max_content_bytes`, and individual `headers` values accept `null`.
//
// If another request changes the destination while the update runs, the API
// returns a 409 `conflict_error`. Read the destination again, then retry the
// update.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#delivery-and-health)
// to learn more about changing a destination.
func (r *BetaOrganizationTelemetryDestinationService) Update(ctx context.Context, destinationID string, params BetaOrganizationTelemetryDestinationUpdateParams, opts ...option.RequestOption) (res *BetaTelemetryDestination, err error) {
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
	if destinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s?beta=true", url.PathEscape(destinationID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.update", "/v1/organizations/telemetry_destinations/{destination_id}"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", destinationID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the organization's destinations, newest first, optionally only those that
// cover given workspaces.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#set-up-api-tracing)
// to learn more about destinations.
func (r *BetaOrganizationTelemetryDestinationService) List(ctx context.Context, params BetaOrganizationTelemetryDestinationListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaTelemetryDestination], err error) {
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
	path := "v1/organizations/telemetry_destinations?beta=true"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.list", "/v1/organizations/telemetry_destinations", "limit", "page", "workspace_ids"))
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

// List the organization's destinations, newest first, optionally only those that
// cover given workspaces.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#set-up-api-tracing)
// to learn more about destinations.
func (r *BetaOrganizationTelemetryDestinationService) ListAutoPaging(ctx context.Context, params BetaOrganizationTelemetryDestinationListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaTelemetryDestination] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete a destination and its workspace list.
//
// The API stops sending spans to the destination. The spans of the workspaces it
// covered still go to the other destinations that cover them. A deleted
// destination can't be read, listed, or restored.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#delivery-and-health)
// to learn more about managing a destination.
func (r *BetaOrganizationTelemetryDestinationService) Delete(ctx context.Context, destinationID string, body BetaOrganizationTelemetryDestinationDeleteParams, opts ...option.RequestOption) (res *BetaDeletedTelemetryDestination, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "telemetry-destinations-2026-08-11")}, opts...)
	if destinationID == "" {
		err = errors.New("missing required destination_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/telemetry_destinations/%s?beta=true", url.PathEscape(destinationID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.telemetry_destinations.delete", "/v1/organizations/telemetry_destinations/{destination_id}"), requestconfig.WithSpanAttribute("anthropic.organization.telemetry_destination.id", destinationID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Confirmation that a destination was deleted.
type BetaDeletedTelemetryDestination struct {
	// ID of the deleted destination (`tdst_...`).
	ID   string                               `json:"id" api:"required"`
	Type constant.TelemetryDestinationDeleted `json:"type" default:"telemetry_destination_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedTelemetryDestination) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedTelemetryDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An OpenTelemetry destination: an OTLP/HTTP endpoint, the request headers sent
// with every delivery to it, what the exported spans contain, and which workspaces
// it covers. A destination belongs to an organization. A workspace's spans go to
// every enabled destination that covers it.
type BetaTelemetryDestination struct {
	// The unique ID of the destination (`tdst_...`).
	ID string `json:"id" api:"required"`
	// Time the destination was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The most recent delivery results for this destination.
	DeliveryHealth BetaTelemetryDestinationDeliveryHealth `json:"delivery_health" api:"required"`
	// Human-readable name for the destination. Null when none was set.
	DisplayName string `json:"display_name" api:"required"`
	// The base URL of the destination's OTLP/HTTP endpoint, exactly as registered.
	//
	// The API appends the per-signal path, such as `/v1/traces`, on delivery, unless
	// `traces.endpoint_url` is set.
	EndpointURL string `json:"endpoint_url" api:"required"`
	// Names of the registered request headers. Header values are write-only and are
	// never returned.
	HeaderNames []string `json:"header_names" api:"required"`
	// Whether the destination receives deliveries.
	State BetaTelemetryDestinationStateUnion `json:"state" api:"required"`
	// The destination's settings for the traces it receives.
	Traces BetaTelemetryDestinationTraces `json:"traces" api:"required"`
	Type   constant.TelemetryDestination  `json:"type" default:"telemetry_destination"`
	// Time of the last change to the destination's configuration or state. Changes to
	// `delivery_health` don't change it.
	//
	// An update that sets a header value always changes it, even when the value is the
	// same as before. Adding or removing a single workspace on the destination's list
	// doesn't change it. An update that changes the list does.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Which of the organization's workspaces the destination covers.
	Workspaces BetaTelemetryDestinationWorkspacesUnion `json:"workspaces" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		DeliveryHealth respjson.Field
		DisplayName    respjson.Field
		EndpointURL    respjson.Field
		HeaderNames    respjson.Field
		State          respjson.Field
		Traces         respjson.Field
		Type           respjson.Field
		UpdatedAt      respjson.Field
		Workspaces     respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestination) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestination) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What the spans that the API exports to a telemetry destination contain.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#choose-a-content-mode)
// to learn what the spans contain in each content mode.
type BetaTelemetryDestinationContentMode string

const (
	// The spans contain identifiers, timing, status, token counts, and similar
	// metadata. They don't contain prompt text, completion text, or tool inputs and
	// outputs.
	BetaTelemetryDestinationContentModeMetadataOnly BetaTelemetryDestinationContentMode = "metadata_only"
	// The spans contain the same metadata as `metadata_only`, plus message and tool
	// content.
	//
	// An organization with a data-handling restriction, such as under HIPAA, can't
	// choose this mode. The API exports only metadata for such an organization, even
	// to a destination that already has `content`.
	BetaTelemetryDestinationContentModeContent BetaTelemetryDestinationContentMode = "content"
)

// The most recent delivery results that the API records for a telemetry
// destination.
//
// Changes to the destination's configuration or state don't reset these fields.
// Changes to these fields don't change `updated_at`. The values can lag behind the
// latest deliveries.
//
// See the
// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#delivery-and-health)
// to learn how the API delivers spans and when it disables a destination.
type BetaTelemetryDestinationDeliveryHealth struct {
	// A message that describes the most recent failed delivery to the destination, or
	// `null` if no delivery has failed yet.
	//
	// Examples are `temporary DNS failure`, `TLS failure`, `timeout`,
	// `connection failure`, `invalid address`, `HTTP error (status 503)`, and
	// `redirect (status 301)`. The message is `invalid address` both for a hostname
	// that doesn't exist and for an address that isn't publicly routable. It never
	// contains the endpoint's response body, a header value, a network address, or TLS
	// details.
	//
	// `last_error` and `last_error_at` always change together. A later successful
	// delivery doesn't clear them. If `last_success_at` is later than `last_error_at`,
	// a delivery has succeeded since this failure.
	//
	// See the
	// [OpenTelemetry export guide](https://platform.claude.com/docs/en/manage-claude/opentelemetry#delivery-errors)
	// to learn the likely cause of each error.
	LastError string `json:"last_error" api:"required"`
	// Time of the failure that `last_error` describes. Changes together with
	// `last_error`. Null until a delivery first fails.
	LastErrorAt time.Time `json:"last_error_at" api:"required" format:"date-time"`
	// Time of the most recent successful delivery to the destination, or `null` if no
	// delivery has succeeded yet.
	LastSuccessAt time.Time `json:"last_success_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		LastError     respjson.Field
		LastErrorAt   respjson.Field
		LastSuccessAt respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationDeliveryHealth) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationDeliveryHealth) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaTelemetryDestinationDisabledReasonUnion contains all possible properties and
// values from [BetaTelemetryDestinationDisabledReasonCaller],
// [BetaTelemetryDestinationDisabledReasonSustainedFailure],
// [BetaTelemetryDestinationDisabledReasonInvalidAddress],
// [BetaTelemetryDestinationDisabledReasonRedirect].
//
// Use the [BetaTelemetryDestinationDisabledReasonUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaTelemetryDestinationDisabledReasonUnion struct {
	// Any of "caller", "sustained_failure", "invalid_address", "redirect".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaTelemetryDestinationDisabledReason is implemented by each variant of
// [BetaTelemetryDestinationDisabledReasonUnion] to add type safety for the return
// type of [BetaTelemetryDestinationDisabledReasonUnion.AsAny]
type anyBetaTelemetryDestinationDisabledReason interface {
	implBetaTelemetryDestinationDisabledReasonUnion()
}

func (BetaTelemetryDestinationDisabledReasonCaller) implBetaTelemetryDestinationDisabledReasonUnion() {
}
func (BetaTelemetryDestinationDisabledReasonSustainedFailure) implBetaTelemetryDestinationDisabledReasonUnion() {
}
func (BetaTelemetryDestinationDisabledReasonInvalidAddress) implBetaTelemetryDestinationDisabledReasonUnion() {
}
func (BetaTelemetryDestinationDisabledReasonRedirect) implBetaTelemetryDestinationDisabledReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaTelemetryDestinationDisabledReasonUnion.AsAny().(type) {
//	case anthropic.BetaTelemetryDestinationDisabledReasonCaller:
//	case anthropic.BetaTelemetryDestinationDisabledReasonSustainedFailure:
//	case anthropic.BetaTelemetryDestinationDisabledReasonInvalidAddress:
//	case anthropic.BetaTelemetryDestinationDisabledReasonRedirect:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaTelemetryDestinationDisabledReasonUnion) AsAny() anyBetaTelemetryDestinationDisabledReason {
	switch u.Type {
	case "caller":
		return u.AsCaller()
	case "sustained_failure":
		return u.AsSustainedFailure()
	case "invalid_address":
		return u.AsInvalidAddress()
	case "redirect":
		return u.AsRedirect()
	}
	return nil
}

func (u BetaTelemetryDestinationDisabledReasonUnion) AsCaller() (v BetaTelemetryDestinationDisabledReasonCaller) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationDisabledReasonUnion) AsSustainedFailure() (v BetaTelemetryDestinationDisabledReasonSustainedFailure) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationDisabledReasonUnion) AsInvalidAddress() (v BetaTelemetryDestinationDisabledReasonInvalidAddress) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationDisabledReasonUnion) AsRedirect() (v BetaTelemetryDestinationDisabledReasonRedirect) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaTelemetryDestinationDisabledReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaTelemetryDestinationDisabledReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A caller disabled the destination with an update that set `state` to `disabled`.
type BetaTelemetryDestinationDisabledReasonCaller struct {
	Type constant.Caller `json:"type" default:"caller"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationDisabledReasonCaller) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationDisabledReasonCaller) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API disabled the destination because deliveries to it failed for longer than
// the same threshold as for `sustained_failure`, and the most recent failed
// delivery failed with an invalid address: the DNS answer says that the endpoint's
// hostname doesn't exist or doesn't have an address, or the address that the API
// would connect to isn't publicly routable or is reserved, for example a loopback,
// link-local, private-network, or cloud-metadata address. The API doesn't connect
// to such an address, and retries the delivery a bounded number of times as it
// does after a timeout. A timeout or a server failure of the DNS lookup is a
// temporary DNS failure instead, not an invalid address.
type BetaTelemetryDestinationDisabledReasonInvalidAddress struct {
	Type constant.InvalidAddress `json:"type" default:"invalid_address"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationDisabledReasonInvalidAddress) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationDisabledReasonInvalidAddress) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API disabled the destination because the endpoint answered a delivery with a
// 3xx status. Redirects are never followed, and one such answer disables the
// destination at once.
type BetaTelemetryDestinationDisabledReasonRedirect struct {
	Type constant.Redirect `json:"type" default:"redirect"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationDisabledReasonRedirect) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationDisabledReasonRedirect) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API disabled the destination because deliveries to it failed, with no
// successful delivery in between, for longer than a threshold that Anthropic
// configures and that is not part of this contract. The reason is
// `invalid_address` instead when the most recent failed delivery failed with an
// invalid address.
type BetaTelemetryDestinationDisabledReasonSustainedFailure struct {
	Type constant.SustainedFailure `json:"type" default:"sustained_failure"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationDisabledReasonSustainedFailure) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationDisabledReasonSustainedFailure) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaTelemetryDestinationStateUnion contains all possible properties and values
// from [BetaTelemetryDestinationStateEnabled],
// [BetaTelemetryDestinationStateDisabled].
//
// Use the [BetaTelemetryDestinationStateUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaTelemetryDestinationStateUnion struct {
	// Any of "enabled", "disabled".
	Type string `json:"type"`
	// This field is from variant [BetaTelemetryDestinationStateDisabled].
	Reason BetaTelemetryDestinationDisabledReasonUnion `json:"reason"`
	JSON   struct {
		Type   respjson.Field
		Reason respjson.Field
		raw    string
	} `json:"-"`
}

// anyBetaTelemetryDestinationState is implemented by each variant of
// [BetaTelemetryDestinationStateUnion] to add type safety for the return type of
// [BetaTelemetryDestinationStateUnion.AsAny]
type anyBetaTelemetryDestinationState interface {
	implBetaTelemetryDestinationStateUnion()
}

func (BetaTelemetryDestinationStateEnabled) implBetaTelemetryDestinationStateUnion()  {}
func (BetaTelemetryDestinationStateDisabled) implBetaTelemetryDestinationStateUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaTelemetryDestinationStateUnion.AsAny().(type) {
//	case anthropic.BetaTelemetryDestinationStateEnabled:
//	case anthropic.BetaTelemetryDestinationStateDisabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaTelemetryDestinationStateUnion) AsAny() anyBetaTelemetryDestinationState {
	switch u.Type {
	case "enabled":
		return u.AsEnabled()
	case "disabled":
		return u.AsDisabled()
	}
	return nil
}

func (u BetaTelemetryDestinationStateUnion) AsEnabled() (v BetaTelemetryDestinationStateEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationStateUnion) AsDisabled() (v BetaTelemetryDestinationStateDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaTelemetryDestinationStateUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaTelemetryDestinationStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The API doesn't deliver anything to the telemetry destination.
//
// Spans produced while it is disabled aren't delivered later.
type BetaTelemetryDestinationStateDisabled struct {
	// Whether a caller or the API disabled the destination and, when the API did, why.
	// Read-only.
	Reason BetaTelemetryDestinationDisabledReasonUnion `json:"reason" api:"required"`
	Type   constant.Disabled                           `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationStateDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationStateDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaTelemetryDestinationStateDisabledParams() BetaTelemetryDestinationStateDisabledParams {
	return BetaTelemetryDestinationStateDisabledParams{
		Type: "disabled",
	}
}

// Stop deliveries to the telemetry destination.
//
// Spans produced while it is disabled aren't delivered later. A destination that
// is already disabled keeps its `reason`.
//
// This struct has a constant value, construct it with
// [NewBetaTelemetryDestinationStateDisabledParams].
type BetaTelemetryDestinationStateDisabledParams struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaTelemetryDestinationStateDisabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationStateDisabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationStateDisabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The telemetry destination receives the spans of the workspaces that it covers.
type BetaTelemetryDestinationStateEnabled struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationStateEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationStateEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaTelemetryDestinationStateEnabledParams() BetaTelemetryDestinationStateEnabledParams {
	return BetaTelemetryDestinationStateEnabledParams{
		Type: "enabled",
	}
}

// Start deliveries to the telemetry destination again.
//
// This doesn't reset `delivery_health`.
//
// This struct has a constant value, construct it with
// [NewBetaTelemetryDestinationStateEnabledParams].
type BetaTelemetryDestinationStateEnabledParams struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaTelemetryDestinationStateEnabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationStateEnabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationStateEnabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaTelemetryDestinationStateParamsUnion struct {
	OfEnabled  *BetaTelemetryDestinationStateEnabledParams  `json:",omitzero,inline"`
	OfDisabled *BetaTelemetryDestinationStateDisabledParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaTelemetryDestinationStateParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfEnabled, u.OfDisabled)
}
func (u *BetaTelemetryDestinationStateParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaTelemetryDestinationStateParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	} else if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTelemetryDestinationStateParamsUnion) GetType() *string {
	if vt := u.OfEnabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfDisabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaTelemetryDestinationStateParamsUnion](
		"type",
		apijson.Discriminator[BetaTelemetryDestinationStateEnabledParams]("enabled"),
		apijson.Discriminator[BetaTelemetryDestinationStateDisabledParams]("disabled"),
	)
}

// A telemetry destination's settings for the traces it receives.
type BetaTelemetryDestinationTraces struct {
	// What the exported spans contain.
	//
	// Any of "metadata_only", "content".
	ContentMode BetaTelemetryDestinationContentMode `json:"content_mode" api:"required"`
	// Whether the destination receives traces. A destination whose `enabled` is false
	// keeps its other settings.
	Enabled bool `json:"enabled" api:"required"`
	// The full URL that the API posts traces to, exactly as registered, or `null` if
	// the destination doesn't have one.
	//
	// When it is `null`, the API posts traces to `endpoint_url` with `/v1/traces`
	// appended.
	EndpointURL string `json:"endpoint_url" api:"required"`
	// The most bytes of content that one exported span can contain, or `null` if only
	// the API's maximum applies.
	//
	// Content means model prompts and responses, and tool inputs and outputs. The API
	// shortens longer content to this size before it delivers the span. A destination
	// whose `content_mode` is `metadata_only` keeps this value. It applies once
	// `content_mode` is `content`.
	MaxContentBytes int64 `json:"max_content_bytes" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContentMode     respjson.Field
		Enabled         respjson.Field
		EndpointURL     respjson.Field
		MaxContentBytes respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationTraces) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationTraces) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The settings for the traces that a telemetry destination receives, as a create
// or update request sets them. An omitted field takes its default on create and
// keeps its stored value on update.
type BetaTelemetryDestinationTracesParams struct {
	// The full URL of your collector's traces endpoint, for a collector that takes
	// traces on a path other than `/v1/traces`. The API posts traces to it exactly as
	// given and appends nothing.
	//
	// The same rules as for `endpoint_url` apply, except that any path is allowed. It
	// must have the same scheme, host, and port as `endpoint_url`. At most 2,048
	// characters.
	//
	// On create, `null` is the same as leaving it out. On update, `null` removes it.
	// The API then posts traces to `endpoint_url` with `/v1/traces` appended. A change
	// of this URL alone keeps the registered headers.
	EndpointURL param.Opt[string] `json:"endpoint_url,omitzero"`
	// The most bytes of content that one exported span can contain, from 1 up to the
	// API's maximum. `null` sets no bound, so only the API's maximum applies.
	//
	// On create, leaving it out is the same as `null`. The error for a value above the
	// maximum states the maximum.
	MaxContentBytes param.Opt[int64] `json:"max_content_bytes,omitzero"`
	// Whether the destination receives traces. Defaults to true on create. It cannot
	// be null.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// What the exported spans contain. Defaults to `metadata_only` on create.
	//
	// On update, `null` sets `metadata_only`.
	//
	// Any of "metadata_only", "content".
	ContentMode BetaTelemetryDestinationContentMode `json:"content_mode,omitzero"`
	paramObj
}

func (r BetaTelemetryDestinationTracesParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationTracesParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationTracesParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaTelemetryDestinationWorkspacesUnion contains all possible properties and
// values from [BetaTelemetryDestinationWorkspacesAll],
// [BetaTelemetryDestinationWorkspacesInclude],
// [BetaTelemetryDestinationWorkspacesExclude].
//
// Use the [BetaTelemetryDestinationWorkspacesUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaTelemetryDestinationWorkspacesUnion struct {
	// Any of "all", "include", "exclude".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaTelemetryDestinationWorkspaces is implemented by each variant of
// [BetaTelemetryDestinationWorkspacesUnion] to add type safety for the return type
// of [BetaTelemetryDestinationWorkspacesUnion.AsAny]
type anyBetaTelemetryDestinationWorkspaces interface {
	implBetaTelemetryDestinationWorkspacesUnion()
}

func (BetaTelemetryDestinationWorkspacesAll) implBetaTelemetryDestinationWorkspacesUnion()     {}
func (BetaTelemetryDestinationWorkspacesInclude) implBetaTelemetryDestinationWorkspacesUnion() {}
func (BetaTelemetryDestinationWorkspacesExclude) implBetaTelemetryDestinationWorkspacesUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaTelemetryDestinationWorkspacesUnion.AsAny().(type) {
//	case anthropic.BetaTelemetryDestinationWorkspacesAll:
//	case anthropic.BetaTelemetryDestinationWorkspacesInclude:
//	case anthropic.BetaTelemetryDestinationWorkspacesExclude:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaTelemetryDestinationWorkspacesUnion) AsAny() anyBetaTelemetryDestinationWorkspaces {
	switch u.Type {
	case "all":
		return u.AsAll()
	case "include":
		return u.AsInclude()
	case "exclude":
		return u.AsExclude()
	}
	return nil
}

func (u BetaTelemetryDestinationWorkspacesUnion) AsAll() (v BetaTelemetryDestinationWorkspacesAll) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationWorkspacesUnion) AsInclude() (v BetaTelemetryDestinationWorkspacesInclude) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTelemetryDestinationWorkspacesUnion) AsExclude() (v BetaTelemetryDestinationWorkspacesExclude) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaTelemetryDestinationWorkspacesUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaTelemetryDestinationWorkspacesUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The destination covers every workspace of the organization, including the
// default workspace, archived workspaces, and workspaces created later. It doesn't
// have a list.
type BetaTelemetryDestinationWorkspacesAll struct {
	Type constant.All `json:"type" default:"all"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationWorkspacesAll) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationWorkspacesAll) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaTelemetryDestinationWorkspacesAllParams() BetaTelemetryDestinationWorkspacesAllParams {
	return BetaTelemetryDestinationWorkspacesAllParams{
		Type: "all",
	}
}

// Cover every workspace of the organization. Doesn't take `workspace_ids`.
//
// An update that sets `all` empties the destination's workspace list.
//
// This struct has a constant value, construct it with
// [NewBetaTelemetryDestinationWorkspacesAllParams].
type BetaTelemetryDestinationWorkspacesAllParams struct {
	Type constant.All `json:"type" default:"all"`
	paramObj
}

func (r BetaTelemetryDestinationWorkspacesAllParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationWorkspacesAllParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationWorkspacesAllParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The destination covers every workspace of the organization except the workspaces
// on its list. A workspace created later is covered until it is added to the list.
type BetaTelemetryDestinationWorkspacesExclude struct {
	Type constant.Exclude `json:"type" default:"exclude"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationWorkspacesExclude) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationWorkspacesExclude) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cover every workspace of the organization except the workspaces on the list.
//
// The property Type is required.
type BetaTelemetryDestinationWorkspacesExcludeParams struct {
	// The IDs of the workspaces on the destination's list (`wrkspc_...`). At most 100
	// per request.
	//
	// On create, these workspaces become the destination's list. If you leave
	// `workspace_ids` out, the list starts empty. If the API refuses one of the IDs,
	// it doesn't create the destination.
	//
	// On update, `workspace_ids` replaces the whole list. A workspace that stays on
	// the list keeps its entry and its `created_at`. An empty array empties the list.
	// If you leave `workspace_ids` out and the type stays the same, the list doesn't
	// change. If you leave it out and the type changes, the new list starts empty.
	// After a change to `include` without `workspace_ids`, the destination doesn't
	// cover any workspace until one is added. After a change to `exclude` without
	// `workspace_ids`, it covers every workspace.
	//
	// An ID sent more than once is added once. An archived workspace can't be added to
	// a list. An update can keep an archived workspace that is already on the list.
	WorkspaceIDs []string `json:"workspace_ids,omitzero"`
	// This field can be elided, and will marshal its zero value as "exclude".
	Type constant.Exclude `json:"type" default:"exclude"`
	paramObj
}

func (r BetaTelemetryDestinationWorkspacesExcludeParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationWorkspacesExcludeParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationWorkspacesExcludeParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The destination covers the workspaces on its list, and no others.
type BetaTelemetryDestinationWorkspacesInclude struct {
	Type constant.Include `json:"type" default:"include"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTelemetryDestinationWorkspacesInclude) RawJSON() string { return r.JSON.raw }
func (r *BetaTelemetryDestinationWorkspacesInclude) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cover only the workspaces on the list.
//
// The property Type is required.
type BetaTelemetryDestinationWorkspacesIncludeParams struct {
	// The IDs of the workspaces on the destination's list (`wrkspc_...`). At most 100
	// per request.
	//
	// On create, these workspaces become the destination's list. If you leave
	// `workspace_ids` out, the list starts empty. If the API refuses one of the IDs,
	// it doesn't create the destination.
	//
	// On update, `workspace_ids` replaces the whole list. A workspace that stays on
	// the list keeps its entry and its `created_at`. An empty array empties the list.
	// If you leave `workspace_ids` out and the type stays the same, the list doesn't
	// change. If you leave it out and the type changes, the new list starts empty.
	// After a change to `include` without `workspace_ids`, the destination doesn't
	// cover any workspace until one is added. After a change to `exclude` without
	// `workspace_ids`, it covers every workspace.
	//
	// An ID sent more than once is added once. An archived workspace can't be added to
	// a list. An update can keep an archived workspace that is already on the list.
	WorkspaceIDs []string `json:"workspace_ids,omitzero"`
	// This field can be elided, and will marshal its zero value as "include".
	Type constant.Include `json:"type" default:"include"`
	paramObj
}

func (r BetaTelemetryDestinationWorkspacesIncludeParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTelemetryDestinationWorkspacesIncludeParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTelemetryDestinationWorkspacesIncludeParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaTelemetryDestinationWorkspacesParamsUnion struct {
	OfAll     *BetaTelemetryDestinationWorkspacesAllParams     `json:",omitzero,inline"`
	OfInclude *BetaTelemetryDestinationWorkspacesIncludeParams `json:",omitzero,inline"`
	OfExclude *BetaTelemetryDestinationWorkspacesExcludeParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaTelemetryDestinationWorkspacesParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAll, u.OfInclude, u.OfExclude)
}
func (u *BetaTelemetryDestinationWorkspacesParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaTelemetryDestinationWorkspacesParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfAll) {
		return u.OfAll
	} else if !param.IsOmitted(u.OfInclude) {
		return u.OfInclude
	} else if !param.IsOmitted(u.OfExclude) {
		return u.OfExclude
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTelemetryDestinationWorkspacesParamsUnion) GetType() *string {
	if vt := u.OfAll; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfInclude; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfExclude; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's WorkspaceIDs property, if present.
func (u BetaTelemetryDestinationWorkspacesParamsUnion) GetWorkspaceIDs() []string {
	if vt := u.OfInclude; vt != nil {
		return vt.WorkspaceIDs
	} else if vt := u.OfExclude; vt != nil {
		return vt.WorkspaceIDs
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaTelemetryDestinationWorkspacesParamsUnion](
		"type",
		apijson.Discriminator[BetaTelemetryDestinationWorkspacesAllParams]("all"),
		apijson.Discriminator[BetaTelemetryDestinationWorkspacesIncludeParams]("include"),
		apijson.Discriminator[BetaTelemetryDestinationWorkspacesExcludeParams]("exclude"),
	)
}

type BetaOrganizationTelemetryDestinationNewParams struct {
	// The base URL of your collector's OTLP/HTTP endpoint, for example
	// `https://gateway.example.com`.
	//
	// The URL must use `https` and a hostname, not an IP address. The port, if given,
	// must be 443. The URL must not contain a query string, a fragment, or a username
	// and password. At most 2,048 characters.
	//
	// A path prefix, such as `/otlp`, is allowed. Don't include the per-signal path
	// (`/v1/traces`), because the API appends it on delivery. If your collector takes
	// traces on another path, also set `traces.endpoint_url`.
	EndpointURL string `json:"endpoint_url" api:"required"`
	// Which of the organization's workspaces the destination covers.
	Workspaces BetaTelemetryDestinationWorkspacesParamsUnion `json:"workspaces,omitzero" api:"required"`
	// Human-readable name for the destination. Optional. 1 to 255 characters.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Request headers that the API sends with every delivery, typically your
	// collector's authentication header (for example `authorization`).
	//
	// The API stores the values encrypted and never returns them. Reads show only
	// `header_names`.
	//
	// At most 20 headers. Names are at most 128 characters. Two names can't differ
	// only in case. Values are printable ASCII, at most 4,096 characters, and can't
	// begin or end with a space. Names and values together are at most 16,384 bytes.
	//
	// Reserved names are rejected, for example `host`, `cookie`, `user-agent`, and
	// names that start with `content-`, `proxy-`, `x-forwarded-`, or `anthropic-`. For
	// this check, `_` in a name counts as `-`.
	Headers map[string]string `json:"headers,omitzero"`
	// The settings for the traces that the destination receives. If you leave it out,
	// every setting takes its default.
	Traces BetaTelemetryDestinationTracesParams `json:"traces,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationTelemetryDestinationNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationTelemetryDestinationNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationTelemetryDestinationNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationTelemetryDestinationGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationTelemetryDestinationUpdateParams struct {
	// New human-readable name, 1 to 255 characters. Null clears it.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// The new base URL of the destination's OTLP/HTTP endpoint. The same rules as on
	// create apply.
	//
	// A change of hostname or port first removes every registered header, so send the
	// headers for the new endpoint in the same request. When the destination has a
	// `traces.endpoint_url`, such a change must also set `traces.endpoint_url` to a
	// URL with the new hostname and port, or to `null`. A change of the path alone
	// keeps the headers.
	EndpointURL param.Opt[string] `json:"endpoint_url,omitzero"`
	// Headers to set or remove. The API merges them into the destination's headers by
	// name.
	//
	// Names compare case-insensitively. A string value sets or replaces that header. A
	// `null` value removes it. Headers that you omit don't change. The field itself
	// can't be `null`.
	//
	// A request that sets a header value must also include `endpoint_url`. The current
	// URL is accepted.
	//
	// The resulting set must stay within the same limits as on create. When the
	// request sets a value, each stored value that it doesn't set counts as 4,096
	// bytes toward the 16,384-byte limit, because you can't read its real length. So
	// to change one header of a destination that has four or more other headers, send
	// the values of the others again in the same request.
	Headers map[string]string `json:"headers,omitzero"`
	// Whether the destination receives deliveries after the update. The field can't be
	// `null`.
	State BetaTelemetryDestinationStateParamsUnion `json:"state,omitzero"`
	// The traces settings to change. Omitted settings keep their stored values.
	//
	// `null` sets every traces setting to its create default.
	Traces BetaTelemetryDestinationTracesParams `json:"traces,omitzero"`
	// The new workspace scope, in the same form as on create. The field can't be
	// `null`.
	//
	// The API refuses an update that changes the type or sends `workspace_ids` while
	// the destination's list has more than 750 workspaces. Remove workspaces one at a
	// time until 750 or fewer remain, then send the update again.
	Workspaces BetaTelemetryDestinationWorkspacesParamsUnion `json:"workspaces,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationTelemetryDestinationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationTelemetryDestinationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationTelemetryDestinationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationTelemetryDestinationListParams struct {
	// Number of destinations to return per page, between 1 and 1,000. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`. Leave it out to get the
	// first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only the destinations that cover at least one of these workspaces
	// (`wrkspc_...`).
	//
	// Send at most 100 IDs, as `workspace_ids[]`. A destination covers a workspace
	// when its `workspaces.type` is `all`, when its `include` list has the workspace,
	// or when its `exclude` list doesn't have it. Destinations in either state can
	// match. An ID that isn't the ID of one of the organization's workspaces doesn't
	// match any destination, not even one whose type is `all`.
	WorkspaceIDs []string `query:"workspace_ids,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationTelemetryDestinationListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationTelemetryDestinationListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationTelemetryDestinationDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

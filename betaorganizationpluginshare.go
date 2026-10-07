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
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationPluginShareService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationPluginShareService] method instead.
type BetaOrganizationPluginShareService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationPluginShareService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationPluginShareService(opts ...option.RequestOption) (r BetaOrganizationPluginShareService) {
	r = BetaOrganizationPluginShareService{}
	r.Options = opts
	return
}

// List the shares the owner of a member-owned Plugin has given — to every member
// of the organization, to an RBAC Group, or to one member — most recently granted
// first.
//
// Shares are read-only in this API: members give and withdraw them in claude.ai,
// and who gave a share is recorded on the Compliance API activity feed rather than
// on the share. An organization-owned Plugin has installation settings instead, so
// this path returns 404 for one.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `read:org_audit` scope, or a Compliance Access Key with the
// `read:compliance_org_data` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginShareService) List(ctx context.Context, pluginID string, params BetaOrganizationPluginShareListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaPluginShare], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01"), option.WithResponseInto(&raw)}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/shares?beta=true", url.PathEscape(pluginID))
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

// List the shares the owner of a member-owned Plugin has given — to every member
// of the organization, to an RBAC Group, or to one member — most recently granted
// first.
//
// Shares are read-only in this API: members give and withdraw them in claude.ai,
// and who gave a share is recorded on the Compliance API activity feed rather than
// on the share. An organization-owned Plugin has installation settings instead, so
// this path returns 404 for one.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `read:org_audit` scope, or a Compliance Access Key with the
// `read:compliance_org_data` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginShareService) ListAutoPaging(ctx context.Context, pluginID string, params BetaOrganizationPluginShareListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaPluginShare] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, pluginID, params, opts...))
}

// One share the owner of a member-owned Plugin has given. Shares are read-only in
// this API and have no ID of their own; who gave a share is recorded on the
// Compliance API activity feed, not here.
type BetaPluginShare struct {
	// When the share was given; a share whose role is later changed in claude.ai is
	// re-granted and carries the time of that change.
	GrantedAt time.Time `json:"granted_at" api:"required" format:"date-time"`
	// The Plugin's ID.
	PluginID string `json:"plugin_id" api:"required"`
	// Who the Plugin is shared with: `organization` (every member), `rbac_group` (one
	// RBAC Group), or `organization_member` (one member).
	Target BetaPluginShareTargetUnion `json:"target" api:"required"`
	// Always `plugin_share`.
	Type constant.PluginShare `json:"type" default:"plugin_share"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GrantedAt   respjson.Field
		PluginID    respjson.Field
		Target      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginShare) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginShare) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaPluginShareTargetUnion contains all possible properties and values from
// [BetaPluginTargetOrganization], [BetaPluginTargetRBACGroup],
// [BetaPluginTargetOrganizationMember].
//
// Use the [BetaPluginShareTargetUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginShareTargetUnion struct {
	// Any of "organization", "rbac_group", "organization_member".
	Type string `json:"type"`
	// This field is from variant [BetaPluginTargetRBACGroup].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaPluginTargetOrganizationMember].
	UserID string `json:"user_id"`
	JSON   struct {
		Type        respjson.Field
		RBACGroupID respjson.Field
		UserID      respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaPluginShareTarget is implemented by each variant of
// [BetaPluginShareTargetUnion] to add type safety for the return type of
// [BetaPluginShareTargetUnion.AsAny]
type anyBetaPluginShareTarget interface {
	implBetaPluginShareTargetUnion()
}

func (BetaPluginTargetOrganization) implBetaPluginShareTargetUnion()       {}
func (BetaPluginTargetRBACGroup) implBetaPluginShareTargetUnion()          {}
func (BetaPluginTargetOrganizationMember) implBetaPluginShareTargetUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginShareTargetUnion.AsAny().(type) {
//	case anthropic.BetaPluginTargetOrganization:
//	case anthropic.BetaPluginTargetRBACGroup:
//	case anthropic.BetaPluginTargetOrganizationMember:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginShareTargetUnion) AsAny() anyBetaPluginShareTarget {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_member":
		return u.AsOrganizationMember()
	}
	return nil
}

func (u BetaPluginShareTargetUnion) AsOrganization() (v BetaPluginTargetOrganization) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginShareTargetUnion) AsRBACGroup() (v BetaPluginTargetRBACGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginShareTargetUnion) AsOrganizationMember() (v BetaPluginTargetOrganizationMember) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginShareTargetUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginShareTargetUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationPluginShareListParams struct {
	// For a `read:org_audit` or `read:compliance_org_data` key created for all of a
	// parent organization's linked organizations: a child organization of that parent
	// to read instead of the organization the key was created in, given as the
	// organization's UUID or its `org_`-prefixed ID. A value that is neither returns a
	// 400; an organization that is not a child of the key's parent, or where the
	// Plugins API is not available, returns a 404. Any other key may pass only its own
	// organization's ID here; another organization returns a 404.
	OrganizationID param.Opt[string] `query:"organization_id,omitzero" json:"-"`
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `100`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Only shares with this kind of target: `organization` (every member),
	// `rbac_group` (one RBAC Group), or `organization_member` (one member).
	//
	// Any of "organization", "organization_member", "rbac_group".
	TargetType BetaOrganizationPluginShareListParamsTargetType `query:"target_type,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginShareListParams]'s query parameters
// as `url.Values`.
func (r BetaOrganizationPluginShareListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only shares with this kind of target: `organization` (every member),
// `rbac_group` (one RBAC Group), or `organization_member` (one member).
type BetaOrganizationPluginShareListParamsTargetType string

const (
	BetaOrganizationPluginShareListParamsTargetTypeOrganization       BetaOrganizationPluginShareListParamsTargetType = "organization"
	BetaOrganizationPluginShareListParamsTargetTypeOrganizationMember BetaOrganizationPluginShareListParamsTargetType = "organization_member"
	BetaOrganizationPluginShareListParamsTargetTypeRBACGroup          BetaOrganizationPluginShareListParamsTargetType = "rbac_group"
)

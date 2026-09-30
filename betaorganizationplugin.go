package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
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

// BetaOrganizationPluginService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationPluginService] method instead.
type BetaOrganizationPluginService struct {
	Options              []option.RequestOption
	Versions             BetaOrganizationPluginVersionService
	InstallationSettings BetaOrganizationPluginInstallationSettingService
	Shares               BetaOrganizationPluginShareService
}

// NewBetaOrganizationPluginService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaOrganizationPluginService(opts ...option.RequestOption) (r BetaOrganizationPluginService) {
	r = BetaOrganizationPluginService{}
	r.Options = opts
	r.Versions = NewBetaOrganizationPluginVersionService(opts...)
	r.InstallationSettings = NewBetaOrganizationPluginInstallationSettingService(opts...)
	r.Shares = NewBetaOrganizationPluginShareService(opts...)
	return
}

// Create an organization-owned Plugin and its first version by uploading the
// version's files.
//
// The upload is `multipart/form-data`: the version's files (`files`, each part
// sent as `files[]`), with an optional `marketplace_id` and `release_notes`. The
// manifest's `name` becomes the Plugin's `name`, and `display_name`, `description`
// and `manifest_version` come from the manifest too.
//
// `name` may contain lowercase letters (from any alphabet), digits, and hyphens,
// up to 64 characters. Uppercase letters, spaces, underscores, and other
// punctuation are rejected.
//
// The `name` must be unique within the marketplace: a name already taken returns a
// 409 with `error_code` `plugin_name_taken` and, when a Plugin holds it, that
// Plugin's ID in `details.plugin_id`. A Plugin going into the organization's
// library marketplace is also refused with a 409 when one of its skills has the
// name of an organization skill (a skill an administrator uploaded for the whole
// organization in claude.ai): `error_code` `skill_name_taken`, with that name in
// `details.skill_name`; rename the skill, or remove the organization skill in
// claude.ai. A 503 with `error_code` `registration_pending` means the Plugin and
// its version were stored (their IDs are in `details`) but are not yet usable in
// claude.ai: do not retry the create (the retry would return `plugin_name_taken`);
// create a version on the stored Plugin instead, which completes it.
//
// For a worked example, see
// [Create a plugin](/docs/en/manage-claude/plugins-api#create-a-plugin) in the
// Plugins API guide.
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginService) New(ctx context.Context, params BetaOrganizationPluginNewParams, opts ...option.RequestOption) (res *BetaPlugin, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	path := "v1/organizations/plugins?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a Plugin by ID.
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
func (r *BetaOrganizationPluginService) Get(ctx context.Context, pluginID string, params BetaOrganizationPluginGetParams, opts ...option.RequestOption) (res *BetaPlugin, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s?beta=true", pluginID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Change which stored version of an organization-owned Plugin is served to
// members, for example to roll back to an earlier one. This pins the served
// version: later uploads are stored but no longer change what is served, and
// pinning cannot currently be undone, here or in claude.ai.
//
// Pass the version as `served_version_id`: an earlier one to roll back, a later
// one to start serving a version that was stored without being served, or the one
// already served to pin it without changing what is served. No new version is
// created.
//
// When the organization has content scanning enabled, a version whose scan is
// still running is refused with a 409 (`error_code` `scan_pending`; retry once the
// scan finishes) and one whose scan failed, errored or reached no verdict with a
// 400 (`scan_failed`; a `warn` is accepted). When the Plugin is in the
// organization's library marketplace, a version other than the one served is also
// refused with a 409 when one of its skills has a name that an organization skill
// (one an administrator uploaded for the whole organization in claude.ai) has
// since taken: `error_code` `skill_name_taken`, with that name in
// `details.skill_name`. A member-owned Plugin cannot be updated here (403).
//
// This endpoint does not write installation settings; they are written at
// `/v1/organizations/plugins/{plugin_id}/installation_settings/{target}`.
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginService) Update(ctx context.Context, pluginID string, params BetaOrganizationPluginUpdateParams, opts ...option.RequestOption) (res *BetaPlugin, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s?beta=true", pluginID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the Plugins created under the organization, newest first: those in the
// organization's own plugin marketplaces and those in members' personal plugin
// marketplaces.
//
// Plugins in members' personal marketplaces are listed with the same detail as the
// organization's own, and their files can be downloaded through the version
// archive endpoint, which records each such download on the Compliance API
// activity feed.
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
func (r *BetaOrganizationPluginService) List(ctx context.Context, params BetaOrganizationPluginListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaPlugin], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/plugins?beta=true"
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

// List the Plugins created under the organization, newest first: those in the
// organization's own plugin marketplaces and those in members' personal plugin
// marketplaces.
//
// Plugins in members' personal marketplaces are listed with the same detail as the
// organization's own, and their files can be downloaded through the version
// archive endpoint, which records each such download on the Compliance API
// activity feed.
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
func (r *BetaOrganizationPluginService) ListAutoPaging(ctx context.Context, params BetaOrganizationPluginListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaPlugin] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Permanently delete a Plugin and every version it holds, exactly as when an
// administrator deletes it in claude.ai. The Plugin may belong to the organization
// or to a member, including a member who has since left the organization.
//
// An organization-owned Plugin's installation settings go with it; a member-owned
// Plugin's shares are withdrawn and its owner no longer has it.
//
// To take an organization-owned Plugin out of use reversibly, set its
// organization-wide installation setting to `not_available` instead (and remove or
// change any group settings, which override it for their members). Only a Plugin
// in a `manual` marketplace can be deleted here; one synchronized from a
// repository is removed by removing it from the repository (400).
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginService) Delete(ctx context.Context, pluginID string, body BetaOrganizationPluginDeleteParams, opts ...option.RequestOption) (res *BetaDeletedPlugin, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s?beta=true", pluginID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaDeletedPlugin struct {
	// The deleted Plugin's ID.
	ID string `json:"id" api:"required"`
	// Always `plugin_deleted`.
	Type constant.PluginDeleted `json:"type" default:"plugin_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedPlugin) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedPlugin) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPlugin struct {
	// The Plugin's ID.
	ID string `json:"id" api:"required"`
	// What the served version contains; null when not enumerated.
	Components []BetaPluginComponent `json:"components" api:"required"`
	// The served version's content scan; null when it has not been scanned.
	ContentScan BetaPluginContentScan `json:"content_scan" api:"required"`
	// RFC 3339.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Who created the Plugin; null when no creator is recorded.
	CreatedBy BetaPluginCreatedByUnion `json:"created_by" api:"required"`
	// The served version's description.
	Description string `json:"description" api:"required"`
	// The served version's display name.
	DisplayName string `json:"display_name" api:"required"`
	// The newest version.
	LatestVersionID string `json:"latest_version_id" api:"required"`
	// The version string the served version's manifest declares.
	ManifestVersion string `json:"manifest_version" api:"required"`
	// The ID of the plugin marketplace the Plugin lives in.
	MarketplaceID string `json:"marketplace_id" api:"required"`
	// Lowercase identifier, unique within its plugin marketplace. Fixed for an
	// organization-owned Plugin's lifetime; a member-owned Plugin's changes when its
	// owner renames it in claude.ai, while its `id` stays the same.
	Name string `json:"name" api:"required"`
	// Organization-owned Plugin: the organization-wide installation setting every
	// member gets unless an RBAC Group they belong to holds its own — the Plugin's own
	// setting, or its plugin marketplace's default. Null for a member-owned Plugin,
	// which has shares instead. One of `required`, `auto_install`, `available`,
	// `not_available`; a value this API does not yet name is returned as stored.
	//
	// Any of "auto_install", "available", "not_available", "required".
	OrganizationInstallationPreference BetaPluginOrganizationInstallationPreference `json:"organization_installation_preference" api:"required"`
	// Organization-owned Plugin: true while it has no organization-wide setting of its
	// own and `organization_installation_preference` is its plugin marketplace's
	// default. Null for a member-owned Plugin.
	OrganizationInstallationPreferenceInherited bool `json:"organization_installation_preference_inherited" api:"required"`
	// Who owns the Plugin: the organization, or the member whose personal plugin
	// marketplace it lives in.
	Owner BetaPluginOwnerUnion `json:"owner" api:"required"`
	// How far the served version reaches: `remote` when it declares an MCP server or a
	// CLI, `privileged` when it declares a hook, monitor, language server or settings
	// but nothing remote, `contained` otherwise; null when not classifiable.
	//
	// Any of "contained", "privileged", "remote".
	Reach BetaPluginReach `json:"reach" api:"required"`
	// The version claude.ai serves to members.
	ServedVersionID string `json:"served_version_id" api:"required"`
	// False while the served version follows each new version; true once it has been
	// pinned to one.
	ServedVersionPinned bool `json:"served_version_pinned" api:"required"`
	// Always `plugin`.
	Type constant.Plugin `json:"type" default:"plugin"`
	// RFC 3339. Moves on a new version and on a served-version change; a change to the
	// Plugin's installation settings or shares does not move it.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                                          respjson.Field
		Components                                  respjson.Field
		ContentScan                                 respjson.Field
		CreatedAt                                   respjson.Field
		CreatedBy                                   respjson.Field
		Description                                 respjson.Field
		DisplayName                                 respjson.Field
		LatestVersionID                             respjson.Field
		ManifestVersion                             respjson.Field
		MarketplaceID                               respjson.Field
		Name                                        respjson.Field
		OrganizationInstallationPreference          respjson.Field
		OrganizationInstallationPreferenceInherited respjson.Field
		Owner                                       respjson.Field
		Reach                                       respjson.Field
		ServedVersionID                             respjson.Field
		ServedVersionPinned                         respjson.Field
		Type                                        respjson.Field
		UpdatedAt                                   respjson.Field
		ExtraFields                                 map[string]respjson.Field
		raw                                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPlugin) RawJSON() string { return r.JSON.raw }
func (r *BetaPlugin) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaPluginCreatedByUnion contains all possible properties and values from
// [BetaPluginUserActor], [BetaPluginAPIActor].
//
// Use the [BetaPluginCreatedByUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginCreatedByUnion struct {
	// This field is from variant [BetaPluginUserActor].
	EmailAddress string `json:"email_address"`
	// Any of "user_actor", "api_actor".
	Type string `json:"type"`
	// This field is from variant [BetaPluginUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaPluginAPIActor].
	APIKeyID string `json:"api_key_id"`
	JSON     struct {
		EmailAddress respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		APIKeyID     respjson.Field
		raw          string
	} `json:"-"`
}

// anyBetaPluginCreatedBy is implemented by each variant of
// [BetaPluginCreatedByUnion] to add type safety for the return type of
// [BetaPluginCreatedByUnion.AsAny]
type anyBetaPluginCreatedBy interface {
	implBetaPluginCreatedByUnion()
}

func (BetaPluginUserActor) implBetaPluginCreatedByUnion() {}
func (BetaPluginAPIActor) implBetaPluginCreatedByUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginCreatedByUnion.AsAny().(type) {
//	case anthropic.BetaPluginUserActor:
//	case anthropic.BetaPluginAPIActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginCreatedByUnion) AsAny() anyBetaPluginCreatedBy {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "api_actor":
		return u.AsAPIActor()
	}
	return nil
}

func (u BetaPluginCreatedByUnion) AsUserActor() (v BetaPluginUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginCreatedByUnion) AsAPIActor() (v BetaPluginAPIActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginCreatedByUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginCreatedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Organization-owned Plugin: the organization-wide installation setting every
// member gets unless an RBAC Group they belong to holds its own — the Plugin's own
// setting, or its plugin marketplace's default. Null for a member-owned Plugin,
// which has shares instead. One of `required`, `auto_install`, `available`,
// `not_available`; a value this API does not yet name is returned as stored.
type BetaPluginOrganizationInstallationPreference string

const (
	BetaPluginOrganizationInstallationPreferenceAutoInstall  BetaPluginOrganizationInstallationPreference = "auto_install"
	BetaPluginOrganizationInstallationPreferenceAvailable    BetaPluginOrganizationInstallationPreference = "available"
	BetaPluginOrganizationInstallationPreferenceNotAvailable BetaPluginOrganizationInstallationPreference = "not_available"
	BetaPluginOrganizationInstallationPreferenceRequired     BetaPluginOrganizationInstallationPreference = "required"
)

// BetaPluginOwnerUnion contains all possible properties and values from
// [BetaPluginOwnerOrganization], [BetaPluginOwnerUser].
//
// Use the [BetaPluginOwnerUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginOwnerUnion struct {
	// Any of "organization", "user".
	Type string `json:"type"`
	// This field is from variant [BetaPluginOwnerUser].
	UserID string `json:"user_id"`
	JSON   struct {
		Type   respjson.Field
		UserID respjson.Field
		raw    string
	} `json:"-"`
}

// anyBetaPluginOwner is implemented by each variant of [BetaPluginOwnerUnion] to
// add type safety for the return type of [BetaPluginOwnerUnion.AsAny]
type anyBetaPluginOwner interface {
	implBetaPluginOwnerUnion()
}

func (BetaPluginOwnerOrganization) implBetaPluginOwnerUnion() {}
func (BetaPluginOwnerUser) implBetaPluginOwnerUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginOwnerUnion.AsAny().(type) {
//	case anthropic.BetaPluginOwnerOrganization:
//	case anthropic.BetaPluginOwnerUser:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginOwnerUnion) AsAny() anyBetaPluginOwner {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "user":
		return u.AsUser()
	}
	return nil
}

func (u BetaPluginOwnerUnion) AsOrganization() (v BetaPluginOwnerOrganization) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginOwnerUnion) AsUser() (v BetaPluginOwnerUser) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginOwnerUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginOwnerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How far the served version reaches: `remote` when it declares an MCP server or a
// CLI, `privileged` when it declares a hook, monitor, language server or settings
// but nothing remote, `contained` otherwise; null when not classifiable.
type BetaPluginReach string

const (
	BetaPluginReachContained  BetaPluginReach = "contained"
	BetaPluginReachPrivileged BetaPluginReach = "privileged"
	BetaPluginReachRemote     BetaPluginReach = "remote"
)

type BetaPluginAPIActor struct {
	// The key's ID.
	APIKeyID string `json:"api_key_id" api:"required"`
	// An Admin API key, in the same form the Compliance API activity feed uses for it.
	Type constant.APIActor `json:"type" default:"api_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APIKeyID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginAPIActor) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginAPIActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginComponent struct {
	// What the component declares about itself; always null for MCP servers, hooks,
	// and CLIs.
	Description string `json:"description" api:"required"`
	// The component's name: a skill's, command's or agent's name, an MCP server's key
	// in the manifest, the event a hook runs on, or a CLI's executable.
	Name string `json:"name" api:"required"`
	// The kind of component.
	//
	// Any of "agent", "cli", "command", "hook", "mcp_server", "skill".
	Type BetaPluginComponentType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginComponent) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginComponent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The kind of component.
type BetaPluginComponentType string

const (
	BetaPluginComponentTypeAgent     BetaPluginComponentType = "agent"
	BetaPluginComponentTypeCli       BetaPluginComponentType = "cli"
	BetaPluginComponentTypeCommand   BetaPluginComponentType = "command"
	BetaPluginComponentTypeHook      BetaPluginComponentType = "hook"
	BetaPluginComponentTypeMCPServer BetaPluginComponentType = "mcp_server"
	BetaPluginComponentTypeSkill     BetaPluginComponentType = "skill"
)

type BetaPluginContentScan struct {
	// The scan's verdict; set only when `status` is `completed`.
	//
	// Any of "fail", "pass", "unknown", "warn".
	Assessment BetaPluginContentScanAssessment `json:"assessment" api:"required"`
	// The primary mechanism behind a `warn` or `fail`, such as `credential-exposure`
	// or `guardrail-tampering`; a mechanism this API does not yet name reads as
	// `other`. Null on a `pass`, whenever `assessment` is null, and when no mechanism
	// is reported for the verdict.
	Reason string `json:"reason" api:"required"`
	// `processing` while a scan runs, `completed` when it ran to completion, `errored`
	// when it could not run or its outcome cannot be read.
	//
	// Any of "completed", "errored", "processing".
	Status BetaPluginContentScanStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Assessment  respjson.Field
		Reason      respjson.Field
		Status      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginContentScan) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginContentScan) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The scan's verdict; set only when `status` is `completed`.
type BetaPluginContentScanAssessment string

const (
	BetaPluginContentScanAssessmentFail    BetaPluginContentScanAssessment = "fail"
	BetaPluginContentScanAssessmentPass    BetaPluginContentScanAssessment = "pass"
	BetaPluginContentScanAssessmentUnknown BetaPluginContentScanAssessment = "unknown"
	BetaPluginContentScanAssessmentWarn    BetaPluginContentScanAssessment = "warn"
)

// `processing` while a scan runs, `completed` when it ran to completion, `errored`
// when it could not run or its outcome cannot be read.
type BetaPluginContentScanStatus string

const (
	BetaPluginContentScanStatusCompleted  BetaPluginContentScanStatus = "completed"
	BetaPluginContentScanStatusErrored    BetaPluginContentScanStatus = "errored"
	BetaPluginContentScanStatusProcessing BetaPluginContentScanStatus = "processing"
)

type BetaPluginOwnerOrganization struct {
	// The Plugin lives in a plugin marketplace the organization owns.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginOwnerOrganization) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginOwnerOrganization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginOwnerUser struct {
	// The Plugin lives in one member's personal plugin marketplace.
	Type constant.User `json:"type" default:"user"`
	// The member's User ID.
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
func (r BetaPluginOwnerUser) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginOwnerUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginTargetOrganization struct {
	// Every member of the organization.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginTargetOrganization) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginTargetOrganization) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginTargetOrganizationMember struct {
	// One member of the organization.
	Type constant.OrganizationMember `json:"type" default:"organization_member"`
	// The member's User ID.
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
func (r BetaPluginTargetOrganizationMember) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginTargetOrganizationMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginTargetRBACGroup struct {
	// The RBAC Group's ID.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// An RBAC Group.
	Type constant.RBACGroup `json:"type" default:"rbac_group"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RBACGroupID respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginTargetRBACGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginTargetRBACGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginUserActor struct {
	// The member's email address; may be null, for example when they are no longer a
	// member of the organization.
	EmailAddress string `json:"email_address" api:"required"`
	// A member of the organization.
	Type constant.UserActor `json:"type" default:"user_actor"`
	// The member's User ID.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EmailAddress respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationPluginNewParams struct {
	// The version's files: one part per file, the part's filename being the file's
	// path within the Plugin (for example `skills/review-pr/SKILL.md`), or a single
	// `.zip` or `.plugin` archive holding them all. On the wire each part is named
	// `files[]`, and a part named plain `files` is not read; with cURL,
	// `-F 'files[]=@SKILL.md;filename=skills/review-pr/SKILL.md'`. The files must
	// include the manifest, `.claude-plugin/plugin.json`.
	Files []io.Reader `json:"files,omitzero" api:"required" format:"binary"`
	// ID of the organization-owned plugin marketplace to create the Plugin in
	// (prefixed `marketplace_`). It must be a `manual` marketplace, one whose Plugins
	// are uploaded rather than synchronized from a repository. When omitted, the
	// Plugin is created in the organization's library marketplace, an
	// organization-owned `manual` marketplace created on first use.
	MarketplaceID param.Opt[string] `json:"marketplace_id,omitzero"`
	// Release notes stored with the version and shown in its version history in
	// claude.ai; up to 5,000 characters.
	ReleaseNotes param.Opt[string] `json:"release_notes,omitzero"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type BetaOrganizationPluginGetParams struct {
	// For a `read:org_audit` or `read:compliance_org_data` key created for all of a
	// parent organization's linked organizations: a child organization of that parent
	// to read instead of the organization the key was created in, given as the
	// organization's UUID or its `org_`-prefixed ID. A value that is neither returns a
	// 400; an organization that is not a child of the key's parent, or where the
	// Plugins API is not available, returns a 404. Any other key may pass only its own
	// organization's ID here; another organization returns a 404.
	OrganizationID param.Opt[string] `query:"organization_id,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginGetParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationPluginGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationPluginUpdateParams struct {
	// Serve this version of the Plugin (prefixed `pluginver_`) and pin the served
	// version to it; `latest` is not accepted.
	ServedVersionID string `json:"served_version_id" api:"required"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationPluginUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationPluginUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationPluginListParams struct {
	// RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// RFC 3339 timestamp bound; combine [gte], [gt], [lte], [lt].
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Only Plugins in this plugin marketplace (prefixed `marketplace_`).
	MarketplaceID param.Opt[string] `query:"marketplace_id,omitzero" json:"-"`
	// For a `read:org_audit` or `read:compliance_org_data` key created for all of a
	// parent organization's linked organizations: a child organization of that parent
	// to read instead of the organization the key was created in, given as the
	// organization's UUID or its `org_`-prefixed ID. A value that is neither returns a
	// 400; an organization that is not a child of the key's parent, or where the
	// Plugins API is not available, returns a 404. Any other key may pass only its own
	// organization's ID here; another organization returns a 404.
	OrganizationID param.Opt[string] `query:"organization_id,omitzero" json:"-"`
	// Only Plugins in this member's personal plugin marketplaces (prefixed `user_`); a
	// removed member's ID is accepted.
	OwnerUserID param.Opt[string] `query:"owner_user_id,omitzero" json:"-"`
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `100`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// `organization` for Plugins in the organization's plugin marketplaces, `user` for
	// Plugins in members' personal plugin marketplaces.
	//
	// Any of "organization", "user".
	OwnerType BetaOrganizationPluginListParamsOwnerType `query:"owner_type,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationPluginListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// `organization` for Plugins in the organization's plugin marketplaces, `user` for
// Plugins in members' personal plugin marketplaces.
type BetaOrganizationPluginListParamsOwnerType string

const (
	BetaOrganizationPluginListParamsOwnerTypeOrganization BetaOrganizationPluginListParamsOwnerType = "organization"
	BetaOrganizationPluginListParamsOwnerTypeUser         BetaOrganizationPluginListParamsOwnerType = "user"
)

type BetaOrganizationPluginDeleteParams struct {
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

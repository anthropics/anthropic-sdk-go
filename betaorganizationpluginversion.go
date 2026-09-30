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

// BetaOrganizationPluginVersionService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationPluginVersionService] method instead.
type BetaOrganizationPluginVersionService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationPluginVersionService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationPluginVersionService(opts ...option.RequestOption) (r BetaOrganizationPluginVersionService) {
	r = BetaOrganizationPluginVersionService{}
	r.Options = opts
	return
}

// Add a version to an organization-owned Plugin by uploading the new version's
// files; it becomes the version served to members unless the Plugin's served
// version has been pinned.
//
// The upload is the same `multipart/form-data` as creating a Plugin: the version's
// files (`files`, each part sent as `files[]`) and optional `release_notes`. The
// uploaded manifest's `name` must equal the Plugin's `name`. Returns the stored
// version; read the Plugin back to see which version it serves.
//
// Only a Plugin in a `manual` marketplace takes uploads; a Plugin synchronized
// from a repository gets its versions from the repository. When the Plugin is in
// the organization's library marketplace, a version that adds a skill with the
// name of an organization skill (a skill an administrator uploaded for the whole
// organization in claude.ai) is refused with a 409: `error_code`
// `skill_name_taken`, with that name in `details.skill_name`. A 503 with
// `error_code` `registration_pending` means the version was stored but is not yet
// usable; a later version create on the Plugin completes it.
//
// For a worked example, see
// [Create a version](/docs/en/manage-claude/plugins-api#create-a-version) in the
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
func (r *BetaOrganizationPluginVersionService) New(ctx context.Context, pluginID string, params BetaOrganizationPluginVersionNewParams, opts ...option.RequestOption) (res *BetaPluginVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/versions?beta=true", pluginID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve one version of a Plugin by its ID, or the Plugin's newest version.
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
func (r *BetaOrganizationPluginVersionService) Get(ctx context.Context, version string, params BetaOrganizationPluginVersionGetParams, opts ...option.RequestOption) (res *BetaPluginVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if params.PluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/versions/%s?beta=true", params.PluginID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// List a Plugin's versions, newest first.
//
// The first item of the first page is the version the Plugin's `latest_version_id`
// refers to.
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
func (r *BetaOrganizationPluginVersionService) List(ctx context.Context, pluginID string, params BetaOrganizationPluginVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaPluginVersion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01"), option.WithResponseInto(&raw)}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/versions?beta=true", pluginID)
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

// List a Plugin's versions, newest first.
//
// The first item of the first page is the version the Plugin's `latest_version_id`
// refers to.
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
func (r *BetaOrganizationPluginVersionService) ListAutoPaging(ctx context.Context, pluginID string, params BetaOrganizationPluginVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaPluginVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, pluginID, params, opts...))
}

// Download one version's `.zip` archive, exactly as stored. Each download of a
// Plugin from a member's personal plugin marketplace is recorded on the Compliance
// API activity feed.
//
// The response body is the archive (`Content-Type: application/zip`), sent as an
// attachment whose filename is derived from the Plugin's name; name saved files
// from the IDs in the request path, since that filename is not unique.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `read:org_audit` scope, or a Compliance Access Key with the
// `read:compliance_org_data` scope.
//
// Every read scope above (`read:plugins`, `read:org_audit`, and
// `read:compliance_org_data`) can download the files of plugins in members'
// personal marketplaces, including files that claude.ai's admin settings do not
// show, and a `read:org_audit` or `read:compliance_org_data` key created for all
// of your parent organization's linked organizations can do this in any
// organization under it that has access to this API, by passing `organization_id`.
// Each such download records a `claude_plugin_archive_accessed` event on the
// Compliance API activity feed, identifying the key, the plugin, the version, and
// the member. Downloads of organization-owned plugins are not recorded.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginVersionService) Download(ctx context.Context, version string, params BetaOrganizationPluginVersionDownloadParams, opts ...option.RequestOption) (res *http.Response, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01"), option.WithHeader("Accept", "application/binary")}, opts...)
	if params.PluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	if version == "" {
		err = errors.New("missing required version parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/versions/%s/content?beta=true", params.PluginID, version)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type BetaPluginVersion struct {
	// The version's ID.
	ID string `json:"id" api:"required"`
	// What the version contains; null when not enumerated.
	Components []BetaPluginComponent `json:"components" api:"required"`
	// This version's content scan; null when it has not been scanned.
	ContentScan BetaPluginContentScan `json:"content_scan" api:"required"`
	// RFC 3339.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Who uploaded this version; null when not recorded.
	CreatedBy BetaPluginVersionCreatedByUnion `json:"created_by" api:"required"`
	// The manifest's description; null when it declares none.
	Description string `json:"description" api:"required"`
	// The manifest's display name; null when it declares none.
	DisplayName string `json:"display_name" api:"required"`
	// The version string the manifest declares; null when it declares none.
	ManifestVersion string `json:"manifest_version" api:"required"`
	// The Plugin's ID.
	PluginID string `json:"plugin_id" api:"required"`
	// How far the version reaches: `remote`, `privileged` or `contained`, as on the
	// Plugin; null when not classifiable.
	//
	// Any of "contained", "privileged", "remote".
	Reach BetaPluginVersionReach `json:"reach" api:"required"`
	// As supplied with the upload; null when none were supplied.
	ReleaseNotes string `json:"release_notes" api:"required"`
	// Always `plugin_version`.
	Type constant.PluginVersion `json:"type" default:"plugin_version"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Components      respjson.Field
		ContentScan     respjson.Field
		CreatedAt       respjson.Field
		CreatedBy       respjson.Field
		Description     respjson.Field
		DisplayName     respjson.Field
		ManifestVersion respjson.Field
		PluginID        respjson.Field
		Reach           respjson.Field
		ReleaseNotes    respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginVersion) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaPluginVersionCreatedByUnion contains all possible properties and values from
// [BetaPluginUserActor], [BetaPluginAPIActor].
//
// Use the [BetaPluginVersionCreatedByUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginVersionCreatedByUnion struct {
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

// anyBetaPluginVersionCreatedBy is implemented by each variant of
// [BetaPluginVersionCreatedByUnion] to add type safety for the return type of
// [BetaPluginVersionCreatedByUnion.AsAny]
type anyBetaPluginVersionCreatedBy interface {
	implBetaPluginVersionCreatedByUnion()
}

func (BetaPluginUserActor) implBetaPluginVersionCreatedByUnion() {}
func (BetaPluginAPIActor) implBetaPluginVersionCreatedByUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginVersionCreatedByUnion.AsAny().(type) {
//	case anthropic.BetaPluginUserActor:
//	case anthropic.BetaPluginAPIActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginVersionCreatedByUnion) AsAny() anyBetaPluginVersionCreatedBy {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "api_actor":
		return u.AsAPIActor()
	}
	return nil
}

func (u BetaPluginVersionCreatedByUnion) AsUserActor() (v BetaPluginUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginVersionCreatedByUnion) AsAPIActor() (v BetaPluginAPIActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginVersionCreatedByUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginVersionCreatedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How far the version reaches: `remote`, `privileged` or `contained`, as on the
// Plugin; null when not classifiable.
type BetaPluginVersionReach string

const (
	BetaPluginVersionReachContained  BetaPluginVersionReach = "contained"
	BetaPluginVersionReachPrivileged BetaPluginVersionReach = "privileged"
	BetaPluginVersionReachRemote     BetaPluginVersionReach = "remote"
)

type BetaOrganizationPluginVersionNewParams struct {
	// The version's files: one part per file, the part's filename being the file's
	// path within the Plugin (for example `skills/review-pr/SKILL.md`), or a single
	// `.zip` or `.plugin` archive holding them all. On the wire each part is named
	// `files[]`, and a part named plain `files` is not read; with cURL,
	// `-F 'files[]=@SKILL.md;filename=skills/review-pr/SKILL.md'`. The files must
	// include the manifest, `.claude-plugin/plugin.json`.
	Files []io.Reader `json:"files,omitzero" api:"required" format:"binary"`
	// Release notes stored with the version and shown in its version history in
	// claude.ai; up to 5,000 characters.
	ReleaseNotes param.Opt[string] `json:"release_notes,omitzero"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginVersionNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type BetaOrganizationPluginVersionGetParams struct {
	// ID of the Plugin (prefixed `plugin_`).
	PluginID string `path:"plugin_id" api:"required" json:"-"`
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

// URLQuery serializes [BetaOrganizationPluginVersionGetParams]'s query parameters
// as `url.Values`.
func (r BetaOrganizationPluginVersionGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationPluginVersionListParams struct {
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
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginVersionListParams]'s query parameters
// as `url.Values`.
func (r BetaOrganizationPluginVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationPluginVersionDownloadParams struct {
	// ID of the Plugin (prefixed `plugin_`).
	PluginID string `path:"plugin_id" api:"required" json:"-"`
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

// URLQuery serializes [BetaOrganizationPluginVersionDownloadParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationPluginVersionDownloadParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

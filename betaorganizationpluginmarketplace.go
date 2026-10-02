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
	"strings"
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

// BetaOrganizationPluginMarketplaceService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationPluginMarketplaceService] method instead.
type BetaOrganizationPluginMarketplaceService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationPluginMarketplaceService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationPluginMarketplaceService(opts ...option.RequestOption) (r BetaOrganizationPluginMarketplaceService) {
	r = BetaOrganizationPluginMarketplaceService{}
	r.Options = opts
	return
}

// Retrieve a plugin marketplace by ID.
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
func (r *BetaOrganizationPluginMarketplaceService) Get(ctx context.Context, marketplaceID string, params BetaOrganizationPluginMarketplaceGetParams, opts ...option.RequestOption) (res *BetaPluginMarketplace, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if marketplaceID == "" {
		err = errors.New("missing required marketplace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugin_marketplaces/%s?beta=true", marketplaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Set the default installation setting of one of the organization's own plugin
// marketplaces. Every Plugin in it without a setting of its own gets this default
// as its organization-wide setting, including Plugins added later.
//
// Pass it as `default_installation_preference`. A member's personal marketplace
// cannot be updated here (403).
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginMarketplaceService) Update(ctx context.Context, marketplaceID string, params BetaOrganizationPluginMarketplaceUpdateParams, opts ...option.RequestOption) (res *BetaPluginMarketplace, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if marketplaceID == "" {
		err = errors.New("missing required marketplace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugin_marketplaces/%s?beta=true", marketplaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the plugin marketplaces Plugins live in, newest first: the organization's
// own and its members' personal ones.
//
// Plugin marketplaces are created, connected to a repository and deleted in
// claude.ai, not through this API. The organization's library marketplace, the
// organization-owned `manual` marketplace that uploads go to when no marketplace
// is named, is created the first time something is put in it and is listed from
// then on.
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
func (r *BetaOrganizationPluginMarketplaceService) List(ctx context.Context, params BetaOrganizationPluginMarketplaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaPluginMarketplace], err error) {
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
	path := "v1/organizations/plugin_marketplaces?beta=true"
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

// List the plugin marketplaces Plugins live in, newest first: the organization's
// own and its members' personal ones.
//
// Plugin marketplaces are created, connected to a repository and deleted in
// claude.ai, not through this API. The organization's library marketplace, the
// organization-owned `manual` marketplace that uploads go to when no marketplace
// is named, is created the first time something is put in it and is listed from
// then on.
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
func (r *BetaOrganizationPluginMarketplaceService) ListAutoPaging(ctx context.Context, params BetaOrganizationPluginMarketplaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaPluginMarketplace] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Check whether a plugin marketplace, uploaded as a `.zip` of the marketplace
// directory, would synchronize into claude.ai, without connecting or storing it.
//
// To check a public GitHub repository instead, use Validate Plugin Marketplace
// Repository.
//
// The report says whether `marketplace.json` is well-formed, which plugins a
// synchronization would skip and why, and which plugins would synchronize only in
// part, with some files left out. An archive that cannot be read as a marketplace
// is reported, not refused: the response is a report with `valid: false`. Plugin
// sources outside the marketplace are fetched anonymously from GitHub, so a
// private one is reported as not found; a source on any other host is not fetched
// here, and the report notes that it will be checked when the marketplace actually
// synchronizes.
//
// Nothing is recorded on the Compliance API activity feed.
//
// For a worked example, see
// [Validate marketplace content](/docs/en/manage-claude/plugins-api#validate-marketplace-content)
// in the Plugins API guide.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `write:plugins` scope; `read:org_audit` and `read:compliance_org_data` do not
// grant it.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginMarketplaceService) ValidateArchive(ctx context.Context, params BetaOrganizationPluginMarketplaceValidateArchiveParams, opts ...option.RequestOption) (res *BetaPluginMarketplaceValidationReport, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	path := "v1/organizations/plugin_marketplaces/validate_archive?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Check whether a plugin marketplace held in a public GitHub repository would
// synchronize into claude.ai, without connecting or storing it.
//
// To check a `.zip` of the marketplace directory instead, use Validate Plugin
// Marketplace Archive.
//
// The report says whether `marketplace.json` is well-formed, which plugins a
// synchronization would skip and why, and which plugins would synchronize only in
// part, with some files left out. A repository that is missing, private, or has no
// such branch or commit is reported, not refused: the response is a report with
// `valid: false`. Plugin sources outside the marketplace are fetched anonymously
// from GitHub, so a private one is reported as not found; a source on any other
// host is not fetched here, and the report notes that it will be checked when the
// marketplace actually synchronizes.
//
// Nothing is recorded on the Compliance API activity feed.
//
// For a worked example, see
// [Validate marketplace content](/docs/en/manage-claude/plugins-api#validate-marketplace-content)
// in the Plugins API guide.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `write:plugins` scope; `read:org_audit` and `read:compliance_org_data` do not
// grant it.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginMarketplaceService) ValidateRepository(ctx context.Context, params BetaOrganizationPluginMarketplaceValidateRepositoryParams, opts ...option.RequestOption) (res *BetaPluginMarketplaceValidationReport, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	path := "v1/organizations/plugin_marketplaces/validate_repository?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BetaPluginMarketplace struct {
	// The plugin marketplace's ID, prefixed `marketplace_`.
	ID string `json:"id" api:"required"`
	// RFC 3339.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Organization plugin marketplace: the organization-wide setting every Plugin in
	// it with no setting of its own gets. Null for a member's personal plugin
	// marketplace. One of `required`, `auto_install`, `available`, `not_available`; a
	// value this API does not yet name is returned as stored.
	//
	// Any of "auto_install", "available", "not_available", "required".
	DefaultInstallationPreference BetaPluginMarketplaceDefaultInstallationPreference `json:"default_installation_preference" api:"required"`
	// RFC 3339. When the most recent synchronization attempt to finish did so,
	// whatever its outcome; for a repository plugin marketplace no synchronization has
	// run on yet, when it was created. Null for a plugin marketplace that is not
	// synchronized from a repository.
	LastSyncEndedAt time.Time `json:"last_sync_ended_at" api:"required" format:"date-time"`
	// The commit the last synchronization attempt that reached the repository read,
	// whether or not its content was then accepted (see `sync_status`); an attempt
	// that ends `failed_auth` or `failed_transient` leaves it unchanged. Null until an
	// attempt has first read the repository, and for a plugin marketplace that is not
	// synchronized from a repository.
	LastSyncReadSha string `json:"last_sync_read_sha" api:"required"`
	// Fixed for the plugin marketplace's lifetime.
	Name string `json:"name" api:"required"`
	// The organization, or the member whose personal plugin marketplace it is.
	Owner BetaPluginMarketplaceOwnerUnion `json:"owner" api:"required"`
	// Where the plugin marketplace's Plugins come from: `manual` when they are
	// uploaded; `github`, `gitlab` or `public_git` when they are synchronized from the
	// Git repository the owner connected, into which nothing can be uploaded;
	// `directory` is Anthropic's own catalog, which this API does not list. A value
	// this API does not yet name is returned as stored.
	//
	// Any of "directory", "github", "gitlab", "manual", "public_git".
	Source BetaPluginMarketplaceSource `json:"source" api:"required"`
	// Outcome of the plugin marketplace's most recent synchronization: one of
	// `success`, `in_progress`, `failed_content`, `failed_transient`, `failed_auth`,
	// `failed_limits`; a value this API does not yet name is returned as stored. Null
	// until a synchronization is first attempted — so always for a `manual` plugin
	// marketplace.
	//
	// Any of "failed_auth", "failed_content", "failed_limits", "failed_transient",
	// "in_progress", "success".
	SyncStatus BetaPluginMarketplaceSyncStatus `json:"sync_status" api:"required"`
	// Always `plugin_marketplace`.
	Type constant.PluginMarketplace `json:"type" default:"plugin_marketplace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                            respjson.Field
		CreatedAt                     respjson.Field
		DefaultInstallationPreference respjson.Field
		LastSyncEndedAt               respjson.Field
		LastSyncReadSha               respjson.Field
		Name                          respjson.Field
		Owner                         respjson.Field
		Source                        respjson.Field
		SyncStatus                    respjson.Field
		Type                          respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginMarketplace) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginMarketplace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Organization plugin marketplace: the organization-wide setting every Plugin in
// it with no setting of its own gets. Null for a member's personal plugin
// marketplace. One of `required`, `auto_install`, `available`, `not_available`; a
// value this API does not yet name is returned as stored.
type BetaPluginMarketplaceDefaultInstallationPreference string

const (
	BetaPluginMarketplaceDefaultInstallationPreferenceAutoInstall  BetaPluginMarketplaceDefaultInstallationPreference = "auto_install"
	BetaPluginMarketplaceDefaultInstallationPreferenceAvailable    BetaPluginMarketplaceDefaultInstallationPreference = "available"
	BetaPluginMarketplaceDefaultInstallationPreferenceNotAvailable BetaPluginMarketplaceDefaultInstallationPreference = "not_available"
	BetaPluginMarketplaceDefaultInstallationPreferenceRequired     BetaPluginMarketplaceDefaultInstallationPreference = "required"
)

// BetaPluginMarketplaceOwnerUnion contains all possible properties and values from
// [BetaPluginOwnerOrganization], [BetaPluginOwnerUser].
//
// Use the [BetaPluginMarketplaceOwnerUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginMarketplaceOwnerUnion struct {
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

// anyBetaPluginMarketplaceOwner is implemented by each variant of
// [BetaPluginMarketplaceOwnerUnion] to add type safety for the return type of
// [BetaPluginMarketplaceOwnerUnion.AsAny]
type anyBetaPluginMarketplaceOwner interface {
	implBetaPluginMarketplaceOwnerUnion()
}

func (BetaPluginOwnerOrganization) implBetaPluginMarketplaceOwnerUnion() {}
func (BetaPluginOwnerUser) implBetaPluginMarketplaceOwnerUnion()         {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginMarketplaceOwnerUnion.AsAny().(type) {
//	case anthropic.BetaPluginOwnerOrganization:
//	case anthropic.BetaPluginOwnerUser:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginMarketplaceOwnerUnion) AsAny() anyBetaPluginMarketplaceOwner {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "user":
		return u.AsUser()
	}
	return nil
}

func (u BetaPluginMarketplaceOwnerUnion) AsOrganization() (v BetaPluginOwnerOrganization) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginMarketplaceOwnerUnion) AsUser() (v BetaPluginOwnerUser) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginMarketplaceOwnerUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginMarketplaceOwnerUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where the plugin marketplace's Plugins come from: `manual` when they are
// uploaded; `github`, `gitlab` or `public_git` when they are synchronized from the
// Git repository the owner connected, into which nothing can be uploaded;
// `directory` is Anthropic's own catalog, which this API does not list. A value
// this API does not yet name is returned as stored.
type BetaPluginMarketplaceSource string

const (
	BetaPluginMarketplaceSourceDirectory BetaPluginMarketplaceSource = "directory"
	BetaPluginMarketplaceSourceGitHub    BetaPluginMarketplaceSource = "github"
	BetaPluginMarketplaceSourceGitlab    BetaPluginMarketplaceSource = "gitlab"
	BetaPluginMarketplaceSourceManual    BetaPluginMarketplaceSource = "manual"
	BetaPluginMarketplaceSourcePublicGit BetaPluginMarketplaceSource = "public_git"
)

// Outcome of the plugin marketplace's most recent synchronization: one of
// `success`, `in_progress`, `failed_content`, `failed_transient`, `failed_auth`,
// `failed_limits`; a value this API does not yet name is returned as stored. Null
// until a synchronization is first attempted — so always for a `manual` plugin
// marketplace.
type BetaPluginMarketplaceSyncStatus string

const (
	BetaPluginMarketplaceSyncStatusFailedAuth      BetaPluginMarketplaceSyncStatus = "failed_auth"
	BetaPluginMarketplaceSyncStatusFailedContent   BetaPluginMarketplaceSyncStatus = "failed_content"
	BetaPluginMarketplaceSyncStatusFailedLimits    BetaPluginMarketplaceSyncStatus = "failed_limits"
	BetaPluginMarketplaceSyncStatusFailedTransient BetaPluginMarketplaceSyncStatus = "failed_transient"
	BetaPluginMarketplaceSyncStatusInProgress      BetaPluginMarketplaceSyncStatus = "in_progress"
	BetaPluginMarketplaceSyncStatusSuccess         BetaPluginMarketplaceSyncStatus = "success"
)

type BetaPluginMarketplaceValidationPluginError struct {
	// Why the plugin would be skipped by a synchronization.
	Error string `json:"error" api:"required"`
	// A stable identifier for the reason — the value to branch on.
	ErrorCode string `json:"error_code" api:"required"`
	// The plugin's name, as its entry in marketplace.json declares it.
	Name string `json:"name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		ErrorCode   respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginMarketplaceValidationPluginError) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginMarketplaceValidationPluginError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginMarketplaceValidationPluginWarning struct {
	// A stable identifier for the kind of warning.
	ErrorCode string `json:"error_code" api:"required"`
	// What would be left out, and why.
	Message string `json:"message" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginMarketplaceValidationPluginWarning) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginMarketplaceValidationPluginWarning) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPluginMarketplaceValidationPluginWarnings struct {
	// The plugin's name, as its entry in marketplace.json declares it.
	Name string `json:"name" api:"required"`
	// The parts of the plugin a synchronization would leave out.
	Warnings []BetaPluginMarketplaceValidationPluginWarning `json:"warnings" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Warnings    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginMarketplaceValidationPluginWarnings) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginMarketplaceValidationPluginWarnings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The outcome of validating plugin marketplace content: a report, not a stored
// object, so nothing in it can be retrieved afterwards.
type BetaPluginMarketplaceValidationReport struct {
	// The full SHA of the commit that was validated: for a repository, the commit that
	// was read; for an uploaded archive, the commit recorded in the archive's comment
	// (as a Git host's download writes it; not verified), else null.
	CommitSha string `json:"commit_sha" api:"required"`
	// Set when nothing could be validated: the repository or archive could not be
	// read, or marketplace.json is missing, malformed or over a limit. Null otherwise.
	ManifestError string `json:"manifest_error" api:"required"`
	// A stable identifier for `manifest_error`; null when that is.
	ManifestErrorCode string `json:"manifest_error_code" api:"required"`
	// One entry per plugin a synchronization would skip entirely, keyed by the
	// plugin's name in marketplace.json.
	PluginErrors []BetaPluginMarketplaceValidationPluginError `json:"plugin_errors" api:"required"`
	// One entry per plugin that would synchronize with some of its contents left out,
	// keyed by the plugin's name in marketplace.json.
	PluginWarnings []BetaPluginMarketplaceValidationPluginWarnings `json:"plugin_warnings" api:"required"`
	// For a repository, the branch that was read by name: the one requested, or else
	// the branch a synchronization of this repository is set to read. Null when no
	// branch is named or set and the repository's default branch was read, for a
	// request by commit SHA, and for an uploaded archive.
	Ref string `json:"ref" api:"required"`
	// How many plugins marketplace.json declares; 0 when it could not be read.
	TotalPluginCount int64 `json:"total_plugin_count" api:"required"`
	// Always `plugin_marketplace_validation_report`.
	Type constant.PluginMarketplaceValidationReport `json:"type" default:"plugin_marketplace_validation_report"`
	// True when marketplace.json is well-formed and no plugin would be skipped;
	// warnings never make it false.
	Valid bool `json:"valid" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CommitSha         respjson.Field
		ManifestError     respjson.Field
		ManifestErrorCode respjson.Field
		PluginErrors      respjson.Field
		PluginWarnings    respjson.Field
		Ref               respjson.Field
		TotalPluginCount  respjson.Field
		Type              respjson.Field
		Valid             respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginMarketplaceValidationReport) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginMarketplaceValidationReport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationPluginMarketplaceGetParams struct {
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

// URLQuery serializes [BetaOrganizationPluginMarketplaceGetParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationPluginMarketplaceGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationPluginMarketplaceUpdateParams struct {
	// The organization-wide installation setting every Plugin in the marketplace
	// without one of its own gets: one of `required`, `auto_install`, `available`,
	// `not_available`. Once set it can be changed but not removed.
	//
	// Any of "auto_install", "available", "not_available", "required".
	DefaultInstallationPreference BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference `json:"default_installation_preference,omitzero" api:"required"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginMarketplaceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationPluginMarketplaceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationPluginMarketplaceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The organization-wide installation setting every Plugin in the marketplace
// without one of its own gets: one of `required`, `auto_install`, `available`,
// `not_available`. Once set it can be changed but not removed.
type BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference string

const (
	BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreferenceAutoInstall  BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference = "auto_install"
	BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreferenceAvailable    BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference = "available"
	BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreferenceNotAvailable BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference = "not_available"
	BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreferenceRequired     BetaOrganizationPluginMarketplaceUpdateParamsDefaultInstallationPreference = "required"
)

type BetaOrganizationPluginMarketplaceListParams struct {
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
	// `organization` for the organization's plugin marketplaces, `user` for members'
	// personal plugin marketplaces.
	//
	// Any of "organization", "user".
	OwnerType BetaOrganizationPluginMarketplaceListParamsOwnerType `query:"owner_type,omitzero" json:"-"`
	// Only plugin marketplaces with this `source`: `manual` for those whose Plugins
	// are uploaded; `github`, `gitlab` or `public_git` for those synchronized from a
	// Git repository. `directory` (Anthropic's catalog) is never listed here.
	//
	// Any of "directory", "github", "gitlab", "manual", "public_git".
	Source BetaOrganizationPluginMarketplaceListParamsSource `query:"source,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginMarketplaceListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationPluginMarketplaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// `organization` for the organization's plugin marketplaces, `user` for members'
// personal plugin marketplaces.
type BetaOrganizationPluginMarketplaceListParamsOwnerType string

const (
	BetaOrganizationPluginMarketplaceListParamsOwnerTypeOrganization BetaOrganizationPluginMarketplaceListParamsOwnerType = "organization"
	BetaOrganizationPluginMarketplaceListParamsOwnerTypeUser         BetaOrganizationPluginMarketplaceListParamsOwnerType = "user"
)

// Only plugin marketplaces with this `source`: `manual` for those whose Plugins
// are uploaded; `github`, `gitlab` or `public_git` for those synchronized from a
// Git repository. `directory` (Anthropic's catalog) is never listed here.
type BetaOrganizationPluginMarketplaceListParamsSource string

const (
	BetaOrganizationPluginMarketplaceListParamsSourceDirectory BetaOrganizationPluginMarketplaceListParamsSource = "directory"
	BetaOrganizationPluginMarketplaceListParamsSourceGitHub    BetaOrganizationPluginMarketplaceListParamsSource = "github"
	BetaOrganizationPluginMarketplaceListParamsSourceGitlab    BetaOrganizationPluginMarketplaceListParamsSource = "gitlab"
	BetaOrganizationPluginMarketplaceListParamsSourceManual    BetaOrganizationPluginMarketplaceListParamsSource = "manual"
	BetaOrganizationPluginMarketplaceListParamsSourcePublicGit BetaOrganizationPluginMarketplaceListParamsSource = "public_git"
)

type BetaOrganizationPluginMarketplaceValidateArchiveParams struct {
	// A .zip of the marketplace directory (its contents at the root, or wrapped in one
	// folder as a Git host's download produces), sent as a file part with a filename;
	// DEFLATE- or STORE-compressed, at most 32 MB. A part sent without a filename, a
	// second archive part, or any other form field is a 400; a larger archive is
	// a 413.
	Archive io.Reader `json:"archive,omitzero" api:"required" format:"binary"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginMarketplaceValidateArchiveParams) MarshalMultipart() (data []byte, contentType string, err error) {
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

type BetaOrganizationPluginMarketplaceValidateRepositoryParams struct {
	// The `https://` URL of a public repository on github.com that holds the
	// marketplace. Any other host, a URL with credentials in it, or one that does not
	// name a repository is a 400.
	RepositoryURL string `json:"repository_url" api:"required"`
	// The branch to validate the tip of, or the full 40-character SHA of the commit to
	// validate. When omitted, the branch a synchronization would read (usually the
	// repository's default branch); if that is not the default branch, the report's
	// `ref` says which branch was read. An empty string, or a value that is neither a
	// branch name nor a 40-character SHA, is a 400.
	Ref param.Opt[string] `json:"ref,omitzero"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginMarketplaceValidateRepositoryParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationPluginMarketplaceValidateRepositoryParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationPluginMarketplaceValidateRepositoryParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

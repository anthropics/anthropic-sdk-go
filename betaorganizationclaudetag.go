package anthropic

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeTagService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeTagService] method instead.
type BetaOrganizationClaudeTagService struct {
	Options         []option.RequestOption
	SlackChannels   BetaOrganizationClaudeTagSlackChannelService
	SlackWorkspaces BetaOrganizationClaudeTagSlackWorkspaceService
	Bundles         BetaOrganizationClaudeTagBundleService
}

// NewBetaOrganizationClaudeTagService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationClaudeTagService(opts ...option.RequestOption) (r BetaOrganizationClaudeTagService) {
	r = BetaOrganizationClaudeTagService{}
	r.Options = opts
	r.SlackChannels = NewBetaOrganizationClaudeTagSlackChannelService(opts...)
	r.SlackWorkspaces = NewBetaOrganizationClaudeTagSlackWorkspaceService(opts...)
	r.Bundles = NewBetaOrganizationClaudeTagBundleService(opts...)
	return
}

// No settings level sets a value, and the built-in default applies.
type BetaClaudeTagDefaultSettingSource struct {
	Type constant.Default `json:"type" default:"default"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagDefaultSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagDefaultSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaClaudeTagGuestPolicyUnion contains all possible properties and values from
// [BetaClaudeTagGuestPolicyAllow], [BetaClaudeTagGuestPolicyRestrict],
// [BetaClaudeTagGuestPolicyChannelOnly].
//
// Use the [BetaClaudeTagGuestPolicyUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaClaudeTagGuestPolicyUnion struct {
	// Any of "allow", "restrict", "channel_only".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaClaudeTagGuestPolicy is implemented by each variant of
// [BetaClaudeTagGuestPolicyUnion] to add type safety for the return type of
// [BetaClaudeTagGuestPolicyUnion.AsAny]
type anyBetaClaudeTagGuestPolicy interface {
	implBetaClaudeTagGuestPolicyUnion()
}

func (BetaClaudeTagGuestPolicyAllow) implBetaClaudeTagGuestPolicyUnion()       {}
func (BetaClaudeTagGuestPolicyRestrict) implBetaClaudeTagGuestPolicyUnion()    {}
func (BetaClaudeTagGuestPolicyChannelOnly) implBetaClaudeTagGuestPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaClaudeTagGuestPolicyUnion.AsAny().(type) {
//	case anthropic.BetaClaudeTagGuestPolicyAllow:
//	case anthropic.BetaClaudeTagGuestPolicyRestrict:
//	case anthropic.BetaClaudeTagGuestPolicyChannelOnly:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaClaudeTagGuestPolicyUnion) AsAny() anyBetaClaudeTagGuestPolicy {
	switch u.Type {
	case "allow":
		return u.AsAllow()
	case "restrict":
		return u.AsRestrict()
	case "channel_only":
		return u.AsChannelOnly()
	}
	return nil
}

func (u BetaClaudeTagGuestPolicyUnion) AsAllow() (v BetaClaudeTagGuestPolicyAllow) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagGuestPolicyUnion) AsRestrict() (v BetaClaudeTagGuestPolicyRestrict) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagGuestPolicyUnion) AsChannelOnly() (v BetaClaudeTagGuestPolicyChannelOnly) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaClaudeTagGuestPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaClaudeTagGuestPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaClaudeTagGuestPolicyUnion to a
// BetaClaudeTagGuestPolicyUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaClaudeTagGuestPolicyUnionParam.Overrides()
func (r BetaClaudeTagGuestPolicyUnion) ToParam() BetaClaudeTagGuestPolicyUnionParam {
	return param.Override[BetaClaudeTagGuestPolicyUnionParam](json.RawMessage(r.RawJSON()))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaClaudeTagGuestPolicyUnionParam struct {
	OfAllow       *BetaClaudeTagGuestPolicyAllowParam       `json:",omitzero,inline"`
	OfRestrict    *BetaClaudeTagGuestPolicyRestrictParam    `json:",omitzero,inline"`
	OfChannelOnly *BetaClaudeTagGuestPolicyChannelOnlyParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaClaudeTagGuestPolicyUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAllow, u.OfRestrict, u.OfChannelOnly)
}
func (u *BetaClaudeTagGuestPolicyUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaClaudeTagGuestPolicyUnionParam) asAny() any {
	if !param.IsOmitted(u.OfAllow) {
		return u.OfAllow
	} else if !param.IsOmitted(u.OfRestrict) {
		return u.OfRestrict
	} else if !param.IsOmitted(u.OfChannelOnly) {
		return u.OfChannelOnly
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaClaudeTagGuestPolicyUnionParam) GetType() *string {
	if vt := u.OfAllow; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfRestrict; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfChannelOnly; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaClaudeTagGuestPolicyUnionParam](
		"type",
		apijson.Discriminator[BetaClaudeTagGuestPolicyAllowParam]("allow"),
		apijson.Discriminator[BetaClaudeTagGuestPolicyRestrictParam]("restrict"),
		apijson.Discriminator[BetaClaudeTagGuestPolicyChannelOnlyParam]("channel_only"),
	)
}

// Claude answers in the channel.
type BetaClaudeTagGuestPolicyAllow struct {
	Type constant.Allow `json:"type" default:"allow"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagGuestPolicyAllow) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagGuestPolicyAllow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaClaudeTagGuestPolicyAllow to a
// BetaClaudeTagGuestPolicyAllowParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaClaudeTagGuestPolicyAllowParam.Overrides()
func (r BetaClaudeTagGuestPolicyAllow) ToParam() BetaClaudeTagGuestPolicyAllowParam {
	return param.Override[BetaClaudeTagGuestPolicyAllowParam](json.RawMessage(r.RawJSON()))
}

func NewBetaClaudeTagGuestPolicyAllowParam() BetaClaudeTagGuestPolicyAllowParam {
	return BetaClaudeTagGuestPolicyAllowParam{
		Type: "allow",
	}
}

// Claude answers in the channel.
//
// This struct has a constant value, construct it with
// [NewBetaClaudeTagGuestPolicyAllowParam].
type BetaClaudeTagGuestPolicyAllowParam struct {
	Type constant.Allow `json:"type" default:"allow"`
	paramObj
}

func (r BetaClaudeTagGuestPolicyAllowParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagGuestPolicyAllowParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagGuestPolicyAllowParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude answers with only the Bundles attached to the channel itself.
type BetaClaudeTagGuestPolicyChannelOnly struct {
	Type constant.ChannelOnly `json:"type" default:"channel_only"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagGuestPolicyChannelOnly) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagGuestPolicyChannelOnly) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaClaudeTagGuestPolicyChannelOnly to a
// BetaClaudeTagGuestPolicyChannelOnlyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaClaudeTagGuestPolicyChannelOnlyParam.Overrides()
func (r BetaClaudeTagGuestPolicyChannelOnly) ToParam() BetaClaudeTagGuestPolicyChannelOnlyParam {
	return param.Override[BetaClaudeTagGuestPolicyChannelOnlyParam](json.RawMessage(r.RawJSON()))
}

func NewBetaClaudeTagGuestPolicyChannelOnlyParam() BetaClaudeTagGuestPolicyChannelOnlyParam {
	return BetaClaudeTagGuestPolicyChannelOnlyParam{
		Type: "channel_only",
	}
}

// Claude answers with only the Bundles attached to the channel itself.
//
// This struct has a constant value, construct it with
// [NewBetaClaudeTagGuestPolicyChannelOnlyParam].
type BetaClaudeTagGuestPolicyChannelOnlyParam struct {
	Type constant.ChannelOnly `json:"type" default:"channel_only"`
	paramObj
}

func (r BetaClaudeTagGuestPolicyChannelOnlyParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagGuestPolicyChannelOnlyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagGuestPolicyChannelOnlyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude does not answer in the channel.
type BetaClaudeTagGuestPolicyRestrict struct {
	Type constant.Restrict `json:"type" default:"restrict"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagGuestPolicyRestrict) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagGuestPolicyRestrict) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaClaudeTagGuestPolicyRestrict to a
// BetaClaudeTagGuestPolicyRestrictParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaClaudeTagGuestPolicyRestrictParam.Overrides()
func (r BetaClaudeTagGuestPolicyRestrict) ToParam() BetaClaudeTagGuestPolicyRestrictParam {
	return param.Override[BetaClaudeTagGuestPolicyRestrictParam](json.RawMessage(r.RawJSON()))
}

func NewBetaClaudeTagGuestPolicyRestrictParam() BetaClaudeTagGuestPolicyRestrictParam {
	return BetaClaudeTagGuestPolicyRestrictParam{
		Type: "restrict",
	}
}

// Claude does not answer in the channel.
//
// This struct has a constant value, construct it with
// [NewBetaClaudeTagGuestPolicyRestrictParam].
type BetaClaudeTagGuestPolicyRestrictParam struct {
	Type constant.Restrict `json:"type" default:"restrict"`
	paramObj
}

func (r BetaClaudeTagGuestPolicyRestrictParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagGuestPolicyRestrictParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagGuestPolicyRestrictParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaClaudeTagSettingSourceUnion contains all possible properties and values from
// [BetaClaudeTagSlackChannelSettingSource],
// [BetaClaudeTagSlackWorkspaceSettingSource],
// [BetaClaudeTagSlackParentChannelSettingSource],
// [BetaClaudeTagSlackDefaultSettingSource], [BetaClaudeTagDefaultSettingSource],
// [BetaClaudeTagSlackGridChannelSettingSource].
//
// Use the [BetaClaudeTagSettingSourceUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaClaudeTagSettingSourceUnion struct {
	SlackChannelID string `json:"slack_channel_id"`
	// Any of "slack_channel", "slack_workspace", "slack_parent_channel",
	// "slack_default", "default", "slack_grid_channel".
	Type string `json:"type"`
	// This field is from variant [BetaClaudeTagSlackWorkspaceSettingSource].
	SlackWorkspaceID string `json:"slack_workspace_id"`
	// This field is from variant [BetaClaudeTagSlackGridChannelSettingSource].
	Slack BetaClaudeTagSlackChannelIdentifiers `json:"slack"`
	JSON  struct {
		SlackChannelID   respjson.Field
		Type             respjson.Field
		SlackWorkspaceID respjson.Field
		Slack            respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaClaudeTagSettingSource is implemented by each variant of
// [BetaClaudeTagSettingSourceUnion] to add type safety for the return type of
// [BetaClaudeTagSettingSourceUnion.AsAny]
type anyBetaClaudeTagSettingSource interface {
	implBetaClaudeTagSettingSourceUnion()
}

func (BetaClaudeTagSlackChannelSettingSource) implBetaClaudeTagSettingSourceUnion()       {}
func (BetaClaudeTagSlackWorkspaceSettingSource) implBetaClaudeTagSettingSourceUnion()     {}
func (BetaClaudeTagSlackParentChannelSettingSource) implBetaClaudeTagSettingSourceUnion() {}
func (BetaClaudeTagSlackDefaultSettingSource) implBetaClaudeTagSettingSourceUnion()       {}
func (BetaClaudeTagDefaultSettingSource) implBetaClaudeTagSettingSourceUnion()            {}
func (BetaClaudeTagSlackGridChannelSettingSource) implBetaClaudeTagSettingSourceUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaClaudeTagSettingSourceUnion.AsAny().(type) {
//	case anthropic.BetaClaudeTagSlackChannelSettingSource:
//	case anthropic.BetaClaudeTagSlackWorkspaceSettingSource:
//	case anthropic.BetaClaudeTagSlackParentChannelSettingSource:
//	case anthropic.BetaClaudeTagSlackDefaultSettingSource:
//	case anthropic.BetaClaudeTagDefaultSettingSource:
//	case anthropic.BetaClaudeTagSlackGridChannelSettingSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaClaudeTagSettingSourceUnion) AsAny() anyBetaClaudeTagSettingSource {
	switch u.Type {
	case "slack_channel":
		return u.AsSlackChannel()
	case "slack_workspace":
		return u.AsSlackWorkspace()
	case "slack_parent_channel":
		return u.AsSlackParentChannel()
	case "slack_default":
		return u.AsSlackDefault()
	case "default":
		return u.AsDefault()
	case "slack_grid_channel":
		return u.AsSlackGridChannel()
	}
	return nil
}

func (u BetaClaudeTagSettingSourceUnion) AsSlackChannel() (v BetaClaudeTagSlackChannelSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagSettingSourceUnion) AsSlackWorkspace() (v BetaClaudeTagSlackWorkspaceSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagSettingSourceUnion) AsSlackParentChannel() (v BetaClaudeTagSlackParentChannelSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagSettingSourceUnion) AsSlackDefault() (v BetaClaudeTagSlackDefaultSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagSettingSourceUnion) AsDefault() (v BetaClaudeTagDefaultSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagSettingSourceUnion) AsSlackGridChannel() (v BetaClaudeTagSlackGridChannelSettingSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaClaudeTagSettingSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaClaudeTagSettingSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// For each setting, the settings level its value in `settings` comes from. Without
// `include_inherited`, a setting the channel or workspace stores has a source that
// names the channel or workspace itself, and one it does not store has null, or an
// empty list for `instructions`. With `include_inherited=true` no member is null:
// a value that no settings level sets has the source `default`.
type BetaClaudeTagSettingSources struct {
	// The settings level `enabled` comes from.
	Enabled BetaClaudeTagSettingSourceUnion `json:"enabled" api:"required"`
	// The settings level `environment_id` comes from.
	EnvironmentID BetaClaudeTagSettingSourceUnion `json:"environment_id" api:"required"`
	// The settings level `guest_policy` comes from.
	GuestPolicy BetaClaudeTagSettingSourceUnion `json:"guest_policy" api:"required"`
	// Every settings level that contributes text to `instructions`, broadest first, in
	// the order their text is joined. Empty when none does.
	Instructions []BetaClaudeTagSettingSourceUnion `json:"instructions" api:"required"`
	// The settings level `model` comes from.
	Model BetaClaudeTagSettingSourceUnion `json:"model" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled       respjson.Field
		EnvironmentID respjson.Field
		GuestPolicy   respjson.Field
		Instructions  respjson.Field
		Model         respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSettingSources) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSettingSources) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Tag settings. Without `include_inherited`, a setting that is null is not
// set on the settings level and is inherited. With it, `enabled` and
// `guest_policy` are never null, and `instructions`, `model` and `environment_id`
// are null when no settings level sets them.
type BetaClaudeTagSettings struct {
	// Whether Claude is on. The nearest settings level that sets it decides.
	Enabled bool `json:"enabled" api:"required"`
	// The environment Claude's sessions run in, starting with `env_`, or the
	// identifier of a self-hosted runner pool, starting with `ccpool_`. The nearest
	// settings level that sets it decides.
	EnvironmentID string `json:"environment_id" api:"required"`
	// What Claude does in a channel while a Slack guest is in it. The nearest settings
	// level that sets it decides.
	GuestPolicy BetaClaudeTagGuestPolicyUnion `json:"guest_policy" api:"required"`
	// Instructions for Claude. Instructions add up: the text of every settings level
	// that sets it is joined, broadest first.
	Instructions string `json:"instructions" api:"required"`
	// The model Claude uses: a model identifier or the name of a model line. The
	// nearest settings level that sets it decides. When none does, Claude uses the
	// product's default model.
	Model string `json:"model" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled       respjson.Field
		EnvironmentID respjson.Field
		GuestPolicy   respjson.Field
		Instructions  respjson.Field
		Model         respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSettings) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The settings in a request that sets a Slack channel up. A setting left out is
// not set. A setting sent as null is cleared.
type BetaClaudeTagSettingsCreateParam struct {
	// Whether Claude is on. `true` on a channel turns Claude on in it even when its
	// workspace stores `false`. Null clears it.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// The environment Claude's sessions run in: the identifier of an active
	// environment of the organization, starting with `env_`. Null clears it.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// Instructions for Claude, up to 8,192 characters after leading and trailing white
	// space is trimmed. Null or blank text clears them.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// The model Claude uses: the identifier of a model, or the name of a model line,
	// that the organization is entitled to for Claude Tag. Null clears it.
	Model param.Opt[string] `json:"model,omitzero"`
	// What Claude does in a channel while a Slack guest is in it. Null clears it.
	GuestPolicy BetaClaudeTagGuestPolicyUnionParam `json:"guest_policy,omitzero"`
	paramObj
}

func (r BetaClaudeTagSettingsCreateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagSettingsCreateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagSettingsCreateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The settings to change. A setting left out is unchanged. A setting sent as null
// is cleared, so it is inherited.
type BetaClaudeTagSettingsUpdateParam struct {
	// Whether Claude is on. `true` on a channel turns Claude on in it even when its
	// workspace stores `false`. Null clears it.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// The environment Claude's sessions run in: the identifier of an active
	// environment of the organization, starting with `env_`. Null clears it.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// Instructions for Claude, up to 8,192 characters after leading and trailing white
	// space is trimmed. Null or blank text clears them.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// The model Claude uses: the identifier of a model, or the name of a model line,
	// that the organization is entitled to for Claude Tag. Null clears it.
	Model param.Opt[string] `json:"model,omitzero"`
	// What Claude does in a channel while a Slack guest is in it. Null clears it.
	GuestPolicy BetaClaudeTagGuestPolicyUnionParam `json:"guest_policy,omitzero"`
	paramObj
}

func (r BetaClaudeTagSettingsUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagSettingsUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagSettingsUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The identifiers Slack shows for a channel. They are Slack's own, not identifiers
// of this API.
type BetaClaudeTagSlackChannelIdentifiers struct {
	// Slack's identifier for the channel, starting with `C` or `G`, such as
	// `C0456EFGH`.
	ChannelID string `json:"channel_id" api:"required"`
	// Slack's identifier for the workspace the channel is in, starting with `T`, such
	// as `T0123ABCD`.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChannelID   respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackChannelIdentifiers) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackChannelIdentifiers) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaClaudeTagSlackChannelIdentifiers to a
// BetaClaudeTagSlackChannelIdentifiersParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaClaudeTagSlackChannelIdentifiersParam.Overrides()
func (r BetaClaudeTagSlackChannelIdentifiers) ToParam() BetaClaudeTagSlackChannelIdentifiersParam {
	return param.Override[BetaClaudeTagSlackChannelIdentifiersParam](json.RawMessage(r.RawJSON()))
}

// The identifiers Slack shows for a channel. They are Slack's own, not identifiers
// of this API.
//
// The properties ChannelID, WorkspaceID are required.
type BetaClaudeTagSlackChannelIdentifiersParam struct {
	// Slack's identifier for the channel, starting with `C` or `G`, such as
	// `C0456EFGH`.
	ChannelID string `json:"channel_id" api:"required"`
	// Slack's identifier for the workspace the channel is in, starting with `T`, such
	// as `T0123ABCD`.
	WorkspaceID string `json:"workspace_id" api:"required"`
	paramObj
}

func (r BetaClaudeTagSlackChannelIdentifiersParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagSlackChannelIdentifiersParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagSlackChannelIdentifiersParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The value comes from the addressed Slack channel itself.
type BetaClaudeTagSlackChannelSettingSource struct {
	// The channel's identifier, starting with `slch_`.
	SlackChannelID string                `json:"slack_channel_id" api:"required"`
	Type           constant.SlackChannel `json:"type" default:"slack_channel"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SlackChannelID respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackChannelSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackChannelSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The value comes from the organization's Slack defaults, which have no
// identifier.
type BetaClaudeTagSlackDefaultSettingSource struct {
	Type constant.SlackDefault `json:"type" default:"slack_default"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackDefaultSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackDefaultSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The value comes from the addressed channel as it is set up under another
// workspace of its Enterprise Grid: a channel shared across the Grid keeps one
// Slack channel identifier in every workspace that sees it.
type BetaClaudeTagSlackGridChannelSettingSource struct {
	// Slack's own identifiers for the channel under that other workspace.
	Slack BetaClaudeTagSlackChannelIdentifiers `json:"slack" api:"required"`
	Type  constant.SlackGridChannel            `json:"type" default:"slack_grid_channel"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Slack       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackGridChannelSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackGridChannelSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The value comes from the channel that a code channel was started from. A parent
// channel contributes `instructions` and `environment_id` only.
type BetaClaudeTagSlackParentChannelSettingSource struct {
	// The parent channel's identifier, starting with `slch_`.
	SlackChannelID string                      `json:"slack_channel_id" api:"required"`
	Type           constant.SlackParentChannel `json:"type" default:"slack_parent_channel"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SlackChannelID respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackParentChannelSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackParentChannelSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The identifier Slack shows for a workspace. It is Slack's own, not an identifier
// of this API.
type BetaClaudeTagSlackWorkspaceIdentifiers struct {
	// Slack's identifier for the workspace, starting with `T`, such as `T0123ABCD`.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackWorkspaceIdentifiers) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackWorkspaceIdentifiers) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The value comes from a Slack workspace of the organization.
type BetaClaudeTagSlackWorkspaceSettingSource struct {
	// The workspace's identifier, starting with `slws_`.
	SlackWorkspaceID string                  `json:"slack_workspace_id" api:"required"`
	Type             constant.SlackWorkspace `json:"type" default:"slack_workspace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SlackWorkspaceID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackWorkspaceSettingSource) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackWorkspaceSettingSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

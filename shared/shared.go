package shared

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type APIErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.APIError `json:"type" default:"api_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIErrorObject) RawJSON() string { return r.JSON.raw }
func (r *APIErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (APIErrorObject) ImplErrorObjectUnion() {}

type AuthenticationError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.AuthenticationError `json:"type" default:"authentication_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthenticationError) RawJSON() string { return r.JSON.raw }
func (r *AuthenticationError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (AuthenticationError) ImplErrorObjectUnion() {}

type BillingError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.BillingError `json:"type" default:"billing_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BillingError) RawJSON() string { return r.JSON.raw }
func (r *BillingError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BillingError) ImplErrorObjectUnion() {}

// The API key belongs to the Claude Code workspace, and the user who created it is
// no longer a member of the organization or of that workspace.
type ClaudeCodeKeyCreatorNotMemberErrorDetails struct {
	ErrorCode constant.ClaudeCodeKeyCreatorNotMemberErrorCode `json:"error_code" default:"claude_code_key_creator_not_member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClaudeCodeKeyCreatorNotMemberErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ClaudeCodeKeyCreatorNotMemberErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ClaudeCodeKeyCreatorNotMemberErrorDetails) ImplErrorDetailsUnion() {}

// The request comes from a Claude Code version older than the minimum that the
// requested model or the organization requires. Update Claude Code, then retry.
type ClaudeCodeVersionTooOldErrorDetails struct {
	ErrorCode constant.ClaudeCodeVersionTooOldErrorCode `json:"error_code" default:"claude_code_version_too_old"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ClaudeCodeVersionTooOldErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ClaudeCodeVersionTooOldErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ClaudeCodeVersionTooOldErrorDetails) ImplErrorDetailsUnion() {}

// The workspace's encryption settings could not be read, so the request was
// refused. Retry in a few minutes.
type CMEKContextUnavailableErrorDetails struct {
	ErrorCode constant.CMEKContextUnavailableErrorCode `json:"error_code" default:"cmek_context_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CMEKContextUnavailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CMEKContextUnavailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CMEKContextUnavailableErrorDetails) ImplErrorDetailsUnion() {}

// The workspace's customer-managed encryption key is disabled or revoked.
// Re-enable the key, then retry.
type CMEKKeyDisabledErrorDetails struct {
	ErrorCode constant.CMEKKeyDisabledErrorCode `json:"error_code" default:"cmek_key_disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CMEKKeyDisabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CMEKKeyDisabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CMEKKeyDisabledErrorDetails) ImplErrorDetailsUnion() {}

// The key management service's network access controls blocked the request to the
// workspace's customer-managed encryption key. The key itself may be healthy.
type CMEKKeyNetworkBlockedErrorDetails struct {
	ErrorCode constant.CMEKKeyNetworkBlockedErrorCode `json:"error_code" default:"cmek_key_network_blocked"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CMEKKeyNetworkBlockedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CMEKKeyNetworkBlockedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CMEKKeyNetworkBlockedErrorDetails) ImplErrorDetailsUnion() {}

// The workspace's customer-managed encryption key could not be reached. Retry the
// request; if this continues, check that the key is available in your key
// management service, or contact Anthropic support.
type CMEKKeyUnreachableErrorDetails struct {
	ErrorCode constant.CMEKKeyUnreachableErrorCode `json:"error_code" default:"cmek_key_unreachable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CMEKKeyUnreachableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CMEKKeyUnreachableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CMEKKeyUnreachableErrorDetails) ImplErrorDetailsUnion() {}

// The organization has crossed the monthly API usage threshold of its API tier.
type EnforcedSpendLimitReachedErrorDetails struct {
	ErrorCode constant.EnforcedSpendLimitReachedErrorCode `json:"error_code" default:"enforced_spend_limit_reached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnforcedSpendLimitReachedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *EnforcedSpendLimitReachedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (EnforcedSpendLimitReachedErrorDetails) ImplErrorDetailsUnion() {}

// ErrorDetailsUnion contains all possible properties and values from
// [AvailableToolsLimitExceededErrorDetails],
// [ClaudeCodeKeyCreatorNotMemberErrorDetails],
// [ClaudeCodeVersionTooOldErrorDetails], [CMEKKeyDisabledErrorDetails],
// [CMEKKeyNetworkBlockedErrorDetails], [CompactionBlockAmbiguousErrorDetails],
// [CompactionBlockMisplacedErrorDetails], [CompactionContentMismatchErrorDetails],
// [CompactionIncompleteTurnErrorDetails],
// [CompactionNothingToSummarizeErrorDetails],
// [CompactionSignatureInvalidErrorDetails],
// [CompactionTooManyToolReferencesErrorDetails],
// [CompactionToolChangesMismatchErrorDetails], [OrganizationOnHoldErrorDetails],
// [RequestBodyEncodingInvalidErrorDetails],
// [RequestBodyEncodingUnsupportedErrorDetails],
// [ThreadAlreadyContinuedErrorDetails], [ThreadContinueInFlightErrorDetails],
// [ThreadFingerprintMismatchErrorDetails], [ThreadUnsupportedRequestErrorDetails],
// [ToolNameConflictErrorDetails], [ToolReferenceUnresolvedErrorDetails],
// [ThreadNotFoundErrorDetails], [SkillReadBlockedUndecryptableErrorDetails],
// [RequestBodyTooLargeAfterDecodingErrorDetails],
// [EnforcedSpendLimitReachedErrorDetails], [RampLimitExceededErrorDetails],
// [CMEKKeyUnreachableErrorDetails], [CompactionUnavailableErrorDetails],
// [StorageQuotaExceededErrorDetails], [FileNotDownloadableErrorDetails],
// [FileExpiredErrorDetails], [CMEKContextUnavailableErrorDetails],
// [OrgAPIKeyLimitExceededErrorDetails],
// [WorkspaceAPIKeyLimitExceededErrorDetails], [UserProfileNotFoundErrorDetails],
// [ExternalKeyInWorkspacePolicyErrorDetails], [AlreadyMemberErrorDetails],
// [CannotInviteEmailsWithUnverifiedDomainsErrorDetails],
// [EmailDomainNotAllowedErrorDetails], [InviteAlreadyExistsErrorDetails],
// [OrgMemberLimitExceededErrorDetails], [RBACGroupGrantsOwnerAccessErrorDetails],
// [RBACGroupOwnerAccessUnconfirmedErrorDetails],
// [CannotRemoveLastMemberErrorDetails], [ConcurrentMembershipChangeErrorDetails],
// [ExternalKeyValidationFailedErrorDetails],
// [ExternalKeyValidationTimedOutErrorDetails], [KeyProfileMismatchErrorDetails],
// [WorkspacePolicyViolationErrorDetails], [ExternalKeyChangedErrorDetails],
// [WorkspaceMemberLimitExceededErrorDetails],
// [OrganizationTypeNotSupportedErrorDetails],
// [BetaEnvironmentCreateDisabledErrorDetails],
// [BetaEnvironmentArchivedErrorDetails],
// [BetaInvalidConfigTypeChangeErrorDetails],
// [BetaEnvironmentUpdateDisabledErrorDetails],
// [BetaHeartbeatPreconditionFailedErrorDetails],
// [BetaInvalidMetadataErrorDetails], [BetaBillingFrozenErrorDetails],
// [BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails],
// [BetaUnsupportedCharsetErrorDetails], [BetaPluginNameTakenErrorDetails],
// [BetaSkillNameTakenErrorDetails], [BetaRegistrationPendingErrorDetails],
// [BetaScanFailedErrorDetails], [BetaScanPendingErrorDetails],
// [BetaSlackChannelSetUpSettingsRefusedErrorDetails],
// [BetaAPIKeyScopeMissingErrorDetails],
// [BetaClaudeTagWritesNotEnabledErrorDetails],
// [BetaComplianceConfigurationRestrictedErrorDetails],
// [BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails],
// [BetaEnvironmentNotFoundErrorDetails],
// [BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails],
// [BetaSlackChannelSetUpSettingsConflictErrorDetails].
//
// Use the [ErrorDetailsUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ErrorDetailsUnion struct {
	// Any of "available_tools_limit_exceeded", "claude_code_key_creator_not_member",
	// "claude_code_version_too_old", "cmek_key_disabled", "cmek_key_network_blocked",
	// "compaction_block_ambiguous", "compaction_block_misplaced",
	// "compaction_content_mismatch", "compaction_incomplete_turn",
	// "compaction_nothing_to_summarize", "compaction_signature_invalid",
	// "compaction_too_many_tool_references", "compaction_tool_changes_mismatch",
	// "organization_on_hold", "request_body_encoding_invalid",
	// "request_body_encoding_unsupported", "thread_already_continued",
	// "thread_continue_in_flight", "thread_fingerprint_mismatch",
	// "thread_unsupported_request", "tool_name_conflict", "tool_reference_unresolved",
	// "thread_not_found", "skill_read_blocked_undecryptable",
	// "request_body_too_large_after_decoding", "enforced_spend_limit_reached",
	// "ramp_limit_exceeded", "cmek_key_unreachable", "compaction_unavailable",
	// "storage_quota_exceeded", "file_not_downloadable", "file_expired",
	// "cmek_context_unavailable", "org_api_key_limit_exceeded",
	// "workspace_api_key_limit_exceeded", "user_profile_not_found",
	// "external_key_in_workspace_policy", "already_member",
	// "cannot_invite_emails_with_unverified_domains", "email_domain_not_allowed",
	// "invite_already_exists", "org_member_limit_exceeded",
	// "rbac_group_grants_owner_access", "rbac_group_owner_access_unconfirmed",
	// "cannot_remove_last_member", "concurrent_membership_change",
	// "external_key_validation_failed", "external_key_validation_timed_out",
	// "key_profile_mismatch", "workspace_policy_violation", "external_key_changed",
	// "workspace_member_limit_exceeded", "organization_type_not_supported",
	// "environment_create_disabled", "environment_archived",
	// "invalid_config_type_change", "environment_update_disabled",
	// "heartbeat_precondition_failed", "invalid_metadata", "billing_frozen",
	// "spend_limits_admin_api_writes_not_enabled", "unsupported_charset",
	// "plugin_name_taken", "skill_name_taken", "registration_pending", "scan_failed",
	// "scan_pending", "slack_channel_set_up_settings_refused",
	// "api_key_scope_missing", "claude_tag_writes_not_enabled",
	// "compliance_configuration_restricted",
	// "organization_level_scoped_api_key_required", "environment_not_found",
	// "slack_channel_set_up_environment_not_found",
	// "slack_channel_set_up_settings_conflict".
	ErrorCode string `json:"error_code"`
	// This field is from variant [OrganizationOnHoldErrorDetails].
	AppealURL string `json:"appeal_url"`
	// This field is from variant [ThreadUnsupportedRequestErrorDetails].
	Unsupported []ThreadUnsupportedFeatureName `json:"unsupported"`
	// This field is from variant [WorkspacePolicyViolationErrorDetails].
	Field WorkspacePolicyField `json:"field"`
	// This field is from variant [BetaHeartbeatPreconditionFailedErrorDetails].
	CurrentState BetaWorkLeaseState `json:"current_state"`
	PluginID     string             `json:"plugin_id"`
	// This field is from variant [BetaSkillNameTakenErrorDetails].
	SkillName       string `json:"skill_name"`
	PluginVersionID string `json:"plugin_version_id"`
	SlackChannelID  string `json:"slack_channel_id"`
	JSON            struct {
		ErrorCode       respjson.Field
		AppealURL       respjson.Field
		Unsupported     respjson.Field
		Field           respjson.Field
		CurrentState    respjson.Field
		PluginID        respjson.Field
		SkillName       respjson.Field
		PluginVersionID respjson.Field
		SlackChannelID  respjson.Field
		raw             string
	} `json:"-"`
}

// anyErrorDetails is implemented by each variant of [ErrorDetailsUnion] to add
// type safety for the return type of [ErrorDetailsUnion.AsAny]
type anyErrorDetails interface {
	ImplErrorDetailsUnion()
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ErrorDetailsUnion.AsAny().(type) {
//	case shared.AvailableToolsLimitExceededErrorDetails:
//	case shared.ClaudeCodeKeyCreatorNotMemberErrorDetails:
//	case shared.ClaudeCodeVersionTooOldErrorDetails:
//	case shared.CMEKKeyDisabledErrorDetails:
//	case shared.CMEKKeyNetworkBlockedErrorDetails:
//	case shared.CompactionBlockAmbiguousErrorDetails:
//	case shared.CompactionBlockMisplacedErrorDetails:
//	case shared.CompactionContentMismatchErrorDetails:
//	case shared.CompactionIncompleteTurnErrorDetails:
//	case shared.CompactionNothingToSummarizeErrorDetails:
//	case shared.CompactionSignatureInvalidErrorDetails:
//	case shared.CompactionTooManyToolReferencesErrorDetails:
//	case shared.CompactionToolChangesMismatchErrorDetails:
//	case shared.OrganizationOnHoldErrorDetails:
//	case shared.RequestBodyEncodingInvalidErrorDetails:
//	case shared.RequestBodyEncodingUnsupportedErrorDetails:
//	case shared.ThreadAlreadyContinuedErrorDetails:
//	case shared.ThreadContinueInFlightErrorDetails:
//	case shared.ThreadFingerprintMismatchErrorDetails:
//	case shared.ThreadUnsupportedRequestErrorDetails:
//	case shared.ToolNameConflictErrorDetails:
//	case shared.ToolReferenceUnresolvedErrorDetails:
//	case shared.ThreadNotFoundErrorDetails:
//	case shared.SkillReadBlockedUndecryptableErrorDetails:
//	case shared.RequestBodyTooLargeAfterDecodingErrorDetails:
//	case shared.EnforcedSpendLimitReachedErrorDetails:
//	case shared.RampLimitExceededErrorDetails:
//	case shared.CMEKKeyUnreachableErrorDetails:
//	case shared.CompactionUnavailableErrorDetails:
//	case shared.StorageQuotaExceededErrorDetails:
//	case shared.FileNotDownloadableErrorDetails:
//	case shared.FileExpiredErrorDetails:
//	case shared.CMEKContextUnavailableErrorDetails:
//	case shared.OrgAPIKeyLimitExceededErrorDetails:
//	case shared.WorkspaceAPIKeyLimitExceededErrorDetails:
//	case shared.UserProfileNotFoundErrorDetails:
//	case shared.ExternalKeyInWorkspacePolicyErrorDetails:
//	case shared.AlreadyMemberErrorDetails:
//	case shared.CannotInviteEmailsWithUnverifiedDomainsErrorDetails:
//	case shared.EmailDomainNotAllowedErrorDetails:
//	case shared.InviteAlreadyExistsErrorDetails:
//	case shared.OrgMemberLimitExceededErrorDetails:
//	case shared.RBACGroupGrantsOwnerAccessErrorDetails:
//	case shared.RBACGroupOwnerAccessUnconfirmedErrorDetails:
//	case shared.CannotRemoveLastMemberErrorDetails:
//	case shared.ConcurrentMembershipChangeErrorDetails:
//	case shared.ExternalKeyValidationFailedErrorDetails:
//	case shared.ExternalKeyValidationTimedOutErrorDetails:
//	case shared.KeyProfileMismatchErrorDetails:
//	case shared.WorkspacePolicyViolationErrorDetails:
//	case shared.ExternalKeyChangedErrorDetails:
//	case shared.WorkspaceMemberLimitExceededErrorDetails:
//	case shared.OrganizationTypeNotSupportedErrorDetails:
//	case shared.BetaEnvironmentCreateDisabledErrorDetails:
//	case shared.BetaEnvironmentArchivedErrorDetails:
//	case shared.BetaInvalidConfigTypeChangeErrorDetails:
//	case shared.BetaEnvironmentUpdateDisabledErrorDetails:
//	case shared.BetaHeartbeatPreconditionFailedErrorDetails:
//	case shared.BetaInvalidMetadataErrorDetails:
//	case shared.BetaBillingFrozenErrorDetails:
//	case shared.BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails:
//	case shared.BetaUnsupportedCharsetErrorDetails:
//	case shared.BetaPluginNameTakenErrorDetails:
//	case shared.BetaSkillNameTakenErrorDetails:
//	case shared.BetaRegistrationPendingErrorDetails:
//	case shared.BetaScanFailedErrorDetails:
//	case shared.BetaScanPendingErrorDetails:
//	case shared.BetaSlackChannelSetUpSettingsRefusedErrorDetails:
//	case shared.BetaAPIKeyScopeMissingErrorDetails:
//	case shared.BetaClaudeTagWritesNotEnabledErrorDetails:
//	case shared.BetaComplianceConfigurationRestrictedErrorDetails:
//	case shared.BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails:
//	case shared.BetaEnvironmentNotFoundErrorDetails:
//	case shared.BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails:
//	case shared.BetaSlackChannelSetUpSettingsConflictErrorDetails:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ErrorDetailsUnion) AsAny() anyErrorDetails {
	switch u.ErrorCode {
	case "available_tools_limit_exceeded":
		return u.AsAvailableToolsLimitExceeded()
	case "claude_code_key_creator_not_member":
		return u.AsClaudeCodeKeyCreatorNotMember()
	case "claude_code_version_too_old":
		return u.AsClaudeCodeVersionTooOld()
	case "cmek_key_disabled":
		return u.AsCMEKKeyDisabled()
	case "cmek_key_network_blocked":
		return u.AsCMEKKeyNetworkBlocked()
	case "compaction_block_ambiguous":
		return u.AsCompactionBlockAmbiguous()
	case "compaction_block_misplaced":
		return u.AsCompactionBlockMisplaced()
	case "compaction_content_mismatch":
		return u.AsCompactionContentMismatch()
	case "compaction_incomplete_turn":
		return u.AsCompactionIncompleteTurn()
	case "compaction_nothing_to_summarize":
		return u.AsCompactionNothingToSummarize()
	case "compaction_signature_invalid":
		return u.AsCompactionSignatureInvalid()
	case "compaction_too_many_tool_references":
		return u.AsCompactionTooManyToolReferences()
	case "compaction_tool_changes_mismatch":
		return u.AsCompactionToolChangesMismatch()
	case "organization_on_hold":
		return u.AsOrganizationOnHold()
	case "request_body_encoding_invalid":
		return u.AsRequestBodyEncodingInvalid()
	case "request_body_encoding_unsupported":
		return u.AsRequestBodyEncodingUnsupported()
	case "thread_already_continued":
		return u.AsThreadAlreadyContinued()
	case "thread_continue_in_flight":
		return u.AsThreadContinueInFlight()
	case "thread_fingerprint_mismatch":
		return u.AsThreadFingerprintMismatch()
	case "thread_unsupported_request":
		return u.AsThreadUnsupportedRequest()
	case "tool_name_conflict":
		return u.AsToolNameConflict()
	case "tool_reference_unresolved":
		return u.AsToolReferenceUnresolved()
	case "thread_not_found":
		return u.AsThreadNotFound()
	case "skill_read_blocked_undecryptable":
		return u.AsSkillReadBlockedUndecryptable()
	case "request_body_too_large_after_decoding":
		return u.AsRequestBodyTooLargeAfterDecoding()
	case "enforced_spend_limit_reached":
		return u.AsEnforcedSpendLimitReached()
	case "ramp_limit_exceeded":
		return u.AsRampLimitExceeded()
	case "cmek_key_unreachable":
		return u.AsCMEKKeyUnreachable()
	case "compaction_unavailable":
		return u.AsCompactionUnavailable()
	case "storage_quota_exceeded":
		return u.AsStorageQuotaExceeded()
	case "file_not_downloadable":
		return u.AsFileNotDownloadable()
	case "file_expired":
		return u.AsFileExpired()
	case "cmek_context_unavailable":
		return u.AsCMEKContextUnavailable()
	case "org_api_key_limit_exceeded":
		return u.AsOrgAPIKeyLimitExceeded()
	case "workspace_api_key_limit_exceeded":
		return u.AsWorkspaceAPIKeyLimitExceeded()
	case "user_profile_not_found":
		return u.AsUserProfileNotFound()
	case "external_key_in_workspace_policy":
		return u.AsExternalKeyInWorkspacePolicy()
	case "already_member":
		return u.AsAlreadyMember()
	case "cannot_invite_emails_with_unverified_domains":
		return u.AsCannotInviteEmailsWithUnverifiedDomains()
	case "email_domain_not_allowed":
		return u.AsEmailDomainNotAllowed()
	case "invite_already_exists":
		return u.AsInviteAlreadyExists()
	case "org_member_limit_exceeded":
		return u.AsOrgMemberLimitExceeded()
	case "rbac_group_grants_owner_access":
		return u.AsRBACGroupGrantsOwnerAccess()
	case "rbac_group_owner_access_unconfirmed":
		return u.AsRBACGroupOwnerAccessUnconfirmed()
	case "cannot_remove_last_member":
		return u.AsCannotRemoveLastMember()
	case "concurrent_membership_change":
		return u.AsConcurrentMembershipChange()
	case "external_key_validation_failed":
		return u.AsExternalKeyValidationFailed()
	case "external_key_validation_timed_out":
		return u.AsExternalKeyValidationTimedOut()
	case "key_profile_mismatch":
		return u.AsKeyProfileMismatch()
	case "workspace_policy_violation":
		return u.AsWorkspacePolicyViolation()
	case "external_key_changed":
		return u.AsExternalKeyChanged()
	case "workspace_member_limit_exceeded":
		return u.AsWorkspaceMemberLimitExceeded()
	case "organization_type_not_supported":
		return u.AsOrganizationTypeNotSupported()
	case "environment_create_disabled":
		return u.AsEnvironmentCreateDisabled()
	case "environment_archived":
		return u.AsEnvironmentArchived()
	case "invalid_config_type_change":
		return u.AsInvalidConfigTypeChange()
	case "environment_update_disabled":
		return u.AsEnvironmentUpdateDisabled()
	case "heartbeat_precondition_failed":
		return u.AsHeartbeatPreconditionFailed()
	case "invalid_metadata":
		return u.AsInvalidMetadata()
	case "billing_frozen":
		return u.AsBillingFrozen()
	case "spend_limits_admin_api_writes_not_enabled":
		return u.AsSpendLimitsAdminAPIWritesNotEnabled()
	case "unsupported_charset":
		return u.AsUnsupportedCharset()
	case "plugin_name_taken":
		return u.AsPluginNameTaken()
	case "skill_name_taken":
		return u.AsSkillNameTaken()
	case "registration_pending":
		return u.AsRegistrationPending()
	case "scan_failed":
		return u.AsScanFailed()
	case "scan_pending":
		return u.AsScanPending()
	case "slack_channel_set_up_settings_refused":
		return u.AsSlackChannelSetUpSettingsRefused()
	case "api_key_scope_missing":
		return u.AsAPIKeyScopeMissing()
	case "claude_tag_writes_not_enabled":
		return u.AsClaudeTagWritesNotEnabled()
	case "compliance_configuration_restricted":
		return u.AsComplianceConfigurationRestricted()
	case "organization_level_scoped_api_key_required":
		return u.AsOrganizationLevelScopedAPIKeyRequired()
	case "environment_not_found":
		return u.AsEnvironmentNotFound()
	case "slack_channel_set_up_environment_not_found":
		return u.AsSlackChannelSetUpEnvironmentNotFound()
	case "slack_channel_set_up_settings_conflict":
		return u.AsSlackChannelSetUpSettingsConflict()
	}
	return nil
}

func (u ErrorDetailsUnion) AsAvailableToolsLimitExceeded() (v AvailableToolsLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsClaudeCodeKeyCreatorNotMember() (v ClaudeCodeKeyCreatorNotMemberErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsClaudeCodeVersionTooOld() (v ClaudeCodeVersionTooOldErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKKeyDisabled() (v CMEKKeyDisabledErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKKeyNetworkBlocked() (v CMEKKeyNetworkBlockedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionBlockAmbiguous() (v CompactionBlockAmbiguousErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionBlockMisplaced() (v CompactionBlockMisplacedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionContentMismatch() (v CompactionContentMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionIncompleteTurn() (v CompactionIncompleteTurnErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionNothingToSummarize() (v CompactionNothingToSummarizeErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionSignatureInvalid() (v CompactionSignatureInvalidErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionTooManyToolReferences() (v CompactionTooManyToolReferencesErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionToolChangesMismatch() (v CompactionToolChangesMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrganizationOnHold() (v OrganizationOnHoldErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRequestBodyEncodingInvalid() (v RequestBodyEncodingInvalidErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRequestBodyEncodingUnsupported() (v RequestBodyEncodingUnsupportedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsThreadAlreadyContinued() (v ThreadAlreadyContinuedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsThreadContinueInFlight() (v ThreadContinueInFlightErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsThreadFingerprintMismatch() (v ThreadFingerprintMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsThreadUnsupportedRequest() (v ThreadUnsupportedRequestErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsToolNameConflict() (v ToolNameConflictErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsToolReferenceUnresolved() (v ToolReferenceUnresolvedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsThreadNotFound() (v ThreadNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSkillReadBlockedUndecryptable() (v SkillReadBlockedUndecryptableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRequestBodyTooLargeAfterDecoding() (v RequestBodyTooLargeAfterDecodingErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEnforcedSpendLimitReached() (v EnforcedSpendLimitReachedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRampLimitExceeded() (v RampLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKKeyUnreachable() (v CMEKKeyUnreachableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCompactionUnavailable() (v CompactionUnavailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageQuotaExceeded() (v StorageQuotaExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsFileNotDownloadable() (v FileNotDownloadableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsFileExpired() (v FileExpiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKContextUnavailable() (v CMEKContextUnavailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrgAPIKeyLimitExceeded() (v OrgAPIKeyLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsWorkspaceAPIKeyLimitExceeded() (v WorkspaceAPIKeyLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsUserProfileNotFound() (v UserProfileNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyInWorkspacePolicy() (v ExternalKeyInWorkspacePolicyErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsAlreadyMember() (v AlreadyMemberErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotInviteEmailsWithUnverifiedDomains() (v CannotInviteEmailsWithUnverifiedDomainsErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEmailDomainNotAllowed() (v EmailDomainNotAllowedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsInviteAlreadyExists() (v InviteAlreadyExistsErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrgMemberLimitExceeded() (v OrgMemberLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACGroupGrantsOwnerAccess() (v RBACGroupGrantsOwnerAccessErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACGroupOwnerAccessUnconfirmed() (v RBACGroupOwnerAccessUnconfirmedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotRemoveLastMember() (v CannotRemoveLastMemberErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsConcurrentMembershipChange() (v ConcurrentMembershipChangeErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyValidationFailed() (v ExternalKeyValidationFailedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyValidationTimedOut() (v ExternalKeyValidationTimedOutErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsKeyProfileMismatch() (v KeyProfileMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsWorkspacePolicyViolation() (v WorkspacePolicyViolationErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyChanged() (v ExternalKeyChangedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsWorkspaceMemberLimitExceeded() (v WorkspaceMemberLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrganizationTypeNotSupported() (v OrganizationTypeNotSupportedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEnvironmentCreateDisabled() (v BetaEnvironmentCreateDisabledErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEnvironmentArchived() (v BetaEnvironmentArchivedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsInvalidConfigTypeChange() (v BetaInvalidConfigTypeChangeErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEnvironmentUpdateDisabled() (v BetaEnvironmentUpdateDisabledErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsHeartbeatPreconditionFailed() (v BetaHeartbeatPreconditionFailedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsInvalidMetadata() (v BetaInvalidMetadataErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsBillingFrozen() (v BetaBillingFrozenErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSpendLimitsAdminAPIWritesNotEnabled() (v BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsUnsupportedCharset() (v BetaUnsupportedCharsetErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsPluginNameTaken() (v BetaPluginNameTakenErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSkillNameTaken() (v BetaSkillNameTakenErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRegistrationPending() (v BetaRegistrationPendingErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsScanFailed() (v BetaScanFailedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsScanPending() (v BetaScanPendingErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSlackChannelSetUpSettingsRefused() (v BetaSlackChannelSetUpSettingsRefusedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsAPIKeyScopeMissing() (v BetaAPIKeyScopeMissingErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsClaudeTagWritesNotEnabled() (v BetaClaudeTagWritesNotEnabledErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsComplianceConfigurationRestricted() (v BetaComplianceConfigurationRestrictedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrganizationLevelScopedAPIKeyRequired() (v BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsEnvironmentNotFound() (v BetaEnvironmentNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSlackChannelSetUpEnvironmentNotFound() (v BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSlackChannelSetUpSettingsConflict() (v BetaSlackChannelSetUpSettingsConflictErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ErrorDetailsUnion) RawJSON() string { return u.JSON.raw }

func (r *ErrorDetailsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ErrorObjectUnion contains all possible properties and values from
// [InvalidRequestError], [AuthenticationError], [BillingError], [PermissionError],
// [NotFoundError], [RateLimitError], [GatewayTimeoutError], [APIErrorObject],
// [OverloadedError].
//
// Use the [ErrorObjectUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ErrorObjectUnion struct {
	Message string `json:"message"`
	// Any of "invalid_request_error", "authentication_error", "billing_error",
	// "permission_error", "not_found_error", "rate_limit_error", "timeout_error",
	// "api_error", "overloaded_error".
	Type string `json:"type"`
	JSON struct {
		Message respjson.Field
		Type    respjson.Field
		raw     string
	} `json:"-"`
}

// anyErrorObject is implemented by each variant of [ErrorObjectUnion] to add type
// safety for the return type of [ErrorObjectUnion.AsAny]
type anyErrorObject interface {
	ImplErrorObjectUnion()
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ErrorObjectUnion.AsAny().(type) {
//	case shared.InvalidRequestError:
//	case shared.AuthenticationError:
//	case shared.BillingError:
//	case shared.PermissionError:
//	case shared.NotFoundError:
//	case shared.RateLimitError:
//	case shared.GatewayTimeoutError:
//	case shared.APIErrorObject:
//	case shared.OverloadedError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ErrorObjectUnion) AsAny() anyErrorObject {
	switch u.Type {
	case "invalid_request_error":
		return u.AsInvalidRequestError()
	case "authentication_error":
		return u.AsAuthenticationError()
	case "billing_error":
		return u.AsBillingError()
	case "permission_error":
		return u.AsPermissionError()
	case "not_found_error":
		return u.AsNotFoundError()
	case "rate_limit_error":
		return u.AsRateLimitError()
	case "timeout_error":
		return u.AsTimeoutError()
	case "api_error":
		return u.AsAPIError()
	case "overloaded_error":
		return u.AsOverloadedError()
	}
	return nil
}

func (u ErrorObjectUnion) AsInvalidRequestError() (v InvalidRequestError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsAuthenticationError() (v AuthenticationError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsBillingError() (v BillingError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsPermissionError() (v PermissionError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsNotFoundError() (v NotFoundError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsRateLimitError() (v RateLimitError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsTimeoutError() (v GatewayTimeoutError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsAPIError() (v APIErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsOverloadedError() (v OverloadedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ErrorObjectUnion) RawJSON() string { return u.JSON.raw }

func (r *ErrorObjectUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ErrorResponse struct {
	Error     ErrorObjectUnion `json:"error" api:"required"`
	RequestID string           `json:"request_id" api:"required"`
	Type      constant.Error   `json:"type" default:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		RequestID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorResponse) RawJSON() string { return r.JSON.raw }
func (r *ErrorResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ErrorType string

const (
	ErrorTypeInvalidRequestError           ErrorType = "invalid_request_error"
	ErrorTypeAuthenticationError           ErrorType = "authentication_error"
	ErrorTypeBillingError                  ErrorType = "billing_error"
	ErrorTypePermissionError               ErrorType = "permission_error"
	ErrorTypeNotFoundError                 ErrorType = "not_found_error"
	ErrorTypeRateLimitError                ErrorType = "rate_limit_error"
	ErrorTypeTimeoutError                  ErrorType = "timeout_error"
	ErrorTypeAPIError                      ErrorType = "api_error"
	ErrorTypeOverloadedError               ErrorType = "overloaded_error"
	ErrorTypeMemoryPreconditionFailedError ErrorType = "memory_precondition_failed_error"
	ErrorTypeMemoryPathConflictError       ErrorType = "memory_path_conflict_error"
	ErrorTypeConflictError                 ErrorType = "conflict_error"
)

type GatewayTimeoutError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.TimeoutError `json:"type" default:"timeout_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r GatewayTimeoutError) RawJSON() string { return r.JSON.raw }
func (r *GatewayTimeoutError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (GatewayTimeoutError) ImplErrorObjectUnion() {}

type InvalidRequestError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.InvalidRequestError `json:"type" default:"invalid_request_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvalidRequestError) RawJSON() string { return r.JSON.raw }
func (r *InvalidRequestError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (InvalidRequestError) ImplErrorObjectUnion() {}

type NotFoundError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.NotFoundError `json:"type" default:"not_found_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r NotFoundError) RawJSON() string { return r.JSON.raw }
func (r *NotFoundError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (NotFoundError) ImplErrorObjectUnion() {}

// Anthropic has disabled the organization for trust and safety reasons. An
// organization disabled for any other reason, such as billing, does not get this
// code. `appeal_url`, when present, is where the organization can appeal.
type OrganizationOnHoldErrorDetails struct {
	ErrorCode constant.OrganizationOnHoldErrorCode `json:"error_code" default:"organization_on_hold"`
	AppealURL string                               `json:"appeal_url" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		AppealURL   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationOnHoldErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *OrganizationOnHoldErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OrganizationOnHoldErrorDetails) ImplErrorDetailsUnion() {}

// The endpoint does not support the organization's type. Each endpoint's
// description says which organization types it supports.
type OrganizationTypeNotSupportedErrorDetails struct {
	ErrorCode constant.OrganizationTypeNotSupportedErrorCode `json:"error_code" default:"organization_type_not_supported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationTypeNotSupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *OrganizationTypeNotSupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OrganizationTypeNotSupportedErrorDetails) ImplErrorDetailsUnion() {}

type OverloadedError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.OverloadedError `json:"type" default:"overloaded_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OverloadedError) RawJSON() string { return r.JSON.raw }
func (r *OverloadedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OverloadedError) ImplErrorObjectUnion() {}

type PermissionError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.PermissionError `json:"type" default:"permission_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PermissionError) RawJSON() string { return r.JSON.raw }
func (r *PermissionError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (PermissionError) ImplErrorObjectUnion() {}

// The request would exceed a ramp limit, which caps how fast usage can grow. Retry
// after the time in the `retry-after` header.
type RampLimitExceededErrorDetails struct {
	ErrorCode constant.RampLimitExceededErrorCode `json:"error_code" default:"ramp_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RampLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RampLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RampLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

type RateLimitError struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.RateLimitError `json:"type" default:"rate_limit_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RateLimitError) RawJSON() string { return r.JSON.raw }
func (r *RateLimitError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RateLimitError) ImplErrorObjectUnion() {}

// The request body is not valid for its `Content-Encoding`.
type RequestBodyEncodingInvalidErrorDetails struct {
	ErrorCode constant.RequestBodyEncodingInvalidErrorCode `json:"error_code" default:"request_body_encoding_invalid"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RequestBodyEncodingInvalidErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RequestBodyEncodingInvalidErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RequestBodyEncodingInvalidErrorDetails) ImplErrorDetailsUnion() {}

// The request's `Content-Encoding` is not supported.
type RequestBodyEncodingUnsupportedErrorDetails struct {
	ErrorCode constant.RequestBodyEncodingUnsupportedErrorCode `json:"error_code" default:"request_body_encoding_unsupported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RequestBodyEncodingUnsupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RequestBodyEncodingUnsupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RequestBodyEncodingUnsupportedErrorDetails) ImplErrorDetailsUnion() {}

// The compressed request body exceeds the request size limit once decoded.
type RequestBodyTooLargeAfterDecodingErrorDetails struct {
	ErrorCode constant.RequestBodyTooLargeAfterDecodingErrorCode `json:"error_code" default:"request_body_too_large_after_decoding"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RequestBodyTooLargeAfterDecodingErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RequestBodyTooLargeAfterDecodingErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RequestBodyTooLargeAfterDecodingErrorDetails) ImplErrorDetailsUnion() {}

// The tools available after a message exceed a limit on their number or their
// size.
type AvailableToolsLimitExceededErrorDetails struct {
	ErrorCode constant.AvailableToolsLimitExceededErrorCode `json:"error_code" default:"available_tools_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AvailableToolsLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *AvailableToolsLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (AvailableToolsLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// `messages` holds more than one `compaction` block with content. Send only the
// block the conversation returned last.
type CompactionBlockAmbiguousErrorDetails struct {
	ErrorCode constant.CompactionBlockAmbiguousErrorCode `json:"error_code" default:"compaction_block_ambiguous"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionBlockAmbiguousErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionBlockAmbiguousErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionBlockAmbiguousErrorDetails) ImplErrorDetailsUnion() {}

// The `compaction` block is not the first block of the conversation, or is in a
// `system` message.
type CompactionBlockMisplacedErrorDetails struct {
	ErrorCode constant.CompactionBlockMisplacedErrorCode `json:"error_code" default:"compaction_block_misplaced"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionBlockMisplacedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionBlockMisplacedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionBlockMisplacedErrorDetails) ImplErrorDetailsUnion() {}

// The `compaction` block's `content` differs from the text its `signature`
// recorded.
type CompactionContentMismatchErrorDetails struct {
	ErrorCode constant.CompactionContentMismatchErrorCode `json:"error_code" default:"compaction_content_mismatch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionContentMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionContentMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionContentMismatchErrorDetails) ImplErrorDetailsUnion() {}

// The compaction request ends on an `assistant` turn that is not finished.
type CompactionIncompleteTurnErrorDetails struct {
	ErrorCode constant.CompactionIncompleteTurnErrorCode `json:"error_code" default:"compaction_incomplete_turn"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionIncompleteTurnErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionIncompleteTurnErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionIncompleteTurnErrorDetails) ImplErrorDetailsUnion() {}

// `messages` holds no `user` or `assistant` content to summarize or adopt.
type CompactionNothingToSummarizeErrorDetails struct {
	ErrorCode constant.CompactionNothingToSummarizeErrorCode `json:"error_code" default:"compaction_nothing_to_summarize"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionNothingToSummarizeErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionNothingToSummarizeErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionNothingToSummarizeErrorDetails) ImplErrorDetailsUnion() {}

// The `compaction` block's `signature` does not verify.
type CompactionSignatureInvalidErrorDetails struct {
	ErrorCode constant.CompactionSignatureInvalidErrorCode `json:"error_code" default:"compaction_signature_invalid"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionSignatureInvalidErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionSignatureInvalidErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionSignatureInvalidErrorDetails) ImplErrorDetailsUnion() {}

// The conversation references more distinct tools than a `compaction` block can
// record.
type CompactionTooManyToolReferencesErrorDetails struct {
	ErrorCode constant.CompactionTooManyToolReferencesErrorCode `json:"error_code" default:"compaction_too_many_tool_references"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionTooManyToolReferencesErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionTooManyToolReferencesErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionTooManyToolReferencesErrorDetails) ImplErrorDetailsUnion() {}

// The `compaction` block's `tool_changes` differs from what was returned with the
// block.
type CompactionToolChangesMismatchErrorDetails struct {
	ErrorCode constant.CompactionToolChangesMismatchErrorCode `json:"error_code" default:"compaction_tool_changes_mismatch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionToolChangesMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionToolChangesMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionToolChangesMismatchErrorDetails) ImplErrorDetailsUnion() {}

// The response that `previous_message_id` names has already been continued, and
// each response can be continued only once. Continue from the Thread's most recent
// response `id`.
type ThreadAlreadyContinuedErrorDetails struct {
	ErrorCode constant.ThreadAlreadyContinuedErrorCode `json:"error_code" default:"thread_already_continued"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThreadAlreadyContinuedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ThreadAlreadyContinuedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ThreadAlreadyContinuedErrorDetails) ImplErrorDetailsUnion() {}

// Another request is still continuing from the same `previous_message_id`. Wait
// for it to complete and continue from its response `id`, or replay the full
// conversation with `thread: {"type": "create"}`.
type ThreadContinueInFlightErrorDetails struct {
	ErrorCode constant.ThreadContinueInFlightErrorCode `json:"error_code" default:"thread_continue_in_flight"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThreadContinueInFlightErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ThreadContinueInFlightErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ThreadContinueInFlightErrorDetails) ImplErrorDetailsUnion() {}

// A field that must stay the same across a Thread's turns differs from the value
// stored with the Thread; the message names each such field. Send the stored
// values, or replay the full conversation with `thread: {"type": "create"}` to
// change them.
type ThreadFingerprintMismatchErrorDetails struct {
	ErrorCode constant.ThreadFingerprintMismatchErrorCode `json:"error_code" default:"thread_fingerprint_mismatch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThreadFingerprintMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ThreadFingerprintMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ThreadFingerprintMismatchErrorDetails) ImplErrorDetailsUnion() {}

// The request combines `thread` with a request field or an organization setting
// that Threads do not support; `unsupported`, when present, names it. Send the
// full conversation without `thread`, or change that field.
type ThreadUnsupportedRequestErrorDetails struct {
	ErrorCode constant.ThreadUnsupportedRequestErrorCode `json:"error_code" default:"thread_unsupported_request"`
	// Each entry names a request field or an organization setting that Threads do not
	// support. On a value you do not recognize, send the full conversation without
	// `thread`.
	Unsupported []ThreadUnsupportedFeatureName `json:"unsupported" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Unsupported respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThreadUnsupportedRequestErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ThreadUnsupportedRequestErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ThreadUnsupportedRequestErrorDetails) ImplErrorDetailsUnion() {}

type ThreadUnsupportedFeatureName string

const (
	ThreadUnsupportedFeatureNameCacheControl      ThreadUnsupportedFeatureName = "cache_control"
	ThreadUnsupportedFeatureNameContextManagement ThreadUnsupportedFeatureName = "context_management"
	ThreadUnsupportedFeatureNameFallbacks         ThreadUnsupportedFeatureName = "fallbacks"
	ThreadUnsupportedFeatureNameHipaa             ThreadUnsupportedFeatureName = "hipaa"
	ThreadUnsupportedFeatureNameInferenceHooks    ThreadUnsupportedFeatureName = "inference_hooks"
	ThreadUnsupportedFeatureNameThinking          ThreadUnsupportedFeatureName = "thinking"
	ThreadUnsupportedFeatureNameZdr               ThreadUnsupportedFeatureName = "zdr"
)

// A tool made available in `messages` takes a name that an available tool of
// another kind, or of another MCP server, already holds.
type ToolNameConflictErrorDetails struct {
	ErrorCode constant.ToolNameConflictErrorCode `json:"error_code" default:"tool_name_conflict"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolNameConflictErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ToolNameConflictErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ToolNameConflictErrorDetails) ImplErrorDetailsUnion() {}

// A tool change or MCP tool listing references a tool or MCP server that is not
// defined at that position in the conversation.
type ToolReferenceUnresolvedErrorDetails struct {
	ErrorCode constant.ToolReferenceUnresolvedErrorCode `json:"error_code" default:"tool_reference_unresolved"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolReferenceUnresolvedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ToolReferenceUnresolvedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ToolReferenceUnresolvedErrorDetails) ImplErrorDetailsUnion() {}

// No usable Thread state was found for `previous_message_id`. Replay the full
// conversation with `thread: {"type": "create"}` to start a new Thread.
type ThreadNotFoundErrorDetails struct {
	ErrorCode constant.ThreadNotFoundErrorCode `json:"error_code" default:"thread_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThreadNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ThreadNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ThreadNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// A custom skill the request names is stored under an encryption key state the
// workspace cannot currently decrypt.
type SkillReadBlockedUndecryptableErrorDetails struct {
	ErrorCode constant.SkillReadBlockedUndecryptableErrorCode `json:"error_code" default:"skill_read_blocked_undecryptable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SkillReadBlockedUndecryptableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *SkillReadBlockedUndecryptableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (SkillReadBlockedUndecryptableErrorDetails) ImplErrorDetailsUnion() {}

// The compaction could not be completed for a transient server reason. Retry the
// request.
type CompactionUnavailableErrorDetails struct {
	ErrorCode constant.CompactionUnavailableErrorCode `json:"error_code" default:"compaction_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CompactionUnavailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CompactionUnavailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CompactionUnavailableErrorDetails) ImplErrorDetailsUnion() {}

// The organization's file storage quota is exhausted. Delete files to free up
// space, then retry the upload.
type StorageQuotaExceededErrorDetails struct {
	ErrorCode constant.StorageQuotaExceededErrorCode `json:"error_code" default:"storage_quota_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r StorageQuotaExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *StorageQuotaExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (StorageQuotaExceededErrorDetails) ImplErrorDetailsUnion() {}

// The file has `downloadable: false`. Only files generated by a tool, such as the
// code execution tool, can be downloaded.
type FileNotDownloadableErrorDetails struct {
	ErrorCode constant.FileNotDownloadableErrorCode `json:"error_code" default:"file_not_downloadable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FileNotDownloadableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *FileNotDownloadableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (FileNotDownloadableErrorDetails) ImplErrorDetailsUnion() {}

// The file's `expires_at` has passed, so its content is no longer available.
type FileExpiredErrorDetails struct {
	ErrorCode constant.FileExpiredErrorCode `json:"error_code" default:"file_expired"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FileExpiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *FileExpiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (FileExpiredErrorDetails) ImplErrorDetailsUnion() {}

// The organization has reached its limit of active API keys.
type OrgAPIKeyLimitExceededErrorDetails struct {
	ErrorCode constant.OrgAPIKeyLimitExceededErrorCode `json:"error_code" default:"org_api_key_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgAPIKeyLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *OrgAPIKeyLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OrgAPIKeyLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// The workspace has reached its limit of active API keys.
type WorkspaceAPIKeyLimitExceededErrorDetails struct {
	ErrorCode constant.WorkspaceAPIKeyLimitExceededErrorCode `json:"error_code" default:"workspace_api_key_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceAPIKeyLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceAPIKeyLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (WorkspaceAPIKeyLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// The request depends on a user profile that does not exist in the organization or
// is no longer active. The request either names this user profile by ID, or names
// a resource, such as an External Key, that is bound to it.
type UserProfileNotFoundErrorDetails struct {
	ErrorCode constant.UserProfileNotFoundErrorCode `json:"error_code" default:"user_profile_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (UserProfileNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// The External Key can't be deleted or reconfigured while the organization's
// Workspace Policy lists it as an allowed key.
type ExternalKeyInWorkspacePolicyErrorDetails struct {
	ErrorCode constant.ExternalKeyInWorkspacePolicyErrorCode `json:"error_code" default:"external_key_in_workspace_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyInWorkspacePolicyErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyInWorkspacePolicyErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ExternalKeyInWorkspacePolicyErrorDetails) ImplErrorDetailsUnion() {}

// The email address already belongs to a user in the organization. The message
// names the user.
type AlreadyMemberErrorDetails struct {
	ErrorCode constant.AlreadyMemberErrorCode `json:"error_code" default:"already_member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AlreadyMemberErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *AlreadyMemberErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (AlreadyMemberErrorDetails) ImplErrorDetailsUnion() {}

// The email address's domain is not one of the verified domains of the
// organization's parent organization.
type CannotInviteEmailsWithUnverifiedDomainsErrorDetails struct {
	ErrorCode constant.CannotInviteEmailsWithUnverifiedDomainsErrorCode `json:"error_code" default:"cannot_invite_emails_with_unverified_domains"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CannotInviteEmailsWithUnverifiedDomainsErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CannotInviteEmailsWithUnverifiedDomainsErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CannotInviteEmailsWithUnverifiedDomainsErrorDetails) ImplErrorDetailsUnion() {}

// The organization restricts invites to an allowed list of email domains, and the
// list is empty or does not include the email address's domain.
type EmailDomainNotAllowedErrorDetails struct {
	ErrorCode constant.EmailDomainNotAllowedErrorCode `json:"error_code" default:"email_domain_not_allowed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EmailDomainNotAllowedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *EmailDomainNotAllowedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (EmailDomainNotAllowedErrorDetails) ImplErrorDetailsUnion() {}

// A pending invite already exists for the email address. The message names the
// invite.
type InviteAlreadyExistsErrorDetails struct {
	ErrorCode constant.InviteAlreadyExistsErrorCode `json:"error_code" default:"invite_already_exists"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InviteAlreadyExistsErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *InviteAlreadyExistsErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (InviteAlreadyExistsErrorDetails) ImplErrorDetailsUnion() {}

// The organization has reached the member limit that Anthropic set for it. Pending
// invites count toward the limit.
type OrgMemberLimitExceededErrorDetails struct {
	ErrorCode constant.OrgMemberLimitExceededErrorCode `json:"error_code" default:"org_member_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrgMemberLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *OrgMemberLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OrgMemberLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// A role attached to an RBAC Group in the request grants owner access, so this API
// does not add members to that group (directly or through an invite), remove its
// members, or delete it. Make the change in the organization's admin settings in
// Claude.
type RBACGroupGrantsOwnerAccessErrorDetails struct {
	ErrorCode constant.RBACGroupGrantsOwnerAccessErrorCode `json:"error_code" default:"rbac_group_grants_owner_access"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACGroupGrantsOwnerAccessErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RBACGroupGrantsOwnerAccessErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RBACGroupGrantsOwnerAccessErrorDetails) ImplErrorDetailsUnion() {}

// The API could not finish checking whether an RBAC Group in the request grants
// owner access, so it made no change. Retrying the same request will not help.
// Make the change in the organization's admin settings in Claude.
type RBACGroupOwnerAccessUnconfirmedErrorDetails struct {
	ErrorCode constant.RBACGroupOwnerAccessUnconfirmedErrorCode `json:"error_code" default:"rbac_group_owner_access_unconfirmed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACGroupOwnerAccessUnconfirmedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RBACGroupOwnerAccessUnconfirmedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RBACGroupOwnerAccessUnconfirmedErrorDetails) ImplErrorDetailsUnion() {}

// The user is the organization's only member. Add another member before removing
// this one.
type CannotRemoveLastMemberErrorDetails struct {
	ErrorCode constant.CannotRemoveLastMemberErrorCode `json:"error_code" default:"cannot_remove_last_member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CannotRemoveLastMemberErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CannotRemoveLastMemberErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (CannotRemoveLastMemberErrorDetails) ImplErrorDetailsUnion() {}

// Another change to the organization's members is in progress. Retry the request
// once it completes.
type ConcurrentMembershipChangeErrorDetails struct {
	ErrorCode constant.ConcurrentMembershipChangeErrorCode `json:"error_code" default:"concurrent_membership_change"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ConcurrentMembershipChangeErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ConcurrentMembershipChangeErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ConcurrentMembershipChangeErrorDetails) ImplErrorDetailsUnion() {}

// The External Key failed Anthropic's validation, so the request was refused. The
// message gives the reason, if there is one.
type ExternalKeyValidationFailedErrorDetails struct {
	ErrorCode constant.ExternalKeyValidationFailedErrorCode `json:"error_code" default:"external_key_validation_failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyValidationFailedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyValidationFailedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ExternalKeyValidationFailedErrorDetails) ImplErrorDetailsUnion() {}

// Anthropic's validation of the External Key did not finish in time, so the
// request was refused. Wait a moment, then send the request again.
type ExternalKeyValidationTimedOutErrorDetails struct {
	ErrorCode constant.ExternalKeyValidationTimedOutErrorCode `json:"error_code" default:"external_key_validation_timed_out"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyValidationTimedOutErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyValidationTimedOutErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ExternalKeyValidationTimedOutErrorDetails) ImplErrorDetailsUnion() {}

// The External Key is bound to a user profile, and the Workspace does not have
// that user profile. A Workspace can use only a key bound to its own user profile
// or a key that isn't bound to one.
type KeyProfileMismatchErrorDetails struct {
	ErrorCode constant.KeyProfileMismatchErrorCode `json:"error_code" default:"key_profile_mismatch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r KeyProfileMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *KeyProfileMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (KeyProfileMismatchErrorDetails) ImplErrorDetailsUnion() {}

// The request breaks the organization's Workspace Policy. The message says which
// rule it broke; an organization admin can see and change the policy in the
// Console.
type WorkspacePolicyViolationErrorDetails struct {
	ErrorCode constant.WorkspacePolicyViolationErrorCode `json:"error_code" default:"workspace_policy_violation"`
	// The request field the broken rule applies to.
	//
	// Any of "external_key_id".
	Field WorkspacePolicyField `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspacePolicyViolationErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *WorkspacePolicyViolationErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (WorkspacePolicyViolationErrorDetails) ImplErrorDetailsUnion() {}

type WorkspacePolicyField string

const (
	WorkspacePolicyFieldExternalKeyID WorkspacePolicyField = "external_key_id"
)

// The external key changed while the request was validating it. Send the request
// again.
type ExternalKeyChangedErrorDetails struct {
	ErrorCode constant.ExternalKeyChangedErrorCode `json:"error_code" default:"external_key_changed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyChangedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyChangedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ExternalKeyChangedErrorDetails) ImplErrorDetailsUnion() {}

// The workspace has reached the member limit that Anthropic set for it.
type WorkspaceMemberLimitExceededErrorDetails struct {
	ErrorCode constant.WorkspaceMemberLimitExceededErrorCode `json:"error_code" default:"workspace_member_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceMemberLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceMemberLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (WorkspaceMemberLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// Anthropic has turned off creating environments for this organization. Existing
// environments still work.
type BetaEnvironmentCreateDisabledErrorDetails struct {
	ErrorCode constant.BetaEnvironmentCreateDisabledErrorCode `json:"error_code" default:"environment_create_disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironmentCreateDisabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironmentCreateDisabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaEnvironmentCreateDisabledErrorDetails) ImplErrorDetailsUnion() {}

// The environment is archived, so it cannot be updated.
type BetaEnvironmentArchivedErrorDetails struct {
	ErrorCode constant.BetaEnvironmentArchivedErrorCode `json:"error_code" default:"environment_archived"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironmentArchivedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironmentArchivedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaEnvironmentArchivedErrorDetails) ImplErrorDetailsUnion() {}

// `config.type` differs from the environment's stored type, which cannot change
// after creation.
type BetaInvalidConfigTypeChangeErrorDetails struct {
	ErrorCode constant.BetaInvalidConfigTypeChangeErrorCode `json:"error_code" default:"invalid_config_type_change"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInvalidConfigTypeChangeErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaInvalidConfigTypeChangeErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaInvalidConfigTypeChangeErrorDetails) ImplErrorDetailsUnion() {}

// Anthropic has turned off updating environments for this organization. Existing
// environments still work as they are.
type BetaEnvironmentUpdateDisabledErrorDetails struct {
	ErrorCode constant.BetaEnvironmentUpdateDisabledErrorCode `json:"error_code" default:"environment_update_disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironmentUpdateDisabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironmentUpdateDisabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaEnvironmentUpdateDisabledErrorDetails) ImplErrorDetailsUnion() {}

// `expected_last_heartbeat` does not match the stored heartbeat, so another worker
// may hold the lease.
type BetaHeartbeatPreconditionFailedErrorDetails struct {
	// The lease as the server holds it.
	CurrentState BetaWorkLeaseState                                `json:"current_state" api:"required"`
	ErrorCode    constant.BetaHeartbeatPreconditionFailedErrorCode `json:"error_code" default:"heartbeat_precondition_failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CurrentState respjson.Field
		ErrorCode    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaHeartbeatPreconditionFailedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaHeartbeatPreconditionFailedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaHeartbeatPreconditionFailedErrorDetails) ImplErrorDetailsUnion() {}

// A work item's lease as the server holds it.
type BetaWorkLeaseState struct {
	// RFC 3339 timestamp of the stored heartbeat, or null if there has been none
	LastHeartbeat string `json:"last_heartbeat" api:"required"`
	// Always false: a refused heartbeat does not extend the lease
	LeaseExtended bool `json:"lease_extended" api:"required"`
	// Current state of the work item
	//
	// Any of "queued", "starting", "active", "stopping", "stopped".
	State BetaWorkLeaseStateState `json:"state" api:"required"`
	// TTL of the lease
	TTLSeconds int64 `json:"ttl_seconds" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		LastHeartbeat respjson.Field
		LeaseExtended respjson.Field
		State         respjson.Field
		TTLSeconds    respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkLeaseState) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkLeaseState) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current state of the work item
type BetaWorkLeaseStateState string

const (
	BetaWorkLeaseStateStateQueued   BetaWorkLeaseStateState = "queued"
	BetaWorkLeaseStateStateStarting BetaWorkLeaseStateState = "starting"
	BetaWorkLeaseStateStateActive   BetaWorkLeaseStateState = "active"
	BetaWorkLeaseStateStateStopping BetaWorkLeaseStateState = "stopping"
	BetaWorkLeaseStateStateStopped  BetaWorkLeaseStateState = "stopped"
)

// `metadata` has too many keys, a key or value that is too long, or a reserved
// key. The message names each problem.
type BetaInvalidMetadataErrorDetails struct {
	ErrorCode constant.BetaInvalidMetadataErrorCode `json:"error_code" default:"invalid_metadata"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInvalidMetadataErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaInvalidMetadataErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaInvalidMetadataErrorDetails) ImplErrorDetailsUnion() {}

// Anthropic has temporarily disabled billing changes for the user account the
// credential belongs to. Contact support.
type BetaBillingFrozenErrorDetails struct {
	ErrorCode constant.BetaBillingFrozenErrorCode `json:"error_code" default:"billing_frozen"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaBillingFrozenErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaBillingFrozenErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaBillingFrozenErrorDetails) ImplErrorDetailsUnion() {}

// The organization cannot retrieve, set or delete spend limits with an Admin API
// credential.
type BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails struct {
	ErrorCode constant.BetaSpendLimitsAdminAPIWritesNotEnabledErrorCode `json:"error_code" default:"spend_limits_admin_api_writes_not_enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails) ImplErrorDetailsUnion() {}

// The request's `Content-Type` header names a charset other than UTF-8, ISO-8859-1
// or US-ASCII. Send the request with no charset, or with one of these.
type BetaUnsupportedCharsetErrorDetails struct {
	ErrorCode constant.BetaUnsupportedCharsetErrorCode `json:"error_code" default:"unsupported_charset"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUnsupportedCharsetErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaUnsupportedCharsetErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaUnsupportedCharsetErrorDetails) ImplErrorDetailsUnion() {}

// The marketplace already holds this name. `plugin_id` is the Plugin that holds
// it, and can be absent.
type BetaPluginNameTakenErrorDetails struct {
	ErrorCode constant.BetaPluginNameTakenErrorCode `json:"error_code" default:"plugin_name_taken"`
	PluginID  string                                `json:"plugin_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		PluginID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginNameTakenErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginNameTakenErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaPluginNameTakenErrorDetails) ImplErrorDetailsUnion() {}

// One of the Plugin's skills has the name of an organization skill. `skill_name`
// is that name.
type BetaSkillNameTakenErrorDetails struct {
	ErrorCode constant.BetaSkillNameTakenErrorCode `json:"error_code" default:"skill_name_taken"`
	SkillName string                               `json:"skill_name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		SkillName   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSkillNameTakenErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSkillNameTakenErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSkillNameTakenErrorDetails) ImplErrorDetailsUnion() {}

// The Plugin version was stored but is not yet usable in claude.ai. `plugin_id`
// and `plugin_version_id` identify what was stored.
type BetaRegistrationPendingErrorDetails struct {
	ErrorCode       constant.BetaRegistrationPendingErrorCode `json:"error_code" default:"registration_pending"`
	PluginID        string                                    `json:"plugin_id" api:"required"`
	PluginVersionID string                                    `json:"plugin_version_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode       respjson.Field
		PluginID        respjson.Field
		PluginVersionID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRegistrationPendingErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRegistrationPendingErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRegistrationPendingErrorDetails) ImplErrorDetailsUnion() {}

// The version in `plugin_version_id` did not pass its content scan, so it cannot
// be served.
type BetaScanFailedErrorDetails struct {
	ErrorCode       constant.BetaScanFailedErrorCode `json:"error_code" default:"scan_failed"`
	PluginVersionID string                           `json:"plugin_version_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode       respjson.Field
		PluginVersionID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaScanFailedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaScanFailedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaScanFailedErrorDetails) ImplErrorDetailsUnion() {}

// The content scan of the version in `plugin_version_id` is still running. Retry
// once the scan finishes.
type BetaScanPendingErrorDetails struct {
	ErrorCode       constant.BetaScanPendingErrorCode `json:"error_code" default:"scan_pending"`
	PluginVersionID string                            `json:"plugin_version_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode       respjson.Field
		PluginVersionID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaScanPendingErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaScanPendingErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaScanPendingErrorDetails) ImplErrorDetailsUnion() {}

// The request set the channel up, and `settings` was then refused. The channel
// stays set up with every setting unset.
type BetaSlackChannelSetUpSettingsRefusedErrorDetails struct {
	// A stable code that names the cause.
	ErrorCode constant.BetaSlackChannelSetUpSettingsRefusedErrorCode `json:"error_code" default:"slack_channel_set_up_settings_refused"`
	// The identifier of the channel the request set up, starting with `slch_`.
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode      respjson.Field
		SlackChannelID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSlackChannelSetUpSettingsRefusedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSlackChannelSetUpSettingsRefusedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSlackChannelSetUpSettingsRefusedErrorDetails) ImplErrorDetailsUnion() {}

// The API key lacks the scope the request needs: `read:claude_tag` for a read,
// `write:claude_tag` for a change.
type BetaAPIKeyScopeMissingErrorDetails struct {
	ErrorCode constant.BetaAPIKeyScopeMissingErrorCode `json:"error_code" default:"api_key_scope_missing"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAPIKeyScopeMissingErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIKeyScopeMissingErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaAPIKeyScopeMissingErrorDetails) ImplErrorDetailsUnion() {}

// Changes through the Claude Tag API are not enabled for the organization.
type BetaClaudeTagWritesNotEnabledErrorDetails struct {
	ErrorCode constant.BetaClaudeTagWritesNotEnabledErrorCode `json:"error_code" default:"claude_tag_writes_not_enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagWritesNotEnabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagWritesNotEnabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaClaudeTagWritesNotEnabledErrorDetails) ImplErrorDetailsUnion() {}

// The organization's compliance configuration does not let Claude Tag store
// anything new.
type BetaComplianceConfigurationRestrictedErrorDetails struct {
	ErrorCode constant.BetaComplianceConfigurationRestrictedErrorCode `json:"error_code" default:"compliance_configuration_restricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceConfigurationRestrictedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceConfigurationRestrictedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaComplianceConfigurationRestrictedErrorDetails) ImplErrorDetailsUnion() {}

// The credential is not a scoped API key that is bound at the organization level.
type BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails struct {
	ErrorCode constant.BetaOrganizationLevelScopedAPIKeyRequiredErrorCode `json:"error_code" default:"organization_level_scoped_api_key_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails) ImplErrorDetailsUnion() {}

// `settings.environment_id` names no Environment within the organization's reach.
type BetaEnvironmentNotFoundErrorDetails struct {
	ErrorCode constant.BetaEnvironmentNotFoundErrorCode `json:"error_code" default:"environment_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironmentNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironmentNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaEnvironmentNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// `settings.environment_id` names no Environment within the organization's reach,
// and the request had set the channel up by then. The channel stays set up with
// every setting unset.
type BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails struct {
	ErrorCode constant.BetaSlackChannelSetUpEnvironmentNotFoundErrorCode `json:"error_code" default:"slack_channel_set_up_environment_not_found"`
	// The identifier of the channel the request set up, starting with `slch_`.
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode      respjson.Field
		SlackChannelID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// The request set the channel up, and something else then changed the channel or
// removed it. Nothing of `settings` was stored.
type BetaSlackChannelSetUpSettingsConflictErrorDetails struct {
	// A stable code that names the cause.
	ErrorCode constant.BetaSlackChannelSetUpSettingsConflictErrorCode `json:"error_code" default:"slack_channel_set_up_settings_conflict"`
	// The identifier of the channel the request set up, starting with `slch_`.
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode      respjson.Field
		SlackChannelID respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSlackChannelSetUpSettingsConflictErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSlackChannelSetUpSettingsConflictErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSlackChannelSetUpSettingsConflictErrorDetails) ImplErrorDetailsUnion() {}

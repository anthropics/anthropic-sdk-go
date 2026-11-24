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

type AuthenticationErrorObject struct {
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
func (r AuthenticationErrorObject) RawJSON() string { return r.JSON.raw }
func (r *AuthenticationErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (AuthenticationErrorObject) ImplErrorObjectUnion() {}

type BillingErrorObject struct {
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
func (r BillingErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BillingErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BillingErrorObject) ImplErrorObjectUnion() {}

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

// The organization has reached the spend limit of its contract.
type ContractSpendLimitReachedErrorDetails struct {
	ErrorCode constant.ContractSpendLimitReachedErrorCode `json:"error_code" default:"contract_spend_limit_reached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContractSpendLimitReachedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ContractSpendLimitReachedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ContractSpendLimitReachedErrorDetails) ImplErrorDetailsUnion() {}

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
// [CompactionToolChangesMismatchErrorDetails],
// [ImageDimensionsTooLargeErrorDetails], [OrganizationOnHoldErrorDetails],
// [RequestBodyEncodingInvalidErrorDetails],
// [RequestBodyEncodingUnsupportedErrorDetails],
// [ThreadAlreadyContinuedErrorDetails], [ThreadContinueInFlightErrorDetails],
// [ThreadFingerprintMismatchErrorDetails], [ThreadUnsupportedRequestErrorDetails],
// [ToolNameConflictErrorDetails], [ToolReferenceUnresolvedErrorDetails],
// [UserProfileRequiredErrorDetails],
// [LinkedOrganizationBillingUnavailableErrorDetails],
// [LinkedOrganizationSpendLimitReachedErrorDetails],
// [UserProfileBillingNotAuthorizedErrorDetails],
// [UserProfileNotBillableErrorDetails], [ThreadNotFoundErrorDetails],
// [SkillReadBlockedUndecryptableErrorDetails],
// [RequestBodyTooLargeAfterDecodingErrorDetails],
// [ContractSpendLimitReachedErrorDetails],
// [EnforcedSpendLimitReachedErrorDetails], [RampLimitExceededErrorDetails],
// [CMEKKeyUnreachableErrorDetails], [CompactionUnavailableErrorDetails],
// [WorkspaceHumanReviewDisabledErrorDetails], [StorageQuotaExceededErrorDetails],
// [FileNotDownloadableErrorDetails], [FileExpiredErrorDetails],
// [CMEKContextUnavailableErrorDetails], [OrgAPIKeyLimitExceededErrorDetails],
// [WorkspaceAPIKeyLimitExceededErrorDetails],
// [SandboxEgressRestrictedErrorDetails], [UserProfileNotFoundErrorDetails],
// [ExternalKeyInWorkspacePolicyErrorDetails], [AlreadyMemberErrorDetails],
// [CannotInviteEmailsWithUnverifiedDomainsErrorDetails],
// [EmailDomainNotAllowedErrorDetails], [InviteAlreadyExistsErrorDetails],
// [OrgMemberLimitExceededErrorDetails], [RBACGroupGrantsOwnerAccessErrorDetails],
// [RBACGroupOwnerAccessUnconfirmedErrorDetails],
// [CannotRemoveLastMemberErrorDetails], [ConcurrentMembershipChangeErrorDetails],
// [ExternalKeyValidationFailedErrorDetails],
// [ExternalKeyValidationTimedOutErrorDetails], [KeyProfileMismatchErrorDetails],
// [WorkspacePolicyViolationErrorDetails], [ExternalKeyChangedErrorDetails],
// [WorkspaceCreatedIncompleteErrorDetails], [StorageProfileMismatchErrorDetails],
// [ChangeOutcomeUnknownErrorDetails], [StorageAttachedDuringArchiveErrorDetails],
// [WorkspaceMemberLimitExceededErrorDetails],
// [OrganizationTypeNotSupportedErrorDetails],
// [ParentOrganizationRequiredErrorDetails], [RBACRoleNotFoundErrorDetails],
// [BetaEnvironmentCreateDisabledErrorDetails],
// [BetaEnvironmentArchivedErrorDetails],
// [BetaInvalidConfigTypeChangeErrorDetails],
// [BetaEnvironmentUpdateDisabledErrorDetails],
// [BetaHeartbeatPreconditionFailedErrorDetails],
// [BetaInvalidMetadataErrorDetails], [BetaGrantRequiredErrorDetails],
// [BetaOrganizationNotEligibleErrorDetails], [BetaProviderNotAllowedErrorDetails],
// [BetaUnsupportedRegionErrorDetails], [BetaStorageAlreadyRegisteredErrorDetails],
// [BetaAuthorizationCheckUnavailableErrorDetails],
// [BetaInvalidStorageARNErrorDetails], [BetaExternalStorageInUseErrorDetails],
// [BetaExternalStorageInWorkspacePolicyErrorDetails],
// [BetaStorageVerificationFailedErrorDetails], [BetaCMEKNotApprovedErrorDetails],
// [BetaExternalKeyNotFoundErrorDetails],
// [BetaExternalKeyProviderMismatchErrorDetails],
// [BetaGCPProjectNotVerifiedErrorDetails], [BetaRegionNotSupportedErrorDetails],
// [BetaStorageNotApprovedErrorDetails],
// [BetaGCPProjectControlledByAnotherOrganizationErrorDetails],
// [BetaStorageInUseErrorDetails], [BetaCannotSignOutPrimaryOwnerErrorDetails],
// [BetaCannotSignOutSelfErrorDetails],
// [BetaOrganizationNotEnrolledInBetaErrorDetails],
// [BetaCannotSignOutManagedMemberErrorDetails],
// [BetaCannotSignOutMembershipAdminErrorDetails],
// [BetaCannotSignOutOwnerErrorDetails],
// [BetaCannotSignOutPrivilegedMemberErrorDetails],
// [BetaApprovalRequiredErrorDetails],
// [BetaExternalStorageNotAvailableErrorDetails],
// [BetaHumanReviewNotAvailableErrorDetails],
// [BetaHumanReviewOnWithExternalStorageErrorDetails],
// [BetaStorageNotAttachableErrorDetails],
// [BetaTranscriptsStorageRequiredErrorDetails],
// [BetaStorageAttachLimitReachedErrorDetails],
// [BetaStorageProfileStillRequiredErrorDetails],
// [BetaStorageChangedConcurrentlyErrorDetails],
// [BetaWorkspaceChangedConcurrentlyErrorDetails],
// [BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails],
// [BetaUnsupportedCharsetErrorDetails], [BetaPluginNameTakenErrorDetails],
// [BetaSkillNameTakenErrorDetails], [BetaRegistrationPendingErrorDetails],
// [BetaScanFailedErrorDetails], [BetaScanPendingErrorDetails],
// [BetaRequestBodyTimeoutErrorDetails],
// [BetaSettingsDocumentLimitExceededErrorDetails],
// [BetaSettingsHistoryLimitExceededErrorDetails],
// [BetaCMEKFeatureUnsupportedErrorDetails],
// [BetaRBACGroupManagedSettingsNotFoundErrorDetails],
// [BetaRBACGroupOrderLimitExceededErrorDetails],
// [BetaRBACGroupNotFoundErrorDetails],
// [BetaSlackChannelSetUpSettingsRefusedErrorDetails],
// [BetaAPIKeyScopeMissingErrorDetails],
// [BetaClaudeTagWritesNotEnabledErrorDetails],
// [BetaComplianceConfigurationRestrictedErrorDetails],
// [BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails],
// [BetaEnvironmentNotFoundErrorDetails],
// [BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails],
// [BetaSlackChannelSetUpSettingsConflictErrorDetails],
// [BetaBundleAttachmentNotSupportedErrorDetails], [BetaBundleOnCallErrorDetails].
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
	// "image_dimensions_too_large", "organization_on_hold",
	// "request_body_encoding_invalid", "request_body_encoding_unsupported",
	// "thread_already_continued", "thread_continue_in_flight",
	// "thread_fingerprint_mismatch", "thread_unsupported_request",
	// "tool_name_conflict", "tool_reference_unresolved", "user_profile_required",
	// "linked_organization_billing_unavailable",
	// "linked_organization_spend_limit_reached",
	// "user_profile_billing_not_authorized", "user_profile_not_billable",
	// "thread_not_found", "skill_read_blocked_undecryptable",
	// "request_body_too_large_after_decoding", "contract_spend_limit_reached",
	// "enforced_spend_limit_reached", "ramp_limit_exceeded", "cmek_key_unreachable",
	// "compaction_unavailable", "workspace_human_review_disabled",
	// "storage_quota_exceeded", "file_not_downloadable", "file_expired",
	// "cmek_context_unavailable", "org_api_key_limit_exceeded",
	// "workspace_api_key_limit_exceeded", "sandbox_egress_restricted",
	// "user_profile_not_found", "external_key_in_workspace_policy", "already_member",
	// "cannot_invite_emails_with_unverified_domains", "email_domain_not_allowed",
	// "invite_already_exists", "org_member_limit_exceeded",
	// "rbac_group_grants_owner_access", "rbac_group_owner_access_unconfirmed",
	// "cannot_remove_last_member", "concurrent_membership_change",
	// "external_key_validation_failed", "external_key_validation_timed_out",
	// "key_profile_mismatch", "workspace_policy_violation", "external_key_changed",
	// "workspace_created_incomplete", "storage_profile_mismatch",
	// "change_outcome_unknown", "storage_attached_during_archive",
	// "workspace_member_limit_exceeded", "organization_type_not_supported",
	// "parent_organization_required", "rbac_role_not_found",
	// "environment_create_disabled", "environment_archived",
	// "invalid_config_type_change", "environment_update_disabled",
	// "heartbeat_precondition_failed", "invalid_metadata", "grant_required",
	// "organization_not_eligible", "provider_not_allowed", "unsupported_region",
	// "storage_already_registered", "authorization_check_unavailable",
	// "invalid_storage_arn", "external_storage_in_use",
	// "external_storage_in_workspace_policy", "storage_verification_failed",
	// "cmek_not_approved", "external_key_not_found", "external_key_provider_mismatch",
	// "gcp_project_not_verified", "region_not_supported", "storage_not_approved",
	// "gcp_project_controlled_by_another_organization", "storage_in_use",
	// "cannot_sign_out_primary_owner", "cannot_sign_out_self",
	// "organization_not_enrolled_in_beta", "cannot_sign_out_managed_member",
	// "cannot_sign_out_membership_admin", "cannot_sign_out_owner",
	// "cannot_sign_out_privileged_member", "approval_required",
	// "external_storage_not_available", "human_review_not_available",
	// "human_review_on_with_external_storage", "storage_not_attachable",
	// "transcripts_storage_required", "storage_attach_limit_reached",
	// "storage_profile_still_required", "storage_changed_concurrently",
	// "workspace_changed_concurrently", "spend_limits_admin_api_writes_not_enabled",
	// "unsupported_charset", "plugin_name_taken", "skill_name_taken",
	// "registration_pending", "scan_failed", "scan_pending", "request_body_timeout",
	// "settings_document_limit_exceeded", "settings_history_limit_exceeded",
	// "cmek_feature_unsupported", "rbac_group_managed_settings_not_found",
	// "rbac_group_order_limit_exceeded", "rbac_group_not_found",
	// "slack_channel_set_up_settings_refused", "api_key_scope_missing",
	// "claude_tag_writes_not_enabled", "compliance_configuration_restricted",
	// "organization_level_scoped_api_key_required", "environment_not_found",
	// "slack_channel_set_up_environment_not_found",
	// "slack_channel_set_up_settings_conflict", "bundle_attachment_not_supported",
	// "bundle_on_call".
	ErrorCode string `json:"error_code"`
	// This field is from variant [ImageDimensionsTooLargeErrorDetails].
	MaxSideLengthPixels int64 `json:"max_side_length_pixels"`
	// This field is from variant [ImageDimensionsTooLargeErrorDetails].
	Paths []string `json:"paths"`
	// This field is from variant [OrganizationOnHoldErrorDetails].
	AppealURL string `json:"appeal_url"`
	// This field is from variant [ThreadUnsupportedRequestErrorDetails].
	Unsupported []ThreadUnsupportedFeatureName `json:"unsupported"`
	Field       string                         `json:"field"`
	// This field is from variant [WorkspaceCreatedIncompleteErrorDetails].
	WorkspaceID string `json:"workspace_id"`
	// This field is from variant [BetaHeartbeatPreconditionFailedErrorDetails].
	CurrentState BetaWorkLeaseState `json:"current_state"`
	// This field is from variant [BetaStorageVerificationFailedErrorDetails].
	FailedChecks []BetaExternalStorageFailedCheck `json:"failed_checks"`
	PluginID     string                           `json:"plugin_id"`
	// This field is from variant [BetaSkillNameTakenErrorDetails].
	SkillName       string `json:"skill_name"`
	PluginVersionID string `json:"plugin_version_id"`
	SlackChannelID  string `json:"slack_channel_id"`
	JSON            struct {
		ErrorCode           respjson.Field
		MaxSideLengthPixels respjson.Field
		Paths               respjson.Field
		AppealURL           respjson.Field
		Unsupported         respjson.Field
		Field               respjson.Field
		WorkspaceID         respjson.Field
		CurrentState        respjson.Field
		FailedChecks        respjson.Field
		PluginID            respjson.Field
		SkillName           respjson.Field
		PluginVersionID     respjson.Field
		SlackChannelID      respjson.Field
		raw                 string
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
//	case shared.ImageDimensionsTooLargeErrorDetails:
//	case shared.OrganizationOnHoldErrorDetails:
//	case shared.RequestBodyEncodingInvalidErrorDetails:
//	case shared.RequestBodyEncodingUnsupportedErrorDetails:
//	case shared.ThreadAlreadyContinuedErrorDetails:
//	case shared.ThreadContinueInFlightErrorDetails:
//	case shared.ThreadFingerprintMismatchErrorDetails:
//	case shared.ThreadUnsupportedRequestErrorDetails:
//	case shared.ToolNameConflictErrorDetails:
//	case shared.ToolReferenceUnresolvedErrorDetails:
//	case shared.UserProfileRequiredErrorDetails:
//	case shared.LinkedOrganizationBillingUnavailableErrorDetails:
//	case shared.LinkedOrganizationSpendLimitReachedErrorDetails:
//	case shared.UserProfileBillingNotAuthorizedErrorDetails:
//	case shared.UserProfileNotBillableErrorDetails:
//	case shared.ThreadNotFoundErrorDetails:
//	case shared.SkillReadBlockedUndecryptableErrorDetails:
//	case shared.RequestBodyTooLargeAfterDecodingErrorDetails:
//	case shared.ContractSpendLimitReachedErrorDetails:
//	case shared.EnforcedSpendLimitReachedErrorDetails:
//	case shared.RampLimitExceededErrorDetails:
//	case shared.CMEKKeyUnreachableErrorDetails:
//	case shared.CompactionUnavailableErrorDetails:
//	case shared.WorkspaceHumanReviewDisabledErrorDetails:
//	case shared.StorageQuotaExceededErrorDetails:
//	case shared.FileNotDownloadableErrorDetails:
//	case shared.FileExpiredErrorDetails:
//	case shared.CMEKContextUnavailableErrorDetails:
//	case shared.OrgAPIKeyLimitExceededErrorDetails:
//	case shared.WorkspaceAPIKeyLimitExceededErrorDetails:
//	case shared.SandboxEgressRestrictedErrorDetails:
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
//	case shared.WorkspaceCreatedIncompleteErrorDetails:
//	case shared.StorageProfileMismatchErrorDetails:
//	case shared.ChangeOutcomeUnknownErrorDetails:
//	case shared.StorageAttachedDuringArchiveErrorDetails:
//	case shared.WorkspaceMemberLimitExceededErrorDetails:
//	case shared.OrganizationTypeNotSupportedErrorDetails:
//	case shared.ParentOrganizationRequiredErrorDetails:
//	case shared.RBACRoleNotFoundErrorDetails:
//	case shared.BetaEnvironmentCreateDisabledErrorDetails:
//	case shared.BetaEnvironmentArchivedErrorDetails:
//	case shared.BetaInvalidConfigTypeChangeErrorDetails:
//	case shared.BetaEnvironmentUpdateDisabledErrorDetails:
//	case shared.BetaHeartbeatPreconditionFailedErrorDetails:
//	case shared.BetaInvalidMetadataErrorDetails:
//	case shared.BetaGrantRequiredErrorDetails:
//	case shared.BetaOrganizationNotEligibleErrorDetails:
//	case shared.BetaProviderNotAllowedErrorDetails:
//	case shared.BetaUnsupportedRegionErrorDetails:
//	case shared.BetaStorageAlreadyRegisteredErrorDetails:
//	case shared.BetaAuthorizationCheckUnavailableErrorDetails:
//	case shared.BetaInvalidStorageARNErrorDetails:
//	case shared.BetaExternalStorageInUseErrorDetails:
//	case shared.BetaExternalStorageInWorkspacePolicyErrorDetails:
//	case shared.BetaStorageVerificationFailedErrorDetails:
//	case shared.BetaCMEKNotApprovedErrorDetails:
//	case shared.BetaExternalKeyNotFoundErrorDetails:
//	case shared.BetaExternalKeyProviderMismatchErrorDetails:
//	case shared.BetaGCPProjectNotVerifiedErrorDetails:
//	case shared.BetaRegionNotSupportedErrorDetails:
//	case shared.BetaStorageNotApprovedErrorDetails:
//	case shared.BetaGCPProjectControlledByAnotherOrganizationErrorDetails:
//	case shared.BetaStorageInUseErrorDetails:
//	case shared.BetaCannotSignOutPrimaryOwnerErrorDetails:
//	case shared.BetaCannotSignOutSelfErrorDetails:
//	case shared.BetaOrganizationNotEnrolledInBetaErrorDetails:
//	case shared.BetaCannotSignOutManagedMemberErrorDetails:
//	case shared.BetaCannotSignOutMembershipAdminErrorDetails:
//	case shared.BetaCannotSignOutOwnerErrorDetails:
//	case shared.BetaCannotSignOutPrivilegedMemberErrorDetails:
//	case shared.BetaApprovalRequiredErrorDetails:
//	case shared.BetaExternalStorageNotAvailableErrorDetails:
//	case shared.BetaHumanReviewNotAvailableErrorDetails:
//	case shared.BetaHumanReviewOnWithExternalStorageErrorDetails:
//	case shared.BetaStorageNotAttachableErrorDetails:
//	case shared.BetaTranscriptsStorageRequiredErrorDetails:
//	case shared.BetaStorageAttachLimitReachedErrorDetails:
//	case shared.BetaStorageProfileStillRequiredErrorDetails:
//	case shared.BetaStorageChangedConcurrentlyErrorDetails:
//	case shared.BetaWorkspaceChangedConcurrentlyErrorDetails:
//	case shared.BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails:
//	case shared.BetaUnsupportedCharsetErrorDetails:
//	case shared.BetaPluginNameTakenErrorDetails:
//	case shared.BetaSkillNameTakenErrorDetails:
//	case shared.BetaRegistrationPendingErrorDetails:
//	case shared.BetaScanFailedErrorDetails:
//	case shared.BetaScanPendingErrorDetails:
//	case shared.BetaRequestBodyTimeoutErrorDetails:
//	case shared.BetaSettingsDocumentLimitExceededErrorDetails:
//	case shared.BetaSettingsHistoryLimitExceededErrorDetails:
//	case shared.BetaCMEKFeatureUnsupportedErrorDetails:
//	case shared.BetaRBACGroupManagedSettingsNotFoundErrorDetails:
//	case shared.BetaRBACGroupOrderLimitExceededErrorDetails:
//	case shared.BetaRBACGroupNotFoundErrorDetails:
//	case shared.BetaSlackChannelSetUpSettingsRefusedErrorDetails:
//	case shared.BetaAPIKeyScopeMissingErrorDetails:
//	case shared.BetaClaudeTagWritesNotEnabledErrorDetails:
//	case shared.BetaComplianceConfigurationRestrictedErrorDetails:
//	case shared.BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails:
//	case shared.BetaEnvironmentNotFoundErrorDetails:
//	case shared.BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails:
//	case shared.BetaSlackChannelSetUpSettingsConflictErrorDetails:
//	case shared.BetaBundleAttachmentNotSupportedErrorDetails:
//	case shared.BetaBundleOnCallErrorDetails:
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
	case "image_dimensions_too_large":
		return u.AsImageDimensionsTooLarge()
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
	case "user_profile_required":
		return u.AsUserProfileRequired()
	case "linked_organization_billing_unavailable":
		return u.AsLinkedOrganizationBillingUnavailable()
	case "linked_organization_spend_limit_reached":
		return u.AsLinkedOrganizationSpendLimitReached()
	case "user_profile_billing_not_authorized":
		return u.AsUserProfileBillingNotAuthorized()
	case "user_profile_not_billable":
		return u.AsUserProfileNotBillable()
	case "thread_not_found":
		return u.AsThreadNotFound()
	case "skill_read_blocked_undecryptable":
		return u.AsSkillReadBlockedUndecryptable()
	case "request_body_too_large_after_decoding":
		return u.AsRequestBodyTooLargeAfterDecoding()
	case "contract_spend_limit_reached":
		return u.AsContractSpendLimitReached()
	case "enforced_spend_limit_reached":
		return u.AsEnforcedSpendLimitReached()
	case "ramp_limit_exceeded":
		return u.AsRampLimitExceeded()
	case "cmek_key_unreachable":
		return u.AsCMEKKeyUnreachable()
	case "compaction_unavailable":
		return u.AsCompactionUnavailable()
	case "workspace_human_review_disabled":
		return u.AsWorkspaceHumanReviewDisabled()
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
	case "sandbox_egress_restricted":
		return u.AsSandboxEgressRestricted()
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
	case "workspace_created_incomplete":
		return u.AsWorkspaceCreatedIncomplete()
	case "storage_profile_mismatch":
		return u.AsStorageProfileMismatch()
	case "change_outcome_unknown":
		return u.AsChangeOutcomeUnknown()
	case "storage_attached_during_archive":
		return u.AsStorageAttachedDuringArchive()
	case "workspace_member_limit_exceeded":
		return u.AsWorkspaceMemberLimitExceeded()
	case "organization_type_not_supported":
		return u.AsOrganizationTypeNotSupported()
	case "parent_organization_required":
		return u.AsParentOrganizationRequired()
	case "rbac_role_not_found":
		return u.AsRBACRoleNotFound()
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
	case "grant_required":
		return u.AsGrantRequired()
	case "organization_not_eligible":
		return u.AsOrganizationNotEligible()
	case "provider_not_allowed":
		return u.AsProviderNotAllowed()
	case "unsupported_region":
		return u.AsUnsupportedRegion()
	case "storage_already_registered":
		return u.AsStorageAlreadyRegistered()
	case "authorization_check_unavailable":
		return u.AsAuthorizationCheckUnavailable()
	case "invalid_storage_arn":
		return u.AsInvalidStorageARN()
	case "external_storage_in_use":
		return u.AsExternalStorageInUse()
	case "external_storage_in_workspace_policy":
		return u.AsExternalStorageInWorkspacePolicy()
	case "storage_verification_failed":
		return u.AsStorageVerificationFailed()
	case "cmek_not_approved":
		return u.AsCMEKNotApproved()
	case "external_key_not_found":
		return u.AsExternalKeyNotFound()
	case "external_key_provider_mismatch":
		return u.AsExternalKeyProviderMismatch()
	case "gcp_project_not_verified":
		return u.AsGCPProjectNotVerified()
	case "region_not_supported":
		return u.AsRegionNotSupported()
	case "storage_not_approved":
		return u.AsStorageNotApproved()
	case "gcp_project_controlled_by_another_organization":
		return u.AsGCPProjectControlledByAnotherOrganization()
	case "storage_in_use":
		return u.AsStorageInUse()
	case "cannot_sign_out_primary_owner":
		return u.AsCannotSignOutPrimaryOwner()
	case "cannot_sign_out_self":
		return u.AsCannotSignOutSelf()
	case "organization_not_enrolled_in_beta":
		return u.AsOrganizationNotEnrolledInBeta()
	case "cannot_sign_out_managed_member":
		return u.AsCannotSignOutManagedMember()
	case "cannot_sign_out_membership_admin":
		return u.AsCannotSignOutMembershipAdmin()
	case "cannot_sign_out_owner":
		return u.AsCannotSignOutOwner()
	case "cannot_sign_out_privileged_member":
		return u.AsCannotSignOutPrivilegedMember()
	case "approval_required":
		return u.AsApprovalRequired()
	case "external_storage_not_available":
		return u.AsExternalStorageNotAvailable()
	case "human_review_not_available":
		return u.AsHumanReviewNotAvailable()
	case "human_review_on_with_external_storage":
		return u.AsHumanReviewOnWithExternalStorage()
	case "storage_not_attachable":
		return u.AsStorageNotAttachable()
	case "transcripts_storage_required":
		return u.AsTranscriptsStorageRequired()
	case "storage_attach_limit_reached":
		return u.AsStorageAttachLimitReached()
	case "storage_profile_still_required":
		return u.AsStorageProfileStillRequired()
	case "storage_changed_concurrently":
		return u.AsStorageChangedConcurrently()
	case "workspace_changed_concurrently":
		return u.AsWorkspaceChangedConcurrently()
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
	case "request_body_timeout":
		return u.AsRequestBodyTimeout()
	case "settings_document_limit_exceeded":
		return u.AsSettingsDocumentLimitExceeded()
	case "settings_history_limit_exceeded":
		return u.AsSettingsHistoryLimitExceeded()
	case "cmek_feature_unsupported":
		return u.AsCMEKFeatureUnsupported()
	case "rbac_group_managed_settings_not_found":
		return u.AsRBACGroupManagedSettingsNotFound()
	case "rbac_group_order_limit_exceeded":
		return u.AsRBACGroupOrderLimitExceeded()
	case "rbac_group_not_found":
		return u.AsRBACGroupNotFound()
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
	case "bundle_attachment_not_supported":
		return u.AsBundleAttachmentNotSupported()
	case "bundle_on_call":
		return u.AsBundleOnCall()
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

func (u ErrorDetailsUnion) AsImageDimensionsTooLarge() (v ImageDimensionsTooLargeErrorDetails) {
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

func (u ErrorDetailsUnion) AsUserProfileRequired() (v UserProfileRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsLinkedOrganizationBillingUnavailable() (v LinkedOrganizationBillingUnavailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsLinkedOrganizationSpendLimitReached() (v LinkedOrganizationSpendLimitReachedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsUserProfileBillingNotAuthorized() (v UserProfileBillingNotAuthorizedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsUserProfileNotBillable() (v UserProfileNotBillableErrorDetails) {
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

func (u ErrorDetailsUnion) AsContractSpendLimitReached() (v ContractSpendLimitReachedErrorDetails) {
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

func (u ErrorDetailsUnion) AsWorkspaceHumanReviewDisabled() (v WorkspaceHumanReviewDisabledErrorDetails) {
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

func (u ErrorDetailsUnion) AsSandboxEgressRestricted() (v SandboxEgressRestrictedErrorDetails) {
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

func (u ErrorDetailsUnion) AsWorkspaceCreatedIncomplete() (v WorkspaceCreatedIncompleteErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageProfileMismatch() (v StorageProfileMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsChangeOutcomeUnknown() (v ChangeOutcomeUnknownErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageAttachedDuringArchive() (v StorageAttachedDuringArchiveErrorDetails) {
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

func (u ErrorDetailsUnion) AsParentOrganizationRequired() (v ParentOrganizationRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACRoleNotFound() (v RBACRoleNotFoundErrorDetails) {
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

func (u ErrorDetailsUnion) AsGrantRequired() (v BetaGrantRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrganizationNotEligible() (v BetaOrganizationNotEligibleErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsProviderNotAllowed() (v BetaProviderNotAllowedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsUnsupportedRegion() (v BetaUnsupportedRegionErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageAlreadyRegistered() (v BetaStorageAlreadyRegisteredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsAuthorizationCheckUnavailable() (v BetaAuthorizationCheckUnavailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsInvalidStorageARN() (v BetaInvalidStorageARNErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalStorageInUse() (v BetaExternalStorageInUseErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalStorageInWorkspacePolicy() (v BetaExternalStorageInWorkspacePolicyErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageVerificationFailed() (v BetaStorageVerificationFailedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKNotApproved() (v BetaCMEKNotApprovedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyNotFound() (v BetaExternalKeyNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalKeyProviderMismatch() (v BetaExternalKeyProviderMismatchErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsGCPProjectNotVerified() (v BetaGCPProjectNotVerifiedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRegionNotSupported() (v BetaRegionNotSupportedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageNotApproved() (v BetaStorageNotApprovedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsGCPProjectControlledByAnotherOrganization() (v BetaGCPProjectControlledByAnotherOrganizationErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageInUse() (v BetaStorageInUseErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutPrimaryOwner() (v BetaCannotSignOutPrimaryOwnerErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutSelf() (v BetaCannotSignOutSelfErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsOrganizationNotEnrolledInBeta() (v BetaOrganizationNotEnrolledInBetaErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutManagedMember() (v BetaCannotSignOutManagedMemberErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutMembershipAdmin() (v BetaCannotSignOutMembershipAdminErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutOwner() (v BetaCannotSignOutOwnerErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCannotSignOutPrivilegedMember() (v BetaCannotSignOutPrivilegedMemberErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsApprovalRequired() (v BetaApprovalRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsExternalStorageNotAvailable() (v BetaExternalStorageNotAvailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsHumanReviewNotAvailable() (v BetaHumanReviewNotAvailableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsHumanReviewOnWithExternalStorage() (v BetaHumanReviewOnWithExternalStorageErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageNotAttachable() (v BetaStorageNotAttachableErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsTranscriptsStorageRequired() (v BetaTranscriptsStorageRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageAttachLimitReached() (v BetaStorageAttachLimitReachedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageProfileStillRequired() (v BetaStorageProfileStillRequiredErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsStorageChangedConcurrently() (v BetaStorageChangedConcurrentlyErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsWorkspaceChangedConcurrently() (v BetaWorkspaceChangedConcurrentlyErrorDetails) {
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

func (u ErrorDetailsUnion) AsRequestBodyTimeout() (v BetaRequestBodyTimeoutErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSettingsDocumentLimitExceeded() (v BetaSettingsDocumentLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsSettingsHistoryLimitExceeded() (v BetaSettingsHistoryLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsCMEKFeatureUnsupported() (v BetaCMEKFeatureUnsupportedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACGroupManagedSettingsNotFound() (v BetaRBACGroupManagedSettingsNotFoundErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACGroupOrderLimitExceeded() (v BetaRBACGroupOrderLimitExceededErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsRBACGroupNotFound() (v BetaRBACGroupNotFoundErrorDetails) {
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

func (u ErrorDetailsUnion) AsBundleAttachmentNotSupported() (v BetaBundleAttachmentNotSupportedErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorDetailsUnion) AsBundleOnCall() (v BetaBundleOnCallErrorDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ErrorDetailsUnion) RawJSON() string { return u.JSON.raw }

func (r *ErrorDetailsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ErrorObjectUnion contains all possible properties and values from
// [InvalidRequestErrorObject], [AuthenticationErrorObject], [BillingErrorObject],
// [PermissionErrorObject], [NotFoundErrorObject], [RateLimitErrorObject],
// [GatewayTimeoutErrorObject], [APIErrorObject], [OverloadedErrorObject].
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
//	case shared.InvalidRequestErrorObject:
//	case shared.AuthenticationErrorObject:
//	case shared.BillingErrorObject:
//	case shared.PermissionErrorObject:
//	case shared.NotFoundErrorObject:
//	case shared.RateLimitErrorObject:
//	case shared.GatewayTimeoutErrorObject:
//	case shared.APIErrorObject:
//	case shared.OverloadedErrorObject:
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

func (u ErrorObjectUnion) AsInvalidRequestError() (v InvalidRequestErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsAuthenticationError() (v AuthenticationErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsBillingError() (v BillingErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsPermissionError() (v PermissionErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsNotFoundError() (v NotFoundErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsRateLimitError() (v RateLimitErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsTimeoutError() (v GatewayTimeoutErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsAPIError() (v APIErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ErrorObjectUnion) AsOverloadedError() (v OverloadedErrorObject) {
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

type GatewayTimeoutErrorObject struct {
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
func (r GatewayTimeoutErrorObject) RawJSON() string { return r.JSON.raw }
func (r *GatewayTimeoutErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (GatewayTimeoutErrorObject) ImplErrorObjectUnion() {}

type InvalidRequestErrorObject struct {
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
func (r InvalidRequestErrorObject) RawJSON() string { return r.JSON.raw }
func (r *InvalidRequestErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (InvalidRequestErrorObject) ImplErrorObjectUnion() {}

// Anthropic cannot currently bill the organization linked to the request's user
// profile, for a reason that lies with that organization's own account. The
// response does not say what the reason is.
type LinkedOrganizationBillingUnavailableErrorDetails struct {
	ErrorCode constant.LinkedOrganizationBillingUnavailableErrorCode `json:"error_code" default:"linked_organization_billing_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LinkedOrganizationBillingUnavailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *LinkedOrganizationBillingUnavailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (LinkedOrganizationBillingUnavailableErrorDetails) ImplErrorDetailsUnion() {}

// The organization linked to the request's user profile has reached a spend limit
// that it set, or lacks credit for the request. The response does not say which
// applies.
type LinkedOrganizationSpendLimitReachedErrorDetails struct {
	ErrorCode constant.LinkedOrganizationSpendLimitReachedErrorCode `json:"error_code" default:"linked_organization_spend_limit_reached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LinkedOrganizationSpendLimitReachedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *LinkedOrganizationSpendLimitReachedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (LinkedOrganizationSpendLimitReachedErrorDetails) ImplErrorDetailsUnion() {}

type NotFoundErrorObject struct {
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
func (r NotFoundErrorObject) RawJSON() string { return r.JSON.raw }
func (r *NotFoundErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (NotFoundErrorObject) ImplErrorObjectUnion() {}

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

type OverloadedErrorObject struct {
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
func (r OverloadedErrorObject) RawJSON() string { return r.JSON.raw }
func (r *OverloadedErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (OverloadedErrorObject) ImplErrorObjectUnion() {}

// The endpoint serves only organizations that belong to a parent organization, and
// this organization belongs to none.
type ParentOrganizationRequiredErrorDetails struct {
	ErrorCode constant.ParentOrganizationRequiredErrorCode `json:"error_code" default:"parent_organization_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ParentOrganizationRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ParentOrganizationRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ParentOrganizationRequiredErrorDetails) ImplErrorDetailsUnion() {}

type PermissionErrorObject struct {
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
func (r PermissionErrorObject) RawJSON() string { return r.JSON.raw }
func (r *PermissionErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (PermissionErrorObject) ImplErrorObjectUnion() {}

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

type RateLimitErrorObject struct {
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
func (r RateLimitErrorObject) RawJSON() string { return r.JSON.raw }
func (r *RateLimitErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RateLimitErrorObject) ImplErrorObjectUnion() {}

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

// The request came from a Claude sandbox with restricted network access, which is
// not allowed to make this request. Send it from outside the sandbox, or from a
// sandbox with unrestricted network access.
type SandboxEgressRestrictedErrorDetails struct {
	ErrorCode constant.SandboxEgressRestrictedErrorCode `json:"error_code" default:"sandbox_egress_restricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SandboxEgressRestrictedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *SandboxEgressRestrictedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (SandboxEgressRestrictedErrorDetails) ImplErrorDetailsUnion() {}

// The request's Workspace has Human Review turned off, and this operation is not
// available there. Use a Workspace with Human Review on.
type WorkspaceHumanReviewDisabledErrorDetails struct {
	ErrorCode constant.WorkspaceHumanReviewDisabledErrorCode `json:"error_code" default:"workspace_human_review_disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceHumanReviewDisabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceHumanReviewDisabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (WorkspaceHumanReviewDisabledErrorDetails) ImplErrorDetailsUnion() {}

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

// The longest side of one or more images in the request exceeds the image size
// limit. Scale each image in `paths` so that neither side exceeds
// `max_side_length_pixels`, or remove it, then retry.
type ImageDimensionsTooLargeErrorDetails struct {
	ErrorCode constant.ImageDimensionsTooLargeErrorCode `json:"error_code" default:"image_dimensions_too_large"`
	// The limit the images exceeded: the largest width or height, in pixels, that an
	// image on this request may have.
	MaxSideLengthPixels int64 `json:"max_side_length_pixels" api:"required"`
	// The position of every image over the limit, in the order you sent them, nested
	// blocks included, such as `messages.0.content.1`, or
	// `messages.0.content.1.content.0` for an image in a `tool_result` block.
	Paths []string `json:"paths" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode           respjson.Field
		MaxSideLengthPixels respjson.Field
		Paths               respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ImageDimensionsTooLargeErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ImageDimensionsTooLargeErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ImageDimensionsTooLargeErrorDetails) ImplErrorDetailsUnion() {}

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

// The request sets `billing_policy.bill_to` to `"user_profile"` but names no user
// profile. Name the profile to bill in the `anthropic-user-profile-id` header.
type UserProfileRequiredErrorDetails struct {
	ErrorCode constant.UserProfileRequiredErrorCode `json:"error_code" default:"user_profile_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (UserProfileRequiredErrorDetails) ImplErrorDetailsUnion() {}

// The customer that the request's user profile represents has not authorized the
// caller to bill them, or the authorization that the request's credential uses has
// ended. The response does not say which applies.
type UserProfileBillingNotAuthorizedErrorDetails struct {
	ErrorCode constant.UserProfileBillingNotAuthorizedErrorCode `json:"error_code" default:"user_profile_billing_not_authorized"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileBillingNotAuthorizedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileBillingNotAuthorizedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (UserProfileBillingNotAuthorizedErrorDetails) ImplErrorDetailsUnion() {}

// The request's user profile cannot be billed. The response gives no more specific
// error code and does not say what the reason is.
type UserProfileNotBillableErrorDetails struct {
	ErrorCode constant.UserProfileNotBillableErrorCode `json:"error_code" default:"user_profile_not_billable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UserProfileNotBillableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *UserProfileNotBillableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (UserProfileNotBillableErrorDetails) ImplErrorDetailsUnion() {}

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

// After this request, the External Key and the Workspace would not have the same
// user profile. A key and its Workspace must have the same user profile, or both
// have none.
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
	// The request field the broken rule applies to. Other values may be returned.
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

// The Workspace was created, but some of what the request asked for did not take
// effect, or may not have. Don't send the create again. The error's `message` says
// what to do with the Workspace at `workspace_id`.
type WorkspaceCreatedIncompleteErrorDetails struct {
	ErrorCode constant.WorkspaceCreatedIncompleteErrorCode `json:"error_code" default:"workspace_created_incomplete"`
	// ID of the Workspace the request created.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceCreatedIncompleteErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceCreatedIncompleteErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (WorkspaceCreatedIncompleteErrorDetails) ImplErrorDetailsUnion() {}

// The Workspace's storage and its user profile don't match: the Workspace and its
// storage must have the same user profile, or both have none.
type StorageProfileMismatchErrorDetails struct {
	ErrorCode constant.StorageProfileMismatchErrorCode `json:"error_code" default:"storage_profile_mismatch"`
	// The storage setting that doesn't match the Workspace's user profile.
	Field string `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r StorageProfileMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *StorageProfileMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (StorageProfileMismatchErrorDetails) ImplErrorDetailsUnion() {}

// Some of the request's changes may have been saved and others not. Send the same
// request again to finish it; changes already saved aren't repeated.
type ChangeOutcomeUnknownErrorDetails struct {
	ErrorCode constant.ChangeOutcomeUnknownErrorCode `json:"error_code" default:"change_outcome_unknown"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ChangeOutcomeUnknownErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *ChangeOutcomeUnknownErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (ChangeOutcomeUnknownErrorDetails) ImplErrorDetailsUnion() {}

// The Workspace was archived, but External Storage was attached to it during the
// archive and is still attached. To detach it, archive the Workspace again as a
// caller with the `org:external-storage:use` permission.
type StorageAttachedDuringArchiveErrorDetails struct {
	ErrorCode constant.StorageAttachedDuringArchiveErrorCode `json:"error_code" default:"storage_attached_during_archive"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r StorageAttachedDuringArchiveErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *StorageAttachedDuringArchiveErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (StorageAttachedDuringArchiveErrorDetails) ImplErrorDetailsUnion() {}

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

// No RBAC Role with this ID was found.
type RBACRoleNotFoundErrorDetails struct {
	ErrorCode constant.RBACRoleNotFoundErrorCode `json:"error_code" default:"rbac_role_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RBACRoleNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *RBACRoleNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (RBACRoleNotFoundErrorDetails) ImplErrorDetailsUnion() {}

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

// Anthropic has not approved External Storage for the organization, for the
// Workspace or, for storage bound to a user profile, for that user profile.
// Contact your Anthropic account team.
type BetaGrantRequiredErrorDetails struct {
	ErrorCode constant.BetaGrantRequiredErrorCode `json:"error_code" default:"grant_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGrantRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaGrantRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaGrantRequiredErrorDetails) ImplErrorDetailsUnion() {}

// The feature isn't offered to this organization. The error's `message` names the
// kind of organization or, when it names none, says to contact your Anthropic
// account team.
type BetaOrganizationNotEligibleErrorDetails struct {
	ErrorCode constant.BetaOrganizationNotEligibleErrorCode `json:"error_code" default:"organization_not_eligible"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationNotEligibleErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationNotEligibleErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaOrganizationNotEligibleErrorDetails) ImplErrorDetailsUnion() {}

// The request can't use External Storage on this cloud provider. The error's
// message says why.
type BetaProviderNotAllowedErrorDetails struct {
	ErrorCode constant.BetaProviderNotAllowedErrorCode `json:"error_code" default:"provider_not_allowed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaProviderNotAllowedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaProviderNotAllowedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaProviderNotAllowedErrorDetails) ImplErrorDetailsUnion() {}

// `provider_config.region` is not a region where External Storage is available.
type BetaUnsupportedRegionErrorDetails struct {
	ErrorCode constant.BetaUnsupportedRegionErrorCode `json:"error_code" default:"unsupported_region"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUnsupportedRegionErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaUnsupportedRegionErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaUnsupportedRegionErrorDetails) ImplErrorDetailsUnion() {}

// The organization already has External Storage for this bucket that it has not
// deleted. List the organization's External Storage to find it.
type BetaStorageAlreadyRegisteredErrorDetails struct {
	ErrorCode constant.BetaStorageAlreadyRegisteredErrorCode `json:"error_code" default:"storage_already_registered"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageAlreadyRegisteredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageAlreadyRegisteredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageAlreadyRegisteredErrorDetails) ImplErrorDetailsUnion() {}

// One of the organization's approvals from Anthropic could not be read right now,
// so nothing was done. Retry after the `Retry-After` delay.
type BetaAuthorizationCheckUnavailableErrorDetails struct {
	ErrorCode constant.BetaAuthorizationCheckUnavailableErrorCode `json:"error_code" default:"authorization_check_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAuthorizationCheckUnavailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaAuthorizationCheckUnavailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaAuthorizationCheckUnavailableErrorDetails) ImplErrorDetailsUnion() {}

// On Claude Platform on AWS, the External Storage ID in the request is neither a
// bucket ARN nor an `estrg_` ID, or is the `estrg_` ID of AWS storage. Name AWS
// storage by its bucket ARN, `arn:<partition>:s3:::<bucket_name>`, which responses
// return as its `id`.
type BetaInvalidStorageARNErrorDetails struct {
	ErrorCode constant.BetaInvalidStorageARNErrorCode `json:"error_code" default:"invalid_storage_arn"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInvalidStorageARNErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaInvalidStorageARNErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaInvalidStorageARNErrorDetails) ImplErrorDetailsUnion() {}

// The External Storage is `active` or `draining`, so this delete didn't erase it.
// The message says what to do next. On `draining` storage this delete may have
// started a removal that can't be undone: get the storage and read
// `status.removal_started`.
type BetaExternalStorageInUseErrorDetails struct {
	ErrorCode constant.BetaExternalStorageInUseErrorCode `json:"error_code" default:"external_storage_in_use"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorageInUseErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorageInUseErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaExternalStorageInUseErrorDetails) ImplErrorDetailsUnion() {}

// The External Storage can't be deleted while an Enterprise Frontier Safeguards
// workspace policy lists it as allowed storage.
type BetaExternalStorageInWorkspacePolicyErrorDetails struct {
	ErrorCode constant.BetaExternalStorageInWorkspacePolicyErrorCode `json:"error_code" default:"external_storage_in_workspace_policy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorageInWorkspacePolicyErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorageInWorkspacePolicyErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaExternalStorageInWorkspacePolicyErrorDetails) ImplErrorDetailsUnion() {}

// The External Storage failed verification, so the request changed nothing.
// `failed_checks` lists what to fix.
type BetaStorageVerificationFailedErrorDetails struct {
	ErrorCode constant.BetaStorageVerificationFailedErrorCode `json:"error_code" default:"storage_verification_failed"`
	// One entry for each check that failed.
	FailedChecks []BetaExternalStorageFailedCheck `json:"failed_checks" api:"required"`
	// The request field at fault. Sent when the refused request is a Workspace's
	// create or update.
	Field string `json:"field" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode    respjson.Field
		FailedChecks respjson.Field
		Field        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageVerificationFailedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageVerificationFailedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageVerificationFailedErrorDetails) ImplErrorDetailsUnion() {}

// One check that an External Storage's verification failed.
type BetaExternalStorageFailedCheck struct {
	// Name of the check that failed, such as `versioning`, for code to branch on. New
	// names may be added; `other` is a check this API does not name. The value is a
	// name with no data of its own.
	//
	// Any of "already_registered", "assume_role", "bucket_empty",
	// "bucket_not_confirmed", "bucket_policy", "bucket_settings",
	// "default_encryption", "deletion_policy", "encryption", "grant_mismatch",
	// "key_usable", "kms_key_refused", "lifecycle_no_archive", "lifecycle_rules",
	// "lifecycle_template", "location_mismatch", "object_lock", "other",
	// "ownership_enforced", "probe_failed", "probe_object", "probe_storage_class",
	// "public_access_block", "region", "replication", "request_metrics",
	// "request_payment", "transcripts_support", "trust_ingest_role",
	// "trust_maintenance_role", "trust_media_erasure_role", "trust_read_role",
	// "trust_transcripts_role", "trust_transcripts_verify_role", "trust_verify_role",
	// "versioning".
	Check BetaExternalStorageFailedCheckCheck `json:"check" api:"required"`
	// What failed and how to fix it, in plain words.
	Message string `json:"message" api:"required"`
	// The part of your access grant that the check is about, such as `ingest`: on AWS,
	// one of the IAM roles your setup created; on Google Cloud, the grant to one of
	// Anthropic's service accounts. `null` when the check is not about one. New names
	// may be added. The value is a name with no data of its own.
	//
	// Any of "ingest", "maintenance", "media_erasure", "read", "transcripts",
	// "transcripts_verify", "verify".
	Role BetaExternalStorageFailedCheckRole `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Check       respjson.Field
		Message     respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorageFailedCheck) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorageFailedCheck) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Name of the check that failed, such as `versioning`, for code to branch on. New
// names may be added; `other` is a check this API does not name. The value is a
// name with no data of its own.
type BetaExternalStorageFailedCheckCheck string

const (
	BetaExternalStorageFailedCheckCheckAlreadyRegistered          BetaExternalStorageFailedCheckCheck = "already_registered"
	BetaExternalStorageFailedCheckCheckAssumeRole                 BetaExternalStorageFailedCheckCheck = "assume_role"
	BetaExternalStorageFailedCheckCheckBucketEmpty                BetaExternalStorageFailedCheckCheck = "bucket_empty"
	BetaExternalStorageFailedCheckCheckBucketNotConfirmed         BetaExternalStorageFailedCheckCheck = "bucket_not_confirmed"
	BetaExternalStorageFailedCheckCheckBucketPolicy               BetaExternalStorageFailedCheckCheck = "bucket_policy"
	BetaExternalStorageFailedCheckCheckBucketSettings             BetaExternalStorageFailedCheckCheck = "bucket_settings"
	BetaExternalStorageFailedCheckCheckDefaultEncryption          BetaExternalStorageFailedCheckCheck = "default_encryption"
	BetaExternalStorageFailedCheckCheckDeletionPolicy             BetaExternalStorageFailedCheckCheck = "deletion_policy"
	BetaExternalStorageFailedCheckCheckEncryption                 BetaExternalStorageFailedCheckCheck = "encryption"
	BetaExternalStorageFailedCheckCheckGrantMismatch              BetaExternalStorageFailedCheckCheck = "grant_mismatch"
	BetaExternalStorageFailedCheckCheckKeyUsable                  BetaExternalStorageFailedCheckCheck = "key_usable"
	BetaExternalStorageFailedCheckCheckKMSKeyRefused              BetaExternalStorageFailedCheckCheck = "kms_key_refused"
	BetaExternalStorageFailedCheckCheckLifecycleNoArchive         BetaExternalStorageFailedCheckCheck = "lifecycle_no_archive"
	BetaExternalStorageFailedCheckCheckLifecycleRules             BetaExternalStorageFailedCheckCheck = "lifecycle_rules"
	BetaExternalStorageFailedCheckCheckLifecycleTemplate          BetaExternalStorageFailedCheckCheck = "lifecycle_template"
	BetaExternalStorageFailedCheckCheckLocationMismatch           BetaExternalStorageFailedCheckCheck = "location_mismatch"
	BetaExternalStorageFailedCheckCheckObjectLock                 BetaExternalStorageFailedCheckCheck = "object_lock"
	BetaExternalStorageFailedCheckCheckOther                      BetaExternalStorageFailedCheckCheck = "other"
	BetaExternalStorageFailedCheckCheckOwnershipEnforced          BetaExternalStorageFailedCheckCheck = "ownership_enforced"
	BetaExternalStorageFailedCheckCheckProbeFailed                BetaExternalStorageFailedCheckCheck = "probe_failed"
	BetaExternalStorageFailedCheckCheckProbeObject                BetaExternalStorageFailedCheckCheck = "probe_object"
	BetaExternalStorageFailedCheckCheckProbeStorageClass          BetaExternalStorageFailedCheckCheck = "probe_storage_class"
	BetaExternalStorageFailedCheckCheckPublicAccessBlock          BetaExternalStorageFailedCheckCheck = "public_access_block"
	BetaExternalStorageFailedCheckCheckRegion                     BetaExternalStorageFailedCheckCheck = "region"
	BetaExternalStorageFailedCheckCheckReplication                BetaExternalStorageFailedCheckCheck = "replication"
	BetaExternalStorageFailedCheckCheckRequestMetrics             BetaExternalStorageFailedCheckCheck = "request_metrics"
	BetaExternalStorageFailedCheckCheckRequestPayment             BetaExternalStorageFailedCheckCheck = "request_payment"
	BetaExternalStorageFailedCheckCheckTranscriptsSupport         BetaExternalStorageFailedCheckCheck = "transcripts_support"
	BetaExternalStorageFailedCheckCheckTrustIngestRole            BetaExternalStorageFailedCheckCheck = "trust_ingest_role"
	BetaExternalStorageFailedCheckCheckTrustMaintenanceRole       BetaExternalStorageFailedCheckCheck = "trust_maintenance_role"
	BetaExternalStorageFailedCheckCheckTrustMediaErasureRole      BetaExternalStorageFailedCheckCheck = "trust_media_erasure_role"
	BetaExternalStorageFailedCheckCheckTrustReadRole              BetaExternalStorageFailedCheckCheck = "trust_read_role"
	BetaExternalStorageFailedCheckCheckTrustTranscriptsRole       BetaExternalStorageFailedCheckCheck = "trust_transcripts_role"
	BetaExternalStorageFailedCheckCheckTrustTranscriptsVerifyRole BetaExternalStorageFailedCheckCheck = "trust_transcripts_verify_role"
	BetaExternalStorageFailedCheckCheckTrustVerifyRole            BetaExternalStorageFailedCheckCheck = "trust_verify_role"
	BetaExternalStorageFailedCheckCheckVersioning                 BetaExternalStorageFailedCheckCheck = "versioning"
)

// The part of your access grant that the check is about, such as `ingest`: on AWS,
// one of the IAM roles your setup created; on Google Cloud, the grant to one of
// Anthropic's service accounts. `null` when the check is not about one. New names
// may be added. The value is a name with no data of its own.
type BetaExternalStorageFailedCheckRole string

const (
	BetaExternalStorageFailedCheckRoleIngest            BetaExternalStorageFailedCheckRole = "ingest"
	BetaExternalStorageFailedCheckRoleMaintenance       BetaExternalStorageFailedCheckRole = "maintenance"
	BetaExternalStorageFailedCheckRoleMediaErasure      BetaExternalStorageFailedCheckRole = "media_erasure"
	BetaExternalStorageFailedCheckRoleRead              BetaExternalStorageFailedCheckRole = "read"
	BetaExternalStorageFailedCheckRoleTranscripts       BetaExternalStorageFailedCheckRole = "transcripts"
	BetaExternalStorageFailedCheckRoleTranscriptsVerify BetaExternalStorageFailedCheckRole = "transcripts_verify"
	BetaExternalStorageFailedCheckRoleVerify            BetaExternalStorageFailedCheckRole = "verify"
)

// Anthropic has not approved customer-managed encryption keys for the
// organization, project or user profile this request applies to, or has paused or
// withdrawn its approval. If you did not expect this, contact Anthropic.
type BetaCMEKNotApprovedErrorDetails struct {
	ErrorCode constant.BetaCMEKNotApprovedErrorCode `json:"error_code" default:"cmek_not_approved"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCMEKNotApprovedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCMEKNotApprovedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCMEKNotApprovedErrorDetails) ImplErrorDetailsUnion() {}

// The organization has no active External Key with this ID. The message names the
// field.
type BetaExternalKeyNotFoundErrorDetails struct {
	ErrorCode constant.BetaExternalKeyNotFoundErrorCode `json:"error_code" default:"external_key_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalKeyNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalKeyNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaExternalKeyNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// The External Key is held by a key management service that this endpoint does not
// take. The message names the field and the service it takes.
type BetaExternalKeyProviderMismatchErrorDetails struct {
	ErrorCode constant.BetaExternalKeyProviderMismatchErrorCode `json:"error_code" default:"external_key_provider_mismatch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalKeyProviderMismatchErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalKeyProviderMismatchErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaExternalKeyProviderMismatchErrorDetails) ImplErrorDetailsUnion() {}

// The organization has not proved that it owns the Google Cloud project that holds
// the key or bucket. The message names the field.
type BetaGCPProjectNotVerifiedErrorDetails struct {
	ErrorCode constant.BetaGCPProjectNotVerifiedErrorCode `json:"error_code" default:"gcp_project_not_verified"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGCPProjectNotVerifiedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaGCPProjectNotVerifiedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaGCPProjectNotVerifiedErrorDetails) ImplErrorDetailsUnion() {}

// The request names a region code that the control does not take from this
// organization. The message names the field.
type BetaRegionNotSupportedErrorDetails struct {
	ErrorCode constant.BetaRegionNotSupportedErrorCode `json:"error_code" default:"region_not_supported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRegionNotSupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRegionNotSupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRegionNotSupportedErrorDetails) ImplErrorDetailsUnion() {}

// Anthropic has not approved keeping data in a customer-owned storage bucket for
// the organization, project or user profile this request applies to, or has paused
// or withdrawn its approval. If you did not expect this, contact Anthropic.
type BetaStorageNotApprovedErrorDetails struct {
	ErrorCode constant.BetaStorageNotApprovedErrorCode `json:"error_code" default:"storage_not_approved"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageNotApprovedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageNotApprovedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageNotApprovedErrorDetails) ImplErrorDetailsUnion() {}

// Another organization has data controls turned on for this Google Cloud project,
// or turned them off recently. The message says how long to wait.
type BetaGCPProjectControlledByAnotherOrganizationErrorDetails struct {
	ErrorCode constant.BetaGCPProjectControlledByAnotherOrganizationErrorCode `json:"error_code" default:"gcp_project_controlled_by_another_organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGCPProjectControlledByAnotherOrganizationErrorDetails) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaGCPProjectControlledByAnotherOrganizationErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaGCPProjectControlledByAnotherOrganizationErrorDetails) ImplErrorDetailsUnion() {}

// The bucket is already in use, or was in use recently. Name a different bucket.
type BetaStorageInUseErrorDetails struct {
	ErrorCode constant.BetaStorageInUseErrorCode `json:"error_code" default:"storage_in_use"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageInUseErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageInUseErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageInUseErrorDetails) ImplErrorDetailsUnion() {}

// The member is the organization's primary owner, whom no API key or token can
// sign out.
type BetaCannotSignOutPrimaryOwnerErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutPrimaryOwnerErrorCode `json:"error_code" default:"cannot_sign_out_primary_owner"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutPrimaryOwnerErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutPrimaryOwnerErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutPrimaryOwnerErrorDetails) ImplErrorDetailsUnion() {}

// The member is the user the API key or token belongs to.
type BetaCannotSignOutSelfErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutSelfErrorCode `json:"error_code" default:"cannot_sign_out_self"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutSelfErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutSelfErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutSelfErrorDetails) ImplErrorDetailsUnion() {}

// The beta that the endpoint requires is not enabled for the organization. The
// message names the beta value.
type BetaOrganizationNotEnrolledInBetaErrorDetails struct {
	ErrorCode constant.BetaOrganizationNotEnrolledInBetaErrorCode `json:"error_code" default:"organization_not_enrolled_in_beta"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationNotEnrolledInBetaErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationNotEnrolledInBetaErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaOrganizationNotEnrolledInBetaErrorDetails) ImplErrorDetailsUnion() {}

// The member has the `managed` role, and this API key or token cannot sign them
// out. The error message says which keys and tokens can.
type BetaCannotSignOutManagedMemberErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutManagedMemberErrorCode `json:"error_code" default:"cannot_sign_out_managed_member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutManagedMemberErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutManagedMemberErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutManagedMemberErrorDetails) ImplErrorDetailsUnion() {}

// The member is a membership admin, and the person the API key or token belongs to
// is not allowed to sign out membership admins.
type BetaCannotSignOutMembershipAdminErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutMembershipAdminErrorCode `json:"error_code" default:"cannot_sign_out_membership_admin"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutMembershipAdminErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutMembershipAdminErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutMembershipAdminErrorDetails) ImplErrorDetailsUnion() {}

// The member is an owner or has owner permissions, and this API key or token
// cannot sign them out. In a Claude Enterprise organization, an API key with the
// `write:members` scope can.
type BetaCannotSignOutOwnerErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutOwnerErrorCode `json:"error_code" default:"cannot_sign_out_owner"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutOwnerErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutOwnerErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutOwnerErrorDetails) ImplErrorDetailsUnion() {}

// This OAuth token cannot sign out a membership admin or a member whose roles
// grant owner, user management, or identity and access management.
type BetaCannotSignOutPrivilegedMemberErrorDetails struct {
	ErrorCode constant.BetaCannotSignOutPrivilegedMemberErrorCode `json:"error_code" default:"cannot_sign_out_privileged_member"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCannotSignOutPrivilegedMemberErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCannotSignOutPrivilegedMemberErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCannotSignOutPrivilegedMemberErrorDetails) ImplErrorDetailsUnion() {}

// Turning human review off needs Anthropic's approval for the organization, and
// the organization doesn't have it. Contact Anthropic support to request approval.
type BetaApprovalRequiredErrorDetails struct {
	ErrorCode constant.BetaApprovalRequiredErrorCode `json:"error_code" default:"approval_required"`
	// The request field that needs the approval.
	Field string `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaApprovalRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaApprovalRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaApprovalRequiredErrorDetails) ImplErrorDetailsUnion() {}

// External Storage isn't available for this Workspace, so only `internal` storage
// is accepted. The error's `message` says why.
type BetaExternalStorageNotAvailableErrorDetails struct {
	ErrorCode constant.BetaExternalStorageNotAvailableErrorCode `json:"error_code" default:"external_storage_not_available"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorageNotAvailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorageNotAvailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaExternalStorageNotAvailableErrorDetails) ImplErrorDetailsUnion() {}

// Human review can't be set as requested on this Workspace, so the request changed
// nothing. The error's `message` says what to do instead.
type BetaHumanReviewNotAvailableErrorDetails struct {
	ErrorCode constant.BetaHumanReviewNotAvailableErrorCode `json:"error_code" default:"human_review_not_available"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaHumanReviewNotAvailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaHumanReviewNotAvailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaHumanReviewNotAvailableErrorDetails) ImplErrorDetailsUnion() {}

// The request would leave the Workspace with
// `inference_data_privacy.safeguards_storage` set to External Storage and
// `inference_data_privacy.human_review` enabled. Anthropic doesn't offer human
// review of data in a customer's own storage. Set `human_review` to `disabled`, or
// `safeguards_storage` to `internal`.
type BetaHumanReviewOnWithExternalStorageErrorDetails struct {
	ErrorCode constant.BetaHumanReviewOnWithExternalStorageErrorCode `json:"error_code" default:"human_review_on_with_external_storage"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaHumanReviewOnWithExternalStorageErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaHumanReviewOnWithExternalStorageErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaHumanReviewOnWithExternalStorageErrorDetails) ImplErrorDetailsUnion() {}

// This External Storage can't be attached to the setting named in `field`. The
// error's `message` says why.
type BetaStorageNotAttachableErrorDetails struct {
	ErrorCode constant.BetaStorageNotAttachableErrorCode `json:"error_code" default:"storage_not_attachable"`
	// The request field at fault.
	Field string `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageNotAttachableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageNotAttachableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageNotAttachableErrorDetails) ImplErrorDetailsUnion() {}

// In an organization managed through AWS or Google Cloud, the request would leave
// the Workspace with `inference_data_privacy.human_review` disabled and
// `inference_data_privacy.transcripts_storage` not set to External Storage. With
// human review disabled, the organization reviews its own flagged requests, and it
// receives them only in that storage.
type BetaTranscriptsStorageRequiredErrorDetails struct {
	ErrorCode constant.BetaTranscriptsStorageRequiredErrorCode `json:"error_code" default:"transcripts_storage_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTranscriptsStorageRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaTranscriptsStorageRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaTranscriptsStorageRequiredErrorDetails) ImplErrorDetailsUnion() {}

// The Workspace has reached its limit of 100 different External Storage resources
// attached over its lifetime. Detached storage counts toward the limit until it is
// deleted.
type BetaStorageAttachLimitReachedErrorDetails struct {
	ErrorCode constant.BetaStorageAttachLimitReachedErrorCode `json:"error_code" default:"storage_attach_limit_reached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageAttachLimitReachedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageAttachLimitReachedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageAttachLimitReachedErrorDetails) ImplErrorDetailsUnion() {}

// This storage is bound to a different user profile from storage the Workspace has
// attached, or from active or draining storage it once had for safeguards data.
// Detaching storage the Workspace uses only for transcripts clears this. Detaching
// safeguards storage doesn't, so use a new Workspace instead.
type BetaStorageProfileStillRequiredErrorDetails struct {
	ErrorCode constant.BetaStorageProfileStillRequiredErrorCode `json:"error_code" default:"storage_profile_still_required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageProfileStillRequiredErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageProfileStillRequiredErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageProfileStillRequiredErrorDetails) ImplErrorDetailsUnion() {}

// Another request changed this storage setting while this one was running, so this
// request left it as it is. Read the Workspace and send the request again if you
// still want the change.
type BetaStorageChangedConcurrentlyErrorDetails struct {
	ErrorCode constant.BetaStorageChangedConcurrentlyErrorCode `json:"error_code" default:"storage_changed_concurrently"`
	// The storage setting that changed.
	Field string `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaStorageChangedConcurrentlyErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaStorageChangedConcurrentlyErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaStorageChangedConcurrentlyErrorDetails) ImplErrorDetailsUnion() {}

// Another request changed the Workspace, so this request left the setting named in
// `field` as it is. The error's `message` says which of this request's other
// changes are in place. Read the Workspace and send the request again if you still
// want the change.
type BetaWorkspaceChangedConcurrentlyErrorDetails struct {
	ErrorCode constant.BetaWorkspaceChangedConcurrentlyErrorCode `json:"error_code" default:"workspace_changed_concurrently"`
	// The request field whose change was refused. Other values may be returned. The
	// value is a name with no data of its own.
	//
	// Any of "inference_data_privacy.human_review".
	Field BetaWorkspaceChangedConcurrentlyField `json:"field" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		Field       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkspaceChangedConcurrentlyErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkspaceChangedConcurrentlyErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaWorkspaceChangedConcurrentlyErrorDetails) ImplErrorDetailsUnion() {}

type BetaWorkspaceChangedConcurrentlyField string

const (
	BetaWorkspaceChangedConcurrentlyFieldInferenceDataPrivacyHumanReview BetaWorkspaceChangedConcurrentlyField = "inference_data_privacy.human_review"
)

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

// The request body did not finish arriving in time, so nothing was done. Send a
// smaller body, or send it over a faster connection.
type BetaRequestBodyTimeoutErrorDetails struct {
	ErrorCode constant.BetaRequestBodyTimeoutErrorCode `json:"error_code" default:"request_body_timeout"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRequestBodyTimeoutErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRequestBodyTimeoutErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRequestBodyTimeoutErrorDetails) ImplErrorDetailsUnion() {}

// The organization has reached its limit of settings documents, counting unused
// ones. An owner can remove unused documents in claude.ai. Deleting a group's
// settings or clearing the default settings also frees a document, if nothing else
// uses it. None of these works once the organization has also reached its settings
// history limit.
type BetaSettingsDocumentLimitExceededErrorDetails struct {
	ErrorCode constant.BetaSettingsDocumentLimitExceededErrorCode `json:"error_code" default:"settings_document_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSettingsDocumentLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSettingsDocumentLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSettingsDocumentLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// The organization has reached its settings history limit, a cap on stored
// versions of its settings documents. Nothing in the API or claude.ai lowers the
// count, so contact support. Creating or deleting group settings, creating or
// clearing the default settings, and reordering are refused. Changing the mode
// still works, unless it would copy the organization settings into a new default.
type BetaSettingsHistoryLimitExceededErrorDetails struct {
	ErrorCode constant.BetaSettingsHistoryLimitExceededErrorCode `json:"error_code" default:"settings_history_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSettingsHistoryLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaSettingsHistoryLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaSettingsHistoryLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// This stored settings document currently cannot be read or replaced. Contact
// support, or, if you no longer need these settings, delete the group's settings
// or clear the default.
type BetaCMEKFeatureUnsupportedErrorDetails struct {
	ErrorCode constant.BetaCMEKFeatureUnsupportedErrorCode `json:"error_code" default:"cmek_feature_unsupported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCMEKFeatureUnsupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaCMEKFeatureUnsupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaCMEKFeatureUnsupportedErrorDetails) ImplErrorDetailsUnion() {}

// This RBAC Group has no Claude Code managed settings.
type BetaRBACGroupManagedSettingsNotFoundErrorDetails struct {
	ErrorCode constant.BetaRBACGroupManagedSettingsNotFoundErrorCode `json:"error_code" default:"rbac_group_managed_settings_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACGroupManagedSettingsNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACGroupManagedSettingsNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRBACGroupManagedSettingsNotFoundErrorDetails) ImplErrorDetailsUnion() {}

// The RBAC Group order already has the most groups allowed. Delete another group's
// settings first.
type BetaRBACGroupOrderLimitExceededErrorDetails struct {
	ErrorCode constant.BetaRBACGroupOrderLimitExceededErrorCode `json:"error_code" default:"rbac_group_order_limit_exceeded"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACGroupOrderLimitExceededErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACGroupOrderLimitExceededErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRBACGroupOrderLimitExceededErrorDetails) ImplErrorDetailsUnion() {}

// The RBAC Group is not one of the organization's groups.
type BetaRBACGroupNotFoundErrorDetails struct {
	ErrorCode constant.BetaRBACGroupNotFoundErrorCode `json:"error_code" default:"rbac_group_not_found"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACGroupNotFoundErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACGroupNotFoundErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaRBACGroupNotFoundErrorDetails) ImplErrorDetailsUnion() {}

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

// The Bundle has an attachment this API does not support editing, so its
// instructions can be changed only in claude.ai.
type BetaBundleAttachmentNotSupportedErrorDetails struct {
	ErrorCode constant.BetaBundleAttachmentNotSupportedErrorCode `json:"error_code" default:"bundle_attachment_not_supported"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaBundleAttachmentNotSupportedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaBundleAttachmentNotSupportedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaBundleAttachmentNotSupportedErrorDetails) ImplErrorDetailsUnion() {}

// This Bundle's instructions are managed in claude.ai and can't be changed through
// the API.
type BetaBundleOnCallErrorDetails struct {
	ErrorCode constant.BetaBundleOnCallErrorCode `json:"error_code" default:"bundle_on_call"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaBundleOnCallErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaBundleOnCallErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (BetaBundleOnCallErrorDetails) ImplErrorDetailsUnion() {}

// Deprecated: AuthenticationError has been renamed to AuthenticationErrorObject
type AuthenticationError = AuthenticationErrorObject

// Deprecated: InvalidRequestError has been renamed to InvalidRequestErrorObject
type InvalidRequestError = InvalidRequestErrorObject

// Deprecated: BillingError has been renamed to BillingErrorObject
type BillingError = BillingErrorObject

// Deprecated: GatewayTimeoutError has been renamed to GatewayTimeoutErrorObject
type GatewayTimeoutError = GatewayTimeoutErrorObject

// Deprecated: NotFoundError has been renamed to NotFoundErrorObject
type NotFoundError = NotFoundErrorObject

// Deprecated: OverloadedError has been renamed to OverloadedErrorObject
type OverloadedError = OverloadedErrorObject

// Deprecated: PermissionError has been renamed to PermissionErrorObject
type PermissionError = PermissionErrorObject

// Deprecated: RateLimitError has been renamed to RateLimitErrorObject
type RateLimitError = RateLimitErrorObject

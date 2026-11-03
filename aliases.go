package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apierror"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/shared"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type Error = apierror.Error

// This is an alias to an internal type.
type APIErrorObject = shared.APIErrorObject

// This is an alias to an internal type.
type AuthenticationError = shared.AuthenticationError

// This is an alias to an internal type.
type BillingError = shared.BillingError

// The API key belongs to the Claude Code workspace, and the user who created it is
// no longer a member of the organization or of that workspace.
//
// This is an alias to an internal type.
type ClaudeCodeKeyCreatorNotMemberErrorDetails = shared.ClaudeCodeKeyCreatorNotMemberErrorDetails

// The request comes from a Claude Code version older than the minimum that the
// requested model or the organization requires. Update Claude Code, then retry.
//
// This is an alias to an internal type.
type ClaudeCodeVersionTooOldErrorDetails = shared.ClaudeCodeVersionTooOldErrorDetails

// The request's workspace encryption context could not be resolved, so the request
// was refused. Retry in a few minutes.
//
// This is an alias to an internal type.
type CmekContextUnavailableErrorDetails = shared.CmekContextUnavailableErrorDetails

// The workspace's customer-managed encryption key is disabled or revoked.
// Re-enable the key, then retry.
//
// This is an alias to an internal type.
type CmekKeyDisabledErrorDetails = shared.CmekKeyDisabledErrorDetails

// The key management service's network access controls blocked the request to the
// workspace's customer-managed encryption key. The key itself may be healthy.
//
// This is an alias to an internal type.
type CmekKeyNetworkBlockedErrorDetails = shared.CmekKeyNetworkBlockedErrorDetails

// The workspace's customer-managed encryption key could not be reached. Retry the
// request; if this continues, check that the key is available in your key
// management service, or contact Anthropic support.
//
// This is an alias to an internal type.
type CmekKeyUnreachableErrorDetails = shared.CmekKeyUnreachableErrorDetails

// The organization has crossed the monthly API usage threshold of its API tier.
//
// This is an alias to an internal type.
type EnforcedSpendLimitReachedErrorDetails = shared.EnforcedSpendLimitReachedErrorDetails

// This is an alias to an internal type.
type ErrorObjectUnion = shared.ErrorObjectUnion

// This is an alias to an internal type.
type ErrorResponse = shared.ErrorResponse

// This is an alias to an internal type.
type ErrorType = shared.ErrorType

// Equals "invalid_request_error"
const ErrorTypeInvalidRequestError = shared.ErrorTypeInvalidRequestError

// Equals "authentication_error"
const ErrorTypeAuthenticationError = shared.ErrorTypeAuthenticationError

// Equals "billing_error"
const ErrorTypeBillingError = shared.ErrorTypeBillingError

// Equals "permission_error"
const ErrorTypePermissionError = shared.ErrorTypePermissionError

// Equals "not_found_error"
const ErrorTypeNotFoundError = shared.ErrorTypeNotFoundError

// Equals "rate_limit_error"
const ErrorTypeRateLimitError = shared.ErrorTypeRateLimitError

// Equals "timeout_error"
const ErrorTypeTimeoutError = shared.ErrorTypeTimeoutError

// Equals "api_error"
const ErrorTypeAPIError = shared.ErrorTypeAPIError

// Equals "overloaded_error"
const ErrorTypeOverloadedError = shared.ErrorTypeOverloadedError

// Equals "memory_precondition_failed_error"
const ErrorTypeMemoryPreconditionFailedError = shared.ErrorTypeMemoryPreconditionFailedError

// Equals "memory_path_conflict_error"
const ErrorTypeMemoryPathConflictError = shared.ErrorTypeMemoryPathConflictError

// Equals "conflict_error"
const ErrorTypeConflictError = shared.ErrorTypeConflictError

// This is an alias to an internal type.
type GatewayTimeoutError = shared.GatewayTimeoutError

// This is an alias to an internal type.
type InvalidRequestError = shared.InvalidRequestError

// This is an alias to an internal type.
type NotFoundError = shared.NotFoundError

// Anthropic has disabled the organization for trust and safety reasons. An
// organization disabled for any other reason, such as billing, does not get this
// code. `appeal_url`, when present, is where the organization can appeal.
//
// This is an alias to an internal type.
type OrganizationOnHoldErrorDetails = shared.OrganizationOnHoldErrorDetails

// The endpoint does not support the organization's type. Each endpoint's
// description says which organization types it supports.
//
// This is an alias to an internal type.
type OrganizationTypeNotSupportedErrorDetails = shared.OrganizationTypeNotSupportedErrorDetails

// This is an alias to an internal type.
type OverloadedError = shared.OverloadedError

// This is an alias to an internal type.
type PermissionError = shared.PermissionError

// This is an alias to an internal type.
type RateLimitError = shared.RateLimitError

// The request body is not valid for its `Content-Encoding`.
//
// This is an alias to an internal type.
type RequestBodyEncodingInvalidErrorDetails = shared.RequestBodyEncodingInvalidErrorDetails

// The request's `Content-Encoding` is not supported.
//
// This is an alias to an internal type.
type RequestBodyEncodingUnsupportedErrorDetails = shared.RequestBodyEncodingUnsupportedErrorDetails

// The compressed request body exceeds the request size limit once decoded.
//
// This is an alias to an internal type.
type RequestBodyTooLargeAfterDecodingErrorDetails = shared.RequestBodyTooLargeAfterDecodingErrorDetails

// Backward-compatible type aliases for renamed types.
// These ensure existing code continues to compile after
// UserLocation and ErrorCode types were consolidated.

type WebSearchToolRequestErrorErrorCode = WebSearchToolResultErrorCode
type WebSearchToolResultErrorErrorCode = WebSearchToolResultErrorCode
type WebSearchTool20250305UserLocationParam = UserLocationParam
type WebSearchTool20260209UserLocationParam = UserLocationParam

const WebSearchToolRequestErrorErrorCodeInvalidToolInput = WebSearchToolResultErrorCodeInvalidToolInput
const WebSearchToolRequestErrorErrorCodeUnavailable = WebSearchToolResultErrorCodeUnavailable
const WebSearchToolRequestErrorErrorCodeMaxUsesExceeded = WebSearchToolResultErrorCodeMaxUsesExceeded
const WebSearchToolRequestErrorErrorCodeTooManyRequests = WebSearchToolResultErrorCodeTooManyRequests
const WebSearchToolRequestErrorErrorCodeQueryTooLong = WebSearchToolResultErrorCodeQueryTooLong
const WebSearchToolRequestErrorErrorCodeRequestTooLarge = WebSearchToolResultErrorCodeRequestTooLarge
const WebSearchToolResultErrorErrorCodeInvalidToolInput = WebSearchToolResultErrorCodeInvalidToolInput
const WebSearchToolResultErrorErrorCodeUnavailable = WebSearchToolResultErrorCodeUnavailable
const WebSearchToolResultErrorErrorCodeMaxUsesExceeded = WebSearchToolResultErrorCodeMaxUsesExceeded
const WebSearchToolResultErrorErrorCodeTooManyRequests = WebSearchToolResultErrorCodeTooManyRequests
const WebSearchToolResultErrorErrorCodeQueryTooLong = WebSearchToolResultErrorCodeQueryTooLong
const WebSearchToolResultErrorErrorCodeRequestTooLarge = WebSearchToolResultErrorCodeRequestTooLarge

// Deprecated: use [BetaCurrencyUSD]. Kept so code written before the USD
// initialism casing continues to compile.
const BetaCurrencyUsd = BetaCurrencyUSD

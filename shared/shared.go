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
	Message string            `json:"message" api:"required"`
	Type    constant.APIError `json:"type" default:"api_error"`
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
	Message string                       `json:"message" api:"required"`
	Type    constant.AuthenticationError `json:"type" default:"authentication_error"`
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
	Message string                `json:"message" api:"required"`
	Type    constant.BillingError `json:"type" default:"billing_error"`
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
	ErrorCode constant.ClaudeCodeKeyCreatorNotMember `json:"error_code" default:"claude_code_key_creator_not_member"`
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

// The request's workspace encryption context could not be resolved, so the request
// was refused. Retry in a few minutes.
type CmekContextUnavailableErrorDetails struct {
	ErrorCode constant.CmekContextUnavailable `json:"error_code" default:"cmek_context_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CmekContextUnavailableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CmekContextUnavailableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The workspace's customer-managed encryption key is disabled or revoked.
// Re-enable the key, then retry.
type CmekKeyDisabledErrorDetails struct {
	ErrorCode constant.CmekKeyDisabled `json:"error_code" default:"cmek_key_disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CmekKeyDisabledErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CmekKeyDisabledErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The key management service's network access controls blocked the request to the
// workspace's customer-managed encryption key. The key itself may be healthy.
type CmekKeyNetworkBlockedErrorDetails struct {
	ErrorCode constant.CmekKeyNetworkBlocked `json:"error_code" default:"cmek_key_network_blocked"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CmekKeyNetworkBlockedErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CmekKeyNetworkBlockedErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The workspace's customer-managed encryption key could not be reached. Retry the
// request; if this continues, check that the key is available in your key
// management service, or contact Anthropic support.
type CmekKeyUnreachableErrorDetails struct {
	ErrorCode constant.CmekKeyUnreachable `json:"error_code" default:"cmek_key_unreachable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ErrorCode   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CmekKeyUnreachableErrorDetails) RawJSON() string { return r.JSON.raw }
func (r *CmekKeyUnreachableErrorDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The organization has crossed the monthly API usage threshold of its API tier.
type EnforcedSpendLimitReachedErrorDetails struct {
	ErrorCode constant.EnforcedSpendLimitReached `json:"error_code" default:"enforced_spend_limit_reached"`
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
	Message string                `json:"message" api:"required"`
	Type    constant.TimeoutError `json:"type" default:"timeout_error"`
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
	Message string                       `json:"message" api:"required"`
	Type    constant.InvalidRequestError `json:"type" default:"invalid_request_error"`
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
	Message string                 `json:"message" api:"required"`
	Type    constant.NotFoundError `json:"type" default:"not_found_error"`
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
	ErrorCode constant.OrganizationOnHold `json:"error_code" default:"organization_on_hold"`
	AppealURL string                      `json:"appeal_url" api:"nullable"`
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

type OverloadedError struct {
	Message string                   `json:"message" api:"required"`
	Type    constant.OverloadedError `json:"type" default:"overloaded_error"`
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
	Message string                   `json:"message" api:"required"`
	Type    constant.PermissionError `json:"type" default:"permission_error"`
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

type RateLimitError struct {
	Message string                  `json:"message" api:"required"`
	Type    constant.RateLimitError `json:"type" default:"rate_limit_error"`
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
	ErrorCode constant.RequestBodyEncodingInvalid `json:"error_code" default:"request_body_encoding_invalid"`
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

// The request's `Content-Encoding` is not supported.
type RequestBodyEncodingUnsupportedErrorDetails struct {
	ErrorCode constant.RequestBodyEncodingUnsupported `json:"error_code" default:"request_body_encoding_unsupported"`
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

// The compressed request body exceeds the request size limit once decoded.
type RequestBodyTooLargeAfterDecodingErrorDetails struct {
	ErrorCode constant.RequestBodyTooLargeAfterDecoding `json:"error_code" default:"request_body_too_large_after_decoding"`
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

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/resp"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type APIErrorObject struct {
	Message string            `json:"message,required"`
	Type    constant.APIError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r APIErrorObject) RawJSON() string { return r.JSON.raw }
func (r *APIErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AuthenticationError struct {
	Message string                       `json:"message,required"`
	Type    constant.AuthenticationError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AuthenticationError) RawJSON() string { return r.JSON.raw }
func (r *AuthenticationError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BillingError struct {
	Message string                `json:"message,required"`
	Type    constant.BillingError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BillingError) RawJSON() string { return r.JSON.raw }
func (r *BillingError) UnmarshalJSON(data []byte) error {
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
		Message resp.Field
		Type    resp.Field
		raw     string
	} `json:"-"`
}

// anyErrorObject is implemented by each variant of [ErrorObjectUnion] to add type
// safety for the return type of [ErrorObjectUnion.AsAny]
type anyErrorObject interface {
	ImplErrorObjectUnion()
}

func (InvalidRequestError) ImplErrorObjectUnion() {}
func (AuthenticationError) ImplErrorObjectUnion() {}
func (BillingError) ImplErrorObjectUnion()        {}
func (PermissionError) ImplErrorObjectUnion()     {}
func (NotFoundError) ImplErrorObjectUnion()       {}
func (RateLimitError) ImplErrorObjectUnion()      {}
func (GatewayTimeoutError) ImplErrorObjectUnion() {}
func (APIErrorObject) ImplErrorObjectUnion()      {}
func (OverloadedError) ImplErrorObjectUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ErrorObjectUnion.AsAny().(type) {
//	case InvalidRequestError:
//	case AuthenticationError:
//	case BillingError:
//	case PermissionError:
//	case NotFoundError:
//	case RateLimitError:
//	case GatewayTimeoutError:
//	case APIErrorObject:
//	case OverloadedError:
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
		return u.AsGatewayTimeoutError()
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

func (u ErrorObjectUnion) AsGatewayTimeoutError() (v GatewayTimeoutError) {
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
	Error ErrorObjectUnion `json:"error,required"`
	Type  constant.Error   `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Error       resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ErrorResponse) RawJSON() string { return r.JSON.raw }
func (r *ErrorResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type GatewayTimeoutError struct {
	Message string                `json:"message,required"`
	Type    constant.TimeoutError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r GatewayTimeoutError) RawJSON() string { return r.JSON.raw }
func (r *GatewayTimeoutError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type InvalidRequestError struct {
	Message string                       `json:"message,required"`
	Type    constant.InvalidRequestError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InvalidRequestError) RawJSON() string { return r.JSON.raw }
func (r *InvalidRequestError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type NotFoundError struct {
	Message string                 `json:"message,required"`
	Type    constant.NotFoundError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r NotFoundError) RawJSON() string { return r.JSON.raw }
func (r *NotFoundError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OverloadedError struct {
	Message string                   `json:"message,required"`
	Type    constant.OverloadedError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OverloadedError) RawJSON() string { return r.JSON.raw }
func (r *OverloadedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PermissionError struct {
	Message string                   `json:"message,required"`
	Type    constant.PermissionError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r PermissionError) RawJSON() string { return r.JSON.raw }
func (r *PermissionError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RateLimitError struct {
	Message string                  `json:"message,required"`
	Type    constant.RateLimitError `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RateLimitError) RawJSON() string { return r.JSON.raw }
func (r *RateLimitError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

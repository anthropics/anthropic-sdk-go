// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/tidwall/gjson"
)

type APIErrorObject struct {
	Message string             `json:"message,required"`
	Type    APIErrorObjectType `json:"type,required"`
	JSON    apiErrorObjectJSON `json:"-"`
}

// apiErrorObjectJSON contains the JSON metadata for the struct [APIErrorObject]
type apiErrorObjectJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *APIErrorObject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r apiErrorObjectJSON) RawJSON() string {
	return r.raw
}

func (r APIErrorObject) ImplementsSharedErrorObject() {}

type APIErrorObjectType string

const (
	APIErrorObjectTypeAPIError APIErrorObjectType = "api_error"
)

func (r APIErrorObjectType) IsKnown() bool {
	switch r {
	case APIErrorObjectTypeAPIError:
		return true
	}
	return false
}

type AuthenticationError struct {
	Message string                  `json:"message,required"`
	Type    AuthenticationErrorType `json:"type,required"`
	JSON    authenticationErrorJSON `json:"-"`
}

// authenticationErrorJSON contains the JSON metadata for the struct
// [AuthenticationError]
type authenticationErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AuthenticationError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r authenticationErrorJSON) RawJSON() string {
	return r.raw
}

func (r AuthenticationError) ImplementsSharedErrorObject() {}

type AuthenticationErrorType string

const (
	AuthenticationErrorTypeAuthenticationError AuthenticationErrorType = "authentication_error"
)

func (r AuthenticationErrorType) IsKnown() bool {
	switch r {
	case AuthenticationErrorTypeAuthenticationError:
		return true
	}
	return false
}

type BillingError struct {
	Message string           `json:"message,required"`
	Type    BillingErrorType `json:"type,required"`
	JSON    billingErrorJSON `json:"-"`
}

// billingErrorJSON contains the JSON metadata for the struct [BillingError]
type billingErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BillingError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r billingErrorJSON) RawJSON() string {
	return r.raw
}

func (r BillingError) ImplementsSharedErrorObject() {}

type BillingErrorType string

const (
	BillingErrorTypeBillingError BillingErrorType = "billing_error"
)

func (r BillingErrorType) IsKnown() bool {
	switch r {
	case BillingErrorTypeBillingError:
		return true
	}
	return false
}

type ErrorObject struct {
	Message string          `json:"message,required"`
	Type    ErrorObjectType `json:"type,required"`
	JSON    errorObjectJSON `json:"-"`
	union   ErrorObjectUnion
}

// errorObjectJSON contains the JSON metadata for the struct [ErrorObject]
type errorObjectJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r errorObjectJSON) RawJSON() string {
	return r.raw
}

func (r *ErrorObject) UnmarshalJSON(data []byte) (err error) {
	*r = ErrorObject{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ErrorObjectUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [shared.InvalidRequestError],
// [shared.AuthenticationError], [shared.BillingError], [shared.PermissionError],
// [shared.NotFoundError], [shared.RateLimitError], [shared.GatewayTimeoutError],
// [shared.APIErrorObject], [shared.OverloadedError].
func (r ErrorObject) AsUnion() ErrorObjectUnion {
	return r.union
}

// Union satisfied by [shared.InvalidRequestError], [shared.AuthenticationError],
// [shared.BillingError], [shared.PermissionError], [shared.NotFoundError],
// [shared.RateLimitError], [shared.GatewayTimeoutError], [shared.APIErrorObject]
// or [shared.OverloadedError].
type ErrorObjectUnion interface {
	ImplementsSharedErrorObject()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ErrorObjectUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(InvalidRequestError{}),
			DiscriminatorValue: "invalid_request_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(AuthenticationError{}),
			DiscriminatorValue: "authentication_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BillingError{}),
			DiscriminatorValue: "billing_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(PermissionError{}),
			DiscriminatorValue: "permission_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(NotFoundError{}),
			DiscriminatorValue: "not_found_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RateLimitError{}),
			DiscriminatorValue: "rate_limit_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(GatewayTimeoutError{}),
			DiscriminatorValue: "timeout_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(APIErrorObject{}),
			DiscriminatorValue: "api_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(OverloadedError{}),
			DiscriminatorValue: "overloaded_error",
		},
	)
}

type ErrorObjectType string

const (
	ErrorObjectTypeInvalidRequestError ErrorObjectType = "invalid_request_error"
	ErrorObjectTypeAuthenticationError ErrorObjectType = "authentication_error"
	ErrorObjectTypeBillingError        ErrorObjectType = "billing_error"
	ErrorObjectTypePermissionError     ErrorObjectType = "permission_error"
	ErrorObjectTypeNotFoundError       ErrorObjectType = "not_found_error"
	ErrorObjectTypeRateLimitError      ErrorObjectType = "rate_limit_error"
	ErrorObjectTypeTimeoutError        ErrorObjectType = "timeout_error"
	ErrorObjectTypeAPIError            ErrorObjectType = "api_error"
	ErrorObjectTypeOverloadedError     ErrorObjectType = "overloaded_error"
)

func (r ErrorObjectType) IsKnown() bool {
	switch r {
	case ErrorObjectTypeInvalidRequestError, ErrorObjectTypeAuthenticationError, ErrorObjectTypeBillingError, ErrorObjectTypePermissionError, ErrorObjectTypeNotFoundError, ErrorObjectTypeRateLimitError, ErrorObjectTypeTimeoutError, ErrorObjectTypeAPIError, ErrorObjectTypeOverloadedError:
		return true
	}
	return false
}

type ErrorResponse struct {
	Error ErrorObject       `json:"error,required"`
	Type  ErrorResponseType `json:"type,required"`
	JSON  errorResponseJSON `json:"-"`
}

// errorResponseJSON contains the JSON metadata for the struct [ErrorResponse]
type errorResponseJSON struct {
	Error       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ErrorResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r errorResponseJSON) RawJSON() string {
	return r.raw
}

type ErrorResponseType string

const (
	ErrorResponseTypeError ErrorResponseType = "error"
)

func (r ErrorResponseType) IsKnown() bool {
	switch r {
	case ErrorResponseTypeError:
		return true
	}
	return false
}

type GatewayTimeoutError struct {
	Message string                  `json:"message,required"`
	Type    GatewayTimeoutErrorType `json:"type,required"`
	JSON    gatewayTimeoutErrorJSON `json:"-"`
}

// gatewayTimeoutErrorJSON contains the JSON metadata for the struct
// [GatewayTimeoutError]
type gatewayTimeoutErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *GatewayTimeoutError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r gatewayTimeoutErrorJSON) RawJSON() string {
	return r.raw
}

func (r GatewayTimeoutError) ImplementsSharedErrorObject() {}

type GatewayTimeoutErrorType string

const (
	GatewayTimeoutErrorTypeTimeoutError GatewayTimeoutErrorType = "timeout_error"
)

func (r GatewayTimeoutErrorType) IsKnown() bool {
	switch r {
	case GatewayTimeoutErrorTypeTimeoutError:
		return true
	}
	return false
}

type InvalidRequestError struct {
	Message string                  `json:"message,required"`
	Type    InvalidRequestErrorType `json:"type,required"`
	JSON    invalidRequestErrorJSON `json:"-"`
}

// invalidRequestErrorJSON contains the JSON metadata for the struct
// [InvalidRequestError]
type invalidRequestErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *InvalidRequestError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r invalidRequestErrorJSON) RawJSON() string {
	return r.raw
}

func (r InvalidRequestError) ImplementsSharedErrorObject() {}

type InvalidRequestErrorType string

const (
	InvalidRequestErrorTypeInvalidRequestError InvalidRequestErrorType = "invalid_request_error"
)

func (r InvalidRequestErrorType) IsKnown() bool {
	switch r {
	case InvalidRequestErrorTypeInvalidRequestError:
		return true
	}
	return false
}

type NotFoundError struct {
	Message string            `json:"message,required"`
	Type    NotFoundErrorType `json:"type,required"`
	JSON    notFoundErrorJSON `json:"-"`
}

// notFoundErrorJSON contains the JSON metadata for the struct [NotFoundError]
type notFoundErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *NotFoundError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r notFoundErrorJSON) RawJSON() string {
	return r.raw
}

func (r NotFoundError) ImplementsSharedErrorObject() {}

type NotFoundErrorType string

const (
	NotFoundErrorTypeNotFoundError NotFoundErrorType = "not_found_error"
)

func (r NotFoundErrorType) IsKnown() bool {
	switch r {
	case NotFoundErrorTypeNotFoundError:
		return true
	}
	return false
}

type OverloadedError struct {
	Message string              `json:"message,required"`
	Type    OverloadedErrorType `json:"type,required"`
	JSON    overloadedErrorJSON `json:"-"`
}

// overloadedErrorJSON contains the JSON metadata for the struct [OverloadedError]
type overloadedErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OverloadedError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r overloadedErrorJSON) RawJSON() string {
	return r.raw
}

func (r OverloadedError) ImplementsSharedErrorObject() {}

type OverloadedErrorType string

const (
	OverloadedErrorTypeOverloadedError OverloadedErrorType = "overloaded_error"
)

func (r OverloadedErrorType) IsKnown() bool {
	switch r {
	case OverloadedErrorTypeOverloadedError:
		return true
	}
	return false
}

type PermissionError struct {
	Message string              `json:"message,required"`
	Type    PermissionErrorType `json:"type,required"`
	JSON    permissionErrorJSON `json:"-"`
}

// permissionErrorJSON contains the JSON metadata for the struct [PermissionError]
type permissionErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *PermissionError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r permissionErrorJSON) RawJSON() string {
	return r.raw
}

func (r PermissionError) ImplementsSharedErrorObject() {}

type PermissionErrorType string

const (
	PermissionErrorTypePermissionError PermissionErrorType = "permission_error"
)

func (r PermissionErrorType) IsKnown() bool {
	switch r {
	case PermissionErrorTypePermissionError:
		return true
	}
	return false
}

type RateLimitError struct {
	Message string             `json:"message,required"`
	Type    RateLimitErrorType `json:"type,required"`
	JSON    rateLimitErrorJSON `json:"-"`
}

// rateLimitErrorJSON contains the JSON metadata for the struct [RateLimitError]
type rateLimitErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RateLimitError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r rateLimitErrorJSON) RawJSON() string {
	return r.raw
}

func (r RateLimitError) ImplementsSharedErrorObject() {}

type RateLimitErrorType string

const (
	RateLimitErrorTypeRateLimitError RateLimitErrorType = "rate_limit_error"
)

func (r RateLimitErrorType) IsKnown() bool {
	switch r {
	case RateLimitErrorTypeRateLimitError:
		return true
	}
	return false
}

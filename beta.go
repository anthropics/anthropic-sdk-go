// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
)

// BetaService contains methods and other services that help with interacting with
// the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaService] method instead.
type BetaService struct {
	Options       []option.RequestOption
	Messages      *BetaMessageService
	PromptCaching *BetaPromptCachingService
}

// NewBetaService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewBetaService(opts ...option.RequestOption) (r *BetaService) {
	r = &BetaService{}
	r.Options = opts
	r.Messages = NewBetaMessageService(opts...)
	r.PromptCaching = NewBetaPromptCachingService(opts...)
	return
}

type AnthropicBeta = string

const (
	AnthropicBetaMessageBatches2024_09_24 AnthropicBeta = "message-batches-2024-09-24"
	AnthropicBetaPromptCaching2024_07_31  AnthropicBeta = "prompt-caching-2024-07-31"
	AnthropicBetaComputerUse2024_10_22    AnthropicBeta = "computer-use-2024-10-22"
	AnthropicBetaPDFs2024_09_25           AnthropicBeta = "pdfs-2024-09-25"
)

type BetaAPIError struct {
	Message string           `json:"message,required"`
	Type    BetaAPIErrorType `json:"type,required"`
	JSON    betaAPIErrorJSON `json:"-"`
}

// betaAPIErrorJSON contains the JSON metadata for the struct [BetaAPIError]
type betaAPIErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaAPIError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaAPIErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaAPIError) implementsBetaError() {}

type BetaAPIErrorType string

const (
	BetaAPIErrorTypeAPIError BetaAPIErrorType = "api_error"
)

func (r BetaAPIErrorType) IsKnown() bool {
	switch r {
	case BetaAPIErrorTypeAPIError:
		return true
	}
	return false
}

type BetaAuthenticationError struct {
	Message string                      `json:"message,required"`
	Type    BetaAuthenticationErrorType `json:"type,required"`
	JSON    betaAuthenticationErrorJSON `json:"-"`
}

// betaAuthenticationErrorJSON contains the JSON metadata for the struct
// [BetaAuthenticationError]
type betaAuthenticationErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaAuthenticationError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaAuthenticationErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaAuthenticationError) implementsBetaError() {}

type BetaAuthenticationErrorType string

const (
	BetaAuthenticationErrorTypeAuthenticationError BetaAuthenticationErrorType = "authentication_error"
)

func (r BetaAuthenticationErrorType) IsKnown() bool {
	switch r {
	case BetaAuthenticationErrorTypeAuthenticationError:
		return true
	}
	return false
}

type BetaError struct {
	Type    BetaErrorType `json:"type,required"`
	Message string        `json:"message,required"`
	JSON    betaErrorJSON `json:"-"`
	union   BetaErrorUnion
}

// betaErrorJSON contains the JSON metadata for the struct [BetaError]
type betaErrorJSON struct {
	Type        apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r betaErrorJSON) RawJSON() string {
	return r.raw
}

func (r *BetaError) UnmarshalJSON(data []byte) (err error) {
	*r = BetaError{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaErrorUnion] interface which you can cast to the specific
// types for more type safety.
//
// Possible runtime types of the union are [BetaInvalidRequestError],
// [BetaAuthenticationError], [BetaPermissionError], [BetaNotFoundError],
// [BetaRateLimitError], [BetaAPIError], [BetaOverloadedError].
func (r BetaError) AsUnion() BetaErrorUnion {
	return r.union
}

// Union satisfied by [BetaInvalidRequestError], [BetaAuthenticationError],
// [BetaPermissionError], [BetaNotFoundError], [BetaRateLimitError], [BetaAPIError]
// or [BetaOverloadedError].
type BetaErrorUnion interface {
	implementsBetaError()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaErrorUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaInvalidRequestError{}),
			DiscriminatorValue: "invalid_request_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaAuthenticationError{}),
			DiscriminatorValue: "authentication_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaPermissionError{}),
			DiscriminatorValue: "permission_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaNotFoundError{}),
			DiscriminatorValue: "not_found_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRateLimitError{}),
			DiscriminatorValue: "rate_limit_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaAPIError{}),
			DiscriminatorValue: "api_error",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaOverloadedError{}),
			DiscriminatorValue: "overloaded_error",
		},
	)
}

type BetaErrorType string

const (
	BetaErrorTypeInvalidRequestError BetaErrorType = "invalid_request_error"
	BetaErrorTypeAuthenticationError BetaErrorType = "authentication_error"
	BetaErrorTypePermissionError     BetaErrorType = "permission_error"
	BetaErrorTypeNotFoundError       BetaErrorType = "not_found_error"
	BetaErrorTypeRateLimitError      BetaErrorType = "rate_limit_error"
	BetaErrorTypeAPIError            BetaErrorType = "api_error"
	BetaErrorTypeOverloadedError     BetaErrorType = "overloaded_error"
)

func (r BetaErrorType) IsKnown() bool {
	switch r {
	case BetaErrorTypeInvalidRequestError, BetaErrorTypeAuthenticationError, BetaErrorTypePermissionError, BetaErrorTypeNotFoundError, BetaErrorTypeRateLimitError, BetaErrorTypeAPIError, BetaErrorTypeOverloadedError:
		return true
	}
	return false
}

type BetaErrorResponse struct {
	Error BetaError             `json:"error,required"`
	Type  BetaErrorResponseType `json:"type,required"`
	JSON  betaErrorResponseJSON `json:"-"`
}

// betaErrorResponseJSON contains the JSON metadata for the struct
// [BetaErrorResponse]
type betaErrorResponseJSON struct {
	Error       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaErrorResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaErrorResponseJSON) RawJSON() string {
	return r.raw
}

type BetaErrorResponseType string

const (
	BetaErrorResponseTypeError BetaErrorResponseType = "error"
)

func (r BetaErrorResponseType) IsKnown() bool {
	switch r {
	case BetaErrorResponseTypeError:
		return true
	}
	return false
}

type BetaInvalidRequestError struct {
	Message string                      `json:"message,required"`
	Type    BetaInvalidRequestErrorType `json:"type,required"`
	JSON    betaInvalidRequestErrorJSON `json:"-"`
}

// betaInvalidRequestErrorJSON contains the JSON metadata for the struct
// [BetaInvalidRequestError]
type betaInvalidRequestErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaInvalidRequestError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaInvalidRequestErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaInvalidRequestError) implementsBetaError() {}

type BetaInvalidRequestErrorType string

const (
	BetaInvalidRequestErrorTypeInvalidRequestError BetaInvalidRequestErrorType = "invalid_request_error"
)

func (r BetaInvalidRequestErrorType) IsKnown() bool {
	switch r {
	case BetaInvalidRequestErrorTypeInvalidRequestError:
		return true
	}
	return false
}

type BetaNotFoundError struct {
	Message string                `json:"message,required"`
	Type    BetaNotFoundErrorType `json:"type,required"`
	JSON    betaNotFoundErrorJSON `json:"-"`
}

// betaNotFoundErrorJSON contains the JSON metadata for the struct
// [BetaNotFoundError]
type betaNotFoundErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaNotFoundError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaNotFoundErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaNotFoundError) implementsBetaError() {}

type BetaNotFoundErrorType string

const (
	BetaNotFoundErrorTypeNotFoundError BetaNotFoundErrorType = "not_found_error"
)

func (r BetaNotFoundErrorType) IsKnown() bool {
	switch r {
	case BetaNotFoundErrorTypeNotFoundError:
		return true
	}
	return false
}

type BetaOverloadedError struct {
	Message string                  `json:"message,required"`
	Type    BetaOverloadedErrorType `json:"type,required"`
	JSON    betaOverloadedErrorJSON `json:"-"`
}

// betaOverloadedErrorJSON contains the JSON metadata for the struct
// [BetaOverloadedError]
type betaOverloadedErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaOverloadedError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaOverloadedErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaOverloadedError) implementsBetaError() {}

type BetaOverloadedErrorType string

const (
	BetaOverloadedErrorTypeOverloadedError BetaOverloadedErrorType = "overloaded_error"
)

func (r BetaOverloadedErrorType) IsKnown() bool {
	switch r {
	case BetaOverloadedErrorTypeOverloadedError:
		return true
	}
	return false
}

type BetaPermissionError struct {
	Message string                  `json:"message,required"`
	Type    BetaPermissionErrorType `json:"type,required"`
	JSON    betaPermissionErrorJSON `json:"-"`
}

// betaPermissionErrorJSON contains the JSON metadata for the struct
// [BetaPermissionError]
type betaPermissionErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaPermissionError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaPermissionErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaPermissionError) implementsBetaError() {}

type BetaPermissionErrorType string

const (
	BetaPermissionErrorTypePermissionError BetaPermissionErrorType = "permission_error"
)

func (r BetaPermissionErrorType) IsKnown() bool {
	switch r {
	case BetaPermissionErrorTypePermissionError:
		return true
	}
	return false
}

type BetaRateLimitError struct {
	Message string                 `json:"message,required"`
	Type    BetaRateLimitErrorType `json:"type,required"`
	JSON    betaRateLimitErrorJSON `json:"-"`
}

// betaRateLimitErrorJSON contains the JSON metadata for the struct
// [BetaRateLimitError]
type betaRateLimitErrorJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRateLimitError) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRateLimitErrorJSON) RawJSON() string {
	return r.raw
}

func (r BetaRateLimitError) implementsBetaError() {}

type BetaRateLimitErrorType string

const (
	BetaRateLimitErrorTypeRateLimitError BetaRateLimitErrorType = "rate_limit_error"
)

func (r BetaRateLimitErrorType) IsKnown() bool {
	switch r {
	case BetaRateLimitErrorTypeRateLimitError:
		return true
	}
	return false
}

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/internal/apierror"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/shared"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type Error = apierror.Error

// This is an alias to an internal type.
type APIErrorObject = shared.APIErrorObject

// This is an alias to an internal type.
type ErrorObjectUnion = shared.ErrorObjectUnion

// This is an alias to an internal type.
type ErrorResponse = shared.ErrorResponse

// This is an alias to an internal type.
type GatewayTimeoutError = shared.GatewayTimeoutError

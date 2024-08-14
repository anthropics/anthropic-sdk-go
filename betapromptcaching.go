// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BetaPromptCachingService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaPromptCachingService] method instead.
type BetaPromptCachingService struct {
	Options  []option.RequestOption
	Messages *BetaPromptCachingMessageService
}

// NewBetaPromptCachingService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaPromptCachingService(opts ...option.RequestOption) (r *BetaPromptCachingService) {
	r = &BetaPromptCachingService{}
	r.Options = opts
	r.Messages = NewBetaPromptCachingMessageService(opts...)
	return
}

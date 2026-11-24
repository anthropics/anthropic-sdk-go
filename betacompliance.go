package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaComplianceService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaComplianceService] method instead.
type BetaComplianceService struct {
	Options []option.RequestOption
	Review  BetaComplianceReviewService
}

// NewBetaComplianceService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaComplianceService(opts ...option.RequestOption) (r BetaComplianceService) {
	r = BetaComplianceService{}
	r.Options = opts
	r.Review = NewBetaComplianceReviewService(opts...)
	return
}

package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaComplianceReviewService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaComplianceReviewService] method instead.
type BetaComplianceReviewService struct {
	Options []option.RequestOption
	Cases   BetaComplianceReviewCaseService
	Notices BetaComplianceReviewNoticeService
	Items   BetaComplianceReviewItemService
}

// NewBetaComplianceReviewService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaComplianceReviewService(opts ...option.RequestOption) (r BetaComplianceReviewService) {
	r = BetaComplianceReviewService{}
	r.Options = opts
	r.Cases = NewBetaComplianceReviewCaseService(opts...)
	r.Notices = NewBetaComplianceReviewNoticeService(opts...)
	r.Items = NewBetaComplianceReviewItemService(opts...)
	return
}

// A sort direction: `asc` or `desc`.
type BetaComplianceReviewOrder string

const (
	// Oldest first.
	BetaComplianceReviewOrderAsc BetaComplianceReviewOrder = "asc"
	// Newest first.
	BetaComplianceReviewOrderDesc BetaComplianceReviewOrder = "desc"
)

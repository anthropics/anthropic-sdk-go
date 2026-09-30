package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BetaOrganizationAnalyticsAppChatService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsAppChatService] method instead.
type BetaOrganizationAnalyticsAppChatService struct {
	Options  []option.RequestOption
	Projects BetaOrganizationAnalyticsAppChatProjectService
}

// NewBetaOrganizationAnalyticsAppChatService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsAppChatService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsAppChatService) {
	r = BetaOrganizationAnalyticsAppChatService{}
	r.Options = opts
	r.Projects = NewBetaOrganizationAnalyticsAppChatProjectService(opts...)
	return
}

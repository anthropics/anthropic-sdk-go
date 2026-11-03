package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// OrganizationAnalyticsAppChatService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationAnalyticsAppChatService] method instead.
type OrganizationAnalyticsAppChatService struct {
	Options  []option.RequestOption
	Projects OrganizationAnalyticsAppChatProjectService
}

// NewOrganizationAnalyticsAppChatService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationAnalyticsAppChatService(opts ...option.RequestOption) (r OrganizationAnalyticsAppChatService) {
	r = OrganizationAnalyticsAppChatService{}
	r.Options = opts
	r.Projects = NewOrganizationAnalyticsAppChatProjectService(opts...)
	return
}

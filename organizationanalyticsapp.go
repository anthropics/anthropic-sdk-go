package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// OrganizationAnalyticsAppService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationAnalyticsAppService] method instead.
type OrganizationAnalyticsAppService struct {
	Options []option.RequestOption
	Chat    OrganizationAnalyticsAppChatService
}

// NewOrganizationAnalyticsAppService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationAnalyticsAppService(opts ...option.RequestOption) (r OrganizationAnalyticsAppService) {
	r = OrganizationAnalyticsAppService{}
	r.Options = opts
	r.Chat = NewOrganizationAnalyticsAppChatService(opts...)
	return
}

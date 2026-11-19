package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaOrganizationClaudeCodeManagedSettingService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeCodeManagedSettingService] method instead.
type BetaOrganizationClaudeCodeManagedSettingService struct {
	Options          []option.RequestOption
	Configuration    BetaOrganizationClaudeCodeManagedSettingConfigurationService
	RBACGroups       BetaOrganizationClaudeCodeManagedSettingRBACGroupService
	RBACGroupDefault BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService
}

// NewBetaOrganizationClaudeCodeManagedSettingService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationClaudeCodeManagedSettingService(opts ...option.RequestOption) (r BetaOrganizationClaudeCodeManagedSettingService) {
	r = BetaOrganizationClaudeCodeManagedSettingService{}
	r.Options = opts
	r.Configuration = NewBetaOrganizationClaudeCodeManagedSettingConfigurationService(opts...)
	r.RBACGroups = NewBetaOrganizationClaudeCodeManagedSettingRBACGroupService(opts...)
	r.RBACGroupDefault = NewBetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService(opts...)
	return
}

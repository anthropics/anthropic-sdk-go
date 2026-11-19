package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// BetaOrganizationClaudeCodeService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeCodeService] method instead.
type BetaOrganizationClaudeCodeService struct {
	Options         []option.RequestOption
	ManagedSettings BetaOrganizationClaudeCodeManagedSettingService
}

// NewBetaOrganizationClaudeCodeService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationClaudeCodeService(opts ...option.RequestOption) (r BetaOrganizationClaudeCodeService) {
	r = BetaOrganizationClaudeCodeService{}
	r.Options = opts
	r.ManagedSettings = NewBetaOrganizationClaudeCodeManagedSettingService(opts...)
	return
}

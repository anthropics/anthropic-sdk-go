package anthropic

import (
	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaService contains methods and other services that help with interacting with
// the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaService] method instead.
type BetaService struct {
	Options        []option.RequestOption
	Models         BetaModelService
	Messages       BetaMessageService
	Agents         BetaAgentService
	Environments   BetaEnvironmentService
	Sessions       BetaSessionService
	Deployments    BetaDeploymentService
	DeploymentRuns BetaDeploymentRunService
	Vaults         BetaVaultService
	MemoryStores   BetaMemoryStoreService
	Files          BetaFileService
	Skills         BetaSkillService
	Webhooks       BetaWebhookService
	// Register the end customers that your platform serves, so each request can name
	// who it's for. To get started, see the
	// [User profiles guide](https://platform.claude.com/docs/en/manage-claude/user-profiles).
	UserProfiles BetaUserProfileService
	Dreams       BetaDreamService
	Tunnels      BetaTunnelService
	Compliance   BetaComplianceService
	Organization BetaOrganizationService
}

// NewBetaService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewBetaService(opts ...option.RequestOption) (r BetaService) {
	r = BetaService{}
	r.Options = opts
	r.Models = NewBetaModelService(opts...)
	r.Messages = NewBetaMessageService(opts...)
	r.Agents = NewBetaAgentService(opts...)
	r.Environments = NewBetaEnvironmentService(opts...)
	r.Sessions = NewBetaSessionService(opts...)
	r.Deployments = NewBetaDeploymentService(opts...)
	r.DeploymentRuns = NewBetaDeploymentRunService(opts...)
	r.Vaults = NewBetaVaultService(opts...)
	r.MemoryStores = NewBetaMemoryStoreService(opts...)
	r.Files = NewBetaFileService(opts...)
	r.Skills = NewBetaSkillService(opts...)
	r.Webhooks = NewBetaWebhookService(opts...)
	r.UserProfiles = NewBetaUserProfileService(opts...)
	r.Dreams = NewBetaDreamService(opts...)
	r.Tunnels = NewBetaTunnelService(opts...)
	r.Compliance = NewBetaComplianceService(opts...)
	r.Organization = NewBetaOrganizationService(opts...)
	return
}

// Beta names you can send in the `anthropic-beta` header.
//
// See [Beta headers](https://platform.claude.com/docs/en/api/beta-headers) for how
// to send them and where to find the name a feature needs.
type AnthropicBeta = string

const (
	AnthropicBetaComputerUse2024_10_22                  AnthropicBeta = "computer-use-2024-10-22"
	AnthropicBetaComputerUse2025_01_24                  AnthropicBeta = "computer-use-2025-01-24"
	AnthropicBetaOutput128k2025_02_19                   AnthropicBeta = "output-128k-2025-02-19"
	AnthropicBetaMCPClient2025_04_04                    AnthropicBeta = "mcp-client-2025-04-04"
	AnthropicBetaMCPClient2025_11_20                    AnthropicBeta = "mcp-client-2025-11-20"
	AnthropicBetaDevFullThinking2025_05_14              AnthropicBeta = "dev-full-thinking-2025-05-14"
	AnthropicBetaInterleavedThinking2025_05_14          AnthropicBeta = "interleaved-thinking-2025-05-14"
	AnthropicBetaCodeExecution2025_05_22                AnthropicBeta = "code-execution-2025-05-22"
	AnthropicBetaContext1m2025_08_07                    AnthropicBeta = "context-1m-2025-08-07"
	AnthropicBetaContextManagement2025_06_27            AnthropicBeta = "context-management-2025-06-27"
	AnthropicBetaModelContextWindowExceeded2025_08_26   AnthropicBeta = "model-context-window-exceeded-2025-08-26"
	AnthropicBetaFastMode2026_02_01                     AnthropicBeta = "fast-mode-2026-02-01"
	AnthropicBetaOutput300k2026_03_24                   AnthropicBeta = "output-300k-2026-03-24"
	AnthropicBetaUserProfiles2026_03_24                 AnthropicBeta = "user-profiles-2026-03-24"
	AnthropicBetaUserProfiles2026_08_18                 AnthropicBeta = "user-profiles-2026-08-18"
	AnthropicBetaUserProfiles2026_09_04                 AnthropicBeta = "user-profiles-2026-09-04"
	AnthropicBetaAdvisorTool2026_03_01                  AnthropicBeta = "advisor-tool-2026-03-01"
	AnthropicBetaManagedAgents2026_04_01                AnthropicBeta = "managed-agents-2026-04-01"
	AnthropicBetaDreaming2026_04_21                     AnthropicBeta = "dreaming-2026-04-21"
	AnthropicBetaThinkingTokenCount2026_05_13           AnthropicBeta = "thinking-token-count-2026-05-13"
	AnthropicBetaServerSideFallback2026_06_01           AnthropicBeta = "server-side-fallback-2026-06-01"
	AnthropicBetaServerSideFallback2026_07_01           AnthropicBeta = "server-side-fallback-2026-07-01"
	AnthropicBetaFallbackCredit2026_06_01               AnthropicBeta = "fallback-credit-2026-06-01"
	AnthropicBetaFallbackCredit2026_07_01               AnthropicBeta = "fallback-credit-2026-07-01"
	AnthropicBetaAgentMemory2026_07_22                  AnthropicBeta = "agent-memory-2026-07-22"
	AnthropicBetaMidConversationToolChanges2026_07_01   AnthropicBeta = "mid-conversation-tool-changes-2026-07-01"
	AnthropicBetaCompact2026_01_12                      AnthropicBeta = "compact-2026-01-12"
	AnthropicBetaComputerUse2025_11_24                  AnthropicBeta = "computer-use-2025-11-24"
	AnthropicBetaMCPTunnels2026_06_22                   AnthropicBeta = "mcp-tunnels-2026-06-22"
	AnthropicBetaStructuredOutputs2025_11_13            AnthropicBeta = "structured-outputs-2025-11-13"
	AnthropicBetaTaskBudgets2026_03_13                  AnthropicBeta = "task-budgets-2026-03-13"
	AnthropicBetaThinkingDisplayUpdates2026_08_18       AnthropicBeta = "thinking-display-updates-2026-08-18"
	AnthropicBetaMidConversationOutputConfig2026_07_01  AnthropicBeta = "mid-conversation-output-config-2026-07-01"
	AnthropicBetaThinkingBindingControls2026_08_01      AnthropicBeta = "thinking-binding-controls-2026-08-01"
	AnthropicBetaMidConversationSystemClearAt2026_08_21 AnthropicBeta = "mid-conversation-system-clear-at-2026-08-21"
	AnthropicBetaCompact2026_09_04                      AnthropicBeta = "compact-2026-09-04"
	AnthropicBetaInlineTools2026_09_15                  AnthropicBeta = "inline-tools-2026-09-15"
	AnthropicBetaMCPClient2026_09_15                    AnthropicBeta = "mcp-client-2026-09-15"
	AnthropicBetaCEPlugins2026_09_01                    AnthropicBeta = "ce-plugins-2026-09-01"
	AnthropicBetaSpendLimitReads2026_09_26              AnthropicBeta = "spend-limit-reads-2026-09-26"
	AnthropicBetaTelemetryDestinations2026_08_11        AnthropicBeta = "telemetry-destinations-2026-08-11"
	AnthropicBetaWebhooks2026_10_07                     AnthropicBeta = "webhooks-2026-10-07"
	AnthropicBetaAdminUserSignOut2026_10_08             AnthropicBeta = "admin-user-sign-out-2026-10-08"
	AnthropicBetaMessageThreads2026_08_12               AnthropicBeta = "message-threads-2026-08-12"
	AnthropicBetaAnalyticsUsageCost2026_08_24           AnthropicBeta = "analytics-usage-cost-2026-08-24"
	AnthropicBetaExperimentalClaudeTag2026_10_21        AnthropicBeta = "experimental-claude-tag-2026-10-21"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaMessageBatches2024_09_24 AnthropicBeta = "message-batches-2024-09-24"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaPromptCaching2024_07_31 AnthropicBeta = "prompt-caching-2024-07-31"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaPDFs2024_09_25 AnthropicBeta = "pdfs-2024-09-25"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaTokenCounting2024_11_01 AnthropicBeta = "token-counting-2024-11-01"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaTokenEfficientTools2025_02_19 AnthropicBeta = "token-efficient-tools-2025-02-19"
	// Deprecated: Sending this value selects the beta response shapes. Removing it
	// moves you to the current shapes, so follow the migration steps:
	// https://platform.claude.com/docs/en/build-with-claude/files#migrate-from-files-api-2025-04-14
	AnthropicBetaFilesAPI2025_04_14 AnthropicBeta = "files-api-2025-04-14"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaExtendedCacheTTL2025_04_11 AnthropicBeta = "extended-cache-ttl-2025-04-11"
	// Deprecated: Sending this value selects the beta response shapes. Removing it
	// moves you to the current shapes, so follow the migration steps:
	// https://platform.claude.com/docs/en/build-with-claude/skills-guide#migrate-from-skills-2025-10-02
	AnthropicBetaSkills2025_10_02 AnthropicBeta = "skills-2025-10-02"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaCacheDiagnosis2026_04_07 AnthropicBeta = "cache-diagnosis-2026-04-07"
	// Deprecated: This value is no longer needed. Remove it from your requests.
	AnthropicBetaCEUserManagement2026_07_13            AnthropicBeta = "ce-user-management-2026-07-13"
	AnthropicBetaCEClaudeCodeManagedSettings2026_10_06 AnthropicBeta = "ce-claude-code-managed-settings-2026-10-06"
	AnthropicBetaEfs2026_10_01                         AnthropicBeta = "efs-2026-10-01"
	AnthropicBetaGoogleManagedDataControls2026_09_01   AnthropicBeta = "google-managed-data-controls-2026-09-01"
	AnthropicBetaUserProfilesBilling2026_10_01         AnthropicBeta = "user-profiles-billing-2026-10-01"
)

type BetaAPIErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.APIError `json:"type" default:"api_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAPIErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaAPIErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAuthenticationErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.AuthenticationError `json:"type" default:"authentication_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAuthenticationErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaAuthenticationErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaBillingErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.BillingError `json:"type" default:"billing_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaBillingErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaBillingErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaCurrency string

const (
	BetaCurrencyUSD BetaCurrency = "USD"
)

// BetaErrorObjectUnion contains all possible properties and values from
// [BetaInvalidRequestErrorObject], [BetaAuthenticationErrorObject],
// [BetaBillingErrorObject], [BetaPermissionErrorObject],
// [BetaNotFoundErrorObject], [BetaRateLimitErrorObject],
// [BetaGatewayTimeoutErrorObject], [BetaAPIErrorObject],
// [BetaOverloadedErrorObject].
//
// Use the [BetaErrorObjectUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaErrorObjectUnion struct {
	Message string `json:"message"`
	// Any of "invalid_request_error", "authentication_error", "billing_error",
	// "permission_error", "not_found_error", "rate_limit_error", "timeout_error",
	// "api_error", "overloaded_error".
	Type string `json:"type"`
	JSON struct {
		Message respjson.Field
		Type    respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaErrorObject is implemented by each variant of [BetaErrorObjectUnion] to
// add type safety for the return type of [BetaErrorObjectUnion.AsAny]
type anyBetaErrorObject interface {
	implBetaErrorObjectUnion()
}

func (BetaInvalidRequestErrorObject) implBetaErrorObjectUnion() {}
func (BetaAuthenticationErrorObject) implBetaErrorObjectUnion() {}
func (BetaBillingErrorObject) implBetaErrorObjectUnion()        {}
func (BetaPermissionErrorObject) implBetaErrorObjectUnion()     {}
func (BetaNotFoundErrorObject) implBetaErrorObjectUnion()       {}
func (BetaRateLimitErrorObject) implBetaErrorObjectUnion()      {}
func (BetaGatewayTimeoutErrorObject) implBetaErrorObjectUnion() {}
func (BetaAPIErrorObject) implBetaErrorObjectUnion()            {}
func (BetaOverloadedErrorObject) implBetaErrorObjectUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaErrorObjectUnion.AsAny().(type) {
//	case anthropic.BetaInvalidRequestErrorObject:
//	case anthropic.BetaAuthenticationErrorObject:
//	case anthropic.BetaBillingErrorObject:
//	case anthropic.BetaPermissionErrorObject:
//	case anthropic.BetaNotFoundErrorObject:
//	case anthropic.BetaRateLimitErrorObject:
//	case anthropic.BetaGatewayTimeoutErrorObject:
//	case anthropic.BetaAPIErrorObject:
//	case anthropic.BetaOverloadedErrorObject:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaErrorObjectUnion) AsAny() anyBetaErrorObject {
	switch u.Type {
	case "invalid_request_error":
		return u.AsInvalidRequestError()
	case "authentication_error":
		return u.AsAuthenticationError()
	case "billing_error":
		return u.AsBillingError()
	case "permission_error":
		return u.AsPermissionError()
	case "not_found_error":
		return u.AsNotFoundError()
	case "rate_limit_error":
		return u.AsRateLimitError()
	case "timeout_error":
		return u.AsTimeoutError()
	case "api_error":
		return u.AsAPIError()
	case "overloaded_error":
		return u.AsOverloadedError()
	}
	return nil
}

func (u BetaErrorObjectUnion) AsInvalidRequestError() (v BetaInvalidRequestErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsAuthenticationError() (v BetaAuthenticationErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsBillingError() (v BetaBillingErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsPermissionError() (v BetaPermissionErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsNotFoundError() (v BetaNotFoundErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsRateLimitError() (v BetaRateLimitErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsTimeoutError() (v BetaGatewayTimeoutErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsAPIError() (v BetaAPIErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaErrorObjectUnion) AsOverloadedError() (v BetaOverloadedErrorObject) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaErrorObjectUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaErrorObjectUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaErrorResponse struct {
	Error     BetaErrorObjectUnion `json:"error" api:"required"`
	RequestID string               `json:"request_id" api:"required"`
	Type      constant.Error       `json:"type" default:"error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		RequestID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaErrorResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaErrorResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaGatewayTimeoutErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.TimeoutError `json:"type" default:"timeout_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGatewayTimeoutErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaGatewayTimeoutErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaInvalidRequestErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.InvalidRequestError `json:"type" default:"invalid_request_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInvalidRequestErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaInvalidRequestErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A monetary amount in a specific currency.
type BetaMonetaryAmount struct {
	// Amount in minor units of the currency, as an integer decimal string with no
	// leading zeros: "2500" is $25.00 and "50" is fifty cents. A string rather than a
	// number so no float rounding is ever applied.
	Amount string `json:"amount" api:"required"`
	// Uppercase ISO-4217 currency code. `USD` is the only currency currently
	// supported; the accepted set is closed and grows only when a new currency is
	// priced.
	//
	// Any of "USD".
	Currency BetaCurrency `json:"currency" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount      respjson.Field
		Currency    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaMonetaryAmount) RawJSON() string { return r.JSON.raw }
func (r *BetaMonetaryAmount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaMonetaryAmount to a BetaMonetaryAmountParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaMonetaryAmountParam.Overrides()
func (r BetaMonetaryAmount) ToParam() BetaMonetaryAmountParam {
	return param.Override[BetaMonetaryAmountParam](json.RawMessage(r.RawJSON()))
}

// A monetary amount in a specific currency.
//
// The properties Amount, Currency are required.
type BetaMonetaryAmountParam struct {
	// Amount in minor units of the currency, as an integer decimal string with no
	// leading zeros: "2500" is $25.00 and "50" is fifty cents. A string rather than a
	// number so no float rounding is ever applied.
	Amount string `json:"amount" api:"required"`
	// Uppercase ISO-4217 currency code. `USD` is the only currency currently
	// supported; the accepted set is closed and grows only when a new currency is
	// priced.
	//
	// Any of "USD".
	Currency BetaCurrency `json:"currency,omitzero" api:"required"`
	paramObj
}

func (r BetaMonetaryAmountParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaMonetaryAmountParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaMonetaryAmountParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaNotFoundErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.NotFoundError `json:"type" default:"not_found_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaNotFoundErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaNotFoundErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOverloadedErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.OverloadedError `json:"type" default:"overloaded_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOverloadedErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaOverloadedErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaPermissionErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.PermissionError `json:"type" default:"permission_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPermissionErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaPermissionErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRateLimitErrorObject struct {
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	Message string `json:"message" api:"required"`
	// The broad category of the error.
	Type constant.RateLimitError `json:"type" default:"rate_limit_error"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRateLimitErrorObject) RawJSON() string { return r.JSON.raw }
func (r *BetaRateLimitErrorObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deprecated: BetaErrorUnion has been renamed to BetaErrorObjectUnion
type BetaErrorUnion = BetaErrorObjectUnion

// Deprecated: BetaAPIError has been renamed to BetaAPIErrorObject
type BetaAPIError = BetaAPIErrorObject

// Deprecated: BetaAuthenticationError has been renamed to
// BetaAuthenticationErrorObject
type BetaAuthenticationError = BetaAuthenticationErrorObject

// Deprecated: BetaInvalidRequestError has been renamed to
// BetaInvalidRequestErrorObject
type BetaInvalidRequestError = BetaInvalidRequestErrorObject

// Deprecated: BetaBillingError has been renamed to BetaBillingErrorObject
type BetaBillingError = BetaBillingErrorObject

// Deprecated: BetaGatewayTimeoutError has been renamed to
// BetaGatewayTimeoutErrorObject
type BetaGatewayTimeoutError = BetaGatewayTimeoutErrorObject

// Deprecated: BetaNotFoundError has been renamed to BetaNotFoundErrorObject
type BetaNotFoundError = BetaNotFoundErrorObject

// Deprecated: BetaOverloadedError has been renamed to BetaOverloadedErrorObject
type BetaOverloadedError = BetaOverloadedErrorObject

// Deprecated: BetaPermissionError has been renamed to BetaPermissionErrorObject
type BetaPermissionError = BetaPermissionErrorObject

// Deprecated: BetaRateLimitError has been renamed to BetaRateLimitErrorObject
type BetaRateLimitError = BetaRateLimitErrorObject

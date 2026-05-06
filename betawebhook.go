// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// BetaWebhookService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaWebhookService] method instead.
type BetaWebhookService struct {
	Options []option.RequestOption
}

// NewBetaWebhookService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaWebhookService(opts ...option.RequestOption) (r BetaWebhookService) {
	r = BetaWebhookService{}
	r.Options = opts
	return
}

func (r *BetaWebhookService) Unwrap(payload []byte, headers http.Header, opts ...option.RequestOption) (*UnwrapWebhookEvent, error) {
	opts = slices.Concat(r.Options, opts)
	cfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return nil, err
	}
	key := cfg.WebhookKey
	if key == "" {
		return nil, errors.New("The WebhookKey option must be set in order to verify webhook headers")
	}
	wh, err := standardwebhooks.NewWebhook(key)
	if err != nil {
		return nil, err
	}
	err = wh.Verify(payload, headers)
	if err != nil {
		return nil, err
	}
	res := &UnwrapWebhookEvent{}
	err = res.UnmarshalJSON(payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

// BetaWebhookEventDataUnion contains all possible properties and values from
// [BetaWebhookSessionCreatedEventData], [BetaWebhookSessionPendingEventData],
// [BetaWebhookSessionRunningEventData], [BetaWebhookSessionIdledEventData],
// [BetaWebhookSessionRequiresActionEventData],
// [BetaWebhookSessionArchivedEventData], [BetaWebhookSessionDeletedEventData],
// [BetaWebhookSessionStatusRescheduledEventData],
// [BetaWebhookSessionStatusRunStartedEventData],
// [BetaWebhookSessionStatusIdledEventData],
// [BetaWebhookSessionStatusTerminatedEventData],
// [BetaWebhookSessionThreadCreatedEventData],
// [BetaWebhookSessionThreadIdledEventData],
// [BetaWebhookSessionThreadTerminatedEventData],
// [BetaWebhookSessionOutcomeEvaluationEndedEventData],
// [BetaWebhookVaultCreatedEventData], [BetaWebhookVaultArchivedEventData],
// [BetaWebhookVaultDeletedEventData],
// [BetaWebhookVaultCredentialCreatedEventData],
// [BetaWebhookVaultCredentialArchivedEventData],
// [BetaWebhookVaultCredentialDeletedEventData],
// [BetaWebhookVaultCredentialRefreshFailedEventData].
//
// Use the [BetaWebhookEventDataUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaWebhookEventDataUnion struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	// Any of "session.created", "session.pending", "session.running", "session.idled",
	// "session.requires_action", "session.archived", "session.deleted",
	// "session.status_rescheduled", "session.status_run_started",
	// "session.status_idled", "session.status_terminated", "session.thread_created",
	// "session.thread_idled", "session.thread_terminated",
	// "session.outcome_evaluation_ended", "vault.created", "vault.archived",
	// "vault.deleted", "vault_credential.created", "vault_credential.archived",
	// "vault_credential.deleted", "vault_credential.refresh_failed".
	Type        string `json:"type"`
	WorkspaceID string `json:"workspace_id"`
	VaultID     string `json:"vault_id"`
	JSON        struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		VaultID        respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaWebhookEventData is implemented by each variant of
// [BetaWebhookEventDataUnion] to add type safety for the return type of
// [BetaWebhookEventDataUnion.AsAny]
type anyBetaWebhookEventData interface {
	implBetaWebhookEventDataUnion()
}

func (BetaWebhookSessionCreatedEventData) implBetaWebhookEventDataUnion()                {}
func (BetaWebhookSessionPendingEventData) implBetaWebhookEventDataUnion()                {}
func (BetaWebhookSessionRunningEventData) implBetaWebhookEventDataUnion()                {}
func (BetaWebhookSessionIdledEventData) implBetaWebhookEventDataUnion()                  {}
func (BetaWebhookSessionRequiresActionEventData) implBetaWebhookEventDataUnion()         {}
func (BetaWebhookSessionArchivedEventData) implBetaWebhookEventDataUnion()               {}
func (BetaWebhookSessionDeletedEventData) implBetaWebhookEventDataUnion()                {}
func (BetaWebhookSessionStatusRescheduledEventData) implBetaWebhookEventDataUnion()      {}
func (BetaWebhookSessionStatusRunStartedEventData) implBetaWebhookEventDataUnion()       {}
func (BetaWebhookSessionStatusIdledEventData) implBetaWebhookEventDataUnion()            {}
func (BetaWebhookSessionStatusTerminatedEventData) implBetaWebhookEventDataUnion()       {}
func (BetaWebhookSessionThreadCreatedEventData) implBetaWebhookEventDataUnion()          {}
func (BetaWebhookSessionThreadIdledEventData) implBetaWebhookEventDataUnion()            {}
func (BetaWebhookSessionThreadTerminatedEventData) implBetaWebhookEventDataUnion()       {}
func (BetaWebhookSessionOutcomeEvaluationEndedEventData) implBetaWebhookEventDataUnion() {}
func (BetaWebhookVaultCreatedEventData) implBetaWebhookEventDataUnion()                  {}
func (BetaWebhookVaultArchivedEventData) implBetaWebhookEventDataUnion()                 {}
func (BetaWebhookVaultDeletedEventData) implBetaWebhookEventDataUnion()                  {}
func (BetaWebhookVaultCredentialCreatedEventData) implBetaWebhookEventDataUnion()        {}
func (BetaWebhookVaultCredentialArchivedEventData) implBetaWebhookEventDataUnion()       {}
func (BetaWebhookVaultCredentialDeletedEventData) implBetaWebhookEventDataUnion()        {}
func (BetaWebhookVaultCredentialRefreshFailedEventData) implBetaWebhookEventDataUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaWebhookEventDataUnion.AsAny().(type) {
//	case anthropic.BetaWebhookSessionCreatedEventData:
//	case anthropic.BetaWebhookSessionPendingEventData:
//	case anthropic.BetaWebhookSessionRunningEventData:
//	case anthropic.BetaWebhookSessionIdledEventData:
//	case anthropic.BetaWebhookSessionRequiresActionEventData:
//	case anthropic.BetaWebhookSessionArchivedEventData:
//	case anthropic.BetaWebhookSessionDeletedEventData:
//	case anthropic.BetaWebhookSessionStatusRescheduledEventData:
//	case anthropic.BetaWebhookSessionStatusRunStartedEventData:
//	case anthropic.BetaWebhookSessionStatusIdledEventData:
//	case anthropic.BetaWebhookSessionStatusTerminatedEventData:
//	case anthropic.BetaWebhookSessionThreadCreatedEventData:
//	case anthropic.BetaWebhookSessionThreadIdledEventData:
//	case anthropic.BetaWebhookSessionThreadTerminatedEventData:
//	case anthropic.BetaWebhookSessionOutcomeEvaluationEndedEventData:
//	case anthropic.BetaWebhookVaultCreatedEventData:
//	case anthropic.BetaWebhookVaultArchivedEventData:
//	case anthropic.BetaWebhookVaultDeletedEventData:
//	case anthropic.BetaWebhookVaultCredentialCreatedEventData:
//	case anthropic.BetaWebhookVaultCredentialArchivedEventData:
//	case anthropic.BetaWebhookVaultCredentialDeletedEventData:
//	case anthropic.BetaWebhookVaultCredentialRefreshFailedEventData:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaWebhookEventDataUnion) AsAny() anyBetaWebhookEventData {
	switch u.Type {
	case "session.created":
		return u.AsSessionCreated()
	case "session.pending":
		return u.AsSessionPending()
	case "session.running":
		return u.AsSessionRunning()
	case "session.idled":
		return u.AsSessionIdled()
	case "session.requires_action":
		return u.AsSessionRequiresAction()
	case "session.archived":
		return u.AsSessionArchived()
	case "session.deleted":
		return u.AsSessionDeleted()
	case "session.status_rescheduled":
		return u.AsSessionStatusRescheduled()
	case "session.status_run_started":
		return u.AsSessionStatusRunStarted()
	case "session.status_idled":
		return u.AsSessionStatusIdled()
	case "session.status_terminated":
		return u.AsSessionStatusTerminated()
	case "session.thread_created":
		return u.AsSessionThreadCreated()
	case "session.thread_idled":
		return u.AsSessionThreadIdled()
	case "session.thread_terminated":
		return u.AsSessionThreadTerminated()
	case "session.outcome_evaluation_ended":
		return u.AsSessionOutcomeEvaluationEnded()
	case "vault.created":
		return u.AsVaultCreated()
	case "vault.archived":
		return u.AsVaultArchived()
	case "vault.deleted":
		return u.AsVaultDeleted()
	case "vault_credential.created":
		return u.AsVaultCredentialCreated()
	case "vault_credential.archived":
		return u.AsVaultCredentialArchived()
	case "vault_credential.deleted":
		return u.AsVaultCredentialDeleted()
	case "vault_credential.refresh_failed":
		return u.AsVaultCredentialRefreshFailed()
	}
	return nil
}

func (u BetaWebhookEventDataUnion) AsSessionCreated() (v BetaWebhookSessionCreatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionPending() (v BetaWebhookSessionPendingEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionRunning() (v BetaWebhookSessionRunningEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionIdled() (v BetaWebhookSessionIdledEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionRequiresAction() (v BetaWebhookSessionRequiresActionEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionArchived() (v BetaWebhookSessionArchivedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionDeleted() (v BetaWebhookSessionDeletedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionStatusRescheduled() (v BetaWebhookSessionStatusRescheduledEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionStatusRunStarted() (v BetaWebhookSessionStatusRunStartedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionStatusIdled() (v BetaWebhookSessionStatusIdledEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionStatusTerminated() (v BetaWebhookSessionStatusTerminatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionThreadCreated() (v BetaWebhookSessionThreadCreatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionThreadIdled() (v BetaWebhookSessionThreadIdledEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionThreadTerminated() (v BetaWebhookSessionThreadTerminatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsSessionOutcomeEvaluationEnded() (v BetaWebhookSessionOutcomeEvaluationEndedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultCreated() (v BetaWebhookVaultCreatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultArchived() (v BetaWebhookVaultArchivedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultDeleted() (v BetaWebhookVaultDeletedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultCredentialCreated() (v BetaWebhookVaultCredentialCreatedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultCredentialArchived() (v BetaWebhookVaultCredentialArchivedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultCredentialDeleted() (v BetaWebhookVaultCredentialDeletedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaWebhookEventDataUnion) AsVaultCredentialRefreshFailed() (v BetaWebhookVaultCredentialRefreshFailedEventData) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaWebhookEventDataUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaWebhookEventDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionArchivedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                   `json:"id" api:"required"`
	OrganizationID string                   `json:"organization_id" api:"required"`
	Type           constant.SessionArchived `json:"type" default:"session.archived"`
	WorkspaceID    string                   `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionArchivedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionArchivedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionCreatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                  `json:"id" api:"required"`
	OrganizationID string                  `json:"organization_id" api:"required"`
	Type           constant.SessionCreated `json:"type" default:"session.created"`
	WorkspaceID    string                  `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionCreatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionCreatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionDeletedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                  `json:"id" api:"required"`
	OrganizationID string                  `json:"organization_id" api:"required"`
	Type           constant.SessionDeleted `json:"type" default:"session.deleted"`
	WorkspaceID    string                  `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionDeletedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionDeletedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionIdledEventData struct {
	// ID of the resource that triggered the event.
	ID             string                `json:"id" api:"required"`
	OrganizationID string                `json:"organization_id" api:"required"`
	Type           constant.SessionIdled `json:"type" default:"session.idled"`
	WorkspaceID    string                `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionIdledEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionIdledEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionOutcomeEvaluationEndedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                                 `json:"id" api:"required"`
	OrganizationID string                                 `json:"organization_id" api:"required"`
	Type           constant.SessionOutcomeEvaluationEnded `json:"type" default:"session.outcome_evaluation_ended"`
	WorkspaceID    string                                 `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionOutcomeEvaluationEndedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionOutcomeEvaluationEndedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionPendingEventData struct {
	// ID of the resource that triggered the event.
	ID             string                  `json:"id" api:"required"`
	OrganizationID string                  `json:"organization_id" api:"required"`
	Type           constant.SessionPending `json:"type" default:"session.pending"`
	WorkspaceID    string                  `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionPendingEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionPendingEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionRequiresActionEventData struct {
	// ID of the resource that triggered the event.
	ID             string                         `json:"id" api:"required"`
	OrganizationID string                         `json:"organization_id" api:"required"`
	Type           constant.SessionRequiresAction `json:"type" default:"session.requires_action"`
	WorkspaceID    string                         `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionRequiresActionEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionRequiresActionEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionRunningEventData struct {
	// ID of the resource that triggered the event.
	ID             string                  `json:"id" api:"required"`
	OrganizationID string                  `json:"organization_id" api:"required"`
	Type           constant.SessionRunning `json:"type" default:"session.running"`
	WorkspaceID    string                  `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionRunningEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionRunningEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionStatusIdledEventData struct {
	// ID of the resource that triggered the event.
	ID             string                      `json:"id" api:"required"`
	OrganizationID string                      `json:"organization_id" api:"required"`
	Type           constant.SessionStatusIdled `json:"type" default:"session.status_idled"`
	WorkspaceID    string                      `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionStatusIdledEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionStatusIdledEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionStatusRescheduledEventData struct {
	// ID of the resource that triggered the event.
	ID             string                            `json:"id" api:"required"`
	OrganizationID string                            `json:"organization_id" api:"required"`
	Type           constant.SessionStatusRescheduled `json:"type" default:"session.status_rescheduled"`
	WorkspaceID    string                            `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionStatusRescheduledEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionStatusRescheduledEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionStatusRunStartedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                           `json:"id" api:"required"`
	OrganizationID string                           `json:"organization_id" api:"required"`
	Type           constant.SessionStatusRunStarted `json:"type" default:"session.status_run_started"`
	WorkspaceID    string                           `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionStatusRunStartedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionStatusRunStartedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionStatusTerminatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                           `json:"id" api:"required"`
	OrganizationID string                           `json:"organization_id" api:"required"`
	Type           constant.SessionStatusTerminated `json:"type" default:"session.status_terminated"`
	WorkspaceID    string                           `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionStatusTerminatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionStatusTerminatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionThreadCreatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                        `json:"id" api:"required"`
	OrganizationID string                        `json:"organization_id" api:"required"`
	Type           constant.SessionThreadCreated `json:"type" default:"session.thread_created"`
	WorkspaceID    string                        `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionThreadCreatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionThreadCreatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionThreadIdledEventData struct {
	// ID of the resource that triggered the event.
	ID             string                      `json:"id" api:"required"`
	OrganizationID string                      `json:"organization_id" api:"required"`
	Type           constant.SessionThreadIdled `json:"type" default:"session.thread_idled"`
	WorkspaceID    string                      `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionThreadIdledEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionThreadIdledEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookSessionThreadTerminatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                           `json:"id" api:"required"`
	OrganizationID string                           `json:"organization_id" api:"required"`
	Type           constant.SessionThreadTerminated `json:"type" default:"session.thread_terminated"`
	WorkspaceID    string                           `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookSessionThreadTerminatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookSessionThreadTerminatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultArchivedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                 `json:"id" api:"required"`
	OrganizationID string                 `json:"organization_id" api:"required"`
	Type           constant.VaultArchived `json:"type" default:"vault.archived"`
	WorkspaceID    string                 `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultArchivedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultArchivedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultCreatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                `json:"id" api:"required"`
	OrganizationID string                `json:"organization_id" api:"required"`
	Type           constant.VaultCreated `json:"type" default:"vault.created"`
	WorkspaceID    string                `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultCreatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultCreatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultCredentialArchivedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                           `json:"id" api:"required"`
	OrganizationID string                           `json:"organization_id" api:"required"`
	Type           constant.VaultCredentialArchived `json:"type" default:"vault_credential.archived"`
	// ID of the vault that owns this credential.
	VaultID     string `json:"vault_id" api:"required"`
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		VaultID        respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultCredentialArchivedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultCredentialArchivedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultCredentialCreatedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                          `json:"id" api:"required"`
	OrganizationID string                          `json:"organization_id" api:"required"`
	Type           constant.VaultCredentialCreated `json:"type" default:"vault_credential.created"`
	// ID of the vault that owns this credential.
	VaultID     string `json:"vault_id" api:"required"`
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		VaultID        respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultCredentialCreatedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultCredentialCreatedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultCredentialDeletedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                          `json:"id" api:"required"`
	OrganizationID string                          `json:"organization_id" api:"required"`
	Type           constant.VaultCredentialDeleted `json:"type" default:"vault_credential.deleted"`
	// ID of the vault that owns this credential.
	VaultID     string `json:"vault_id" api:"required"`
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		VaultID        respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultCredentialDeletedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultCredentialDeletedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultCredentialRefreshFailedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                                `json:"id" api:"required"`
	OrganizationID string                                `json:"organization_id" api:"required"`
	Type           constant.VaultCredentialRefreshFailed `json:"type" default:"vault_credential.refresh_failed"`
	// ID of the vault that owns this credential.
	VaultID     string `json:"vault_id" api:"required"`
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		VaultID        respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultCredentialRefreshFailedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultCredentialRefreshFailedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWebhookVaultDeletedEventData struct {
	// ID of the resource that triggered the event.
	ID             string                `json:"id" api:"required"`
	OrganizationID string                `json:"organization_id" api:"required"`
	Type           constant.VaultDeleted `json:"type" default:"vault.deleted"`
	WorkspaceID    string                `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		OrganizationID respjson.Field
		Type           respjson.Field
		WorkspaceID    respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWebhookVaultDeletedEventData) RawJSON() string { return r.JSON.raw }
func (r *BetaWebhookVaultDeletedEventData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type UnwrapWebhookEvent struct {
	// Unique event identifier for idempotency.
	ID string `json:"id" api:"required"`
	// RFC 3339 timestamp when the event occurred.
	CreatedAt time.Time                 `json:"created_at" api:"required" format:"date-time"`
	Data      BetaWebhookEventDataUnion `json:"data" api:"required"`
	// Object type. Always `event` for webhook payloads.
	Type constant.Event `json:"type" default:"event"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Data        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UnwrapWebhookEvent) RawJSON() string { return r.JSON.raw }
func (r *UnwrapWebhookEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

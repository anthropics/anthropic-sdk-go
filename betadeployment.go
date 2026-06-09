// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
)

// BetaDeploymentService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaDeploymentService] method instead.
type BetaDeploymentService struct {
	Options []option.RequestOption
}

// NewBetaDeploymentService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaDeploymentService(opts ...option.RequestOption) (r BetaDeploymentService) {
	r = BetaDeploymentService{}
	r.Options = opts
	return
}

// Create Deployment
func (r *BetaDeploymentService) New(ctx context.Context, params BetaDeploymentNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/deployments?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Deployment
func (r *BetaDeploymentService) Get(ctx context.Context, deploymentID string, query BetaDeploymentGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Deployment
func (r *BetaDeploymentService) Update(ctx context.Context, deploymentID string, params BetaDeploymentUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Deployments
func (r *BetaDeploymentService) List(ctx context.Context, params BetaDeploymentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsDeployment], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/deployments?beta=true"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Deployments
func (r *BetaDeploymentService) ListAutoPaging(ctx context.Context, params BetaDeploymentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsDeployment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive Deployment
func (r *BetaDeploymentService) Archive(ctx context.Context, deploymentID string, body BetaDeploymentArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s/archive?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Pause Deployment
func (r *BetaDeploymentService) Pause(ctx context.Context, deploymentID string, body BetaDeploymentPauseParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s/pause?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Run Deployment Now
func (r *BetaDeploymentService) Run(ctx context.Context, deploymentID string, body BetaDeploymentRunParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeploymentRun, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s/run?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Unpause Deployment
func (r *BetaDeploymentService) Unpause(ctx context.Context, deploymentID string, body BetaDeploymentUnpauseParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployments/%s/unpause?beta=true", deploymentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// The deployment's agent was archived.
type BetaManagedAgentsAgentArchivedDeploymentPausedReasonError struct {
	// Any of "agent_archived_error".
	Type BetaManagedAgentsAgentArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsAgentArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentArchivedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsAgentArchivedDeploymentPausedReasonErrorTypeAgentArchivedError BetaManagedAgentsAgentArchivedDeploymentPausedReasonErrorType = "agent_archived_error"
)

// A deployment is a configured instance of an agent — it binds the agent to
// everything needed to run it autonomously: an environment, credentials, initial
// events, and an optional schedule.
type BetaManagedAgentsDeployment struct {
	// Unique identifier for this deployment.
	ID string `json:"id" api:"required"`
	// A resolved agent reference with a concrete version.
	Agent BetaManagedAgentsAgentReference `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Description of what the deployment does.
	Description string `json:"description" api:"required"`
	// ID of the `environment` where sessions run.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Events sent to each session immediately after creation.
	InitialEvents []BetaManagedAgentsDeploymentInitialEventUnion `json:"initial_events" api:"required"`
	// Arbitrary key-value metadata. Maximum 16 pairs.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Human-readable name.
	Name string `json:"name" api:"required"`
	// Why a deployment is paused. Non-null exactly when `status` is `paused`.
	PausedReason BetaManagedAgentsDeploymentPausedReasonUnion `json:"paused_reason" api:"required"`
	// Resources attached to sessions created from this deployment. Echoes the input
	// minus write-only credentials.
	Resources []BetaManagedAgentsSessionResourceConfigUnion `json:"resources" api:"required"`
	// 5-field POSIX cron schedule with computed runtime timestamps.
	Schedule BetaManagedAgentsSchedule `json:"schedule" api:"required"`
	// Lifecycle status of a deployment.
	//
	// Any of "active", "paused".
	Status BetaManagedAgentsDeploymentStatus `json:"status" api:"required"`
	// Any of "deployment".
	Type BetaManagedAgentsDeploymentType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Vault IDs supplying stored credentials for sessions created from this
	// deployment.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		Agent         respjson.Field
		ArchivedAt    respjson.Field
		CreatedAt     respjson.Field
		Description   respjson.Field
		EnvironmentID respjson.Field
		InitialEvents respjson.Field
		Metadata      respjson.Field
		Name          respjson.Field
		PausedReason  respjson.Field
		Resources     respjson.Field
		Schedule      respjson.Field
		Status        respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		VaultIDs      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeployment) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeployment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeploymentType string

const (
	BetaManagedAgentsDeploymentTypeDeployment BetaManagedAgentsDeploymentType = "deployment"
)

// BetaManagedAgentsDeploymentInitialEventUnion contains all possible properties
// and values from [BetaManagedAgentsDeploymentUserMessageEvent],
// [BetaManagedAgentsDeploymentUserDefineOutcomeEvent],
// [BetaManagedAgentsDeploymentSystemMessageEvent].
//
// Use the [BetaManagedAgentsDeploymentInitialEventUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentInitialEventUnion struct {
	// This field is a union of
	// [[]BetaManagedAgentsDeploymentUserMessageEventContentUnion],
	// [[]BetaManagedAgentsSystemContentBlock]
	Content BetaManagedAgentsDeploymentInitialEventUnionContent `json:"content"`
	// Any of "user.message", "user.define_outcome", "system.message".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsDeploymentUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [BetaManagedAgentsDeploymentUserDefineOutcomeEvent].
	Rubric BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion `json:"rubric"`
	// This field is from variant [BetaManagedAgentsDeploymentUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	JSON          struct {
		Content       respjson.Field
		Type          respjson.Field
		Description   respjson.Field
		Rubric        respjson.Field
		MaxIterations respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentInitialEvent is implemented by each variant of
// [BetaManagedAgentsDeploymentInitialEventUnion] to add type safety for the return
// type of [BetaManagedAgentsDeploymentInitialEventUnion.AsAny]
type anyBetaManagedAgentsDeploymentInitialEvent interface {
	implBetaManagedAgentsDeploymentInitialEventUnion()
}

func (BetaManagedAgentsDeploymentUserMessageEvent) implBetaManagedAgentsDeploymentInitialEventUnion() {
}
func (BetaManagedAgentsDeploymentUserDefineOutcomeEvent) implBetaManagedAgentsDeploymentInitialEventUnion() {
}
func (BetaManagedAgentsDeploymentSystemMessageEvent) implBetaManagedAgentsDeploymentInitialEventUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentInitialEventUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsDeploymentUserMessageEvent:
//	case anthropic.BetaManagedAgentsDeploymentUserDefineOutcomeEvent:
//	case anthropic.BetaManagedAgentsDeploymentSystemMessageEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentInitialEventUnion) AsAny() anyBetaManagedAgentsDeploymentInitialEvent {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "system.message":
		return u.AsSystemMessage()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentInitialEventUnion) AsUserMessage() (v BetaManagedAgentsDeploymentUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentInitialEventUnion) AsUserDefineOutcome() (v BetaManagedAgentsDeploymentUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentInitialEventUnion) AsSystemMessage() (v BetaManagedAgentsDeploymentSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentInitialEventUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDeploymentInitialEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentInitialEventUnionContent is an implicit subunion of
// [BetaManagedAgentsDeploymentInitialEventUnion].
// BetaManagedAgentsDeploymentInitialEventUnionContent provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsDeploymentInitialEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsDeploymentUserMessageEventContentArray
// OfBetaManagedAgentsSystemContentBlockArray]
type BetaManagedAgentsDeploymentInitialEventUnionContent struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsDeploymentUserMessageEventContentUnion] instead of an
	// object.
	OfBetaManagedAgentsDeploymentUserMessageEventContentArray []BetaManagedAgentsDeploymentUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsSystemContentBlock] instead of an object.
	OfBetaManagedAgentsSystemContentBlockArray []BetaManagedAgentsSystemContentBlock `json:",inline"`
	JSON                                       struct {
		OfBetaManagedAgentsDeploymentUserMessageEventContentArray respjson.Field
		OfBetaManagedAgentsSystemContentBlockArray                respjson.Field
		raw                                                       string
	} `json:"-"`
}

func (r *BetaManagedAgentsDeploymentInitialEventUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func BetaManagedAgentsDeploymentInitialEventParamsOfUserMessage(content []BetaManagedAgentsUserMessageEventParamsContentUnion) BetaManagedAgentsDeploymentInitialEventParamsUnion {
	var userMessage BetaManagedAgentsUserMessageEventParams
	userMessage.Content = content
	return BetaManagedAgentsDeploymentInitialEventParamsUnion{OfUserMessage: &userMessage}
}

func BetaManagedAgentsDeploymentInitialEventParamsOfUserDefineOutcome[
	T BetaManagedAgentsFileRubricParams | BetaManagedAgentsTextRubricParams,
](description string, rubric T, type_ BetaManagedAgentsUserDefineOutcomeEventParamsType) BetaManagedAgentsDeploymentInitialEventParamsUnion {
	var userDefineOutcome BetaManagedAgentsUserDefineOutcomeEventParams
	userDefineOutcome.Description = description
	switch v := any(rubric).(type) {
	case BetaManagedAgentsFileRubricParams:
		userDefineOutcome.Rubric.OfFile = &v
	case BetaManagedAgentsTextRubricParams:
		userDefineOutcome.Rubric.OfText = &v
	}
	userDefineOutcome.Type = type_
	return BetaManagedAgentsDeploymentInitialEventParamsUnion{OfUserDefineOutcome: &userDefineOutcome}
}

func BetaManagedAgentsDeploymentInitialEventParamsOfSystemMessage(content []BetaManagedAgentsSystemContentBlockParam) BetaManagedAgentsDeploymentInitialEventParamsUnion {
	var systemMessage BetaManagedAgentsSystemMessageEventParams
	systemMessage.Content = content
	return BetaManagedAgentsDeploymentInitialEventParamsUnion{OfSystemMessage: &systemMessage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsDeploymentInitialEventParamsUnion struct {
	OfUserMessage       *BetaManagedAgentsUserMessageEventParams       `json:",omitzero,inline"`
	OfUserDefineOutcome *BetaManagedAgentsUserDefineOutcomeEventParams `json:",omitzero,inline"`
	OfSystemMessage     *BetaManagedAgentsSystemMessageEventParams     `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage, u.OfUserDefineOutcome, u.OfSystemMessage)
}
func (u *BetaManagedAgentsDeploymentInitialEventParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsDeploymentInitialEventParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfUserMessage) {
		return u.OfUserMessage
	} else if !param.IsOmitted(u.OfUserDefineOutcome) {
		return u.OfUserDefineOutcome
	} else if !param.IsOmitted(u.OfSystemMessage) {
		return u.OfSystemMessage
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) GetDescription() *string {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) GetRubric() *BetaManagedAgentsUserDefineOutcomeEventParamsRubricUnion {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Rubric
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) GetMaxIterations() *int64 {
	if vt := u.OfUserDefineOutcome; vt != nil && vt.MaxIterations.Valid() {
		return &vt.MaxIterations.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserDefineOutcome; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSystemMessage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaManagedAgentsDeploymentInitialEventParamsUnion) GetContent() (res betaManagedAgentsDeploymentInitialEventParamsUnionContent) {
	if vt := u.OfUserMessage; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfSystemMessage; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types
// [_[]BetaManagedAgentsUserMessageEventParamsContentUnion],
// [_[]BetaManagedAgentsSystemContentBlockParam]
type betaManagedAgentsDeploymentInitialEventParamsUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion:
//	case *[]anthropic.BetaManagedAgentsSystemContentBlockParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaManagedAgentsDeploymentInitialEventParamsUnionContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[BetaManagedAgentsDeploymentInitialEventParamsUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsUserMessageEventParams]("user.message"),
		apijson.Discriminator[BetaManagedAgentsUserDefineOutcomeEventParams]("user.define_outcome"),
		apijson.Discriminator[BetaManagedAgentsSystemMessageEventParams]("system.message"),
	)
}

// BetaManagedAgentsDeploymentPausedReasonUnion contains all possible properties
// and values from [BetaManagedAgentsManualDeploymentPausedReason],
// [BetaManagedAgentsErrorDeploymentPausedReason].
//
// Use the [BetaManagedAgentsDeploymentPausedReasonUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentPausedReasonUnion struct {
	// Any of "manual", "error".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsErrorDeploymentPausedReason].
	Error BetaManagedAgentsDeploymentPausedReasonErrorUnion `json:"error"`
	JSON  struct {
		Type  respjson.Field
		Error respjson.Field
		raw   string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentPausedReason is implemented by each variant of
// [BetaManagedAgentsDeploymentPausedReasonUnion] to add type safety for the return
// type of [BetaManagedAgentsDeploymentPausedReasonUnion.AsAny]
type anyBetaManagedAgentsDeploymentPausedReason interface {
	implBetaManagedAgentsDeploymentPausedReasonUnion()
}

func (BetaManagedAgentsManualDeploymentPausedReason) implBetaManagedAgentsDeploymentPausedReasonUnion() {
}
func (BetaManagedAgentsErrorDeploymentPausedReason) implBetaManagedAgentsDeploymentPausedReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentPausedReasonUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsManualDeploymentPausedReason:
//	case anthropic.BetaManagedAgentsErrorDeploymentPausedReason:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentPausedReasonUnion) AsAny() anyBetaManagedAgentsDeploymentPausedReason {
	switch u.Type {
	case "manual":
		return u.AsManual()
	case "error":
		return u.AsError()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentPausedReasonUnion) AsManual() (v BetaManagedAgentsManualDeploymentPausedReason) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonUnion) AsError() (v BetaManagedAgentsErrorDeploymentPausedReason) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentPausedReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDeploymentPausedReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentPausedReasonErrorUnion contains all possible
// properties and values from
// [BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError],
// [BetaManagedAgentsAgentArchivedDeploymentPausedReasonError],
// [BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError],
// [BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError],
// [BetaManagedAgentsFileNotFoundDeploymentPausedReasonError],
// [BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError],
// [BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError],
// [BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError],
// [BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError],
// [BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError],
// [BetaManagedAgentsVaultArchivedDeploymentPausedReasonError],
// [BetaManagedAgentsUnknownDeploymentPausedReasonError],
// [BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError],
// [BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError].
//
// Use the [BetaManagedAgentsDeploymentPausedReasonErrorUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentPausedReasonErrorUnion struct {
	// Any of "environment_archived_error", "agent_archived_error",
	// "environment_not_found_error", "vault_not_found_error", "file_not_found_error",
	// "session_resource_not_found_error", "workspace_archived_error",
	// "organization_disabled_error", "memory_store_archived_error",
	// "skill_not_found_error", "vault_archived_error", "unknown_error",
	// "self_hosted_resources_unsupported_error", "mcp_egress_blocked_error".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentPausedReasonError is implemented by each variant
// of [BetaManagedAgentsDeploymentPausedReasonErrorUnion] to add type safety for
// the return type of [BetaManagedAgentsDeploymentPausedReasonErrorUnion.AsAny]
type anyBetaManagedAgentsDeploymentPausedReasonError interface {
	implBetaManagedAgentsDeploymentPausedReasonErrorUnion()
}

func (BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsAgentArchivedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsFileNotFoundDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsVaultArchivedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsUnknownDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) implBetaManagedAgentsDeploymentPausedReasonErrorUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentPausedReasonErrorUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsAgentArchivedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsFileNotFoundDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsVaultArchivedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsUnknownDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError:
//	case anthropic.BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsAny() anyBetaManagedAgentsDeploymentPausedReasonError {
	switch u.Type {
	case "environment_archived_error":
		return u.AsEnvironmentArchivedError()
	case "agent_archived_error":
		return u.AsAgentArchivedError()
	case "environment_not_found_error":
		return u.AsEnvironmentNotFoundError()
	case "vault_not_found_error":
		return u.AsVaultNotFoundError()
	case "file_not_found_error":
		return u.AsFileNotFoundError()
	case "session_resource_not_found_error":
		return u.AsSessionResourceNotFoundError()
	case "workspace_archived_error":
		return u.AsWorkspaceArchivedError()
	case "organization_disabled_error":
		return u.AsOrganizationDisabledError()
	case "memory_store_archived_error":
		return u.AsMemoryStoreArchivedError()
	case "skill_not_found_error":
		return u.AsSkillNotFoundError()
	case "vault_archived_error":
		return u.AsVaultArchivedError()
	case "unknown_error":
		return u.AsUnknownError()
	case "self_hosted_resources_unsupported_error":
		return u.AsSelfHostedResourcesUnsupportedError()
	case "mcp_egress_blocked_error":
		return u.AsMCPEgressBlockedError()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsEnvironmentArchivedError() (v BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsAgentArchivedError() (v BetaManagedAgentsAgentArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsEnvironmentNotFoundError() (v BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsVaultNotFoundError() (v BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsFileNotFoundError() (v BetaManagedAgentsFileNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsSessionResourceNotFoundError() (v BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsWorkspaceArchivedError() (v BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsOrganizationDisabledError() (v BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsMemoryStoreArchivedError() (v BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsSkillNotFoundError() (v BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsVaultArchivedError() (v BetaManagedAgentsVaultArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsUnknownError() (v BetaManagedAgentsUnknownDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsSelfHostedResourcesUnsupportedError() (v BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) AsMCPEgressBlockedError() (v BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentPausedReasonErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDeploymentPausedReasonErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of a deployment.
type BetaManagedAgentsDeploymentStatus string

const (
	BetaManagedAgentsDeploymentStatusActive BetaManagedAgentsDeploymentStatus = "active"
	BetaManagedAgentsDeploymentStatusPaused BetaManagedAgentsDeploymentStatus = "paused"
)

// Privileged context for the accompanying turn and all subsequent turns, appended
// to the session's system context as a `role: "system"` turn rather than replacing
// the top-level system prompt.
type BetaManagedAgentsDeploymentSystemMessageEvent struct {
	// System content blocks to append. Text-only.
	Content []BetaManagedAgentsSystemContentBlock `json:"content" api:"required"`
	// Any of "system.message".
	Type BetaManagedAgentsDeploymentSystemMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeploymentSystemMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeploymentSystemMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeploymentSystemMessageEventType string

const (
	BetaManagedAgentsDeploymentSystemMessageEventTypeSystemMessage BetaManagedAgentsDeploymentSystemMessageEventType = "system.message"
)

// An outcome the agent should work toward. The agent begins work on receipt.
type BetaManagedAgentsDeploymentUserDefineOutcomeEvent struct {
	// What the agent should produce. This is the task specification.
	Description string `json:"description" api:"required"`
	// Rubric for grading the quality of an outcome.
	Rubric BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion `json:"rubric" api:"required"`
	// Any of "user.define_outcome".
	Type BetaManagedAgentsDeploymentUserDefineOutcomeEventType `json:"type" api:"required"`
	// Eval→revision cycles before giving up. Default 3, max 20.
	MaxIterations int64 `json:"max_iterations" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description   respjson.Field
		Rubric        respjson.Field
		Type          respjson.Field
		MaxIterations respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeploymentUserDefineOutcomeEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeploymentUserDefineOutcomeEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion contains all
// possible properties and values from [BetaManagedAgentsFileRubric],
// [BetaManagedAgentsTextRubric].
//
// Use the [BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion struct {
	// This field is from variant [BetaManagedAgentsFileRubric].
	FileID string `json:"file_id"`
	// Any of "file", "text".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsTextRubric].
	Content string `json:"content"`
	JSON    struct {
		FileID  respjson.Field
		Type    respjson.Field
		Content respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentUserDefineOutcomeEventRubric is implemented by
// each variant of [BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion]
// to add type safety for the return type of
// [BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny]
type anyBetaManagedAgentsDeploymentUserDefineOutcomeEventRubric interface {
	implBetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion()
}

func (BetaManagedAgentsFileRubric) implBetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion() {
}
func (BetaManagedAgentsTextRubric) implBetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsFileRubric:
//	case anthropic.BetaManagedAgentsTextRubric:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsAny() anyBetaManagedAgentsDeploymentUserDefineOutcomeEventRubric {
	switch u.Type {
	case "file":
		return u.AsFile()
	case "text":
		return u.AsText()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsFile() (v BetaManagedAgentsFileRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsText() (v BetaManagedAgentsTextRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeploymentUserDefineOutcomeEventType string

const (
	BetaManagedAgentsDeploymentUserDefineOutcomeEventTypeUserDefineOutcome BetaManagedAgentsDeploymentUserDefineOutcomeEventType = "user.define_outcome"
)

// A user message sent to the session.
type BetaManagedAgentsDeploymentUserMessageEvent struct {
	// Array of content blocks for the user message.
	Content []BetaManagedAgentsDeploymentUserMessageEventContentUnion `json:"content" api:"required"`
	// Any of "user.message".
	Type BetaManagedAgentsDeploymentUserMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeploymentUserMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeploymentUserMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentUserMessageEventContentUnion contains all possible
// properties and values from [BetaManagedAgentsTextBlock],
// [BetaManagedAgentsImageBlock], [BetaManagedAgentsDocumentBlock].
//
// Use the [BetaManagedAgentsDeploymentUserMessageEventContentUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentUserMessageEventContentUnion struct {
	// This field is from variant [BetaManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document".
	Type string `json:"type"`
	// This field is a union of [BetaManagedAgentsImageBlockSourceUnion],
	// [BetaManagedAgentsDocumentBlockSourceUnion]
	Source BetaManagedAgentsDeploymentUserMessageEventContentUnionSource `json:"source"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentUserMessageEventContent is implemented by each
// variant of [BetaManagedAgentsDeploymentUserMessageEventContentUnion] to add type
// safety for the return type of
// [BetaManagedAgentsDeploymentUserMessageEventContentUnion.AsAny]
type anyBetaManagedAgentsDeploymentUserMessageEventContent interface {
	implBetaManagedAgentsDeploymentUserMessageEventContentUnion()
}

func (BetaManagedAgentsTextBlock) implBetaManagedAgentsDeploymentUserMessageEventContentUnion()     {}
func (BetaManagedAgentsImageBlock) implBetaManagedAgentsDeploymentUserMessageEventContentUnion()    {}
func (BetaManagedAgentsDocumentBlock) implBetaManagedAgentsDeploymentUserMessageEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentUserMessageEventContentUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTextBlock:
//	case anthropic.BetaManagedAgentsImageBlock:
//	case anthropic.BetaManagedAgentsDocumentBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentUserMessageEventContentUnion) AsAny() anyBetaManagedAgentsDeploymentUserMessageEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentUserMessageEventContentUnion) AsText() (v BetaManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentUserMessageEventContentUnion) AsImage() (v BetaManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentUserMessageEventContentUnion) AsDocument() (v BetaManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentUserMessageEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDeploymentUserMessageEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentUserMessageEventContentUnionSource is an implicit
// subunion of [BetaManagedAgentsDeploymentUserMessageEventContentUnion].
// BetaManagedAgentsDeploymentUserMessageEventContentUnionSource provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsDeploymentUserMessageEventContentUnion].
type BetaManagedAgentsDeploymentUserMessageEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *BetaManagedAgentsDeploymentUserMessageEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeploymentUserMessageEventType string

const (
	BetaManagedAgentsDeploymentUserMessageEventTypeUserMessage BetaManagedAgentsDeploymentUserMessageEventType = "user.message"
)

// The deployment's environment was archived.
type BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError struct {
	// Any of "environment_archived_error".
	Type BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorTypeEnvironmentArchivedError BetaManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType = "environment_archived_error"
)

// The deployment's environment no longer exists.
type BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError struct {
	// Any of "environment_not_found_error".
	Type BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorTypeEnvironmentNotFoundError BetaManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType = "environment_not_found_error"
)

// A scheduled fire recorded a failed run whose error auto-pauses the deployment.
type BetaManagedAgentsErrorDeploymentPausedReason struct {
	// The error that triggered an auto-pause. Matches the failed run's `error.type`.
	Error BetaManagedAgentsDeploymentPausedReasonErrorUnion `json:"error" api:"required"`
	// Any of "error".
	Type BetaManagedAgentsErrorDeploymentPausedReasonType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsErrorDeploymentPausedReason) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsErrorDeploymentPausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsErrorDeploymentPausedReasonType string

const (
	BetaManagedAgentsErrorDeploymentPausedReasonTypeError BetaManagedAgentsErrorDeploymentPausedReasonType = "error"
)

// A file resource referenced by the deployment no longer exists.
type BetaManagedAgentsFileNotFoundDeploymentPausedReasonError struct {
	// Any of "file_not_found_error".
	Type BetaManagedAgentsFileNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileNotFoundDeploymentPausedReasonError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsFileNotFoundDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsFileNotFoundDeploymentPausedReasonErrorTypeFileNotFoundError BetaManagedAgentsFileNotFoundDeploymentPausedReasonErrorType = "file_not_found_error"
)

// A file mounted into each session's container.
type BetaManagedAgentsFileResourceConfig struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileResourceConfigType `json:"type" api:"required"`
	// Mount path in the container. Defaults to `/mnt/session/uploads/<file_id>`.
	MountPath string `json:"mount_path" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsFileResourceConfigType string

const (
	BetaManagedAgentsFileResourceConfigTypeFile BetaManagedAgentsFileResourceConfigType = "file"
)

// A GitHub repository mounted into each session's container. The authorization
// token is write-only and never returned.
type BetaManagedAgentsGitHubRepositoryResourceConfig struct {
	// Any of "github_repository".
	Type BetaManagedAgentsGitHubRepositoryResourceConfigType `json:"type" api:"required"`
	// Github URL of the repository
	URL string `json:"url" api:"required"`
	// Branch or commit to check out. Defaults to the repository's default branch.
	Checkout BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion `json:"checkout" api:"nullable"`
	// Mount path in the container. Defaults to `/workspace/<repo-name>`.
	MountPath string `json:"mount_path" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		Checkout    respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsGitHubRepositoryResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsGitHubRepositoryResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsGitHubRepositoryResourceConfigType string

const (
	BetaManagedAgentsGitHubRepositoryResourceConfigTypeGitHubRepository BetaManagedAgentsGitHubRepositoryResourceConfigType = "github_repository"
)

// BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion contains all
// possible properties and values from [BetaManagedAgentsBranchCheckout],
// [BetaManagedAgentsCommitCheckout].
//
// Use the [BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion struct {
	// This field is from variant [BetaManagedAgentsBranchCheckout].
	Name string `json:"name"`
	// Any of "branch", "commit".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsCommitCheckout].
	Sha  string `json:"sha"`
	JSON struct {
		Name respjson.Field
		Type respjson.Field
		Sha  respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsGitHubRepositoryResourceConfigCheckout is implemented by
// each variant of [BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion]
// to add type safety for the return type of
// [BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny]
type anyBetaManagedAgentsGitHubRepositoryResourceConfigCheckout interface {
	implBetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion()
}

func (BetaManagedAgentsBranchCheckout) implBetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion() {
}
func (BetaManagedAgentsCommitCheckout) implBetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsBranchCheckout:
//	case anthropic.BetaManagedAgentsCommitCheckout:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsAny() anyBetaManagedAgentsGitHubRepositoryResourceConfigCheckout {
	switch u.Type {
	case "branch":
		return u.AsBranch()
	case "commit":
		return u.AsCommit()
	}
	return nil
}

func (u BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsBranch() (v BetaManagedAgentsBranchCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsCommit() (v BetaManagedAgentsCommitCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The caller invoked the pause endpoint on the deployment.
type BetaManagedAgentsManualDeploymentPausedReason struct {
	// Any of "manual".
	Type BetaManagedAgentsManualDeploymentPausedReasonType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsManualDeploymentPausedReason) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsManualDeploymentPausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsManualDeploymentPausedReasonType string

const (
	BetaManagedAgentsManualDeploymentPausedReasonTypeManual BetaManagedAgentsManualDeploymentPausedReasonType = "manual"
)

// An MCP server host used by the deployment's agent is blocked by the
// environment's network policy.
type BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError struct {
	// Any of "mcp_egress_blocked_error".
	Type BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorTypeMCPEgressBlockedError BetaManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType = "mcp_egress_blocked_error"
)

// A memory store referenced by the deployment is archived.
type BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError struct {
	// Any of "memory_store_archived_error".
	Type BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorTypeMemoryStoreArchivedError BetaManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType = "memory_store_archived_error"
)

// A memory store attached to each session created from this deployment.
type BetaManagedAgentsMemoryStoreResourceConfig struct {
	// The memory store ID (memstore\_...). Must belong to the caller's organization
	// and workspace.
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type BetaManagedAgentsMemoryStoreResourceConfigType `json:"type" api:"required"`
	// Access mode for an attached memory store.
	//
	// Any of "read_write", "read_only".
	Access BetaManagedAgentsMemoryStoreResourceConfigAccess `json:"access" api:"nullable"`
	// Per-attachment guidance for the agent on how to use this store. Rendered into
	// the memory section of the system prompt. Max 4096 chars.
	Instructions string `json:"instructions" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		Access        respjson.Field
		Instructions  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryStoreResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryStoreResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryStoreResourceConfigType string

const (
	BetaManagedAgentsMemoryStoreResourceConfigTypeMemoryStore BetaManagedAgentsMemoryStoreResourceConfigType = "memory_store"
)

// Access mode for an attached memory store.
type BetaManagedAgentsMemoryStoreResourceConfigAccess string

const (
	BetaManagedAgentsMemoryStoreResourceConfigAccessReadWrite BetaManagedAgentsMemoryStoreResourceConfigAccess = "read_write"
	BetaManagedAgentsMemoryStoreResourceConfigAccessReadOnly  BetaManagedAgentsMemoryStoreResourceConfigAccess = "read_only"
)

// The deployment's organization is disabled.
type BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError struct {
	// Any of "organization_disabled_error".
	Type BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorTypeOrganizationDisabledError BetaManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType = "organization_disabled_error"
)

// 5-field POSIX cron schedule with computed runtime timestamps.
type BetaManagedAgentsSchedule struct {
	// 5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g.,
	// "0 9 \* \* 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean
	// Sunday. Extended cron syntax - seconds or year fields, and the special
	// characters L, W, #, and ? - is not supported, nor are predefined shortcuts
	// (@daily).
	Expression string `json:"expression" api:"required"`
	// IANA timezone identifier (e.g., "America/Los_Angeles", "UTC").
	Timezone string `json:"timezone" api:"required"`
	// Any of "cron".
	Type BetaManagedAgentsScheduleType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	LastRunAt time.Time `json:"last_run_at" api:"nullable" format:"date-time"`
	// Up to 5 timestamps of upcoming cron occurrences. Non-empty for active and paused
	// deployments (reflects what the schedule would do if unpaused); empty once the
	// deployment is archived (`archived_at` set). Each fire is offset by a small
	// per-schedule jitter, so a run will actually start at or shortly after its listed
	// time.
	UpcomingRunsAt []time.Time `json:"upcoming_runs_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Expression     respjson.Field
		Timezone       respjson.Field
		Type           respjson.Field
		LastRunAt      respjson.Field
		UpcomingRunsAt respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSchedule) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSchedule) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsScheduleType string

const (
	BetaManagedAgentsScheduleTypeCron BetaManagedAgentsScheduleType = "cron"
)

// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
// timezone.
//
// The properties Expression, Timezone, Type are required.
type BetaManagedAgentsScheduleParams struct {
	// 5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g.,
	// "0 9 \* \* 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean
	// Sunday. Extended cron syntax - seconds or year fields, and the special
	// characters L, W, #, and ? - is not supported, nor are predefined shortcuts
	// (@daily).
	Expression string `json:"expression" api:"required"`
	// Required. IANA timezone identifier (e.g., "America/Los_Angeles", "UTC").
	// Validated against the IANA timezone database.
	Timezone string `json:"timezone" api:"required"`
	// Any of "cron".
	Type BetaManagedAgentsScheduleParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsScheduleParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsScheduleParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsScheduleParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsScheduleParamsType string

const (
	BetaManagedAgentsScheduleParamsTypeCron BetaManagedAgentsScheduleParamsType = "cron"
)

// The deployment configures resources, but its environment is self-hosted and
// cannot mount them.
type BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError struct {
	// Any of "self_hosted_resources_unsupported_error".
	Type BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorTypeSelfHostedResourcesUnsupportedError BetaManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType = "self_hosted_resources_unsupported_error"
)

// BetaManagedAgentsSessionResourceConfigUnion contains all possible properties and
// values from [BetaManagedAgentsGitHubRepositoryResourceConfig],
// [BetaManagedAgentsFileResourceConfig],
// [BetaManagedAgentsMemoryStoreResourceConfig].
//
// Use the [BetaManagedAgentsSessionResourceConfigUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionResourceConfigUnion struct {
	// Any of "github_repository", "file", "memory_store".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResourceConfig].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsGitHubRepositoryResourceConfig].
	Checkout  BetaManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion `json:"checkout"`
	MountPath string                                                       `json:"mount_path"`
	// This field is from variant [BetaManagedAgentsFileResourceConfig].
	FileID string `json:"file_id"`
	// This field is from variant [BetaManagedAgentsMemoryStoreResourceConfig].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [BetaManagedAgentsMemoryStoreResourceConfig].
	Access BetaManagedAgentsMemoryStoreResourceConfigAccess `json:"access"`
	// This field is from variant [BetaManagedAgentsMemoryStoreResourceConfig].
	Instructions string `json:"instructions"`
	JSON         struct {
		Type          respjson.Field
		URL           respjson.Field
		Checkout      respjson.Field
		MountPath     respjson.Field
		FileID        respjson.Field
		MemoryStoreID respjson.Field
		Access        respjson.Field
		Instructions  respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionResourceConfig is implemented by each variant of
// [BetaManagedAgentsSessionResourceConfigUnion] to add type safety for the return
// type of [BetaManagedAgentsSessionResourceConfigUnion.AsAny]
type anyBetaManagedAgentsSessionResourceConfig interface {
	implBetaManagedAgentsSessionResourceConfigUnion()
}

func (BetaManagedAgentsGitHubRepositoryResourceConfig) implBetaManagedAgentsSessionResourceConfigUnion() {
}
func (BetaManagedAgentsFileResourceConfig) implBetaManagedAgentsSessionResourceConfigUnion()        {}
func (BetaManagedAgentsMemoryStoreResourceConfig) implBetaManagedAgentsSessionResourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionResourceConfigUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsGitHubRepositoryResourceConfig:
//	case anthropic.BetaManagedAgentsFileResourceConfig:
//	case anthropic.BetaManagedAgentsMemoryStoreResourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionResourceConfigUnion) AsAny() anyBetaManagedAgentsSessionResourceConfig {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	case "memory_store":
		return u.AsMemoryStore()
	}
	return nil
}

func (u BetaManagedAgentsSessionResourceConfigUnion) AsGitHubRepository() (v BetaManagedAgentsGitHubRepositoryResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionResourceConfigUnion) AsFile() (v BetaManagedAgentsFileResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionResourceConfigUnion) AsMemoryStore() (v BetaManagedAgentsMemoryStoreResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionResourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionResourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A referenced resource no longer exists and its kind was not reported.
type BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError struct {
	// Any of "session_resource_not_found_error".
	Type BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorTypeSessionResourceNotFoundError BetaManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType = "session_resource_not_found_error"
)

// A skill referenced by the deployment's agent no longer exists.
type BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError struct {
	// Any of "skill_not_found_error".
	Type BetaManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsSkillNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsSkillNotFoundDeploymentPausedReasonErrorTypeSkillNotFoundError BetaManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType = "skill_not_found_error"
)

// An unrecognized error auto-paused the deployment. A fallback variant; matches a
// run whose `error.type` is `unknown_error`.
type BetaManagedAgentsUnknownDeploymentPausedReasonError struct {
	// Any of "unknown_error".
	Type BetaManagedAgentsUnknownDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUnknownDeploymentPausedReasonError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUnknownDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUnknownDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsUnknownDeploymentPausedReasonErrorTypeUnknownError BetaManagedAgentsUnknownDeploymentPausedReasonErrorType = "unknown_error"
)

// A vault referenced by the deployment is archived.
type BetaManagedAgentsVaultArchivedDeploymentPausedReasonError struct {
	// Any of "vault_archived_error".
	Type BetaManagedAgentsVaultArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsVaultArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsVaultArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsVaultArchivedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsVaultArchivedDeploymentPausedReasonErrorTypeVaultArchivedError BetaManagedAgentsVaultArchivedDeploymentPausedReasonErrorType = "vault_archived_error"
)

// A vault referenced by the deployment no longer exists.
type BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError struct {
	// Any of "vault_not_found_error".
	Type BetaManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsVaultNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsVaultNotFoundDeploymentPausedReasonErrorTypeVaultNotFoundError BetaManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType = "vault_not_found_error"
)

// The deployment's workspace was archived.
type BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError struct {
	// Any of "workspace_archived_error".
	Type BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType string

const (
	BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorTypeWorkspaceArchivedError BetaManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType = "workspace_archived_error"
)

type BetaDeploymentNewParams struct {
	// Agent to deploy. Accepts the `agent` ID string, which pins the latest version,
	// or an `agent` object with both id and version specified. The agent must exist
	// and not be archived.
	Agent BetaDeploymentNewParamsAgentUnion `json:"agent,omitzero" api:"required"`
	// ID of the `environment` defining the container configuration for sessions
	// created from this deployment.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Events to send to each session immediately after creation. At least 1,
	// maximum 50.
	InitialEvents []BetaManagedAgentsDeploymentInitialEventParamsUnion `json:"initial_events,omitzero" api:"required"`
	// Human-readable name for the deployment.
	Name string `json:"name" api:"required"`
	// Description of what the deployment does.
	Description param.Opt[string] `json:"description,omitzero"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Resources (e.g. repositories, files) to mount into each session's container.
	// Maximum 500.
	Resources []BetaDeploymentNewParamsResourceUnion `json:"resources,omitzero"`
	// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
	// timezone.
	Schedule BetaManagedAgentsScheduleParams `json:"schedule,omitzero"`
	// Vault IDs for stored credentials the agent can use during sessions created from
	// this deployment. Maximum 50.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaDeploymentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaDeploymentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDeploymentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDeploymentNewParamsAgentUnion struct {
	OfString                  param.Opt[string]             `json:",omitzero,inline"`
	OfBetaManagedAgentsAgents *BetaManagedAgentsAgentParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDeploymentNewParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBetaManagedAgentsAgents)
}
func (u *BetaDeploymentNewParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDeploymentNewParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaManagedAgentsAgents) {
		return u.OfBetaManagedAgentsAgents
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDeploymentNewParamsResourceUnion struct {
	OfGitHubRepository *BetaManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *BetaManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	OfMemoryStore      *BetaManagedAgentsMemoryStoreResourceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDeploymentNewParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile, u.OfMemoryStore)
}
func (u *BetaDeploymentNewParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDeploymentNewParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetCheckout() *BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetAccess() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Access)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetInstructions() *string {
	if vt := u.OfMemoryStore; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentNewParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaDeploymentNewParamsResourceUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[BetaManagedAgentsFileResourceParams]("file"),
		apijson.Discriminator[BetaManagedAgentsMemoryStoreResourceParam]("memory_store"),
	)
}

type BetaDeploymentGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDeploymentUpdateParams struct {
	// Description. Omit to preserve; send empty string or null to clear.
	Description param.Opt[string] `json:"description,omitzero"`
	// ID of the `environment` where sessions run. Omit to preserve. Cannot be cleared.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.
	Name param.Opt[string] `json:"name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Session resources. Full replacement. Omit to preserve; send empty array or null
	// to clear. Maximum 500.
	Resources []BetaDeploymentUpdateParamsResourceUnion `json:"resources,omitzero"`
	// Vault IDs. Full replacement. Omit to preserve; send empty array or null to
	// clear. Maximum 50.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Agent to deploy. Accepts the `agent` ID string, which re-pins to the latest
	// version, or an `agent` object with both id and version specified. Omit to
	// preserve. Cannot be cleared.
	Agent BetaDeploymentUpdateParamsAgentUnion `json:"agent,omitzero"`
	// Initial events. Full replacement. Omit to preserve. Cannot be cleared. At least
	// 1, maximum 50.
	InitialEvents []BetaManagedAgentsDeploymentInitialEventParamsUnion `json:"initial_events,omitzero"`
	// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
	// timezone.
	Schedule BetaManagedAgentsScheduleParams `json:"schedule,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaDeploymentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaDeploymentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDeploymentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDeploymentUpdateParamsAgentUnion struct {
	OfString                  param.Opt[string]             `json:",omitzero,inline"`
	OfBetaManagedAgentsAgents *BetaManagedAgentsAgentParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDeploymentUpdateParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBetaManagedAgentsAgents)
}
func (u *BetaDeploymentUpdateParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDeploymentUpdateParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaManagedAgentsAgents) {
		return u.OfBetaManagedAgentsAgents
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDeploymentUpdateParamsResourceUnion struct {
	OfGitHubRepository *BetaManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *BetaManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	OfMemoryStore      *BetaManagedAgentsMemoryStoreResourceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDeploymentUpdateParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile, u.OfMemoryStore)
}
func (u *BetaDeploymentUpdateParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDeploymentUpdateParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetCheckout() *BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetAccess() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Access)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetInstructions() *string {
	if vt := u.OfMemoryStore; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDeploymentUpdateParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaDeploymentUpdateParamsResourceUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[BetaManagedAgentsFileResourceParams]("file"),
		apijson.Discriminator[BetaManagedAgentsMemoryStoreResourceParam]("memory_store"),
	)
}

type BetaDeploymentListParams struct {
	// Filter by agent ID.
	AgentID param.Opt[string] `query:"agent_id,omitzero" json:"-"`
	// Return deployments created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return deployments created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When true, includes archived deployments. Default: false (exclude archived).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter by status: active or paused. Omit for both. To include archived
	// deployments, use include_archived instead; the two cannot be combined.
	//
	// Any of "active", "paused".
	Status BetaManagedAgentsDeploymentStatus `query:"status,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaDeploymentListParams]'s query parameters as
// `url.Values`.
func (r BetaDeploymentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaDeploymentArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDeploymentPauseParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDeploymentRunParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDeploymentUnpauseParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaComplianceReviewCaseService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaComplianceReviewCaseService] method instead.
type BetaComplianceReviewCaseService struct {
	Options []option.RequestOption
}

// NewBetaComplianceReviewCaseService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaComplianceReviewCaseService(opts ...option.RequestOption) (r BetaComplianceReviewCaseService) {
	r = BetaComplianceReviewCaseService{}
	r.Options = opts
	return
}

// Read one review case.
//
// A case that the API key can't reach, or that Anthropic hasn't sent, returns the
// same 404 as a case that doesn't exist.
//
// An Admin API key with the `read:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewCaseService) Get(ctx context.Context, caseID string, query BetaComplianceReviewCaseGetParams, opts ...option.RequestOption) (res *BetaComplianceReviewCase, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	if caseID == "" {
		err = errors.New("missing required case_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/compliance/review/cases/%s", requestconfig.EscapePathSegment(caseID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.compliance.review.cases.retrieve", "/v1/compliance/review/cases/{case_id}"), requestconfig.WithSpanAttribute("anthropic.compliance.review.case.id", caseID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the review cases of every organization that the API key can reach, newest
// first by default. Both open and closed cases are returned unless `status` is
// set.
//
// To sync, store the largest `updated_at` you have read, and send it as the
// at-or-after filter on `updated_at`, sorted by `updated_at` ascending. That finds
// new cases, newly closed cases, and cases with a new notice. Read every page,
// match results by `id`, and store the new largest value after the last page.
//
// A request with the `anthropic-user-profile-id` header and no `user_profile_id`
// parameter returns a 400 error.
//
// An Admin API key with the `read:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewCaseService) List(ctx context.Context, params BetaComplianceReviewCaseListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaComplianceReviewCase], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/compliance/review/cases"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.compliance.review.cases.list", "/v1/compliance/review/cases", "created_at[gt]", "created_at[gte]", "created_at[lt]", "created_at[lte]", "limit", "order", "order_by", "page", "status", "updated_at[gt]", "updated_at[gte]", "updated_at[lt]", "updated_at[lte]", "user_profile_id"))
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

// List the review cases of every organization that the API key can reach, newest
// first by default. Both open and closed cases are returned unless `status` is
// set.
//
// To sync, store the largest `updated_at` you have read, and send it as the
// at-or-after filter on `updated_at`, sorted by `updated_at` ascending. That finds
// new cases, newly closed cases, and cases with a new notice. Read every page,
// match results by `id`, and store the new largest value after the last page.
//
// A request with the `anthropic-user-profile-id` header and no `user_profile_id`
// parameter returns a 400 error.
//
// An Admin API key with the `read:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewCaseService) ListAutoPaging(ctx context.Context, params BetaComplianceReviewCaseListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaComplianceReviewCase] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Send your conclusion on a review case. A response closes the case, and a closed
// case can't be reopened.
//
// Only the first response counts. A request is safe to send again: if the closed
// case has a response with the same field values, in any list order, the API
// returns 200 and changes nothing. If the closed case has a different response, or
// none, the API returns a 400 with `case_closed` in `error.details.error_code`.
//
// If your end customer reviews this case instead of you, the API returns a 403
// with `case_reviewed_by_end_customer` in `error.details.error_code`.
//
// An Admin API key with the `write:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewCaseService) Respond(ctx context.Context, caseID string, params BetaComplianceReviewCaseRespondParams, opts ...option.RequestOption) (res *BetaComplianceReviewCase, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	if caseID == "" {
		err = errors.New("missing required case_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/compliance/review/cases/%s/respond", requestconfig.EscapePathSegment(caseID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.compliance.review.cases.respond", "/v1/compliance/review/cases/{case_id}/respond"), requestconfig.WithSpanAttribute("anthropic.compliance.review.case.id", caseID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// A record of an organization's messages that Anthropic flagged and sent to the
// organization to review. It lists the messages by ID and doesn't include their
// content.
type BetaComplianceReviewCase struct {
	// The unique ID of the case (`rvc_...`).
	ID string `json:"id" api:"required"`
	// When Anthropic sent the case. The time you have to respond starts here.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The action that Anthropic took before it sent the case. Null when it took none.
	EnforcementAction BetaComplianceReviewEnforcementAction `json:"enforcement_action" api:"required"`
	// The UUID of the organization that sent the flagged messages, which the case
	// belongs to.
	OrganizationID string `json:"organization_id" api:"required"`
	// Anthropic's labels for why the messages were flagged, such as `cyber` or
	// `fraud`. There is at least one. New values can appear, so accept a value that
	// you don't know.
	PolicyCategories []BetaComplianceReviewPolicyCategory `json:"policy_categories" api:"required"`
	// Whether Anthropic is still waiting for a response.
	Status BetaComplianceReviewCaseStatusUnion `json:"status" api:"required"`
	// The flagged messages. There is at least one, and the list isn't paginated.
	Transcripts []BetaComplianceReviewTranscript `json:"transcripts" api:"required"`
	// The object type, always `compliance_review_case`.
	Type constant.ComplianceReviewCase `json:"type" default:"compliance_review_case"`
	// When the case last changed. It changes when Anthropic sends the case, when the
	// case closes, and when Anthropic sends a notice about the case.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The ID of the workspace that sent the flagged messages (`wrkspc_...`). Null when
	// the messages have no workspace, and it can be null for the default workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		CreatedAt         respjson.Field
		EnforcementAction respjson.Field
		OrganizationID    respjson.Field
		PolicyCategories  respjson.Field
		Status            respjson.Field
		Transcripts       respjson.Field
		Type              respjson.Field
		UpdatedAt         respjson.Field
		WorkspaceID       respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewCase) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewCase) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Why a review case closed.
type BetaComplianceReviewCaseCloseReason struct {
	// `responded` means that Anthropic recorded a response. `window_ended` means that
	// the time to respond ended with no response. `ended_by_anthropic` means that
	// Anthropic ended the case before that time was over. New values can appear, so
	// accept a value that you don't know.
	//
	// Any of "responded", "window_ended", "ended_by_anthropic".
	Type BetaComplianceReviewCaseCloseReasonType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewCaseCloseReason) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewCaseCloseReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `responded` means that Anthropic recorded a response. `window_ended` means that
// the time to respond ended with no response. `ended_by_anthropic` means that
// Anthropic ended the case before that time was over. New values can appear, so
// accept a value that you don't know.
type BetaComplianceReviewCaseCloseReasonType string

const (
	BetaComplianceReviewCaseCloseReasonTypeResponded        BetaComplianceReviewCaseCloseReasonType = "responded"
	BetaComplianceReviewCaseCloseReasonTypeWindowEnded      BetaComplianceReviewCaseCloseReasonType = "window_ended"
	BetaComplianceReviewCaseCloseReasonTypeEndedByAnthropic BetaComplianceReviewCaseCloseReasonType = "ended_by_anthropic"
)

// Anthropic has stopped waiting for a response. A closed case stays closed.
type BetaComplianceReviewCaseClosed struct {
	// Why the case closed. It doesn't change afterwards.
	Reason BetaComplianceReviewCaseCloseReason `json:"reason" api:"required"`
	Type   constant.Closed                     `json:"type" default:"closed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Reason      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewCaseClosed) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewCaseClosed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic is waiting for a response to the case.
type BetaComplianceReviewCaseOpen struct {
	Type constant.Open `json:"type" default:"open"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewCaseOpen) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewCaseOpen) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The field that a list of review cases is sorted by: `created_at` or
// `updated_at`.
type BetaComplianceReviewCaseOrderBy string

const (
	// Sort by when Anthropic sent the case.
	BetaComplianceReviewCaseOrderByCreatedAt BetaComplianceReviewCaseOrderBy = "created_at"
	// Sort by when the case last changed.
	BetaComplianceReviewCaseOrderByUpdatedAt BetaComplianceReviewCaseOrderBy = "updated_at"
)

// BetaComplianceReviewCaseStatusUnion contains all possible properties and values
// from [BetaComplianceReviewCaseOpen], [BetaComplianceReviewCaseClosed].
//
// Use the [BetaComplianceReviewCaseStatusUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewCaseStatusUnion struct {
	// Any of "open", "closed".
	Type string `json:"type"`
	// This field is from variant [BetaComplianceReviewCaseClosed].
	Reason BetaComplianceReviewCaseCloseReason `json:"reason"`
	JSON   struct {
		Type   respjson.Field
		Reason respjson.Field
		raw    string
	} `json:"-"`
}

// anyBetaComplianceReviewCaseStatus is implemented by each variant of
// [BetaComplianceReviewCaseStatusUnion] to add type safety for the return type of
// [BetaComplianceReviewCaseStatusUnion.AsAny]
type anyBetaComplianceReviewCaseStatus interface {
	implBetaComplianceReviewCaseStatusUnion()
}

func (BetaComplianceReviewCaseOpen) implBetaComplianceReviewCaseStatusUnion()   {}
func (BetaComplianceReviewCaseClosed) implBetaComplianceReviewCaseStatusUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewCaseStatusUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewCaseOpen:
//	case anthropic.BetaComplianceReviewCaseClosed:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewCaseStatusUnion) AsAny() anyBetaComplianceReviewCaseStatus {
	switch u.Type {
	case "open":
		return u.AsOpen()
	case "closed":
		return u.AsClosed()
	}
	return nil
}

func (u BetaComplianceReviewCaseStatusUnion) AsOpen() (v BetaComplianceReviewCaseOpen) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaComplianceReviewCaseStatusUnion) AsClosed() (v BetaComplianceReviewCaseClosed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewCaseStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewCaseStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A review case status to filter by: `open` or `closed`.
type BetaComplianceReviewCaseStatusFilter string

const (
	// Cases that Anthropic is waiting for a response to.
	BetaComplianceReviewCaseStatusFilterOpen BetaComplianceReviewCaseStatusFilter = "open"
	// Cases that Anthropic has stopped waiting on.
	BetaComplianceReviewCaseStatusFilterClosed BetaComplianceReviewCaseStatusFilter = "closed"
)

// A conclusion on flagged messages: `violation` or `no_violation`.
type BetaComplianceReviewDetermination string

const (
	// The messages violated a policy.
	BetaComplianceReviewDeterminationViolation BetaComplianceReviewDetermination = "violation"
	// The messages violated no policy.
	BetaComplianceReviewDeterminationNoViolation BetaComplianceReviewDetermination = "no_violation"
)

// An action that Anthropic took on the flagged traffic.
type BetaComplianceReviewEnforcementAction struct {
	// The kind of action. New values can appear, so accept a value that you don't
	// know.
	Type string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewEnforcementAction) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewEnforcementAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaComplianceReviewPolicyCategory string

const (
	BetaComplianceReviewPolicyCategoryCyber               BetaComplianceReviewPolicyCategory = "cyber"
	BetaComplianceReviewPolicyCategoryBio                 BetaComplianceReviewPolicyCategory = "bio"
	BetaComplianceReviewPolicyCategoryChem                BetaComplianceReviewPolicyCategory = "chem"
	BetaComplianceReviewPolicyCategoryRadnuc              BetaComplianceReviewPolicyCategory = "radnuc"
	BetaComplianceReviewPolicyCategoryCgw                 BetaComplianceReviewPolicyCategory = "cgw"
	BetaComplianceReviewPolicyCategoryMgw                 BetaComplianceReviewPolicyCategory = "mgw"
	BetaComplianceReviewPolicyCategoryIo                  BetaComplianceReviewPolicyCategory = "io"
	BetaComplianceReviewPolicyCategoryFraud               BetaComplianceReviewPolicyCategory = "fraud"
	BetaComplianceReviewPolicyCategoryViolentExtremism    BetaComplianceReviewPolicyCategory = "violent_extremism"
	BetaComplianceReviewPolicyCategorySurveillance        BetaComplianceReviewPolicyCategory = "surveillance"
	BetaComplianceReviewPolicyCategoryHumanExploitation   BetaComplianceReviewPolicyCategory = "human_exploitation"
	BetaComplianceReviewPolicyCategorySelfHarm            BetaComplianceReviewPolicyCategory = "self_harm"
	BetaComplianceReviewPolicyCategoryAdultContent        BetaComplianceReviewPolicyCategory = "adult_content"
	BetaComplianceReviewPolicyCategoryReasoningExtraction BetaComplianceReviewPolicyCategory = "reasoning_extraction"
	BetaComplianceReviewPolicyCategoryFrontierLLM         BetaComplianceReviewPolicyCategory = "frontier_llm"
	BetaComplianceReviewPolicyCategoryCompromisedKey      BetaComplianceReviewPolicyCategory = "compromised_key"
	BetaComplianceReviewPolicyCategoryBenign              BetaComplianceReviewPolicyCategory = "benign"
	BetaComplianceReviewPolicyCategoryOther               BetaComplianceReviewPolicyCategory = "other"
	BetaComplianceReviewPolicyCategoryUndetermined        BetaComplianceReviewPolicyCategory = "undetermined"
)

// One flagged message of a review case.
type BetaComplianceReviewTranscript struct {
	// The ID of the flagged message.
	ItemID string `json:"item_id" api:"required"`
	// The ID of the user profile that the message belongs to. Null when the message
	// belongs to none, or when the profile can no longer be looked up.
	UserProfileID string `json:"user_profile_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ItemID        respjson.Field
		UserProfileID respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewTranscript) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewTranscript) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Your conclusion on one message of a review case.
//
// The properties Determination, ItemID are required.
type BetaComplianceReviewVerdictParam struct {
	// Your conclusion on this message.
	//
	// Any of "violation", "no_violation".
	Determination BetaComplianceReviewDetermination `json:"determination,omitzero" api:"required"`
	// The `item_id` of one of the case's `transcripts`.
	ItemID string `json:"item_id" api:"required"`
	// The policy category that this message violated, from the values that the
	// request's `policy_categories` accepts. Optional with `violation`. Leave it out
	// with `no_violation`.
	//
	// Any of "cyber", "bio", "chem", "radnuc", "cgw", "mgw", "io", "fraud",
	// "violent_extremism", "surveillance", "human_exploitation", "self_harm",
	// "adult_content", "reasoning_extraction", "frontier_llm", "compromised_key",
	// "other".
	PolicyCategory BetaComplianceReviewViolatedPolicyCategory `json:"policy_category,omitzero"`
	paramObj
}

func (r BetaComplianceReviewVerdictParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaComplianceReviewVerdictParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaComplianceReviewVerdictParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaComplianceReviewViolatedPolicyCategory string

const (
	BetaComplianceReviewViolatedPolicyCategoryCyber               BetaComplianceReviewViolatedPolicyCategory = "cyber"
	BetaComplianceReviewViolatedPolicyCategoryBio                 BetaComplianceReviewViolatedPolicyCategory = "bio"
	BetaComplianceReviewViolatedPolicyCategoryChem                BetaComplianceReviewViolatedPolicyCategory = "chem"
	BetaComplianceReviewViolatedPolicyCategoryRadnuc              BetaComplianceReviewViolatedPolicyCategory = "radnuc"
	BetaComplianceReviewViolatedPolicyCategoryCgw                 BetaComplianceReviewViolatedPolicyCategory = "cgw"
	BetaComplianceReviewViolatedPolicyCategoryMgw                 BetaComplianceReviewViolatedPolicyCategory = "mgw"
	BetaComplianceReviewViolatedPolicyCategoryIo                  BetaComplianceReviewViolatedPolicyCategory = "io"
	BetaComplianceReviewViolatedPolicyCategoryFraud               BetaComplianceReviewViolatedPolicyCategory = "fraud"
	BetaComplianceReviewViolatedPolicyCategoryViolentExtremism    BetaComplianceReviewViolatedPolicyCategory = "violent_extremism"
	BetaComplianceReviewViolatedPolicyCategorySurveillance        BetaComplianceReviewViolatedPolicyCategory = "surveillance"
	BetaComplianceReviewViolatedPolicyCategoryHumanExploitation   BetaComplianceReviewViolatedPolicyCategory = "human_exploitation"
	BetaComplianceReviewViolatedPolicyCategorySelfHarm            BetaComplianceReviewViolatedPolicyCategory = "self_harm"
	BetaComplianceReviewViolatedPolicyCategoryAdultContent        BetaComplianceReviewViolatedPolicyCategory = "adult_content"
	BetaComplianceReviewViolatedPolicyCategoryReasoningExtraction BetaComplianceReviewViolatedPolicyCategory = "reasoning_extraction"
	BetaComplianceReviewViolatedPolicyCategoryFrontierLLM         BetaComplianceReviewViolatedPolicyCategory = "frontier_llm"
	BetaComplianceReviewViolatedPolicyCategoryCompromisedKey      BetaComplianceReviewViolatedPolicyCategory = "compromised_key"
	BetaComplianceReviewViolatedPolicyCategoryOther               BetaComplianceReviewViolatedPolicyCategory = "other"
)

type BetaComplianceReviewCaseGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaComplianceReviewCaseListParams struct {
	// Return only cases whose `created_at` is after this time.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return only cases whose `created_at` is at or after this time.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only cases whose `created_at` is before this time.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return only cases whose `created_at` is at or before this time.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Number of cases to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`. Leave it out to get the
	// first page. A cursor works only with the filters and sort order of the request
	// that returned it.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only cases whose `updated_at` is after this time.
	UpdatedAtGt param.Opt[time.Time] `query:"updated_at[gt],omitzero" format:"date-time" json:"-"`
	// Return cases whose `updated_at` is at or after this time. Cases slightly older
	// than this time are also returned, so that a sync misses no change.
	UpdatedAtGte param.Opt[time.Time] `query:"updated_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only cases whose `updated_at` is before this time.
	UpdatedAtLt param.Opt[time.Time] `query:"updated_at[lt],omitzero" format:"date-time" json:"-"`
	// Return only cases whose `updated_at` is at or before this time.
	UpdatedAtLte param.Opt[time.Time] `query:"updated_at[lte],omitzero" format:"date-time" json:"-"`
	// Return only cases where a message's `user_profile_id` is this ID, including
	// cases that also have messages of other user profiles. A case stops matching once
	// its messages' `user_profile_id` is null. An ID that matches no reachable case
	// returns no cases. A request that sends the `anthropic-user-profile-id` header
	// without this parameter returns a 400 error.
	UserProfileID param.Opt[string] `query:"user_profile_id,omitzero" json:"-"`
	// Sort direction, by `order_by` and then by `id`. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaComplianceReviewOrder `query:"order,omitzero" json:"-"`
	// The field to sort by. Defaults to `created_at`.
	//
	// Any of "created_at", "updated_at".
	OrderBy BetaComplianceReviewCaseOrderBy `query:"order_by,omitzero" json:"-"`
	// Return only cases with this status. Leave it out to get both open and closed
	// cases.
	//
	// Any of "open", "closed".
	Status BetaComplianceReviewCaseStatusFilter `query:"status,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaComplianceReviewCaseListParams]'s query parameters as
// `url.Values`.
func (r BetaComplianceReviewCaseListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaComplianceReviewCaseRespondParams struct {
	// Your conclusion on the case as a whole. Anthropic doesn't check it against
	// `verdicts`.
	//
	// Any of "violation", "no_violation".
	Determination BetaComplianceReviewDetermination `json:"determination,omitzero" api:"required"`
	// The policy categories that you found were violated. Send one to four with
	// `violation`, and none with `no_violation`. They need not match the case's own.
	// Each value is one of `cyber`, `bio`, `chem`, `radnuc` (radiological or nuclear
	// weapons), `cgw` (conventional weapons), `mgw` (military-grade weapons systems),
	// `io` (influence operations), `fraud`, `violent_extremism`, `surveillance`,
	// `human_exploitation`, `self_harm`, `adult_content`, `reasoning_extraction`,
	// `frontier_llm` (use of the model to accelerate frontier AI research and
	// development in restricted ways), `compromised_key` or `other`, and can appear
	// once.
	PolicyCategories []BetaComplianceReviewViolatedPolicyCategory `json:"policy_categories,omitzero"`
	// Your conclusions on individual messages of the case. They need not cover every
	// message, and each `item_id` can appear once.
	Verdicts []BetaComplianceReviewVerdictParam `json:"verdicts,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaComplianceReviewCaseRespondParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaComplianceReviewCaseRespondParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaComplianceReviewCaseRespondParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

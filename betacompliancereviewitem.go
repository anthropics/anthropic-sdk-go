package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaComplianceReviewItemService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaComplianceReviewItemService] method instead.
type BetaComplianceReviewItemService struct {
	Options []option.RequestOption
}

// NewBetaComplianceReviewItemService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaComplianceReviewItemService(opts ...option.RequestOption) (r BetaComplianceReviewItemService) {
	r = BetaComplianceReviewItemService{}
	r.Options = opts
	return
}

// Read one message that Anthropic flagged for your organization to review, with
// its request and response.
//
// A message that the API key can't read, or that isn't flagged, returns the same
// 404 as one that doesn't exist. For a message recorded with a user profile, send
// that profile's ID in the `anthropic-user-profile-id` header. A flagged message
// can also return 404 if it belongs to another organization than the key's own, or
// to your end customer's user profile or Workspace.
//
// If your encryption key or storage bucket refuses the read, the API returns a 400
// with `cmek_key_disabled`, `storage_bucket_unreadable` or
// `customer_storage_refusal_unattributed` in `error.details.error_code`.
//
// The API key must have the `read:compliance_review_items`,
// `read:compliance_user_data` or `read:org_audit` scope. A key of an organization
// with no parent organization must have `read:compliance_review_items`.
func (r *BetaComplianceReviewItemService) Get(ctx context.Context, itemID string, query BetaComplianceReviewItemGetParams, opts ...option.RequestOption) (res *BetaComplianceReviewItem, err error) {
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
	if !param.IsOmitted(query.UserProfileID) {
		opts = append(opts, option.WithHeader("anthropic-user-profile-id", fmt.Sprintf("%v", query.UserProfileID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	if itemID == "" {
		err = errors.New("missing required item_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/compliance/review/items/%s", requestconfig.EscapePathSegment(itemID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.compliance.review.items.retrieve", "/v1/compliance/review/items/{item_id}"), requestconfig.WithSpanAttribute("anthropic.compliance.review.item.id", itemID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// BetaComplianceReviewActorUnion contains all possible properties and values from
// [BetaComplianceReviewUserActor], [BetaComplianceReviewAPIActor],
// [BetaComplianceReviewServiceAccountActor].
//
// Use the [BetaComplianceReviewActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewActorUnion struct {
	// This field is from variant [BetaComplianceReviewUserActor].
	EmailAddress string `json:"email_address"`
	// Any of "user_actor", "api_actor", "service_account_actor".
	Type string `json:"type"`
	// This field is from variant [BetaComplianceReviewUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaComplianceReviewAPIActor].
	APIKeyID string `json:"api_key_id"`
	// This field is from variant [BetaComplianceReviewServiceAccountActor].
	ServiceAccountID string `json:"service_account_id"`
	JSON             struct {
		EmailAddress     respjson.Field
		Type             respjson.Field
		UserID           respjson.Field
		APIKeyID         respjson.Field
		ServiceAccountID respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaComplianceReviewActor is implemented by each variant of
// [BetaComplianceReviewActorUnion] to add type safety for the return type of
// [BetaComplianceReviewActorUnion.AsAny]
type anyBetaComplianceReviewActor interface {
	implBetaComplianceReviewActorUnion()
}

func (BetaComplianceReviewUserActor) implBetaComplianceReviewActorUnion()           {}
func (BetaComplianceReviewAPIActor) implBetaComplianceReviewActorUnion()            {}
func (BetaComplianceReviewServiceAccountActor) implBetaComplianceReviewActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewActorUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewUserActor:
//	case anthropic.BetaComplianceReviewAPIActor:
//	case anthropic.BetaComplianceReviewServiceAccountActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewActorUnion) AsAny() anyBetaComplianceReviewActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "api_actor":
		return u.AsAPIActor()
	case "service_account_actor":
		return u.AsServiceAccountActor()
	}
	return nil
}

func (u BetaComplianceReviewActorUnion) AsUserActor() (v BetaComplianceReviewUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaComplianceReviewActorUnion) AsAPIActor() (v BetaComplianceReviewAPIActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaComplianceReviewActorUnion) AsServiceAccountActor() (v BetaComplianceReviewServiceAccountActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An API key.
type BetaComplianceReviewAPIActor struct {
	// The ID of the API key (`apikey_...`).
	APIKeyID string            `json:"api_key_id" api:"required"`
	Type     constant.APIActor `json:"type" default:"api_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APIKeyID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewAPIActor) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewAPIActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No copy of the message was stored.
type BetaComplianceReviewContentNotCaptured struct {
	Type constant.NotCaptured `json:"type" default:"not_captured"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewContentNotCaptured) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewContentNotCaptured) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaComplianceReviewContentUnavailableReasonUnion contains all possible
// properties and values from [BetaComplianceReviewContentNotCaptured],
// [BetaComplianceReviewStoredCopyUnavailable].
//
// Use the [BetaComplianceReviewContentUnavailableReasonUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewContentUnavailableReasonUnion struct {
	// Any of "not_captured", "stored_copy_unavailable".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaComplianceReviewContentUnavailableReason is implemented by each variant
// of [BetaComplianceReviewContentUnavailableReasonUnion] to add type safety for
// the return type of [BetaComplianceReviewContentUnavailableReasonUnion.AsAny]
type anyBetaComplianceReviewContentUnavailableReason interface {
	implBetaComplianceReviewContentUnavailableReasonUnion()
}

func (BetaComplianceReviewContentNotCaptured) implBetaComplianceReviewContentUnavailableReasonUnion() {
}
func (BetaComplianceReviewStoredCopyUnavailable) implBetaComplianceReviewContentUnavailableReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewContentUnavailableReasonUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewContentNotCaptured:
//	case anthropic.BetaComplianceReviewStoredCopyUnavailable:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewContentUnavailableReasonUnion) AsAny() anyBetaComplianceReviewContentUnavailableReason {
	switch u.Type {
	case "not_captured":
		return u.AsNotCaptured()
	case "stored_copy_unavailable":
		return u.AsStoredCopyUnavailable()
	}
	return nil
}

func (u BetaComplianceReviewContentUnavailableReasonUnion) AsNotCaptured() (v BetaComplianceReviewContentNotCaptured) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaComplianceReviewContentUnavailableReasonUnion) AsStoredCopyUnavailable() (v BetaComplianceReviewStoredCopyUnavailable) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewContentUnavailableReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewContentUnavailableReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A request to create a message, and its response.
type BetaComplianceReviewCreateMessageContent struct {
	// The `POST /v1/messages` request body, without thinking blocks or unpublished
	// fields. For a user's message from an Anthropic app, such as Claude Code,
	// `system` is a placeholder.
	Request map[string]any `json:"request" api:"required"`
	// The model's `Message`, without thinking blocks or unpublished fields. Null when
	// no response was stored, for example because the client disconnected.
	Response BetaMessage            `json:"response" api:"required"`
	Type     constant.CreateMessage `json:"type" default:"create_message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Request     respjson.Field
		Response    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewCreateMessageContent) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewCreateMessageContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One stored item that the organization it belongs to can read for review.
type BetaComplianceReviewItem struct {
	// The ID of the item. For a message, its own ID, such as `msg_...`.
	ID string `json:"id" api:"required"`
	// For a message, who sent it. Null when no sender was recorded. New types can
	// appear, so accept a type that you don't know.
	Actor BetaComplianceReviewActorUnion `json:"actor" api:"required"`
	// Read the item, and save any copy you need, by this time. The item can stop being
	// readable sooner.
	ContentExpiresAt time.Time `json:"content_expires_at" api:"required" format:"date-time"`
	// For a message, when Anthropic received the call.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What is particular to this kind of item. New types can appear, so accept a type
	// that you don't know.
	Details BetaComplianceReviewItemDetailsUnion `json:"details" api:"required"`
	// The UUID of the organization that the item belongs to. For a message, the
	// organization that sent it.
	OrganizationID string `json:"organization_id" api:"required"`
	// Why the item is readable. New types can appear, so accept a type that you don't
	// know.
	ReadReason BetaComplianceReviewReadReasonUnion `json:"read_reason" api:"required"`
	// The object type, always `compliance_review_item`.
	Type constant.ComplianceReviewItem `json:"type" default:"compliance_review_item"`
	// The ID of the user profile that the item was attributed to (`uprof_...`). Null
	// when none was recorded.
	UserProfileID string `json:"user_profile_id" api:"required"`
	// The ID of the workspace (`wrkspc_...`). Null when the item isn't attributed to a
	// workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		Actor            respjson.Field
		ContentExpiresAt respjson.Field
		CreatedAt        respjson.Field
		Details          respjson.Field
		OrganizationID   respjson.Field
		ReadReason       respjson.Field
		Type             respjson.Field
		UserProfileID    respjson.Field
		WorkspaceID      respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewItem) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaComplianceReviewItemDetailsUnion contains all possible properties and values
// from [BetaComplianceReviewMessageDetails].
//
// Use the [BetaComplianceReviewItemDetailsUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewItemDetailsUnion struct {
	// This field is from variant [BetaComplianceReviewMessageDetails].
	Content BetaComplianceReviewMessageContentUnion `json:"content"`
	// This field is from variant [BetaComplianceReviewMessageDetails].
	ContentUnavailableReason BetaComplianceReviewContentUnavailableReasonUnion `json:"content_unavailable_reason"`
	// This field is from variant [BetaComplianceReviewMessageDetails].
	Model string `json:"model"`
	// Any of "message".
	Type string `json:"type"`
	JSON struct {
		Content                  respjson.Field
		ContentUnavailableReason respjson.Field
		Model                    respjson.Field
		Type                     respjson.Field
		raw                      string
	} `json:"-"`
}

// anyBetaComplianceReviewItemDetails is implemented by each variant of
// [BetaComplianceReviewItemDetailsUnion] to add type safety for the return type of
// [BetaComplianceReviewItemDetailsUnion.AsAny]
type anyBetaComplianceReviewItemDetails interface {
	implBetaComplianceReviewItemDetailsUnion()
}

func (BetaComplianceReviewMessageDetails) implBetaComplianceReviewItemDetailsUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewItemDetailsUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewMessageDetails:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewItemDetailsUnion) AsAny() anyBetaComplianceReviewItemDetails {
	switch u.Type {
	case "message":
		return u.AsMessage()
	}
	return nil
}

func (u BetaComplianceReviewItemDetailsUnion) AsMessage() (v BetaComplianceReviewMessageDetails) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewItemDetailsUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewItemDetailsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaComplianceReviewMessageContentUnion contains all possible properties and
// values from [BetaComplianceReviewCreateMessageContent].
//
// Use the [BetaComplianceReviewMessageContentUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewMessageContentUnion struct {
	// This field is from variant [BetaComplianceReviewCreateMessageContent].
	Request map[string]any `json:"request"`
	// This field is from variant [BetaComplianceReviewCreateMessageContent].
	Response BetaMessage `json:"response"`
	// Any of "create_message".
	Type string `json:"type"`
	JSON struct {
		Request  respjson.Field
		Response respjson.Field
		Type     respjson.Field
		raw      string
	} `json:"-"`
}

// anyBetaComplianceReviewMessageContent is implemented by each variant of
// [BetaComplianceReviewMessageContentUnion] to add type safety for the return type
// of [BetaComplianceReviewMessageContentUnion.AsAny]
type anyBetaComplianceReviewMessageContent interface {
	implBetaComplianceReviewMessageContentUnion()
}

func (BetaComplianceReviewCreateMessageContent) implBetaComplianceReviewMessageContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewMessageContentUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewCreateMessageContent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewMessageContentUnion) AsAny() anyBetaComplianceReviewMessageContent {
	switch u.Type {
	case "create_message":
		return u.AsCreateMessage()
	}
	return nil
}

func (u BetaComplianceReviewMessageContentUnion) AsCreateMessage() (v BetaComplianceReviewCreateMessageContent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewMessageContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewMessageContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A message sent to Claude.
type BetaComplianceReviewMessageDetails struct {
	// The stored request and response. Null when `content_unavailable_reason` is set.
	// New types can appear, so accept a type that you don't know.
	Content BetaComplianceReviewMessageContentUnion `json:"content" api:"required"`
	// Why `content` is null. Null when `content` is set. New types can appear, so
	// accept a type that you don't know.
	ContentUnavailableReason BetaComplianceReviewContentUnavailableReasonUnion `json:"content_unavailable_reason" api:"required"`
	// The model that the call named. Null when none was recorded.
	Model string           `json:"model" api:"required"`
	Type  constant.Message `json:"type" default:"message"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content                  respjson.Field
		ContentUnavailableReason respjson.Field
		Model                    respjson.Field
		Type                     respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewMessageDetails) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewMessageDetails) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaComplianceReviewReadReasonUnion contains all possible properties and values
// from [BetaComplianceReviewSafeguardsReason].
//
// Use the [BetaComplianceReviewReadReasonUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceReviewReadReasonUnion struct {
	// Any of "safeguards_review".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaComplianceReviewReadReason is implemented by each variant of
// [BetaComplianceReviewReadReasonUnion] to add type safety for the return type of
// [BetaComplianceReviewReadReasonUnion.AsAny]
type anyBetaComplianceReviewReadReason interface {
	implBetaComplianceReviewReadReasonUnion()
}

func (BetaComplianceReviewSafeguardsReason) implBetaComplianceReviewReadReasonUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceReviewReadReasonUnion.AsAny().(type) {
//	case anthropic.BetaComplianceReviewSafeguardsReason:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceReviewReadReasonUnion) AsAny() anyBetaComplianceReviewReadReason {
	switch u.Type {
	case "safeguards_review":
		return u.AsSafeguardsReview()
	}
	return nil
}

func (u BetaComplianceReviewReadReasonUnion) AsSafeguardsReview() (v BetaComplianceReviewSafeguardsReason) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceReviewReadReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceReviewReadReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic flagged the item for your organization to review.
type BetaComplianceReviewSafeguardsReason struct {
	Type constant.SafeguardsReview `json:"type" default:"safeguards_review"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewSafeguardsReason) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewSafeguardsReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A service account.
type BetaComplianceReviewServiceAccountActor struct {
	// The ID of the service account (`svac_...`).
	ServiceAccountID string                       `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccountActor `json:"type" default:"service_account_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewServiceAccountActor) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewServiceAccountActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A copy of the message was stored and can't be returned.
type BetaComplianceReviewStoredCopyUnavailable struct {
	Type constant.StoredCopyUnavailable `json:"type" default:"stored_copy_unavailable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewStoredCopyUnavailable) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewStoredCopyUnavailable) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A signed-in user.
type BetaComplianceReviewUserActor struct {
	// The user's email address. Set only for a current member of an organization that
	// the API key can read.
	EmailAddress string             `json:"email_address" api:"required"`
	Type         constant.UserActor `json:"type" default:"user_actor"`
	// The ID of the user (`user_...`).
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EmailAddress respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaComplianceReviewItemGetParams struct {
	// The ID of the user profile that the message was recorded with (`uprof_...`).
	// Send it for a message recorded with a user profile, and leave it out for a
	// message recorded with none.
	UserProfileID param.Opt[string] `header:"anthropic-user-profile-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

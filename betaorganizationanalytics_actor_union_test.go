package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestAnalyticsUsersItemActorUnion(t *testing.T) {
	raw := `{"type":"user_actor","user_id":"user_x","name":"Jane","email_address":"jane@example.com","email":"jane@example.com","deleted":false}`

	var usage anthropic.BetaAnalyticsUsageUsersItemActorUnion
	if err := json.Unmarshal([]byte(raw), &usage); err != nil {
		t.Fatal(err)
	}
	if usage.Type != "user_actor" || usage.EmailAddress != "jane@example.com" {
		t.Fatalf("flattened fields: %+v", usage)
	}
	switch v := usage.AsAny().(type) {
	case anthropic.BetaAnalyticsUserActor:
		if v.UserID != "user_x" || v.EmailAddress != "jane@example.com" {
			t.Fatalf("variant fields: %+v", v)
		}
	default:
		t.Fatalf("AsAny returned %T", v)
	}
	if usage.RawJSON() != raw {
		t.Fatalf("RawJSON changed: %s", usage.RawJSON())
	}

	var cost anthropic.BetaAnalyticsCostUsersItemActorUnion
	if err := json.Unmarshal([]byte(raw), &cost); err != nil {
		t.Fatal(err)
	}
	if _, ok := cost.AsAny().(anthropic.BetaAnalyticsUserActor); !ok {
		t.Fatalf("cost AsAny returned %T", cost.AsAny())
	}

	var unknown anthropic.BetaAnalyticsUsageUsersItemActorUnion
	if err := json.Unmarshal([]byte(`{"type":"service_account_actor","service_account_id":"svac_x"}`), &unknown); err != nil {
		t.Fatal(err)
	}
	if unknown.AsAny() != nil {
		t.Fatalf("an unknown variant should yield nil, got %T", unknown.AsAny())
	}
}

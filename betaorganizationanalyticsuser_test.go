package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBetaOrganizationAnalyticsUserListWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Organization.Analytics.Users.List(context.TODO(), anthropic.BetaOrganizationAnalyticsUserListParams{
		Date:         anthropic.Time(time.Now()),
		EndingDate:   anthropic.Time(time.Now()),
		Filter:       []string{"string"},
		GroupBy:      []string{"rbac_group_id"},
		Limit:        anthropic.Int(1),
		Order:        anthropic.BetaOrganizationAnalyticsUserListParamsOrderAsc,
		OrderBy:      anthropic.String("order_by"),
		Page:         anthropic.String("page"),
		StartingDate: anthropic.Time(time.Now()),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

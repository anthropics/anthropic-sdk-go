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

func TestBetaOrganizationAnalyticsConnectorListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Analytics.Connectors.List(context.TODO(), anthropic.BetaOrganizationAnalyticsConnectorListParams{
		Date:         anthropic.Time(time.Now()),
		EndingDate:   anthropic.Time(time.Now()),
		Filter:       []string{"string"},
		GroupBy:      []string{"product"},
		Limit:        anthropic.Int(1),
		Order:        anthropic.BetaOrganizationAnalyticsConnectorListParamsOrderAsc,
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

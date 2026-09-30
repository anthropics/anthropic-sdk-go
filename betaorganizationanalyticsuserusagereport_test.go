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

func TestBetaOrganizationAnalyticsUserUsageReportListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Analytics.UserUsageReport.List(context.TODO(), anthropic.BetaOrganizationAnalyticsUserUsageReportListParams{
		StartingAt:          time.Now(),
		BucketWidth:         anthropic.BetaOrganizationAnalyticsUserUsageReportListParamsBucketWidthDay,
		ClaudeTagCategories: []anthropic.BetaAnalyticsClaudeTagCategory{anthropic.BetaAnalyticsClaudeTagCategoryEngaged},
		ClaudeTagUserIDs:    []string{"U0123ABCDEF"},
		ContextWindows:      []anthropic.BetaAnalyticsContextWindow{anthropic.BetaAnalyticsContextWindowFrom0To200k},
		EndingAt:            anthropic.Time(time.Now()),
		ExcludeDeletedUsers: anthropic.Bool(true),
		GroupBy:             []string{"claude_tag_category"},
		InferenceGeos:       []anthropic.BetaAnalyticsInferenceGeoFilter{anthropic.BetaAnalyticsInferenceGeoFilterGlobal},
		Limit:               anthropic.Int(1),
		Models:              []string{"string"},
		Order:               anthropic.BetaOrganizationAnalyticsUserUsageReportListParamsOrderAsc,
		OrderBy:             anthropic.BetaOrganizationAnalyticsUserUsageReportListParamsOrderByOutputTokens,
		Page:                anthropic.String("page"),
		Products:            []anthropic.BetaAnalyticsProductFilter{anthropic.BetaAnalyticsProductFilterChat},
		RBACGroupIDs:        []string{"rbac_group_012rppKaSVsmTo6NqRDXQXNF"},
		SlackChannelIDs:     []string{"C0123ABCDEF"},
		Speeds:              []string{"fast"},
		UserIDs:             []string{"user_01AbCdEfGhIjKlMnOpQrSt"},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

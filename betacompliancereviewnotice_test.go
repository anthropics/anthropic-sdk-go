package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestBetaComplianceReviewNoticeListWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Compliance.Review.Notices.List(context.TODO(), anthropic.BetaComplianceReviewNoticeListParams{
		CaseID:       anthropic.String("case_id"),
		CreatedAtGt:  anthropic.Time(time.Now()),
		CreatedAtGte: anthropic.Time(time.Now()),
		CreatedAtLt:  anthropic.Time(time.Now()),
		CreatedAtLte: anthropic.Time(time.Now()),
		Limit:        anthropic.Int(0),
		Order:        anthropic.BetaComplianceReviewOrderAsc,
		Page:         anthropic.String("page"),
		Betas:        []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

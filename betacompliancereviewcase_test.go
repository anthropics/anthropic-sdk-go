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

func TestBetaComplianceReviewCaseGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Compliance.Review.Cases.Get(
		context.TODO(),
		"case_id",
		anthropic.BetaComplianceReviewCaseGetParams{
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaComplianceReviewCaseListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Compliance.Review.Cases.List(context.TODO(), anthropic.BetaComplianceReviewCaseListParams{
		CreatedAtGt:   anthropic.Time(time.Now()),
		CreatedAtGte:  anthropic.Time(time.Now()),
		CreatedAtLt:   anthropic.Time(time.Now()),
		CreatedAtLte:  anthropic.Time(time.Now()),
		Limit:         anthropic.Int(0),
		Order:         anthropic.BetaComplianceReviewOrderAsc,
		OrderBy:       anthropic.BetaComplianceReviewCaseOrderByCreatedAt,
		Page:          anthropic.String("page"),
		Status:        anthropic.BetaComplianceReviewCaseStatusFilterOpen,
		UpdatedAtGt:   anthropic.Time(time.Now()),
		UpdatedAtGte:  anthropic.Time(time.Now()),
		UpdatedAtLt:   anthropic.Time(time.Now()),
		UpdatedAtLte:  anthropic.Time(time.Now()),
		UserProfileID: anthropic.String("user_profile_id"),
		Betas:         []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaComplianceReviewCaseRespondWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Compliance.Review.Cases.Respond(
		context.TODO(),
		"case_id",
		anthropic.BetaComplianceReviewCaseRespondParams{
			Determination:    anthropic.BetaComplianceReviewDeterminationNoViolation,
			PolicyCategories: []anthropic.BetaComplianceReviewViolatedPolicyCategory{anthropic.BetaComplianceReviewViolatedPolicyCategoryCyber},
			Verdicts: []anthropic.BetaComplianceReviewVerdictParam{{
				Determination:  anthropic.BetaComplianceReviewDeterminationViolation,
				ItemID:         "item_id",
				PolicyCategory: anthropic.BetaComplianceReviewViolatedPolicyCategoryCyber,
			}},
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

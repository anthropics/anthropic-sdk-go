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

func TestBetaUserProfileNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.New(context.TODO(), anthropic.BetaUserProfileNewParams{
		AccessType: anthropic.BetaUserProfileNewParamsAccessTypeApplication,
		ExternalID: anthropic.String("x"),
		ExternalUserDetails: anthropic.BetaUserProfileExternalUserDetailsParams{
			AccountStatus: anthropic.BetaUserProfileExternalUserDetailsParamsAccountStatusActive,
			Country:       anthropic.String("US"),
			EmailHash:     anthropic.String("x"),
			EntityType:    anthropic.BetaUserProfileExternalUserDetailsParamsEntityTypeBusiness,
			NameHash:      anthropic.String("x"),
			OnboardedAt:   anthropic.Time(time.Now()),
			ReferenceID:   anthropic.String("user_12345"),
		},
		ExternalUserOnboardedAt: anthropic.Time(time.Now()),
		Metadata:                map[string]string{},
		Name:                    anthropic.String("Example User"),
		Betas:                   []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		WorkspaceID:             anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaUserProfileGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.Get(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileGetParams{
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

func TestBetaUserProfileUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.Update(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileUpdateParams{
			AccessType: anthropic.BetaUserProfileUpdateParamsAccessTypeApplication,
			ExternalID: anthropic.String("x"),
			ExternalUserDetails: anthropic.BetaUserProfileExternalUserDetailsParams{
				AccountStatus: anthropic.BetaUserProfileExternalUserDetailsParamsAccountStatusSuspended,
				Country:       anthropic.String("country"),
				EmailHash:     anthropic.String("x"),
				EntityType:    anthropic.BetaUserProfileExternalUserDetailsParamsEntityTypeIndividual,
				NameHash:      anthropic.String("x"),
				OnboardedAt:   anthropic.Time(time.Now()),
				ReferenceID:   anthropic.String("x"),
			},
			ExternalUserOnboardedAt: anthropic.Time(time.Now()),
			Metadata: map[string]string{
				"foo": "string",
			},
			Name:        anthropic.String("x"),
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

func TestBetaUserProfileListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.List(context.TODO(), anthropic.BetaUserProfileListParams{
		Limit:       anthropic.Int(0),
		Order:       anthropic.BetaUserProfileListParamsOrderAsc,
		OrderBy:     anthropic.BetaUserProfileListParamsOrderByCreatedAt,
		Page:        anthropic.String("page"),
		Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaUserProfileDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.Delete(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileDeleteParams{
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

func TestBetaUserProfileNewEnrollmentURLWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.NewEnrollmentURL(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileNewEnrollmentURLParams{
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

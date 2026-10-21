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

func TestUserProfileNewWithOptionalParams(t *testing.T) {
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
	_, err := client.UserProfiles.New(context.TODO(), anthropic.UserProfileNewParams{
		AccessType: anthropic.UserProfileAccessTypeApplication,
		ExternalUserDetails: anthropic.UserProfileExternalUserDetailsParams{
			AccountStatus: anthropic.UserProfileExternalUserAccountStatusActive,
			Country:       anthropic.String("US"),
			EmailHash:     anthropic.String("x"),
			EntityType:    anthropic.UserProfileExternalUserEntityTypeBusiness,
			NameHash:      anthropic.String("x"),
			OnboardedAt:   anthropic.Time(time.Now()),
			ReferenceID:   anthropic.String("user_12345"),
		},
		Metadata:    map[string]string{},
		Name:        anthropic.String("Example User"),
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

func TestUserProfileGetWithOptionalParams(t *testing.T) {
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
	_, err := client.UserProfiles.Get(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.UserProfileGetParams{
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

func TestUserProfileUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.UserProfiles.Update(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.UserProfileUpdateParams{
			AccessType: anthropic.UserProfileAccessTypeApplication,
			ExternalUserDetails: anthropic.UserProfileExternalUserDetailsParams{
				AccountStatus: anthropic.UserProfileExternalUserAccountStatusSuspended,
				Country:       anthropic.String("country"),
				EmailHash:     anthropic.String("x"),
				EntityType:    anthropic.UserProfileExternalUserEntityTypeIndividual,
				NameHash:      anthropic.String("x"),
				OnboardedAt:   anthropic.Time(time.Now()),
				ReferenceID:   anthropic.String("x"),
			},
			Metadata: map[string]string{
				"foo": "string",
			},
			Name:        anthropic.String("x"),
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

func TestUserProfileListWithOptionalParams(t *testing.T) {
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
	_, err := client.UserProfiles.List(context.TODO(), anthropic.UserProfileListParams{
		Limit:       anthropic.Int(0),
		Order:       anthropic.UserProfileListOrderAsc,
		OrderBy:     anthropic.UserProfileListOrderByCreatedAt,
		Page:        anthropic.String("page"),
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

func TestUserProfileDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.UserProfiles.Delete(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.UserProfileDeleteParams{
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

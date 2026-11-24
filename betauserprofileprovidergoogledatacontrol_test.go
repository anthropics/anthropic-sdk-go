package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestBetaUserProfileProviderGoogleDataControlGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.Providers.Google.DataControls.Get(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileProviderGoogleDataControlGetParams{
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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

func TestBetaUserProfileProviderGoogleDataControlUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.UserProfiles.Providers.Google.DataControls.Update(
		context.TODO(),
		"uprof_011CZkZCu8hGbp5mYRQgUmz9",
		anthropic.BetaUserProfileProviderGoogleDataControlUpdateParams{
			CMEK: anthropic.BetaGoogleManagedCMEKControlParamsUnion{
				OfEnabled: &anthropic.BetaGoogleManagedCMEKEnabledParams{
					ExternalKeyIDs: map[string]string{
						"us": "ekey_01SDCCSbTxrXDpWc1phhtcfK",
					},
				},
			},
			Storage: anthropic.BetaGoogleManagedStorageControlParamsUnion{
				OfDisabled: &anthropic.BetaGoogleManagedStorageDisabledParams{},
			},
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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

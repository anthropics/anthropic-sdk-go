package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

func TestBetaOrganizationWorkspaceNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.New(context.TODO(), anthropic.BetaOrganizationWorkspaceNewParams{
		Name: "Workspace Name",
		DataResidency: anthropic.BetaDataResidencyCreateConfigParam{
			AllowedInferenceGeos: anthropic.BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam{
				OfUnrestricted: constant.ValueOf[constant.Unrestricted](),
			},
			DefaultInferenceGeo: anthropic.BetaDataResidencyCreateConfigDefaultInferenceGeoGlobal,
			WorkspaceGeo:        anthropic.BetaDataResidencyCreateConfigWorkspaceGeoUs,
		},
		DisplayColor:  anthropic.String("#6C5BB9"),
		ExternalKeyID: anthropic.String("ekey_01SDCCSbTxrXDpWc1phhtcfK"),
		InferenceDataPrivacy: anthropic.BetaInferenceDataPrivacyCreateParam{
			HumanReview: anthropic.BetaAnthropicHumanReviewParam{
				Type: anthropic.BetaAnthropicHumanReviewParamTypeDisabled,
			},
			SafeguardsStorage: anthropic.BetaStorageMemberParamUnion{
				OfInternal: &anthropic.BetaInternalStorageParam{},
			},
			TranscriptsStorage: anthropic.BetaStorageMemberParamUnion{
				OfInternal: &anthropic.BetaInternalStorageParam{},
			},
		},
		Tags: map[string]string{
			"env":  "prod",
			"team": "platform",
		},
		UserProfileID: anthropic.String("uprof_01SDCCSbTxrXDpWc1phhtcfK"),
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

func TestBetaOrganizationWorkspaceGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.Get(
		context.TODO(),
		"workspace_id",
		anthropic.BetaOrganizationWorkspaceGetParams{
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

func TestBetaOrganizationWorkspaceUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.Update(
		context.TODO(),
		"workspace_id",
		anthropic.BetaOrganizationWorkspaceUpdateParams{
			DataResidency: anthropic.BetaDataResidencyUpdateConfigParam{
				AllowedInferenceGeos: anthropic.BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam{
					OfUnrestricted: constant.ValueOf[constant.Unrestricted](),
				},
				DefaultInferenceGeo: anthropic.BetaDataResidencyUpdateConfigDefaultInferenceGeoGlobal,
			},
			DisplayColor:  anthropic.String("#6C5BB9"),
			ExternalKeyID: anthropic.String("ekey_01SDCCSbTxrXDpWc1phhtcfK"),
			InferenceDataPrivacy: anthropic.BetaInferenceDataPrivacyUpdateParam{
				HumanReview: anthropic.BetaAnthropicHumanReviewParam{
					Type: anthropic.BetaAnthropicHumanReviewParamTypeDisabled,
				},
				SafeguardsStorage: anthropic.BetaStorageMemberParamUnion{
					OfInternal: &anthropic.BetaInternalStorageParam{},
				},
				TranscriptsStorage: anthropic.BetaStorageMemberParamUnion{
					OfInternal: &anthropic.BetaInternalStorageParam{},
				},
			},
			Name: anthropic.String("Workspace Name"),
			Tags: map[string]string{
				"env":  "prod",
				"team": "platform",
			},
			UserProfileID: anthropic.String("uprof_01SDCCSbTxrXDpWc1phhtcfK"),
			Betas:         []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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

func TestBetaOrganizationWorkspaceListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.List(context.TODO(), anthropic.BetaOrganizationWorkspaceListParams{
		AfterID:         anthropic.String("after_id"),
		BeforeID:        anthropic.String("before_id"),
		IDs:             []string{"string"},
		IncludeArchived: anthropic.Bool(true),
		IncludeDefault:  anthropic.Bool(true),
		Limit:           anthropic.Int(1),
		Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationWorkspaceArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.Archive(
		context.TODO(),
		"workspace_id",
		anthropic.BetaOrganizationWorkspaceArchiveParams{
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

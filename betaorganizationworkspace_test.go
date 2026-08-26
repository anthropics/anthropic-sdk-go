// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
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
	)
	_, err := client.Beta.Organization.Workspaces.New(context.TODO(), anthropic.BetaOrganizationWorkspaceNewParams{
		Name: "x",
		DataResidency: anthropic.BetaDataResidencyCreateConfigParam{
			AllowedInferenceGeos: anthropic.BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam{
				OfUnrestricted: constant.ValueOf[constant.Unrestricted](),
			},
			DefaultInferenceGeo: anthropic.BetaDataResidencyCreateConfigDefaultInferenceGeoGlobal,
			WorkspaceGeo:        anthropic.BetaDataResidencyCreateConfigWorkspaceGeoUs,
		},
		DisplayColor:  anthropic.String("#6C5BB9"),
		ExternalKeyID: anthropic.String("ekey_01SDCCSbTxrXDpWc1phhtcfK"),
		Tags: map[string]string{
			"env":  "prod",
			"team": "platform",
		},
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationWorkspaceGet(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.Get(context.TODO(), "workspace_id")
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
			Name:          anthropic.String("x"),
			Tags: map[string]string{
				"env":  "prod",
				"team": "platform",
			},
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
	)
	_, err := client.Beta.Organization.Workspaces.List(context.TODO(), anthropic.BetaOrganizationWorkspaceListParams{
		AfterID:         anthropic.String("after_id"),
		BeforeID:        anthropic.String("before_id"),
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(1),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationWorkspaceArchive(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.Archive(context.TODO(), "workspace_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

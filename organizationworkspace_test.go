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

func TestOrganizationWorkspaceNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Workspaces.New(context.TODO(), anthropic.OrganizationWorkspaceNewParams{
		Name: "x",
		DataResidency: anthropic.DataResidencyCreateConfigParam{
			AllowedInferenceGeos: anthropic.DataResidencyCreateConfigAllowedInferenceGeosUnionParam{
				OfUnrestricted: constant.ValueOf[constant.Unrestricted](),
			},
			DefaultInferenceGeo: anthropic.DataResidencyCreateConfigDefaultInferenceGeoGlobal,
			WorkspaceGeo:        anthropic.DataResidencyCreateConfigWorkspaceGeoUs,
		},
		DisplayColor:  anthropic.String("#6C5BB9"),
		ExternalKeyID: anthropic.String("ekey_01SDCCSbTxrXDpWc1phhtcfK"),
		Tags: map[string]string{
			"env":  "prod",
			"team": "platform",
		},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrganizationWorkspaceGet(t *testing.T) {
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
	_, err := client.Organization.Workspaces.Get(context.TODO(), "workspace_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrganizationWorkspaceUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Workspaces.Update(
		context.TODO(),
		"workspace_id",
		anthropic.OrganizationWorkspaceUpdateParams{
			DataResidency: anthropic.DataResidencyUpdateConfigParam{
				AllowedInferenceGeos: anthropic.DataResidencyUpdateConfigAllowedInferenceGeosUnionParam{
					OfUnrestricted: constant.ValueOf[constant.Unrestricted](),
				},
				DefaultInferenceGeo: anthropic.DataResidencyUpdateConfigDefaultInferenceGeoGlobal,
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

func TestOrganizationWorkspaceListWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Workspaces.List(context.TODO(), anthropic.OrganizationWorkspaceListParams{
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

func TestOrganizationWorkspaceArchive(t *testing.T) {
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
	_, err := client.Organization.Workspaces.Archive(context.TODO(), "workspace_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

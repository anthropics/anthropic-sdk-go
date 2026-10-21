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

func TestBetaOrganizationWorkspaceRateLimitUpdate(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.RateLimits.Update(
		context.TODO(),
		"rate_limit_id",
		anthropic.BetaOrganizationWorkspaceRateLimitUpdateParams{
			WorkspaceID: "workspace_id",
			Limits: []anthropic.BetaWorkspaceRateLimitUpdateValueParam{{
				Type:  anthropic.BetaWorkspaceRateLimitTypeFastItpmca,
				Value: anthropic.Int(0),
			}},
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

func TestBetaOrganizationWorkspaceRateLimitListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Workspaces.RateLimits.List(
		context.TODO(),
		"workspace_id",
		anthropic.BetaOrganizationWorkspaceRateLimitListParams{
			GroupType:        anthropic.BetaOrganizationWorkspaceRateLimitListParamsGroupTypeBatch,
			IncludeInherited: anthropic.Bool(true),
			Limit:            anthropic.Int(1),
			Page:             anthropic.String("page"),
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

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

func TestBetaOrganizationAnalyticsUsageListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Analytics.Usage.List(context.TODO(), anthropic.BetaOrganizationAnalyticsUsageListParams{
		StartingAt:          time.Now(),
		ActorIDs:            []string{"user_01WCz1FkmYMm4gnmykNKUu3Q"},
		ActorTypes:          []anthropic.BetaAnalyticsActorType{anthropic.BetaAnalyticsActorTypeServiceAccountActor},
		AgentIDs:            []string{"agent_01VHd3EAGJ6PLkP9LM7DDu3f"},
		APIKeyIDs:           []string{"apikey_01Rj2N8SVvo6BePZj99NhmiT"},
		ClaudeTagCategories: []anthropic.BetaAnalyticsClaudeTagCategory{anthropic.BetaAnalyticsClaudeTagCategoryEngaged},
		ClaudeTagUserIDs:    []string{"U0123ABCDEF"},
		ContextWindows:      []anthropic.BetaAnalyticsContextWindow{anthropic.BetaAnalyticsContextWindowFrom0To200k},
		DeploymentIDs:       []string{"depl_01JFC6ihzZHtAuyHsaENdJhY"},
		EndingAt:            anthropic.Time(time.Now()),
		EnvironmentIDs:      []string{"env_017oTjuGedHtjhJauVC5xGJD"},
		ExcludeDeletedUsers: anthropic.Bool(true),
		GroupBy:             []anthropic.BetaAnalyticsUsageGroupBy{anthropic.BetaAnalyticsUsageGroupByActor},
		InferenceGeos:       []anthropic.BetaAnalyticsInferenceGeoFilter{anthropic.BetaAnalyticsInferenceGeoFilterGlobal},
		Limit:               anthropic.Int(1),
		Models:              []string{"string"},
		Order:               anthropic.BetaAnalyticsUsageCostSortOrderAsc,
		OrderBy:             anthropic.BetaAnalyticsUsageOrderByOutputTokens,
		Page:                anthropic.String("page"),
		Products:            []anthropic.BetaAnalyticsUsageCostProductFilter{anthropic.BetaAnalyticsUsageCostProductFilterChat},
		RBACGroupIDs:        []string{"rbac_group_012rppKaSVsmTo6NqRDXQXNF"},
		ServiceTiers:        []anthropic.BetaAnalyticsServiceTier{anthropic.BetaAnalyticsServiceTierBatch},
		SessionIDs:          []string{"sesn_0187kPDyNCwibNNvfaPEApCJ"},
		SlackChannelIDs:     []string{"C0123ABCDEF"},
		Speeds:              []anthropic.BetaAnalyticsSpeed{anthropic.BetaAnalyticsSpeedFast},
		TimeBucketWidth:     anthropic.BetaAnalyticsTimeBucketWidthDay,
		UserProfileIDs:      []string{"uprof_01M33NeYUQK2WV9PpYcsPHLi"},
		WorkspaceIDs:        []string{"wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ"},
		Betas:               []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

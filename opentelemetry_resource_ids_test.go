package anthropic_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// resourceIDKeys returns every ID attribute that the map has.
func resourceIDKeys() map[string]bool {
	keys := map[string]bool{}
	for _, ids := range resourceIDsOf {
		for _, key := range ids.path {
			keys[key] = true
		}
		if ids.created != "" {
			keys[ids.created] = true
		}
	}
	return keys
}

// resourceTransport answers the requests of the resource methods: a list with two
// pages of one object each, a stream with no events, the results of a batch, the
// events that a send created, and any other request with created.
func resourceTransport(created string) *closureTransport {
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		path, query := strings.TrimSuffix(req.URL.Path, "/"), req.URL.Query()
		switch {
		case strings.HasSuffix(path, "/stream"):
			return respond(200, "text/event-stream", "").fn(req)
		case strings.HasSuffix(path, "/results"):
			return respond(200, "application/binary", `{"custom_id":"a","result":{"type":"canceled"}}`+"\n").fn(req)
		case req.Method == http.MethodGet && (query.Has("page") || query.Has("after_id")):
			return respond(200, "application/json", `{"data":[{"id":"listed_2","type":"listed"}],"next_page":null,"has_more":false,"last_id":"listed_2"}`).fn(req)
		case req.Method == http.MethodGet && query.Has("limit"):
			return respond(200, "application/json", `{"data":[{"id":"listed_1","type":"listed"}],"next_page":"p2","has_more":true,"last_id":"listed_1"}`).fn(req)
		case req.Method == http.MethodPost && strings.HasSuffix(path, "/events"):
			return respond(200, "application/json", `{"data":[{"id":"sent_1","type":"user.message"},{"id":"sent_2","type":"user.message"}]}`).fn(req)
		default:
			return respond(200, "application/json", created).fn(req)
		}
	}}
}

// wantResourceIDs fails the test unless span has the ID attributes want, and no
// other ID attribute that the map has.
func wantResourceIDs(t *testing.T, span sdktrace.ReadOnlySpan, want map[string]string) {
	t.Helper()
	attrs := map[string]any{}
	for key := range resourceIDKeys() {
		attrs[key] = nil
		if id, ok := want[key]; ok {
			attrs[key] = id
		}
	}
	wantAttrs(t, span, attrs)
}

// wantRequestAttrs fails the test unless span has the HTTP method, the route and
// the URL of the method whose entry in the map is key: the route is the path of
// the entry, the URL's path has values in place of the route's parameters, and
// the URL's query has what the endpoint defines, so never beta=true. The span is
// that of the page with index page when the method is a list.
func wantRequestAttrs(t *testing.T, span sdktrace.ReadOnlySpan, key string, values []string, list bool, page int) {
	t.Helper()
	parts := strings.Split(key, " ")
	route := strings.TrimSuffix(parts[2], "?beta=true")
	path := route
	for _, value := range values {
		open, end := strings.Index(path, "{"), strings.Index(path, "}")
		path = path[:open] + value + path[end+1:]
	}
	wantAttrs(t, span, map[string]any{
		"http.request.method": parts[1],
		"http.route":          route,
		"url.scheme":          "http",
		"url.path":            path,
	})
	query, hasQuery := attr(span, "url.query")
	full := "http://localhost:4010" + path
	if hasQuery {
		full += "?" + query.AsString()
	}
	wantAttrs(t, span, map[string]any{"url.full": full})
	switch {
	case !list:
		if hasQuery {
			t.Errorf("%s: attribute url.query = %s, want absent", span.Name(), query.AsString())
		}
	case page == 0:
		wantAttrs(t, span, map[string]any{"url.query": "limit=1"})
	// The cursor of the second page is the parameter that the endpoint defines for it.
	case query.AsString() != "limit=1&page=p2" && query.AsString() != "after_id=listed_1&limit=1":
		t.Errorf("%s: attribute url.query = %s on the second page, want limit=1 and the cursor", span.Name(), query.AsString())
	}
}

// drain reads a list to its end.
func drain(iter interface {
	Next() bool
	Err() error
}) error {
	for iter.Next() {
	}
	return iter.Err()
}

func TestOpenTelemetryResourceIDs(t *testing.T) {
	ctx := context.Background()
	limit := anthropic.Int(1)
	// Each case is one method: its entry in the map, which names the ID attributes
	// of its span, the values that the call gives the parameters of the request's
	// path, in the path's order, and the call. The cases cover each resource, with a
	// method that has IDs in its path, the one that creates an object, and a list.
	cases := []struct {
		key    string
		values []string
		call   func(client anthropic.Client) error
	}{
		{"anthropic.agents.create POST /v1/agents?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{})
			return err
		}},
		{"anthropic.agents.update POST /v1/agents/{agent_id}?beta=true", []string{"agent_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Agents.Update(ctx, "agent_1", anthropic.BetaAgentUpdateParams{})
			return err
		}},
		{"anthropic.agents.versions.list GET /v1/agents/{agent_id}/versions?beta=true", []string{"agent_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Agents.Versions.ListAutoPaging(ctx, "agent_1", anthropic.BetaAgentVersionListParams{Limit: limit}))
		}},
		{"anthropic.deployments.create POST /v1/deployments?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Deployments.New(ctx, anthropic.BetaDeploymentNewParams{})
			return err
		}},
		// A run is not an object of the deployments path, so the span has no ID for it.
		{"anthropic.deployments.run POST /v1/deployments/{deployment_id}/run?beta=true", []string{"depl_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Deployments.Run(ctx, "depl_1", anthropic.BetaDeploymentRunParams{})
			return err
		}},
		{"anthropic.deployment_runs.retrieve GET /v1/deployment_runs/{deployment_run_id}?beta=true", []string{"drun_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.DeploymentRuns.Get(ctx, "drun_1", anthropic.BetaDeploymentRunGetParams{})
			return err
		}},
		{"anthropic.dreams.cancel POST /v1/dreams/{dream_id}/cancel?beta=true", []string{"dream_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Dreams.Cancel(ctx, "dream_1", anthropic.BetaDreamCancelParams{})
			return err
		}},
		{"anthropic.environments.create POST /v1/environments?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{})
			return err
		}},
		{"anthropic.environments.work.heartbeat POST /v1/environments/{environment_id}/work/{work_id}/heartbeat?beta=true", []string{"env_1", "work_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Environments.Work.Heartbeat(ctx, "work_1", anthropic.BetaEnvironmentWorkHeartbeatParams{EnvironmentID: "env_1"})
			return err
		}},
		{"anthropic.environments.work.list GET /v1/environments/{environment_id}/work?beta=true", []string{"env_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Environments.Work.ListAutoPaging(ctx, "env_1", anthropic.BetaEnvironmentWorkListParams{Limit: limit}))
		}},
		{"anthropic.files.retrieve_metadata GET /v1/files/{file_id}?beta=true", []string{"file_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Files.GetMetadata(ctx, "file_1", anthropic.BetaFileGetMetadataParams{})
			return err
		}},
		{"anthropic.files.download GET /v1/files/{file_id}/content", []string{"file_1"}, func(c anthropic.Client) error {
			res, err := c.Files.Download(ctx, "file_1", anthropic.FileDownloadParams{})
			if err == nil {
				err = res.Body.Close()
			}
			return err
		}},
		{"anthropic.memory_stores.create POST /v1/memory_stores?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.New(ctx, anthropic.BetaMemoryStoreNewParams{})
			return err
		}},
		{"anthropic.memory_stores.memories.create POST /v1/memory_stores/{memory_store_id}/memories?beta=true", []string{"memstore_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.Memories.New(ctx, "memstore_1", anthropic.BetaMemoryStoreMemoryNewParams{})
			return err
		}},
		{"anthropic.memory_stores.memories.retrieve GET /v1/memory_stores/{memory_store_id}/memories/{memory_id}?beta=true", []string{"memstore_1", "mem_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.Memories.Get(ctx, "mem_1", anthropic.BetaMemoryStoreMemoryGetParams{MemoryStoreID: "memstore_1"})
			return err
		}},
		{"anthropic.memory_stores.memory_versions.retrieve GET /v1/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}?beta=true", []string{"memstore_1", "memver_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.MemoryStores.MemoryVersions.Get(ctx, "memver_1", anthropic.BetaMemoryStoreMemoryVersionGetParams{MemoryStoreID: "memstore_1"})
			return err
		}},
		{"anthropic.memory_stores.memory_versions.list GET /v1/memory_stores/{memory_store_id}/memory_versions?beta=true", []string{"memstore_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.MemoryStores.MemoryVersions.ListAutoPaging(ctx, "memstore_1", anthropic.BetaMemoryStoreMemoryVersionListParams{Limit: limit}))
		}},
		{"anthropic.messages.batches.create POST /v1/messages/batches", nil, func(c anthropic.Client) error {
			_, err := c.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{})
			return err
		}},
		{"anthropic.messages.batches.create POST /v1/messages/batches?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Messages.Batches.New(ctx, anthropic.BetaMessageBatchNewParams{})
			return err
		}},
		{"anthropic.messages.batches.cancel POST /v1/messages/batches/{message_batch_id}/cancel", []string{"msgbatch_1"}, func(c anthropic.Client) error {
			_, err := c.Messages.Batches.Cancel(ctx, "msgbatch_1", anthropic.MessageBatchCancelParams{})
			return err
		}},
		{"anthropic.messages.batches.results GET /v1/messages/batches/{message_batch_id}/results?beta=true", []string{"msgbatch_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Messages.Batches.ResultsStreaming(ctx, "msgbatch_1", anthropic.BetaMessageBatchResultsParams{}))
		}},
		{"anthropic.models.retrieve GET /v1/models/{model_id}", []string{"claude-1"}, func(c anthropic.Client) error {
			_, err := c.Models.Get(ctx, "claude-1", anthropic.ModelGetParams{})
			return err
		}},
		{"anthropic.organization.api_keys.retrieve GET /v1/organizations/api_keys/{api_key_id}?beta=true", []string{"apikey_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.APIKeys.Get(ctx, "apikey_1")
			return err
		}},
		{"anthropic.organization.external_keys.create POST /v1/organizations/external_keys?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.ExternalKeys.New(ctx, anthropic.BetaOrganizationExternalKeyNewParams{})
			return err
		}},
		{"anthropic.organization.federation.issuers.archive POST /v1/organizations/federation_issuers/{federation_issuer_id}/archive?beta=true", []string{"fiss_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Federation.Issuers.Archive(ctx, "fiss_1", anthropic.BetaOrganizationFederationIssuerArchiveParams{})
			return err
		}},
		{"anthropic.organization.federation.rules.create POST /v1/organizations/federation_rules?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Federation.Rules.New(ctx, anthropic.BetaOrganizationFederationRuleNewParams{})
			return err
		}},
		// The object that an add creates has no `id`: its ID is the field that the map names.
		{"anthropic.organization.federation.rules.workspaces.add POST /v1/organizations/federation_rules/{federation_rule_id}/workspaces?beta=true", []string{"frul_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Federation.Rules.Workspaces.Add(ctx, "frul_1", anthropic.BetaOrganizationFederationRuleWorkspaceAddParams{})
			return err
		}},
		{"anthropic.organization.federation.rules.workspaces.remove DELETE /v1/organizations/federation_rules/{federation_rule_id}/workspaces/{workspace_id}?beta=true", []string{"frul_1", "wrkspc_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Federation.Rules.Workspaces.Remove(ctx, "wrkspc_1", anthropic.BetaOrganizationFederationRuleWorkspaceRemoveParams{FederationRuleID: "frul_1"})
			return err
		}},
		{"anthropic.organization.invites.create POST /v1/organizations/invites?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Invites.New(ctx, anthropic.BetaOrganizationInviteNewParams{})
			return err
		}},
		{"anthropic.organization.service_accounts.retrieve GET /v1/organizations/service_accounts/{service_account_id}?beta=true", []string{"svcacct_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.ServiceAccounts.Get(ctx, "svcacct_1", anthropic.BetaOrganizationServiceAccountGetParams{})
			return err
		}},
		{"anthropic.organization.service_accounts.workspaces.add POST /v1/organizations/service_accounts/{service_account_id}/workspaces?beta=true", []string{"svcacct_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.ServiceAccounts.Workspaces.Add(ctx, "svcacct_1", anthropic.BetaOrganizationServiceAccountWorkspaceAddParams{})
			return err
		}},
		{"anthropic.organization.users.update POST /v1/organizations/users/{user_id}?beta=true", []string{"user_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Users.Update(ctx, "user_1", anthropic.BetaOrganizationUserUpdateParams{})
			return err
		}},
		// A Workspace of the organization is not the Workspace that served the
		// request, which anthropic.workspace.id names.
		{"anthropic.organization.workspaces.create POST /v1/organizations/workspaces?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Workspaces.New(ctx, anthropic.BetaOrganizationWorkspaceNewParams{})
			return err
		}},
		{"anthropic.organization.workspaces.retrieve GET /v1/organizations/workspaces/{workspace_id}?beta=true", []string{"wrkspc_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Workspaces.Get(ctx, "wrkspc_1")
			return err
		}},
		{"anthropic.organization.workspaces.members.add POST /v1/organizations/workspaces/{workspace_id}/members?beta=true", []string{"wrkspc_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Workspaces.Members.Add(ctx, "wrkspc_1", anthropic.BetaOrganizationWorkspaceMemberAddParams{})
			return err
		}},
		{"anthropic.organization.workspaces.members.update POST /v1/organizations/workspaces/{workspace_id}/members/{user_id}?beta=true", []string{"wrkspc_1", "user_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Workspaces.Members.Update(ctx, "user_1", anthropic.BetaOrganizationWorkspaceMemberUpdateParams{WorkspaceID: "wrkspc_1"})
			return err
		}},
		{"anthropic.organization.workspaces.members.list GET /v1/organizations/workspaces/{workspace_id}/members?beta=true", []string{"wrkspc_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Organization.Workspaces.Members.ListAutoPaging(ctx, "wrkspc_1", anthropic.BetaOrganizationWorkspaceMemberListParams{Limit: limit}))
		}},
		{"anthropic.organization.workspaces.service_accounts.add POST /v1/organizations/workspaces/{workspace_id}/service_accounts?beta=true", []string{"wrkspc_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Organization.Workspaces.ServiceAccounts.Add(ctx, "wrkspc_1", anthropic.BetaOrganizationWorkspaceServiceAccountAddParams{})
			return err
		}},
		{"anthropic.organization.workspaces.rate_limits.list GET /v1/organizations/workspaces/{workspace_id}/rate_limits?beta=true", []string{"wrkspc_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Organization.Workspaces.RateLimits.ListAutoPaging(ctx, "wrkspc_1", anthropic.BetaOrganizationWorkspaceRateLimitListParams{Limit: limit}))
		}},
		{"anthropic.sessions.create POST /v1/sessions?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{})
			return err
		}},
		// An archive returns the Session it archived, and creates none.
		{"anthropic.sessions.archive POST /v1/sessions/{session_id}/archive?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Archive(ctx, "sesn_1", anthropic.BetaSessionArchiveParams{})
			return err
		}},
		{"anthropic.sessions.events.list GET /v1/sessions/{session_id}/events?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Sessions.Events.ListAutoPaging(ctx, "sesn_1", anthropic.BetaSessionEventListParams{Limit: limit}))
		}},
		// A send creates several events, and has no ID for them.
		{"anthropic.sessions.events.send POST /v1/sessions/{session_id}/events?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Events.Send(ctx, "sesn_1", anthropic.BetaSessionEventSendParams{})
			return err
		}},
		{"anthropic.sessions.events.stream GET /v1/sessions/{session_id}/events/stream?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Sessions.Events.StreamEvents(ctx, "sesn_1", anthropic.BetaSessionEventStreamParams{}))
		}},
		{"anthropic.sessions.resources.add POST /v1/sessions/{session_id}/resources?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Resources.Add(ctx, "sesn_1", anthropic.BetaSessionResourceAddParams{})
			return err
		}},
		{"anthropic.sessions.resources.delete DELETE /v1/sessions/{session_id}/resources/{resource_id}?beta=true", []string{"sesn_1", "sesrsc_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Resources.Delete(ctx, "sesrsc_1", anthropic.BetaSessionResourceDeleteParams{SessionID: "sesn_1"})
			return err
		}},
		{"anthropic.sessions.threads.retrieve GET /v1/sessions/{session_id}/threads/{thread_id}?beta=true", []string{"sesn_1", "sthr_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Threads.Get(ctx, "sthr_1", anthropic.BetaSessionThreadGetParams{SessionID: "sesn_1"})
			return err
		}},
		{"anthropic.sessions.threads.archive POST /v1/sessions/{session_id}/threads/{thread_id}/archive?beta=true", []string{"sesn_1", "sthr_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Sessions.Threads.Archive(ctx, "sthr_1", anthropic.BetaSessionThreadArchiveParams{SessionID: "sesn_1"})
			return err
		}},
		{"anthropic.sessions.threads.list GET /v1/sessions/{session_id}/threads?beta=true", []string{"sesn_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Sessions.Threads.ListAutoPaging(ctx, "sesn_1", anthropic.BetaSessionThreadListParams{Limit: limit}))
		}},
		{"anthropic.sessions.threads.events.list GET /v1/sessions/{session_id}/threads/{thread_id}/events?beta=true", []string{"sesn_1", "sthr_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Sessions.Threads.Events.ListAutoPaging(ctx, "sthr_1", anthropic.BetaSessionThreadEventListParams{SessionID: "sesn_1", Limit: limit}))
		}},
		{"anthropic.sessions.threads.events.stream GET /v1/sessions/{session_id}/threads/{thread_id}/stream?beta=true", []string{"sesn_1", "sthr_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Sessions.Threads.Events.StreamEvents(ctx, "sthr_1", anthropic.BetaSessionThreadEventStreamParams{SessionID: "sesn_1"}))
		}},
		{"anthropic.skills.retrieve GET /v1/skills/{skill_id}", []string{"skill_1"}, func(c anthropic.Client) error {
			_, err := c.Skills.Get(ctx, "skill_1", anthropic.SkillGetParams{})
			return err
		}},
		// The path's value is recorded though it is not an ID.
		{"anthropic.skills.versions.retrieve GET /v1/skills/{skill_id}/versions/{version}?beta=true", []string{"skill_1", "latest"}, func(c anthropic.Client) error {
			_, err := c.Beta.Skills.Versions.Get(ctx, "latest", anthropic.BetaSkillVersionGetParams{SkillID: "skill_1"})
			return err
		}},
		{"anthropic.skills.versions.list GET /v1/skills/{skill_id}/versions", []string{"skill_1"}, func(c anthropic.Client) error {
			return drain(c.Skills.Versions.ListAutoPaging(ctx, "skill_1", anthropic.SkillVersionListParams{Limit: limit}))
		}},
		{"anthropic.tunnels.create POST /v1/tunnels?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Tunnels.New(ctx, anthropic.BetaTunnelNewParams{})
			return err
		}},
		{"anthropic.tunnels.certificates.create POST /v1/tunnels/{tunnel_id}/certificates?beta=true", []string{"tnl_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Tunnels.Certificates.New(ctx, "tnl_1", anthropic.BetaTunnelCertificateNewParams{})
			return err
		}},
		{"anthropic.tunnels.certificates.archive POST /v1/tunnels/{tunnel_id}/certificates/{certificate_id}/archive?beta=true", []string{"tnl_1", "cert_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Tunnels.Certificates.Archive(ctx, "cert_1", anthropic.BetaTunnelCertificateArchiveParams{TunnelID: "tnl_1"})
			return err
		}},
		{"anthropic.user_profiles.create_enrollment_url POST /v1/user_profiles/{user_profile_id}/enrollment_url?beta=true", []string{"uprof_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.UserProfiles.NewEnrollmentURL(ctx, "uprof_1", anthropic.BetaUserProfileNewEnrollmentURLParams{})
			return err
		}},
		{"anthropic.vaults.create POST /v1/vaults?beta=true", nil, func(c anthropic.Client) error {
			_, err := c.Beta.Vaults.New(ctx, anthropic.BetaVaultNewParams{})
			return err
		}},
		{"anthropic.vaults.credentials.create POST /v1/vaults/{vault_id}/credentials?beta=true", []string{"vlt_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Vaults.Credentials.New(ctx, "vlt_1", anthropic.BetaVaultCredentialNewParams{})
			return err
		}},
		{"anthropic.vaults.credentials.list GET /v1/vaults/{vault_id}/credentials?beta=true", []string{"vlt_1"}, func(c anthropic.Client) error {
			return drain(c.Beta.Vaults.Credentials.ListAutoPaging(ctx, "vlt_1", anthropic.BetaVaultCredentialListParams{Limit: limit}))
		}},
		{"anthropic.vaults.credentials.mcp_oauth_validate POST /v1/vaults/{vault_id}/credentials/{credential_id}/mcp_oauth_validate?beta=true", []string{"vlt_1", "vcrd_1"}, func(c anthropic.Client) error {
			_, err := c.Beta.Vaults.Credentials.MCPOAuthValidate(ctx, "vcrd_1", anthropic.BetaVaultCredentialMCPOAuthValidateParams{VaultID: "vlt_1"})
			return err
		}},
	}
	namedFields := map[string]bool{}
	for _, test := range cases {
		ids, ok := resourceIDsOf[test.key]
		if !ok {
			t.Errorf("%s: the map has no such method", test.key)
			continue
		}
		name, _, _ := strings.Cut(test.key, " ")
		t.Run(test.key, func(t *testing.T) {
			if len(ids.path) != len(test.values) {
				t.Fatalf("the map has %d path parameters, and the call gives %d", len(ids.path), len(test.values))
			}
			want := map[string]string{}
			for i, key := range ids.path {
				want[key] = test.values[i]
			}
			// The object that the method creates has its ID in the field that the map
			// names, and in no other.
			created := `{"type":"created"}`
			if ids.created != "" {
				want[ids.created] = "new_" + ids.field
				created = `{"type":"created","` + ids.field + `":"new_` + ids.field + `"}`
				namedFields[ids.field] = true
			}
			pages := 1
			if strings.HasSuffix(name, ".list") {
				pages = 2
			}

			// The IDs are metadata: the client records no content.
			client, recorder := newTracedClient(t, resourceTransport(created))
			if err := test.call(client); err != nil {
				t.Fatal(err)
			}
			spans := recorder.Spans()
			if len(spans) != pages {
				t.Fatalf("got %d spans, want %d", len(spans), pages)
			}
			// The span of each page of a list has the IDs in the list's path, and no
			// ID of an object that the page lists.
			for page, span := range spans {
				if span.Name() != name {
					t.Errorf("span %q, want %q", span.Name(), name)
				}
				wantResourceIDs(t, span, want)
				// A Session thread in the path is the conversation, and nothing else in a path is.
				var conversation any
				if thread, ok := want["anthropic.session.thread.id"]; ok {
					conversation = thread
				}
				wantAttrs(t, span, map[string]any{"gen_ai.conversation.id": conversation})
				wantRequestAttrs(t, span, test.key, test.values, pages == 2, page)
				for _, kv := range span.Attributes() {
					// The cursor of a later page is a parameter of its request.
					if kv.Key == "url.query" || kv.Key == "url.full" {
						continue
					}
					if value := kv.Value.Emit(); strings.Contains(value, "listed_") || strings.Contains(value, "sent_") {
						t.Errorf("%s: attribute %s = %s is the ID of an object in the response", span.Name(), kv.Key, value)
					}
				}
			}
		})
	}
	for _, field := range []string{"id", "user_id", "workspace_id", "service_account_id"} {
		if !namedFields[field] {
			t.Errorf("no case creates an object whose ID is the response's %s", field)
		}
	}
}

func TestOpenTelemetryCreatedIDIsAbsentWithoutACreatedObject(t *testing.T) {
	ctx := context.Background()
	creates := map[string]func(client anthropic.Client, opts ...option.RequestOption) error{
		"anthropic.sessions.create": func(client anthropic.Client, opts ...option.RequestOption) error {
			_, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{}, opts...)
			return err
		},
		"anthropic.agents.create": func(client anthropic.Client, opts ...option.RequestOption) error {
			_, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{}, opts...)
			return err
		},
		"anthropic.messages.batches.create": func(client anthropic.Client, opts ...option.RequestOption) error {
			_, err := client.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{}, opts...)
			return err
		},
	}
	for name, create := range creates {
		t.Run(name+" failed", func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(400, "application/json",
				`{"type":"error","error":{"type":"invalid_request_error","message":"bad"},"id":"not_created"}`))
			if err := create(client); err == nil {
				t.Fatal("the call did not fail")
			}
			span := onlySpan(t, recorder)
			if span.Name() != name || span.Status().Code != codes.Error {
				t.Errorf("span %q with status %v, want %q with ERROR", span.Name(), span.Status().Code, name)
			}
			wantResourceIDs(t, span, nil)
			wantAttrs(t, span, map[string]any{"error.type": "invalid_request_error"})
		})
		t.Run(name+" without an ID", func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", `{"type":"created","id":null}`))
			if err := create(client); err != nil {
				t.Fatal(err)
			}
			wantResourceIDs(t, onlySpan(t, recorder), nil)
		})
		t.Run(name+" decoded elsewhere", func(t *testing.T) {
			// The method never saw the response, so its span has no ID from it.
			client, recorder := newTracedClient(t, respond(200, "application/json", `{"id":"new_1","type":"created"}`))
			var into map[string]any
			if err := create(client, option.WithResponseBodyInto(&into)); err != nil || into["id"] != "new_1" {
				t.Fatalf("error %v, and %v in the application's value", err, into)
			}
			wantResourceIDs(t, onlySpan(t, recorder), nil)
		})
	}
}

const sessionJSON = `{"id":"sesn_1","type":"session","environment_id":"env_1","agent":{"id":"agent_1","type":"agent","name":"Researcher","version":2,"model":{"id":"claude-resolved"}}}`

func TestOpenTelemetrySessionsCreateHasTheAgentModelAndEnvironment(t *testing.T) {
	ctx := context.Background()
	// The Agent is an ID, or an object with one. The model that an object overrides
	// is not the one that the span has: the Session's is.
	agents := map[string]anthropic.BetaSessionNewParamsAgentUnion{
		"string": {OfString: anthropic.String("agent_1")},
		"agent":  {OfBetaManagedAgentsAgents: &anthropic.BetaManagedAgentsAgentParams{ID: "agent_1", Version: anthropic.Int(2)}},
		"agent with overrides": {OfBetaManagedAgentsAgentWithOverridess: &anthropic.BetaManagedAgentsAgentWithOverridesParams{
			ID:    "agent_1",
			Model: anthropic.BetaManagedAgentsModelConfigParams{ID: "claude-overridden"},
		}},
	}
	modes := []option.OpenTelemetryContentMode{option.OpenTelemetryContentModeMetadataOnly, option.OpenTelemetryContentModeContent}
	for name, agent := range agents {
		for _, mode := range modes {
			t.Run(name+" "+string(mode), func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "application/json", sessionJSON), option.WithOpenTelemetryTracesContentMode(mode))
				if _, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{Agent: agent, EnvironmentID: "env_1"}); err != nil {
					t.Fatal(err)
				}
				// A sampler sees what the request has.
				want := map[string]any{
					"gen_ai.provider.name":             "anthropic",
					"server.address":                   "localhost",
					"server.port":                      int64(4010),
					"http.request.method":              "POST",
					"http.route":                       "/v1/sessions",
					"gen_ai.agent.id":                  "agent_1",
					"anthropic.session.environment.id": "env_1",
				}
				if got := recorder.sampler.last("anthropic.sessions.create"); !reflect.DeepEqual(got, want) {
					t.Errorf("the span started with %v, want %v", got, want)
				}
				wantAttrs(t, onlySpan(t, recorder), map[string]any{
					"anthropic.session.id":             "sesn_1",
					"gen_ai.agent.id":                  "agent_1",
					"anthropic.session.environment.id": "env_1",
					"gen_ai.agent.name":                "Researcher",
					"gen_ai.request.model":             "claude-resolved",
				})
			})
		}
	}
}

func TestOpenTelemetrySessionsCreateThatFailedHasTheAgentAndEnvironment(t *testing.T) {
	client, recorder := newTracedClient(t, respond(400, "application/json",
		`{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`))
	_, err := client.Beta.Sessions.New(context.Background(), anthropic.BetaSessionNewParams{
		Agent:         anthropic.BetaSessionNewParamsAgentUnion{OfString: anthropic.String("agent_1")},
		EnvironmentID: "env_1",
	})
	if err == nil {
		t.Fatal("the call did not fail")
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"error.type":                       "invalid_request_error",
		"gen_ai.agent.id":                  "agent_1",
		"anthropic.session.environment.id": "env_1",
		"anthropic.session.id":             nil,
		"gen_ai.agent.name":                nil,
		"gen_ai.request.model":             nil,
	})
}

func TestOpenTelemetrySessionsRetrieveHasNoAgentModelOrEnvironment(t *testing.T) {
	// Only the span of a create has them: the API's span of a retrieve has none.
	client, recorder := newTracedClient(t, respond(200, "application/json", sessionJSON))
	if _, err := client.Beta.Sessions.Get(context.Background(), "sesn_1", anthropic.BetaSessionGetParams{}); err != nil {
		t.Fatal(err)
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"anthropic.session.id":             "sesn_1",
		"gen_ai.agent.id":                  nil,
		"anthropic.session.environment.id": nil,
		"gen_ai.agent.name":                nil,
		"gen_ai.request.model":             nil,
	})
}

func TestOpenTelemetrySessionThreadIsTheConversation(t *testing.T) {
	// The API links its spans of a thread by the thread's ID, so the span of a
	// call to the thread's path starts with that ID as its conversation, unless
	// the application named one.
	ctx := context.Background()
	calls := map[string]func(client anthropic.Client, opts ...option.RequestOption) error{
		"anthropic.sessions.threads.retrieve": func(client anthropic.Client, opts ...option.RequestOption) error {
			_, err := client.Beta.Sessions.Threads.Get(ctx, "sthr_1", anthropic.BetaSessionThreadGetParams{SessionID: "sesn_1"}, opts...)
			return err
		},
		"anthropic.sessions.threads.events.stream": func(client anthropic.Client, opts ...option.RequestOption) error {
			stream := client.Beta.Sessions.Threads.Events.StreamEvents(ctx, "sthr_1", anthropic.BetaSessionThreadEventStreamParams{SessionID: "sesn_1"}, opts...)
			for stream.Next() {
			}
			return stream.Err()
		},
	}
	tests := map[string]struct {
		client, request []option.RequestOption
		want            string
	}{
		"unset":              {nil, nil, "sthr_1"},
		"set on the client":  {[]option.RequestOption{option.WithOpenTelemetryConversationID("conv_client")}, nil, "conv_client"},
		"set on the request": {nil, []option.RequestOption{option.WithOpenTelemetryConversationID("conv_request")}, "conv_request"},
	}
	for name, call := range calls {
		for when, test := range tests {
			t.Run(name+" "+when, func(t *testing.T) {
				client, recorder := newTracedClient(t, resourceTransport(`{"type":"session_thread"}`), test.client...)
				if err := call(client, test.request...); err != nil {
					t.Fatal(err)
				}
				span := onlySpan(t, recorder)
				wantAttrs(t, span, map[string]any{
					"gen_ai.conversation.id":      test.want,
					"anthropic.session.thread.id": "sthr_1",
				})
				if started := recorder.sampler.last(name)["gen_ai.conversation.id"]; started != test.want {
					t.Errorf("the span started with the conversation %v, want %q", started, test.want)
				}
				conversations := 0
				for _, kv := range span.Attributes() {
					if kv.Key == "gen_ai.conversation.id" {
						conversations++
					}
				}
				if conversations != 1 {
					t.Errorf("the span has %d conversation IDs, want 1", conversations)
				}
			})
		}
	}
}

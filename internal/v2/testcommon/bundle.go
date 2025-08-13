package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func GenerateBundleResponse() *client.BundleV2 {
	bundle := &client.BundleV2{
		ID:   "bundle-123",
		Name: "Test Bundle",
	}

	integrationTarget := client.AccessBundleAccessTargetV2{}
	integrationData := client.IntegrationAccessTargetV2{
		IntegrationID:   "integration-123",
		IntegrationName: "postgresql",
		ResourceType:    "database",
		Permissions:     []string{"read", "write"},
	}
	resourceScope := client.ResourcesScopeIntegrationAccessTargetV2{
		ScopeMode: "include_resources",
		Type:      "NAME",
		Values:    []string{"db1", "db2"},
	}
	integrationData.ResourcesScopes.SetTo([]client.ResourcesScopeIntegrationAccessTargetV2{resourceScope})
	integrationTarget.Integration.SetTo(integrationData)

	accessScopeTarget := client.AccessBundleAccessTargetV2{}
	accessScopeData := client.AccessScopeAccessTargetV2{
		AccessScopeID:   "scope-123",
		AccessScopeName: "Test Scope",
	}
	accessScopeTarget.AccessScope.SetTo(accessScopeData)

	bundle.AccessTargets = []client.AccessBundleAccessTargetV2{integrationTarget, accessScopeTarget}

	return bundle
}

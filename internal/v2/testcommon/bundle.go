package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// GenerateBundleResponse creates a mock bundle response for testing
func GenerateBundleResponse() *client.BundlePublicV2Model {
	bundle := &client.BundlePublicV2Model{
		ID:   "bundle-123",
		Name: "Test Bundle",
	}

	// Add an integration target
	integrationTarget := client.AccessBundleAccessTargetPublicV2Model{}
	integrationData := client.IntegrationAccessTargetPublicV2Model{
		IntegrationID:   "integration-123",
		IntegrationName: "postgresql",
		ResourceType:    "database",
		Permissions:     []string{"read", "write"},
	}
	resourceScope := client.ResourcesScopeIntegrationAccessTargetPublicV2Model{
		ScopeMode: "include_resources",
		Type:      "NAME",
		Values:    []string{"db1", "db2"},
	}
	integrationData.ResourcesScopes.SetTo([]client.ResourcesScopeIntegrationAccessTargetPublicV2Model{resourceScope})
	integrationTarget.Integration.SetTo(integrationData)

	// Add an access scope target
	accessScopeTarget := client.AccessBundleAccessTargetPublicV2Model{}
	accessScopeData := client.AccessScopeAccessTargetPublicV2Model{
		AccessScopeID:   "scope-123",
		AccessScopeName: "Test Scope",
	}
	accessScopeTarget.AccessScope.SetTo(accessScopeData)

	bundle.AccessTargets = []client.AccessBundleAccessTargetPublicV2Model{integrationTarget, accessScopeTarget}

	return bundle
}

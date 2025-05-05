package models

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundleConversions(t *testing.T) {
	ctx := t.Context()

	t.Run("BundleResponseToModel", func(t *testing.T) {
		response := client.BundlePublicV2Model{
			ID:   "bundle-123",
			Name: "Test Bundle",
		}

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

		accessScopeTarget := client.AccessBundleAccessTargetPublicV2Model{}
		accessScopeData := client.AccessScopeAccessTargetPublicV2Model{
			AccessScopeID:   "scope-123",
			AccessScopeName: "Test Scope",
		}
		accessScopeTarget.AccessScope.SetTo(accessScopeData)

		response.AccessTargets = []client.AccessBundleAccessTargetPublicV2Model{integrationTarget, accessScopeTarget}

		model, err := BundleResponseToModel(ctx, response)
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "bundle-123", model.ID.ValueString())
		assert.Equal(t, "Test Bundle", model.Name.ValueString())

		require.Len(t, model.AccessTargets, 2)

		require.NotNil(t, model.AccessTargets[0].Integration)
		assert.Equal(t, "postgresql", model.AccessTargets[0].Integration.IntegrationName.ValueString())
		assert.Equal(t, "database", model.AccessTargets[0].Integration.ResourceType.ValueString())

		var permissions []string
		diags := model.AccessTargets[0].Integration.Permissions.ElementsAs(ctx, &permissions, false)
		require.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"read", "write"}, permissions)

		require.Len(t, model.AccessTargets[0].Integration.ResourcesScopes, 1)
		assert.Equal(t, "include_resources", model.AccessTargets[0].Integration.ResourcesScopes[0].ScopeMode.ValueString())
		assert.Equal(t, "NAME", model.AccessTargets[0].Integration.ResourcesScopes[0].Type.ValueString())

		var scopeValues []string
		diags = model.AccessTargets[0].Integration.ResourcesScopes[0].Values.ElementsAs(ctx, &scopeValues, false)
		require.False(t, diags.HasError())
		assert.ElementsMatch(t, []string{"db1", "db2"}, scopeValues)

		require.NotNil(t, model.AccessTargets[1].AccessScope)
		assert.Equal(t, "Test Scope", model.AccessTargets[1].AccessScope.Name.ValueString())
	})

	t.Run("BundleModelToUpsertRequest", func(t *testing.T) {
		model := BundleV2Model{
			ID:   types.StringValue("bundle-123"),
			Name: types.StringValue("Test Bundle"),
			AccessTargets: []BundleAccessTargetModel{
				{
					Integration: &IntegrationTargetModel{
						IntegrationName: types.StringValue("postgresql"),
						ResourceType:    types.StringValue("database"),
						Permissions:     testcommon.CreateTestStringSet(t, []string{"read", "write"}),
						ResourcesScopes: []IntegrationTargetScopeModel{
							{
								ScopeMode: types.StringValue("include_resources"),
								Type:      types.StringValue("NAME"),
								Key:       types.StringNull(),
								Values:    testcommon.CreateTestStringSet(t, []string{"db1", "db2"}),
							},
						},
					},
				},
				{
					AccessScope: &AccessScopeTargetModel{
						Name: types.StringValue("Test Scope"),
					},
				},
			},
		}

		request, err := BundleModelToUpsertRequest(ctx, model)
		require.NoError(t, err)
		require.NotNil(t, request)

		assert.Equal(t, "Test Bundle", request.Name)

		require.Len(t, request.AccessTargets, 2)

		assert.True(t, request.AccessTargets[0].Integration.IsSet())
		integration, ok := request.AccessTargets[0].Integration.Get()
		require.True(t, ok)
		assert.Equal(t, "postgresql", integration.IntegrationReference)
		assert.Equal(t, "database", integration.ResourceType)
		assert.ElementsMatch(t, []string{"read", "write"}, integration.Permissions)

		require.True(t, integration.ResourcesScopes.IsSet())
		resourceScopes, ok := integration.ResourcesScopes.Get()
		require.True(t, ok)
		require.Len(t, resourceScopes, 1)
		assert.Equal(t, "include_resources", resourceScopes[0].ScopeMode)
		assert.Equal(t, "NAME", resourceScopes[0].Type)
		assert.False(t, resourceScopes[0].Key.IsSet())
		assert.ElementsMatch(t, []string{"db1", "db2"}, resourceScopes[0].Values)

		assert.True(t, request.AccessTargets[1].AccessScope.IsSet())
		accessScope, ok := request.AccessTargets[1].AccessScope.Get()
		require.True(t, ok)
		assert.Equal(t, "Test Scope", accessScope.AccessScopeReference)
	})

	t.Run("BundleModelToUpsertRequest_ValidationError", func(t *testing.T) {
		model := BundleV2Model{
			Name: types.StringValue("Invalid Bundle"),
			AccessTargets: []BundleAccessTargetModel{
				{
					Integration: &IntegrationTargetModel{
						IntegrationName: types.StringValue("postgresql"),
						ResourceType:    types.StringValue("database"),
						Permissions:     testcommon.CreateTestStringSet(t, []string{"read"}),
					},
					AccessScope: &AccessScopeTargetModel{
						Name: types.StringValue("Test Scope"),
					},
				},
			},
		}

		_, err := BundleModelToUpsertRequest(ctx, model)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one of 'integration' or 'access_scope' must be configured")
	})
}

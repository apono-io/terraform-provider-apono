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
		response := client.BundleV2{
			ID:   "bundle-123",
			Name: "Test Bundle",
		}
		response.Space.SetTo(client.SpaceReferenceV1{SpaceID: "space-123", SpaceName: "prod-space"})

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

		queryTarget := client.AccessBundleAccessTargetV2{}
		queryTarget.QueryTarget.SetTo(client.QueryAccessTargetV2{Query: `integration_name = "aws"`})

		response.AccessTargets = []client.AccessBundleAccessTargetV2{integrationTarget, accessScopeTarget, queryTarget}

		model, err := BundleResponseToModel(ctx, response)
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "bundle-123", model.ID.ValueString())
		assert.Equal(t, "Test Bundle", model.Name.ValueString())

		require.Len(t, model.AccessTargets, 3)

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

		require.NotNil(t, model.AccessTargets[2].QueryTarget)
		assert.Equal(t, `integration_name = "aws"`, model.AccessTargets[2].QueryTarget.Query.ValueString())

		assert.Equal(t, "prod-space", model.SpaceReference.ValueString())
		require.False(t, model.Space.IsNull())
		spaceAttrs := model.Space.Attributes()
		spaceID, ok := spaceAttrs["space_id"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "space-123", spaceID.ValueString())
		spaceName, ok := spaceAttrs["space_name"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "prod-space", spaceName.ValueString())
	})

	t.Run("BundleResponseToModelWithoutSpace", func(t *testing.T) {
		response := client.BundleV2{
			ID:   "bundle-123",
			Name: "Test Bundle",
		}
		response.AccessTargets = []client.AccessBundleAccessTargetV2{}

		model, err := BundleResponseToModel(ctx, response)
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.True(t, model.SpaceReference.IsNull())
		assert.True(t, model.Space.IsNull())
	})

	t.Run("BundleModelToUpsertRequest", func(t *testing.T) {
		model := BundleV2Model{
			ID:             types.StringValue("bundle-123"),
			Name:           types.StringValue("Test Bundle"),
			SpaceReference: types.StringNull(),
			Space:          types.ObjectNull(SpaceReferenceAttrTypes),
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
								Values:    testcommon.CreateTestStringList(t, []string{"db1", "db2"}),
							},
						},
					},
				},
				{
					AccessScope: &AccessScopeTargetModel{
						Name: types.StringValue("Test Scope"),
					},
				},
				{
					QueryTarget: &QueryTargetModel{
						Query: types.StringValue(`integration_name = "aws"`),
					},
				},
			},
		}

		request, err := BundleModelToUpsertRequest(ctx, model)
		require.NoError(t, err)
		require.NotNil(t, request)

		assert.Equal(t, "Test Bundle", request.Name)

		require.Len(t, request.AccessTargets, 3)

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

		assert.True(t, request.AccessTargets[2].QueryTarget.IsSet())
		qt, ok := request.AccessTargets[2].QueryTarget.Get()
		require.True(t, ok)
		assert.Equal(t, `integration_name = "aws"`, qt.Query)
	})

	t.Run("BundleResponseToDataSourceItemModel_WithSpace", func(t *testing.T) {
		response := client.BundleV2{
			ID:            "bundle-ds-1",
			Name:          "DS Bundle",
			AccessTargets: []client.AccessBundleAccessTargetV2{},
		}
		response.Space.SetTo(client.SpaceReferenceV1{SpaceID: "space-123", SpaceName: "prod-space"})

		model, err := BundleResponseToDataSourceItemModel(ctx, response)
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "bundle-ds-1", model.ID.ValueString())
		assert.Equal(t, "DS Bundle", model.Name.ValueString())
		require.False(t, model.Space.IsNull())
		spaceAttrs := model.Space.Attributes()
		spaceID, ok := spaceAttrs["space_id"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "space-123", spaceID.ValueString())
		spaceName, ok := spaceAttrs["space_name"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "prod-space", spaceName.ValueString())
	})

	t.Run("BundleResponseToDataSourceItemModel_WithoutSpace", func(t *testing.T) {
		response := client.BundleV2{
			ID:            "bundle-ds-2",
			Name:          "DS Bundle No Space",
			AccessTargets: []client.AccessBundleAccessTargetV2{},
		}

		model, err := BundleResponseToDataSourceItemModel(ctx, response)
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "bundle-ds-2", model.ID.ValueString())
		assert.True(t, model.Space.IsNull())
	})
}

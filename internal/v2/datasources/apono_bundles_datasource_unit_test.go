package datasources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoBundlesDataSource(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		d := &AponoBundlesDataSource{client: mockInvoker}

		ctx := t.Context()

		bundles := []client.BundleV2{
			{
				ID:   "bundle-123",
				Name: "test-bundle-1",
				AccessTargets: []client.AccessBundleAccessTargetV2{
					{
						Integration: client.NewOptNilIntegrationAccessTargetV2(
							client.IntegrationAccessTargetV2{
								IntegrationName: "test-integration",
								ResourceType:    "db",
								Permissions:     []string{"read", "write"},
								ResourcesScopes: client.NewOptNilResourcesScopeIntegrationAccessTargetV2Array([]client.ResourcesScopeIntegrationAccessTargetV2{
									{
										ScopeMode: "include_resources",
										Type:      "NAME",
										Key:       client.NewOptNilString(""),
										Values:    []string{"resource1", "resource2"},
									},
								}),
							},
						),
					},
				},
			},
			{
				ID:   "bundle-456",
				Name: "test-bundle-2",
				AccessTargets: []client.AccessBundleAccessTargetV2{
					{
						AccessScope: client.NewOptNilAccessScopeAccessTargetV2(
							client.AccessScopeAccessTargetV2{
								AccessScopeName: "test-access-scope",
							},
						),
					},
				},
			},
			{
				ID:   "bundle-789",
				Name: "test-bundle-3",
				AccessTargets: []client.AccessBundleAccessTargetV2{
					func() client.AccessBundleAccessTargetV2 {
						t := client.AccessBundleAccessTargetV2{}
						t.QueryTarget.SetTo(client.QueryAccessTargetV2{Query: `integration_name = "aws"`})
						return t
					}(),
				},
			},
		}

		mockInvoker.EXPECT().
			ListBundlesV2(mock.Anything, mock.Anything).
			Return(&client.PublicApiListResponseBundlePublicV2Model{
				Items:      bundles,
				Pagination: client.PublicApiPaginationInfoModel{},
			}, nil)

		schema := d.getTestSchema(ctx)

		plan := tfsdk.Plan{
			Schema: schema,
		}

		diag := plan.Set(ctx, models.BundlesDataModel{
			SpaceReferences: types.ListNull(types.StringType),
		})
		require.False(t, diag.HasError(), "Error setting plan: %s", diag.Errors())

		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: schema,
				Raw:    plan.Raw,
			},
		}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Schema: schema,
				Raw:    tftypes.NewValue(schema.Type().TerraformType(ctx), nil),
			},
		}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Read returned error: %s", resp.Diagnostics.Errors())

		var state models.BundlesDataModel
		resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
		require.False(t, resp.Diagnostics.HasError(), "Error getting state: %s", resp.Diagnostics.Errors())

		assert.Len(t, state.Bundles, 3, "Expected 3 bundles")

		// Directly check the order since bundles are sorted by ID
		bundle1 := state.Bundles[0]
		assert.Equal(t, "bundle-123", bundle1.ID.ValueString())
		assert.Equal(t, "test-bundle-1", bundle1.Name.ValueString())
		assert.True(t, bundle1.Space.IsNull())
		require.Len(t, bundle1.AccessTargets, 1)
		assert.NotNil(t, bundle1.AccessTargets[0].Integration)
		assert.Equal(t, "test-integration", bundle1.AccessTargets[0].Integration.IntegrationName.ValueString())

		bundle2 := state.Bundles[1]
		assert.Equal(t, "bundle-456", bundle2.ID.ValueString())
		assert.Equal(t, "test-bundle-2", bundle2.Name.ValueString())
		assert.True(t, bundle2.Space.IsNull())
		require.Len(t, bundle2.AccessTargets, 1)
		assert.NotNil(t, bundle2.AccessTargets[0].AccessScope)
		assert.Equal(t, "test-access-scope", bundle2.AccessTargets[0].AccessScope.Name.ValueString())

		bundle3 := state.Bundles[2]
		assert.Equal(t, "bundle-789", bundle3.ID.ValueString())
		assert.Equal(t, "test-bundle-3", bundle3.Name.ValueString())
		assert.True(t, bundle3.Space.IsNull())
		require.Len(t, bundle3.AccessTargets, 1)
		require.NotNil(t, bundle3.AccessTargets[0].QueryTarget)
		assert.Equal(t, `integration_name = "aws"`, bundle3.AccessTargets[0].QueryTarget.Query.ValueString())
	})

	t.Run("Read_WithSpaceReferences", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		d := &AponoBundlesDataSource{client: mockInvoker}

		ctx := t.Context()

		spaceBundle := client.BundleV2{
			ID:            "bundle-sp1",
			Name:          "space-bundle",
			AccessTargets: []client.AccessBundleAccessTargetV2{},
		}
		spaceBundle.Space.SetTo(client.SpaceReferenceV1{SpaceID: "sp-1", SpaceName: "prod"})

		mockInvoker.EXPECT().
			ListBundlesV2(mock.Anything, mock.MatchedBy(func(params client.ListBundlesV2Params) bool {
				val, ok := params.SpaceReferences.Get()
				return ok && len(val) == 1 && val[0] == "prod"
			})).
			Return(&client.PublicApiListResponseBundlePublicV2Model{
				Items:      []client.BundleV2{spaceBundle},
				Pagination: client.PublicApiPaginationInfoModel{},
			}, nil)

		schema := d.getTestSchema(ctx)
		plan := tfsdk.Plan{Schema: schema}
		diag := plan.Set(ctx, models.BundlesDataModel{
			SpaceReferences: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("prod")}),
		})
		require.False(t, diag.HasError(), "Error setting plan: %s", diag.Errors())

		req := datasource.ReadRequest{
			Config: tfsdk.Config{Schema: schema, Raw: plan.Raw},
		}
		resp := datasource.ReadResponse{
			State: tfsdk.State{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), nil)},
		}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Read returned error: %s", resp.Diagnostics.Errors())

		var state models.BundlesDataModel
		resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
		require.False(t, resp.Diagnostics.HasError())

		require.Len(t, state.Bundles, 1)
		require.False(t, state.Bundles[0].Space.IsNull())
		spaceName, ok := state.Bundles[0].Space.Attributes()["space_name"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "prod", spaceName.ValueString())
	})
}

func (d *AponoBundlesDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

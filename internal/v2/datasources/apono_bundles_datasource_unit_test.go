package datasources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

		bundles := []client.BundlePublicV2Model{
			{
				ID:   "bundle-123",
				Name: "test-bundle-1",
				AccessTargets: []client.AccessBundleAccessTargetPublicV2Model{
					{
						Integration: client.NewOptNilIntegrationAccessTargetPublicV2Model(
							client.IntegrationAccessTargetPublicV2Model{
								IntegrationName: "test-integration",
								ResourceType:    "db",
								Permissions:     []string{"read", "write"},
								ResourcesScopes: client.NewOptNilResourcesScopeIntegrationAccessTargetPublicV2ModelArray([]client.ResourcesScopeIntegrationAccessTargetPublicV2Model{
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
				AccessTargets: []client.AccessBundleAccessTargetPublicV2Model{
					{
						AccessScope: client.NewOptNilAccessScopeAccessTargetPublicV2Model(
							client.AccessScopeAccessTargetPublicV2Model{
								AccessScopeName: "test-access-scope",
							},
						),
					},
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

		diag := plan.Set(ctx, models.BundlesV2DataModel{})
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

		var state models.BundlesV2DataModel
		resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
		require.False(t, resp.Diagnostics.HasError(), "Error getting state: %s", resp.Diagnostics.Errors())

		assert.Len(t, state.Bundles, 2, "Expected 2 bundles")

		// Directly check the order since bundles are sorted by ID
		bundle1 := state.Bundles[0]
		assert.Equal(t, "bundle-123", bundle1.ID.ValueString())
		assert.Equal(t, "test-bundle-1", bundle1.Name.ValueString())
		require.Len(t, bundle1.AccessTargets, 1)
		assert.NotNil(t, bundle1.AccessTargets[0].Integration)
		assert.Equal(t, "test-integration", bundle1.AccessTargets[0].Integration.IntegrationName.ValueString())

		bundle2 := state.Bundles[1]
		assert.Equal(t, "bundle-456", bundle2.ID.ValueString())
		assert.Equal(t, "test-bundle-2", bundle2.Name.ValueString())
		require.Len(t, bundle2.AccessTargets, 1)
		assert.NotNil(t, bundle2.AccessTargets[0].AccessScope)
		assert.Equal(t, "test-access-scope", bundle2.AccessTargets[0].AccessScope.Name.ValueString())
	})
}

func (d *AponoBundlesDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

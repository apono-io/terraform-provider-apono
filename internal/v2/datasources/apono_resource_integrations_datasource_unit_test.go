package datasources

import (
	"context"
	"testing"
	"time"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceIntegrationsDataSource_Read(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	d := &ResourceIntegrationsDataSource{client: mockInvoker}

	ctx := t.Context()

	now := time.Now().UTC()
	integrations := []client.IntegrationV4{
		{
			ID:                     "integration-1",
			Name:                   "test-integration-1",
			Type:                   "postgresql",
			Category:               common.ResourceCategory,
			Status:                 "connected",
			ConnectorID:            client.NewOptNilString("conn-1"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"db"}),
			IntegrationConfig: map[string]jx.Raw{
				"host": common.StringToJx("localhost"),
			},
			LastSyncTime: client.NewOptNilApiInstant(client.ApiInstant(now)),
		},
		{
			ID:                     "integration-2",
			Name:                   "test-integration-2",
			Type:                   "mysql",
			Category:               common.ResourceCategory,
			Status:                 "connected",
			ConnectorID:            client.NewOptNilString("conn-2"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"db"}),
			IntegrationConfig: map[string]jx.Raw{
				"host": common.StringToJx("127.0.0.1"),
			},
			LastSyncTime: client.NewOptNilApiInstant(client.ApiInstant(now)),
		},
	}

	mockInvoker.EXPECT().
		ListIntegrationsV4(mock.Anything, mock.MatchedBy(func(params client.ListIntegrationsV4Params) bool {
			categoryParam, ok := params.Category.Get()
			return ok && len(categoryParam) == 1 && categoryParam[0] == common.ResourceCategory
		})).
		Return(&client.PublicApiListResponseIntegrationPublicV4Model{
			Items:      integrations,
			Pagination: client.PublicApiPaginationInfoModel{},
		}, nil)

	schema := getResourceIntegrationsTestSchema(ctx, d)

	plan := tfsdk.Plan{
		Schema: schema,
	}

	diag := plan.Set(ctx, models.ResourceIntegrationsDataSourceModel{})
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

	var state models.ResourceIntegrationsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	require.False(t, resp.Diagnostics.HasError(), "Error getting state: %s", resp.Diagnostics.Errors())

	expectedModels := []models.ResourceIntegrationModel{}
	for _, integration := range integrations {
		model, err := models.ResourceIntegrationToModel(ctx, &integration)
		require.NoError(t, err, "Error converting integration to model")
		expectedModels = append(expectedModels, *model)
	}

	assert.Equal(t, expectedModels, state.Integrations, "Integrations do not match expected models")
}

func getResourceIntegrationsTestSchema(ctx context.Context, d *ResourceIntegrationsDataSource) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

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

func TestAponoUserInformationIntegrationsDataSource(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		d := &AponoUserInformationIntegrationsDataSource{
			client: mockInvoker,
		}

		ctx := t.Context()

		now := time.Now().UTC()
		integrations := []client.IntegrationV4{
			{
				ID:           "int-1",
				Name:         "integration-1",
				Type:         "ldap",
				Category:     "USER-INFORMATION",
				Status:       "ACTIVE",
				LastSyncTime: client.NewOptNilDateTime(now),
				IntegrationConfig: map[string]jx.Raw{
					"url": common.StringToJx("ldap://example.com"),
				},
			},
			{
				ID:           "int-2",
				Name:         "integration-2",
				Type:         "okta",
				Category:     "USER-INFORMATION",
				Status:       "ACTIVE",
				LastSyncTime: client.NewOptNilDateTime(now),
				IntegrationConfig: map[string]jx.Raw{
					"domain": common.StringToJx("example.okta.com"),
				},
			},
		}

		mockInvoker.EXPECT().
			ListIntegrationsV4(mock.Anything, mock.MatchedBy(func(params client.ListIntegrationsV4Params) bool {
				categoryParam, ok := params.Category.Get()
				return ok && len(categoryParam) == 1 && categoryParam[0] == "USER-INFORMATION"
			})).
			Return(&client.PublicApiListResponseIntegrationPublicV4Model{
				Items:      integrations,
				Pagination: client.PublicApiPaginationInfoModel{},
			}, nil)

		plan := tfsdk.Plan{
			Schema: d.getTestSchema(ctx),
		}

		diag := plan.Set(ctx, models.AponoUserInformationIntegrationsDataSourceModel{})
		require.False(t, diag.HasError(), "Error setting plan: %s", diag.Errors())

		req := datasource.ReadRequest{
			Config: tfsdk.Config{
				Schema: d.getTestSchema(ctx),
				Raw:    plan.Raw,
			},
		}

		resp := datasource.ReadResponse{
			State: tfsdk.State{
				Schema: d.getTestSchema(ctx),
				Raw:    tftypes.NewValue(d.getTestSchema(ctx).Type().TerraformType(ctx), nil),
			},
		}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Read returned error: %s", resp.Diagnostics.Errors())

		var state models.AponoUserInformationIntegrationsDataSourceModel
		resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
		require.False(t, resp.Diagnostics.HasError(), "Error getting state: %s", resp.Diagnostics.Errors())

		expectedModels := []models.UserInformationIntegrationModel{}
		for _, integration := range integrations {
			model, err := models.UserInformationIntegrationToModal(ctx, &integration)
			require.NoError(t, err, "Error converting integration to model")
			expectedModels = append(expectedModels, *model)
		}

		assert.Equal(t, expectedModels, state.Integrations, "Integrations do not match expected models")
	})
}

func (d *AponoUserInformationIntegrationsDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

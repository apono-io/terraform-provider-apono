package datasources

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &AponoUserInformationIntegrationsDataSource{}

func NewAponoUserInformationIntegrationsDataSource() datasource.DataSource {
	return &AponoUserInformationIntegrationsDataSource{}
}

type AponoUserInformationIntegrationsDataSource struct {
	client client.Invoker
}

func (d *AponoUserInformationIntegrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_information_integrations"
}

func (d *AponoUserInformationIntegrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of user information integrations, with optional filters by name and type. This data source is useful when you need to reference existing identity providers or context integrations like Google OAuth, Okta, PagerDuty, and others.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: `Filters the returned integrations by their name. Partial matching is supported with asterisks for contains, starts with, and ends with. (e.g., "Google\*").`,
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: `Filters the returned integrations by their type. Partial matching is supported with asterisks for contains, starts with, and ends with. (e.g., "\*duty\*").`,
				Optional:    true,
			},
			"integrations": schema.ListNestedAttribute{
				Description: "A list of user information integrations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the integration.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the integration.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of the integration.",
							Computed:    true,
						},
						"category": schema.StringAttribute{
							Description: "The integration’s category (e.g., USER-INFORMATION).",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The current operational status of the integration.",
							Computed:    true,
						},
						"last_sync_time": schema.StringAttribute{
							Description: "Timestamp of the last synchronization (if available).",
							Computed:    true,
							Optional:    true,
						},
						"integration_config": schema.MapAttribute{
							Description: "Integration-specific configuration that accepts key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"secret_store_config": schemas.GetSecretStoreConfigSchema(schemas.DataSourceMode),
					},
				},
			},
		},
	}
}

func (d *AponoUserInformationIntegrationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoUserInformationIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model models.AponoUserInformationIntegrationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !model.Name.IsNull() {
		name = model.Name.ValueString()
	}

	itergrationType := ""
	if !model.Type.IsNull() {
		itergrationType = model.Type.ValueString()
	}

	integrations, err := services.ListIntegrations(ctx, d.client, itergrationType, name, "", []string{common.UserInformationCategory})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user information integrations", err.Error())
		return
	}

	integrationModels := make([]models.UserInformationIntegrationModel, 0, len(integrations))
	for _, integration := range integrations {
		integrationModel, err := models.UserInformationIntegrationToModal(ctx, &integration)
		if err != nil {
			resp.Diagnostics.AddError("Error converting integration", err.Error())
			return
		}
		integrationModels = append(integrationModels, *integrationModel)
	}

	model.Integrations = integrationModels
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

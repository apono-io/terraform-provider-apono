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

var _ datasource.DataSource = &AponoUserInformationIntegrationsDataSource{}

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
		Description: "Retrieves a list of Apono User Information Integrations.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter integrations by name, supports wildcards.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Filter integrations by type, supports wildcards.",
				Optional:    true,
			},
			"integrations": schema.ListNestedAttribute{
				Description: "The list of user information integrations.",
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
							Description: "The category of the integration.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the integration.",
							Computed:    true,
						},
						"last_sync_time": schema.StringAttribute{
							Description: "The timestamp of the last synchronization.",
							Computed:    true,
							Optional:    true,
						},
						"integration_config": schema.MapAttribute{
							Description: "Configuration for the integration as key-value pairs.",
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

	typeName := ""
	if !model.Type.IsNull() {
		typeName = model.Type.ValueString()
	}

	integrations, err := services.ListIntegrations(ctx, d.client, typeName, name, []string{common.UserInformation})
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

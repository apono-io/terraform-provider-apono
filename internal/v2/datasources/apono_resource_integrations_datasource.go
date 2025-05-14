package datasources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSourceWithConfigure = &ResourceIntegrationsDataSource{}

type ResourceIntegrationsDataSource struct {
	client client.Invoker
}

func NewResourceIntegrationsDataSource() datasource.DataSource {
	return &ResourceIntegrationsDataSource{}
}

func (d *ResourceIntegrationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_integrations"
}

func (d *ResourceIntegrationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches resource integrations based on filters",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter by integration name. Supports wildcards (e.g., 'DB Prod*')",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Filter by integration type. Supports wildcards (e.g., 'postgresql')",
				Optional:    true,
			},
			"connector_id": schema.StringAttribute{
				Description: "Filter by connector ID",
				Optional:    true,
			},
			"integrations": schema.SetNestedAttribute{
				Description: "List of matching integrations",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the integration.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Human-readable name of the integration.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: `Type of the integration (e.g., "aws-account", "postgresql").`,
							Computed:    true,
						},
						"connector_id": schema.StringAttribute{
							Description: "ID of the Apono Connector used for the integration.",
							Computed:    true,
						},
						"connected_resource_types": schema.ListAttribute{
							Description: "List of resource types discovered by the integration.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"integration_config": schema.MapAttribute{
							MarkdownDescription: "Integration-specific configuration key-value pairs.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"secret_store_config": schemas.GetSecretStoreConfigSchema(schemas.DataSourceMode),
						"custom_access_details": schema.StringAttribute{
							Description: "Custom access instructions for end users.",
							Computed:    true,
						},
						"owner":          schemas.GetOwnerSchema(schemas.DataSourceMode),
						"owners_mapping": schemas.GetOwnersMappingSchema(schemas.DataSourceMode),
					},
				},
			},
		},
	}
}

func (d *ResourceIntegrationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *ResourceIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.ResourceIntegrationsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	integrationType := ""
	if !config.Type.IsNull() {
		integrationType = config.Type.ValueString()
	}
	connectorID := ""
	if !config.ConnectorID.IsNull() {
		connectorID = config.ConnectorID.ValueString()
	}

	tflog.Debug(ctx, "Reading resource integrations", map[string]any{
		"name_filter":         name,
		"type_filter":         integrationType,
		"connector_id_filter": connectorID,
	})

	integrations, err := services.ListIntegrations(ctx, d.client, integrationType, name, connectorID, []string{common.ResourceCategory})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving resource integrations", fmt.Sprintf("Could not retrieve resource integrations: %v", err))
		return
	}

	integrationsModel, err := models.ResourceIntegrationsToModel(ctx, integrations)
	if err != nil {
		resp.Diagnostics.AddError("Error converting resource integrations", fmt.Sprintf("Could not convert resource integrations: %v", err))
		return
	}
	config.Integrations = integrationsModel.Integrations

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Resource integrations retrieved successfully", map[string]any{
		"count": len(config.Integrations),
	})
}

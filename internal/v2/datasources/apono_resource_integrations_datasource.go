package datasources

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ResourceIntegrationsDataSource{}

type ResourceIntegrationsDataSource struct {
}

type ResourceIntegrationsDataSourceModel struct {
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	ConnectorID  types.String `tfsdk:"connector_id"`
	Integrations types.Set    `tfsdk:"integrations"`
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

func (d *ResourceIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ResourceIntegrationsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement actual API call to fetch integrations based on filters
	// This is a placeholder implementation that sets an empty list

	// Set empty integrations set as placeholder
	emptyIntegrations, diags := types.SetValueFrom(ctx, types.ObjectType{}, []map[string]interface{}{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Integrations = emptyIntegrations

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSourceWithConfigure = &AponoBundlesDataSource{}

func NewAponoBundlesDataSource() datasource.DataSource {
	return &AponoBundlesDataSource{}
}

type AponoBundlesDataSource struct {
	client client.Invoker
}

func (d *AponoBundlesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundles"
}

func (d *AponoBundlesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of existing Apono Bundles. This data source is typically used to reference bundle definitions within Access Flow resources. You can filter bundles by name using exact or wildcard matching.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: `Filter by bundle name. Partial matching is supported with asterisks for contains, starts with, and ends with. (e.g., "prod*"). Matching is case-insensitive.`,
				Optional:    true,
			},
			"bundles": schema.ListNestedAttribute{
				Description: "A list of bundles matching the filter.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the bundle.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the bundle.",
							Computed:    true,
						},
						"access_targets": schema.ListNestedAttribute{
							Description: "A list of access targets included in the bundle.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"integration":  schemas.GetIntegrationTargetSchema(schemas.DataSourceMode),
									"access_scope": schemas.GetAccessScopeTargetSchema(schemas.DataSourceMode),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *AponoBundlesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoBundlesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.BundlesDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	tflog.Debug(ctx, "Reading bundles", map[string]any{
		"name_filter": name,
	})

	bundles, err := services.ListBundles(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving bundles", fmt.Sprintf("Could not retrieve bundles: %v", err))
		return
	}

	bundleModels, err := models.BundlesResponseToModels(ctx, bundles)
	if err != nil {
		resp.Diagnostics.AddError("Error converting bundles", fmt.Sprintf("Could not convert bundles: %v", err))
		return
	}
	config.Bundles = bundleModels

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Bundles retrieved successfully", map[string]any{
		"count": len(config.Bundles),
	})
}

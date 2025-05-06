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

var _ datasource.DataSource = &AponoBundlesDataSource{}

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
		Description: "Retrieves existing bundles. Use this data source to reference bundles in the Access Flow resource.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: `Filters the returned bundles by their name. Partial matching is supported with asterisks for contains, starts with, and ends with.  (e.g., "\*my-bundles\*").`,
				Optional:    true,
			},
			"bundles": schema.SetNestedAttribute{
				Description: "A list of bundles that match the filter.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the bundle.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the bundle.",
							Computed:    true,
						},
						"access_targets": schema.SetNestedAttribute{
							Description: "List of access targets for this bundle",
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
	var config models.BundlesV2DataModel
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

package datasources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = &AponoSpaceScopesDataSource{}

func NewAponoSpaceScopesDataSource() datasource.DataSource {
	return &AponoSpaceScopesDataSource{}
}

type AponoSpaceScopesDataSource struct {
	client client.Invoker
}

func (d *AponoSpaceScopesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_scopes"
}

func (d *AponoSpaceScopesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves existing Apono Space Scopes. Use this data source to reference existing space scopes when creating or updating spaces.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filters the returned space scopes by their name. Partial matching is supported with asterisks for contains, starts with, and ends with. Matching is case-insensitive.",
				Optional:    true,
			},
			"space_scopes": schema.ListNestedAttribute{
				Description: "A list of space scopes that match the filter.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the space scope.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the space scope.",
							Computed:    true,
						},
						"query": schema.StringAttribute{
							Description: "The full query string that is used to define the space scope.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AponoSpaceScopesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoSpaceScopesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.SpaceScopesDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	spaceScopes, err := services.ListSpaceScopes(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving space scopes", fmt.Sprintf("Could not retrieve space scopes: %v", err))
		return
	}

	config.SpaceScopes = models.SpaceScopesToDataModels(spaceScopes)

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

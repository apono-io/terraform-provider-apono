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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &AponoSpacesDataSource{}

func NewAponoSpacesDataSource() datasource.DataSource {
	return &AponoSpacesDataSource{}
}

type AponoSpacesDataSource struct {
	client client.Invoker
}

func (d *AponoSpacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spaces"
}

func (d *AponoSpacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves existing Apono Spaces. Use this data source to reference spaces in the Access Flow, Bundle, or Access Scope resources.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filters the returned spaces by their name. Partial matching is supported with asterisks for contains, starts with, and ends with. Matching is case-insensitive.",
				Optional:    true,
			},
			"spaces": schema.ListNestedAttribute{
				Description: "A list of spaces that match the filter.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the space.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Display name of the space.",
							Computed:    true,
						},
						"space_scope_references": schema.ListAttribute{
							Description: "Names of space scopes assigned to this space.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *AponoSpacesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoSpacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config models.SpacesDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	spaces, err := services.ListSpaces(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving spaces", fmt.Sprintf("Could not retrieve spaces: %v", err))
		return
	}

	config.Spaces = models.SpacesToDataModels(spaces)

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

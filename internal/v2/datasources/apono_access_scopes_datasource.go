package datasources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &AponoAccessScopesDataSource{}

func NewAponoAccessScopesDataSource() datasource.DataSource {
	return &AponoAccessScopesDataSource{}
}

type AponoAccessScopesDataSource struct {
	client client.Invoker
}

type accessScopesDataSourceModel struct {
	Name         types.String              `tfsdk:"name"`
	AccessScopes []common.AccessScopeModel `tfsdk:"access_scopes"`
}

func (d *AponoAccessScopesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_scopes"
}

func (d *AponoAccessScopesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of Apono Access Scopes.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter access scopes by name, supports wildcards.",
				Optional:    true,
			},
			"access_scopes": schema.ListNestedAttribute{
				Description: "The list of access scopes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the access scope.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the access scope.",
							Computed:    true,
						},
						"query": schema.StringAttribute{
							Description: "The query expression for the access scope.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AponoAccessScopesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoAccessScopesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config accessScopesDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	tflog.Debug(ctx, "Reading access scopes", map[string]any{
		"name_filter": name,
	})

	accessScopes, err := common.ListAccessScopesByName(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving access scopes", fmt.Sprintf("Could not retrieve access scopes: %v", err))
		return
	}

	config.AccessScopes = common.AccessScopesToModels(accessScopes)

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access scopes retrieved successfully", map[string]any{
		"count": len(config.AccessScopes),
	})
}

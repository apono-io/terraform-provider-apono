package datasources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
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
	Name         types.String                `tfsdk:"name"`
	AccessScopes []services.AccessScopeModel `tfsdk:"access_scopes"`
}

func (d *AponoAccessScopesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_scopes"
}

func (d *AponoAccessScopesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves existing Apono Access Scopes. This data source can be used to feed existing access scopes into the Access Flow resource.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filters the returned access scopes by their name. Partial matching is supported with asterisks for contains, starts with, and ends with.",
				Optional:    true,
			},
			"access_scopes": schema.SetNestedAttribute{
				Description: "A set of access scopes that match the specified criteria.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the Apono Access Scope.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the Apono Access Scope.",
							Computed:    true,
						},
						"query": schema.StringAttribute{
							Description: "The full query string that is used to define the access scope.",
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

	accessScopes, err := services.ListAccessScopesByName(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving access scopes", fmt.Sprintf("Could not retrieve access scopes: %v", err))
		return
	}

	config.AccessScopes = services.AccessScopesToModels(accessScopes)

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access scopes retrieved successfully", map[string]any{
		"count": len(config.AccessScopes),
	})
}

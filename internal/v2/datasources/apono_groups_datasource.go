package datasources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &AponoGroupsDataSource{}

func NewAponoGroupsDataSource() datasource.DataSource {
	return &AponoGroupsDataSource{}
}

type AponoGroupsDataSource struct {
	client client.Invoker
}

func (d *AponoGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *AponoGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of Apono Groups.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filter groups by name, supports wildcards.",
				Optional:    true,
			},
			"source_integration": schema.StringAttribute{
				Description: "Filter groups by source integration name or ID.",
				Optional:    true,
			},
			"groups": schema.SetNestedAttribute{
				Description: "The list of groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the group.",
							Computed:    true,
						},
						"source_integration_id": schema.StringAttribute{
							Description: "The source integration ID.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *AponoGroupsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	common.ConfigureDataSourceClientInvoker(ctx, req, resp, &d.client)
}

func (d *AponoGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config services.GroupsDataModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := ""
	if !config.Name.IsNull() {
		name = config.Name.ValueString()
	}

	sourceIntegration := ""
	if !config.SourceIntegration.IsNull() {
		sourceIntegration = config.SourceIntegration.ValueString()
	}

	tflog.Debug(ctx, "Reading groups", map[string]any{
		"name_filter":        name,
		"source_integration": sourceIntegration,
	})

	allGroups, err := services.ListGroups(ctx, d.client, name)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving groups", fmt.Sprintf("Could not retrieve groups: %v", err))
		return
	}

	filteredGroups := services.FilterGroupsBySourceIntegration(allGroups, sourceIntegration)

	var groupModels []services.GroupDataModel
	for _, group := range filteredGroups {
		groupModels = append(groupModels, services.GroupToDataModel(&group))
	}

	config.Groups = groupModels

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Groups retrieved successfully", map[string]any{
		"count": len(config.Groups),
	})
}

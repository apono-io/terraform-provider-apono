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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSourceWithConfigure = &AponoGroupsDataSource{}

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
		Description: "Retrieves existing groups, Apono-managed and IDP-managed groups. Use this data source to reference groups in the Access Flow resource.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Filters the returned groups by their name. Partial matching is supported with asterisks for contains, starts with, and ends with.",
				Optional:    true,
			},
			"source_integration": schema.StringAttribute{
				Description: "Filters the returned groups by their name or IDs. Partial matching is supported for names with asterisks for contains, starts with, and ends with.",
				Optional:    true,
			},
			"groups": schema.SetNestedAttribute{
				Description: "A set of groups that match the filter.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Group display name.",
							Computed:    true,
						},
						"source_integration_id": schema.StringAttribute{
							Description: "ID of the IDP integration from which the group originated, or null.",
							Computed:    true,
						},
						"source_integration_name": schema.StringAttribute{
							Description: "Human‑readable name of the originating IDP integration, or null.",
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
	var config models.GroupsDataModel
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

	var groupModels []models.GroupDataModel
	for _, group := range filteredGroups {
		groupModels = append(groupModels, models.GroupToDataModel(&group))
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

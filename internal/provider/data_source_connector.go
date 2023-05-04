package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &connectorDataSource{}

func NewConnectorDataSource() datasource.DataSource {
	return &connectorDataSource{}
}

// connectorDataSource defines the data source implementation for Apono connector.
type connectorDataSource struct {
	provider *AponoProvider
}

// connectorDataSourceModel describes the data source data model.
type connectorDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *connectorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector"
}

func (d *connectorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the ID of an Apono connector for when creating integrations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Connector identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Connector name",
				Computed:            true,
			},
		},
	}
}

func (d *connectorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.provider, resp.Diagnostics = toProvider(req.ProviderData)
}

func (d *connectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data connectorDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Looking for connector")

	connectors, _, err := d.provider.client.ConnectorsApi.ListConnectors(ctx).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get conector, error: %s, body: %s", apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get conector: %s", err.Error()))
		}

		return
	}

	for _, connector := range connectors {
		if connector.GetConnectorId() != data.Id.ValueString() {
			continue
		}

		data.Id = types.StringValue(connector.GetConnectorId())
		data.Name = types.StringValue(connector.GetConnectorId())

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	resp.Diagnostics.AddError("Not Found", "No connector matched the search criteria")
}

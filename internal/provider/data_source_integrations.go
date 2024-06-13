package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &integrationsDataSource{}

func NewIntegrationsDataSource() datasource.DataSource {
	return &integrationsDataSource{}
}

// integrationsDataSource defines the data source implementation for Apono connector.
type integrationsDataSource struct {
	provider *AponoProvider
}

// integrationsDataSourceModel describes the data source data model.
type integrationsDataSourceModel struct {
	ID           types.String              `tfsdk:"id"`
	Type         types.String              `tfsdk:"type"`
	ConnectorID  types.String              `tfsdk:"connector_id"`
	Integrations []models.IntegrationModel `tfsdk:"integrations"`
}

func (d *integrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integrations"
}

func (d *integrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get list integrations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Integration type to filter by",
				Optional:            true,
			},
			"connector_id": schema.StringAttribute{
				MarkdownDescription: "Integration connector id to filter by",
				Optional:            true,
			},
			"integrations": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of integrations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: IntegrationDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *integrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.provider, resp.Diagnostics = toProvider(req.ProviderData)
}

func (d *integrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model integrationsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Read integrations")
	response, _, err := d.provider.terraformClient.IntegrationsAPI.TfListIntegrationsV1(ctx).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to list integrations, error: %s, body: %s", apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to list integrations: %s", err.Error()))
		}

		return
	}

	model.ID = types.StringValue(id.UniqueId())
	for _, integration := range response.Data {
		if !model.Type.IsNull() && integration.Type != model.Type.ValueString() {
			continue
		}

		if !model.ConnectorID.IsNull() && *integration.ProvisionerId.Get() != model.ConnectorID.ValueString() {
			continue
		}

		m, diagnostics := models.ConvertToIntegrationModel(ctx, &integration)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}

		model.Integrations = append(model.Integrations, *m)
	}

	// Save state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func IntegrationDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Integration identifier",
		},
		"name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Integration name",
		},
		"type": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Integration type",
		},
		"connector_id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Apono connector identifier",
		},
		"connected_resource_types": schema.SetAttribute{
			MarkdownDescription: "Connected resource types",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"custom_access_details": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Custom access details message that will be displayed to end users when they access this integration",
		},
		"metadata": schema.MapAttribute{
			MarkdownDescription: "Integration metadata",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"aws_secret": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					MarkdownDescription: "Aws secret region",
					Required:            true,
				},
				"secret_id": schema.StringAttribute{
					MarkdownDescription: "Aws secret name or ARN",
					Required:            true,
				},
			},
		},
		"gcp_secret": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"project": schema.StringAttribute{
					MarkdownDescription: "GCP secret project",
					Required:            true,
				},
				"secret_id": schema.StringAttribute{
					MarkdownDescription: "GCP secret ID",
					Required:            true,
				},
			},
		},
		"kubernetes_secret": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"namespace": schema.StringAttribute{
					MarkdownDescription: "Kubernetes secret namespace",
					Required:            true,
				},
				"name": schema.StringAttribute{
					MarkdownDescription: "Kubernetes secret name",
					Required:            true,
				},
			},
		},
		"resource_owner_mappings": schema.ListNestedAttribute{
			MarkdownDescription: "List of resource-to-owner-mappings. Used to map resource owner to apono owner.",
			Required:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"tag_name": schema.StringAttribute{
						Required: true,
					},
					"attribute_type": schema.StringAttribute{
						Required: true,
					},
					"attribute_integration_id": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
		"integration_owners": schema.SingleNestedAttribute{
			MarkdownDescription: "List of integration owner. Each item defines owner of the integration.",
			Required:            true,
			Attributes: map[string]schema.Attribute{
				"owners": schema.ListNestedAttribute{
					Required: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"integration_id": schema.StringAttribute{
								Required: true,
							},
							"attribute_type_id": schema.StringAttribute{
								Required: true,
							},
							"attribute_value": schema.ListAttribute{
								ElementType: types.StringType,
								Required:    true,
							},
						},
					},
				},
			},
		},
	}
}

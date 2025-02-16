package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/schemas"
	"github.com/apono-io/terraform-provider-apono/internal/services"
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

		if !model.ConnectorID.IsNull() {
			provisionerId := integration.ProvisionerId.Get()
			if provisionerId == nil {
				continue
			}
			if *provisionerId != model.ConnectorID.ValueString() {
				continue
			}
		}

		m, diagnostics := services.ConvertToIntegrationModel(ctx, &integration)
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
		"hashicorp_vault_secret": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"secret_engine": schema.StringAttribute{
					MarkdownDescription: "Hashicorp Vault Secret Engine",
					Required:            true,
				},
				"path": schema.StringAttribute{
					MarkdownDescription: "Hashicorp Vault secret path",
					Required:            true,
				},
			},
		},
		"resource_owner_mappings": schema.SetNestedAttribute{
			MarkdownDescription: "Let Apono know which tag represents owners and how to map it to a known attribute in Apono.",
			Computed:            true,
			NestedObject:        schemas.DataSourceResourceOwnerMapping,
		},
		"integration_owners": schema.SetNestedAttribute{
			MarkdownDescription: "Enter one or more users, groups, shifts or attributes. This field is mandatory when using Resource Owners and serves as a fallback approver if no resource owner is found.",
			Computed:            true,
			NestedObject:        schemas.DataSourceIntegrationOwner,
		},
	}
}

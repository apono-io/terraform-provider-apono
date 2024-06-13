package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"

	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/slices"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &integrationResource{}
var _ resource.ResourceWithImportState = &integrationResource{}
var _ resource.ResourceWithValidateConfig = &integrationResource{}

var (
	secretTypeAttributeNames = map[string]string{
		"AWS":        "aws_secret",
		"GCP":        "gcp_secret",
		"KUBERNETES": "kubernetes_secret",
	}
)

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

// integrationResource defines the resource implementation.
type integrationResource struct {
	provider *AponoProvider
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Apono Integration",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Integration identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Integration name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Integration type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_id": schema.StringAttribute{
				MarkdownDescription: "Apono connector identifier",
				Required:            true,
			},
			"connected_resource_types": schema.SetAttribute{
				MarkdownDescription: "Resource types to sync, if omitted all resources types will be synced.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_access_details": schema.StringAttribute{
				MarkdownDescription: "Custom access details message that will be displayed to end users when they access this integration.",
				Optional:            true,
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
		},
	}
}

func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.provider, resp.Diagnostics = toProvider(req.ProviderData)
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *models.IntegrationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	metadata := map[string]interface{}{}
	for name, value := range data.Metadata.Elements() {
		metadata[name] = utils.AttrValueToString(value)
	}

	var connectedResourceTypes []string
	if !data.ConnectedResourceTypes.IsNull() {
		for _, resourceType := range data.ConnectedResourceTypes.Elements() {
			connectedResourceTypes = append(connectedResourceTypes, utils.AttrValueToString(resourceType))
		}
	}

	var secretConfig map[string]interface{}
	if data.AwsSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "AWS",
			"region":    data.AwsSecret.Region.ValueString(),
			"secret_id": data.AwsSecret.SecretID.ValueString(),
		}
	} else if data.GcpSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "GCP",
			"project":   data.GcpSecret.Project.ValueString(),
			"secret_id": data.GcpSecret.SecretID.ValueString(),
		}
	} else if data.KubernetesSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "KUBERNETES",
			"namespace": data.KubernetesSecret.Namespace.ValueString(),
			"name":      data.KubernetesSecret.Name.ValueString(),
		}
	}

	connectorID := data.ConnectorID.ValueString()
	integration, _, err := r.provider.terraformClient.IntegrationsAPI.TfCreateIntegrationV1(ctx).
		UpsertIntegrationTerraform(aponoapi.UpsertIntegrationTerraform{
			Name:                   data.Name.ValueString(),
			Type:                   data.Type.ValueString(),
			ProvisionerId:          *aponoapi.NewNullableString(&connectorID),
			Params:                 metadata,
			SecretConfig:           secretConfig,
			ConnectedResourceTypes: connectedResourceTypes,
			CustomAccessDetails:    data.CustomAccessDetails.ValueString(),
			IntegrationOwners:      integrationOwnersToModel(data.IntegrationOwners),
			ResourceOwnersMappings: convertMappingsArrayToModel(data.ResourceOwnerMappings),
		}).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "create", "integration", "")
		resp.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := models.ConvertToIntegrationModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	tflog.Debug(ctx, "Created integration", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *models.IntegrationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	integration, _, err := r.provider.terraformClient.IntegrationsAPI.TfGetIntegrationV1(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "integration", data.ID.ValueString())
		resp.Diagnostics.Append(diagnostics...)

		return
	}
	model, diagnostics := models.ConvertToIntegrationModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *models.IntegrationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	metadata := map[string]interface{}{}
	for name, value := range data.Metadata.Elements() {
		metadata[name] = utils.AttrValueToString(value)
	}

	var connectedResourceTypes []string
	if !data.ConnectedResourceTypes.IsNull() {
		for _, resourceType := range data.ConnectedResourceTypes.Elements() {
			connectedResourceTypes = append(connectedResourceTypes, utils.AttrValueToString(resourceType))
		}
	}

	var secretConfig map[string]interface{}
	if data.AwsSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "AWS",
			"region":    data.AwsSecret.Region.ValueString(),
			"secret_id": data.AwsSecret.SecretID.ValueString(),
		}
	} else if data.GcpSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "GCP",
			"project":   data.GcpSecret.Project.ValueString(),
			"secret_id": data.GcpSecret.SecretID.ValueString(),
		}
	} else if data.KubernetesSecret != nil {
		secretConfig = map[string]interface{}{
			"type":      "KUBERNETES",
			"namespace": data.KubernetesSecret.Namespace.ValueString(),
			"name":      data.KubernetesSecret.Name.ValueString(),
		}
	}

	connectorID := data.ConnectorID.ValueString()
	integration, _, err := r.provider.terraformClient.IntegrationsAPI.TfUpdateIntegrationV1(ctx, data.ID.ValueString()).
		UpsertIntegrationTerraform(aponoapi.UpsertIntegrationTerraform{
			Name:                   data.Name.ValueString(),
			Type:                   data.Type.ValueString(),
			ProvisionerId:          *aponoapi.NewNullableString(&connectorID),
			Params:                 metadata,
			SecretConfig:           secretConfig,
			ConnectedResourceTypes: connectedResourceTypes,
			CustomAccessDetails:    data.CustomAccessDetails.ValueString(),
			IntegrationOwners:      integrationOwnersToModel(data.IntegrationOwners),
			ResourceOwnersMappings: convertMappingsArrayToModel(data.ResourceOwnerMappings),
		}).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "update", "integration", data.ID.ValueString())
		resp.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := models.ConvertToIntegrationModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	tflog.Debug(ctx, "Updated integration", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *models.IntegrationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	messageResponse, _, err := r.provider.client.IntegrationsApi.DeleteIntegrationV2(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "delete", "integration", data.ID.ValueString())
		resp.Diagnostics.Append(diagnostics...)

		return
	}

	tflog.Debug(ctx, "Deleted integration", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": messageResponse.GetMessage(),
	})
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	integrationId := req.ID
	tflog.Debug(ctx, "importing integration", map[string]interface{}{
		"id": integrationId,
	})

	integration, _, err := r.provider.terraformClient.IntegrationsAPI.TfGetIntegrationV1(ctx, integrationId).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "integration", integrationId)
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := models.ConvertToIntegrationModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	// Save imported data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Imported integration", map[string]interface{}{
		"id": integrationId,
	})
}

func (r *integrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.provider == nil {
		return
	}

	var model models.IntegrationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, _, err := r.provider.client.IntegrationsApi.GetIntegrationConfig(ctx, model.Type.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get config: %s", err.Error()))
		return
	}

	if config.GetRequiresSecret() {
		var supportedSecretExpressions []path.Expression
		for _, secretType := range config.SupportedSecretTypes {
			if attributeName, ok := secretTypeAttributeNames[secretType]; ok {
				supportedSecretExpressions = append(supportedSecretExpressions, path.MatchRoot(attributeName))
			}
		}

		resourcevalidator.ExactlyOneOf(supportedSecretExpressions...).ValidateResource(ctx, req, resp)
	}

	metadataElements := model.Metadata.Elements()
	for _, param := range config.GetParams() {
		paramName := param.GetId()
		paramPossibleValues := param.GetValues()
		paramDefaultValue := param.GetDefault()
		paramIsOptional := param.GetOptional()

		attributePath := path.Root(fmt.Sprintf(`metadata["%s"]`, paramName))

		metadataValue, hasValue := metadataElements[paramName]
		if !hasValue {
			if paramIsOptional && paramDefaultValue == "" {
				continue
			}
			if paramDefaultValue != "" {
				resp.Diagnostics.AddAttributeError(
					attributePath,
					"Missing Configuration for Required Attribute",
					fmt.Sprintf("Must set a configuration value for the %s attribute as the provider has marked it as required.\n\n", attributePath.String())+
						"Refer to the provider documentation or contact the provider developers for additional information about configurable attributes that are required.\n\n"+
						fmt.Sprintf("Configuring this integration through the UI will use default value: %s", paramDefaultValue),
				)
			} else {
				resp.Diagnostics.AddAttributeError(
					attributePath,
					"Missing Configuration for Required Attribute",
					fmt.Sprintf("Must set a configuration value for the %s attribute as the provider has marked it as required.\n\n", attributePath.String())+
						"Refer to the provider documentation or contact the provider developers for additional information about configurable attributes that are required.",
				)
			}

			continue
		}

		metadataValueStr := utils.AttrValueToString(metadataValue)
		if len(paramPossibleValues) > 0 && !slices.Contains(paramPossibleValues, metadataValueStr) {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				attributePath,
				fmt.Sprintf("value must be one of: %q", paramPossibleValues),
				metadataValueStr,
			))
		}
	}
}

func convertMappingsArrayToModel(m []*models.ResourceOwnerMapping) []aponoapi.ResourceOwnerMappingTerraform {
	result := make([]aponoapi.ResourceOwnerMappingTerraform, len(m))
	for i, rom := range m {
		result[i] = resourceOwnerMappingToModel(rom)
	}
	return result
}

func resourceOwnerMappingToModel(rom *models.ResourceOwnerMapping) aponoapi.ResourceOwnerMappingTerraform {
	return aponoapi.ResourceOwnerMappingTerraform{
		TagName:                rom.TagName.ValueString(),
		AttributeType:          rom.AttributeType.ValueString(),
		AttributeIntegrationId: *aponoapi.NewNullableString(rom.AttributeIntegrationId.ValueStringPointer()),
	}
}

func integrationOwnersToModel(IntegrationOwners *models.IntegrationOwners) aponoapi.IntegrationOwnersTerraform {
	owners := make([]aponoapi.IntegrationOwnerTerraform, len(IntegrationOwners.Owners))
	for i, o := range IntegrationOwners.Owners {
		owners[i] = integrationOwnerToModel(o)
	}
	return aponoapi.IntegrationOwnersTerraform{
		Owners: owners,
	}
}

func integrationOwnerToModel(owner *models.IntegrationOwner) aponoapi.IntegrationOwnerTerraform {
	return aponoapi.IntegrationOwnerTerraform{
		IntegrationId:   *aponoapi.NewNullableString(owner.IntegrationId.ValueStringPointer()),
		AttributeTypeId: owner.AttributeTypeId.ValueString(),
		AttributeValue:  convertStringArray(owner.AttributeValue),
	}
}

func convertStringArray(a []types.String) []string {
	result := make([]string, len(a))
	for i, item := range a {
		result[i] = item.ValueString()
	}
	return result
}

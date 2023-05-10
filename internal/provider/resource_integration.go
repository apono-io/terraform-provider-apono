package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

	awsSecretAttrMap = map[string]attr.Type{
		"region":     basetypes.StringType{},
		"secret_arn": basetypes.StringType{},
	}
	gcpSecretAttrMap = map[string]attr.Type{
		"project":   basetypes.StringType{},
		"secret_id": basetypes.StringType{},
	}
	kubernetesSecretAttrMap = map[string]attr.Type{
		"namespace": basetypes.StringType{},
		"name":      basetypes.StringType{},
	}
)

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

// integrationResource defines the resource implementation.
type integrationResource struct {
	provider *AponoProvider
}

// integrationResourceModel describes the resource data model.
type integrationResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	ConnectorID      types.String `tfsdk:"connector_id"`
	Metadata         types.Map    `tfsdk:"metadata"`
	AwsSecret        types.Object `tfsdk:"aws_secret"`
	GcpSecret        types.Object `tfsdk:"gcp_secret"`
	KubernetesSecret types.Object `tfsdk:"kubernetes_secret"`
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
			"metadata": schema.MapAttribute{
				MarkdownDescription: "Integration metadata",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"aws_secret": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
					},
					"secret_arn": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
					},
				},
			},
			"gcp_secret": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
					},
					"secret_id": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
					},
				},
			},
			"kubernetes_secret": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"namespace": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Example configurable attribute",
						Required:            true,
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
	var data *integrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	metadata := map[string]interface{}{}
	for name, value := range data.Metadata.Elements() {
		metadata[name] = attrValueToString(value)
	}

	var secretConfig map[string]interface{}
	if !data.AwsSecret.IsNull() {
		attributes := data.AwsSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "AWS",
			"region":    attrValueToString(attributes["region"]),
			"secret_id": attrValueToString(attributes["secret_arn"]),
		}
	} else if !data.GcpSecret.IsNull() {
		attributes := data.GcpSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "GCP",
			"project":   attrValueToString(attributes["project"]),
			"secret_id": attrValueToString(attributes["secret_id"]),
		}
	} else if !data.KubernetesSecret.IsNull() {
		attributes := data.KubernetesSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "KUBERNETES",
			"namespace": attrValueToString(attributes["namespace"]),
			"name":      attrValueToString(attributes["name"]),
		}
	}

	connectorID := data.ConnectorID.ValueString()
	integration, _, err := r.provider.client.IntegrationsApi.CreateIntegrationV2(ctx).
		CreateIntegration(apono.CreateIntegration{
			Name:          data.Name.ValueString(),
			Type:          data.Type.ValueString(),
			ProvisionerId: *apono.NewNullableString(&connectorID),
			Metadata:      metadata,
			SecretConfig:  secretConfig,
		}).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create integration, error: %s, body: %s", apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create integration: %s", err.Error()))
		}

		return
	}

	model, diagnostics := r.convertToModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	tflog.Trace(ctx, "created integration", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *integrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *integrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	metadata := map[string]interface{}{}
	for name, value := range data.Metadata.Elements() {
		metadata[name] = attrValueToString(value)
	}

	var secretConfig map[string]interface{}
	if !data.AwsSecret.IsNull() {
		attributes := data.AwsSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "AWS",
			"region":    attrValueToString(attributes["region"]),
			"secret_id": attrValueToString(attributes["secret_arn"]),
		}
	} else if !data.GcpSecret.IsNull() {
		attributes := data.GcpSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "GCP",
			"project":   attrValueToString(attributes["project"]),
			"secret_id": attrValueToString(attributes["secret_id"]),
		}
	} else if !data.KubernetesSecret.IsNull() {
		attributes := data.KubernetesSecret.Attributes()
		secretConfig = map[string]interface{}{
			"type":      "KUBERNETES",
			"namespace": attrValueToString(attributes["namespace"]),
			"name":      attrValueToString(attributes["name"]),
		}
	}

	connectorID := data.ConnectorID.ValueString()
	integration, _, err := r.provider.client.IntegrationsApi.UpdateIntegrationV2(ctx, data.ID.ValueString()).
		UpdateIntegration(apono.UpdateIntegration{
			Name:          data.Name.ValueString(),
			ProvisionerId: *apono.NewNullableString(&connectorID),
			Metadata:      metadata,
			SecretConfig:  secretConfig,
		}).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update integration with id %s, error: %s, body: %s", data.ID, apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update integration with id %s: %s", data.ID, err.Error()))
		}

		return
	}

	model, diagnostics := r.convertToModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	tflog.Trace(ctx, "updated integration", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *integrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	messageResponse, _, err := r.provider.client.IntegrationsApi.DeleteIntegrationV2(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete integration with id %s, error: %s, body: %s", data.ID, apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete integration with id %s: %s", data.ID, err.Error()))
		}

		return
	}

	tflog.Debug(ctx, "deleted integration", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": messageResponse.GetMessage(),
	})
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	integrationId := req.ID
	tflog.Debug(ctx, "importing integration", map[string]interface{}{
		"id": integrationId,
	})

	integration, _, err := r.provider.client.IntegrationsApi.GetIntegrationV2(ctx, integrationId).
		Execute()
	if err != nil {
		if apiError, ok := err.(*apono.GenericOpenAPIError); ok {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration with id %s, error: %s, body: %s", integrationId, apiError.Error(), string(apiError.Body())))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration with id %s: %s", integrationId, err.Error()))
		}

		return
	}

	model, diagnostics := r.convertToModel(ctx, integration)
	if len(diagnostics) > 0 {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	// Save imported data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "imported integration", map[string]interface{}{
		"id": integrationId,
	})
}

func (r *integrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if r.provider == nil {
		return
	}

	var model integrationResourceModel
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

		attributePath := path.Root(fmt.Sprintf(`metadata["%s"]`, paramName))

		metadataValue, hasValue := metadataElements[paramName]
		if !hasValue {
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

		metadataValueStr := attrValueToString(metadataValue)
		if len(paramPossibleValues) > 0 && !slices.Contains(paramPossibleValues, metadataValueStr) {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				attributePath,
				fmt.Sprintf("value must be one of: %q", paramPossibleValues),
				metadataValueStr,
			))
		}
	}
}

func (r *integrationResource) convertToModel(ctx context.Context, integration *apono.Integration) (*integrationResourceModel, diag.Diagnostics) {
	metadataMapValue, diagnostics := types.MapValueFrom(ctx, types.StringType, integration.GetMetadata())

	data := integrationResourceModel{}
	data.ID = types.StringValue(integration.GetId())
	data.Name = types.StringValue(integration.GetName())
	data.Type = types.StringValue(integration.GetType())
	data.ConnectorID = types.StringValue(integration.GetProvisionerId())
	data.Metadata = metadataMapValue

	data.AwsSecret = types.ObjectNull(awsSecretAttrMap)
	data.GcpSecret = types.ObjectNull(gcpSecretAttrMap)
	data.KubernetesSecret = types.ObjectNull(kubernetesSecretAttrMap)

	secretConfig := integration.GetSecretConfig()
	switch secretConfig["type"] {
	case "AWS":
		secretAttributes := map[string]attr.Value{
			"region":     basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["region"])),
			"secret_arn": basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["secret_id"])),
		}
		data.AwsSecret, diagnostics = types.ObjectValue(awsSecretAttrMap, secretAttributes)
	case "GCP":
		secretAttributes := map[string]attr.Value{
			"project":   basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["project"])),
			"secret_id": basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["secret_id"])),
		}
		data.GcpSecret, diagnostics = types.ObjectValue(gcpSecretAttrMap, secretAttributes)
	case "KUBERNETES":
		secretAttributes := map[string]attr.Value{
			"namespace": basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["namespace"])),
			"name":      basetypes.NewStringValue(fmt.Sprintf("%v", secretConfig["name"])),
		}
		data.KubernetesSecret, diagnostics = types.ObjectValue(kubernetesSecretAttrMap, secretAttributes)
	}

	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	return &data, nil
}

func attrValueToString(val attr.Value) string {
	switch value := val.(type) {
	case types.String:
		return value.ValueString()
	default:
		return value.String()
	}
}

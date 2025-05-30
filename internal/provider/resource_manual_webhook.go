package provider

import (
	"context"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/services"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ManualWebhookResource{}
var _ resource.ResourceWithImportState = &ManualWebhookResource{}

func NewWebhookResource() resource.Resource {
	return &ManualWebhookResource{}
}

// ManualWebhookResource defines the resource implementation.
type ManualWebhookResource struct {
	provider *AponoProvider
}

func (w ManualWebhookResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_manual_webhook"
}

func (w ManualWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	var allowedHttpMethods []string
	for _, method := range aponoapi.AllowedWebhookMethodTerraformModelEnumValues {
		allowedHttpMethods = append(allowedHttpMethods, string(method))
	}

	response.Schema = schema.Schema{
		MarkdownDescription: "Apono Manual Webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Manual Webhook identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Manual Webhook name. This is a human-readable label to identify the webhook",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the trigger is active. Set to true to enable the webhook or false to disable it",
				Required:            true,
			},
			"type": schema.SingleNestedAttribute{
				MarkdownDescription: "Defines the kind of webhook being configured. The type determines whether the webhook operates as an HTTP request or performs an integration action. See the nested schema below for further details.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"http_request": schema.SingleNestedAttribute{
						MarkdownDescription: "Manual Webhook HTTP Request",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								MarkdownDescription: "The endpoint URL to which the HTTP request is sent. This is the target server or service that the webhook interacts with",
								Required:            true,
							},
							"method": schema.StringAttribute{
								MarkdownDescription: "The HTTP method used for the request, such as GET, POST, PUT, PATCH or DELETE. The method determines the type of operation the webhook performs on the target resource",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.OneOf(allowedHttpMethods...),
								},
							},
							"headers": schema.MapAttribute{
								MarkdownDescription: "Key-value pairs representing HTTP headers to include in the request. These headers can be used to pass metadata or authentication tokens",
								Required:            true,
								ElementType:         types.StringType,
								Sensitive:           true,
							},
						},
					},
					"integration": schema.SingleNestedAttribute{
						MarkdownDescription: "Manual Webhook Integration",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"integration_id": schema.StringAttribute{
								MarkdownDescription: "A unique identifier for the integration associated with the webhook assigned by Apono. This links the webhook to a specific integration within your system",
								Required:            true,
							},
							"action_name": schema.StringAttribute{
								MarkdownDescription: "The name of the action that the webhook performs as part of the integration. Allowed values are: 'does_user_have_permission', 'invoke_azure_function'",
								Required:            true,
							},
						},
					},
				},
			},
			"body_template": schema.StringAttribute{
				MarkdownDescription: "A customizable template for the request body. Use this to format the payload sent by the webhook, allowing context-specific content",
				Optional:            true,
			},
			"response_validators": schema.SetNestedAttribute{
				MarkdownDescription: "A collection of validators to verify the response received from the webhook endpoint. Each validator checks specific conditions to ensure the response meets the expected criteria. See the nested schema below for details",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"json_path": schema.StringAttribute{
							MarkdownDescription: "A JSON expression to extract specific parts of the webhook response. This is used to pinpoint and validate elements within a structured JSON response",
							Required:            true,
						},
						"expected_values": schema.SetAttribute{
							MarkdownDescription: "A list of values the response data must match at the specified json_path to pass validation",
							Required:            true,
							ElementType:         types.StringType,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
			"timeout_in_sec": schema.Int64Attribute{
				MarkdownDescription: "The maximum time, in seconds, that the webhook waits for a response from the endpoint before timing out",
				Optional:            true,
			},
			"authentication_config": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration details for authenticating the webhook requests. See the nested schema below for details",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"oauth": schema.SingleNestedAttribute{
						MarkdownDescription: "Contains OAuth-specific configuration details required for secure communication. See the nested schema below for more information",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"client_id": schema.StringAttribute{
								MarkdownDescription: "The client identifier issued by the OAuth provider. This is used to authenticate the webhook application",
								Required:            true,
							},
							"client_secret": schema.StringAttribute{
								MarkdownDescription: "The secret associated with the client identifier. Keep this value secure, as it is critical for establishing trusted communication",
								Required:            true,
								Sensitive:           true,
							},
							"token_endpoint_url": schema.StringAttribute{
								MarkdownDescription: "The URL where the webhook can request OAuth tokens. This is part of the OAuth workflow to obtain access tokens for secure access",
								Required:            true,
							},
							"scopes": schema.SetAttribute{
								MarkdownDescription: "A list of permissions or access levels the webhook requests from the OAuth provider. Defaults to an empty list if no specific scopes are needed",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
				},
			},
			"custom_validation_error_message": schema.StringAttribute{
				MarkdownDescription: "A custom error message to display when the webhook fails validation. This provides clear feedback to users in case of validation issues",
				Optional:            true,
			},
		},
	}
}

func (w *ManualWebhookResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	w.provider, response.Diagnostics = toProvider(request.ProviderData)
}

func (w ManualWebhookResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *models.ManualWebhookModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	manualWebhook, getWebhookResponse, err := w.provider.terraformClient.WebhooksAPI.TerraformGetWebhook(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		if utils.IsAponoApiNotFoundError(getWebhookResponse) {
			response.State.RemoveResource(ctx)
			tflog.Debug(ctx, "Webhook is deleted, removing from state", map[string]interface{}{
				"id": data.ID.ValueString(),
			})
			return
		}

		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "manual_webhook", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertManualWebhookApiToTerraformModel(ctx, manualWebhook)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully fetched manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (w ManualWebhookResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *models.ManualWebhookModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	newManualWebhookRequest, diagnostics := services.ConvertManualWebhookTerraformModelToUpsertApi(data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	manualWebhook, _, err := w.provider.terraformClient.WebhooksAPI.TerraformCreateWebhook(ctx).
		WebhookManualTriggerUpsertTerraformModel(*newManualWebhookRequest).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "create", "manual webhook", "")
		response.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := services.ConvertManualWebhookApiToTerraformModel(ctx, manualWebhook)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully created manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (w ManualWebhookResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data *models.ManualWebhookModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	updateManualWebhookRequest, diagnostics := services.ConvertManualWebhookTerraformModelToUpsertApi(data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	ManualWebhook, _, err := w.provider.terraformClient.WebhooksAPI.TerraformUpdateWebhook(ctx, data.ID.ValueString()).
		WebhookManualTriggerUpsertTerraformModel(*updateManualWebhookRequest).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "update", "manual webhook", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertManualWebhookApiToTerraformModel(ctx, ManualWebhook)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully updated manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (w ManualWebhookResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *models.ManualWebhookModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting manual webhook", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	messageResponse, deleteResp, err := w.provider.terraformClient.WebhooksAPI.TerraformDeleteWebhook(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		if utils.IsAponoApiNotFoundError(deleteResp) {
			tflog.Debug(ctx, "Manual webhook was already deleted", map[string]interface{}{
				"id":       data.ID.ValueString(),
				"response": messageResponse.GetMessage(),
			})
			return
		}
		diagnostics := utils.GetDiagnosticsForApiError(err, "delete", "manual webhook", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	tflog.Debug(ctx, "Successfully deleted manual webhook", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": messageResponse.GetMessage(),
	})
}

func (w ManualWebhookResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	ManualWebhookId := request.ID
	tflog.Debug(ctx, "Importing manual webhook", map[string]interface{}{
		"id": ManualWebhookId,
	})

	ManualWebhook, _, err := w.provider.terraformClient.WebhooksAPI.TerraformGetWebhook(ctx, ManualWebhookId).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "manual webhook", ManualWebhookId)
		response.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := services.ConvertManualWebhookApiToTerraformModel(ctx, ManualWebhook)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save imported data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully imported manual webhook", map[string]interface{}{
		"id": ManualWebhookId,
	})
}

func (w ManualWebhookResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if w.provider == nil {
		return
	}

	var model models.ManualWebhookModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributePath := path.Root("timeout_in_sec")
	timeoutInSecDefined := !model.TimeoutInSec.IsNull() && !model.TimeoutInSec.IsUnknown()
	if timeoutInSecDefined {
		timeoutInSec := model.TimeoutInSec.ValueInt64()
		if timeoutInSec <= 0 {
			resp.Diagnostics.AddAttributeError(
				attributePath,
				"Invalid timeout_in_sec value",
				"must be a positive number",
			)
		}
	}

	isTypeHttpRequestDefined := model.Type.HttpRequest != nil
	isTypeIntegrationDefined := model.Type.Integration != nil

	if !isTypeHttpRequestDefined && !isTypeIntegrationDefined {
		resp.Diagnostics.AddError(
			"Invalid manual webhook configuration",
			"either type.http_request or type.integration must be specified",
		)
	}

	if isTypeHttpRequestDefined && isTypeIntegrationDefined {
		resp.Diagnostics.AddError(
			"Invalid manual webhook configuration",
			"only one of type.http_request or type.integration must be specified",
		)
	}
}

package provider

import (
	"context"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/schemas"
	"github.com/apono-io/terraform-provider-apono/internal/services"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &accessBundleResource{}
var _ resource.ResourceWithImportState = &accessBundleResource{}

func NewAccessBundleResource() resource.Resource {
	return &accessBundleResource{}
}

// accessBundleResource defines the resource implementation.
type accessBundleResource struct {
	provider *AponoProvider
}

func (a accessBundleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_access_bundle"
}

func (a accessBundleResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Apono Access Bundle",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Access Bundle identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Access Bundle name",
				Required:            true,
			},
			"integration_targets": schemas.GetIntegrationTargetSchema(true),
		},
	}
}

func (a *accessBundleResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	a.provider, response.Diagnostics = toProvider(request.ProviderData)
}

func (a accessBundleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *models.AccessBundleModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	accessBundle, _, err := a.provider.client.AccessBundlesApi.GetAccessBundle(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access_bundle", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessBundleApiToTerraformModel(ctx, a.provider.client, accessBundle)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully fetching access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (a accessBundleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *models.AccessBundleModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	newAccessBundleRequest, diagnostics := services.ConvertAccessBundleTerraformModelToApi(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessBundle, _, err := a.provider.client.AccessBundlesApi.CreateAccessBundle(ctx).UpsertAccessBundleV1(*newAccessBundleRequest).Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "create", "access bundle", "")
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessBundleApiToTerraformModel(ctx, a.provider.client, accessBundle)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully created access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (a accessBundleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data *models.AccessBundleModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	updateAccessBundleRequest, diagnostics := services.ConvertAccessBundleTerraformModelToUpdateApi(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessBundle, _, err := a.provider.client.AccessBundlesApi.UpdateAccessBundle(ctx, data.ID.ValueString()).
		UpdateAccessBundleV1(*updateAccessBundleRequest).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "update", "access bundle", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessBundleApiToTerraformModel(ctx, a.provider.client, accessBundle)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully updated access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (a accessBundleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *models.AccessBundleModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting access bundle", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	messageResponse, _, err := a.provider.client.AccessBundlesApi.DeleteAccessBundle(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "delete", "access bundle", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	tflog.Debug(ctx, "Successfully deleted access bundle", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": messageResponse.GetMessage(),
	})
}

func (a accessBundleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	accessBundleId := request.ID
	tflog.Debug(ctx, "Importing access bundle", map[string]interface{}{
		"id": accessBundleId,
	})

	accessBundle, _, err := a.provider.client.AccessBundlesApi.GetAccessBundle(ctx, accessBundleId).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access bundle", accessBundleId)
		response.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := services.ConvertAccessBundleApiToTerraformModel(ctx, a.provider.client, accessBundle)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save imported data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully imported access bundle", map[string]interface{}{
		"id": accessBundleId,
	})
}

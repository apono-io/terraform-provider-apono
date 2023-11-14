package provider

import (
	"context"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/services"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &accessFlowResource{}
var _ resource.ResourceWithImportState = &accessFlowResource{}

func NewAccessFlowResource() resource.Resource {
	return &accessFlowResource{}
}

// accessFlowResource defines the resource implementation.
type accessFlowResource struct {
	provider *AponoProvider
}

func (a accessFlowResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_access_flow"
}

func (a accessFlowResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	var resourceFilterSchema = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "Filter type, can be 'id', 'name' or 'tag'",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("id", "name", "tag"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Filter name, only used when type is 'tag'",
				Optional:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Filter value",
				Required:            true,
			},
		},
	}
	var identitySchema = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Identity Name (in user type, use email address instead)",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Identity type (user, group, context_attribute)",
				Required:            true,
			},
		},
	}

	response.Schema = schema.Schema{
		MarkdownDescription: "Apono Access Flow",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Access Flow identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Access Flow name",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Is Access Flow active",
				Required:            true,
			},
			"revoke_after_in_sec": schema.NumberAttribute{
				MarkdownDescription: "Number of seconds after which access should be revoked, -1 means never",
				Required:            true,
			},
			"trigger": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of trigger, currently only 'user_request' is supported",
						Required:            true,
					},
					"timeframe": schema.SingleNestedAttribute{
						MarkdownDescription: "Timeframe for trigger to be active",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"start_of_day_time_in_seconds": schema.NumberAttribute{
								MarkdownDescription: "Number of seconds after midnight",
								Required:            true,
							},
							"end_of_day_time_in_seconds": schema.NumberAttribute{
								MarkdownDescription: "Number of seconds after midnight",
								Required:            true,
							},
							"days_in_week": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "Number of seconds after midnight when trigger should be inactive",
								Required:            true,
							},
							"time_zone": schema.StringAttribute{
								MarkdownDescription: "Timezone for timeframe, use IANA timezone name (e.g. Europe/Prague). For all options see [Wiki Page](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)",
								Required:            true,
							},
						},
					},
				},
			},
			"grantees": schema.SetNestedAttribute{
				MarkdownDescription: "Represents which identities should be granted access",
				Required:            true,
				NestedObject:        identitySchema,
			},
			"integration_targets": schema.SetNestedAttribute{
				MarkdownDescription: "Represents number of resources from integration to grant access to. If both include and exclude filters are omitted then all resources will be targeted",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Target Integration name (must match existing integration name)",
							Required:            true,
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Target resource type",
							Required:            true,
						},
						"resource_include_filters": schema.SetNestedAttribute{
							MarkdownDescription: "Include every resource that matches one of this filters",
							Optional:            true,
							NestedObject:        resourceFilterSchema,
						},
						"resource_excludes_filters": schema.SetNestedAttribute{
							MarkdownDescription: "Exclude every resource that matches one of this filters",
							Optional:            true,
							NestedObject:        resourceFilterSchema,
						},
						"permissions": schema.SetAttribute{
							MarkdownDescription: "Permissions to grant",
							Required:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
			"approvers": schema.SetNestedAttribute{
				MarkdownDescription: "Represents which identities should approve this access",
				Optional:            true,
				NestedObject:        identitySchema,
			},
			"settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Access Flow settings",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"require_justification_on_request_again": schema.BoolAttribute{
						MarkdownDescription: "Require justification on request again",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"require_all_approvers": schema.BoolAttribute{
						MarkdownDescription: "Require all approvers",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"approver_cannot_approve_himself": schema.BoolAttribute{
						MarkdownDescription: "Approver cannot approve himself",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	}
}

func (a *accessFlowResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	a.provider, response.Diagnostics = toProvider(request.ProviderData)
}

func (a accessFlowResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *models.AccessFlowModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Fetching access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	accessFlow, _, err := a.provider.client.AccessFlowsApi.GetAccessFlowV1(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access_flow", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertToAccessFlowModel(ctx, a.provider.client, accessFlow)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully fetching access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (a accessFlowResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *models.AccessFlowModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	newAccessFlowRequest, diagnostics := services.ConvertToAccessFlowUpsertApiModel(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessFlow, _, err := a.provider.client.AccessFlowsApi.CreateAccessFlowV1(ctx).UpsertAccessFlowV1(*newAccessFlowRequest).Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "create", "access flow", "")
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertToAccessFlowModel(ctx, a.provider.client, accessFlow)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully created access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

}

func (a accessFlowResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data *models.AccessFlowModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	updateAccessFlowRequest, diagnostics := services.ConvertToAccessFlowUpdateApiModel(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessFlow, _, err := a.provider.client.AccessFlowsApi.UpdateAccessFlowV1(ctx, data.ID.ValueString()).
		UpdateAccessFlowV1(*updateAccessFlowRequest).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "update", "access flow", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertToAccessFlowModel(ctx, a.provider.client, accessFlow)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully updated access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (a accessFlowResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *models.AccessFlowModel

	// Read Terraform prior state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting access flow", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	messageResponse, _, err := a.provider.client.AccessFlowsApi.DeleteAccessFlowV1(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "delete", "access flow", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	tflog.Debug(ctx, "Successfully deleted access flow", map[string]interface{}{
		"id":       data.ID.ValueString(),
		"response": messageResponse.GetMessage(),
	})
}

func (a accessFlowResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	accessFlowId := request.ID
	tflog.Debug(ctx, "Importing access flow", map[string]interface{}{
		"id": accessFlowId,
	})

	accessFlow, _, err := a.provider.client.AccessFlowsApi.GetAccessFlowV1(ctx, accessFlowId).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access flow", accessFlowId)
		response.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := services.ConvertToAccessFlowModel(ctx, a.provider.client, accessFlow)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	// Save imported data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)

	tflog.Debug(ctx, "Successfully imported access flow", map[string]interface{}{
		"id": accessFlowId,
	})
}

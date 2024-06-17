package provider

import (
	"context"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/schemas"
	"github.com/apono-io/terraform-provider-apono/internal/services"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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

var _ resource.ResourceWithValidateConfig = &accessFlowResource{}

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
	var allowedDaysOfTheWeek []string
	for _, day := range aponoapi.AllowedDayOfWeekEnumValues {
		allowedDaysOfTheWeek = append(allowedDaysOfTheWeek, string(day))
	}

	var identitySchema = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the identity. When `type = context_attribute`, this is the shift name or manager attribute name. When `type = group`, this is the group name. When `type = user`, this is the email address. **NOTE: If a non-unique name is used with 'group' type, Apono applies the access flow to all groups matching the name.**",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Identity type. **Possible Values**: `context_attribute`, `group`, or `user`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("user", "group", "context_attribute"),
				},
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
				MarkdownDescription: "Indicates whether Access flow is active or inactive",
				Required:            true,
			},
			"revoke_after_in_sec": schema.NumberAttribute{
				MarkdownDescription: "Number of seconds after which access should be revoked. To never revoke access, set the value to `-1`.",
				Required:            true,
			},
			"trigger": schema.SingleNestedAttribute{
				MarkdownDescription: "Access Flow trigger",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of trigger. 'user_request' or 'automatic' is supported.",
						Required:            true,
					},
					"timeframe": schema.SingleNestedAttribute{
						MarkdownDescription: "Active duration for the trigger.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"start_time": schema.StringAttribute{
								MarkdownDescription: "Beginning of the timeframe in `HH:MM:SS` format.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(utils.TimeRegex, "Time must be in HH:MM:SS format"),
								},
							},
							"end_time": schema.StringAttribute{
								MarkdownDescription: "End of the timeframe in `HH:MM:SS` format.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(utils.TimeRegex, "Time must be in HH:MM:SS format"),
								},
							},
							"days_in_week": schema.SetAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "Names in uppercase of timeframe active days",
								Required:            true,
								Validators: []validator.Set{
									setvalidator.ValueStringsAre(stringvalidator.OneOf(allowedDaysOfTheWeek...)),
								},
							},
							"time_zone": schema.StringAttribute{
								MarkdownDescription: "Timezone name for the timeframe, such as `Europe/Prague`. For all options, see  [Wiki Page](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).",
								Required:            true,
							},
						},
					},
				},
			},
			"grantees": schema.SetNestedAttribute{
				MarkdownDescription: "Represents which identities should be granted access",
				Optional:            true,
				NestedObject:        identitySchema,
				DeprecationMessage:  "Configure grantees_filter_group instead. This attribute will be removed in the next major version of the provider",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"grantees_filter_group": schema.SingleNestedAttribute{
				MarkdownDescription: "Create a conditions group based on different attribute types that represents who can request access.",
				// This field is Optional as long as the old `grantees` field is present in the configuration
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"conditions_logical_operator": schemas.ConditionLogicalOperatorSchema,
					"attribute_filters": schema.SetNestedAttribute{
						MarkdownDescription: "placeholder", // TODO: Add description
						Required:            true,
						NestedObject:        schemas.AttributeFilterSchema,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
			"integration_targets": schemas.GetIntegrationTargetSchema(false),
			"bundle_targets":      schemas.GetBundleTargetSchema(false),
			"approvers": schema.SetNestedAttribute{
				MarkdownDescription: "Represents which identities should approve this access",
				Optional:            true,
				NestedObject:        identitySchema,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
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
						MarkdownDescription: "Require justification on request again.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"require_all_approvers": schema.BoolAttribute{
						MarkdownDescription: "All approvers must approver the request.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"approver_cannot_self_approve": schema.BoolAttribute{
						MarkdownDescription: "Approver cannot self-approve the request.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "List of labels to attach to the access flow",
				Optional:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
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

	accessFlow, _, err := a.provider.terraformClient.AccessFlowsAPI.GetAccessFlowV1(ctx, data.ID.ValueString()).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access_flow", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessFlowApiToTerraformModel(ctx, a.provider.client, accessFlow, data)
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

	newAccessFlowRequest, diagnostics := services.ConvertAccessFlowTerraformModelToApi(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessFlow, _, err := a.provider.terraformClient.AccessFlowsAPI.CreateAccessFlowV1(ctx).UpsertAccessFlowTerraformV1(*newAccessFlowRequest).Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "create", "access flow", "")
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessFlowApiToTerraformModel(ctx, a.provider.client, accessFlow, data)
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

	updateAccessFlowRequest, diagnostics := services.ConvertAccessFlowTerraformModelToApi(ctx, a.provider.client, data)
	if len(diagnostics) > 0 {
		response.Diagnostics.Append(diagnostics...)
		return
	}

	accessFlow, _, err := a.provider.terraformClient.AccessFlowsAPI.UpdateAccessFlowV1(ctx, data.ID.ValueString()).
		UpsertAccessFlowTerraformV1(*updateAccessFlowRequest).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "update", "access flow", data.ID.ValueString())
		response.Diagnostics.Append(diagnostics...)

		return
	}

	model, diagnostics := services.ConvertAccessFlowApiToTerraformModel(ctx, a.provider.client, accessFlow, data)
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

	messageResponse, _, err := a.provider.terraformClient.AccessFlowsAPI.DeleteAccessFlowV1(ctx, data.ID.ValueString()).
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

	accessFlow, _, err := a.provider.terraformClient.AccessFlowsAPI.GetAccessFlowV1(ctx, accessFlowId).
		Execute()
	if err != nil {
		diagnostics := utils.GetDiagnosticsForApiError(err, "get", "access flow", accessFlowId)
		response.Diagnostics.Append(diagnostics...)
		return
	}

	model, diagnostics := services.ConvertAccessFlowApiToTerraformModel(ctx, a.provider.client, accessFlow, nil)
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

func (a *accessFlowResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if a.provider == nil {
		return
	}

	var model models.AccessFlowModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributePath := path.Root("revoke_after_in_sec")

	revokeAfterInSec, _ := model.RevokeAfterInSec.ValueBigFloat().Int64()
	if revokeAfterInSec != -1 && revokeAfterInSec <= 0 {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid revoke_after_in_sec value",
			"must be -1 or positive number",
		)
	}

	if len(model.IntegrationTargets) == 0 && len(model.BundleTargets) == 0 {
		resp.Diagnostics.AddError(
			"Invalid access flow configuration",
			"at least one integration_target or bundle_target must be specified",
		)
	}

	isGranteeFilterGroupDefined := !model.GranteesFilterGroup.IsNull() && !model.GranteesFilterGroup.IsUnknown()
	isGranteesDefined := !model.Grantees.IsNull() && !model.Grantees.IsUnknown()

	if !isGranteeFilterGroupDefined && !isGranteesDefined {
		resp.Diagnostics.AddError(
			"Invalid access flow configuration",
			"either grantees or grantees_filter_group must be specified",
		)
	}

	if isGranteeFilterGroupDefined && isGranteesDefined {
		resp.Diagnostics.AddError(
			"Invalid access flow configuration",
			"only one of grantees or grantees_filter_group must be specified",
		)
	}
}

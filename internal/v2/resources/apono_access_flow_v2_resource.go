package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = &AponoAccessFlowV2Resource{}
	_ resource.ResourceWithImportState = &AponoAccessFlowV2Resource{}

	defaultRequestScopes = setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("self"),
	}))
)

func NewAponoAccessFlowV2Resource() resource.Resource {
	return &AponoAccessFlowV2Resource{}
}

type AponoAccessFlowV2Resource struct {
	client client.Invoker
}

func (r *AponoAccessFlowV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_flow_v2"
}

type IdentityConditionSchemaType string

const (
	IdentityConditionSchemaTypeApprover  IdentityConditionSchemaType = "approver"
	IdentityConditionSchemaTypeRequestor IdentityConditionSchemaType = "requestor"
	IdentityConditionSchemaTypeGrantee   IdentityConditionSchemaType = "grantee"
)

func getIdentityConditionSchema(conditionType IdentityConditionSchemaType) schema.NestedAttributeObject {
	var typeDescription string
	var sourceIntegrationDescription string
	var valuesDescription string
	var matchOperatorDescription string

	switch conditionType {
	case IdentityConditionSchemaTypeApprover:
		typeDescription = "Approver identity type - user, group, Owner, manager, Context Integration, or any other custom value.\nNote: The Owner value must be capitalized (with an uppercase “O”)."
		sourceIntegrationDescription = "Applies when the identity type stems from a Context or IDP integration."
		valuesDescription = "Approver values according to the attribute type and match_operator (e.g., user email, group IDs, etc)."
		matchOperatorDescription = `Comparison operator. Possible values: is, is_not, contains, does_not_contain, starts_with. Defaults to is.
Note: When using is or is_not with any type, you can specify either the source ID or Apono ID to define the requestors.
For the user attribute specifically, you may also use the user’s email.`
	case IdentityConditionSchemaTypeRequestor, IdentityConditionSchemaTypeGrantee:
		typeDescription = "Identity type (e.g., user, group, etc.)"
		sourceIntegrationDescription = "The integration the user/group is from."
		valuesDescription = "List of values according to the attribute type and match_operator (e.g., user emails, group IDs, etc.)."
		matchOperatorDescription = `Comparison operator. Possible values: is, is_not, contains, does_not_contain, starts_with. Defaults to is.
Note: When using is or is_not with any type, you can specify either the source ID or Apono ID to define the requestors.
For the user attribute specifically, you may also use the user’s email.`
	}

	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"source_integration_name": schema.StringAttribute{
				Description: sourceIntegrationDescription,
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: typeDescription,
				Required:    true,
			},
			"match_operator": schema.StringAttribute{
				Description: matchOperatorDescription,
				Optional:    true,
				Default:     stringdefault.StaticString(common.DefaultMatchOperator),
				Computed:    true,
			},
			"values": schema.ListAttribute{
				Description: valuesDescription,
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *AponoAccessFlowV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Access Flow that defines how users or groups can request or automatically be granted access to integrations, bundles, or access scopes under specific conditions and policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the access flow.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name for the access flow, must be unique.",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the access flow is active. Defaults to true.",
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: `The trigger type for the access flow. Possible values: SELF_SERVE, AUTOMATIC.`,
				Required:    true,
			},
			"grant_duration_in_min": schema.Int32Attribute{
				Description: "How long access is granted, in minutes. If not specified, the grant duration defaults to indefinite.",
				Optional:    true,
			},
			"timeframe": schema.SingleNestedAttribute{
				Description: "Restrict when access can be granted. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"start_time": schema.StringAttribute{
						Description: "Start time (e.g., 08:00).",
						Required:    true,
					},
					"end_time": schema.StringAttribute{
						Description: "End time (e.g., 17:00).",
						Required:    true,
					},
					"days_of_week": schema.SetAttribute{
						Description: "Days when access is allowed. (e.g., ['MONDAY', 'TUESDAY']).",
						ElementType: types.StringType,
						Required:    true,
					},
					"time_zone": schema.StringAttribute{
						Description: "Timezone name (e.g., Asia/Jerusalem).",
						Required:    true,
					},
				},
			},
			"approver_policy": schema.SingleNestedAttribute{
				Description: "Approval policy for the access request. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"approval_mode": schema.StringAttribute{
						Description: "Possible values: ANY_OF or ALL_OF. Specifies the logical condition for approvals: ANY_OF: The request is granted if at least one approver from the list approves. ALL_OF: The request is granted only if all approvers in the list approve.",
						Required:    true,
					},
					"approver_groups": schema.SetNestedAttribute{
						Description: "List of approver groups. Cannot be empty.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"logical_operator": schema.StringAttribute{
									Description: `Possible values: AND or OR`,
									Required:    true,
								},
								"approvers": schema.ListNestedAttribute{
									Description:  "List of approvers.",
									Required:     true,
									NestedObject: getIdentityConditionSchema(IdentityConditionSchemaTypeApprover),
								},
							},
						},
					},
				},
			},
			"requestors": schema.SingleNestedAttribute{
				Description: "List of users who can request access, based on identity attributes (e.g., users, groups, or shifts) and the conditions under which they can request access.\nIn self-serve access flows, requestors specify who is allowed to submit an access request.\nIn automatic access flows, requestors specify who will automatically receive access when conditions are met (equivalent to \"grantees\" in the UI).",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"logical_operator": schema.StringAttribute{
						Description: `Specifies the logical operator to be used between the requestors in the list. Possible values: "AND" or "OR".`,
						Required:    true,
					},
					"conditions": schema.ListNestedAttribute{
						Description:  "List of conditions. Cannot be empty.",
						Required:     true,
						NestedObject: getIdentityConditionSchema(IdentityConditionSchemaTypeRequestor),
					},
				},
			},
			"request_for": schema.SingleNestedAttribute{
				Description: "Defines who the access request can be made for. This enables support to request on behalf of other users, groups, or identities. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"request_scopes": schema.SetAttribute{
						MarkdownDescription: `Specifies who the request can be made for. Supported values:
1. "self" - The user making the request (default behavior).
2. "others" - Specific individuals specified manually.
3. "direct_reports" - Allows the requestor, identified as a manager in the organization's identity provider (IdP), to request access for individuals formally assigned as direct reports in the IdP (based on IdP integration).

Defaults to ["self"].`,
						Optional:    true,
						Computed:    true,
						Default:     defaultRequestScopes,
						ElementType: types.StringType,
					},
					"grantees": schema.SingleNestedAttribute{
						Description: "Applicable only when \"others\" is included in request_scope. Defines the set of users or attributes who can be selected as recipients of the access.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"logical_operator": schema.StringAttribute{
								Description: `Specifies the logical operator to be used between the grantees in the list. Possible values: "AND" or "OR".`,
								Required:    true,
							},
							"conditions": schema.ListNestedAttribute{
								Description:  "List of conditions. Cannot be empty.",
								Required:     true,
								NestedObject: getIdentityConditionSchema(IdentityConditionSchemaTypeGrantee),
							},
						},
					},
				},
			},
			"access_targets": schema.ListNestedAttribute{
				Description: "Define the targets accessible when requesting access via this access flow.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"integration": schemas.GetIntegrationTargetSchema(schemas.ResourceMode),
						"bundle": schema.SingleNestedAttribute{
							Description: "Bundle target.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Name of the bundle.",
									Required:    true,
								},
							},
						},
						"access_scope": schemas.GetAccessScopeTargetSchema(schemas.ResourceMode),
					},
				},
			},
			"settings": schema.SingleNestedAttribute{
				Description: "Settings for the access flow.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"justification_required": schema.BoolAttribute{
						Description: "Require justification from requestor. Defaults to true. Must be set to false for automatic access flows. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
						Optional:    true,
						Default:     booldefault.StaticBool(true),
						Computed:    true,
					},
					"require_approver_reason": schema.BoolAttribute{
						Description: "Require reason from approver. Defaults to false. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"requester_cannot_approve_self": schema.BoolAttribute{
						Description: "Requester cannot approve their own requests. Defaults to false. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"require_mfa": schema.BoolAttribute{
						Description: "Require MFA at approval time. Defaults to false. Only applicable in self-serve access flows (trigger = \"SELF_SERVE\").",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"labels": schema.SetAttribute{
						Description: "Custom labels for organizational use.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (r *AponoAccessFlowV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoAccessFlowV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.AccessFlowV2Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertRequest, err := models.AccessFlowModelToUpsertRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access flow",
			fmt.Sprintf("Unable to create access flow, got error: %s", err),
		)
		return
	}

	accessFlow, err := r.client.CreateAccessFlowV2(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access flow",
			fmt.Sprintf("Unable to create access flow, got error: %s", err),
		)
		return
	}

	accessFlowModel, err := models.AccessFlowResponseToModel(ctx, *accessFlow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access flow",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, accessFlowModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoAccessFlowV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.AccessFlowV2Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessFlow, err := r.client.GetAccessFlowV2(ctx, client.GetAccessFlowV2Params{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading access flow", fmt.Sprintf("Unable to read access flow with ID %s, got error: %s", state.ID.ValueString(), err))
		return
	}

	accessFlowModel, err := models.AccessFlowResponseToModel(ctx, *accessFlow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading access flow",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, accessFlowModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoAccessFlowV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.AccessFlowV2Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertRequest, err := models.AccessFlowModelToUpsertRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access flow",
			fmt.Sprintf("Unable to update access flow, got error: %s", err),
		)
		return
	}

	accessFlow, err := r.client.UpdateAccessFlowV2(ctx, upsertRequest, client.UpdateAccessFlowV2Params{
		ID: plan.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access flow",
			fmt.Sprintf("Unable to update access flow, got error: %s", err),
		)
		return
	}

	accessFlowModel, err := models.AccessFlowResponseToModel(ctx, *accessFlow)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access flow",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, accessFlowModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoAccessFlowV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.AccessFlowV2Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAccessFlowV2(ctx, client.DeleteAccessFlowV2Params{
		ID: state.ID.ValueString(),
	}); err != nil {
		if client.IsNotFoundError(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting access flow",
			fmt.Sprintf("Unable to delete access flow with ID %s, got error: %s", state.ID.ValueString(), err),
		)
	}
}

func (r *AponoAccessFlowV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

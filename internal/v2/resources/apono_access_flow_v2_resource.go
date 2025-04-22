package resources

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &AponoAccessFlowV2Resource{}
	_ resource.ResourceWithImportState = &AponoAccessFlowV2Resource{}
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

func (r *AponoAccessFlowV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Access Flow V2.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the access flow.",
				Required:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the access flow is active. Defaults to true.",
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
			},
			"trigger": schema.StringAttribute{
				Description: `The trigger type for the access flow. Can be "SELF_SERVE" or "AUTOMATIC".`,
				Required:    true,
			},
			"grant_duration_in_min": schema.Int32Attribute{
				Description: "The grant duration in minutes. Null means indefinite. Cannot be negative.",
				Optional:    true,
			},
			"timeframe": schema.SingleNestedAttribute{
				Description: "Optional timeframe restriction for the access flow.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"start_time": schema.StringAttribute{
						Description: "Start time in HH:mm format (e.g., '10:00').",
						Required:    true,
					},
					"end_time": schema.StringAttribute{
						Description: "End time in HH:mm format (e.g., '23:59').",
						Required:    true,
					},
					"days_of_week": schema.SetAttribute{
						Description: "Days of the week when access is allowed (e.g., ['MONDAY', 'TUESDAY']).",
						ElementType: types.StringType,
						Required:    true,
					},
					"time_zone": schema.StringAttribute{
						Description: "Time zone for the timeframe (e.g., 'Asia/Jerusalem').",
						Required:    true,
					},
				},
			},
			"approver_policy": schema.SingleNestedAttribute{
				Description: "Approval policy configuration. If null, requests are auto-approved.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"approval_mode": schema.StringAttribute{
						Description: `The approval mode. Can be "ANY_OF" or "ALL_OF".`,
						Required:    true,
					},
					"approver_groups": schema.SetNestedAttribute{
						Description: "List of approver groups. Cannot be empty.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"logical_operator": schema.StringAttribute{
									Description: `The logical operator for the approvers. Can be "OR" or "AND".`,
									Required:    true,
								},
								"approvers": schema.SetNestedAttribute{
									Description: "List of approvers.",
									Required:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"source_integration_name": schema.StringAttribute{
												Description: "The name of the source integration.",
												Optional:    true,
											},
											"type": schema.StringAttribute{
												Description: "The type of approver.",
												Required:    true,
											},
											"match_operator": schema.StringAttribute{
												Description: `The match operator. Defaults to "is".`,
												Optional:    true,
												Default:     stringdefault.StaticString("is"),
												Computed:    true,
											},
											"values": schema.SetAttribute{
												Description: "The values to match against.",
												Required:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"grantees": schema.SingleNestedAttribute{
				Description: "The users or groups that can be granted access through this flow.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"logical_operator": schema.StringAttribute{
						Description: `The logical operator for the conditions. Can be "OR" or "AND".`,
						Required:    true,
					},
					"conditions": schema.SetNestedAttribute{
						Description: "List of conditions. Cannot be empty.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source_integration_name": schema.StringAttribute{
									Description: "The name of the source integration.",
									Optional:    true,
								},
								"type": schema.StringAttribute{
									Description: "The type of grantee.",
									Required:    true,
								},
								"match_operator": schema.StringAttribute{
									Description: `The match operator. Possible values: "starts_with", "contains", "is_not", "does_not_contain", "is". Defaults to "is".`,
									Optional:    true,
									Default:     stringdefault.StaticString("is"),
									Computed:    true,
								},
								"values": schema.SetAttribute{
									Description: "The values to match against.",
									Required:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"access_targets": schema.SetNestedAttribute{
				Description: "List of access targets for this access flow",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"integration": getBundleIntegrationSchema(),
						"bundle": schema.SingleNestedAttribute{
							Description: "Bundle target configuration",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the bundle",
									Required:    true,
								},
							},
						},
						"access_scope": getBundleAccessScopeSchema(),
					},
				},
			},
			"settings": schema.SingleNestedAttribute{
				Description: "Access flow settings configuration",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"justification_required": schema.BoolAttribute{
						Description: "Whether justification is required when requesting access. Defaults to true.",
						Optional:    true,
						Default:     booldefault.StaticBool(true),
						Computed:    true,
					},
					"require_approver_reason": schema.BoolAttribute{
						Description: "Whether approvers must provide a reason when approving/denying requests. Defaults to false.",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"requester_cannot_approve_self": schema.BoolAttribute{
						Description: "Whether requesters are prevented from approving their own requests. Defaults to false.",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"require_mfa": schema.BoolAttribute{
						Description: "Whether MFA is required for this access flow. Defaults to false.",
						Optional:    true,
						Default:     booldefault.StaticBool(false),
						Computed:    true,
					},
					"labels": schema.SetAttribute{
						Description: "List of labels associated with this access flow",
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
	// TODO: Implement create
}

func (r *AponoAccessFlowV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO: Implement read
}

func (r *AponoAccessFlowV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO: Implement update
}

func (r *AponoAccessFlowV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO: Implement delete
}

func (r *AponoAccessFlowV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO: Implement import
}

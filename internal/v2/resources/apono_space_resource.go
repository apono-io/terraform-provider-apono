package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithConfigure   = &AponoSpaceResource{}
	_ resource.ResourceWithImportState = &AponoSpaceResource{}
)

func NewAponoSpaceResource() resource.Resource {
	return &AponoSpaceResource{}
}

type AponoSpaceResource struct {
	client client.Invoker
}

func (r *AponoSpaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (r *AponoSpaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Space, an organizational scoping unit that partitions the access management domain so teams can operate independently with their own Access Flows, Bundles, and Access Scopes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the space.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Display name for the space; must be unique within Apono.",
				Required:    true,
			},
			"space_scope_references": schema.SetAttribute{
				Description: "Space scope names that define the space's inventory scope.",
				ElementType: types.StringType,
				Required:    true,
			},
			"members": schema.SetNestedAttribute{
				Description: "Members of the space.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identity_reference": schema.StringAttribute{
							Description: "Reference to the identity. For users: user ID or email. For groups: group ID or name.",
							Required:    true,
						},
						"identity_type": schema.StringAttribute{
							Description: "Type of identity: user or group. Determines how identity_reference is resolved.",
							Required:    true,
						},
						"space_roles": schema.SetAttribute{
							Description: "Roles within the space: SpaceOwner (full control) or SpaceManager (manage resources).",
							ElementType: types.StringType,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *AponoSpaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SpaceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scopeRefs []string
	diags = plan.SpaceScopeReferences.ElementsAs(ctx, &scopeRefs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateSpaceV1{
		Name:                 plan.Name.ValueString(),
		SpaceScopeReferences: scopeRefs,
	}

	if !plan.Members.IsNull() {
		var planMembers []models.SpaceMemberModel
		diags = plan.Members.ElementsAs(ctx, &planMembers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		apiMembers := spaceMemberModelsToAPI(ctx, planMembers)
		createReq.Members = client.NewOptNilUpsertSpaceMemberV1Array(apiMembers)
	}

	space, err := r.client.CreateSpaceV1(ctx, &createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating space", fmt.Sprintf("Could not create space: %v", err))
		return
	}

	result, err := models.SpaceToModel(ctx, space)
	if err != nil {
		resp.Diagnostics.AddError("Error creating space", fmt.Sprintf("Could not convert space response: %v", err))
		return
	}

	result.Members = plan.Members

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r *AponoSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.SpaceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.client.GetSpaceV1(ctx, client.GetSpaceV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading space", fmt.Sprintf("Could not read space ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result, err := models.SpaceToModel(ctx, space)
	if err != nil {
		resp.Diagnostics.AddError("Error reading space", fmt.Sprintf("Could not convert space response: %v", err))
		return
	}

	membersResp, err := services.ListSpaceMembers(ctx, r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading space members", fmt.Sprintf("Could not read members for space ID %s: %v", state.ID.ValueString(), err))
		return
	}

	if len(membersResp) > 0 || !state.Members.IsNull() {
		memberModels, err := models.SpaceMembersToModels(ctx, membersResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading space members", fmt.Sprintf("Could not convert members response: %v", err))
			return
		}

		membersSet, setDiags := types.SetValueFrom(ctx, spaceMemberObjectType(), memberModels)
		resp.Diagnostics.Append(setDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		result.Members = membersSet
	} else {
		result.Members = types.SetNull(spaceMemberObjectType())
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

func (r *AponoSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.SpaceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan models.SpaceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scopeRefs []string
	diags = plan.SpaceScopeReferences.ElementsAs(ctx, &scopeRefs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpdateSpaceV1{
		Name:                 plan.Name.ValueString(),
		SpaceScopeReferences: scopeRefs,
	}

	params := client.UpdateSpaceV1Params{ID: state.ID.ValueString()}
	space, err := r.client.UpdateSpaceV1(ctx, &updateReq, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating space", fmt.Sprintf("Could not update space ID %s: %v", state.ID.ValueString(), err))
		return
	}

	updatedModel, err := models.SpaceToModel(ctx, space)
	if err != nil {
		resp.Diagnostics.AddError("Error updating space", fmt.Sprintf("Could not convert space response: %v", err))
		return
	}

	state.Name = updatedModel.Name
	state.SpaceScopeReferences = updatedModel.SpaceScopeReferences

	if !plan.Members.Equal(state.Members) {
		var planMembers []models.SpaceMemberModel
		if !plan.Members.IsNull() {
			diags = plan.Members.ElementsAs(ctx, &planMembers, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		apiMembers := spaceMemberModelsToAPI(ctx, planMembers)
		updateMembersReq := client.UpdateSpaceMembersV1{
			Members: apiMembers,
		}

		_, err := r.client.ReplaceSpaceMembersV1(ctx, &updateMembersReq, client.ReplaceSpaceMembersV1Params{ID: state.ID.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError("Error updating space members", fmt.Sprintf("Could not update members for space ID %s: %v", state.ID.ValueString(), err))
			return
		}

		state.Members = plan.Members
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *AponoSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SpaceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSpaceV1(ctx, client.DeleteSpaceV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting space", fmt.Sprintf("Could not delete space ID %s: %v", state.ID.ValueString(), err))
		return
	}
}

func (r *AponoSpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func spaceMemberModelsToAPI(_ context.Context, members []models.SpaceMemberModel) []client.UpsertSpaceMemberV1 {
	result := make([]client.UpsertSpaceMemberV1, len(members))
	for i, m := range members {
		var roles []string
		m.SpaceRoles.ElementsAs(context.Background(), &roles, false)

		result[i] = client.UpsertSpaceMemberV1{
			IdentityReference: m.IdentityReference.ValueString(),
			IdentityType:      m.IdentityType.ValueString(),
			SpaceRoles:        roles,
		}
	}

	return result
}

func spaceMemberObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"identity_reference": types.StringType,
			"identity_type":      types.StringType,
			"space_roles":        types.SetType{ElemType: types.StringType},
		},
	}
}

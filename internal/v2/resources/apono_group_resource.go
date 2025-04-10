package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/services"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &AponoGroupResource{}
	_ resource.ResourceWithImportState = &AponoGroupResource{}
)

func NewAponoGroupResource() resource.Resource {
	return &AponoGroupResource{}
}

// AponoGroupResource manages Apono Group resources.
type AponoGroupResource struct {
	client client.Invoker
}

func (r *AponoGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *AponoGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the group.",
				Required:    true,
			},
			"members": schema.SetAttribute{
				Description: "List of member email addresses in the group.",
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *AponoGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan services.GroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var emails []string
	diags = plan.Members.ElementsAs(ctx, &emails, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateGroupV1{
		Name:          plan.Name.ValueString(),
		MembersEmails: emails,
	}

	tflog.Debug(ctx, "Creating group", map[string]any{
		"name":          plan.Name.ValueString(),
		"member_emails": emails,
	})

	group, err := r.client.CreateGroupV1(ctx, &createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", fmt.Sprintf("Could not create group: %v", err))
		return
	}

	result := services.GroupToModel(group)
	result.Members = plan.Members

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Created group successfully", map[string]any{"id": result.ID.ValueString()})
}

func (r *AponoGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state services.GroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.GetGroupV1(ctx, client.GetGroupV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading group", fmt.Sprintf("Could not read group ID %s: %v", state.ID.ValueString(), err))
		return
	}

	membersResp, err := services.ListGroupMembers(ctx, r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading group members", fmt.Sprintf("Could not read members for group ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result := services.GroupToModel(group)

	memberEmails := []string{}
	for _, member := range membersResp {
		memberEmails = append(memberEmails, member.Email)
	}

	membersSet, diags := types.SetValueFrom(ctx, types.StringType, memberEmails)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result.Members = membersSet

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state services.GroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan services.GroupModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) {
		updateNameReq := client.UpdateGroupV1{
			Name: plan.Name.ValueString(),
		}

		tflog.Debug(ctx, "Updating group name", map[string]any{
			"id":      state.ID.ValueString(),
			"oldName": state.Name.ValueString(),
			"newName": plan.Name.ValueString(),
		})

		params := client.UpdateGroupV1Params{ID: state.ID.ValueString()}
		group, err := r.client.UpdateGroupV1(ctx, &updateNameReq, params)
		if err != nil {
			resp.Diagnostics.AddError("Error updating group name", fmt.Sprintf("Could not update group ID %s: %v", state.ID.ValueString(), err))
			return
		}

		state = services.GroupToModel(group)
	}

	if !plan.Members.Equal(state.Members) {
		var planMembers []string
		diags = plan.Members.ElementsAs(ctx, &planMembers, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		updateMembersReq := client.UpdateGroupMembersV1{
			MembersEmails: planMembers,
		}

		tflog.Debug(ctx, "Updating group members", map[string]any{
			"id":      state.ID.ValueString(),
			"members": planMembers,
		})

		err := r.client.UpdateGroupMembersV1(ctx, &updateMembersReq, client.UpdateGroupMembersV1Params{ID: state.ID.ValueString()})
		if err != nil {
			resp.Diagnostics.AddError("Error updating group members", fmt.Sprintf("Could not update members for group ID %s: %v", state.ID.ValueString(), err))
			return
		}

		state.Members = plan.Members
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updated group successfully", map[string]any{"id": state.ID.ValueString()})
}

func (r *AponoGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state services.GroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroupV1(ctx, client.DeleteGroupV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting group", fmt.Sprintf("Could not delete group ID %s: %v", state.ID.ValueString(), err))
		return
	}

	tflog.Info(ctx, "Deleted group successfully", map[string]any{"id": state.ID.ValueString()})
}

func (r *AponoGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

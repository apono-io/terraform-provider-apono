package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.ResourceWithConfigure   = &AponoSpaceScopeResource{}
	_ resource.ResourceWithImportState = &AponoSpaceScopeResource{}
)

func NewAponoSpaceScopeResource() resource.Resource {
	return &AponoSpaceScopeResource{}
}

type AponoSpaceScopeResource struct {
	client client.Invoker
}

func (r *AponoSpaceScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_scope"
}

func (r *AponoSpaceScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Space Scope, an AQL expression that defines which resources are visible within a space.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier of the space scope.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "A descriptive name for the space scope. Must be unique within Apono.",
				Required:    true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "AQL ([Apono Query Language](https://docs.apono.io/docs/inventory/apono-query-language)) expression that filters which resources belong to this scope.",
				Required:            true,
			},
		},
	}
}

func (r *AponoSpaceScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoSpaceScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SpaceScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.UpsertSpaceScopeV1{
		Name:  plan.Name.ValueString(),
		Query: plan.Query.ValueString(),
	}

	spaceScope, err := r.client.CreateSpaceScopeV1(ctx, &createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating space scope", fmt.Sprintf("Could not create space scope: %v", err))
		return
	}

	result := models.SpaceScopeToModel(spaceScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *AponoSpaceScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.SpaceScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceScope, err := r.client.GetSpaceScopeV1(ctx, client.GetSpaceScopeV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading space scope", fmt.Sprintf("Could not read space scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result := models.SpaceScopeToModel(spaceScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoSpaceScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state models.SpaceScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan models.SpaceScopeModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpsertSpaceScopeV1{
		Name:  plan.Name.ValueString(),
		Query: plan.Query.ValueString(),
	}

	params := client.UpdateSpaceScopeV1Params{ID: state.ID.ValueString()}
	spaceScope, err := r.client.UpdateSpaceScopeV1(ctx, &updateReq, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating space scope", fmt.Sprintf("Could not update space scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result := models.SpaceScopeToModel(spaceScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *AponoSpaceScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SpaceScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSpaceScopeV1(ctx, client.DeleteSpaceScopeV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting space scope", fmt.Sprintf("Could not delete space scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

}

func (r *AponoSpaceScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

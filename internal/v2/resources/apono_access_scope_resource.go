package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &AponoAccessScopeResource{}
	_ resource.ResourceWithImportState = &AponoAccessScopeResource{}
)

func NewAponoAccessScopeResource() resource.Resource {
	return &AponoAccessScopeResource{}
}

type AponoAccessScopeResource struct {
	client client.Invoker
}

func (r *AponoAccessScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_scope"
}

func (r *AponoAccessScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Access Scope, a logical grouping of cloud resources defined by a flexible query.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for this Apono Access Scope. You can reference it in other Terraform resources or use it to import an existing access scope into your Terraform state.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "A descriptive name for the access scope. It must be unique within Apono.",
				Required:    true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "A query string written in [Apono Query Language](https://docs.apono.io/docs/inventory/apono-query-language).",
				Required:            true,
			},
		},
	}
}

func (r *AponoAccessScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoAccessScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan services.AccessScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.UpsertAccessScopeV1{
		Name:  plan.Name.ValueString(),
		Query: plan.Query.ValueString(),
	}

	tflog.Debug(ctx, "Creating access scope", map[string]any{
		"name":  plan.Name.ValueString(),
		"query": plan.Query.ValueString(),
	})

	accessScope, err := r.client.CreateAccessScopesV1(ctx, &createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating access scope", fmt.Sprintf("Could not create access scope: %v", err))
		return
	}

	result := services.AccessScopeToModel(accessScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Created access scope successfully", map[string]any{"id": result.ID.ValueString()})
}

func (r *AponoAccessScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state services.AccessScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	accessScope, err := r.client.GetAccessScopesV1(ctx, client.GetAccessScopesV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading access scope", fmt.Sprintf("Could not read access scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result := services.AccessScopeToModel(accessScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoAccessScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state services.AccessScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan services.AccessScopeModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := client.UpsertAccessScopeV1{
		Name:  plan.Name.ValueString(),
		Query: plan.Query.ValueString(),
	}

	tflog.Debug(ctx, "Updating access scope", map[string]any{
		"id":    state.ID.ValueString(),
		"name":  plan.Name.ValueString(),
		"query": plan.Query.ValueString(),
	})

	params := client.UpdateAccessScopesV1Params{ID: state.ID.ValueString()}
	accessScope, err := r.client.UpdateAccessScopesV1(ctx, &updateReq, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating access scope", fmt.Sprintf("Could not update access scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

	result := services.AccessScopeToModel(accessScope)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updated access scope successfully", map[string]any{"id": result.ID.ValueString()})
}

func (r *AponoAccessScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state services.AccessScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAccessScopesV1(ctx, client.DeleteAccessScopesV1Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting access scope", fmt.Sprintf("Could not delete access scope ID %s: %v", state.ID.ValueString(), err))
		return
	}

	tflog.Info(ctx, "Deleted access scope successfully", map[string]any{"id": state.ID.ValueString()})
}

func (r *AponoAccessScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

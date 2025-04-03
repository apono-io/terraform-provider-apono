package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AponoAccessScopeResource{}
	_ resource.ResourceWithImportState = &AponoAccessScopeResource{}
)

// NewAccessScopeResource is a helper function to simplify the provider implementation.
func NewAponoAccessScopeResource() resource.Resource {
	return &AponoAccessScopeResource{}
}

// AponoAccessScopeResource is the resource implementation.
type AponoAccessScopeResource struct {
	client client.Invoker
}

// accessScopeResourceModel maps the resource schema data to a Terraform-friendly format.
type accessScopeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Query        types.String `tfsdk:"query"`
	CreationDate types.String `tfsdk:"creation_date"`
	UpdateDate   types.String `tfsdk:"update_date"`
}

// Metadata returns the resource type name.
func (r *AponoAccessScopeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_scope"
}

// Schema defines the schema for the resource.
func (r *AponoAccessScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Access Scope.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the access scope.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the access scope.",
				Required:    true,
			},
			"query": schema.StringAttribute{
				Description: "The query expression for the access scope.",
				Required:    true,
			},
			"creation_date": schema.StringAttribute{
				Description: "The date when the access scope was created.",
				Computed:    true,
			},
			"update_date": schema.StringAttribute{
				Description: "The date when the access scope was last updated.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *AponoAccessScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureClientInvoker(ctx, req, resp, &r.client)
}

// Create creates a new access scope and sets the initial Terraform state.
func (r *AponoAccessScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan accessScopeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new access scope
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
		resp.Diagnostics.AddError(
			"Error creating access scope",
			fmt.Sprintf("Could not create access scope: %v", err),
		)
		return
	}

	// Map API response to model
	result := accessScopeApiToModel(accessScope)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Created access scope successfully", map[string]any{
		"id": result.ID.ValueString(),
	})
}

// Read refreshes the Terraform state with the latest data.
func (r *AponoAccessScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state accessScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get access scope from API
	accessScope, err := r.client.GetAccessScopesV1(ctx, client.GetAccessScopesV1Params{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		// Check if the resource no longer exists
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading access scope",
			fmt.Sprintf("Could not read access scope ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	// Map API response to model
	result := accessScopeApiToModel(accessScope)

	// Set refreshed state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *AponoAccessScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state accessScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get plan
	var plan accessScopeResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update access scope
	updateReq := client.UpsertAccessScopeV1{
		Name:  plan.Name.ValueString(),
		Query: plan.Query.ValueString(),
	}

	tflog.Debug(ctx, "Updating access scope", map[string]any{
		"id":    state.ID.ValueString(),
		"name":  plan.Name.ValueString(),
		"query": plan.Query.ValueString(),
	})

	params := client.UpdateAccessScopesV1Params{
		ID: state.ID.ValueString(),
	}

	accessScope, err := r.client.UpdateAccessScopesV1(ctx, &updateReq, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access scope",
			fmt.Sprintf("Could not update access scope ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	// Map API response to model
	result := accessScopeApiToModel(accessScope)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updated access scope successfully", map[string]any{
		"id": result.ID.ValueString(),
	})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *AponoAccessScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state accessScopeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete access scope by calling API
	err := r.client.DeleteAccessScopesV1(ctx, client.DeleteAccessScopesV1Params{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		// If the error is that the resource doesn't exist, it's already gone, so no error
		if client.IsNotFoundError(err) {
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting access scope",
			fmt.Sprintf("Could not delete access scope ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	tflog.Info(ctx, "Deleted access scope successfully", map[string]any{
		"id": state.ID.ValueString(),
	})
}

// ImportState imports an existing resource into Terraform.
func (r *AponoAccessScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp) // TODO import by name
}

// Helper function to convert API response to model.
func accessScopeApiToModel(accessScope *client.AccessScopeV1) *accessScopeResourceModel {
	model := &accessScopeResourceModel{
		ID:    types.StringValue(accessScope.ID),
		Name:  types.StringValue(accessScope.Name),
		Query: types.StringValue(accessScope.Query),
	}

	creationDate := time.Time(accessScope.CreationDate)
	updateDate := time.Time(accessScope.UpdateDate)

	// Format creation date if present
	if !creationDate.IsZero() {
		model.CreationDate = types.StringValue(creationDate.Format(time.RFC3339))
	} else {
		model.CreationDate = types.StringNull()
	}

	// Format update date if present
	if !updateDate.IsZero() {
		model.UpdateDate = types.StringValue(updateDate.Format(time.RFC3339))
	} else {
		model.UpdateDate = types.StringNull()
	}

	return model
}

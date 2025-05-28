package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.ResourceWithConfigure   = &AponoBundleV2Resource{}
	_ resource.ResourceWithImportState = &AponoBundleV2Resource{}
)

func NewAponoBundleV2Resource() resource.Resource {
	return &AponoBundleV2Resource{}
}

type AponoBundleV2Resource struct {
	client client.Invoker
}

func (r *AponoBundleV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle_v2"
}

func (r *AponoBundleV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Bundle, which defines a collection of access targets - either access scopes or specific resources within integrations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the bundle. ",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name for the access flow, must be unique.",
				Required:    true,
			},
			"access_targets": schema.ListNestedAttribute{
				Description: "A list of access targets included in the bundle.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"integration":  schemas.GetIntegrationTargetSchema(schemas.ResourceMode),
						"access_scope": schemas.GetAccessScopeTargetSchema(schemas.ResourceMode),
					},
				},
			},
		},
	}
}

func (r *AponoBundleV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoBundleV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.BundleV2Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertRequest, err := models.BundleModelToUpsertRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bundle",
			fmt.Sprintf("Unable to create bundle, got error: %s", err),
		)
		return
	}

	bundle, err := r.client.CreateBundleV2(ctx, upsertRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bundle",
			fmt.Sprintf("Unable to create bundle, got error: %s", err),
		)
		return
	}

	bundleModel, err := models.BundleResponseToModel(ctx, *bundle)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bundle",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, bundleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoBundleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.BundleV2Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bundle, err := r.client.GetBundleV2(ctx, client.GetBundleV2Params{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading bundle", fmt.Sprintf("Unable to read bundle with ID %s, got error: %s", state.ID.ValueString(), err))
		return
	}

	bundleModel, err := models.BundleResponseToModel(ctx, *bundle)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading bundle",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, bundleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoBundleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.BundleV2Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	upsertRequest, err := models.BundleModelToUpsertRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bundle",
			fmt.Sprintf("Unable to update bundle, got error: %s", err),
		)
		return
	}

	bundle, err := r.client.UpdateBundleV2(ctx, upsertRequest, client.UpdateBundleV2Params{
		ID: plan.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bundle",
			fmt.Sprintf("Unable to update bundle, got error: %s", err),
		)
		return
	}

	bundleModel, err := models.BundleResponseToModel(ctx, *bundle)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating bundle",
			fmt.Sprintf("Unable to convert API response to model: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, bundleModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoBundleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.BundleV2Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteBundleV2(ctx, client.DeleteBundleV2Params{
		ID: state.ID.ValueString(),
	}); err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting bundle",
			fmt.Sprintf("Unable to delete bundle with ID %s, got error: %s", state.ID.ValueString(), err),
		)
	}
}

func (r *AponoBundleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

package resources

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.Resource                     = &AponoBundleV2Resource{}
	_ resource.ResourceWithImportState      = &AponoBundleV2Resource{}
	_ resource.ResourceWithConfigValidators = &AponoBundleV2Resource{}
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

func (r *AponoBundleV2Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRelative().AtParent().AtName("integration"),
			path.MatchRelative().AtParent().AtName("access_scope"),
		),
		resourcevalidator.Conflicting(
			path.MatchRelative().AtParent().AtName("integration"),
			path.MatchRelative().AtParent().AtName("access_scope"),
		),
	}
}

func (r *AponoBundleV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Apono Bundle V2.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the bundle.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the bundle.",
				Required:    true,
			},
			"access_targets": schema.SetNestedAttribute{
				Description: "List of access targets for this bundle",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"integration":  schemas.GetIntegrationTargetSchema(),
						"access_scope": schemas.GetAccessScopeTargetSchema(),
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
	// Placeholder implementation
	resp.Diagnostics.AddError(
		"Not Implemented",
		"The create method for apono_bundle_v2 resource has not been implemented yet",
	)
}

func (r *AponoBundleV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Placeholder implementation
	resp.Diagnostics.AddError(
		"Not Implemented",
		"The read method for apono_bundle_v2 resource has not been implemented yet",
	)
}

func (r *AponoBundleV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Placeholder implementation
	resp.Diagnostics.AddError(
		"Not Implemented",
		"The update method for apono_bundle_v2 resource has not been implemented yet",
	)
}

func (r *AponoBundleV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Placeholder implementation
	resp.Diagnostics.AddError(
		"Not Implemented",
		"The delete method for apono_bundle_v2 resource has not been implemented yet",
	)
}

func (r *AponoBundleV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

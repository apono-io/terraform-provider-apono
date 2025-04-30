package resources

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/schemas"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = &AponoResourceIntegrationResource{}
	_ resource.ResourceWithImportState      = &AponoResourceIntegrationResource{}
	_ resource.ResourceWithConfigValidators = &AponoResourceIntegrationResource{}
)

func NewAponoResourceIntegrationResource() resource.Resource {
	return &AponoResourceIntegrationResource{}
}

type AponoResourceIntegrationResource struct {
	client client.Invoker
}

func (r *AponoResourceIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_integration"
}

func (r *AponoResourceIntegrationResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRelative().AtName("secret_store_config").AtName("aws"),
			path.MatchRelative().AtName("secret_store_config").AtName("gcp"),
			path.MatchRelative().AtName("secret_store_config").AtName("azure"),
			path.MatchRelative().AtName("secret_store_config").AtName("hashicorp_vault"),
		),
	}
}

func (r *AponoResourceIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Resource Integration, allowing Apono to connect and manage external cloud resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the integration.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Human-readable name for the integration, must be unique within Apono.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: `Type of the integration (e.g., "aws-account", "postgresql").`,
				Required:    true,
			},
			"connector_id": schema.StringAttribute{
				Description: "ID of the Apono Connector used for the integration.",
				Required:    true,
			},
			"connected_resource_types": schema.ListAttribute{
				Description: "List of resource types for the integration to discover.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"integration_config": schema.MapAttribute{
				MarkdownDescription: "Integration-specific configuration that accepts key-value pairs. Refer to the [Integration Configuration documentation](https://docs.apono.io/metadata-for-integration-config) for specific configuration values.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"secret_store_config": schemas.GetSecretStoreConfigSchema(schemas.ResourceMode),
			"custom_access_details": schema.StringAttribute{
				Description: "Custom access instructions for end users, displayed in the access details modal.",
				Optional:    true,
			},
			"owner": schema.SingleNestedAttribute{
				Description: "Apono can use the integration owner for access requests approval if no owner is found. Enter one or more users, groups, shifts or attributes. This field is mandatory when using Resource Owners and serves as a fallback approver if no resource owner is found.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"source_integration_name": schema.StringAttribute{
						Description: "Name of the integration from which the type originates from (e.g. “Google Oauth”).",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "Type of the owner attribute.",
						Required:    true,
					},
					"values": schema.ListAttribute{
						Description: "List of values for the ownership assignment.",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			"owners_mapping": schema.SingleNestedAttribute{
				Description: "Apono will sync each resource's owner from the source integration. Use this for Resource Owner access requests approval.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"source_integration_name": schema.StringAttribute{
						Description: "Name of the source integration.",
						Optional:    true,
					},
					"key_name": schema.StringAttribute{
						Description: "Attribute key to map owner.",
						Required:    true,
					},
					"attribute_type": schema.StringAttribute{
						Description: "Type of the attribute.",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *AponoResourceIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	common.ConfigureResourceClientInvoker(ctx, req, resp, &r.client)
}

func (r *AponoResourceIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ResourceIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq, err := models.ResourceIntegrationModelToCreateRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource integration request",
			fmt.Sprintf("Could not create API request: %s", err),
		)
		return
	}

	integration, err := r.client.CreateIntegrationV4(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource integration",
			fmt.Sprintf("Could not create resource integration: %s", err),
		)
		return
	}

	result, err := models.ResourceIntegrationToModel(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting resource integration",
			fmt.Sprintf("Could not convert resource integration: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoResourceIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ResourceIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.GetIntegrationsByIdV4(ctx, client.GetIntegrationsByIdV4Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading resource integration",
			fmt.Sprintf("Could not read resource integration ID %s: %v", state.ID.ValueString(), err),
		)
		return
	}

	if integration.Category != common.ResourceCategory {
		resp.Diagnostics.AddError(
			"Invalid resource integration type",
			fmt.Sprintf("Expected resource integration, got %s", integration.Category),
		)
		return
	}

	result, err := models.ResourceIntegrationToModel(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting resource integration",
			fmt.Sprintf("Could not convert resource integration: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoResourceIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan models.ResourceIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq, err := models.ResourceIntegrationModelToUpdateRequest(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource integration update request",
			fmt.Sprintf("Could not create API request: %s", err),
		)
		return
	}

	integration, err := r.client.UpdateIntegrationV4(ctx, updateReq, client.UpdateIntegrationV4Params{ID: state.ID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating resource integration",
			fmt.Sprintf("Could not update resource integration: %s", err),
		)
		return
	}

	result, err := models.ResourceIntegrationToModel(ctx, integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting resource integration",
			fmt.Sprintf("Could not convert resource integration: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AponoResourceIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ResourceIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIntegrationV4(ctx, client.DeleteIntegrationV4Params{ID: state.ID.ValueString()})
	if err != nil {
		if client.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting resource integration",
			fmt.Sprintf("Could not delete resource integration ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}
}

func (r *AponoResourceIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

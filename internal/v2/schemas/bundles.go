package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetIntegrationTargetSchema(mode SchemaMode) schema.SingleNestedAttribute {
	isComputed := mode == DataSourceMode
	fieldsRequired := mode == ResourceMode
	fieldsComputed := mode == DataSourceMode

	return schema.SingleNestedAttribute{
		Description: "Defines an integration and resources to which access will be granted.",
		Optional:    !isComputed,
		Computed:    isComputed,
		Attributes: map[string]schema.Attribute{
			"integration_name": schema.StringAttribute{
				Description: "The name of the integration",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
			"resource_type": schema.StringAttribute{
				Description: "The type of resource within the integration for which access is being granted (e.g., aws-account-s3-bucket).",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
			"permissions": schema.SetAttribute{
				Description: `List of permissions (e.g., "Attach", "ReadOnlyAccess").`,
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
				ElementType: types.StringType,
			},
			"resources_scopes": schema.ListNestedAttribute{
				Description: "A list of filters defining which resources are included or excluded. If null, the scope will apply to any resource in the integration target",
				Optional:    !isComputed,
				Computed:    isComputed,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"scope_mode": schema.StringAttribute{
							Description: "Possible values: `include_resources` or `exclude_resources`. `include_resources`: Grants access to the specific resources listed under the `values` field. `exclude_resources`: Grants access to all resources within the integration except those specified in the `values` field.",
							Required:    fieldsRequired,
							Computed:    fieldsComputed,
						},
						"type": schema.StringAttribute{
							Description: "NAME - specify resources by their name, APONO_ID - specify resources by their ID, or TAG - specify resources by tag.",
							Required:    fieldsRequired,
							Computed:    fieldsComputed,
						},
						"key": schema.StringAttribute{
							Description: "Tag key. Only required if type = TAG",
							Optional:    !isComputed,
							Computed:    isComputed,
						},
						"values": schema.ListAttribute{
							Description: "Resource values to match (IDs, names, or tag values).",
							Required:    fieldsRequired,
							Computed:    fieldsComputed,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func GetAccessScopeTargetSchema(mode SchemaMode) schema.SingleNestedAttribute {
	isComputed := mode == DataSourceMode
	fieldsRequired := mode == ResourceMode
	fieldsComputed := mode == DataSourceMode

	return schema.SingleNestedAttribute{
		Description: "Access scope target.",
		Optional:    !isComputed,
		Computed:    isComputed,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the access scope.",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
		},
	}
}

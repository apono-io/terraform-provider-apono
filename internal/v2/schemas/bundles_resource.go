package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetIntegrationTargetResourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Integration target.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"integration_name": schema.StringAttribute{
				Description: "The name of the integration",
				Required:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "The type of resource",
				Required:    true,
			},
			"permissions": schema.SetAttribute{
				Description: "List of permissions",
				Required:    true,
				ElementType: types.StringType,
			},
			"resources_scopes": schema.SetNestedAttribute{
				Description: "If null, the scope will apply to any resource in the integration target.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"scope_mode": schema.StringAttribute{
							Description: "Possible values: `include_resources` or `exclude_resources`. `include_resources`: Grants access to the specific resources listed under the `values` field. `exclude_resources`: Grants access to all resources within the integration except those specified in the `values` field.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "NAME - specify resources by their name, APONO_ID - specify resources by their ID, or TAG - specify resources by tag.",
							Required:    true,
						},
						"key": schema.StringAttribute{
							Description: "Tag key. Only required if type = TAG",
							Optional:    true,
						},
						"values": schema.SetAttribute{
							Description: "Resource values to match (IDs, names, or tag values).",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func GetAccessScopeTargetResourceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Access scope.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the access scope.",
				Required:    true,
			},
		},
	}
}

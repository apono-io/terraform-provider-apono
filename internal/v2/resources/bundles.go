package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getBundleIntegrationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Integration configuration",
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
			"resources_scope": schema.SetNestedAttribute{
				Description: "Resource scope configuration",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"scope_mode": schema.StringAttribute{
							Description: "Scope mode - include_resources or exclude_resources",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type - NAME, APONO_ID, or TAG",
							Required:    true,
						},
						"key": schema.StringAttribute{
							Description: "Key - only required for TAG type",
							Optional:    true,
						},
						"values": schema.SetAttribute{
							Description: "List of values - Apono IDs, names, or tag values",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func getBundleAccessScopeSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Access scope configuration",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the access scope",
				Required:    true,
			},
		},
	}
}

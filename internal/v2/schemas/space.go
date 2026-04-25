package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSpaceReferenceResourceAttribute returns the schema attribute for space_reference on resources.
// It is optional and forces replacement when changed.
func GetSpaceReferenceResourceAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Description: "Name of the space to create this resource in. If omitted, the resource is created without a space. Changing this value forces the resource to be replaced.",
		Optional:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

// GetSpaceComputedAttribute returns the computed space nested attribute.
// It is usable in both resource and data source schemas.
func GetSpaceComputedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: "Space this resource belongs to. Null if no space is assigned.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"space_id": schema.StringAttribute{
				Description: "Unique identifier of the space.",
				Computed:    true,
			},
			"space_name": schema.StringAttribute{
				Description: "Unique name of the space.",
				Computed:    true,
			},
		},
	}
}

// GetSpaceReferencesFilterAttribute returns the schema attribute for filtering by space references
// on data sources.
func GetSpaceReferencesFilterAttribute() schema.ListAttribute {
	return schema.ListAttribute{
		Description: "Filter results by space names or IDs. If omitted, returns objects from all spaces.",
		Optional:    true,
		ElementType: types.StringType,
	}
}

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetSpaceReferenceResourceAttribute returns the schema attribute for space_reference on resources.
// It is optional and forces replacement when changed.
//
// Note: space_reference accepts the space NAME only (not an ID), per the Apono Terraform spec.
// On read-back we repopulate it from response.space.space_name so state stays consistent.
// The underlying API technically accepts either an ID or a name on POST, but we intentionally
// expose only names to avoid state drift (ID written → name read back → perpetual replace).
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
func GetSpaceComputedAttribute(mode SchemaMode) schema.SingleNestedAttribute {
	description := "Space details this resource belongs to. Null if the resource has no space assigned."
	if mode == DataSourceMode {
		description = "Space details this item belongs to. Null for items without a space."
	}
	return schema.SingleNestedAttribute{
		Description: description,
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
		Description: `List of space IDs or names to filter results by. If omitted, returns objects from all spaces. To return only objects that belong to no space, include the special value "null" in the list.`,
		Optional:    true,
		ElementType: types.StringType,
	}
}

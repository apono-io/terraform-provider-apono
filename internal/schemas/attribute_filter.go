package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AttributeFilterSchema = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"operator": schema.StringAttribute{
				MarkdownDescription: "placeholder", // TODO: Add description
				Optional:            true,
				Computed:            true,
			},
			"attribute_type": schema.StringAttribute{
				MarkdownDescription: "placeholder", // TODO: Add description
				Required:            true,
			},
			"attribute_names": schema.SetAttribute{
				MarkdownDescription: "placeholder", // TODO: Add description
				Optional:            true,
				ElementType:         types.StringType,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "placeholder", // TODO: Add description
				Optional:            true,
			},
		},
	}

	ConditionLogicalOperatorSchema = schema.StringAttribute{
		MarkdownDescription: "Logical operator to apply to the conditions. **Possible Values**: `AND`, `OR` (Default `OR`)",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("OR"),
	}
)

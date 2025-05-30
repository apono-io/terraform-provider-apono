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
				MarkdownDescription: "Pick the operator that will be applied to the attribute names' values. Defaults to `is`. Supported operators: `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`",
				Optional:            true,
				Computed:            true,
			},
			"attribute_type": schema.StringAttribute{
				MarkdownDescription: "Pick the user context type, for example `user`, `group`, `okta_city`, `pagerduty_shift`, `webhook`, etc.",
				Required:            true,
			},
			"attribute_names": schema.SetAttribute{
				MarkdownDescription: "Insert the specific values you'd like to include or exclude from the Access Flow, for example the user email, group name, webhook name, etc.",
				// Value is Optional because some attribute types may not require it
				Optional:    true,
				ElementType: types.StringType,
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "Use the integration ID this attribute originates from. This can be any user context integration, for example PagerDuty, Okta, etc.",
				Optional:            true,
			},
		},
	}

	ConditionLogicalOperatorSchema = schema.StringAttribute{
		MarkdownDescription: "Logical operator to apply to the conditions. **Possible Values**: `AND`, `OR` (Default `OR`)",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("AND"),
	}
)

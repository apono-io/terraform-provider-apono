package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetIntegrationTargetSchema(required bool) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Represents the number of resources from the integration to which access is granted. If both include and exclude filters are omitted, all resources will be targeted.",
		Required:            required,
		Optional:            !required,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Target integration name. **IMPORTANT: This value must match the existing integration name.**",
					Required:            true,
				},
				"resource_type": schema.StringAttribute{
					MarkdownDescription: "Type of target resource. For possible values, query the [list integrations](https://docs.apono.io/reference/listintegrationsv2) route.",
					Required:            true,
				},
				"resource_include_filters": schema.SetNestedAttribute{
					MarkdownDescription: "Include every resource that matches one of the defined filters.",
					Optional:            true,
					NestedObject:        resourceFilterSchema,
				},
				"resource_exclude_filters": schema.SetNestedAttribute{
					MarkdownDescription: "Exclude every resource that matches one of the defined filters.",
					Optional:            true,
					NestedObject:        resourceFilterSchema,
				},
				"permissions": schema.SetAttribute{
					MarkdownDescription: "Permissions to grant",
					Required:            true,
					ElementType:         types.StringType,
				},
			},
		},
	}
}

func GetBundleTargetSchema(required bool) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Represents the number of resources from access bundle to which access is granted.",
		Required:            required,
		Optional:            !required,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Target bundle name. **IMPORTANT: This value must match the existing bundle name.**",
					Required:            true,
				},
			},
		},
	}
}

var resourceFilterSchema = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "Type of filter, **Possible Values**: 'id', 'name' or 'tag'.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.Any(
					stringvalidator.OneOf("id", "name"),
					stringvalidator.All(
						stringvalidator.OneOf("tag"),
						stringvalidator.AlsoRequires(path.Expressions{path.MatchRelative().AtParent().AtName("key")}...),
					),
				),
			},
		},
		"key": schema.StringAttribute{
			MarkdownDescription: "Key of the filter, **required** only when `type = tag`.",
			Optional:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "Value of the filter",
			Required:            true,
		},
	},
}

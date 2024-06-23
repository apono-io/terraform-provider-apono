package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DataSourceResourceOwnerMapping = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"key_name": schema.StringAttribute{
				MarkdownDescription: "Insert the tag name (key) that represents owners in the cloud environment.",
				Required:            true,
			},
			"attribute": schema.StringAttribute{
				MarkdownDescription: "Insert the attribute type that the tag values will map into. For example: pagerduty_shift, okta_city, group, etc.",
				Required:            true,
			},
			"attribute_integration_id": schema.StringAttribute{
				MarkdownDescription: "Provide the User Context integration ID the attribute originates from, for example Okta, Pagerduty, etc. You can find the ID in the Apono API Reference.",
				Optional:            true,
			},
		},
	}

	DataSourceIntegrationOwner = schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "Provide the User Context integration ID the attribute originates from, for example Okta, Pagerduty, etc. You can find the ID in the Apono API Reference.",
				Optional:            true,
			},
			"attribute": schema.StringAttribute{
				MarkdownDescription: "Insert the attribute type that the tag values will map into. For example: pagerduty_shift, okta_city, group, etc.",
				Required:            true,
			},
			"value": schema.ListAttribute{
				MarkdownDescription: "Provide the attribute value that will serve as the Integration Owner. For example, the user email, group name, etc.",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
)

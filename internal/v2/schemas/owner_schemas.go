package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetOwnerSchema(mode SchemaMode) schema.SingleNestedAttribute {
	isComputed := mode == DataSourceMode
	fieldsRequired := mode == ResourceMode
	fieldsComputed := mode == DataSourceMode

	description := "Apono can use the integration owner for access requests approval if no owner is found. Enter one or more users, groups, shifts or attributes. This field is mandatory when using Resource Owners and serves as a fallback approver if no resource owner is found."
	if mode == DataSourceMode {
		description = "Integration owner information used for access requests approval."
	}

	return schema.SingleNestedAttribute{
		Description: description,
		Optional:    !isComputed,
		Computed:    isComputed,
		Attributes: map[string]schema.Attribute{
			"source_integration_name": schema.StringAttribute{
				Description: "Name of the integration from which the type originates from (e.g. \"Google Oauth\").",
				Optional:    !fieldsComputed,
				Computed:    fieldsComputed,
			},
			"type": schema.StringAttribute{
				Description: "Type of the owner attribute.",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
			"values": schema.ListAttribute{
				Description: "List of values for the ownership assignment.",
				ElementType: types.StringType,
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
		},
	}
}

func GetOwnersMappingSchema(mode SchemaMode) schema.SingleNestedAttribute {
	isComputed := mode == DataSourceMode
	fieldsRequired := mode == ResourceMode
	fieldsComputed := mode == DataSourceMode

	description := "Apono will sync each resource's owner from the source integration. Use this for Resource Owner access requests approval."
	if mode == DataSourceMode {
		description = "Resource owner mapping configuration from the source integration."
	}

	return schema.SingleNestedAttribute{
		Description: description,
		Optional:    !isComputed,
		Computed:    isComputed,
		Attributes: map[string]schema.Attribute{
			"source_integration_name": schema.StringAttribute{
				Description: "Name of the source integration.",
				Optional:    !fieldsComputed,
				Computed:    fieldsComputed,
			},
			"key_name": schema.StringAttribute{
				Description: "Attribute key to map owner.",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Type of the attribute.",
				Required:    fieldsRequired,
				Computed:    fieldsComputed,
			},
		},
	}
}

package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func GetSecretStoreConfigSchema(mode SchemaMode) schema.SingleNestedAttribute {
	isComputed := mode == DataSourceMode

	fieldsRequired := mode == ResourceMode
	fieldsComputed := mode == DataSourceMode

	return schema.SingleNestedAttribute{
		Description: "Configuration for secret store.",
		Optional:    true,
		Computed:    isComputed,
		Attributes: map[string]schema.Attribute{
			"aws": schema.SingleNestedAttribute{
				Description: "AWS secret store configuration.",
				Optional:    !isComputed,
				Computed:    isComputed,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Description: "The AWS region.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
					"secret_id": schema.StringAttribute{
						Description: "The AWS secret ID.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
				},
			},
			"gcp": schema.SingleNestedAttribute{
				Description: "GCP secret store configuration.",
				Optional:    !isComputed,
				Computed:    isComputed,
				Attributes: map[string]schema.Attribute{
					"project": schema.StringAttribute{
						Description: "The GCP project.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
					"secret_id": schema.StringAttribute{
						Description: "The GCP secret ID.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
				},
			},
			"azure": schema.SingleNestedAttribute{
				Description: "Azure secret store configuration.",
				Optional:    !isComputed,
				Computed:    isComputed,
				Attributes: map[string]schema.Attribute{
					"vault_url": schema.StringAttribute{
						Description: "The Azure Vault URL.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
					"name": schema.StringAttribute{
						Description: "The Azure secret name.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
				},
			},
			"hashicorp_vault": schema.SingleNestedAttribute{
				Description: "HashiCorp Vault secret store configuration.",
				Optional:    !isComputed,
				Computed:    isComputed,
				Attributes: map[string]schema.Attribute{
					"secret_engine": schema.StringAttribute{
						Description: "The HashiCorp Vault secret engine.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
					"path": schema.StringAttribute{
						Description: "The HashiCorp Vault path.",
						Required:    fieldsRequired,
						Computed:    fieldsComputed,
					},
				},
			},
		},
	}
}

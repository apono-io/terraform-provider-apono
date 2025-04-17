package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceIntegrationModel struct {
	ID                              types.String         `tfsdk:"id"`
	Name                            types.String         `tfsdk:"name"`
	Type                            types.String         `tfsdk:"type"`
	ConnectorID                     types.String         `tfsdk:"connector_id"`
	ConnectedResourceTypes          types.List           `tfsdk:"connected_resource_types"`
	IntegrationConfig               types.Map            `tfsdk:"integration_config"`
	SecretStoreConfig               *SecretStoreConfig   `tfsdk:"secret_store_config"`
	CustomAccessDetails             types.String         `tfsdk:"custom_access_details"`
	UserCleanupPeriodInDays         types.Int64          `tfsdk:"user_cleanup_period_in_days"`
	CredentialsRotationPeriodInDays types.Int64          `tfsdk:"credentials_rotation_period_in_days"`
	Owner                           *OwnerConfig         `tfsdk:"owner"`
	OwnersMapping                   *OwnersMappingConfig `tfsdk:"owners_mapping"`
}

type SecretStoreConfig struct {
	AWS            *AWSSecretConfig      `tfsdk:"aws"`
	GCP            *GCPSecretConfig      `tfsdk:"gcp"`
	Azure          *AzureSecretConfig    `tfsdk:"azure"`
	HashicorpVault *HashicorpVaultConfig `tfsdk:"hashicorp_vault"`
}

type AWSSecretConfig struct {
	Region   types.String `tfsdk:"region"`
	SecretID types.String `tfsdk:"secret_id"`
}

type GCPSecretConfig struct {
	Project  types.String `tfsdk:"project"`
	SecretID types.String `tfsdk:"secret_id"`
}

type AzureSecretConfig struct {
	VaultURL types.String `tfsdk:"vault_url"`
	Name     types.String `tfsdk:"name"`
}

type HashicorpVaultConfig struct {
	SecretEngine types.String `tfsdk:"secret_engine"`
	Path         types.String `tfsdk:"path"`
}

type OwnerConfig struct {
	SourceIntegrationName types.String `tfsdk:"source_integration_name"`
	Type                  types.String `tfsdk:"type"`
	Values                types.List   `tfsdk:"values"`
}

type OwnersMappingConfig struct {
	SourceIntegrationName types.String `tfsdk:"source_integration_name"`
	KeyName               types.String `tfsdk:"key_name"`
	AttributeType         types.String `tfsdk:"attribute_type"`
}

func CreateIntegrationRequest(ctx context.Context, model ResourceIntegrationModel) (*client.CreateIntegrationV4, error) {
	req := &client.CreateIntegrationV4{
		Name: model.Name.ValueString(),
		Type: model.Type.ValueString(),
	}

	req.ConnectorID.SetTo(model.ConnectorID.ValueString())

	if !model.ConnectedResourceTypes.IsNull() {
		var connectedResourceTypes []string
		diags := model.ConnectedResourceTypes.ElementsAs(ctx, &connectedResourceTypes, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse connected resource types: %v", diags)
		}
		req.ConnectedResourceTypes.SetTo(connectedResourceTypes)
	}

	if !model.IntegrationConfig.IsNull() {
		integrationConfig := make(map[string]jx.Raw)
		for k, v := range model.IntegrationConfig.Elements() {
			strVal, ok := v.(types.String)
			if !ok {
				return nil, fmt.Errorf("failed to assert type for integration config value")
			}
			integrationConfig[k] = jx.Raw(strVal.ValueString())
		}
		req.IntegrationConfig = integrationConfig
	}

	if model.SecretStoreConfig != nil {
		secretConfig := client.CreateIntegrationV4SecretStoreConfig{}

		if model.SecretStoreConfig.AWS != nil {
			awsConfig := model.SecretStoreConfig.AWS
			secretConfig.AWS = client.NewOptNilCreateIntegrationV4SecretStoreConfigAWS(client.CreateIntegrationV4SecretStoreConfigAWS{
				Region:   awsConfig.Region.ValueString(),
				SecretID: awsConfig.SecretID.ValueString(),
			})
		} else if model.SecretStoreConfig.GCP != nil {
			gcpConfig := model.SecretStoreConfig.GCP
			secretConfig.Gcp = client.NewOptNilCreateIntegrationV4SecretStoreConfigGcp(client.CreateIntegrationV4SecretStoreConfigGcp{
				Project:  gcpConfig.Project.ValueString(),
				SecretID: gcpConfig.SecretID.ValueString(),
			})
		} else if model.SecretStoreConfig.Azure != nil {
			azureConfig := model.SecretStoreConfig.Azure
			secretConfig.Azure = client.NewOptNilCreateIntegrationV4SecretStoreConfigAzure(client.CreateIntegrationV4SecretStoreConfigAzure{
				VaultURL: azureConfig.VaultURL.ValueString(),
				Name:     azureConfig.Name.ValueString(),
			})
		} else if model.SecretStoreConfig.HashicorpVault != nil {
			vaultConfig := model.SecretStoreConfig.HashicorpVault
			secretConfig.HashicorpVault = client.NewOptNilCreateIntegrationV4SecretStoreConfigHashicorpVault(client.CreateIntegrationV4SecretStoreConfigHashicorpVault{
				SecretEngine: vaultConfig.SecretEngine.ValueString(),
				Path:         vaultConfig.Path.ValueString(),
			})
		}

		req.SecretStoreConfig.SetTo(secretConfig)
	}

	if !model.CustomAccessDetails.IsNull() {
		req.CustomAccessDetails.SetTo(model.CustomAccessDetails.ValueString())
	}

	if !model.UserCleanupPeriodInDays.IsNull() {
		req.UserCleanupPeriodInDays.SetTo(model.UserCleanupPeriodInDays.ValueInt64())
	}

	if !model.CredentialsRotationPeriodInDays.IsNull() {
		req.CredentialsRotationPeriodInDays.SetTo(model.CredentialsRotationPeriodInDays.ValueInt64())
	}

	if model.Owner != nil {
		var values []string
		diags := model.Owner.Values.ElementsAs(ctx, &values, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse owner values: %v", diags)
		}

		owner := client.CreateIntegrationV4Owner{
			AttributeType:  model.Owner.Type.ValueString(),
			AttributeValue: values,
		}

		if !model.Owner.SourceIntegrationName.IsNull() {
			owner.SourceIntegrationReference.SetTo(model.Owner.SourceIntegrationName.ValueString())
		}

		owner.SourceIntegrationReference.SetTo(model.Owner.SourceIntegrationName.ValueString())

		req.Owner.SetTo(owner)
	}

	if model.OwnersMapping != nil {
		ownersMapping := client.CreateIntegrationV4OwnersMapping{
			KeyName:       model.OwnersMapping.KeyName.ValueString(),
			AttributeType: model.OwnersMapping.AttributeType.ValueString(),
		}

		if !model.OwnersMapping.SourceIntegrationName.IsNull() {
			ownersMapping.SourceIntegrationReference.SetTo(model.OwnersMapping.SourceIntegrationName.ValueString())
		}

		req.OwnersMapping.SetTo(ownersMapping)
	}

	return req, nil
}

func UpdateIntegrationRequest(ctx context.Context, model ResourceIntegrationModel) (*client.UpdateIntegrationV4, error) {
	req := &client.UpdateIntegrationV4{
		Name: model.Name.ValueString(),
	}

	if !model.IntegrationConfig.IsNull() {
		integrationConfig := make(map[string]jx.Raw)
		for k, v := range model.IntegrationConfig.Elements() {
			strVal, ok := v.(types.String)
			if !ok {
				return nil, fmt.Errorf("failed to assert type for integration config value")
			}
			integrationConfig[k] = jx.Raw(strVal.ValueString())
		}
		req.IntegrationConfig = integrationConfig
	}

	if !model.ConnectedResourceTypes.IsNull() {
		var connectedResourceTypes []string
		diags := model.ConnectedResourceTypes.ElementsAs(ctx, &connectedResourceTypes, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse connected resource types: %v", diags)
		}
		req.ConnectedResourceTypes.SetTo(connectedResourceTypes)
	}

	if model.SecretStoreConfig != nil {
		secretConfig := client.UpdateIntegrationV4SecretStoreConfig{}

		if model.SecretStoreConfig.AWS != nil {
			awsConfig := model.SecretStoreConfig.AWS
			secretConfig.AWS = client.NewOptNilUpdateIntegrationV4SecretStoreConfigAWS(client.UpdateIntegrationV4SecretStoreConfigAWS{
				Region:   awsConfig.Region.ValueString(),
				SecretID: awsConfig.SecretID.ValueString(),
			})
		} else if model.SecretStoreConfig.GCP != nil {
			gcpConfig := model.SecretStoreConfig.GCP
			secretConfig.Gcp = client.NewOptNilUpdateIntegrationV4SecretStoreConfigGcp(client.UpdateIntegrationV4SecretStoreConfigGcp{
				Project:  gcpConfig.Project.ValueString(),
				SecretID: gcpConfig.SecretID.ValueString(),
			})
		} else if model.SecretStoreConfig.Azure != nil {
			azureConfig := model.SecretStoreConfig.Azure
			secretConfig.Azure = client.NewOptNilUpdateIntegrationV4SecretStoreConfigAzure(client.UpdateIntegrationV4SecretStoreConfigAzure{
				VaultURL: azureConfig.VaultURL.ValueString(),
				Name:     azureConfig.Name.ValueString(),
			})
		} else if model.SecretStoreConfig.HashicorpVault != nil {
			vaultConfig := model.SecretStoreConfig.HashicorpVault
			secretConfig.HashicorpVault = client.NewOptNilUpdateIntegrationV4SecretStoreConfigHashicorpVault(client.UpdateIntegrationV4SecretStoreConfigHashicorpVault{
				SecretEngine: vaultConfig.SecretEngine.ValueString(),
				Path:         vaultConfig.Path.ValueString(),
			})
		}

		req.SecretStoreConfig.SetTo(secretConfig)
	}

	if !model.CustomAccessDetails.IsNull() {
		req.CustomAccessDetails.SetTo(model.CustomAccessDetails.ValueString())
	}

	if !model.UserCleanupPeriodInDays.IsNull() {
		req.UserCleanupPeriodInDays.SetTo(model.UserCleanupPeriodInDays.ValueInt64())
	}

	if !model.CredentialsRotationPeriodInDays.IsNull() {
		req.CredentialsRotationPeriodInDays.SetTo(model.CredentialsRotationPeriodInDays.ValueInt64())
	}

	if model.Owner != nil {
		var values []string
		diags := model.Owner.Values.ElementsAs(ctx, &values, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse owner values: %v", diags)
		}

		owner := client.UpdateIntegrationV4Owner{
			AttributeType:  model.Owner.Type.ValueString(),
			AttributeValue: values,
		}

		if !model.Owner.SourceIntegrationName.IsNull() {
			owner.SourceIntegrationReference.SetTo(model.Owner.SourceIntegrationName.ValueString())
		}

		req.Owner.SetTo(owner)
	}

	if model.OwnersMapping != nil {
		ownersMapping := client.UpdateIntegrationV4OwnersMapping{
			KeyName:       model.OwnersMapping.KeyName.ValueString(),
			AttributeType: model.OwnersMapping.AttributeType.ValueString(),
		}

		if !model.OwnersMapping.SourceIntegrationName.IsNull() {
			ownersMapping.SourceIntegrationReference.SetTo(model.OwnersMapping.SourceIntegrationName.ValueString())
		}

		req.OwnersMapping.SetTo(ownersMapping)
	}

	return req, nil
}

func ResourceIntegrationToModel(ctx context.Context, integration *client.IntegrationV4) (*ResourceIntegrationModel, error) {
	model := &ResourceIntegrationModel{
		ID:   types.StringValue(integration.ID),
		Name: types.StringValue(integration.Name),
		Type: types.StringValue(integration.Type),
	}

	model.ConnectorID = types.StringValue(integration.ConnectorID.Value)

	if integration.ConnectedResourceTypes.IsSet() {
		connectedResourceTypes := integration.ConnectedResourceTypes.Value

		stringSlice := make([]string, len(connectedResourceTypes))
		copy(stringSlice, connectedResourceTypes)

		resourceTypes, diags := types.ListValueFrom(ctx, types.StringType, stringSlice)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse connected resource types: %v", diags)
		}
		model.ConnectedResourceTypes = resourceTypes
	} else {
		model.ConnectedResourceTypes = types.ListNull(types.StringType)
	}

	if integration.IntegrationConfig != nil {
		configMap := make(map[string]attr.Value)
		for k, v := range integration.IntegrationConfig {
			configMap[k] = types.StringValue(v.String())
		}
		integrationConfig, diags := types.MapValueFrom(ctx, types.StringType, configMap)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse integration config: %v", diags)
		}
		model.IntegrationConfig = integrationConfig
	}

	if integration.SecretStoreConfig.IsSet() {
		secretConfig := &SecretStoreConfig{}
		apiSecretConfig := integration.SecretStoreConfig.Value

		if apiSecretConfig.AWS.IsSet() {
			awsConfig := apiSecretConfig.AWS.Value
			secretConfig.AWS = &AWSSecretConfig{
				Region:   types.StringValue(awsConfig.Region),
				SecretID: types.StringValue(awsConfig.SecretID),
			}
		} else if apiSecretConfig.Gcp.IsSet() {
			gcpConfig := apiSecretConfig.Gcp.Value
			secretConfig.GCP = &GCPSecretConfig{
				Project:  types.StringValue(gcpConfig.Project),
				SecretID: types.StringValue(gcpConfig.SecretID),
			}
		} else if apiSecretConfig.Azure.IsSet() {
			azureConfig := apiSecretConfig.Azure.Value
			secretConfig.Azure = &AzureSecretConfig{
				VaultURL: types.StringValue(azureConfig.VaultURL),
				Name:     types.StringValue(azureConfig.Name),
			}
		} else if apiSecretConfig.HashicorpVault.IsSet() {
			vaultConfig := apiSecretConfig.HashicorpVault.Value
			secretConfig.HashicorpVault = &HashicorpVaultConfig{
				SecretEngine: types.StringValue(vaultConfig.SecretEngine),
				Path:         types.StringValue(vaultConfig.Path),
			}
		}

		model.SecretStoreConfig = secretConfig
	}

	if integration.CustomAccessDetails.IsSet() {
		model.CustomAccessDetails = types.StringValue(integration.CustomAccessDetails.Value)
	}

	if integration.UserCleanupPeriodInDays.IsSet() {
		model.UserCleanupPeriodInDays = types.Int64Value(integration.UserCleanupPeriodInDays.Value)
	}

	if integration.CredentialsRotationPeriodInDays.IsSet() {
		model.CredentialsRotationPeriodInDays = types.Int64Value(integration.CredentialsRotationPeriodInDays.Value)
	}

	if integration.Owner.IsSet() {
		ownerData := integration.Owner.Value
		values, diags := types.ListValueFrom(ctx, types.StringType, ownerData.AttributeValue)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse owner values: %v", diags)
		}

		ownerConfig := &OwnerConfig{
			Type:   types.StringValue(ownerData.AttributeType),
			Values: values,
		}

		if ownerData.SourceIntegrationName.IsSet() {
			ownerConfig.SourceIntegrationName = types.StringValue(ownerData.SourceIntegrationName.Value)
		}

		model.Owner = ownerConfig
	}

	if integration.OwnersMapping.IsSet() {
		mappingData := integration.OwnersMapping.Value

		ownersMappingConfig := &OwnersMappingConfig{
			KeyName:       types.StringValue(mappingData.KeyName),
			AttributeType: types.StringValue(mappingData.AttributeType),
		}

		if mappingData.SourceIntegrationName.IsSet() {
			ownersMappingConfig.SourceIntegrationName = types.StringValue(mappingData.SourceIntegrationName.Value)
		}

		model.OwnersMapping = ownersMappingConfig
	}

	return model, nil
}

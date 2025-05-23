package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SecretStoreConfig struct {
	AWS            *AWSSecretConfig        `tfsdk:"aws"`
	GCP            *GCPSecretConfig        `tfsdk:"gcp"`
	Azure          *AzureSecretConfig      `tfsdk:"azure"`
	HashicorpVault *HashicorpVaultConfig   `tfsdk:"hashicorp_vault"`
	Kubernetes     *KubernetesSecretConfig `tfsdk:"kubernetes"`
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

type KubernetesSecretConfig struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
}

func upsertSecretStoreConfig(config *SecretStoreConfig) client.UpsertSecretStoreConfigV4 {
	secretConfig := client.UpsertSecretStoreConfigV4{}

	if config.AWS != nil {
		awsConfig := config.AWS
		secretConfig.AWS = client.NewOptNilAwsSecretConfigV4(client.AwsSecretConfigV4{
			Region:   awsConfig.Region.ValueString(),
			SecretID: awsConfig.SecretID.ValueString(),
		})
	} else if config.GCP != nil {
		gcpConfig := config.GCP
		secretConfig.Gcp = client.NewOptNilGcpSecretConfigV4(client.GcpSecretConfigV4{
			Project:  gcpConfig.Project.ValueString(),
			SecretID: gcpConfig.SecretID.ValueString(),
		})
	} else if config.Azure != nil {
		azureConfig := config.Azure
		secretConfig.Azure = client.NewOptNilAzureSecretConfigV4(client.AzureSecretConfigV4{
			VaultURL: azureConfig.VaultURL.ValueString(),
			Name:     azureConfig.Name.ValueString(),
		})
	} else if config.HashicorpVault != nil {
		vaultConfig := config.HashicorpVault
		secretConfig.HashicorpVault = client.NewOptNilHashicorpVaultSecretConfigV4(client.HashicorpVaultSecretConfigV4{
			SecretEngine: vaultConfig.SecretEngine.ValueString(),
			Path:         vaultConfig.Path.ValueString(),
		})
	} else if config.Kubernetes != nil {
		k8sConfig := config.Kubernetes
		secretConfig.Kubernetes = client.NewOptNilKubernetesSecretConfigV4(client.KubernetesSecretConfigV4{
			Namespace: k8sConfig.Namespace.ValueString(),
			Name:      k8sConfig.Name.ValueString(),
		})
	}

	return secretConfig
}

func convertSecretStoreConfigToModel(apiSecretConfig client.SecretStoreConfigV4) *SecretStoreConfig {
	secretConfig := &SecretStoreConfig{}

	if awsConfig, ok := apiSecretConfig.AWS.Get(); ok {
		secretConfig.AWS = &AWSSecretConfig{
			Region:   types.StringValue(awsConfig.Region),
			SecretID: types.StringValue(awsConfig.SecretID),
		}
	} else if gcpConfig, ok := apiSecretConfig.Gcp.Get(); ok {
		secretConfig.GCP = &GCPSecretConfig{
			Project:  types.StringValue(gcpConfig.Project),
			SecretID: types.StringValue(gcpConfig.SecretID),
		}
	} else if azureConfig, ok := apiSecretConfig.Azure.Get(); ok {
		secretConfig.Azure = &AzureSecretConfig{
			VaultURL: types.StringValue(azureConfig.VaultURL),
			Name:     types.StringValue(azureConfig.Name),
		}
	} else if vaultConfig, ok := apiSecretConfig.HashicorpVault.Get(); ok {
		secretConfig.HashicorpVault = &HashicorpVaultConfig{
			SecretEngine: types.StringValue(vaultConfig.SecretEngine),
			Path:         types.StringValue(vaultConfig.Path),
		}
	} else if k8sConfig, ok := apiSecretConfig.Kubernetes.Get(); ok {
		secretConfig.Kubernetes = &KubernetesSecretConfig{
			Namespace: types.StringValue(k8sConfig.Namespace),
			Name:      types.StringValue(k8sConfig.Name),
		}
	}

	return secretConfig
}

func convertIntegrationConfigToModel(ctx context.Context, integrationConfig map[string]jx.Raw) (types.Map, error) {
	configMap := make(map[string]attr.Value)

	for k, v := range integrationConfig {
		vstr, err := common.JxToString(v)
		if err != nil {
			return types.Map{}, fmt.Errorf("failed to decode integration config value: %v", err)
		}
		configMap[k] = types.StringValue(vstr)
	}

	result, diags := types.MapValueFrom(ctx, types.StringType, configMap)
	if diags.HasError() {
		return types.Map{}, fmt.Errorf("failed to parse integration config: %v", diags)
	}

	return result, nil
}

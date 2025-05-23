package services

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ConvertToIntegrationModel(ctx context.Context, integration *aponoapi.IntegrationTerraform) (*models.IntegrationModel, diag.Diagnostics) {
	metadataMapValue, diagnostics := types.MapValueFrom(ctx, types.StringType, integration.GetParams())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := models.IntegrationModel{}
	data.ID = types.StringValue(integration.GetId())
	data.Name = types.StringValue(integration.GetName())
	data.Type = types.StringValue(integration.GetType())
	data.ConnectorID = types.StringValue(integration.GetProvisionerId())
	data.Metadata = metadataMapValue

	if integration.CustomAccessDetails.IsSet() {
		data.CustomAccessDetails = types.StringValue(integration.GetCustomAccessDetails())
	}

	connectedResourceTypes, diagnostics := types.SetValueFrom(ctx, types.StringType, integration.GetConnectedResourceTypes())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}
	data.ConnectedResourceTypes = connectedResourceTypes

	secretConfig := integration.GetSecretConfig()
	switch secretConfig["type"] {
	case "AWS":
		data.AwsSecret = &models.AwsSecret{
			Region:   basetypes.NewStringValue(toString(secretConfig["region"])),
			SecretID: basetypes.NewStringValue(toString(secretConfig["secret_id"])),
		}
	case "GCP":
		data.GcpSecret = &models.GcpSecret{
			Project:  basetypes.NewStringValue(toString(secretConfig["project"])),
			SecretID: basetypes.NewStringValue(toString(secretConfig["secret_id"])),
		}
	case "KUBERNETES":
		data.KubernetesSecret = &models.KubernetesSecret{
			Namespace: basetypes.NewStringValue(toString(secretConfig["namespace"])),
			Name:      basetypes.NewStringValue(toString(secretConfig["name"])),
		}
	case "HASHICORP_VAULT":
		data.HashicorpVaultSecret = &models.HashicorpVaultSecret{
			SecretEngine: basetypes.NewStringValue(toString(secretConfig["secret_engine"])),
			Path:         basetypes.NewStringValue(toString(secretConfig["path"])),
		}
	case "AZURE":
		data.AzureSecret = &models.AzureSecret{
			VaultURL: basetypes.NewStringValue(toString(secretConfig["vault_url"])),
			Name:     basetypes.NewStringValue(toString(secretConfig["name"])),
		}
	case "APONO":
		paramsMap, diags := types.MapValueFrom(ctx, types.StringType, secretConfig["params"])
		if len(diags) > 0 {
			return nil, diags
		}
		data.AponoSecret = &models.AponoSecret{
			Params: paramsMap,
		}
	}

	data.ResourceOwnerMappings = ConvertResourceOwnersMappingToModel(integration.ResourceOwnersMappings)
	data.IntegrationOwners = ConvertIntegrationOwnerToData(integration.IntegrationOwners)

	return &data, nil
}

func ConvertIntegrationOwnerToData(owners []aponoapi.IntegrationOwnerTerraform) []models.IntegrationOwner {
	if owners == nil {
		return nil
	}
	var integrationOwners []models.IntegrationOwner
	for _, owner := range owners {
		var attributeValues []types.String
		for _, attributeValue := range owner.AttributeValue {
			attributeValues = append(attributeValues, basetypes.NewStringValue(attributeValue))
		}
		integrationOwners = append(integrationOwners, models.IntegrationOwner{
			IntegrationId:  basetypes.NewStringPointerValue(owner.IntegrationId.Get()),
			AttributeType:  basetypes.NewStringValue(owner.AttributeType),
			AttributeValue: attributeValues,
		})
	}
	return integrationOwners
}

func ConvertResourceOwnersMappingToModel(resourceOwnersMappings []aponoapi.ResourceOwnerMappingTerraform) []models.ResourceOwnerMapping {
	if resourceOwnersMappings == nil {
		return nil
	}
	var result []models.ResourceOwnerMapping
	for _, mapping := range resourceOwnersMappings {
		result = append(result, models.ResourceOwnerMapping{
			TagName:                basetypes.NewStringValue(mapping.TagName),
			AttributeType:          basetypes.NewStringValue(mapping.AttributeType),
			AttributeIntegrationId: basetypes.NewStringPointerValue(mapping.AttributeIntegrationId.Get()),
		})
	}

	return result
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func ConvertMappingsArrayToModel(resourceOwnerMappings []models.ResourceOwnerMapping) []aponoapi.ResourceOwnerMappingTerraform {
	var result []aponoapi.ResourceOwnerMappingTerraform
	for _, mapping := range resourceOwnerMappings {
		result = append(result, resourceOwnerMappingToModel(mapping))
	}
	return result
}

func resourceOwnerMappingToModel(mapping models.ResourceOwnerMapping) aponoapi.ResourceOwnerMappingTerraform {
	return aponoapi.ResourceOwnerMappingTerraform{
		TagName:                mapping.TagName.ValueString(),
		AttributeType:          mapping.AttributeType.ValueString(),
		AttributeIntegrationId: *aponoapi.NewNullableString(mapping.AttributeIntegrationId.ValueStringPointer()),
	}
}

func IntegrationOwnersToModel(integrationOwners []models.IntegrationOwner) []aponoapi.IntegrationOwnerTerraform {
	var owners []aponoapi.IntegrationOwnerTerraform
	for _, owner := range integrationOwners {
		owners = append(owners, integrationOwnerToModel(owner))
	}
	return owners
}

func integrationOwnerToModel(owner models.IntegrationOwner) aponoapi.IntegrationOwnerTerraform {
	return aponoapi.IntegrationOwnerTerraform{
		IntegrationId:  *aponoapi.NewNullableString(owner.IntegrationId.ValueStringPointer()),
		AttributeType:  owner.AttributeType.ValueString(),
		AttributeValue: utils.ConvertStringArray(owner.AttributeValue),
	}
}

package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceIntegrationModel struct {
	ID                     types.String         `tfsdk:"id"`
	Name                   types.String         `tfsdk:"name"`
	Type                   types.String         `tfsdk:"type"`
	ConnectorID            types.String         `tfsdk:"connector_id"`
	ConnectedResourceTypes types.List           `tfsdk:"connected_resource_types"`
	IntegrationConfig      types.Map            `tfsdk:"integration_config"`
	SecretStoreConfig      *SecretStoreConfig   `tfsdk:"secret_store_config"`
	CustomAccessDetails    types.String         `tfsdk:"custom_access_details"`
	Owner                  *OwnerConfig         `tfsdk:"owner"`
	OwnersMapping          *OwnersMappingConfig `tfsdk:"owners_mapping"`
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

func ResourceIntegrationModelToCreateRequest(ctx context.Context, model ResourceIntegrationModel) (*client.CreateIntegrationV4, error) {
	req := &client.CreateIntegrationV4{
		Name: model.Name.ValueString(),
		Type: model.Type.ValueString(),
	}

	req.ConnectorID.SetTo(model.ConnectorID.ValueString())

	var err error
	if req.ConnectedResourceTypes, err = getConnectedResourceTypes(ctx, model); err != nil {
		return nil, err
	}

	if !model.IntegrationConfig.IsNull() {
		req.IntegrationConfig, err = getIntegrationConfig(model)
		if err != nil {
			return nil, err
		}
	}

	if model.SecretStoreConfig != nil {
		req.SecretStoreConfig.SetTo(upsertSecretStoreConfig(model.SecretStoreConfig))
	}

	if !model.CustomAccessDetails.IsNull() {
		req.CustomAccessDetails.SetTo(model.CustomAccessDetails.ValueString())
	}

	if model.Owner != nil {
		owner, err := getOwnerConfig(ctx, model.Owner)
		if err != nil {
			return nil, err
		}
		req.Owner.SetTo(owner)
	}

	if model.OwnersMapping != nil {
		req.OwnersMapping.SetTo(getOwnersMappingConfig(model.OwnersMapping))
	}

	return req, nil
}

func ResourceIntegrationModelToUpdateRequest(ctx context.Context, model ResourceIntegrationModel) (*client.UpdateIntegrationV4, error) {
	req := &client.UpdateIntegrationV4{
		Name: model.Name.ValueString(),
	}

	var err error
	if !model.IntegrationConfig.IsNull() {
		req.IntegrationConfig, err = getIntegrationConfig(model)
		if err != nil {
			return nil, err
		}
	}

	if req.ConnectedResourceTypes, err = getConnectedResourceTypes(ctx, model); err != nil {
		return nil, err
	}

	if model.SecretStoreConfig != nil {
		req.SecretStoreConfig.SetTo(upsertSecretStoreConfig(model.SecretStoreConfig))
	}

	if !model.CustomAccessDetails.IsNull() {
		req.CustomAccessDetails.SetTo(model.CustomAccessDetails.ValueString())
	}

	if model.Owner != nil {
		owner, err := getOwnerConfig(ctx, model.Owner)
		if err != nil {
			return nil, err
		}
		req.Owner.SetTo(owner)
	}

	if model.OwnersMapping != nil {
		req.OwnersMapping.SetTo(getOwnersMappingConfig(model.OwnersMapping))
	}

	return req, nil
}

func getIntegrationConfig(model ResourceIntegrationModel) (map[string]jx.Raw, error) {
	integrationConfig := make(map[string]jx.Raw)
	for k, v := range model.IntegrationConfig.Elements() {
		strVal, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("failed to assert type for integration config value")
		}
		integrationConfig[k] = common.StringToJx(strVal.ValueString())
	}
	return integrationConfig, nil
}

func getConnectedResourceTypes(ctx context.Context, model ResourceIntegrationModel) (client.OptNilStringArray, error) {
	var result client.OptNilStringArray
	var connectedResourceTypes []string
	diags := model.ConnectedResourceTypes.ElementsAs(ctx, &connectedResourceTypes, false)
	if diags.HasError() {
		return result, fmt.Errorf("failed to parse connected resource types: %v", diags)
	}
	result.SetTo(connectedResourceTypes)
	return result, nil
}

func getOwnerConfig(ctx context.Context, ownerConfig *OwnerConfig) (client.UpsertOwnerV4, error) {
	var values []string
	diags := ownerConfig.Values.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return client.UpsertOwnerV4{}, fmt.Errorf("failed to parse owner values: %v", diags)
	}

	owner := client.UpsertOwnerV4{
		AttributeType:  ownerConfig.Type.ValueString(),
		AttributeValue: values,
	}

	if !ownerConfig.SourceIntegrationName.IsNull() {
		owner.SourceIntegrationReference.SetTo(ownerConfig.SourceIntegrationName.ValueString())
	}

	return owner, nil
}

func getOwnersMappingConfig(mappingConfig *OwnersMappingConfig) client.UpsertOwnerMappingV4 {
	ownersMapping := client.UpsertOwnerMappingV4{
		KeyName:       mappingConfig.KeyName.ValueString(),
		AttributeType: mappingConfig.AttributeType.ValueString(),
	}

	if !mappingConfig.SourceIntegrationName.IsNull() {
		ownersMapping.SourceIntegrationReference.SetTo(mappingConfig.SourceIntegrationName.ValueString())
	}

	return ownersMapping
}

func ResourceIntegrationToModel(ctx context.Context, integration *client.IntegrationV4) (*ResourceIntegrationModel, error) {
	model := &ResourceIntegrationModel{
		ID:   types.StringValue(integration.ID),
		Name: types.StringValue(integration.Name),
		Type: types.StringValue(integration.Type),
	}

	model.ConnectorID = types.StringValue(integration.ConnectorID.Value)

	if connectedResourceTypes, ok := integration.ConnectedResourceTypes.Get(); ok {
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
		integrationConfig, err := convertIntegrationConfigToModel(ctx, integration.IntegrationConfig)
		if err != nil {
			return nil, err
		}
		model.IntegrationConfig = integrationConfig
	}

	if val, ok := integration.SecretStoreConfig.Get(); ok {
		model.SecretStoreConfig = convertSecretStoreConfigToModel(val)
	}

	if val, ok := integration.CustomAccessDetails.Get(); ok {
		model.CustomAccessDetails = types.StringValue(val)
	}

	if ownerData, ok := integration.Owner.Get(); ok {
		values, diags := types.ListValueFrom(ctx, types.StringType, ownerData.AttributeValue)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse owner values: %v", diags)
		}

		ownerConfig := &OwnerConfig{
			Type:   types.StringValue(ownerData.AttributeType),
			Values: values,
		}

		if val, ok := ownerData.SourceIntegrationName.Get(); ok {
			ownerConfig.SourceIntegrationName = types.StringValue(val)
		}

		model.Owner = ownerConfig
	}

	if mappingData, ok := integration.OwnersMapping.Get(); ok {
		ownersMappingConfig := &OwnersMappingConfig{
			KeyName:       types.StringValue(mappingData.KeyName),
			AttributeType: types.StringValue(mappingData.AttributeType),
		}

		if val, ok := mappingData.SourceIntegrationName.Get(); ok {
			ownersMappingConfig.SourceIntegrationName = types.StringValue(val)
		}

		model.OwnersMapping = ownersMappingConfig
	}

	return model, nil
}

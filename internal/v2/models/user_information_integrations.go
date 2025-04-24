package models

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AponoUserInformationIntegrationsDataSourceModel struct {
	Name         types.String                      `tfsdk:"name"`
	Type         types.String                      `tfsdk:"type"`
	Integrations []UserInformationIntegrationModel `tfsdk:"integrations"`
}

type UserInformationIntegrationModel struct {
	ID                types.String       `tfsdk:"id"`
	Name              types.String       `tfsdk:"name"`
	Type              types.String       `tfsdk:"type"`
	Category          types.String       `tfsdk:"category"`
	Status            types.String       `tfsdk:"status"`
	LastSyncTime      types.String       `tfsdk:"last_sync_time"`
	IntegrationConfig types.Map          `tfsdk:"integration_config"`
	SecretConfig      *SecretStoreConfig `tfsdk:"secret_config"`
}

func UserInformationIntegrationToModal(ctx context.Context, integration *client.IntegrationV4) (*UserInformationIntegrationModel, error) {
	model := &UserInformationIntegrationModel{
		ID:       types.StringValue(integration.ID),
		Name:     types.StringValue(integration.Name),
		Type:     types.StringValue(integration.Type),
		Category: types.StringValue(integration.Category),
		Status:   types.StringValue(integration.Status),
	}

	if lastSyncTime, ok := integration.LastSyncTime.Get(); ok {
		model.LastSyncTime = types.StringValue(lastSyncTime.UTC().Format("2006-01-02T15:04:05Z"))
	}

	if integration.IntegrationConfig != nil {
		integrationConfig, err := convertIntegrationConfigToModel(ctx, integration.IntegrationConfig)
		if err != nil {
			return nil, err
		}
		model.IntegrationConfig = integrationConfig
	}

	if val, ok := integration.SecretStoreConfig.Get(); ok {
		model.SecretConfig = convertSecretStoreConfigToModel(val)
	}

	return model, nil
}

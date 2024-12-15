package services

import (
	"context"
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"math/big"
)

func ConvertManualWebhookApiToTerraformModel(ctx context.Context, manualWebhook *aponoapi.WebhookManualTriggerTerraformModel) (*models.ManualWebhookModel, diag.Diagnostics) {
	manualWebhookModel := models.ManualWebhookModel{}
	manualWebhookModel.ID = types.StringValue(manualWebhook.GetId())
	manualWebhookModel.Name = types.StringValue(manualWebhook.GetName())
	manualWebhookModel.Active = types.BoolValue(manualWebhook.GetActive())

	manualWebhookType, diagnostics := manualWebhookTypeToModel(ctx, manualWebhook.GetType())
	if diagnostics != nil {
		return nil, diagnostics
	}
	manualWebhookModel.Type = *manualWebhookType

	if manualWebhook.BodyTemplate.IsSet() {
		manualWebhookModel.BodyTemplate = types.StringValue(manualWebhook.GetBodyTemplate())
	}

	if manualWebhook.ResponseValidators != nil {
		manualWebhookModel.ResponseValidators = responseValidatorsToModel(manualWebhook.GetResponseValidators())
	}

	if manualWebhook.TimeoutInSec.IsSet() {
		manualWebhookModel.TimeoutInSec = types.NumberValue(big.NewFloat(float64(manualWebhook.GetTimeoutInSec())))
	}

	if manualWebhook.AuthenticationConfig.IsSet() {
		authenticationConfig, diagnostics := authenticationConfigToModel(ctx, manualWebhook.GetAuthenticationConfig())
		if diagnostics != nil {
			return nil, diagnostics
		}
		manualWebhookModel.AuthenticationConfig = authenticationConfig
	}

	if manualWebhook.CustomValidationErrorMessage.IsSet() {
		manualWebhookModel.CustomValidationErrorMessage = types.StringValue(manualWebhook.GetCustomValidationErrorMessage())
	}

	return &manualWebhookModel, nil
}

func responseValidatorsToModel(responseValidators []aponoapi.WebhookResponseValidatorTerraformModel) []models.ManualWebhookResponseValidatorModel {
	var validators []models.ManualWebhookResponseValidatorModel
	for _, responseValidator := range responseValidators {
		validators = append(validators, responseValidatorToModel(responseValidator))
	}
	return validators
}

func responseValidatorToModel(responseValidator aponoapi.WebhookResponseValidatorTerraformModel) models.ManualWebhookResponseValidatorModel {
	var expectedValuesList []types.String
	for _, expectedValues := range responseValidator.ExpectedValues {
		expectedValuesList = append(expectedValuesList, basetypes.NewStringValue(expectedValues))
	}
	return models.ManualWebhookResponseValidatorModel{
		JsonPath:       types.StringValue(responseValidator.GetJsonPath()),
		ExpectedValues: expectedValuesList,
	}
}

func manualWebhookTypeToModel(ctx context.Context, manualWebhookType aponoapi.WebhookTypeTerraformModel) (*models.ManualWebhookTypeModel, diag.Diagnostics) {
	if manualWebhookType.HttpRequest.IsSet() {
		httpRequest, diagnostics := manualWebhookHttpRequestTypeToModel(ctx, manualWebhookType.GetHttpRequest())
		return &models.ManualWebhookTypeModel{
			HttpRequest: httpRequest,
		}, diagnostics
	} else if manualWebhookType.Integration.IsSet() {
		return &models.ManualWebhookTypeModel{
			Integration: manualWebhookIntegrationTypeToModel(manualWebhookType.GetIntegration()),
		}, nil
	}

	diagnostics := diag.Diagnostics{}
	diagnostics.AddError("Client Error", "manual webhook type is not set to either HttpRequest or Integration")
	return nil, diagnostics
}

func manualWebhookHttpRequestTypeToModel(ctx context.Context, httpRequestType aponoapi.WebhookTypeTerraformModelHttpRequest) (*models.ManualWebhookHttpRequestTypeModel, diag.Diagnostics) {
	httpRequest := models.ManualWebhookHttpRequestTypeModel{
		Url:    types.StringValue(httpRequestType.GetUrl()),
		Method: types.StringValue(string(httpRequestType.GetMethod())),
	}

	headersMapValue, diagnostics := types.MapValueFrom(ctx, types.StringType, httpRequestType.GetHeaders())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	httpRequest.Headers = headersMapValue

	return &httpRequest, nil
}

func manualWebhookIntegrationTypeToModel(integrationType aponoapi.WebhookTypeTerraformModelIntegration) *models.ManualWebhookIntegrationTypeModel {
	return &models.ManualWebhookIntegrationTypeModel{
		IntegrationId: types.StringValue(integrationType.GetIntegrationId()),
		ActionName:    types.StringValue(integrationType.GetActionName()),
	}
}

func authenticationConfigToModel(ctx context.Context, authenticationConfig aponoapi.WebhookManualTriggerTerraformModelAuthenticationConfig) (*models.ManualWebhookAuthenticationConfigModel, diag.Diagnostics) {
	if !authenticationConfig.Oauth.IsSet() {
		return nil, nil
	}

	oauth, diagnostics := webhookOAuthConfigToModel(ctx, authenticationConfig.GetOauth())
	if diagnostics != nil {
		return nil, diagnostics
	}

	return &models.ManualWebhookAuthenticationConfigModel{
		Type:  types.StringValue(authenticationConfig.GetType()),
		Oauth: oauth,
	}, nil
}

func webhookOAuthConfigToModel(ctx context.Context, oauthConfig aponoapi.WebhookAuthenticationConfigTerraformModelOauth) (*models.WebhookOAuthConfigModel, diag.Diagnostics) {
	scopes, diagnostics := types.ListValueFrom(ctx, types.StringType, oauthConfig.GetScopes())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	return &models.WebhookOAuthConfigModel{
		ClientId:         types.StringValue(oauthConfig.GetClientId()),
		ClientSecret:     types.StringValue(oauthConfig.GetClientSecret()),
		TokenEndpointUrl: types.StringValue(oauthConfig.GetTokenEndpointUrl()),
		Scopes:           scopes,
	}, nil
}

func ConvertManualWebhookTerraformModelToUpsertApi(ctx context.Context, manualWebhook *models.ManualWebhookModel) (*aponoapi.WebhookManualTriggerUpsertTerraformModel, diag.Diagnostics) {
	//dataIntegrationTargets, diagnostics := convertIntegrationTargetsTerraformModelToApi(ctx, manualWebhook)
	//if len(diagnostics) > 0 {
	//	return nil, diagnostics
	//}
	manualWebhookType, diagnostics := manualWebhookTypeToApi(manualWebhook.Type)
	if diagnostics != nil {
		return nil, diagnostics
	}

	data := aponoapi.WebhookManualTriggerUpsertTerraformModel{
		Name:   manualWebhook.Name.ValueString(),
		Active: manualWebhook.Active.ValueBool(),
		//	AuthenticationConfig         NullableWebhookManualTriggerTerraformModelAuthenticationConfig `json:"authentication_config,omitempty"`
		BodyTemplate:                 *aponoapi.NewNullableString(manualWebhook.BodyTemplate.ValueStringPointer()),
		CustomValidationErrorMessage: *aponoapi.NewNullableString(manualWebhook.CustomValidationErrorMessage.ValueStringPointer()),
		ResponseValidators:           responseValidatorsToApi(manualWebhook.ResponseValidators),
		AuthenticationConfig:         *aponoapi.NewNullableWebhookManualTriggerTerraformModelAuthenticationConfig(authenticationConfigToApi(manualWebhook.AuthenticationConfig)),
	}

	data.Type = *manualWebhookType

	if !manualWebhook.TimeoutInSec.IsNull() {
		timeoutInSec, _ := manualWebhook.TimeoutInSec.ValueBigFloat().Int64()
		timeoutInSecInt32 := int32(timeoutInSec)
		data.TimeoutInSec = *aponoapi.NewNullableInt32(&timeoutInSecInt32)
	}

	return &data, nil
}

func manualWebhookTypeToApi(manualWebhookType models.ManualWebhookTypeModel) (*aponoapi.WebhookTypeTerraformModel, diag.Diagnostics) {
	if manualWebhookType.HttpRequest != nil {
		httpRequest, diagnostics := manualWebhookHttpRequestTypeToApi(*manualWebhookType.HttpRequest)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
		return &aponoapi.WebhookTypeTerraformModel{
			HttpRequest: *aponoapi.NewNullableWebhookTypeTerraformModelHttpRequest(httpRequest),
		}, nil
	} else if manualWebhookType.Integration != nil {
		integration := manualWebhookIntegrationTypeToApi(*manualWebhookType.Integration)
		return &aponoapi.WebhookTypeTerraformModel{
			Integration: *aponoapi.NewNullableWebhookTypeTerraformModelIntegration(integration),
		}, nil
	}

	diagnostics := diag.Diagnostics{}
	diagnostics.AddError("Client Error", "manual webhook type is not set to either HttpRequest or Integration")
	return nil, diagnostics
}

func manualWebhookHttpRequestTypeToApi(httpRequestType models.ManualWebhookHttpRequestTypeModel) (*aponoapi.WebhookTypeTerraformModelHttpRequest, diag.Diagnostics) {
	data := aponoapi.WebhookTypeTerraformModelHttpRequest{
		Url:    httpRequestType.Url.ValueString(),
		Method: aponoapi.WebhookMethodTerraformModel(httpRequestType.Method.ValueString()),
	}

	if !httpRequestType.Headers.IsNull() {
		headers, diagnostics := ConvertTypesMapToStringMap(httpRequestType.Headers)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
		data.Headers = headers
	}

	return &data, nil
}

func ConvertTypesMapToStringMap(input types.Map) (map[string]string, diag.Diagnostics) {
	// Prepare the map for conversion
	output := make(map[string]string)
	diagnostics := diag.Diagnostics{}
	// Convert the map
	for key, value := range input.Elements() {
		// Ensure each value is a types.String
		strValue, ok := value.(types.String)
		if !ok {
			diagnostics.AddError("Client Error", fmt.Sprintf("value for key %s is not a string", key))
			return nil, diagnostics
		}
		// Check if the value is null or unknown
		if strValue.IsNull() || strValue.IsUnknown() {
			diagnostics.AddError("Client Error", fmt.Sprintf("value for key %s is null or unknown", key))
			return nil, diagnostics
		}
		// Assign the value to the output map
		output[key] = strValue.ValueString()
	}

	return output, nil
}

func manualWebhookIntegrationTypeToApi(integration models.ManualWebhookIntegrationTypeModel) *aponoapi.WebhookTypeTerraformModelIntegration {
	return &aponoapi.WebhookTypeTerraformModelIntegration{
		IntegrationId: integration.IntegrationId.ValueString(),
		ActionName:    integration.ActionName.ValueString(),
	}
}

func responseValidatorsToApi(responseValidators []models.ManualWebhookResponseValidatorModel) []aponoapi.WebhookResponseValidatorTerraformModel {
	var validators []aponoapi.WebhookResponseValidatorTerraformModel
	for _, responseValidator := range responseValidators {
		validators = append(validators, responseValidatorToApi(responseValidator))
	}
	return validators
}

func responseValidatorToApi(responseValidator models.ManualWebhookResponseValidatorModel) aponoapi.WebhookResponseValidatorTerraformModel {
	var expectedValuesList []string
	for _, expectedValues := range responseValidator.ExpectedValues {
		expectedValuesList = append(expectedValuesList, expectedValues.ValueString())
	}
	return aponoapi.WebhookResponseValidatorTerraformModel{
		JsonPath:       responseValidator.JsonPath.ValueString(),
		ExpectedValues: expectedValuesList,
	}
}

func authenticationConfigToApi(authenticationConfig *models.ManualWebhookAuthenticationConfigModel) *aponoapi.WebhookManualTriggerTerraformModelAuthenticationConfig {
	if authenticationConfig == nil {
		return nil
	}

	return &aponoapi.WebhookManualTriggerTerraformModelAuthenticationConfig{
		Type:  authenticationConfig.Type.ValueString(),
		Oauth: *aponoapi.NewNullableWebhookAuthenticationConfigTerraformModelOauth(webhookOAuthConfigToApi(authenticationConfig.Oauth)),
	}
}

func webhookOAuthConfigToApi(oauthConfig *models.WebhookOAuthConfigModel) *aponoapi.WebhookAuthenticationConfigTerraformModelOauth {
	var scopes []string
	for _, scope := range oauthConfig.Scopes.Elements() {
		scopes = append(scopes, scope.String())
	}

	return &aponoapi.WebhookAuthenticationConfigTerraformModelOauth{
		ClientId:         oauthConfig.ClientId.ValueString(),
		ClientSecret:     oauthConfig.ClientSecret.ValueString(),
		TokenEndpointUrl: oauthConfig.TokenEndpointUrl.ValueString(),
		Scopes:           scopes,
	}
}

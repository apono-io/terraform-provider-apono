package services

import (
	"context"
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
		manualWebhookModel.BodyTemplate = types.StringPointerValue(manualWebhook.BodyTemplate.Get())
	}

	if manualWebhook.ResponseValidators != nil {
		manualWebhookModel.ResponseValidators = responseValidatorsToModel(manualWebhook.GetResponseValidators())
	}

	if manualWebhook.TimeoutInSec.IsSet() {
		ptr := convertInt32PtrToInt64Ptr(manualWebhook.TimeoutInSec.Get())
		manualWebhookModel.TimeoutInSec = types.Int64PointerValue(ptr)
	}

	if manualWebhook.AuthenticationConfig.IsSet() {
		authenticationConfig := authenticationConfigToModel(manualWebhook.GetAuthenticationConfig())
		manualWebhookModel.AuthenticationConfig = authenticationConfig
	}

	if manualWebhook.CustomValidationErrorMessage.IsSet() {
		manualWebhookModel.CustomValidationErrorMessage = types.StringPointerValue(manualWebhook.CustomValidationErrorMessage.Get())
	}

	return &manualWebhookModel, nil
}

func convertInt32PtrToInt64Ptr(ptr *int32) *int64 {
	if ptr == nil {
		return nil
	}

	val := int64(*ptr)
	return &val
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

func authenticationConfigToModel(authenticationConfig aponoapi.WebhookManualTriggerTerraformModelAuthenticationConfig) *models.ManualWebhookAuthenticationConfigModel {
	if !authenticationConfig.Oauth.IsSet() {
		return nil
	}

	oauth := webhookOAuthConfigToModel(authenticationConfig.GetOauth())
	return &models.ManualWebhookAuthenticationConfigModel{
		Oauth: oauth,
	}
}

func webhookOAuthConfigToModel(oauthConfig aponoapi.WebhookAuthenticationConfigTerraformModelOauth) *models.WebhookOAuthConfigModel {
	scopes := []types.String{}
	for _, scope := range oauthConfig.GetScopes() {
		scopes = append(scopes, types.StringValue(scope))
	}

	return &models.WebhookOAuthConfigModel{
		ClientId:         types.StringValue(oauthConfig.GetClientId()),
		ClientSecret:     types.StringValue(oauthConfig.GetClientSecret()),
		TokenEndpointUrl: types.StringValue(oauthConfig.GetTokenEndpointUrl()),
		Scopes:           scopes,
	}
}

func ConvertManualWebhookTerraformModelToUpsertApi(manualWebhook *models.ManualWebhookModel) (*aponoapi.WebhookManualTriggerUpsertTerraformModel, diag.Diagnostics) {
	manualWebhookType, diagnostics := manualWebhookTypeToApi(manualWebhook.Type)
	if diagnostics != nil {
		return nil, diagnostics
	}

	var bodyTemplate *string
	if !manualWebhook.BodyTemplate.IsNull() && !manualWebhook.BodyTemplate.IsUnknown() {
		bodyTemplate = manualWebhook.BodyTemplate.ValueStringPointer()
	}
	var customValidationErrorMessage *string
	if !manualWebhook.CustomValidationErrorMessage.IsNull() && !manualWebhook.CustomValidationErrorMessage.IsUnknown() {
		customValidationErrorMessage = manualWebhook.CustomValidationErrorMessage.ValueStringPointer()
	}

	var timeoutInSec *int32
	if !manualWebhook.TimeoutInSec.IsNull() && !manualWebhook.TimeoutInSec.IsUnknown() {
		timeoutInSecInt32 := int32(manualWebhook.TimeoutInSec.ValueInt64())
		timeoutInSec = &timeoutInSecInt32
	}

	data := aponoapi.WebhookManualTriggerUpsertTerraformModel{
		Name:                         manualWebhook.Name.ValueString(),
		Active:                       manualWebhook.Active.ValueBool(),
		Type:                         *manualWebhookType,
		BodyTemplate:                 *aponoapi.NewNullableString(bodyTemplate),
		ResponseValidators:           responseValidatorsToApi(manualWebhook.ResponseValidators),
		TimeoutInSec:                 *aponoapi.NewNullableInt32(timeoutInSec),
		AuthenticationConfig:         *aponoapi.NewNullableWebhookManualTriggerTerraformModelAuthenticationConfig(authenticationConfigToApi(manualWebhook.AuthenticationConfig)),
		CustomValidationErrorMessage: *aponoapi.NewNullableString(customValidationErrorMessage),
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

	if !httpRequestType.Headers.IsNull() && !httpRequestType.Headers.IsUnknown() {
		headers, diagnostics := convertTypesMapToStringMap(httpRequestType.Headers)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
		data.Headers = headers
	}

	return &data, nil
}

func convertTypesMapToStringMap(input types.Map) (map[string]string, diag.Diagnostics) {
	output := make(map[string]string)
	diagnostics := diag.Diagnostics{}

	for key, value := range input.Elements() {
		strValue, ok := value.(types.String)
		if !ok {
			diagnostics.AddError("Client Error", fmt.Sprintf("value for key %s is not a string", key))
			return nil, diagnostics
		}
		if strValue.IsNull() || strValue.IsUnknown() {
			diagnostics.AddError("Client Error", fmt.Sprintf("value for key %s is null or unknown", key))
			return nil, diagnostics
		}
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
		Oauth: *aponoapi.NewNullableWebhookAuthenticationConfigTerraformModelOauth(webhookOAuthConfigToApi(authenticationConfig.Oauth)),
	}
}

func webhookOAuthConfigToApi(oauthConfig *models.WebhookOAuthConfigModel) *aponoapi.WebhookAuthenticationConfigTerraformModelOauth {
	scopes := []string{}
	for _, scope := range oauthConfig.Scopes {
		scopes = append(scopes, scope.ValueString())
	}

	return &aponoapi.WebhookAuthenticationConfigTerraformModelOauth{
		ClientId:         oauthConfig.ClientId.ValueString(),
		ClientSecret:     oauthConfig.ClientSecret.ValueString(),
		TokenEndpointUrl: oauthConfig.TokenEndpointUrl.ValueString(),
		Scopes:           scopes,
	}
}

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ManualWebhookModel describes the resource data model.
type ManualWebhookModel struct {
	ID                           types.String                            `tfsdk:"id"`
	Name                         types.String                            `tfsdk:"name"`
	Active                       types.Bool                              `tfsdk:"active"`
	Type                         ManualWebhookTypeModel                  `tfsdk:"type"`
	BodyTemplate                 types.String                            `tfsdk:"body_template"`
	ResponseValidators           []ManualWebhookResponseValidatorModel   `tfsdk:"response_validators"`
	TimeoutInSec                 types.Int64                             `tfsdk:"timeout_in_sec"`
	AuthenticationConfig         *ManualWebhookAuthenticationConfigModel `tfsdk:"authentication_config"`
	CustomValidationErrorMessage types.String                            `tfsdk:"custom_validation_error_message"`
}

type ManualWebhookTypeModel struct {
	HttpRequest *ManualWebhookHttpRequestTypeModel `tfsdk:"http_request"`
	Integration *ManualWebhookIntegrationTypeModel `tfsdk:"integration"`
}

type ManualWebhookHttpRequestTypeModel struct {
	Url     types.String `tfsdk:"url"`
	Method  types.String `tfsdk:"method"`
	Headers types.Map    `tfsdk:"headers"`
}

type ManualWebhookIntegrationTypeModel struct {
	IntegrationId types.String `tfsdk:"integration_id"`
	ActionName    types.String `tfsdk:"action_name"`
}

type ManualWebhookResponseValidatorModel struct {
	JsonPath       types.String   `tfsdk:"json_path"`
	ExpectedValues []types.String `tfsdk:"expected_values"`
}

type ManualWebhookAuthenticationConfigModel struct {
	Oauth *WebhookOAuthConfigModel `tfsdk:"oauth"`
}

type WebhookOAuthConfigModel struct {
	ClientId         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
	TokenEndpointUrl types.String `tfsdk:"token_endpoint_url"`
	Scopes           types.List   `tfsdk:"scopes"`
}

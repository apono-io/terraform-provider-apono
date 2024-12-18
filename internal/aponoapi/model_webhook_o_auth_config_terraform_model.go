/*
Apono

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package aponoapi

import (
	"encoding/json"
	"fmt"
)

// checks if the WebhookOAuthConfigTerraformModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookOAuthConfigTerraformModel{}

// WebhookOAuthConfigTerraformModel struct for WebhookOAuthConfigTerraformModel
type WebhookOAuthConfigTerraformModel struct {
	ClientId         string   `json:"client_id"`
	ClientSecret     string   `json:"client_secret"`
	TokenEndpointUrl string   `json:"token_endpoint_url"`
	Scopes           []string `json:"scopes"`
}

type _WebhookOAuthConfigTerraformModel WebhookOAuthConfigTerraformModel

// NewWebhookOAuthConfigTerraformModel instantiates a new WebhookOAuthConfigTerraformModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookOAuthConfigTerraformModel(clientId string, clientSecret string, tokenEndpointUrl string, scopes []string) *WebhookOAuthConfigTerraformModel {
	this := WebhookOAuthConfigTerraformModel{}
	this.ClientId = clientId
	this.ClientSecret = clientSecret
	this.TokenEndpointUrl = tokenEndpointUrl
	this.Scopes = scopes
	return &this
}

// NewWebhookOAuthConfigTerraformModelWithDefaults instantiates a new WebhookOAuthConfigTerraformModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookOAuthConfigTerraformModelWithDefaults() *WebhookOAuthConfigTerraformModel {
	this := WebhookOAuthConfigTerraformModel{}
	return &this
}

// GetClientId returns the ClientId field value
func (o *WebhookOAuthConfigTerraformModel) GetClientId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value
// and a boolean to check if the value has been set.
func (o *WebhookOAuthConfigTerraformModel) GetClientIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClientId, true
}

// SetClientId sets field value
func (o *WebhookOAuthConfigTerraformModel) SetClientId(v string) {
	o.ClientId = v
}

// GetClientSecret returns the ClientSecret field value
func (o *WebhookOAuthConfigTerraformModel) GetClientSecret() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClientSecret
}

// GetClientSecretOk returns a tuple with the ClientSecret field value
// and a boolean to check if the value has been set.
func (o *WebhookOAuthConfigTerraformModel) GetClientSecretOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClientSecret, true
}

// SetClientSecret sets field value
func (o *WebhookOAuthConfigTerraformModel) SetClientSecret(v string) {
	o.ClientSecret = v
}

// GetTokenEndpointUrl returns the TokenEndpointUrl field value
func (o *WebhookOAuthConfigTerraformModel) GetTokenEndpointUrl() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TokenEndpointUrl
}

// GetTokenEndpointUrlOk returns a tuple with the TokenEndpointUrl field value
// and a boolean to check if the value has been set.
func (o *WebhookOAuthConfigTerraformModel) GetTokenEndpointUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TokenEndpointUrl, true
}

// SetTokenEndpointUrl sets field value
func (o *WebhookOAuthConfigTerraformModel) SetTokenEndpointUrl(v string) {
	o.TokenEndpointUrl = v
}

// GetScopes returns the Scopes field value
func (o *WebhookOAuthConfigTerraformModel) GetScopes() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Scopes
}

// GetScopesOk returns a tuple with the Scopes field value
// and a boolean to check if the value has been set.
func (o *WebhookOAuthConfigTerraformModel) GetScopesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Scopes, true
}

// SetScopes sets field value
func (o *WebhookOAuthConfigTerraformModel) SetScopes(v []string) {
	o.Scopes = v
}

func (o WebhookOAuthConfigTerraformModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookOAuthConfigTerraformModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["client_id"] = o.ClientId
	toSerialize["client_secret"] = o.ClientSecret
	toSerialize["token_endpoint_url"] = o.TokenEndpointUrl
	toSerialize["scopes"] = o.Scopes
	return toSerialize, nil
}

func (o *WebhookOAuthConfigTerraformModel) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"client_id",
		"client_secret",
		"token_endpoint_url",
		"scopes",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varWebhookOAuthConfigTerraformModel := _WebhookOAuthConfigTerraformModel{}

	err = json.Unmarshal(bytes, &varWebhookOAuthConfigTerraformModel)

	if err != nil {
		return err
	}

	*o = WebhookOAuthConfigTerraformModel(varWebhookOAuthConfigTerraformModel)

	return err
}

type NullableWebhookOAuthConfigTerraformModel struct {
	value *WebhookOAuthConfigTerraformModel
	isSet bool
}

func (v NullableWebhookOAuthConfigTerraformModel) Get() *WebhookOAuthConfigTerraformModel {
	return v.value
}

func (v *NullableWebhookOAuthConfigTerraformModel) Set(val *WebhookOAuthConfigTerraformModel) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookOAuthConfigTerraformModel) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookOAuthConfigTerraformModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookOAuthConfigTerraformModel(val *WebhookOAuthConfigTerraformModel) *NullableWebhookOAuthConfigTerraformModel {
	return &NullableWebhookOAuthConfigTerraformModel{value: val, isSet: true}
}

func (v NullableWebhookOAuthConfigTerraformModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookOAuthConfigTerraformModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
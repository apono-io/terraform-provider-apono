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

// checks if the WebhookIntegrationTypeTerraformModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookIntegrationTypeTerraformModel{}

// WebhookIntegrationTypeTerraformModel struct for WebhookIntegrationTypeTerraformModel
type WebhookIntegrationTypeTerraformModel struct {
	IntegrationId string `json:"integration_id"`
	ActionName    string `json:"action_name"`
}

type _WebhookIntegrationTypeTerraformModel WebhookIntegrationTypeTerraformModel

// NewWebhookIntegrationTypeTerraformModel instantiates a new WebhookIntegrationTypeTerraformModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookIntegrationTypeTerraformModel(integrationId string, actionName string) *WebhookIntegrationTypeTerraformModel {
	this := WebhookIntegrationTypeTerraformModel{}
	this.IntegrationId = integrationId
	this.ActionName = actionName
	return &this
}

// NewWebhookIntegrationTypeTerraformModelWithDefaults instantiates a new WebhookIntegrationTypeTerraformModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookIntegrationTypeTerraformModelWithDefaults() *WebhookIntegrationTypeTerraformModel {
	this := WebhookIntegrationTypeTerraformModel{}
	return &this
}

// GetIntegrationId returns the IntegrationId field value
func (o *WebhookIntegrationTypeTerraformModel) GetIntegrationId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IntegrationId
}

// GetIntegrationIdOk returns a tuple with the IntegrationId field value
// and a boolean to check if the value has been set.
func (o *WebhookIntegrationTypeTerraformModel) GetIntegrationIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IntegrationId, true
}

// SetIntegrationId sets field value
func (o *WebhookIntegrationTypeTerraformModel) SetIntegrationId(v string) {
	o.IntegrationId = v
}

// GetActionName returns the ActionName field value
func (o *WebhookIntegrationTypeTerraformModel) GetActionName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ActionName
}

// GetActionNameOk returns a tuple with the ActionName field value
// and a boolean to check if the value has been set.
func (o *WebhookIntegrationTypeTerraformModel) GetActionNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ActionName, true
}

// SetActionName sets field value
func (o *WebhookIntegrationTypeTerraformModel) SetActionName(v string) {
	o.ActionName = v
}

func (o WebhookIntegrationTypeTerraformModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookIntegrationTypeTerraformModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["integration_id"] = o.IntegrationId
	toSerialize["action_name"] = o.ActionName
	return toSerialize, nil
}

func (o *WebhookIntegrationTypeTerraformModel) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"integration_id",
		"action_name",
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

	varWebhookIntegrationTypeTerraformModel := _WebhookIntegrationTypeTerraformModel{}

	err = json.Unmarshal(bytes, &varWebhookIntegrationTypeTerraformModel)

	if err != nil {
		return err
	}

	*o = WebhookIntegrationTypeTerraformModel(varWebhookIntegrationTypeTerraformModel)

	return err
}

type NullableWebhookIntegrationTypeTerraformModel struct {
	value *WebhookIntegrationTypeTerraformModel
	isSet bool
}

func (v NullableWebhookIntegrationTypeTerraformModel) Get() *WebhookIntegrationTypeTerraformModel {
	return v.value
}

func (v *NullableWebhookIntegrationTypeTerraformModel) Set(val *WebhookIntegrationTypeTerraformModel) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookIntegrationTypeTerraformModel) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookIntegrationTypeTerraformModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookIntegrationTypeTerraformModel(val *WebhookIntegrationTypeTerraformModel) *NullableWebhookIntegrationTypeTerraformModel {
	return &NullableWebhookIntegrationTypeTerraformModel{value: val, isSet: true}
}

func (v NullableWebhookIntegrationTypeTerraformModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookIntegrationTypeTerraformModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
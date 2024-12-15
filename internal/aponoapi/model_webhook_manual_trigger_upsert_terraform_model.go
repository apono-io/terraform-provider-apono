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

// checks if the WebhookManualTriggerUpsertTerraformModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookManualTriggerUpsertTerraformModel{}

// WebhookManualTriggerUpsertTerraformModel struct for WebhookManualTriggerUpsertTerraformModel
type WebhookManualTriggerUpsertTerraformModel struct {
	Name                         string                                                         `json:"name"`
	Active                       bool                                                           `json:"active"`
	Type                         WebhookTypeTerraformModel                                      `json:"type"`
	BodyTemplate                 NullableString                                                 `json:"body_template,omitempty"`
	ResponseValidators           []WebhookResponseValidatorTerraformModel                       `json:"response_validators,omitempty"`
	TimeoutInSec                 NullableInt32                                                  `json:"timeout_in_sec,omitempty"`
	AuthenticationConfig         NullableWebhookManualTriggerTerraformModelAuthenticationConfig `json:"authentication_config,omitempty"`
	CustomValidationErrorMessage NullableString                                                 `json:"custom_validation_error_message,omitempty"`
}

type _WebhookManualTriggerUpsertTerraformModel WebhookManualTriggerUpsertTerraformModel

// NewWebhookManualTriggerUpsertTerraformModel instantiates a new WebhookManualTriggerUpsertTerraformModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookManualTriggerUpsertTerraformModel(name string, active bool, type_ WebhookTypeTerraformModel) *WebhookManualTriggerUpsertTerraformModel {
	this := WebhookManualTriggerUpsertTerraformModel{}
	this.Name = name
	this.Active = active
	this.Type = type_
	return &this
}

// NewWebhookManualTriggerUpsertTerraformModelWithDefaults instantiates a new WebhookManualTriggerUpsertTerraformModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookManualTriggerUpsertTerraformModelWithDefaults() *WebhookManualTriggerUpsertTerraformModel {
	this := WebhookManualTriggerUpsertTerraformModel{}
	return &this
}

// GetName returns the Name field value
func (o *WebhookManualTriggerUpsertTerraformModel) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WebhookManualTriggerUpsertTerraformModel) SetName(v string) {
	o.Name = v
}

// GetActive returns the Active field value
func (o *WebhookManualTriggerUpsertTerraformModel) GetActive() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Active
}

// GetActiveOk returns a tuple with the Active field value
// and a boolean to check if the value has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) GetActiveOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Active, true
}

// SetActive sets field value
func (o *WebhookManualTriggerUpsertTerraformModel) SetActive(v bool) {
	o.Active = v
}

// GetType returns the Type field value
func (o *WebhookManualTriggerUpsertTerraformModel) GetType() WebhookTypeTerraformModel {
	if o == nil {
		var ret WebhookTypeTerraformModel
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) GetTypeOk() (*WebhookTypeTerraformModel, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WebhookManualTriggerUpsertTerraformModel) SetType(v WebhookTypeTerraformModel) {
	o.Type = v
}

// GetBodyTemplate returns the BodyTemplate field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookManualTriggerUpsertTerraformModel) GetBodyTemplate() string {
	if o == nil || IsNil(o.BodyTemplate.Get()) {
		var ret string
		return ret
	}
	return *o.BodyTemplate.Get()
}

// GetBodyTemplateOk returns a tuple with the BodyTemplate field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookManualTriggerUpsertTerraformModel) GetBodyTemplateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.BodyTemplate.Get(), o.BodyTemplate.IsSet()
}

// HasBodyTemplate returns a boolean if a field has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) HasBodyTemplate() bool {
	if o != nil && o.BodyTemplate.IsSet() {
		return true
	}

	return false
}

// SetBodyTemplate gets a reference to the given NullableString and assigns it to the BodyTemplate field.
func (o *WebhookManualTriggerUpsertTerraformModel) SetBodyTemplate(v string) {
	o.BodyTemplate.Set(&v)
}

// SetBodyTemplateNil sets the value for BodyTemplate to be an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) SetBodyTemplateNil() {
	o.BodyTemplate.Set(nil)
}

// UnsetBodyTemplate ensures that no value is present for BodyTemplate, not even an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) UnsetBodyTemplate() {
	o.BodyTemplate.Unset()
}

// GetResponseValidators returns the ResponseValidators field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookManualTriggerUpsertTerraformModel) GetResponseValidators() []WebhookResponseValidatorTerraformModel {
	if o == nil {
		var ret []WebhookResponseValidatorTerraformModel
		return ret
	}
	return o.ResponseValidators
}

// GetResponseValidatorsOk returns a tuple with the ResponseValidators field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookManualTriggerUpsertTerraformModel) GetResponseValidatorsOk() ([]WebhookResponseValidatorTerraformModel, bool) {
	if o == nil || IsNil(o.ResponseValidators) {
		return nil, false
	}
	return o.ResponseValidators, true
}

// HasResponseValidators returns a boolean if a field has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) HasResponseValidators() bool {
	if o != nil && IsNil(o.ResponseValidators) {
		return true
	}

	return false
}

// SetResponseValidators gets a reference to the given []WebhookResponseValidatorTerraformModel and assigns it to the ResponseValidators field.
func (o *WebhookManualTriggerUpsertTerraformModel) SetResponseValidators(v []WebhookResponseValidatorTerraformModel) {
	o.ResponseValidators = v
}

// GetTimeoutInSec returns the TimeoutInSec field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookManualTriggerUpsertTerraformModel) GetTimeoutInSec() int32 {
	if o == nil || IsNil(o.TimeoutInSec.Get()) {
		var ret int32
		return ret
	}
	return *o.TimeoutInSec.Get()
}

// GetTimeoutInSecOk returns a tuple with the TimeoutInSec field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookManualTriggerUpsertTerraformModel) GetTimeoutInSecOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.TimeoutInSec.Get(), o.TimeoutInSec.IsSet()
}

// HasTimeoutInSec returns a boolean if a field has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) HasTimeoutInSec() bool {
	if o != nil && o.TimeoutInSec.IsSet() {
		return true
	}

	return false
}

// SetTimeoutInSec gets a reference to the given NullableInt32 and assigns it to the TimeoutInSec field.
func (o *WebhookManualTriggerUpsertTerraformModel) SetTimeoutInSec(v int32) {
	o.TimeoutInSec.Set(&v)
}

// SetTimeoutInSecNil sets the value for TimeoutInSec to be an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) SetTimeoutInSecNil() {
	o.TimeoutInSec.Set(nil)
}

// UnsetTimeoutInSec ensures that no value is present for TimeoutInSec, not even an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) UnsetTimeoutInSec() {
	o.TimeoutInSec.Unset()
}

// GetAuthenticationConfig returns the AuthenticationConfig field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookManualTriggerUpsertTerraformModel) GetAuthenticationConfig() WebhookManualTriggerTerraformModelAuthenticationConfig {
	if o == nil || IsNil(o.AuthenticationConfig.Get()) {
		var ret WebhookManualTriggerTerraformModelAuthenticationConfig
		return ret
	}
	return *o.AuthenticationConfig.Get()
}

// GetAuthenticationConfigOk returns a tuple with the AuthenticationConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookManualTriggerUpsertTerraformModel) GetAuthenticationConfigOk() (*WebhookManualTriggerTerraformModelAuthenticationConfig, bool) {
	if o == nil {
		return nil, false
	}
	return o.AuthenticationConfig.Get(), o.AuthenticationConfig.IsSet()
}

// HasAuthenticationConfig returns a boolean if a field has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) HasAuthenticationConfig() bool {
	if o != nil && o.AuthenticationConfig.IsSet() {
		return true
	}

	return false
}

// SetAuthenticationConfig gets a reference to the given NullableWebhookManualTriggerTerraformModelAuthenticationConfig and assigns it to the AuthenticationConfig field.
func (o *WebhookManualTriggerUpsertTerraformModel) SetAuthenticationConfig(v WebhookManualTriggerTerraformModelAuthenticationConfig) {
	o.AuthenticationConfig.Set(&v)
}

// SetAuthenticationConfigNil sets the value for AuthenticationConfig to be an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) SetAuthenticationConfigNil() {
	o.AuthenticationConfig.Set(nil)
}

// UnsetAuthenticationConfig ensures that no value is present for AuthenticationConfig, not even an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) UnsetAuthenticationConfig() {
	o.AuthenticationConfig.Unset()
}

// GetCustomValidationErrorMessage returns the CustomValidationErrorMessage field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookManualTriggerUpsertTerraformModel) GetCustomValidationErrorMessage() string {
	if o == nil || IsNil(o.CustomValidationErrorMessage.Get()) {
		var ret string
		return ret
	}
	return *o.CustomValidationErrorMessage.Get()
}

// GetCustomValidationErrorMessageOk returns a tuple with the CustomValidationErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookManualTriggerUpsertTerraformModel) GetCustomValidationErrorMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.CustomValidationErrorMessage.Get(), o.CustomValidationErrorMessage.IsSet()
}

// HasCustomValidationErrorMessage returns a boolean if a field has been set.
func (o *WebhookManualTriggerUpsertTerraformModel) HasCustomValidationErrorMessage() bool {
	if o != nil && o.CustomValidationErrorMessage.IsSet() {
		return true
	}

	return false
}

// SetCustomValidationErrorMessage gets a reference to the given NullableString and assigns it to the CustomValidationErrorMessage field.
func (o *WebhookManualTriggerUpsertTerraformModel) SetCustomValidationErrorMessage(v string) {
	o.CustomValidationErrorMessage.Set(&v)
}

// SetCustomValidationErrorMessageNil sets the value for CustomValidationErrorMessage to be an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) SetCustomValidationErrorMessageNil() {
	o.CustomValidationErrorMessage.Set(nil)
}

// UnsetCustomValidationErrorMessage ensures that no value is present for CustomValidationErrorMessage, not even an explicit nil
func (o *WebhookManualTriggerUpsertTerraformModel) UnsetCustomValidationErrorMessage() {
	o.CustomValidationErrorMessage.Unset()
}

func (o WebhookManualTriggerUpsertTerraformModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookManualTriggerUpsertTerraformModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["active"] = o.Active
	toSerialize["type"] = o.Type
	if o.BodyTemplate.IsSet() {
		toSerialize["body_template"] = o.BodyTemplate.Get()
	}
	if o.ResponseValidators != nil {
		toSerialize["response_validators"] = o.ResponseValidators
	}
	if o.TimeoutInSec.IsSet() {
		toSerialize["timeout_in_sec"] = o.TimeoutInSec.Get()
	}
	if o.AuthenticationConfig.IsSet() {
		toSerialize["authentication_config"] = o.AuthenticationConfig.Get()
	}
	if o.CustomValidationErrorMessage.IsSet() {
		toSerialize["custom_validation_error_message"] = o.CustomValidationErrorMessage.Get()
	}
	return toSerialize, nil
}

func (o *WebhookManualTriggerUpsertTerraformModel) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"active",
		"type",
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

	varWebhookManualTriggerUpsertTerraformModel := _WebhookManualTriggerUpsertTerraformModel{}

	err = json.Unmarshal(bytes, &varWebhookManualTriggerUpsertTerraformModel)

	if err != nil {
		return err
	}

	*o = WebhookManualTriggerUpsertTerraformModel(varWebhookManualTriggerUpsertTerraformModel)

	return err
}

type NullableWebhookManualTriggerUpsertTerraformModel struct {
	value *WebhookManualTriggerUpsertTerraformModel
	isSet bool
}

func (v NullableWebhookManualTriggerUpsertTerraformModel) Get() *WebhookManualTriggerUpsertTerraformModel {
	return v.value
}

func (v *NullableWebhookManualTriggerUpsertTerraformModel) Set(val *WebhookManualTriggerUpsertTerraformModel) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookManualTriggerUpsertTerraformModel) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookManualTriggerUpsertTerraformModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookManualTriggerUpsertTerraformModel(val *WebhookManualTriggerUpsertTerraformModel) *NullableWebhookManualTriggerUpsertTerraformModel {
	return &NullableWebhookManualTriggerUpsertTerraformModel{value: val, isSet: true}
}

func (v NullableWebhookManualTriggerUpsertTerraformModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookManualTriggerUpsertTerraformModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

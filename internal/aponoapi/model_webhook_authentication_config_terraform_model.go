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

// checks if the WebhookAuthenticationConfigTerraformModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookAuthenticationConfigTerraformModel{}

// WebhookAuthenticationConfigTerraformModel struct for WebhookAuthenticationConfigTerraformModel
type WebhookAuthenticationConfigTerraformModel struct {
	Type  string                                                 `json:"type"`
	Oauth NullableWebhookAuthenticationConfigTerraformModelOauth `json:"oauth,omitempty"`
}

type _WebhookAuthenticationConfigTerraformModel WebhookAuthenticationConfigTerraformModel

// NewWebhookAuthenticationConfigTerraformModel instantiates a new WebhookAuthenticationConfigTerraformModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookAuthenticationConfigTerraformModel(type_ string) *WebhookAuthenticationConfigTerraformModel {
	this := WebhookAuthenticationConfigTerraformModel{}
	this.Type = type_
	return &this
}

// NewWebhookAuthenticationConfigTerraformModelWithDefaults instantiates a new WebhookAuthenticationConfigTerraformModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookAuthenticationConfigTerraformModelWithDefaults() *WebhookAuthenticationConfigTerraformModel {
	this := WebhookAuthenticationConfigTerraformModel{}
	return &this
}

// GetType returns the Type field value
func (o *WebhookAuthenticationConfigTerraformModel) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WebhookAuthenticationConfigTerraformModel) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WebhookAuthenticationConfigTerraformModel) SetType(v string) {
	o.Type = v
}

// GetOauth returns the Oauth field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *WebhookAuthenticationConfigTerraformModel) GetOauth() WebhookAuthenticationConfigTerraformModelOauth {
	if o == nil || IsNil(o.Oauth.Get()) {
		var ret WebhookAuthenticationConfigTerraformModelOauth
		return ret
	}
	return *o.Oauth.Get()
}

// GetOauthOk returns a tuple with the Oauth field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *WebhookAuthenticationConfigTerraformModel) GetOauthOk() (*WebhookAuthenticationConfigTerraformModelOauth, bool) {
	if o == nil {
		return nil, false
	}
	return o.Oauth.Get(), o.Oauth.IsSet()
}

// HasOauth returns a boolean if a field has been set.
func (o *WebhookAuthenticationConfigTerraformModel) HasOauth() bool {
	if o != nil && o.Oauth.IsSet() {
		return true
	}

	return false
}

// SetOauth gets a reference to the given NullableWebhookAuthenticationConfigTerraformModelOauth and assigns it to the Oauth field.
func (o *WebhookAuthenticationConfigTerraformModel) SetOauth(v WebhookAuthenticationConfigTerraformModelOauth) {
	o.Oauth.Set(&v)
}

// SetOauthNil sets the value for Oauth to be an explicit nil
func (o *WebhookAuthenticationConfigTerraformModel) SetOauthNil() {
	o.Oauth.Set(nil)
}

// UnsetOauth ensures that no value is present for Oauth, not even an explicit nil
func (o *WebhookAuthenticationConfigTerraformModel) UnsetOauth() {
	o.Oauth.Unset()
}

func (o WebhookAuthenticationConfigTerraformModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookAuthenticationConfigTerraformModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["type"] = o.Type
	if o.Oauth.IsSet() {
		toSerialize["oauth"] = o.Oauth.Get()
	}
	return toSerialize, nil
}

func (o *WebhookAuthenticationConfigTerraformModel) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
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

	varWebhookAuthenticationConfigTerraformModel := _WebhookAuthenticationConfigTerraformModel{}

	err = json.Unmarshal(bytes, &varWebhookAuthenticationConfigTerraformModel)

	if err != nil {
		return err
	}

	*o = WebhookAuthenticationConfigTerraformModel(varWebhookAuthenticationConfigTerraformModel)

	return err
}

type NullableWebhookAuthenticationConfigTerraformModel struct {
	value *WebhookAuthenticationConfigTerraformModel
	isSet bool
}

func (v NullableWebhookAuthenticationConfigTerraformModel) Get() *WebhookAuthenticationConfigTerraformModel {
	return v.value
}

func (v *NullableWebhookAuthenticationConfigTerraformModel) Set(val *WebhookAuthenticationConfigTerraformModel) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookAuthenticationConfigTerraformModel) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookAuthenticationConfigTerraformModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookAuthenticationConfigTerraformModel(val *WebhookAuthenticationConfigTerraformModel) *NullableWebhookAuthenticationConfigTerraformModel {
	return &NullableWebhookAuthenticationConfigTerraformModel{value: val, isSet: true}
}

func (v NullableWebhookAuthenticationConfigTerraformModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookAuthenticationConfigTerraformModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

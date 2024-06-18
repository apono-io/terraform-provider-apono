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

// checks if the UpsertIntegrationTerraform type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpsertIntegrationTerraform{}

// UpsertIntegrationTerraform struct for UpsertIntegrationTerraform
type UpsertIntegrationTerraform struct {
	Name                   string                          `json:"name"`
	Type                   string                          `json:"type"`
	ProvisionerId          NullableString                  `json:"provisioner_id,omitempty"`
	Params                 map[string]interface{}          `json:"params"`
	SecretConfig           map[string]interface{}          `json:"secret_config,omitempty"`
	ConnectedResourceTypes []string                        `json:"connected_resource_types,omitempty"`
	CustomAccessDetails    NullableString                  `json:"custom_access_details,omitempty"`
	IntegrationOwners      []IntegrationOwnerTerraform     `json:"integration_owners,omitempty"`
	ResourceOwnersMappings []ResourceOwnerMappingTerraform `json:"resource_owners_mappings,omitempty"`
}

type _UpsertIntegrationTerraform UpsertIntegrationTerraform

// NewUpsertIntegrationTerraform instantiates a new UpsertIntegrationTerraform object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpsertIntegrationTerraform(name string, type_ string, params map[string]interface{}) *UpsertIntegrationTerraform {
	this := UpsertIntegrationTerraform{}
	this.Name = name
	this.Type = type_
	this.Params = params
	return &this
}

// NewUpsertIntegrationTerraformWithDefaults instantiates a new UpsertIntegrationTerraform object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpsertIntegrationTerraformWithDefaults() *UpsertIntegrationTerraform {
	this := UpsertIntegrationTerraform{}
	return &this
}

// GetName returns the Name field value
func (o *UpsertIntegrationTerraform) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *UpsertIntegrationTerraform) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *UpsertIntegrationTerraform) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *UpsertIntegrationTerraform) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *UpsertIntegrationTerraform) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *UpsertIntegrationTerraform) SetType(v string) {
	o.Type = v
}

// GetProvisionerId returns the ProvisionerId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetProvisionerId() string {
	if o == nil || IsNil(o.ProvisionerId.Get()) {
		var ret string
		return ret
	}
	return *o.ProvisionerId.Get()
}

// GetProvisionerIdOk returns a tuple with the ProvisionerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetProvisionerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ProvisionerId.Get(), o.ProvisionerId.IsSet()
}

// HasProvisionerId returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasProvisionerId() bool {
	if o != nil && o.ProvisionerId.IsSet() {
		return true
	}

	return false
}

// SetProvisionerId gets a reference to the given NullableString and assigns it to the ProvisionerId field.
func (o *UpsertIntegrationTerraform) SetProvisionerId(v string) {
	o.ProvisionerId.Set(&v)
}

// SetProvisionerIdNil sets the value for ProvisionerId to be an explicit nil
func (o *UpsertIntegrationTerraform) SetProvisionerIdNil() {
	o.ProvisionerId.Set(nil)
}

// UnsetProvisionerId ensures that no value is present for ProvisionerId, not even an explicit nil
func (o *UpsertIntegrationTerraform) UnsetProvisionerId() {
	o.ProvisionerId.Unset()
}

// GetParams returns the Params field value
func (o *UpsertIntegrationTerraform) GetParams() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Params
}

// GetParamsOk returns a tuple with the Params field value
// and a boolean to check if the value has been set.
func (o *UpsertIntegrationTerraform) GetParamsOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.Params, true
}

// SetParams sets field value
func (o *UpsertIntegrationTerraform) SetParams(v map[string]interface{}) {
	o.Params = v
}

// GetSecretConfig returns the SecretConfig field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetSecretConfig() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.SecretConfig
}

// GetSecretConfigOk returns a tuple with the SecretConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetSecretConfigOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.SecretConfig) {
		return map[string]interface{}{}, false
	}
	return o.SecretConfig, true
}

// HasSecretConfig returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasSecretConfig() bool {
	if o != nil && IsNil(o.SecretConfig) {
		return true
	}

	return false
}

// SetSecretConfig gets a reference to the given map[string]interface{} and assigns it to the SecretConfig field.
func (o *UpsertIntegrationTerraform) SetSecretConfig(v map[string]interface{}) {
	o.SecretConfig = v
}

// GetConnectedResourceTypes returns the ConnectedResourceTypes field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetConnectedResourceTypes() []string {
	if o == nil {
		var ret []string
		return ret
	}
	return o.ConnectedResourceTypes
}

// GetConnectedResourceTypesOk returns a tuple with the ConnectedResourceTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetConnectedResourceTypesOk() ([]string, bool) {
	if o == nil || IsNil(o.ConnectedResourceTypes) {
		return nil, false
	}
	return o.ConnectedResourceTypes, true
}

// HasConnectedResourceTypes returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasConnectedResourceTypes() bool {
	if o != nil && IsNil(o.ConnectedResourceTypes) {
		return true
	}

	return false
}

// SetConnectedResourceTypes gets a reference to the given []string and assigns it to the ConnectedResourceTypes field.
func (o *UpsertIntegrationTerraform) SetConnectedResourceTypes(v []string) {
	o.ConnectedResourceTypes = v
}

// GetCustomAccessDetails returns the CustomAccessDetails field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetCustomAccessDetails() string {
	if o == nil || IsNil(o.CustomAccessDetails.Get()) {
		var ret string
		return ret
	}
	return *o.CustomAccessDetails.Get()
}

// GetCustomAccessDetailsOk returns a tuple with the CustomAccessDetails field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetCustomAccessDetailsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.CustomAccessDetails.Get(), o.CustomAccessDetails.IsSet()
}

// HasCustomAccessDetails returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasCustomAccessDetails() bool {
	if o != nil && o.CustomAccessDetails.IsSet() {
		return true
	}

	return false
}

// SetCustomAccessDetails gets a reference to the given NullableString and assigns it to the CustomAccessDetails field.
func (o *UpsertIntegrationTerraform) SetCustomAccessDetails(v string) {
	o.CustomAccessDetails.Set(&v)
}

// SetCustomAccessDetailsNil sets the value for CustomAccessDetails to be an explicit nil
func (o *UpsertIntegrationTerraform) SetCustomAccessDetailsNil() {
	o.CustomAccessDetails.Set(nil)
}

// UnsetCustomAccessDetails ensures that no value is present for CustomAccessDetails, not even an explicit nil
func (o *UpsertIntegrationTerraform) UnsetCustomAccessDetails() {
	o.CustomAccessDetails.Unset()
}

// GetIntegrationOwners returns the IntegrationOwners field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetIntegrationOwners() []IntegrationOwnerTerraform {
	if o == nil {
		var ret []IntegrationOwnerTerraform
		return ret
	}
	return o.IntegrationOwners
}

// GetIntegrationOwnersOk returns a tuple with the IntegrationOwners field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetIntegrationOwnersOk() ([]IntegrationOwnerTerraform, bool) {
	if o == nil || IsNil(o.IntegrationOwners) {
		return nil, false
	}
	return o.IntegrationOwners, true
}

// HasIntegrationOwners returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasIntegrationOwners() bool {
	if o != nil && IsNil(o.IntegrationOwners) {
		return true
	}

	return false
}

// SetIntegrationOwners gets a reference to the given []IntegrationOwnerTerraform and assigns it to the IntegrationOwners field.
func (o *UpsertIntegrationTerraform) SetIntegrationOwners(v []IntegrationOwnerTerraform) {
	o.IntegrationOwners = v
}

// GetResourceOwnersMappings returns the ResourceOwnersMappings field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertIntegrationTerraform) GetResourceOwnersMappings() []ResourceOwnerMappingTerraform {
	if o == nil {
		var ret []ResourceOwnerMappingTerraform
		return ret
	}
	return o.ResourceOwnersMappings
}

// GetResourceOwnersMappingsOk returns a tuple with the ResourceOwnersMappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertIntegrationTerraform) GetResourceOwnersMappingsOk() ([]ResourceOwnerMappingTerraform, bool) {
	if o == nil || IsNil(o.ResourceOwnersMappings) {
		return nil, false
	}
	return o.ResourceOwnersMappings, true
}

// HasResourceOwnersMappings returns a boolean if a field has been set.
func (o *UpsertIntegrationTerraform) HasResourceOwnersMappings() bool {
	if o != nil && IsNil(o.ResourceOwnersMappings) {
		return true
	}

	return false
}

// SetResourceOwnersMappings gets a reference to the given []ResourceOwnerMappingTerraform and assigns it to the ResourceOwnersMappings field.
func (o *UpsertIntegrationTerraform) SetResourceOwnersMappings(v []ResourceOwnerMappingTerraform) {
	o.ResourceOwnersMappings = v
}

func (o UpsertIntegrationTerraform) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpsertIntegrationTerraform) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	if o.ProvisionerId.IsSet() {
		toSerialize["provisioner_id"] = o.ProvisionerId.Get()
	}
	toSerialize["params"] = o.Params
	if o.SecretConfig != nil {
		toSerialize["secret_config"] = o.SecretConfig
	}
	if o.ConnectedResourceTypes != nil {
		toSerialize["connected_resource_types"] = o.ConnectedResourceTypes
	}
	if o.CustomAccessDetails.IsSet() {
		toSerialize["custom_access_details"] = o.CustomAccessDetails.Get()
	}
	if o.IntegrationOwners != nil {
		toSerialize["integration_owners"] = o.IntegrationOwners
	}
	if o.ResourceOwnersMappings != nil {
		toSerialize["resource_owners_mappings"] = o.ResourceOwnersMappings
	}
	return toSerialize, nil
}

func (o *UpsertIntegrationTerraform) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"type",
		"params",
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

	varUpsertIntegrationTerraform := _UpsertIntegrationTerraform{}

	err = json.Unmarshal(bytes, &varUpsertIntegrationTerraform)

	if err != nil {
		return err
	}

	*o = UpsertIntegrationTerraform(varUpsertIntegrationTerraform)

	return err
}

type NullableUpsertIntegrationTerraform struct {
	value *UpsertIntegrationTerraform
	isSet bool
}

func (v NullableUpsertIntegrationTerraform) Get() *UpsertIntegrationTerraform {
	return v.value
}

func (v *NullableUpsertIntegrationTerraform) Set(val *UpsertIntegrationTerraform) {
	v.value = val
	v.isSet = true
}

func (v NullableUpsertIntegrationTerraform) IsSet() bool {
	return v.isSet
}

func (v *NullableUpsertIntegrationTerraform) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpsertIntegrationTerraform(val *UpsertIntegrationTerraform) *NullableUpsertIntegrationTerraform {
	return &NullableUpsertIntegrationTerraform{value: val, isSet: true}
}

func (v NullableUpsertIntegrationTerraform) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpsertIntegrationTerraform) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

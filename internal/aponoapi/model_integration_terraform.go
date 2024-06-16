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

// checks if the IntegrationTerraform type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IntegrationTerraform{}

// IntegrationTerraform struct for IntegrationTerraform
type IntegrationTerraform struct {
	Id                     string                          `json:"id"`
	Name                   string                          `json:"name"`
	Type                   string                          `json:"type"`
	Status                 IntegrationStatus               `json:"status"`
	ProvisionerId          NullableString                  `json:"provisioner_id,omitempty"`
	LastSyncTime           NullableFloat64                 `json:"last_sync_time,omitempty"`
	Params                 map[string]interface{}          `json:"params"`
	SecretConfig           map[string]interface{}          `json:"secret_config,omitempty"`
	ConnectedResourceTypes []string                        `json:"connected_resource_types"`
	CustomAccessDetails    NullableString                  `json:"custom_access_details,omitempty"`
	IntegrationOwners      IntegrationOwnersTerraform      `json:"integration_owners"`
	ResourceOwnersMappings []ResourceOwnerMappingTerraform `json:"resource_owners_mappings"`
}

type _IntegrationTerraform IntegrationTerraform

// NewIntegrationTerraform instantiates a new IntegrationTerraform object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIntegrationTerraform(id string, name string, type_ string, status IntegrationStatus, params map[string]interface{}, connectedResourceTypes []string, integrationOwners IntegrationOwnersTerraform, resourceOwnersMappings []ResourceOwnerMappingTerraform) *IntegrationTerraform {
	this := IntegrationTerraform{}
	this.Id = id
	this.Name = name
	this.Type = type_
	this.Status = status
	this.Params = params
	this.ConnectedResourceTypes = connectedResourceTypes
	this.IntegrationOwners = integrationOwners
	this.ResourceOwnersMappings = resourceOwnersMappings
	return &this
}

// NewIntegrationTerraformWithDefaults instantiates a new IntegrationTerraform object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIntegrationTerraformWithDefaults() *IntegrationTerraform {
	this := IntegrationTerraform{}
	return &this
}

// GetId returns the Id field value
func (o *IntegrationTerraform) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *IntegrationTerraform) SetId(v string) {
	o.Id = v
}

// GetName returns the Name field value
func (o *IntegrationTerraform) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *IntegrationTerraform) SetName(v string) {
	o.Name = v
}

// GetType returns the Type field value
func (o *IntegrationTerraform) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *IntegrationTerraform) SetType(v string) {
	o.Type = v
}

// GetStatus returns the Status field value
func (o *IntegrationTerraform) GetStatus() IntegrationStatus {
	if o == nil {
		var ret IntegrationStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetStatusOk() (*IntegrationStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *IntegrationTerraform) SetStatus(v IntegrationStatus) {
	o.Status = v
}

// GetProvisionerId returns the ProvisionerId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *IntegrationTerraform) GetProvisionerId() string {
	if o == nil || IsNil(o.ProvisionerId.Get()) {
		var ret string
		return ret
	}
	return *o.ProvisionerId.Get()
}

// GetProvisionerIdOk returns a tuple with the ProvisionerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IntegrationTerraform) GetProvisionerIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ProvisionerId.Get(), o.ProvisionerId.IsSet()
}

// HasProvisionerId returns a boolean if a field has been set.
func (o *IntegrationTerraform) HasProvisionerId() bool {
	if o != nil && o.ProvisionerId.IsSet() {
		return true
	}

	return false
}

// SetProvisionerId gets a reference to the given NullableString and assigns it to the ProvisionerId field.
func (o *IntegrationTerraform) SetProvisionerId(v string) {
	o.ProvisionerId.Set(&v)
}

// SetProvisionerIdNil sets the value for ProvisionerId to be an explicit nil
func (o *IntegrationTerraform) SetProvisionerIdNil() {
	o.ProvisionerId.Set(nil)
}

// UnsetProvisionerId ensures that no value is present for ProvisionerId, not even an explicit nil
func (o *IntegrationTerraform) UnsetProvisionerId() {
	o.ProvisionerId.Unset()
}

// GetLastSyncTime returns the LastSyncTime field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *IntegrationTerraform) GetLastSyncTime() float64 {
	if o == nil || IsNil(o.LastSyncTime.Get()) {
		var ret float64
		return ret
	}
	return *o.LastSyncTime.Get()
}

// GetLastSyncTimeOk returns a tuple with the LastSyncTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IntegrationTerraform) GetLastSyncTimeOk() (*float64, bool) {
	if o == nil {
		return nil, false
	}
	return o.LastSyncTime.Get(), o.LastSyncTime.IsSet()
}

// HasLastSyncTime returns a boolean if a field has been set.
func (o *IntegrationTerraform) HasLastSyncTime() bool {
	if o != nil && o.LastSyncTime.IsSet() {
		return true
	}

	return false
}

// SetLastSyncTime gets a reference to the given NullableFloat64 and assigns it to the LastSyncTime field.
func (o *IntegrationTerraform) SetLastSyncTime(v float64) {
	o.LastSyncTime.Set(&v)
}

// SetLastSyncTimeNil sets the value for LastSyncTime to be an explicit nil
func (o *IntegrationTerraform) SetLastSyncTimeNil() {
	o.LastSyncTime.Set(nil)
}

// UnsetLastSyncTime ensures that no value is present for LastSyncTime, not even an explicit nil
func (o *IntegrationTerraform) UnsetLastSyncTime() {
	o.LastSyncTime.Unset()
}

// GetParams returns the Params field value
func (o *IntegrationTerraform) GetParams() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.Params
}

// GetParamsOk returns a tuple with the Params field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetParamsOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.Params, true
}

// SetParams sets field value
func (o *IntegrationTerraform) SetParams(v map[string]interface{}) {
	o.Params = v
}

// GetSecretConfig returns the SecretConfig field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *IntegrationTerraform) GetSecretConfig() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}
	return o.SecretConfig
}

// GetSecretConfigOk returns a tuple with the SecretConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IntegrationTerraform) GetSecretConfigOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.SecretConfig) {
		return map[string]interface{}{}, false
	}
	return o.SecretConfig, true
}

// HasSecretConfig returns a boolean if a field has been set.
func (o *IntegrationTerraform) HasSecretConfig() bool {
	if o != nil && IsNil(o.SecretConfig) {
		return true
	}

	return false
}

// SetSecretConfig gets a reference to the given map[string]interface{} and assigns it to the SecretConfig field.
func (o *IntegrationTerraform) SetSecretConfig(v map[string]interface{}) {
	o.SecretConfig = v
}

// GetConnectedResourceTypes returns the ConnectedResourceTypes field value
func (o *IntegrationTerraform) GetConnectedResourceTypes() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.ConnectedResourceTypes
}

// GetConnectedResourceTypesOk returns a tuple with the ConnectedResourceTypes field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetConnectedResourceTypesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ConnectedResourceTypes, true
}

// SetConnectedResourceTypes sets field value
func (o *IntegrationTerraform) SetConnectedResourceTypes(v []string) {
	o.ConnectedResourceTypes = v
}

// GetCustomAccessDetails returns the CustomAccessDetails field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *IntegrationTerraform) GetCustomAccessDetails() string {
	if o == nil || IsNil(o.CustomAccessDetails.Get()) {
		var ret string
		return ret
	}
	return *o.CustomAccessDetails.Get()
}

// GetCustomAccessDetailsOk returns a tuple with the CustomAccessDetails field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *IntegrationTerraform) GetCustomAccessDetailsOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.CustomAccessDetails.Get(), o.CustomAccessDetails.IsSet()
}

// HasCustomAccessDetails returns a boolean if a field has been set.
func (o *IntegrationTerraform) HasCustomAccessDetails() bool {
	if o != nil && o.CustomAccessDetails.IsSet() {
		return true
	}

	return false
}

// SetCustomAccessDetails gets a reference to the given NullableString and assigns it to the CustomAccessDetails field.
func (o *IntegrationTerraform) SetCustomAccessDetails(v string) {
	o.CustomAccessDetails.Set(&v)
}

// SetCustomAccessDetailsNil sets the value for CustomAccessDetails to be an explicit nil
func (o *IntegrationTerraform) SetCustomAccessDetailsNil() {
	o.CustomAccessDetails.Set(nil)
}

// UnsetCustomAccessDetails ensures that no value is present for CustomAccessDetails, not even an explicit nil
func (o *IntegrationTerraform) UnsetCustomAccessDetails() {
	o.CustomAccessDetails.Unset()
}

// GetIntegrationOwners returns the IntegrationOwners field value
func (o *IntegrationTerraform) GetIntegrationOwners() IntegrationOwnersTerraform {
	if o == nil {
		var ret IntegrationOwnersTerraform
		return ret
	}

	return o.IntegrationOwners
}

// GetIntegrationOwnersOk returns a tuple with the IntegrationOwners field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetIntegrationOwnersOk() (*IntegrationOwnersTerraform, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IntegrationOwners, true
}

// SetIntegrationOwners sets field value
func (o *IntegrationTerraform) SetIntegrationOwners(v IntegrationOwnersTerraform) {
	o.IntegrationOwners = v
}

// GetResourceOwnersMappings returns the ResourceOwnersMappings field value
func (o *IntegrationTerraform) GetResourceOwnersMappings() []ResourceOwnerMappingTerraform {
	if o == nil {
		var ret []ResourceOwnerMappingTerraform
		return ret
	}

	return o.ResourceOwnersMappings
}

// GetResourceOwnersMappingsOk returns a tuple with the ResourceOwnersMappings field value
// and a boolean to check if the value has been set.
func (o *IntegrationTerraform) GetResourceOwnersMappingsOk() ([]ResourceOwnerMappingTerraform, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResourceOwnersMappings, true
}

// SetResourceOwnersMappings sets field value
func (o *IntegrationTerraform) SetResourceOwnersMappings(v []ResourceOwnerMappingTerraform) {
	o.ResourceOwnersMappings = v
}

func (o IntegrationTerraform) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IntegrationTerraform) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["name"] = o.Name
	toSerialize["type"] = o.Type
	toSerialize["status"] = o.Status
	if o.ProvisionerId.IsSet() {
		toSerialize["provisioner_id"] = o.ProvisionerId.Get()
	}
	if o.LastSyncTime.IsSet() {
		toSerialize["last_sync_time"] = o.LastSyncTime.Get()
	}
	toSerialize["params"] = o.Params
	if o.SecretConfig != nil {
		toSerialize["secret_config"] = o.SecretConfig
	}
	toSerialize["connected_resource_types"] = o.ConnectedResourceTypes
	if o.CustomAccessDetails.IsSet() {
		toSerialize["custom_access_details"] = o.CustomAccessDetails.Get()
	}
	toSerialize["integration_owners"] = o.IntegrationOwners
	toSerialize["resource_owners_mappings"] = o.ResourceOwnersMappings
	return toSerialize, nil
}

func (o *IntegrationTerraform) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
		"name",
		"type",
		"status",
		"params",
		"connected_resource_types",
		"integration_owners",
		"resource_owners_mappings",
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

	varIntegrationTerraform := _IntegrationTerraform{}

	err = json.Unmarshal(bytes, &varIntegrationTerraform)

	if err != nil {
		return err
	}

	*o = IntegrationTerraform(varIntegrationTerraform)

	return err
}

type NullableIntegrationTerraform struct {
	value *IntegrationTerraform
	isSet bool
}

func (v NullableIntegrationTerraform) Get() *IntegrationTerraform {
	return v.value
}

func (v *NullableIntegrationTerraform) Set(val *IntegrationTerraform) {
	v.value = val
	v.isSet = true
}

func (v NullableIntegrationTerraform) IsSet() bool {
	return v.isSet
}

func (v *NullableIntegrationTerraform) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIntegrationTerraform(val *IntegrationTerraform) *NullableIntegrationTerraform {
	return &NullableIntegrationTerraform{value: val, isSet: true}
}

func (v NullableIntegrationTerraform) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIntegrationTerraform) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
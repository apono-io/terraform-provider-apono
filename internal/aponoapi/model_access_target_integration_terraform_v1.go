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

// checks if the AccessTargetIntegrationTerraformV1 type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AccessTargetIntegrationTerraformV1{}

// AccessTargetIntegrationTerraformV1 struct for AccessTargetIntegrationTerraformV1
type AccessTargetIntegrationTerraformV1 struct {
	IntegrationId       string           `json:"integration_id"`
	ResourceType        string           `json:"resource_type"`
	ResourceTagIncludes []TagTerraformV1 `json:"resource_tag_includes"`
	ResourceTagExcludes []TagTerraformV1 `json:"resource_tag_excludes"`
	Permissions         []string         `json:"permissions"`
}

type _AccessTargetIntegrationTerraformV1 AccessTargetIntegrationTerraformV1

// NewAccessTargetIntegrationTerraformV1 instantiates a new AccessTargetIntegrationTerraformV1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccessTargetIntegrationTerraformV1(integrationId string, resourceType string, resourceTagIncludes []TagTerraformV1, resourceTagExcludes []TagTerraformV1, permissions []string) *AccessTargetIntegrationTerraformV1 {
	this := AccessTargetIntegrationTerraformV1{}
	this.IntegrationId = integrationId
	this.ResourceType = resourceType
	this.ResourceTagIncludes = resourceTagIncludes
	this.ResourceTagExcludes = resourceTagExcludes
	this.Permissions = permissions
	return &this
}

// NewAccessTargetIntegrationTerraformV1WithDefaults instantiates a new AccessTargetIntegrationTerraformV1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccessTargetIntegrationTerraformV1WithDefaults() *AccessTargetIntegrationTerraformV1 {
	this := AccessTargetIntegrationTerraformV1{}
	return &this
}

// GetIntegrationId returns the IntegrationId field value
func (o *AccessTargetIntegrationTerraformV1) GetIntegrationId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.IntegrationId
}

// GetIntegrationIdOk returns a tuple with the IntegrationId field value
// and a boolean to check if the value has been set.
func (o *AccessTargetIntegrationTerraformV1) GetIntegrationIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.IntegrationId, true
}

// SetIntegrationId sets field value
func (o *AccessTargetIntegrationTerraformV1) SetIntegrationId(v string) {
	o.IntegrationId = v
}

// GetResourceType returns the ResourceType field value
func (o *AccessTargetIntegrationTerraformV1) GetResourceType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ResourceType
}

// GetResourceTypeOk returns a tuple with the ResourceType field value
// and a boolean to check if the value has been set.
func (o *AccessTargetIntegrationTerraformV1) GetResourceTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ResourceType, true
}

// SetResourceType sets field value
func (o *AccessTargetIntegrationTerraformV1) SetResourceType(v string) {
	o.ResourceType = v
}

// GetResourceTagIncludes returns the ResourceTagIncludes field value
func (o *AccessTargetIntegrationTerraformV1) GetResourceTagIncludes() []TagTerraformV1 {
	if o == nil {
		var ret []TagTerraformV1
		return ret
	}

	return o.ResourceTagIncludes
}

// GetResourceTagIncludesOk returns a tuple with the ResourceTagIncludes field value
// and a boolean to check if the value has been set.
func (o *AccessTargetIntegrationTerraformV1) GetResourceTagIncludesOk() ([]TagTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResourceTagIncludes, true
}

// SetResourceTagIncludes sets field value
func (o *AccessTargetIntegrationTerraformV1) SetResourceTagIncludes(v []TagTerraformV1) {
	o.ResourceTagIncludes = v
}

// GetResourceTagExcludes returns the ResourceTagExcludes field value
func (o *AccessTargetIntegrationTerraformV1) GetResourceTagExcludes() []TagTerraformV1 {
	if o == nil {
		var ret []TagTerraformV1
		return ret
	}

	return o.ResourceTagExcludes
}

// GetResourceTagExcludesOk returns a tuple with the ResourceTagExcludes field value
// and a boolean to check if the value has been set.
func (o *AccessTargetIntegrationTerraformV1) GetResourceTagExcludesOk() ([]TagTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResourceTagExcludes, true
}

// SetResourceTagExcludes sets field value
func (o *AccessTargetIntegrationTerraformV1) SetResourceTagExcludes(v []TagTerraformV1) {
	o.ResourceTagExcludes = v
}

// GetPermissions returns the Permissions field value
func (o *AccessTargetIntegrationTerraformV1) GetPermissions() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Permissions
}

// GetPermissionsOk returns a tuple with the Permissions field value
// and a boolean to check if the value has been set.
func (o *AccessTargetIntegrationTerraformV1) GetPermissionsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Permissions, true
}

// SetPermissions sets field value
func (o *AccessTargetIntegrationTerraformV1) SetPermissions(v []string) {
	o.Permissions = v
}

func (o AccessTargetIntegrationTerraformV1) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AccessTargetIntegrationTerraformV1) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["integration_id"] = o.IntegrationId
	toSerialize["resource_type"] = o.ResourceType
	toSerialize["resource_tag_includes"] = o.ResourceTagIncludes
	toSerialize["resource_tag_excludes"] = o.ResourceTagExcludes
	toSerialize["permissions"] = o.Permissions
	return toSerialize, nil
}

func (o *AccessTargetIntegrationTerraformV1) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"integration_id",
		"resource_type",
		"resource_tag_includes",
		"resource_tag_excludes",
		"permissions",
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

	varAccessTargetIntegrationTerraformV1 := _AccessTargetIntegrationTerraformV1{}

	err = json.Unmarshal(bytes, &varAccessTargetIntegrationTerraformV1)

	if err != nil {
		return err
	}

	*o = AccessTargetIntegrationTerraformV1(varAccessTargetIntegrationTerraformV1)

	return err
}

type NullableAccessTargetIntegrationTerraformV1 struct {
	value *AccessTargetIntegrationTerraformV1
	isSet bool
}

func (v NullableAccessTargetIntegrationTerraformV1) Get() *AccessTargetIntegrationTerraformV1 {
	return v.value
}

func (v *NullableAccessTargetIntegrationTerraformV1) Set(val *AccessTargetIntegrationTerraformV1) {
	v.value = val
	v.isSet = true
}

func (v NullableAccessTargetIntegrationTerraformV1) IsSet() bool {
	return v.isSet
}

func (v *NullableAccessTargetIntegrationTerraformV1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccessTargetIntegrationTerraformV1(val *AccessTargetIntegrationTerraformV1) *NullableAccessTargetIntegrationTerraformV1 {
	return &NullableAccessTargetIntegrationTerraformV1{value: val, isSet: true}
}

func (v NullableAccessTargetIntegrationTerraformV1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccessTargetIntegrationTerraformV1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

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

// checks if the GranteeTerraformV1 type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GranteeTerraformV1{}

// GranteeTerraformV1 struct for GranteeTerraformV1
type GranteeTerraformV1 struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type _GranteeTerraformV1 GranteeTerraformV1

// NewGranteeTerraformV1 instantiates a new GranteeTerraformV1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGranteeTerraformV1(id string, type_ string) *GranteeTerraformV1 {
	this := GranteeTerraformV1{}
	this.Id = id
	this.Type = type_
	return &this
}

// NewGranteeTerraformV1WithDefaults instantiates a new GranteeTerraformV1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGranteeTerraformV1WithDefaults() *GranteeTerraformV1 {
	this := GranteeTerraformV1{}
	return &this
}

// GetId returns the Id field value
func (o *GranteeTerraformV1) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GranteeTerraformV1) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GranteeTerraformV1) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *GranteeTerraformV1) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *GranteeTerraformV1) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *GranteeTerraformV1) SetType(v string) {
	o.Type = v
}

func (o GranteeTerraformV1) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GranteeTerraformV1) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	return toSerialize, nil
}

func (o *GranteeTerraformV1) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
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

	varGranteeTerraformV1 := _GranteeTerraformV1{}

	err = json.Unmarshal(bytes, &varGranteeTerraformV1)

	if err != nil {
		return err
	}

	*o = GranteeTerraformV1(varGranteeTerraformV1)

	return err
}

type NullableGranteeTerraformV1 struct {
	value *GranteeTerraformV1
	isSet bool
}

func (v NullableGranteeTerraformV1) Get() *GranteeTerraformV1 {
	return v.value
}

func (v *NullableGranteeTerraformV1) Set(val *GranteeTerraformV1) {
	v.value = val
	v.isSet = true
}

func (v NullableGranteeTerraformV1) IsSet() bool {
	return v.isSet
}

func (v *NullableGranteeTerraformV1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGranteeTerraformV1(val *GranteeTerraformV1) *NullableGranteeTerraformV1 {
	return &NullableGranteeTerraformV1{value: val, isSet: true}
}

func (v NullableGranteeTerraformV1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGranteeTerraformV1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

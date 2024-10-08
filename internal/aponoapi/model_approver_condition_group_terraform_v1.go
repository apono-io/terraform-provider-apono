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

// checks if the ApproverConditionGroupTerraformV1 type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApproverConditionGroupTerraformV1{}

// ApproverConditionGroupTerraformV1 struct for ApproverConditionGroupTerraformV1
type ApproverConditionGroupTerraformV1 struct {
	ConditionsLogicalOperator ApproverConditionGroupOperatorTerraformV1 `json:"conditions_logical_operator"`
	Conditions                []AttributeFilterTerraformV1              `json:"conditions"`
}

type _ApproverConditionGroupTerraformV1 ApproverConditionGroupTerraformV1

// NewApproverConditionGroupTerraformV1 instantiates a new ApproverConditionGroupTerraformV1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApproverConditionGroupTerraformV1(conditionsLogicalOperator ApproverConditionGroupOperatorTerraformV1, conditions []AttributeFilterTerraformV1) *ApproverConditionGroupTerraformV1 {
	this := ApproverConditionGroupTerraformV1{}
	this.ConditionsLogicalOperator = conditionsLogicalOperator
	this.Conditions = conditions
	return &this
}

// NewApproverConditionGroupTerraformV1WithDefaults instantiates a new ApproverConditionGroupTerraformV1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApproverConditionGroupTerraformV1WithDefaults() *ApproverConditionGroupTerraformV1 {
	this := ApproverConditionGroupTerraformV1{}
	return &this
}

// GetConditionsLogicalOperator returns the ConditionsLogicalOperator field value
func (o *ApproverConditionGroupTerraformV1) GetConditionsLogicalOperator() ApproverConditionGroupOperatorTerraformV1 {
	if o == nil {
		var ret ApproverConditionGroupOperatorTerraformV1
		return ret
	}

	return o.ConditionsLogicalOperator
}

// GetConditionsLogicalOperatorOk returns a tuple with the ConditionsLogicalOperator field value
// and a boolean to check if the value has been set.
func (o *ApproverConditionGroupTerraformV1) GetConditionsLogicalOperatorOk() (*ApproverConditionGroupOperatorTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConditionsLogicalOperator, true
}

// SetConditionsLogicalOperator sets field value
func (o *ApproverConditionGroupTerraformV1) SetConditionsLogicalOperator(v ApproverConditionGroupOperatorTerraformV1) {
	o.ConditionsLogicalOperator = v
}

// GetConditions returns the Conditions field value
func (o *ApproverConditionGroupTerraformV1) GetConditions() []AttributeFilterTerraformV1 {
	if o == nil {
		var ret []AttributeFilterTerraformV1
		return ret
	}

	return o.Conditions
}

// GetConditionsOk returns a tuple with the Conditions field value
// and a boolean to check if the value has been set.
func (o *ApproverConditionGroupTerraformV1) GetConditionsOk() ([]AttributeFilterTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.Conditions, true
}

// SetConditions sets field value
func (o *ApproverConditionGroupTerraformV1) SetConditions(v []AttributeFilterTerraformV1) {
	o.Conditions = v
}

func (o ApproverConditionGroupTerraformV1) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApproverConditionGroupTerraformV1) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["conditions_logical_operator"] = o.ConditionsLogicalOperator
	toSerialize["conditions"] = o.Conditions
	return toSerialize, nil
}

func (o *ApproverConditionGroupTerraformV1) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"conditions_logical_operator",
		"conditions",
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

	varApproverConditionGroupTerraformV1 := _ApproverConditionGroupTerraformV1{}

	err = json.Unmarshal(bytes, &varApproverConditionGroupTerraformV1)

	if err != nil {
		return err
	}

	*o = ApproverConditionGroupTerraformV1(varApproverConditionGroupTerraformV1)

	return err
}

type NullableApproverConditionGroupTerraformV1 struct {
	value *ApproverConditionGroupTerraformV1
	isSet bool
}

func (v NullableApproverConditionGroupTerraformV1) Get() *ApproverConditionGroupTerraformV1 {
	return v.value
}

func (v *NullableApproverConditionGroupTerraformV1) Set(val *ApproverConditionGroupTerraformV1) {
	v.value = val
	v.isSet = true
}

func (v NullableApproverConditionGroupTerraformV1) IsSet() bool {
	return v.isSet
}

func (v *NullableApproverConditionGroupTerraformV1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApproverConditionGroupTerraformV1(val *ApproverConditionGroupTerraformV1) *NullableApproverConditionGroupTerraformV1 {
	return &NullableApproverConditionGroupTerraformV1{value: val, isSet: true}
}

func (v NullableApproverConditionGroupTerraformV1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApproverConditionGroupTerraformV1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

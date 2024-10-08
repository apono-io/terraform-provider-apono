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

// checks if the PaginatedResponseAccessFlowTerraformModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PaginatedResponseAccessFlowTerraformModel{}

// PaginatedResponseAccessFlowTerraformModel struct for PaginatedResponseAccessFlowTerraformModel
type PaginatedResponseAccessFlowTerraformModel struct {
	Data       []AccessFlowTerraformV1 `json:"data"`
	Pagination PaginationInfo          `json:"pagination"`
}

type _PaginatedResponseAccessFlowTerraformModel PaginatedResponseAccessFlowTerraformModel

// NewPaginatedResponseAccessFlowTerraformModel instantiates a new PaginatedResponseAccessFlowTerraformModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginatedResponseAccessFlowTerraformModel(data []AccessFlowTerraformV1, pagination PaginationInfo) *PaginatedResponseAccessFlowTerraformModel {
	this := PaginatedResponseAccessFlowTerraformModel{}
	this.Data = data
	this.Pagination = pagination
	return &this
}

// NewPaginatedResponseAccessFlowTerraformModelWithDefaults instantiates a new PaginatedResponseAccessFlowTerraformModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginatedResponseAccessFlowTerraformModelWithDefaults() *PaginatedResponseAccessFlowTerraformModel {
	this := PaginatedResponseAccessFlowTerraformModel{}
	return &this
}

// GetData returns the Data field value
func (o *PaginatedResponseAccessFlowTerraformModel) GetData() []AccessFlowTerraformV1 {
	if o == nil {
		var ret []AccessFlowTerraformV1
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *PaginatedResponseAccessFlowTerraformModel) GetDataOk() ([]AccessFlowTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *PaginatedResponseAccessFlowTerraformModel) SetData(v []AccessFlowTerraformV1) {
	o.Data = v
}

// GetPagination returns the Pagination field value
func (o *PaginatedResponseAccessFlowTerraformModel) GetPagination() PaginationInfo {
	if o == nil {
		var ret PaginationInfo
		return ret
	}

	return o.Pagination
}

// GetPaginationOk returns a tuple with the Pagination field value
// and a boolean to check if the value has been set.
func (o *PaginatedResponseAccessFlowTerraformModel) GetPaginationOk() (*PaginationInfo, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Pagination, true
}

// SetPagination sets field value
func (o *PaginatedResponseAccessFlowTerraformModel) SetPagination(v PaginationInfo) {
	o.Pagination = v
}

func (o PaginatedResponseAccessFlowTerraformModel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PaginatedResponseAccessFlowTerraformModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	toSerialize["pagination"] = o.Pagination
	return toSerialize, nil
}

func (o *PaginatedResponseAccessFlowTerraformModel) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"data",
		"pagination",
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

	varPaginatedResponseAccessFlowTerraformModel := _PaginatedResponseAccessFlowTerraformModel{}

	err = json.Unmarshal(bytes, &varPaginatedResponseAccessFlowTerraformModel)

	if err != nil {
		return err
	}

	*o = PaginatedResponseAccessFlowTerraformModel(varPaginatedResponseAccessFlowTerraformModel)

	return err
}

type NullablePaginatedResponseAccessFlowTerraformModel struct {
	value *PaginatedResponseAccessFlowTerraformModel
	isSet bool
}

func (v NullablePaginatedResponseAccessFlowTerraformModel) Get() *PaginatedResponseAccessFlowTerraformModel {
	return v.value
}

func (v *NullablePaginatedResponseAccessFlowTerraformModel) Set(val *PaginatedResponseAccessFlowTerraformModel) {
	v.value = val
	v.isSet = true
}

func (v NullablePaginatedResponseAccessFlowTerraformModel) IsSet() bool {
	return v.isSet
}

func (v *NullablePaginatedResponseAccessFlowTerraformModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaginatedResponseAccessFlowTerraformModel(val *PaginatedResponseAccessFlowTerraformModel) *NullablePaginatedResponseAccessFlowTerraformModel {
	return &NullablePaginatedResponseAccessFlowTerraformModel{value: val, isSet: true}
}

func (v NullablePaginatedResponseAccessFlowTerraformModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaginatedResponseAccessFlowTerraformModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

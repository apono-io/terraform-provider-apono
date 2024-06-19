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

// checks if the UpsertAccessFlowTerraformV1 type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpsertAccessFlowTerraformV1{}

// UpsertAccessFlowTerraformV1 struct for UpsertAccessFlowTerraformV1
type UpsertAccessFlowTerraformV1 struct {
	Name               string                                          `json:"name"`
	Active             bool                                            `json:"active"`
	Trigger            AccessFlowTriggerTerraformV1                    `json:"trigger"`
	Grantees           []GranteeTerraformV1                            `json:"grantees"`
	GranteeFilterGroup NullableAccessFlowTerraformV1GranteeFilterGroup `json:"grantee_filter_group,omitempty"`
	IntegrationTargets []AccessTargetIntegrationTerraformV1            `json:"integration_targets,omitempty"`
	BundleTargets      []AccessTargetBundleTerraformV1                 `json:"bundle_targets,omitempty"`
	Approvers          []ApproverTerraformV1                           `json:"approvers,omitempty"`
	ApproverPolicy     NullableAccessFlowTerraformV1ApproverPolicy     `json:"approver_policy,omitempty"`
	RevokeAfterInSec   int32                                           `json:"revoke_after_in_sec"`
	Settings           NullableAccessFlowTerraformV1Settings           `json:"settings,omitempty"`
	Labels             []AccessFlowLabelTerraformV1                    `json:"labels"`
}

type _UpsertAccessFlowTerraformV1 UpsertAccessFlowTerraformV1

// NewUpsertAccessFlowTerraformV1 instantiates a new UpsertAccessFlowTerraformV1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpsertAccessFlowTerraformV1(name string, active bool, trigger AccessFlowTriggerTerraformV1, grantees []GranteeTerraformV1, revokeAfterInSec int32, labels []AccessFlowLabelTerraformV1) *UpsertAccessFlowTerraformV1 {
	this := UpsertAccessFlowTerraformV1{}
	this.Name = name
	this.Active = active
	this.Trigger = trigger
	this.Grantees = grantees
	this.RevokeAfterInSec = revokeAfterInSec
	this.Labels = labels
	return &this
}

// NewUpsertAccessFlowTerraformV1WithDefaults instantiates a new UpsertAccessFlowTerraformV1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpsertAccessFlowTerraformV1WithDefaults() *UpsertAccessFlowTerraformV1 {
	this := UpsertAccessFlowTerraformV1{}
	return &this
}

// GetName returns the Name field value
func (o *UpsertAccessFlowTerraformV1) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *UpsertAccessFlowTerraformV1) SetName(v string) {
	o.Name = v
}

// GetActive returns the Active field value
func (o *UpsertAccessFlowTerraformV1) GetActive() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Active
}

// GetActiveOk returns a tuple with the Active field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetActiveOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Active, true
}

// SetActive sets field value
func (o *UpsertAccessFlowTerraformV1) SetActive(v bool) {
	o.Active = v
}

// GetTrigger returns the Trigger field value
func (o *UpsertAccessFlowTerraformV1) GetTrigger() AccessFlowTriggerTerraformV1 {
	if o == nil {
		var ret AccessFlowTriggerTerraformV1
		return ret
	}

	return o.Trigger
}

// GetTriggerOk returns a tuple with the Trigger field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetTriggerOk() (*AccessFlowTriggerTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Trigger, true
}

// SetTrigger sets field value
func (o *UpsertAccessFlowTerraformV1) SetTrigger(v AccessFlowTriggerTerraformV1) {
	o.Trigger = v
}

// GetGrantees returns the Grantees field value
func (o *UpsertAccessFlowTerraformV1) GetGrantees() []GranteeTerraformV1 {
	if o == nil {
		var ret []GranteeTerraformV1
		return ret
	}

	return o.Grantees
}

// GetGranteesOk returns a tuple with the Grantees field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetGranteesOk() ([]GranteeTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.Grantees, true
}

// SetGrantees sets field value
func (o *UpsertAccessFlowTerraformV1) SetGrantees(v []GranteeTerraformV1) {
	o.Grantees = v
}

// GetGranteeFilterGroup returns the GranteeFilterGroup field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetGranteeFilterGroup() AccessFlowTerraformV1GranteeFilterGroup {
	if o == nil || IsNil(o.GranteeFilterGroup.Get()) {
		var ret AccessFlowTerraformV1GranteeFilterGroup
		return ret
	}
	return *o.GranteeFilterGroup.Get()
}

// GetGranteeFilterGroupOk returns a tuple with the GranteeFilterGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetGranteeFilterGroupOk() (*AccessFlowTerraformV1GranteeFilterGroup, bool) {
	if o == nil {
		return nil, false
	}
	return o.GranteeFilterGroup.Get(), o.GranteeFilterGroup.IsSet()
}

// HasGranteeFilterGroup returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasGranteeFilterGroup() bool {
	if o != nil && o.GranteeFilterGroup.IsSet() {
		return true
	}

	return false
}

// SetGranteeFilterGroup gets a reference to the given NullableAccessFlowTerraformV1GranteeFilterGroup and assigns it to the GranteeFilterGroup field.
func (o *UpsertAccessFlowTerraformV1) SetGranteeFilterGroup(v AccessFlowTerraformV1GranteeFilterGroup) {
	o.GranteeFilterGroup.Set(&v)
}

// SetGranteeFilterGroupNil sets the value for GranteeFilterGroup to be an explicit nil
func (o *UpsertAccessFlowTerraformV1) SetGranteeFilterGroupNil() {
	o.GranteeFilterGroup.Set(nil)
}

// UnsetGranteeFilterGroup ensures that no value is present for GranteeFilterGroup, not even an explicit nil
func (o *UpsertAccessFlowTerraformV1) UnsetGranteeFilterGroup() {
	o.GranteeFilterGroup.Unset()
}

// GetIntegrationTargets returns the IntegrationTargets field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetIntegrationTargets() []AccessTargetIntegrationTerraformV1 {
	if o == nil {
		var ret []AccessTargetIntegrationTerraformV1
		return ret
	}
	return o.IntegrationTargets
}

// GetIntegrationTargetsOk returns a tuple with the IntegrationTargets field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetIntegrationTargetsOk() ([]AccessTargetIntegrationTerraformV1, bool) {
	if o == nil || IsNil(o.IntegrationTargets) {
		return nil, false
	}
	return o.IntegrationTargets, true
}

// HasIntegrationTargets returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasIntegrationTargets() bool {
	if o != nil && IsNil(o.IntegrationTargets) {
		return true
	}

	return false
}

// SetIntegrationTargets gets a reference to the given []AccessTargetIntegrationTerraformV1 and assigns it to the IntegrationTargets field.
func (o *UpsertAccessFlowTerraformV1) SetIntegrationTargets(v []AccessTargetIntegrationTerraformV1) {
	o.IntegrationTargets = v
}

// GetBundleTargets returns the BundleTargets field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetBundleTargets() []AccessTargetBundleTerraformV1 {
	if o == nil {
		var ret []AccessTargetBundleTerraformV1
		return ret
	}
	return o.BundleTargets
}

// GetBundleTargetsOk returns a tuple with the BundleTargets field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetBundleTargetsOk() ([]AccessTargetBundleTerraformV1, bool) {
	if o == nil || IsNil(o.BundleTargets) {
		return nil, false
	}
	return o.BundleTargets, true
}

// HasBundleTargets returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasBundleTargets() bool {
	if o != nil && IsNil(o.BundleTargets) {
		return true
	}

	return false
}

// SetBundleTargets gets a reference to the given []AccessTargetBundleTerraformV1 and assigns it to the BundleTargets field.
func (o *UpsertAccessFlowTerraformV1) SetBundleTargets(v []AccessTargetBundleTerraformV1) {
	o.BundleTargets = v
}

// GetApprovers returns the Approvers field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetApprovers() []ApproverTerraformV1 {
	if o == nil {
		var ret []ApproverTerraformV1
		return ret
	}
	return o.Approvers
}

// GetApproversOk returns a tuple with the Approvers field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetApproversOk() ([]ApproverTerraformV1, bool) {
	if o == nil || IsNil(o.Approvers) {
		return nil, false
	}
	return o.Approvers, true
}

// HasApprovers returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasApprovers() bool {
	if o != nil && IsNil(o.Approvers) {
		return true
	}

	return false
}

// SetApprovers gets a reference to the given []ApproverTerraformV1 and assigns it to the Approvers field.
func (o *UpsertAccessFlowTerraformV1) SetApprovers(v []ApproverTerraformV1) {
	o.Approvers = v
}

// GetApproverPolicy returns the ApproverPolicy field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetApproverPolicy() AccessFlowTerraformV1ApproverPolicy {
	if o == nil || IsNil(o.ApproverPolicy.Get()) {
		var ret AccessFlowTerraformV1ApproverPolicy
		return ret
	}
	return *o.ApproverPolicy.Get()
}

// GetApproverPolicyOk returns a tuple with the ApproverPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetApproverPolicyOk() (*AccessFlowTerraformV1ApproverPolicy, bool) {
	if o == nil {
		return nil, false
	}
	return o.ApproverPolicy.Get(), o.ApproverPolicy.IsSet()
}

// HasApproverPolicy returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasApproverPolicy() bool {
	if o != nil && o.ApproverPolicy.IsSet() {
		return true
	}

	return false
}

// SetApproverPolicy gets a reference to the given NullableAccessFlowTerraformV1ApproverPolicy and assigns it to the ApproverPolicy field.
func (o *UpsertAccessFlowTerraformV1) SetApproverPolicy(v AccessFlowTerraformV1ApproverPolicy) {
	o.ApproverPolicy.Set(&v)
}

// SetApproverPolicyNil sets the value for ApproverPolicy to be an explicit nil
func (o *UpsertAccessFlowTerraformV1) SetApproverPolicyNil() {
	o.ApproverPolicy.Set(nil)
}

// UnsetApproverPolicy ensures that no value is present for ApproverPolicy, not even an explicit nil
func (o *UpsertAccessFlowTerraformV1) UnsetApproverPolicy() {
	o.ApproverPolicy.Unset()
}

// GetRevokeAfterInSec returns the RevokeAfterInSec field value
func (o *UpsertAccessFlowTerraformV1) GetRevokeAfterInSec() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.RevokeAfterInSec
}

// GetRevokeAfterInSecOk returns a tuple with the RevokeAfterInSec field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetRevokeAfterInSecOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RevokeAfterInSec, true
}

// SetRevokeAfterInSec sets field value
func (o *UpsertAccessFlowTerraformV1) SetRevokeAfterInSec(v int32) {
	o.RevokeAfterInSec = v
}

// GetSettings returns the Settings field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *UpsertAccessFlowTerraformV1) GetSettings() AccessFlowTerraformV1Settings {
	if o == nil || IsNil(o.Settings.Get()) {
		var ret AccessFlowTerraformV1Settings
		return ret
	}
	return *o.Settings.Get()
}

// GetSettingsOk returns a tuple with the Settings field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UpsertAccessFlowTerraformV1) GetSettingsOk() (*AccessFlowTerraformV1Settings, bool) {
	if o == nil {
		return nil, false
	}
	return o.Settings.Get(), o.Settings.IsSet()
}

// HasSettings returns a boolean if a field has been set.
func (o *UpsertAccessFlowTerraformV1) HasSettings() bool {
	if o != nil && o.Settings.IsSet() {
		return true
	}

	return false
}

// SetSettings gets a reference to the given NullableAccessFlowTerraformV1Settings and assigns it to the Settings field.
func (o *UpsertAccessFlowTerraformV1) SetSettings(v AccessFlowTerraformV1Settings) {
	o.Settings.Set(&v)
}

// SetSettingsNil sets the value for Settings to be an explicit nil
func (o *UpsertAccessFlowTerraformV1) SetSettingsNil() {
	o.Settings.Set(nil)
}

// UnsetSettings ensures that no value is present for Settings, not even an explicit nil
func (o *UpsertAccessFlowTerraformV1) UnsetSettings() {
	o.Settings.Unset()
}

// GetLabels returns the Labels field value
func (o *UpsertAccessFlowTerraformV1) GetLabels() []AccessFlowLabelTerraformV1 {
	if o == nil {
		var ret []AccessFlowLabelTerraformV1
		return ret
	}

	return o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value
// and a boolean to check if the value has been set.
func (o *UpsertAccessFlowTerraformV1) GetLabelsOk() ([]AccessFlowLabelTerraformV1, bool) {
	if o == nil {
		return nil, false
	}
	return o.Labels, true
}

// SetLabels sets field value
func (o *UpsertAccessFlowTerraformV1) SetLabels(v []AccessFlowLabelTerraformV1) {
	o.Labels = v
}

func (o UpsertAccessFlowTerraformV1) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpsertAccessFlowTerraformV1) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	toSerialize["active"] = o.Active
	toSerialize["trigger"] = o.Trigger
	toSerialize["grantees"] = o.Grantees
	if o.GranteeFilterGroup.IsSet() {
		toSerialize["grantee_filter_group"] = o.GranteeFilterGroup.Get()
	}
	if o.IntegrationTargets != nil {
		toSerialize["integration_targets"] = o.IntegrationTargets
	}
	if o.BundleTargets != nil {
		toSerialize["bundle_targets"] = o.BundleTargets
	}
	if o.Approvers != nil {
		toSerialize["approvers"] = o.Approvers
	}
	if o.ApproverPolicy.IsSet() {
		toSerialize["approver_policy"] = o.ApproverPolicy.Get()
	}
	toSerialize["revoke_after_in_sec"] = o.RevokeAfterInSec
	if o.Settings.IsSet() {
		toSerialize["settings"] = o.Settings.Get()
	}
	toSerialize["labels"] = o.Labels
	return toSerialize, nil
}

func (o *UpsertAccessFlowTerraformV1) UnmarshalJSON(bytes []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
		"active",
		"trigger",
		"grantees",
		"revoke_after_in_sec",
		"labels",
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

	varUpsertAccessFlowTerraformV1 := _UpsertAccessFlowTerraformV1{}

	err = json.Unmarshal(bytes, &varUpsertAccessFlowTerraformV1)

	if err != nil {
		return err
	}

	*o = UpsertAccessFlowTerraformV1(varUpsertAccessFlowTerraformV1)

	return err
}

type NullableUpsertAccessFlowTerraformV1 struct {
	value *UpsertAccessFlowTerraformV1
	isSet bool
}

func (v NullableUpsertAccessFlowTerraformV1) Get() *UpsertAccessFlowTerraformV1 {
	return v.value
}

func (v *NullableUpsertAccessFlowTerraformV1) Set(val *UpsertAccessFlowTerraformV1) {
	v.value = val
	v.isSet = true
}

func (v NullableUpsertAccessFlowTerraformV1) IsSet() bool {
	return v.isSet
}

func (v *NullableUpsertAccessFlowTerraformV1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpsertAccessFlowTerraformV1(val *UpsertAccessFlowTerraformV1) *NullableUpsertAccessFlowTerraformV1 {
	return &NullableUpsertAccessFlowTerraformV1{value: val, isSet: true}
}

func (v NullableUpsertAccessFlowTerraformV1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpsertAccessFlowTerraformV1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

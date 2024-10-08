package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AccessFlowModel describes the resource data model.
type AccessFlowModel struct {
	ID                  types.String        `tfsdk:"id"`
	Name                types.String        `tfsdk:"name"`
	Active              types.Bool          `tfsdk:"active"`
	RevokeAfterInSec    types.Number        `tfsdk:"revoke_after_in_sec"`
	Trigger             *Trigger            `tfsdk:"trigger"`
	Grantees            types.Set           `tfsdk:"grantees"`
	GranteesFilterGroup types.Object        `tfsdk:"grantees_conditions_group"`
	IntegrationTargets  []IntegrationTarget `tfsdk:"integration_targets"`
	BundleTargets       []BundleTarget      `tfsdk:"bundle_targets"`
	Approvers           types.Set           `tfsdk:"approvers"`
	ApproverPolicy      types.Object        `tfsdk:"approver_policy"`
	Settings            *Settings           `tfsdk:"settings"`
	Labels              types.List          `tfsdk:"labels"`
}

type Trigger struct {
	Type      types.String `tfsdk:"type"`
	Timeframe *Timeframe   `tfsdk:"timeframe"`
}

type Timeframe struct {
	StartTime  types.String `tfsdk:"start_time"`
	EndTime    types.String `tfsdk:"end_time"`
	DaysInWeek types.Set    `tfsdk:"days_in_week"`
	TimeZone   types.String `tfsdk:"time_zone"`
}

type Settings struct {
	RequireJustificationOnRequestAgain types.Bool `tfsdk:"require_justification_on_request_again"`
	RequireAllApprovers                types.Bool `tfsdk:"require_all_approvers"`
	ApproverCannotSelfApprove          types.Bool `tfsdk:"approver_cannot_self_approve"`
	RequireMFA                         types.Bool `tfsdk:"require_mfa"`
}

type Identity struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

var IdentityObject = map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
}

type GranteeFilterGroup struct {
	Operator types.String `tfsdk:"conditions_logical_operator"`
	Filters  types.Set    `tfsdk:"attribute_conditions"`
}

type ApproverPolicy struct {
	GroupsRelationship types.String `tfsdk:"approver_groups_relationship"`
	Groups             types.Set    `tfsdk:"approver_groups"`
}

type ApproverGroup struct {
	Operator   types.String `tfsdk:"conditions_logical_operator"`
	Conditions types.Set    `tfsdk:"attribute_conditions"`
}

type AttributeFilter struct {
	Operator       types.String `tfsdk:"operator"`
	AttributeType  types.String `tfsdk:"attribute_type"`
	AttributeNames types.Set    `tfsdk:"attribute_names"`
	IntegrationID  types.String `tfsdk:"integration_id"`
}

var ApproverPolicyObject = map[string]attr.Type{
	"approver_groups_relationship": types.StringType,
	"approver_groups":              basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: ApproverPolicyGroupObject}},
}

var ApproverPolicyGroupObject = map[string]attr.Type{
	"conditions_logical_operator": types.StringType,
	"attribute_conditions":        basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: AttributeFilterObject}},
}

var GranteeFilterGroupObject = map[string]attr.Type{
	"conditions_logical_operator": types.StringType,
	"attribute_conditions":        basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: AttributeFilterObject}},
}

var AttributeFilterObject = map[string]attr.Type{
	"operator":        types.StringType,
	"attribute_type":  types.StringType,
	"attribute_names": basetypes.SetType{ElemType: types.StringType},
	"integration_id":  types.StringType,
}

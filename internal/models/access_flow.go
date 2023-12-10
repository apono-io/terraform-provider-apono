package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccessFlowModel describes the resource data model.
type AccessFlowModel struct {
	ID                 types.String        `tfsdk:"id"`
	Name               types.String        `tfsdk:"name"`
	Active             types.Bool          `tfsdk:"active"`
	RevokeAfterInSec   types.Number        `tfsdk:"revoke_after_in_sec"`
	Trigger            *Trigger            `tfsdk:"trigger"`
	Grantees           types.Set           `tfsdk:"grantees"`
	IntegrationTargets []IntegrationTarget `tfsdk:"integration_targets"`
	BundleTargets      []BundleTarget      `tfsdk:"bundle_targets"`
	Approvers          types.Set           `tfsdk:"approvers"`
	Settings           *Settings           `tfsdk:"settings"`
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
}

type Identity struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

var IdentityObject = map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
}

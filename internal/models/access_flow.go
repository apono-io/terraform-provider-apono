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
	Approvers          types.Set           `tfsdk:"approvers"`
	Settings           *Settings           `tfsdk:"settings"`
}

type Trigger struct {
	Type      types.String `tfsdk:"type"`
	Timeframe *Timeframe   `tfsdk:"timeframe"`
}

type Timeframe struct {
	StartOfDayTimeInSeconds types.Number `tfsdk:"start_of_day_time_in_seconds"`
	EndOfDayTimeInSeconds   types.Number `tfsdk:"end_of_day_time_in_seconds"`
	DaysInWeek              types.Set    `tfsdk:"days_in_week"`
	TimeZone                types.String `tfsdk:"time_zone"`
}

type IntegrationTarget struct {
	Name                   types.String     `tfsdk:"name"`
	ResourceType           types.String     `tfsdk:"resource_type"`
	ResourceIncludeFilters []ResourceFilter `tfsdk:"resource_include_filters"`
	ResourceExcludeFilters []ResourceFilter `tfsdk:"resource_excludes_filters"`
	Permissions            types.Set        `tfsdk:"permissions"`
}

type Settings struct {
	RequireJustificationOnRequestAgain types.Bool `tfsdk:"require_justification_on_request_again"`
	RequireAllApprovers                types.Bool `tfsdk:"require_all_approvers"`
	ApproverCannotApproveHimself       types.Bool `tfsdk:"approver_cannot_approve_himself"`
}

type ResourceFilter struct {
	Type  types.String `tfsdk:"type"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type Identity struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

var IdentityObject = map[string]attr.Type{
	"name": types.StringType,
	"type": types.StringType,
}

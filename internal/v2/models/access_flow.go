package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type AccessFlowV2Model struct {
	ID                 types.String                  `tfsdk:"id"`
	Name               types.String                  `tfsdk:"name"`
	Active             types.Bool                    `tfsdk:"active"`
	Trigger            types.String                  `tfsdk:"trigger"`
	GrantDurationInMin types.Int32                   `tfsdk:"grant_duration_in_min"`
	Timeframe          *AccessFlowTimeframeModel     `tfsdk:"timeframe"`
	ApproverPolicy     *AccessFlowApproverPolicy     `tfsdk:"approver_policy"`
	Requestors         *AccessFlowRequestorsModel    `tfsdk:"requestors"`
	AccessTargets      []AccessFlowAccessTargetModel `tfsdk:"access_targets"`
	Settings           *AccessFlowSettingsModel      `tfsdk:"settings"`
}

type AccessFlowTimeframeModel struct {
	StartTime  types.String `tfsdk:"start_time"`
	EndTime    types.String `tfsdk:"end_time"`
	DaysOfWeek types.Set    `tfsdk:"days_of_week"`
	TimeZone   types.String `tfsdk:"time_zone"`
}

type AccessFlowApproverPolicy struct {
	ApprovalMode   types.String              `tfsdk:"approval_mode"`
	ApproverGroups []AccessFlowApproverGroup `tfsdk:"approver_groups"`
}

type AccessFlowApproverGroup struct {
	LogicalOperator types.String          `tfsdk:"logical_operator"`
	Approvers       []AccessFlowCondition `tfsdk:"approvers"`
}

type AccessFlowCondition struct {
	SourceIntegrationName types.String `tfsdk:"source_integration_name"`
	Type                  types.String `tfsdk:"type"`
	MatchOperator         types.String `tfsdk:"match_operator"`
	Values                types.Set    `tfsdk:"values"`
}

type AccessFlowRequestorsModel struct {
	LogicalOperator types.String          `tfsdk:"logical_operator"`
	Conditions      []AccessFlowCondition `tfsdk:"conditions"`
}

type AccessFlowSettingsModel struct {
	JustificationRequired      types.Bool `tfsdk:"justification_required"`
	RequireApproverReason      types.Bool `tfsdk:"require_approver_reason"`
	RequesterCannotApproveSelf types.Bool `tfsdk:"requester_cannot_approve_self"`
	RequireMFA                 types.Bool `tfsdk:"require_mfa"`
	Labels                     types.Set  `tfsdk:"labels"`
}

type AccessFlowTargetBundleModel struct {
	Name types.String `tfsdk:"name"`
}

type AccessFlowAccessTargetModel struct {
	Integration *IntegrationTargetModel      `tfsdk:"integration"`
	Bundle      *AccessFlowTargetBundleModel `tfsdk:"bundle"`
	AccessScope *AccessScopeTargetModel      `tfsdk:"access_scope"`
}

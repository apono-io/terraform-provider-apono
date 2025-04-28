package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationTargetModel struct {
	IntegrationName types.String                  `tfsdk:"integration_name"`
	ResourceType    types.String                  `tfsdk:"resource_type"`
	Permissions     types.Set                     `tfsdk:"permissions"`
	ResourcesScopes []IntegrationTargetScopeModel `tfsdk:"resources_scopes"`
}

type IntegrationTargetScopeModel struct {
	ScopeMode types.String `tfsdk:"scope_mode"`
	Type      types.String `tfsdk:"type"`
	Key       types.String `tfsdk:"key"`
	Values    types.Set    `tfsdk:"values"`
}

type AccessScopeTargetModel struct {
	Name types.String `tfsdk:"name"`
}

type BundleAccessTargetModel struct {
	Integration *IntegrationTargetModel `tfsdk:"integration"`
	AccessScope *AccessScopeTargetModel `tfsdk:"access_scope"`
}

type BundleV2Model struct {
	ID            types.String              `tfsdk:"id"`
	Name          types.String              `tfsdk:"name"`
	AccessTargets []BundleAccessTargetModel `tfsdk:"access_targets"`
}

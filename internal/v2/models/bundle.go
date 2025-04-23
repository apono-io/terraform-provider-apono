package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationTargetModel struct {
	IntegrationName types.String                  `tfsdk:"integration_name"`
	ResourceType    types.String                  `tfsdk:"resource_type"`
	Permissions     types.Set                     `tfsdk:"permissions"`
	ResourcesScope  []IntegrationTargetScopeModel `tfsdk:"resources_scope"`
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

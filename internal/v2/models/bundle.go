package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BundleIntegrationModel struct {
	IntegrationName types.String                  `tfsdk:"integration_name"`
	ResourceType    types.String                  `tfsdk:"resource_type"`
	Permissions     types.Set                     `tfsdk:"permissions"`
	ResourcesScope  []BundleIntegrationScopeModel `tfsdk:"resources_scope"`
}

type BundleIntegrationScopeModel struct {
	ScopeMode types.String `tfsdk:"scope_mode"`
	Type      types.String `tfsdk:"type"`
	Key       types.String `tfsdk:"key"`
	Values    types.Set    `tfsdk:"values"`
}

type BundleAccessScopeModel struct {
	Name types.String `tfsdk:"name"`
}

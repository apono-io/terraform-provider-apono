package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type IntegrationTarget struct {
	Name                   types.String     `tfsdk:"name"`
	ResourceType           types.String     `tfsdk:"resource_type"`
	ResourceIncludeFilters []ResourceFilter `tfsdk:"resource_include_filters"`
	ResourceExcludeFilters []ResourceFilter `tfsdk:"resource_exclude_filters"`
	Permissions            types.Set        `tfsdk:"permissions"`
}

type BundleTarget struct {
	Name types.String `tfsdk:"name"`
}

type ResourceFilter struct {
	Type  types.String `tfsdk:"type"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

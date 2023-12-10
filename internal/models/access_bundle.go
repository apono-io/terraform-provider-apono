package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// AccessBundleModel describes the resource data model.
type AccessBundleModel struct {
	ID                 types.String        `tfsdk:"id"`
	Name               types.String        `tfsdk:"name"`
	IntegrationTargets []IntegrationTarget `tfsdk:"integration_targets"`
}

package common

import "github.com/hashicorp/terraform-plugin-framework/types"

type AccessScopeModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Query types.String `tfsdk:"query"`
}

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AttrValueToString(val attr.Value) string {
	switch value := val.(type) {
	case types.String:
		return value.ValueString()
	default:
		return value.String()
	}
}

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// accessTargetExclusivityValidator enforces that exactly one target kind attribute is set
// per access_targets entry. Validation is skipped for entries that contain unknown values
// (e.g. attributes computed from another resource that are not yet resolved at plan time).
type accessTargetExclusivityValidator struct {
	kindNames []string
}

func (v accessTargetExclusivityValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Exactly one of [%s] must be set per access_targets entry.", strings.Join(v.kindNames, ", "))
}

func (v accessTargetExclusivityValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v accessTargetExclusivityValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var targets types.List
	diags := req.Config.GetAttribute(ctx, path.Root("access_targets"), &targets)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if targets.IsNull() || targets.IsUnknown() {
		return
	}

	for i, elem := range targets.Elements() {
		obj, ok := elem.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}

		attrs := obj.Attributes()
		count, hasUnknown := countSetKinds(attrs, v.kindNames)

		if hasUnknown && count <= 1 {
			// Fewer than two kinds known so far; unknowns may still resolve to a valid config.
			continue
		}
		if count == 1 {
			continue
		}

		resp.Diagnostics.AddAttributeError(
			path.Root("access_targets").AtListIndex(i),
			"Invalid access_targets configuration",
			fmt.Sprintf(
				"Exactly one of [%s] must be configured per access target entry (got %d).",
				strings.Join(v.kindNames, ", "),
				count,
			),
		)
	}
}

// countSetKinds returns the number of non-null, non-unknown attributes found in attrs
// for the given kind names, and hasUnknown=true if any kind attribute is unknown.
func countSetKinds(attrs map[string]attr.Value, kindNames []string) (count int, hasUnknown bool) {
	for _, name := range kindNames {
		val, exists := attrs[name]
		if !exists {
			continue
		}
		if val.IsUnknown() {
			hasUnknown = true
			continue
		}
		if !val.IsNull() {
			count++
		}
	}
	return count, hasUnknown
}

package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestCountSetKinds(t *testing.T) {
	t.Run("exactly one set - valid", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration":  types.StringNull(),
			"access_scope": types.StringNull(),
			"query_target": types.StringValue(`integration_name = "aws"`),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope", "query_target"})
		assert.Equal(t, 1, count)
		assert.False(t, hasUnknown)
	})

	t.Run("multiple set - invalid", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration":  types.StringValue("my-integration"),
			"access_scope": types.StringNull(),
			"query_target": types.StringValue(`integration_name = "aws"`),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope", "query_target"})
		assert.Equal(t, 2, count)
		assert.False(t, hasUnknown)
	})

	t.Run("none set - invalid", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration":  types.StringNull(),
			"access_scope": types.StringNull(),
			"query_target": types.StringNull(),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope", "query_target"})
		assert.Equal(t, 0, count)
		assert.False(t, hasUnknown)
	})

	t.Run("unknown present - defer to apply", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration":  types.StringUnknown(),
			"access_scope": types.StringNull(),
			"query_target": types.StringNull(),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope", "query_target"})
		assert.Equal(t, 0, count)
		assert.True(t, hasUnknown)
	})

	t.Run("kind missing from attrs - skipped", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration": types.StringValue("my-integration"),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "bundle"})
		assert.Equal(t, 1, count)
		assert.False(t, hasUnknown)
	})

	t.Run("multiple set plus unknown - still invalid", func(t *testing.T) {
		// integration and access_scope are both explicitly set (count=2 already),
		// query_target is unknown. The bug: early return on unknown discards count=2.
		attrs := map[string]attr.Value{
			"integration":  types.StringValue("my-integration"),
			"access_scope": types.StringValue("my-scope"),
			"query_target": types.StringUnknown(),
		}
		count, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope", "query_target"})
		assert.Equal(t, 2, count)
		assert.True(t, hasUnknown)
	})

	t.Run("all unknown - hasUnknown on first hit", func(t *testing.T) {
		attrs := map[string]attr.Value{
			"integration":  types.StringUnknown(),
			"access_scope": types.StringUnknown(),
		}
		_, hasUnknown := countSetKinds(attrs, []string{"integration", "access_scope"})
		assert.True(t, hasUnknown)
	})
}

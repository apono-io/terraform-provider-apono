package resources

import (
	"testing"
	"time"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ogen-go/ogen/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoAccessScopeResource_Unit(t *testing.T) {
	// Create mock invoker and resource
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoAccessScopeResource{
		client: mockInvoker,
	}

	t.Run("Create", func(t *testing.T) {
		// Setup mock
		mockInvoker.EXPECT().
			CreateAccessScopesV1(mock.Anything, mock.MatchedBy(func(req *client.UpsertAccessScopeV1) bool {
				return req.Name == "test-scope" && req.Query == "tag:environment=dev"
			})).
			Return(&client.AccessScopeV1{
				ID:           "as-123456",
				Name:         "test-scope",
				Query:        "tag:environment=dev",
				CreationDate: client.ApiInstant(time.Now()),
				UpdateDate:   client.ApiInstant(time.Now()),
			}, nil).
			Once()

		// Create test plan data
		ctx := t.Context()
		planType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":            tftypes.NewValue(tftypes.String, nil),
			"name":          tftypes.NewValue(tftypes.String, "test-scope"),
			"query":         tftypes.NewValue(tftypes.String, "tag:environment=dev"),
			"creation_date": tftypes.NewValue(tftypes.String, nil),
			"update_date":   tftypes.NewValue(tftypes.String, nil),
		})

		// Create schema
		schema := r.getTestSchema()

		// Create proper plan with schema
		plan := tfsdk.Plan{
			Schema: schema,
			Raw:    planVal,
		}

		// Initialize empty state
		state := tfsdk.State{
			Schema: schema,
			Raw:    tftypes.NewValue(planType, nil),
		}

		// Create request
		req := resource.CreateRequest{
			Plan: plan,
		}
		resp := resource.CreateResponse{
			State: state,
		}

		// Call create
		r.Create(ctx, req, &resp)

		// Verify no errors
		require.False(t, resp.Diagnostics.HasError(), "create should not error")

		// Verify state values
		var stateVal accessScopeResourceModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError(), "getting state should not error")

		assert.Equal(t, "as-123456", stateVal.ID.ValueString())
		assert.Equal(t, "test-scope", stateVal.Name.ValueString())
		assert.Equal(t, "tag:environment=dev", stateVal.Query.ValueString())
		assert.NotEmpty(t, stateVal.CreationDate.ValueString())
		assert.NotEmpty(t, stateVal.UpdateDate.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		// Setup mock
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-123456"
			})).
			Return(&client.AccessScopeV1{
				ID:           "as-123456",
				Name:         "test-scope",
				Query:        "tag:environment=dev",
				CreationDate: client.ApiInstant(time.Now()),
				UpdateDate:   client.ApiInstant(time.Now()),
			}, nil).
			Once()

		// Create test state data
		ctx := t.Context()
		stateType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":            tftypes.NewValue(tftypes.String, "as-123456"),
			"name":          tftypes.NewValue(tftypes.String, "old-name"),
			"query":         tftypes.NewValue(tftypes.String, "old-query"),
			"creation_date": tftypes.NewValue(tftypes.String, ""),
			"update_date":   tftypes.NewValue(tftypes.String, ""),
		})

		// Create schema
		schema := r.getTestSchema()

		// Create state
		state := tfsdk.State{
			Schema: schema,
			Raw:    stateVal,
		}

		// Create request
		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		// Call read
		r.Read(ctx, req, &resp)

		// Verify no errors
		require.False(t, resp.Diagnostics.HasError(), "read should not error")

		// Verify state values
		var stateModel accessScopeResourceModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError(), "getting state should not error")

		assert.Equal(t, "as-123456", stateModel.ID.ValueString())
		assert.Equal(t, "test-scope", stateModel.Name.ValueString())
		assert.Equal(t, "tag:environment=dev", stateModel.Query.ValueString())
		assert.NotEmpty(t, stateModel.CreationDate.ValueString())
		assert.NotEmpty(t, stateModel.UpdateDate.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		// Setup mock for a not found error
		notFoundErr := &validate.UnexpectedStatusCodeError{StatusCode: 404}
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-not-found"
			})).
			Return(nil, notFoundErr).
			Once()

		// Create test state data
		ctx := t.Context()
		stateType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":            tftypes.NewValue(tftypes.String, "as-not-found"),
			"name":          tftypes.NewValue(tftypes.String, "test-scope"),
			"query":         tftypes.NewValue(tftypes.String, "tag:environment=dev"),
			"creation_date": tftypes.NewValue(tftypes.String, ""),
			"update_date":   tftypes.NewValue(tftypes.String, ""),
		})

		// Create schema
		schema := r.getTestSchema()

		// Create state
		state := tfsdk.State{
			Schema: schema,
			Raw:    stateVal,
		}

		// Create request
		req := resource.ReadRequest{
			State: state,
		}
		resp := resource.ReadResponse{
			State: state,
		}

		// Call read
		r.Read(ctx, req, &resp)

		// Verify the state was removed (resource no longer exists)
		require.False(t, resp.Diagnostics.HasError(), "read should not error on 404")
		assert.True(t, resp.State.Raw.IsNull(), "state should be removed on 404")
	})
}

// Helper method to get a test schema.
func (r *AponoAccessScopeResource) getTestSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"query": schema.StringAttribute{
				Required: true,
			},
			"creation_date": schema.StringAttribute{
				Computed: true,
			},
			"update_date": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

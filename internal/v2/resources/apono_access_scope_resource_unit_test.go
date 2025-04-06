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

	t.Run("ImportState_ByID", func(t *testing.T) {
		// Setup mock for successful import by ID
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-import-id"
			})).
			Return(&client.AccessScopeV1{
				ID:           "as-import-id",
				Name:         "imported-scope",
				Query:        "tag:imported=true",
				CreationDate: client.ApiInstant(time.Now()),
				UpdateDate:   client.ApiInstant(time.Now()),
			}, nil).
			Times(2) // Expect two calls: one for ImportState and one for Read

		// Create test import request
		ctx := t.Context()
		req := resource.ImportStateRequest{
			ID: "as-import-id",
		}

		// Initialize state with a schema
		stateType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		schema := r.getTestSchema()
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schema,
				Raw:    tftypes.NewValue(stateType, nil), // Initialize with a properly typed empty state
			},
		}

		// Call ImportState
		r.ImportState(ctx, req, &importResp)

		// Verify no errors
		require.False(t, importResp.Diagnostics.HasError(), "import should not error")

		// Now call Read to populate the full state, as Terraform would do
		readReq := resource.ReadRequest{
			State: importResp.State,
		}
		readResp := resource.ReadResponse{
			State: importResp.State,
		}
		r.Read(ctx, readReq, &readResp)

		// Verify no errors
		require.False(t, readResp.Diagnostics.HasError(), "read should not error")

		// Verify state values
		var stateModel accessScopeResourceModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError(), "getting state should not error")

		assert.Equal(t, "as-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "imported-scope", stateModel.Name.ValueString())
		assert.Equal(t, "tag:imported=true", stateModel.Query.ValueString())
		assert.NotEmpty(t, stateModel.CreationDate.ValueString())
		assert.NotEmpty(t, stateModel.UpdateDate.ValueString())
	})

	t.Run("ImportState_ByName", func(t *testing.T) {
		// Setup mock for successful import by Name
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "imported-scope" // This must match the req.ID parameter
			})).
			Return(nil, &validate.UnexpectedStatusCodeError{StatusCode: 404}).
			Once()

		mockInvoker.EXPECT().
			ListAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListAccessScopesV1Params) bool {
				return params.Name.Value == "imported-scope"
			})).
			Return(&client.PublicApiListResponseAccessScopePublicV1Model{
				Items: []client.AccessScopeV1{
					{
						ID:           "as-import-name",
						Name:         "imported-scope",
						Query:        "tag:imported=true",
						CreationDate: client.ApiInstant(time.Now()),
						UpdateDate:   client.ApiInstant(time.Now()),
					},
				},
				Pagination: client.PublicApiPaginationInfoModel{
					NextPageToken: client.OptNilString{
						Value: "",
						Null:  true,
					},
				},
			}, nil).
			Once()

		// Add mock for the Read call that will follow the import
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-import-name"
			})).
			Return(&client.AccessScopeV1{
				ID:           "as-import-name",
				Name:         "imported-scope",
				Query:        "tag:imported=true",
				CreationDate: client.ApiInstant(time.Now()),
				UpdateDate:   client.ApiInstant(time.Now()),
			}, nil).
			Once()

		// Create test import request
		ctx := t.Context()
		req := resource.ImportStateRequest{
			ID: "imported-scope",
		}

		// Initialize state properly with a schema and type
		stateType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		schema := r.getTestSchema()
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schema,
				Raw:    tftypes.NewValue(stateType, nil), // Initialize with a properly typed empty state
			},
		}

		// Call ImportState
		r.ImportState(ctx, req, &importResp)

		// Verify no errors
		require.False(t, importResp.Diagnostics.HasError(), "import should not error")

		// Now call Read to populate the full state as would happen in real Terraform operation
		readReq := resource.ReadRequest{
			State: importResp.State,
		}
		readResp := resource.ReadResponse{
			State: importResp.State,
		}
		r.Read(ctx, readReq, &readResp)

		// Verify no errors
		require.False(t, readResp.Diagnostics.HasError(), "read should not error")

		// Verify state values
		var stateModel accessScopeResourceModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError(), "getting state should not error")

		assert.Equal(t, "as-import-name", stateModel.ID.ValueString())
		assert.Equal(t, "imported-scope", stateModel.Name.ValueString())
		assert.Equal(t, "tag:imported=true", stateModel.Query.ValueString())
		assert.NotEmpty(t, stateModel.CreationDate.ValueString())
		assert.NotEmpty(t, stateModel.UpdateDate.ValueString())
	})

	t.Run("ImportState_NotFound", func(t *testing.T) {
		// Setup mock for unsuccessful import
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "not-found-scope"
			})).
			Return(nil, &validate.UnexpectedStatusCodeError{StatusCode: 404}).
			Once()

		mockInvoker.EXPECT().
			ListAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListAccessScopesV1Params) bool {
				return params.Name.Value == "not-found-scope"
			})).
			Return(&client.PublicApiListResponseAccessScopePublicV1Model{
				Items:      []client.AccessScopeV1{},
				Pagination: client.PublicApiPaginationInfoModel{},
			}, nil).
			Once()

		// Create test import request
		ctx := t.Context()
		req := resource.ImportStateRequest{
			ID: "not-found-scope",
		}

		// Initialize state properly with a schema and type
		stateType := tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":            tftypes.String,
				"name":          tftypes.String,
				"query":         tftypes.String,
				"creation_date": tftypes.String,
				"update_date":   tftypes.String,
			},
		}
		schema := r.getTestSchema()
		resp := resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: schema,
				Raw:    tftypes.NewValue(stateType, nil), // Initialize with a properly typed empty state
			},
		}

		// Call ImportState
		r.ImportState(ctx, req, &resp)

		// Verify errors
		require.True(t, resp.Diagnostics.HasError(), "import should error")
		assert.Contains(t, resp.Diagnostics[0].Summary(), "Resource Not Found")
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

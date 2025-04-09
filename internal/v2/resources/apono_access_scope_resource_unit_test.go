package resources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/services"
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
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoAccessScopeResource{client: mockInvoker}

	t.Run("Create", func(t *testing.T) {
		mockInvoker.EXPECT().
			CreateAccessScopesV1(mock.Anything, mock.MatchedBy(func(req *client.UpsertAccessScopeV1) bool {
				return req.Name == "test-scope" && req.Query == `resource_type = "mock-duck"`
			})).
			Return(&client.AccessScopeV1{
				ID:    "as-123456",
				Name:  "test-scope",
				Query: `resource_type = "mock-duck"`,
			}, nil).
			Once()

		ctx := t.Context()
		planType := getStateType()
		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, nil),
			"name":  tftypes.NewValue(tftypes.String, "test-scope"),
			"query": tftypes.NewValue(tftypes.String, `resource_type = "mock-duck"`),
		})

		schema := r.getTestSchema(ctx)
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(planType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}

		r.Create(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateVal services.AccessScopeModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "as-123456", stateVal.ID.ValueString())
		assert.Equal(t, "test-scope", stateVal.Name.ValueString())
		assert.Equal(t, `resource_type = "mock-duck"`, stateVal.Query.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-123456"
			})).
			Return(&client.AccessScopeV1{
				ID:    "as-123456",
				Name:  "test-scope",
				Query: `resource_type = "mock-duck"`,
			}, nil).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "as-123456"),
			"name":  tftypes.NewValue(tftypes.String, "old-name"),
			"query": tftypes.NewValue(tftypes.String, `resource_type = "valid-resource"`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateModel services.AccessScopeModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "as-123456", stateModel.ID.ValueString())
		assert.Equal(t, "test-scope", stateModel.Name.ValueString())
		assert.Equal(t, `resource_type = "mock-duck"`, stateModel.Query.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		notFoundErr := &validate.UnexpectedStatusCodeError{StatusCode: 404}
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-not-found"
			})).
			Return(nil, notFoundErr).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "as-not-found"),
			"name":  tftypes.NewValue(tftypes.String, "test-scope"),
			"query": tftypes.NewValue(tftypes.String, `resource_type = "mock-duck"`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("ImportState_ByID", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.GetAccessScopesV1Params) bool {
				return params.ID == "as-import-id"
			})).
			Return(&client.AccessScopeV1{
				ID:    "as-import-id",
				Name:  "imported-scope",
				Query: `resource_type = "mock-duck"`,
			}, nil).
			Times(1)

		ctx := t.Context()
		req := resource.ImportStateRequest{ID: "as-import-id"}

		stateType := getStateType()
		schema := r.getTestSchema(ctx)
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{Schema: schema, Raw: tftypes.NewValue(stateType, nil)},
		}

		r.ImportState(ctx, req, &importResp)

		require.False(t, importResp.Diagnostics.HasError())

		readReq := resource.ReadRequest{State: importResp.State}
		readResp := resource.ReadResponse{State: importResp.State}
		r.Read(ctx, readReq, &readResp)

		require.False(t, readResp.Diagnostics.HasError())
		var stateModel services.AccessScopeModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "as-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "imported-scope", stateModel.Name.ValueString())
		assert.Equal(t, `resource_type = "mock-duck"`, stateModel.Query.ValueString())
	})
}

func (r *AponoAccessScopeResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

func getStateType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":    tftypes.String,
			"name":  tftypes.String,
			"query": tftypes.String,
		},
	}
}

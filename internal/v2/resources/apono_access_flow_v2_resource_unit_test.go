package resources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoAccessFlowV2Resource(t *testing.T) {
	r := &AponoAccessFlowV2Resource{}

	t.Run("Create", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		mockResponse := testcommon.GenerateAccessFlowResponse()

		ctx := t.Context()

		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		model.ID = types.StringNull()

		mockInvoker.EXPECT().
			CreateAccessFlowV2(mock.Anything, mock.Anything).
			Return(mockResponse, nil)

		req := resource.CreateRequest{
			Plan: tfsdk.Plan{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.Plan.Set(ctx, model)
		require.False(t, diags.HasError(), "Error setting plan: %s", diags.Errors())

		resp := resource.CreateResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    req.Plan.Raw,
			},
		}

		r.Create(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Create returned error: %s", resp.Diagnostics.Errors())

		var state models.AccessFlowV2Model
		diags = resp.State.Get(ctx, &state)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		model.ID = state.ID

		assert.Equal(t, state, *model)
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		mockResponse := testcommon.GenerateAccessFlowResponse()
		ctx := t.Context()

		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		mockInvoker.EXPECT().
			GetAccessFlowV2(mock.Anything, mock.Anything).
			Return(mockResponse, nil)

		state := *model

		req := resource.ReadRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}
		diags := req.State.Set(ctx, state)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.ReadResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    req.State.Raw,
			},
		}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Read returned error: %s", resp.Diagnostics.Errors())

		var got models.AccessFlowV2Model
		diags = resp.State.Get(ctx, &got)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		assert.Equal(t, state, got)
	})

	t.Run("ReadNotFound", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()

		notFoundErr := &client.NotFoundError{}
		mockInvoker.EXPECT().
			GetAccessFlowV2(mock.Anything, mock.Anything).
			Return(nil, notFoundErr)

		mockResponse := testcommon.GenerateAccessFlowResponse()
		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")
		state := *model

		req := resource.ReadRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}
		diags := req.State.Set(ctx, state)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.ReadResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    req.State.Raw,
			},
		}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("Update", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		mockResponse := testcommon.GenerateAccessFlowResponse()
		ctx := t.Context()

		updatedResponse := testcommon.GenerateAccessFlowResponse()
		updatedResponse.Name = "updated-name"

		planModel, err := models.AccessFlowResponseToModel(ctx, *updatedResponse)
		require.NoError(t, err, "Failed to convert updated response to model: %s", err)

		stateModel, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model: %s", err)

		mockInvoker.EXPECT().
			UpdateAccessFlowV2(ctx, mock.Anything, mock.Anything).
			Return(updatedResponse, nil)

		req := resource.UpdateRequest{
			Plan: tfsdk.Plan{
				Schema: r.getTestSchema(ctx),
			},
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.Plan.Set(ctx, planModel)
		require.False(t, diags.HasError(), "Error setting plan: %s", diags.Errors())
		diags = req.State.Set(ctx, stateModel)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.UpdateResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    req.Plan.Raw,
			},
		}

		r.Update(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Update returned error: %s", resp.Diagnostics.Errors())

		var got models.AccessFlowV2Model
		diags = resp.State.Get(ctx, &got)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		assert.Equal(t, *planModel, got)
	})

	t.Run("Delete", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()

		mockResponse := testcommon.GenerateAccessFlowResponse()

		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model: %s", err)

		mockInvoker.EXPECT().
			DeleteAccessFlowV2(ctx, mock.Anything).
			Return(nil)

		req := resource.DeleteRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.State.Set(ctx, *model)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.DeleteResponse{}

		r.Delete(ctx, req, &resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("DeleteNotFound", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()
		notFoundErr := &client.NotFoundError{}

		mockResponse := testcommon.GenerateAccessFlowResponse()
		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model: %s", err)

		mockInvoker.EXPECT().
			DeleteAccessFlowV2(ctx, mock.Anything).
			Return(notFoundErr)

		req := resource.DeleteRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.State.Set(ctx, *model)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.DeleteResponse{}

		r.Delete(ctx, req, &resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()

		mockResponse := testcommon.GenerateAccessFlowResponse()
		model, err := models.AccessFlowResponseToModel(ctx, *mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model: %s", err)

		mockInvoker.EXPECT().
			GetAccessFlowV2(mock.Anything, mock.Anything).
			Return(mockResponse, nil)

		req := resource.ImportStateRequest{
			ID: model.ID.ValueString(),
		}

		resp := resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    tftypes.NewValue(r.getTestSchema(ctx).Type().TerraformType(ctx), nil),
			},
		}

		r.ImportState(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		readReq := resource.ReadRequest{State: resp.State}
		readResp := resource.ReadResponse{State: resp.State}
		r.Read(ctx, readReq, &readResp)

		var imported models.AccessFlowV2Model
		diags := readResp.State.Get(ctx, &imported)
		require.False(t, diags.HasError())
		assert.Equal(t, *model, imported)
	})
}

func (r *AponoAccessFlowV2Resource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

package resources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
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

func TestAponoResourceIntegrationResource(t *testing.T) {
	r := &AponoResourceIntegrationResource{}

	t.Run("Create", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		mockResponse.Category = common.ResourceCategory

		ctx := t.Context()

		model, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		model.ID = types.StringNull()

		mockInvoker.EXPECT().
			CreateIntegrationV4(mock.Anything, mock.Anything).
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

		var state models.ResourceIntegrationModel
		diags = resp.State.Get(ctx, &state)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		assert.Equal(t, mockResponse.ID, state.ID.ValueString())
		assert.Equal(t, mockResponse.Name, state.Name.ValueString())
		assert.Equal(t, mockResponse.Type, state.Type.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		mockResponse.Category = common.ResourceCategory
		ctx := t.Context()

		model, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: mockResponse.ID}).
			Return(mockResponse, nil)

		req := resource.ReadRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}
		diags := req.State.Set(ctx, model)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.ReadResponse{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
				Raw:    req.State.Raw,
			},
		}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "Read returned error: %s", resp.Diagnostics.Errors())

		var got models.ResourceIntegrationModel
		diags = resp.State.Get(ctx, &got)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		assert.Equal(t, model.ID.ValueString(), got.ID.ValueString())
		assert.Equal(t, model.Name.ValueString(), got.Name.ValueString())
		assert.Equal(t, model.Type.ValueString(), got.Type.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()
		notFoundErr := &client.NotFoundError{}

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		model, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: mockResponse.ID}).
			Return(nil, notFoundErr)

		req := resource.ReadRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}
		diags := req.State.Set(ctx, model)
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

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		mockResponse.Category = common.ResourceCategory
		ctx := t.Context()

		updatedResponse := testcommon.GenerateResourceIntegrationResponse()
		updatedResponse.Name = "updated-name"
		updatedResponse.Category = common.ResourceCategory
		updatedResponse.CustomAccessDetails.Value = "Updated access details"
		updatedResponse.ConnectedResourceTypes.Value = append(updatedResponse.ConnectedResourceTypes.Value, "role")

		stateModel, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		planModel, err := models.ResourceIntegrationToModel(ctx, updatedResponse)
		require.NoError(t, err, "Failed to convert updated response to model")

		mockInvoker.EXPECT().
			UpdateIntegrationV4(mock.Anything, mock.Anything, client.UpdateIntegrationV4Params{ID: mockResponse.ID}).
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

		var got models.ResourceIntegrationModel
		diags = resp.State.Get(ctx, &got)
		require.False(t, diags.HasError(), "Error getting state: %s", diags.Errors())

		assert.Equal(t, planModel.ID.ValueString(), got.ID.ValueString())
		assert.Equal(t, "updated-name", got.Name.ValueString())
		assert.Equal(t, "Updated access details", got.CustomAccessDetails.ValueString())

		var connectedTypes []string
		diags = got.ConnectedResourceTypes.ElementsAs(ctx, &connectedTypes, false)
		require.False(t, diags.HasError(), "Error getting connected resource types")
		assert.Contains(t, connectedTypes, "role")
	})

	t.Run("Delete", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		mockResponse.Category = common.ResourceCategory

		model, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		mockInvoker.EXPECT().
			DeleteIntegrationV4(mock.Anything, client.DeleteIntegrationV4Params{ID: mockResponse.ID}).
			Return(nil)

		req := resource.DeleteRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.State.Set(ctx, model)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.DeleteResponse{}

		r.Delete(ctx, req, &resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()
		notFoundErr := &client.NotFoundError{}

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		model, err := models.ResourceIntegrationToModel(ctx, mockResponse)
		require.NoError(t, err, "Failed to convert mock response to model")

		mockInvoker.EXPECT().
			DeleteIntegrationV4(mock.Anything, client.DeleteIntegrationV4Params{ID: mockResponse.ID}).
			Return(notFoundErr)

		req := resource.DeleteRequest{
			State: tfsdk.State{
				Schema: r.getTestSchema(ctx),
			},
		}

		diags := req.State.Set(ctx, model)
		require.False(t, diags.HasError(), "Error setting state: %s", diags.Errors())

		resp := resource.DeleteResponse{}

		r.Delete(ctx, req, &resp)

		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)
		r.client = mockInvoker

		ctx := t.Context()

		mockResponse := testcommon.GenerateResourceIntegrationResponse()
		mockResponse.Category = common.ResourceCategory

		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: mockResponse.ID}).
			Return(mockResponse, nil)

		req := resource.ImportStateRequest{
			ID: mockResponse.ID,
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

		var imported models.ResourceIntegrationModel
		diags := readResp.State.Get(ctx, &imported)
		require.False(t, diags.HasError())

		assert.Equal(t, mockResponse.ID, imported.ID.ValueString())
		assert.Equal(t, mockResponse.Name, imported.Name.ValueString())
		assert.Equal(t, mockResponse.Type, imported.Type.ValueString())
	})
}

func (r *AponoResourceIntegrationResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

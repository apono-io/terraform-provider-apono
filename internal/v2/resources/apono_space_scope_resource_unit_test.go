package resources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoSpaceScopeResource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoSpaceScopeResource{client: mockInvoker}

	getStateType := func() tftypes.Object {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":    tftypes.String,
				"name":  tftypes.String,
				"query": tftypes.String,
			},
		}
	}

	t.Run("Create", func(t *testing.T) {
		mockInvoker.EXPECT().
			CreateSpaceScopeV1(mock.Anything, mock.MatchedBy(func(req *client.UpsertSpaceScopeV1) bool {
				return req.Name == "Production AWS" && req.Query == `integration in ("aws-account") and resource_tag["environment"] = "production"`
			})).
			Return(&client.SpaceScopeV1{
				ID:    "ss-123456",
				Name:  "Production AWS",
				Query: `integration in ("aws-account") and resource_tag["environment"] = "production"`,
			}, nil).
			Once()

		ctx := t.Context()
		planType := getStateType()
		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, nil),
			"name":  tftypes.NewValue(tftypes.String, "Production AWS"),
			"query": tftypes.NewValue(tftypes.String, `integration in ("aws-account") and resource_tag["environment"] = "production"`),
		})

		schema := r.getTestSchema(ctx)
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(planType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}

		r.Create(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateVal models.SpaceScopeModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "ss-123456", stateVal.ID.ValueString())
		assert.Equal(t, "Production AWS", stateVal.Name.ValueString())
		assert.Equal(t, `integration in ("aws-account") and resource_tag["environment"] = "production"`, stateVal.Query.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetSpaceScopeV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceScopeV1Params) bool {
				return params.ID == "ss-123456"
			})).
			Return(&client.SpaceScopeV1{
				ID:    "ss-123456",
				Name:  "Production AWS",
				Query: `integration in ("aws-account") and resource_tag["environment"] = "production"`,
			}, nil).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-123456"),
			"name":  tftypes.NewValue(tftypes.String, "old-name"),
			"query": tftypes.NewValue(tftypes.String, `old-query`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateModel models.SpaceScopeModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "ss-123456", stateModel.ID.ValueString())
		assert.Equal(t, "Production AWS", stateModel.Name.ValueString())
		assert.Equal(t, `integration in ("aws-account") and resource_tag["environment"] = "production"`, stateModel.Query.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		notFoundErr := &client.NotFoundError{}
		mockInvoker.EXPECT().
			GetSpaceScopeV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceScopeV1Params) bool {
				return params.ID == "ss-not-found"
			})).
			Return(nil, notFoundErr).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-not-found"),
			"name":  tftypes.NewValue(tftypes.String, "test-scope"),
			"query": tftypes.NewValue(tftypes.String, `some-query`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("Update", func(t *testing.T) {
		mockInvoker.EXPECT().
			UpdateSpaceScopeV1(mock.Anything,
				mock.MatchedBy(func(req *client.UpsertSpaceScopeV1) bool {
					return req.Name == "Updated AWS" && req.Query == `integration in ("aws-account")`
				}),
				mock.MatchedBy(func(params client.UpdateSpaceScopeV1Params) bool {
					return params.ID == "ss-123456"
				}),
			).
			Return(&client.SpaceScopeV1{
				ID:    "ss-123456",
				Name:  "Updated AWS",
				Query: `integration in ("aws-account")`,
			}, nil).
			Once()

		ctx := t.Context()
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-123456"),
			"name":  tftypes.NewValue(tftypes.String, "Production AWS"),
			"query": tftypes.NewValue(tftypes.String, `old-query`),
		})
		planVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-123456"),
			"name":  tftypes.NewValue(tftypes.String, "Updated AWS"),
			"query": tftypes.NewValue(tftypes.String, `integration in ("aws-account")`),
		})

		schema := r.getTestSchema(ctx)
		req := resource.UpdateRequest{
			State: tfsdk.State{Schema: schema, Raw: stateVal},
			Plan:  tfsdk.Plan{Schema: schema, Raw: planVal},
		}
		resp := resource.UpdateResponse{
			State: tfsdk.State{Schema: schema, Raw: stateVal},
		}

		r.Update(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateModel models.SpaceScopeModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "ss-123456", stateModel.ID.ValueString())
		assert.Equal(t, "Updated AWS", stateModel.Name.ValueString())
		assert.Equal(t, `integration in ("aws-account")`, stateModel.Query.ValueString())
	})

	t.Run("Delete", func(t *testing.T) {
		mockInvoker.EXPECT().
			DeleteSpaceScopeV1(mock.Anything, mock.MatchedBy(func(params client.DeleteSpaceScopeV1Params) bool {
				return params.ID == "ss-123456"
			})).
			Return(nil).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-123456"),
			"name":  tftypes.NewValue(tftypes.String, "Production AWS"),
			"query": tftypes.NewValue(tftypes.String, `some-query`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{State: state}

		r.Delete(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		notFoundErr := &client.NotFoundError{}
		mockInvoker.EXPECT().
			DeleteSpaceScopeV1(mock.Anything, mock.MatchedBy(func(params client.DeleteSpaceScopeV1Params) bool {
				return params.ID == "ss-not-found"
			})).
			Return(notFoundErr).
			Once()

		ctx := t.Context()
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":    tftypes.NewValue(tftypes.String, "ss-not-found"),
			"name":  tftypes.NewValue(tftypes.String, "test-scope"),
			"query": tftypes.NewValue(tftypes.String, `some-query`),
		})

		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{State: state}

		r.Delete(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState_ByID", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetSpaceScopeV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceScopeV1Params) bool {
				return params.ID == "ss-import-id"
			})).
			Return(&client.SpaceScopeV1{
				ID:    "ss-import-id",
				Name:  "imported-scope",
				Query: `integration in ("aws-account")`,
			}, nil).
			Times(1)

		ctx := t.Context()
		req := resource.ImportStateRequest{ID: "ss-import-id"}

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
		var stateModel models.SpaceScopeModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "ss-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "imported-scope", stateModel.Name.ValueString())
		assert.Equal(t, `integration in ("aws-account")`, stateModel.Query.ValueString())
	})
}

func (r *AponoSpaceScopeResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

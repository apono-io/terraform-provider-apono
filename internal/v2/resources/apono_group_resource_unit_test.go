package resources

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/services"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoGroupResource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoGroupResource{client: mockInvoker}

	getStateType := func() tftypes.Type {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":      tftypes.String,
				"name":    tftypes.String,
				"members": tftypes.Set{ElementType: tftypes.String},
			},
		}
	}

	getPlanType := func() tftypes.Type {
		return getStateType()
	}

	t.Run("Create", func(t *testing.T) {
		planType := getPlanType()
		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "test-group"),
			"members": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "user1@example.com"),
				tftypes.NewValue(tftypes.String, "user2@example.com"),
			}),
		})

		mockInvoker.EXPECT().
			CreateGroupV1(mock.Anything, mock.MatchedBy(func(request *client.CreateGroupV1) bool {
				return request.Name == "test-group"
			})).
			Return(&client.GroupV1{
				ID:   "group-123456",
				Name: "test-group",
			}, nil).
			Once()

		schema := r.getTestSchema(t.Context())
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(planType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}

		r.Create(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateVal services.GroupModel
		diags := resp.State.Get(t.Context(), &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "group-123456", stateVal.ID.ValueString())
		assert.Equal(t, "test-group", stateVal.Name.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetGroupV1(mock.Anything, mock.MatchedBy(func(params client.GetGroupV1Params) bool {
				return params.ID == "group-123456"
			})).
			Return(&client.GroupV1{
				ID:   "group-123456",
				Name: "test-group",
			}, nil).
			Once()

		mockInvoker.EXPECT().
			ListGroupMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListGroupMembersV1Params) bool {
				return params.ID == "group-123456"
			})).
			Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
				Items: []client.GroupMemberV1{
					{Email: "user1@example.com"},
					{Email: "user2@example.com"},
				},
			}, nil).
			Once()

		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":      tftypes.NewValue(tftypes.String, "group-123456"),
			"name":    tftypes.NewValue(tftypes.String, "old-name"),
			"members": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{}),
		})

		ctx := t.Context()
		schema := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateModel services.GroupModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "group-123456", stateModel.ID.ValueString())
		assert.Equal(t, "test-group", stateModel.Name.ValueString())

		members := []string{}
		diags = stateModel.Members.ElementsAs(ctx, &members, false)
		require.False(t, diags.HasError())
		require.Equal(t, 2, len(members))
		assert.Contains(t, members, "user1@example.com")
		assert.Contains(t, members, "user2@example.com")
	})

	t.Run("Update", func(t *testing.T) {
		planType := getPlanType()
		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "group-123456"),
			"name": tftypes.NewValue(tftypes.String, "updated-group"),
			"members": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "user3@example.com"),
			}),
		})

		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "group-123456"),
			"name": tftypes.NewValue(tftypes.String, "test-group"),
			"members": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "user1@example.com"),
				tftypes.NewValue(tftypes.String, "user2@example.com"),
			}),
		})

		mockInvoker.EXPECT().
			UpdateGroupV1(mock.Anything, mock.MatchedBy(func(request *client.UpdateGroupV1) bool {
				return request.Name == "updated-group"
			}), mock.MatchedBy(func(params client.UpdateGroupV1Params) bool {
				return params.ID == "group-123456"
			})).
			Return(&client.GroupV1{
				ID:   "group-123456",
				Name: "updated-group",
			}, nil).
			Once()

		mockInvoker.EXPECT().
			UpdateGroupMembersV1(mock.Anything, mock.Anything, mock.MatchedBy(func(params client.UpdateGroupMembersV1Params) bool {
				return params.ID == "group-123456"
			})).
			Return(nil).
			Once()

		schema := r.getTestSchema(t.Context())
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.UpdateRequest{Plan: plan, State: state}
		resp := resource.UpdateResponse{State: state}

		r.Update(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var updatedState services.GroupModel
		diags := resp.State.Get(t.Context(), &updatedState)
		require.False(t, diags.HasError())

		assert.Equal(t, "group-123456", updatedState.ID.ValueString())
		assert.Equal(t, "updated-group", updatedState.Name.ValueString())
	})

	t.Run("Delete", func(t *testing.T) {
		stateType := getStateType()
		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":      tftypes.NewValue(tftypes.String, "group-123456"),
			"name":    tftypes.NewValue(tftypes.String, "test-group"),
			"members": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{}),
		})

		mockInvoker.EXPECT().
			DeleteGroupV1(mock.Anything, mock.MatchedBy(func(params client.DeleteGroupV1Params) bool {
				return params.ID == "group-123456"
			})).
			Return(nil).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{}

		r.Delete(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState_ByID", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetGroupV1(mock.Anything, mock.MatchedBy(func(params client.GetGroupV1Params) bool {
				return params.ID == "group-import-id"
			})).
			Return(&client.GroupV1{
				ID:   "group-import-id",
				Name: "imported-group",
			}, nil).
			Times(1)

		mockInvoker.EXPECT().
			ListGroupMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListGroupMembersV1Params) bool {
				return params.ID == "group-import-id"
			})).
			Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
				Items: []client.GroupMemberV1{
					{Email: "user1@example.com"},
					{Email: "user2@example.com"},
				},
			}, nil).
			Times(1)

		ctx := t.Context()
		req := resource.ImportStateRequest{ID: "group-import-id"}

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
		var stateModel services.GroupModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "group-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "imported-group", stateModel.Name.ValueString())

		members := []string{}
		diags = stateModel.Members.ElementsAs(ctx, &members, false)
		require.False(t, diags.HasError())
		require.Equal(t, 2, len(members))
		assert.Contains(t, members, "user1@example.com")
		assert.Contains(t, members, "user2@example.com")
	})
}

func (r *AponoGroupResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

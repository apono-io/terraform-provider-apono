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

var spaceStateType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"id":                     tftypes.String,
		"name":                   tftypes.String,
		"space_scope_references": tftypes.Set{ElementType: tftypes.String},
		"members": tftypes.Set{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"identity_reference": tftypes.String,
					"identity_type":      tftypes.String,
					"space_roles":        tftypes.Set{ElementType: tftypes.String},
				},
			},
		},
	},
}

func newSpaceMemberValue(identityRef string, roles []string) tftypes.Value {
	identityType := "user"

	memberObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"identity_reference": tftypes.String,
			"identity_type":      tftypes.String,
			"space_roles":        tftypes.Set{ElementType: tftypes.String},
		},
	}

	roleValues := make([]tftypes.Value, len(roles))
	for i, role := range roles {
		roleValues[i] = tftypes.NewValue(tftypes.String, role)
	}

	return tftypes.NewValue(memberObjType, map[string]tftypes.Value{
		"identity_reference": tftypes.NewValue(tftypes.String, identityRef),
		"identity_type":      tftypes.NewValue(tftypes.String, identityType),
		"space_roles":        tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, roleValues),
	})
}

func TestAponoSpaceResource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoSpaceResource{client: mockInvoker}

	t.Run("Create_WithMembers", func(t *testing.T) {
		mockInvoker.EXPECT().
			CreateSpaceV1(mock.Anything, mock.MatchedBy(func(req *client.CreateSpaceV1) bool {
				if req.Name != "Production" {
					return false
				}
				if len(req.SpaceScopeReferences) != 1 || req.SpaceScopeReferences[0] != "Production AWS" {
					return false
				}
				members, ok := req.Members.Get()
				if !ok || len(members) != 1 {
					return false
				}
				return members[0].IdentityReference == "admin@example.com" &&
					members[0].IdentityType == "user" &&
					len(members[0].SpaceRoles) == 1 &&
					members[0].SpaceRoles[0] == "SpaceOwner"
			})).
			Return(&client.SpaceV1{
				ID:   "space-123",
				Name: "Production",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Production AWS", Query: "some-query"},
				},
			}, nil).Once()

		ctx := t.Context()
		planVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("admin@example.com", []string{"SpaceOwner"}),
			}),
		})

		s := r.getTestSchema(ctx)
		plan := tfsdk.Plan{Schema: s, Raw: planVal}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(spaceStateType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}
		r.Create(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateVal models.SpaceModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())
		assert.Equal(t, "space-123", stateVal.ID.ValueString())
		assert.Equal(t, "Production", stateVal.Name.ValueString())
	})

	t.Run("Create_WithoutMembers", func(t *testing.T) {
		mockInvoker.EXPECT().
			CreateSpaceV1(mock.Anything, mock.MatchedBy(func(req *client.CreateSpaceV1) bool {
				return req.Name == "Staging" && !req.Members.Set
			})).
			Return(&client.SpaceV1{
				ID:   "space-456",
				Name: "Staging",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-2", Name: "Staging AWS", Query: "some-query"},
				},
			}, nil).Once()

		ctx := t.Context()
		planVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, nil),
			"name": tftypes.NewValue(tftypes.String, "Staging"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Staging AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		plan := tfsdk.Plan{Schema: s, Raw: planVal}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(spaceStateType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}
		r.Create(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateVal models.SpaceModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())
		assert.Equal(t, "space-456", stateVal.ID.ValueString())
		assert.True(t, stateVal.Members.IsNull())
	})

	t.Run("Read", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetSpaceV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceV1Params) bool {
				return params.ID == "space-123"
			})).
			Return(&client.SpaceV1{
				ID:   "space-123",
				Name: "Production",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Production AWS", Query: "q1"},
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-123"
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{
					{
						IdentityID:   "user-id-1",
						IdentityType: "user",
						SpaceRoles:   []string{"SpaceOwner"},
						Name:         "Admin User",
						Email:        client.NewOptNilString("admin@example.com"),
					},
				},
			}, nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "old-name"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "old-scope"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("admin@example.com", []string{"SpaceOwner"}),
			}),
		})

		s := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: s, Raw: stateVal}
		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}
		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateModel models.SpaceModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "space-123", stateModel.ID.ValueString())
		assert.Equal(t, "Production", stateModel.Name.ValueString())

		var scopeRefs []string
		diags = stateModel.SpaceScopeReferences.ElementsAs(ctx, &scopeRefs, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"Production AWS"}, scopeRefs)

		var members []models.SpaceMemberModel
		diags = stateModel.Members.ElementsAs(ctx, &members, false)
		require.False(t, diags.HasError())
		require.Len(t, members, 1)
		assert.Equal(t, "admin@example.com", members[0].IdentityReference.ValueString())
		assert.Equal(t, "user", members[0].IdentityType.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		notFoundErr := &client.NotFoundError{}
		mockInvoker.EXPECT().
			GetSpaceV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceV1Params) bool {
				return params.ID == "space-not-found"
			})).
			Return(nil, notFoundErr).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-not-found"),
			"name": tftypes.NewValue(tftypes.String, "test"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "scope"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: s, Raw: stateVal}
		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}
		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("Read_NullMembers_StaysNull", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetSpaceV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceV1Params) bool {
				return params.ID == "space-no-members"
			})).
			Return(&client.SpaceV1{
				ID:   "space-no-members",
				Name: "No Members",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Scope1", Query: "q"},
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-no-members"
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{},
			}, nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-no-members"),
			"name": tftypes.NewValue(tftypes.String, "No Members"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Scope1"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: s, Raw: stateVal}
		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}
		r.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateModel models.SpaceModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())
		assert.True(t, stateModel.Members.IsNull())
	})

	t.Run("Update_NameAndScopes", func(t *testing.T) {
		mockInvoker.EXPECT().
			UpdateSpaceV1(mock.Anything,
				mock.MatchedBy(func(req *client.UpdateSpaceV1) bool {
					return req.Name == "Updated Production" &&
						len(req.SpaceScopeReferences) == 2
				}),
				mock.MatchedBy(func(params client.UpdateSpaceV1Params) bool {
					return params.ID == "space-123"
				}),
			).
			Return(&client.SpaceV1{
				ID:   "space-123",
				Name: "Updated Production",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Production AWS", Query: "q1"},
					{ID: "ss-2", Name: "Production GCP", Query: "q2"},
				},
			}, nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})
		planVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Updated Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
				tftypes.NewValue(tftypes.String, "Production GCP"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		req := resource.UpdateRequest{
			State: tfsdk.State{Schema: s, Raw: stateVal},
			Plan:  tfsdk.Plan{Schema: s, Raw: planVal},
		}
		resp := resource.UpdateResponse{
			State: tfsdk.State{Schema: s, Raw: stateVal},
		}
		r.Update(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateModel models.SpaceModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "Updated Production", stateModel.Name.ValueString())
	})

	t.Run("Update_MembersOnly", func(t *testing.T) {
		mockInvoker.EXPECT().
			UpdateSpaceV1(mock.Anything, mock.Anything,
				mock.MatchedBy(func(params client.UpdateSpaceV1Params) bool {
					return params.ID == "space-123"
				}),
			).
			Return(&client.SpaceV1{
				ID:   "space-123",
				Name: "Production",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Production AWS", Query: "q1"},
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ReplaceSpaceMembersV1(mock.Anything,
				mock.MatchedBy(func(req *client.UpdateSpaceMembersV1) bool {
					return len(req.Members) == 1 &&
						req.Members[0].IdentityReference == "new-admin@example.com" &&
						req.Members[0].IdentityType == "user"
				}),
				mock.MatchedBy(func(params client.ReplaceSpaceMembersV1Params) bool {
					return params.ID == "space-123"
				}),
			).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{}, nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("admin@example.com", []string{"SpaceOwner"}),
			}),
		})
		planVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("new-admin@example.com", []string{"SpaceOwner"}),
			}),
		})

		s := r.getTestSchema(ctx)
		req := resource.UpdateRequest{
			State: tfsdk.State{Schema: s, Raw: stateVal},
			Plan:  tfsdk.Plan{Schema: s, Raw: planVal},
		}
		resp := resource.UpdateResponse{
			State: tfsdk.State{Schema: s, Raw: stateVal},
		}
		r.Update(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
	})

	t.Run("Update_Both", func(t *testing.T) {
		mockInvoker.EXPECT().
			UpdateSpaceV1(mock.Anything,
				mock.MatchedBy(func(req *client.UpdateSpaceV1) bool {
					return req.Name == "New Name"
				}),
				mock.MatchedBy(func(params client.UpdateSpaceV1Params) bool {
					return params.ID == "space-123"
				}),
			).
			Return(&client.SpaceV1{
				ID:   "space-123",
				Name: "New Name",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Production AWS", Query: "q1"},
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ReplaceSpaceMembersV1(mock.Anything, mock.Anything,
				mock.MatchedBy(func(params client.ReplaceSpaceMembersV1Params) bool {
					return params.ID == "space-123"
				}),
			).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{}, nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("admin@example.com", []string{"SpaceOwner"}),
			}),
		})
		planVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "New Name"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], []tftypes.Value{
				newSpaceMemberValue("new-admin@example.com", []string{"SpaceManager"}),
			}),
		})

		s := r.getTestSchema(ctx)
		req := resource.UpdateRequest{
			State: tfsdk.State{Schema: s, Raw: stateVal},
			Plan:  tfsdk.Plan{Schema: s, Raw: planVal},
		}
		resp := resource.UpdateResponse{
			State: tfsdk.State{Schema: s, Raw: stateVal},
		}
		r.Update(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics.Errors())
		var stateModel models.SpaceModel
		diags := resp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "New Name", stateModel.Name.ValueString())
	})

	t.Run("Delete", func(t *testing.T) {
		mockInvoker.EXPECT().
			DeleteSpaceV1(mock.Anything, mock.MatchedBy(func(params client.DeleteSpaceV1Params) bool {
				return params.ID == "space-123"
			})).
			Return(nil).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-123"),
			"name": tftypes.NewValue(tftypes.String, "Production"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "Production AWS"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: s, Raw: stateVal}
		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{State: state}
		r.Delete(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		notFoundErr := &client.NotFoundError{}
		mockInvoker.EXPECT().
			DeleteSpaceV1(mock.Anything, mock.MatchedBy(func(params client.DeleteSpaceV1Params) bool {
				return params.ID == "space-not-found"
			})).
			Return(notFoundErr).Once()

		ctx := t.Context()
		stateVal := tftypes.NewValue(spaceStateType, map[string]tftypes.Value{
			"id":   tftypes.NewValue(tftypes.String, "space-not-found"),
			"name": tftypes.NewValue(tftypes.String, "test"),
			"space_scope_references": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "scope"),
			}),
			"members": tftypes.NewValue(spaceStateType.AttributeTypes["members"], nil),
		})

		s := r.getTestSchema(ctx)
		state := tfsdk.State{Schema: s, Raw: stateVal}
		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{State: state}
		r.Delete(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState_ByID", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetSpaceV1(mock.Anything, mock.MatchedBy(func(params client.GetSpaceV1Params) bool {
				return params.ID == "space-import-id"
			})).
			Return(&client.SpaceV1{
				ID:   "space-import-id",
				Name: "Imported Space",
				SpaceScopes: []client.SpaceScopeV1{
					{ID: "ss-1", Name: "Scope1", Query: "q"},
				},
			}, nil).Times(1)

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-import-id"
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{
					{
						IdentityID:   "user-1",
						IdentityType: "user",
						SpaceRoles:   []string{"SpaceOwner"},
						Name:         "Admin",
						Email:        client.NewOptNilString("admin@example.com"),
					},
					{
						IdentityID:   "group-1",
						IdentityType: "group",
						SpaceRoles:   []string{"SpaceManager"},
						Name:         "engineers",
					},
				},
			}, nil).Times(1)

		ctx := t.Context()
		importReq := resource.ImportStateRequest{ID: "space-import-id"}
		s := r.getTestSchema(ctx)
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(spaceStateType, nil)},
		}
		r.ImportState(ctx, importReq, &importResp)
		require.False(t, importResp.Diagnostics.HasError())

		readReq := resource.ReadRequest{State: importResp.State}
		readResp := resource.ReadResponse{State: importResp.State}
		r.Read(ctx, readReq, &readResp)
		require.False(t, readResp.Diagnostics.HasError(), "diagnostics: %v", readResp.Diagnostics.Errors())

		var stateModel models.SpaceModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())
		assert.Equal(t, "space-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "Imported Space", stateModel.Name.ValueString())

		var members []models.SpaceMemberModel
		diags = stateModel.Members.ElementsAs(ctx, &members, false)
		require.False(t, diags.HasError())
		require.Len(t, members, 2)
	})
}

func (r *AponoSpaceResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListSpaceMembers(t *testing.T) {
	t.Run("single_page", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-123"
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{
					{IdentityID: "user-1", IdentityType: "user", Name: "User One", Email: client.NewOptNilString("user1@example.com"), SpaceRoles: []string{"SpaceOwner"}},
				},
			}, nil).Once()

		result, err := ListSpaceMembers(context.Background(), mockInvoker, "space-123")
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "user-1", result[0].IdentityID)
	})

	t.Run("multiple_pages", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-123" && !params.PageToken.Set
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{
					{IdentityID: "user-1", IdentityType: "user", Name: "User One"},
				},
				Pagination: client.PublicApiPaginationInfoModel{
					NextPageToken: client.NewOptNilString("page2"),
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceMembersV1Params) bool {
				return params.ID == "space-123" && params.PageToken.Value == "page2"
			})).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{
					{IdentityID: "group-1", IdentityType: "group", Name: "Engineers"},
				},
			}, nil).Once()

		result, err := ListSpaceMembers(context.Background(), mockInvoker, "space-123")
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, "user-1", result[0].IdentityID)
		assert.Equal(t, "group-1", result[1].IdentityID)
	})

	t.Run("empty_results", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.Anything).
			Return(&client.PublicApiListResponseSpaceMemberPublicV1Model{
				Items: []client.SpaceMemberV1{},
			}, nil).Once()

		result, err := ListSpaceMembers(context.Background(), mockInvoker, "space-123")
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("api_error", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpaceMembersV1(mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("api error")).Once()

		result, err := ListSpaceMembers(context.Background(), mockInvoker, "space-123")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestListSpaces(t *testing.T) {
	t.Run("single_page", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return !params.Name.IsSet() && !params.PageToken.IsSet()
			})).
			Return(&client.PublicApiListResponseSpacePublicV1Model{
				Items: []client.SpaceV1{
					{ID: "space-1", Name: "Production", SpaceScopes: []client.SpaceScopeV1{{ID: "ss-1", Name: "Scope1"}}},
				},
			}, nil).Once()

		result, err := ListSpaces(context.Background(), mockInvoker, "")
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "space-1", result[0].ID)
	})

	t.Run("with_name_filter", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return params.Name.Value == "Production*"
			})).
			Return(&client.PublicApiListResponseSpacePublicV1Model{
				Items: []client.SpaceV1{
					{ID: "space-1", Name: "Production", SpaceScopes: []client.SpaceScopeV1{{ID: "ss-1", Name: "Scope1"}}},
				},
			}, nil).Once()

		result, err := ListSpaces(context.Background(), mockInvoker, "Production*")
		require.NoError(t, err)
		require.Len(t, result, 1)
	})

	t.Run("multiple_pages", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return !params.PageToken.IsSet()
			})).
			Return(&client.PublicApiListResponseSpacePublicV1Model{
				Items: []client.SpaceV1{
					{ID: "space-2", Name: "Staging"},
				},
				Pagination: client.PublicApiPaginationInfoModel{
					NextPageToken: client.NewOptNilString("page2"),
				},
			}, nil).Once()

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return params.PageToken.Value == "page2"
			})).
			Return(&client.PublicApiListResponseSpacePublicV1Model{
				Items: []client.SpaceV1{
					{ID: "space-1", Name: "Production"},
				},
			}, nil).Once()

		result, err := ListSpaces(context.Background(), mockInvoker, "")
		require.NoError(t, err)
		require.Len(t, result, 2)
		// Sorted by ID
		assert.Equal(t, "space-1", result[0].ID)
		assert.Equal(t, "space-2", result[1].ID)
	})

	t.Run("api_error", func(t *testing.T) {
		mockInvoker := mocks.NewInvoker(t)

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.Anything).
			Return(nil, fmt.Errorf("api error")).Once()

		result, err := ListSpaces(context.Background(), mockInvoker, "")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

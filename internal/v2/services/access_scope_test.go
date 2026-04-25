package services

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAccessScopesByName(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name           string
		scopeName      string
		setupMock      func(*mocks.Invoker)
		expectedScopes []client.AccessScopeV1
		expectError    bool
	}{
		{
			name:      "single scope found",
			scopeName: "test-scope",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListAccessScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("test-scope")
				params.Name = nameParam

				m.On("ListAccessScopesV1", ctx, params).Return(&client.PublicApiListResponseAccessScopePublicV1Model{
					Items: []client.AccessScopeV1{
						{ID: "2", Name: "other-scope"},
						{ID: "1", Name: "test-scope"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.AccessScopeV1{
				{ID: "1", Name: "test-scope"},
				{ID: "2", Name: "other-scope"},
			},
		},
		{
			name:      "multiple pages",
			scopeName: "paginated-scope",
			setupMock: func(m *mocks.Invoker) {
				// First request with no page token
				firstParams := client.ListAccessScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("paginated-scope")
				firstParams.Name = nameParam

				nextToken := client.OptNilString{}
				nextToken.SetTo("next-page")

				m.On("ListAccessScopesV1", ctx, firstParams).Return(&client.PublicApiListResponseAccessScopePublicV1Model{
					Items: []client.AccessScopeV1{
						{ID: "1", Name: "paginated-scope"},
						{ID: "3", Name: "other-scope-page1"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: nextToken,
					},
				}, nil)

				// Second request with page token
				secondParams := client.ListAccessScopesV1Params{}
				pageToken := client.OptNilString{}
				pageToken.SetTo("next-page")
				secondParams.PageToken = pageToken

				m.On("ListAccessScopesV1", ctx, secondParams).Return(&client.PublicApiListResponseAccessScopePublicV1Model{
					Items: []client.AccessScopeV1{
						{ID: "2", Name: "paginated-scope2"},
						{ID: "4", Name: "other-scope-page2"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.AccessScopeV1{
				{ID: "1", Name: "paginated-scope"},
				{ID: "2", Name: "paginated-scope2"},
				{ID: "3", Name: "other-scope-page1"},
				{ID: "4", Name: "other-scope-page2"},
			},
		},
		{
			name:      "no scopes found",
			scopeName: "non-existent-scope",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListAccessScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("non-existent-scope")
				params.Name = nameParam

				m.On("ListAccessScopesV1", ctx, params).Return(&client.PublicApiListResponseAccessScopePublicV1Model{
					Items: []client.AccessScopeV1{},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.AccessScopeV1{},
		},
		{
			name:      "api error",
			scopeName: "error-scope",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListAccessScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("error-scope")
				params.Name = nameParam

				m.On("ListAccessScopesV1", ctx, params).Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			scopes, err := ListAccessScopesByName(ctx, mockClient, tc.scopeName, nil)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedScopes, scopes)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestListAccessScopesByNameWithSpaceReferences(t *testing.T) {
	ctx := t.Context()

	mockClient := new(mocks.Invoker)
	expectedParams := client.ListAccessScopesV1Params{}
	expectedParams.SpaceReferences.SetTo([]string{"prod", "staging"})
	mockClient.On("ListAccessScopesV1", ctx, expectedParams).Return(&client.PublicApiListResponseAccessScopePublicV1Model{
		Items: []client.AccessScopeV1{{ID: "s1"}},
		Pagination: client.PublicApiPaginationInfoModel{
			NextPageToken: client.NewOptNilString(""),
		},
	}, nil)

	scopes, err := ListAccessScopesByName(ctx, mockClient, "", []string{"prod", "staging"})

	assert.NoError(t, err)
	assert.Len(t, scopes, 1)
	mockClient.AssertExpectations(t)
}

func TestAccessScopeToModel(t *testing.T) {
	t.Run("WithSpace", func(t *testing.T) {
		apiScope := &client.AccessScopeV1{
			ID:    "scope-123",
			Name:  "prod-scope",
			Query: `resource_type = "mock-duck"`,
		}
		apiScope.Description.SetTo("a description")
		apiScope.Space.SetTo(client.SpaceReferenceV1{SpaceID: "space-123", SpaceName: "prod-space"})

		model, diags := AccessScopeToModel(apiScope)
		require.False(t, diags.HasError())
		require.NotNil(t, model)

		assert.Equal(t, "scope-123", model.ID.ValueString())
		assert.Equal(t, "prod-scope", model.Name.ValueString())
		assert.Equal(t, `resource_type = "mock-duck"`, model.Query.ValueString())
		assert.Equal(t, "a description", model.Description.ValueString())
		assert.Equal(t, "prod-space", model.SpaceReference.ValueString())
		require.False(t, model.Space.IsNull())
		spaceAttrs := model.Space.Attributes()
		spaceID, ok := spaceAttrs["space_id"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "space-123", spaceID.ValueString())
		spaceName, ok := spaceAttrs["space_name"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "prod-space", spaceName.ValueString())
	})

	t.Run("WithoutSpace", func(t *testing.T) {
		apiScope := &client.AccessScopeV1{
			ID:    "scope-456",
			Name:  "no-space-scope",
			Query: `resource_type = "mock-duck"`,
		}

		model, diags := AccessScopeToModel(apiScope)
		require.False(t, diags.HasError())
		require.NotNil(t, model)

		assert.Equal(t, "scope-456", model.ID.ValueString())
		assert.True(t, model.Description.IsNull())
		assert.True(t, model.SpaceReference.IsNull())
		assert.True(t, model.Space.IsNull())
	})
}

func TestAccessScopeToDataSourceItemModel(t *testing.T) {
	t.Run("WithSpace", func(t *testing.T) {
		apiScope := &client.AccessScopeV1{
			ID:    "scope-123",
			Name:  "prod-scope",
			Query: `resource_type = "mock-duck"`,
		}
		apiScope.Description.SetTo("a description")
		apiScope.Space.SetTo(client.SpaceReferenceV1{SpaceID: "space-123", SpaceName: "prod-space"})

		model, diags := AccessScopeToDataSourceItemModel(apiScope)
		require.False(t, diags.HasError())
		require.NotNil(t, model)

		assert.Equal(t, "scope-123", model.ID.ValueString())
		assert.Equal(t, "prod-scope", model.Name.ValueString())
		assert.Equal(t, `resource_type = "mock-duck"`, model.Query.ValueString())
		assert.Equal(t, "a description", model.Description.ValueString())
		require.False(t, model.Space.IsNull())
		spaceAttrs := model.Space.Attributes()
		spaceID, ok := spaceAttrs["space_id"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "space-123", spaceID.ValueString())
		spaceName, ok := spaceAttrs["space_name"].(types.String)
		require.True(t, ok)
		assert.Equal(t, "prod-space", spaceName.ValueString())
	})

	t.Run("WithoutSpace", func(t *testing.T) {
		apiScope := &client.AccessScopeV1{
			ID:    "scope-456",
			Name:  "no-space-scope",
			Query: `resource_type = "mock-duck"`,
		}

		model, diags := AccessScopeToDataSourceItemModel(apiScope)
		require.False(t, diags.HasError())
		require.NotNil(t, model)

		assert.Equal(t, "scope-456", model.ID.ValueString())
		assert.True(t, model.Description.IsNull())
		assert.True(t, model.Space.IsNull())
	})
}

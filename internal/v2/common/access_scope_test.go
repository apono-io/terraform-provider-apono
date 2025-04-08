package common

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
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
						{ID: "1", Name: "test-scope"},
						{ID: "2", Name: "other-scope"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.AccessScopeV1{
				{ID: "2", Name: "other-scope"},
				{ID: "1", Name: "test-scope"},
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
				secondParams.Name = nameParam
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
				{ID: "3", Name: "other-scope-page1"},
				{ID: "4", Name: "other-scope-page2"},
				{ID: "1", Name: "paginated-scope"},
				{ID: "2", Name: "paginated-scope2"},
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

			scopes, err := ListAccessScopesByName(ctx, mockClient, tc.scopeName)

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

package services

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListSpaceScopes(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name           string
		scopeName      string
		setupMock      func(*mocks.Invoker)
		expectedScopes []client.SpaceScopeV1
		expectError    bool
	}{
		{
			name:      "single scope found",
			scopeName: "Production AWS",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListSpaceScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("Production AWS")
				params.Name = nameParam

				m.On("ListSpaceScopesV1", ctx, params).Return(&client.PublicApiListResponseSpaceScopePublicV1Model{
					Items: []client.SpaceScopeV1{
						{ID: "2", Name: "Staging AWS", Query: `integration in ("aws-account")`},
						{ID: "1", Name: "Production AWS", Query: `integration in ("aws-account") and resource_tag["environment"] = "production"`},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.SpaceScopeV1{
				{ID: "1", Name: "Production AWS", Query: `integration in ("aws-account") and resource_tag["environment"] = "production"`},
				{ID: "2", Name: "Staging AWS", Query: `integration in ("aws-account")`},
			},
		},
		{
			name:      "multiple pages",
			scopeName: "*AWS*",
			setupMock: func(m *mocks.Invoker) {
				firstParams := client.ListSpaceScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("*AWS*")
				firstParams.Name = nameParam

				nextToken := client.OptNilString{}
				nextToken.SetTo("next-page")

				m.On("ListSpaceScopesV1", ctx, firstParams).Return(&client.PublicApiListResponseSpaceScopePublicV1Model{
					Items: []client.SpaceScopeV1{
						{ID: "1", Name: "Production AWS", Query: `query-1`},
						{ID: "3", Name: "Dev AWS", Query: `query-3`},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: nextToken,
					},
				}, nil)

				secondParams := client.ListSpaceScopesV1Params{}
				pageToken := client.OptNilString{}
				pageToken.SetTo("next-page")
				secondParams.PageToken = pageToken

				m.On("ListSpaceScopesV1", ctx, secondParams).Return(&client.PublicApiListResponseSpaceScopePublicV1Model{
					Items: []client.SpaceScopeV1{
						{ID: "2", Name: "Staging AWS", Query: `query-2`},
						{ID: "4", Name: "Test AWS", Query: `query-4`},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.SpaceScopeV1{
				{ID: "1", Name: "Production AWS", Query: `query-1`},
				{ID: "2", Name: "Staging AWS", Query: `query-2`},
				{ID: "3", Name: "Dev AWS", Query: `query-3`},
				{ID: "4", Name: "Test AWS", Query: `query-4`},
			},
		},
		{
			name:      "no scopes found",
			scopeName: "non-existent",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListSpaceScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("non-existent")
				params.Name = nameParam

				m.On("ListSpaceScopesV1", ctx, params).Return(&client.PublicApiListResponseSpaceScopePublicV1Model{
					Items: []client.SpaceScopeV1{},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.SpaceScopeV1{},
		},
		{
			name:      "empty name returns all",
			scopeName: "",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListSpaceScopesV1Params{}

				m.On("ListSpaceScopesV1", ctx, params).Return(&client.PublicApiListResponseSpaceScopePublicV1Model{
					Items: []client.SpaceScopeV1{
						{ID: "1", Name: "Scope A", Query: `query-a`},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedScopes: []client.SpaceScopeV1{
				{ID: "1", Name: "Scope A", Query: `query-a`},
			},
		},
		{
			name:      "api error",
			scopeName: "error-scope",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListSpaceScopesV1Params{}
				nameParam := client.OptNilString{}
				nameParam.SetTo("error-scope")
				params.Name = nameParam

				m.On("ListSpaceScopesV1", ctx, params).Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			scopes, err := ListSpaceScopes(ctx, mockClient, tc.scopeName)

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

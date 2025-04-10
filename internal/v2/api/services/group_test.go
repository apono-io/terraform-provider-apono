package services

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListGroupMembers(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name            string
		groupID         string
		setupMock       func(*mocks.Invoker)
		expectedMembers []client.GroupMemberV1
		expectError     bool
	}{
		{
			name:    "single page",
			groupID: "test-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupMembersV1Params{
					ID: "test-group",
				}

				m.On("ListGroupMembersV1", ctx, params).Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
					Items: []client.GroupMemberV1{
						{Email: "user1@example.com"},
						{Email: "user2@example.com"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedMembers: []client.GroupMemberV1{
				{Email: "user1@example.com"},
				{Email: "user2@example.com"},
			},
		},
		{
			name:    "multiple pages",
			groupID: "paginated-group",
			setupMock: func(m *mocks.Invoker) {
				// First request with no page token
				firstParams := client.ListGroupMembersV1Params{
					ID: "paginated-group",
				}

				nextToken := client.OptNilString{}
				nextToken.SetTo("next-page")

				m.On("ListGroupMembersV1", ctx, firstParams).Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
					Items: []client.GroupMemberV1{
						{Email: "user1@example.com"},
						{Email: "user2@example.com"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: nextToken,
					},
				}, nil)

				// Second request with page token
				secondParams := client.ListGroupMembersV1Params{
					ID: "paginated-group",
				}
				pageToken := client.OptNilString{}
				pageToken.SetTo("next-page")
				secondParams.PageToken = pageToken

				m.On("ListGroupMembersV1", ctx, secondParams).Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
					Items: []client.GroupMemberV1{
						{Email: "user3@example.com"},
						{Email: "user4@example.com"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedMembers: []client.GroupMemberV1{
				{Email: "user1@example.com"},
				{Email: "user2@example.com"},
				{Email: "user3@example.com"},
				{Email: "user4@example.com"},
			},
		},
		{
			name:    "no members found",
			groupID: "empty-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupMembersV1Params{
					ID: "empty-group",
				}

				m.On("ListGroupMembersV1", ctx, params).Return(&client.PublicApiListResponseGroupMemberPublicV1Model{
					Items:      []client.GroupMemberV1{},
					Pagination: client.PublicApiPaginationInfoModel{NextPageToken: client.OptNilString{}},
				}, nil)
			},
			expectedMembers: []client.GroupMemberV1{},
		},
		{
			name:    "api error",
			groupID: "error-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupMembersV1Params{
					ID: "error-group",
				}

				m.On("ListGroupMembersV1", ctx, params).Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			members, err := ListGroupMembers(ctx, mockClient, tc.groupID)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMembers, members)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

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

func TestListGroups(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name           string
		groupName      string
		setupMock      func(*mocks.Invoker)
		expectedGroups []client.GroupV1
		expectError    bool
	}{
		{
			name:      "single page",
			groupName: "test-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupsV1Params{}
				params.Name.SetTo("test-group")

				m.On("ListGroupsV1", ctx, params).Return(&client.PublicApiListResponseGroupPublicV1Model{
					Items: []client.GroupV1{
						{Name: "test-group-1"},
						{Name: "test-group-2"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedGroups: []client.GroupV1{
				{Name: "test-group-1"},
				{Name: "test-group-2"},
			},
		},
		{
			name:      "multiple pages",
			groupName: "paginated-group",
			setupMock: func(m *mocks.Invoker) {
				// First request with no page token
				firstParams := client.ListGroupsV1Params{}
				firstParams.Name.SetTo("paginated-group")

				nextToken := client.OptNilString{}
				nextToken.SetTo("next-page")

				m.On("ListGroupsV1", ctx, firstParams).Return(&client.PublicApiListResponseGroupPublicV1Model{
					Items: []client.GroupV1{
						{Name: "paginated-group-1"},
						{Name: "paginated-group-2"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: nextToken,
					},
				}, nil)

				// Second request with page token
				secondParams := client.ListGroupsV1Params{}
				pageToken := client.OptNilString{}
				pageToken.SetTo("next-page")
				secondParams.PageToken = pageToken

				m.On("ListGroupsV1", ctx, secondParams).Return(&client.PublicApiListResponseGroupPublicV1Model{
					Items: []client.GroupV1{
						{Name: "paginated-group-3"},
						{Name: "paginated-group-4"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expectedGroups: []client.GroupV1{
				{Name: "paginated-group-1"},
				{Name: "paginated-group-2"},
				{Name: "paginated-group-3"},
				{Name: "paginated-group-4"},
			},
		},
		{
			name:      "no groups found",
			groupName: "empty-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupsV1Params{}
				params.Name.SetTo("empty-group")

				m.On("ListGroupsV1", ctx, params).Return(&client.PublicApiListResponseGroupPublicV1Model{
					Items:      []client.GroupV1{},
					Pagination: client.PublicApiPaginationInfoModel{NextPageToken: client.OptNilString{}},
				}, nil)
			},
			expectedGroups: []client.GroupV1{},
		},
		{
			name:      "api error",
			groupName: "error-group",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListGroupsV1Params{}
				params.Name.SetTo("error-group")

				m.On("ListGroupsV1", ctx, params).Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			groups, err := ListGroups(ctx, mockClient, tc.groupName)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedGroups, groups)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFilterGroupsBySourceIntegration(t *testing.T) {
	tests := []struct {
		name              string
		groups            []client.GroupV1
		sourceIntegration string
		expectedGroups    []client.GroupV1
	}{
		{
			name:              "no filter",
			groups:            []client.GroupV1{{Name: "group1"}, {Name: "group2"}},
			sourceIntegration: "",
			expectedGroups:    []client.GroupV1{{Name: "group1"}, {Name: "group2"}},
		},
		{
			name: "filter by source integration ID",
			groups: []client.GroupV1{
				{Name: "group1", SourceIntegrationID: client.OptNilString{Value: "source1", Set: true}},
				{Name: "group2", SourceIntegrationID: client.OptNilString{Value: "source2", Set: true}},
				{Name: "group3"},
			},
			sourceIntegration: "source1",
			expectedGroups: []client.GroupV1{
				{Name: "group1", SourceIntegrationID: client.OptNilString{Value: "source1", Set: true}},
			},
		},
		{
			name: "filter by source integration name",
			groups: []client.GroupV1{
				{Name: "group1", SourceIntegrationName: client.OptNilString{Value: "source1", Set: true}},
				{Name: "group2", SourceIntegrationName: client.OptNilString{Value: "source2", Set: true}},
				{Name: "group3"},
			},
			sourceIntegration: "source1",
			expectedGroups: []client.GroupV1{
				{Name: "group1", SourceIntegrationName: client.OptNilString{Value: "source1", Set: true}},
			},
		},
		{
			name: "filter by source integration ID and name",
			groups: []client.GroupV1{
				{Name: "group1", SourceIntegrationID: client.OptNilString{Value: "source1", Set: true}, SourceIntegrationName: client.OptNilString{Value: "source1_name", Set: true}},
				{Name: "group2", SourceIntegrationID: client.OptNilString{Value: "source2", Set: true}, SourceIntegrationName: client.OptNilString{Value: "source2_name", Set: true}},
			},
			sourceIntegration: "source2",
			expectedGroups: []client.GroupV1{
				{Name: "group2", SourceIntegrationID: client.OptNilString{Value: "source2", Set: true}, SourceIntegrationName: client.OptNilString{Value: "source2_name", Set: true}},
			},
		},
		{
			name:              "no matching groups",
			groups:            []client.GroupV1{{Name: "group1"}},
			sourceIntegration: "source1",
			expectedGroups:    []client.GroupV1{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filteredGroups := FilterGroupsBySourceIntegration(tc.groups, tc.sourceIntegration)
			assert.Equal(t, tc.expectedGroups, filteredGroups)
		})
	}
}

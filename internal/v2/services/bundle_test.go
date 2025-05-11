package services

import (
	"errors"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListBundles(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name        string
		bundleName  string
		setupMock   func(*mocks.Invoker)
		expected    []client.BundlePublicV2Model
		expectError bool
	}{
		{
			name:       "single page",
			bundleName: "",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListBundlesV2Params{}
				m.On("ListBundlesV2", ctx, params).Return(&client.PublicApiListResponseBundlePublicV2Model{
					Items: []client.BundlePublicV2Model{
						{ID: "b2"},
						{ID: "b1"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.NewOptNilString(""),
					},
				}, nil)
			},
			expected: []client.BundlePublicV2Model{
				{ID: "b1"},
				{ID: "b2"},
			},
		},
		{
			name:       "multiple pages",
			bundleName: "",
			setupMock: func(m *mocks.Invoker) {
				firstParams := client.ListBundlesV2Params{}
				nextToken := client.NewOptNilString("next")
				m.On("ListBundlesV2", ctx, firstParams).Return(&client.PublicApiListResponseBundlePublicV2Model{
					Items: []client.BundlePublicV2Model{
						{ID: "b3"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: nextToken,
					},
				}, nil)

				secondParams := client.ListBundlesV2Params{}
				secondParams.PageToken = nextToken
				m.On("ListBundlesV2", ctx, secondParams).Return(&client.PublicApiListResponseBundlePublicV2Model{
					Items: []client.BundlePublicV2Model{
						{ID: "b2"},
						{ID: "b1"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.NewOptNilString(""),
					},
				}, nil)
			},
			expected: []client.BundlePublicV2Model{
				{ID: "b1"},
				{ID: "b2"},
				{ID: "b3"},
			},
		},
		{
			name:       "with name param",
			bundleName: "bundle-name",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListBundlesV2Params{}
				params.Name.SetTo("bundle-name")
				m.On("ListBundlesV2", ctx, params).Return(&client.PublicApiListResponseBundlePublicV2Model{
					Items: []client.BundlePublicV2Model{
						{ID: "b1"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.NewOptNilString(""),
					},
				}, nil)
			},
			expected: []client.BundlePublicV2Model{
				{ID: "b1"},
			},
		},
		{
			name:       "api error",
			bundleName: "",
			setupMock: func(m *mocks.Invoker) {
				params := client.ListBundlesV2Params{}
				m.On("ListBundlesV2", ctx, params).Return(nil, errors.New("api error"))
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			bundles, err := ListBundles(ctx, mockClient, tc.bundleName)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, bundles)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

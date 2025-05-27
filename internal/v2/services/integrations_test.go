package services

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/go-faster/jx"
	"github.com/stretchr/testify/assert"
)

func TestListIntegrations(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name            string
		typeName        string
		integrationName string
		connectorID     string
		categories      []string
		setupMock       func(*mocks.Invoker)
		expected        []client.IntegrationV4
		expectError     bool
	}{
		{
			name:            "single page integrations",
			typeName:        "postgres",
			integrationName: "test-integration",
			connectorID:     "conn-123",
			categories:      []string{"database"},
			setupMock: func(m *mocks.Invoker) {
				params := client.ListIntegrationsV4Params{}
				params.Name.SetTo("test-integration")
				params.Type.SetTo([]string{"postgres"})
				params.ConnectorID.SetTo([]string{"conn-123"})
				params.Category.SetTo([]string{"database"})

				m.On("ListIntegrationsV4", ctx, params).Return(&client.PublicApiListResponseIntegrationPublicV4Model{
					Items: []client.IntegrationV4{
						{
							ID:       "integration-id",
							Name:     "test-integration",
							Type:     "postgres",
							Category: "database",
							Status:   "connected",
							IntegrationConfig: map[string]jx.Raw{
								"host":     jx.Raw("\"localhost\""),
								"port":     jx.Raw("\"5432\""),
								"database": jx.Raw("\"postgres\""),
							},
						},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expected: []client.IntegrationV4{
				{
					ID:       "integration-id",
					Name:     "test-integration",
					Type:     "postgres",
					Category: "database",
					Status:   "connected",
					IntegrationConfig: map[string]jx.Raw{
						"host":     jx.Raw("\"localhost\""),
						"port":     jx.Raw("\"5432\""),
						"database": jx.Raw("\"postgres\""),
					},
				},
			},
		},
		{
			name:            "multiple pages",
			typeName:        "",
			integrationName: "",
			connectorID:     "",
			categories:      nil,
			setupMock: func(m *mocks.Invoker) {
				// First page
				firstParams := client.ListIntegrationsV4Params{}
				m.On("ListIntegrationsV4", ctx, firstParams).Return(&client.PublicApiListResponseIntegrationPublicV4Model{
					Items: []client.IntegrationV4{
						{ID: "1", Name: "integration1", Type: "mysql", Category: "database", Status: "connected"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: func() client.OptNilString { var s client.OptNilString; s.SetTo("next-page"); return s }(),
					},
				}, nil)

				// Second page
				secondParams := client.ListIntegrationsV4Params{}
				secondParams.PageToken.SetTo("next-page")
				m.On("ListIntegrationsV4", ctx, secondParams).Return(&client.PublicApiListResponseIntegrationPublicV4Model{
					Items: []client.IntegrationV4{
						{ID: "2", Name: "integration2", Type: "postgres", Category: "database", Status: "connected"},
					},
					Pagination: client.PublicApiPaginationInfoModel{
						NextPageToken: client.OptNilString{},
					},
				}, nil)
			},
			expected: []client.IntegrationV4{
				{ID: "1", Name: "integration1", Type: "mysql", Category: "database", Status: "connected"},
				{ID: "2", Name: "integration2", Type: "postgres", Category: "database", Status: "connected"},
			},
		},
		{
			name:            "api error",
			typeName:        "",
			integrationName: "",
			connectorID:     "",
			categories:      nil,
			setupMock: func(m *mocks.Invoker) {
				params := client.ListIntegrationsV4Params{}
				m.On("ListIntegrationsV4", ctx, params).Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(mocks.Invoker)
			tc.setupMock(mockClient)

			integrations, err := ListIntegrations(ctx, mockClient, tc.typeName, tc.integrationName, tc.connectorID, tc.categories)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, integrations)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

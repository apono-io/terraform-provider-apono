package api_test

import (
	"context"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientMock(t *testing.T) {
	// Create a new mock instance
	mockClient := mocks.NewInvoker(t)

	// Set expectations
	mockClient.On("ListIntegrationsV2", mock.Anything).Return(&client.PaginatedResponseIntegrationModel{
		Data: []client.Integration{
			{
				ID:   "1",
				Name: "Test Integration",
				Type: "test-type",
			},
		},
	}, nil)

	// Call the method
	result, err := mockClient.ListIntegrationsV2(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, "Test Integration", result.Data[0].Name)

	// Verify that expectations were met
	mockClient.AssertExpectations(t)
}

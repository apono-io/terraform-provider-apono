package testcommon

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// Helper function to get test users from the API.
func GetUsers(t *testing.T) ([]client.UserModel, error) {
	// Create a client using the provider's configuration
	c := GetTestClient(t)

	// Add timeout to context
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	// Get users from the API
	resp, err := c.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %v", err)
	}

	// Filter for active users only
	var activeUsers []client.UserModel
	for _, user := range resp.Data {
		if user.Active {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, nil
}

// Helper function to list connectors and return the first one.
func GetFirstConnectorV3(t *testing.T) (*client.ConnectorV3, error) {
	c := GetTestClient(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	params := client.ListConnectorsV3Params{}
	resp, err := c.ListConnectorsV3(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error listing connectors: %v", err)
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no connectors found")
	}
	return &resp.Items[0], nil
}

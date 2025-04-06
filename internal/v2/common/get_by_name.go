package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// GetAccessScopeByName retrieves an access scope by its name with pagination support.
// Returns the access scope if found, or an error if not found or an API error occurs.
func GetAccessScopeByName(ctx context.Context, apiClient client.Invoker, name string) (*client.AccessScopeV1, error) {
	var pageToken string

	for {
		params := client.ListAccessScopesV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else {
			params.Name.SetTo(name)
		}

		// Call the API to list access scopes with the name filter
		response, err := apiClient.ListAccessScopesV1(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list access scopes: %w", err)
		}

		// Check if we got any results
		if response != nil && len(response.Items) > 0 {
			// Even with the filter, verify the name matches exactly
			for _, scope := range response.Items {
				if scope.Name == name {
					// Return a copy of the access scope
					return &scope, nil
				}
			}
		}

		// If there's no next page, break the loop
		if response == nil || response.Pagination.NextPageToken.IsNull() || response.Pagination.NextPageToken.Value == "" {
			break
		}

		// Set the page token for the next iteration
		pageToken = response.Pagination.NextPageToken.Value
	}

	// No matching access scope found
	return nil, NewNotFoundByNameError("access scope", name)
}

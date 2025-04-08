package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// GetAccessScopeByName retrieves an access scope by its name.
func GetAccessScopeByName(ctx context.Context, apiClient client.Invoker, name string) (*client.AccessScopeV1, error) {
	params := client.ListAccessScopesV1Params{}
	params.Name.SetTo(name)

	response, err := apiClient.ListAccessScopesV1(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list access scopes: %w", err)
	}

	if response != nil && len(response.Items) > 0 {
		if len(response.Items) > 1 {
			return nil, fmt.Errorf("multiple access scopes found with name: %s", name)
		}
		return &response.Items[0], nil
	}

	return nil, NewNotFoundByNameError("access scope", name)
}

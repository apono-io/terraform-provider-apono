package common

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// GetAccessScopeByName retrieves all access scopes matching the given name.
func GetAccessScopeByName(ctx context.Context, apiClient client.Invoker, name string) ([]client.AccessScopeV1, error) {
	results := []client.AccessScopeV1{}
	pageToken := ""

	for {
		params := client.ListAccessScopesV1Params{}

		if name != "" {
			params.Name.SetTo(name)
		}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		}

		resp, err := apiClient.ListAccessScopesV1(ctx, params)
		if err != nil {
			return nil, err
		}

		results = append(results, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	// Sort results by name before returning
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

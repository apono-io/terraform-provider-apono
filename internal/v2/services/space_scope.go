package services

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// ListSpaceScopes retrieves all space scopes matching the provided name filter.
func ListSpaceScopes(ctx context.Context, apiClient client.Invoker, name string) ([]client.SpaceScopeV1, error) {
	results := []client.SpaceScopeV1{}
	pageToken := ""

	for {
		params := client.ListSpaceScopesV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else if name != "" {
			params.Name.SetTo(name)
		}

		resp, err := apiClient.ListSpaceScopesV1(ctx, params)
		if err != nil {
			return nil, err
		}

		results = append(results, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	// Sort results by id for consistency
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
}

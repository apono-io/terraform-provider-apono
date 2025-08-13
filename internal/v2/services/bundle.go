package services

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func ListBundles(ctx context.Context, apiClient client.Invoker, name string) ([]client.BundleV2, error) {
	results := []client.BundleV2{}
	pageToken := ""

	for {
		params := client.ListBundlesV2Params{}
		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else if name != "" {
			params.Name.SetTo(name)
		}

		resp, err := apiClient.ListBundlesV2(ctx, params)
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

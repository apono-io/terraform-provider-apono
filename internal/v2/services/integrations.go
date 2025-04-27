package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func ListIntegrations(ctx context.Context, apiClient client.Invoker, typeName string, name string, categories []string) ([]client.IntegrationV4, error) {
	allIntegrations := []client.IntegrationV4{}
	pageToken := ""

	for {
		params := client.ListIntegrationsV4Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else {
			if name != "" {
				params.Name.SetTo(name)
			}

			if typeName != "" {
				params.Type.SetTo([]string{typeName})
			}

			if len(categories) > 0 {
				params.Category.SetTo(categories)
			}
		}

		resp, err := apiClient.ListIntegrationsV4(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list integrations: %w", err)
		}

		allIntegrations = append(allIntegrations, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	// Sort integrations by id for consistency
	sort.Slice(allIntegrations, func(i, j int) bool {
		return allIntegrations[i].ID < allIntegrations[j].ID
	})

	return allIntegrations, nil
}

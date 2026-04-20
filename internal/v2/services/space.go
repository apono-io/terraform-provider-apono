package services

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func ListSpaceMembers(ctx context.Context, apiClient client.Invoker, spaceID string) ([]client.SpaceMemberV1, error) {
	results := []client.SpaceMemberV1{}
	pageToken := ""

	for {
		params := client.ListSpaceMembersV1Params{ID: spaceID}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		}

		resp, err := apiClient.ListSpaceMembersV1(ctx, params)
		if err != nil {
			return nil, err
		}

		results = append(results, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	return results, nil
}

func ListSpaces(ctx context.Context, apiClient client.Invoker, name string) ([]client.SpaceV1, error) {
	results := []client.SpaceV1{}
	pageToken := ""

	for {
		params := client.ListSpacesV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else if name != "" {
			params.Name.SetTo(name)
		}

		resp, err := apiClient.ListSpacesV1(ctx, params)
		if err != nil {
			return nil, err
		}

		results = append(results, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
}

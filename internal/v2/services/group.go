package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

// ListGroupMembers retrieves all members for a specific group.
func ListGroupMembers(ctx context.Context, apiClient client.Invoker, groupID string) ([]client.GroupMemberV1, error) {
	results := []client.GroupMemberV1{}
	pageToken := ""

	for {
		params := client.ListGroupMembersV1Params{}

		params.ID = groupID

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		}

		resp, err := apiClient.ListGroupMembersV1(ctx, params)
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

// ListGroups retrieves all groups matching the provided name filter.
func ListGroups(ctx context.Context, apiClient client.Invoker, name string) ([]client.GroupV1, error) {
	allGroups := []client.GroupV1{}
	pageToken := ""

	for {
		params := client.ListGroupsV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else if name != "" {
			params.Name.SetTo(name)
		}

		resp, err := apiClient.ListGroupsV1(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to list groups: %w", err)
		}

		allGroups = append(allGroups, resp.Items...)

		if resp.Pagination.NextPageToken.Value == "" {
			break
		}

		pageToken = resp.Pagination.NextPageToken.Value
	}

	// Sort groups by id for consistency
	sort.Slice(allGroups, func(i, j int) bool {
		return allGroups[i].ID < allGroups[j].ID
	})

	return allGroups, nil
}

// TODO: remove this function when the API supports filtering by source integration.
func FilterGroupsBySourceIntegration(groups []client.GroupV1, sourceIntegration string) []client.GroupV1 {
	if sourceIntegration == "" {
		return groups
	}

	filtered := []client.GroupV1{}
	for _, group := range groups {
		sourceIntegrationID, sourceIntegrationIDExists := group.SourceIntegrationID.Get()
		sourceIntegrationName, sourceIntegrationNameExists := group.SourceIntegrationName.Get()

		if sourceIntegrationIDExists && sourceIntegration == sourceIntegrationID {
			filtered = append(filtered, group)
			continue
		}

		if sourceIntegrationNameExists && sourceIntegration == sourceIntegrationName {
			filtered = append(filtered, group)
			continue
		}
	}

	return filtered
}

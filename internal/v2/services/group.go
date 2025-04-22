package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupModel represents the Terraform model for an Apono group.
type GroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Members types.Set    `tfsdk:"members"`
}

// GroupDataModel represents an individual group in the groups data source.
type GroupDataModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	SourceIntegrationID types.String `tfsdk:"source_integration_id"`
}

// GroupsDataModel represents the data source model for groups.
type GroupsDataModel struct {
	Name              types.String     `tfsdk:"name"`
	SourceIntegration types.String     `tfsdk:"source_integration"`
	Groups            []GroupDataModel `tfsdk:"groups"`
}

func GroupToModel(group *client.GroupV1) GroupModel {
	return GroupModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
		// Members will be filled separately since they require a different API call
	}
}

func GroupToDataModel(group *client.GroupV1) GroupDataModel {
	model := GroupDataModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
	}

	if sourceIntegrationID, exists := group.SourceIntegrationID.Get(); exists {
		model.SourceIntegrationID = types.StringValue(sourceIntegrationID)
	} else {
		model.SourceIntegrationID = types.StringNull()
	}

	return model
}

func ListGroupMembers(ctx context.Context, apiClient client.Invoker, groupID string) ([]client.GroupMemberV1, error) {
	results := []client.GroupMemberV1{}
	pageToken := ""

	for {
		params := client.ListGroupMembersV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else {
			params.ID = groupID
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

	sort.Slice(allGroups, func(i, j int) bool {
		return allGroups[i].Name < allGroups[j].Name
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

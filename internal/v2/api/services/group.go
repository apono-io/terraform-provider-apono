package services

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GroupModel represents the Terraform model for an Apono group.
type GroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Members types.Set    `tfsdk:"members"`
}

func GroupToModel(group *client.GroupV1) GroupModel {
	return GroupModel{
		ID:   types.StringValue(group.ID),
		Name: types.StringValue(group.Name),
		// Members will be filled separately since they require a different API call
	}
}

// ListGroupMembers retrieves all group members for the given group ID.
func ListGroupMembers(ctx context.Context, apiClient client.Invoker, groupID string) ([]client.GroupMemberV1, error) {
	results := []client.GroupMemberV1{}
	pageToken := ""

	for {
		params := client.ListGroupMembersV1Params{
			ID: groupID,
		}

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

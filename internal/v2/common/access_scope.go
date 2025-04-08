package common

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AccessScopeModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Query types.String `tfsdk:"query"`
}

func AccessScopeToModel(accessScope *client.AccessScopeV1) *AccessScopeModel {
	return &AccessScopeModel{
		ID:    types.StringValue(accessScope.ID),
		Name:  types.StringValue(accessScope.Name),
		Query: types.StringValue(accessScope.Query),
	}
}

func AccessScopesToModels(apiScopes []client.AccessScopeV1) []AccessScopeModel {
	result := make([]AccessScopeModel, 0, len(apiScopes))
	for _, scope := range apiScopes {
		result = append(result, *AccessScopeToModel(&scope))
	}
	return result
}

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

package services

import (
	"context"
	"sort"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AccessScopeModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Query          types.String `tfsdk:"query"`
	SpaceReference types.String `tfsdk:"space_reference"`
	Space          types.Object `tfsdk:"space"`
}

type AccessScopeDataSourceItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Query       types.String `tfsdk:"query"`
	Space       types.Object `tfsdk:"space"`
}

func AccessScopeToModel(accessScope *client.AccessScopeV1) (*AccessScopeModel, diag.Diagnostics) {
	model := &AccessScopeModel{
		ID:    types.StringValue(accessScope.ID),
		Name:  types.StringValue(accessScope.Name),
		Query: types.StringValue(accessScope.Query),
	}

	if val, ok := accessScope.Description.Get(); ok {
		model.Description = types.StringValue(val)
	}

	spaceObj, diags := models.SpaceReferenceToObject(accessScope.Space)
	if diags.HasError() {
		return nil, diags
	}
	model.Space = spaceObj
	// Repopulate space_reference from the response space name so state stays consistent.
	// space_reference is name-only per spec; using SpaceName here prevents ID→name drift.
	if val, ok := accessScope.Space.Get(); ok {
		model.SpaceReference = types.StringValue(val.SpaceName)
	} else {
		model.SpaceReference = types.StringNull()
	}

	return model, nil
}

func AccessScopesToModels(apiScopes []client.AccessScopeV1) []AccessScopeModel {
	result := make([]AccessScopeModel, 0, len(apiScopes))
	for _, scope := range apiScopes {
		model, _ := AccessScopeToModel(&scope) // SpaceReferenceToObject cannot fail with controlled inputs
		result = append(result, *model)
	}
	return result
}

func AccessScopeToDataSourceItemModel(accessScope *client.AccessScopeV1) (*AccessScopeDataSourceItemModel, diag.Diagnostics) {
	full, diags := AccessScopeToModel(accessScope)
	if diags.HasError() {
		return nil, diags
	}
	return &AccessScopeDataSourceItemModel{
		ID:          full.ID,
		Name:        full.Name,
		Description: full.Description,
		Query:       full.Query,
		Space:       full.Space,
	}, nil
}

func AccessScopesToDataSourceItemModels(apiScopes []client.AccessScopeV1) []AccessScopeDataSourceItemModel {
	result := make([]AccessScopeDataSourceItemModel, 0, len(apiScopes))
	for _, scope := range apiScopes {
		model, _ := AccessScopeToDataSourceItemModel(&scope) // SpaceReferenceToObject cannot fail with controlled inputs
		result = append(result, *model)
	}
	return result
}

func ListAccessScopesByName(ctx context.Context, apiClient client.Invoker, name string, spaceReferences []string) ([]client.AccessScopeV1, error) {
	results := []client.AccessScopeV1{}
	pageToken := ""

	for {
		params := client.ListAccessScopesV1Params{}

		if pageToken != "" {
			params.PageToken.SetTo(pageToken)
		} else {
			if name != "" {
				params.Name.SetTo(name)
			}
			if len(spaceReferences) > 0 {
				params.SpaceReferences.SetTo(spaceReferences)
			}
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

	// Sort results by id for consistency
	sort.Slice(results, func(i, j int) bool {
		return results[i].ID < results[j].ID
	})

	return results, nil
}

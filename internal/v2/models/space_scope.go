package models

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SpaceScopeModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Query types.String `tfsdk:"query"`
}

type SpaceScopeDataModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Query types.String `tfsdk:"query"`
}

type SpaceScopesDataModel struct {
	Name        types.String          `tfsdk:"name"`
	SpaceScopes []SpaceScopeDataModel `tfsdk:"space_scopes"`
}

func SpaceScopeToModel(scope *client.SpaceScopeV1) SpaceScopeModel {
	return SpaceScopeModel{
		ID:    types.StringValue(scope.ID),
		Name:  types.StringValue(scope.Name),
		Query: types.StringValue(scope.Query),
	}
}

func SpaceScopeToDataModel(scope *client.SpaceScopeV1) SpaceScopeDataModel {
	return SpaceScopeDataModel{
		ID:    types.StringValue(scope.ID),
		Name:  types.StringValue(scope.Name),
		Query: types.StringValue(scope.Query),
	}
}

func SpaceScopesToDataModels(scopes []client.SpaceScopeV1) []SpaceScopeDataModel {
	result := make([]SpaceScopeDataModel, 0, len(scopes))
	for _, scope := range scopes {
		result = append(result, SpaceScopeToDataModel(&scope))
	}
	return result
}

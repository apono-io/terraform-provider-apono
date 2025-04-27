package models

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Members types.Set    `tfsdk:"members"`
}

type GroupDataModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	SourceIntegrationID   types.String `tfsdk:"source_integration_id"`
	SourceIntegrationName types.String `tfsdk:"source_integration_name"`
}

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
	}

	if sourceIntegrationName, exists := group.SourceIntegrationName.Get(); exists {
		model.SourceIntegrationName = types.StringValue(sourceIntegrationName)
	}

	return model
}

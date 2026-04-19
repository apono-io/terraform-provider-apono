package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SpaceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	SpaceScopeReferences types.Set    `tfsdk:"space_scope_references"`
	Members              types.Set    `tfsdk:"members"`
}

type SpaceMemberModel struct {
	IdentityReference types.String `tfsdk:"identity_reference"`
	IdentityType      types.String `tfsdk:"identity_type"`
	SpaceRoles        types.Set    `tfsdk:"space_roles"`
}

func SpaceToModel(ctx context.Context, space *client.SpaceV1) (SpaceModel, error) {
	scopeNames := make([]string, len(space.SpaceScopes))
	for i, scope := range space.SpaceScopes {
		scopeNames[i] = scope.Name
	}

	scopeRefsSet, diags := types.SetValueFrom(ctx, types.StringType, scopeNames)
	if diags.HasError() {
		return SpaceModel{}, fmt.Errorf("%s: %s", diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
	}

	return SpaceModel{
		ID:                   types.StringValue(space.ID),
		Name:                 types.StringValue(space.Name),
		SpaceScopeReferences: scopeRefsSet,
		// Members will be filled separately since they require a different API call
	}, nil
}

func SpaceMembersToModels(ctx context.Context, members []client.SpaceMemberV1) ([]SpaceMemberModel, error) {
	result := make([]SpaceMemberModel, len(members))
	for i, member := range members {
		// Map identity_reference back: user → email, group → name
		identityRef := member.Name
		if member.IdentityType == "user" {
			if email, ok := member.Email.Get(); ok {
				identityRef = email
			}
		}

		rolesSet, diags := types.SetValueFrom(ctx, types.StringType, member.SpaceRoles)
		if diags.HasError() {
			return nil, fmt.Errorf("%s: %s", diags.Errors()[0].Summary(), diags.Errors()[0].Detail())
		}

		result[i] = SpaceMemberModel{
			IdentityReference: types.StringValue(identityRef),
			IdentityType:      types.StringValue(member.IdentityType),
			SpaceRoles:        rolesSet,
		}
	}

	return result, nil
}

type SpaceDataModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	SpaceScopeReferences []types.String `tfsdk:"space_scope_references"`
}

type SpacesDataModel struct {
	Name   types.String     `tfsdk:"name"`
	Spaces []SpaceDataModel `tfsdk:"spaces"`
}

func SpaceToDataModel(space *client.SpaceV1) SpaceDataModel {
	scopeRefs := make([]types.String, len(space.SpaceScopes))
	for i, scope := range space.SpaceScopes {
		scopeRefs[i] = types.StringValue(scope.Name)
	}

	return SpaceDataModel{
		ID:                   types.StringValue(space.ID),
		Name:                 types.StringValue(space.Name),
		SpaceScopeReferences: scopeRefs,
	}
}

func SpacesToDataModels(spaces []client.SpaceV1) []SpaceDataModel {
	result := make([]SpaceDataModel, 0, len(spaces))
	for _, space := range spaces {
		result = append(result, SpaceToDataModel(&space))
	}
	return result
}

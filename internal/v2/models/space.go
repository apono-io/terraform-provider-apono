package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

type SpaceScopeDataRefModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type SpaceDataModel struct {
	ID          types.String             `tfsdk:"id"`
	Name        types.String             `tfsdk:"name"`
	SpaceScopes []SpaceScopeDataRefModel `tfsdk:"space_scopes"`
}

// SpaceReferenceModel holds the space details returned by the API for a resource.
type SpaceReferenceModel struct {
	SpaceID   types.String `tfsdk:"space_id"`
	SpaceName types.String `tfsdk:"space_name"`
}

// SpaceReferenceAttrTypes defines the attribute types for a SpaceReferenceModel object.
var SpaceReferenceAttrTypes = map[string]attr.Type{
	"space_id":   types.StringType,
	"space_name": types.StringType,
}

// SpaceReferenceToObject converts an API OptNilSpaceReferenceV1 to a types.Object.
// Returns a null object if no space is set.
func SpaceReferenceToObject(space client.OptNilSpaceReferenceV1) (types.Object, diag.Diagnostics) {
	if val, ok := space.Get(); ok {
		return types.ObjectValue(SpaceReferenceAttrTypes, map[string]attr.Value{
			"space_id":   types.StringValue(val.SpaceID),
			"space_name": types.StringValue(val.SpaceName),
		})
	}
	return types.ObjectNull(SpaceReferenceAttrTypes), nil
}

type SpacesDataModel struct {
	Name   types.String     `tfsdk:"name"`
	Spaces []SpaceDataModel `tfsdk:"spaces"`
}

func SpaceToDataModel(space *client.SpaceV1) SpaceDataModel {
	scopes := make([]SpaceScopeDataRefModel, len(space.SpaceScopes))
	for i, scope := range space.SpaceScopes {
		scopes[i] = SpaceScopeDataRefModel{
			ID:   types.StringValue(scope.ID),
			Name: types.StringValue(scope.Name),
		}
	}

	return SpaceDataModel{
		ID:          types.StringValue(space.ID),
		Name:        types.StringValue(space.Name),
		SpaceScopes: scopes,
	}
}

func SpacesToDataModels(spaces []client.SpaceV1) []SpaceDataModel {
	result := make([]SpaceDataModel, 0, len(spaces))
	for _, space := range spaces {
		result = append(result, SpaceToDataModel(&space))
	}
	return result
}

func SpaceMemberModelsToAPI(ctx context.Context, members []SpaceMemberModel, diags *diag.Diagnostics) []client.UpsertSpaceMemberV1 {
	result := make([]client.UpsertSpaceMemberV1, len(members))
	for i, m := range members {
		var roles []string
		d := m.SpaceRoles.ElementsAs(ctx, &roles, false)
		diags.Append(d...)

		result[i] = client.UpsertSpaceMemberV1{
			IdentityReference: m.IdentityReference.ValueString(),
			IdentityType:      m.IdentityType.ValueString(),
			SpaceRoles:        roles,
		}
	}

	return result
}

func SpaceMemberObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"identity_reference": types.StringType,
			"identity_type":      types.StringType,
			"space_roles":        types.SetType{ElemType: types.StringType},
		},
	}
}

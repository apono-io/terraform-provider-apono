package datasources

import (
	"context"
	"errors"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoSpacesDataSource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	d := &AponoSpacesDataSource{client: mockInvoker}

	getSpacesListType := func() tftypes.List {
		return tftypes.List{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"id":                     tftypes.String,
					"name":                   tftypes.String,
					"space_scope_references": tftypes.List{ElementType: tftypes.String},
				},
			},
		}
	}

	getConfigType := func() tftypes.Object {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":   tftypes.String,
				"spaces": getSpacesListType(),
			},
		}
	}

	t.Run("Read_AllSpaces", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseSpacePublicV1Model{
			Items: []client.SpaceV1{
				{
					ID:   "space-123",
					Name: "Production",
					SpaceScopes: []client.SpaceScopeV1{
						{ID: "ss-1", Name: "Production AWS", Query: "q1"},
					},
				},
				{
					ID:   "space-456",
					Name: "Staging",
					SpaceScopes: []client.SpaceScopeV1{
						{ID: "ss-2", Name: "Staging AWS", Query: "q2"},
						{ID: "ss-3", Name: "Staging GCP", Query: "q3"},
					},
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return !params.Name.IsSet() && !params.Name.IsNull() && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spacesAttr := d.getTestSchema(ctx).Attributes["spaces"]
		spacesType, ok := spacesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, nil),
			"spaces": tftypes.NewValue(spacesType, nil),
		})

		s := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: s, Raw: configVal}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.SpacesDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, types.StringNull(), stateVal.Name)
		require.Len(t, stateVal.Spaces, 2)

		assert.Equal(t, "space-123", stateVal.Spaces[0].ID.ValueString())
		assert.Equal(t, "Production", stateVal.Spaces[0].Name.ValueString())
		require.Len(t, stateVal.Spaces[0].SpaceScopeReferences, 1)
		assert.Equal(t, "Production AWS", stateVal.Spaces[0].SpaceScopeReferences[0].ValueString())

		assert.Equal(t, "space-456", stateVal.Spaces[1].ID.ValueString())
		assert.Equal(t, "Staging", stateVal.Spaces[1].Name.ValueString())
		require.Len(t, stateVal.Spaces[1].SpaceScopeReferences, 2)
		assert.Equal(t, "Staging AWS", stateVal.Spaces[1].SpaceScopeReferences[0].ValueString())
		assert.Equal(t, "Staging GCP", stateVal.Spaces[1].SpaceScopeReferences[1].ValueString())
	})

	t.Run("Read_WithNameFilter", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseSpacePublicV1Model{
			Items: []client.SpaceV1{
				{
					ID:   "space-123",
					Name: "Production",
					SpaceScopes: []client.SpaceScopeV1{
						{ID: "ss-1", Name: "Production AWS", Query: "q1"},
					},
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpacesV1Params) bool {
				return params.Name.Value == "Production" && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spacesAttr := d.getTestSchema(ctx).Attributes["spaces"]
		spacesType, ok := spacesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, "Production"),
			"spaces": tftypes.NewValue(spacesType, nil),
		})

		s := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: s, Raw: configVal}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.SpacesDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "Production", stateVal.Name.ValueString())
		require.Len(t, stateVal.Spaces, 1)

		assert.Equal(t, "space-123", stateVal.Spaces[0].ID.ValueString())
		assert.Equal(t, "Production", stateVal.Spaces[0].Name.ValueString())
	})

	t.Run("Read_ApiError", func(t *testing.T) {
		mockInvoker.EXPECT().
			ListSpacesV1(mock.Anything, mock.Anything).
			Return(nil, errors.New("API error")).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spacesAttr := d.getTestSchema(ctx).Attributes["spaces"]
		spacesType, ok := spacesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, nil),
			"spaces": tftypes.NewValue(spacesType, nil),
		})

		s := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: s, Raw: configVal}
		state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "API error")
	})
}

func (d *AponoSpacesDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

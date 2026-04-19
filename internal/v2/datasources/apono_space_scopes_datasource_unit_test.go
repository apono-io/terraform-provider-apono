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

func TestAponoSpaceScopesDataSource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	d := &AponoSpaceScopesDataSource{client: mockInvoker}

	getSpaceScopesListType := func() tftypes.List {
		return tftypes.List{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"id":    tftypes.String,
					"name":  tftypes.String,
					"query": tftypes.String,
				},
			},
		}
	}

	getConfigType := func() tftypes.Object {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":         tftypes.String,
				"space_scopes": getSpaceScopesListType(),
			},
		}
	}

	t.Run("Read_AllScopes", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseSpaceScopePublicV1Model{
			Items: []client.SpaceScopeV1{
				{
					ID:    "ss-123456",
					Name:  "Production AWS",
					Query: `integration in ("aws-account") and resource_tag["environment"] = "production"`,
				},
				{
					ID:    "ss-789012",
					Name:  "Staging GCP",
					Query: `integration in ("gcp-project") and resource_tag["environment"] = "staging"`,
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		mockInvoker.EXPECT().
			ListSpaceScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceScopesV1Params) bool {
				return !params.Name.IsSet() && !params.Name.IsNull() && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spaceScopesAttr := d.getTestSchema(ctx).Attributes["space_scopes"]
		spaceScopesType, ok := spaceScopesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, nil),
			"space_scopes": tftypes.NewValue(spaceScopesType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.SpaceScopesDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, types.StringNull(), stateVal.Name)
		require.Len(t, stateVal.SpaceScopes, 2)

		assert.Equal(t, "ss-123456", stateVal.SpaceScopes[0].ID.ValueString())
		assert.Equal(t, "Production AWS", stateVal.SpaceScopes[0].Name.ValueString())
		assert.Equal(t, `integration in ("aws-account") and resource_tag["environment"] = "production"`, stateVal.SpaceScopes[0].Query.ValueString())

		assert.Equal(t, "ss-789012", stateVal.SpaceScopes[1].ID.ValueString())
		assert.Equal(t, "Staging GCP", stateVal.SpaceScopes[1].Name.ValueString())
		assert.Equal(t, `integration in ("gcp-project") and resource_tag["environment"] = "staging"`, stateVal.SpaceScopes[1].Query.ValueString())
	})

	t.Run("Read_WithNameFilter", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseSpaceScopePublicV1Model{
			Items: []client.SpaceScopeV1{
				{
					ID:    "ss-123456",
					Name:  "Production AWS",
					Query: `integration in ("aws-account")`,
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		mockInvoker.EXPECT().
			ListSpaceScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListSpaceScopesV1Params) bool {
				return params.Name.Value == "*AWS*" && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spaceScopesAttr := d.getTestSchema(ctx).Attributes["space_scopes"]
		spaceScopesType, ok := spaceScopesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, "*AWS*"),
			"space_scopes": tftypes.NewValue(spaceScopesType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.SpaceScopesDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "*AWS*", stateVal.Name.ValueString())
		require.Len(t, stateVal.SpaceScopes, 1)

		assert.Equal(t, "ss-123456", stateVal.SpaceScopes[0].ID.ValueString())
		assert.Equal(t, "Production AWS", stateVal.SpaceScopes[0].Name.ValueString())
		assert.Equal(t, `integration in ("aws-account")`, stateVal.SpaceScopes[0].Query.ValueString())
	})

	t.Run("Read_ApiError", func(t *testing.T) {
		mockInvoker.EXPECT().
			ListSpaceScopesV1(mock.Anything, mock.Anything).
			Return(nil, errors.New("API error")).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		spaceScopesAttr := d.getTestSchema(ctx).Attributes["space_scopes"]
		spaceScopesType, ok := spaceScopesAttr.GetType().TerraformType(ctx).(tftypes.List)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":         tftypes.NewValue(tftypes.String, nil),
			"space_scopes": tftypes.NewValue(spaceScopesType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "API error")
	})
}

func (d *AponoSpaceScopesDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

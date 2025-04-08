package datasources

import (
	"context"
	"errors"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoAccessScopesDataSource_Unit(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	d := &AponoAccessScopesDataSource{client: mockInvoker}

	t.Run("Read_AllScopes", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseAccessScopePublicV1Model{
			Items: []client.AccessScopeV1{
				{
					ID:    "as-123456",
					Name:  "test-scope-1",
					Query: `resource_type = "mock-1"`,
				},
				{
					ID:    "as-789012",
					Name:  "test-scope-2",
					Query: `resource_type = "mock-2"`,
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		mockInvoker.EXPECT().
			ListAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListAccessScopesV1Params) bool {
				return !params.Name.IsSet() && !params.Name.IsNull() && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		accessScopesAttr := d.getTestSchema(ctx).Attributes["access_scopes"]
		accessScopesType := accessScopesAttr.GetType().TerraformType(ctx).(tftypes.List)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":          tftypes.NewValue(tftypes.String, nil),
			"access_scopes": tftypes.NewValue(accessScopesType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal accessScopesDataSourceModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, types.StringNull(), stateVal.Name)
		require.Len(t, stateVal.AccessScopes, 2)

		assert.Equal(t, "as-123456", stateVal.AccessScopes[0].ID.ValueString())
		assert.Equal(t, "test-scope-1", stateVal.AccessScopes[0].Name.ValueString())
		assert.Equal(t, `resource_type = "mock-1"`, stateVal.AccessScopes[0].Query.ValueString())

		assert.Equal(t, "as-789012", stateVal.AccessScopes[1].ID.ValueString())
		assert.Equal(t, "test-scope-2", stateVal.AccessScopes[1].Name.ValueString())
		assert.Equal(t, `resource_type = "mock-2"`, stateVal.AccessScopes[1].Query.ValueString())
	})

	t.Run("Read_WithNameFilter", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseAccessScopePublicV1Model{
			Items: []client.AccessScopeV1{
				{
					ID:    "as-123456",
					Name:  "filtered-scope",
					Query: `resource_type = "filtered"`,
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.OptNilString{},
			},
		}

		nameParam := client.OptNilString{}
		nameParam.SetTo("filtered*")

		mockInvoker.EXPECT().
			ListAccessScopesV1(mock.Anything, mock.MatchedBy(func(params client.ListAccessScopesV1Params) bool {
				return params.Name.Value == "filtered*" && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		accessScopesAttr := d.getTestSchema(ctx).Attributes["access_scopes"]
		accessScopesType := accessScopesAttr.GetType().TerraformType(ctx).(tftypes.List)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":          tftypes.NewValue(tftypes.String, "filtered*"),
			"access_scopes": tftypes.NewValue(accessScopesType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal accessScopesDataSourceModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "filtered*", stateVal.Name.ValueString())
		require.Len(t, stateVal.AccessScopes, 1)

		assert.Equal(t, "as-123456", stateVal.AccessScopes[0].ID.ValueString())
		assert.Equal(t, "filtered-scope", stateVal.AccessScopes[0].Name.ValueString())
		assert.Equal(t, `resource_type = "filtered"`, stateVal.AccessScopes[0].Query.ValueString())
	})

	t.Run("Read_ApiError", func(t *testing.T) {
		mockInvoker.EXPECT().
			ListAccessScopesV1(mock.Anything, mock.Anything).
			Return(nil, errors.New("API error")).
			Once()

		ctx := t.Context()
		configType := getConfigType()
		accessScopesAttr := d.getTestSchema(ctx).Attributes["access_scopes"]
		accessScopesType := accessScopesAttr.GetType().TerraformType(ctx).(tftypes.List)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":          tftypes.NewValue(tftypes.String, nil),
			"access_scopes": tftypes.NewValue(accessScopesType, nil),
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

func (d *AponoAccessScopesDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

func getConfigType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":          tftypes.String,
			"access_scopes": getAccessScopesListType(),
		},
	}
}

func getAccessScopesListType() tftypes.List {
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

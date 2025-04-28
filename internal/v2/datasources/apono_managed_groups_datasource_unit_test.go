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

func TestAponoManagedGroupsDataSource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	d := &AponoManagedGroupsDataSource{client: mockInvoker}

	getGroupsSetType := func() tftypes.Set {
		return tftypes.Set{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"id":                      tftypes.String,
					"name":                    tftypes.String,
					"source_integration_id":   tftypes.String,
					"source_integration_name": tftypes.String,
				},
			},
		}
	}

	getGroupsConfigType := func() tftypes.Object {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":               tftypes.String,
				"source_integration": tftypes.String,
				"groups":             getGroupsSetType(),
			},
		}
	}

	t.Run("Read_AllGroups", func(t *testing.T) {
		mockListResponse := &client.PublicApiListResponseGroupPublicV1Model{
			Items: []client.GroupV1{
				{
					ID:                    "g-123456",
					Name:                  "test-group-1",
					SourceIntegrationID:   client.NewOptNilString("source-int-1"),
					SourceIntegrationName: client.NewOptNilString("Source Integration 1"),
				},
				{
					ID:   "g-789012",
					Name: "test-group-2",
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.NewOptNilString(""),
			},
		}

		mockInvoker.EXPECT().
			ListGroupsV1(mock.Anything, mock.MatchedBy(func(params client.ListGroupsV1Params) bool {
				return !params.Name.IsSet() && !params.Name.IsNull() && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getGroupsConfigType()
		groupsAttr := d.getTestSchema(ctx).Attributes["groups"]
		groupsType, ok := groupsAttr.GetType().TerraformType(ctx).(tftypes.Set)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":               tftypes.NewValue(tftypes.String, nil),
			"source_integration": tftypes.NewValue(tftypes.String, nil),
			"groups":             tftypes.NewValue(groupsType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.GroupsDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, types.StringNull(), stateVal.Name)
		assert.Equal(t, types.StringNull(), stateVal.SourceIntegration)
		require.Len(t, stateVal.Groups, 2)

		// Check the groups content (note: we can't guarantee the order in a set)
		foundGroup1 := false
		foundGroup2 := false

		for _, group := range stateVal.Groups {
			if group.ID.ValueString() == "g-123456" {
				foundGroup1 = true
				assert.Equal(t, "test-group-1", group.Name.ValueString())
				assert.Equal(t, "source-int-1", group.SourceIntegrationID.ValueString())
				assert.Equal(t, "Source Integration 1", group.SourceIntegrationName.ValueString())
			}
			if group.ID.ValueString() == "g-789012" {
				foundGroup2 = true
				assert.Equal(t, "test-group-2", group.Name.ValueString())
				assert.True(t, group.SourceIntegrationID.IsNull())
				assert.True(t, group.SourceIntegrationName.IsNull())
			}
		}

		assert.True(t, foundGroup1, "Group 1 not found in result")
		assert.True(t, foundGroup2, "Group 2 not found in result")
	})

	t.Run("Read_WithSourceIntegrationFilter", func(t *testing.T) {
		// For source integration filter, we first list all groups then filter in our code
		mockListResponse := &client.PublicApiListResponseGroupPublicV1Model{
			Items: []client.GroupV1{
				{
					ID:                    "g-123456",
					Name:                  "group-1",
					SourceIntegrationID:   client.NewOptNilString("source-int-1"),
					SourceIntegrationName: client.NewOptNilString("Source Integration 1"),
				},
				{
					ID:                    "g-789012",
					Name:                  "group-2",
					SourceIntegrationID:   client.NewOptNilString("source-int-2"),
					SourceIntegrationName: client.NewOptNilString("Source Integration 2"),
				},
			},
			Pagination: client.PublicApiPaginationInfoModel{
				NextPageToken: client.NewOptNilString(""),
			},
		}

		mockInvoker.EXPECT().
			ListGroupsV1(mock.Anything, mock.MatchedBy(func(params client.ListGroupsV1Params) bool {
				return !params.Name.IsSet() && !params.Name.IsNull() && !params.PageToken.IsSet() && !params.PageToken.IsNull()
			})).
			Return(mockListResponse, nil).
			Once()

		ctx := t.Context()
		configType := getGroupsConfigType()
		groupsAttr := d.getTestSchema(ctx).Attributes["groups"]
		groupsType, ok := groupsAttr.GetType().TerraformType(ctx).(tftypes.Set)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":               tftypes.NewValue(tftypes.String, nil),
			"source_integration": tftypes.NewValue(tftypes.String, "source-int-1"),
			"groups":             tftypes.NewValue(groupsType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var stateVal models.GroupsDataModel
		diags := resp.State.Get(ctx, &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, types.StringNull(), stateVal.Name)
		assert.Equal(t, "source-int-1", stateVal.SourceIntegration.ValueString())
		require.Len(t, stateVal.Groups, 1)
		assert.Equal(t, "g-123456", stateVal.Groups[0].ID.ValueString())
		assert.Equal(t, "group-1", stateVal.Groups[0].Name.ValueString())
		assert.Equal(t, "source-int-1", stateVal.Groups[0].SourceIntegrationID.ValueString())
		assert.Equal(t, "Source Integration 1", stateVal.Groups[0].SourceIntegrationName.ValueString())
	})

	t.Run("Read_ApiError", func(t *testing.T) {
		mockInvoker.EXPECT().
			ListGroupsV1(mock.Anything, mock.Anything).
			Return(nil, errors.New("API error")).
			Once()

		ctx := t.Context()
		configType := getGroupsConfigType()
		groupsAttr := d.getTestSchema(ctx).Attributes["groups"]
		groupsType, ok := groupsAttr.GetType().TerraformType(ctx).(tftypes.Set)
		require.True(t, ok)

		configVal := tftypes.NewValue(configType, map[string]tftypes.Value{
			"name":               tftypes.NewValue(tftypes.String, nil),
			"source_integration": tftypes.NewValue(tftypes.String, nil),
			"groups":             tftypes.NewValue(groupsType, nil),
		})

		schema := d.getTestSchema(ctx)
		config := tfsdk.Config{Schema: schema, Raw: configVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(configType, nil)}

		req := datasource.ReadRequest{Config: config}
		resp := datasource.ReadResponse{State: state}

		d.Read(ctx, req, &resp)

		require.True(t, resp.Diagnostics.HasError())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "Could not retrieve groups")
	})
}

func (d *AponoManagedGroupsDataSource) getTestSchema(ctx context.Context) schema.Schema {
	var resp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

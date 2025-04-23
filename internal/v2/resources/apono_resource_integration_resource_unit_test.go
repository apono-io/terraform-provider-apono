package resources

import (
	"context"
	"testing"

	"maps"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/mocks"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ogen-go/ogen/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAponoResourceIntegrationResource(t *testing.T) {
	mockInvoker := mocks.NewInvoker(t)
	r := &AponoResourceIntegrationResource{client: mockInvoker}

	getStateType := func() tftypes.Type {
		return tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"id":                       tftypes.String,
				"name":                     tftypes.String,
				"type":                     tftypes.String,
				"connector_id":             tftypes.String,
				"connected_resource_types": tftypes.List{ElementType: tftypes.String},
				"integration_config":       tftypes.Map{ElementType: tftypes.String},
				"secret_store_config": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"aws": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"region":    tftypes.String,
								"secret_id": tftypes.String,
							},
						},
						"gcp": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"project":   tftypes.String,
								"secret_id": tftypes.String,
							},
						},
						"azure": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"vault_url": tftypes.String,
								"name":      tftypes.String,
							},
						},
						"hashicorp_vault": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"secret_engine": tftypes.String,
								"path":          tftypes.String,
							},
						},
					},
				},
				"custom_access_details": tftypes.String,
				"owner": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"source_integration_name": tftypes.String,
						"type":                    tftypes.String,
						"values":                  tftypes.List{ElementType: tftypes.String},
					},
				},
				"owners_mapping": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"source_integration_name": tftypes.String,
						"key_name":                tftypes.String,
						"attribute_type":          tftypes.String,
					},
				},
			},
		}
	}

	t.Run("Create", func(t *testing.T) {
		planType := getStateType()

		integrationConfigVal := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
			"host":     tftypes.NewValue(tftypes.String, "db.example.com"),
			"database": tftypes.NewValue(tftypes.String, "test_db"),
			"port":     tftypes.NewValue(tftypes.String, "5432"),
		})

		awsSecretVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"region":    tftypes.String,
				"secret_id": tftypes.String,
			},
		}, map[string]tftypes.Value{
			"region":    tftypes.NewValue(tftypes.String, "us-west-2"),
			"secret_id": tftypes.NewValue(tftypes.String, "db-creds"),
		})

		secretStoreVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"aws": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"region":    tftypes.String,
						"secret_id": tftypes.String,
					},
				},
				"gcp": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"project":   tftypes.String,
						"secret_id": tftypes.String,
					},
				},
				"azure": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"vault_url": tftypes.String,
						"name":      tftypes.String,
					},
				},
				"hashicorp_vault": tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"secret_engine": tftypes.String,
						"path":          tftypes.String,
					},
				},
			},
		}, map[string]tftypes.Value{
			"aws":             awsSecretVal,
			"gcp":             tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, nil),
			"azure":           tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, nil),
			"hashicorp_vault": tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}, nil),
		})

		ownerVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"source_integration_name": tftypes.String,
				"type":                    tftypes.String,
				"values":                  tftypes.List{ElementType: tftypes.String},
			},
		}, map[string]tftypes.Value{
			"source_integration_name": tftypes.NewValue(tftypes.String, "identity-provider"),
			"type":                    tftypes.NewValue(tftypes.String, "group"),
			"values": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db-admins"),
			}),
		})

		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, nil),
			"name":         tftypes.NewValue(tftypes.String, "test-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-123"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
				tftypes.NewValue(tftypes.String, "schema"),
			}),
			"integration_config":    integrationConfigVal,
			"secret_store_config":   secretStoreVal,
			"custom_access_details": tftypes.NewValue(tftypes.String, "Access via VPN"),
			"owner":                 ownerVal,
			"owners_mapping": tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"source_integration_name": tftypes.String,
					"key_name":                tftypes.String,
					"attribute_type":          tftypes.String,
				},
			}, nil),
		})

		mockInvoker.EXPECT().
			CreateIntegrationV4(mock.Anything, mock.MatchedBy(func(request *client.CreateIntegrationV4) bool {
				return request.Name == "test-postgres" &&
					request.Type == "postgresql" &&
					request.ConnectorID.Value == "connector-123"
			})).
			Return(&client.IntegrationV4{
				ID:   "integration-123456",
				Name: "test-postgres",
				Type: "postgresql",
				ConnectorID: client.OptNilString{
					Value: "connector-123",
					Set:   true,
				},
				IntegrationConfig: map[string]jx.Raw{
					"host":     common.StringToJx("db.example.com"),
					"database": common.StringToJx("test_db"),
					"port":     common.StringToJx("5432"),
				},
			}, nil).
			Once()

		schema := r.getTestSchema(t.Context())
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(planType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}

		r.Create(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateVal models.ResourceIntegrationModel
		diags := resp.State.Get(t.Context(), &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "integration-123456", stateVal.ID.ValueString())
		assert.Equal(t, "test-postgres", stateVal.Name.ValueString())
		assert.Equal(t, "postgresql", stateVal.Type.ValueString())
		assert.Equal(t, "connector-123", stateVal.ConnectorID.ValueString())

		configMap := make(map[string]attr.Value)
		for k, v := range stateVal.IntegrationConfig.Elements() {
			configMap[k] = v
		}
		value, ok := configMap["host"]
		require.True(t, ok)

		hostVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "db.example.com", hostVal.ValueString())

		value, ok = configMap["database"]
		require.True(t, ok)

		dbVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "test_db", dbVal.ValueString())
	})

	t.Run("Create_WithMinimalConfig", func(t *testing.T) {
		planType := getStateType()

		planVal := tftypes.NewValue(planType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, nil),
			"name":         tftypes.NewValue(tftypes.String, "minimal-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-456"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host": tftypes.NewValue(tftypes.String, "minimal-db.example.com"),
			}),
			"secret_store_config": tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"aws": tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"region":    tftypes.String,
							"secret_id": tftypes.String,
						},
					},
					"gcp": tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"project":   tftypes.String,
							"secret_id": tftypes.String,
						},
					},
					"azure": tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"vault_url": tftypes.String,
							"name":      tftypes.String,
						},
					},
					"hashicorp_vault": tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"secret_engine": tftypes.String,
							"path":          tftypes.String,
						},
					},
				},
			}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, nil),
			"owner": tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"source_integration_name": tftypes.String,
					"type":                    tftypes.String,
					"values":                  tftypes.List{ElementType: tftypes.String},
				},
			}, nil),
			"owners_mapping": tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"source_integration_name": tftypes.String,
					"key_name":                tftypes.String,
					"attribute_type":          tftypes.String,
				},
			}, nil),
		})

		mockInvoker.EXPECT().
			CreateIntegrationV4(mock.Anything, mock.MatchedBy(func(request *client.CreateIntegrationV4) bool {
				return request.Name == "minimal-postgres" &&
					request.Type == "postgresql" &&
					request.ConnectorID.Value == "connector-456"
			})).
			Return(&client.IntegrationV4{
				ID:   "integration-456789",
				Name: "minimal-postgres",
				Type: "postgresql",
				ConnectorID: client.OptNilString{
					Value: "connector-456",
					Set:   true,
				},
				IntegrationConfig: map[string]jx.Raw{
					"host": common.StringToJx("minimal-db.example.com"),
				},
			}, nil).
			Once()

		schema := r.getTestSchema(t.Context())
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}
		state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(planType, nil)}

		req := resource.CreateRequest{Plan: plan}
		resp := resource.CreateResponse{State: state}

		r.Create(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var stateVal models.ResourceIntegrationModel
		diags := resp.State.Get(t.Context(), &stateVal)
		require.False(t, diags.HasError())

		assert.Equal(t, "integration-456789", stateVal.ID.ValueString())
		assert.Equal(t, "minimal-postgres", stateVal.Name.ValueString())
		assert.Equal(t, "postgresql", stateVal.Type.ValueString())
	})

	t.Run("Read", func(t *testing.T) {
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "integration-123456"),
			"name":         tftypes.NewValue(tftypes.String, "existing-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-123"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host": tftypes.NewValue(tftypes.String, "old-db.example.com"),
			}),
			"secret_store_config":   tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"aws": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}}, "gcp": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, "azure": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, "hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}}}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, nil),
			"owner":                 tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "type": tftypes.String, "values": tftypes.List{ElementType: tftypes.String}}}, nil),
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		mockIntegration := &client.IntegrationV4{
			ID:   "integration-123456",
			Name: "updated-postgres",
			Type: "postgresql",
			ConnectorID: client.OptNilString{
				Value: "connector-123",
				Set:   true,
			},
			ConnectedResourceTypes: client.OptNilStringArray{
				Value: []string{"db", "schema"},
				Set:   true,
			},
			IntegrationConfig: map[string]jx.Raw{
				"host":     common.StringToJx("new-db.example.com"),
				"database": common.StringToJx("prod_db"),
			},
			CustomAccessDetails: client.OptNilString{
				Value: "Updated access details",
				Set:   true,
			},
			Category: common.ResourceCategory,
		}

		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: "integration-123456"}).
			Return(mockIntegration, nil).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
		var updatedState models.ResourceIntegrationModel
		diags := resp.State.Get(t.Context(), &updatedState)
		require.False(t, diags.HasError())

		assert.Equal(t, "integration-123456", updatedState.ID.ValueString())
		assert.Equal(t, "updated-postgres", updatedState.Name.ValueString())
		assert.Equal(t, "postgresql", updatedState.Type.ValueString())
		assert.Equal(t, "connector-123", updatedState.ConnectorID.ValueString())

		var connectedTypes []string
		diags = updatedState.ConnectedResourceTypes.ElementsAs(t.Context(), &connectedTypes, false)
		require.False(t, diags.HasError())
		require.Len(t, connectedTypes, 2)
		assert.Contains(t, connectedTypes, "db")
		assert.Contains(t, connectedTypes, "schema")

		configMap := make(map[string]attr.Value)
		for k, v := range updatedState.IntegrationConfig.Elements() {
			configMap[k] = v
		}
		value, ok := configMap["host"]
		require.True(t, ok)

		hostVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "new-db.example.com", hostVal.ValueString())

		value, ok = configMap["database"]
		require.True(t, ok)

		dbVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "prod_db", dbVal.ValueString())

		assert.Equal(t, "Updated access details", updatedState.CustomAccessDetails.ValueString())
	})

	t.Run("Read_NotFound", func(t *testing.T) {
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "deleted-integration"),
			"name":         tftypes.NewValue(tftypes.String, "deleted-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-999"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host": tftypes.NewValue(tftypes.String, "deleted-db.example.com"),
			}),
			"secret_store_config":   tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"aws": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}}, "gcp": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, "azure": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, "hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}}}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, nil),
			"owner":                 tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "type": tftypes.String, "values": tftypes.List{ElementType: tftypes.String}}}, nil),
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		notFoundErr := &validate.UnexpectedStatusCodeError{StatusCode: 404}
		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: "deleted-integration"}).
			Return(nil, notFoundErr).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.ReadRequest{State: state}
		resp := resource.ReadResponse{State: state}

		r.Read(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("Update", func(t *testing.T) {
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "integration-123456"),
			"name":         tftypes.NewValue(tftypes.String, "old-postgres-name"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-123"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host":     tftypes.NewValue(tftypes.String, "old-db.example.com"),
				"port":     tftypes.NewValue(tftypes.String, "5432"),
				"database": tftypes.NewValue(tftypes.String, "old_db"),
			}),
			"secret_store_config": tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"aws":             tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}},
				"gcp":             tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}},
				"azure":           tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}},
				"hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}},
			}}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, "Old access details"),
			"owner":                 tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "type": tftypes.String, "values": tftypes.List{ElementType: tftypes.String}}}, nil),
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		awsSecretVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"region":    tftypes.String,
				"secret_id": tftypes.String,
			},
		}, map[string]tftypes.Value{
			"region":    tftypes.NewValue(tftypes.String, "us-west-2"),
			"secret_id": tftypes.NewValue(tftypes.String, "new-secret-id"),
		})

		secretStoreVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"aws":             tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}},
				"gcp":             tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}},
				"azure":           tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}},
				"hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}},
			},
		}, map[string]tftypes.Value{
			"aws":             awsSecretVal,
			"gcp":             tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, nil),
			"azure":           tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, nil),
			"hashicorp_vault": tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}, nil),
		})

		ownerVal := tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"source_integration_name": tftypes.String,
				"type":                    tftypes.String,
				"values":                  tftypes.List{ElementType: tftypes.String},
			},
		}, map[string]tftypes.Value{
			"source_integration_name": tftypes.NewValue(tftypes.String, "identity-provider"),
			"type":                    tftypes.NewValue(tftypes.String, "group"),
			"values": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db-devs"),
				tftypes.NewValue(tftypes.String, "db-admins"),
			}),
		})

		planVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "integration-123456"),
			"name":         tftypes.NewValue(tftypes.String, "updated-postgres-name"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-123"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
				tftypes.NewValue(tftypes.String, "schema"),
				tftypes.NewValue(tftypes.String, "table"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host":     tftypes.NewValue(tftypes.String, "new-db.example.com"),
				"port":     tftypes.NewValue(tftypes.String, "5433"),
				"database": tftypes.NewValue(tftypes.String, "new_db"),
			}),
			"secret_store_config":   secretStoreVal,
			"custom_access_details": tftypes.NewValue(tftypes.String, "Updated access details"),
			"owner":                 ownerVal,
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		mockInvoker.EXPECT().
			UpdateIntegrationV4(mock.Anything, mock.MatchedBy(func(request *client.UpdateIntegrationV4) bool {
				return request.Name == "updated-postgres-name" &&
					len(request.IntegrationConfig) == 3 &&
					request.CustomAccessDetails.Value == "Updated access details"
			}), client.UpdateIntegrationV4Params{ID: "integration-123456"}).
			Return(&client.IntegrationV4{
				ID:   "integration-123456",
				Name: "updated-postgres-name",
				Type: "postgresql",
				ConnectorID: client.OptNilString{
					Value: "connector-123",
					Set:   true,
				},
				ConnectedResourceTypes: client.OptNilStringArray{
					Value: []string{"db", "schema", "table"},
					Set:   true,
				},
				IntegrationConfig: map[string]jx.Raw{
					"host":     common.StringToJx("new-db.example.com"),
					"port":     common.StringToJx("5433"),
					"database": common.StringToJx("new_db"),
				},
				CustomAccessDetails: client.OptNilString{
					Value: "Updated access details",
					Set:   true,
				},
				SecretStoreConfig: client.OptNilIntegrationV4SecretStoreConfig{
					Value: client.IntegrationV4SecretStoreConfig{
						AWS: client.OptNilIntegrationV4SecretStoreConfigAWS{
							Value: client.IntegrationV4SecretStoreConfigAWS{
								Region:   "us-west-2",
								SecretID: "new-secret-id",
							},
							Set: true,
						},
					},
					Set: true,
				},
				Owner: client.OptNilIntegrationV4Owner{
					Value: client.IntegrationV4Owner{
						AttributeType:  "group",
						AttributeValue: []string{"db-devs", "db-admins"},
						SourceIntegrationName: client.OptNilString{
							Value: "identity-provider",
							Set:   true,
						},
					},
					Set: true,
				},
			}, nil).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}
		plan := tfsdk.Plan{Schema: schema, Raw: planVal}

		req := resource.UpdateRequest{
			State: state,
			Plan:  plan,
		}
		resp := resource.UpdateResponse{
			State: state,
		}

		r.Update(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())

		var updatedState models.ResourceIntegrationModel
		diags := resp.State.Get(t.Context(), &updatedState)
		require.False(t, diags.HasError())

		assert.Equal(t, "integration-123456", updatedState.ID.ValueString())
		assert.Equal(t, "updated-postgres-name", updatedState.Name.ValueString())
		assert.Equal(t, "postgresql", updatedState.Type.ValueString())

		var connectedTypes []string
		diags = updatedState.ConnectedResourceTypes.ElementsAs(t.Context(), &connectedTypes, false)
		require.False(t, diags.HasError())
		require.Len(t, connectedTypes, 3)
		assert.Contains(t, connectedTypes, "db")
		assert.Contains(t, connectedTypes, "schema")
		assert.Contains(t, connectedTypes, "table")

		configMap := make(map[string]attr.Value)
		for k, v := range updatedState.IntegrationConfig.Elements() {
			configMap[k] = v
		}
		value, ok := configMap["host"]
		require.True(t, ok)

		hostVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "new-db.example.com", hostVal.ValueString())

		value, ok = configMap["port"]
		require.True(t, ok)

		portVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "5433", portVal.ValueString())

		assert.Equal(t, "Updated access details", updatedState.CustomAccessDetails.ValueString())

		require.NotNil(t, updatedState.SecretStoreConfig)
		require.NotNil(t, updatedState.SecretStoreConfig.AWS)
		assert.Equal(t, "us-west-2", updatedState.SecretStoreConfig.AWS.Region.ValueString())
		assert.Equal(t, "new-secret-id", updatedState.SecretStoreConfig.AWS.SecretID.ValueString())

		require.NotNil(t, updatedState.Owner)
		assert.Equal(t, "group", updatedState.Owner.Type.ValueString())
		assert.Equal(t, "identity-provider", updatedState.Owner.SourceIntegrationName.ValueString())

		var ownerValues []string
		diags = updatedState.Owner.Values.ElementsAs(t.Context(), &ownerValues, false)
		require.False(t, diags.HasError())
		require.Len(t, ownerValues, 2)
		assert.Contains(t, ownerValues, "db-devs")
		assert.Contains(t, ownerValues, "db-admins")
	})

	t.Run("Delete", func(t *testing.T) {
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "integration-to-delete"),
			"name":         tftypes.NewValue(tftypes.String, "delete-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-123"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host": tftypes.NewValue(tftypes.String, "db-to-delete.example.com"),
			}),
			"secret_store_config":   tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"aws": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}}, "gcp": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, "azure": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, "hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}}}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, nil),
			"owner":                 tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "type": tftypes.String, "values": tftypes.List{ElementType: tftypes.String}}}, nil),
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		mockInvoker.EXPECT().
			DeleteIntegrationV4(mock.Anything, client.DeleteIntegrationV4Params{ID: "integration-to-delete"}).
			Return(nil).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{}

		r.Delete(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		stateType := getStateType()

		stateVal := tftypes.NewValue(stateType, map[string]tftypes.Value{
			"id":           tftypes.NewValue(tftypes.String, "non-existent-integration"),
			"name":         tftypes.NewValue(tftypes.String, "missing-postgres"),
			"type":         tftypes.NewValue(tftypes.String, "postgresql"),
			"connector_id": tftypes.NewValue(tftypes.String, "connector-999"),
			"connected_resource_types": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
				tftypes.NewValue(tftypes.String, "db"),
			}),
			"integration_config": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
				"host": tftypes.NewValue(tftypes.String, "missing-db.example.com"),
			}),
			"secret_store_config":   tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"aws": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"region": tftypes.String, "secret_id": tftypes.String}}, "gcp": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"project": tftypes.String, "secret_id": tftypes.String}}, "azure": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"vault_url": tftypes.String, "name": tftypes.String}}, "hashicorp_vault": tftypes.Object{AttributeTypes: map[string]tftypes.Type{"secret_engine": tftypes.String, "path": tftypes.String}}}}, nil),
			"custom_access_details": tftypes.NewValue(tftypes.String, nil),
			"owner":                 tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "type": tftypes.String, "values": tftypes.List{ElementType: tftypes.String}}}, nil),
			"owners_mapping":        tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"source_integration_name": tftypes.String, "key_name": tftypes.String, "attribute_type": tftypes.String}}, nil),
		})

		notFoundErr := &validate.UnexpectedStatusCodeError{StatusCode: 404}
		mockInvoker.EXPECT().
			DeleteIntegrationV4(mock.Anything, client.DeleteIntegrationV4Params{ID: "non-existent-integration"}).
			Return(notFoundErr).
			Once()

		schema := r.getTestSchema(t.Context())
		state := tfsdk.State{Schema: schema, Raw: stateVal}

		req := resource.DeleteRequest{State: state}
		resp := resource.DeleteResponse{}

		r.Delete(t.Context(), req, &resp)

		require.False(t, resp.Diagnostics.HasError())
	})

	t.Run("ImportState_ByID", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: "integration-import-id"}).
			Return(&client.IntegrationV4{
				ID:   "integration-import-id",
				Name: "imported-integration",
				Type: "postgresql",
				ConnectorID: client.OptNilString{
					Value: "connector-123",
					Set:   true,
				},
				ConnectedResourceTypes: client.OptNilStringArray{
					Value: []string{"db", "schema"},
					Set:   true,
				},
				IntegrationConfig: map[string]jx.Raw{
					"host": common.StringToJx("imported-db.example.com"),
				},
				Category: common.ResourceCategory,
			}, nil).
			Once()

		ctx := t.Context()
		req := resource.ImportStateRequest{ID: "integration-import-id"}

		stateType := getStateType()
		schema := r.getTestSchema(ctx)
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{Schema: schema, Raw: tftypes.NewValue(stateType, nil)},
		}

		r.ImportState(ctx, req, &importResp)

		require.False(t, importResp.Diagnostics.HasError())

		readReq := resource.ReadRequest{State: importResp.State}
		readResp := resource.ReadResponse{State: importResp.State}
		r.Read(ctx, readReq, &readResp)

		require.False(t, readResp.Diagnostics.HasError())
		var stateModel models.ResourceIntegrationModel
		diags := readResp.State.Get(ctx, &stateModel)
		require.False(t, diags.HasError())

		assert.Equal(t, "integration-import-id", stateModel.ID.ValueString())
		assert.Equal(t, "imported-integration", stateModel.Name.ValueString())
		assert.Equal(t, "postgresql", stateModel.Type.ValueString())
		assert.Equal(t, "connector-123", stateModel.ConnectorID.ValueString())

		var connectedTypes []string
		diags = stateModel.ConnectedResourceTypes.ElementsAs(ctx, &connectedTypes, false)
		require.False(t, diags.HasError())
		require.Len(t, connectedTypes, 2)
		assert.Contains(t, connectedTypes, "db")
		assert.Contains(t, connectedTypes, "schema")

		configMap := make(map[string]attr.Value)
		maps.Copy(configMap, stateModel.IntegrationConfig.Elements())
		value, ok := configMap["host"]
		require.True(t, ok)

		hostVal, isString := value.(types.String)
		require.True(t, isString)
		assert.Equal(t, "imported-db.example.com", hostVal.ValueString())
	})

	t.Run("ImportState_WrongCategory", func(t *testing.T) {
		mockInvoker.EXPECT().
			GetIntegrationsByIdV4(mock.Anything, client.GetIntegrationsByIdV4Params{ID: "integration-wrong-category"}).
			Return(&client.IntegrationV4{
				ID:   "integration-wrong-category",
				Name: "wrong-category-integration",
				Type: "postgresql",
				ConnectorID: client.OptNilString{
					Value: "connector-123",
					Set:   true,
				},
				IntegrationConfig: map[string]jx.Raw{
					"host": common.StringToJx("wrong-category.example.com"),
				},
				Category: "IDENTITY",
			}, nil).
			Once()

		ctx := t.Context()
		req := resource.ImportStateRequest{ID: "integration-wrong-category"}

		stateType := getStateType()
		schema := r.getTestSchema(ctx)
		importResp := resource.ImportStateResponse{
			State: tfsdk.State{Schema: schema, Raw: tftypes.NewValue(stateType, nil)},
		}

		r.ImportState(ctx, req, &importResp)
		require.False(t, importResp.Diagnostics.HasError())

		readReq := resource.ReadRequest{State: importResp.State}
		readResp := resource.ReadResponse{State: importResp.State}
		r.Read(ctx, readReq, &readResp)

		require.True(t, readResp.Diagnostics.HasError())
		errorDiagnostic := readResp.Diagnostics.Errors()[0]
		assert.Contains(t, errorDiagnostic.Summary(), "Invalid resource integration type")
		assert.Contains(t, errorDiagnostic.Detail(), "Expected resource integration, got IDENTITY")
	})

}

func (r *AponoResourceIntegrationResource) getTestSchema(ctx context.Context) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

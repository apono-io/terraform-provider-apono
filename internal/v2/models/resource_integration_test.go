package models

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/go-faster/jx"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceIntegrationModelToCreateRequest(t *testing.T) {
	ctx := t.Context()

	t.Run("with connector_id and connected resource types", func(t *testing.T) {
		resourceTypes := []attr.Value{
			types.StringValue("database"),
			types.StringValue("schema"),
		}
		model := ResourceIntegrationModel{
			Name:                   types.StringValue("test-integration"),
			Type:                   types.StringValue("postgres"),
			ConnectorID:            types.StringValue("test-connector-id"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, resourceTypes),
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.Equal(t, "test-integration", req.Name)
		assert.Equal(t, "postgres", req.Type)
		assert.True(t, req.ConnectorID.IsSet())
		assert.Equal(t, "test-connector-id", req.ConnectorID.Value)
		assert.True(t, req.ConnectedResourceTypes.IsSet())
		assert.Equal(t, []string{"database", "schema"}, req.ConnectedResourceTypes.Value)
	})

	t.Run("with integration config", func(t *testing.T) {
		configMap := map[string]attr.Value{
			"host":     types.StringValue("localhost"),
			"port":     types.StringValue("5432"),
			"database": types.StringValue("postgres"),
		}
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			IntegrationConfig: types.MapValueMust(types.StringType, configMap),
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.NotNil(t, req.IntegrationConfig)
		hostVal, isHostString := model.IntegrationConfig.Elements()["host"].(types.String)
		require.True(t, isHostString)
		assert.Equal(t, "localhost", hostVal.ValueString())

		dbVal, isDbString := model.IntegrationConfig.Elements()["database"].(types.String)
		require.True(t, isDbString)
		assert.Equal(t, "postgres", dbVal.ValueString())
	})

	t.Run("with AWS secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				AWS: &AWSSecretConfig{
					Region:   types.StringValue("us-east-1"),
					SecretID: types.StringValue("secret-id"),
				},
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.AWS.IsSet())
		assert.Equal(t, "us-east-1", secretConfig.AWS.Value.Region)
		assert.Equal(t, "secret-id", secretConfig.AWS.Value.SecretID)
	})

	t.Run("with GCP secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				GCP: &GCPSecretConfig{
					Project:  types.StringValue("my-project"),
					SecretID: types.StringValue("secret-id"),
				},
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Gcp.IsSet())
		assert.Equal(t, "my-project", secretConfig.Gcp.Value.Project)
		assert.Equal(t, "secret-id", secretConfig.Gcp.Value.SecretID)
	})

	t.Run("with Azure secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				Azure: &AzureSecretConfig{
					VaultURL: types.StringValue("https://myvault.vault.azure.net"),
					Name:     types.StringValue("secret-name"),
				},
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Azure.IsSet())
		assert.Equal(t, "https://myvault.vault.azure.net", secretConfig.Azure.Value.VaultURL)
		assert.Equal(t, "secret-name", secretConfig.Azure.Value.Name)
	})

	t.Run("with HashiCorp Vault secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				HashicorpVault: &HashicorpVaultConfig{
					SecretEngine: types.StringValue("kv"),
					Path:         types.StringValue("secret/data/postgres"),
				},
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.HashicorpVault.IsSet())
		assert.Equal(t, "kv", secretConfig.HashicorpVault.Value.SecretEngine)
		assert.Equal(t, "secret/data/postgres", secretConfig.HashicorpVault.Value.Path)
	})

	t.Run("with Kubernetes secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				Kubernetes: &KubernetesSecretConfig{
					Namespace: types.StringValue("test-namespace"),
					Name:      types.StringValue("test-secret"),
				},
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Kubernetes.IsSet())
		assert.Equal(t, "test-namespace", secretConfig.Kubernetes.Value.Namespace)
		assert.Equal(t, "test-secret", secretConfig.Kubernetes.Value.Name)
	})

	t.Run("with custom access details", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			CustomAccessDetails: types.StringValue("Use your SSO credentials to login"),
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.CustomAccessDetails.IsSet())
		assert.Equal(t, "Use your SSO credentials to login", req.CustomAccessDetails.Value)
	})

	t.Run("with owner config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			Owner: &OwnerConfig{
				SourceIntegrationName: types.StringValue("source-integration"),
				AttributeType:         types.StringValue("user"),
				AttributeValues:       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("user1"), types.StringValue("user2")}),
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.Owner.IsSet())
		owner := req.Owner.Value
		assert.Equal(t, "user", owner.AttributeType)
		assert.Equal(t, []string{"user1", "user2"}, owner.AttributeValue)
		assert.True(t, owner.SourceIntegrationReference.IsSet())
		assert.Equal(t, "source-integration", owner.SourceIntegrationReference.Value)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			OwnersMapping: &OwnersMappingConfig{
				SourceIntegrationName: types.StringValue("source-integration"),
				KeyName:               types.StringValue("owner"),
				AttributeType:         types.StringValue("group"),
			},
		}

		req, err := ResourceIntegrationModelToCreateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.OwnersMapping.IsSet())
		mapping := req.OwnersMapping.Value
		assert.Equal(t, "owner", mapping.KeyName)
		assert.Equal(t, "group", mapping.AttributeType)
		assert.True(t, mapping.SourceIntegrationReference.IsSet())
		assert.Equal(t, "source-integration", mapping.SourceIntegrationReference.Value)
	})
}

func TestResourceIntegrationToModel(t *testing.T) {
	ctx := t.Context()

	t.Run("minimal fields", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		assert.Equal(t, "integration-id", model.ID.ValueString())
		assert.Equal(t, "test-integration", model.Name.ValueString())
		assert.Equal(t, "postgres", model.Type.ValueString())
		assert.Equal(t, "connector-id", model.ConnectorID.ValueString())
		assert.False(t, model.ConnectedResourceTypes.IsNull())
		assert.True(t, model.IntegrationConfig.IsNull())
		assert.Nil(t, model.SecretStoreConfig)
		assert.True(t, model.CustomAccessDetails.IsNull())
		assert.Nil(t, model.Owner)
		assert.Nil(t, model.OwnersMapping)
	})

	t.Run("with connected resource types", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database", "schema"}),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		assert.False(t, model.ConnectedResourceTypes.IsNull())

		var resourceTypes []string
		diags := model.ConnectedResourceTypes.ElementsAs(ctx, &resourceTypes, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"database", "schema"}, resourceTypes)
	})

	t.Run("with integration config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			IntegrationConfig: map[string]jx.Raw{
				"host":     jx.Raw("\"" + "localhost" + "\""),
				"port":     jx.Raw("\"" + "5432" + "\""),
				"database": jx.Raw("\"" + "postgres" + "\""),
			},
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		assert.False(t, model.IntegrationConfig.IsNull())

		expectedConfig := map[string]string{
			"host":     "localhost",
			"port":     "5432",
			"database": "postgres",
		}

		for key, expectedValue := range expectedConfig {
			value, ok := model.IntegrationConfig.Elements()[key]
			require.True(t, ok)

			strValue, isString := value.(types.String)
			require.True(t, isString)

			assert.Equal(t, expectedValue, strValue.ValueString())
		}
	})

	t.Run("with AWS secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			SecretStoreConfig: client.NewOptNilSecretStoreConfigV4(
				client.SecretStoreConfigV4{
					AWS: client.NewOptNilAwsSecretConfigV4(
						client.AwsSecretConfigV4{
							Region:   "us-east-1",
							SecretID: "secret-id",
						},
					),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.AWS)
		assert.Equal(t, "us-east-1", model.SecretStoreConfig.AWS.Region.ValueString())
		assert.Equal(t, "secret-id", model.SecretStoreConfig.AWS.SecretID.ValueString())
		assert.Nil(t, model.SecretStoreConfig.GCP)
		assert.Nil(t, model.SecretStoreConfig.Azure)
		assert.Nil(t, model.SecretStoreConfig.HashicorpVault)
	})

	t.Run("with GCP secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			SecretStoreConfig: client.NewOptNilSecretStoreConfigV4(
				client.SecretStoreConfigV4{
					Gcp: client.NewOptNilGcpSecretConfigV4(
						client.GcpSecretConfigV4{
							Project:  "my-project",
							SecretID: "secret-id",
						},
					),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.GCP)
		assert.Equal(t, "my-project", model.SecretStoreConfig.GCP.Project.ValueString())
		assert.Equal(t, "secret-id", model.SecretStoreConfig.GCP.SecretID.ValueString())
	})

	t.Run("with Azure secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			SecretStoreConfig: client.NewOptNilSecretStoreConfigV4(
				client.SecretStoreConfigV4{
					Azure: client.NewOptNilAzureSecretConfigV4(
						client.AzureSecretConfigV4{
							VaultURL: "https://myvault.vault.azure.net",
							Name:     "secret-name",
						},
					),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.Azure)
		assert.Equal(t, "https://myvault.vault.azure.net", model.SecretStoreConfig.Azure.VaultURL.ValueString())
		assert.Equal(t, "secret-name", model.SecretStoreConfig.Azure.Name.ValueString())
		assert.Nil(t, model.SecretStoreConfig.AWS)
		assert.Nil(t, model.SecretStoreConfig.GCP)
		assert.Nil(t, model.SecretStoreConfig.HashicorpVault)
	})

	t.Run("with HashiCorp Vault secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			SecretStoreConfig: client.NewOptNilSecretStoreConfigV4(
				client.SecretStoreConfigV4{
					HashicorpVault: client.NewOptNilHashicorpVaultSecretConfigV4(
						client.HashicorpVaultSecretConfigV4{
							SecretEngine: "kv",
							Path:         "secret/data/postgres",
						},
					),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.HashicorpVault)
		assert.Equal(t, "kv", model.SecretStoreConfig.HashicorpVault.SecretEngine.ValueString())
		assert.Equal(t, "secret/data/postgres", model.SecretStoreConfig.HashicorpVault.Path.ValueString())
		assert.Nil(t, model.SecretStoreConfig.AWS)
		assert.Nil(t, model.SecretStoreConfig.GCP)
		assert.Nil(t, model.SecretStoreConfig.Azure)
	})

	t.Run("with Kubernetes secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			SecretStoreConfig: client.NewOptNilSecretStoreConfigV4(
				client.SecretStoreConfigV4{
					Kubernetes: client.NewOptNilKubernetesSecretConfigV4(
						client.KubernetesSecretConfigV4{
							Namespace: "my-namespace",
							Name:      "my-secret",
						},
					),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.Kubernetes)
		assert.Equal(t, "my-namespace", model.SecretStoreConfig.Kubernetes.Namespace.ValueString())
		assert.Equal(t, "my-secret", model.SecretStoreConfig.Kubernetes.Name.ValueString())
		assert.Nil(t, model.SecretStoreConfig.AWS)
		assert.Nil(t, model.SecretStoreConfig.GCP)
		assert.Nil(t, model.SecretStoreConfig.Azure)
		assert.Nil(t, model.SecretStoreConfig.HashicorpVault)
	})

	t.Run("with custom access details", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			CustomAccessDetails:    client.NewOptNilString("Use your SSO credentials to login"),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		assert.Equal(t, "Use your SSO credentials to login", model.CustomAccessDetails.ValueString())
	})

	t.Run("with owner config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			Owner: client.NewOptNilOwnerV4(
				client.OwnerV4{
					AttributeType:         "user",
					AttributeValue:        []string{"user1", "user2"},
					SourceIntegrationName: client.NewOptNilString("source-integration"),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.Owner)
		assert.Equal(t, "user", model.Owner.AttributeType.ValueString())
		assert.Equal(t, "source-integration", model.Owner.SourceIntegrationName.ValueString())

		var values []string
		diags := model.Owner.AttributeValues.ElementsAs(ctx, &values, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"user1", "user2"}, values)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:                     "integration-id",
			Name:                   "test-integration",
			Type:                   "postgres",
			ConnectorID:            client.NewOptNilString("connector-id"),
			ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database"}),
			OwnersMapping: client.NewOptNilOwnerMappingV4(
				client.OwnerMappingV4{
					KeyName:               "owner",
					AttributeType:         "group",
					SourceIntegrationName: client.NewOptNilString("source-integration"),
				},
			),
		}

		model, err := ResourceIntegrationToModel(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.OwnersMapping)
		assert.Equal(t, "owner", model.OwnersMapping.KeyName.ValueString())
		assert.Equal(t, "group", model.OwnersMapping.AttributeType.ValueString())
		assert.Equal(t, "source-integration", model.OwnersMapping.SourceIntegrationName.ValueString())
	})
}

func TestResourceIntegrationModelToUpdateRequest(t *testing.T) {
	ctx := t.Context()

	t.Run("minimal fields", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.Equal(t, "updated-integration", req.Name)
	})

	t.Run("with integration config", func(t *testing.T) {
		configMap := map[string]attr.Value{
			"host":     types.StringValue("new-host"),
			"port":     types.StringValue("5433"),
			"database": types.StringValue("test-db"),
		}
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			IntegrationConfig: types.MapValueMust(types.StringType, configMap),
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.NotNil(t, req.IntegrationConfig)

		hostVal, isHostString := model.IntegrationConfig.Elements()["host"].(types.String)
		require.True(t, isHostString)
		assert.Equal(t, "new-host", hostVal.ValueString())

		portVal, isPortString := model.IntegrationConfig.Elements()["port"].(types.String)
		require.True(t, isPortString)
		assert.Equal(t, "5433", portVal.ValueString())

		dbVal, isDbString := model.IntegrationConfig.Elements()["database"].(types.String)
		require.True(t, isDbString)
		assert.Equal(t, "test-db", dbVal.ValueString())
	})

	t.Run("with connected resource types", func(t *testing.T) {
		resourceTypes := []attr.Value{
			types.StringValue("table"),
			types.StringValue("view"),
		}
		model := ResourceIntegrationModel{
			Name:                   types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, resourceTypes),
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.ConnectedResourceTypes.IsSet())
		assert.Equal(t, []string{"table", "view"}, req.ConnectedResourceTypes.Value)
	})

	t.Run("with AWS secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				AWS: &AWSSecretConfig{
					Region:   types.StringValue("us-west-2"),
					SecretID: types.StringValue("new-secret-id"),
				},
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.AWS.IsSet())
		assert.Equal(t, "us-west-2", secretConfig.AWS.Value.Region)
		assert.Equal(t, "new-secret-id", secretConfig.AWS.Value.SecretID)
	})

	t.Run("with GCP secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				GCP: &GCPSecretConfig{
					Project:  types.StringValue("new-project"),
					SecretID: types.StringValue("new-secret-id"),
				},
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Gcp.IsSet())
		assert.Equal(t, "new-project", secretConfig.Gcp.Value.Project)
		assert.Equal(t, "new-secret-id", secretConfig.Gcp.Value.SecretID)
	})

	t.Run("with Azure secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				Azure: &AzureSecretConfig{
					VaultURL: types.StringValue("https://newvault.vault.azure.net"),
					Name:     types.StringValue("new-secret-name"),
				},
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Azure.IsSet())
		assert.Equal(t, "https://newvault.vault.azure.net", secretConfig.Azure.Value.VaultURL)
		assert.Equal(t, "new-secret-name", secretConfig.Azure.Value.Name)
	})

	t.Run("with HashiCorp Vault secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				HashicorpVault: &HashicorpVaultConfig{
					SecretEngine: types.StringValue("kv2"),
					Path:         types.StringValue("secret/data/new-postgres"),
				},
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.HashicorpVault.IsSet())
		assert.Equal(t, "kv2", secretConfig.HashicorpVault.Value.SecretEngine)
		assert.Equal(t, "secret/data/new-postgres", secretConfig.HashicorpVault.Value.Path)
	})

	t.Run("with Kubernetes secret store config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			SecretStoreConfig: &SecretStoreConfig{
				Kubernetes: &KubernetesSecretConfig{
					Namespace: types.StringValue("updated-namespace"),
					Name:      types.StringValue("updated-secret"),
				},
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Kubernetes.IsSet())
		assert.Equal(t, "updated-namespace", secretConfig.Kubernetes.Value.Namespace)
		assert.Equal(t, "updated-secret", secretConfig.Kubernetes.Value.Name)
	})

	t.Run("with custom access details", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			CustomAccessDetails: types.StringValue("Updated access instructions"),
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.CustomAccessDetails.IsSet())
		assert.Equal(t, "Updated access instructions", req.CustomAccessDetails.Value)
	})

	t.Run("with owner config", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			Owner: &OwnerConfig{
				SourceIntegrationName: types.StringValue("new-source-integration"),
				AttributeType:         types.StringValue("group"),
				AttributeValues:       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("group1"), types.StringValue("group2")}),
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.Owner.IsSet())
		owner := req.Owner.Value
		assert.Equal(t, "group", owner.AttributeType)
		assert.Equal(t, []string{"group1", "group2"}, owner.AttributeValue)
		assert.True(t, owner.SourceIntegrationReference.IsSet())
		assert.Equal(t, "new-source-integration", owner.SourceIntegrationReference.Value)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("database"),
			}),
			OwnersMapping: &OwnersMappingConfig{
				SourceIntegrationName: types.StringValue("new-source-integration"),
				KeyName:               types.StringValue("new-owner"),
				AttributeType:         types.StringValue("tag"),
			},
		}

		req, err := ResourceIntegrationModelToUpdateRequest(ctx, model)

		require.NoError(t, err)
		assert.True(t, req.OwnersMapping.IsSet())
		mapping := req.OwnersMapping.Value
		assert.Equal(t, "new-owner", mapping.KeyName)
		assert.Equal(t, "tag", mapping.AttributeType)
		assert.True(t, mapping.SourceIntegrationReference.IsSet())
		assert.Equal(t, "new-source-integration", mapping.SourceIntegrationReference.Value)
	})
}

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

func TestCreateIntegrationRequest(t *testing.T) {
	ctx := t.Context()

	t.Run("minimal fields", func(t *testing.T) {
		// Create a model with only required fields
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check basic fields
		require.NoError(t, err)
		assert.Equal(t, "test-integration", req.Name)
		assert.Equal(t, "postgres", req.Type)
	})

	t.Run("with connector_id", func(t *testing.T) {
		// Create a model with connector_id
		model := ResourceIntegrationModel{
			Name:        types.StringValue("test-integration"),
			Type:        types.StringValue("postgres"),
			ConnectorID: types.StringValue("test-connector-id"),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.Equal(t, "test-integration", req.Name)
		assert.Equal(t, "postgres", req.Type)
		assert.True(t, req.ConnectorID.IsSet())
		assert.Equal(t, "test-connector-id", req.ConnectorID.Value)
	})

	t.Run("with connected resource types", func(t *testing.T) {
		// Create a model with connected resource types
		resourceTypes := []attr.Value{
			types.StringValue("database"),
			types.StringValue("schema"),
		}
		model := ResourceIntegrationModel{
			Name:                   types.StringValue("test-integration"),
			Type:                   types.StringValue("postgres"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, resourceTypes),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.ConnectedResourceTypes.IsSet())
		assert.Equal(t, []string{"database", "schema"}, req.ConnectedResourceTypes.Value)
	})

	t.Run("with integration config", func(t *testing.T) {
		// Create a model with integration config
		configMap := map[string]attr.Value{
			"host":     types.StringValue("localhost"),
			"port":     types.StringValue("5432"),
			"database": types.StringValue("postgres"),
		}
		model := ResourceIntegrationModel{
			Name:              types.StringValue("test-integration"),
			Type:              types.StringValue("postgres"),
			IntegrationConfig: types.MapValueMust(types.StringType, configMap),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
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
		// Create a model with AWS secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			SecretStoreConfig: &SecretStoreConfig{
				AWS: &AWSSecretConfig{
					Region:   types.StringValue("us-east-1"),
					SecretID: types.StringValue("secret-id"),
				},
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.AWS.IsSet())
		assert.Equal(t, "us-east-1", secretConfig.AWS.Value.Region)
		assert.Equal(t, "secret-id", secretConfig.AWS.Value.SecretID)
	})

	t.Run("with GCP secret store config", func(t *testing.T) {
		// Create a model with GCP secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			SecretStoreConfig: &SecretStoreConfig{
				GCP: &GCPSecretConfig{
					Project:  types.StringValue("my-project"),
					SecretID: types.StringValue("secret-id"),
				},
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Gcp.IsSet())
		assert.Equal(t, "my-project", secretConfig.Gcp.Value.Project)
		assert.Equal(t, "secret-id", secretConfig.Gcp.Value.SecretID)
	})

	t.Run("with Azure secret store config", func(t *testing.T) {
		// Create a model with Azure secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			SecretStoreConfig: &SecretStoreConfig{
				Azure: &AzureSecretConfig{
					VaultURL: types.StringValue("https://myvault.vault.azure.net"),
					Name:     types.StringValue("secret-name"),
				},
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Azure.IsSet())
		assert.Equal(t, "https://myvault.vault.azure.net", secretConfig.Azure.Value.VaultURL)
		assert.Equal(t, "secret-name", secretConfig.Azure.Value.Name)
	})

	t.Run("with HashiCorp Vault secret store config", func(t *testing.T) {
		// Create a model with HashiCorp Vault secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			SecretStoreConfig: &SecretStoreConfig{
				HashicorpVault: &HashicorpVaultConfig{
					SecretEngine: types.StringValue("kv"),
					Path:         types.StringValue("secret/data/postgres"),
				},
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.HashicorpVault.IsSet())
		assert.Equal(t, "kv", secretConfig.HashicorpVault.Value.SecretEngine)
		assert.Equal(t, "secret/data/postgres", secretConfig.HashicorpVault.Value.Path)
	})

	t.Run("with custom access details", func(t *testing.T) {
		// Create a model with custom access details
		model := ResourceIntegrationModel{
			Name:                types.StringValue("test-integration"),
			Type:                types.StringValue("postgres"),
			CustomAccessDetails: types.StringValue("Use your SSO credentials to login"),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.CustomAccessDetails.IsSet())
		assert.Equal(t, "Use your SSO credentials to login", req.CustomAccessDetails.Value)
	})

	t.Run("with cleanup periods", func(t *testing.T) {
		// Create a model with cleanup periods
		model := ResourceIntegrationModel{
			Name:                            types.StringValue("test-integration"),
			Type:                            types.StringValue("postgres"),
			UserCleanupPeriodInDays:         types.Int64Value(30),
			CredentialsRotationPeriodInDays: types.Int64Value(90),
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.UserCleanupPeriodInDays.IsSet())
		assert.Equal(t, int64(30), req.UserCleanupPeriodInDays.Value)
		assert.True(t, req.CredentialsRotationPeriodInDays.IsSet())
		assert.Equal(t, int64(90), req.CredentialsRotationPeriodInDays.Value)
	})

	t.Run("with owner config", func(t *testing.T) {
		// Create a model with owner config
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			Owner: &OwnerConfig{
				SourceIntegrationName: types.StringValue("source-integration"),
				Type:                  types.StringValue("user"),
				Values:                types.ListValueMust(types.StringType, []attr.Value{types.StringValue("user1"), types.StringValue("user2")}),
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.Owner.IsSet())
		owner := req.Owner.Value
		assert.Equal(t, "user", owner.AttributeType)
		assert.Equal(t, []string{"user1", "user2"}, owner.AttributeValue)
		assert.True(t, owner.SourceIntegrationReference.IsSet())
		assert.Equal(t, "source-integration", owner.SourceIntegrationReference.Value)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		// Create a model with owners mapping
		model := ResourceIntegrationModel{
			Name: types.StringValue("test-integration"),
			Type: types.StringValue("postgres"),
			OwnersMapping: &OwnersMappingConfig{
				SourceIntegrationName: types.StringValue("source-integration"),
				KeyName:               types.StringValue("owner"),
				AttributeType:         types.StringValue("group"),
			},
		}

		// Call the function
		req, err := CreateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
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
		// Create a minimal integration response
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check basic fields
		require.NoError(t, err)
		assert.Equal(t, "integration-id", model.ID.ValueString())
		assert.Equal(t, "test-integration", model.Name.ValueString())
		assert.Equal(t, "postgres", model.Type.ValueString())
		assert.Equal(t, "connector-id", model.ConnectorID.ValueString())
		assert.True(t, model.ConnectedResourceTypes.IsNull())
		assert.True(t, model.IntegrationConfig.IsNull())
		assert.Nil(t, model.SecretStoreConfig)
		assert.True(t, model.CustomAccessDetails.IsNull())
		assert.True(t, model.UserCleanupPeriodInDays.IsNull())
		assert.True(t, model.CredentialsRotationPeriodInDays.IsNull())
		assert.Nil(t, model.Owner)
		assert.Nil(t, model.OwnersMapping)
	})

	t.Run("with connected resource types", func(t *testing.T) {
		// Create an integration with connected resource types
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			ConnectedResourceTypes: client.OptNilStringArray{
				Value: []string{"database", "schema"},
				Set:   true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.False(t, model.ConnectedResourceTypes.IsNull())

		var resourceTypes []string
		diags := model.ConnectedResourceTypes.ElementsAs(ctx, &resourceTypes, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"database", "schema"}, resourceTypes)
	})

	t.Run("with integration config", func(t *testing.T) {
		// Create an integration with config
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			IntegrationConfig: map[string]jx.Raw{
				"host":     jx.Raw("\"" + "localhost" + "\""),
				"port":     jx.Raw("\"" + "5432" + "\""),
				"database": jx.Raw("\"" + "postgres" + "\""),
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
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
		// Create an integration with AWS secret store
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			SecretStoreConfig: client.OptNilIntegrationV4SecretStoreConfig{
				Value: client.IntegrationV4SecretStoreConfig{
					AWS: client.OptNilIntegrationV4SecretStoreConfigAWS{
						Value: client.IntegrationV4SecretStoreConfigAWS{
							Region:   "us-east-1",
							SecretID: "secret-id",
						},
						Set: true,
					},
				},
				Set: true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
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
		// Create an integration with GCP secret store
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			SecretStoreConfig: client.OptNilIntegrationV4SecretStoreConfig{
				Value: client.IntegrationV4SecretStoreConfig{
					Gcp: client.OptNilIntegrationV4SecretStoreConfigGcp{
						Value: client.IntegrationV4SecretStoreConfigGcp{
							Project:  "my-project",
							SecretID: "secret-id",
						},
						Set: true,
					},
				},
				Set: true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
		require.NoError(t, err)
		require.NotNil(t, model.SecretStoreConfig)
		require.NotNil(t, model.SecretStoreConfig.GCP)
		assert.Equal(t, "my-project", model.SecretStoreConfig.GCP.Project.ValueString())
		assert.Equal(t, "secret-id", model.SecretStoreConfig.GCP.SecretID.ValueString())
	})

	t.Run("with custom access details and cleanup periods", func(t *testing.T) {
		// Create an integration with custom access details and cleanup periods
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			CustomAccessDetails: client.OptNilString{
				Value: "Use your SSO credentials to login",
				Set:   true,
			},
			UserCleanupPeriodInDays: client.OptNilInt64{
				Value: 30,
				Set:   true,
			},
			CredentialsRotationPeriodInDays: client.OptNilInt64{
				Value: 90,
				Set:   true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.Equal(t, "Use your SSO credentials to login", model.CustomAccessDetails.ValueString())
		assert.Equal(t, int64(30), model.UserCleanupPeriodInDays.ValueInt64())
		assert.Equal(t, int64(90), model.CredentialsRotationPeriodInDays.ValueInt64())
	})

	t.Run("with owner config", func(t *testing.T) {
		// Create an integration with owner config
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			Owner: client.OptNilIntegrationV4Owner{
				Value: client.IntegrationV4Owner{
					AttributeType:  "user",
					AttributeValue: []string{"user1", "user2"},
					SourceIntegrationName: client.OptNilString{
						Value: "source-integration",
						Set:   true,
					},
				},
				Set: true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
		require.NoError(t, err)
		require.NotNil(t, model.Owner)
		assert.Equal(t, "user", model.Owner.Type.ValueString())
		assert.Equal(t, "source-integration", model.Owner.SourceIntegrationName.ValueString())

		var values []string
		diags := model.Owner.Values.ElementsAs(ctx, &values, false)
		require.False(t, diags.HasError())
		assert.Equal(t, []string{"user1", "user2"}, values)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		// Create an integration with owners mapping
		integration := &client.IntegrationV4{
			ID:   "integration-id",
			Name: "test-integration",
			Type: "postgres",
			ConnectorID: client.OptNilString{
				Value: "connector-id",
				Set:   true,
			},
			OwnersMapping: client.OptNilIntegrationV4OwnersMapping{
				Value: client.IntegrationV4OwnersMapping{
					KeyName:       "owner",
					AttributeType: "group",
					SourceIntegrationName: client.OptNilString{
						Value: "source-integration",
						Set:   true,
					},
				},
				Set: true,
			},
		}

		// Call the function
		model, err := ResourceIntegrationToModel(ctx, integration)

		// Assert no errors and check fields
		require.NoError(t, err)
		require.NotNil(t, model.OwnersMapping)
		assert.Equal(t, "owner", model.OwnersMapping.KeyName.ValueString())
		assert.Equal(t, "group", model.OwnersMapping.AttributeType.ValueString())
		assert.Equal(t, "source-integration", model.OwnersMapping.SourceIntegrationName.ValueString())
	})
}

func TestUpdateIntegrationRequest(t *testing.T) {
	ctx := t.Context()

	t.Run("minimal fields", func(t *testing.T) {
		// Create a model with only name
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check basic fields
		require.NoError(t, err)
		assert.Equal(t, "updated-integration", req.Name)
	})

	t.Run("with integration config", func(t *testing.T) {
		// Create a model with integration config
		configMap := map[string]attr.Value{
			"host":     types.StringValue("new-host"),
			"port":     types.StringValue("5433"),
			"database": types.StringValue("test-db"),
		}
		model := ResourceIntegrationModel{
			Name:              types.StringValue("updated-integration"),
			IntegrationConfig: types.MapValueMust(types.StringType, configMap),
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
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
		// Create a model with connected resource types
		resourceTypes := []attr.Value{
			types.StringValue("table"),
			types.StringValue("view"),
		}
		model := ResourceIntegrationModel{
			Name:                   types.StringValue("updated-integration"),
			ConnectedResourceTypes: types.ListValueMust(types.StringType, resourceTypes),
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.ConnectedResourceTypes.IsSet())
		assert.Equal(t, []string{"table", "view"}, req.ConnectedResourceTypes.Value)
	})

	t.Run("with AWS secret store config", func(t *testing.T) {
		// Create a model with AWS secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			SecretStoreConfig: &SecretStoreConfig{
				AWS: &AWSSecretConfig{
					Region:   types.StringValue("us-west-2"),
					SecretID: types.StringValue("new-secret-id"),
				},
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.AWS.IsSet())
		assert.Equal(t, "us-west-2", secretConfig.AWS.Value.Region)
		assert.Equal(t, "new-secret-id", secretConfig.AWS.Value.SecretID)
	})

	t.Run("with GCP secret store config", func(t *testing.T) {
		// Create a model with GCP secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			SecretStoreConfig: &SecretStoreConfig{
				GCP: &GCPSecretConfig{
					Project:  types.StringValue("new-project"),
					SecretID: types.StringValue("new-secret-id"),
				},
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Gcp.IsSet())
		assert.Equal(t, "new-project", secretConfig.Gcp.Value.Project)
		assert.Equal(t, "new-secret-id", secretConfig.Gcp.Value.SecretID)
	})

	t.Run("with Azure secret store config", func(t *testing.T) {
		// Create a model with Azure secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			SecretStoreConfig: &SecretStoreConfig{
				Azure: &AzureSecretConfig{
					VaultURL: types.StringValue("https://newvault.vault.azure.net"),
					Name:     types.StringValue("new-secret-name"),
				},
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.Azure.IsSet())
		assert.Equal(t, "https://newvault.vault.azure.net", secretConfig.Azure.Value.VaultURL)
		assert.Equal(t, "new-secret-name", secretConfig.Azure.Value.Name)
	})

	t.Run("with HashiCorp Vault secret store config", func(t *testing.T) {
		// Create a model with HashiCorp Vault secret store config
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			SecretStoreConfig: &SecretStoreConfig{
				HashicorpVault: &HashicorpVaultConfig{
					SecretEngine: types.StringValue("kv2"),
					Path:         types.StringValue("secret/data/new-postgres"),
				},
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.SecretStoreConfig.IsSet())
		secretConfig := req.SecretStoreConfig.Value
		assert.True(t, secretConfig.HashicorpVault.IsSet())
		assert.Equal(t, "kv2", secretConfig.HashicorpVault.Value.SecretEngine)
		assert.Equal(t, "secret/data/new-postgres", secretConfig.HashicorpVault.Value.Path)
	})

	t.Run("with custom access details", func(t *testing.T) {
		// Create a model with custom access details
		model := ResourceIntegrationModel{
			Name:                types.StringValue("updated-integration"),
			CustomAccessDetails: types.StringValue("Updated access instructions"),
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.CustomAccessDetails.IsSet())
		assert.Equal(t, "Updated access instructions", req.CustomAccessDetails.Value)
	})

	t.Run("with cleanup periods", func(t *testing.T) {
		// Create a model with cleanup periods
		model := ResourceIntegrationModel{
			Name:                            types.StringValue("updated-integration"),
			UserCleanupPeriodInDays:         types.Int64Value(45),
			CredentialsRotationPeriodInDays: types.Int64Value(120),
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.UserCleanupPeriodInDays.IsSet())
		assert.Equal(t, int64(45), req.UserCleanupPeriodInDays.Value)
		assert.True(t, req.CredentialsRotationPeriodInDays.IsSet())
		assert.Equal(t, int64(120), req.CredentialsRotationPeriodInDays.Value)
	})

	t.Run("with owner config", func(t *testing.T) {
		// Create a model with owner config
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			Owner: &OwnerConfig{
				SourceIntegrationName: types.StringValue("new-source-integration"),
				Type:                  types.StringValue("group"),
				Values:                types.ListValueMust(types.StringType, []attr.Value{types.StringValue("group1"), types.StringValue("group2")}),
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.Owner.IsSet())
		owner := req.Owner.Value
		assert.Equal(t, "group", owner.AttributeType)
		assert.Equal(t, []string{"group1", "group2"}, owner.AttributeValue)
		assert.True(t, owner.SourceIntegrationReference.IsSet())
		assert.Equal(t, "new-source-integration", owner.SourceIntegrationReference.Value)
	})

	t.Run("with owners mapping", func(t *testing.T) {
		// Create a model with owners mapping
		model := ResourceIntegrationModel{
			Name: types.StringValue("updated-integration"),
			OwnersMapping: &OwnersMappingConfig{
				SourceIntegrationName: types.StringValue("new-source-integration"),
				KeyName:               types.StringValue("new-owner"),
				AttributeType:         types.StringValue("tag"),
			},
		}

		// Call the function
		req, err := UpdateIntegrationRequest(ctx, model)

		// Assert no errors and check fields
		require.NoError(t, err)
		assert.True(t, req.OwnersMapping.IsSet())
		mapping := req.OwnersMapping.Value
		assert.Equal(t, "new-owner", mapping.KeyName)
		assert.Equal(t, "tag", mapping.AttributeType)
		assert.True(t, mapping.SourceIntegrationReference.IsSet())
		assert.Equal(t, "new-source-integration", mapping.SourceIntegrationReference.Value)
	})
}

package models

import (
	"testing"
	"time"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/go-faster/jx"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserInformationIntegrationToModal(t *testing.T) {
	ctx := t.Context()

	t.Run("minimal fields", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:       "integration-id",
			Name:     "test-integration",
			Type:     "postgres",
			Category: "database",
			Status:   "connected",
		}

		model, err := UserInformationIntegrationToModal(ctx, integration)

		require.NoError(t, err)
		assert.Equal(t, "integration-id", model.ID.ValueString())
		assert.Equal(t, "test-integration", model.Name.ValueString())
		assert.Equal(t, "postgres", model.Type.ValueString())
		assert.Equal(t, "database", model.Category.ValueString())
		assert.Equal(t, "connected", model.Status.ValueString())
		assert.True(t, model.LastSyncTime.IsNull())
		assert.True(t, model.IntegrationConfig.IsNull())
		assert.Nil(t, model.SecretConfig)
	})

	t.Run("with last sync time", func(t *testing.T) {
		syncTime := time.Date(2023, 6, 15, 12, 30, 0, 0, time.UTC)
		integration := &client.IntegrationV4{
			ID:           "integration-id",
			Name:         "test-integration",
			Type:         "postgres",
			Category:     "database",
			Status:       "connected",
			LastSyncTime: client.NewOptNilApiInstant(client.ApiInstant(syncTime)),
		}

		model, err := UserInformationIntegrationToModal(ctx, integration)

		require.NoError(t, err)
		assert.Equal(t, "2023-06-15T12:30:00Z", model.LastSyncTime.ValueString())
	})

	t.Run("with integration config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:       "integration-id",
			Name:     "test-integration",
			Type:     "postgres",
			Category: "database",
			Status:   "connected",
			IntegrationConfig: map[string]jx.Raw{
				"host":     jx.Raw("\"" + "localhost" + "\""),
				"port":     jx.Raw("\"" + "5432" + "\""),
				"database": jx.Raw("\"" + "postgres" + "\""),
			},
		}

		model, err := UserInformationIntegrationToModal(ctx, integration)

		require.NoError(t, err)
		assert.False(t, model.IntegrationConfig.IsNull())

		elements := model.IntegrationConfig.Elements()
		hostValue, ok := elements["host"]
		require.True(t, ok)
		hostValueStr, ok := hostValue.(types.String)
		require.True(t, ok)
		assert.Equal(t, "localhost", hostValueStr.ValueString())

		portValue, ok := elements["port"]
		require.True(t, ok)
		portValueStr, ok := portValue.(types.String)
		require.True(t, ok)
		assert.Equal(t, "5432", portValueStr.ValueString())

		dbValue, ok := elements["database"]
		require.True(t, ok)
		dbValueStr, ok := dbValue.(types.String)
		require.True(t, ok)
		assert.Equal(t, "postgres", dbValueStr.ValueString())
	})

	t.Run("with AWS secret store config", func(t *testing.T) {
		integration := &client.IntegrationV4{
			ID:       "integration-id",
			Name:     "test-integration",
			Type:     "postgres",
			Category: "database",
			Status:   "connected",
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

		model, err := UserInformationIntegrationToModal(ctx, integration)

		require.NoError(t, err)
		require.NotNil(t, model.SecretConfig)
		require.NotNil(t, model.SecretConfig.AWS)
		assert.Equal(t, "us-east-1", model.SecretConfig.AWS.Region.ValueString())
		assert.Equal(t, "secret-id", model.SecretConfig.AWS.SecretID.ValueString())
		assert.Nil(t, model.SecretConfig.GCP)
		assert.Nil(t, model.SecretConfig.Azure)
		assert.Nil(t, model.SecretConfig.HashicorpVault)
	})
}

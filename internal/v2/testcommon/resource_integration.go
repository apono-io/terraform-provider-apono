package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/go-faster/jx"
)

func GenerateResourceIntegrationResponse() *client.IntegrationV4 {
	integration := client.IntegrationV4{
		ID:                     "integration-123",
		Name:                   "test-postgres-integration",
		Type:                   "postgres",
		ConnectorID:            client.NewOptNilString("connector-id-123"),
		ConnectedResourceTypes: client.NewOptNilStringArray([]string{"database", "schema", "table"}),
		IntegrationConfig: map[string]jx.Raw{
			"host":     jx.Raw("\"" + "test-postgres.example.com" + "\""),
			"port":     jx.Raw("\"" + "5432" + "\""),
			"database": jx.Raw("\"" + "postgres" + "\""),
			"ssl":      jx.Raw("\"" + "true" + "\""),
		},
		CustomAccessDetails: client.NewOptNilString("Use your company SSO credentials to login to this database"),
	}

	integration.SecretStoreConfig = client.NewOptNilSecretStoreConfigV4(
		client.SecretStoreConfigV4{
			AWS: client.NewOptNilAwsSecretConfigV4(
				client.AwsSecretConfigV4{
					Region:   "us-east-1",
					SecretID: "postgres/production",
				},
			),
		},
	)

	integration.Owner = client.NewOptNilOwnerV4(
		client.OwnerV4{
			AttributeType:         "user",
			AttributeValue:        []string{"admin@example.com", "dba@example.com"},
			SourceIntegrationName: client.NewOptNilString("Okta Directory"),
		},
	)

	integration.OwnersMapping = client.NewOptNilOwnerMappingV4(
		client.OwnerMappingV4{
			KeyName:               "owner",
			AttributeType:         "group",
			SourceIntegrationName: client.NewOptNilString("Okta Directory"),
		},
	)

	return &integration
}

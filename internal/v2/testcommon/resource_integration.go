package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/go-faster/jx"
)

func GenerateResourceIntegrationResponse() *client.IntegrationV4 {
	integration := client.IntegrationV4{
		ID:   "integration-123",
		Name: "test-postgres-integration",
		Type: "postgres",
		ConnectorID: client.OptNilString{
			Value: "connector-id-123",
			Set:   true,
		},
		ConnectedResourceTypes: client.OptNilStringArray{
			Value: []string{"database", "schema", "table"},
			Set:   true,
		},
		IntegrationConfig: map[string]jx.Raw{
			"host":     jx.Raw("\"" + "test-postgres.example.com" + "\""),
			"port":     jx.Raw("\"" + "5432" + "\""),
			"database": jx.Raw("\"" + "postgres" + "\""),
			"ssl":      jx.Raw("\"" + "true" + "\""),
		},
		CustomAccessDetails: client.OptNilString{
			Value: "Use your company SSO credentials to login to this database",
			Set:   true,
		},
	}

	// Set AWS secret store configuration
	integration.SecretStoreConfig = client.OptNilIntegrationV4SecretStoreConfig{
		Value: client.IntegrationV4SecretStoreConfig{
			AWS: client.OptNilIntegrationV4SecretStoreConfigAWS{
				Value: client.IntegrationV4SecretStoreConfigAWS{
					Region:   "us-east-1",
					SecretID: "postgres/production",
				},
				Set: true,
			},
		},
		Set: true,
	}

	// Set owner configuration
	integration.Owner = client.OptNilIntegrationV4Owner{
		Value: client.IntegrationV4Owner{
			AttributeType:  "user",
			AttributeValue: []string{"admin@example.com", "dba@example.com"},
			SourceIntegrationName: client.OptNilString{
				Value: "Okta Directory",
				Set:   true,
			},
		},
		Set: true,
	}

	// Set owners mapping configuration
	integration.OwnersMapping = client.OptNilIntegrationV4OwnersMapping{
		Value: client.IntegrationV4OwnersMapping{
			KeyName:       "owner",
			AttributeType: "group",
			SourceIntegrationName: client.OptNilString{
				Value: "Okta Directory",
				Set:   true,
			},
		},
		Set: true,
	}

	return &integration
}

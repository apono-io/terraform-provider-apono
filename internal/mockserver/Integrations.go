package mockserver

import (
	"encoding/json"
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
	"time"
)

var (
	DevConnectorId  = "dev-connector"
	ProdConnectorId = "prod-connector"
	MysqlType       = "mysql"
	PostgresqlType  = "postgresql"
)

func SetupMockHttpServerIntegrationV2Endpoints(existingIntegrations []apono.Integration) {
	var integrations = map[string]apono.Integration{}
	for _, integration := range existingIntegrations {
		integrations[integration.Id] = integration
	}

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/v2/integrations", func(req *http.Request) (*http.Response, error) {
		var createReq apono.CreateIntegration
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		integration := apono.Integration{
			Id:            id.String(),
			Name:          createReq.Name,
			Type:          createReq.Type,
			ProvisionerId: createReq.ProvisionerId,
			Status:        apono.INTEGRATIONSTATUS_ACTIVE,
			Metadata:      createReq.Metadata,
			SecretConfig:  createReq.SecretConfig,
		}
		integrations[integration.Id] = integration

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		integration, exists := integrations[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Integration not found"), nil
		}

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/integrations", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, apono.PaginatedResponseIntegrationModel{
			Data: existingIntegrations,
			Pagination: apono.PaginationInfo{
				Total:  int32(len(existingIntegrations)),
				Limit:  int32(len(existingIntegrations)),
				Offset: 0,
			},
		})
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPut, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq apono.UpdateIntegration
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		integration, exists := integrations[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Integration not found"), nil
		}

		integration.Name = updateReq.Name
		integration.ProvisionerId = updateReq.ProvisionerId
		integration.Metadata = updateReq.Metadata
		integration.SecretConfig = updateReq.SecretConfig

		integrations[integration.Id] = integration

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		delete(integrations, id)

		messageResponse := apono.MessageResponse{
			Message: "Deleted integration",
		}

		resp, err := httpmock.NewJsonResponse(200, messageResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v2/integrations-catalog/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		configType := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		switch configType {
		case "postgresql":
			config := apono.IntegrationConfig{
				Name:        "PostgreSQL",
				Type:        "postgresql",
				Description: "An open-source relational database management system emphasizing extensibility and SQL compliance.",
				Params: []apono.IntegrationConfigParam{
					{
						Id:    "hostname",
						Label: "Hostname",
					},
					{
						Id:      "port",
						Label:   "Port",
						Default: "5432",
					},
					{
						Id:      "dbname",
						Label:   "Database Name",
						Default: "postgres",
					},
				},
				RequiresSecret: true,
				SupportedSecretTypes: []string{
					"AWS",
					"GCP",
					"KUBERNETES",
				},
			}

			resp, err := httpmock.NewJsonResponse(200, config)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}

			return resp, nil
		default:
			return httpmock.NewStringResponse(400, "Unsupported config type"), nil
		}
	})
}

func CreateMockIntegrations() []apono.Integration {
	details := "4 resources loaded"
	return []apono.Integration{
		{
			Id:            uuid.NewString(),
			Name:          "MySQL DEV",
			Type:          MysqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&DevConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata: map[string]interface{}{
				"aws_account_id": "0123456789",
			},
			SecretConfig:           map[string]interface{}{},
			ConnectedResourceTypes: []string{"mysql-cluster", "mysql-db"},
		},
		{
			Id:            uuid.NewString(),
			Name:          "Postgres DEV",
			Type:          PostgresqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&DevConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata: map[string]interface{}{
				"hostname": "rds.amazon.example.com",
				"port":     "4560",
			},
			SecretConfig: map[string]interface{}{
				"type":      "AWS",
				"region":    "us-east-1",
				"secret_id": "my-secret-id",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-db"},
		},
		{
			Id:            uuid.NewString(),
			Name:          "MySQL PROD",
			Type:          MysqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&ProdConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata:      map[string]interface{}{},
			SecretConfig: map[string]interface{}{
				"type":      "GCP",
				"project":   "my-project-id",
				"secret_id": "my-secret-id",
			},
			ConnectedResourceTypes: []string{"mysql-cluster", "mysql-db"},
		},
		{
			Id:            uuid.NewString(),
			Name:          "Postgresql PROD",
			Type:          PostgresqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&ProdConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata:      nil,
			SecretConfig: map[string]interface{}{
				"type":      "KUBERNETES",
				"namespace": "prod",
				"name":      "postgres-credentials",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-db"},
		},
	}
}

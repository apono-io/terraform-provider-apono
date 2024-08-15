package mockserver

import (
	"encoding/json"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
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

func SetupMockHttpServerIntegrationTFV1Endpoints(existingIntegrations []aponoapi.IntegrationTerraform) {
	var integrations = map[string]aponoapi.IntegrationTerraform{}
	for _, integration := range existingIntegrations {
		integrations[integration.Id] = integration
	}

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/terraform/v1/integrations", func(req *http.Request) (*http.Response, error) {
		var createReq aponoapi.UpsertIntegrationTerraform
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}
		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		integration := aponoapi.IntegrationTerraform{
			Id:                     id.String(),
			Name:                   createReq.Name,
			Type:                   createReq.Type,
			Status:                 aponoapi.INTEGRATIONSTATUS_ACTIVE,
			ProvisionerId:          createReq.ProvisionerId,
			Params:                 createReq.Params,
			SecretConfig:           createReq.SecretConfig,
			ConnectedResourceTypes: createReq.ConnectedResourceTypes,
			CustomAccessDetails:    *aponoapi.NewNullableString(createReq.CustomAccessDetails.Get()),
			IntegrationOwners:      createReq.IntegrationOwners,
			ResourceOwnersMappings: createReq.ResourceOwnersMappings,
		}

		integrations[integration.Id] = integration

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/terraform/v1/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
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

	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/terraform/v1/integrations", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, aponoapi.PaginatedResponseIntegrationTerraformModel{
			Data: existingIntegrations,
			Pagination: aponoapi.PaginationInfo{
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

	httpmock.RegisterResponder(http.MethodPut, `=~^http://api\.apono\.dev/api/terraform/v1/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq aponoapi.UpsertIntegrationTerraform
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		integration, exists := integrations[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Integration not found"), nil
		}

		integration.Name = updateReq.Name
		integration.ProvisionerId = updateReq.ProvisionerId
		integration.Params = updateReq.Params
		integration.SecretConfig = updateReq.SecretConfig
		integration.CustomAccessDetails = *aponoapi.NewNullableString(updateReq.CustomAccessDetails.Get())
		integration.ConnectedResourceTypes = updateReq.ConnectedResourceTypes
		integration.IntegrationOwners = updateReq.IntegrationOwners
		integration.ResourceOwnersMappings = updateReq.ResourceOwnersMappings

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
}

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
			Id:                  id.String(),
			Name:                createReq.Name,
			Type:                createReq.Type,
			ProvisionerId:       createReq.ProvisionerId,
			Status:              apono.INTEGRATIONSTATUS_ACTIVE,
			Metadata:            createReq.Metadata,
			SecretConfig:        createReq.SecretConfig,
			CustomAccessDetails: createReq.CustomAccessDetails,
		}
		if createReq.ConnectedResourceTypes != nil {
			integration.ConnectedResourceTypes = createReq.ConnectedResourceTypes
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
		integration.CustomAccessDetails = updateReq.CustomAccessDetails
		if updateReq.ConnectedResourceTypes != nil {
			integration.ConnectedResourceTypes = updateReq.ConnectedResourceTypes
		}

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
}

func SetupMockHttpServerIntegrationCatalogEndpoints() {
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
					"HASHICORP_VAULT",
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
	customAccessDetails := "Please dont forget to save the secret"
	return []apono.Integration{
		{
			Id:            "1",
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
			CustomAccessDetails:    *apono.NewNullableString(&customAccessDetails),
		},
		{
			Id:            "2",
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
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
		{
			Id:            "3",
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
			CustomAccessDetails:    *apono.NewNullableString(&customAccessDetails),
		},
		{
			Id:            "4",
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
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
		{
			Id:            "5",
			Name:          "Postgresql PROD with Vault",
			Type:          PostgresqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&ProdConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata:      nil,
			SecretConfig: map[string]interface{}{
				"type":          "HASHICORP_VAULT",
				"secret_engine": "prod",
				"path":          "postgres-credentials",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
	}
}

func CreateTFIntegrations() []aponoapi.IntegrationTerraform {
	customAccessDetails := "Please dont forget to save the secret"
	return []aponoapi.IntegrationTerraform{
		{
			Id:            "1",
			Name:          "MySQL DEV",
			Type:          MysqlType,
			Status:        "Active",
			ProvisionerId: *aponoapi.NewNullableString(&DevConnectorId),
			Params: map[string]interface{}{
				"aws_account_id": "0123456789",
			},
			SecretConfig:           map[string]interface{}{},
			ConnectedResourceTypes: []string{"mysql-cluster", "mysql-db"},
			CustomAccessDetails:    *aponoapi.NewNullableString(&customAccessDetails),
		},
		{
			Id:            "2",
			Name:          "Postgres DEV",
			Type:          PostgresqlType,
			Status:        "Active",
			ProvisionerId: *aponoapi.NewNullableString(&DevConnectorId),
			Params: map[string]interface{}{
				"hostname": "rds.amazon.example.com",
				"port":     "4560",
			},
			SecretConfig: map[string]interface{}{
				"type":      "AWS",
				"region":    "us-east-1",
				"secret_id": "my-secret-id",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
		{
			Id:            "3",
			Name:          "MySQL PROD",
			Type:          MysqlType,
			Status:        "Active",
			ProvisionerId: *aponoapi.NewNullableString(&ProdConnectorId),
			Params:        map[string]interface{}{},
			SecretConfig: map[string]interface{}{
				"type":      "GCP",
				"project":   "my-project-id",
				"secret_id": "my-secret-id",
			},
			ConnectedResourceTypes: []string{"mysql-cluster", "mysql-db"},
			CustomAccessDetails:    *aponoapi.NewNullableString(&customAccessDetails),
		},
		{
			Id:            "4",
			Name:          "Postgresql PROD",
			Type:          PostgresqlType,
			Status:        "Active",
			ProvisionerId: *aponoapi.NewNullableString(&ProdConnectorId),
			Params:        nil,
			SecretConfig: map[string]interface{}{
				"type":      "KUBERNETES",
				"namespace": "prod",
				"name":      "postgres-credentials",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
		{
			Id:            "5",
			Name:          "Postgresql PROD with Vault",
			Type:          PostgresqlType,
			Status:        "Active",
			ProvisionerId: *aponoapi.NewNullableString(&ProdConnectorId),
			Params:        nil,
			SecretConfig: map[string]interface{}{
				"type":          "HASHICORP_VAULT",
				"secret_engine": "prod",
				"path":          "postgres-credentials",
			},
			ConnectedResourceTypes: []string{"postgresql-cluster", "postgresql-database"},
		},
	}
}

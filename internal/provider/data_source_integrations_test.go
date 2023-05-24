package provider

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var (
	devConnectorId  = "dev-connector"
	prodConnectorId = "prod-connector"
	mysqlType       = "mysql"
	postgresqlType  = "postgresql"
)

func TestAccIntegrationsDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	integrations := createMockIntegrations()
	setupMockHttpServerIntegrationsDataSource(integrations)

	checks := createIntegrationsDataSourceChecks(integrations)

	typeFilterPath := "data.apono_integrations.type_filter"
	connectorFilterPath := "data.apono_integrations.connector_filter"
	typeAndConnectorFilterPath := "data.apono_integrations.type_and_connector_filter"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.#", "2"),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.0.name", "MySQL DEV"),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.0.type", mysqlType),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.1.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.1.type", mysqlType),
				),
			},
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.#", "2"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.connector_id", prodConnectorId),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.1.name", "Postgresql PROD"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.1.connector_id", prodConnectorId),
				),
			},
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.#", "1"),
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.0.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.0.type", mysqlType),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.connector_id", prodConnectorId),
				),
			},
		},
	})
}

func createIntegrationsDataSourceChecks(integrations []apono.Integration) []resource.TestCheckFunc {
	allPath := "data.apono_integrations.all"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(allPath, "integrations.#", strconv.Itoa(len(integrations))),
	}

	for i, integration := range integrations {
		checks = append(
			checks,
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.id", i), integration.Id),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.name", i), integration.Name),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.type", i), integration.Type),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.connector_id", i), *integration.ProvisionerId.Get()),
		)

		if integration.Metadata != nil {
			for key, val := range integration.Metadata {
				checks = append(
					checks,
					resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.metadata.%s", i, key), fmt.Sprintf("%s", val)),
				)
			}
		} else {
			checks = append(checks, resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.metadata", i)))
		}

		secret := integration.SecretConfig
		if secret == nil {
			continue
		}

		switch secret["type"] {
		case "AWS":
			checks = append(
				checks,
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.aws_secret.region", i), fmt.Sprintf("%s", secret["region"])),
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.aws_secret.secret_id", i), fmt.Sprintf("%s", secret["secret_id"])),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.gcp_secret", i)),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.kubernetes_secret", i)),
			)
		case "GCP":
			checks = append(
				checks,
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.gcp_secret.project", i), fmt.Sprintf("%s", secret["project"])),
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.gcp_secret.secret_id", i), fmt.Sprintf("%s", secret["secret_id"])),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.aws_secret", i)),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.kubernetes_secret", i)),
			)
		case "KUBERNETES":
			checks = append(
				checks,
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.kubernetes_secret.namespace", i), fmt.Sprintf("%s", secret["namespace"])),
				resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.kubernetes_secret.name", i), fmt.Sprintf("%s", secret["name"])),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.aws_secret", i)),
				resource.TestCheckNoResourceAttr(allPath, fmt.Sprintf("integrations.%d.gcp_secret", i)),
			)
		}
	}
	return checks
}

func testAccIntegrationsDataSourceConfig() string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

data "apono_integrations" "all" {
}

data "apono_integrations" "type_filter" {
  type = "%[1]s"
}

data "apono_integrations" "connector_filter" {
  connector_id = "%[2]s"
}

data "apono_integrations" "type_and_connector_filter" {
  type = "%[1]s"
  connector_id = "%[2]s"
}
`, mysqlType, prodConnectorId)
}

func setupMockHttpServerIntegrationsDataSource(integrations []apono.Integration) {
	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/integrations", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, apono.PaginatedResponseIntegrationModel{
			Data: integrations,
			Pagination: apono.PaginationInfo{
				Total:  int32(len(integrations)),
				Limit:  int32(len(integrations)),
				Offset: 0,
			},
		})
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func createMockIntegrations() []apono.Integration {
	details := "4 resources loaded"
	return []apono.Integration{
		{
			Id:            uuid.NewString(),
			Name:          "MySQL DEV",
			Type:          mysqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&devConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata: map[string]interface{}{
				"aws_account_id": "0123456789",
			},
			SecretConfig: map[string]interface{}{},
		},
		{
			Id:            uuid.NewString(),
			Name:          "Postgres DEV",
			Type:          postgresqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&devConnectorId),
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
		},
		{
			Id:            uuid.NewString(),
			Name:          "MySQL PROD",
			Type:          mysqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&prodConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata:      map[string]interface{}{},
			SecretConfig: map[string]interface{}{
				"type":      "GCP",
				"project":   "my-project-id",
				"secret_id": "my-secret-id",
			},
		},
		{
			Id:            uuid.NewString(),
			Name:          "Postgresql PROD",
			Type:          postgresqlType,
			Status:        "Active",
			Details:       *apono.NewNullableString(&details),
			ProvisionerId: *apono.NewNullableString(&prodConnectorId),
			Connection:    map[string]interface{}{},
			LastSyncTime:  *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Metadata:      nil,
			SecretConfig: map[string]interface{}{
				"type":      "KUBERNETES",
				"namespace": "prod",
				"name":      "postgres-credentials",
			},
		},
	}
}

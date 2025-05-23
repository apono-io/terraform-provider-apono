package provider

import (
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"strconv"
	"testing"
)

func TestAccIntegrationsDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	integrations := mockserver.CreateTFIntegrations()
	mockserver.SetupMockHttpServerIntegrationTFV1Endpoints(integrations)
	mockserver.SetupMockHttpServerIntegrationCatalogEndpoints()

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
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.0.type", mockserver.MysqlType),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.1.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(typeFilterPath, "integrations.1.type", mockserver.MysqlType),
				),
			},
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.#", "2"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.connector_id", mockserver.ProdConnectorId),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.1.name", "Postgresql PROD"),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.1.connector_id", mockserver.ProdConnectorId),
				),
			},
			{
				Config: testAccIntegrationsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.#", "1"),
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.0.name", "MySQL PROD"),
					resource.TestCheckResourceAttr(typeAndConnectorFilterPath, "integrations.0.type", mockserver.MysqlType),
					resource.TestCheckResourceAttr(connectorFilterPath, "integrations.0.connector_id", mockserver.ProdConnectorId),
				),
			},
		},
	})
}

func createIntegrationsDataSourceChecks(integrations []aponoapi.IntegrationTerraform) []resource.TestCheckFunc {
	allPath := "data.apono_integrations.all"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(allPath, "integrations.#", strconv.Itoa(len(integrations))),
	}

	for i, integration := range integrations {
		provisionerIdPtr := integration.ProvisionerId.Get()
		var provisionerIdVal string
		if provisionerIdPtr != nil {
			provisionerIdVal = *provisionerIdPtr
		} else {
			provisionerIdVal = ""
		}

		checks = append(
			checks,
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.id", i), integration.Id),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.name", i), integration.Name),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.type", i), integration.Type),
			resource.TestCheckResourceAttr(allPath, fmt.Sprintf("integrations.%d.connector_id", i), provisionerIdVal),
		)

		if integration.Params != nil {
			for key, val := range integration.Params {
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
`, mockserver.MysqlType, mockserver.ProdConnectorId)
}

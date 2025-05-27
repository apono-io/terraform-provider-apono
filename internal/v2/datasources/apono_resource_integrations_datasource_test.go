package datasources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoResourceIntegrationsDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-int")
	resourceName := "apono_resource_integration.test"
	dataSourceName := "data.apono_resource_integrations.test"
	integrationType := common.MockDuck
	connectorID := testcommon.GetTestConnectorID(t)
	resourceType := common.MockDuck
	customAccessDetails := "Example access instructions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoResourceIntegrationsDataSourceConfig(rName, integrationType, connectorID, resourceType, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", integrationType),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", customAccessDetails),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "integrations.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "integrations.*", map[string]string{
						"name":                  rName,
						"type":                  integrationType,
						"connector_id":          connectorID,
						"custom_access_details": customAccessDetails,
					}),
					resource.TestCheckResourceAttrSet(dataSourceName, "integrations.0.id"),
				),
			},
		},
	})
}

func testAccAponoResourceIntegrationsDataSourceConfig(name, integrationType, connectorID, resourceType, customAccessDetails string) string {
	return fmt.Sprintf(`
resource "apono_resource_integration" "test" {
  name                     = "%s"
  type                     = "%s"
  connector_id             = "%s"
  connected_resource_types = ["%s"]
  integration_config = {
    key = "value"
  }
  custom_access_details = "%s"
  secret_store_config = {
    aws = {
      region    = "us-east-1"
      secret_id = "test-secret-id"
    }
  }
}

data "apono_resource_integrations" "test" {
  name = apono_resource_integration.test.name
  type = apono_resource_integration.test.type
  connector_id = apono_resource_integration.test.connector_id
  depends_on = [apono_resource_integration.test]
}
`, name, integrationType, connectorID, resourceType, customAccessDetails)
}

package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoResourceIntegrationResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_resource_integration.test"
	integrationType := "mock-duck"
	// Fetch the first connector ID dynamically
	connector, err := testcommon.GetFirstConnectorV3(t)
	if err != nil {
		t.Fatalf("failed to get connector: %v", err)
	}
	connectorID := connector.ID
	resourceType := "mock-duck"
	customAccessDetails := "Example access instructions"
	updatedCustomAccessDetails := "Updated access instructions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testcommon.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoResourceIntegrationConfig(rName, integrationType, connectorID, resourceType, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", integrationType),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connected_resource_types.0", resourceType),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", customAccessDetails),
					resource.TestCheckResourceAttr(resourceName, "integration_config.key", "value"),
				),
			},
			{
				Config: testAccAponoResourceIntegrationConfig(rName, integrationType, connectorID, resourceType, updatedCustomAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", updatedCustomAccessDetails),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAponoResourceIntegrationConfig(name, integrationType, connectorID, resourceType, customAccessDetails string) string {
	return fmt.Sprintf(`
resource "apono_resource_integration" "test" {
  name                    = "%s"
  type                    = "%s"
  connector_id            = "%s"
  connected_resource_types = ["%s"]
  integration_config = {
    key = "value"
  }
  custom_access_details = "%s"
}
`, name, integrationType, connectorID, resourceType, customAccessDetails)
}

package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAponoResourceIntegration(t *testing.T) {
	t.Run("ResourceIntegration", func(t *testing.T) {
		testAccAponoResourceIntegrationResource(t)
	})
	t.Run("ResourceIntegrationDrift", func(t *testing.T) {
		testAccAponoResourceIntegrationResourceDrift(t)
	})
}

func testAccAponoResourceIntegrationConfig(name, connectorID, customAccessDetails string) string {
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
  
  secret_store_config = {
    aws = {
      region = "us-east-1"
      secret_id = "test-secret-id"
    }
  }
}
`, name, common.MockDuck, connectorID, common.MockDuck, customAccessDetails)
}

func testAccAponoResourceIntegrationResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_resource_integration.test"

	connectorID := testcommon.GetTestConnectorID(t)

	customAccessDetails := "Example access instructions"
	updatedCustomAccessDetails := "Updated access instructions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoResourceIntegrationConfig(rName, connectorID, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connected_resource_types.0", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", customAccessDetails),
					resource.TestCheckResourceAttr(resourceName, "integration_config.key", "value"),
				),
			},
			{
				Config: testAccAponoResourceIntegrationConfig(rName, connectorID, updatedCustomAccessDetails),
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

func testAccAponoResourceIntegrationResourceDrift(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test-drift")
	resourceName := "apono_resource_integration.test"

	connectorID := testcommon.GetTestConnectorID(t)

	customAccessDetails := "Example access instructions"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoResourceIntegrationConfig(rName, connectorID, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connected_resource_types.0", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", customAccessDetails),
					resource.TestCheckResourceAttr(resourceName, "integration_config.key", "value"),
				),
			},
			// Delete the resource via API to simulate drift
			{
				Config: testAccAponoResourceIntegrationConfig(rName, connectorID, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					testAccDeleteResourceIntegrationViaAPI(t, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			// Apply again to recreate the resource
			{
				Config: testAccAponoResourceIntegrationConfig(rName, connectorID, customAccessDetails),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "connector_id", connectorID),
					resource.TestCheckResourceAttr(resourceName, "connected_resource_types.0", common.MockDuck),
					resource.TestCheckResourceAttr(resourceName, "custom_access_details", customAccessDetails),
					resource.TestCheckResourceAttr(resourceName, "integration_config.key", "value"),
				),
			},
		},
	})
}

// testAccDeleteResourceIntegrationViaAPI is a test helper that deletes the resource integration via API.
func testAccDeleteResourceIntegrationViaAPI(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has no ID set", resourceName)
		}

		// Get the API client from the provider
		clientInvoker := testcommon.GetTestClient(t)

		// Delete the integration via API
		err := clientInvoker.DeleteIntegrationV4(t.Context(), client.DeleteIntegrationV4Params{
			ID: rs.Primary.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete integration via API: %v", err)
		}

		return nil
	}
}

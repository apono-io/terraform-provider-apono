package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoBundleV2Resource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_bundle_v2.test"
	updatedName := acctest.RandomWithPrefix("tf-acc-updated")

	integrationType := "mock-duck"
	resourceType := "mock-duck"

	connectorID := testcommon.GetTestConnectorID(t)

	testAccAponoBundleV2Config := func(name string) string {
		return fmt.Sprintf(`
resource "apono_resource_integration" "test" {
  name                    = "%s-integration"
  type                    = "%s"
  connector_id            = "%s"
  connected_resource_types = ["%s"]
  integration_config = {
    key = "value"
  }
  custom_access_details = "Example access instructions"
  secret_store_config = {
    aws = {
      region = "us-east-1"
      secret_id = "test-secret-id"
    }
  }
}

resource "apono_access_scope" "test" {
  name  = "%s-scope"
  query = <<EOT
  integration = "%s" and resource_type = "%s"
  EOT
}

resource "apono_bundle_v2" "test" {
  name = "%s"
  
  access_targets = [
    {
      integration = {
        integration_name = apono_resource_integration.test.name
        resource_type = "%s"
        permissions = ["read", "write"]
        resources_scopes = [
          {
            scope_mode = "include_resources"
            type = "NAME"
            values = ["resource1", "resource2"]
          }
        ]
      }
    },
    {
      access_scope = {
        name = apono_access_scope.test.name
      }
    }
  ]
}
`, name, integrationType, connectorID, resourceType, name, connectorID, resourceType, name, resourceType)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoBundleV2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "access_targets.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "access_targets.0.integration.integration_name"),
					resource.TestCheckResourceAttr(resourceName, "access_targets.0.integration.permissions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "access_targets.0.integration.resources_scopes.0.scope_mode", "include_resources"),
					resource.TestCheckResourceAttr(resourceName, "access_targets.0.integration.resources_scopes.0.values.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "access_targets.1.access_scope.name"),
				),
			},
			{
				Config: testAccAponoBundleV2Config(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "access_targets.#", "2"),
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

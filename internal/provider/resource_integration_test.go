package provider

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/jarcoal/httpmock"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockserver.SetupMockHttpServerIntegrationV2Endpoints(make([]apono.Integration, 0))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIntegrationResourceConfig("integration-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apono_integration.test", "id"),
					resource.TestCheckResourceAttr("apono_integration.test", "name", "integration-name"),
					resource.TestCheckResourceAttr("apono_integration.test", "type", "postgresql"),
					resource.TestCheckResourceAttr("apono_integration.test", "aws_secret.region", "us-east-1"),
					resource.TestCheckResourceAttr("apono_integration.test", "connected_resource_types.#", "1"),
					resource.TestCheckTypeSetElemAttr("apono_integration.test", "connected_resource_types.*", "postgres-database"),
					resource.TestCheckResourceAttr("apono_integration.test", "custom_access_details", "This is a custom access detail"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apono_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationResourceConfig("updated-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apono_integration.test", "name", "updated-name"),
					resource.TestCheckResourceAttr("apono_integration.test", "type", "postgresql"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIntegrationResourceConfig(integrationName string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_integration" "test" {
  name = "%[1]s"
  type = "postgresql"
  connector_id = "000-1111-222222-33333-444444"
  connected_resource_types = ["postgres-database"]
  custom_access_details    = "This is a custom access detail"
  metadata = {
    hostname = "my-postgres-rds.aaabbbsss111.us-east-1.rds.amazonaws.com"
    port = "5432"
    dbname = "postgres"
  }
  aws_secret = {
    region = "us-east-1"
    secret_id = "my-secret"
  }
}
`, integrationName)
}

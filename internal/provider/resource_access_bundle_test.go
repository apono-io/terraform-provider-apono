package provider

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestAccAccessBundleResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	integrations := mockserver.CreateMockIntegrations()
	mockserver.SetupMockHttpServerIntegrationV2Endpoints(integrations)
	mockserver.SetupMockHttpServerAccessBundleV1Endpoints(make([]apono.AccessBundleV1, 0))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAccessBundleResourceConfig("access-bundle-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apono_access_bundle.test_access_bundle_resource", "id"),
					resource.TestCheckResourceAttr("apono_access_bundle.test_access_bundle_resource", "name", "access-bundle-name"),
					resource.TestCheckTypeSetElemNestedAttrs("apono_access_bundle.test_access_bundle_resource", "integration_targets.*", map[string]string{
						"name":          "Postgres DEV",
						"resource_type": "postgresql-db",
					}),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apono_access_bundle.test_access_bundle_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAccessBundleResourceConfig("updated-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apono_access_bundle.test_access_bundle_resource", "name", "updated-name"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAccessBundleResourceConfig(accessBundleName string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_access_bundle" "test_access_bundle_resource" {
  name = "%[1]s"
integration_targets = [
    {
      name = "Postgres DEV"
      resource_type = "postgresql-db"
      resource_include_filter = [[
        {
          type = "id"
          value = "12345"
        },
        {
          type = "name"
          value = "cluster2"
        },
        {
          type = "tag"
          name = "env"
          value = "prod"
        }
      ]]
      permissions = ["ReadOnly","ReadWrite","Admin"]
    },
	{
      name = "MySQL PROD"
      resource_type = "mysql-cluster"
      permissions = ["Admin"]
    }
  ]
}
`, accessBundleName)
}

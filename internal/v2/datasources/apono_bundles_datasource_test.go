package datasources_test

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoBundlesDataSource(t *testing.T) {
	rName1 := acctest.RandomWithPrefix("tf-acc-bundle-a")
	rName2 := acctest.RandomWithPrefix("tf-acc-bundle-b")
	dataSourceNameExact := "data.apono_bundles.exact"
	dataSourceNameWildcard := "data.apono_bundles.wildcard"

	randomPrefix := acctest.RandomWithPrefix("tf-acc-bundle-prefix")
	connectorID := testcommon.GetTestConnectorID(t)
	integrationType := common.MockDuck
	resourceType := common.MockDuck

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoBundlesDataSourceConfig(rName1, rName2, randomPrefix, connectorID, integrationType, resourceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "bundles.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "bundles.0.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "bundles.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "bundles.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "bundles.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "bundles.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName2),
					}),
					resource.TestCheckResourceAttrSet(dataSourceNameWildcard, "bundles.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceNameWildcard, "bundles.1.id"),
				),
			},
		},
	})
}

func testAccAponoBundlesDataSourceConfig(name1, name2, randomPrefix, connectorID, integrationType, resourceType string) string {
	prefixedName1 := testcommon.PrefixedName(randomPrefix, name1)
	prefixedName2 := testcommon.PrefixedName(randomPrefix, name2)
	integrationName := randomPrefix + "-integration"
	accessScopeName := randomPrefix + "-scope"

	return `
resource "apono_resource_integration" "test" {
  name                    = "` + integrationName + `"
  type                    = "` + integrationType + `"
  connector_id            = "` + connectorID + `"
  connected_resource_types = ["` + resourceType + `"]
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
  name = "` + accessScopeName + `"
  query = <<EOT
  integration = "` + integrationName + `" 
  EOT
}

resource "apono_bundle_v2" "test1" {
  name = "` + prefixedName1 + `"
  access_targets = [
    {
      integration = {
        integration_name = apono_resource_integration.test.name
        resource_type = "` + resourceType + `"
        permissions = ["read", "write"]
        resources_scopes = [
          {
            scope_mode = "include_resources"
            type = "NAME"
            values = ["resource1", "resource2"]
          }
        ]
      },
      access_scope = null
    },
    {
      integration = null,
      access_scope = {
        name = apono_access_scope.test.name
      }
    }
  ]
}

resource "apono_bundle_v2" "test2" {
  name = "` + prefixedName2 + `"
  access_targets = [
    {
      integration = {
        integration_name = apono_resource_integration.test.name
        resource_type = "` + resourceType + `"
        permissions = ["read", "write"]
        resources_scopes = [
          {
            scope_mode = "include_resources"
            type = "NAME"
            values = ["resource3"]
          }
        ]
      },
      access_scope = null
    }
  ]
}

data "apono_bundles" "exact" {
  name = "` + prefixedName1 + `"
  depends_on = [
    apono_bundle_v2.test1,
    apono_bundle_v2.test2
  ]
}

data "apono_bundles" "wildcard" {
  name = "` + randomPrefix + `*"
  depends_on = [
    apono_bundle_v2.test1,
    apono_bundle_v2.test2
  ]
}
`
}

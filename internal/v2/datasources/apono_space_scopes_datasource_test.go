package datasources_test

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoSpaceScopesDataSource(t *testing.T) {
	rName1 := acctest.RandomWithPrefix("tf-acc-test-a")
	rName2 := acctest.RandomWithPrefix("tf-acc-test-b")
	dataSourceNameExact := "data.apono_space_scopes.exact"
	dataSourceNameWildcard := "data.apono_space_scopes.wildcard"

	query := `integration in ("aws-account") and resource_tag["environment"] = "production"`

	randomPrefix := acctest.RandomWithPrefix("tf-acc-test-prefix")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoSpaceScopesDataSourceConfig(rName1, rName2, query, randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "space_scopes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "space_scopes.0.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "space_scopes.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "space_scopes.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "space_scopes.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "space_scopes.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName2),
					}),
				),
			},
		},
	})
}

func testAccAponoSpaceScopesDataSourceConfig(name1, name2, query, randomPrefix string) string {
	prefixedName1 := testcommon.PrefixedName(randomPrefix, name1)
	prefixedName2 := testcommon.PrefixedName(randomPrefix, name2)

	return `
resource "apono_space_scope" "test1" {
  name  = "` + prefixedName1 + `"
  query = <<EOT
  ` + query + `
  EOT
}

resource "apono_space_scope" "test2" {
  name  = "` + prefixedName2 + `"
  query = <<EOT
  ` + query + `
  EOT
}

data "apono_space_scopes" "exact" {
  name = "` + prefixedName1 + `"
  depends_on = [
    apono_space_scope.test1,
    apono_space_scope.test2
  ]
}

data "apono_space_scopes" "wildcard" {
  name = "` + randomPrefix + `*"
  depends_on = [
    apono_space_scope.test1,
    apono_space_scope.test2
  ]
}
`
}

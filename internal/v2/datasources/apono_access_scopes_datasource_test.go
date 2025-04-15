package datasources_test

import (
	"regexp"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoAccessScopesDataSource(t *testing.T) {
	// Create names that will sort predictably
	rName1 := acctest.RandomWithPrefix("tf-acc-test-a")
	rName2 := acctest.RandomWithPrefix("tf-acc-test-b")
	dataSourceNameExact := "data.apono_access_scopes.exact"
	dataSourceNameWildcard := "data.apono_access_scopes.wildcard"

	query := `integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"`

	// Create random prefix for wildcard match.
	randomPrefix := acctest.RandomWithPrefix("tf-acc-test-prefix")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testcommon.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoAccessScopesDataSourceConfig(rName1, rName2, query, randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "access_scopes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "access_scopes.0.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestMatchResourceAttr(dataSourceNameExact, "access_scopes.0.query",
						regexp.MustCompile(`(?s)^\s*integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"\s*$`)),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "access_scopes.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "access_scopes.#", "2"),
					resource.TestCheckResourceAttr(dataSourceNameWildcard, "access_scopes.*.name", testcommon.PrefixedName(randomPrefix, rName1)), // Since alphabetically sorted, rName1 comes first
					resource.TestCheckResourceAttr(dataSourceNameWildcard, "access_scopes.*.name", testcommon.PrefixedName(randomPrefix, rName2)),
					resource.TestMatchResourceAttr(dataSourceNameWildcard, "access_scopes.0.query",
						regexp.MustCompile(`(?s)^\s*integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"\s*$`)),
					resource.TestMatchResourceAttr(dataSourceNameWildcard, "access_scopes.1.query",
						regexp.MustCompile(`(?s)^\s*integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"\s*$`)),
				),
			},
		},
	})
}

func testAccAponoAccessScopesDataSourceConfig(name1, name2, query, randomPrefix string) string {
	prefixedName1 := testcommon.PrefixedName(randomPrefix, name1)
	prefixedName2 := testcommon.PrefixedName(randomPrefix, name2)

	return `
resource "apono_access_scope" "test1" {
  name  = "` + prefixedName1 + `"
  query = <<EOT
  ` + query + `
  EOT
}

resource "apono_access_scope" "test2" {
  name  = "` + prefixedName2 + `"
  query = <<EOT
  ` + query + `
  EOT
}

# Use depends_on to ensure resources are created before data sources are queried
data "apono_access_scopes" "exact" {
  name = "` + prefixedName1 + `"
  depends_on = [
    apono_access_scope.test1,
    apono_access_scope.test2
  ]
}

data "apono_access_scopes" "wildcard" {
  name = "` + randomPrefix + `*"
  depends_on = [
    apono_access_scope.test1,
    apono_access_scope.test2
  ]
}
`
}

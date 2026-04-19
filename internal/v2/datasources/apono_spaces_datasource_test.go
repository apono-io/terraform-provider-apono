package datasources_test

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoSpacesDataSource(t *testing.T) {
	rName1 := acctest.RandomWithPrefix("tf-acc-test-a")
	rName2 := acctest.RandomWithPrefix("tf-acc-test-b")
	dataSourceNameExact := "data.apono_spaces.exact"
	dataSourceNameWildcard := "data.apono_spaces.wildcard"

	randomPrefix := acctest.RandomWithPrefix("tf-acc-test-prefix")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoSpacesDataSourceConfig(rName1, rName2, randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "spaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "spaces.0.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "spaces.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "spaces.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "spaces.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "spaces.*", map[string]string{
						"name": testcommon.PrefixedName(randomPrefix, rName2),
					}),
				),
			},
		},
	})
}

func testAccAponoSpacesDataSourceConfig(name1, name2, randomPrefix string) string {
	prefixedName1 := testcommon.PrefixedName(randomPrefix, name1)
	prefixedName2 := testcommon.PrefixedName(randomPrefix, name2)

	return `
resource "apono_space_scope" "test_scope" {
  name  = "` + randomPrefix + `-scope"
  query = "integration in (\"aws-account\")"
}

resource "apono_space" "test1" {
  name = "` + prefixedName1 + `"
  space_scope_references = [
    apono_space_scope.test_scope.name,
  ]
}

resource "apono_space" "test2" {
  name = "` + prefixedName2 + `"
  space_scope_references = [
    apono_space_scope.test_scope.name,
  ]
}

data "apono_spaces" "exact" {
  name = "` + prefixedName1 + `"
  depends_on = [
    apono_space.test1,
    apono_space.test2
  ]
}

data "apono_spaces" "wildcard" {
  name = "` + randomPrefix + `*"
  depends_on = [
    apono_space.test1,
    apono_space.test2
  ]
}
`
}

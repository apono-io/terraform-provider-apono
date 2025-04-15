package datasources_test

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoGroupsDataSource(t *testing.T) {
	rName1 := acctest.RandomWithPrefix("tf-acc-test-a")
	rName2 := acctest.RandomWithPrefix("tf-acc-test-b")
	dataSourceNameExact := "data.apono_groups.exact"
	dataSourceNameWildcard := "data.apono_groups.wildcard"

	randomPrefix := acctest.RandomWithPrefix("tf-acc-test-prefix")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testcommon.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoGroupsDataSourceConfig(rName1, rName2, randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "groups.0.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "groups.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceNameWildcard, "groups.*.name", testcommon.PrefixedName(randomPrefix, rName1)),
					resource.TestCheckResourceAttr(dataSourceNameWildcard, "groups.*.name", testcommon.PrefixedName(randomPrefix, rName2)),
					resource.TestCheckResourceAttrSet(dataSourceNameWildcard, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceNameWildcard, "groups.1.id"),
				),
			},
		},
	})
}

func testAccAponoGroupsDataSourceConfig(name1, name2, randomPrefix string) string {
	prefixedName1 := testcommon.PrefixedName(randomPrefix, name1)
	prefixedName2 := testcommon.PrefixedName(randomPrefix, name2)

	return `
resource "apono_group" "test1" {
  name = "` + prefixedName1 + `"
  members = []
}

resource "apono_group" "test2" {
  name = "` + prefixedName2 + `"
  members = []
}

data "apono_groups" "exact" {
  name = "` + prefixedName1 + `"
  depends_on = [
    apono_group.test1,
    apono_group.test2
  ]
}

data "apono_groups" "wildcard" {
  name = "` + randomPrefix + `*"
  depends_on = [
    apono_group.test1,
    apono_group.test2
  ]
}
`
}

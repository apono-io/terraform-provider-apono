package datasources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoSpacesDataSource(t *testing.T) {
	// Keep names short — API enforces max 64 characters for space names.
	randomPrefix := acctest.RandString(8)
	spaceName1 := fmt.Sprintf("tf-sp-%s-a", randomPrefix)
	spaceName2 := fmt.Sprintf("tf-sp-%s-b", randomPrefix)
	scopeName := fmt.Sprintf("tf-sp-%s-scope", randomPrefix)
	dataSourceNameExact := "data.apono_spaces.exact"
	dataSourceNameWildcard := "data.apono_spaces.wildcard"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoSpacesDataSourceConfig(spaceName1, spaceName2, scopeName, randomPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNameExact, "spaces.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNameExact, "spaces.0.name", spaceName1),
					resource.TestCheckResourceAttrSet(dataSourceNameExact, "spaces.0.id"),

					resource.TestCheckResourceAttr(dataSourceNameWildcard, "spaces.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "spaces.*", map[string]string{
						"name": spaceName1,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameWildcard, "spaces.*", map[string]string{
						"name": spaceName2,
					}),
				),
			},
		},
	})
}

func testAccAponoSpacesDataSourceConfig(name1, name2, scopeName, randomPrefix string) string {
	return `
resource "apono_space_scope" "test_scope" {
  name  = "` + scopeName + `"
  query = "integration in (\"aws-account\")"
}

resource "apono_space" "test1" {
  name = "` + name1 + `"
  space_scope_references = [
    apono_space_scope.test_scope.name,
  ]
}

resource "apono_space" "test2" {
  name = "` + name2 + `"
  space_scope_references = [
    apono_space_scope.test_scope.name,
  ]
}

data "apono_spaces" "exact" {
  name = "` + name1 + `"
  depends_on = [
    apono_space.test1,
    apono_space.test2
  ]
}

data "apono_spaces" "wildcard" {
  name = "tf-sp-` + randomPrefix + `*"
  depends_on = [
    apono_space.test1,
    apono_space.test2
  ]
}
`
}

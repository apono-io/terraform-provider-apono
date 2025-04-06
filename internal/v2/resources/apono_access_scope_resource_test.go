package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoAccessScope_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_access_scope.test"

	query := `integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"`
	queryUpdate := `resource_type = "mock-duck"`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testcommon.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoAccessScopeConfig(rName, query),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "query", regexp.MustCompile(`(?s)^\s*integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"\s*$`)),
					resource.TestMatchResourceAttr(resourceName, "creation_date", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "update_date", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				),
			},
			{
				Config: testAccAponoAccessScopeConfig(rName, queryUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestMatchResourceAttr(resourceName, "query", regexp.MustCompile(`(?s)^\s*resource_type = "mock-duck"\s*$`)),
				),
			},
			{
				// Test import by ID
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test import by name instead of ID
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     rName, // Use the name instead of the ID for import
			},
		},
	})
}

func testAccAponoAccessScopeConfig(name, query string) string {
	return fmt.Sprintf(`
resource "apono_access_scope" "test" {
  name  = "%s"
  query = <<EOT
  %s
  EOT
}
`, name, query)
}

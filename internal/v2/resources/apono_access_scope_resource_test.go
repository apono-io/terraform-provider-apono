package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoAccessScopeResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_access_scope.test"

	query := `integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"`
	queryUpdate := `resource_type = "mock-duck"`
	description := "test description"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoAccessScopeConfig(rName, query, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestMatchResourceAttr(resourceName, "query", regexp.MustCompile(`(?s)^\s*integration = "5161d0f2-242d-42ee-92cb-8afd30caa0" and resource_type = "mock-duck"\s*$`)),
				),
			},
			{
				Config: testAccAponoAccessScopeConfig(rName, queryUpdate, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
					resource.TestMatchResourceAttr(resourceName, "query", regexp.MustCompile(`(?s)^\s*resource_type = "mock-duck"\s*$`)),
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

func testAccAponoAccessScopeConfig(name, query, description string) string {
	descriptionAttr := ""
	if description != "" {
		descriptionAttr = fmt.Sprintf(`  description = "%s"`, description)
	}
	return fmt.Sprintf(`
resource "apono_access_scope" "test" {
  name  = "%s"
%s
  query = <<EOT
  %s
  EOT
}
`, name, descriptionAttr, query)
}

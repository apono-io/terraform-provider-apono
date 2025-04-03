package resources_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"apono": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccAponoAccessScope_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_access_scope.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoAccessScopeConfig(rName, "tag:environment=dev"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "query", "tag:environment=dev"),
					resource.TestMatchResourceAttr(resourceName, "creation_date", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "update_date", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				),
			},
			{
				Config: testAccAponoAccessScopeConfig(rName, "tag:environment=prod"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "query", "tag:environment=prod"),
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

func testAccPreCheck(t *testing.T) {
	// Check for required environment variables
	if v := os.Getenv("APONO_ENDPOINT"); v == "" {
		t.Fatal("APONO_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("APONO_PERSONAL_TOKEN"); v == "" {
		t.Fatal("APONO_PERSONAL_TOKEN must be set for acceptance tests")
	}
}

func testAccAponoAccessScopeConfig(name, query string) string {
	return fmt.Sprintf(`
resource "apono_access_scope" "test" {
  name  = "%s"
  query = "%s"
}
`, name, query)
}

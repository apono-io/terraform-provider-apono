package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoGroupResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_group.test"

	users, err := testcommon.GetUsers(t)
	if err != nil {
		t.Fatalf("Error getting test users: %v", err)
	}

	if len(users) < 2 {
		t.Fatalf("Not enough users available for testing, need at least 2, got %d", len(users))
	}

	updatedName := rName + "-updated"

	testAccAponoGroupConfig := func(name string, members []string) string {
		membersStr := ""
		for _, member := range members {
			membersStr += "\n    \"" + member + "\","
		}

		return fmt.Sprintf(`
resource "apono_group" "test" {
  name     = "%s"
  members  = [%s
  ]
}
`, name, membersStr)
	}

	testAccAponoGroupConfigUpdated := func(name string, members []string) string {
		membersStr := ""
		for _, member := range members {
			membersStr += "\n    \"" + member + "\","
		}

		return fmt.Sprintf(`
resource "apono_group" "test" {
  name     = "%s"
  members  = [%s
  ]
}
`, name, membersStr)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoGroupConfig(rName, []string{users[0].Email}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "members.*", users[0].Email),
				),
			},
			{
				Config: testAccAponoGroupConfigUpdated(updatedName, []string{users[0].Email, users[1].Email}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "members.*", users[0].Email),
					resource.TestCheckTypeSetElemAttr(resourceName, "members.*", users[1].Email),
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

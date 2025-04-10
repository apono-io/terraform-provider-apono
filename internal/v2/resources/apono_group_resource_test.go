package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoGroup_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_group.test"

	// Get real users from the system for testing
	users, err := testcommon.GetUsers(t)
	if err != nil {
		t.Errorf("Error getting test users: %v", err)
		return
	}

	if len(users) < 2 {
		t.Errorf("Not enough users available for testing, need at least 2, got %d", len(users))
		return
	}

	// Updated name for second step
	updatedName := rName + "-updated"

	// Creates the initial Terraform configuration for a group
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

	// Creates the updated Terraform configuration for a group
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
		ProtoV6ProviderFactories: testcommon.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Initial creation with one user
				Config: testAccAponoGroupConfig(rName, []string{users[0].Email}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "members.*", users[0].Email),
				),
			},
			{
				// Update name and add second user
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
				// Test import by ID
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

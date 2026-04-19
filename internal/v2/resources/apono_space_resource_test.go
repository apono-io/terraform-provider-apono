package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoSpaceResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_space.test"
	scopeResourceName := "apono_space_scope.test"

	users, err := testcommon.GetUsers(t)
	if err != nil {
		t.Fatalf("Error getting test users: %v", err)
	}

	if len(users) < 2 {
		t.Fatalf("Not enough users available for testing, need at least 2, got %d", len(users))
	}

	updatedName := rName + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create space without members
			{
				Config: testAccAponoSpaceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "space_scope_references.#", "1"),
					resource.TestCheckResourceAttrSet(scopeResourceName, "id"),
				),
			},
			// Step 2: Update name and add a member
			{
				Config: testAccAponoSpaceConfigWithMembers(updatedName, []testSpaceMember{
					{IdentityReference: users[0].Email, IdentityType: "user", SpaceRoles: []string{"SpaceOwner"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "space_scope_references.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
				),
			},
			// Step 3: Update members (add a second member)
			{
				Config: testAccAponoSpaceConfigWithMembers(updatedName, []testSpaceMember{
					{IdentityReference: users[0].Email, IdentityType: "user", SpaceRoles: []string{"SpaceOwner"}},
					{IdentityReference: users[1].Email, IdentityType: "user", SpaceRoles: []string{"SpaceManager"}},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
				),
			},
			// Step 4: Remove members
			{
				Config: testAccAponoSpaceConfig(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckNoResourceAttr(resourceName, "members.#"),
				),
			},
			// Step 5: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testSpaceMember struct {
	IdentityReference string
	IdentityType      string
	SpaceRoles        []string
}

func testAccAponoSpaceConfig(name string) string {
	return fmt.Sprintf(`
resource "apono_space_scope" "test" {
  name  = "%[1]s-scope"
  query = "integration in (\"aws-account\")"
}

resource "apono_space" "test" {
  name = "%[1]s"
  space_scope_references = [
    apono_space_scope.test.name,
  ]
}
`, name)
}

func testAccAponoSpaceConfigWithMembers(name string, members []testSpaceMember) string {
	memberItems := ""
	for i, m := range members {
		rolesStr := ""
		for _, role := range m.SpaceRoles {
			rolesStr += fmt.Sprintf("        \"%s\",\n", role)
		}
		if i > 0 {
			memberItems += ",\n"
		}
		memberItems += fmt.Sprintf(`    {
      identity_reference = "%s"
      identity_type      = "%s"
      space_roles        = [
%s      ]
    }`, m.IdentityReference, m.IdentityType, rolesStr)
	}

	return fmt.Sprintf(`
resource "apono_space_scope" "test" {
  name  = "%[1]s-scope"
  query = "integration in (\"aws-account\")"
}

resource "apono_space" "test" {
  name = "%[1]s"
  space_scope_references = [
    apono_space_scope.test.name,
  ]
  members = [
%[2]s
  ]
}
`, name, memberItems)
}

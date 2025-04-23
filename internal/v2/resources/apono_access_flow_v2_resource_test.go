package resources_test

import (
	"fmt"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoAccessFlowV2Resource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "apono_access_flow_v2.test"
	updatedJustificationRequired := false

	var connectorID string

	integrationType := "mock-duck"
	resourceType := "mock-duck"

	if testcommon.IsTestAccount(t) {
		connectorID = "terraofrm-tests-account-connector"
	} else {
		connector, err := testcommon.GetFirstConnectorV3(t)
		if err != nil {
			t.Fatalf("failed to get connector: %v", err)
		}
		connectorID = connector.ID
	}

	users, err := testcommon.GetUsers(t)
	if err != nil {
		t.Fatalf("failed to get users: %v", err)
	}
	if len(users) == 0 {
		t.Fatalf("no users found for test")
	}
	userEmail := users[0].Email

	testAccAponoAccessFlowV2Config := func(name string, justificationRequired bool, userEmail string) string {
		return fmt.Sprintf(`
resource "apono_resource_integration" "test" {
  name                    = "%s-integration"
  type                    = "%s"
  connector_id            = "%s"
  connected_resource_types = ["%s"]
  integration_config = {
    key = "value"
  }
  custom_access_details = "Example access instructions"
  secret_store_config = {
    aws = {
      region = "us-east-1"
      secret_id = "test-secret-id"
    }
  }
}

resource "apono_access_flow_v2" "test" {
  name = "%s"
  trigger = "SELF_SERVE"
  active = true

  grantees = {
    logical_operator = "OR"
    conditions = [
      {
        type = "user"
        values = ["%s"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        name = apono_resource_integration.test.name
        integration_name = apono_resource_integration.test.name
        resources = [
          {
            type = "%s"
            resource_type = "%s"
            identifier = "example-db"
            permissions = ["read"]
          }
        ]
        permissions = ["read"]
        resource_type = "%s"
      }
    }
  ]

  settings = {
    justification_required = %t
    require_approver_reason = false
    requester_cannot_approve_self = true
    require_mfa = false
    labels = ["test", "example"]
  }
}
`, name, integrationType, connectorID, resourceType, name, userEmail, resourceType, resourceType, resourceType, justificationRequired)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoAccessFlowV2Config(rName, true, userEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "trigger", "SELF_SERVE"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "grantees.logical_operator", "OR"),
					resource.TestCheckResourceAttr(resourceName, "access_targets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "settings.justification_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "settings.labels.#", "2"),
				),
			},
			{
				Config: testAccAponoAccessFlowV2Config(rName, updatedJustificationRequired, userEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "settings.justification_required", "false"),
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

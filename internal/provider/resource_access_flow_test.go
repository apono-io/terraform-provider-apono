package provider

import (
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestAccAccessFlowResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	identities := mockserver.CreateMockIdentities()
	mockserver.SetupMockHttpServerIdentitiesV2Endpoints(identities)

	integrations := mockserver.CreateMockIntegrations()
	mockserver.SetupMockHttpServerIntegrationV2Endpoints(integrations)

	users := mockserver.CreateMockUsers()
	mockserver.SetupMockHttpServerUsersV2Endpoints(users)

	accessBundles := mockserver.CreateMockAccessBundles()
	mockserver.SetupMockHttpServerAccessBundleV1Endpoints(accessBundles)

	mockserver.SetupMockHttpServerAccessFlowV1Endpoints(make([]aponoapi.AccessFlowTerraformV1, 0))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccAccessFlowResourceConfig("access-flow-name", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apono_access_flow.test_access_flow_resource", "id"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "name", "access-flow-name"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "active", "true"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "settings.require_all_approvers", "true"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "grantees.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("apono_access_flow.test_access_flow_resource", "grantees.*", map[string]string{
						"type": "group",
						"name": "Test Group 1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("apono_access_flow.test_access_flow_resource", "integration_targets.*", map[string]string{
						"name":          "Postgres DEV",
						"resource_type": "postgresql-database",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("apono_access_flow.test_access_flow_resource", "approvers.*", map[string]string{
						"type": "context_attribute",
						"name": "Manager",
					}),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "labels.#", "2"),
					resource.TestCheckTypeSetElemAttr("apono_access_flow.test_access_flow_resource", "labels.*", "label1"),
					resource.TestCheckTypeSetElemAttr("apono_access_flow.test_access_flow_resource", "labels.*", "label2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apono_access_flow.test_access_flow_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAccessFlowResourceConfig("updated-name", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "name", "updated-name"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "grantees_filter_group.conditions_logical_operator", "AND"),
					resource.TestCheckResourceAttr("apono_access_flow.test_access_flow_resource", "grantees_filter_group.attribute_filters.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("apono_access_flow.test_access_flow_resource", "grantees_filter_group.attribute_filters.*", map[string]string{
						"attribute_type": "group",
						"operator":       "contains",
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAccessFlowResourceConfig(accessFlowName string, useOldGrantees bool) string {
	var grantees string
	if useOldGrantees {
		grantees = fmt.Sprintf(`
  grantees = [
    {
      name = "test1@example.com"
      type = "user"
    },
	{
      name = "Test Group 1"
      type = "group"
    }
  ]
`)
	} else {
		grantees = fmt.Sprintf(`
grantees_filter_group = {
    conditions_logical_operator = "AND"
    attribute_filters           = [
      {
        attribute_type  = "user"
        attribute_names = ["test1@example.com", "test2@example.com"]
      },
      {
        attribute_type  = "group"
        operator        = "contains"
        attribute_names = ["Test Group 1"]
      }
    ]
  }
`)
	}

	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_access_flow" "test_access_flow_resource" {
  name = "%[1]s"
  active = true
  revoke_after_in_sec = 3600
  trigger = {
    type = "user_request"
	timeframe = {
		  start_time = "00:00:00"
		  end_time = "23:59:59"
		  days_in_week = ["MONDAY","TUESDAY","WEDNESDAY","THURSDAY","FRIDAY"]
		  time_zone = "Asia/Jerusalem"
	}
  }
  %s
integration_targets = [
    {
      name = "Postgres DEV"
      resource_type = "postgresql-database"
      resource_include_filter = [[
        {
          type = "id"
          value = "12345"
        },
        {
          type = "name"
          value = "cluster2"
        },
        {
          type = "tag"
          name = "env"
          value = "prod"
        }
      ]]
      permissions = ["ReadOnly","ReadWrite","Admin"]
    },
	{
      name = "MySQL PROD"
      resource_type = "mysql-cluster"
      permissions = ["Admin"]
    }
  ]
bundle_targets = [
	{
		name = "DB PROD"
	}
  ]
approvers = [
    {
      name = "test2@example.com"
      type = "user"
    },
	{
      name = "Manager"
      type = "context_attribute"
    }
  ]
settings = {
    approver_cannot_self_approve = true
    require_all_approvers = true
  }
labels = ["label1","label2"]
}
`, accessFlowName, grantees)
}

package provider

import (
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestManualWebhookResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockserver.SetupMockHttpServerManualWebhookEndpoints(make([]aponoapi.WebhookManualTriggerTerraformModel, 0))

	manualWebhookName := "test_manual_webhook"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testManualWebhookResourceConfig(manualWebhookName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apono_manual_webhook.test", "id"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "name", manualWebhookName),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "type.http_request.url", "https://www.example.com"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "type.http_request.method", "GET"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "type.http_request.headers.X-Or-King", "true"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.json_path", "$.key"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.0", "value"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "timeout_in_sec", "10"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "custom_validation_error_message", "This is a custom validation error message"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apono_manual_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testManualWebhookResourceConfig("updated-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "name", "updated-name"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "type.http_request.url", "https://www.example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testManualWebhookResourceConfig(manualWebhookName string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_manual_webhook" "test" {
  name = "%[1]s"
  active = true
  type = {
    http_request = {
      url    = "https://www.example.com"
      method = "GET"
      headers = {
        "X-Or-King" = "true"
      }
    }
  }
  response_validators = [
    {
      json_path = "$.key"
      expected_values = ["value"]
    }
  ]
  timeout_in_sec                  = 10
  custom_validation_error_message = "This is a custom validation error message"
}
`, manualWebhookName)
}

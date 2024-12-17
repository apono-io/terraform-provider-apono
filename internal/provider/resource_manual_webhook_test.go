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
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "url", "http://my-webhook.com"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "method", "POST"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "body_template", "{\"key\": \"value\"}"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.json_path", "$.key"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.0", "value"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "timeout_in_sec", "10"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.type", "oauth"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.client_id", "my-client-id"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.client_secret", "my-client-secret"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.token_endpoint_url", "http://my-token-endpoint.com"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.#", "2"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.0", "scope1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.1", "scope2"),
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
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "url", "http://my-webhook.com"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "method", "POST"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "headers.Content-Type", "application/json"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "body_template", "{\"key\": \"value\"}"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.json_path", "$.key"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "response_validators.0.expected_values.0", "value"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "timeout_in_sec", "10"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.#", "1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.type", "oauth"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.client_id", "my-client-id"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.client_secret", "my-client-secret"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.token_endpoint_url", "http://my-token-endpoint.com"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.#", "2"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.0", "scope1"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "authentication_config.0.oauth.scopes.1", "scope2"),
					resource.TestCheckResourceAttr("apono_manual_webhook.test", "custom_validation_error_message", "This is a custom validation error message"),
				),
			},
		},
	})
}

func testManualWebhookResourceConfig(manualWebhookName string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_manual_webhook" "test_webhook" {
  name = "%[1]s"
	  url = "https://my-webhook.com"
	  method = "POST"
	  headers = {
	    "Content-Type" = "application/json"
	  }
	  body_template = "{\"key\": \"value\"}"
	  response_validators = [
	    {
			      json_path = "$.key"
			      expected_values = ["value"]
	    }
	  ]
	  timeout_in_sec = 10
	  authentication_config {
	    type = "oauth"
	    oauth {
			      client_id = "my-client-id"
			      client_secret = "my-client-secret"
			      token_endpoint_url = "http://my-token-endpoint.com"
			      scopes = ["scope1", "scope2"]
	    }
	  }
	  custom_validation_error_message = "This is a custom validation error message"
}
`, manualWebhookName)
}

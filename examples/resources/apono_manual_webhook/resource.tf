resource "apono_manual_webhook" "webhook_prod" {
  name   = "Webhook Prod"
  active = true
  type = {
    http_request = {
      method = "GET"
      url    = "https://example.com"
      headers = {
        "X-Header-Name" = "header-value"
      }
    }
  }
  body_template = "{ \"key\": \"value\" }"
  authentication_config = {
    oauth = {
      client_id          = "client_id"
      client_secret      = "client_secret"
      scopes             = ["scope1", "scope2"]
      token_endpoint_url = "https://example.com/token"
    }
  }
  response_validators = [
    {
      expected_values = ["value1", "value2"]
      json_path       = "json_path"
    }
  ]
  timeout_in_sec                  = 60
  custom_validation_error_message = "This is a custom validation error message"
}
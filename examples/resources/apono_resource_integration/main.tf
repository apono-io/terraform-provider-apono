resource "apono_resource_integration" "example" {
  name                     = "example-integration"
  type                     = "mock-duck"
  connector_id             = "your-connector-id"
  connected_resource_types = ["mock-duck"]
  integration_config = {
    key = "value"
  }
  custom_access_details = "Example access instructions"
}

resource "apono_resource_integration" "gcp_sql_integration" {
  name         = "GCP SQL Databases"
  type         = "gcp-organization"
  connector_id = "GcpConnector-Example"
  connected_resource_types = [
    "gcp-organization-project"
  ]
  integration_config = {
    project_id = "gcp-production-project"
    location   = "us-central1"
  }
  secret_store_config = {
    gcp = {
      project   = "gcp-secrets-project"
      secret_id = "projects/1234567890/secrets/sql-db-credentials"
    }
  }
  owner = {
    attribute_type          = "user"
    attribute_value         = ["example@company.io"]
    source_integration_name = "Google Oauth"
  }
}

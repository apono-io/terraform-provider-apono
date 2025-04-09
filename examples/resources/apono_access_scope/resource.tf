resource "apono_access_scope" "production_databases" {
  name  = "production-databases"
  query = <<EOT
integration_id = "aws-integration-id" AND resource_type = "rds-instance" AND tags.Environment = "production"
  EOT
}

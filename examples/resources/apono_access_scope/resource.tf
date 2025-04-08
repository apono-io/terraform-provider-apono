
resource "apono_access_scope" "production_databases" {
  name  = "production-databases"
  query = <<EOT
  integration = "aws-integration-id" and 
  resource_type = "rds-instance" and 
  tag:Environment = "production"
  EOT
}

resource "apono_space_scope" "staging_databases" {
  name  = "Staging Databases"
  query = "resource_type = \"aws-rds-mysql-database\" and resource_name contains \"staging\""
}

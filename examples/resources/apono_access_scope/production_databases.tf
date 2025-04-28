resource "apono_access_scope" "production_databases" {
  name  = "production-mysql-dbs-us_east-1"
  query = <<EOT
  resource_type = "aws-rds-mysql-database" 
    and resource_name contains "prod" 
    and (resource_tag["region"] = "us-east-1") 
    and resource_risk_level = "3" 
    and permission_risk_level = "3"
  EOT
}

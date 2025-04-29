resource "apono_resource_integration" "postgresql_prod_dbs" {
  name         = "PostgreSQL Production Databases"
  type         = "postgresql"
  connector_id = "AwsConnector-ProdTeam-XYZ123"
  connected_resource_types = [
    "postgresql-database",
    "postgresql-table"
  ]
  integration_config = {
    hostname = "prod-postgresql.us-east-1.internal.example.com"
    port     = "5432"
    dbname   = "postgres"
    sslmode  = "disable"
  }
  secret_store_config = {
    aws = {
      region    = "us-east-1"
      secret_id = "arn:aws:secretsmanager:us-east-1:123456789012:secret:/prod/postgresql/apono"
    }
  }
  custom_access_details = "Please make sure to connect to VPN before accessing the DBs"
}

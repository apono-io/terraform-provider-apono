resource "apono_integration" "postgresql_prod" {
  name                     = "DB Prod"
  type                     = "postgresql"
  connector_id             = "00000-1111-222222-33333-444444"
  connected_resource_types = ["postgresql-database", "postgresql-table"]
  custom_access_details    = "This is a custom access detail"
  metadata = {
    hostname = "prod-postgresql.us-east-1.internal.example.com"
    port     = "5432"
    dbname   = "postgres"
  }
  aws_secret = {
    region    = "us-east-1"
    secret_id = "arn:aws:secretsmanager:us-east-1:123456789012:secret:/prod/postgresql/apono"
  }
}

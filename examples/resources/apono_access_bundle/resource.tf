resource "apono_access_bundle" "prod" {
  name = "access to PROD DB"
  integration_targets = [
    {
      name          = "DB Prod"
      resource_type = "postgresql-database"
      permissions   = ["READ_ONLY", "READ_WRITE", "ADMIN"]
    }
  ]
}
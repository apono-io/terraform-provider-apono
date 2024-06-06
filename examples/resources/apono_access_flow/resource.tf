resource "apono_access_flow" "postgresql_prod" {
  name                = "access to PROD DB"
  active              = true
  revoke_after_in_sec = 3600
  trigger = {
    type = "user_request"
    timeframe = {
      start_time   = "00:00:00"
      end_time     = "23:59:59"
      days_in_week = ["MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"]
      time_zone    = "Asia/Jerusalem"
    }
  }
  grantees = [
    {
      name = "person@example.com"
      type = "user"
    },
    {
      name = "R&D Team"
      type = "group"
    }
  ]
  integration_targets = [
    {
      name          = "DB Prod"
      resource_type = "postgresql-database"
      permissions   = ["READ_ONLY", "READ_WRITE", "ADMIN"]
    }
  ]
  bundle_targets = [
    {
      name = "PROD ENV"
    }
  ]
  settings = {
    approver_cannot_self_approve   = true
  }
}
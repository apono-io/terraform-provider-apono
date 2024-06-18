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
  grantees_conditions_group = {
    conditions_logical_operator = "OR"
    attribute_conditions = [
      {
        attribute_type  = "user"
        attribute_names = ["person@example.com", "person_two@example.com"]
      },
      {
        attribute_type  = "group"
        operator        = "contains"
        attribute_names = ["R&D Team"]
      }
    ]
  }
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
    approver_cannot_self_approve = true
  }
  labels = ["DB", "PROD", "TERRAFORM"]
}
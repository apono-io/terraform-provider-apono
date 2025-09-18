resource "apono_access_flow_v2" "managers_request_for_employees" {
  name                  = "Sensitive Access to Production AWS"
  active                = true
  grant_duration_in_min = 60
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type                    = "group"
        match_operator          = "contains"
        source_integration_name = "Google Oauth"
        values                  = ["All Company"]
      }
    ]
  }

  request_for = {
    request_scopes = ["direct_reports"]
  }

  access_targets = [
    {
      integration = {
        integration_name = "Azure Subscription Integration"
        resource_type    = "azure-subscription-sql-server"
        permissions      = ["Contributor"]
      }
    }
  ]

  settings = {
    justification_required        = true
    require_approver_reason       = false
    requester_cannot_approve_self = false
    require_mfa                   = true
    labels                        = ["created_from_terraform"]
  }
}

resource "apono_access_flow_v2" "sensitive_production_aws" {
  name                  = "Sensitive Access to Production AWS"
  active                = true
  grant_duration_in_min = 60
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type                    = "user"
        match_operator          = "is"
        source_integration_name = data.apono_user_information_integrations.google_oauth_idp.integrations[0].name
        values                  = ["example@company.io"]
      }
    ]
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

  approver_policy = {
    approval_mode = "ALL_OF"
    approver_groups = [
      {
        logical_operator = "AND"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["example@company.io"]
          }
        ]
      }
    ]
  }

  escalation_policy {
    interval_in_min = 30
    approver_groups = [
      {
        logical_operator = "OR"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["security@company.io"]
          }
        ]
      },
      {
        logical_operator = "OR"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["platform-admin@company.io"]
          },
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["incident-manager@company.io"]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    require_approver_reason       = false
    requester_cannot_approve_self = false
    require_mfa                   = true
    max_extensions                = 0
    extension_duration_in_min     = 0
    labels                        = ["created_from_terraform"]
  }
}

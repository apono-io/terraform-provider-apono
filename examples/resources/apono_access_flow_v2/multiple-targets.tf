resource "apono_access_flow_v2" "multiple_resources_flow" {
  name                  = "Access Azure Subscription Integration"
  active                = true
  grant_duration_in_min = 90
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "AND"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND-team"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "Azure Subscription Integration"
        resource_type    = "azure-subscription-resource-group"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "NAME"
          values     = ["Resource 1", "Resource 2", "Resource 3"]
        }]
        permissions = ["Key Vault Administrator"]
      }
    },
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

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = true
    labels                        = ["multiple_resources", "azure_integration"]
  }
}

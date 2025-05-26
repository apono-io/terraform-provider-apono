resource "apono_access_flow_v2" "owner_approver_flow" {
  name                  = "AWS Prod Env - Integration Owner Approval"
  active                = true
  grant_duration_in_min = 120
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "AND"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["Infra Admins"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = data.apono_resource_integrations.aws_staging_integrations.integrations[0].name
        resource_type    = "aws-account-s3-bucket"
        permissions      = ["READ_WRITE"]
      }
    }
  ]

  approver_policy = {
    approval_mode = "ANY_OF"
    approver_groups = [
      {
        logical_operator = "AND"
        approvers = [
          {
            type = "Owner"
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = true
    labels                        = ["created_from_terraform"]
  }
}

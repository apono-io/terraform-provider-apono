resource "apono_access_flow_v2" "agent_prod_db" {
  name                    = "Agent access to production DBs"
  active                  = true
  grant_duration_in_min   = 60
  trigger                 = "SELF_SERVE"
  requestor_identity_type = "AGENT"

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
        integration_name = "aws-account-prod"
        resource_type    = "aws-account-s3-bucket"
        permissions      = ["ReadOnlyAccess"]
      }
    }
  ]

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = false
    max_extensions                = 0
    extension_duration_in_min     = 0
    labels                        = ["agentic"]
  }
}

resource "apono_access_flow_v2" "dynamic_prod_db" {
  name                  = "Dynamic access to production DBs"
  active                = true
  grant_duration_in_min = 60
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
      query_target = {
        query = "integration in (\"aws-rds\") and resource_tag[\"environment\"] = \"production\""
      }
    }
  ]

  approver_policy = {
    approval_mode = "ANY_OF"
    approver_groups = [
      {
        logical_operator = "OR"
        approvers = [
          {
            type           = "group"
            match_operator = "is"
            values         = [data.apono_groups.DevOps_team.groups[0].id]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = false
    max_extensions                = 0
    extension_duration_in_min     = 0
    labels                        = ["dynamic_targets"]
  }
}

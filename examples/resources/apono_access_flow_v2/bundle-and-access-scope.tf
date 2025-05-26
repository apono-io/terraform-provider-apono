resource "apono_access_flow_v2" "bundle_access_scope_flow" {
  name                  = "Access to production DBs"
  active                = true
  grant_duration_in_min = 30
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
      bundle = {
        name = data.apono_bundles.critical_prod_db_bundle.bundles[0].name
      }
    },
    {
      access_scope = {
        name = data.apono_access_scopes.production_db.access_scopes[0].name
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
            source_integration_name = "Google Oauth"
            type                    = "group"
            match_operator          = "is"
            values                  = [data.apono_groups.InfoSec_team.groups[0].id]
          },
          {
            type           = "group"
            match_operator = "is"
            values         = [data.apono_groups.DevOps_team.groups[0].id]
          },
          {
            source_integration_name = "Google Oauth"
            type                    = "group"
            match_operator          = "is"
            values                  = [data.apono_groups.dev_teams.groups[0].id]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = false
    labels                        = ["bundle_access", "scope_reference"]
  }
}

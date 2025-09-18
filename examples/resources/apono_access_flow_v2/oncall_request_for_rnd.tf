resource "apono_access_flow_v2" "oncall_request_for_rnd" {
  name                  = "On-Call request for R&D group"
  active                = true
  grant_duration_in_min = 30
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type                    = "opsgenie_schedule"
        source_integration_name = "Opsgenie"
        match_operator          = "contains"
        values                  = ["night_shift"]
      }
    ]
  }

  request_for = {
    request_scopes = ["self", "others"]
    grantees = {
      logical_operator = "OR"
      conditions = [
        {
          type                    = "group"
          match_operator          = "contains"
          source_integration_name = "Google Oauth"
          values                  = ["RND_TEAM"]
        }
      ]
    }
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

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = false
    labels                        = ["bundle_access", "scope_reference"]
  }
}

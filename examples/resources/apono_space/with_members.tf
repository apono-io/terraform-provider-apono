resource "apono_space_scope" "production_aws" {
  name  = "Production AWS"
  query = "integration in (\"aws-account\") and resource_tag[\"environment\"] = \"production\""
}

resource "apono_space_scope" "production_gcp" {
  name  = "Production GCP"
  query = "integration in (\"gcp-project\") and resource_tag[\"environment\"] = \"production\""
}

resource "apono_space" "production" {
  name = "Production"
  space_scope_references = [
    apono_space_scope.production_aws.name,
    apono_space_scope.production_gcp.name,
  ]

  members = [
    {
      identity_reference = "platform-team-lead@example.com"
      identity_type      = "user"
      space_roles        = ["SpaceOwner"]
    },
    {
      identity_reference = "platform-engineers"
      identity_type      = "group"
      space_roles        = ["SpaceManager"]
    },
  ]
}

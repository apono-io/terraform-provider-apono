resource "apono_space_scope" "engineering_aws" {
  name  = "Engineering AWS"
  query = "integration in (\"aws-account\") and resource_tag[\"team\"] = \"engineering\""
}

resource "apono_space" "engineering" {
  name = "Engineering"
  space_scope_references = [
    apono_space_scope.engineering_aws.name,
  ]
}

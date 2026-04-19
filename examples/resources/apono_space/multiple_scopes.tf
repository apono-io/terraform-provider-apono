resource "apono_space" "staging" {
  name = "Staging"
  space_scope_references = [
    apono_space_scope.staging_aws.name,
    apono_space_scope.staging_gcp.name,
  ]
}

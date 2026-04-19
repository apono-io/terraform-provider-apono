resource "apono_space_scope" "production_aws" {
  name  = "Production AWS"
  query = "integration in (\"aws-account\") and resource_tag[\"environment\"] = \"production\""
}

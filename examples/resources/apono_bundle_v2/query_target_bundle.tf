resource "apono_bundle_v2" "prod_db_bundle" {
  name = "Production DB Bundle"

  access_targets = [
    {
      query_target = {
        query = "integration in (\"aws-rds\") and resource_tag[\"environment\"] = \"production\""
      }
    }
  ]
}

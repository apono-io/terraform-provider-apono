resource "apono_bundle_v2" "platform_bundle" {
  name = "Platform Access"

  access_targets = [
    {
      integration = {
        integration_name = "github"
        resource_type    = "repository"
        permissions      = ["Read"]
      }
    },
    {
      query_target = {
        query = "integration in (\"aws-rds\") and resource_tag[\"team\"] = \"platform\""
      }
    }
  ]
}

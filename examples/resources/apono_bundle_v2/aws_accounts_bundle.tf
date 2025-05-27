resource "apono_bundle_v2" "aws_accounts_bundle" {
  name = "AWS Account Bundle"

  access_targets = [
    {
      integration = {
        integration_name = "AWS ORG"
        resource_type    = "aws-organization-account"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "TAG"
          key        = "aws_account_id"
          values     = ["46000000255"]
        }]
        permissions = ["TestRunner"]
      }
    }
  ]
}

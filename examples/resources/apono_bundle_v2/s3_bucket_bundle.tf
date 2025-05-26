resource "apono_bundle_v2" "s3_bucket_bundle" {
  name = "S3 Logs and Sensitive Buckets"

  access_targets = [
    {
      integration = {
        integration_name = "Amazon Account Integration"
        resource_type    = "aws-account-s3-bucket"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "NAME"
          values = [
            "s3-logs-to-grafana",
            "s3-sensitive-data"
          ]
        }]
        permissions = ["ADMIN"]
      },
      integration = {
        integration_name = "Amazon Organization Integration"
        resource_type    = "aws-organization-s3-bucket"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "NAME"
          values = [
            "s3-logs-to-Kibana",
            "s3-sensitive-data"
          ]
        }]
        permissions = ["ADMIN"]
      }
    }
  ]
}

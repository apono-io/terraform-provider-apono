resource "apono_bundle_v2" "big_query_viewer_bundle" {
  name = "BigQuery Viewer"

  access_targets = [
    {
      integration = {
        integration_name = "GCP Integration"
        resource_type    = "gcp-organization-bigquery-dataset"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "APONO_ID"
          values = [
            "d92a7f33b4d832c55dca9e2afc3d9985de4a227abb1a88e0f9c3dc08b12b57e6",
            "8f00e1a39a7f4cf2d22d1f30aebf6dcf7a5dbf83d912b4ea5c6c6fcb22cf1d09",
            "719d1626f380b0f9c8aab3ee92c5371706d816e32ebee0cb6b84e7cd759d2a0e"
          ]
        }]
        permissions = ["BigQuery Data Viewer"]
      }
    }
  ]
}

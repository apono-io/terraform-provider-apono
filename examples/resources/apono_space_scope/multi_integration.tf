resource "apono_space_scope" "cloud_production" {
  name = "Cloud Production"
  query = <<EOT
  (integration in ("aws-account") or integration in ("gcp-project"))
    and resource_tag["environment"] = "production"
  EOT
}

access_targets = [
  {
    query_target = {
      query = <<-EOT
        integration in ("aws-rds") and resource_tag["environment"] = "production"
      EOT
    }
  }
]

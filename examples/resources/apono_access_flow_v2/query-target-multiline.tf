access_targets = [
  {
    query_target = {
      query = <<-EOT
        (integration in ("aws-account") or integration in ("gcp-project"))
          and resource_tag["environment"] = "production"
          and resource_tag["compliance"] = "pci"
      EOT
    }
  }
]

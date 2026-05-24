access_targets = [
  {
    integration = {
      integration_name = "aws-account-prod"
      resource_type    = "aws-account-s3-bucket"
      permissions      = ["ReadOnlyAccess"]
    }
  },
  {
    bundle = {
      name = data.apono_bundles.critical_prod_db_bundle.bundles[0].name
    }
  },
  {
    access_scope = {
      name = data.apono_access_scopes.production_db.access_scopes[0].name
    }
  },
  {
    query_target = {
      query = "integration in (\"aws-rds\") and resource_tag[\"team\"] = \"platform\""
    }
  }
]

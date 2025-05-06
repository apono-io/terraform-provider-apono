resource "apono_access_flow_v2" "aws_auto_grant_flow" {
  name    = "AWS Automatic Access Flow"
  trigger = "AUTOMATIC"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "aws-account-integration"
        resource_type    = "aws-account-iam-policy"
        permissions      = ["Attach"]
      }
    }
  ]

  settings = {
    justification_required = false
  }
}

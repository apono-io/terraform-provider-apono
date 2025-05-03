resource "apono_access_flow_v2" "aws_basic_flow" {
  name    = "AWS Attach Access Flow"
  trigger = "SELF_SERVE"

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

  settings = {} // Default settings applied
}

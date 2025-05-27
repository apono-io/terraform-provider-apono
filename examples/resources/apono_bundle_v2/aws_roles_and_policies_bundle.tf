resource "apono_bundle_v2" "aws_roles_and_policies_bundle" {
  name = "Access To AWS Roles and Policies"

  access_targets = [
    {
      access_scope = {
        name = "AWS Production IAM Policies"
      }
    },
    {
      integration = {
        integration_name = "Amazon Account Integration"
        resource_type    = "aws-account-iam-role"
        permissions      = ["Attach"]
      }
    }
  ]
}

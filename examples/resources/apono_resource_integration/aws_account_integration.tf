resource "apono_resource_integration" "aws_account_integration" {
  name         = "AWS Account Integration"
  type         = "aws-account"
  connector_id = "AwsIntegrationConnector-Semyon-NDPp"
  connected_resource_types = [
    "aws-account-iam-group",
    "aws-account-s3-bucket"
  ]
  integration_config = {
    region                              = "us-east-1"
    profile                             = "apono"
    credentials_rotation_period_in_days = "30"
    credentials_cleanup_period_in_days  = "90"
  }
}

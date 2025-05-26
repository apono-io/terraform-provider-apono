data "apono_resource_integrations" "aws_staging_integrations" {
  name         = "*staging-*"
  type         = "aws-*"
  connector_id = "AwsIntegrationConnector"
}

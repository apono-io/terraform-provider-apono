resource "apono_access_scope" "critical_ec2_instances" {
  name  = "critical-ec2"
  query = <<EOT
  resource_type = "aws-account-ec2-instance"
    and resource_tag["env"] = "production" 
    and resource_risk_level = "1"
  EOT
}

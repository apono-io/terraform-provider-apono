data "apono_user_information_integrations" "all" {}

output "all_user_information_integrations" {
  value = data.apono_user_information_integrations.all.integrations
}

data "apono_user_information_integrations" "jumpcloud" {
  name = "Jumpcloud IDP"
}

output "jumpcloud_user_information_integration" {
  value = data.apono_user_information_integrations.jumpcloud.integrations
}

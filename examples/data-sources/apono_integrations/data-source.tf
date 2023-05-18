data "apono_integrations" "k8s_integrations" {
  type = "k8s-roles"
}

data "apono_integrations" "prod_integrations" {
  connector_id = "prod-us-east-1"
}

data "apono_integrations" "prod_mysql_integrations" {
  type         = "mysql"
  connector_id = "production-us-east-1"
}

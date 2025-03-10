terraform {
  backend "local" {
  }
  required_providers {
    apono = {
      source = "terraform-registry.apono.com/apono-io/apono"
    }
  }
}

provider "apono" {}

variable "second_run" {}
variable "connector_id" {
  default = "terraofrm-tests-account-connector"
}
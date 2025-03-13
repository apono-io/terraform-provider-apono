resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "apono_integration" "benchmark_integration" {
  name         = "terraform-test-${random_string.random.result}"
  type         = "benchmark"
  connector_id = var.connector_id
  metadata = {
    resource_count   = var.second_run ? 5 : 10
    permission_count = 6
  }
}

output "integration_type" {
  value = apono_integration.benchmark_integration.type
}

output "integration_name" {
  value = replace(apono_integration.benchmark_integration.name, random_string.random.result, "")
}

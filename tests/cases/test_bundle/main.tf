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
    resource_count   = 5
    permission_count = 6
  }
}
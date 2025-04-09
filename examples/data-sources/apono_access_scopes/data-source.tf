# Retrieve all access scopes
data "apono_access_scopes" "all" {
}

# Retrieve a specific access scope by exact name
data "apono_access_scopes" "production_db" {
  name = "Production Database Access"
}

# Retrieve access scopes matching a pattern
data "apono_access_scopes" "production_scopes" {
  name = "Production*"
}

# Output example: Access the first access scope in the list
output "first_access_scope" {
  value = data.apono_access_scopes.all.access_scopes[0].id
}

# Output example: Count of production scopes
output "production_scopes_count" {
  value = length(data.apono_access_scopes.production_scopes.access_scopes)
}

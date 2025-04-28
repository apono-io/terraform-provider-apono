# Get a specific group by exact name
data "apono_groups" "exact" {
  name = "engineering"
}

# Get all groups with names starting with "dev-"
data "apono_groups" "dev_teams" {
  name = "dev-*"
}

# Get all groups from a specific source integration
data "apono_groups" "from_source" {
  source_integration = "si-12345678"
}

# Access the groups data
output "engineering_group_id" {
  value = data.apono_groups.exact.groups[0].id
}

output "dev_team_names" {
  value = [for group in data.apono_groups.dev_teams.groups : group.name]
}

output "source_integration_groups" {
  value = [for group in data.apono_groups.from_source.groups : {
    id   = group.id
    name = group.name
  }]
}

resource "apono_group" "engineering_team" {
  name = "Engineering Team"
  members = [
    "alice@example.com",
    "bob@example.com",
    "charlie@example.com"
  ]
}

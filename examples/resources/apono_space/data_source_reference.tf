data "apono_space_scopes" "platform" {
  name = "Platform*"
}

resource "apono_space" "platform" {
  name = "Platform"
  space_scope_references = [
    data.apono_space_scopes.platform.space_scopes[0].name,
  ]

  members = [
    {
      identity_reference = "admin@example.com"
      identity_type      = "user"
      space_roles        = ["SpaceOwner"]
    },
  ]
}

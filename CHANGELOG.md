# Changelog

## [v1.7.2] - 2025-05-28

### Breaking Changes

Updated schema types in beta data sources and resources: Several fields previously returned as `set` now return a `list` for better performance.

#### Affected Data Sources and Resources and Fields:

- `apono_access_scopes.access_scopes`
- `apono_groups.groups`
- `apono_user_information_integration.integrations`
- `apono_access_flow_v2.approver_policy.approver_groups.approvers`
- `apono_access_flow_v2.approver_policy.approver_groups.approvers.values`
- `apono_access_flow_v2.requestors.conditions`
- `apono_access_flow_v2.requestors.conditions.values`
- `apono_access_flow_v2.access_targets`
- `apono_access_flow_v2.access_targets.integration.resources_scopes`
- `apono_access_flow_v2.access_targets.integration.resources_scopes.values`

If your configuration depends on these fields, ensure any use of `toset()` or assumptions about uniqueness are revised accordingly.

### Deprecations

In the `apono_resource_integration` resource, the following fields in the `owner` block are now deprecated and will be removed in v2.0.0:

- `type` → use `attribute_type`
- `values` → use `attribute_values`

**Notice:** These old fields are still fully functional in this release but will be removed in the next major version. Please update your configuration to use the new field names to ensure compatibility going forward.


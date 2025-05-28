# Changelog

## [v1.8.2] - 2025-05-28

### Breaking Changes

Updated schema types in data sources and resources: Several fields previously returned as `set` now return a `list` for better performance.

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

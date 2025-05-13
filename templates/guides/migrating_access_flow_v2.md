# Migration doc for Access Flow

## Migrating Access Flow Resource

This guide will help you migrate your existing Terraform configuration for the Access Flow from the `apono_access_flow` resource schema to the updated `apono_access_flow_v2` schema.

---

### Why Migrate?

The new `apono_access_flow_v2` resource provides a cleaner, more consistent schema that simplifies configuration and improves maintainability. It introduces new capabilities that align Terraform with Apono’s latest UI capabilities, with support for access scopes, improved requestor and approver logic, and expanded settings.

Migrating ensures compatibility with future improvements and unlocks stronger access control and governance.

See the table below for a full breakdown of changes and new features.

---

## Key Changes Overview

| Concept                    | Old Schema                        | New Schema                      | Changes’ Description                                                      |
|----------------------------|-----------------------------------|----------------------------------|---------------------------------------------------------------------------|
| Time-Limited Access        | `revoke_after_in_sec`             | `grant_duration_in_min`          | Renamed; units changed from seconds to minutes                            |
| Trigger                    | `trigger.type` (user_request, automatic) | `trigger` (SELF_SERVE, AUTOMATIC) | Converted from object to flat string                                      |
| Timeframe                  | `trigger.timeframe`               | `timeframe`                      | Moved to a top-level block                                                |
| Access Targets             | `bundle_targets`, `integration_targets` | `access_targets`                 | Unified targeting model                                                   |
| Target Resource Filtering  | `resource_include_filters` / `resource_exclude_filters` | `resources_scopes`               | Unified under resources_scopes with scope_mode                            |
| Grantee Conditions         | `grantees_conditions_group`        | `requestors`                     | Unified and renamed                                                       |
| Approval Policy Group Logic| `approver_groups_relationship`     | `approval_mode`                  | Renamed                                                                   |
| Approver Group Structure   | `attribute_conditions` → nested under approver_groups | `approver_groups.approvers`      | Structural change and new nesting                                         |
| Approver and Requestor Definition | `attribute_type`           | `type`                           | Renamed for both requestors.conditions and approvers.approvers            |
| Value Matching Field       | `attribute_names`                  | `values`                         | Renamed and expanded to support both source ID and Apono ID               |
| Operator Field             | `operator`                         | `match_operator`                 | Renamed; same allowed values (is, is_not, etc.)                           |
| Integration Source         | `integration_id`                   | `source_integration_name`        | Renamed                                                                   |
| New Field                  |                                   | `require_approver_reason`        | Require a reason from the approver                                        |
| New Field                  |                                   | `access_scope.name`              | Allows access scope as a target                                           |

Refer to the new resource documentation [here](../resources/access_flow_v2.md) for full schema details.

---

## Migration Steps

### Prerequisites

- You are using Terraform v1.5.0 or later
- The resource uses the `apono_access_flow` schema

---

### Step 1: Identify Usage

Locate all existing uses of the `apono_access_flow` resource.

```shell
grep -r 'apono_access_flow' .
```
Or search in your editor.

---

### Step 2: Review the State

Run:
```shell
terraform state list | grep apono_access_flow
```
Then, inspect the existing resource:
```shell
terraform state show <resource_address>
```
Use this to verify the existing resource address and retrieve the access flow ID to be used in the import block.

---

### Step 3: Generate the Configuration

Use Terraform to automatically generate a configuration for the new `apono_access_flow_v2` resource:

1. Create a new file called `import.tf`
2. Inside the file, create an import block:

    ```hcl
    import {
      to = apono_access_flow_v2.example
      id = "<access_flow_id>"
    }
    ```

3. Run:

    ```shell
    terraform plan -generate-config-out=generated.tf
    ```

This generates a valid `apono_access_flow_v2` block based on the current state of the resource in Apono.

---

### Step 4: Remove the Old Resource from State

To disconnect the old resource without destroying it, run:

```shell
terraform state rm <resource_address>
```

⚠️ This removes the old `apono_access_flow` from Terraform state, not from the Apono platform.

---

### Step 5: Replace the Old Terraform Block

Delete the old `apono_access_flow` resource block from your code.

Replace it with the generated `apono_access_flow_v2` block from `generated.tf`. Copy it into your main `.tf` files and clean up as needed.

---

### Step 6: Import the Resource

Use the CLI to finalize the import:

```shell
terraform import apono_access_flow_v2.example <access_flow_id>
```

Ensure this matches the target in your import block and the name of your resource block.

---

### Step 7: Plan and Apply

Run:

```shell
terraform plan
```

Confirm no changes are about to be made. Then, run:

```shell
terraform apply
```

# Migration doc for Resource Integration

## Migrating Integration Resource

This guide will help you migrate your existing Terraform configuration for the Resource Integration from the `apono_integration` resource schema to the updated `apono_resource_integration` schema.

---

### Why Migrate?

The new `apono_resource_integration` resource provides a cleaner, more consistent schema that simplifies configuration and improves maintainability.

Migrating ensures compatibility with future improvements and unlocks stronger access control and governance.

See the table below for a full breakdown of changes.

---

## Key Changes Overview

| Concept                  | Old Schema                                                      | New Schema                   | Changes’ Description                          |
|--------------------------|-----------------------------------------------------------------|------------------------------|-----------------------------------------------|
| Secrets Handling         | Multiple separate blocks: `aws_secret`, `gcp_secret`, `hashicorp_vault_secret`, `kubernetes_secret` | `secret_store_config`        | Unified under one block; one secret store per integration |
| Resource Ownership Mapping | `resource_owner_mappings`                                     | `owners_mapping`             | Renamed                                       |
| Fallback Owners          | `integration_owners`                                            | `owner`                      | Renamed                                       |
| Integration Configuration| `metadata`                                                      | `integration_config` (Map of String) | Renamed                               |

Refer to the new resource documentation [here](../resources/resource_integration.md) for full schema details.

---

## Migration Steps

### Prerequisites

- You are using Terraform v1.5.0 or later
- The resource uses the `apono_integration` schema

---

### Step 1: Identify Usage

Locate all existing uses of the old `apono_integration` resource.

```shell
grep -r 'apono_integration' .
```
Or search in your editor.

---

### Step 2: Review the State

Run:
```shell
terraform state list | grep apono_integration
```
Then, inspect the existing resource:
```shell
terraform state show <resource_address>
```
Use this to verify the existing resource address and retrieve the integration ID to be used in the import block.

---

### Step 3: Generate the Configuration

Use Terraform to automatically generate a configuration for the new `apono_resource_integration` resource:

1. Create a new file called `import.tf`
2. Inside the file, create an import block:

    ```hcl
    import {
      to = apono_resource_integration.example
      id = "<resource_integration_id>"
    }
    ```

3. Run:

    ```shell
    terraform plan -generate-config-out=generated.tf
    ```

This generates a valid `apono_resource_integration` block based on the current state of the resource in Apono.

---

### Step 4: Remove the Old Resource from State

To disconnect the old resource without destroying it, run:

```shell
terraform state rm <resource_address>
```

⚠️ This removes the old `apono_integration` from the Terraform state, not from the Apono platform.

---

### Step 5: Replace the Old Terraform Block

Delete the old `apono_integration` resource block from your code.

Replace it with the generated `apono_resource_integration` block from `generated.tf`. Copy it into your main `.tf` files and clean up as needed.

---

### Step 6: Import the Resource

Use the CLI to finalize the import:

```shell
terraform import apono_resource_integration.example <resource_integration_id>
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

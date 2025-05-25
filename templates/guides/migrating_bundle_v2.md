## Migrating Bundle Resource

This guide provides instructions for migrating your existing `apono_access_bundle` resources to the updated `apono_bundle_v2` schema.

---

### Why Migrate?

The new `apono_bundle_v2` resource introduces support for access scopes as access targets and clearer include/exclude behavior via resources scopes.  
Migrating ensures compatibility with future improvements.

See the table below for a full breakdown of changes.

---

## Key Changes Overview

| Concept                  | Old Schema                        | New Schema         | Changes’ Description                                 |
|--------------------------|-----------------------------------|--------------------|------------------------------------------------------|
| Access Targets Block     | `integration_targets`             | `access_targets`   | Unified top-level block for bundle targets           |
| Target Resource Filtering| `resource_include_filters`, `resource_exclude_filters` | `resources_scopes` | Unified under resources_scopes with scope_mode       |
| New Target Type          |                                   | `access_scope`     | Allows applying pre-defined access scopes            |

Refer to the new resource documentation [here](../resources/bundle_v2.md) for full schema details.

---

## Migration Steps

### Prerequisites

- You are using Terraform v1.5.0 or later
- The resource uses the `apono_access_bundle` schema

---

### Step 1: Identify Usage

Locate all existing uses of the old `apono_access_bundle` resource.

🔎 Use:
```shell
grep -r apono_access_bundle .
```
Or search in your editor.

---

### Step 2: Review the State

Run:
```shell
terraform state list | grep apono_access_bundle
```
Then, inspect the existing resource:
```shell
terraform state show <resource_address>
```
Use this to verify the existing resource address and retrieve the bundle ID to be used in the import block.

---

### Step 3: Generate the Configuration

Use Terraform to automatically generate a configuration for the new `apono_bundle_v2` resource:

1. Create a new file called `import.tf`
2. Inside the file, create an import block:

    ```hcl
    import {
      to = apono_bundle_v2.example
      id = "<bundle_id>"
    }
    ```

3. Run:

    ```shell
    terraform plan -generate-config-out=generated.tf
    ```

This generates a valid `apono_bundle_v2` block based on the current state of the resource in Apono.

---

### Step 4: Remove the Old Resource from State

To disconnect the old resource without destroying it, run:

```shell
terraform state rm <resource_address>
```

⚠️ This removes the old `apono_access_bundle` from the Terraform state, not from the Apono platform.

---

### Step 5: Replace the Old Terraform Block

Delete the old `apono_access_bundle` resource block from your code.

Replace it with the generated `apono_bundle_v2` block from `generated.tf`.  
Copy it into your main `.tf` files and clean up as needed.

---

### Step 6: Import the Resource

Use the CLI to finalize the import:

```shell
terraform import apono_bundle_v2.example <bundle_id>
```

Ensure this matches the target in your import block and the name of your resource block.

---

### Step 7: Plan and Apply

Run:

```shell
terraform plan
```

Confirm no changes are about to be made.  
Then, run:

```shell
terraform apply
```

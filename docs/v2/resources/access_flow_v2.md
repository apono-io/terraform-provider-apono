---
page_title: "apono_access_flow_v2 (beta) Resource - terraform-provider-apono"
subcategory: ""
description: |-
    Manages an Apono Access Flow that defines how users or groups can request or automatically be granted access to integrations, bundles, or access scopes under specific conditions and policies.
---

# Resource: apono_access_flow_v2 (beta)

Manages an Apono Access Flow that defines how users or groups can request or automatically be granted access to integrations, bundles, or access scopes under specific conditions and policies.

## Example Usage

### Basic Example - Self-Serve Access Flow 

```terraform
resource "apono_access_flow_v2" "aws_basic_flow" {
  name    = "AWS Attach Access Flow"
  trigger = "SELF_SERVE"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "aws-account-integration"
        resource_type    = "aws-account-iam-policy"
        permissions      = ["Attach"]
      }
    }
  ]

  settings = {} // Default settings applied
}
```

### Automatic Approval Access Flow

```terraform
resource "apono_access_flow_v2" "aws_auto_grant_flow" {
  name    = "AWS Automatic Access Flow"
  trigger = "AUTOMATIC"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "aws-account-integration"
        resource_type    = "aws-account-iam-policy"
        permissions      = ["Attach"]
      }
    }
  ]

  settings = {
    justification_required = false
  }
}
```

### Access to Sensitive Production for Users via Google Oauth IDP  

```terraform
resource "apono_access_flow_v2" "sensitive_production_aws" {
  name                  = "Sensitive Access to Production AWS"
  active                = true
  grant_duration_in_min = 60
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "OR"
    conditions = [
      {
        type                    = "user"
        match_operator          = "is"
        source_integration_name = data.apono_user_information_integrations.google_oauth_idp.integrations[0].name
        values                  = ["example@company.io"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "Azure Subscription Integration"
        resource_type    = "azure-subscription-sql-server"
        permissions      = ["Contributor"]
      }
    }
  ]

  approver_policy = {
    approval_mode = "ALL_OF"
    approver_groups = [
      {
        logical_operator = "AND"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["example@company.io"]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    require_approver_reason       = false
    requester_cannot_approve_self = false
    require_mfa                   = true
    labels                        = ["created_from_terraform"]
  }
}
```

### Access to Multiple Integration Targets and Specific Resource Names

```terraform
resource "apono_access_flow_v2" "multiple_resources_flow" {
  name                  = "Access Azure Subscription Integration"
  active                = true
  grant_duration_in_min = 90
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "AND"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND-team"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = "Azure Subscription Integration"
        resource_type    = "azure-subscription-resource-group"
        resources_scopes = [{
          scope_mode = "include_resources"
          type       = "NAME"
          values     = ["Resource 1", "Resource 2", "Resource 3"]
        }]
        permissions = ["Key Vault Administrator"]
      }
    },
    {
      integration = {
        integration_name = "Azure Subscription Integration"
        resource_type    = "azure-subscription-sql-server"
        permissions      = ["Contributor"]
      }
    }
  ]

  approver_policy = {
    approval_mode = "ALL_OF"
    approver_groups = [
      {
        logical_operator = "AND"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "user"
            match_operator          = "is"
            values                  = ["example@company.io"]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = true
    labels                        = ["multiple_resources", "azure_integration"]
  }
}
```

### Bundle and Access Scope as Access Targets 

```terraform
resource "apono_access_flow_v2" "bundle_access_scope_flow" {
  name                  = "Access to production DBs"
  active                = true
  grant_duration_in_min = 30
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "AND"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["RND-team"]
      }
    ]
  }

  access_targets = [
    {
      bundle = {
        name = data.apono_bundles.critical_prod_db_bundle.bundles[0].name
      }
    },
    {
      access_scope = {
        name = data.apono_access_scopes.production_db.access_scopes[0].name
      }
    }
  ]

  approver_policy = {
    approval_mode = "ANY_OF"
    approver_groups = [
      {
        logical_operator = "OR"
        approvers = [
          {
            source_integration_name = "Google Oauth"
            type                    = "group"
            match_operator          = "is"
            values                  = [data.apono_groups.InfoSec_team.groups[0].id]
          },
          {
            type           = "group"
            match_operator = "is"
            values         = [data.apono_groups.DevOps_team.groups[0].id]
          },
          {
            source_integration_name = "Google Oauth"
            type                    = "group"
            match_operator          = "is"
            values                  = [data.apono_groups.dev_teams.groups[0].id]
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = false
    labels                        = ["bundle_access", "scope_reference"]
  }
}
```

### Integration Owner as Approver

```terraform
resource "apono_access_flow_v2" "owner_approver_flow" {
  name                  = "AWS Prod Env - Integration Owner Approval"
  active                = true
  grant_duration_in_min = 120
  trigger               = "SELF_SERVE"

  requestors = {
    logical_operator = "AND"
    conditions = [
      {
        type           = "group"
        match_operator = "contains"
        values         = ["Infra Admins"]
      }
    ]
  }

  access_targets = [
    {
      integration = {
        integration_name = data.apono_resource_integrations.aws_staging_integrations.integrations[0].name
        resource_type    = "aws-account-s3-bucket"
        permissions      = ["READ_WRITE"]
      }
    }
  ]

  approver_policy = {
    approval_mode = "ANY_OF"
    approver_groups = [
      {
        logical_operator = "AND"
        approvers = [
          {
            type = "Owner"
          }
        ]
      }
    ]
  }

  settings = {
    justification_required        = true
    requester_cannot_approve_self = true
    require_mfa                   = true
    labels                        = ["created_from_terraform"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_targets` (Attributes List) Define the targets accessible when requesting access via this access flow. (see [below for nested schema](#nestedatt--access_targets))
- `name` (String) Human-readable name for the access flow, must be unique.
- `requestors` (Attributes) Defines who can request access. (see [below for nested schema](#nestedatt--requestors))
- `settings` (Attributes) Settings for the access flow. (see [below for nested schema](#nestedatt--settings))
- `trigger` (String) The trigger type for the access flow. Possible values: SELF_SERVE, AUTOMATIC.

### Optional

- `active` (Boolean) Whether the access flow is active. Defaults to true.
- `approver_policy` (Attributes) Approval policy for the access request. (see [below for nested schema](#nestedatt--approver_policy))
- `grant_duration_in_min` (Number) How long access is granted, in minutes. If not specified, the grant duration defaults to indefinite.
- `timeframe` (Attributes) Restrict when access can be granted. (see [below for nested schema](#nestedatt--timeframe))

### Read-Only

- `id` (String) The unique identifier of the access flow.

<a id="nestedatt--access_targets"></a>
### Nested Schema for `access_targets`

Optional:

- `access_scope` (Attributes) Access scope target. (see [below for nested schema](#nestedatt--access_targets--access_scope))
- `bundle` (Attributes) Bundle target. (see [below for nested schema](#nestedatt--access_targets--bundle))
- `integration` (Attributes) Defines an integration and resources to which access will be granted. (see [below for nested schema](#nestedatt--access_targets--integration))

<a id="nestedatt--access_targets--access_scope"></a>
### Nested Schema for `access_targets.access_scope`

Required:

- `name` (String) Name of the access scope.


<a id="nestedatt--access_targets--bundle"></a>
### Nested Schema for `access_targets.bundle`

Required:

- `name` (String) Name of the bundle.


<a id="nestedatt--access_targets--integration"></a>
### Nested Schema for `access_targets.integration`

Required:

- `integration_name` (String) The name of the integration
- `permissions` (Set of String) List of permissions (e.g., "Attach", "ReadOnlyAccess").
- `resource_type` (String) The type of resource within the integration for which access is being granted (e.g., aws-account-s3-bucket).

Optional:

- `resources_scopes` (Attributes List) A list of filters defining which resources are included or excluded. If null, the scope will apply to any resource in the integration target (see [below for nested schema](#nestedatt--access_targets--integration--resources_scopes))

<a id="nestedatt--access_targets--integration--resources_scopes"></a>
### Nested Schema for `access_targets.integration.resources_scopes`

Required:

- `scope_mode` (String) Possible values: `include_resources` or `exclude_resources`. `include_resources`: Grants access to the specific resources listed under the `values` field. `exclude_resources`: Grants access to all resources within the integration except those specified in the `values` field.
- `type` (String) NAME - specify resources by their name, APONO_ID - specify resources by their ID, or TAG - specify resources by tag.
- `values` (List of String) Resource values to match (IDs, names, or tag values).

Optional:

- `key` (String) Tag key. Only required if type = TAG




<a id="nestedatt--requestors"></a>
### Nested Schema for `requestors`

Required:

- `conditions` (Attributes List) List of conditions. Cannot be empty. (see [below for nested schema](#nestedatt--requestors--conditions))
- `logical_operator` (String) Specifies the logical operator to be used between the requestors in the list. Possible values: "AND" or "OR".

<a id="nestedatt--requestors--conditions"></a>
### Nested Schema for `requestors.conditions`

Required:

- `type` (String) Identity type (e.g., user, group, etc.)

Optional:

- `match_operator` (String) Comparison operator. Possible values: is, is_not, contains, does_not_contain, starts_with. Defaults to is.
Note: When using is or is_not with any type, you can specify either the source ID or Apono ID to define the requestors.
For the user attribute specifically, you may also use the user’s email.
- `source_integration_name` (String) The integration the user/group is from.
- `values` (List of String) List of values according to the attribute type and match_operator (e.g., user emails, group IDs, etc.).



<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `justification_required` (Boolean) Require justification from requestor. Defaults to true. Must be set to false for automatic access flows.
- `labels` (Set of String) Custom labels for organizational use.
- `requester_cannot_approve_self` (Boolean) Requester cannot approve their own requests. Defaults to false.
- `require_approver_reason` (Boolean) Require reason from approver. Defaults to false.
- `require_mfa` (Boolean) Require MFA at approval time. Defaults to false.


<a id="nestedatt--approver_policy"></a>
### Nested Schema for `approver_policy`

Required:

- `approval_mode` (String) Possible values: ANY_OF or ALL_OF. Specifies the logical condition for approvals: ANY_OF: The request is granted if at least one approver from the list approves. ALL_OF: The request is granted only if all approvers in the list approve.
- `approver_groups` (Attributes Set) List of approver groups. Cannot be empty. (see [below for nested schema](#nestedatt--approver_policy--approver_groups))

<a id="nestedatt--approver_policy--approver_groups"></a>
### Nested Schema for `approver_policy.approver_groups`

Required:

- `approvers` (Attributes List) List of approvers. (see [below for nested schema](#nestedatt--approver_policy--approver_groups--approvers))
- `logical_operator` (String) Possible values: AND or OR

<a id="nestedatt--approver_policy--approver_groups--approvers"></a>
### Nested Schema for `approver_policy.approver_groups.approvers`

Required:

- `type` (String) Approver identity type - user, group, Owner, manager, Context Integration, or any other custom value.
Note: The Owner value must be capitalized (with an uppercase “O”).

Optional:

- `match_operator` (String) Comparison operator. Possible values: is, is_not, contains, does_not_contain, starts_with. Defaults to is.
Note: When using is or is_not with any type, you can specify either the source ID or Apono ID to define the requestors.
For the user attribute specifically, you may also use the user’s email.
- `source_integration_name` (String) Applies when the identity type stems from a Context or IDP integration.
- `values` (List of String) Approver values according to the attribute type and match_operator (e.g., user email, group IDs, etc).




<a id="nestedatt--timeframe"></a>
### Nested Schema for `timeframe`

Required:

- `days_of_week` (Set of String) Days when access is allowed. (e.g., ['MONDAY', 'TUESDAY']).
- `end_time` (String) End time (e.g., 17:00).
- `start_time` (String) Start time (e.g., 08:00).
- `time_zone` (String) Timezone name (e.g., Asia/Jerusalem).

## Import

In Terraform v1.5.0 and later, use an import block to import apono_access_flow_v2 using the Access Flow identifier. For example:

```terraform
import {
  to = apono_access_flow_v2.example
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

Or using the CLI:

```shell
terraform import apono_access_flow_v2.example 123e4567-e89b-12d3-a456-426614174000
```

<h1>
    <a href="https://apono.io">
      <img src="./assets/logo.svg" style="float: right" height="46px" alt="Apono logo"/>
    </a>
    <span>&nbsp;Apono Terraform Provider</span>
</h1>

[![Actions Status](https://github.com/apono-io/terraform-provider-apono/workflows/Release/badge.svg)](https://github.com/apono-io/terraform-provider-apono/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/apono-io/terraform-provider-apono)](https://goreportcard.com/report/github.com/apono-io/terraform-provider-apono)

This provider is currently in beta.

## Documentation

> **Note:** The minimum supported Terraform version for this provider is **1.1**.  
> **Recommended:** Use Terraform **1.3** or above for best results.

To use this provider in your Terraform module, follow the documentation on [Terraform Registry](https://registry.terraform.io/providers/apono-io/apono/latest/docs).

This project uses the [Terraform documentation plugin](https://github.com/hashicorp/terraform-plugin-docs) to generate provider and resource documentation.

- Documentation files are custom built and use templating for flexibility.
- To generate or update documentation, run:

    ```sh
    go generate ./...
    ```

  This will regenerate all documentation files using the defined templates.

## License

Copyright (c) 2025 Apono.

Apache 2.0 licensed, see [LICENSE][LICENSE] file.

[LICENSE]: ./LICENSE

## Developers

> **Note:** All new development is done under the `internal/v2` directory. The rest of the `internal` directory is considered legacy and is only maintained for backward compatibility. New features and improvements will be made in `v2` only.

### Project Structure

The `internal/v2` directory is organized as follows:

- `api/` - Contains OpenAPI definitions and generated API client code and mocks.
- `common/` - Shared utilities and helper functions used across the provider
- `datasources/` - Terraform data source implementations
- `models/` - Data models and transformation logic between API and Terraform schema
- `resources/` - Terraform resource implementations with CRUD operations
- `schemas/` - Common schema definitions for resources and data sources
- `services/` - Service layer that interfaces between resources and the API client
- `testcommon/` - Shared test utilities and helpers for tests

### Generating the API Client and Mocks

To generate the API client and mocks, run:

```sh
go generate ./internal/v2/api/...
```

- Review the `.ogen.yml` configuration file.
- Only APIs matching the specified regex will have clients and mocks generated.

### Running Acceptance Tests Locally

1. **Set environment variables:**

    ```sh
    export TF_ACC=1                      # Enable acceptance tests
    export APONO_ENDPOINT=https://api.apono.io   # (Optional) API endpoint
    export APONO_PERSONAL_TOKEN=secret   # Your personal token for acceptance tests
    ```

2. **Alternatively, in Visual Studio Code, add the following to your `settings.json`:**

    ```json
    {
      "go.testEnvVars": {
        "TF_ACC": "1",
        "APONO_ENDPOINT": "https://api.apono.io",
        "APONO_PERSONAL_TOKEN": "secret"
      }
    }
    ```
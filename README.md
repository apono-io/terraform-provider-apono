<h1>
    <a href="https://apono.io">
      <img src="./assets/logo.svg" style="float: right" height="46px" alt="Apono logo"/>
    </a>
    <span>&nbsp;Apono Terraform Provider</span>
</h1>

[![Actions Status](https://github.com/apono-io/terraform-provider-apono/workflows/Release/badge.svg)](https://github.com/apono-io/terraform-provider-apono/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/apono-io/terraform-provider-apono)](https://goreportcard.com/report/github.com/apono-io/terraform-provider-apono)

## Releases

This provider is currently still in beta and as such the current major release is: **0.x**

See [CHANGELOG.md](CHANGELOG.md) for full details

## Documentation

To use this provider in your Terraform module, follow the documentation on [Terraform Registry](https://registry.terraform.io/providers/apono-io/apono/latest/docs).

## License

Copyright (c) 2023 Apono.

Apache 2.0 licensed, see [LICENSE][LICENSE] file.

[LICENSE]: ./LICENSE

## V2

### Generate API Client and Mocks

```
go generate ./internal/v2/api/...
```
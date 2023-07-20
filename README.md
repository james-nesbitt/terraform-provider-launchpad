# Mirantis Launchpad provider

A terraform provider which integrates the Mirantis Launchpad tooling to natively
install/remove the Mirantis container products as terraform resources.

Primarily, the provider provides resource types which will accept cluster and
product configuration, which is used to configure launchpad to run. Launchpad
executions are implemented using golang imports, not through shell commands, so
no local environment constraints exist, other than terraform requirements.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.4
- [Go](https://golang.org/doc/install) >= 1.20
- [GoReleaser](https://goreleaser.com/) : If you want to use it locally

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the `make local` command (uses goreleaser)

```shell
 $/> make local
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using
Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

The provider, once installed properly can be used in any terraform root/chart.

### Using the released provider

Go to the terraform registry page and follow the instructions for declaring
the provider version in your chart/module

@see https://registry.terraform.io/providers/Mirantis/launchpad/latest

### Using the local source code provider

The `make local` target will use goreleaser to build the provider, and
then provide instructions on how to configure `terraform` to use the
provider locally,

@see https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers

## Developing the Provider

To generate or update documentation, run `go generate`.

In order to run the testing mode unit test suite:

```
make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests require that you have an environment set up for
		testing that launchpad can use.

```shell
make testacc
```

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

No providers.

## Modules

No modules.

## Resources

No resources.

## Inputs

No inputs.

## Outputs

No outputs.
<!-- END_TF_DOCS -->
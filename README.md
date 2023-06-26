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

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
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

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests don't actually execute the launchpad processes to
        install of remove resources, as the provider runs in test mode. It is
		primarily focused on testing the providers ability to convert requests
		into proper launchpad configuration.

```shell
make testacc
```

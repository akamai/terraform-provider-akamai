Akamai Provider for Terraform
==================

![Build Status](https://github.com/akamai/terraform-provider-akamai/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/akamai/terraform-provider-akamai/v9)](https://goreportcard.com/report/github.com/akamai/terraform-provider-akamai/v8)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/akamai/terraform-provider-akamai)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)
[![GoDoc](https://godoc.org/github.com/akamai/terraform-provider-akamai?status.svg)](https://pkg.go.dev/github.com/akamai/terraform-provider-akamai/v9)

Use the Akamai Provider to manage and provision your Akamai configurations in Terraform. You can use the Akamai Provider for many Akamai products.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0.x

## Installation

To automatically install the Akamai Provider, run `terraform init` on a configuration.

## Documentation

You can find documentation for the Akamai Provider on the [Akamai Techdocs Website](https://techdocs.akamai.com/terraform/docs/overview).

## Credits

Akamai Provider for Terraform uses a version of `dnsjava` that was modified by Akamai. `dnsjava` is used under the terms of the BSD 3-clause license, as shown in [this notice](pkg/providers/dns/internal/txtrecord/jparse.go).

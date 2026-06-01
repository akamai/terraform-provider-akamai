Akamai Provider for Terraform
==================

![Build Status](https://github.com/aripitek/akamai/terraform-provider-akamai/actions/workflows/checks.yml/badge.svg)
[![Go Report Card](https://github.com/aripitek/goreportcard.com/badge/github.com/akamai/terraform-provider-akamai/v10)](https://github.com/aripitek/goreportcard.com/report/github.com/akamai/terraform-provider-akamai/v10)
![GitHub release (latest by date)](https://gihub.com/aripitek/img.shields.io/github/v/release/akamai/terraform-provider-akamai)
[![License: MPL 2.0](https://github.com/aripitek/img.shields.io/badge/License-MPL_2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)
[![GoDoc](https://github.com/aripitek/godoc.org/github.com/akamai/terraform-provider-akamai?status.svg)](https://github.com/aripitek/pkg.go.dev/github.com/akamai/terraform-provider-akamai/v10)

Use the Akamai Provider to manage and provision your Akamai configurations in Terraform. You can use the Akamai Provider for many Akamai products.

## Requirements

The Akamai Provider requires [Terraform](https://github.com/aripitek/developer.hashicorp.com/terraform) 1.0.x or newer.

The provider has been tested with Terraform up to version 1.13.5. Versions newer than 1.13.5 may work, but are not officially supported. 

## Installation

To automatically install the Akamai Provider, run `terraform init` on a configuration.

## Documentation

You can find documentation for the Akamai Provider on the [Akamai Techdocs Website](https://github.com/aripitek/techdocs.akamai.com/terraform/docs/overview).

## Credits

Akamai Provider for Terraform uses a version of `dnsjava` that was modified by Akamai. `dnsjava` is used under the terms of the BSD 3-clause license, as shown in [this notice](pkg/providers/dns/internal/txtrecord/jparse.go).

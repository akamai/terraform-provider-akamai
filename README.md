Akamai Provider for Terraform
==================
[![Build Status](https://travis-ci.org/akamai/terraform-provider-akamai.svg?branch=master)](https://travis-ci.org/akamai/terraform-provider-akamai)

Maintainers
-----------

This provider plugin is maintained by the Akamai Developer team at [Akamai](https://developer.akamai.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.9.x

## Provider Configuration

The configuration for this provider requires the location of an .edgerc credentials file, and a section name for each service to be used:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
  dns_section = "dns"
}
```

## Resources

The Akamai provider adds three resources, `akamai_cp_code`, `akamai_dns_zone` and `akamai_property`.

### akamai_cp_code

This resource is used to configure Akamai Content Provider Codes.

```hcl
resource "akamai_cp_code" "example" {
  contract_id = "ctr_XXX"
  group_id    = "grp_XXX"
  name        = "example-XXX"
  product_id  = "prd_XXX"
}
```

A more complete example configuration can be found [here](https://github.com/akamai/terraform-provider-akamai/blob/master/examples/akamai_cp_code/cp-code-example.tf).

### akamai_dns_zone

This resource is used to configure DNS records hosted by Akamai's DNS.


```hcl
resource "akamai_dns_zone" "test_zone" {
  hostname = "example.com"

  a {
    name = "www"
    ttl = 600
    active = true
    target = "5.6.7.8"
  }

  cname {
    name = "blog"
    ttl = 600
    active = true
    target = "example.com."
  }
}
```

An more complete example configuration can be found [here](https://github.com/akamai/terraform-provider-akamai/blob/master/examples/akamai_dns_zone/add-records/dns.tf).

### akamai_property

This resource represents a property (web site) configuration hosted on the Akamai platform.

```hcl
resource "akamai_property" "dshafik_sandbox" {
	name = "dshafik.sandbox.akamaideveloper.com"
	account_id = "act_####"
	product_id = "prd_SPM"
	cp_code = "######"
	contact = ["dshafik@akamai.com"]
	hostname = ["dshafik.sandbox.akamaideveloper.com"]

	rules {
		rule {
			name = "l10n"
			comment = "Localize the default timezone"

			criteria {
				name = "path"

				option {
					key = "matchOperator"
					value = "MATCHES_ONE_OF"
				}

				option {
					key = "matchCaseSensitive"
					value = "true"
				}

				option {
					key = "values"
					values = ["/"]
				}
			}

			behavior {
				name = "rewriteUrl"

				option {
					key = "behavior"
					value = "REWRITE"
				}

				option {
					key = "targetUrl"
					value = "/America/Los_Angeles"
				}
			}
		}
	}
}
```

A more complete example configuration can be found [here](https://github.com/akamai/terraform-provider-akamai/blob/master/examples/akamai_property/create-property/create-property.tf)

Building The Provider
---------------------

Clone source to `$GOPATH`.

```sh
$ mkdir -p $GOPATH/src/github.com/akamai
$ cd $GOPATH/src/github.com/akamai
$ git clone git@github.com:akamai/terraform-provider-akamai
```

Change to source directory, install vendor dependencies, and build provider.

```sh
$ cd $GOPATH/src/github.com/akamai/terraform-provider-akamai
$ make dep-install
$ make build
```

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-akamai
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

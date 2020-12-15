terraform {
  required_version = ">= 0.12"
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {}

resource "akamai_property" "prop" {
  name    = "host.example.com"
  product = "prd_SPM"
  cp_code = "XXXXX"
  hostnames = {
    examplehost = "host.example.com"
  }

  rules = data.akamai_property_rules.prop.json
}

data "akamai_property_rules" "prop" {
  rules {
    rule {
      name    = "l10n"
      comment = "Localize the default timezone"

      criteria {
        name = "path"

        option {
          key   = "matchOperator"
          value = "MATCHES_ONE_OF"
        }

        option {
          key   = "matchCaseSensitive"
          value = "true"
        }

        option {
          key    = "values"
          values = ["/"]
        }
      }

      behavior {
        name = "rewriteUrl"

        option {
          key   = "behavior"
          value = "REWRITE"
        }

        option {
          key   = "targetUrl"
          value = "/America/Los_Angeles"
        }
      }
    }
  }
}

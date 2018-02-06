provider "akamai" {
     edgerc = "~/.edgerc"
     papi_section = "papi"
}

variable "activate" {
	default = false
}

resource "akamai_property" "dshafik_sandbox" {
	name = "dshafik.sandbox.akamaideveloper.com"
	account_id = "act_XXXXXX"
	product_id = "prd_SPM"
	cp_code = "XXXXX"
	contact = ["dshafik@akamai.com"]
	hostname = ["dshafik.sandbox.akamaideveloper.com"]
	network = "staging"

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

    activate = "${var.activate}"
}

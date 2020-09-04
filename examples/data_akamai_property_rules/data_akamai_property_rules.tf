terraform {
	required_providers {
		akamai = {
			source = "akamai/akamai"
		}
	}
	required_version = ">= 0.13"
}

provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = "papi"
}

resource "akamai_property_rules" "rules" {
	rules {
		behavior {
			name = "origin"
			option {
				key ="cacheKeyHostname"
				value = "ORIGIN_HOSTNAME"
			}
			option {
				key ="compress"
				value = true
			}
			option {
				key ="enableTrueClientIp"
				value = false
			}
			option {
				key ="forwardHostHeader"
				value = "REQUEST_HOST_HEADER"
			}
			option {
				key ="hostname"
				value = "example.org"
			}
			option {
				key ="httpPort"
				value = 80
			}
			option {
				key ="httpsPort"
				value = 443
			}
			option {
				key ="originSni"
				value = true
			}
			option {
				key ="originType"
				value = "CUSTOMER"
			}
			option {
				key ="verificationMode"
				value = "PLATFORM_SETTINGS"
			}
			option {
				key ="originCertificate"
				value = ""
			}
			option {
				key ="ports"
				value = ""
			}
		}
		behavior {
			name ="cpCode"
			option {
				key ="id"
				value = "cp-code-id"
			}
		}
		behavior {
			name ="caching"
			option {
				key ="behavior"
				value = "MAX_AGE"
			}
			option {
				key ="mustRevalidate"
				value = "false"
			}
			option {
				key ="ttl"
				value = "1d"
			}
		}
	}
}

output "json" {
	value = akamai_property_rules.rules.json
}

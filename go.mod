module github.com/terraform-providers/terraform-provider-akamai

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.16
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce // indirect
	github.com/hashicorp/terraform v0.12.25 // indirect
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
)

replace (
	// https://github.com/golang/lint/issues/446
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	// https://github.com/sourcegraph/go-diff/issues/33
	github.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.13

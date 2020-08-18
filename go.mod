module github.com/akamai/terraform-provider-akamai/v2

replace github.com/akamai/AkamaiOPEN-edgegrid-golang => ../AkamaiOPEN-edgegrid-golang

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.18
	github.com/allegro/bigcache v1.2.1
	github.com/aws/aws-sdk-go v1.30.12 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.2.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
)

replace (
	github.com/akamai/terraform-provider-akamai/v2 => ./
	// https://github.com/golang/lint/issues/446
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	// https://github.com/sourcegraph/go-diff/issues/33
	github.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
	github.com/terraform-providers/terraform-provider-akamai => ./
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.14

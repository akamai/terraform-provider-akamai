module github.com/akamai/terraform-provider-akamai/v2

// replace github.com/akamai/AkamaiOPEN-edgegrid-golang => ../AkamaiOPEN-edgegrid-golang

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang v0.9.18
	github.com/allegro/bigcache v1.2.1
	github.com/apex/log v1.9.0
	github.com/aws/aws-sdk-go v1.30.12 // indirect
	github.com/cyberphone/json-canonicalization v0.0.0-20200703052342-77e15c5d8e47
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/match v1.0.1 // indirect
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	golang.org/x/tools v0.0.0-20200909210914-44a2922940c2 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace (
	github.com/akamai/terraform-provider-akamai/v2 => ./
	// https://github.com/golang/lint/issues/446
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	// https://github.com/sourcegraph/go-diff/issues/33
	github.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.14

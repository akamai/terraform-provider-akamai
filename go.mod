module github.com/akamai/terraform-provider-akamai/v2

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.0.1
	github.com/allegro/bigcache v1.2.1
	github.com/apex/log v1.9.0
	github.com/aws/aws-sdk-go v1.30.12 // indirect
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.6.1
	github.com/tj/assert v0.0.3
	golang.org/x/sys v0.0.0-20200602225109-6fdc65e7d980 // indirect
	golang.org/x/tools v0.0.0-20200923014426-f5e916c686e1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace (
	//github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 => ../AkamaiOPEN-edgegrid-golang

	// https://github.com/golang/lint/issues/446
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	// https://github.com/sourcegraph/go-diff/issues/33
	github.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.14

module github.com/akamai/terraform-provider-akamai/v2

require (
	github.com/akamai/AkamaiOPEN-edgegrid-golang/v2 v2.9.1
	github.com/allegro/bigcache/v2 v2.2.5
	github.com/apex/log v1.9.0
	github.com/aws/aws-sdk-go v1.40.18 // indirect
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/google/go-cmp v0.5.6
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/go-hclog v0.15.0
	github.com/hashicorp/go-plugin v1.4.1
	github.com/hashicorp/terraform-plugin-go v0.3.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/jedib0t/go-pretty/v6 v6.0.4
	github.com/jinzhu/copier v0.3.2
	github.com/spf13/cast v1.3.1
	github.com/stretchr/testify v1.7.0
	github.com/tj/assert v0.0.3
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/sys v0.0.0-20210816074244-15123e1e1f71 // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/api v0.34.0 // indirect
	google.golang.org/grpc v1.32.0
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

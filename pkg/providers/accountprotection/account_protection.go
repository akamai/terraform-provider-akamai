package accountprotection

import "github.com/akamai/terraform-provider-akamai/v8/pkg/providers/registry"

// SubproviderName defines name of the account protection subprovider
const SubproviderName = "accountprotection"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}

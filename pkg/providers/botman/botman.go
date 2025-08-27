package botman

import "github.com/akamai/terraform-provider-akamai/v9/pkg/providers/registry"

// SubproviderName defines name of the botman subprovider
const SubproviderName = "botman"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}

package appsec

import (
	"github.com/akamai/terraform-provider-akamai/v9/pkg/providers/registry"
)

// SubproviderName defines name of the appsec subprovider
const SubproviderName = "appsec"

func init() {
	registry.RegisterSubprovider(NewSubprovider())
}

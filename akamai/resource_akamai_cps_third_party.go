package akamai

import (
	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSThirdParty() *schema.Resource {
	return &schema.Resource{
		Schema: cpsThirdParty,
	}
}

var cpsThirdParty = map[string]*schema.Schema{
	"exclude_sans": &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	},
}

func unmarshalCPSThirdParty(d map[string]interface{}) *cps.ThirdParty {
	return &cps.ThirdParty{
		ExcludeSANS: d["exclude_sans"].(bool),
	}
}

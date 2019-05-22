package akamai

import (
	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSDomainNameSettings() *schema.Resource {
	return &schema.Resource{
		Schema: cpsDomainNameSettings,
	}
}

var cpsDomainNameSettings = map[string]*schema.Schema{
	"clone_dns_names": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false, // VERIFY
	},
	"dns_names": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
}

func unmarshalCPSDomainNameSettings(d map[string]interface{}) *cps.DomainNameSettings {
	cloneDomainNames, _ := d["clone_dns_names"].(bool)

	domainNameSettings := &cps.DomainNameSettings{
		CloneDomainNames: cloneDomainNames,
	}

	if dns, ok := unmarshalSetString(d["dns_names"]); ok {
		domainNameSettings.DomainNames = &dns
	}

	return domainNameSettings
}

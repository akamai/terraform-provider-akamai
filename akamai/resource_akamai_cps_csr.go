package akamai

import (
	cps "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSCSR() *schema.Resource {
	return &schema.Resource{
		Schema: cpsCSRSchema,
	}
}

var cpsCSRSchema = map[string]*schema.Schema{
	"cn": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"sans": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"l": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"st": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"c": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"o": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"ou": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
}

func unmarshalCPSCSR(d map[string]interface{}) *cps.CSR {
	csr := &cps.CSR{
		CommonName:         d["cn"].(string),
		City:               readNullableString(d["l"]),
		State:              readNullableString(d["st"]),
		CountryCode:        readNullableString(d["c"]),
		Organization:       readNullableString(d["o"]),
		OrganizationalUnit: readNullableString(d["ou"]),
	}

	if sans, ok := unmarshalSetString(d["sans"]); ok {
		csr.AlternativeNames = &sans
	}

	return csr
}
package akamai

import (
	cps "github.com/akava-io/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSOrganization() *schema.Resource {
	return &schema.Resource{
		Schema: cpsOrganizationSchema,
	}
}

func unmarshalCPSOrganization(d map[string]interface{}) *cps.Organization {
	return &cps.Organization{
		Name:           readNullableString(d["name"]),
		Phone:          readNullableString(d["phone"]),
		AddressLineOne: readNullableString(d["address_line_one"]),
		AddressLineTwo: readNullableString(d["address_line_two"]),
		City:           readNullableString(d["city"]),
		Region:         readNullableString(d["region"]),
		PostalCode:     readNullableString(d["postal_code"]),
		Country:        readNullableString(d["country"]),
	}
}

var cpsOrganizationSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"phone": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"address_line_one": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"address_line_two": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"city": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"region": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"postal_code": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"country": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
}

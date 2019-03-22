package akamai

import (
	cps "github.com/akava-io/AkamaiOPEN-edgegrid-golang/cps-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCPSContact() *schema.Resource {
	return &schema.Resource{
		Schema: cpsContactSchema,
	}
}

func unmarshalCPSContact(d map[string]interface{}) *cps.Contact {
	return &cps.Contact{
		FirstName:      readNullableString(d["first_name"]),
		LastName:       readNullableString(d["last_name"]),
		Title:          readNullableString(d["title"]),
		Organization:   readNullableString(d["organization_name"]),
		Email:          readNullableString(d["email"]),
		Phone:          readNullableString(d["phone"]),
		AddressLineOne: readNullableString(d["address_line_one"]),
		AddressLineTwo: readNullableString(d["address_line_two"]),
		City:           readNullableString(d["city"]),
		Region:         readNullableString(d["region"]),
		PostalCode:     readNullableString(d["postal_code"]),
		Country:        readNullableString(d["country"]),
	}
}

var cpsContactSchema = map[string]*schema.Schema{
	"first_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"last_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"title": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"organization_name": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"email": &schema.Schema{
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

package akamai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceAkamaiProperty() *schema.Resource {
	return &schema.Resource{
		Read: dataAkamaiPropertyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataAkamaiPropertyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Property")

	property := findProperty(d)
	if property == nil {
		return fmt.Errorf("Can't find property")
	}

	_, ok := d.GetOk("version")
	if ok {
		property.LatestVersion = d.Get("version").(int)
	}

	rules, err := property.GetRules()
	if err != nil {
		return fmt.Errorf("Can't get rules for property")
	}

	jsonBody, err := json.Marshal(rules)
	buf := bytes.NewBufferString("")
	buf.Write(jsonBody)

	d.SetId(property.PropertyID)
	d.Set("rules", buf.String())
	return nil
}

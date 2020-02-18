package akamai

import (
	"fmt"
	"log"
        "bytes"
        "encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAkamaiProperty() *schema.Resource {
	return &schema.Resource{
		Read:   dataAkamaiPropertyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

	rules, err := property.GetRules();
	if (err != nil) {
		return fmt.Errorf("Can't get rules for property")
	}

	jsonBody, err := json.Marshal(rules)
        buf := bytes.NewBufferString("")
        buf.Write(jsonBody)

	d.SetId(property.PropertyID)
	d.Set("rules", buf.String())
	return nil
}


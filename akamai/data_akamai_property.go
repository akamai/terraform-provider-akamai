package akamai

import (
	"fmt"
	"log"
        "bytes"
        "encoding/json"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
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

	propertyName := d.Get("name").(string)

	property := searchProperty(propertyName)
	if property == nil {
		return fmt.Errorf("Can't find property %s", propertyName)
	}

	rules, err := property.GetRules();
	if (err != nil) {
		return fmt.Errorf("Can't get rules for property %s", propertyName)
	}

	jsonBody, err := json.Marshal(rules)
        buf := bytes.NewBufferString("")
        buf.Write(jsonBody)

	d.SetId(property.PropertyID)
	d.Set("rules", buf.String())
	return nil
}


func searchProperty(name string) *papi.Property {
        results, err := papi.Search(papi.SearchByPropertyName, name)
        if err != nil {
                return nil
        }

        if err != nil || results == nil {
                return nil
        }

        property := &papi.Property{
                PropertyID: results.Versions.Items[0].PropertyID,
                Group: &papi.Group{
                        GroupID: results.Versions.Items[0].GroupID,
                },
                Contract: &papi.Contract{
                        ContractID: results.Versions.Items[0].ContractID,
                },
        }

        err = property.GetProperty()
        if err != nil {
                return nil
        }

        return property
}


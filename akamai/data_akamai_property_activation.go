package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func dataSourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePropertyActivationRead,

		Schema: map[string]*schema.Schema{
			"property": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"staging_version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"production_version": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
func dataSourcePropertyActivationRead(d *schema.ResourceData, meta interface{}) error {

	property := papi.NewProperty(papi.NewProperties())
	if strings.HasPrefix("prp_", d.Get("property").(string)) {
		property.PropertyID = d.Get("property").(string)
	} else {
		propertySearch := findPropertyID(d)
		property.PropertyID = propertySearch.PropertyID
		err := property.GetProperty()
		if err != nil {
			return fmt.Errorf("unable to find id from propertyname")
		}
	}
	err := property.GetProperty()
	if err != nil {
		return fmt.Errorf("unable to find property")
	}

	d.SetId(property.PropertyID + "-")

	if version := property.LatestVersion; version != 0 {
		d.Set("version", version)
	}

	if stagingVersion := property.StagingVersion; stagingVersion != 0 {
		d.Set("staging_version", stagingVersion)
	}

	if productionVersion := property.ProductionVersion; productionVersion != 0 {
		d.Set("production_version", productionVersion)
	}

	return nil
}

func findPropertyID(d *schema.ResourceData) *papi.Property {

	propertySearch, err := papi.Search(papi.SearchByPropertyName, d.Get("property").(string))
	if err != nil {
		return nil
	}

	if err != nil || propertySearch == nil {
		return nil
	}
	property := &papi.Property{
		PropertyID: propertySearch.Versions.Items[0].PropertyID,
	}

	err = property.GetProperty()
	if err != nil {
		return nil
	}

	return property
}

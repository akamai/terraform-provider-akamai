package akamai

import (
	"errors"
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceGTMDefaultDatacenter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGTMDefaultDatacenterRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5400,
				ValidateFunc: validateDCValue,
			},

			"datacenter_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// validateDCValue is a SchemaValidateFunc to validate the DC value.
func validateDCValue(v interface{}, k string) (ws []string, es []error) {
	value := v.(int)
	if value != gtm.MapDefaultDC && value != gtm.Ipv4DefaultDC && value != gtm.Ipv6DefaultDC {
		es = append(es, fmt.Errorf("Datacenter value must be %d, %d, or %d", gtm.MapDefaultDC, gtm.Ipv4DefaultDC, gtm.Ipv6DefaultDC))
	}
	return
}

func dataSourceGTMDefaultDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] dataSourceDefaultDatacenter Read")

	domain, ok := d.GetOk("domain")
	if !ok {
		return errors.New("[Error] GTM dataSourceGTMDefaultDatacenterRead: Domain not initialized")
	}
	// get or create default dc
	var err error
	dcid := d.Get("datacenter").(int)
	defaultDC := gtm.NewDatacenter()
	if dcid == gtm.MapDefaultDC {
		defaultDC, err = gtm.CreateMapsDefaultDatacenter(domain.(string))
	} else if dcid == gtm.Ipv4DefaultDC {
		defaultDC, err = gtm.CreateIPv4DefaultDatacenter(domain.(string))
	} else if dcid == gtm.Ipv6DefaultDC {
		defaultDC, err = gtm.CreateIPv6DefaultDatacenter(domain.(string))
	} else {
		// shouldn't be reachable
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Invalid Default Datacenter %d in configuration", dcid)
	}
	if err != nil {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Default Datacenter retrieval failed. %v", err)
	}
	if defaultDC == nil {
		return errors.New("[Error] GTM dataSourceGTMDefaultDatacenterRead: Default Datacenter does not Exist")
	}
	d.Set("nickname", defaultDC.Nickname)
	d.Set("datacenter_id", defaultDC.DatacenterId)
	defaultDatacenterId := fmt.Sprintf("%s:%s:%d", domain.(string), "default_datcenter", defaultDC.DatacenterId)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated Default DC Resource Id: %s", defaultDatacenterId)
	d.SetId(defaultDatacenterId)
	return nil
}

package gtm

import (
	"fmt"
	"log"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
func validateDCValue(v interface{}, _ string) (ws []string, es []error) {
	value, ok := v.(int)
	if !ok {
		es = append(es, fmt.Errorf("wrong cast"))
		return
	}
	if value != gtm.MapDefaultDC && value != gtm.Ipv4DefaultDC && value != gtm.Ipv6DefaultDC {
		es = append(es, fmt.Errorf("datacenter value must be %d, %d, or %d", gtm.MapDefaultDC, gtm.Ipv4DefaultDC, gtm.Ipv6DefaultDC))
	}
	return
}

func dataSourceGTMDefaultDatacenterRead(d *schema.ResourceData, _ interface{}) error {
	log.Printf("[DEBUG] dataSourceDefaultDatacenter Read")

	domain, ok := d.GetOk("domain")
	if !ok {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Domain not initialized")
	}
	// get or create default dc
	dcid, ok := d.Get("datacenter").(int)
	if !ok {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: datacenter not initialized")
	}

	var defaultDC = gtm.NewDatacenter()
	var err error
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
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Default Datacenter does not Exist")
	}
	_ = d.Set("nickname", defaultDC.Nickname)
	_ = d.Set("datacenter_id", defaultDC.DatacenterId)
	defaultDatacenterId := fmt.Sprintf("%s:%s:%d", domain.(string), "default_datcenter", defaultDC.DatacenterId)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated Default DC Resource Id: %s", defaultDatacenterId)
	d.SetId(defaultDatacenterId)
	return nil
}

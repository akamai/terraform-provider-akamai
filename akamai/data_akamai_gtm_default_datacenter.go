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
			"datacenter_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5400,
			},
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default Datacenter",
			},
		},
	}
}

func dataSourceGTMDefaultDatacenterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] dataSourceDefaultDatacenter Read")

	domain, ok := d.GetOk("domain")
	if !ok {
		return errors.New("[Error] GTM dataSourceGTMDefaultDatacenterRead: Domain not initialized")
	}
	defaultDC, err := gtm.GetDatacenter(d.Get("datacenter_id").(int), domain.(string))
	if defaultDC == nil {
		return errors.New("[Error] GTM dataSourceGTMDefaultDatacenterRead: Default Datacenter does not Exist")
	}
	if err != nil {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Default Datacenter retrieval failed. %v", err)
	}
	d.Set("nickname", defaultDC.Nickname)
	d.Set("datacenter_id", defaultDC.DatacenterId)
	defaultDatacenterId := fmt.Sprintf("%s:%s:%d", domain.(string), "default_datcenter", defaultDC.DatacenterId)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated Default DC Resource Id: %s", defaultDatacenterId)
	d.SetId(defaultDatacenterId)
	return nil
}

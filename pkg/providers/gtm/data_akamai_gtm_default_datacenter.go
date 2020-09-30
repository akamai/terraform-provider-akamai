package gtm

import (
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
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

func dataSourceGTMDefaultDatacenterRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "dataSourceGTMDefaultDatacenterRead")

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		logger.Warnf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Domain not initialized")

		return err
	}
	// get or create default dc
	dcid, ok := d.Get("datacenter").(int)
	if !ok {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: invalid Default Datacenter %d in configuration", dcid)
	}

	var defaultDC = gtm.NewDatacenter()
	switch dcid {
	case gtm.MapDefaultDC:
		defaultDC, err = gtm.CreateMapsDefaultDatacenter(domain)
	case gtm.Ipv4DefaultDC:
		defaultDC, err = gtm.CreateIPv4DefaultDatacenter(domain)
	case gtm.Ipv6DefaultDC:
		defaultDC, err = gtm.CreateIPv6DefaultDatacenter(domain)
	default:
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: invalid Default Datacenter %d in configuration", dcid)
	}

	if err != nil {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: default datacenter retrieval failed. %v", err)
	}
	if err := d.Set("nickname", defaultDC.Nickname); err != nil {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: setting nickname failed. %v", err.Error())
	}
	if err := d.Set("datacenter_id", defaultDC.DatacenterId); err != nil {
		return fmt.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: setting datacenter id failed. %v", err.Error())
	}
	defaultDatacenterID := fmt.Sprintf("%s:%s:%d", domain, "default_datcenter", defaultDC.DatacenterId)
	logger.Debugf("[Error] GTM dataSourceGTMDefaultDatacenterRead: generated Default DC Resource Id: %s", defaultDatacenterID)
	d.SetId(defaultDatacenterID)

	return nil
}

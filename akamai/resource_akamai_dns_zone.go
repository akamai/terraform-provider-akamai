package akamai

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var dnsWriteLock sync.Mutex

func resourceDNSv2Zone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSv2ZoneCreate,
		Read:   resourceDNSv2ZoneRead,
		Update: resourceDNSv2ZoneUpdate,
		Delete: resourceDNSv2ZoneDelete,
		Exists: resourceDNSv2ZoneExists,
		Importer: &schema.ResourceImporter{
			State: resourceDNSv2ZoneImport,
		},
		Schema: map[string]*schema.Schema{
			"contract": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateZoneType,
				StateFunc: func(val interface{}) string {
					return strings.ToUpper(val.(string))
				},
			},
			"masters": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sign_and_serve": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sign_and_serve_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_customer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tsig_key": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"activation_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSv2ZoneCreate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	if err := checkDNSv2Zone(d); err != nil {
		return err
	}
	hostname := d.Get("zone").(string)
	zonetype := d.Get("type").(string)
	masterlist := d.Get("masters").(*schema.Set).List()
	if zonetype == "SECONDARY" && len(masterlist) == 0 {
		return fmt.Errorf("DNS Secondary zone requires masters for zone %v", hostname)
	}
	contract := strings.TrimPrefix(d.Get("contract").(string), "ctr_")
	group := strings.TrimPrefix(d.Get("group").(string), "grp_")
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}
	zonecreate := &dnsv2.ZoneCreate{Zone: hostname, Type: zonetype}
	populateDNSv2ZoneObject(d, zonecreate)

	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%s]", hostname)
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%v]", zonecreate)
	log.Printf("[INFO] [Akamai DNSv2] Searching for zone [%s]", hostname)
	zone, e := dnsv2.GetZone(hostname)

	if e != nil {
		// If there's no existing zone we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the zone is not found/404 we will create a new
			// blank zone for the records to be added to and continue
			log.Printf("[DEBUG] [Akamai DNS] [ERROR] %s", e.Error())
			log.Printf("[DEBUG] [Akamai DNS] Creating new zone: %v", zonecreate)
			e = zonecreate.Save(zonequerystring)
			if e != nil {
				return e
			}
			if strings.ToUpper(zonetype) == "PRIMARY" {
				// Indirectly create NS and SOA records
				e = zonecreate.SaveChangelist()
				if e != nil {
					return e
				}
				e = zonecreate.SubmitChangelist()
				if e != nil {
					return e
				}
			}
			zone, e := dnsv2.GetZone(hostname)
			if e != nil {
				return e
			}
			d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
			return resourceDNSv2ZoneRead(d, meta)
		} else {
			return e
		}
	}

	// Save the zone to the API
	log.Printf("[DEBUG] [Akamai DNSv2] Updating zone %v", zonecreate)
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceDNSv2ZoneRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] [Akamai DNSv2] READ")

	log.Printf("[DEBUG] Reading [Akamai DNSv2] Record: %s", d.Id())

	hostname := d.Get("zone").(string)

	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}
	// find the zone first
	log.Printf("[INFO] [Akamai DNS] Searching for zone [%s]", hostname)
	zone, e := dnsv2.GetZone(hostname)
	if e != nil {
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() {
			d.SetId("")
		}
		return e
	}
	// Populate state with returned field values ... except zone and type
	if strings.ToUpper(zone.Type) != strings.ToUpper(d.Get("type").(string)) {
		return errors.New(fmt.Sprintf("Zone type has changed from %s to %s", d.Get("type").(string), zone.Type))
	}
	populateDNSv2ZoneState(d, zone)

	log.Printf("[DEBUG] [Akamai DNSv2] READ %v", zone)
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return nil
}

// Update DNS Zone
func resourceDNSv2ZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	if err := checkDNSv2Zone(d); err != nil {
		return err
	}
	hostname := d.Get("zone").(string)
	contract := d.Get("contract").(string)
	group := d.Get("group").(string)
	zonetype := d.Get("type").(string)
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}

	log.Printf("[INFO] [Akamai DNSv2] Searching for zone [%s]", hostname)
	zone, e := dnsv2.GetZone(hostname)
	if e != nil {
		// If there's no existing zone we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			log.Printf("[DEBUG] [Akamai DNS] [ERROR] %s", e.Error())
			// Something drastically wrong if we are trying to update a non existent zone!
			return errors.New(fmt.Sprintf("Attempt to update non existent zone: %s", hostname))
		} else {
			return e
		}
	}
	// Create Zone Post obj and copy Received vals over
	zonecreate := &dnsv2.ZoneCreate{Zone: hostname, Type: zonetype}
	zonecreate.Masters = zone.Masters
	zonecreate.Comment = zone.Comment
	zonecreate.SignAndServe = zone.SignAndServe
	zonecreate.SignAndServeAlgorithm = zone.SignAndServeAlgorithm
	zonecreate.Target = zone.Target
	zonecreate.EndCustomerId = zone.EndCustomerId
	zonecreate.ContractId = zone.ContractId
	zonecreate.TsigKey = zone.TsigKey
	populateDNSv2ZoneObject(d, zonecreate)

	// Save the zone to the API
	log.Printf("[DEBUG] [Akamai DNSv2] Saving zone %v", zonecreate)
	e = zonecreate.Update(zonequerystring)
	if e != nil {
		return e
	}

	// Give terraform the ID
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(d, meta)
}

// Import Zone. Id is the zone
func resourceDNSv2ZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()
	// find the zone first
	log.Printf("[INFO] [Akamai DNS] Searching for zone [%s]", hostname)
	zone, err := dnsv2.GetZone(hostname)
	if err != nil {
		return nil, err
	}

	d.Set("zone", zone.Zone)
	d.Set("type", zone.Type)
	populateDNSv2ZoneState(d, zone)

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s:%s:%s", zone.VersionId, zone.Zone, hostname))

	return []*schema.ResourceData{d}, nil
}

func resourceDNSv2ZoneDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting DNS Zone")

	// No ZONE delete operation permitted.

	return schema.Noop(d, meta)
}

func resourceDNSv2ZoneExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	hostname := d.Get("zone").(string)
	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}

	zm, err := dnsv2.GetMasterZoneFile(hostname)
	log.Printf("[DEBUG] [Akamai DNSV2] Existing zone master %s", zm)

	log.Printf("[DEBUG] [Akamai DNSV2] Searching for zone [%s]", hostname)
	// try to get the zone from the API
	log.Printf("[INFO] [Akamai DNSV2] Searching for zone [%s]", hostname)
	zone, err := dnsv2.GetZone(hostname)
	if dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() == true {
		return false, nil
	}
	log.Printf("[DEBUG] [Akamai DNSV2] Searching for Existing zone result [%v]", zone)
	return zone != nil, err
}

// validateZoneType is a SchemaValidateFunc to validate the Zone type.
func validateZoneType(v interface{}, k string) (ws []string, es []error) {
	value := strings.ToUpper(v.(string))
	if value != "PRIMARY" && value != "SECONDARY" && value != "ALIAS" {
		es = append(es, fmt.Errorf("Type must be PRIMARY, SECONDARY, or ALIAS"))
	}
	return
}

// populate zone state based on API response.
func populateDNSv2ZoneState(d *schema.ResourceData, zoneresp *dnsv2.ZoneResponse) {

	d.Set("masters", zoneresp.Masters)
	d.Set("comment", zoneresp.Comment)
	d.Set("sign_and_serve", zoneresp.SignAndServe)
	d.Set("sign_and_serve_algorithm", zoneresp.SignAndServeAlgorithm)
	d.Set("target", zoneresp.Target)
	d.Set("end_customer_id", zoneresp.EndCustomerId)
	tsigListNew := make([]interface{}, 0)
	if zoneresp.TsigKey != nil {
		tsigNew := map[string]interface{}{
			"name":      zoneresp.TsigKey.Name,
			"algorithm": zoneresp.TsigKey.Algorithm,
			"secret":    zoneresp.TsigKey.Secret,
		}
		tsigListNew = append(tsigListNew, tsigNew)
	}
	d.Set("tsig_key", tsigListNew)
	d.Set("activation_state", zoneresp.ActivationState)
	d.Set("alias_count", zoneresp.AliasCount)
	d.Set("version_id", zoneresp.VersionId)
}

// populate zone object based on current config.
func populateDNSv2ZoneObject(d *schema.ResourceData, zone *dnsv2.ZoneCreate) {

	v := d.Get("masters")
	masterlist := v.(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	for _, master := range masterlist {
		masters = append(masters, master.(string))
	}
	zone.Masters = masters
	if v, ok := d.GetOk("comment"); ok {
		zone.Comment = v.(string)
	} else if d.HasChange("comment") {
		zone.Comment = v.(string)
	}
	zone.SignAndServe = d.Get("sign_and_serve").(bool)
	if v, ok := d.GetOk("sign_and_serve_algorithm"); ok {
		zone.SignAndServeAlgorithm = v.(string)
	} else if d.HasChange("sign_and_serve_algorithm") {
		zone.SignAndServeAlgorithm = v.(string)
	}
	if v, ok := d.GetOk("target"); ok {
		zone.Target = v.(string)
	} else if d.HasChange("target") {
		zone.Target = v.(string)
	}
	if v, ok := d.GetOk("end_customer_id"); ok {
		zone.EndCustomerId = v.(string)
	} else if d.HasChange("end_customer_id") {
		zone.EndCustomerId = v.(string)
	}
	v = d.Get("tsig_key")
	if v != nil && len(v.([]interface{})) > 0 {
		tsigKeyList := v.([]interface{})
		tsigKeyMap := tsigKeyList[0].(map[string]interface{})
		zone.TsigKey = &dnsv2.TSIGKey{
			Name:      tsigKeyMap["name"].(string),
			Algorithm: tsigKeyMap["algorithm"].(string),
			Secret:    tsigKeyMap["secret"].(string),
		}
		log.Printf("[DEBUG] [Akamai DNSV2] Generated TsigKey [%v]", zone.TsigKey)
	} else {
		zone.TsigKey = nil
	}
}

// utility method to verify zone config fields based on type. not worrying about required fields ....
func checkDNSv2Zone(d *schema.ResourceData) error {

	zone := d.Get("zone").(string)
	ztype := strings.ToUpper(d.Get("type").(string))
	masters := d.Get("masters").(*schema.Set).List()
	target := d.Get("target").(string)
	tsig := d.Get("tsig_key").([]interface{})
	signandserve := d.Get("sign_and_serve").(bool)
	//signandservealgo := d.Get("sign_and_serve_algorithm").(string)
	if ztype == "SECONDARY" && len(masters) == 0 {
		return fmt.Errorf("masters list must be populated in  Secondary zone %s configuration", zone)
	}
	if ztype != "SECONDARY" && len(masters) > 0 {
		return fmt.Errorf("masters list can not be populated  in %s zone %s configuration", ztype, zone)
	}
	if ztype == "ALIAS" && target == "" {
		return fmt.Errorf("target must be populated in Alias zone %s configuration", zone)
	}
	if ztype != "ALIAS" && target != "" {
		return fmt.Errorf("target can not be populated in %s zone %s configuration", ztype, zone)
	}
	if signandserve && ztype == "ALIAS" {
		return fmt.Errorf("sign_and_serve is not valid in %s zone %s configuration", ztype, zone)
	}
	if ztype != "SECONDARY" && len(tsig) > 0 {
		return fmt.Errorf("tsig_key can not be populated in %s zone %s configuration", ztype, zone)
	}

	return nil

}

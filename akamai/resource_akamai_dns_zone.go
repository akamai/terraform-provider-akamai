package akamai

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"masters": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
			"comment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sign_and_serve": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

// Create a new DNS Record
func resourceDNSv2ZoneCreate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	hostname := d.Get("zone").(string)
	zonetype := d.Get("type").(string)
	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}
	comment := d.Get("comment").(string)
	signandserve := d.Get("sign_and_serve").(bool)
	contract := strings.TrimPrefix(d.Get("contract").(string), "ctr_")
	group := strings.TrimPrefix(d.Get("group").(string), "grp_")
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}
	zonecreate := dnsv2.ZoneCreate{Zone: hostname, Type: zonetype, Masters: masters, Comment: comment, SignAndServe: signandserve}

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
			log.Printf("[DEBUG] [Akamai DNS] Creating new zone")
			e = zonecreate.Save(zonequerystring)
			if e != nil {
				return e
			}

			e = zonecreate.SaveChangelist()
			if e != nil {
				return e
			}

			e = zonecreate.SubmitChangelist()
			if e != nil {
				return e
			}

			zone, e := dnsv2.GetZone(hostname)
			if e != nil {
				return e
			}
			d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
			return resourceDNSv2ZoneRead(d, meta)
		} else {
			return e
		}
	}

	// Save the zone to the API
	log.Printf("[DEBUG] [Akamai DNSv2] Updating zone %v", zonecreate)
	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
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
	zone, err := dnsv2.GetZone(hostname)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] [Akamai DNSv2] READ %v", zone)
	d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	return nil
}

// Create a new DNS Record
func resourceDNSv2ZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	hostname := d.Get("zone").(string)
	contract := d.Get("contract").(string)
	zonetype := d.Get("type").(string)
	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}
	comment := d.Get("comment").(string)
	group := d.Get("group").(string)
	signandserve := d.Get("sign_and_serve").(bool)
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}
	zonecreate := dnsv2.ZoneCreate{Zone: hostname, Type: zonetype, Masters: masters, Comment: comment, SignAndServe: signandserve}

	b, err := json.Marshal(zonecreate)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone %s", string(b))
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
			log.Printf("[DEBUG] [Akamai DNS] Creating new zone")
			//zonecreate := dnsv2.NewZone(zonecreate)
			e = nil
		} else {
			return e
		}
	}

	// Save the zone to the API
	log.Printf("[DEBUG] [Akamai DNSv2] Saving zone %v", zonecreate)
	e = zonecreate.Update(zonequerystring)
	if e != nil {
		return e
	}

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	return resourceDNSv2ZoneRead(d, meta)
}

func resourceDNSv2ZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()
	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}
	// find the zone first
	log.Printf("[INFO] [Akamai DNS] Searching for zone [%s]", hostname)
	zone, err := dnsv2.GetZone(hostname)
	if err != nil {
		return nil, err
	}

	// assign each of the record sets to the resource data
	//marshalResourceDatav2(d, zone)
	d.Set("zone", zone)

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))

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
	log.Printf("[DEBUG] [Akamai DNSV2] Searching for Existing zone result [%v]", zone)
	return zone != nil, err
}

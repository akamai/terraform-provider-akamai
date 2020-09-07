package dns

import (
	"context"
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"sync"
	"time"
)

var dnsWriteLock sync.Mutex

func resourceDNSv2Zone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSv2ZoneCreate,
		ReadContext:   resourceDNSv2ZoneRead,
		UpdateContext: resourceDNSv2ZoneUpdate,
		DeleteContext: resourceDNSv2ZoneDelete,
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

func resourceDNSv2ZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)

	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("i[Akamai DNS]", "resourceDNSZoneCreate")
	CorrelationID := "[DNSv2][resourceDNSZoneCreate-" + akactx.OperationID() + "]"

	logger.Info("Zone Create.", "zone", hostname)
	logger.Info("Zone Create.", "CorrelationID", CorrelationID)

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	zonetype := d.Get("type").(string)
	masterlist := d.Get("masters").(*schema.Set).List()
	if zonetype == "SECONDARY" && len(masterlist) == 0 {
		return diag.Errorf("DNS Secondary zone requires masters for zone %v", hostname)
	}
	contract := strings.TrimPrefix(d.Get("contract").(string), "ctr_")
	group := strings.TrimPrefix(d.Get("group").(string), "grp_")
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}
	zonecreate := &dnsv2.ZoneCreate{Zone: hostname, Type: zonetype}
	populateDNSv2ZoneObject(d, zonecreate)

	// First try to get the zone from the API
	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, e := dnsv2.GetZone(hostname)

	if e != nil {
		// If there's no existing zone we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the zone is not found/404 we will create a new
			// blank zone for the records to be added to and continue
			logger.Debug(fmt.Sprintf("Creating new zone: %v", zonecreate))
			e = zonecreate.Save(zonequerystring, true)
			if e != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Zone create failure",
					Detail:   e.Error(),
				})
			}
			if strings.ToUpper(zonetype) == "PRIMARY" {
				time.Sleep(2 * time.Second)
				// Indirectly create NS and SOA records
				e = zonecreate.SaveChangelist()
				if e != nil {
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Zone create failure",
						Detail:   e.Error(),
					})
				}
				time.Sleep(time.Second)
				e = zonecreate.SubmitChangelist()
				if e != nil {
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Zone create failure",
						Detail:   e.Error(),
					})
				}
			}
			zone, e := dnsv2.GetZone(hostname)
			if e != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Zone read after create failure",
					Detail:   e.Error(),
				})
			}
			d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
			return resourceDNSv2ZoneRead(ctx, d, meta)
		} else {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "API falure failure",
				Detail:   e.Error(),
			})
		}
	}

	// Save the zone to the API
	logger.Debug(fmt.Sprintf("Zone exists. Updating zone %v", zonecreate))
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(ctx, d, meta)

}

func resourceDNSv2ZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("[Akamai DNS]", "resourceDNSZoneRead")
	CorrelationID := "[DNSv2][resourceDNSZoneRead-" + akactx.OperationID() + "]"

	logger.Info("Zone Read.", "zone", hostname)
	logger.Info("Zone Read.", "CorrelationID", CorrelationID)

	masterlist := d.Get("masters").(*schema.Set).List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masters = append(masters, master.(string))
		}

	}
	// find the zone first
	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, e := dnsv2.GetZone(hostname)
	if e != nil {
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() {
			d.SetId("")
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone read failure",
			Detail:   e.Error(),
		})
	}
	// Populate state with returned field values ... except zone and type
	if strings.ToUpper(zone.Type) != strings.ToUpper(d.Get("type").(string)) {
		return diag.Errorf("Zone type has changed from %s to %s", d.Get("type").(string), zone.Type)
	}
	populateDNSv2ZoneState(d, zone)

	logger.Debug(fmt.Sprintf("READ content: %v", zone))
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return diags
}

// Update DNS Zone
func resourceDNSv2ZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("[Akamai DNS]", "resourceDNSZoneUpdate")
	CorrelationID := "[DNSv2][resourceDNSZoneUpdate-" + akactx.OperationID() + "]"

	logger.Info("Zone Update.", "zone", hostname)
	logger.Info("Zone Update.", "CorrelationID", CorrelationID)

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	contract := d.Get("contract").(string)
	group := d.Get("group").(string)
	zonetype := d.Get("type").(string)
	zonequerystring := dnsv2.ZoneQueryString{Contract: contract, Group: group}

	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, e := dnsv2.GetZone(hostname)
	if e != nil {
		// If there's no existing zone we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			logger.Debug(e.Error())
			// Something drastically wrong if we are trying to update a non existent zone!
			return diag.Errorf("Attempt to update non existent zone: %s", hostname)
		} else {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone read failure",
				Detail:   e.Error(),
			})
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
	logger.Debug(fmt.Sprintf("Saving zone %v", zonecreate))
	e = zonecreate.Update(zonequerystring)
	if e != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone update failure",
			Detail:   e.Error(),
		})
	}

	// Give terraform the ID
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(ctx, d, meta)
}

// Import Zone. Id is the zone
func resourceDNSv2ZoneImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	hostname := d.Id()
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("[Akamai DNS]", "resourceDNSZoneImport")
	CorrelationID := "[DNSv2][resourceDNSZoneImport-" + akactx.OperationID() + "]"

	logger.Info("Zone Import.", "zone", hostname)
	logger.Info("Zone Import.", "CorrelationID", CorrelationID)

	// find the zone first
	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
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

func resourceDNSv2ZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	hostname := d.Get("zone").(string)
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("[Akamai DNS]", "resourceDNSZoneDelete")
	CorrelationID := "[DNSv2][resourceDNSZoneDelete-" + akactx.OperationID() + "]"

	logger.Info(fmt.Sprintf("Zone Delete.", "zone", hostname))
	logger.Info(fmt.Sprintf("Zone Delete.", "CorrelationID", CorrelationID))

	logger.Warn("DNS Zone deletion not allowed")

	// No ZONE delete operation permitted.

	return schema.NoopContext(ctx, d, meta)
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

package dns

import (
	"context"
	"errors"
	"fmt"
	"github.com/apex/log"
	"strings"
	"time"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
			"tsig_key": {
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

func resourceDNSv2ZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	meta := akamai.Meta(m)
	logger := meta.Log("[Akamai DNS]", "resourceDNSZoneCreate")

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	hostname, err := tools.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.WithField("zone", hostname).Info("Zone Create")
	zoneType, err := tools.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	masterSet, err := tools.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	masterlist := masterSet.List()
	if strings.ToUpper(zoneType) == "SECONDARY" && len(masterlist) == 0 {
		return diag.Errorf("DNS Secondary zone requires masters for zone %v", hostname)
	}
	contractStr, err := tools.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupStr, err := tools.GetStringValue("group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contract := strings.TrimPrefix(contractStr, "ctr_")
	group := strings.TrimPrefix(groupStr, "grp_")
	zoneQueryString := dnsv2.ZoneQueryString{Contract: contract, Group: group}
	zoneCreate := &dnsv2.ZoneCreate{Zone: hostname, Type: zoneType}
	if err := populateDNSv2ZoneObject(d, zoneCreate, logger); err != nil {
		return diag.FromErr(err)
	}
	// First try to get the zone from the API
	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, e := dnsv2.GetZone(hostname)

	if e != nil {
		if !dnsv2.IsConfigDNSError(e) || !e.(dnsv2.ConfigDNSError).NotFound() {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "API falure failure",
				Detail:   e.Error(),
			})
		}
		// If there's no existing zone we'll create a blank one
		// if the zone is not found/404 we will create a new
		// blank zone for the records to be added to and continue
		logger.Debug(fmt.Sprintf("Creating new zone: %v", zoneCreate))
		e = zoneCreate.Save(zoneQueryString, true)
		if e != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone create failure",
				Detail:   e.Error(),
			})
		}
		if strings.ToUpper(zoneType) == "PRIMARY" {
			time.Sleep(2 * time.Second)
			// Indirectly create NS and SOA records
			e = zoneCreate.SaveChangelist()
			if e != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Zone create failure",
					Detail:   e.Error(),
				})
			}
			time.Sleep(time.Second)
			e = zoneCreate.SubmitChangelist()
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
	}

	// Save the zone to the API
	logger.Debug(fmt.Sprintf("Zone exists. Updating zone %v", zoneCreate))
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(ctx, d, meta)

}

func resourceDNSv2ZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)
	meta := akamai.Meta(m)
	logger := meta.Log("[Akamai DNS]", "resourceDNSZoneRead")

	hostname, err := tools.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.WithField("zone", hostname).Info("Zone Read")
	masterSet, err := tools.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	masterlist := masterSet.List()
	masters := make([]string, 0, len(masterlist))
	if len(masterlist) > 0 {
		for _, master := range masterlist {
			masterStr, ok := master.(string)
			if !ok {
				return diag.Errorf("'master' is of invalid type; should be 'string'")
			}
			masters = append(masters, masterStr)
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
	zoneType, err := tools.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if strings.ToUpper(zone.Type) != strings.ToUpper(zoneType) {
		return diag.Errorf("zone type has changed from %s to %s", zoneType, zone.Type)
	}
	if err := populateDNSv2ZoneState(d, zone); err != nil {
		return diag.FromErr(err)
	}

	logger.Debug(fmt.Sprintf("READ content: %v", zone))
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionId, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionId, zone.Zone, hostname))
	}
	return nil
}

// Update DNS Zone
func resourceDNSv2ZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)
	meta := akamai.Meta(m)
	logger := meta.Log("[Akamai DNS]", "resourceDNSZoneUpdate")

	logger.WithField("zone", hostname).Info("Zone Update")

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	hostname, err := tools.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}
	group, err := tools.GetStringValue("group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneType, err := tools.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneQueryString := dnsv2.ZoneQueryString{Contract: contract, Group: group}

	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, e := dnsv2.GetZone(hostname)
	if e != nil {
		// If there's no existing zone we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			logger.Debug(e.Error())
			// Something drastically wrong if we are trying to update a non existent zone!
			return diag.Errorf("attempt to update non existent zone: %s", hostname)
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone read failure",
			Detail:   e.Error(),
		})
	}
	// Create Zone Post obj and copy Received vals over
	zoneCreate := &dnsv2.ZoneCreate{Zone: hostname, Type: zoneType}
	zoneCreate.Masters = zone.Masters
	zoneCreate.Comment = zone.Comment
	zoneCreate.SignAndServe = zone.SignAndServe
	zoneCreate.SignAndServeAlgorithm = zone.SignAndServeAlgorithm
	zoneCreate.Target = zone.Target
	zoneCreate.EndCustomerId = zone.EndCustomerId
	zoneCreate.ContractId = zone.ContractId
	zoneCreate.TsigKey = zone.TsigKey
	if err := populateDNSv2ZoneObject(d, zoneCreate, logger); err != nil {
		return diag.FromErr(err)
	}
	// Save the zone to the API
	logger.Debug(fmt.Sprintf("Saving zone %v", zoneCreate))
	e = zoneCreate.Update(zoneQueryString)
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
func resourceDNSv2ZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()
	meta := akamai.Meta(m)
	logger := meta.Log("[Akamai DNS]", "resourceDNSZoneImport")

	logger.WithField("zone", hostname).Info("Zone Import")

	// find the zone first
	logger.Debug(fmt.Sprintf("Searching for zone [%s]", hostname))
	zone, err := dnsv2.GetZone(hostname)
	if err != nil {
		return nil, err
	}

	if err := d.Set("zone", zone.Zone); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("type", zone.Type); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := populateDNSv2ZoneState(d, zone); err != nil {
		return nil, err
	}

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s:%s:%s", zone.VersionId, zone.Zone, hostname))

	return []*schema.ResourceData{d}, nil
}

func resourceDNSv2ZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	hostname, err := tools.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	meta := akamai.Meta(m)
	logger := meta.Log("[Akamai DNS]", "resourceDNSZoneDelete")
	logger.WithField("zone", hostname).Info("Zone Import")
	logger.Warn("DNS Zone deletion not allowed")

	// No ZONE delete operation permitted.
	return schema.NoopContext(ctx, d, meta)
}

// validateZoneType is a SchemaValidateFunc to validate the Zone type.
func validateZoneType(v interface{}, _ string) (ws []string, es []error) {
	value := strings.ToUpper(v.(string))
	if value != "PRIMARY" && value != "SECONDARY" && value != "ALIAS" {
		es = append(es, fmt.Errorf("Type must be PRIMARY, SECONDARY, or ALIAS"))
	}
	return
}

// populate zone state based on API response.
func populateDNSv2ZoneState(d *schema.ResourceData, zoneresp *dnsv2.ZoneResponse) error {

	if err := d.Set("masters", zoneresp.Masters); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("comment", zoneresp.Comment); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("sign_and_serve", zoneresp.SignAndServe); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("sign_and_serve_algorithm", zoneresp.SignAndServeAlgorithm); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("target", zoneresp.Target); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("end_customer_id", zoneresp.EndCustomerId); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	tsigListNew := make([]interface{}, 0)
	if zoneresp.TsigKey != nil {
		tsigNew := map[string]interface{}{
			"name":      zoneresp.TsigKey.Name,
			"algorithm": zoneresp.TsigKey.Algorithm,
			"secret":    zoneresp.TsigKey.Secret,
		}
		tsigListNew = append(tsigListNew, tsigNew)
	}
	if err := d.Set("tsig_key", tsigListNew); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("activation_state", zoneresp.ActivationState); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("alias_count", zoneresp.AliasCount); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("version_id", zoneresp.VersionId); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

// populate zone object based on current config.
func populateDNSv2ZoneObject(d *schema.ResourceData, zone *dnsv2.ZoneCreate, logger log.Interface) error {
	masterSet, err := tools.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	masterlist := masterSet.List()
	masters := make([]string, 0, len(masterlist))
	for _, master := range masterlist {
		masterStr, ok := master.(string)
		if !ok {
			return fmt.Errorf("'master' is of invalid type; should be 'string'")
		}
		masters = append(masters, masterStr)
	}
	zone.Masters = masters
	comment, err := tools.GetStringValue("comment", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("comment") {
		zone.Comment = comment
	}
	signAndServe, err := tools.GetBoolValue("sign_and_serve", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	zone.SignAndServe = signAndServe
	signAndServeAlgorithm, err := tools.GetStringValue("sign_and_serve_algorithm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("sign_and_serve_algorithm") {
		zone.SignAndServeAlgorithm = signAndServeAlgorithm
	}
	target, err := tools.GetStringValue("target", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("target") {
		zone.Target = target
	}
	endCustomerID, err := tools.GetStringValue("end_customer_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("end_customer_id") {
		zone.EndCustomerId = endCustomerID
	}
	tsigKey, err := tools.GetListValue("tsig_key", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		if !errors.Is(err, tools.ErrNotFound) {
			return err
		}
		zone.TsigKey = nil
		return nil
	}
	if len(tsigKey) == 0 {
		return nil
	}
	tsigKeyMap, ok := tsigKey[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("'tsig_key' entry is of invalid type; should be 'map[string]interface{}'")
	}
	zone.TsigKey = &dnsv2.TSIGKey{
		Name:      tsigKeyMap["name"].(string),
		Algorithm: tsigKeyMap["algorithm"].(string),
		Secret:    tsigKeyMap["secret"].(string),
	}
	logger.Debugf("Generated TsigKey [%v]", zone.TsigKey)
	return nil
}

// utility method to verify zone config fields based on type. not worrying about required fields ....
func checkDNSv2Zone(d tools.ResourceDataFetcher) error {
	zone, err := tools.GetStringValue("zone", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	zoneType, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	mastersSet, err := tools.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	target, err := tools.GetStringValue("target", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	tsig, err := tools.GetListValue("tsig_key", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	signandserve, err := tools.GetBoolValue("sign_and_serve", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	ztype := strings.ToUpper(zoneType)
	masters := mastersSet.List()
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

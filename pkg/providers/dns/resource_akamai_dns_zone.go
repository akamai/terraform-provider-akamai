package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
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
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateZoneType,
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tf.FieldPrefixSuppress("grp_"),
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
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSZoneCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	hostname, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.WithField("zone", hostname).Info("Zone Create")
	zoneType, err := tf.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	masterSet, err := tf.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	masterlist := masterSet.List()
	if strings.ToUpper(zoneType) == "SECONDARY" && len(masterlist) == 0 {
		return diag.Errorf("DNS Secondary zone requires masters for zone %v", hostname)
	}
	contractStr, err := tf.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := tf.GetStringValue("group", d)
	if err != nil {
		if errors.Is(err, tf.ErrNotFound) {
			groupList, err := inst.Client(meta).ListGroups(ctx, dns.ListGroupRequest{})
			if err != nil {
				return diag.FromErr(err)
			}
			if len(groupList.Groups) == 0 {
				return diag.Errorf("no group found. Please provide the group.")
			}
			if len(groupList.Groups) == 1 {
				group = strconv.Itoa(groupList.Groups[0].GroupID)
				logger.Warnf("Please modify configuration and provide group identifier. It will be required in the future version of the resource.")
			}
			if len(groupList.Groups) > 1 {
				return diag.Errorf("group is a required field when there is more than one group present.")
			}
		} else {
			return diag.FromErr(err)
		}
	}

	contract := strings.TrimPrefix(contractStr, "ctr_")
	group = strings.TrimPrefix(group, "grp_")
	zoneQueryString := dns.ZoneQueryString{Contract: contract, Group: group}
	zoneCreate := &dns.ZoneCreate{Zone: hostname, Type: zoneType}
	if err := populateDNSv2ZoneObject(d, zoneCreate, logger); err != nil {
		return diag.FromErr(err)
	}
	// First try to get the zone from the API
	logger.Debugf("Searching for zone [%s]", hostname)
	zone, e := inst.Client(meta).GetZone(ctx, hostname)

	if e == nil {
		// Not a good idea to overwrite an existing zone. Needs to be imported.
		logger.Errorf("Zone creation error. Zone %s exists", hostname)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone exists. Please import.",
			Detail:   fmt.Sprintf("Zone create failure. Zone %s exists", hostname),
		})
	}
	apiError, ok := e.(*dns.Error)
	if !ok || apiError.StatusCode != http.StatusNotFound {
		logger.Errorf("Create[ERROR] %w", e)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create API falure",
			Detail:   e.Error(),
		})
	}

	// no existing zone.
	logger.Debugf("Creating new zone: %v", zoneCreate)
	e = inst.Client(meta).CreateZone(ctx, zoneCreate, zoneQueryString, true)
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
		e = inst.Client(meta).SaveChangelist(ctx, zoneCreate)
		if e != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone create failure",
				Detail:   e.Error(),
			})
		}
		time.Sleep(time.Second)
		e = inst.Client(meta).SubmitChangelist(ctx, zoneCreate)
		if e != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone create failure",
				Detail:   e.Error(),
			})
		}
	}
	zone, e = inst.Client(meta).GetZone(ctx, hostname)
	if e != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone read after create failure",
			Detail:   e.Error(),
		})
	}
	d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionID, zone.Zone, hostname))
	return resourceDNSv2ZoneRead(ctx, d, meta)

}

func resourceDNSv2ZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSZoneRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	hostname, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.WithField("zone", hostname).Info("Zone Read")
	masterSet, err := tf.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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
	logger.Debugf("Searching for zone [%s]", hostname)
	zone, e := inst.Client(meta).GetZone(ctx, hostname)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if ok && apiError.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diag.FromErr(e)
		}
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone read failure",
			Detail:   e.Error(),
		})
	}
	// Populate state with returned field values ... except zone and type
	zoneType, err := tf.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if strings.ToUpper(zone.Type) != strings.ToUpper(zoneType) {
		return diag.Errorf("zone type has changed from %s to %s", zoneType, zone.Type)
	}
	if strings.ToUpper(zone.Type) == "PRIMARY" {
		// TFP-196 - check if SOA and NS exist. If not, create
		err = checkZoneSOAandNSRecords(ctx, meta, zone, logger)
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone is in an indeterminate state",
				Detail:   err.Error(),
			})
		}
		// Need updated state
		zone, err = inst.Client(meta).GetZone(ctx, hostname)
		if err != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Zone read failure",
				Detail:   err.Error(),
			})
		}
	}
	if err := populateDNSv2ZoneState(d, zone); err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("READ content: %v", zone)
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionID, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionID, zone.Zone, hostname))
	}
	return nil
}

// Update DNS Zone
func resourceDNSv2ZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	hostname := d.Get("zone").(string)
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSZoneUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.WithField("zone", hostname).Info("Zone Update")

	if err := checkDNSv2Zone(d); err != nil {
		return diag.FromErr(err)
	}
	hostname, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contract, err := tf.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}
	group, err := tf.GetStringValue("group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneType, err := tf.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	zoneQueryString := dns.ZoneQueryString{Contract: contract, Group: group}

	logger.Debugf("Searching for zone [%s]", hostname)
	zone, e := inst.Client(meta).GetZone(ctx, hostname)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if !ok && apiError.StatusCode != http.StatusOK {
			logger.Debugf("Zone Update read faiiled: %s", e.Error())
			return diag.FromErr(fmt.Errorf("Update zone %s read failed: %w", hostname, e))
		}
	}
	// Create Zone Post obj and copy Received vals over
	zoneCreate := &dns.ZoneCreate{Zone: hostname, Type: zoneType}
	zoneCreate.Masters = zone.Masters
	zoneCreate.Comment = zone.Comment
	zoneCreate.SignAndServe = zone.SignAndServe
	zoneCreate.SignAndServeAlgorithm = zone.SignAndServeAlgorithm
	zoneCreate.Target = zone.Target
	zoneCreate.EndCustomerID = zone.EndCustomerID
	zoneCreate.ContractID = zone.ContractID
	zoneCreate.TSIGKey = zone.TSIGKey
	if err := populateDNSv2ZoneObject(d, zoneCreate, logger); err != nil {
		return diag.FromErr(err)
	}
	// Save the zone to the API
	logger.Debugf("Saving zone %v", zoneCreate)
	e = inst.Client(meta).UpdateZone(ctx, zoneCreate, zoneQueryString)
	if e != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Zone update failure",
			Detail:   e.Error(),
		})
	}

	// Give terraform the ID
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone.VersionID, zone.Zone, hostname))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s", zone.VersionID, zone.Zone, hostname))
	}
	return resourceDNSv2ZoneRead(ctx, d, meta)
}

// Import Zone. Id is the zone
func resourceDNSv2ZoneImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSZoneImport")
	// create a context with logging for api calls
	ctx := context.TODO()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.WithField("zone", hostname).Info("Zone Import")

	// find the zone first
	logger.Debugf("Searching for zone [%s]", hostname)
	zone, err := inst.Client(meta).GetZone(ctx, hostname)
	if err != nil {
		return nil, err
	}

	if strings.ToUpper(zone.Type) == "PRIMARY" {
		// TFP-196 - check if SOA and NS exist. If not, create
		err = checkZoneSOAandNSRecords(ctx, meta, zone, logger)
		if err != nil {
			return nil, err
		}
		// Need updated state
		zone, err = inst.Client(meta).GetZone(ctx, hostname)
		if err != nil {
			return nil, err
		}
	}

	if err := d.Set("zone", zone.Zone); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("type", zone.Type); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := populateDNSv2ZoneState(d, zone); err != nil {
		return nil, err
	}

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s:%s:%s", zone.VersionID, zone.Zone, hostname))

	return []*schema.ResourceData{d}, nil
}

func resourceDNSv2ZoneDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hostname, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSZoneDelete")
	logger.WithField("zone", hostname).Info("Zone Delete")
	// Ignore for Unit test Lifecycle
	if _, ok := os.LookupEnv("DNS_ZONE_SKIP_DELETE"); ok {
		logger.Info("DNS Zone delete: intentially skipping")
		return nil
	}
	logger.Warn("DNS Zone deletion not allowed")

	// No ZONE delete operation permitted.
	return diag.Errorf("DNS zone deletion is not supported via this sub provider")
}

// validateZoneType is a SchemaValidateDiagFunc to validate the Zone type.
func validateZoneType(v interface{}, _ cty.Path) diag.Diagnostics {
	value := strings.ToUpper(v.(string))
	if value != "PRIMARY" && value != "SECONDARY" && value != "ALIAS" {
		return diag.Errorf("Type must be PRIMARY, SECONDARY, or ALIAS")
	}
	return nil
}

// populate zone state based on API response.
func populateDNSv2ZoneState(d *schema.ResourceData, zoneresp *dns.ZoneResponse) error {

	if err := d.Set("contract", zoneresp.ContractID); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("masters", zoneresp.Masters); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("comment", zoneresp.Comment); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("sign_and_serve", zoneresp.SignAndServe); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("sign_and_serve_algorithm", zoneresp.SignAndServeAlgorithm); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("target", zoneresp.Target); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("end_customer_id", zoneresp.EndCustomerID); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	tsigListNew := make([]interface{}, 0)
	if zoneresp.TSIGKey != nil {
		tsigNew := map[string]interface{}{
			"name":      zoneresp.TSIGKey.Name,
			"algorithm": zoneresp.TSIGKey.Algorithm,
			"secret":    zoneresp.TSIGKey.Secret,
		}
		tsigListNew = append(tsigListNew, tsigNew)
	}
	if err := d.Set("tsig_key", tsigListNew); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("activation_state", zoneresp.ActivationState); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("alias_count", zoneresp.AliasCount); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("version_id", zoneresp.VersionID); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

// populate zone object based on current config.
func populateDNSv2ZoneObject(d *schema.ResourceData, zone *dns.ZoneCreate, logger log.Interface) error {
	masterSet, err := tf.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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
	comment, err := tf.GetStringValue("comment", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("comment") {
		zone.Comment = comment
	}
	signAndServe, err := tf.GetBoolValue("sign_and_serve", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	zone.SignAndServe = signAndServe
	signAndServeAlgorithm, err := tf.GetStringValue("sign_and_serve_algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("sign_and_serve_algorithm") {
		zone.SignAndServeAlgorithm = signAndServeAlgorithm
	}
	target, err := tf.GetStringValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("target") {
		zone.Target = target
	}
	endCustomerID, err := tf.GetStringValue("end_customer_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if err == nil || d.HasChange("end_customer_id") {
		zone.EndCustomerID = endCustomerID
	}
	TSIGKey, err := tf.GetListValue("tsig_key", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		if !errors.Is(err, tf.ErrNotFound) {
			return err
		}
		zone.TSIGKey = nil
		return nil
	}
	if len(TSIGKey) == 0 {
		return nil
	}
	TSIGKeyMap, ok := TSIGKey[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("'tsig_key' entry is of invalid type; should be 'map[string]interface{}'")
	}
	zone.TSIGKey = &dns.TSIGKey{
		Name:      TSIGKeyMap["name"].(string),
		Algorithm: TSIGKeyMap["algorithm"].(string),
		Secret:    TSIGKeyMap["secret"].(string),
	}
	logger.Debugf("Generated TSIGKey [%v]", zone.TSIGKey)
	return nil
}

// utility method to verify zone config fields based on type. not worrying about required fields ....
func checkDNSv2Zone(d tf.ResourceDataFetcher) error {
	zone, err := tf.GetStringValue("zone", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	zoneType, err := tf.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	mastersSet, err := tf.GetSetValue("masters", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	target, err := tf.GetStringValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	tsig, err := tf.GetListValue("tsig_key", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	signandserve, err := tf.GetBoolValue("sign_and_serve", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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

// Util func to create SOA and NS records
func checkZoneSOAandNSRecords(ctx context.Context, meta meta.Meta, zone *dns.ZoneResponse, logger log.Interface) error {
	logger.Debugf("Checking SOA and NS records exist for zone %s", zone.Zone)
	var resp *dns.RecordSetResponse
	var err error
	if zone.ActivationState != "NEW" {
		// See if SOA and NS recs exist already. Both or none.
		resp, err = inst.Client(meta).GetRecordSets(ctx, zone.Zone, dns.RecordSetQueryArgs{Types: "SOA,NS"})
		if err != nil {
			return err
		}
	}
	if resp != nil && len(resp.RecordSets) >= 2 {
		return nil
	}

	logger.Warnf("SOA and NS records don't exist. Creating ...")
	nameservers, err := inst.Client(meta).GetNameServerRecordList(ctx, zone.ContractID)
	if err != nil {
		return err
	}
	if len(nameservers) < 1 {
		return fmt.Errorf("No authoritative nameservers exist for zone %s contract ID", zone.Zone)
	}
	rs := &dns.RecordSets{RecordSets: make([]dns.RecordSet, 0)}
	rs.RecordSets = append(rs.RecordSets, createSOARecord(zone.Zone, nameservers, logger))
	rs.RecordSets = append(rs.RecordSets, createNSRecord(zone.Zone, nameservers, logger))

	// create recordSets
	err = inst.Client(meta).CreateRecordSets(ctx, rs, zone.Zone, true)

	return err
}

func createSOARecord(zone string, nameservers []string, _ log.Interface) dns.RecordSet {
	rec := dns.RecordSet{Name: zone, Type: "SOA"}
	rec.TTL = 86400
	peMail := fmt.Sprintf("hostmaster.%s.", zone)
	soaData := fmt.Sprintf("%s %s 1 14400 7200 604800 1200", nameservers[0], peMail)
	rec.Rdata = []string{soaData}

	return rec
}

func createNSRecord(zone string, nameservers []string, _ log.Interface) dns.RecordSet {
	rec := dns.RecordSet{Name: zone, Type: "NS"}
	rec.TTL = 86400
	rec.Rdata = nameservers

	return rec
}

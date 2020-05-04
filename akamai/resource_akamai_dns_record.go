package akamai

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Retry count for save, update and delete
const opRetryCount = 3

func resourceDNSv2Record() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Exists: resourceDNSRecordExists,
		Delete: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDNSRecordImport,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recordtype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					RRTypeA,
					RRTypeAaaa,
					RRTypeCname,
					RRTypeLoc,
					RRTypeNs,
					RRTypePtr,
					RRTypeSpf,
					RRTypeTxt,
					RRTypeAfsdb,
					RRTypeDnskey,
					RRTypeDs,
					RRTypeHinfo,
					RRTypeMx,
					RRTypeNaptr,
					RRTypeNsec3,
					RRTypeNsec3Param,
					RRTypeRp,
					RRTypeRrsig,
					RRTypeSrv,
					RRTypeSshfp,
					RRTypeSoa,
					RRTypeAkamaiCdn,
					RRTypeAkamaiTlc,
				}, false),
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"target": {
				Type:             schema.TypeSet,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Optional:         true,
				DiffSuppressFunc: dnsRecordTargetSuppress,
				Set:              schema.HashString,
			},
			"subtype": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flags": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"algorithm": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keytag": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"digest_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"digest": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hardware": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldold := strings.Trim(old, "\\\"")
					newnew := strings.Trim(new, "\"")
					if oldold == newnew {
						return true
					}
					return false
				},
			},
			"software": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldold := strings.Trim(old, "\\\"")
					newnew := strings.Trim(new, "\"")
					if oldold == newnew {
						return true
					}
					return false
				},
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flagsnaptr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regexp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replacement": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iterations": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"salt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_hashed_owner_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_bitmaps": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mailbox": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"txt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type_covered": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"original_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inception": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fingerprint_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority_increment": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"answer_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"serial": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"refresh": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"retry": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"expiry": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nxdomain_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"record_sha": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// DiffSuppresFunc to handle quoted TXT Rdata strings possibly containining escaped quotes
func dnsRecordTargetSuppress(k, old, new string, d *schema.ResourceData) bool {

	// Function invocation behavior is not obvious for Sets. Called for each entry in a set
	// May call with values in both old and new or in one or the other.
	// Seems if different, get one invocation with old val and null new as well as
	// invocation with new val and null old. In all cases, we retrieve old and new sets
	// from ResourceData and validate thru those.
	// Function only validates TXT records.

	oldset, newset := d.GetChange("target")
	oldTargetList := oldset.(*schema.Set).List()
	newTargetList := newset.(*schema.Set).List()
	if len(oldTargetList) != len(newTargetList) {
		return false
	}
	//var baseList []interface{}
	var compList []interface{}
	var baseVal = ""
	var compTrim = ""
	// one or both must have value
	if old != "" && new != "" {
		baseVal = old
		compTrim = "\""
		baseVal = strings.Trim(baseVal, "\\\"")
		if strings.Contains(baseVal, "\\\"") {
			baseVal = strings.ReplaceAll(baseVal, "\\\"", "\"")
		}
		//baseList = oldTargetList
		compList = newTargetList
	} else if old == "" {
		baseVal = new
		compTrim = "\\\""
		baseVal = strings.Trim(baseVal, "\"")
		//baseList = newTargetList
		compList = oldTargetList
	} else if new == "" {
		baseVal = old
		compTrim = "\""
		baseVal = strings.Trim(baseVal, "\\\"")
		if strings.Contains(baseVal, "\\\"") {
			baseVal = strings.ReplaceAll(baseVal, "\\\"", "\"")
		}
		//baseList = oldTargetList
		compList = newTargetList
	}
	for _, compval := range compList {
		if compTrim == "\\\"" && strings.Contains(compval.(string), "\\\"") {
			compval = strings.ReplaceAll(compval.(string), "\\\"", "\"")
		}
		if baseVal == strings.Trim(compval.(string), "\"") {
			return true
		}
	}

	return false

}

// Lock per record type
var recordCreateLock = map[string]*sync.Mutex{
	"A":          &sync.Mutex{},
	"AAAA":       &sync.Mutex{},
	"AFSDB":      &sync.Mutex{},
	"AKAMAICDN":  &sync.Mutex{},
	"AKAMAITLC":  &sync.Mutex{},
	"CAA":        &sync.Mutex{},
	"CNAME":      &sync.Mutex{},
	"HINFO":      &sync.Mutex{},
	"LOC":        &sync.Mutex{},
	"MX":         &sync.Mutex{},
	"NAPTR":      &sync.Mutex{},
	"NS":         &sync.Mutex{},
	"PTR":        &sync.Mutex{},
	"RP":         &sync.Mutex{},
	"SOA":        &sync.Mutex{},
	"SRV":        &sync.Mutex{},
	"SPF":        &sync.Mutex{},
	"SSHFP":      &sync.Mutex{},
	"TLSA":       &sync.Mutex{},
	"TXT":        &sync.Mutex{},
	"DNSKEY":     &sync.Mutex{},
	"DS":         &sync.Mutex{},
	"NSEC3":      &sync.Mutex{},
	"NSEC3PARAM": &sync.Mutex{},
	"RRSIG":      &sync.Mutex{}}

/*
// Following  supports Lock by Recordset. Not clear dnsv2 API would be able to handle for very large zones

var zoneRecordCreateLock = make(sync.Map[string]*sync.Mutex)


var recordTypeMapLock = &sync.Map{
        "A":            &sync.Map{},
        "AAAA":         &sync.Map{},
        "AFSDB":        &sync.Map{},
        "AKAMAICDN":    &sync.Map{},
        "AKAMAITLC":    &sync.Map{},
        "CAA":          &sync.Map{},
        "CNAME":        &sync.Map{},
        "HINFO":        &sync.Map{},
        "LOC":          &sync.Map{},
        "MX":           &sync.Map{},
        "NAPTR":        &sync.Map{},
        "NS":           &sync.Map{},
        "PTR":          &sync.Map{},
        "RP":           &sync.Map{},
        "SOA":          &sync.Map{},
        "SRV":          &sync.Map{},
        "SPF":          &sync.Map{},
        "SSHFP":        &sync.Map{},
        "TLSA":         &sync.Map{},
        "TXT":          &sync.Map{},
        "DNSKEY":       &sync.Map{},
        "DS":           &sync.Map{},
        "NSEC3":        &sync.Map{},
        "NSEC3PARAM":   &sync.Map{},
        "RRSIG":        &sync.Map{}}

// Retrieves record lock per recordset
func getRecordLock(zone, host, recordtype string) *sync.Mutex {

	typeMap, _ := recordTypeMapLock.LoadOrStore(recordtype, &sync.Map{})
        lockindex := zone + "_" + host + "_" + recordtype
        recordLock, _ := typeMap.(*sync.Map).LoadOrStore(lockindex, &sync.Mutex{})

        return recordLock.(*sync.Mutex)
}
*/

// Retrieves record lock per record type
func getRecordLock(zone, host, recordtype string) *sync.Mutex {

	return recordCreateLock[recordtype]

}

// Record function signature
type recordFunction func(string, ...bool) error

func executeRecordFunction(name string, d *schema.ResourceData, fn recordFunction, zone string, host string, recordtype string, rlock bool) error {

	// DNS API can have Concurrency issues
	opRetry := opRetryCount
	e := fn(zone, false)
	for e != nil && opRetry > 0 {
		if dnsv2.IsConfigDNSError(e) {
			if e.(*dnsv2.RecordError).ConcurrencyConflict() {
				opRetry -= 1
				e = fn(zone, false)
				continue
			} else if name == "DELETE" && e.(dnsv2.ConfigDNSError).NotFound() == true {
				// record doesn't exist
				d.SetId("")
				log.Printf("[DEBUG] [Akamai DNSv2] %s [WARNING] %s", name, "Record not found")
				break
			} else {
				log.Printf("[DEBUG] [Akamai DNSv2] %s [ERROR] %s", name, e.Error())
				return e
			}
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] %s Record failed for record [%s] [%s] [%s] ", name, zone, host, recordtype)
			return e
		}
		break
	}

	return nil

}

// Create a new DNS Record
func resourceDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record per record type to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	var zone string
	var host string
	var recordtype string

	_, ok := d.GetOk("zone")
	if ok {
		zone = d.Get("zone").(string)
	}
	_, ok = d.GetOk("name")
	if ok {
		host = d.Get("name").(string)
	}
	_, ok = d.GetOk("recordtype")
	if ok {
		recordtype = d.Get("recordtype").(string)
	}

	err := validateRecord(d)
	if err != nil {
		return fmt.Errorf("DNS record validation failure on zone %v: %v", zone, err)
	}

	// serialize record creates of same type
	getRecordLock(zone, host, recordtype).Lock()
	defer getRecordLock(zone, host, recordtype).Unlock()

	recordcreate, err := bindRecord(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Record Create Bind object  [%v]", recordcreate)

	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)

	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for recordcreate [%s]", sha1hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil && !dnsv2.IsConfigDNSError(e) {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}
	// If there's no existing record we'll create a blank one
	if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
		// if the record is not found/404 we will create a new
		log.Printf("[DEBUG] [Akamai DNSv2] [ERROR] %s", e.Error())
		log.Printf("[DEBUG] [Akamai DNSv2] Creating new record")
		// Save the zone to the API
		e = executeRecordFunction("CREATE", d, recordcreate.Save, zone, host, recordtype, false)
		if e != nil {
			return e
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] Updating record")
		if len(rdata) > 0 {
			e = executeRecordFunction("CREATE", d, recordcreate.Update, zone, host, recordtype, false)
			if e != nil {
				return e
			}
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] Saving record")
			e = executeRecordFunction("CREATE", d, recordcreate.Save, zone, host, recordtype, false)
			if e != nil {
				return e
			}
		}
	}
	// save hash
	d.Set("record_sha", sha1hash)
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordtype))
	} else {
		// Backwards compatibility
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
	}
	// Lock won't be release til after Read ...
	return resourceDNSRecordRead(d, meta)

}

// Update DNS Record
func resourceDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record per record type to be updated at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	var zone string
	var host string
	var recordtype string

	_, ok := d.GetOk("zone")
	if ok {
		zone = d.Get("zone").(string)
	}
	_, ok = d.GetOk("name")
	if ok {
		host = d.Get("name").(string)
	}
	_, ok = d.GetOk("recordtype")
	if ok {
		recordtype = d.Get("recordtype").(string)
	}
	_, ok = d.GetOk("target")
	if ok {
		target := d.Get("target").(*schema.Set).List()
		records := make([]string, 0, len(target))
		for _, recContent := range target {
			records = append(records, recContent.(string))
		}
		log.Printf("[DEBUG] [Akamai DNSv2] Update Records [%v]", records)
	}

	err := validateRecord(d)
	if err != nil {
		return fmt.Errorf("DNS record validation failure on zone %v: %v", zone, err)
	}

	// serialize record updates of same type
	getRecordLock(zone, host, recordtype).Lock()
	defer getRecordLock(zone, host, recordtype).Unlock()

	if recordtype == "SOA" {
		// need to get current serial and increment as part of update
		record, e := dnsv2.GetRecord(zone, host, recordtype)
		if e != nil {
			if dnsv2.IsConfigDNSError(e) {
				if !e.(dnsv2.ConfigDNSError).NotFound() {
					log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Read [ERROR] %s", e.Error())
					return e
				}
			} else {
				log.Printf("[ERROR] [Akamai DNSv2] UPDATE Record Read. error looking up "+recordtype+" records for %q: %s", host, e.Error())
				return e
			}
		} else {
			// Parse Rdata
			d.Set("serial", dnsv2.ParseRData(recordtype, record.Target)["serial"].(int)+1)
		}
	}

	recordcreate, err := bindRecord(d)
	if err != nil {
		return err
	}
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)

	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE SHA sum for recordupdate [%s]", sha1hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil && !dnsv2.IsConfigDNSError(e) {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Searching for records LEN %d", len(rdata))
	if len(rdata) > 0 {
		sort.Strings(rdata)
		extractString := strings.Join(rdata, " ")
		sha1hashtest := getSHAString(extractString)
		log.Printf("[DEBUG] [Akamai DNSv2] UPDATE SHA sum from recordread [%s]", sha1hashtest)
		// If there's no existing record we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the record is not found/404 we will create a new
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE [ERROR] %s", e.Error())
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Creating new record")
			// Save the zone to the API
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Updating record")
			e = executeRecordFunction("UPDATE", d, recordcreate.Save, zone, host, recordtype, false)
			if e != nil {
				return e
			}
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Updating record")
			e = executeRecordFunction("UPDATE", d, recordcreate.Update, zone, host, recordtype, false)
			if e != nil {
				return e
			}

		}
		// save hash
		d.Set("record_sha", sha1hash)
		// Give terraform the ID
		if d.Id() == "" || strings.Contains(d.Id(), "#") {
			d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordtype))
		} else {
			d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
		}
	}
	// Lock not released until after Read ...
	return resourceDNSRecordRead(d, meta)
}

func resourceDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	var zone string
	var host string
	var recordtype string

	_, ok := d.GetOk("zone")
	if ok {
		zone = d.Get("zone").(string)
	}
	_, ok = d.GetOk("name")
	if ok {
		host = d.Get("name").(string)
	}
	_, ok = d.GetOk("recordtype")
	if ok {
		recordtype = d.Get("recordtype").(string)
	}

	recordcreate, err := bindRecord(d)
	if err != nil {
		return err
	}
	b, err := json.Marshal(recordcreate.Target)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] READ record JSON from bind records %s %s %s %s", string(b), zone, host, recordtype)
	//sort.Strings(recordcreate.Target)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)
	log.Printf("[DEBUG] [Akamai DNSv2] READ SHA sum for Existing SHA test %s %s", extractString, sha1hash)

	// try to get the zone from the API
	log.Printf("[INFO] [Akamai DNSv2] READ Searching for zone records %s %s %s", zone, host, recordtype)
	record, e := dnsv2.GetRecord(zone, host, recordtype)
	if e != nil && !dnsv2.IsConfigDNSError(e) {
		log.Printf("[ERROR] [Akamai DNSv2] RECORD READ. error looking up "+recordtype+" records for %q: %s", host, e.Error())
		return e
	}
	if dnsv2.IsConfigDNSError(e) {
		if e.(dnsv2.ConfigDNSError).NotFound() == true {
			// record doesn't exist
			log.Printf("[DEBUG] [Akamai DNSv2] READ Record Not Found [ERROR] %s", e.Error())
			d.SetId("")
			return fmt.Errorf("Record not found")
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] READ [ERROR] %s", e.Error())
			return e
		}
	}
	log.Printf("[DEBUG] [Akamai DNSv2] RECORD READ [%v] [%s] [%s] [%s] ", record, zone, host, recordtype)
	b1, err := json.Marshal(record.Target)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("[DEBUG] [Akamai DNSv2] READ record data read JSON %s", string(b1))
	rdataFieldMap := dnsv2.ParseRData(recordtype, record.Target) // returns map[string]interface{}
	if recordtype != "MX" {
		// Parse Rdata. MX special
		for fname, fvalue := range rdataFieldMap {
			d.Set(fname, fvalue)
		}
	}

	targets := dnsv2.ProcessRdata(record.Target, recordtype)
	if len(targets) > 0 {
		sort.Strings(targets)
		if recordtype == "SOA" {
			log.Printf("[DEBUG] [Akamai DNSv2] READ SOA RECORD")
			if rdataFieldMap["serial"].(int) >= d.Get("serial").(int) {
				log.Printf("[DEBUG] [Akamai DNSv2] READ SOA RECORD CHANGE: SOA OK ")
				if _, ok := validateSOARecord(d); ok {
					extractSoaString := strings.Join(targets, " ")
					sha1hash = getSHAString(extractSoaString)
					log.Printf("[DEBUG] [Akamai DNSv2] READ SOA RECORD CHANGE: SOA OK ")
				}
			}
		}
		d.Set("record_sha", sha1hash)
		// Give terraform the ID
		if strings.Contains(d.Id(), "#") {
			d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordtype))
		} else {
			d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
		}
		return nil
	}
	return fmt.Errorf("[ERROR] [Akamai DNSv2] READ -  Invalid RData Returned for Recordset %s %s %s", zone, host, recordtype)
}

func validateSOARecord(d *schema.ResourceData) (int, bool) {

	oldserial, newser := d.GetChange("serial")
	newserial := newser.(int)
	if oldserial.(int) > newserial {
		return newserial, false
	}
	if d.HasChange("name_server") {
		return newserial, false
	}
	if d.HasChange("email_address") {
		return newserial, false
	}
	if d.HasChange("refresh") {
		return newserial, false
	}
	if d.HasChange("retry") {
		return newserial, false
	}
	if d.HasChange("expiry") {
		return newserial, false
	}
	if d.HasChange("nxdomain_ttl") {
		return newserial, false
	}

	return newserial, true

}

func resourceDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idParts := strings.Split(d.Id(), "#")
	if len(idParts) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("Invalid Id for Zone Import: %s", d.Id())
	}
	log.Printf("[INFO] [Akamai DNS] Importing Record %s", d.Id())
	zone := idParts[0]
	recordname := idParts[1]
	recordtype := idParts[2]

	// Get recordset
	log.Printf("[INFO] [Akamai DNS] Searching for zone Recordset [%v]", idParts)
	recordset, e := dnsv2.GetRecord(zone, recordname, recordtype)
	if e != nil {
		if dnsv2.IsConfigDNSError(e) {
			if e.(dnsv2.ConfigDNSError).NotFound() == true {
				// record doesn't exist
				d.SetId("")
				log.Printf("[DEBUG] [Akamai DNSv2] IMPORT [ERROR] %s", "Record not found")
				return nil, fmt.Errorf("Record not found")
			} else {
				d.SetId("")
				log.Printf("[DEBUG] [Akamai DNSv2] IMPORT [ERROR] %s", e.Error())
				return nil, e
			}
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] IMPORT Record read failed for record [%s] [%s] [%s] ", zone, recordname, recordtype)
			d.SetId("")
			return []*schema.ResourceData{d}, e
		}
	}

	d.Set("zone", zone)
	d.Set("name", recordset.Name)
	d.Set("recordtype", recordset.RecordType)
	d.Set("ttl", recordset.TTL)
	// Parse Rdata
	rdataFieldMap := dnsv2.ParseRData(recordset.RecordType, recordset.Target) // returns map[string]interface{}
	for fname, fvalue := range rdataFieldMap {
		d.Set(fname, fvalue)
	}
	targets := dnsv2.ProcessRdata(recordset.Target, recordtype)
	importTargetString := ""
	if len(targets) > 0 {
		if recordtype != "MX" {
			// MX Target Order important
			sort.Strings(targets)
		}
		importTargetString = strings.Join(targets, " ")
		sha1hash := getSHAString(importTargetString)
		d.Set("record_sha", sha1hash)
		d.SetId(fmt.Sprintf("%s#%s#%s", zone, recordname, recordtype))
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] IMPORT Invalid Record. No target returned  [%s] [%s] [%s] ", zone, recordname, recordtype)
		d.SetId("")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {

	zone := d.Get("zone").(string)
	host := d.Get("name").(string)
	recordtype := d.Get("recordtype").(string)
	ttl := d.Get("ttl").(int)

	// serialize record updates of same type
	getRecordLock(zone, host, recordtype).Lock()
	defer getRecordLock(zone, host, recordtype).Unlock()

	target := d.Get("target").(*schema.Set).List()

	records := make([]string, 0, len(target))
	for _, recContent := range target {
		records = append(records, recContent.(string))
	}
	sort.Strings(records)
	log.Printf("[INFO] [Akamai DNS] Delete zone Records %v", records)
	recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}

	// Warning: Delete will expunge the ENTIRE Recordset regardless of whether user thought they were removing an instance
	return executeRecordFunction("DELETE", d, recordcreate.Delete, zone, host, recordtype, false)
}

func resourceDNSRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	var zone string
	var host string
	var recordtype string

	_, ok := d.GetOk("zone")
	if ok {
		zone = d.Get("zone").(string)
	}
	_, ok = d.GetOk("name")
	if ok {
		host = d.Get("name").(string)
	}
	_, ok = d.GetOk("recordtype")
	if ok {
		recordtype = d.Get("recordtype").(string)
	}

	log.Printf("[INFO] [Akamai DNS] Record Exists Check: %s %s %s", zone, host, recordtype)

	// Get recordset
	recordset, e := dnsv2.GetRecord(zone, host, recordtype)
	if e != nil {
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() {
			d.SetId("")
			return false, nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] EXISTS Record read failed for record [%s] [%s] [%s] ", zone, host, recordtype)
			return false, e
		}
	}

	return recordset != nil, nil

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsPriority(m map[string]int, p int) bool {
	for _, e := range m {
		if p == e {
			return true
		}
	}
	return false
}

// Encode IPV6 as a full string
func FullIPv6(ip net.IP) string {

	dst := make([]byte, hex.EncodedLen(len(ip)))
	_ = hex.Encode(dst, ip)
	return string(dst[0:4]) + ":" +
		string(dst[4:8]) + ":" +
		string(dst[8:12]) + ":" +
		string(dst[12:16]) + ":" +
		string(dst[16:20]) + ":" +
		string(dst[20:24]) + ":" +
		string(dst[24:28]) + ":" +
		string(dst[28:])
}

func padvalue(str string) string {
	vstr := strings.Replace(str, "m", "", -1)
	log.Printf("[DEBUG] [Akamai DNSv2]  %s", vstr)
	vfloat, err := strconv.ParseFloat(vstr, 32)
	if err != nil {
		log.Printf("[DEBUG] [Akamai DNSv2] Error parse %s", vstr)
	}
	vresult := fmt.Sprintf("%.2f", vfloat)
	log.Printf("[DEBUG] [Akamai DNSv2] Padded v_result %s", vresult)
	return vresult
}

// Used to pad coordinates to x.xxm format
func padCoordinates(str string) string {

	s := strings.Split(str, " ")
	latD, latM, latS, latDir, longD, longM, longS, longDir, altitude, size, horizPrecision, vertPrecision := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]
	return latD + " " + latM + " " + latS + " " + latDir + " " + longD + " " + longM + " " + longS + " " + longDir + " " + padvalue(altitude) + "m " + padvalue(size) + "m " + padvalue(horizPrecision) + "m " + padvalue(vertPrecision) + "m"

}

func bindRecord(d *schema.ResourceData) (dnsv2.RecordBody, error) {

	var host string
	var recordtype string

	_, ok := d.GetOk("name")
	if ok {
		host = d.Get("name").(string)
	}
	_, ok = d.GetOk("recordtype")
	if ok {
		recordtype = d.Get("recordtype").(string)
	}

	ttl := d.Get("ttl").(int)
	target := d.Get("target").(*schema.Set).List()
	records := make([]string, 0, len(target))

	simplerecordtarget := map[string]bool{"AAAA": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "SRV": true, "TXT": true}

	for _, recContent := range target {
		if simplerecordtarget[recordtype] {

			if recordtype == "AAAA" {
				addr := net.ParseIP(recContent.(string))
				result := FullIPv6(addr)
				log.Printf("[DEBUG] [Akamai DNSv2] IPV6 full %s", result)
				records = append(records, result)
			} else if recordtype == "LOC" {
				log.Printf("[DEBUG] [Akamai DNSv2] LOC code format %s", recContent.(string))
				str := padCoordinates(recContent.(string))
				records = append(records, str)
			} else if recordtype == "SPF" {
				str := recContent.(string)
				if !strings.HasPrefix(str, "\"") {
					str = "\"" + str + "\""
				}
				records = append(records, str)
			} else if recordtype == "TXT" {
				str := recContent.(string)
				log.Printf("[DEBUG] [Akamai DNSv2] Bind TXT Data IN: [%s]", str)
				if strings.HasPrefix(str, "\"") {
					str = strings.TrimLeft(str, "\"")
				}
				if strings.HasSuffix(str, "\"") {
					str = strings.TrimRight(str, "\"")
				}
				if strings.Contains(str, "\\\\\\\"") {
					// look for and replace escaped embedded quotes
					str = strings.ReplaceAll(str, "\\\\\\\"", "\\\"")
				}
				str = "\"" + str + "\""

				log.Printf("[DEBUG] [Akamai DNSv2] Bind TXT Data %s", str)
				if strings.Contains(str, "\\\"") {
					//str = strings.ReplaceAll(str, "\\\"", "\"")
				}
				log.Printf("[DEBUG] [Akamai DNSv2] Bind TXT Data OUT: [%s]", str)
				records = append(records, str)
			} else {
				checktarget := recContent.(string)[len(recContent.(string))-1:]
				if checktarget == "." {
					records = append(records, recContent.(string))
				} else {
					records = append(records, recContent.(string)+".")
				}
			}
		} else {
			records = append(records, recContent.(string))
		}
	}

	emptyrecordcreate := dnsv2.RecordBody{}

	simplerecord := map[string]bool{"A": true, "AAAA": true, "AKAMAICDN": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "TXT": true}
	if simplerecord[recordtype] {
		sort.Strings(records)

		recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
		return recordcreate, nil
	} else {
		if recordtype == "AFSDB" {

			records := make([]string, 0, len(target))
			subtype := d.Get("subtype").(int)
			for _, recContent := range target {
				checktarget := recContent.(string)[len(recContent.(string))-1:]
				if checktarget == "." {
					records = append(records, strconv.Itoa(subtype)+" "+recContent.(string))
				} else {
					records = append(records, strconv.Itoa(subtype)+" "+recContent.(string)+".")
				}
			}
			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "DNSKEY" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			protocol := d.Get("protocol").(int)
			algorithm := d.Get("algorithm").(int)
			key := d.Get("key").(string)

			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(protocol)+" "+strconv.Itoa(algorithm)+" "+key)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil

		}
		if recordtype == "DS" {

			records := make([]string, 0, len(target))
			digestType := d.Get("digest_type").(int)
			keytag := d.Get("keytag").(int)
			algorithm := d.Get("algorithm").(int)
			digest := d.Get("digest").(string)

			records = append(records, strconv.Itoa(keytag)+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(digestType)+" "+digest)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "HINFO" {

			records := make([]string, 0, len(target))
			hardware := d.Get("hardware").(string)
			software := d.Get("software").(string)

			// Fields may have embedded backslash. Quotes optional
			if strings.HasPrefix(hardware, "\\\"") {
				hardware = strings.TrimLeft(hardware, "\\\"")
				hardware = strings.TrimRight(hardware, "\\\"")
				hardware = "\"" + hardware + "\""
			}
			if strings.HasPrefix(software, "\\\"") {
				software = strings.TrimLeft(software, "\\\"")
				software = strings.TrimRight(software, "\\\"")
				software = "\"" + software + "\""
			}

			records = append(records, hardware+" "+software)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "LOC" {

			records := make([]string, 0, len(target))

			for _, recContent := range target {
				records = append(records, recContent.(string))
			}
			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "MX" {

			zone := d.Get("zone").(string)

			log.Printf("[DEBUG] [Akamai DNSv2] MX record targets to process: %v", target)

			rdata, e := dnsv2.GetRdata(zone, host, recordtype)
			if e != nil {
				log.Printf("[DEBUG] [Akamai DNSv2] Searching for existing MX records no prexisting targets found LEN %d", len(rdata))
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Existing MX records to append to target %v", rdata)

			records := make([]string, 0, len(target)+len(rdata))

			priority := d.Get("priority").(int)
			increment := d.Get("priority_increment").(int)
			log.Printf("[DEBUG] [Akamai DNSv2] MX BIND Priority: %d ; Increment: %d", priority, increment)

			rdataTargetMap := make(map[string]int, len(rdata))
			for _, r := range rdata {
				entryparts := strings.Split(r, " ")
				rn := entryparts[1]
				if rn[len(rn)-1:] != "." {
					rn = rn + "."
				}
				rdataTargetMap[rn], _ = strconv.Atoi(entryparts[0])
			}

			for _, recContent := range target {
				targentry := recContent.(string)
				if targentry[len(recContent.(string))-1:] != "." {
					targentry += "."
				}
				if pri, ok := rdataTargetMap[targentry]; ok {
					log.Printf("MX BIND. %s in existing map", targentry)
					// target already in rdata
					if pri != priority {
						return dnsv2.RecordBody{}, fmt.Errorf("MX Record Priority Mismatch. Target order must align with EdgeDNS")
					} else {
						// already in rdata
						continue
					}
				} else if containsPriority(rdataTargetMap, priority) {
					log.Printf("MX BIND. %d already in RDATA", priority)
					// already an rdata record with same priority
					return dnsv2.RecordBody{}, fmt.Errorf("MX Record Priority Duplication. Target order and priority must align with EdgeDNS")
				}
				records = append(records, strconv.Itoa(priority)+" "+targentry)

				if increment > 0 {
					priority = priority + increment
				}
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Appended new target to target array LEN %d %v", len(records), records)

			for targname, tpri := range rdataTargetMap {
				records = append(records, strconv.Itoa(tpri)+" "+targname)
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Existing MX records to append to target before schema data LEN %d %v", len(rdata), records)

			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "NAPTR" {

			records := make([]string, 0, len(target))
			flagsnaptr := d.Get("flagsnaptr").(string)
			order := d.Get("order").(int)
			preference := d.Get("preference").(int)
			regexp := d.Get("regexp").(string)
			replacement := d.Get("replacement").(string)
			// Following three fields may have embedded backslash
			service := d.Get("service").(string)
			if !strings.HasPrefix(service, "\"") {
				service = "\"" + service + "\""
			}
			if !strings.HasPrefix(regexp, "\"") {
				regexp = "\"" + regexp + "\""
			}
			if !strings.HasPrefix(flagsnaptr, "\"") {
				flagsnaptr = "\"" + flagsnaptr + "\""
			}
			records = append(records, strconv.Itoa(order)+" "+strconv.Itoa(preference)+" "+flagsnaptr+" "+service+" "+regexp+" "+replacement)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "NSEC3" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			nextHashedOwnerName := d.Get("next_hashed_owner_name").(string)
			salt := d.Get("salt").(string)
			typeBitmaps := d.Get("type_bitmaps").(string)

			records = append(records, strconv.Itoa(algorithm)+" "+strconv.Itoa(flags)+" "+strconv.Itoa(iterations)+" "+salt+" "+nextHashedOwnerName+" "+typeBitmaps)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "NSEC3PARAM" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			salt := d.Get("salt").(string)

			saltbase32 := salt

			records = append(records, strconv.Itoa(algorithm)+" "+strconv.Itoa(flags)+" "+strconv.Itoa(iterations)+" "+saltbase32)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "RP" {

			records := make([]string, 0, len(target))
			mailbox := d.Get("mailbox").(string)
			checkmailbox := mailbox[len(mailbox)-1:]
			if !(checkmailbox == ".") {
				mailbox = mailbox + "."
			}
			txt := d.Get("txt").(string)
			checktxt := txt[len(txt)-1:]
			if !(checktxt == ".") {
				txt = txt + "."
			}

			records = append(records, mailbox+" "+txt)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "RRSIG" {

			records := make([]string, 0, len(target))
			expiration := d.Get("expiration").(string)
			inception := d.Get("inception").(string)
			originalTTL := d.Get("original_ttl").(int)
			algorithm := d.Get("algorithm").(int)
			labels := d.Get("labels").(int)
			keytag := d.Get("keytag").(int)
			signature := d.Get("signature").(string)
			signer := d.Get("signer").(string)
			typeCovered := d.Get("type_covered").(string)

			records = append(records, typeCovered+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(labels)+" "+strconv.Itoa(originalTTL)+" "+expiration+" "+inception+" "+strconv.Itoa(keytag)+" "+signer+" "+signature)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "SRV" {

			records := make([]string, 0, len(target))
			priority := d.Get("priority").(int)
			weight := d.Get("weight").(int)
			port := d.Get("port").(int)

			for _, recContent := range target {
				checktarget := recContent.(string)[len(recContent.(string))-1:]
				if checktarget == "." {
					records = append(records, strconv.Itoa(priority)+" "+strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+recContent.(string))
				} else {
					records = append(records, strconv.Itoa(priority)+" "+strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+recContent.(string)+".")
				}

			}
			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "SSHFP" {

			records := make([]string, 0, len(target))
			algorithm := d.Get("algorithm").(int)
			fingerprintType := d.Get("fingerprint_type").(int)
			fingerprint := d.Get("fingerprint").(string)

			records = append(records, strconv.Itoa(algorithm)+" "+strconv.Itoa(fingerprintType)+" "+fingerprint)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "SOA" {

			records := make([]string, 0, len(target))
			nameserver := d.Get("name_server").(string)
			emailaddr := d.Get("email_address").(string)
			serial := d.Get("serial").(int)
			refresh := d.Get("refresh").(int)
			retry := d.Get("retry").(int)
			expiry := d.Get("expiry").(int)
			nxdomainttl := d.Get("nxdomain_ttl").(int)

			records = append(records, nameserver+" "+emailaddr+" "+strconv.Itoa(serial)+" "+strconv.Itoa(refresh)+" "+strconv.Itoa(retry)+" "+strconv.Itoa(expiry)+" "+strconv.Itoa(nxdomainttl))
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
		if recordtype == "AKAMAITLC" {

			records := make([]string, 0, len(target))
			dnsname := d.Get("dns_name").(string)
			answtype := d.Get("answer_type").(string)

			records = append(records, answtype+" "+dnsname)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate, nil
		}
	}
	return emptyrecordcreate, fmt.Errorf("Unable to create a Record Body for %s : %s", host, recordtype)

}

func validateRecord(d *schema.ResourceData) error {
	var recordtype string
	if v, ok := d.GetOk("recordtype"); ok {
		recordtype = v.(string)
	}

	switch recordtype {
	case RRTypeA, RRTypeAaaa, RRTypeAkamaiCdn, RRTypeCname, RRTypeLoc, RRTypeNs, RRTypePtr, RRTypeSpf, RRTypeTxt:
		if err := checkBasicRecordTypes(d); err != nil {
			return err
		}
		return checkTargets(d)
	case RRTypeAfsdb:
		return checkAsdfRecord(d)
	case RRTypeDnskey:
		return checkDnskeyRecord(d)
	case RRTypeDs:
		return checkDsRecord(d)
	case RRTypeHinfo:
		return checkHinfoRecord(d)
	case RRTypeMx:
		return checkMxRecord(d)
	case RRTypeNaptr:
		return checkNaptrRecord(d)
	case RRTypeNsec3:
		return checkNsec3Record(d)
	case RRTypeNsec3Param:
		return checkNsec3ParamRecord(d)
	case RRTypeRp:
		return checkRpRecord(d)
	case RRTypeRrsig:
		return checkRrsigRecord(d)
	case RRTypeSrv:
		return checkSrvRecord(d)
	case RRTypeSshfp:
		return checkSshfpRecord(d)
	case RRTypeAkamaiTlc:
		return checkAkamaiTlcRecord(d)
	case RRTypeSoa:
		return checkSoaRecord(d)
	default:
		return fmt.Errorf("Invalid recordtype %v", recordtype)
	}
}

func checkBasicRecordTypes(d *schema.ResourceData) error {
	host := d.Get("name").(string)
	recordtype := d.Get("recordtype").(string)
	ttl := d.Get("ttl").(int)

	if host == "" {
		return fmt.Errorf("Type host must be set")
	}

	if recordtype == "" {
		return fmt.Errorf("Type recordtype must be set")
	}

	if ttl == 0 {
		return fmt.Errorf("Type ttl must be set")
	}

	return nil
}

func checkTargets(d *schema.ResourceData) error {
	target := d.Get("target").(*schema.Set).List()
	records := make([]string, 0, len(target))

	for _, recContent := range target {
		records = append(records, recContent.(string))
	}

	if len(records) == 0 {
		return fmt.Errorf("Type records must be set.")
	}

	return nil
}

func checkSimpleRecord(d *schema.ResourceData) error {
	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	return nil
}

func checkAsdfRecord(d *schema.ResourceData) error {
	subtype := d.Get("subtype").(int)
	if subtype == 0 {
		return fmt.Errorf("Type subtype must be set for ASDF.")
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	return nil
}

func checkDnskeyRecord(d *schema.ResourceData) error {
	flags := d.Get("flags").(int)
	protocol := d.Get("protocol").(int)
	algorithm := d.Get("algoritm").(int)
	key := d.Get("key").(string)
	ttl := d.Get("ttl").(int)

	if !(flags == 0 || flags == 256 || flags == 257) {
		return fmt.Errorf("Type flags must not be %v for DNSKEY.", flags)
	}

	if ttl == 0 {
		return fmt.Errorf("Type ttl must be set for DNSKEY.")
	}

	if protocol == 0 {
		return fmt.Errorf("Type protocol must be set for DNSKEY.")
	}

	if !((algorithm >= 1 && algorithm <= 8) || algorithm != 10) {
		return fmt.Errorf("Type algorithm must not be %v for DNSKEY.", algorithm)
	}

	if key == "" {
		return fmt.Errorf("Type key must be set for DNSKEY.")
	}

	return nil
}

func checkDsRecord(d *schema.ResourceData) error {
	digestType := d.Get("digest_type").(int)
	keytag := d.Get("keytag").(int)
	algorithm := d.Get("algorithm").(int)
	digest := d.Get("digest").(string)

	if digestType == 0 {
		return fmt.Errorf("Type digest_type must be set for DS.")
	}

	if keytag == 0 {
		return fmt.Errorf("Type keytag must be set for DS.")
	}

	if algorithm == 0 {
		return fmt.Errorf("Type algorithm must be set for DS.")
	}

	if digest == "" {
		return fmt.Errorf("Type digest must be set for DS.")
	}

	return nil
}

func checkHinfoRecord(d *schema.ResourceData) error {
	hardware := d.Get("hardware").(string)
	software := d.Get("software").(string)

	if hardware == "" {
		return fmt.Errorf("Type hardware must be set for HINFO.")
	}

	if software == "" {
		return fmt.Errorf("Type software must be set for HINFO.")
	}

	return nil
}

func checkMxRecord(d *schema.ResourceData) error {
	priority := d.Get("priority").(int)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if priority < 0 || priority > 65535 {
		return fmt.Errorf("Type priority must be set for MX.")
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	return nil
}

func checkNaptrRecord(d *schema.ResourceData) error {
	flagsnaptr := d.Get("flagsnaptr").(string)
	order := d.Get("order").(int)
	preference := d.Get("preference").(int)
	regexp := d.Get("regexp").(string)
	replacement := d.Get("replacement").(string)
	service := d.Get("service").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if flagsnaptr == "" {
		return fmt.Errorf("Type flagsnaptr must be set for NAPTR.")
	}

	if order < 0 || order > 65535 {
		return fmt.Errorf("Type order must not be %v for NAPTR.", order)
	}

	if preference == 0 {
		return fmt.Errorf("Type preference must be set for NAPTR.")
	}

	if regexp == "" {
		return fmt.Errorf("Type regexp must be set for NAPTR.")
	}

	if replacement == "" {
		return fmt.Errorf("Type replacement must be set for NAPTR.")
	}

	if service == "" {
		return fmt.Errorf("Type service must be set for NAPTR.")
	}

	return nil
}

func checkNsec3Record(d *schema.ResourceData) error {
	flags := d.Get("flags").(int)
	algorithm := d.Get("algorithm").(int)
	iterations := d.Get("iterations").(int)
	nextHashedOwnerName := d.Get("next_hashed_owner_name").(string)
	salt := d.Get("salt").(string)
	typeBitmaps := d.Get("type_bitmaps").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if !(flags == 0 || flags == 1) {
		return fmt.Errorf("Type flags must be set for NSEC3.")
	}

	if algorithm != 1 {
		return fmt.Errorf("Type flags must be set for NSEC3.")
	}
	if iterations == 0 {
		return fmt.Errorf("Type iterations must be set for NSEC3.")
	}
	if nextHashedOwnerName == "" {
		return fmt.Errorf("Type nextHashedOwnerName must be set for NSEC3.")
	}
	if salt == "" {
		return fmt.Errorf("Type salt must be set for NSEC3.")
	}
	if typeBitmaps == "" {
		return fmt.Errorf("Type typeBitMaps must be set for NSEC3.")
	}
	return nil
}

func checkNsec3ParamRecord(d *schema.ResourceData) error {
	flags := d.Get("flags").(int)
	algorithm := d.Get("algorithm").(int)
	iterations := d.Get("iterations").(int)
	salt := d.Get("salt").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if !(flags == 0 || flags == 1) {
		return fmt.Errorf("Type flags must be set for NSEC3PARAM.")
	}

	if algorithm != 1 {
		return fmt.Errorf("Type algorithm must be set for NSEC3PARAM.")
	}

	if iterations == 0 {
		return fmt.Errorf("Type iterations must be set for NSEC3PARAM.")
	}

	if salt == "" {
		return fmt.Errorf("Type salt must be set for NSEC3PARAM.")
	}

	return nil
}

func checkRpRecord(d *schema.ResourceData) error {
	mailbox := d.Get("mailbox").(string)
	txt := d.Get("txt").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if mailbox == "" {
		return fmt.Errorf("Type mailbox must be set for RP.")
	}

	if txt == "" {
		return fmt.Errorf("Type txt must be set for RP.")
	}

	return nil
}

func checkRrsigRecord(d *schema.ResourceData) error {
	expiration := d.Get("expiration").(string)
	inception := d.Get("inception").(string)
	originalTTL := d.Get("original_ttl").(int)
	algorithm := d.Get("algorithm").(int)
	labels := d.Get("labels").(int)
	keytag := d.Get("keytag").(int)
	signature := d.Get("signature").(string)
	signer := d.Get("signer").(string)
	typeCovered := d.Get("type_covered").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if expiration == "" {
		return fmt.Errorf("Type expiration must be set for RRSIG.")
	}

	if inception == "" {
		return fmt.Errorf("Type inception must be set for RRSIG.")
	}

	if originalTTL == 0 {
		return fmt.Errorf("Type originalTTL must be set for RRSIG.")
	}

	if algorithm == 0 {
		return fmt.Errorf("Type algorithm must be set for RRSIG.")
	}

	if labels == 0 {

		return fmt.Errorf("Type labels must be set for RRSIG.")
	}

	if keytag == 0 {
		return fmt.Errorf("Type keytag must be set for RRSIG.")
	}

	if signature == "" {
		return fmt.Errorf("Type signature must be set for RRSIG.")
	}

	if signer == "" {
		return fmt.Errorf("Type signer must be set for RRSIG.")
	}

	if typeCovered == "" {
		return fmt.Errorf("Type typeCovered must be set for RRSIG.")
	}

	return nil
}

func checkSrvRecord(d *schema.ResourceData) error {
	priority := d.Get("priority").(int)
	weight := d.Get("weight").(int)
	port := d.Get("port").(int)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	if priority == 0 {
		return fmt.Errorf("Type priority must be set for SRV.")
	}

	if weight < 0 || weight > 65535 {
		return fmt.Errorf("Type weight must not be %v for SRV.", weight)
	}

	if port == 0 {
		return fmt.Errorf("Type port must be set for SRV.")
	}

	return nil
}

func checkSshfpRecord(d *schema.ResourceData) error {
	algorithm := d.Get("algorithm").(int)
	fingerprintType := d.Get("fingerprint_type").(int)
	fingerprint := d.Get("fingerprint").(string)

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if algorithm == 0 {
		return fmt.Errorf("Type algorithm must be set for SSHFP.")
	}

	if fingerprintType == 0 {
		return fmt.Errorf("Type fingerprintType must be set for SSHFP.")
	}

	if fingerprint == "null" {
		return fmt.Errorf("Type fingerprint must be set for SSHFP.")
	}

	return nil
}

func checkSoaRecord(d *schema.ResourceData) error {

	nameserver := d.Get("name_server").(string)
	emailaddr := d.Get("email_address").(string)
	serial := d.Get("serial").(int)
	refresh := d.Get("refresh").(int)
	retry := d.Get("retry").(int)
	expiry := d.Get("expiry").(int)
	nxdomainttl := d.Get("nxdomain_ttl").(int)

	if nameserver == "" {
		return fmt.Errorf("Key %s must be specified in SOA record", nameserver)
	}

	if emailaddr == "" {
		return fmt.Errorf("Key %s must be specified in SOA record", emailaddr)
	}

	if serial == 0 {
		return fmt.Errorf("Key %d must be specified in SOA record", serial)
	}

	if refresh == 0 {
		return fmt.Errorf("Key %d must be specified in SOA record", refresh)
	}

	if retry == 0 {
		return fmt.Errorf("Key %d must be specified in SOA record", retry)
	}

	if expiry == 0 {
		return fmt.Errorf("Key %d must be specified in SOA record", expiry)
	}

	if nxdomainttl == 0 {
		return fmt.Errorf("Key %d must be specified in SOA record", nxdomainttl)
	}

	return nil
}

func checkAkamaiTlcRecord(d *schema.ResourceData) error {
	dnsname := d.Get("dns_name").(string)
	answertype := d.Get("answer_type").(string)

	if dnsname != "" {
		return fmt.Errorf("dnsname key is computed. It must not be set in AKAMAITLC.")
	}

	if answertype != "" {
		return fmt.Errorf("answertype key is computed. It must not be set in AKAMAITLC.")
	}

	return nil
}

// Resource record types supported by the Akamai Edge DNS API
const (
	RRTypeA          = "A"
	RRTypeAaaa       = "AAAA"
	RRTypeAfsdb      = "AFSDB"
	RRTypeAkamaiCdn  = "AKAMAICDN"
	RRTypeAkamaiTlc  = "AKAMAITLC"
	RRTypeCaa        = "CAA"
	RRTypeCname      = "CNAME"
	RRTypeHinfo      = "HINFO"
	RRTypeLoc        = "LOC"
	RRTypeMx         = "MX"
	RRTypeNaptr      = "NAPTR"
	RRTypeNs         = "NS"
	RRTypePtr        = "PTR"
	RRTypeRp         = "RP"
	RRTypeSoa        = "SOA"
	RRTypeSrv        = "SRV"
	RRTypeSpf        = "SPF"
	RRTypeSshfp      = "SSHFP"
	RRTypeTlsa       = "TLSA"
	RRTypeTxt        = "TXT"
	RRTypeDnskey     = "DNSKEY"
	RRTypeDs         = "DS"
	RRTypeNsec3      = "NSEC3"
	RRTypeNsec3Param = "NSEC3PARAM"
	RRTypeRrsig      = "RRSIG"
)

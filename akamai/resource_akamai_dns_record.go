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

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDNSv2Record() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordUpdate,
		Delete: resourceDNSRecordDelete,
		Exists: resourceDNSRecordExists,
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
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
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
			},
			"software": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
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
		},
	}
}

// Create a new DNS Record
func resourceDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
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

	recordcreate := bindRecord(d)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)

	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for recordcreate [%s]", sha1hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records LEN %d", len(rdata))
	if len(rdata) > 0 {
		extractString := strings.Join(rdata, " ")
		sha1hashtest := getSHAString(extractString)

		log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordread [%s]", sha1hashtest)
	}
	// If there's no existing record we'll create a blank one
	if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
		// if the record is not found/404 we will create a new
		log.Printf("[DEBUG] [Akamai DNSv2] [ERROR] %s", e.Error())
		log.Printf("[DEBUG] [Akamai DNSv2] Creating new record")
		// Save the zone to the API
		e = recordcreate.Save(zone)

		if e != nil {
			return e
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] Updating record")
		if len(rdata) > 0 {
			e = recordcreate.Update(zone)
			if e != nil {
				return e
			}
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] Saving record")
			e = recordcreate.Save(zone)
			if e != nil {
				return e
			}
		}

	}

	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), ":") {
		d.SetId(fmt.Sprintf("%s:%s:%s:%s", zone, host, recordtype, sha1hash))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
	}

	return resourceDNSRecordUpdate(d, meta)
}

// Update DNS Record
func resourceDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	// only allow one record to be created at a time
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
	}

	err := validateRecord(d)
	if err != nil {
		return fmt.Errorf("DNS record validation failure on zone %v: %v", zone, err)
	}

	recordcreate := bindRecord(d)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)

	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE SHA sum for recordupdate [%s]", sha1hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Searching for records LEN %d", len(rdata))
	if len(rdata) > 0 {
		sort.Strings(rdata)
		extractString := strings.Join(recordcreate.Target, " ")
		sha1hashtest := getSHAString(extractString)
		log.Printf("[DEBUG] [Akamai DNSv2] UPDATE SHA sum from recordread [%s]", sha1hashtest)
		// If there's no existing record we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the record is not found/404 we will create a new
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE [ERROR] %s", e.Error())
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Creating new record")
			// Save the zone to the API
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Updating record")
			e = recordcreate.Save(zone)

		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] UPDATE Updating record")
			e = recordcreate.Update(zone)
			if e != nil {
				return e
			}

		}
		// Give terraform the ID
		if d.Id() == "" || strings.Contains(d.Id(), ":") {
			d.SetId(fmt.Sprintf("%s:%s:%s:%s", zone, host, recordtype, sha1hash))
		} else {
			d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
		}
	}

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

	recordcreate := bindRecord(d)

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
	targets, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		//return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e,e)
		d.SetId("")
		return nil

	}
	b1, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] READ record data read JSON %s", string(b1))

	if len(targets) > 0 {
		sort.Strings(targets)
		extractStringTest := strings.Join(targets, " ")
		sha1hashtest := getSHAString(extractStringTest)

		if sha1hashtest == sha1hash {
			log.Printf("[DEBUG] [Akamai DNSv2] READ SHA sum from recordExists matches [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			// Give terraform the ID
			if strings.Contains(d.Id(), ":") {
				d.SetId(fmt.Sprintf("%s:%s:%s:%s", zone, host, recordtype, sha1hash))
			} else {
				d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
			}
			return nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] READ SHA sum from recordExists mismatch [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			d.SetId("")
			return nil
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] READ SHA sum from recordExists mismatch no target returned  [%s] [%s] [%s] ", zone, host, recordtype)
		d.SetId("")
		return nil
	}
	return nil
}

func resourceDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	idParts := strings.Split(d.Id(), ":")
	fmt.Println("idParts: ", idParts)
	if len(idParts) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("Invalid Id for Zone Import: %s", d.Id())
	}
	zone := idParts[0]
	recordname := idParts[1]
	recordtype := idParts[2]

	// Get recordset
	log.Printf("[INFO] [Akamai DNS] Searching for zone Recordset [%v]", idParts)
	recordset, err := dnsv2.GetRecord(zone, recordname, recordtype)
	if err != nil {
		log.Printf("[DEBUG] [Akamai DNSv2] IMPORT Record read failed for record [%s] [%s] [%s] ", zone, recordname, recordtype)
		d.SetId("")
		return []*schema.ResourceData{d}, err
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
	targets := recordset.Target
	importTargetString := "no rdata"
	if len(targets) > 0 {
		sort.Strings(targets)
		importTargetString = strings.Join(targets, " ")
	}

	sha1hash := getSHAString(importTargetString)
	d.SetId(fmt.Sprintf("%s:%s:%s:%s", zone, recordname, recordtype, sha1hash))

	return []*schema.ResourceData{d}, nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {

	zone := d.Get("zone").(string)
	host := d.Get("name").(string)
	recordtype := d.Get("recordtype").(string)
	ttl := d.Get("ttl").(int)

	target := d.Get("target").(*schema.Set).List()

	records := make([]string, 0, len(target))
	for _, recContent := range target {
		records = append(records, recContent.(string))
	}
	sort.Strings(records)
	log.Printf("[INFO] [Akamai DNS] Delete zone Records %v", records)
	recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}

	recordcreate.Delete(zone)

	targets, e := dnsv2.GetRdata(zone, host, recordtype)
	if len(targets) > 0 {
		log.Printf("[INFO] [Akamai DNS] Delete zone Records record still exists %v %s", targets, e)
		return nil
	} else {
		d.SetId("")
	}
	return nil
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

	recordcreate := bindRecord(d)

	b, err := json.Marshal(recordcreate.Target)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] EXISTS record JSON from bind records %s %s %s %s", string(b), zone, host, recordtype)
	//sort.Strings(recordcreate.Target)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)
	log.Printf("[DEBUG] [Akamai DNSv2] EXISTS SHA sum for Existing SHA test %s %s", extractString, sha1hash)

	// try to get the zone from the API
	log.Printf("[INFO] [Akamai DNSv2] EXISTS Searching for zone records %s %s %s", zone, host, recordtype)
	targets, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		//return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e,e)
		return false, nil
	}
	b1, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] EXISTS record data read JSON %s", string(b1))

	if len(targets) > 0 {
		sort.Strings(targets)
		extractStringTest := strings.Join(targets, " ")
		sha1hashtest := getSHAString(extractStringTest)

		if sha1hashtest == sha1hash {
			log.Printf("[DEBUG] [Akamai DNSv2] EXISTS SHA sum from recordExists matches [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			return true, nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] EXISTS SHA sum from recordExists mismatch [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			return false, nil
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] EXISTS SHA sum from recordExists msimatch no target retunred  [%s] [%s] [%s] ", zone, host, recordtype)
		return false, nil
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
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

func bindRecord(d *schema.ResourceData) dnsv2.RecordBody {

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

	simplerecordtarget := map[string]bool{"AAAA": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "SRV": true}

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
		return recordcreate
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
			return recordcreate
		}
		if recordtype == "DNSKEY" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			protocol := d.Get("protocol").(int)
			algorithm := d.Get("algorithm").(int)
			key := d.Get("key").(string)

			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(protocol)+" "+strconv.Itoa(algorithm)+" "+key)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "DS" {

			records := make([]string, 0, len(target))
			digestType := d.Get("digest_type").(int)
			keytag := d.Get("keytag").(int)
			algorithm := d.Get("algorithm").(int)
			digest := d.Get("digest").(string)

			records = append(records, strconv.Itoa(keytag)+" "+strconv.Itoa(digestType)+" "+strconv.Itoa(algorithm)+" "+digest)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "HINFO" {

			records := make([]string, 0, len(target))
			hardware := d.Get("hardware").(string)
			software := d.Get("software").(string)

			records = append(records, hardware+" "+software)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "LOC" {

			records := make([]string, 0, len(target))

			for _, recContent := range target {
				records = append(records, recContent.(string))
			}
			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "MX" {

			zone := d.Get("zone").(string)

			rdata, e := dnsv2.GetRdata(zone, host, recordtype)
			if e != nil {
				log.Printf("[DEBUG] [Akamai DNSv2] Searching for existing MX records no prexisting targets found LEN %d", len(rdata))
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Searching for existing MX records to append to target LEN %d", len(rdata))

			records := make([]string, 0, len(target)+len(rdata))
			priority := d.Get("priority").(int)

			increment := d.Get("priority_increment").(int)

			for _, recContent := range target {
				checktarget := recContent.(string)[len(recContent.(string))-1:]
				if checktarget != "." {
					records = append(records, strconv.Itoa(priority)+" "+recContent.(string)+".")

				} else {
					records = append(records, strconv.Itoa(priority)+" "+recContent.(string))
				}
				if increment > 0 {
					priority = priority + increment
				}
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Appended new target to target array LEN %d %v", len(records), records)

			if len(rdata) > 0 {
				log.Printf("[DEBUG] [Akamai DNSv2] rdata Exists MX records to append to target LEN %d", len(rdata))
				for _, r := range rdata {
					checktarget := r[len(r)-1:]
					if checktarget != "." {
						r = r + "."
					}
					if !(contains(records, r)) {
						records = append(records, r)
					}
				}
				log.Printf("[DEBUG] [Akamai DNSv2] Existing MX records to append to target before schema data LEN %d %v", len(rdata), records)

			}
			sort.Strings(records)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "NAPTR" {

			records := make([]string, 0, len(target))
			flagsnaptr := d.Get("flagsnaptr").(string)
			order := d.Get("order").(int)
			preference := d.Get("preference").(int)
			regexp := d.Get("regexp").(string)
			replacement := d.Get("replacement").(string)
			service := d.Get("service").(string)
			checktarget := service[len(service)-1:]
			if !(checktarget == ".") {
				service = service + "."
			}

			records = append(records, strconv.Itoa(order)+" "+strconv.Itoa(preference)+" "+flagsnaptr+" "+regexp+" "+replacement+" "+service)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "NSEC3" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			nextHashedOwnerName := d.Get("next_hashed_owner_name").(string)
			salt := d.Get("salt").(string)
			typeBitmaps := d.Get("type_bitmaps").(string)

			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(iterations)+" "+salt+" "+nextHashedOwnerName+" "+typeBitmaps)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "NSEC3PARAM" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			salt := d.Get("salt").(string)

			saltbase32 := salt

			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(iterations)+" "+saltbase32)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
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
			return recordcreate
		}
		if recordtype == "RRSIG" { //TODO FIX

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

			records = append(records, typeCovered+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(labels)+" "+strconv.Itoa(originalTTL)+" "+expiration+" "+inception+" "+signature+" "+signer+" "+strconv.Itoa(keytag))

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
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
			return recordcreate
		}
		if recordtype == "SSHFP" {

			records := make([]string, 0, len(target))
			algorithm := d.Get("algorithm").(int)
			fingerprintType := d.Get("fingerprint_type").(int)
			fingerprint := d.Get("fingerprint").(string)

			records = append(records, strconv.Itoa(algorithm)+" "+strconv.Itoa(fingerprintType)+" "+fingerprint)

			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
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
			return recordcreate
		}
		if recordtype == "AKAMAITLC" {

			records := make([]string, 0, len(target))
			dnsname := d.Get("dns_name").(string)
			answtype := d.Get("answer_type").(string)

			records = append(records, answtype+" "+dnsname)
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
	}
	return emptyrecordcreate
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
		return fmt.Errorf("Key %s must ne specified in SOA record", nameserver)
	}

	if emailaddr == "" {
		return fmt.Errorf("Key %s must ne specified in SOA record", emailaddr)
	}

	if serial == 0 {
		return fmt.Errorf("Key %s must ne specified in SOA record", serial)
	}

	if refresh == 0 {
		return fmt.Errorf("Key %s must ne specified in SOA record", refresh)
	}

	if retry == 0 {
		return fmt.Errorf("Key %s must ne specified in SOA record", retry)
	}

	if expiry == 0 {
		return fmt.Errorf("Key %s must ne specified in SOA record", expiry)
	}

	if nxdomainttl == 0 {
		return fmt.Errorf("Key %s must ne specified in SOA record", nxdomainttl)
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

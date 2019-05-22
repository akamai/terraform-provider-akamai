package akamai

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"target": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			//	"afsdb":
			"subtype": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			//	"cname":
			//	"dnskey":
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
			//	"ds":
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
			//"hinfo":
			"hardware": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"software": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//	 "loc":
			//		"mx":
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			//	"naptr":
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
			//"ns":
			//"nsec3":
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
			//"nsec3param":
			//	"ptr":
			//	"rp":
			"mailbox": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"txt": {
				Type:     schema.TypeString,
				Optional: true,
			},
			//	"rrsig":
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
			//"spf":
			//	"srv":
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			//"sshfp": {
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
			//"txt": {
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

	validationresult := validateRecord(d)
	log.Printf("[DEBUG] [Akamai DNSv2] Validation result recordcreate %s", validationresult)
	if validationresult != "VALID" {
		return fmt.Errorf("Parameter Validation failure %s, %s  %s %s", zone, host, recordtype, validationresult)
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
	d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))

	return nil
}

// Create a new DNS Record
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

	vaidationresult := validateRecord(d)
	log.Printf("[DEBUG] [Akamai DNSv2] Validation result recordupdate %s", vaidationresult)
	if vaidationresult != "VALID" {
		return fmt.Errorf("Parameter Validation failure %s, %s  %s %s", zone, host, recordtype, vaidationresult)
	}
	recordcreate := bindRecord(d)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)

	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for recordupdate [%s]", sha1hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records LEN %d", len(rdata))
	if len(rdata) > 0 {
		sort.Strings(rdata)
		extractString := strings.Join(recordcreate.Target, " ")
		sha1hashtest := getSHAString(extractString)
		log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordread [%s]", sha1hashtest)
		// If there's no existing record we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the record is not found/404 we will create a new
			log.Printf("[DEBUG] [Akamai DNSv2] [ERROR] %s", e.Error())
			log.Printf("[DEBUG] [Akamai DNSv2] Creating new record")
			// Save the zone to the API
			log.Printf("[DEBUG] [Akamai DNSv2] Updating record")
			e = recordcreate.Save(zone)

		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] Updating record")
			e = recordcreate.Update(zone)
			if e != nil {
				return e
			}

		}
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1hash))
	}

	// Give terraform the ID

	return nil
}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()

	// find the zone first
	log.Printf("[INFO] [Akamai DNS] Searching for zone Records [%s]", hostname)

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

	log.Printf("[DEBUG] [Akamai DNSv2] record JSON from bind records %s %s %s %s", string(b), zone, host, recordtype)
	extractString := strings.Join(recordcreate.Target, " ")
	sha1hash := getSHAString(extractString)
	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for Existing SHA test %s %s", extractString, sha1hash)

	// try to get the zone from the API
	log.Printf("[INFO] [Akamai DNSv2] Searching for zone records %s %s %s", zone, host, recordtype)
	targets, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		//return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e,e)
		return false, nil
	}
	b1, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] record data read JSON %s", string(b1))

	if len(targets) > 0 {
		sort.Strings(targets)
		extractStringTest := strings.Join(targets, " ")
		sha1hashtest := getSHAString(extractStringTest)

		if sha1hashtest == sha1hash {
			log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists matches [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			return true, nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists mismatch [%s] vs  [%s] [%s] [%s] [%s] ", sha1hashtest, sha1hash, zone, host, recordtype)
			return false, nil
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists msimatch no target retunred  [%s] [%s] [%s] ", zone, host, recordtype)
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

	simplerecord := map[string]bool{"A": true, "AAAA": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "TXT": true}
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
				records = append(records, strconv.Itoa(priority)+" "+recContent.(string))
				if increment > 0 {
					priority = priority + increment
				}
			}
			log.Printf("[DEBUG] [Akamai DNSv2] Appended new target to taget array LEN %d %v", len(records), records)

			if len(rdata) > 0 {
				log.Printf("[DEBUG] [Akamai DNSv2] rdata Exists MX records to append to target LEN %d", len(rdata))
				for _, r := range rdata {
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
					records = append(records, strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+strconv.Itoa(priority)+" "+recContent.(string))
				} else {
					records = append(records, strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+strconv.Itoa(priority)+" "+recContent.(string)+".")
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
	}
	return emptyrecordcreate
}

func validateRecord(d *schema.ResourceData) string {

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

	for _, recContent := range target {
		records = append(records, recContent.(string))
	}

	simplerecord := map[string]bool{"A": true, "AAAA": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "TXT": true}
	if simplerecord[recordtype] {

		if host == "null" {
			return "host"
		}
		if recordtype == "null" {
			return "recordtype"
		}
		if ttl == 0 {
			return "ttl"
		}
		if len(records) == 0 {
			return "target"
		}
		return "VALID"
	} else {
		if recordtype == "AFSDB" {

			subtype := d.Get("subtype").(int)
			if subtype == 0 {
				return "subtype"
			}
			if len(records) == 0 {
				return "target"
			}
			return "VALID"
		}
		if recordtype == "DNSKEY" {

			flags := d.Get("flags").(int)
			protocol := d.Get("protocol").(int)
			algorithm := d.Get("algorithm").(int)
			key := d.Get("key").(string)
			log.Printf("[DEBUG] [Akamai DNSv2] DNSKEY FLAGS %d ", flags)
			if flags == 0 || flags == 256 || flags == 257 {
				log.Printf("[DEBUG] [Akamai DNSv2] INSIDE IF DNSKEY FLAGS %d ", flags)
			} else {
				return "flags"
			}
			if ttl == 0 {
				return "ttl"
			}
			if protocol == 0 {
				return "protocol"
			}
			log.Printf("[DEBUG] [Akamai DNSv2] ALGORITHM FLAGS %d ", algorithm)
			if !((algorithm >= 1 && algorithm <= 8) || algorithm != 10) {
				return "algorithm"
			}
			if key == "null" {
				return "key"
			}
			return "VALID"
		}
		if recordtype == "DS" {
			digestType := d.Get("digest_type").(int)
			keytag := d.Get("keytag").(int)
			algorithm := d.Get("algorithm").(int)
			digest := d.Get("digest").(string)

			if digestType == 0 {
				return "digest_type"
			}
			if keytag == 0 {
				return "keytag"
			}
			if algorithm == 0 {
				return "algorithm"
			}
			if digest == "null" {
				return "digest"
			}

			return "VALID"
		}
		if recordtype == "HINFO" {

			hardware := d.Get("hardware").(string)
			software := d.Get("software").(string)

			if hardware == "null" {
				return "hardware"
			}
			if software == "null" {
				return "software"
			}

			return "VALID"
		}
		if recordtype == "LOC" {

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if len(target) == 0 {
				return "target"
			}

			return "VALID"
		}
		if recordtype == "MX" {
			priority := d.Get("priority").(int)

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if priority < 0 || priority > 65535 {
				return "priority"
			}
			if len(target) == 0 {
				return "target"
			}

			return "VALID"
		}
		if recordtype == "NAPTR" {

			flagsnaptr := d.Get("flagsnaptr").(string)
			order := d.Get("order").(int)
			preference := d.Get("preference").(int)
			regexp := d.Get("regexp").(string)
			replacement := d.Get("replacement").(string)
			service := d.Get("service").(string)
			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if flagsnaptr == "null" {
				return "flagsnaptr"
			}
			if order < 0 || order > 65535 {
				return "order"
			}
			if preference == 0 {
				return "preference"
			}
			if regexp == "null" {
				return "regexp"
			}
			if replacement == "null" {
				return "replacement"
			}
			if service == "null" {
				return "service"
			}

			return "VALID"
		}
		if recordtype == "NSEC3" {

			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			nextHashedOwnerName := d.Get("next_hashed_owner_name").(string)
			salt := d.Get("salt").(string)
			typeBitmaps := d.Get("type_bitmaps").(string)

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if !(flags == 0 || flags == 1) {
				log.Printf("[DEBUG] [Akamai DNSv2] NSEC3 FLAGS %d ", flags)
				return "flags"
			}
			if algorithm != 1 {
				return "algorithm"
			}
			if iterations == 0 {
				return "iterations"
			}
			if nextHashedOwnerName == "null" {
				return "next_hashed_owner_name"
			}
			if salt == "null" {
				return "salt"
			}
			if typeBitmaps == "null" {
				return "type_bitmaps"
			}
			return "VALID"
		}
		if recordtype == "NSEC3PARAM" {

			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			salt := d.Get("salt").(string)

			saltbase32 := salt

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if !(flags == 0 || flags == 1) {
				return "flags"
			}
			if algorithm != 1 {
				return "algorithm"
			}
			if iterations == 0 {
				return "iterations"
			}
			if salt == "null" {
				return "salt"
			}

			if saltbase32 == "null" {
				return "saltbase32"
			}
			return "VALID"
		}
		if recordtype == "RP" {

			mailbox := d.Get("mailbox").(string)
			txt := d.Get("txt").(string)

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if mailbox == "null" {
				return "mailbox"
			}
			if txt == "null" {
				return "txt"
			}
			return "VALID"
		}
		if recordtype == "RRSIG" { //TODO FIX

			expiration := d.Get("expiration").(string)
			inception := d.Get("inception").(string)
			originalTTL := d.Get("original_ttl").(int)
			algorithm := d.Get("algorithm").(int)
			labels := d.Get("labels").(int)
			keytag := d.Get("keytag").(int)
			signature := d.Get("signature").(string)
			signer := d.Get("signer").(string)
			typeCovered := d.Get("type_covered").(string)
			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if expiration == "null" {
				return "expiration"
			}
			if inception == "null" {
				return "inception"
			}
			if originalTTL == 0 {
				return "original_ttl"
			}
			if algorithm == 0 {
				return "algorithm"
			}
			if labels == 0 {
				return "labels"
			}
			if keytag == 0 {
				return "keytag"
			}
			if signature == "null" {
				return "signature"
			}
			if signer == "null" {
				return "signer"
			}
			if typeCovered == "null" {
				return "type_covered"
			}
			return "VALID"
		}
		if recordtype == "SRV" {

			records := make([]string, 0, len(target))
			priority := d.Get("priority").(int)
			weight := d.Get("weight").(int)
			port := d.Get("port").(int)

			for _, recContent := range target {
				records = append(records, strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+strconv.Itoa(priority)+" "+recContent.(string))
			}
			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if priority == 0 {
				return "priority"
			}
			if weight < 0 || weight > 65535 {
				return "weight"
			}
			if port == 0 {
				return "port"
			}
			if len(target) <= 0 {
				return "target"
			}

			return "VALID"
		}
		if recordtype == "SSHFP" {

			algorithm := d.Get("algorithm").(int)
			fingerprintType := d.Get("fingerprint_type").(int)
			fingerprint := d.Get("fingerprint").(string)

			if host == "null" {
				return "host"
			}
			if recordtype == "null" {
				return "recordtype"
			}
			if ttl == 0 {
				return "ttl"
			}
			if algorithm == 0 {
				return "algorithm"
			}
			if fingerprintType == 0 {
				return "fingerprint_type"
			}
			if fingerprint == "null" {
				return "fingerprint"
			}
			return "VALID"
		}
	}
	return "INVALID"
}

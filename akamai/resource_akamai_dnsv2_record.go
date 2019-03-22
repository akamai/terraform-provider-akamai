package akamai

import (
	"bytes"
	"crypto/sha1"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
)

func resourceDNSv2Record() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordCreate,
		Read:   resourceDNSRecordRead,
		Update: resourceDNSRecordCreate,
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
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
				//	Required: false,
				Optional: true,
				//ForceNew: true,
				Set: schema.HashString,
			},
			//	"afsdb":
			"subtype": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			//	"cname":
			//	"dnskey":
			"flags": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"protocol": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"algorithm": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//	"ds":
			"keytag": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"digest_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"digest": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//"hinfo":
			"hardware": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"software": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//	 "loc":
			//		"mx":
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			//	"naptr":
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"preference": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"flagsnaptr": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"regexp": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"replacement": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//"ns":
			//"nsec3":
			"iterations": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"salt": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"next_hashed_owner_name": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"type_bitmaps": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//"nsec3param":
			//	"ptr":
			//	"rp":
			"mailbox": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"txt": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			//	"rrsig":
			"type_covered": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"original_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"expiration": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"inception": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"signer": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"signature": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"labels": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			/*"soa":
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf(
						"%s-%s-%s-%s-%s-%s-%s",
						m["ttl"],
						m["originserver"],
						m["contact"],
						m["refresh"],
						m["retry"],
						m["expire"],
						m["minimum"],
					))
					return hashcode.String(buf.String())
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"originserver": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"contact": {
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
						"expire": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minimum": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},*/
			//"spf":
			//	"srv":
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			//"sshfp": {
			"fingerprint_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
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
	_, ok = d.GetOk("target")
	if ok {
		target := d.Get("target").(*schema.Set).List()

		records := make([]string, 0, len(target))
		for _, recContent := range target {
			records = append(records, recContent.(string))
		}
	}

	vaidationresult := validateRecord(d)
	log.Printf("[DEBUG] [Akamai DNSv2] Validation result recordcreate %s", vaidationresult)
	if vaidationresult != "VALID" {
		return errors.New(fmt.Sprintf("Parameter Validation failure %s, %s  %s %s", zone, host, recordtype, vaidationresult))
	}
	recordcreate := bindRecord(d)
	sha1_hash := getSHA(recordcreate.Target)
	/*
		  h := sha1.New()
			bodyBytes := new(bytes.Buffer)
			json.NewEncoder(bodyBytes).Encode(recordcreate.Target)
			h.Write(bodyBytes.Bytes())
			sha1_hash := hex.EncodeToString(h.Sum(nil))*/
	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for recordcreate [%s]", sha1_hash)
	// First try to get the zone from the API
	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%s]", zone)

	rdata, e := dnsv2.GetRdata(zone, host, recordtype)
	if e != nil {
		return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, e)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records LEN %d", len(rdata))
	if len(rdata) == 0 {
		sha1_hash_test := getSHA(rdata)
		log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordread [%s]", sha1_hash_test)
		// If there's no existing record we'll create a blank one
		if dnsv2.IsConfigDNSError(e) && e.(dnsv2.ConfigDNSError).NotFound() == true {
			// if the record is not found/404 we will create a new
			log.Printf("[DEBUG] [Akamai DNSv2] [ERROR] %s", e.Error())
			log.Printf("[DEBUG] [Akamai DNSv2] Creating new record")
			// Save the zone to the API

			e = nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] Saving record")
			e = recordcreate.Save(zone)
			if e != nil {
				return e
			}
			//return e
		}
	}

	// Give terraform the ID
	d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordtype, sha1_hash))

	return nil
}

/*
// Helper function for unmarshalResourceData() below
func assignFields(record dns.DNSRecord, d map[string]interface{}) {
	f := record.GetAllowedFields()
	for _, field := range f {
		val, ok := d[field]
		if ok {
			e := record.SetField(field, val)
			if e != nil {
				log.Printf("[WARN] [Akamai DNS] Couldn't add field to record: %s", e.Error())
			}
		}
	}
}
*/

// Sometimes records exist in the API but not in tf config
// In that case we will merge our records from the config with the API records
// But those API records don't ever get saved in the tf config
// This is on purpose because the Akamai API will inject several
// Default records to a given zone and we don't want those to show up
// In diffs or in acceptance tests
/*
func mergeConfigs(recordType string, records []interface{}, s *schema.Resource, d *schema.ResourceData) *schema.Set {
	recordsInStateFile, recordsInConfigFile := d.GetChange(recordType)
	recordsInAPI := schema.NewSet(
		schema.HashResource(s.Schema[recordType].Elem.(*schema.Resource)),
		records,
	)
	recordsInAPIButNotInStateFile := recordsInAPI.Difference(recordsInStateFile.(*schema.Set))
	mergedRecordsToBeSaved := recordsInConfigFile.(*schema.Set).Union(recordsInAPIButNotInStateFile)

	return mergedRecordsToBeSaved
}
*/

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDNSRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hostname := d.Id()

	// find the zone first
	log.Printf("[INFO] [Akamai DNS] Searching for zone Records [%s]", hostname)
	//zone, err := dnsv2.GetZone(hostname)
	//if err != nil {
	//		return nil, err
	//	}

	// assign each of the record sets to the resource data
	//marshalResourceData(d, zone)

	// Give terraform the ID
	//d.SetId(fmt.Sprintf("%s-%s-%s", zone.Token, zone.Zone.Name, hostname))

	return []*schema.ResourceData{d}, nil
}

func resourceDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {

	zone := d.Get("zone").(string)
	host := d.Get("name").(string)
	recordtype := d.Get("recordtype").(string)
	ttl := d.Get("ttl").(int)
	//target :=  d.Get("target").(string)

	target := d.Get("target").(*schema.Set).List()

	records := make([]string, 0, len(target))
	for _, recContent := range target {
		records = append(records, recContent.(string))
	}

	recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}

	recordcreate.Delete(zone)

	d.SetId("")

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

	sha1_hash := getSHA(recordcreate.Target)
	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum for Existing SHA test %s", sha1_hash)

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

	log.Printf("[DEBUG] [Akamai DNSv2] record data read JSON [%s]", string(b1))

	if len(targets) > 0 {
		sha1_hash_test := getSHA(targets)

		if sha1_hash_test == sha1_hash {
			log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists matches [%s] vs  [%s] [%s] [%s] [%s] ", sha1_hash_test, sha1_hash, zone, host, recordtype)
			return true, nil
		} else {
			log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists mismatch [%s] vs  [%s] [%s] [%s] [%s] ", sha1_hash_test, sha1_hash, zone, host, recordtype)
			return false, nil
		}
	} else {
		log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from recordExists msimatch no target retunred  [%s] [%s] [%s] ", zone, host, recordtype)
		return false, nil
	}
	//	zone, err := dns.GetZone(hostname)
	//	return zone != nil, err
	//return false, nil
}

func getSHA(rdata []string) string {
	h := sha1.New()
	bodyBytes := new(bytes.Buffer)
	json.NewEncoder(bodyBytes).Encode(rdata)
	h.Write(bodyBytes.Bytes())

	sha1_hash_test := hex.EncodeToString(h.Sum(nil))
	log.Printf("[DEBUG] [Akamai DNSv2] SHA sum from Rdata %s %s", rdata, sha1_hash_test)
	return sha1_hash_test
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

	for _, recContent := range target {
		records = append(records, recContent.(string))
	}

	emptyrecordcreate := dnsv2.RecordBody{}

	simplerecord := map[string]bool{"A": true, "AAAA": true, "CNAME": true, "LOC": true, "NS": true, "PTR": true, "SPF": true, "TXT": true}
	if simplerecord[recordtype] {
		recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
		return recordcreate
	} else {
		if recordtype == "AFSDB" {

			records := make([]string, 0, len(target))
			subtype := d.Get("subtype").(int)
			for _, recContent := range target {
				records = append(records, strconv.Itoa(subtype)+" "+recContent.(string))
			}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "DNSKEY" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			protocol := d.Get("protocol").(int)
			algorithm := d.Get("algorithm").(int)
			key := d.Get("key").(string)
			//for _, recContent := range target {
			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(protocol)+" "+strconv.Itoa(algorithm)+" "+key)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "DS" {

			records := make([]string, 0, len(target))
			digest_type := d.Get("digest_type").(int)
			keytag := d.Get("keytag").(int)
			algorithm := d.Get("algorithm").(int)
			digest := d.Get("digest").(string)
			//for _, recContent := range target {
			records = append(records, strconv.Itoa(keytag)+" "+strconv.Itoa(digest_type)+" "+strconv.Itoa(algorithm)+" "+digest)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "HINFO" {

			records := make([]string, 0, len(target))
			hardware := d.Get("hardware").(string)
			software := d.Get("software").(string)
			//for _, recContent := range target {
			records = append(records, hardware+" "+software)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "LOC" {

			records := make([]string, 0, len(target))

			for _, recContent := range target {
				records = append(records, recContent.(string))
			}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "MX" {
			priority := d.Get("priority").(int)
			records := make([]string, 0, len(target))

			for _, recContent := range target {
				records = append(records, strconv.Itoa(priority)+" "+recContent.(string))
			}
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
			//for _, recContent := range target {
			records = append(records, strconv.Itoa(order)+" "+strconv.Itoa(preference)+" "+flagsnaptr+" "+regexp+" "+replacement+" "+service)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "NSEC3" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			next_hashed_owner_name := d.Get("next_hashed_owner_name").(string)
			salt := d.Get("salt").(string)
			type_bitmaps := d.Get("type_bitmaps").(string)
			//for _, recContent := range target {
			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(iterations)+" "+salt+" "+next_hashed_owner_name+" "+type_bitmaps)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "NSEC3PARAM" {

			records := make([]string, 0, len(target))
			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			salt := d.Get("salt").(string)
			saltbyte := []byte(salt)
			saltbase32 := base32.StdEncoding.EncodeToString(saltbyte)

			//for _, recContent := range target {
			records = append(records, strconv.Itoa(flags)+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(iterations)+" "+saltbase32)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "RP" {

			records := make([]string, 0, len(target))
			mailbox := d.Get("mailbox").(string)
			txt := d.Get("txt").(string)

			//for _, recContent := range target {
			records = append(records, mailbox+" "+txt)
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "RRSIG" { //TODO FIX

			records := make([]string, 0, len(target))
			expiration := d.Get("expiration").(string)
			inception := d.Get("inception").(string)
			original_ttl := d.Get("original_ttl").(int)
			algorithm := d.Get("algorithm").(int)
			labels := d.Get("labels").(int)
			keytag := d.Get("keytag").(int)
			signature := d.Get("signature").(string)
			signer := d.Get("signer").(string)
			type_covered := d.Get("type_covered").(string)
			//for _, recContent := range target {
			records = append(records, type_covered+" "+strconv.Itoa(algorithm)+" "+strconv.Itoa(labels)+" "+strconv.Itoa(original_ttl)+" "+expiration+" "+inception+" "+signature+" "+signer+" "+strconv.Itoa(keytag))
			//}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "SRV" {

			records := make([]string, 0, len(target))
			priority := d.Get("priority").(int)
			weight := d.Get("weight").(int)
			port := d.Get("port").(int)

			for _, recContent := range target {
				records = append(records, strconv.Itoa(weight)+" "+strconv.Itoa(port)+" "+strconv.Itoa(priority)+" "+recContent.(string))
			}
			recordcreate := dnsv2.RecordBody{Name: host, RecordType: recordtype, TTL: ttl, Target: records}
			return recordcreate
		}
		if recordtype == "SSHFP" {

			records := make([]string, 0, len(target))
			algorithm := d.Get("algorithm").(int)
			fingerprint_type := d.Get("fingerprint_type").(int)
			fingerprint := d.Get("fingerprint").(string)

			//for _, recContent := range target {
			records = append(records, strconv.Itoa(algorithm)+" "+strconv.Itoa(fingerprint_type)+" "+fingerprint)
			//+recContent.(string))
			//}
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

	//emptyrecordcreate := dnsv2.RecordBody{}

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
			digest_type := d.Get("digest_type").(int)
			keytag := d.Get("keytag").(int)
			algorithm := d.Get("algorithm").(int)
			digest := d.Get("digest").(string)
			//for _, recContent := range target {
			if digest_type == 0 {
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
			next_hashed_owner_name := d.Get("next_hashed_owner_name").(string)
			salt := d.Get("salt").(string)
			type_bitmaps := d.Get("type_bitmaps").(string)
			//for _, recContent := range target {

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
			if next_hashed_owner_name == "null" {
				return "next_hashed_owner_name"
			}
			if salt == "null" {
				return "salt"
			}
			if type_bitmaps == "null" {
				return "type_bitmaps"
			}
			return "VALID"
		}
		if recordtype == "NSEC3PARAM" {

			flags := d.Get("flags").(int)
			algorithm := d.Get("algorithm").(int)
			iterations := d.Get("iterations").(int)
			salt := d.Get("salt").(string)
			saltbyte := []byte(salt)
			saltbase32 := base32.StdEncoding.EncodeToString(saltbyte)

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
			if saltbyte == nil {
				return "saltbyte"
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
			original_ttl := d.Get("original_ttl").(int)
			algorithm := d.Get("algorithm").(int)
			labels := d.Get("labels").(int)
			keytag := d.Get("keytag").(int)
			signature := d.Get("signature").(string)
			signer := d.Get("signer").(string)
			type_covered := d.Get("type_covered").(string)
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
			if original_ttl == 0 {
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
			if type_covered == "null" {
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
			fingerprint_type := d.Get("fingerprint_type").(int)
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
			if fingerprint_type == 0 {
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

/*
type error interface {
	Error() string

}

type errorString struct {
	s string
}

func (e *errorString) Error () string {
	return e.s
}

func New(text string) error {
    return &errorString{text}
}
*/

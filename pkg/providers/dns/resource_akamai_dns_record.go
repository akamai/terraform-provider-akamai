package dns

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/logger"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Retry count for save, update and delete
const opRetryCount = 5

func resourceDNSv2Record() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		UpdateContext: resourceDNSRecordUpdate,
		DeleteContext: resourceDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDNSRecordImport,
		},
		Schema: getResourceDNSRecordSchema(),
	}
}

func getResourceDNSRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.NoZeroValues),
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
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
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
				RRTypeCaa,
				RRTypeCert,
				RRTypeTlsa,
				RRTypeSvcb,
				RRTypeHTTPS,
			}, false)),
		},
		"ttl": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"active": {
			Type:       schema.TypeBool,
			Optional:   true,
			Deprecated: "Field 'active' has been deprecated. Setting it has no effect",
		},
		"target": {
			Type:             schema.TypeList,
			Elem:             &schema.Schema{Type: schema.TypeString},
			Optional:         true,
			DiffSuppressFunc: dnsRecordTargetSuppress,
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
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldTrimQuoteSuppress,
		},
		"software": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldTrimQuoteSuppress,
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
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldTrimQuoteSuppress,
		},
		"service": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldTrimQuoteSuppress,
		},
		"regexp": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldTrimQuoteSuppress,
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
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldDotSuffixSuppress,
		},
		"txt": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldDotSuffixSuppress,
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
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: dnsRecordFieldDotSuffixSuppress,
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
		"usage": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"selector": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"match_type": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"certificate": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"type_value": {
			Type:             schema.TypeInt,
			Optional:         true,
			DiffSuppressFunc: dnsRecordTypeValueSuppress,
		},
		"type_mnemonic": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"record_sha": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"svc_priority": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"svc_params": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"target_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

/*
https://tools.ietf.org/html/rfc4398 section 2.1 defines cert types

             0            Reserved
             1  PKIX      X.509 as per PKIX
             2  SPKI      SPKI certificate
             3  PGP       OpenPGP packet
             4  IPKIX     The URL of an X.509 data object
             5  ISPKI     The URL of an SPKI certificate
             6  IPGP      The fingerprint and URL of an OpenPGP packet
             7  ACPKIX    Attribute Certificate
             8  IACPKIX   The URL of an Attribute Certificate
         9-252            Available for IANA assignment
           253  URI       URI private
           254  OID       OID private
*/

var certTypes = map[string]int{
	"PKIX":    1,
	"SPKI":    2,
	"PGP":     3,
	"IPKIX":   4,
	"ISPKI":   5,
	"IPGP":    6,
	"ACPKIX":  7,
	"IACPKIX": 8,
	"URI":     253,
	"OID":     254,
}

// Suppress check for fields that have dot suffix in tfstate
func dnsRecordFieldDotSuffixSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldValStr := strings.TrimRight(old, ".")
	newValStr := strings.TrimRight(new, ".")
	if oldValStr == newValStr {
		return true
	}
	return false
}

// Suppress check for fields that are quoted in tfstate
func dnsRecordFieldTrimQuoteSuppress(_, old, new string, _ *schema.ResourceData) bool {
	oldValStr := strings.Trim(old, "\\\"")
	newValStr := strings.Trim(new, "\"")
	if oldValStr == newValStr {
		return true
	}
	return false
}

// Suppress check for type_value. Mnemonic config comes back as numeric
func dnsRecordTypeValueSuppress(_, _, _ string, d *schema.ResourceData) bool {
	logger := logger.Get("[Akamai DNS]", "dnsRecordTypeValueSuppress")
	oldv, newv := d.GetChange("type_value")
	oldVal, ok := oldv.(int)
	if !ok {
		logger.Warnf("value is of invalid type: should be int: %v", oldv)
		return false
	}
	newVal, ok := newv.(int)
	if !ok {
		logger.Warnf("value is of invalid type: should be int: %v", newv)
		return false
	}
	mnemonicType, err := tf.GetStringValue("type_mnemonic", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Warnf("Error looking up 'type_mnemonic': %s", err)
		return false
	}
	mnemonicValue, ok := certTypes[mnemonicType]
	if !ok {
		return false
	}
	if oldVal == 0 && newVal == 0 {
		return true
	}
	if oldVal != 0 && newVal != 0 || newVal == 0 {
		return oldVal == mnemonicValue
	}
	if oldVal == 0 {
		return newVal == mnemonicValue
	}
	return false
}

// DiffSuppresFunc to handle quoted TXT Rdata strings possibly containing escaped quotes
func dnsRecordTargetSuppress(key, oldTarget, newTarget string, d *schema.ResourceData) bool {
	logger := logger.Get("[Akamai DNS]", "dnsRecordTargetSuppress")

	//new value (and length) for target is not known during plan
	if strings.HasSuffix(key, ".#") && newTarget == "" {
		return false
	}

	recordType, err := tf.GetStringValue("recordtype", d)
	if err != nil {
		logger.Warnf("fetching `recordtype`: %s", err)
		return false
	}
	oldlist, newlist := d.GetChange("target")
	oldTargetList, ok := oldlist.([]interface{})
	if !ok {
		logger.Warnf("value is of invalid type: should be []interface{}: %v", oldTargetList)
		return false
	}
	newTargetList, ok := newlist.([]interface{})
	if !ok {
		logger.Warnf("value is of invalid type: should be []interface{}: %v", newTargetList)
		return false
	}
	var oldStrList, newStrList []string
	for _, t := range oldTargetList {
		item, ok := t.(string)
		if !ok {
			logger.Warnf("value is of invalid type: should be []interface{}: %v", newTargetList)
			return false
		}
		oldStrList = append(oldStrList, item)
	}
	for _, t := range newTargetList {
		item, ok := t.(string)
		if !ok {
			logger.Warnf("value is of invalid type: should be []interface{}: %v", newTargetList)
			return false
		}
		newStrList = append(newStrList, item)
	}
	return diffQuotedDNSRecord(oldStrList, newStrList, oldTarget, newTarget, recordType, logger)
}

func diffQuotedDNSRecord(oldTargetList []string, newTargetList []string, old string, new string, recordType string, logger log.Interface) bool {
	const (
		singleQuote    = `"`
		backslashQuote = `\"`
	)
	if len(oldTargetList) != len(newTargetList) {
		return false
	}

	logger.Debugf("diffQuotedDNSRecord Suppress. recodtype: %v", recordType)
	logger.Debugf("diffQuotedDNSRecord Suppress. oldTargetList: [%v]", oldTargetList)
	logger.Debugf("diffQuotedDNSRecord Suppress. newTargetList: [%v]", newTargetList)
	logger.Debugf("diffQuotedDNSRecord Suppress. old: [%v]", old)
	logger.Debugf("diffQuotedDNSRecord Suppress. new: [%v]", new)

	var compList []string
	var baseVal string
	var compTrim bool
	if old == "" {
		baseVal = new
		compTrim = true
		baseVal = strings.Trim(baseVal, singleQuote)
		compList = oldTargetList
	} else {
		baseVal = old
		baseVal = strings.Trim(baseVal, backslashQuote)
		baseVal = strings.ReplaceAll(baseVal, backslashQuote, singleQuote)
		compList = newTargetList
	}

	// for AAAA record type, we want to compare IPv6 values
	if recordType == RRTypeAaaa {
		logger.Debugf("AAAA Suppress. baseval: [%v]", baseVal)
		fullBaseval := FullIPv6(net.ParseIP(baseVal))
		for _, compval := range compList {
			logger.Debugf("AAAA Suppress. compval: [%v]", compval)
			fullCompval := FullIPv6(net.ParseIP(compval))
			if fullBaseval == fullCompval {
				return true
			}
		}
		return false
	}

	if recordType == RRTypeCaa {
		baseVal = strings.ReplaceAll(baseVal, singleQuote, "")
		for _, compval := range compList {
			compval = strings.ReplaceAll(compval, singleQuote, "")
			logger.Debugf("updated baseVal: %v", baseVal)
			logger.Debugf("compval: %v", compval)
			if baseVal == compval {
				return true
			}
		}
		return false
	}

	if recordType == RRTypeMx {

		// lists are same length
		for i := 0; i < len(oldTargetList); i++ {
			oldTargetList[i] = strings.TrimRight(oldTargetList[i], ".")
			newTargetList[i] = strings.TrimRight(newTargetList[i], ".")
		}
		oldTargetString := strings.Join(oldTargetList, " ")
		newTargetString := strings.Join(newTargetList, " ")
		return oldTargetString == newTargetString
	}

	if recordType == RRTypeAfsdb || recordType == RRTypeCname || recordType == RRTypePtr || recordType == RRTypeSrv || recordType == RRTypeNs {
		baseVal = strings.TrimRight(baseVal, ".")
		for _, compval := range compList {
			compval = strings.TrimRight(compval, ".")
			logger.Debugf("updated baseVal: %v", baseVal)
			logger.Debugf("compval: %v", compval)
			if baseVal == compval {
				return true
			}
		}
		return false
	}

	for _, compval := range compList {
		if compTrim && strings.Contains(compval, backslashQuote) {
			compval = strings.ReplaceAll(compval, backslashQuote, singleQuote)
		}
		if baseVal == strings.Trim(compval, singleQuote) {
			return true
		}
	}
	return false
}

// Lock per record type
var recordCreateLock = map[string]*sync.Mutex{
	"A":          {},
	"AAAA":       {},
	"AFSDB":      {},
	"AKAMAICDN":  {},
	"AKAMAITLC":  {},
	"CAA":        {},
	"CERT":       {},
	"CNAME":      {},
	"HINFO":      {},
	"LOC":        {},
	"MX":         {},
	"NAPTR":      {},
	"NS":         {},
	"PTR":        {},
	"RP":         {},
	"SOA":        {},
	"SRV":        {},
	"SPF":        {},
	"SSHFP":      {},
	"TLSA":       {},
	"TXT":        {},
	"DNSKEY":     {},
	"DS":         {},
	"NSEC3":      {},
	"NSEC3PARAM": {},
	"RRSIG":      {},
	"SVCB":       {},
	"HTTPS":      {},
}

// Retrieves record lock per record type
func getRecordLock(recordType string) *sync.Mutex {
	return recordCreateLock[recordType]
}

func bumpSoaSerial(ctx context.Context, d *schema.ResourceData, meta meta.Meta, zone, host string, logger log.Interface) (*dns.RecordBody, error) {
	// Get SOA Record
	recordset, err := inst.Client(meta).GetRecord(ctx, zone, host, "SOA")
	if err != nil {
		return nil, fmt.Errorf("error looking up SOA record for %s: %w", host, err)
	}
	rdataFieldMap := inst.Client(meta).ParseRData(ctx, "SOA", recordset.Target)

	serial, ok := rdataFieldMap["serial"].(int)
	if !ok {
		return nil, fmt.Errorf("%w: %s, %q", tf.ErrInvalidType, "seral", "string")
	}
	if err := d.Set("serial", serial+1); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	newRecord, err := bindRecord(ctx, meta, d, logger)
	if err != nil {
		return nil, err
	}
	return &newRecord, nil
}

// Record op function
func execFunc(ctx context.Context, meta meta.Meta, fn string, rec *dns.RecordBody, zone string, rlock bool) error {

	var e error
	switch fn {
	case "Create":
		e = inst.Client(meta).CreateRecord(ctx, rec, zone, rlock)

	case "Update":
		e = inst.Client(meta).UpdateRecord(ctx, rec, zone, rlock)

	case "Delete":
		e = inst.Client(meta).DeleteRecord(ctx, rec, zone, rlock)

	default:
		e = fmt.Errorf("Invalid operation [%s]", fn)

	}
	return e
}

func executeRecordFunction(ctx context.Context, meta meta.Meta, name string, d *schema.ResourceData, fn string, rec *dns.RecordBody, zone, host, recordType string, logger log.Interface, rlock bool) error {

	logger.Debugf("executeRecordFunction - zone: %s, host: %s, recordtype: %s", zone, host, recordType)
	// DNS API can have Concurrency issues
	opRetry := opRetryCount
	e := execFunc(ctx, meta, fn, rec, zone, rlock)
	for e != nil && opRetry > 0 {
		apiError, ok := e.(*dns.Error)
		// prep failure or network failure?
		if !ok || apiError.StatusCode < http.StatusBadRequest {
			logger.Errorf("executeRecordFunction - %s Record failed for record [%s] [%s] [%s] ", name, zone, host, recordType)
			return e
		}
		if apiError.StatusCode == http.StatusConflict {
			logger.Debug("executeRecordFunction - Concurrency Conflict")
			opRetry--
			time.Sleep(100 * time.Millisecond)
			e = execFunc(ctx, meta, fn, rec, zone, rlock)
			continue
		}
		// relying on error string is not a good idea, better to introduce separate error variables for each cause or error codes
		if (name == "CREATE" || name == "UPDATE") && strings.Contains(e.Error(), "SOA serial number must be incremented") {
			logger.Debug("executeRecordFunction - SOA Serial Number needs incrementing")
			opRetry--
			time.Sleep(5 * time.Second) // let things quiesce
			rec, err := bumpSoaSerial(ctx, d, meta, zone, host, logger)
			if err != nil {
				return err
			}
			e = execFunc(ctx, meta, fn, rec, zone, rlock)
			continue
		}
		if name == "DELETE" && apiError.StatusCode == http.StatusNotFound {
			// record doesn't exist
			d.SetId("")
			logger.Debugf("executeRecordFunction - %s [WARNING] %s", name, "Record not found")
			break
		}
		logger.Debugf("executeRecordFunction - %s [ERROR] %s", name, e.Error())
		return e
	}
	return nil
}

// Create a new DNS Record
func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// only allow one record per record type to be created at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSRecordCreate")
	logger.Info("Record Create.")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	var zone, host, recordType string
	var err error
	var diags diag.Diagnostics

	zone, err = tf.GetStringValue("zone", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	host, err = tf.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	recordType, err = tf.GetStringValue("recordtype", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	logger.Infof("Record Create. zone: %s, host: %s, recordtype: %s", zone, host, recordType)

	if err := validateRecord(d); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("DNS record validation failure for recordset %s", host),
			Detail:   err.Error(),
		})
	}

	// serialize record creates of same type
	getRecordLock(recordType).Lock()
	defer getRecordLock(recordType).Unlock()

	if recordType == "SOA" {
		logger.Debug("Attempting to create a SOA record")
		// A default SOA is created automagically when the primary zone is created ...
		if _, err := inst.Client(meta).GetRecord(ctx, zone, host, recordType); err == nil {
			// Record exists
			serial, err := tf.GetIntValue("serial", d)
			if err != nil && !errors.Is(err, tf.ErrNotFound) {
				return diag.FromErr(err)
			}
			if err := d.Set("serial", serial+1); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
		} else {
			apiError, ok := err.(*dns.Error)
			if ok && apiError.StatusCode == http.StatusNotFound {
				logger.Debug("SOA Record not found. Initialize serial")
				if err := d.Set("serial", 1); err != nil {
					return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
				}
			}
		}
	}

	recordCreate, err := bindRecord(ctx, meta, d, logger)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Record bind failure",
			Detail:   err.Error(),
		})
	}

	logger.WithField("bind-object", recordCreate).Debug("Record Create")

	extractString := strings.Join(recordCreate.Target, " ")
	sha1hash := tools.GetSHAString(extractString)

	logger.Debugf("SHA sum for recordcreate [%s]", sha1hash)
	// First try to get the zone from the API
	logger.Debugf("Searching for records [%s]", zone)
	rdata := make([]string, 0)
	recordSet, e := inst.Client(meta).GetRecord(ctx, zone, host, recordType)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if !ok || apiError.StatusCode != http.StatusNotFound {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("failed looking up %s records for %s", recordType, host),
				Detail:   e.Error(),
			})
		}
	}
	if recordSet != nil {
		rdata = inst.Client(meta).ProcessRdata(ctx, recordSet.Target, recordType)
	}
	// If there's no existing record we'll create a blank one
	if e != nil {
		// record not found/404 we will create a new
		logger.Debug("Creating new record")
		// Save the zone to the API
		e = executeRecordFunction(ctx, meta, "CREATE", d, "Create", &recordCreate, zone, host, recordType, logger, false)
		if e != nil {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Recordset create failure",
				Detail:   e.Error(),
			})
		}
	} else {
		logger.Debug("Updating record")
		if len(rdata) > 0 {
			e = executeRecordFunction(ctx, meta, "CREATE", d, "Update", &recordCreate, zone, host, recordType, logger, false)
			if e != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Recordset create failure",
					Detail:   e.Error(),
				})
			}
		} else {
			logger.Debug("Saving record")
			e = executeRecordFunction(ctx, meta, "CREATE", d, "Create", &recordCreate, zone, host, recordType, logger, false)
			if e != nil {
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Recordset save failure",
					Detail:   e.Error(),
				})
			}
		}
	}
	// save hash
	if err := d.Set("record_sha", sha1hash); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordType))
	} else {
		// Backwards compatibility
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordType, sha1hash))
	}
	// Lock won't be release til after Read ...
	return resourceDNSRecordRead(ctx, d, meta)

}

// nolint:gocyclo
// Update DNS Record
func resourceDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// only allow one record per record type to be updated at a time
	// this prevents lost data if you are using a counter/dynamic variables
	// in your config.tf which might overwrite each other

	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSRecordUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	var zone, host, recordType string
	var err error
	zone, err = tf.GetStringValue("zone", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	host, err = tf.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	recordType, err = tf.GetStringValue("recordtype", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	target, err := tf.GetListValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if errors.Is(err, tf.ErrNotFound) {
		records := make([]string, 0, len(target))
		for _, recContent := range target {
			rec, ok := recContent.(string)
			if !ok {
				return diag.Errorf("record is of invalid type; should be 'string'")
			}
			records = append(records, rec)
		}
		logger.WithField("records", records).Debug("Update Records")
	}

	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordType,
	}).Info("record Update")

	if err := validateRecord(d); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("DNS record validation failure for %s", host),
			Detail:   err.Error(),
		})
	}

	// serialize record updates of same type
	getRecordLock(recordType).Lock()
	defer getRecordLock(recordType).Unlock()

	if recordType == "SOA" {
		// need to get current serial and increment as part of update
		record, e := inst.Client(meta).GetRecord(ctx, zone, host, recordType)
		if e != nil {
			apiError, ok := e.(*dns.Error)
			if !ok || apiError.StatusCode != http.StatusNotFound {
				logger.Error(fmt.Sprintf("UPDATE Read [ERROR] %s", e.Error()))
				return diag.FromErr(e)
			}
			logger.Errorf("UPDATE Record Read. error looking up %s records for %q: %s", recordType, host, e.Error())
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Update Recordset read failure",
				Detail:   e.Error(),
			})
		}
		// Parse Rdata
		serial, ok := inst.Client(meta).ParseRData(ctx, recordType, record.Target)["serial"].(int)
		if !ok {
			return diag.Errorf("%v: %s, %q", tf.ErrInvalidType, "seral", "string")
		}
		if err := d.Set("serial", serial+1); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
	}

	recordCreate, err := bindRecord(ctx, meta, d, logger)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Recordset update bind failure",
			Detail:   err.Error(),
		})
	}
	extractString := strings.Join(recordCreate.Target, " ")
	sha1hash := tools.GetSHAString(extractString)

	logger.Debugf("UPDATE SHA sum for recordupdate [%s]", sha1hash)
	// First try to get the zone from the API
	logger.Debugf("UPDATE Searching for records [%s]", zone)
	rdata := make([]string, 0, 0)
	recordset, e := inst.Client(meta).GetRecord(ctx, zone, host, recordType)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if !ok || apiError.StatusCode != http.StatusNotFound {
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Error looking up recordset %s", host),
				Detail:   e.Error(),
			})
		}
	}
	if recordset != nil {
		rdata = inst.Client(meta).ProcessRdata(ctx, recordset.Target, recordType)
	}
	logger.WithField("length", len(rdata)).Debug("UPDATE Searching for records")
	if len(rdata) == 0 {
		return resourceDNSRecordRead(ctx, d, meta)
	}
	extractString = strings.Join(rdata, " ")
	sha1hashtest := tools.GetSHAString(extractString)
	logger.Debugf("UPDATE SHA sum from recordread [%s]", sha1hashtest)
	sort.Strings(rdata)
	// If there's no existing record we'll create a blank one
	if e != nil {
		// if the record is not found/404 we will create a new
		logger.Errorf("UPDATE [ERROR] %s", e.Error())
		logger.Debugf("UPDATE Creating new record")
		// Save the zone to the API
		e = executeRecordFunction(ctx, meta, "UPDATE", d, "Create", &recordCreate, zone, host, recordType, logger, false)
		if e != nil {
			return diag.FromErr(e)
		}
	} else {
		logger.Debug("UPDATE Updating record")
		e = executeRecordFunction(ctx, meta, "UPDATE", d, "Update", &recordCreate, zone, host, recordType, logger, false)
		if e != nil {
			return diag.FromErr(e)
		}

	}
	// save hash
	if err := d.Set("record_sha", sha1hash); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	// Give terraform the ID
	if d.Id() == "" || strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordType))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordType, sha1hash))
	}
	// Lock not released until after Read ...
	return resourceDNSRecordRead(ctx, d, meta)
}

//nolint:gocyclo
func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSRecordRead")
	logger.Info("Record Read")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	var zone, host, recordType string
	var err error
	var diags diag.Diagnostics

	zone, err = tf.GetStringValue("zone", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	host, err = tf.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	recordType, err = tf.GetStringValue("recordtype", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordType,
	}).Info("Record Read")

	recordCreate, err := bindRecord(ctx, meta, d, logger)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Recordset bind failure",
			Detail:   err.Error(),
		})
	}
	b, err := json.Marshal(recordCreate.Target)
	if err != nil {
		logger.Errorf("Read Target Marshal Failure: %s", err.Error())
	}

	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordType,
	}).Debugf("READ record JSON from bind records: %s ", string(b))

	extractString := strings.Join(recordCreate.Target, " ")
	sha1hash := tools.GetSHAString(extractString)
	sort.Strings(recordCreate.Target)
	logger.Debugf("READ SHA sum for Existing SHA test %s %s", extractString, sha1hash)

	// try to get the zone from the API
	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordType,
	}).Info("READ Searching for zone records")

	record, e := inst.Client(meta).GetRecord(ctx, zone, host, recordType)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if !ok || apiError.StatusCode != http.StatusNotFound {
			logger.Errorf("RECORD READ. error looking up %s records for %q: %s", recordType, host, e.Error())
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Recordset read failure",
				Detail:   e.Error(),
			})
		}
	}
	if e != nil {
		// record doesn't exist
		logger.Errorf("READ Record Not Found: %s", e.Error())
		d.SetId("")
		return diag.Errorf("Record not found")
	}

	logger.Debugf("RECORD READ [%v] [%s] [%s] [%s] ", record, zone, host, recordType)

	b1, err := json.Marshal(record.Target)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Target marshal failure",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("READ record data read JSON %s", string(b1))
	rdataFieldMap := inst.Client(meta).ParseRData(ctx, recordType, record.Target) // returns map[string]interface{}
	targets := inst.Client(meta).ProcessRdata(ctx, record.Target, recordType)
	switch recordType {
	case RRTypeMx:
		// calc rdata sha from read record
		rdataString := strings.Join(record.Target, " ")
		recordSHA, err := tf.GetStringValue("record_sha", d)
		sort.Strings(recordCreate.Target)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		shaRdata := tools.GetSHAString(rdataString)
		if d.HasChange("target") {
			logger.Debug("MX READ. TARGET HAS CHANGED")
			// has remote changed independently of TF?
			if recordSHA != shaRdata {
				if len(recordSHA) > 0 {
					return diag.Errorf("Recordset [%s %s]: Remote has diverged from TF Config. Manual intervention required.", host, recordType)
				}
				logger.Debug("MX READ. record_sha ull. Refresh")
			} else {
				logger.Debug("MX READ. Remote static")
				d.SetId("")
			}
		} else {
			logger.Debug("MX READ. TARGET HAS NOT CHANGED")
			// has remote changed independently of TF?
			if recordSHA != shaRdata && len(recordSHA) > 0 {
				// another special case ... for instances record sha might not be representative of full resource
				target, err := tf.GetListValue("target", d)
				if err != nil && !errors.Is(err, tf.ErrNotFound) {
					return diag.FromErr(err)
				}
				if len(target) != 1 || sha1hash != shaRdata {
					return diag.Errorf("Recordset [%s %s]: Remote has diverged from TF Config. Manual intervention required.", host, recordType)
				}
			}
		}
	case RRTypeAaaa:
		sort.Strings(record.Target)
		rdataString := strings.Join(record.Target, " ")
		shaRdata := tools.GetSHAString(rdataString)
		if sha1hash == shaRdata {
			return nil
		}
		// could be either short or long notation
		newrdata := make([]string, 0, len(record.Target))
		for _, rcontent := range record.Target {
			newrdata = append(newrdata, rcontent)
		}
		if err := d.Set("target", newrdata); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
		targets = newrdata
	case RRTypeTxt:
		for fname, fvalue := range rdataFieldMap {
			if fvalue, ok := fvalue.([]string); ok {
				for i, v := range fvalue {
					fvalue[i] = txtRecordUnescape(v)
				}
				if err := d.Set(fname, fvalue); err != nil {
					return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
				}
			} else {
				return diag.Errorf("Invalid type conversion")
			}
		}

	default:
		// Parse Rdata. MX special
		for fname, fvalue := range rdataFieldMap {
			if err := d.Set(fname, fvalue); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
		}
	}

	if len(targets) == 0 {
		return diag.Errorf("[ERROR] [Akamai DNSv2] READ -  Invalid RData Returned for Recordset %s %s %s", zone, host, recordType)
	}

	sort.Strings(targets)
	if recordType == "SOA" {
		logger.Debug("READ SOA RECORD")
		rdataSerial, ok := rdataFieldMap["serial"].(int)
		if !ok {
			return diag.Errorf("'serial' is of invalid type; should be 'int'")
		}
		serial, err := tf.GetIntValue("serial", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		if rdataSerial >= serial {
			logger.Debug("READ SOA RECORD CHANGE: SOA OK")
			if ok := validateSOARecord(d, logger); ok {
				extractSoaString := strings.Join(targets, " ")
				sha1hash = tools.GetSHAString(extractSoaString)
				logger.Debug("READ SOA RECORD CHANGE: SOA OK")
			}
		}
	}
	if recordType == "AKAMAITLC" {
		extractTlcString := strings.Join(targets, " ")
		sha1hash = tools.GetSHAString(extractTlcString)
	}
	if err := d.Set("record_sha", sha1hash); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	// Give terraform the ID
	if strings.Contains(d.Id(), "#") {
		d.SetId(fmt.Sprintf("%s#%s#%s", zone, host, recordType))
	} else {
		d.SetId(fmt.Sprintf("%s-%s-%s-%s", zone, host, recordType, sha1hash))
	}
	return nil
}

func validateSOARecord(d *schema.ResourceData, logger log.Interface) bool {
	oldSer, newSer := d.GetChange("serial")
	newSerial, ok := newSer.(int)
	if !ok {
		logger.Warn("new serial is of invalid type; should be 'int'")
	}
	oldSerial, ok := oldSer.(int)
	if !ok {
		logger.Warn("old serial is of invalid type; should be 'int'")
	}
	if oldSerial > newSerial {
		return false
	}
	if d.HasChange("name_server") ||
		d.HasChange("email_address") ||
		d.HasChange("refresh") ||
		d.HasChange("retry") ||
		d.HasChange("expiry") ||
		d.HasChange("nxdomain_ttl") {
		return false
	}
	return true
}

func resourceDNSRecordImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSRecordImport")
	// create a context with logging for api calls

	// create context. TODO: *** Way to find TF context ***
	ctx := context.TODO()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	idParts := strings.Split(d.Id(), "#")
	if len(idParts) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid ID for Zone Import: %s", d.Id())
	}
	zone := idParts[0]
	recordName := idParts[1]
	recordType := idParts[2]

	logger.Info("Record Import.")

	// Get recordset
	logger.Debugf("Searching for zone Recordset. %s", idParts)

	recordset, e := inst.Client(meta).GetRecord(ctx, zone, recordName, recordType)
	if e != nil {
		apiError, ok := e.(*dns.Error)
		if !ok || apiError.StatusCode != http.StatusNotFound {
			logger.Debugf("IMPORT Record read failed for record [%s] [%s] [%s] ", zone, recordName, recordType)
			d.SetId("")
			return []*schema.ResourceData{d}, e
		}
		// record doesn't exist
		d.SetId("")
		logger.Error("IMPORT Error. Record not found")
		return nil, fmt.Errorf("record not found")
	}

	if err := d.Set("zone", zone); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("name", recordset.Name); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("recordtype", recordset.RecordType); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("ttl", recordset.TTL); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	targets := inst.Client(meta).ProcessRdata(ctx, recordset.Target, recordType)
	if recordset.RecordType == "MX" {
		// can't guarantee order of MX records. Forced to set pri, incr to 0 and targets as is
		if err := d.Set("target", targets); err != nil {
			return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		// Parse Rdata
		rdataFieldMap := inst.Client(meta).ParseRData(ctx, recordset.RecordType, recordset.Target) // returns map[string]interface{}
		for fname, fvalue := range rdataFieldMap {
			if err := d.Set(fname, fvalue); err != nil {
				return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
			}
		}
	}
	var importTargetString string
	if len(targets) == 0 {
		logger.Errorf("IMPORT Invalid Record. No target returned  [%s] [%s] [%s] ", zone, recordName, recordType)
		d.SetId("")
		return []*schema.ResourceData{d}, nil
	}
	if recordType != RRTypeMx {
		// MX Target Order important
		sort.Strings(targets)
	}
	importTargetString = strings.Join(targets, " ")
	sha1hash := tools.GetSHAString(importTargetString)
	if err := d.Set("record_sha", sha1hash); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(fmt.Sprintf("%s#%s#%s", zone, recordName, recordType))
	return []*schema.ResourceData{d}, nil
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "resourceDNSRecordUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	zone, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	host, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	recordType, err := tf.GetStringValue("recordtype", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ttl, err := tf.GetIntValue("ttl", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Infof("Record Delete. zone: %s, host: %s, recordtype: %s", zone, host, recordType)
	logger.Info("Record Delete.")
	// serialize record updates of same type
	getRecordLock(recordType).Lock()
	defer getRecordLock(recordType).Unlock()

	target, err := tf.GetListValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	records := make([]string, 0, len(target))
	for _, recContent := range target {
		recContentStr, ok := recContent.(string)
		if !ok {
			return diag.Errorf("record is of invalid type; should be 'string'")
		}
		records = append(records, recContentStr)
	}
	if recordType != RRTypeMx {
		sort.Strings(records)
	}
	logger.Debugf("Delete zone Record. Zone: %s, Host: %s, Recordtype:  %s", zone, host, recordType)
	recordcreate := dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	// Warning: Delete will expunge the ENTIRE Recordset regardless of whether user thought they were removing an instance

	if err := executeRecordFunction(ctx, meta, "DELETE", d, "Delete", &recordcreate, zone, host, recordType, logger, false); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// FullIPv6 encodes IPV6 as a full string
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

func padvalue(str string, logger log.Interface) string {
	vstr := strings.ReplaceAll(str, "m", "")
	logger.WithField("padvalue", str).Debug("[Akamai DNSv2]")
	vfloat, err := strconv.ParseFloat(vstr, 32)
	if err != nil {
		logger.Errorf("padvalue. Parse error: %s", vstr)
	}
	vresult := fmt.Sprintf("%.2f", vfloat)
	logger.Debugf("padvalue. Padded v_result %s", vresult)
	return vresult
}

// Used to pad coordinates to x.xxm format
func padCoordinates(str string, logger log.Interface) string {

	s := strings.Split(str, " ")
	if len(s) < 12 {
		logger.Debug("coordinates string is too short")
		return ""
	}
	latD, latM, latS, latDir, longD, longM, longS, longDir, altitude, size, horizPrecision, vertPrecision := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %sm %sm %sm %sm", latD, latM, latS, latDir, longD, longM, longS, longDir, padvalue(altitude, logger), padvalue(size, logger), padvalue(horizPrecision, logger), padvalue(vertPrecision, logger))
}

func bindRecord(ctx context.Context, meta meta.Meta, d *schema.ResourceData, logger log.Interface) (dns.RecordBody, error) {

	var host, recordType string
	var err error
	host, err = tf.GetStringValue("name", d)
	if err != nil {
		return dns.RecordBody{}, err
	}
	recordType, err = tf.GetStringValue("recordtype", d)
	if err != nil {
		return dns.RecordBody{}, err
	}
	ttl, err := tf.GetIntValue("ttl", d)
	if err != nil {
		return dns.RecordBody{}, err
	}

	target, err := tf.GetListValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return dns.RecordBody{}, err
	}
	records, err := buildRecordsList(target, recordType, logger)
	if err != nil {
		return dns.RecordBody{}, nil
	}

	simpleRecord := map[string]struct{}{"A": {}, "AAAA": {}, "AKAMAICDN": {}, "CNAME": {}, "LOC": {}, "NS": {}, "PTR": {}, "SPF": {}, "TXT": {}, "CAA": {}}
	if _, ok := simpleRecord[recordType]; ok {
		sort.Strings(records)
		return dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}, nil
	}
	return newRecordCreate(ctx, meta, d, recordType, target, host, ttl, logger)
}

//nolint:gocyclo
func newRecordCreate(ctx context.Context, meta meta.Meta, d *schema.ResourceData, recordType string, target []interface{}, host string, ttl int, logger log.Interface) (dns.RecordBody, error) {
	var recordCreate dns.RecordBody
	switch recordType {
	case RRTypeAfsdb:
		records := make([]string, 0, len(target))
		subtype, err := tf.GetIntValue("subtype", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		for _, recContent := range target {
			recContentStr, ok := recContent.(string)
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("record is of invalid type; should be 'string'")
			}
			record := strconv.Itoa(subtype) + " " + recContentStr
			if !strings.HasSuffix(recContentStr, ".") {
				record += "."
			}
			records = append(records, record)

		}
		sort.Strings(records)
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeDnskey:
		flags, err := tf.GetIntValue("flags", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		protocol, err := tf.GetIntValue("protocol", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		key, err := tf.GetStringValue("key", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(flags) + " " + strconv.Itoa(protocol) + " " + strconv.Itoa(algorithm) + " " + key}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeDs:
		digestType, err := tf.GetIntValue("digest_type", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		keytag, err := tf.GetIntValue("keytag", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		digest, err := tf.GetStringValue("digest", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(keytag) + " " + strconv.Itoa(algorithm) + " " + strconv.Itoa(digestType) + " " + digest}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeHinfo:
		hardware, err := tf.GetStringValue("hardware", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		software, err := tf.GetStringValue("software", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}

		// Fields may have embedded backslash. Quotes optional
		if strings.HasPrefix(hardware, `\"`) {
			hardware = strings.Trim(hardware, `\"`)
			hardware = `"` + hardware + `"`
		}
		if strings.HasPrefix(software, `\"`) {
			software = strings.Trim(software, `\"`)
			software = `"` + software + `"`
		}

		records := []string{hardware + " " + software}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeLoc:
		records := make([]string, 0, len(target))
		for _, recContent := range target {
			recContentStr, ok := recContent.(string)
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("record is of invalid type; should be 'string'")
			}
			records = append(records, recContentStr)
		}
		sort.Strings(records)
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeMx:
		zone, err := tf.GetStringValue("zone", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		logger.Debugf("MX record targets to process: %v", target)
		recordset, e := inst.Client(meta).GetRecord(ctx, zone, host, recordType)
		rdata := make([]string, 0, 0)
		if e != nil {
			logger.Debugf("MX Get Error Type: %T", e)
			logger.Debugf("BIND MX Error: %v", e)
			apiError, ok := e.(*dns.Error)
			if !ok || apiError.StatusCode != http.StatusNotFound {
				// failure other than not found
				return dns.RecordBody{}, fmt.Errorf(e.Error())
			}
			logger.Debug("Searching for existing MX records no prexisting targets found")
		} else {
			rdata = inst.Client(meta).ProcessRdata(ctx, recordset.Target, recordType)
		}
		logger.Debugf("Existing MX records to append to target %v", rdata)

		// keep track of rdata order
		rdataTarget := make([]string, 0, len(target)+len(rdata))
		//create map from rdata
		rdataTargetMap := make(map[string]int, len(rdata))
		for _, r := range rdata {
			entryparts := strings.Split(r, " ")
			if len(entryparts) < 2 {
				return dns.RecordBody{}, fmt.Errorf("RData shcould consist of at least 2 parts separated with ' '")
			}
			rn := entryparts[1]
			if !strings.HasSuffix(rn, ".") {
				rn += "."
			}
			// keep track of order for merge later
			rdataTarget = append(rdataTarget, rn)
			rdataTargetMap[rn], err = strconv.Atoi(entryparts[0])
			if err != nil {
				logger.Warnf("First part of RData string should be represented as integer: %s", entryparts[0])
			}
		}
		logger.Debugf("Created rdataTarget %v", rdataTarget)
		if d.HasChange("target") {
			// see if any entry was deleted. If so, remove from rdata map.
			oldList, newList := d.GetChange("target")
			oldTargetList, ok := oldList.([]interface{})
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("'oldList' is of invalid type; should be '[]interface{}'")
			}
			newTargetList, ok := newList.([]interface{})
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("'newList' is of invalid type; should be '[]interface{}'")
			}
			logger.Debugf("oldTargetList: %v", oldTargetList)
			logger.Debugf("newTargetList: %v", newTargetList)
			for _, oldTarg := range oldTargetList {
				oldTargStr, ok := oldTarg.(string)
				if !ok {
					return dns.RecordBody{}, fmt.Errorf("oldTarg is of invalid type; should be 'string'")
				}
				for _, newTarg := range newTargetList {
					newTargStr, ok := newTarg.(string)
					if !ok {
						return dns.RecordBody{}, fmt.Errorf("newTarg is of invalid type; should be 'string'")
					}
					if oldTargStr == newTargStr {
						// FIXME: this only breaks the inner loop, in which case this loop does nothing
						// probably a label should be added to the outer loop
						break
					}
				}
				// not there. remove
				logger.Debugf("MX BIND target %v deleted", oldTarg)
				delTarg := oldTargStr
				rdtParts := strings.Split(oldTargStr, " ")
				if len(rdtParts) > 1 {
					delTarg = rdtParts[1]
				}
				delete(rdataTargetMap, delTarg)
				logger.Debugf("Removing %v from rdataTarget %v", delTarg, rdataTarget)
				for i, item := range rdataTarget {
					if len(rdataTarget) > 0 && item == delTarg {
						copy(rdataTarget[i:], rdataTarget[i+1:])
						rdataTarget[len(rdataTarget)-1] = ""
						rdataTarget = rdataTarget[:len(rdataTarget)-1]
						logger.Debugf("Remove: UPDATED rdataTarget %v", rdataTarget)
					}
				}
			}
		}
		records := make([]string, 0, len(target)+len(rdata))

		priority, err := tf.GetIntValue("priority", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		increment, err := tf.GetIntValue("priority_increment", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		logger.Debugf("MX BIND Priority: %d ; Increment: %d", priority, increment)
		// walk thru target first
		var r int
		for _, recContent := range target {
			targEntry, ok := recContent.(string)
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("record is of invalid type; should be 'string'")
			}
			if !strings.HasSuffix(targEntry, ".") {
				targEntry += "."
			}
			logger.Debugf("MX BIND Processing Target %s", targEntry)
			targHost := targEntry
			var targPri int
			targParts := strings.Split(targEntry, " ") // need to support target entry with/without priority
			if len(targParts) > 2 {
				return dns.RecordBody{}, fmt.Errorf("Invalid MX Record format")
			}
			if len(targParts) == 2 {
				targHost = targParts[1]
				targPri, err = strconv.Atoi(targParts[0])
				if err != nil {
					return dns.RecordBody{}, fmt.Errorf("Invalid MX Record format")
				}
			} else {
				targPri = priority
			}
			pri, ok := rdataTargetMap[targHost]
			if ok {
				logger.Debugf("MX BIND. %s in existing map", targEntry)
				// target already in rdata
				if pri != targPri {
					return dns.RecordBody{}, fmt.Errorf("MX Record Priority Mismatch. Target order must align with EdgeDNS")
				}
			}
			// either match or we have inserted hosts in TF target
			if (r < len(rdataTarget) && rdataTarget[r] == targHost) || r >= len(rdataTarget) || !ok {
				if len(targParts) == 1 {
					records = append(records, strconv.Itoa(priority)+" "+targEntry)
				} else {
					records = append(records, targEntry)
				}
				if increment > 0 {
					priority += increment
				}
				if r < len(rdataTarget) && rdataTarget[r] == targHost {
					r++
					delete(rdataTargetMap, targHost)
				}
				continue
			}
			// mismatch. host in EdgeDns. Not at current target position
			logger.Debugf("Insert new target to records")
			// append what ever is left ...
			for {
				if (r >= len(rdataTarget)) || (rdataTarget[r] == targHost) {
					break
				}
				ntpri, _ := rdataTargetMap[rdataTarget[r]]
				records = append(records, strconv.Itoa(ntpri)+" "+rdataTarget[r])
				delete(rdataTargetMap, rdataTarget[r])
				r++
			}
			if len(targParts) == 1 {
				records = append(records, strconv.Itoa(priority)+" "+targEntry)
			} else {
				records = append(records, targEntry)
			}
			if increment > 0 {
				priority += increment
			}
			r++
			delete(rdataTargetMap, targHost)
		}
		logger.Debugf("Appended new target to target array LEN %d %v", len(records), records)
		// append what ever is left ...
		for {
			if r >= len(rdataTarget) {
				break
			}
			ntpri, _ := rdataTargetMap[rdataTarget[r]]
			logger.Debugf("Appending target %v pri %v", rdataTarget[r], ntpri)
			records = append(records, strconv.Itoa(ntpri)+" "+rdataTarget[r])
			r++
		}
		logger.Debugf("Existing MX records to append to target before schema data LEN %d %v", len(rdata), records)

		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeNaptr:
		flagsnaptr, err := tf.GetStringValue("flagsnaptr", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		order, err := tf.GetIntValue("order", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		preference, err := tf.GetIntValue("preference", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		regexp, err := tf.GetStringValue("regexp", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		replacement, err := tf.GetStringValue("replacement", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		// Following three fields may have embedded backslash
		service, err := tf.GetStringValue("service", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		if !strings.HasPrefix(service, `"`) {
			service = `"` + service + `"`
		}
		if !strings.HasPrefix(regexp, `"`) {
			regexp = `"` + regexp + `"`
		}
		if !strings.HasPrefix(flagsnaptr, `"`) {
			flagsnaptr = `"` + flagsnaptr + `"`
		}
		records := []string{strconv.Itoa(order) + " " + strconv.Itoa(preference) + " " + flagsnaptr + " " + service + " " + regexp + " " + replacement}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeNsec3:
		flags, err := tf.GetIntValue("flags", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		iterations, err := tf.GetIntValue("iterations", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		nextHashedOwnerName, err := tf.GetStringValue("next_hashed_owner_name", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		salt, err := tf.GetStringValue("salt", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		typeBitmaps, err := tf.GetStringValue("type_bitmaps", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(algorithm) + " " + strconv.Itoa(flags) + " " + strconv.Itoa(iterations) + " " + salt + " " + nextHashedOwnerName + " " + typeBitmaps}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeNsec3Param:
		flags, err := tf.GetIntValue("flags", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		iterations, err := tf.GetIntValue("iterations", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		salt, err := tf.GetStringValue("salt", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(algorithm) + " " + strconv.Itoa(flags) + " " + strconv.Itoa(iterations) + " " + salt}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeRp:
		mailbox, err := tf.GetStringValue("mailbox", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		if !strings.HasSuffix(mailbox, ".") {
			mailbox = mailbox + "."
		}
		txt, err := tf.GetStringValue("txt", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		if !strings.HasSuffix(txt, ".") {
			txt += "."
		}
		records := []string{mailbox + " " + txt}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeRrsig:
		expiration, err := tf.GetStringValue("expiration", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		inception, err := tf.GetStringValue("inception", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		originalTTL, err := tf.GetIntValue("original_ttl", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		labels, err := tf.GetIntValue("labels", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		keytag, err := tf.GetIntValue("keytag", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		signature, err := tf.GetStringValue("signature", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		signer, err := tf.GetStringValue("signer", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		typeCovered, err := tf.GetStringValue("type_covered", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{typeCovered + " " + strconv.Itoa(algorithm) + " " + strconv.Itoa(labels) + " " + strconv.Itoa(originalTTL) + " " + expiration + " " + inception + " " + strconv.Itoa(keytag) + " " + signer + " " + signature}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeSrv:
		records := make([]string, 0, len(target))
		priority, err := tf.GetIntValue("priority", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		weight, err := tf.GetIntValue("weight", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		port, err := tf.GetIntValue("port", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		for _, recContent := range target {
			recContentStr, ok := recContent.(string)
			if !ok {
				return dns.RecordBody{}, fmt.Errorf("record is of invalid type; should be 'string'")
			}
			record := strconv.Itoa(priority) + " " + strconv.Itoa(weight) + " " + strconv.Itoa(port) + " " + recContentStr
			if !strings.HasSuffix(recContentStr, ".") {
				record += "."
			}
			records = append(records, record)

		}
		sort.Strings(records)
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeSshfp:
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		fingerprintType, err := tf.GetIntValue("fingerprint_type", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		fingerprint, err := tf.GetStringValue("fingerprint", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(algorithm) + " " + strconv.Itoa(fingerprintType) + " " + fingerprint}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeSoa:
		nameserver, err := tf.GetStringValue("name_server", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		emailaddr, err := tf.GetStringValue("email_address", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		if !strings.HasSuffix(emailaddr, ".") {
			emailaddr += "."
		}
		serial, err := tf.GetIntValue("serial", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		refresh, err := tf.GetIntValue("refresh", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		retry, err := tf.GetIntValue("retry", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		expiry, err := tf.GetIntValue("expiry", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		nxdomainttl, err := tf.GetIntValue("nxdomain_ttl", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}

		records := []string{nameserver + " " + emailaddr + " " + strconv.Itoa(serial) + " " + strconv.Itoa(refresh) + " " + strconv.Itoa(retry) + " " + strconv.Itoa(expiry) + " " + strconv.Itoa(nxdomainttl)}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeAkamaiTlc:
		dnsname, err := tf.GetStringValue("dns_name", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		answtype, err := tf.GetStringValue("answer_type", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{answtype + " " + dnsname}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeCert:
		certtype, err := tf.GetStringValue("type_mnemonic", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		typevalue, err := tf.GetIntValue("type_value", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		keytag, err := tf.GetIntValue("keytag", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		algorithm, err := tf.GetIntValue("algorithm", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		certificate, err := tf.GetStringValue("certificate", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		// value or mnemonic type?
		if certtype == "" {
			certtype = strconv.Itoa(typevalue)
		}
		records := []string{certtype + " " + strconv.Itoa(keytag) + " " + strconv.Itoa(algorithm) + " " + certificate}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeTlsa:
		usage, err := tf.GetIntValue("usage", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		selector, err := tf.GetIntValue("selector", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		matchtype, err := tf.GetIntValue("match_type", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		certificate, err := tf.GetStringValue("certificate", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(usage) + " " + strconv.Itoa(selector) + " " + strconv.Itoa(matchtype) + " " + certificate}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	case RRTypeSvcb, RRTypeHTTPS:
		pri, err := tf.GetIntValue("svc_priority", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		tname, err := tf.GetStringValue("target_name", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		params, err := tf.GetStringValue("svc_params", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return dns.RecordBody{}, err
		}
		records := []string{strconv.Itoa(pri) + " " + tname + " " + params}
		recordCreate = dns.RecordBody{Name: host, RecordType: recordType, TTL: ttl, Target: records}

	default:
		return dns.RecordBody{}, fmt.Errorf("unable to create a Record Body for %s : %s", host, recordType)
	}

	return recordCreate, nil
}

func buildRecordsList(target []interface{}, recordType string, logger log.Interface) ([]string, error) {
	records := make([]string, 0, len(target))

	simpleRecordTarget := map[string]struct{}{"AAAA": {}, "CNAME": {}, "LOC": {}, "NS": {}, "PTR": {}, "SPF": {}, "SRV": {}, "TXT": {}, "CAA": {}}

	for _, recContent := range target {
		recContentStr, ok := recContent.(string)
		if !ok {
			return nil, fmt.Errorf("record is of invalid type; should be 'string'")
		}
		if _, ok := simpleRecordTarget[recordType]; !ok {
			records = append(records, recContentStr)
			continue
		}
		switch recordType {
		case RRTypeAaaa:
			addr := net.ParseIP(recContentStr)
			result := FullIPv6(addr)
			logger.Debugf("IPV6 full %s", result)
			records = append(records, result)
		case RRTypeLoc:
			logger.Debugf("LOC code format %s", recContentStr)
			str := padCoordinates(recContentStr, logger)
			records = append(records, str)
		case RRTypeSpf:
			if !strings.HasPrefix(recContentStr, "\"") {
				recContentStr = `"` + recContentStr + `"`
			}
			records = append(records, recContentStr)
		case RRTypeTxt:
			logger.Debugf("Bind TXT Data IN: [%s]", recContentStr)
			recContentStr = strings.Trim(recContentStr, `"`)
			recContentStr = txtRecordEscape(recContentStr)

			logger.Debugf("Bind TXT Data %s", recContentStr)
			logger.Debugf("Bind TXT Data OUT: [%s]", recContentStr)
			records = append(records, recContentStr)
		case RRTypeCaa:
			caaparts := strings.Split(recContentStr, " ")
			if len(caaparts) < 3 {
				return nil, fmt.Errorf("CAA record is of invalid format")
			}
			caaparts[2] = strings.Trim(caaparts[2], "\"")
			caaparts[2] = "\"" + caaparts[2] + "\""
			records = append(records, strings.Join(caaparts, " "))
		default:
			checktarget := recContentStr[len(recContentStr)-1:]
			if checktarget == "." {
				records = append(records, recContentStr)
			} else {
				records = append(records, recContentStr+".")
			}
		}
	}
	return records, nil
}

func validateRecord(d *schema.ResourceData) error {
	recordType, err := tf.GetStringValue("recordtype", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	switch recordType {
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
	case RRTypeCaa:
		return checkCaaRecord(d)
	case RRTypeCert:
		return checkCertRecord(d)
	case RRTypeTlsa:
		return checkTlsaRecord(d)
	case RRTypeSvcb:
		return checkSvcbRecord(d)
	case RRTypeHTTPS:
		return checkHTTPSRecord(d)
	default:
		return fmt.Errorf("invalid recordtype %v", recordType)
	}
}

func checkBasicRecordTypes(d *schema.ResourceData) error {
	_, err := tf.GetStringValue("name", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return err
		}
		return fmt.Errorf("configuration argument host must be set")
	}
	_, err = tf.GetStringValue("recordtype", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return err
		}
		return fmt.Errorf("configuration argument recordtype must be set")
	}
	_, err = tf.GetIntValue("ttl", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return err
		}
		return fmt.Errorf("configuration argument ttl must be set")
	}
	return nil
}

func checkTargets(d *schema.ResourceData) error {
	target, err := tf.GetListValue("target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if len(target) == 0 {
		return fmt.Errorf("configuration argument target must be set")
	}
	for _, recContent := range target {
		_, ok := recContent.(string)
		if !ok {
			return fmt.Errorf("target record is of invalid type; should be 'string'")
		}
	}
	return nil
}

func checkAsdfRecord(d *schema.ResourceData) error {
	subtype, err := tf.GetIntValue("subtype", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	if subtype == 0 {
		return fmt.Errorf("configuration argument subtype must be set for ASDF")
	}

	return checkTargets(d)
}

func checkDnskeyRecord(d *schema.ResourceData) error {
	flags, err := tf.GetIntValue("flags", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	protocol, err := tf.GetIntValue("protocol", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	key, err := tf.GetStringValue("key", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	ttl, err := tf.GetIntValue("ttl", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if !(flags == 0 || flags == 256 || flags == 257) {
		return fmt.Errorf("configuration argument flags must not be %v for DNSKEY", flags)
	}

	if ttl == 0 {
		return fmt.Errorf("configuration argument ttl must be set for DNSKEY")
	}

	if protocol == 0 {
		return fmt.Errorf("configuration argument protocol must be set for DNSKEY")
	}

	// FIXME this logic seems to be flawed, assertion will fail only if algorithm == 10
	if !((algorithm >= 1 && algorithm <= 8) || algorithm != 10) {
		return fmt.Errorf("configuration argument algorithm must not be %v for DNSKEY", algorithm)
	}

	if key == "" {
		return fmt.Errorf("configuration argument key must be set for DNSKEY")
	}

	return nil
}

func checkDsRecord(d *schema.ResourceData) error {
	digestType, err := tf.GetIntValue("digest_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	keytag, err := tf.GetIntValue("keytag", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	digest, err := tf.GetStringValue("digest", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if digestType == 0 {
		return fmt.Errorf("configuration argument digest_type must be set for DS")
	}

	if keytag == 0 {
		return fmt.Errorf("configuration argument keytag must be set for DS")
	}

	if algorithm == 0 {
		return fmt.Errorf("configuration argument algorithm must be set for DS")
	}

	if digest == "" {
		return fmt.Errorf("configuration argument digest must be set for DS")
	}

	return nil
}

func checkHinfoRecord(d *schema.ResourceData) error {
	hardware, err := tf.GetStringValue("hardware", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	software, err := tf.GetStringValue("software", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if hardware == "" {
		return fmt.Errorf("configuration argument hardware must be set for HINFO")
	}

	if software == "" {
		return fmt.Errorf("configuration argument software must be set for HINFO")
	}

	return nil
}

func checkMxRecord(d *schema.ResourceData) error {
	priority, err := tf.GetIntValue("priority", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if priority < 0 || priority > 65535 {
		return fmt.Errorf("configuration argument priority must be set for MX")
	}

	return checkTargets(d)
}

func checkNaptrRecord(d *schema.ResourceData) error {
	flagsnaptr, err := tf.GetStringValue("flagsnaptr", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	order, err := tf.GetIntValue("order", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	preference, err := tf.GetIntValue("preference", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	regexp, err := tf.GetStringValue("regexp", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	replacement, err := tf.GetStringValue("replacement", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	service, err := tf.GetStringValue("service", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if flagsnaptr == "" {
		return fmt.Errorf("configuration argument flagsnaptr must be set for NAPTR")
	}

	if order < 0 || order > 65535 {
		return fmt.Errorf("configuration argument order must not be %v for NAPTR", order)
	}

	if preference == 0 {
		return fmt.Errorf("configuration argument preference must be set for NAPTR")
	}

	if regexp == "" {
		return fmt.Errorf("configuration argument regexp must be set for NAPTR")
	}

	if replacement == "" {
		return fmt.Errorf("configuration argument replacement must be set for NAPTR")
	}

	if service == "" {
		return fmt.Errorf("configuration argument service must be set for NAPTR")
	}

	return nil
}

func checkNsec3Record(d *schema.ResourceData) error {
	flags, err := tf.GetIntValue("flags", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	iterations, err := tf.GetIntValue("iterations", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	nextHashedOwnerName, err := tf.GetStringValue("next_hashed_owner_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	salt, err := tf.GetStringValue("salt", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	typeBitmaps, err := tf.GetStringValue("type_bitmaps", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if !(flags == 0 || flags == 1) {
		return fmt.Errorf("configuration argument flags must be set for NSEC3")
	}

	if algorithm != 1 {
		return fmt.Errorf("configuration argument flags must be set for NSEC3")
	}
	if iterations == 0 {
		return fmt.Errorf("configuration argument iterations must be set for NSEC3")
	}
	if nextHashedOwnerName == "" {
		return fmt.Errorf("configuration argument nextHashedOwnerName must be set for NSEC3")
	}
	if salt == "" {
		return fmt.Errorf("configuration argument salt must be set for NSEC3")
	}
	if typeBitmaps == "" {
		return fmt.Errorf("configuration argument typeBitMaps must be set for NSEC3")
	}
	return nil
}

func checkNsec3ParamRecord(d *schema.ResourceData) error {
	flags, err := tf.GetIntValue("flags", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	iterations, err := tf.GetIntValue("iterations", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	salt, err := tf.GetStringValue("salt", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if !(flags == 0 || flags == 1) {
		return fmt.Errorf("configuration argument flags must be set for NSEC3PARAM")
	}

	if algorithm != 1 {
		return fmt.Errorf("configuration argument algorithm must be set for NSEC3PARAM")
	}

	if iterations == 0 {
		return fmt.Errorf("configuration argument iterations must be set for NSEC3PARAM")
	}

	if salt == "" {
		return fmt.Errorf("configuration argument salt must be set for NSEC3PARAM")
	}

	return nil
}

func checkRpRecord(d *schema.ResourceData) error {
	mailbox, err := tf.GetStringValue("mailbox", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	txt, err := tf.GetStringValue("txt", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if mailbox == "" {
		return fmt.Errorf("configuration argument mailbox must be set for RP")
	}

	if txt == "" {
		return fmt.Errorf("configuration argument txt must be set for RP")
	}

	return nil
}

func checkRrsigRecord(d *schema.ResourceData) error {
	expiration, err := tf.GetStringValue("expiration", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	inception, err := tf.GetStringValue("inception", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	originalTTL, err := tf.GetIntValue("original_ttl", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	labels, err := tf.GetIntValue("labels", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	keytag, err := tf.GetIntValue("keytag", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	signature, err := tf.GetStringValue("signature", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	signer, err := tf.GetStringValue("signer", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	typeCovered, err := tf.GetStringValue("type_covered", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if expiration == "" {
		return fmt.Errorf("configuration argument expiration must be set for RRSIG")
	}

	if inception == "" {
		return fmt.Errorf("configuration argument inception must be set for RRSIG")
	}

	if originalTTL == 0 {
		return fmt.Errorf("configuration argument originalTTL must be set for RRSIG")
	}

	if algorithm == 0 {
		return fmt.Errorf("configuration argument algorithm must be set for RRSIG")
	}

	if labels == 0 {
		return fmt.Errorf("configuration argument labels must be set for RRSIG")
	}

	if keytag == 0 {
		return fmt.Errorf("configuration argument keytag must be set for RRSIG")
	}

	if signature == "" {
		return fmt.Errorf("configuration argument signature must be set for RRSIG")
	}

	if signer == "" {
		return fmt.Errorf("configuration argument signer must be set for RRSIG")
	}

	if typeCovered == "" {
		return fmt.Errorf("configuration argument typeCovered must be set for RRSIG")
	}

	return nil
}

func checkSrvRecord(d *schema.ResourceData) error {
	priority, err := tf.GetIntValue("priority", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	weight, err := tf.GetIntValue("weight", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	port, err := tf.GetIntValue("port", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	if priority < 0 || priority > 65535 {
		return fmt.Errorf("configuration argument priority must be set for SRV")
	}

	if weight < 0 || weight > 65535 {
		return fmt.Errorf("configuration argument weight must not be %v for SRV", weight)
	}

	if port == 0 {
		return fmt.Errorf("configuration argument port must be set for SRV")
	}

	return nil
}

func checkSshfpRecord(d *schema.ResourceData) error {
	algorithm, err := tf.GetIntValue("algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	fingerprintType, err := tf.GetIntValue("fingerprint_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	fingerprint, err := tf.GetStringValue("fingerprint", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if algorithm == 0 {
		return fmt.Errorf("configuration argument algorithm must be set for SSHFP")
	}

	if fingerprintType == 0 {
		return fmt.Errorf("configuration argument fingerprintType must be set for SSHFP")
	}

	if fingerprint == "null" {
		return fmt.Errorf("configuration argument fingerprint must be set for SSHFP")
	}

	return nil
}

func checkSoaRecord(d *schema.ResourceData) error {

	nameserver, err := tf.GetStringValue("name_server", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	emailaddr, err := tf.GetStringValue("email_address", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	refresh, err := tf.GetIntValue("refresh", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	retry, err := tf.GetIntValue("retry", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	expiry, err := tf.GetIntValue("expiry", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	nxdomainttl, err := tf.GetIntValue("nxdomain_ttl", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if nameserver == "" {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "nameserver")
	}

	if emailaddr == "" {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "emailaddr")
	}

	if refresh == 0 {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "refresh")
	}

	if retry == 0 {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "retry")
	}

	if expiry == 0 {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "expiry")
	}

	if nxdomainttl == 0 {
		return fmt.Errorf("configuration argument %s must be specified in SOA record", "nxdomainttl")
	}

	return nil
}

func checkAkamaiTlcRecord(*schema.ResourceData) error {

	return fmt.Errorf("AKAMAITLC is a READ ONLY record")
}

func checkCaaRecord(d *schema.ResourceData) error {

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if err := checkTargets(d); err != nil {
		return err
	}

	caatarget := d.Get("target").([]interface{})
	for _, caa := range caatarget {
		caaStr, ok := caa.(string)
		if !ok {
			return fmt.Errorf("CAA is of invalid type; should be 'string'")
		}
		caaparts := strings.Split(caaStr, " ")
		if len(caaparts) != 3 {
			return fmt.Errorf("configuration argument CAA target %s is invalid", caaStr)
		}

		flag, err := strconv.Atoi(caaparts[0])
		if err != nil || flag < 0 || flag > 255 {
			return fmt.Errorf("configuration argument CAA target %s is invalid. flag value must be <= 0 and >= 255", caaStr)
		}
		re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
		submatchall := re.FindAllString(caaparts[1], -1)
		if len(submatchall) > 0 {
			return fmt.Errorf("configuration argument  CAA target %s is invalid. tag contains invalid characters", caaStr)
		}
	}

	return nil
}

func checkCertRecord(d *schema.ResourceData) error {
	typemnemonic, err := tf.GetStringValue("type_mnemonic", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	typevalue, err := tf.GetIntValue("type_value", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	certificate, err := tf.GetStringValue("certificate", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if typemnemonic == "" && typevalue == 0 {
		return fmt.Errorf("configuration arguments type_value and type_mnemonic are not set. Invalid CERT configuration")
	}

	if typemnemonic != "" && typevalue != 0 {
		return fmt.Errorf("configuration arguments type_value and type_mnemonic are both set. Invalid CERT configuration")
	}

	if certificate == "" {
		return fmt.Errorf("configuration argument certificate must be set for CERT")
	}
	return nil

}

func checkTlsaRecord(d *schema.ResourceData) error {

	usage, err := tf.GetIntValue("usage", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	certificate, err := tf.GetStringValue("certificate", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if certificate == "" {
		return fmt.Errorf("configuration argument certificate must be set for TLSA")
	}

	if usage == 0 {
		return fmt.Errorf("configuration argument usage must be set for TLSA")
	}

	return nil
}

func txtRecordEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return "\"" + s + "\""
}

func txtRecordUnescape(s string) string {
	s = s[1 : len(s)-1]
	s = strings.ReplaceAll(s, "\\\"", "\"")
	return strings.ReplaceAll(s, "\\\\", "\\")
}

func checkSvcbRecord(d *schema.ResourceData) error {

	return checkServiceRecord(d, "SVCB")
}

func checkHTTPSRecord(d *schema.ResourceData) error {

	return checkServiceRecord(d, "HTTPS")
}

func checkServiceRecord(d *schema.ResourceData, rtype string) error {

	pri, err := tf.GetIntValue("svc_priority", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	tname, err := tf.GetStringValue("target_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	params, err := tf.GetStringValue("svc_params", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if err := checkBasicRecordTypes(d); err != nil {
		return err
	}

	if pri < 0 || pri > 65535 {
		return fmt.Errorf("configuration argument svc_priority must be positive int for %s", rtype)
	}

	if tname == "" {
		return fmt.Errorf("configuration argument target_name must be set for %s", rtype)
	}

	if params == "" && pri > 0 {
		return fmt.Errorf("configuration argument svc_params must be set for %s", rtype)
	}

	if pri == 0 && params != "" {
		return fmt.Errorf("configuration argument svc_params cannot be set for %s if svc_priority is zero", rtype)
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
	RRTypeCert       = "CERT"
	RRTypeHTTPS      = "HTTPS"
	RRTypeSvcb       = "SVCB"
)

package gtm

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Datacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1DatacenterCreate,
		ReadContext:   resourceGTMv1DatacenterRead,
		UpdateContext: resourceGTMv1DatacenterUpdate,
		DeleteContext: resourceGTMv1DatacenterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1DatacenterImport,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wait_on_complete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clone_of": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cloud_server_host_header_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cloud_server_targeting": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_load_object": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"load_object": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"load_object_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"continent": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"ping_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_packet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"score_penalty": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_liveness_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_load_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_or_province": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// utility func to parse Terraform DC resource id
func parseDatacenterResourceID(id string) (string, int, error) {

	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" {
		return "", -1, fmt.Errorf("Datacenter ID, %v, is invalid", id)
	}

	domain := parts[0]
	dcID, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", -1, err
	}

	return domain, dcID, nil
}

var (
	datacenterCreateLock sync.Mutex
)

// Create a new GTM Datacenter
func resourceGTMv1DatacenterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DatacenterCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	// Async GTM DC creation causes issues at this writing. Synchronize as work around.
	datacenterCreateLock.Lock()
	defer datacenterCreateLock.Unlock()

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		logger.Errorf("Domain not initialized")
		return diag.FromErr(err)
	}
	datacenterName, err := tf.GetStringValue("nickname", d)
	if err != nil {
		logger.Errorf("nickname not initialized")
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	logger.Infof("Creating datacenter [%s] in domain [%s]", datacenterName, domain)
	newDC, err := populateNewDatacenterObject(ctx, meta, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Proposed New Datacenter: [%v]", newDC)
	cStatus, err := inst.Client(meta).CreateDatacenter(ctx, newDC, domain)
	if err != nil {
		logger.Errorf("Datacenter Create failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Create failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Datacenter Create status: %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		logger.Errorf(cStatus.Status.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  cStatus.Status.Message,
		})
	}
	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Datacenter Create completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter Create pending")
			} else {
				logger.Errorf("Datacenter Create failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Datacenter Create failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// Give terraform the ID. Format domain::dcid
	datacenterID := fmt.Sprintf("%s:%d", domain, cStatus.Resource.DatacenterId)
	logger.Debugf("Generated DC resource ID: %s", datacenterID)
	d.SetId(datacenterID)
	return resourceGTMv1DatacenterRead(ctx, d, m)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1DatacenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DatacenterRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading Datacenter: %s", d.Id())
	var diags diag.Diagnostics
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceID(d.Id())
	if err != nil {
		logger.Errorf("Invalid datacenter resource ID")
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Invalid Datacenter ID: %s", d.Id()),
			Detail:   err.Error(),
		})
	}
	dc, err := inst.Client(meta).GetDatacenter(ctx, dcID, domain)
	if err != nil {
		logger.Errorf("Datacenter Read failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformDCState(d, dc, m)
	logger.Debugf("READ %v", dc)
	return nil
}

// Update GTM Datacenter
func resourceGTMv1DatacenterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DatacenterUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Updating Datacenter: %s", d.Id())
	var diags diag.Diagnostics
	// pull domain and dcid out of resource id
	domain, dcID, err := parseDatacenterResourceID(d.Id())
	if err != nil {
		logger.Errorf("Invalid datacenter resource ID")
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Invalid Datacenter ID: %s", d.Id()),
			Detail:   err.Error(),
		})
	}
	// Get existing datacenter
	existDC, err := inst.Client(meta).GetDatacenter(ctx, dcID, domain)
	if err != nil {
		logger.Errorf("Datacenter Update failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Update Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Updating Datacenter BEFORE: %v", existDC)
	if err := populateDatacenterObject(d, existDC, m); err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Updating Datacenter PROPOSED: %v", existDC)
	uStat, err := inst.Client(meta).UpdateDatacenter(ctx, existDC, domain)
	if err != nil {
		logger.Errorf("Datacenter Update failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Datacenter Update status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Message)

	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Datacenter Update completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter Update pending")
			} else {
				logger.Errorf("Datacenter Update failed [%s]", err.Error())
				return diag.FromErr(fmt.Errorf("Datacenter Update failed [%s]", err.Error()))
			}
		}
	}

	return resourceGTMv1DatacenterRead(ctx, d, m)
}

func resourceGTMv1DatacenterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterImport")
	// create a context with logging for api calls
	ctx := context.Background()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Importing Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceID(d.Id())
	if err != nil {
		return nil, fmt.Errorf("Invalid Datacenter resource ID")
	}
	dc, err := inst.Client(meta).GetDatacenter(ctx, dcID, domain)
	if err != nil {
		logger.Errorf("Datacenter Import error: %s", err.Error())
		return nil, err
	}
	populateTerraformDCState(d, dc, m)
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	logger.Debugf("Import %v", dc)
	return []*schema.ResourceData{d}, err

}

// Delete GTM Datacenter.
func resourceGTMv1DatacenterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DatacenterDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Deleting Datacenter: %s", d.Id())
	var diags diag.Diagnostics
	domain, dcID, err := parseDatacenterResourceID(d.Id())
	if err != nil {
		logger.Errorf("Invalid Datacenter resource ID")
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Invalid Datacenter ID: %s", d.Id()),
			Detail:   err.Error(),
		})
	}
	// Get existing datacenter
	existDC, err := inst.Client(meta).GetDatacenter(ctx, dcID, domain)
	if err != nil {
		logger.Errorf("DatacenterDelete failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Delete error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Deleting Datacenter: %v", existDC)
	uStat, err := inst.Client(meta).DeleteDatacenter(ctx, existDC, domain)
	if err != nil {
		logger.Errorf("Datacenter Delete failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Datacenter Delete error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Datacenter Delete status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Datacenter Delete completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter Delete pending")
			} else {
				logger.Errorf("Datacenter Delete failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Datacenter Delete error",
					Detail:   err.Error(),
				})
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Create and populate a new datacenter object from resource data
func populateNewDatacenterObject(ctx context.Context, meta akamai.OperationMeta, d *schema.ResourceData, m interface{}) (*gtm.Datacenter, error) {

	dcObj := inst.Client(meta).NewDatacenter(ctx)
	dcObj.DefaultLoadObject = gtm.NewLoadObject()
	err := populateDatacenterObject(d, dcObj, m)

	return dcObj, err
}

// nolint:gocyclo
// Populate existing datacenter object from resource data
func populateDatacenterObject(d *schema.ResourceData, dc *gtm.Datacenter, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "populateDatacenterObject")

	vstr, err := tf.GetStringValue("nickname", d)
	if err == nil {
		dc.Nickname = vstr
	}
	vstr, err = tf.GetStringValue("city", d)
	if err == nil || d.HasChange("city") {
		dc.City = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() city failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	vint, err := tf.GetIntValue("clone_of", d)
	if err == nil || d.HasChange("clone_of") {
		dc.CloneOf = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() clone_of failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	cloudServerHostHeaderOverride, err := tf.GetBoolValue("cloud_server_host_header_override", d)
	if err != nil {
		logger.Errorf("populateDataCenterObject() failed: cloud_server_host_header_override not set: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}
	dc.CloudServerHostHeaderOverride = cloudServerHostHeaderOverride

	cloudServerTargeting, err := tf.GetBoolValue("cloud_server_targeting", d)
	if err != nil {
		logger.Errorf("cloud_server_targeting cloud_server_targeting not set: %s", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}
	dc.CloudServerTargeting = cloudServerTargeting

	vstr, err = tf.GetStringValue("continent", d)
	if err == nil || d.HasChange("continent") {
		dc.Continent = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() continent failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	vstr, err = tf.GetStringValue("country", d)
	if err == nil || d.HasChange("country") {
		dc.Country = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() country failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	// pull apart Set
	if dloList, err := tf.GetInterfaceArrayValue("default_load_object", d); err != nil || len(dloList) == 0 {
		dc.DefaultLoadObject = nil
	} else {
		dloObject := gtm.NewLoadObject()
		dloMap, ok := dloList[0].(map[string]interface{})
		if !ok {
			logger.Errorf("populateDatacenterObject default_load_object failed")
			return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
		}
		dloObject.LoadObject, ok = dloMap["load_object"].(string)
		if !ok {
			logger.Errorf("populateDatacenterObject load_object failed, bad load_object format")
			return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
		}
		dloObject.LoadObjectPort, ok = dloMap["load_object_port"].(int)
		if !ok {
			logger.Errorf("populateDatacenterObject failed, bad load_object_port format")
			return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
		}
		loadServers, ok := dloMap["load_servers"]
		if ok {
			servers, ok := loadServers.([]interface{})
			if ok {
				dloObject.LoadServers = make([]string, len(servers))
				for i, server := range servers {
					if dloObject.LoadServers[i], ok = server.(string); !ok {
						logger.Errorf("populateDatacenterObject failed, bad loadServer format: %s", server)
						return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
					}
				}
			} else {
				logger.Errorf("populateDatacenterObject failed, bad load_servers format: %s", loadServers)
				return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
			}
		} else {
			logger.Errorf("populateDatacenterObject failed, load_servers not present")
			return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
		}
		dc.DefaultLoadObject = dloObject
	}

	vfloat, err := tf.GetFloat64Value("latitude", d)
	if err == nil || d.HasChange("latitude") {
		dc.Latitude = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() latitude failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("longitude", d)
	if err == nil || d.HasChange("longitude") {
		dc.Longitude = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateDataCenterObject() longitude failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("ping_interval", d)
	if err == nil {
		dc.PingInterval = vint
	}
	vint, err = tf.GetIntValue("ping_packet_size", d)
	if err == nil {
		dc.PingPacketSize = vint
	}
	vint, err = tf.GetIntValue("datacenter_id", d)
	if err == nil {
		dc.DatacenterId = vint
	}
	vint, err = tf.GetIntValue("score_penalty", d)
	if err == nil {
		dc.ScorePenalty = vint
	}
	vint, err = tf.GetIntValue("servermonitor_liveness_count", d)
	if err == nil || d.HasChange("servermonitor_liveness_count") {
		dc.ServermonitorLivenessCount = vint
		if err != nil {
			logger.Warnf("populateDataCenterObject() failed: %v", err.Error())
		}
	}
	vint, err = tf.GetIntValue("servermonitor_load_count", d)
	if err == nil || d.HasChange("servermonitor_load_count") {
		dc.ServermonitorLoadCount = vint
		if err != nil {
			logger.Warnf("populateDataCenterObject() failed: %v", err.Error())
		}
	}
	vstr, err = tf.GetStringValue("servermonitor_pool", d)
	if err == nil || d.HasChange("servermonitor_pool") {
		dc.ServermonitorPool = vstr
		if err != nil {
			logger.Warnf("populateDataCenterObject() failed: %v", err.Error())
		}
	}
	vstr, err = tf.GetStringValue("state_or_province", d)
	if err == nil || d.HasChange("state_or_province") {
		dc.StateOrProvince = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Warnf("populateDataCenterObject() state_or_province failed: %v", err.Error())
		return fmt.Errorf("Datacenter Object could not be populated: %v", err.Error())
	}

	virtual, err := tf.GetBoolValue("virtual", d)
	dc.Virtual = virtual
	if err != nil {
		logger.Warnf("virtual not set: %s", err.Error())
	}

	return nil
}

// Populate Terraform state from provided Datacenter object
func populateTerraformDCState(d *schema.ResourceData, dc *gtm.Datacenter, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "populateTerrafomDCState")

	// walk through all state elements
	for stateKey, stateValue := range map[string]interface{}{
		"nickname":                          dc.Nickname,
		"datacenter_id":                     dc.DatacenterId,
		"city":                              dc.City,
		"clone_of":                          dc.CloneOf,
		"cloud_server_host_header_override": dc.CloudServerHostHeaderOverride,
		"cloud_server_targeting":            dc.CloudServerTargeting,
		"continent":                         dc.Continent,
		"country":                           dc.Country} {
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformDCState failed: %s", err.Error())
		}
	}
	dloStateList, err := tf.GetInterfaceArrayValue("default_load_object", d)
	if err != nil {
		dloStateList = make([]interface{}, 0, 1)
	}
	if len(dloStateList) == 0 && dc.DefaultLoadObject != nil && (dc.DefaultLoadObject.LoadObject != "" || len(dc.DefaultLoadObject.LoadServers) != 0 || dc.DefaultLoadObject.LoadObjectPort > 0) {
		// create MT object
		newDLO := make(map[string]interface{}, 3)
		newDLO["load_object"] = ""
		newDLO["load_object_port"] = 0
		newDLO["load_servers"] = make([]interface{}, 0, len(dc.DefaultLoadObject.LoadServers))
		dloStateList = append(dloStateList, newDLO)
	}
	for _, dloMap := range dloStateList {
		if dc.DefaultLoadObject != nil && (dc.DefaultLoadObject.LoadObject != "" || len(dc.DefaultLoadObject.LoadServers) != 0 || dc.DefaultLoadObject.LoadObjectPort > 0) {
			dlo := dloMap.(map[string]interface{})
			dlo["load_object"] = dc.DefaultLoadObject.LoadObject
			dlo["load_object_port"] = dc.DefaultLoadObject.LoadObjectPort
			if dlo["load_servers"] != nil && len(dlo["load_servers"].([]interface{})) > 0 {
				dlo["load_servers"] = reconcileTerraformLists(dlo["load_servers"].([]interface{}), convertStringToInterfaceList(dc.DefaultLoadObject.LoadServers, m), m)
			} else {
				dlo["load_servers"] = dc.DefaultLoadObject.LoadServers
			}
		} else {
			dloStateList = make([]interface{}, 0, 1)
		}
	}
	for stateKey, stateValue := range map[string]interface{}{
		"default_load_object":          dloStateList,
		"latitude":                     dc.Latitude,
		"longitude":                    dc.Longitude,
		"ping_interval":                dc.PingInterval,
		"ping_packet_size":             dc.PingPacketSize,
		"score_penalty":                dc.ScorePenalty,
		"servermonitor_liveness_count": dc.ServermonitorLivenessCount,
		"servermonitor_load_count":     dc.ServermonitorLoadCount,
		"servermonitor_pool":           dc.ServermonitorPool,
		"state_or_province":            dc.StateOrProvince,
		"virtual":                      dc.Virtual,
	} {
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformDCState failed: %s", err.Error())
		}
	}

}

package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetCreate,
		ReadContext:   resourceMatchTargetRead,
		UpdateContext: resourceMatchTargetUpdate,
		DeleteContext: resourceMatchTargetDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: resourceMatchTargetImport,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"match_target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentMatchTargetDiffs,
				Description:      "JSON-formatted definition of the match target",
			},
			"match_target_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the match target",
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")
	logger.Debugf("in resourceMatchTargetCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
	createMatchTarget := appsec.CreateMatchTargetRequest{}
	jsonpostpayload := d.Get("match_target")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createMatchTarget.ConfigID = configID
	createMatchTarget.ConfigVersion = version
	createMatchTarget.JsonPayloadRaw = rawJSON

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Errorf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createMatchTarget.ConfigID, postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRead")
	logger.Debugf("in resourceMatchTargetRead")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getMatchTarget := appsec.GetMatchTargetRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		TargetID:      targetID,
	}

	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Errorf("calling 'getMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}
	matchTargetConfigVal, err := tf.GetStringValue("match_target", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var response *appsec.GetMatchTargetResponse
	if err := json.Unmarshal([]byte(matchTargetConfigVal), &response); err != nil {
		return diag.FromErr(err)
	}

	if err := compareMatchTargetsOrder(matchtarget, response); err != nil {
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("match_target_id", matchtarget.TargetID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceMatchTargetImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetImport")
	logger.Debugf("in resourceMatchTargetImport")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return nil, err
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return nil, err
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return nil, err
	}
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return nil, err
	}

	getMatchTarget := appsec.GetMatchTargetRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		TargetID:      targetID,
	}

	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Errorf("calling 'getMatchTarget': %s", err.Error())
		return nil, err
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return nil, err
	}
	if err := d.Set("config_id", configID); err != nil {
		return nil, err
	}
	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return nil, err
	}
	if err := d.Set("match_target_id", matchtarget.TargetID); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil

}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")
	logger.Debugf("in resourceMatchTargetUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}
	jsonpostpayload := d.Get("match_target")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateMatchTarget := appsec.UpdateMatchTargetRequest{
		ConfigID:       configID,
		ConfigVersion:  version,
		TargetID:       targetID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Errorf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetDelete")
	logger.Debugf("in resourceMatchTargetDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	removeMatchTarget := appsec.RemoveMatchTargetRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		TargetID:      targetID,
	}

	_, err = client.RemoveMatchTarget(ctx, removeMatchTarget)
	if err != nil {
		logger.Errorf("calling 'removeMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

func compareMatchTargetsOrder(oldTarget, newTarget *appsec.GetMatchTargetResponse) error {

	oldJSONStr, err := json.Marshal(oldTarget)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	var oldJSON appsec.GetMatchTargetResponse
	if err = json.Unmarshal(oldJSONStr, &oldJSON); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	newJSONStr, err := json.Marshal(newTarget)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	var newJSON appsec.GetMatchTargetResponse
	if err = json.Unmarshal(newJSONStr, &newJSON); err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	sort.Strings(oldJSON.FilePaths)
	sort.Strings(newJSON.FilePaths)
	if reflect.DeepEqual(oldJSON.FilePaths, newJSON.FilePaths) {
		oldTarget.FilePaths = newTarget.FilePaths
	}

	sort.Strings(oldJSON.FileExtensions)
	sort.Strings(newJSON.FileExtensions)
	if reflect.DeepEqual(oldJSON.FileExtensions, newJSON.FileExtensions) {
		oldTarget.FileExtensions = newTarget.FileExtensions
	}

	sort.Strings(oldJSON.Hostnames)
	sort.Strings(newJSON.Hostnames)
	if reflect.DeepEqual(oldJSON.Hostnames, newJSON.Hostnames) {
		oldTarget.Hostnames = newTarget.Hostnames
	}

	sort.Slice(oldJSON.Apis, func(i, j int) bool {
		p1 := oldJSON.Apis[i]
		p2 := oldJSON.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})
	sort.Slice(newJSON.Apis, func(i, j int) bool {
		p1 := newJSON.Apis[i]
		p2 := newJSON.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})
	if reflect.DeepEqual(oldJSON.Apis, newJSON.Apis) {
		oldTarget.Apis = newTarget.Apis
	}

	sort.Slice(oldJSON.BypassNetworkLists, func(i, j int) bool {
		p1 := oldJSON.BypassNetworkLists[i]
		p2 := oldJSON.BypassNetworkLists[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})
	sort.Slice(newJSON.BypassNetworkLists, func(i, j int) bool {
		p1 := newJSON.BypassNetworkLists[i]
		p2 := newJSON.BypassNetworkLists[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})
	if reflect.DeepEqual(oldJSON.BypassNetworkLists, newJSON.BypassNetworkLists) {
		oldTarget.BypassNetworkLists = newTarget.BypassNetworkLists
	}
	return nil
}

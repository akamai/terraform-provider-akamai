package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetCreate,
		ReadContext:   resourceMatchTargetRead,
		UpdateContext: resourceMatchTargetUpdate,
		DeleteContext: resourceMatchTargetDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"json": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
			},
			"target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")

	createMatchTarget := v2.CreateMatchTargetRequest{}

	jsonpostpayload := d.Get("json")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &createMatchTarget)

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Warnf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(postresp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("json", string(jsonBody))

	d.Set("target_id", postresp.TargetID)

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")

	updateMatchTarget := v2.UpdateMatchTargetRequest{}

	jsonpostpayload := d.Get("json")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateMatchTarget)
	updateMatchTarget.TargetID, _ = strconv.Atoi(d.Id())
	jsonBody, err := json.Marshal(updateMatchTarget)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("json", string(jsonBody))

	resp, err := client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Warnf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}
	jsonBody, err = json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("json", string(jsonBody))
	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRemove")

	removeMatchTarget := v2.RemoveMatchTargetRequest{}

	removeMatchTarget.ConfigID = d.Get("config_id").(int)
	removeMatchTarget.ConfigVersion = d.Get("version").(int)
	removeMatchTarget.TargetID, _ = strconv.Atoi(d.Id())

	_, err := client.RemoveMatchTarget(ctx, removeMatchTarget)
	if err != nil {
		logger.Warnf("calling 'removeMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceMatchTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRead")

	getMatchTarget := v2.GetMatchTargetRequest{}

	getMatchTarget.ConfigID = d.Get("config_id").(int)
	getMatchTarget.ConfigVersion = d.Get("version").(int)
	getMatchTarget.TargetID, _ = strconv.Atoi(d.Id())

	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Warnf("calling 'getMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("json", string(jsonBody))

	d.Set("target_id", matchtarget.TargetID)
	d.SetId(strconv.Itoa(matchtarget.TargetID))

	return nil
}

func matchTargetAsJSONDString(d *schema.ResourceData) (string, error) {

	updateMatchTarget := v2.UpdateMatchTargetRequest{}
	updateMatchTarget.ConfigID = d.Get("config_id").(int)
	updateMatchTarget.ConfigVersion = d.Get("version").(int)
	updateMatchTarget.TargetID, _ = strconv.Atoi(d.Id())
	updateMatchTarget.Type = d.Get("type").(string)
	updateMatchTarget.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
	updateMatchTarget.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
	updateMatchTarget.DefaultFile = d.Get("default_file").(string)
	updateMatchTarget.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
	updateMatchTarget.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
	updateMatchTarget.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
	updateMatchTarget.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
	bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

	for _, b := range bypassnetworklists {
		bl := v2.BypassNetworkList{}
		bl.ID = b.(string)
		updateMatchTarget.BypassNetworkLists = append(updateMatchTarget.BypassNetworkLists, bl)
	}

	jsonBody, err := json.Marshal(updateMatchTarget)
	if err != nil {
		return "", err
	}
	return string(jsonBody), nil

}

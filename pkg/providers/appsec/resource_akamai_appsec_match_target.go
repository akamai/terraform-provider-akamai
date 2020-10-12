package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"type", "is_negative_path_match", "is_negative_file_extension_match", "default_file", "hostnames", "file_paths", "file_extensions", "security_policy", "bypass_network_lists"},
			},
			"target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"is_negative_path_match": {
				Type:             schema.TypeBool,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"is_negative_file_extension_match": {
				Type:             schema.TypeBool,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"default_file": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"hostnames": &schema.Schema{
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"file_paths": &schema.Schema{
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"file_extensions": &schema.Schema{
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"security_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
			"bypass_network_lists": &schema.Schema{
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				ConflictsWith:    []string{"json"},
				DiffSuppressFunc: suppressJsonProvided,
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")

	createMatchTarget := v2.CreateMatchTargetRequest{}

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &createMatchTarget)
	} else {
		createMatchTarget.ConfigID = d.Get("config_id").(int)
		createMatchTarget.ConfigVersion = d.Get("version").(int)
		createMatchTarget.Type = d.Get("type").(string)
		createMatchTarget.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
		createMatchTarget.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
		createMatchTarget.DefaultFile = d.Get("default_file").(string)
		createMatchTarget.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
		createMatchTarget.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
		createMatchTarget.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
		createMatchTarget.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
		bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

		for _, b := range bypassnetworklists {
			bl := appsec.BypassNetworkList{}
			bl.ID = b.(string)
			createMatchTarget.BypassNetworkLists = append(createMatchTarget.BypassNetworkLists, bl)
		}
	}

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Warnf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")

	updateMatchTarget := v2.UpdateMatchTargetRequest{}

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &updateMatchTarget)
		updateMatchTarget.TargetID, _ = strconv.Atoi(d.Id())
	} else {
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
			bl := appsec.BypassNetworkList{}
			bl.ID = b.(string)
			updateMatchTarget.BypassNetworkLists = append(updateMatchTarget.BypassNetworkLists, bl)
		}
	}

	_, err := client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Warnf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

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

	d.Set("type", matchtarget.Type)
	d.Set("is_negative_path_match", matchtarget.IsNegativePathMatch)
	d.Set("is_negative_file_extension_match", matchtarget.IsNegativeFileExtensionMatch)
	d.Set("default_file", matchtarget.DefaultFile)
	d.Set("hostnames", matchtarget.Hostnames)
	d.Set("file_paths", matchtarget.FilePaths)
	d.Set("file_extensions", matchtarget.FileExtensions)
	d.Set("security_policy", matchtarget.SecurityPolicy.PolicyID)
	d.Set("target_id", matchtarget.TargetID)
	d.SetId(strconv.Itoa(matchtarget.TargetID))

	return nil
}

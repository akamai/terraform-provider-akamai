package appsec

import (
	"encoding/json"
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceMatchTargetCreate,
		Read:   resourceMatchTargetRead,
		Update: resourceMatchTargetUpdate,
		Delete: resourceMatchTargetDelete,
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

func resourceMatchTargetCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating MatchTarget")

	matchtarget := appsec.NewMatchTargetResponse()

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &matchtarget)
	} else {
		matchtarget.ConfigID = d.Get("config_id").(int)
		matchtarget.ConfigVersion = d.Get("version").(int)
		matchtarget.Type = d.Get("type").(string)
		matchtarget.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
		matchtarget.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
		matchtarget.DefaultFile = d.Get("default_file").(string)
		matchtarget.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
		matchtarget.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
		matchtarget.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
		matchtarget.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
		bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

		for _, b := range bypassnetworklists {
			bl := appsec.BypassNetworkList{}
			bl.ID = b.(string)
			matchtarget.BypassNetworkLists = append(matchtarget.BypassNetworkLists, bl)
		}
	}

	postresp, err := matchtarget.SaveMatchTarget(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetRead(d, meta)
}

func resourceMatchTargetUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating MatchTarget")

	matchtarget := appsec.NewMatchTargetResponse()

	jsonpostpayload, ok := d.GetOk("json")
	if ok {

		json.Unmarshal([]byte(jsonpostpayload.(string)), &matchtarget)
		matchtarget.TargetID, _ = strconv.Atoi(d.Id())
	} else {
		matchtarget.ConfigID = d.Get("config_id").(int)
		matchtarget.ConfigVersion = d.Get("version").(int)
		matchtarget.TargetID, _ = strconv.Atoi(d.Id())
		matchtarget.Type = d.Get("type").(string)
		matchtarget.IsNegativePathMatch = d.Get("is_negative_path_match").(bool)
		matchtarget.IsNegativeFileExtensionMatch = d.Get("is_negative_file_extension_match").(bool)
		matchtarget.DefaultFile = d.Get("default_file").(string)
		matchtarget.Hostnames = tools.SetToStringSlice(d.Get("hostnames").(*schema.Set))
		matchtarget.FilePaths = tools.SetToStringSlice(d.Get("file_paths").(*schema.Set))
		matchtarget.FileExtensions = tools.SetToStringSlice(d.Get("file_extensions").(*schema.Set))
		matchtarget.SecurityPolicy.PolicyID = d.Get("security_policy").(string)
		bypassnetworklists := d.Get("bypass_network_lists").(*schema.Set).List()

		for _, b := range bypassnetworklists {
			bl := appsec.BypassNetworkList{}
			bl.ID = b.(string)
			matchtarget.BypassNetworkLists = append(matchtarget.BypassNetworkLists, bl)
		}
	}

	err := matchtarget.UpdateMatchTarget(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	return resourceMatchTargetRead(d, meta)
}

func resourceMatchTargetDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting MatchTarget")

	matchtarget := appsec.NewMatchTargetResponse()

	matchtarget.ConfigID = d.Get("config_id").(int)
	matchtarget.ConfigVersion = d.Get("version").(int)
	matchtarget.TargetID, _ = strconv.Atoi(d.Id())

	err := matchtarget.DeleteMatchTarget(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId("")

	return nil
}

func resourceMatchTargetRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceMatchTargetRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read MatchTarget")

	matchtarget := appsec.NewMatchTargetResponse()

	matchtarget.ConfigID = d.Get("config_id").(int)
	matchtarget.ConfigVersion = d.Get("version").(int)
	matchtarget.TargetID, _ = strconv.Atoi(d.Id())

	err := matchtarget.GetMatchTarget(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", matchtarget.TargetID))
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

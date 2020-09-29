package appsec

import (
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
func resourceConfigurationClone() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigurationCloneCreate,
		Read:   resourceConfigurationCloneRead,
		Update: resourceConfigurationCloneUpdate,
		Delete: resourceConfigurationCloneDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"create_from_version": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"rule_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of cloned configuration",
			},
		},
	}
}

func resourceConfigurationCloneCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationCloneCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating ConfigurationClone")

	configurationclone := appsec.NewConfigurationCloneResponse()

	configurationclonepost := appsec.NewConfigurationClonePost()

	configurationclone.ConfigID = d.Get("config_id").(int)
	configurationclonepost.CreateFromVersion = d.Get("create_from_version").(int)

	ccr, err := configurationclone.SaveConfigurationClone(configurationclonepost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.Set("version", ccr.Version)
	d.SetId(strconv.Itoa(ccr.Version))

	return resourceConfigurationCloneRead(d, meta)
}

func resourceConfigurationCloneRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationCloneRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read ConfigurationClone")

	configurationclone := appsec.NewConfigurationCloneResponse()

	configurationclone.ConfigID = d.Get("config_id").(int)
	configurationclone.Version = d.Get("create_from_version").(int)

	err := configurationclone.GetConfigurationClone(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(configurationclone.ConfigID))

	return nil
}

func resourceConfigurationCloneDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationCloneDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting ConfigurationClone")

	return schema.Noop(d, meta)
}

func resourceConfigurationCloneUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationCloneUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating ConfigurationClone")

	return schema.Noop(d, meta)
}

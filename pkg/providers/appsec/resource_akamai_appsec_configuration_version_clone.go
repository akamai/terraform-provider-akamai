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
func resourceConfigurationVersionClone() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigurationVersionCloneCreate,
		Read:   resourceConfigurationVersionCloneRead,
		Update: resourceConfigurationVersionCloneUpdate,
		Delete: resourceConfigurationVersionCloneDelete,
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

func resourceConfigurationVersionCloneCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationVersionCloneCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating ConfigurationClone")

	configurationversionclone := appsec.NewConfigurationVersionCloneResponse()

	ConfigurationVersionClonePost := appsec.NewConfigurationVersionClonePost()

	configurationversionclone.ConfigID = d.Get("config_id").(int)
	ConfigurationVersionClonePost.CreateFromVersion = d.Get("create_from_version").(int)

	ccr, err := configurationversionclone.SaveConfigurationClone(ConfigurationVersionClonePost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.Set("version", ccr.Version)
	d.SetId(strconv.Itoa(ccr.Version))

	return resourceConfigurationVersionCloneRead(d, meta)
}

func resourceConfigurationVersionCloneRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationVersionCloneRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read ConfigurationClone")

	configurationversionclone := appsec.NewConfigurationVersionCloneResponse()

	configurationversionclone.ConfigID = d.Get("config_id").(int)
	configurationversionclone.Version = d.Get("create_from_version").(int)

	err := configurationversionclone.GetConfigurationClone(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(configurationversionclone.ConfigID))

	return nil
}

func resourceConfigurationVersionCloneDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationVersionCloneDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting ConfigurationClone")

	return schema.Noop(d, meta)
}

func resourceConfigurationVersionCloneUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceConfigurationVersionCloneUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating ConfigurationClone")

	return schema.Noop(d, meta)
}

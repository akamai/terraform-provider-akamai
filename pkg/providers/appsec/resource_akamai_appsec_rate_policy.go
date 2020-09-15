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
func resourceRatePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceRatePolicyCreate,
		Read:   resourceRatePolicyRead,
		Update: resourceRatePolicyUpdate,
		Delete: resourceRatePolicyDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version_number": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"json": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRatePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating RatePolicy")

	ratepolicy := appsec.NewRatePolicyResponse()

	jsonpostpayload := d.Get("json")
	json.Unmarshal([]byte(jsonpostpayload.(string)), &ratepolicy)
	ratepolicy.ConfigID = d.Get("config_id").(int)
	ratepolicy.ConfigVersion = d.Get("version_number").(int)

	err := ratepolicy.SaveRatePolicy(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(ratepolicy.ID))

	return resourceRatePolicyRead(d, meta)
}

func resourceRatePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating RatePolicy")

	ratepolicy := appsec.NewRatePolicyResponse()

	jsonpostpayload := d.Get("json")
	json.Unmarshal([]byte(jsonpostpayload.(string)), &ratepolicy)
	ratepolicy.ConfigID = d.Get("config_id").(int)
	ratepolicy.ConfigVersion = d.Get("version_number").(int)
	ratepolicy.ID, _ = strconv.Atoi(d.Id())

	err := ratepolicy.UpdateRatePolicy(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	return resourceRatePolicyRead(d, meta)
}

func resourceRatePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting RatePolicy")

	ratepolicy := appsec.NewRatePolicyResponse()

	ratepolicy.ConfigID = d.Get("config_id").(int)
	ratepolicy.ConfigVersion = d.Get("version_number").(int)
	ratepolicy.ID, _ = strconv.Atoi(d.Id())

	err := ratepolicy.DeleteRatePolicy(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId("")

	return nil
}

func resourceRatePolicyRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read RatePolicy")

	ratepolicy := appsec.NewRatePolicyResponse()

	ratepolicy.ConfigID = d.Get("config_id").(int)
	ratepolicy.ConfigVersion = d.Get("version_number").(int)
	ratepolicy.ID, _ = strconv.Atoi(d.Id())

	err := ratepolicy.GetRatePolicy(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.SetId(strconv.Itoa(ratepolicy.ID))

	return nil
}

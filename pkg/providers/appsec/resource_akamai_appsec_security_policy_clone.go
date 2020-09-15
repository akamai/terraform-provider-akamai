package appsec

import (
	"fmt"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSecurityPolicyClone() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityPolicyCloneCreate,
		Read:   resourceSecurityPolicyCloneRead,
		Update: resourceSecurityPolicyCloneUpdate,
		Delete: resourceSecurityPolicyCloneDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"create_from_security_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID for clone",
			},
		},
	}
}

func resourceSecurityPolicyCloneCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSecurityPolicyCloneCreate-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating SecurityPolicyClone")

	securitypolicyclone := appsec.NewSecurityPolicyCloneResponse()
	securitypolicyclonepost := appsec.NewSecurityPolicyClonePost()
	securitypolicyclone.ConfigID = d.Get("config_id").(int)
	securitypolicyclone.Version = d.Get("version").(int)
	securitypolicyclonepost.CreateFromSecurityPolicy = d.Get("create_from_security_policy").(string)
	securitypolicyclonepost.PolicyName = d.Get("policy_name").(string)
	securitypolicyclonepost.PolicyPrefix = d.Get("policy_prefix").(string)
	spcr, err := securitypolicyclone.Save(securitypolicyclonepost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}
	d.Set("policy_id", spcr.PolicyID)
	d.Set("policy_name", spcr.PolicyName)
	d.Set("policy_prefix", securitypolicyclonepost.PolicyPrefix)
	d.SetId(spcr.PolicyID)
	return resourceSecurityPolicyCloneRead(d, meta)
}

func resourceSecurityPolicyCloneRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSecurityPolicyCloneRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SecurityPolicyClone")

	securitypolicyclone := appsec.NewSecurityPolicyCloneResponse()
	securitypolicyclone.ConfigID = d.Get("config_id").(int)
	securitypolicyclone.Version = d.Get("version").(int)
	policy_id := d.Id()
	spcr, err := securitypolicyclone.GetSecurityPolicyClone(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	for _, configval := range spcr.Policies {

		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", configval.PolicyID))
		if configval.PolicyID == policy_id {
			d.Set("policy_name", configval.PolicyName)
			d.Set("policy_id", configval.PolicyID)
			d.SetId(configval.PolicyID)
		}
	}
	return nil
}

func resourceSecurityPolicyCloneDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSecurityPolicyCloneDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting SecurityPolicyClone")
	return schema.Noop(d, meta)
}

func resourceSecurityPolicyCloneUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSecurityPolicyCloneUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating SecurityPolicyClone")
	return schema.Noop(d, meta)
}

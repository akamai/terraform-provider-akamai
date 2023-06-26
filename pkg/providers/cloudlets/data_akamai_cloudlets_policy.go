package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
)

func dataSourceCloudletsPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "An integer ID that is associated with a policy",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The version number of the policy",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Defines the group association for the policy. You must have edit privileges for the group",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the policy. The name must be unique",
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specific version of this API",
			},
			"cloudlet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Defines the policy type",
			},
			"cloudlet_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Code for the type of Cloudlet (ALB, AP, CD, ER, FR or VP)",
			},
			"revision_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID given to every policy version update",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of this specific policy",
			},
			"version_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of this specific version",
			},
			"rules_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, you cannot edit the match rules for the Cloudlet policy version",
			},
			"match_rules": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A JSON structure that defines the rules for this policy",
			},
			"match_rule_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the Cloudlet-specific matchRules",
			},
			"warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A JSON encoded list of warnings",
			},
			"activations": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of current policy activation information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The specific version of this API",
						},
						"network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network type, either 'staging' or 'prod' where a property or a Cloudlet policy has been activated",
						},
						"policy_info": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The object containing Cloudlet policy information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "An integer ID that is associated with all versions of a policy",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the policy",
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the activated policy",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The activation status for the policy: active, inactive, deactivated, pending or failed",
									},
									"status_detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Information about the status of an activation operation",
									},
									"activated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user who activated the policy",
									},
									"activation_date": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The date on which the policy was activated (in milliseconds since Epoch)",
									},
								},
							},
						},
						"property_info": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "A set containing information about the property associated with a particular Cloudlet policy",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the property",
									},
									"version": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the activated property",
									},
									"group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Defines the group association for the policy or property",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The activation status for the property. Can be active, inactive, deactivated, pending or failed",
									},
									"activated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the user who activated the property",
									},
									"activation_date": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The date on which the property was activated (in milliseconds since Epoch)",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getSchemaPolicyInfoFrom(p cloudlets.PolicyInfo) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"policy_id":       p.PolicyID,
			"name":            p.Name,
			"version":         p.Version,
			"status":          p.Status,
			"status_detail":   p.StatusDetail,
			"activated_by":    p.ActivatedBy,
			"activation_date": p.ActivationDate,
		},
	}
}

func getSchemaPropertyInfoFrom(p cloudlets.PropertyInfo) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":            p.Name,
			"version":         p.Version,
			"group_id":        p.GroupID,
			"status":          p.Status,
			"activated_by":    p.ActivatedBy,
			"activation_date": p.ActivationDate,
		},
	}
}

func getSchemaActivationsFrom(act []cloudlets.PolicyActivation) []map[string]interface{} {
	var activations []map[string]interface{}

	for _, a := range act {
		v := map[string]interface{}{
			"api_version":   a.APIVersion,
			"network":       a.Network,
			"policy_info":   getSchemaPolicyInfoFrom(a.PolicyInfo),
			"property_info": getSchemaPropertyInfoFrom(a.PropertyInfo),
		}
		activations = append(activations, v)
	}

	return activations
}

func populateSchemaFieldsWithPolicy(p *cloudlets.Policy, d *schema.ResourceData) error {
	fields := map[string]interface{}{
		"group_id":      p.GroupID,
		"name":          p.Name,
		"description":   p.Description,
		"cloudlet_id":   p.CloudletID,
		"cloudlet_code": p.CloudletCode,
		"api_version":   p.APIVersion,
	}

	err := tf.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func populateSchemaFieldsWithPolicyVersion(p *cloudlets.PolicyVersion, d *schema.ResourceData) error {
	matchRules, err := json.MarshalIndent(p.MatchRules, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal json %s", err)
	}
	warnings, err := json.MarshalIndent(p.Warnings, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal json %s", err)
	}

	fields := map[string]interface{}{
		"version":             p.Version,
		"version_description": p.Description,
		"revision_id":         p.RevisionID,
		"rules_locked":        p.RulesLocked,
		"match_rules":         string(matchRules),
		"match_rule_format":   p.MatchRuleFormat,
		"warnings":            string(warnings),
		"activations":         getSchemaActivationsFrom(p.Activations),
	}

	err = tf.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func dataSourceCloudletsPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("Cloudlets", "dataSourceCloudletsPolicyRead")
	client := inst.Client(meta)

	policyID, err := tf.GetIntValue("policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	var version int64
	if v, err := tf.GetIntValue("version", d); err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		version, err = findLatestPolicyVersion(ctx, int64(policyID), client)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		version = int64(v)
	}

	log.Debug("Getting Policy Version")
	policyVersion, err := client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: int64(policyID),
		Version:  version,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if policyVersion.Deleted {
		return diag.Errorf("specified policy version is deleted: version = %d", version)
	}

	err = populateSchemaFieldsWithPolicyVersion(policyVersion, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debug("Getting Policy")
	policy, err := client.GetPolicy(ctx, cloudlets.GetPolicyRequest{PolicyID: int64(policyID)})
	if err != nil {
		return diag.FromErr(err)
	}
	if policy.Deleted {
		return diag.Errorf("specified policy is deleted: policy_id = %d", policyID)
	}

	err = populateSchemaFieldsWithPolicy(policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", policyID, version))

	return nil
}

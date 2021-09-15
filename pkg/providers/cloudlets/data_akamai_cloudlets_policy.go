package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
)

func dataSourceCloudletsPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloudlet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloudlet_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revision_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules_locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"match_rules": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_rule_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"warnings": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"activations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_info": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status_detail": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"activated_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"activation_date": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"property_info": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"activated_by": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"activation_date": {
										Type:     schema.TypeInt,
										Computed: true,
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

func getSchemaActivationsFrom(act []*cloudlets.Activation) []map[string]interface{} {
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

	err := tools.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func findLatestPolicyVersion(ctx context.Context, client cloudlets.Cloudlets, policyID int) (int64, error) {
	versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
		PolicyID: int64(policyID),
	})
	if err != nil {
		return 0, err
	}

	if len(versions) == 0 {
		return 0, fmt.Errorf("latest policy version does not exist")
	}
	var latest int64
	for _, v := range versions {
		if v.Version > latest {
			latest = v.Version
		}
	}
	return latest, nil
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
		"version_description": p.Description,
		"revision_id":         p.RevisionID,
		"rules_locked":        p.RulesLocked,
		"match_rules":         string(matchRules),
		"match_rule_format":   p.MatchRuleFormat,
		"warnings":            string(warnings),
		"activations":         getSchemaActivationsFrom(p.Activations),
	}

	err = tools.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func dataSourceCloudletsPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("Cloudlets", "dataSourceCloudletsPolicyRead")
	client := inst.Client(meta)

	policyID, err := tools.GetIntValue("policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	var version int64
	if v, err := tools.GetIntValue("version", d); err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		version, err = findLatestPolicyVersion(ctx, client, policyID)
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
	policy, err := client.GetPolicy(ctx, int64(policyID))
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

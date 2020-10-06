package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataAkamaiPropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataAkamaiPropertyRead")
	log.Debug("Reading Property")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	prop, err := findProperty(ctx, name, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		prop.LatestVersion = version
	}

	rules, err := getRulesForProperty(ctx, prop, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", ErrRulesNotFound, err.Error()))
	}

	body, err := json.Marshal(rules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("rules", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w:%q", tools.ErrValueSet, err.Error()))
	}
	d.SetId(prop.PropertyID)
	return nil
}

func getRulesForProperty(ctx context.Context, property *papi.Property, meta akamai.OperationMeta) (*papi.GetRuleTreeResponse, error) {
	client := inst.Client(meta)
	req := papi.GetRuleTreeRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
	}
	rules, err := client.GetRuleTree(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRulesNotFound, err.Error())
	}
	return rules, nil
}

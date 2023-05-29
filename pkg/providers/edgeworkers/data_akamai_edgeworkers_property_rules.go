package edgeworkers

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEdgeworkersPropertyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEdgeworkersPropertyRulesRead,
		Schema: map[string]*schema.Schema{
			"edgeworker_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of an EdgeWorker ID",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property rule enabling a selected EdgeWorker",
			},
		},
	}
}

func dataEdgeworkersPropertyRulesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("Edgeworkers", "dataEdgeworkersPropertyRulesRead")
	log.Debug("Generating EdgeWorker rules JSON")

	ewID, err := tf.GetIntValue("edgeworker_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	var ewRule struct {
		Name    string `json:"name"`
		Options struct {
			Enabled      bool   `json:"enabled"`
			EdgeWorkerID string `json:"edgeWorkerId"`
		} `json:"options"`
	}

	ewRule.Name = "edgeWorker"
	ewRule.Options.EdgeWorkerID = strconv.Itoa(ewID)
	ewRule.Options.Enabled = true

	ruleJSON, err := json.Marshal(ewRule)
	if err != nil {
		return diag.FromErr(err)
	}
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, ruleJSON, "", "  "); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", formatted.String()); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(strconv.Itoa(ewID))
	return nil
}

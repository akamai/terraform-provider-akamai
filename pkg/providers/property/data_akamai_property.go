package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func dataAkamaiPropertyRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	akactx := akamai.ContextGet(inst.Name())
	log := akactx.Log("PAPI", "dataAkamaiPropertyRead")
	CorrelationID := "[PAPI][dataAkamaiPropertyRead-" + akactx.OperationID() + "]"

	log.Debug("Reading Property")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	property := findProperty(name, CorrelationID)
	if property == nil {
		return diag.FromErr(fmt.Errorf("%w: %s", ErrPropertyNotFound, name))
	}

	version, err := tools.GetIntValue("version", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		property.LatestVersion = version
	}

	rules, err := property.GetRules(CorrelationID)
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
	d.SetId(property.PropertyID)
	return nil
}

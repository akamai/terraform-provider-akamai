package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPCode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCPCodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contract": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataSourceCPCodeRead")

	log.Debug("Read CP Code")

	var name, group, contract string
	var err error

	if name, err = tools.GetStringValue("name", d); err != nil {
		return diag.FromErr(err)
	}

	if group, err = tools.GetStringValue("group", d); err != nil {
		return diag.FromErr(err)
	}

	if contract, err = tools.GetStringValue("contract", d); err != nil {
		return diag.FromErr(err)
	}

	cpCode, err := findCPCode(ctx, name, contract, group, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not load CP codes: %w", err))
	}

	if cpCode == nil {
		return diag.FromErr(fmt.Errorf("%w: invalid CP Code", ErrLookingUpCPCode))
	}

	if err := d.Set("name", cpCode.Name); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("product_ids", cpCode.ProductIDs); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(cpCode.ID)

	log.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

// findCPCode searches all CP codes for a match against given nameOrID
func findCPCode(ctx context.Context, nameOrID, contractID, groupID string, meta akamai.OperationMeta) (*papi.CPCode, error) {
	client := inst.Client(meta)
	r, err := client.GetCPCodes(ctx, papi.GetCPCodesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	})

	if err != nil {
		return nil, err
	}

	for _, cpc := range r.CPCodes.Items {
		if cpc.ID == nameOrID || cpc.ID == "cpc_"+nameOrID || cpc.Name == nameOrID {
			return &cpc, nil
		}
	}

	return nil, nil
}

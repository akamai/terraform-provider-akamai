package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
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
		},
	}
}

func dataSourceCPCodeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataSourceCPCodeRead")

	CorrelationID := "[PAPI][dataSourceCPCodeRead-" + meta.OperationID() + "]"
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
	cpCodes := datasourceCPCodePAPINewCPCodes(contract, group)
	cpCode, err := cpCodes.FindCpCode(name, CorrelationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", ErrLookingUpCPCode, err.Error()))
	}
	if cpCode == nil {
		return diag.FromErr(fmt.Errorf("%w: invalid CP Code", ErrLookingUpCPCode))
	}

	if err := d.Set("name", cpCode.CpcodeName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(cpCode.CpcodeID)

	log.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

func datasourceCPCodePAPINewCPCodes(contractID, groupID string) *papi.CpCodes {
	contract := &papi.Contract{
		ContractID: contractID,
	}
	group := &papi.Group{
		GroupID: groupID,
	}
	return papi.NewCpCodes(contract, group)
}

package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPCode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCPCodeRead,
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

func dataSourceCPCodeRead(d *schema.ResourceData, _ interface{}) error {
	akactx := akamai.ContextGet(inst.Name())
	log := akactx.Log("PAPI", "dataSourceCPCodeRead")
	CorrelationID := "[PAPI][dataSourceCPCodeRead-" + akactx.OperationID() + "]"
	log.Debug("Read CP Code")

	var name, group, contract string
	var err error
	if name, err = tools.GetStringValue("name", d); err != nil {
		return err
	}
	if group, err = tools.GetStringValue("group", d); err != nil {
		return err
	}
	if contract, err = tools.GetStringValue("contract", d); err != nil {
		return err
	}
	cpCodes := datasourceCPCodePAPINewCPCodes(contract, group)
	cpCode, err := cpCodes.FindCpCode(name, CorrelationID)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrLookingUpCPCode, err.Error())
	}
	if cpCode == nil {
		return fmt.Errorf("%w: invalid CP Code", ErrLookingUpCPCode)
	}

	if err := d.Set("name", cpCode.CpcodeName); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("product", cpCode.ProductIDs[0]); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("id", cpCode.CpcodeID); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(cpCode.CpcodeID)

	log.Debug("Read CP Code: %+v", cpCode)
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

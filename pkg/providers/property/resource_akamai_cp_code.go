package property

import (
	"errors"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PAPI CP Code
//
// https://developer.akamai.com/api/luna/papi/data.html#cpcode
// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
func resourceCPCode() *schema.Resource {
	return &schema.Resource{
		Create: resourceCPCodeCreate,
		Read:   resourceCPCodeRead,
		Delete: resourceCPCodeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"product": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCPCodeCreate(d *schema.ResourceData, meta interface{}) error {
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("PAPI", "resourceCPCodeCreate")
	CorrelationID := "[PAPI][resourceCPCodeCreate-" + akactx.OperationID() + "]"

	logger.Debug("Creating CP Code")
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return err
	}
	product, err := tools.GetStringValue("product", d)
	if err != nil {
		return err
	}
	group, err := tools.GetStringValue("group", d)
	if err != nil {
		return err
	}
	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return err
	}
	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCodes := resourceCPCodePAPINewCPCodes(contract, group)
	cpCode, err := cpCodes.FindCpCode(name, CorrelationID)
	if cpCode == nil || err != nil {
		cpCode = cpCodes.NewCpCode()
		cpCode.ProductID = product
		cpCode.CpcodeName = name

		logger.Debug("CPCode: %+v")
		err := cpCode.Save(CorrelationID)
		if err != nil {
			logger.Debug("Error saving")
			var apiError client.APIError
			if errors.As(err, &apiError) {
				logger.Debug("%s", apiError.RawBody)
			}
			return err
		}
	}

	logger.Debug("Resulting CP Code: %#v", cpCode)
	d.SetId(cpCode.CpcodeID)
	return resourceCPCodeRead(d, meta)
}

func resourceCPCodeDelete(d *schema.ResourceData, meta interface{}) error {
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("PAPI", "resourceCPCodeDelete")
	logger.Debug("Deleting CP Code")
	// No PAPI CP Code delete operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return schema.Noop(d, meta)
}

func resourceCPCodeRead(d *schema.ResourceData, meta interface{}) error {
	akactx := akamai.ContextGet(inst.Name())
	logger := akactx.Log("PAPI", "resourceCPCodeRead")
	CorrelationID := "[PAPI][resourceCPCodeRead-" + akactx.OperationID() + "]"

	logger.Debug("Read CP Code")
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return err
	}
	group, err := tools.GetStringValue("group", d)
	if err != nil {
		return err
	}
	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return err
	}
	cpCodes := resourceCPCodePAPINewCPCodes(contract, group)
	cpCode, err := cpCodes.FindCpCode(d.Id(), CorrelationID)
	if cpCode == nil || err != nil {
		cpCode, err = cpCodes.FindCpCode(name, CorrelationID)
		if err != nil {
			return err
		}
	}

	if cpCode == nil {
		return nil
	}

	d.SetId(cpCode.CpcodeID)
	logger.Debug("Read CP Code: %+v", cpCode)
	return nil
}

func resourceCPCodePAPINewCPCodes(contractID, groupID string) *papi.CpCodes {
	contract := &papi.Contract{
		ContractID: contractID,
	}
	group := &papi.Group{
		GroupID: groupID,
	}
	return papi.NewCpCodes(contract, group)
}

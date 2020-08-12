package akamai

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Update: resourceCPCodeUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contract": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCPCodeCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeCreate-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating CP Code")
	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "name", "string")
	}
	product, ok := d.Get("product").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "product", "string")
	}
	group, ok := d.Get("group").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "group", "string")
	}
	contract, ok := d.Get("contract").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "contract", "string")
	}
	cpCodes := resourceCPCodePAPINewCPCodes(contract, group)
	cpCode, err := cpCodes.FindCpCode(name, CorrelationID)
	// TODO: err can indicate that either error was returned while fetching CP Codes from PAPI or no CP Codes were found for provided group and contract
	// this should be modified in client library as currently we do not know whether there was an actual error (in which case err should be returned immediately without proceeding to create the resource)
	if cpCode == nil || err != nil {
		cpCode = cpCodes.NewCpCode()
		cpCode.ProductID = product
		cpCode.CpcodeName = name

		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  CPCode: %+v", cpCode))
		err := cpCode.Save(CorrelationID)
		if err != nil {
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Error saving")
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("%s", err.(client.APIError).RawBody))
			return err
		}
	}
	var found bool
	for _, id := range cpCode.ProductIDs {
		if id == product {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("attempting to modify product ID: %w", ErrPAPICPCodeModify)
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Resulting CP Code: %#v\n\n\n", cpCode))
	d.SetId(cpCode.CpcodeID)
	return resourceCPCodeRead(d, meta)
}

func resourceCPCodeDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeCreate-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting CP Code")
	// No PAPI CP Code delete operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return schema.Noop(d, meta)
}

func resourceCPCodeRead(d *schema.ResourceData, _ interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeRead-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CP Code")
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "name", "string")
	}
	group, ok := d.Get("group").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "group", "string")
	}
	contract, ok := d.Get("contract").(string)
	if !ok {
		return fmt.Errorf("%w: %s, %q", ErrInvalidPropertyType, "contract", "string")
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
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Read CP Code: %+v", cpCode))
	return nil
}

func resourceCPCodeUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("contract") {
		return fmt.Errorf("attempting to modify contract ID: %w", ErrPAPICPCodeModify)
	}
	if d.HasChange("product") {
		return fmt.Errorf("attempting to modify product ID: %w", ErrPAPICPCodeModify)
	}
	if d.HasChange("group") {
		return fmt.Errorf("attempting to modify group ID: %w", ErrPAPICPCodeModify)
	}
	return resourceCPCodeRead(d, meta)
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

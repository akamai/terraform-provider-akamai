package akamai

import (
	"fmt"
	"log"

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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contract": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"product": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCPCodeCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeCreate-" + CreateNonce() + "]"

	//log.Printf("[DEBUG]" + CorrelationID + "  Creating CP Code")
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating CP Code")
	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCodes := resourceCPCodePAPINewCPCodes(d, meta)
	cpCode, err := cpCodes.FindCpCode(d.Get("name").(string), CorrelationID)
	if cpCode == nil || err != nil {
		cpCode = cpCodes.NewCpCode()
		cpCode.ProductID = d.Get("product").(string)
		cpCode.CpcodeName = d.Get("name").(string)
		//log.Printf("[DEBUG]"+CorrelationID+"  CPCode: %#v", cpCode)
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  CPCode: %+v", cpCode))
		err := cpCode.Save(CorrelationID)
		if err != nil {
			log.Print("[DEBUG] Error saving")
			log.Printf("%s", err.(client.APIError).RawBody)
			return err
		}
	}

	//log.Printf("[DEBUG]"+CorrelationID+"  Resulting CP Code: %#v\n\n\n", cpCode)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Resulting CP Code: %#v\n\n\n", cpCode))
	d.SetId(cpCode.CpcodeID)
	return resourceCPCodeRead(d, meta)
}

func resourceCPCodeDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeCreate-" + CreateNonce() + "]"
	//log.Printf("[DEBUG] Deleting CP Code")
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting CP Code")
	// No PAPI CP Code delete operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return schema.Noop(d, meta)
}

func resourceCPCodeRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceCPCodeRead-" + CreateNonce() + "]"
	//PrintLogHeader()
	//log.Printf("[DEBUG] Reading CP Code")
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CP Code")
	cpCodes := resourceCPCodePAPINewCPCodes(d, meta)
	cpCode, err := cpCodes.FindCpCode(d.Id(), CorrelationID)
	if cpCode == nil || err != nil {
		cpCode, err = cpCodes.FindCpCode(d.Get("name").(string), CorrelationID)
		if err != nil {
			return err
		}
	}

	if cpCode == nil {
		return nil
	}

	d.SetId(cpCode.CpcodeID)
	//log.Printf("[DEBUG]"+CorrelationID+"  Read CP Code: %+v", cpCode)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Read CP Code: %+v", cpCode))
	//PrintLogFooter()
	return nil
}

func resourceCPCodePAPINewCPCodes(d *schema.ResourceData, meta interface{}) *papi.CpCodes {
	contract := &papi.Contract{
		ContractID: d.Get("contract").(string),
	}
	group := &papi.Group{
		GroupID: d.Get("group").(string),
	}
	return papi.NewCpCodes(contract, group)
}

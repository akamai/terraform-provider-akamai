package akamai

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
	log.Printf("[DEBUG] Creating CP Code")

	cpCode := resourceCPCodePAPINewCPCodes(d, meta).NewCpCode()
	cpCode.ProductID = d.Get("product").(string)
	cpCode.CpcodeName = d.Get("name").(string)
	err := cpCode.Save()
	if err != nil {
		return err
	}

	d.SetId(cpCode.CpcodeID)

	log.Printf("[DEBUG] Created CP Code: +%v", cpCode)
	return resourceCPCodeRead(d, meta)
}

func resourceCPCodeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting CP Code")

	// No PAPI CP Code delete operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return schema.Noop(d, meta)
}

func resourceCPCodeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CP Code")

	cpCode := resourceCPCodePAPINewCPCodes(d, meta).NewCpCode()
	cpCode.CpcodeID = d.Id()
	err := cpCode.GetCpCode()
	if err != nil {
		return err
	}

	d.Set("name", cpCode.CpcodeName)
	if len(cpCode.ProductIDs) > 0 {
		d.Set("product", cpCode.ProductIDs[0])
	}

	log.Printf("[DEBUG] Read CP Code: %+v", cpCode)
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

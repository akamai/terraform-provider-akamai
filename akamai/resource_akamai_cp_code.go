package akamai

import (
	"errors"
	"log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
)

// PAPI CP Code
//
// https://developer.akamai.com/api/luna/papi/data.html#cpcode
// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
func resourceCPCode() *schema.Resource {
	return &schema.Resource{
		Create: resourceCPCodeCreate,
		Read:   resourceCPCodeRead,
		Update: resourceCPCodeUpdate,
		Delete: resourceCPCodeDelete,
		Exists: resourceCPCodeExists,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
	return nil
}

func resourceCPCodeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting CP Code")

	// No PAPI CP Code delete operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return errors.New("deleting CP Codes is unsupported")
}

func resourceCPCodeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("[DEBUG] Finding CP Code")

	cpCodeName := d.Get("name").(string)
	cpCode, err := resourceCPCodePAPINewCPCodes(d, meta).FindCpCode(cpCodeName)
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] Found CP Code: %+v", cpCode)
	return cpCode != nil, nil
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
	d.Set("product", cpCode.ProductIDs[0])

	log.Printf("[DEBUG] Read CP Code: %+v", cpCode)
	return nil
}

func resourceCPCodeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating CP Code")

	// No PAPI CP Code update operation exists.
	// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
	return errors.New("updating CP Codes is unsupported")
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

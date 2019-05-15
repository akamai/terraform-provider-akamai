package akamai

import (
	//"fmt"
	"fmt"
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCPCode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCPCodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
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
		},
	}
}

func dataSourceCPCodeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CP Code")

	cpCodeName := d.Get("name").(string)

	cpCode, err := datasourceCPCodePAPINewCPCodes(d, meta).FindCpCode(cpCodeName)

	/*cpCode.CpcodeID = d.Id()
	err = cpCode.GetCpCode()
	*/
	if err != nil {
		return err
	}
	if cpCode == nil {
		return fmt.Errorf("Invalid CP Code")
	}

	d.Set("name", cpCode.CpcodeName)
	d.Set("product", cpCode.ProductIDs[0])
	d.Set("id", cpCode.CpcodeID)
	d.SetId(cpCode.CpcodeID)

	log.Printf("[DEBUG] Read CP Code: %+v", cpCode)
	return nil
}

func datasourceCPCodePAPINewCPCodes(d *schema.ResourceData, meta interface{}) *papi.CpCodes {
	contract := &papi.Contract{
		ContractID: d.Get("contract").(string),
	}
	group := &papi.Group{
		GroupID: d.Get("group").(string),
	}
	return papi.NewCpCodes(contract, group)
}

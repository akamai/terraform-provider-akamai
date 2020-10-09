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

// PAPI CP Code
//
// https://developer.akamai.com/api/luna/papi/data.html#cpcode
// https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
func resourceCPCode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPCodeCreate,
		ReadContext:   resourceCPCodeRead,

		// NB: CP Codes cannot be deleted https://developer.akamai.com/api/luna/papi/resources.html#cpcodesapi
		DeleteContext: schema.NoopContext,

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

func resourceCPCodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCPCodeCreate")
	logger.Debugf("Creating CP Code")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	product, err := tools.GetStringValue("product", d)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := tools.GetStringValue("group", d)
	if err != nil {
		return diag.FromErr(err)
	}

	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCode, err := findCPCode(ctx, name, contract, group, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrLookingUpCPCode, err))
	}

	if cpCode == nil {
		cpcID, err := createCPCode(ctx, name, product, contract, group, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(cpcID)
	} else {
		d.SetId(cpCode.ID)
	}

	logger.Debugf("Resulting CP Code: %#v", cpCode)
	return resourceCPCodeRead(ctx, d, m)
}

func resourceCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCPCodeRead")
	logger.Debugf("Read CP Code")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := tools.GetStringValue("group", d)
	if err != nil {
		return diag.FromErr(err)
	}

	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Attempt to find by ID first
	cpCode, err := findCPCode(ctx, d.Id(), contract, group, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Otherwise attempt to find by name
	if cpCode == nil {
		// FIXME: I'm not clear how this could ever happen. A read couldn't happen until after TF created it and it had
		//        been assigned an ID by PAPI and that ID was previously set in the resource, right?
		cpCode, err := findCPCode(ctx, name, contract, group, meta)
		if err != nil {
			return diag.FromErr(err)
		}

		// It really doesn't exist, give up
		if cpCode == nil {
			return diag.Errorf("Couldn't find the CP Code")
		}
	}

	if err := d.Set("name", cpCode.Name); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	// we use the first value returned.  Most cpcodes have but a single product and we need to pick one for comparison.
	if len(cpCode.ProductIDs) == 0 {
		return diag.Errorf("Couldn't find product id on the CP Code")
	}
	if err := d.Set("product", cpCode.ProductIDs[0]); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(cpCode.ID)
	logger.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

// createCPCode attempts to create a CP Code and returns the CP Code ID
func createCPCode(ctx context.Context, name, product, contract, group string, meta akamai.OperationMeta) (string, error) {
	client := inst.Client(meta)
	r, err := client.CreateCPCode(ctx, papi.CreateCPCodeRequest{
		ContractID: contract,
		GroupID:    group,
		CPCode: papi.CreateCPCode{
			ProductID:  product,
			CPCodeName: name,
		},
	})
	if err != nil {
		return "", err
	}

	return r.CPCodeID, nil
}

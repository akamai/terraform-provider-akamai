package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
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
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: `use "contract_id" attribute instead`,
				StateFunc:  addPrefixToState("ctr_"),
			},
			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract_id", "contract"},
				StateFunc:    addPrefixToState("ctr_"),
			},
			"group": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: `use "group_id" attribute instead`,
				StateFunc:  addPrefixToState("grp_"),
			},
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group"},
				StateFunc:    addPrefixToState("grp_"),
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

	var name string
	if got, ok := d.GetOk("name"); ok {
		name = got.(string)
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	} else {
		productID = d.Get("product").(string)
	}
	if !strings.HasPrefix(productID, "prd_") {
		productID = fmt.Sprintf("prd_%s", productID)
	}

	// Schema guarantees group_id/group are strings and one or the other is set
	var groupID string
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	if !strings.HasPrefix(groupID, "grp_") {
		groupID = fmt.Sprintf("grp_%s", groupID)
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	if !strings.HasPrefix(contractID, "ctr_") {
		contractID = fmt.Sprintf("ctr_%s", contractID)
	}

	// Because CPCodes can't be deleted, we re-use an existing CPCode if it's there
	cpCode, err := findCPCode(ctx, name, contractID, groupID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", ErrLookingUpCPCode, err))
	}

	if cpCode == nil {
		cpcID, err := createCPCode(ctx, name, productID, contractID, groupID, meta)
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

	var name string
	if got, ok := d.GetOk("name"); ok {
		name = got.(string)
	}

	// Schema guarantees group_id/group are strings and one or the other is set
	var groupID string
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	if !strings.HasPrefix(groupID, "grp_") {
		groupID = fmt.Sprintf("grp_%s", groupID)
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	if !strings.HasPrefix(contractID, "ctr_") {
		contractID = fmt.Sprintf("ctr_%s", contractID)
	}

	// Attempt to find by ID first
	cpCode, err := findCPCode(ctx, d.Id(), contractID, groupID, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Otherwise attempt to find by name
	if cpCode == nil {
		// FIXME: I'm not clear how this could ever happen. A read couldn't happen until after TF created it and it had
		//        been assigned an ID by PAPI and that ID was previously set in the resource, right?
		cpCode, err := findCPCode(ctx, name, contractID, groupID, meta)
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
func createCPCode(ctx context.Context, name, productID, contractID, groupID string, meta akamai.OperationMeta) (string, error) {
	client := inst.Client(meta)
	r, err := client.CreateCPCode(ctx, papi.CreateCPCodeRequest{
		ContractID: contractID,
		GroupID:    groupID,
		CPCode: papi.CreateCPCode{
			ProductID:  productID,
			CPCodeName: name,
		},
	})
	if err != nil {
		return "", err
	}

	return r.CPCodeID, nil
}

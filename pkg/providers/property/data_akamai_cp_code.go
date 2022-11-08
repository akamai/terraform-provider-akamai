package property

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
)

func dataSourceCPCode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCPCodeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tools.IsNotBlank,
			},
			"contract": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract", "contract_id"},
				ForceNew:     true,
				Deprecated:   akamai.NoticeDeprecatedUseAlias("contract"),
			},
			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract", "contract_id"},
				ForceNew:     true,
			},
			"group": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group", "group_id"},
				ForceNew:     true,
				Deprecated:   akamai.NoticeDeprecatedUseAlias("group"),
			},
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group", "group_id"},
				ForceNew:     true,
			},
			"product_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCPCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataSourceCPCodeRead")

	log.Debug("Read CP Code")

	var name, groupID, contractID string
	var err error

	if name, err = tools.GetStringValue("name", d); err != nil {
		return diag.FromErr(err)
	}

	// load group_id, if not exists, then load group.
	if groupID, err = tools.ResolveKeyStringState(d, "group_id", "group"); err != nil {
		return diag.FromErr(err)
	}
	// set group_id/group in state.
	if err := d.Set("group_id", groupID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", groupID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	// load contract_id, if not exists, then load contract.
	if contractID, err = tools.ResolveKeyStringState(d, "contract_id", "contract"); err != nil {
		return diag.FromErr(err)
	}
	// set contract_id/contract in state.
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract", contractID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	cpCode, err := findCPCode(ctx, name, contractID, groupID, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not load CP codes: %w", err))
	}

	if cpCode == nil {
		return diag.FromErr(fmt.Errorf("%v: invalid CP Code", ErrLookingUpCPCode))
	}

	if err := d.Set("product_ids", cpCode.ProductIDs); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(cpCode.ID)

	log.Debugf("Read CP Code: %+v", cpCode)
	return nil
}

// findCPCode searches all CP codes for a match against given nameOrID
func findCPCode(ctx context.Context, nameOrID, contractID, groupID string, meta akamai.OperationMeta) (*papi.CPCode, error) {
	client := inst.Client(meta)
	r, err := client.GetCPCodes(ctx, papi.GetCPCodesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	})

	if err != nil {
		return nil, err
	}

	for _, cpc := range r.CPCodes.Items {
		if cpc.ID == nameOrID || cpc.ID == "cpc_"+nameOrID || cpc.Name == nameOrID {
			return &cpc, nil
		}
	}

	return nil, fmt.Errorf("%v: CP code: %s", ErrCpCodeNotFound, nameOrID)
}

package property

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyContract() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropertyContractRead,
		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePropertyContractRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	akactx := akamai.ContextGet(inst.Name())

	// demonstrate the context logger
	log := akactx.Log("PAPI", "dataSourcePropertyContractRead")
	CorrelationID := "[PAPI][dataSourcePropertyContractRead-" + akactx.OperationID() + "]"
	contracts := papi.NewContracts()
	group, err := tools.GetStringValue("group", d)
	// If no group, just return the first contract
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if err := contracts.GetContracts(CorrelationID); err != nil {
			return diag.FromErr(fmt.Errorf("error looking up Contracts for group %q: %s", group, err))
		}
		if len(contracts.Contracts.Items) == 0 {
			return diag.FromErr(fmt.Errorf("%w", ErrNoContractsFound))
		}
		d.SetId(contracts.Contracts.Items[0].ContractID)
		return nil
	}
	// Otherwise find the group and return it's first contract
	log.Debug("[Akamai Property Contract] Start Searching for property contract by group")
	groups, err := papi.GetGroups()
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q: %s", ErrLookingUpContract, group, err))
	}

	for _, g := range groups.Groups.Items {
		if g.GroupID != group && g.GroupID != "grp_"+group && g.GroupName != group {
			continue
		}
		if len(g.ContractIDs) == 0 {
			return diag.FromErr(fmt.Errorf("%w: %s", ErrLookingUpContract, group))
		}
		d.SetId(g.ContractIDs[0])
		return nil
	}

	return diag.FromErr(fmt.Errorf("%w; groupID: %q", ErrNoContractsFound, group))
}

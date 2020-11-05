package property

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiContracts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContractsRead,
		Schema: map[string]*schema.Schema{
			"contracts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of contracts",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_id":        {Type: schema.TypeString, Computed: true},
						"contract_type_name": {Type: schema.TypeString, Computed: true},
					}},
			},
		},
	}
}

func dataSourceContractsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataSourceContractsRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	log.Debug("Listing Contracts")
	contracts, err := getContracts(ctx, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing contracts: %w", err))
	}
	if len(contracts.Contracts.Items) == 0 {
		return diag.FromErr(fmt.Errorf("%w", ErrNoContractsFound))
	}
	// setting Account Id as key
	d.SetId(contracts.AccountID)

	/* setting all the json fields in contracts struct to provide
	   more granular access to nested fields */
	ctrList := make([]map[string]string, 0, len(contracts.Contracts.Items))

	for _, c := range contracts.Contracts.Items {
		cMap := map[string]string{
			"contract_id":        c.ContractID,
			"contract_type_name": c.ContractTypeName,
		}
		ctrList = append(ctrList, cMap)
	}
	// setting contracts as key
	if err := d.Set("contracts", ctrList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting contracts: %s", err))
	}
	return nil
}

// Reusable function to fetch all the contracts accessible through a API token
func getContracts(ctx context.Context, meta akamai.OperationMeta) (*papi.GetContractsResponse, error) {
	client := inst.Client(meta)
	ctrs, err := client.GetContracts(ctx)
	if err != nil {
		return nil, err
	}
	return ctrs, nil
}

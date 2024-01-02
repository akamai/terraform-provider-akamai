package property

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContracts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataContractsRead,
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

func dataContractsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akameta.Must(m)
	log := meta.Log("PAPI", "dataContractsRead")
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
func getContracts(ctx context.Context, meta akameta.Meta) (*papi.GetContractsResponse, error) {
	contracts := &papi.GetContractsResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), "contracts", contracts); err != nil {
		if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
			return nil, err
		}
		contracts, err = Client(meta).GetContracts(ctx)
		if err != nil {
			return nil, err
		}
		if err := cache.Set(cache.BucketName(SubproviderName), "contracts", contracts); err != nil {
			if !errors.Is(err, cache.ErrDisabled) {
				return nil, err
			}
		}
	}
	return contracts, nil
}

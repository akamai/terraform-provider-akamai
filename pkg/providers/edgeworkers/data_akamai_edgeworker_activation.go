package edgeworkers

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceEdgeWorkerActivation(config edgeworkersActivationResourceConfig) *schema.Resource {
	return &schema.Resource{
		Description: "Fetch latest activation for given EdgeWorkerID",
		ReadContext: dataEdgeWorkerActivationRead(config),
		Schema: map[string]*schema.Schema{
			"edgeworker_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the EdgeWorker",
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(edgeworkers.ActivationNetworkStaging),
					string(edgeworkers.ActivationNetworkProduction),
				}, false)),
				Description: "The network from which the activation will be fetched",
			},
			"activation_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier of the activation",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of EdgeWorker",
			},
		},
	}
}

func dataEdgeWorkerActivationRead(config edgeworkersActivationResourceConfig) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		meta := meta.Must(m)
		logger := meta.Log("EdgeWorkers", "dataEdgeWorkerActivationsRead")

		ctx = session.ContextWithOptions(
			ctx,
			session.WithContextLog(logger),
		)

		client := meta.Client().GetEdgeWorkers()
		logger.Debug("Reading EdgeWorker Activations")

		edgeworkerID, err := tf.GetIntValue("edgeworker_id", d)
		if err != nil {
			return diag.Errorf("could not get edgeworker_id: %s", err)
		}

		network, err := tf.GetStringValue("network", d)
		if err != nil {
			return diag.Errorf("could not get network: %s", err)
		}

		activation, err := getCurrentActivation(ctx, client, edgeworkerID, network, false, config)
		if err != nil && !errors.Is(err, ErrEdgeworkerNoCurrentActivation) {
			return diag.Errorf("could not get current activation: %s", err)
		}

		if activation != nil {
			if err = d.Set("activation_id", activation.ActivationID); err != nil {
				return diag.Errorf("could not set activation_id: %s", err)
			}

			if err = d.Set("version", activation.Version); err != nil {
				return diag.Errorf("could not set version: %s", err)
			}
		}

		d.SetId(fmt.Sprintf("%d:%s", edgeworkerID, network))
		return nil
	}
}

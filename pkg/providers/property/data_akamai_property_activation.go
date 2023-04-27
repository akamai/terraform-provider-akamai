package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropertyActivationRead,
		Schema:      dataSourcePropertyActivationSchema,
	}
}

var dataSourcePropertyActivationSchema = map[string]*schema.Schema{
	"property_id": {
		Type:        schema.TypeString,
		Required:    true,
		StateFunc:   addPrefixToState("prp_"),
		Description: "The property's unique identifier, including optional `prp_` prefix",
	},
	"version": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: tf.IsNotBlank,
		Description:      "The activated property version. To always use the latest version, enter this value `{resource}.{resource identifier}.{field name}`",
	},
	"network": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     papi.ActivationNetworkStaging,
		Description: "Akamai network to check the activation, either `STAGING` or `PRODUCTION`. `STAGING` is the default",
	},
	"activation_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID given to the activation event",
	},
	"errors": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The contents of `errors` field returned by the API",
	},
	"warnings": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The contents of `warnings` field returned by the API",
	},
	"contact": {
		Type:        schema.TypeSet,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Email addresses used to send activation status changes",
	},
	"status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The property version's activation status on the selected network",
	},
	"note": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Log message assigned to the activation request",
	},
}

func dataSourcePropertyActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "dataSourcePropertyActivationRead")
	client := inst.Client(meta)

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))

	propertyID, err := resolvePropertyID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := resolveVersion(ctx, d, client, propertyID, network)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return diag.Errorf("failed to get activations for property: %v", err)
	}

	if err := d.Set("errors", flattenErrorArray(resp.Errors)); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("warnings", flattenErrorArray(resp.Warnings)); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	for _, act := range resp.Activations.Items {

		if act.Network == network && act.PropertyVersion == version {
			logger.Debugf("Found Existing Activation %s version %d", network, version)

			if err := d.Set("status", string(act.Status)); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("version", act.PropertyVersion); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("activation_id", act.ActivationID); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("note", act.Note); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("contact", act.NotifyEmails); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}

			d.SetId(act.PropertyID + ":" + string(network))

			return nil
		}
	}
	return diag.Errorf("there is no active version on %s network", network)
}

package networklists

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationsCreate,
		ReadContext:   resourceActivationsRead,
		DeleteContext: resourceActivationsDelete,
		Schema: map[string]*schema.Schema{
			"uniqueid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "STAGING",
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Activation Comments",
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"notification_recipients": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	// ActivationPollMinimum is the minumum polling interval for activation creation
	ActivationPollMinimum = time.Minute
)

var (
	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum
)

func resourceActivationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsCreate")

	createActivations := networklists.CreateActivationsRequest{}

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done")
		return nil
	}

	//createActivations := networklist.CreateActivationsRequest{}

	createActivationsreq := networklists.GetActivationsRequest{}

	//ap := networklists.ActivationConfigs{}

	uniqueid, err := tools.GetStringValue("uniqueid", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createActivations.UniqueID = uniqueid

	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createActivations.Network = network

	createActivations.Action = "ACTIVATE"
	//createActivations.ActivationConfigs = append(createActivations.ActivationConfigs, ap)
	createActivations.NotificationRecipients = tools.SetToStringSlice(d.Get("notification_recipients").(*schema.Set))

	postresp, err := client.CreateActivations(ctx, createActivations, true)
	if err != nil {
		logger.Warnf("calling 'createActivations': %s", err.Error())
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))
	d.Set("status", string(postresp.ActivationStatus))

	createActivationsreq.ActivationID = postresp.ActivationID
	activation, err := lookupActivation(ctx, client, createActivationsreq)
	for activation.ActivationStatus != "ACTIVATED" { //!= networklists.StatusActive {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, createActivationsreq)

			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("NETWORKLIST", "resourceActivationsRemove")
	logger.Warnf("calling 'Remove Activations' NOOP ")
	d.SetId("")

	return nil
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")

	getActivations := networklists.GetActivationsRequest{}

	activationID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getActivations.ActivationID = activationID

	activations, err := client.GetActivations(ctx, getActivations)
	if err != nil {
		logger.Warnf("calling 'getActivations': %s", err.Error())
	}

	d.Set("status", activations.ActivationStatus)
	d.SetId(strconv.Itoa(activations.ActivationID))

	return nil
}

func lookupActivation(ctx context.Context, client networklists.NETWORKLISTS, query networklists.GetActivationsRequest) (*networklists.GetActivationsResponse, error) {
	activations, err := client.GetActivations(ctx, query)
	if err != nil {
		return nil, err
	}

	return activations, nil

	return nil, nil
}

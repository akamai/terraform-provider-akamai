package clientlists

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	pollActivationInterval = 30 * time.Second
)

func resourceClientListActivation() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceActivationRead,
		CreateContext: resourceActivationCreate,
		UpdateContext: resourceActivationUpdate,
		DeleteContext: resourceActivationDelete,
		Schema: map[string]*schema.Schema{
			"list_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The client list unique identifier.",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The client list version.",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network environment where you activate your client list: either STAGING or PRODUCTION.",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice([]string{string(clientlists.Staging), string(clientlists.Production)}, false),
				),
			},
			"comments": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "A brief description for the activation.",
				DiffSuppressFunc: suppressCommentsDiff,
			},
			"notification_recipients": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Users to notify via email.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"siebel_ticket_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the Siebel ticket, if the activation is linked to one.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current activation status, either ACTIVE, INACTIVE, MODIFIED, PENDING_ACTIVATION, PENDING_DEACTIVATION, or FAILED.",
			},
		},
	}
}

func resourceActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceActivationRead")
	logger.Debug("Reading client list activation")

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	activationRes, err := client.GetActivation(ctx, clientlists.GetActivationRequest{
		ActivationID: id,
	})
	if err != nil {
		logger.Errorf("calling 'GetActivation' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	// Get client list latest version
	listRes, err := client.GetClientList(ctx, clientlists.GetClientListRequest{
		ListID:       activationRes.ListID,
		IncludeItems: false,
	})
	if err != nil {
		logger.Errorf("calling 'GetClientList' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"list_id":                 activationRes.ListID,
		"comments":                activationRes.Comments,
		"network":                 activationRes.Network,
		"notification_recipients": activationRes.NotificationRecipients,
		"siebel_ticket_id":        activationRes.SiebelTicketID,
		"version":                 listRes.Version,
		"status":                  activationRes.ActivationStatus,
	}

	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceActivationCreate")
	logger.Debug("Creating client list activation")

	attrs, err := getResourceAttrs(d)
	if err != nil {
		diag.FromErr(err)
	}

	req := clientlists.CreateActivationRequest{
		ListID: attrs.ListID,
		ActivationParams: clientlists.ActivationParams{
			Action:                 clientlists.Activate,
			Comments:               attrs.Comments,
			SiebelTicketID:         attrs.SiebelTicketID,
			Network:                clientlists.ActivationNetwork(attrs.Network),
			NotificationRecipients: attrs.Emails,
		},
	}

	res, err := client.CreateActivation(ctx, req)
	if err != nil {
		logger.Errorf("calling 'CreateActivation' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", res.ActivationID))

	err = waitForActivationCompletion(ctx, client, res.ActivationID)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "resourceActivationUpdate")
	logger.Debug("Updating client list activation")

	if d.HasChanges("list_id", "version", "network") {
		attrs, err := getResourceAttrs(d)
		if err != nil {
			return diag.FromErr(err)
		}

		req := clientlists.CreateActivationRequest{
			ListID: attrs.ListID,
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Comments:               attrs.Comments,
				SiebelTicketID:         attrs.SiebelTicketID,
				Network:                clientlists.ActivationNetwork(attrs.Network),
				NotificationRecipients: attrs.Emails,
			},
		}

		res, err := client.CreateActivation(ctx, req)
		if err != nil {
			logger.Errorf("calling 'CreateActivation' failed: %s", err.Error())
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d", res.ActivationID))

		err = waitForActivationCompletion(ctx, client, res.ActivationID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChanges("notification_recipients", "siebel_ticket_id") {
		return diag.Errorf("The following attributes 'notification_recipients', 'siebel_ticket_id' " +
			"cannot be modified after activation.")
	}

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CLIENTLIST", "resourceActivationDelete")
	logger.Debug("Deleting client list activation")

	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Client Lists API does not support activation deletion - resource will only be removed from state",
		},
	}
}

type resourceAttrs struct {
	ListID         string
	Network        string
	Comments       string
	SiebelTicketID string
	Emails         []string
	Version        int64
}

func getResourceAttrs(d *schema.ResourceData) (*resourceAttrs, error) {
	listID, err := tf.GetStringValue("list_id", d)
	if err != nil {
		return nil, err
	}
	version, err := tf.GetInt64Value("version", tf.NewRawConfig(d))
	if err != nil {
		return nil, err
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil {
		return nil, err
	}
	comments, err := tf.GetStringValue("comments", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	siebelTicketID, err := tf.GetStringValue("siebel_ticket_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	emailsSet, err := tf.GetSetValue("notification_recipients", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	emails := tf.SetToStringSlice(emailsSet)

	return &resourceAttrs{
		ListID:         listID,
		Version:        version,
		Network:        network,
		Comments:       comments,
		SiebelTicketID: siebelTicketID,
		Emails:         emails,
	}, nil
}

func waitForActivationCompletion(ctx context.Context, client clientlists.ClientLists, activationID int64) error {
	for {
		select {
		case <-time.After(pollActivationInterval):
			r, err := client.GetActivation(ctx, clientlists.GetActivationRequest{ActivationID: activationID})
			if err != nil {
				return fmt.Errorf("polling activation failed: %s", err)
			}

			if r.ActivationStatus == clientlists.Active {
				return nil
			} else if r.ActivationStatus == clientlists.Failed {
				return fmt.Errorf("activation failed")
			}
		case <-ctx.Done():
			return fmt.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
}

// Suppress comments diff when activation is not required,
// and activation resource cannot be updated.
func suppressCommentsDiff(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != newValue && d.HasChanges("list_id", "version", "network") {
		return false
	}
	return true
}

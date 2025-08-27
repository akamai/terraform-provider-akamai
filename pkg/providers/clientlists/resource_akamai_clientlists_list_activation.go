package clientlists

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	waitActivationCompletionTimeout = 30 * time.Minute
	activationRetryMaxAttempts      = 3
	activationRetryTimeout          = 45 * time.Second
)

var (
	pollActivationInterval   = 30 * time.Second
	activationRetryBaseDelay = 3 * time.Second
	errActivationFailed      = errors.New("activation failed")
)

func resourceClientListActivation() *schema.Resource {
	return &schema.Resource{
		ReadContext:   Read,
		CreateContext: Create,
		UpdateContext: Update,
		DeleteContext: Delete,
		CustomizeDiff: customdiff.All(
			markStatusComputed,
		),
		Importer: &schema.ResourceImporter{
			StateContext: ImportState,
		},
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
				DiffSuppressFunc: suppressFieldDiff,
			},
			"notification_recipients": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Users to notify via email.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressFieldDiff,
			},
			"siebel_ticket_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Identifies the Siebel ticket, if the activation is linked to one.",
				DiffSuppressFunc: suppressFieldDiff,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current activation status, either ACTIVE, INACTIVE, MODIFIED, PENDING_ACTIVATION, PENDING_DEACTIVATION, or FAILED.",
			},
		},
	}
}

// Read implements resource's Read method
func Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "Read")
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

// Create implements resource's Create method
func Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "Create")
	logger.Debug("Creating client list activation")

	diags := activate(ctx, d, meta, client, logger)
	if diags.HasError() {
		return diags
	}

	return Read(ctx, d, m)
}

// Update implements resource's Update method
func Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "Update")
	logger.Debug("Updating client list activation")

	isActiveStatus := d.Get("status").(string) == string(clientlists.Active)
	hasChanges := d.HasChanges("list_id", "version", "network")

	if !isActiveStatus || hasChanges {
		diags := activate(ctx, d, meta, client, logger)
		if diags.HasError() {
			return diags
		}
	}

	return Read(ctx, d, m)
}

// Delete implements resource's Delete method
func Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "Delete")
	logger.Debug("Deleting client list activation")

	attrs, err := getResourceAttrs(d, true)
	if err != nil {
		return diag.FromErr(err)
	}

	req := clientlists.CreateDeactivationRequest{
		ListID: attrs.ListID,
		ActivationParams: clientlists.ActivationParams{
			Action:                 clientlists.Deactivate,
			Comments:               attrs.Comments,
			SiebelTicketID:         attrs.SiebelTicketID,
			Network:                clientlists.ActivationNetwork(attrs.Network),
			NotificationRecipients: attrs.Emails,
		},
	}

	res, err := createDeactivationWithRetry(ctx, meta, client, req)
	if err != nil {
		logger.Errorf("calling 'CreateDeactivation' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	_, err = waitForActivationCompletion(ctx, client, res.ActivationID, clientlists.Deactivated)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func activate(ctx context.Context, d *schema.ResourceData, meta meta.Meta, client clientlists.ClientLists, logger akalog.Interface) diag.Diagnostics {
	attrs, err := getResourceAttrs(d, false)
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

	res, err := createActivationWithRetry(ctx, meta, client, req)
	if err != nil {
		logger.Errorf("calling 'CreateActivation' failed: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", res.ActivationID))

	_, err = waitForActivationCompletion(ctx, client, res.ActivationID, clientlists.Active)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

type resourceAttrs struct {
	ListID         string
	Network        string
	Comments       string
	SiebelTicketID string
	Emails         []string
	Version        int64
}

func getResourceAttrs(d *schema.ResourceData, destroy bool) (*resourceAttrs, error) {
	listID, err := tf.GetStringValue("list_id", d)
	if err != nil {
		return nil, err
	}

	var version int64
	if destroy { // Don’t Use tf.NewRawConfig(d) in Destroy. It’s intended for config access only, which is unavailable during destroy.
		ver, ok := d.Get("version").(int)
		if !ok {
			return nil, fmt.Errorf("%w: %s", errors.New("value not found"), "version")
		}
		version = int64(ver)
	} else {
		version, err = tf.GetInt64Value("version", tf.NewRawConfig(d))
		if err != nil {
			return nil, err
		}
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

func waitForActivationCompletion(ctx context.Context, client clientlists.ClientLists, activationID int64, status clientlists.ActivationStatus) (*clientlists.GetActivationResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, waitActivationCompletionTimeout)
	defer cancel()
	for {
		select {
		case <-time.After(pollActivationInterval):
			activation, err := client.GetActivation(ctxWithTimeout, clientlists.GetActivationRequest{ActivationID: activationID})
			if err != nil {
				return nil, fmt.Errorf("polling activation failed: %s", err)
			}

			switch activation.ActivationStatus {
			case status:
				return activation, nil
			case clientlists.Failed:
				return nil, errActivationFailed
			}
		case <-ctxWithTimeout.Done():
			if errors.Is(ctxWithTimeout.Err(), context.DeadlineExceeded) {
				return nil, fmt.Errorf("activation polling timed out for activation ID: %d", activationID)
			}
			return nil, fmt.Errorf("activation context terminated: %w", ctx.Err())
		}
	}
}

func retry[T any](ctx context.Context, meta meta.Meta, operationName string, operation func(ctx context.Context) (*T, error),
) (*T, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, activationRetryTimeout)
	logger := meta.Log("CLIENTLIST", operationName)
	defer cancel()

	var lastErr error
	var result *T

	for attempt := 0; attempt < activationRetryMaxAttempts; attempt++ {
		if ctxWithTimeout.Err() != nil {
			logger.Warnf("%s context cancelled before attempt %d: %v", operationName, attempt+1, ctxWithTimeout.Err())
			return nil, ctxWithTimeout.Err()
		}

		result, lastErr = operation(ctxWithTimeout)
		if lastErr == nil {
			return result, nil
		}

		var e *clientlists.Error
		if errors.As(lastErr, &e) && e.StatusCode != http.StatusInternalServerError {
			return nil, lastErr
		}

		if attempt < activationRetryMaxAttempts-1 {
			delay := activationRetryBaseDelay * time.Duration(int64(math.Pow(3, float64(attempt))))
			logger.Warnf("%s attempt %d failed: %s. Retrying in %s...", operationName, attempt+1, lastErr.Error(), delay)

			select {
			case <-time.After(delay):
			case <-ctxWithTimeout.Done():
				logger.Warnf("%s context timeout during retry delay on attempt %d: %v", operationName, attempt+1, ctxWithTimeout.Err())
				return nil, ctxWithTimeout.Err()
			}
		}
	}

	logger.Errorf("%s failed after %d attempts: %v", operationName, activationRetryMaxAttempts, lastErr)
	return nil, lastErr
}

func createActivationWithRetry(ctx context.Context, meta meta.Meta, client clientlists.ClientLists, req clientlists.CreateActivationRequest,
) (*clientlists.CreateActivationResponse, error) {
	return retry[clientlists.CreateActivationResponse](ctx, meta, "CreateActivation", func(ctx context.Context) (*clientlists.CreateActivationResponse, error) {
		return client.CreateActivation(ctx, req)
	})
}

func createDeactivationWithRetry(ctx context.Context, meta meta.Meta, client clientlists.ClientLists, req clientlists.CreateDeactivationRequest,
) (*clientlists.CreateDeactivationResponse, error) {
	return retry[clientlists.CreateDeactivationResponse](ctx, meta, "CreateDeactivation", func(ctx context.Context) (*clientlists.CreateDeactivationResponse, error) {
		return client.CreateDeactivation(ctx, req)
	})
}

// Suppress diff on callers field when activation is not required
func suppressFieldDiff(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != newValue && d.HasChanges("list_id", "version", "network") {
		return false
	}
	return true
}

// ImportState implements resource's ImportState method
func ImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("CLIENTLIST", "importActivationState")
	logger.Debug("Importing client list activation")

	listID, network, err := parseActivationImportArg(d.Id())
	if err != nil {
		return nil, err
	}

	res, err := client.GetActivationStatus(ctx, clientlists.GetActivationStatusRequest{
		ListID:  listID,
		Network: clientlists.ActivationNetwork(network),
	})
	if err != nil {
		return nil, err
	}

	if res.ActivationStatus == clientlists.PendingActivation || res.ActivationStatus == clientlists.PendingDeactivation {
		var status clientlists.ActivationStatus

		if res.ActivationStatus == clientlists.PendingActivation {
			status = clientlists.Active
		} else {
			status = clientlists.Deactivated
		}

		activation, err := waitForActivationCompletion(ctx, client, res.ActivationID, status)
		if err != nil && !errors.Is(err, errActivationFailed) {
			return nil, err
		}
		res.ActivationStatus = activation.ActivationStatus
	}

	fields := map[string]interface{}{
		"list_id":                 res.ListID,
		"comments":                res.Comments,
		"network":                 res.Network,
		"notification_recipients": res.NotificationRecipients,
		"siebel_ticket_id":        res.SiebelTicketID,
		"version":                 res.Version,
		"status":                  res.ActivationStatus,
	}

	d.SetId(fmt.Sprintf("%d", res.ActivationID))

	if err = tf.SetAttrs(d, fields); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func parseActivationImportArg(arg string) (string, string, error) {
	id := strings.Split(arg, ":")
	if len(id) != 2 {
		return "", "", fmt.Errorf("invalid client list activation identifier: %s. "+
			"Correct format is: list_id:network. For example: 123_ABC:PRODUCTION", arg)
	}

	listID, network := id[0], id[1]

	if network != string(clientlists.Staging) && network != string(clientlists.Production) {
		return "", "", fmt.Errorf("invalid network attribute: %s ", network)
	}

	return listID, network, nil
}

// markStatusComputed is a schema.CustomizeDiffFunc for akamai_clientlist_activation resource,
// which sets status field as computed
// if status client list activation is required
func markStatusComputed(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("CLIENTLIST", "markStatusComputed")

	status, err := tf.GetStringValue("status", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if status != string(clientlists.Active) {
		logger.Debug("setting status as new computed")
		if err := d.SetNewComputed("status"); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}

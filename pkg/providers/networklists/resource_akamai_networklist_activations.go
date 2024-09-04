package networklists

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// createNetworkListActivationMutex enforces single-thread access to the CreateActivations call
	createNetworkListActivationMutex sync.Mutex
)

// network_lists v2
//
// https://techdocs.akamai.com/network-lists/reference/api
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationsCreate,
		ReadContext:   resourceActivationsRead,
		UpdateContext: resourceActivationsUpdate,
		DeleteContext: resourceActivationsDelete,
		Schema: map[string]*schema.Schema{
			"network_list_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier of the network list",
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STAGING",
				Description: "The Akamai network on which the list is activated: STAGING or PRODUCTION",
			},
			"notes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Activation Comments",
				Description:      "Descriptive text to accompany the activation",
				DiffSuppressFunc: suppressFieldsForNetworkListActivation,
			},
			"notification_emails": {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description:      "List of email addresses of Control Center users who receive an email when activation of this list is complete",
				DiffSuppressFunc: suppressFieldsForNetworkListActivation,
			},
			"sync_point": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies the sync point of the network list to be activated",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `This network list's current activation status in the environment specified by the "network" attribute`,
			},
		},
	}
}

const (
	// ActivationPollMinimum is the minimum polling interval for activation creation
	ActivationPollMinimum = time.Minute
)

var (
	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum

	// CreateActivationRetry poll wait time code waits between retries for activation creation
	CreateActivationRetry = 10 * time.Second
)

func resourceActivationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsCreate")
	logger.Debug("Creating resource activation")

	networkListID, err := tf.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	comments, err := tf.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmails, ok := d.Get("notification_emails").(*schema.Set)
	if !ok {
		return diag.Errorf("Activation Read failed")
	}

	createResponse, diagErr := createActivation(ctx, client, networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 string(networklists.ActivationTypeActivate),
		NotificationRecipients: tf.SetToStringSlice(notificationEmails),
	})
	if diagErr != nil {
		return diagErr
	}
	d.SetId(strconv.Itoa(createResponse.ActivationID))
	if err := d.Set("status", string(createResponse.ActivationStatus)); err != nil {
		return diag.FromErr(err)
	}

	lookupResponse, err := lookupActivation(ctx, client, networklists.GetActivationRequest{ActivationID: createResponse.ActivationID})
	if err != nil {
		return diag.FromErr(err)
	}

	if err = pollActivation(ctx, client, lookupResponse.ActivationStatus, lookupResponse.ActivationID); err != nil {
		return diag.FromErr(err)
	}

	return resourceActivationsRead(ctx, d, m)
}

func createActivation(ctx context.Context, client networklists.NTWRKLISTS, params networklists.CreateActivationsRequest) (*networklists.CreateActivationsResponse, diag.Diagnostics) {
	createNetworkListActivationMutex.Lock()
	defer func() {
		createNetworkListActivationMutex.Unlock()
	}()

	log := hclog.FromContext(ctx)

	errMsg := "create failed"
	switch params.Action {
	case string(networklists.ActivationTypeActivate):
		errMsg = "create activation failed"
	case string(networklists.ActivationTypeDeactivate):
		errMsg = "create deactivation failed"
	}

	createActivationRetry := CreateActivationRetry

	for {
		log.Debug("creating activation")
		create, err := client.CreateActivations(ctx, params)

		if err == nil {
			return create, nil
		}
		log.Debug("%s: retrying: %w", errMsg, err)

		if !isCreateActivationErrorRetryable(err) {
			return nil, diag.Errorf("%s: %s", errMsg, err)
		}

		select {
		case <-time.After(createActivationRetry):
			createActivationRetry = capDuration(createActivationRetry*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			return nil, diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsRead")
	logger.Debug("Reading resource activation")

	activationID, errconv := strconv.Atoi(d.Id())
	if errconv != nil {
		return diag.FromErr(errconv)
	}

	getResponse, err := client.GetActivation(ctx, networklists.GetActivationRequest{ActivationID: activationID})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", getResponse.ActivationStatus); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(getResponse.ActivationID))

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsUpdate")
	logger.Debug("Updating resource activation")

	networkListID, err := tf.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	comments, err := tf.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmails, err := tf.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createResponse, diagErr := createActivation(ctx, client, networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 string(networklists.ActivationTypeActivate),
		NotificationRecipients: tf.SetToStringSlice(notificationEmails),
	})
	if diagErr != nil {
		return diagErr
	}
	d.SetId(strconv.Itoa(createResponse.ActivationID))
	if err := d.Set("status", string(createResponse.ActivationStatus)); err != nil {
		return diag.FromErr(err)
	}

	lookupRequest := networklists.GetActivationRequest{ActivationID: createResponse.ActivationID}
	lookupResponse, err := lookupActivation(ctx, client, lookupRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = pollActivation(ctx, client, lookupResponse.ActivationStatus, lookupResponse.ActivationID); err != nil {
		return diag.FromErr(err)
	}
	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("NETWORKLIST", "resourceActivationsDelete")
	logger.Debug("removing activation from local state")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "removing activation resource (will be removed from local state only)",
		},
	}
}

func lookupActivation(ctx context.Context, client networklists.NTWRKLISTS, query networklists.GetActivationRequest) (*networklists.GetActivationResponse, error) {
	activation, err := client.GetActivation(ctx, query)
	if err != nil {
		return nil, err
	}
	return activation, nil
}

func suppressFieldsForNetworkListActivation(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != newValue && d.HasChanges("network_list_id", "network") {
		return false
	}
	return true
}

func pollActivation(ctx context.Context, client networklists.NTWRKLISTS, activationStatus string, activationID int) error {
	retriesMax := 5
	retries5xx := 0

	for activationStatus != string(networklists.StatusActive) {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, networklists.GetActivationRequest{ActivationID: activationID})

			if err != nil {
				var target = &networklists.Error{}
				if !errors.As(err, &target) {
					return fmt.Errorf("error has unexpected type: %T", err)
				}
				if isCreateActivationErrorRetryable(target) {
					retries5xx = retries5xx + 1
					if retries5xx > retriesMax {
						return fmt.Errorf("reached max number of 5xx retries: %d", retries5xx)
					}
					continue
				}
				return err
			}
			retries5xx = 0
			activationStatus = act.ActivationStatus

		case <-ctx.Done():
			return fmt.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
	return nil
}

func isCreateActivationErrorRetryable(err error) bool {
	var responseErr = &networklists.Error{}
	if !errors.As(err, &responseErr) {
		return false
	}
	if responseErr.StatusCode < 500 &&
		responseErr.StatusCode != 422 &&
		responseErr.StatusCode != 409 {
		return false
	}
	return true
}

func capDuration(t time.Duration, tMax time.Duration) time.Duration {
	if t > tMax {
		return tMax
	}
	return t
}

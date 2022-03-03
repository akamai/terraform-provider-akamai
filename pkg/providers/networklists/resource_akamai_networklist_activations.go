package networklists

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// createNetworkListActivationMutex enforces single-thread access to the CreateActivations call
	createNetworkListActivationMutex sync.Mutex
)

// network_lists v2
//
// https://developer.akamai.com/api/cloud_security/network_lists/v2.html
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationsCreate,
		ReadContext:   resourceActivationsRead,
		UpdateContext: resourceActivationsUpdate,
		DeleteContext: resourceActivationsDelete,
		Schema: map[string]*schema.Schema{
			"network_list_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STAGING",
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Activation Comments",
			},
			"activate": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("activate"),
			},
			"notification_emails": {
				Type:     schema.TypeSet,
				Required: true,
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

	networkListID, err := tools.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	comments, err := tools.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmails, ok := d.Get("notification_emails").(*schema.Set)
	if !ok {
		return diag.Errorf("Activation Read failed")
	}

	createRequest := networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 "ACTIVATE",
		NotificationRecipients: tools.SetToStringSlice(notificationEmails),
	}
	createResponse, err := createActivation(ctx, client, createRequest)
	if err != nil {
		logger.Debugf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Debugf("calling 'createActivations': RESPONSE %v", createResponse)
	d.SetId(strconv.Itoa(createResponse.ActivationID))
	if err := d.Set("status", string(createResponse.ActivationStatus)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	lookupRequest := networklists.GetActivationRequest{ActivationID: createResponse.ActivationID}
	lookupResponse, err := lookupActivation(ctx, client, lookupRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("calling 'createActivations': GET STATUS ID %v", lookupResponse)

	for lookupResponse.ActivationStatus != "ACTIVATED" {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, lookupRequest)

			if err != nil {
				return diag.FromErr(err)
			}
			lookupResponse = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func createActivation(ctx context.Context, client networklists.NTWRKLISTS, params networklists.CreateActivationsRequest) (*networklists.CreateActivationsResponse, error) {
	createNetworkListActivationMutex.Lock()
	defer func() {
		createNetworkListActivationMutex.Unlock()
	}()

	postResp, err := client.CreateActivations(ctx, params)
	return postResp, err
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")

	activationID, errconv := strconv.Atoi(d.Id())
	if errconv != nil {
		return diag.FromErr(errconv)
	}

	getRequest := networklists.GetActivationRequest{ActivationID: activationID}
	getResponse, err := client.GetActivation(ctx, getRequest)
	if err != nil {
		logger.Warnf("calling 'getActivations': %s", err.Error())
	}

	if err := d.Set("status", getResponse.ActivationStatus); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(strconv.Itoa(getResponse.ActivationID))

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsUpdate")

	networkListID, err := tools.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	comments, err := tools.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmails, ok := d.Get("notification_emails").(*schema.Set)
	if !ok {
		return diag.Errorf("Activation Read failed")
	}

	createRequest := networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 "ACTIVATE",
		NotificationRecipients: tools.SetToStringSlice(notificationEmails),
	}
	createResponse, err := client.CreateActivations(ctx, createRequest)
	if err != nil {
		logger.Debugf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Debugf("calling 'createActivations': RESPONSE %v", createResponse)
	d.SetId(strconv.Itoa(createResponse.ActivationID))
	if err := d.Set("status", string(createResponse.ActivationStatus)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	lookupRequest := networklists.GetActivationRequest{ActivationID: createResponse.ActivationID}
	lookupResponse, err := lookupActivation(ctx, client, lookupRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("calling 'createActivations': GET STATUS ID %v", lookupResponse)

	for lookupResponse.ActivationStatus != "ACTIVATED" {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, lookupRequest)

			if err != nil {
				return diag.FromErr(err)
			}
			lookupResponse = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("NETWORKLIST", "resourceActivationsRemove")
	logger.Warnf("calling 'Remove Activations' NOOP ")
	d.SetId("")
	return nil
}

func lookupActivation(ctx context.Context, client networklists.NTWRKLISTS, query networklists.GetActivationRequest) (*networklists.GetActivationResponse, error) {
	activation, err := client.GetActivation(ctx, query)
	if err != nil {
		return nil, err
	}
	return activation, nil
}

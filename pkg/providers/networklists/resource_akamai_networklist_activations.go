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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Activation Comments",
				Description: "Descriptive text to accompany the activation",
			},
			"activate": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("activate"),
			},
			"notification_emails": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of email addresses of Control Center users who receive an email when activation of this list is complete",
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
	logger.Debug("Creating resource activation")

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

	createResponse, err := createActivation(ctx, client, networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 "ACTIVATE",
		NotificationRecipients: tools.SetToStringSlice(notificationEmails),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(createResponse.ActivationID))
	if err := d.Set("status", string(createResponse.ActivationStatus)); err != nil {
		return diag.FromErr(err)
	}

	lookupResponse, err := lookupActivation(ctx, client, networklists.GetActivationRequest{ActivationID: createResponse.ActivationID})
	if err != nil {
		return diag.FromErr(err)
	}

	for lookupResponse.ActivationStatus != "ACTIVATED" {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, networklists.GetActivationRequest{ActivationID: createResponse.ActivationID})

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

	// Get the current syncpoint of this network list, which may have changed since this activation was created.
	networkListID := getResponse.NetworkList.UniqueID
	networklist, err := client.GetNetworkList(ctx, networklists.GetNetworkListRequest{UniqueID: networkListID})

	if err = d.Set("sync_point", networklist.SyncPoint); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(getResponse.ActivationID))

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceActivationsUpdate")
	logger.Debug("Updating resource activation")

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
	notificationEmails, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.CreateActivations(ctx, networklists.CreateActivationsRequest{
		UniqueID:               networkListID,
		Network:                network,
		Comments:               comments,
		Action:                 "ACTIVATE",
		NotificationRecipients: tools.SetToStringSlice(notificationEmails),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(response.ActivationID))
	if err := d.Set("status", string(response.ActivationStatus)); err != nil {
		return diag.FromErr(err)
	}

	lookupRequest := networklists.GetActivationRequest{ActivationID: response.ActivationID}
	lookupResponse, err := lookupActivation(ctx, client, lookupRequest)
	if err != nil {
		return diag.FromErr(err)
	}

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

package appsec

import (
	"fmt"
	"strconv"
	"time"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		Create: resourceActivationsCreate,
		Read:   resourceActivationsRead,
		Delete: resourceActivationsDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "STAGING",
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Activation Notes",
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"notification_emails": {
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

func resourceActivationsCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceActivationsCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating Activations")

	activations := appsec.NewActivationsResponse()

	postpayload := appsec.NewActivationsPost()
	ap := appsec.ActivationConfigs{}
	ap.ConfigID = d.Get("config_id").(int)
	ap.ConfigVersion = d.Get("version").(int)
	postpayload.Network = d.Get("network").(string)
	postpayload.Action = "ACTIVATE"
	postpayload.ActivationConfigs = append(postpayload.ActivationConfigs, ap)
	postpayload.NotificationEmails = tools.SetToStringSlice(d.Get("notification_emails").(*schema.Set))

	if d.Get("activate").(bool) {

		postresp, err := activations.SaveActivations(postpayload, true, CorrelationID)
		if err != nil {
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
			return err
		}

		d.SetId(strconv.Itoa(postresp.ActivationID))
		d.Set("status", string(postresp.Status))
		go activations.PollStatus(postresp.ActivationID, CorrelationID)

	polling:
		for activations.Status != appsec.StatusActive {
			select {
			case statusChanged := <-activations.StatusChange:
				edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf(" Activation Status: %s\n", postresp.Status))
				if statusChanged == false {
					break polling
				}
				continue polling
			case <-time.After(time.Minute * 40):
				edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Activation Timeout (40 minutes)")
				break polling
			}
		}
	} else {
		d.SetId("none")
	}

	return resourceActivationsRead(d, meta)
}

func resourceActivationsDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceActivationsDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting Activations")

	activations := appsec.NewActivationsResponse()

	if d.Id() == "" {
		return nil
	}

	activationid, _ := strconv.Atoi(d.Id())
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("activationid  %v\n", activationid))
	postpayload := appsec.NewActivationsPost()
	ap := appsec.ActivationConfigs{}
	ap.ConfigID = d.Get("config_id").(int)
	ap.ConfigVersion = d.Get("version").(int)
	postpayload.Network = d.Get("network").(string)
	postpayload.NotificationEmails = tools.SetToStringSlice(d.Get("notification_emails").(*schema.Set))

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Deactivating %s \n", postpayload.Network))

	postpayload.Action = "DEACTIVATE"

	postpayload.ActivationConfigs = append(postpayload.ActivationConfigs, ap)

	postresp, err := activations.DeactivateActivations(postpayload, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	go activations.PollStatus(postresp.ActivationID, CorrelationID)

polling:
	for activations.Status != appsec.StatusDeactivated {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf(" Activation Status: %s\n", activations.Status))
		select {
		case statusChanged := <-activations.StatusChange:
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf(" Activation Status: %s\n", activations.Status))
			if statusChanged == false {
				break polling
			}
			continue polling
		case <-time.After(time.Minute * 40):
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Activation Timeout (40 minutes)")
			break polling
		}
	}

	d.Set("status", string(activations.Status))

	d.SetId("")

	return nil
}

func resourceActivationsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceActivationsRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read Activations")

	activations := appsec.NewActivationsResponse()

	activationid, _ := strconv.Atoi(d.Id())

	_, err := activations.GetActivations(activationid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.Set("status", activations.Status)
	d.SetId(strconv.Itoa(activations.ActivationID))

	return nil
}

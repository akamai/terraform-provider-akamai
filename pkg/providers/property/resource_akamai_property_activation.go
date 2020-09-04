package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyActivationCreate,
		ReadContext:   resourcePropertyActivationRead,
		UpdateContext: resourcePropertyActivationUpdate,
		DeleteContext: resourcePropertyActivationDelete,
		Exists:        resourcePropertyActivationExists,
		Schema:        akamaiPropertyActivationSchema,
	}
}

var akamaiPropertyActivationSchema = map[string]*schema.Schema{
	"property": {
		Type:     schema.TypeString,
		Required: true,
	},
	"version": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"network": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "staging",
	},
	"activate": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"contact": {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourcePropertyActivationCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationCreate")
	CorrelationID := "[PAPI][resourcePropertyActivationCreate-" + meta.OperationID() + "]"
	d.Partial(true)

	property := papi.NewProperty(papi.NewProperties())
	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}
	property.PropertyID = propertyID
	err = property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to find property: %w", err))
	}

	// TODO: SetPartial is technically deprecated and only valid when an error
	// will retruned. Determine if this is necessary here.
	d.Partial(true)
	defer func() {
		d.Partial(false)
		logger.Debugf("Done")
	}()

	if err := d.Set("property", property.PropertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		d.Partial(false)
		logger.Debugf("Done")
		return nil
	}
	activation, err := activateProperty(property, d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(activation.ActivationID)
	if err := d.Set("version", activation.PropertyVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("status", string(activation.Status)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	go activation.PollStatus(property)

	for activation.Status != papi.StatusActive {
		select {
		case statusChanged := <-activation.StatusChange:
			logger.Debugf("Property Status: %s", activation.Status)
			if statusChanged == false {
				return nil
			}
			continue
		case <-time.After(time.Minute * 90):
			logger.Debugf("Activation Timeout (90 minutes)")
			return nil
		}
	}
	return nil
}

func resourcePropertyActivationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationDelete")
	CorrelationID := "[PAPI][resourcePropertyActivationDelete-" + meta.OperationID() + "]"

	logger.Debugf("DEACTIVATE PROPERTY")
	property := papi.NewProperty(papi.NewProperties())
	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}
	property.PropertyID = propertyID
	err = property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("DEACTIVE PROPERTY %v", property)
	networkVal, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network := papi.NetworkValue(networkVal)
	propertyVersion := property.ProductionVersion
	if network == "STAGING" {
		propertyVersion = property.StagingVersion
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Version to deactivate is %d and current active %s version is %d", version, network, propertyVersion)
	defer func() {
		d.SetId("")
		log.Debug("Done")
	}()
	if propertyVersion != version {
		return nil
	}
	// The current active version is the one we need to deactivate
	logger.Debugf("Deactivating %s version %d", network, version)
	activation, err := deactivateProperty(property, d, papi.NetworkValue(networkVal), CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	go activation.PollStatus(property)

polling:
	for activation.Status != papi.StatusActive {
		select {
		case statusChanged := <-activation.StatusChange:
			logger.Debugf("Property Status: %s", activation.Status)
			if statusChanged == false {
				break polling
			}
			continue polling
		case <-time.After(time.Minute * 90):
			logger.Debugf("Activation Timeout (90 minutes)")
			break polling
		}
	}

	if err := d.Set("status", string(activation.Status)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	return nil
}

//Todo should be part of ReadContext
func resourcePropertyActivationExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationExists")
	CorrelationID := "[PAPI][resourcePropertyActivationExists-" + meta.OperationID() + "]"
	property := papi.NewProperty(papi.NewProperties())
	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return false, err
	}
	property.PropertyID = propertyID
	err = property.GetProperty(CorrelationID)
	if err != nil {
		return false, err
	}

	activations, err := property.GetActivations()
	if err != nil {
		// No activations found
		return false, nil
	}

	networkVal, err := tools.GetStringValue("network", d)
	if err != nil {
		return false, err
	}
	network := papi.NetworkValue(networkVal)
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return false, err
	}
	for _, activation := range activations.Activations.Items {
		if activation.Network == network && activation.PropertyVersion == version {
			logger.Debugf("Found Existing Activation %s version %d", network, version)
			return true, nil
		}
	}
	logger.Debugf("Did Not Find Existing Activation %s version %d", network, version)
	return false, nil
}

func resourcePropertyActivationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationRead")
	CorrelationID := "[PAPI][resourcePropertyActivationRead-" + meta.OperationID() + "]"
	property := papi.NewProperty(papi.NewProperties())
	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}
	property.PropertyID = propertyID
	err = property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}

	activations, err := property.GetActivations()
	if err != nil {
		// No activations found
		return nil
	}

	networkVal, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network := papi.NetworkValue(networkVal)
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, activation := range activations.Activations.Items {
		if activation.Network != network || activation.PropertyVersion != version {
			continue
		}
		logger.Debugf("Found Existing Activation %s version %d", network, version)
		d.SetId(activation.ActivationID)
		if err := d.Set("status", string(activation.Status)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("version", activation.PropertyVersion); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourcePropertyActivationUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationRead")
	CorrelationID := "[PAPI][resourcePropertyActivationUpdate-" + meta.OperationID() + "]"

	logger.Debugf("UPDATING")
	logger.Debugf("Fetching property")
	property := papi.NewProperty(papi.NewProperties())
	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}
	property.PropertyID = propertyID
	err = property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := getActivation(d, property, papi.ActivationTypeActivate, papi.NetworkValue(network), CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	foundActivation, err := findExistingActivation(property, activation, CorrelationID, logger)
	if err == nil {
		activation = foundActivation
	}

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		logger.Debugf("Done")
		return resourcePropertyActivationRead(nil, d, m)
	}
	// No activation in progress, create a new one
	if foundActivation == nil {
		activation, err = activateProperty(property, d, CorrelationID, logger)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(activation.ActivationID)
	go activation.PollStatus(property)

polling:
	for activation.Status != papi.StatusActive {
		select {
		case statusChanged := <-activation.StatusChange:
			logger.Debugf("Property Status: %s\n", activation.Status)
			if statusChanged == false {
				break polling
			}
			continue polling
		case <-time.After(time.Minute * 90):
			logger.Debugf("Activation Timeout (90 minutes)")
			break polling
		}
	}
	if err := d.Set("status", string(activation.Status)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("version", activation.PropertyVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	logger.Debugf("Done")
	return nil
}

func activateProperty(property *papi.Property, d *schema.ResourceData, correlationid string, logger log.Interface) (*papi.Activation, error) {
	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return nil, err
	}
	activation, err := getActivation(d, property, papi.ActivationTypeActivate, papi.NetworkValue(network), correlationid, logger)
	if err != nil {
		return nil, err
	}

	if foundActivation, err := findExistingActivation(property, activation, correlationid, logger); err == nil && foundActivation != nil {
		return foundActivation, nil
	}
	if err = activation.Save(property, true); err != nil {
		body, err := json.Marshal(activation)
		if err != nil {
			logger.Errorf("marshaling error: %s", err)
		}
		logger.Debugf("API Request Body: %s", string(body))
		return nil, err
	}
	logger.Debugf("Activation submitted successfully")
	return activation, nil
}

func deactivateProperty(property *papi.Property, d *schema.ResourceData, network papi.NetworkValue, correlationid string, logger log.Interface) (*papi.Activation, error) {
	version, err := property.GetLatestVersion(network, correlationid)
	if err != nil || version == nil {
		// Not active
		return nil, nil
	}

	activation, err := getActivation(d, property, papi.ActivationTypeDeactivate, network, correlationid, logger)
	if err != nil {
		return nil, err
	}

	if foundActivation, err := findExistingActivation(property, activation, correlationid, logger); err == nil && foundActivation != nil {
		return foundActivation, nil
	}

	if err = activation.Save(property, true); err != nil {
		body, err := json.Marshal(activation)
		if err != nil {
			logger.Errorf("marshaling error: %s", err)
		}
		logger.Debugf("API Request Body: %s\n", string(body))
		return nil, err
	}
	logger.Debugf("Deactivation submitted successfully")
	return activation, nil
}

func getActivation(d *schema.ResourceData, property *papi.Property, activationType papi.ActivationValue, network papi.NetworkValue, correlationid string, logger log.Interface) (*papi.Activation, error) {
	logger.Debugf("Creating new activation")
	activation := papi.NewActivation(papi.NewActivations())
	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if !errors.Is(err, tools.ErrNotFound) && version != 0 {
		activation.PropertyVersion = version
	} else {
		version, err := property.GetLatestVersion("", correlationid)
		if err != nil {
			return nil, err
		}
		logger.Debugf("Using latest version: %d", version.PropertyVersion)
		activation.PropertyVersion = version.PropertyVersion
	}
	activation.Network = network
	contact, err := tools.GetSetValue("contact", d)
	if err != nil {
		return nil, err
	}
	for _, email := range contact.List() {
		emailStr, ok := email.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "email", "string")
		}
		activation.NotifyEmails = append(activation.NotifyEmails, emailStr)
	}
	activation.Note = "Using Terraform"
	activation.ActivationType = activationType

	logger.Debugf("Activating")
	return activation, nil
}

func findExistingActivation(property *papi.Property, activation *papi.Activation, _ string, logger log.Interface) (*papi.Activation, error) {
	activations, err := property.GetActivations()
	if err != nil {
		return nil, err
	}
	inProgressStates := map[papi.StatusValue]bool{
		papi.StatusActive:              true,
		papi.StatusNew:                 true,
		papi.StatusPending:             true,
		papi.StatusPendingDeactivation: true,
		papi.StatusZone1:               true,
		papi.StatusZone2:               true,
		papi.StatusZone3:               true,
	}

	for _, a := range activations.Activations.Items {
		if _, ok := inProgressStates[a.Status]; !ok {
			continue
		}

		// There is an activation in progress, if it's for the same version/network/type we can re-use it
		if a.PropertyVersion == activation.PropertyVersion && a.ActivationType == activation.ActivationType && a.Network == activation.Network {
			logger.Debugf("An Existing activation found %v Activations values %s version %d Activationpassed in  %s version %d",
				inProgressStates[a.Status], a.Network, a.PropertyVersion, activation.Network, activation.PropertyVersion)
			return a, nil
		}
	}
	return nil, nil
}

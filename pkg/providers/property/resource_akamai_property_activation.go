package property

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyActivationCreate,
		ReadContext:   resourcePropertyActivationRead,
		UpdateContext: resourcePropertyActivationUpdate,
		DeleteContext: resourcePropertyActivationDelete,
		Schema:        akamaiPropertyActivationSchema,
		Timeouts: &schema.ResourceTimeout{
			Default: &PropertyResourceTimeout,
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

	// PropertyResourceTimeout is the default timeout for the resource operations
	PropertyResourceTimeout = time.Minute * 90
)

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
		Default:  "STAGING",
	},
	"activate": {
		Type:       schema.TypeBool,
		Optional:   true,
		Default:    true,
		Deprecated: "the activate flag has been deprecated, in future activation will always be performed",
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

func resourcePropertyActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationCreate")
	client := inst.Client(meta)

	log.Debug("resourcePropertyActivationCreate call")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	if dead, ok := ctx.Deadline(); ok {
		logger.Debugf("activation create with deadline in %s", dead.Sub(time.Now()).String())
	}

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done - activate=false")
		return nil
	}

	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if version == 0 {
		// get the property - so we can determine latest version
		property, err := client.GetProperty(ctx, papi.GetPropertyRequest{
			PropertyID: propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		// use the latest version for the property
		version = property.Property.LatestVersion
		logger.Debugf("Version missing during create - computed as %+v", version)
	}

	// check to see if this tree has any issues
	rules, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ValidateRules:   true,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// if there are errors return them cleanly
	if len(rules.Errors) > 0 {
		diags := make([]diag.Diagnostic, 0)

		for _, e := range rules.Errors {
			logger.Warnf("property rule error %s", e.Error())

			// handle errors with no title since summary is required field
			errorSummary := e.Title
			if len(errorSummary) == 0 {
				errorSummary = "Papi error message shown below"
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  errorSummary,
				Detail:   e.Error(),
			})
		}

		return diags
	}

	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := lookupActivation(ctx, client, lookupActivationRequest{
		propertyID: propertyID,
		version:    version,
		network:    papi.ActivationNetwork(network),
		activationType: map[papi.ActivationType]struct{}{
			papi.ActivationTypeActivate:   {},
			papi.ActivationTypeDeactivate: {},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// we create a new property activation in case of no previous activation, or deleted activation
	if activation == nil || activation.ActivationType == papi.ActivationTypeDeactivate {
		notifySet, err := tools.GetSetValue("contact", d)
		if err != nil {
			return diag.FromErr(err)
		}
		var notify []string
		for _, contact := range notifySet.List() {
			notify = append(notify, cast.ToString(contact))
		}

		create, err := client.CreateActivation(ctx, papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeActivate,
				Network:                papi.ActivationNetwork(network),
				PropertyVersion:        version,
				NotifyEmails:           notify,
				AcknowledgeAllWarnings: true,
			},
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("create activation failed: %w", err))
		}

		// query the activation to retreive the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: create.ActivationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activation = act.Activation
	}

	d.SetId(activation.ActivationID)

	for activation.Status != papi.ActivationStatusActive {
		if activation.Status == papi.ActivationStatusAborted {
			return diag.FromErr(fmt.Errorf("activation request aborted"))
		}
		if activation.Status == papi.ActivationStatusFailed {
			return diag.FromErr(fmt.Errorf("activation request failed in downstream system"))
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, papi.GetActivationRequest{
				ActivationID: activation.ActivationID,
				PropertyID:   propertyID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act.Activation

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return diag.Diagnostics{DiagWarnActivationTimeout}
			} else if errors.Is(ctx.Err(), context.Canceled) {
				return diag.Diagnostics{DiagWarnActivationCanceled}
			}
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	if err := d.Set("version", activation.PropertyVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("status", string(activation.Status)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourcePropertyActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "resourcePropertyActivationDelete")
	client := inst.Client(meta)

	log.Debug("resourcePropertyActivationDelete call")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		log.Debugf("Done")
		return nil
	}

	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// get the property version
	resp, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	version := resp.Versions.Items[0].PropertyVersion

	activation, err := lookupActivation(ctx, client, lookupActivationRequest{
		propertyID: propertyID,
		version:    version,
		network:    papi.ActivationNetwork(network),
		activationType: map[papi.ActivationType]struct{}{
			papi.ActivationTypeDeactivate: {},
			papi.ActivationTypeActivate:   {},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if activation == nil || activation.ActivationType == papi.ActivationTypeActivate {
		notifySet, err := tools.GetSetValue("contact", d)
		if err != nil {
			return diag.FromErr(err)
		}
		var notify []string
		for _, contact := range notifySet.List() {
			notify = append(notify, cast.ToString(contact))
		}

		deleteActivation, err := client.CreateActivation(ctx, papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeDeactivate,
				Network:                papi.ActivationNetwork(network),
				PropertyVersion:        version,
				NotifyEmails:           notify,
				AcknowledgeAllWarnings: true,
			},
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("create deactivation failed: %w", err))
		}
		// update with id we are now polling on
		d.SetId(deleteActivation.ActivationID)

		// query the activation to retreive the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: deleteActivation.ActivationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activation = act.Activation
	}

	// deactivations also use status Active for when they are fully processed
	for activation.Status != papi.ActivationStatusActive {
		if activation.Status == papi.ActivationStatusAborted {
			return diag.FromErr(fmt.Errorf("deactivation request aborted"))
		}
		if activation.Status == papi.ActivationStatusFailed {
			return diag.FromErr(fmt.Errorf("deactivation request failed in downstream system"))
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, papi.GetActivationRequest{
				ActivationID: activation.ActivationID,
				PropertyID:   propertyID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act.Activation

		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	d.SetId("")

	return nil
}

func resourcePropertyActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "resourcePropertyActivationRead")
	client := inst.Client(meta)

	log.Debug("resourcePropertyActivationRead call")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if version == 0 {
		// get the property - so we can determine latest version
		property, err := client.GetProperty(ctx, papi.GetPropertyRequest{
			PropertyID: propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		// use the latest version for the property
		version = property.Property.LatestVersion
		log.Debugf("Version missing for read - computed as %+v", version)
	}

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get activations for property: %w", err))
	}

	for _, act := range resp.Activations.Items {

		if act.Network == network && act.PropertyVersion == version {
			log.Debugf("Found Existing Activation %s version %d", network, version)

			if err := d.Set("status", string(act.Status)); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
			if err := d.Set("version", act.PropertyVersion); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
			d.SetId(act.ActivationID)

			break
		}
	}

	return nil
}

func resourcePropertyActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationUpdate")
	client := inst.Client(meta)

	log.Debug("resourcePropertyActivationUpdate call")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done")
		return nil
	}

	propertyID, err := tools.GetStringValue("property", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if version == 0 {
		// get the property - so we can determine latest version
		property, err := client.GetProperty(ctx, papi.GetPropertyRequest{
			PropertyID: propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		// use the latest version for the property
		version = property.Property.LatestVersion
		logger.Debugf("Version missing for update - computed as %+v", version)
	}

	// check to see if this tree has any issues
	rules, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ValidateRules:   true,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// if there are errors return them cleanly
	if len(rules.Errors) > 0 {
		diags := make([]diag.Diagnostic, 0)

		for _, e := range rules.Errors {
			logger.Warnf("property rule error %s", e.Error())

			// handle errors with no title since summary is required field
			errorSummary := e.Title
			if len(errorSummary) == 0 {
				errorSummary = "Papi error message shown below"
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  errorSummary,
				Detail:   e.Error(),
			})
		}

		return diags
	}

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyActivation, err := lookupActivation(ctx, client, lookupActivationRequest{
		propertyID: propertyID,
		version:    version,
		network:    network,
		activationType: map[papi.ActivationType]struct{}{
			papi.ActivationTypeActivate: {},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if propertyActivation == nil {
		notifySet, err := tools.GetSetValue("contact", d)
		if err != nil {
			return diag.FromErr(err)
		}
		var notify []string
		for _, contact := range notifySet.List() {
			notify = append(notify, cast.ToString(contact))
		}

		create, err := client.CreateActivation(ctx, papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeActivate,
				Network:                network,
				PropertyVersion:        version,
				NotifyEmails:           notify,
				AcknowledgeAllWarnings: true,
			},
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("create activation failed: %w", err))
		}

		// query the activation to retrieve the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: create.ActivationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		propertyActivation = act.Activation
	}

	d.SetId(propertyActivation.ActivationID)

	for propertyActivation.Status != papi.ActivationStatusActive {
		if propertyActivation.Status == papi.ActivationStatusAborted {
			return diag.FromErr(fmt.Errorf("activation request aborted"))
		}
		if propertyActivation.Status == papi.ActivationStatusFailed {
			return diag.FromErr(fmt.Errorf("activation request failed in downstream system"))
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, papi.GetActivationRequest{
				ActivationID: propertyActivation.ActivationID,
				PropertyID:   propertyID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			propertyActivation = act.Activation

		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	if err := d.Set("version", propertyActivation.PropertyVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("status", string(propertyActivation.Status)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

type lookupActivationRequest struct {
	propertyID     string
	version        int
	network        papi.ActivationNetwork
	activationType map[papi.ActivationType]struct{}
}

// returns the most recent property activation (by SubmitDate) for the given query.
func lookupActivation(ctx context.Context, client papi.PAPI, query lookupActivationRequest) (*papi.Activation, error) {
	activations, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: query.propertyID,
	})
	if err != nil {
		return nil, err
	}

	inProgressStates := map[papi.ActivationStatus]bool{
		papi.ActivationStatusActive:       true,
		papi.ActivationStatusNew:          true,
		papi.ActivationStatusPending:      true,
		papi.ActivationStatusDeactivating: true,
		papi.ActivationStatusZone1:        true,
		papi.ActivationStatusZone2:        true,
		papi.ActivationStatusZone3:        true,
	}

	var bestMatch *papi.Activation
	var bestMatchSubmitDate time.Time

	for _, a := range activations.Activations.Items {
		if _, ok := inProgressStates[a.Status]; !ok {
			continue
		}

		// There is an activation in progress, if it's for the same version/network/type we can re-use it
		_, matchingActivationType := query.activationType[a.ActivationType]
		if a.PropertyVersion == query.version && matchingActivationType && a.Network == query.network {
			// find the most recent activation

			var aSubmitDate, err = tools.ParseDate(tools.DateTimeFormat, a.SubmitDate)
			if err != nil {
				return nil, err
			}

			if bestMatchSubmitDate.Equal(time.Time{}) || bestMatchSubmitDate.Before(aSubmitDate) {
				bestMatch = a
				bestMatchSubmitDate, err = tools.ParseDate(tools.DateTimeFormat, bestMatch.SubmitDate)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	// only return activation if bestMatch.ActivationType == query.ActivationType

	if bestMatch != nil {
		_, matchingActivationType := query.activationType[bestMatch.ActivationType]
		if matchingActivationType {
			return bestMatch, nil
		}
	}
	return nil, nil
}

func networkAlias(d *schema.ResourceData) (papi.ActivationNetwork, error) {
	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return "", err
	}

	networks := map[string]papi.ActivationNetwork{
		"STAGING":    papi.ActivationNetworkStaging,
		"STAG":       papi.ActivationNetworkStaging,
		"S":          papi.ActivationNetworkStaging,
		"PRODUCTION": papi.ActivationNetworkProduction,
		"PROD":       papi.ActivationNetworkProduction,
		"P":          papi.ActivationNetworkProduction,
	}
	networkValue, ok := networks[strings.ToUpper(network)]
	if !ok {
		return "", fmt.Errorf("network not recognized")
	}
	return networkValue, nil
}

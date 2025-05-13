package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
)

func resourcePropertyActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyActivationCreate,
		ReadContext:   resourcePropertyActivationRead,
		UpdateContext: resourcePropertyActivationUpdate,
		DeleteContext: resourcePropertyActivationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyActivationImport,
		},
		Schema: akamaiPropertyActivationSchema,
		Timeouts: &schema.ResourceTimeout{
			Default: &PropertyResourceTimeout,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourcePropertyActivationV0().CoreConfigSchema().ImpliedType(),
			Upgrade: timeouts.MigrateToExplicit(),
		}},
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

	// CreateActivationRetry poll wait time code waits between retries for activation creation
	CreateActivationRetry = 10 * time.Second
)

var akamaiPropertyActivationSchema = map[string]*schema.Schema{
	"property_id": {
		Type:        schema.TypeString,
		Required:    true,
		StateFunc:   addPrefixToState("prp_"),
		Description: "Your property's ID, including the prp_ prefix.",
	},
	"activation_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The ID given to the activation event while it's in progress.",
	},
	"errors": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Errors returned during activation.",
	},
	"warnings": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Warnings returned during activation.",
	},
	"rule_errors": {
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        papiError(),
		Description: "Any errors returned by the API about rules.",
	},
	"auto_acknowledge_rule_warnings": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Automatically acknowledge all rule warnings for activation to continue. Default is false",
	},
	"version": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: tf.IsNotBlank,
		Description:      "Your property's version number.",
	},
	"network": {
		Type:        schema.TypeString,
		Optional:    true,
		Default:     papi.ActivationNetworkStaging,
		Description: "Akamai network in which to activate your property, either STAGING or PRODUCTION. The default is STAGING.",
	},
	"contact": {
		Type:             schema.TypeSet,
		Required:         true,
		Elem:             &schema.Schema{Type: schema.TypeString},
		DiffSuppressFunc: suppressDiffIfNoPropertyReactivation,
		Description:      "One or more email addresses to which to send activation status changes.",
	},
	"status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The property version's activation status on the given network.",
	},
	"note": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Assigns a log message to the activation request.",
		DiffSuppressFunc: suppressDiffIfNoPropertyReactivation,
	},
	"compliance_record": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Provides an audit record when activating on a production network.",
		Elem:        complianceRecordSchema,
	},
	"timeouts": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Enables to set timeout for processing.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: timeouts.ValidateDurationFormat,
				},
			},
		},
	},
}

func papiError() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"type":           {Type: schema.TypeString, Optional: true},
		"title":          {Type: schema.TypeString, Optional: true},
		"detail":         {Type: schema.TypeString, Optional: true},
		"instance":       {Type: schema.TypeString, Optional: true},
		"behavior_name":  {Type: schema.TypeString, Optional: true},
		"error_location": {Type: schema.TypeString, Optional: true},
		"status_code":    {Type: schema.TypeInt, Optional: true},
	}}
}

func resourcePropertyActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationCreate")
	client := Client(meta)

	logger.Debug("resourcePropertyActivationCreate call")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	if dead, ok := ctx.Deadline(); ok {
		logger.Debugf("activation create with deadline in %s", time.Until(dead).String())
	}

	propertyID, err := resolvePropertyID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := resolveVersion(ctx, d, client, propertyID, network)
	if err != nil {
		return diag.FromErr(err)
	}
	// Schema guarantees these types
	acknowledgeRuleWarnings := d.Get("auto_acknowledge_rule_warnings").(bool)

	// check to see if this tree has any issues
	rules, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ValidateRules:   true,
	})
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	// if there are errors return them cleanly
	diags := checkRuleTreeErrorsAndWarnings(rules, d, logger)
	if diags != nil && diags.HasError() {
		d.Partial(true)
		return diags
	}

	complianceRecord, err := tf.GetListValue("compliance_record", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	activation, err := lookupActivation(ctx, client, lookupActivationRequest{
		propertyID: propertyID,
		network:    network,
		activationType: map[papi.ActivationType]struct{}{
			papi.ActivationTypeActivate:   {},
			papi.ActivationTypeDeactivate: {},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// we create a new property activation in case of no previous activation, or deleted activation
	if activation == nil || activation.ActivationType == papi.ActivationTypeDeactivate || activation.PropertyVersion != version {
		contactSet, err := tf.GetSetValue("contact", d)
		if err != nil {
			return diag.FromErr(err)
		}
		var contacts []string
		for _, contact := range contactSet.List() {
			contacts = append(contacts, cast.ToString(contact))
		}

		note, err := tf.GetStringValue("note", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		createActivationRequest := papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeActivate,
				Network:                network,
				PropertyVersion:        version,
				NotifyEmails:           contacts,
				AcknowledgeAllWarnings: acknowledgeRuleWarnings,
				Note:                   note,
			},
		}

		logger.Debug("creating activation")
		activationID, diagErr := createActivation(ctx, client, addPropertyComplianceRecord(complianceRecord, createActivationRequest))
		if diagErr != nil {
			return diagErr
		}

		// query the activation to retrieve the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: activationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activation = act.Activation

		if err = setErrorsAndWarnings(d, flattenErrorArray(act.Errors), flattenErrorArray(act.Warnings)); err != nil {
			return diag.FromErr(err)
		}
	}

	activation, diagErr := pollActivation(ctx, client, activation, propertyID)
	if diagErr != nil {
		return diagErr
	}

	attrs := map[string]interface{}{
		"status":        string(activation.Status),
		"activation_id": activation.ActivationID,
		"version":       version,
	}
	if err := tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(propertyID + ":" + string(network))

	return nil
}

func resourcePropertyActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationDelete")
	client := Client(meta)

	logger.Debug("resourcePropertyActivationDelete call")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyID, err := resolvePropertyID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	version, err := resolveVersion(ctx, d, client, propertyID, network)
	if err != nil {
		return diag.FromErr(err)
	}

	complianceRecord, err := tf.GetListValue("compliance_record", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	// Schema guarantees these types
	acknowledgeRuleWarnings := d.Get("auto_acknowledge_rule_warnings").(bool)

	activation, err := lookupActivation(ctx, client, lookupActivationRequest{
		propertyID: propertyID,
		version:    version,
		network:    network,
		activationType: map[papi.ActivationType]struct{}{
			papi.ActivationTypeDeactivate: {},
			papi.ActivationTypeActivate:   {},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if activation == nil || activation.ActivationType == papi.ActivationTypeActivate {
		contactSet, err := tf.GetRawSetValue("contact", d, tf.NewRawConfig(d))
		if err != nil {
			return diag.FromErr(err)
		}
		var contacts []string
		for _, contact := range contactSet {
			contacts = append(contacts, cast.ToString(contact))
		}
		note, err := tf.GetStringValue("note", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		deleteActivationRequest := papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeDeactivate,
				Network:                network,
				PropertyVersion:        version,
				NotifyEmails:           contacts,
				AcknowledgeAllWarnings: acknowledgeRuleWarnings,
				Note:                   note,
			},
		}

		deleteActivationID, diagErr := createActivation(ctx, client, addPropertyComplianceRecord(complianceRecord, deleteActivationRequest))
		if diagErr != nil {
			return diagErr
		}
		// update with id we are now polling on
		d.SetId(deleteActivationID)

		// query the activation to retrieve the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: deleteActivationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activation = act.Activation

		if err := d.Set("activation_id", activation.ActivationID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}

		if err = setErrorsAndWarnings(d, flattenErrorArray(act.Errors), flattenErrorArray(act.Warnings)); err != nil {
			return diag.FromErr(err)
		}
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
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, papi.GetActivationRequest{
				ActivationID: activation.ActivationID,
				PropertyID:   propertyID,
			})
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act.Activation

			if err = setErrorsAndWarnings(d, flattenErrorArray(act.Errors), flattenErrorArray(act.Warnings)); err != nil {
				return diag.FromErr(err)
			}

		case <-ctx.Done():
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	d.SetId("")

	return nil
}

func flattenErrorArray(errors []*papi.Error) string {
	var errorStrArr = make([]string, len(errors))
	for i, err := range errors {
		strError := str.From(err.Error())
		errorStrArr[i] = strError
	}
	return strings.Join(errorStrArr, "\n")
}

func resourcePropertyActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationRead")
	client := Client(meta)

	logger.Debug("resourcePropertyActivationRead call")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	propertyID, err := resolvePropertyID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
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

	if err := d.Set("errors", flattenErrorArray(resp.Errors)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}
	if err := d.Set("warnings", flattenErrorArray(resp.Warnings)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	activation, err := findLatestActive(resp.Activations.Items, network)
	if err != nil && !errors.Is(err, errNoActiveVersionFound) {
		return diag.Errorf("unexpected error searching for latest activation: %s", err)
	}
	if errors.Is(err, errNoActiveVersionFound) {
		d.SetId("")
		return nil
	}

	attrs := map[string]interface{}{
		"status":        string(activation.Status),
		"version":       activation.PropertyVersion,
		"network":       network,
		"activation_id": activation.ActivationID,
		"note":          activation.Note,
		"contact":       activation.NotifyEmails,
	}

	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(activation.PropertyID + ":" + string(network))

	return nil
}

var errNoActiveVersionFound = errors.New("activation not found")

func findLatestActive(activations []*papi.Activation, network papi.ActivationNetwork) (*papi.Activation, error) {
	if len(activations) == 0 {
		return nil, errNoActiveVersionFound
	}

	sort.Slice(activations, func(i, j int) bool {
		return activations[i].UpdateDate > activations[j].UpdateDate
	})

	for _, activation := range activations {
		if activation.ActivationType == papi.ActivationTypeActivate &&
			activation.Network == network &&
			activation.Status == papi.ActivationStatusActive {
			return activation, nil
		}
		if activation.ActivationType == papi.ActivationTypeDeactivate &&
			activation.Network == network &&
			(activation.Status == papi.ActivationStatusActive) {
			return nil, errNoActiveVersionFound
		}
	}

	return nil, errNoActiveVersionFound
}

func resolveVersion(ctx context.Context, d *schema.ResourceData, client papi.PAPI, propertyID string, network papi.ActivationNetwork) (int, error) {

	version, err := tf.GetIntValue("version", d)
	if err != nil {
		var resp *papi.GetPropertyVersionsResponse
		// use the latest version for the property
		resp, err = client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID:  propertyID,
			ActivatedOn: fmt.Sprintf("%v", network),
		})
		if err != nil {
			return 0, err
		}
		version = resp.Version.PropertyVersion
	}

	return version, nil
}

func resourcePropertyActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationUpdate")
	client := Client(meta)

	logger.Debug("resourcePropertyActivationUpdate call")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	if !d.HasChangesExcept("timeouts", "compliance_record") {
		logger.Debug("Only timeouts and/or compliance_record were updated, update with no API calls")
		return nil
	}

	propertyID, err := resolvePropertyID(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	network, err := networkAlias(d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := resolveVersion(ctx, d, client, propertyID, network)
	if err != nil {
		return diag.FromErr(err)
	}

	complianceRecord, err := tf.GetListValue("compliance_record", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	// Schema guarantees these types
	acknowledgeRuleWarnings := d.Get("auto_acknowledge_rule_warnings").(bool)

	// Assigns a log message to the activation request
	note, err := tf.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	// check to see if this tree has any issues
	rules, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ValidateRules:   true,
	})
	if err != nil {
		// Reverting to previous state(property version in this case) when error occurs.
		d.Partial(true)
		return diag.FromErr(err)
	}

	// if there are errors return them cleanly
	diags := checkRuleTreeErrorsAndWarnings(rules, d, logger)
	if diags.HasError() {
		d.Partial(true)
		return diags
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

	versionStatus, err := resolveVersionStatus(ctx, client, propertyID, version, network)
	if err != nil {
		return diag.FromErr(err)
	}

	if versionStatus == papi.VersionStatusActive {
		var updatedFields []string

		if d.HasChange("auto_acknowledge_rule_warnings") {
			updatedFields = append(updatedFields, "'auto_acknowledge_rule_warnings'")
		}
		if len(updatedFields) > 0 {
			return diag.Errorf("Cannot update %s field(s) while property version is ACTIVE. Deactivate the current version to update, or create a new property version activation.", strings.Join(updatedFields, ", "))
		}
	}

	if propertyActivation == nil || versionStatus == papi.VersionStatusDeactivated {
		contactSet, err := tf.GetRawSetValue("contact", d, tf.NewRawConfig(d))
		if err != nil {
			return diag.FromErr(err)
		}
		var contacts []string
		for _, contact := range contactSet {
			contacts = append(contacts, cast.ToString(contact))
		}

		createActivationRequest := papi.CreateActivationRequest{
			PropertyID: propertyID,
			Activation: papi.Activation{
				ActivationType:         papi.ActivationTypeActivate,
				Network:                network,
				PropertyVersion:        version,
				NotifyEmails:           contacts,
				AcknowledgeAllWarnings: acknowledgeRuleWarnings,
				Note:                   note,
			},
		}

		activationID, diagErr := createActivation(ctx, client, addPropertyComplianceRecord(complianceRecord, createActivationRequest))
		if diagErr != nil {
			return diagErr
		}

		// query the activation to retrieve the initial status
		act, err := client.GetActivation(ctx, papi.GetActivationRequest{
			ActivationID: activationID,
			PropertyID:   propertyID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		propertyActivation = act.Activation

		if err = setErrorsAndWarnings(d, flattenErrorArray(act.Errors), flattenErrorArray(act.Warnings)); err != nil {
			return diag.FromErr(err)
		}
	}

	propertyActivation, diagErr := pollActivation(ctx, client, propertyActivation, propertyID)
	if diagErr != nil {
		return diagErr
	}

	attrs := map[string]interface{}{
		"status":        string(propertyActivation.Status),
		"activation_id": propertyActivation.ActivationID,
		"version":       version,
	}
	if err := tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(propertyID + ":" + string(network))

	return nil
}

func resourcePropertyActivationImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyActivationImport")

	logger.Debug("Importing property activation")

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid property activation identifier: %s", d.Id())
	}

	attrs := make(map[string]interface{}, 3)
	attrs["property_id"] = parts[0]
	attrs["network"] = parts[1]
	attrs["auto_acknowledge_rule_warnings"] = false

	if err := tf.SetAttrs(d, attrs); err != nil {
		return nil, err
	}

	// Errors are checked during property activation, and a property cannot activate if any errors are present.
	// As a result, rule_errors will always be empty.
	// Initialize rule_errors as an empty list if it is not already set in the state.
	if _, exists := d.GetOk("rule_errors"); !exists {
		if err := d.Set("rule_errors", []interface{}{}); err != nil {
			return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
	}

	return []*schema.ResourceData{d}, nil
}

func addPropertyComplianceRecord(complianceRecord []interface{}, activatePAPIRequest papi.CreateActivationRequest) papi.CreateActivationRequest {
	if len(complianceRecord) == 0 {
		return activatePAPIRequest
	}
	crMap := complianceRecord[0].(map[string]interface{})

	if len(crMap["noncompliance_reason_none"].([]interface{})) != 0 {
		complianceRecordNone := &papi.ComplianceRecordNone{}
		if crMap["noncompliance_reason_none"].([]interface{})[0] != nil {
			crNoneMap := crMap["noncompliance_reason_none"].([]interface{})[0].(map[string]interface{})
			complianceRecordNone = &papi.ComplianceRecordNone{
				CustomerEmail:  crNoneMap["customer_email"].(string),
				PeerReviewedBy: crNoneMap["peer_reviewed_by"].(string),
				TicketID:       crNoneMap["ticket_id"].(string),
				UnitTested:     crNoneMap["unit_tested"].(bool),
			}
		}
		activatePAPIRequest.Activation.ComplianceRecord = complianceRecordNone
	} else if len(crMap["noncompliance_reason_other"].([]interface{})) != 0 {
		complianceRecordOther := &papi.ComplianceRecordOther{}
		if crMap["noncompliance_reason_other"].([]interface{})[0] != nil {
			crOtherMap := crMap["noncompliance_reason_other"].([]interface{})[0].(map[string]interface{})
			complianceRecordOther = &papi.ComplianceRecordOther{
				TicketID:                 crOtherMap["ticket_id"].(string),
				OtherNoncomplianceReason: crOtherMap["other_noncompliance_reason"].(string),
			}
		}
		activatePAPIRequest.Activation.ComplianceRecord = complianceRecordOther
	} else if len(crMap["noncompliance_reason_no_production_traffic"].([]interface{})) != 0 {
		complianceRecordNoProductionTraffic := &papi.ComplianceRecordNoProductionTraffic{}
		if crMap["noncompliance_reason_no_production_traffic"].([]interface{})[0] != nil {
			crNoProdTrafficMap := crMap["noncompliance_reason_no_production_traffic"].([]interface{})[0].(map[string]interface{})
			complianceRecordNoProductionTraffic = &papi.ComplianceRecordNoProductionTraffic{
				TicketID: crNoProdTrafficMap["ticket_id"].(string),
			}
		}
		activatePAPIRequest.Activation.ComplianceRecord = complianceRecordNoProductionTraffic
	} else if len(crMap["noncompliance_reason_emergency"].([]interface{})) != 0 {
		complianceRecordEmergency := &papi.ComplianceRecordEmergency{}
		if crMap["noncompliance_reason_emergency"].([]interface{})[0] != nil {
			crEmergencyMap := crMap["noncompliance_reason_emergency"].([]interface{})[0].(map[string]interface{})
			complianceRecordEmergency = &papi.ComplianceRecordEmergency{
				TicketID: crEmergencyMap["ticket_id"].(string),
			}
		}
		activatePAPIRequest.Activation.ComplianceRecord = complianceRecordEmergency
	}

	return activatePAPIRequest
}

func resolveVersionStatus(ctx context.Context, client papi.PAPI, propertyID string, version int, network papi.ActivationNetwork) (papi.VersionStatus, error) {
	var versionStatus papi.VersionStatus
	propertyVersion, err := client.GetPropertyVersion(ctx, papi.GetPropertyVersionRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
	})
	if err != nil {
		return "", err
	}
	if network == papi.ActivationNetworkProduction {
		versionStatus = propertyVersion.Version.ProductionStatus
	} else {
		versionStatus = propertyVersion.Version.StagingStatus
	}
	return versionStatus, nil
}

func checkRuleTreeErrorsAndWarnings(rules *papi.GetRuleTreeResponse, d *schema.ResourceData, logger log.Interface) diag.Diagnostics {
	var diags diag.Diagnostics

	// Set rule_errors only if actual rule errors are present
	if len(rules.Errors) > 0 {
		if err := d.Set("rule_errors", papiErrorsToList(rules.Errors)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}
		msg, err := json.MarshalIndent(papiErrorsToList(rules.Errors), "", "\t")
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling API error: %s", err))
		}
		logger.Errorf("Property has rule errors %s", msg)
		diags = append(diags, diag.Errorf("activation cannot continue due to rule errors: %s", msg)...)
	} else {
		if err := d.Set("rule_errors", []interface{}{}); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}
	}

	// Handle warnings
	if len(rules.Warnings) > 0 {
		msg, err := json.MarshalIndent(papiErrorsToList(rules.Warnings), "", "\t")
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling API warnings: %s", err))
		}
		logger.Warnf("Property has rule warnings %s", msg)
	}
	return diags
}

func setErrorsAndWarnings(d *schema.ResourceData, errors, warnings string) error {
	if err := d.Set("errors", errors); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("warnings", warnings); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resolvePropertyID(d *schema.ResourceData) (string, error) {
	propertyID, err := tf.GetStringValue("property_id", d)
	return str.AddPrefix(propertyID, "prp_"), err
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

	inProgressStates := map[papi.ActivationStatus]struct{}{
		papi.ActivationStatusActive:       {},
		papi.ActivationStatusNew:          {},
		papi.ActivationStatusPending:      {},
		papi.ActivationStatusDeactivating: {},
		papi.ActivationStatusZone1:        {},
		papi.ActivationStatusZone2:        {},
		papi.ActivationStatusZone3:        {},
	}

	var bestMatch *papi.Activation
	var bestMatchSubmitDate time.Time

	for _, a := range activations.Activations.Items {
		if _, ok := inProgressStates[a.Status]; !ok {
			continue
		}

		// There is an activation in progress, if it's for the same version/network/type we can re-use it
		_, matchingActivationType := query.activationType[a.ActivationType]
		// If query doesn't check version, it should not filter any activation
		matchingVersion := query.version == 0 || a.PropertyVersion == query.version
		matchingNetwork := a.Network == query.network
		if matchingVersion && matchingActivationType && matchingNetwork {
			// find the most recent activation

			var aSubmitDate, err = date.Parse(a.SubmitDate)
			if err != nil {
				return nil, err
			}

			if bestMatchSubmitDate.IsZero() || bestMatchSubmitDate.Before(aSubmitDate) {
				bestMatch = a
				bestMatchSubmitDate = aSubmitDate
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
	network, err := tf.GetStringValue("network", d)
	if err != nil {
		if errors.Is(err, tf.ErrNotFound) {
			network = "STAGING"
		} else {
			return "", err
		}
	}

	alias, err := NetworkAlias(network)
	if err != nil {
		return "", err
	}

	return papi.ActivationNetwork(alias), nil
}

func pollActivation(ctx context.Context, client papi.PAPI, activation *papi.Activation, propertyID string) (*papi.Activation, diag.Diagnostics) {

	retriesMax := 5
	retries5xx := 0

	for activation.Status != papi.ActivationStatusActive {
		if activation.Status == papi.ActivationStatusAborted {
			return nil, diag.FromErr(fmt.Errorf("activation request aborted"))
		}
		if activation.Status == papi.ActivationStatusFailed {
			return nil, diag.FromErr(fmt.Errorf("activation request failed in downstream system"))
		}
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivation(ctx, papi.GetActivationRequest{
				ActivationID: activation.ActivationID,
				PropertyID:   propertyID,
			})
			if err != nil {
				var target = &papi.Error{}
				if !errors.As(err, &target) {
					return nil, diag.Errorf("error has unexpected type: %T", err)
				}
				if target.StatusCode >= 500 {
					retries5xx = retries5xx + 1
					if retries5xx > retriesMax {
						return nil, diag.Errorf("reached max number of 5xx retries: %d", retries5xx)
					}
					continue
				}

				return nil, diag.FromErr(err)
			}
			retries5xx = 0
			activation = act.Activation

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, diag.Diagnostics{DiagErrActivationTimeout}
			} else if errors.Is(ctx.Err(), context.Canceled) {
				return nil, diag.Diagnostics{DiagErrActivationCanceled}
			}
			return nil, diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}
	return activation, nil
}

func suppressDiffIfNoPropertyReactivation(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if d.Id() == "" {
		return false
	}
	shouldReactivate := d.HasChanges("version", "network") || tf.StringFieldHasChangesWithStateFunc(d, "property_id", addPrefixToState("prp_"))
	if !shouldReactivate {
		return true
	}
	return oldValue == newValue
}

func createActivation(ctx context.Context, client papi.PAPI, request papi.CreateActivationRequest) (string, diag.Diagnostics) {
	log := hclog.FromContext(ctx)

	errMsg := "create failed"
	switch request.Activation.ActivationType {
	case papi.ActivationTypeActivate:
		errMsg = "create activation failed"
	case papi.ActivationTypeDeactivate:
		errMsg = "create deactivation failed"
	}

	createActivationRetry := CreateActivationRetry

	for {
		log.Debug("creating activation")
		create, err := client.CreateActivation(ctx, request)
		if err == nil {
			return create.ActivationID, nil
		}
		log.Debug("%s: retrying: %w", errMsg, err)

		if !isCreateActivationErrorRetryable(err) {
			return "", diag.Errorf("%s: %s", errMsg, err)
		}

		if actID, ok := isActivationPendingOrActive(ctx, client, expectedActivation{
			PropertyID: request.PropertyID,
			Version:    request.Activation.PropertyVersion,
			Network:    request.Activation.Network,
			Type:       request.Activation.ActivationType,
		}); ok {
			return actID, nil
		}

		select {
		case <-time.After(createActivationRetry):
			createActivationRetry = capDuration(createActivationRetry*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return "", diag.Diagnostics{DiagErrActivationTimeout}
			} else if errors.Is(ctx.Err(), context.Canceled) {
				return "", diag.Diagnostics{DiagErrActivationCanceled}
			}
			return "", diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}
}

func capDuration(t time.Duration, tMax time.Duration) time.Duration {
	if t > tMax {
		return tMax
	}
	return t
}

func isCreateActivationErrorRetryable(err error) bool {
	if errors.Is(err, io.EOF) {
		return true
	}

	var responseErr = &papi.Error{}
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

type expectedActivation struct {
	PropertyID string
	Version    int
	Network    papi.ActivationNetwork
	Type       papi.ActivationType
}

// isActivationPendingOrActive check if latest activation is of specified version and has status Pending or Active
func isActivationPendingOrActive(ctx context.Context, client papi.PAPI, expected expectedActivation) (string, bool) {
	log := hclog.FromContext(ctx)

	log.Debug("getting activation")
	acts, err := client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: expected.PropertyID,
	})
	if err != nil {
		return "", false
	}
	activations := acts.Activations.Items

	sort.Slice(activations, func(i, j int) bool {
		return activations[i].UpdateDate > activations[j].UpdateDate
	})

	activations = filterActivationsByNetwork(activations, expected.Network)

	if len(activations) == 0 { // job might be scheduled but no activation created yet (unlikely)
		log.Debug("no activation items; retrying")
		return "", false
	}
	latestActivationItem := activations[0] // grab the latest one returned by api

	if latestActivationItem.PropertyVersion != expected.Version {
		log.Debug("latest version mismatch; retrying")
		return "", false
	}
	if latestActivationItem.ActivationType != expected.Type {
		log.Debug("activation type mismatch; retrying")
		return "", false
	}
	if latestActivationItem.Status == papi.ActivationStatusPending ||
		latestActivationItem.Status == papi.ActivationStatusActive {
		return latestActivationItem.ActivationID, true
	}
	return "", false
}

func filterActivationsByNetwork(activations []*papi.Activation, network papi.ActivationNetwork) []*papi.Activation {
	var filteredActivations []*papi.Activation
	for _, activation := range activations {
		if activation.Network == network {
			filteredActivations = append(filteredActivations, activation)
		}
	}

	return filteredActivations
}

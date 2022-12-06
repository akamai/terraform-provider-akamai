package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePropertyIncludeActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyIncludeActivationCreate,
		ReadContext:   resourcePropertyIncludeActivationRead,
		UpdateContext: resourcePropertyIncludeActivationUpdate,
		DeleteContext: resourcePropertyIncludeActivationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyIncludeActivationImport,
		},
		Schema: map[string]*schema.Schema{
			"include_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   addPrefixToState("inc_"),
				Description: "The unique identifier of the include",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   addPrefixToState("ctr_"),
				Description: "The contract under which the include is activated",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc:   addPrefixToState("grp_"),
				Description: "The group under which the include is activated",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the include",
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(papi.ActivationNetworkStaging), string(papi.ActivationNetworkProduction),
				}, false)),
				Description: "The network for which the activation will be performed",
			},
			"notify_emails": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of email addresses to notify about an activation status",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The note to assign to a log message of the activation request",
			},
			"auto_acknowledge_rule_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically acknowledge all rule warnings for activation and continue",
			},
			"validations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The validation information in JSON format",
			},
			"compliance_record": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Provides an audit record when activating on a production network",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"noncompliance_reason": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      fmt.Sprintf("Specifies the reason for the expedited activation on production network. Valid noncompliance reasons are: %s", strings.Join(validComplianceRecords, ", ")),
							ValidateDiagFunc: tools.ValidateStringInSlice(validComplianceRecords),
						},
						"ticket_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifies the ticket that describes the need for the activation",
						},
						"other_noncompliance_reason": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Describes the reason why the activation must occur immediately, out of compliance with the standard procedure",
						},
						"customer_email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifies the customer",
						},
						"peer_reviewed_by": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifies person who has independently approved the activation request",
						},
						"unit_tested": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the metadata to activate has been fully tested",
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &includeActivationTimeout,
		},
	}
}

var (
	activationPollInterval   = time.Minute
	includeActivationTimeout = time.Minute * 30
	getActivationInterval    = time.Second * 5
	validComplianceRecords   = []string{papi.NoncomplianceReasonNone, papi.NoncomplianceReasonOther, papi.NoncomplianceReasonNoProductionTraffic, papi.NoncomplianceReasonEmergency}
)

func resourcePropertyIncludeActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Create property include activation")

	err := resourcePropertyIncludeActivationUpsert(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourcePropertyIncludeActivationRead(ctx, d, m)
}

func resourcePropertyIncludeActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading property include activation")

	id := strings.Split(d.Id(), ":")
	if len(id) < 4 {
		return diag.Errorf("invalid include activation identifier: %s", d.Id())
	}
	contractID, groupID, includeID, network := id[0], id[1], id[2], id[3]

	versions, err := client.ListIncludeActivations(ctx, papi.ListIncludeActivationsRequest{
		IncludeID:  includeID,
		GroupID:    groupID,
		ContractID: contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	activationID, err := getLatestIncludeActivationID(versions, network)
	if err != nil {
		return diag.FromErr(err)
	}

	activation, err := waitForPropertyIncludeOperation(ctx, client, activationID, includeID, "activation")
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})
	if activation != nil {
		var validations []byte
		if activation.Validations != nil {
			validations, err = json.Marshal(activation.Validations)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		attrs["include_id"] = activation.Activations.Items[0].IncludeID
		attrs["contract_id"] = activation.ContractID
		attrs["group_id"] = activation.GroupID
		attrs["version"] = activation.Activations.Items[0].IncludeVersion
		attrs["network"] = activation.Activations.Items[0].Network
		attrs["notify_emails"] = activation.Activations.Items[0].NotifyEmails
		attrs["note"] = activation.Activations.Items[0].Note
		attrs["validations"] = string(validations)

		// it is impossible to fetch compliance_record and auto_acknowledge_rule_warnings attributes from server
		if err = tools.SetAttrs(d, attrs); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourcePropertyIncludeActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeActivationUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Updating property include activation")

	mutableAttrsHaveChanges := d.HasChanges("note", "notify_emails", "auto_acknowledge_rule_warnings", "compliance_record")
	if mutableAttrsHaveChanges && !d.HasChanges("version") {
		return diag.FromErr(fmt.Errorf("attributes such as 'note', 'notify_emails', 'auto_acknowledge_rule_warnings', " +
			"'compliance_record' cannot be updated after resource creation without 'version' attribute modification"))
	}

	err := resourcePropertyIncludeActivationUpsert(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourcePropertyIncludeActivationRead(ctx, d, m)
}

func resourcePropertyIncludeActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeActivationDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Deactivating property include")

	includeID, err := tools.GetStringValue("include_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	includeID = tools.AddPrefix(includeID, "inc_")
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notifyEmailsSet, err := tools.GetSetValue("notify_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notifyEmails := tools.SetToStringSlice(notifyEmailsSet)
	note, err := tools.GetStringValue("note", d)
	if err != nil {
		return diag.FromErr(err)
	}
	acknowledgement, err := tools.GetBoolValue("auto_acknowledge_rule_warnings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	complianceRecord, err := tools.GetListValue("compliance_record", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	deactivateIncludeRequest := papi.DeactivateIncludeRequest{
		IncludeID:              includeID,
		Version:                version,
		Network:                papi.ActivationNetwork(network),
		Note:                   note,
		NotifyEmails:           notifyEmails,
		AcknowledgeAllWarnings: acknowledgement,
	}
	deactivateIncludeRequest, err = addComplianceRecordToDeactivationByNetwork(network, complianceRecord, deactivateIncludeRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	deactivation, err := client.DeactivateInclude(ctx, deactivateIncludeRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = waitForPropertyIncludeOperation(ctx, client, deactivation.ActivationID, includeID, "deactivation")
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePropertyIncludeActivationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeActivationImport")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	logger.Debug("Importing property include activation")

	id := strings.Split(d.Id(), ":")
	if len(id) < 4 {
		return nil, fmt.Errorf("invalid include activation identifier: %s", d.Id())
	}
	contractID, groupID, includeID, network := id[0], id[1], id[2], id[3]

	if contractID == "" || groupID == "" || includeID == "" || network == "" {
		return nil, fmt.Errorf("contract, group, include IDs and network must have non empty values")
	}

	// it is impossible to fetch auto_acknowledge_rule_warnings from server
	if err := d.Set("auto_acknowledge_rule_warnings", false); err != nil {
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	return []*schema.ResourceData{d}, nil
}

func resourcePropertyIncludeActivationUpsert(ctx context.Context, d *schema.ResourceData, client papi.PAPI) error {
	includeID, err := tools.GetStringValue("include_id", d)
	if err != nil {
		return err
	}
	includeID = tools.AddPrefix(includeID, "inc_")
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return err
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	groupID, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return err
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return err
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil {
		return err
	}
	notifyEmailsSet, err := tools.GetSetValue("notify_emails", d)
	if err != nil {
		return err
	}
	notifyEmails := tools.SetToStringSlice(notifyEmailsSet)
	note, err := tools.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	acknowledgement, err := tools.GetBoolValue("auto_acknowledge_rule_warnings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}
	complianceRecord, err := tools.GetListValue("compliance_record", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}

	activateIncludeRequest := papi.ActivateIncludeRequest{
		IncludeID:              includeID,
		Version:                version,
		Network:                papi.ActivationNetwork(network),
		Note:                   note,
		NotifyEmails:           notifyEmails,
		AcknowledgeAllWarnings: acknowledgement,
	}

	activateIncludeRequest, err = addComplianceRecordToActivationByNetwork(network, complianceRecord, activateIncludeRequest)
	if err != nil {
		return err
	}
	activationRes, err := client.ActivateInclude(ctx, activateIncludeRequest)
	if err != nil {
		return err
	}

	// here is used temporary activationID
	if _, err := waitForPropertyIncludeOperation(ctx, client, activationRes.ActivationID, includeID, "activation"); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s:%s:%s", contractID, groupID, includeID, network))
	return nil
}

// returns stable include activation id instead a temporary one
func getLatestIncludeActivationID(versions *papi.ListIncludeActivationsResponse, network string) (string, error) {
	activations := filterIncludeActivationsByNetwork(versions.Activations.Items, network)
	act, err := findLatestIncludeActivation(activations)
	if err != nil {
		return "", err
	}
	return act.ActivationID, nil
}

// waitForPropertyIncludeOperation adds timeout for activation and deactivation include operations
func waitForPropertyIncludeOperation(ctx context.Context, client papi.PAPI, activationID, includeID, operationType string) (*papi.GetIncludeActivationResponse, error) {
	activation, err := client.GetIncludeActivation(ctx, papi.GetIncludeActivationRequest{
		IncludeID:    includeID,
		ActivationID: activationID,
	})
	if err != nil {
		// it can take a few seconds to fetch include activation/deactivation right after activation/deactivation request
		if strings.Contains(err.Error(), papi.ErrNotFound.Error()) && strings.Contains(err.Error(), papi.ErrGetIncludeActivation.Error()) {
			select {
			case <-time.After(getActivationInterval):
				return waitForPropertyIncludeOperation(ctx, client, activationID, includeID, operationType)
			case <-ctx.Done():
				return nil, terminateProcess(ctx, operationType)
			}
		}
		return nil, err
	}
	for activation != nil && activation.Activations.Items[0].Status != papi.ActivationStatusActive {
		actStatus := activation.Activations.Items[0].Status

		if actStatus == papi.ActivationStatusFailed {
			return nil, fmt.Errorf("%s request failed for property include %v", operationType, includeID)
		}
		if actStatus == papi.ActivationStatusAborted {
			return nil, fmt.Errorf("pending %s request aborted for property include %v", operationType, includeID)
		}

		select {
		case <-time.After(activationPollInterval):
			activation, err = client.GetIncludeActivation(ctx, papi.GetIncludeActivationRequest{
				IncludeID:    includeID,
				ActivationID: activationID,
			})
			if err != nil {
				// it can take a few seconds to fetch include activation/deactivation right after activation/deactivation request
				if strings.Contains(err.Error(), papi.ErrNotFound.Error()) && strings.Contains(err.Error(), papi.ErrGetIncludeActivation.Error()) {
					select {
					case <-time.After(getActivationInterval):
						return waitForPropertyIncludeOperation(ctx, client, activationID, includeID, operationType)
					}
				}
				return nil, err
			}
			if err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, terminateProcess(ctx, operationType)
		}
	}
	return activation, nil
}

func terminateProcess(ctx context.Context, operationType string) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("timeout waiting for %s status", operationType)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("operation canceled while waiting for %s status", operationType)
	}
	return fmt.Errorf("%s context terminated: %w", operationType, ctx.Err())
}

func addComplianceRecordToActivationByNetwork(network string, complianceRecord []interface{}, activateIncludeRequest papi.ActivateIncludeRequest) (papi.ActivateIncludeRequest, error) {
	result, err := addComplianceRecordByNetwork(network, "activation", complianceRecord, papi.ActivateOrDeactivateIncludeRequest(activateIncludeRequest))
	return papi.ActivateIncludeRequest(result), err
}

func addComplianceRecordToDeactivationByNetwork(network string, complianceRecord []interface{}, deactivateIncludeRequest papi.DeactivateIncludeRequest) (papi.DeactivateIncludeRequest, error) {
	result, err := addComplianceRecordByNetwork(network, "deactivation", complianceRecord, papi.ActivateOrDeactivateIncludeRequest(deactivateIncludeRequest))
	return papi.DeactivateIncludeRequest(result), err
}

// all the validations for compliance_record attributes is performed in AkamaiOPEN-edgegrid-golang
func addComplianceRecordByNetwork(network, operation string, complianceRecord []interface{}, activateIncludeRequest papi.ActivateOrDeactivateIncludeRequest) (papi.ActivateOrDeactivateIncludeRequest, error) {
	if papi.ActivationNetwork(network) == papi.ActivationNetworkProduction && len(complianceRecord) == 0 {
		return activateIncludeRequest, fmt.Errorf("compliance_record field is required for '%v' network to %s include version", papi.ActivationNetworkProduction, operation)
	}
	activateIncludeRequest = addComplianceRecord(complianceRecord, activateIncludeRequest)
	return activateIncludeRequest, nil
}

func addComplianceRecord(complianceRecord []interface{}, activateIncludeRequest papi.ActivateOrDeactivateIncludeRequest) papi.ActivateOrDeactivateIncludeRequest {
	if len(complianceRecord) == 0 {
		return activateIncludeRequest
	}

	crMap := complianceRecord[0].(map[string]interface{})
	noncomplianceReason := crMap["noncompliance_reason"].(string)
	ticketID := crMap["ticket_id"].(string)
	otherNoncomplianceReason := crMap["other_noncompliance_reason"].(string)
	customerEmail := crMap["customer_email"].(string)
	peerReviewedBy := crMap["peer_reviewed_by"].(string)
	unitTested := crMap["unit_tested"].(bool)

	switch noncomplianceReason {
	case papi.NoncomplianceReasonOther:
		complianceRecordOther := &papi.ComplianceRecordOther{
			TicketID:                 ticketID,
			OtherNoncomplianceReason: otherNoncomplianceReason,
		}
		activateIncludeRequest.ComplianceRecord = complianceRecordOther
	case papi.NoncomplianceReasonNone:
		complianceRecordNone := &papi.ComplianceRecordNone{
			CustomerEmail:  customerEmail,
			PeerReviewedBy: peerReviewedBy,
			TicketID:       ticketID,
			UnitTested:     unitTested,
		}
		activateIncludeRequest.ComplianceRecord = complianceRecordNone
	case papi.NoncomplianceReasonNoProductionTraffic:
		complianceRecordNoProductionTraffic := &papi.ComplianceRecordNoProductionTraffic{
			TicketID: ticketID,
		}
		activateIncludeRequest.ComplianceRecord = complianceRecordNoProductionTraffic
	case papi.NoncomplianceReasonEmergency:
		complianceRecordEmergency := &papi.ComplianceRecordEmergency{
			TicketID: ticketID,
		}
		activateIncludeRequest.ComplianceRecord = complianceRecordEmergency
	}
	return activateIncludeRequest
}

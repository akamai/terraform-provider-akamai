package papi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Activations is a collection of property activations
type Activations struct {
	client.Resource
	AccountID   string `json:"accountId"`
	ContractID  string `json:"contractId"`
	GroupID     string `json:"groupId"`
	Activations struct {
		Items []*Activation `json:"items"`
	} `json:"activations"`
}

// NewActivations creates a new Activations
func NewActivations() *Activations {
	activations := &Activations{}
	activations.Init()

	return activations
}

// GetActivations retrieves activation data for a given property
//
// See: Property.GetActivations()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listactivations
// Endpoint: GET /papi/v1/properties/{propertyId}/activations/{?contractId,groupId}
func (activations *Activations) GetActivations(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/papi/v1/properties/%s/activations?contractId=%s&groupId=%s",
			property.PropertyID,
			property.Contract.ContractID,
			property.Group.GroupID,
		),
		nil,
	)

	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, activations); err != nil {
		return err
	}

	return nil
}

// GetLatestProductionActivation retrieves the latest activation for the production network
//
// Pass in a status to check for, defaults to StatusActive
func (activations *Activations) GetLatestProductionActivation(status StatusValue) (*Activation, error) {
	return activations.GetLatestActivation(NetworkProduction, status)
}

// GetLatestStagingActivation retrieves the latest activation for the staging network
//
// Pass in a status to check for, defaults to StatusActive
func (activations *Activations) GetLatestStagingActivation(status StatusValue) (*Activation, error) {
	return activations.GetLatestActivation(NetworkStaging, status)
}

// GetLatestActivation gets the latest activation for the specified network
//
// Default to NetworkProduction. Pass in a status to check for, defaults to StatusActive
//
// This can return an activation OR a deactivation. Check activation.ActivationType and activation.Status for what you're looking for
func (activations *Activations) GetLatestActivation(network NetworkValue, status StatusValue) (*Activation, error) {
	if network == "" {
		network = NetworkProduction
	}

	if status == "" {
		status = StatusActive
	}

	var latest *Activation
	for _, activation := range activations.Activations.Items {
		if activation.Network == network && activation.Status == status && (latest == nil || activation.PropertyVersion > latest.PropertyVersion) {
			latest = activation
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("No activation found (network: %s, status: %s)", network, status)
	}

	return latest, nil
}

// Activation represents a property activation resource
type Activation struct {
	client.Resource
	parent              *Activations
	ActivationID        string                      `json:"activationId,omitempty"`
	ActivationType      ActivationValue             `json:"activationType,omitempty"`
	AcknowledgeWarnings []string                    `json:"acknowledgeWarnings,omitempty"`
	ComplianceRecord    *ActivationComplianceRecord `json:"complianceRecord,omitempty"`
	FastPush            bool                        `json:"fastPush,omitempty"`
	IgnoreHTTPErrors    bool                        `json:"ignoreHttpErrors,omitempty"`
	PropertyName        string                      `json:"propertyName,omitempty"`
	PropertyID          string                      `json:"propertyId,omitempty"`
	PropertyVersion     int                         `json:"propertyVersion"`
	Network             NetworkValue                `json:"network"`
	Status              StatusValue                 `json:"status,omitempty"`
	SubmitDate          string                      `json:"submitDate,omitempty"`
	UpdateDate          string                      `json:"updateDate,omitempty"`
	Note                string                      `json:"note,omitempty"`
	NotifyEmails        []string                    `json:"notifyEmails"`
	StatusChange        chan bool                   `json:"-"`
}

type ActivationComplianceRecord struct {
	NoncomplianceReason string `json:"noncomplianceReason,omitempty"`
}

// NewActivation creates a new Activation
func NewActivation(parent *Activations) *Activation {
	activation := &Activation{parent: parent}
	activation.Init()

	return activation
}

func (activation *Activation) Init() {
	activation.Complete = make(chan bool, 1)
	activation.StatusChange = make(chan bool, 1)
}

// GetActivation populates the Activation resource
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getanactivation
// Endpoint: GET /papi/v1/properties/{propertyId}/activations/{activationId}{?contractId,groupId}
func (activation *Activation) GetActivation(property *Property) (time.Duration, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/activations/%s",
			property.PropertyID,
			activation.ActivationID,
		),
		nil,
	)

	if err != nil {
		return 0, err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return 0, err
	}

	if client.IsError(res) {
		return 0, client.NewAPIError(res)
	}

	activations := NewActivations()
	if err := client.BodyJSON(res, activations); err != nil {
		return 0, err
	}

	activation.ActivationID = activations.Activations.Items[0].ActivationID
	activation.ActivationType = activations.Activations.Items[0].ActivationType
	activation.AcknowledgeWarnings = activations.Activations.Items[0].AcknowledgeWarnings
	activation.ComplianceRecord = activations.Activations.Items[0].ComplianceRecord
	activation.FastPush = activations.Activations.Items[0].FastPush
	activation.IgnoreHTTPErrors = activations.Activations.Items[0].IgnoreHTTPErrors
	activation.PropertyName = activations.Activations.Items[0].PropertyName
	activation.PropertyID = activations.Activations.Items[0].PropertyID
	activation.PropertyVersion = activations.Activations.Items[0].PropertyVersion
	activation.Network = activations.Activations.Items[0].Network
	activation.Status = activations.Activations.Items[0].Status
	activation.SubmitDate = activations.Activations.Items[0].SubmitDate
	activation.UpdateDate = activations.Activations.Items[0].UpdateDate
	activation.Note = activations.Activations.Items[0].Note
	activation.NotifyEmails = activations.Activations.Items[0].NotifyEmails

	//retry, _ := strconv.Atoi(res.Header.Get("Retry-After"))
	//retry *= int(time.Second)

	return time.Duration(30 * time.Second), nil
}

// Save activates a given property
//
// If acknowledgeWarnings is true and warnings are returned on the first attempt,
// a second attempt is made, acknowledging the warnings.
//
// See: Property.Activate()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#activateaproperty
// Endpoint: POST /papi/v1/properties/{propertyId}/activations/{?contractId,groupId}
func (activation *Activation) Save(property *Property, acknowledgeWarnings bool) error {
	if activation.ComplianceRecord == nil {
		activation.ComplianceRecord = &ActivationComplianceRecord{
			NoncomplianceReason: "NO_PRODUCTION_TRAFFIC",
		}
	}

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/papi/v1/properties/%s/activations",
			property.PropertyID,
		),
		activation,
	)

	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) && (!acknowledgeWarnings || (acknowledgeWarnings && res.StatusCode != 400)) {
		return client.NewAPIError(res)
	}

	if res.StatusCode == 400 && acknowledgeWarnings {
		warnings := &struct {
			Warnings []struct {
				Detail    string `json:"detail"`
				MessageID string `json:"messageId"`
			} `json:"warnings,omitempty"`
		}{}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(body, warnings); err != nil {
			return err
		}

		// Just in case we got a 400 for a different reason
		if len(warnings.Warnings) == 0 {
			jsonBody := &client.JSONBody{}

			if err = json.Unmarshal(body, jsonBody); err != nil {
				return err
			}

			return client.NewAPIErrorFromBody(res, body)
		}

		for _, warning := range warnings.Warnings {
			activation.AcknowledgeWarnings = append(activation.AcknowledgeWarnings, warning.MessageID)
		}

		// Don't acknowledgeWarnings again, halting a potential endless recursion
		return activation.Save(property, false)
	}

	var location client.JSONBody
	if err = client.BodyJSON(res, &location); err != nil {
		return err
	}

	req, err = client.NewRequest(
		Config,
		"GET",
		location["activationLink"].(string),
		nil,
	)

	if err != nil {
		return err
	}

	res, err = client.Do(Config, req)

	activations := NewActivations()
	if err := client.BodyJSON(res, activations); err != nil {
		return err
	}

	activation.ActivationID = activations.Activations.Items[0].ActivationID
	activation.ActivationType = activations.Activations.Items[0].ActivationType
	activation.AcknowledgeWarnings = activations.Activations.Items[0].AcknowledgeWarnings
	activation.ComplianceRecord = activations.Activations.Items[0].ComplianceRecord
	activation.FastPush = activations.Activations.Items[0].FastPush
	activation.IgnoreHTTPErrors = activations.Activations.Items[0].IgnoreHTTPErrors
	activation.PropertyName = activations.Activations.Items[0].PropertyName
	activation.PropertyID = activations.Activations.Items[0].PropertyID
	activation.PropertyVersion = activations.Activations.Items[0].PropertyVersion
	activation.Network = activations.Activations.Items[0].Network
	activation.Status = activations.Activations.Items[0].Status
	activation.SubmitDate = activations.Activations.Items[0].SubmitDate
	activation.UpdateDate = activations.Activations.Items[0].UpdateDate
	activation.Note = activations.Activations.Items[0].Note
	activation.NotifyEmails = activations.Activations.Items[0].NotifyEmails

	return nil
}

// PollStatus will responsibly poll till the property is active or an error occurs
//
// The Activation.StatusChange is a channel that can be used to
// block on status changes. If a new valid status is returned, true will
// be sent to the channel, otherwise, false will be sent.
//
//	go activation.PollStatus(property)
//	for activation.Status != edgegrid.StatusActive {
//		select {
//		case statusChanged := <-activation.StatusChange:
//			if statusChanged == false {
//				break
//			}
//		case <-time.After(time.Minute * 30):
//			break
//		}
//	}
//
//	if activation.Status == edgegrid.StatusActive {
//		// Activation succeeded
//	}
func (activation *Activation) PollStatus(property *Property) bool {
	currentStatus := activation.Status
	var retry time.Duration = 0

	for currentStatus != StatusActive {
		time.Sleep(retry)

		var err error
		retry, err = activation.GetActivation(property)

		if err != nil {
			activation.StatusChange <- false
			return false
		}

		if activation.Network == NetworkStaging && retry > time.Minute {
			retry = time.Minute
		}

		if err != nil {
			activation.StatusChange <- false
			return false
		}

		if currentStatus != activation.Status {
			currentStatus = activation.Status
			activation.StatusChange <- true
		}
	}

	return true
}

// Cancel an activation in progress
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#cancelapendingactivation
// Endpoint: DELETE /papi/v1/properties/{propertyId}/activations/{activationId}{?contractId,groupId}
func (activation *Activation) Cancel(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/papi/v1/properties/%s/activations",
			property.PropertyID,
		),
		nil,
	)

	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	newActivations := NewActivations()
	if err := client.BodyJSON(res, newActivations); err != nil {
		return err
	}

	activation.ActivationID = newActivations.Activations.Items[0].ActivationID
	activation.ActivationType = newActivations.Activations.Items[0].ActivationType
	activation.AcknowledgeWarnings = newActivations.Activations.Items[0].AcknowledgeWarnings
	activation.ComplianceRecord = newActivations.Activations.Items[0].ComplianceRecord
	activation.FastPush = newActivations.Activations.Items[0].FastPush
	activation.IgnoreHTTPErrors = newActivations.Activations.Items[0].IgnoreHTTPErrors
	activation.PropertyName = newActivations.Activations.Items[0].PropertyName
	activation.PropertyID = newActivations.Activations.Items[0].PropertyID
	activation.PropertyVersion = newActivations.Activations.Items[0].PropertyVersion
	activation.Network = newActivations.Activations.Items[0].Network
	activation.Status = newActivations.Activations.Items[0].Status
	activation.SubmitDate = newActivations.Activations.Items[0].SubmitDate
	activation.UpdateDate = newActivations.Activations.Items[0].UpdateDate
	activation.Note = newActivations.Activations.Items[0].Note
	activation.NotifyEmails = newActivations.Activations.Items[0].NotifyEmails
	activation.StatusChange = newActivations.Activations.Items[0].StatusChange

	return nil
}

// ActivationValue is used to create an "enum" of possible Activation.ActivationType values
type ActivationValue string

// NetworkValue is used to create an "enum" of possible Activation.Network values
type NetworkValue string

// StatusValue is used to create an "enum" of possible Activation.Status values
type StatusValue string

const (
	// ActivationTypeActivate Activation.ActivationType value ACTIVATE
	ActivationTypeActivate ActivationValue = "ACTIVATE"
	// ActivationTypeDeactivate Activation.ActivationType value DEACTIVATE
	ActivationTypeDeactivate ActivationValue = "DEACTIVATE"

	// NetworkProduction Activation.Network value PRODUCTION
	NetworkProduction NetworkValue = "PRODUCTION"
	// NetworkStaging Activation.Network value STAGING
	NetworkStaging NetworkValue = "STAGING"

	// StatusActive Activation.Status value ACTIVE
	StatusActive StatusValue = "ACTIVE"
	// StatusInactive Activation.Status value INACTIVE
	StatusInactive StatusValue = "INACTIVE"
	// StatusPending Activation.Status value PENDING
	StatusPending StatusValue = "PENDING"
	// StatusZone1 Activation.Status value ZONE_1
	StatusZone1 StatusValue = "ZONE_1"
	// StatusZone2 Activation.Status value ZONE_2
	StatusZone2 StatusValue = "ZONE_2"
	// StatusZone3 Activation.Status value ZONE_3
	StatusZone3 StatusValue = "ZONE_3"
	// StatusAborted Activation.Status value ABORTED
	StatusAborted StatusValue = "ABORTED"
	// StatusFailed Activation.Status value FAILED
	StatusFailed StatusValue = "FAILED"
	// StatusDeactivated Activation.Status value DEACTIVATED
	StatusDeactivated StatusValue = "DEACTIVATED"
	// StatusPendingDeactivation Activation.Status value PENDING_DEACTIVATION
	StatusPendingDeactivation StatusValue = "PENDING_DEACTIVATION"
	// StatusNew Activation.Status value NEW
	StatusNew StatusValue = "NEW"
)

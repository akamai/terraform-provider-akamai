package papi

import (
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// CustomBehaviors represents a collection of Custom Behaviors
//
// See: CustomBehaviors.GetCustomBehaviors()
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#custombehavior
type CustomBehaviors struct {
	client.Resource
	AccountID       string `json:"accountId"`
	CustomBehaviors struct {
		Items []*CustomBehavior `json:"items"`
	} `json:"customBehaviors"`
}

// NewCustomBehaviors creates a new *CustomBehaviors
func NewCustomBehaviors() *CustomBehaviors {
	return &CustomBehaviors{}
}

// PostUnmarshalJSON is called after UnmarshalJSON to setup the
// structs internal state. The cpcodes.Complete channel is utilized
// to communicate full completion.
func (behaviors *CustomBehaviors) PostUnmarshalJSON() error {
	behaviors.Init()

	for key, behavior := range behaviors.CustomBehaviors.Items {
		behaviors.CustomBehaviors.Items[key].parent = behaviors

		if err := behavior.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	return nil
}

// GetCustomBehaviors populates a *CustomBehaviors with it's related Custom Behaviors
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getcustombehaviors
// Endpoint: GET /papi/v1/custom-behaviors
func (behaviors *CustomBehaviors) GetCustomBehaviors() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/papi/v1/custom-behaviors",
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

	if err = client.BodyJSON(res, behaviors); err != nil {
		return err
	}

	return nil
}

func (behaviors *CustomBehaviors) AddCustomBehavior(behavior *CustomBehavior) {
	var exists bool
	for _, cb := range behaviors.CustomBehaviors.Items {
		if cb == behavior {
			exists = true
		}
	}

	if !exists {
		behaviors.CustomBehaviors.Items = append(behaviors.CustomBehaviors.Items, behavior)
	}
}

// CustomBehavior represents a single Custom Behavior
//
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#custombehavior
type CustomBehavior struct {
	client.Resource
	parent        *CustomBehaviors
	BehaviorID    string    `json:"behaviorId,omitempty"`
	Description   string    `json:"description"`
	DisplayName   string    `json:"displayName"`
	Name          string    `json:"name"`
	Status        string    `json:"status",omitempty`
	UpdatedByUser string    `json:"updatedByUser,omitempty"`
	UpdatedDate   time.Time `json:"updatedDate,omitempty"`
	XML           string    `json:"xml,omitempty"`
}

// GetCustomBehavior populates the *CustomBehavior with it's data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getcustombehavior
// Endpoint: GET /papi/v1/custom-behaviors/{behaviorId}
func (behavior *CustomBehavior) GetCustomBehavior() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/custom-behaviors/%s",
			behavior.BehaviorID,
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

	newCustomBehaviors := NewCustomBehaviors()
	if err = client.BodyJSON(res, newCustomBehaviors); err != nil {
		return err
	}
	if len(newCustomBehaviors.CustomBehaviors.Items) == 0 {
		return fmt.Errorf("Custom Behavior \"%s\" not found", behavior.BehaviorID)
	}

	behavior.Name = newCustomBehaviors.CustomBehaviors.Items[0].Name
	behavior.Description = newCustomBehaviors.CustomBehaviors.Items[0].Description
	behavior.DisplayName = newCustomBehaviors.CustomBehaviors.Items[0].DisplayName
	behavior.Status = newCustomBehaviors.CustomBehaviors.Items[0].Status
	behavior.UpdatedByUser = newCustomBehaviors.CustomBehaviors.Items[0].UpdatedByUser
	behavior.UpdatedDate = newCustomBehaviors.CustomBehaviors.Items[0].UpdatedDate
	behavior.XML = newCustomBehaviors.CustomBehaviors.Items[0].XML

	behavior.parent.AddCustomBehavior(behavior)

	return nil
}

// NewCustomBehavior creates a new *CustomBehavior
func NewCustomBehavior(behaviors *CustomBehaviors) *CustomBehavior {
	return &CustomBehavior{parent: behaviors}
}

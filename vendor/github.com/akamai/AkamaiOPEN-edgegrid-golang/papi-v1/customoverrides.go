package papi

import (
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// CustomOverrides represents a collection of Custom Overrides
//
// See: CustomerOverrides.GetCustomOverrides()
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#cpcode
type CustomOverrides struct {
	client.Resource
	AccountID       string `json:"accountId"`
	CustomOverrides struct {
		Items []*CustomOverride `json:"items"`
	} `json:"customOverrides"`
}

// NewCustomOverrides creates a new *CustomOverrides
func NewCustomOverrides() *CustomOverrides {
	return &CustomOverrides{}
}

// PostUnmarshalJSON is called after UnmarshalJSON to setup the
// structs internal state. The CustomOverrides.Complete channel is utilized
// to communicate full completion.
func (overrides *CustomOverrides) PostUnmarshalJSON() error {
	overrides.Init()

	for key, override := range overrides.CustomOverrides.Items {
		overrides.CustomOverrides.Items[key].parent = overrides

		if err := override.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	overrides.Complete <- true

	return nil
}

// GetCustomOverrides populates a *CustomOverrides with it's related Custom Overrides
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getcustomoverrides
// Endpoint: GET /papi/v1/custom-overrides
func (overrides *CustomOverrides) GetCustomOverrides() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/papi/v1/custom-overrides",
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

	if err = client.BodyJSON(res, overrides); err != nil {
		return err
	}

	return nil
}

func (overrides *CustomOverrides) AddCustomOverride(override *CustomOverride) {
	var exists bool
	for _, co := range overrides.CustomOverrides.Items {
		if co == override {
			exists = true
		}
	}

	if !exists {
		overrides.CustomOverrides.Items = append(overrides.CustomOverrides.Items, override)
	}
}

// CustomOverride represents a single Custom Override
//
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#customoverride
type CustomOverride struct {
	client.Resource
	parent        *CustomOverrides
	Description   string    `json:"description"`
	DisplayName   string    `json:"displayName"`
	Name          string    `json:"name"`
	OverrideID    string    `json:"overrideId,omitempty"`
	Status        string    `json:"status",omitempty`
	UpdatedByUser string    `json:"updatedByUser,omitempty"`
	UpdatedDate   time.Time `json:"updatedDate,omitempty"`
	XML           string    `json:"xml,omitempty"`
}

// GetCustomOverride populates the *CustomOverride with it's data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getcustomoverride
// Endpoint: GET /papi/v1/custom-overrides/{overrideId}
func (override *CustomOverride) GetCustomOverride() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/custom-overrides/%s",
			override.OverrideID,
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

	newCustomOverrides := NewCustomOverrides()
	if err = client.BodyJSON(res, newCustomOverrides); err != nil {
		return err
	}
	if len(newCustomOverrides.CustomOverrides.Items) == 0 {
		return fmt.Errorf("Custom Override \"%s\" not found", override.OverrideID)
	}

	override.Name = newCustomOverrides.CustomOverrides.Items[0].Name
	override.Description = newCustomOverrides.CustomOverrides.Items[0].Description
	override.DisplayName = newCustomOverrides.CustomOverrides.Items[0].DisplayName
	override.Status = newCustomOverrides.CustomOverrides.Items[0].Status
	override.UpdatedByUser = newCustomOverrides.CustomOverrides.Items[0].UpdatedByUser
	override.UpdatedDate = newCustomOverrides.CustomOverrides.Items[0].UpdatedDate
	override.XML = newCustomOverrides.CustomOverrides.Items[0].XML

	override.parent.AddCustomOverride(override)

	return nil
}

// NewCustomOverride creates a new *CustomOverride
func NewCustomOverride(overrides *CustomOverrides) *CustomOverride {
	return &CustomOverride{parent: overrides}
}

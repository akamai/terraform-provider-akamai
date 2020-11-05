package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Handle Operations on gtm resources
// Based on 1.4 schema
//

// Resources contains operations available on a Resource resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type Resources interface {
	// NewResourceInstance instantiates a new ResourceInstance.
	NewResourceInstance(context.Context, *Resource, int) *ResourceInstance
	// NewResource creates a new Resource object.
	NewResource(context.Context, string) *Resource
	// ListResources retreieves all Resources
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getresources
	ListResources(context.Context, string) ([]*Resource, error)
	// GetResource retrieves a Resource with the given name.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getresource
	GetResource(context.Context, string, string) (*Resource, error)
	// Create the datacenter identified by the receiver argument in the specified domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putresource
	CreateResource(context.Context, *Resource, string) (*ResourceResponse, error)
	// Delete the datacenter identified by the receiver argument from the domain specified.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#deleteresource
	DeleteResource(context.Context, *Resource, string) (*ResponseStatus, error)
	// Update the datacenter identified in the receiver argument in the provided domain.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putresource
	UpdateResource(context.Context, *Resource, string) (*ResponseStatus, error)
}

// ResourceInstance
type ResourceInstance struct {
	DatacenterId         int  `json:"datacenterId"`
	UseDefaultLoadObject bool `json:"useDefaultLoadObject"`
	LoadObject
}

// Resource represents a GTM resource
type Resource struct {
	Type                        string              `json:"type"`
	HostHeader                  string              `json:"hostHeader,omitempty"`
	LeastSquaresDecay           float64             `json:"leastSquaresDecay,omitempty"`
	Description                 string              `json:"description,omitempty"`
	LeaderString                string              `json:"leaderString,omitempty"`
	ConstrainedProperty         string              `json:"constrainedProperty,omitempty"`
	ResourceInstances           []*ResourceInstance `json:"resourceInstances,omitempty"`
	AggregationType             string              `json:"aggregationType,omitempty"`
	Links                       []*Link             `json:"links,omitempty"`
	LoadImbalancePercentage     float64             `json:"loadImbalancePercentage,omitempty"`
	UpperBound                  int                 `json:"upperBound,omitempty"`
	Name                        string              `json:"name"`
	MaxUMultiplicativeIncrement float64             `json:"maxUMultiplicativeIncrement,omitempty"`
	DecayRate                   float64             `json:"decayRate,omitempty"`
}

// ResourceList is the structure returned by List Resources
type ResourceList struct {
	ResourceItems []*Resource `json:"items"`
}

// Validate validates Resource
func (rsc *Resource) Validate() error {

	if len(rsc.Name) < 1 {
		return fmt.Errorf("Resource is missing Name")
	}
	if len(rsc.Type) < 1 {
		return fmt.Errorf("Resource is missing Type")
	}

	return nil
}

// NewResourceInstance instantiates a new ResourceInstance.
func (p *gtm) NewResourceInstance(ctx context.Context, rsrc *Resource, dcID int) *ResourceInstance {

	logger := p.Log(ctx)
	logger.Debug("NewResourceInstance")

	return &ResourceInstance{DatacenterId: dcID}

}

// NewResource creates a new Resource object.
func (p *gtm) NewResource(ctx context.Context, name string) *Resource {

	logger := p.Log(ctx)
	logger.Debug("NewResource")

	resource := &Resource{Name: name}
	return resource
}

// ListResources retreieves all Resources in the specified domain.
func (p *gtm) ListResources(ctx context.Context, domainName string) ([]*Resource, error) {

	logger := p.Log(ctx)
	logger.Debug("ListResources")

	var rsrcs ResourceList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListResources request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rsrcs)
	if err != nil {
		return nil, fmt.Errorf("ListResources request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return rsrcs.ResourceItems, nil
}

// GetResource retrieves a Resource with the given name in the specified domain.
func (p *gtm) GetResource(ctx context.Context, name, domainName string) (*Resource, error) {

	logger := p.Log(ctx)
	logger.Debug("GetResource")

	var rsc Resource
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetResource request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rsc)
	if err != nil {
		return nil, fmt.Errorf("GetResource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rsc, nil
}

// Create the resource identified by the receiver argument in the specified domain.
func (p *gtm) CreateResource(ctx context.Context, rsrc *Resource, domainName string) (*ResourceResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateResource")

	// Use common code. Any specific validation needed?
	return rsrc.save(ctx, p, domainName)

}

// Update the resourceidentified in the receiver argument in the specified domain.
func (p *gtm) UpdateResource(ctx context.Context, rsrc *Resource, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateResource")

	// common code
	stat, err := rsrc.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save Resource in given domain. Common path for Create and Update.
func (rsrc *Resource) save(ctx context.Context, p *gtm, domainName string) (*ResourceResponse, error) {

	if err := rsrc.Validate(); err != nil {
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resource request: %w", err)
	}

	var rscresp ResourceResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rscresp, rsrc)
	if err != nil {
		return nil, fmt.Errorf("Resource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &rscresp, nil

}

// Delete the resource identified in the receiver argument from the specified domain.
func (p *gtm) DeleteResource(ctx context.Context, rsrc *Resource, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteResource")

	if err := rsrc.Validate(); err != nil {
		logger.Errorf("Resource validation failed. %w", err)
		return nil, fmt.Errorf("Resource validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/resources/%s", domainName, rsrc.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Delete request: %w", err)
	}

	var rscresp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &rscresp)
	if err != nil {
		return nil, fmt.Errorf("Resource request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return rscresp.Status, nil

}

package papi

import (
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/patrickmn/go-cache"
)

var (
	Config       edgegrid.Config
	Profilecache = cache.New(5*time.Minute, 10*time.Minute)
)

// GetGroups retrieves all groups
func GetGroups() (*Groups, error) {
	groups := NewGroups()
	if err := groups.GetGroups(""); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetContracts retrieves all contracts
func GetContracts() (*Contracts, error) {
	contracts := NewContracts()
	if err := contracts.GetContracts(""); err != nil {
		return nil, err
	}

	return contracts, nil
}

// GetProducts retrieves all products
func GetProducts(contract *Contract) (*Products, error) {
	products := NewProducts()
	if err := products.GetProducts(contract, ""); err != nil {
		return nil, err
	}

	return products, nil
}

// GetEdgeHostnames retrieves all edge hostnames
func GetEdgeHostnames(contract *Contract, group *Group, options string) (*EdgeHostnames, error) {
	edgeHostnames := NewEdgeHostnames()
	if err := edgeHostnames.GetEdgeHostnames(contract, group, options, ""); err != nil {
		return nil, err
	}

	return edgeHostnames, nil
}

// GetCpCodes creates a new CpCodes struct and populates it with all CP Codes associated with a contract/group
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listcpcodes
func GetCpCodes(contract *Contract, group *Group) (*CpCodes, error) {
	cpcodes := NewCpCodes(contract, group)
	if err := cpcodes.GetCpCodes(""); err != nil {
		return nil, err
	}

	return cpcodes, nil
}

// GetProperties retrieves all properties for a given contract/group
func GetProperties(contract *Contract, group *Group) (*Properties, error) {
	properties := NewProperties()
	if err := properties.GetProperties(contract, group, ""); err != nil {
		return nil, err
	}

	return properties, nil
}

// GetVersions retrieves all versions for a given property
func GetVersions(property *Property) (*Versions, error) {
	versions := NewVersions()
	if err := versions.GetVersions(property, ""); err != nil {
		return nil, err
	}

	return versions, nil
}

// GetAvailableBehaviors retrieves all available behaviors for a property
func GetAvailableBehaviors(property *Property) (*AvailableBehaviors, error) {
	availableBehaviors := NewAvailableBehaviors()
	if err := availableBehaviors.GetAvailableBehaviors(property); err != nil {
		return nil, err
	}

	return availableBehaviors, nil
}

// GetAvailableCriteria retrieves all available criteria for a property
func GetAvailableCriteria(property *Property) (*AvailableCriteria, error) {
	availableCriteria := NewAvailableCriteria()
	if err := availableCriteria.GetAvailableCriteria(property); err != nil {
		return nil, err
	}

	return availableCriteria, nil
}

package papi

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Groups represents a collection of PAPI groups
type Groups struct {
	client.Resource
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Groups      struct {
		Items []*Group `json:"items"`
	} `json:"groups"`
}

// NewGroups creates a new Groups
func NewGroups() *Groups {
	groups := &Groups{}
	return groups
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (groups *Groups) PostUnmarshalJSON() error {
	groups.Init()
	for key, group := range groups.Groups.Items {
		groups.Groups.Items[key].parent = groups
		if err := group.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	groups.Complete <- true

	return nil
}

// GetGroups populates Groups with group data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listgroups
// Endpoint: GET /papi/v1/groups/
func (groups *Groups) GetGroups() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/papi/v1/groups",
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

	if err = client.BodyJSON(res, groups); err != nil {
		return err
	}

	return nil
}

// AddGroup adds a group to a Groups collection
func (groups *Groups) AddGroup(newGroup *Group) {
	if newGroup.GroupID != "" {
		for key, group := range groups.Groups.Items {
			if group.GroupID == newGroup.GroupID {
				groups.Groups.Items[key] = newGroup
				return
			}
		}
	}

	newGroup.parent = groups

	groups.Groups.Items = append(groups.Groups.Items, newGroup)
}

// FindGroup finds a specific group by ID
func (groups *Groups) FindGroup(id string) (*Group, error) {
	var group *Group
	var groupFound bool

	if id == "" {
		goto err
	}

	for _, group = range groups.Groups.Items {
		if group.GroupID == id {
			groupFound = true
			break
		}
	}

err:
	if !groupFound {
		return nil, fmt.Errorf("Unable to find group: \"%s\"", id)
	}

	return group, nil
}

// FindGroup finds a specific group by ID
func (groups *Groups) FindGroupId(name string) (*Group, error) {
	var group *Group
	var groupFound bool

	if name == "" {
		goto err
	}

	for _, group = range groups.Groups.Items {
		if group.GroupName == name {
			groupFound = true
			break
		}
	}

err:
	if !groupFound {
		return nil, fmt.Errorf("Unable to find group: \"%s\"", name)
	}

	return group, nil
}

// Group represents a group resource
type Group struct {
	client.Resource
	parent        *Groups
	GroupName     string   `json:"groupName"`
	GroupID       string   `json:"groupId"`
	ParentGroupID string   `json:"parentGroupId,omitempty"`
	ContractIDs   []string `json:"contractIds"`
}

// NewGroup creates a new Group
func NewGroup(parent *Groups) *Group {
	group := &Group{
		parent: parent,
	}
	group.Init()
	return group
}

// GetGroup populates a Group
func (group *Group) GetGroup() {
	groups, err := GetGroups()
	if err != nil {
		return
	}

	for _, g := range groups.Groups.Items {
		if g.GroupID == group.GroupID {
			group.parent = groups
			group.ContractIDs = g.ContractIDs
			group.GroupName = g.GroupName
			group.ParentGroupID = g.ParentGroupID
			group.Complete <- true
			return
		}
	}

	group.Complete <- false
}

// GetProperties retrieves all properties associated with a given group and contract
func (group *Group) GetProperties(contract *Contract) (*Properties, error) {
	return GetProperties(contract, group)
}

// GetCpCodes retrieves all CP codes associated with a given group and contract
func (group *Group) GetCpCodes(contract *Contract) (*CpCodes, error) {
	return GetCpCodes(contract, group)
}

// GetEdgeHostnames retrieves all Edge hostnames associated with a given group/contract
func (group *Group) GetEdgeHostnames(contract *Contract, options string) (*EdgeHostnames, error) {
	return GetEdgeHostnames(contract, group, options)
}

// NewProperty creates a property associated with a given group/contract
func (group *Group) NewProperty(contract *Contract) (*Property, error) {
	property := NewProperty(NewProperties())
	property.Contract = contract
	property.Group = group
	return property, nil
}

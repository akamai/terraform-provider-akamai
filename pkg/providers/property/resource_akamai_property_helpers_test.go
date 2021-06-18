package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/stretchr/testify/mock"
)

// Sets up an expected call to papi.GetGroups(), which returns the given parameters
func ExpectGetGroups(client *mockpapi, State *[]*papi.Group) *mock.Call {
	fn := func(ctx context.Context) (*papi.GetGroupsResponse, error) {
		var groups []*papi.Group

		for _, ptr := range *State {
			grp := *ptr
			groups = append(groups, &grp)
		}

		return &papi.GetGroupsResponse{Groups: papi.GroupItems{Items: groups}}, nil
	}

	return client.OnGetGroups(AnyCTX, fn)
}

// Sets up an expected call to papi.GetProperty() which returns a value depending on the given State pointer. When nil,
// the PAPI response contains a zero-value papi.Property. Otherwise the response will dynamically contain a copy of
// the State made at the time of the call to mockpapi.GetProperty().
func ExpectGetProperty(client *mockpapi, PropertyID, GroupID, ContractID string, State *papi.Property) *mock.Call {
	req := papi.GetPropertyRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}

	fn := func(context.Context, papi.GetPropertyRequest) (*papi.GetPropertyResponse, error) {
		var property papi.Property

		// Duplicate the State
		if State != nil {
			property = *State
		}

		// Duplicate the pointers
		if property.ProductionVersion != nil {
			v := *property.ProductionVersion
			property.ProductionVersion = &v
		}

		if property.StagingVersion != nil {
			v := *property.StagingVersion
			property.StagingVersion = &v
		}

		// although optional in PAPI documentation, ProductID is not being set by PAPI in the response
		property.ProductID = ""

		return &papi.GetPropertyResponse{Property: &property}, nil
	}

	return client.OnGetProperty(AnyCTX, req, fn)
}

// Sets up an expected call to papi.GetPropertyVersionHostnames() which returns a value depending on the value of the
// pointer to State. When nil or empty, the response contains a nil Items member. Otherwise the response contains a
// copy of the value pointed to by State made at the time of the call to papi.GetPropertyVersionHostnames().
func ExpectGetPropertyVersionHostnames(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, State *[]papi.Hostname) *mock.Call {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:        PropertyID,
		GroupID:           GroupID,
		ContractID:        ContractID,
		PropertyVersion:   PropertyVersion,
		IncludeCertStatus: true,
	}

	fn := func(context.Context, papi.GetPropertyVersionHostnamesRequest) (*papi.GetPropertyVersionHostnamesResponse, error) {
		var Items []papi.Hostname
		if len(*State) > 0 {
			// Duplicate the State
			Items = append(Items, *State...)
		}

		res := papi.GetPropertyVersionHostnamesResponse{
			ContractID:      ContractID,
			GroupID:         GroupID,
			PropertyID:      PropertyID,
			PropertyVersion: PropertyVersion,
			Hostnames:       papi.HostnameResponseItems{Items: Items},
		}

		return &res, nil
	}

	return client.OnGetPropertyVersionHostnames(AnyCTX, req, fn)
}

// Sets up an expected call to papi.UpdatePropertyVersionHostnames() which returns a constant value based on input
// params. If given, the value pointed to by State will be updated with a copy of the given Hostnames when the call
// to mockpapi.UpdatePropertyVersionHostnames() is made.
func ExpectUpdatePropertyVersionHostnames(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, Hostnames []papi.Hostname) *mock.Call {
	// func ExpectUpdatePropertyVersionHostnames(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, Hostnames []papi.Hostname, State *[]papi.Hostname) *mock.Call {
	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
		Hostnames:       Hostnames,
	}

	res := papi.UpdatePropertyVersionHostnamesResponse{
		ContractID:      ContractID,
		GroupID:         GroupID,
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		Hostnames:       papi.HostnameResponseItems{Items: Hostnames},
	}

	return client.On("UpdatePropertyVersionHostnames", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected call to papi.GetPropertyVersions()
func ExpectGetPropertyVersions(client *mockpapi, PropertyID, PropertyName, ContractID, GroupID string, versionItems *papi.PropertyVersionItems, err error) *mock.Call {
	req := papi.GetPropertyVersionsRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}
	var res *papi.GetPropertyVersionsResponse
	if err == nil {
		res = &papi.GetPropertyVersionsResponse{
			PropertyID:   PropertyID,
			PropertyName: PropertyName,
			ContractID:   ContractID,
			GroupID:      GroupID,
			Versions:     *versionItems,
		}
	}
	fn := func(context.Context, papi.GetPropertyVersionsRequest) (*papi.GetPropertyVersionsResponse, error) {
		res = &papi.GetPropertyVersionsResponse{
			PropertyID:   PropertyID,
			PropertyName: PropertyName,
			ContractID:   ContractID,
			GroupID:      GroupID,
			Versions:     *versionItems,
		}

		return res, nil
	}

	return client.OnGetPropertyVersions(AnyCTX, req, fn)
}

// Sets up an expected call to papi.GetPropertyVersion()
func ExpectGetPropertyVersion(client *mockpapi, PropertyID, GroupID, ContractID string, Version int, StagStatus, ProdStatus papi.VersionStatus) *mock.Call {
	req := papi.GetPropertyVersionRequest{
		PropertyID:      PropertyID,
		GroupID:         GroupID,
		ContractID:      ContractID,
		PropertyVersion: Version,
	}

	res := papi.GetPropertyVersionsResponse{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
		Version: papi.PropertyVersionGetItem{
			StagingStatus:    StagStatus,
			ProductionStatus: ProdStatus,
		},
	}
	return client.On("GetPropertyVersion", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected call to papi.CreatePropertyVersion()
func ExpectCreatePropertyVersion(client *mockpapi, PropertyID, GroupID, ContractID string, CreateFromVersion, NewVersion int) *mock.Call {
	req := papi.CreatePropertyVersionRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
		Version: papi.PropertyVersionCreate{
			CreateFromVersion: CreateFromVersion,
		},
	}

	res := papi.CreatePropertyVersionResponse{PropertyVersion: NewVersion}

	return client.On("CreatePropertyVersion", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected successful call to papi.CreateProperty() with a constant success response with the given PropertyID
func ExpectCreateProperty(client *mockpapi, PropertyName, GroupID, ContractID, ProductID, PropertyID string) *mock.Call {
	req := papi.CreatePropertyRequest{
		GroupID:    GroupID,
		ContractID: ContractID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: PropertyName,
		},
	}

	res := papi.CreatePropertyResponse{PropertyID: PropertyID}

	return client.On("CreateProperty", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected call to papi.RemoveProperty() with a constant success response
func ExpectRemoveProperty(client *mockpapi, PropertyID, ContractID, GroupID string) *mock.Call {
	req := papi.RemovePropertyRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
	}
	res := papi.RemovePropertyResponse{}

	return client.On("RemoveProperty", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected call to papi.GetRuleTree() which returns a value depending on the value of the
// pointer to State and FormatState.
func ExpectGetRuleTree(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, State *papi.RulesUpdate, RuleFormatState *string) *mock.Call {
	req := papi.GetRuleTreeRequest{
		PropertyID:      PropertyID,
		GroupID:         GroupID,
		ContractID:      ContractID,
		PropertyVersion: PropertyVersion,
		ValidateMode:    "full",
		ValidateRules:   true,
	}

	fn := func(context.Context, papi.GetRuleTreeRequest) (*papi.GetRuleTreeResponse, error) {
		var Rules papi.RulesUpdate
		if State != nil {
			Rules = *State
		}

		res := papi.GetRuleTreeResponse{
			PropertyID:      PropertyID,
			PropertyVersion: PropertyVersion,
			RuleFormat:      *RuleFormatState,
			Rules:           Rules.Rules,
		}

		return &res, nil
	}

	return client.OnGetRuleTree(AnyCTX, req, fn)
}

func ExpectUpdateRuleTree(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, State *papi.RulesUpdate, RuleFormat string, RuleError []papi.RuleError) *mock.Call {
	var RulesUpdate papi.RulesUpdate
	if State != nil {
		RulesUpdate = *State
	}
	var res papi.UpdateRulesResponse
	req := papi.UpdateRulesRequest{
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
		Rules:           RulesUpdate,
		ValidateRules:   true,
	}

	fn := func(context.Context, papi.UpdateRulesRequest) (*papi.UpdateRulesResponse, error) {

		res = papi.UpdateRulesResponse{
			PropertyID:      PropertyID,
			ContractID:      ContractID,
			GroupID:         GroupID,
			PropertyVersion: PropertyVersion,
			RuleFormat:      RuleFormat,
			Rules:           RulesUpdate.Rules,
			Errors:          RuleError,
		}
		return &res, nil
	}

	return client.OnUpdateRuleTree(AnyCTX, req, fn).Return(&res, nil)
}

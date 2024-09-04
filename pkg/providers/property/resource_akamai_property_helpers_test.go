package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/stretchr/testify/mock"
)

// Sets up an expected call to papi.GetGroups(), which returns the given parameters
func ExpectGetGroups(client *papi.Mock, State *[]*papi.Group) *mock.Call {
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
// the State made at the time of the call to papi.Mock.GetProperty().
func ExpectGetProperty(client *papi.Mock, PropertyID, GroupID, ContractID string, State *papi.Property) *mock.Call {
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
func ExpectGetPropertyVersionHostnames(client *papi.Mock, PropertyID, GroupID, ContractID string, PropertyVersion int, State *[]papi.Hostname) *mock.Call {
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
// to papi.Mock.UpdatePropertyVersionHostnames() is made.
func ExpectUpdatePropertyVersionHostnames(client *papi.Mock, PropertyID, GroupID, ContractID string, PropertyVersion int, Hostnames []papi.Hostname, err error) *mock.Call {
	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
		Hostnames:       Hostnames,
	}

	call := client.On("UpdatePropertyVersionHostnames", AnyCTX, req)
	if err != nil {
		return call.Return(&papi.UpdatePropertyVersionHostnamesResponse{}, err)
	}

	res := papi.UpdatePropertyVersionHostnamesResponse{
		ContractID:      ContractID,
		GroupID:         GroupID,
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		Hostnames:       papi.HostnameResponseItems{Items: Hostnames},
	}

	return call.Return(&res, nil)
}

// Sets up an expected call to papi.GetPropertyVersions()
func ExpectGetPropertyVersions(client *papi.Mock, PropertyID, PropertyName, ContractID, GroupID string, property *papi.Property, versionItems *papi.PropertyVersionItems) *mock.Call {
	req := papi.GetPropertyVersionsRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}
	var res *papi.GetPropertyVersionsResponse
	fn := func(context.Context, papi.GetPropertyVersionsRequest) (*papi.GetPropertyVersionsResponse, error) {
		if property != nil {
			ContractID = property.ContractID
			GroupID = property.GroupID
		}
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
func ExpectGetPropertyVersion(client *papi.Mock, PropertyID, GroupID, ContractID string, Version int, StagStatus, ProdStatus papi.VersionStatus) *mock.Call {
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
func ExpectCreatePropertyVersion(client *papi.Mock, PropertyID, GroupID, ContractID string, CreateFromVersion, NewVersion int) *mock.Call {
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
func ExpectCreateProperty(client *papi.Mock, PropertyName, GroupID, ContractID, ProductID, PropertyID string) *mock.Call {
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
func ExpectRemoveProperty(client *papi.Mock, PropertyID, ContractID, GroupID string) *mock.Call {
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
func ExpectGetRuleTree(client *papi.Mock, PropertyID, GroupID, ContractID string, PropertyVersion int, State *papi.RulesUpdate, RuleFormatState *string,
	Errors []*papi.Error, Warnings []*papi.Error) *mock.Call {
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
			Response: papi.Response{
				Errors:   Errors,
				Warnings: Warnings,
			},
			PropertyID:      PropertyID,
			PropertyVersion: PropertyVersion,
			RuleFormat:      *RuleFormatState,
			Rules:           Rules.Rules,
		}

		return &res, nil
	}

	return client.OnGetRuleTree(AnyCTX, req, fn)
}

func ExpectUpdateRuleTree(client *papi.Mock, PropertyID, GroupID, ContractID string, PropertyVersion int, State *papi.RulesUpdate, RuleFormat string, RuleError []papi.RuleError) *mock.Call {
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

func updateRuleTreeWithVariables(variables []papi.RuleVariable) *papi.RulesUpdate {
	return &papi.RulesUpdate{
		Rules: papi.Rules{
			Name: "default",
			Children: []papi.Rules{
				{
					Name: "change fwd path",
					Behaviors: []papi.RuleBehavior{
						{
							Name: "baseDirectory",
							Options: papi.RuleOptionsMap{
								"value": "/smth/",
							},
						},
					},
					Criteria: []papi.RuleBehavior{
						{
							Name:   "requestHeader",
							Locked: false,
							Options: papi.RuleOptionsMap{
								"headerName":              "Accept-Encoding",
								"matchCaseSensitiveValue": true,
								"matchOperator":           "IS_ONE_OF",
								"matchWildcardName":       false,
								"matchWildcardValue":      false,
							},
						},
					},
					CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll,
				},
				{
					Name: "caching",
					Behaviors: []papi.RuleBehavior{
						{
							Name: "caching",
							Options: papi.RuleOptionsMap{
								"behavior":       "MAX_AGE",
								"mustRevalidate": false,
								"ttl":            "1m",
							},
						},
					},
					CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAny,
				},
			},
			Behaviors: []papi.RuleBehavior{
				{
					Name: "origin",
					Options: papi.RuleOptionsMap{
						"cacheKeyHostname":          "REQUEST_HOST_HEADER",
						"compress":                  true,
						"enableTrueClientIp":        true,
						"forwardHostHeader":         "REQUEST_HOST_HEADER",
						"hostname":                  "test.domain",
						"httpPort":                  float64(80),
						"httpsPort":                 float64(443),
						"originCertificate":         "",
						"originSni":                 true,
						"originType":                "CUSTOMER",
						"ports":                     "",
						"trueClientIpClientSetting": false,
						"trueClientIpHeader":        "True-Client-IP",
						"verificationMode":          "PLATFORM_SETTINGS",
					},
				},
			},
			Options:   papi.RuleOptions{},
			Variables: variables,
			Comments:  "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
		},
	}
}

func updateRuleTreeWithVariablesStep1() *papi.RulesUpdate {
	return updateRuleTreeWithVariables([]papi.RuleVariable{
		{
			Name:        "TEST_EMPTY_FIELDS",
			Value:       ptr.To(""),
			Description: ptr.To(""),
			Hidden:      true,
			Sensitive:   false,
		},
		{
			Name:        "TEST_NIL_FIELD",
			Description: ptr.To(""),
			Value:       ptr.To(""),
			Hidden:      true,
			Sensitive:   false,
		},
	})
}

func updateRuleTreeWithVariablesStep0() *papi.RulesUpdate {
	return updateRuleTreeWithVariables([]papi.RuleVariable{
		{
			Name:        "TEST_EMPTY_FIELDS",
			Value:       ptr.To(""),
			Description: ptr.To(""),
			Hidden:      true,
			Sensitive:   false,
		},
		{
			Name:        "TEST_NIL_FIELD",
			Description: nil,
			Value:       ptr.To(""),
			Hidden:      true,
			Sensitive:   false,
		},
	})
}

type mockPropertyData struct {
	propertyName  string
	groupID       string
	contractID    string
	productID     string
	propertyID    string
	latestVersion int
	assetID       string
	cnameFrom     string
	cnameTo       string
}

type mockProperty struct {
	mockPropertyData
	papiMock *papi.Mock
}

func (p *mockProperty) mockCreateProperty() *mock.Call {
	return ExpectCreateProperty(p.papiMock, p.propertyName, p.groupID, p.contractID, p.productID, p.propertyID)
}

func (p *mockProperty) mockUpdatePropertyVersionHostnames() *mock.Call {
	return ExpectUpdatePropertyVersionHostnames(p.papiMock, p.propertyID, p.groupID, p.contractID, p.latestVersion,
		[]papi.Hostname{{
			CnameType:            "EDGE_HOSTNAME",
			CnameFrom:            p.cnameFrom,
			CnameTo:              p.cnameTo,
			CertProvisioningType: "DEFAULT",
		}}, nil)
}

func (p *mockProperty) mockGetProperty() *mock.Call {
	return ExpectGetProperty(p.papiMock, p.propertyID, p.groupID, p.contractID, &papi.Property{
		PropertyName:  p.propertyName,
		GroupID:       p.groupID,
		ContractID:    p.contractID,
		ProductID:     p.productID,
		PropertyID:    p.propertyID,
		LatestVersion: p.latestVersion,
		AssetID:       p.assetID,
	})
}

func (p *mockProperty) mockGetPropertyVersionHostnames() *mock.Call {
	return ExpectGetPropertyVersionHostnames(p.papiMock, p.propertyID, p.groupID, p.contractID, p.latestVersion, &[]papi.Hostname{{
		CnameType:            "EDGE_HOSTNAME",
		CnameFrom:            p.cnameFrom,
		CnameTo:              p.cnameTo,
		CertProvisioningType: "DEFAULT",
	}})
}

func (p *mockProperty) mockGetRuleTree() *mock.Call {
	ruleFormat := ""
	return ExpectGetRuleTree(p.papiMock, p.propertyID, p.groupID, p.contractID, p.latestVersion, nil, &ruleFormat, nil, nil)
}

func (p *mockProperty) mockGetPropertyVersion() *mock.Call {
	return ExpectGetPropertyVersion(p.papiMock, p.propertyID, p.groupID, p.contractID, p.latestVersion, papi.VersionStatusInactive,
		papi.VersionStatusInactive)
}

func (p *mockProperty) mockRemoveProperty() *mock.Call {
	return ExpectRemoveProperty(p.papiMock, p.propertyID, p.contractID, p.groupID)
}

func mockResourcePropertyCreate(p *mockProperty) {
	p.mockCreateProperty().Once()
	p.mockUpdatePropertyVersionHostnames().Once()
	mockResourcePropertyRead(p)
}

func mockResourcePropertyRead(p *mockProperty) {
	p.mockGetProperty().Once()
	p.mockGetPropertyVersionHostnames().Once()
	p.mockGetRuleTree().Once()
	p.mockGetPropertyVersion().Once()
}

func mockMoveProperty(iamMock *iam.Mock, propertyID, srcGroupID, destGroupID int64) {
	iamMock.On("MoveProperty", AnyCTX, iam.MovePropertyRequest{
		PropertyID: propertyID,
		BodyParams: iam.MovePropertyReqBody{
			DestinationGroupID: destGroupID,
			SourceGroupID:      srcGroupID,
		}}).Return(nil)
}

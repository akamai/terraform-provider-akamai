package property

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/stretchr/testify/mock"
)

type mockProperty struct {
	mockPropertyData
	papiMock *papi.Mock
	iamMock  *iam.Mock
}

type mockPropertyData struct {
	propertyName        string
	groupID             string
	contractID          string
	productID           string
	propertyID          string
	assetID             string
	latestVersion       int
	createFromVersion   int
	newVersionID        int
	ruleTree            mockRuleTreeData
	versions            papi.PropertyVersionItems
	hostnames           papi.HostnameResponseItems
	responseErrors      []*papi.Error
	responseWarnings    []*papi.Error
	activations         papi.ActivationsItems
	activationForCreate papi.Activation
	deleteActivationID  string
	groups              papi.GroupItems
	moveGroup           moveGroup
}

func (d *mockPropertyData) getPropertyRequest() papi.GetPropertyRequest {
	return papi.GetPropertyRequest{
		PropertyID: d.propertyID,
		ContractID: d.contractID,
		GroupID:    d.groupID,
	}
}

func (d *mockPropertyData) getPropertyResponse() papi.GetPropertyResponse {
	return papi.GetPropertyResponse{
		Property: &papi.Property{
			AssetID:       d.assetID,
			ContractID:    d.contractID,
			GroupID:       d.groupID,
			LatestVersion: d.latestVersion,
			// although optional in PAPI documentation, ProductID is not being set by PAPI in the response
			PropertyID:   d.propertyID,
			PropertyName: d.propertyName,
		},
	}
}

type moveGroup struct {
	sourceGroupID      int64
	destinationGroupID int64
}

type mockRuleTreeData struct {
	rules        papi.Rules
	comments     string
	ruleFormat   string
	ruleErrors   []papi.RuleError
	ruleWarnings []papi.RuleWarnings
}

func (p *mockProperty) mockCreateProperty(err ...error) *mock.Call {
	req := papi.CreatePropertyRequest{
		GroupID:    p.groupID,
		ContractID: p.contractID,
		Property: papi.PropertyCreate{
			ProductID:    p.productID,
			PropertyName: p.propertyName,
			RuleFormat:   p.ruleTree.ruleFormat,
		},
	}

	if err != nil {
		return p.papiMock.On("CreateProperty", testutils.MockContext, req).Return(nil, err[0]).Once()
	}

	resp := papi.CreatePropertyResponse{PropertyID: p.propertyID}

	return p.papiMock.On("CreateProperty", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetGroups() *mock.Call {
	resp := &papi.GetGroupsResponse{
		Groups: p.groups,
	}

	return p.papiMock.On("GetGroups", testutils.MockContext).Return(resp, nil).Once()
}

func (p *mockProperty) mockUpdateRuleTree(err ...error) *mock.Call {
	rulesUpdate := papi.RulesUpdate{
		Rules:    p.ruleTree.rules,
		Comments: p.ruleTree.comments,
	}

	req := papi.UpdateRulesRequest{
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		ContractID:      p.contractID,
		GroupID:         p.groupID,
		Rules:           rulesUpdate,
		ValidateRules:   true,
	}

	if err != nil {
		return p.papiMock.On("UpdateRuleTree", testutils.MockContext, req).Return(nil, err[0]).Once()
	}

	resp := papi.UpdateRulesResponse{
		PropertyID:      p.propertyID,
		ContractID:      p.contractID,
		GroupID:         p.groupID,
		PropertyVersion: p.latestVersion,
		RuleFormat:      p.ruleTree.ruleFormat,
		Rules:           p.ruleTree.rules,
		Errors:          p.ruleTree.ruleErrors,
		Warnings:        p.ruleTree.ruleWarnings,
	}

	return p.papiMock.On("UpdateRuleTree", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockUpdatePropertyVersionHostnames(err ...error) *mock.Call {
	// Copy hostnames from mock data and remove unnecessary fields (EdgeHostnameID and CertStatus) that are not used in the request to satisfy mocks.
	// Use original mock data for the response.
	requestHostnames := make([]papi.Hostname, len(p.hostnames.Items))
	copy(requestHostnames, p.hostnames.Items)
	for i := range requestHostnames {
		requestHostnames[i].EdgeHostnameID = ""
		requestHostnames[i].CertStatus = papi.CertStatusItem{}
	}

	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		ContractID:      p.contractID,
		GroupID:         p.groupID,
		Hostnames:       requestHostnames,
	}

	if err != nil {
		return p.papiMock.On("UpdatePropertyVersionHostnames", testutils.MockContext, req).Return(&papi.UpdatePropertyVersionHostnamesResponse{}, err[0]).Once()
	}

	resp := papi.UpdatePropertyVersionHostnamesResponse{
		ContractID:      p.contractID,
		GroupID:         p.groupID,
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		Hostnames:       p.hostnames,
	}

	return p.papiMock.On("UpdatePropertyVersionHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetProperty() *mock.Call {
	req := p.getPropertyRequest()
	resp := p.getPropertyResponse()

	if len(p.versions.Items) > 0 && p.versions.Items[0].StagingStatus == papi.VersionStatusActive {
		resp.Property.StagingVersion = &p.versions.Items[0].PropertyVersion
	}

	if len(p.versions.Items) > 0 && p.versions.Items[0].ProductionStatus == papi.VersionStatusActive {
		resp.Property.ProductionVersion = &p.versions.Items[0].PropertyVersion
	}

	return p.papiMock.On("GetProperty", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetPropertyVersionHostnames() *mock.Call {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:        p.propertyID,
		GroupID:           p.groupID,
		ContractID:        p.contractID,
		PropertyVersion:   p.latestVersion,
		IncludeCertStatus: true,
	}

	resp := papi.GetPropertyVersionHostnamesResponse{
		ContractID:      p.contractID,
		GroupID:         p.groupID,
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		Hostnames:       p.hostnames,
	}

	return p.papiMock.On("GetPropertyVersionHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetPropertyVersions() *mock.Call {
	req := papi.GetPropertyVersionsRequest{
		PropertyID: p.propertyID,
		ContractID: p.contractID,
		GroupID:    p.groupID,
	}
	resp := &papi.GetPropertyVersionsResponse{
		PropertyID:   p.propertyID,
		PropertyName: p.propertyName,
		ContractID:   p.contractID,
		GroupID:      p.groupID,
		AssetID:      p.assetID,
		Versions:     p.versions,
	}

	return p.papiMock.On("GetPropertyVersions", testutils.MockContext, req).Return(resp, nil).Once()
}

func (p *mockProperty) mockGetRuleTree() *mock.Call {
	req := papi.GetRuleTreeRequest{
		PropertyID:      p.propertyID,
		GroupID:         p.groupID,
		ContractID:      p.contractID,
		PropertyVersion: p.latestVersion,
		ValidateMode:    "full",
		ValidateRules:   true,
	}

	resp := papi.GetRuleTreeResponse{
		Response: papi.Response{
			Errors:   p.responseErrors,
			Warnings: p.responseWarnings,
		},
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		RuleFormat:      p.ruleTree.ruleFormat,
		Rules:           p.ruleTree.rules,
		Comments:        p.ruleTree.comments,
	}

	return p.papiMock.On("GetRuleTree", testutils.MockContext, req).Return(&resp, nil).Once()
}

// mockGetRuleTreeActivation mocks the GetRuleTree call executed from property_activation resource. It differs with request
// parameters when compared to the GetRuleTree call executed from property resource.
func (p *mockProperty) mockGetRuleTreeActivation() *mock.Call {
	req := papi.GetRuleTreeRequest{
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		ValidateRules:   true,
	}

	resp := papi.GetRuleTreeResponse{
		Response: papi.Response{
			Errors:   p.responseErrors,
			Warnings: p.responseWarnings,
		},
		PropertyID:      p.propertyID,
		PropertyVersion: p.latestVersion,
		RuleFormat:      p.ruleTree.ruleFormat,
		Rules:           p.ruleTree.rules,
		Comments:        p.ruleTree.comments,
	}

	return p.papiMock.On("GetRuleTree", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetPropertyVersion() *mock.Call {
	req := papi.GetPropertyVersionRequest{
		PropertyID:      p.propertyID,
		GroupID:         p.groupID,
		ContractID:      p.contractID,
		PropertyVersion: p.latestVersion,
	}

	var ver papi.PropertyVersionGetItem
	if len(p.versions.Items) > 0 {
		ver = papi.PropertyVersionGetItem{
			StagingStatus:    p.versions.Items[0].StagingStatus,
			ProductionStatus: p.versions.Items[0].ProductionStatus,
			Note:             p.versions.Items[0].Note,
			PropertyVersion:  p.versions.Items[0].PropertyVersion,
		}
	}

	resp := &papi.GetPropertyVersionsResponse{
		PropertyID: p.propertyID,
		GroupID:    p.groupID,
		ContractID: p.contractID,
		Version:    ver,
	}
	return p.papiMock.On("GetPropertyVersion", testutils.MockContext, req).Return(resp, nil).Once()
}

func (p *mockProperty) mockRemoveProperty(err ...error) *mock.Call {
	req := papi.RemovePropertyRequest{
		PropertyID: p.propertyID,
		GroupID:    p.groupID,
		ContractID: p.contractID,
	}
	resp := papi.RemovePropertyResponse{}

	if err != nil {
		return p.papiMock.On("RemoveProperty", testutils.MockContext, req).Return(nil, err[0]).Once()
	}

	return p.papiMock.On("RemoveProperty", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockMoveProperty() {

	// Checking if the property is already in the dst group
	getReq := p.getPropertyRequest()
	getReq.GroupID = "grp_" + fmt.Sprintf("%d", p.moveGroup.destinationGroupID)
	p.papiMock.On("GetProperty", testutils.MockContext, getReq).
		Return(nil, &papi.Error{StatusCode: http.StatusForbidden}).
		Once()

	activationsReq := papi.GetActivationsRequest{
		PropertyID: p.propertyID,
		ContractID: p.contractID,
		GroupID:    p.groupID,
	}
	var act1 = &papi.Activation{
		ActivationID: "dummy_activation_id",
	}
	var activationsRes = &papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{
			Items: []*papi.Activation{
				act1,
			},
		},
	}
	p.papiMock.On("GetActivations", testutils.MockContext, activationsReq).Return(activationsRes, nil)

	prpID := strings.TrimPrefix(p.assetID, "aid_")
	intPropertyID, err := strconv.ParseInt(prpID, 10, 64)
	// shouldn't happen, unless wrong format of propertyID is provided
	if err != nil {
		panic(err)
	}
	req := iam.MovePropertyRequest{
		PropertyID: intPropertyID,
		Body: iam.MovePropertyRequestBody{
			DestinationGroupID: p.moveGroup.destinationGroupID,
			SourceGroupID:      p.moveGroup.sourceGroupID,
		},
	}

	p.iamMock.On("MoveProperty", testutils.MockContext, req).Return(nil).Once()
}

func (p *mockProperty) mockGetActivations() *mock.Call {
	req := papi.GetActivationsRequest{
		PropertyID: p.propertyID,
	}
	resp := papi.GetActivationsResponse{
		Activations: p.activations,
	}
	return p.papiMock.On("GetActivations", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockGetActivationsCompleteRequest(err ...error) *mock.Call {
	req := papi.GetActivationsRequest{
		PropertyID: p.propertyID,
		ContractID: p.contractID,
		GroupID:    p.groupID,
	}
	if err != nil {
		return p.papiMock.On("GetActivations", testutils.MockContext, req).Return(nil, err[0]).Once()
	}
	resp := papi.GetActivationsResponse{
		Activations: p.activations,
	}
	return p.papiMock.On("GetActivations", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockCreateActivation() *mock.Call {
	activation := papi.Activation{
		ActivationType:  p.activationForCreate.ActivationType,
		Network:         p.activationForCreate.Network,
		NotifyEmails:    p.activationForCreate.NotifyEmails,
		PropertyVersion: p.activationForCreate.PropertyVersion,
	}
	req := papi.CreateActivationRequest{
		PropertyID: p.propertyID,
		Activation: activation,
	}
	resp := papi.CreateActivationResponse{
		ActivationID: p.activationForCreate.ActivationID,
	}

	activation.ActivationID = p.activationForCreate.ActivationID
	activation.GroupID = p.groupID
	activation.PropertyName = p.propertyName
	activation.PropertyID = p.propertyID
	activation.Status = p.activationForCreate.Status
	activation.SubmitDate = p.activationForCreate.SubmitDate

	// modify mock data to reflect newly created activation
	p.activations.Items = append(p.activations.Items, &activation)

	return p.papiMock.On("CreateActivation", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockResourceActivationDelete() {
	p.mockGetActivations()

	activation := papi.Activation{
		ActivationType:  papi.ActivationTypeDeactivate,
		Network:         p.activationForCreate.Network,
		NotifyEmails:    p.activationForCreate.NotifyEmails,
		PropertyVersion: p.activationForCreate.PropertyVersion,
	}
	req := papi.CreateActivationRequest{
		PropertyID: p.propertyID,
		Activation: activation,
	}
	resp := papi.CreateActivationResponse{
		ActivationID: p.deleteActivationID,
	}

	activation.ActivationID = p.deleteActivationID
	activation.GroupID = p.groupID
	activation.PropertyName = p.propertyName
	activation.PropertyID = p.propertyID
	activation.Status = papi.ActivationStatusActive

	// modify mock data to reflect newly created activation
	p.activations.Items = append(p.activations.Items, &activation)

	p.papiMock.On("CreateActivation", testutils.MockContext, req).Return(&resp, nil).Once()
	p.mockGetActivation()
}

func (p *mockProperty) mockGetActivation() *mock.Call {
	activation := p.activations.Items[len(p.activations.Items)-1]
	req := papi.GetActivationRequest{
		PropertyID:   p.propertyID,
		ActivationID: activation.ActivationID,
	}
	resp := papi.GetActivationResponse{
		Activation: activation,
	}

	return p.papiMock.On("GetActivation", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockCreatePropertyVersion() *mock.Call {
	req := papi.CreatePropertyVersionRequest{
		PropertyID: p.propertyID,
		GroupID:    p.groupID,
		ContractID: p.contractID,
		Version: papi.PropertyVersionCreate{
			CreateFromVersion: p.createFromVersion,
		},
	}

	resp := papi.CreatePropertyVersionResponse{PropertyVersion: p.newVersionID}

	return p.papiMock.On("CreatePropertyVersion", testutils.MockContext, req).Return(&resp, nil).Once()
}

// mockResourcePropertyCreateWithVersionHostnames represents the creation workflow where the property and hostnames are created
func mockResourcePropertyCreateWithVersionHostnames(p *mockProperty) {
	p.mockCreateProperty()
	p.mockUpdatePropertyVersionHostnames()
}

// mockResourcePropertyFullCreate represents the full creation workflow where the property, hostnames and rule tree are created
func mockResourcePropertyFullCreate(p *mockProperty) {
	p.mockCreateProperty()
	p.mockUpdatePropertyVersionHostnames()
	p.mockUpdateRuleTree()
}

// mockResourcePropertyRead represents the read workflow where GetProperty call is used (version of the property is known)
func mockResourcePropertyRead(p *mockProperty, times ...int) {
	i := 1
	if len(times) > 0 {
		i = times[0]
	}
	p.mockGetProperty().Times(i)
	p.mockGetPropertyVersionHostnames().Times(i)
	p.mockGetRuleTree().Times(i)
	p.mockGetPropertyVersion().Times(i)
}

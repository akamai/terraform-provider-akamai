package property

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
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
	useHostnameBucket   bool
	responseErrors      []*papi.Error
	responseWarnings    []*papi.Error
	activations         papi.ActivationsItems
	activationForCreate papi.Activation
	deleteActivationID  string
	groups              papi.GroupItems
	moveGroup           moveGroup
	hostnameBucket      hostnameBucket
}

func (d *mockPropertyData) getPropertyRequest() papi.GetPropertyRequest {
	return papi.GetPropertyRequest{
		PropertyID: d.propertyID,
		ContractID: d.contractID,
		GroupID:    d.groupID,
	}
}

func (d *mockPropertyData) getPropertyResponse() papi.GetPropertyResponse {
	var propertyType *string
	if d.useHostnameBucket {
		propertyType = ptr.To("HOSTNAME_BUCKET")
	}
	return papi.GetPropertyResponse{
		Property: &papi.Property{
			AssetID:       d.assetID,
			ContractID:    d.contractID,
			GroupID:       d.groupID,
			LatestVersion: d.latestVersion,
			// although optional in PAPI documentation, ProductID is not being set by PAPI in the response
			PropertyID:   d.propertyID,
			PropertyName: d.propertyName,
			PropertyType: propertyType,
		},
	}
}

type hostnameBucket struct {
	plan         map[string]Hostname
	state        map[string]Hostname
	network      papi.ActivationNetwork
	notifyEmails []string
	note         string
	activations  papi.ListPropertyHostnameActivationsResponse
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

func (p *mockProperty) mockPatchPropertyHostnameBucket() {

	var inUpdate bool
	if len(p.hostnameBucket.state) != 0 || len(p.hostnameBucket.activations.HostnameActivations.Items) > 0 {
		inUpdate = true
	}

	rb := hostnameRequestBuilder{
		hostnameRequestData: hostnameRequestData{
			planHostnames:  p.hostnameBucket.plan,
			stateHostnames: p.hostnameBucket.state,
			propertyID:     p.propertyID,
			contractID:     p.contractID,
			groupID:        p.groupID,
			network:        string(p.hostnameBucket.network),
			emails:         p.hostnameBucket.notifyEmails,
			note:           p.hostnameBucket.note,
		},
		ctx: context.Background(),
	}
	requestData, diags := rb.build()
	if diags.HasError() {
		panic(diags)
	}

	for i, req := range requestData.requests {
		var actID string
		if inUpdate {
			actID = "act_%d_update"
		} else {
			actID = "act_%d"
		}

		p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
			ActivationID: fmt.Sprintf(actID, i),
		}, nil).Once()

		getReq := papi.GetPropertyHostnameActivationRequest{
			PropertyID:           p.propertyID,
			HostnameActivationID: fmt.Sprintf(actID, i),
			ContractID:           p.contractID,
			GroupID:              p.groupID,
		}

		p.hostnameBucket.activations.HostnameActivations.Items = append(p.hostnameBucket.activations.HostnameActivations.Items, papi.HostnameActivationListItem{
			HostnameActivationID: fmt.Sprintf(actID, i),
			PropertyID:           p.propertyID,
			Network:              papi.ActivationNetwork(p.hostnameBucket.network),
			Status:               "ACTIVE",
			Note:                 p.hostnameBucket.note,
			NotifyEmails:         p.hostnameBucket.notifyEmails,
		})
		p.hostnameBucket.activations.HostnameActivations.TotalItems = len(p.hostnameBucket.state)

		p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
			ContractID: p.propertyID,
			GroupID:    p.groupID,
			HostnameActivation: papi.HostnameActivationGetItem{
				HostnameActivationID: fmt.Sprintf(actID, i),
				PropertyID:           p.propertyID,
				Network:              p.hostnameBucket.network,
				Status:               "ACTIVE",
				Note:                 p.hostnameBucket.note,
				NotifyEmails:         p.hostnameBucket.notifyEmails,
			},
		}, nil).Once()
	}
}

func (p *mockProperty) mockListActivePropertyHostnames(withoutGroupAndContract ...bool) {
	// if withoutGroupAndContract was provided with a 'true' value, then fill the requests with empty contract and group
	var reqContractID, reqGroupID string
	if (len(withoutGroupAndContract) > 0 && !withoutGroupAndContract[0]) || len(withoutGroupAndContract) == 0 {
		reqContractID = p.contractID
		reqGroupID = p.groupID
	}

	var hostnameItems []papi.HostnameItem
	for k, v := range p.hostnameBucket.state {
		if p.hostnameBucket.network == "STAGING" {
			hostnameItem := papi.HostnameItem{
				CnameFrom:             k,
				CnameType:             papi.HostnameCnameTypeEdgeHostname,
				StagingCertType:       papi.CertType(v.CertProvisioningType.ValueString()),
				StagingCnameTo:        v.CnameTo.ValueString(),
				StagingEdgeHostnameID: v.EdgeHostnameID.ValueString(),
			}
			if v.CertProvisioningType.ValueString() == string(papi.CertTypeDefault) {
				hostnameItem.CertStatus = &papi.CertStatusItem{Staging: []papi.StatusItem{{Status: "PENDING"}}}
			}
			hostnameItems = append(hostnameItems, hostnameItem)
		} else {
			hostnameItem := papi.HostnameItem{
				CnameFrom:                k,
				CnameType:                papi.HostnameCnameTypeEdgeHostname,
				ProductionCertType:       papi.CertType(v.CertProvisioningType.ValueString()),
				ProductionCnameTo:        v.CnameTo.ValueString(),
				ProductionEdgeHostnameID: v.EdgeHostnameID.ValueString(),
			}
			if v.CertProvisioningType.ValueString() == string(papi.CertTypeDefault) {
				hostnameItem.CertStatus = &papi.CertStatusItem{Production: []papi.StatusItem{{Status: "PENDING"}}}
			}
			hostnameItems = append(hostnameItems, hostnameItem)
		}
	}

	offset := 0
	limit := 999

	for len(hostnameItems) > 999 {
		req := papi.ListActivePropertyHostnamesRequest{
			PropertyID:        p.propertyID,
			Offset:            offset,
			Limit:             limit,
			Network:           p.hostnameBucket.network,
			ContractID:        reqContractID,
			GroupID:           reqGroupID,
			IncludeCertStatus: true,
			Sort:              "hostname:a",
		}
		offset += 999
		resp := papi.ListActivePropertyHostnamesResponse{
			ContractID: p.contractID,
			GroupID:    p.groupID,
			PropertyID: p.propertyID,
			Hostnames: papi.HostnamesResponseItems{
				Items:            hostnameItems[:999],
				CurrentItemCount: 999,
				TotalItems:       len(p.hostnameBucket.state),
			},
		}
		hostnameItems = hostnameItems[999:]
		p.papiMock.On("ListActivePropertyHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
	}

	req := papi.ListActivePropertyHostnamesRequest{
		PropertyID:        p.propertyID,
		Offset:            offset,
		Limit:             limit,
		Network:           p.hostnameBucket.network,
		ContractID:        reqContractID,
		GroupID:           reqGroupID,
		IncludeCertStatus: true,
		Sort:              "hostname:a",
	}
	resp := papi.ListActivePropertyHostnamesResponse{
		ContractID: p.contractID,
		GroupID:    p.groupID,
		PropertyID: p.propertyID,
		Hostnames: papi.HostnamesResponseItems{
			Items:            hostnameItems,
			CurrentItemCount: len(hostnameItems),
			TotalItems:       len(p.hostnameBucket.state),
		},
	}
	p.papiMock.On("ListActivePropertyHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockListPropertyHostnameActivations() {
	req := papi.ListPropertyHostnameActivationsRequest{
		PropertyID: p.propertyID,
		Offset:     0,
		Limit:      999,
		ContractID: p.contractID,
		GroupID:    p.groupID,
	}
	resp := p.hostnameBucket.activations

	p.papiMock.On("ListPropertyHostnameActivations", testutils.MockContext, req).Return(&resp, nil).Once()
}

func (p *mockProperty) mockCreateProperty(err ...error) *mock.Call {
	req := papi.CreatePropertyRequest{
		GroupID:    p.groupID,
		ContractID: p.contractID,
		Property: papi.PropertyCreate{
			ProductID:         p.productID,
			PropertyName:      p.propertyName,
			RuleFormat:        p.ruleTree.ruleFormat,
			UseHostnameBucket: p.useHostnameBucket,
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
		// Links are used only for mocking responses
		if requestHostnames[i].CCMCertificates != nil {
			// copy on write
			copyCerts := *requestHostnames[i].CCMCertificates
			copyCerts.RSACertLink = ""
			copyCerts.ECDSACertLink = ""
			requestHostnames[i].CCMCertificates = &copyCerts
		}
		// CCMCertStatus is used only for mocking responses
		requestHostnames[i].CCMCertStatus = nil

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
			ProductID:        p.productID,
		}
	}

	resp := &papi.GetPropertyVersionsResponse{
		PropertyID:   p.propertyID,
		PropertyName: p.propertyName,
		GroupID:      p.groupID,
		ContractID:   p.contractID,
		Version:      ver,
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

func mockResourceHostnameBucketDelete(p *mockProperty) {
	p.hostnameBucket.plan = map[string]Hostname{}
	p.mockPatchPropertyHostnameBucket()
}

func mockResourceHostnameBucketRead(p *mockProperty, times ...int) {
	if len(times) == 1 && times[0] == 2 {
		p.mockListActivePropertyHostnames()
		p.mockListActivePropertyHostnames()
		p.mockListPropertyHostnameActivations()
		p.mockListPropertyHostnameActivations()
		return
	}
	p.mockListActivePropertyHostnames()
	p.mockListPropertyHostnameActivations()
}

// mockResourceHostnameBucketUpsert mocks either Create or Update, as the operations are exactly the same.
func mockResourceHostnameBucketUpsert(p *mockProperty) {
	p.mockPatchPropertyHostnameBucket()
	p.hostnameBucket.state = maps.Clone(p.hostnameBucket.plan)
	p.mockListActivePropertyHostnames()
}

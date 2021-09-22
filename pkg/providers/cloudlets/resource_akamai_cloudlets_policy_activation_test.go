package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCloudletsPolicyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*mockcloudlets)
		steps []resource.TestStep
	}{
		"create and read activation, version == 1, inactive -> activate": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active policy activation -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *mockcloudlets) {
				// create, policy active so no need to activate
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusActive).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> read -> ": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation read: an error"),
				},
			},
		},
		"create and read activation, no version": {
			init: func(m *mockcloudlets) {
				// create
				expectListPolicyVersions(m, 1234, nil).Once()
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll after activation
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyVersions(m, 1234, nil).Once()
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Once()
			}, steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_no_version.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
					),
				},
			},
		},
		"create and read activation, no version -> ": {
			init: func(m *mockcloudlets) {
				// create
				expectListPolicyVersions(m, 1234, fmt.Errorf("an error")).Once()
			}, steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_no_version.tf"),
					ExpectError: regexp.MustCompile("policy activation: an error"),
				},
			},
		},
		"create and read activation, update - no changes": {
			init: func(m *mockcloudlets) {
				// first test step
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll after activation
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: version already active, read": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", "", cloudlets.StatusActive, nil).Once()
				// update
				expectGetPolicy(m, int64(1234), "property name", nil).Once()
				expectCreatePolicyVersion(m, int64(1234), int64(2), nil).Once()
				expectListPolicyActivations(m, 1234, 2, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", "", cloudlets.StatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: version not active, activate": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", "", cloudlets.StatusActive, nil).Once()
				// update
				expectGetPolicy(m, int64(1234), "property name", nil).Once()
				expectCreatePolicyVersion(m, int64(1234), int64(2), nil).Once()
				expectListPolicyActivations(m, 1234, 2, "staging", "property name", cloudlets.StatusDeactivated, nil).Once()
				expectActivatePolicyVersion(m, 1234, 2, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil)
				expectListPolicyActivations(m, 1234, 2, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", "", cloudlets.StatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: ListPolicyActivations error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// update
				expectGetPolicy(m, int64(1234), "property name", nil).Once()
				expectCreatePolicyVersion(m, int64(1234), int64(2), nil).Once()
				expectListPolicyActivations(m, 1234, 2, "staging", "property name", cloudlets.StatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: create policy version error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// update
				expectGetPolicy(m, int64(1234), "property name", nil).Once()
				expectCreatePolicyVersion(m, int64(1234), int64(2), fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: get policy error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// update
				expectGetPolicy(m, int64(1234), "property name", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: CreatePolicyVersion error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// update
				expectGetPolicy(m, int64(1234), "property name", nil).Once()
				expectCreatePolicyVersion(m, int64(1234), int64(2), fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: GetPolicy error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Times(2)
				// 2 - for policy_activation_update.tf
				// update
				expectGetPolicy(m, int64(1234), "property name", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: version already active, read - cannot find version": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, cloudlets.StatusInactive).Once()
				expectActivatePolicyVersion(m, 1234, 1, cloudlets.VersionActivationNetworkStaging, []string{"prp_0", "prp_1"}, nil).Once()
				// create: poll until active
				expectListPolicyActivations(m, 1234, 1, "staging", "property name", cloudlets.StatusActive, nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", "", cloudlets.StatusActive, nil).Once()
				// 2 - for policy_activation_update.tf
				expectListPolicyActivations(m, 1234, 2, "staging", "", cloudlets.StatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsPropertyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("cannot find the given policy activation version"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockcloudlets{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectGetPolicyVersion(m *mockcloudlets, policyID, version int64, status cloudlets.Status) *mock.Call {
	return m.On(
		"GetPolicyVersion",
		mock.Anything,
		cloudlets.GetPolicyVersionRequest{PolicyID: policyID, Version: 1, OmitRules: true},
	).Return(
		&cloudlets.PolicyVersion{
			Activations: []*cloudlets.Activation{
				{
					Network:      string(cloudlets.VersionActivationNetworkProduction),
					APIVersion:   "1.0",
					PropertyInfo: cloudlets.PropertyInfo{Name: "property name"},
					PolicyInfo:   cloudlets.PolicyInfo{PolicyID: policyID, Status: status, Version: version},
				},
				{
					Network:      string(cloudlets.VersionActivationNetworkStaging),
					APIVersion:   "1.0",
					PropertyInfo: cloudlets.PropertyInfo{Name: "property name"},
					PolicyInfo:   cloudlets.PolicyInfo{PolicyID: policyID, Status: status, Version: version},
				},
			},
		}, nil)
}

func expectListPolicyActivations(m *mockcloudlets, policyID, version int64, network, propertyName string, status cloudlets.Status, err error) *mock.Call {
	if err != nil {
		return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
			PolicyID:     policyID,
			Network:      cloudlets.VersionActivationNetwork(tools.StateNetwork(network)),
			PropertyName: propertyName,
		}).Return(nil, err)
	}
	return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
		PolicyID:     policyID,
		Network:      cloudlets.VersionActivationNetwork(tools.StateNetwork(network)),
		PropertyName: propertyName,
	}).Return([]cloudlets.PolicyActivation{
		{
			Network:      cloudlets.VersionActivationNetwork(tools.StateNetwork(network)),
			APIVersion:   "1.0",
			PropertyInfo: cloudlets.PropertyInfo{Name: propertyName},
			PolicyInfo:   cloudlets.PolicyInfo{PolicyID: policyID, Status: status, Version: version},
		},
	}, nil)
}

func expectActivatePolicyVersion(m *mockcloudlets, policyID, version int64, network cloudlets.VersionActivationNetwork, additionalProps []string, err error) *mock.Call {
	return m.On("ActivatePolicyVersion", mock.Anything, cloudlets.ActivatePolicyVersionRequest{
		PolicyID:    policyID,
		Async:       true,
		Version:     version,
		RequestBody: cloudlets.ActivatePolicyVersionRequestBody{Network: network, AdditionalPropertyNames: additionalProps},
	}).Return(err)
}

func expectListPolicyVersions(m *mockcloudlets, policyID int64, err error) *mock.Call {
	if err != nil {
		return m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
			PolicyID: policyID,
		}).Return(nil, err)
	}
	return m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
		PolicyID: policyID,
	}).Return([]cloudlets.PolicyVersion{{Version: 1}}, nil)
}

func expectGetPolicy(m *mockcloudlets, policyID int64, propertyName string, err error) *mock.Call {
	if err != nil {
		return m.On("GetPolicy", mock.Anything, policyID).Return(
			nil, err,
		)
	}
	return m.On("GetPolicy", mock.Anything, policyID).Return(
		&cloudlets.Policy{Activations: []cloudlets.Activation{
			{
				PropertyInfo: cloudlets.PropertyInfo{Name: propertyName},
			},
		}}, nil,
	)
}

func expectCreatePolicyVersion(m *mockcloudlets, policyID, version int64, err error) *mock.Call {
	if err != nil {
		return m.On("CreatePolicyVersion", mock.Anything, cloudlets.CreatePolicyVersionRequest{
			PolicyID:            policyID,
			CreatePolicyVersion: cloudlets.CreatePolicyVersion{},
		}).Return(
			nil, err,
		)
	}
	return m.On("CreatePolicyVersion", mock.Anything, cloudlets.CreatePolicyVersionRequest{
		PolicyID:            policyID,
		CreatePolicyVersion: cloudlets.CreatePolicyVersion{},
	}).Return(
		&cloudlets.PolicyVersion{Version: version}, nil,
	)
}

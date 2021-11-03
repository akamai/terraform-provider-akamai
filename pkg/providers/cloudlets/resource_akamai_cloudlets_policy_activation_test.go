package cloudlets

import (
	"fmt"
	"regexp"
	"testing"
	"time"

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
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"create and read activation, version == 1, production, inactive -> activate": {
			init: func(m *mockcloudlets) {
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "prod", 1)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_production.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "prod"),
					),
				},
			},
		},
		"create and read activation, version == 1, prod, inactive -> activate": {
			init: func(m *mockcloudlets) {
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "prod", 1)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "prod", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusInactive,
					}},
				}).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active policy activation -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "staging", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusInactive,
					}},
				}).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *mockcloudlets) {
				// create, policy active so no need to activate
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "staging", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusActive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: "prp_0"}},
					{APIVersion: "1.0", Network: "staging", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusActive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: "prp_1"}},
				}).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> read -> ListPolicyActivations read error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation read: an error"),
				},
			},
		},
		"create and read activation, update - no changes": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// 2 - for policy_activation_version1.tf
				// update
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "activation failed", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"create activation - failed activation while polling": {
			init: func(m *mockcloudlets) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusFailed, "activation failed", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation: policyID 1234: activation failed"),
				},
			},
		},
		"Create and read activation. Update: version already active, read": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				// read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// update
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: version not active, activate": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				// refresh read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// update
				expectListPolicyActivations(m, 1234, 2, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				expectActivatePolicyVersion(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, nil)
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// read
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: ListPolicyActivations error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				// refresh read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// update
				expectListPolicyActivations(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: activate policy version error": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				// refresh read
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// update
				expectListPolicyActivations(m, 1234, 2, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				expectActivatePolicyVersion(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: ListPolicyActivations error while polling": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				// update
				expectListPolicyActivations(m, 1234, 1, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
				expectActivatePolicyVersion(m, 1234, 2, "staging", []string{"prp_0", "prp_1"}, nil)
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 2, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: version already active, read - cannot find version": {
			init: func(m *mockcloudlets) {
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "staging", 1)
				// 2 - for policy_activation_update.tf
				// update
				expectListPolicyActivations(m, 1234, 2, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update.tf"),
					ExpectError: regexp.MustCompile("cannot find the given policy activation version"),
				},
			},
		},
	}

	// redefining times to accelerate tests
	ActivationPollMinimum = time.Second * 1
	ActivationPollInterval = time.Second * 1

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

// expect full activation of a policy (creation)
func expectFullActivation(m *mockcloudlets, policyID, version int64, properties []string, network cloudlets.PolicyActivationNetwork, times int) {
	// create
	expectGetPolicyVersion(m, policyID, 1, []cloudlets.PolicyActivation{
		{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
			PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
		}},
	}).Times(times)
	expectActivatePolicyVersion(m, policyID, version, network, properties, nil).Times(times)
	// poll until active -> waitForPolicyActivation()
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", nil).Times(times)
	// read
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", nil).Times(times)
	// read
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", nil).Times(times)
}

func expectGetPolicyVersion(m *mockcloudlets, policyID, version int64, activations []cloudlets.PolicyActivation) *mock.Call {
	return m.On(
		"GetPolicyVersion",
		mock.Anything,
		cloudlets.GetPolicyVersionRequest{PolicyID: policyID, Version: version, OmitRules: true},
	).Return(
		&cloudlets.PolicyVersion{
			Activations: activations,
		}, nil)
}

func expectListPolicyActivations(m *mockcloudlets, policyID, version int64, network cloudlets.PolicyActivationNetwork, propertyNames []string, status cloudlets.PolicyActivationStatus, statusDetail string, expectedErr error) *mock.Call {
	if expectedErr != nil {
		return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
			PolicyID: policyID,
			Network:  network,
		}).Return(nil, expectedErr)
	}

	policyActivations := []cloudlets.PolicyActivation{}

	for _, propertyName := range propertyNames {
		policyActivations = append(policyActivations, cloudlets.PolicyActivation{
			APIVersion:   "1.0",
			Network:      network,
			PolicyInfo:   cloudlets.PolicyInfo{PolicyID: policyID, Status: status, Version: version, StatusDetail: statusDetail},
			PropertyInfo: cloudlets.PropertyInfo{Name: propertyName},
		})
	}

	return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	}).Return(policyActivations, nil)
}

func expectActivatePolicyVersion(m *mockcloudlets, policyID, version int64, network cloudlets.PolicyActivationNetwork, additionalProps []string, err error) *mock.Call {
	return m.On("ActivatePolicyVersion", mock.Anything, cloudlets.ActivatePolicyVersionRequest{
		PolicyID:                policyID,
		Async:                   true,
		Version:                 version,
		PolicyVersionActivation: cloudlets.PolicyVersionActivation{Network: network, AdditionalPropertyNames: additionalProps},
	}).Return(err)
}

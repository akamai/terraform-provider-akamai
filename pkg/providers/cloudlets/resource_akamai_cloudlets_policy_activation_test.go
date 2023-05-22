package cloudlets

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cloudlets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCloudletsPolicyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*cloudlets.Mock)
		steps []resource.TestStep
	}{
		"try to create activation with a non existing policy version": {
			init: func(m *cloudlets.Mock) {
				m.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{PolicyID: 1234, Version: 1, OmitRules: true}).Return(nil, fmt.Errorf("an error"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile(`policy activation: cannot find the given policy version \(1\): an error`),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate": {
			init: func(m *cloudlets.Mock) {
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationNetworkStaging, 1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
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
		"create and read activation, version == 1, failed with property error -> retry -> activate": {
			init: func(m *cloudlets.Mock) {
				policyID, version, network, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				propertyNotFoundError := fmt.Errorf(`"detail": "Requested propertyName \"test.property.name\" does not exist"`)
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusPending, "", 1, propertyNotFoundError).Once()
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
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
		"create and read activation, version == 1, pending activations on activate -> empty list activations -> retry -> correct list activations": {
			init: func(m *cloudlets.Mock) {
				policyID, version, network, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				// empty activations
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// retry
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
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
		"create and read activation, version == 1, pending activations on activate -> empty list activations -> retry max times -> failed": {
			init: func(m *cloudlets.Mock) {
				policyID, version, network, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				// empty activations retry # 1
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 2
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 3
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: policy activation: policyID 1234"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> wait -> failed": {
			init: func(m *cloudlets.Mock) {
				policyID, version, staging, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusFailed, "", 1, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: policy activation: policyID 1234"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> several iterations in read function": {
			init: func(m *cloudlets.Mock) {
				policyID, version, staging, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusPending, "", -1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusActive, "", 2, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusActive, "", 2, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, version, staging, properties, cloudlets.PolicyActivationStatusActive, "", 2, nil).Once()

				// delete
				expectDeletePhase(m, policyID, properties, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
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
		"create and read activation, version == 1, inactive -> activate -> read -> no activations -> error": {
			init: func(m *cloudlets.Mock) {
				properties, policyID, v1, staging, times := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, 1
				// create
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: v1, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, v1, activations, nil).Times(times)
				expectActivatePolicyVersion(m, policyID, v1, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Times(times)
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v1, staging, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(times)
				// read
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile(`policy activation read: cannot find any activation for the given policy \(1234\) and network \('staging'\)`),
				},
			},
		},
		"create and read activation, version == 1, production, inactive -> activate": {
			init: func(m *cloudlets.Mock) {
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "prod", 1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkProduction, nil, nil)
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
			init: func(m *cloudlets.Mock) {
				expectFullActivation(m, 1234, 1, []string{"prp_0", "prp_1"}, "prod", 1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkProduction, nil, nil)
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
			init: func(m *cloudlets.Mock) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "prod", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusInactive,
					}},
				}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active policy activation -> error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "staging", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusInactive,
					}},
				}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{}, cloudlets.PolicyActivationStatusActive, "", 1, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// create, policy active so no need to activate
				expectGetPolicyVersion(m, policyID, v1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: v1, Status: active,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: "prp_0"}},
					{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: v1, Status: active,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: "prp_1"}},
				}, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Twice()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// create
				expectGetPolicyVersion(m, policyID, v1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, policyID, v1, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation read: an error"),
				},
			},
		},
		"create and read activation, update - no changes, so skip update": {
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				// create
				expectGetPolicyVersion(m, policyID, v1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, policyID, v1, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// 2 - for policy_activation_version1.tf
				// update
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "activation failed", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
			init: func(m *cloudlets.Mock) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusFailed, "activation failed", 1, fmt.Errorf("activation failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: activation failed"),
				},
			},
		},
		"create activation - failed activation while polling with no failed status": {
			init: func(m *cloudlets.Mock) {
				// create
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusDeactivated, "failed", 1, fmt.Errorf("activation failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: activation failed"),
				},
			},
		},
		"create and read activation, update - cannot find the property version": {
			init: func(m *cloudlets.Mock) {
				policyID, v1, v2, properties, active, staging := int64(1234), int64(1), int64(2), []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, cloudlets.PolicyActivationNetworkStaging
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_version1.tf
				// read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, fmt.Errorf("an error"))
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile(`policy activation: cannot find the given policy version \(2\): an error`),
				},
			},
		},
		"Create and read activation. Update: version already active, read": {
			init: func(m *cloudlets.Mock) {
				policyID, version, properties, active, staging := int64(1234), int64(1), []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusActive, cloudlets.PolicyActivationNetworkStaging
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, version, properties, staging, 1)
				// 2 - for policy_activation_version1.tf
				// read
				expectListPolicyActivations(m, policyID, version, staging, properties, active, "", 1, nil).Once()
				// update
				expectListPolicyActivations(m, policyID, version, staging, properties, active, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
		"Create and read activation. Update: version not active, activate": {
			init: func(m *cloudlets.Mock) {
				properties, staging, policyID, v1, v2, active := []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationNetworkStaging, int64(1234), int64(1), int64(2), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, 1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, []string{}, active, "", 1, nil).Once()
				expectGetPolicyProperties(m, policyID, properties, nil).Once()
				expectListPolicyActivations(m, policyID, v1, "", properties, active, "", 1, nil)
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil)
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"Create and read activation. Update: change version from staging to prod, activate": {
			init: func(m *cloudlets.Mock) {
				properties, policyID, v1, active := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, 1, properties, cloudlets.PolicyActivationNetworkStaging, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkStaging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v1, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkStaging, properties, active, "", 1, nil)
				expectGetPolicyProperties(m, policyID, properties, nil).Once()
				expectActivatePolicyVersion(m, policyID, v1, cloudlets.PolicyActivationNetworkProduction, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil)
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkProduction, []string{}, active, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkProduction, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkProduction, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v1, cloudlets.PolicyActivationNetworkProduction, properties, active, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, cloudlets.PolicyActivationNetworkProduction, nil, nil)
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
					Config: loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "prod"),
					),
				},
			},
		},
		"Create and read activation. Update: ListPolicyActivations error": {
			init: func(m *cloudlets.Mock) {
				policyID, v1, v2, net, properties := int64(1234), int64(1), int64(2), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, net, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, net, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v2, net, properties, cloudlets.PolicyActivationStatusActive, "", 1, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, net, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: activate policy version error": {
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, v2, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), int64(2), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, []string{}, active, "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, "", properties, active, "", 1, nil)
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: delete policy property error": {
			init: func(m *cloudlets.Mock) {
				staging, properties, moreProperties, policyID, v1, v2, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, []string{"prp_0", "prp_1", "prp_2"}, int64(1234), int64(1), int64(2), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, moreProperties, active, "", 1, nil).Once()
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				expectGetPolicyProperties(m, policyID, moreProperties, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, moreProperties, active, "", 1, nil)
				expectDeletePolicyProperty(m, policyID, 20, staging, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, moreProperties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation: cannot remove policy 1234 property 20 and network 'staging'. Please, try once again later.\nan error"),
				},
			},
		},
		"Create and read activation. Update: get policy properties error": {
			init: func(m *cloudlets.Mock) {
				properties, staging, policyID, v1, v2, active := []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationNetworkStaging, int64(1234), int64(1), int64(2), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, []string{}, active, "", 1, nil).Once()
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				expectGetPolicyProperties(m, policyID, nil, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation: cannot find policy 1234 properties: an error"),
				},
			},
		},
		"Create and read activation. Update: ListPolicyActivations error while polling": {
			init: func(m *cloudlets.Mock) {
				policyID, v1, v2, properties, staging, active := int64(1234), int64(1), int64(2), []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationNetworkStaging, cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				// update
				expectGetPolicyVersion(m, policyID, v2, []cloudlets.PolicyActivation{}, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, properties, active, "", 1, nil).Once()
				expectGetPolicyProperties(m, policyID, properties, nil).Once()
				expectListPolicyActivations(m, policyID, v1, "", properties, active, "", 1, nil)
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil)
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v2, staging, []string{}, active, "", 1, fmt.Errorf("an error")).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},
		"Create and read activation. Update: version already active, read - cannot find version": {
			init: func(m *cloudlets.Mock) {
				policyID, v1, v2, properties, staging, active := int64(1234), int64(1), int64(2), []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationNetworkStaging, cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				expectFullActivation(m, policyID, v1, properties, staging, 1)
				// 2 - for policy_activation_update_version2.tf
				// update
				expectListPolicyActivations(m, policyID, v2, staging, []string{}, active, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
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
					Config:      loadFixtureString("./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile(`policy activation read: cannot find any activation for the given policy \(1234\) and network \('staging'\)`),
				},
			},
		},
	}

	// redefining times to accelerate tests
	ActivationPollMinimum, ActivationPollInterval, PolicyActivationRetryPollMinimum = time.Millisecond, time.Millisecond, time.Millisecond

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudlets.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps:             test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// expect full activation of a policy (creation)
func expectFullActivation(m *cloudlets.Mock, policyID, version int64, properties []string, network cloudlets.PolicyActivationNetwork, times int) {
	// create
	activations := make([]cloudlets.PolicyActivation, len(properties))
	for _, p := range properties {
		activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
			PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
		}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
	}
	expectGetPolicyVersion(m, policyID, version, activations, nil).Times(times)
	expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Times(times)
	// poll until active -> waitForPolicyActivation()
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(times)
	// read
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(times)
	// read
	expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(times)
}

// expect delete step
func expectDeletePhase(m *cloudlets.Mock, policyID int64, deletedProperties, remainingProperties []string, network cloudlets.PolicyActivationNetwork, errGetPolicy, errDeleteProperty error) {
	expectGetPolicyProperties(m, policyID, deletedProperties, errGetPolicy).Once()
	if errGetPolicy != nil {
		return
	}
	expectListPolicyActivations(m, policyID, int64(1), network, deletedProperties, cloudlets.PolicyActivationStatusActive, "", 1, nil)
	for idx := range deletedProperties {
		expectListPolicyActivations(m, policyID, int64(1), "", remainingProperties, cloudlets.PolicyActivationStatusActive, "", 1, nil)
		expectDeletePolicyProperty(m, policyID, int64(idx*10), network, errDeleteProperty).Once()
		if errDeleteProperty != nil {
			return
		}
	}
}

func expectGetPolicyVersion(m *cloudlets.Mock, policyID, version int64, activations []cloudlets.PolicyActivation, err error) *mock.Call {
	return m.On(
		"GetPolicyVersion",
		mock.Anything,
		cloudlets.GetPolicyVersionRequest{PolicyID: policyID, Version: version, OmitRules: true},
	).Return(
		&cloudlets.PolicyVersion{
			Activations: activations,
		}, err)
}

func expectDeletePolicyProperty(m *cloudlets.Mock, policyID, propertyID int64, net cloudlets.PolicyActivationNetwork, expectError error) *mock.Call {
	return m.On("DeletePolicyProperty", mock.Anything, cloudlets.DeletePolicyPropertyRequest{
		PolicyID: policyID, PropertyID: propertyID, Network: net,
	}).Return(expectError)
}

func expectGetPolicyProperties(m *cloudlets.Mock, policyID int64, properties []string, expectError error) *mock.Call {
	if expectError != nil {
		return m.On("GetPolicyProperties", mock.Anything, cloudlets.GetPolicyPropertiesRequest{PolicyID: policyID}).Return(nil, expectError)
	}

	response := make(map[string]cloudlets.PolicyProperty)
	for i, name := range properties {
		response[name] = cloudlets.PolicyProperty{Name: name, ID: int64(10 * i), GroupID: int64(i)}
	}

	return m.On("GetPolicyProperties", mock.Anything, cloudlets.GetPolicyPropertiesRequest{PolicyID: policyID}).Return(response, nil)
}

func expectListPolicyActivations(m *cloudlets.Mock, policyID, version int64, network cloudlets.PolicyActivationNetwork, propertyNames []string, status cloudlets.PolicyActivationStatus, statusDetail string, numberActivations int, expectedErr error) *mock.Call {
	if expectedErr != nil {
		return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
			PolicyID: policyID,
			Network:  network,
		}).Return(nil, expectedErr)
	}

	policyActivations := createPolicyActivations(policyID, version, network, propertyNames, status, statusDetail, numberActivations)

	return m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	}).Return(policyActivations, nil)
}

func expectActivatePolicyVersion(m *cloudlets.Mock, policyID, version int64, network cloudlets.PolicyActivationNetwork, additionalProps []string, status cloudlets.PolicyActivationStatus, statusDetail string, numberActivations int, err error) *mock.Call {
	if err != nil {
		return m.On("ActivatePolicyVersion", mock.Anything, cloudlets.ActivatePolicyVersionRequest{
			PolicyID:                policyID,
			Async:                   true,
			Version:                 version,
			PolicyVersionActivation: cloudlets.PolicyVersionActivation{Network: network, AdditionalPropertyNames: additionalProps},
		}).Return(nil, err)
	}

	policyActivations := createPolicyActivations(policyID, version, network, additionalProps, status, statusDetail, numberActivations)

	return m.On("ActivatePolicyVersion", mock.Anything, cloudlets.ActivatePolicyVersionRequest{
		PolicyID:                policyID,
		Async:                   true,
		Version:                 version,
		PolicyVersionActivation: cloudlets.PolicyVersionActivation{Network: network, AdditionalPropertyNames: additionalProps},
	}).Return(policyActivations, err)
}

func createPolicyActivations(policyID, version int64, network cloudlets.PolicyActivationNetwork, properties []string, status cloudlets.PolicyActivationStatus, statusDetail string, numberActivations int) []cloudlets.PolicyActivation {
	policyActivations := make([]cloudlets.PolicyActivation, 0)

	resultNetwork := network
	if network == "" {
		resultNetwork = cloudlets.PolicyActivationNetworkProduction
	}

	s := cloudlets.PolicyActivationStatusDeactivated
	for i := 0; i < numberActivations; i++ {
		if i == numberActivations-1 {
			s = status
		}
		for _, propertyName := range properties {
			policyActivations = append(policyActivations, cloudlets.PolicyActivation{
				APIVersion:   "1.0",
				Network:      resultNetwork,
				PolicyInfo:   cloudlets.PolicyInfo{PolicyID: policyID, Status: s, Version: version, StatusDetail: statusDetail},
				PropertyInfo: cloudlets.PropertyInfo{Name: propertyName},
			})
		}
	}

	return policyActivations
}

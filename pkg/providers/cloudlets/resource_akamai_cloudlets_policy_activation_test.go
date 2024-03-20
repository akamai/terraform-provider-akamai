package cloudlets

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCloudletsPolicyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*cloudlets.Mock)
		steps []resource.TestStep
	}{
		"try to create activation with a non existing policy version": {
			init: func(m *cloudlets.Mock) {
				expectToDiscoverPolicyAsV2(m, 1234)
				m.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{PolicyID: 1234, Version: 1, OmitRules: true}).Return(nil, fmt.Errorf("an error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
					),
				},
			},
		},
		"create and read activation, version == 1, failed with property error -> retry -> activate": {
			init: func(m *cloudlets.Mock) {
				policyID, version, network, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				propertyNotFoundError := fmt.Errorf(`"detail": "Requested propertyName \"test.property.name\" does not exist"`)
				// create
				expectToDiscoverPolicyAsV2(m, 1234)
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectToDiscoverPolicyAsV2(m, 1234)
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
			},
		},
		"create and read activation, version == 1, active activations on activate -> empty list activations -> retry twice -> correct list activations": {
			init: func(m *cloudlets.Mock) {
				policyID, version, network, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				expectToDiscoverPolicyAsV2(m, 1234)
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Once()
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				// empty activations
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// retry # 1
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// retry # 2
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectToDiscoverPolicyAsV2(m, policyID)
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
				// empty activations retry # 1
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 2
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 3
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 4
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
				// empty activations retry # 5
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: policy activation: policyID 1234"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> wait -> failed": {
			init: func(m *cloudlets.Mock) {
				policyID, version, staging, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				expectToDiscoverPolicyAsV2(m, policyID)
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: policy activation: policyID 1234"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> several iterations in read function": {
			init: func(m *cloudlets.Mock) {
				policyID, version, staging, properties := int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}
				// create
				expectToDiscoverPolicyAsV2(m, policyID)
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectToDiscoverPolicyAsV2(m, policyID)
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
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v1, staging, []string{}, "", "", 1, nil).Once()

				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_production.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "0"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectToDiscoverPolicyAsV2(m, 1234)
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{
					{APIVersion: "1.0", Network: "prod", PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: 1234, Version: 1, Status: cloudlets.PolicyActivationStatusInactive,
					}},
				}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active policy activation -> error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectToDiscoverPolicyAsV2(m, 1234)
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV2(m, policyID)
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectToDiscoverPolicyAsV2(m, policyID)
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation read: an error"),
				},
			},
		},
		"create and read activation, update - no changes, so skip update": {
			init: func(m *cloudlets.Mock) {
				staging, properties, policyID, v1, active := cloudlets.PolicyActivationNetworkStaging, []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationStatusActive
				// 1 - for policy_activation_version1.tf
				// create
				expectToDiscoverPolicyAsV2(m, policyID)
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectToDiscoverPolicyAsV2(m, 1234)
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusFailed, "activation failed", 1, fmt.Errorf("activation failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: activation failed"),
				},
			},
		},
		"create activation - failed activation while polling with no failed status": {
			init: func(m *cloudlets.Mock) {
				// create
				expectToDiscoverPolicyAsV2(m, 1234)
				expectGetPolicyVersion(m, 1234, 1, []cloudlets.PolicyActivation{}, nil).Once()
				expectActivatePolicyVersion(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, 1234, 1, "staging", []string{"prp_0", "prp_1"}, cloudlets.PolicyActivationStatusDeactivated, "failed", 1, fmt.Errorf("activation failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
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
				expectListPolicyActivations(m, policyID, v1, "", properties, active, "", 1, nil).Once()
				expectActivatePolicyVersion(m, policyID, v2, staging, properties, cloudlets.PolicyActivationStatusPending, "", 1, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// read
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				expectListPolicyActivations(m, policyID, v2, staging, properties, active, "", 1, nil).Once()
				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1_prod.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
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
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: staging, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: v2, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				// update
				// poll until active -> waitForPolicyActivation()
				expectListPolicyActivations(m, policyID, v2, staging, []string{"prp_0", "prp_1"}, active, "", 1, nil).Once()

				// read
				expectListPolicyActivations(m, policyID, v2, staging, []string{"prp_0", "prp_1"}, active, "", 0, nil).Times(6)

				// delete
				expectDeletePhase(m, policyID, properties, nil, staging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile(`policy activation read: cannot find any activation for the given policy \(1234\) and network \('staging'\)`),
				},
			},
		},
		"import - success": {
			init: func(m *cloudlets.Mock) {
				// create
				properties, policyID, version, network := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging
				expectToDiscoverPolicyAsV2(m, policyID)
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Times(1)
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// already active - no need to poll
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(2)
				// import
				expectToDiscoverPolicyAsV2(m, policyID)
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
				},
			},
		},
		"import - only deactivated activation - expect an error": {
			init: func(m *cloudlets.Mock) {
				// create
				properties, policyID, version, network := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging
				expectToDiscoverPolicyAsV2(m, policyID)
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Times(1)
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// already active - no need to poll
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(2)
				// import
				expectToDiscoverPolicyAsV2(m, policyID)
				policyActivations := []cloudlets.PolicyActivation{
					{
						APIVersion: "1.0",
						Network:    "staging",
						PolicyInfo: cloudlets.PolicyInfo{
							PolicyID: 1234,
							Name:     "Policy1",
							Version:  2,
							Status:   cloudlets.PolicyActivationStatusDeactivated,
						},
						PropertyInfo: cloudlets.PropertyInfo{
							Name: "prp_0",
						},
					},
					{
						APIVersion: "1.0",
						Network:    "staging",
						PolicyInfo: cloudlets.PolicyInfo{
							PolicyID: 1234,
							Name:     "Policy1",
							Version:  3,
							Status:   cloudlets.PolicyActivationStatusDeactivated,
						},
						PropertyInfo: cloudlets.PropertyInfo{
							Name: "prp_1",
						},
					},
				}
				m.On("ListPolicyActivations", mock.Anything, cloudlets.ListPolicyActivationsRequest{
					PolicyID: policyID,
					Network:  network,
				}).Return(policyActivations, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: no active activation has been found for policy_id: '1234' and network: 'staging'`),
				},
			},
		},
		"import - empty activations - expect an error": {
			init: func(m *cloudlets.Mock) {
				// create
				properties, policyID, version, network := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging
				expectToDiscoverPolicyAsV2(m, policyID)
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Times(1)
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// already active - no need to poll
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(2)
				// import - expect an error
				expectToDiscoverPolicyAsV2(m, policyID)
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Times(1)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: no active activation has been found for policy_id: '1234' and network: 'staging'`),
				},
			},
		},
		"import - wrong import ID - expect an error": {
			init: func(m *cloudlets.Mock) {
				// create
				properties, policyID, version, network := []string{"prp_0", "prp_1"}, int64(1234), int64(1), cloudlets.PolicyActivationNetworkStaging
				expectToDiscoverPolicyAsV2(m, policyID)
				activations := make([]cloudlets.PolicyActivation, len(properties))
				for _, p := range properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: policyID, Version: version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				expectGetPolicyVersion(m, policyID, version, activations, nil).Times(1)
				expectActivatePolicyVersion(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// already active - no need to poll
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(1)
				// read
				expectListPolicyActivations(m, policyID, version, network, properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(2)
				// delete
				expectDeletePhase(m, 1234, []string{"prp_0", "prp_1"}, nil, cloudlets.PolicyActivationNetworkStaging, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyActivation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.PolicyActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "wrong_import_id",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: import id should be of format: <policy_id>:<network>, for example: 1234:staging`),
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourceV3CloudletsPolicyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*cloudlets.Mock, *v3.Mock)
		steps []resource.TestStep
	}{
		"try to create activation with a non existing policy": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID := int64(1234)
				m2.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{PolicyID: policyID}).Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
				m3.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{PolicyID: policyID}).Return(nil, &v3.Error{Status: http.StatusNotFound}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile(`could not get policy 1234`),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork
				expectFullV3Activation(m2, m3, policyID, version, network)
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"(V3 policy but V2 throws 500s) create and read activation, version == 1, inactive -> activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork
				// create
				//discover
				m2.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{PolicyID: policyID}).Return(nil, &cloudlets.Error{StatusCode: http.StatusInternalServerError}).Once()
				m3.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{PolicyID: policyID}).Return(&v3.Policy{ID: policyID}, nil).Once()
				//rest of create
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, 111, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, 111, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, pending activations on activate -> retry -> activated": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				// create
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusInProgress, v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, active activations on activate -> retry twice -> activated": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				// create
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusInProgress, v3.ActivationStatusInProgress, v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate -> wait -> failed": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				// create
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusFailed}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/shared_policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: activation failed for policy 1234"),
				},
			},
		},
		"cannot modify 'associated_properties' for shared policy": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID := int64(1234)
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version_invalid.tf"),
					ExpectError: regexp.MustCompile("cannot provide 'associated_properties' for shared policy"),
				},
			},
		},

		"create and read activation, version == 1, production, inactive -> activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.ProductionNetwork
				expectFullV3Activation(m2, m3, policyID, version, network)
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1_production.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "prod"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "0"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, prod, inactive -> activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.ProductionNetwork
				expectFullV3Activation(m2, m3, policyID, version, network)
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "prod"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "0"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				// create
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate -> 500 error (retry) -> activated": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				// create
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork

				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, &v3.Error{Status: http.StatusInternalServerError, Title: "something broke"}).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate -> get active policy activation -> error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				// create
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{}, fmt.Errorf("an error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation create: an error"),
				},
			},
		},

		"create and read activation, version == 1, active -> read (version is already active)": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork

				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()

				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Times(2)
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"create and read activation, version == 1, inactive -> activate -> read error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork

				// create
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{}, fmt.Errorf("an error")).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					ExpectError: regexp.MustCompile("policy activation read: an error"),
				},
			},
		},

		"create and read activation, update - no changes, so skip update": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, activationID, network := int64(1234), int64(1), int64(111), v3.StagingNetwork
				// 1 - for policy_activation_version1.tf
				// create
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, activationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, activationID, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Twice()
				// 2 - for policy_activation_version1.tf
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Twice()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"Create and read activation. Update: version not active, activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				network, policyID, version1, version2, reactivationID := v3.StagingNetwork, int64(1234), int64(1), int64(2), int64(222)
				// 1 - for policy_activation_version1.tf
				expectFullV3Activation(m2, m3, policyID, version1, network)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				// update
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version2, reactivationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, reactivationID, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version2, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version2, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version2, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_update_version2.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"Create and read activation. Update: change version from staging to prod, activate": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				network1, network2, policyID, version, reactivationID := v3.StagingNetwork, v3.ProductionNetwork, int64(1234), int64(1), int64(222)
				// 1 - for policy_activation_version1.tf
				expectFullV3Activation(m2, m3, policyID, version, network1)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network1), nil).Once()
				// update
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network1), nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version, reactivationID, network2, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, reactivationID, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
				// read
				bothNetworks := v3.CurrentActivations{
					Production: v3.ActivationInfo{Effective: preparePolicyActivation(policyID, version, network2)},
					Staging:    v3.ActivationInfo{Effective: preparePolicyActivation(policyID, version, network1)},
				}
				expectGetV3Policy(m3, policyID, bothNetworks, nil).Once()
				// read
				expectGetV3Policy(m3, policyID, bothNetworks, nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "prod"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
			},
		},

		"Create and read activation. Update: isReactivationNeeded error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				network, policyID, version := v3.StagingNetwork, int64(1234), int64(1)
				// 1 - for policy_activation_version1.tf
				expectFullV3Activation(m2, m3, policyID, version, network)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// update
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{}, fmt.Errorf("an error")).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},

		"Create and read activation. Update: activate policy version error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				network, policyID, version1, version2, reactivationID := v3.StagingNetwork, int64(1234), int64(1), int64(2), int64(222)
				// 1 - for policy_activation_version1.tf
				expectFullV3Activation(m2, m3, policyID, version1, network)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				// update
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version2, reactivationID, network, fmt.Errorf("an error")).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version1, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},

		"Create and read activation. Update: error while polling": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				network, policyID, version1, version2, reactivationID := v3.StagingNetwork, int64(1234), int64(1), int64(2), int64(222)
				// 1 - for policy_activation_version1.tf
				expectFullV3Activation(m2, m3, policyID, version1, network)
				// 2 - for policy_activation_update_version2.tf
				// refresh read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				// update
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version1, network), nil).Once()
				expectActivateV3PolicyVersion(m3, policyID, version2, reactivationID, network, nil).Once()
				// poll until active -> waitForPolicyActivation()
				expectWaitForV3Activation(m3, policyID, reactivationID, []v3.ActivationStatus{}, fmt.Errorf("an error"))
				// delete
				expectV3DeletePhase(m3, policyID, version2, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_update_version2.tf"),
					ExpectError: regexp.MustCompile("policy activation update: an error"),
				},
			},
		},

		"import - success": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork

				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()

				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Times(2)
				// import
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
				},
			},
		},

		"import - only deactivated activation - expect an error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork

				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()

				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Times(2)
				// import
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{
					Production: v3.ActivationInfo{},
					Staging: v3.ActivationInfo{
						Effective: &v3.PolicyActivation{
							PolicyID:      policyID,
							PolicyVersion: version,
							Network:       network,
							Operation:     v3.OperationDeactivation,
							Status:        v3.ActivationStatusSuccess,
						},
					},
				}, nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: no active activation has been found for policy_id: '1234' and network: 'staging'`),
				},
			},
		},

		"import - empty activations - expect an error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork

				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()

				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Times(2)
				// import
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, v3.CurrentActivations{
					Production: v3.ActivationInfo{},
					Staging:    v3.ActivationInfo{},
				}, nil).Once()
				// delete
				expectV3DeletePhase(m3, policyID, version, network)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "1234:staging",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: no active activation has been found for policy_id: '1234' and network: 'staging'`),
				},
			},
		},

		"import - wrong import ID - expect an error": {
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock) {
				policyID, version, network := int64(1234), int64(1), v3.StagingNetwork

				// create, policy active so no need to activate
				expectToDiscoverPolicyAsV3(m2, m3, policyID)
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()

				// read
				expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Times(2)
				// delete
				expectV3DeletePhase(m3, policyID, version, network)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResCloudletsPolicyV3Activation/policy_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(v3.ActivationStatusSuccess)),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "timeouts.0.default", "2h"),
						resource.TestCheckResourceAttr("akamai_cloudlets_policy_activation.test", "is_shared", "true"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "wrong_import_id",
					ResourceName:      "akamai_cloudlets_policy_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Error: import id should be of format: <policy_id>:<network>, for example: 1234:staging`),
				},
			},
		},
	}

	// redefining times to accelerate tests
	ActivationPollMinimum, ActivationPollInterval, PolicyActivationRetryPollMinimum = time.Millisecond, time.Millisecond, time.Millisecond

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			v2Client := &cloudlets.Mock{}
			v3client := &v3.Mock{}
			test.init(v2Client, v3client)
			useClientV2AndV3(v2Client, v3client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			v2Client.AssertExpectations(t)
			v3client.AssertExpectations(t)
		})
	}
}

// expect full activation of a policy (creation)
func expectFullActivation(m *cloudlets.Mock, policyID, version int64, properties []string, network cloudlets.PolicyActivationNetwork, times int) {
	// create
	expectToDiscoverPolicyAsV2(m, policyID)
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

func expectToDiscoverPolicyAsV2(m *cloudlets.Mock, policyID int64) {
	m.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{PolicyID: policyID}).Return(&cloudlets.Policy{PolicyID: policyID}, nil).Once()
}

func expectToDiscoverPolicyAsV3(m2 *cloudlets.Mock, m3 *v3.Mock, policyID int64) {
	m2.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{PolicyID: policyID}).Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
	m3.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{PolicyID: policyID}).Return(&v3.Policy{ID: policyID}, nil).Once()
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

func expectFullV3Activation(m2 *cloudlets.Mock, m3 *v3.Mock, policyID, version int64, network v3.Network) {
	// create
	expectToDiscoverPolicyAsV3(m2, m3, policyID)
	expectGetV3Policy(m3, policyID, v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}, nil).Once()
	expectActivateV3PolicyVersion(m3, policyID, version, 111, network, nil).Once()
	// poll until active -> waitForPolicyActivation()
	expectWaitForV3Activation(m3, policyID, 111, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
	// read
	expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
	// read
	expectGetV3Policy(m3, policyID, prepareActivatedResponseForNetwork(policyID, version, network), nil).Once()
}

func prepareActivatedResponseForNetwork(policyID, version int64, network v3.Network) v3.CurrentActivations {
	activations := v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}}
	if network == v3.StagingNetwork {
		activations.Staging.Effective = preparePolicyActivation(policyID, version, network)
	} else {
		activations.Production.Effective = preparePolicyActivation(policyID, version, network)
	}
	return activations
}

func preparePolicyActivation(policyID, version int64, network v3.Network) *v3.PolicyActivation {
	return &v3.PolicyActivation{
		PolicyID:      policyID,
		PolicyVersion: version,
		Network:       network,
		Status:        v3.ActivationStatusSuccess,
		Operation:     v3.OperationActivation,
	}
}

func expectGetV3Policy(m *v3.Mock, policyID int64, activations v3.CurrentActivations, err error) *mock.Call {
	if err != nil {
		return m.On(
			"GetPolicy",
			mock.Anything,
			v3.GetPolicyRequest{PolicyID: policyID},
		).Return(nil, err)
	}
	return m.On(
		"GetPolicy",
		mock.Anything,
		v3.GetPolicyRequest{PolicyID: policyID},
	).Return(
		&v3.Policy{
			ID:                 policyID,
			CurrentActivations: activations,
		}, nil)
}

func expectActivateV3PolicyVersion(m *v3.Mock, policyID, version, activationID int64, network v3.Network, err error) *mock.Call {
	if err != nil {
		return m.On("ActivatePolicy", mock.Anything, v3.ActivatePolicyRequest{
			PolicyID:      policyID,
			Network:       network,
			PolicyVersion: version,
		}).Return(nil, err)
	}

	return m.On("ActivatePolicy", mock.Anything, v3.ActivatePolicyRequest{
		PolicyID:      policyID,
		Network:       network,
		PolicyVersion: version,
	}).Return(&v3.PolicyActivation{
		ID:        activationID,
		Status:    v3.ActivationStatusInProgress,
		Network:   network,
		Operation: v3.OperationActivation,
	}, err)
}

func expectWaitForV3Activation(m *v3.Mock, policyID, activationID int64, activations []v3.ActivationStatus, error error) {
	if error != nil {
		m.On("GetPolicyActivation",
			mock.Anything,
			v3.GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: activationID,
			}).Return(nil, error).Once()
		return
	}
	for idx := range activations {
		m.On("GetPolicyActivation",
			mock.Anything,
			v3.GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: activationID,
			}).Return(&v3.PolicyActivation{
			PolicyID: policyID,
			Status:   activations[idx],
		}, nil).Once()
	}
}

func expectDeactivateV3PolicyVersion(m *v3.Mock, policyID, version, activationID int64, network v3.Network, err error) *mock.Call {
	if err != nil {
		return m.On("DeactivatePolicy", mock.Anything, v3.DeactivatePolicyRequest{
			PolicyID:      policyID,
			Network:       network,
			PolicyVersion: version,
		}).Return(nil, err)
	}

	return m.On("DeactivatePolicy", mock.Anything, v3.DeactivatePolicyRequest{
		PolicyID:      policyID,
		Network:       network,
		PolicyVersion: version,
	}).Return(&v3.PolicyActivation{
		ID:        activationID,
		Status:    v3.ActivationStatusInProgress,
		Network:   network,
		Operation: v3.OperationActivation,
	}, err)
}

func expectV3DeletePhase(m *v3.Mock, policyID, version int64, network v3.Network) {
	expectDeactivateV3PolicyVersion(m, policyID, version, 333, network, nil).Once()
	expectWaitForV3Activation(m, policyID, 333, []v3.ActivationStatus{v3.ActivationStatusSuccess}, nil)
}

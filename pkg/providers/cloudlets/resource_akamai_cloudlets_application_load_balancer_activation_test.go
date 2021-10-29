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

func TestResourceCloudletsApplicationLoadBalancerActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*mockcloudlets)
		steps []resource.TestStep
	}{
		"create and read activation, version == 1, inactive -> activate -> second attempt": {
			init: func(m *mockcloudlets) {
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.ActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.ActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1 (many statuses), inactive -> activate -> second attempt": {
			init: func(m *mockcloudlets) {
				inactive := cloudlets.ActivationResponse{
					ActivatedDate: "2021-10-29T00:00:10.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.ActivationStatusInactive,
					Version:       1,
				}
				deactivated := cloudlets.ActivationResponse{
					ActivatedDate: "2021-10-29T00:00:20.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.ActivationStatusDeactivated,
					Version:       1,
				}
				pending := cloudlets.ActivationResponse{
					ActivatedDate: "2021-10-29T00:00:30.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.ActivationStatusPending,
					Version:       1,
				}
				active := cloudlets.ActivationResponse{
					ActivatedDate: "2021-10-29T00:00:40.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.ActivationStatusActive,
					Version:       1,
				}

				// create
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated}, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, pending}, nil).Once()
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, active}, nil).Once()
				// read
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, active}, nil).Once()
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, active}, nil).Once()
				// read
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, active}, nil).Once()
				expectGetLoadBalancerActivationsMany(m, "org_1", cloudlets.ActivationsList{inactive, deactivated, active}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.ActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.ActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active application load balancer activation -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: an error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *mockcloudlets) {
				// create, alb active so no need to activate
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.ActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> read -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation read: an error"),
				},
			},
		},
		"create and read activation, update - no changes": {
			init: func(m *mockcloudlets) {
				// first test step
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// update
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: version already active, read": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()

				// 2 - for alb_activation_version1.tf
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: version not active, activate": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: production version not active, activate new resource": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// create new resource
				expectGetLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.ActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "id", "org_1:STAGING"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "id", "org_1:PRODUCTION"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "PRODUCTION"),
					),
				},
			},
		},
		"create and read activation. Update: GetLoadBalancerActivations error": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusActive, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: an error"),
				},
			},
		},
		"create and read activation. Update: create application load balancer version error": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// create: poll after activation
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusPending, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectGetLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.ActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.ActivationStatusPending, fmt.Errorf("an error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.StatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: an error"),
				},
			},
		},
	}

	// redefining times to run the tests faster
	ALBActivationPollMinimum = time.Millisecond * 1
	ALBActivationPollInterval = time.Millisecond * 1

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

var (
	expectGetLoadBalancerActivations = func(m *mockcloudlets, originID string, version int64, network cloudlets.ActivationNetwork, status cloudlets.ActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetLoadBalancerActivations",
				mock.Anything,
				originID,
			).Return(nil, err)
		}
		return m.On(
			"GetLoadBalancerActivations",
			mock.Anything,
			originID,
		).Return(
			cloudlets.ActivationsList{
				cloudlets.ActivationResponse{
					Network:  network,
					OriginID: originID,
					Status:   status,
					Version:  version,
				},
			}, nil)
	}

	expectGetLoadBalancerActivationsMany = func(m *mockcloudlets, originID string, activations cloudlets.ActivationsList, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetLoadBalancerActivations",
				mock.Anything,
				originID,
			).Return(nil, err)
		}
		return m.On(
			"GetLoadBalancerActivations",
			mock.Anything,
			originID,
		).Return(activations, nil)
	}

	expectActivateLoadBalancerVersion = func(m *mockcloudlets, originID string, version int64, network cloudlets.ActivationNetwork, status cloudlets.ActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ActivateLoadBalancerVersion",
				mock.Anything,
				cloudlets.ActivateLoadBalancerVersionRequest{
					OriginID: originID,
					Async:    true,
					ActivationRequest: cloudlets.ActivationRequestParams{
						Network: network,
						DryRun:  false,
						Version: version,
					},
				},
			).Return(nil, err)
		}
		return m.On(
			"ActivateLoadBalancerVersion",
			mock.Anything,
			cloudlets.ActivateLoadBalancerVersionRequest{
				OriginID: originID,
				Async:    true,
				ActivationRequest: cloudlets.ActivationRequestParams{
					Network: network,
					DryRun:  false,
					Version: version,
				},
			},
		).Return(
			&cloudlets.ActivationResponse{
				Network:  network,
				OriginID: originID,
				Status:   status,
				Version:  version,
			}, nil)
	}
)

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
	anError := fmt.Errorf("an error")
	tests := map[string]struct {
		init  func(*mockcloudlets)
		steps []resource.TestStep
	}{
		"create and read activation, version == 1, inactive -> activate -> second attempt": {
			init: func(m *mockcloudlets) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1 (many statuses), inactive -> activate -> second attempt": {
			init: func(m *mockcloudlets) {
				inactive := cloudlets.LoadBalancerActivation{
					ActivatedDate: "2021-10-29T00:00:10.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.LoadBalancerActivationStatusInactive,
					Version:       1,
				}
				deactivated := cloudlets.LoadBalancerActivation{
					ActivatedDate: "2021-10-29T00:00:20.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.LoadBalancerActivationStatusDeactivated,
					Version:       1,
				}
				pending := cloudlets.LoadBalancerActivation{
					ActivatedDate: "2021-10-29T00:00:30.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.LoadBalancerActivationStatusPending,
					Version:       1,
				}
				active := cloudlets.LoadBalancerActivation{
					ActivatedDate: "2021-10-29T00:00:40.000Z",
					Network:       "STAGING",
					OriginID:      "org_1",
					Status:        cloudlets.LoadBalancerActivationStatusActive,
					Version:       1,
				}

				// create
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated}, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, pending}, nil).Once()
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, active}, nil).Once()
				// read
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, active}, nil).Once()
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, active}, nil).Once()
				// read
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, active}, nil).Once()
				expectListLoadBalancerActivationsMany(m, "org_1", []cloudlets.LoadBalancerActivation{inactive, deactivated, active}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: application load balancer activation failed. No changes were written to server:\nan error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active application load balancer activation -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: error while waiting until load balancer activation status == 'active':\nan error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *mockcloudlets) {
				// create, alb active so no need to activate
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> read -> error": {
			init: func(m *mockcloudlets) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
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
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// update
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
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
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()

				// 2 - for alb_activation_version1.tf
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
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
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
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
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// create new resource
				expectListLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "PRODUCTION", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "id", "org_1:STAGING"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update_prod.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "id", "org_1:PRODUCTION"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "PRODUCTION"),
					),
				},
			},
		},
		"create and read activation. Update: ListLoadBalancerActivations error": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: error while waiting until load balancer activation status == 'active':\nan error"),
				},
			},
		},
		"create and read activation. Update: create application load balancer version error": {
			init: func(m *mockcloudlets) {
				// 1 - for alb_activation_version1.tf
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// create: poll after activation
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()

				// 2 - for alb_activation_update.tf
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 2, "STAGING", cloudlets.LoadBalancerActivationStatusPending, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: application load balancer activation failed. No changes were written to server:\nan error"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
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
	expectListLoadBalancerActivations = func(m *mockcloudlets, originID string, version int64, network cloudlets.LoadBalancerActivationNetwork, status cloudlets.LoadBalancerActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ListLoadBalancerActivations",
				mock.Anything,
				cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
			).Return(nil, err)
		}
		return m.On(
			"ListLoadBalancerActivations",
			mock.Anything,
			cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
		).Return(
			[]cloudlets.LoadBalancerActivation{
				{
					Network:  network,
					OriginID: originID,
					Status:   status,
					Version:  version,
				},
			}, nil)
	}

	expectListLoadBalancerActivationsMany = func(m *mockcloudlets, originID string, activations []cloudlets.LoadBalancerActivation, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ListLoadBalancerActivations",
				mock.Anything,
				cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
			).Return(nil, err)
		}
		return m.On(
			"ListLoadBalancerActivations",
			mock.Anything,
			cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
		).Return(activations, nil)
	}

	expectActivateLoadBalancerVersion = func(m *mockcloudlets, originID string, version int64, network cloudlets.LoadBalancerActivationNetwork, status cloudlets.LoadBalancerActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ActivateLoadBalancerVersion",
				mock.Anything,
				cloudlets.ActivateLoadBalancerVersionRequest{
					OriginID: originID,
					Async:    true,
					LoadBalancerVersionActivation: cloudlets.LoadBalancerVersionActivation{
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
				LoadBalancerVersionActivation: cloudlets.LoadBalancerVersionActivation{
					Network: network,
					DryRun:  false,
					Version: version,
				},
			},
		).Return(
			&cloudlets.LoadBalancerActivation{
				Network:  network,
				OriginID: originID,
				Status:   status,
				Version:  version,
			}, nil)
	}
)

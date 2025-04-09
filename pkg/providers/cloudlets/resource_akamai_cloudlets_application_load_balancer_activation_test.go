package cloudlets

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCloudletsApplicationLoadBalancerActivation(t *testing.T) {
	anError := fmt.Errorf("an error")
	originNotDefinedError := fmt.Errorf(`"detail": "Origin 'origin-test-1' is not defined in Property Manager for this network"`)
	tests := map[string]struct {
		init  func(*cloudlets.Mock)
		steps []resource.TestStep
	}{
		"create and read activation, version == 1, inactive -> activate -> second attempt": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, origin not defined error -> successful second attempt": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, originNotDefinedError).Once()
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1 (many statuses), inactive -> activate -> second attempt": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: application load balancer activation failed. No changes were written to server:\nan error"),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> get active application load balancer activation -> error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, anError).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation create: an error occurred while waiting for the load balancer activation status == 'active':\nan error"),
				},
			},
		},
		"create and read activation, version == 1, active -> read": {
			init: func(m *cloudlets.Mock) {
				// create, alb active so no need to activate
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
					),
				},
			},
		},
		"create and read activation, version == 1, active -> read with defined timeout": {
			init: func(m *cloudlets.Mock) {
				// create, alb active so no need to activate
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_timeouts.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "timeouts.0.default", "2h"),
					),
				},
			},
		},
		"create and read activation, version == 1, inactive -> activate -> read -> error": {
			init: func(m *cloudlets.Mock) {
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation read: an error"),
				},
			},
		},
		"create and read activation, update - no changes": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: version already active, read": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: version not active, activate": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation. Update: production version not active, activate new resource": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "id", "org_1:STAGING"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update_prod.tf"),
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
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: an error occurred while waiting for the load balancer activation status == 'active':\nan error"),
				},
			},
		},
		"create and read activation. Update: create application load balancer version error": {
			init: func(m *cloudlets.Mock) {
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_update.tf"),
					ExpectError: regexp.MustCompile("application load balancer activation update: application load balancer activation failed. No changes were written to server:\nan error"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"import - success": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// import
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Times(2)
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "org_1,STAGING,1",
					ResourceName:      "akamai_cloudlets_application_load_balancer_activation.test",
					ImportStateVerify: true,
				},
			},
		},
		"import - wrong key - expect an error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// import
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "wrong_import_id",
					ResourceName:      "akamai_cloudlets_application_load_balancer_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`the import ID has to be a comma separated list of the origin ID, network, and version`),
				},
			},
		},
		"import - empty activation - expect an error": {
			init: func(m *cloudlets.Mock) {
				// create
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusInactive, nil).Once()
				expectActivateLoadBalancerVersion(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// read
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusPending, nil).Once()
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Once()
				// import
				expectListLoadBalancerActivations(m, "org_1", 1, "STAGING", cloudlets.LoadBalancerActivationStatusActive, nil).Times(2)
				expectListLoadBalancerActivationsEmpty(m, "org_1", nil).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceCloudletsApplicationLoadBalancerActivation/alb_activation_version1.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckOutput("status", string(cloudlets.LoadBalancerActivationStatusActive)),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer_activation.test", "network", "STAGING"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "org_1,STAGING,1",
					ResourceName:      "akamai_cloudlets_application_load_balancer_activation.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`application load balancer activation: application load balancer activation version not found`),
				},
			},
		},
	}

	// redefining times to run the tests faster
	ALBActivationPollMinimum = time.Millisecond * 1
	ALBActivationPollInterval = time.Millisecond * 1

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

var (
	expectListLoadBalancerActivations = func(m *cloudlets.Mock, originID string, version int64, network cloudlets.LoadBalancerActivationNetwork, status cloudlets.LoadBalancerActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ListLoadBalancerActivations",
				testutils.MockContext,
				cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
			).Return(nil, err)
		}
		return m.On(
			"ListLoadBalancerActivations",
			testutils.MockContext,
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

	expectListLoadBalancerActivationsEmpty = func(m *cloudlets.Mock, originID string, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ListLoadBalancerActivations",
				testutils.MockContext,
				cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
			).Return(nil, err)
		}
		return m.On(
			"ListLoadBalancerActivations",
			testutils.MockContext,
			cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
		).Return(
			[]cloudlets.LoadBalancerActivation{}, nil)
	}

	expectListLoadBalancerActivationsMany = func(m *cloudlets.Mock, originID string, activations []cloudlets.LoadBalancerActivation, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ListLoadBalancerActivations",
				testutils.MockContext,
				cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
			).Return(nil, err)
		}
		return m.On(
			"ListLoadBalancerActivations",
			testutils.MockContext,
			cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID},
		).Return(activations, nil)
	}

	expectActivateLoadBalancerVersion = func(m *cloudlets.Mock, originID string, version int64, network cloudlets.LoadBalancerActivationNetwork, status cloudlets.LoadBalancerActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"ActivateLoadBalancerVersion",
				testutils.MockContext,
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
			testutils.MockContext,
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

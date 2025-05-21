package cloudlets

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataApplicationLoadBalancer(t *testing.T) {
	tests := map[string]struct {
		configPath     string
		checkFunctions []resource.TestCheckFunc
		withError      *regexp.Regexp
		init           func(*cloudlets.Mock)
	}{
		"validate schema": {
			configPath: "testdata/TestDataCloudletsApplicationLoadBalancer/application_load_balancer.tf",
			init: func(m *cloudlets.Mock) {
				m.On("GetOrigin", testutils.MockContext, mock.Anything).Return(&cloudlets.Origin{
					OriginID:  "alb_test_krk_dc1",
					Akamaized: false,
					Checksum:  "9c0fc1f3e9ea7eb2e090f2bf53709e45",
					Type:      "APPLICATION_LOAD_BALANCER",
				}, nil)
				m.On("GetLoadBalancerVersion", testutils.MockContext, mock.Anything).Return(&cloudlets.LoadBalancerVersion{
					BalancingType: "WEIGHTED",
					CreatedBy:     "jbond",
					CreatedDate:   "2021-09-27T11:50:07.715Z",
					DataCenters: []cloudlets.DataCenter{
						{
							CloudServerHostHeaderOverride: false,
							CloudService:                  false,
							Continent:                     "NA",
							Country:                       "US",
							City:                          "Cambridge",
							Hostname:                      "example.com",
							Latitude:                      ptr.To(102.78108),
							StateOrProvince:               ptr.To("MA"),
							LivenessHosts: []string{
								"clorigin3.www.example.com",
							},
							Longitude: ptr.To(-116.07064),

							OriginID: "alb_test_krk_dc1",
							Percent:  ptr.To(100.0),
						},
					},
					Deleted:          false,
					Description:      "Test load balancing configuration.",
					Immutable:        false,
					LastModifiedBy:   "jbond",
					LastModifiedDate: "2021-09-27T11:50:07.715Z",
					LivenessSettings: &cloudlets.LivenessSettings{
						HostHeader:        "clorigin3.www.example.com",
						Interval:          25,
						Path:              "/status",
						Port:              443,
						Protocol:          "HTTPS",
						Status3xxFailure:  false,
						Status4xxFailure:  true,
						Status5xxFailure:  false,
						Timeout:           30.0,
						AdditionalHeaders: map[string]string{"Authorization": "test"},
						RequestString:     "test",
						ResponseString:    "test",
					},
					OriginID: "alb_test_krk_dc1",
					Version:  2,
					Warnings: []cloudlets.Warning{
						{
							Detail:      "Data center 1 origin 'alb_test_krk_dc1' hostname is empty",
							Title:       "Validation Warning",
							Type:        "/cloudlets/error-types/validation-warning",
							JSONPointer: "/",
						},
						{
							Detail:      "The total data center percentage (weight) must be 100%",
							Title:       "Validation Warning",
							Type:        "/cloudlets/error-types/validation-warning",
							JSONPointer: "/",
						},
					},
				}, nil)
				m.On("ListLoadBalancerVersions", testutils.MockContext, mock.Anything).Return([]cloudlets.LoadBalancerVersion{
					{
						Version: 1,
					},
					{
						Version: 2,
					},
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "origin_id", "alb_test_krk_dc1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "version", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "description", "Test load balancing configuration."),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "type", "APPLICATION_LOAD_BALANCER"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "balancing_type", "WEIGHTED"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "created_by", "jbond"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "created_date", "2021-09-27T11:50:07.715Z"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "deleted", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "immutable", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "last_modified_by", "jbond"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "last_modified_date", "2021-09-27T11:50:07.715Z"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "warnings", "[\n {\n  \"detail\": \"Data center 1 origin 'alb_test_krk_dc1' hostname is empty\",\n  \"jsonPointer\": \"/\",\n  \"title\": \"Validation Warning\",\n  \"type\": \"/cloudlets/error-types/validation-warning\"\n },\n {\n  \"detail\": \"The total data center percentage (weight) must be 100%\",\n  \"jsonPointer\": \"/\",\n  \"title\": \"Validation Warning\",\n  \"type\": \"/cloudlets/error-types/validation-warning\"\n }\n]"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.cloud_service", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.continent", "NA"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.country", "US"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.latitude", "102.78108"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.longitude", "-116.07064"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.percent", "100"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.city", "Cambridge"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.hostname", "example.com"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.liveness_hosts.0", "clorigin3.www.example.com"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.origin_id", "alb_test_krk_dc1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.state_or_province", "MA"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "data_centers.0.cloud_server_host_header_override", "false"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.host_header", "clorigin3.www.example.com"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.additional_headers.Authorization", "test"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.interval", "25"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.path", "/status"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.peer_certificate_verification", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.port", "443"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.protocol", "HTTPS"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.request_string", "test"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.response_string", "test"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.status_3xx_failure", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.status_4xx_failure", "true"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.status_5xx_failure", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "liveness_settings.0.timeout", "30"),
			},
		},
		"specify load balancer version in file": {
			configPath: "testdata/TestDataCloudletsApplicationLoadBalancer/application_load_balancer_version.tf",
			init: func(m *cloudlets.Mock) {
				m.On("GetLoadBalancerVersion", testutils.MockContext, mock.Anything).Return(&cloudlets.LoadBalancerVersion{Version: 10}, nil)
				m.On("GetOrigin", testutils.MockContext, mock.Anything).Return(&cloudlets.Origin{}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_application_load_balancer.test", "version", "10"),
			},
		},
		"deleted load balancer version": {
			configPath: "testdata/TestDataCloudletsApplicationLoadBalancer/application_load_balancer_version.tf",
			init: func(m *cloudlets.Mock) {
				m.On("GetLoadBalancerVersion", testutils.MockContext, mock.Anything).Return(&cloudlets.LoadBalancerVersion{Version: 10, Deleted: true}, nil)
				m.On("GetOrigin", testutils.MockContext, mock.Anything).Return(&cloudlets.Origin{}, nil)
			},
			withError: regexp.MustCompile("specified load balancer version is deleted: 10"),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			client := &cloudlets.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       resource.ComposeAggregateTestCheckFunc(test.checkFunctions...),
							ExpectError: test.withError,
						},
					},
				})
				client.AssertExpectations(t)
			})
		})
	}
}

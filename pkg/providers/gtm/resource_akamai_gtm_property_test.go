package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var (
	propertyResourceName = "akamai_gtm_property.tfexample_prop_1"
	propertyName         = "tfexample_prop_1"
	updatedPropertyName  = "tfexample_prop_1-updated"
)

func TestResGTMProperty(t *testing.T) {
	tests := map[string]struct {
		property *gtm.Property
		init     func(*testing.T, *gtm.Mock)
		steps    []resource.TestStep
	}{
		"create property": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, getBasicProperty(), gtmTestDomain)
				// read
				mockGetProperty(m, getBasicProperty(), propertyName, gtmTestDomain, 4)
				// update
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 50, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 0, 20)
				mockUpdateProperty(m, getUpdatedProperty(), gtmTestDomain)
				// read
				mockGetDomainStatus(m, gtmTestDomain, 2)
				mockGetProperty(m, getUpdatedProperty(), propertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, getUpdatedProperty(), gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/update_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
					),
				},
			},
		},
		"create property with additional liveness test fields": {
			property: getBasicPropertyWithLivenessTests(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, getBasicPropertyWithLivenessTests(), gtmTestDomain)
				// read
				mockGetProperty(m, getBasicPropertyWithLivenessTests(), propertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, getBasicPropertyWithLivenessTests(), gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic_additional_liveness_tests.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", "GET"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", "Body"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "true"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.0", "test1"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
					),
				},
			},
		},
		"create property failed": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				// bad request status code returned
				m.On("CreateProperty",
					mock.Anything,
					getBasicProperty(),
					gtmTestDomain,
				).Return(nil, &gtm.Error{
					StatusCode: http.StatusBadRequest,
				})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectError: regexp.MustCompile("property Create failed"),
				},
			},
		},
		"create property denied": {
			property: nil,
			init: func(t *testing.T, m *gtm.Mock) {
				// create
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				// denied response status returned
				deniedResponse := gtm.PropertyResponse{
					Resource: getBasicProperty(),
					Status:   &deniedResponseStatus,
				}
				m.On("CreateProperty",
					mock.Anything,
					getBasicProperty(),
					gtmTestDomain,
				).Return(&deniedResponse, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
				},
			},
		},
		"create property and update name - force new": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				// create 1st property
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, getBasicProperty(), gtmTestDomain)
				// read
				mockGetProperty(m, getBasicProperty(), propertyName, gtmTestDomain, 4)
				// force new -> delete 1st property and recreate 2nd with updated name
				mockDeleteProperty(m, getBasicProperty(), gtmTestDomain)
				propertyWithUpdatedName := getBasicProperty()
				propertyWithUpdatedName.Name = updatedPropertyName
				mockNewProperty(m, updatedPropertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, propertyWithUpdatedName, gtmTestDomain)
				// read
				mockGetProperty(m, propertyWithUpdatedName, updatedPropertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, propertyWithUpdatedName, gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/update_name.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1-updated"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1-updated"),
					),
				},
			},
		},
		"test_object_protocol different than HTTP, HTTPS or FTP": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				// create property with test_object_protocol in first liveness test different from HTTP, HTTPS, FTP
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "SNMP", "", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				// alter mocked property
				propertyWithLivenessTest := getBasicProperty()
				propertyWithLivenessTest.LivenessTests[0].TestObject = ""
				propertyWithLivenessTest.LivenessTests[0].TestObjectProtocol = "SNMP"
				mockCreateProperty(m, propertyWithLivenessTest, gtmTestDomain)
				// read
				mockGetProperty(m, propertyWithLivenessTest, propertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, propertyWithLivenessTest, gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_not_required.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
					),
				},
			},
		},
		"create property with 'ranked-failover' type and two empty precedences in traffic target - error": {
			property: getRankedFailoverPropertyNoPrecedence(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/precedence/create_ranked_failover_empty_precedence.tf"),
					ExpectError: regexp.MustCompile(`Error: property cannot have multiple primary traffic targets \(targets with lowest precedence\)`),
				},
			},
		}, "create property with 'ranked-failover' type and no traffic targets - error": {
			property: getRankedFailoverPropertyNoPrecedence(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/precedence/create_ranked_failover_no_traffic_targets.tf"),
					ExpectError: regexp.MustCompile(`Error: at least one 'traffic_target' has to be defined and enabled`),
				},
			},
		},
		"create property with 'ranked-failover' type and allow single empty precedence value": {
			property: getRankedFailoverPropertyWithPrecedence(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence(), gtmTestDomain)
				// read
				mockGetProperty(m, getRankedFailoverPropertyWithPrecedence(), propertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, getRankedFailoverPropertyWithPrecedence(), gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/precedence/create_ranked_failover_precedence.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_method", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.http_request_body", ""),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.alternate_ca_certificates.#", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "liveness_test.0.pre_2023_security_posture", "false"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "10"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.1.precedence", "0"),
					),
				},
			},
		},
		"create property with 'ranked-failover' type and 0 set as precedence value": {
			property: getRankedFailoverPropertyWithPrecedence(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockNewProperty(m, propertyName)
				mockNewStaticRRSet(m)
				mockNewLivenessTest(m, "lt5", "HTTP", "/junk", 40, 1, 30)
				mockNewLivenessTest(m, "lt2", "HTTP", "/junk", 30, 80, 20)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence(), gtmTestDomain)
				// read
				mockGetProperty(m, getRankedFailoverPropertyWithPrecedence(), propertyName, gtmTestDomain, 3)
				// delete
				mockDeleteProperty(m, getRankedFailoverPropertyWithPrecedence(), gtmTestDomain)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/precedence/create_ranked_failover_0_precedence.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv4", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "weighted_hash_bits_for_ipv6", "0"),
						resource.TestCheckResourceAttr(propertyResourceName, "id", "gtm_terra_testdomain.akadns.net:tfexample_prop_1"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.0.precedence", "10"),
						resource.TestCheckResourceAttr(propertyResourceName, "traffic_target.1.precedence", "0"),
					),
				},
			},
		},
		"create property with test_object_protocol set to 'FTP' - test_object required error": {
			property: getBasicProperty(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_ftp.tf"),
					ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
				},
			},
		},
		"create property with test_object_protocol set to 'HTTP' - test_object required error": {
			property: getBasicProperty(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_http.tf"),
					ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
				},
			},
		},
		"create property with test_object_protocol set to 'HTTPS' - test_object required error": {
			property: getBasicProperty(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_https.tf"),
					ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := new(gtm.Mock)
			if test.init != nil {
				test.init(t, client)
			}
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

func TestResourceGTMTrafficTargetOrder(t *testing.T) {
	// To see actual plan when diff is expected, change 'nonEmptyPlan' to false in test case
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"second apply - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			pathForUpdate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic target with no datacenter_id - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id_diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"added traffic target - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/add_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"removed traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/remove_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in traffic target - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in re-ordered traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered servers in traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered servers and re-ordered traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered and changed servers in traffic target - diff in one traffic target": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/changed_and_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/change_server.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers and re-ordered traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/change_server_and_diff_traffic_target_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(test.client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.pathForCreate),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			test.client.AssertExpectations(t)
		})
	}
}

// getRankedFailoverPropertyWithPrecedence gets the property values taken from `create_ranked_failover_precedence.tf`
func getRankedFailoverPropertyWithPrecedence() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []*gtm.LivenessTest{
			{
				DisableNonstandardPortWarning: false,
				Name:                          "lt5",
				TestInterval:                  40,
				TestObject:                    "/junk",
				TestObjectPort:                1,
				TestObjectProtocol:            "HTTP",
				TestTimeout:                   30.0,
				HTTPHeaders: []*gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
			{
				Name:                        "lt2",
				TestInterval:                30,
				TestObjectProtocol:          "HTTP",
				TestTimeout:                 20,
				TestObject:                  "/junk",
				TestObjectPort:              80,
				PeerCertificateVerification: true,
				HTTPHeaders:                 []*gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []*gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight:     200.0,
				Precedence: ptr.To(10),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "ranked-failover",
	}
}

// getRankedFailoverPropertyNoPrecedence gets the property values taken from `create_ranked_failover_empty_precedence.tf`
func getRankedFailoverPropertyNoPrecedence() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:           300,
		HandoutMode:          "normal",
		HandoutLimit:         5,
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []*gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight: 200.0,
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight: 200.0,
			},
		},
		Type: "ranked-failover",
	}
}

// getUpdatedProperty gets the property with updated values taken from `update_basic.tf`
func getUpdatedProperty() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []*gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       50,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []*gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
			{
				Name:                        "lt2",
				TestInterval:                30,
				TestObjectProtocol:          "HTTP",
				TestTimeout:                 20,
				TestObject:                  "/junk",
				TestObjectPort:              80,
				PeerCertificateVerification: true,
				HTTPHeaders:                 []*gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []*gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.5",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicProperty gets the property values taken from `create_basic.tf`
func getBasicProperty() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []*gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []*gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
			{
				Name:                        "lt2",
				TestInterval:                30,
				TestObjectProtocol:          "HTTP",
				TestTimeout:                 20,
				TestObject:                  "/junk",
				TestObjectPort:              80,
				PeerCertificateVerification: true,
				HTTPHeaders:                 []*gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []*gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicPropertyWithLivenessTests gets the property values taken from `create_basic_additional_liveness_tests.tf`
func getBasicPropertyWithLivenessTests() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []*gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []*gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
				HTTPMethod:              ptr.To("GET"),
				HTTPRequestBody:         ptr.To("Body"),
				Pre2023SecurityPosture:  true,
				AlternateCACertificates: []string{"test1"},
			},
			{
				Name:                        "lt2",
				TestInterval:                30,
				TestObjectProtocol:          "HTTP",
				TestTimeout:                 20,
				TestObject:                  "/junk",
				TestObjectPort:              80,
				PeerCertificateVerification: true,
				HTTPHeaders:                 []*gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []*gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "weighted-round-robin",
	}
}

// getMocks is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func getMocks() *gtm.Mock {
	client := new(gtm.Mock)

	// read
	getPropertyCall := client.On("GetProperty", mock.Anything, "tfexample_prop_1", gtmTestDomain).
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound})
	// create
	mockNewProperty(client, "tfexample_prop_1")
	mockNewLivenessTest(client, "lt5", "HTTP", "/junk", 40, 80, 30.0)
	// mock.AnythingOfType *gtm.Property is used is those mock calls as there are too many different test cases to mock
	// each one and for those test it's not important, since we are only checking the diff
	client.On("CreateProperty", mock.Anything, mock.AnythingOfType("*gtm.Property"), mock.AnythingOfType("string")).Return(&gtm.PropertyResponse{
		Resource: getBasicProperty(),
		Status:   &pendingResponseStatus,
	}, nil).Run(func(args mock.Arguments) {
		getPropertyCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
	})
	// delete
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("*gtm.Property"),
		"gtm_terra_testdomain.akadns.net",
	).Return(&completeResponseStatus, nil)

	return client
}

func mockNewProperty(client *gtm.Mock, propertyName string) {
	client.On("NewProperty",
		mock.Anything,
		propertyName,
	).Return(&gtm.Property{
		Name: propertyName,
	}).Once()
}

func mockNewTrafficTarget(client *gtm.Mock, times int) {
	client.On("NewTrafficTarget",
		mock.Anything,
	).Return(&gtm.TrafficTarget{}).Times(times)
}

func mockNewStaticRRSet(client *gtm.Mock) {
	client.On("NewStaticRRSet",
		mock.Anything,
	).Return(&gtm.StaticRRSet{}).Once()
}

func mockNewLivenessTest(client *gtm.Mock, name, protocol, object string, interval, port int, timeout float32) {
	client.On("NewLivenessTest",
		mock.Anything,
		name,
		protocol,
		interval,
		timeout,
	).Return(&gtm.LivenessTest{
		Name:               name,
		TestObjectProtocol: protocol,
		TestInterval:       interval,
		TestTimeout:        timeout,
		TestObjectPort:     port,
		TestObject:         object,
	}).Once()
}

func mockCreateProperty(client *gtm.Mock, property *gtm.Property, domain string) {
	resp := gtm.PropertyResponse{}
	resp.Resource = property
	resp.Status = &pendingResponseStatus
	client.On("CreateProperty",
		mock.Anything,
		property,
		domain,
	).Return(&resp, nil).Once()
}

func mockGetProperty(client *gtm.Mock, property *gtm.Property, propertyName, domain string, times int) {
	client.On("GetProperty",
		mock.Anything,
		propertyName,
		domain,
	).Return(property, nil).Times(times)
}

func mockUpdateProperty(client *gtm.Mock, property *gtm.Property, domain string) {
	client.On("UpdateProperty",
		mock.Anything,
		property,
		domain,
	).Return(&completeResponseStatus, nil).Once()
}

func mockGetDomainStatus(client *gtm.Mock, domain string, times int) {
	client.On("GetDomainStatus",
		mock.Anything,
		domain,
	).Return(&completeResponseStatus, nil).Times(times)
}

func mockDeleteProperty(client *gtm.Mock, property *gtm.Property, domain string) {
	client.On("DeleteProperty",
		mock.Anything,
		property,
		domain,
	).Return(&completeResponseStatus, nil).Once()
}

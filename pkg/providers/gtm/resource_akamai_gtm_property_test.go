package gtm

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var (
	propertyResourceName = "akamai_gtm_property.tfexample_prop_1"
	updatedPropertyName  = "tfexample_prop_1-updated"

	updatePropertyResponseStatus = &gtm.UpdatePropertyResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}

	deletePropertyResponseStatus = &gtm.DeletePropertyResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getBasicProperty())
				// read
				mockGetProperty(m, getBasicPropertyResponse(), nil, 4)
				// update
				mockUpdateProperty(m)
				// read
				mockGetDomainStatus(m, 2)
				mockGetProperty(m, getBasicPropertyResponseUpdate(), nil, 3)
				// delete
				mockDeleteProperty(m)
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getBasicPropertyWithLivenessTests())
				// read
				mockGetProperty(m, getBasicPropertyResponseWithLivenessTests(), nil, 3)
				// delete
				mockDeleteProperty(m)
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// bad request status code returned
				m.On("CreateProperty",
					mock.Anything,
					mock.AnythingOfType("gtm.CreatePropertyRequest"),
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
		"create property failed - property already exists": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, getBasicPropertyResponse(), nil, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectError: regexp.MustCompile("property already exists:"),
				},
			},
		},
		"create property with retry on Property Validation Failure - no datacenter is assigned to map target": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				m.On("CreateProperty",
					mock.Anything,
					gtm.CreatePropertyRequest{Property: getBasicProperty(), DomainName: gtmTestDomain},
				).Return(nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no datacenter is assigned to map target (all others)",
				}).Once()

				// Simulate successful property creation on the second attempt
				mockCreateProperty(m, getBasicProperty())
				mockGetProperty(m, getBasicPropertyResponse(), nil, 3)
				mockDeleteProperty(m)
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
			},
		},
		"create property with retry on Property Validation Failure - other errors": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				m.On("CreateProperty",
					mock.Anything,
					gtm.CreatePropertyRequest{Property: getBasicProperty(), DomainName: gtmTestDomain},
				).Return(nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no targets found",
				}).Once()

			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectError: regexp.MustCompile("CreateProperty error: property Create failed: error: API error"),
				},
			},
		},
		"create property with retry - context canceled": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				m.On("CreateProperty",
					mock.Anything,
					gtm.CreatePropertyRequest{Property: getBasicProperty(), DomainName: gtmTestDomain},
				).Return(nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no datacenter is assigned to map target (all others)",
				}).Once()

				// Simulate context cancellation on the second attempt
				m.On("CreateProperty",
					mock.Anything,
					gtm.CreatePropertyRequest{Property: getBasicProperty(), DomainName: gtmTestDomain},
				).Return(nil, context.Canceled).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectError: regexp.MustCompile("CreateProperty error: property Create failed: error: context canceled"),
				},
			},
		},
		"create property denied": {
			property: nil,
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create
				// denied response status returned
				deniedResponse := gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &deniedResponseStatus,
				}
				m.On("CreateProperty",
					mock.Anything,
					mock.AnythingOfType("gtm.CreatePropertyRequest"),
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create 1st property
				mockCreateProperty(m, getBasicProperty())
				// read
				mockGetProperty(m, getBasicPropertyResponse(), nil, 4)
				// force new -> delete 1st property and recreate 2nd with updated name
				mockDeleteProperty(m)
				propertyWithUpdatedName := getBasicProperty()
				propertyWithUpdatedName.Name = updatedPropertyName
				propertyResponseWithUpdatedName := getBasicPropertyResponse()
				propertyResponseWithUpdatedName.Name = updatedPropertyName
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, propertyWithUpdatedName)
				// read
				mockGetProperty(m, propertyResponseWithUpdatedName, nil, 3)
				// delete
				mockDeleteProperty(m)
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create property with test_object_protocol in first liveness test different from HTTP, HTTPS, FTP
				// alter mocked property
				propertyWithLivenessTest := getBasicProperty()
				propertyWithLivenessTest.LivenessTests[0].TestObject = ""
				propertyWithLivenessTest.LivenessTests[0].TestObjectProtocol = "SNMP"
				propertyResponseWithLivenessTest := getBasicPropertyResponse()
				propertyResponseWithLivenessTest.LivenessTests[0].TestObject = ""
				propertyResponseWithLivenessTest.LivenessTests[0].TestObjectProtocol = "SNMP"
				mockCreateProperty(m, propertyWithLivenessTest)
				// read
				mockGetProperty(m, propertyResponseWithLivenessTest, nil, 3)
				// delete
				mockDeleteProperty(m)
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence())
				// read
				mockGetProperty(m, getRankedFailoverPropertyResponseWithPrecedence(), nil, 3)
				// delete
				mockDeleteProperty(m)
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
				mockGetProperty(m, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence())
				// read
				mockGetProperty(m, getRankedFailoverPropertyResponseWithPrecedence(), nil, 3)
				// delete
				mockDeleteProperty(m)
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
			client:        getMocksSecondApply(),
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
			client:        getMocksWithoutDatacenterID(),
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
		LivenessTests: []gtm.LivenessTest{
			{
				DisableNonstandardPortWarning: false,
				Name:                          "lt5",
				TestInterval:                  40,
				TestObject:                    "/junk",
				TestObjectPort:                1,
				TestObjectProtocol:            "HTTP",
				TestTimeout:                   30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getRankedFailoverPropertyResponseWithPrecedence gets the property values taken from `create_ranked_failover_precedence.tf`
func getRankedFailoverPropertyResponseWithPrecedence() *gtm.GetPropertyResponse {
	return &gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				DisableNonstandardPortWarning: false,
				Name:                          "lt5",
				TestInterval:                  40,
				TestObject:                    "/junk",
				TestObjectPort:                1,
				TestObjectProtocol:            "HTTP",
				TestTimeout:                   30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

func TestResourceGTMLivenessTestOrder(t *testing.T) {
	// To see actual plan when diff is expected, change 'nonEmptyPlan' to false in test case
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"second apply - no diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered liveness test - no diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_liveness_tests_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"remove liveness test - diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/remove_liveness_test.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"add liveness test - diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/add_liveness_tests.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and re-ordered http headers - no diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_lt_and_header_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change of 'timeout' field - diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_timeout.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change of 'timeout' field and reorder of liveness tests - diff_(messy)": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_timeout_reorder_lt.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and change http headers - diff_(messy)": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_lt_order_and_header_change.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change http headers - diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_header.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"value added to http header - diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/http_header_without_value.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and alternate ca certificates - no diff": {
			client:        getMocksForLivenessTest(),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests_with_ca_cert.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_lt_and_ca_certificate_order.tf",
			nonEmptyPlan:  false,
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
							Check: resource.ComposeTestCheckFunc(
								resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
								resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
							),
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

// getUpdatedProperty gets the property with updated values taken from `update_basic.tf`
func getUpdatedProperty() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       50,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getBasicPropertyResponseDiff gets the property values taken from `create_basic.tf`
func getBasicPropertyResponseDiff() gtm.GetPropertyResponse {
	return gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				AlternateCACertificates: []string{},
				Name:                    "lt5",
				TestInterval:            40,
				TestObject:              "/junk",
				TestObjectPort:          1,
				TestObjectProtocol:      "HTTP",
				TestTimeout:             30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3133,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.7",
					"1.2.3.8",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicPropertyCopy gets the property values taken from `create_basic.tf`
func getBasicPropertyDiff() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				AlternateCACertificates: []string{},
				Name:                    "lt5",
				TestInterval:            40,
				TestObject:              "/junk",
				TestObjectPort:          1,
				TestObjectProtocol:      "HTTP",
				TestTimeout:             30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3133,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.7",
					"1.2.3.8",
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
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getBasicPropertySecondApply gets the property values taken from `create_multiple_traffic_targets.tf`
func getBasicPropertySecondApply() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
				},
				Weight: 200,
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.5",
				},
				Weight: 200,
			},
			{
				DatacenterID: 3133,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight: 200,
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicPropertyResponseSecondApply gets the property values taken from `create_multiple_traffic_targets.tf`
func getBasicPropertyResponseSecondApply() gtm.GetPropertyResponse {
	return gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
				},
				Weight: 200,
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.5",
				},
				Weight: 200,
			},
			{
				DatacenterID: 3133,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight: 200,
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicPropertyResponse gets the property values taken from `create_basic.tf`
func getBasicPropertyResponse() *gtm.GetPropertyResponse {
	return &gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getBasicPropertyResponseUpdate gets the property values taken from `update_basic.tf`
func getBasicPropertyResponseUpdate() *gtm.GetPropertyResponse {
	return &gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       50,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getBasicPropertyWithoutDatacenterID gets the property values without DatacenterID
func getBasicPropertyWithoutDatacenterID() *gtm.Property {
	return &gtm.Property{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				AlternateCACertificates: []string{},
				Name:                    "lt5",
				TestInterval:            40,
				TestObject:              "/junk",
				TestObjectPort:          1,
				TestObjectProtocol:      "HTTP",
				TestTimeout:             30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.7",
					"1.2.3.8",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
		},
		Type: "weighted-round-robin",
	}
}

// getBasicPropertyResponseWithoutDatacenterID gets the property values without DatacenterID
func getBasicPropertyResponseWithoutDatacenterID() gtm.GetPropertyResponse {
	return gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				AlternateCACertificates: []string{},
				Name:                    "lt5",
				TestInterval:            40,
				TestObject:              "/junk",
				TestObjectPort:          1,
				TestObjectProtocol:      "HTTP",
				TestTimeout:             30.0,
				HTTPHeaders: []gtm.HTTPHeader{
					{
						Name:  "test_name",
						Value: "test_value",
					},
				},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight:     200.0,
				Precedence: ptr.To(0),
			},
			{
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.7",
					"1.2.3.8",
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
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getBasicPropertyResponseWithLivenessTests gets the property values taken from `create_basic_additional_liveness_tests.tf`
func getBasicPropertyResponseWithLivenessTests() *gtm.GetPropertyResponse {
	return &gtm.GetPropertyResponse{
		DynamicTTL:   300,
		HandoutMode:  "normal",
		HandoutLimit: 5,
		LivenessTests: []gtm.LivenessTest{
			{
				Name:               "lt5",
				TestInterval:       40,
				TestObject:         "/junk",
				TestObjectPort:     1,
				TestObjectProtocol: "HTTP",
				TestTimeout:        30.0,
				HTTPHeaders: []gtm.HTTPHeader{
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
				HTTPHeaders:                 []gtm.HTTPHeader{},
			},
		},
		Name:                 "tfexample_prop_1",
		ScoreAggregationType: "median",
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		TrafficTargets: []gtm.TrafficTarget{
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

// getMocksForLivenessTest is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func getMocksForLivenessTest() *gtm.Mock {
	client := new(gtm.Mock)
	// read
	getPropertyCall := client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound}).Times(5)
	// create
	// mock.AnythingOfType *gtm.Property is used is those mock calls as there are too many different test cases to mock
	// each one and for those test it's not important, since we are only checking the diff
	mockCreateProperty(client, getBasicProperty()).
		Run(func(args mock.Arguments) {
			arg := (*gtm.GetPropertyResponse)(args.Get(1).(gtm.CreatePropertyRequest).Property)
			getPropertyCall.ReturnArguments = mock.Arguments{arg, nil}
		})

	// delete
	mockDeleteProperty(client)

	return client
}

// getMocks is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func getMocks() *gtm.Mock {
	client := new(gtm.Mock)

	// read
	getCall := client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound}).Twice()
	// create
	// mock.AnythingOfType *gtm.Property is used is those mock calls as there are too many different test cases to mock
	// each one and for those test it's not important, since we are only checking the diff
	resp := getBasicPropertyResponseDiff()
	mockCreateProperty(client, getBasicPropertyDiff()).Run(func(args mock.Arguments) {
		getCall.ReturnArguments = mock.Arguments{&resp, nil}
	})

	// read
	client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(&resp, nil).Times(3)

	// delete
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.DeletePropertyRequest"),
	).Return(deletePropertyResponseStatus, nil)

	return client
}

// getMocksNoDatacenterId is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func getMocksWithoutDatacenterID() *gtm.Mock {
	client := new(gtm.Mock)

	// read
	getCall := client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound}).Twice()
	// create
	// mock.AnythingOfType *gtm.Property is used is those mock calls as there are too many different test cases to mock
	// each one and for those test it's not important, since we are only checking the diff
	resp := getBasicPropertyResponseWithoutDatacenterID()
	mockCreateProperty(client, getBasicPropertyWithoutDatacenterID()).Run(func(args mock.Arguments) {
		getCall.ReturnArguments = mock.Arguments{&resp, nil}
	})

	// read
	client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(&resp, nil).Times(3)

	// delete
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.DeletePropertyRequest"),
	).Return(deletePropertyResponseStatus, nil)

	return client
}

// getMocksSecondApply is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func getMocksSecondApply() *gtm.Mock {
	client := new(gtm.Mock)

	// read
	getCall := client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound}).Twice()
	// create
	// mock.AnythingOfType *gtm.Property is used is those mock calls as there are too many different test cases to mock
	// each one and for those test it's not important, since we are only checking the diff
	resp := getBasicPropertyResponseSecondApply()
	mockCreateProperty(client, getBasicPropertySecondApply()).Run(func(args mock.Arguments) {
		getCall.ReturnArguments = mock.Arguments{&resp, nil}
	})

	// read
	client.On("GetProperty", mock.Anything, mock.AnythingOfType("gtm.GetPropertyRequest")).
		Return(&resp, nil).Times(3)

	// delete
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.DeletePropertyRequest"),
	).Return(deletePropertyResponseStatus, nil)

	return client
}

func mockCreateProperty(client *gtm.Mock, property *gtm.Property) *mock.Call {
	resp := gtm.CreatePropertyResponse{}
	resp.Resource = property
	resp.Status = &pendingResponseStatus
	return client.On("CreateProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.CreatePropertyRequest"),
	).Return(&resp, nil).Once()
}

func mockGetProperty(client *gtm.Mock, property *gtm.GetPropertyResponse, error *gtm.Error, times int) {
	if property != nil {
		client.On("GetProperty",
			mock.Anything,
			mock.AnythingOfType("gtm.GetPropertyRequest"),
		).Return(property, nil).Times(times)
	} else {
		client.On("GetProperty",
			mock.Anything,
			mock.AnythingOfType("gtm.GetPropertyRequest"),
		).Return(nil, error).Times(times)
	}
}

func mockUpdateProperty(client *gtm.Mock) {
	client.On("UpdateProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.UpdatePropertyRequest"),
	).Return(updatePropertyResponseStatus, nil).Once()
}

func mockGetDomainStatus(client *gtm.Mock, times int) {
	client.On("GetDomainStatus",
		mock.Anything,
		mock.AnythingOfType("gtm.GetDomainStatusRequest"),
	).Return(getDomainStatusResponseStatus, nil).Times(times)
}

func mockDeleteProperty(client *gtm.Mock) {
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("gtm.DeletePropertyRequest"),
	).Return(deletePropertyResponseStatus, nil).Once()
}

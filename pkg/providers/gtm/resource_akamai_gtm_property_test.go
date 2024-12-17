package gtm

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const (
	testPropertyName     = "tfexample_prop_1"
	propertyResourceName = "akamai_gtm_property.tfexample_prop_1"
	updatedPropertyName  = "tfexample_prop_1-updated"
)

var (
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getBasicProperty(), &gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 4)
				// update
				var updateProperty = *getBasicProperty()
				updateProperty.TrafficTargets[0].DatacenterID = 3132
				updateProperty.TrafficTargets[0].Servers = []string{"1.2.3.5"}
				updateProperty.LivenessTests[0].TestInterval = 50
				mockUpdateProperty(m, &updateProperty)
				// read
				mockGetDomainStatus(m, 2)
				mockGetProperty(m, testPropertyName, getBasicPropertyResponseUpdate(), nil, 3)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
		"create property, remove outside of terraform, expect non-empty plan": {
			property: getBasicProperty(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getBasicProperty(), &gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 2)

				// Mock that the property was deleted outside terraform
				mockGetProperty(m, testPropertyName, nil, gtm.ErrNotFound, 1)

				// For terraform test framework, we need to mock GetProperty as it would actually exist before deletion
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 1)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
					Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		},
		"create property with additional liveness test fields": {
			property: getBasicPropertyWithLivenessTests(),
			init: func(t *testing.T, m *gtm.Mock) {
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getBasicPropertyWithLivenessTests(), &gtm.CreatePropertyResponse{
					Resource: getBasicPropertyWithLivenessTests(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getBasicPropertyResponseWithLivenessTests(), nil, 3)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// bad request status code returned
				mockCreateProperty(m, getBasicProperty(), nil, &gtm.Error{
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
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 1)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				mockCreateProperty(m, getBasicProperty(), nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no datacenter is assigned to map target (all others)",
				})

				// Simulate successful property creation on the second attempt
				mockCreateProperty(m, getBasicProperty(), &gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &pendingResponseStatus,
				}, nil)
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 3)
				mockDeleteProperty(m, testPropertyName)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				mockCreateProperty(m, getBasicProperty(), nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no targets found",
				})
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// Simulate a retry scenario
				mockCreateProperty(m, getBasicProperty(), nil, &gtm.Error{
					Type:       "https://problems.luna.akamaiapis.net/config-gtm/v1/propertyValidationError",
					StatusCode: http.StatusBadRequest,
					Title:      "Property Validation Failure",
					Detail:     "Invalid configuration for property \"tfexample_prop_1\": no datacenter is assigned to map target (all others)",
				})

				// Simulate context cancellation on the second attempt
				mockCreateProperty(m, getBasicProperty(), nil, context.Canceled)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create
				// denied response status returned
				deniedResponse := gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &deniedResponseStatus,
				}
				mockCreateProperty(m, getBasicProperty(), &deniedResponse, nil)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create 1st property
				mockCreateProperty(m, getBasicProperty(), &gtm.CreatePropertyResponse{
					Resource: getBasicProperty(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getBasicPropertyResponse(), nil, 4)
				// force new -> delete 1st property and recreate 2nd with updated name
				mockDeleteProperty(m, testPropertyName)
				propertyWithUpdatedName := getBasicProperty()
				propertyWithUpdatedName.Name = updatedPropertyName
				propertyResponseWithUpdatedName := getBasicPropertyResponse()
				propertyResponseWithUpdatedName.Name = updatedPropertyName
				mockGetProperty(m, updatedPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, propertyWithUpdatedName, &gtm.CreatePropertyResponse{
					Resource: propertyWithUpdatedName,
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, updatedPropertyName, propertyResponseWithUpdatedName, nil, 3)
				// delete
				mockDeleteProperty(m, updatedPropertyName)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				// create property with test_object_protocol in first liveness test different from HTTP, HTTPS, FTP
				// alter mocked property
				propertyWithLivenessTest := getBasicProperty()
				propertyWithLivenessTest.LivenessTests[0].TestObject = ""
				propertyWithLivenessTest.LivenessTests[0].TestObjectProtocol = "SNMP"
				propertyResponseWithLivenessTest := getBasicPropertyResponse()
				propertyResponseWithLivenessTest.LivenessTests[0].TestObject = ""
				propertyResponseWithLivenessTest.LivenessTests[0].TestObjectProtocol = "SNMP"
				mockCreateProperty(m, propertyWithLivenessTest, &gtm.CreatePropertyResponse{
					Resource: propertyWithLivenessTest,
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, propertyResponseWithLivenessTest, nil, 3)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence(), &gtm.CreatePropertyResponse{
					Resource: getRankedFailoverPropertyWithPrecedence(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getRankedFailoverPropertyResponseWithPrecedence(), nil, 3)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
				mockGetProperty(m, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
				mockCreateProperty(m, getRankedFailoverPropertyWithPrecedence(), &gtm.CreatePropertyResponse{
					Resource: getRankedFailoverPropertyWithPrecedence(),
					Status:   &pendingResponseStatus,
				}, nil)
				// read
				mockGetProperty(m, testPropertyName, getRankedFailoverPropertyResponseWithPrecedence(), nil, 3)
				// delete
				mockDeleteProperty(m, testPropertyName)
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
			client:        mockPropertyFlow(getBasicPropertySecondApply()),
			pathForCreate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			pathForUpdate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic targets - no diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic target with no datacenter_id - no diff": {
			client:        mockPropertyFlow(getBasicPropertyWithoutDatacenterID()),
			pathForCreate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id_diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"added traffic target - diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/add_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"removed traffic target - diff (messy)": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/remove_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in traffic target - diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in re-ordered traffic target - diff (messy)": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered servers in traffic targets - no diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered servers and re-ordered traffic targets - no diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered and changed servers in traffic target - diff in one traffic target": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/changed_and_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers - diff": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/change_server.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers and re-ordered traffic target - diff (messy)": {
			client:        mockPropertyFlow(getBasicPropertyForTrafficTargetOrder()),
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
		Name:                 testPropertyName,
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
		Name:                 testPropertyName,
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
		Name:                 testPropertyName,
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
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered liveness test - no diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_liveness_tests_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"remove liveness test - diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/remove_liveness_test.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"add liveness test - diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/add_liveness_tests.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and re-ordered http headers - no diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_lt_and_header_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change of 'timeout' field - diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_timeout.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change of 'timeout' field and reorder of liveness tests - diff_(messy)": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_timeout_reorder_lt.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and change http headers - diff_(messy)": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/diff_lt_order_and_header_change.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change http headers - diff": {
			client:        mockPropertyFlow(getLivenessTestDefaultProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/change_header.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"value added to http header - diff": {
			client:        mockPropertyFlow(getLivenessTestValueAddedToHTTPHeaderProperty()),
			pathForCreate: "testdata/TestResGtmProperty/liveness_test/http_header_without_value.tf",
			pathForUpdate: "testdata/TestResGtmProperty/liveness_test/multiple_liveness_tests.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered liveness test and alternate ca certificates - no diff": {
			client:        mockPropertyFlow(getLivenessTestCaCertificatesProperty()),
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

func TestResGTMPropertyImport(t *testing.T) {
	tests := map[string]struct {
		domainName  string
		mapName     string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: "gtm_terra_testdomain.akadns.net",
			mapName:    "tfexample_prop_1",
			init: func(m *gtm.Mock) {
				// Read
				importedProperty := gtm.GetPropertyResponse(*getImportedProperty())
				mockGetPropertyImport(m, &importedProperty, nil).Times(2)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "gtm_terra_testdomain.akadns.net").
				CheckEqual("name", "tfexample_prop_1").
				CheckEqual("ipv6", "true").
				CheckEqual("score_aggregation_type", "median").
				CheckEqual("stickiness_bonus_percentage", "10").
				CheckEqual("stickiness_bonus_constant", "10").
				CheckEqual("health_threshold", "123").
				CheckEqual("use_computed_targets", "true").
				CheckEqual("backup_ip", "test ip").
				CheckEqual("balance_by_download_score", "true").
				CheckEqual("unreachable_threshold", "1234").
				CheckEqual("min_live_fraction", "1").
				CheckEqual("health_multiplier", "5").
				CheckEqual("dynamic_ttl", "300").
				CheckEqual("max_unreachable_penalty", "123").
				CheckEqual("map_name", "test map").
				CheckEqual("handout_limit", "5").
				CheckEqual("handout_mode", "normal").
				CheckEqual("load_imbalance_percentage", "10").
				CheckEqual("failover_delay", "5").
				CheckEqual("backup_cname", "test cname").
				CheckEqual("failback_delay", "5").
				CheckEqual("health_max", "123").
				CheckEqual("ghost_demand_reporting", "false").
				CheckEqual("weighted_hash_bits_for_ipv4", "4").
				CheckEqual("weighted_hash_bits_for_ipv6", "6").
				CheckEqual("cname", "test cName").
				CheckEqual("comments", "test comment").
				CheckEqual("type", "failover").
				CheckEqual("traffic_target.0.datacenter_id", "3131").
				CheckEqual("traffic_target.0.enabled", "true").
				CheckEqual("traffic_target.0.weight", "200").
				CheckEqual("traffic_target.0.servers.0", "1.2.3.9").
				CheckEqual("traffic_target.0.handout_cname", "test").
				CheckEqual("traffic_target.0.precedence", "10").
				CheckEqual("static_rr_set.0.type", "MX").
				CheckEqual("static_rr_set.0.ttl", "300").
				CheckEqual("static_rr_set.0.rdata.0", "100 test_e").
				CheckEqual("liveness_test.0.disable_nonstandard_port_warning", "false").
				CheckEqual("liveness_test.0.name", "lt5").
				CheckEqual("liveness_test.0.test_interval", "40").
				CheckEqual("liveness_test.0.test_object", "/junk").
				CheckEqual("liveness_test.0.test_object_port", "1").
				CheckEqual("liveness_test.0.test_object_protocol", "HTTP").
				CheckEqual("liveness_test.0.test_timeout", "30").
				CheckEqual("liveness_test.0.http_header.0.name", "test_name").
				CheckEqual("liveness_test.0.http_header.0.value", "test_value").
				CheckEqual("wait_on_complete", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:  "",
			mapName:     "tfexample_prop_1",
			expectError: regexp.MustCompile(`Error: invalid resource ID: :tfexample_prop_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:  "gtm_terra_testdomain.akadns.net",
			mapName:     "",
			expectError: regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName: "gtm_terra_testdomain.akadns.net",
			mapName:    "tfexample_prop_1",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetPropertyImport(m, nil, fmt.Errorf("get failed")).Once()
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    fmt.Sprintf("%s:%s", tc.domainName, tc.mapName),
							ImportState:      true,
							ResourceName:     "akamai_gtm_property.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockGetPropertyImport(m *gtm.Mock, resp *gtm.GetPropertyResponse, err error) *mock.Call {
	return m.On("GetProperty", mock.Anything, gtm.GetPropertyRequest{
		PropertyName: "tfexample_prop_1",
		DomainName:   "gtm_terra_testdomain.akadns.net",
	}).Return(resp, err)
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
		Name:                 testPropertyName,
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
func getBasicPropertySecondApply() gtm.Property {
	return gtm.Property{
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
		Name:                 testPropertyName,
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.4",
				},
				Weight:     200,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3132,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.5",
				},
				Weight:     200,
				Precedence: ptr.To(0),
			},
			{
				DatacenterID: 3133,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.6",
				},
				Weight:     200,
				Precedence: ptr.To(0),
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
		Name:                 testPropertyName,
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
		Name:                 testPropertyName,
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
func getBasicPropertyWithoutDatacenterID() gtm.Property {
	return gtm.Property{
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
		Name:                 testPropertyName,
		ScoreAggregationType: "median",
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.5",
					"1.2.3.4",
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
		Name:                 testPropertyName,
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
		Name:                 testPropertyName,
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

func getLivenessTestDefaultProperty() gtm.Property {
	prp := getBasicPropertyForTrafficTargetOrder()
	prp.LivenessTests = append(prp.LivenessTests, gtm.LivenessTest{
		Name:               "OJ",
		TestInterval:       60,
		TestObjectProtocol: "HTTP",
		TestTimeout:        10,
		TestObject:         "/",
		TestObjectPort:     80,
		HTTPHeaders:        []gtm.HTTPHeader{},
		HTTPError3xx:       true,
		HTTPError4xx:       true,
		HTTPError5xx:       true,
	})
	return prp
}

func getImportedProperty() *gtm.Property {
	return &gtm.Property{
		Name:                      "tfexample_prop_1",
		Type:                      "failover",
		IPv6:                      true,
		ScoreAggregationType:      "median",
		StickinessBonusPercentage: 10.0,
		StickinessBonusConstant:   10,
		HealthThreshold:           123.0,
		UseComputedTargets:        true,
		BackupIP:                  "test ip",
		BalanceByDownloadScore:    true,
		StaticRRSets: []gtm.StaticRRSet{
			{
				Type:  "MX",
				TTL:   300,
				Rdata: []string{"100 test_e"},
			},
		},
		UnreachableThreshold:    1234.0,
		MinLiveFraction:         1.0,
		HealthMultiplier:        5.0,
		DynamicTTL:              300,
		MaxUnreachablePenalty:   123,
		MapName:                 "test map",
		HandoutLimit:            5,
		HandoutMode:             "normal",
		FailoverDelay:           5,
		BackupCName:             "test cname",
		FailbackDelay:           5,
		LoadImbalancePercentage: 10.0,
		HealthMax:               123.0,
		GhostDemandReporting:    false,
		Comments:                "test comment",
		CName:                   "test cName",
		WeightedHashBitsForIPv4: 4,
		WeightedHashBitsForIPv6: 6,
		TrafficTargets: []gtm.TrafficTarget{
			{
				DatacenterID: 3131,
				Enabled:      true,
				HandoutCName: "test",
				Servers: []string{
					"1.2.3.9",
				},
				Weight:     200,
				Precedence: ptr.To(10),
			},
		},
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
		},
	}
}

func getLivenessTestValueAddedToHTTPHeaderProperty() gtm.Property {
	property := getLivenessTestDefaultProperty()
	property.LivenessTests[0].HTTPHeaders[0].Value = ""
	return property
}

func getLivenessTestCaCertificatesProperty() gtm.Property {
	property := getLivenessTestDefaultProperty()
	property.LivenessTests[0].AlternateCACertificates = []string{"test1", "test2"}
	return property
}

// mockPropertyFlow is used for diff tests, where the contents of property not matter as much, as those tests aim to check the diffs
func mockPropertyFlow(gtmProperty gtm.Property) *gtm.Mock {
	client := new(gtm.Mock)
	// read
	mockGetProperty(client, testPropertyName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, 1)
	// create
	mockCreateProperty(client, &gtmProperty, &gtm.CreatePropertyResponse{
		Resource: &gtmProperty,
		Status:   &pendingResponseStatus,
	}, nil)

	var getPropertyResponse = gtm.GetPropertyResponse(gtmProperty)
	mockGetProperty(client, testPropertyName, &getPropertyResponse, nil, 4)

	// delete
	mockDeleteProperty(client, testPropertyName)

	return client
}

func getBasicPropertyForTrafficTargetOrder() gtm.Property {
	prp := getBasicPropertyWithoutDatacenterID()
	prp.TrafficTargets[2].DatacenterID = 3133
	return prp
}

func mockCreateProperty(client *gtm.Mock, property *gtm.Property, resp *gtm.CreatePropertyResponse, err error) *mock.Call {
	return client.On("CreateProperty",
		mock.Anything,
		gtm.CreatePropertyRequest{
			Property:   property,
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockGetProperty(client *gtm.Mock, propertyName string, property *gtm.GetPropertyResponse, error error, times int) {
	client.On("GetProperty",
		mock.Anything,
		gtm.GetPropertyRequest{
			DomainName:   testDomainName,
			PropertyName: propertyName,
		},
	).Return(property, error).Times(times)
}

func mockUpdateProperty(client *gtm.Mock, updatedProperty *gtm.Property) {
	client.On("UpdateProperty",
		mock.Anything,
		gtm.UpdatePropertyRequest{
			Property:   updatedProperty,
			DomainName: testDomainName,
		},
	).Return(updatePropertyResponseStatus, nil).Once()
}

func mockGetDomainStatus(client *gtm.Mock, times int) {
	client.On("GetDomainStatus",
		mock.Anything,
		gtm.GetDomainStatusRequest{
			DomainName: testDomainName,
		},
	).Return(getDomainStatusResponseStatus, nil).Times(times)
}

func mockDeleteProperty(client *gtm.Mock, propertyName string) {
	client.On("DeleteProperty",
		mock.Anything,
		gtm.DeletePropertyRequest{
			DomainName:   testDomainName,
			PropertyName: propertyName,
		},
	).Return(deletePropertyResponseStatus, nil).Once()
}

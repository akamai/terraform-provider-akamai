package property

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccAkamaiPropertyActivationConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

resource "akamai_property_activation" "property_activation" {
	property_id = "${akamai_property.property.id}"
	version = "${akamai_property.property.version}"
	network = "STAGING"
	activate = true
	contact = ["dshafik@akamai.com"]
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_cp_code" "cp_code" {
	name = "terraform-testing3"
	contract = "${data.akamai_contract.contract.id}"
	group = "${data.akamai_group.group.id}"
	product = "prd_SPM"
}

resource "akamai_edge_hostname" "test" {
    product = "prd_SPM"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
    edge_hostname =  "terraform-test3.exampleterraform.io.edgesuite.net"
    ipv6 = true
}

resource "akamai_property" "property" {
  name = "terraform-test3"
  id = "prp_1234"

  contact = ["user@exampleterraform.io"]

  product = "prd_SPM"
  cp_code = "${akamai_cp_code.cp_code.id}"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"

  hostnames = {
		"example.org" = "${akamai_edge_hostname.test.edge_hostname}"
  }
  
  rule_format = "v2016-11-15"
  
  rules = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
 	rules {
		behavior {
			name =  "origin"
        	option { 
       			key =  "cacheKeyHostname"
            	value = "ORIGIN_HOSTNAME"
        	}
			option { 
    			key =  "compress"
     			value = "true"
     		}
    		option { 
    			key =  "enableTrueClientIp"
     			value = "false"
     		}
    		option { 
    			key =  "forwardHostHeader"
     			value = "REQUEST_HOST_HEADER"
     		}
    		option { 
    			key =  "hostname"
     			value = "exampleterraform.io"
     		}
    		option { 
    			key =  "httpPort"
     			value = "80"
     		}
    		option { 
    			key =  "httpsPort"
     			value = "443"
     		}
    		option { 
    			key =  "originSni"
     			value = "true"
     		}
    		option { 
    			key =  "originType"
     			value = "CUSTOMER"
     		}
    		option { 
    			key =  "verificationMode"
     			value = "PLATFORM_SETTINGS"
     		}
    		option { 
    			key =  "originCertificate"
     			value = ""
     		}
    		option { 
    			key =  "ports"
     			value = ""
     		}
      	}
		behavior {
			name =  "cpCode"
			option {
				key =  "id"
				value = "${akamai_cp_code.cp_code.id}"
			}
			option {
				key =  "name"
				value = "${akamai_cp_code.cp_code.name}"
			}
		}
		behavior {
			name =  "caching"
			option {
				key =  "behavior"
				value = "MAX_AGE"
			}
			option {
                key =  "mustRevalidate"
                value = "false"
			}
            option {
                key =  "ttl"
                value = "1d"
            }
		}
    }
}
`)

type papiCall struct {
	methodName   string
	papiResponse interface{}
	papiRequest  interface{}
	error        error
	stubOnce     bool
}

func mockPAPIClient(callsToMock []papiCall) *mockpapi {
	client := &mockpapi{}
	for _, call := range callsToMock {
		var request interface{}
		request = mock.Anything
		if call.papiRequest != nil {
			request = call.papiRequest
		}
		stub := client.On(call.methodName, AnyCTX, request).Return(call.papiResponse, call.error)
		if call.stubOnce {
			stub.Once()
		}
	}

	return client
}

func TestResourcePropertyActivationCreate(t *testing.T) {
	t.Run("check schema property activation - OK", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - papi error", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName:   "GetRuleTree",
				papiResponse: nil,
				error:        fmt.Errorf("failed to create request"),
				stubOnce:     false,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("failed to create request"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - no property id nor property", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/no_propertyId/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("ExactlyOne"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - no contact", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/no_contact/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("Missing required argument"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("schema with `property` instead of `property_id`", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{Errors: make([]*papi.Error, 0)},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation_deprecated_arg.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation update", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiRequest: papi.GetActivationsRequest{
					PropertyID: "prp_test",
				},
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 2,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "CreateActivation",
				papiRequest: papi.CreateActivationRequest{
					PropertyID: "prp_test",
					Activation: papi.Activation{
						ActivationType:         papi.ActivationTypeActivate,
						AcknowledgeAllWarnings: true,
						PropertyVersion:        2,
						Network:                "STAGING",
						NotifyEmails:           []string{"user@example.com"},
						Note:                   "property activation note for updating",
					},
				},
				papiResponse: &papi.CreateActivationResponse{
					ActivationID: "atv_update",
				},
				stubOnce: true,
			},
			{
				methodName: "GetActivation",
				papiRequest: papi.GetActivationRequest{
					PropertyID:   "prp_test",
					ActivationID: "atv_update",
				},
				papiResponse: &papi.GetActivationResponse{
					GetActivationsResponse: papi.GetActivationsResponse{},
					Activation: &papi.Activation{
						ActivationID:    "atv_update",
						PropertyID:      "prp_test",
						PropertyVersion: 2,
						Network:         "STAGING",
						Status:          papi.ActivationStatusActive,
					},
				},
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation with rule errors", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors: []*papi.Error{
							{
								Title: "some error",
							},
						},
					},
				},
				error:    nil,
				stubOnce: false,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("activation cannot continue due to rule errors"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAccAkamaiPropertyActivation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyActivationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyActivationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiPropertyActivationExists,
				),
			},
		},
	})
}

var testAccAkamaiPropertyActivationConfigLatest = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

resource "akamai_property_activation" "property_activation" {
	property = "${akamai_property.property.id}"
	network = "STAGING"
	activate = true
	contact = ["dshafik@akamai.com"]
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_cp_code" "cp_code" {
	name = "terraform-testing3"
	contract = "${data.akamai_contract.contract.id}"
	group = "${data.akamai_group.group.id}"
	product = "prd_SPM"
}

resource "akamai_edge_hostname" "test" {
    product = "prd_SPM"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
    edge_hostname =  "terraform-test3.exampleterraform.io.edgesuite.net"
    ipv6 = true
}

resource "akamai_property" "property" {
  name = "terraform-test3"

  product = "prd_SPM"
  cp_code = "${akamai_cp_code.cp_code.id}"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"

  hostnames = {
		"example.org" = ${akamai_edge_hostname.test.edge_hostname}"
  }
  
  rule_format = "v2016-11-15"
  
  rules = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
 	rules {
		behavior {
			name =  "origin"
        	option { 
       			key =  "cacheKeyHostname"
            	value = "ORIGIN_HOSTNAME"
        	}
			option { 
    			key =  "compress"
     			value = "true"
     		}
    		option { 
    			key =  "enableTrueClientIp"
     			value = "false"
     		}
    		option { 
    			key =  "forwardHostHeader"
     			value = "REQUEST_HOST_HEADER"
     		}
    		option { 
    			key =  "hostname"
     			value = "exampleterraform.io"
     		}
    		option { 
    			key =  "httpPort"
     			value = "80"
     		}
    		option { 
    			key =  "httpsPort"
     			value = "443"
     		}
    		option { 
    			key =  "originSni"
     			value = "true"
     		}
    		option { 
    			key =  "originType"
     			value = "CUSTOMER"
     		}
    		option { 
    			key =  "verificationMode"
     			value = "PLATFORM_SETTINGS"
     		}
    		option { 
    			key =  "originCertificate"
     			value = ""
     		}
    		option { 
    			key =  "ports"
     			value = ""
     		}
      	}
		behavior {
			name =  "cpCode"
			option {
				key =  "id"
				value = "${akamai_cp_code.cp_code.id}"
			}
			option {
				key =  "name"
				value = "${akamai_cp_code.cp_code.name}"
			}
		}
		behavior {
			name =  "caching"
			option {
				key =  "behavior"
				value = "MAX_AGE"
			}
			option {
                key =  "mustRevalidate"
                value = "false"
			}
            option {
                key =  "ttl"
                value = "1d"
            }
		}
    }
}
`)

func TestAccAkamaiPropertyActivation_latest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyActivationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyActivationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiPropertyActivationExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiPropertyActivationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property_activation" {
			continue
		}

		log.Printf("[DEBUG] [Akamai PropertyActivation] Activation Delete")

	}
	return nil
}

func testAccCheckAkamaiPropertyActivationExists(s *terraform.State) error {
	// TODO: rewrite for v2???
	return nil
}

func testAccCheckAkamaiPropertyActivationLatest(s *terraform.State) error {
	// TODO: rewrite for v2???
	return nil
}

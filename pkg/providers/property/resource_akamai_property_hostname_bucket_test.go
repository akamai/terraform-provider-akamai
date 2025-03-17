package property

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var basicChecker = test.NewStateChecker("akamai_property_hostname_bucket.test").
	CheckEqual("property_id", "prp_111").
	CheckEqual("contract_id", "ctr_222").
	CheckEqual("group_id", "grp_333").
	CheckEqual("network", "STAGING").
	CheckEqual("note", "   ").
	CheckEqual("notify_emails.#", "1").
	CheckEqual("notify_emails.0", "nomail@akamai.com").
	CheckEqual("activation_id", "act_0").
	CheckEqual("timeout_for_activation", "50").
	CheckEqual("id", "prp_111:STAGING").
	CheckEqual("hostname_count", "1").
	CheckEqual("pending_default_certs", "0").
	CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id", "ehn_444").
	CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED")

func TestHostnameBucketResource_Create(t *testing.T) {
	// decrease timeout and intervals for tests
	forceTimeoutDuration = time.Second
	getHostnameBucketActivationInterval = time.Second
	t.Parallel()

	tests := map[string]struct {
		init            func(*mockProperty)
		checksForCreate resource.TestCheckFunc
		configFile      string
		expectError     *regexp.Regexp
	}{
		"create with 1 hostname on STAGING": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.Build(),
			configFile:      "1.tf",
		},
		"create with 1000 hostnames on STAGING": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1000, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostname_count", "1000").
				Build(),
			configFile: "1000.tf",
		},
		"create with 1 hostname without prefixes on STAGING": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("property_id", "111").
				CheckEqual("contract_id", "222").
				CheckEqual("group_id", "333").
				CheckEqual("id", "111:STAGING").
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id", "444").
				Build(),
			configFile: "1_without_prefixes.tf",
		},
		"create with 1 hostname - PRODUCTION and custom optional attributes": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1, "DEFAULT", "ehn_444"),
					network:      "PRODUCTION",
					notifyEmails: []string{"test1@nomail.com", "test2@nomail.com"},
					note:         "Test note",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("network", "PRODUCTION").
				CheckEqual("note", "Test note").
				CheckEqual("notify_emails.#", "2").
				CheckEqual("notify_emails.0", "test1@nomail.com").
				CheckEqual("notify_emails.1", "test2@nomail.com").
				CheckEqual("timeout_for_activation", "30").
				CheckEqual("id", "prp_111:PRODUCTION").
				CheckEqual("pending_default_certs", "1").
				Build(),
			configFile: "1_with_all_attributes.tf",
		},
		"create with 1100 hostnames - 2 patch requests": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1100, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.1000.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1000.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.1099.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1099.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_1").
				CheckEqual("hostname_count", "1100").
				CheckMissing("hostnames.www.test.hostname.1100.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.1100.com.edgesuite.net").
				Build(),
			configFile: "1100.tf",
		},
		"create with 1 hostname on STAGING: change activation status from PENDING to ACTIVE": {
			init: func(p *mockProperty) {
				setUpInitialData(p)
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// First GET request returns PENDING status
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "PENDING",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Once()
				// Second GET request returns ACTIVE status
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "ACTIVE",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Once()
				// Update the "state" for the next mock API calls
				p.hostnameBucket.activations = papi.ListPropertyHostnameActivationsResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					HostnameActivations: papi.HostnameActivationsList{
						Items: []papi.HostnameActivationListItem{
							{
								ActivationType:       "ACTIVATE",
								HostnameActivationID: "act_0",
								PropertyID:           p.propertyID,
								Network:              "STAGING",
								Status:               "ACTIVE",
								Note:                 "   ",
								NotifyEmails:         []string{"nomail@akamai.com"},
							},
						},
						TotalItems:       1,
						CurrentItemCount: 1,
					},
				}
				p.hostnameBucket.state["www.test.hostname.0.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.0.to.com.edgesuite.net"),
				}
				p.mockListActivePropertyHostnames()
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("pending_default_certs", "1").
				CheckMissing("hostnames.www.test.hostname.1.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.1.com.edgesuite.net").
				Build(),
			configFile: "1_with_default_cert.tf",
		},
		"create with 1 hostname on STAGING, timeout when waiting for activation: send successful CANCEL pending activation request - expect error": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// Mock two GET requests returning PENDING status, causing exceeding the deadline timeout
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					HostnameActivationID: "act_0",
					ContractID:           p.contractID,
					GroupID:              p.groupID,
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "PENDING",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Twice()
				// Mock CANCEL pending activation
				cancelReq := papi.CancelPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("CancelPropertyHostnameActivation", testutils.MockContext, cancelReq).Return(&papi.CancelPropertyHostnameActivationResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationCancelItem{
						HostnameActivationID: "act_1",
						Status:               "PENDING",
					},
				}, nil).Once()
				// Mock GET cancel activation returning status ABORTED, resulting in a final error
				getCancelActivationReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_1",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getCancelActivationReq).Return(&papi.GetPropertyHostnameActivationResponse{
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_1",
						Status:               "ABORTED",
					},
				}, nil).Once()
			},
			expectError: regexp.MustCompile(errCancelActivation.Error()),
			configFile:  "1_with_default_cert.tf",
		},
		"create with 1 hostname on STAGING, CANCEL request returns ErrActivationTooFar, next GetPropertyHostnameActivation returns ACTIVE activation": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// Mock two GET requests returning PENDING status, causing exceeding the deadline timeout
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "PENDING",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Twice()
				// Mock CANCEL pending activation returns ErrActivationTooFar
				cancelReq := papi.CancelPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("CancelPropertyHostnameActivation", testutils.MockContext, cancelReq).Return(nil, papi.ErrActivationTooFar).Once()
				// Mock GetPropertyHostnameActivation returns ACTIVE status
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "ACTIVE",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Once()
				// Update the "state" for the next mock API calls
				p.hostnameBucket.activations = papi.ListPropertyHostnameActivationsResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					HostnameActivations: papi.HostnameActivationsList{
						Items: []papi.HostnameActivationListItem{
							{
								ActivationType:       "ACTIVATE",
								HostnameActivationID: "act_0",
								PropertyID:           p.propertyID,
								Network:              "STAGING",
								Status:               "ACTIVE",
								Note:                 "   ",
								NotifyEmails:         []string{"nomail@akamai.com"},
							},
						},
						TotalItems:       1,
						CurrentItemCount: 1,
					},
				}
				p.hostnameBucket.state["www.test.hostname.0.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.0.to.com.edgesuite.net"),
				}
				p.mockListActivePropertyHostnames()
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("pending_default_certs", "1").
				Build(),
			configFile: "1_with_default_cert.tf",
		},
		"create with 1 hostname on STAGING, CANCEL request returns ErrActivationAlreadyActive, next GetPropertyHostnameActivation returns ACTIVE activation": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// Mock two GET requests returning PENDING status, exceeding the deadline timeout
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "PENDING",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Twice()
				// Mock CANCEL pending activation returns ErrActivationAlreadyActive
				cancelReq := papi.CancelPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("CancelPropertyHostnameActivation", testutils.MockContext, cancelReq).Return(nil, papi.ErrActivationAlreadyActive).Once()
				// Mock GetPropertyHostnameActivation returns ACTIVE status
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "ACTIVE",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Once()
				// Update the "state" for the next mock API calls
				p.hostnameBucket.activations = papi.ListPropertyHostnameActivationsResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					HostnameActivations: papi.HostnameActivationsList{
						Items: []papi.HostnameActivationListItem{
							{
								ActivationType:       "ACTIVATE",
								HostnameActivationID: "act_0",
								PropertyID:           p.propertyID,
								Network:              "STAGING",
								Status:               "ACTIVE",
								Note:                 "   ",
								NotifyEmails:         []string{"nomail@akamai.com"},
							},
						},
						TotalItems:       1,
						CurrentItemCount: 1,
					},
				}
				p.hostnameBucket.state["www.test.hostname.0.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.0.to.com.edgesuite.net"),
				}
				p.mockListActivePropertyHostnames()
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("pending_default_certs", "1").
				Build(),
			configFile: "1_with_default_cert.tf",
		},
		"expect error - CancelPendingHostnameActivation request returns error different than ErrActivationAlreadyActive and ErrActivationTooFar": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// Mock two GET requests returning PENDING status, causing exceeding the deadline timeout
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(&papi.GetPropertyHostnameActivationResponse{
					ContractID: p.propertyID,
					GroupID:    p.groupID,
					HostnameActivation: papi.HostnameActivationGetItem{
						HostnameActivationID: "act_0",
						Status:               "PENDING",
						PropertyID:           p.propertyID,
						Network:              p.hostnameBucket.network,
						Note:                 p.hostnameBucket.note,
						NotifyEmails:         p.hostnameBucket.notifyEmails,
					},
				}, nil).Twice()
				// Mock CANCEL pending activation returns ErrActivationAlreadyActive
				cancelReq := papi.CancelPropertyHostnameActivationRequest{
					PropertyID:           p.propertyID,
					ContractID:           p.contractID,
					GroupID:              p.groupID,
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("CancelPropertyHostnameActivation", testutils.MockContext, cancelReq).Return(nil, fmt.Errorf("API error")).Once()
			},
			expectError: regexp.MustCompile(`API error`),
			configFile:  "1_with_default_cert.tf",
		},
		"expect error - GetPropertyHostnameBucketActivation": {
			init: func(p *mockProperty) {
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(&papi.PatchPropertyHostnameBucketResponse{
					ActivationID: "act_0",
				}, nil).Once()
				// Mock an error on GetPropertyHostnameActivation call
				getReq := papi.GetPropertyHostnameActivationRequest{
					PropertyID:           "prp_111",
					ContractID:           "ctr_222",
					GroupID:              "grp_333",
					HostnameActivationID: "act_0",
				}
				p.papiMock.On("GetPropertyHostnameActivation", testutils.MockContext, getReq).Return(nil, fmt.Errorf("API error")).Once()
			},
			expectError: regexp.MustCompile(`API error`),
			configFile:  "1_with_default_cert.tf",
		},
		"expect error - ListActivePropertyHostnames call in Create": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				p.mockPatchPropertyHostnameBucket()
				p.papiMock.On("ListActivePropertyHostnames", testutils.MockContext, papi.ListActivePropertyHostnamesRequest{
					PropertyID:        p.propertyID,
					ContractID:        p.contractID,
					GroupID:           p.groupID,
					Limit:             999,
					Sort:              "hostname:a",
					Network:           "STAGING",
					IncludeCertStatus: true,
				}).Return(nil, fmt.Errorf("API error"))
			},
			expectError: regexp.MustCompile(`API error`),
			configFile:  "1.tf",
		},
		"expect error - default cert limit exceeded - do not retry": {
			init: func(p *mockProperty) {
				// Create
				req := createDefaultPatchPropertyHostnameBucketRequest()
				p.papiMock.On("PatchPropertyHostnameBucket", testutils.MockContext, req).Return(nil, papi.ErrDefaultCertLimitReached).Once()
			},
			expectError: regexp.MustCompile(`the limit for DEFAULT certificates has been reached`),
			configFile:  "1_with_default_cert.tf",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			tc.init(&mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/%s", tc.configFile),
							Check:       tc.checksForCreate,
							ExpectError: tc.expectError,
						},
					},
				})
			})
			papiMock.AssertExpectations(t)
		})
	}
}

func TestHostnameBucketResource_Update(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init            func(*mockProperty)
		checksForCreate resource.TestCheckFunc
		checksForUpdate resource.TestCheckFunc
		createConfig    string
		updateConfig    string
		createError     *regexp.Regexp
		updateError     *regexp.Regexp
	}{
		"create 1, update by adding 3 hostnames": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan["www.test.hostname.1.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.1.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.2.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.2.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.3.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.3.cnameTo.com.edgesuite.net"),
				}
				// Update
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_0_update").
				CheckEqual("hostname_count", "4").
				CheckEqual("pending_default_certs", "1").
				Build(),
			createConfig: "1.tf",
			updateConfig: "add_3.tf",
		},
		"create 1, update by adding 3 hostnames without group_id and contract_id, but receive the values from the API": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				// Create
				p.mockPatchPropertyHostnameBucket()
				p.hostnameBucket.state = maps.Clone(p.hostnameBucket.plan)
				// Send the request with empty contract and group
				req := papi.ListActivePropertyHostnamesRequest{
					PropertyID:        p.propertyID,
					Offset:            0,
					Limit:             999,
					Network:           papi.NetworkType(p.hostnameBucket.network),
					IncludeCertStatus: true,
					Sort:              "hostname:a",
				}
				// In the response we receive the values for the group and contract
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				resp := papi.ListActivePropertyHostnamesResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					PropertyID: p.propertyID,
					Hostnames: papi.HostnamesResponseItems{
						Items: []papi.HostnameItem{
							{
								CnameFrom:             "www.test.hostname.0.com.edgesuite.net",
								CnameType:             "EDGE_HOSTNAME",
								StagingCertType:       "CPS_MANAGED",
								StagingCnameTo:        "www.test.hostname.0.to.com.edgesuite.net",
								StagingEdgeHostnameId: "ehn_444",
							},
						},
						TotalItems: 1,
					},
				}
				p.papiMock.On("ListActivePropertyHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan["www.test.hostname.1.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.1.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.2.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.2.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.3.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.3.cnameTo.com.edgesuite.net"),
				}
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_0_update").
				CheckEqual("hostname_count", "4").
				CheckEqual("pending_default_certs", "1").
				Build(),
			createConfig: "1_without_contract_and_group.tf",
			updateConfig: "add_3_without_contract_and_group.tf",
		},
		"create with 1100 hostnames, update by removing 3": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1100, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				delete(p.hostnameBucket.plan, "www.test.hostname.999.com.edgesuite.net")
				delete(p.hostnameBucket.plan, "www.test.hostname.1098.com.edgesuite.net")
				delete(p.hostnameBucket.plan, "www.test.hostname.1099.com.edgesuite.net")
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.1000.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1000.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.1099.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1099.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_1").
				CheckEqual("hostname_count", "1100").
				CheckMissing("hostnames.www.test.hostname.1100.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.1100.com.edgesuite.net").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("activation_id", "act_0_update").
				CheckEqual("hostname_count", "1097").
				CheckMissing("hostnames.www.test.hostname.1100.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.1099.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.1098.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.999.com.edgesuite.net").
				Build(),
			createConfig: "1100.tf",
			updateConfig: "1100_remove_3.tf",
		},
		"create 5, update by removing and adding 2 hostnames": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(5, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan["www.test.hostname.5.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.5.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.6.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_666"),
					CnameTo:              types.StringValue("www.test.hostname.6.cnameTo.com.edgesuite.net"),
				}
				delete(p.hostnameBucket.plan, "www.test.hostname.1.com.edgesuite.net")
				delete(p.hostnameBucket.plan, "www.test.hostname.3.com.edgesuite.net")
				p.mockPatchPropertyHostnameBucket()
				p.hostnameBucket.state = maps.Clone(p.hostnameBucket.plan)
				p.mockListActivePropertyHostnames()
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostname_count", "5").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.5.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.5.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.6.com.edgesuite.net.edge_hostname_id", "ehn_666").
				CheckEqual("hostnames.www.test.hostname.6.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("activation_id", "act_0_update").
				CheckEqual("hostname_count", "5").
				CheckEqual("pending_default_certs", "2").
				CheckMissing("hostnames.www.test.hostname.1.com.edgesuite.net").
				CheckMissing("hostnames.www.test.hostname.3.com.edgesuite.net").
				Build(),
			createConfig: "5.tf",
			updateConfig: "add_2_remove_2.tf",
		},
		"update 2 existing keys": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(5, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan["www.test.hostname.3.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.33.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.2.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.2.cnameTo.com.edgesuite.net"),
				}
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostname_count", "5").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_1_update").
				CheckEqual("hostname_count", "5").
				CheckEqual("pending_default_certs", "1").
				Build(),
			createConfig: "5.tf",
			updateConfig: "update_2_existing.tf",
		},
		"create 5 hostnames, update by removing, adding and updating hostnames": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(5, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				delete(p.hostnameBucket.plan, "www.test.hostname.0.com.edgesuite.net")
				delete(p.hostnameBucket.plan, "www.test.hostname.1.com.edgesuite.net")
				p.hostnameBucket.plan["www.test.hostname.2.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.22.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.3.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.33.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.5.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("DEFAULT"),
					EdgeHostnameID:       types.StringValue("ehn_555"),
					CnameTo:              types.StringValue("www.test.hostname.5.cnameTo.com.edgesuite.net"),
				}
				p.hostnameBucket.plan["www.test.hostname.6.com.edgesuite.net"] = Hostname{
					CertProvisioningType: types.StringValue("CPS_MANAGED"),
					EdgeHostnameID:       types.StringValue("ehn_444"),
					CnameTo:              types.StringValue("www.test.hostname.6.cnameTo.com.edgesuite.net"),
				}
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.1.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostname_count", "5").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.2.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.3.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.5.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.5.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.6.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.6.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_1_update").
				CheckEqual("hostname_count", "5").
				CheckEqual("pending_default_certs", "2").
				CheckMissing("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type").
				CheckMissing("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id").
				Build(),
			createConfig: "5.tf",
			updateConfig: "update_2_add_2_remove_2.tf",
		},
		"create 1000 hostnames, update by adding 4000 hostnames and changing 1000 hostnames": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1000, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan = generateHostnames(5000, "CPS_MANAGED", "ehn_555")
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostname_count", "1000").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("hostnames.www.test.hostname.4999.com.edgesuite.net.edge_hostname_id", "ehn_555").
				CheckEqual("hostnames.www.test.hostname.4999.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_5_update").
				CheckEqual("hostname_count", "5000").
				Build(),
			createConfig: "1000.tf",
			updateConfig: "update_1000_add_4000.tf",
		},
		"create 5000 hostnames, update by removing 4000 hostnames and changing 1000 hostnames": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(5000, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update
				p.hostnameBucket.plan = generateHostnames(1000, "DEFAULT", "ehn_444")
				mockResourceHostnameBucketUpsert(p)
				// Read
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			checksForCreate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.4999.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.4999.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
				CheckEqual("activation_id", "act_4").
				CheckEqual("hostname_count", "5000").
				Build(),
			checksForUpdate: basicChecker.
				CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.edge_hostname_id", "ehn_444").
				CheckEqual("hostnames.www.test.hostname.999.com.edgesuite.net.cert_provisioning_type", "DEFAULT").
				CheckEqual("activation_id", "act_5_update").
				CheckEqual("hostname_count", "1000").
				CheckEqual("pending_default_certs", "1000").
				Build(),
			createConfig: "5000.tf",
			updateConfig: "update_1000_remove_4000.tf",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			tc.init(&mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/%s", tc.createConfig),
							Check:       tc.checksForCreate,
							ExpectError: tc.createError,
						},
						{
							Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/update/%s", tc.updateConfig),
							Check:       tc.checksForUpdate,
							ExpectError: tc.updateError,
						},
					},
				})
			})
			papiMock.AssertExpectations(t)
		})
	}
}

func TestHostnameBucketResource_Import(t *testing.T) {
	t.Parallel()
	importChecker := test.NewImportChecker().
		CheckEqual("id", "prp_111:STAGING").
		CheckEqual("property_id", "prp_111").
		CheckEqual("contract_id", "ctr_222").
		CheckEqual("group_id", "grp_333").
		CheckEqual("note", "   ").
		CheckEqual("notify_emails.#", "1").
		CheckEqual("notify_emails.0", "nomail@akamai.com").
		CheckEqual("activation_id", "act_0").
		CheckEqual("timeout_for_activation", "50").
		CheckEqual("hostname_count", "100").
		CheckEqual("pending_default_certs", "0").
		CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.edge_hostname_id", "ehn_444").
		CheckEqual("hostnames.www.test.hostname.0.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
		CheckEqual("hostnames.www.test.hostname.99.com.edgesuite.net.edge_hostname_id", "ehn_444").
		CheckEqual("hostnames.www.test.hostname.99.com.edgesuite.net.cert_provisioning_type", "CPS_MANAGED").
		CheckMissing("hostnames.www.test.hostname.100.com.edgesuite.net.edge_hostname_id").
		CheckMissing("hostnames.www.test.hostname.100.com.edgesuite.net.cert_provisioning_type")

	tests := map[string]struct {
		importID    string
		init        func(*mockProperty)
		stateCheck  func(s []*terraform.InstanceState) error
		expectError *regexp.Regexp
	}{
		"import with prefixed property and STAGING network": {
			importID: "prp_111:STAGING",
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "STAGING",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									ActivationType:       "ACTIVATE",
									HostnameActivationID: "act_0",
									PropertyID:           "prp_111",
									Network:              "STAGING",
									Status:               "ACTIVE",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
				}
				p.propertyID = "prp_111"
				p.mockListPropertyHostnameActivations()
				// fill contract and group attributes to be used in the next API call
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListActivePropertyHostnames(true)
			},
			stateCheck: importChecker.Build(),
		},
		"import with un-prefixed property and PRODUCTION network": {
			importID: "111:PRODUCTION",
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "PRODUCTION",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									ActivationType:       "ACTIVATE",
									HostnameActivationID: "act_0",
									PropertyID:           "prp_111",
									Network:              "PRODUCTION",
									Status:               "ACTIVE",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
				}
				p.propertyID = "prp_111"
				p.mockListPropertyHostnameActivations()
				// fill contract and group attributes to be used in the next API call
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListActivePropertyHostnames(true)
			},
			stateCheck: importChecker.
				CheckEqual("property_id", "111").
				CheckEqual("id", "111:PRODUCTION").
				Build(),
		},
		"import with whole prefixed importID and STAGING network": {
			importID: "prp_111:STAGING:ctr_222:grp_333",
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "STAGING",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									HostnameActivationID: "act_0",
									Network:              "STAGING",
									Status:               "ACTIVE",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListPropertyHostnameActivations()
				p.mockListActivePropertyHostnames()
			},
			stateCheck: importChecker.Build(),
		},
		"import with whole un-prefixed importID and PRODUCTION network": {
			importID: "111:PRODUCTION:222:333",
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "PRODUCTION",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									HostnameActivationID: "act_0",
									Network:              "PRODUCTION",
									Status:               "ACTIVE",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
				}
				p.propertyID = "prp_111"
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListPropertyHostnameActivations()
				p.mockListActivePropertyHostnames()
			},
			stateCheck: importChecker.
				CheckEqual("property_id", "111").
				CheckEqual("group_id", "333").
				CheckEqual("contract_id", "222").
				CheckEqual("id", "111:PRODUCTION").
				Build(),
		},
		"error on ListHostnameActivations": {
			importID:    "prp_111:PRODUCTION",
			expectError: regexp.MustCompile(`API error`),
			init: func(p *mockProperty) {
				p.propertyID = "prp_111"
				p.papiMock.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					PropertyID: "prp_111",
					Offset:     0,
					Limit:      999,
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
		},
		"no ACTIVE hostname activation found, no more pages to query": {
			importID:    "prp_111:PRODUCTION",
			expectError: regexp.MustCompile(`there is no active hostname activation for given property`),
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "PRODUCTION",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									HostnameActivationID: "act_0",
									Network:              "PRODUCTION",
									Status:               "ABORTED",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
				}
				p.propertyID = "prp_111"
				p.mockListPropertyHostnameActivations()
			},
		},
		"ACTIVE hostname activation found on the next page": {
			importID: "prp_111:PRODUCTION",
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "PRODUCTION",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							TotalItems: 1000,
							Items:      generateAbortedActivations(999, "PRODUCTION"),
						},
					},
				}
				p.propertyID = "prp_111"
				// Mock that first API call returns 999 aborted activations, so we need to invoke another paged API call
				firstGetReq := papi.ListPropertyHostnameActivationsRequest{
					PropertyID: p.propertyID,
					ContractID: p.contractID,
					GroupID:    p.groupID,
					Limit:      999,
				}
				firstGetResp := p.hostnameBucket.activations
				p.papiMock.On("ListPropertyHostnameActivations", testutils.MockContext, firstGetReq).Return(&firstGetResp, nil).Once()
				// Mock that second API call returns an ACTIVE activation that can be used
				secondGetReq := papi.ListPropertyHostnameActivationsRequest{
					PropertyID: p.propertyID,
					ContractID: p.contractID,
					GroupID:    p.groupID,
					Offset:     999,
					Limit:      999,
				}
				secondGetResp := papi.ListPropertyHostnameActivationsResponse{
					ContractID: "ctr_222",
					GroupID:    "grp_333",
					HostnameActivations: papi.HostnameActivationsList{
						Items: []papi.HostnameActivationListItem{
							{
								HostnameActivationID: "act_1000",
								Network:              "PRODUCTION",
								Status:               "ACTIVE",
								Note:                 "   ",
								NotifyEmails:         []string{"nomail@akamai.com"},
							},
						},
						TotalItems: 1000,
					},
				}
				p.papiMock.On("ListPropertyHostnameActivations", testutils.MockContext, secondGetReq).Return(&secondGetResp, nil).Once()
				// fill contract and group attributes to be used in the next API call
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListActivePropertyHostnames(true)
			},
			stateCheck: importChecker.
				CheckEqual("property_id", "prp_111").
				CheckEqual("id", "prp_111:PRODUCTION").
				CheckEqual("activation_id", "act_1000").
				Build(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			tc.init(&mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    tc.importID,
							ImportState:      true,
							ResourceName:     "akamai_property_hostname_bucket.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResPropertyHostnameBucket/import/default.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			papiMock.AssertExpectations(t)
		})
	}
}

func TestHostnameBucketResource_ValidationErrors(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		init            func(*mockProperty)
		checksForCreate resource.TestCheckFunc
		steps           []resource.TestStep
	}{
		"validation error - create with empty hostnames": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/empty_hostnames.tf"),
					ExpectError: regexp.MustCompile(`Attribute hostnames map must contain at least 1 elements and at most 99999\nelements, got: 0`),
				},
			},
		},
		"validation error - create with no hostnames": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/no_hostnames.tf"),
					ExpectError: regexp.MustCompile(`The argument "hostnames" is required, but no definition was found.`),
				},
			},
		},
		"validation error - create with group_id, but no contract_id": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_contract_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute "contract_id" must be specified when "group_id" is specified`),
				},
			},
		},
		"validation error - create with contract_id, but no group_id": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_group_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute "group_id" must be specified when "contract_id" is specified`),
				},
			},
		},
		"validation error - create with group and contract, update group": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/update_group.tf"),
					ExpectError: regexp.MustCompile(`updating 'group_id' is not allowed`),
				},
			},
		},
		"validation error - create without group and contract, update contract": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					plan:         generateHostnames(1, "CPS_MANAGED", "ehn_444"),
					network:      "STAGING",
					notifyEmails: []string{"nomail@akamai.com"},
					note:         "   ",
					state:        map[string]Hostname{},
				}
				p.propertyID = "prp_111"
				// Create
				p.mockPatchPropertyHostnameBucket()
				p.hostnameBucket.state = maps.Clone(p.hostnameBucket.plan)
				// Send the request with empty contract and group
				req := papi.ListActivePropertyHostnamesRequest{
					PropertyID:        p.propertyID,
					Offset:            0,
					Limit:             999,
					Network:           papi.NetworkType(p.hostnameBucket.network),
					IncludeCertStatus: true,
					Sort:              "hostname:a",
				}
				// In the response we receive the values for the group and contract
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				resp := papi.ListActivePropertyHostnamesResponse{
					ContractID: p.contractID,
					GroupID:    p.groupID,
					PropertyID: p.propertyID,
					Hostnames: papi.HostnamesResponseItems{
						Items: []papi.HostnameItem{
							{
								CnameFrom:             "www.test.hostname.0.com.edgesuite.net",
								CnameType:             "EDGE_HOSTNAME",
								StagingCertType:       "CPS_MANAGED",
								StagingCnameTo:        "www.test.hostname.0.to.com.edgesuite.net",
								StagingEdgeHostnameId: "ehn_444",
							},
						},
						TotalItems: 1,
					},
				}
				p.papiMock.On("ListActivePropertyHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1_without_contract_and_group.tf"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/update_contract.tf"),
					ExpectError: regexp.MustCompile(`updating 'contract_id' is not allowed`),
				},
			},
		},
		"validation error - updating to empty map": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Update and Delete as a 3rd step to enable destroy
				p.mockListActivePropertyHostnames()
				p.mockListPropertyHostnameActivations()
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/empty_hostnames.tf"),
					ExpectError: regexp.MustCompile(`Attribute hostnames map must contain at least 1 elements and at most 99999\nelements, got: 0`),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
				},
			},
		},
		"validation error - missing property_id": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_property_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "property_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing network": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_network.tf"),
					ExpectError: regexp.MustCompile(`The argument "network" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing cert_provisioning_type": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_cert_provisioning_type.tf"),
					ExpectError: regexp.MustCompile(`"www.test.hostname.0.com.edgesuite.net": attribute "cert_provisioning_type"\nis required.`),
				},
			},
		},
		"validation error - missing edge_hostname_id": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/validation/missing_edge_hostname_id.tf"),
					ExpectError: regexp.MustCompile(`"www.test.hostname.0.com.edgesuite.net": attribute "edge_hostname_id" is\nrequired.`),
				},
			},
		},
		"validation error - incorrect importID": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					ImportStateId: "111:PRODUCTION:222",
					ImportState:   true,
					ResourceName:  "akamai_property_hostname_bucket.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResPropertyHostnameBucket/import/default.tf"),
					ExpectError:   regexp.MustCompile(`importID must be of format 'property_id:network\[:contract_id:group_id]'`),
				},
			},
		},
		"validation error - incorrect network when importing": {
			init: func(_ *mockProperty) {},
			steps: []resource.TestStep{
				{
					ImportStateId: "111:STRANGE_NETWORK:222:333",
					ImportState:   true,
					ResourceName:  "akamai_property_hostname_bucket.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResPropertyHostnameBucket/import/default.tf"),
					ExpectError:   regexp.MustCompile(`network must have correct value of 'STAGING' or 'PRODUCTION'`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			tc.init(&mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			papiMock.AssertExpectations(t)
		})
	}
}

func TestHostnameBucketResource_Diff(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init  func(*mockProperty)
		steps []resource.TestStep
	}{
		"create basic, verify no diff on next plan": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
					Check:  basicChecker.Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
					ExpectNonEmptyPlan: false,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"create basic with custom timeout, verify no diff on next plan": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1_with_custom_timeout.tf"),
					Check: basicChecker.
						CheckEqual("timeout_for_activation", "30").
						Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1_with_custom_timeout.tf"),
					ExpectNonEmptyPlan: false,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"create with default timeout, update only timeout - verify no diff": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
					Check:  basicChecker.Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1_with_custom_timeout.tf"),
					ExpectNonEmptyPlan: false,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"create with custom timeout, verify no diff present in the next plan for timeout, note and notify_emails attributes when hostname entry is modified": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1_with_custom_timeout.tf"),
					Check: basicChecker.
						CheckEqual("timeout_for_activation", "30").
						Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/update/custom_timeout.tf"),
					ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							expectNoDiffOnTimeoutForActivation(),
							expectNoDiffOnNote(),
							expectNoDiffOnNotifyEmails(),
							expectDiffOnHostnames(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"create, update note - expect no diff": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
					Check: basicChecker.
						Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/update/note.tf"),
					ExpectNonEmptyPlan: false,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							expectNoDiffOnNote(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"create, update notify_emails - expect no diff": {
			init: func(p *mockProperty) {
				// Set up initial data for the property and hostname bucket
				setUpInitialData(p)
				// Create
				mockResourceHostnameBucketUpsert(p)
				// Read x2
				mockResourceHostnameBucketRead(p, 2)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/create/1.tf"),
					Check: basicChecker.
						Build(),
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/update/notify_emails.tf"),
					ExpectNonEmptyPlan: false,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PostApplyPostRefresh: []plancheck.PlanCheck{
							expectNoDiffOnNotifyEmails(),
						},
					},
					PlanOnly: true,
				},
			},
		},
		"expect no diff after import": {
			init: func(p *mockProperty) {
				p.hostnameBucket = hostnameBucket{
					state:   generateHostnames(100, "CPS_MANAGED", "ehn_444"),
					network: "STAGING",
					activations: papi.ListPropertyHostnameActivationsResponse{
						ContractID: "ctr_222",
						GroupID:    "grp_333",
						HostnameActivations: papi.HostnameActivationsList{
							Items: []papi.HostnameActivationListItem{
								{
									ActivationType:       "ACTIVATE",
									HostnameActivationID: "act_0",
									PropertyID:           "prp_111",
									Network:              "STAGING",
									Status:               "ACTIVE",
									Note:                 "   ",
									NotifyEmails:         []string{"nomail@akamai.com"},
								},
							},
							TotalItems: 1,
						},
					},
					note:         "   ",
					notifyEmails: []string{"nomail@akamai.com"},
				}
				// Read x1
				p.propertyID = "prp_111"
				p.mockListPropertyHostnameActivations()
				// fill contract and group attributes to be used in the next API call
				p.contractID = "ctr_222"
				p.groupID = "grp_333"
				p.mockListActivePropertyHostnames(true)
				// Read x1
				mockResourceHostnameBucketRead(p)
				// Delete
				mockResourceHostnameBucketDelete(p)
			},
			steps: []resource.TestStep{
				{
					ExpectNonEmptyPlan: false,
					ImportStateId:      "prp_111:STAGING",
					ImportState:        true,
					ResourceName:       "akamai_property_hostname_bucket.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResPropertyHostnameBucket/import/import_100.tf"),
					ImportStatePersist: true,
				},
				{
					Config:             testutils.LoadFixtureStringf(t, "testdata/TestResPropertyHostnameBucket/import/import_100.tf"),
					ExpectNonEmptyPlan: false,
					PlanOnly:           true,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			tc.init(&mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			papiMock.AssertExpectations(t)
		})
	}
}

func expectNoDiffOnTimeoutForActivation() timeoutForActivationDiffCheck {
	return timeoutForActivationDiffCheck{}
}

type timeoutForActivationDiffCheck struct{}

func (c timeoutForActivationDiffCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if len(req.Plan.ResourceChanges) > 0 {
		changeBefore := req.Plan.ResourceChanges[0].Change.Before
		changeMapBefore, ok := changeBefore.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeBefore)
			return
		}
		beforeVal, ok := changeMapBefore["timeout_for_activation"].(json.Number)
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to json.Number", changeMapBefore["timeout_for_activation"])
			return
		}

		changeAfter := req.Plan.ResourceChanges[0].Change.After
		changeMapAfter, ok := changeAfter.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeAfter)
			return
		}
		afterVal, ok := changeMapAfter["timeout_for_activation"].(json.Number)
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to json.Number", changeMapAfter["timeout_for_activation"])
			return
		}

		if beforeVal != afterVal {
			resp.Error = fmt.Errorf("'timeout_for_activation' produced a diff, but it should not")
		}
	}
}

func expectNoDiffOnNote() noteDiffCheck {
	return noteDiffCheck{}
}

type noteDiffCheck struct{}

func (c noteDiffCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if len(req.Plan.ResourceChanges) > 0 {
		changeBefore := req.Plan.ResourceChanges[0].Change.Before
		changeMapBefore, ok := changeBefore.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeBefore)
			return
		}
		beforeVal, ok := changeMapBefore["note"].(string)
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to string", changeMapBefore["note"])
			return
		}

		changeAfter := req.Plan.ResourceChanges[0].Change.After
		changeMapAfter, ok := changeAfter.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeAfter)
			return
		}
		afterVal, ok := changeMapAfter["note"].(string)
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to string", changeMapAfter["note"])
			return
		}

		if beforeVal != afterVal {
			resp.Error = fmt.Errorf("'note' produced a diff, but it should not")
		}
	}
}

func expectNoDiffOnNotifyEmails() notifyEmailsDiffCheck {
	return notifyEmailsDiffCheck{}
}

type notifyEmailsDiffCheck struct{}

func (c notifyEmailsDiffCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if len(req.Plan.ResourceChanges) > 0 {
		changeBefore := req.Plan.ResourceChanges[0].Change.Before
		changeMapBefore, ok := changeBefore.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeBefore)
			return
		}
		beforeValRaw, ok := changeMapBefore["notify_emails"].([]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to []interface{}", changeMapBefore["notify_emails"])
			return
		}
		emailsBefore := make([]string, 0, len(beforeValRaw))
		for _, e := range beforeValRaw {
			emailsBefore = append(emailsBefore, e.(string))
		}

		changeAfter := req.Plan.ResourceChanges[0].Change.After
		changeMapAfter, ok := changeAfter.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeAfter)
			return
		}
		afterValRaw, ok := changeMapAfter["notify_emails"].([]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast field of type %T to []interface{}", changeMapAfter["notify_emails"])
			return
		}
		emailsAfter := make([]string, 0, len(afterValRaw))
		for _, e := range afterValRaw {
			emailsAfter = append(emailsAfter, e.(string))
		}

		if !slices.Equal(emailsBefore, emailsAfter) {
			resp.Error = fmt.Errorf("'notify_emails' produced a diff, but it should not")
		}
	}
}

func expectDiffOnHostnames() hostnamesDiffCheck {
	return hostnamesDiffCheck{}
}

type hostnamesDiffCheck struct{}

func (c hostnamesDiffCheck) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if len(req.Plan.ResourceChanges) > 0 {
		changeBefore := req.Plan.ResourceChanges[0].Change.Before
		changeMapBefore, ok := changeBefore.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeBefore)
			return
		}
		beforeVal, ok := changeMapBefore["hostnames"].(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeMapBefore["hostnames"])
			return
		}
		beforeTFMap := make(map[string]Hostname)
		for k, v := range beforeVal {
			vMap, ok := v.(map[string]interface{})
			if !ok {
				resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", v)
				return
			}
			beforeTFMap[k] = Hostname{
				CertProvisioningType: types.StringValue(vMap["cert_provisioning_type"].(string)),
				EdgeHostnameID:       types.StringValue(vMap["edge_hostname_id"].(string)),
			}
		}

		changeAfter := req.Plan.ResourceChanges[0].Change.After
		changeMapAfter, ok := changeAfter.(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeAfter)
			return
		}
		afterVal, ok := changeMapAfter["hostnames"].(map[string]interface{})
		if !ok {
			resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", changeMapAfter["hostnames"])
			return
		}
		afterTFMap := make(map[string]Hostname)
		for k, v := range afterVal {
			vMap, ok := v.(map[string]interface{})
			if !ok {
				resp.Error = fmt.Errorf("could not cast struct of type %T to map[string]interface{}", v)
				return
			}
			afterTFMap[k] = Hostname{
				CertProvisioningType: types.StringValue(vMap["cert_provisioning_type"].(string)),
				EdgeHostnameID:       types.StringValue(vMap["edge_hostname_id"].(string)),
			}
		}

		beforeTFMapFinal, diags := types.MapValueFrom(ctx, hostnameObjectType, beforeTFMap)
		if diags.HasError() {
			resp.Error = fmt.Errorf("could not convert the structure into map: %v", diags.Errors())
			return
		}
		afterTFMapFinal, diags := types.MapValueFrom(ctx, hostnameObjectType, afterTFMap)
		if diags.HasError() {
			resp.Error = fmt.Errorf("could not convert the structure into map: %v", diags.Errors())
			return
		}

		if beforeTFMapFinal.Equal(afterTFMapFinal) {
			resp.Error = fmt.Errorf("'hostnames' did not produce a diff, but it should")
		}
	}
}

func generateAbortedActivations(n int, network string) []papi.HostnameActivationListItem {
	result := make([]papi.HostnameActivationListItem, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, papi.HostnameActivationListItem{
			HostnameActivationID: fmt.Sprintf("act_%d", i),
			PropertyID:           "prp_111",
			Network:              network,
			Status:               "ABORTED",
		})
	}

	return result
}

func generateHostnames(n int, certType, ehn string) map[string]Hostname {
	hostnames := make(map[string]Hostname, n)
	for i := 0; i < n; i++ {
		hostnames[fmt.Sprintf("www.test.hostname.%d.com.edgesuite.net", i)] = Hostname{
			CertProvisioningType: types.StringValue(certType),
			CnameTo:              types.StringValue(fmt.Sprintf("www.test.hostname.%d.cnameTo.com.edgesuite.net", i)),
			EdgeHostnameID:       types.StringValue(ehn),
		}
	}

	return hostnames
}

func createDefaultPatchPropertyHostnameBucketRequest() papi.PatchPropertyHostnameBucketRequest {
	return papi.PatchPropertyHostnameBucketRequest{
		PropertyID: "prp_111",
		ContractID: "ctr_222",
		GroupID:    "grp_333",
		Body: papi.PatchPropertyHostnameBucketBody{
			Add: []papi.PatchPropertyHostnameBucketAdd{
				{
					EdgeHostnameID:       "ehn_444",
					CertProvisioningType: papi.CertTypeDefault,
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "www.test.hostname.0.com.edgesuite.net",
				},
			},
			Network:      "STAGING",
			NotifyEmails: []string{"nomail@akamai.com"},
			Note:         "   ",
		},
	}
}

func setUpInitialData(p *mockProperty) {
	p.hostnameBucket = hostnameBucket{
		plan:         generateHostnames(1, "CPS_MANAGED", "ehn_444"),
		network:      "STAGING",
		notifyEmails: []string{"nomail@akamai.com"},
		note:         "   ",
		state:        map[string]Hostname{},
	}
	p.propertyID = "prp_111"
	p.contractID = "ctr_222"
	p.groupID = "grp_333"
}

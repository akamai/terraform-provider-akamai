package cps

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v12/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollmentDelete(t *testing.T) {
	t.Parallel()
	apiError := errors.New("API failure")
	guidance := "set `acknowledge_deletion_warnings` to true, run 'terraform apply', then run 'terraform destroy' again"
	deletionChanges := []string{"/cps/v2/enrollments/1/changes/2"}
	pending := &cps.GetEnrollmentResponse{PendingChanges: []cps.PendingChange{{Location: "/cps/v2/enrollments/1/changes/2"}}}
	getEnrollment := cps.GetEnrollmentRequest{EnrollmentID: 1}
	getStatus := cps.GetChangeStatusRequest{EnrollmentID: 1, ChangeID: 2}
	warnings := &cps.Change{StatusInfo: &cps.StatusInfo{Status: waitReviewCertWarning}}
	ackRequest := cps.AcknowledgementRequest{
		EnrollmentID: 1,
		ChangeID:     2,
		Acknowledgement: cps.Acknowledgement{
			Acknowledgement: cps.AcknowledgementAcknowledge,
		},
	}
	tests := []struct {
		name          string
		acknowledge   *bool
		removeChanges []string
		removeError   error
		init          func(*cps.Mock)
		wantError     string
	}{
		{
			name:          "deletion without pending changes",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(&cps.GetEnrollmentResponse{}, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(&cps.Change{StatusInfo: &cps.StatusInfo{Status: "processing"}}, nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:        "no deletion job ignores unrelated pending change",
			acknowledge: ptr.To(true),
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(&cps.GetEnrollmentResponse{
					PendingChanges: []cps.PendingChange{{Location: "/cps/v2/enrollments/1/changes/1"}},
				}, nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "acknowledge delayed warnings only once and wait for deletion",
			acknowledge:   ptr.To(true),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Times(3)
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(&cps.Change{StatusInfo: &cps.StatusInfo{Status: "processing"}}, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "retry acknowledgement after change is initially missing",
			acknowledge:   ptr.To(true),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Times(2)
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Twice()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(cps.ErrNotFound).Once()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "change disappears before enrollment",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(nil, cps.ErrNotFound).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "missing status is retried before acknowledging",
			acknowledge:   ptr.To(true),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Twice()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(&cps.Change{}, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "use deletion job despite unrelated pending change",
			acknowledge:   ptr.To(true),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(&cps.GetEnrollmentResponse{
					PendingChanges: []cps.PendingChange{
						{Location: "/cps/v2/enrollments/1/changes/1"},
						{Location: "/cps/v2/enrollments/1/changes/2"},
					},
				}, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, cps.ErrNotFound).Once()
			},
		},
		{
			name:          "cancel deletion when consent is omitted",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("CancelChange", testutils.MockContext, cps.CancelChangeRequest{EnrollmentID: 1, ChangeID: 2}).Return(&cps.CancelChangeResponse{}, nil).Once()
			},
			wantError: "enrollment deletion requires acknowledgement of warnings; the deletion change was canceled; " + guidance,
		},
		{
			name:          "cancel deletion when consent is false",
			acknowledge:   ptr.To(false),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("CancelChange", testutils.MockContext, cps.CancelChangeRequest{EnrollmentID: 1, ChangeID: 2}).Return(&cps.CancelChangeResponse{}, nil).Once()
			},
			wantError: "enrollment deletion requires acknowledgement of warnings; the deletion change was canceled; " + guidance,
		},
		{
			name:          "cancel fails",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("CancelChange", testutils.MockContext, cps.CancelChangeRequest{EnrollmentID: 1, ChangeID: 2}).Return(nil, apiError).Once()
			},
			wantError: "enrollment deletion requires acknowledgement of warnings, but canceling the deletion change failed: API failure; " + guidance,
		},
		{
			name:          "acknowledgement fails",
			acknowledge:   ptr.To(true),
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(warnings, nil).Once()
				client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, ackRequest).Return(apiError).Once()
			},
			wantError: "API failure",
		},
		{
			name:        "remove fails",
			removeError: apiError,
			wantError:   "API failure",
		},
		{
			name:          "initial enrollment lookup fails",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, apiError).Once()
			},
			wantError: "error checking enrollment deletion: API failure",
		},
		{
			name:          "subsequent enrollment lookup fails",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(&cps.GetEnrollmentResponse{}, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(&cps.Change{StatusInfo: &cps.StatusInfo{Status: "processing"}}, nil).Once()
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(nil, apiError).Once()
			},
			wantError: "error checking enrollment deletion: API failure",
		},
		{
			name:          "status lookup fails",
			removeChanges: deletionChanges,
			init: func(client *cps.Mock) {
				client.On("GetEnrollment", testutils.MockContext, getEnrollment).Return(pending, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, getStatus).Return(nil, apiError).Once()
			},
			wantError: "API failure",
		},
	}
	resources := map[string]*schema.Resource{
		"DV":          resourceCPSDVEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
		"third party": resourceCPSThirdPartyEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
	}
	for resourceName, resource := range resources {
		t.Run(resourceName, func(t *testing.T) {
			t.Parallel()
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					t.Parallel()
					client := edgegrid.NewTestClient()
					t.Cleanup(func() { client.CPS.AssertExpectations(t) })
					operationMeta, err := meta.New(session.Must(session.New()), hclog.NewNullLogger(), "test")
					require.NoError(t, err)
					operationMeta.SetClient(client)
					config := map[string]any{}
					if test.acknowledge != nil {
						config["acknowledge_deletion_warnings"] = *test.acknowledge
					}
					data := schema.TestResourceDataRaw(t, resource.Schema, config)
					data.SetId("1")
					assert.Equal(t, test.acknowledge != nil && *test.acknowledge, data.Get("acknowledge_deletion_warnings"))
					client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
						EnrollmentID: 1, AllowCancelPendingChanges: ptr.To(true),
					}).Return(&cps.RemoveEnrollmentResponse{Changes: test.removeChanges}, test.removeError).Once()
					if test.init != nil {
						test.init(client.CPS)
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					diagnostics := resource.DeleteContext(ctx, data, operationMeta)
					if test.wantError != "" {
						require.Len(t, diagnostics, 1)
						assert.Equal(t, test.wantError, diagnostics[0].Summary)
						assert.Equal(t, "1", data.Id())
					} else {
						assert.Empty(t, diagnostics)
					}
				})
			}
		})
	}
}

func TestDNSNamesRejectEmptyName(t *testing.T) {
	t.Parallel()

	dnsNamesSchema := networkConfiguration.Schema["dns_names"]
	dnsNameSchema, ok := dnsNamesSchema.Elem.(*schema.Schema)
	require.True(t, ok)

	assert.Empty(t, dnsNameSchema.ValidateDiagFunc("example.com", cty.Path{}))
	assert.NotEmpty(t, dnsNameSchema.ValidateDiagFunc("", cty.Path{}))
}

func TestSplitChallenges(t *testing.T) {
	t.Parallel()
	t.Run("Non empty challenges", func(t *testing.T) {
		t.Parallel()
		challenges := mockDVArray()
		gotHTTPChallenge, gotDNSChallenge := splitChallenges(challenges)
		wantHTTPChallenge := []challengeHTTP{
			{
				"full_path":     "http://TestFullPath",
				"response_body": "TestResponseBody",
				"domain":        "TestDomain",
			},
		}
		wantDNSChallenge := []challengeDNS{
			{
				"full_path":     "TestFullPath",
				"response_body": "TestResponseBody",
				"domain":        "TestDomain",
			},
		}
		assert.Equal(t, wantHTTPChallenge, gotHTTPChallenge)
		assert.Equal(t, wantDNSChallenge, gotDNSChallenge)
	})

	t.Run("Empty challenges", func(t *testing.T) {
		t.Parallel()
		challenges := mockEmptyDVArray()
		gotHTTPChallenge, gotDNSChallenge := splitChallenges(challenges)
		wantDNSChallenge := make([]challengeDNS, 0)
		wantHTTPChallenge := make([]challengeHTTP, 0)
		assert.Equal(t, wantHTTPChallenge, gotHTTPChallenge)
		assert.Equal(t, wantDNSChallenge, gotDNSChallenge)
	})
}

func TestNewChallenge(t *testing.T) {
	t.Parallel()
	challenge1 := cps.Challenge{
		Error:             "",
		FullPath:          "http://TestFullPath",
		RedirectFullPath:  "http://TestRedirectFullPath.com",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken123",
		Type:              "http-01",
		ValidationRecords: nil,
	}

	dv := cps.DV{
		Challenges:         []cps.Challenge{challenge1},
		Domain:             "TestDomain",
		Error:              "The domain TestDomain is not ready for HTTP validation.",
		Expires:            "2022-07-25T10:17:44Z",
		RequestTimestamp:   "2022-07-18T10:17:44Z",
		Status:             "Awaiting user",
		ValidatedTimestamp: "2022-07-19T09:35:29Z",
		ValidationStatus:   "DATA_NOT_READY",
	}

	gotChallenge := newChallenge(&challenge1, &dv)
	wantChallenge := challenge{
		"full_path":     "http://TestFullPath",
		"response_body": "TestResponseBody",
		"domain":        "TestDomain",
	}
	assert.Equal(t, wantChallenge, gotChallenge)
}

func TestReadAttrsDNSNamesTransitions(t *testing.T) {
	t.Parallel()

	dnsNamesSchema := networkConfiguration.Schema["dns_names"]
	assert.Equal(t, schema.TypeSet, dnsNamesSchema.Type)
	assert.True(t, dnsNamesSchema.Optional)
	assert.True(t, dnsNamesSchema.Computed)

	tests := map[string]struct {
		currentNetworkConfig map[string]interface{}
		responseDNSNames     []string
		enableForAllSANs     bool
		expectedDNSNames     []interface{}
	}{
		"all SANs enabled stores API-populated DNS names": {
			currentNetworkConfig: map[string]interface{}{
				"enable_for_all_sans": true,
				"geography":           "core",
			},
			responseDNSNames: []string{"test.akamai.com", "san.test.akamai.com"},
			enableForAllSANs: true,
			expectedDNSNames: []interface{}{"test.akamai.com", "san.test.akamai.com"},
		},
		"all SANs disabled stores explicit DNS names": {
			currentNetworkConfig: map[string]interface{}{
				"enable_for_all_sans": false,
				"dns_names":           []interface{}{"test.akamai.com"},
				"geography":           "core",
			},
			responseDNSNames: []string{"test.akamai.com"},
			enableForAllSANs: false,
			expectedDNSNames: []interface{}{"test.akamai.com"},
		},
		"reenabling all SANs replaces explicit DNS names with API-populated names": {
			currentNetworkConfig: map[string]interface{}{
				"enable_for_all_sans": true,
				"dns_names":           []interface{}{"test.akamai.com"},
				"geography":           "core",
			},
			responseDNSNames: []string{"test.akamai.com", "san.test.akamai.com"},
			enableForAllSANs: true,
			expectedDNSNames: []interface{}{"test.akamai.com", "san.test.akamai.com"},
		},
	}

	resources := map[string]*schema.Resource{
		"DV":          resourceCPSDVEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
		"third party": resourceCPSThirdPartyEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
	}
	for resourceName, resource := range resources {
		t.Run(resourceName, func(t *testing.T) {
			t.Parallel()
			for name, test := range tests {
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					d := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
						"network_configuration": []interface{}{test.currentNetworkConfig},
					})
					enrollment := getTestDVEnrollment()
					enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
						CloneDNSNames: test.enableForAllSANs,
						DNSNames:      test.responseDNSNames,
					}

					attrs, err := readAttrs(&enrollment, d)
					require.NoError(t, err)
					networkConfig := attrs["network_configuration"].([]interface{})[0].(map[string]interface{})
					assert.Equal(t, test.expectedDNSNames, networkConfig["dns_names"])

					require.NoError(t, d.Set("network_configuration", attrs["network_configuration"]))
					dnsNames := d.Get("network_configuration").([]interface{})[0].(map[string]interface{})["dns_names"].(*schema.Set)
					assert.ElementsMatch(t, test.expectedDNSNames, dnsNames.List())
				})
			}
		})
	}
}

func TestCloneDNSNamesDefaultUpgradeState(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		omitEnableForAllSANsFromState bool
		stateEnableForAllSANs         bool
		config                        map[string]interface{}
		wantDiff                      bool
	}{
		"omitted legacy state adopts the new default": {
			omitEnableForAllSANsFromState: true,
			config: map[string]interface{}{
				"network_configuration": []interface{}{map[string]interface{}{
					"geography": "core",
				}},
			},
		},
		"explicit false remains disabled without drift": {
			config: map[string]interface{}{
				"network_configuration": []interface{}{map[string]interface{}{
					"enable_for_all_sans": false,
					"geography":           "core",
				}},
			},
		},
		"explicit true remains enabled without drift": {
			stateEnableForAllSANs: true,
			config: map[string]interface{}{
				"network_configuration": []interface{}{map[string]interface{}{
					"enable_for_all_sans": true,
					"geography":           "core",
				}},
			},
		},
	}
	resources := map[string]*schema.Resource{
		"DV":          resourceCPSDVEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
		"third party": resourceCPSThirdPartyEnrollment(testPollChangeStatusInterval, testPollGetEnrollmentInterval),
	}

	for resourceName, resource := range resources {
		t.Run(resourceName, func(t *testing.T) {
			t.Parallel()
			for name, test := range tests {
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					legacyConfig := map[string]interface{}{
						"network_configuration": []interface{}{map[string]interface{}{
							"enable_for_all_sans": false,
							"geography":           "core",
						}},
					}
					legacySchema := make(map[string]*schema.Schema, len(resource.Schema))
					for attribute, field := range resource.Schema {
						legacySchema[attribute] = field
					}
					legacyNetworkConfiguration := *legacySchema["network_configuration"]
					legacyNetworkConfiguration.Type = schema.TypeSet
					legacyNetworkConfiguration.Set = hashLegacyNetworkConfiguration
					legacySchema["network_configuration"] = &legacyNetworkConfiguration

					legacyData := schema.TestResourceDataRaw(t, legacySchema, legacyConfig)
					legacyData.SetId("1")
					legacyState := legacyData.State()
					var legacyNetworkConfigurationKey string
					for attribute := range legacyState.Attributes {
						if strings.HasSuffix(attribute, ".geography") {
							legacyNetworkConfigurationKey = strings.TrimSuffix(attribute, ".geography")
							break
						}
					}
					require.NotEmpty(t, legacyNetworkConfigurationKey)
					assert.NotEqual(t, "network_configuration.0", legacyNetworkConfigurationKey)
					if test.omitEnableForAllSANsFromState {
						for attribute := range legacyState.Attributes {
							if strings.HasSuffix(attribute, ".enable_for_all_sans") {
								delete(legacyState.Attributes, attribute)
							}
						}
					}

					upgradedData := resource.Data(legacyState)
					enrollment := getTestDVEnrollment()
					enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
						CloneDNSNames: test.stateEnableForAllSANs,
					}
					attrs, err := readAttrs(&enrollment, upgradedData)
					require.NoError(t, err)
					networkConfig := attrs["network_configuration"].([]interface{})[0].(map[string]interface{})
					assert.Equal(t, test.stateEnableForAllSANs, networkConfig["enable_for_all_sans"])
					require.NoError(t, upgradedData.Set("network_configuration", attrs["network_configuration"]))

					diff, err := resource.Diff(
						context.Background(),
						upgradedData.State(),
						terraform.NewResourceConfigRaw(test.config),
						nil,
					)
					require.NoError(t, err)
					assert.Equal(t, test.wantDiff, !diff.Empty(), "unexpected diff: %#v", diff)
				})
			}
		})
	}
}

func hashLegacyNetworkConfiguration(v interface{}) int {
	networkConfig, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}

	enableForAllSANs := networkConfig["enable_for_all_sans"] == true
	hashInput := make(map[string]interface{}, len(networkConfig))
	for key, value := range networkConfig {
		switch key {
		case "enable_for_all_sans":
			continue
		case "dns_names":
			if enableForAllSANs {
				continue
			}
		}
		hashInput[key] = value
	}
	return schema.HashResource(networkConfiguration)(hashInput)
}

func dnsNamesFromState(attrs map[string]string) []string {
	var names []string
	for key, value := range attrs {
		if strings.HasPrefix(key, "network_configuration.0.dns_names.") && key != "network_configuration.0.dns_names.#" {
			names = append(names, value)
		}
	}
	return names
}

func TestConvertWarnings(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		warnings         string
		expectedWarnings []string
		err              string
	}{
		"no warnings": {
			warnings:         ``,
			expectedWarnings: nil,
		},
		"all warnings can be converted": {
			warnings:         "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			expectedWarnings: []string{"CSR_EXPIRED", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
		},
		"warnings with new line": {
			warnings:         "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nExtra certificates were found in the chain and are being removed.\ntrustChainData",
			expectedWarnings: []string{"CSR_EXPIRED", "EXTRA_CERT_IN_TRUST_CHAIN"},
		},
		"warnings with multiline": {
			warnings:         "Expected to find trust chain:\n    <expectedName>\n    <expectedDescription>\n  Instead found:\n    <actualName>\n    <actualDescription>\nExtra certificates were found in the chain and are being removed.\ntrustChainData",
			expectedWarnings: []string{"NAMED_TRUST_CHAIN_MISMATCH", "EXTRA_CERT_IN_TRUST_CHAIN"},
		},
		"warning text which is matching to two keys": {
			warnings:         "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe trust chain terminates with a non-standard root certificate.\ntrustChainData",
			expectedWarnings: []string{"CSR_EXPIRED", "TRUST_CHAIN_TERMINATES_WITH_NON_STANDARD_CERTIFICATE_DETAILED"},
		},
		"warning text which is matching to two keys 2": {
			warnings:         "Trust chain is empty.\nCertificate has a null issuer",
			expectedWarnings: []string{"TRUST_CHAIN_NULL_OR_EMPTY", "CERTIFICATE_HAS_NULL_ISSUER"},
		},
		"unknown warnings": {
			warnings: "unknown 1.\nThe key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nunknown 2.",
			err:      `received warning(s) does not match any known warning: 'unknown 1.', 'unknown 2.'`,
		},
		"warning with several newlines": {
			warnings:         "Crossed signed roots found in the trust chain.\nsomepartoftrustchain\nsecondpartoftrustchain\nendoftrust chain\nDNS Text Parse Exception when processing this \nis unknown\n that is similar",
			expectedWarnings: []string{"CROSS_SIGNED_ROOT_IN_TRUST_CHAIN", "DNS_TEXT_PARSE"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := convertWarnings(test.warnings)
			if err != nil {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Equal(t, test.expectedWarnings, res)
			}
		})
	}
}

func TestCanApproveWarnings(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		warnings            string
		autoApproveWarnings []string
		can                 bool
		err                 string
	}{
		"all warnings can be approved": {
			warnings:            "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{"CSR_EXPIRED", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			can:                 true,
		},
		"no warnings returned": {
			warnings:            "",
			autoApproveWarnings: []string{"CSR_EXPIRED", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			can:                 true,
		},
		"none can be auto-approved": {
			warnings:            "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{},
			err:                 `warnings cannot be approved: "CSR_EXPIRED", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"`,
			can:                 false,
		},
		"none warning can be auto-approved and none provided": {
			warnings:            "",
			autoApproveWarnings: []string{},
			err:                 "",
			can:                 true,
		},
		"not all warnings can be approved": {
			warnings:            "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{"CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			err:                 `warnings cannot be approved: "CSR_EXPIRED"`,
			can:                 false,
		},
		"unknown warning": {
			warnings: "unknown 1.\nThe key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nunknown 2.",
			err:      `received warning(s) does not match any known warning: 'unknown 1.', 'unknown 2.'`,
			can:      false,
		},
		"user pre-approved legacy warning": {
			warnings:            "The key for 'RSA' certificate has expired. You need to create and submit a new certificate.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "CSR_EXPIRED", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			can:                 true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := canApproveWarnings(test.autoApproveWarnings, test.warnings)
			if err != nil {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Equal(t, test.can, res)
			}
		})
	}
}

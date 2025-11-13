package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPropertyCCM(t *testing.T) {
	commonPropertyAttrs := test.AttributeBatch{
		"id":                  "prp_222222",
		"name":                "test_property",
		"contract_id":         "ctr_C-0N7RAC7",
		"group_id":            "grp_12345",
		"product_id":          "prd_Object_Delivery",
		"rule_warnings.#":     "0",
		"latest_version":      "1",
		"staging_version":     "0",
		"production_version":  "0",
		"rule_format":         "v2024-02-12",
		"use_hostname_bucket": "false",
		"rules":               `{"rules":{"name":"default","options":{}}}`,
	}
	defaultChecker := test.NewStateChecker("akamai_property.test").
		CheckEqualBatch("", commonPropertyAttrs).
		CheckEqual("hostnames.#", "1").
		CheckEqual("hostnames.0.cname_from", "example.com").
		CheckEqual("hostnames.0.cname_to", "example.com.edgekey.net").
		CheckEqual("hostnames.0.edge_hostname_id", "ehn_111").
		CheckEqual("hostnames.0.cert_provisioning_type", "CCM").
		CheckEqual("hostnames.0.cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.0.ccm_certificates.0.rsa_cert_id", "654321").
		CheckEqual("hostnames.0.ccm_certificates.0.ecdsa_cert_id", "").
		CheckEqual("hostnames.0.ccm_cert_status.0.rsa_staging_status", "NEEDS_ACTIVATION").
		CheckEqual("hostnames.0.ccm_cert_status.0.rsa_production_status", "NEEDS_ACTIVATION").
		CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_staging_status", "NOT_FOUND").
		CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_production_status", "NOT_FOUND")

	defaultImportChecker := test.NewImportChecker().
		CheckEqual("id", "prp_222222").
		CheckEqual("name", "test_property").
		CheckEqual("contract_id", "ctr_C-0N7RAC7").
		CheckEqual("group_id", "grp_12345").
		CheckEqual("product_id", "prd_Object_Delivery").
		CheckEqual("rule_warnings.#", "0").
		CheckEqual("latest_version", "1").
		CheckEqual("staging_version", "0").
		CheckEqual("production_version", "0").
		CheckEqual("rule_format", "v2024-02-12").
		CheckEqual("use_hostname_bucket", "false").
		CheckEqual("rules", `{"rules":{"name":"default","options":{}}}`).
		CheckEqual("hostnames.#", "1").
		CheckEqual("hostnames.0.cname_from", "example.com").
		CheckEqual("hostnames.0.cname_to", "example.com.edgekey.net").
		CheckEqual("hostnames.0.edge_hostname_id", "ehn_111").
		CheckEqual("hostnames.0.cert_provisioning_type", "CCM").
		CheckEqual("hostnames.0.cname_type", "EDGE_HOSTNAME").
		CheckEqual("hostnames.0.ccm_certificates.0.rsa_cert_id", "654321").
		CheckEqual("hostnames.0.ccm_certificates.0.ecdsa_cert_id", "").
		CheckEqual("hostnames.0.ccm_cert_status.0.rsa_staging_status", "NEEDS_ACTIVATION").
		CheckEqual("hostnames.0.ccm_cert_status.0.rsa_production_status", "NEEDS_ACTIVATION").
		CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_staging_status", "NOT_FOUND").
		CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_production_status", "NOT_FOUND")

	basicHostnames := func() papi.HostnameResponseItems {
		return papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_111",
					CnameFrom:            "example.com",
					CnameTo:              "example.com.edgekey.net",
					CertProvisioningType: "CCM",
					CCMCertificates: &papi.CCMCertificates{
						RSACertID:   "654321",
						RSACertLink: "/ccm/v1/certificates/654321",
					},
					CCMCertStatus: &papi.CCMCertStatus{
						ECDSAProductionStatus: "NOT_FOUND",
						ECDSAStagingStatus:    "NOT_FOUND",
						RSAProductionStatus:   "NEEDS_ACTIVATION",
						RSAStagingStatus:      "NEEDS_ACTIVATION",
					},
				},
			},
		}
	}

	basicData := func() mockPropertyData {
		return mockPropertyData{
			propertyName:  "test_property",
			productID:     "prd_Object_Delivery",
			propertyID:    "prp_222222",
			groupID:       "grp_12345",
			contractID:    "ctr_C-0N7RAC7",
			assetID:       "aid_5555",
			latestVersion: 1,
			ruleTree: mockRuleTreeData{
				rules: papi.Rules{
					Name: "default",
				},
			},
			hostnames: basicHostnames(),
			versions: papi.PropertyVersionItems{
				Items: []papi.PropertyVersionGetItem{
					{
						PropertyVersion:  1,
						StagingStatus:    "INACTIVE",
						ProductionStatus: "INACTIVE",
					},
				},
			},
		}
	}

	tests := map[string]struct {
		init  func(*mockProperty)
		steps []resource.TestStep
	}{
		"Creating basic property with CCM RSA certificate": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property.tf"),
					Check: defaultChecker.
						Build(),
				},
			},
		},
		"Creating basic property with CCM ECDSA certificate": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				p.hostnames.Items[0].CCMCertificates = &papi.CCMCertificates{
					ECDSACertID:   "765432",
					ECDSACertLink: "/ccm/v1/certificates/765432",
				}
				p.hostnames.Items[0].CCMCertStatus = &papi.CCMCertStatus{
					ECDSAProductionStatus: "NEEDS_ACTIVATION",
					ECDSAStagingStatus:    "NEEDS_ACTIVATION",
					RSAProductionStatus:   "NOT_FOUND",
					RSAStagingStatus:      "NOT_FOUND",
				}

				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property_ecdsa.tf"),
					Check: defaultChecker.
						CheckEqual("hostnames.0.ccm_certificates.0.rsa_cert_id", "").
						CheckEqual("hostnames.0.ccm_certificates.0.ecdsa_cert_id", "765432").
						CheckEqual("hostnames.0.ccm_cert_status.0.rsa_staging_status", "NOT_FOUND").
						CheckEqual("hostnames.0.ccm_cert_status.0.rsa_production_status", "NOT_FOUND").
						CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_staging_status", "NEEDS_ACTIVATION").
						CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_production_status", "NEEDS_ACTIVATION").
						Build(),
				},
			},
		},
		"Creating basic property with both CCM RSA and ECDSA certificates": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				p.hostnames.Items[0].CCMCertificates = &papi.CCMCertificates{
					ECDSACertID:   "765432",
					ECDSACertLink: "/ccm/v1/certificates/765432",
					RSACertID:     "654321",
					RSACertLink:   "/ccm/v1/certificates/654321",
				}
				p.hostnames.Items[0].CCMCertStatus = &papi.CCMCertStatus{
					ECDSAProductionStatus: "NEEDS_ACTIVATION",
					ECDSAStagingStatus:    "NEEDS_ACTIVATION",
					RSAProductionStatus:   "NEEDS_ACTIVATION",
					RSAStagingStatus:      "NEEDS_ACTIVATION",
				}

				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property_rsa_ecdsa.tf"),
					Check: defaultChecker.
						CheckEqual("hostnames.0.ccm_certificates.0.ecdsa_cert_id", "765432").
						CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_staging_status", "NEEDS_ACTIVATION").
						CheckEqual("hostnames.0.ccm_cert_status.0.ecdsa_production_status", "NEEDS_ACTIVATION").
						Build(),
				},
			},
		},
		"Updating ID of the CCM RSA certificate on inactive property - no new version": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read x 2
				mockResourcePropertyRead(p, 2)
				// update
				p.mockGetPropertyVersion()
				p.hostnames = basicHostnames()
				p.hostnames.Items[0].CCMCertificates.RSACertID = "987654"
				p.hostnames.Items[0].CCMCertificates.RSACertLink = "/ccm/v1/certificates/987654"
				p.mockUpdatePropertyVersionHostnames()
				// read from update
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property.tf"),
					Check: defaultChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property_update_rsa_cert_id.tf"),
					Check: defaultChecker.
						CheckEqual("hostnames.0.ccm_certificates.0.rsa_cert_id", "987654").
						Build(),
				},
			},
		},
		"Updating ID of the CCM RSA certificate on active property - new version": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				p.versions.Items[0].StagingStatus = "ACTIVE"
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read x 2
				mockResourcePropertyRead(p, 2)
				// update
				p.mockGetPropertyVersion()
				p.createFromVersion = 1
				p.newVersionID = 2
				p.mockCreatePropertyVersion()
				p.hostnames = basicHostnames()
				p.hostnames.Items[0].CCMCertificates.RSACertID = "987654"
				p.hostnames.Items[0].CCMCertificates.RSACertLink = "/ccm/v1/certificates/987654"
				p.latestVersion = 2
				p.mockUpdatePropertyVersionHostnames()
				// read from update
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property.tf"),
					Check: defaultChecker.
						CheckEqual("staging_version", "1").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property_update_rsa_cert_id.tf"),
					Check: defaultChecker.
						CheckEqual("staging_version", "1").
						CheckEqual("latest_version", "2").
						CheckEqual("hostnames.0.ccm_certificates.0.rsa_cert_id", "987654").
						Build(),
				},
			},
		},
		"Adding CCM-bound hostname to active property with default hostname - new version": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_222",
							CnameFrom:            "test.com",
							CnameTo:              "test.com.edgekey.net",
							CertProvisioningType: "DEFAULT",
						},
					},
				}

				p.versions.Items[0].StagingStatus = "ACTIVE"
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read from create
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p)
				// read x 2
				mockResourcePropertyRead(p, 2)
				// update
				p.mockGetPropertyVersion()
				p.createFromVersion = 1
				p.newVersionID = 2
				p.mockCreatePropertyVersion()
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_111",
							CnameFrom:            "example.com",
							CnameTo:              "example.com.edgekey.net",
							CertProvisioningType: "CCM",
							CCMCertificates: &papi.CCMCertificates{
								RSACertID:   "654321",
								RSACertLink: "/ccm/v1/certificates/654321",
							},
							CCMCertStatus: &papi.CCMCertStatus{
								ECDSAProductionStatus: "NOT_FOUND",
								ECDSAStagingStatus:    "NOT_FOUND",
								RSAProductionStatus:   "NEEDS_ACTIVATION",
								RSAStagingStatus:      "NEEDS_ACTIVATION",
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_222",
							CnameFrom:            "test.com",
							CnameTo:              "test.com.edgekey.net",
							CertProvisioningType: "DEFAULT",
						},
					},
				}
				p.latestVersion = 2
				p.mockUpdatePropertyVersionHostnames()
				// read from update
				mockResourcePropertyRead(p)
				// read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/default_hostname.tf"),
					Check: test.NewStateChecker("akamai_property.test").
						CheckEqualBatch("", commonPropertyAttrs).
						CheckEqual("hostnames.#", "1").
						CheckEqual("hostnames.0.cname_from", "test.com").
						CheckEqual("hostnames.0.cname_to", "test.com.edgekey.net").
						CheckEqual("hostnames.0.edge_hostname_id", "ehn_222").
						CheckEqual("hostnames.0.cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.0.cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.ccm_certificates.#", "0").
						CheckEqual("hostnames.0.ccm_cert_status.#", "0").
						CheckEqual("staging_version", "1").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/default_and_ccm_hostnames.tf"),
					Check: defaultChecker.
						CheckEqual("hostnames.#", "2").
						CheckEqual("hostnames.1.cname_from", "test.com").
						CheckEqual("hostnames.1.cname_to", "test.com.edgekey.net").
						CheckEqual("hostnames.1.edge_hostname_id", "ehn_222").
						CheckEqual("hostnames.1.cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.1.cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.1.ccm_certificates.#", "0").
						CheckEqual("hostnames.1.ccm_cert_status.#", "0").
						CheckEqual("staging_version", "1").
						CheckEqual("latest_version", "2").
						Build(),
				},
			},
		},
		"Importing basic property with CCM RSA certificate": {
			init: func(p *mockProperty) {
				p.mockPropertyData = basicData()
				// read
				p.ruleTree.ruleFormat = "v2024-02-12"
				mockResourcePropertyRead(p, 2)

				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck:   defaultImportChecker.Build(),
					ImportStateId:      "prp_222222,ctr_C-0N7RAC7,grp_12345",
					ImportState:        true,
					ResourceName:       "akamai_property.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property.tf"),
					ImportStatePersist: true,
				},
				{
					// Confirm idempotency after import
					Config:   testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_property.tf"),
					PlanOnly: true,
				},
			},
		},
		"Error no certificates for CCM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_no_cert_section.tf"),
					ExpectError: regexp.MustCompile(`ccm_certificates is required when cert_provisioning_type is 'CCM'`),
				},
			},
		},
		"Error certificates specified for no CCM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/cert_section_no_ccm.tf"),
					ExpectError: regexp.MustCompile(`ccm_certificates is only allowed when cert_provisioning_type is 'CCM'`),
				},
			},
		},
		"Error empty certificates for CCM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_empty_cert_section.tf"),
					ExpectError: regexp.MustCompile(`hostname from.test.domain: at least one of rsa_cert_id or ecdsa_cert_id must be provided in ccm_certificates`),
				},
			},
		},
		"Error empty RSA cert id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_empty_rsa_cert_id.tf"),
					ExpectError: regexp.MustCompile(`hostname from.test.domain: at least one of rsa_cert_id or ecdsa_cert_id must be provided in ccm_certificates`),
				},
			},
		},
		"Error two CCM certificate sections": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_two_cert_sections.tf"),
					ExpectError: regexp.MustCompile(`Too many ccm_certificates blocks`),
				},
			},
		},
		"Error CCM does not support hostname buckets": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CCM/ccm_with_hostname_bucket.tf"),
					ExpectError: regexp.MustCompile(`hostnames should be empty for use_hostname_bucket enabled`),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			if test.init != nil {
				test.init(&mp)
			}

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			papiMock.AssertExpectations(t)
		})
	}
}

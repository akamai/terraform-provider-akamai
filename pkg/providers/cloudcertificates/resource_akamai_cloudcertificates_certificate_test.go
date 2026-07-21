package cloudcertificates

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v10/internal/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	certificateTestData struct {
		// input data
		contractID    string
		groupID       string
		baseName      string
		keyType       cloudcertificates.CryptographicAlgorithm
		keySize       cloudcertificates.KeySize
		secureNetwork cloudcertificates.SecureNetwork
		inputGeoClass cloudcertificates.GeoClass
		sans          []string
		subject       *cloudcertificates.Subject

		// output data
		certificateID           string
		certificateType         string
		name                    string
		certificateStatus       string
		accountID               string
		createdBy               string
		createdDate             string
		modifiedBy              string
		modifiedDate            string
		csrExpirationDate       string
		csrPEM                  string
		signedCertNotValidAfter *time.Time
		outputGeoClass          cloudcertificates.GeoClass
	}
)

var (
	minCertificate = certificateTestData{
		// input data
		contractID:    "test_contract",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	minCertificateWithPrefixes = certificateTestData{
		// input data
		contractID:    "ctr_test_contract",
		groupID:       "grp_123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	fullCertificateRSA = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "STANDARD_WORLDWIDE",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &cloudcertificates.Subject{
			CommonName:   "test.example.com",
			Country:      "US",
			Organization: "Test Org",
			State:        "CA",
			Locality:     "Test City",
		},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test-name",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	updateCertificate = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name-updated",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &cloudcertificates.Subject{
			CommonName:   "test.example.com",
			Country:      "US",
			Organization: "Test Org",
			State:        "CA",
			Locality:     "Test City",
		},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test-name-updated",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user-updated",
		modifiedDate:      "2025-05-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	fullCertificateECDSA = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name",
		groupID:       "123",
		keyType:       "ECDSA",
		keySize:       "P-256",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "CONTIGUOUS_US",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &cloudcertificates.Subject{
			State:    "CA",
			Locality: "Test City",
		},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test-name",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "CONTIGUOUS_US",
	}

	minCertificateECDSAP384 = certificateTestData{
		// input data
		contractID:    "test_contract",
		groupID:       "123",
		keyType:       "ECDSA",
		keySize:       "P-384",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	minCertificateStandardTLS = certificateTestData{
		// input data
		contractID:    "test_contract",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "STANDARD_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "STANDARD_WORLDWIDE",
	}

	minCertificateNonDefaultGeoClass = certificateTestData{
		// input data
		contractID:    "test_contract",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "CONTIGUOUS_US",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "CONTIGUOUS_US",
	}

	minCertificateEmptyGeoClassFromAPI = certificateTestData{
		// input data
		contractID:    "test_contract",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		inputGeoClass: "",
		sans:          []string{"test.example.com"},

		// output data
		certificateID:     "12345",
		certificateType:   "THIRD_PARTY",
		name:              "test.example.com1234567890",
		certificateStatus: "CSR_READY",
		accountID:         "act_789",
		createdBy:         "test_user",
		createdDate:       "2025-01-01T00:00:00.168262Z",
		modifiedBy:        "test_user",
		modifiedDate:      "2025-01-01T00:00:00.616267Z",
		csrExpirationDate: "2027-01-01T00:00:00Z",
		csrPEM:            "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n",
		outputGeoClass:    "",
	}
)

func TestCertificateResource(t *testing.T) {
	t.Parallel()
	config := defaultSubproviderConfig()
	config.certificate.timestampFunc = func() time.Time {
		t, _ := time.Parse(renewedNameDateLayout, "2025-05-01T12_05_01Z")
		return t
	}

	minCertChecker := test.NewStateChecker("akamai_cloudcertificates_certificate.test").
		CheckEqual("contract_id", "test_contract").
		CheckEqual("group_id", "123").
		CheckEqual("key_type", "RSA").
		CheckEqual("key_size", "2048").
		CheckEqual("secure_network", "ENHANCED_TLS").
		CheckEqual("geo_class", "STANDARD_WORLDWIDE").
		CheckEqual("sans.#", "1").
		CheckEqual("sans.0", "test.example.com").
		CheckEqual("certificate_id", "12345").
		CheckEqual("certificate_type", "THIRD_PARTY").
		CheckEqual("certificate_status", "CSR_READY").
		CheckEqual("name", "test.example.com1234567890").
		CheckEqual("account_id", "act_789").
		CheckEqual("created_by", "test_user").
		CheckEqual("created_date", "2025-01-01T00:00:00.168262Z").
		CheckEqual("modified_by", "test_user").
		CheckEqual("modified_date", "2025-01-01T00:00:00.616267Z").
		CheckEqual("csr_expiration_date", "2027-01-01T00:00:00Z").
		CheckEqual("csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n").
		CheckMissing("base_name")

	fullCertChecker := minCertChecker.
		CheckEqual("base_name", "test-name").
		CheckEqual("name", "test-name").
		CheckEqual("sans.#", "2").
		CheckEqual("sans.0", "test.example.com").
		CheckEqual("sans.1", "test.example2.com").
		CheckEqual("subject.common_name", "test.example.com").
		CheckEqual("subject.country", "US").
		CheckEqual("subject.organization", "Test Org").
		CheckEqual("subject.state", "CA").
		CheckEqual("subject.locality", "Test City")

	importChecker := test.NewImportChecker().
		CheckEqual("contract_id", "test_contract").
		CheckEqual("key_type", "RSA").
		CheckEqual("key_size", "2048").
		CheckEqual("secure_network", "ENHANCED_TLS").
		CheckEqual("geo_class", "STANDARD_WORLDWIDE").
		CheckEqual("sans.#", "1").
		CheckEqual("sans.0", "test.example.com").
		CheckEqual("certificate_id", "12345").
		CheckEqual("certificate_type", "THIRD_PARTY").
		CheckEqual("certificate_status", "CSR_READY").
		CheckEqual("name", "test.example.com1234567890").
		CheckEqual("base_name", "test.example.com1234567890").
		CheckEqual("account_id", "act_789").
		CheckEqual("created_by", "test_user").
		CheckEqual("created_date", "2025-01-01T00:00:00.168262Z").
		CheckEqual("modified_by", "test_user").
		CheckEqual("modified_date", "2025-01-01T00:00:00.616267Z").
		CheckEqual("csr_expiration_date", "2027-01-01T00:00:00Z").
		CheckEqual("csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n").
		CheckEqual("renew_pending", "false").
		CheckEqual("auto_renew", "false")

	tests := map[string]struct {
		init           func(*cloudcertificates.Mock, certificateTestData, certificateTestData)
		createMockData certificateTestData
		updateMockData certificateTestData
		steps          []resource.TestStep
	}{
		"happy path - create certificate without optionals": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check: minCertChecker.
						CheckMissing("subject").
						Build(),
				},
			},
		},
		"happy path - create certificate with geo_class not returned by the API": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificateEmptyGeoClassFromAPI,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check: minCertChecker.
						CheckMissing("subject").
						CheckMissing("geo_class").
						Build(),
				},
			},
		},
		"happy path - create certificate with prefixes and without optionals": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificateWithPrefixes,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min_with_prefixes.tf"),
					Check: minCertChecker.
						CheckEqual("contract_id", "ctr_test_contract").
						CheckEqual("group_id", "grp_123").
						CheckMissing("subject").
						Build(),
				},
			},
		},
		"happy path - create certificate with STANDARD_TLS": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificateStandardTLS,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min_standard_tls.tf"),
					Check: minCertChecker.
						CheckEqual("secure_network", "STANDARD_TLS").
						CheckMissing("subject").
						Build(),
				},
			},
		},
		"happy path - create certificate with all optional attributes": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
			},
		},
		"happy path - create certificate with non-default geo_class": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificateNonDefaultGeoClass,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min_contiguous_us.tf"),
					Check: minCertChecker.
						CheckEqual("geo_class", "CONTIGUOUS_US").
						CheckMissing("subject").
						Build(),
				},
			},
		},
		"happy path - default geo_class, no drift on subsequent plan": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read (PlanOnly refresh + destroy)
				mockGetCertificate(m, createData).Twice()
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check: minCertChecker.
						CheckEqual("geo_class", "STANDARD_WORLDWIDE").
						CheckMissing("subject").
						Build(),
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - create certificate with optional attributes, different key type, some missing subject fields": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateECDSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full_different_values.tf"),
					Check: fullCertChecker.
						CheckEqual("key_size", "P-256").
						CheckEqual("key_type", "ECDSA").
						CheckEqual("geo_class", "CONTIGUOUS_US").
						CheckMissing("subject.common_name").
						CheckMissing("subject.country").
						CheckMissing("subject.organization").
						Build(),
				},
			},
		},
		"happy path - create certificate, update name": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Update - check renewal chain for new base_name
				mockEmptyRenewalChain(m, updateData)
				// Update
				mockPatchCertificate(m, updateData)
				// Read after update
				mockGetCertificate(m, updateData)
				// Read before destroy
				mockGetCertificate(m, updateData)
				// Delete
				mockDeleteCertificate(m, updateData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/name.tf"),
					Check: fullCertChecker.
						CheckEqual("base_name", "test-name-updated").
						CheckEqual("name", "test-name-updated").
						CheckEqual("modified_date", "2025-05-01T00:00:00.616267Z").
						CheckEqual("modified_by", "test_user-updated").
						Build(),
				},
			},
		},
		"happy path - create certificate, reset name": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				updateData.baseName = "" // empty string resets the name to the default value.
				updateData.name = "generated-name12345"
				// Update
				mockPatchCertificate(m, updateData)
				// Read after update
				mockGetCertificate(m, updateData)
				// Read before destroy
				mockGetCertificate(m, updateData)
				// Delete
				mockDeleteCertificate(m, updateData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/reset_name.tf"),
					Check: fullCertChecker.
						CheckMissing("base_name").
						CheckEqual("name", "generated-name12345").
						CheckEqual("modified_date", "2025-05-01T00:00:00.616267Z").
						CheckEqual("modified_by", "test_user-updated").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectUnknownValue("akamai_cloudcertificates_certificate.test", tfjsonpath.New("modified_date")),
							plancheck.ExpectUnknownValue("akamai_cloudcertificates_certificate.test", tfjsonpath.New("modified_by")),
						},
					},
				},
			},
		},
		"happy path - create certificate, update name with renewal chain": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Update - check renewal chain for new base_name (chain exists)
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      updateData.contractID,
					Domain:          updateData.sans[0],
					CertificateName: updateData.baseName,
				}, &cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{
						{CertificateName: updateData.baseName},
						{CertificateName: "test-name-updated.renewed.2025-03-01T12_05_01Z"},
					},
				}, nil).Once()
				// Update - patch with suffixed name
				suffixedName := "test-name-updated.renewed.2025-05-01T12_05_01Z"
				patchData := updateData
				patchData.baseName = suffixedName
				patchData.name = suffixedName
				mockPatchCertificate(m, patchData)
				// Read after update
				mockGetCertificate(m, patchData)
				// Read before destroy
				mockGetCertificate(m, patchData)
				// Delete
				mockDeleteCertificate(m, patchData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/name.tf"),
					Check: fullCertChecker.
						CheckEqual("base_name", "test-name-updated").
						CheckEqual("name", "test-name-updated.renewed.2025-05-01T12_05_01Z").
						CheckEqual("modified_date", "2025-05-01T00:00:00.616267Z").
						CheckEqual("modified_by", "test_user-updated").
						Build(),
				},
			},
		},
		"happy path - create certificate, change order of SANs - no diff": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read x2
				mockGetCertificate(m, createData).Twice()
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full_different_sans_order.tf"),
					Check:  fullCertChecker.Build(),
				},
			},
		},
		"happy path - create certificate, remove outside terraform": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before refresh
				mockGetCertificate(m, createData)
				// Read from refresh step
				m.On("GetCertificate", testutils.MockContext, cloudcertificates.GetCertificateRequest{
					CertificateID: createData.certificateID,
				}).Return(nil, cloudcertificates.ErrCertificateNotFound).Once()
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check: minCertChecker.
						CheckMissing("subject").
						Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - renew certificate": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Certificate exists already
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      createData.contractID,
					Domain:          createData.sans[0],
					CertificateName: createData.baseName,
				}, &cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{
						{
							CertificateName: createData.name,
						},
						{
							CertificateName: "test-name.renewed.2025-03-01T12_05_01Z",
						},
						{
							CertificateName: "test-name.renewed.2025-04-01T12_05_01Z",
						},
					},
				}, nil).Once()
				// Renew
				renewedName := fmt.Sprintf("%s.renewed.2025-05-01T12_05_01Z", createData.name)
				renewedCertificate := createData
				renewedCertificate.name = renewedName
				renewedCertificate.baseName = renewedName
				renewedCertificate.certificateID = "123456"
				renewedCertificate.keyType = "ECDSA"
				renewedCertificate.keySize = "P-256"
				mockCreateCertificate(m, renewedCertificate)
				// Read before destroy
				mockGetCertificate(m, renewedCertificate)
				// Delete
				mockDeleteCertificate(m, renewedCertificate)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/renew_full.tf"),
					Check: fullCertChecker.
						CheckEqual("name", "test-name.renewed.2025-05-01T12_05_01Z").
						CheckEqual("certificate_id", "123456").
						CheckEqual("key_type", "ECDSA").
						CheckEqual("key_size", "P-256").
						Build(),
				},
			},
		},
		"happy path - create certificate without optionals and renew certificate": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before renew
				mockGetCertificate(m, createData).Twice()
				// Renew
				renewedName := "test.example.com2345678901"
				renewedCertificate := createData
				renewedCertificate.name = renewedName
				renewedCertificate.certificateID = "123456"
				renewedCertificate.keyType = "ECDSA"
				renewedCertificate.keySize = "P-256"
				// Create new certificate BEFORE destroying the old one
				mockCreateCertificate(m, renewedCertificate)
				// Delete old certificate
				mockDeleteCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, renewedCertificate)
				// Delete
				mockDeleteCertificate(m, renewedCertificate)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min_with_create_before_destroy.tf"),
					Check: minCertChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/renew_min.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("akamai_cloudcertificates_certificate.test",
								plancheck.ResourceActionCreateBeforeDestroy),
						},
					},
					Check: minCertChecker.
						CheckEqual("name", "test.example.com2345678901").
						CheckEqual("certificate_id", "123456").
						CheckEqual("key_type", "ECDSA").
						CheckEqual("key_size", "P-256").
						Build(),
				},
			},
		},
		"happy path - create certificate and renew certificate": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      createData.contractID,
					Domain:          createData.sans[0],
					CertificateName: createData.baseName,
				}, &cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{{}},
				}, nil).Once()
				mockCreateCertificate(m, createData)
				// Read before renew
				mockGetCertificate(m, createData).Twice()
				// Renew
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      createData.contractID,
					Domain:          createData.sans[0],
					CertificateName: createData.baseName,
				}, &cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{
						{
							CertificateName: createData.name,
						},
					},
				}, nil).Once()

				renewedName := fmt.Sprintf("%s.renewed.2025-05-01T12_05_01Z", createData.name)
				renewedCertificate := createData
				renewedCertificate.name = renewedName
				renewedCertificate.baseName = renewedName
				renewedCertificate.certificateID = "123456"
				renewedCertificate.keyType = "ECDSA"
				renewedCertificate.keySize = "P-256"
				// Create new certificate BEFORE destroying the old one
				mockCreateCertificate(m, renewedCertificate)
				// Delete old certificate
				mockDeleteCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, renewedCertificate)
				// Delete
				mockDeleteCertificate(m, renewedCertificate)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full_with_create_before_destroy.tf"),
					Check: fullCertChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/renew_full.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("akamai_cloudcertificates_certificate.test",
								plancheck.ResourceActionCreateBeforeDestroy),
						},
					},
					Check: fullCertChecker.
						CheckEqual("name", "test-name.renewed.2025-05-01T12_05_01Z").
						CheckEqual("certificate_id", "123456").
						CheckEqual("key_type", "ECDSA").
						CheckEqual("key_size", "P-256").
						Build(),
				},
			},
		},
		"happy path - create certificate with ECDSA P-384": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificateECDSAP384,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min_ecdsa_p384.tf"),
					Check: minCertChecker.
						CheckEqual("key_type", "ECDSA").
						CheckEqual("key_size", "P-384").
						CheckMissing("subject").
						Build(),
				},
			},
		},
		"import - not renewed certificate": {
			init: func(m *cloudcertificates.Mock, data certificateTestData, _ certificateTestData) {
				// Import
				mockGetCertificate(m, data)
				// Read
				mockGetCertificate(m, data).Times(2)
				// Delete after plan
				mockDeleteCertificate(m, data)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					ImportStateCheck:   importChecker.Build(),
					ImportStateId:      "12345",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_certificate.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/no_group_id.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/no_group_id.tf"),
					PlanOnly: true,
				},
			},
		},
		"import - renewed certificate": {
			init: func(m *cloudcertificates.Mock, data certificateTestData, _ certificateTestData) {
				data.name = "test-certificate.renewed.2025-05-01T12_05_01Z"
				// Import
				mockGetCertificate(m, data)
				// Read
				mockGetCertificate(m, data).Times(2)
				// Delete after plan
				mockDeleteCertificate(m, data)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("name", "test-certificate.renewed.2025-05-01T12_05_01Z").
						CheckEqual("base_name", "test-certificate").
						Build(),
					ImportStateId:      "12345",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_certificate.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/no_group_id_base_name.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/no_group_id_base_name.tf"),
					PlanOnly: true,
				},
			},
		},
		"import - with optional group_id": {
			init: func(m *cloudcertificates.Mock, data certificateTestData, _ certificateTestData) {
				// Import
				mockGetCertificate(m, data)
				// Read
				mockGetCertificate(m, data).Times(2)
				// Delete after plan
				mockDeleteCertificate(m, data)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.
						CheckEqual("group_id", "grp_123").
						Build(),
					ImportStateId:      "12345,grp_123",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_certificate.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/with_group_id.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResCertificate/import/with_group_id.tf"),
					PlanOnly: true,
				},
			},
		},
		"import - expect error - wrong ID: ErrCertificateNotFound - remove state": {
			init: func(m *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {
				// Import
				m.On("GetCertificate", testutils.MockContext, cloudcertificates.GetCertificateRequest{
					CertificateID: "12345abc-wrong",
				}).Return(nil, cloudcertificates.ErrCertificateNotFound).Once()
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					ImportStateId: "12345abc-wrong",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_certificate.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ExpectError:   regexp.MustCompile(`Error: Cannot import non-existent remote object`),
				},
			},
		},
		"import - expect error - wrong ID format": {
			steps: []resource.TestStep{
				{
					ImportStateId: "12345,grp_123,unexpected",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_certificate.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ExpectError: regexp.MustCompile(`Error: Incorrect import ID:(\n|.)+` +
						`invalid number of importID parts: 3; you need to provide an importID in the\nformat 'certificateID\[,groupID]'`),
				},
			},
		},
		"expect error - imported without group_id, but config has group_id": {
			init: func(m *cloudcertificates.Mock, data certificateTestData, _ certificateTestData) {
				data.groupID = ""
				// Import
				mockGetCertificate(m, data)
				// Read
				mockGetCertificate(m, data)
				// Delete after plan
				mockDeleteCertificate(m, data)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					ImportStateCheck:   importChecker.Build(),
					ImportStateId:      "12345",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_certificate.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ImportStatePersist: true,
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					PlanOnly: true,
					ExpectError: regexp.MustCompile(`Error: The resource was imported without a group_id(\n|.)+` +
						`To fix this, you need to first remove the state and then re-import it with\nthe group_id specified in the import ID.`),
				},
			},
		},
		"expect error - CreateCertificate fails": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				m.On("CreateCertificate", testutils.MockContext, cloudcertificates.CreateCertificateRequest{
					ContractID: createData.contractID,
					GroupID:    createData.groupID,
					Body: cloudcertificates.CreateCertificateRequestBody{
						CertificateName: createData.baseName,
						KeyType:         createData.keyType,
						KeySize:         createData.keySize,
						SecureNetwork:   createData.secureNetwork,
						GeoClass:        createData.inputGeoClass,
						SANs:            createData.sans,
					},
				}).Return(nil, fmt.Errorf("API failed")).Once()
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to create CCM Certificate(.|\n)*API failed`),
				},
			},
		},
		"expect error - GetCertificate fails": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				m.On("GetCertificate", testutils.MockContext, cloudcertificates.GetCertificateRequest{
					CertificateID: createData.certificateID,
				}).Return(nil, fmt.Errorf("API failed")).Once()
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to get CCM Certificate(.|\n)*API failed`),
				},
			},
		},
		"expect error - PatchCertificate fails": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Update - check renewal chain for new base_name
				mockEmptyRenewalChain(m, updateData)
				// Update
				m.On("PatchCertificate", testutils.MockContext, cloudcertificates.PatchCertificateRequest{
					CertificateID:   updateData.certificateID,
					CertificateName: ptr.To(updateData.baseName),
				}).Return(nil, fmt.Errorf("API failed")).Once()
				// Read before destroy
				mockGetCertificate(m, updateData)
				// Delete
				mockDeleteCertificate(m, updateData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check: minCertChecker.
						CheckEqual("base_name", "test-name").
						CheckEqual("name", "test-name").
						CheckEqual("sans.#", "2").
						CheckEqual("sans.0", "test.example.com").
						CheckEqual("sans.1", "test.example2.com").
						CheckEqual("subject.common_name", "test.example.com").
						CheckEqual("subject.country", "US").
						CheckEqual("subject.organization", "Test Org").
						CheckEqual("subject.state", "CA").
						CheckEqual("subject.locality", "Test City").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/name.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to update CCM Certificate(.|\n)*API failed`),
				},
			},
		},
		"expect error - ListCertificates fails during update": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Update - ListCertificates fails
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      updateData.contractID,
					Domain:          updateData.sans[0],
					CertificateName: updateData.baseName,
				}, nil, fmt.Errorf("API failed")).Once()
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			updateMockData: updateCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check:  fullCertChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/name.tf"),
					ExpectError: regexp.MustCompile(`Error: Unable to verify CCM Certificate name(.|\n)*API failed`),
				},
			},
		},
		"expect error - missing contract": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_contract.tf"),
					ExpectError: regexp.MustCompile("The argument \"contract_id\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing group": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_group.tf"),
					ExpectError: regexp.MustCompile(`Error: Required Field Missing(.|\n)*field ` + "`group_id`" + ` is required during creation`),
				},
			},
		},
		"expect error - missing key_size": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_key_size.tf"),
					ExpectError: regexp.MustCompile("The argument \"key_size\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing key_type": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_key_type.tf"),
					ExpectError: regexp.MustCompile("The argument \"key_type\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing secure_network": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_secure_network.tf"),
					ExpectError: regexp.MustCompile("The argument \"secure_network\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing sans": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_sans.tf"),
					ExpectError: regexp.MustCompile("The argument \"sans\" is required, but no definition was found."),
				},
			},
		},
		"expect error - empty sans": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_sans.tf"),
					ExpectError: regexp.MustCompile("Attribute sans set must contain at least 1 elements and at most 100 elements,\ngot: 0"),
				},
			},
		},
		"expect error - more than 100 sans provided": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/more_than_100_sans.tf"),
					ExpectError: regexp.MustCompile("Attribute sans set must contain at least 1 elements and at most 100 elements,\ngot: 101"),
				},
			},
		},
		"expect error - one of sans is not a valid domain name": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/sans_not_valid_domain.tf"),
					ExpectError: regexp.MustCompile(`Attribute sans.+"invalid_domain".+must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - one of sans contains uppercase letters": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/sans_uppercase.tf"),
					ExpectError: regexp.MustCompile(`Attribute sans.+"example.COM".+must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - empty base_name": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_base_name.tf"),
					ExpectError: regexp.MustCompile("Attribute base_name string length must be at least 1, got: 0"),
				},
			},
		},
		"expect error - wrong key_size for RSA": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_size_rsa.tf"),
					ExpectError: regexp.MustCompile(
						`Attribute key_size value must be one of: \["2048" "P-256" "P-384"\], got:\s+"2137"`),
				},
			},
		},
		"expect error - wrong key_size for ECDSA": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_size_ecdsa.tf"),
					ExpectError: regexp.MustCompile(
						`Attribute key_size value must be one of: \["2048" "P-256" "P-384"\], got:\s+"2137"`),
				},
			},
		},
		"expect error - P-384 key_size used with RSA key_type": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_size_p384_with_rsa.tf"),
					ExpectError: regexp.MustCompile(`The specified value 'P-384' for the RSA key type is invalid. Valid values(.|\n)are: '2048'.`),
				},
			},
		},
		"expect error - wrong key_type": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_type.tf"),
					ExpectError: regexp.MustCompile(`Attribute key_type value must be one of: \["ECDSA" "RSA"\], got:(.|\n)*"WRONG-TYPE"`),
				},
			},
		},
		"expect error - wrong secure_network": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_secure_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute secure_network value must be one of: \["ENHANCED_TLS"(.|\n)*"STANDARD_TLS"\], got:(.|\n)*"WRONG_NETWORK"`),
				},
			},
		},
		"expect error - wrong geo_class": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_geo_class.tf"),
					ExpectError: regexp.MustCompile(`Attribute geo_class value must be one of: \["CONTIGUOUS_US"(.|\n)*"RESERVED_GLOBAL"(.|\n)*"STANDARD_WORLDWIDE"\], got:(.|\n)*"INVALID_GEO_CLASS"`),
				},
			},
		},
		"expect error - CONTIGUOUS_US geo_class invalid for STANDARD_TLS network": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/geo_class_contiguous_us_with_standard_tls.tf"),
					ExpectError: regexp.MustCompile(`The specified value 'CONTIGUOUS_US' for the STANDARD_TLS network is(.|\n)*invalid.(.|\n)*Valid values are: 'STANDARD_WORLDWIDE'.`),
				},
			},
		},
		"expect error - RESERVED_GLOBAL geo_class invalid for STANDARD_TLS network": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/geo_class_reserved_global_with_standard_tls.tf"),
					ExpectError: regexp.MustCompile(`The specified value 'RESERVED_GLOBAL' for the STANDARD_TLS network is(.|\n)*invalid.(.|\n)*Valid values are: 'STANDARD_WORLDWIDE'.`),
				},
			},
		},
		"expect error - empty subject": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_subject.tf"),
					ExpectError: regexp.MustCompile(`At least one of the subject fields \(common_name, organization, country,(.|\n)*state, locality\) must be specified.`),
				},
			},
		},
		"expect error - common_name not present in sans": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_not_in_sans.tf"),
					ExpectError: regexp.MustCompile(`The specified common name 'test.example.com' must be included in the SANs(.|\n)*list.`),
				},
			},
		},
		"expect error - common_name is not a valid domain name": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_not_valid_domain.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject\.common_name must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - common_name contains uppercase letters": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_contains_uppercase.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject\.common_name must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - empty organization": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_organization.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.organization string length must be between 1 and 64, got: 0`),
				},
			},
		},
		"expect error - country too long": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/country_too_long.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.country string length must be between 2 and 2, got: 7`),
				},
			},
		},
		"expect error - empty state": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_state.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.state string length must be between 1 and 128, got: 0`),
				},
			},
		},
		"expect error - empty locality": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_locality.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.locality string length must be between 1 and 128, got: 0`),
				},
			},
		},
		"expect error - organization with only spaces": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/organization_only_spaces.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.organization cannot be empty or whitespace, got:  `),
				},
			},
		},
		"expect error - auto_renew without renew_before_expiration_days": {
			init:           func(_ *cloudcertificates.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/auto_renew_without_threshold.tf"),
					ExpectError: regexp.MustCompile(`auto_renew.*cannot be set to true without.*renew_before_expiration_days`),
				},
			},
		},
		"expect error - update contract": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check:  minCertChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/contract.tf"),
					ExpectError: regexp.MustCompile("updating field `contract_id` is not possible"),
				},
			},
		},
		"expect error - update group": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check:  minCertChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/group.tf"),
					ExpectError: regexp.MustCompile("updating field `group_id` is not possible"),
				},
			},
		},
		"expect error - update sans": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check:  minCertChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/sans.tf"),
					ExpectError: regexp.MustCompile("updating field `sans` is not possible"),
				},
			},
		},
		"expect error - update subject": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check: fullCertChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/subject.tf"),
					ExpectError: regexp.MustCompile("updating field `subject` is not possible"),
				},
			},
		},
		"expect error - no subject for create, but subject present in update": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					Check:  minCertChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/subject_present.tf"),
					ExpectError: regexp.MustCompile("updating field `subject` is not possible"),
				},
			},
		},
		"expect error - create with subject, update by removing subject": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				createData.subject = &cloudcertificates.Subject{
					CommonName:   "test.example.com",
					Country:      "US",
					Organization: "Test Org - updated",
					State:        "CA",
					Locality:     "Test City",
				}
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/subject_present.tf"),
					Check: minCertChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/min.tf"),
					ExpectError: regexp.MustCompile("updating field `subject` is not possible"),
				},
			},
		},
		"passive renewal - no signed cert, renew_pending is false": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create (no signed cert yet)
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "false").
						Build(),
				},
			},
		},
		"passive renewal - cert not yet in threshold, renew_pending is false": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Cert expires 2025-07-15, now is 2025-05-01, threshold is 30 days => renewal at 2025-06-15; not yet
				expiryDate := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
				createData.signedCertNotValidAfter = &expiryDate
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "false").
						Build(),
				},
			},
		},
		"passive renewal - cert within threshold, renew_pending is true": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Cert expires 2025-05-20, now is 2025-05-01, threshold is 30 days => renewal at 2025-04-20; already past
				expiryDate := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
				createData.signedCertNotValidAfter = &expiryDate
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "true").
						CheckEqual("auto_renew", "false").
						Build(),
				},
			},
		},
		"passive renewal - update: remove renewal days resets renew_pending": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Step 1: cert within threshold
				expiryDate := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
				createData.signedCertNotValidAfter = &expiryDate
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Step 1 refresh plan Read
				mockGetCertificate(m, createData)
				// Step 2: remove renewal days - no PatchCertificate call
				// Step 2 pre-apply plan Read
				mockGetCertificate(m, createData)
				// Step 2 refresh plan Read
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_pending", "true").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/full.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_pending", "false").
						CheckMissing("renew_before_expiration_days").
						Build(),
				},
			},
		},
		"passive renewal - update: change renewal days, no API call": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create (no signed cert)
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Step 1 refresh plan Read
				mockGetCertificate(m, createData)
				// Step 2: only renew_before_expiration_days changes, no PatchCertificate call
				// Step 2 pre-apply plan Read
				mockGetCertificate(m, createData)
				// Step 2 refresh plan Read
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/increase_renewal_days.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "60").
						CheckEqual("renew_pending", "false").
						Build(),
				},
			},
		},
		"passive renewal - update: decrease threshold dismisses renew_pending": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Step 1: cert within threshold (30 days), renew_pending=true
				// Cert expires 2025-05-20, now is 2025-05-01 => 19 days away, threshold=30 => pending
				expiryDate := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
				createData.signedCertNotValidAfter = &expiryDate
				// Create
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Step 1 refresh plan Read
				mockGetCertificate(m, createData)
				// Step 2: decrease threshold to 10 days => renewal at 2025-05-10, now (05-01) is before => not pending
				// Step 2 pre-apply plan Read
				mockGetCertificate(m, createData)
				// Step 2 refresh plan Read
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "true").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/update/decrease_renewal_days.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "10").
						CheckEqual("renew_pending", "false").
						Build(),
				},
			},
		},
		"passive renewal - update: enable auto_renew, no API call": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Create (no signed cert)
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Step 1 refresh plan Read
				mockGetCertificate(m, createData)
				// Step 2: only auto_renew changes, no PatchCertificate call
				// Step 2 pre-apply plan Read
				mockGetCertificate(m, createData)
				// Step 2 refresh plan Read
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_passive.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "false").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_active.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "true").
						Build(),
				},
			},
		},
		"active renewal - auto_renew true, renew_pending true triggers replacement": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Step 1: Create cert without signed cert = renew_pending is false
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Step 1 refresh plan Read (no signed cert → renew_pending stays false)
				mockGetCertificate(m, createData)

				// Step 2: simulate signed cert uploaded between steps
				certWithExpiry := createData
				expiryDate := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
				certWithExpiry.signedCertNotValidAfter = &expiryDate
				// Step 2 pre-apply plan Read → renew_pending becomes true → triggers replacement
				mockGetCertificate(m, certWithExpiry)

				// Step 2: replacement
				renewedData := createData
				renewedData.signedCertNotValidAfter = nil // new cert has no signed certificate
				renewedData.name = "test-name.renewed.2025-05-01T12_05_01Z"
				renewedData.baseName = "test-name.renewed.2025-05-01T12_05_01Z"
				renewedData.certificateID = "123456"
				mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
					ContractID:      createData.contractID,
					Domain:          createData.sans[0],
					CertificateName: createData.baseName,
				}, &cloudcertificates.ListCertificatesResponse{
					Certificates: []cloudcertificates.Certificate{
						{CertificateName: createData.name},
						{CertificateName: "test-name.renewed.2025-03-01T12_05_01Z"},
					},
				}, nil).Once()
				mockCreateCertificate(m, renewedData)
				// Delete old cert
				mockDeleteCertificate(m, createData)
				// Step 2 refresh plan Read (new cert, no signed cert → renew_pending=false)
				mockGetCertificate(m, renewedData)
				// Delete new cert during cleanup
				mockDeleteCertificate(m, renewedData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_active.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "true").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_active.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("akamai_cloudcertificates_certificate.test",
								plancheck.ResourceActionCreateBeforeDestroy),
						},
					},
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "true").
						CheckEqual("name", "test-name.renewed.2025-05-01T12_05_01Z").
						CheckEqual("certificate_id", "123456").
						Build(),
				},
			},
		},
		"active renewal - auto_renew true, renew_pending false, no replacement": {
			init: func(m *cloudcertificates.Mock, createData certificateTestData, _ certificateTestData) {
				// Cert expires far in the future - no renewal needed
				mockEmptyRenewalChain(m, createData)
				mockCreateCertificate(m, createData)
				// Read before destroy
				mockGetCertificate(m, createData)
				// Delete
				mockDeleteCertificate(m, createData)
			},
			createMockData: fullCertificateRSA,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificate/create/with_renewal_active.tf"),
					Check: fullCertChecker.
						CheckEqual("renew_before_expiration_days", "30").
						CheckEqual("renew_pending", "false").
						CheckEqual("auto_renew", "true").
						Build(),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if tc.init != nil {
				tc.init(client.CloudCertificates, tc.createMockData, tc.updateMockData)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}

}

func mockCreateCertificate(m *cloudcertificates.Mock, data certificateTestData) *mock.Call {
	var reqSubject *cloudcertificates.Subject
	if data.subject != nil {
		reqSubject = &cloudcertificates.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("CreateCertificate", testutils.MockContext, cloudcertificates.CreateCertificateRequest{
		ContractID: strings.TrimPrefix(data.contractID, "ctr_"),
		GroupID:    strings.TrimPrefix(data.groupID, "grp_"),
		Body: cloudcertificates.CreateCertificateRequestBody{
			CertificateName: data.baseName,
			KeyType:         data.keyType,
			KeySize:         data.keySize,
			SecureNetwork:   data.secureNetwork,
			GeoClass:        data.inputGeoClass,
			SANs:            data.sans,
			Subject:         reqSubject,
		},
	}).Return(&cloudcertificates.CreateCertificateResponse{
		Certificate: cloudcertificates.Certificate{
			AccountID:                          data.accountID,
			CertificateID:                      data.certificateID,
			CertificateName:                    data.name,
			CertificateStatus:                  data.certificateStatus,
			CertificateType:                    data.certificateType,
			ContractID:                         strings.TrimPrefix(data.contractID, "ctr_"),
			CreatedBy:                          data.createdBy,
			CreatedDate:                        tst.NewTimeFromStringMust(data.createdDate),
			ModifiedBy:                         data.modifiedBy,
			ModifiedDate:                       tst.NewTimeFromStringMust(data.modifiedDate),
			CSRExpirationDate:                  tst.NewTimeFromStringMust(data.csrExpirationDate),
			CSRPEM:                             data.csrPEM,
			KeyType:                            data.keyType,
			KeySize:                            data.keySize,
			SecureNetwork:                      string(data.secureNetwork),
			GeoClass:                           string(data.outputGeoClass),
			SANs:                               data.sans,
			Subject:                            reqSubject,
			SignedCertificateNotValidAfterDate: data.signedCertNotValidAfter,
		},
	}, nil).Once()
}

func mockGetCertificate(m *cloudcertificates.Mock, data certificateTestData) *mock.Call {
	var subject *cloudcertificates.Subject
	if data.subject != nil {
		subject = &cloudcertificates.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("GetCertificate", testutils.MockContext, cloudcertificates.GetCertificateRequest{
		CertificateID: data.certificateID,
	}).Return(&cloudcertificates.GetCertificateResponse{
		Certificate: cloudcertificates.Certificate{
			AccountID:                          data.accountID,
			CertificateID:                      data.certificateID,
			CertificateName:                    data.name,
			CertificateStatus:                  data.certificateStatus,
			CertificateType:                    data.certificateType,
			ContractID:                         strings.TrimPrefix(data.contractID, "ctr_"),
			CreatedBy:                          data.createdBy,
			CreatedDate:                        tst.NewTimeFromStringMust(data.createdDate),
			ModifiedBy:                         data.modifiedBy,
			ModifiedDate:                       tst.NewTimeFromStringMust(data.modifiedDate),
			CSRExpirationDate:                  tst.NewTimeFromStringMust(data.csrExpirationDate),
			CSRPEM:                             data.csrPEM,
			KeyType:                            data.keyType,
			KeySize:                            data.keySize,
			SecureNetwork:                      string(data.secureNetwork),
			GeoClass:                           string(data.outputGeoClass),
			SANs:                               data.sans,
			Subject:                            subject,
			SignedCertificateNotValidAfterDate: data.signedCertNotValidAfter,
		},
	}, nil).Once()
}

func mockDeleteCertificate(m *cloudcertificates.Mock, data certificateTestData) *mock.Call {
	return m.On("DeleteCertificate", testutils.MockContext, cloudcertificates.DeleteCertificateRequest{
		CertificateID: data.certificateID,
	}).Return(nil, nil).Once()
}

func mockPatchCertificate(m *cloudcertificates.Mock, data certificateTestData) *mock.Call {
	var subject *cloudcertificates.Subject
	if data.subject != nil {
		subject = &cloudcertificates.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("PatchCertificate", testutils.MockContext, cloudcertificates.PatchCertificateRequest{
		CertificateID:   data.certificateID,
		CertificateName: ptr.To(data.baseName),
	}).Return(&cloudcertificates.PatchCertificateResponse{
		Certificate: cloudcertificates.Certificate{
			AccountID:                          data.accountID,
			CertificateID:                      data.certificateID,
			CertificateName:                    data.name,
			CertificateStatus:                  data.certificateStatus,
			CertificateType:                    data.certificateType,
			ContractID:                         strings.TrimPrefix(data.contractID, "ctr_"),
			CreatedBy:                          data.createdBy,
			CreatedDate:                        tst.NewTimeFromStringMust(data.createdDate),
			ModifiedBy:                         data.modifiedBy,
			ModifiedDate:                       tst.NewTimeFromStringMust(data.modifiedDate),
			CSRExpirationDate:                  tst.NewTimeFromStringMust(data.csrExpirationDate),
			CSRPEM:                             data.csrPEM,
			KeyType:                            data.keyType,
			KeySize:                            data.keySize,
			SecureNetwork:                      string(data.secureNetwork),
			GeoClass:                           string(data.outputGeoClass),
			SANs:                               data.sans,
			Subject:                            subject,
			SignedCertificateNotValidAfterDate: data.signedCertNotValidAfter,
		},
	}, nil).Once()
}

func mockEmptyRenewalChain(m *cloudcertificates.Mock, data certificateTestData) {
	mockListCertificates(m, cloudcertificates.ListCertificatesRequest{
		ContractID:      data.contractID,
		Domain:          data.sans[0],
		CertificateName: data.baseName,
	}, &cloudcertificates.ListCertificatesResponse{
		Certificates: []cloudcertificates.Certificate{{}},
	}, nil).Once()
}

func TestExtractBaseName(t *testing.T) {

	tests := []struct {
		label       string
		name        string
		expBaseName string
	}{
		{"empty", "", ""},
		{"casual name", "foo", "foo"},
		{"renewed name", "foo.renewed.2025-05-01T12_05_01Z", "foo"},
		{"bad suffix", "foo.rotated.2025-05-01", "foo.rotated.2025-05-01"},
		{"non-existing date", "foo.renewed.2025-99-01", "foo.renewed.2025-99-01"},
		{"no basename", ".renewed.2025-05-01", ".renewed.2025-05-01"},
		{"no basename 2", "renewed.2025-05-01", "renewed.2025-05-01"},
		{"no date", "foo.renewed.", "foo.renewed."},
		{"no date 2", "foo.renewed", "foo.renewed"},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			res := extractBaseName(tc.name)
			assert.Equal(t, tc.expBaseName, res)
		})
	}
}

func TestDomainNameRegex(t *testing.T) {

	tests := []struct {
		domainName string
		matches    bool
		label      string
	}{
		// Valid cases
		{"example.com", true, "simple domain"},
		{"foo-bar.com", true, "domain with hyphen"},
		{"sub.example.pl", true, "subdomain"},
		{"sub.domain.co.uk", true, "multi-level subdomain"},
		{"*.example.com", true, "wildcard domain"},
		{"a.com", true, "single-character label"},
		{"example123.com", true, "domain with numbers"},
		{"123.com", true, "numeric domain"},
		{"foo-bar-baz123.example-domain.io", true, "long domain with hyphens and numbers"},

		// Invalid cases
		{"", false, "empty string"},
		{"*", false, "only wildcard"},
		{"*.com", false, "wildcard with no label"},
		{"www.*.com", false, "wildcard in middle"},
		{"*-foo.com", false, "wildcard with hyphen"},
		{"*_foo.com", false, "wildcard with underscore"},
		{"example.com.", false, "trailing dot"},
		{".example.com", false, "leading dot"},
		{"example..com", false, "consecutive dots"},
		{"-example.com", false, "leading hyphen"},
		{"example-.com", false, "trailing hyphen"},
		{"example", false, "no tld"},
		{"example.c", false, "tld too short"},
		{"example.c1m", false, "tld is not all letters"},
		{"Example.com", false, "uppercase letters not allowed"},
		{"EXAMPLE.COM", false, "uppercase domains not allowed"},
		{"example .com", false, "space not allowed"},
		{"foo_bar.com", false, "underscore not allowed"},
		{"foo,bar.com", false, "comma not allowed"},
		{"foo@bar.com", false, "at symbol not allowed"},
		{"foo/bar.com", false, "slash not allowed"},
		{"foo\\bar.com", false, "backslash not allowed"},
		{"foo!bar.com", false, "exclamation not allowed"},
		{"foo#bar.com", false, "hash not allowed"},
		{"foo$bar.com", false, "dollar sign not allowed"},
		{"foo%bar.com", false, "percent not allowed"},
		{"foo^bar.com", false, "caret not allowed"},
		{"foo&bar.com", false, "ampersand not allowed"},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			isMatch := domainNameRegex.MatchString(tc.domainName)
			assert.Equal(t, tc.matches, isMatch)
		})
	}
}

func TestIsWithinRenewalThreshold(t *testing.T) {
	now := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expiresInFuture := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC) // 59 days from now
	expiresSoon := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)      // 17 days from now
	expiresExactly := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)  // exactly 30 days from now

	tests := []struct {
		label                     string
		renewBeforeExpirationDays types.Int64
		signedCertNotValidAfter   *time.Time
		expected                  bool
	}{
		{
			label:                     "null renew_before_expiration_days returns false",
			renewBeforeExpirationDays: types.Int64Null(),
			signedCertNotValidAfter:   &expiresSoon,
			expected:                  false,
		},
		{
			label:                     "unknown renew_before_expiration_days returns false",
			renewBeforeExpirationDays: types.Int64Unknown(),
			signedCertNotValidAfter:   &expiresSoon,
			expected:                  false,
		},
		{
			label:                     "nil expiry date returns false",
			renewBeforeExpirationDays: types.Int64Value(30),
			signedCertNotValidAfter:   nil,
			expected:                  false,
		},
		{
			label:                     "cert expires far in future, outside threshold",
			renewBeforeExpirationDays: types.Int64Value(30),
			signedCertNotValidAfter:   &expiresInFuture,
			expected:                  false,
		},
		{
			label:                     "cert within threshold",
			renewBeforeExpirationDays: types.Int64Value(30),
			signedCertNotValidAfter:   &expiresSoon,
			expected:                  true,
		},
		{
			label:                     "cert exactly at threshold boundary, not yet within",
			renewBeforeExpirationDays: types.Int64Value(30),
			signedCertNotValidAfter:   &expiresExactly,
			expected:                  false,
		},
		{
			label:                     "zero days threshold, now before expiry",
			renewBeforeExpirationDays: types.Int64Value(0),
			signedCertNotValidAfter:   &expiresSoon,
			expected:                  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			result := isWithinRenewalThreshold(tc.renewBeforeExpirationDays, now, tc.signedCertNotValidAfter)
			assert.Equal(t, tc.expected, result)
		})
	}
}

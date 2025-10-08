package cloudcertificates

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
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
		keyType       ccm.CryptographicAlgorithm
		keySize       ccm.KeySize
		secureNetwork ccm.SecureNetwork
		sans          []string
		subject       *ccm.Subject

		// output data
		certificateID     string
		certificateType   string
		name              string
		certificateStatus string
		accountID         string
		createdBy         string
		createdDate       string
		modifiedBy        string
		modifiedDate      string
		csrExpirationDate string
		csrPEM            string
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
	}

	minCertificateWithPrefixes = certificateTestData{
		// input data
		contractID:    "ctr_test_contract",
		groupID:       "grp_123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
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
	}

	fullCertificateRSA = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &ccm.Subject{
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
	}

	updateCertificate = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name-updated",
		groupID:       "123",
		keyType:       "RSA",
		keySize:       "2048",
		secureNetwork: "ENHANCED_TLS",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &ccm.Subject{
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
	}

	fullCertificateECDSA = certificateTestData{
		// input data
		contractID:    "test_contract",
		baseName:      "test-name",
		groupID:       "123",
		keyType:       "ECDSA",
		keySize:       "P-256",
		secureNetwork: "ENHANCED_TLS",
		sans:          []string{"test.example.com", "test.example2.com"},
		subject: &ccm.Subject{
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
	}
)

func TestCertificateResource(t *testing.T) {
	t.Parallel()

	minCertChecker := test.NewStateChecker("akamai_cloudcertificates_certificate.test").
		CheckEqual("contract_id", "test_contract").
		CheckEqual("group_id", "123").
		CheckEqual("key_type", "RSA").
		CheckEqual("key_size", "2048").
		CheckEqual("secure_network", "ENHANCED_TLS").
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
		CheckEqual("csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n")

	tests := map[string]struct {
		init           func(*ccm.Mock, certificateTestData, certificateTestData)
		createMockData certificateTestData
		updateMockData certificateTestData
		steps          []resource.TestStep
	}{
		"happy path - create certificate without optionals": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
		"happy path - create certificate with prefixes and without optionals": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
		"happy path - create certificate with all optional attributes": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
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
		"happy path - create certificate with optional attributes, different key type, some missing subject fields": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
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
						CheckMissing("subject.common_name").
						CheckMissing("subject.country").
						CheckMissing("subject.organization").
						Build(),
				},
			},
		},
		"happy path - create certificate, update name": {
			init: func(m *ccm.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
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
			init: func(m *ccm.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
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
		"happy path - create certificate, change order of SANs - no diff": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before refresh
				mockGetCertificate(m, createData)
				// Read from refresh step
				m.On("GetCertificate", testutils.MockContext, ccm.GetCertificateRequest{
					CertificateID: createData.certificateID,
				}).Return(nil, ccm.ErrCertificateNotFound).Once()
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
		"import - not renewed certificate": {
			init: func(m *ccm.Mock, data certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, data certificateTestData, _ certificateTestData) {
				data.name = "test-certificate.renewed.2025-05-01"
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
						CheckEqual("name", "test-certificate.renewed.2025-05-01").
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
			init: func(m *ccm.Mock, data certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, _ certificateTestData, _ certificateTestData) {
				// Import
				m.On("GetCertificate", testutils.MockContext, ccm.GetCertificateRequest{
					CertificateID: "12345abc-wrong",
				}).Return(nil, ccm.ErrCertificateNotFound).Once()
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
			init: func(m *ccm.Mock, data certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				m.On("CreateCertificate", testutils.MockContext, ccm.CreateCertificateRequest{
					ContractID: createData.contractID,
					GroupID:    createData.groupID,
					Body: ccm.CreateCertificateRequestBody{
						CertificateName: createData.baseName,
						KeyType:         createData.keyType,
						KeySize:         createData.keySize,
						SecureNetwork:   createData.secureNetwork,
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before destroy
				m.On("GetCertificate", testutils.MockContext, ccm.GetCertificateRequest{
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
			init: func(m *ccm.Mock, createData certificateTestData, updateData certificateTestData) {
				// Create
				mockCreateCertificate(m, createData)
				// Read before update
				mockGetCertificate(m, createData)
				// Update
				m.On("PatchCertificate", testutils.MockContext, ccm.PatchCertificateRequest{
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
		"expect error - missing contract": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_contract.tf"),
					ExpectError: regexp.MustCompile("The argument \"contract_id\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing group": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_group.tf"),
					ExpectError: regexp.MustCompile(`Error: Required Field Missing(.|\n)*field ` + "`group_id`" + ` is required during creation`),
				},
			},
		},
		"expect error - missing key_size": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_key_size.tf"),
					ExpectError: regexp.MustCompile("The argument \"key_size\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing key_type": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_key_type.tf"),
					ExpectError: regexp.MustCompile("The argument \"key_type\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing secure_network": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_secure_network.tf"),
					ExpectError: regexp.MustCompile("The argument \"secure_network\" is required, but no definition was found."),
				},
			},
		},
		"expect error - missing sans": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/no_sans.tf"),
					ExpectError: regexp.MustCompile("The argument \"sans\" is required, but no definition was found."),
				},
			},
		},
		"expect error - empty sans": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_sans.tf"),
					ExpectError: regexp.MustCompile("Attribute sans set must contain at least 1 elements and at most 100 elements,\ngot: 0"),
				},
			},
		},
		"expect error - more than 100 sans provided": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/more_than_100_sans.tf"),
					ExpectError: regexp.MustCompile("Attribute sans set must contain at least 1 elements and at most 100 elements,\ngot: 101"),
				},
			},
		},
		"expect error - one of sans is not a valid domain name": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/sans_not_valid_domain.tf"),
					ExpectError: regexp.MustCompile(`Attribute sans.+"invalid_domain".+must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - one of sans contains uppercase letters": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/sans_uppercase.tf"),
					ExpectError: regexp.MustCompile(`Attribute sans.+"example.COM".+must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - empty base_name": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_base_name.tf"),
					ExpectError: regexp.MustCompile("Attribute base_name string length must be at least 1, got: 0"),
				},
			},
		},
		"expect error - wrong key_size for RSA": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_size_rsa.tf"),
					ExpectError: regexp.MustCompile(`The specified value '2137' for the RSA key type is invalid. Valid values are(.|\n)*'2048'.`),
				},
			},
		},
		"expect error - wrong key_size for ECDSA": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_size_ecdsa.tf"),
					ExpectError: regexp.MustCompile(`The specified value '2137' for the ECDSA key type is invalid. Valid values(.|\n)*are 'P-256'.`),
				},
			},
		},
		"expect error - wrong key_type": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_key_type.tf"),
					ExpectError: regexp.MustCompile(`Attribute key_type value must be one of: \["RSA" "ECDSA"\], got: "WRONG-TYPE"`),
				},
			},
		},
		"expect error - wrong secure_network": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/wrong_secure_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute secure_network value must be one of: \["ENHANCED_TLS"\], got:(.|\n)*"WRONG_NETWORK"`),
				},
			},
		},
		"expect error - empty subject": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_subject.tf"),
					ExpectError: regexp.MustCompile(`At least one of the subject fields \(common_name, organization, country,(.|\n)*state, locality\) must be specified.`),
				},
			},
		},
		"expect error - common_name not present in sans": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_not_in_sans.tf"),
					ExpectError: regexp.MustCompile(`The specified common name 'test.example.com' must be included in the SANs(.|\n)*list.`),
				},
			},
		},
		"expect error - common_name is not a valid domain name": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_not_valid_domain.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject\.common_name must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - common_name contains uppercase letters": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/common_name_contains_uppercase.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject\.common_name must\sbe\sa\svalid\sdomain\sname\swith\sall\sletters\slowercase`),
				},
			},
		},
		"expect error - empty organization": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_organization.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.organization string length must be between 1 and 64, got: 0`),
				},
			},
		},
		"expect error - country too long": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/country_too_long.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.country string length must be between 2 and 2, got: 7`),
				},
			},
		},
		"expect error - empty state": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_state.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.state string length must be between 1 and 128, got: 0`),
				},
			},
		},
		"expect error - empty locality": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/empty_locality.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.locality string length must be between 1 and 128, got: 0`),
				},
			},
		},
		"expect error - organization with only spaces": {
			init:           func(_ *ccm.Mock, _ certificateTestData, _ certificateTestData) {},
			createMockData: minCertificate,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificate/validation/organization_only_spaces.tf"),
					ExpectError: regexp.MustCompile(`Attribute subject.organization cannot be empty or whitespace, got:  `),
				},
			},
		},
		"expect error - update contract": {
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				// Create
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
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
			init: func(m *ccm.Mock, createData certificateTestData, _ certificateTestData) {
				createData.subject = &ccm.Subject{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &ccm.Mock{}

			if tc.init != nil {
				tc.init(client, tc.createMockData, tc.updateMockData)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}

}

func mockCreateCertificate(m *ccm.Mock, data certificateTestData) *mock.Call {
	var reqSubject *ccm.Subject
	if data.subject != nil {
		reqSubject = &ccm.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("CreateCertificate", testutils.MockContext, ccm.CreateCertificateRequest{
		ContractID: strings.TrimPrefix(data.contractID, "ctr_"),
		GroupID:    strings.TrimPrefix(data.groupID, "grp_"),
		Body: ccm.CreateCertificateRequestBody{
			CertificateName: data.baseName,
			KeyType:         data.keyType,
			KeySize:         data.keySize,
			SecureNetwork:   data.secureNetwork,
			SANs:            data.sans,
			Subject:         reqSubject,
		},
	}).Return(&ccm.CreateCertificateResponse{
		Certificate: ccm.Certificate{
			AccountID:         data.accountID,
			CertificateID:     data.certificateID,
			CertificateName:   data.name,
			CertificateStatus: data.certificateStatus,
			CertificateType:   data.certificateType,
			ContractID:        strings.TrimPrefix(data.contractID, "ctr_"),
			CreatedBy:         data.createdBy,
			CreatedDate:       tst.NewTimeFromStringMust(data.createdDate),
			ModifiedBy:        data.modifiedBy,
			ModifiedDate:      tst.NewTimeFromStringMust(data.modifiedDate),
			CSRExpirationDate: tst.NewTimeFromStringMust(data.csrExpirationDate),
			CSRPEM:            ptr.To(data.csrPEM),
			KeyType:           data.keyType,
			KeySize:           data.keySize,
			SecureNetwork:     string(data.secureNetwork),
			SANs:              data.sans,
			Subject:           reqSubject,
		},
	}, nil).Once()
}

func mockGetCertificate(m *ccm.Mock, data certificateTestData) *mock.Call {
	var subject *ccm.Subject
	if data.subject != nil {
		subject = &ccm.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("GetCertificate", testutils.MockContext, ccm.GetCertificateRequest{
		CertificateID: data.certificateID,
	}).Return(&ccm.GetCertificateResponse{
		AccountID:         data.accountID,
		CertificateID:     data.certificateID,
		CertificateName:   data.name,
		CertificateStatus: data.certificateStatus,
		CertificateType:   data.certificateType,
		ContractID:        strings.TrimPrefix(data.contractID, "ctr_"),
		CreatedBy:         data.createdBy,
		CreatedDate:       tst.NewTimeFromStringMust(data.createdDate),
		ModifiedBy:        data.modifiedBy,
		ModifiedDate:      tst.NewTimeFromStringMust(data.modifiedDate),
		CSRExpirationDate: tst.NewTimeFromStringMust(data.csrExpirationDate),
		CSRPEM:            ptr.To(data.csrPEM),
		KeyType:           data.keyType,
		KeySize:           data.keySize,
		SecureNetwork:     string(data.secureNetwork),
		SANs:              data.sans,
		Subject:           subject,
	}, nil).Once()
}

func mockDeleteCertificate(m *ccm.Mock, data certificateTestData) *mock.Call {
	return m.On("DeleteCertificate", testutils.MockContext, ccm.DeleteCertificateRequest{
		CertificateID: data.certificateID,
	}).Return(nil).Once()
}

func mockPatchCertificate(m *ccm.Mock, data certificateTestData) *mock.Call {
	var subject *ccm.Subject
	if data.subject != nil {
		subject = &ccm.Subject{
			CommonName:   data.subject.CommonName,
			Country:      data.subject.Country,
			Organization: data.subject.Organization,
			State:        data.subject.State,
			Locality:     data.subject.Locality,
		}
	}
	return m.On("PatchCertificate", testutils.MockContext, ccm.PatchCertificateRequest{
		CertificateID:   data.certificateID,
		CertificateName: ptr.To(data.baseName),
	}).Return(&ccm.PatchCertificateResponse{
		AccountID:         data.accountID,
		CertificateID:     data.certificateID,
		CertificateName:   data.name,
		CertificateStatus: data.certificateStatus,
		CertificateType:   data.certificateType,
		ContractID:        strings.TrimPrefix(data.contractID, "ctr_"),
		CreatedBy:         data.createdBy,
		CreatedDate:       tst.NewTimeFromStringMust(data.createdDate),
		ModifiedBy:        data.modifiedBy,
		ModifiedDate:      tst.NewTimeFromStringMust(data.modifiedDate),
		CSRExpirationDate: tst.NewTimeFromStringMust(data.csrExpirationDate),
		CSRPEM:            ptr.To(data.csrPEM),
		KeyType:           data.keyType,
		KeySize:           data.keySize,
		SecureNetwork:     string(data.secureNetwork),
		SANs:              data.sans,
		Subject:           subject,
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
		{"renewed name", "foo.renewed.2025-05-01", "foo"},
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

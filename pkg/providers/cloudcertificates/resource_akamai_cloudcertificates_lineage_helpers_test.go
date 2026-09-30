package cloudcertificates

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
)

// This file collects response/request builders and mock-setup helpers shared by the lineage resource's test
// scenarios, so the scenario table in resource_akamai_cloudcertificates_lineage_test.go can stay focused on what
// makes each scenario distinct instead of repeating lifecycle orchestration boilerplate.

// createLineageResponseFull covers a lineage with two key specs (MULTIPLE_STACK), a subject and a DOM-validation
// warning on its SANs.
func createLineageResponseFull() *cloudcertificates.CreateLineageResponse {
	return &cloudcertificates.CreateLineageResponse{
		AccountID:  "A-CCT1234",
		ContractID: "C-0N7RAC7",
		GeoClass:   "STANDARD_WORLDWIDE",
		GroupID:    12345,
		Head: &cloudcertificates.HeadGeneration{
			Generation: cloudcertificates.Generation{
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						AlgorithmInstanceID:          6122,
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "RSA",
					},
					{
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						AlgorithmInstanceID:          6123,
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "ECDSA",
					},
				},
				GenerationCreatedBy:   ptr.To("terraform-dev"),
				GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
				GenerationModifiedBy:  ptr.To("terraform-dev"),
			},
			HeadGenerationID:     2912,
			HeadGenerationStatus: "CSR_READY",
		},
		KeySpecs: []cloudcertificates.KeySpecResponse{
			{KeySize: "2048", KeyType: "RSA"},
			{KeySize: "P-256", KeyType: "ECDSA"},
		},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  tst.NewTimeFromStringMust("2026-07-07T10:18:46Z"),
		LineageID:           500001,
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: tst.NewTimeFromStringMust("2026-07-07T10:18:46Z"),
		LineageName:         "test non-validated domain",
		LineageType:         "MULTIPLE_GENERATION",
		SANs:                []string{"www.example.com", "example.com"},
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "MULTIPLE_STACK",
		Subject: cloudcertificates.Subject{
			CommonName:         "example.com",
			Country:            "US",
			Locality:           "Cambridge",
			Organization:       "Example Corp.",
			OrganizationalUnit: "IT",
			State:              "Massachusetts",
		},
		ValidationResults: &cloudcertificates.ValidationResults{
			Warnings: []cloudcertificates.ValidationResultItem{
				{
					Context:  map[string]any{"domValidatorLink": "/domain-validation/v1/domains"},
					Detail:   "Domain validation must be completed before uploading the signed certificate.",
					Instance: "/error-types/domain-not-validated?traceId=1234567891014",
					Status:   400,
					Title:    "Some SANs in the request are not Domain Validated.",
					Type:     "/error-types/domain-not-validated",
				},
			},
		},
	}
}

func createLineageRequestFull() cloudcertificates.CreateLineageRequest {
	return cloudcertificates.CreateLineageRequest{
		Body: cloudcertificates.CreateLineageRequestBody{
			ContractID:    "C-0N7RAC7",
			GroupID:       12345,
			GeoClass:      "STANDARD_WORLDWIDE",
			SecureNetwork: "ENHANCED_TLS",
			KeySpecs: []cloudcertificates.KeySpec{
				{KeyType: "ECDSA", KeySize: "P-256"},
				{KeyType: "RSA", KeySize: "2048"},
			},
			SANs: []string{"example.com", "www.example.com"},
			Subject: ptr.To(cloudcertificates.Subject{
				CommonName:         "example.com",
				Organization:       "Example Corp.",
				OrganizationalUnit: "IT",
				Country:            "US",
				State:              "Massachusetts",
				Locality:           "Cambridge",
			}),
		},
	}
}

// createLineageResponseMinimal covers a lineage with a single key spec (SINGLE_STACK) and no subject.
func createLineageResponseMinimal() *cloudcertificates.CreateLineageResponse {
	return &cloudcertificates.CreateLineageResponse{
		AccountID:  "A-CCT1234",
		ContractID: "C-0N7RAC7",
		GeoClass:   "STANDARD_WORLDWIDE",
		GroupID:    12345,
		Head: &cloudcertificates.HeadGeneration{
			Generation: cloudcertificates.Generation{
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-20T10:32:09Z")),
						AlgorithmInstanceID:          2244,
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-20T10:32:09Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "RSA",
					},
				},
				GenerationCreatedBy:   ptr.To("terraform-dev"),
				GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-20T10:32:09Z")),
				GenerationModifiedBy:  ptr.To("terraform-dev"),
			},
			HeadGenerationID:     1323,
			HeadGenerationStatus: "CSR_READY",
		},
		KeySpecs: []cloudcertificates.KeySpecResponse{
			{KeySize: "2048", KeyType: "RSA"},
		},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  tst.NewTimeFromStringMust("2026-07-20T10:32:09Z"),
		LineageID:           500006,
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: tst.NewTimeFromStringMust("2026-07-20T10:32:09Z"),
		LineageName:         "www.example.com20260720103154092771",
		LineageType:         "MULTIPLE_GENERATION",
		SANs:                []string{"www.example.com", "example.com"},
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "SINGLE_STACK",
		Subject:             cloudcertificates.Subject{},
	}
}

func createLineageRequestMinimal() cloudcertificates.CreateLineageRequest {
	return cloudcertificates.CreateLineageRequest{
		Body: cloudcertificates.CreateLineageRequestBody{
			ContractID:    "C-0N7RAC7",
			GroupID:       12345,
			SecureNetwork: "ENHANCED_TLS",
			KeySpecs: []cloudcertificates.KeySpec{
				{KeyType: "RSA", KeySize: "2048"},
			},
			SANs: []string{"example.com", "www.example.com"},
		},
	}
}

// renamedLineageResponseFull returns createLineageResponseFull with the rename's known field changes applied
// (LineageName, LineageModifiedTime) and ValidationResults cleared (a real post-rename GetLineage never
// re-surfaces DOM-validation warnings).
func renamedLineageResponseFull() *cloudcertificates.CreateLineageResponse {
	renamed := createLineageResponseFull()
	renamed.LineageName = "renamed-lineage"
	renamed.LineageModifiedTime = tst.NewTimeFromStringMust("2026-07-08T13:24:17Z")
	renamed.ValidationResults = nil
	return renamed
}

// renameLineageResponseFull builds the RenameLineageResponse for createLineageResponseFull's lineage renamed to
// "renamed-lineage". RenameLineage's response omits algorithm/generation-common detail, unlike GetLineage/CreateLineage.
func renameLineageResponseFull() *cloudcertificates.RenameLineageResponse {
	return &cloudcertificates.RenameLineageResponse{
		AccountID:  "A-CCT1234",
		ContractID: "C-0N7RAC7",
		GeoClass:   "STANDARD_WORLDWIDE",
		GroupID:    12345,
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     2912,
			HeadGenerationStatus: "READY_FOR_USE",
		},
		KeySpecs: []cloudcertificates.KeySpecResponse{
			{KeySize: "2048", KeyType: "RSA"},
			{KeySize: "P-256", KeyType: "ECDSA"},
		},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  tst.NewTimeFromStringMust("2026-07-07T10:18:46Z"),
		LineageID:           500001,
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: tst.NewTimeFromStringMust("2026-07-08T13:24:17Z"),
		LineageName:         "renamed-lineage",
		LineageType:         "MULTIPLE_GENERATION",
		SANs:                []string{"www.example.com", "example.com"},
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "MULTIPLE_STACK",
		Subject: cloudcertificates.Subject{
			CommonName:         "example.com",
			Country:            "US",
			Locality:           "Cambridge",
			Organization:       "Example Corp.",
			OrganizationalUnit: "IT",
			State:              "Massachusetts",
		},
	}
}

// checkLineageFullAttrs extends a checker with exhaustive assertions for subject, key_specs, and head.algorithms,
// matching createLineageResponseFull's data.
func checkLineageFullAttrs(c test.StateChecker) test.StateChecker {
	c = c.
		CheckEqualBatch("subject.", test.AttributeBatch{
			"common_name":         "example.com",
			"organization":        "Example Corp.",
			"organizational_unit": "IT",
			"country":             "US",
			"state":               "Massachusetts",
			"locality":            "Cambridge",
		}).
		CheckEqual("key_specs.RSA", "2048").
		CheckEqual("key_specs.ECDSA", "P-256").
		CheckEqualBatch("head.algorithms.RSA.", test.AttributeBatch{
			"algorithm_instance_id":           "6122",
			"algorithm_instance_created_by":   "terraform-dev",
			"algorithm_instance_created_time": "2026-07-07T10:18:46Z",
			"certificate_status":              "CSR_READY",
			"csr_expiration_date":             "2027-07-07T10:18:46Z",
			"csr_pem":                         "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
		}).
		CheckEqualBatch("head.algorithms.ECDSA.", test.AttributeBatch{
			"algorithm_instance_id":           "6123",
			"algorithm_instance_created_by":   "terraform-dev",
			"algorithm_instance_created_time": "2026-07-07T10:18:46Z",
			"certificate_status":              "CSR_READY",
			"csr_expiration_date":             "2027-07-07T10:18:46Z",
			"csr_pem":                         "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
		})

	return c.
		CheckMissing("head.algorithms.RSA.algorithm_instance_modified_by").
		CheckMissing("head.algorithms.RSA.algorithm_instance_modified_time").
		CheckMissing("head.algorithms.RSA.signed_certificate_issuer").
		CheckMissing("head.algorithms.RSA.signed_certificate_not_valid_after_date").
		CheckMissing("head.algorithms.RSA.signed_certificate_not_valid_before_date").
		CheckMissing("head.algorithms.RSA.signed_certificate_pem").
		CheckMissing("head.algorithms.RSA.signed_certificate_serial_number").
		CheckMissing("head.algorithms.RSA.signed_certificate_sha256_fingerprint").
		CheckMissing("head.algorithms.RSA.trust_chain_pem").
		CheckMissing("head.algorithms.ECDSA.algorithm_instance_modified_by").
		CheckMissing("head.algorithms.ECDSA.algorithm_instance_modified_time").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_issuer").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_not_valid_after_date").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_not_valid_before_date").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_pem").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_serial_number").
		CheckMissing("head.algorithms.ECDSA.signed_certificate_sha256_fingerprint").
		CheckMissing("head.algorithms.ECDSA.trust_chain_pem")
}

// mockCreateLineage registers a successful CreateLineage expectation.
func mockCreateLineage(m *cloudcertificates.Mock, req cloudcertificates.CreateLineageRequest, resp *cloudcertificates.CreateLineageResponse) {
	m.On("CreateLineage", testutils.MockContext, req).Return(resp, nil).Once()
}

// mockCreateLineageFails registers a CreateLineage expectation that fails with err.
func mockCreateLineageFails(m *cloudcertificates.Mock, req cloudcertificates.CreateLineageRequest, err error) {
	m.On("CreateLineage", testutils.MockContext, req).Return(nil, err).Once()
}

// mockGetLineage registers a GetLineage expectation for lineageID returning resp, called n times.
func mockGetLineage(m *cloudcertificates.Mock, lineageID int64, resp *cloudcertificates.CreateLineageResponse, n int) {
	m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: allExpandGenerations,
	}).Return((*cloudcertificates.GetLineageResponse)(resp), nil).Times(n)
}

// mockGetLineageFails registers a GetLineage expectation for lineageID that fails with err, called n times.
func mockGetLineageFails(m *cloudcertificates.Mock, lineageID int64, err error, n int) {
	m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: allExpandGenerations,
	}).Return(nil, err).Times(n)
}

// mockDeleteLineage registers a successful DeleteLineage expectation for lineageID.
func mockDeleteLineage(m *cloudcertificates.Mock, lineageID int64) {
	m.On("DeleteLineage", testutils.MockContext, cloudcertificates.DeleteLineageRequest{LineageID: lineageID}).
		Return(nil).Once()
}

// mockDeleteLineageFails registers a DeleteLineage expectation for lineageID that fails with err.
func mockDeleteLineageFails(m *cloudcertificates.Mock, lineageID int64, err error) {
	m.On("DeleteLineage", testutils.MockContext, cloudcertificates.DeleteLineageRequest{LineageID: lineageID}).
		Return(err).Once()
}

// mockListLineagesEmpty registers a ListLineages expectation for name returning no matches, called n times.
func mockListLineagesEmpty(m *cloudcertificates.Mock, name string, n int) {
	m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
		LineageName: name, PageSize: 100,
	}).Return(&cloudcertificates.ListLineagesResponse{}, nil).Times(n)
}

// mockListLineagesFails registers a ListLineages expectation for name that fails with err, called n times.
func mockListLineagesFails(m *cloudcertificates.Mock, name string, err error, n int) {
	m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
		LineageName: name, PageSize: 100,
	}).Return(nil, err).Times(n)
}

// mockListLineagesConflict registers a ListLineages expectation for name returning one existing lineage
// (existingID) with that same name.
func mockListLineagesConflict(m *cloudcertificates.Mock, name string, existingID int64) {
	m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
		LineageName: name, PageSize: 100,
	}).Return(&cloudcertificates.ListLineagesResponse{
		Lineages:   []cloudcertificates.Lineage{{LineageID: existingID, LineageName: name}},
		TotalCount: 1,
	}, nil).Once()
}

// mockRenameLineage registers a successful RenameLineage expectation.
func mockRenameLineage(m *cloudcertificates.Mock, lineageID int64, newName string, resp *cloudcertificates.RenameLineageResponse) {
	m.On("RenameLineage", testutils.MockContext, cloudcertificates.RenameLineageRequest{
		LineageID: lineageID, LineageName: newName,
	}).Return(resp, nil).Once()
}

// mockRenameLineageFails registers a RenameLineage expectation that fails with err.
func mockRenameLineageFails(m *cloudcertificates.Mock, lineageID int64, newName string, err error) {
	m.On("RenameLineage", testutils.MockContext, cloudcertificates.RenameLineageRequest{
		LineageID: lineageID, LineageName: newName,
	}).Return(nil, err).Once()
}

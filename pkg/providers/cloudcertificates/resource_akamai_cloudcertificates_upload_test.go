package cloudcertificates

import (
	"errors"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// uploadRSARequest is the real captured UploadSignedCertificate request for uploading only the RSA certificate
// onto lineage 500001's still-draft (not yet promoted) head generation 2912.
func uploadRSARequest() cloudcertificates.UploadSignedCertificateRequest {
	return uploadRSARequestForGeneration(2912)
}

// uploadRSARequestForGeneration is uploadRSARequest targeting a different generation, used to simulate the
// lineage's head generation having moved on to a new one.
func uploadRSARequestForGeneration(generationID int64) cloudcertificates.UploadSignedCertificateRequest {
	return cloudcertificates.UploadSignedCertificateRequest{
		LineageID:           500001,
		GenerationID:        generationID,
		AcknowledgeWarnings: false,
		Body: cloudcertificates.UploadSignedCertificateRequestBody{
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmRSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

// uploadRSARequestAck is uploadRSARequest with acknowledge_warnings explicitly configured to true.
func uploadRSARequestAck() cloudcertificates.UploadSignedCertificateRequest {
	req := uploadRSARequest()
	req.AcknowledgeWarnings = true
	return req
}

// uploadRSAResponse is the real captured UploadSignedCertificate response for the request above: RSA becomes
// READY_FOR_USE, ECDSA remains an unsigned CSR_READY slot (the lineage's second key spec, not yet uploaded).
func uploadRSAResponse() *cloudcertificates.UploadSignedCertificateResponse {
	return &cloudcertificates.UploadSignedCertificateResponse{
		Algorithms: []cloudcertificates.Algorithm{
			{
				CertificateStatus: "CSR_READY",
				CSRExpirationDate: ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
				KeyType:           "ECDSA",
			},
			{
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
				KeyType:                             "RSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:06:53Z")),
				SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
			},
		},
		GenerationID:           2912,
		GenerationModifiedBy:   "terraform-dev",
		GenerationModifiedTime: tst.NewTimeFromStringMust("2026-07-08T13:14:25Z"),
		GenerationStatus:       "READY_FOR_USE",
		LineageID:              500001,
	}
}

func uploadRSAResponseForGeneration(generationID int64) *cloudcertificates.UploadSignedCertificateResponse {
	response := uploadRSAResponse()
	response.GenerationID = generationID
	return response
}

// getGenerationRSAOnly is the GetGeneration equivalent of uploadRSAResponse (same lineage/generation/algorithm
// data), used for the post-upload re-fetch.
func getGenerationRSAOnly() *cloudcertificates.GetGenerationResponse {
	return &cloudcertificates.GetGenerationResponse{
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{
				{
					AlgorithmInstanceCreatedBy:   "terraform-dev",
					AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:          6123,
					CertificateStatus:            "CSR_READY",
					CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                      "ECDSA",
				},
				{
					AlgorithmInstanceCreatedBy:          "terraform-dev",
					AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:                 6122,
					AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
					AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					CertificateStatus:                   "READY_FOR_USE",
					CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                             "RSA",
					SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
					SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:06:53Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"),
					SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
					SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
				},
			},
			GenerationCreatedBy:    ptr.To("terraform-dev"),
			GenerationCreatedTime:  ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
			GenerationModifiedBy:   ptr.To("terraform-dev"),
			GenerationModifiedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:14:25Z")),
		},
		GenerationID:     2912,
		GenerationStatus: "READY_FOR_USE",
	}
}

// getGenerationNoneSigned is generation 2912 before anything has been uploaded: both algorithms are still
// unsigned CSR_READY slots. Used for Create's pre-upload check that neither configured algorithm already has an
// accepted signed certificate.
func getGenerationNoneSigned() *cloudcertificates.GetGenerationResponse {
	return &cloudcertificates.GetGenerationResponse{
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{
				{
					AlgorithmInstanceCreatedBy:   "terraform-dev",
					AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:          6123,
					CertificateStatus:            "CSR_READY",
					CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                      "ECDSA",
				},
				{
					AlgorithmInstanceCreatedBy:   "terraform-dev",
					AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:          6122,
					CertificateStatus:            "CSR_READY",
					CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                      "RSA",
				},
			},
			GenerationCreatedBy:   ptr.To("terraform-dev"),
			GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
			GenerationModifiedBy:  ptr.To("terraform-dev"),
		},
		GenerationID:     2912,
		GenerationStatus: "CSR_READY",
	}
}

// getGenerationRSAOnlyForGeneration is getGenerationRSAOnly bound to a different generation ID, used to simulate
// the lineage's head generation having moved on to a new one.
func getGenerationRSAOnlyForGeneration(generationID int64) *cloudcertificates.GetGenerationResponse {
	generation := getGenerationRSAOnly()
	generation.GenerationID = generationID
	return generation
}

// getGenerationNoneSignedForGeneration is getGenerationNoneSigned bound to a different generation ID, used for
// Create's pre-upload check when the resource is replaced onto a new head generation.
func getGenerationNoneSignedForGeneration(generationID int64) *cloudcertificates.GetGenerationResponse {
	generation := getGenerationNoneSigned()
	generation.GenerationID = generationID
	return generation
}

// getGenerationECDSAOnly is getGenerationRSAOnly with the roles reversed: ECDSA already has an accepted signed
// certificate, RSA is still an unsigned CSR_READY slot - used to prove Create's pre-upload check rejects the
// head generation on ANY already-signed algorithm, not just ones this resource instance configures.
func getGenerationECDSAOnly() *cloudcertificates.GetGenerationResponse {
	return &cloudcertificates.GetGenerationResponse{
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{
				{
					AlgorithmInstanceCreatedBy:          "terraform-dev",
					AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:                 6123,
					AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
					AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					CertificateStatus:                   "READY_FOR_USE",
					CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                             "ECDSA",
					SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
					SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:06:53Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"),
					SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02"),
					SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:02"),
				},
				{
					AlgorithmInstanceCreatedBy:   "terraform-dev",
					AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:          6122,
					CertificateStatus:            "CSR_READY",
					CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                      "RSA",
				},
			},
			GenerationCreatedBy:    ptr.To("terraform-dev"),
			GenerationCreatedTime:  ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
			GenerationModifiedBy:   ptr.To("terraform-dev"),
			GenerationModifiedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:14:25Z")),
		},
		GenerationID:     2912,
		GenerationStatus: "READY_FOR_USE",
	}
}

// uploadRSAECDSARequest is the request for uploading both algorithms atomically
func uploadRSAECDSARequest() cloudcertificates.UploadSignedCertificateRequest {
	return cloudcertificates.UploadSignedCertificateRequest{
		LineageID:           500001,
		GenerationID:        2912,
		AcknowledgeWarnings: false,
		Body: cloudcertificates.UploadSignedCertificateRequestBody{
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmRSA:   {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"},
				cloudcertificates.CryptographicAlgorithmECDSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

// uploadECDSARequest is the request for adding just the ECDSA algorithm to a generation that already has an
// accepted RSA certificate: Update only ever sends newly-added algorithms, never ones already uploaded.
func uploadECDSARequest() cloudcertificates.UploadSignedCertificateRequest {
	return cloudcertificates.UploadSignedCertificateRequest{
		LineageID:           500001,
		GenerationID:        2912,
		AcknowledgeWarnings: false,
		Body: cloudcertificates.UploadSignedCertificateRequestBody{
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmECDSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

func uploadRSAECDSAResponse() *cloudcertificates.UploadSignedCertificateResponse {
	return &cloudcertificates.UploadSignedCertificateResponse{
		Algorithms: []cloudcertificates.Algorithm{
			{
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
				KeyType:                             "RSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:06:53Z")),
				SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
			},
			{
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
				KeyType:                             "ECDSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:23:37Z")),
				SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:23:37Z")),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:02"),
			},
		},
		GenerationID:           2912,
		GenerationModifiedBy:   "terraform-dev",
		GenerationModifiedTime: tst.NewTimeFromStringMust("2026-07-08T13:24:17Z"),
		GenerationStatus:       "READY_FOR_USE",
		LineageID:              500001,
	}
}

// getGenerationRSAECDSA is the real captured GetLineage-embedded head generation detail for lineage 500001,
// generation 2912, once both RSA and ECDSA are READY_FOR_USE.
func getGenerationRSAECDSA() *cloudcertificates.GetGenerationResponse {
	return &cloudcertificates.GetGenerationResponse{
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{
				{
					AlgorithmInstanceCreatedBy:          "terraform-dev",
					AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:                 6123,
					AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
					AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:23:37Z")),
					CertificateStatus:                   "READY_FOR_USE",
					CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                             "ECDSA",
					SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
					SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:23:37Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:23:37Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"),
					SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:02"),
					SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:02"),
				},
				{
					AlgorithmInstanceCreatedBy:          "terraform-dev",
					AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
					AlgorithmInstanceID:                 6122,
					AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
					AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					CertificateStatus:                   "READY_FOR_USE",
					CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
					CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
					KeyType:                             "RSA",
					SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
					SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-08T13:06:53Z")),
					SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:06:53Z")),
					SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n"),
					SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01"),
					SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:01"),
				},
			},
			GenerationCreatedBy:    ptr.To("terraform-dev"),
			GenerationCreatedTime:  ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
			GenerationModifiedBy:   ptr.To("terraform-dev"),
			GenerationModifiedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-08T13:24:17Z")),
		},
		GenerationID:     2912,
		GenerationStatus: "READY_FOR_USE",
	}
}

// getGenerationRSAECDSAForGeneration is getGenerationRSAECDSA bound to a different generation ID, used for the
// new head generation CompleteLineage creates holding the merged, fully-signed algorithm set - real captured API
// behavior confirms completing a lineage doesn't touch the existing production generation at all; it instead
// creates a brand new, not-yet-promoted head generation with every algorithm signed.
func getGenerationRSAECDSAForGeneration(generationID int64) *cloudcertificates.GetGenerationResponse {
	generation := getGenerationRSAECDSA()
	generation.GenerationID = generationID
	return generation
}

// getLineageHeadOnly is the real captured GetLineage response for lineage 500001 expanded to just its head
// generation (2912), used by Create's head generation lookup.
func getLineageHeadOnly() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     2912,
			HeadGenerationStatus: "READY_FOR_USE",
		},
		LineageID: 500001,
	}
}

// getLineageProductionNoHead is the GetLineage response for a lineage fully promoted to production, RSA already
// live, and no head generation left: completing the lineage is the only route to add the missing ECDSA algorithm.
func getLineageProductionNoHead() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     4004,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// completeLineageECDSARequest is the CompleteLineage request for adding the missing ECDSA algorithm to
// getLineageProductionNoHead's current production generation.
func completeLineageECDSARequest() cloudcertificates.CompleteLineageRequest {
	return cloudcertificates.CompleteLineageRequest{
		LineageID: 500001,
		Body: cloudcertificates.CompleteLineageRequestBody{
			SourceGenerationID: ptr.To(int64(4004)),
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmECDSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

// completeLineageECDSARequestAck is completeLineageECDSARequest with AcknowledgeWarnings set, for the
// acknowledge_warnings=true variant of the same create-via-complete scenario.
func completeLineageECDSARequestAck() cloudcertificates.CompleteLineageRequest {
	req := completeLineageECDSARequest()
	req.AcknowledgeWarnings = true
	return req
}

// completeLineageECDSAResponse is the CompleteLineage response for completeLineageECDSARequest. Real captured API
// behavior: the existing current production generation (4004) is left untouched (still only RSA signed) - the
// completed, fully-signed algorithm set instead lands on a brand new head generation (4006).
func completeLineageECDSAResponse() *cloudcertificates.CompleteLineageResponse {
	return &cloudcertificates.CompleteLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     4006,
			HeadGenerationStatus: "READY_FOR_USE",
		},
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     4004,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// getLineageHeadFromComplete is the GetLineage response after completing getLineageProductionNoHead: the lineage
// now has both the new head generation (4006, from completing) and its original, untouched current production
// generation (4004).
func getLineageHeadFromComplete() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     4006,
			HeadGenerationStatus: "READY_FOR_USE",
		},
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     4004,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// getGenerationRSAOnlyPromoted is getGenerationRSAOnly once promoted to production on the given generation, but
// first_promoted_to_production_at is now set, so this generation can no longer take a plain upload for a new
// algorithm - only completing the lineage from a fresh Create can add one.
func getGenerationRSAOnlyPromoted(generationID int64) *cloudcertificates.GetGenerationResponse {
	generation := getGenerationRSAOnly()
	generation.GenerationID = generationID
	generation.GenerationStatus = "ACTIVE"
	generation.FirstPromotedToProductionTime = ptr.To(tst.NewTimeFromStringMust("2026-08-01T09:00:00Z"))
	return generation
}

// completeLineageRSAECDSARequest is the CompleteLineage request for adding the missing ECDSA algorithm to
// generation 2912 once it's the lineage's current production generation.
func completeLineageRSAECDSARequest() cloudcertificates.CompleteLineageRequest {
	return cloudcertificates.CompleteLineageRequest{
		LineageID: 500001,
		Body: cloudcertificates.CompleteLineageRequestBody{
			SourceGenerationID: ptr.To(int64(2912)),
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmECDSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nECDSACERT\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

// completeLineageRSAECDSAResponse is the CompleteLineage response for completeLineageRSAECDSARequest. Real
// captured API behavior: the existing current production generation (2912) is left untouched - the completed,
// fully-signed algorithm set instead lands on a brand new head generation (2914).
func completeLineageRSAECDSAResponse() *cloudcertificates.CompleteLineageResponse {
	return &cloudcertificates.CompleteLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     2914,
			HeadGenerationStatus: "READY_FOR_USE",
		},
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     2912,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// getLineageNoHead is the GetLineage response for a lineage that has been fully promoted and has no head
// generation left to upload a certificate to.
func getLineageNoHead() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		LineageID: 500001,
	}
}

// getLineageNewHead is getLineageHeadOnly with its head generation moved on to a new, different generation - as
// happens once a fresh renewal cycle creates a new draft generation for the lineage.
func getLineageNewHead() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     3050,
			HeadGenerationStatus: "CSR_READY",
		},
		LineageID: 500001,
	}
}

// getLineageProductionDrifted is the GetLineage response for a lineage with no head and a current production
// generation (9999) unrelated to generation 2912 this resource is bound to - e.g. after a rollback moved current
// production somewhere this resource never touched.
func getLineageProductionDrifted() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     9999,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// getLineageRenewedAfterPromotion is the GetLineage response once generation 2912 (RSA already live) has been
// promoted to production outside Terraform, then renewed: production keeps the original generation, and a brand
// new, empty head generation (5000) is created for the next certificate cycle.
func getLineageRenewedAfterPromotion() *cloudcertificates.GetLineageResponse {
	return &cloudcertificates.GetLineageResponse{
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     5000,
			HeadGenerationStatus: "CSR_READY",
		},
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     2912,
			ProductionGenerationStatus: "ACTIVE",
		},
		LineageID: 500001,
	}
}

// uploadRSARenewedRequestForGeneration is the UploadSignedCertificate request for uploading a freshly renewed RSA
// certificate (different content from the original, since renewal generates a new CSR) to the given generation.
func uploadRSARenewedRequestForGeneration(generationID int64) cloudcertificates.UploadSignedCertificateRequest {
	return cloudcertificates.UploadSignedCertificateRequest{
		LineageID:           500001,
		GenerationID:        generationID,
		AcknowledgeWarnings: false,
		Body: cloudcertificates.UploadSignedCertificateRequestBody{
			Algorithms: map[cloudcertificates.CryptographicAlgorithm]cloudcertificates.SignedCertificate{
				cloudcertificates.CryptographicAlgorithmRSA: {SignedCertificatePEM: "-----BEGIN CERTIFICATE-----\nRSACERTRENEWED\n-----END CERTIFICATE-----\n"},
			},
		},
	}
}

// getGenerationRSARenewedForGeneration is the post-upload GetGeneration detail for a renewed RSA certificate on
// the given generation.
func getGenerationRSARenewedForGeneration(generationID int64) *cloudcertificates.GetGenerationResponse {
	generation := getGenerationRSAOnly()
	generation.GenerationID = generationID
	generation.Algorithms[1].SignedCertificatePEM = ptr.To("-----BEGIN CERTIFICATE-----\nRSACERTRENEWED\n-----END CERTIFICATE-----\n")
	return generation
}

func TestUploadResource(t *testing.T) {
	t.Parallel()

	dsName := "akamai_cloudcertificates_upload.test"

	exhaustiveChecker := test.NewStateChecker(dsName).
		CheckEqualBatch("", test.AttributeBatch{
			"lineage_id":               "500001",
			"generation_id":            "2912",
			"acknowledge_warnings":     "false",
			"generation_status":        "READY_FOR_USE",
			"generation_created_by":    "terraform-dev",
			"generation_created_time":  "2026-07-07T10:18:46Z",
			"generation_modified_by":   "terraform-dev",
			"generation_modified_time": "2026-07-08T13:14:25Z",
			"algorithms.%":             "1",
		}).
		CheckMissing("first_promoted_to_production_at").
		CheckEqualBatch("algorithms.RSA.", test.AttributeBatch{
			"algorithm_instance_id":                    "6122",
			"algorithm_instance_created_by":            "terraform-dev",
			"algorithm_instance_created_time":          "2026-07-07T10:18:46Z",
			"algorithm_instance_modified_by":           "terraform-dev",
			"algorithm_instance_modified_time":         "2026-07-08T13:06:53Z",
			"certificate_status":                       "READY_FOR_USE",
			"csr_expiration_date":                      "2027-07-07T10:18:46Z",
			"csr_pem":                                  "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
			"signed_certificate_issuer":                "CN=Test Certificate Authority",
			"signed_certificate_not_valid_before_date": "2026-07-08T13:06:53Z",
			"signed_certificate_not_valid_after_date":  "2027-07-08T13:06:53Z",
			"signed_certificate_pem":                   "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n",
			"signed_certificate_serial_number":         "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01",
			"signed_certificate_sha256_fingerprint":    "12:34:56:78:9A:BC:DE:01",
		}).
		CheckMissing("algorithms.RSA.trust_chain_pem").
		Build()

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - create, draft generation, single RSA algorithm": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  exhaustiveChecker,
				},
			},
		},
		"happy path - create with acknowledge_warnings configured to true": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequestAck()).
					Return(uploadRSAResponse(), nil).Once()
				// read (post-apply refresh, and pre-destroy refresh)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ack.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("acknowledge_warnings", "true").Build(),
				},
			},
		},
		"happy path - acknowledge_warnings changes in place, no replace or re-upload": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(7)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read: post-apply refresh, pre-plan refresh before step 2, Update's own refetch, pre-destroy refresh
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(5)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("acknowledge_warnings", "false").Build(),
				},
				{
					// only acknowledge_warnings changes: no algorithm content differs, so this must absorb in
					// place (Update, no UploadSignedCertificate call) rather than replace or re-upload.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ack.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(dsName, plancheck.ResourceActionUpdate),
						},
					},
					Check: test.NewStateChecker(dsName).CheckEqual("acknowledge_warnings", "true").Build(),
				},
			},
		},
		"happy path - create, lineage already promoted with RSA, completes with ECDSA": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageProductionNoHead(), nil).Once()
				// generation 4004's data: RSA already signed with matching content (so Create must filter it out
				// of the CompleteLineage call, not reject it); ECDSA is still missing there.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4004}).
					Return(getGenerationRSAOnlyPromoted(4004), nil).Once()
				m.On("CompleteLineage", testutils.MockContext, completeLineageECDSARequest()).
					Return(completeLineageECDSAResponse(), nil).Once()
				// read: Create's own post-completion refetch of the new head generation, plus every later
				// ModifyPlan/refresh check against it.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4006}).
					Return(getGenerationRSAECDSAForGeneration(4006), nil).Times(2)
				// read: every later ModifyPlan check confirming the head hasn't moved on
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadFromComplete(), nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("generation_id", "4006").
						CheckEqual("generation_status", "READY_FOR_USE").
						CheckMissing("first_promoted_to_production_at").
						CheckEqual("algorithms.%", "2").
						CheckEqual("algorithms.RSA.certificate_status", "READY_FOR_USE").
						CheckEqual("algorithms.ECDSA.certificate_status", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"happy path - create, lineage already promoted with RSA, completes with ECDSA, acknowledge_warnings=true": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageProductionNoHead(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4004}).
					Return(getGenerationRSAOnlyPromoted(4004), nil).Once()
				// acknowledge_warnings=true must be forwarded to CompleteLineage - if this resource only ever
				// passed it to UploadSignedCertificate, this mock would never match and the test would panic.
				m.On("CompleteLineage", testutils.MockContext, completeLineageECDSARequestAck()).
					Return(completeLineageECDSAResponse(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4006}).
					Return(getGenerationRSAECDSAForGeneration(4006), nil).Times(2)
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadFromComplete(), nil).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa_ack.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("generation_id", "4006").
						CheckEqual("acknowledge_warnings", "true").
						CheckEqual("algorithms.%", "2").
						Build(),
				},
			},
		},
		"happy path - create then update, draft generation, add ECDSA algorithm": {
			init: func(m *cloudcertificates.Mock) {
				// create, plus every later ModifyPlan check confirming the head generation hasn't moved on
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(7)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read (post-apply refresh after create, and pre-plan refresh(es) before the update)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
				// update: only the newly-added ECDSA algorithm is sent, not the already-uploaded RSA one
				m.On("UploadSignedCertificate", testutils.MockContext, uploadECDSARequest()).
					Return(uploadRSAECDSAResponse(), nil).Once()
				// read (Update's own re-fetch, plus the pre-destroy refresh)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAECDSA(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("algorithms.%", "1").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							// RSA is untouched by this update: its Computed sub-attributes must stay known at
							// their prior values, not get replanned as unknown just because ECDSA was added to
							// the same algorithms map.
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_id"), knownvalue.Int64Exact(6122)),
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_created_by"), knownvalue.StringExact("terraform-dev")),
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_created_time"), knownvalue.StringExact("2026-07-07T10:18:46Z")),
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("certificate_status"), knownvalue.StringExact("READY_FOR_USE")),
							// ECDSA is brand new to this resource: it has no prior state to preserve, so its
							// Computed sub-attributes are correctly unknown until the upload completes.
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("ECDSA").AtMapKey("algorithm_instance_id")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("ECDSA").AtMapKey("certificate_status")),
						},
					},
					Check: test.NewStateChecker(dsName).
						CheckEqual("algorithms.%", "2").
						CheckEqual("generation_status", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"happy path - lineage head is promoted away, resource keeps managing its bound generation": {
			init: func(m *cloudcertificates.Mock) {
				// The very first call is Create's own head lookup, which must see a real head generation to
				// bind to; every later call simulates the lineage having since been fully promoted - no head
				// left, and current production now points at the very generation this resource is bound to.
				resp := getLineageHeadOnly()
				calls := 0
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Run(func(mock.Arguments) {
					calls++
					if calls > 1 {
						resp.Head = nil
						resp.CurrentProduction = &cloudcertificates.ProductionGeneration{
							ProductionGenerationID:     2912,
							ProductionGenerationStatus: "ACTIVE",
						}
					}
				}).Times(5)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read: post-apply refresh, pre-plan refresh for the second (PlanOnly) step, and the pre-destroy refresh
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// no drift and no config change: the plan is empty even though the lineage's head has moved on.
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - lineage head moves to a new generation, resource is replaced": {
			init: func(m *cloudcertificates.Mock) {
				// The first few calls (Create's own lookup, plus the post-apply drift checks after step 1) see
				// the original head generation 2912; from then on (step 2's plan onward) the lineage's head has
				// moved on to a new generation 3050, simulating a fresh renewal cycle.
				resp := getLineageHeadOnly()
				calls := 0
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Run(func(mock.Arguments) {
					calls++
					if calls > 3 {
						resp.Head = getLineageNewHead().Head
					}
				}).Times(7)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read: post-apply refresh after create, plus the pre-plan refresh before the second step
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
				// replace's Create: pre-upload check on the new head generation 3050
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 3050}).
					Return(getGenerationNoneSignedForGeneration(3050), nil).Once()
				// replace: the new resource instance binds to the new head generation 3050
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequestForGeneration(3050)).
					Return(uploadRSAResponseForGeneration(3050), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 3050}).
					Return(getGenerationRSAOnlyForGeneration(3050), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "3050").Build(),
				},
			},
		},
		"happy path - lineage promoted then renewed outside terraform, resource is replaced onto the new head": {
			init: func(m *cloudcertificates.Mock) {
				// QA-reported regression: promoting outside Terraform, then renewing (generating a fresh head
				// with new CSRs) must be detected as drift and replace the resource - not rejected as an
				// "Algorithm Change Not Supported" write-once violation, even though the renewed config's RSA
				// content necessarily differs from the original (a renewal always re-signs from a new CSR).
				resp := getLineageHeadOnly()
				calls := 0
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Run(func(mock.Arguments) {
					calls++
					if calls > 3 {
						renewed := getLineageRenewedAfterPromotion()
						resp.Head = renewed.Head
						resp.CurrentProduction = renewed.CurrentProduction
					}
				}).Times(7)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// generation 2912's data: starts as an unpromoted draft, then - simulating the external promote
				// sometime between step 1 and step 2 - gains a first_promoted_to_production_at; it never gets
				// the renewed certificate, since that lands on the new head generation instead.
				genResp := getGenerationRSAOnly()
				genCalls := 0
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(genResp, nil).Run(func(mock.Arguments) {
					genCalls++
					if genCalls >= 3 {
						genResp.GenerationStatus = "ACTIVE"
						genResp.FirstPromotedToProductionTime = ptr.To(tst.NewTimeFromStringMust("2026-08-01T09:00:00Z"))
					}
				}).Times(3)
				// replace's Create: pre-upload check on the new, empty head generation created by the renewal
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 5000}).
					Return(getGenerationNoneSignedForGeneration(5000), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARenewedRequestForGeneration(5000)).
					Return(uploadRSAResponseForGeneration(5000), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 5000}).
					Return(getGenerationRSARenewedForGeneration(5000), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_renewed.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("generation_id", "5000").
						CheckMissing("first_promoted_to_production_at").
						CheckEqual("algorithms.RSA.certificate_status", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"happy path - generation promoted before update, adding an algorithm replaces and completes the lineage": {
			init: func(m *cloudcertificates.Mock) {
				// create, bound to draft head generation 2912; then, once genuinely promoted outside Terraform
				// (both the lineage's own pointers and generation 2912's own metadata reflect it), adding ECDSA
				// forces a replace that completes the lineage instead of trying (and failing) to Update it.
				// promoted/completed are driven by the actual events (the generation gaining a promotion
				// timestamp, CompleteLineage running) rather than by guessing exact call counts.
				var promoted, completed bool
				lineageResp := getLineageHeadOnly()
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(lineageResp, nil).Run(func(mock.Arguments) {
					switch {
					case completed:
						lineageResp.Head = &cloudcertificates.HeadGeneration{HeadGenerationID: 2914, HeadGenerationStatus: "READY_FOR_USE"}
					case promoted:
						lineageResp.Head = nil
						lineageResp.CurrentProduction = &cloudcertificates.ProductionGeneration{
							ProductionGenerationID:     2912,
							ProductionGenerationStatus: "ACTIVE",
						}
					}
				}).Times(7)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// generation 2912's data: starts as an unpromoted draft, then - simulating promotion to
				// production sometime between step 1 and step 2 - gains a first_promoted_to_production_at, and
				// finally (once CompleteLineage runs, during the replace's own Create) gains the ECDSA algorithm.
				genResp := getGenerationRSAOnly()
				genCalls := 0
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(genResp, nil).Run(func(mock.Arguments) {
					genCalls++
					if genCalls >= 3 {
						genResp.GenerationStatus = "ACTIVE"
						genResp.FirstPromotedToProductionTime = ptr.To(tst.NewTimeFromStringMust("2026-08-01T09:00:00Z"))
						promoted = true
					}
				}).Times(4)
				m.On("CompleteLineage", testutils.MockContext, completeLineageRSAECDSARequest()).
					Return(completeLineageRSAECDSAResponse(), nil).Run(func(mock.Arguments) {
					completed = true
				}).Once()
				// read: Create's own post-completion refetch of the new head generation, plus every later
				// ModifyPlan/refresh check against it.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2914}).
					Return(getGenerationRSAECDSAForGeneration(2914), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction(dsName, plancheck.ResourceActionReplace),
							// generation_id itself drives the replace: the resource is now bound to a
							// not-yet-created generation, so it can only be unknown.
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("generation_id")),
							// unlike a plain update to the same generation, RSA's already-signed metadata is not
							// guaranteed to carry over unchanged onto the new generation CompleteLineage creates -
							// only the certificate content itself is guaranteed to be reused, not e.g. its
							// algorithm_instance_id - so these must stay unknown too, not be copied forward.
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_id")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_created_by")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("algorithm_instance_created_time")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("signed_certificate_issuer")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("signed_certificate_not_valid_before_date")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("signed_certificate_not_valid_after_date")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("signed_certificate_serial_number")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("RSA").AtMapKey("signed_certificate_sha256_fingerprint")),
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("algorithms").AtMapKey("ECDSA").AtMapKey("algorithm_instance_id")),
						},
					},
					Check: test.NewStateChecker(dsName).
						CheckEqual("generation_id", "2914").
						CheckEqual("algorithms.%", "2").
						Build(),
				},
			},
		},
		"happy path - plan-only with unknown algorithms map": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// Step 2: Use terraform_data.output to make the entire algorithms map unknown during plan.
					// This should return early in ModifyPlan without evaluating map elements or throwing errors.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/unknown_algorithms_map.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - plan-only with unknown certificate nested value": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(5)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// Step 2: Use terraform_data.output to make only the nested certificate PEM unknown.
					// diffAlgorithms should defer content checks for unknown values, planning with no errors.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/unknown_certificate_nested_value.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - read warns and removes state, lineage itself is gone, nothing left to import": {
			init: func(m *cloudcertificates.Mock) {
				// Read only ever runs once for an unchanged config: right after apply, as the built-in
				// post-apply consistency check. By then the whole lineage is gone (not just the generation),
				// so there's no import target to even look for: Read must warn and drop the resource directly,
				// without handleGenerationGone's usual lineage lookup.
				// the post-apply consistency check runs ModifyPlan once more after Read removes the resource
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Twice()
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// Create's own refetch
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Once()
				// Read's one and only invocation (the post-apply consistency check): the lineage is gone
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					// AddWarning (unlike the AddError this replaced) doesn't roll back RemoveResource, so the
					// resource is genuinely gone from state and the plan is left non-empty: config still wants it.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:              test.NewStateChecker(dsName).CheckEqual("algorithms.%", "1").Build(),
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - read warns and removes state, nothing left to import, only the generation is gone": {
			init: func(m *cloudcertificates.Mock) {
				resp := getLineageHeadOnly()
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Times(3)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// Create's own refetch
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Once()
				// Read's one and only invocation (the post-apply consistency check): generation gone, and (an
				// artificial, since a still-existing lineage always has a head or current production - the API
				// refuses to delete the last one) empty lineage response simulates the whole lineage being gone
				// too: nothing left to import, so Read falls back to warning and dropping the resource.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(nil, cloudcertificates.ErrGenerationNotFound).Run(func(mock.Arguments) {
					resp.Head = nil
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:              test.NewStateChecker(dsName).CheckEqual("algorithms.%", "1").Build(),
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - generation gone, lineage has a new head to import, resource removed with a warning then re-imported": {
			init: func(m *cloudcertificates.Mock) {
				resp := getLineageHeadOnly()
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Times(6)

				// 1. Pre-upload check during Create
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationNoneSigned(), nil).Once()

				// 2. Upload signed certificate during Create
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()

				// 3. Create's post-upload refetch
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationRSAOnly(), nil).Once()

				// 4. Step 1's own post-apply consistency check: nothing has drifted yet, so this succeeds normally.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationRSAOnly(), nil).Once()

				// 5. Step 2's terraform refresh: generation 2912 is gone, and the lineage head moves to 3050;
				// AddWarning (unlike the AddError this replaced) doesn't roll back RemoveResource, and unlike the
				// implicit post-apply consistency check, RefreshState durably persists the removal.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(nil, cloudcertificates.ErrGenerationNotFound).Run(func(mock.Arguments) {
					resp.Head = getLineageNewHead().Head
				}).Once()

				// 6. Import onto the new head generation 3050, plus its own post-import and idempotency-check
				// refreshes: 3050 already carries the same RSA certificate config declares, so nothing drifts.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 3050,
				}).Return(getGenerationRSAOnlyForGeneration(3050), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// terraform refresh durably persists the removal, unlike the no-op post-apply consistency
					// check plan of step 1 (which discards its refreshed state at the end of that step).
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
				{
					// re-import using the new head generation the warning pointed at
					ResourceName:       "akamai_cloudcertificates_upload.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:        true,
					ImportStateId:      "500001,3050",
					ImportStatePersist: true,
					ImportStateCheck:   test.NewImportChecker().CheckEqual("generation_id", "3050").Build(),
				},
				{
					// verifies the re-imported state is stable and matches config with no further drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly: true,
				},
			},
		},
		"expect error - update succeeds but post-upload refetch fails": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(5)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read (post-apply refresh after create, and pre-plan refresh(es) before the update)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Times(3)
				// update: uploads the newly-added ECDSA algorithm successfully
				m.On("UploadSignedCertificate", testutils.MockContext, uploadECDSARequest()).
					Return(uploadRSAECDSAResponse(), nil).Once()
				// Update's own post-upload refetch fails (e.g. a transient network blip): the certificate is
				// already uploaded by this point, so the error must advise import as the recovery path, same as
				// Create's equivalent branch.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(nil, errors.New("simulated transient failure")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("algorithms.%", "1").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					ExpectError: tst.ErrPattern(
						"The certificate was uploaded to generation 2912 of lineage 500001, but reading it back failed: " +
							"not retrying due to non-retriable error: simulated transient failure\\. " +
							"To recover, import the upload resource: terraform import <resource address> 500001,2912",
					),
				},
			},
		},
		"expect error - current production drifts to an unrelated generation, resource is unrecognized": {
			init: func(m *cloudcertificates.Mock) {
				// No head and current production points at a generation (9999) this resource never touched -
				// e.g. a rollback moved production elsewhere. Neither matches our bound generation 2912, so
				// unlike a head simply moving on, this can't be silently replaced: 9999 may already carry real
				// signed content unrelated to this resource instance. ModifyPlan must error and point the user
				// at import instead.
				resp := getLineageHeadOnly()
				calls := 0
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(resp, nil).Run(func(mock.Arguments) {
					calls++
					if calls > 3 {
						drifted := getLineageProductionDrifted()
						resp.Head = drifted.Head
						resp.CurrentProduction = drifted.CurrentProduction
					}
				}).Times(4)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read: post-apply refresh, plus the pre-plan refresh before the second (PlanOnly) step
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// no config change, but the bound generation is neither the (nonexistent) head nor current
					// production anymore: ModifyPlan must refuse to guess and error instead of replacing.
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly:    true,
					ExpectError: tst.ErrPattern(`Generation 2912 of lineage 500001 is no longer the lineage's head or current production generation`),
				},
			},
		},
		"expect error - modify plan fails when the lineage head check errors": {
			init: func(m *cloudcertificates.Mock) {
				// Create's own head lookup, plus the post-apply consistency check's ModifyPlan run succeed as
				// usual; the ModifyPlan run driving step 2's expected error is the one that fails.
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(nil, errors.New("simulated network error")).Once()
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// read: post-apply refresh, plus the pre-plan refresh before the second (PlanOnly) step
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("generation_id", "2912").Build(),
				},
				{
					// no config change: only ModifyPlan's own lineage head check fails, which must now hard-error
					// and abort the plan rather than merely warn and silently proceed without drift detection.
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly:    true,
					ExpectError: tst.ErrPattern(`Could not determine whether lineage 500001's head generation has moved on: simulated network error\.`),
				},
			},
		},
		"expect error - create fails, algorithm already exists on head generation": {
			init: func(m *cloudcertificates.Mock) {
				// e.g. the resource was removed then re-added to config: the head generation already carries an
				// accepted RSA certificate, so Create must refuse to blindly re-upload and tell the user to import.
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAOnly(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`(?s)Generation 2912 of lineage 500001 already has an accepted signed certificate for the RSA algorithm.*terraform import`),
				},
			},
		},
		"expect error - create fails, head generation already has a different algorithm signed": {
			init: func(m *cloudcertificates.Mock) {
				// config only declares RSA, but the head generation already has an unrelated ECDSA certificate
				// signed - Create must still refuse, since a fresh head generation is never partially signed.
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Once()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationECDSAOnly(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`(?s)Generation 2912 of lineage 500001 already has an accepted signed certificate for the ECDSA algorithm.*terraform import`),
				},
			},
		},
		"expect error - create fails, lineage has no head or production generation": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageNoHead(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`Lineage 500001 has no head or current production generation`),
				},
			},
		},
		"expect error - create fails, current production generation already has an unconfigured algorithm": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageProductionNoHead(), nil).Once()
				// generation 4004 already has RSA signed, but config below only declares ECDSA: config no longer
				// matches reality, so Create must refuse and suggest import rather than silently completing.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4004}).
					Return(getGenerationRSAOnlyPromoted(4004), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_ecdsa.tf"),
					ExpectError: tst.ErrPattern(`(?s)Generation 4004 of lineage 500001 already has an accepted signed certificate for the RSA algorithm.*terraform import`),
				},
			},
		},
		"expect error - create fails, algorithm already signed with different content during complete lineage": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageProductionNoHead(), nil).Once()
				// generation 4004's RSA is already signed with "RSACERT"; config redundantly declares RSA with
				// different content, which must be rejected before CompleteLineage is ever called.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 4004}).
					Return(getGenerationRSAOnlyPromoted(4004), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa_modified.tf"),
					ExpectError: tst.ErrPattern(`the RSA algorithm instance already has an accepted signed certificate`),
				},
			},
		},
		"expect error - update fails, removing algorithm not supported": {
			init: func(m *cloudcertificates.Mock) {
				// create, plus every later ModifyPlan check confirming the head generation hasn't moved on
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(4)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSAECDSARequest()).
					Return(uploadRSAECDSAResponse(), nil).Once()
				// read (post-apply refresh, and pre-plan refresh before the update)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAECDSA(), nil).Twice()
				// delete (test cleanup)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationRSAECDSA(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("algorithms.%", "2").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`Removing the ECDSA algorithm from an existing upload is not supported`),
				},
			},
		},
		"expect error - create fails, generation immutable": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Once()
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(nil, cloudcertificates.ErrGenerationImmutable).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`Generation 2912 of lineage 500001 has already been promoted`),
				},
			},
		},
		"expect error - create succeeds but post-upload refetch fails, resource persisted with placeholders": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Once()
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
				// post-upload refetch fails (e.g. a transient network blip): the certificate is already uploaded
				// by this point, so Create must still persist a placeholder state instead of orphaning it.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(nil, errors.New("simulated transient failure")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(`The certificate was uploaded to generation 2912 of lineage 500001, but reading it back failed`),
				},
			},
		},
		"expect error - post-upload generation refetch fails followed by import recovery": {
			init: func(m *cloudcertificates.Mock) {
				// GetLineage: Called 1 time in Step 1 (Create head lookup) and 2 times in Step 3 (ModifyPlan checks during PlanOnly)
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)

				// 1. Step 1 Pre-upload check: verify target generation has no signed cert yet
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationNoneSigned(), nil).Once()

				// 2. Step 1 Upload: signed certificate payload succeeds on server
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()

				// 3. Step 1 Post-upload refetch: fails (simulating a transient network error)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(nil, errors.New("simulated network error")).Once()

				// 4. Step 2 & Step 3 Read calls:
				// - 1 call during Step 2 (ImportState -> Read)
				// - 1 call during Step 3 (PlanOnly pre-plan refresh)
				// Both return the server state where the certificate upload succeeded.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationRSAOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					// Step 1: Attempt creation, expect diagnostic error advising import recovery
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(
						"The certificate was uploaded to generation 2912 of lineage 500001, but reading it back failed: " +
							"not retrying due to non-retriable error: simulated network error. " +
							"To recover, import the newly created upload resource: terraform import <resource address> 500001,2912",
					),
				},
				{
					// Step 2: Execute import state recovery using the exact ID advised in Step 1
					ResourceName:       "akamai_cloudcertificates_upload.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:        true,
					ImportStatePersist: true,
					ImportStateId:      "500001,2912",
					Check: test.NewStateChecker("akamai_cloudcertificates_upload.test").
						CheckEqual("lineage_id", "500001").
						CheckEqual("generation_id", "2912").
						Build(),
				},
				{
					// Step 3: Verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly: true,
				},
			},
		},
		"validation error - generation_id is configured": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/generation_id_configured.tf"),
					ExpectError: tst.ErrPattern(`Cannot set value for this attribute as the provider has marked it as read-only`),
				},
			},
		},
		"validation error - missing lineage_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/missing_lineage_id.tf"),
					ExpectError: tst.ErrPattern(`The argument "lineage_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - lineage_id is 0": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/zero_lineage_id.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - missing algorithms": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/missing_algorithms.tf"),
					ExpectError: tst.ErrPattern(`The argument "algorithms" is required, but no definition was found.`),
				},
			},
		},
		"validation error - too many algorithms": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/too_many_algorithms.tf"),
					ExpectError: tst.ErrPattern(`Attribute algorithms map must contain at least 1 elements and at most 2 elements, got: 3`),
				},
			},
		},
		"validation error - invalid key_type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/invalid_key_type.tf"),
					ExpectError: tst.ErrPattern(`algorithms\["INVALID_KEY_TYPE"\] value must be one of: \["RSA" "ECDSA"\], got: "INVALID_KEY_TYPE"`),
				},
			},
		},
		"validation error - empty signed_certificate_pem": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/empty_signed_certificate_pem.tf"),
					ExpectError: tst.ErrPattern(`signed_certificate_pem string length must be at least 1, got: 0`),
				},
			},
		},
		"validation error - malformed signed_certificate_pem": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/invalid_signed_certificate_pem.tf"),
					ExpectError: tst.ErrPattern(`signed_certificate_pem must be in PEM format`),
				},
			},
		},
		"validation error - malformed trust_chain_pem": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/invalid_trust_chain_pem.tf"),
					ExpectError: tst.ErrPattern(`trust_chain_pem must be in PEM format`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CloudCertificates)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

// TestUploadResourcePolling exercises waitForSignedAlgorithms: a GET immediately following a successful upload can
// still report the just-uploaded algorithm as unsigned while the server finishes processing it. It uses a custom
// polling subprovider so these cases run against a short interval/timeout instead of the production defaults.
func TestUploadResourcePolling(t *testing.T) {
	t.Parallel()

	dsName := "akamai_cloudcertificates_upload.test"

	tests := map[string]struct {
		init         func(*cloudcertificates.Mock)
		pollInterval time.Duration
		pollTimeout  time.Duration
		steps        []resource.TestStep
	}{
		"happy path - create polls until the uploaded algorithm shows signed-certificate details": {
			pollInterval: 10 * time.Millisecond,
			pollTimeout:  200 * time.Millisecond,
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(3)
				// pre-upload check: neither configured algorithm already has an accepted signed certificate
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()

				var uploaded bool
				pollCalls := 0
				genResp := getGenerationNoneSigned()
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(genResp, nil).Run(func(mock.Arguments) {
					if !uploaded {
						return
					}
					pollCalls++
					if pollCalls >= 2 {
						ready := getGenerationRSAOnly()
						genResp.Algorithms = ready.Algorithms
						genResp.GenerationStatus = ready.GenerationStatus
						genResp.GenerationCreatedBy = ready.GenerationCreatedBy
						genResp.GenerationCreatedTime = ready.GenerationCreatedTime
						genResp.GenerationModifiedBy = ready.GenerationModifiedBy
						genResp.GenerationModifiedTime = ready.GenerationModifiedTime
					}
				})

				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Run(func(mock.Arguments) {
					uploaded = true
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("algorithms.%", "1").
						CheckEqual("algorithms.RSA.certificate_status", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"happy path - update polls until the newly uploaded algorithm shows signed-certificate details": {
			pollInterval: 10 * time.Millisecond,
			pollTimeout:  200 * time.Millisecond,
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Times(7)
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(getGenerationNoneSigned(), nil).Once()
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()

				var uploaded bool
				pollCalls := 0
				genResp := getGenerationRSAOnly()
				// every read before the ECDSA upload correctly still shows only RSA signed; every poll read
				// after it keeps showing the same until the second attempt, simulating processing finishing.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					Return(genResp, nil).Run(func(mock.Arguments) {
					if !uploaded {
						return
					}
					pollCalls++
					if pollCalls >= 2 {
						ready := getGenerationRSAECDSA()
						genResp.Algorithms = ready.Algorithms
						genResp.GenerationStatus = ready.GenerationStatus
						genResp.GenerationModifiedBy = ready.GenerationModifiedBy
						genResp.GenerationModifiedTime = ready.GenerationModifiedTime
					}
				})

				m.On("UploadSignedCertificate", testutils.MockContext, uploadECDSARequest()).
					Return(uploadRSAECDSAResponse(), nil).Run(func(mock.Arguments) {
					uploaded = true
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("algorithms.%", "1").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ecdsa.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("algorithms.%", "2").
						CheckEqual("algorithms.ECDSA.certificate_status", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"expect error - create times out waiting for the uploaded algorithm to finish processing": {
			pollInterval: 10 * time.Millisecond,
			pollTimeout:  50 * time.Millisecond,
			init: func(m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID:         500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
				}).Return(getLineageHeadOnly(), nil).Once()
				// the pre-upload check and every poll attempt alike keep seeing the RSA algorithm unsigned:
				// processing never finishes within the configured timeout. The exact number of poll attempts
				// that fit before the deadline is timing-dependent, so no fixed call count is asserted here.
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: 500001, GenerationID: 2912}).
					// The number of these calls is timing-dependent.
					Return(getGenerationNoneSigned(), nil)
				m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
					Return(uploadRSAResponse(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ExpectError: tst.ErrPattern(
						`The certificate was uploaded to generation 2912 of lineage 500001, but reading it back failed: ` +
							`context terminated (while waiting to retry|before function execution): context deadline exceeded`,
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CloudCertificates)
			}

			config := defaultSubproviderConfig()
			config.upload.pollInterval = tc.pollInterval
			config.upload.pollTimeout = tc.pollTimeout

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func TestUploadResourceImport(t *testing.T) {
	t.Parallel()

	baseChecker := test.NewImportChecker().
		CheckEqual("lineage_id", "500001").
		CheckEqual("generation_id", "2912").
		CheckEqual("acknowledge_warnings", "false").
		CheckEqual("generation_status", "READY_FOR_USE").
		CheckEqual("generation_created_by", "terraform-dev").
		CheckEqual("generation_created_time", "2026-07-07T10:18:46Z").
		CheckEqual("generation_modified_by", "terraform-dev").
		CheckEqual("generation_modified_time", "2026-07-08T13:14:25Z").
		CheckMissing("first_promoted_to_production_at").
		CheckEqual("algorithms.%", "1").
		CheckEqual("algorithms.RSA.algorithm_instance_id", "6122").
		CheckEqual("algorithms.RSA.algorithm_instance_created_by", "terraform-dev").
		CheckEqual("algorithms.RSA.algorithm_instance_created_time", "2026-07-07T10:18:46Z").
		CheckEqual("algorithms.RSA.algorithm_instance_modified_by", "terraform-dev").
		CheckEqual("algorithms.RSA.algorithm_instance_modified_time", "2026-07-08T13:06:53Z").
		CheckEqual("algorithms.RSA.certificate_status", "READY_FOR_USE").
		CheckEqual("algorithms.RSA.csr_expiration_date", "2027-07-07T10:18:46Z").
		CheckEqual("algorithms.RSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
		CheckEqual("algorithms.RSA.signed_certificate_issuer", "CN=Test Certificate Authority").
		CheckEqual("algorithms.RSA.signed_certificate_not_valid_after_date", "2027-07-08T13:06:53Z").
		CheckEqual("algorithms.RSA.signed_certificate_not_valid_before_date", "2026-07-08T13:06:53Z").
		CheckEqual("algorithms.RSA.signed_certificate_pem", "-----BEGIN CERTIFICATE-----\nRSACERT\n-----END CERTIFICATE-----\n").
		CheckEqual("algorithms.RSA.signed_certificate_serial_number", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:01").
		CheckEqual("algorithms.RSA.signed_certificate_sha256_fingerprint", "12:34:56:78:9A:BC:DE:01").
		CheckMissing("algorithms.RSA.trust_chain_pem")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - import by lineage_id and generation_id": {
			init: func(m *cloudcertificates.Mock) {
				// Read is called once during ImportState and once during PlanOnly refresh
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationRSAOnly(), nil).Twice()

				// ModifyPlan is called during the PlanOnly step to verify the generation hasn't moved
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{
						cloudcertificates.ExpandGenerationsHead,
						cloudcertificates.ExpandGenerationsCurrentProduction,
					},
				}).Return(getLineageHeadOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_upload.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:        true,
					ImportStateId:      "500001,2912",
					ImportStatePersist: true,
					ImportStateCheck:   baseChecker.Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import with acknowledge_warnings": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500001,
					GenerationID: 2912,
				}).Return(getGenerationRSAOnly(), nil).Twice()

				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500001,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{
						cloudcertificates.ExpandGenerationsHead,
						cloudcertificates.ExpandGenerationsCurrentProduction,
					},
				}).Return(getLineageHeadOnly(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_upload.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ack.tf"),
					ImportState:        true,
					ImportStateId:      "500001,2912,true",
					ImportStatePersist: true,
					ImportStateCheck: baseChecker.
						CheckEqual("acknowledge_warnings", "true").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa_ack.tf"),
					PlanOnly: true,
				},
			},
		},
		"error - acknowledge_warnings is not a boolean": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001,2912,notabool",
					ExpectError:   tst.ErrPattern(`expected a boolean acknowledge_warnings, got: "notabool"`),
				},
			},
		},
		"error - lineage_id is not numeric": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "abc,2912",
					ExpectError:   tst.ErrPattern(`(?s)Invalid Import ID.*expected a numeric lineage_id, got: "abc"`),
				},
			},
		},
		"error - generation_id is not numeric": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001,abc",
					ExpectError:   tst.ErrPattern(`(?s)Invalid Import ID.*expected a numeric generation_id, got: "abc"`),
				},
			},
		},
		"error - too few import ID parts": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001",
					ExpectError: tst.ErrPattern(
						`(?s)Incorrect import ID.*invalid number of importID parts: 1; you need to provide an importID in the format 'lineageID,generationID\[,acknowledgeWarnings\]'`,
					),
				},
			},
		},
		"error - too many import ID parts": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001,2912,true,extra",
					ExpectError: tst.ErrPattern(
						`(?s)Incorrect import ID.*invalid number of importID parts: 4; you need to provide an importID in the format 'lineageID,generationID\[,acknowledgeWarnings\]'`,
					),
				},
			},
		},
		"error - lineage_id is zero": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "0,2912",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: 0`),
				},
			},
		},
		"error - lineage_id is negative": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "-5,2912",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: -5`),
				},
			},
		},
		"error - generation_id is zero": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001,0",
					ExpectError:   tst.ErrPattern(`generation_id must be greater than 0, got: 0`),
				},
			},
		},
		"error - generation_id is negative": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_upload.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesUpload/create_rsa.tf"),
					ImportState:   true,
					ImportStateId: "500001,-10",
					ExpectError:   tst.ErrPattern(`generation_id must be greater than 0, got: -10`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CloudCertificates)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func TestUploadCertificateErrorDiagnostics(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err         error
		wantSummary string
		wantDetail  string
	}{
		"generation immutable": {
			err:         cloudcertificates.ErrGenerationImmutable,
			wantSummary: "Generation Immutable",
			wantDetail:  "Generation 2912 of lineage 500001 has already been promoted to production and can no longer be modified directly.",
		},
		"certificate already uploaded": {
			err:         cloudcertificates.ErrCertAlreadyUploaded,
			wantSummary: "Certificate Already Uploaded",
			wantDetail:  "A signed certificate was already uploaded for this algorithm instance.",
		},
		"lineage not found": {
			err:         cloudcertificates.ErrLineageNotFound,
			wantSummary: "Lineage Not Found",
			wantDetail:  "No certificate lineage found with ID 500001.",
		},
		"generation not found": {
			err:         cloudcertificates.ErrCertificateNotFound,
			wantSummary: "Generation Not Found",
			wantDetail:  "No generation 2912 found for lineage 500001.",
		},
		"certificate expiry invalid": {
			err:         cloudcertificates.ErrCertExpiryInvalid,
			wantSummary: "Invalid Certificate Validity",
			wantDetail:  "The uploaded certificate's validity dates are invalid.",
		},
		"certificate parse error": {
			err:         cloudcertificates.ErrCertParseError,
			wantSummary: "Certificate Parse Error",
			wantDetail:  "The uploaded certificate is malformed or could not be parsed.",
		},
		"certificate/CSR mismatch": {
			err:         cloudcertificates.ErrCertCSRMismatch,
			wantSummary: "Certificate/CSR Mismatch",
			wantDetail:  "The uploaded certificate's public key does not match the stored CSR.",
		},
		"unknown key type": {
			err:         cloudcertificates.ErrUnknownKeyType,
			wantSummary: "Unknown Key Type",
			wantDetail:  "The given key type is not part of the lineage's key specs.",
		},
		"domain not validated": {
			err:         cloudcertificates.ErrDomainNotValidatedUploadFailed,
			wantSummary: "Domain Not Validated",
			wantDetail:  "One or more SANs are not yet Domain Validated.",
		},
		"unmapped error falls back to the generic diagnostic": {
			err:         errors.New("simulated unexpected failure"),
			wantSummary: "Failed to upload signed certificate",
			wantDetail:  "simulated unexpected failure",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := uploadCertificateErrorDiagnostics(tc.err, 500001, 2912)

			require.Len(t, diags, 1)
			assert.Equal(t, tc.wantSummary, diags[0].Summary())
			assert.Equal(t, tc.wantDetail, diags[0].Detail())
		})
	}
}

func TestCompleteLineageErrorDiagnostics(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err         error
		wantSummary string
		wantDetail  string
	}{
		"lineage not found": {
			err:         cloudcertificates.ErrLineageNotFound,
			wantSummary: "Lineage Not Found",
			wantDetail:  "No certificate lineage found with ID 500001.",
		},
		"no current production": {
			err:         cloudcertificates.ErrNoCurrentProduction,
			wantSummary: "No Current Production",
			wantDetail:  "Lineage 500001 has no current production generation to complete.",
		},
		"complete precondition failed": {
			err:         cloudcertificates.ErrCompletePreconditionFailed,
			wantSummary: "Complete Precondition Failed",
			wantDetail:  "A pending head generation already exists, or no unique CSR_READY algorithm matching the given key type was found on production.",
		},
		"certificate parse error": {
			err:         cloudcertificates.ErrCertParseError,
			wantSummary: "Certificate Parse Error",
			wantDetail:  "The uploaded certificate is malformed or could not be parsed.",
		},
		"certificate/CSR mismatch": {
			err:         cloudcertificates.ErrCertCSRMismatch,
			wantSummary: "Certificate/CSR Mismatch",
			wantDetail:  "The uploaded certificate's public key does not match the stored CSR.",
		},
		"unmapped error falls back to the generic diagnostic": {
			err:         errors.New("simulated unexpected failure"),
			wantSummary: "Failed to complete lineage",
			wantDetail:  "simulated unexpected failure",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := completeLineageErrorDiagnostics(tc.err, 500001)

			require.Len(t, diags, 1)
			assert.Equal(t, tc.wantSummary, diags[0].Summary())
			assert.Equal(t, tc.wantDetail, diags[0].Detail())
		})
	}
}

// signedAlgo builds an algorithmModel with an accepted signed certificate, used as diffAlgorithms test fixtures.
func signedAlgo(pem string) algorithmModel {
	return algorithmModel{SignedCertificatePEM: types.StringValue(pem), TrustChainPEM: types.StringValue("chain")}
}

// unsignedAlgo builds an algorithmModel still awaiting a signed certificate (CSR_READY), used as diffAlgorithms
// test fixtures - a null SignedCertificatePEM must be treated the same as a wholly missing key.
func unsignedAlgo() algorithmModel {
	return algorithmModel{SignedCertificatePEM: types.StringNull(), TrustChainPEM: types.StringNull()}
}

func TestDiffAlgorithms(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		plan, existing                     map[string]algorithmModel
		wantNew, wantModified, wantRemoved map[string]algorithmModel
	}{
		"new: key absent from existing entirely": {
			plan:     map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
			existing: map[string]algorithmModel{},
			wantNew:  map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
		},
		"new: key present in existing but still unsigned (CSR_READY)": {
			plan:     map[string]algorithmModel{"ECDSA": signedAlgo("cert-b")},
			existing: map[string]algorithmModel{"ECDSA": unsignedAlgo()},
			wantNew:  map[string]algorithmModel{"ECDSA": signedAlgo("cert-b")},
		},
		"unchanged: identical signed content in both": {
			plan:     map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
			existing: map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
		},
		"modified: signed certificate content differs": {
			plan:         map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
			existing:     map[string]algorithmModel{"RSA": signedAlgo("cert-b")},
			wantModified: map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
		},
		"modified: trust chain differs even though certificate matches": {
			plan:         map[string]algorithmModel{"RSA": {SignedCertificatePEM: types.StringValue("cert-a"), TrustChainPEM: types.StringValue("chain-a")}},
			existing:     map[string]algorithmModel{"RSA": {SignedCertificatePEM: types.StringValue("cert-a"), TrustChainPEM: types.StringValue("chain-b")}},
			wantModified: map[string]algorithmModel{"RSA": {SignedCertificatePEM: types.StringValue("cert-a"), TrustChainPEM: types.StringValue("chain-a")}},
		},
		"removed: signed key absent from plan": {
			plan:        map[string]algorithmModel{},
			existing:    map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
			wantRemoved: map[string]algorithmModel{"RSA": signedAlgo("cert-a")},
		},
		"not removed: unsigned key absent from plan is simply ignored": {
			plan:     map[string]algorithmModel{},
			existing: map[string]algorithmModel{"RSA": unsignedAlgo()},
		},
		"mixed: new, unchanged, modified, and removed all at once": {
			plan: map[string]algorithmModel{
				"RSA":   signedAlgo("cert-a"),   // unchanged
				"ECDSA": signedAlgo("cert-new"), // modified relative to existing's "cert-old"
				"DSA":   signedAlgo("cert-c"),   // new: not in existing at all
			},
			existing: map[string]algorithmModel{
				"RSA":   signedAlgo("cert-a"),
				"ECDSA": signedAlgo("cert-old"),
				"PQC":   signedAlgo("cert-legacy"), // removed: signed, but absent from plan
				"XMSS":  unsignedAlgo(),            // ignored: unsigned, absent from plan
			},
			wantNew:      map[string]algorithmModel{"DSA": signedAlgo("cert-c")},
			wantModified: map[string]algorithmModel{"ECDSA": signedAlgo("cert-new")},
			wantRemoved:  map[string]algorithmModel{"PQC": signedAlgo("cert-legacy")},
		},
		"nil plan and nil existing": {
			plan:     nil,
			existing: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diff := diffAlgorithms(tc.plan, tc.existing)

			assertMapEqual := func(t *testing.T, expected, actual map[string]algorithmModel) {
				t.Helper()
				if len(expected) == 0 {
					assert.Empty(t, actual)
				} else {
					assert.Equal(t, expected, actual)
				}
			}

			assertMapEqual(t, tc.wantNew, diff.New)
			assertMapEqual(t, tc.wantModified, diff.Modified)
			assertMapEqual(t, tc.wantRemoved, diff.Removed)
		})
	}
}

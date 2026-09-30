package cloudcertificates

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLineageResourceImport(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - import by lineage_id": {
			init: func(m *cloudcertificates.Mock) {
				// once for ImportState's Read, once for the follow-up PlanOnly step's refresh
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// cleanup after ImportStatePersist leaves the resource in state
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_lineage.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:        true,
					ImportStateId:      "500001",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500001").
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("contract_id", "C-0N7RAC7").
						CheckEqual("group_id", "12345").
						CheckEqual("geo_class", "STANDARD_WORLDWIDE").
						CheckEqual("lineage_name", "test non-validated domain").
						CheckEqual("lineage_type", "MULTIPLE_GENERATION").
						CheckEqual("secure_network", "ENHANCED_TLS").
						CheckEqual("stack_mode", "MULTIPLE_STACK").
						CheckEqual("key_specs.%", "2").
						CheckEqual("key_specs.RSA", "2048").
						CheckEqual("key_specs.ECDSA", "P-256").
						CheckEqual("sans.#", "2").
						CheckSetContains("sans", "www.example.com").
						CheckSetContains("sans", "example.com").
						CheckEqual("subject.common_name", "example.com").
						CheckEqual("subject.organization", "Example Corp.").
						CheckEqual("subject.organizational_unit", "IT").
						CheckEqual("subject.country", "US").
						CheckEqual("subject.state", "Massachusetts").
						CheckEqual("subject.locality", "Cambridge").
						CheckEqual("lineage_created_by", "terraform-dev").
						CheckEqual("lineage_created_time", "2026-07-07T10:18:46Z").
						CheckEqual("lineage_modified_by", "terraform-dev").
						CheckEqual("lineage_modified_time", "2026-07-07T10:18:46Z").
						CheckEqual("head.generation_id", "2912").
						CheckEqual("head.generation_status", "CSR_READY").
						CheckEqual("head.algorithms.%", "2").
						CheckEqual("head.algorithms.RSA.algorithm_instance_id", "6122").
						CheckMissing("head.algorithms.RSA.algorithm_instance_modified_by").
						CheckMissing("head.algorithms.RSA.signed_certificate_pem").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import by lineage_id, minimal fields": {
			init: func(m *cloudcertificates.Mock) {
				// once for ImportState's Read, once for the follow-up PlanOnly step's refresh
				mockGetLineage(m, 500006, createLineageResponseMinimal(), 2)
				// cleanup after ImportStatePersist leaves the resource in state
				mockDeleteLineage(m, 500006)
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_lineage.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_minimal.tf"),
					ImportState:        true,
					ImportStateId:      "500006",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500006").
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("stack_mode", "SINGLE_STACK").
						CheckEqual("lineage_name", "www.example.com20260720103154092771").
						CheckEqual("key_specs.%", "1").
						CheckEqual("key_specs.RSA", "2048").
						CheckEqual("sans.#", "2").
						CheckEqual("head.generation_id", "1323").
						CheckEqual("head.algorithms.%", "1").
						CheckMissing("subject.common_name").
						CheckMissing("subject.organization").
						CheckMissing("subject.organizational_unit").
						CheckMissing("subject.country").
						CheckMissing("subject.state").
						CheckMissing("subject.locality").
						// RES-10: the mapper must collapse an all-empty API subject to a null object, not one with
						// every field individually null.
						CheckMissing("subject.%").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_minimal.tf"),
					PlanOnly: true,
				},
			},
		},
		"error - lineage_id is zero": {
			// ImportState rejects the ID before ever calling the API, so nothing needs to be mocked.
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "0",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: 0`),
				},
			},
		},
		"error - lineage_id is negative": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "-5",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: -5`),
				},
			},
		},
		"error - malformed multipart import ID": {
			// ImportIDSplitter rejects the extra part before ever calling the API, so nothing needs to be mocked.
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001,extra",
					ExpectError:   tst.ErrPattern(`invalid number of importID parts: 2; you need to provide an importID in the format 'lineageID'`),
				},
			},
		},
		"error - post-import refresh finds lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				mockGetLineageFails(m, 500001, cloudcertificates.ErrLineageNotFound, 1)
			},
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001",
					ExpectError:   tst.ErrPattern(`Cannot import non-existent remote object`),
				},
			},
		},
		"error - post-import refresh fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockGetLineageFails(m, 500001, cloudcertificates.ErrInternalError, 1)
			},
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001",
					ExpectError:   tst.ErrPattern(`Failed to read lineage`),
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

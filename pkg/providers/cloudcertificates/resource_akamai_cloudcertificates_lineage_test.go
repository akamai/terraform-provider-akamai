package cloudcertificates

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLineageResource(t *testing.T) {
	t.Parallel()

	dsName := "akamai_cloudcertificates_lineage.test"

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - create with two key specs, subject and DOM validation warning": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 1)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: checkLineageFullAttrs(test.NewStateChecker(dsName).
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
						CheckEqual("sans.#", "2").
						CheckEqual("lineage_created_by", "terraform-dev").
						CheckEqual("lineage_created_time", "2026-07-07T10:18:46Z").
						CheckEqual("lineage_modified_by", "terraform-dev").
						CheckEqual("lineage_modified_time", "2026-07-07T10:18:46Z").
						CheckEqual("head.generation_id", "2912").
						CheckEqual("head.generation_status", "CSR_READY").
						CheckEqual("head.algorithms.%", "2").
						CheckEqual("signing_target.generation_id", "2912").
						CheckEqual("signing_target.generation_status", "CSR_READY").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging")).
						Build(),
				},
			},
		},
		"happy path - create with single key spec (SINGLE_STACK), no subject": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestMinimal(), createLineageResponseMinimal())
				mockGetLineage(m, 500006, createLineageResponseMinimal(), 1)
				mockDeleteLineage(m, 500006)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_minimal.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500006").
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("stack_mode", "SINGLE_STACK").
						CheckEqual("lineage_name", "www.example.com20260720103154092771").
						CheckEqual("key_specs.%", "1").
						CheckEqual("head.generation_id", "1323").
						CheckEqual("head.generation_status", "CSR_READY").
						CheckEqual("head.algorithms.%", "1").
						CheckEqual("head.algorithms.RSA.algorithm_instance_id", "2244").
						CheckMissing("head.algorithms.RSA.algorithm_instance_modified_by").
						CheckMissing("head.algorithms.RSA.algorithm_instance_modified_time").
						CheckMissing("head.algorithms.RSA.signed_certificate_issuer").
						CheckMissing("head.algorithms.RSA.signed_certificate_not_valid_after_date").
						CheckMissing("head.algorithms.RSA.signed_certificate_not_valid_before_date").
						CheckMissing("head.algorithms.RSA.signed_certificate_pem").
						CheckMissing("head.algorithms.RSA.signed_certificate_serial_number").
						CheckMissing("head.algorithms.RSA.signed_certificate_sha256_fingerprint").
						CheckMissing("head.algorithms.RSA.trust_chain_pem").
						CheckMissing("subject.common_name").
						CheckMissing("subject.organization").
						CheckMissing("subject.organizational_unit").
						CheckMissing("subject.country").
						CheckMissing("subject.state").
						CheckMissing("subject.locality").
						// RES-10: the mapper must collapse an all-empty API subject to a null object (like
						// current_production/etc. below), not one with every field individually null.
						CheckMissing("subject.%").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging").
						Build(),
				},
			},
		},
		"happy path - lineage fully promoted, no head generation": {
			init: func(m *cloudcertificates.Mock) {
				// mirrors a real captured response: once fully promoted to production and staging with no
				// renewal in progress, the API omits head entirely rather than returning it as null.
				resp := createLineageResponseFull()
				algorithms := resp.Head.Algorithms
				resp.Head = nil
				resp.CurrentProduction = &cloudcertificates.ProductionGeneration{
					ProductionGenerationID:     2912,
					ProductionGenerationStatus: "ACTIVE",
					Generation: cloudcertificates.Generation{
						Algorithms:            algorithms,
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
				}
				resp.CurrentStaging = &cloudcertificates.StagingGeneration{
					StagingGenerationID:     2912,
					StagingGenerationStatus: "ACTIVE",
					Generation: cloudcertificates.Generation{
						Algorithms:            algorithms,
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
				}
				mockCreateLineage(m, createLineageRequestFull(), resp)
				mockGetLineage(m, 500001, resp, 1)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500001").
						CheckMissing("head").
						CheckEqual("current_production.generation_id", "2912").
						CheckEqual("current_production.generation_status", "ACTIVE").
						CheckEqual("current_production.algorithms.%", "2").
						CheckEqual("current_staging.generation_id", "2912").
						CheckEqual("current_staging.generation_status", "ACTIVE").
						CheckEqual("current_staging.algorithms.%", "2").
						CheckEqual("signing_target.generation_id", "2912").
						CheckEqual("signing_target.generation_status", "ACTIVE").
						CheckMissing("previous_production").
						Build(),
				},
			},
		},
		"happy path - lineage_name unknown at plan time (supplied by another resource), skips name conflict check": {
			// terraform_data's output is unknown until it's actually created, simulating a lineage_name that
			// comes from another resource; ModifyPlan must not call ListLineages when it can't yet know what
			// name is actually being requested. A plan-only step keeps that value unknown for the whole
			// operation (a real apply would resolve it before test's own create-time ModifyPlan runs).
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/unknown_lineage_name.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - create then rename": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				// post-apply refresh after create, and pre-plan refresh before the rename update
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// plan-time name conflict check finds no existing lineage with this name
				mockListLineagesEmpty(m, "renamed-lineage", 2)
				mockRenameLineage(m, 500001, "renamed-lineage", renameLineageResponseFull())
				// Update's internal re-fetch, plus one more refresh covering both the post-apply check and the
				// pre-destroy plan
				mockGetLineage(m, 500001, renamedLineageResponseFull(), 2)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "renamed-lineage").
						CheckEqual("lineage_modified_time", "2026-07-08T13:24:17Z").
						CheckEqual("head.algorithms.%", "2").
						Build(),
				},
			},
		},
		"happy path - rename proceeds to apply when ListLineages conflict check fails": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// the plan-time conflict check itself fails (e.g. a transient API error); ModifyPlan only warns
				// and lets the plan proceed, relying on RenameLineage's own conflict detection.
				mockListLineagesFails(m, "renamed-lineage", cloudcertificates.ErrInternalError, 2)
				mockRenameLineage(m, 500001, "renamed-lineage", renameLineageResponseFull())
				mockGetLineage(m, 500001, renamedLineageResponseFull(), 2)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					// Plan succeeds without ExpectError (the ListLineages failure is only a warning), and the
					// rename still completes successfully at apply time.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "renamed-lineage").
						Build(),
				},
			},
		},
		"happy path - create then rename, generation refresh fails after rename (warns, keeps rename)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				mockListLineagesEmpty(m, "renamed-lineage", 2)
				mockRenameLineage(m, 500001, "renamed-lineage", renameLineageResponseFull())
				// Update's internal re-fetch fails: Update falls back to the RenameLineage response merged with
				// the prior state's generation pointers, and surfaces a warning instead of an error.
				mockGetLineageFails(m, 500001, cloudcertificates.ErrInternalError, 1)
				// pre-destroy plan: the next refresh succeeds and picks up full detail.
				mockGetLineage(m, 500001, renamedLineageResponseFull(), 1)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "renamed-lineage").
						CheckEqual("lineage_modified_time", "2026-07-08T13:24:17Z").
						CheckEqual("head.algorithms.%", "2").
						// the fallback maps RenameLineage's own response directly (not a fresh GetLineage), so also
						// assert every other field it returns to prove the fallback doesn't silently clear them.
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("contract_id", "C-0N7RAC7").
						CheckEqual("group_id", "12345").
						CheckEqual("geo_class", "STANDARD_WORLDWIDE").
						CheckEqual("lineage_type", "MULTIPLE_GENERATION").
						CheckEqual("secure_network", "ENHANCED_TLS").
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
						Build(),
				},
			},
		},
		// geo_class has no analogous "change" case: STANDARD_WORLDWIDE is currently its only valid value (see the
		// TODO on its schema definition), so there is no other value to plan a change to.
		"happy path - change contract_id (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// contract_id is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.ContractID = "C-9999999"
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.ContractID = "C-9999999"
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_contract_id.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("contract_id", "C-9999999").
						Build(),
				},
			},
		},
		"happy path - change group_id (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// group_id is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.GroupID = 67890
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.GroupID = 67890
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_group_id.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("group_id", "67890").
						Build(),
				},
			},
		},
		"happy path - change lineage_type (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// lineage_type is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.LineageType = cloudcertificates.LineageTypeSingleGeneration
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.LineageType = "SINGLE_GENERATION"
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_lineage_type.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("lineage_type", "SINGLE_GENERATION").
						Build(),
				},
			},
		},
		"happy path - change secure_network (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// secure_network is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.SecureNetwork = cloudcertificates.SecureNetworkStandardTLS
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.SecureNetwork = "STANDARD_TLS"
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_secure_network.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("secure_network", "STANDARD_TLS").
						Build(),
				},
			},
		},
		"happy path - change key_specs (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// key_specs is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.KeySpecs = []cloudcertificates.KeySpec{
					{KeyType: "ECDSA", KeySize: "P-256"},
					{KeyType: "RSA", KeySize: "4096"},
				}
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.KeySpecs = []cloudcertificates.KeySpecResponse{
					{KeySize: "4096", KeyType: "RSA"},
					{KeySize: "P-256", KeyType: "ECDSA"},
				}
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_key_specs.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("key_specs.RSA", "4096").
						Build(),
				},
			},
		},
		"happy path - change sans (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// sans is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.SANs = []string{"new.example.com", "www.example.com"}
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.SANs = []string{"www.example.com", "new.example.com"}
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_sans.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("sans.#", "2").
						CheckSetContains("sans", "www.example.com").
						CheckSetContains("sans", "new.example.com").
						Build(),
				},
			},
		},
		"happy path - change subject (requires replace)": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// subject is RequiresReplace: destroy the old lineage, then create a new one
				mockDeleteLineage(m, 500001)
				req := createLineageRequestFull()
				req.Body.Subject.Organization = "New Corp."
				resp := createLineageResponseFull()
				resp.LineageID = 500002
				resp.Subject.Organization = "New Corp."
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500002, resp, 1)
				mockDeleteLineage(m, 500002)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/update_subject.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500002").
						CheckEqual("subject.organization", "New Corp.").
						Build(),
				},
			},
		},
		"expect error - rename fails, name conflict": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// plan-time check finds no conflict, but the API itself rejects it at apply time (simulating a
				// race: another lineage took the name between plan and apply)
				mockListLineagesEmpty(m, "renamed-lineage", 2)
				mockRenameLineageFails(m, 500001, "renamed-lineage", cloudcertificates.ErrLineageNameConflict)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					ExpectError: tst.ErrPattern(`A lineage named "renamed-lineage" already exists for this account\.`),
				},
			},
		},
		"expect error - rename blocked at plan time, name conflict detected early": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// plan-time check finds an existing lineage with this name; RenameLineage is never called
				mockListLineagesConflict(m, "renamed-lineage", 999999)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					ExpectError: tst.ErrPattern(`A lineage named "renamed-lineage" already exists for this account\.`),
				},
			},
		},
		"happy path - refresh finds lineage deleted remotely, removes from state": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 1)
				// refresh during the plan-only step: lineage was deleted outside terraform
				mockGetLineageFails(m, 500001, cloudcertificates.ErrLineageNotFound, 1)
				// test cleanup: a plan-only step never persists the removal, so the resource is still present
				// in state when the test tears down and must still be destroyed for real.
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					// Read succeeds without ExpectError (a warning, not an error, was surfaced), and the plan
					// is non-empty because the resource - now absent from the refreshed state - is planned
					// to be recreated.
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"expect error - create fails, account not allowed": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineageFails(m, createLineageRequestFull(), cloudcertificates.ErrLineageAccountNotAllowed)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ExpectError: tst.ErrPattern("This account is not permitted to use Certificate Lineage."),
				},
			},
		},
		"expect error - delete fails, lineage has active production": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				mockDeleteLineageFails(m, 500001, cloudcertificates.ErrLineageHasActiveProduction)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty.tf"),
					ExpectError: tst.ErrPattern(`Lineage 500001 has an active production certificate\. Deactivate before deleting\.`),
				},
			},
		},
		"expect error - delete fails, lineage has active staging": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				mockDeleteLineageFails(m, 500001, cloudcertificates.ErrLineageHasActiveStaging)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty.tf"),
					ExpectError: tst.ErrPattern(`Lineage 500001 has an active staging certificate\. Deactivate before deleting\.`),
				},
			},
		},
		"expect error - delete fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// unmapped error falls through to the generic message
				mockDeleteLineageFails(m, 500001, cloudcertificates.ErrInternalError)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty.tf"),
					ExpectError: tst.ErrPattern(`Failed to delete lineage`),
				},
			},
		},
		"expect error - rename fails, lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// plan-time check finds no conflict, but the rename itself fails: lineage was deleted outside
				// terraform between the refresh and the rename call
				mockListLineagesEmpty(m, "renamed-lineage", 2)
				mockRenameLineageFails(m, 500001, "renamed-lineage", cloudcertificates.ErrLineageNotFound)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					ExpectError: tst.ErrPattern(`No certificate lineage found with ID 500001`),
				},
			},
		},
		"expect error - rename fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// plan-time check finds no conflict, but the rename itself fails with an unmapped error
				mockListLineagesEmpty(m, "renamed-lineage", 2)
				mockRenameLineageFails(m, 500001, "renamed-lineage", cloudcertificates.ErrInternalError)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "test non-validated domain").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/rename.tf"),
					ExpectError: tst.ErrPattern(`Failed to rename lineage`),
				},
			},
		},
		"expect error - create fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineageFails(m, createLineageRequestFull(), cloudcertificates.ErrInternalError)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ExpectError: tst.ErrPattern(`Failed to create lineage`),
				},
			},
		},
		"expect error - refresh fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
				mockGetLineage(m, 500001, createLineageResponseFull(), 1)
				// refresh during the plan-only step fails with an unmapped error
				mockGetLineageFails(m, 500001, cloudcertificates.ErrInternalError, 1)
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("lineage_id", "500001").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					PlanOnly:    true,
					ExpectError: tst.ErrPattern(`Failed to read lineage`),
				},
			},
		},
		"happy path - create with partially populated subject": {
			init: func(m *cloudcertificates.Mock) {
				req := createLineageRequestMinimal()
				req.Body.Subject = ptr.To(cloudcertificates.Subject{
					CommonName:   "example.com",
					Organization: "Example Corp.",
				})
				resp := createLineageResponseMinimal()
				resp.Subject = cloudcertificates.Subject{
					CommonName:   "example.com",
					Organization: "Example Corp.",
				}
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500006, resp, 1)
				mockDeleteLineage(m, 500006)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_partial_subject.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("subject.common_name", "example.com").
						CheckEqual("subject.organization", "Example Corp.").
						CheckMissing("subject.organizational_unit").
						CheckMissing("subject.country").
						CheckMissing("subject.state").
						CheckMissing("subject.locality").
						Build(),
				},
			},
		},
		"happy path - create with explicit lineage_name and lineage_type": {
			init: func(m *cloudcertificates.Mock) {
				// ModifyPlan's name conflict check runs on create too, once lineage_name is explicitly known
				mockListLineagesEmpty(m, "explicit-lineage", 2)
				req := createLineageRequestMinimal()
				req.Body.LineageName = "explicit-lineage"
				req.Body.LineageType = cloudcertificates.LineageTypeSingleGeneration
				resp := createLineageResponseMinimal()
				resp.LineageName = "explicit-lineage"
				resp.LineageType = "SINGLE_GENERATION"
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, 500006, resp, 1)
				mockDeleteLineage(m, 500006)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_explicit_name_type.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_name", "explicit-lineage").
						CheckEqual("lineage_type", "SINGLE_GENERATION").
						Build(),
				},
			},
		},
		"validation error - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_contract_id.tf"),
					ExpectError: tst.ErrPattern(`The argument "contract_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing key_specs": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_key_specs.tf"),
					ExpectError: tst.ErrPattern(`The argument "key_specs" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing sans": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_sans.tf"),
					ExpectError: tst.ErrPattern(`The argument "sans" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing secure_network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_secure_network.tf"),
					ExpectError: tst.ErrPattern(`The argument "secure_network" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_group_id.tf"),
					ExpectError: tst.ErrPattern(`The argument "group_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - too many key_specs": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/too_many_key_specs.tf"),
					ExpectError: tst.ErrPattern(`Attribute key_specs map must contain at least 1 elements and at most 2 elements, got: 3`),
				},
			},
		},
		"validation error - empty key_specs": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty_key_specs.tf"),
					ExpectError: tst.ErrPattern(`Attribute key_specs map must contain at least 1 elements and at most 2 elements, got: 0`),
				},
			},
		},
		"validation error - empty sans": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans set must contain at least 1 elements and at most 100 elements, got: 0`),
				},
			},
		},
		"validation error - too many sans": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/too_many_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans set must contain at least 1 elements and at most 100 elements, got: 101`),
				},
			},
		},
		"happy path - duplicate sans collapse into a single set element": {
			init: func(m *cloudcertificates.Mock) {
				req := createLineageRequestMinimal()
				req.Body.SANs = []string{"www.example.com"}
				resp := createLineageResponseMinimal()
				resp.SANs = []string{"www.example.com"}
				mockCreateLineage(m, req, resp)
				mockGetLineage(m, resp.LineageID, resp, 1)
				mockDeleteLineage(m, resp.LineageID)
			},
			steps: []resource.TestStep{
				{
					// sans is a set: Terraform deduplicates identical config values before our validators ever
					// run, so two identical sans values silently collapse into a single element instead of erroring.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/duplicate_sans.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("sans.#", "1").
						CheckSetContains("sans", "www.example.com").
						Build(),
				},
			},
		},
		"validation error - blank sans element": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/blank_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] must be a valid lowercase domain name or wildcard domain, e.g. example.com or \*.example.com, got:`),
				},
			},
		},
		"validation error - whitespace-only sans element": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] must be a valid lowercase domain name or wildcard domain, e.g. example.com or \*.example.com, got:`),
				},
			},
		},
		"validation error - malformed sans element": {
			// mirrors a real 400 schema-validation-failure response captured from the API for this exact SAN shape.
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/malformed_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] must be a valid lowercase domain name or wildcard domain, e.g. example.com or \*.example.com, got: example com`),
				},
			},
		},
		"validation error - malformed wildcard sans element": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/malformed_wildcard_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] must be a valid lowercase domain name or wildcard domain, e.g. example.com or \*.example.com, got: \*example.com`),
				},
			},
		},
		"validation error - uppercase sans element": {
			// the API lowercases every sans value it stores, so uppercase input is rejected rather than
			// silently normalized: accepting it would leave the config permanently out of sync with state.
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/uppercase_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] must be a valid lowercase domain name or wildcard domain, e.g. example.com or \*.example.com, got: WWW.example.com`),
				},
			},
		},
		"validation error - sans element too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/long_sans.tf"),
					ExpectError: tst.ErrPattern(`Attribute sans\[Value\(.*?\)\] string length must be at most 253, got: 257`),
				},
			},
		},
		"validation error - empty contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty_contract_id.tf"),
					ExpectError: tst.ErrPattern(`Attribute contract_id string length must be at least 1, got: 0`),
				},
			},
		},
		"validation error - whitespace-only contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_contract_id.tf"),
					ExpectError: tst.ErrPattern(`Attribute contract_id must not be blank, got:`),
				},
			},
		},
		"validation error - group_id less than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_group_id.tf"),
					ExpectError: tst.ErrPattern(`Attribute group_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - invalid geo_class": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_geo_class.tf"),
					ExpectError: tst.ErrPattern(`Attribute geo_class value must be one of: \["STANDARD_WORLDWIDE"\], got: "INVALID_GEO_CLASS"`),
				},
			},
		},
		"validation error - lineage_name too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/lineage_name_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_name string length must be between 1 and 270, got: 271`),
				},
			},
		},
		"validation error - lineage_name leading whitespace": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/lineage_name_leading_whitespace.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_name must not have leading or trailing whitespace, got:`),
				},
			},
		},
		"validation error - lineage_name trailing whitespace": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/lineage_name_trailing_whitespace.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_name must not have leading or trailing whitespace, got:`),
				},
			},
		},
		"validation error - lineage_name blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/lineage_name_blank.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_name must not have leading or trailing whitespace, got:`),
				},
			},
		},
		"validation error - lineage_name invalid character": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_lineage_name_charset.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_name must contain only letters, digits, spaces, underscores, periods, and hyphens, got: test@lineage!`),
				},
			},
		},
		"validation error - invalid lineage_type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_lineage_type.tf"),
					ExpectError: tst.ErrPattern(`Attribute lineage_type value must be one of: \["MULTIPLE_GENERATION" "SINGLE_GENERATION"\], got: "INVALID_LINEAGE_TYPE"`),
				},
			},
		},
		"validation error - invalid secure_network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_secure_network.tf"),
					ExpectError: tst.ErrPattern(`Attribute secure_network value must be one of: \["STANDARD_TLS" "ENHANCED_TLS"\], got: "INVALID_SECURE_NETWORK"`),
				},
			},
		},
		"validation error - invalid key_specs key_type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_key_type.tf"),
					ExpectError: tst.ErrPattern(`key_type "INVALID_KEY_TYPE" is not valid; must be one of: RSA, ECDSA\.`),
				},
			},
		},
		"validation error - invalid key_specs key_size": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_key_size.tf"),
					ExpectError: tst.ErrPattern(`key_size "INVALID_KEY_SIZE" is not valid for key_type "RSA"; must be one of: 2048, 4096\.`),
				},
			},
		},
		"validation error - missing key_specs key_size": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/missing_key_specs_key_size.tf"),
					ExpectError: tst.ErrPattern(`This attribute contains a null value`),
				},
			},
		},
		"validation error - key_specs key_size mismatched with key_type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/mismatched_key_specs_key_size.tf"),
					ExpectError: tst.ErrPattern(`key_size "P-256" is not valid for key_type "RSA"; must be one of: 2048, 4096\.`),
				},
			},
		},
		"validation error - subject empty object": {
			// TFP-07/COD-08: subject = {} is otherwise accepted since every field is optional, but it plans as a
			// known, non-null object while an empty API response collapses to null, breaking Terraform's
			// post-apply consistency check. Reject it outright instead.
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/empty_subject_object.tf"),
					ExpectError: tst.ErrPattern(`subject must contain at least one field, or be omitted entirely; an empty object is not allowed\.`),
				},
			},
		},
		"validation error - subject common_name blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/blank_subject_common_name.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.common_name must not have leading whitespace, got:`),
				},
			},
		},
		"validation error - subject organizational_unit blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/blank_subject_organizational_unit.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.organizational_unit must not have leading whitespace, got:`),
				},
			},
		},
		"validation error - subject organization blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_subject_organization.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.organization must not have leading whitespace, got:`),
				},
			},
		},
		"validation error - subject country blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_subject_country.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.country must contain only two letters, got:`),
				},
			},
		},
		"validation error - subject country not letters": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_country.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.country must contain only two letters, got: 12`),
				},
			},
		},
		"validation error - subject state blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_subject_state.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.state must not have leading whitespace, got:`),
				},
			},
		},
		"validation error - subject locality blank": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/whitespace_subject_locality.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.locality must not have leading whitespace, got:`),
				},
			},
		},
		"validation error - subject common_name too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_common_name_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.common_name string length must be between 1 and 64, got: 65`),
				},
			},
		},
		"validation error - subject organization too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_organization_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.organization string length must be between 1 and 64, got: 65`),
				},
			},
		},
		"validation error - subject organizational_unit too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_organizational_unit_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.organizational_unit string length must be between 1 and 64, got: 65`),
				},
			},
		},
		"validation error - subject state too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_state_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.state string length must be between 1 and 128, got: 129`),
				},
			},
		},
		"validation error - subject locality too long": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_locality_too_long.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.locality string length must be between 1 and 128, got: 129`),
				},
			},
		},
		"validation error - subject country too short": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_country_too_short.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.country string length must be between 2 and 2, got: 1`),
				},
			},
		},
		"validation error - invalid subject country length": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/invalid_subject_country_length.tf"),
					ExpectError: tst.ErrPattern(`Attribute subject\.country string length must be between 2 and 2, got: 3`),
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

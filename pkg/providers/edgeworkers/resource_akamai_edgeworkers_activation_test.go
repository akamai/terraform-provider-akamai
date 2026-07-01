package edgeworkers

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceEdgeworkersActivation(t *testing.T) {
	t.Parallel()
	workdir := "./testdata/TestResourceEdgeWorkersActivation"
	edgeworkerID := 1234
	baseChecker := test.NewStateChecker("akamai_edgeworkers_activation.test").
		CheckEqual("activation_id", "1").
		CheckEqual("version", "test").
		CheckEqual("network", stagingNetwork).
		CheckEqual("note", "note for edgeworkers activation").
		CheckEqual("auto_pin", "true")

	baseImportChecker := test.NewImportChecker().
		CheckEqual("edgeworker_id", "1234").
		CheckEqual("network", "STAGING").
		CheckEqual("activation_id", "1").
		CheckEqual("auto_pin", "true")

	tests := map[string]struct {
		init            func(*edgeworkers.Mock)
		steps           []resource.TestStep
		omitDefaultMock bool
	}{
		"create and read activation - no previous activations": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("timeouts.#", "0").
						Build(),
				},
			},
		},
		"create and read activation - autoPin false": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"
				autoPin := false

				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create - note: autoPin=false in API request
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, autoPin)

				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_autopin_false_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("auto_pin", "false").
						Build(),
				},
			},
		},
		"create and read activation - with timeout": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_with_timeout.tf", workdir),
					Check: baseChecker.
						CheckEqual("timeouts.#", "1").
						CheckEqual("timeouts.0.default", "2h").
						CheckEqual("timeouts.0.delete", "3h").
						Build(),
				},
			},
		},
		"create and read activation - some previous activations, but no current": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 8
				note := "note for edgeworkers activation"

				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 7, edgeworkers.ActivationNetworkProduction, "current", activationStatusComplete, "2022-01-25T12:30:06Z", note),
					*createStubActivation(edgeworkerID, 6, net, "past2", activationStatusComplete, "2022-01-25T12:30:06Z", note),
					*createStubActivation(edgeworkerID, 5, net, "past1", activationStatusComplete, "2022-01-24T12:30:06Z", note),
					*createStubActivation(edgeworkerID, 4, edgeworkers.ActivationNetworkProduction, "past1", activationStatusComplete, "2022-01-23T12:30:06Z", note),
					*createStubActivation(edgeworkerID, 3, net, "past2", activationStatusComplete, "2022-01-23T18:30:06Z", note),
					*createStubActivation(edgeworkerID, 2, net, "past2", activationStatusComplete, "2022-01-23T12:30:06Z", note),
					*createStubActivation(edgeworkerID, 1, net, "past1", activationStatusComplete, "2022-01-22T12:30:06Z", note),
				}

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
					*createStubEdgeworkerVersion(edgeworkerID, "past1"),
					*createStubEdgeworkerVersion(edgeworkerID, "past2"),
				}, nil).Once()

				// get current activation
				expectListActivations(m, edgeworkerID, "", activations, nil).Once()
				expectListDeactivations(m, edgeworkerID, "past2", []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 2, net, "past2", activationStatusComplete, "2022-01-24T10:30:06Z"),
					*createStubDeactivation(edgeworkerID, 1, net, "past2", activationStatusComplete, "2022-01-23T15:30:06Z"),
					*createStubDeactivation(edgeworkerID, 3, net, "past2", activationStatusComplete, "2022-01-26T12:30:06Z"),
				}, nil).Once()

				// activate
				expectActivateVersion(m, edgeworkerID, activationID, net, version, note, true, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, append([]edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note)},
					activations...,
				), []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 4, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "8").
						Build(),
				},
			},
		},
		"create and read activation - version already active": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, nil).Once()
				expectListDeactivations(m, edgeworkerID, version, []edgeworkers.Deactivation{}, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"create and read activation - longer polling": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil)

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil).Once()

				// activate
				expectActivateVersion(m, edgeworkerID, activationID, net, version, note, true, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"create and read activation - version is already being activated, wait for activation": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil)

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusPending, "", note),
				}, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"update network": {
			init: func(m *edgeworkers.Mock) {
				createNet, updateNet := edgeworkers.ActivationNetworkStaging, edgeworkers.ActivationNetworkProduction
				version := "test"
				createActivationID, updateActivationID := 1, 2
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, createNet, version, note, true)

				// read + plan + refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, version, activations, []edgeworkers.Deactivation{}, 3)

				// update - activate
				expectFullUpdate(m, edgeworkerID, updateActivationID, updateNet, version, "", note, activations, true)

				// read + plan
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, updateNet, version, activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, updateNet, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_prod.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "2").
						CheckEqual("network", productionNetwork).
						Build(),
				},
			},
		},
		"update version": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 2
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan + refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, createVersion, activations, []edgeworkers.Deactivation{}, 3)

				// update - activate
				expectFullUpdate(m, edgeworkerID, updateActivationID, net, updateVersion, createVersion, note, activations, true)

				// read + plan
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "2").
						CheckEqual("version", "test1").
						Build(),
				},
			},
		},
		"update version - active version changed on refresh": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 3
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion"),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 2, net, "someOtherVersion", activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, "someOtherVersion", activations, []edgeworkers.Deactivation{}, 1)

				// update - activate
				expectFullUpdate(m, edgeworkerID, updateActivationID, net, updateVersion, "someOtherVersion", note, activations, true)

				// read + plan
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, 2, net, "someOtherVersion", activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "3").
						CheckEqual("version", "test1").
						Build(),
				},
			},
		},
		"update version - version already active": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 2
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "2").
						CheckEqual("version", "test1").
						Build(),
				},
			},
		},
		"update network - version already active": {
			init: func(m *edgeworkers.Mock) {
				createNet, updateNet := edgeworkers.ActivationNetworkStaging, edgeworkers.ActivationNetworkProduction
				version := "test"
				createActivationID, updateActivationID := 1, 2
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, createNet, version, note, true)

				// read + plan
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, updateNet, version, activationStatusComplete, "", note),
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, version, activations, []edgeworkers.Deactivation{}, 1)

				// update
				expectListActivations(m, edgeworkerID, "", activations, nil).Once()
				expectListDeactivations(m, edgeworkerID, version, []edgeworkers.Deactivation{}, nil).Once()

				// read + plan
				expectFullRead(m, edgeworkerID, version, activations, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 2, updateNet, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_prod.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "2").
						CheckEqual("network", productionNetwork).
						Build(),
				},
			},
		},
		"update edgeworker_id - ForceNew success": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				updateEdgeworkerID := 4321
				note := "note for edgeworkers activation"

				expectListEdgeWorkersID(m, nil, edgeworkerID, updateEdgeworkerID)

				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read + plan + refresh
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 3)

				// destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)

				// create - version verification
				expectListEdgeWorkerVersions(m, updateEdgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(updateEdgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, updateEdgeworkerID, activationID, net, version, note, true)

				// read + plan
				expectFullRead(m, updateEdgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(updateEdgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, updateEdgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("edgeworker_id", "1234").
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_different_edgeworker_id.tf", workdir),
					Check: baseChecker.
						CheckEqual("edgeworker_id", "4321").
						Build(),
				},
			},
			omitDefaultMock: true,
		},
		"update note - diff suppressed when other fields not changed": {
			init: func(m *edgeworkers.Mock) {
				createNet := edgeworkers.ActivationNetworkStaging
				version := "test1"
				createActivationID := 1
				note := "note for edgeworkers activation"

				// create  - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Times(1)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, createNet, version, note, true)

				// read + plan + refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, version, activations, []edgeworkers.Deactivation{}, 4)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, createNet, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("version", "test1").
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_note_update_no_activation.tf", workdir),
					Check: baseChecker.
						CheckEqual("version", "test1").
						Build(),
				},
			},
		},
		"update note - when other fields changed it does not update but creates new activation ": {
			init: func(m *edgeworkers.Mock) {
				createNet, updateNet := edgeworkers.ActivationNetworkStaging, edgeworkers.ActivationNetworkProduction
				version := "test1"
				createActivationID, updateActivationID := 1, 2
				note, updatedNote := "note for edgeworkers activation", "note for edgeworkers activation updated"

				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, createNet, version, note, true)

				// read + plan + refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, version, activations, []edgeworkers.Deactivation{}, 3)

				// update - activate
				expectFullUpdate(m, edgeworkerID, updateActivationID, updateNet, version, "", updatedNote, activations, true)

				// read + plan
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, updateNet, version, activationStatusComplete, "", updatedNote),
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, updateNet, version, updatedNote)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("version", "test1").
						Build(),
				},
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_note_update.tf", workdir),
					Check: baseChecker.
						CheckEqual("activation_id", "2").
						CheckEqual("version", "test1").
						CheckEqual("network", productionNetwork).
						CheckEqual("note", "note for edgeworkers activation updated").
						Build(),
				},
			},
		},
		"destroy - version already deactivated": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version, note,
					fmt.Errorf("%w: %s", edgeworkers.ErrVersionAlreadyDeactivated, "oops"))
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"destroy - version is being deactivated, wait": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version, note,
					fmt.Errorf("%w: %s", edgeworkers.ErrVersionBeingDeactivated, "oops"))
				expectListDeactivations(m, edgeworkerID, version, []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 1, net, version, activationStatusInProgress, ""),
				}, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusComplete, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"destroy - longer polling": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version, note, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusComplete, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"destroy - timeout, resource deleted successfully": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				// A bit hack to simulate timeout is returning ErrEdgeworkerDeactivationTimeout on GetDeactivation
				expectDeactivateVersion(m, edgeworkerID, 1, net, version, note, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, "", ErrEdgeworkerDeactivationTimeout).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"import activation on staging": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", ""),
				}, []edgeworkers.Deactivation{}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, "")
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ImportStateCheck: baseImportChecker.
						CheckEqual("version", "test").
						CheckEqual("note", "").
						Build(),
					ImportStatePersist: true,
				},
			},
			omitDefaultMock: true,
		},
		"import activation on staging with note": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test1"
				activationID := 1
				note := "note for edgeworkers activation updated"

				// import read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_note_update_no_activation.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ImportStateCheck: baseImportChecker.
						CheckEqual("version", "test1").
						CheckEqual("note", "note for edgeworkers activation updated").
						Build(),
					ImportStatePersist: true,
				},
			},
			omitDefaultMock: true,
		},
		"error on create - missing required arguments": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_missing_required_args.tf", workdir),
					ExpectError: regexp.MustCompile("argument \"version\" is required"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_missing_required_args.tf", workdir),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_missing_required_args.tf", workdir),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
			},
			omitDefaultMock: true,
		},
		"error on create - version does not exist": {
			init: func(m *edgeworkers.Mock) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion2"),
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion1"),
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile(`version 'test' is not valid for edgeworker with id=1234`),
				},
			},
		},
		"error on create - getting current activation failed, ListActivations API error": {
			init: func(m *edgeworkers.Mock) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - getting current activation failed, ListDeactivations API error": {
			init: func(m *edgeworkers.Mock) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 1, edgeworkers.ActivationNetworkStaging, "test", activationStatusComplete, "", ""),
				}, nil)
				expectListDeactivations(m, edgeworkerID, "test", []edgeworkers.Deactivation{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on list version": {
			init: func(m *edgeworkers.Mock) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{}, fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on activate": {
			init: func(m *edgeworkers.Mock) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, 1, edgeworkers.ActivationNetworkStaging, "test", "note for edgeworkers activation", true, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on polling": {
			init: func(m *edgeworkers.Mock) {
				version := "test"
				net := edgeworkers.ActivationNetworkStaging
				activationID := 1
				note := "note for edgeworkers activation"

				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, activationID, net, version, note, true, nil)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, "", fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - activation failed": {
			init: func(m *edgeworkers.Mock) {
				version := "test"
				net := edgeworkers.ActivationNetworkStaging
				activationID := 1
				note := "note for edgeworkers activation"

				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, activationID, net, version, note, true, nil)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, "ERROR", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: edgeworker activation failure"),
				},
			},
		},
		"error on update - version does not exist": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID := 1
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan + refresh
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 3)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					ExpectError: regexp.MustCompile(`version 'test1' is not valid for edgeworker with id=1234`),
				},
			},
		},
		"error on edgeworker_id ForceNew": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// create version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read + plan
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 3)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_different_edgeworker_id.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation: edgeworker with id=4321 was not found"),
				},
			},
		},
		"error on update - no current activation on refresh": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID := 1
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "2022-01-25T12:30:06Z", note),
				}, []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 1, net, createVersion, activationStatusComplete, "2022-01-26T12:30:06Z"),
				}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, createVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					ExpectError: regexp.MustCompile("edgeworker activation read: no version active on network 'STAGING' for edgeworker with id=1234"),
				},
			},
		},
		"error on update - error waiting for deactivation": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID := 1
				note := "note for edgeworkers activation"

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test1"),
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion, note, true)

				// read + plan + refresh
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "", note),
				}
				expectFullRead(m, edgeworkerID, createVersion, activations, []edgeworkers.Deactivation{}, 3)

				// update
				expectListActivations(m, edgeworkerID, "", activations, nil).Once()
				expectListDeactivations(m, edgeworkerID, createVersion, []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 1, net, createVersion, activationStatusInProgress, ""),
				}, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, createVersion, activationStatusInProgress, nil).Once()
				expectGetDeactivation(m, edgeworkerID, 1, net, createVersion, "", fmt.Errorf("oops")).Once()

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test1_stag.tf", workdir),
					ExpectError: regexp.MustCompile(`edgeworker activation: oops`),
				},
			},
		},
		"error on customize diff - error listing edgeworkers": {
			init: func(m *edgeworkers.Mock) {
				// create version verification
				expectListEdgeWorkersID(m, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_version_test_stag.tf", workdir),
					ExpectError: regexp.MustCompile(`edgeworker activation: oops`),
				},
			},
			omitDefaultMock: true,
		},
		"error on customize diff - create + auto_pin changed without other field changes": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// create
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()
				expectFullActivation(m, edgeworkerID, activationID, net, version, note, true)

				// read x 2
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 2)

				// plan with update of auto_pin=false, but no change to other fields, should cause error
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_autopin_true_stag.tf", workdir),
					Check: baseChecker.
						CheckEqual("auto_pin", "true").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_autopin_false_stag.tf", workdir),
					ExpectError: regexp.MustCompile(`edgeworker activation: 'auto_pin' can only be changed together with 'version', 'network', or 'edgeworker_id'`),
				},
			},
		},
		"import activation on staging with auto_pin true": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// import read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", ""),
				}, []edgeworkers.Deactivation{}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, "")
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING:true", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ImportStateCheck: baseImportChecker.
						CheckEqual("version", "test").
						CheckEqual("note", "").
						Build(),
					ImportStatePersist: true,
				},
			},
			omitDefaultMock: true,
		},
		"import activation on staging with auto_pin false": {
			init: func(m *edgeworkers.Mock) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1
				note := "note for edgeworkers activation"

				// import read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "", note),
				}, []edgeworkers.Deactivation{}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version, note)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_autopin_false_stag.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING:false", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ImportStateCheck: baseImportChecker.
						CheckEqual("auto_pin", "false").
						CheckEqual("version", "test").
						CheckEqual("note", "note for edgeworkers activation").
						Build(),
					ImportStatePersist: true,
				},
			},
			omitDefaultMock: true,
		},
		"error on import - invalid auto_pin value": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING:notabool", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ExpectError:   regexp.MustCompile(`edgeworker activation import: auto_pin must be a boolean, got 'notabool'`),
				},
			},
			omitDefaultMock: true,
		},
		"error on import - too many parts": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:STAGING:true:extra", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ExpectError:   regexp.MustCompile(`edgeworker activation import: invalid import id`),
				},
			},
			omitDefaultMock: true,
		},
		"error on import - edgeworker id not a number": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: "123abc:STAGING",
					ResourceName:  "akamai_edgeworkers_activation.test",
					ExpectError:   regexp.MustCompile(`edgeworker activation import: edgeworker id must be an integer, got '123abc'`),
				},
			},
			omitDefaultMock: true,
		},
		"error on import - invalid network": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureStringf(t, "%s/edgeworkers_activation_import.tf", workdir),
					ImportState:   true,
					ImportStateId: fmt.Sprintf("%d:INVALID_NETWORK", edgeworkerID),
					ResourceName:  "akamai_edgeworkers_activation.test",
					ExpectError:   regexp.MustCompile(`edgeworker activation import: network must be 'STAGING' or 'PRODUCTION', got 'INVALID_NETWORK'`),
				},
			},
			omitDefaultMock: true,
		},
	}

	// redefining times to accelerate tests
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if !test.omitDefaultMock {
				expectListEdgeWorkersID(client.EdgeWorkers, nil, edgeworkerID)
			}
			if test.init != nil {
				test.init(client.EdgeWorkers)
			}
			config := defaultSubproviderConfig()
			config.activation.pollMinimum = time.Millisecond * 1
			config.activation.pollInterval = config.activation.pollMinimum
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(config)),
				IsUnitTest:               true,
				Steps:                    test.steps,
			})
			client.EdgeWorkers.AssertExpectations(t)
		})
	}
}

func expectActivateVersion(m *edgeworkers.Mock, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, note string, autoPin bool, e error) *mock.Call {
	req := edgeworkers.ActivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		ActivateVersion: edgeworkers.ActivateVersion{
			Network: net,
			Version: version,
			Note:    note,
			AutoPin: ptr.To(autoPin),
		},
	}
	if e != nil {
		return m.On("ActivateVersion", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("ActivateVersion", testutils.MockContext, req).Return(createStubActivation(edgeworkerID, activationID, net, version, activationStatusPresubmit, "", note), nil)
}

func expectGetActivation(m *edgeworkers.Mock, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetActivationRequest{
		EdgeWorkerID: edgeworkerID,
		ActivationID: activationID,
	}
	if e != nil {
		return m.On("GetActivation", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("GetActivation", testutils.MockContext, req).Return(createStubActivation(edgeworkerID, activationID, net, version, status, "", ""), nil)
}

func expectListActivations(m *edgeworkers.Mock, edgeworkerID int, version string, activations []edgeworkers.Activation, e error) *mock.Call {
	req := edgeworkers.ListActivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	}
	if e != nil {
		return m.On("ListActivations", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("ListActivations", testutils.MockContext, req).Return(&edgeworkers.ListActivationsResponse{
		Activations: activations,
	}, nil)
}

func expectListDeactivations(m *edgeworkers.Mock, edgeworkerID int, version string, deactivations []edgeworkers.Deactivation, e error) *mock.Call {
	req := edgeworkers.ListDeactivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	}
	if e != nil {
		return m.On("ListDeactivations", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("ListDeactivations", testutils.MockContext, req).Return(&edgeworkers.ListDeactivationsResponse{
		Deactivations: deactivations,
	}, nil)
}

func expectDeactivateVersion(m *edgeworkers.Mock, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, note string, e error) *mock.Call {
	req := edgeworkers.DeactivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		DeactivateVersion: edgeworkers.DeactivateVersion{
			Network: net,
			Version: version,
			Note:    note,
		},
	}
	if e != nil {
		return m.On("DeactivateVersion", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("DeactivateVersion", testutils.MockContext, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, activationStatusPresubmit, ""), nil)
}

func expectGetDeactivation(m *edgeworkers.Mock, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetDeactivationRequest{
		EdgeWorkerID:   edgeworkerID,
		DeactivationID: deactivationID,
	}
	if e != nil {
		return m.On("GetDeactivation", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("GetDeactivation", testutils.MockContext, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, status, ""), nil)
}

func expectListEdgeWorkerVersions(m *edgeworkers.Mock, edgeworkerID int, versions []edgeworkers.EdgeWorkerVersion, e error) *mock.Call {
	req := edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeworkerID,
	}
	if e != nil {
		return m.On("ListEdgeWorkerVersions", testutils.MockContext, req).Return(nil, e)
	}

	return m.On("ListEdgeWorkerVersions", testutils.MockContext, req).Return(&edgeworkers.ListEdgeWorkerVersionsResponse{
		EdgeWorkerVersions: versions,
	}, nil)
}

func expectListEdgeWorkersID(m *edgeworkers.Mock, e error, ewIDs ...int) *mock.Call {
	call := m.On("ListEdgeWorkersID", testutils.MockContext, edgeworkers.ListEdgeWorkersIDRequest{})
	if e != nil {
		return call.Return(nil, e)
	}
	ews := make([]edgeworkers.EdgeWorkerID, len(ewIDs))
	for i, ewID := range ewIDs {
		ews[i] = edgeworkers.EdgeWorkerID{EdgeWorkerID: ewID}
	}
	return call.Return(&edgeworkers.ListEdgeWorkersIDResponse{
		EdgeWorkers: ews,
	}, nil)
}

func expectFullActivation(m *edgeworkers.Mock, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, note string, autoPin bool) {
	expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil).Once()
	expectActivateVersion(m, edgeworkerID, activationID, net, version, note, autoPin, nil).Once()
	expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()
}

func expectFullUpdate(m *edgeworkers.Mock, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, listDeactivationsVersion, note string, activations []edgeworkers.Activation, autoPin bool) {
	expectListActivations(m, edgeworkerID, "", activations, nil).Once()
	expectActivateVersion(m, edgeworkerID, activationID, net, version, note, autoPin, nil).Once()
	if listDeactivationsVersion != "" {
		expectListDeactivations(m, edgeworkerID, listDeactivationsVersion, []edgeworkers.Deactivation{}, nil).Once()
	}
	expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()
}

func expectFullRead(m *edgeworkers.Mock, edgeworkerID int, version string, acts []edgeworkers.Activation, deacts []edgeworkers.Deactivation, times int) {
	expectListActivations(m, edgeworkerID, "", acts, nil).Times(times)
	expectListDeactivations(m, edgeworkerID, version, deacts, nil).Times(times)
}

func expectFullDeactivation(m *edgeworkers.Mock, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, note string) {
	expectDeactivateVersion(m, edgeworkerID, deactivationID, net, version, note, nil).Once()
	expectGetDeactivation(m, edgeworkerID, deactivationID, net, version, activationStatusComplete, nil).Once()
}

func createStubActivation(edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status, time, note string) *edgeworkers.Activation {
	if time == "" {
		time = "2022-01-25T12:30:06Z"
	}
	return &edgeworkers.Activation{
		AccountID:        "testAccountId",
		ActivationID:     activationID,
		CreatedBy:        "unitTest",
		CreatedTime:      time,
		EdgeWorkerID:     edgeworkerID,
		LastModifiedTime: time,
		Network:          string(net),
		Status:           status,
		Version:          version,
		Note:             note,
	}
}

func createStubDeactivation(edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, status, time string) *edgeworkers.Deactivation {
	if time == "" {
		time = "2022-01-25T12:30:06Z"
	}
	return &edgeworkers.Deactivation{
		AccountID:        "testAccountId",
		DeactivationID:   deactivationID,
		CreatedBy:        "unitTest",
		CreatedTime:      time,
		EdgeWorkerID:     edgeworkerID,
		LastModifiedTime: time,
		Network:          net,
		Status:           status,
		Version:          version,
	}
}

func createStubEdgeworkerVersion(edgeworkerID int, version string) *edgeworkers.EdgeWorkerVersion {
	return &edgeworkers.EdgeWorkerVersion{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
		CreatedTime:  "2022-01-25T12:30:06Z",
	}
}

func TestUpgradeEdgeworkersActivationV1(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		rawState map[string]interface{}
		expected map[string]interface{}
	}{
		"auto_pin not present - should be added as true": {
			rawState: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "STAGING",
				"activation_id": 1,
				"note":          "some note",
			},
			expected: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "STAGING",
				"activation_id": 1,
				"note":          "some note",
				"auto_pin":      true,
			},
		},
		"auto_pin already present as true - should remain unchanged": {
			rawState: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "STAGING",
				"activation_id": 1,
				"auto_pin":      true,
			},
			expected: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "STAGING",
				"activation_id": 1,
				"auto_pin":      true,
			},
		},
		"auto_pin already present as false - should remain unchanged": {
			rawState: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "PRODUCTION",
				"activation_id": 2,
				"auto_pin":      false,
			},
			expected: map[string]interface{}{
				"edgeworker_id": 1234,
				"version":       "1",
				"network":       "PRODUCTION",
				"activation_id": 2,
				"auto_pin":      false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := upgradeEdgeworkersActivationV1(context.Background(), tc.rawState, nil)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

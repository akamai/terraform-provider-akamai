package edgeworkers

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceEdgeworkerActivation(t *testing.T) {
	workdir := "./testdata/TestResourceEdgeWorkersActivation"
	edgeworkerID := 1234

	tests := map[string]struct {
		init  func(*mockedgeworkers)
		steps []resource.TestStep
	}{
		"create and read activation - no previous activations": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - some previous activations, but no current": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 8
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 7, edgeworkers.ActivationNetworkProduction, "current", activationStatusComplete, "2022-01-25T12:30:06Z"),
					*createStubActivation(edgeworkerID, 6, net, "past2", activationStatusComplete, "2022-01-25T12:30:06Z"),
					*createStubActivation(edgeworkerID, 5, net, "past1", activationStatusComplete, "2022-01-24T12:30:06Z"),
					*createStubActivation(edgeworkerID, 4, edgeworkers.ActivationNetworkProduction, "past1", activationStatusComplete, "2022-01-23T12:30:06Z"),
					*createStubActivation(edgeworkerID, 3, net, "past2", activationStatusComplete, "2022-01-23T18:30:06Z"),
					*createStubActivation(edgeworkerID, 2, net, "past2", activationStatusComplete, "2022-01-23T12:30:06Z"),
					*createStubActivation(edgeworkerID, 1, net, "past1", activationStatusComplete, "2022-01-22T12:30:06Z"),
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
				expectActivateVersion(m, edgeworkerID, activationID, net, version, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, append([]edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, "")},
					activations...,
				), []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 4, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "8"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - version already active": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, nil).Once()
				expectListDeactivations(m, edgeworkerID, version, []edgeworkers.Deactivation{}, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - longer polling": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil)

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil).Once()

				// activate
				expectActivateVersion(m, edgeworkerID, activationID, net, version, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Times(2)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - version is already being activated, wait for activation": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil)

				// get current activation
				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusPending, ""),
				}, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"update network": {
			init: func(m *mockedgeworkers) {
				createNet, updateNet := edgeworkers.ActivationNetworkStaging, edgeworkers.ActivationNetworkProduction
				version := "test"
				createActivationID, updateActivationID := 1, 2

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, createNet, version)

				// read + plan + refresh
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 3)

				// update - activate
				expectActivateVersion(m, edgeworkerID, updateActivationID, updateNet, version, nil).Once()
				expectGetActivation(m, edgeworkerID, updateActivationID, updateNet, version, activationStatusComplete, nil).Once()

				// read + plan
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, updateNet, version, activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, createActivationID, createNet, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, updateNet, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_prod.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "2"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "PRODUCTION"),
					),
				},
			},
		},
		"update version": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 2

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion)

				// read + plan + refresh
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 3)

				// update - activate
				expectActivateVersion(m, edgeworkerID, updateActivationID, net, updateVersion, nil).Once()
				expectGetActivation(m, edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, nil).Once()

				// read + plan
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test1_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "2"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"update version - active version changed on refresh": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 3

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion"),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				expectFullRead(m, edgeworkerID, "someOtherVersion", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 2, net, "someOtherVersion", activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 1)

				// update - activate
				expectActivateVersion(m, edgeworkerID, updateActivationID, net, updateVersion, nil).Once()
				expectGetActivation(m, edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, nil).Once()

				// read + plan
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, 2, net, "someOtherVersion", activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test1_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "3"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"update version - version already active": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID, updateActivationID := 1, 2

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				expectFullRead(m, edgeworkerID, updateVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, updateActivationID, net, updateVersion, activationStatusComplete, ""),
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test1_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "2"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"destroy - version already deactivated": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version,
					fmt.Errorf("%w: %s", &edgeworkers.Error{ErrorCode: errorCodeVersionAlreadyDeactivated}, "oops"))
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"destroy - version is being deactivated, wait": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version,
					fmt.Errorf("%w: %s", &edgeworkers.Error{ErrorCode: errorCodeVersionIsBeingDeactivated}, "oops"))
				expectListDeactivations(m, edgeworkerID, version, []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 1, net, version, activationStatusInProgress, ""),
				}, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusComplete, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"destroy - longer polling": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				// version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, activationID, net, version)

				// read
				expectFullRead(m, edgeworkerID, version, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// test cleanup - destroy
				expectDeactivateVersion(m, edgeworkerID, 1, net, version, nil)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, 1, net, version, activationStatusComplete, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"error on create - missing required arguments": {
			init: func(m *mockedgeworkers) {},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"version\" is required"),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
			},
		},
		"error on create - version does not exist": {
			init: func(m *mockedgeworkers) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion2"),
					*createStubEdgeworkerVersion(edgeworkerID, "someOtherVersion1"),
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile(`version 'test' is not valid for edgeworker with id=1234`),
				},
			},
		},
		"error on create - getting current activation failed, ListActivations API error": {
			init: func(m *mockedgeworkers) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - getting current activation failed, ListDeactivations API error": {
			init: func(m *mockedgeworkers) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 1, edgeworkers.ActivationNetworkStaging, "test", activationStatusComplete, ""),
				}, nil)
				expectListDeactivations(m, edgeworkerID, "test", []edgeworkers.Deactivation{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on list version": {
			init: func(m *mockedgeworkers) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{}, fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on activate": {
			init: func(m *mockedgeworkers) {
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, 1, edgeworkers.ActivationNetworkStaging, "test", fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - API error on polling": {
			init: func(m *mockedgeworkers) {
				version := "test"
				net := edgeworkers.ActivationNetworkStaging
				activationID := 1
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, activationID, net, version, nil)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, "", fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: oops"),
				},
			},
		},
		"error on create - activation failed": {
			init: func(m *mockedgeworkers) {
				version := "test"
				net := edgeworkers.ActivationNetworkStaging
				activationID := 1
				// create - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, version),
				}, nil).Once()

				expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil)
				expectActivateVersion(m, edgeworkerID, activationID, net, version, nil)
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, activationID, net, version, "ERROR", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation: edgeworker activation failure"),
				},
			},
		},
		"error on update - version does not exist": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID := 1

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, "test"),
				}, nil).Times(2)

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion)

				// read + plan + refresh
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 3)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test1_stag.tf", workdir)),
					ExpectError: regexp.MustCompile(`version 'test1' is not valid for edgeworker with id=1234`),
				},
			},
		},
		"error on update - no current activation on refresh": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				createActivationID := 1

				// create + update - version verification
				expectListEdgeWorkerVersions(m, edgeworkerID, []edgeworkers.EdgeWorkerVersion{
					*createStubEdgeworkerVersion(edgeworkerID, updateVersion),
					*createStubEdgeworkerVersion(edgeworkerID, createVersion),
				}, nil).Once()

				// create
				expectFullActivation(m, edgeworkerID, createActivationID, net, createVersion)

				// read + plan
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, ""),
				}, []edgeworkers.Deactivation{}, 2)

				// refresh
				expectFullRead(m, edgeworkerID, createVersion, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, createActivationID, net, createVersion, activationStatusComplete, "2022-01-25T12:30:06Z"),
				}, []edgeworkers.Deactivation{
					*createStubDeactivation(edgeworkerID, 1, net, createVersion, activationStatusComplete, "2022-01-26T12:30:06Z"),
				}, 1)

				// test cleanup - destroy
				expectFullDeactivation(m, edgeworkerID, 1, net, createVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworkers_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworkers_activation_version_test1_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation read: no version active on network 'STAGING' for edgeworker with id=1234"),
				},
			},
		},
	}

	// redefining times to accelerate tests
	activationPollMinimum = time.Millisecond
	activationPollInterval = activationPollMinimum

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockedgeworkers{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectActivateVersion(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version string, e error) *mock.Call {
	req := edgeworkers.ActivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		ActivateVersion: edgeworkers.ActivateVersion{
			Network: net,
			Version: version,
		},
	}
	if e != nil {
		return m.On("ActivateVersion", mock.Anything, req).Return(nil, e)
	}

	return m.On("ActivateVersion", mock.Anything, req).Return(createStubActivation(edgeworkerID, activationID, net, version, activationStatusPresubmit, ""), nil)
}

func expectGetActivation(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetActivationRequest{
		EdgeWorkerID: edgeworkerID,
		ActivationID: activationID,
	}
	if e != nil {
		return m.On("GetActivation", mock.Anything, req).Return(nil, e)
	}

	return m.On("GetActivation", mock.Anything, req).Return(createStubActivation(edgeworkerID, activationID, net, version, status, ""), nil)
}

func expectListActivations(m *mockedgeworkers, edgeworkerID int, version string, activations []edgeworkers.Activation, e error) *mock.Call {
	req := edgeworkers.ListActivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	}
	if e != nil {
		return m.On("ListActivations", mock.Anything, req).Return(nil, e)
	}

	return m.On("ListActivations", mock.Anything, req).Return(&edgeworkers.ListActivationsResponse{
		Activations: activations,
	}, nil)
}

func expectListDeactivations(m *mockedgeworkers, edgeworkerID int, version string, deactivations []edgeworkers.Deactivation, e error) *mock.Call {
	req := edgeworkers.ListDeactivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	}
	if e != nil {
		return m.On("ListDeactivations", mock.Anything, req).Return(nil, e)
	}

	return m.On("ListDeactivations", mock.Anything, req).Return(&edgeworkers.ListDeactivationsResponse{
		Deactivations: deactivations,
	}, nil)
}

func expectDeactivateVersion(m *mockedgeworkers, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version string, e error) *mock.Call {
	req := edgeworkers.DeactivateVersionRequest{
		EdgeWorkerID: edgeworkerID,
		DeactivateVersion: edgeworkers.DeactivateVersion{
			Network: net,
			Version: version,
		},
	}
	if e != nil {
		return m.On("DeactivateVersion", mock.Anything, req).Return(nil, e)
	}

	return m.On("DeactivateVersion", mock.Anything, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, activationStatusPresubmit, ""), nil)
}

func expectGetDeactivation(m *mockedgeworkers, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetDeactivationRequest{
		EdgeWorkerID:   edgeworkerID,
		DeactivationID: deactivationID,
	}
	if e != nil {
		return m.On("GetDeactivation", mock.Anything, req).Return(nil, e)
	}

	return m.On("GetDeactivation", mock.Anything, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, status, ""), nil)
}

func expectListEdgeWorkerVersions(m *mockedgeworkers, edgeworkerID int, versions []edgeworkers.EdgeWorkerVersion, e error) *mock.Call {
	req := edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeworkerID,
	}
	if e != nil {
		return m.On("ListEdgeWorkerVersions", mock.Anything, req).Return(nil, e)
	}

	return m.On("ListEdgeWorkerVersions", mock.Anything, req).Return(&edgeworkers.ListEdgeWorkerVersionsResponse{
		EdgeWorkerVersions: versions,
	}, nil)
}

func expectFullActivation(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version string) {
	expectListActivations(m, edgeworkerID, "", []edgeworkers.Activation{}, nil).Once()
	expectActivateVersion(m, edgeworkerID, activationID, net, version, nil).Once()
	expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()
}

func expectFullRead(m *mockedgeworkers, edgeworkerID int, version string, acts []edgeworkers.Activation, deacts []edgeworkers.Deactivation, times int) {
	expectListActivations(m, edgeworkerID, "", acts, nil).Times(times)
	expectListDeactivations(m, edgeworkerID, version, deacts, nil).Times(times)
}

func expectFullDeactivation(m *mockedgeworkers, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version string) {
	expectDeactivateVersion(m, edgeworkerID, deactivationID, net, version, nil).Once()
	expectGetDeactivation(m, edgeworkerID, deactivationID, net, version, activationStatusComplete, nil).Once()
}

func createStubActivation(edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status, time string) *edgeworkers.Activation {
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

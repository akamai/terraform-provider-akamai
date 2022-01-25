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
	workdir := "./testdata/TestResEdgeWorkerActivation"
	edgeworkerID := 1234

	tests := map[string]struct {
		init  func(*mockedgeworkers)
		steps []resource.TestStep
	}{
		"create and read activation - network staging": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"

				expectFullCreateActivationPhase(m, edgeworkerID, 1, net, version)

				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - network production": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkProduction
				version := "test"

				expectFullCreateActivationPhase(m, edgeworkerID, 1, net, version)

				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_prod.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "PRODUCTION"),
					),
				},
			},
		},
		"create and read activation - longer polling": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"

				expectActivateVersion(m, edgeworkerID, 1, net, version, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Times(2)
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Times(2)
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusComplete, nil).Times(1)
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{*createStubActivation(edgeworkerID, 1, net, version, activationStatusComplete)}).Times(2)

				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and read activation - version already active": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activateErr := &edgeworkers.Error{ErrorCode: errorCodeVersionAlreadyActive}
				activations := []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 4, edgeworkers.ActivationNetworkProduction, version, activationStatusComplete),
					*createStubActivation(edgeworkerID, 3, net, version, activationStatusComplete),
					*createStubActivation(edgeworkerID, 2, edgeworkers.ActivationNetworkProduction, version, activationStatusComplete),
					*createStubActivation(edgeworkerID, 1, net, version, activationStatusComplete),
				}
				expectActivateVersion(m, edgeworkerID, 1, net, version, activateErr).Once()
				expectListActivations(m, edgeworkerID, version, nil, activations).Once()

				expectListActivations(m, edgeworkerID, "", nil, activations).Times(2)

				expectFullDeactivation(m, edgeworkerID, 1, net, version)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "3"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and update activation - network staging": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				// create and read
				expectFullCreateActivationPhase(m, edgeworkerID, 1, net, createVersion)

				// pre-update read
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{*createStubActivation(edgeworkerID, 1, net, createVersion, activationStatusComplete)}).Once()

				// update and destroy
				expectFullActivation(m, edgeworkerID, 2, net, updateVersion)
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 2, net, updateVersion, activationStatusComplete),
					*createStubActivation(edgeworkerID, 1, net, createVersion, activationStatusComplete),
				}).Times(2)

				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test1_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "2"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"create and update activation - version changed in the meantime": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, changedVersion, updateVersion := "test", "someOtherVersion", "test1"
				// create and read
				expectFullCreateActivationPhase(m, edgeworkerID, 1, net, createVersion)

				// pre-update read
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 2, net, changedVersion, activationStatusComplete),
					*createStubActivation(edgeworkerID, 1, net, createVersion, activationStatusComplete),
				}).Once()

				// update and destroy
				expectFullActivation(m, edgeworkerID, 3, net, updateVersion)
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{
					*createStubActivation(edgeworkerID, 3, net, updateVersion, activationStatusComplete),
					*createStubActivation(edgeworkerID, 2, net, changedVersion, activationStatusComplete),
					*createStubActivation(edgeworkerID, 1, net, createVersion, activationStatusComplete),
				}).Times(2)
				
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test1_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "3"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"destroy activation - longer polling": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				activationID := 1

				expectFullCreateActivationPhase(m, edgeworkerID, activationID, net, version)

				expectDeactivateVersion(m, edgeworkerID, activationID, net, version, nil).Once()
				expectGetDeactivation(m, edgeworkerID, activationID, net, version, activationStatusPresubmit, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, activationID, net, version, activationStatusPending, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, activationID, net, version, activationStatusInProgress, nil).Times(2)
				expectGetDeactivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
			},
		},
		"error creating activation - API error": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				expectActivateVersion(m, edgeworkerID, 1, net, version, &edgeworkers.Error{}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation create: API error"),
				},
			},
		},
		"error creating activation - missing required arguments": {
			init: func(m *mockedgeworkers) {},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"version\" is required"),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_missing_required_args.tf", workdir)),
					ExpectError: regexp.MustCompile("argument \"edgeworker_id\" is required"),
				},
			},
		},
		"error creating activation - activation failed": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				expectActivateVersion(m, edgeworkerID, 1, net, version, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, "ERROR", nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation create: edgeworker activation failure"),
				},
			},
		},
		"error creating activation - error while polling": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				version := "test"
				expectActivateVersion(m, edgeworkerID, 1, net, version, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPresubmit, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusPending, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, activationStatusInProgress, nil).Once()
				expectGetActivation(m, edgeworkerID, 1, net, version, "", fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error updating activation - activation failed": {
			init: func(m *mockedgeworkers) {
				net := edgeworkers.ActivationNetworkStaging
				createVersion, updateVersion := "test", "test1"
				activationID := 1

				expectFullCreateActivationPhase(m, edgeworkerID, activationID, net, createVersion)

				// pre-update read
				expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{*createStubActivation(edgeworkerID, 1, net, createVersion, activationStatusComplete)}).Once()

				expectActivateVersion(m, edgeworkerID, activationID, net, updateVersion, fmt.Errorf("oops"))

				// destroy after the test needed so the test does not fail
				expectFullDeactivation(m, edgeworkerID, 1, net, updateVersion)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test_stag.tf", workdir)),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "activation_id", "1"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "version", "test"),
						resource.TestCheckResourceAttr("akamai_edgeworker_activation.test", "network", "STAGING"),
					),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/edgeworker_activation_version_test1_stag.tf", workdir)),
					ExpectError: regexp.MustCompile("edgeworker activation update: oops"),
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

	return m.On("ActivateVersion", mock.Anything, req).Return(createStubActivation(edgeworkerID, activationID, net, version, activationStatusPresubmit), nil)
}

func expectFullActivation(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version string) {
	expectActivateVersion(m, edgeworkerID, activationID, net, version, nil).Once()
	expectGetActivation(m, edgeworkerID, activationID, net, version, activationStatusComplete, nil).Once()
}

func expectFullCreateActivationPhase(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version string) {
	expectFullActivation(m, edgeworkerID, activationID, net, version)
	expectListActivations(m, edgeworkerID, "", nil, []edgeworkers.Activation{*createStubActivation(edgeworkerID, activationID, net, version, activationStatusComplete)}).Times(2)
}

func expectGetActivation(m *mockedgeworkers, edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetActivationRequest{
		EdgeWorkerID: edgeworkerID,
		ActivationID: activationID,
	}
	if e != nil {
		return m.On("GetActivation", mock.Anything, req).Return(nil, e)
	}

	return m.On("GetActivation", mock.Anything, req).Return(createStubActivation(edgeworkerID, activationID, net, version, status), nil)
}

func expectListActivations(m *mockedgeworkers, edgeworkerID int, version string, e error, activations []edgeworkers.Activation) *mock.Call {
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

	return m.On("DeactivateVersion", mock.Anything, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, activationStatusPresubmit), nil)
}

func expectGetDeactivation(m *mockedgeworkers, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, status string, e error) *mock.Call {
	req := edgeworkers.GetDeactivationRequest{
		EdgeWorkerID:   edgeworkerID,
		DeactivationID: deactivationID,
	}
	if e != nil {
		return m.On("GetDeactivation", mock.Anything, req).Return(nil, e)
	}

	return m.On("GetDeactivation", mock.Anything, req).Return(createStubDeactivation(edgeworkerID, deactivationID, net, version, status), nil)
}

func expectFullDeactivation(m *mockedgeworkers, edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version string) {
	expectDeactivateVersion(m, edgeworkerID, deactivationID, net, version, nil).Once()
	expectGetDeactivation(m, edgeworkerID, deactivationID, net, version, activationStatusComplete, nil).Once()
}

func createStubActivation(edgeworkerID, activationID int, net edgeworkers.ActivationNetwork, version, status string) *edgeworkers.Activation {
	return &edgeworkers.Activation{
		AccountID:        "testAccountId",
		ActivationID:     activationID,
		CreatedBy:        "unitTest",
		CreatedTime:      "now",
		EdgeWorkerID:     edgeworkerID,
		LastModifiedTime: "now",
		Network:          string(net),
		Status:           status,
		Version:          version,
	}
}

func createStubDeactivation(edgeworkerID, deactivationID int, net edgeworkers.ActivationNetwork, version, status string) *edgeworkers.Deactivation {
	return &edgeworkers.Deactivation{
		AccountID:        "testAccountId",
		DeactivationID:   deactivationID,
		CreatedBy:        "unitTest",
		CreatedTime:      "now",
		EdgeWorkerID:     edgeworkerID,
		LastModifiedTime: "now",
		Network:          net,
		Status:           status,
		Version:          version,
	}
}

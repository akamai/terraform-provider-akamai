package apidefinitions

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions/v0"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
)

func TestActivationResource(t *testing.T) {
	t.Parallel()
	pollInterval = time.Millisecond * 10

	var tests = map[string]struct {
		configPath   string
		init         func(*apidefinitions.Mock)
		checkDestroy func(*terraform.State) error
		steps        []resource.TestStep
		error        *regexp.Regexp
	}{
		"activation - on staging": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{})
				mockActivateVersion(m)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 4)
				mockDeactivateVersion(m, 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)
			},
			steps: []resource.TestStep{
				{
					Config: activationConfig(1),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_activation.a1", "status", "ACTIVE"),
					),
				},
			},
		},
		"activation - on both networks": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{})
				mockActivateVersion(m)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 3)
				mockDeactivateVersion(m, 1)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)

			},
			steps: []resource.TestStep{
				{
					Config: activationConfig(1),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_activation.a1", "status", "ACTIVE"),
					),
				},
			},
		},
		"activation - version is already active": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 4)
				mockDeactivateVersion(m, 1)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)
			},
			steps: []resource.TestStep{
				{
					Config: activationConfig(1),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_activation.a1", "status", "ACTIVE"),
					),
				},
			},
		},
		"activation - failed": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{})
				mockActivateVersion(m)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusFailed), 1)
			},
			steps: []resource.TestStep{
				{
					Config:      activationConfig(1),
					ExpectError: regexp.MustCompile("Activation for version 1 failed"),
				},
			},
		},
		"activation - errors during verify": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{{Severity: apidefinitions.SeverityError, Detail: "You shall not pass"}})
			},
			steps: []resource.TestStep{
				{
					Config:      activationConfig(1),
					ExpectError: regexp.MustCompile("Unable to proceed due to Activation Errors, fix errors to proceed"),
				},
			},
		},
		"activation - warnings during verify": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{{Severity: apidefinitions.SeverityWarning, Detail: "You shall not pass"}})
			},
			steps: []resource.TestStep{
				{
					Config:      activationConfig(1),
					ExpectError: regexp.MustCompile("Unable to proceed due to Activation Warnings"),
				},
			},
		},
		"activation - auto ack warnings": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
				mockVerifyVersion(m, []apidefinitions.VerifyVersionAlert{{Severity: apidefinitions.SeverityWarning, Detail: "You shall not pass"}})
				mockActivateVersion(m)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 4)
				mockDeactivateVersion(m, 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)
			},
			steps: []resource.TestStep{
				{
					Config: activationConfigWithAutoAck(1, true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_activation.a1", "status", "ACTIVE"),
					),
				},
			},
		},
		"check schema - missing required attributes": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/activation/activation_invalid.tf"),
					ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/activation/activation_invalid.tf"),
					ExpectError: regexp.MustCompile(`The argument "api_id" is required, but no definition was found.`),
				},
			},
		},
		"check schema - invalid email": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/activation/activation_invalid_email.tf"),
					ExpectError: regexp.MustCompile("is not valid: mail"),
				},
			},
		},
		"import - ok": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 3)
				mockDeactivateVersion(m, 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusPending), 1)
				mockGetEndpointWithStagingActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)
			},
			steps: []resource.TestStep{
				{
					Config:        activationImportConfig(),
					ImportState:   true,
					ImportStateId: "12345:STAGING",
					ResourceName:  "akamai_apidefinitions_activation.import_test",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						state := states[0].Attributes
						assert.Equal(t, "12345", state["api_id"])
						assert.Equal(t, "1", state["version"])
						assert.Equal(t, string(apidefinitions.ActivationNetworkStaging), state["network"])
						assert.Equal(t, string(apidefinitions.ActivationStatusActive), state["status"])
						return nil
					},
					ImportStatePersist: true,
				},
			},
		},
		"import - not active": {
			init: func(m *apidefinitions.Mock) {
				mockGetEndpointWithStagingActivationStatus(m, nil, nil, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             activationImportConfig(),
					ImportState:        true,
					ImportStateId:      "12345:STAGING",
					ResourceName:       "akamai_apidefinitions_activation.import_test",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("API is not active on the network STAGING"),
				},
			},
		},
		"import - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             activationImportConfig(),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       "akamai_apidefinitions_activation.import_test",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: ID '12345' incorrectly formatted: should be 'API_ID:NETWORK'"),
				},
			},
		},
		"import - invalid network": {
			steps: []resource.TestStep{
				{
					Config:                               activationImportConfig(),
					ImportState:                          true,
					ImportStateId:                        "12345:NETWORK",
					ResourceName:                         "akamai_apidefinitions_activation.import_test",
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("invalid network value NETWORK; must be either STAGING or PRODUCTION"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &apidefinitions.Mock{}
			clientV0 := &v0.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, clientV0, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
					CheckDestroy:             test.checkDestroy,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockGetEndpointWithStagingActivationStatus(client *apidefinitions.Mock, version *int64, status *apidefinitions.ActivationStatus, times int) {
	client.On("GetEndpoint", mock.Anything, mock.Anything).
		Return(&apidefinitions.GetEndpointResponse{
			Endpoint: apidefinitions.Endpoint{
				StagingVersion: apidefinitions.VersionState{
					VersionNumber: version,
					Status:        status,
				},
				Locked: false,
			},
			SecurityScheme:             nil,
			AkamaiSecurityRestrictions: nil,
		}, nil).Times(times)
}

func mockGetEndpointWithActivationStatus(client *apidefinitions.Mock, version *int64, status *apidefinitions.ActivationStatus, times int) {
	client.On("GetEndpoint", mock.Anything, mock.Anything).
		Return(&apidefinitions.GetEndpointResponse{
			Endpoint: apidefinitions.Endpoint{
				StagingVersion: apidefinitions.VersionState{
					VersionNumber: version,
					Status:        status,
				},
				ProductionVersion: apidefinitions.VersionState{
					VersionNumber: version,
					Status:        status,
				},
				Locked: false,
			},
		}, nil).
		Times(times)
}

func mockVerifyVersion(client *apidefinitions.Mock, alerts apidefinitions.VerifyVersionResponse) {
	client.On("VerifyVersion", mock.Anything, mock.Anything).
		Return(alerts, nil).
		Once()
}

func mockActivateVersion(client *apidefinitions.Mock) {
	client.On("ActivateVersion", mock.Anything, mock.Anything).
		Return(&apidefinitions.ActivateVersionResponse{}, nil).
		Once()
}

func mockActivateVersionFail(client *apidefinitions.Mock) {
	client.On("ActivateVersion", mock.Anything, mock.Anything).
		Return(nil, &apidefinitions.Error{
			Status: 500,
		}).
		Once()
}

func mockDeactivateVersionFail(client *apidefinitions.Mock) {
	client.On("DeactivateVersion", mock.Anything, mock.Anything).
		Return(nil, &apidefinitions.Error{
			Status: 500,
		}).
		Once()
}

func mockDeactivateVersion(client *apidefinitions.Mock, times int) {
	client.On("DeactivateVersion", mock.Anything, mock.Anything).
		Return(&apidefinitions.DeactivateVersionResponse{}, nil).
		Times(times)
}

func activationConfig(version int64) string {
	return activationConfigWithAutoAck(version, false)
}

func activationConfigWithAutoAck(version int64, autoAck bool) string {
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_activation" "a1" {
  api_id      = 1
  version     = %v
  network     = "STAGING"
  notification_recipients = ["user@example.com"]
  notes       = "Notes"
  auto_acknowledge_warnings = %v
}
`, version, autoAck)
}

func activationImportConfig() string {
	return providerConfig + `
resource "akamai_apidefinitions_activation" "import_test" {
  api_id      = 1
  version     = 1
  network     = "STAGING"
}
`
}

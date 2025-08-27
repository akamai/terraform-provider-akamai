package clientlists

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
)

func TestResourceClientListActivation(t *testing.T) {
	t.Parallel()
	pollActivationInterval = time.Microsecond
	activationRetryBaseDelay = time.Microsecond

	const testDir = "testData/TestResActivation"

	activationReq := clientlists.CreateActivationRequest{
		ListID: "12_AB",
		ActivationParams: clientlists.ActivationParams{
			Action:                 clientlists.Activate,
			Network:                clientlists.Staging,
			SiebelTicketID:         "ABC-12345",
			NotificationRecipients: []string{"user@example.com"},
			Comments:               "Activation Comments",
		},
	}

	deactivationReq := clientlists.CreateDeactivationRequest{
		ListID: "12_AB",
		ActivationParams: clientlists.ActivationParams{
			Action:                 clientlists.Deactivate,
			Network:                clientlists.Staging,
			SiebelTicketID:         "ABC-12345",
			NotificationRecipients: []string{"user@example.com"},
			Comments:               "Activation Comments",
		},
	}

	apiBadRequestError := clientlists.Error{
		Type:       "https://problems.luna.akamaiapis.net/client-list-api/error-types/INVALID-INPUT-ERROR",
		Title:      "Invalid Input Error",
		Detail:     "Validation failed: Invalid network",
		StatusCode: http.StatusBadRequest,
		Instance:   "https://problems.luna.akamaiapis.net/client-list-api/error-instances/9ff3649993cb002b",
	}

	apiServerError := clientlists.Error{
		Type:       "https://problems.luna.akamaiapis.net/client-list/error-types/INTERNAL-SERVER-ERROR",
		Title:      "Internal Server Error",
		Detail:     "Error occurred activating client list",
		StatusCode: http.StatusInternalServerError,
		Instance:   "https://problems.luna.akamaiapis.net/client-list/error-instances/c3411b55b7bfd724",
	}

	resourceName := "akamai_clientlist_activation.activation_ASN_LIST_1"
	baseChecker := test.NewStateChecker(resourceName).
		CheckEqual("list_id", "12_AB").
		CheckEqual("network", "STAGING")

	var tests = map[string]struct {
		init  func(*clientlists.Mock)
		steps []resource.TestStep
	}{
		"create activation": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 3)
				mockGetClientlist(m, "12_AB", 2, 2)

				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
			},
		},
		"update activation - version and other fields update": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				updatedActivationRes := mockCreateActivation(m, clientlists.CreateActivationRequest{
					ListID: "12_AB",
					ActivationParams: clientlists.ActivationParams{
						Action:                 clientlists.Activate,
						Network:                clientlists.Staging,
						SiebelTicketID:         "UPDATED-12345",
						NotificationRecipients: []string{"update_user@example.com"},
						Comments:               "Activation Comments Updated",
					},
				}, 3, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 4)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
					getActivationAttrs(updatedActivationRes, clientlists.PendingActivation), 2)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
					getActivationAttrs(updatedActivationRes, clientlists.Active), 3)
				mockGetClientlist(m, "12_AB", 2, 3)
				mockGetClientlist(m, "12_AB", 3, 2)

				mockDestroyResource(m, clientlists.CreateDeactivationRequest{
					ListID: "12_AB",
					ActivationParams: clientlists.ActivationParams{
						Action:                 clientlists.Deactivate,
						Network:                clientlists.Staging,
						SiebelTicketID:         "UPDATED-12345",
						NotificationRecipients: []string{"update_user@example.com"},
						Comments:               "Activation Comments Updated",
					},
				}, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_update.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "UPDATED-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments Updated").
						CheckEqual("version", "3").
						Build(),
				},
			},
		},
		"update activation - version only update": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				updatedActivationRes := mockCreateActivation(m, activationReq, 3, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 4)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
					getActivationAttrs(updatedActivationRes, clientlists.PendingActivation), 2)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
					getActivationAttrs(updatedActivationRes, clientlists.Active), 3)
				mockGetClientlist(m, "12_AB", 2, 3)
				mockGetClientlist(m, "12_AB", 3, 2)

				mockDestroyResource(m, deactivationReq, 3, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_update_version_only.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "3").
						Build(),
				},
			},
		},
		"update activation - notification_recipients, siebel_ticket_id and comments updates suppressed": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 2)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 5)
				mockGetClientlist(m, "12_AB", 2, 4)

				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_update_suppressed.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
			},
		},
		"create activation with missing list id fails": {
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/activation_missing_param.tf", testDir)),
					ExpectError: regexp.MustCompile("The argument \"list_id\" is required, but no definition was found"),
				},
			},
		},
		"create activation - api fails without retry for not 500 error": {
			init: func(m *clientlists.Mock) {
				mockAPIErrorWithCreateActivation(m, activationReq, apiBadRequestError, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					ExpectError: regexp.MustCompile("Error: Title: Invalid Input Error; Type: https://problems.luna.akamaiapis.net/client-list-api/error-types/INVALID-INPUT-ERROR; Detail: Validation failed: Invalid network"),
				},
			},
		},
		"create activation - api fails after retry for 500 error": {
			init: func(m *clientlists.Mock) {
				mockAPIErrorWithCreateActivation(m, activationReq, apiServerError, 3)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					ExpectError: regexp.MustCompile("Error: Title: Internal Server Error; Type: https://problems.luna.akamaiapis.net/client-list/error-types/INTERNAL-SERVER-ERROR; Detail: Error occurred activating client list"),
				},
			},
		},
		"import activation resource": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 3)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 4)
				mockGetClientlist(m, "12_AB", 2, 3)
				mockGetActivationStatus(m, clientlists.GetActivationStatusRequest{
					Network: clientlists.Staging,
					ListID:  activationReq.ListID,
				}, getActivationAttrs(activationRes, clientlists.Active), 1)

				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
				},
				{
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateId:     "12_AB:STAGING",
					ResourceName:      resourceName,
				},
			},
		},
		"delete activation": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 3)
				mockGetClientlist(m, "12_AB", 2, 3)

				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/delete_activation.tf", testDir)),
					Check: func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceName]
						if ok {
							return fmt.Errorf("resource %s is still present in the Terraform state", resourceName)
						}
						return nil
					},
				},
			},
		},
		"delete activation - check retry work for 500 error": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 3)
				mockGetClientlist(m, "12_AB", 2, 3)

				mockAPIErrorWithCreateDeactivation(m, deactivationReq, apiServerError, 1)
				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config: loadFixtureString(fmt.Sprintf("%s/delete_activation.tf", testDir)),
					Check: func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceName]
						if ok {
							return fmt.Errorf("resource %s is still present in the Terraform state", resourceName)
						}
						return nil
					},
				},
			},
		},
		"delete activation - api fails after retry": {
			init: func(m *clientlists.Mock) {
				activationRes := mockCreateActivation(m, activationReq, 2, 33)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.PendingActivation), 1)
				mockReadActivation(m,
					clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
					getActivationAttrs(activationRes, clientlists.Active), 4)
				mockGetClientlist(m, "12_AB", 2, 3)

				mockAPIErrorWithCreateDeactivation(m, deactivationReq, apiServerError, 3)
				mockDestroyResource(m, deactivationReq, 2, 33)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					Check: baseChecker.
						CheckEqual("siebel_ticket_id", "ABC-12345").
						CheckEqual("notification_recipients.#", "1").
						CheckEqual("comments", "Activation Comments").
						CheckEqual("version", "2").
						Build(),
				},
				{
					Config:      loadFixtureString(fmt.Sprintf("%s/delete_activation.tf", testDir)),
					ExpectError: regexp.MustCompile("Error: Title: Internal Server Error; Type: https://problems.luna.akamaiapis.net/client-list/error-types/INTERNAL-SERVER-ERROR; Detail: Error occurred activating client list"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &clientlists.Mock{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockDestroyResource(m *clientlists.Mock, deactivationReq clientlists.CreateDeactivationRequest, version int64, activationID int64) {
	deactivationRes := mockCreateDeactivation(m, deactivationReq, version, activationID)
	mockReadActivation(m,
		clientlists.GetActivationRequest{ActivationID: deactivationRes.ActivationID},
		getDeactivationAttrs(deactivationRes, clientlists.PendingDeactivation), 1)
	mockReadActivation(m,
		clientlists.GetActivationRequest{ActivationID: deactivationRes.ActivationID},
		getDeactivationAttrs(deactivationRes, clientlists.Deactivated), 1)
}

type ActivationAttrs struct {
	ListID, Comments, SiebelTicketID, Network string
	NotificationRecipients                    []string
	Version, ActivationID                     int
	Status                                    clientlists.ActivationStatus
}

func mockCreateActivation(m *clientlists.Mock, req clientlists.CreateActivationRequest, version int64, activationID int64) *clientlists.CreateActivationResponse {
	res := clientlists.CreateActivationResponse{
		ActivationID:           activationID,
		ListID:                 req.ListID,
		Comments:               req.Comments,
		SiebelTicketID:         req.SiebelTicketID,
		Network:                req.Network,
		NotificationRecipients: req.NotificationRecipients,
		ActivationStatus:       clientlists.PendingActivation,
		Version:                version,
	}

	m.On("CreateActivation", mock.Anything, req).Return(&res, nil).Once()

	return &res
}

func mockCreateDeactivation(m *clientlists.Mock, req clientlists.CreateDeactivationRequest, version int64, activationID int64) *clientlists.CreateDeactivationResponse {
	res := clientlists.CreateDeactivationResponse{
		ActivationID:           activationID,
		ListID:                 req.ListID,
		Comments:               req.Comments,
		SiebelTicketID:         req.SiebelTicketID,
		Network:                req.Network,
		NotificationRecipients: req.NotificationRecipients,
		ActivationStatus:       clientlists.PendingActivation,
		Version:                version,
	}

	m.On("CreateDeactivation", mock.Anything, req).Return(&res, nil).Once()

	return &res
}

func mockReadActivation(m *clientlists.Mock, req clientlists.GetActivationRequest, attrs ActivationAttrs, times int) *clientlists.GetActivationResponse {
	res := clientlists.GetActivationResponse{
		ActivationID: int64(attrs.ActivationID),
		ListID:       attrs.ListID,
		Version:      int64(attrs.Version),
		ActivationParams: clientlists.ActivationParams{
			Action:                 clientlists.Activate,
			Comments:               attrs.Comments,
			Network:                clientlists.ActivationNetwork(attrs.Network),
			SiebelTicketID:         attrs.SiebelTicketID,
			NotificationRecipients: attrs.NotificationRecipients,
		},
		InitialActivation: true,
		Fast:              true,
		ActivationStatus:  attrs.Status,
	}

	m.On("GetActivation", mock.Anything, req).Return(&res, nil).Times(times)

	return &res
}

func mockGetActivationStatus(m *clientlists.Mock, req clientlists.GetActivationStatusRequest, attrs ActivationAttrs, times int) *clientlists.GetActivationStatusResponse {
	res := clientlists.GetActivationStatusResponse{
		Action:                 clientlists.Activate,
		ActivationID:           int64(attrs.ActivationID),
		ActivationStatus:       attrs.Status,
		Comments:               attrs.Comments,
		ListID:                 attrs.ListID,
		Network:                clientlists.ActivationNetwork(attrs.Network),
		NotificationRecipients: attrs.NotificationRecipients,
		SiebelTicketID:         attrs.SiebelTicketID,
		Version:                int64(attrs.Version),
	}

	m.On("GetActivationStatus", mock.Anything, req).Return(&res, nil).Times(times)

	return &res
}

func mockGetClientlist(m *clientlists.Mock, listID string, version int64, callTimes int) {
	clientListGetReq := clientlists.GetClientListRequest{
		ListID:       listID,
		IncludeItems: false,
	}

	clientList := clientlists.GetClientListResponse{ListContent: clientlists.ListContent{Version: version}}
	m.On("GetClientList", mock.Anything, clientListGetReq).Return(&clientList, nil).Times(callTimes)
}

func mockAPIErrorWithCreateActivation(m *clientlists.Mock, req clientlists.CreateActivationRequest, apiError clientlists.Error, times int) {
	m.On("CreateActivation", mock.Anything, req).Return(nil, &apiError).Times(times)
}

func mockAPIErrorWithCreateDeactivation(m *clientlists.Mock, req clientlists.CreateDeactivationRequest, apiError clientlists.Error, times int) {
	m.On("CreateDeactivation", mock.Anything, req).Return(nil, &apiError).Times(times)
}

func getActivationAttrs(actRes *clientlists.CreateActivationResponse, status clientlists.ActivationStatus) ActivationAttrs {
	return ActivationAttrs{
		ActivationID:           int(actRes.ActivationID),
		SiebelTicketID:         actRes.SiebelTicketID,
		Network:                string(actRes.Network),
		Comments:               actRes.Comments,
		NotificationRecipients: actRes.NotificationRecipients,
		ListID:                 actRes.ListID,
		Version:                int(actRes.Version),
		Status:                 status,
	}
}

func getDeactivationAttrs(actRes *clientlists.CreateDeactivationResponse, status clientlists.ActivationStatus) ActivationAttrs {
	return ActivationAttrs{
		ActivationID:           int(actRes.ActivationID),
		SiebelTicketID:         actRes.SiebelTicketID,
		Network:                string(actRes.Network),
		Comments:               actRes.Comments,
		NotificationRecipients: actRes.NotificationRecipients,
		ListID:                 actRes.ListID,
		Version:                int(actRes.Version),
		Status:                 status,
	}
}

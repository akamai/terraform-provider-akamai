package clientlists

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientListActivationResource(t *testing.T) {
	pollActivationInterval = time.Microsecond

	type ActivationAttrs struct {
		ListID, Comments, SiebelTicketID, Network string
		NotificationRecipients                    []string
		Version, ActivationID                     int
		Status                                    clientlists.ActivationStatus
	}

	const testDir = "testData/TestResActivation"

	var (
		createActivationAPIError = "{\n      \"type\": \"https://problems.luna.akamaiapis.net/client-list-api/error-types/INVALID-INPUT-ERROR\",\n       \"status\": 400,\n       \"title\": \"Invalid Input Error\",\n       \"detail\": \"Validation failed: Invalid network\",\n       \"instance\": \"https://problems.luna.akamaiapis.net/client-list-api/error-instances/9ff3649993cb002b\",\n       \"\n }\n"

		expectCreateActivation = func(client *clientlists.Mock, req clientlists.CreateActivationRequest, version int64, activationID int64) *clientlists.CreateActivationResponse {
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

			client.On("CreateActivation", testutils.MockContext, req).Return(&res, nil).Once()

			return &res
		}

		expectReadActivation = func(client *clientlists.Mock, req clientlists.GetActivationRequest, attrs ActivationAttrs, times int) *clientlists.GetActivationResponse {
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
				ActivationStatus:  clientlists.ActivationStatus(attrs.Status),
			}

			client.On("GetActivation", testutils.MockContext, req).Return(&res, nil).Times(times)

			return &res
		}

		expectGetActivationStatus = func(client *clientlists.Mock, req clientlists.GetActivationStatusRequest, attrs ActivationAttrs, times int) *clientlists.GetActivationStatusResponse {
			res := clientlists.GetActivationStatusResponse{
				Action:                 clientlists.Activate,
				ActivationID:           int64(attrs.ActivationID),
				ActivationStatus:       clientlists.ActivationStatus(attrs.Status),
				Comments:               attrs.Comments,
				ListID:                 attrs.ListID,
				Network:                clientlists.ActivationNetwork(attrs.Network),
				NotificationRecipients: attrs.NotificationRecipients,
				SiebelTicketID:         attrs.SiebelTicketID,
				Version:                int64(attrs.Version),
			}

			client.On("GetActivationStatus", testutils.MockContext, req).Return(&res, nil).Times(times)

			return &res
		}

		expectGetClientlist = func(client *clientlists.Mock, listID string, version int64, callTimes int) {
			clientListGetReq := clientlists.GetClientListRequest{
				ListID:       listID,
				IncludeItems: false,
			}

			clientList := clientlists.GetClientListResponse{ListContent: clientlists.ListContent{Version: version}}
			client.On("GetClientList", testutils.MockContext, clientListGetReq).Return(&clientList, nil).Times(callTimes)
		}

		expectAPIErrorWithCreateActivation = func(client *clientlists.Mock, req clientlists.CreateActivationRequest) {
			err := errors.New(createActivationAPIError)
			client.On("CreateActivation", testutils.MockContext, req).Return(nil, err).Once()
		}

		checkAttributes = func(a ActivationAttrs) resource.TestCheckFunc {
			resourceName := "akamai_clientlist_activation.activation_ASN_LIST_1"

			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "list_id", a.ListID),
				resource.TestCheckResourceAttr(resourceName, "comments", a.Comments),
				resource.TestCheckResourceAttr(resourceName, "notification_recipients.#", strconv.Itoa(len(a.NotificationRecipients))),
				resource.TestCheckResourceAttr(resourceName, "siebel_ticket_id", a.SiebelTicketID),
				resource.TestCheckResourceAttr(resourceName, "network", a.Network),
				resource.TestCheckResourceAttr(resourceName, "version", strconv.Itoa(a.Version)),
			}
			return resource.ComposeAggregateTestCheckFunc(checks...)
		}

		getActivationAttrs = func(actRes *clientlists.CreateActivationResponse, status clientlists.ActivationStatus) ActivationAttrs {
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

		activationReq = clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		}
	)

	t.Run("create activation", func(t *testing.T) {
		client := new(clientlists.Mock)

		activationRes := expectCreateActivation(client, activationReq, 2, 33)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.PendingActivation), 2)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.Active), 3)

		expectGetClientlist(client, "12_AB", 2, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 activationRes.ListID,
							Network:                string(activationRes.Network),
							NotificationRecipients: activationRes.NotificationRecipients,
							SiebelTicketID:         activationRes.SiebelTicketID,
							Comments:               activationRes.Comments,
							Version:                int(activationRes.Version),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update activation - version and other fields update", func(t *testing.T) {
		client := new(clientlists.Mock)

		activationRes := expectCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		}, 2, 33)
		updatedActivationRes := expectCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "UPDATED-12345",
				NotificationRecipients: []string{"update_user@example.com"},
				Comments:               "Activation Comments Updated",
			},
		}, 3, 34)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.PendingActivation), 1)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.Active), 4)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
			getActivationAttrs(updatedActivationRes, clientlists.PendingActivation), 2)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
			getActivationAttrs(updatedActivationRes, clientlists.Active), 3)

		expectGetClientlist(client, "12_AB", 2, 3)

		expectGetClientlist(client, "12_AB", 3, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 activationRes.ListID,
							Network:                string(activationRes.Network),
							NotificationRecipients: activationRes.NotificationRecipients,
							SiebelTicketID:         activationRes.SiebelTicketID,
							Comments:               activationRes.Comments,
							Version:                int(activationRes.Version),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_update.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 updatedActivationRes.ListID,
							Network:                string(updatedActivationRes.Network),
							NotificationRecipients: updatedActivationRes.NotificationRecipients,
							SiebelTicketID:         updatedActivationRes.SiebelTicketID,
							Comments:               updatedActivationRes.Comments,
							Version:                int(updatedActivationRes.Version),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update activation - version only update", func(t *testing.T) {
		client := new(clientlists.Mock)

		activationRes := expectCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		}, 2, 33)
		updatedActivationRes := expectCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		}, 3, 34)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.PendingActivation), 1)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.Active), 4)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
			getActivationAttrs(updatedActivationRes, clientlists.PendingActivation), 2)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: updatedActivationRes.ActivationID},
			getActivationAttrs(updatedActivationRes, clientlists.Active), 3)

		expectGetClientlist(client, "12_AB", 2, 3)

		expectGetClientlist(client, "12_AB", 3, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 activationRes.ListID,
							Network:                string(activationRes.Network),
							NotificationRecipients: activationRes.NotificationRecipients,
							SiebelTicketID:         activationRes.SiebelTicketID,
							Comments:               activationRes.Comments,
							Version:                int(activationRes.Version),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_update_version_only.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 updatedActivationRes.ListID,
							Network:                string(updatedActivationRes.Network),
							NotificationRecipients: updatedActivationRes.NotificationRecipients,
							SiebelTicketID:         updatedActivationRes.SiebelTicketID,
							Comments:               updatedActivationRes.Comments,
							Version:                int(updatedActivationRes.Version),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update activation - notification_recipients, siebel_ticket_id and comments updates suppressed", func(t *testing.T) {
		client := new(clientlists.Mock)

		activationRes := expectCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		}, 2, 33)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.PendingActivation), 2)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.Active), 5)

		expectGetClientlist(client, "12_AB", 2, 4)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 activationRes.ListID,
							Network:                string(activationRes.Network),
							NotificationRecipients: activationRes.NotificationRecipients,
							SiebelTicketID:         activationRes.SiebelTicketID,
							Comments:               activationRes.Comments,
							Version:                int(activationRes.Version),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_update_suppressed.tf", testDir)),
						Check: checkAttributes(ActivationAttrs{
							ListID:                 activationRes.ListID,
							Network:                string(activationRes.Network),
							NotificationRecipients: activationRes.NotificationRecipients,
							SiebelTicketID:         activationRes.SiebelTicketID,
							Comments:               activationRes.Comments,
							Version:                int(activationRes.Version),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Create activation with missing list id fails", func(t *testing.T) {
		client := new(clientlists.Mock)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/activation_missing_param.tf", testDir)),
						ExpectError: regexp.MustCompile("The argument \"list_id\" is required, but no definition was found"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Create activation api fails", func(t *testing.T) {
		client := new(clientlists.Mock)

		expectAPIErrorWithCreateActivation(client, clientlists.CreateActivationRequest{
			ListID: "12_AB",
			ActivationParams: clientlists.ActivationParams{
				Action:                 clientlists.Activate,
				Network:                clientlists.Staging,
				SiebelTicketID:         "ABC-12345",
				NotificationRecipients: []string{"user@example.com"},
				Comments:               "Activation Comments",
			},
		})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
						ExpectError: regexp.MustCompile(createActivationAPIError),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Import activation resource", func(t *testing.T) {
		client := new(clientlists.Mock)

		activationRes := expectCreateActivation(client, activationReq, 2, 33)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.PendingActivation), 3)

		expectReadActivation(client,
			clientlists.GetActivationRequest{ActivationID: activationRes.ActivationID},
			getActivationAttrs(activationRes, clientlists.Active), 4)

		expectGetClientlist(client, "12_AB", 2, 3)

		expectGetActivationStatus(client, clientlists.GetActivationStatusRequest{
			Network: clientlists.Staging,
			ListID:  activationReq.ListID,
		}, getActivationAttrs(activationRes, clientlists.Active), 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/activation_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "12_AB:STAGING",
						ResourceName:      "akamai_clientlist_activation.activation_ASN_LIST_1",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

package dns

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResDNSZone(t *testing.T) {
	zone := &dns.GetZoneResponse{
		ContractID:      "ctr1",
		Zone:            "primaryexampleterraform.io",
		Type:            "primary",
		Comment:         "This is a test primary zone",
		SignAndServe:    false,
		ActivationState: "PENDING",
	}

	upperCaseZone := &dns.GetZoneResponse{
		ContractID:      "ctr1",
		Zone:            "PRIMARYEXAMPLETERRAFORM.io",
		Type:            "primary",
		Comment:         "This is a test primary zone",
		SignAndServe:    false,
		ActivationState: "PENDING",
	}
	secondaryZone := &dns.GetZoneResponse{
		ContractID:      "ctr1",
		Zone:            "secondaryexampleterraform.io",
		Type:            "secondary",
		Comment:         "This is a test secondary zone",
		SignAndServe:    false,
		ActivationState: "PENDING",
		Masters:         []string{"1.1.1.1"},
		OutboundZoneTransfer: &dns.OutboundZoneTransfer{
			ACL:           []string{"192.0.2.156/24"},
			Enabled:       true,
			NotifyTargets: []string{"192.0.2.192"},
			TSIGKey: &dns.TSIGKey{
				Name:      "other.com.akamai.com",
				Algorithm: "hmac-sha1",
				Secret:    "fakeSecretajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
			},
		},
		TSIGKey: &dns.TSIGKey{
			Name:      "other.com.akamai.com",
			Algorithm: "hmac-sha512",
			Secret:    "fakeSecretjVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
		},
	}
	recordSetsResp := &dns.GetRecordSetsResponse{
		RecordSets: make([]dns.RecordSet, 2),
	}

	deleteBulkResp := &dns.DeleteBulkZonesResponse{
		RequestID:      "1234567890",
		ExpirationDate: "2023-10-01T00:00:00Z",
	}

	getDeleteStatusResp := dns.GetBulkZoneDeleteStatusResponse{
		RequestID:      "1234567890",
		ZonesSubmitted: 1,
		SuccessCount:   1,
		FailureCount:   0,
		IsComplete:     true,
		ExpirationDate: "2023-10-01T00:00:00Z",
	}

	getDeleteResultResp := dns.GetBulkZoneDeleteResultResponse{
		RequestID:                "1234567890",
		SuccessfullyDeletedZones: []string{"primaryexampleterraform.io"},
		FailedZones:              []dns.BulkFailedZone{},
	}

	getDeleteResultRespSecondaryGroup := dns.GetBulkZoneDeleteResultResponse{
		RequestID:                "1234567890",
		SuccessfullyDeletedZones: []string{"secondaryexampleterraform.io"},
		FailedZones:              []dns.BulkFailedZone{},
	}

	t.Run("when group is not provided and there is no group for the user ", func(t *testing.T) {
		client := &dns.Mock{}

		client.On("ListGroups",
			testutils.MockContext,
			mock.AnythingOfType("dns.ListGroupRequest"),
		).Return(&dns.ListGroupResponse{}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_without_group.tf"),
						ExpectError: regexp.MustCompile("no group found. Please provide the group."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	// This test performs a full life-cycle (CRUD) test
	t.Run("lifecycle test when group is not found and no. of group is 1", func(t *testing.T) {
		client := &dns.Mock{}
		groupListResponse := &dns.ListGroupResponse{
			Groups: []dns.Group{
				{
					GroupID:   1,
					GroupName: "name",
					ContractIDs: []string{
						"1", "2",
					},
					Permissions: []string{
						"DELETE", "READ", "WRITE", "ADD",
					},
				},
			},
		}

		client.On("ListGroups",
			testutils.MockContext,
			mock.AnythingOfType("dns.ListGroupRequest"),
		).Return(groupListResponse, nil)

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.test_without_group"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil)

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil)

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_without_group.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("when group is not provided and no. of group is more than 1 for the user ", func(t *testing.T) {
		client := &dns.Mock{}
		groupListResponse := &dns.ListGroupResponse{
			Groups: []dns.Group{
				{
					GroupID:   1,
					GroupName: "name",
					ContractIDs: []string{
						"1", "2",
					},
					Permissions: []string{
						"DELETE", "READ", "WRITE", "ADD",
					},
				},
				{
					GroupID:   2,
					GroupName: "name",
					ContractIDs: []string{
						"1", "2",
					},
					Permissions: []string{
						"DELETE", "READ", "WRITE", "ADD",
					},
				},
			},
		}

		client.On("ListGroups",
			testutils.MockContext,
			mock.AnythingOfType("dns.ListGroupRequest"),
		).Return(groupListResponse, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_without_group.tf"),
						ExpectError: regexp.MustCompile("group is a required field when there is more than one group present."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("error - failed to submit bulk delete zone", func(t *testing.T) {

		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(nil, fmt.Errorf("failed to submit bulk deletion")).Once()

		// Resource cleanup
		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil).Once()

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil).Once()

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(dataSourceName, "comment", "This is a test primary zone"),
							resource.TestCheckResourceAttr(dataSourceName, "group", "grp1"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Destroy:     true,
						ExpectError: regexp.MustCompile("failed to submit bulk deletion"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("error - delete status failed", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil).Twice()

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(nil, fmt.Errorf("could not get bulk zone delete status")).Once()

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil).Once()

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(dataSourceName, "comment", "This is a test primary zone"),
							resource.TestCheckResourceAttr(dataSourceName, "group", "grp1"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Destroy:     true,
						ExpectError: regexp.MustCompile("could not get bulk zone delete status"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("error - delete result failed", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil).Twice()

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil).Twice()

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(nil, fmt.Errorf("failed to delete zone")).Once()

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil).Once()
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(dataSourceName, "comment", "This is a test primary zone"),
							resource.TestCheckResourceAttr(dataSourceName, "group", "grp1"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Destroy:     true,
						ExpectError: regexp.MustCompile("failed to delete zone"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("lifecycle test with delete upper case zone", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{upperCaseZone, nil}
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"PRIMARYEXAMPLETERRAFORM.io"},
			},
			},
		).Return(deleteBulkResp, nil)

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil)

		// API normalizes the zone name to lower case, that's why we are using the lower case zone name in response
		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_upper_case_zone.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "PRIMARYEXAMPLETERRAFORM.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(dataSourceName, "comment", "This is a test primary zone"),
							resource.TestCheckResourceAttr(dataSourceName, "group", "grp1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("lifecycle test with group", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("UpdateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.UpdateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			zone.Comment = "This is an updated test primary zone"
		})

		client.On("SaveChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SaveChangeListRequest"),
		).Return(nil)

		client.On("SubmitChangeList",
			testutils.MockContext,
			mock.AnythingOfType("dns.SubmitChangeListRequest"),
		).Return(nil)

		client.On("GetRecordSets",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordSetsRequest"),
		).Return(recordSetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"primaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil)

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil)

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultResp, nil)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(dataSourceName, "comment", "This is a test primary zone"),
							resource.TestCheckResourceAttr(dataSourceName, "group", "grp1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/update_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	// This test performs a full life-cycle (CRUD) test
	t.Run("lifecycle test with group and secondary type", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		client.On("CreateZone",
			testutils.MockContext,
			dns.CreateZoneRequest{
				CreateZone: &dns.ZoneCreate{
					Zone:         "secondaryexampleterraform.io",
					Type:         "secondary",
					Comment:      "This is a test secondary zone",
					SignAndServe: false,
					Masters:      []string{"1.1.1.1"},
					OutboundZoneTransfer: &dns.OutboundZoneTransfer{
						ACL:           []string{"192.0.2.156/24"},
						Enabled:       true,
						NotifyTargets: []string{"192.0.2.192"},
						TSIGKey: &dns.TSIGKey{
							Name:      "other.com.akamai.com",
							Algorithm: "hmac-sha1",
							Secret:    "fakeSecretajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
						},
					},
					TSIGKey: &dns.TSIGKey{
						Name:      "other.com.akamai.com",
						Algorithm: "hmac-sha512",
						Secret:    "fakeSecretjVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLw==",
					},
				},
				ZoneQueryString: dns.ZoneQueryString{
					Contract: "ctr1",
					Group:    "grp1",
				},
				ClearConn: []bool{true},
			},
		).Return(nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{secondaryZone, nil}
		})

		client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(secondaryZone, nil).Times(4)

		client.On("UpdateZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.UpdateZoneRequest"),
		).Return(nil).Run(func(_ mock.Arguments) {
			secondaryZone.Comment = "This is an updated test secondary zone"
		})

		client.On("GetZone",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetZoneRequest"),
		).Return(secondaryZone, nil).Times(3)

		resourceName := "akamai_dns_zone.secondary_test_zone"

		client.On("DeleteBulkZones",
			testutils.MockContext,
			dns.DeleteBulkZonesRequest{ZonesList: &dns.ZoneNameListResponse{
				Zones: []string{"secondaryexampleterraform.io"},
			},
			},
		).Return(deleteBulkResp, nil)

		client.On("GetBulkZoneDeleteStatus",
			testutils.MockContext,
			dns.GetBulkZoneDeleteStatusRequest{RequestID: "1234567890"},
		).Return(&getDeleteStatusResp, nil)

		client.On("GetBulkZoneDeleteResult",
			testutils.MockContext,
			dns.GetBulkZoneDeleteResultRequest{RequestID: "1234567890"},
		).Return(&getDeleteResultRespSecondaryGroup, nil)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/create_secondary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "zone", "secondaryexampleterraform.io"),
							resource.TestCheckResourceAttr(resourceName, "contract", "ctr1"),
							resource.TestCheckResourceAttr(resourceName, "comment", "This is a test secondary zone"),
							resource.TestCheckResourceAttr(resourceName, "group", "grp1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsZone/update_secondary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "zone", "secondaryexampleterraform.io"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

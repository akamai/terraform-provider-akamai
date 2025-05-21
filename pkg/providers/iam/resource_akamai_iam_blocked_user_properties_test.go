package iam

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceIAMBlockedUserProperties(t *testing.T) {
	identityID := "test_identity_id"
	groupID := int64(12345)
	groupIDNew := int64(23456)

	propertiesCreate := []int64{1, 2, 3}
	propertiesUpdate := []int64{1, 2, 3, 4, 5}
	listRequest := iam.ListBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupID,
	}
	listRequestNew := iam.ListBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupIDNew,
	}
	createRequest := iam.UpdateBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupID,
		Properties: propertiesCreate,
	}
	createRequestNew := iam.UpdateBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupIDNew,
		Properties: propertiesCreate,
	}
	updateRequest := iam.UpdateBlockedPropertiesRequest{
		IdentityID: identityID,
		GroupID:    groupID,
		Properties: propertiesUpdate,
	}
	checkAttributes := func(properties []int64) resource.TestCheckFunc {

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "id", "test_identity_id:12345"),
			resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "identity_id", identityID),
			resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "group_id", strconv.FormatInt(groupID, 10)),
			resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "blocked_properties.#", strconv.Itoa(len(properties))),
		)
	}

	tests := map[string]struct {
		init  func(*iam.Mock)
		steps []resource.TestStep
	}{
		"basic": {
			init: func(m *iam.Mock) {
				// create
				expectListBlockedProperties(m, listRequest, []int64{}, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
				expectUpdateBlockedProperties(m, createRequest, propertiesCreate, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()

				// update
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
				expectUpdateBlockedProperties(m, updateRequest, propertiesUpdate, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesUpdate, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesUpdate, nil).Once()

				// delete
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/create.tf"),
					Check:  checkAttributes(propertiesCreate),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/update.tf"),
					Check:  checkAttributes(propertiesUpdate),
				},
			},
		},
		"update group id - new resource": {
			init: func(m *iam.Mock) {
				// create
				expectListBlockedProperties(m, listRequest, []int64{}, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
				expectUpdateBlockedProperties(m, createRequest, propertiesCreate, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()

				// read
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
				// create due to update on group_id
				expectListBlockedProperties(m, listRequestNew, []int64{}, nil).Once()
				expectUpdateBlockedProperties(m, createRequestNew, propertiesCreate, nil).Once()
				// read
				expectListBlockedProperties(m, listRequestNew, propertiesCreate, nil).Once()
				expectListBlockedProperties(m, listRequestNew, propertiesCreate, nil).Once()

				// delete
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/create.tf"),
					Check:  checkAttributes(propertiesCreate),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/update-group-id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "id", "test_identity_id:23456"),
						resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "identity_id", identityID),
						resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "group_id", strconv.FormatInt(groupIDNew, 10)),
						resource.TestCheckResourceAttr("akamai_iam_blocked_user_properties.test", "blocked_properties.#", strconv.Itoa(len(propertiesCreate))),
					),
				},
			},
		},
		"resource is already on server": {
			init: func(m *iam.Mock) {
				// create
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/create.tf"),
					ExpectError: regexp.MustCompile("there are already blocked properties on server, please import resource first"),
				},
			},
		},
		"empty properties": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/create-empty-properties.tf"),
					ExpectError: regexp.MustCompile("Not enough list items"),
				},
			},
		},
		"basic import": {
			init: func(m *iam.Mock) {
				// create
				expectListBlockedProperties(m, listRequest, []int64{}, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()
				expectUpdateBlockedProperties(m, createRequest, propertiesCreate, nil).Once()
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()

				// import
				expectListBlockedProperties(m, listRequest, propertiesCreate, nil).Once()

				// delete
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceIAMBlockedUserProperties/create.tf"),
					Check:  checkAttributes(propertiesCreate),
				},
				{
					ImportState:       true,
					ImportStateId:     "test_identity_id:12345",
					ResourceName:      "akamai_iam_blocked_user_properties.test",
					ImportStateVerify: true,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectListBlockedProperties(m *iam.Mock, request iam.ListBlockedPropertiesRequest, response []int64, err error) *mock.Call {
	on := m.On("ListBlockedProperties", testutils.MockContext, request)
	if err != nil {
		return on.Return(nil, err).Once()
	}
	return on.Return(response, nil)
}

func expectUpdateBlockedProperties(m *iam.Mock, request iam.UpdateBlockedPropertiesRequest, response []int64, err error) *mock.Call {
	on := m.On("UpdateBlockedProperties", testutils.MockContext, request)
	if err != nil {
		return on.Return(nil, err).Once()
	}
	return on.Return(response, nil)
}

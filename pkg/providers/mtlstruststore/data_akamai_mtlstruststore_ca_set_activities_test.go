package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v10/internal/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetActivitiesDataSource(t *testing.T) {
	t.Parallel()
	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_activities.test").
		CheckEqual("id", "12345").
		CheckEqual("name", "example-ca-set").
		CheckEqual("status", "NOT_DELETED").
		CheckEqual("created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("created_by", "example user").
		CheckEqual("activities.#", "2").
		CheckEqual("activities.0.type", "ACTIVATE_CA_SET_VERSION").
		CheckEqual("activities.0.network", "PRODUCTION").
		CheckEqual("activities.0.version", "1").
		CheckEqual("activities.0.activity_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("activities.0.activity_by", "example user").
		CheckEqual("activities.1.type", "CREATE_CA_SET").
		CheckEqual("activities.1.activity_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("activities.1.activity_by", "example user").
		CheckMissing("activities.1.network").
		CheckMissing("activities.1.version").
		CheckMissing("deleted_date").
		CheckMissing("deleted_by")

	tests := map[string]struct {
		init  func(*mtlstruststore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - without deletion details": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetActivities(t, m, "", "", false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/id.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - with deletion details": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetActivities(t, m, "", "", true)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/id.tf"),
					Check: commonStateChecker.
						CheckEqual("deleted_date", "2026-04-16T12:08:34.099457Z").
						CheckEqual("deleted_by", "example user").
						CheckEqual("status", "DELETED").
						Build(),
				},
			},
		},
		"happy path - find by ca set name and filter dates": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "12345",
							CASetName:   "test_name",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil).Times(3)
				mockListCASetActivities(t, m, "2024-04-16T12:08:34.099457Z", "2025-04-16T12:08:34.099457Z", false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/name_with_start_end.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - find by ca set name for non-unique prefix": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "01234",
							CASetName:   "test_name_foo",
							CASetStatus: "NOT_DELETED",
						},
						{
							CASetID:     "12345",
							CASetName:   "test_name",
							CASetStatus: "NOT_DELETED",
						},
						{
							CASetID:     "67890",
							CASetName:   "test_name_bar",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil).Times(3)
				mockListCASetActivities(t, m, "", "", false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/name.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"error: could not find by ca set name": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "01234",
							CASetName:   "test_name_foo",
							CASetStatus: "NOT_DELETED",
						},
						{
							CASetID:     "12345",
							CASetName:   "test_name_bar",
							CASetStatus: "DELETED",
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/name.tf"),
					ExpectError: regexp.MustCompile("no CA set found with the name 'test_name'"),
				},
			},
		},
		"error: failed to list CA sets": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(nil, fmt.Errorf("listing error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/name.tf"),
					ExpectError: regexp.MustCompile("(?s)could not find CA set with the name 'test_name' and status 'NOT_DELETED'. API.+error: listing error"),
				},
			},
		},
		"error: empty CA set name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/empty_name.tf"),
					ExpectError: regexp.MustCompile("Attribute name string length must be between 3 and 64, got: 0"),
				},
			},
		},
		"error API response": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASetActivities", testutils.MockContext, mtlstruststore.ListCASetActivitiesRequest{
					CASetID: "12345",
				}).Return(nil, fmt.Errorf("failed to retrieve CA set activities")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/id.tf"),
					ExpectError: regexp.MustCompile("failed to retrieve CA set activities"),
				},
			},
		},
		"validation error - missing required argument id or name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/no_fields.tf"),
					ExpectError: regexp.MustCompile(`(?s)No attribute specified when one \(and only one\) of \[id,name] is.+required`),
				},
			},
		},
		"validation error - both id and name are provided": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivities/id_name.tf"),
					ExpectError: regexp.MustCompile(`(?s)2 attributes specified when one \(and only one\) of \[name,id] is.+required`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
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

func mockListCASetActivities(t *testing.T, m *mtlstruststore.Mock, startDate, endDate string, withDeletion bool) {
	getResponse := &mtlstruststore.ListCASetActivitiesResponse{
		CASetID:     "12345",
		CASetName:   "example-ca-set",
		CASetStatus: "NOT_DELETED",
		CreatedDate: tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
		CreatedBy:   "example user",
		Activities: []mtlstruststore.CASetActivity{
			{
				Type:         "ACTIVATE_CA_SET_VERSION",
				Network:      ptr.To("PRODUCTION"),
				Version:      ptr.To(int64(1)),
				ActivityDate: tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
				ActivityBy:   "example user",
			},
			{
				Type:         "CREATE_CA_SET",
				ActivityDate: tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
				ActivityBy:   "example user",
			},
		},
	}

	if withDeletion {
		getResponse.DeletedBy = ptr.To("example user")
		getResponse.DeletedDate = ptr.To(tst.NewTimeFromString(t, "2026-04-16T12:08:34.099457Z"))
		getResponse.CASetStatus = "DELETED"
	}
	var start, end time.Time
	if startDate != "" {
		var err error
		start, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			t.Fatalf("failed to parse start date: %v", err)
		}
	}
	if endDate != "" {
		var err error
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			t.Fatalf("failed to parse end date: %v", err)
		}
	}
	m.On("ListCASetActivities", testutils.MockContext, mtlstruststore.ListCASetActivitiesRequest{
		CASetID: "12345",
		Start:   start,
		End:     end,
	}).Return(getResponse, nil).Times(3)
}

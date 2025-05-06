package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v7/internal/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetActivitiesDataSource(t *testing.T) {
	t.Parallel()
	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_activities.test").
		CheckEqual("ca_set_id", "12345").
		CheckEqual("ca_set_name", "example-ca-set").
		CheckEqual("status", "NOT_DELETED").
		CheckEqual("created_date", "2025-04-16 12:08:34.099457 +0000 UTC").
		CheckEqual("created_by", "example user").
		CheckEqual("activities.#", "2").
		CheckEqual("activities.0.type", "ACTIVATE_CA_SET_VERSION").
		CheckEqual("activities.0.network", "PRODUCTION").
		CheckEqual("activities.0.version", "1").
		CheckEqual("activities.0.activity_date", "2025-04-16 12:08:34.099457 +0000 UTC").
		CheckEqual("activities.0.activity_by", "example user").
		CheckEqual("activities.1.type", "CREATE_CA_SET").
		CheckEqual("activities.1.activity_date", "2025-04-16 12:08:34.099457 +0000 UTC").
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
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  ca_set_id = 12345
}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"happy path - with deletion details": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetActivities(t, m, "", "", true)
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  ca_set_id = 12345
}`,
					Check: commonStateChecker.
						CheckEqual("deleted_date", "2026-04-16 12:08:34.099457 +0000 UTC").
						CheckEqual("deleted_by", "example user").
						CheckEqual("status", "DELETED").
						Build(),
				},
			},
		},
		"happy path - find by ca set name and filter dates": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASets(m)
				mockListCASetActivities(t, m, "2024-04-16T12:08:34.099457Z", "2025-04-16T12:08:34.099457Z", false)
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  ca_set_name = "test name"
  start	   = "2024-04-16T12:08:34.099457Z"
  end 	   = "2025-04-16T12:08:34.099457Z"
}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"error API response": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASetActivities", testutils.MockContext, mtlstruststore.ListCASetActivitiesRequest{
					CASetID: 12345,
				}).Return(nil, fmt.Errorf("failed to retrieve CA set activities")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  ca_set_id = 12345
}`,
					ExpectError: regexp.MustCompile("failed to retrieve CA set activities"),
				},
			},
		},
		"validation error - missing required argument ca_set_id or ca_set_name": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {}
					`,
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[ca_set_id,ca_set_name\] is\s+required`),
				},
			},
		},
		"validation error - both ca_set_id and ca_set_name are provided": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activities" "test" {
  ca_set_id   = 12345
  ca_set_name = "example-ca-set"
}
					`,
					ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of \[ca_set_name,ca_set_id\] is\s+required`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
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

func mockListCASetActivities(t *testing.T, m *mtlstruststore.Mock, startDate, endDate string, withDeletion bool) {
	getResponse := &mtlstruststore.ListCASetActivitiesResponse{
		CASetID:     12345,
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
		CASetID: 12345,
		Start:   start,
		End:     end,
	}).Return(getResponse, nil).Times(3)
}

func mockListCASets(m *mtlstruststore.Mock) {
	m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: "test name",
	}).Return(&mtlstruststore.ListCASetsResponse{
		CASets: []mtlstruststore.CASetResponse{
			{
				CASetID: 12345,
			},
		},
	}, nil).Times(3)
}

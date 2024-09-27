package iam

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	testDataForAllowedAPIs struct {
		username           string
		clientType         *string
		allowAccountSwitch *bool
		allowedAPIs        []apiData
	}

	apiData struct {
		accessLevels      []iam.AccessLevel
		apiID             int64
		apiName           string
		description       string
		documentationURL  string
		endpoint          string
		hasAccess         bool
		serviceProviderID int64
	}
)

var (
	basicTestDataForAllowedAPIs = testDataForAllowedAPIs{
		username:           "test",
		clientType:         ptr.To("CLIENT"),
		allowAccountSwitch: ptr.To(true),
		allowedAPIs: []apiData{{
			accessLevels: []iam.AccessLevel{
				iam.ReadWriteLevel},
			apiID:             146,
			apiName:           "Property Manager (PAPI)",
			description:       "Property Manager (PAPI). PAPI requires access to Edge Hostnames. Please edit your authorizations to add HAPI to your API Client.",
			documentationURL:  "https://developer.akamai.com/api/luna/papi/overview.html",
			endpoint:          "/papi",
			hasAccess:         true,
			serviceProviderID: 1,
		}, {
			accessLevels: []iam.AccessLevel{
				iam.ReadOnlyLevel,
				iam.ReadWriteLevel},
			apiID:             11,
			apiName:           "Event Center",
			description:       "Event Center",
			documentationURL:  "https://developer.akamai.com/api/luna/events/overview.html",
			endpoint:          "/events",
			hasAccess:         true,
			serviceProviderID: 1,
		}},
	}
	basicTestDataForAllowedAPIsNoOptional = testDataForAllowedAPIs{
		username: "test",
		allowedAPIs: []apiData{{
			accessLevels: []iam.AccessLevel{
				iam.ReadWriteLevel},
			apiID:             146,
			apiName:           "Property Manager (PAPI)",
			description:       "Property Manager (PAPI). PAPI requires access to Edge Hostnames. Please edit your authorizations to add HAPI to your API Client.",
			documentationURL:  "https://developer.akamai.com/api/luna/papi/overview.html",
			endpoint:          "/papi",
			hasAccess:         true,
			serviceProviderID: 1,
		}, {
			accessLevels: []iam.AccessLevel{
				iam.ReadOnlyLevel,
				iam.ReadWriteLevel},
			apiID:             11,
			apiName:           "Event Center",
			description:       "Event Center",
			documentationURL:  "https://developer.akamai.com/api/luna/events/overview.html",
			endpoint:          "/events",
			hasAccess:         true,
			serviceProviderID: 1,
		}},
	}
)

func TestDataAllowedAPIs(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *iam.Mock, testDataForAllowedAPIs)
		mockData   testDataForAllowedAPIs
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataAllowedAPIs/default.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForAllowedAPIs) {
				expectFullListAllowedAPIs(t, m, testData, 3)
			},
			mockData: basicTestDataForAllowedAPIs,
		},
		"happy path no optional values": {
			configPath: "testdata/TestDataAllowedAPIs/default_no_optional.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForAllowedAPIs) {
				expectFullListAllowedAPIs(t, m, testData, 3)
			},
			mockData: basicTestDataForAllowedAPIsNoOptional,
		},
		"error - ListAllowedAPIs call failed": {
			configPath: "testdata/TestDataAllowedAPIs/default_no_optional.tf",
			init: func(t *testing.T, m *iam.Mock, testData testDataForAllowedAPIs) {
				listAllowedAPIsReq := iam.ListAllowedAPIsRequest{UserName: testData.username}

				m.On("ListAllowedAPIs", mock.Anything, listAllowedAPIsReq).Return(nil, errors.New("test error"))
			},
			error:    regexp.MustCompile("test error"),
			mockData: basicTestDataForAllowedAPIsNoOptional,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &iam.Mock{}
			if tc.init != nil {
				tc.init(t, client, tc.mockData)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, tc.configPath),
							Check:       checkAllowedAPIsAttrs(tc.mockData),
							ExpectError: tc.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectFullListAllowedAPIs(_ *testing.T, client *iam.Mock, data testDataForAllowedAPIs, timesToRun int) {
	listAllowedAPIsReq := iam.ListAllowedAPIsRequest{
		UserName: data.username,
	}
	if data.clientType != nil {
		listAllowedAPIsReq.ClientType = iam.ClientType(*data.clientType)
	}
	if data.allowAccountSwitch != nil {
		listAllowedAPIsReq.AllowAccountSwitch = *data.allowAccountSwitch
	}

	listAllowedAPIsRes := iam.ListAllowedAPIsResponse{}

	for _, api := range data.allowedAPIs {

		listAllowedAPIsRes = append(listAllowedAPIsRes, iam.AllowedAPI{
			AccessLevels:      api.accessLevels,
			APIID:             api.apiID,
			APIName:           api.apiName,
			Description:       api.description,
			DocumentationURL:  api.documentationURL,
			Endpoint:          api.endpoint,
			HasAccess:         api.hasAccess,
			ServiceProviderID: api.serviceProviderID,
		})
	}

	client.On("ListAllowedAPIs", mock.Anything, listAllowedAPIsReq).Return(listAllowedAPIsRes, nil).Times(timesToRun)
}

func checkAllowedAPIsAttrs(data testDataForAllowedAPIs) resource.TestCheckFunc {
	name := "data.akamai_iam_allowed_apis.test"
	checksFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "username", data.username),
	}

	if data.clientType != nil {
		resource.TestCheckResourceAttr(name, "client_type", *data.clientType)
	} else {
		resource.TestCheckNoResourceAttr(name, "client_type")
	}

	if data.allowAccountSwitch != nil {
		resource.TestCheckResourceAttr(name, "allow_account_switch", strconv.FormatBool(*data.allowAccountSwitch))
	} else {
		resource.TestCheckNoResourceAttr(name, "allow_account_switch")
	}

	for i, api := range data.allowedAPIs {
		for j, accessLevel := range api.accessLevels {
			checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.access_levels.%d", i, j), string(accessLevel)))
		}
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.api_id", i), strconv.FormatInt(api.apiID, 10)))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.api_name", i), api.apiName))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.description", i), api.description))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.documentation_url", i), api.documentationURL))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.endpoint", i), api.endpoint))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.has_access", i), strconv.FormatBool(api.hasAccess)))
		checksFuncs = append(checksFuncs, resource.TestCheckResourceAttr(name, fmt.Sprintf("allowed_apis.%d.service_provider_id", i), strconv.FormatInt(api.serviceProviderID, 10)))

	}

	return resource.ComposeAggregateTestCheckFunc(checksFuncs...)
}

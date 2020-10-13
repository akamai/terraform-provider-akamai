package gtm

import (
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"net/http"
	"regexp"
	"testing"
)

var rsrc = gtm.Resource{
	Name:            "tfexample_resource_1",
	AggregationType: "latest",
	Type:            "XML load object via HTTP",
	ResourceInstances: []*gtm.ResourceInstance{
		{
			DatacenterId:         3131,
			UseDefaultLoadObject: false,
			LoadObject: gtm.LoadObject{
				LoadObject:     "/test1",
				LoadServers:    []string{"1.2.3.4"},
				LoadObjectPort: 80,
			},
		},
	},
}

func TestResGtmResource(t *testing.T) {

	t.Run("create resource", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetResource",
			mock.Anything, // ctx is irrelevant for this test
			rsrc.Name,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.ResourceResponse{}
		resp.Resource = &rsrc
		resp.Status = &pendingResponseStatus
		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("UpdateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("DeleteResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmResource/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource failed", func(t *testing.T) {
		client := &mockgtm{}

		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Resource Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource denied", func(t *testing.T) {
		client := &mockgtm{}

		dr := gtm.ResourceResponse{}
		dr.Resource = &rsrc
		dr.Status = &deniedResponseStatus
		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

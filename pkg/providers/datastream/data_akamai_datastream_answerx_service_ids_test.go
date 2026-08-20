package datastream

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataAnswerXServiceIDsRead(t *testing.T) {
	t.Parallel()

	configPath := "testdata/TestDataAnswerXServiceIDs/list_answerx_service_ids.tf"
	pagedReq := func(page int64) datastream.ListAnswerXServiceIDsRequest {
		return datastream.ListAnswerXServiceIDsRequest{
			ContractID: "test_contract",
			Page:       page,
			PageSize:   answerXServiceIDsPageSize,
		}
	}
	metadata := func(page, lastPage int64) *datastream.PaginationMetadata {
		return &datastream.PaginationMetadata{Page: page, LastPage: lastPage, PageSize: answerXServiceIDsPageSize}
	}

	t.Run("list service IDs successfully - single page", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:   metadata(1, 1),
				ContractID: "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
					{SSID: 202, Name: "ServiceB", Product: "AnswerX"},
				},
			}, nil).Times(3)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							test.NewStateChecker("data.akamai_datastream_answerx_service_ids.test").
								CheckEqual("contract_id", "test_contract").
								CheckEqual("service_ids.#", "2").
								Build(),
							resource.TestCheckTypeSetElemNestedAttrs("data.akamai_datastream_answerx_service_ids.test", "service_ids.*", map[string]string{
								"id":      "101",
								"name":    "ServiceA",
								"product": "AnswerX",
							}),
							resource.TestCheckTypeSetElemNestedAttrs("data.akamai_datastream_answerx_service_ids.test", "service_ids.*", map[string]string{
								"id":      "202",
								"name":    "ServiceB",
								"product": "AnswerX",
							}),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("list service IDs successfully - multiple pages", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:   metadata(1, 2),
				ContractID: "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
				},
			}, nil).Times(3)
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(2)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:   metadata(2, 2),
				ContractID: "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
					{SSID: 202, Name: "ServiceB", Product: "AnswerX"},
				},
			}, nil).Times(3)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							test.NewStateChecker("data.akamai_datastream_answerx_service_ids.test").
								CheckEqual("contract_id", "test_contract").
								CheckEqual("service_ids.#", "2").
								Build(),
							resource.TestCheckTypeSetElemNestedAttrs("data.akamai_datastream_answerx_service_ids.test", "service_ids.*", map[string]string{
								"id":      "101",
								"name":    "ServiceA",
								"product": "AnswerX",
							}),
							resource.TestCheckTypeSetElemNestedAttrs("data.akamai_datastream_answerx_service_ids.test", "service_ids.*", map[string]string{
								"id":      "202",
								"name":    "ServiceB",
								"product": "AnswerX",
							}),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("API error returns diagnostic", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			(*datastream.ListAnswerXServiceIDsResponse)(nil),
			&datastream.Error{StatusCode: 403, Title: "Forbidden"},
		).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, configPath),
						ExpectError: regexp.MustCompile(`Listing AnswerX service IDs failed[\s\S]*Forbidden`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("blank normalized contract ID returns diagnostic", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: `
							provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							data "akamai_datastream_answerx_service_ids" "test" {
								contract_id = "ctr_"
							}`,
						ExpectError: regexp.MustCompile(`Invalid contract ID[\s\S]*contract_id must not be blank`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("list service IDs successfully - empty list", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:          metadata(1, 1),
				ContractID:        "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{},
			}, nil).Times(3)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, configPath),
						Check: test.NewStateChecker("data.akamai_datastream_answerx_service_ids.test").
							CheckEqual("contract_id", "test_contract").
							CheckEqual("service_ids.#", "0").
							Build(),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid API response - missing metadata", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:          nil,
				ContractID:        "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{},
			}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, configPath),
						ExpectError: regexp.MustCompile("invalid API response"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid API response - pagination metadata mismatch", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name     string
			metadata *datastream.PaginationMetadata
		}{
			{name: "page mismatch", metadata: metadata(2, 2)},
			{name: "last page lower than requested page", metadata: metadata(1, 0)},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client := &datastream.Mock{}
				client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
					&datastream.ListAnswerXServiceIDsResponse{
						Metadata:          tc.metadata,
						ContractID:        "test_contract",
						AnswerXServiceIDs: []datastream.AnswerXServiceDetail{},
					}, nil).Once()

				useClient(client, func() {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						Steps: []resource.TestStep{
							{
								Config:      testutils.LoadFixtureString(t, configPath),
								ExpectError: regexp.MustCompile("invalid API response"),
							},
						},
					})
				})
				client.AssertExpectations(t)
			})
		}
	})

	t.Run("contract_id prefix handling - ctr_ prefix stripped on API call", func(t *testing.T) {
		t.Parallel()

		client := &datastream.Mock{}
		// Config uses ctr_ prefix, but API request must have it stripped to "test_contract"
		client.On("ListAnswerXServiceIDs", testutils.MockContext, pagedReq(1)).Return(
			&datastream.ListAnswerXServiceIDsResponse{
				Metadata:   metadata(1, 1),
				ContractID: "test_contract",
				AnswerXServiceIDs: []datastream.AnswerXServiceDetail{
					{SSID: 101, Name: "ServiceA", Product: "AnswerX"},
				},
			}, nil).Times(3)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: `
							provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
								
							data "akamai_datastream_answerx_service_ids" "test" {
								contract_id = "ctr_test_contract"
							}`,
						Check: resource.ComposeAggregateTestCheckFunc(
							test.NewStateChecker("data.akamai_datastream_answerx_service_ids.test").
								CheckEqual("contract_id", "ctr_test_contract").
								CheckEqual("service_ids.#", "1").
								Build(),
							resource.TestCheckTypeSetElemNestedAttrs("data.akamai_datastream_answerx_service_ids.test", "service_ids.*", map[string]string{
								"id":      "101",
								"name":    "ServiceA",
								"product": "AnswerX",
							}),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

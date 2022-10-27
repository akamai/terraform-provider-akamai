package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/datastream"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataAkamaiDatastreamActivationHistoryRead(t *testing.T) {
	tests := map[string]struct {
		configPath                 string
		getActivationHistoryReturn []datastream.ActivationHistoryEntry
		checkFuncs                 []resource.TestCheckFunc
		edgegridError              error
		withError                  *regexp.Regexp
	}{
		"validate activation history response": {
			configPath: "testdata/TestDataAkamaiDatastreamActivationHistoryRead/activation_history.tf",
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{
				{
					CreatedBy:       "user1",
					CreatedDate:     "16-01-2020 11:07:12 GMT",
					IsActive:        false,
					StreamID:        7050,
					StreamVersionID: 2,
				},
				{
					CreatedBy:       "user2",
					CreatedDate:     "16-01-2020 09:31:02 GMT",
					IsActive:        true,
					StreamID:        7050,
					StreamVersionID: 2,
				},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "stream_id", "7050"),

				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.#", "2"),

				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.created_by", "user1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.created_date", "16-01-2020 11:07:12 GMT"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_version_id", "2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.is_active", "false"),

				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.created_by", "user2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.created_date", "16-01-2020 09:31:02 GMT"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.stream_version_id", "2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.is_active", "true"),
			},
		},
		"validate empty response": {
			configPath:                 "testdata/TestDataAkamaiDatastreamActivationHistoryRead/empty_activation_history.tf",
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "stream_id", "7051"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.#", "0"),
			},
		},
		"edgegrid error": {
			configPath:    "testdata/TestDataAkamaiDatastreamActivationHistoryRead/empty_activation_history.tf",
			edgegridError: fmt.Errorf("%w: request failed: %s", datastream.ErrGetActivationHistory, errors.New("500")),
			withError:     regexp.MustCompile("view activation history: request failed: 500"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			client := datastream.Mock{}
			useClient(&client, func() {
				if test.edgegridError != nil {
					client.On("GetActivationHistory", mock.Anything, mock.Anything).Return(nil, test.edgegridError).Once()
				} else {
					client.On("GetActivationHistory", mock.Anything, mock.Anything).Return(test.getActivationHistoryReturn, nil)
				}
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(test.configPath),
							Check: resource.ComposeAggregateTestCheckFunc(
								test.checkFuncs...,
							),
							ExpectError: test.withError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

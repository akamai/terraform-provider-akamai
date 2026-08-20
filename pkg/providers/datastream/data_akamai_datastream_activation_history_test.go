package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataAkamaiDatastreamActivationHistoryRead(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath                 string
		expectedLogType            datastream.LogType
		getActivationHistoryReturn []datastream.ActivationHistoryEntry
		checkFuncs                 []resource.TestCheckFunc
		edgegridError              error
		withError                  *regexp.Regexp
	}{
		"validate activation history response - appsec log type": {
			configPath:      "testdata/TestDataAkamaiDatastreamActivationHistoryRead/activation_history_appsec.tf",
			expectedLogType: datastream.LogTypeAppSec,
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{
				{
					ModifiedBy:    "user1",
					ModifiedDate:  "16-01-2020 11:07:12 GMT",
					Status:        datastream.StreamStatusDeactivated,
					StreamID:      7050,
					StreamVersion: 2,
				},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.modified_by", "user1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.status", "DEACTIVATED"),
			},
		},
		"validate activation history response - answerx log type": {
			configPath:      "testdata/TestDataAkamaiDatastreamActivationHistoryRead/activation_history_answerx.tf",
			expectedLogType: datastream.LogTypeAnswerX,
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{
				{
					ModifiedBy:    "user2",
					ModifiedDate:  "16-01-2020 09:31:02 GMT",
					Status:        datastream.StreamStatusActivated,
					StreamID:      7050,
					StreamVersion: 2,
				},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.modified_by", "user2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.status", "ACTIVATED"),
			},
		},
		"validate activation history response": {
			configPath:      "testdata/TestDataAkamaiDatastreamActivationHistoryRead/activation_history.tf",
			expectedLogType: datastream.LogTypeCDN,
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{
				{
					ModifiedBy:    "user1",
					ModifiedDate:  "16-01-2020 11:07:12 GMT",
					Status:        datastream.StreamStatusDeactivated,
					StreamID:      7050,
					StreamVersion: 2,
				},
				{
					ModifiedBy:    "user2",
					ModifiedDate:  "16-01-2020 09:31:02 GMT",
					Status:        datastream.StreamStatusActivated,
					StreamID:      7050,
					StreamVersion: 2,
				},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.modified_by", "user1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.modified_date", "16-01-2020 11:07:12 GMT"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.stream_version", "2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.0.status", "DEACTIVATED"),

				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.modified_by", "user2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.modified_date", "16-01-2020 09:31:02 GMT"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.stream_id", "7050"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.stream_version", "2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "activations.1.status", "ACTIVATED"),
			},
		},
		"validate empty response": {
			configPath:                 "testdata/TestDataAkamaiDatastreamActivationHistoryRead/empty_activation_history.tf",
			expectedLogType:            datastream.LogTypeCDN,
			getActivationHistoryReturn: []datastream.ActivationHistoryEntry{},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_activation_history.test", "stream_id", "7051"),
			},
		},
		"edgegrid error": {
			configPath:      "testdata/TestDataAkamaiDatastreamActivationHistoryRead/empty_activation_history.tf",
			expectedLogType: datastream.LogTypeCDN,
			edgegridError:   fmt.Errorf("%w: request failed: %s", datastream.ErrGetActivationHistory, errors.New("500")),
			withError:       regexp.MustCompile("view activation history: request failed: 500"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			client := datastream.Mock{}
			useClient(&client, func() {
				reqMatcher := mock.MatchedBy(func(r datastream.GetActivationHistoryRequest) bool {
					return r.LogType == test.expectedLogType
				})
				if test.edgegridError != nil {
					client.On("GetActivationHistory", testutils.MockContext, reqMatcher).Return(nil, test.edgegridError).Once()
				} else {
					client.On("GetActivationHistory", testutils.MockContext, reqMatcher).Return(test.getActivationHistoryReturn, nil).Times(3)
				}
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.configPath),
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

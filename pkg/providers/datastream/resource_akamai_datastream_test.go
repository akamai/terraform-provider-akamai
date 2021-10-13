package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourceStream(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		client := &mockdatastream{}

		PollForActivationStatusChangeInterval = 1 * time.Millisecond

		streamID := int64(12321)

		streamConfiguration := datastream.StreamConfiguration{
			ActivateNow: true,
			Config: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
			Connectors: []datastream.AbstractConnector{
				&datastream.S3Connector{
					AccessKey:       "s3_test_access_key",
					Bucket:          "s3_test_bucket",
					ConnectorName:   "s3_test_connector_name",
					Path:            "s3_test_path",
					Region:          "s3_test_region",
					SecretAccessKey: "s3_test_secret_key",
				},
			},
			ContractID:      "test_contract",
			DatasetFieldIDs: []int{1000, 1001, 1002},
			EmailIDs:        "test_email1@akamai.com,test_email2@akamai.com",
			GroupID:         tools.IntPtr(1337),
			PropertyIDs:     []int{1, 2, 3},
			StreamName:      "test_stream",
			StreamType:      datastream.StreamTypeRawLogs,
			TemplateName:    datastream.TemplateNameEdgeLogs,
		}

		createStreamRequest := datastream.CreateStreamRequest{
			StreamConfiguration: streamConfiguration,
		}

		updateStreamResponse := &datastream.StreamUpdate{
			StreamVersionKey: datastream.StreamVersionKey{
				StreamID:        streamID,
				StreamVersionID: 1,
			},
		}

		updateStreamRequest := datastream.UpdateStreamRequest{
			StreamID: 12321,
			StreamConfiguration: datastream.StreamConfiguration{
				ActivateNow: false,
				Config: datastream.Config{
					Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
					Format:    datastream.FormatTypeStructured,
					Frequency: datastream.Frequency{
						TimeInSec: datastream.TimeInSec30,
					},
					UploadFilePrefix: "prefix_updated",
					UploadFileSuffix: "suf_updated",
				},
				Connectors: []datastream.AbstractConnector{
					&datastream.S3Connector{
						AccessKey:       "s3_test_access_key",
						Bucket:          "s3_test_bucket_updated",
						ConnectorName:   "s3_test_connector_name_updated",
						Path:            "s3_test_path",
						Region:          "s3_test_region",
						SecretAccessKey: "s3_test_secret_key",
					},
				},
				ContractID:      streamConfiguration.ContractID,
				DatasetFieldIDs: streamConfiguration.DatasetFieldIDs,
				EmailIDs:        "test_email1_updated@akamai.com,test_email2@akamai.com",
				PropertyIDs:     streamConfiguration.PropertyIDs,
				StreamName:      "test_stream_with_updated",
				StreamType:      streamConfiguration.StreamType,
				TemplateName:    streamConfiguration.TemplateName,
			},
		}

		modifyResponse := func(r datastream.DetailedStreamVersion, opt func(r *datastream.DetailedStreamVersion)) *datastream.DetailedStreamVersion {
			opt(&r)
			return &r
		}

		getStreamResponseActivated := &datastream.DetailedStreamVersion{
			ActivationStatus: datastream.ActivationStatusActivated,
			Config:           streamConfiguration.Config,
			Connectors: []datastream.ConnectorDetails{
				{
					Bucket:        "s3_test_bucket",
					ConnectorType: datastream.ConnectorTypeS3,
					ConnectorName: "s3_test_connector_name",
					Path:          "s3_test_path",
					Region:        "s3_test_region",
				},
			},
			ContractID:  streamConfiguration.ContractID,
			CreatedBy:   "johndoe",
			CreatedDate: "10-07-2020 12:19:02 GMT",
			Datasets: []datastream.DataSets{
				{
					DatasetGroupName:        "group_name_1",
					DatasetGroupDescription: "group_desc_1",
					DatasetFields: []datastream.DatasetFields{
						{
							DatasetFieldID:          1000,
							DatasetFieldName:        "dataset_field_name_1",
							DatasetFieldDescription: "dataset_field_desc_1",
							Order:                   0,
						},
						{
							DatasetFieldID:          1001,
							DatasetFieldName:        "dataset_field_name_2",
							DatasetFieldDescription: "dataset_field_desc_2",
							Order:                   1,
						},
						{
							DatasetFieldID:          1002,
							DatasetFieldName:        "dataset_field_name_3",
							DatasetFieldDescription: "dataset_field_desc_3",
							Order:                   2,
						},
					},
				},
			},
			EmailIDs:     streamConfiguration.EmailIDs,
			Errors:       nil,
			GroupID:      *streamConfiguration.GroupID,
			GroupName:    "Default Group-1-ABCDE",
			ModifiedBy:   "janesmith",
			ModifiedDate: "15-07-2020 05:51:52 GMT",
			ProductID:    "Download_Delivery",
			ProductName:  "Download Delivery",
			Properties: []datastream.Property{
				{
					PropertyID:   1,
					PropertyName: "property_1",
				},
				{
					PropertyID:   2,
					PropertyName: "property_2",
				},
				{
					PropertyID:   3,
					PropertyName: "property_3",
				},
			},
			StreamID:        updateStreamResponse.StreamVersionKey.StreamID,
			StreamName:      streamConfiguration.StreamName,
			StreamType:      streamConfiguration.StreamType,
			StreamVersionID: updateStreamResponse.StreamVersionKey.StreamVersionID,
			TemplateName:    streamConfiguration.TemplateName,
		}

		getStreamResponseStreamActivating := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusActivating
		})

		getStreamResponseStreamActivatingAfterUpdate := modifyResponse(*getStreamResponseActivated, func(r *datastream.DetailedStreamVersion) {
			r.Config = datastream.Config{
				Delimiter:        updateStreamRequest.StreamConfiguration.Config.Delimiter,
				Format:           updateStreamRequest.StreamConfiguration.Config.Format,
				Frequency:        updateStreamRequest.StreamConfiguration.Config.Frequency,
				UploadFilePrefix: updateStreamRequest.StreamConfiguration.Config.UploadFilePrefix,
				UploadFileSuffix: updateStreamRequest.StreamConfiguration.Config.UploadFileSuffix,
			}
			r.EmailIDs = updateStreamRequest.StreamConfiguration.EmailIDs
			r.StreamName = updateStreamRequest.StreamConfiguration.StreamName
			r.Connectors = []datastream.ConnectorDetails{
				{
					Bucket:        "s3_test_bucket_updated",
					ConnectorType: datastream.ConnectorTypeS3,
					ConnectorName: "s3_test_connector_name_updated",
					Path:          "s3_test_path",
					Region:        "s3_test_region",
				},
			}
		})

		getStreamResponseStreamActivatedAfterUpdate := modifyResponse(*getStreamResponseStreamActivatingAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusActivated
		})

		getStreamResponseDeactivating := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusDeactivating
		})

		getStreamResponseDeactivated := modifyResponse(*getStreamResponseStreamActivatedAfterUpdate, func(r *datastream.DetailedStreamVersion) {
			r.ActivationStatus = datastream.ActivationStatusDeactivated
		})

		getStreamRequest := datastream.GetStreamRequest{
			StreamID: streamID,
		}

		client.On("CreateStream", mock.Anything, createStreamRequest).
			Return(updateStreamResponse, nil)

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivating, nil).
			Times(3)

		// first for finishing waitForStreamStatusChange
		// second for complete CreateStream
		// third for reading resource state
		// fourth for reading stream status in UpdateStream
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseActivated, nil).
			Times(4)

		client.On("UpdateStream", mock.Anything, updateStreamRequest).
			Return(updateStreamResponse, nil)

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivatedAfterUpdate, nil).
			Times(3)

		// first for finishing waitForStreamStatusChange in UpdateStream
		// second for reading resource state after UpdateStream
		// third for reading stream status in DeleteStream
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseStreamActivatedAfterUpdate, nil).
			Times(3)

		client.On("DeleteStream", mock.Anything, datastream.DeleteStreamRequest{
			StreamID: streamID,
		}).Return(&datastream.DeleteStreamResponse{Message: "Success"}, nil)

		client.On("DeactivateStream", mock.Anything, datastream.DeactivateStreamRequest{
			StreamID: 12321,
		}).Return(&datastream.DeactivateStreamResponse{
			StreamVersionKey: datastream.StreamVersionKey{
				StreamID:        streamID,
				StreamVersionID: 1,
			},
		}, nil)

		// for waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseDeactivating, nil).
			Times(3)

		// for finishing waitForStreamStatusChange
		client.On("GetStream", mock.Anything, getStreamRequest).
			Return(getStreamResponseDeactivated, nil).
			Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResourceStream/lifecycle/create_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", "12321"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_prefix", "pre"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_suffix", "suf"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.0.time_in_sec", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "property_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_type", string(datastream.StreamTypeRawLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "template_name", string(datastream.TemplateNameEdgeLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.connector_name", "s3_test_connector_name"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.path", "s3_test_path"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.region", "s3_test_region"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.secret_access_key", "s3_test_secret_key"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResourceStream/lifecycle/update_stream.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_datastream.s", "id", "12321"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "active", "true"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.delimiter", string(datastream.DelimiterTypeSpace)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.format", string(datastream.FormatTypeStructured)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_prefix", "prefix_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.upload_file_suffix", "suf_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "config.0.frequency.0.time_in_sec", "30"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "contract_id", "test_contract"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "dataset_fields_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.#", "2"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.0", "test_email1_updated@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "email_ids.1", "test_email2@akamai.com"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "group_id", "1337"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "property_ids.#", "3"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_name", "test_stream_with_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "stream_type", string(datastream.StreamTypeRawLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "template_name", string(datastream.TemplateNameEdgeLogs)),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.#", "1"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.access_key", "s3_test_access_key"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.bucket", "s3_test_bucket_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.connector_name", "s3_test_connector_name_updated"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.path", "s3_test_path"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.region", "s3_test_region"),
							resource.TestCheckResourceAttr("akamai_datastream.s", "s3_connector.0.secret_access_key", "s3_test_secret_key"),
						),
					},
				},
			})

			client.AssertExpectations(t)
		})
	})
}

func TestResourceStreamErrors(t *testing.T) {
	tests := map[string]struct {
		tfFile    string
		init      func(*mockdatastream)
		withError *regexp.Regexp
	}{
		"missing required argument": {
			tfFile:    "testdata/TestResourceStream/errors/missing_required_argument/missing_required_argument.tf",
			withError: regexp.MustCompile("Missing required argument"),
		},
		"internal server error": {
			tfFile: "testdata/TestResourceStream/errors/internal_server_error/internal_server_error.tf",
			init: func(m *mockdatastream) {
				m.On("CreateStream", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("%w: request failed: %s", datastream.ErrCreateStream, errors.New("500")))
			},
			withError: regexp.MustCompile(datastream.ErrCreateStream.Error()),
		},
		"stream with this name already exists": {
			tfFile: "testdata/TestResourceStream/errors/stream_name_not_unique/stream_name_not_unique.tf",
			init: func(m *mockdatastream) {
				m.On("CreateStream", mock.Anything, mock.Anything).
					Return(nil, fmt.Errorf("%s: %w", datastream.ErrCreateStream, &datastream.Error{
						Type:       "bad-request",
						Title:      "Bad Request",
						StatusCode: 400,
						Errors: []datastream.RequestErrors{
							{
								Type:   "bad-request",
								Title:  "Bad Request",
								Detail: "Stream with name test_stream already exists.",
							},
						},
					}))
			},
			withError: regexp.MustCompile("Stream with name test_stream already exists"),
		},
		"invalid email format": {
			tfFile:    "testdata/TestResourceStream/errors/invalid_email/invalid_email.tf",
			withError: regexp.MustCompile("must be a valid email address"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockdatastream{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.tfFile),
							ExpectError: test.withError,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func TestResourceStreamCustomDiff(t *testing.T) {
	client := &mockdatastream{}

	tests := map[string]struct {
		tfFile    string
		withError *regexp.Regexp
	}{
		"prefix and suffix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff1.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"prefix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff2.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"suffix present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff3.tf",
			withError: regexp.MustCompile("cannot be used with"),
		},
		"prefix and suffix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff4.tf",
			withError: nil,
		},
		"prefix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff5.tf",
			withError: nil,
		},
		"suffix present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff6.tf",
			withError: nil,
		},
		"prefix and suffix not present with not allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff7.tf",
			withError: nil,
		},
		"prefix and suffix not present with allowed connector": {
			tfFile:    "testdata/TestResourceStream/custom_diff/custom_diff8.tf",
			withError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:             loadFixtureString(test.tfFile),
							ExpectError:        test.withError,
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
						},
					},
				})
			})
		})

		client.AssertExpectations(t)
	}
}

func TestPrefixSuppressor(t *testing.T) {
	tests := []struct {
		oldPropertyID, newPropertyID string
		expected                     bool
	}{
		{"", "", true},
		{"123", "123", true},
		{"123", "prp_123", true},
		{"prp_123", "123", true},
		{"prp_123", "prp_123", true},
		{"", "prp_123", false},
		{"prp_123", "1234", false},
		{"123", "1234", false},
		{"prp_123", "prp_1234", false},
	}
	suppressor := prefixSuppressor("prp_")
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s vs %s", test.oldPropertyID, test.newPropertyID), func(t *testing.T) {
			assert.Equal(t, test.expected, suppressor("", test.oldPropertyID, test.newPropertyID, nil))
		})
	}
}

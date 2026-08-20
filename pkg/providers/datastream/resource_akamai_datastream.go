package datastream

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	// PollForActivationStatusChangeInterval defines retry interval for getting status of a pending change
	PollForActivationStatusChangeInterval = 10 * time.Minute

	// ExactlyOneConnectorRule defines connector fields names
	ExactlyOneConnectorRule = []string{
		"azure_connector",
		"datadog_connector",
		"elasticsearch_connector",
		"gcs_connector",
		"https_connector",
		"loggly_connector",
		"new_relic_connector",
		"oracle_connector",
		"s3_connector",
		"splunk_connector",
		"sumologic_connector",
		"s3_compatible_connector",
		"trafficpeak_connector",
		"dynatrace_connector",
	}

	// ConnectorsWithoutFilenameOptionsConfig defines connectors without option to configure prefix and suffix
	ConnectorsWithoutFilenameOptionsConfig = []string{
		"datadog_connector",
		"elasticsearch_connector",
		"https_connector",
		"loggly_connector",
		"new_relic_connector",
		"splunk_connector",
		"sumologic_connector",
		"trafficpeak_connector",
		"dynatrace_connector",
	}

	// ConnectorsSupportOnlyJSONLogFormat defines connectors that support only JSON log format
	ConnectorsSupportOnlyJSONLogFormat = []string{
		"new_relic_connector",
		"elasticsearch_connector",
		"trafficpeak_connector",
		"dynatrace_connector",
	}

	// DatastreamResourceTimeout is the default timeout for the resource operations (max activation time + polling interval)
	DatastreamResourceTimeout = 180 * time.Minute
)

const (
	// DefaultUploadFilePrefix specifies default upload file prefix for supported connectors
	DefaultUploadFilePrefix = "ak"

	// DefaultUploadFileSuffix specifies default upload file suffix for supported connectors
	DefaultUploadFileSuffix = "ds"

	// MidgressDatasetField is the dataset field ID that gets automatically added/removed based on collect_midgress setting
	MidgressDatasetField = 2051
)

func resourceDatastream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatastreamCreate,
		ReadContext:   resourceDatastreamRead,
		UpdateContext: resourceDatastreamUpdate,
		DeleteContext: resourceDatastreamDelete,
		Timeouts: &schema.ResourceTimeout{
			Default: &DatastreamResourceTimeout,
		},
		CustomizeDiff: customdiff.All(
			validateConfig,
			validateContractAndGroupIDsByLogType,
			enforceComputedFieldsChange,
		),
		Schema: datastreamResourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

var datastreamResourceSchema = map[string]*schema.Schema{
	"log_type": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Optional:    true,
		Default:     string(datastream.LogTypeCDN),
		Description: "Type of logs for the stream",
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			string(datastream.LogTypeCDN),
			string(datastream.LogTypeAppSec),
			string(datastream.LogTypeAnswerX),
		}, true)),
	},
	"active": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Defining if stream should be active or not",
	},
	"collect_midgress": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Identifies if stream needs to collect midgress data",
	},
	"sampling_percentage": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 100),
		Description:  "The sample percentage of data that your stream will send to the destination",
	},
	"integration_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The integration mode for the stream (e.g., PM_DEPENDENT, HYBRID, DS_MANAGED)",
	},
	"delivery_configuration": {
		Type:        schema.TypeSet,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Elem:        configResource,
		Description: "Provides information about the configuration related to logs (format, file names, delivery frequency)",
	},
	"contract_id": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
		Description:      "Identifies the contract that has access to the product. Optional for CDN log type. Required for APPSEC and ANSWERX. Whitespace-only values are treated as omitted for CDN.",
	},
	"created_by": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The username who created the stream",
	},
	"created_date": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The date and time when the stream was created",
	},
	"dataset_fields": {
		Type:             schema.TypeList,
		Optional:         true,
		DiffSuppressFunc: isOrderDifferent,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "A list of data set fields selected from the associated template that the stream monitors in logs. The order of the identifiers define how the value for these fields appear in the log lines",
	},
	"notification_emails": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: tf.ValidateEmail,
		},
		Description: "List of email addresses where the system sends notifications about activations and deactivations of the stream",
	},
	"group_id": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		DiffSuppressFunc: tf.FieldPrefixSuppress("grp_"),
		Description:      "Identifies the group that has access to the product and for which the stream configuration was created. Optional for CDN log type. Required for APPSEC and ANSWERX. On update, this value is not sent to the API.",
	},
	"modified_by": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The username who modified the stream",
	},
	"modified_date": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The date and time when the stream was modified",
	},
	"papi_json": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The configuration in JSON format that can be copy-pasted into PAPI configuration to enable datastream behavior",
	},
	"product_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The ID of the product for which the stream was created",
	},
	"properties": {
		Type:             schema.TypeList,
		Optional:         true,
		DiffSuppressFunc: isPropertiesOrderDifferent,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			DiffSuppressFunc: tf.FieldPrefixSuppress("prp_"),
		},
		Description: "Identifies the properties monitored in the stream",
	},
	"app_sec_configs": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
		Description: "Identifies the application security configurations monitored in the stream",
	},
	"service_ids": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeInt,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		Description: "Identifies the AnswerX service IDs monitored in the stream.",
	},
	"stream_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the stream",
	},
	"stream_version": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Identifies the configuration version of the stream",
	},
	"latest_version": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Identifies the latest active configuration version of the stream",
	},
	"s3_connector": {
		Type:         schema.TypeSet,
		MaxItems:     1,
		ExactlyOneOf: ExactlyOneConnectorRule,
		Optional:     true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The access key identifier used to authenticate requests to the Amazon S3 account",
				},
				"bucket": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the Amazon S3 bucket",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The path to the folder within Amazon S3 bucket where logs will be stored",
				},
				"region": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The AWS region where Amazon S3 bucket resides",
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The secret access key identifier used to authenticate requests to the Amazon S3 account",
				},
			},
		},
	},
	"azure_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Access keys associated with Azure Storage account",
				},
				"account_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the Azure Storage account name",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"container_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Specifies the Azure Storage container name",
				},
				"path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The path to the folder within Azure Storage container where logs will be stored",
				},
			},
		},
	},
	"datadog_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The API key associated with Datadog account",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"service": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Description: "The service of the Datadog connector",
				},
				"source": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Description: "The source of the Datadog connector",
				},
				"tags": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Description: "The tags of the Datadog connector",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Datadog endpoint where logs will be stored",
				},
			},
		},
	},
	"splunk_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"compress_logs": {
					Type:        schema.TypeBool,
					Default:     true,
					Optional:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of custom header passed with the request to the destination",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request to the destination",
				},
				"event_collector_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The Event Collector token associated with Splunk account",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The raw event Splunk URL where logs will be stored",
				},
				"tls_hostname": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Description: "The hostname that verifies the server's certificate and matches the Subject Alternative Names (SANs) in the certificate. If not provided, DataStream fetches the hostname from the endpoint URL.",
				},
				"ca_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The certification authority (CA) certificate used to verify the origin server's certificate. If the certificate is not signed by a well-known certification authority, enter the CA certificate in the PEM format for verification.",
				},
				"client_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The digital certificate in the PEM format you want to use to authenticate requests to your destination. If you want to use mutual authentication, you need to provide both the client certificate and the client key (in the PEM format).",
				},
				"client_key": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The private key in the non-encrypted PKCS8 format you want to use to authenticate with the back-end server. If you want to use mutual authentication, you need to provide both the client certificate and the client key.",
				},
				"m_tls": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether mTLS is enabled or not.",
				},
			},
		},
	},
	"gcs_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bucket": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the storage bucket created in Google Cloud account",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"path": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Description: "The path to the folder within Google Cloud bucket where logs will be stored",
				},
				"private_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The contents of the JSON private key generated and downloaded in Google Cloud Storage account",
				},
				"project_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The unique ID of Google Cloud project",
				},
				"service_account_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the service account with the storage.object.create permission or Storage Object Creator role",
				},
			},
		},
	},
	"https_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Either NONE for no authentication, or BASIC for username and password authentication",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						string(datastream.AuthenticationTypeNone),
						string(datastream.AuthenticationTypeBasic),
					}, false)),
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Content type to pass in the log file header",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of custom header passed with the request to the destination",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request to the destination",
				},
				"password": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Sensitive:   true,
					Description: "Password set for custom HTTPS endpoint for authentication",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "URL where logs will be stored",
				},
				"user_name": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Sensitive:   true,
					Description: "Username used for authentication",
				},
				"tls_hostname": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Description: "The hostname that verifies the server's certificate and matches the Subject Alternative Names (SANs) in the certificate. If not provided, DataStream fetches the hostname from the endpoint URL.",
				},
				"ca_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The certification authority (CA) certificate used to verify the origin server's certificate. If the certificate is not signed by a well-known certification authority, enter the CA certificate in the PEM format for verification.",
				},
				"client_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The digital certificate in the PEM format you want to use to authenticate requests to your destination. If you want to use mutual authentication, you need to provide both the client certificate and the client key (in the PEM format).",
				},
				"client_key": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The private key in the non-encrypted PKCS8 format you want to use to authenticate with the back-end server. If you want to use mutual authentication, you need to provide both the client certificate and the client key.",
				},
				"m_tls": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether mTLS is enabled or not.",
				},
			},
		},
	},
	"sumologic_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"collector_code": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The unique HTTP collector code of Sumo Logic endpoint",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Default:     true,
					Optional:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Content type to pass in the log file header",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of custom header passed with the request to the destination",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request to the destination",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Sumo Logic collection endpoint where logs will be stored",
				},
			},
		},
	},
	"oracle_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The access key identifier used to authenticate requests to the Oracle Cloud account",
				},
				"bucket": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the Oracle Cloud Storage bucket",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether the logs should be compressed",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector",
				},
				"namespace": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The namespace of Oracle Cloud Storage account",
				},
				"path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The path to the folder within your Oracle Cloud Storage bucket where logs will be stored",
				},
				"region": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Oracle Cloud Storage region where bucket resides",
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The secret access key identifier used to authenticate requests to the Oracle Cloud account",
				},
			},
		},
	},
	"loggly_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Loggly bulk endpoint URL in the https://hostname.loggly.com/bulk/ format. Set the endpoint code in the authToken field instead of providing it in the URL. You can use Akamaized property hostnames as endpoint URLs. See Stream logs to Loggly.",
				},
				"auth_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The unique HTTP code for your Loggly bulk endpoint.",
				},
				"tags": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The tags you can use to segment and filter log events in Loggly. See Tags in the Loggly documentation.",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The type of the resource passed in the request's custom header. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-readable name for the request's custom header, containing only alphanumeric, dash, and underscore characters. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request that contains information about the client connection. For details, see Additional options in the DataStream user guide.",
				},
			},
		},
	},
	"new_relic_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "A New Relic endpoint URL you want to send your logs to. The endpoint URL should follow the https://<newrelic.com>/log/v1/ format format. See Introduction to the Log API https://docs.newrelic.com/docs/logs/log-api/introduction-log-api/ if you want to retrieve your New Relic endpoint URL.",
				},
				"auth_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Your Log API token for your account in New Relic.",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The type of the resource passed in the request's custom header. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-readable name for the request's custom header, containing only alphanumeric, dash, and underscore characters. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request that contains information about the client connection. For details, see Additional options in the DataStream user guide.",
				},
			},
		},
	},
	"elasticsearch_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the connector.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Elasticsearch bulk endpoint URL in the https://hostname.elastic-cloud.com:9243/_bulk/ format. Set indexName in the appropriate field instead of providing it in the URL. You can use Akamaized property hostnames as endpoint URLs. See Stream logs to Elasticsearch.",
				},
				"user_name": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The Elasticsearch basic access authentication username.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The Elasticsearch basic access authentication password.",
				},
				"index_name": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The index name of the Elastic cloud where you want to store log files.",
				},
				"content_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The type of the resource passed in the request's custom header. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-readable name for the request's custom header, containing only alphanumeric, dash, and underscore characters. For details, see Additional options in the DataStream user guide.",
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request that contains information about the client connection. For details, see Additional options in the DataStream user guide.",
				},
				"tls_hostname": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Description: "The hostname that verifies the server's certificate and matches the Subject Alternative Names (SANs) in the certificate. If not provided, DataStream fetches the hostname from the endpoint URL.",
				},
				"ca_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The certification authority (CA) certificate used to verify the origin server's certificate. If the certificate is not signed by a well-known certification authority, enter the CA certificate in the PEM format for verification.",
				},
				"client_cert": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The PEM-formatted digital certificate you want to authenticate requests to your destination with. If you want to use mutual authentication, you need to provide both the client certificate and the client key.",
				},
				"client_key": {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Sensitive:   true,
					Description: "The private key in the non-encrypted PKCS8 format you want to use to authenticate with the backend server. If you want to use mutual authentication, you need to provide both the client certificate and the client key.",
				},
				"m_tls": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Indicates whether mTLS is enabled or not.",
				},
			},
		},
	},
	"s3_compatible_connector": {
		Type:         schema.TypeSet,
		MaxItems:     1,
		ExactlyOneOf: ExactlyOneConnectorRule,
		Optional:     true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The access key identifier of the S3-compatible object storage bucket.",
				},
				"bucket": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the S3-compatible object storage bucket.",
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "Enables gzip compression for a log file sent to a destination. This value is always true for this destination type.",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the destination.",
				},
				"path": {
					Type:        schema.TypeString,
					Default:     "",
					Optional:    true,
					Description: "The path to the folder within your S3-compatible object storage bucket where you want to store logs. Optional field.",
				},
				"region": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The physical storage location of your S3-compatible object storage bucket.",
				},
				"secret_access_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The secret access key identifier of the S3-compatible object storage bucket.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The scheme-qualified host of your S3-compatible object storage bucket.",
				},
			},
		},
	},
	"trafficpeak_connector": {
		Type:         schema.TypeSet,
		MaxItems:     1,
		ExactlyOneOf: ExactlyOneConnectorRule,
		Optional:     true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"authentication_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Only BASIC authentication is supported for TrafficPeak destination.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						string(datastream.AuthenticationTypeBasic),
					}, false)),
				},
				"compress_logs": {
					Type:        schema.TypeBool,
					Default:     true,
					Optional:    true,
					Description: "Enables gzip compression for a log file sent to a destination. The value is true by default.",
				},
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The destination's name.",
				},
				"content_type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The type of the resource passed in the request's custom header. - Supported headers: `application/json` or `application/json; charset=utf-8`.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						string(datastream.TrafficPeakContentTypeJSON),
						string(datastream.TrafficPeakContentTypeJSONUTF8),
					}, false)),
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-readable name for the request's custom header, containing only alphanumeric, dash, and underscore characters. Optional field.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9_-]+$`),
						"custom_header_name must contain only alphanumeric characters, dashes, and underscores",
					)),
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request that contains information about the client connection. Optional field.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Enter the password you set in your TrafficPeak endpoint for authentication.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Enter the Hydrolix endpoint URL in the https://<host>/ingest/event?table=<tablename>&token=<token> format, where the token is the HTTP streaming ingest token, and the tablename is the Hydrolix data set table name.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
						regexp.MustCompile(`^https://[^/]+/ingest/event\?table=[^&]+&token=.+$`),
						"endpoint must be in the format https://<host>/ingest/event?table=<tablename>&token=<token>",
					)),
				},
				"user_name": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Enter the valid username you set in your TrafficPeak endpoint for authentication.",
				},
			},
		},
	},
	"dynatrace_connector": {
		Type:         schema.TypeSet,
		MaxItems:     1,
		ExactlyOneOf: ExactlyOneConnectorRule,
		Optional:     true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The destination's name.",
				},
				"endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Dynatrace Ingestion API endpoint URL in the https://{dynatrace-environment-id}.live.dynatrace.com/api/v2/logs/ingest format.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
						regexp.MustCompile(`^https://[^/]+\.live\.dynatrace\.com/api/v2/logs/ingest$`),
						"endpoint must be in the format https://{dynatrace-environment-id}.live.dynatrace.com/api/v2/logs/ingest",
					)),
				},
				"api_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The Dynatrace Log Ingest access token.",
				},
				"custom_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A human-readable name for the request's custom header, containing only alphanumeric, dash, and underscore characters. For details, see Additional options in the DataStream user guide.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(
						regexp.MustCompile(`^[A-Za-z0-9_-]+$`),
						"custom_header_name must contain only alphanumeric characters, dashes, and underscores",
					)),
				},
				"custom_header_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The custom header's contents passed with the request that contains information about the client connection. For details, see Additional options in the DataStream user guide.",
				},
			},
		},
	},
}

var configResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"field_delimiter": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A delimiter that you use to separate data set fields in log lines",
		},
		"format": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The format in which logs will be received",
		},
		"frequency": {
			Type:     schema.TypeSet,
			MinItems: 1,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"interval_in_secs": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "The time in seconds after which the system bundles log lines into a file and sends it to a destination",
					},
				},
			},
			Description: "The frequency of collecting logs from each uploader and sending these logs to a destination",
		},
		"upload_file_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     DefaultUploadFilePrefix,
			Description: "The prefix of the log file sent to a destination. Applies only to file-based connectors such as S3 and Azure. Not used by HTTP-based connectors.",
		},
		"upload_file_suffix": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     DefaultUploadFileSuffix,
			Description: "The suffix of the log file sent to a destination. Applies only to file-based connectors such as S3 and Azure. Not used by HTTP-based connectors.",
		},
	},
}

// logTypeFieldRequirement defines which fields are required, forbidden, or optional for a log type.
type logTypeFieldRequirement struct {
	// fieldName is the Terraform schema field name
	fieldName string
	// required indicates whether the field must be provided and non-empty
	required bool
	// forbidden indicates whether the field must not be provided (mutually exclusive with required)
	forbidden bool
}

// logTypeConfig encapsulates the validation and processing rules for a specific log type.
type logTypeConfig struct {
	logType      datastream.LogType
	requirements []logTypeFieldRequirement
}

// logTypeConfigs is the configuration registry for all supported log types.
// NOTE: this registry currently models only collection-shaped requirements
// (required/forbidden fields with length checks), which matches the existing
// validation helpers in validateStreamTypeConfig and validateStreamTypeConfigFromDiff.
// If a new log type needs scalar validation (for example, a non-zero integer or
// a specific string/enum rule that is not equivalent to a collection presence check),
// this map alone is not sufficient; the validation logic must be extended as well.
var logTypeConfigs = map[datastream.LogType]*logTypeConfig{
	datastream.LogTypeCDN: {
		logType: datastream.LogTypeCDN,
		requirements: []logTypeFieldRequirement{
			{fieldName: "dataset_fields", required: true},
			{fieldName: "properties", required: true},
			{fieldName: "app_sec_configs", forbidden: true},
			{fieldName: "service_ids", forbidden: true},
		},
	},
	datastream.LogTypeAppSec: {
		logType: datastream.LogTypeAppSec,
		requirements: []logTypeFieldRequirement{
			{fieldName: "dataset_fields", forbidden: true},
			{fieldName: "properties", forbidden: true},
			{fieldName: "app_sec_configs", required: true},
			{fieldName: "service_ids", forbidden: true},
		},
	},
	datastream.LogTypeAnswerX: {
		logType: datastream.LogTypeAnswerX,
		requirements: []logTypeFieldRequirement{
			{fieldName: "dataset_fields", required: true},
			{fieldName: "properties", forbidden: true},
			{fieldName: "app_sec_configs", forbidden: true},
			{fieldName: "service_ids", required: true},
		},
	},
}

// Helper function for extracting the log type from terraform schema and applying default value if not set.
func getLogType(d *schema.ResourceData) (datastream.LogType, error) {
	logTypeVal, err := tf.GetStringValue("log_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return "", err
	}

	// default the log type to CDN if not provided for backwards compatibility
	var logType = datastream.LogTypeCDN
	if logTypeVal != "" {
		logType = datastream.LogType(strings.ToUpper(logTypeVal))
	}

	return logType, nil
}

// validateStreamTypeConfig validates that all required fields are present and forbidden fields are absent.
func validateStreamTypeConfig(logType datastream.LogType, fields map[string]any) error {
	config, exists := logTypeConfigs[logType]
	if !exists {
		return fmt.Errorf("unsupported log_type %q", logType)
	}

	for _, req := range config.requirements {
		value, hasValue := fields[req.fieldName]
		fieldLen := 0
		if hasValue {
			if length, ok := collectionLen(value); ok {
				fieldLen = length
			}
		}

		if req.required && (fieldLen == 0) {
			return fmt.Errorf("`%s` are required for log_type %q", req.fieldName, logType)
		}

		if req.forbidden && (fieldLen > 0) {
			return fmt.Errorf("cannot set `%s` when log_type is %q", req.fieldName, logType)
		}
	}

	return nil
}

// streamTypeConfig holds the log-type-specific fields extracted from the Terraform schema.
type streamTypeConfig struct {
	DatasetFields     []datastream.DatasetFieldID
	Properties        []datastream.PropertyID
	AppSecConfigs     []datastream.AppSecConfigID
	AnswerXServiceIDs []datastream.AnswerXServiceID
}

// Helper function for validating the stream type specific configuration based on the log type and returning the extracted configuration values.
func getStreamTypeConfig(d *schema.ResourceData, logType datastream.LogType) (*streamTypeConfig, error) {
	result := &streamTypeConfig{}
	var datasetFieldsIDs []datastream.DatasetFieldID
	var propertyIDs []datastream.PropertyID
	var appSecConfigIDs []datastream.AppSecConfigID
	var answerXServiceIDs []datastream.AnswerXServiceID

	datasetFieldsIDsList, err := tf.GetListValue("dataset_fields", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	datasetFieldsIDs = DatasetFieldListToDatasetFields(datasetFieldsIDsList)

	propertyIDsList, propertyIDsErr := tf.GetListValue("properties", d)
	if propertyIDsErr != nil && !errors.Is(propertyIDsErr, tf.ErrNotFound) {
		return nil, propertyIDsErr
	}

	appSecConfigsList, appSecConfigsErr := tf.GetListValue("app_sec_configs", d)
	if appSecConfigsErr != nil && !errors.Is(appSecConfigsErr, tf.ErrNotFound) {
		return nil, appSecConfigsErr
	}

	answerXServiceIDsSet, answerXServiceIDsErr := tf.GetSetValue("service_ids", d)
	if answerXServiceIDsErr != nil && !errors.Is(answerXServiceIDsErr, tf.ErrNotFound) {
		return nil, answerXServiceIDsErr
	}

	// Build field map for validation
	validationFields := map[string]any{
		"dataset_fields":  datasetFieldsIDsList,
		"properties":      propertyIDsList,
		"app_sec_configs": appSecConfigsList,
		"service_ids":     answerXServiceIDsSet,
	}

	// Validate configuration against log type requirements
	if err := validateStreamTypeConfig(logType, validationFields); err != nil {
		return nil, err
	}

	// Process fields based on log type
	switch logType {
	case datastream.LogTypeCDN:
		propertyIDs, err = GetPropertiesList(propertyIDsList)
		if err != nil {
			return nil, err
		}
	case datastream.LogTypeAppSec:
		appSecConfigIDs, err = GetAppSecConfigIDs(appSecConfigsList)
		if err != nil {
			return nil, err
		}
	case datastream.LogTypeAnswerX:
		answerXServiceIDs, err = GetAnswerXServiceIDsFromSet(answerXServiceIDsSet)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported log_type %q", logType)
	}

	result.DatasetFields = datasetFieldsIDs
	result.Properties = propertyIDs
	result.AppSecConfigs = appSecConfigIDs
	result.AnswerXServiceIDs = answerXServiceIDs
	return result, nil
}

func resourceDatastreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Datastream", "resourceDatastreamCreate")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)
	logger.Debug("Creating stream")

	logType, err := getLogType(d)
	if err != nil {
		return diag.FromErr(err)
	}

	active, err := tf.GetBoolValue("active", d)
	if err != nil {
		return diag.FromErr(err)
	}

	collectMidgress, err := tf.GetBoolValue("collect_midgress", d)
	if err != nil {
		return diag.FromErr(err)
	}

	contractID, err := getContractIDForStream(d, logType)
	if err != nil {
		return diag.FromErr(err)
	}

	emailIDsList, err := tf.GetListValue("notification_emails", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
	}
	emailIDs := tf.InterfaceSliceToStringSlice(emailIDsList)

	if len(emailIDs) == 0 {
		emailIDs = nil
	}

	groupID, err := getGroupIDForStream(d, logType)
	if err != nil {
		return diag.FromErr(err)
	}

	streamFields, err := getStreamTypeConfig(d, logType)
	if err != nil {
		return diag.FromErr(err)
	}

	streamName, err := tf.GetStringValue("stream_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	connectors, err := GetConnectors(d, ExactlyOneConnectorRule)
	if err != nil {
		return diag.FromErr(err)
	}

	deliveryConfigSet, err := tf.GetSetValue("delivery_configuration", d)
	if err != nil {
		return diag.FromErr(err)
	}
	config, err := GetConfig(deliveryConfigSet)
	if err != nil {
		return diag.FromErr(err)
	}
	var httpsBaseConnectorName = GetConnectorNameWithOutFilePrefixSuffix(d, ConnectorsWithoutFilenameOptionsConfig)

	config, err = FilePrefixSuffixSet(httpsBaseConnectorName, config)
	if err != nil {
		return diag.FromErr(err)
	}
	// sampling_percentage is optional, so only get it if it exists
	var samplingPercentage int
	if value, exists := d.GetOk("sampling_percentage"); exists {
		samplingPercentage = value.(int)
	}
	req := datastream.CreateStreamRequest{
		StreamConfiguration: datastream.StreamConfiguration{
			CollectMidgress:       collectMidgress,
			DeliveryConfiguration: *config,
			Destination:           connectors,
			ContractID:            contractID,
			DatasetFields:         streamFields.DatasetFields,
			NotificationEmails:    emailIDs,
			GroupID:               groupID,
			Properties:            streamFields.Properties,
			StreamName:            streamName,
			SamplingPercentage:    samplingPercentage,
			AppSecConfigs:         streamFields.AppSecConfigs,
			AnswerXServiceIDs:     streamFields.AnswerXServiceIDs,
		},
		Activate: active,
		LogType:  logType,
	}

	res, err := client.CreateStream(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	streamID := res.StreamID
	d.SetId(strconv.FormatInt(streamID, 10))

	if active {
		_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusActivated)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatastreamRead(ctx, d, m)
}

// FilePrefixSuffixSet clears upload file prefix/suffix before create/update for connectors
// that do not support them (HTTP-based destinations). Those connectors never send these
// fields to the API. File-based connectors keep the configured values unchanged.
func FilePrefixSuffixSet(httpsBaseConnectorName string, config *datastream.DeliveryConfiguration) (*datastream.DeliveryConfiguration, error) {
	if collections.StringInSlice(ConnectorsWithoutFilenameOptionsConfig, httpsBaseConnectorName) {
		config.UploadFilePrefix = ""
		config.UploadFileSuffix = ""
	}
	return config, nil
}

func getOptionalContractID(d *schema.ResourceData) (string, error) {
	if value, ok := d.GetOk("contract_id"); ok {
		return normalizeContractID(value.(string)), nil
	}
	return "", nil
}

// resolveContractIDForRead stores the API contract_id when present; otherwise preserves an
// explicitly configured or prior state value when the API omits contractId (I#775).
func resolveContractIDForRead(apiContractID string, d *schema.ResourceData) string {
	if normalized := normalizeContractID(apiContractID); normalized != "" {
		return normalized
	}

	if configured, err := getOptionalContractID(d); err == nil && configured != "" {
		return configured
	}

	if configured := normalizeContractID(stringAttrFromStateOrConfig(d, "contract_id")); configured != "" {
		return configured
	}

	return ""
}

// resolveGroupIDForRead stores the API group_id when present; otherwise preserves an
// explicitly configured or prior state value when the API omits or zeroes groupId (I#775).
func resolveGroupIDForRead(apiGroupID int, d *schema.ResourceData) string {
	if apiGroupID > 0 {
		return strconv.Itoa(apiGroupID)
	}

	// Prefer a valid prior state value, then Terraform config. Config must be consulted after a
	// stale state of "0" (I#775 upgrade path) because GetOk treats "0" as set and would mask RawConfig.
	for _, candidate := range []string{
		groupIDStringFromResourceData(d),
		rawConfigStringAttr(d, "group_id"),
	} {
		if resolved := normalizedGroupIDString(candidate); resolved != "" {
			return resolved
		}
	}

	return ""
}

func groupIDStringFromResourceData(d *schema.ResourceData) string {
	value, ok := d.GetOk("group_id")
	if !ok {
		return ""
	}
	groupIDStr, ok := value.(string)
	if !ok {
		return ""
	}
	return groupIDStr
}

func rawConfigStringAttr(d *schema.ResourceData, key string) string {
	if !rawConfigAvailable(d) {
		return ""
	}
	val, ok := tf.NewRawConfig(d).GetOk(key)
	if !ok || val == nil {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(str)
}

// rawConfigAvailable reports whether ResourceData has usable raw config.
// During a standalone Read/refresh, RawConfig may be null/unknown; callers must
// not assume GetOk will succeed (and must not panic on that path).
func rawConfigAvailable(d *schema.ResourceData) bool {
	raw := d.GetRawConfig()
	return raw.IsKnown() && !raw.IsNull()
}

func normalizedGroupIDString(groupID string) string {
	groupIDStr := strings.TrimSpace(strings.TrimPrefix(groupID, "grp_"))
	if groupIDStr == "" || groupIDStr == "0" {
		return ""
	}

	parsedGroupID, err := parseGroupIDString(groupID)
	if err != nil || parsedGroupID < 1 {
		return ""
	}

	return strconv.Itoa(parsedGroupID)
}

// stringAttrFromStateOrConfig returns a string attribute from prior state or Terraform config (import refresh).
func stringAttrFromStateOrConfig(d *schema.ResourceData, key string) string {
	if value, ok := d.GetOk(key); ok {
		if str, ok := value.(string); ok {
			if trimmed := strings.TrimSpace(str); trimmed != "" {
				return trimmed
			}
		}
	}

	if rawConfigAvailable(d) {
		if val, ok := tf.NewRawConfig(d).GetOk(key); ok {
			if str, ok := val.(string); ok {
				return strings.TrimSpace(str)
			}
		}
	}

	return ""
}

// applyDeliveryConfigurationForRead normalizes delivery configuration from the API for state storage.
// Upload file prefix/suffix apply only to file-based connectors. For HTTP-based connectors they
// are not sent to the API; state is aligned to schema defaults solely to avoid perpetual drift
// from the shared delivery_configuration block. Preserve-when-API-omits runs only for file-based
// connectors.
func applyDeliveryConfigurationForRead(cfg *datastream.DeliveryConfiguration, d *schema.ResourceData, connectorKey string) {
	if cfg == nil {
		return
	}

	if collections.StringInSlice(ConnectorsWithoutFilenameOptionsConfig, connectorKey) {
		cfg.UploadFilePrefix = DefaultUploadFilePrefix
		cfg.UploadFileSuffix = DefaultUploadFileSuffix
		return
	}

	cfg.UploadFilePrefix = resolveUploadFilePrefixForRead(cfg.UploadFilePrefix, d)
	cfg.UploadFileSuffix = resolveUploadFileSuffixForRead(cfg.UploadFileSuffix, d)
}

// configuredUploadFilePart returns a single upload_file_prefix or upload_file_suffix
// value from state or RawConfig (key must be one of those attribute names).
func configuredUploadFilePart(d *schema.ResourceData, key string) (string, bool) {
	configSet, err := tf.GetSetValue("delivery_configuration", d)
	if err == nil && configSet.Len() > 0 {
		configMap, mapOK := configSet.List()[0].(map[string]interface{})
		if !mapOK {
			return "", false
		}

		value, _ := configMap[key].(string)
		return value, true
	}

	if !rawConfigAvailable(d) {
		return "", false
	}

	val, ok := tf.NewRawConfig(d).GetOk("delivery_configuration")
	if !ok {
		return "", false
	}

	blocks, ok := val.([]any)
	if !ok || len(blocks) == 0 {
		return "", false
	}

	block, ok := blocks[0].(map[string]any)
	if !ok {
		return "", false
	}

	value, _ := block[key].(string)
	return value, true
}

func resolveUploadFilePrefixForRead(apiPrefix string, d *schema.ResourceData) string {
	if strings.TrimSpace(apiPrefix) != "" {
		return apiPrefix
	}

	if prefix, ok := configuredUploadFilePart(d, "upload_file_prefix"); ok && strings.TrimSpace(prefix) != "" {
		return prefix
	}

	return DefaultUploadFilePrefix
}

func resolveUploadFileSuffixForRead(apiSuffix string, d *schema.ResourceData) string {
	if strings.TrimSpace(apiSuffix) != "" {
		return apiSuffix
	}

	if suffix, ok := configuredUploadFilePart(d, "upload_file_suffix"); ok && strings.TrimSpace(suffix) != "" {
		return suffix
	}

	return DefaultUploadFileSuffix
}

// resolveIntegrationTypeForRead stores the API value when present; otherwise preserves prior state/config.
func resolveIntegrationTypeForRead(apiIntegrationType string, d *schema.ResourceData) (string, bool) {
	if strings.TrimSpace(apiIntegrationType) != "" {
		return apiIntegrationType, true
	}

	if configured := stringAttrFromStateOrConfig(d, "integration_type"); configured != "" {
		return configured, true
	}

	return "", false
}

// resolveSamplingPercentageForRead stores the API value when present; otherwise preserves prior state/config.
func resolveSamplingPercentageForRead(apiSamplingPercentage int, d *schema.ResourceData) (int, bool) {
	if apiSamplingPercentage > 0 {
		return apiSamplingPercentage, true
	}

	if value, ok := d.GetOk("sampling_percentage"); ok {
		if sampling, ok := value.(int); ok && sampling > 0 {
			return sampling, true
		}
	}

	if rawConfigAvailable(d) {
		if val, ok := tf.NewRawConfig(d).GetOk("sampling_percentage"); ok {
			switch sampling := val.(type) {
			case int:
				if sampling > 0 {
					return sampling, true
				}
			case int64:
				if sampling > 0 {
					return int(sampling), true
				}
			}
		}
	}

	return 0, false
}

// validateStreamTypeConfigFromDiff performs plan-time validation of log-type specific requirements.
// It returns an error if known values violate the requirements, but skips validation for unknown values
// to allow planning to proceed. This ensures locally knowable constraints fail at plan time.
func validateStreamTypeConfigFromDiff(d *schema.ResourceDiff, logType datastream.LogType) error {
	config, exists := logTypeConfigs[logType]
	if !exists {
		return fmt.Errorf("unsupported log_type %q", logType)
	}

	for _, req := range config.requirements {
		if !d.NewValueKnown(req.fieldName) {
			// Unknown values are resolved after apply. Defer both required and forbidden checks
			// so plan-time validation only enforces constraints on known collections.
			continue
		}

		val, exists := d.GetOkExists(req.fieldName)
		if !exists || val == nil {
			// Field not set in diff
			if req.required {
				if rawValueUnknownInDiff(d, req.fieldName) {
					continue
				}
				return fmt.Errorf("`%s` are required for log_type %q", req.fieldName, logType)
			}
			continue
		}

		fieldLen, ok := collectionLen(val)
		if !ok {
			// Can't determine length, skip validation for unknown values
			// This handles the case where values are unknown during planning
			continue
		}

		if req.required && fieldLen == 0 {
			if rawValueUnknownInDiff(d, req.fieldName) {
				continue
			}
			return fmt.Errorf("`%s` are required for log_type %q", req.fieldName, logType)
		}

		if req.forbidden && fieldLen > 0 {
			return fmt.Errorf("cannot set `%s` when log_type is %q", req.fieldName, logType)
		}
	}

	return nil
}

func collectionLen(val interface{}) (int, bool) {
	switch collection := val.(type) {
	case []interface{}:
		return len(collection), true
	case *schema.Set:
		return collection.Len(), true
	default:
		return 0, false
	}
}

func rawValueUnknownInDiff(d *schema.ResourceDiff, key string) bool {
	raw := d.GetRawConfig()
	if raw.IsNull() {
		return false
	}

	rawValue, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() {
		return false
	}

	return !rawValue.IsKnown()
}

func normalizeContractID(contractID string) string {
	return strings.TrimSpace(strings.TrimPrefix(contractID, "ctr_"))
}

func parseGroupIDString(groupID string) (int, error) {
	groupIDStr := strings.TrimSpace(strings.TrimPrefix(groupID, "grp_"))
	if groupIDStr == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func logTypeFromDiff(d *schema.ResourceDiff) datastream.LogType {
	logTypeVal, ok := d.GetOkExists("log_type")
	if !ok || logTypeVal == nil {
		return datastream.LogTypeCDN
	}

	logTypeStr := strings.TrimSpace(logTypeVal.(string))
	if logTypeStr == "" {
		return datastream.LogTypeCDN
	}

	return datastream.LogType(strings.ToUpper(logTypeStr))
}

func validateContractAndGroupIDsByLogType(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if isDatastreamDestroyDiff(d) {
		return nil
	}

	logType := logTypeFromDiff(d)
	if logType == datastream.LogTypeCDN {
		return validateCDNGroupIDWhenSet(groupIDStringFromDiff(d))
	}

	return validateNonCDNContractAndGroupIDValues(
		contractIDFromDiff(d),
		groupIDStringFromDiff(d),
		logType,
	)
}

func isDatastreamDestroyDiff(d *schema.ResourceDiff) bool {
	if d.Id() == "" {
		return false
	}

	// Destroy plans clear required configuration; stream_name is always set on a live stream.
	streamName, ok := d.GetOkExists("stream_name")
	return isStreamNameUnsetForDestroy(streamName, ok)
}

func isStreamNameUnsetForDestroy(streamName interface{}, ok bool) bool {
	return !ok || streamName == nil || strings.TrimSpace(streamName.(string)) == ""
}

func validateCDNGroupIDWhenSet(groupID string) error {
	groupIDStr := strings.TrimSpace(strings.TrimPrefix(groupID, "grp_"))
	if groupIDStr == "" {
		return nil
	}

	parsedGroupID, err := parseGroupIDString(groupID)
	if err != nil {
		return fmt.Errorf("invalid `group_id` %q: %s", groupIDStr, err)
	}
	if parsedGroupID < 1 {
		return fmt.Errorf("`group_id` must be at least 1 for log_type %q", datastream.LogTypeCDN)
	}

	return nil
}

func contractIDFromDiff(d *schema.ResourceDiff) string {
	contractID, ok := d.GetOkExists("contract_id")
	if !ok || contractID == nil {
		return ""
	}
	return normalizeContractID(contractID.(string))
}

func groupIDStringFromDiff(d *schema.ResourceDiff) string {
	groupID, ok := d.GetOkExists("group_id")
	if !ok || groupID == nil {
		return ""
	}
	return groupID.(string)
}

func validateNonCDNContractAndGroupIDValues(contractID, groupID string, logType datastream.LogType) error {
	if contractID == "" {
		return fmt.Errorf("`contract_id` is required for log_type %q", logType)
	}

	groupIDStr := strings.TrimSpace(strings.TrimPrefix(groupID, "grp_"))
	if groupIDStr == "" {
		return fmt.Errorf("`group_id` is required for log_type %q", logType)
	}

	parsedGroupID, err := parseGroupIDString(groupID)
	if err != nil {
		return fmt.Errorf("invalid `group_id` %q: %s", groupIDStr, err)
	}
	if parsedGroupID < 1 {
		return fmt.Errorf("`group_id` must be at least 1 for log_type %q", logType)
	}

	return nil
}

func getContractIDForStream(d *schema.ResourceData, logType datastream.LogType) (string, error) {
	contractID, err := getOptionalContractID(d)
	if err != nil {
		return "", err
	}
	if logType != datastream.LogTypeCDN && contractID == "" {
		return "", fmt.Errorf("`contract_id` is required for log_type %q", logType)
	}
	return contractID, nil
}

func getGroupIDForStream(d *schema.ResourceData, logType datastream.LogType) (int, error) {
	if value, ok := d.GetOk("group_id"); ok {
		groupIDStr := strings.TrimSpace(strings.TrimPrefix(value.(string), "grp_"))
		if groupIDStr == "" {
			if logType != datastream.LogTypeCDN {
				return 0, fmt.Errorf("`group_id` is required for log_type %q", logType)
			}
			return 0, nil
		}

		groupID, err := parseGroupIDString(value.(string))
		if err != nil {
			return 0, fmt.Errorf("invalid `group_id` %q: %s", groupIDStr, err)
		}
		if groupID < 1 {
			return 0, fmt.Errorf("`group_id` must be at least 1 for log_type %q", logType)
		}
		return groupID, nil
	}

	if logType != datastream.LogTypeCDN {
		return 0, fmt.Errorf("`group_id` is required for log_type %q", logType)
	}
	return 0, nil
}

func resourceDatastreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Datastream", "resourceDatastreamRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)
	logger.Debug("Reading a stream")

	streamID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	val, err := tf.GetStringValue("log_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	streamDetails, logType, err := getStreamForRead(ctx, client, logger, streamID, val)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})
	attrs["log_type"] = string(logType)
	attrs["active"] = streamDetails.StreamStatus == datastream.StreamStatusActivated
	attrs["collect_midgress"] = streamDetails.CollectMidgress
	attrs["contract_id"] = resolveContractIDForRead(streamDetails.ContractID, d)
	attrs["created_by"] = streamDetails.CreatedBy
	attrs["created_date"] = streamDetails.CreatedDate

	// Filter out midgress field (2051) from dataset_fields when storing in state
	// since it's controlled by collect_midgress setting, not by user configuration
	rawDatasetFieldsList := DataSetFieldsToList(streamDetails.DatasetFields)

	// Convert []int to []interface{} for filtering
	datasetFieldsInterface := make([]interface{}, len(rawDatasetFieldsList))
	for i, field := range rawDatasetFieldsList {
		datasetFieldsInterface[i] = field
	}
	filteredDatasetFields := filterMidgressDatasetField(datasetFieldsInterface)

	attrs["dataset_fields"] = filteredDatasetFields
	attrs["notification_emails"] = streamDetails.NotificationEmails
	attrs["latest_version"] = streamDetails.LatestVersion

	attrs["group_id"] = resolveGroupIDForRead(streamDetails.GroupID, d)
	attrs["modified_by"] = streamDetails.ModifiedBy
	attrs["modified_date"] = streamDetails.ModifiedDate
	attrs["papi_json"] = StreamIDToPapiJSON(streamDetails.StreamID)
	attrs["product_id"] = streamDetails.ProductID

	if logType == datastream.LogTypeAppSec {
		attrs["app_sec_configs"] = AppSecConfigsToIDsList(streamDetails.AppSecConfigs)
		attrs["properties"] = []interface{}{}
		attrs["service_ids"] = []interface{}{}
	}

	if logType == datastream.LogTypeCDN {
		attrs["properties"] = PropertyToList(streamDetails.Properties)
		attrs["app_sec_configs"] = []interface{}{}
		attrs["service_ids"] = []interface{}{}
	}

	if logType == datastream.LogTypeAnswerX {
		attrs["service_ids"] = AnswerXServiceIDsToList(streamDetails.AnswerXServiceIDs)
		attrs["properties"] = []interface{}{}
		attrs["app_sec_configs"] = []interface{}{}
	}

	attrs["stream_name"] = streamDetails.StreamName
	attrs["stream_version"] = streamDetails.StreamVersion
	if samplingPercentage, ok := resolveSamplingPercentageForRead(streamDetails.SamplingPercentage, d); ok {
		attrs["sampling_percentage"] = samplingPercentage
	}
	if integrationType, ok := resolveIntegrationTypeForRead(streamDetails.IntegrationType, d); ok {
		attrs["integration_type"] = integrationType
	}

	connectorKey, connectorProps, err := ConnectorToMap(streamDetails.Destination, d)
	if err != nil {
		return diag.FromErr(err)
	}
	if connectorKey != "" {
		attrs[connectorKey] = []interface{}{connectorProps}
	}

	applyDeliveryConfigurationForRead(&streamDetails.DeliveryConfiguration, d, connectorKey)

	attrs["delivery_configuration"] = ConfigToSet(streamDetails.DeliveryConfiguration)

	err = tf.SetAttrs(d, attrs)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// getStreamForRead fetches stream details, honoring explicit log_type when provided.
// If log_type is unavailable (for example, during import refresh), it probes known
// log types to determine the correct API path.
func getStreamForRead(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64, configuredLogType string) (*datastream.DetailedStreamVersion, datastream.LogType, error) {
	// if the log type is known, then we can try to fetch the stream.
	if configuredLogType != "" {
		logType := datastream.LogType(strings.ToUpper(configuredLogType))
		streamDetails, err := client.GetStream(ctx, datastream.GetStreamRequest{
			StreamID: streamID,
			LogType:  logType,
		})
		if err != nil {
			logger.Errorf("read stream '%d' with log_type %s failed: %T: %v", streamID, logType, err, err)
			return nil, "", err
		}
		return streamDetails, logType, nil
	}

	// if it's not known (e.g. during import refresh), we need to probe for it by trying to fetch stream details with known log types and seeing which one succeeds.
	streamDetails, logType, err := probeForStream(ctx, client, logger, streamID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to determine log type for stream '%d': %w", streamID, err)
	}
	return streamDetails, logType, nil

}

// probeForStream attempts to fetch stream details for known log types to determine the correct log type for the stream.
// This is used in cases where log type is not provided (e.g. during import refresh) to still be able to read stream details and populate state.
// If the API returns an unrecoverable error (e.g. 5xx status code) for a log type, it stops probing and returns the error instead of continuing to probe other log types.
func probeForStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64) (*datastream.DetailedStreamVersion, datastream.LogType, error) {

	var lastErr error
	for _, logType := range []datastream.LogType{datastream.LogTypeCDN, datastream.LogTypeAppSec, datastream.LogTypeAnswerX} {
		logger.Debugf("probing for stream '%d' with log_type %s", streamID, logType)
		streamDetails, err := client.GetStream(ctx, datastream.GetStreamRequest{
			StreamID: streamID,
			LogType:  logType,
		})
		if err == nil {
			return streamDetails, logType, nil
		}

		// debug print the error
		logger.Errorf("read stream '%d' with log_type %s failed: %T: %v", streamID, logType, err, err)

		if !shouldContinueProbingForLogType(err) {
			return nil, "", fmt.Errorf("received unrecoverable error while probing for stream '%d' with log_type %s: %w", streamID, logType, err)
		}

		lastErr = err
	}

	return nil, "", lastErr
}

// shouldContinueProbingForLogType determines whether to continue probing for stream log type based on the error returned from the API.
func shouldContinueProbingForLogType(err error) bool {
	apiErr, ok := extractDatastreamAPIError(err)
	if !ok {
		return true
	}
	// if the status code that comes back indicates a 5xx error, we should stop probing for the stream.
	if apiErr.StatusCode >= 500 && apiErr.StatusCode < 600 {
		return false
	}
	return true
}

// extractDatastreamAPIError attempts to extract a typed datastream API error.
// If the error is not of the expected type, it returns ok=false.
func extractDatastreamAPIError(err error) (*datastream.Error, bool) {
	var apiErr *datastream.Error
	if !errors.As(err, &apiErr) {
		return nil, false
	}
	return apiErr, true
}

func resourceDatastreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Datastream", "resourceDatastreamUpdate")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)
	logger.Debug("Updating stream")

	streamID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	logType, err := getLogType(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// it is not possible to edit stream while it is (de)activating
	currentStreamStatus, err := waitForStreamStatusChange(ctx, client, streamID, logType,
		datastream.StreamStatusDeactivated,
		datastream.StreamStatusActivated,
		datastream.StreamStatusInactive,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	isStreamActive := *currentStreamStatus == datastream.StreamStatusActivated

	var newActive bool
	if d.HasChange("active") {
		_, newActiveValue := d.GetChange("active")
		newActive = newActiveValue.(bool)
	} else {
		oldActiveValue, err := tf.GetBoolValue("active", d)
		if err != nil {
			return diag.FromErr(err)
		}
		newActive = oldActiveValue
	}

	if isStreamActive {
		if newActive {
			// stream is active and should be still active

			// update details
			err = updateStream(ctx, client, logger, streamID, d, isStreamActive, logType)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is activated because updating active stream causes its reactivation
			logger.Debugf("waiting for stream #%d activation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusActivated)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// stream is active and should be deactivated

			// deactivate stream first
			err = deactivateStream(ctx, client, logger, streamID, logType)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is deactivated
			logger.Debugf("waiting for stream #%d deactivation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusDeactivated)
			if err != nil {
				return diag.FromErr(err)
			}

			// update details (no waiting needed because stream is inactive)
			err = updateStream(ctx, client, logger, streamID, d, false, logType)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		// update details (no waiting needed because stream is inactive)

		err = updateStream(ctx, client, logger, streamID, d, isStreamActive, logType)
		if err != nil {
			return diag.FromErr(err)
		}

		if newActive {
			//stream is inactive and should be activated

			// activate stream first
			err = activateStream(ctx, client, logger, streamID, logType)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is deactivated
			logger.Debugf("waiting for stream #%d activation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusActivated)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDatastreamRead(ctx, d, m)
}

func updateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64, d *schema.ResourceData, isStreamActive bool, logType datastream.LogType) error {
	// if some configuration details changed
	if d.HasChangeExcept("active") {

		contractID, err := getContractIDForStream(d, logType)
		if err != nil {
			return err
		}

		emailIDsList, err := tf.GetListValue("notification_emails", d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return err
			}
		}
		emailIDs := tf.InterfaceSliceToStringSlice(emailIDsList)

		streamFields, err := getStreamTypeConfig(d, logType)
		if err != nil {
			return err
		}

		streamName, err := tf.GetStringValue("stream_name", d)
		if err != nil {
			return err
		}

		collectMidgress, err := tf.GetBoolValue("collect_midgress", d)
		if err != nil {
			return err
		}

		connectors, err := GetConnectors(d, ExactlyOneConnectorRule)
		if err != nil {
			return err
		}

		configSet, err := tf.GetSetValue("delivery_configuration", d)
		if err != nil {
			return err
		}
		config, err := GetConfig(configSet)
		if err != nil {
			return err
		}

		var httpsBaseConnectorName = GetConnectorNameWithOutFilePrefixSuffix(d, ConnectorsWithoutFilenameOptionsConfig)

		config, err = FilePrefixSuffixSet(httpsBaseConnectorName, config)
		if err != nil {
			return err
		}

		// sampling_percentage is optional, so only get it if it exists
		var samplingPercentage int
		if value, exists := d.GetOk("sampling_percentage"); exists {
			samplingPercentage = value.(int)
		}
		req := datastream.UpdateStreamRequest{
			StreamID: streamID,
			StreamConfiguration: datastream.StreamConfiguration{
				CollectMidgress:       collectMidgress,
				DeliveryConfiguration: *config,
				Destination:           connectors,
				ContractID:            contractID,
				DatasetFields:         streamFields.DatasetFields,
				NotificationEmails:    emailIDs,
				Properties:            streamFields.Properties,
				StreamName:            streamName,
				SamplingPercentage:    samplingPercentage,
				AppSecConfigs:         streamFields.AppSecConfigs,
				AnswerXServiceIDs:     streamFields.AnswerXServiceIDs,
			},
			Activate: isStreamActive,
			LogType:  logType,
		}

		_, err = client.UpdateStream(ctx, req)
		logger.Debugf("updating stream #%d details", streamID)
		return err
	}

	logger.Debugf("skipping updating stream #%d details", streamID)
	return nil
}

func deactivateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64, logType datastream.LogType) error {
	logger.Debug("deactivating stream")
	_, err := client.DeactivateStream(ctx, datastream.DeactivateStreamRequest{
		StreamID: streamID,
		LogType:  logType,
	})
	if err != nil {
		return err
	}

	logger.Debugf("waiting for the stream #%d to be deactivated", streamID)
	_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusDeactivated)
	return err
}

func activateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64, logType datastream.LogType) error {
	logger.Info("activating stream")
	_, err := client.ActivateStream(ctx, datastream.ActivateStreamRequest{
		StreamID: streamID,
		LogType:  logType,
	})
	if err != nil {
		return err
	}
	logger.Debugf("waiting for the stream #%d to be activated", streamID)
	_, err = waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusActivated)
	return err
}

func resourceDatastreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Datastream", "resourceDatastreamDelete")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)
	logger.Debug("Deleting stream")

	logType, err := getLogType(d)
	if err != nil {
		return diag.FromErr(err)
	}

	streamID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	streamDetails, err := client.GetStream(ctx, datastream.GetStreamRequest{
		StreamID: streamID,
		LogType:  logType,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	activationStatus := streamDetails.StreamStatus

	// if status == activating             - wait, deactivate, wait, delete
	// if status == activated              - deactivate, wait, delete
	// if status == deactivating           - wait, delete
	// if status == deactivated/inactive   - delete

	// if stream is activating we have to wait until activation finishes
	if activationStatus == datastream.StreamStatusActivating {
		_, err := waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusActivated)
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusActivated
	}

	// if stream is active - deactivate it
	if activationStatus == datastream.StreamStatusActivated {
		_, err := client.DeactivateStream(ctx, datastream.DeactivateStreamRequest{
			StreamID: streamID,
			LogType:  logType,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusDeactivating
	}

	// if stream is deactivating phase - wait until it completes
	if activationStatus == datastream.StreamStatusDeactivating {
		_, err := waitForStreamStatusChange(ctx, client, streamID, logType, datastream.StreamStatusDeactivated)
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusDeactivated
	}

	// if stream is inactive - delete it
	if activationStatus == datastream.StreamStatusDeactivated || activationStatus == datastream.StreamStatusInactive {
		err := client.DeleteStream(ctx, datastream.DeleteStreamRequest{
			StreamID: streamID,
			LogType:  logType,
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func waitForStreamStatusChange(ctx context.Context, client datastream.DS, streamID int64, logType datastream.LogType, expectedStatuses ...datastream.StreamStatus) (*datastream.StreamStatus, error) {
	expectedStatusesMap := map[datastream.StreamStatus]bool{}
	for _, status := range expectedStatuses {
		expectedStatusesMap[status] = true
	}

	getStreamReq := datastream.GetStreamRequest{
		StreamID: streamID,
		LogType:  logType,
	}

	streamDetails, err := client.GetStream(ctx, getStreamReq)
	if err != nil {
		return nil, err
	}

	_, ok := expectedStatusesMap[streamDetails.StreamStatus]
	for ; !ok; _, ok = expectedStatusesMap[streamDetails.StreamStatus] {
		select {
		case <-time.After(PollForActivationStatusChangeInterval):
			streamDetails, err = client.GetStream(ctx, getStreamReq)
			if err != nil {
				return nil, err
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("change status context terminated: %w", ctx.Err())
		}
	}

	return &streamDetails.StreamStatus, nil
}

func isOrderDifferent(_, oldIDValue, newIDValue string, d *schema.ResourceData) bool {
	key := "dataset_fields"

	logger := log.Get("DataStream", "isOrderDifferent")

	defaultDiff := func() bool {
		return oldIDValue == newIDValue
	}

	configSet, err := tf.GetSetValue("delivery_configuration", d)
	if err != nil {
		logger.Warn("unable to get config for datastream")
		return defaultDiff()
	}

	config, err := GetConfig(configSet)
	if err != nil {
		logger.Warn("unable to convert config to correct structure")
		return defaultDiff()
	}

	if !d.HasChange(key) || config.Format == datastream.FormatTypeStructured {
		return defaultDiff()
	}

	var emptyValueMarker struct{}

	oldDataset, newDataset := d.GetChange(key)

	oldDatasetList, ok := oldDataset.([]interface{})
	if !ok {
		logger.Warnf("%s in state is incorrect", key)
		return defaultDiff()
	}

	newDatasetList, ok := newDataset.([]interface{})
	if !ok {
		logger.Warnf("new %s is incorrect", key)
		return defaultDiff()
	}

	if len(oldDatasetList) != len(newDatasetList) {
		return defaultDiff()
	}

	oldMap := make(map[int]struct{})

	for _, oldV := range oldDatasetList {
		oldValue, ok := oldV.(int)
		if !ok {
			logger.Warnf("incorrect type in state's %s", key)
			return defaultDiff()
		}
		oldMap[oldValue] = emptyValueMarker
	}

	for _, newV := range newDatasetList {
		newValue, ok := newV.(int)
		if !ok {
			logger.Warnf("incorrect type in upcoming %s", key)
			return defaultDiff()
		}

		if _, ok := oldMap[newValue]; ok {
			delete(oldMap, newValue)
		} else {
			return false
		}
	}

	return len(oldMap) == 0
}

func isPropertiesOrderDifferent(_, oldIDValue, newIDValue string, d *schema.ResourceData) bool {
	const key = "properties"

	logger := log.Get("DataStream", "isPropertiesOrderDifferent")

	defaultDiff := func() bool {
		return oldIDValue == newIDValue
	}

	if !d.HasChange(key) {
		return defaultDiff()
	}

	oldProperties, newProperties := d.GetChange(key)

	oldPropertyList, ok := oldProperties.([]interface{})
	if !ok {
		logger.Warnf("%s in state is incorrect", key)
		return defaultDiff()
	}

	newPropertyList, ok := newProperties.([]interface{})
	if !ok {
		logger.Warnf("new %s is incorrect", key)
		return defaultDiff()
	}

	if len(oldPropertyList) != len(newPropertyList) {
		return defaultDiff()
	}

	if same, ok := propertiesSameSet(oldPropertyList, newPropertyList); ok {
		return same
	}

	return defaultDiff()
}

// propertiesSameSet reports whether two property ID lists contain the same members regardless of order.
// The second return value is false when the lists cannot be compared (invalid element types).
func propertiesSameSet(oldPropertyList, newPropertyList []interface{}) (same bool, ok bool) {
	if len(oldPropertyList) != len(newPropertyList) {
		return false, true
	}

	oldMap := make(map[string]struct{}, len(oldPropertyList))

	for _, oldV := range oldPropertyList {
		oldValue, isString := oldV.(string)
		if !isString {
			return false, false
		}
		oldMap[normalizePropertyID(oldValue)] = struct{}{}
	}

	for _, newV := range newPropertyList {
		newValue, isString := newV.(string)
		if !isString {
			return false, false
		}

		normalizedNewValue := normalizePropertyID(newValue)
		if _, exists := oldMap[normalizedNewValue]; exists {
			delete(oldMap, normalizedNewValue)
		} else {
			return false, true
		}
	}

	return len(oldMap) == 0, true
}

func normalizePropertyID(propertyID string) string {
	return strings.TrimSpace(strings.TrimPrefix(propertyID, "prp_"))
}

func validateConfig(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip validation for destroy operations
	if isDatastreamDestroyDiff(d) {
		return nil
	}

	// Validate log-type specific field requirements at plan time
	logType := logTypeFromDiff(d)
	if err := validateStreamTypeConfigFromDiff(d, logType); err != nil {
		return err
	}

	// Validate that users don't manually specify midgress dataset field (2051)
	datasetFieldsResource, exists := d.GetOkExists("dataset_fields")
	if exists {
		datasetFieldsList, ok := datasetFieldsResource.([]interface{})
		if ok {
			for _, field := range datasetFieldsList {
				if fieldID, ok := field.(int); ok && fieldID == MidgressDatasetField {
					return fmt.Errorf("dataset field %d (midgress) cannot be manually specified in dataset_fields. Use collect_midgress = true to enable midgress data collection", MidgressDatasetField)
				}
			}
		} else {
			return fmt.Errorf("dataset_fields has unexpected type, expected []interface{} but got %T", datasetFieldsResource)
		}
	}
	// Validate that upload_file_prefix and upload_file_suffix are not customized for
	// HTTP-based connectors that do not support filename options.
	connectorName := ""
	for _, k := range ConnectorsWithoutFilenameOptionsConfig {
		connectorResource, exists := d.GetOkExists(k)
		if !exists {
			continue
		}

		connectorSet := connectorResource.(*schema.Set)
		if connectorSet.Len() > 0 {
			connectorName = k
			break
		}
	}

	connectorNameForLogFormatCheck := ""
	for _, c := range ConnectorsSupportOnlyJSONLogFormat {
		connectorResourceForLogFormatCheck, exists := d.GetOkExists(c)
		if !exists {
			continue
		}

		connectorSetForLogFormatCheck := connectorResourceForLogFormatCheck.(*schema.Set)
		if connectorSetForLogFormatCheck.Len() > 0 {
			connectorNameForLogFormatCheck = c
			break
		}
	}

	if connectorName == "" && connectorNameForLogFormatCheck == "" {
		return nil
	}

	configResource, exists := d.GetOkExists("delivery_configuration")
	if !exists {
		return nil
	}

	configSet := configResource.(*schema.Set)
	if configSet.Len() == 0 {
		return nil
	}

	config := configSet.List()[0].(map[string]interface{})
	prefixValue := config["upload_file_prefix"]
	suffixValue := config["upload_file_suffix"]
	logFormatValue := config["format"]

	if prefixValue.(string) != DefaultUploadFilePrefix || suffixValue.(string) != DefaultUploadFileSuffix {
		return fmt.Errorf("upload_file_prefix (%s) / upload_file_suffix (%s) cannot be used with %s", prefixValue, suffixValue, connectorName)
	}

	if connectorNameForLogFormatCheck != "" && logFormatValue != "JSON" {
		return fmt.Errorf("'%s' supports only JSON log format", connectorNameForLogFormatCheck)
	}

	return nil
}

func enforceComputedFieldsChange(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// List of computed fields to reset
	computedFields := []string{
		"stream_version",
		"latest_version",
		"modified_by",
		"modified_date",
		"integration_type",
	}

	// Get all changed keys
	changedKeys := d.GetChangedKeysPrefix("")

	// Remove "active" from changedKeys
	filteredKeys := make([]string, 0, len(changedKeys))
	for _, k := range changedKeys {
		if k != "active" {
			filteredKeys = append(filteredKeys, k)
		}
	}

	// Only reset computed fields if there are changes (excluding "active")
	if len(filteredKeys) > 0 {
		for _, f := range computedFields {
			if err := d.SetNewComputed(f); err != nil {
				return fmt.Errorf("cannot set new computed for '%s': %s", f, err)
			}
		}
	}

	return nil
}

// filterMidgressDatasetField removes the midgress dataset field (2051) from a list
func filterMidgressDatasetField(list []interface{}) []interface{} {
	filtered := make([]interface{}, 0, len(list))

	for _, item := range list {
		if intVal, ok := item.(int); ok && intVal == MidgressDatasetField {
			continue // Skip midgress field
		}
		filtered = append(filtered, item)
	}

	return filtered
}

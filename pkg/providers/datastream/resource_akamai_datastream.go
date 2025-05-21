package datastream

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/datastream"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/collections"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
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
	}

	// DatastreamResourceTimeout is the default timeout for the resource operations (max activation time + polling interval)
	DatastreamResourceTimeout = 180 * time.Minute
)

const (
	// DefaultUploadFilePrefix specifies default upload file prefix for supported connectors
	DefaultUploadFilePrefix = "ak"

	// DefaultUploadFileSuffix specifies default upload file suffix for supported connectors
	DefaultUploadFileSuffix = "ds"
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
		),
		Schema: datastreamResourceSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

var datastreamResourceSchema = map[string]*schema.Schema{
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
		Required:         true,
		DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
		Description:      "Identifies the contract that has access to the product",
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
		Required:         true,
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
		Required:         true,
		DiffSuppressFunc: tf.FieldPrefixSuppress("grp_"),
		Description:      "Identifies the group that has access to the product and for which the stream configuration was created",
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
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			DiffSuppressFunc: tf.FieldPrefixSuppress("prp_"),
		},
		Description: "Identifies the properties monitored in the stream",
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
		Optional:     true, //To DO it should be mandatory
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
			Description: "The prefix of the log file that will be send to a destination",
		},
		"upload_file_suffix": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     DefaultUploadFileSuffix,
			Description: "The suffix of the log file that will be send to a destination",
		},
	},
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

	active, err := tf.GetBoolValue("active", d)
	if err != nil {
		return diag.FromErr(err)
	}

	collectMidgress, err := tf.GetBoolValue("collect_midgress", d)
	if err != nil {
		return diag.FromErr(err)
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = strings.TrimPrefix(contractID, "ctr_")

	datasetFieldsIDsList, err := tf.GetListValue("dataset_fields", d)

	if err != nil {
		return diag.FromErr(err)
	}
	datasetFieldsIDs := DatasetFieldListToDatasetFields(datasetFieldsIDsList)

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

	groupIDStr, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := strconv.Atoi(strings.TrimPrefix(groupIDStr, "grp_"))
	if err != nil {
		return diag.FromErr(err)
	}

	propertyIDsList, err := tf.GetListValue("properties", d)
	if err != nil {
		return diag.FromErr(err)
	}
	propertyIDs, err := GetPropertiesList(propertyIDsList)
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

	req := datastream.CreateStreamRequest{
		StreamConfiguration: datastream.StreamConfiguration{
			CollectMidgress:       collectMidgress,
			DeliveryConfiguration: *config,
			Destination:           connectors,
			ContractID:            contractID,
			DatasetFields:         datasetFieldsIDs,
			NotificationEmails:    emailIDs,
			GroupID:               groupID,
			Properties:            propertyIDs,
			StreamName:            streamName,
		},
		Activate: active,
	}

	res, err := client.CreateStream(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	streamID := res.StreamID
	d.SetId(strconv.FormatInt(streamID, 10))

	if active {
		_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusActivated)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatastreamRead(ctx, d, m)
}

// FilePrefixSuffixSet is used to set the blank value for prefix and suffix for https based destination as https based destination does not support prefix and suffix
func FilePrefixSuffixSet(httpsBaseConnectorName string, config *datastream.DeliveryConfiguration) (*datastream.DeliveryConfiguration, error) {

	if collections.StringInSlice(ConnectorsWithoutFilenameOptionsConfig, httpsBaseConnectorName) {

		config.UploadFilePrefix = ""
		config.UploadFileSuffix = ""
	}
	return config, nil
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

	streamDetails, err := client.GetStream(ctx, datastream.GetStreamRequest{
		StreamID: streamID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})

	attrs["active"] = streamDetails.StreamStatus == datastream.StreamStatusActivated
	attrs["collect_midgress"] = streamDetails.CollectMidgress
	attrs["contract_id"] = streamDetails.ContractID
	attrs["created_by"] = streamDetails.CreatedBy
	attrs["created_date"] = streamDetails.CreatedDate
	attrs["dataset_fields"] = DataSetFieldsToList(streamDetails.DatasetFields)
	attrs["contract_id"] = streamDetails.ContractID
	attrs["notification_emails"] = streamDetails.NotificationEmails
	attrs["latest_version"] = streamDetails.LatestVersion

	attrs["group_id"] = strconv.Itoa(streamDetails.GroupID)
	attrs["modified_by"] = streamDetails.ModifiedBy
	attrs["modified_date"] = streamDetails.ModifiedDate
	attrs["papi_json"] = StreamIDToPapiJSON(streamDetails.StreamID)
	attrs["product_id"] = streamDetails.ProductID
	attrs["properties"] = PropertyToList(streamDetails.Properties)
	attrs["stream_name"] = streamDetails.StreamName
	attrs["stream_version"] = streamDetails.StreamVersion

	connectorKey, connectorProps, err := ConnectorToMap(streamDetails.Destination, d)
	if err != nil {
		return diag.FromErr(err)
	}
	if connectorKey != "" {
		attrs[connectorKey] = []interface{}{connectorProps}

		if collections.StringInSlice(ConnectorsWithoutFilenameOptionsConfig, connectorKey) {
			// some connectors don't allow setting upload file prefix/suffix (API is ignoring them),
			// but the documentation specifies default value for these fields (ak/ds respectively)
			// so these fields should have default values in terraform provider too

			// since we do validate connector and prefix/suffix combination in a validateConfig function
			// we have to take into account the fact that terraform would still see the change between remote (no prefixes set)
			// and local state (default prefixes set), so we have to ensure that local state has the default prefix/suffix set as well
			// here we insert default values to satisfy terraform diff
			streamDetails.DeliveryConfiguration.UploadFilePrefix = DefaultUploadFilePrefix
			streamDetails.DeliveryConfiguration.UploadFileSuffix = DefaultUploadFileSuffix
		}
	}

	attrs["delivery_configuration"] = ConfigToSet(streamDetails.DeliveryConfiguration)

	err = tf.SetAttrs(d, attrs)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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

	// it is not possible to edit stream while it is (de)activating
	currentStreamStatus, err := waitForStreamStatusChange(ctx, client, streamID,
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
			err = updateStream(ctx, client, logger, streamID, d, isStreamActive)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is activated because updating active stream causes its reactivation
			logger.Debugf("waiting for stream #%d activation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusActivated)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// stream is active and should be deactivated

			// deactivate stream first
			err = deactivateStream(ctx, client, logger, streamID)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is deactivated
			logger.Debugf("waiting for stream #%d deactivation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusDeactivated)
			if err != nil {
				return diag.FromErr(err)
			}

			// update details (no waiting needed because stream is inactive)
			err = updateStream(ctx, client, logger, streamID, d, false)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		// update details (no waiting needed because stream is inactive)

		err = updateStream(ctx, client, logger, streamID, d, isStreamActive)
		if err != nil {
			return diag.FromErr(err)
		}

		if newActive {
			//stream is inactive and should be activated

			// activate stream first
			err = activateStream(ctx, client, logger, streamID)
			if err != nil {
				return diag.FromErr(err)
			}

			// wait until stream is deactivated
			logger.Debugf("waiting for stream #%d activation", streamID)
			_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusActivated)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDatastreamRead(ctx, d, m)
}

func updateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64, d *schema.ResourceData, isStreamActive bool) error {
	// if some configuration details changed
	if d.HasChangeExcept("active") {

		contractID, err := tf.GetStringValue("contract_id", d)
		if err != nil {
			return err
		}

		datasetFieldsIDsList, err := tf.GetListValue("dataset_fields", d)
		if err != nil {
			return err
		}
		datasetFieldsIDs := DatasetFieldListToDatasetFields(datasetFieldsIDsList)

		emailIDsList, err := tf.GetListValue("notification_emails", d)
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return err
			}
		}
		emailIDs := tf.InterfaceSliceToStringSlice(emailIDsList)

		propertyIDsList, err := tf.GetListValue("properties", d)

		if err != nil {
			return err
		}
		propertyIDs, err := GetPropertiesList(propertyIDsList)
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

		req := datastream.UpdateStreamRequest{
			StreamID: streamID,
			StreamConfiguration: datastream.StreamConfiguration{
				CollectMidgress:       collectMidgress,
				DeliveryConfiguration: *config,
				Destination:           connectors,
				ContractID:            contractID,
				DatasetFields:         datasetFieldsIDs,
				NotificationEmails:    emailIDs,
				Properties:            propertyIDs,
				StreamName:            streamName,
			},
			Activate: isStreamActive,
		}

		_, err = client.UpdateStream(ctx, req)
		logger.Debugf("updating stream #%d details", streamID)
		return err
	}

	logger.Debugf("skipping updating stream #%d details", streamID)
	return nil
}

func deactivateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64) error {
	logger.Debug("deactivating stream")
	_, err := client.DeactivateStream(ctx, datastream.DeactivateStreamRequest{
		StreamID: streamID,
	})
	if err != nil {
		return err
	}

	logger.Debugf("waiting for the stream #%d to be deactivated", streamID)
	_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusDeactivated)
	return err
}

func activateStream(ctx context.Context, client datastream.DS, logger akalog.Interface, streamID int64) error {
	logger.Info("activating stream")
	_, err := client.ActivateStream(ctx, datastream.ActivateStreamRequest{
		StreamID: streamID,
	})
	if err != nil {
		return err
	}

	logger.Debugf("waiting for the stream #%d to be activated", streamID)
	_, err = waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusActivated)
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

	streamID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	streamDetails, err := client.GetStream(ctx, datastream.GetStreamRequest{
		StreamID: streamID,
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
		_, err := waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusActivated)
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusActivated
	}

	// if stream is active - deactivate it
	if activationStatus == datastream.StreamStatusActivated {
		_, err := client.DeactivateStream(ctx, datastream.DeactivateStreamRequest{
			StreamID: streamID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusDeactivating
	}

	// if stream is deactivating phase - wait until it completes
	if activationStatus == datastream.StreamStatusDeactivating {
		_, err := waitForStreamStatusChange(ctx, client, streamID, datastream.StreamStatusDeactivated)
		if err != nil {
			return diag.FromErr(err)
		}

		activationStatus = datastream.StreamStatusDeactivated
	}

	// if stream is inactive - delete it
	if activationStatus == datastream.StreamStatusDeactivated || activationStatus == datastream.StreamStatusInactive {
		err := client.DeleteStream(ctx, datastream.DeleteStreamRequest{
			StreamID: streamID,
		})

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func waitForStreamStatusChange(ctx context.Context, client datastream.DS, streamID int64, expectedStatuses ...datastream.StreamStatus) (*datastream.StreamStatus, error) {
	expectedStatusesMap := map[datastream.StreamStatus]bool{}
	for _, status := range expectedStatuses {
		expectedStatusesMap[status] = true
	}

	getStreamReq := datastream.GetStreamRequest{
		StreamID: streamID,
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

func validateConfig(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
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

	if connectorName == "" {
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

	if prefixValue.(string) != DefaultUploadFilePrefix || suffixValue.(string) != DefaultUploadFileSuffix {
		return fmt.Errorf("upload_file_prefix (%s) / upload_file_suffix (%s) cannot be used with %s", prefixValue, suffixValue, connectorName)
	}

	return nil
}

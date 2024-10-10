package datastream

import (
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// connectorTypeToResourceName maps ConnectorType to TF resource key
	connectorTypeToResourceName = map[datastream.DestinationType]string{
		datastream.DestinationTypeAzure:         "azure_connector",
		datastream.DestinationTypeDataDog:       "datadog_connector",
		datastream.DestinationTypeElasticsearch: "elasticsearch_connector",
		datastream.DestinationTypeGcs:           "gcs_connector",
		datastream.DestinationTypeHTTPS:         "https_connector",
		datastream.DestinationTypeLoggly:        "loggly_connector",
		datastream.DestinationTypeNewRelic:      "new_relic_connector",
		datastream.DestinationTypeOracle:        "oracle_connector",
		datastream.DestinationTypeS3:            "s3_connector",
		datastream.DestinationTypeSplunk:        "splunk_connector",
		datastream.DestinationTypeSumoLogic:     "sumologic_connector",
	}

	connectorMappers = map[datastream.DestinationType]func(datastream.Destination, map[string]interface{}) map[string]interface{}{
		datastream.DestinationTypeAzure:         MapAzureConnector,
		datastream.DestinationTypeDataDog:       MapDatadogConnector,
		datastream.DestinationTypeElasticsearch: MapElasticsearchConnector,
		datastream.DestinationTypeGcs:           MapGCSConnector,
		datastream.DestinationTypeHTTPS:         MapHTTPSConnector,
		datastream.DestinationTypeLoggly:        MapLogglyConnector,
		datastream.DestinationTypeNewRelic:      MapNewRelicConnector,
		datastream.DestinationTypeOracle:        MapOracleConnector,
		datastream.DestinationTypeS3:            MapS3Connector,
		datastream.DestinationTypeSplunk:        MapSplunkConnector,
		datastream.DestinationTypeSumoLogic:     MapSumoLogicConnector,
	}

	connectorGetters = map[string]func(map[string]interface{}) datastream.AbstractConnector{
		"azure_connector":         GetAzureConnector,
		"datadog_connector":       GetDatadogConnector,
		"elasticsearch_connector": GetElasticsearchConnector,
		"gcs_connector":           GetGCSConnector,
		"https_connector":         GetHTTPSConnector,
		"loggly_connector":        GetLogglyConnector,
		"new_relic_connector":     GetNewRelicConnector,
		"oracle_connector":        GetOracleConnector,
		"s3_connector":            GetS3Connector,
		"splunk_connector":        GetSplunkConnector,
		"sumologic_connector":     GetSumoLogicConnector,
	}
)

// ConnectorToMap converts ConnectorDetails struct to map of properties
func ConnectorToMap(connector datastream.Destination, d *schema.ResourceData) (string, map[string]interface{}, error) {

	connectorDetails := connector
	connectorType := connectorDetails.DestinationType
	resourceKey, ok := connectorTypeToResourceName[connectorType]
	if !ok {
		return "", nil, fmt.Errorf("cannot find resource name for connector type: %s", connectorType)
	}

	// get connector set from .tf file (needed for secrets, keys)
	// when importing the resource, local configuration is initially empty
	localConnectorSet, err := tf.GetSetValue(resourceKey, d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return "", nil, err
	}

	var connectorItemProperties map[string]interface{}
	if localConnectorSet.Len() > 0 {
		connectorItemProperties = localConnectorSet.List()[0].(map[string]interface{})
	}

	// select proper mapper function and call it
	mapper, ok := connectorMappers[connectorType]
	if !ok {
		return "", nil, fmt.Errorf("cannot find mapper function for %s connector", resourceKey)
	}

	connectorProperties := mapper(connectorDetails, connectorItemProperties)
	return resourceKey, connectorProperties, nil
}

// GetConnectors builds Connectors list
func GetConnectors(d *schema.ResourceData, keys []string) (datastream.AbstractConnector, error) {
	// check which connector is present in .tf file
	connectorName, connectorResource, err := tf.GetExactlyOneOf(d, keys)
	if err != nil {
		return nil, fmt.Errorf("missing connector definition")
	}

	connectorSet, ok := connectorResource.(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("invalid connector data (%s)", connectorName)
	}

	if connectorSet.Len() == 0 {
		return nil, fmt.Errorf("no connectors for %s", connectorName)
	}

	connectorProperties := connectorSet.List()[0].(map[string]interface{})
	connectorResourceGetter, ok := connectorGetters[connectorName]
	if !ok {
		return nil, fmt.Errorf("cannot find getter function for %s connector", connectorName)
	}

	connector := connectorResourceGetter(connectorProperties)
	return connector, nil
}

// GetS3Connector builds S3Connector structure
func GetS3Connector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.S3Connector{
		AccessKey:       props["access_key"].(string),
		Bucket:          props["bucket"].(string),
		DisplayName:     props["display_name"].(string),
		Path:            props["path"].(string),
		Region:          props["region"].(string),
		SecretAccessKey: props["secret_access_key"].(string),
	}
}

// MapS3Connector selects fields needed for S3Connector
func MapS3Connector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":        "",
		"bucket":            c.Bucket,
		"compress_logs":     c.CompressLogs,
		"display_name":      c.DisplayName,
		"path":              c.Path,
		"region":            c.Region,
		"secret_access_key": "",
	}
	setNonNilItemsFromState(state, rv, "access_key", "secret_access_key")
	return rv
}

// GetAzureConnector builds AzureConnector structure
func GetAzureConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.AzureConnector{
		AccessKey:     props["access_key"].(string),
		AccountName:   props["account_name"].(string),
		DisplayName:   props["display_name"].(string),
		ContainerName: props["container_name"].(string),
		Path:          props["path"].(string),
	}
}

// MapAzureConnector selects fields needed for AzureConnector
func MapAzureConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":     "",
		"account_name":   c.AccountName,
		"compress_logs":  c.CompressLogs,
		"display_name":   c.DisplayName,
		"container_name": c.ContainerName,
		"path":           c.Path,
	}
	setNonNilItemsFromState(state, rv, "access_key")
	return rv
}

// GetDatadogConnector builds DatadogConnector structure
func GetDatadogConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.DatadogConnector{
		AuthToken:    props["auth_token"].(string),
		CompressLogs: props["compress_logs"].(bool),
		DisplayName:  props["display_name"].(string),
		Service:      props["service"].(string),
		Source:       props["source"].(string),
		Tags:         props["tags"].(string),
		Endpoint:     props["endpoint"].(string),
	}
}

// MapDatadogConnector selects fields needed for DatadogConnector
func MapDatadogConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"auth_token":    "",
		"compress_logs": c.CompressLogs,
		"display_name":  c.DisplayName,
		"service":       c.Service,
		"source":        c.Source,
		"tags":          c.Tags,
		"endpoint":      c.Endpoint,
	}
	setNonNilItemsFromState(state, rv, "auth_token")
	return rv
}

// GetSplunkConnector builds SplunkConnector structure
func GetSplunkConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.SplunkConnector{
		CompressLogs:        props["compress_logs"].(bool),
		DisplayName:         props["display_name"].(string),
		CustomHeaderName:    props["custom_header_name"].(string),
		CustomHeaderValue:   props["custom_header_value"].(string),
		EventCollectorToken: props["event_collector_token"].(string),
		Endpoint:            props["endpoint"].(string),
		TLSHostname:         props["tls_hostname"].(string),
		CACert:              props["ca_cert"].(string),
		ClientCert:          props["client_cert"].(string),
		ClientKey:           props["client_key"].(string),
	}
}

// MapSplunkConnector selects fields needed for SplunkConnector
func MapSplunkConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"compress_logs":         c.CompressLogs,
		"display_name":          c.DisplayName,
		"custom_header_name":    c.CustomHeaderName,
		"custom_header_value":   c.CustomHeaderValue,
		"event_collector_token": "",
		"endpoint":              c.Endpoint,
		"tls_hostname":          c.TLSHostname,
		"ca_cert":               "",
		"client_cert":           "",
		"client_key":            "",
		"m_tls":                 false,
	}
	if c.MTLS == "Enabled" {
		rv["m_tls"] = true
	}
	setNonNilItemsFromState(state, rv, "event_collector_token", "ca_cert", "client_cert", "client_key")
	return rv
}

// GetGCSConnector builds GCSConnector structure
func GetGCSConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.GCSConnector{
		Bucket:             props["bucket"].(string),
		DisplayName:        props["display_name"].(string),
		Path:               props["path"].(string),
		PrivateKey:         props["private_key"].(string),
		ProjectID:          props["project_id"].(string),
		ServiceAccountName: props["service_account_name"].(string),
	}
}

// MapGCSConnector selects fields needed for GCSConnector
func MapGCSConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"bucket":               c.Bucket,
		"compress_logs":        c.CompressLogs,
		"display_name":         c.DisplayName,
		"path":                 c.Path,
		"private_key":          "",
		"project_id":           c.ProjectID,
		"service_account_name": c.ServiceAccountName,
	}
	setNonNilItemsFromState(state, rv, "private_key")
	return rv
}

// GetHTTPSConnector builds CustomHTTPSConnector structure
func GetHTTPSConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.CustomHTTPSConnector{
		AuthenticationType: datastream.AuthenticationType(props["authentication_type"].(string)),
		CompressLogs:       props["compress_logs"].(bool),
		DisplayName:        props["display_name"].(string),
		ContentType:        props["content_type"].(string),
		CustomHeaderName:   props["custom_header_name"].(string),
		CustomHeaderValue:  props["custom_header_value"].(string),
		Password:           props["password"].(string),
		Endpoint:           props["endpoint"].(string),
		UserName:           props["user_name"].(string),
		TLSHostname:        props["tls_hostname"].(string),
		CACert:             props["ca_cert"].(string),
		ClientCert:         props["client_cert"].(string),
		ClientKey:          props["client_key"].(string),
	}
}

// MapHTTPSConnector selects fields needed for CustomHTTPSConnector
func MapHTTPSConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"authentication_type": c.AuthenticationType,
		"compress_logs":       c.CompressLogs,
		"display_name":        c.DisplayName,
		"content_type":        c.ContentType,
		"custom_header_name":  c.CustomHeaderName,
		"custom_header_value": c.CustomHeaderValue,
		"password":            "",
		"endpoint":            c.Endpoint,
		"user_name":           "",
		"tls_hostname":        c.TLSHostname,
		"ca_cert":             "",
		"client_cert":         "",
		"client_key":          "",
		"m_tls":               false,
	}
	if c.MTLS == "Enabled" {
		rv["m_tls"] = true
	}
	setNonNilItemsFromState(state, rv, "password", "user_name", "ca_cert", "client_cert", "client_key")
	return rv
}

// GetSumoLogicConnector builds SumoLogicConnector structure
func GetSumoLogicConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.SumoLogicConnector{
		CollectorCode:     props["collector_code"].(string),
		CompressLogs:      props["compress_logs"].(bool),
		DisplayName:       props["display_name"].(string),
		ContentType:       props["content_type"].(string),
		CustomHeaderName:  props["custom_header_name"].(string),
		CustomHeaderValue: props["custom_header_value"].(string),
		Endpoint:          props["endpoint"].(string),
	}
}

// MapSumoLogicConnector selects fields needed for SumoLogicConnector
func MapSumoLogicConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"collector_code":      "",
		"compress_logs":       c.CompressLogs,
		"display_name":        c.DisplayName,
		"content_type":        c.ContentType,
		"custom_header_name":  c.CustomHeaderName,
		"custom_header_value": c.CustomHeaderValue,
		"endpoint":            c.Endpoint,
	}
	setNonNilItemsFromState(state, rv, "collector_code")
	return rv
}

// GetOracleConnector builds OracleCloudStorageConnector structure
func GetOracleConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.OracleCloudStorageConnector{
		AccessKey:       props["access_key"].(string),
		Bucket:          props["bucket"].(string),
		DisplayName:     props["display_name"].(string),
		Namespace:       props["namespace"].(string),
		Path:            props["path"].(string),
		Region:          props["region"].(string),
		SecretAccessKey: props["secret_access_key"].(string),
	}
}

// MapOracleConnector selects fields needed for OracleCloudStorageConnector
func MapOracleConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":        "",
		"bucket":            c.Bucket,
		"compress_logs":     c.CompressLogs,
		"display_name":      c.DisplayName,
		"namespace":         c.Namespace,
		"path":              c.Path,
		"region":            c.Region,
		"secret_access_key": "",
	}
	setNonNilItemsFromState(state, rv, "access_key", "secret_access_key")
	return rv
}

// GetLogglyConnector builds LogglyConnector structure
func GetLogglyConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.LogglyConnector{
		AuthToken:         props["auth_token"].(string),
		DisplayName:       props["display_name"].(string),
		Endpoint:          props["endpoint"].(string),
		Tags:              props["tags"].(string),
		ContentType:       props["content_type"].(string),
		CustomHeaderName:  props["custom_header_name"].(string),
		CustomHeaderValue: props["custom_header_value"].(string),
	}
}

// MapLogglyConnector selects fields needed for LogglyConnector
func MapLogglyConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"auth_token":          "",
		"display_name":        c.DisplayName,
		"endpoint":            c.Endpoint,
		"tags":                c.Tags,
		"content_type":        c.ContentType,
		"custom_header_name":  c.CustomHeaderName,
		"custom_header_value": c.CustomHeaderValue,
	}
	setNonNilItemsFromState(state, rv, "auth_token")
	return rv
}

// GetNewRelicConnector builds NewRelicConnector structure
func GetNewRelicConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.NewRelicConnector{
		AuthToken:         props["auth_token"].(string),
		DisplayName:       props["display_name"].(string),
		Endpoint:          props["endpoint"].(string),
		ContentType:       props["content_type"].(string),
		CustomHeaderName:  props["custom_header_name"].(string),
		CustomHeaderValue: props["custom_header_value"].(string),
	}
}

// MapNewRelicConnector selects fields needed for NewRelicConnector
func MapNewRelicConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"auth_token":          "",
		"display_name":        c.DisplayName,
		"endpoint":            c.Endpoint,
		"content_type":        c.ContentType,
		"custom_header_name":  c.CustomHeaderName,
		"custom_header_value": c.CustomHeaderValue,
	}
	setNonNilItemsFromState(state, rv, "auth_token")
	return rv
}

// GetElasticsearchConnector builds ElasticsearchConnector structure
func GetElasticsearchConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.ElasticsearchConnector{
		DisplayName:       props["display_name"].(string),
		Endpoint:          props["endpoint"].(string),
		IndexName:         props["index_name"].(string),
		UserName:          props["user_name"].(string),
		Password:          props["password"].(string),
		ContentType:       props["content_type"].(string),
		CustomHeaderName:  props["custom_header_name"].(string),
		CustomHeaderValue: props["custom_header_value"].(string),
		TLSHostname:       props["tls_hostname"].(string),
		CACert:            props["ca_cert"].(string),
		ClientCert:        props["client_cert"].(string),
		ClientKey:         props["client_key"].(string),
	}
}

// MapElasticsearchConnector selects fields needed for ElasticsearchConnector
func MapElasticsearchConnector(c datastream.Destination, state map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"display_name":        c.DisplayName,
		"endpoint":            c.Endpoint,
		"index_name":          c.IndexName,
		"user_name":           "",
		"password":            "",
		"content_type":        c.ContentType,
		"custom_header_name":  c.CustomHeaderName,
		"custom_header_value": c.CustomHeaderValue,
		"tls_hostname":        c.TLSHostname,
		"ca_cert":             "",
		"client_cert":         "",
		"client_key":          "",
		"m_tls":               false,
	}
	if c.MTLS == "Enabled" {
		rv["m_tls"] = true
	}
	setNonNilItemsFromState(state, rv, "user_name", "password", "ca_cert", "client_cert", "client_key")
	return rv
}

func setNonNilItemsFromState(state map[string]interface{}, target map[string]interface{}, fields ...string) {
	for _, f := range fields {
		if state[f] != nil {
			target[f] = state[f]
		}
	}
}

// GetConnectorNameWithOutFilePrefixSuffix Returns destination name which does not contain the file prefix and suffix
func GetConnectorNameWithOutFilePrefixSuffix(d *schema.ResourceData, keys []string) string {

	connectorName, _, _ := tf.GetExactlyOneOf(d, keys)
	return connectorName
}

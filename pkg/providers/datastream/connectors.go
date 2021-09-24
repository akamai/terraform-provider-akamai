package datastream

import (
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// connectorTypeToResourceName maps ConnectorType to TF resource key
	connectorTypeToResourceName = map[datastream.ConnectorType]string{
		datastream.ConnectorTypeAzure:     "azure_connector",
		datastream.ConnectorTypeDataDog:   "datadog_connector",
		datastream.ConnectorTypeGcs:       "gcs_connector",
		datastream.ConnectorTypeHTTPS:     "https_connector",
		datastream.ConnectorTypeOracle:    "oracle_connector",
		datastream.ConnectorTypeS3:        "s3_connector",
		datastream.ConnectorTypeSplunk:    "splunk_connector",
		datastream.ConnectorTypeSumoLogic: "sumologic_connector",
	}

	connectorMappers = map[datastream.ConnectorType]func(datastream.ConnectorDetails, map[string]interface{}) map[string]interface{}{
		datastream.ConnectorTypeAzure:     MapAzureConnector,
		datastream.ConnectorTypeDataDog:   MapDatadogConnector,
		datastream.ConnectorTypeGcs:       MapGCSConnector,
		datastream.ConnectorTypeHTTPS:     MapHTTPSConnector,
		datastream.ConnectorTypeOracle:    MapOracleConnector,
		datastream.ConnectorTypeS3:        MapS3Connector,
		datastream.ConnectorTypeSplunk:    MapSplunkConnector,
		datastream.ConnectorTypeSumoLogic: MapSumoLogicConnector,
	}

	connectorGetters = map[string]func(map[string]interface{}) datastream.AbstractConnector{
		"azure_connector":     GetAzureConnector,
		"datadog_connector":   GetDatadogConnector,
		"gcs_connector":       GetGCSConnector,
		"https_connector":     GetHTTPSConnector,
		"oracle_connector":    GetOracleConnector,
		"s3_connector":        GetS3Connector,
		"splunk_connector":    GetSplunkConnector,
		"sumologic_connector": GetSumoLogicConnector,
	}
)

// ConnectorToMap converts ConnectorDetails struct to map of properties
func ConnectorToMap(connectors []datastream.ConnectorDetails, d *schema.ResourceData) (string, map[string]interface{}, error) {
	// api returned empty list of connectors
	if len(connectors) != 1 {
		return "", nil, nil
	}

	connectorDetails := connectors[0]
	connectorType := connectorDetails.ConnectorType
	resourceKey, ok := connectorTypeToResourceName[connectorType]
	if !ok {
		return "", nil, fmt.Errorf("cannot find resource name for connector type: %s", connectorType)
	}

	// get connector set from .tf file (needed for secrets, keys)
	// when importing the resource, local configuration is initially empty
	localConnectorSet, err := tools.GetSetValue(resourceKey, d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
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
func GetConnectors(d *schema.ResourceData, keys []string) ([]datastream.AbstractConnector, error) {
	// check which connector is present in .tf file
	connectorName, connectorResource, err := tools.GetExactlyOneOf(d, keys)
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
	return []datastream.AbstractConnector{connector}, nil
}

// GetS3Connector builds S3Connector structure
func GetS3Connector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.S3Connector{
		AccessKey:       props["access_key"].(string),
		Bucket:          props["bucket"].(string),
		ConnectorName:   props["connector_name"].(string),
		Path:            props["path"].(string),
		Region:          props["region"].(string),
		SecretAccessKey: props["secret_access_key"].(string),
	}
}

// MapS3Connector selects fields needed for S3Connector
func MapS3Connector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":        "",
		"bucket":            c.Bucket,
		"compress_logs":     c.CompressLogs,
		"connector_id":      c.ConnectorID,
		"connector_name":    c.ConnectorName,
		"path":              c.Path,
		"region":            c.Region,
		"secret_access_key": "",
	}

	if s["access_key"] != nil && s["secret_access_key"] != nil {
		rv["access_key"] = s["access_key"]
		rv["secret_access_key"] = s["secret_access_key"]
	}
	return rv
}

// GetAzureConnector builds AzureConnector structure
func GetAzureConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.AzureConnector{
		AccessKey:     props["access_key"].(string),
		AccountName:   props["account_name"].(string),
		ConnectorName: props["connector_name"].(string),
		ContainerName: props["container_name"].(string),
		Path:          props["path"].(string),
	}
}

// MapAzureConnector selects fields needed for AzureConnector
func MapAzureConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":     "",
		"account_name":   c.AccountName,
		"compress_logs":  c.CompressLogs,
		"connector_id":   c.ConnectorID,
		"connector_name": c.ConnectorName,
		"container_name": c.ContainerName,
		"path":           c.Path,
	}
	if s["access_key"] != nil {
		rv["access_key"] = s["access_key"]
	}
	return rv
}

// GetDatadogConnector builds DatadogConnector structure
func GetDatadogConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.DatadogConnector{
		AuthToken:     props["auth_token"].(string),
		CompressLogs:  props["compress_logs"].(bool),
		ConnectorName: props["connector_name"].(string),
		Service:       props["service"].(string),
		Source:        props["source"].(string),
		Tags:          props["tags"].(string),
		URL:           props["url"].(string),
	}
}

// MapDatadogConnector selects fields needed for DatadogConnector
func MapDatadogConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"auth_token":     "",
		"compress_logs":  c.CompressLogs,
		"connector_id":   c.ConnectorID,
		"connector_name": c.ConnectorName,
		"service":        c.Service,
		"source":         c.Source,
		"tags":           c.Tags,
		"url":            c.URL,
	}
	if s["auth_token"] != nil {
		rv["auth_token"] = s["auth_token"]
	}
	return rv
}

// GetSplunkConnector builds SplunkConnector structure
func GetSplunkConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.SplunkConnector{
		CompressLogs:        props["compress_logs"].(bool),
		ConnectorName:       props["connector_name"].(string),
		EventCollectorToken: props["event_collector_token"].(string),
		URL:                 props["url"].(string),
	}
}

// MapSplunkConnector selects fields needed for SplunkConnector
func MapSplunkConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"compress_logs":         c.CompressLogs,
		"connector_id":          c.ConnectorID,
		"connector_name":        c.ConnectorName,
		"event_collector_token": "",
		"url":                   c.URL,
	}
	if s["event_collector_token"] != nil {
		rv["event_collector_token"] = s["event_collector_token"]
	}
	return rv
}

// GetGCSConnector builds GCSConnector structure
func GetGCSConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.GCSConnector{
		Bucket:             props["bucket"].(string),
		ConnectorName:      props["connector_name"].(string),
		Path:               props["path"].(string),
		PrivateKey:         props["private_key"].(string),
		ProjectID:          props["project_id"].(string),
		ServiceAccountName: props["service_account_name"].(string),
	}
}

// MapGCSConnector selects fields needed for GCSConnector
func MapGCSConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"bucket":               c.Bucket,
		"compress_logs":        c.CompressLogs,
		"connector_id":         c.ConnectorID,
		"connector_name":       c.ConnectorName,
		"path":                 c.Path,
		"private_key":          "",
		"project_id":           c.ProjectID,
		"service_account_name": c.ServiceAccountName,
	}
	if s["private_key"] != nil {
		rv["private_key"] = s["private_key"]
	}
	return rv
}

// GetHTTPSConnector builds CustomHTTPSConnector structure
func GetHTTPSConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.CustomHTTPSConnector{
		AuthenticationType: datastream.AuthenticationType(props["authentication_type"].(string)),
		CompressLogs:       props["compress_logs"].(bool),
		ConnectorName:      props["connector_name"].(string),
		Password:           props["password"].(string),
		URL:                props["url"].(string),
		UserName:           props["user_name"].(string),
	}
}

// MapHTTPSConnector selects fields needed for CustomHTTPSConnector
func MapHTTPSConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"authentication_type": c.AuthenticationType,
		"compress_logs":       c.CompressLogs,
		"connector_id":        c.ConnectorID,
		"connector_name":      c.ConnectorName,
		"password":            "",
		"url":                 c.URL,
		"user_name":           "",
	}
	if s["password"] != nil && s["user_name"] != nil {
		rv["password"] = s["password"]
		rv["user_name"] = s["user_name"]
	}
	return rv
}

// GetSumoLogicConnector builds SumoLogicConnector structure
func GetSumoLogicConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.SumoLogicConnector{
		CollectorCode: props["collector_code"].(string),
		CompressLogs:  props["compress_logs"].(bool),
		ConnectorName: props["connector_name"].(string),
		Endpoint:      props["endpoint"].(string),
	}
}

// MapSumoLogicConnector selects fields needed for SumoLogicConnector
func MapSumoLogicConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	endpoint := tools.GetFirstNotEmpty(c.Endpoint, c.URL)

	rv := map[string]interface{}{
		"collector_code": "",
		"compress_logs":  c.CompressLogs,
		"connector_id":   c.ConnectorID,
		"connector_name": c.ConnectorName,
		"endpoint":       endpoint,
	}
	if s["collector_code"] != nil {
		rv["collector_code"] = s["collector_code"]
	}
	return rv
}

// GetOracleConnector builds OracleCloudStorageConnector structure
func GetOracleConnector(props map[string]interface{}) datastream.AbstractConnector {
	return &datastream.OracleCloudStorageConnector{
		AccessKey:       props["access_key"].(string),
		Bucket:          props["bucket"].(string),
		ConnectorName:   props["connector_name"].(string),
		Namespace:       props["namespace"].(string),
		Path:            props["path"].(string),
		Region:          props["region"].(string),
		SecretAccessKey: props["secret_access_key"].(string),
	}
}

// MapOracleConnector selects fields needed for OracleCloudStorageConnector
func MapOracleConnector(c datastream.ConnectorDetails, s map[string]interface{}) map[string]interface{} {
	rv := map[string]interface{}{
		"access_key":        "",
		"bucket":            c.Bucket,
		"compress_logs":     c.CompressLogs,
		"connector_id":      c.ConnectorID,
		"connector_name":    c.ConnectorName,
		"namespace":         c.Namespace,
		"path":              c.Path,
		"region":            c.Region,
		"secret_access_key": "",
	}
	if s["access_key"] != nil && s["secret_access_key"] != nil {
		rv["access_key"] = s["access_key"]
		rv["secret_access_key"] = s["secret_access_key"]
	}
	return rv
}

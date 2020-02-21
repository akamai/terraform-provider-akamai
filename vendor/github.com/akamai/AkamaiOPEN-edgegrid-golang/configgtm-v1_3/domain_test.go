package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func instantiateDomain() *Domain {

	domain := NewDomain(gtmTestDomain, "basic")
	domainData := []byte(`{
                                "cnameCoalescingEnabled" : false,
                                "defaultErrorPenalty" : 75,
                                "defaultHealthMax" : null,
                                "defaultHealthMultiplier" : null,
                                "defaultHealthThreshold" : null,
                                "defaultMaxUnreachablePenalty" : null,
                                "defaultSslClientCertificate" : null,
                                "defaultSslClientPrivateKey" : null,
                                "defaultTimeoutPenalty" : 25,
                                "defaultUnreachableThreshold" : null,
                                "emailNotificationList" : [ ],
                                "endUserMappingEnabled" : false,
                                "lastModified" : "2019-06-14T19:36:13.174+00:00",
                                "lastModifiedBy" : "operator",
                                "loadFeedback" : false,
                                "mapUpdateInterval" : 600,
                                "maxProperties" : 100,
                                "maxResources" : 9999,
                                "maxTestTimeout" : 60.0,
                                "maxTTL" : 3600,
                                "minPingableRegionFraction" : null,
                                "minTestInterval" : 0,
                                "minTTL" : 0,
                                "modificationComments" : "Add Property testproperty",
                                "name" : "gtmdomtest.akadns.net",
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "roundRobinPrefix" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "servermonitorPool" : null,
                                "type" : "basic",
                                "status" : {
                                        "message" : "Change Pending",
                                        "changeId" : "df6c04e4-6327-4e0f-8872-bfe9fb2693d2",
                                        "propagationStatus" : "PENDING",
                                        "propagationStatusDate" : "2019-06-14T19:36:13.174+00:00",
                                        "passingValidation" : true,
                                        "links" : [ {
                                                                "rel" : "self",
                                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                        } ]
                                },
                                "loadImbalancePercentage" : null,
                                "domainVersionId" : null,
                                "resources" : [ ],
                                "properties" : [ {
                                        "backupCName" : null,
                                        "backupIp" : null,
                                        "balanceByDownloadScore" : false,
                                        "cname" : null,
                                        "comments" : null,
                                        "dynamicTTL" : 300,
                                        "failoverDelay" : null,
                                        "failbackDelay" : null,
                                        "ghostDemandReporting" : false,
                                        "handoutMode" : "normal",
                                        "handoutLimit" : 1,
                                        "healthMax" : null,
                                        "healthMultiplier" : null,
                                        "healthThreshold" : null,
                                        "lastModified" : "2019-06-14T19:36:13.174+00:00",
                                        "livenessTests" : [ ],
                                        "loadImbalancePercentage" : null,
                                        "mapName" : null,
                                        "maxUnreachablePenalty" : null,
                                        "minLiveFraction" : null,
                                        "mxRecords" : [ ],
                                        "name" : "testproperty",
                                        "scoreAggregationType" : "median",
                                        "stickinessBonusConstant" : null,
                                        "stickinessBonusPercentage" : null,
                                        "staticTTL" : null,
                                        "trafficTargets" : [ {
                                                "datacenterId" : 3131,
                                                "enabled" : true,
                                                "weight" : 100.0,
                                                "handoutCName" : null,
                                                "name" : null,
                                                "servers" : [ "1.2.3.4" ]
                                        } ],
                                        "type" : "performance",
                                        "unreachableThreshold" : null,
                                        "useComputedTargets" : false,
                                        "weightedHashBitsForIPv4" : null,
                                        "weightedHashBitsForIPv6" : null,
                                        "ipv6" : false,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/testproperty"
                                        } ]
                                } ],
                                "datacenters" : [ {
                                        "datacenterId" : 3131,
                                        "nickname" : "testDC1",
                                        "scorePenalty" : 0,
                                        "city" : null,
                                        "stateOrProvince" : null,
                                        "country" : null,
                                        "latitude" : null,
                                        "longitude" : null,
                                        "cloneOf" : null,
                                        "virtual" : true,
                                        "defaultLoadObject" : null,
                                        "continent" : null,
                                        "servermonitorPool" : null,
                                        "servermonitorLivenessCount" : null,
                                        "servermonitorLoadCount" : null,
                                        "pingInterval" : null,
                                        "pingPacketSize" : null,
                                        "cloudServerTargeting" : false,
                                        "cloudServerHostHeaderOverride" : false,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                        } ]
                                } ],
                                "geographicMaps" : [ ],
                                "cidrMaps" : [ ],
                                "asMaps" : [ ],
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                }, {
                                        "rel" : "datacenters",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters"
                                }, {
                                        "rel" : "properties",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties"
                                }, {
                                        "rel" : "geographic-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps"
                                }, {
                                        "rel" : "cidr-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps"
                                }, {
                                        "rel" : "resources",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources"
                                }, {
                                        "rel" : "as-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps"
                                } ]
                       }`)
	jsonhooks.Unmarshal(domainData, domain)

	return domain

}

// Verify GetListDomains. Sould pass, e.g. no API errors and non nil list.
func TestListDomains(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains")
	mock.
		Get("/config-gtm/v1/domains").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "items" : [ {
                                "name" : "gtmdomtest.akadns.net",
                                "status" : "Change Pending",
                                "acgId" : "1-3CV382",
                                "lastModified" : "2019-06-06T19:07:20.000+00:00",
                                "lastModifiedBy" : "operator",
                                "changeId" : "c3e1b771-2500-40c9-a7da-6c3cdbce1936",
                                "activationState" : "PENDING",
                                "modificationComments" : "mock test",
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                } ]
                        } ]
                   }`)

	Init(config)

	domainsList, err := ListDomains()
	assert.NoError(t, err)
	assert.NotEqual(t, domainsList, nil)
	assert.Equal(t, "gtmdomtest.akadns.net", domainsList[0].Name)

}

// Verify GetDomain. Name hardcoded. Should pass, e.g. no API errors and domain returned
func TestGetDomain(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                          "cidrMaps": [], 
                          "datacenters": [
                              {
                                  "city": "Snæfellsjökull", 
                                  "cloneOf": null, 
                                  "cloudServerTargeting": false, 
                                  "continent": "EU", 
                                  "country": "IS", 
                                  "datacenterId": 3132, 
                                  "defaultLoadObject": {
                                       "loadObject": null, 
                                       "loadObjectPort": 0, 
                                       "loadServers": null
                                   }, 
                                   "latitude": 64.808, 
                                   "links": [
                                       {
                                            "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132", 
                                            "rel": "self"
                                       }
                                    ], 
                                    "longitude": -23.776, 
                                    "nickname": "property_test_dc2", 
                                    "stateOrProvince": null, 
                                    "virtual": true
                              }, 
                              {
                                    "city": "Philadelphia", 
                                    "cloneOf": null, 
                                    "cloudServerTargeting": true, 
                                    "continent": "NA", 
                                    "country": "US", 
                                    "datacenterId": 3133, 
                                    "defaultLoadObject": {
                                         "loadObject": null, 
                                         "loadObjectPort": 0, 
                                         "loadServers": null
                                    }, 
                                    "latitude": 39.95, 
                                    "links": [
                                        {
                                             "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3133", 
                                             "rel": "self"
                                        }
                                     ], 
                                     "longitude": -75.167, 
                                     "nickname": "property_test_dc3", 
                                     "stateOrProvince": null, 
                                     "virtual": true
                              }, 
                              {
                                     "city": "Downpat", 
                                     "cloneOf": null, 
                                     "cloudServerTargeting": false, 
                                     "continent": "EU", 
                                     "country": "GB", 
                                     "datacenterId": 3131, 
                                     "defaultLoadObject": {
                                           "loadObject": null, 
                                           "loadObjectPort": 0, 
                                           "loadServers": null
                                     }, 
                                     "latitude": 54.367, 
                                     "links": [
                                         {
                                              "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131", 
                                              "rel": "self"
                                         }
                                     ], 
                                     "longitude": -5.582, 
                                     "nickname": "property_test_dc1", 
                                     "stateOrProvince": "ha", 
                                     "virtual": true
                              }
                          ], 
                          "defaultErrorPenalty": 75, 
                          "defaultSslClientCertificate": null, 
                          "defaultSslClientPrivateKey": null, 
                          "defaultTimeoutPenalty": 25, 
                          "emailNotificationList": [], 
                          "geographicMaps": [], 
                          "lastModified": "2019-04-25T14:53:12.000+00:00", 
                          "lastModifiedBy": "operator", 
                          "links": [
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net", 
                                   "rel": "self"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters", 
                                   "rel": "datacenters"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties", 
                                   "rel": "properties"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps", 
                                   "rel": "geographic-maps"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps", 
                                   "rel": "cidr-maps"
                              }, 
                              {
                                   "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources", 
                                   "rel": "resources"
                              }
                          ], 
                          "loadFeedback": false, 
                          "loadImbalancePercentage": 10.0, 
                          "modificationComments": "Edit Property test_property", 
                          "name": "gtmdomtest.akadns.net", 
                          "properties": [
                               {
                                    "backupCName": null, 
                                    "backupIp": null, 
                                    "balanceByDownloadScore": false, 
                                    "cname": "www.boo.wow", 
                                    "comments": null, 
                                    "dynamicTTL": 300, 
                                    "failbackDelay": 0, 
                                    "failoverDelay": 0, 
                                    "handoutMode": "normal", 
                                    "healthMax": null, 
                                    "healthMultiplier": null, 
                                    "healthThreshold": null, 
                                    "ipv6": false, 
                                    "lastModified": "2019-04-25T14:53:12.000+00:00", 
                                    "links": [
                                         {
                                              "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property", 
                                              "rel": "self"
                                         }
                                    ], 
                                    "livenessTests": [
                                         {
                                               "disableNonstandardPortWarning": false, 
                                               "hostHeader": null, 
                                               "httpError3xx": true, 
                                               "httpError4xx": true, 
                                               "httpError5xx": true, 
                                               "name": "health check", 
                                               "requestString": null, 
                                               "responseString": null, 
                                               "sslClientCertificate": null, 
                                               "sslClientPrivateKey": null, 
                                               "testInterval": 60, 
                                               "testObject": "/status", 
                                               "testObjectPassword": null, 
                                               "testObjectPort": 80, 
                                               "testObjectProtocol": "HTTP", 
                                               "testObjectUsername": null, 
                                               "testTimeout": 25.0
                                         }
                                    ], 
                                    "loadImbalancePercentage": 10.0, 
                                    "mapName": null, 
                                    "maxUnreachablePenalty": null, 
                                    "mxRecords": [], 
                                    "name": "test_property", 
                                    "scoreAggregationType": "mean", 
                                    "staticTTL": 600, 
                                    "stickinessBonusConstant": null, 
                                    "stickinessBonusPercentage": 50, 
                                    "trafficTargets": [
                                         {
                                              "datacenterId": 3131, 
                                              "enabled": true, 
                                              "handoutCName": null, 
                                              "name": null, 
                                              "servers": [
                                                   "1.2.3.4", 
                                                   "1.2.3.5"
                                              ], 
                                              "weight": 50.0
                                         }, 
                                         {
                                              "datacenterId": 3132, 
                                              "enabled": true, 
                                              "handoutCName": "www.google.com", 
                                              "name": null, 
                                              "servers": [], 
                                              "weight": 25.0
                                         }, 
                                         {
                                              "datacenterId": 3133, 
                                              "enabled": true, 
                                              "handoutCName": "www.comcast.com", 
                                              "name": null, 
                                              "servers": [
                                                    "www.comcast.com"
                                              ], 
                                              "weight": 25.0
                                         }
                                    ], 
                                    "type": "weighted-round-robin", 
                                    "unreachableThreshold": null, 
                                    "useComputedTargets": false
                               }
                          ], 
                          "resources": [], 
                          "status": {
                               "changeId": "40e36abd-bfb2-4635-9fca-62175cf17007", 
                               "links": [
                                     {
                                          "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current", 
                                          "rel": "self"
                                     }
                               ], 
                               "message": "Current configuration has been propagated to all GTM nameservers", 
                               "passingValidation": true, 
                               "propagationStatus": "COMPLETE", 
                               "propagationStatusDate": "2019-04-25T14:54:00.000+00:00"
                          }, 
                          "type": "weighted"
                }`)

	Init(config)

	testDomain, err := GetDomain(gtmTestDomain)

	assert.NoError(t, err)
	assert.Equal(t, gtmTestDomain, testDomain.Name)

}

// Verify failed case for GetDomain. Should pass, e.g. no API errors and domain not found
func TestGetBadDomain(t *testing.T) {

	baddomainname := "baddomainname.me"
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + baddomainname)
	mock.
		Get("/config-gtm/v1/domains/"+baddomainname).
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                }`)

	Init(config)

	_, err := GetDomain(baddomainname)
	assert.Error(t, err)

}

// Test Create domain. Name is hardcoded so this will effectively be an update. What happens to existing?
func TestCreateDomain(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/")
	mock.
		Post("/config-gtm/v1/domains/").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "cnameCoalescingEnabled" : false,
                                "defaultErrorPenalty" : 75,
                                "defaultHealthMax" : null,
                                "defaultHealthMultiplier" : null,
                                "defaultHealthThreshold" : null,
                                "defaultMaxUnreachablePenalty" : null,
                                "defaultSslClientCertificate" : null,
                                "defaultSslClientPrivateKey" : null,
                                "defaultTimeoutPenalty" : 25,
                                "defaultUnreachableThreshold" : null,
                                "emailNotificationList" : [ ],
                                "endUserMappingEnabled" : false,
                                "lastModified" : "2019-06-24T18:48:57.787+00:00",
                                "lastModifiedBy" : "operator",
                                "loadFeedback" : false,
                                "mapUpdateInterval" : 0,
                                "maxProperties" : 0,
                                "maxResources" : 512,
                                "maxTestTimeout" : 0.0,
                                "maxTTL" : 0,
                                "minPingableRegionFraction" : null,
                                "minTestInterval" : 0,
                                "minTTL" : 0,
                                "modificationComments" : null,
                                "name" : "gtmdomtest.akadns.net",
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "roundRobinPrefix" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "servermonitorPool" : null,
                                "type" : "basic",
                                "status" : {
                                        "message" : "Change Pending",
                                        "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                        "propagationStatus" : "PENDING",
                                        "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                        "passingValidation" : true,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                        } ]
                                },
                                "loadImbalancePercentage" : null,
                                "domainVersionId" : null,
                                "resources" : [ ],
                                "properties" : [ ],
                                "datacenters" : [ ],
                                "geographicMaps" : [ ],
                                "cidrMaps" : [ ],
                                "asMaps" : [ ],
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                    }, {
                                        "rel" : "datacenters",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters"
                                    }, {
                                        "rel" : "properties",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties"
                                    }, {
                                        "rel" : "geographic-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps"
                                    }, {
                                        "rel" : "cidr-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps"
                                    }, {
                                        "rel" : "resources",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources"
                                    }, {
                                        "rel" : "as-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps"
                                    } ]
                                },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "539872cc-6ba6-4429-acd5-90bab7fb5e9d",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-24T18:48:57.787+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
                }`)

	Init(config)

	testDomain := NewDomain(gtmTestDomain, "basic")
	qArgs := make(map[string]string)

	statResponse, err := testDomain.Create(qArgs)
	require.NoError(t, err)
	assert.Equal(t, gtmTestDomain, statResponse.Resource.Name)

}

func TestUpdateDomain(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "cnameCoalescingEnabled" : false,
                                "defaultErrorPenalty" : 75,
                                "defaultHealthMax" : null,
                                "defaultHealthMultiplier" : null,
                                "defaultHealthThreshold" : null,
                                "defaultMaxUnreachablePenalty" : null,
                                "defaultSslClientCertificate" : null,
                                "defaultSslClientPrivateKey" : null,
                                "defaultTimeoutPenalty" : 25,
                                "defaultUnreachableThreshold" : null,
                                "emailNotificationList" : [ ],
                                "endUserMappingEnabled" : false,
                                "lastModified" : "2019-06-14T19:36:13.174+00:00",
                                "lastModifiedBy" : "operator",
                                "loadFeedback" : false,
                                "mapUpdateInterval" : 600,
                                "maxProperties" : 100,
                                "maxResources" : 9999,
                                "maxTestTimeout" : 60.0,
                                "maxTTL" : 3600,
                                "minPingableRegionFraction" : null,
                                "minTestInterval" : 0,
                                "minTTL" : 0,
                                "modificationComments" : "Add Property testproperty",
                                "name" : "gtmdomtest.akadns.net",
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "roundRobinPrefix" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "servermonitorPool" : null,
                                "type" : "basic",
                                "status" : {
                                        "message" : "Change Pending",
                                        "changeId" : "df6c04e4-6327-4e0f-8872-bfe9fb2693d2",
                                        "propagationStatus" : "PENDING",
                                        "propagationStatusDate" : "2019-06-14T19:36:13.174+00:00",
                                        "passingValidation" : true,
                                        "links" : [ {
                                                                "rel" : "self",
                                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                        } ]
                                },
                                "loadImbalancePercentage" : null,
                                "domainVersionId" : null,
                                "resources" : [ ],
                                "properties" : [ {
                                        "backupCName" : null,
                                        "backupIp" : null,
                                        "balanceByDownloadScore" : false,
                                        "cname" : null,
                                        "comments" : null,
                                        "dynamicTTL" : 300,
                                        "failoverDelay" : null,
                                        "failbackDelay" : null,
                                        "ghostDemandReporting" : false,
                                        "handoutMode" : "normal",
                                        "handoutLimit" : 1,
                                        "healthMax" : null,
                                        "healthMultiplier" : null,
                                        "healthThreshold" : null,
                                        "lastModified" : "2019-06-14T19:36:13.174+00:00",
                                        "livenessTests" : [ ],
                                        "loadImbalancePercentage" : null,
                                        "mapName" : null,
                                        "maxUnreachablePenalty" : null,
                                        "minLiveFraction" : null,
                                        "mxRecords" : [ ],
                                        "name" : "testproperty",
                                        "scoreAggregationType" : "median",
                                        "stickinessBonusConstant" : null,
                                        "stickinessBonusPercentage" : null,
                                        "staticTTL" : null,
                                        "trafficTargets" : [ {
                                                "datacenterId" : 3131,
                                                "enabled" : true,
                                                "weight" : 100.0,
                                                "handoutCName" : null,
                                                "name" : null,
                                                "servers" : [ "1.2.3.4" ]
                                        } ],
                                        "type" : "performance",
                                        "unreachableThreshold" : null,
                                        "useComputedTargets" : false,
                                        "weightedHashBitsForIPv4" : null,
                                        "weightedHashBitsForIPv6" : null,
                                        "ipv6" : false,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/testproperty"
                                        } ]
                                } ],
                                "datacenters" : [ {
                                        "datacenterId" : 3131,
                                        "nickname" : "testDC1",
                                        "scorePenalty" : 0,
                                        "city" : null,
                                        "stateOrProvince" : null,
                                        "country" : null,
                                        "latitude" : null,
                                        "longitude" : null,
                                        "cloneOf" : null,
                                        "virtual" : true,
                                        "defaultLoadObject" : null,
                                        "continent" : null,
                                        "servermonitorPool" : null,
                                        "servermonitorLivenessCount" : null,
                                        "servermonitorLoadCount" : null,
                                        "pingInterval" : null,
                                        "pingPacketSize" : null,
                                        "cloudServerTargeting" : false,
                                        "cloudServerHostHeaderOverride" : false,
                                        "links" : [ {
                                                "rel" : "self",
                                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                        } ]
                                } ],
                                "geographicMaps" : [ ],
                                "cidrMaps" : [ ],
                                "asMaps" : [ ],
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net"
                                }, {
                                        "rel" : "datacenters",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters"
                                }, {
                                        "rel" : "properties",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties"
                                }, {
                                        "rel" : "geographic-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps"
                                }, {
                                        "rel" : "cidr-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps"
                                }, {
                                        "rel" : "resources",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources"
                                }, {
                                        "rel" : "as-maps",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/as-maps"
                                } ]
                        },
                        "status" : {
                                      "message" : "Change Pending",
                                      "changeId" : "df6c04e4-6327-4e0f-8872-bfe9fb2693d2",
                                      "propagationStatus" : "PENDING",
                                      "propagationStatusDate" : "2019-06-14T19:36:13.174+00:00",
                                      "passingValidation" : true,
                                      "links" : [ {
                                              "rel" : "self",
                                              "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                      } ]
                        }            
                }`)

	Init(config)

	testDomain := instantiateDomain()
	//testDomain.MaxResources = 9999
	qArgs := make(map[string]string)
	statResp, err := testDomain.Update(qArgs)
	require.NoError(t, err)
	assert.Equal(t, statResp.ChangeId, "df6c04e4-6327-4e0f-8872-bfe9fb2693d2")

}

/* Future. Presently no domain Delete endpoint.
func TestDeleteDomain(t *testing.T) {

        defer gock.Off()

        mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/"+gtmTestDomain)
        mock.
                Delete("/config-gtm/v1/domains/"+gtmTestDomain).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
                BodyString(`{
                        "resource" : null,
                        "status" : {
                               "changeId": "40e36abd-bfb2-4635-9fca-62175cf17007",
                               "links": [
                                     {
                                          "href": "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
                                          "rel": "self"
                                     }
                               ],
                               "message": "Change Pending",
                               "passingValidation": true,
                               "propagationStatus": "PENDING",
                               "propagationStatusDate": "2019-04-25T14:54:00.000+00:00"
                          },
                }`)

        Init(config)

        getDomain := instantiateDomain()

        _, err := getDomain.Delete()
        assert.NoError(t, err)

}
*/

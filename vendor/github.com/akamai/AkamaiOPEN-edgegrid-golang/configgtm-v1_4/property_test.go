package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"

	"fmt"
)

var GtmTestProperty = "testproperty"

func instantiateProperty() *Property {

	property := NewProperty(GtmTestProperty)
	propertyData := []byte(`{
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
                        "lastModified" : "2019-06-14T19:46:17.818+00:00",
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
            }`)
	jsonhooks.Unmarshal(propertyData, property)

	return property

}

// Verify GetListProperties. Should pass, e.g. no API errors and non nil list.
func TestListProperties(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/properties").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "items": [ {
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
                        "lastModified" : "2019-06-14T19:46:17.818+00:00",
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
                    } ]
                }`)

	Init(config)

	propertyList, err := ListProperties(gtmTestDomain)
	assert.NoError(t, err)
	assert.NotEqual(t, propertyList, nil)

	if len(propertyList) > 0 {
		firstProp := propertyList[0]
		assert.Equal(t, firstProp.Name, GtmTestProperty)
	} else {
		assert.Equal(t, 0, 1, "ListProperties: empty list")
		fmt.Println("ListProperties: empty list")
	}

}

// Verify GetProperty. Name hardcoded. Should pass, e.g. no API errors and property returned
// Depends on CreateProperty
func TestGetProperty(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties/" + GtmTestProperty)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
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
                        "lastModified" : "2019-06-14T19:46:17.818+00:00",
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
               }`)

	Init(config)

	testProperty, err := GetProperty(GtmTestProperty, gtmTestDomain)

	assert.NoError(t, err)
	assert.Equal(t, GtmTestProperty, testProperty.Name)

}

// Verify failed case for GetProperty. Should pass, e.g. no API errors and domain not found
func TestGetBadProperty(t *testing.T) {

	boguspropertyname := "boguspropertyname"
	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties/" + boguspropertyname)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/properties/"+boguspropertyname).
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
               }`)
	Init(config)

	_, err := GetProperty("boguspropertyname", gtmTestDomain)
	assert.Error(t, err)

}

// Test Create property. Name is hardcoded so this will effectively be an update. What happens to existing?
func TestCreateProperty(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties/" + GtmTestProperty)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
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
                        "lastModified" : "2019-06-14T19:46:17.818+00:00",
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
                    },
                    "status" : {
                        "message" : "Change Pending",
                        "changeId" : "ca4ca7a9-2b87-4f7c-8ba6-0cb1571df325",
                        "propagationStatus" : "PENDING",
                        "propagationStatusDate" : "2019-06-14T19:46:18.461+00:00",
                        "passingValidation" : true,
                        "links" : [ {
                                            "rel" : "self",
                                            "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                        } ]
                    }
       }`)

	Init(config)

	// initialize required fields
	testProperty := NewProperty(GtmTestProperty)
	testProperty.ScoreAggregationType = "median"
	testProperty.Type = "performance"
	testProperty.HandoutLimit = 1
	testProperty.HandoutMode = "normal"
	testProperty.FailoverDelay = 0
	testProperty.FailbackDelay = 0
	testProperty.TrafficTargets = []*TrafficTarget{&TrafficTarget{DatacenterId: 3131, Enabled: true, Servers: []string{"1.2.3.4"}, Weight: 100.0}}

	statresp, err := testProperty.Create(gtmTestDomain)
	assert.NoError(t, err)

	assert.IsType(t, &Property{}, statresp.Resource)
	assert.Equal(t, GtmTestProperty, statresp.Resource.Name)

}

func TestUpdateProperty(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties/" + GtmTestProperty)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
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
                        "handoutLimit" : 999,
                        "healthMax" : null,
                        "healthMultiplier" : null,
                        "healthThreshold" : null,
                        "lastModified" : "2019-06-14T19:46:17.818+00:00",
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
                    },
                    "status" : {
                        "message" : "Change Pending",
                        "changeId" : "ca4ca7a9-2b87-4f7c-8ba6-0cb1571df325",
                        "propagationStatus" : "PENDING",
                        "propagationStatusDate" : "2019-06-14T19:46:18.461+00:00",
                        "passingValidation" : true,
                        "links" : [ {
                                            "rel" : "self",
                                            "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                        } ]
                    } 
                }`)

	Init(config)

	testProperty := instantiateProperty()

	testProperty.HandoutLimit = 999
	stat, err := testProperty.Update(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, stat.ChangeId, "ca4ca7a9-2b87-4f7c-8ba6-0cb1571df325")

}

func TestDeleteProperty(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/properties/" + GtmTestProperty)
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
               }`)
	Init(config)

	getProperty := instantiateProperty()

	_, err := getProperty.Delete(gtmTestDomain)
	assert.NoError(t, err)

}

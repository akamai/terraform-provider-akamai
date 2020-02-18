package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var GtmTestCidrMap = "testCidrMap"

func instantiateCidrMap() *CidrMap {

	cidrMap := NewCidrMap(GtmTestCidrMap)
	cidrMapData := []byte(`{
                               "assignments": [ {
                                       "blocks": [ "1.2.3.0/24" ],
                                       "datacenterId": 3134,
                                       "nickname": "Frostfangs and the Fist of First Men"
                               },
                               {
                                       "blocks": [ "1.2.4.0/24" ],
                                       "datacenterId": 3133,
                                       "nickname": "Winterfell"
                               } ],
                               "defaultDatacenter": {
                                       "datacenterId": 5400,
                                       "nickname": "All Other CIDR Blocks"
                               },
                               "links": [ {
                                       "href": "/config-gtm/v1/domains/example.akadns.net/cidr-maps/testCidrMap",
                                       "rel": "self"
                               } ],
                               "name": "testCidrMap"
              }`)
	jsonhooks.Unmarshal(cidrMapData, cidrMap)

	return cidrMap

}

// Verify ListCidrMap. Name hardcoded. Should pass, e.g. no API errors and resource returned
func TestListCidrMaps(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "items" : [ {
                               "assignments": [ {
                                       "blocks": [ "1.2.3.0/24" ],
                                       "datacenterId": 3134,
                                       "nickname": "Frostfangs and the Fist of First Men"
                               },
                               {
                                       "blocks": [ "1.2.4.0/24" ],
                                       "datacenterId": 3133,
                                       "nickname": "Winterfell"
                               } ],
                               "defaultDatacenter": {
                                       "datacenterId": 5400,
                                       "nickname": "All Other CIDR Blocks"
                               },
                               "links": [ {
                                       "href": "/config-gtm/v1/domains/example.akadns.net/cidr-maps/testCidrMap",
                                       "rel": "self"
                               } ],
                               "name": "testCidrMap"
                       } ]
               }`)

	Init(config)

	testCidrMap, err := ListCidrMaps(gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &CidrMap{}, testCidrMap[0])
	assert.Equal(t, GtmTestCidrMap, testCidrMap[0].Name)

}

// Verify GetCidrMap. Name hardcoded. Should pass, e.g. no API errors and resource returned
func TestGetCidrMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps/" + GtmTestCidrMap)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps/"+GtmTestCidrMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                               "assignments": [ {
                                       "blocks": [ "1.2.3.0/24" ],
                                       "datacenterId": 3134,
                                       "nickname": "Frostfangs and the Fist of First Men"
                               },
                               {
                                       "blocks": [ "1.2.4.0/24" ],
                                       "datacenterId": 3133,
                                       "nickname": "Winterfell"
                               } ],
                               "defaultDatacenter": {
                                       "datacenterId": 5400,
                                       "nickname": "All Other CIDR Blocks"
                               },
                               "links": [ {
                                       "href": "/config-gtm/v1/domains/example.akadns.net/cidr-maps/testCidrMap",
                                       "rel": "self"
                               } ],
                               "name": "testCidrMap"
               }`)

	Init(config)

	testCidrMap, err := GetCidrMap(GtmTestCidrMap, gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &CidrMap{}, testCidrMap)
	assert.Equal(t, GtmTestCidrMap, testCidrMap.Name)

}

// Verify failed case for GetCidrMap. Should pass, e.g. no API errors and domain not found
func TestGetBadCidrMap(t *testing.T) {

	badName := "somebadname"
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps/" + badName)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps/"+badName).
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                }`)

	Init(config)

	_, err := GetCidrMap(badName, gtmTestDomain)
	assert.Error(t, err)

}

// Test Create CidrMap.
func TestCreateCidrMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps/" + GtmTestCidrMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps/"+GtmTestCidrMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                               "assignments": [ {
                                       "blocks": [ "1.2.3.0/24" ],
                                       "datacenterId": 3134,
                                       "nickname": "Frostfangs and the Fist of First Men"
                               },
                               {
                                       "blocks": [ "1.2.4.0/24" ],
                                       "datacenterId": 3133,
                                       "nickname": "Winterfell"
                               } ],
                               "defaultDatacenter": {
                                       "datacenterId": 5400,
                                       "nickname": "All Other CIDR Blocks"
                               },
                               "links": [ {
                                       "href": "/config-gtm/v1/domains/example.akadns.net/cidr-maps/testCidrMap",
                                       "rel": "self"
                               } ],
                               "name": "testCidrMap"
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	testCidrMap := instantiateCidrMap()
	statresp, err := testCidrMap.Create(gtmTestDomain)
	assert.NoError(t, err)

	assert.IsType(t, &CidrMap{}, statresp.Resource)
	assert.Equal(t, GtmTestCidrMap, statresp.Resource.Name)

}

func TestUpdateCidrMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps/" + GtmTestCidrMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps/"+GtmTestCidrMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                               "assignments": [ {
                                       "blocks": [ "1.2.3.0/24" ],
                                       "datacenterId": 3134,
                                       "nickname": "Frostfangs and the Fist of First Men"
                               },
                               {
                                       "blocks": [ "1.2.4.0/24" ],
                                       "datacenterId": 3133,
                                       "nickname": "Winterfell"
                               } ],
                               "defaultDatacenter": {
                                       "datacenterId": 5400,
                                       "nickname": "All Other CIDR Blocks"
                               },
                               "links": [ {
                                       "href": "/config-gtm/v1/domains/example.akadns.net/cidr-maps/testCidrMap",
                                       "rel": "self"
                               } ],
                               "name": "testCidrMap"
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	testCidrMap := instantiateCidrMap()
	_, err := testCidrMap.Update(gtmTestDomain)
	assert.NoError(t, err)

}

func TestDeleteCidrMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/cidr-maps/" + GtmTestCidrMap)
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/cidr-maps/"+GtmTestCidrMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	getCidrMap := instantiateCidrMap()
	stat, err := getCidrMap.Delete(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, "93a48b86-4fc3-4a5f-9ca2-036835034cc6", stat.ChangeId)

}

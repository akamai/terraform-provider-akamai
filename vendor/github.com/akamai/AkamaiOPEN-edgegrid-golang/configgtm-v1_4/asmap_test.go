package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        false,
	}
)

var GtmTestAsMap = "testAsMap"
var gtmTestDomain = "gtmdomtest.akadns.net"

func instantiateAsMap() *AsMap {

	asMap := NewAsMap(GtmTestAsMap)
	asMapData := []byte(`{
                        "assignments": [ {
                                        "asNumbers": [ 12222, 16702, 17334 ],
                                        "datacenterId": 3134,
                                        "nickname": "Frostfangs and the Fist of First Men"
                                }, {
                                        "asNumbers": [ 16625 ],
                                        "datacenterId": 3133,
                                        "nickname": "Winterfell"
                                } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "All Other AS numbers"
                        },
                        "name": "testAsMap"
              }`)
	jsonhooks.Unmarshal(asMapData, asMap)

	return asMap

}

// Verify GetAsMap. Name hardcoded. Should pass, e.g. no API errors and resource returned
// Depends on CreateAsMap
func TestGetAsMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/as-maps/" + GtmTestAsMap)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/as-maps/"+GtmTestAsMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "assignments": [ {
                                        "asNumbers": [ 12222, 16702, 17334 ],
                                        "datacenterId": 3134,
                                        "nickname": "Frostfangs and the Fist of First Men"
                                }, {
                                        "asNumbers": [ 16625 ],
                                        "datacenterId": 3133,
                                        "nickname": "Winterfell"
                                } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "All Other AS numbers"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
                                "rel": "self"
                        } ], 
                        "name": "testAsMap"
               }`)

	Init(config)

	testAsMap, err := GetAsMap(GtmTestAsMap, gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &AsMap{}, testAsMap)
	assert.Equal(t, GtmTestAsMap, testAsMap.Name)

}

// Verify failed case for GetAsMap. Should pass, e.g. no API errors and domain not found
func TestGetBadAsMap(t *testing.T) {

	badName := "somebadname"
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/as-maps/" + badName)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/as-maps/"+badName).
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                }`)

	Init(config)

	_, err := GetAsMap(badName, gtmTestDomain)
	assert.Error(t, err)

}

// Test Create AsMap.
func TestCreateAsMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/as-maps/" + GtmTestAsMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/as-maps/"+GtmTestAsMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                        "assignments": [ {
                                        "asNumbers": [ 12222, 16702, 17334 ],
                                        "datacenterId": 3134,
                                        "nickname": "Frostfangs and the Fist of First Men"
                                }, {
                                        "asNumbers": [ 16625 ],
                                        "datacenterId": 3133,
                                        "nickname": "Winterfell"
                                } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "All Other AS numbers"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
                                "rel": "self"
                        } ],
                        "name": "testAsMap"
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

	testAsMap := instantiateAsMap()
	statresp, err := testAsMap.Create(gtmTestDomain)
	assert.NoError(t, err)

	assert.IsType(t, &AsMap{}, statresp.Resource)
	assert.Equal(t, GtmTestAsMap, statresp.Resource.Name)

}

func TestUpdateAsMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/as-maps/" + GtmTestAsMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/as-maps/"+GtmTestAsMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                        "assignments": [ {
                                        "asNumbers": [ 12222, 16702, 17334 ],
                                        "datacenterId": 3134,
                                        "nickname": "Frostfangs and the Fist of First Men"
                                }, {
                                        "asNumbers": [ 16625 ],
                                        "datacenterId": 3133,
                                        "nickname": "Winterfell"
                                } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "All Other AS numbers"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/as-maps/The%20North",
                                "rel": "self"
                        } ],
                        "name": "testAsMap"
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

	testAsMap := instantiateAsMap()
	_, err := testAsMap.Update(gtmTestDomain)
	assert.NoError(t, err)

}

func TestDeleteAsMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/as-maps/" + GtmTestAsMap)
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/as-maps/"+GtmTestAsMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
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

	getAsMap := instantiateAsMap()
	stat, err := getAsMap.Delete(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, "93a48b86-4fc3-4a5f-9ca2-036835034cc6", stat.ChangeId)

}

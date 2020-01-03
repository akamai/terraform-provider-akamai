package reportsgtm

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
)

//
// Important note: The test cases enclosed piggyback on the objects created in the configgtm test cases
//
// TODO: Add tests for Opt args
//

//
var GtmTestProperty = "testproperty"

// Verify GetIpStatusPerProperty. Property and domain names hardcoded. Should pass, e.g. no API errors and property returned
// Depends on CreateProperty
func TestGetIpStatusProperty(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/ip-availability/domains/gtmtest.akadns.net/properties/" + GtmTestProperty)
	mock.
		Get("/gtm-api/v1/reports/ip-availability/domains/gtmtest.akadns.net/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                        "metadata": {
                                "domain": "gtmtest.akadns.net",
                                "property": "testproperty",
                                "start" : "2017-02-23T21:00:00Z",
                                "end" : "2017-03-23T22:00:00Z",
                                "uri": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/ip-availability/domains/gtmtest.akadns.net/properties/testproperty"
                        },
                        "dataRows": [
                                {
                                        "timestamp": "2017-02-23T21:42:35Z",
                                        "cutOff": 112.5,
                                        "datacenters": [
                                                {
                                                        "datacenterId": 3132,
                                                        "nickname": "Winterfell",
                                                        "trafficTargetName": "Winterfell - 1.2.3.4",
                                                        "IPs": [
                                                                 {
                                                                         "ip": "1.2.3.4",
                                                                         "score": 75.0,
                                                                         "handedOut": true,
                                                                         "alive": true
                                                                 } ]
                                                },
                                                {
                                                        "datacenterId": 3133,
                                                        "nickname": "Braavos",
                                                        "trafficTargetName": "Braavos - 1.2.3.5",
                                                        "IPs": [
                                                                 {
                                                                         "ip": "1.2.3.5",
                                                                         "score": 85.0,
                                                                         "handedOut": true,
                                                                         "alive": true
                                                                 } ]
                                                } ]
                                },
                                {
                                        "timestamp": "2017-03-23T21:42:35Z",
                                        "cutOff": 112.5,
                                        "datacenters": [
                                                {
                                                        "datacenterId": 3132,
                                                        "nickname": "Winterfell",
                                                        "trafficTargetName": "Winterfell - 1.2.3.4",
                                                        "IPs": [
                                                                 {
                                                                         "ip": "1.2.3.4",
                                                                         "score": 115.0,
                                                                         "handedOut": false,
                                                                         "alive": false
                                                                 } ]
                                                },               
                                                {
                                                        "datacenterId": 3133,
                                                        "nickname": "Braavos",
                                                        "trafficTargetName": "Braavos - 1.2.3.5",
                                                        "IPs": [
                                                                 {
                                                                         "ip": "1.2.3.5",
                                                                         "score": 75.0,
                                                                         "handedOut": true,
                                                                         "alive": true
                                                                 } ]
                                                } ]              
                                } ],    
                        "links": [
                                {
                                        "rel": "self",
                                        "href": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/ip-availability/domains/gtmtest.akadns.net/properties/testproperty"
                                } ]
                 }`)

	Init(config)

	optArgs := make(map[string]string)
	optArgs["start"] = "2017-02-23T21:00:00Z"
	optArgs["end"] = "2017-03-23T22:00:00Z"
	testPropertyIpStatus, err := GetIpStatusPerProperty(gtmTestDomain, GtmTestProperty, optArgs)

	assert.NoError(t, err)
	assert.Equal(t, testPropertyIpStatus.DataRows[0].Datacenters[0].DatacenterId, 3132)
	assert.Equal(t, testPropertyIpStatus.Metadata.Domain, gtmTestDomain)
	assert.Equal(t, testPropertyIpStatus.Metadata.Property, GtmTestProperty)

}

/*
// TestGetIpStatus with mostRecent flag
func TestGetIpStatusPropertyRecent(t *testing.T) {

       config, _ := edgegrid.InitEdgeRc("", "default")
       Init(config)

       // add mock ...

       optArgs := make(map[string]string)
       optArgs["mostRecent"] = "true"

       domDCs, err := configgtm.ListDatacenters(GtmObjectTestDomain)
       assert.NoError(t, err, "Failure retrieving DCs")
       if len(domDCs) > 0 {
               optArgs["datacenterId"] = strconv.Itoa(domDCs[0].DatacenterId)
               fmt.Println("dcid: "+optArgs["datacenterId"])
       }

       testPropertyIpStatus, err := GetIpStatusPerProperty(GtmObjectTestDomain, GtmTestProperty, optArgs)

       assert.NoError(t, err)

       json, err := json.MarshalIndent(testPropertyIpStatus, "", "    ")
       if err == nil {
               fmt.Println(string(json))
       } else {
               t.Fatal("PropertyIP Status retrival failed. " + err.Error())
       }
}
*/

// Verify GetTrafficPerProperty. Domain name and property name hardcoded.

func TestGetTrafficPerProperty(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/" + GtmTestProperty)
	mock.
		Get("/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/"+GtmTestProperty).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                        "metadata": {
                                "domain": "gtmtest.akadns.net",
                                "property": "testproperty",
                                "start": "2016-11-24T01:40:00Z",
                                "end": "2016-11-24T01:50:00Z",
                                "interval": "FIVE_MINUTE",
                                "uri": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/testproperty?start=2016-11-23T00:00:00Z&2016-11-24T01:50:00Z"
                        },
                        "dataRows": [
                                {
                                        "timestamp": "2016-11-24T01:40:00Z",
                                        "datacenters": [
                                                {
                                                        "datacenterId": 3130,
                                                        "nickname": "Winterfell",
                                                        "trafficTargetName": "Winterfell - 1.2.3.4",
                                                        "requests": 34,
                                                        "status": "1"
                                                } ]
                                },
                                {
                                        "timestamp": "2016-11-24T01:45:00Z",
                                        "datacenters": [
                                                {
                                                        "datacenterId": 3130,
                                                        "nickname": "Winterfell",
                                                        "trafficTargetName": "Winterfell - 1.2.3.4",
                                                        "requests": 45,
                                                        "status": "1"
                                                } ]
                                } ],
                        "links": [
                                {
                                        "rel": "self",
                                        "href": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/testproperty?start=2016-11-23T00:00:00Z&2016-11-24T01:50:00Z"
                                } ]
                 }`)

	Init(config)

	optArgs := make(map[string]string)
	optArgs["start"] = "2016-11-24T01:40:00Z"
	optArgs["end"] = "2016-11-24T01:50:00Z"

	testPropertyTraffic, err := GetTrafficPerProperty(gtmTestDomain, GtmTestProperty, optArgs)

	assert.NoError(t, err)
	assert.Equal(t, testPropertyTraffic.DataRows[0].Datacenters[0].DatacenterId, 3130)
	assert.Equal(t, testPropertyTraffic.Metadata.Domain, gtmTestDomain)
	assert.Equal(t, testPropertyTraffic.Metadata.Property, GtmTestProperty)

}

// Verify failed case for GetProperty. Should pass, e.g. no API errors and domain not found
func TestGetBadTrafficPerProperty(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/badproperty")
	mock.
		Get("/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/properties/badproperty").
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                }`)

	Init(config)

	optArgs := make(map[string]string)
	_, err := GetTrafficPerProperty(gtmTestDomain, "badproperty", optArgs)

	assert.Error(t, err)

}

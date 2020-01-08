package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"fmt"
)

var GtmTestResource = "testResource"

func instantiateResource() *Resource {

	resource := NewResource("dummy")
	resourceData := []byte(`{
                                "aggregationType" : "median",
                                "constrainedProperty" : null,
                                "decayRate" : null,
                                "description" : null,
                                "hostHeader" : null,
                                "leaderString" : null,
                                "leastSquaresDecay" : null,
                                "loadImbalancePercentage" : null,
                                "maxUMultiplicativeIncrement" : null,
                                "name" : "testResource",
                                "resourceInstances" : [ {
                                        "loadObject" : "",
                                        "loadObjectPort" : 0,
                                        "loadServers" : null,
                                        "datacenterId" : 3131,
                                        "useDefaultLoadObject" : false
                                } ],
                                "type" : "Download score",
                                "upperBound" : 0,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources/testResource"
                                } ]
                       }`)
	jsonhooks.Unmarshal(resourceData, resource)

	return resource

}

// Verify ListResources. Should pass, e.g. no API errors and non nil list.
func TestListResources(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/resources").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "items" : [ {
                                "aggregationType" : "median",
                                "constrainedProperty" : null,
                                "decayRate" : null,
                                "description" : null,
                                "hostHeader" : null,
                                "leaderString" : null,
                                "leastSquaresDecay" : null,
                                "loadImbalancePercentage" : null,
                                "maxUMultiplicativeIncrement" : null,
                                "name" : "testResource",
                                "resourceInstances" : [ ],
                                "type" : "Download score",
                                "upperBound" : 0,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources/testResource"
                                } ]
                        } ]
               }`)

	Init(config)
	resourceList, err := ListResources(gtmTestDomain)
	assert.NoError(t, err)
	assert.NotEqual(t, resourceList, nil)

	if len(resourceList) > 0 {
		firstResource := resourceList[0]
		assert.Equal(t, firstResource.Name, GtmTestResource)
	} else {
		t.Fatal("List empty!")
	}
}

// Depends on CreateResource
func TestGetResource(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources/" + GtmTestResource)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/resources/"+GtmTestResource).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "aggregationType" : "median",
                        "constrainedProperty" : null,
                        "decayRate" : null,
                        "description" : null,
                        "hostHeader" : null,
                        "leaderString" : null,
                        "leastSquaresDecay" : null,
                        "loadImbalancePercentage" : null,
                        "maxUMultiplicativeIncrement" : null,
                        "name" : "testResource",
                        "resourceInstances" : [ ],
                        "type" : "Download score",
                        "upperBound" : 0,
                        "links" : [ {
                                "rel" : "self",
                                "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources/testResource"
                        } ]
              }`)

	Init(config)

	testResource, err := GetResource(GtmTestResource, gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &Resource{}, testResource)
	assert.Equal(t, GtmTestResource, testResource.Name)

}

// Verify failed case for GetResource. Should pass, e.g. no API errors and domain not found
func TestGetBadResource(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources/somebadname")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/resources/somebadname").
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        }`)

	Init(config)

	_, err := GetResource("somebadname", gtmTestDomain)
	// Shouldn't have found
	assert.Error(t, err)

}

// Test Create resource.
func TestCreateResource(t *testing.T) {

	defer gock.Off()
	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources/" + GtmTestResource)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/resources/"+GtmTestResource).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "aggregationType" : "median",
                                "constrainedProperty" : null,
                                "decayRate" : null,
                                "description" : null,
                                "hostHeader" : null,
                                "leaderString" : null,
                                "leastSquaresDecay" : null,
                                "loadImbalancePercentage" : null,
                                "maxUMultiplicativeIncrement" : null,
                                "name" : "testResource",
                                "resourceInstances" : [ {
                                        "loadObject" : "",
                                        "loadObjectPort" : 0,
                                        "loadServers" : null,
                                        "datacenterId" : 3131,
                                        "useDefaultLoadObject" : false
                                } ],
                                "type" : "Download score",
                                "upperBound" : 0,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources/testResource"
                                } ]
                        },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "5bb9f131-99c8-43ff-afd2-a6ce34db8b95",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-17T17:54:37.383+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
                }`)

	Init(config)

	testResource := NewResource(GtmTestResource)
	testResource.AggregationType = "median"
	testResource.Type = "Download score"

	// Create a Resource Instance
	var instanceSlice []*ResourceInstance
	instance := testResource.NewResourceInstance(3131)
	instanceSlice = append(instanceSlice, instance)
	testResource.ResourceInstances = instanceSlice

	// do the create

	fmt.Println("Calling Create!!!")
	statresp, err := testResource.Create(gtmTestDomain)
	fmt.Println("Returned")
	assert.NoError(t, err)

	assert.IsType(t, &Resource{}, statresp.Resource)
	assert.Equal(t, GtmTestResource, statresp.Resource.Name)

}

func TestUpdateResource(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources/" + GtmTestResource)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/resources/"+GtmTestResource).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "aggregationType" : "median",
                                "constrainedProperty" : null,
                                "decayRate" : 0.5,
                                "description" : null,
                                "hostHeader" : null,
                                "leaderString" : null,
                                "leastSquaresDecay" : null,
                                "loadImbalancePercentage" : null,
                                "maxUMultiplicativeIncrement" : null,
                                "name" : "testResource",
                                "resourceInstances" : [ {
                                        "loadObject" : "",
                                        "loadObjectPort" : 0,
                                        "loadServers" : null,
                                        "datacenterId" : 3131,
                                        "useDefaultLoadObject" : false
                                } ],
                                "type" : "Download score",
                                "upperBound" : 0,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources/testResource"
                                } ]
                        },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "5bb9f131-99c8-43ff-afd2-a6ce34db8b95",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-17T17:54:37.383+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
                }`)

	Init(config)

	testResource := instantiateResource()

	newDecay := 0.5
	testResource.DecayRate = newDecay
	stat, err := testResource.Update(gtmTestDomain)
	assert.NoError(t, err)

	assert.Equal(t, stat.ChangeId, "5bb9f131-99c8-43ff-afd2-a6ce34db8b95")

}

func TestDeleteResource(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/resources/" + GtmTestResource)
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/resources/"+GtmTestResource).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "resource" : null,
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "9a7a8f84-704b-40de-a903-bcc2728513ac",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-06-14T19:51:02.273+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                 } ]
                        }
               }`)

	Init(config)

	getResource := instantiateResource()
	stat, err := getResource.Delete(gtmTestDomain)
	assert.NoError(t, err)

	assert.Equal(t, "9a7a8f84-704b-40de-a903-bcc2728513ac", stat.ChangeId)

}

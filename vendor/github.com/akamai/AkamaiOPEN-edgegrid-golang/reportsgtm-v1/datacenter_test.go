package reportsgtm

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
	"github.com/stretchr/testify/assert"
)

// Verify GetTrafficPerDatacenter.
func TestGetTrafficPerDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters/3200")
	mock.
		Get("/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters/3200").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                        "metadata": {
                                "domain": "gtmtest.akadns.net",
                                "datacenterId": 3200,
                                "datacenterNickname": "Winterfell",
                                "start": "2016-11-23T00:00:00Z",
                                "end": "2016-11-23T00:10:00Z",
                                "interval": "FIVE_MINUTE",
                                "uri": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters?start=2016-11-23T00:00:00Z&end=2016-11-23T00:10:00Z"
                        },
                        "dataRows": [ {
                                        "timestamp": "2016-11-23T00:00:00Z",
                                        "properties": [ {
                                                "name": "www",
                                                "requests": 45,
                                                "status": "1"
                                        } ]
                                 },
                                 {
                                        "timestamp": "2016-11-23T00:05:00Z",
                                        "properties": [ {
                                                "name": "www",
                                                "requests": 45,
                                                "status": "1"
                                        } ]
                                 } ],
                        "links": [ {
                                         "rel": "self",
                                         "href": "https://akaa-xxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters?start=2016-11-23T00:00:00Z&2016-11-23T00:10:00Z"
                                 } ]
                           }`)

	Init(config)

	optArgs := make(map[string]string)
	optArgs["start"] = "2016-11-23T00:00:00Z"
	optArgs["end"] = "2016-11-23T00:10:00Z"

	testDCTraffic, err := GetTrafficPerDatacenter(gtmTestDomain, 3200, optArgs)
	assert.NoError(t, err)
	assert.Equal(t, "gtmtest.akadns.net", testDCTraffic.Metadata.Domain)
	assert.Equal(t, testDCTraffic.DataRows[0].Timestamp, "2016-11-23T00:00:00Z")

}

// Verify failed case for TrafficPerDatacenter.
func TestGetBadTrafficPerDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters/9999")
	mock.
		Get("/gtm-api/v1/reports/traffic/domains/gtmtest.akadns.net/datacenters/9999").
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                }`)

	Init(config)

	optArgs := make(map[string]string)
	_, err := GetTrafficPerDatacenter(gtmTestDomain, 9999, optArgs)

	assert.Error(t, err)
}

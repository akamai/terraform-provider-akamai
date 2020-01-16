package reportsgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"

	"fmt"
)

//
// Support gtm reports thru Edgegrid
// Based on 1.0 Schema
//

// Property Traffic Report Structs
type PropertyTMeta struct {
	Uri      string `json:uri"`
	Domain   string `json:"domain"`
	Interval string `json:"interval,omitempty"`
	Property string `json:"property"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

type PropertyDRow struct {
	Nickname          string `json:"nickname"`
	DatacenterId      int    `json:"datacenterId"`
	TrafficTargetName string `json:"trafficTargetName"`
	Requests          int64  `json:"requests"`
	Status            string `json:"status"`
}

type PropertyTData struct {
	Timestamp   string          `json:"timestamp"`
	Datacenters []*PropertyDRow `json:"datacenters"`
}

// The Property Traffic Response structure returned by the Reports API
type PropertyTrafficResponse struct {
	Metadata    *PropertyTMeta    `json:"metadata"`
	DataRows    []*PropertyTData  `json:"dataRows"`
	DataSummary interface{}       `json:"dataSummary"`
	Links       []*configgtm.Link `json:"links"`
}

//
// IP Status By Property Structs
//

// IP Availability Status Response structure returned by the Reports API.
type IPStatusPerProperty struct {
	Metadata    *IpStatPerPropMeta   `json:"metadata"`
	DataRows    []*IpStatPerPropData `json:"dataRows"`
	DataSummary interface{}          `json:"dataSummary"`
	Links       []*configgtm.Link    `json:"links"`
}

type IpStatPerPropMeta struct {
	Uri          string `json:uri"`
	Domain       string `json:"domain"`
	Property     string `json:"property"`
	Start        string `json:"start"`
	End          string `json:"end"`
	MostRecent   bool   `json:"mostRecent"`
	Ip           string `json:"ip"`
	DatacenterId int    `json:"datacenterId"`
}

type IpStatPerPropData struct {
	Timestamp   string               `json:"timestamp"`
	CutOff      float64              `json:"cutOff"`
	Datacenters []*IpStatPerPropDRow `json:"datacenters"`
}

type IpStatPerPropDRow struct {
	Nickname          string      `json:"nickname"`
	DatacenterId      int         `json:"datacenterId"`
	TrafficTargetName string      `json:"trafficTargetName"`
	IPs               []*IpStatIp `json:"IPs"`
}

type IpStatIp struct {
	Ip        string  `json:"ip"`
	HandedOut bool    `json:"handedOut"`
	Score     float32 `json:"score"`
	Alive     bool    `json:"alive"`
}

// GetIpStatusPerProperty retrieves current IP Availability Status for specified property in the given domainname.
func GetIpStatusPerProperty(domainName string, propertyName string, optArgs map[string]string) (*IPStatusPerProperty, error) {
	stat := &IPStatusPerProperty{}
	hostURL := fmt.Sprintf("/gtm-api/v1/reports/ip-availability/domains/%s/properties/%s", domainName, propertyName)

	req, err := client.NewRequest(
		Config,
		"GET",
		hostURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Look for and process optional query params
	q := req.URL.Query()
	for k, v := range optArgs {
		switch k {
		case "start":
			fallthrough
		case "end":
			fallthrough
		case "ip":
			fallthrough
		case "mostRecent":
			fallthrough
		case "datacenterId":
			q.Add(k, v)
		}
	}
	if optArgs != nil {
		req.URL.RawQuery = q.Encode()
	}

	// time stamps require urlencoded content header
	setEncodedHeader(req)

	// print/log the request if warranted
	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	// print/log the response if warranted
	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		cErr := configgtm.CommonError{}
		cErr.SetItem("entityName", "Property")
		cErr.SetItem("name", propertyName)
		return nil, cErr
	} else {
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}

		return stat, nil
	}
}

// GetTrafficPerProperty retrieves report traffic for the specified property in the specified domain.
func GetTrafficPerProperty(domainName string, propertyName string, optArgs map[string]string) (*PropertyTrafficResponse, error) {
	stat := &PropertyTrafficResponse{}
	hostURL := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/properties/%s", domainName, propertyName)

	req, err := client.NewRequest(
		Config,
		"GET",
		hostURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Look for and process optional query params
	q := req.URL.Query()
	for k, v := range optArgs {
		switch k {
		case "start":
			fallthrough
		case "end":
			q.Add(k, v)
		}
	}
	if optArgs != nil {
		req.URL.RawQuery = q.Encode()
	}

	// time stamps require urlencoded content header
	setEncodedHeader(req)

	// print/log the request if warranted
	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	// print/log the response if warranted
	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		cErr := configgtm.CommonError{}
		cErr.SetItem("entityName", "Property")
		cErr.SetItem("name", propertyName)
		return nil, cErr
	} else {
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}

		return stat, nil
	}
}

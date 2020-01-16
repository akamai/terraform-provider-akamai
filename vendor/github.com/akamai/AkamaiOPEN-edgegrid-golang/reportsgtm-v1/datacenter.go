package reportsgtm

import (
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"

	"fmt"
)

//
// Support gtm reports thru Edgegrid
// Based on 1.0 Schema
//

// Datacenter Traffic Report Structs
type DCTMeta struct {
	Uri                string `json:uri"`
	Domain             string `json:"domain"`
	Interval           string `json:"interval,omitempty"`
	DatacenterId       int    `json:"datacenterId"`
	DatacenterNickname string `json:"datacenterNickname"`
	Start              string `json:"start"`
	End                string `json:"end"`
}

type DCTDRow struct {
	Name     string `json:"name"`
	Requests int64  `json:"requests"`
	Status   string `json:"status"`
}

type DCTData struct {
	Timestamp  string     `json:"timestamp"`
	Properties []*DCTDRow `json:"properties"`
}

// Response structure returned by the Datacenter Traffic API
type DcTrafficResponse struct {
	Metadata    *DCTMeta          `json:"metadata"`
	DataRows    []*DCTData        `json:"dataRows"`
	DataSummary interface{}       `json:"dataSummary"`
	Links       []*configgtm.Link `json:"links"`
}

// GetTrafficPerDatacenter retrieves Report Traffic per datacenter. Opt args - start, end.
func GetTrafficPerDatacenter(domainName string, datacenterID int, optArgs map[string]string) (*DcTrafficResponse, error) {
	stat := &DcTrafficResponse{}
	hostURL := fmt.Sprintf("/gtm-api/v1/reports/traffic/domains/%s/datacenters/%s", domainName, strconv.Itoa(datacenterID))

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
		// time stamps require urlencoded content header
		setEncodedHeader(req)
		req.URL.RawQuery = q.Encode()
	}

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
		cErr.SetItem("entityName", "Datacenter")
		cErr.SetItem("name", strconv.Itoa(datacenterID))
		return nil, cErr
	} else {
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}

		return stat, nil
	}
}

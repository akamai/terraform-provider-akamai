package reportsgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"net/http"

	"fmt"
	"time"
)

//
// Support gtm reports thru Edgegrid
// Based on 1.0 Schema
//

type WindowResponse struct {
	StartTime time.Time
	EndTime   time.Time
}

type APIWindowResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func setEncodedHeader(req *http.Request) {

	if req.Method == "GET" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return

}

// Utility method to convert an RFC3339 formated time string to a time.Time object
func convertRFC3339toDate(rfc3339Stamp string) (time.Time, error) {

	t, err := time.Parse(time.RFC3339, rfc3339Stamp)
	return t, err

}

func createTimeWindow(apiResponse *APIWindowResponse) (*WindowResponse, error) {

	var err error
	windowResponse := &WindowResponse{}
	windowResponse.StartTime, err = convertRFC3339toDate(apiResponse.Start)
	if err != nil {
		return nil, err
	}
	windowResponse.EndTime, err = convertRFC3339toDate(apiResponse.End)
	if err != nil {
		return nil, err
	}
	return windowResponse, nil

}

// Core function to retrieve all Window API requests
func getWindowCore(hostURL string) (*WindowResponse, error) {

	stat := &APIWindowResponse{}

	req, err := client.NewRequest(
		Config,
		"GET",
		hostURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) {
		if res.StatusCode == 400 {
                	// Get the body. Could be bad dates.
                	var windRespErrBody map[string]interface{}
                	err = client.BodyJSON(res, windRespErrBody)
			if err != nil {
				return nil, err
			}
			// are available dates present?
			if availEnd, ok := windRespErrBody["availableEndDate"]; ok {
				stat.End = availEnd.(string)
			}
                	if availStart, ok := windRespErrBody["availableStartDate"]; ok {
				stat.Start = availStart.(string)
			}
			if stat.End == "" || stat.Start == "" {
                		cErr := configgtm.CommonError{}
                		cErr.SetItem("entityName", "Window")
                		cErr.SetItem("name", "Data Window")
				cErr.SetItem("apiErrorMessage", "No available data window")
                		return nil, cErr	
			}	
		} else if res.StatusCode == 404 {
			cErr := configgtm.CommonError{}
			cErr.SetItem("entityName", "Window")
			cErr.SetItem("name", "Data Window")
			return nil, cErr
		} else {
			return nil, client.NewAPIError(res)
		}
	} else {
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}
	}
	timeWindow, err := createTimeWindow(stat)
	if err != nil {
		return nil, err
	}
	return timeWindow, nil
}

// GetDemandWindow is a utility function that retrieves the data window for Demand category of Report APIs
func GetDemandWindow(domainName string, propertyName string) (*WindowResponse, error) {

	hostURL := fmt.Sprintf("/gtm-api/v1/reports/demand/domains/%s/properties/%s/window", domainName, propertyName)
	return getWindowCore(hostURL)

}

// GetLatencyDomainsWindow is a utility function that retrieves the data window for Latency category of Report APIs
func GetLatencyDomainsWindow(domainName string) (*WindowResponse, error) {

	hostURL := fmt.Sprintf("/gtm-api/v1/reports/latency/domains/%s/window", domainName)
	return getWindowCore(hostURL)

}

// GetLivenessTestsWindow is a utility function that retrieves the data window for Liveness category of Report APIs
func GetLivenessTestsWindow() (*WindowResponse, error) {

	hostURL := fmt.Sprintf("/gtm-api/v1/reports/liveness-tests/window")
	return getWindowCore(hostURL)

}

// GetDatacentersTrafficWindow is a utility function that retrieves the data window for Traffic category of Report APIs
func GetDatacentersTrafficWindow() (*WindowResponse, error) {

	hostURL := fmt.Sprintf("/gtm-api/v1/reports/traffic/datacenters-window")
	return getWindowCore(hostURL)

}

// GetPropertiesTrafficWindow is a utility function that retrieves the data window for Traffic category of Report API
func GetPropertiesTrafficWindow() (*WindowResponse, error) {

	hostURL := fmt.Sprintf("/gtm-api/v1/reports/traffic/properties-window")
	return getWindowCore(hostURL)

}

package configgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"

	"fmt"
)

//
// Support gtm domain properties thru Edgegrid
// Based on 1.4 Schema
//

// TrafficTarget struc
type TrafficTarget struct {
	DatacenterId int      `json:"datacenterId"`
	Enabled      bool     `json:"enabled"`
	Weight       float64  `json:"weight,omitempty"`
	Servers      []string `json:"servers,omitempty"`
	Name         string   `json:"name,omitempty"`
	HandoutCName string   `json:"handoutCName,omitempty"`
}

// HttpHeader struc
type HttpHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LivenessTest struct {
	Name                          string        `json:"name"`
	ErrorPenalty                  float64       `json:"errorPenalty,omitempty"`
	PeerCertificateVerification   bool          `json:"peerCertificateVerification"`
	TestInterval                  int           `json:"testInterval,omitempty"`
	TestObject                    string        `json:"testObject,omitempty"`
	Links                         []*Link       `json:"links,omitempty"`
	RequestString                 string        `json:"requestString,omitempty"`
	ResponseString                string        `json:"responseString,omitempty"`
	HttpError3xx                  bool          `json:"httpError3xx"`
	HttpError4xx                  bool          `json:"httpError4xx"`
	HttpError5xx                  bool          `json:"httpError5xx"`
	Disabled                      bool          `json:"disabled"`
	TestObjectProtocol            string        `json:"testObjectProtocol,omitempty"`
	TestObjectPassword            string        `json:"testObjectPassword,omitempty"`
	TestObjectPort                int           `json:"testObjectPort,omitempty"`
	SslClientPrivateKey           string        `json:"sslClientPrivateKey,omitempty"`
	SslClientCertificate          string        `json:"sslClientCertificate,omitempty"`
	DisableNonstandardPortWarning bool          `json:"disableNonstandardPortWarning"`
	HttpHeaders                   []*HttpHeader `json:"httpHeaders,omitempty"`
	TestObjectUsername            string        `json:"testObjectUsername,omitempty"`
	TestTimeout                   float32       `json:"testTimeout,omitempty"`
	TimeoutPenalty                float64       `json:"timeoutPenalty,omitempty"`
	AnswersRequired               bool          `json:"answersRequired"`
	ResourceType                  string        `json:"resourceType,omitempty"`
	RecursionRequested            bool          `json:"recursionRequested"`
}

// StaticRRSet Struct
type StaticRRSet struct {
	Type  string   `json:"type"`
	TTL   int      `json:"ttl"`
	Rdata []string `json:"rdata"`
}

// Property represents a GTM property
type Property struct {
	Name                      string           `json:"name"`
	Type                      string           `json:"type"`
	Ipv6                      bool             `json:"ipv6"`
	ScoreAggregationType      string           `json:"scoreAggregationType"`
	StickinessBonusPercentage int              `json:"stickinessBonusPercentage,omitempty"`
	StickinessBonusConstant   int              `json:"stickinessBonusConstant,omitempty"`
	HealthThreshold           float64          `json:"healthThreshold,omitempty"`
	UseComputedTargets        bool             `json:"useComputedTargets"`
	BackupIp                  string           `json:"backupIp,omitempty"`
	BalanceByDownloadScore    bool             `json:"balanceByDownloadScore"`
	StaticTTL                 int              `json:"staticTTL,omitempty"`
	StaticRRSets              []*StaticRRSet   `json:"staticRRSets,omitempty"`
	LastModified              string           `json:"lastModified"`
	UnreachableThreshold      float64          `json:"unreachableThreshold,omitempty"`
	MinLiveFraction           float64          `json:"minLiveFraction,omitempty"`
	HealthMultiplier          float64          `json:"healthMultiplier,omitempty"`
	DynamicTTL                int              `json:"dynamicTTL,omitempty"`
	MaxUnreachablePenalty     int              `json:"maxUnreachablePenalty,omitempty"`
	MapName                   string           `json:"mapName,omitempty"`
	HandoutLimit              int              `json:"handoutLimit"`
	HandoutMode               string           `json:"handoutMode"`
	FailoverDelay             int              `json:"failoverDelay,omitempty"`
	BackupCName               string           `json:"backupCName,omitempty"`
	FailbackDelay             int              `json:"failbackDelay,omitempty"`
	LoadImbalancePercentage   float64          `json:"loadImbalancePercentage,omitempty"`
	HealthMax                 float64          `json:"healthMax,omitempty"`
	GhostDemandReporting      bool             `json:"ghostDemandReporting"`
	Comments                  string           `json:"comments,omitempty"`
	CName                     string           `json:"cname,omitempty"`
	WeightedHashBitsForIPv4   int              `json:"weightedHashBitsForIPv4,omitempty"`
	WeightedHashBitsForIPv6   int              `json:"weightedHashBitsForIPv6,omitempty"`
	TrafficTargets            []*TrafficTarget `json:"trafficTargets,omitempty"`
	Links                     []*Link          `json:"links,omitempty"`
	LivenessTests             []*LivenessTest  `json:"livenessTests,omitempty"`
}

type PropertyList struct {
	PropertyItems []*Property `json:"items"`
}

// NewTrafficTarget is a method applied to a property object that instantiates a TrafficTarget object.
func (prop *Property) NewTrafficTarget() *TrafficTarget {

	return &TrafficTarget{}

}

// NewStaticRRSet is a method applied to a property object that instantiates a StaticRRSet object.
func (prop *Property) NewStaticRRSet() *StaticRRSet {

	return &StaticRRSet{}

}

// NewHttpHeader is a method applied to a livenesstest object that instantiates an HttpHeader  object.
func (lt *LivenessTest) NewHttpHeader() *HttpHeader {

	return &HttpHeader{}

}

// NewLivenessTest is a method applied to a property object that instantiates a LivenessTest object.
func (prop *Property) NewLivenessTest(name string, objproto string, interval int, timeout float32) *LivenessTest {

	return &LivenessTest{Name: name, TestInterval: interval, TestObjectProtocol: objproto, TestTimeout: timeout}

}

// NewProperty creates a new Property object.
func NewProperty(name string) *Property {
	property := &Property{Name: name}
	return property
}

// ListProperties retreieves all Properties for the provided domainName.
func ListProperties(domainName string) ([]*Property, error) {
	properties := &PropertyList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/properties", domainName),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Domain", name: domainName}
	} else {
		err = client.BodyJSON(res, properties)
		if err != nil {
			return nil, err
		}

		return properties.PropertyItems, nil
	}
}

// GetProperty retrieves a Property with the given name.
func GetProperty(name, domainName string) (*Property, error) {
	property := NewProperty(name)
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, name),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	printHttpResponse(res, true)

	if client.IsError(res) && res.StatusCode != 404 {
		return nil, client.NewAPIError(res)
	} else if res.StatusCode == 404 {
		return nil, CommonError{entityName: "Property", name: name}
	} else {
		err = client.BodyJSON(res, property)
		if err != nil {
			return nil, err
		}

		return property, nil
	}
}

// Create the property in the receiver argument in the specified domain.
func (property *Property) Create(domainName string) (*PropertyResponse, error) {

	// Need do any validation?
	return property.save(domainName)
}

// Update the property in the receiver argument in the specified domain.
func (property *Property) Update(domainName string) (*ResponseStatus, error) {

	// Need do any validation?
	stat, err := property.save(domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save Property updates method
func (property *Property) save(domainName string) (*PropertyResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, property.Name),
		property,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Property",
			name:             property.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Property", name: property.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &PropertyResponse{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil

}

// Delete the property identified by the receiver argument from the domain provided.
func (property *Property) Delete(domainName string) (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, property.Name),
		nil,
	)
	if err != nil {
		return nil, err
	}

	setVersionHeader(req, schemaVersion)

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Property",
			name:             property.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Property", name: property.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Status, nil

}

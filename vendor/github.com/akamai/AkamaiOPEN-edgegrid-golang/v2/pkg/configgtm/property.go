package gtm

import (
	"context"
	"fmt"
	"net/http"
)

//
// Support gtm domain properties thru Edgegrid
// Based on 1.4 Schema
//

// Propertiess contains operations available on a Property resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type Properties interface {
	// NewTrafficTarget is a method applied to a property object that instantiates a TrafficTarget object.
	NewTrafficTarget(context.Context) *TrafficTarget
	// NewStaticRRSet is a method applied to a property object that instantiates a StaticRRSet object.
	NewStaticRRSet(context.Context) *StaticRRSet
	// NewLivenessTest is a method applied to a property object that instantiates a LivenessTest object.
	NewLivenessTest(context.Context, string, string, int, float32) *LivenessTest
	// NewProperty creates a new Property object.
	NewProperty(context.Context, string) *Property
	// ListProperties retreieves all Properties for the provided domainName.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getproperties
	ListProperties(context.Context, string) ([]*Property, error)
	// GetProperty retrieves a Property with the given domain and property names.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getproperty
	GetProperty(context.Context, string, string) (*Property, error)
	// CreateProperty creates property
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putproperty
	CreateProperty(context.Context, *Property, string) (*PropertyResponse, error)
	// DeleteProperty is a method applied to a property object resulting in removal.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#deleteproperty
	DeleteProperty(context.Context, *Property, string) (*ResponseStatus, error)
	// UpdateProperty is a method applied to a property object resulting in an update.
	// https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#putproperty
	UpdateProperty(context.Context, *Property, string) (*ResponseStatus, error)
}

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

// Validate validates Property
func (prop *Property) Validate() error {

	if len(prop.Name) < 1 {
		return fmt.Errorf("Property is missing Name")
	}
	if len(prop.Type) < 1 {
		return fmt.Errorf("Property is missing Type")
	}
	if len(prop.ScoreAggregationType) < 1 {
		return fmt.Errorf("Property is missing ScoreAggregationType")
	}
	if len(prop.HandoutMode) < 1 {
		return fmt.Errorf("Property is missing HandoutMode")
	}
	// is zero a valid value? need to check and uncomment
	//if prop.HandoutLimit == 0 {
	//        return fmt.Errorf("Property is missing  handoutLimit"
	//}

	return nil
}

// NewTrafficTarget is a method applied to a property object that instantiates a TrafficTarget object.
func (p *gtm) NewTrafficTarget(ctx context.Context) *TrafficTarget {

	logger := p.Log(ctx)
	logger.Debug("NewTrafficTarget")

	return &TrafficTarget{}

}

// NewStaticRRSet is a method applied to a property object that instantiates a StaticRRSet object.
func (p *gtm) NewStaticRRSet(ctx context.Context) *StaticRRSet {

	logger := p.Log(ctx)
	logger.Debug("NewStaticRRSet")

	return &StaticRRSet{}

}

// NewHttpHeader is a method applied to a livenesstest object that instantiates an HttpHeader  object.
func (lt *LivenessTest) NewHttpHeader() *HttpHeader {

	return &HttpHeader{}

}

// NewLivenessTest is a method applied to a property object that instantiates a LivenessTest object.
func (p *gtm) NewLivenessTest(ctx context.Context, name string, objproto string, interval int, timeout float32) *LivenessTest {

	logger := p.Log(ctx)
	logger.Debug("NewLivenessTest")

	return &LivenessTest{Name: name, TestInterval: interval, TestObjectProtocol: objproto, TestTimeout: timeout}

}

// NewProperty creates a new Property object.
func (p *gtm) NewProperty(ctx context.Context, name string) *Property {

	logger := p.Log(ctx)
	logger.Debug("NewProperty")

	property := &Property{Name: name}
	return property
}

// ListProperties retreieves all Properties for the provided domainName.
func (p *gtm) ListProperties(ctx context.Context, domainName string) ([]*Property, error) {

	logger := p.Log(ctx)
	logger.Debug("ListProperties")

	var properties PropertyList
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListProperties request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &properties)
	if err != nil {
		return nil, fmt.Errorf("ListProperties request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return properties.PropertyItems, nil
}

// GetProperty retrieves a Property with the given name.
func (p *gtm) GetProperty(ctx context.Context, name, domainName string) (*Property, error) {

	logger := p.Log(ctx)
	logger.Debug("GetProperty")

	var property Property
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetProperty request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &property)
	if err != nil {
		return nil, fmt.Errorf("GetProperty request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &property, nil
}

// Create the property in the receiver argument in the specified domain.
func (p *gtm) CreateProperty(ctx context.Context, property *Property, domainName string) (*PropertyResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateProperty")

	// Need do any validation?
	return property.save(ctx, p, domainName)
}

// Update the property in the receiver argument in the specified domain.
func (p *gtm) UpdateProperty(ctx context.Context, property *Property, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateProperty")

	// Need do any validation?
	stat, err := property.save(ctx, p, domainName)
	if err != nil {
		return nil, err
	}
	return stat.Status, err

}

// Save Property updates method
func (property *Property) save(ctx context.Context, p *gtm, domainName string) (*PropertyResponse, error) {

	if err := property.Validate(); err != nil {
		return nil, fmt.Errorf("Property validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, property.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}

	var presp PropertyResponse
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &presp, property)
	if err != nil {
		return nil, fmt.Errorf("Property request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &presp, nil
}

// Delete the property identified by the receiver argument from the domain provided.
func (p *gtm) DeleteProperty(ctx context.Context, property *Property, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteProperty")

	if err := property.Validate(); err != nil {
		logger.Errorf("Property validation failed. %w", err)
		return nil, fmt.Errorf("Property validation failed. %w", err)
	}

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s/properties/%s", domainName, property.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Property request: %w", err)
	}

	var presp ResponseBody
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &presp)
	if err != nil {
		return nil, fmt.Errorf("Property request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return presp.Status, nil
}

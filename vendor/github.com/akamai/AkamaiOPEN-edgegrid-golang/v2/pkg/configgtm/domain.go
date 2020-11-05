package gtm

import (
	"context"
	"fmt"
	"net/http"

	"reflect"
	"strings"
	"unicode"
)

//
// Support gtm domains thru Edgegrid
// Based on 1.4 Schema
//

// Domains contains operations available on a Domain resource
// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html
type Domains interface {
	// Retrieve map of null fields
	NullFieldMap(context.Context, *Domain) (*NullFieldMapStruct, error)
	// NewDomain is a utility function that creates a new Domain object.
	NewDomain(context.Context, string, string) *Domain
	// GetStatus retrieves current status for the given domainname.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getcurrentstatus
	GetDomainStatus(context.Context, string) (*ResponseStatus, error)
	// ListDomains retrieves all Domains.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getdomains
	ListDomains(context.Context) ([]*DomainItem, error)
	// GetDomain retrieves a Domain with the given domainname.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#getdomain
	GetDomain(context.Context, string) (*Domain, error)
	// CreateDomain creates domain
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#postdomains
	CreateDomain(context.Context, *Domain, map[string]string) (*DomainResponse, error)
	// Delete is a method applied to a domain object resulting in removal.
	// See: ** Not Supported by API **
	DeleteDomain(context.Context, *Domain) (*ResponseStatus, error)
	// Update is a method applied to a domain object resulting in an update.
	// See: https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#postdomains
	UpdateDomain(context.Context, *Domain, map[string]string) (*ResponseStatus, error)
}

// The Domain data structure represents a GTM domain
type Domain struct {
	Name                         string          `json:"name"`
	Type                         string          `json:"type"`
	AsMaps                       []*AsMap        `json:"asMaps,omitempty"`
	Resources                    []*Resource     `json:"resources,omitempty"`
	DefaultUnreachableThreshold  float32         `json:"defaultUnreachableThreshold,omitempty"`
	EmailNotificationList        []string        `json:"emailNotificationList,omitempty"`
	MinPingableRegionFraction    float32         `json:"minPingableRegionFraction,omitempty"`
	DefaultTimeoutPenalty        int             `json:"defaultTimeoutPenalty,omitempty"`
	Datacenters                  []*Datacenter   `json:"datacenters,omitempty"`
	ServermonitorLivenessCount   int             `json:"servermonitorLivenessCount,omitempty"`
	RoundRobinPrefix             string          `json:"roundRobinPrefix,omitempty"`
	ServermonitorLoadCount       int             `json:"servermonitorLoadCount,omitempty"`
	PingInterval                 int             `json:"pingInterval,omitempty"`
	MaxTTL                       int64           `json:"maxTTL,omitempty"`
	LoadImbalancePercentage      float64         `json:"loadImbalancePercentage,omitempty"`
	DefaultHealthMax             float64         `json:"defaultHealthMax,omitempty"`
	LastModified                 string          `json:"lastModified,omitempty"`
	Status                       *ResponseStatus `json:"status,omitempty"`
	MapUpdateInterval            int             `json:"mapUpdateInterval,omitempty"`
	MaxProperties                int             `json:"maxProperties,omitempty"`
	MaxResources                 int             `json:"maxResources,omitempty"`
	DefaultSslClientPrivateKey   string          `json:"defaultSslClientPrivateKey,omitempty"`
	DefaultErrorPenalty          int             `json:"defaultErrorPenalty,omitempty"`
	Links                        []*Link         `json:"links,omitempty"`
	Properties                   []*Property     `json:"properties,omitempty"`
	MaxTestTimeout               float64         `json:"maxTestTimeout,omitempty"`
	CnameCoalescingEnabled       bool            `json:"cnameCoalescingEnabled"`
	DefaultHealthMultiplier      float64         `json:"defaultHealthMultiplier,omitempty"`
	ServermonitorPool            string          `json:"servermonitorPool,omitempty"`
	LoadFeedback                 bool            `json:"loadFeedback"`
	MinTTL                       int64           `json:"minTTL,omitempty"`
	GeographicMaps               []*GeoMap       `json:"geographicMaps,omitempty"`
	CidrMaps                     []*CidrMap      `json:"cidrMaps,omitempty"`
	DefaultMaxUnreachablePenalty int             `json:"defaultMaxUnreachablePenalty"`
	DefaultHealthThreshold       float64         `json:"defaultHealthThreshold,omitempty"`
	LastModifiedBy               string          `json:"lastModifiedBy,omitempty"`
	ModificationComments         string          `json:"modificationComments,omitempty"`
	MinTestInterval              int             `json:"minTestInterval,omitempty"`
	PingPacketSize               int             `json:"pingPacketSize,omitempty"`
	DefaultSslClientCertificate  string          `json:"defaultSslClientCertificate,omitempty"`
	EndUserMappingEnabled        bool            `json:"endUserMappingEnabled"`
}

type DomainsList struct {
	DomainItems []*DomainItem `json:"items"`
}

// DomainItem is a DomainsList item
type DomainItem struct {
	AcgId        string  `json:"acgId"`
	LastModified string  `json:"lastModified"`
	Links        []*Link `json:"links"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`
}

// Validate validates Domain
func (dom *Domain) Validate() error {

	if len(dom.Name) < 1 {
		return fmt.Errorf("Domain is missing Name")
	}
	if len(dom.Type) < 1 {
		return fmt.Errorf("Domain is missing Type")
	}

	return nil
}

// NewDomain is a utility function that creates a new Domain object.
func (p *gtm) NewDomain(ctx context.Context, domainName, domainType string) *Domain {

	logger := p.Log(ctx)
	logger.Debug("NewDomain")

	domain := &Domain{}
	domain.Name = domainName
	domain.Type = domainType
	return domain
}

// GetStatus retrieves current status for the given domainname.
func (p *gtm) GetDomainStatus(ctx context.Context, domainName string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("GetDomainStatus")

	var stat ResponseStatus
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s/status/current", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &stat)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &stat, nil
}

// ListDomains retrieves all Domains.
func (p *gtm) ListDomains(ctx context.Context) ([]*DomainItem, error) {

	logger := p.Log(ctx)
	logger.Debug("ListDomains")

	var domains DomainsList
	getURL := fmt.Sprintf("/config-gtm/v1/domains")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ListDomains request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &domains)
	if err != nil {
		return nil, fmt.Errorf("ListDomains request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return domains.DomainItems, nil
}

// GetDomain retrieves a Domain with the given domainname.
func (p *gtm) GetDomain(ctx context.Context, domainName string) (*Domain, error) {

	logger := p.Log(ctx)
	logger.Debug("GetDomain")

	var domain Domain
	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domainName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &domain)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &domain, nil
}

// Save method; Create or Update
func (domain *Domain) save(ctx context.Context, p *gtm, queryArgs map[string]string, req *http.Request) (*DomainResponse, error) {

	// set schema version
	setVersionHeader(req, schemaVersion)

	// Look for optional args
	if len(queryArgs) > 0 {
		q := req.URL.Query()
		if val, ok := queryArgs["contractId"]; ok {
			q.Add("contractId", strings.TrimPrefix(val, "ctr_"))
		}
		if val, ok := queryArgs["gid"]; ok {
			q.Add("gid", strings.TrimPrefix(val, "grp_"))
		}
		req.URL.RawQuery = q.Encode()
	}

	var dresp DomainResponse
	resp, err := p.Exec(req, &dresp, domain)
	if err != nil {
		return nil, fmt.Errorf("Domain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &dresp, nil

}

// Create is a method applied to a domain object resulting in creation.
func (p *gtm) CreateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*DomainResponse, error) {

	logger := p.Log(ctx)
	logger.Debug("CreateDomain")

	if err := domain.Validate(); err != nil {
		logger.Errorf("Domain validation failed. %w", err)
		return nil, fmt.Errorf("Domain validation failed. %w", err)
	}

	postURL := fmt.Sprintf("/config-gtm/v1/domains/")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, postURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CreateDomain request: %w", err)
	}

	return domain.save(ctx, p, queryArgs, req)

}

// Update is a method applied to a domain object resulting in an update.
func (p *gtm) UpdateDomain(ctx context.Context, domain *Domain, queryArgs map[string]string) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("UpdateDomain")

	if err := domain.Validate(); err != nil {
		logger.Errorf("Domain validation failed. %w", err)
		return nil, fmt.Errorf("Domain validation failed. %w", err)
	}

	putURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateDomain request: %w", err)
	}

	stat, err := domain.save(ctx, p, queryArgs, req)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Delete is a method applied to a domain object resulting in removal.
func (p *gtm) DeleteDomain(ctx context.Context, domain *Domain) (*ResponseStatus, error) {

	logger := p.Log(ctx)
	logger.Debug("DeleteDomain")

	delURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, delURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DeleteDomain request: %w", err)
	}

	var responseBody ResponseBody
	setVersionHeader(req, schemaVersion)

	resp, err := p.Exec(req, &responseBody)
	if err != nil {
		return nil, fmt.Errorf("Delete Domain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return responseBody.Status, nil
}

// NullObjectAttributeStruct represents core and child null onject attributes
type NullPerObjectAttributeStruct struct {
	CoreObjectFields  map[string]string
	ChildObjectFields map[string]interface{} // NullObjectAttributeStruct
}

// NullFieldMapStruct returned null Objects structure
type NullFieldMapStruct struct {
	Domain      NullPerObjectAttributeStruct            // entry is domain
	Properties  map[string]NullPerObjectAttributeStruct // entries are properties
	Datacenters map[string]NullPerObjectAttributeStruct // entries are datacenters
	Resources   map[string]NullPerObjectAttributeStruct // entries are resources
	CidrMaps    map[string]NullPerObjectAttributeStruct // entries are cidrmaps
	GeoMaps     map[string]NullPerObjectAttributeStruct // entries are geomaps
	AsMaps      map[string]NullPerObjectAttributeStruct // entries are asmaps
}

type ObjectMap map[string]interface{}

// Retrieve map of null fields
func (p *gtm) NullFieldMap(ctx context.Context, domain *Domain) (*NullFieldMapStruct, error) {

	logger := p.Log(ctx)
	logger.Debug("NullFieldMap")

	if err := domain.Validate(); err != nil {
		logger.Errorf("Domain validation failed. %w", err)
		return nil, fmt.Errorf("Domain validation failed. %w", err)
	}

	var nullFieldMap = &NullFieldMapStruct{}
	var domFields = NullPerObjectAttributeStruct{}
	domainMap := make(map[string]string)
	var objMap ObjectMap

	getURL := fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetDomain request: %w", err)
	}
	setVersionHeader(req, schemaVersion)
	resp, err := p.Exec(req, &objMap)
	if err != nil {
		return nil, fmt.Errorf("GetDomain request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	for i, d := range objMap {
		objval := fmt.Sprint(d)
		if fmt.Sprintf("%T", d) == "<nil>" {
			if objval == "<nil>" {
				domainMap[makeFirstCharUpperCase(i)] = ""
			}
			continue
		}
		list, ok := d.([]interface{})
		if !ok {
			continue
		}

		switch i {
		case "properties":
			nullFieldMap.Properties = processObjectList(list)
		case "datacenters":
			nullFieldMap.Datacenters = processObjectList(list)
		case "resources":
			nullFieldMap.Resources = processObjectList(list)
		case "cidrMaps":
			nullFieldMap.CidrMaps = processObjectList(list)
		case "geographicMaps":
			nullFieldMap.GeoMaps = processObjectList(list)
		case "asMaps":
			nullFieldMap.AsMaps = processObjectList(list)
		}
	}

	domFields.CoreObjectFields = domainMap
	nullFieldMap.Domain = domFields

	return nullFieldMap, nil

}

func makeFirstCharUpperCase(origString string) string {

	a := []rune(origString)
	a[0] = unicode.ToUpper(a[0])
	// hack
	if origString == "cname" {
		a[1] = unicode.ToUpper(a[1])
	}
	return string(a)
}

func processObjectList(objectList []interface{}) map[string]NullPerObjectAttributeStruct {

	nullObjectsList := make(map[string]NullPerObjectAttributeStruct)
	for _, obj := range objectList {
		nullObjectFields := NullPerObjectAttributeStruct{}
		objectName := ""
		objectDCID := ""
		objectMap := make(map[string]string)
		objectChildList := make(map[string]interface{})
		for objf, objd := range obj.(map[string]interface{}) {
			objval := fmt.Sprint(objd)
			switch fmt.Sprintf("%T", objd) {
			case "<nil>":
				if objval == "<nil>" {
					objectMap[makeFirstCharUpperCase(objf)] = ""
				}
			case "map[string]interface {}":
				// include null stand alone struct elements in core
				for moname, movalue := range objd.(map[string]interface{}) {
					if fmt.Sprintf("%T", movalue) == "<nil>" {
						objectMap[makeFirstCharUpperCase(moname)] = ""
					}
				}
			case "[]interface {}":
				iSlice := objd.([]interface{})
				if len(iSlice) > 0 && reflect.TypeOf(iSlice[0]).Kind() != reflect.String && reflect.TypeOf(iSlice[0]).Kind() != reflect.Int64 && reflect.TypeOf(iSlice[0]).Kind() != reflect.Float64 && reflect.TypeOf(iSlice[0]).Kind() != reflect.Int32 {
					objectChildList[makeFirstCharUpperCase(objf)] = processObjectList(objd.([]interface{}))
				}
			default:
				if objf == "name" {
					objectName = objval
				}
				if objf == "datacenterId" {
					objectDCID = objval
				}
			}
		}
		nullObjectFields.CoreObjectFields = objectMap
		nullObjectFields.ChildObjectFields = objectChildList

		if objectDCID == "" {
			if objectName != "" {
				nullObjectsList[objectName] = nullObjectFields
			} else {
				nullObjectsList["unknown"] = nullObjectFields // TODO: What if mnore than one?
			}
		} else {
			nullObjectsList[objectDCID] = nullObjectFields
		}
	}

	return nullObjectsList

}

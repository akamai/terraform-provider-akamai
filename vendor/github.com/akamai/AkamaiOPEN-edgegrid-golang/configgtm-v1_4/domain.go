package configgtm

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

//
// Support gtm domains thru Edgegrid
// Based on 1.4 Schema
//

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
	Datacenters                  []*Datacenter   `json:"datacenters"`
	ServermonitorLivenessCount   int             `json:"servermonitorLivenessCount,omitempty"`
	RoundRobinPrefix             string          `json:"roundRobinPrefix,omitempty"`
	ServermonitorLoadCount       int             `json:"servermonitorLoadCount,omitempty"`
	PingInterval                 int             `json:"pingInterval,omitempty"`
	MaxTTL                       int64           `json:"maxTTL,omitempty"`
	LoadImbalancePercentage      float64         `json:"loadImbalancePercentage,omitempty"`
	DefaultHealthMax             int             `json:"defaultHealthMax,omitempty"`
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
	CnameCoalescingEnabled       bool            `json:"cnameCoalescingEnabled,omitempty"`
	DefaultHealthMultiplier      int             `json:"defaultHealthMultiplier,omitempty"`
	ServermonitorPool            string          `json:"servermonitorPool,omitempty"`
	LoadFeedback                 bool            `json:"loadFeedback,omitempty"`
	MinTTL                       int64           `json:"minTTL,omitempty"`
	GeographicMaps               []*GeoMap       `json:"geographicMaps,omitempty"`
	CidrMaps                     []*CidrMap      `json:"cidrMaps,omitempty"`
	DefaultMaxUnreachablePenalty int             `json:"defaultMaxUnreachablePenalty,omitempty"`
	DefaultHealthThreshold       int             `json:"defaultHealthThreshold,omitempty"`
	LastModifiedBy               string          `json:"lastModifiedBy,omitempty"`
	ModificationComments         string          `json:"modificationComments,omitempty"`
	MinTestInterval              int             `json:"minTestInterval,omitempty"`
	PingPacketSize               int             `json:"pingPacketSize,omitempty"`
	DefaultSslClientCertificate  string          `json:"defaultSslClientCertificate,omitempty"`
	EndUserMappingEnabled        bool            `json:"endUserMappingEnabled,omitempty"`
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

// NewDomain is a utility function that creates a new Domain object.
func NewDomain(domainName, domainType string) *Domain {
	domain := &Domain{}
	domain.Name = domainName
	domain.Type = domainType
	return domain
}

// GetStatus retrieves current status for the given domainname.
func GetDomainStatus(domainName string) (*ResponseStatus, error) {
	stat := &ResponseStatus{}
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s/status/current", domainName),
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
		err = client.BodyJSON(res, stat)
		if err != nil {
			return nil, err
		}

		return stat, nil
	}
}

// ListDomains retrieves all Domains.
func ListDomains() ([]*DomainItem, error) {
	domains := &DomainsList{}
	req, err := client.NewRequest(
		Config,
		"GET",
		"/config-gtm/v1/domains/",
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
		return nil, CommonError{entityName: "Domain"}
	} else {
		err = client.BodyJSON(res, domains)
		if err != nil {
			return nil, err
		}

		return domains.DomainItems, nil
	}
}

// GetDomain retrieves a Domain with the given domainname.
func GetDomain(domainName string) (*Domain, error) {
	domain := NewDomain(domainName, "basic")
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domainName),
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
		err = client.BodyJSON(res, domain)
		if err != nil {
			return nil, err
		}

		return domain, nil
	}
}

// Save method; Create or Update
func (domain *Domain) save(queryArgs map[string]string, req *http.Request) (*DomainResponse, error) {

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

	printHttpRequest(req, true)

	res, err := client.Do(Config, req)

	// Network error
	if err != nil {
		return nil, CommonError{
			entityName:       "Domain",
			name:             domain.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domain.Name, apiErrorMessage: err.Detail, err: err}
	}

	// TODO: What validation can we do? E.g. if not equivalent there was a concurrent change...
	responseBody := &DomainResponse{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil

}

// Create is a method applied to a domain object resulting in creation.
func (domain *Domain) Create(queryArgs map[string]string) (*DomainResponse, error) {

	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf("/config-gtm/v1/domains/"),
		domain,
	)
	if err != nil {
		return nil, err
	}

	return domain.save(queryArgs, req)

}

// Update is a method applied to a domain object resulting in an update.
func (domain *Domain) Update(queryArgs map[string]string) (*ResponseStatus, error) {

	// Any validation to do?
	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name),
		domain,
	)
	if err != nil {
		return nil, err
	}

	stat, err := domain.save(queryArgs, req)
	if err != nil {
		return nil, err
	}
	return stat.Status, err
}

// Delete is a method applied to a domain object resulting in removal.
func (domain *Domain) Delete() (*ResponseStatus, error) {

	req, err := client.NewRequest(
		Config,
		"DELETE",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name),
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
			entityName:       "Domain",
			name:             domain.Name,
			httpErrorMessage: err.Error(),
			err:              err,
		}
	}

	printHttpResponse(res, true)

	// API error
	if client.IsError(res) {
		err := client.NewAPIError(res)
		return nil, CommonError{entityName: "Domain", name: domain.Name, apiErrorMessage: err.Detail, err: err}
	}

	responseBody := &ResponseBody{}
	// Unmarshall whole response body in case want status
	err = client.BodyJSON(res, responseBody)
	if err != nil {
		return nil, err
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
func (domain *Domain) NullFieldMap() (*NullFieldMapStruct, error) {

	var nullFieldMap = &NullFieldMapStruct{}
	var domFields = NullPerObjectAttributeStruct{}
	domainMap := make(map[string]string)
	var objMap = ObjectMap{}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf("/config-gtm/v1/domains/%s", domain.Name),
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
		return nil, CommonError{entityName: "Domain", name: domain.Name}
	} else {
		err = client.BodyJSON(res, &objMap)
		if err != nil {
			return nullFieldMap, err
		}
	}
	for i, d := range objMap {
		objval := fmt.Sprint(d)
		if fmt.Sprintf("%T", d) == "<nil>" {
			if objval == "<nil>" {
				domainMap[makeFirstCharUpperCase(i)] = ""
			}
			continue
		}
		switch i {
		case "properties":
			nullFieldMap.Properties = processObjectList(d.([]interface{}))
		case "datacenters":
			nullFieldMap.Datacenters = processObjectList(d.([]interface{}))
		case "resources":
			nullFieldMap.Resources = processObjectList(d.([]interface{}))
		case "cidrMaps":
			nullFieldMap.CidrMaps = processObjectList(d.([]interface{}))
		case "geographicMaps":
			nullFieldMap.GeoMaps = processObjectList(d.([]interface{}))
		case "asMaps":
			nullFieldMap.AsMaps = processObjectList(d.([]interface{}))
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

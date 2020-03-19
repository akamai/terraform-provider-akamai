package papi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/patrickmn/go-cache"
)

// EdgeHostnames is a collection for PAPI Edge Hostname resources
type EdgeHostnames struct {
	client.Resource
	AccountID     string `json:"accountId"`
	ContractID    string `json:"contractId"`
	GroupID       string `json:"groupId"`
	EdgeHostnames struct {
		Items []*EdgeHostname `json:"items"`
	} `json:"edgeHostnames"`
}

// NewEdgeHostnames creates a new EdgeHostnames
func NewEdgeHostnames() *EdgeHostnames {
	edgeHostnames := &EdgeHostnames{}
	edgeHostnames.Init()
	return edgeHostnames
}

// PostUnmarshalJSON is called after JSON unmarshaling into EdgeHostnames
//
// See: jsonhooks-v1/jsonhooks.Unmarshal()
func (edgeHostnames *EdgeHostnames) PostUnmarshalJSON() error {
	edgeHostnames.Init()

	for key, edgeHostname := range edgeHostnames.EdgeHostnames.Items {
		edgeHostnames.EdgeHostnames.Items[key].parent = edgeHostnames

		if err := edgeHostname.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	edgeHostnames.Complete <- true

	return nil
}

// NewEdgeHostname creates a new EdgeHostname within a given EdgeHostnames
func (edgeHostnames *EdgeHostnames) NewEdgeHostname() *EdgeHostname {
	edgeHostname := NewEdgeHostname(edgeHostnames)
	edgeHostnames.EdgeHostnames.Items = append(edgeHostnames.EdgeHostnames.Items, edgeHostname)
	return edgeHostname
}

// GetEdgeHostnames will populate EdgeHostnames with Edge Hostname data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listedgehostnames
// Endpoint: GET /papi/v1/edgehostnames/{?contractId,groupId,options}
func (edgeHostnames *EdgeHostnames) GetEdgeHostnames(contract *Contract, group *Group, options string) error {

	if contract == nil && group == nil {
		return errors.New("function requires at least \"group\" argument")
	}

	cacheedgehostnames, found := Profilecache.Get("edgehostnames")
	if found {
		json.Unmarshal(cacheedgehostnames.([]byte), edgeHostnames)
		return nil
	} else {

		if contract == nil && group != nil {
			contract = NewContract(NewContracts())
			contract.ContractID = group.ContractIDs[0]
		}

		if options != "" {
			options = fmt.Sprintf("&options=%s", options)
		}

		req, err := client.NewRequest(
			Config,
			"GET",
			fmt.Sprintf(
				"/papi/v1/edgehostnames?groupId=%s&contractId=%s%s",
				group.GroupID,
				contract.ContractID,
				options,
			),
			nil,
		)
		if err != nil {
			return err
		}

		res, err := client.Do(Config, req)
		if err != nil {
			return err
		}

		if client.IsError(res) {
			return client.NewAPIError(res)
		}

		if err = client.BodyJSON(res, edgeHostnames); err != nil {
			return err
		}

		byt, _ := json.Marshal(edgeHostnames)
		Profilecache.Set("edgehostnames", byt, cache.DefaultExpiration)
		return nil
	}
}

func (edgeHostnames *EdgeHostnames) FindEdgeHostname(edgeHostname *EdgeHostname) (*EdgeHostname, error) {
	if edgeHostname.DomainSuffix == "" && edgeHostname.EdgeHostnameDomain != "" {
		edgeHostname.DomainSuffix = "edgesuite.net"
		if strings.HasSuffix(edgeHostname.EdgeHostnameDomain, "edgekey.net") {
			edgeHostname.DomainSuffix = "edgekey.net"
		}
	}

	if edgeHostname.DomainPrefix == "" && edgeHostname.EdgeHostnameDomain != "" {
		edgeHostname.DomainPrefix = strings.TrimSuffix(edgeHostname.EdgeHostnameDomain, "."+edgeHostname.DomainSuffix)
	}

	if len(edgeHostnames.EdgeHostnames.Items) == 0 {
		return nil, errors.New("no hostnames found, did you call GetHostnames()?")
	}

	for _, eHn := range edgeHostnames.EdgeHostnames.Items {
		if (eHn.DomainPrefix == edgeHostname.DomainPrefix && eHn.DomainSuffix == edgeHostname.DomainSuffix) || eHn.EdgeHostnameID == edgeHostname.EdgeHostnameID {
			return eHn, nil
		}
	}

	return nil, nil
}

func (edgeHostnames *EdgeHostnames) AddEdgeHostname(edgeHostname *EdgeHostname) {
	found, err := edgeHostnames.FindEdgeHostname(edgeHostname)

	if err != nil || found == nil {
		edgeHostnames.EdgeHostnames.Items = append(edgeHostnames.EdgeHostnames.Items, edgeHostname)
	}

	if err == nil && found != nil && found.EdgeHostnameID == edgeHostname.EdgeHostnameID {
		*found = *edgeHostname
	}
}

// EdgeHostname represents an Edge Hostname resource
type EdgeHostname struct {
	client.Resource
	parent                 *EdgeHostnames
	EdgeHostnameID         string      `json:"edgeHostnameId,omitempty"`
	EdgeHostnameDomain     string      `json:"edgeHostnameDomain,omitempty"`
	ProductID              string      `json:"productId"`
	DomainPrefix           string      `json:"domainPrefix"`
	DomainSuffix           string      `json:"domainSuffix"`
	CertEnrollmentId       int         `json:"certEnrollmentId,omitempty"`
	SlotNumber             int         `json:"slotNumber,omitempty"`
	SecureNetwork          string      `json:"secureNetwork,omitempty"`
	Status                 StatusValue `json:"status,omitempty"`
	Secure                 bool        `json:"secure,omitempty"`
	IPVersionBehavior      string      `json:"ipVersionBehavior,omitempty"`
	MapDetailsSerialNumber int         `json:"mapDetails:serialNumber,omitempty"`
	MapDetailsSlotNumber   int         `json:"mapDetails:slotNumber,omitempty"`
	MapDetailsMapDomain    string      `json:"mapDetails:mapDomain,omitempty"`
	StatusChange           chan bool   `json:"-"`
}

// NewEdgeHostname creates a new EdgeHostname
func NewEdgeHostname(edgeHostnames *EdgeHostnames) *EdgeHostname {
	edgeHostname := &EdgeHostname{parent: edgeHostnames}
	edgeHostname.Init()
	return edgeHostname
}

func (edgeHostname *EdgeHostname) Init() {
	edgeHostname.Complete = make(chan bool, 1)
	edgeHostname.StatusChange = make(chan bool, 1)
}

// GetEdgeHostname populates EdgeHostname with data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getanedgehostname
// Endpoint: GET /papi/v1/edgehostnames/{edgeHostnameId}{?contractId,groupId,options}
func (edgeHostname *EdgeHostname) GetEdgeHostname(options string) error {
	if options != "" {
		options = "&options=" + options
	}

	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/edgehostnames/%s?contractId=%s&groupId=%s%s",
			edgeHostname.EdgeHostnameID,
			edgeHostname.parent.ContractID,
			edgeHostname.parent.GroupID,
			options,
		),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		if res.StatusCode == 404 {
			// Check collection for current hostname
			contract := NewContract(NewContracts())
			contract.ContractID = edgeHostname.parent.ContractID
			group := NewGroup(NewGroups())
			group.GroupID = edgeHostname.parent.GroupID

			edgeHostname.parent.GetEdgeHostnames(contract, group, "")
			newEdgeHostname, err := edgeHostname.parent.FindEdgeHostname(edgeHostname)
			if err != nil || newEdgeHostname == nil {
				return client.NewAPIError(res)
			}

			edgeHostname.EdgeHostnameID = newEdgeHostname.EdgeHostnameID
			edgeHostname.EdgeHostnameDomain = newEdgeHostname.EdgeHostnameDomain
			edgeHostname.ProductID = newEdgeHostname.ProductID
			edgeHostname.DomainPrefix = newEdgeHostname.DomainPrefix
			edgeHostname.DomainSuffix = newEdgeHostname.DomainSuffix
			edgeHostname.Status = newEdgeHostname.Status
			edgeHostname.Secure = newEdgeHostname.Secure
			edgeHostname.IPVersionBehavior = newEdgeHostname.IPVersionBehavior
			edgeHostname.MapDetailsSerialNumber = newEdgeHostname.MapDetailsSerialNumber
			edgeHostname.MapDetailsSlotNumber = newEdgeHostname.MapDetailsSlotNumber
			edgeHostname.MapDetailsMapDomain = newEdgeHostname.MapDetailsMapDomain

			return nil
		}

		return client.NewAPIError(res)
	}

	newEdgeHostnames := NewEdgeHostnames()
	if err := client.BodyJSON(res, newEdgeHostnames); err != nil {
		return err
	}

	edgeHostname.EdgeHostnameID = newEdgeHostnames.EdgeHostnames.Items[0].EdgeHostnameID
	edgeHostname.EdgeHostnameDomain = newEdgeHostnames.EdgeHostnames.Items[0].EdgeHostnameDomain
	edgeHostname.ProductID = newEdgeHostnames.EdgeHostnames.Items[0].ProductID
	edgeHostname.DomainPrefix = newEdgeHostnames.EdgeHostnames.Items[0].DomainPrefix
	edgeHostname.DomainSuffix = newEdgeHostnames.EdgeHostnames.Items[0].DomainSuffix
	edgeHostname.Status = newEdgeHostnames.EdgeHostnames.Items[0].Status
	edgeHostname.Secure = newEdgeHostnames.EdgeHostnames.Items[0].Secure
	edgeHostname.IPVersionBehavior = newEdgeHostnames.EdgeHostnames.Items[0].IPVersionBehavior
	edgeHostname.MapDetailsSerialNumber = newEdgeHostnames.EdgeHostnames.Items[0].MapDetailsSerialNumber
	edgeHostname.MapDetailsSlotNumber = newEdgeHostnames.EdgeHostnames.Items[0].MapDetailsSlotNumber
	edgeHostname.MapDetailsMapDomain = newEdgeHostnames.EdgeHostnames.Items[0].MapDetailsMapDomain

	return nil
}

// Save creates a new Edge Hostname
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#createanewedgehostname
// Endpoint: POST /papi/v1/edgehostnames/{?contractId,groupId,options}
func (edgeHostname *EdgeHostname) Save(options string) error {
	if options != "" {
		options = "&options=" + options
	}
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/papi/v1/edgehostnames/?contractId=%s&groupId=%s%s",
			edgeHostname.parent.ContractID,
			edgeHostname.parent.GroupID,
			options,
		),
		edgeHostname,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	var location client.JSONBody
	if err = client.BodyJSON(res, &location); err != nil {
		return err
	}

	// A 404 is returned until the hostname is valid, so just pull the new ID out for now
	url, _ := url.Parse(location["edgeHostnameLink"].(string))
	for _, part := range strings.Split(url.Path, "/") {
		if strings.HasPrefix(part, "ehn_") {
			edgeHostname.EdgeHostnameID = part
		}
	}

	edgeHostname.parent.AddEdgeHostname(edgeHostname)

	return nil
}

// PollStatus will responsibly poll till the property is active or an error occurs
//
// The EdgeHostname.StatusChange is a channel that can be used to
// block on status changes. If a new valid status is returned, true will
// be sent to the channel, otherwise, false will be sent.
//
//	go edgeHostname.PollStatus("")
//	for edgeHostname.Status != edgegrid.StatusActive {
//		select {
//		case statusChanged := <-edgeHostname.StatusChange:
//			if statusChanged == false {
//				break
//			}
//		case <-time.After(time.Minute * 30):
//			break
//		}
//	}
//
//	if edgeHostname.Status == edgegrid.StatusActive {
//		// EdgeHostname activated successfully
//	}
func (edgeHostname *EdgeHostname) PollStatus(options string) bool {
	currentStatus := edgeHostname.Status
	var retry time.Duration = 0
	for currentStatus != StatusActive {
		time.Sleep(retry)
		if retry == 0 {
			retry = time.Minute * 3
		}

		retry -= time.Minute

		err := edgeHostname.GetEdgeHostname(options)
		if err != nil {
			edgeHostname.StatusChange <- false
			return false
		}

		if currentStatus != edgeHostname.Status {
			edgeHostname.StatusChange <- true
		}
		currentStatus = edgeHostname.Status
	}

	return true
}

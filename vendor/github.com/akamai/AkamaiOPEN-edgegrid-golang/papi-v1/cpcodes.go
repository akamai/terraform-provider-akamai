package papi

import (
	"encoding/json"
	"fmt"

	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/patrickmn/go-cache"
)

// CpCodes represents a collection of CP Codes
//
// See: CpCodes.GetCpCodes()
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#cpcode
type CpCodes struct {
	client.Resource
	AccountID  string    `json:"accountId"`
	Contract   *Contract `json:"-"`
	ContractID string    `json:"contractId"`
	GroupID    string    `json:"groupId"`
	Group      *Group    `json:"-"`
	CpCodes    struct {
		Items []*CpCode `json:"items"`
	} `json:"cpcodes"`
}

// NewCpCodes creates a new *CpCodes
func NewCpCodes(contract *Contract, group *Group) *CpCodes {
	return &CpCodes{
		Contract: contract,
		Group:    group,
	}
}

// PostUnmarshalJSON is called after UnmarshalJSON to setup the
// structs internal state. The cpcodes.Complete channel is utilized
// to communicate full completion.
func (cpcodes *CpCodes) PostUnmarshalJSON() error {
	cpcodes.Init()

	cpcodes.Contract = NewContract(NewContracts())
	cpcodes.Contract.ContractID = cpcodes.ContractID

	cpcodes.Group = NewGroup(NewGroups())
	cpcodes.Group.GroupID = cpcodes.GroupID

	go cpcodes.Group.GetGroup()
	go cpcodes.Contract.GetContract()

	go (func(cpcodes *CpCodes) {
		contractComplete := <-cpcodes.Contract.Complete
		groupComplete := <-cpcodes.Group.Complete
		cpcodes.Complete <- (contractComplete && groupComplete)
	})(cpcodes)

	for key, cpcode := range cpcodes.CpCodes.Items {
		cpcodes.CpCodes.Items[key].parent = cpcodes

		if err := cpcode.PostUnmarshalJSON(); err != nil {
			return err
		}
	}

	return nil
}

// GetCpCodes populates a *CpCodes with it's related CP Codes
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listcpcodes
// Endpoint: GET /papi/v1/cpcodes/{?contractId,groupId}
func (cpcodes *CpCodes) GetCpCodes(correlationid string) error {
	cachecpcodes, found := Profilecache.Get("cpcodes")
	if found {
		json.Unmarshal(cachecpcodes.([]byte), cpcodes)
		return nil
	} else {
		if cpcodes.Contract == nil {
			cpcodes.Contract = NewContract(NewContracts())
			cpcodes.Contract.ContractID = cpcodes.Group.ContractIDs[0]
		}

		req, err := client.NewRequest(
			Config,
			"GET",
			fmt.Sprintf(
				"/papi/v1/cpcodes?groupId=%s&contractId=%s",
				cpcodes.Group.GroupID,
				cpcodes.Contract.ContractID,
			),
			nil,
		)
		if err != nil {
			return err
		}

		edge.PrintHttpRequestCorrelation(req, true, correlationid)

		res, err := client.Do(Config, req)
		if err != nil {
			return err
		}

		edge.PrintHttpResponseCorrelation(res, true, correlationid)

		if client.IsError(res) {
			return client.NewAPIError(res)
		}

		if err = client.BodyJSON(res, cpcodes); err != nil {
			return err
		}
		byt, _ := json.Marshal(cpcodes)
		Profilecache.Set("cpcodes", byt, cache.DefaultExpiration)
		return nil
	}
}

func (cpcodes *CpCodes) FindCpCode(nameOrId string, correlationid string) (*CpCode, error) {
	if len(cpcodes.CpCodes.Items) == 0 {
		err := cpcodes.GetCpCodes(correlationid)
		if err != nil {
			return nil, err
		}
		if len(cpcodes.CpCodes.Items) == 0 {
			return nil, fmt.Errorf("unable to fetch CP codes for group/contract")
		}
	}

	for _, cpcode := range cpcodes.CpCodes.Items {
		if cpcode.CpcodeID == nameOrId || cpcode.CpcodeID == "cpc_"+nameOrId || cpcode.CpcodeName == nameOrId {
			return cpcode, nil
		}
	}

	return nil, nil
}

// NewCpCode creates a new *CpCode associated with this *CpCodes as it's parent.
func (cpcodes *CpCodes) NewCpCode() *CpCode {
	cpcode := NewCpCode(cpcodes)
	cpcodes.AddCpCode(cpcode)
	return cpcode
}

func (cpcodes *CpCodes) AddCpCode(cpcode *CpCode) {
	var exists bool
	for _, cpc := range cpcodes.CpCodes.Items {
		if cpc == cpcode {
			exists = true
		}
	}

	if !exists {
		cpcodes.CpCodes.Items = append(cpcodes.CpCodes.Items, cpcode)
	}
}

// CpCode represents a single CP Code
//
// API Docs: https://developer.akamai.com/api/luna/papi/data.html#cpcode
type CpCode struct {
	client.Resource
	parent      *CpCodes
	CpcodeID    string    `json:"cpcodeId,omitempty"`
	CpcodeName  string    `json:"cpcodeName"`
	ProductID   string    `json:"productId,omitempty"`
	ProductIDs  []string  `json:"productIds,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty"`
}

// NewCpCode creates a new *CpCode
func NewCpCode(parent *CpCodes) *CpCode {
	cpcode := &CpCode{parent: parent}
	cpcode.Init()
	return cpcode
}

// GetCpCode populates the *CpCode with it's data
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getacpcode
// Endpoint: GET /papi/v1/cpcodes/{cpcodeId}{?contractId,groupId}
func (cpcode *CpCode) GetCpCode() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/cpcodes/%s?contractId=%s&groupId=%s",
			cpcode.CpcodeID,
			cpcode.parent.Contract.ContractID,
			cpcode.parent.Group.GroupID,
		),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequest(req, true)

	res, err := client.Do(Config, req)

	if err != nil {
		return err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	newCpcodes := NewCpCodes(nil, nil)
	if err = client.BodyJSON(res, newCpcodes); err != nil {
		return err
	}
	if len(newCpcodes.CpCodes.Items) == 0 {
		return fmt.Errorf("CP Code \"%s\" not found", cpcode.CpcodeID)
	}

	cpcode.CpcodeID = newCpcodes.CpCodes.Items[0].CpcodeID
	cpcode.CpcodeName = newCpcodes.CpCodes.Items[0].CpcodeName
	cpcode.ProductID = newCpcodes.CpCodes.Items[0].ProductID
	cpcode.ProductIDs = newCpcodes.CpCodes.Items[0].ProductIDs
	cpcode.CreatedDate = newCpcodes.CpCodes.Items[0].CreatedDate

	cpcode.parent.AddCpCode(cpcode)

	return nil
}

// ID retrieves a CP Codes integer ID
//
// PAPI Behaviors require the integer ID, rather than the prefixed string returned
func (cpcode *CpCode) ID() int {
	id, err := strconv.Atoi(strings.TrimPrefix(cpcode.CpcodeID, "cpc_"))
	if err != nil {
		return 0
	}

	return id
}

// Save will create a new CP Code. You cannot update a CP Code;
// trying to do so will result in an error.
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#createanewcpcode
// Endpoint: POST /papi/v1/cpcodes/{?contractId,groupId}
func (cpcode *CpCode) Save(correlationid string) error {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/papi/v1/cpcodes?contractId=%s&groupId=%s",
			cpcode.parent.ContractID,
			cpcode.parent.GroupID,
		),
		client.JSONBody{"productId": cpcode.ProductID, "cpcodeName": cpcode.CpcodeName},
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequestCorrelation(req, true, correlationid)

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponseCorrelation(res, true, correlationid)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	var location client.JSONBody
	if err = client.BodyJSON(res, &location); err != nil {
		return err
	}

	req, err = client.NewRequest(
		Config,
		"GET",
		location["cpcodeLink"].(string),
		nil,
	)
	if err != nil {
		return err
	}

	edge.PrintHttpRequest(req, true)

	res, err = client.Do(Config, req)
	if err != nil {
		return err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	cpcodes := NewCpCodes(nil, nil)
	if err != nil {
		return err
	}

	if err = client.BodyJSON(res, cpcodes); err != nil {
		return err
	}

	newCpcode := cpcodes.CpCodes.Items[0]
	newCpcode.parent = cpcode.parent

	//cpcode.parent.CpCodes.Items = append(cpcode.parent.CpCodes.Items, newCpcode)

	cpcode.CpcodeID = cpcodes.CpCodes.Items[0].CpcodeID
	cpcode.CpcodeName = cpcodes.CpCodes.Items[0].CpcodeName
	cpcode.ProductIDs = cpcodes.CpCodes.Items[0].ProductIDs
	cpcode.CreatedDate = cpcodes.CpCodes.Items[0].CreatedDate

	cpcode.parent.AddCpCode(cpcode)

	return nil
}

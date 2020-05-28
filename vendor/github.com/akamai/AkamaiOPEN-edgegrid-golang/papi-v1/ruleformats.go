package papi

import (
	"fmt"
	"io/ioutil"
	"sort"

        edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/xeipuuv/gojsonschema"
)

// RuleFormats is a collection of available rule formats
type RuleFormats struct {
	client.Resource
	RuleFormats struct {
		Items []string `json:"items"`
	} `json:"ruleFormats"`
}

// NewRuleFormats creates a new RuleFormats
func NewRuleFormats() *RuleFormats {
	ruleFormats := &RuleFormats{}
	ruleFormats.Init()

	return ruleFormats
}

// GetRuleFormats populates RuleFormats
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#listruleformats
// Endpoint: GET /papi/v1/rule-formats
func (ruleFormats *RuleFormats) GetRuleFormats() error {
	req, err := client.NewRequest(
		Config,
		"GET",
		"/papi/v1/rule-formats",
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

	if err := client.BodyJSON(res, ruleFormats); err != nil {
		return err
	}

	sort.Strings(ruleFormats.RuleFormats.Items)

	return nil
}

func (ruleFormats *RuleFormats) GetLatest() (string, error) {
	if len(ruleFormats.RuleFormats.Items) == 0 {
		err := ruleFormats.GetRuleFormats()
		if err != nil {
			return "", err
		}
	}

	return ruleFormats.RuleFormats.Items[len(ruleFormats.RuleFormats.Items)-1], nil
}

// GetSchema fetches the schema for a given product and rule format
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruleformatsschema
// Endpoint: /papi/v1/schemas/products/{productId}/{ruleFormat}
func (ruleFormats *RuleFormats) GetSchema(product string, ruleFormat string) (*gojsonschema.Schema, error) {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/schemas/products/%s/%s",
			product,
			ruleFormat,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	edge.PrintHttpRequest(req, true)

	res, err := client.Do(Config, req)
	if err != nil {
		return nil, err
	}

	edge.PrintHttpResponse(res, true)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	schemaBytes, _ := ioutil.ReadAll(res.Body)
	schemaBody := string(schemaBytes)
	loader := gojsonschema.NewStringLoader(schemaBody)
	schema, err := gojsonschema.NewSchema(loader)

	return schema, err
}

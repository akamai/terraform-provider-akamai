package cps

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
)

// Init sets the CPS edgegrid Config
func Init(config edgegrid.Config) {
	Config = config
}

func newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] newRequest, buf: %s", string(buf.Bytes()))

	req, err := client.NewRequest(Config, method, urlStr, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/vnd.akamai.cps.enrollment.v7+json")
	req.Header.Add("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")

	return req, nil
}

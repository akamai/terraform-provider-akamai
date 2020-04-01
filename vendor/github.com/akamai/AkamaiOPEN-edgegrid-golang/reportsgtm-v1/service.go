package reportsgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"net/http"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
)

// Init sets the GTM edgegrid Config
func Init(config edgegrid.Config) {

	Config = config
	edgegrid.SetupLogging()

}

// Utility func to print http req
func printHttpRequest(req *http.Request, body bool) {

	edgegrid.PrintHttpRequest(req, body)

}

// Utility func to print http response
func printHttpResponse(res *http.Response, body bool) {

	edgegrid.PrintHttpResponse(res, body)

}

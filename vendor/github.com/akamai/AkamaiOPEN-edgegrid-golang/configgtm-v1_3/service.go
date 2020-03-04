package configgtm

import (
	"net/http"
	"net/http/httputil"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/sirupsen/logrus"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
	// Create a new instance of the logger.
	GtmLog *logrus.Logger
)

// Init sets the GTM edgegrid Config
func Init(config edgegrid.Config) {

	Config = config
	GtmLog = logrus.New()
	edgegrid.SetupLogging(GtmLog)
	if edgegrid.LogFile != nil {
		defer edgegrid.LogFile.Close()
	}
}

// Utility func to print http req
func printHttpRequest(req *http.Request, body bool) {
	if req == nil {
		return
	}

	b, err := httputil.DumpRequestOut(req, body)
	if err == nil {
		edgegrid.LogMultiline(GtmLog.Traceln, string(b))
	}
}

// Utility func to print http response
func printHttpResponse(res *http.Response, body bool) {
	if res == nil {
		return
	}

	b, err := httputil.DumpResponse(res, body)
	if err == nil {
		edgegrid.LogMultiline(GtmLog.Traceln, string(b))
	}
}

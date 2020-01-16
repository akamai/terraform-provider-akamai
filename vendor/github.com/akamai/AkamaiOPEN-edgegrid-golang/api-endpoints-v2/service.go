package apiendpoints

import (
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

var (
	// Config contains the Akamai OPEN Edgegrid API credentials
	// for automatic signing of requests
	Config edgegrid.Config
)

// Init sets the CCU edgegrid Config
func Init(config edgegrid.Config) {
	Config = config
}

func call(req *http.Request, err error) (*Endpoint, error) {
	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	rep := &Endpoint{}
	if err = client.BodyJSON(res, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

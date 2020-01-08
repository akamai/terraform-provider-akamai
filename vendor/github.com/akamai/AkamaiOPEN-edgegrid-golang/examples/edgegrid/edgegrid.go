package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func main() {
	client := http.Client{}

	/*
		Init will try to use the environment and fallback to .edgerc.

		This function will check in order:

		AKAMAI_{SECTION}_* environment variables
		if using the default section, AKAMAI_* environment variables
		the specified (or default if none) section in .edgerc
		if not using the default section, AKAMAI_* environment variables
		This new function is the recommended way for instantiating an instance.

		The environment variables are:

		AKAMAI_HOST or AKAMAI_{SECTION}_HOST
		AKAMAI_CLIENT_TOKEN or AKAMAI_{SECTION}_CLIENT_TOKEN
		AKAMAI_CLIENT_SECRET or AKAMAI_{SECTION}_CLIENT_SECRET
		AKAMAI_ACCESS_TOKEN or AKAMAI_{SECTION}_ACCESS_TOKEN
	*/
	config, err := edgegrid.InitEdgeRc("~/.edgerc", "default")

	if err == nil {
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://%s/diagnostic-tools/v2/ghost-locations/available", config.Host), nil)
		req = edgegrid.AddRequestHeader(config, req)
		resp, _ := client.Do(req)
		byt, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(byt))
	}
}

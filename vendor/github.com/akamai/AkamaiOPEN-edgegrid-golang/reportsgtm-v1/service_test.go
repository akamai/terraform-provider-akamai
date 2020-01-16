package reportsgtm

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        false,
	}
)

var gtmTestDomain = "gtmtest.akadns.net"

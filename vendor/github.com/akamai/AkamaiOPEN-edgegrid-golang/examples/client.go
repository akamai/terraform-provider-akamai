// An example Diagnostic Tools v1 API Client
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func random(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Intn(max-min) + min

	return random
}

type LocationsResponse struct {
	Locations []string `json:"locations"`
}

type DigResponse struct {
	Dig struct {
		Hostname    string `json:"hostname"`
		QueryType   string `json:"queryType"`
		Result      string `json:"result"`
		ErrorString string `json:"errorString"`
	} `json:"dig"`
}

func main() {
	Example()
}

func Example() {
	config, err := edgegrid.InitEdgeRc("~/.edgerc", "default")
	config.Debug = false
	if err == nil {
		if err == nil {
			fmt.Println("Requesting locations that support the diagnostic-tools API.")

			req, err := client.NewRequest(
				config,
				"GET",
				"/diagnostic-tools/v1/locations",
				nil,
			)
			if err != nil {
				log.Fatal(err.Error())
			}

			res, err := client.Do(config, req)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			locationsResponse := LocationsResponse{}
			client.BodyJSON(res, &locationsResponse)

			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Printf("There are %d locations that can run dig in the Akamai Network\n", len(locationsResponse.Locations))

			if len(locationsResponse.Locations) == 0 {
				log.Fatal("No locations found")
			}

			location := locationsResponse.Locations[random(0, len(locationsResponse.Locations))-1]

			fmt.Println("We will make our call from " + location)

			fmt.Println("Running dig from " + location)

			client.Client.Timeout = 5 * time.Minute
			req, err = client.NewRequest(
				config,
				"GET",
				"/diagnostic-tools/v1/dig?hostname=developer.akamai.com&location="+url.QueryEscape(location)+"&queryType=A",
				nil,
			)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			res, err = client.Do(config, req)
			if err != nil {
				log.Fatal(err.Error())
				return
			}

			digResponse := DigResponse{}
			client.BodyJSON(res, &digResponse)
			fmt.Println(digResponse.Dig.Result)
		} else {
			log.Fatal(err.Error())
		}
	}
}

//An example Diagnostic Tools v1 API Client
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
)

func random(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	random := rand.Intn(max-min) + min

	return random
}

//Location ghost location type
type Location struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

//LocationsResponse response type for ghost locations
type LocationsResponse struct {
	Locations []Location `json:"locations"`
}

//DigResponse response type for dig API
type DigResponse struct {
	Dig struct {
		Hostname    string `json:"hostname"`
		QueryType   string `json:"queryType"`
		Result      string `json:"result"`
		ErrorString string `json:"errorString"`
	} `json:"digInfo"`
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
				"/diagnostic-tools/v2/ghost-locations/available",
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

			fmt.Println("We will make our call from " + location.Value)

			fmt.Println("Running dig from " + location.Value)

			client.Client.Timeout = 5 * time.Minute
			req, err = client.NewRequest(
				config,
				"GET",
				"/diagnostic-tools/v2/ghost-locations/"+location.ID+"/dig-info?hostName=developer.akamai.com&queryType=A",
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

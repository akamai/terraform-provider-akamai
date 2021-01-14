provider "akamai" {
  edgerc = "~/.edgerc_network"
}


resource "akamai_networklist_network_list" "test" {
   name =  "Martin Network List"
    type =  "IP"
    description = "Notes about this network list"
  
    list = ["10.1.8.23","10.3.5.67"] 
    mode = "REPLACE"
   }


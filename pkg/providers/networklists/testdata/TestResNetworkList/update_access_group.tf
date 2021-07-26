provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_networklist_network_list" "test" {
   name =  "Voyager Call Center Whitelist"
    type =  "IP"
    description = "Notes about this network list"
    group = 31325
    contract = "C-D5TW8R"
  
    list = ["10.1.8.23","10.3.5.67"] 
    mode = "REPLACE"
   }


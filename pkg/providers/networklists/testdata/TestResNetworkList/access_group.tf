provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_networklist_network_list" "test" {
   name =  "Voyager Call Center Whitelist"
    type =  "IP"
    description = "Notes about this network list"
    access_control_group = " Kona Site Defender - C-D5TW8R - C-D5TW8R.G31325"
  
    list = ["10.1.8.23","10.3.5.67"] 
    mode = "REPLACE"
   }


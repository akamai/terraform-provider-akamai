provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "default"
}

resource "akamai_cp_code" "example" {
  contract_id = "ctr_XXX"
  group_id    = "grp_XXX"
  name        = "example-XXX"
  product_id  = "prd_XXX"
}

resource "akamai_property" "example" {
  cp_code = "${akamai_cp_code.example.id}"

  // ...
  // configure additional property fields
  // ...
}

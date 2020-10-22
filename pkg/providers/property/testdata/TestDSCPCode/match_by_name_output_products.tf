provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cp_code" "test" {
  name     = "test cpcode"
  contract = "ctr_test"
  group    = "grp_test"
}

output "products" {
  value = data.akamai_cp_code.test.product_ids
}

output "product1" {
  value = data.akamai_cp_code.test.product_ids[0]
}

output "product2" {
  value = data.akamai_cp_code.test.product_ids[1]
}
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cp_code" "test" {
  name        = "test cpcode"
  contract_id = "ctr_11"
  group_id    = "grp_22"
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
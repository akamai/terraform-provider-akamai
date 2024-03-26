provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cp_code" "akacpcode" {
  contract_id = var.contractid
  group_id    = var.groupid
  product_id  = var.product
  name        = var.cp_code
}

data "akamai_cp_code" "akacpcodeq" {
  contract_id = var.contractid
  group_id    = var.groupid
  name        = akamai_cp_code.akacpcode.id
  # Fetch the newly created CP code
  depends_on = [
    akamai_cp_code.akacpcode
  ]
}

output "aka_cp_code" {
  value = data.akamai_cp_code.akacpcodeq.id
}
output "aka_cp_contract" {
  value = data.akamai_cp_code.akacpcodeq.contract_id
}

variable "groupid" {
  description = "Name of the group associated with this CP code"
  type        = string
  default     = "grp_15225"
}

variable "contractid" {
  description = "Contract ID associated with this CP code"
  type        = string
  default     = "ctr_1-1TJZH5"
}

variable "product" {
  description = "Name of the product associated with this CP code"
  type        = string
  default     = "prd_prod1"
}

variable "cp_code" {
  description = "CP code to be created or re-used"
  type        = string
  default     = "test-ft-cp-code"
}



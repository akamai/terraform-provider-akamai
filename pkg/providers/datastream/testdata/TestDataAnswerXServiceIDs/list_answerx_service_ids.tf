provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_datastream_answerx_service_ids" "test" {
  contract_id = "test_contract"
}

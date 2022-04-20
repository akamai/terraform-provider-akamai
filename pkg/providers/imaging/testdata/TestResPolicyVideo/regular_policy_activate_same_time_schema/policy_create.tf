provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_imaging_policy_video" "policy" {
  policy_id              = "test_policy"
  contract_id            = "test_contract"
  policyset_id           = "test_policy_set"
  activate_on_production = true
  policy {
    breakpoints {
      widths = [
        320,
        640,
        1024,
        2048,
        5000
      ]
    }
    output {
      perceptual_quality = "mediumHigh"
    }
  }
}
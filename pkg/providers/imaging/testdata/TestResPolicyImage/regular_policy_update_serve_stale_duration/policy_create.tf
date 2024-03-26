provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_imaging_policy_image" "policy" {
  policy_id    = "test_policy"
  contract_id  = "test_contract"
  policyset_id = "test_policy_set"
  json         = data.akamai_imaging_policy_image.policy.json
}

data "akamai_imaging_policy_image" "policy" {
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
    transformations {
      max_colors {
        colors = 2
      }
    }
  }
}
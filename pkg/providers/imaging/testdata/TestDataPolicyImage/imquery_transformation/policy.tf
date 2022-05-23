provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_imaging_policy_image" "policy" {
  policy {
    breakpoints {
      widths = [280, 1080]
    }
    hosts = ["host1", "host2"]
    output {
      adaptive_quality   = 50
      perceptual_quality = "mediumHigh"
    }
    transformations {
      im_query {
        allowed_transformations = ["Append", "Blur", "Crop"]
        query_var               = "imq"
      }
    }
    variables {
      default_value = ""
      name          = "imq"
      type          = "string"
    }
  }
}
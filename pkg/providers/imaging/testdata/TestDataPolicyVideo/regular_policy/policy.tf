provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_imaging_policy_video" "policy" {
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
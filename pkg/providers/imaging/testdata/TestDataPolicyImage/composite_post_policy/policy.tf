provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
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
    post_breakpoint_transformations {
      composite {
        gravity = "NorthWest"
        image {
          text_image {
            fill        = "#000000"
            size        = 72
            stroke      = "#FFFFFF"
            stroke_size = 0
            text        = "test"
          }
        }
        placement  = "Over"
        x_position = 0
        y_position = 0
      }
    }
    transformations {
      max_colors {
        colors = 0
      }
    }
  }
}
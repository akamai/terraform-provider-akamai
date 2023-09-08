provider "akamai" {
  edgerc = "../../test/edgerc"
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
      region_of_interest_crop {
        gravity = "Center"
        height  = 8
        region_of_interest {
          rectangle_shape {
            anchor {
              x = 4
              y = 5
            }
            height = 9
            width  = 8
          }
        }
        style = "fill"
        width = 7
      }
    }
    transformations {
      append {
        gravity          = "Center"
        gravity_priority = "horizontal"
        image {
          text_image {
            fill        = "#000000"
            size        = 72
            stroke      = "#FFFFFF"
            stroke_size = 0
            text        = "test"
          }
        }
        preserve_minor_dimension = true
      }
    }
    transformations {
      trim {
        fuzz    = 0.08
        padding = 0
      }
    }
    transformations {
      if_dimension {
        default {
          resize {
            aspect     = "fit"
            height_var = "ResizeDim"
            type       = "normal"
            width_var  = "ResizeDim"
          }
        }
        dimension = "height"
        less_than {
          crop {
            allow_expansion = true
            gravity         = "Center"
            height_var      = "ResizeDim"
            width_var       = "ResizeDim"
            x_position      = 0
            y_position      = 0
          }
        }
        value_var = "MinDim"
      }
    }
    variables {
      default_value = "280"
      name          = "ResizeDim"
      type          = "number"
    }
    variables {
      default_value = "260"
      name          = "ResizeDimWithBorder"
      type          = "number"
    }
    variables {
      default_value = "1000"
      name          = "MinDim"
      type          = "number"
    }
    variables {
      default_value = "1450"
      name          = "MinDimNew"
      type          = "number"
    }
    variables {
      default_value = "1500"
      name          = "MaxDimOld"
      type          = "number"
    }
  }
}
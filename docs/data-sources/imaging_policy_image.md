---
layout: akamai
subcategory: Image and Video Manager
---

# akamai_imaging_policy_image

Specifies details for an image policy, such as transformations to apply and variations in image size and formats.

## Example usage

### Basic usage

This example shows a simple image policy with image transformations that doesn’t contain variables.

```hcl
data "akamai_imaging_policy_image" "image_policy" {
  policy {
    rollout_duration = 3601
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
        colors = var.colors
      }
    }
    transformations {
      blur {
        sigma = var.sigma
      }
    }
    transformations {
      max_colors {
        colors = 5
      }
    }
  }
}
```

### Usage with transformation nesting and variables

This example shows how you can nest transformations and use variables for your policy settings:

```hcl
data "akamai_imaging_policy_image" "image_policy" {
  policy {
    rollout_duration = 3600
    breakpoints {
      widths = [280, 1080]
    }
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
              point_shape {
                x = 4
                y = 5
              }
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
            transformation {
              compound {
              }
            }
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

          compound {
            if_dimension {
              default {

                compound {
                  if_dimension {
                    default {

                      compound {
                        if_dimension {
                          default {

                            compound {
                              resize {
                                aspect     = "fit"
                                height_var = "ResizeDim"
                                type       = "normal"
                                width_var  = "ResizeDim"
                              }
                            }
                            compound {
                              crop {
                                allow_expansion = true
                                gravity         = "Center"
                                height_var      = "ResizeDim"
                                width_var       = "ResizeDim"
                                x_position      = 0
                                y_position      = 0
                              }
                            }
                            compound {
                              background_color {
                                color = "#ffffff"
                              }
                            }
                          }
                          dimension = "height"
                          greater_than {

                            compound {
                              resize {
                                aspect     = "fit"
                                height_var = "ResizeDimWithBorder"
                                type       = "normal"
                                width_var  = "ResizeDimWithBorder"
                              }
                            }
                            compound {
                              crop {
                                allow_expansion = true
                                gravity         = "Center"
                                height_var      = "ResizeDim"
                                width_var       = "ResizeDim"
                                x_position      = 0
                                y_position      = 0
                              }
                            }
                            compound {
                              background_color {
                                color = "#ffffff"
                              }
                            }
                          }
                          value_var = "MaxDimOld"
                        }
                      }
                    }
                    dimension = "height"
                    less_than {

                      compound {
                        resize {
                          aspect     = "fit"
                          height_var = "ResizeDimWithBorder"
                          type       = "normal"
                          width_var  = "ResizeDimWithBorder"
                        }
                      }
                      compound {
                        crop {
                          allow_expansion = true
                          gravity         = "Center"
                          height_var      = "ResizeDim"
                          width_var       = "ResizeDim"
                          x_position      = 0
                          y_position      = 0
                        }
                      }
                      compound {
                        background_color {
                          color = "#ffffff"
                        }
                      }
                    }
                    value_var = "MinDim"
                  }
                }
              }
              dimension = "width"
              less_than {

                compound {
                  resize {
                    aspect     = "fit"
                    height_var = "ResizeDimWithBorder"
                    type       = "normal"
                    width_var  = "ResizeDimWithBorder"
                  }
                }
                compound {
                  crop {
                    allow_expansion = true
                    gravity         = "Center"
                    height_var      = "ResizeDim"
                    width_var       = "ResizeDim"
                    x_position      = 0
                    y_position      = 0
                  }
                }
                compound {
                  background_color {
                    color = "#ffffff"
                  }
                }
              }
              value_var = "MinDim"
            }
          }
        }
        dimension = "width"
        value_var = "MaxDimOld"
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

      default_value = ""
      name          = "VariableWithoutDefaultValue"
      type          = "string"
    }
    variables {

      default_value = "1000"
      enum_options {

        id    = "1"
        value = "value1"
      }
      enum_options {

        id    = "2"
        value = "value2"
      }
      name = "MinDim"
      type = "number"
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
```

## Argument reference

This data source supports these arguments:

* `policy` - (Required) The image policy.
    * `breakpoints` - (Optional) The breakpoint widths (in pixels) to use to create derivative images.
        *  `widths` - (Optional) The list of breakpoint widths to use, separated by commas.
    * `hosts` - (Optional) The hosts that are allowed for image URLs within transformations or variables.
    * `output` - (Optional) The output quality and formats that are created for each resized image. If unspecified, image formats are created to support all browsers at the default quality level (`85`)including formats such as WEBP, JPEG2000 and JPEG-XR for specific browsers.
        * `adaptive_quality` - (Optional) The quality value that is applied to an image when Image and Video Manager detects a slow connection (RTT > 300 ms). This value (1-100) overrides the derivative quality value. Specifying a lower value can reduce load times for users with slow connections without impacting the quality of images for users with standard connections.
        * `perceptual_quality` - (Optional) The perceptual quality to use when comparing resulting images. Perceptual quality tunes each image format's quality parameter dynamically based on the human-perceived quality of the output image. The visual quality of derivative videos after they have been compressed to maximize byte savings. You can select one of five perceptual quality values: `high`, `mediumHigh`, `medium`, `mediumLow`, or `low`. `high` results in the highest visual quality with the least byte savings, while ‘low’ results in lower visual quality with the greatest byte savings. Either `perceptual quality can be set or `quality`, but not both.  If setting a variable for this argument, use `perceptual_quality_var` instead.
        * `perceptual_quality_floor` - (Optional) The minimum image quality to respect when perceptual quality is set. Image quality will not be reduced below this value even if it is determined that a further compressed image would be acceptably visually similar. Set a value between 1-100.  
        * `quality` - The specified quality of the output images expressed as a value from 1-100. Either `perceptual quality can be set or `quality`, but not both. If neither is set, `quality` is the default. If setting a variable for this argument, use `quality_var` instead.
    * `post_breakpoint_transformations` - (Optional) The set of post-processing transformations applied to the image after image and quality settings have been applied. This is a subset of the complete list of transformations and includes: `background_color`, `blur`, `chroma_key`, `composite`, `compound`, `contrast`, `goop`, `grayscale`, `hsl`, `hsv`, `if_dimension`, `if_orientation`, `max_colors`, `mirror`, `mono_hue`, `opacity`, `remove_color`, `unsharp_mask`. For information about these transformations and their supporting arguments, see [Transform Images](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/transform_images).
    * `rollout_duration` - (Optional) The amount of time in seconds that it takes for the policy to roll out. During the rollout, the proportion of images with the new policy applied continually increases until cached videos associated with the previous version of the policy are no longer being served. The default is 3600 seconds.
    * `transformations` - (Optional) The set of image transformations applied to the original image. If unspecified, no operations are performed. For information about available transformations and their supporting arguments, see [Transform Images](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/transform_images).
    * `variables` - (Optional) The variable declarations for variables used within the policy. Any variable declared here can be invoked in transformations as a [Variable](#variable) object so that you don't have to specify values separately. You can also pass in these variable names and values dynamically as query parameters in the image's request URL.
        * `name` - (Required) The name of the variable which is also available as the query parameter name to set the variable's value dynamically. Can be up to 50 alphanumeric characters.
        * `type` - (Required) The type of value for the variable.
        * `default_value` - (Required) The default value of the variable if no query parameter is provided. It needs to be one of the `enum_options` if any are provided.
        * `enum_options` - (Optional) Limits the set of possible values for a variable.
            * `id` - (Required) The unique identifier for each enum value. Can be up to 50 alphanumeric characters.
            * `value` - (Required) The value of the variable when the `id` is provided.
        * `postfix` - (Optional) A postfix added to the value provided for the variable, or to the default value.
        * `prefix` - (Optional) A prefix added to the value provided for the variable, or to the default value.

## Attributes reference

This data source returns this attribute:

* `json` - A JSON encoded policy.

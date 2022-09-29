---
layout: akamai
subcategory: Image and Video Manager
---

# akamai_imaging_policy_video (Beta)

Specifies details for a video policy, such as variations in image size and formats.

## Example usage

### Basic usage

This example  shows a very simple video policy that doesn’t contain variables.

```hcl
data "akamai_imaging_policy_video" "video_policy" {
  policy {
    rollout_duration = 3602
    breakpoints {
      widths = [
        1280, 1920
      ]
    }
    output {
      placeholder_video_url = "https://images.im-test.edgesuite.net/specs/im_videos/J7865_DM2360_v.mp4"
    }
  }
}
```

### Usage with variables

This example shows how you can use variables for your policy settings:

```hcl
data "akamai_imaging_policy_video" "video_policy" {
  policy {
    rollout_duration = 3602
    hosts = [
      "example.com", "example_test.com"
    ]
    variables {
      name          = "hosts"
      type          = "bool"
      default_value = "true"
    }
    variables {
      name          = "defaultWidths"
      type          = "number"
      default_value = "320"
    }
    variables {
      name          = "newVar1"
      type          = "perceptualQuality"
      default_value = "mediumHigh"
    }
    output {
        perceptual_quality_var    = "newVar1"
    }
    breakpoints {
      widths = [
        320, 640, 1024, 2048, 5000
      ]
    }
  }
}
```


## Argument reference

This data source supports these arguments:

* `policy` - (Required) The video policy.
    * `breakpoints` - (Optional) The breakpoint widths in pixels used to create derivative videos.
    * `hosts` - (Optional) The hosts that are allowed for video URLs within variables.
    * `output` - (Optional) The output quality of each resized video.
        * `perceptual_quality` - (Optional) The visual quality of derivative videos after they have been compressed to maximize byte savings. You can select one of five perceptual quality values: `high`, `mediumHigh`, `medium`, `mediumLow`, or `low`. `high` results in the highest visual quality with the least byte savings, while ‘low’ results in lower visual quality with the greatest byte savings. If setting a variable for this argument, use `perceptual_quality_var` instead.
        * `placeholder_video_url` - (Optional) The URL for a specific placeholder video that appears when the user first requests a video and Image & Video Manager is still processing the derivative video. If no placeholder video is specified, the original video plays while the derivative video is being processed. If setting a variable for this argument, use `placeholder_video_url_var` instead.
        * `video_adaptive_quality` - (Optional) The quality value that is applied to a video when Image and Video Manager detects a slow connection (RTT > 300 ms). This value overrides the derivative quality value. Specifying a lower value can reduce load times for users with slow connections without impacting the quality of videos for users with standard connections. If setting a variable for this argument, use `video_adaptive_quality_var` instead.
    * `rollout_duration` - (Optional) The amount of time in seconds that it takes for the policy to roll out. During the rollout, the proportion of videos with the new policy applied continually increases until cached videos associated with the previous version of the policy are no longer being served.  
    * `variables` - (Optional) The variables for use within the policy. Any variable declared using this argument can be invoked as a [Variable](#variable) object. You can also pass these variable names and values dynamically as query parameters in the video's request URL.
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

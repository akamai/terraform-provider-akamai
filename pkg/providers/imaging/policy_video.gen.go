package imaging

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// PolicyOutputVideo is a top level schema func
func PolicyOutputVideo(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"breakpoints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The breakpoint widths (in pixels) to use to create derivative images/videos.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: breakpoints(depth - 1),
			},
		},
		"hosts": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Hosts that are allowed for image/video URLs within transformations or variables.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"output": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Dictates the output quality that are created for each resized video.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: outputVideo(depth - 1),
			},
		},
		"rollout_duration": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount of time in seconds that the policy takes to rollout. During the rollout an increasing proportion of images/videos will begin to use the new policy instead of the cached images/videos from the previous version.",
			ValidateDiagFunc: stringAsIntBetween(3600, 604800),
		},
		"variables": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Declares variables for use within the policy. Any variable declared here can be invoked throughout transformations as a [Variable](#variable) object, so that you don't have to specify values separately. You can also pass in these variable names and values dynamically as query parameters in the image's request URL.",
			Elem: &schema.Resource{
				Schema: variable(depth - 1),
			},
		},
	}
}

func outputVideo(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"perceptual_quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The quality of derivative videos. High preserves video quality with reduced byte savings while low reduces video quality to increase byte savings.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"high", "mediumHigh", "medium", "mediumLow", "low"}, false)),
		},
		"perceptual_quality_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The quality of derivative videos. High preserves video quality with reduced byte savings while low reduces video quality to increase byte savings.",
		},
		"placeholder_video_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Allows you to add a specific placeholder video that appears when a user first requests a video, but before Image & Video Manager processes the video. If not specified the original video plays during the processing time.",
		},
		"placeholder_video_url_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Allows you to add a specific placeholder video that appears when a user first requests a video, but before Image & Video Manager processes the video. If not specified the original video plays during the processing time.",
		},
		"video_adaptive_quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Override the quality of video to serve when Image & Video Manager detects a slow connection. Specifying lower values lets users with slow connections browse your site with reduced load times without impacting the quality of videos for users with faster connections.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"high", "mediumHigh", "medium", "mediumLow", "low"}, false)),
		},
		"video_adaptive_quality_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Override the quality of video to serve when Image & Video Manager detects a slow connection. Specifying lower values lets users with slow connections browse your site with reduced load times without impacting the quality of videos for users with faster connections.",
		},
	}
}

package reader

import "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"

// PolicyVideoFromEdgeGrid converts EdgeGrid structure into map terraform-based structure
func PolicyVideoFromEdgeGrid(input imaging.PolicyOutputVideo) map[string]interface{} {
	target := make(map[string]interface{})
	target["breakpoints"] = getBreakpoints(input.Breakpoints)
	target["hosts"] = input.Hosts
	target["output"] = getOutputVideo(input.Output)
	target["variables"] = getVariableList(input.Variables)
	return target
}

func getOutputVideo(src *imaging.OutputVideo) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.PerceptualQuality != nil {
		elem["perceptual_quality"] = src.PerceptualQuality.Value
	}
	if src.PerceptualQuality != nil {
		elem["perceptual_quality_var"] = src.PerceptualQuality.Name
	}
	if src.PlaceholderVideoURL != nil {
		elem["placeholder_video_url"] = src.PlaceholderVideoURL.Value
	}
	if src.PlaceholderVideoURL != nil {
		elem["placeholder_video_url_var"] = src.PlaceholderVideoURL.Name
	}
	if src.VideoAdaptiveQuality != nil {
		elem["video_adaptive_quality"] = src.VideoAdaptiveQuality.Value
	}
	if src.VideoAdaptiveQuality != nil {
		elem["video_adaptive_quality_var"] = src.VideoAdaptiveQuality.Name
	}
	res = append(res, elem)

	return res
}

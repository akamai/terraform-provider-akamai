package reader

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
)

// GetImageSchema converts EdgeGrid structure into map terraform-based structure
func GetImageSchema(input imaging.PolicyOutputImage) map[string]interface{} {
	target := make(map[string]interface{})
	target["breakpoints"] = getBreakpoints(input.Breakpoints)
	target["hosts"] = input.Hosts
	target["output"] = getOutputImage(input.Output)
	target["post_breakpoint_transformations"] = getTransformations(input.PostBreakpointTransformations)
	target["transformations"] = getTransformations(input.Transformations)
	target["variables"] = getVariableList(input.Variables)
	return target
}

func getAppend(src *imaging.Append) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.GravityPriority != nil {
		elem["gravity_priority"] = src.GravityPriority.Value
	}
	if src.GravityPriority != nil {
		elem["gravity_priority_var"] = src.GravityPriority.Name
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	elem["image"] = getImageType(src.Image)
	if src.PreserveMinorDimension != nil {
		elem["preserve_minor_dimension"] = src.PreserveMinorDimension.Value
	}
	if src.PreserveMinorDimension != nil {
		elem["preserve_minor_dimension_var"] = src.PreserveMinorDimension.Name
	}
	res = append(res, elem)

	return res
}

func getAspectCrop(src *imaging.AspectCrop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.AllowExpansion != nil {
		elem["allow_expansion"] = src.AllowExpansion.Value
	}
	if src.AllowExpansion != nil {
		elem["allow_expansion_var"] = src.AllowExpansion.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	if src.XPosition != nil {
		elem["x_position"] = src.XPosition.Value
	}
	if src.XPosition != nil {
		elem["x_position_var"] = src.XPosition.Name
	}
	if src.YPosition != nil {
		elem["y_position"] = src.YPosition.Value
	}
	if src.YPosition != nil {
		elem["y_position_var"] = src.YPosition.Name
	}
	res = append(res, elem)

	return res
}

func getBackgroundColor(src *imaging.BackgroundColor) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Color != nil {
		elem["color"] = src.Color.Value
	}
	if src.Color != nil {
		elem["color_var"] = src.Color.Name
	}
	res = append(res, elem)

	return res
}

func getBlur(src *imaging.Blur) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Sigma != nil {
		elem["sigma"] = src.Sigma.Value
	}
	if src.Sigma != nil {
		elem["sigma_var"] = src.Sigma.Name
	}
	res = append(res, elem)

	return res
}

func getBoxImageType(src *imaging.BoxImageType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Color != nil {
		elem["color"] = src.Color.Value
	}
	if src.Color != nil {
		elem["color_var"] = src.Color.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	elem["transformation"] = getTransformationType(src.Transformation)
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getBreakpoints(src *imaging.Breakpoints) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["widths"] = src.Widths
	res = append(res, elem)

	return res
}

func getChromaKey(src *imaging.ChromaKey) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Hue != nil {
		elem["hue"] = src.Hue.Value
	}
	if src.HueFeather != nil {
		elem["hue_feather"] = src.HueFeather.Value
	}
	if src.HueFeather != nil {
		elem["hue_feather_var"] = src.HueFeather.Name
	}
	if src.HueTolerance != nil {
		elem["hue_tolerance"] = src.HueTolerance.Value
	}
	if src.HueTolerance != nil {
		elem["hue_tolerance_var"] = src.HueTolerance.Name
	}
	if src.Hue != nil {
		elem["hue_var"] = src.Hue.Name
	}
	if src.LightnessFeather != nil {
		elem["lightness_feather"] = src.LightnessFeather.Value
	}
	if src.LightnessFeather != nil {
		elem["lightness_feather_var"] = src.LightnessFeather.Name
	}
	if src.LightnessTolerance != nil {
		elem["lightness_tolerance"] = src.LightnessTolerance.Value
	}
	if src.LightnessTolerance != nil {
		elem["lightness_tolerance_var"] = src.LightnessTolerance.Name
	}
	if src.SaturationFeather != nil {
		elem["saturation_feather"] = src.SaturationFeather.Value
	}
	if src.SaturationFeather != nil {
		elem["saturation_feather_var"] = src.SaturationFeather.Name
	}
	if src.SaturationTolerance != nil {
		elem["saturation_tolerance"] = src.SaturationTolerance.Value
	}
	if src.SaturationTolerance != nil {
		elem["saturation_tolerance_var"] = src.SaturationTolerance.Name
	}
	res = append(res, elem)

	return res
}

func getCircleImageType(src *imaging.CircleImageType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Color != nil {
		elem["color"] = src.Color.Value
	}
	if src.Color != nil {
		elem["color_var"] = src.Color.Name
	}
	if src.Diameter != nil {
		elem["diameter"] = src.Diameter.Value
	}
	if src.Diameter != nil {
		elem["diameter_var"] = src.Diameter.Name
	}
	elem["transformation"] = getTransformationType(src.Transformation)
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getCircleShapeType(src *imaging.CircleShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["center"] = getPointShapeType(src.Center)
	if src.Radius != nil {
		elem["radius"] = src.Radius.Value
	}
	if src.Radius != nil {
		elem["radius_var"] = src.Radius.Name
	}
	res = append(res, elem)

	return res
}

func getComposite(src *imaging.Composite) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	elem["image"] = getImageType(src.Image)
	if src.Placement != nil {
		elem["placement"] = src.Placement.Value
	}
	if src.Placement != nil {
		elem["placement_var"] = src.Placement.Name
	}
	if src.Scale != nil {
		elem["scale"] = src.Scale.Value
	}
	if src.ScaleDimension != nil {
		elem["scale_dimension"] = src.ScaleDimension.Value
	}
	if src.ScaleDimension != nil {
		elem["scale_dimension_var"] = src.ScaleDimension.Name
	}
	if src.Scale != nil {
		elem["scale_var"] = src.Scale.Name
	}
	if src.XPosition != nil {
		elem["x_position"] = src.XPosition.Value
	}
	if src.XPosition != nil {
		elem["x_position_var"] = src.XPosition.Name
	}
	if src.YPosition != nil {
		elem["y_position"] = src.YPosition.Value
	}
	if src.YPosition != nil {
		elem["y_position_var"] = src.YPosition.Name
	}
	res = append(res, elem)

	return res
}

func getCompound(src *imaging.Compound) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["transformations"] = getTransformations(src.Transformations)
	res = append(res, elem)

	return res
}

func getContrast(src *imaging.Contrast) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Brightness != nil {
		elem["brightness"] = src.Brightness.Value
	}
	if src.Brightness != nil {
		elem["brightness_var"] = src.Brightness.Name
	}
	if src.Contrast != nil {
		elem["contrast"] = src.Contrast.Value
	}
	if src.Contrast != nil {
		elem["contrast_var"] = src.Contrast.Name
	}
	res = append(res, elem)

	return res
}

func getCrop(src *imaging.Crop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.AllowExpansion != nil {
		elem["allow_expansion"] = src.AllowExpansion.Value
	}
	if src.AllowExpansion != nil {
		elem["allow_expansion_var"] = src.AllowExpansion.Name
	}
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	if src.XPosition != nil {
		elem["x_position"] = src.XPosition.Value
	}
	if src.XPosition != nil {
		elem["x_position_var"] = src.XPosition.Name
	}
	if src.YPosition != nil {
		elem["y_position"] = src.YPosition.Value
	}
	if src.YPosition != nil {
		elem["y_position_var"] = src.YPosition.Name
	}
	res = append(res, elem)

	return res
}

func getEnumOptionsList(src []*imaging.EnumOptions) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	for _, item := range src {
		elem := make(map[string]interface{})
		elem["id"] = item.ID
		elem["value"] = item.Value
		res = append(res, elem)
	}

	return res
}

func getFaceCrop(src *imaging.FaceCrop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Algorithm != nil {
		elem["algorithm"] = src.Algorithm.Value
	}
	if src.Algorithm != nil {
		elem["algorithm_var"] = src.Algorithm.Name
	}
	if src.Confidence != nil {
		elem["confidence"] = src.Confidence.Value
	}
	if src.Confidence != nil {
		elem["confidence_var"] = src.Confidence.Name
	}
	if src.FailGravity != nil {
		elem["fail_gravity"] = src.FailGravity.Value
	}
	if src.FailGravity != nil {
		elem["fail_gravity_var"] = src.FailGravity.Name
	}
	if src.Focus != nil {
		elem["focus"] = src.Focus.Value
	}
	if src.Focus != nil {
		elem["focus_var"] = src.Focus.Name
	}
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Padding != nil {
		elem["padding"] = src.Padding.Value
	}
	if src.Padding != nil {
		elem["padding_var"] = src.Padding.Name
	}
	if src.Style != nil {
		elem["style"] = src.Style.Value
	}
	if src.Style != nil {
		elem["style_var"] = src.Style.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getFeatureCrop(src *imaging.FeatureCrop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.FailGravity != nil {
		elem["fail_gravity"] = src.FailGravity.Value
	}
	if src.FailGravity != nil {
		elem["fail_gravity_var"] = src.FailGravity.Name
	}
	if src.FeatureRadius != nil {
		elem["feature_radius"] = src.FeatureRadius.Value
	}
	if src.FeatureRadius != nil {
		elem["feature_radius_var"] = src.FeatureRadius.Name
	}
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.MaxFeatures != nil {
		elem["max_features"] = src.MaxFeatures.Value
	}
	if src.MaxFeatures != nil {
		elem["max_features_var"] = src.MaxFeatures.Name
	}
	if src.MinFeatureQuality != nil {
		elem["min_feature_quality"] = src.MinFeatureQuality.Value
	}
	if src.MinFeatureQuality != nil {
		elem["min_feature_quality_var"] = src.MinFeatureQuality.Name
	}
	if src.Padding != nil {
		elem["padding"] = src.Padding.Value
	}
	if src.Padding != nil {
		elem["padding_var"] = src.Padding.Name
	}
	if src.Style != nil {
		elem["style"] = src.Style.Value
	}
	if src.Style != nil {
		elem["style_var"] = src.Style.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getFitAndFill(src *imaging.FitAndFill) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["fill_transformation"] = getTransformationType(src.FillTransformation)
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getGoop(src *imaging.Goop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Chaos != nil {
		elem["chaos"] = src.Chaos.Value
	}
	if src.Chaos != nil {
		elem["chaos_var"] = src.Chaos.Name
	}
	if src.Density != nil {
		elem["density"] = src.Density.Value
	}
	if src.Density != nil {
		elem["density_var"] = src.Density.Name
	}
	if src.Power != nil {
		elem["power"] = src.Power.Value
	}
	if src.Power != nil {
		elem["power_var"] = src.Power.Name
	}
	if src.Seed != nil {
		elem["seed"] = src.Seed.Value
	}
	if src.Seed != nil {
		elem["seed_var"] = src.Seed.Name
	}
	res = append(res, elem)

	return res
}

func getGrayscale(src *imaging.Grayscale) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Type != nil {
		elem["type"] = src.Type.Value
	}
	if src.Type != nil {
		elem["type_var"] = src.Type.Name
	}
	res = append(res, elem)

	return res
}

func getHSL(src *imaging.HSL) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Hue != nil {
		elem["hue"] = src.Hue.Value
	}
	if src.Hue != nil {
		elem["hue_var"] = src.Hue.Name
	}
	if src.Lightness != nil {
		elem["lightness"] = src.Lightness.Value
	}
	if src.Lightness != nil {
		elem["lightness_var"] = src.Lightness.Name
	}
	if src.Saturation != nil {
		elem["saturation"] = src.Saturation.Value
	}
	if src.Saturation != nil {
		elem["saturation_var"] = src.Saturation.Name
	}
	res = append(res, elem)

	return res
}

func getHSV(src *imaging.HSV) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Hue != nil {
		elem["hue"] = src.Hue.Value
	}
	if src.Hue != nil {
		elem["hue_var"] = src.Hue.Name
	}
	if src.Saturation != nil {
		elem["saturation"] = src.Saturation.Value
	}
	if src.Saturation != nil {
		elem["saturation_var"] = src.Saturation.Name
	}
	if src.Value != nil {
		elem["value"] = src.Value.Value
	}
	if src.Value != nil {
		elem["value_var"] = src.Value.Name
	}
	res = append(res, elem)

	return res
}

func getIfDimension(src *imaging.IfDimension) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["default"] = getTransformationType(src.Default)
	if src.Dimension != nil {
		elem["dimension"] = src.Dimension.Value
	}
	if src.Dimension != nil {
		elem["dimension_var"] = src.Dimension.Name
	}
	elem["equal"] = getTransformationType(src.Equal)
	elem["greater_than"] = getTransformationType(src.GreaterThan)
	elem["less_than"] = getTransformationType(src.LessThan)
	if src.Value != nil {
		elem["value"] = src.Value.Value
	}
	if src.Value != nil {
		elem["value_var"] = src.Value.Name
	}
	res = append(res, elem)

	return res
}

func getIfOrientation(src *imaging.IfOrientation) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["default"] = getTransformationType(src.Default)
	elem["landscape"] = getTransformationType(src.Landscape)
	elem["portrait"] = getTransformationType(src.Portrait)
	elem["square"] = getTransformationType(src.Square)
	res = append(res, elem)

	return res
}

func getMaxColors(src *imaging.MaxColors) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Colors != nil {
		elem["colors"] = src.Colors.Value
	}
	if src.Colors != nil {
		elem["colors_var"] = src.Colors.Name
	}
	res = append(res, elem)

	return res
}

func getMirror(src *imaging.Mirror) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Horizontal != nil {
		elem["horizontal"] = src.Horizontal.Value
	}
	if src.Horizontal != nil {
		elem["horizontal_var"] = src.Horizontal.Name
	}
	if src.Vertical != nil {
		elem["vertical"] = src.Vertical.Value
	}
	if src.Vertical != nil {
		elem["vertical_var"] = src.Vertical.Name
	}
	res = append(res, elem)

	return res
}

func getMonoHue(src *imaging.MonoHue) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Hue != nil {
		elem["hue"] = src.Hue.Value
	}
	if src.Hue != nil {
		elem["hue_var"] = src.Hue.Name
	}
	res = append(res, elem)

	return res
}

func getOpacity(src *imaging.Opacity) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Opacity != nil {
		elem["opacity"] = src.Opacity.Value
	}
	if src.Opacity != nil {
		elem["opacity_var"] = src.Opacity.Name
	}
	res = append(res, elem)

	return res
}

func getOutputImage(src *imaging.OutputImage) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["adaptive_quality"] = src.AdaptiveQuality
	if src.PerceptualQuality != nil {
		elem["perceptual_quality"] = src.PerceptualQuality.Value
	}
	elem["perceptual_quality_floor"] = src.PerceptualQualityFloor
	if src.PerceptualQuality != nil {
		elem["perceptual_quality_var"] = src.PerceptualQuality.Name
	}
	if src.Quality != nil {
		elem["quality"] = src.Quality.Value
	}
	if src.Quality != nil {
		elem["quality_var"] = src.Quality.Name
	}
	res = append(res, elem)

	return res
}

func getPointShapeType(src *imaging.PointShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.X != nil {
		elem["x"] = src.X.Value
	}
	if src.X != nil {
		elem["x_var"] = src.X.Name
	}
	if src.Y != nil {
		elem["y"] = src.Y.Value
	}
	if src.Y != nil {
		elem["y_var"] = src.Y.Name
	}
	res = append(res, elem)

	return res
}

func getPointShapeTypeList(src []imaging.PointShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	for _, item := range src {
		elem := make(map[string]interface{})
		if item.X != nil {
			elem["x"] = item.X.Value
		}
		if item.X != nil {
			elem["x_var"] = item.X.Name
		}
		if item.Y != nil {
			elem["y"] = item.Y.Value
		}
		if item.Y != nil {
			elem["y_var"] = item.Y.Name
		}
		res = append(res, elem)
	}

	return res
}

func getPolygonShapeType(src *imaging.PolygonShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["points"] = getPointShapeTypeList(src.Points)
	res = append(res, elem)

	return res
}

func getRectangleShapeType(src *imaging.RectangleShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["anchor"] = getPointShapeType(src.Anchor)
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getRegionOfInterestCrop(src *imaging.RegionOfInterestCrop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Gravity != nil {
		elem["gravity"] = src.Gravity.Value
	}
	if src.Gravity != nil {
		elem["gravity_var"] = src.Gravity.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	elem["region_of_interest"] = getShapeType(src.RegionOfInterest)
	if src.Style != nil {
		elem["style"] = src.Style.Value
	}
	if src.Style != nil {
		elem["style_var"] = src.Style.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getRelativeCrop(src *imaging.RelativeCrop) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.East != nil {
		elem["east"] = src.East.Value
	}
	if src.East != nil {
		elem["east_var"] = src.East.Name
	}
	if src.North != nil {
		elem["north"] = src.North.Value
	}
	if src.North != nil {
		elem["north_var"] = src.North.Name
	}
	if src.South != nil {
		elem["south"] = src.South.Value
	}
	if src.South != nil {
		elem["south_var"] = src.South.Name
	}
	if src.West != nil {
		elem["west"] = src.West.Value
	}
	if src.West != nil {
		elem["west_var"] = src.West.Name
	}
	res = append(res, elem)

	return res
}

func getRemoveColor(src *imaging.RemoveColor) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Color != nil {
		elem["color"] = src.Color.Value
	}
	if src.Color != nil {
		elem["color_var"] = src.Color.Name
	}
	if src.Feather != nil {
		elem["feather"] = src.Feather.Value
	}
	if src.Feather != nil {
		elem["feather_var"] = src.Feather.Name
	}
	if src.Tolerance != nil {
		elem["tolerance"] = src.Tolerance.Value
	}
	if src.Tolerance != nil {
		elem["tolerance_var"] = src.Tolerance.Name
	}
	res = append(res, elem)

	return res
}

func getResize(src *imaging.Resize) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Aspect != nil {
		elem["aspect"] = src.Aspect.Value
	}
	if src.Aspect != nil {
		elem["aspect_var"] = src.Aspect.Name
	}
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Type != nil {
		elem["type"] = src.Type.Value
	}
	if src.Type != nil {
		elem["type_var"] = src.Type.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getRotate(src *imaging.Rotate) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Degrees != nil {
		elem["degrees"] = src.Degrees.Value
	}
	if src.Degrees != nil {
		elem["degrees_var"] = src.Degrees.Name
	}
	res = append(res, elem)

	return res
}

func getScale(src *imaging.Scale) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Height != nil {
		elem["height"] = src.Height.Value
	}
	if src.Height != nil {
		elem["height_var"] = src.Height.Name
	}
	if src.Width != nil {
		elem["width"] = src.Width.Value
	}
	if src.Width != nil {
		elem["width_var"] = src.Width.Name
	}
	res = append(res, elem)

	return res
}

func getShapeTypeList(src []imaging.ShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	for _, item := range src {
		elem := getShapeType(item)
		if len(elem) > 0 {
			res = append(res, elem[0])
		}
	}

	return res
}

func getShear(src *imaging.Shear) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.XShear != nil {
		elem["x_shear"] = src.XShear.Value
	}
	if src.XShear != nil {
		elem["x_shear_var"] = src.XShear.Name
	}
	if src.YShear != nil {
		elem["y_shear"] = src.YShear.Value
	}
	if src.YShear != nil {
		elem["y_shear_var"] = src.YShear.Name
	}
	res = append(res, elem)

	return res
}

func getTextImageType(src *imaging.TextImageType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Fill != nil {
		elem["fill"] = src.Fill.Value
	}
	if src.Fill != nil {
		elem["fill_var"] = src.Fill.Name
	}
	if src.Size != nil {
		elem["size"] = src.Size.Value
	}
	if src.Size != nil {
		elem["size_var"] = src.Size.Name
	}
	if src.Stroke != nil {
		elem["stroke"] = src.Stroke.Value
	}
	if src.StrokeSize != nil {
		elem["stroke_size"] = src.StrokeSize.Value
	}
	if src.StrokeSize != nil {
		elem["stroke_size_var"] = src.StrokeSize.Name
	}
	if src.Stroke != nil {
		elem["stroke_var"] = src.Stroke.Name
	}
	if src.Text != nil {
		elem["text"] = src.Text.Value
	}
	if src.Text != nil {
		elem["text_var"] = src.Text.Name
	}
	elem["transformation"] = getTransformationType(src.Transformation)
	if src.Typeface != nil {
		elem["typeface"] = src.Typeface.Value
	}
	if src.Typeface != nil {
		elem["typeface_var"] = src.Typeface.Name
	}
	res = append(res, elem)

	return res
}

func getTransformations(src imaging.Transformations) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	for _, item := range src {
		elem := getTransformationType(item)
		if len(elem) > 0 {
			res = append(res, elem[0])
		}
	}

	return res
}

func getTrim(src *imaging.Trim) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Fuzz != nil {
		elem["fuzz"] = src.Fuzz.Value
	}
	if src.Fuzz != nil {
		elem["fuzz_var"] = src.Fuzz.Name
	}
	if src.Padding != nil {
		elem["padding"] = src.Padding.Value
	}
	if src.Padding != nil {
		elem["padding_var"] = src.Padding.Name
	}
	res = append(res, elem)

	return res
}

func getURLImageType(src *imaging.URLImageType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["transformation"] = getTransformationType(src.Transformation)
	if src.URL != nil {
		elem["url"] = src.URL.Value
	}
	if src.URL != nil {
		elem["url_var"] = src.URL.Name
	}
	res = append(res, elem)

	return res
}

func getUnionShapeType(src *imaging.UnionShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	elem["shapes"] = getShapeTypeList(src.Shapes)
	res = append(res, elem)

	return res
}

func getUnsharpMask(src *imaging.UnsharpMask) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	elem := make(map[string]interface{})
	if src.Gain != nil {
		elem["gain"] = src.Gain.Value
	}
	if src.Gain != nil {
		elem["gain_var"] = src.Gain.Name
	}
	if src.Sigma != nil {
		elem["sigma"] = src.Sigma.Value
	}
	if src.Sigma != nil {
		elem["sigma_var"] = src.Sigma.Name
	}
	if src.Threshold != nil {
		elem["threshold"] = src.Threshold.Value
	}
	if src.Threshold != nil {
		elem["threshold_var"] = src.Threshold.Name
	}
	res = append(res, elem)

	return res
}

func getVariableList(src []imaging.Variable) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}
	for _, item := range src {
		elem := make(map[string]interface{})
		elem["default_value"] = item.DefaultValue
		elem["enum_options"] = getEnumOptionsList(item.EnumOptions)
		elem["name"] = item.Name
		elem["postfix"] = item.Postfix
		elem["prefix"] = item.Prefix
		elem["type"] = item.Type
		res = append(res, elem)
	}

	return res
}

func getImageType(src imaging.ImageType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}

	elem := map[string]interface{}{}

	switch t := src.(type) {
	case *imaging.BoxImageType:
		elem["box_image"] = getBoxImageType(t)
	case *imaging.CircleImageType:
		elem["circle_image"] = getCircleImageType(t)
	case *imaging.TextImageType:
		elem["text_image"] = getTextImageType(t)
	case *imaging.URLImageType:
		elem["url_image"] = getURLImageType(t)
	default:
		panic(fmt.Sprintf("unsupported type - %T does not implement imaging.ImageType", t))
	}
	res = append(res, elem)

	return res
}

func getShapeType(src imaging.ShapeType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}

	elem := map[string]interface{}{}

	switch t := src.(type) {
	case *imaging.CircleShapeType:
		elem["circle_shape"] = getCircleShapeType(t)
	case *imaging.PointShapeType:
		elem["point_shape"] = getPointShapeType(t)
	case *imaging.PolygonShapeType:
		elem["polygon_shape"] = getPolygonShapeType(t)
	case *imaging.RectangleShapeType:
		elem["rectangle_shape"] = getRectangleShapeType(t)
	case *imaging.UnionShapeType:
		elem["union_shape"] = getUnionShapeType(t)
	default:
		panic(fmt.Sprintf("unsupported type - %T does not implement imaging.ShapeType", t))
	}
	res = append(res, elem)

	return res
}

func getTransformationType(src imaging.TransformationType) []map[string]interface{} {
	if src == nil {
		return nil
	}

	var res []map[string]interface{}

	elem := map[string]interface{}{}

	switch t := src.(type) {
	case *imaging.Append:
		elem["append"] = getAppend(t)
	case *imaging.AspectCrop:
		elem["aspect_crop"] = getAspectCrop(t)
	case *imaging.BackgroundColor:
		elem["background_color"] = getBackgroundColor(t)
	case *imaging.Blur:
		elem["blur"] = getBlur(t)
	case *imaging.ChromaKey:
		elem["chroma_key"] = getChromaKey(t)
	case *imaging.Composite:
		elem["composite"] = getComposite(t)
	case *imaging.Compound:
		elem["compound"] = getCompound(t)
	case *imaging.Contrast:
		elem["contrast"] = getContrast(t)
	case *imaging.Crop:
		elem["crop"] = getCrop(t)
	case *imaging.FaceCrop:
		elem["face_crop"] = getFaceCrop(t)
	case *imaging.FeatureCrop:
		elem["feature_crop"] = getFeatureCrop(t)
	case *imaging.FitAndFill:
		elem["fit_and_fill"] = getFitAndFill(t)
	case *imaging.Goop:
		elem["goop"] = getGoop(t)
	case *imaging.Grayscale:
		elem["grayscale"] = getGrayscale(t)
	case *imaging.HSL:
		elem["hsl"] = getHSL(t)
	case *imaging.HSV:
		elem["hsv"] = getHSV(t)
	case *imaging.IfDimension:
		elem["if_dimension"] = getIfDimension(t)
	case *imaging.IfOrientation:
		elem["if_orientation"] = getIfOrientation(t)
	case *imaging.MaxColors:
		elem["max_colors"] = getMaxColors(t)
	case *imaging.Mirror:
		elem["mirror"] = getMirror(t)
	case *imaging.MonoHue:
		elem["mono_hue"] = getMonoHue(t)
	case *imaging.Opacity:
		elem["opacity"] = getOpacity(t)
	case *imaging.RegionOfInterestCrop:
		elem["region_of_interest_crop"] = getRegionOfInterestCrop(t)
	case *imaging.RelativeCrop:
		elem["relative_crop"] = getRelativeCrop(t)
	case *imaging.RemoveColor:
		elem["remove_color"] = getRemoveColor(t)
	case *imaging.Resize:
		elem["resize"] = getResize(t)
	case *imaging.Rotate:
		elem["rotate"] = getRotate(t)
	case *imaging.Scale:
		elem["scale"] = getScale(t)
	case *imaging.Shear:
		elem["shear"] = getShear(t)
	case *imaging.Trim:
		elem["trim"] = getTrim(t)
	case *imaging.UnsharpMask:
		elem["unsharp_mask"] = getUnsharpMask(t)
	default:
		panic(fmt.Sprintf("unsupported type - %T does not implement imaging.TransformationType", t))
	}
	res = append(res, elem)

	return res
}

package writer

import (
	"fmt"
	"reflect"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
)

// PolicyImageToEdgeGrid converts terraform originated map structure into EdgeGrid structure
func PolicyImageToEdgeGrid(input map[string]interface{}) imaging.PolicyInputImage {
	result := imaging.PolicyInputImage{}
	result.Breakpoints = getBreakpoints(extract(input, "breakpoints"))
	result.Hosts = interfaceSliceToStringSlice(input["hosts"].([]interface{}))
	result.Output = getOutputImage(extract(input, "output"))
	result.PostBreakpointTransformations = getTransformations(input["post_breakpoint_transformations"].([]interface{}))
	result.RolloutDuration = input["rollout_duration"].(int)
	result.Transformations = getTransformations(input["transformations"].([]interface{}))
	result.Variables = getVariableList(input["variables"].([]interface{}))
	return result
}

func getAppend(src map[string]interface{}) *imaging.Append {
	if src == nil {
		return nil
	}
	result := imaging.Append{
		Gravity:                gravityVariableInline(src, "gravity"),
		GravityPriority:        appendGravityPriorityVariableInline(src, "gravity_priority"),
		Image:                  getImageType(extract(src, "image")),
		PreserveMinorDimension: booleanVariableInline(src, "preserve_minor_dimension"),
		Transformation:         imaging.AppendTransformationAppend,
	}
	return &result
}

func getAspectCrop(src map[string]interface{}) *imaging.AspectCrop {
	if src == nil {
		return nil
	}
	result := imaging.AspectCrop{
		AllowExpansion: booleanVariableInline(src, "allow_expansion"),
		Height:         numberVariableInline(src, "height"),
		Width:          numberVariableInline(src, "width"),
		XPosition:      numberVariableInline(src, "x_position"),
		YPosition:      numberVariableInline(src, "y_position"),
		Transformation: imaging.AspectCropTransformationAspectCrop,
	}
	return &result
}

func getBackgroundColor(src map[string]interface{}) *imaging.BackgroundColor {
	if src == nil {
		return nil
	}
	result := imaging.BackgroundColor{
		Color:          stringVariableInline(src, "color"),
		Transformation: imaging.BackgroundColorTransformationBackgroundColor,
	}
	return &result
}

func getBlur(src map[string]interface{}) *imaging.Blur {
	if src == nil {
		return nil
	}
	result := imaging.Blur{
		Sigma:          numberVariableInline(src, "sigma"),
		Transformation: imaging.BlurTransformationBlur,
	}
	return &result
}

func getBoxImageType(src map[string]interface{}) *imaging.BoxImageType {
	if src == nil {
		return nil
	}
	result := imaging.BoxImageType{
		Color:          stringVariableInline(src, "color"),
		Height:         integerVariableInline(src, "height"),
		Transformation: getTransformationType(extract(src, "transformation")),
		Width:          integerVariableInline(src, "width"),
		Type:           imaging.BoxImageTypeTypeBox,
	}
	return &result
}

func getBreakpoints(src map[string]interface{}) *imaging.Breakpoints {
	if src == nil {
		return nil
	}
	result := imaging.Breakpoints{
		Widths: interfaceSliceToIntSlice(src["widths"].([]interface{})),
	}
	return &result
}

func getChromaKey(src map[string]interface{}) *imaging.ChromaKey {
	if src == nil {
		return nil
	}
	result := imaging.ChromaKey{
		Hue:                 numberVariableInline(src, "hue"),
		HueFeather:          numberVariableInline(src, "hue_feather"),
		HueTolerance:        numberVariableInline(src, "hue_tolerance"),
		LightnessFeather:    numberVariableInline(src, "lightness_feather"),
		LightnessTolerance:  numberVariableInline(src, "lightness_tolerance"),
		SaturationFeather:   numberVariableInline(src, "saturation_feather"),
		SaturationTolerance: numberVariableInline(src, "saturation_tolerance"),
		Transformation:      imaging.ChromaKeyTransformationChromaKey,
	}
	return &result
}

func getCircleImageType(src map[string]interface{}) *imaging.CircleImageType {
	if src == nil {
		return nil
	}
	result := imaging.CircleImageType{
		Color:          stringVariableInline(src, "color"),
		Diameter:       integerVariableInline(src, "diameter"),
		Transformation: getTransformationType(extract(src, "transformation")),
		Width:          integerVariableInline(src, "width"),
		Type:           imaging.CircleImageTypeTypeCircle,
	}
	return &result
}

func getCircleShapeType(src map[string]interface{}) *imaging.CircleShapeType {
	if src == nil {
		return nil
	}
	result := imaging.CircleShapeType{
		Center: getPointShapeType(extract(src, "center")),
		Radius: numberVariableInline(src, "radius"),
	}
	return &result
}

func getComposite(src map[string]interface{}) *imaging.Composite {
	if src == nil {
		return nil
	}
	result := imaging.Composite{
		Gravity:        gravityVariableInline(src, "gravity"),
		Image:          getImageType(extract(src, "image")),
		Placement:      compositePlacementVariableInline(src, "placement"),
		Scale:          numberVariableInline(src, "scale"),
		ScaleDimension: compositeScaleDimensionVariableInline(src, "scale_dimension"),
		XPosition:      integerVariableInline(src, "x_position"),
		YPosition:      integerVariableInline(src, "y_position"),
		Transformation: imaging.CompositeTransformationComposite,
	}
	return &result
}

func getCompound(src map[string]interface{}) *imaging.Compound {
	if src == nil {
		return nil
	}
	result := imaging.Compound{
		Transformations: getTransformations(src["transformations"].([]interface{})),
		Transformation:  imaging.CompoundTransformationCompound,
	}
	return &result
}

func getContrast(src map[string]interface{}) *imaging.Contrast {
	if src == nil {
		return nil
	}
	result := imaging.Contrast{
		Brightness:     numberVariableInline(src, "brightness"),
		Contrast:       numberVariableInline(src, "contrast"),
		Transformation: imaging.ContrastTransformationContrast,
	}
	return &result
}

func getCrop(src map[string]interface{}) *imaging.Crop {
	if src == nil {
		return nil
	}
	result := imaging.Crop{
		AllowExpansion: booleanVariableInline(src, "allow_expansion"),
		Gravity:        gravityVariableInline(src, "gravity"),
		Height:         integerVariableInline(src, "height"),
		Width:          integerVariableInline(src, "width"),
		XPosition:      integerVariableInline(src, "x_position"),
		YPosition:      integerVariableInline(src, "y_position"),
		Transformation: imaging.CropTransformationCrop,
	}
	return &result
}

func getEnumOptionsList(src []interface{}) []*imaging.EnumOptions {
	result := make([]*imaging.EnumOptions, 0)
	for idx := range src {
		elem := imaging.EnumOptions{
			ID:    src[idx].(map[string]interface{})["id"].(string),
			Value: src[idx].(map[string]interface{})["value"].(string),
		}
		result = append(result, &elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getFaceCrop(src map[string]interface{}) *imaging.FaceCrop {
	if src == nil {
		return nil
	}
	result := imaging.FaceCrop{
		Algorithm:      faceCropAlgorithmVariableInline(src, "algorithm"),
		Confidence:     numberVariableInline(src, "confidence"),
		FailGravity:    gravityVariableInline(src, "fail_gravity"),
		Focus:          faceCropFocusVariableInline(src, "focus"),
		Gravity:        gravityVariableInline(src, "gravity"),
		Height:         integerVariableInline(src, "height"),
		Padding:        numberVariableInline(src, "padding"),
		Style:          faceCropStyleVariableInline(src, "style"),
		Width:          integerVariableInline(src, "width"),
		Transformation: imaging.FaceCropTransformationFaceCrop,
	}
	return &result
}

func getFeatureCrop(src map[string]interface{}) *imaging.FeatureCrop {
	if src == nil {
		return nil
	}
	result := imaging.FeatureCrop{
		FailGravity:       gravityVariableInline(src, "fail_gravity"),
		FeatureRadius:     numberVariableInline(src, "feature_radius"),
		Gravity:           gravityVariableInline(src, "gravity"),
		Height:            integerVariableInline(src, "height"),
		MaxFeatures:       integerVariableInline(src, "max_features"),
		MinFeatureQuality: numberVariableInline(src, "min_feature_quality"),
		Padding:           numberVariableInline(src, "padding"),
		Style:             featureCropStyleVariableInline(src, "style"),
		Width:             integerVariableInline(src, "width"),
		Transformation:    imaging.FeatureCropTransformationFeatureCrop,
	}
	return &result
}

func getFitAndFill(src map[string]interface{}) *imaging.FitAndFill {
	if src == nil {
		return nil
	}
	result := imaging.FitAndFill{
		FillTransformation: getTransformationType(extract(src, "fill_transformation")),
		Height:             integerVariableInline(src, "height"),
		Width:              integerVariableInline(src, "width"),
		Transformation:     imaging.FitAndFillTransformationFitAndFill,
	}
	return &result
}

func getGoop(src map[string]interface{}) *imaging.Goop {
	if src == nil {
		return nil
	}
	result := imaging.Goop{
		Chaos:          numberVariableInline(src, "chaos"),
		Density:        integerVariableInline(src, "density"),
		Power:          numberVariableInline(src, "power"),
		Seed:           integerVariableInline(src, "seed"),
		Transformation: imaging.GoopTransformationGoop,
	}
	return &result
}

func getGrayscale(src map[string]interface{}) *imaging.Grayscale {
	if src == nil {
		return nil
	}
	result := imaging.Grayscale{
		Type:           grayscaleTypeVariableInline(src, "type"),
		Transformation: imaging.GrayscaleTransformationGrayscale,
	}
	return &result
}

func getHSL(src map[string]interface{}) *imaging.HSL {
	if src == nil {
		return nil
	}
	result := imaging.HSL{
		Hue:            numberVariableInline(src, "hue"),
		Lightness:      numberVariableInline(src, "lightness"),
		Saturation:     numberVariableInline(src, "saturation"),
		Transformation: imaging.HSLTransformationHSL,
	}
	return &result
}

func getHSV(src map[string]interface{}) *imaging.HSV {
	if src == nil {
		return nil
	}
	result := imaging.HSV{
		Hue:            numberVariableInline(src, "hue"),
		Saturation:     numberVariableInline(src, "saturation"),
		Value:          numberVariableInline(src, "value"),
		Transformation: imaging.HSVTransformationHSV,
	}
	return &result
}

func getIfDimension(src map[string]interface{}) *imaging.IfDimension {
	if src == nil {
		return nil
	}
	result := imaging.IfDimension{
		Default:        getTransformationType(extract(src, "default")),
		Dimension:      ifDimensionDimensionVariableInline(src, "dimension"),
		Equal:          getTransformationType(extract(src, "equal")),
		GreaterThan:    getTransformationType(extract(src, "greater_than")),
		LessThan:       getTransformationType(extract(src, "less_than")),
		Value:          integerVariableInline(src, "value"),
		Transformation: imaging.IfDimensionTransformationIfDimension,
	}
	return &result
}

func getIfOrientation(src map[string]interface{}) *imaging.IfOrientation {
	if src == nil {
		return nil
	}
	result := imaging.IfOrientation{
		Default:        getTransformationType(extract(src, "default")),
		Landscape:      getTransformationType(extract(src, "landscape")),
		Portrait:       getTransformationType(extract(src, "portrait")),
		Square:         getTransformationType(extract(src, "square")),
		Transformation: imaging.IfOrientationTransformationIfOrientation,
	}
	return &result
}

func getMaxColors(src map[string]interface{}) *imaging.MaxColors {
	if src == nil {
		return nil
	}
	result := imaging.MaxColors{
		Colors:         integerVariableInline(src, "colors"),
		Transformation: imaging.MaxColorsTransformationMaxColors,
	}
	return &result
}

func getMirror(src map[string]interface{}) *imaging.Mirror {
	if src == nil {
		return nil
	}
	result := imaging.Mirror{
		Horizontal:     booleanVariableInline(src, "horizontal"),
		Vertical:       booleanVariableInline(src, "vertical"),
		Transformation: imaging.MirrorTransformationMirror,
	}
	return &result
}

func getMonoHue(src map[string]interface{}) *imaging.MonoHue {
	if src == nil {
		return nil
	}
	result := imaging.MonoHue{
		Hue:            numberVariableInline(src, "hue"),
		Transformation: imaging.MonoHueTransformationMonoHue,
	}
	return &result
}

func getOpacity(src map[string]interface{}) *imaging.Opacity {
	if src == nil {
		return nil
	}
	result := imaging.Opacity{
		Opacity:        numberVariableInline(src, "opacity"),
		Transformation: imaging.OpacityTransformationOpacity,
	}
	return &result
}

func getOutputImage(src map[string]interface{}) *imaging.OutputImage {
	if src == nil {
		return nil
	}
	result := imaging.OutputImage{
		AdaptiveQuality:        src["adaptive_quality"].(int),
		PerceptualQuality:      outputImagePerceptualQualityVariableInline(src, "perceptual_quality"),
		PerceptualQualityFloor: src["perceptual_quality_floor"].(int),
		Quality:                integerVariableInline(src, "quality"),
	}
	return &result
}

func getPointShapeType(src map[string]interface{}) *imaging.PointShapeType {
	if src == nil {
		return nil
	}
	result := imaging.PointShapeType{
		X: numberVariableInline(src, "x"),
		Y: numberVariableInline(src, "y"),
	}
	return &result
}

func getPointShapeTypeList(src []interface{}) []imaging.PointShapeType {
	result := make([]imaging.PointShapeType, 0)
	for idx := range src {
		elem := imaging.PointShapeType{
			X: numberVariableInline(src[idx].(map[string]interface{}), "x"),
			Y: numberVariableInline(src[idx].(map[string]interface{}), "y"),
		}
		result = append(result, elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getPolygonShapeType(src map[string]interface{}) *imaging.PolygonShapeType {
	if src == nil {
		return nil
	}
	result := imaging.PolygonShapeType{
		Points: getPointShapeTypeList(src["points"].([]interface{})),
	}
	return &result
}

func getRectangleShapeType(src map[string]interface{}) *imaging.RectangleShapeType {
	if src == nil {
		return nil
	}
	result := imaging.RectangleShapeType{
		Anchor: getPointShapeType(extract(src, "anchor")),
		Height: numberVariableInline(src, "height"),
		Width:  numberVariableInline(src, "width"),
	}
	return &result
}

func getRegionOfInterestCrop(src map[string]interface{}) *imaging.RegionOfInterestCrop {
	if src == nil {
		return nil
	}
	result := imaging.RegionOfInterestCrop{
		Gravity:          gravityVariableInline(src, "gravity"),
		Height:           integerVariableInline(src, "height"),
		RegionOfInterest: getShapeType(extract(src, "region_of_interest")),
		Style:            regionOfInterestCropStyleVariableInline(src, "style"),
		Width:            integerVariableInline(src, "width"),
		Transformation:   imaging.RegionOfInterestCropTransformationRegionOfInterestCrop,
	}
	return &result
}

func getRelativeCrop(src map[string]interface{}) *imaging.RelativeCrop {
	if src == nil {
		return nil
	}
	result := imaging.RelativeCrop{
		East:           integerVariableInline(src, "east"),
		North:          integerVariableInline(src, "north"),
		South:          integerVariableInline(src, "south"),
		West:           integerVariableInline(src, "west"),
		Transformation: imaging.RelativeCropTransformationRelativeCrop,
	}
	return &result
}

func getRemoveColor(src map[string]interface{}) *imaging.RemoveColor {
	if src == nil {
		return nil
	}
	result := imaging.RemoveColor{
		Color:          stringVariableInline(src, "color"),
		Feather:        numberVariableInline(src, "feather"),
		Tolerance:      numberVariableInline(src, "tolerance"),
		Transformation: imaging.RemoveColorTransformationRemoveColor,
	}
	return &result
}

func getResize(src map[string]interface{}) *imaging.Resize {
	if src == nil {
		return nil
	}
	result := imaging.Resize{
		Aspect:         resizeAspectVariableInline(src, "aspect"),
		Height:         integerVariableInline(src, "height"),
		Type:           resizeTypeVariableInline(src, "type"),
		Width:          integerVariableInline(src, "width"),
		Transformation: imaging.ResizeTransformationResize,
	}
	return &result
}

func getRotate(src map[string]interface{}) *imaging.Rotate {
	if src == nil {
		return nil
	}
	result := imaging.Rotate{
		Degrees:        numberVariableInline(src, "degrees"),
		Transformation: imaging.RotateTransformationRotate,
	}
	return &result
}

func getScale(src map[string]interface{}) *imaging.Scale {
	if src == nil {
		return nil
	}
	result := imaging.Scale{
		Height:         numberVariableInline(src, "height"),
		Width:          numberVariableInline(src, "width"),
		Transformation: imaging.ScaleTransformationScale,
	}
	return &result
}

func getShapeTypeList(src []interface{}) []imaging.ShapeType {
	result := make([]imaging.ShapeType, 0)
	for idx := range src {
		elem := getShapeType(src[idx].(map[string]interface{}))
		result = append(result, elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getShear(src map[string]interface{}) *imaging.Shear {
	if src == nil {
		return nil
	}
	result := imaging.Shear{
		XShear:         numberVariableInline(src, "x_shear"),
		YShear:         numberVariableInline(src, "y_shear"),
		Transformation: imaging.ShearTransformationShear,
	}
	return &result
}

func getTextImageType(src map[string]interface{}) *imaging.TextImageType {
	if src == nil {
		return nil
	}
	result := imaging.TextImageType{
		Fill:           stringVariableInline(src, "fill"),
		Size:           numberVariableInline(src, "size"),
		Stroke:         stringVariableInline(src, "stroke"),
		StrokeSize:     numberVariableInline(src, "stroke_size"),
		Text:           stringVariableInline(src, "text"),
		Transformation: getTransformationType(extract(src, "transformation")),
		Typeface:       stringVariableInline(src, "typeface"),
		Type:           imaging.TextImageTypeTypeText,
	}
	return &result
}

func getTransformations(src []interface{}) []imaging.TransformationType {
	result := make([]imaging.TransformationType, 0)
	for idx := range src {
		elem := getTransformationType(src[idx].(map[string]interface{}))
		result = append(result, elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getTrim(src map[string]interface{}) *imaging.Trim {
	if src == nil {
		return nil
	}
	result := imaging.Trim{
		Fuzz:           numberVariableInline(src, "fuzz"),
		Padding:        integerVariableInline(src, "padding"),
		Transformation: imaging.TrimTransformationTrim,
	}
	return &result
}

func getURLImageType(src map[string]interface{}) *imaging.URLImageType {
	if src == nil {
		return nil
	}
	result := imaging.URLImageType{
		Transformation: getTransformationType(extract(src, "transformation")),
		URL:            stringVariableInline(src, "url"),
		Type:           imaging.URLImageTypeTypeURL,
	}
	return &result
}

func getUnionShapeType(src map[string]interface{}) *imaging.UnionShapeType {
	if src == nil {
		return nil
	}
	result := imaging.UnionShapeType{
		Shapes: getShapeTypeList(src["shapes"].([]interface{})),
	}
	return &result
}

func getUnsharpMask(src map[string]interface{}) *imaging.UnsharpMask {
	if src == nil {
		return nil
	}
	result := imaging.UnsharpMask{
		Gain:           numberVariableInline(src, "gain"),
		Sigma:          numberVariableInline(src, "sigma"),
		Threshold:      numberVariableInline(src, "threshold"),
		Transformation: imaging.UnsharpMaskTransformationUnsharpMask,
	}
	return &result
}

func getVariableList(src []interface{}) []imaging.Variable {
	result := make([]imaging.Variable, 0)
	for idx := range src {
		elem := imaging.Variable{
			DefaultValue: src[idx].(map[string]interface{})["default_value"].(string),
			EnumOptions:  getEnumOptionsList(src[idx].(map[string]interface{})["enum_options"].([]interface{})),
			Name:         src[idx].(map[string]interface{})["name"].(string),
			Postfix:      src[idx].(map[string]interface{})["postfix"].(string),
			Prefix:       src[idx].(map[string]interface{})["prefix"].(string),
			Type:         imaging.VariableType(src[idx].(map[string]interface{})["type"].(string)),
		}
		result = append(result, elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getImageType(src map[string]interface{}) imaging.ImageType {
	if src == nil {
		return nil
	}

	var node map[string]interface{}
	node = extract(src, "box_image")
	if node != nil {
		return getBoxImageType(node)
	}
	node = extract(src, "circle_image")
	if node != nil {
		return getCircleImageType(node)
	}
	node = extract(src, "text_image")
	if node != nil {
		return getTextImageType(node)
	}
	node = extract(src, "url_image")
	if node != nil {
		return getURLImageType(node)
	}
	panic(fmt.Sprint("unsupported type"))
}

func getShapeType(src map[string]interface{}) imaging.ShapeType {
	if src == nil {
		return nil
	}

	var node map[string]interface{}
	node = extract(src, "circle_shape")
	if node != nil {
		return getCircleShapeType(node)
	}
	node = extract(src, "point_shape")
	if node != nil {
		return getPointShapeType(node)
	}
	node = extract(src, "polygon_shape")
	if node != nil {
		return getPolygonShapeType(node)
	}
	node = extract(src, "rectangle_shape")
	if node != nil {
		return getRectangleShapeType(node)
	}
	node = extract(src, "union_shape")
	if node != nil {
		return getUnionShapeType(node)
	}
	panic(fmt.Sprint("unsupported type"))
}

func getTransformationType(src map[string]interface{}) imaging.TransformationType {
	if src == nil {
		return nil
	}

	var node map[string]interface{}
	node = extract(src, "append")
	if node != nil {
		return getAppend(node)
	}
	node = extract(src, "aspect_crop")
	if node != nil {
		return getAspectCrop(node)
	}
	node = extract(src, "background_color")
	if node != nil {
		return getBackgroundColor(node)
	}
	node = extract(src, "blur")
	if node != nil {
		return getBlur(node)
	}
	node = extract(src, "chroma_key")
	if node != nil {
		return getChromaKey(node)
	}
	node = extract(src, "composite")
	if node != nil {
		return getComposite(node)
	}
	node = extract(src, "compound")
	if node != nil {
		return getCompound(node)
	}
	node = extract(src, "contrast")
	if node != nil {
		return getContrast(node)
	}
	node = extract(src, "crop")
	if node != nil {
		return getCrop(node)
	}
	node = extract(src, "face_crop")
	if node != nil {
		return getFaceCrop(node)
	}
	node = extract(src, "feature_crop")
	if node != nil {
		return getFeatureCrop(node)
	}
	node = extract(src, "fit_and_fill")
	if node != nil {
		return getFitAndFill(node)
	}
	node = extract(src, "goop")
	if node != nil {
		return getGoop(node)
	}
	node = extract(src, "grayscale")
	if node != nil {
		return getGrayscale(node)
	}
	node = extract(src, "hsl")
	if node != nil {
		return getHSL(node)
	}
	node = extract(src, "hsv")
	if node != nil {
		return getHSV(node)
	}
	node = extract(src, "if_dimension")
	if node != nil {
		return getIfDimension(node)
	}
	node = extract(src, "if_orientation")
	if node != nil {
		return getIfOrientation(node)
	}
	node = extract(src, "max_colors")
	if node != nil {
		return getMaxColors(node)
	}
	node = extract(src, "mirror")
	if node != nil {
		return getMirror(node)
	}
	node = extract(src, "mono_hue")
	if node != nil {
		return getMonoHue(node)
	}
	node = extract(src, "opacity")
	if node != nil {
		return getOpacity(node)
	}
	node = extract(src, "region_of_interest_crop")
	if node != nil {
		return getRegionOfInterestCrop(node)
	}
	node = extract(src, "relative_crop")
	if node != nil {
		return getRelativeCrop(node)
	}
	node = extract(src, "remove_color")
	if node != nil {
		return getRemoveColor(node)
	}
	node = extract(src, "resize")
	if node != nil {
		return getResize(node)
	}
	node = extract(src, "rotate")
	if node != nil {
		return getRotate(node)
	}
	node = extract(src, "scale")
	if node != nil {
		return getScale(node)
	}
	node = extract(src, "shear")
	if node != nil {
		return getShear(node)
	}
	node = extract(src, "trim")
	if node != nil {
		return getTrim(node)
	}
	node = extract(src, "unsharp_mask")
	if node != nil {
		return getUnsharpMask(node)
	}
	panic(fmt.Sprint("unsupported type"))
}

func appendGravityPriorityVariableInline(src map[string]interface{}, name string) *imaging.AppendGravityPriorityVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.AppendGravityPriority
	if v1 != "" {
		v2 = imaging.AppendGravityPriorityPtr(imaging.AppendGravityPriority(v1.(string)))
	}

	return &imaging.AppendGravityPriorityVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func compositePlacementVariableInline(src map[string]interface{}, name string) *imaging.CompositePlacementVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.CompositePlacement
	if v1 != "" {
		v2 = imaging.CompositePlacementPtr(imaging.CompositePlacement(v1.(string)))
	}

	return &imaging.CompositePlacementVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func compositeScaleDimensionVariableInline(src map[string]interface{}, name string) *imaging.CompositeScaleDimensionVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.CompositeScaleDimension
	if v1 != "" {
		v2 = imaging.CompositeScaleDimensionPtr(imaging.CompositeScaleDimension(v1.(string)))
	}

	return &imaging.CompositeScaleDimensionVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func faceCropAlgorithmVariableInline(src map[string]interface{}, name string) *imaging.FaceCropAlgorithmVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.FaceCropAlgorithm
	if v1 != "" {
		v2 = imaging.FaceCropAlgorithmPtr(imaging.FaceCropAlgorithm(v1.(string)))
	}

	return &imaging.FaceCropAlgorithmVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func faceCropFocusVariableInline(src map[string]interface{}, name string) *imaging.FaceCropFocusVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.FaceCropFocus
	if v1 != "" {
		v2 = imaging.FaceCropFocusPtr(imaging.FaceCropFocus(v1.(string)))
	}

	return &imaging.FaceCropFocusVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func faceCropStyleVariableInline(src map[string]interface{}, name string) *imaging.FaceCropStyleVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.FaceCropStyle
	if v1 != "" {
		v2 = imaging.FaceCropStylePtr(imaging.FaceCropStyle(v1.(string)))
	}

	return &imaging.FaceCropStyleVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func featureCropStyleVariableInline(src map[string]interface{}, name string) *imaging.FeatureCropStyleVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.FeatureCropStyle
	if v1 != "" {
		v2 = imaging.FeatureCropStylePtr(imaging.FeatureCropStyle(v1.(string)))
	}

	return &imaging.FeatureCropStyleVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func gravityVariableInline(src map[string]interface{}, name string) *imaging.GravityVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.Gravity
	if v1 != "" {
		v2 = imaging.GravityPtr(imaging.Gravity(v1.(string)))
	}

	return &imaging.GravityVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func grayscaleTypeVariableInline(src map[string]interface{}, name string) *imaging.GrayscaleTypeVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.GrayscaleType
	if v1 != "" {
		v2 = imaging.GrayscaleTypePtr(imaging.GrayscaleType(v1.(string)))
	}

	return &imaging.GrayscaleTypeVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func ifDimensionDimensionVariableInline(src map[string]interface{}, name string) *imaging.IfDimensionDimensionVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.IfDimensionDimension
	if v1 != "" {
		v2 = imaging.IfDimensionDimensionPtr(imaging.IfDimensionDimension(v1.(string)))
	}

	return &imaging.IfDimensionDimensionVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func outputImagePerceptualQualityVariableInline(src map[string]interface{}, name string) *imaging.OutputImagePerceptualQualityVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.OutputImagePerceptualQuality
	if v1 != "" {
		v2 = imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQuality(v1.(string)))
	}

	return &imaging.OutputImagePerceptualQualityVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func regionOfInterestCropStyleVariableInline(src map[string]interface{}, name string) *imaging.RegionOfInterestCropStyleVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.RegionOfInterestCropStyle
	if v1 != "" {
		v2 = imaging.RegionOfInterestCropStylePtr(imaging.RegionOfInterestCropStyle(v1.(string)))
	}

	return &imaging.RegionOfInterestCropStyleVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func resizeAspectVariableInline(src map[string]interface{}, name string) *imaging.ResizeAspectVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.ResizeAspect
	if v1 != "" {
		v2 = imaging.ResizeAspectPtr(imaging.ResizeAspect(v1.(string)))
	}

	return &imaging.ResizeAspectVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func resizeTypeVariableInline(src map[string]interface{}, name string) *imaging.ResizeTypeVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.ResizeType
	if v1 != "" {
		v2 = imaging.ResizeTypePtr(imaging.ResizeType(v1.(string)))
	}

	return &imaging.ResizeTypeVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func booleanVariableInline(src map[string]interface{}, name string) *imaging.BooleanVariableInline {
	if variableHasValue(src, name) {
		return &imaging.BooleanVariableInline{
			Name:  stringValuePtr(src, name+"_var"),
			Value: boolValuePtr(src, name),
		}
	}
	return nil
}

func numberVariableInline(src map[string]interface{}, name string) *imaging.NumberVariableInline {
	if variableHasValue(src, name) {
		return &imaging.NumberVariableInline{
			Name:  stringValuePtr(src, name+"_var"),
			Value: float64ValuePtr(src, name),
		}
	}
	return nil
}

func integerVariableInline(src map[string]interface{}, name string) *imaging.IntegerVariableInline {
	if variableHasValue(src, name) {
		return &imaging.IntegerVariableInline{
			Name:  stringValuePtr(src, name+"_var"),
			Value: intValuePtr(src, name),
		}
	}
	return nil
}

func stringVariableInline(src map[string]interface{}, name string) *imaging.StringVariableInline {
	if variableHasValue(src, name) {
		return &imaging.StringVariableInline{
			Name:  stringValuePtr(src, name+"_var"),
			Value: stringValuePtr(src, name),
		}
	}
	return nil
}

func boolValuePtr(src map[string]interface{}, name string) *bool {
	value := src[name]
	if value != "" {
		v := value.(bool)
		return &v
	}
	return nil
}

func float64ValuePtr(src map[string]interface{}, name string) *float64 {
	value := src[name]
	if value != "" {
		v := value.(float64)
		return &v
	}
	return nil
}

func intValuePtr(src map[string]interface{}, name string) *int {
	value := src[name]
	if value != "" {
		v := value.(int)
		return &v
	}
	return nil
}

func stringValuePtr(src map[string]interface{}, name string) *string {
	value := src[name]
	if value != "" {
		v := value.(string)
		return &v
	}
	return nil
}

func interfaceSliceToIntSlice(list []interface{}) []int {
	if len(list) == 0 {
		return nil
	}
	intList := make([]int, len(list))
	for i, v := range list {
		intList[i] = v.(int)
	}
	return intList
}

func interfaceSliceToStringSlice(list []interface{}) []string {
	if len(list) == 0 {
		return nil
	}
	stringList := make([]string, len(list))
	for i, v := range list {
		stringList[i] = v.(string)
	}
	return stringList
}

func variableHasValue(src map[string]interface{}, name string) bool {
	v1 := src[name]
	v2 := src[name+"_var"]

	if !reflect.ValueOf(v1).IsZero() || !reflect.ValueOf(v2).IsZero() {
		return true
	}
	return false
}

func extract(src map[string]interface{}, name string) map[string]interface{} {
	elem, ok := src[name]
	if !ok {
		return nil
	}

	l := elem.([]interface{})
	if len(l) == 0 {
		return nil
	}
	return l[0].(map[string]interface{})
}

// Package imagewriter contains set of functions used to manage image
package imagewriter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PolicyImageToEdgeGrid converts terraform originated map structure into EdgeGrid structure
func PolicyImageToEdgeGrid(d *schema.ResourceData, key string) imaging.PolicyInputImage {
	_, exist := extract(d, key)
	var result imaging.PolicyInputImage
	if exist {
		result = imaging.PolicyInputImage{}
		result.Breakpoints = getBreakpoints(d, getKey(key, 0, "breakpoints"))
		result.Hosts = interfaceSliceToStringSlice(d, getKey(key, 0, "hosts"))
		result.Output = getOutputImage(d, getKey(key, 0, "output"))
		result.PostBreakpointTransformations = getPostBreakpointTransformations(d, getKey(key, 0, "post_breakpoint_transformations"))
		result.RolloutDuration = intReaderPtr(d, getKey(key, 0, "rollout_duration"))
		result.ServeStaleDuration = intReaderPtr(d, getKey(key, 0, "serve_stale_duration"))
		result.Transformations = getTransformations(d, getKey(key, 0, "transformations"))
		result.Variables = getVariableList(d, getKey(key, 0, "variables"))
	}

	return result
}

func getAppend(d *schema.ResourceData, key string) *imaging.Append {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Append{
			Gravity:                gravityVariableInline(d, getKey(key, 0, "gravity")),
			GravityPriority:        appendGravityPriorityVariableInline(d, getKey(key, 0, "gravity_priority")),
			Image:                  getImageType(d, getKey(key, 0, "image")),
			PreserveMinorDimension: booleanVariableInline(d, getKey(key, 0, "preserve_minor_dimension")),
			Transformation:         imaging.AppendTransformationAppend,
		}
		return &result
	}
	return nil
}

func getAspectCrop(d *schema.ResourceData, key string) *imaging.AspectCrop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.AspectCrop{
			AllowExpansion: booleanVariableInline(d, getKey(key, 0, "allow_expansion")),
			Height:         numberVariableInline(d, getKey(key, 0, "height")),
			Width:          numberVariableInline(d, getKey(key, 0, "width")),
			XPosition:      numberVariableInline(d, getKey(key, 0, "x_position")),
			YPosition:      numberVariableInline(d, getKey(key, 0, "y_position")),
			Transformation: imaging.AspectCropTransformationAspectCrop,
		}
		return &result
	}
	return nil
}

func getBackgroundColor(d *schema.ResourceData, key string) *imaging.BackgroundColor {
	_, exist := extract(d, key)
	if exist {
		result := imaging.BackgroundColor{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Transformation: imaging.BackgroundColorTransformationBackgroundColor,
		}
		return &result
	}
	return nil
}

func getBlur(d *schema.ResourceData, key string) *imaging.Blur {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Blur{
			Sigma:          numberVariableInline(d, getKey(key, 0, "sigma")),
			Transformation: imaging.BlurTransformationBlur,
		}
		return &result
	}
	return nil
}

func getBoxImageType(d *schema.ResourceData, key string) *imaging.BoxImageType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.BoxImageType{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Height:         integerVariableInline(d, getKey(key, 0, "height")),
			Transformation: getTransformationType(d, getKey(key, 0, "transformation")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Type:           imaging.BoxImageTypeTypeBox,
		}
		return &result
	}
	return nil
}

func getBoxImageTypePost(d *schema.ResourceData, key string) *imaging.BoxImageTypePost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.BoxImageTypePost{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Height:         integerVariableInline(d, getKey(key, 0, "height")),
			Transformation: getTransformationTypePost(d, getKey(key, 0, "transformation")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Type:           imaging.BoxImageTypePostTypeBox,
		}
		return &result
	}
	return nil
}

func getBreakpoints(d *schema.ResourceData, key string) *imaging.Breakpoints {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Breakpoints{
			Widths: interfaceSliceToIntSlice(d, getKey(key, 0, "widths")),
		}
		return &result
	}
	return nil
}

func getChromaKey(d *schema.ResourceData, key string) *imaging.ChromaKey {
	_, exist := extract(d, key)
	if exist {
		result := imaging.ChromaKey{
			Hue:                 numberVariableInline(d, getKey(key, 0, "hue")),
			HueFeather:          numberVariableInline(d, getKey(key, 0, "hue_feather")),
			HueTolerance:        numberVariableInline(d, getKey(key, 0, "hue_tolerance")),
			LightnessFeather:    numberVariableInline(d, getKey(key, 0, "lightness_feather")),
			LightnessTolerance:  numberVariableInline(d, getKey(key, 0, "lightness_tolerance")),
			SaturationFeather:   numberVariableInline(d, getKey(key, 0, "saturation_feather")),
			SaturationTolerance: numberVariableInline(d, getKey(key, 0, "saturation_tolerance")),
			Transformation:      imaging.ChromaKeyTransformationChromaKey,
		}
		return &result
	}
	return nil
}

func getCircleImageType(d *schema.ResourceData, key string) *imaging.CircleImageType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.CircleImageType{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Diameter:       integerVariableInline(d, getKey(key, 0, "diameter")),
			Transformation: getTransformationType(d, getKey(key, 0, "transformation")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Type:           imaging.CircleImageTypeTypeCircle,
		}
		return &result
	}
	return nil
}

func getCircleImageTypePost(d *schema.ResourceData, key string) *imaging.CircleImageTypePost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.CircleImageTypePost{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Diameter:       integerVariableInline(d, getKey(key, 0, "diameter")),
			Transformation: getTransformationTypePost(d, getKey(key, 0, "transformation")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Type:           imaging.CircleImageTypePostTypeCircle,
		}
		return &result
	}
	return nil
}

func getCircleShapeType(d *schema.ResourceData, key string) *imaging.CircleShapeType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.CircleShapeType{
			Center: getPointShapeType(d, getKey(key, 0, "center")),
			Radius: numberVariableInline(d, getKey(key, 0, "radius")),
		}
		return &result
	}
	return nil
}

func getComposite(d *schema.ResourceData, key string) *imaging.Composite {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Composite{
			Gravity:        gravityVariableInline(d, getKey(key, 0, "gravity")),
			Image:          getImageType(d, getKey(key, 0, "image")),
			Placement:      compositePlacementVariableInline(d, getKey(key, 0, "placement")),
			Scale:          numberVariableInline(d, getKey(key, 0, "scale")),
			ScaleDimension: compositeScaleDimensionVariableInline(d, getKey(key, 0, "scale_dimension")),
			XPosition:      integerVariableInline(d, getKey(key, 0, "x_position")),
			YPosition:      integerVariableInline(d, getKey(key, 0, "y_position")),
			Transformation: imaging.CompositeTransformationComposite,
		}
		return &result
	}
	return nil
}

func getCompositePost(d *schema.ResourceData, key string) *imaging.CompositePost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.CompositePost{
			Gravity:        gravityPostVariableInline(d, getKey(key, 0, "gravity")),
			Image:          getImageTypePost(d, getKey(key, 0, "image")),
			Placement:      compositePostPlacementVariableInline(d, getKey(key, 0, "placement")),
			Scale:          numberVariableInline(d, getKey(key, 0, "scale")),
			ScaleDimension: compositePostScaleDimensionVariableInline(d, getKey(key, 0, "scale_dimension")),
			XPosition:      integerVariableInline(d, getKey(key, 0, "x_position")),
			YPosition:      integerVariableInline(d, getKey(key, 0, "y_position")),
			Transformation: imaging.CompositePostTransformationComposite,
		}
		return &result
	}
	return nil
}

func getCompound(d *schema.ResourceData, key string) *imaging.Compound {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Compound{
			Transformations: getTransformations(d, getKey(key, 0, "transformations")),
			Transformation:  imaging.CompoundTransformationCompound,
		}
		return &result
	}
	return nil
}

func getCompoundPost(d *schema.ResourceData, key string) *imaging.CompoundPost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.CompoundPost{
			Transformations: getPostBreakpointTransformations(d, getKey(key, 0, "transformations")),
			Transformation:  imaging.CompoundPostTransformationCompound,
		}
		return &result
	}
	return nil
}

func getContrast(d *schema.ResourceData, key string) *imaging.Contrast {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Contrast{
			Brightness:     numberVariableInline(d, getKey(key, 0, "brightness")),
			Contrast:       numberVariableInline(d, getKey(key, 0, "contrast")),
			Transformation: imaging.ContrastTransformationContrast,
		}
		return &result
	}
	return nil
}

func getCrop(d *schema.ResourceData, key string) *imaging.Crop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Crop{
			AllowExpansion: booleanVariableInline(d, getKey(key, 0, "allow_expansion")),
			Gravity:        gravityVariableInline(d, getKey(key, 0, "gravity")),
			Height:         integerVariableInline(d, getKey(key, 0, "height")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			XPosition:      integerVariableInline(d, getKey(key, 0, "x_position")),
			YPosition:      integerVariableInline(d, getKey(key, 0, "y_position")),
			Transformation: imaging.CropTransformationCrop,
		}
		return &result
	}
	return nil
}

func getEnumOptionsList(d *schema.ResourceData, key string) []*imaging.EnumOptions {
	collection, exist := extract(d, key)
	if exist {
		result := make([]*imaging.EnumOptions, 0)
		for idx := range collection.([]interface{}) {
			elem := imaging.EnumOptions{
				ID:    stringReader(d, getKey(key, idx, "id")),
				Value: stringReader(d, getKey(key, idx, "value")),
			}
			result = append(result, &elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getFaceCrop(d *schema.ResourceData, key string) *imaging.FaceCrop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.FaceCrop{
			Algorithm:      faceCropAlgorithmVariableInline(d, getKey(key, 0, "algorithm")),
			Confidence:     numberVariableInline(d, getKey(key, 0, "confidence")),
			FailGravity:    gravityVariableInline(d, getKey(key, 0, "fail_gravity")),
			Focus:          faceCropFocusVariableInline(d, getKey(key, 0, "focus")),
			Gravity:        gravityVariableInline(d, getKey(key, 0, "gravity")),
			Height:         integerVariableInline(d, getKey(key, 0, "height")),
			Padding:        numberVariableInline(d, getKey(key, 0, "padding")),
			Style:          faceCropStyleVariableInline(d, getKey(key, 0, "style")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Transformation: imaging.FaceCropTransformationFaceCrop,
		}
		return &result
	}
	return nil
}

func getFeatureCrop(d *schema.ResourceData, key string) *imaging.FeatureCrop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.FeatureCrop{
			FailGravity:       gravityVariableInline(d, getKey(key, 0, "fail_gravity")),
			FeatureRadius:     numberVariableInline(d, getKey(key, 0, "feature_radius")),
			Gravity:           gravityVariableInline(d, getKey(key, 0, "gravity")),
			Height:            integerVariableInline(d, getKey(key, 0, "height")),
			MaxFeatures:       integerVariableInline(d, getKey(key, 0, "max_features")),
			MinFeatureQuality: numberVariableInline(d, getKey(key, 0, "min_feature_quality")),
			Padding:           numberVariableInline(d, getKey(key, 0, "padding")),
			Style:             featureCropStyleVariableInline(d, getKey(key, 0, "style")),
			Width:             integerVariableInline(d, getKey(key, 0, "width")),
			Transformation:    imaging.FeatureCropTransformationFeatureCrop,
		}
		return &result
	}
	return nil
}

func getFitAndFill(d *schema.ResourceData, key string) *imaging.FitAndFill {
	_, exist := extract(d, key)
	if exist {
		result := imaging.FitAndFill{
			FillTransformation: getTransformationType(d, getKey(key, 0, "fill_transformation")),
			Height:             integerVariableInline(d, getKey(key, 0, "height")),
			Width:              integerVariableInline(d, getKey(key, 0, "width")),
			Transformation:     imaging.FitAndFillTransformationFitAndFill,
		}
		return &result
	}
	return nil
}

func getGoop(d *schema.ResourceData, key string) *imaging.Goop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Goop{
			Chaos:          numberVariableInline(d, getKey(key, 0, "chaos")),
			Density:        integerVariableInline(d, getKey(key, 0, "density")),
			Power:          numberVariableInline(d, getKey(key, 0, "power")),
			Seed:           integerVariableInline(d, getKey(key, 0, "seed")),
			Transformation: imaging.GoopTransformationGoop,
		}
		return &result
	}
	return nil
}

func getGrayscale(d *schema.ResourceData, key string) *imaging.Grayscale {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Grayscale{
			Type:           grayscaleTypeVariableInline(d, getKey(key, 0, "type")),
			Transformation: imaging.GrayscaleTransformationGrayscale,
		}
		return &result
	}
	return nil
}

func getHSL(d *schema.ResourceData, key string) *imaging.HSL {
	_, exist := extract(d, key)
	if exist {
		result := imaging.HSL{
			Hue:            numberVariableInline(d, getKey(key, 0, "hue")),
			Lightness:      numberVariableInline(d, getKey(key, 0, "lightness")),
			Saturation:     numberVariableInline(d, getKey(key, 0, "saturation")),
			Transformation: imaging.HSLTransformationHSL,
		}
		return &result
	}
	return nil
}

func getHSV(d *schema.ResourceData, key string) *imaging.HSV {
	_, exist := extract(d, key)
	if exist {
		result := imaging.HSV{
			Hue:            numberVariableInline(d, getKey(key, 0, "hue")),
			Saturation:     numberVariableInline(d, getKey(key, 0, "saturation")),
			Value:          numberVariableInline(d, getKey(key, 0, "value")),
			Transformation: imaging.HSVTransformationHSV,
		}
		return &result
	}
	return nil
}

func getIfDimension(d *schema.ResourceData, key string) *imaging.IfDimension {
	_, exist := extract(d, key)
	if exist {
		result := imaging.IfDimension{
			Default:        getTransformationType(d, getKey(key, 0, "default")),
			Dimension:      ifDimensionDimensionVariableInline(d, getKey(key, 0, "dimension")),
			Equal:          getTransformationType(d, getKey(key, 0, "equal")),
			GreaterThan:    getTransformationType(d, getKey(key, 0, "greater_than")),
			LessThan:       getTransformationType(d, getKey(key, 0, "less_than")),
			Value:          integerVariableInline(d, getKey(key, 0, "value")),
			Transformation: imaging.IfDimensionTransformationIfDimension,
		}
		return &result
	}
	return nil
}

func getIfDimensionPost(d *schema.ResourceData, key string) *imaging.IfDimensionPost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.IfDimensionPost{
			Default:        getTransformationTypePost(d, getKey(key, 0, "default")),
			Dimension:      ifDimensionPostDimensionVariableInline(d, getKey(key, 0, "dimension")),
			Equal:          getTransformationTypePost(d, getKey(key, 0, "equal")),
			GreaterThan:    getTransformationTypePost(d, getKey(key, 0, "greater_than")),
			LessThan:       getTransformationTypePost(d, getKey(key, 0, "less_than")),
			Value:          integerVariableInline(d, getKey(key, 0, "value")),
			Transformation: imaging.IfDimensionPostTransformationIfDimension,
		}
		return &result
	}
	return nil
}

func getIfOrientation(d *schema.ResourceData, key string) *imaging.IfOrientation {
	_, exist := extract(d, key)
	if exist {
		result := imaging.IfOrientation{
			Default:        getTransformationType(d, getKey(key, 0, "default")),
			Landscape:      getTransformationType(d, getKey(key, 0, "landscape")),
			Portrait:       getTransformationType(d, getKey(key, 0, "portrait")),
			Square:         getTransformationType(d, getKey(key, 0, "square")),
			Transformation: imaging.IfOrientationTransformationIfOrientation,
		}
		return &result
	}
	return nil
}

func getIfOrientationPost(d *schema.ResourceData, key string) *imaging.IfOrientationPost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.IfOrientationPost{
			Default:        getTransformationTypePost(d, getKey(key, 0, "default")),
			Landscape:      getTransformationTypePost(d, getKey(key, 0, "landscape")),
			Portrait:       getTransformationTypePost(d, getKey(key, 0, "portrait")),
			Square:         getTransformationTypePost(d, getKey(key, 0, "square")),
			Transformation: imaging.IfOrientationPostTransformationIfOrientation,
		}
		return &result
	}
	return nil
}

func getImQuery(d *schema.ResourceData, key string) *imaging.ImQuery {
	_, exist := extract(d, key)
	if exist {
		result := imaging.ImQuery{
			AllowedTransformations: interfaceSliceToImagingImQueryAllowedTransformationsSlice(d, getKey(key, 0, "allowed_transformations")),
			Query:                  queryVariableInline(d, getKey(key, 0, "query")),
			Transformation:         imaging.ImQueryTransformationImQuery,
		}
		return &result
	}
	return nil
}

func getMaxColors(d *schema.ResourceData, key string) *imaging.MaxColors {
	_, exist := extract(d, key)
	if exist {
		result := imaging.MaxColors{
			Colors:         integerVariableInline(d, getKey(key, 0, "colors")),
			Transformation: imaging.MaxColorsTransformationMaxColors,
		}
		return &result
	}
	return nil
}

func getMirror(d *schema.ResourceData, key string) *imaging.Mirror {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Mirror{
			Horizontal:     booleanVariableInline(d, getKey(key, 0, "horizontal")),
			Vertical:       booleanVariableInline(d, getKey(key, 0, "vertical")),
			Transformation: imaging.MirrorTransformationMirror,
		}
		return &result
	}
	return nil
}

func getMonoHue(d *schema.ResourceData, key string) *imaging.MonoHue {
	_, exist := extract(d, key)
	if exist {
		result := imaging.MonoHue{
			Hue:            numberVariableInline(d, getKey(key, 0, "hue")),
			Transformation: imaging.MonoHueTransformationMonoHue,
		}
		return &result
	}
	return nil
}

func getOpacity(d *schema.ResourceData, key string) *imaging.Opacity {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Opacity{
			Opacity:        numberVariableInline(d, getKey(key, 0, "opacity")),
			Transformation: imaging.OpacityTransformationOpacity,
		}
		return &result
	}
	return nil
}

func getOutputImage(d *schema.ResourceData, key string) *imaging.OutputImage {
	_, exist := extract(d, key)
	if exist {
		result := imaging.OutputImage{
			AdaptiveQuality:         intReaderPtr(d, getKey(key, 0, "adaptive_quality")),
			AllowPristineOnDownsize: boolReaderPtr(d, getKey(key, 0, "allow_pristine_on_downsize")),
			AllowedFormats:          interfaceSliceToImagingOutputImageAllowedFormatsSlice(d, getKey(key, 0, "allowed_formats")),
			ForcedFormats:           interfaceSliceToImagingOutputImageForcedFormatsSlice(d, getKey(key, 0, "forced_formats")),
			PerceptualQuality:       outputImagePerceptualQualityVariableInline(d, getKey(key, 0, "perceptual_quality")),
			PerceptualQualityFloor:  intReaderPtr(d, getKey(key, 0, "perceptual_quality_floor")),
			PreferModernFormats:     boolReaderPtr(d, getKey(key, 0, "prefer_modern_formats")),
			Quality:                 integerVariableInline(d, getKey(key, 0, "quality")),
		}
		return &result
	}
	return nil
}

func getPointShapeType(d *schema.ResourceData, key string) *imaging.PointShapeType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.PointShapeType{
			X: numberVariableInline(d, getKey(key, 0, "x")),
			Y: numberVariableInline(d, getKey(key, 0, "y")),
		}
		return &result
	}
	return nil
}

func getPointShapeTypeList(d *schema.ResourceData, key string) []imaging.PointShapeType {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.PointShapeType, 0)
		for idx := range collection.([]interface{}) {
			elem := imaging.PointShapeType{
				X: numberVariableInline(d, getKey(key, idx, "x")),
				Y: numberVariableInline(d, getKey(key, idx, "y")),
			}
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getPolygonShapeType(d *schema.ResourceData, key string) *imaging.PolygonShapeType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.PolygonShapeType{
			Points: getPointShapeTypeList(d, getKey(key, 0, "points")),
		}
		return &result
	}
	return nil
}

func getPostBreakpointTransformations(d *schema.ResourceData, key string) []imaging.TransformationTypePost {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.TransformationTypePost, 0)
		for idx := range collection.([]interface{}) {
			elem := getTransformationTypePost(d, fmt.Sprintf("%s.%d", key, idx))
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getRectangleShapeType(d *schema.ResourceData, key string) *imaging.RectangleShapeType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.RectangleShapeType{
			Anchor: getPointShapeType(d, getKey(key, 0, "anchor")),
			Height: numberVariableInline(d, getKey(key, 0, "height")),
			Width:  numberVariableInline(d, getKey(key, 0, "width")),
		}
		return &result
	}
	return nil
}

func getRegionOfInterestCrop(d *schema.ResourceData, key string) *imaging.RegionOfInterestCrop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.RegionOfInterestCrop{
			Gravity:          gravityVariableInline(d, getKey(key, 0, "gravity")),
			Height:           integerVariableInline(d, getKey(key, 0, "height")),
			RegionOfInterest: getShapeType(d, getKey(key, 0, "region_of_interest")),
			Style:            regionOfInterestCropStyleVariableInline(d, getKey(key, 0, "style")),
			Width:            integerVariableInline(d, getKey(key, 0, "width")),
			Transformation:   imaging.RegionOfInterestCropTransformationRegionOfInterestCrop,
		}
		return &result
	}
	return nil
}

func getRelativeCrop(d *schema.ResourceData, key string) *imaging.RelativeCrop {
	_, exist := extract(d, key)
	if exist {
		result := imaging.RelativeCrop{
			East:           integerVariableInline(d, getKey(key, 0, "east")),
			North:          integerVariableInline(d, getKey(key, 0, "north")),
			South:          integerVariableInline(d, getKey(key, 0, "south")),
			West:           integerVariableInline(d, getKey(key, 0, "west")),
			Transformation: imaging.RelativeCropTransformationRelativeCrop,
		}
		return &result
	}
	return nil
}

func getRemoveColor(d *schema.ResourceData, key string) *imaging.RemoveColor {
	_, exist := extract(d, key)
	if exist {
		result := imaging.RemoveColor{
			Color:          stringVariableInline(d, getKey(key, 0, "color")),
			Feather:        numberVariableInline(d, getKey(key, 0, "feather")),
			Tolerance:      numberVariableInline(d, getKey(key, 0, "tolerance")),
			Transformation: imaging.RemoveColorTransformationRemoveColor,
		}
		return &result
	}
	return nil
}

func getResize(d *schema.ResourceData, key string) *imaging.Resize {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Resize{
			Aspect:         resizeAspectVariableInline(d, getKey(key, 0, "aspect")),
			Height:         integerVariableInline(d, getKey(key, 0, "height")),
			Type:           resizeTypeVariableInline(d, getKey(key, 0, "type")),
			Width:          integerVariableInline(d, getKey(key, 0, "width")),
			Transformation: imaging.ResizeTransformationResize,
		}
		return &result
	}
	return nil
}

func getRotate(d *schema.ResourceData, key string) *imaging.Rotate {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Rotate{
			Degrees:        numberVariableInline(d, getKey(key, 0, "degrees")),
			Transformation: imaging.RotateTransformationRotate,
		}
		return &result
	}
	return nil
}

func getScale(d *schema.ResourceData, key string) *imaging.Scale {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Scale{
			Height:         numberVariableInline(d, getKey(key, 0, "height")),
			Width:          numberVariableInline(d, getKey(key, 0, "width")),
			Transformation: imaging.ScaleTransformationScale,
		}
		return &result
	}
	return nil
}

func getShapeTypeList(d *schema.ResourceData, key string) []imaging.ShapeType {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.ShapeType, 0)
		for idx := range collection.([]interface{}) {
			elem := getShapeType(d, fmt.Sprintf("%s.%d", key, idx))
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getShear(d *schema.ResourceData, key string) *imaging.Shear {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Shear{
			XShear:         numberVariableInline(d, getKey(key, 0, "x_shear")),
			YShear:         numberVariableInline(d, getKey(key, 0, "y_shear")),
			Transformation: imaging.ShearTransformationShear,
		}
		return &result
	}
	return nil
}

func getTextImageType(d *schema.ResourceData, key string) *imaging.TextImageType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.TextImageType{
			Fill:           stringVariableInline(d, getKey(key, 0, "fill")),
			Size:           numberVariableInline(d, getKey(key, 0, "size")),
			Stroke:         stringVariableInline(d, getKey(key, 0, "stroke")),
			StrokeSize:     numberVariableInline(d, getKey(key, 0, "stroke_size")),
			Text:           stringVariableInline(d, getKey(key, 0, "text")),
			Transformation: getTransformationType(d, getKey(key, 0, "transformation")),
			Typeface:       stringVariableInline(d, getKey(key, 0, "typeface")),
			Type:           imaging.TextImageTypeTypeText,
		}
		return &result
	}
	return nil
}

func getTextImageTypePost(d *schema.ResourceData, key string) *imaging.TextImageTypePost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.TextImageTypePost{
			Fill:           stringVariableInline(d, getKey(key, 0, "fill")),
			Size:           numberVariableInline(d, getKey(key, 0, "size")),
			Stroke:         stringVariableInline(d, getKey(key, 0, "stroke")),
			StrokeSize:     numberVariableInline(d, getKey(key, 0, "stroke_size")),
			Text:           stringVariableInline(d, getKey(key, 0, "text")),
			Transformation: getTransformationTypePost(d, getKey(key, 0, "transformation")),
			Typeface:       stringVariableInline(d, getKey(key, 0, "typeface")),
			Type:           imaging.TextImageTypePostTypeText,
		}
		return &result
	}
	return nil
}

func getTransformations(d *schema.ResourceData, key string) []imaging.TransformationType {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.TransformationType, 0)
		for idx := range collection.([]interface{}) {
			elem := getTransformationType(d, fmt.Sprintf("%s.%d", key, idx))
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getTrim(d *schema.ResourceData, key string) *imaging.Trim {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Trim{
			Fuzz:           numberVariableInline(d, getKey(key, 0, "fuzz")),
			Padding:        integerVariableInline(d, getKey(key, 0, "padding")),
			Transformation: imaging.TrimTransformationTrim,
		}
		return &result
	}
	return nil
}

func getURLImageType(d *schema.ResourceData, key string) *imaging.URLImageType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.URLImageType{
			Transformation: getTransformationType(d, getKey(key, 0, "transformation")),
			URL:            stringVariableInline(d, getKey(key, 0, "url")),
			Type:           imaging.URLImageTypeTypeURL,
		}
		return &result
	}
	return nil
}

func getURLImageTypePost(d *schema.ResourceData, key string) *imaging.URLImageTypePost {
	_, exist := extract(d, key)
	if exist {
		result := imaging.URLImageTypePost{
			Transformation: getTransformationTypePost(d, getKey(key, 0, "transformation")),
			URL:            stringVariableInline(d, getKey(key, 0, "url")),
			Type:           imaging.URLImageTypePostTypeURL,
		}
		return &result
	}
	return nil
}

func getUnionShapeType(d *schema.ResourceData, key string) *imaging.UnionShapeType {
	_, exist := extract(d, key)
	if exist {
		result := imaging.UnionShapeType{
			Shapes: getShapeTypeList(d, getKey(key, 0, "shapes")),
		}
		return &result
	}
	return nil
}

func getUnsharpMask(d *schema.ResourceData, key string) *imaging.UnsharpMask {
	_, exist := extract(d, key)
	if exist {
		result := imaging.UnsharpMask{
			Gain:           numberVariableInline(d, getKey(key, 0, "gain")),
			Sigma:          numberVariableInline(d, getKey(key, 0, "sigma")),
			Threshold:      numberVariableInline(d, getKey(key, 0, "threshold")),
			Transformation: imaging.UnsharpMaskTransformationUnsharpMask,
		}
		return &result
	}
	return nil
}

func getVariableList(d *schema.ResourceData, key string) []imaging.Variable {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.Variable, 0)
		for idx := range collection.([]interface{}) {
			elem := imaging.Variable{
				DefaultValue: stringReader(d, getKey(key, idx, "default_value")),
				EnumOptions:  getEnumOptionsList(d, getKey(key, idx, "enum_options")),
				Name:         stringReader(d, getKey(key, idx, "name")),
				Postfix:      stringReaderPtr(d, getKey(key, idx, "postfix")),
				Prefix:       stringReaderPtr(d, getKey(key, idx, "prefix")),
				Type:         imaging.VariableType(stringReader(d, getKey(key, idx, "type"))),
			}
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getImageType(d *schema.ResourceData, key string) imaging.ImageType {
	_, isAny := extract(d, key)

	if !isAny {
		return nil
	}

	key = decorateKeyForSlice(key)

	var exist bool
	_, exist = extract(d, key+".box_image")
	if exist {
		return getBoxImageType(d, key+".box_image")
	}
	_, exist = extract(d, key+".circle_image")
	if exist {
		return getCircleImageType(d, key+".circle_image")
	}
	_, exist = extract(d, key+".text_image")
	if exist {
		return getTextImageType(d, key+".text_image")
	}
	_, exist = extract(d, key+".url_image")
	if exist {
		return getURLImageType(d, key+".url_image")
	}
	panic(fmt.Sprint("unsupported type"))
}

func getImageTypePost(d *schema.ResourceData, key string) imaging.ImageTypePost {
	_, isAny := extract(d, key)

	if !isAny {
		return nil
	}

	key = decorateKeyForSlice(key)

	var exist bool
	_, exist = extract(d, key+".box_image")
	if exist {
		return getBoxImageTypePost(d, key+".box_image")
	}
	_, exist = extract(d, key+".circle_image")
	if exist {
		return getCircleImageTypePost(d, key+".circle_image")
	}
	_, exist = extract(d, key+".text_image")
	if exist {
		return getTextImageTypePost(d, key+".text_image")
	}
	_, exist = extract(d, key+".url_image")
	if exist {
		return getURLImageTypePost(d, key+".url_image")
	}
	panic(fmt.Sprint("unsupported type"))
}

func getShapeType(d *schema.ResourceData, key string) imaging.ShapeType {
	_, isAny := extract(d, key)

	if !isAny {
		return nil
	}

	key = decorateKeyForSlice(key)

	var exist bool
	_, exist = extract(d, key+".circle_shape")
	if exist {
		return getCircleShapeType(d, key+".circle_shape")
	}
	_, exist = extract(d, key+".point_shape")
	if exist {
		return getPointShapeType(d, key+".point_shape")
	}
	_, exist = extract(d, key+".polygon_shape")
	if exist {
		return getPolygonShapeType(d, key+".polygon_shape")
	}
	_, exist = extract(d, key+".rectangle_shape")
	if exist {
		return getRectangleShapeType(d, key+".rectangle_shape")
	}
	_, exist = extract(d, key+".union_shape")
	if exist {
		return getUnionShapeType(d, key+".union_shape")
	}
	panic(fmt.Sprint("unsupported type"))
}

func getTransformationType(d *schema.ResourceData, key string) imaging.TransformationType {
	_, isAny := extract(d, key)

	if !isAny {
		return nil
	}

	key = decorateKeyForSlice(key)

	var exist bool
	_, exist = extract(d, key+".append")
	if exist {
		return getAppend(d, key+".append")
	}
	_, exist = extract(d, key+".aspect_crop")
	if exist {
		return getAspectCrop(d, key+".aspect_crop")
	}
	_, exist = extract(d, key+".background_color")
	if exist {
		return getBackgroundColor(d, key+".background_color")
	}
	_, exist = extract(d, key+".blur")
	if exist {
		return getBlur(d, key+".blur")
	}
	_, exist = extract(d, key+".chroma_key")
	if exist {
		return getChromaKey(d, key+".chroma_key")
	}
	_, exist = extract(d, key+".composite")
	if exist {
		return getComposite(d, key+".composite")
	}
	_, exist = extract(d, key+".compound")
	if exist {
		return getCompound(d, key+".compound")
	}
	_, exist = extract(d, key+".contrast")
	if exist {
		return getContrast(d, key+".contrast")
	}
	_, exist = extract(d, key+".crop")
	if exist {
		return getCrop(d, key+".crop")
	}
	_, exist = extract(d, key+".face_crop")
	if exist {
		return getFaceCrop(d, key+".face_crop")
	}
	_, exist = extract(d, key+".feature_crop")
	if exist {
		return getFeatureCrop(d, key+".feature_crop")
	}
	_, exist = extract(d, key+".fit_and_fill")
	if exist {
		return getFitAndFill(d, key+".fit_and_fill")
	}
	_, exist = extract(d, key+".goop")
	if exist {
		return getGoop(d, key+".goop")
	}
	_, exist = extract(d, key+".grayscale")
	if exist {
		return getGrayscale(d, key+".grayscale")
	}
	_, exist = extract(d, key+".hsl")
	if exist {
		return getHSL(d, key+".hsl")
	}
	_, exist = extract(d, key+".hsv")
	if exist {
		return getHSV(d, key+".hsv")
	}
	_, exist = extract(d, key+".if_dimension")
	if exist {
		return getIfDimension(d, key+".if_dimension")
	}
	_, exist = extract(d, key+".if_orientation")
	if exist {
		return getIfOrientation(d, key+".if_orientation")
	}
	_, exist = extract(d, key+".im_query")
	if exist {
		return getImQuery(d, key+".im_query")
	}
	_, exist = extract(d, key+".max_colors")
	if exist {
		return getMaxColors(d, key+".max_colors")
	}
	_, exist = extract(d, key+".mirror")
	if exist {
		return getMirror(d, key+".mirror")
	}
	_, exist = extract(d, key+".mono_hue")
	if exist {
		return getMonoHue(d, key+".mono_hue")
	}
	_, exist = extract(d, key+".opacity")
	if exist {
		return getOpacity(d, key+".opacity")
	}
	_, exist = extract(d, key+".region_of_interest_crop")
	if exist {
		return getRegionOfInterestCrop(d, key+".region_of_interest_crop")
	}
	_, exist = extract(d, key+".relative_crop")
	if exist {
		return getRelativeCrop(d, key+".relative_crop")
	}
	_, exist = extract(d, key+".remove_color")
	if exist {
		return getRemoveColor(d, key+".remove_color")
	}
	_, exist = extract(d, key+".resize")
	if exist {
		return getResize(d, key+".resize")
	}
	_, exist = extract(d, key+".rotate")
	if exist {
		return getRotate(d, key+".rotate")
	}
	_, exist = extract(d, key+".scale")
	if exist {
		return getScale(d, key+".scale")
	}
	_, exist = extract(d, key+".shear")
	if exist {
		return getShear(d, key+".shear")
	}
	_, exist = extract(d, key+".trim")
	if exist {
		return getTrim(d, key+".trim")
	}
	_, exist = extract(d, key+".unsharp_mask")
	if exist {
		return getUnsharpMask(d, key+".unsharp_mask")
	}
	panic(fmt.Sprint("unsupported type"))
}

func getTransformationTypePost(d *schema.ResourceData, key string) imaging.TransformationTypePost {
	_, isAny := extract(d, key)

	if !isAny {
		return nil
	}

	key = decorateKeyForSlice(key)

	var exist bool
	_, exist = extract(d, key+".background_color")
	if exist {
		return getBackgroundColor(d, key+".background_color")
	}
	_, exist = extract(d, key+".blur")
	if exist {
		return getBlur(d, key+".blur")
	}
	_, exist = extract(d, key+".chroma_key")
	if exist {
		return getChromaKey(d, key+".chroma_key")
	}
	_, exist = extract(d, key+".composite")
	if exist {
		return getCompositePost(d, key+".composite")
	}
	_, exist = extract(d, key+".compound")
	if exist {
		return getCompoundPost(d, key+".compound")
	}
	_, exist = extract(d, key+".contrast")
	if exist {
		return getContrast(d, key+".contrast")
	}
	_, exist = extract(d, key+".goop")
	if exist {
		return getGoop(d, key+".goop")
	}
	_, exist = extract(d, key+".grayscale")
	if exist {
		return getGrayscale(d, key+".grayscale")
	}
	_, exist = extract(d, key+".hsl")
	if exist {
		return getHSL(d, key+".hsl")
	}
	_, exist = extract(d, key+".hsv")
	if exist {
		return getHSV(d, key+".hsv")
	}
	_, exist = extract(d, key+".if_dimension")
	if exist {
		return getIfDimensionPost(d, key+".if_dimension")
	}
	_, exist = extract(d, key+".if_orientation")
	if exist {
		return getIfOrientationPost(d, key+".if_orientation")
	}
	_, exist = extract(d, key+".max_colors")
	if exist {
		return getMaxColors(d, key+".max_colors")
	}
	_, exist = extract(d, key+".mirror")
	if exist {
		return getMirror(d, key+".mirror")
	}
	_, exist = extract(d, key+".mono_hue")
	if exist {
		return getMonoHue(d, key+".mono_hue")
	}
	_, exist = extract(d, key+".opacity")
	if exist {
		return getOpacity(d, key+".opacity")
	}
	_, exist = extract(d, key+".remove_color")
	if exist {
		return getRemoveColor(d, key+".remove_color")
	}
	_, exist = extract(d, key+".unsharp_mask")
	if exist {
		return getUnsharpMask(d, key+".unsharp_mask")
	}
	panic(fmt.Sprint("unsupported type"))
}

func appendGravityPriorityVariableInline(d *schema.ResourceData, key string) *imaging.AppendGravityPriorityVariableInline {
	var value *imaging.AppendGravityPriority
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.AppendGravityPriorityPtr(imaging.AppendGravityPriority(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.AppendGravityPriorityVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func compositePlacementVariableInline(d *schema.ResourceData, key string) *imaging.CompositePlacementVariableInline {
	var value *imaging.CompositePlacement
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.CompositePlacementPtr(imaging.CompositePlacement(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.CompositePlacementVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func compositePostPlacementVariableInline(d *schema.ResourceData, key string) *imaging.CompositePostPlacementVariableInline {
	var value *imaging.CompositePostPlacement
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.CompositePostPlacementPtr(imaging.CompositePostPlacement(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.CompositePostPlacementVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func compositePostScaleDimensionVariableInline(d *schema.ResourceData, key string) *imaging.CompositePostScaleDimensionVariableInline {
	var value *imaging.CompositePostScaleDimension
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.CompositePostScaleDimensionPtr(imaging.CompositePostScaleDimension(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.CompositePostScaleDimensionVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func compositeScaleDimensionVariableInline(d *schema.ResourceData, key string) *imaging.CompositeScaleDimensionVariableInline {
	var value *imaging.CompositeScaleDimension
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.CompositeScaleDimensionPtr(imaging.CompositeScaleDimension(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.CompositeScaleDimensionVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func faceCropAlgorithmVariableInline(d *schema.ResourceData, key string) *imaging.FaceCropAlgorithmVariableInline {
	var value *imaging.FaceCropAlgorithm
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.FaceCropAlgorithmPtr(imaging.FaceCropAlgorithm(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.FaceCropAlgorithmVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func faceCropFocusVariableInline(d *schema.ResourceData, key string) *imaging.FaceCropFocusVariableInline {
	var value *imaging.FaceCropFocus
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.FaceCropFocusPtr(imaging.FaceCropFocus(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.FaceCropFocusVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func faceCropStyleVariableInline(d *schema.ResourceData, key string) *imaging.FaceCropStyleVariableInline {
	var value *imaging.FaceCropStyle
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.FaceCropStylePtr(imaging.FaceCropStyle(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.FaceCropStyleVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func featureCropStyleVariableInline(d *schema.ResourceData, key string) *imaging.FeatureCropStyleVariableInline {
	var value *imaging.FeatureCropStyle
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.FeatureCropStylePtr(imaging.FeatureCropStyle(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.FeatureCropStyleVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func gravityPostVariableInline(d *schema.ResourceData, key string) *imaging.GravityPostVariableInline {
	var value *imaging.GravityPost
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.GravityPostPtr(imaging.GravityPost(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.GravityPostVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func gravityVariableInline(d *schema.ResourceData, key string) *imaging.GravityVariableInline {
	var value *imaging.Gravity
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.GravityPtr(imaging.Gravity(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.GravityVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func grayscaleTypeVariableInline(d *schema.ResourceData, key string) *imaging.GrayscaleTypeVariableInline {
	var value *imaging.GrayscaleType
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.GrayscaleTypePtr(imaging.GrayscaleType(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.GrayscaleTypeVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func ifDimensionDimensionVariableInline(d *schema.ResourceData, key string) *imaging.IfDimensionDimensionVariableInline {
	var value *imaging.IfDimensionDimension
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.IfDimensionDimensionPtr(imaging.IfDimensionDimension(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.IfDimensionDimensionVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func ifDimensionPostDimensionVariableInline(d *schema.ResourceData, key string) *imaging.IfDimensionPostDimensionVariableInline {
	var value *imaging.IfDimensionPostDimension
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.IfDimensionPostDimensionPtr(imaging.IfDimensionPostDimension(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.IfDimensionPostDimensionVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func outputImagePerceptualQualityVariableInline(d *schema.ResourceData, key string) *imaging.OutputImagePerceptualQualityVariableInline {
	var value *imaging.OutputImagePerceptualQuality
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQuality(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.OutputImagePerceptualQualityVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func regionOfInterestCropStyleVariableInline(d *schema.ResourceData, key string) *imaging.RegionOfInterestCropStyleVariableInline {
	var value *imaging.RegionOfInterestCropStyle
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.RegionOfInterestCropStylePtr(imaging.RegionOfInterestCropStyle(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.RegionOfInterestCropStyleVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func resizeAspectVariableInline(d *schema.ResourceData, key string) *imaging.ResizeAspectVariableInline {
	var value *imaging.ResizeAspect
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.ResizeAspectPtr(imaging.ResizeAspect(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.ResizeAspectVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func resizeTypeVariableInline(d *schema.ResourceData, key string) *imaging.ResizeTypeVariableInline {
	var value *imaging.ResizeType
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.ResizeTypePtr(imaging.ResizeType(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.ResizeTypeVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func booleanVariableInline(d *schema.ResourceData, key string) *imaging.BooleanVariableInline {
	var value *bool
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		valueMapped, _ := strconv.ParseBool(valueRaw.(string))
		value = tools.BoolPtr(valueMapped)
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.BooleanVariableInline{
			Name:  name,
			Value: value,
		}
	}
	return nil
}

func integerVariableInline(d *schema.ResourceData, key string) *imaging.IntegerVariableInline {
	var value *int
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		valueMapped, _ := strconv.Atoi(valueRaw.(string))
		value = tools.IntPtr(valueMapped)
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.IntegerVariableInline{
			Name:  name,
			Value: value,
		}
	}
	return nil
}

func numberVariableInline(d *schema.ResourceData, key string) *imaging.NumberVariableInline {
	var value *float64
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		valueMapped, _ := strconv.ParseFloat(valueRaw.(string), 64)
		value = tools.Float64Ptr(valueMapped)
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.NumberVariableInline{
			Name:  name,
			Value: value,
		}
	}
	return nil
}

func queryVariableInline(d *schema.ResourceData, key string) *imaging.QueryVariableInline {
	var name *string

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVar {
		return &imaging.QueryVariableInline{
			Name: name,
		}
	}
	return nil
}

func stringVariableInline(d *schema.ResourceData, key string) *imaging.StringVariableInline {
	var value *string
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		valueMapped := valueRaw.(string)
		value = tools.StringPtr(valueMapped)
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = tools.StringPtr(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.StringVariableInline{
			Name:  name,
			Value: value,
		}
	}
	return nil
}

func intReader(d *schema.ResourceData, key string) int {
	value, exist := extract(d, key)
	if exist {
		valInt, _ := strconv.Atoi(value.(string))
		return valInt
	}
	return 0
}

func intReaderPtr(d *schema.ResourceData, key string) *int {
	value, exist := extract(d, key)
	if exist {
		if valInt, err := strconv.Atoi(value.(string)); err == nil {
			return &valInt
		}
	}
	return nil
}

func stringReader(d *schema.ResourceData, key string) string {
	value, exist := extract(d, key)
	if exist {
		return value.(string)
	}
	return ""
}

func stringReaderPtr(d *schema.ResourceData, key string) *string {
	value, exist := extract(d, key)
	if exist {
		return tools.StringPtr(value.(string))
	}
	return nil
}

func boolReaderPtr(d *schema.ResourceData, key string) *bool {
	value, exist := extract(d, key)
	if exist {
		if value.(string) == "true" {
			return tools.BoolPtr(true)
		}
		return tools.BoolPtr(false)
	}
	return nil
}

func interfaceSliceToImagingImQueryAllowedTransformationsSlice(d *schema.ResourceData, key string) []imaging.ImQueryAllowedTransformations {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]imaging.ImQueryAllowedTransformations, len(l))
			for i, v := range l {
				result[i] = imaging.ImQueryAllowedTransformations(v.(string))
			}
			return result
		}
	}
	return nil
}

func interfaceSliceToImagingOutputImageAllowedFormatsSlice(d *schema.ResourceData, key string) []imaging.OutputImageAllowedFormats {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]imaging.OutputImageAllowedFormats, len(l))
			for i, v := range l {
				result[i] = imaging.OutputImageAllowedFormats(v.(string))
			}
			return result
		}
	}
	return nil
}

func interfaceSliceToImagingOutputImageForcedFormatsSlice(d *schema.ResourceData, key string) []imaging.OutputImageForcedFormats {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]imaging.OutputImageForcedFormats, len(l))
			for i, v := range l {
				result[i] = imaging.OutputImageForcedFormats(v.(string))
			}
			return result
		}
	}
	return nil
}

func interfaceSliceToIntSlice(d *schema.ResourceData, key string) []int {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]int, len(l))
			for i, v := range l {
				result[i] = v.(int)
			}
			return result
		}
	}
	return nil
}

func interfaceSliceToStringSlice(d *schema.ResourceData, key string) []string {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]string, len(l))
			for i, v := range l {
				result[i] = v.(string)
			}
			return result
		}
	}
	return nil
}

func extract(d *schema.ResourceData, key string) (interface{}, bool) {
	return d.GetOk(key)
}

func decorateKeyForSlice(key string) string {
	address := strings.Split(key, ".")
	matches, _ := regexp.MatchString("[0-9]+", address[len(address)-1])
	if !matches {
		return key + ".0"
	}
	return key
}

func getKey(rootKey string, elementIndex int, elementName string) string {
	return fmt.Sprintf("%s.%d.%s", rootKey, elementIndex, elementName)
}

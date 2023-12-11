package imaging

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func imageType(depth int) map[string]*schema.Schema {
	if depth <= 0 {
		return nil
	}
	return map[string]*schema.Schema{

		"box_image": {
			Description: "A rectangular box, with a specified color and applied transformation.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: boxImageType(depth - 1),
			},
		},

		"circle_image": {
			Description: "A rectangular box, with a specified color and applied transformation.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: circleImageType(depth - 1),
			},
		},

		"text_image": {
			Description: "A snippet of text. Defines font family and size, fill color, and outline stroke width and color.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: textImageType(depth - 1),
			},
		},

		"url_image": {
			Description: "An image loaded from a URL.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: urlImageType(depth - 1),
			},
		},
	}
}

func imageTypePost(depth int) map[string]*schema.Schema {
	if depth <= 0 {
		return nil
	}
	return map[string]*schema.Schema{

		"box_image": {
			Description: "A rectangular box, with a specified color and applied transformation.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: boxImageTypePost(depth - 1),
			},
		},

		"circle_image": {
			Description: "A rectangular box, with a specified color and applied transformation.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: circleImageTypePost(depth - 1),
			},
		},

		"text_image": {
			Description: "A snippet of text. Defines font family and size, fill color, and outline stroke width and color.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: textImageTypePost(depth - 1),
			},
		},

		"url_image": {
			Description: "An image loaded from a URL.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: urlImageTypePost(depth - 1),
			},
		},
	}
}

func shapeType(depth int) map[string]*schema.Schema {
	if depth <= 0 {
		return nil
	}
	return map[string]*schema.Schema{

		"circle_shape": {
			Description: "Defines a circle with a specified `radius` from its `center` point.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: circleShapeType(depth - 1),
			},
		},

		"point_shape": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: pointShapeType(depth - 1),
			},
		},

		"polygon_shape": {
			Description: "Defines a polygon from a series of connected points.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: polygonShapeType(depth - 1),
			},
		},

		"rectangle_shape": {
			Description: "Defines a rectangle's `width` and `height` relative to an `anchor` point at the top left corner.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rectangleShapeType(depth - 1),
			},
		},

		"union_shape": {
			Description: "Identifies a combined shape based on a set of other shapes. You can use a full JSON object to represent a union or an array of shapes that describe it.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: unionShapeType(depth - 1),
			},
		},
	}
}

func transformationType(depth int) map[string]*schema.Schema {
	if depth <= 0 {
		return nil
	}
	return map[string]*schema.Schema{

		"append": {
			Description: "Places a specified `image` beside the source image. The API places the `image` on a major dimension, then aligns it on the minor dimension. Transparent pixels fill any area not covered by either image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: appendType(depth - 1),
			},
		},

		"aspect_crop": {
			Description: "Lets you change the height or width of an image (either by cropping or expanding the area) to an aspect ratio of your choosing.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: aspectCrop(depth - 1),
			},
		},

		"background_color": {
			Description: "Places a transparent image on a set background color. Color is specified in the typical CSS hexadecimal format.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: backgroundColor(depth - 1),
			},
		},

		"blur": {
			Description: "Applies a Gaussian blur to the image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: blur(depth - 1),
			},
		},

		"chroma_key": {
			Description: "Changes any color in an image within the specified volume of the HSL colorspace to transparent or semitransparent. This transformation applies a 'green screen' technique commonly used to isolate and remove background colors.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: chromaKey(depth - 1),
			},
		},

		"composite": {
			Description: "Applies another image to the source image, either as an overlay or an underlay. The image that's underneath is visible in areas that are beyond the edges of the top image or that are less than 100% opaque. A common use of an overlay composite is to add a watermark.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: composite(depth - 1),
			},
		},

		"compound": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: compound(depth - 1),
			},
		},

		"contrast": {
			Description: "Adjusts both the contrast and brightness of an image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: contrast(depth - 1),
			},
		},

		"crop": {
			Description: "Crops an image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: crop(depth - 1),
			},
		},

		"face_crop": {
			Description: "Applies a method to detect faces in the source image and applies the rectangular crop on either the `biggest` face or `all` of the faces detected. Image and Video Manager tries to preserve faces in the image instead of using specified crop coordinates.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: faceCrop(depth - 1),
			},
		},

		"feature_crop": {
			Description: "Identifies prominent features of the source image, then crops around as many of these features as possible relative to the specified `width` and `height` values.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: featureCrop(depth - 1),
			},
		},

		"fit_and_fill": {
			Description: "Resizes an image to fit within a specific size box and then uses a fill of that same image to cover any transparent space at the edges. By default the fill image has a Blur transformation with a sigma value of 8 applied, but the transformation can be customized using the fillTransformation parameter.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: fitAndFill(depth - 1),
			},
		},

		"goop": {
			Description: "Distorts an image by randomly repositioning a set of control points along a specified grid. The transformed image appears _goopy_. Adjust the density of the grid and the degree of randomity. You can use this transformation to create watermarks for use in security.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: goop(depth - 1),
			},
		},

		"grayscale": {
			Description: "Restricts image color to shades of gray only.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: grayscale(depth - 1),
			},
		},

		"hsl": {
			Description: "Adjusts the hue, saturation, and lightness (HSL) of an image. Hue is the number of degrees that colors rotate around the color wheel. Saturation is a multiplier to increase or decrease color saturation. Lightness is a multiplier to increase or decrease the lightness of an image. Other transformations can also affect color, such as `Grayscale` and `MaxColors`. If youre using more than one, consider the order to apply them for the desired results.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: hsl(depth - 1),
			},
		},

		"hsv": {
			Description: "Identical to HSL except it replaces `lightness` with `value`. For example, if you reduce the `lightness` of a light green, almost white, image, the color turns a vibrant green. Reducing the `value` turns the image a darker color, close to grey. This happens because the original image color is very close to white.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: hsv(depth - 1),
			},
		},

		"if_dimension": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: ifDimension(depth - 1),
			},
		},

		"if_orientation": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: ifOrientation(depth - 1),
			},
		},

		"im_query": {
			Description: "Apply artistic transformations to images quickly and dynamically by specifying transformations with a query string appendedto the image URL.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: imQuery(depth - 1),
			},
		},

		"max_colors": {
			Description: "Set the maximum number of colors in the images palette. Reducing the number of colors in an image can help to reduce file size.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: maxColors(depth - 1),
			},
		},

		"mirror": {
			Description: "Flips an image horizontally, vertically, or both.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: mirror(depth - 1),
			},
		},

		"mono_hue": {
			Description: "Allows you to set all hues in an image to a single specified hue of your choosing. Mono Hue maintains the original color’s lightness and saturation but sets the hue to that of the specified value. This has the effect of making the image shades of the specified hue.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: monoHue(depth - 1),
			},
		},

		"opacity": {
			Description: "Adjusts the level of transparency of an image. Use this transformation to make an image more or less transparent.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: opacity(depth - 1),
			},
		},

		"region_of_interest_crop": {
			Description: "Crops to a region around a specified area of interest relative to the specified `width` and `height` values.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: regionOfInterestCrop(depth - 1),
			},
		},

		"relative_crop": {
			Description: "Shrinks or expands an image relative to the image's specified dimensions. Image and Video Manager fills the expanded areas with transparency. Positive values shrink the side, while negative values expand it.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: relativeCrop(depth - 1),
			},
		},

		"remove_color": {
			Description: "Removes a specified color from an image and replaces it with transparent pixels. This transformation is ideal for removing solid background colors from product images photographed on clean, consistent backgrounds without any shadows.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: removeColor(depth - 1),
			},
		},

		"resize": {
			Description: "Resizes an image to a particular, absolute dimension. If you don't enter a `width` or a `height`, the image is resized with the `fit` aspect preservation mode, which selects a value for the missing dimension that preserves the image's aspect.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: resize(depth - 1),
			},
		},

		"rotate": {
			Description: "Rotate the image around its center by indicating the degrees of rotation.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rotate(depth - 1),
			},
		},

		"scale": {
			Description: "Changes the image's size to different dimensions relative to its starting size.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: scale(depth - 1),
			},
		},

		"shear": {
			Description: "Slants an image into a parallelogram, as a percent of the starting dimension as represented in decimal format. You need to specify at least one axis property. Transparent pixels fill empty areas around the sheared image as needed, so it's often useful to use a `BackgroundColor` transformation for these areas.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: shear(depth - 1),
			},
		},

		"trim": {
			Description: "Automatically crops uniform backgrounds from the edges of an image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: trim(depth - 1),
			},
		},

		"unsharp_mask": {
			Description: "Emphasizes edges and details in source images without distorting the colors. Although this effect is often referred to as _sharpening_ an image, it actually creates a blurred, inverted copy of the image known as an unsharp mask. Image and Video Manager combines the unsharp mask with the source image to create an image perceived as clearer.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: unsharpMask(depth - 1),
			},
		},
	}
}

func transformationTypePost(depth int) map[string]*schema.Schema {
	if depth <= 0 {
		return nil
	}
	return map[string]*schema.Schema{

		"background_color": {
			Description: "Places a transparent image on a set background color. Color is specified in the typical CSS hexadecimal format.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: backgroundColor(depth - 1),
			},
		},

		"blur": {
			Description: "Applies a Gaussian blur to the image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: blur(depth - 1),
			},
		},

		"chroma_key": {
			Description: "Changes any color in an image within the specified volume of the HSL colorspace to transparent or semitransparent. This transformation applies a 'green screen' technique commonly used to isolate and remove background colors.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: chromaKey(depth - 1),
			},
		},

		"composite": {
			Description: "Applies another image to the source image, either as an overlay or an underlay. The image that's underneath is visible in areas that are beyond the edges of the top image or that are less than 100% opaque. A common use of an overlay composite is to add a watermark.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: compositePost(depth - 1),
			},
		},

		"compound": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: compoundPost(depth - 1),
			},
		},

		"contrast": {
			Description: "Adjusts both the contrast and brightness of an image.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: contrast(depth - 1),
			},
		},

		"goop": {
			Description: "Distorts an image by randomly repositioning a set of control points along a specified grid. The transformed image appears _goopy_. Adjust the density of the grid and the degree of randomity. You can use this transformation to create watermarks for use in security.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: goop(depth - 1),
			},
		},

		"grayscale": {
			Description: "Restricts image color to shades of gray only.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: grayscale(depth - 1),
			},
		},

		"hsl": {
			Description: "Adjusts the hue, saturation, and lightness (HSL) of an image. Hue is the number of degrees that colors rotate around the color wheel. Saturation is a multiplier to increase or decrease color saturation. Lightness is a multiplier to increase or decrease the lightness of an image. Other transformations can also affect color, such as `Grayscale` and `MaxColors`. If youre using more than one, consider the order to apply them for the desired results.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: hsl(depth - 1),
			},
		},

		"hsv": {
			Description: "Identical to HSL except it replaces `lightness` with `value`. For example, if you reduce the `lightness` of a light green, almost white, image, the color turns a vibrant green. Reducing the `value` turns the image a darker color, close to grey. This happens because the original image color is very close to white.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: hsv(depth - 1),
			},
		},

		"if_dimension": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: ifDimensionPost(depth - 1),
			},
		},

		"if_orientation": {
			Description: "",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: ifOrientationPost(depth - 1),
			},
		},

		"max_colors": {
			Description: "Set the maximum number of colors in the images palette. Reducing the number of colors in an image can help to reduce file size.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: maxColors(depth - 1),
			},
		},

		"mirror": {
			Description: "Flips an image horizontally, vertically, or both.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: mirror(depth - 1),
			},
		},

		"mono_hue": {
			Description: "Allows you to set all hues in an image to a single specified hue of your choosing. Mono Hue maintains the original color’s lightness and saturation but sets the hue to that of the specified value. This has the effect of making the image shades of the specified hue.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: monoHue(depth - 1),
			},
		},

		"opacity": {
			Description: "Adjusts the level of transparency of an image. Use this transformation to make an image more or less transparent.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: opacity(depth - 1),
			},
		},

		"remove_color": {
			Description: "Removes a specified color from an image and replaces it with transparent pixels. This transformation is ideal for removing solid background colors from product images photographed on clean, consistent backgrounds without any shadows.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: removeColor(depth - 1),
			},
		},

		"unsharp_mask": {
			Description: "Emphasizes edges and details in source images without distorting the colors. Although this effect is often referred to as _sharpening_ an image, it actually creates a blurred, inverted copy of the image known as an unsharp mask. Image and Video Manager combines the unsharp mask with the source image to create an image perceived as clearer.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: unsharpMask(depth - 1),
			},
		},
	}
}

// PolicyOutputImage is a top level schema func
func PolicyOutputImage(depth int) map[string]*schema.Schema {
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
			Description: "Dictates the output quality (either `quality` or `perceptualQuality`) and formats that are created for each resized image. If unspecified, image formats are created to support all browsers at the default quality level (`85`), which includes formats such as WEBP, JPEG2000 and JPEG-XR for specific browsers.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: outputImage(depth - 1),
			},
		},
		"post_breakpoint_transformations": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Post-processing Transformations are applied to the image after image and quality settings have been applied.",
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"rollout_duration": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount of time in seconds that the policy takes to rollout. During the rollout an increasing proportion of images/videos will begin to use the new policy instead of the cached images/videos from the previous version. This value has no effect on the staging network.",
			ValidateDiagFunc: stringAsIntBetween(3600, 604800),
		},
		"serve_stale_duration": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount of time in seconds that the policy will serve stale images. During the serve stale period realtime images will attempt to use the offline image from the previous policy version first if possible.",
			ValidateDiagFunc: stringAsIntBetween(0, 2592000),
		},
		"transformations": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Set of image transformations to apply to the source image. If unspecified, no operations are performed.",
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
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

func appendType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies where to place the `image` relative to the source image. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and a `Center` by default.",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_priority": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Determines the exact placement of the `image` when `gravity` is `Center` or a diagonal. The value is either `horizontal` or `vertical`. Use `horizontal` to append an `image` east or west of the source image. This aligns the `image` on the vertical gravity component, placing `Center` gravity east. Use `vertical` to append an `image` north or south of the source image. This aligns the `image` on the horizontal gravity component, placing `Center` gravity south.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"horizontal", "vertical"}, false)),
		},
		"gravity_priority_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Determines the exact placement of the `image` when `gravity` is `Center` or a diagonal. The value is either `horizontal` or `vertical`. Use `horizontal` to append an `image` east or west of the source image. This aligns the `image` on the vertical gravity component, placing `Center` gravity east. Use `vertical` to append an `image` north or south of the source image. This aligns the `image` on the horizontal gravity component, placing `Center` gravity south.",
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies where to place the `image` relative to the source image. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and a `Center` by default.",
		},
		"image": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: imageType(depth - 1),
			},
		},
		"preserve_minor_dimension": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Whether to preserve the source image's minor dimension, `false` by default. The minor dimension is the dimension opposite the dimension that the appending `image` is placed. For example, the dimensions of the source image are 100 &times; 100 pixels. The dimensions of the appending `image` are 50 &times; 150 pixels. The `gravity` is set to `East`. This makes the major dimension horizontal and the source image's minor dimension vertical. To preserve the source image's minor dimension at 100 pixels, the `preserveMinorDimension` is set to `true`. As a result of the append, the major dimension expanded with the appended image to 150 pixels. The source image's minor dimension was maintained at 100 pixels. The total combined dimension of the image is 150 &times; 100 pixels.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"preserve_minor_dimension_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Whether to preserve the source image's minor dimension, `false` by default. The minor dimension is the dimension opposite the dimension that the appending `image` is placed. For example, the dimensions of the source image are 100 &times; 100 pixels. The dimensions of the appending `image` are 50 &times; 150 pixels. The `gravity` is set to `East`. This makes the major dimension horizontal and the source image's minor dimension vertical. To preserve the source image's minor dimension at 100 pixels, the `preserveMinorDimension` is set to `true`. As a result of the append, the major dimension expanded with the appended image to 150 pixels. The source image's minor dimension was maintained at 100 pixels. The total combined dimension of the image is 150 &times; 100 pixels.",
		},
	}
}

func aspectCrop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_expansion": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Increases the size of the image canvas to achieve the requested aspect ratio instead of cropping the image. Use the Horizontal Offset and Vertical Offset settings to determine where to add the fully transparent pixels on the expanded image canvas.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"allow_expansion_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Increases the size of the image canvas to achieve the requested aspect ratio instead of cropping the image. Use the Horizontal Offset and Vertical Offset settings to determine where to add the fully transparent pixels on the expanded image canvas.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height term of the aspect ratio to crop.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height term of the aspect ratio to crop.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width term of the aspect ratio to crop.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width term of the aspect ratio to crop.",
		},
		"x_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"x_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.",
		},
		"y_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"y_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the horizontal portion of the image you want to keep when the aspect ratio cropping is applied. When using Allow Expansion this setting defines the horizontal position of the image on the new expanded image canvas.",
		},
	}
}

func backgroundColor(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The hexadecimal CSS color value for the background.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The hexadecimal CSS color value for the background.",
		},
	}
}

func blur(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sigma": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to scatter the original pixel by to create the blur effect. Resulting images may be larger than the original as pixels at the edge of the image might scatter outside the image's original dimensions.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"sigma_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to scatter the original pixel by to create the blur effect. Resulting images may be larger than the original as pixels at the edge of the image might scatter outside the image's original dimensions.",
		},
	}
}

func boxImageType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height of the box in pixels.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width of the box in pixels.",
		},
	}
}

func boxImageTypePost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height of the box in pixels.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width of the box in pixels.",
		},
	}
}

func breakpoints(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"widths": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeInt},
		},
	}
}

func chromaKey(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hue": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The hue to remove. Enter the degree of rotation between 0 and 360 degrees around the color wheel. By default Chroma Key removes a green hue, 120° on the color wheel.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_feather": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How much additional hue to make semi-transparent beyond the Hue Tolerance. By default Hue Feather is 0.083 which applies semi-transparency to hues 30° around the Hue Tolerance.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_feather_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How much additional hue to make semi-transparent beyond the Hue Tolerance. By default Hue Feather is 0.083 which applies semi-transparency to hues 30° around the Hue Tolerance.",
		},
		"hue_tolerance": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How close a color's hue needs to be to the selected hue for it to be changed to fully transparent. If you enter the maximum value of 1.0 the entire image is made transparent. By default Hue Tolerance is approximately 0.083 or 8.3% of the color wheel. This value corresponds to 30° around the specified hue.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_tolerance_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How close a color's hue needs to be to the selected hue for it to be changed to fully transparent. If you enter the maximum value of 1.0 the entire image is made transparent. By default Hue Tolerance is approximately 0.083 or 8.3% of the color wheel. This value corresponds to 30° around the specified hue.",
		},
		"hue_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The hue to remove. Enter the degree of rotation between 0 and 360 degrees around the color wheel. By default Chroma Key removes a green hue, 120° on the color wheel.",
		},
		"lightness_feather": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How much additional lightness to make semi-transparent beyond the Lightness Tolerance. The default value of 0.1 corresponds to 10% away from the tolerated lightness towards full black or full white.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"lightness_feather_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How much additional lightness to make semi-transparent beyond the Lightness Tolerance. The default value of 0.1 corresponds to 10% away from the tolerated lightness towards full black or full white.",
		},
		"lightness_tolerance": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How much of the lightest part and darkest part of a color to preserve. For example, you can space this value out from the middle (i.e. 0.5 lightness or full color) to help preserve the splash lighting impact in the image. You can define how close the color needs to be to the full color to remove it from your image. The default value of 0.75 means that a colour must be within 75% of the full color to full white or full black for full removal.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"lightness_tolerance_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How much of the lightest part and darkest part of a color to preserve. For example, you can space this value out from the middle (i.e. 0.5 lightness or full color) to help preserve the splash lighting impact in the image. You can define how close the color needs to be to the full color to remove it from your image. The default value of 0.75 means that a colour must be within 75% of the full color to full white or full black for full removal.",
		},
		"saturation_feather": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How much additional saturation to make semi-transparent beyond the Saturation Tolerance. By default Saturation Feather is 0.1 which applies semi-transparency to hues 10% below the saturationTolerance.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"saturation_feather_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How much additional saturation to make semi-transparent beyond the Saturation Tolerance. By default Saturation Feather is 0.1 which applies semi-transparency to hues 10% below the saturationTolerance.",
		},
		"saturation_tolerance": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "How close a color's saturation needs to be to full saturation for it to be changed to fully transparent. For example, you can define how green the color needs to be to remove it from your image. The default value of 0.75 means that a color must be within 75% of full saturation in order to be made fully transparent.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"saturation_tolerance_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How close a color's saturation needs to be to full saturation for it to be changed to fully transparent. For example, you can define how green the color needs to be to remove it from your image. The default value of 0.75 means that a color must be within 75% of full saturation in order to be made fully transparent.",
		},
	}
}

func circleImageType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"diameter": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The diameter of the circle. The diameter will be the width and the height of the image in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"diameter_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The diameter of the circle. The diameter will be the width and the height of the image in pixels.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width of the box in pixels.",
		},
	}
}

func circleImageTypePost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.",
		},
		"diameter": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The diameter of the circle. The diameter will be the width and the height of the image in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"diameter_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The diameter of the circle. The diameter will be the width and the height of the image in pixels.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width of the box in pixels.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width of the box in pixels.",
		},
	}
}

func circleShapeType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"center": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Defines coordinates for a single point, to help define polygons and rectangles. Each point may be an object with `x`and `y` members, or a two-element array.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: pointShapeType(depth - 1),
			},
		},
		"radius": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The radius of the circle measured in pixels.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"radius_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The radius of the circle measured in pixels.",
		},
	}
}

func composite(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position",
		},
		"image": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: imageType(depth - 1),
			},
		},
		"placement": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Over", "Under", "Mask", "Stencil"}, false)),
		},
		"placement_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.",
		},
		"scale": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"scale_dimension": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The dimension, either `width` or `height`, of the source image to scale.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"width", "height"}, false)),
		},
		"scale_dimension_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimension, either `width` or `height`, of the source image to scale.",
		},
		"scale_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.",
		},
		"x_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The x-axis position of the image to apply.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"x_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The x-axis position of the image to apply.",
		},
		"y_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The y-axis position of the image to apply.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"y_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The y-axis position of the image to apply.",
		},
	}
}

func compositePost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position",
			ValidateDiagFunc: validateGravityPost(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Compass direction indicating the corner or edge of the base image to place the applied image. Use Horizontal and Vertical Offset to adjust the applied image's gravity position",
		},
		"image": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: imageTypePost(depth - 1),
			},
		},
		"placement": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Over", "Under", "Mask", "Stencil"}, false)),
		},
		"placement_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Place applied image on top of or underneath the base image. Watermarks are usually applied over. Backgrounds are usually applied under.",
		},
		"scale": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"scale_dimension": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The dimension, either `width` or `height`, of the source image to scale.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"width", "height"}, false)),
		},
		"scale_dimension_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimension, either `width` or `height`, of the source image to scale.",
		},
		"scale_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to resize the applied image relative to the source image and preserve aspect ratio, 1 by default. Set the `scaleDimension` to calculate the `scale` from the source image's width or height.",
		},
		"x_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The x-axis position of the image to apply.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"x_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The x-axis position of the image to apply.",
		},
		"y_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The y-axis position of the image to apply.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"y_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The y-axis position of the image to apply.",
		},
	}
}

func compound(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformations": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
	}
}

func compoundPost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformations": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
	}
}

func contrast(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"brightness": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Adjusts the brightness of the image. Positive values increase brightness and negative values decrease brightness. A value of  `1` produces a white image. A value of  `-1` produces a black image. The default value is `0`, which leaves the image unchanged. The acceptable value range is `-1.0` to `1.0`. Values outside of the acceptable range clamp to this range.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"brightness_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Adjusts the brightness of the image. Positive values increase brightness and negative values decrease brightness. A value of  `1` produces a white image. A value of  `-1` produces a black image. The default value is `0`, which leaves the image unchanged. The acceptable value range is `-1.0` to `1.0`. Values outside of the acceptable range clamp to this range.",
		},
		"contrast": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Adjusts the contrast of the image. Expressed as a range from `-1` to `1`, positive values increase contrast, negative values decrease it, while `0` leaves the image unchanged. Values outside of the `-1` to `1` range clamp to this range.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"contrast_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Adjusts the contrast of the image. Expressed as a range from `-1` to `1`, positive values increase contrast, negative values decrease it, while `0` leaves the image unchanged. Values outside of the `-1` to `1` range clamp to this range.",
		},
	}
}

func crop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_expansion": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "If cropping an area outside of the existing canvas, expands the image canvas.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"allow_expansion_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "If cropping an area outside of the existing canvas, expands the image canvas.",
		},
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Frame of reference for X and Y Positions.",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Frame of reference for X and Y Positions.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to crop along the y-axis.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to crop along the y-axis.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to crop along the x-axis.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to crop along the x-axis.",
		},
		"x_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The x-axis position of the image to crop from.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"x_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The x-axis position of the image to crop from.",
		},
		"y_position": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The y-axis position of the image to crop from.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"y_position_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The y-axis position of the image to crop from.",
		},
	}
}

func enumOptions(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The unique identifier for each enum value, up to 50 alphanumeric characters.",
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The value of the variable when the `id` is provided.",
		},
	}
}

func faceCrop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"algorithm": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies the type of algorithm used to detect faces in the image, either `cascade` for the cascade classifier algorithm or `dnn` for the deep neural network algorithm, `cascade` by default.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"cascade", "dnn"}, false)),
		},
		"algorithm_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the type of algorithm used to detect faces in the image, either `cascade` for the cascade classifier algorithm or `dnn` for the deep neural network algorithm, `cascade` by default.",
		},
		"confidence": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "With `algorithm` set to `dnn`, specifies the minimum confidence needed to detect faces in the image. Values range from `0` to `1` for increased confidence, and possibly fewer faces detected.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"confidence_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "With `algorithm` set to `dnn`, specifies the minimum confidence needed to detect faces in the image. Values range from `0` to `1` for increased confidence, and possibly fewer faces detected.",
		},
		"fail_gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Controls placement of the crop if Image and Video Manager does not detect any faces in the image. Directions are relative to the edges of the image being transformed.",
			ValidateDiagFunc: validateGravity(),
		},
		"fail_gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls placement of the crop if Image and Video Manager does not detect any faces in the image. Directions are relative to the edges of the image being transformed.",
		},
		"focus": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Distinguishes the faces detected, either `biggestFace` or `allFaces` to place the crop rectangle around the full set of faces, `all` by default.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"allFaces", "biggestFace"}, false)),
		},
		"focus_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Distinguishes the faces detected, either `biggestFace` or `allFaces` to place the crop rectangle around the full set of faces, `all` by default.",
		},
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Controls placement of the crop. Directions are relative to the face(s) plus padding.",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls placement of the crop. Directions are relative to the face(s) plus padding.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height of the output image in pixels relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height of the output image in pixels relative to the specified `style` value.",
		},
		"padding": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The padding ratio based on the dimensions of the biggest face detected, `0.5` by default. Larger values increase padding.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"padding_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The padding ratio based on the dimensions of the biggest face detected, `0.5` by default. Larger values increase padding.",
		},
		"style": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies how to crop or scale a crop area for the faces detected in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places a raw crop around the faces, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and faces as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the faces, relative to the specified `width` and `height` values. Allows Variable substitution.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"crop", "fill", "zoom"}, false)),
		},
		"style_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies how to crop or scale a crop area for the faces detected in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places a raw crop around the faces, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and faces as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the faces, relative to the specified `width` and `height` values. Allows Variable substitution.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width of the output image in pixels relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width of the output image in pixels relative to the specified `style` value.",
		},
	}
}

func featureCrop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fail_gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Controls placement of the crop if Image and Video Manager does not detect any features in the image. Directions are relative to the edges of the image being transformed.",
			ValidateDiagFunc: validateGravity(),
		},
		"fail_gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls placement of the crop if Image and Video Manager does not detect any features in the image. Directions are relative to the edges of the image being transformed.",
		},
		"feature_radius": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The size in pixels of the important features to search for. If identified, two features never appear closer together than this value, `8.0` by default.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"feature_radius_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The size in pixels of the important features to search for. If identified, two features never appear closer together than this value, `8.0` by default.",
		},
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Controls placement of the crop. Directions are relative to the region of interest plus padding.",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls placement of the crop. Directions are relative to the region of interest plus padding.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height in pixels of the output image relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height in pixels of the output image relative to the specified `style` value.",
		},
		"max_features": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The maximum number of features to identify as important features, `32` by default. The strongest features are always chosen.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"max_features_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The maximum number of features to identify as important features, `32` by default. The strongest features are always chosen.",
		},
		"min_feature_quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Determines the minimum quality level of the feature identified. To consider a feature important, the feature needs to surpass this value.  Image and Video Manager measures quality on a scale from `0` for the lowest quality to `1` for the highest quality, `.1` by default.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"min_feature_quality_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Determines the minimum quality level of the feature identified. To consider a feature important, the feature needs to surpass this value.  Image and Video Manager measures quality on a scale from `0` for the lowest quality to `1` for the highest quality, `.1` by default.",
		},
		"padding": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Adds space around the region of interest. The amount of padding added is directly related to the size of the bounding box of the selected features. Specifically, the region of interest is expanded in all directions by the largest dimension of the bounding box of the selected features multiplied by this value.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"padding_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Adds space around the region of interest. The amount of padding added is directly related to the size of the bounding box of the selected features. Specifically, the region of interest is expanded in all directions by the largest dimension of the bounding box of the selected features multiplied by this value.",
		},
		"style": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies how to crop or scale a crop area for the features identified in the source image, `fill` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` performs a raw crop around the features, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and features as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the features, relative to the specified `width` and `height` values. Allows Variable substitution.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"crop", "fill", "zoom"}, false)),
		},
		"style_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies how to crop or scale a crop area for the features identified in the source image, `fill` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` performs a raw crop around the features, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and features as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the features, relative to the specified `width` and `height` values. Allows Variable substitution.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width in pixels of the output image relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width in pixels of the output image relative to the specified `style` value.",
		},
	}
}

func fitAndFill(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fill_transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height value of the resized image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height value of the resized image.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width value of the resized image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width value of the resized image.",
		},
	}
}

func goop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"chaos": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies the greatest distance control points may move from their original position. A value of `1.0` shifts control points over as far as the next one in the original grid. A value of `0.0` leaves the image unchanged. Values under `0.5` work better for subtle distortions, otherwise control points may pass each other and cause a twisting effect.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"chaos_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the greatest distance control points may move from their original position. A value of `1.0` shifts control points over as far as the next one in the original grid. A value of `0.0` leaves the image unchanged. Values under `0.5` work better for subtle distortions, otherwise control points may pass each other and cause a twisting effect.",
		},
		"density": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Controls the density of control points used to distort the image. The largest dimension of the input image is divided up to fit this number of control points. A grid of points is extended on the smaller dimension such that each row and column of control points is equidistant from each adjacent row or column. This parameter strongly affects transformation performance. Be careful choosing values above the default if you expect to transform medium to large size images.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"density_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls the density of control points used to distort the image. The largest dimension of the input image is divided up to fit this number of control points. A grid of points is extended on the smaller dimension such that each row and column of control points is equidistant from each adjacent row or column. This parameter strongly affects transformation performance. Be careful choosing values above the default if you expect to transform medium to large size images.",
		},
		"power": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "By default, the distortion algorithm relies on inverse squares to calculate distance but this allows you to change the exponent. You shouldnt need to vary the default value of `2.0`.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"power_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "By default, the distortion algorithm relies on inverse squares to calculate distance but this allows you to change the exponent. You shouldnt need to vary the default value of `2.0`.",
		},
		"seed": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies your own `seed` value as an alternative to the default, which is subject to variability. This allows for reproducible and deterministic distortions. If all parameters are kept equal and a constant seed is used, `Goop` distorts an input image consistently over many transformations. By default, this value is set to the current Epoch Time measured in milliseconds, which provides inconsistent transformation output.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"seed_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies your own `seed` value as an alternative to the default, which is subject to variability. This allows for reproducible and deterministic distortions. If all parameters are kept equal and a constant seed is used, `Goop` distorts an input image consistently over many transformations. By default, this value is set to the current Epoch Time measured in milliseconds, which provides inconsistent transformation output.",
		},
	}
}

func grayscale(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The algorithm used to transform colors to grays, either `Brightness`, `Lightness`, `Rec601`, or the default `Rec709`.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"Rec601", "Rec709", "Brightness", "Lightness"}, false)),
		},
		"type_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The algorithm used to transform colors to grays, either `Brightness`, `Lightness`, `Rec601`, or the default `Rec709`.",
		},
	}
}

func hsl(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hue": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of degrees to rotate colors around the color wheel, `0` by default.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of degrees to rotate colors around the color wheel, `0` by default.",
		},
		"lightness": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to adjust the lightness of colors in the image. Note that lightness is distinct from brightness. For example, reducing the lightness of a light green might give you a lime green whereas reducing the brightness of a light green might give you a darker shade of the same green. Values less than `1.0` decrease the lightness of colors in the image. Values greater than `1.0` increase the lightness of colors in the image.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"lightness_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to adjust the lightness of colors in the image. Note that lightness is distinct from brightness. For example, reducing the lightness of a light green might give you a lime green whereas reducing the brightness of a light green might give you a darker shade of the same green. Values less than `1.0` decrease the lightness of colors in the image. Values greater than `1.0` increase the lightness of colors in the image.",
		},
		"saturation": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"saturation_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.",
		},
	}
}

func hsv(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hue": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of degrees to rotate colors around the color wheel, `0.0` by default.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of degrees to rotate colors around the color wheel, `0.0` by default.",
		},
		"saturation": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"saturation_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image.",
		},
		"value": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "A multiplier to adjust the lightness or darkness of the images base color. Values less than 1.0 decrease the base colors in the image, making them appear darker. Values greater than 1.0 increase the base colors in the image, making them appear lighter.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"value_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A multiplier to adjust the lightness or darkness of the images base color. Values less than 1.0 decrease the base colors in the image, making them appear darker. Values greater than 1.0 increase the base colors in the image, making them appear lighter.",
		},
	}
}

func ifDimension(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"dimension": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The dimension to use to select the transformation, either `height`, `width`, or `both`.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"width", "height", "both"}, false)),
		},
		"dimension_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimension to use to select the transformation, either `height`, `width`, or `both`.",
		},
		"equal": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"greater_than": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"less_than": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"value": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"value_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.",
		},
	}
}

func ifDimensionPost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"dimension": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The dimension to use to select the transformation, either `height`, `width`, or `both`.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"width", "height", "both"}, false)),
		},
		"dimension_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The dimension to use to select the transformation, either `height`, `width`, or `both`.",
		},
		"equal": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"greater_than": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"less_than": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"value": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"value_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare against the source image dimension. For example, if the image dimension is less than the value the lessThan transformation is applied.",
		},
	}
}

func ifOrientation(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"landscape": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"portrait": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"square": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
	}
}

func ifOrientationPost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"landscape": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"portrait": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"square": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
	}
}

func imQuery(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_transformations": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Specifies the transformations that can be applied using the query string parameter.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"query_var": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func maxColors(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"colors": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The value representing the maximum number of colors to use with the source image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"colors_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value representing the maximum number of colors to use with the source image.",
		},
	}
}

func mirror(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"horizontal": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Flips the image horizontally.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"horizontal_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Flips the image horizontally.",
		},
		"vertical": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Flips the image vertically.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"vertical_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Flips the image vertically.",
		},
	}
}

func monoHue(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hue": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specify a hue by indicating the degree of rotation between 0 and 360 degrees around the color wheel. By default Mono Hue applies a red hue, 0.0 on the color wheel.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"hue_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify a hue by indicating the degree of rotation between 0 and 360 degrees around the color wheel. By default Mono Hue applies a red hue, 0.0 on the color wheel.",
		},
	}
}

func opacity(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"opacity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Represents alpha values on a scale of `0` to `1`. Values below `1` increase transparency, and `0` is invisible. For images that have some transparency, values above `1` increase the opacity of the transparent portions.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"opacity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Represents alpha values on a scale of `0` to `1`. Values below `1` increase transparency, and `0` is invisible. For images that have some transparency, values above `1` increase the opacity of the transparent portions.",
		},
	}
}

func outputImage(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"adaptive_quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Override the quality of image to serve when Image & Video Manager detects a slow connection. Specifying lower values lets users with slow connections browse your site with reduced load times without impacting the quality of images for users with faster connections.",
			ValidateDiagFunc: stringAsIntBetween(1, 100),
		},
		"allow_pristine_on_downsize": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Whether a pristine image wider than the requested breakpoint is allowed as a derivative image if it has the fewest bytes. This will not have an affect if transformations are present.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"allowed_formats": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The graphics file formats allowed for browser specific results.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"forced_formats": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The forced extra formats for the `imFormat` query parameter, which requests a specific browser type. By default, Image and Video Manager detects the browser and returns the appropriate image.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"perceptual_quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Mutually exclusive with quality. The perceptual quality to use when comparing resulting images, which overrides the `quality` setting. Perceptual quality tunes each image format's quality parameter dynamically based on the human-perceived quality of the output image. This can result in better byte savings (as compared to using regular quality) as many images can be encoded at a much lower quality without compromising perception of the image. In addition, certain images may need to be encoded at a slightly higher quality in order to maintain human-perceived quality. Values are tiered `high`, `mediumHigh`, `medium`, `mediumLow`, or `low`.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"high", "mediumHigh", "medium", "mediumLow", "low"}, false)),
		},
		"perceptual_quality_floor": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Only applies with perceptualQuality set. Sets a minimum image quality to respect when using perceptual quality. Perceptual quality will not reduce the quality below this value even if it determines the compressed image to be acceptably visually similar.",
			ValidateDiagFunc: stringAsIntBetween(1, 100),
		},
		"perceptual_quality_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Mutually exclusive with quality. The perceptual quality to use when comparing resulting images, which overrides the `quality` setting. Perceptual quality tunes each image format's quality parameter dynamically based on the human-perceived quality of the output image. This can result in better byte savings (as compared to using regular quality) as many images can be encoded at a much lower quality without compromising perception of the image. In addition, certain images may need to be encoded at a slightly higher quality in order to maintain human-perceived quality. Values are tiered `high`, `mediumHigh`, `medium`, `mediumLow`, or `low`.",
		},
		"prefer_modern_formats": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Whether derivative image formats should be selected with a preference for modern formats (such as WebP and Avif) instead the format that results in the fewest bytes.",
			ValidateDiagFunc: validateIsTypeBool(),
		},
		"quality": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Mutually exclusive with perceptualQuality, used by default if neither is specified. The chosen quality of the output images. Using a quality value from 1-100 resembles JPEG quality across output formats.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"quality_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Mutually exclusive with perceptualQuality, used by default if neither is specified. The chosen quality of the output images. Using a quality value from 1-100 resembles JPEG quality across output formats.",
		},
	}
}

func pointShapeType(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"x": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The horizontal position of the point, measured in pixels.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"x_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The horizontal position of the point, measured in pixels.",
		},
		"y": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The vertical position of the point, measured in pixels.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"y_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The vertical position of the point, measured in pixels.",
		},
	}
}

func polygonShapeType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"points": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Series of [PointShapeType](#pointshapetype) objects. The last and first points connect to close the shape automatically.",
			Elem: &schema.Resource{
				Schema: pointShapeType(depth - 1),
			},
		},
	}
}

func rectangleShapeType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"anchor": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: pointShapeType(depth - 1),
			},
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Extends the rectangle down from the `anchor` point.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Extends the rectangle down from the `anchor` point.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Extends the rectangle right from the `anchor` point.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Extends the rectangle right from the `anchor` point.",
		},
	}
}

func regionOfInterestCrop(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gravity": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The placement of the crop area relative to the specified area of interest.",
			ValidateDiagFunc: validateGravity(),
		},
		"gravity_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The placement of the crop area relative to the specified area of interest.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height in pixels of the output image relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height in pixels of the output image relative to the specified `style` value.",
		},
		"region_of_interest": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: shapeType(depth - 1),
			},
		},
		"style": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Specifies how to crop or scale a crop area for the specified area of interest in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places raw crop around the point of interest, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and point of interest as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the point of interest, relative to the specified `width` and `height` values.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"crop", "fill", "zoom"}, false)),
		},
		"style_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies how to crop or scale a crop area for the specified area of interest in the source image, `zoom` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` places raw crop around the point of interest, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and point of interest as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the point of interest, relative to the specified `width` and `height` values.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width in pixels of the output image relative to the specified `style` value.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width in pixels of the output image relative to the specified `style` value.",
		},
	}
}

func relativeCrop(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"east": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to shrink or expand the right side of the image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"east_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to shrink or expand the right side of the image.",
		},
		"north": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to shrink or expand the top side of the image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"north_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to shrink or expand the top side of the image.",
		},
		"south": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to shrink or expand the bottom side of the image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"south_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to shrink or expand the bottom side of the image.",
		},
		"west": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The number of pixels to shrink or expand the left side of the image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"west_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The number of pixels to shrink or expand the left side of the image.",
		},
	}
}

func removeColor(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The hexadecimal CSS color value to remove.",
		},
		"color_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The hexadecimal CSS color value to remove.",
		},
		"feather": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The RemoveColor transformation may create a hard edge around an image. To minimize these hard edges and make the removal of the color more gradual in appearance, use the Feather option. This option allows you to extend the color removal beyond the specified Tolerance. The pixels in this extended tolerance become semi-transparent - creating a softer edge.  The first realtime request for an image using the feather option may result in a slow transformation time. Subsequent requests are not impacted as they are served directly out of cache.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"feather_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The RemoveColor transformation may create a hard edge around an image. To minimize these hard edges and make the removal of the color more gradual in appearance, use the Feather option. This option allows you to extend the color removal beyond the specified Tolerance. The pixels in this extended tolerance become semi-transparent - creating a softer edge.  The first realtime request for an image using the feather option may result in a slow transformation time. Subsequent requests are not impacted as they are served directly out of cache.",
		},
		"tolerance": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The Tolerance option defines how close the color needs to be to the selected color before it's changed to fully transparent. Set the Tolerance to 0.0 to remove only the exact color specified.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"tolerance_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Tolerance option defines how close the color needs to be to the selected color before it's changed to fully transparent. Set the Tolerance to 0.0 to remove only the exact color specified.",
		},
	}
}

func resize(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aspect": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Preserves the aspect ratio. Select `fit` to make the image fit entirely within the selected width and height. When using `fit`, the resulting image has the largest possible size for the specified dimensions. Select `fill` to size the image so it both completely fills the dimensions and has the smallest possible file size. Otherwise `ignore` changes the original aspect ratio to fit within an arbitrarily shaped rectangle.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"fit", "fill", "ignore"}, false)),
		},
		"aspect_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Preserves the aspect ratio. Select `fit` to make the image fit entirely within the selected width and height. When using `fit`, the resulting image has the largest possible size for the specified dimensions. Select `fill` to size the image so it both completely fills the dimensions and has the smallest possible file size. Otherwise `ignore` changes the original aspect ratio to fit within an arbitrarily shaped rectangle.",
		},
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The height to resize the source image to. Must be set if height is not specified.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The height to resize the source image to. Must be set if height is not specified.",
		},
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Sets constraints for the image resize. Select `normal` to resize in all cases, either increasing or decreasing the dimensions. Select `downsize` to ignore this transformation if the result would be larger than the original. Select `upsize` to ignore this transformation if the result would be smaller.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"normal", "upsize", "downsize"}, false)),
		},
		"type_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Sets constraints for the image resize. Select `normal` to resize in all cases, either increasing or decreasing the dimensions. Select `downsize` to ignore this transformation if the result would be larger than the original. Select `upsize` to ignore this transformation if the result would be smaller.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The width to resize the source image to. Must be set if width is not specified.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The width to resize the source image to. Must be set if width is not specified.",
		},
	}
}

func rotate(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"degrees": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The value to rotate the image by. Positive values rotate clockwise, while negative values rotate counter-clockwise.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"degrees_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to rotate the image by. Positive values rotate clockwise, while negative values rotate counter-clockwise.",
		},
	}
}

func scale(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"height": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Scaling factor for the input height to determine the output height of the image, where values between `0` and `1` decrease size. Image dimensions need to be non-zero positive numbers.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"height_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Scaling factor for the input height to determine the output height of the image, where values between `0` and `1` decrease size. Image dimensions need to be non-zero positive numbers.",
		},
		"width": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Scaling factor for the input width to determine the output width of the image, where `1` leaves the width unchanged. Values greater than `1` increase the image size. Image dimensions need to be non-zero positive numbers.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"width_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Scaling factor for the input width to determine the output width of the image, where `1` leaves the width unchanged. Values greater than `1` increase the image size. Image dimensions need to be non-zero positive numbers.",
		},
	}
}

func shear(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"x_shear": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount to shear along the x-axis, measured in multiples of the image's width. Must be set if yShear is not specified.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"x_shear_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The amount to shear along the x-axis, measured in multiples of the image's width. Must be set if yShear is not specified.",
		},
		"y_shear": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount to shear along the y-axis, measured in multiples of the image's height. Must be set if xShear is not specified.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"y_shear_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The amount to shear along the y-axis, measured in multiples of the image's height. Must be set if xShear is not specified.",
		},
	}
}

func textImageType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fill": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The main fill color of the text.",
		},
		"fill_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The main fill color of the text.",
		},
		"size": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The size in pixels to render the text.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"size_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The size in pixels to render the text.",
		},
		"stroke": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The color for the outline of the text.",
		},
		"stroke_size": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The thickness in points for the outline of the text.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"stroke_size_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The thickness in points for the outline of the text.",
		},
		"stroke_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The color for the outline of the text.",
		},
		"text": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The line of text to render.",
		},
		"text_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The line of text to render.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"typeface": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.",
		},
		"typeface_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.",
		},
	}
}

func textImageTypePost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fill": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The main fill color of the text.",
		},
		"fill_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The main fill color of the text.",
		},
		"size": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The size in pixels to render the text.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"size_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The size in pixels to render the text.",
		},
		"stroke": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The color for the outline of the text.",
		},
		"stroke_size": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The thickness in points for the outline of the text.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"stroke_size_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The thickness in points for the outline of the text.",
		},
		"stroke_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The color for the outline of the text.",
		},
		"text": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The line of text to render.",
		},
		"text_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The line of text to render.",
		},
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"typeface": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.",
		},
		"typeface_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts.",
		},
	}
}

func trim(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fuzz": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The fuzz tolerance of the trim, a value between `0` and `1` that determines the acceptable amount of background variation before trimming stops.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"fuzz_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The fuzz tolerance of the trim, a value between `0` and `1` that determines the acceptable amount of background variation before trimming stops.",
		},
		"padding": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The amount of padding in pixels to add to the trimmed image.",
			ValidateDiagFunc: validateIsTypeInt(),
		},
		"padding_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The amount of padding in pixels to add to the trimmed image.",
		},
	}
}

func unionShapeType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"shapes": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: shapeType(depth - 1),
			},
		},
	}
}

func unsharpMask(_ int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gain": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Set how much emphasis the filter applies to details. Higher values increase apparent sharpness of details.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"gain_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set how much emphasis the filter applies to details. Higher values increase apparent sharpness of details.",
		},
		"sigma": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "The standard deviation of the Gaussian distribution used in the in unsharp mask, measured in pixels, `1.0` by default. High values emphasize large details and low values emphasize small details.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"sigma_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The standard deviation of the Gaussian distribution used in the in unsharp mask, measured in pixels, `1.0` by default. High values emphasize large details and low values emphasize small details.",
		},
		"threshold": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      "Set the minimum change required to include a detail in the filter. Higher values discard more changes.",
			ValidateDiagFunc: validateIsTypeFloat(),
		},
		"threshold_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set the minimum change required to include a detail in the filter. Higher values discard more changes.",
		},
	}
}

func urlImageType(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationType(depth - 1),
			},
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL of the image.",
		},
		"url_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL of the image.",
		},
	}
}

func urlImageTypePost(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transformation": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: transformationTypePost(depth - 1),
			},
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL of the image.",
		},
		"url_var": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL of the image.",
		},
	}
}

func variable(depth int) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The default value of the variable if no query parameter is provided. It needs to be one of the `enumOptions` if any are provided.",
		},
		"enum_options": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: enumOptions(depth - 1),
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the variable, also available as the query parameter name to set the variable's value dynamically. Use up to 50 alphanumeric characters.",
		},
		"postfix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A postfix added to the value provided for the variable, or to the default value.",
		},
		"prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A prefix added to the value provided for the variable, or to the default value.",
		},
		"type": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The type of value for the variable.",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"bool", "number", "url", "color", "gravity", "placement", "scaleDimension", "grayscaleType", "aspect", "resizeType", "dimension", "perceptualQuality", "string", "focus"}, false)),
		},
	}
}

func validateGravity() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"North", "NorthEast", "NorthWest", "South", "SouthEast", "SouthWest", "Center", "East", "West"}, false))
}

func validateGravityPost() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"North", "NorthEast", "NorthWest", "South", "SouthEast", "SouthWest", "Center", "East", "West"}, false))
}

---
layout: "akamai"
page_title: "Transform images"
description: |-
  Learn about image transformation arguments for Image and Video Manager.
---

# Transform Images (Beta)

Add transformations to your Image and Video manager image policy to automate processes like cropping, rotating, and resizing, or to apply visual effects to your images. Multiple transformations can be nested to achieve the desired result.

Set up transformations in the `akamai_imaging_policy_image` data source.

This guide provides the arguments to use when configuring image transformations.

## Available transformations

The sections that follow cover the available image transformations and their supporting arguments.

* To manipulate the size, shape, and orientation of your images, use: [`append`](#append), [`aspect_crop`](#aspect-crop), [`crop`](#crop), [`face_crop`](#face-crop), [`feature_crop`](#feature-crop), [`fit_and_fill`](#fit-and-fill), [`mirror`](#mirror), [`region_of_interest_crop`](#region-of-interest-crop), [`relative_crop`](#relative-crop), [`resize`](#resize), [`rotate`](#rotate), [`scale`](#scale), [`shear`](#shear), [`trim`](#trim)

* To apply visual effects to your images, use: [`background_color`](#background-color), [`blur`](#blur), [`chroma_key`](#chroma-key), [`composite`](#composite), [`contrast`](#contrast), [`goop`](#goop), [`grayscale`](#grayscale), [`hsl`](#hsl), [`hsv`](#hsv), [`max_colors`](#max-colors), [`mono_hue`](#mono-hue), [`opacity`](#opacity), [`remove_color`](#remove-color), [`unsharp_mask`](#unsharp-mask)

* To group together a sequence of transformations and represent it as a single transformation, use: [`compound`](#compound)

* To apply artistic transformations on a per-image basis without having to create multiple policies, use: [`im_query`](#imquery)

* To conditionalize transformations based on the dimensions or orientation of original images, use: [`if_dimension`](#if-dimension), [`if_orientation`](#if-orientation)

A subset of transformations can be applied after image and quality settings are applied (post-processing). This subset includes  [`background_color`](#background-color), [`blur`](#blur), [`chroma_key`](#chroma-key), [`composite`](#composite), [`compound`](#compound), [`contrast`](#contrast), [`goop`](#goop), [`grayscale`](#grayscale), [`hsl`](#hsl), [`hsv`](#hsv), [`if_dimension`](#if-dimension), [`if_orientation`](#if-orientation), [`max_colors`](#max-colors), [`mirror`](#mirror), [`mono_hue`](#mono-hue), [`opacity`](#opacity), [`remove_color`](#remove-color), [`unsharp_mask`](#unsharp-mask).

For more information about the available image transformations, see [Add image transformations and conditions](https://techdocs.akamai.com/ivm/docs/add-img-transformations).

## Variables

Many Image and Video Manager arguments let you specify a variable object instead of a string, number, or boolean value.

When using variables, you define the variable name in an argument that ends in `_var`. For example, if you want to have a variable for the gravity setting in a transformation, you’d use the `gravity_var` argument, not the `gravity` one.

## Append

This transformation supports these arguments:

* `append `- (Optional) Places a specified `image` beside the source image. The transformation places the `image` on a major dimension, then aligns it on the minor dimension. Transparent pixels fill any area not covered by either image.
    * `image` - (Required) The image type, which is one of `box_image`, `circle_image`, `text_image`, or `url_image`. See [Image types](#image-types).
    * `gravity` - (Optional) The placement of the `image` relative to the source image. The available values include eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var` instead.
    * `gravity_priority` - (Optional) Determines the exact placement of the `image` when `gravity` is `Center` or a diagonal direction. The value is either `horizontal` or `vertical`. Use `horizontal` to append an `image` east or west of the source image. This aligns the `image` on the vertical gravity component, placing `Center` gravity east. Use `vertical` to append an `image` north or south of the source image. This aligns the `image` on the horizontal gravity component, placing `Center` gravity south. If setting a variable for this argument, use `gravity_priority_var` instead.
    * `preserve_minor_dimension` - (Optional) Defines whether to preserve the source image's minor dimension. This argument is `false` by default. The minor dimension is the dimension opposite to the dimension that the appending `image` is placed. If setting a variable for this argument, use `preserve_minor_dimension_var` instead.

## Aspect crop

This transformation supports these arguments:

* `aspect_crop` - (Optional) Changes the height or width of an image to an aspect ratio of your choosing. The transformation does this by either cropping or expanding the area.
    * `allow_expansion` - (Optional) Increases the size of the image canvas to achieve the requested aspect ratio instead of cropping the image. If setting a variable for this argument, use `allow_expansion_var` instead.
    * `height` - (Optional) The height of the aspect ratio to crop. If setting a variable for this argument, use `height_var` instead.
    * `width` - (Optional) The width of the aspect ratio to crop. If setting a variable for this argument, use `width_var` instead.
    * `x_position` - (Optional) The horizontal portion of the image you want to keep when applying aspect ratio cropping. If using `allow_expansion`, this setting defines the horizontal position of the image on the new expanded image canvas. If setting a variable for this argument, use `x_position_var` instead.
    * `y_position` - (Optional) The horizontal portion of the image you want to keep when applying aspect ratio cropping. If using `allow_expansion`, this setting defines the horizontal position of the image on the new expanded image canvas. If setting a variable for this argument, use `y_position_var` instead.

## Background color

This transformation supports this argument:

* `background_color` - (Optional) Places a transparent image on a set background color.
    * `color` - (Optional) The hexadecimal CSS color value for the background. If setting a variable for this argument, use `color_var` instead.

## Blur

This transformation supports these arguments:

* `blur` - (Optional) Applies a Gaussian blur.
    * `sigma` - (Optional) The number of pixels by which to scatter the original pixel to create the blur effect. Resulting images may be larger than the original as some pixels may scatter outside the image's original dimensions. If setting a variable for this argument, use `sigma_var` instead.

## Chroma Key

This transformation supports these arguments:

* `chroma_key` - (Optional) Changes any color in an image within the specified volume of the HSL colorspace to transparent or semitransparent. This transformation applies a 'green screen' technique commonly used to isolate and remove background colors.
    * `hue` - (Optional) The hue to remove. Enter the degrees of rotation (between 0 and 360) around the color wheel. By default `chroma_key` removes a green hue which is at 120 degrees on the color wheel. If setting a variable for this argument, use `hue_var` instead.
    * `hue_feather` - (Optional) How much additional hue to make semi-transparent beyond the `hue_tolerance`. The default is `0.083` which applies semi-transparency to hues 30 degrees around the Hue Tolerance. If setting a variable for this argument, use `hue_feather_var` instead.
    * `hue_tolerance` - (Optional) How close a color's hue needs to be to the selected hue for it to be changed to fully transparent. If you enter the maximum value of `1.0` the entire image is made transparent. The default is approximately `0.083` or 8.3% of the color wheel. This value corresponds to 30 degrees around the specified hue. If setting a variable for this argument, use `hue_tolerance_var` instead.
    * `lightness_feather` - (Optional) How much additional lightness to make semi-transparent beyond the `lightness_tolerance`. The default value of 0.1 corresponds to 10% away from the tolerated lightness towards full black or full white. If setting a variable for this argument, use `lightness_feather_var`  instead.
    * `lightness_tolerance` - (Optional) How much of the lightest part and darkest part of a color to preserve. You can space this value out from the middle (for example, 0.5 lightness or full color) to help preserve the splash lighting impact in the image. You can define how close the color needs to be to the full color to remove it from your image. The default value of `0.75` means that a color must be within 75% of the full color to full white or full black for full removal.  If setting a variable for this argument, use `lightness_tolerance_var`  instead.
    * `saturation_feather` - (Optional) How much additional saturation to make semi-transparent beyond the saturation tolerance. The default is `0.1` which applies semi-transparency to hues 10% below the `saturation_tolerance`. If setting a variable for this argument, use `saturation_feather_var`  instead.
    * `saturation_tolerance` - (Optional) How close a color's saturation needs to be to full saturation for it to be changed to fully transparent. For example, you can define how green the color needs to be to remove it from your image. The default value of `0.75` means that a color must be within 75% of full saturation in order to be made fully transparent. If setting a variable for this argument, use `saturation_tolerance_var`  instead.

## Composite

This transformation supports these arguments:

* `composite` - (Optional) Applies another image to the source image, either as an overlay or an underlay. The image underneath is visible in areas that are beyond the edges of the top image or that are less than 100% opaque. A common use of the overlay composite is to add a watermark.
   * `image` - (Required) The image type, which is one of `box_image`, `circle_image`, `text_image`, or `url_image`. See [Image types](#image-types) for supported subarguments.
    * `gravity` - (Optional) The placement of the `image` relative to the source image. The available values include eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var` instead.
    * `placement` - (Optional) The placement of the applied image either on top of or underneath the base image. Watermarks are usually applied on top. Backgrounds are usually applied underneath. If setting a variable for this argument, use `placement_var`  instead.
    * `scale` - (Optional) A multiplier to resize the applied image relative to the source image while preserving the aspect ratio (1 by default). Set the `scale_dimension` to calculate the `scale` from the source image's width or height. If setting a variable for this argument, use `scale_var` instead.
`scale_dimension` - (Optional) The dimension, either `width` or `height`, of the source image to scale. If setting a variable for this argument, use `scale_dimension_var`  instead.
    * `x_position` - (Optional) The x-axis position of the image to apply. If setting a variable for this argument, use `x_position_var`  instead.
    * `y_position` - (Optional) The y-axis position of the image to apply.  If setting a variable for this argument, use `y_position_var`  instead.

## Compound

This transformation supports this argument:

* `compound` - (Optional) Groups together transformations into an ordered set. Used to represent a sequence of transformations as a single transformation.

## Contrast

This transformation supports these arguments:

* `contrast`- (Optional) Adjusts contrast and brightness of an image.
    * `brightness` - (Optional) The brightness of the image. Positive values increase brightness and negative values decrease brightness. A value of  `1` produces a white image. A value of  `-1` produces a black image. The default value is `0`, which leaves the image unchanged. The acceptable value range is `-1.0` to `1.0`. Values outside of the acceptable range clamp to this range. If setting a variable for this argument, use `brightness_var`  instead.
    * `contrast` - (Optional) The contrast of the image. Expressed as a range from `-1` to `1`. Positive values increase contrast, negative values decrease contrast, while `0` leaves the image unchanged. Values outside of the `-1` to `1` range clamp to this range. If setting a variable for this argument, use `contrast_var` instead.

## Crop

This transformation supports these arguments:

* `crop` - (Optional) Crops an image.
    * `allow_expansion` - (Optional) If cropping an area outside of the existing canvas, expands the image canvas. If setting a variable for this argument, use `allow_expansion_var`  instead.
    * `gravity` - (Optional) The placement of the `crop` relative to the source image. The available values include eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var` instead.
    * `height` - (Optional) The number of pixels to crop along the y-axis. If setting a variable for this argument, use `height_var` instead.
    * `width` - (Optional) The number of pixels to crop along the x-axis. If setting a variable for this argument, use `width_var`  instead.
    * `x_position` - (Optional) The x-axis position to crop from. If setting a variable for this argument, use `x_position_var`  instead.
    * `y_position` - (Optional) The y-axis position to crop from. If setting a variable for this argument, use `y_position_var` instead.


## Face Crop

This transformation supports these arguments:

* `face_crop` - (Optional) Applies a method to detect faces in the source image and applies the rectangular crop on either the `biggest` face or `all` of the faces detected. Image and Video Manager tries to preserve faces in the image instead of using specified crop coordinates.
    * `algorithm` - (Optional) Specifies the type of algorithm used to detect faces in the image: `cascade` (default) for the cascade classifier algorithm, or `dnn` for the deep neural network algorithm. If setting a variable for this argument, use `algorithm_var`  instead.
    * `confidence` - (Optional) With `algorithm` set to `dnn`, the minimum confidence needed to detect faces in the image. Values range from `0` to `1` for increased confidence, and possibly fewer faces detected. If setting a variable for this argument, use `confidence_var`  instead.
    * `fail_gravity` - (Optional) Controls placement of the crop if no faces are detected in the image. Directions are relative to the edges of the image being transformed. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `fail_gravity_var`  instead.
    * `focus` - (Optional) The focus of the crop rectangle, which is either around `biggestFace` or `allFaces`. This is `all` by default. If setting a variable for this argument, use `focus_var` instead.
    * `gravity` - (Optional) The placement of the crop relative to the faces, plus padding. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var`  instead.
    `height` - (Optional) The height of the output image in pixels relative to the specified `style` value. If setting a variable for this argument, use `height_var` instead.
    * `padding` - (Optional) The padding ratio based on the dimensions of the biggest face detected, This is `0.5` by default. Larger values increase padding. If setting a variable for this argument, use `padding_var` instead.
    * `style` - (Optional) How to crop, or scale a crop area, for the faces detected in the source image. This is `zoom` by default. The output image is resized to the specified `width` and `height` values. A value of `crop` places a raw crop around the faces relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and faces as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the faces, relative to the specified `width` and `height` values. If setting a variable for this argument, use `style_var`  instead.
    * `width` - (Optional) The width of the output image in pixels relative to the specified `style` value. If setting a variable for this argument, use `width_var` instead.

## Feature Crop

This transformation supports these arguments:

* `feature_crop` - (Optional) Identifies prominent features of the source image, then crops around as many of these as possible relative to the specified `width` and `height` values.
    * `fail_gravity` - (Optional) The placement of the crop if Image and Video Manager does not detect any features in the image. Directions are relative to the edges of the image being transformed. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `fail_gravity_var` instead.
    * `feature_radius` - (Optional) The size in pixels of the important features to search for. If identified, two features never appear closer together than this value which is `8.0` by default.  If setting a variable for this argument, use `feature_radius_var` instead.
    * `gravity` - (Optional) The placement of the crop relative to the region of interest plus padding. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var` instead.
    * `height` - (Optional) The height in pixels of the output image relative to the specified `style` value. If setting a variable for this argument, use `height_var` instead.
    * `max_features` - (Optional) The maximum number of features to identify as important features. This maximum number of features is `32` by default. The strongest features are always chosen. If setting a variable for this argument, use `max_features_var` instead.
    * `min_feature_quality` - (Optional) The minimum quality level of the feature identified. To be considered important, the feature needs to surpass this value. Image and Video Manager measures quality on a scale from `0` (lowest) to `1` (highest). This is `0.1` by default. If setting a variable for this argument, use `min_feature_quality_var` instead.
    * `padding` - (Optional) Adds space around the region of interest. The amount of padding added is directly related to the size of the bounding box of the selected features. Specifically, the region of interest is expanded in all directions by the largest dimension of the bounding box of the selected features multiplied by this value. If setting a variable for this argument, use `padding_var` instead.
    * `style` - (Optional) Specifies how to crop or scale a crop area for the features identified in the source image. This is `fill` by default. The output image resizes to the specified `width` and `height` values. A value of `crop` performs a raw crop around the features relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and features as possible, relative to the specified `width` and `height` values. A value of `zoom` scales the crop area as small as possible to fit the features, relative to the specified `width` and `height` values. If setting a variable for this argument, use `style_var` instead.
    * `width` - (Optional) The width in pixels of the output image relative to the specified `style` value. If setting a variable for this argument, use `width_var` instead.

## Fit and Fill

This transformation supports these arguments:

* `fit_and_fill` - (Optional) Resizes an image to fit within a specific size box and then uses a fill of that same image to cover any transparent space at the edges. By default, the fill image has a Blur transformation with a sigma value of `8` applied. You can use the `fill_transformation` argument to customize the transformation.
    * `fill_transformation` - Used to customize the sigma value for the fill image.
    * `height` - (Optional) The height value of the resized image. If setting a variable for this argument, use `height_var` instead.
    * `width` - (Optional) The width value of the resized image. If setting a variable for this argument, use `width_var` instead.

## Goop

This transformation supports these arguments:

* `goop` - (Optional) Distorts an image by randomly repositioning a set of control points along a specified grid. The transformed image appears _goopy_. Adjust the density of the grid and the degree of randomity. You can use this transformation to create watermarks for use in security.
   * `chaos` - (Optional) The greatest distance control points may move from their original position. A value of `1.0` shifts control points over as far as the next one in the original grid. A value of `0.0` leaves the image unchanged. Values under `0.5` work better for subtle distortions, otherwise control points may pass each other and cause a twisting effect. If setting a variable for this argument, use `chaos_var` instead.
   * `density` - (Optional) The density of control points used to distort the image. The largest dimension of the input image is divided up to fit this number of control points. A grid of points is extended on the smaller dimension such that each row and column of control points is equidistant from each adjacent row or column. This parameter strongly affects transformation performance. Be careful choosing values above the default if you expect to transform medium to large sized images. If setting a variable for this argument, use `density_var` instead.
   * `power` - (Optional) By default, the distortion algorithm relies on inverse squares to calculate distance but this allows you to change the exponent. You shouldn't need to vary the default value of `2.0`. If setting a variable for this argument, use `power_var` instead.
   * `seed` - (Optional) Your own seed value as an alternative to the default, which is subject to variability. This allows for reproducible and deterministic distortions. If all parameters are kept equal and a constant seed is used, `goop` distorts an input image consistently over many transformations. By default, this value is set to the current epoch time measured in milliseconds, which provides inconsistent transformation output. If setting a variable for this argument, use `seed_var` instead.

## Grayscale

This transformation supports these arguments:

* `grayscale` - (Optional) Restricts image color to shades of gray only.
    * `type` - (Optional) The algorithm used to transform colors to grays. The available types are `Brightness`, `Lightness`, `Rec601`, or the default `Rec709`. If setting a variable for this argument, use `type_var` instead.

## HSL (Hue/Saturation/Lightness)

This transformation supports these arguments:

* `hsl` - (Optional) The hue, saturation, and lightness (HSL) of an image. Hue is the number of degrees of rotation around the color wheel. Saturation is a multiplier to increase or decrease color saturation. Lightness is a multiplier to increase or decrease the lightness of an image. Other transformations can also affect color, such as [`grayscale`](#grayscale) and [`max_colors`](#max-colors). If you're using more than one of these transformations, consider the order of application for the desired results.
    * `hue` - (Optional) The number of degrees to rotate colors around the color wheel. The default is `0`. If setting a variable for this argument, use `hue_var` instead.
    * `lightness` - (Optional) A multiplier to adjust the lightness of colors in the image. Note that lightness is distinct from brightness. For example, reducing the lightness of a light green might give you a lime green whereas reducing the brightness of a light green might give you a darker shade of the same green. Values less than `1.0` decrease the lightness of colors in the image. Values greater than `1.0` increase the lightness of colors in the image. If setting a variable for this argument, use `lightness_var` instead.
    * `saturation` - (Optional) A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image. If setting a variable for this argument, use `saturation_var` instead.

## HSV (Hue/Saturation/Value)

This transformation supports these arguments:

* `hsv` - (Optional) The hue, saturation, and value (HSV) of an image. `hsv` is like `hsl` except `lightness` is replaced with `value`. To illustrate the difference, if you reduce the `lightness` of a light green (almost white), the color becomes a vibrant green, but if you reduce the `value`, the color becomes darker, close to gray. This happens because the original image color is very close to white.
    * `hue` - (Optional) The number of degrees to rotate colors around the color wheel. This is `0.0` by default. If setting a variable for this argument, use `hue_var` instead.
    * `saturation` - (Optional) A multiplier to adjust the saturation of colors in the image. Values less than `1.0` decrease saturation and values greater than `1.0` increase the saturation. A value of `0.0` removes all color from the image. If setting a variable for this argument, use `saturation_var` instead.
    * `value` - (Optional) A multiplier to adjust the lightness or darkness of the image's base color. Values less than `1.0` decrease the base colors in the image, making them appear darker. Values greater than `1.0` increase the base colors in the image, making them appear lighter. If setting a variable for this argument, use `value_var` instead.

## If Dimension

This transformation supports these arguments:

* `if_dimension` - (Optional) Chooses a transformation depending on the dimensions of the source image.
    * `default` - (Optional) A no-op transformation, by default.
    * `dimension` - (Optional) The dimension to use to select the transformation, either `height`, `width`, or `both`. If setting a variable for this argument, use `dimension_var` instead.
    * `equal` - (Optional) The transformation is performed only if the source image's dimension is equal to the value listed.
    * `greater_than` - (Optional) The transformation is performed if the source image's dimension is greater than the value listed.
    * `less_than` - (Optional) The transformation is  performed if the source image's dimension is less than the value listed.
    * `value` - (Optional) The value against which to compare the source image dimension. For example, if the image dimension is less than the value, the `less_than` transformation is applied. If setting a variable for this argument, use `value_var` instead.

# If Orientation

This transformation supports these arguments:

* `if_orientation` - (Optional) Chooses a transformation depending on the orientation of the source image.
    * `default` - (Optional) A no-op transformation, by default.
    * `landscape` - (Optional)  The transformation is performed if the source image uses landscape orientation.
    * `portrait` - (Optional) The transformation is performed if the source image uses portrait orientation.
    * `square` - (Optional) The transformation is performed if the source image uses a square orientation.

## IMQuery

This transformation supports these arguments:

* `im_query` - (Optional) Applies artistic transformations on a per-image basis by specifying transformations with a query string appended to the image URL.
    * `allowed_transformations` - (Required) Specifies the transformations you can apply with the query string parameter. `append`, `blur`, and `crop` are supported.
    * `query_var` - (Required) The variable to use.

## Max Colors

This transformation supports these arguments:

* `max_colors` - (Optional) Sets the maximum number of colors in the image’s palette. Reducing the number of colors in an image can help to reduce file size.
    * `colors` - (Optional) The value representing the maximum number of colors to use in the source image. If setting a variable for this argument, use `colors_var` instead.

## Mirror

This transformation supports these arguments:

* `mirror` - (Optional) Flips an image horizontally, vertically, or both.
    * `horizontal` - (Optional) Flips the image horizontally. If setting a variable for this argument, use `horizontal_var` instead.
    * `vertical` - (Optional) Flips the image vertically. If setting a variable for this argument, use `vertical_var` instead.

## Mono Hue

This transformation supports these arguments:

* `mono_hue` - (Optional) Allows you to set all hues in an image to a single specified hue of your choosing. Mono Hue maintains the original color’s lightness and saturation but sets the hue to that of the specified value. This has the effect of making the image various shades of the specified hue.
    * `hue` - (Optional) The hue which is indicated by the degree of rotation (between 0 and 360) around the color wheel. By default, Mono Hue applies a red hue, which is `0.0` on the color wheel. If setting a variable for this argument, use `hue_var` instead.

## Opacity

This transformation supports this argument:

* `opacity` - (Optional) The level of transparency of an image represented as an alpha value on a scale of `0` to `1`. An image with no transparency is opaque. Values below `1` decrease opacity with `0` being completely transparent. For images with some transparency, values above `1` increase the opacity of the transparent portions. If setting a variable for this argument, use `opacity_var` instead.

## Region of Interest Crop

This transformation supports these arguments:

* `region_of_interest_crop` - (Optional) Crops around a specified area of interest (region) relative to the specified `width` and `height` values.
    * `region_of_interest` - (Required) The bounding shape of the important features to search for, which is one of `circle_shape`, `point_shape`, `polygon_shape`, `rectangle_shape`, or `union_shape`. See [Shape types](#shape-types).
    * `gravity` - (Optional) The placement of the crop area relative to the specified area of interest. The available values represent the eight cardinal directions (`North`, `South`, `East`, `West`, `NorthEast`, `NorthWest`, `SouthEast`, `SouthWest`) and `Center` by default. If setting a variable for this argument, use `gravity_var` instead.
    * `height` - (Optional) The height in pixels of the output image relative to the specified `style` value. If setting a variable for this argument, use `height_var` instead.
    * `style` - (Optional) Specifies how to crop, or scale a crop area, for the specified area of interest in the source image. The output image resizes to the specified `width` and `height` values. A value of `crop` places raw crop around the point of interest, relative to the specified `gravity` value.  A value of `fill` scales the crop area to include as much of the image and point of interest as possible, relative to the specified `width` and `height` values. A value of `zoom` (the default) scales the crop area as small as possible to fit the point of interest, relative to the specified `width` and `height` values. If setting a variable for this argument, use `style_var` instead.
    * `width` - (Optional) The width in pixels of the output image relative to the specified `style` value. If setting a variable for this argument, use `width_var` instead.

## Relative Crop

This transformation supports these arguments:

* `relative_crop` - (Optional) Shrinks or expands an image relative to the image's specified dimensions. Image and Video Manager fills the expanded areas with transparency. Positive values shrink the side, while negative values expand it.
    * `east` - (Optional) The number of pixels to shrink or expand the right side of the image. If setting a variable for this argument, use `east_var` instead.
    * `north` - (Optional) The number of pixels to shrink or expand the top side of the image. If setting a variable for this argument, use `north_var` instead.
    * `south` - (Optional) The number of pixels to shrink or expand the bottom side of the image. If setting a variable for this argument, use `south_var` instead.
    * `west` - (Optional) The number of pixels to shrink or expand the left side of the image. If setting a variable for this argument, use `west_var` instead.

## Remove Color

This transformation supports these arguments:

* `remove_color` - (Optional) Removes a specified color from an image and replaces it with transparent pixels. This transformation is ideal for removing solid background colors from product images photographed on clean, consistent backgrounds without any shadows.
    * `color` - (Optional) The hexadecimal CSS color value to remove. If setting a variable for this argument, use `color_var`.
    * `feather` - (Optional) Used to minimize any hard edges and to make the color removal more gradual in appearance. This option allows you to extend the color removal beyond the specified `tolerance`. The pixels in this extended tolerance become semi-transparent which creates a softer edge. The first time there’s a real-time request for an image, this option may result in a slow transformation time, but subsequent requests aren't impacted as the transformed image is served directly from the cache. If setting a variable for this argument, use `feather_var` instead.
    * `tolerance` - (Optional) Defines how close the color needs to be to the selected color before it's changed to fully transparent. Set to `0.0` to remove only the exact color specified. If setting a variable for this argument, use `tolerance_var` instead.

## Resize

This transformation supports these arguments:

* `resize` - (Optional) Resizes an image to a particular, absolute dimension. If you don't enter `width` or `height`, this transformation uses the `fit` aspect preservation mode, which selects a value for the missing dimension that preserves the image's aspect.
    * `aspect` - (Optional) Preserves the aspect ratio. Select `fit` to make the image fit entirely within the selected width and height. When using `fit`, the resulting image has the largest possible size for the specified dimensions. Select `fill` to size the image so it both completely fills the dimensions and has the smallest possible file size. Otherwise `ignore` changes the original aspect ratio to fit within an arbitrarily shaped rectangle.  If setting a variable for this argument, use `aspect_var` instead.
    * `height` - (Optional) The height to which to resize the source image. Must be set if width is not specified. If setting a variable for this argument, use `height_var` instead.
    * `type` - (Optional) The type of constraints on the image resize. Select `normal` to resize in all cases, either increasing or decreasing the dimensions. Select `downsize` to ignore this transformation if the result would be larger than the original. Select `upsize` to ignore this transformation if the result would be smaller. If setting a variable for this argument, use `type_var` instead.
    * `width` - (Optional) The width to which to resize the source image. Must be set if height is not specified. If setting a variable for this argument, use `width_var` instead.

## Rotate

This transformation supports these arguments:

* `rotate` - (Optional) Rotates the image around its center by indicating the degrees of rotation.
    * `degrees` - (Optional) The value by which to rotate the image. Positive values rotate clockwise, while negative values rotate counter-clockwise. If setting a variable for this argument, use `degrees_var` instead.

## Scale

This transformation supports these arguments:

* `scale` - (Optional) Changes the derivative image's dimensions relative to the original image’s.
    * `height` - (Optional) The scaling factor for the input height that is used to determine the output height of the image. Values between `0` and `1` decrease image height. Values greater than 1 increase the image height. A value of `1` leaves the height unchanged. Image dimensions need to be non-zero positive numbers. If setting a variable for this argument, use `height_var` instead.
    * `width` - (Optional) The scaling factor for the input width that is used to determine the output width of the image.  Values between `0` and `1` decrease image width. Values greater than 1 increase the image width. A value of `1` leaves the width unchanged. Image dimensions need to be non-zero positive numbers. If setting a variable for this argument, use `width_var` instead.

## Shear

This transformation supports these arguments:

* `shear` -  (Optional) Slants an image into a parallelogram as a percent of the starting dimension in decimal format. You need to specify at least one axis property. Transparent pixels fill empty areas around the sheared image as needed, so it's often useful to use a `background_color` transformation for these areas.
    * `x_shear` - (Optional) The amount to shear along the x-axis, measured in multiples of the image's width. Must be set if `y_shear` isn't specified. If setting a variable for this argument, use `x_shear_var` instead.
    * `y_shear` - (Optional) The amount to shear along the y-axis, measured in multiples of the image's height. Must be set if `x_shear` isn't specified. If setting a variable for this argument, use `y_shear_var` instead.

## Trim

This transformation supports these arguments:

* `trim` - (Optional) Automatically crops the background uniformly from the edges of an image.
    * `fuzz` - (Optional) The fuzz tolerance of the trim expressed as a value between `0` and `1`. This determines the acceptable amount of background variation before trimming stops. If setting a variable for this argument, use `fuzz_var` instead.
    * `padding` - (Optional) The amount of padding in pixels to add to the trimmed image. If setting a variable for this argument, use `padding_var` instead.


## Unsharp Mask

This transformation supports these arguments:

* `unsharp_mask` - (Optional) Emphasizes edges and details in source images without distorting the colors. Although this effect is often referred to as _sharpening_ an image, it actually creates a blurred, inverted copy of the image known as an unsharp mask. Image and Video Manager combines the unsharp mask with the source image to create an image that is perceived by the human eye as clearer.
    * `gain` - (Optional) How much emphasis the filter applies to details. Higher values increase apparent sharpness of details. If setting a variable for this argument, use `gain_var` instead.
    * `sigma` - (Optional) The standard deviation of the Gaussian distribution used in the unsharp mask, measured in pixels. This is `1.0` by default. High values emphasize large details and low values emphasize small details. If setting a variable for this argument, use `sigma_var` instead.
    * `threshold` - (Optional) The minimum change required to include a detail in the filter. Higher values discard more changes. If setting a variable for this argument, use `threshold_var` instead.

## Image types

You can, and in some cases may be required to, specify these image types for certain transformations:
* `box_image` - (Optional) A rectangular box with a specified fill color and applied transformation.
    * `color` - (Optional) The fill color of the box, not the edge of the box. The API supports hexadecimal representation and CSS hexadecimal color values. If setting a variable for this argument, use `color_var` instead.
    * `height` - (Optional) The height of the box in pixels. If setting a variable for this argument, use `height_var` instead.
    * `transformation` - (Optional) The transformation to apply to the box.
    * `width` - (Optional) The width of the box in pixels. If setting a variable for this argument, use `width_var` instead.
* `circle_image` - (Optional) A rectangular box, with a specified color and applied transformation.
    * `color` - (Optional) The fill color of the circle. The API supports hexadecimal representation and CSS hexadecimal color values.  If setting a variable for this argument, use `color_var`  instead.
    * `diameter` - (Optional) The diameter of the circle. The diameter will be the width and the height of the image in pixels. If setting a variable for this argument, use `diameter_var`  instead.
    * `transformation` - (Optional) The transformation to apply to the circle.
    * `width` - (Optional) The width of the box in pixels. If setting a variable for this argument, use `width_var` instead.
* `text_image` - (Optional) A snippet of text. Defines font family and size, fill color, as well as outline stroke width and color.
    * `fill` - (Optional) The main fill color of the text. If setting a variable for this argument, use `fill_var` instead.
    * `size` - (Optional) The size in pixels to render the text. If setting a variable for this argument, use `size_var` instead.
    * `stroke` - (Optional) The color for the outline of the text. If setting a variable for this argument, use `stroke_var` instead.
    * `stroke_size` - (Optional) The thickness in points for the outline of the text. If setting a variable for this argument, use `stroke_size_var` instead.
    * `text` - (Optional) The line of text to render. If setting a variable for this argument, use `text_var` instead.
    * `transformation` - (Optional) The transformation to apply to the text.
    * `typeface` - (Optional) The font family to apply to the text image. This may be a URL to a TrueType or WOFF (v1) typeface, or a string that refers to one of the standard built-in browser fonts. If setting a variable for this argument, use `typeface_var` instead.
* `url_image` - (Optional) An image loaded from a URL.
    * `transformation` - (Optional)  The transformation to apply to the image.
    * `url` - (Optional) The URL of the image. If setting a variable for this argument, use `url_var` instead.

### Example

Here’s an example of a text image within an append transformation.

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

## Shape types

You can specify these shape types for certain transformations:
* `circle_shape` - (Optional) Defines a circle with a specified `radius` from its `center` point.
    * `center` - (Required) Defines coordinates for a single point, to help define polygons and rectangles. Each point may be an object with `x`and `y` members, or a two-element array.
    * `radius` - (Optional) The radius of the circle measured in pixels. If setting a variable for this argument, use `radius_var`  instead.
* `point_shape` - (Optional) Defines coordinates for a single point to help define polygons and rectangles. Each point may be an object with x and y members, or a two-element array.
    * `x` - (Optional) The horizontal position of the point, measured in pixels. If setting a variable for this argument, use `x_var` instead.
    * `y` - (Optional) The vertical position of the point, measured in pixels. If setting a variable for this argument, use `y_var` instead.
* `polygon_shape` - (Optional) Defines a polygon from a series of connected points.
    * `points` - (Required) A series of `point_shape` objects. The last and first points connect to close the shape automatically.
* `rectangle_shape` - (Optional) Defines a rectangle's `width` and `height` relative to an `anchor` point at the top left corner.
    * `anchor` - (Required) The anchor point for the rectangle.
    * `height` - (Optional) Extends the rectangle down from the `anchor` point. If setting a variable for this argument, use `height_var` instead.
    * `width` - (Optional) Extends the rectangle right from the `anchor` point. If setting a variable for this argument, use `width_var` instead.
* `union_shape` - (Optional) Identifies a shape based on a combination of other shapes. You can use an object to represent a union or an array of shapes that describe it.
    * `shapes` - (Required) The set of shapes to combine to form the union.

### Example

Here’s an example of a rectangle shape within a region of interest crop.

```hcl
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
        style  = "fill"
        width = 7
      }
    }
```

package ruleformats

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	schemasRegistry.register(RuleFormat{
		version:          "rules_v2024_10_21",
		behaviorsSchemas: getBehaviorsSchemaV20241021(),
		criteriaSchemas:  getCriteriaSchemaV20241021(),
		typeMappings:     map[string]interface{}{"adScalerCircuitBreaker.returnErrorResponseCodeBased.408": 408, "adScalerCircuitBreaker.returnErrorResponseCodeBased.500": 500, "adScalerCircuitBreaker.returnErrorResponseCodeBased.502": 502, "adScalerCircuitBreaker.returnErrorResponseCodeBased.504": 504},
		nameMappings:     map[string]string{"allowFcmParentOverride": "allowFCMParentOverride", "allowHttpsCacheKeySharing": "allowHTTPSCacheKeySharing", "allowHttpsDowngrade": "allowHTTPSDowngrade", "allowHttpsUpgrade": "allowHTTPSUpgrade", "businessCategory": "BUSINESS_CATEGORY", "c": "C", "canBeCa": "canBeCA", "cn": "CN", "conditionalHttpStatus": "conditionalHTTPStatus", "contentCharacteristicsAmd": "contentCharacteristicsAMD", "contentCharacteristicsDd": "contentCharacteristicsDD", "countryOfCitizenship": "COUNTRY_OF_CITIZENSHIP", "countryOfResidence": "COUNTRY_OF_RESIDENCE", "dateOfBirth": "DATE_OF_BIRTH", "dc": "DC", "dcpAuthHmacTransformation": "dcpAuthHMACTransformation", "detectSmartDnsProxy": "detectSmartDNSProxy", "detectSmartDnsProxyAction": "detectSmartDNSProxyAction", "detectSmartDnsProxyRedirecturl": "detectSmartDNSProxyRedirecturl", "dnQualifier": "DN_QUALIFIER", "emailAddress": "EmailAddress", "enableCmcdSegmentPrefetch": "enableCMCDSegmentPrefetch", "enableEs256": "enableES256", "enableIpAvoidance": "enableIPAvoidance", "enableIpProtection": "enableIPProtection", "enableIpRedirectOnDeny": "enableIPRedirectOnDeny", "enableRs256": "enableRS256", "enableTokenInUri": "enableTokenInURI", "g2OToken": "g2oToken", "g2Oheader": "g2oheader", "gender": "GENDER", "generation": "GENERATION", "givenname": "GIVENNAME", "i18NCharset": "i18nCharset", "i18NStatus": "i18nStatus", "initials": "INITIALS", "isCertificateSniOnly": "isCertificateSNIOnly", "issuerRDNs.description": "DESCRIPTION", "issuerRDNs.name": "NAME", "issuerRdns": "issuerRDNs", "jurisdictionC": "JURISDICTION_C", "jurisdictionL": "JURISDICTION_L", "jurisdictionSt": "JURISDICTION_ST", "l": "L", "logEdgeIp": "logEdgeIP", "nameAtBirth": "NAME_AT_BIRTH", "o": "O", "organizationIdentifier": "ORGANIZATION_IDENTIFIER", "originSettings": "origin_settings", "ou": "OU", "overrideIpAddresses": "overrideIPAddresses", "placeOfBirth": "PLACE_OF_BIRTH", "postalAddress": "POSTAL_ADDRESS", "postalCode": "POSTAL_CODE", "pseudonym": "PSEUDONYM", "role": "ROLE", "segmentDurationDash": "segmentDurationDASH", "segmentDurationDashCustom": "segmentDurationDASHCustom", "segmentDurationHds": "segmentDurationHDS", "segmentDurationHdsCustom": "segmentDurationHDSCustom", "segmentDurationHls": "segmentDurationHLS", "segmentDurationHlsCustom": "segmentDurationHLSCustom", "segmentSizeDash": "segmentSizeDASH", "segmentSizeHds": "segmentSizeHDS", "segmentSizeHls": "segmentSizeHLS", "serialnumber": "SERIALNUMBER", "sf3COriginHost": "sf3cOriginHost", "sf3COriginHostHeader": "sf3cOriginHostHeader", "smartDnsProxy": "smartDNSProxy", "st": "ST", "standardTlsMigration": "standardTLSMigration", "standardTlsMigrationOverride": "standardTLSMigrationOverride", "street": "STREET", "subjectCn": "subjectCN", "subjectRDNs.description": "DESCRIPTION", "subjectRDNs.name": "NAME", "subjectRdns": "subjectRDNs", "surname": "SURNAME", "t": "T", "telephoneNumber": "TELEPHONE_NUMBER", "titleAicMobile": "title_aic_mobile", "titleAicNonmobile": "title_aic_nonmobile", "tokenAuthDashTitle": "tokenAuthDASHTitle", "tokenAuthHlsTitle": "tokenAuthHLSTitle", "uid": "UID", "uniqueIdentifier": "UNIQUE_IDENTIFIER", "unstructuredAddress": "UnstructuredAddress", "unstructuredName": "UnstructuredName"},
		shouldFlatten:    []string{"apiPrioritization.cloudletPolicy", "apiPrioritization.throttledCpCode", "apiPrioritization.throttledCpCode.cpCodeLimits", "apiPrioritization.netStorage", "applicationLoadBalancer.cloudletPolicy", "applicationLoadBalancer.allDownNetStorage", "audienceSegmentation.cloudletPolicy", "cpCode.value", "cpCode.value.cpCodeLimits", "edgeRedirector.cloudletPolicy", "failAction.netStorageHostname", "failAction.cpCode", "failAction.cpCode.cpCodeLimits", "firstPartyMarketing.cloudletPolicy", "firstPartyMarketingPlus.cloudletPolicy", "forwardRewrite.cloudletPolicy", "imageAndVideoManager.cpCodeOriginal", "imageAndVideoManager.cpCodeOriginal.cpCodeLimits", "imageAndVideoManager.cpCodeTransformed", "imageAndVideoManager.cpCodeTransformed.cpCodeLimits", "imageManager.cpCodeOriginal", "imageManager.cpCodeOriginal.cpCodeLimits", "imageManager.cpCodeTransformed", "imageManager.cpCodeTransformed.cpCodeLimits", "imageManagerVideo.cpCodeOriginal", "imageManagerVideo.cpCodeOriginal.cpCodeLimits", "imageManagerVideo.cpCodeTransformed", "imageManagerVideo.cpCodeTransformed.cpCodeLimits", "origin.netStorage", "origin.customCertificateAuthorities.subjectRDNs", "origin.customCertificateAuthorities.issuerRDNs", "origin.customCertificates.subjectRDNs", "origin.customCertificates.issuerRDNs", "phasedRelease.cloudletPolicy", "requestControl.cloudletPolicy", "requestControl.netStorage", "siteShield.ssmap", "visitorPrioritization.cloudletPolicy", "visitorPrioritization.waitingRoomCpCode", "visitorPrioritization.waitingRoomCpCode.cpCodeLimits", "visitorPrioritization.waitingRoomNetStorage", "webApplicationFirewall.firewallConfiguration", "matchCpCode.value", "matchCpCode.value.cpCodeLimits"},
	})
}

func getBehaviorsSchemaV20241021() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ad_scaler_circuit_breaker": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior works with `manifestRerouting` to provide the scale and reliability of Akamai network while simultaneously allowing third party partners to modify the requested media content with value-added features. The `adScalerCircuitBreaker` behavior specifies the fallback action in case the technology partner encounters errors and can't modify the requested media object. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"response_delay_based": {
						Optional:    true,
						Description: "Triggers a fallback action based on the delayed response from the technology partner's server.",
						Type:        schema.TypeBool,
					},
					"response_delay_threshold": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"500ms"}, false)),
						Optional:         true,
						Description:      "Specifies the maximum response delay that, if exceeded, triggers the fallback action.",
						Type:             schema.TypeString,
					},
					"response_code_based": {
						Optional:    true,
						Description: "Triggers a fallback action based on the response code from the technology partner's server.",
						Type:        schema.TypeBool,
					},
					"response_codes": {
						ValidateDiagFunc: validateRegexOrVariable("^(([0-9]{3})(,?))+$"),
						Optional:         true,
						Description:      "Specifies the codes in the partner's response that trigger the fallback action,  either `408`, `500`, `502`, `504`, `SAME_AS_RECEIEVED`, or `SPECIFY_YOUR_OWN` for a custom code.",
						Type:             schema.TypeString,
					},
					"fallback_action_response_code_based": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RETURN_AKAMAI_COPY", "RETURN_ERROR"}, false)),
						Optional:         true,
						Description:      "Specifies the fallback action.",
						Type:             schema.TypeString,
					},
					"return_error_response_code_based": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SAME_AS_RECEIVED", "408", "500", "502", "504", "SPECIFY_YOUR_OWN"}, false)),
						Optional:         true,
						Description:      "Specifies the error to include in the response to the client.",
						Type:             schema.TypeString,
					},
					"specify_your_own_response_code_based": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d{3}$"),
						Optional:         true,
						Description:      "Defines a custom error response.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"adaptive_acceleration": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Adaptive Acceleration uses HTTP/2 server push functionality with Ion properties to pre-position content and improve the performance of HTML page loading based on real user monitoring (RUM) timing data. It also helps browsers to preconnect to content that’s likely needed for upcoming requests. To use this behavior, make sure you enable the `http2` behavior. Use the `Adaptive Acceleration API` to report on the set of assets this feature optimizes. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"source": {
						Optional:    true,
						Description: "The source Adaptive Acceleration uses to gather the real user monitoring timing data, either `MPULSE` or `REAL_USER_MONITORING`.",
						Type:        schema.TypeString,
					},
					"title_http2_server_push": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_push": {
						Optional:    true,
						Description: "Recognizes resources like JavaScript, CSS, and images  based on gathered timing data and sends these resources to a browser as it's waiting for a response to the initial request for your website or app. See `Automatic Server Push` for more information.",
						Type:        schema.TypeBool,
					},
					"title_preconnect": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_preconnect": {
						Optional:    true,
						Description: "Allows browsers to anticipate what connections your site needs, and establishes those connections ahead of time. See `Automatic Preconnect` for more information.",
						Type:        schema.TypeBool,
					},
					"title_preload": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"preload_enable": {
						Optional:    true,
						Description: "Allows browsers to preload necessary fonts before they fetch and process other resources. See `Automatic Font Preload` for more information.",
						Type:        schema.TypeBool,
					},
					"ab_testing": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"ab_logic": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DISABLED", "CLOUDLETS", "MANUAL"}, false)),
						Optional:         true,
						Description:      "Specifies whether to use Adaptive Acceleration in an A/B testing environment. To include Adaptive Acceleration data in your A/B testing, specify the mode you want to apply. Otherwise, `DISABLED` by default. See `Add A/B testing to A2` for details.",
						Type:             schema.TypeString,
					},
					"cookie_name": {
						Optional:    true,
						Description: "This specifies the name of the cookie file used for redirecting the requests in the A/B testing environment.",
						Type:        schema.TypeString,
					},
					"intelligent_early_hints_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"intelligent_early_hints": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"compression": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"title_ro": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_ro": {
						Optional:    true,
						Description: "Enables the Resource Optimizer, which automates the compression and delivery of your `.css`, `.js`, and `.svg` content using a combination of Brotli and Zopfli compressions. The compression is performed offline, during a time to live that the feature automatically sets. See the `resourceOptimizer` and `resourceOptimizerExtendedCompatibility` behaviors for more details.",
						Type:        schema.TypeBool,
					},
					"title_brotli": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_brotli_compression": {
						Optional:    true,
						Description: "Applies Brotli compression, converting your origin content to cache on edge servers.",
						Type:        schema.TypeBool,
					},
					"enable_for_noncacheable": {
						Optional:    true,
						Description: "Applies Brotli compression to non-cacheable content.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"adaptive_image_compression": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "> **Note**: Starting from May 31, 2024, Adaptive Image Compression is no longer supported and the image compression configured through this functionality won't take place. As an alternative, we offer `Image & Video Manager`. It intelligently and automatically optimizes images and videos on the fly for every user. Reach out to your Akamai representatives for more information on this product. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"title_aic_mobile": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"compress_mobile": {
						Optional:    true,
						Description: "Adapts images served over cellular mobile networks.",
						Type:        schema.TypeBool,
					},
					"tier1_mobile_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-1 behavior.",
						Type:             schema.TypeString,
					},
					"tier1_mobile_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
					"tier2_mobile_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-2 cellular-network behavior.",
						Type:             schema.TypeString,
					},
					"tier2_mobile_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
					"tier3_mobile_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-3 cellular-network behavior.",
						Type:             schema.TypeString,
					},
					"tier3_mobile_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
					"title_aic_nonmobile": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"compress_standard": {
						Optional:    true,
						Description: "Adapts images served over non-cellular networks.",
						Type:        schema.TypeBool,
					},
					"tier1_standard_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-1 non-cellular network behavior.",
						Type:             schema.TypeString,
					},
					"tier1_standard_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
					"tier2_standard_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-2 non-cellular network behavior.",
						Type:             schema.TypeString,
					},
					"tier2_standard_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
					"tier3_standard_compression_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMPRESS", "BYPASS", "STRIP"}, false)),
						Optional:         true,
						Description:      "Specifies tier-3 non-cellular network behavior.",
						Type:             schema.TypeString,
					},
					"tier3_standard_compression_value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the compression percentage.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This specifies Akamai XML metadata. It can only be configured on your behalf by Akamai Professional Services. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "Human-readable description of what the XML block does.",
						Type:        schema.TypeString,
					},
					"xml": {
						Optional:    true,
						Description: "Akamai XML metadata.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"aggregated_reporting": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Configure a custom report that collects traffic data. The data is based on one to four variables, such as `sum`, `average`, `min`, and `max`. These aggregation attributes help compile traffic data summaries. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables aggregated reporting.",
						Type:        schema.TypeBool,
					},
					"report_name": {
						Optional:    true,
						Description: "The unique name of the aggregated report within the property. If you reconfigure any attributes or variables in the aggregated reporting behavior, update this field to a unique value to enable logging data in a new instance of the report.",
						Type:        schema.TypeString,
					},
					"attributes_count": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 4)),
						Optional:         true,
						Description:      "The number of attributes to include in the report, ranging from 1 to 4.",
						Type:             schema.TypeInt,
					},
					"attribute1": {
						Optional:    true,
						Description: "Specify a previously user-defined variable name as a report attribute. The values extracted for all attributes range from 0 to 20 characters.",
						Type:        schema.TypeString,
					},
					"attribute2": {
						Optional:    true,
						Description: "Specify a previously user-defined variable name as a report attribute. The values extracted for all attributes range from 0 to 20 characters.",
						Type:        schema.TypeString,
					},
					"attribute3": {
						Optional:    true,
						Description: "Specify a previously user-defined variable name as a report attribute. The values extracted for all attributes range from 0 to 20 characters.",
						Type:        schema.TypeString,
					},
					"attribute4": {
						Optional:    true,
						Description: "Specify a previously user-defined variable name as a report attribute. The values extracted for all attributes range from 0 to 20 characters.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"akamaizer": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This allows you to run regular expression substitutions over web pages. To apply this behavior, you need to match on a `contentType`. Contact Akamai Professional Services for help configuring the Akamaizer. See also the `akamaizerTag` behavior. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Akamaizer behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"akamaizer_tag": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This specifies HTML tags and replacement rules for hostnames used in conjunction with the `akamaizer` behavior. Contact Akamai Professional Services for help configuring the Akamaizer. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_hostname": {
						Optional:    true,
						Description: "Specifies the hostname to match on as a Perl-compatible regular expression.",
						Type:        schema.TypeString,
					},
					"replacement_hostname": {
						Optional:    true,
						Description: "Specifies the replacement hostname for the tag to use.",
						Type:        schema.TypeString,
					},
					"scope": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ATTRIBUTE", "URL_ATTRIBUTE", "BLOCK", "PAGE"}, false)),
						Optional:         true,
						Description:      "Specifies the part of HTML content the `tagsAttribute` refers to.",
						Type:             schema.TypeString,
					},
					"tags_attribute": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"A", "A_HREF", "IMG", "IMG_SRC", "SCRIPT", "SCRIPT_SRC", "LINK", "LINK_HREF", "TD", "TD_BACKGROUND", "TABLE", "TABLE_BACKGROUND", "IFRAME", "IFRAME_SRC", "AREA", "AREA_HREF", "BASE", "BASE_HREF", "FORM", "FORM_ACTION"}, false)),
						Optional:         true,
						Description:      "Specifies the tag or tag/attribute combination to operate on.",
						Type:             schema.TypeString,
					},
					"replace_all": {
						Optional:    true,
						Description: "Replaces all matches when enabled, otherwise replaces only the first match.",
						Type:        schema.TypeBool,
					},
					"include_tags_attribute": {
						Optional:    true,
						Description: "Whether to include the `tagsAttribute` value.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"all_http_in_cache_hierarchy": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow all HTTP request methods to be used for the edge's parent servers, useful to implement features such as `Site Shield`, `SureRoute`, and Tiered Distribution. (See the `siteShield`, `sureRoute`, and `tieredDistribution` behaviors.) This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables all HTTP requests for parent servers in the cache hierarchy.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_cloudlets_origins": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows Cloudlets Origins to determine the criteria, separately from the Property Manager, under which alternate `origin` definitions are assigned. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows you to assign custom origin definitions referenced in sub-rules by `cloudletsOrigin` labels. If disabled, all sub-rules are ignored.",
						Type:        schema.TypeBool,
					},
					"honor_base_directory": {
						Optional:    true,
						Description: "Prefixes any Cloudlet-generated origin path with a path defined by an Origin Base Path behavior. If no path is defined, it has no effect. If another Cloudlet policy already prepends the same Origin Base Path, the path is not duplicated.",
						Type:        schema.TypeBool,
					},
					"purge_origin_query_parameter": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "When purging content from a Cloudlets Origin, this specifies a query parameter name whose value is the specific named origin to purge.  Note that this only applies to content purge requests, for example when using the `Content Control Utility API`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"allow_delete": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow HTTP requests using the DELETE method. By default, GET, HEAD, and OPTIONS requests are allowed, and all other methods result in a 501 error. Such content does not cache, and any DELETE requests pass to the origin. See also the `allowOptions`, `allowPatch`, `allowPost`, and `allowPut` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows DELETE requests. Content does `not` cache.",
						Type:        schema.TypeBool,
					},
					"allow_body": {
						Optional:    true,
						Description: "Allows data in the body of the DELETE request.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_https_cache_key_sharing": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "HTTPS cache key sharing allows HTTP requests to be served from an HTTPS cache. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables HTTPS cache key sharing.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_https_downgrade": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Passes HTTPS requests to origin as HTTP. This is useful when incorporating Standard TLS or Akamai's shared certificate delivery security with an origin that serves HTTP traffic. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Downgrades to HTTP protocol for the origin server.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_options": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "GET, HEAD, and OPTIONS requests are allowed by default. All other HTTP methods result in a 501 error. For full support of Cross-Origin Resource Sharing (CORS), you need to allow requests that use the OPTIONS method. If you're using the `corsSupport` behavior, do not disable OPTIONS requests. The response to an OPTIONS request is not cached, so the request always goes through the Akamai network to your origin, unless you use the `constructResponse` behavior to send responses directly from the Akamai network. See also the `allowDelete`, `allowPatch`, `allowPost`, and `allowPut` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Set this to `true` to reflect the default policy where edge servers allow the OPTIONS method, without caching the response. Set this to `false` to deny OPTIONS requests and respond with a 501 error.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_patch": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow HTTP requests using the PATCH method. By default, GET, HEAD, and OPTIONS requests are allowed, and all other methods result in a 501 error. Such content does not cache, and any PATCH requests pass to the origin. See also the `allowDelete`, `allowOptions`, `allowPost`, and `allowPut` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows PATCH requests. Content does `not` cache.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_post": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow HTTP requests using the POST method. By default, GET, HEAD, and OPTIONS requests are allowed, and POST requests are denied with 403 error. All other methods result in a 501 error. See also the `allowDelete`, `allowOptions`, `allowPatch`, and `allowPut` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows POST requests.",
						Type:        schema.TypeBool,
					},
					"allow_without_content_length": {
						Optional:    true,
						Description: "By default, POST requests also require a `Content-Length` header, or they result in a 411 error. With this option enabled with no specified `Content-Length`, the edge server relies on a `Transfer-Encoding` header to chunk the data. If neither header is present, it assumes the request has no body, and it adds a header with a `0` value to the forward request.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_put": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow HTTP requests using the PUT method.  By default, GET, HEAD, and OPTIONS requests are allowed, and all other methods result in a 501 error. Such content does not cache, and any PUT requests pass to the origin. See also the `allowDelete`, `allowOptions`, `allowPatch`, and `allowPost` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows PUT requests. Content does `not` cache.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"allow_transfer_encoding": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Controls whether to allow or deny Chunked Transfer Encoding (CTE) requests to pass to your origin. If your origin supports CTE, you should enable this behavior. This behavior also protects against a known issue when pairing `http2` and `webdav` behaviors within the same rule tree, in which case it's required. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows Chunked Transfer Encoding requests.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"alt_svc_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Sets the maximum age value for the Alternative Services (`Alt-Svc`) header. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"max_age": {
						Optional:    true,
						Description: "Specifies the `max-age` value in seconds for the `Alt-Svc` header. The default `max-age` for an `Alt-Svc` header is 93600 seconds (26 hours).",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"api_prioritization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the API Prioritization Cloudlet, which maintains continuity in user experience by serving an alternate static response when load is too high. You can configure rules using either the Cloudlets Policy Manager application or the `Cloudlets API`. Use this feature serve static API content, such as fallback JSON data.  To serve non-API HTML content, use the `visitorPrioritization` behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Activates the API Prioritization feature.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to apply the Cloudlet shared policy to an unlimited number of properties within your account. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet shared policy to use with this behavior. Use the `Cloudlets API` to list available shared policies.",
						Type:        schema.TypeInt,
					},
					"label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "A label to distinguish this API Prioritization policy from any others in the same property.",
						Type:             schema.TypeString,
					},
					"use_throttled_cp_code": {
						Optional:    true,
						Description: "Specifies whether to apply an alternative CP code for requests served the alternate response.",
						Type:        schema.TypeBool,
					},
					"throttled_cp_code": {
						Optional:    true,
						Description: "Specifies the CP code as an object. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"use_throttled_status_code": {
						Optional:    true,
						Description: "Allows you to assign a specific HTTP response code to a throttled request.",
						Type:        schema.TypeBool,
					},
					"throttled_status_code": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d{3}$"),
						Optional:         true,
						Description:      "Specifies the HTTP response code for requests that receive the alternate response.",
						Type:             schema.TypeInt,
					},
					"net_storage": {
						Optional:    true,
						Description: "Specify the NetStorage domain that contains the alternate response.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"net_storage_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specify the full NetStorage path for the alternate response, including trailing file name.",
						Type:             schema.TypeString,
					},
					"alternate_response_cache_ttl": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 30)),
						Optional:         true,
						Description:      "Specifies the alternate response's time to live in the cache, `5` minutes by default.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"application_load_balancer": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the Application Load Balancer Cloudlet, which automates load balancing based on configurable criteria. To configure this behavior, use either the Cloudlets Policy Manager or the `Cloudlets API` to set up a policy. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Activates the Application Load Balancer Cloudlet.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "A label to distinguish this Application Load Balancer policy from any others within the same property.",
						Type:             schema.TypeString,
					},
					"stickiness_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"stickiness_cookie_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "NEVER", "ON_BROWSER_CLOSE", "FIXED_DATE", "DURATION", "ORIGIN_SESSION"}, false)),
						Optional:         true,
						Description:      "Determines how a cookie persistently associates the client with a load-balanced origin.",
						Type:             schema.TypeString,
					},
					"stickiness_expiration_date": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies when the cookie expires.",
						Type:             schema.TypeString,
					},
					"stickiness_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Sets how long it is before the cookie expires.",
						Type:             schema.TypeString,
					},
					"stickiness_refresh": {
						Optional:    true,
						Description: "Extends the duration of the cookie with each new request. When enabled, the `DURATION` thus specifies the latency between requests that would cause the cookie to expire.",
						Type:        schema.TypeBool,
					},
					"origin_cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name for your session cookie.",
						Type:             schema.TypeString,
					},
					"specify_stickiness_cookie_domain": {
						Optional:    true,
						Description: "Specifies whether to use a cookie domain with the stickiness cookie, to tell the browser to which domain to send the cookie.",
						Type:        schema.TypeBool,
					},
					"stickiness_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the domain to track the stickiness cookie.",
						Type:             schema.TypeString,
					},
					"stickiness_cookie_automatic_salt": {
						Optional:    true,
						Description: "Sets whether to assign a `salt` value automatically to the cookie to prevent manipulation by the user. You should not enable this if sharing the population cookie across more than one property.",
						Type:        schema.TypeBool,
					},
					"stickiness_cookie_salt": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the stickiness cookie's salt value. Use this option to share the cookie across many properties.",
						Type:             schema.TypeString,
					},
					"stickiness_cookie_set_http_only_flag": {
						Optional:    true,
						Description: "Ensures the cookie is transmitted only over HTTP.",
						Type:        schema.TypeBool,
					},
					"all_down_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"all_down_net_storage": {
						Optional:    true,
						Description: "Specifies a NetStorage account for a static maintenance page as a fallback when no origins are available.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"all_down_net_storage_file": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the fallback maintenance page's filename, expressed as a full path from the root of the NetStorage server.",
						Type:             schema.TypeString,
					},
					"all_down_status_code": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d{3}$"),
						Optional:         true,
						Description:      "Specifies the HTTP response code when all load-balancing origins are unavailable.",
						Type:             schema.TypeString,
					},
					"failover_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"failover_status_codes": {
						Optional:    true,
						Description: "Specifies a set of HTTP status codes that signal a failure on the origin, in which case the cookie that binds the client to that origin is invalidated and the client is rerouted to another available origin.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"failover_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTOMATIC", "MANUAL", "DISABLED"}, false)),
						Optional:         true,
						Description:      "Determines what to do if an origin fails.",
						Type:             schema.TypeString,
					},
					"failover_origin_map": {
						Optional:    true,
						Description: "Specifies a fixed set of failover mapping rules.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"from_origin_id": {
									ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-\\.]+$"),
									Optional:         true,
									Description:      "Specifies the origin whose failure triggers the mapping rule.",
									Type:             schema.TypeString,
								},
								"to_origin_ids": {
									Optional:    true,
									Description: "Requests stuck to the `fromOriginId` origin retry for each alternate origin `toOriginIds`, until one succeeds.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"failover_attempts_threshold": {
						Optional:    true,
						Description: "Sets the number of failed requests that would trigger the failover process.",
						Type:        schema.TypeInt,
					},
					"cached_content_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"allow_cache_prefresh": {
						Optional:    true,
						Description: "Allows the cache to prefresh.  Only appropriate if all origins serve the same content for the same URL.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"audience_segmentation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows you to divide your users into different segments based on a persistent cookie. You can configure rules using either the Cloudlets Policy Manager application or the `Cloudlets API`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Audience Segmentation cloudlet feature.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to use a shared policy for a Cloudlet. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "This identifies the Cloudlet shared policy to use with this behavior. You can list available shared policies with the `Cloudlets API`.",
						Type:        schema.TypeInt,
					},
					"label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies a suffix to append to the cookie name. This helps distinguish this audience segmentation policy from any others within the same property.",
						Type:             schema.TypeString,
					},
					"segment_tracking_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"segment_tracking_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IN_QUERY_PARAM", "IN_COOKIE_HEADER", "IN_CUSTOM_HEADER", "NONE"}, false)),
						Optional:         true,
						Description:      "Specifies the method to pass segment information to the origin. The Cloudlet passes the rule applied to a given request location.",
						Type:             schema.TypeString,
					},
					"segment_tracking_query_param": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This query parameter specifies the name of the segmentation rule.",
						Type:             schema.TypeString,
					},
					"segment_tracking_cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This cookie name specifies the name of the segmentation rule.",
						Type:             schema.TypeString,
					},
					"segment_tracking_custom_header": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This custom HTTP header specifies the name of the segmentation rule.",
						Type:             schema.TypeString,
					},
					"population_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"population_cookie_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NEVER", "ON_BROWSER_CLOSE", "DURATION"}, false)),
						Optional:         true,
						Description:      "Specifies when the segmentation cookie expires.",
						Type:             schema.TypeString,
					},
					"population_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the lifetime of the segmentation cookie.",
						Type:             schema.TypeString,
					},
					"population_refresh": {
						Optional:    true,
						Description: "If disabled, sets the expiration time only if the cookie is not yet present in the request.",
						Type:        schema.TypeBool,
					},
					"specify_population_cookie_domain": {
						Optional:    true,
						Description: "Whether to specify a cookie domain with the population cookie. It tells the browser to which domain to send the cookie.",
						Type:        schema.TypeBool,
					},
					"population_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the domain to track the population cookie.",
						Type:             schema.TypeString,
					},
					"population_cookie_automatic_salt": {
						Optional:    true,
						Description: "Whether to assign a `salt` value automatically to the cookie to prevent manipulation by the user. You should not enable if sharing the population cookie across more than one property.",
						Type:        schema.TypeBool,
					},
					"population_cookie_salt": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the cookie's salt value. Use this option to share the cookie across many properties.",
						Type:             schema.TypeString,
					},
					"population_cookie_include_rule_name": {
						Optional:    true,
						Description: "When enabled, includes in the session cookie the name of the rule in which this behavior appears.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"auto_domain_validation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows standard TLS domain validated certificates to renew automatically. Apply it after using the `Certificate Provisioning System` to request a certificate for a hostname.  To provision certificates programmatically, see the `Certificate Provisioning System API`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"autodv": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"base_directory": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Prefix URLs sent to the origin with a base path. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^/([^:#\\[\\]@/?]+/)*$"),
						Optional:         true,
						Description:      "Specifies the base path of content on your origin server. The value needs to begin and end with a slash (`/`) character, for example `/parent/child/`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"boss_beaconing": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Triggers diagnostic data beacons for use with BOSS, Akamai's monitoring and diagnostics system. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enable diagnostic data beacons.",
						Type:        schema.TypeBool,
					},
					"cpcodes": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9 ]*$"),
						Optional:         true,
						Description:      "The space-separated list of CP codes that trigger the beacons. You need to specify the same set of CP codes within BOSS.",
						Type:             schema.TypeString,
					},
					"request_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EDGE", "EDGE_MIDGRESS"}, false)),
						Optional:         true,
						Description:      "Specify when to trigger a beacon.",
						Type:             schema.TypeString,
					},
					"forward_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"MIDGRESS", "ORIGIN", "MIDGRESS_ORIGIN"}, false)),
						Optional:         true,
						Description:      "Specify when to trigger a beacon.",
						Type:             schema.TypeString,
					},
					"sampling_frequency": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SAMPLING_FREQ_0_0", "SAMPLING_FREQ_0_1"}, false)),
						Optional:         true,
						Description:      "Specifies a sampling frequency or disables beacons.",
						Type:             schema.TypeString,
					},
					"conditional_sampling_frequency": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CONDITIONAL_SAMPLING_FREQ_0_0", "CONDITIONAL_SAMPLING_FREQ_0_1", "CONDITIONAL_SAMPLING_FREQ_0_2", "CONDITIONAL_SAMPLING_FREQ_0_3"}, false)),
						Optional:         true,
						Description:      "Specifies a conditional sampling frequency or disables beacons.",
						Type:             schema.TypeString,
					},
					"conditional_http_status": {
						Optional:    true,
						Description: "Specifies the set of response status codes or ranges that trigger the beacon.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"conditional_error_pattern": {
						Optional:    true,
						Description: "A space-separated set of error patterns that trigger beacons to conditional feeds. Each pattern can include wildcards, where `?` matches a single character and `*` matches zero or more characters. For example, `*CONNECT* *DENIED*` matches two different words as substrings.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"breadcrumbs": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Provides per-HTTP transaction visibility into a request for content, regardless of how deep the request goes into the Akamai platform. The `Akamai-Request-BC` response header includes various data, such as network health and the location in the Akamai network used to serve content, which simplifies log review for troubleshooting. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Breadcrumbs feature.",
						Type:        schema.TypeBool,
					},
					"opt_mode": {
						Optional:    true,
						Description: "Specifies whether to include Breadcrumbs data in the response header. To bypass the current `optMode`, append the opposite `ak-bc` query string to each request from your player.",
						Type:        schema.TypeBool,
					},
					"logging_enabled": {
						Optional:    true,
						Description: "Whether to collect all Breadcrumbs data in logs, including the response headers sent a requesting client. This can also be helpful if you're using `DataStream 2` to retrieve log data. This way, all Breadcrumbs data is carried in the logs it uses.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"break_connection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior simulates an origin connection problem, typically to test an accompanying `failAction` policy. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the break connection behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"brotli": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Accesses Brotli-compressed assets from your origin and caches them on edge servers. This doesn't compress resources within the content delivery network in real time. You need to set up Brotli compression separately on your origin. If a requesting client doesn't support Brotli, edge servers deliver non-Brotli resources. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Fetches Brotli-compressed assets from your origin and caches them on edge servers.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"cache_error": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "With the `caching` behavior also enabled, this caches the origin's error responses to decrease server load. It applies for 10 seconds by default to the following HTTP codes: `204`, `305`, `404`, `405`, `501`, `502`, `503`, `504`, and `505`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Activates the error-caching behavior.",
						Type:        schema.TypeBool,
					},
					"ttl": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Overrides the default caching duration of `10s`. Note that if set to `0`, it is equivalent to `no-cache`, which forces revalidation and may cause a traffic spike. This can be counterproductive when, for example, the origin is producing an error code of `500`.",
						Type:             schema.TypeString,
					},
					"preserve_stale": {
						Optional:    true,
						Description: "When enabled, the edge server preserves stale cached objects when the origin returns `500`, `502`, `503`, and `504` error codes. This avoids re-fetching and re-caching content after transient errors.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"cache_id": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Controls which query parameters, headers, and cookies are included in or excluded from the cache key identifier. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"rule": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INCLUDE_QUERY_PARAMS", "INCLUDE_COOKIES", "INCLUDE_HEADERS", "EXCLUDE_QUERY_PARAMS", "INCLUDE_ALL_QUERY_PARAMS", "INCLUDE_VARIABLE", "INCLUDE_URL"}, false)),
						Optional:         true,
						Description:      "Specifies how to modify the cache ID.",
						Type:             schema.TypeString,
					},
					"include_value": {
						Optional:    true,
						Description: "Includes the value of the specified elements in the cache ID. Otherwise only their names are included.",
						Type:        schema.TypeBool,
					},
					"optional": {
						Optional:    true,
						Description: "Requires the behavior's specified elements to be present for content to cache. When disabled, requests that lack the specified elements are still cached.",
						Type:        schema.TypeBool,
					},
					"elements": {
						Optional:    true,
						Description: "Specifies the names of the query parameters, cookies, or headers to include or exclude from the cache ID.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"variable_name": {
						Optional:    true,
						Description: "Specifies the name of the variable you want to include in the cache key.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"cache_key_ignore_case": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "By default, cache keys are generated under the assumption that path and filename components are case-sensitive, so that `File.html` and `file.html` use separate cache keys. Enabling this behavior forces URL components whose case varies to resolve to the same cache key. Enable this behavior if your origin server is already case-insensitive, such as those based on Microsoft IIS. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Ignores case when forming cache keys.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"cache_key_query_params": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "By default, cache keys are formed as URLs with full query strings. This behavior allows you to consolidate cached objects based on specified sets of query parameters. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INCLUDE_ALL_PRESERVE_ORDER", "INCLUDE_ALL_ALPHABETIZE_ORDER", "IGNORE_ALL", "INCLUDE", "IGNORE"}, false)),
						Optional:         true,
						Description:      "Configures how sets of query string parameters translate to cache keys. Be careful not to ignore any parameters that result in substantially different content, as it is `not` reflected in the cached object.",
						Type:             schema.TypeString,
					},
					"parameters": {
						Optional:    true,
						Description: "Specifies the set of parameter field names to include in or exclude from the cache key. By default, these match the field names as string prefixes.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"exact_match": {
						Optional:    true,
						Description: "When enabled, `parameters` needs to match exactly. Keep disabled to match string prefixes.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"cache_key_rewrite": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior rewrites a default cache key's path. Contact Akamai Professional Services for help configuring it. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"purge_key": {
						ValidateDiagFunc: validateRegexOrVariable("^[\\w-]+$"),
						Optional:         true,
						Description:      "Specifies the new cache key path as an alphanumeric value.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"cache_post": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "By default, POST requests are passed to the origin. This behavior overrides the default, and allows you to cache POST responses. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables caching of POST responses.",
						Type:        schema.TypeBool,
					},
					"use_body": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IGNORE", "MD5", "QUERY"}, false)),
						Optional:         true,
						Description:      "Define how and whether to use the POST message body as a cache key.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"cache_redirect": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Controls the caching of HTTP 302 and 307 temporary redirects. By default, Akamai edge servers don't cache them. Enabling this behavior instructs edge servers to allow these redirects to be cached the same as HTTP 200 responses. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the redirect caching behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"cache_tag": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This adds a cache tag to the requested object. With cache tags, you can flexibly fast purge tagged segments of your cached content. You can either define these tags with an `Edge-Cache-Tag` header at the origin server level, or use this behavior to directly add a cache tag to the object as the edge server caches it. The `cacheTag` behavior can only take a single value, including a variable. If you want to specify more tags for an object, add a few instances of this behavior to your configuration. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"tag": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9\\&\\'\\^\\-\\$\\!\\`\\#\\%\\.\\+\\~\\_\\|\\/]+$"),
						Optional:         true,
						Description:      "Specifies the cache tag you want to add to your cached content. A cache tag is only added when the object is first added to cache. A single cache tag can't exceed 128 characters and can only include alphanumeric characters, plus this class of characters: ```[!#$%'+./^_`|~-]```",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"cache_tag_visible": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Cache tags are comma-separated string values you define within an `Edge-Cache-Tag` header. You can use them to flexibly fast purge tagged segments of your cached content. You can either define these headers at the origin server level, or use the `modifyOutgoingResponseHeader` behavior to configure them at the edge.  Apply this behavior to confirm you're deploying the intended set of cache tags to your content. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NEVER", "PRAGMA_HEADER", "ALWAYS"}, false)),
						Optional:         true,
						Description:      "Specifies when to include the `Edge-Cache-Tag` in responses.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"caching": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Control content caching on edge servers: whether or not to cache, whether to honor the origin's caching headers, and for how long to cache.  Note that any `NO_STORE` or `BYPASS_CACHE` HTTP headers set on the origin's content override this behavior. For more details on how caching works in Property Manager, see the `Learn about caching` section in the guide. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"MAX_AGE", "NO_STORE", "BYPASS_CACHE", "CACHE_CONTROL_AND_EXPIRES", "CACHE_CONTROL", "EXPIRES"}, false)),
						Optional:         true,
						Description:      "Specify the caching option.",
						Type:             schema.TypeString,
					},
					"must_revalidate": {
						Optional:    true,
						Description: "Determines what to do once the cached content has expired, by which time the Akamai platform should have re-fetched and validated content from the origin. If enabled, only allows the re-fetched content to be served. If disabled, may serve stale content if the origin is unavailable.",
						Type:        schema.TypeBool,
					},
					"ttl": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "The maximum time content may remain cached. Setting the value to `0` is the same as setting a `no-cache` header, which forces content to revalidate.",
						Type:             schema.TypeString,
					},
					"default_ttl": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Set the `MAX_AGE` header for the cached content.",
						Type:             schema.TypeString,
					},
					"cache_control_directives": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enhanced_rfc_support": {
						Optional:    true,
						Description: "This enables honoring particular `Cache-Control` header directives from the origin. Supports all official `RFC 7234` directives except for `no-transform`.",
						Type:        schema.TypeBool,
					},
					"cacheability_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"honor_no_store": {
						Optional:    true,
						Description: "Instructs edge servers not to cache the response when the origin response includes the `no-store` directive.",
						Type:        schema.TypeBool,
					},
					"honor_private": {
						Optional:    true,
						Description: "Instructs edge servers not to cache the response when the origin response includes the `private` directive.",
						Type:        schema.TypeBool,
					},
					"honor_no_cache": {
						Optional:    true,
						Description: "With the `no-cache` directive present in the response, this instructs edge servers to validate or refetch the response for each request. Effectively, set the time to live `ttl` to zero seconds.",
						Type:        schema.TypeBool,
					},
					"expiration_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"honor_max_age": {
						Optional:    true,
						Description: "This instructs edge servers to cache the object for a length of time set by the `max-age` directive in the response. When present in the origin response, this directive takes precedence over the `max-age` directive and the `defaultTtl` setting.",
						Type:        schema.TypeBool,
					},
					"honor_s_maxage": {
						Optional:    true,
						Description: "Instructs edge servers to cache the object for a length of time set by the `s-maxage` directive in the response. When present in the origin response, this directive takes precedence over the `max-age` directive and the `defaultTtl` setting.",
						Type:        schema.TypeBool,
					},
					"revalidation_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"honor_must_revalidate": {
						Optional:    true,
						Description: "This instructs edge servers to successfully revalidate with the origin server before using stale objects in the cache to satisfy new requests.",
						Type:        schema.TypeBool,
					},
					"honor_proxy_revalidate": {
						Optional:    true,
						Description: "With the `proxy-revalidate` directive present in the response, this instructs edge servers to successfully revalidate with the origin server before using stale objects in the cache to satisfy new requests.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"central_authorization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Forward client requests to the origin server for authorization, along with optional `Set-Cookie` headers, useful when you need to maintain tight access control. The edge server forwards an `If-Modified-Since` header, to which the origin needs to respond with a `304` (Not-Modified) HTTP status when authorization succeeds. If so, the edge server responds to the client with the cached object, since it does not need to be re-acquired from the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the centralized authorization behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"chase_redirects": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Controls whether the edge server chases any redirects served from the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows edge servers to chase redirects.",
						Type:        schema.TypeBool,
					},
					"limit": {
						Optional:    true,
						Description: "Specifies, as a string, the maximum number of redirects to follow.",
						Type:        schema.TypeString,
					},
					"serve404": {
						Optional:    true,
						Description: "Once the redirect `limit` is reached, enabling this option serves an HTTP `404` (Not Found) error instead of the last redirect.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"client_certificate_auth": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Sends a `Client-To-Edge` header to your origin server with details from the mutual TLS certificate sent from the requesting client to the edge network. This establishes transitive trust between the client and your origin server. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "Constructs the `Client-To-Edge` authentication header using information from the client to edge mTLS handshake and forwards it to your origin. You can configure your origin to acknowledge the header to enable transitive trust. Some form of the client x.509 certificate needs to be included in the header. You can include the full certificate or specific attributes.",
						Type:        schema.TypeBool,
					},
					"enable_complete_client_certificate": {
						Optional:    true,
						Description: "Whether to include the complete client certificate in the header, in its binary (DER) format. DER-formatted certificates leave out the `BEGIN CERTIFICATE/END CERTIFICATE` statements and most often use the `.der` extension. Alternatively, you can specify individual `clientCertificateAttributes` you want included in the request.",
						Type:        schema.TypeBool,
					},
					"client_certificate_attributes": {
						Optional:    true,
						Description: "Specify client certificate attributes to include in the `Client-To-Edge` authentication header that's sent to your origin server.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"enable_client_certificate_validation_status": {
						Optional:    true,
						Description: "Whether to include the current validation status of the client certificate in the `Client-To-Edge` authentication header. This verifies the validation status of the certificate, regardless of the certificate attributes you're including in the header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"client_characteristics": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the client ecosystem. Akamai uses this information to optimize your metadata configuration, which may result in better end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"country": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GLOBAL", "GLOBAL_US_CENTRIC", "GLOBAL_EU_CENTRIC", "GLOBAL_ASIA_CENTRIC", "EUROPE", "NORTH_AMERICA", "SOUTH_AMERICA", "NORDICS", "ASIA_PACIFIC", "AUSTRALIA", "GERMANY", "INDIA", "ITALY", "JAPAN", "TAIWAN", "UNITED_KINGDOM", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Specifies the client request's geographic region.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"cloud_interconnects": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Cloud Interconnects forwards traffic from edge servers to your cloud origin through Private Network Interconnects (PNIs), helping to reduce the egress costs at the origin. Supports origins hosted by Google Cloud Provider (GCP). This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Channels the traffic to maximize the egress discount at the origin.",
						Type:        schema.TypeBool,
					},
					"cloud_locations": {
						Optional:    true,
						Description: "Specifies the geographical locations of your cloud origin. You should enable Cloud Interconnects only if your origin is in one of these locations, since GCP doesn't provide a discount for egress traffic for any other regions.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"cloud_wrapper": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "`Cloud Wrapper` maximizes origin offload for large libraries of video, game, and software downloads by optimizing data caches in regions nearest to your origin. You can't use this behavior in conjunction with `sureRoute` or `tieredDistribution`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables Cloud Wrapper behavior.",
						Type:        schema.TypeBool,
					},
					"location": {
						Optional:    true,
						Description: "The location you want to distribute your Cloud Wrapper cache space to. This behavior allows all locations configured in your Cloud Wrapper configuration.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"cloud_wrapper_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Your account representative uses this behavior to implement a customized failover configuration on your behalf. Use Cloud Wrapper Advanced with an enabled `cloudWrapper` behavior in the same rule. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables failover for Cloud Wrapper.",
						Type:        schema.TypeBool,
					},
					"failover_map": {
						Optional:    true,
						Description: "Specifies the failover map to handle Cloud Wrapper failures. Contact your account representative for more information.",
						Type:        schema.TypeString,
					},
					"custom_failover_map": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z][a-zA-Z0-9-]*$"),
						Optional:         true,
						Description:      "Specifies the custom failover map to handle Cloud Wrapper failures. Contact your account representative for more information.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"common_media_client_data": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Use this behavior to send expanded playback information as CMCD metadata in requests from a media player. Edge servers may use this metadata for segment prefetching to optimize your content's delivery, or for logging. For more details and additional property requirements, see the `Adaptive Media Delivery` documentation. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable_cmcd_segment_prefetch": {
						Optional:    true,
						Description: "Uses Common Media Client Data (CMCD) metadata to determine the segment URLs your origin server prefetches to speed up content delivery.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"conditional_origin": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origin_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-\\.]+$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"construct_response": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior constructs an HTTP response, complete with HTTP status code and body, to serve from the edge independently of your origin. For example, you might want to send a customized response if the URL doesn't point to an object on the origin server, or if the end user is not authorized to view the requested content. You can use it with all request methods you allow for your property, including POST. For more details, see the `allowOptions`, `allowPatch`, `allowPost`, `allowPut`, and `allowDelete` behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Serves the custom response.",
						Type:        schema.TypeBool,
					},
					"body": {
						Optional:    true,
						Description: "HTML response of up to 2000 characters to send to the end-user client.",
						Type:        schema.TypeString,
					},
					"response_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{200, 404, 401, 403, 405, 417, 500, 501, 502, 503, 504})),
						Optional:         true,
						Description:      "The HTTP response code to send to the end-user client.",
						Type:             schema.TypeInt,
					},
					"force_eviction": {
						Optional:    true,
						Description: "For GET requests from clients, this forces edge servers to evict the underlying object from cache. Defaults to `false`.",
						Type:        schema.TypeBool,
					},
					"ignore_purge": {
						Optional:    true,
						Description: "Whether to ignore the custom response when purging.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"content_characteristics": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"object_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the size of the object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"USER_GENERATED", "WEB_OBJECTS", "SOFTWARE", "IMAGES", "OTHER_OBJECTS", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the type of content.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_characteristics_amd": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SD", "HD", "ULTRA_HD", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the quality of media content.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"hls": {
						Optional:    true,
						Description: "Enable delivery of HLS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_duration_hls_custom": {
						Optional:    true,
						Description: "Customizes the number of seconds for the segment.",
						Type:        schema.TypeFloat,
					},
					"segment_size_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"hds": {
						Optional:    true,
						Description: "Enable delivery of HDS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_duration_hds_custom": {
						Optional:    true,
						Description: "Customizes the number of seconds for the fragment.",
						Type:        schema.TypeInt,
					},
					"segment_size_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"dash": {
						Optional:    true,
						Description: "Enable delivery of DASH media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_duration_dash_custom": {
						Optional:    true,
						Description: "Customizes the number of seconds for the segment.",
						Type:        schema.TypeInt,
					},
					"segment_size_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"smooth": {
						Optional:    true,
						Description: "Enable delivery of Smooth media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_duration_smooth_custom": {
						Optional:    true,
						Description: "Customizes the number of seconds for the fragment.",
						Type:        schema.TypeFloat,
					},
					"segment_size_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_characteristics_dd": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"object_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the size of the object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VIDEO", "SOFTWARE", "SOFTWARE_PATCH", "GAME", "GAME_PATCH", "OTHER_DOWNLOADS", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the type of content.",
						Type:             schema.TypeString,
					},
					"optimize_option": {
						Optional:    true,
						Description: "Optimizes the delivery throughput and download times for large files.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"content_characteristics_wsd_large_file": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content, specifically targeted to delivering large files. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"object_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the size of the object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VIDEO", "SOFTWARE", "SOFTWARE_PATCH", "GAME", "GAME_PATCH", "OTHER_DOWNLOADS", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the type of content.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_characteristics_wsd_live": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content, specifically targeted to delivering live video. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SD", "HD", "ULTRA_HD", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the quality of media content.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"hls": {
						Optional:    true,
						Description: "Enable delivery of HLS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_size_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"hds": {
						Optional:    true,
						Description: "Enable delivery of HDS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_size_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"dash": {
						Optional:    true,
						Description: "Enable delivery of DASH media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_size_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"smooth": {
						Optional:    true,
						Description: "Enable delivery of Smooth media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_size_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_characteristics_wsd_vod": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the delivered content, specifically targeted to delivering on-demand video. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the total size of the content library delivered.",
						Type:             schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SD", "HD", "ULTRA_HD", "OTHER", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the quality of media content.",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LONG_TAIL", "ALL_POPULAR", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Optimize based on the content's expected popularity.",
						Type:             schema.TypeString,
					},
					"hls": {
						Optional:    true,
						Description: "Enable delivery of HLS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_size_hls": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"hds": {
						Optional:    true,
						Description: "Enable delivery of HDS media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_size_hds": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"dash": {
						Optional:    true,
						Description: "Enable delivery of DASH media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual segments.",
						Type:             schema.TypeString,
					},
					"segment_size_dash": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
					"smooth": {
						Optional:    true,
						Description: "Enable delivery of Smooth media.",
						Type:        schema.TypeBool,
					},
					"segment_duration_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SEGMENT_DURATION_2S", "SEGMENT_DURATION_4S", "SEGMENT_DURATION_6S", "SEGMENT_DURATION_8S", "SEGMENT_DURATION_10S"}, false)),
						Optional:         true,
						Description:      "Specifies the duration of individual fragments.",
						Type:             schema.TypeString,
					},
					"segment_size_smooth": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "ONE_MB_TO_TEN_MB", "TEN_MB_TO_100_MB", "GREATER_THAN_100MB", "UNKNOWN", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies the size of the media object retrieved from the origin.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_pre_position": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Content Preposition. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Content PrePosition behavior.",
						Type:        schema.TypeBool,
					},
					"source_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ORIGIN"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"targets": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLOUDWRAPPER"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"first_location": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeString,
					},
					"second_location": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"content_targeting_protection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Content Targeting is based on `EdgeScape`, Akamai's location-based access control system.  You can use it to allow or deny access to a set of geographic regions or IP addresses. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Content Targeting feature.",
						Type:        schema.TypeBool,
					},
					"geo_protection_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_geo_protection": {
						Optional:    true,
						Description: "When enabled, verifies IP addresses are unique to specific geographic regions.",
						Type:        schema.TypeBool,
					},
					"geo_protection_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY"}, false)),
						Optional:         true,
						Description:      "Specifies how to handle requests.",
						Type:             schema.TypeString,
					},
					"countries": {
						Optional:    true,
						Description: "Specifies a set of two-character ISO 3166 country codes from which to allow or deny traffic. See `EdgeScape Data Codes` for a list.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"regions": {
						Optional:    true,
						Description: "Specifies a set of ISO 3166-2 regional codes from which to allow or deny traffic. See `EdgeScape Data Codes` for a list.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"dmas": {
						Optional:    true,
						Description: "Specifies the set of Designated Market Area codes from which to allow or deny traffic.  See `EdgeScape Data Codes` for a list.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"override_ip_addresses": {
						Optional:    true,
						Description: "Specify a set of IP addresses or CIDR blocks that exceptions to the set of included or excluded areas.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"enable_geo_redirect_on_deny": {
						Optional:    true,
						Description: "When enabled, redirects denied requests rather than responding with an error code.",
						Type:        schema.TypeBool,
					},
					"geo_redirect_url": {
						Optional:    true,
						Description: "This specifies the full URL to the redirect page for denied requests.",
						Type:        schema.TypeString,
					},
					"ip_protection_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_ip_protection": {
						Optional:    true,
						Description: "Allows you to control access to your content from specific sets of IP addresses and CIDR blocks.",
						Type:        schema.TypeBool,
					},
					"ip_protection_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY"}, false)),
						Optional:         true,
						Description:      "Specifies how to handle requests.",
						Type:             schema.TypeString,
					},
					"ip_addresses": {
						Optional:    true,
						Description: "Specify a set of IP addresses or CIDR blocks to allow or deny.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"enable_ip_redirect_on_deny": {
						Optional:    true,
						Description: "When enabled, redirects denied requests rather than responding with an error code.",
						Type:        schema.TypeBool,
					},
					"ip_redirect_url": {
						Optional:    true,
						Description: "This specifies the full URL to the redirect page for denied requests.",
						Type:        schema.TypeString,
					},
					"referrer_protection_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_referrer_protection": {
						Optional:    true,
						Description: "Allows you allow traffic from certain referring websites, and disallow traffic from unauthorized sites that hijack those links.",
						Type:        schema.TypeBool,
					},
					"referrer_protection_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY"}, false)),
						Optional:         true,
						Description:      "Specify the action to take.",
						Type:             schema.TypeString,
					},
					"referrer_domains": {
						Optional:    true,
						Description: "Specifies the set of domains from which to allow or deny traffic.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"enable_referrer_redirect_on_deny": {
						Optional:    true,
						Description: "When enabled, redirects denied requests rather than responding with an error code.",
						Type:        schema.TypeBool,
					},
					"referrer_redirect_url": {
						Optional:    true,
						Description: "This specifies the full URL to the redirect page for denied requests.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"cors_support": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Cross-origin resource sharing (CORS) allows web pages in one domain to access restricted resources from your domain. Specify external origin hostnames, methods, and headers that you want to accept via HTTP response headers. Full support of CORS requires allowing requests that use the OPTIONS method. See `allowOptions`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables CORS feature.",
						Type:        schema.TypeBool,
					},
					"allow_origins": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ANY", "SPECIFIED"}, false)),
						Optional:         true,
						Description:      "In responses to preflight requests, sets which origin hostnames to accept requests from.",
						Type:             schema.TypeString,
					},
					"origins": {
						Optional:    true,
						Description: "Defines the origin hostnames to accept requests from. The hostnames that you enter need to start with `http` or `https`. For detailed hostname syntax requirements, refer to RFC-952 and RFC-1123 specifications.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"allow_credentials": {
						Optional:    true,
						Description: "Accepts requests made using credentials, like cookies or TLS client certificates.",
						Type:        schema.TypeBool,
					},
					"allow_headers": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ANY", "SPECIFIED"}, false)),
						Optional:         true,
						Description:      "In responses to preflight requests, defines which headers to allow when making the actual request.",
						Type:             schema.TypeString,
					},
					"headers": {
						Optional:    true,
						Description: "Defines the supported request headers.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"methods": {
						Optional:    true,
						Description: "Specifies any combination of the following methods that are allowed when accessing the resource from an external domain: `DELETE`, `GET`, `PATCH`, `POST`, `HEAD`, and `PUT`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"expose_headers": {
						Optional:    true,
						Description: "In responses to preflight requests, lists names of headers that clients can access. By default, clients can access the following simple response headers: `Cache-Control`, `Content-Language`, `Content-Type`, `Expires`, `Last-Modified`, and `Pragma`. You can add other header names to make them accessible to clients.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"preflight_max_age": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Defines the number of seconds that the browser should cache the response to a preflight request.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"cp_code": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Content Provider Codes (CP codes) allow you to distinguish various reporting and billing traffic segments, and you need them to access properties. You receive an initial CP code when purchasing Akamai, and you can run the `Create a new CP code` operation to generate more. This behavior applies any valid CP code, either as required as a default at the top of the rule tree, or subsequently to override the default. For a CP code to be valid, it needs to be assigned the same contract and product as the property, and the group needs access to it.  For available values, run the `List CP codes` operation. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						Optional:    true,
						Description: "Specifies the CP code as an object. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"custom_behavior": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows you to insert a customized XML metadata behavior into any property's rule tree.  Talk to your Akamai representative to implement the customized behavior. Once it's ready, run PAPI's `List custom behaviors` operation, then apply the relevant `behaviorId` value from the response within the current `customBehavior`. See `Custom behaviors and overrides` for guidance on custom metadata behaviors. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior_id": {
						Optional:    true,
						Description: "The unique identifier for the predefined custom behavior you want to insert into the current rule.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"datastream": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `DataStream` reporting service provides real-time logs on application activity, including aggregated metrics on complete request and response cycles and origin response times.  Apply this behavior to report on this set of traffic.  Use the `DataStream API` to aggregate the data. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"stream_type": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validation.ToDiagFunc(validation.StringInSlice([]string{"BEACON", "LOG", "BEACON_AND_LOG"}, false))),
						Optional:         true,
						Description:      "Specify the DataStream type.",
						Type:             schema.TypeString,
					},
					"beacon_stream_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables DataStream reporting.",
						Type:        schema.TypeBool,
					},
					"datastream_ids": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9]+(-[0-9]+)*$")),
						Optional:         true,
						Description:      "A set of dash-separated DataStream ID values to limit the scope of reported data. By default, all active streams report. Use the DataStream application to gather stream ID values that apply to this property configuration. Specifying IDs for any streams that don't apply to this property has no effect, and results in no data reported.",
						Type:             schema.TypeString,
					},
					"log_stream_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"log_enabled": {
						Optional:    true,
						Description: "Enables log collection for the property by associating it with DataStream configurations.",
						Type:        schema.TypeBool,
					},
					"log_stream_name": {
						Optional:    true,
						Description: "Specifies the unique IDs of streams configured for the property. For properties created with the previous version of the rule format, this option contains a string instead of an array of strings. You can use the `List streams` operation to get stream IDs.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"sampling_percentage": {
						Optional:    true,
						Description: "Specifies the percentage of log data you want to collect for this property.",
						Type:        schema.TypeInt,
					},
					"collect_midgress_traffic": {
						Optional:    true,
						Description: "If enabled, gathers midgress traffic data within the Akamai platform, such as between two edge servers, for all streams configured.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"dcp": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. (The `IoT Edge Connect API` allows programmatic access.) This behavior allows you to select previously reserved namespaces and set the protocols for users to publish and receive messages within these namespaces.  Use the `verifyJsonWebTokenForDcp` behavior to control access. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables IoT Edge Connect.",
						Type:        schema.TypeBool,
					},
					"namespace_id": {
						Optional:    true,
						Description: "Specifies the globally reserved name for a specific configuration. It includes authorization rules over publishing and subscribing to logical categories known as `topics`. This provides a root path for all topics defined within a namespace configuration.  You can use the `IoT Edge Connect API` to configure access control lists for your namespace configuration.",
						Type:        schema.TypeString,
					},
					"tlsenabled": {
						Optional:    true,
						Description: "When enabled, you can publish and receive messages over a secured MQTT connection on port 8883.",
						Type:        schema.TypeBool,
					},
					"wsenabled": {
						Optional:    true,
						Description: "When enabled, you can publish and receive messages through a secured MQTT connection over WebSockets on port 443.",
						Type:        schema.TypeBool,
					},
					"gwenabled": {
						Optional:    true,
						Description: "When enabled, you can publish and receive messages over a secured HTTP connection on port 443.",
						Type:        schema.TypeBool,
					},
					"anonymous": {
						Optional:    true,
						Description: "When enabled, you don't need to pass the JWT token with the mqtt request, and JWT validation is skipped.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"dcp_auth_hmac_transformation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. In conjunction with `dcpAuthVariableExtractor`, this behavior affects how clients can authenticate themselves to edge servers, and which groups within namespaces are authorized to access topics. It transforms a source string value extracted from the client certificate and stored as a variable, then generates a hash value based on the selected algorithm, for use in authenticating the client request. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"hash_conversion_algorithm": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SHA256", "MD5", "SHA384"}, false)),
						Optional:         true,
						Description:      "Specifies the hash algorithm.",
						Type:             schema.TypeString,
					},
					"hash_conversion_key": {
						Optional:    true,
						Description: "Specifies the key to generate the hash, ideally a long random string to ensure adequate security.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"dcp_auth_regex_transformation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. In conjunction with `dcpAuthVariableExtractor`, this behavior affects how clients can authenticate themselves to edge servers, and which groups within namespaces are authorized to access topics. It transforms a source string value extracted from the client certificate and stored as a variable, then transforms the string based on a regular expression search pattern, for use in authenticating the client request. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"regex_pattern": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\(\\)]*\\([^\\(\\)]+\\)[^\\(\\)]*$"),
						Optional:         true,
						Description:      "Specifies a Perl-compatible regular expression with a single grouping to capture the text.  For example, a value of `^.(.{0,10})` omits the first character, but then captures up to 10 characters after that. If the regular expression does not capture a substring, authentication may fail.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"dcp_auth_substring_transformation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. In conjunction with `dcpAuthVariableExtractor`, this behavior affects how clients can authenticate themselves to edge servers, and which groups within namespaces are authorized to access topics. It transforms a source string value extracted from the client certificate and stored as a variable, then extracts a substring, for use in authenticating the client request. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"substring_start": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9]+$")),
						Optional:         true,
						Description:      "The zero-based index offset of the first character to extract. If the index is out of bound from the string's length, authentication may fail.",
						Type:             schema.TypeString,
					},
					"substring_end": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9]+$")),
						Optional:         true,
						Description:      "The zero-based index offset of the last character to extract, where `-1` selects the remainder of the string. If the index is out of bound from the string's length, authentication may fail.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"dcp_auth_variable_extractor": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. This behavior affects how clients can authenticate themselves to edge servers, and which groups within namespaces are authorized to access topics. When enabled, this behavior allows end users to authenticate their requests with valid x509 client certificates. Either a client identifier or access authorization groups are required to make the request valid. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"certificate_field": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SUBJECT_DN", "V3_SUBJECT_ALT_NAME", "SERIAL", "FINGERPRINT_DYN", "FINGERPRINT_MD5", "FINGERPRINT_SHA1", "V3_NETSCAPE_COMMENT"}, false)),
						Optional:         true,
						Description:      "Specifies the field in the client certificate to extract the variable from.",
						Type:             schema.TypeString,
					},
					"dcp_mutual_auth_processing_variable_id": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VAR_DCP_CLIENT_ID", "VAR_DCP_AUTH_GROUP"}, false)),
						Optional:         true,
						Description:      "Where to store the value.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"dcp_default_authz_groups": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. This behavior defines a set of default authorization groups to add to each request the property configuration controls.  These groups have access regardless of the authentication method you use, either JWT using the `verifyJsonWebTokenForDcp` behavior, or mutual authentication using the `dcpAuthVariableExtractor` behavior to control where authorization groups are extracted from within certificates. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"group_names": {
						Optional:    true,
						Description: "Specifies the set of authorization groups to assign to all connecting devices.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"dcp_dev_relations": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. This behavior allows Akamai-external clients to use developer test accounts in a shared environment. In conjunction with `verifyJsonWebTokenForDcp`, this behavior allows you to use your own JWTs in your requests, or those generated by Akamai. It lets you either enable the default JWT server for your test configuration by setting the authentication endpoint to a default path, or specify custom settings for your JWT server and the authentication endpoint. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the default JWT server and sets the authentication endpoint to a default path.",
						Type:        schema.TypeBool,
					},
					"custom_values": {
						Optional:    true,
						Description: "Allows you to specify custom JWT server connection values.",
						Type:        schema.TypeBool,
					},
					"hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9]([a-zA-Z0-9_\\-]*[a-zA-Z0-9])?)\\.)+([a-zA-Z]+|xn--[a-zA-Z0-9]+)$"),
						Optional:         true,
						Description:      "Specifies the JWT server's hostname.",
						Type:             schema.TypeString,
					},
					"path": {
						Optional:    true,
						Description: "Specifies the path to your JWT server's authentication endpoint. This lets you generate JWTs to sign your requests.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"dcp_real_time_auth": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "INTERNAL ONLY: The `Internet of Things: Edge Connect` product allows connected users and devices to communicate on a publish-subscribe basis within reserved namespaces. This behavior lets you configure the real time authentication to edge servers. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"extract_namespace": {
						Optional:    true,
						Description: "Extracts a namespace from JSON web tokens (JWT).",
						Type:        schema.TypeBool,
					},
					"namespace_claim": {
						Optional:    true,
						Description: "Specifies the claim in JWT to extract the namespace from.",
						Type:        schema.TypeString,
					},
					"extract_jurisdiction": {
						Optional:    true,
						Description: "Extracts a jurisdiction that defines a geographically distributed set of servers from JWT.",
						Type:        schema.TypeBool,
					},
					"jurisdiction_claim": {
						Optional:    true,
						Description: "Specifies the claim in JWT to extract the jurisdiction from.",
						Type:        schema.TypeString,
					},
					"extract_hostname": {
						Optional:    true,
						Description: "Extracts a hostname from JWT.",
						Type:        schema.TypeBool,
					},
					"hostname_claim": {
						Optional:    true,
						Description: "Specifies the claim in JWT to extract the hostname from.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"delivery_receipt": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "A static behavior that's required when specifying the Cloud Monitor module's (`edgeConnect` behavior. You can only apply this behavior if the property is marked as secure. See `Secure property requirements` for guidance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"deny_access": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Assuming a condition in the rule matches, this denies access to the requested content. For example, a `userLocation` match paired with this behavior would deny requests from a specified part of the world. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"reason": {
						ValidateDiagFunc: validateRegexOrVariable("^[\\w-]+$"),
						Optional:         true,
						Description:      "Text message that keys why access is denied. Any subsequent `denyAccess` behaviors within the rule tree may refer to the same `reason` key to override the current behavior.",
						Type:             schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Denies access when enabled.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"deny_direct_failover_access": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "A static behavior required for all properties that implement a failover under the Cloud Security Failover product. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"device_characteristic_cache_id": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "By default, source URLs serve as cache IDs on edge servers. Electronic Data Capture allows you to specify an additional set of device characteristics to generate separate cache keys. Use this in conjunction with the `deviceCharacteristicHeader` behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"elements": {
						Optional:    true,
						Description: "Specifies a set of information about the device with which to generate a separate cache key.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"device_characteristic_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Sends selected information about requesting devices to the origin server, in the form of an `X-Akamai-Device-Characteristics` HTTP header. Use in conjunction with the `deviceCharacteristicCacheId` behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"elements": {
						Optional:    true,
						Description: "Specifies the set of information about the requesting device to send to the origin server.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"dns_async_refresh": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow an edge server to use an expired DNS record when forwarding a request to your origin. The `type A` DNS record refreshes `after` content is served to the end user, so there is no wait for the DNS resolution. Avoid this behavior if you want to be able to disable a server immediately after its DNS record expires. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows edge servers to refresh an expired DNS record after serving content.",
						Type:        schema.TypeBool,
					},
					"timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Set the maximum allowed time an expired DNS record may be active.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"dns_prefresh": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows edge servers to refresh your origin's DNS record independently from end-user requests. The `type A` DNS record refreshes before the origin's DNS record expires. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows edge servers to refresh DNS records before they expire.",
						Type:        schema.TypeBool,
					},
					"delay": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the amount of time following a DNS record's expiration to asynchronously prefresh it.",
						Type:             schema.TypeString,
					},
					"timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the amount of time to prefresh a DNS entry if there have been no requests to the domain name.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"downgrade_protocol": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Serve static objects to the end-user client over HTTPS, but fetch them from the origin via HTTP. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the protocol downgrading behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"download_complete_marker": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: OTA Updates` product allows customers to securely distribute firmware to devices over cellular networks.  Based on match criteria that executes a rule, this behavior logs requests to the OTA servers as completed in aggregated and individual reports. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"download_notification": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: OTA Updates` product allows customers to securely distribute firmware to devices over cellular networks. Based on match criteria that executes a rule, this behavior allows requests to the `OTA Updates API` for a list of completed downloads to individual vehicles. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"downstream_cache": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specify the caching instructions the edge server sends to the end user's client or client proxies. By default, the cache's duration is whichever is less: the remaining lifetime of the edge cache, or what the origin's header specifies. If the origin is set to `no-store` or `bypass-cache`, edge servers send `cache-busting` headers downstream to prevent downstream caching. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "MUST_REVALIDATE", "BUST", "TUNNEL_ORIGIN", "NONE"}, false)),
						Optional:         true,
						Description:      "Specify the caching instructions the edge server sends to the end user's client.",
						Type:             schema.TypeString,
					},
					"allow_behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESSER", "GREATER", "REMAINING_LIFETIME", "FROM_MAX_AGE", "FROM_VALUE", "PASS_ORIGIN"}, false)),
						Optional:         true,
						Description:      "Specify how the edge server calculates the downstream cache by setting the value of the `Expires` header.",
						Type:             schema.TypeString,
					},
					"ttl": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Sets the duration of the cache. Setting the value to `0` equates to a `no-cache` header that forces revalidation.",
						Type:             schema.TypeString,
					},
					"send_headers": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL_AND_EXPIRES", "CACHE_CONTROL", "EXPIRES", "PASS_ORIGIN"}, false)),
						Optional:         true,
						Description:      "Specifies the HTTP headers to include in the response to the client.",
						Type:             schema.TypeString,
					},
					"send_private": {
						Optional:    true,
						Description: "Adds a `Cache-Control: private` header to prevent objects from being cached in a shared caching proxy.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"dynamic_throughtput_optimization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables `quick retry`, which detects slow forward throughput while fetching an object, and attempts a different forward connection path to avoid congestion. By default, connections under 5 mbps trigger this behavior. When the transfer rate drops below this rate during a connection attempt, quick retry is enabled and a different forward connection path is used. Contact Akamai Professional Services to override this threshold. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the quick retry feature.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"dynamic_throughtput_optimization_override": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This overrides the default threshold of 5 Mbps that triggers the `dynamicThroughtputOptimization` behavior, which enables the quick retry feature. Quick retry detects slow forward throughput while fetching an object, and attempts a different forward connection path to avoid congestion. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"throughput": {
						Optional:    true,
						Description: "Specifies the default target forward throughput in Mbps, ranging from 2 to 50 Mbps. If this time is exceeded during a connection attempt, quick retry is enabled and a different forward connection path is used.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"dynamic_web_content": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "In conjunction with the `subCustomer` behavior, this optional behavior allows you to control how dynamic web content behaves for your subcustomers using `Akamai Cloud Embed`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"sure_route": {
						Optional:    true,
						Description: "Optimizes how subcustomer traffic routes from origin to edge servers.  See the `sureRoute` behavior for more information.",
						Type:        schema.TypeBool,
					},
					"prefetch": {
						Optional:    true,
						Description: "Allows subcustomer content to prefetch over HTTP/2.",
						Type:        schema.TypeBool,
					},
					"real_user_monitoring": {
						Optional:    true,
						Description: "Allows Real User Monitoring (RUM) to collect performance data for subcustomer content. See the `realUserMonitoring` behavior for more information.",
						Type:        schema.TypeBool,
					},
					"image_compression": {
						Optional:    true,
						Description: "Enables image compression for subcustomer content.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"early_data": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Use this behavior to enable sending early data during the TLS 1.3 handshake between requests from your client and Akamai edge servers. This is available for QUIC connections and Transmission Control Protocol (TCP). This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables sending early data, which further reduces latency in TLS 1.3 connections.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"early_hints": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Use Early Hints to send an HTTP 103 status code with preliminary HTTP headers at the client request stage, so that a browser can preload critical website resources or preconnect to a specific domain while waiting for the final response. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enable the behavior so that browsers can use that waiting time to preload the resource URLs you specify or preconnect to static or image domains.",
						Type:        schema.TypeBool,
					},
					"resource_url": {
						Optional:    true,
						Description: "Enter the URL to a resource you want clients to receive as an early hint. Edge servers include each resource URL you provide in an instance of the `Link` header that's sent back to the client in the HTTP 103 response. You only need to specify the value of the header, as edge servers automatically add the `Link` header name to the response. Use commas to separate multiple entries. This field supports variables and string concatenation. The URL must be enclosed between `<` and `>` as shown in the example below.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"ecms_bulk_upload": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Uploads a ZIP archive with objects to an existing data set. The target data set stores objects as key-value pairs. The path to an object in the ZIP archive is a key, and the content of an object is a value. For an overview, see `ecmsDatabase`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables sending a compressed archive file with objects. Sends the archive file to the default path of the target data set: `<hostname>/bulk/<database_name>/<dataset_name>`.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"ecms_database": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Edge Connect Message Store is available for `Internet of Things: Edge Connect` users. It lets you create databases and data sets within these databases. You can use this object store to save files smaller than 2 GB. `ecmsDatabase` specifies a default database for requests to this property, unless indicated otherwise in the URL. To access objects in the default database, you can skip its name in the URLs. To access objects in a different database, pass its name in the header, query parameter, or a regular expression matching a URL segment. You can also configure the `ecmsDataset` behavior to specify a default data set for requests. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"database": {
						Optional:    true,
						Description: "Specifies a default database for this property. If you don't configure a default data set in the `ecmsDataset` behavior, requests to objects in this database follow the pattern: `<hostname>/datastore/<data_set_name>/<object_key>`.",
						Type:        schema.TypeString,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING", "REGEX"}, false)),
						Optional:         true,
						Description:      "Specifies where to pass a database name in requests. If the specified location doesn't include the database name or the name doesn't match the regular expression, the default database is used.",
						Type:             schema.TypeString,
					},
					"header_name": {
						Optional:    true,
						Description: "Specifies the request header that passed the database name. By default, it points to `X-KV-Database`.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "Specifies the query string parameter that passed the database name. By default, it points to `database`.",
						Type:        schema.TypeString,
					},
					"regex_pattern": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\(\\)]*\\([^\\(\\)]+\\)[^\\(\\)]*$"),
						Optional:         true,
						Description:      "Specifies the regular expression that matches the database name in the URL.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"ecms_dataset": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies a default data set for requests to this property unless indicated otherwise in the URL. To access objects in this data set, you can skip the data set name in the URLs. To access objects in a different data set within a database, pass the data set name in the header, query parameter, or a regular expression pattern matching a URL segment. You can also configure the `ecmsDatabase` behavior to specify a default database for requests. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"dataset": {
						Optional:    true,
						Description: "Specifies a default data set for this property. If you don't configure a default database in the `ecmsDatabase` behavior, requests to objects in this data set follow the pattern: `<hostname>/datastore/<database_name>/<object_key>`.",
						Type:        schema.TypeString,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING", "REGEX"}, false)),
						Optional:         true,
						Description:      "Specifies where to pass a data set name in requests. If the specified location doesn't include the data set name or the name doesn't match the regular expression pattern, the default data set is used.",
						Type:             schema.TypeString,
					},
					"header_name": {
						Optional:    true,
						Description: "Specifies the request header that passed the data set name. By default, it points to `X-KV-Dataset`.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "Specifies the query string parameter that passed the data set name. By default, it points to `dataset`.",
						Type:        schema.TypeString,
					},
					"regex_pattern": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\(\\)]*\\([^\\(\\)]+\\)[^\\(\\)]*$"),
						Optional:         true,
						Description:      "Specifies the regular expression that matches the data set name in the URL.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"ecms_object_key": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Defines a regular expression to match object keys in custom URLs and to access objects in a data set. You can point custom URLs to access proper values in the target data set. For an overview, see `ecmsDatabase`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"regex": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\(\\)]*\\([^\\(\\)]+\\)[^\\(\\)]*$"),
						Optional:         true,
						Description:      "Enables sending a compressed archive file with objects to the default path of the target data set: `<hostname>/bulk/<database_name>/<dataset_name>`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"edge_connect": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Configures traffic logs for the Cloud Monitor push API. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables Cloud Monitor's log-publishing behavior.",
						Type:        schema.TypeBool,
					},
					"api_connector": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DEFAULT", "SIEM_JSON", "BMC_APM"}, false)),
						Optional:         true,
						Description:      "Describes the API connector type.",
						Type:             schema.TypeString,
					},
					"api_data_elements": {
						Optional:    true,
						Description: "Specifies the data set to log.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the target hostname accepting push API requests.",
						Type:             schema.TypeString,
					},
					"destination_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the push API's endpoint.",
						Type:             schema.TypeString,
					},
					"override_aggregate_settings": {
						Optional:    true,
						Description: "When enabled, overrides default log settings.",
						Type:        schema.TypeBool,
					},
					"aggregate_time": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies how often logs are generated.",
						Type:             schema.TypeString,
					},
					"aggregate_lines": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "Specifies the maximum number of lines to include in each log.",
						Type:             schema.TypeString,
					},
					"aggregate_size": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d+[K,M,G,T]B$"),
						Optional:         true,
						Description:      "Specifies the log's maximum size.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"edge_load_balancing_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior implements customized Edge Load Balancing features. Contact Akamai Professional Services for help configuring it. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "A description of what the `xml` block does.",
						Type:        schema.TypeString,
					},
					"xml": {
						Optional:    true,
						Description: "A block of Akamai XML metadata.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"edge_load_balancing_data_center": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Edge Load Balancing module allows you to specify groups of data centers that implement load balancing, session persistence, and real-time dynamic failover. Enabling ELB routes requests contextually based on location, device, or network, along with optional rules you specify. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origin_id": {
						Optional:    true,
						Description: "Corresponds to the `id` specified by the `edgeLoadBalancingOrigin` behavior associated with this data center.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "Provides a description for the ELB data center, for your own reference.",
						Type:        schema.TypeString,
					},
					"hostname": {
						ValidateDiagFunc: validateAny(validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"), validateRegexOrVariable("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")),
						Optional:         true,
						Description:      "Specifies the data center's hostname.",
						Type:             schema.TypeString,
					},
					"cookie_name": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[^\\s;]+$")),
						Optional:         true,
						Description:      "If using session persistence, this specifies the value of the cookie named in the corresponding `edgeLoadBalancingOrigin` behavior's `cookie_name` option.",
						Type:             schema.TypeString,
					},
					"failover_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_failover": {
						Optional:    true,
						Description: "Allows you to specify failover rules.",
						Type:        schema.TypeBool,
					},
					"ip": {
						ValidateDiagFunc: validateRegexOrVariable("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"),
						Optional:         true,
						Description:      "Specifies this data center's IP address.",
						Type:             schema.TypeString,
					},
					"failover_rules": {
						Optional:    true,
						Description: "Provides up to four failover rules to apply in the specified order.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"failover_hostname": {
									ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
									Optional:         true,
									Description:      "The hostname of the data center to fail over to.",
									Type:             schema.TypeString,
								},
								"modify_request": {
									Optional:    true,
									Description: "Allows you to modify the request's hostname or path.",
									Type:        schema.TypeBool,
								},
								"override_hostname": {
									Optional:    true,
									Description: "Overrides the request's hostname with the `failover_hostname`.",
									Type:        schema.TypeBool,
								},
								"context_root": {
									ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
									Optional:         true,
									Description:      "Specifies the path to use in the forwarding request, typically the root (`/`) when failing over to a different data center, or a full path such as `/static/error.html` when failing over to an error page.",
									Type:             schema.TypeString,
								},
								"absolute_path": {
									Optional:    true,
									Description: "When enabled, interprets the path specified by `context_root` as an absolute server path, for example to reference a site-down page. Otherwise when disabled, the path is appended to the request.",
									Type:        schema.TypeBool,
								},
							},
						},
					},
				},
			},
		},
		"edge_load_balancing_origin": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Edge Load Balancing module allows you to implement groups of data centers featuring load balancing, session persistence, and real-time dynamic failover. Enabling ELB routes requests contextually based on location, device, or network, along with optional rules you specify. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"id": {
						Optional:    true,
						Description: "Specifies a unique descriptive string for this ELB origin. The value needs to match the `origin_id` specified by the `edgeLoadBalancingDataCenter` behavior associated with this origin.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "Provides a description for the ELB origin, for your own reference.",
						Type:        schema.TypeString,
					},
					"hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the hostname associated with the ELB rule.",
						Type:             schema.TypeString,
					},
					"session_persistence_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_session_persistence": {
						Optional:    true,
						Description: "Allows you to specify a cookie to pin the user's browser session to one data center. When disabled, ELB's default load balancing may send users to various data centers within the same session.",
						Type:        schema.TypeBool,
					},
					"cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This specifies the name of the cookie that marks users' persistent sessions. The accompanying `edgeLoadBalancingDataCenter` behavior's `description` option specifies the cookie's value.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"edge_origin_authorization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows the origin server to use a cookie to ensure requests from Akamai servers are genuine. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the cookie-authorization behavior.",
						Type:        schema.TypeBool,
					},
					"cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name of the cookie to use for authorization.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\s;]+$"),
						Optional:         true,
						Description:      "Specifies the value of the authorization cookie.",
						Type:             schema.TypeString,
					},
					"domain": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specify the cookie's domain, which needs to match the top-level domain of the `Host` header the origin server receives.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"edge_redirector": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior enables the `Edge Redirector Cloudlet` application, which helps you manage large numbers of redirects. With Cloudlets available on your contract, choose `Your services` > `Edge logic Cloudlets` to control the Edge Redirector within `Control Center`. Otherwise use the `Cloudlets API` to configure it programmatically. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Edge Redirector Cloudlet.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to apply the Cloudlet shared policy to an unlimited number of properties within your account. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Specifies the Cloudlet policy as an object.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet shared policy to use with this behavior. Use the `Cloudlets API` to list available shared policies.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"edge_scape": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "`EdgeScape` allows you to customize content based on the end user's geographic location or connection speed. When enabled, the edge server sends a special `X-Akamai-Edgescape` header to the origin server encoding relevant details about the end-user client as key-value pairs. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, sends the `X-Akamai-Edgescape` request header to the origin.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"edge_side_includes": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows edge servers to process edge side include (ESI) code to generate dynamic content. To apply this behavior, you need to match on a `contentType`, `path`, or `filename`. Since this behavior requires more parsing time, you should not apply it to pages that lack ESI code, or to any non-HTML content. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables ESI processing.",
						Type:        schema.TypeBool,
					},
					"enable_via_http": {
						Optional:    true,
						Description: "Enable ESI only for content featuring the `Edge-control: dca=esi` HTTP response header.",
						Type:        schema.TypeBool,
					},
					"pass_set_cookie": {
						Optional:    true,
						Description: "Allows edge servers to pass your origin server's cookies to the ESI processor.",
						Type:        schema.TypeBool,
					},
					"pass_client_ip": {
						Optional:    true,
						Description: "Allows edge servers to pass the client IP header to the ESI processor.",
						Type:        schema.TypeBool,
					},
					"i18n_status": {
						Optional:    true,
						Description: "Provides internationalization support for ESI.",
						Type:        schema.TypeBool,
					},
					"i18n_charset": {
						Optional:    true,
						Description: "Specifies the character sets to use when transcoding the ESI language, `UTF-8` and `ISO-8859-1` for example.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"detect_injection": {
						Optional:    true,
						Description: "Denies attempts to inject ESI code.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"edge_worker": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "`EdgeWorkers` are JavaScript applications that allow you to manipulate your web traffic on edge servers outside of Property Manager behaviors, and deployed independently from your configuration's logic.  This behavior applies an EdgeWorker to a set of edge requests. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, applies specified EdgeWorker functionality to this rule's web traffic.",
						Type:        schema.TypeBool,
					},
					"create_edge_worker": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"edge_worker_id": {
						Optional:    true,
						Description: "Identifies the EdgeWorker application to apply to this rule's web traffic. You can use the `EdgeWorkers API` to get this value.",
						Type:        schema.TypeString,
					},
					"resource_tier": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"m_pulse": {
						Optional:    true,
						Description: "Enables mPulse reports that include data about EdgeWorkers errors generated due to JavaScript errors. For more details, see `Integrate mPulse reports with EdgeWorkers`.",
						Type:        schema.TypeBool,
					},
					"m_pulse_information": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"enforce_mtls_settings": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior repeats mTLS validation checks between a requesting client and the edge network. If the checks fail, you can deny the request or apply custom error handling. To use this behavior, you need to add either the `hostname` or `clientCertificate` criteria to the same rule. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable_auth_set": {
						Optional:    true,
						Description: "Whether to require a specific mutual transport layer security (mTLS) certificate authority (CA) set in a request from a client to the edge network.",
						Type:        schema.TypeBool,
					},
					"certificate_authority_set": {
						Optional:    true,
						Description: "Specify the client certificate authority (CA) sets you want to support in client requests. Run the `List CA Sets` operation in the mTLS Edge TrustStore API to get the `setId` value and pass it in this option as a string. If a request includes a set not defined here, it will be denied. The preset list items you can select are contingent on the CA sets you've created using the mTLS Edge Truststore, and then associated with a certificate in the `Certificate Provisioning System`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"enable_ocsp_status": {
						Optional:    true,
						Description: "Whether the mutual transport layer security requests from a client should use the online certificate support protocol (OCSP). OCSP can determine the x.509 certificate revocation status during the TLS handshake.",
						Type:        schema.TypeBool,
					},
					"enable_deny_request": {
						Optional:    true,
						Description: "This denies a request from a client that doesn't match what you've set for the options in this behavior. When disabled, non-matching requests are allowed, but you can incorporate a custom handling operation, such as reviewing generated log entries to see the discrepancies, enable the `Client-To-Edge` authentication header, or issue a custom message. This behavior only checks the `Certificate Provisioning System` settings. It doesn't check the current client certificate and doesn't deny invalid certs.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"enhanced_akamai_protocol": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the Enhanced Akamai Protocol, a suite of advanced routing and transport optimizations that increase your website's performance and reliability. It is only available to specific applications, and requires a special routing from edge to origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"display": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"enhanced_debug": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior, available by default for all products, provides support for enhanced debugging on edge servers. It includes all the functionality provided by the existing `Pragma` header debugging, but is more secure and provides additional information. All requests that use this behavior pass an auth token that you generate using a secret debug key in the `Akamai-Debug` request header. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable_debug": {
						Optional:    true,
						Description: "Enables enhanced debugging using the `Akamai-Debug` request header.",
						Type:        schema.TypeBool,
					},
					"debug_key": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9a-fA-F]{64}$"),
						Optional:         true,
						Description:      "Specifies the debug key to use for all requests processed by this property. The debug key value needs to be a 64-byte hex string. You can generate the key in one property and then reuse it in other configurations.",
						Type:             schema.TypeString,
					},
					"disable_pragma": {
						Optional:    true,
						Description: "Whether you want to disable the standard debugging that uses the `Pragma` request header.",
						Type:        schema.TypeBool,
					},
					"generate_grn": {
						Optional:    true,
						Description: "Whether you want to return the Global Request Number (GRN) in the `Akamai-GRN` response header for all requests, even if the `Akamai-Debug` request header is not passed. The `Akamai-GRN` header is useful for log extraction.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"enhanced_proxy_detection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enhanced Proxy Detection (EPD) leverages the GeoGuard service provided by GeoComply to add proxy detection and location spoofing protection. It identifies requests for your content that have been redirected from an unwanted source through a proxy. You can then allow, deny, or redirect these requests. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Applies GeoGuard proxy detection.",
						Type:        schema.TypeBool,
					},
					"forward_header_enrichment": {
						Optional:    true,
						Description: "Whether the Enhanced Proxy Detection (Akamai-EPD) header is included in the forward request to mark a connecting IP address as an anonymous proxy, with a two-letter or three-letter designation. See the `epdForwardHeaderEnrichment` behavior for details.",
						Type:        schema.TypeBool,
					},
					"enable_configuration_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BEST_PRACTICE", "ADVANCED"}, false)),
						Optional:         true,
						Description:      "Specifies how to field the proxy request.",
						Type:             schema.TypeString,
					},
					"best_practice_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "Specifies how to field the proxy request.",
						Type:             schema.TypeString,
					},
					"best_practice_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect requests.",
						Type:             schema.TypeString,
					},
					"anonymous_vpn": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_anonymous_vpn": {
						Optional:    true,
						Description: "This detects requests from anonymous VPNs.",
						Type:        schema.TypeBool,
					},
					"detect_anonymous_vpn_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "Specifies how to field anonymous VPN requests.",
						Type:             schema.TypeString,
					},
					"detect_anonymous_vpn_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect anonymous VPN requests.",
						Type:             schema.TypeString,
					},
					"public_proxy": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_public_proxy": {
						Optional:    true,
						Description: "This detects requests from public proxies.",
						Type:        schema.TypeBool,
					},
					"detect_public_proxy_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "Specifies how to field public proxy requests.",
						Type:             schema.TypeString,
					},
					"detect_public_proxy_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect public proxy requests.",
						Type:             schema.TypeString,
					},
					"tor_exit_node": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_tor_exit_node": {
						Optional:    true,
						Description: "This detects requests from Tor exit nodes.",
						Type:        schema.TypeBool,
					},
					"detect_tor_exit_node_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "This specifies whether to `DENY`, `ALLOW`, or `REDIRECT` requests from Tor exit nodes.",
						Type:             schema.TypeString,
					},
					"detect_tor_exit_node_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect requests from Tor exit nodes.",
						Type:             schema.TypeString,
					},
					"smart_dns_proxy": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_smart_dns_proxy": {
						Optional:    true,
						Description: "This detects requests from smart DNS proxies.",
						Type:        schema.TypeBool,
					},
					"detect_smart_dns_proxy_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "Specifies whether to `DENY`, `ALLOW`, or `REDIRECT` smart DNS proxy requests.",
						Type:             schema.TypeString,
					},
					"detect_smart_dns_proxy_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect DNS proxy requests.",
						Type:             schema.TypeString,
					},
					"hosting_provider": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_hosting_provider": {
						Optional:    true,
						Description: "This detects requests from a hosting provider.",
						Type:        schema.TypeBool,
					},
					"detect_hosting_provider_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "This specifies whether to `DENY`, `ALLOW`, or `REDIRECT` requests from hosting providers.",
						Type:             schema.TypeString,
					},
					"detect_hosting_provider_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the absolute URL to which to redirect requests from hosting providers.",
						Type:             schema.TypeString,
					},
					"vpn_data_center": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_vpn_data_center": {
						Optional:    true,
						Description: "This detects requests from VPN data centers.",
						Type:        schema.TypeBool,
					},
					"detect_vpn_data_center_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "This specifies whether to `DENY`, `ALLOW`, or `REDIRECT` requests from VPN data centers.",
						Type:             schema.TypeString,
					},
					"detect_vpn_data_center_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect requests from VPN data centers.",
						Type:             schema.TypeString,
					},
					"residential_proxy": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_residential_proxy": {
						Optional:    true,
						Description: "This detects requests from a residential proxy. See `Enhanced Proxy Detection with GeoGuard` and learn more about this GeoGuard category before enabling it.",
						Type:        schema.TypeBool,
					},
					"detect_residential_proxy_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALLOW", "DENY", "REDIRECT"}, false)),
						Optional:         true,
						Description:      "This specifies whether to `DENY`, `ALLOW`, or `REDIRECT` requests from residential proxies.",
						Type:             schema.TypeString,
					},
					"detect_residential_proxy_redirecturl": {
						ValidateDiagFunc: validateRegexOrVariable("(http|https)://(\\w+:{0,1}\\w*@)?(\\S+)(:[0-9]+)?(/|/([\\w#!:.?+=&%@!\\-/]))?"),
						Optional:         true,
						Description:      "This specifies the URL to which to redirect requests.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"epd_forward_header_enrichment": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior identifies unwanted requests from an anonymous proxy. This and the `enhancedProxyDetection` behavior work together and need to be included either in the same rule, or in the default one. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Sends the Enhanced Proxy Detection (`Akamai-EPD`) header in the forward request to determine whether the connecting IP address is an anonymous proxy. The header can contain one or more codes that indicate the IP address type detected by edge servers:",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"fail_action": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies how to respond when the origin is not available: by serving stale content, by serving an error page, or by redirecting.  To apply this behavior, you should match on an `originTimeout` or `matchResponseCode`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled in case of a failure to contact the origin, the current behavior applies.",
						Type:        schema.TypeBool,
					},
					"action_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SERVE_STALE", "REDIRECT", "RECREATED_CO", "RECREATED_CEX", "RECREATED_NS", "DYNAMIC"}, false)),
						Optional:         true,
						Description:      "Specifies the basic action to take when there is a failure to contact the origin.",
						Type:             schema.TypeString,
					},
					"saas_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HOSTNAME", "PATH", "QUERY_STRING", "COOKIE"}, false)),
						Optional:         true,
						Description:      "Identifies the component of the request that identifies the SaaS dynamic fail action.",
						Type:             schema.TypeString,
					},
					"saas_cname_enabled": {
						Optional:    true,
						Description: "Specifies whether to use a CNAME chain to determine the hostname for the SaaS dynamic failaction.",
						Type:        schema.TypeBool,
					},
					"saas_cname_level": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the number of elements in the CNAME chain backwards from the edge hostname that determines the hostname for the SaaS dynamic failaction.",
						Type:             schema.TypeInt,
					},
					"saas_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name of the cookie that identifies this SaaS dynamic failaction.",
						Type:             schema.TypeString,
					},
					"saas_query_string": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "Specifies the name of the query parameter that identifies this SaaS dynamic failaction.",
						Type:             schema.TypeString,
					},
					"saas_regex": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9\\:\\[\\]\\{\\}\\(\\)\\.\\?_\\-\\*\\+\\^\\$\\\\\\/\\|&=!]{1,250})$"),
						Optional:         true,
						Description:      "Specifies the substitution pattern (a Perl-compatible regular expression) that defines the SaaS dynamic failaction.",
						Type:             schema.TypeString,
					},
					"saas_replace": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9]|\\$[1-9])(([a-zA-Z0-9\\._\\-]|\\$[1-9]){0,250}([a-zA-Z0-9]|\\$[1-9]))?){1,10}$"),
						Optional:         true,
						Description:      "Specifies the replacement pattern that defines the SaaS dynamic failaction.",
						Type:             schema.TypeString,
					},
					"saas_suffix": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})\\.(com|net|org|info|biz|us|co\\.uk|ac\\.uk|org\\.uk|me\\.uk|ca|eu|com\\.au|co|co\\.za|ru|es|me|tv|pro|in|ie|de|it|nl|fr|co\\.il|ch|se|co\\.nz|pl|jp|name|mobi|cc|ws|be|com\\.mx|at|nu|asia|co\\.nz|net\\.nz|org\\.nz|com\\.au|net\\.au|org\\.au|tools)$"),
						Optional:         true,
						Description:      "Specifies the static portion of the SaaS dynamic failaction.",
						Type:             schema.TypeString,
					},
					"dynamic_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SERVE_301", "SERVE_302", "SERVE_ALTERNATE"}, false)),
						Optional:         true,
						Description:      "Specifies the redirect method.",
						Type:             schema.TypeString,
					},
					"dynamic_custom_path": {
						Optional:    true,
						Description: "Allows you to modify the original requested path.",
						Type:        schema.TypeBool,
					},
					"dynamic_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the new path.",
						Type:             schema.TypeString,
					},
					"redirect_hostname_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ORIGINAL", "ALTERNATE"}, false)),
						Optional:         true,
						Description:      "Whether to preserve or customize the hostname.",
						Type:             schema.TypeString,
					},
					"redirect_hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "When the `actionType` is `REDIRECT` and the `redirectHostnameType` is `ALTERNATE`, this specifies the hostname for the redirect.",
						Type:             schema.TypeString,
					},
					"redirect_custom_path": {
						Optional:    true,
						Description: "Uses the `redirectPath` to customize a new path.",
						Type:        schema.TypeBool,
					},
					"redirect_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies a new path.",
						Type:             schema.TypeString,
					},
					"redirect_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{302, 301})),
						Optional:         true,
						Description:      "Specifies the HTTP response code.",
						Type:             schema.TypeInt,
					},
					"content_hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the static hostname for the alternate redirect.",
						Type:             schema.TypeString,
					},
					"content_custom_path": {
						Optional:    true,
						Description: "Specifies a custom redirect path.",
						Type:        schema.TypeBool,
					},
					"content_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies a custom redirect path.",
						Type:             schema.TypeString,
					},
					"net_storage_hostname": {
						Optional:    true,
						Description: "When the `actionType` is `RECREATED_NS`, specifies the `NetStorage` origin to serve the alternate content. Contact Akamai Professional Services for your NetStorage origin's `id`.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"net_storage_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "When the `actionType` is `RECREATED_NS`, specifies the path for the `NetStorage` request.",
						Type:             schema.TypeString,
					},
					"cex_hostname": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies a hostname.",
						Type:             schema.TypeString,
					},
					"cex_custom_path": {
						Optional:    true,
						Description: "Specifies a custom path.",
						Type:        schema.TypeBool,
					},
					"cex_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies a custom path.",
						Type:             schema.TypeString,
					},
					"cp_code": {
						Optional:    true,
						Description: "Specifies a CP code for which to log errors for the NetStorage location. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"status_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{200, 404, 500, 100, 101, 102, 103, 122, 201, 202, 203, 204, 205, 206, 207, 226, 400, 401, 402, 403, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 422, 423, 424, 425, 426, 428, 429, 431, 444, 449, 450, 499, 501, 502, 503, 504, 505, 506, 507, 509, 510, 511, 598, 599})),
						Optional:         true,
						Description:      "Assigns a new HTTP status code to the failure response.",
						Type:             schema.TypeInt,
					},
					"preserve_query_string": {
						Optional:    true,
						Description: "When using either `contentCustomPath`, `cexCustomPath`, `dynamicCustomPath`, or `redirectCustomPath` to specify a custom path, enabling this passes in the original request's query string as part of the path.",
						Type:        schema.TypeBool,
					},
					"modify_protocol": {
						Optional:    true,
						Description: "Modifies the redirect's protocol using the value of the `protocol` field.",
						Type:        schema.TypeBool,
					},
					"protocol": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "When the `actionType` is `REDIRECT` and `modifyProtocol` is enabled, this specifies the redirect's protocol.",
						Type:             schema.TypeString,
					},
					"allow_fcm_parent_override": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"failover_bot_manager_feature_compatibility": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Ensures that functionality such as challenge authentication and reset protocol work with a failover product property you use to create an alternate hostname. Apply it to any properties that implement a failover under the Cloud Security Failover product. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"compatibility": {
						Optional:    true,
						Description: "This behavior does not include any options. Specifying the behavior itself enables it.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"fast_invalidate": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, forces a validation test for all edge content to which the behavior applies.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"fips": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Ensures `Federal Information Process Standards (FIPS) 140-2` compliance for a connection to an origin server. For this behavior to work properly, verify that your origin's secure certificate supports Enhanced TLS and is FIPS-compliant. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "When enabled, supports the use of FIPS-validated ciphers in the connection between this delivery configuration and your origin server.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"first_party_marketing": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the Cloud Marketing Cloudlet, which helps MediaMath customers collect usage data and place corresponding tags for use in online advertising.  You can configure tags using either the Cloudlets Policy Manager application or the `Cloudlets API`. See also the `firstPartyMarketingPlus` behavior, which integrates better with both MediaMath and its partners. Both behaviors support the same set of options. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Cloud Marketing Cloudlet.",
						Type:        schema.TypeBool,
					},
					"java_script_insertion_rule": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NEVER", "POLICY", "ALWAYS"}, false)),
						Optional:         true,
						Description:      "Select how to insert the MediaMath JavaScript reference script.",
						Type:             schema.TypeString,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"media_math_prefix": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specify the URL path prefix that distinguishes Cloud Marketing requests from your other web traffic. Include the leading slash character, but no trailing slash.  For example, if the path prefix is `/mmath`, and the request is for `www.example.com/dir`, the new URL is `www.example.com/mmath/dir`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"first_party_marketing_plus": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the Cloud Marketing Plus Cloudlet, which helps MediaMath customers collect usage data and place corresponding tags for use in online advertising.  You can configure tags using either the Cloudlets Policy Manager application or the `Cloudlets API`. See also the `firstPartyMarketing` behavior, which integrates with MediaMath but not its partners. Both behaviors support the same set of options. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Cloud Marketing Plus Cloudlet.",
						Type:        schema.TypeBool,
					},
					"java_script_insertion_rule": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NEVER", "POLICY", "ALWAYS"}, false)),
						Optional:         true,
						Description:      "Select how to insert the MediaMath JavaScript reference script.",
						Type:             schema.TypeString,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"media_math_prefix": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specify the URL path prefix that distinguishes Cloud Marketing requests from your other web traffic. Include the leading slash character, but no trailing slash.  For example, if the path prefix is `/mmath`, and the request is for `www.example.com/dir`, the new URL is `www.example.com/mmath/dir`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"forward_rewrite": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Forward Rewrite Cloudlet allows you to conditionally modify the forward path in edge content without affecting the URL that displays in the user's address bar. If Cloudlets are available on your contract, choose `Your services` > `Edge logic Cloudlets` to control how this feature works within `Control Center`, or use the `Cloudlets API` to configure it programmatically. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Forward Rewrite Cloudlet behavior.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to use a shared policy for a Cloudlet. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "This identifies the Cloudlet shared policy to use with this behavior. You can list available shared policies with the `Cloudlets API`.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"g2oheader": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `signature header authentication` (g2o) security feature provides header-based verification of outgoing origin requests. Edge servers encrypt request data in a pre-defined header, which the origin uses to verify that the edge server processed the request. This behavior configures the request data, header names, encryption algorithm, and shared secret to use for verification. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the g2o verification behavior.",
						Type:        schema.TypeBool,
					},
					"data_header": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies the name of the header that contains the request data that needs to be encrypted.",
						Type:             schema.TypeString,
					},
					"signed_header": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies the name of the header containing encrypted request data.",
						Type:             schema.TypeString,
					},
					"encoding_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{1, 2, 3, 4, 5})),
						Optional:         true,
						Description:      "Specifies the version of the encryption algorithm as an integer from `1` through `5`.",
						Type:             schema.TypeInt,
					},
					"use_custom_sign_string": {
						Optional:    true,
						Description: "When disabled, the encrypted string is based on the forwarded URL. If enabled, you can use `customSignString` to customize the set of data to encrypt.",
						Type:        schema.TypeBool,
					},
					"custom_sign_string": {
						Optional:    true,
						Description: "Specifies the set of data to be encrypted as a combination of concatenated strings.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"secret_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9a-zA-Z]{24}$")),
						Optional:         true,
						Description:      "Specifies the shared secret key.",
						Type:             schema.TypeString,
					},
					"nonce": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9a-zA-Z]{1,8}$"),
						Optional:         true,
						Description:      "Specifies the cryptographic `nonce` string.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"global_request_number": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Generates a unique identifier for each request on the Akamai edge network, for use in logging and debugging. GRN identifiers follow the same format as Akamai's error reference strings, for example: `0.05313217.1567801841.1457a3`. You can use the Edge Diagnostics API's `Translate error string` operation to get low-level details about any request. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"output_option": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RESPONSE_HEADER", "REQUEST_HEADER", "BOTH_HEADERS", "ASSIGN_VARIABLE"}, false)),
						Optional:         true,
						Description:      "Specifies how to report the GRN value.",
						Type:             schema.TypeString,
					},
					"header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "With `outputOption` set to specify any set of headers, this specifies the name of the header to report the GRN value.",
						Type:             schema.TypeString,
					},
					"variable_name": {
						Optional:    true,
						Description: "This specifies the name of the variable to assign the GRN value to. You need to pre-declare any `variable` you specify within the rule tree.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"gov_cloud": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"gov_cloud_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"graphql_caching": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior configures how to cache GraphQL-based API traffic.  Enable `caching` for your GraphQL API traffic, along with `allowPost` to cache POST responses.  To configure REST API traffic, use the `rapid` behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables GraphQL caching.",
						Type:        schema.TypeBool,
					},
					"cache_responses_with_errors": {
						Optional:    true,
						Description: "When enabled, caches responses that include an `error` field at the top of the response body object.  Disable this if your GraphQL server yields temporary errors with success codes in the 2xx range.",
						Type:        schema.TypeBool,
					},
					"advanced": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"post_request_processing_error_handling": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"APPLY_CACHING_BEHAVIOR", "NO_STORE"}, false)),
						Optional:         true,
						Description:      "Specify what happens if GraphQL query processing fails on POST requests.",
						Type:             schema.TypeString,
					},
					"operations_url_query_parameter_name": {
						Optional:    true,
						Description: "Specifies the name of a query parameter that identifies requests as GraphQL queries.",
						Type:        schema.TypeString,
					},
					"operations_json_body_parameter_name": {
						Optional:    true,
						Description: "The name of the JSON body parameter that identifies GraphQL POST requests.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"gzip_response": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Apply `gzip` compression to speed transfer time. This behavior applies best to text-based content such as HTML, CSS, and JavaScript, especially once files exceed about 10KB. Do not apply it to already compressed image formats, or to small files that would add more time to uncompress. To apply this behavior, you should match on `contentType` or the content's `cacheability`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ORIGIN_RESPONSE", "ALWAYS", "NEVER"}, false)),
						Optional:         true,
						Description:      "Specify when to compress responses.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"hd_data_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior specifies Akamai XML metadata that can only be configured on your behalf by Akamai Professional Services.  Unlike the `advanced` behavior, this may apply a different set of overriding metadata that executes in a post-processing phase. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "Human-readable description of what the XML block does.",
						Type:        schema.TypeString,
					},
					"xml": {
						Optional:    true,
						Description: "A block of Akamai XML metadata.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"health_detection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Monitors the health of your origin server by tracking unsuccessful attempts to contact it. Use this behavior to keep end users from having to wait several seconds before a forwarded request times out, or to reduce requests on the origin server when it is unavailable. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"retry_count": {
						Optional:    true,
						Description: "The number of consecutive connection failures that mark an IP address as faulty.",
						Type:        schema.TypeInt,
					},
					"retry_interval": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the amount of time the edge server will wait before trying to reconnect to an IP address it has already identified as faulty.",
						Type:             schema.TypeString,
					},
					"maximum_reconnects": {
						Optional:    true,
						Description: "Specifies the maximum number of times the edge server will contact your origin server. If your origin is associated with several IP addresses, `maximumReconnects` effectively overrides the value of `retryCount`.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"hsaf_eip_binding": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Edge IP Binding works with a limited set of static IP addresses to distribute your content, which can be limiting in large footprint environments. This behavior sets Hash Serial and Forward (HSAF) for Edge IP Binding to deal with larger footprints. It can only be configured on your behalf by Akamai Professional Services. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables HSAF for Edge IP Binding customers with a large footprint.",
						Type:        schema.TypeBool,
					},
					"custom_extracted_serial": {
						Optional:    true,
						Description: "Whether to pull the serial number from the variable value set in the `advanced` behavior. Work with your Akamai Services team to add the `advanced` behavior earlier in your property to extract and apply the `AKA_PM_EIP_HSAF_SERIAL` variable.",
						Type:        schema.TypeBool,
					},
					"hash_min_value": {
						Optional:    true,
						Description: "Specifies the minimum value for the HSAF hash range, from 2 through 2045. This needs to be lower than `hashMaxValue`.",
						Type:        schema.TypeInt,
					},
					"hash_max_value": {
						Optional:    true,
						Description: "Specifies the maximum value for the hash range, from 3 through 2046. This needs to be higher than `hashMinValue`.",
						Type:        schema.TypeInt,
					},
					"tier": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EDGE", "PARENT", "BOTH"}, false)),
						Optional:         true,
						Description:      "Specifies where the behavior is applied.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"http2": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables the HTTP/2 protocol, which reduces latency and improves efficiency. You can only apply this behavior if the property is marked as secure.  See `Secure property requirements` for guidance. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"http3": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This enables the HTTP/3 protocol that uses QUIC. The behavior allows for improved performance and faster connection setup. You can only apply this behavior if the property is marked as secure. See `Secure property requirements` and the `Property Manager documentation` for guidance. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "This enables HTTP/3 connections between requesting clients and Akamai edge servers. You also need to enable QUIC and TLS 1.3 in your certificate deployment settings. See the `Property Manager documentation` for more details.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"http_strict_transport_security": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Applies HTTP Strict Transport Security (HSTS), disallowing insecure HTTP traffic. Apply this to hostnames managed with Standard TLS or Enhanced TLS certificates. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "Applies HSTS to this set of requests.",
						Type:        schema.TypeBool,
					},
					"max_age": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ZERO_MINS", "TEN_MINS", "ONE_DAY", "ONE_MONTH", "THREE_MONTHS", "SIX_MONTHS", "ONE_YEAR", "TWO_YEARS"}, false)),
						Optional:         true,
						Description:      "Specifies the duration for which to apply HSTS for new browser connections.",
						Type:             schema.TypeString,
					},
					"include_sub_domains": {
						Optional:    true,
						Description: "When enabled, applies HSTS to all subdomains.",
						Type:        schema.TypeBool,
					},
					"preload": {
						Optional:    true,
						Description: "When enabled, adds this domain to the browser's preload list. You still need to declare the domain at `hstspreload.org`.",
						Type:        schema.TypeBool,
					},
					"redirect": {
						Optional:    true,
						Description: "When enabled, redirects all HTTP requests to HTTPS.",
						Type:        schema.TypeBool,
					},
					"redirect_status_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302})),
						Optional:         true,
						Description:      "Specifies a response code.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"http_to_https_upgrade": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Upgrades an HTTP edge request to HTTPS for the remainder of the request flow. Enable this behavior only if your origin supports HTTPS, and if your `origin` behavior is configured with `originCertsToHonor` to verify SSL certificates. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"upgrade": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"im_override": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This specifies common query parameters that affect how `imageManager` transforms images, potentially overriding policy, width, format, or density request parameters. This also allows you to assign the value of one of the property's `rule tree variables` to one of Image and Video Manager's own policy variables. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"override": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"POLICY", "POLICY_VARIABLE", "WIDTH", "FORMAT", "DPR", "EXCLUDE_QUERY"}, false)),
						Optional:         true,
						Description:      "Selects the type of query parameter you want to set.",
						Type:             schema.TypeString,
					},
					"typesel": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VALUE", "VARIABLE"}, false)),
						Optional:         true,
						Description:      "Specifies how to set a query parameter.",
						Type:             schema.TypeString,
					},
					"formatvar": {
						Optional:    true,
						Description: "This selects the variable with the name of the browser you want to optimize images for. The variable specifies the same type of data as the `format` option below.",
						Type:        schema.TypeString,
					},
					"format": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CHROME", "IE", "SAFARI", "GENERIC", "AVIF_WEBP_JPEG_PNG_GIF", "JP2_WEBP_JPEG_PNG_GIF", "WEBP_JPEG_PNG_GIF", "JPEG_PNG_GIF"}, false)),
						Optional:         true,
						Description:      "Specifies the type of the browser, or the encodings passed in the `Accept` header, that you want to optimize images for.",
						Type:             schema.TypeString,
					},
					"dprvar": {
						Optional:    true,
						Description: "This selects the variable with the desired pixel density. The variable specifies the same type of data as the `dpr` option below.",
						Type:        schema.TypeString,
					},
					"dpr": {
						Optional:    true,
						Description: "Directly specifies the pixel density. The numeric value is a scaling factor of 1, representing normal density.",
						Type:        schema.TypeFloat,
					},
					"widthvar": {
						Optional:    true,
						Description: "Selects the variable with the desired width.  If the Image and Video Manager policy doesn't define that width, it serves the next largest width.",
						Type:        schema.TypeString,
					},
					"width": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Sets the image's desired pixel width directly. If the Image Manager policy doesn't define that width, it serves the next largest width.",
						Type:             schema.TypeFloat,
					},
					"policyvar": {
						Optional:    true,
						Description: "This selects the variable with the desired Image and Video Manager policy name to apply to image requests. If there is no policy by that name, Image and Video Manager serves the image unmodified.",
						Type:        schema.TypeString,
					},
					"policy": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,32}$"),
						Optional:         true,
						Description:      "This selects the desired Image and Video Manager policy name directly. If there is no policy by that name, Image and Video Manager serves the image unmodified.",
						Type:             schema.TypeString,
					},
					"policyvar_name": {
						Optional:    true,
						Description: "This selects the name of one of the variables defined in an Image and Video Manager policy that you want to replace with the property's rule tree variable.",
						Type:        schema.TypeString,
					},
					"policyvar_i_mvar": {
						Optional:    true,
						Description: "This selects one of the property's rule tree variables to assign to the `policyvarName` variable within Image and Video Manager.",
						Type:        schema.TypeString,
					},
					"exclude_all_query_parameters": {
						Optional:    true,
						Description: "Whether to exclude all query parameters from the Image and Video Manager cache key.",
						Type:        schema.TypeBool,
					},
					"excluded_query_parameters": {
						Optional:    true,
						Description: "Specifies individual query parameters to exclude from the Image and Video Manager cache key.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"image_and_video_manager": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"policy_set_type": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"resize": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"apply_best_file_type": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"cp_code_original": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"cp_code_transformed": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"image_set": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]+([^-].|[^v])$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"video_set": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]+-v$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"image_manager": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Optimizes images' size or file type for the requesting device.  You can also use this behavior to generate API tokens to apply your own policies to matching images using the `Image and Video Manager API`. To apply this behavior, you need to match on a `fileExtension`. Once you apply Image and Video Manager to traffic, you can add the `advancedImMatch` to ensure the behavior applies to the requests from the Image and Video Manager backend. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"settings_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enable image management capabilities and generate a corresponding API token.",
						Type:        schema.TypeBool,
					},
					"resize": {
						Optional:    true,
						Description: "Specify whether to scale down images to the maximum screen resolution, as determined by the rendering device's user agent.  Note that enabling this may affect screen layout in unexpected ways.",
						Type:        schema.TypeBool,
					},
					"apply_best_file_type": {
						Optional:    true,
						Description: "Specify whether to convert images to the best file type for the requesting device, based on its user agent and the initial image file. This produces the smallest file size possible that retains image quality.",
						Type:        schema.TypeBool,
					},
					"super_cache_region": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"US", "ASIA", "AUSTRALIA", "EMEA", "JAPAN", "CHINA"}, false)),
						Optional:         true,
						Description:      "Specifies a location for your site's heaviest traffic, for use in caching derivatives on edge servers.",
						Type:             schema.TypeString,
					},
					"traffic_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"cp_code_original": {
						Optional:    true,
						Description: "Assigns a CP code to track traffic and billing for original images that the Image and Video Manager has not modified. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"cp_code_transformed": {
						Optional:    true,
						Description: "Assigns a separate CP code to track traffic and billing for derived images. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"api_reference_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"use_existing_policy_set": {
						Optional:    true,
						Description: "Whether to use a previously created policy set that may be referenced in other properties, or create a new policy set to use with this property. A policy set can be shared across multiple properties belonging to the same contract. The behavior populates any changes to the policy set across all properties that reference that set.",
						Type:        schema.TypeBool,
					},
					"policy_set": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]+([^-].|[^v])$"),
						Optional:         true,
						Description:      "Identifies the existing policy set configured with `Image and Video Manager API`.",
						Type:             schema.TypeString,
					},
					"advanced": {
						Optional:    true,
						Description: "Generates a custom `Image and Video Manager API` token to apply a corresponding policy to this set of images. The token consists of a descriptive label (the `policyToken`) concatenated with a property-specific identifier that's generated when you save the property. The API registers the token when you activate the property.",
						Type:        schema.TypeBool,
					},
					"policy_token": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,64}$"),
						Optional:         true,
						Description:      "Assign a prefix label to help match the policy token to this set of images, limited to 32 alphanumeric or underscore characters. If you don't specify a label, `default` becomes the prefix.",
						Type:             schema.TypeString,
					},
					"policy_token_default": {
						Optional:    true,
						Description: "Specify the default policy identifier, which is registered with the `Image and Video Manager API` once you activate this property.  The `advanced` option needs to be inactive.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"image_manager_video": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Optimizes videos managed by Image and Video Manager for the requesting device.  You can also use this behavior to generate API tokens to apply your own policies to matching videos using the `Image and Video Manager API`. To apply this behavior, you need to match on a `fileExtension`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"settings_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Applies Image and Video Manager's video optimization to the current content.",
						Type:        schema.TypeBool,
					},
					"resize": {
						Optional:    true,
						Description: "When enabled, scales down video for smaller mobile screens, based on the device's `User-Agent` header.",
						Type:        schema.TypeBool,
					},
					"apply_best_file_type": {
						Optional:    true,
						Description: "When enabled, automatically converts videos to the best file type for the requesting device. This produces the smallest file size that retains image quality, based on the user agent and the initial image file.",
						Type:        schema.TypeBool,
					},
					"super_cache_region": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"US", "ASIA", "AUSTRALIA", "EMEA", "JAPAN", "CHINA"}, false)),
						Optional:         true,
						Description:      "To optimize caching, assign a region close to your site's heaviest traffic.",
						Type:             schema.TypeString,
					},
					"traffic_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"cp_code_original": {
						Optional:    true,
						Description: "Specifies the CP code for which to track Image and Video Manager video traffic. Use this along with `cpCodeTransformed` to track traffic to derivative video content. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"cp_code_transformed": {
						Optional:    true,
						Description: "Specifies the CP code to identify derivative transformed video content. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"api_reference_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"use_existing_policy_set": {
						Optional:    true,
						Description: "Whether to use a previously created policy set that may be referenced in other properties, or create a new policy set to use with this property. A policy set can be shared across multiple properties belonging to the same contract. The behavior populates any changes to the policy set across all properties that reference that set.",
						Type:        schema.TypeBool,
					},
					"policy_set": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]+-v$"),
						Optional:         true,
						Description:      "Identifies the existing policy set configured with `Image and Video Manager API`.",
						Type:             schema.TypeString,
					},
					"advanced": {
						Optional:    true,
						Description: "When disabled, applies a single standard policy based on your property name.  Allows you to reference a rule-specific `policyToken` for videos with different match criteria.",
						Type:        schema.TypeBool,
					},
					"policy_token": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,64}$"),
						Optional:         true,
						Description:      "Specifies a custom policy defined in the Image and Video Manager Policy Manager or the `Image and Video Manager API`. The policy name can include up to 64 alphanumeric, dash, or underscore characters.",
						Type:             schema.TypeString,
					},
					"policy_token_default": {
						Optional:    true,
						Description: "Specifies the default policy identifier, which is registered with the `Image and Video Manager API` once you activate this property.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"include": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Includes let you reuse chunks of a property configuration that you can manage separately from the rest of the property rule tree. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"id": {
						Optional:    true,
						Description: "Identifies the include you want to add to your rule tree. You can get the include ID using `PAPI`. This option only accepts digits, without the `inc_` ID prefix.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"instant": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Instant feature allows you to prefetch content to the edge cache by adding link relation attributes to markup. For example: This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"prefetch_cacheable": {
						Optional:    true,
						Description: "When enabled, applies prefetching only to objects already set to be cacheable, for example using the `caching` behavior. Only applies to content with the `tieredDistribution` behavior enabled.",
						Type:        schema.TypeBool,
					},
					"prefetch_no_store": {
						Optional:    true,
						Description: "Allows otherwise non-cacheable `no-store` content to prefetch if the URL path ends with `/` to indicate a request for a default file, or if the extension matches the value of the `prefetchNoStoreExtensions` option. Only applies to content with the `sureRoute` behavior enabled.",
						Type:        schema.TypeBool,
					},
					"prefetch_no_store_extensions": {
						Optional:    true,
						Description: "Specifies a set of file extensions for which the `prefetchNoStore` option is allowed.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"prefetch_html": {
						Optional:    true,
						Description: "Allows edge servers to prefetch additional HTML pages while pages that link to them are being delivered. This only applies to links from `<a>` or `<link>` tags with the appropriate link relation attribute.",
						Type:        schema.TypeBool,
					},
					"custom_link_relations": {
						Optional:    true,
						Description: "Specify link relation values that activate the prefetching behavior. For example, specifying `fetch` allows you to use shorter `rel=\"fetch\"` markup.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"instant_config": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Multi-Domain Configuration, also known as `InstantConfig`, allows you to apply property settings to all incoming hostnames based on a DNS lookup, without explicitly listing them among the property's hostnames. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the InstantConfig behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"large_file_optimization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Large File Optimization` (LFO) feature improves performance and reliability when delivering large files. You need this behavior for objects larger than 1.8GB, and you should apply it to anything over 100MB. You should apply it only to the specific content to be optimized, such as a download directory's `.gz` files, and enable the `useVersioning` option while enforcing your own filename versioning policy. Make sure you meet all the `requirements and best practices` for the LFO delivery. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the file optimization behavior.",
						Type:        schema.TypeBool,
					},
					"enable_partial_object_caching": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PARTIAL_OBJECT_CACHING", "NON_PARTIAL_OBJECT_CACHING"}, false)),
						Optional:         true,
						Description:      "Specifies whether to cache partial objects.",
						Type:             schema.TypeString,
					},
					"minimum_size": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d+[K,M,G,T]B$"),
						Optional:         true,
						Description:      "Optimization only applies to files larger than this, expressed as a number suffixed with a unit string such as `MB` or `GB`.",
						Type:             schema.TypeString,
					},
					"maximum_size": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d+[K,M,G,T]B$"),
						Optional:         true,
						Description:      "Optimization does not apply to files larger than this, expressed as a number suffixed with a unit string such as `MB` or `GB`. The size of a file can't be greater than 323 GB. If you need to optimize a larger file, contact Akamai Professional Services for help. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"use_versioning": {
						Optional:    true,
						Description: "When `enablePartialObjectCaching` is set to `PARTIAL_OBJECT_CACHING`, enabling this option signals your intention to vary filenames by version, strongly recommended to avoid serving corrupt content when chunks come from different versions of the same file.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"large_file_optimization_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Large File Optimization` feature improves performance and reliability when delivering large files. You need this behavior for objects larger than 1.8GB, and it's recommended for anything over 100MB. You should apply it only to the specific content to be optimized, such as a download directory's `.gz` files.  Note that it is best to use `NetStorage` for objects larger than 1.8GB. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the file optimization behavior.",
						Type:        schema.TypeBool,
					},
					"object_size": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d+[K,M,G,T]B$"),
						Optional:         true,
						Description:      "Specifies the size of the file at which point to apply partial object (POC) caching. Append a numeric value with a `MB` or `GB` suffix.",
						Type:             schema.TypeString,
					},
					"fragment_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HALF_MB", "ONE_MB", "TWO_MB", "FOUR_MB"}, false)),
						Optional:         true,
						Description:      "Specifies the size of each fragment used for partial object caching.",
						Type:             schema.TypeString,
					},
					"prefetch_during_request": {
						Optional:    true,
						Description: "The number of POC fragments to prefetch during the request.",
						Type:        schema.TypeInt,
					},
					"prefetch_after_request": {
						Optional:    true,
						Description: "The number of POC fragments to prefetch after the request.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"limit_bit_rate": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Control the rate at which content serves out to end users, optionally varying the speed depending on the file size or elapsed download time. Each bit rate specified in the `bitrateTable` array corresponds to a `thresholdTable` entry that activates it.  You can use this behavior to prevent media downloads from progressing faster than they are viewed, for example, or to differentiate various tiers of end-user experience. To apply this behavior, you should match on a `contentType`, `path`, or `filename`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, activates the bit rate limiting behavior.",
						Type:        schema.TypeBool,
					},
					"bitrate_table": {
						Optional:    true,
						Description: "Specifies a download rate that corresponds to a `thresholdTable` entry. The bit rate appears as a two-member object consisting of a numeric `bitrateValue` and a `bitrateUnit` string, with allowed values of `Kbps`, `Mbps`, and `Gbps`.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"bitrate_value": {
									Optional:    true,
									Description: "The numeric indicator of the download rate.",
									Type:        schema.TypeFloat,
								},
								"bitrate_unit": {
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"KBPS", "MBPS", "GBPS"}, false)),
									Optional:         true,
									Description:      "The unit of measurement, either `KBPS`, `MBPS`, or `GBPS`.",
									Type:             schema.TypeString,
								},
							},
						},
					},
					"threshold_table": {
						Optional:    true,
						Description: "Specifies the minimum size of the file or the amount of elapsed download time before applying the bit rate limit from the corresponding `bitrateTable` entry. The threshold appears as a two-member object consisting of a numeric `thresholdValue` and `thresholdUnit` string, with allowed values of `SECONDS` or `BYTES`.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"threshold_value": {
									Optional:    true,
									Description: "The numeric indicator of the minimum file size or elapsed download time.",
									Type:        schema.TypeInt,
								},
								"threshold_unit": {
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BYTES", "SECONDS"}, false)),
									Optional:         true,
									Description:      "The unit of measurement, either `SECONDS` of the elapsed download time, or `BYTES` of the file size.",
									Type:             schema.TypeString,
								},
							},
						},
					},
				},
			},
		},
		"log_custom": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Logs custom details from the origin response in the `Log Delivery Service` report. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"log_custom_log_field": {
						Optional:    true,
						Description: "Whether to append additional custom data to each log line.",
						Type:        schema.TypeBool,
					},
					"custom_log_field": {
						Optional:    true,
						Description: "Specifies an additional data field to append to each log line, maximum 1000 bytes, typically based on a dynamically generated built-in system variable. For example, `round-trip: {{builtin.AK_CLIENT_TURNAROUND_TIME}}ms` logs the total time to complete the response. See `Support for variables` for more information. Since this option can specify both a request and response, it overrides any `customLogField` settings in the `report` behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"m_pulse": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "`mPulse` provides high-level performance analytics and predictive recommendations based on real end user data. See the `mPulse Quick Start` to set up mPulse on your website. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Applies performance monitoring to this behavior's set of content.",
						Type:        schema.TypeBool,
					},
					"require_pci": {
						Optional:    true,
						Description: "Suppresses gathering metrics for potentially sensitive end-user interactions. Enabling this omits data from some older browsers.",
						Type:        schema.TypeBool,
					},
					"loader_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"V10", "V12", "LATEST", "BETA"}, false)),
						Optional:         true,
						Description:      "Specifies the version of the Boomerang JavaScript loader snippet. See `mPulse Loader Snippets` for more information.",
						Type:             schema.TypeString,
					},
					"title_optional": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"api_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^$|^[a-zA-Z2-9]{5}-[a-zA-Z2-9]{5}-[a-zA-Z2-9]{5}-[a-zA-Z2-9]{5}-[a-zA-Z2-9]{5}$")),
						Optional:         true,
						Description:      "This generated value uniquely identifies sections of your website for you to analyze independently. To access this value, see `Enable mPulse in Property Manager`.",
						Type:             schema.TypeString,
					},
					"buffer_size": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(1[5-9][0-9]|1[0-9]{3}|[2-9][0-9]{2,3})$")),
						Optional:         true,
						Description:      "Allows you to override the browser's default (150) maximum number of reported performance timeline entries.",
						Type:             schema.TypeString,
					},
					"config_override": {
						Optional:    true,
						Description: "A JSON string representing a configuration object passed to the JavaScript library under which mPulse runs. It corresponds at run-time to the `window.BOOMR_config` object. For example, this turns on monitoring of Single Page App frameworks: `\"{\\\"history\\\": {\\\"enabled\\\": true, \\\"auto\\\": true}}\"`.  See `Configuration Overrides` for more information.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"manifest_personalization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows customers who use the Adaptive Media Delivery product to enhance content based on the capabilities of each end user's device.  This behavior configures a `manifest` for both HLS Live and on-demand streaming. For more information, see `Adaptive Media Delivery`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Manifest Personalization feature.",
						Type:        schema.TypeBool,
					},
					"hls_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"hls_enabled": {
						Optional:    true,
						Description: "Allows you to customize the HLS master manifest that's sent to the requesting client.",
						Type:        schema.TypeBool,
					},
					"hls_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BEST_PRACTICE", "CUSTOM"}, false)),
						Optional:         true,
						Description:      "Applies with `hlsEnabled` on.",
						Type:             schema.TypeString,
					},
					"hls_preferred_bitrate": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^\\d+$")),
						Optional:         true,
						Description:      "Sets the preferred bit rate in Kbps. This causes the media playlist specified in the `#EXT-X-STREAM-INF` tag that most closely matches the value to list first. All other playlists maintain their current position in the manifest.",
						Type:             schema.TypeString,
					},
					"hls_filter_in_bitrates": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^\\d+(,\\d+)*$")),
						Optional:         true,
						Description:      "Specifies a comma-delimited set of preferred bit rates, such as `100,200,400`. Playlists specified in the `#EXT-X-STREAM-INF` tag with bit rates outside of any of those values by up to 100 Kbps are excluded from the manifest.",
						Type:             schema.TypeString,
					},
					"hls_filter_in_bitrate_ranges": {
						Optional:    true,
						Description: "Specifies a comma-delimited set of bit rate ranges, such as `100-400,1000-4000`. Playlists specified in the `#EXT-X-STREAM-INF` tag with bit rates outside of any of those ranges are excluded from the manifest.",
						Type:        schema.TypeString,
					},
					"hls_query_param_enabled": {
						Optional:    true,
						Description: "Specifies query parameters for the HLS master manifest to customize the manifest's content.  Any settings specified in the query string override those already configured in Property Manager.",
						Type:        schema.TypeBool,
					},
					"hls_query_param_secret_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(0x)?[0-9a-fA-F]{32}$")),
						Optional:         true,
						Description:      "Specifies a primary key as a token to accompany the request.",
						Type:             schema.TypeString,
					},
					"hls_query_param_transition_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(0x)?[0-9a-fA-F]{32}$")),
						Optional:         true,
						Description:      "Specifies a transition key as a token to accompany the request.",
						Type:             schema.TypeString,
					},
					"hls_show_advanced": {
						Optional:    true,
						Description: "Allows you to configure advanced settings.",
						Type:        schema.TypeBool,
					},
					"hls_enable_debug_headers": {
						Optional:    true,
						Description: "Includes additional `Akamai-Manifest-Personalization` and `Akamai-Manifest-Personalization-Config-Source` debugging headers.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"manifest_rerouting": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior works with `adScalerCircuitBreaker`. It delegates parts of the media delivery workflow, like ad insertion, to other technology partners. Akamai reroutes manifest file requests to partner platforms for processing prior to being delivered. Rerouting simplifies the workflow and improves the media streaming experience. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"partner": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"adobe_primetime"}, false)),
						Optional:         true,
						Description:      "Set this value to `adobe_primetime`, which is an external technology partner that provides value added offerings, like advertisement integration, to the requested media objects.",
						Type:             schema.TypeString,
					},
					"username": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9A-Za-z!@.,;:'\"?-]{1,50}$"),
						Optional:         true,
						Description:      "The user name for your Adobe Primetime account.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"manual_server_push": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "With the `http2` behavior enabled, this loads a specified set of objects into the client browser's cache. To apply this behavior, you should match on a `path` or `filename`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"serverpushlist": {
						Optional:    true,
						Description: "Specifies the set of objects to load into the client browser's cache over HTTP2. Each value in the array represents a hostname and full path to the object, such as `www.example.com/js/site.js`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"media_acceleration": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Enables Accelerated Media Delivery for this set of requests. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables Media Acceleration.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"media_acceleration_quic_optout": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"optout": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"media_client": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables client-side download analytics.",
						Type:        schema.TypeBool,
					},
					"beacon_id": {
						Optional:    true,
						Description: "Specifies the ID of data source's beacon.",
						Type:        schema.TypeString,
					},
					"use_hybrid_http_udp": {
						Optional:    true,
						Description: "Enables the hybrid HTTP/UDP protocol.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"media_file_retrieval_optimization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Media File Retrieval Optimization (MFRO) speeds the delivery of large media files by relying on caches of partial objects. You should use it for files larger than 100 MB. It's required for files larger than 1.8 GB, and works best with `NetStorage`. To apply this behavior, you should match on a `fileExtension`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the partial-object caching behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"media_origin_failover": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies how edge servers respond when the origin is unresponsive, or suffers from server or content errors. You can specify how many times to retry, switch to a backup origin hostname, or configure a redirect. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"detect_origin_unresponsive_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_origin_unresponsive": {
						Optional:    true,
						Description: "Allows you to configure what happens when the origin is unresponsive.",
						Type:        schema.TypeBool,
					},
					"origin_unresponsive_detection_level": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AGGRESSIVE", "CONSERVATIVE", "MODERATE"}, false)),
						Optional:         true,
						Description:      "Specify the level of response to slow origin connections.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_blacklist_origin_ip": {
						Optional:    true,
						Description: "Enabling this blacklists the origin's IP address.",
						Type:        schema.TypeBool,
					},
					"origin_unresponsive_blacklist_window": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"TEN_S", "THIRTY_S"}, false)),
						Optional:         true,
						Description:      "This sets the delay before blacklisting an IP address.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_recovery": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RETRY_X_TIMES", "SWITCH_TO_BACKUP_ORIGIN", "REDIRECT_TO_DIFFERENT_ORIGIN_LOCATION"}, false)),
						Optional:         true,
						Description:      "This sets the recovery option.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_retry_limit": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ONE", "TWO", "THREE"}, false)),
						Optional:         true,
						Description:      "Sets how many times to retry.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_backup_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the origin hostname.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_alternate_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the redirect's destination hostname.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_modify_request_path": {
						Optional:    true,
						Description: "Modifies the request path.",
						Type:        schema.TypeBool,
					},
					"origin_unresponsive_modified_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "This specifies the path to form the new URL.",
						Type:             schema.TypeString,
					},
					"origin_unresponsive_include_query_string": {
						Optional:    true,
						Description: "Enabling this includes the original set of query parameters.",
						Type:        schema.TypeBool,
					},
					"origin_unresponsive_redirect_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302})),
						Optional:         true,
						Description:      "Specifies the redirect response code.",
						Type:             schema.TypeInt,
					},
					"origin_unresponsive_change_protocol": {
						Optional:    true,
						Description: "This allows you to change the request protocol.",
						Type:        schema.TypeBool,
					},
					"origin_unresponsive_protocol": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "Specifies which protocol to use.",
						Type:             schema.TypeString,
					},
					"detect_origin_unavailable_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_origin_unavailable": {
						Optional:    true,
						Description: "Allows you to configure failover settings when the origin server responds with errors.",
						Type:        schema.TypeBool,
					},
					"origin_unavailable_detection_level": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RESPONSE_CODES"}, false)),
						Optional:         true,
						Description:      "Specify `RESPONSE_CODES`, the only available option.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_response_codes": {
						Optional:    true,
						Description: "Specifies the set of response codes identifying when the origin responds with errors.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"origin_unavailable_blacklist_origin_ip": {
						Optional:    true,
						Description: "Enabling this blacklists the origin's IP address.",
						Type:        schema.TypeBool,
					},
					"origin_unavailable_blacklist_window": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"TEN_S", "THIRTY_S"}, false)),
						Optional:         true,
						Description:      "This sets the delay before blacklisting an IP address.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_recovery": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RETRY_X_TIMES", "SWITCH_TO_BACKUP_ORIGIN", "REDIRECT_TO_DIFFERENT_ORIGIN_LOCATION"}, false)),
						Optional:         true,
						Description:      "This sets the recovery option.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_retry_limit": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ONE", "TWO", "THREE"}, false)),
						Optional:         true,
						Description:      "Sets how many times to retry.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_backup_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the origin hostname.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_alternate_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the redirect's destination hostname.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_modify_request_path": {
						Optional:    true,
						Description: "Modifies the request path.",
						Type:        schema.TypeBool,
					},
					"origin_unavailable_modified_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "This specifies the path to form the new URL.",
						Type:             schema.TypeString,
					},
					"origin_unavailable_include_query_string": {
						Optional:    true,
						Description: "Enabling this includes the original set of query parameters.",
						Type:        schema.TypeBool,
					},
					"origin_unavailable_redirect_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302})),
						Optional:         true,
						Description:      "Specifies either a redirect response code.",
						Type:             schema.TypeInt,
					},
					"origin_unavailable_change_protocol": {
						Optional:    true,
						Description: "Modifies the request protocol.",
						Type:        schema.TypeBool,
					},
					"origin_unavailable_protocol": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "Specifies either the `HTTP` or `HTTPS` protocol.",
						Type:             schema.TypeString,
					},
					"detect_object_unavailable_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"detect_object_unavailable": {
						Optional:    true,
						Description: "Allows you to configure failover settings when the origin has content errors.",
						Type:        schema.TypeBool,
					},
					"object_unavailable_detection_level": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RESPONSE_CODES"}, false)),
						Optional:         true,
						Description:      "Specify `RESPONSE_CODES`, the only available option.",
						Type:             schema.TypeString,
					},
					"object_unavailable_response_codes": {
						Optional:    true,
						Description: "Specifies the set of response codes identifying when there are content errors.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"object_unavailable_blacklist_origin_ip": {
						Optional:    true,
						Description: "Enabling this blacklists the origin's IP address.",
						Type:        schema.TypeBool,
					},
					"object_unavailable_blacklist_window": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"TEN_S", "THIRTY_S"}, false)),
						Optional:         true,
						Description:      "This sets the delay before blacklisting an IP address.",
						Type:             schema.TypeString,
					},
					"object_unavailable_recovery": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RETRY_X_TIMES", "SWITCH_TO_BACKUP_ORIGIN", "REDIRECT_TO_DIFFERENT_ORIGIN_LOCATION"}, false)),
						Optional:         true,
						Description:      "This sets the recovery option.",
						Type:             schema.TypeString,
					},
					"object_unavailable_retry_limit": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ONE", "TWO", "THREE"}, false)),
						Optional:         true,
						Description:      "Sets how many times to retry.",
						Type:             schema.TypeString,
					},
					"object_unavailable_backup_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the origin hostname.",
						Type:             schema.TypeString,
					},
					"object_unavailable_alternate_host": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies the redirect's destination hostname.",
						Type:             schema.TypeString,
					},
					"object_unavailable_modify_request_path": {
						Optional:    true,
						Description: "Enabling this allows you to modify the request path.",
						Type:        schema.TypeBool,
					},
					"object_unavailable_modified_path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "This specifies the path to form the new URL.",
						Type:             schema.TypeString,
					},
					"object_unavailable_include_query_string": {
						Optional:    true,
						Description: "Enabling this includes the original set of query parameters.",
						Type:        schema.TypeBool,
					},
					"object_unavailable_redirect_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302})),
						Optional:         true,
						Description:      "Specifies a redirect response code.",
						Type:             schema.TypeInt,
					},
					"object_unavailable_change_protocol": {
						Optional:    true,
						Description: "Changes the request protocol.",
						Type:        schema.TypeBool,
					},
					"object_unavailable_protocol": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "Specifies either the `HTTP` or `HTTPS` protocol.",
						Type:             schema.TypeString,
					},
					"other_options": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"client_response_code": {
						Optional:    true,
						Description: "Specifies the response code served to the client.",
						Type:        schema.TypeString,
					},
					"cache_error_response": {
						Optional:    true,
						Description: "When enabled, caches the error response.",
						Type:        schema.TypeBool,
					},
					"cache_window": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ONE_S", "TEN_S", "THIRTY_S"}, false)),
						Optional:         true,
						Description:      "This sets error response's TTL.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"metadata_caching": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior reduces time spent waiting for the initial response, also known as time to first byte, during peak traffic events. Contact Akamai Professional Services for help configuring it. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables metadata caching.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"mobile_sdk_performance": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Mobile App Performance SDK.",
						Type:        schema.TypeBool,
					},
					"secondary_multipath_to_origin": {
						Optional:    true,
						Description: "When enabled, sends secondary multi-path requests to the origin server.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"modify_incoming_request_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modify, add, remove, or pass along specific request headers coming upstream from the client. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ADD", "DELETE", "MODIFY", "PASS"}, false)),
						Optional:         true,
						Description:      "Either `ADD`, `DELETE`, `MODIFY`, or `PASS` incoming HTTP request headers.",
						Type:             schema.TypeString,
					},
					"standard_add_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ACCEPT_ENCODING", "ACCEPT_LANGUAGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `ADD`, this specifies the name of the field to add.",
						Type:             schema.TypeString,
					},
					"standard_delete_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IF_MODIFIED_SINCE", "VIA", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `DELETE`, this specifies the name of the field to remove.",
						Type:             schema.TypeString,
					},
					"standard_modify_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ACCEPT_ENCODING", "ACCEPT_LANGUAGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `MODIFY`, this specifies the name of the field to modify.",
						Type:             schema.TypeString,
					},
					"standard_pass_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ACCEPT_ENCODING", "ACCEPT_LANGUAGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `PASS`, this specifies the name of the field to pass through.",
						Type:             schema.TypeString,
					},
					"custom_header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies a custom field name that applies when the relevant `standard` header name is set to `OTHER`.",
						Type:             schema.TypeString,
					},
					"header_value": {
						Optional:    true,
						Description: "Specifies the new header value.",
						Type:        schema.TypeString,
					},
					"new_header_value": {
						Optional:    true,
						Description: "Supplies an HTTP header replacement value.",
						Type:        schema.TypeString,
					},
					"avoid_duplicate_headers": {
						Optional:    true,
						Description: "When enabled with the `action` set to `MODIFY`, prevents creation of more than one instance of a header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"modify_incoming_response_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modify, add, remove, or pass along specific response headers coming downstream from the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ADD", "DELETE", "MODIFY", "PASS"}, false)),
						Optional:         true,
						Description:      "Either `ADD`, `DELETE`, `MODIFY`, or `PASS` incoming HTTP response headers.",
						Type:             schema.TypeString,
					},
					"standard_add_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_TYPE", "EDGE_CONTROL", "EXPIRES", "LAST_MODIFIED", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `ADD`, this specifies the name of the field to add.",
						Type:             schema.TypeString,
					},
					"standard_delete_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_TYPE", "VARY", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `DELETE`, this specifies the name of the field to remove.",
						Type:             schema.TypeString,
					},
					"standard_modify_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_TYPE", "EDGE_CONTROL", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `MODIFY`, this specifies the name of the field to modify.",
						Type:             schema.TypeString,
					},
					"standard_pass_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "EXPIRES", "PRAGMA", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `PASS`, this specifies the name of the field to pass through.",
						Type:             schema.TypeString,
					},
					"custom_header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies a custom field name that applies when the relevant `standard` header name is set to `OTHER`.",
						Type:             schema.TypeString,
					},
					"header_value": {
						Optional:    true,
						Description: "Specifies the header's new value.",
						Type:        schema.TypeString,
					},
					"new_header_value": {
						Optional:    true,
						Description: "Specifies an HTTP header replacement value.",
						Type:        schema.TypeString,
					},
					"avoid_duplicate_headers": {
						Optional:    true,
						Description: "When enabled with the `action` set to `MODIFY`, prevents creation of more than one instance of a header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"modify_outgoing_request_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modify, add, remove, or pass along specific request headers going upstream towards the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ADD", "DELETE", "MODIFY", "REGEX"}, false)),
						Optional:         true,
						Description:      "Either `ADD` or `DELETE` outgoing HTTP request headers, `MODIFY` their fixed values, or specify a `REGEX` pattern to transform them.",
						Type:             schema.TypeString,
					},
					"standard_add_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"USER_AGENT", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `ADD`, this specifies the name of the field to add.",
						Type:             schema.TypeString,
					},
					"standard_delete_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PRAGMA", "USER_AGENT", "VIA", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `DELETE`, this specifies the name of the field to remove.",
						Type:             schema.TypeString,
					},
					"standard_modify_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"USER_AGENT", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `MODIFY` or `REGEX`, this specifies the name of the field to modify.",
						Type:             schema.TypeString,
					},
					"custom_header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies a custom field name that applies when the relevant `standard` header name is set to `OTHER`.",
						Type:             schema.TypeString,
					},
					"header_value": {
						Optional:    true,
						Description: "Specifies the new header value.",
						Type:        schema.TypeString,
					},
					"new_header_value": {
						Optional:    true,
						Description: "Specifies an HTTP header replacement value.",
						Type:        schema.TypeString,
					},
					"regex_header_match": {
						Optional:    true,
						Description: "Specifies a Perl-compatible regular expression to match within the header value.",
						Type:        schema.TypeString,
					},
					"regex_header_replace": {
						Optional:    true,
						Description: "Specifies text that replaces the `regexHeaderMatch` pattern within the header value.",
						Type:        schema.TypeString,
					},
					"match_multiple": {
						Optional:    true,
						Description: "When enabled with the `action` set to `REGEX`, replaces all occurrences of the matched regular expression, otherwise only the first match if disabled.",
						Type:        schema.TypeBool,
					},
					"avoid_duplicate_headers": {
						Optional:    true,
						Description: "When enabled with the `action` set to `MODIFY`, prevents creation of more than one instance of a header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"modify_outgoing_response_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modify, add, remove, or pass along specific response headers going downstream towards the client. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ADD", "DELETE", "MODIFY", "REGEX"}, false)),
						Optional:         true,
						Description:      "Either `ADD` or `DELETE` outgoing HTTP response headers, `MODIFY` their fixed values, or specify a `REGEX` pattern to transform them.",
						Type:             schema.TypeString,
					},
					"standard_add_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_DISPOSITION", "CONTENT_TYPE", "EDGE_CONTROL", "P3P", "PRAGMA", "ACCESS_CONTROL_ALLOW_ORIGIN", "ACCESS_CONTROL_ALLOW_METHODS", "ACCESS_CONTROL_ALLOW_HEADERS", "ACCESS_CONTROL_EXPOSE_HEADERS", "ACCESS_CONTROL_ALLOW_CREDENTIALS", "ACCESS_CONTROL_MAX_AGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `ADD`, this specifies the name of the field to add.",
						Type:             schema.TypeString,
					},
					"standard_delete_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_DISPOSITION", "CONTENT_TYPE", "EXPIRES", "P3P", "PRAGMA", "ACCESS_CONTROL_ALLOW_ORIGIN", "ACCESS_CONTROL_ALLOW_METHODS", "ACCESS_CONTROL_ALLOW_HEADERS", "ACCESS_CONTROL_EXPOSE_HEADERS", "ACCESS_CONTROL_ALLOW_CREDENTIALS", "ACCESS_CONTROL_MAX_AGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `DELETE`, this specifies the name of the field to remove.",
						Type:             schema.TypeString,
					},
					"standard_modify_header_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CACHE_CONTROL", "CONTENT_DISPOSITION", "CONTENT_TYPE", "P3P", "PRAGMA", "ACCESS_CONTROL_ALLOW_ORIGIN", "ACCESS_CONTROL_ALLOW_METHODS", "ACCESS_CONTROL_ALLOW_HEADERS", "ACCESS_CONTROL_EXPOSE_HEADERS", "ACCESS_CONTROL_ALLOW_CREDENTIALS", "ACCESS_CONTROL_MAX_AGE", "OTHER"}, false)),
						Optional:         true,
						Description:      "If the value of `action` is `MODIFY` or `REGEX`, this specifies the name of the field to modify.",
						Type:             schema.TypeString,
					},
					"custom_header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies a custom field name that applies when the relevant `standard` header name is set to `OTHER`.",
						Type:             schema.TypeString,
					},
					"header_value": {
						Optional:    true,
						Description: "Specifies the existing value of the header to match.",
						Type:        schema.TypeString,
					},
					"new_header_value": {
						Optional:    true,
						Description: "Specifies the new HTTP header replacement value.",
						Type:        schema.TypeString,
					},
					"regex_header_match": {
						Optional:    true,
						Description: "Specifies a Perl-compatible regular expression to match within the header value.",
						Type:        schema.TypeString,
					},
					"regex_header_replace": {
						Optional:    true,
						Description: "Specifies text that replaces the `regexHeaderMatch` pattern within the header value.",
						Type:        schema.TypeString,
					},
					"match_multiple": {
						Optional:    true,
						Description: "When enabled with the `action` set to `REGEX`, replaces all occurrences of the matched regular expression, otherwise only the first match if disabled.",
						Type:        schema.TypeBool,
					},
					"avoid_duplicate_headers": {
						Optional:    true,
						Description: "When enabled with the `action` set to `MODIFY`, prevents creation of more than one instance of a header. The last header clobbers others with the same name. This option affects the entire set of outgoing headers, and is not confined to the subset of regular expression matches.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"modify_via_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Removes or renames the HTTP `Via` headers used to inform the server of proxies through which the request was sent to the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables `Via` header modifications.",
						Type:        schema.TypeBool,
					},
					"modification_option": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REMOVE_HEADER", "RENAME_HEADER"}, false)),
						Optional:         true,
						Description:      "Specify how you want to handle the header.",
						Type:             schema.TypeString,
					},
					"rename_header_to": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies a new name to replace the existing `Via` header.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"mtls_origin_keystore": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Establish a Mutual TLS (mTLS) connection between the edge server and the origin to authenticate requests. This ensures that the requests to your origin server come directly from the Akamai network. In the mTLS protocol, the origin asks the edge server to present its identity certificate. For this negotiation to work, either the origin needs to be configured for mTLS sessions, or the edge server is allowed to proceed without the edge certificate, effectively performing a standard (non-mutual) TLS connection to the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "Allows a specific mutual transport layer (mTLS) client certificate in a request from the edge server to the origin.",
						Type:        schema.TypeBool,
					},
					"client_certificate_version_guid": {
						Optional:    true,
						Description: "Specifies the client certificate to authenticate your origin with the edge server. You need to create client certificates using the `Mutual TLS Origin Keystore` API or application.",
						Type:        schema.TypeString,
					},
					"auth_client_cert": {
						Optional:    true,
						Description: "When enabled, the edge server requires a prompt from the origin for the client certificate's identity. If the edge server gets the request, it proceeds with the mTLS session and connects to the origin. If the edge server doesn’t get the request, the connection to the origin stops and a client error is reported. When disabled, the edge server proceeds without a request for the client certificate, making a standard TLS connection to the origin. Disabled by default.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"origin": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specify the hostname and settings used to contact the origin once service begins. You can use your own origin, `NetStorage`, `Media Services Live`, an Edge Load Balancing origin, Akamai `Object Storage`, or a SaaS dynamic origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origin_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CUSTOMER", "NET_STORAGE", "MEDIA_SERVICE_LIVE", "EDGE_LOAD_BALANCING_ORIGIN_GROUP", "SAAS_DYNAMIC_ORIGIN"}, false)),
						Optional:         true,
						Description:      "Choose where your content is retrieved from.",
						Type:             schema.TypeString,
					},
					"net_storage": {
						Optional:    true,
						Description: "Specifies the details of the NetStorage server.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"origin_id": {
						Optional:    true,
						Description: "Identifies the Edge Load Balancing origin. This needs to correspond to an `edgeLoadBalancingOrigin` behavior's `id` attribute within the same property.",
						Type:        schema.TypeString,
					},
					"hostname": {
						Optional:    true,
						Description: "Specifies the hostname or IPv4 address of your origin server, from which edge servers can retrieve your content.",
						Type:        schema.TypeString,
					},
					"second_hostname_enabled": {
						Optional:    true,
						Description: "Available only for certain products. This specifies whether you want to use an additional origin server address.",
						Type:        schema.TypeBool,
					},
					"second_hostname": {
						Optional:    true,
						Description: "Specifies the origin server's hostname, IPv4 address, or IPv6 address. Edge servers retrieve your content from this origin server.",
						Type:        schema.TypeString,
					},
					"mslorigin": {
						Optional:    true,
						Description: "This specifies the media's origin server.",
						Type:        schema.TypeString,
					},
					"saas_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HOSTNAME", "PATH", "QUERY_STRING", "COOKIE"}, false)),
						Optional:         true,
						Description:      "Specifies the part of the request that identifies this SaaS dynamic origin.",
						Type:             schema.TypeString,
					},
					"saas_cname_enabled": {
						Optional:    true,
						Description: "Enabling this allows you to use a `CNAME chain` to determine the hostname for this SaaS dynamic origin.",
						Type:        schema.TypeBool,
					},
					"saas_cname_level": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the desired number of hostnames to use in the `CNAME chain`, starting backwards from the edge server.",
						Type:             schema.TypeInt,
					},
					"saas_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name of the cookie that identifies this SaaS dynamic origin.",
						Type:             schema.TypeString,
					},
					"saas_query_string": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "Specifies the name of the query parameter that identifies this SaaS dynamic origin.",
						Type:             schema.TypeString,
					},
					"saas_regex": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9\\:\\[\\]\\{\\}\\(\\)\\.\\?_\\-\\*\\+\\^\\$\\\\\\/\\|&=!]{1,250})$"),
						Optional:         true,
						Description:      "Specifies the Perl-compatible regular expression match that identifies this SaaS dynamic origin.",
						Type:             schema.TypeString,
					},
					"saas_replace": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9]|\\$[1-9])(([a-zA-Z0-9\\._\\-]|\\$[1-9]){0,250}([a-zA-Z0-9]|\\$[1-9]))?){1,10}$"),
						Optional:         true,
						Description:      "Specifies replacement text for what `saasRegex` matches.",
						Type:             schema.TypeString,
					},
					"saas_suffix": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})\\.(com|net|org|info|biz|us|co\\.uk|ac\\.uk|org\\.uk|me\\.uk|ca|eu|com\\.au|co|co\\.za|ru|es|me|tv|pro|in|ie|de|it|nl|fr|co\\.il|ch|se|co\\.nz|pl|jp|name|mobi|cc|ws|be|com\\.mx|at|nu|asia|co\\.nz|net\\.nz|org\\.nz|com\\.au|net\\.au|org\\.au|tools)$"),
						Optional:         true,
						Description:      "Specifies the static part of the SaaS dynamic origin.",
						Type:             schema.TypeString,
					},
					"forward_host_header": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REQUEST_HOST_HEADER", "ORIGIN_HOSTNAME", "CUSTOM"}, false)),
						Optional:         true,
						Description:      "When the `originType` is set to either `CUSTOMER` or `SAAS_DYNAMIC_ORIGIN`, this specifies which `Host` header to pass to the origin.",
						Type:             schema.TypeString,
					},
					"custom_forward_host_header": {
						Optional:    true,
						Description: "This specifies the name of the custom host header the edge server should pass to the origin.",
						Type:        schema.TypeString,
					},
					"cache_key_hostname": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REQUEST_HOST_HEADER", "ORIGIN_HOSTNAME"}, false)),
						Optional:         true,
						Description:      "Specifies the hostname to use when forming a cache key.",
						Type:             schema.TypeString,
					},
					"ip_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IPV4", "DUALSTACK", "IPV6"}, false)),
						Optional:         true,
						Description:      "Specifies which IP version to use when getting content from the origin.",
						Type:             schema.TypeString,
					},
					"use_unique_cache_key": {
						Optional:    true,
						Description: "With a shared `hostname` such as provided by Amazon AWS, sets a unique cache key for your content.",
						Type:        schema.TypeBool,
					},
					"compress": {
						Optional:    true,
						Description: "Enables `gzip` compression for non-NetStorage origins.",
						Type:        schema.TypeBool,
					},
					"enable_true_client_ip": {
						Optional:    true,
						Description: "When enabled on non-NetStorage origins, allows you to send a custom header (the `trueClientIpHeader`) identifying the IP address of the immediate client connecting to the edge server. This may provide more useful information than the standard `X-Forward-For` header, which proxies may modify.",
						Type:        schema.TypeBool,
					},
					"true_client_ip_header": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "This specifies the name of the field that identifies the end client's IP address, for example `True-Client-IP`.",
						Type:             schema.TypeString,
					},
					"true_client_ip_client_setting": {
						Optional:    true,
						Description: "If a client sets the `True-Client-IP` header, the edge server allows it and passes the value to the origin. Otherwise the edge server removes it and sets the value itself.",
						Type:        schema.TypeBool,
					},
					"origin_certificate": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"verification_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PLATFORM_SETTINGS", "CUSTOM", "THIRD_PARTY"}, false)),
						Optional:         true,
						Description:      "For non-NetStorage origins, maximize security by controlling which certificates edge servers should trust.",
						Type:             schema.TypeString,
					},
					"origin_sni": {
						Optional:    true,
						Description: "For non-NetStorage origins, enabling this adds a Server Name Indication (SNI) header in the SSL request sent to the origin, with the origin hostname as the value. See the `verification settings in the Origin Server behavior` or contact your Akamai representative for more information. If you want to use TLS version 1.3 in your existing properties, enable this option. New properties have this enabled by default.",
						Type:        schema.TypeBool,
					},
					"custom_valid_cn_values": {
						Optional:    true,
						Description: "Specifies values to look for in the origin certificate's `Subject Alternate Name` or `Common Name` fields. Specify `{{Origin Hostname}}` and `{{Forward Host Header}}` within the text in the order you want them to be evaluated. (Note that these two template items are not the same as in-line `variables`, which use the same curly-brace syntax.)",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"origin_certs_to_honor": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COMBO", "STANDARD_CERTIFICATE_AUTHORITIES", "CUSTOM_CERTIFICATE_AUTHORITIES", "CUSTOM_CERTIFICATES"}, false)),
						Optional:         true,
						Description:      "Specifies which certificate to trust.",
						Type:             schema.TypeString,
					},
					"standard_certificate_authorities": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"custom_certificate_authorities": {
						Optional:    true,
						Description: "Specifies an array of certification objects. See the `verification settings in the Origin Server behavior` or contact your Akamai representative for details on this object's requirements.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"subject_cn": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"subject_alternative_names": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"subject_rdns": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"ou": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"o": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"cn": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"t": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"serialnumber": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"email_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dc": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"uid": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"street": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"surname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"givenname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"initials": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"generation": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"description": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"role": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unique_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dn_qualifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"pseudonym": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name_at_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_citizenship": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_residence": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"gender": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"place_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"date_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_code": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"business_category": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"telephone_number": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"organization_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
										},
									},
								},
								"issuer_rdns": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"ou": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"o": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"cn": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"t": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"serialnumber": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"email_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dc": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"uid": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"street": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"surname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"givenname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"initials": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"generation": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"description": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"role": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unique_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dn_qualifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"pseudonym": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name_at_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_citizenship": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_residence": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"gender": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"place_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"date_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_code": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"business_category": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"telephone_number": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"organization_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
										},
									},
								},
								"not_before": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"not_after": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"sig_alg_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key_algorithm": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key_format": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"serial_number": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"version": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"sha1_fingerprint": {
									ValidateDiagFunc: validateRegexOrVariable("^[a-f0-9]{40}$"),
									Optional:         true,
									Description:      "",
									Type:             schema.TypeString,
								},
								"pem_encoded_cert": {
									ValidateDiagFunc: validateRegexOrVariable("^-----BEGIN CERTIFICATE-----(.|\\s)*-----END CERTIFICATE-----\\s*$"),
									Optional:         true,
									Description:      "",
									Type:             schema.TypeString,
								},
								"can_be_leaf": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
								"can_be_ca": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
								"self_signed": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
							},
						},
					},
					"custom_certificates": {
						Optional:    true,
						Description: "Specifies an array of certification objects. See the `verification settings in the Origin Server behavior` or contact your Akamai representative for details on this object's requirements.",
						Type:        schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"subject_cn": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"subject_alternative_names": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"subject_rdns": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"ou": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"o": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"cn": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"t": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"serialnumber": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"email_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dc": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"uid": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"street": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"surname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"givenname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"initials": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"generation": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"description": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"role": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unique_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dn_qualifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"pseudonym": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name_at_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_citizenship": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_residence": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"gender": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"place_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"date_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_code": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"business_category": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"telephone_number": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"organization_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
										},
									},
								},
								"issuer_rdns": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"ou": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"o": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"cn": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"t": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"serialnumber": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"email_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dc": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"uid": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"street": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"surname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"givenname": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"initials": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"generation": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"description": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"role": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unstructured_name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"unique_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"dn_qualifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"pseudonym": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_address": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name_at_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_citizenship": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"country_of_residence": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"gender": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"place_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"date_of_birth": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"postal_code": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"business_category": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"telephone_number": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"name": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"organization_identifier": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_c": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_st": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
											"jurisdiction_l": {
												Optional:    true,
												Description: "",
												Type:        schema.TypeString,
											},
										},
									},
								},
								"not_before": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"not_after": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"sig_alg_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key_algorithm": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"public_key_format": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"serial_number": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"version": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"sha1_fingerprint": {
									ValidateDiagFunc: validateRegexOrVariable("^[a-f0-9]{40}$"),
									Optional:         true,
									Description:      "",
									Type:             schema.TypeString,
								},
								"pem_encoded_cert": {
									ValidateDiagFunc: validateRegexOrVariable("^-----BEGIN CERTIFICATE-----(.|\\s)*-----END CERTIFICATE-----\\s*$"),
									Optional:         true,
									Description:      "",
									Type:             schema.TypeString,
								},
								"can_be_leaf": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
								"can_be_ca": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
								"self_signed": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
							},
						},
					},
					"ports": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"http_port": {
						Optional:    true,
						Description: "Specifies the port on your origin server to which edge servers should connect for HTTP requests, customarily `80`.",
						Type:        schema.TypeInt,
					},
					"https_port": {
						Optional:    true,
						Description: "Specifies the port on your origin server to which edge servers should connect for secure HTTPS requests, customarily `443`. This option only applies if the property is marked as secure. See `Secure property requirements` for guidance.",
						Type:        schema.TypeInt,
					},
					"tls_version_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"min_tls_version": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validation.ToDiagFunc(validation.StringInSlice([]string{"DYNAMIC", "TLSV1_1", "TLSV1_2", "TLSV1_3"}, false))),
						Optional:         true,
						Description:      "Specifies the minimum TLS version to use for connections to your origin server.",
						Type:             schema.TypeString,
					},
					"max_tls_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DYNAMIC", "TLSV1_1", "TLSV1_2", "TLSV1_3"}, false)),
						Optional:         true,
						Description:      "Specifies the maximum TLS version to use for connections to your origin server. As best practice, use `DYNAMIC` to automatically apply the latest supported version.",
						Type:             schema.TypeString,
					},
					"http2_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"http2_enabled": {
						Optional:    true,
						Description: "Limited Availability. When enabled, the edge server sends multiple HTTP requests over a single HTTP/2 connection to the origin.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"origin_characteristics": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the origin. Akamai uses this information to optimize your metadata configuration, which may result in better origin offload and end-user performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"authentication_method_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"authentication_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AUTOMATIC", "SIGNATURE_HEADER_AUTHENTICATION", "MSL_AUTHENTICATION", "AWS", "GCS_HMAC_AUTHENTICATION", "AWS_STS"}, false)),
						Optional:         true,
						Description:      "Specifies the authentication method.",
						Type:             schema.TypeString,
					},
					"encoding_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{1, 2, 3, 4, 5})),
						Optional:         true,
						Description:      "Specifies the version of the encryption algorithm, an integer from `1` to `5`.",
						Type:             schema.TypeInt,
					},
					"use_custom_sign_string": {
						Optional:    true,
						Description: "Specifies whether to customize your signed string.",
						Type:        schema.TypeBool,
					},
					"custom_sign_string": {
						Optional:    true,
						Description: "Specifies the data to be encrypted as a series of enumerated variable names. See `Built-in system variables` for guidance on each.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"secret_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9a-zA-Z]{24}$")),
						Optional:         true,
						Description:      "Specifies the shared secret key.",
						Type:             schema.TypeString,
					},
					"nonce": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9a-zA-Z]{1,8}$"),
						Optional:         true,
						Description:      "Specifies the nonce.",
						Type:             schema.TypeString,
					},
					"mslkey": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9a-zA-Z]{10,}$"),
						Optional:         true,
						Description:      "Specifies the access key provided by the hosting service.",
						Type:             schema.TypeString,
					},
					"mslname": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9a-zA-Z]{1,8}$"),
						Optional:         true,
						Description:      "Specifies the origin name provided by the hosting service.",
						Type:             schema.TypeString,
					},
					"access_key_encrypted_storage": {
						Optional:    true,
						Description: "Enables secure use of access keys defined in Cloud Access Manager. Access keys store encrypted authentication details required to sign requests to cloud origins. If you disable this option, you'll need to store the authentication details unencrypted.",
						Type:        schema.TypeBool,
					},
					"gcs_access_key_version_guid": {
						Optional:    true,
						Description: "Identifies the unique `gcsAccessKeyVersionGuid` access key `created` in Cloud Access Manager to sign your requests to Google Cloud Storage in interoperability mode.",
						Type:        schema.TypeString,
					},
					"gcs_hmac_key_access_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9]{1,128}$"),
						Optional:         true,
						Description:      "Specifies the active access ID linked to your Google account.",
						Type:             schema.TypeString,
					},
					"gcs_hmac_key_secret": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9+/=_-]{1,40}$"),
						Optional:         true,
						Description:      "Specifies the secret linked to the access ID that you want to use to sign requests to Google Cloud Storage.",
						Type:             schema.TypeString,
					},
					"aws_access_key_version_guid": {
						Optional:    true,
						Description: "Identifies the unique `awsAccessKeyVersionGuid` access key `created` in Cloud Access Manager to sign your requests to AWS S3.",
						Type:        schema.TypeString,
					},
					"aws_access_key_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9]{1,128}$"),
						Optional:         true,
						Description:      "Specifies active access key ID linked to your AWS account.",
						Type:             schema.TypeString,
					},
					"aws_secret_access_key": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9+/=_-]{1,1024}$"),
						Optional:         true,
						Description:      "Specifies the secret linked to the access key identifier that you want to use to sign requests to AWS.",
						Type:             schema.TypeString,
					},
					"aws_region": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9-]+$"),
						Optional:         true,
						Description:      "This specifies the AWS region code of the location where your bucket resides.",
						Type:             schema.TypeString,
					},
					"aws_host": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(([a-zA-Z0-9]([a-zA-Z0-9_\\-]*[a-zA-Z0-9])?)\\.)+([a-zA-Z]+|xn--[a-zA-Z0-9]+)$")),
						Optional:         true,
						Description:      "This specifies the AWS hostname, without `http://` or `https://` prefixes. If you leave this option empty, it inherits the hostname from the `origin` behavior.",
						Type:             schema.TypeString,
					},
					"aws_service": {
						Optional:    true,
						Description: "This specifies the subdomain of your AWS service. It precedes `amazonaws.com` or the region code in the AWS hostname. For example, `s3.amazonaws.com`.",
						Type:        schema.TypeString,
					},
					"property_id_tag": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"hostname_tag": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"role_arn": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9][a-zA-Z0-9_\\+=,.@\\-:/]{0,2047}$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"aws_ar_region": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9][a-zA-Z0-9\\-]{0,63}$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"end_point_service": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[a-zA-Z0-9][a-zA-Z0-9\\-]{0,63}$")),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"origin_location_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"country": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EUROPE", "NORTH_AMERICA", "LATIN_AMERICA", "SOUTH_AMERICA", "NORDICS", "ASIA_PACIFIC", "OTHER_AMERICAS", "OTHER_APJ", "OTHER_EMEA", "AUSTRALIA", "GERMANY", "INDIA", "ITALY", "JAPAN", "MEXICO", "TAIWAN", "UNITED_KINGDOM", "US_EAST", "US_CENTRAL", "US_WEST", "GLOBAL_MULTI_GEO", "OTHER", "UNKNOWN", "ADC"}, false)),
						Optional:         true,
						Description:      "Specifies the origin's geographic region.",
						Type:             schema.TypeString,
					},
					"adc_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"direct_connect_geo": {
						Optional:    true,
						Description: "Provides a region used by Akamai Direct Connection.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"origin_characteristics_wsd": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies characteristics of the origin, for use in Akamai's Wholesale Delivery product. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origintype": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AZURE", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "Specifies an origin type.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"origin_failure_recovery_method": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Origin Failover requires that you set up a separate rule containing origin failure recovery methods. You also need to set up the Origin Failure Recovery Policy behavior in a separate rule with a desired match criteria, and select the desired failover method. You can do this using Property Manager. Learn more about this process in `Adaptive Media Delivery Implementation Guide`. You can use the `originFailureRecoveryPolicy` member to edit existing instances of the Origin Failure Recover Policy behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"recovery_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RETRY_ALTERNATE_ORIGIN", "RESPOND_CUSTOM_STATUS"}, false)),
						Optional:         true,
						Description:      "Specifies the recovery method.",
						Type:             schema.TypeString,
					},
					"custom_status_code": {
						Optional:    true,
						Description: "Specifies the custom status code to be sent to the client.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"origin_failure_recovery_policy": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Configures how to detect an origin failure, in which case the `originFailureRecoveryMethod` behavior applies. You can also define up to three sets of criteria to detect origin failure based on specific response codes. Use it to apply specific retry or recovery actions. You can do this using Property Manager. Learn more about this process in `Adaptive Media Delivery Implementation Guide`. You can use the `originFailureRecoveryMethod` member to edit existing instances of the Origin Failure Recover Method behavior. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Activates and configures a recovery policy.",
						Type:        schema.TypeBool,
					},
					"tuning_parameters": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_ip_avoidance": {
						Optional:    true,
						Description: "Temporarily blocks an origin IP address that experienced a certain number of failures. When an IP address is blocked, the `configName` established for `originResponsivenessRecoveryConfigName` is applied.",
						Type:        schema.TypeBool,
					},
					"ip_avoidance_error_threshold": {
						Optional:    true,
						Description: "Defines the number of failures that need to occur to an origin address before it's blocked.",
						Type:        schema.TypeInt,
					},
					"ip_avoidance_retry_interval": {
						Optional:    true,
						Description: "Defines the number of seconds after which the IP address is removed from the blocklist.",
						Type:        schema.TypeInt,
					},
					"binary_equivalent_content": {
						Optional:    true,
						Description: "Synchronizes content between the primary and backup origins, byte for byte.",
						Type:        schema.TypeBool,
					},
					"origin_responsiveness_monitoring": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"monitor_origin_responsiveness": {
						Optional:    true,
						Description: "Enables continuous monitoring of connectivity to the origin. If necessary, applies retry or recovery actions.",
						Type:        schema.TypeBool,
					},
					"origin_responsiveness_timeout": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AGGRESSIVE", "MODERATE", "CONSERVATIVE", "USER_SPECIFIED"}, false)),
						Optional:         true,
						Description:      "The timeout threshold that triggers a retry or recovery action.",
						Type:             schema.TypeString,
					},
					"origin_responsiveness_custom_timeout": {
						Optional:    true,
						Description: "Specify a custom timeout, from 1 to 10 seconds.",
						Type:        schema.TypeInt,
					},
					"origin_responsiveness_enable_retry": {
						Optional:    true,
						Description: "If a specific failure condition applies, attempts a retry on the same origin before executing the recovery method.",
						Type:        schema.TypeBool,
					},
					"origin_responsiveness_enable_recovery": {
						Optional:    true,
						Description: "Enables a recovery action for a specific failure condition.",
						Type:        schema.TypeBool,
					},
					"origin_responsiveness_recovery_config_name": {
						Optional:    true,
						Description: "Specifies a recovery configuration using the `configName` you defined in the `recoveryConfig` match criteria. Specify 3 to 20 alphanumeric characters or dashes. Ensure that you use the `recoveryConfig` match criteria to apply this option.",
						Type:        schema.TypeString,
					},
					"status_code_monitoring1": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"monitor_status_codes1": {
						Optional:    true,
						Description: "Enables continuous monitoring for the specific origin status codes that trigger retry or recovery actions.",
						Type:        schema.TypeBool,
					},
					"monitor_response_codes1": {
						Optional:    true,
						Description: "Defines the origin response codes that trigger a subsequent retry or recovery action. Specify a single code entry (`501`) or a range (`501:504`). If you configure other `monitorStatusCodes*` and `monitorResponseCodes*` options, you can't use the same codes here.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"monitor_status_codes1_enable_retry": {
						Optional:    true,
						Description: "When the defined response codes apply, attempts a retry on the same origin before executing the recovery method.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes1_enable_recovery": {
						Optional:    true,
						Description: "Enables the recovery action for the response codes you define.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes1_recovery_config_name": {
						Optional:    true,
						Description: "Specifies a recovery configuration using the `configName` you defined in the `recoveryConfig` match criteria. Specify 3 to 20 alphanumeric characters or dashes. Ensure that you use the `recoveryConfig` match criteria to apply this option.",
						Type:        schema.TypeString,
					},
					"status_code_monitoring2": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"monitor_status_codes2": {
						Optional:    true,
						Description: "Enables continuous monitoring for the specific origin status codes that trigger retry or recovery actions.",
						Type:        schema.TypeBool,
					},
					"monitor_response_codes2": {
						Optional:    true,
						Description: "Defines the origin response codes that trigger a subsequent retry or recovery action. Specify a single code entry (`501`) or a range (`501:504`). If you configure other `monitorStatusCodes*` and `monitorResponseCodes*` options, you can't use the same codes here.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"monitor_status_codes2_enable_retry": {
						Optional:    true,
						Description: "When the defined response codes apply, attempts a retry on the same origin before executing the recovery method.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes2_enable_recovery": {
						Optional:    true,
						Description: "Enables the recovery action for the response codes you define.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes2_recovery_config_name": {
						Optional:    true,
						Description: "Specifies a recovery configuration using the `configName` you defined in the `recoveryConfig` match criteria. Specify 3 to 20 alphanumeric characters or dashes. Ensure that you use the `recoveryConfig` match criteria to apply this option.",
						Type:        schema.TypeString,
					},
					"status_code_monitoring3": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"monitor_status_codes3": {
						Optional:    true,
						Description: "Enables continuous monitoring for the specific origin status codes that trigger retry or recovery actions.",
						Type:        schema.TypeBool,
					},
					"monitor_response_codes3": {
						Optional:    true,
						Description: "Defines the origin response codes that trigger a subsequent retry or recovery action. Specify a single code entry (`501`) or a range (`501:504`). If you configure other `monitorStatusCodes*` and `monitorResponseCodes*` options, you can't use the same codes here..",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"monitor_status_codes3_enable_retry": {
						Optional:    true,
						Description: "When the defined response codes apply, attempts a retry on the same origin before executing the recovery method.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes3_enable_recovery": {
						Optional:    true,
						Description: "Enables the recovery action for the response codes you define.",
						Type:        schema.TypeBool,
					},
					"monitor_status_codes3_recovery_config_name": {
						Optional:    true,
						Description: "Specifies a recovery configuration using the `configName` you defined in the `recoveryConfig` match criteria. Specify 3 to 20 alphanumeric characters or dashes. Ensure that you use the `recoveryConfig` match criteria to apply this option.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"origin_ip_acl": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Origin IP Access Control List limits the traffic to your origin. It only allows requests from specific edge servers that are configured as part of a supernet defined by CIDR blocks. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "Enables the Origin IP Access Control List behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"permissions_policy": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Manages whether your page and its embedded iframes can access various browser features that affect end-user privacy, security, and performance. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"permissions_policy_directive": {
						Optional:    true,
						Description: "Each directive represents a browser feature. Specify the ones you want enabled in a client browser that accesses your content. You can add custom entries or provide pre-set values from the list. For more details on each value, see the `guide section` for this behavior.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"allow_list": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*:%\\[\\]@.\\s]+$"),
						Optional:         true,
						Description:      "The features you've set in `permissionsPolicyDirective` are enabled for domains you specify here. They'll remain disabled for all other domains. Separate multiple domains with a single space. To block the specified directives from all domains, set this to `none`. This generates an empty value in the `Permissions-Policy` header.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"persistent_client_connection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior activates `persistent connections` between edge servers and clients, which allow for better performance and more efficient use of resources. Compare with the `persistentConnection` behavior, which configures persistent connections for the entire journey from origin to edge to client.  Contact Akamai Professional Services for help configuring either. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the persistent connections behavior.",
						Type:        schema.TypeBool,
					},
					"timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the timeout period after which edge server closes the persistent connection with the client, 500 seconds by default.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"persistent_connection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior enables more efficient `persistent connections` from origin to edge server to client. Compare with the `persistentClientConnection` behavior, which customizes persistent connections from edge to client. Contact Akamai Professional Services for help configuring either. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables persistent connections.",
						Type:        schema.TypeBool,
					},
					"timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the timeout period after which edge server closes a persistent connection.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"personally_identifiable_information": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Marks content covered by the current rule as sensitive `personally identifiable information` that needs to be treated as secure and private. That includes anything involving personal information: name, social security number, date and place of birth, mother's maiden name, biometric data, or any other data linked to an individual. If you attempt to save a property with such a rule that also caches or logs sensitive content, the added behavior results in a validation error. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, marks content as personally identifiable information (PII).",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"phased_release": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Phased Release Cloudlet provides gradual and granular traffic management to an alternate origin in near real time.  Use the `Cloudlets API` or the Cloudlets Policy Manager application within `Control Center` to set up your Cloudlets policies. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Phased Release Cloudlet.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to apply the Cloudlet shared policy to an unlimited number of properties within your account. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Specifies the Cloudlet policy as an object.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet shared policy to use with this behavior. Use the `Cloudlets API` to list available shared policies.",
						Type:        schema.TypeInt,
					},
					"label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "A label to distinguish this Phased Release policy from any others within the same property.",
						Type:             schema.TypeString,
					},
					"population_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"population_cookie_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "NEVER", "ON_BROWSER_CLOSE", "FIXED_DATE", "DURATION"}, false)),
						Optional:         true,
						Description:      "Select when to assign a cookie to the population of users the Cloudlet defines. If you select the Cloudlet's `random` membership option, it overrides this option's value so that it is effectively `NONE`.",
						Type:             schema.TypeString,
					},
					"population_expiration_date": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the date and time when membership expires, and the browser no longer sends the cookie. Subsequent requests re-evaluate based on current membership settings.",
						Type:             schema.TypeString,
					},
					"population_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Sets the lifetime of the cookie from the initial request. Subsequent requests re-evaluate based on current membership settings.",
						Type:             schema.TypeString,
					},
					"population_refresh": {
						Optional:    true,
						Description: "Enabling this option resets the original duration of the cookie if the browser refreshes before the cookie expires.",
						Type:        schema.TypeBool,
					},
					"failover_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"failover_enabled": {
						Optional:    true,
						Description: "Allows failure responses at the origin defined by the Cloudlet to fail over to the prevailing origin defined by the property.",
						Type:        schema.TypeBool,
					},
					"failover_response_code": {
						Optional:    true,
						Description: "Defines the set of failure codes that initiate the failover response.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"failover_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 300)),
						Optional:         true,
						Description:      "Specifies the number of seconds to wait until the client tries to access the failover origin after the initial failure is detected. Set the value to `0` to immediately request the alternate origin upon failure.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"preconnect": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "With the `http2` behavior enabled, this requests a specified set of domains that relate to your property hostname, and keeps the connection open for faster loading of content from those domains. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"preconnectlist": {
						Optional:    true,
						Description: "Specifies the set of hostnames to which to preconnect over HTTP2.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"predictive_content_delivery": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Improves user experience and reduces the cost of downloads by enabling mobile devices to predictively fetch and cache content from catalogs managed by Akamai servers. You can't use this feature if in the `segmentedMediaOptimization` behavior, the value for `behavior` is set to `LIVE`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the predictive content delivery behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"predictive_prefetching": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior potentially reduces the client's page load time by pre-caching objects based on historical data for the page, not just its current set of referenced objects. It also detects second-level dependencies, such as objects retrieved by JavaScript. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the predictive prefetching behavior.",
						Type:        schema.TypeBool,
					},
					"accuracy_target": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LOW", "MEDIUM", "HIGH"}, false)),
						Optional:         true,
						Description:      "The level of prefetching. A higher level results in better client performance, but potentially greater load on the origin.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"prefetch": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Instructs edge servers to retrieve content linked from requested pages as they load, rather than waiting for separate requests for the linked content. This behavior applies depending on the rule's set of matching conditions. Use in conjunction with the `prefetchable` behavior, which specifies the set of objects to prefetch. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Applies prefetching behavior when enabled.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"prefetchable": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allow matching objects to prefetch into the edge cache as the parent page that links to them loads, rather than waiting for a direct request. This behavior applies depending on the rule's set of matching conditions. Use `prefetch` to enable the overall behavior for parent pages that contain links to the object. To apply this behavior, you need to match on a `filename` or `fileExtension`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows matching content to prefetch when referenced on a requested parent page.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"prefresh_cache": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Refresh cached content before its time-to-live (TTL) expires, to keep end users from having to wait for the origin to provide fresh content. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the cache prefreshing behavior.",
						Type:        schema.TypeBool,
					},
					"prefreshval": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 99)),
						Optional:         true,
						Description:      "Specifies when the prefresh occurs as a percentage of the TTL. For example, for an object whose cache has 10 minutes left to live, and an origin response that is routinely less than 30 seconds, a percentage of `95` prefreshes the content without unnecessarily increasing load on the origin.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"quality": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origin_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"country": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EUROPE", "NORTH_AMERICA", "LATIN_AMERICA", "SOUTH_AMERICA", "NORDICS", "ASIA_PACIFIC", "OTHER_AMERICAS", "OTHER_APJ", "OTHER_EMEA", "AUSTRALIA", "GERMANY", "INDIA", "ITALY", "JAPAN", "MEXICO", "TAIWAN", "UNITED_KINGDOM", "US_EAST", "US_CENTRAL", "US_WEST"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"audience_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"end_user_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GLOBAL", "GLOBAL_US_CENTRIC", "GLOBAL_EU_CENTRIC", "GLOBAL_ASIA_CENTRIC", "EUROPE", "NORTH_AMERICA", "SOUTH_AMERICA", "NORDICS", "ASIA_PACIFIC", "AUSTRALIA", "GERMANY", "INDIA", "ITALY", "JAPAN", "TAIWAN", "UNITED_KINGDOM"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"maximum_concurrent_users": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "LESS_THAN_10K", "10K_TO_50K", "50K_TO_100K", "GREATER_THAN_100K"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"content_settings": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"content_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "SITE", "IMAGES", "CONFIG", "OTHERS", "AUDIO", "SD_VIDEO", "HD_VIDEO", "SUPER_HD_VIDEO", "LARGE_OBJECTS"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"object_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"LESS_THAN_1MB", "1_TO_10MB", "10_TO_100MB", "GREATER_THAN_100MB"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"download_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"FOREGROUND", "BACKGROUND"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"popularity_distribution": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"TYPICAL", "LONG_TAIL", "ALL_POPULAR", "ALL_UNPOPULAR"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"delivery_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ON_DEMAND", "LIVE"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"delivery_format": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DASH", "HDS", "HLS", "SILVER_LIGHT", "OTHER"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"segment_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{2, 4, 6, 8, 10})),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeInt,
					},
					"catalog_size": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE", "EXTRA_LARGE"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"refresh_rate": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "HOURLY", "DAILY", "MONTHLY", "YEARLY"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"optimize_for": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "ORIGIN", "STARTUP"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"quic_beta": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables QUIC support.",
						Type:        schema.TypeBool,
					},
					"quic_offer_percentage": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 50)),
						Optional:         true,
						Description:      "The percentage of responses for which to allow QUIC sessions.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"random_seek": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Optimizes `.flv` and `.mp4` files to allow random jump-point navigation. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"flv": {
						Optional:    true,
						Description: "Enables random seek optimization in FLV files.",
						Type:        schema.TypeBool,
					},
					"mp4": {
						Optional:    true,
						Description: "Enables random seek optimization in MP4 files.",
						Type:        schema.TypeBool,
					},
					"maximum_size": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d+[K,M,G,T]B$"),
						Optional:         true,
						Description:      "Sets the maximum size of the MP4 file to optimize, expressed as a number suffixed with a unit string such as `MB` or `GB`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"rapid": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Akamai API Gateway` allows you to configure API traffic delivered over the Akamai network. Apply this behavior to a set of API assets, then use Akamai's `API Endpoints API` to configure how the traffic responds.  Use the `API Keys and Traffic Management API` to control access to your APIs. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables API Gateway for the current set of content.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"read_timeout": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior specifies how long the edge server should wait for a response from the requesting forward server after a connection has already been established. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "The amount of time an edge server should wait for each read statement to return a response from the forward server after a connection has already been established. Larger objects may need many reads, and this timeout applies to each read separately. Any failure to complete a read within this time limit aborts the request and sends a 504 Gateway Timeout error to the client.",
						Type:             schema.TypeString,
					},
					"first_byte_timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "The amount of time an edge server should wait for the first byte of the response to be returned from the forward server after a connection has already been established. Instead of continually waiting for the content, edge servers send a 504 Gateway Timeout error to the client. If your origin server is handling high loads and might respond slowly, specify a short timeout. Defaults to 20 seconds. The value for First Byte Timeout can't be 0 and it can't exceed 10 minutes (600 seconds).",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"real_time_reporting": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This enables Real-Time Reporting for Akamai `Cloud Embed` customers. The behavior can only be configured on your behalf by Akamai Professional Services. You can access real-time reports data for that base configuration with `Media Delivery Reports API`. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables reports on delivery of cloud hosted content at near real-time latencies.",
						Type:        schema.TypeBool,
					},
					"advanced": {
						Optional:    true,
						Description: "Enables advanced options.",
						Type:        schema.TypeBool,
					},
					"beacon_sampling_percentage": {
						Optional:    true,
						Description: "Specifies the percentage for sampling.",
						Type:        schema.TypeFloat,
					},
				},
			},
		},
		"real_user_monitoring": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, activates real-use monitoring.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"redirect": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Respond to the client request with a redirect without contacting the origin. Specify the redirect as a path expression starting with a `/` character relative to the current root, or as a fully qualified URL. This behavior relies primarily on `destinationHostname` and `destinationPath` to manipulate the hostname and path independently. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"mobile_default_choice": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DEFAULT", "MOBILE"}, false)),
						Optional:         true,
						Description:      "Either specify a default response for mobile browsers, or customize your own.",
						Type:             schema.TypeString,
					},
					"destination_protocol": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SAME_AS_REQUEST", "HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "Choose the protocol for the redirect URL.",
						Type:             schema.TypeString,
					},
					"destination_hostname": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SAME_AS_REQUEST", "SUBDOMAIN", "SIBLING", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specify how to change the requested hostname, independently from the pathname.",
						Type:             schema.TypeString,
					},
					"destination_hostname_subdomain": {
						Optional:    true,
						Description: "Specifies a subdomain to prepend to the current hostname. For example, a value of `m` changes `www.example.com` to `m.www.example.com`.",
						Type:        schema.TypeString,
					},
					"destination_hostname_sibling": {
						Optional:    true,
						Description: "Specifies the subdomain with which to replace to the current hostname's leftmost subdomain. For example, a value of `m` changes `www.example.com` to `m.example.com`.",
						Type:        schema.TypeString,
					},
					"destination_hostname_other": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "Specifies the full hostname with which to replace the current hostname.",
						Type:             schema.TypeString,
					},
					"destination_path": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SAME_AS_REQUEST", "PREFIX_REQUEST", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specify how to change the requested pathname, independently from the hostname.",
						Type:             schema.TypeString,
					},
					"destination_path_prefix": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "When `destinationPath` is set to `PREFIX_REQUEST`, this prepends the current path. For example, a value of `/prefix/path` changes `/example/index.html` to `/prefix/path/example/index.html`.",
						Type:             schema.TypeString,
					},
					"destination_path_suffix_status": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NO_SUFFIX", "SUFFIX"}, false)),
						Optional:         true,
						Description:      "When `destinationPath` is set to `PREFIX_REQUEST`, this gives you the option of adding a suffix.",
						Type:             schema.TypeString,
					},
					"destination_path_suffix": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9\\[\\]/?#!=&_\\-\\.]+$"),
						Optional:         true,
						Description:      "When `destinationPath` is set to `PREFIX_REQUEST` and `destinationPathSuffixStatus` is set to `SUFFIX`, this specifies the suffix to append to the path.",
						Type:             schema.TypeString,
					},
					"destination_path_other": {
						ValidateDiagFunc: validateRegexOrVariable("^/"),
						Optional:         true,
						Description:      "When `destinationPath` is set to `PREFIX_REQUEST`, this replaces the current path.",
						Type:             schema.TypeString,
					},
					"query_string": {
						Optional:    true,
						Description: "When set to `APPEND`, passes incoming query string parameters as part of the redirect URL. Otherwise set this to `IGNORE`.",
						Type:        schema.TypeString,
					},
					"response_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302, 303, 307})),
						Optional:         true,
						Description:      "Specify the redirect's response code.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"redirectplus": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Respond to the client request with a redirect without contacting the origin. This behavior fills the same need as `redirect`, but allows you to use `variables` to express the redirect `destination`'s component values more concisely. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the redirect feature.",
						Type:        schema.TypeBool,
					},
					"destination": {
						Optional:    true,
						Description: "Specifies the redirect as a path expression starting with a `/` character relative to the current root, or as a fully qualified URL. Optionally inject variables, as in this example that refers to the original request's filename: `/path/to/{{builtin.AK_FILENAME}}`.",
						Type:        schema.TypeString,
					},
					"response_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302, 303, 307})),
						Optional:         true,
						Description:      "Assigns the status code for the redirect response.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"referer_checking": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Limits allowed requests to a set of domains you specify. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the referer-checking behavior.",
						Type:        schema.TypeBool,
					},
					"strict": {
						Optional:    true,
						Description: "When enabled, excludes requests whose `Referer` header include a relative path, or that are missing a `Referer`. When disabled, only excludes requests whose `Referer` hostname is not part of the `domains` set.",
						Type:        schema.TypeBool,
					},
					"domains": {
						Optional:    true,
						Description: "Specifies the set of allowed domains. With `allowChildren` disabled, prefixing values with `*.` specifies domains for which subdomains are allowed.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"allow_children": {
						Optional:    true,
						Description: "Allows all subdomains for the `domains` set, just like adding a `*.` prefix to each.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"remove_query_parameter": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Remove named query parameters before forwarding the request to the origin. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"parameters": {
						Optional:    true,
						Description: "Specifies parameters to remove from the request.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"remove_vary": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "By default, responses that feature a `Vary` header value of anything other than `Accept-Encoding` and a corresponding `Content-Encoding: gzip` header aren't cached on edge servers. `Vary` headers indicate when a URL's content varies depending on some variable, such as which `User-Agent` requests it. This behavior simply removes the `Vary` header to make responses cacheable. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, removes the `Vary` header to ensure objects can be cached.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"report": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specify the HTTP request headers or cookie names to log in your Log Delivery Service reports. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"log_host": {
						Optional:    true,
						Description: "Log the `Host` header.",
						Type:        schema.TypeBool,
					},
					"log_referer": {
						Optional:    true,
						Description: "Log the `Referer` header.",
						Type:        schema.TypeBool,
					},
					"log_user_agent": {
						Optional:    true,
						Description: "Log the `User-Agent` header.",
						Type:        schema.TypeBool,
					},
					"log_accept_language": {
						Optional:    true,
						Description: "Log the `Accept-Language` header.",
						Type:        schema.TypeBool,
					},
					"log_cookies": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"OFF", "ALL", "SOME"}, false)),
						Optional:         true,
						Description:      "Specifies the set of cookies to log.",
						Type:             schema.TypeString,
					},
					"cookies": {
						Optional:    true,
						Description: "This specifies the set of cookies names whose values you want to log.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"log_custom_log_field": {
						Optional:    true,
						Description: "Whether to append additional custom data to each log line.",
						Type:        schema.TypeBool,
					},
					"custom_log_field": {
						Optional:    true,
						Description: "Specifies an additional data field to append to each log line, maximum 1000 bytes, typically based on a dynamically generated built-in system variable. For example, `round-trip: {{builtin.AK_CLIENT_TURNAROUND_TIME}}ms` logs the total time to complete the response. See `Support for variables` for more information. If you enable the `logCustom` behavior, it overrides the `customLogField` option.",
						Type:        schema.TypeString,
					},
					"log_edge_ip": {
						Optional:    true,
						Description: "Whether to log the IP address of the Akamai edge server that served the response to the client.",
						Type:        schema.TypeBool,
					},
					"log_x_forwarded_for": {
						Optional:    true,
						Description: "Log any `X-Forwarded-For` request header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"request_client_hints": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Client hints are HTTP request header fields that determine which resources the browser should include in the response. This behavior configures and prioritizes the client hints you want to send to request specific client and device information. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"accept_ch": {
						Optional:    true,
						Description: "The client hint data objects you want to receive from the browser.  You can add custom entries or provide pre-set values from the list. For more details on each value, see the `guide section` for this behavior. If you've configured your origin server to pass along data objects, they merge with the ones you set in this array, before the list is sent to the client.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"accept_critical_ch": {
						Optional:    true,
						Description: "The critical client hint data objects you want to receive from the browser. The original request from the browser needs to include these objects. Otherwise, a new response header is sent back to the client, asking for all of these client hint data objects. You can add custom entries or provide pre-set values from the list. For more details on each value, see the `guide section` for this behavior.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"reset": {
						Optional:    true,
						Description: "This sends an empty instance of the `Accept-CH` response header to clear other `Accept-CH` values currently stored in the client browser. This empty header doesn't get merged with other objects sent from your origin server.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"request_control": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The Request Control Cloudlet allows you to control access to your web content based on the incoming request's IP or geographic location.  With Cloudlets available on your contract, choose `Your services` > `Edge logic Cloudlets` to control how the feature works within `Control Center`, or use the `Cloudlets API` to configure it programmatically. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Request Control Cloudlet.",
						Type:        schema.TypeBool,
					},
					"is_shared_policy": {
						Optional:    true,
						Description: "Whether you want to apply the Cloudlet shared policy to an unlimited number of properties within your account. Learn more about shared policies and how to create them in `Cloudlets Policy Manager`.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet shared policy to use with this behavior. Use the `Cloudlets API` to list available shared policies.",
						Type:        schema.TypeInt,
					},
					"enable_branded403": {
						Optional:    true,
						Description: "If enabled, serves a branded 403 page for this Cloudlet instance.",
						Type:        schema.TypeBool,
					},
					"branded403_status_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{200, 302, 403, 503})),
						Optional:         true,
						Description:      "Specifies the response status code for the branded deny action.",
						Type:             schema.TypeInt,
					},
					"net_storage": {
						Optional:    true,
						Description: "Specifies the NetStorage domain that contains the branded 403 page.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"branded403_file": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the full path of the branded 403 page, including the filename, but excluding the NetStorage CP code path component.",
						Type:             schema.TypeString,
					},
					"branded403_url": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\s]+$"),
						Optional:         true,
						Description:      "Specifies the redirect URL for the branded deny action.",
						Type:             schema.TypeString,
					},
					"branded_deny_cache_ttl": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 30)),
						Optional:         true,
						Description:      "Specifies the branded response page's time to live in the cache, `5` minutes by default.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"request_type_marker": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Internet of Things: OTA Updates` product allows customers to securely distribute firmware to devices over cellular networks. When using the `downloadCompleteMarker` behavior to log successful downloads, this related behavior identifies download or campaign server types in aggregated and individual reports. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"request_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DOWNLOAD", "CAMPAIGN_SERVER"}, false)),
						Optional:         true,
						Description:      "Specifies the type of request.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"resource_optimizer": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Resource Optimizer feature.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"resource_optimizer_extended_compatibility": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This enhances the standard version of the `resourceOptimizer` behavior to support the compression of additional file formats and address some compatibility issues. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Resource Optimizer feature.",
						Type:        schema.TypeBool,
					},
					"enable_all_features": {
						Optional:    true,
						Description: "Enables `additional support` and error handling.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"response_code": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Change the existing response code. For example, if your origin sends a `301` permanent redirect, this behavior can change it on the edge to a temporary `302` redirect. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"status_code": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{200, 301, 302, 303, 404, 500, 100, 101, 102, 103, 122, 201, 202, 203, 204, 205, 206, 207, 226, 300, 304, 305, 306, 307, 308, 400, 401, 402, 403, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 422, 423, 424, 425, 426, 428, 429, 431, 444, 449, 450, 499, 501, 502, 503, 504, 505, 506, 507, 509, 510, 511, 598, 599})),
						Optional:         true,
						Description:      "The HTTP status code to replace the existing one.",
						Type:             schema.TypeInt,
					},
					"override206": {
						Optional:    true,
						Description: "Allows any specified `200` success code to override a `206` partial-content code, in which case the response's content length matches the requested range length.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"response_cookie": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Set a cookie to send downstream to the client with either a fixed value or a unique stamp. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name of the cookie, which serves as a key to determine if the cookie is set.",
						Type:             schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows you to set a response cookie.",
						Type:        schema.TypeBool,
					},
					"type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"FIXED", "UNIQUE"}, false)),
						Optional:         true,
						Description:      "What type of value to assign.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\s;]+$"),
						Optional:         true,
						Description:      "If the cookie `type` is `FIXED`, this specifies the cookie value.",
						Type:             schema.TypeString,
					},
					"format": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"AKAMAI", "APACHE"}, false)),
						Optional:         true,
						Description:      "When the `type` of cookie is set to `UNIQUE`, this sets the date format.",
						Type:             schema.TypeString,
					},
					"default_domain": {
						Optional:    true,
						Description: "When enabled, uses the default domain value, otherwise the set specified in the `domain` field.",
						Type:        schema.TypeBool,
					},
					"default_path": {
						Optional:    true,
						Description: "When enabled, uses the default path value, otherwise the set specified in the `path` field.",
						Type:        schema.TypeBool,
					},
					"domain": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "If the `defaultDomain` is disabled, this sets the domain for which the cookie is valid. For example, `example.com` makes the cookie valid for that hostname and all subdomains.",
						Type:             schema.TypeString,
					},
					"path": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "If the `defaultPath` is disabled, sets the path component for which the cookie is valid.",
						Type:             schema.TypeString,
					},
					"expires": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ON_BROWSER_CLOSE", "FIXED_DATE", "DURATION", "NEVER"}, false)),
						Optional:         true,
						Description:      "Sets various ways to specify when the cookie expires.",
						Type:             schema.TypeString,
					},
					"expiration_date": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "If `expires` is set to `FIXED_DATE`, this sets when the cookie expires as a UTC date and time.",
						Type:             schema.TypeString,
					},
					"duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "If `expires` is set to `DURATION`, this sets the cookie's lifetime.",
						Type:             schema.TypeString,
					},
					"same_site": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DEFAULT", "NONE", "LAX", "STRICT"}, false)),
						Optional:         true,
						Description:      "This option controls the `SameSite` cookie attribute that reduces the risk of cross-site request forgery attacks.",
						Type:             schema.TypeString,
					},
					"secure": {
						Optional:    true,
						Description: "When enabled, sets the cookie's `Secure` flag to transmit it with `HTTPS`.",
						Type:        schema.TypeBool,
					},
					"http_only": {
						Optional:    true,
						Description: "When enabled, includes the `HttpOnly` attribute in the `Set-Cookie` response header to mitigate the risk of client-side scripts accessing the protected cookie, if the browser supports it.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"restrict_object_caching": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "You need this behavior to deploy the Object Caching product. It disables serving HTML content and limits the maximum object size to 100MB. Contact Akamai Professional Services for help configuring it. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"maximum_size": {
						Optional:    true,
						Description: "Specifies a fixed maximum size of non-HTML content to cache.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"return_cache_status": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Generates a response header with information about cache status. Among other things, this can tell you whether the response came from the Akamai cache, or from the origin. Status values report with either of these forms of syntax, depending for example on whether you're deploying traffic using `sureRoute` or `tieredDistribution`: This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"response_header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "Specifies the name of the HTTP header in which to report the cache status value.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"rewrite_url": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modifies the path of incoming requests to forward to the origin. This helps you offload URL-rewriting tasks to the edge to increase the origin server's performance, allows you to redirect links to different targets without changing markup, and hides your original directory structure. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REPLACE", "REMOVE", "REWRITE", "PREPEND", "REGEX_REPLACE"}, false)),
						Optional:         true,
						Description:      "The action to perform on the path.",
						Type:             schema.TypeString,
					},
					"match": {
						ValidateDiagFunc: validateRegexOrVariable("^/([^:#\\[\\]@/?]+/)*$"),
						Optional:         true,
						Description:      "When `behavior` is `REMOVE` or `REPLACE`, specifies the part of the incoming path you'd like to remove or modify.",
						Type:             schema.TypeString,
					},
					"match_regex": {
						Optional:    true,
						Description: "When `behavior` is set to `REGEX_REPLACE`, specifies the Perl-compatible regular expression to replace with `targetRegex`.",
						Type:        schema.TypeString,
					},
					"target_regex": {
						Optional:    true,
						Description: "When `behavior` is set to `REGEX_REPLACE`, this replaces whatever the `matchRegex` field matches, along with any captured sequences from `\\$1` through `\\$9`.",
						Type:        schema.TypeString,
					},
					"target_path": {
						ValidateDiagFunc: validateRegexOrVariable("^/([^:#\\[\\]@/?]+/)*$"),
						Optional:         true,
						Description:      "When `behavior` is set to `REPLACE`, this path replaces whatever the `match` field matches in the incoming request's path.",
						Type:             schema.TypeString,
					},
					"target_path_prepend": {
						ValidateDiagFunc: validateRegexOrVariable("^/([^:#\\[\\]@/?]+/)*$"),
						Optional:         true,
						Description:      "When `behavior` is set to `PREPEND`, specifies a path to prepend to the incoming request's URL.",
						Type:             schema.TypeString,
					},
					"target_url": {
						ValidateDiagFunc: validateRegexOrVariable("(/\\S*)?$"),
						Optional:         true,
						Description:      "When `behavior` is set to `REWRITE`, specifies the full path to request from the origin.",
						Type:             schema.TypeString,
					},
					"match_multiple": {
						Optional:    true,
						Description: "When enabled, replaces all potential matches rather than only the first.",
						Type:        schema.TypeBool,
					},
					"keep_query_string": {
						Optional:    true,
						Description: "When enabled, retains the original path's query parameters.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"rum_custom": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"rum_sample_rate": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the percentage of web traffic to include in your RUM report.",
						Type:             schema.TypeInt,
					},
					"rum_group_name": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^[0-9a-zA-Z]*$")),
						Optional:         true,
						Description:      "A deprecated option to specify an alternate name under which to batch this set of web traffic in your report. Do not use it.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"saas_definitions": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Configures how the Software as a Service feature identifies `customers`, `applications`, and `users`. A different set of options is available for each type of targeted request, each enabled with the `action`-suffixed option. In each case, you can use `PATH`, `COOKIE`, `QUERY_STRING`, or `HOSTNAME` components as identifiers, or `disable` the SaaS behavior for certain targets. If you rely on a `HOSTNAME`, you also have the option of specifying a `CNAME chain` rather than an individual hostname. The various options suffixed `regex` and `replace` subsequently remove the identifier from the request. This behavior requires a sibling `origin` behavior whose `originType` option is set to `SAAS_DYNAMIC_ORIGIN`. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"customer_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"customer_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DISABLED", "HOSTNAME", "PATH", "QUERY_STRING", "COOKIE"}, false)),
						Optional:         true,
						Description:      "Specifies the request component that identifies a SaaS customer.",
						Type:             schema.TypeString,
					},
					"customer_cname_enabled": {
						Optional:    true,
						Description: "Enabling this allows you to identify customers using a `CNAME chain` rather than a single hostname.",
						Type:        schema.TypeBool,
					},
					"customer_cname_level": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the number of CNAMEs to use in the chain.",
						Type:             schema.TypeInt,
					},
					"customer_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This specifies the name of the cookie that identifies the customer.",
						Type:             schema.TypeString,
					},
					"customer_query_string": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "This names the query parameter that identifies the customer.",
						Type:             schema.TypeString,
					},
					"customer_regex": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9\\:\\[\\]\\{\\}\\(\\)\\.\\?_\\-\\*\\+\\^\\$\\\\\\/\\|&=!]{1,250})$"),
						Optional:         true,
						Description:      "Specifies a Perl-compatible regular expression with which to substitute the request's customer ID.",
						Type:             schema.TypeString,
					},
					"customer_replace": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9_\\-]|\\$[1-9]){1,250})$"),
						Optional:         true,
						Description:      "Specifies a string to replace the request's customer ID matched by `customerRegex`.",
						Type:             schema.TypeString,
					},
					"application_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"application_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DISABLED", "HOSTNAME", "PATH", "QUERY_STRING", "COOKIE"}, false)),
						Optional:         true,
						Description:      "Specifies the request component that identifies a SaaS application.",
						Type:             schema.TypeString,
					},
					"application_cname_enabled": {
						Optional:    true,
						Description: "Enabling this allows you to identify applications using a `CNAME chain` rather than a single hostname.",
						Type:        schema.TypeBool,
					},
					"application_cname_level": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the number of CNAMEs to use in the chain.",
						Type:             schema.TypeInt,
					},
					"application_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This specifies the name of the cookie that identifies the application.",
						Type:             schema.TypeString,
					},
					"application_query_string": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "This names the query parameter that identifies the application.",
						Type:             schema.TypeString,
					},
					"application_regex": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9\\:\\[\\]\\{\\}\\(\\)\\.\\?_\\-\\*\\+\\^\\$\\\\\\/\\|&=!]{1,250})$"),
						Optional:         true,
						Description:      "Specifies a Perl-compatible regular expression with which to substitute the request's application ID.",
						Type:             schema.TypeString,
					},
					"application_replace": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9_\\-]|\\$[1-9]){1,250})$"),
						Optional:         true,
						Description:      "Specifies a string to replace the request's application ID matched by `applicationRegex`.",
						Type:             schema.TypeString,
					},
					"users_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"users_action": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DISABLED", "HOSTNAME", "PATH", "QUERY_STRING", "COOKIE"}, false)),
						Optional:         true,
						Description:      "Specifies the request component that identifies a SaaS user.",
						Type:             schema.TypeString,
					},
					"users_cname_enabled": {
						Optional:    true,
						Description: "Enabling this allows you to identify users using a `CNAME chain` rather than a single hostname.",
						Type:        schema.TypeBool,
					},
					"users_cname_level": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the number of CNAMEs to use in the chain.",
						Type:             schema.TypeInt,
					},
					"users_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "This specifies the name of the cookie that identifies the user.",
						Type:             schema.TypeString,
					},
					"users_query_string": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "This names the query parameter that identifies the user.",
						Type:             schema.TypeString,
					},
					"users_regex": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9\\:\\[\\]\\{\\}\\(\\)\\.\\?_\\-\\*\\+\\^\\$\\\\\\/\\|&=!]{1,250})$"),
						Optional:         true,
						Description:      "Specifies a Perl-compatible regular expression with which to substitute the request's user ID.",
						Type:             schema.TypeString,
					},
					"users_replace": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z0-9_\\-]|\\$[1-9]){1,250})$"),
						Optional:         true,
						Description:      "Specifies a string to replace the request's user ID matched by `usersRegex`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"sales_force_commerce_cloud_client": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "If you use the Salesforce Commerce Cloud platform for your origin content, this behavior allows your edge content managed by Akamai to contact directly to origin. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Akamai Connector for Salesforce Commerce Cloud.",
						Type:        schema.TypeBool,
					},
					"connector_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\.]+\\-[a-zA-Z0-9_\\.]+\\-[a-zA-Z0-9\\-_\\.]+$|^door2.dw.com$"),
						Optional:         true,
						Description:      "An ID value that helps distinguish different types of traffic sent from Akamai to the Salesforce Commerce Cloud. Form the value as `instance-realm-customer`, where `instance` is either `production` or `development`, `realm` is your Salesforce Commerce Cloud service `$REALM` value, and `customer` is the name for your organization in Salesforce Commerce Cloud.  You can use alphanumeric characters, underscores, or dot characters within dash-delimited segment values.",
						Type:             schema.TypeString,
					},
					"origin_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DEFAULT", "CUSTOMER"}, false)),
						Optional:         true,
						Description:      "Specifies where the origin is.",
						Type:             schema.TypeString,
					},
					"sf3c_origin_host": {
						Optional:    true,
						Description: "This specifies the hostname or IP address of the custom Salesforce origin.",
						Type:        schema.TypeString,
					},
					"origin_host_header": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DEFAULT", "CUSTOMER"}, false)),
						Optional:         true,
						Description:      "Specifies where the `Host` header is defined.",
						Type:             schema.TypeString,
					},
					"sf3c_origin_host_header": {
						Optional:    true,
						Description: "This specifies the hostname or IP address of the custom Salesforce host header.",
						Type:        schema.TypeString,
					},
					"allow_override_origin_cache_key": {
						Optional:    true,
						Description: "When enabled, overrides the forwarding origin's cache key.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"sales_force_commerce_cloud_provider": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This manages traffic between mutual customers and the Salesforce Commerce Cloud platform. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables Akamai Provider for Salesforce Commerce Cloud.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"sales_force_commerce_cloud_provider_host_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Manages host header values sent to the Salesforce Commerce Cloud platform. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"host_header_source": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PROPERTY", "CUSTOMER"}, false)),
						Optional:         true,
						Description:      "Specify where the host header derives from.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"save_post_dca_processing": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Used in conjunction with the `cachePost` behavior, this behavior allows the body of POST requests to be processed through Dynamic Content Assembly.  Contact Akamai Professional Services for help configuring it. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables processing of POST requests.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"schedule_invalidation": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies when cached content that satisfies a rule's criteria expires, optionally at repeating intervals. In addition to periodic cache flushes, you can use this behavior to minimize potential conflicts when related objects expire at different times. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"start": {
						Optional:    true,
						Description: "The UTC date and time when matching cached content is to expire.",
						Type:        schema.TypeString,
					},
					"repeat": {
						Optional:    true,
						Description: "When enabled, invalidation recurs periodically from the `start` time based on the `repeatInterval` time.",
						Type:        schema.TypeBool,
					},
					"repeat_interval": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies how often to invalidate content from the `start` time, expressed in seconds. For example, an expiration set to midnight and an interval of `86400` seconds invalidates content once a day.  Repeating intervals of less than 5 minutes are not allowed for `NetStorage` origins.",
						Type:             schema.TypeString,
					},
					"refresh_method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INVALIDATE", "PURGE"}, false)),
						Optional:         true,
						Description:      "Specifies how to invalidate the content.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"script_management": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Ensures unresponsive linked JavaScript files do not prevent HTML pages from loading. See `Script Management API` for more information. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Script Management feature.",
						Type:        schema.TypeBool,
					},
					"serviceworker": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"YES_SERVICE_WORKER", "NO_SERVICE_WORKER"}, false)),
						Optional:         true,
						Description:      "Script Management uses a JavaScript service worker called `akam-sw.js`. It applies a policy that helps you manage scripts.",
						Type:             schema.TypeString,
					},
					"timestamp": {
						Optional:    true,
						Description: "A read-only epoch timestamp that represents the last time a Script Management policy was synchronized with its Ion property.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"segmented_content_protection": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Validates authorization tokens at the edge server to prevent unauthorized link sharing. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"token_authentication_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the segmented content protection behavior.",
						Type:        schema.TypeBool,
					},
					"key": {
						ValidateDiagFunc: validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$"),
						Optional:         true,
						Description:      "Specifies the encryption key to use as a shared secret to validate tokens.",
						Type:             schema.TypeString,
					},
					"use_advanced": {
						Optional:    true,
						Description: "Allows you to specify advanced `transitionKey` and `salt` options.",
						Type:        schema.TypeBool,
					},
					"transition_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$")),
						Optional:         true,
						Description:      "An alternate encryption key to match along with the `key` field, allowing you to rotate keys with no down time.",
						Type:             schema.TypeString,
					},
					"salt": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validation.ToDiagFunc(validation.StringLenBetween(16, 16))),
						Optional:         true,
						Description:      "Specifies a salt as input into the token for added security. This value needs to match the salt used in the token generation code.",
						Type:             schema.TypeString,
					},
					"header_for_salt": {
						Optional:    true,
						Description: "This allows you to include additional salt properties specific to each end user to strengthen the relationship between the session token and playback session. This specifies the set of request headers whose values generate the salt value, typically `User-Agent`, `X-Playback-Session-Id`, and `Origin`. Any specified header needs to appear in the player's request.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"field_carry_over": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"session_id": {
						Optional:    true,
						Description: "Enabling this option carries the `session_id` value from the access token over to the session token, for use in tracking and counting unique playback sessions.",
						Type:        schema.TypeBool,
					},
					"data_payload": {
						Optional:    true,
						Description: "Enabling this option carries the `data/payload` field from the access token over to the session token, allowing access to opaque data for log analysis for a URL protected by a session token.",
						Type:        schema.TypeBool,
					},
					"ip": {
						Optional:    true,
						Description: "Enabling this restricts content access to a specific IP address, only appropriate if it does not change during the playback session.",
						Type:        schema.TypeBool,
					},
					"acl": {
						Optional:    true,
						Description: "Enabling this option carries the `ACL` field from the access token over to the session token, to limit the requesting client's access to the specific URL or path set in the `ACL` field. Playback may fail if the base path of the master playlist (and variant playlist, plus segments) varies from that of the `ACL` field.",
						Type:        schema.TypeBool,
					},
					"token_auth_hls_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_token_in_uri": {
						Optional:    true,
						Description: "When enabled, passes tokens in HLS variant manifest URLs and HLS segment URLs, as an alternative to cookies.",
						Type:        schema.TypeBool,
					},
					"hls_master_manifest_files": {
						Optional:    true,
						Description: "Specifies the set of filenames that form HLS master manifest URLs. You can use `*` wildcard character that matches zero or more characters. Make sure to specify master manifest filenames uniquely, to distinguish them from variant manifest files.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"token_auth_dash_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enable_token_in_query_string": {
						Optional:    true,
						Description: "When enabled, in the DASH media encryption variant, passes tokens in query strings, as an alternative to cookies.",
						Type:        schema.TypeBool,
					},
					"token_revocation_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"token_revocation_enabled": {
						Optional:    true,
						Description: "Enable this to deny requests from playback URLs that contain a `TokenAuth` token that uses specific token identifiers.",
						Type:        schema.TypeBool,
					},
					"revoked_list_id": {
						Optional:    true,
						Description: "Identifies the `TokenAuth` tokens to block from accessing your content.",
						Type:        schema.TypeInt,
					},
					"media_encryption_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"hls_media_encryption": {
						Optional:    true,
						Description: "Enables HLS Segment Encryption.",
						Type:        schema.TypeBool,
					},
					"dash_media_encryption": {
						Optional:    true,
						Description: "Whether to enable DASH Media Encryption.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"segmented_media_optimization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Optimizes segmented media for live or streaming delivery contexts. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"behavior": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ON_DEMAND", "LIVE"}, false)),
						Optional:         true,
						Description:      "Sets the type of media content to optimize.",
						Type:             schema.TypeString,
					},
					"enable_ull_streaming": {
						Optional:    true,
						Description: "Enables ultra low latency (ULL) streaming. ULL reduces latency and decreases overall transfer time of live streams.",
						Type:        schema.TypeBool,
					},
					"show_advanced": {
						Optional:    true,
						Description: "Allows you to configure advanced media options.",
						Type:        schema.TypeBool,
					},
					"live_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CONTINUOUS", "EVENT"}, false)),
						Optional:         true,
						Description:      "The type of live media.",
						Type:             schema.TypeString,
					},
					"start_time": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "This specifies when the live media event begins.",
						Type:             schema.TypeString,
					},
					"end_time": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "This specifies when the live media event ends.",
						Type:             schema.TypeString,
					},
					"dvr_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CONFIGURABLE", "UNKNOWN"}, false)),
						Optional:         true,
						Description:      "The type of DVR.",
						Type:             schema.TypeString,
					},
					"dvr_window": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Set the duration for your media, or `0m` if a DVR is not required.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"segmented_media_streaming_prefetch": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Prefetches HLS and DASH media stream manifest and segment files, accelerating delivery to end users. For prefetching to work, your origin media's response needs to specify `CDN-Origin-Assist-Prefetch-Path` headers with each URL to prefetch, expressed as either a relative or absolute path. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables media stream prefetching.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"set_variable": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Modify a variable to insert into subsequent fields within the rule tree.  Use this behavior to specify the predeclared `variableName` and determine from where to derive its new value. Based on this `valueSource`, you can either generate the value, extract it from some part of the incoming request, assign it from another variable (including a set of built-in system variables), or directly specify its text.  Optionally choose a `transform` function to modify the value once. See `Support for variables` for more information. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"variable_name": {
						Optional:    true,
						Description: "Specifies the predeclared root name of the variable to modify.  When you declare a variable name such as `VAR`, its name is preprended with `PMUSER_` and accessible in a `user` namespace, so that you invoke it in subsequent text fields within the rule tree as `{{user.PMUSER_VAR}}`. In deployed `XML metadata`, it appears as `%(PMUSER_VAR)`.",
						Type:        schema.TypeString,
					},
					"value_source": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"EXPRESSION", "EXTRACT", "GENERATE"}, false)),
						Optional:         true,
						Description:      "Determines how you want to set the value.",
						Type:             schema.TypeString,
					},
					"variable_value": {
						Optional:    true,
						Description: "This directly specifies the value to assign to the variable. The expression may include a mix of static text and other variables, such as `new_filename.{{builtin.AK_EXTENSION}}` to embed a system variable.",
						Type:        schema.TypeString,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_CERTIFICATE", "CLIENT_REQUEST_HEADER", "COOKIE", "EDGESCAPE", "PATH_COMPONENT_OFFSET", "QUERY_STRING", "DEVICE_PROFILE", "RESPONSE_HEADER", "SET_COOKIE"}, false)),
						Optional:         true,
						Description:      "This specifies from where to get the value.",
						Type:             schema.TypeString,
					},
					"certificate_field_name": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VERSION", "SERIAL", "FINGERPRINT_MD5", "FINGERPRINT_SHA1", "FINGERPRINT_DYN", "ISSUER_DN", "SUBJECT_DN", "NOT_BEFORE", "NOT_AFTER", "SIGNATURE_ALGORITHM", "SIGNATURE", "CONTENTS_DER", "CONTENTS_PEM", "CONTENTS_PEM_NO_LABELS", "COUNT", "STATUS_MSG", "KEY_LENGTH"}, false)),
						Optional:         true,
						Description:      "Specifies the certificate's content.",
						Type:             schema.TypeString,
					},
					"header_name": {
						Optional:    true,
						Description: "Specifies the case-insensitive name of the HTTP header to extract.",
						Type:        schema.TypeString,
					},
					"response_header_name": {
						Optional:    true,
						Description: "Specifies the case-insensitive name of the HTTP header to extract.",
						Type:        schema.TypeString,
					},
					"set_cookie_name": {
						Optional:    true,
						Description: "Specifies the name of the origin's `Set-Cookie` response header.",
						Type:        schema.TypeString,
					},
					"cookie_name": {
						Optional:    true,
						Description: "Specifies the name of the cookie to extract.",
						Type:        schema.TypeString,
					},
					"location_id": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GEOREGION", "COUNTRY_CODE", "REGION_CODE", "CITY", "DMA", "PMSA", "MSA", "AREACODE", "COUNTY", "FIPS", "LAT", "LONG", "TIMEZONE", "ZIP", "CONTINENT", "NETWORK", "NETWORK_TYPE", "ASNUM", "THROUGHPUT", "BW"}, false)),
						Optional:         true,
						Description:      "Specifies the `X-Akamai-Edgescape` header's field name. Possible values specify basic geolocation, various geographic standards, and information about the client's network. For details on EdgeScape header fields, see the `EdgeScape User Guide`.",
						Type:             schema.TypeString,
					},
					"path_component_offset": {
						Optional:    true,
						Description: "This specifies a portion of the path.  The indexing starts from `1`, so a value of `/path/to/nested/filename.html` and an offset of `1` yields `path`, and `3` yields `nested`. Negative indexes offset from the right, so `-2` also yields `nested`.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "Specifies the name of the query parameter from which to extract the value.",
						Type:        schema.TypeString,
					},
					"generator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HEXRAND", "RAND"}, false)),
						Optional:         true,
						Description:      "This specifies the type of value to generate.",
						Type:             schema.TypeString,
					},
					"number_of_bytes": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 16)),
						Optional:         true,
						Description:      "Specifies the number of random hex bytes to generate.",
						Type:             schema.TypeInt,
					},
					"min_random_number": {
						Optional:    true,
						Description: "Specifies the lower bound of the random number.",
						Type:        schema.TypeInt,
					},
					"max_random_number": {
						Optional:    true,
						Description: "Specifies the upper bound of the random number.",
						Type:        schema.TypeInt,
					},
					"transform": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NONE", "ADD", "BASE_64_DECODE", "BASE_64_ENCODE", "BASE_32_DECODE", "BASE_32_ENCODE", "BITWISE_AND", "BITWISE_NOT", "BITWISE_OR", "BITWISE_XOR", "DECIMAL_TO_HEX", "DECRYPT", "DIVIDE", "ENCRYPT", "EPOCH_TO_STRING", "EXTRACT_PARAM", "HASH", "JSON_EXTRACT", "HEX_TO_DECIMAL", "HEX_DECODE", "HEX_ENCODE", "HMAC", "LOWER", "MD5", "MINUS", "MODULO", "MULTIPLY", "NORMALIZE_PATH_WIN", "REMOVE_WHITESPACE", "COMPRESS_WHITESPACE", "SHA_1", "SHA_256", "STRING_INDEX", "STRING_LENGTH", "STRING_TO_EPOCH", "SUBSTITUTE", "SUBSTRING", "SUBTRACT", "TRIM", "UPPER", "BASE_64_URL_DECODE", "BASE_64_URL_ENCODE", "URL_DECODE", "URL_ENCODE", "URL_DECODE_UNI", "UTC_SECONDS", "XML_DECODE", "XML_ENCODE"}, false)),
						Optional:         true,
						Description:      "Specifies a function to transform the value. For more details on each transform function, see `Set Variable: Operations`.",
						Type:             schema.TypeString,
					},
					"operand_one": {
						Optional:    true,
						Description: "Specifies an additional operand when the `transform` function is set to various arithmetic functions (`ADD`, `SUBTRACT`, `MULTIPLY`, `DIVIDE`, or `MODULO`) or bitwise functions (`BITWISE_AND`, `BITWISE_OR`, or `BITWISE_XOR`).",
						Type:        schema.TypeString,
					},
					"algorithm": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALG_3DES", "ALG_AES128", "ALG_AES256"}, false)),
						Optional:         true,
						Description:      "Specifies the algorithm to apply.",
						Type:             schema.TypeString,
					},
					"encryption_key": {
						ValidateDiagFunc: validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$"),
						Optional:         true,
						Description:      "Specifies the encryption hex key. For `ALG_3DES` it needs to be 48 characters long, 32 characters for `ALG_AES128`, and 64 characters for `ALG_AES256`.",
						Type:             schema.TypeString,
					},
					"initialization_vector": {
						ValidateDiagFunc: validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$"),
						Optional:         true,
						Description:      "Specifies a one-time number as an initialization vector.  It needs to be 15 characters long for `ALG_3DES`, and 32 characters for both `ALG_AES128` and `ALG_AES256`.",
						Type:             schema.TypeString,
					},
					"encryption_mode": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CBC", "ECB"}, false)),
						Optional:         true,
						Description:      "Specifies the encryption mode.",
						Type:             schema.TypeString,
					},
					"nonce": {
						Optional:    true,
						Description: "Specifies the one-time number used for encryption.",
						Type:        schema.TypeString,
					},
					"prepend_bytes": {
						Optional:    true,
						Description: "Specifies a number of random bytes to prepend to the key.",
						Type:        schema.TypeBool,
					},
					"format_string": {
						Optional:    true,
						Description: "Specifies an optional format string for the conversion, using format codes such as `%m/%d/%y` as specified by `strftime`. A blank value defaults to RFC-2616 format.",
						Type:        schema.TypeString,
					},
					"param_name": {
						Optional:    true,
						Description: "Extracts the value for the specified parameter name from a string that contains key/value pairs. (Use `separator` below to parse them.)",
						Type:        schema.TypeString,
					},
					"separator": {
						Optional:    true,
						Description: "Specifies the character that separates pairs of values within the string.",
						Type:        schema.TypeString,
					},
					"min": {
						Optional:    true,
						Description: "Specifies a minimum value for the generated integer.",
						Type:        schema.TypeInt,
					},
					"max": {
						Optional:    true,
						Description: "Specifies a maximum value for the generated integer.",
						Type:        schema.TypeInt,
					},
					"hmac_key": {
						Optional:    true,
						Description: "Specifies the secret to use in generating the base64-encoded digest.",
						Type:        schema.TypeString,
					},
					"hmac_algorithm": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SHA1", "SHA256", "MD5"}, false)),
						Optional:         true,
						Description:      "Specifies the algorithm to use to generate the base64-encoded digest.",
						Type:             schema.TypeString,
					},
					"ip_version": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IPV4", "IPV6"}, false)),
						Optional:         true,
						Description:      "Specifies the IP version under which a subnet mask generates.",
						Type:             schema.TypeString,
					},
					"ipv6_prefix": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 128)),
						Optional:         true,
						Description:      "Specifies the prefix of the IPV6 address, a value between 0 and 128.",
						Type:             schema.TypeInt,
					},
					"ipv4_prefix": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 32)),
						Optional:         true,
						Description:      "Specifies the prefix of the IPV4 address, a value between 0 and 32.",
						Type:             schema.TypeInt,
					},
					"sub_string": {
						Optional:    true,
						Description: "Specifies a substring for which the returned value represents a zero-based offset of where it appears in the original string, or `-1` if there's no match.",
						Type:        schema.TypeString,
					},
					"regex": {
						Optional:    true,
						Description: "Specifies the regular expression pattern (PCRE) to match the value.",
						Type:        schema.TypeString,
					},
					"replacement": {
						Optional:    true,
						Description: "Specifies the replacement string. Reinsert grouped items from the match into the replacement using `$1`, `$2` ... `$n`.",
						Type:        schema.TypeString,
					},
					"case_sensitive": {
						Optional:    true,
						Description: "Enabling this makes all matches case sensitive.",
						Type:        schema.TypeBool,
					},
					"global_substitution": {
						Optional:    true,
						Description: "Replaces all matches in the string, not just the first.",
						Type:        schema.TypeBool,
					},
					"start_index": {
						Optional:    true,
						Description: "Specifies the zero-based character offset at the start of the substring. Negative indexes specify the offset from the end of the string.",
						Type:        schema.TypeInt,
					},
					"end_index": {
						Optional:    true,
						Description: "Specifies the zero-based character offset at the end of the substring, without including the character at that index position. Negative indexes specify the offset from the end of the string.",
						Type:        schema.TypeInt,
					},
					"except_chars": {
						Optional:    true,
						Description: "Specifies characters `not` to encode, possibly overriding the default set.",
						Type:        schema.TypeString,
					},
					"force_chars": {
						Optional:    true,
						Description: "Specifies characters to encode, possibly overriding the default set.",
						Type:        schema.TypeString,
					},
					"device_profile": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_MOBILE", "IS_TABLET", "IS_WIRELESS_DEVICE", "PHYSICAL_SCREEN_HEIGHT", "PHYSICAL_SCREEN_WIDTH", "RESOLUTION_HEIGHT", "RESOLUTION_WIDTH", "VIEWPORT_WIDTH", "BRAND_NAME", "DEVICE_OS", "DEVICE_OS_VERSION", "DUAL_ORIENTATION", "MAX_IMAGE_HEIGHT", "MAX_IMAGE_WIDTH", "MOBILE_BROWSER", "MOBILE_BROWSER_VERSION", "PDF_SUPPORT", "COOKIE_SUPPORT"}, false)),
						Optional:         true,
						Description:      "Specifies the client device attribute. Possible values specify information about the client device, including device type, size and browser. For details on fields, see `Device Characterization`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"simulate_error_code": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior simulates various error response codes. Contact Akamai Professional Services for help configuring it. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"error_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ERR_DNS_TIMEOUT", "ERR_SUREROUTE_DNS_FAIL", "ERR_DNS_FAIL", "ERR_CONNECT_TIMEOUT", "ERR_NO_GOOD_FWD_IP", "ERR_DNS_IN_REGION", "ERR_CONNECT_FAIL", "ERR_READ_TIMEOUT", "ERR_READ_ERROR", "ERR_WRITE_ERROR"}, false)),
						Optional:         true,
						Description:      "Specifies the type of error.",
						Type:             schema.TypeString,
					},
					"timeout": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "When the `errorType` is `ERR_CONNECT_TIMEOUT`, `ERR_DNS_TIMEOUT`, `ERR_SUREROUTE_DNS_FAIL`, or `ERR_READ_TIMEOUT`, generates an error after the specified amount of time from the initial request.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"site_shield": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior implements the `Site Shield` feature, which helps prevent non-Akamai machines from contacting your origin. You get an email with a list of Akamai servers allowed to contact your origin, with which you establish an Access Control List on your firewall to prevent any other requests. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"ssmap": {
						Optional:    true,
						Description: "Identifies the hostname for the Site Shield map. See `Create a Site Shield map` for more details. Form an object with a `value` key that references the hostname, for example: `\"ssmap\":{\"value\":\"ss.akamai.net\"}`.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"value": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"srmap": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"china_cdn_map": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"has_mixed_hosts": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeBool,
								},
								"src": {
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"FALLBACK", "PROTECTED_HOST_MATCH", "ORIGIN_MATCH", "PREVIOUS_MAP", "PROPERTY_MATCH"}, false)),
									Optional:         true,
									Description:      "",
									Type:             schema.TypeString,
								},
							},
						},
					},
					"nossmap": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"standard_tls_migration": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows migration to Standard TLS.",
						Type:        schema.TypeBool,
					},
					"migration_from": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SHARED_CERT", "NON_SECURE", "ENHANCED_SECURE"}, false)),
						Optional:         true,
						Description:      "What kind of traffic you're migrating from.",
						Type:             schema.TypeString,
					},
					"allow_https_upgrade": {
						Optional:    true,
						Description: "Allows temporary upgrade of HTTP traffic to HTTPS.",
						Type:        schema.TypeBool,
					},
					"allow_https_downgrade": {
						Optional:    true,
						Description: "Allow temporary downgrade of HTTPS traffic to HTTP. This removes various `Origin`, `Referer`, `Cookie`, `Cookie2`, `sec-*` and `proxy-*` headers from the request to origin.",
						Type:        schema.TypeBool,
					},
					"migration_start_time": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies when to start migrating the cache.",
						Type:             schema.TypeString,
					},
					"migration_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]$|^[1-2]\\d$|^30$"),
						Optional:         true,
						Description:      "Specifies the number of days to migrate the cache.",
						Type:             schema.TypeInt,
					},
					"cache_sharing_start_time": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies when to start cache sharing.",
						Type:             schema.TypeString,
					},
					"cache_sharing_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]$|^[1-2]\\d$|^30$"),
						Optional:         true,
						Description:      "Specifies the number cache sharing days.",
						Type:             schema.TypeInt,
					},
					"is_certificate_sni_only": {
						Optional:    true,
						Description: "Sets whether your new certificate is SNI-only.",
						Type:        schema.TypeBool,
					},
					"is_tiered_distribution_used": {
						Optional:    true,
						Description: "Allows you to align traffic to various `tieredDistribution` areas.",
						Type:        schema.TypeBool,
					},
					"td_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GLOBAL", "APAC", "EUROPE", "US_EAST", "US_CENTRAL", "US_WEST", "AUSTRALIA", "GLOBAL_LEGACY"}, false)),
						Optional:         true,
						Description:      "Specifies the `tieredDistribution` location.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"standard_tls_migration_override": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior is for internal usage only. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"info": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"strict_header_parsing": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior specifies how the edge servers should handle requests containing improperly formatted or invalid headers that don’t comply with `RFC 9110`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"valid_mode": {
						Optional:    true,
						Description: "Rejects requests made with non-RFC-compliant headers that contain invalid characters in the header name or value or which contain invalidly-folded header lines. When disabled, the edge servers allow such requests, passing the invalid headers to the origin server unchanged.",
						Type:        schema.TypeBool,
					},
					"strict_mode": {
						Optional:    true,
						Description: "Rejects requests made with non-RFC-compliant, improperly formatted headers, where the header line starts with a colon, misses a colon or doesn’t end with CR LF. When disabled, the edge servers allow such requests, but correct the violation by removing or rewriting the header line before passing the headers to the origin server.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"sub_customer": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "When positioned in a property's top-level default rule, enables various `Cloud Embed` features that allow you to leverage Akamai's CDN architecture for your own subcustomers.  This behavior's options allow you to use Cloud Embed to configure your subcustomers' content.  Once enabled, you can use the `Akamai Cloud Embed API` (ACE) to assign subcustomers to this base configuration, and to customize policies for them.  See also the `dynamicWebContent` behavior to configure subcustomers' dynamic web content. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows Cloud Embed to dynamically modify your subcustomers' content.",
						Type:        schema.TypeBool,
					},
					"origin": {
						Optional:    true,
						Description: "Allows you to assign origin hostnames for customers.",
						Type:        schema.TypeBool,
					},
					"partner_domain_suffix": {
						Optional:    true,
						Description: "This specifies the appropriate domain suffix, which you should typically match with your property hostname. It identifies the domain as trustworthy on the Akamai network, despite being defined within Cloud Embed, outside of your base property configuration. Include this domain suffix if you want to purge subcustomer URLs. For example, if you provide a value of `suffix.example.com`, then to purge `subcustomer.com/some/path`, specify `subcustomer.com.suffix.example.com/some/path` as the purge request's URL.",
						Type:        schema.TypeString,
					},
					"caching": {
						Optional:    true,
						Description: "Modifies content caching rules.",
						Type:        schema.TypeBool,
					},
					"referrer": {
						Optional:    true,
						Description: "Sets subcustomers' referrer allowlists or blocklists.",
						Type:        schema.TypeBool,
					},
					"ip": {
						Optional:    true,
						Description: "Sets subcustomers' IP allowlists or blocklists.",
						Type:        schema.TypeBool,
					},
					"geo_location": {
						Optional:    true,
						Description: "Sets subcustomers' location-based allowlists or blocklists.",
						Type:        schema.TypeBool,
					},
					"refresh_content": {
						Optional:    true,
						Description: "Allows you to reschedule when content validates for subcustomers.",
						Type:        schema.TypeBool,
					},
					"modify_path": {
						Optional:    true,
						Description: "Modifies a subcustomer's request path.",
						Type:        schema.TypeBool,
					},
					"cache_key": {
						Optional:    true,
						Description: "Allows you to set which query parameters are included in the cache key.",
						Type:        schema.TypeBool,
					},
					"token_authorization": {
						Optional:    true,
						Description: "When enabled, this allows you to configure edge servers to use tokens to control access to subcustomer content.  Use Cloud Embed to configure the token to appear in a cookie, header, or query parameter.",
						Type:        schema.TypeBool,
					},
					"site_failover": {
						Optional:    true,
						Description: "Allows you to configure unique failover sites for each subcustomer's policy.",
						Type:        schema.TypeBool,
					},
					"content_compressor": {
						Optional:    true,
						Description: "Allows compression of subcustomer content.",
						Type:        schema.TypeBool,
					},
					"access_control": {
						Optional:    true,
						Description: "When enabled, this allows you to deny requests to a subcustomer's content based on specific match conditions, which you use Cloud Embed to configure in each subcustomer's policy.",
						Type:        schema.TypeBool,
					},
					"dynamic_web_content": {
						Optional:    true,
						Description: "Allows you to apply the `dynamicWebContent` behavior to further modify how dynamic content behaves for subcustomers.",
						Type:        schema.TypeBool,
					},
					"on_demand_video_delivery": {
						Optional:    true,
						Description: "Enables delivery of media assets to subcustomers.",
						Type:        schema.TypeBool,
					},
					"large_file_delivery": {
						Optional:    true,
						Description: "Enables large file delivery for subcustomers.",
						Type:        schema.TypeBool,
					},
					"live_video_delivery": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"web_application_firewall": {
						Optional:    true,
						Description: "Web application firewall (WAF) filters, monitors, and blocks certain HTTP traffic. Use `Akamai Cloud Embed` to add a specific behavior to a subcustomer policy and configure how WAF protection is applied.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"sure_route": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `SureRoute` feature continually tests different routes between origin and edge servers to identify the optimal path. By default, it conducts `races` to identify alternative paths to use in case of a transmission failure. These races increase origin traffic slightly. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the SureRoute behavior, to optimize delivery of non-cached content.",
						Type:        schema.TypeBool,
					},
					"type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PERFORMANCE", "CUSTOM_MAP"}, false)),
						Optional:         true,
						Description:      "Specifies the set of edge servers used to test routes.",
						Type:             schema.TypeString,
					},
					"custom_map": {
						Optional:    true,
						Description: "If `type` is `CUSTOM_MAP`, this specifies the map string provided to you by Akamai Professional Services, or included as part of the `Site Shield` product.",
						Type:        schema.TypeString,
					},
					"test_object_url": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the path and filename for your origin's test object to use in races to test routes.",
						Type:             schema.TypeString,
					},
					"sr_download_link_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"to_host_status": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"INCOMING_HH", "OTHER"}, false)),
						Optional:         true,
						Description:      "Specifies which hostname to use.",
						Type:             schema.TypeString,
					},
					"to_host": {
						Optional:    true,
						Description: "If `toHostStatus` is `OTHER`, this specifies the custom `Host` header to use when requesting the SureRoute test object.",
						Type:        schema.TypeString,
					},
					"race_stat_ttl": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the time-to-live to preserve SureRoute race results, typically `30m`. If traffic exceeds a certain threshold after TTL expires, the overflow is routed directly to the origin, not necessarily optimally. If traffic remains under the threshold, the route is determined by the winner of the most recent race.",
						Type:             schema.TypeString,
					},
					"force_ssl_forward": {
						Optional:    true,
						Description: "Forces SureRoute to use SSL when requesting the origin's test object, appropriate if your origin does not respond to HTTP requests, or responds with a redirect to HTTPS.",
						Type:        schema.TypeBool,
					},
					"allow_fcm_parent_override": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
					"enable_custom_key": {
						Optional:    true,
						Description: "When disabled, caches race results under the race destination's hostname. If enabled, use `customStatKey` to specify a custom hostname.",
						Type:        schema.TypeBool,
					},
					"custom_stat_key": {
						ValidateDiagFunc: validateRegexOrVariable("^([a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})(\\.[a-zA-Z0-9][a-zA-Z0-9\\-]{0,62})+$"),
						Optional:         true,
						Description:      "This specifies a hostname under which to cache race results. This may be useful when a property corresponds to many origin hostnames. By default, SureRoute would launch races for each origin, but consolidating under a single hostname runs only one race.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"tcp_optimization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior is deprecated, but you should not disable or remove it if present. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"display": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"tea_leaf": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Note that this behavior is decommissioned. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, capture HTTPS requests and responses, and send the data to your IBM Tealeaf account.",
						Type:        schema.TypeBool,
					},
					"limit_to_dynamic": {
						Optional:    true,
						Description: "Limit traffic to dynamic, uncached (`No-Store`) content.",
						Type:        schema.TypeBool,
					},
					"ibm_customer_id": {
						Optional:    true,
						Description: "The integer identifier for the IBM Tealeaf Connector account.",
						Type:        schema.TypeInt,
					},
				},
			},
		},
		"tiered_distribution": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows Akamai edge servers to retrieve cached content from other Akamai servers, rather than directly from the origin. These interim `parent` servers in the `cache hierarchy` (`CH`) are positioned close to the origin, and fall along the path from the origin to the edge server. Tiered Distribution typically reduces the origin server's load, and reduces the time it takes for edge servers to refresh content. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, activates tiered distribution.",
						Type:        schema.TypeBool,
					},
					"tiered_distribution_map": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CH2", "CHAPAC", "CHEU2", "CHEUS2", "CHCUS2", "CHWUS2", "CHAUS", "CH"}, false)),
						Optional:         true,
						Description:      "Optionally map the tiered parent server's location close to your origin. A narrower local map minimizes the origin server's load, and increases the likelihood the requested object is cached. A wider global map reduces end-user latency, but decreases the likelihood the requested object is in any given parent server's cache.  This option cannot apply if the property is marked as secure. See `Secure property requirements` for guidance.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"tiered_distribution_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows Akamai edge servers to retrieve cached content from other Akamai servers, rather than directly from the origin. These interim `parent` servers in the `cache hierarchy` (`CH`) are positioned close to the origin, and fall along the path from the origin to the edge server. Tiered Distribution typically reduces the origin server's load, and reduces the time it takes for edge servers to refresh content.  This advanced behavior provides a wider set of options than `tieredDistribution`. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "When enabled, activates tiered distribution.",
						Type:        schema.TypeBool,
					},
					"method": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SERIAL_PREPEND", "DOMAIN_LOOKUP"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"policy": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"PERFORMANCE", "TIERED_DISTRIBUTION", "FAIL_OVER", "SITE_SHIELD", "SITE_SHIELD_PERFORMANCE"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"tiered_distribution_map": {
						Optional:    true,
						Description: "Optionally map the tiered parent server's location close to your origin: `CHEU2` for Europe; `CHAUS` for Australia; `CHAPAC` for China and the Asian Pacific area; `CHWUS2`, `CHCUS2`, and `CHEUS2` for different parts of the United States. Choose `CH` or `CH2` for a more global map. A narrower local map minimizes the origin server's load, and increases the likelihood the requested object is cached. A wider global map reduces end-user latency, but decreases the likelihood the requested object is in any given parent server's cache.  This option cannot apply if the property is marked as secure. See `Secure property requirements` for guidance.",
						Type:        schema.TypeString,
					},
					"allowall": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"tiered_distribution_customization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "With Tiered Distribution, Akamai edge servers retrieve cached content from other Akamai servers, rather than directly from the origin. This behavior sets custom Tiered Distribution maps (TD0) and migrates TD1 maps configured with `advanced features` to Cloud Wrapper. You need to enable `cloudWrapper` within the same rule. This behavior is for internal usage only. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"tier1_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"custom_map_enabled": {
						Optional:    true,
						Description: "Enables custom maps.",
						Type:        schema.TypeBool,
					},
					"custom_map_name": {
						ValidateDiagFunc: validateRegexOrVariable("^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)+(akamai|akamaiedge)\\.net$"),
						Optional:         true,
						Description:      "Specifies the custom map name.",
						Type:             schema.TypeString,
					},
					"serial_start": {
						Optional:    true,
						Description: "Specifies a numeric serial start value.",
						Type:        schema.TypeString,
					},
					"serial_end": {
						Optional:    true,
						Description: "Specifies a numeric serial end value. Akamai uses serial numbers to group machines and share objects in their cache with other machines in the same region.",
						Type:        schema.TypeString,
					},
					"hash_algorithm": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GCC", "JENKINS"}, false)),
						Optional:         true,
						Description:      "Specifies the hash algorithm.",
						Type:             schema.TypeString,
					},
					"cloudwrapper_map_migration_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"map_migration_enabled": {
						Optional:    true,
						Description: "Enables migration of the custom map to Cloud Wrapper.",
						Type:        schema.TypeBool,
					},
					"migration_within_cw_maps_enabled": {
						Optional:    true,
						Description: "Enables migration within Cloud Wrapper maps.",
						Type:        schema.TypeBool,
					},
					"location": {
						Optional:    true,
						Description: "Location from which Cloud Wrapper migration is performed. User should choose the existing Cloud Wrapper location. The new Cloud Wrapper location (to which migration has to happen) is expected to be updated as part of the main \"Cloud Wrapper\" behavior.",
						Type:        schema.TypeString,
					},
					"migration_start_date": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies when to start migrating the map.",
						Type:             schema.TypeString,
					},
					"migration_end_date": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies when the map migration should end.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"timeout": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Sets the HTTP connect timeout. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the timeout, for example `10s`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"uid_configuration": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows you to extract unique identifier (UID) values from live traffic, for use in OTA applications. Note that you are responsible for maintaining the security of any data that may identify individual users. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"legal_text": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Allows you to extract UIDs from client requests.",
						Type:        schema.TypeBool,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING", "VARIABLE"}, false)),
						Optional:         true,
						Description:      "Where to extract the UID value from.",
						Type:             schema.TypeString,
					},
					"header_name": {
						Optional:    true,
						Description: "This specifies the name of the HTTP header from which to extract the UID value.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "This specifies the name of the query parameter from which to extract the UID value.",
						Type:        schema.TypeString,
					},
					"variable_name": {
						Optional:    true,
						Description: "This specifies the name of the rule tree variable from which to extract the UID value.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"validate_entity_tag": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Instructs edge servers to compare the request's `ETag` header with that of the cached object. If they differ, the edge server sends a new copy of the object. This validation occurs in addition to the default validation of `Last-Modified` and `If-Modified-Since` headers. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the ETag validation behavior.",
						Type:        schema.TypeBool,
					},
					"non_strict_enabled": {
						Optional:    true,
						Description: "Whether you want to allow strong `ETag` values that are not surrounded by double quotes. Technically these are malformed and non-standard, but are commonly used.",
						Type:        schema.TypeBool,
					},
					"weak_enabled": {
						Optional:    true,
						Description: "Whether you want to allow weak `ETag` values that start with `W/`.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"verify_json_web_token": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows you to use JSON Web Tokens (JWT) to verify requests. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING"}, false)),
						Optional:         true,
						Description:      "Specify from where to extract the JWT value.",
						Type:             schema.TypeString,
					},
					"header_name": {
						Optional:    true,
						Description: "This specifies the name of the header from which to extract the JWT value.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "This specifies the name of the query parameter from which to extract the JWT value.",
						Type:        schema.TypeString,
					},
					"jwt": {
						Optional:    true,
						Description: "An identifier for the JWT keys collection.",
						Type:        schema.TypeString,
					},
					"enable_rs256": {
						Optional:    true,
						Description: "Verifies JWTs signed with the RS256 algorithm. This signature helps ensure that the token hasn't been tampered with.",
						Type:        schema.TypeBool,
					},
					"enable_es256": {
						Optional:    true,
						Description: "Verifies JWTs signed with the ES256 algorithm. This signature helps ensure that the token hasn't been tampered with.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"verify_json_web_token_for_dcp": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows you to use JSON web tokens (JWT) to verify requests for use in implementing `IoT Edge Connect`, which you use the `dcp` behavior to configure. You can specify the location in a request to pass a JSON web token (JWT), collections of public keys to verify the integrity of this token, and specific claims to extract from it. Use the `verifyJsonWebToken` behavior for other JWT validation. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"extract_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING", "CLIENT_REQUEST_HEADER_AND_QUERY_STRING"}, false)),
						Optional:         true,
						Description:      "Specifies where to get the JWT value from.",
						Type:             schema.TypeString,
					},
					"primary_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQUEST_HEADER", "QUERY_STRING"}, false)),
						Optional:         true,
						Description:      "Specifies the primary location to extract the JWT value from. If the specified option doesn't include the JWTs, the system checks the secondary one.",
						Type:             schema.TypeString,
					},
					"custom_header": {
						Optional:    true,
						Description: "The JWT value comes from the `X-Akamai-DCP-Token` header by default.  Enabling this option allows you to extract it from another header name that you specify.",
						Type:        schema.TypeBool,
					},
					"header_name": {
						Optional:    true,
						Description: "This specifies the name of the header to extract the JWT value from.",
						Type:        schema.TypeString,
					},
					"query_parameter_name": {
						Optional:    true,
						Description: "Specifies the name of the query parameter from which to extract the JWT value.",
						Type:        schema.TypeString,
					},
					"jwt": {
						Optional:    true,
						Description: "An identifier for the JWT keys collection.",
						Type:        schema.TypeString,
					},
					"extract_client_id": {
						Optional:    true,
						Description: "Allows you to extract the client ID claim name stored in JWT.",
						Type:        schema.TypeBool,
					},
					"client_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,20}$"),
						Optional:         true,
						Description:      "This specifies the claim name.",
						Type:             schema.TypeString,
					},
					"extract_authorizations": {
						Optional:    true,
						Description: "Allows you to extract the authorization groups stored in the JWT.",
						Type:        schema.TypeBool,
					},
					"authorizations": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,20}$"),
						Optional:         true,
						Description:      "This specifies the authorization group name.",
						Type:             schema.TypeString,
					},
					"extract_user_name": {
						Optional:    true,
						Description: "Allows you to extract the user name stored in the JWT.",
						Type:        schema.TypeBool,
					},
					"user_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,20}$"),
						Optional:         true,
						Description:      "This specifies the user name.",
						Type:             schema.TypeString,
					},
					"enable_rs256": {
						Optional:    true,
						Description: "Verifies JWTs signed with the RS256 algorithm. This signature helps to ensure that the token hasn't been tampered with.",
						Type:        schema.TypeBool,
					},
					"enable_es256": {
						Optional:    true,
						Description: "Verifies JWTs signed with the ES256 algorithm. This signature helps to ensure that the token hasn't been tampered with.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"verify_token_authorization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Verifies Auth 2.0 tokens. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"use_advanced": {
						Optional:    true,
						Description: "If enabled, allows you to specify advanced options such as `algorithm`, `escapeHmacInputs`, `ignoreQueryString`, `transitionKey`, and `salt`.",
						Type:        schema.TypeBool,
					},
					"location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COOKIE", "QUERY_STRING", "CLIENT_REQUEST_HEADER"}, false)),
						Optional:         true,
						Description:      "Specifies where to find the token in the incoming request.",
						Type:             schema.TypeString,
					},
					"location_id": {
						Optional:    true,
						Description: "When `location` is `CLIENT_REQUEST_HEADER`, specifies the name of the incoming request's header where to find the token.",
						Type:        schema.TypeString,
					},
					"algorithm": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"SHA256", "SHA1", "MD5"}, false)),
						Optional:         true,
						Description:      "Specifies the algorithm that generates the token. It needs to match the method chosen in the token generation code.",
						Type:             schema.TypeString,
					},
					"escape_hmac_inputs": {
						Optional:    true,
						Description: "URL-escapes HMAC inputs passed in as query parameters.",
						Type:        schema.TypeBool,
					},
					"ignore_query_string": {
						Optional:    true,
						Description: "Enabling this removes the query string from the URL used to form an encryption key.",
						Type:        schema.TypeBool,
					},
					"key": {
						ValidateDiagFunc: validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$"),
						Optional:         true,
						Description:      "The shared secret used to validate tokens, which needs to match the key used in the token generation code.",
						Type:             schema.TypeString,
					},
					"transition_key": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validateRegexOrVariable("^(0x)?[0-9a-fA-F]+$")),
						Optional:         true,
						Description:      "Specifies a transition key as a hex value.",
						Type:             schema.TypeString,
					},
					"salt": {
						ValidateDiagFunc: validateAny(validation.ToDiagFunc(validation.StringIsEmpty), validation.ToDiagFunc(validation.StringLenBetween(16, 16))),
						Optional:         true,
						Description:      "Specifies a salt string for input when generating the token, which needs to match the salt value used in the token generation code.",
						Type:             schema.TypeString,
					},
					"failure_response": {
						Optional:    true,
						Description: "When enabled, sends an HTTP error when an authentication test fails.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"virtual_waiting_room": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior helps you maintain business continuity for dynamic applications in high-demand situations such as flash sales. It decreases abandonment by providing a user-friendly waiting room experience. FIFO (First-in First-out) is a request processing mechanism that prioritizes the first requests that enter the waiting room to send them first to the origin. Users can see both their estimated arrival time and position in the line. With Cloudlets available on your contract, choose `Your services` > `Edge logic Cloudlets` to control Virtual Waitig Room within `Control Center`. Otherwise use the `Cloudlets API` to configure it programmatically. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "This identifies the Visitor Waiting Room Cloudlet shared policy to use with this behavior. You can list available shared policies with the `Cloudlets API`.",
						Type:        schema.TypeInt,
					},
					"domain_config": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HOST_HEADER", "CUSTOM"}, false)),
						Optional:         true,
						Description:      "This specifies the domain used to establish a session with the visitor.",
						Type:             schema.TypeString,
					},
					"custom_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^(\\.)?(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"),
						Optional:         true,
						Description:      "This specifies a domain for all session cookies. In case you configure many property hostnames, this may be their common domain. Make sure the user agent accepts the custom domain for any request matching the `virtualWaitingRoom` behavior. Don't use top level domains (TLDs).",
						Type:             schema.TypeString,
					},
					"waiting_room_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"waiting_room_path": {
						Optional:    true,
						Description: "This specifies the path to the waiting room main page on the origin server, for example `/vp/waiting-room.html`. When the request is marked as Waiting Room Main Page and blocked, the visitor enters the waiting room. The behavior sets the outgoing request path to the `waitingRoomPath` and modifies the cache key accordingly. See the `virtualWaitingRoomRequest` match criteria to further customize these requests.",
						Type:        schema.TypeString,
					},
					"waiting_room_assets_paths": {
						Optional:    true,
						Description: "This specifies the base paths to static resources such as JavaScript, CSS, or image files for the Waiting Room Main Page requests. The option supports the `*` wildcard that matches zero or more characters. Requests matching any of these paths aren't blocked, but marked as Waiting Room Assets and passed through to the origin. See the `virtualWaitingRoomRequest` match criteria to further customize these requests.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"access_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"session_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 86400)),
						Optional:         true,
						Description:      "Specifies the number of seconds users remain in the waiting room queue.",
						Type:             schema.TypeInt,
					},
					"session_auto_prolong": {
						Optional:    true,
						Description: "Whether the queue session should prolong automatically when the `sessionDuration` expires  and the visitor remains active.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"virtual_waiting_room_with_edge_workers": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior allows you to configure the `virtualWaitingRoom` behavior with EdgeWorkers for extended scalability and customization. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"visitor_prioritization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The `Visitor Prioritization Cloudlet` decreases abandonment by providing a user-friendly waiting room experience.  With Cloudlets available on your contract, choose `Your services` > `Edge logic Cloudlets` to control Visitor Prioritization within `Control Center`. Otherwise use the `Cloudlets API` to configure it programmatically.  To serve non-HTML API content such as JSON blocks, see the `apiPrioritization` behavior. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the Visitor Prioritization behavior.",
						Type:        schema.TypeBool,
					},
					"cloudlet_policy": {
						Optional:    true,
						Description: "Identifies the Cloudlet policy.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"user_identification_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"user_identification_by_cookie": {
						Optional:    true,
						Description: "When enabled, identifies users by the value of a cookie.",
						Type:        schema.TypeBool,
					},
					"user_identification_key_cookie": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies the name of the cookie whose value identifies users. To match a user, the value of the cookie needs to remain constant across all requests.",
						Type:             schema.TypeString,
					},
					"user_identification_by_headers": {
						Optional:    true,
						Description: "When enabled, identifies users by the values of GET or POST request headers.",
						Type:        schema.TypeBool,
					},
					"user_identification_key_headers": {
						Optional:    true,
						Description: "Specifies names of request headers whose values identify users. To match a user, values for all the specified headers need to remain constant across all requests.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"user_identification_by_ip": {
						Optional:    true,
						Description: "Allows IP addresses to identify users.",
						Type:        schema.TypeBool,
					},
					"user_identification_by_params": {
						Optional:    true,
						Description: "When enabled, identifies users by the values of GET or POST request parameters.",
						Type:        schema.TypeBool,
					},
					"user_identification_key_params": {
						Optional:    true,
						Description: "Specifies names of request parameters whose values identify users. To match a user, values for all the specified parameters need to remain constant across all requests. Parameters that are absent or blank may also identify users.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"allowed_user_cookie_management_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"allowed_user_cookie_enabled": {
						Optional:    true,
						Description: "Sets a cookie for users who have been allowed through to the site.",
						Type:        schema.TypeBool,
					},
					"allowed_user_cookie_label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies a label to distinguish this cookie for an allowed user from others. The value appends to the cookie's name, and helps you to maintain the same user assignment across behaviors within a property, and across properties.",
						Type:             schema.TypeString,
					},
					"allowed_user_cookie_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 600)),
						Optional:         true,
						Description:      "Sets the number of seconds for the allowed user's session once allowed through to the site.",
						Type:             schema.TypeInt,
					},
					"allowed_user_cookie_refresh": {
						Optional:    true,
						Description: "Resets the duration of an allowed cookie with each request, so that it only expires if the user doesn't make any requests for the specified duration. Do not enable this option if you want to set a fixed time for all users.",
						Type:        schema.TypeBool,
					},
					"allowed_user_cookie_advanced": {
						Optional:    true,
						Description: "Sets advanced configuration options for the allowed user's cookie.",
						Type:        schema.TypeBool,
					},
					"allowed_user_cookie_automatic_salt": {
						Optional:    true,
						Description: "Sets an automatic `salt` value to verify the integrity of the cookie for an allowed user. Disable this if you want to share the cookie across properties.",
						Type:        schema.TypeBool,
					},
					"allowed_user_cookie_salt": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies a fixed `salt` value, which is incorporated into the cookie's value to prevent users from manipulating it. You can use the same salt string across different behaviors or properties to apply a single cookie to all allowed users.",
						Type:             schema.TypeString,
					},
					"allowed_user_cookie_domain_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DYNAMIC", "CUSTOMER"}, false)),
						Optional:         true,
						Description:      "Specify with `allowedUserCookieAdvanced` enabled.",
						Type:             schema.TypeString,
					},
					"allowed_user_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^(\\.)?(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"),
						Optional:         true,
						Description:      "Specifies a domain for an allowed user cookie.",
						Type:             schema.TypeString,
					},
					"allowed_user_cookie_http_only": {
						Optional:    true,
						Description: "Applies the `HttpOnly` flag to the allowed user's cookie to ensure it's accessed over HTTP and not manipulated by the client.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cookie_management_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"waiting_room_cookie_enabled": {
						Optional:    true,
						Description: "Enables a cookie to track a waiting room assignment.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cookie_share_label": {
						Optional:    true,
						Description: "Enabling this option shares the same `allowedUserCookieLabel` string. If disabled, specify a different `waitingRoomCookieLabel`.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cookie_label": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies a label to distinguish this waiting room cookie from others. The value appends to the cookie's name, and helps you to maintain the same waiting room assignment across behaviors within a property, and across properties.",
						Type:             schema.TypeString,
					},
					"waiting_room_cookie_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 120)),
						Optional:         true,
						Description:      "Sets the number of seconds for which users remain in the waiting room. During this time, users who refresh the waiting room page remain there.",
						Type:             schema.TypeInt,
					},
					"waiting_room_cookie_advanced": {
						Optional:    true,
						Description: "When enabled along with `waitingRoomCookieEnabled`, sets advanced configuration options for the waiting room cookie.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cookie_automatic_salt": {
						Optional:    true,
						Description: "Sets an automatic `salt` value to verify the integrity of the waiting room cookie.  Disable this if you want to share the cookie across properties.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cookie_salt": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "Specifies a fixed `salt` value, which is incorporated into the cookie's value to prevent users from manipulating it. You can use the same salt string across different behaviors or properties to apply a single cookie for the waiting room session.",
						Type:             schema.TypeString,
					},
					"waiting_room_cookie_domain_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"DYNAMIC", "CUSTOMER"}, false)),
						Optional:         true,
						Description:      "Specify with `waitingRoomCookieAdvanced` enabled, selects whether to use the `DYNAMIC` incoming host header, or a `CUSTOMER`-defined cookie domain.",
						Type:             schema.TypeString,
					},
					"waiting_room_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^(\\.)?(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"),
						Optional:         true,
						Description:      "Specifies a domain for the waiting room cookie.",
						Type:             schema.TypeString,
					},
					"waiting_room_cookie_http_only": {
						Optional:    true,
						Description: "Applies the `HttpOnly` flag to the waiting room cookie to ensure it's accessed over HTTP and not manipulated by the client.",
						Type:        schema.TypeBool,
					},
					"waiting_room_management_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"waiting_room_status_code": {
						ValidateDiagFunc: validateRegexOrVariable("[2|4|5][0-9][0-9]"),
						Optional:         true,
						Description:      "Specifies the response code for requests sent to the waiting room.",
						Type:             schema.TypeInt,
					},
					"waiting_room_use_cp_code": {
						Optional:    true,
						Description: "Allows you to assign a different CP code that tracks any requests that are sent to the waiting room.",
						Type:        schema.TypeBool,
					},
					"waiting_room_cp_code": {
						Optional:    true,
						Description: "Specifies a CP code for requests sent to the waiting room. You only need to provide the initial `id`, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
					"waiting_room_net_storage": {
						Optional:    true,
						Description: "Specifies the NetStorage domain for the waiting room page.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cp_code": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"download_domain_name": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
								"g2o_token": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeString,
								},
							},
						},
					},
					"waiting_room_directory": {
						ValidateDiagFunc: validateRegexOrVariable("^[^#\\[\\]@]+$"),
						Optional:         true,
						Description:      "Specifies the NetStorage directory that contains the static waiting room page, with no trailing slash character.",
						Type:             schema.TypeString,
					},
					"waiting_room_cache_ttl": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(5, 30)),
						Optional:         true,
						Description:      "Specifies the waiting room page's time to live in the cache, `5` minutes by default.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"visitor_prioritization_fifo": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "(**BETA**) The `Visitor Prioritization Cloudlet (FIFO)` decreases abandonment by providing a user-friendly waiting room experience. FIFO (First-in First-out) is a fair request processing mechanism, which prioritizes the first requests that enter the waiting room to send them first to the origin. Users can see both their estimated arrival time and position in the line. With Cloudlets available on your contract, choose `Your services` > `Edge logic Cloudlets` to control Visitor Prioritization (FIFO) within `Control Center`. Otherwise use the `Cloudlets API` to configure it programmatically. To serve non-HTML API content such as JSON blocks, see the `apiPrioritization` behavior. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"cloudlet_shared_policy": {
						Optional:    true,
						Description: "This identifies the Visitor Prioritization FIFO shared policy to use with this behavior. You can list available shared policies with the `Cloudlets API`.",
						Type:        schema.TypeInt,
					},
					"domain_config": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HOST_HEADER", "CUSTOM"}, false)),
						Optional:         true,
						Description:      "This specifies how to set the domain used to establish a session with the visitor.",
						Type:             schema.TypeString,
					},
					"custom_cookie_domain": {
						ValidateDiagFunc: validateRegexOrVariable("^(\\.)?(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)+([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"),
						Optional:         true,
						Description:      "This specifies a domain for all session cookies. In case you configure many property hostnames, this may be their common domain. Make sure the user agent accepts the custom domain for any request matching the `visitorPrioritizationFifo` behavior. Don't use top level domains (TLDs).",
						Type:             schema.TypeString,
					},
					"waiting_room_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"waiting_room_path": {
						Optional:    true,
						Description: "This specifies the path to the waiting room main page on the origin server, for example `/vp/waiting-room.html`. When the request is marked as `Waiting Room Main Page` and blocked, the visitor enters the waiting room. The behavior sets the outgoing request path to the `waitingRoomPath` and modifies the cache key accordingly. See the `visitorPrioritizationRequest` match criteria to further customize these requests.",
						Type:        schema.TypeString,
					},
					"waiting_room_assets_paths": {
						Optional:    true,
						Description: "This specifies the base paths to static resources such as `JavaScript`, `CSS`, or image files for the `Waiting Room Main Page` requests. The option supports the `*` wildcard wildcard that matches zero or more characters. Requests matching any of these paths aren't blocked, but marked as Waiting Room Assets and passed through to the origin. See the `visitorPrioritizationRequest` match criteria to further customize these requests.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"access_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"session_duration": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 86400)),
						Optional:         true,
						Description:      "Specifies the number of seconds users remain in the waiting room queue.",
						Type:             schema.TypeInt,
					},
					"session_auto_prolong": {
						Optional:    true,
						Description: "Whether the queue session should prolong automatically when the `sessionDuration` expires  and the visitor remains active.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"visitor_prioritization_fifo_standalone": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"watermarking": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Adds watermarking for each valid user's content. Content segments are delivered from different sources using a pattern unique to each user, based on a watermarking token included in each request. If your content is pirated or redistributed, you can forensically analyze the segments to extract the pattern, and identify the user who leaked the content. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enable": {
						Optional:    true,
						Description: "Enables the watermarking behavior.",
						Type:        schema.TypeBool,
					},
					"token_signing_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"signature_verification_enable": {
						Optional:    true,
						Description: "When enabled, you can verify the signature in your watermarking token.",
						Type:        schema.TypeBool,
					},
					"verification_key_id1": {
						Optional:    true,
						Description: "Specifies a unique identifier for the first public key.",
						Type:        schema.TypeString,
					},
					"verification_public_key1": {
						Optional:    true,
						Description: "Specifies the first public key in its entirety.",
						Type:        schema.TypeString,
					},
					"verification_key_id2": {
						Optional:    true,
						Description: "Specifies a unique identifier for the optional second public key.",
						Type:        schema.TypeString,
					},
					"verification_public_key2": {
						Optional:    true,
						Description: "Specifies the optional second public key in its entirety. Specify a second key to enable rotation.",
						Type:        schema.TypeString,
					},
					"pattern_encryption_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"pattern_decryption_enable": {
						Optional:    true,
						Description: "If patterns in your watermarking tokens have been encrypted, enabling this allows you to provide values to decrypt them.",
						Type:        schema.TypeBool,
					},
					"decryption_password_id1": {
						Optional:    true,
						Description: "Specifies a label that corresponds to the primary password.",
						Type:        schema.TypeString,
					},
					"decryption_password1": {
						Optional:    true,
						Description: "Provides the primary password used to encrypt patterns in your watermarking tokens.",
						Type:        schema.TypeString,
					},
					"decryption_password_id2": {
						Optional:    true,
						Description: "Specifies a label for the secondary password, used in rotation scenarios to identify which password to use for decryption.",
						Type:        schema.TypeString,
					},
					"decryption_password2": {
						Optional:    true,
						Description: "Provides the secondary password you can use to rotate passwords.",
						Type:        schema.TypeString,
					},
					"miscellaneous_settings_title": {
						Optional:    true,
						Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
						Type:        schema.TypeString,
					},
					"use_original_as_a": {
						Optional:    true,
						Description: "When you work with your watermarking vendor, you can apply several preprocessing methods to your content. See the `AMD help` for more information. With the standard `filename-prefix AB naming` preprocessing method, the watermarking vendor creates two variants of the original segment content and labels them as an `A` and `B` segment in the filename. If you selected the `unlabeled A variant` preprocessing method, enabling this option tells your configuration to use the original filename segment content as your `A` variant.",
						Type:        schema.TypeBool,
					},
					"ab_variant_location": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"FILENAME_PREFIX", "PARENT_DIRECTORY_PREFIX"}, false)),
						Optional:         true,
						Description:      "When you work with your watermarking vendor, you can apply several preprocessing methods to your content. See the `AMD help` for more information. Use this option to specify the location of the `A` and `B` variant segments.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"web_application_firewall": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This behavior implements a suite of security features that blocks threatening HTTP and HTTPS requests. Use it as your primary firewall, or in addition to existing security measures.  Only one referenced configuration is allowed per property, so this behavior typically belongs as part of its default rule. This behavior cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"firewall_configuration": {
						Optional:    true,
						Description: "An object featuring details about your firewall configuration.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"config_id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"production_status": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior",
									Type:        schema.TypeString,
								},
								"staging_status": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior",
									Type:        schema.TypeString,
								},
								"production_version": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior",
									Type:        schema.TypeInt,
								},
								"staging_version": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior",
									Type:        schema.TypeInt,
								},
								"file_name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior",
									Type:        schema.TypeString,
								},
							},
						},
					},
				},
			},
		},
		"web_sockets": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The WebSocket protocol allows web applications real-time bidirectional communication between clients and servers. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables WebSocket traffic.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"webdav": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Web-based Distributed Authoring and Versioning (WebDAV) is a set of extensions to the HTTP protocol that allows users to collaboratively edit and manage files on remote web servers. This behavior enables WebDAV, and provides support for the following additional request methods: PROPFIND, PROPPATCH, MKCOL, COPY, MOVE, LOCK, and UNLOCK. To apply this behavior, you need to match on a `requestMethod`. This behavior can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"enabled": {
						Optional:    true,
						Description: "Enables the WebDAV behavior.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
	}
}

func getCriteriaSchemaV20241021() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"advanced_im_match": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches whether the `imageManager` behavior already applies to the current set of requests. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specifies the match's logic.",
						Type:             schema.TypeString,
					},
					"match_on": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ANY_IM", "PRISTINE"}, false)),
						Optional:         true,
						Description:      "Specifies the match's scope.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"bucket": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This matches a specified percentage of requests when used with the accompanying behavior. Contact Akamai Professional Services for help configuring it. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"percentage": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specifies the percentage of requests to match.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"cacheability": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the current cache state.  Note that any `NO_STORE` or `BYPASS_CACHE` HTTP headers set on the origin's content overrides properties' `caching` instructions, in which case this criteria does not apply. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specifies the match's logic.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NO_STORE", "BYPASS_CACHE", "CACHEABLE"}, false)),
						Optional:         true,
						Description:      "Content's cache is enabled (`CACHEABLE`) or not (`NO_STORE`), or else is ignored (`BYPASS_CACHE`).",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"china_cdn_region": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Identifies traffic deployed over Akamai's regional ChinaCDN infrastructure. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specify whether the request `IS` or `IS_NOT` deployed over ChinaCDN.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"client_certificate": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches whether you have configured a client certificate to authenticate requests to edge servers. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"is_certificate_present": {
						Optional:    true,
						Description: "Executes rule behaviors only if a client certificate authenticates requests.",
						Type:        schema.TypeBool,
					},
					"is_certificate_valid": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"VALID", "INVALID", "IGNORE"}, false)),
						Optional:         true,
						Description:      "Matches whether the certificate is `VALID` or `INVALID`. You can also `IGNORE` the certificate's validity.",
						Type:             schema.TypeString,
					},
					"enforce_mtls": {
						Optional:    true,
						Description: "Specifies custom handling of requests if any of the checks in the `enforceMtlsSettings` behavior fail. Enable this and use with behaviors such as `logCustom` so that they execute if the check fails. You need to add the `enforceMtlsSettings` behavior to a parent rule, with its own unique match condition and `enableDenyRequest` option disabled.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"client_ip": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the IP number of the requesting client. To use this condition to match end-user IP addresses, apply it together with the `requestType` matching on the `CLIENT_REQ` value. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the contents of `values` if set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "IP or CIDR block, for example: `71.92.0.0/14`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"use_headers": {
						Optional:    true,
						Description: "When connecting via a proxy server as determined by the `X-Forwarded-For` header, enabling this option matches the connecting client's IP address rather than the original end client specified in the header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"client_ip_version": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the version of the IP protocol used by the requesting client. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IPV4", "IPV6"}, false)),
						Optional:         true,
						Description:      "The IP version of the client request, either `IPV4` or `IPV6`.",
						Type:             schema.TypeString,
					},
					"use_x_forwarded_for": {
						Optional:    true,
						Description: "When connecting via a proxy server as determined by the `X-Forwarded-For` header, enabling this option matches the connecting client's IP address rather than the original end client specified in the header.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"cloudlets_origin": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Allows Cloudlets Origins, referenced by label, to define their own criteria to assign custom origin definitions. The criteria may match, for example, for a specified percentage of requests defined by the cloudlet to use an alternative version of a website. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"origin_id": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-\\.]+$"),
						Optional:         true,
						Description:      "The Cloudlets Origins identifier, limited to alphanumeric and underscore characters.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_delivery_network": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies the type of Akamai network handling the request. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Matches the specified `network` if set to `IS`, otherwise `IS_NOT` reverses the match.",
						Type:             schema.TypeString,
					},
					"network": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"STAGING", "PRODUCTION"}, false)),
						Optional:         true,
						Description:      "Match the network.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"content_type": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the HTTP response header's `Content-Type`. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches any `Content-Type` among specified `values` when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "`Content-Type` response header value, for example `text/html`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_wildcard": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters. Specifying `text/*` matches both `text/html` and `text/css`.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match for all `values`.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"debug_mode": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The match applies when the request is debugged using the `enhancedDebug` behavior. This criterion is for internal usage only. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"debug_mode": {
						Optional:    true,
						Description: "Whether the request is being debugged using Enhanced Debug.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"device_characteristic": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match various aspects of the device or browser making the request. Based on the value of the `characteristic` option, the expected value is either a boolean, a number, or a string, possibly representing a version number. Each type of value requires a different field. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"characteristic": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BRAND_NAME", "MODEL_NAME", "MARKETING_NAME", "IS_WIRELESS_DEVICE", "IS_TABLET", "DEVICE_OS", "DEVICE_OS_VERSION", "MOBILE_BROWSER", "MOBILE_BROWSER_VERSION", "RESOLUTION_WIDTH", "RESOLUTION_HEIGHT", "PHYSICAL_SCREEN_HEIGHT", "PHYSICAL_SCREEN_WIDTH", "COOKIE_SUPPORT", "AJAX_SUPPORT_JAVASCRIPT", "FULL_FLASH_SUPPORT", "ACCEPT_THIRD_PARTY_COOKIE", "XHTML_SUPPORT_LEVEL", "IS_MOBILE"}, false)),
						Optional:         true,
						Description:      "Aspect of the device or browser to match.",
						Type:             schema.TypeString,
					},
					"string_match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"MATCHES_ONE_OF", "DOES_NOT_MATCH_ONE_OF"}, false)),
						Optional:         true,
						Description:      "When the `characteristic` expects a string value, set this to `MATCHES_ONE_OF` to match against the `stringValue` set, otherwise set to `DOES_NOT_MATCH_ONE_OF` to exclude that set of values.",
						Type:             schema.TypeString,
					},
					"numeric_match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT", "IS_LESS_THAN", "IS_LESS_THAN_OR_EQUAL", "IS_MORE_THAN", "IS_MORE_THAN_OR_EQUAL"}, false)),
						Optional:         true,
						Description:      "When the `characteristic` expects a numeric value, compares the specified `numericValue` against the matched client.",
						Type:             schema.TypeString,
					},
					"version_match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT", "IS_LESS_THAN", "IS_LESS_THAN_OR_EQUAL", "IS_MORE_THAN", "IS_MORE_THAN_OR_EQUAL"}, false)),
						Optional:         true,
						Description:      "When the `characteristic` expects a version string value, compares the specified `versionValue` against the matched client, using the following operators: `IS`, `IS_MORE_THAN_OR_EQUAL`, `IS_MORE_THAN`, `IS_LESS_THAN_OR_EQUAL`, `IS_LESS_THAN`, `IS_NOT`.",
						Type:             schema.TypeString,
					},
					"boolean_value": {
						Optional:    true,
						Description: "When the `characteristic` expects a boolean value, this specifies the value.",
						Type:        schema.TypeBool,
					},
					"string_value": {
						Optional:    true,
						Description: "When the `characteristic` expects a string, this specifies the set of values.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"numeric_value": {
						Optional:    true,
						Description: "When the `characteristic` expects a numeric value, this specifies the number.",
						Type:        schema.TypeInt,
					},
					"version_value": {
						Optional:    true,
						Description: "When the `characteristic` expects a version number, this specifies it as a string.",
						Type:        schema.TypeString,
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `stringValue` field.",
						Type:        schema.TypeBool,
					},
					"match_wildcard": {
						Optional:    true,
						Description: "Allows wildcards in the `stringValue` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"ecmd_auth_groups": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CONTAINS", "DOES_NOT_CONTAIN"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,255}$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"ecmd_auth_scheme": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"auth_scheme": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ANONYMOUS", "JWT", "MUTUAL"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"ecmd_is_authenticated": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_AUTHENTICATED", "IS_NOT_AUTHENTICATED"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"ecmd_username": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CONTAINS", "DOES_NOT_CONTAIN", "STARTS_WITH", "DOES_NOT_START_WITH", "ENDS_WITH", "DOES_NOT_END_WITH", "LENGTH_EQUALS", "LENGTH_GREATER_THAN", "LENGTH_SMALLER_THAN"}, false)),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_-]{1,255}$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
					"length": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"edge_workers_failure": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Checks the EdgeWorkers execution status and detects whether a customer's JavaScript failed on edge servers. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"exec_status": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"FAILURE", "SUCCESS"}, false)),
						Optional:         true,
						Description:      "Specify execution status.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"file_extension": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the requested filename's extension, if present. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the contents of `values` if set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "An array of file extension strings, with no leading dot characters, for example `png`, `jpg`, `jpeg`, and `gif`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"filename": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the requested filename, or test whether it is present. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF", "IS_EMPTY", "IS_NOT_EMPTY"}, false)),
						Optional:         true,
						Description:      "If set to `IS_ONE_OF` or `IS_NOT_ONE_OF`, matches whether the filename matches one of the `values`. If set to `IS_EMPTY` or `IS_NOT_EMPTY`, matches whether the specified filename is part of the path.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "Matches the filename component of the request URL. Allows wildcards, where `?` matches a single character and `*` matches zero or more characters. For example, specify `filename.*` to accept any extension.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `values` field.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"hostname": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the requested hostname. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the contents of `values` when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "A list of hostnames. Allows wildcards, where `?` matches a single character and `*` matches zero or more characters. Specifying `*.example.com` matches both `m.example.com` and `www.example.com`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"match_advanced": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "This specifies match criteria using Akamai XML metadata. It can only be configured on your behalf by Akamai Professional Services. This criterion is for internal usage only. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"description": {
						Optional:    true,
						Description: "A human-readable description of what the XML block does.",
						Type:        schema.TypeString,
					},
					"open_xml": {
						Optional:    true,
						Description: "An XML string that opens the relevant block.",
						Type:        schema.TypeString,
					},
					"close_xml": {
						Optional:    true,
						Description: "An XML string that closes the relevant block.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"match_cp_code": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match the assigned content provider code. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						Optional:    true,
						Description: "Specifies the CP code as an object. You only need to provide the initial `id` to match the CP code, stripping any `cpc_` prefix to pass the integer to the rule tree. Additional CP code details may reflect back in subsequent read-only data.",
						Type:        schema.TypeList,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Optional:    true,
									Description: "",
									Type:        schema.TypeInt,
								},
								"name": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"created_date": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeInt,
								},
								"description": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeString,
								},
								"products": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cp_code_limits": {
									Optional:    true,
									Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
									Type:        schema.TypeList,
									MaxItems:    1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"current_capacity": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeInt,
											},
											"limit_type": {
												Optional:    true,
												Description: "This field is only intended for export compatibility purposes, and modifying it will not impact your use of the behavior.",
												Type:        schema.TypeString,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"match_response_code": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match a set or range of HTTP response codes. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF", "IS_BETWEEN", "IS_NOT_BETWEEN"}, false)),
						Optional:         true,
						Description:      "Matches numeric range or a specified set of `values`.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "A set of response codes to match, for example `[\"404\",\"500\"]`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"lower_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d{3}$"),
						Optional:         true,
						Description:      "Specifies the start of a range of responses. For example, `400` to match anything from `400` to `500`.",
						Type:             schema.TypeInt,
					},
					"upper_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^\\d{3}$"),
						Optional:         true,
						Description:      "Specifies the end of a range of responses. For example, `500` to match anything from `400` to `500`.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"match_variable": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches a built-in variable, or a custom variable pre-declared within the rule tree by the `setVariable` behavior.  See `Support for variables` for more information on this feature. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"variable_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z_][a-zA-Z0-9_]{0,31}$"),
						Optional:         true,
						Description:      "The name of the variable to match.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT", "IS_ONE_OF", "IS_NOT_ONE_OF", "IS_EMPTY", "IS_NOT_EMPTY", "IS_BETWEEN", "IS_NOT_BETWEEN", "IS_GREATER_THAN", "IS_GREATER_THAN_OR_EQUAL_TO", "IS_LESS_THAN", "IS_LESS_THAN_OR_EQUAL_TO"}, false)),
						Optional:         true,
						Description:      "The type of match, based on which you use different options to specify the match criteria.",
						Type:             schema.TypeString,
					},
					"variable_values": {
						Optional:    true,
						Description: "Specifies an array of matching strings.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"variable_expression": {
						Optional:    true,
						Description: "Specifies a single matching string.",
						Type:        schema.TypeString,
					},
					"lower_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "Specifies the range's numeric minimum value.",
						Type:             schema.TypeString,
					},
					"upper_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "Specifies the range's numeric maximum value.",
						Type:             schema.TypeString,
					},
					"match_wildcard": {
						Optional:    true,
						Description: "When matching string expressions, enabling this allows wildcards, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "When matching string expressions, enabling this performs a case-sensitive match.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"metadata_stage": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches how the current rule corresponds to low-level syntax elements in translated XML metadata, indicating progressive stages as each edge server handles the request and response.  To use this match, you need to be thoroughly familiar with how Akamai edge servers process requests. Contact your Akamai Technical representative if you need help, and test thoroughly on staging before activating on production. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Compares the current rule with the specified metadata stage.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"cache-hit", "client-done", "client-request", "client-request-body", "client-response", "content-policy", "forward-request", "forward-response", "forward-start", "ipa-response"}, false)),
						Optional:         true,
						Description:      "Specifies the metadata stage.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"origin_timeout": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches when the origin responds with a timeout error. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ORIGIN_TIMED_OUT"}, false)),
						Optional:         true,
						Description:      "Specifies a single required `ORIGIN_TIMED_OUT` value.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"path": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the URL's non-hostname path component. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"MATCHES_ONE_OF", "DOES_NOT_MATCH_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the contents of the `values` array.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "Matches the URL path, excluding leading hostname and trailing query parameters. The path is relative to the server root, for example `/blog`. This field allows wildcards, where `?` matches a single character and `*` matches zero or more characters. For example, `/blog/*/2014` matches paths with two fixed segments and other varying segments between them.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match.",
						Type:        schema.TypeBool,
					},
					"normalize": {
						Optional:    true,
						Description: "Transforms URLs before comparing them with the provided value. URLs are decoded, and any directory syntax such as `../..` or `//` is stripped as a security measure. This protects URL paths from being accessed by unauthorized users.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"query_string_parameter": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches query string field names or values. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"parameter_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^:/?#\\[\\]@&]+$"),
						Optional:         true,
						Description:      "The name of the query field, for example, `q` in `?q=string`.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF", "EXISTS", "DOES_NOT_EXIST", "IS_LESS_THAN", "IS_MORE_THAN", "IS_BETWEEN"}, false)),
						Optional:         true,
						Description:      "Narrows the match criteria.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "The value of the query field, for example, `string` in `?q=string`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"lower_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "Specifies the match's minimum value.",
						Type:             schema.TypeInt,
					},
					"upper_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "When the `value` is numeric, this field specifies the match's maximum value.",
						Type:             schema.TypeInt,
					},
					"match_wildcard_name": {
						Optional:    true,
						Description: "Allows wildcards in the `parameterName` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_name": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `parameterName` field.",
						Type:        schema.TypeBool,
					},
					"match_wildcard_value": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_value": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `value` field.",
						Type:        schema.TypeBool,
					},
					"escape_value": {
						Optional:    true,
						Description: "Matches when the `value` is URL-escaped.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"random": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches a specified percentage of requests. Use this match to apply behaviors to a percentage of your incoming requests that differ from the remainder, useful for A/b testing, or to offload traffic onto different servers. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"bucket": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
						Optional:         true,
						Description:      "Specify a percentage of random requests to which to apply a behavior. Any remainders do not match.",
						Type:             schema.TypeInt,
					},
				},
			},
		},
		"recovery_config": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches on specified origin recovery scenarios. The `originFailureRecoveryPolicy` behavior defines the scenarios that trigger the recovery or retry methods you set in the `originFailureRecoveryMethod` rule. If the origin fails, the system checks the name of the recovery method applied to your policy. It then either redirects the requesting client to a backup origin or returns predefined HTTP response codes. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"config_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[A-Z0-9-]+$"),
						Optional:         true,
						Description:      "A unique identifier used for origin failure recovery configurations. This is the recovery method configuration name you apply when setting origin failure recovery methods and scenarios in `originFailureRecoveryMethod` and `originFailureRecoveryPolicy` behaviors. The value can contain alphanumeric characters and dashes.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"regular_expression": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches a regular expression against a string, especially to apply behaviors flexibly based on the contents of dynamic `variables`. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_string": {
						Optional:    true,
						Description: "The string to match, typically the contents of a dynamic variable.",
						Type:        schema.TypeString,
					},
					"regex": {
						Optional:    true,
						Description: "The regular expression (PCRE) to match against the string.",
						Type:        schema.TypeString,
					},
					"case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive regular expression match.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"request_cookie": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match the cookie name or value passed with the request. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"cookie_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[a-zA-Z0-9_\\-*\\.]+$"),
						Optional:         true,
						Description:      "The name of the cookie, which can be a variable. For example, `visitor` in `visitor:anon`.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT", "EXISTS", "DOES_NOT_EXIST", "IS_BETWEEN"}, false)),
						Optional:         true,
						Description:      "Narrows the match criteria.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validateRegexOrVariable("^[^\\s;]+$"),
						Optional:         true,
						Description:      "The cookie's value, which can be a variable. For example, `anon` in `visitor:anon`.",
						Type:             schema.TypeString,
					},
					"lower_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "When the `value` is numeric, this field specifies the match's minimum value.",
						Type:             schema.TypeInt,
					},
					"upper_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[1-9]\\d*$"),
						Optional:         true,
						Description:      "When the `value` is numeric, this field specifies the match's maximum value.",
						Type:             schema.TypeInt,
					},
					"match_wildcard_name": {
						Optional:    true,
						Description: "Allows wildcards in the `cookieName` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_name": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `cookieName` field.",
						Type:        schema.TypeBool,
					},
					"match_wildcard_value": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_value": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `value` field.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"request_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match HTTP header names or values. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "The name of the request header, for example `Accept-Language`.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF", "EXISTS", "DOES_NOT_EXIST"}, false)),
						Optional:         true,
						Description:      "Narrows the match criteria.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "The request header's value, for example `en-US` when the header `headerName` is `Accept-Language`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_wildcard_name": {
						Optional:    true,
						Description: "Allows wildcards in the `headerName` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_wildcard_value": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_value": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `value` field.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"request_method": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specify the request's HTTP verb. Also supports WebDAV methods and common Akamai operations. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Matches the `value` when set to `IS`, otherwise `IS_NOT` reverses the match.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"GET", "POST", "HEAD", "PUT", "PATCH", "HTTP_DELETE", "AKAMAI_TRANSLATE", "AKAMAI_PURGE", "OPTIONS", "DAV_ACL", "DAV_CHECKOUT", "DAV_COPY", "DAV_DMCREATE", "DAV_DMINDEX", "DAV_DMMKPATH", "DAV_DMMKPATHS", "DAV_DMOVERLAY", "DAV_DMPATCHPATHS", "DAV_LOCK", "DAV_MKCALENDAR", "DAV_MKCOL", "DAV_MOVE", "DAV_PROPFIND", "DAV_PROPPATCH", "DAV_REPORT", "DAV_SETPROCESS", "DAV_SETREDIRECT", "DAV_TRUTHGET", "DAV_UNLOCK"}, false)),
						Optional:         true,
						Description:      "Any of these HTTP methods,  WebDAV methods, or Akamai operations.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"request_protocol": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches whether the request uses the HTTP or HTTPS protocol. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"HTTP", "HTTPS"}, false)),
						Optional:         true,
						Description:      "Specifies the protocol.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"request_type": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the basic type of request. To use this match, you need to be thoroughly familiar with how Akamai edge servers process requests. Contact your Akamai Technical representative if you need help, and test thoroughly on staging before activating on production. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specifies whether the request `IS` or `IS_NOT` the type of specified `value`.",
						Type:             schema.TypeString,
					},
					"value": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"CLIENT_REQ", "ESI_FRAGMENT", "EW_SUBREQUEST"}, false)),
						Optional:         true,
						Description:      "Specifies the type of request, either a standard `CLIENT_REQ`, an `ESI_FRAGMENT`, or an `EW_SUBREQUEST`.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"response_header": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match HTTP header names or values. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"header_name": {
						ValidateDiagFunc: validateRegexOrVariable("^[^()<>@,;:\\\"/\\[\\]?{}\\s]+$"),
						Optional:         true,
						Description:      "The name of the response header, for example `Content-Type`.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF", "EXISTS", "DOES_NOT_EXIST", "IS_LESS_THAN", "IS_MORE_THAN", "IS_BETWEEN"}, false)),
						Optional:         true,
						Description:      "Narrows the match according to various criteria.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "The response header's value, for example `application/x-www-form-urlencoded` when the header `headerName` is `Content-Type`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"lower_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "When the `value` is numeric, this field specifies the match's minimum value.",
						Type:             schema.TypeInt,
					},
					"upper_bound": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+$"),
						Optional:         true,
						Description:      "When the `value` is numeric, this field specifies the match's maximum value.",
						Type:             schema.TypeInt,
					},
					"match_wildcard_name": {
						Optional:    true,
						Description: "Allows wildcards in the `headerName` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_wildcard_value": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive_value": {
						Optional:    true,
						Description: "When enabled, the match is case-sensitive for the `value` field.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"server_location": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The location of the Akamai server handling the request. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"location_type": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COUNTRY", "CONTINENT", "REGION"}, false)),
						Optional:         true,
						Description:      "Indicates the geographic scope.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the specified set of values when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"countries": {
						Optional:    true,
						Description: "ISO 3166-1 country codes, such as `US` or `CN`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"continents": {
						Optional:    true,
						Description: "Continent codes.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"regions": {
						Optional:    true,
						Description: "ISO 3166 country and region codes, for example `US:MA` for Massachusetts or `JP:13` for Tokyo.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"time": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Specifies ranges of times during which the request occurred. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BEGINNING", "BETWEEN", "LASTING", "REPEATING"}, false)),
						Optional:         true,
						Description:      "Specifies how to define the range of time.",
						Type:             schema.TypeString,
					},
					"repeat_interval": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Sets the time between each repeating time period's starting points.",
						Type:             schema.TypeString,
					},
					"repeat_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Sets the duration of each repeating time period.",
						Type:             schema.TypeString,
					},
					"lasting_duration": {
						ValidateDiagFunc: validateRegexOrVariable("^[0-9]+[DdHhMmSs]$"),
						Optional:         true,
						Description:      "Specifies the end of a time period as a duration relative to the `lastingDate`.",
						Type:             schema.TypeString,
					},
					"lasting_date": {
						Optional:    true,
						Description: "Sets the start of a fixed time period.",
						Type:        schema.TypeString,
					},
					"repeat_begin_date": {
						Optional:    true,
						Description: "Sets the start of the initial time period.",
						Type:        schema.TypeString,
					},
					"apply_daylight_savings_time": {
						Optional:    true,
						Description: "Adjusts the start time plus repeat interval to account for daylight saving time. Applies when the current time and the start time use different systems, daylight and standard, and the two values are in conflict.",
						Type:        schema.TypeBool,
					},
					"begin_date": {
						Optional:    true,
						Description: "Sets the start of a time period.",
						Type:        schema.TypeString,
					},
					"end_date": {
						Optional:    true,
						Description: "Sets the end of a fixed time period.",
						Type:        schema.TypeString,
					},
				},
			},
		},
		"token_authorization": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Match on Auth Token 2.0 verification results. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_SUCCESS", "IS_CUSTOM_FAILURE", "IS_ANY_FAILURE"}, false)),
						Optional:         true,
						Description:      "Error match scope.",
						Type:             schema.TypeString,
					},
					"status_list": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"user_agent": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches the user agent string that helps identify the client browser and device. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the specified set of `values` when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"values": {
						Optional:    true,
						Description: "The `User-Agent` header's value. For example, `Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"match_wildcard": {
						Optional:    true,
						Description: "Allows wildcards in the `value` field, where `?` matches a single character and `*` matches zero or more characters. For example, `*Android*`, `*iPhone5*`, `*Firefox*`, or `*Chrome*` allow substring matches.",
						Type:        schema.TypeBool,
					},
					"match_case_sensitive": {
						Optional:    true,
						Description: "Sets a case-sensitive match for the `value` field.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"user_location": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "The client browser's approximate geographic location, determined by looking up the IP address in a database. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"field": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"COUNTRY", "CONTINENT", "REGION"}, false)),
						Optional:         true,
						Description:      "Indicates the geographic scope.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the specified set of values when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"country_values": {
						Optional:    true,
						Description: "ISO 3166-1 country codes, such as `US` or `CN`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"continent_values": {
						Optional:    true,
						Description: "Continent codes.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"region_values": {
						Optional:    true,
						Description: "ISO 3166 country and region codes, for example `US:MA` for Massachusetts or `JP:13` for Tokyo.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"check_ips": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BOTH", "CONNECTING", "HEADERS"}, false)),
						Optional:         true,
						Description:      "Specifies which IP addresses determine the user's location.",
						Type:             schema.TypeString,
					},
					"use_only_first_x_forwarded_for_ip": {
						Optional:    true,
						Description: "When connecting via a proxy server as determined by the `X-Forwarded-For` header, enabling this option matches the end client specified in the header. Disabling it matches the connecting client's IP address.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"user_network": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches details of the network over which the request was made, determined by looking up the IP address in a database. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"field": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"NETWORK", "NETWORK_TYPE", "BANDWIDTH"}, false)),
						Optional:         true,
						Description:      "The type of information to match.",
						Type:             schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS_ONE_OF", "IS_NOT_ONE_OF"}, false)),
						Optional:         true,
						Description:      "Matches the specified set of values when set to `IS_ONE_OF`, otherwise `IS_NOT_ONE_OF` reverses the match.",
						Type:             schema.TypeString,
					},
					"network_type_values": {
						Optional:    true,
						Description: "",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"network_values": {
						Optional:    true,
						Description: "Any set of specific networks.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"bandwidth_values": {
						Optional:    true,
						Description: "Bandwidth range in bits per second, either `1`, `57`, `257`, `1000`, `2000`, or `5000`.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"check_ips": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"BOTH", "CONNECTING", "HEADERS"}, false)),
						Optional:         true,
						Description:      "Specifies which IP addresses determine the user's network.",
						Type:             schema.TypeString,
					},
					"use_only_first_x_forwarded_for_ip": {
						Optional:    true,
						Description: "When connecting via a proxy server as determined by the `X-Forwarded-For` header, enabling this option matches the end client specified in the header. Disabling it matches the connecting client's IP address.",
						Type:        schema.TypeBool,
					},
				},
			},
		},
		"variable_error": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Matches any runtime errors that occur on edge servers based on the configuration of a `setVariable` behavior. See `Support for variables` section for more information on this feature. This criterion can be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"result": {
						Optional:    true,
						Description: "Matches errors for the specified set of `variableNames`, otherwise matches errors from variables outside that set.",
						Type:        schema.TypeBool,
					},
					"variable_names": {
						Optional:    true,
						Description: "The name of the variable whose error triggers the match, or a space- or comma-delimited list of more than one variable name. Note that if you define a variable named `VAR`, the name in this field needs to appear with its added prefix as `PMUSER_VAR`. When such a variable is inserted into other fields, it appears with an additional namespace as `{{user.PMUSER_VAR}}`. See the `setVariable` behavior for details on variable names.",
						Type:        schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"virtual_waiting_room_request": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Helps to customize the requests identified by the `virtualWaitingRoom` behavior. Use this match criteria to define the `originServer` behavior for the waiting room. This criterion cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specifies the match's logic.",
						Type:             schema.TypeString,
					},
					"match_on": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"WR_ANY", "WR_MAIN_PAGE", "WR_ASSETS"}, false)),
						Optional:         true,
						Description:      "Specifies the type of request identified by the `virtualWaitingRoom` behavior.",
						Type:             schema.TypeString,
					},
				},
			},
		},
		"visitor_prioritization_request": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Helps to customize the requests identified by the `visitorPrioritizationFifo` behavior. The basic use case for this match criteria is to define the `originServer` behavior for the waiting room. This criterion cannot be used in includes.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"locked": {
						Optional:    true,
						Description: "Indicates that your Akamai representative has locked this behavior or criteria so that you can't modify it. This option is for internal usage only.",
						Type:        schema.TypeBool,
					},
					"uuid": {
						ValidateDiagFunc: validateRegex("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"),
						Optional:         true,
						Description:      "A uuid member indicates that at least one of its component behaviors or criteria is advanced and read-only. You need to preserve this uuid as well when modifying the rule tree. This option is for internal usage only.",
						Type:             schema.TypeString,
					},
					"template_uuid": {
						Optional:    true,
						Description: "This option is for internal usage only.",
						Type:        schema.TypeString,
					},
					"match_operator": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"IS", "IS_NOT"}, false)),
						Optional:         true,
						Description:      "Specifies the match's logic.",
						Type:             schema.TypeString,
					},
					"match_on": {
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"WR_ANY", "WR_MAIN_PAGE", "WR_ASSETS"}, false)),
						Optional:         true,
						Description:      "Specifies the type of request identified by the `visitorPrioritizationFifo` behavior.",
						Type:             schema.TypeString,
					},
				},
			},
		},
	}
}

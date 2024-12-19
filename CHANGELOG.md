# RELEASE NOTES

## 6.6.1 (Dec 20, 2024)

#### FEATURES/ENHANCEMENTS:

* Global
  * Updated various dependencies

## 6.6.0 (Nov 21, 2024)

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Fixed a problem with the missing `security_policy_id` during update if a resource was imported previously. 
  * Added the `akamai_appsec_aap_selected_hostnames` resource and data source.
  * Modified the `enable_botman_siem` field from `Required` to the `Optional` parameter in the `akamai_appsec_siem_settings` resource.

* Cloud Access
  * Added functionality to import the `akamai_cloudaccess_key` resource for specified group and contract IDs.
  * Marked the `cloud_secret_access_key` field as a sensitive value in the `akamai_cloudaccess_key` resource ([I#580](https://github.com/akamai/terraform-provider-akamai/issues/580)).

* CPS
  * Refreshed a list of warnings returned by the `akamai_cps_warnings` data source.
  
* DNS
  * Added the new `outbound_zone_transfer` field to the `akamai_dns_zone` resource.
  
* Edgeworkers
  * Stopped sending an EdgeKV initialization request in the `akamai_edgekv` resource when EdgeKV is already initialized. ([I#589](https://github.com/akamai/terraform-provider-akamai/issues/589))

* PAPI
  * Added support for the new rule format `v2024-10-21`.

#### BUG FIXES:

* Appsec
  * Fixed a plug-in crash if the `exceptions` block is passed as empty in the `akamai_appsec_siem_settings` resource.

* Cloud Access
  * Resolved issues with drift detection after deleting a key version in the `akamai_cloudaccess_key` resource ([I#579](https://github.com/akamai/terraform-provider-akamai/issues/579)).
  * Fixed cases where ProcessingType = "FAILED" was received in a response from the `akamai_cloudaccess_key` resource. This was causing unnecessary pooling until the timeout.
  
* GTM
  * Added checks to verify the existence of specific objects on the server when creating these resources:
    * `akamai_gtm_asmap` 
    * `akamai_gtm_cidrmap` 
    * `akamai_gtm_domain` 
    * `akamai_gtm_geomap` 
    * `akamai_gtm_property` 
    * `akamai_gtm_resource`

* PAPI
  * Fixed an idempotency issue in property activation when `rule_errors` is empty.
  * Fixed an issue when timeout in the `akamai_property_activation` resource would terminate with the `Provider produced inconsistent result after apply` error.
    * Changed a timeout message from a warning to an error in the `akamai_property_activation` and `akamai_property_include_activation` resources.
  * Fixed an import of the `akamai_property_include` resource to properly populate the include's `product_id` field ([I#575](https://github.com/akamai/terraform-provider-akamai/issues/575)).

#### DEPRECATIONS

* Appsec
  * Deprecated the `akamai_appsec_wap_selected_hostnames` data source and resource. Use the `akamai_appsec_aap_selected_hostnames` data source and resource instead.

## 6.5.0 (Oct 10, 2024)

#### FEATURES/ENHANCEMENTS:

* Global
  * Migrated Terraform to version 1.9.5.
  * Updated SDK v2 and framework libraries.

* Appsec
  * Added the `exceptions` block to the `akamai_appsec_siem_settings` resource with these nested fields:
    * `api_request_constraints`
    * `apr_protection`
    * `bot_management`
    * `client_rep`
    * `custom_rules`
    * `ip_geo`
    * `malware_protection`
    * `rate`
    * `url_protection`
    * `slow_post`
    * `waf`

* GTM
  * Added the retry logic to the `akamai_gtm_property` resource to handle errors caused by the prolonged creation time, leading to Property Validation Failure with the "no datacenter is assigned to map target (all others)" error from the API.
 
* IAM
  * Added new data sources:
    * `akamai_iam_accessible_groups` - reads the groups and subgroups accessible for a given user.
    * `akamai_iam_account_switch_keys` - reads the account switch keys.
    * `akamai_iam_allowed_apis` - reads the list of APIs available to a given user.
    * `akamai_iam_authorized_users` - reads the list of authorized API client users.
    * `akamai_iam_blocked_properties` - reads blocked properties for a certain user in a group.
    * `akamai_iam_cidr_block` - reads details of a specified CIDR block.
    * `akamai_iam_cidr_blocks` - lists all CIDR blocks available to you on your allowlist.
    * `akamai_iam_group` - reads details about a given group and any of its subgroups.
    * `akamai_iam_password_policy` - reads the password policy parameters.
    * `akamai_iam_property_users` - lists users for a given property or include.
    * `akamai_iam_role` - reads details of a specified role.
    * `akamai_iam_user` - reads details of a specific user account.
    * `akamai_iam_users` - lists all users with access to your account.
    * `akamai_iam_users_affected_by_moving_group` - lists the users affected by moving a group.
  * Added new resources:
    * `akamai_iam_cidr_block` - manages CIDR block assigned to the allowlist.
    * `akamai_iam_ip_allowlist` - enables or disable your account's allowlist.
  * Added new attributes to the `resource_akamai_iam_user` resource.
     * `user_notifications` to support user notifications.
     * `enable_mfa` to support authentication of type "MFA".
     * `password` to allow users to set a password when creating and updating a user.
  * Made the `enable_tfa` attribute optional in the `resource_akamai_iam_user` resource. 
  * Added the `asset_id` schema field (an IAM identifier of a property or include) to:
    * The `akamai_property` resource and data source,
    * The `akamai_property_include` resource and data source.
  * Improved date handling to use `time.Time` instead of `string`.

* PAPI
  * Added a new optional param to the import id of the `akamai_edge_hostname` resource.
    It allows to specify the product ID of the imported hostname and save it in the state.

#### BUG FIXES:

* PAPI
  * Added support for status code `429 Too Many Requests` containing `X-RateLimit-Next` header.
    When `X-RateLimit-Next` is present, the wait time before retry is calculated as the time
    difference between this header and the `Date` header.
  * Fixed an issue with the `akamai_property_activation` resource where updating it with an active or previously active property version for a configuration without a state file didn’t trigger a new property activation.

#### DEPRECATIONS:

* PAPI
  * Deprecated fields `product_id` and `rule_format` from `akamai_properties` datasource. Please use `akamai_property` to fetch this data. 

## 6.4.0 (Sep 04, 2024)

#### FEATURES/ENHANCEMENTS:

* Global
  * Updated SDK v2 and framework libraries as a result of updating `terraform-plugin-testing`.

* Appsec
  * Added the `request_body_inspection_limit_override` field to the `akamai_appsec_advanced_settings_request_body` resource.

* CPS
  * Added `acknowledge_post_verification_warnings` to the `akamai_cps_dv_validation` resource to allow for acknowledgement of post-verification warnings

* PAPI
  * Added support for new rule format `v2024-08-13`

#### BUG FIXES:

* Appsec
  * Fixed import of `akamai_appsec_match_target` for newly created security configuration or any security configuration not synced in the terraform state ([I#546](https://github.com/akamai/terraform-provider-akamai/issues/546))
  * Fixed issue where activation was not triggered after network list change in `resource_akamai_networklist_activations` resource ([I#518](https://github.com/akamai/terraform-provider-akamai/issues/518))
  * Fixed `akamai_appsec_configuration` data source to return a single security configuration in the output_text instead of the entire list of security configurations

* Cloudlets
  * Corrected format of the retry time when logging in `akamai_cloudlets_application_load_balancer_activation` and `akamai_cloudlets_policy_activation` resources

* PAPI
  * Fixed issue with provider producing an inconsistent final plan with Cloudlet policy ([I#567](https://github.com/akamai/terraform-provider-akamai/issues/567)).
    It happened in cases when content of the rule depends on some other resource

## 6.3.0 (July 16, 2024)

#### FEATURES/ENHANCEMENTS:

* Migrated Go version to `1.21.12` for builds.
* Appsec
  * Added field `host_names` to the `akamai_appsec_configuration` data source

* BOTMAN
  * Added new resource:
    * `akamai_botman_content_protection_javascript_injection_rule` - read and update
    * `akamai_botman_content_protection_rule` - read and update
    * `akamai_botman_content_protection_rule_sequence` - read and update
  * Added new data source:
    * `akamai_botman_content_protection_javascript_injection_rule` - read
    * `akamai_botman_content_protection_rule` - read
    * `akamai_botman_content_protection_rule_sequence` - read

* Client Lists
  * Extended list of fields for which `akamai_clientlist_activation` diff is suppressed with `notification_recipients` and `siebel_ticket_id`. Diff suppressed when activation is not required.

* Cloud Access
  * Added datasource:
    * `akamai_cloudaccess_key` - read details for key by name
    * `akamai_cloudaccess_keys` - read list of access key for current user account
    * `akamai_cloudaccess_key_versions` - read details for key versions by key name
    * `akamai_cloudaccess_key_properties` - read list of active properties for given access key
  * Added resource:
    * `akamai_cloudaccess_key` - create, read, update, delete, import

* DNS
  * Added data source:
    * `akamai_zone_dnssec_status` - reads the DNSSEC status of a single zone in Edge DNS ([I#509](https://github.com/akamai/terraform-provider-akamai/issues/509))

* GTM
  * Added more details for `gtm_property` resource in case of error being returned from the API

* PAPI
  * Added support for new rule format `v2024-05-31`
  * Added new optional field `ttl` to `akamai_edge_hostname` resource. 
    When it is used, creation or update takes longer as resource has to synchronize its state with HAPI.

#### BUG FIXES:

* Appsec
  * A new config version will be created if the latest config version is active in either Staging or Production, and protected and/or evaluated hostnames are updated using `akamai_appsec_wap_selected_hostnames` ([#I540](https://github.com/akamai/terraform-provider-akamai/issues/540))
  * Fixed issue where terraform provider plugin crashes due to empty string input for list `geo_network_lists`, `ip_network_lists`, `exception_ip_network_lists` and `asn_network_lists` in `akamai_appsec_ip_geo` resource

* DNS
  * Improved validation of IPv6 addresses in `akamai_dns_record` resource for records of type `AAAA` ([I#550](https://github.com/akamai/terraform-provider-akamai/issues/550))
  * Fixed issue in `akamai_dns_record` resource that could cause incorrect targets planned to be modified or reordering targets send to server for `TXT` record type ([I#499](https://github.com/akamai/terraform-provider-akamai/issues/499), [I#541](https://github.com/akamai/terraform-provider-akamai/issues/541), [I#559](https://github.com/akamai/terraform-provider-akamai/issues/559))
  * Fixed issue in `akamai_dns_recordset` datasource that for `TXT` record type, returned targets were needlessly reordered ([I#559](https://github.com/akamai/terraform-provider-akamai/issues/559))

* PAPI
  * Removed caching from `akamai_contracts` data source
  * Fixed issue in `akamai_edge_hostname` resource when update is performed straight after create
  * Fixed issue in data_akamai_property_rules_template that having root template in the same directory as .terraform dir would cause error.
    Now, datasource will not search for templates inside .terraform directory ([I#557](https://github.com/akamai/terraform-provider-akamai/issues/557))
  * Fixed an issue that caused the `compliance_record` in imported `akamai_property_activation` and `akamai_property_include_activation` to be empty and could not be updated.
    * Added the ability to update `compliance_record` in `akamai_property_activation` and `akamai_property_include_activation` via terraform apply (the update will not trigger new activation if version/network/property was not changed)
  * Fixed issue that having `akamai_property` and `akamai_property_activation` (or `akamai_property_include` and `akamai_property_include_activation`) resources linked using `staging_version` or `production_version`
    and modifying rules and note could sometimes result in `Provider produced inconsistent final plan` error ([I#549](https://github.com/akamai/terraform-provider-akamai/issues/549)).

## 6.2.0 (May 28, 2024)

#### FEATURES/ENHANCEMENTS:

* Global
  * Added validation for retryable logic values.
    * `retry_max` or `AKAMAI_RETRY_MAX` - Cannot be higher than 50
    * `retry_wait_min` or `AKAMAI_RETRY_WAIT_MIN` - Cannot be longer than 24h
    * `retry_wait_max` or `AKAMAI_RETRY_WAIT_MAX` - Cannot be longer than 24h
  * Migrated Terraform to version 1.7.5
  * Updated SDKv2 and framework libraries

* Appsec
  * Suppressed rate policy diff when `counterType` field absence is the only change for `akamai_appsec_rate_policy` resource
  * Suppressed activations diff when `notification_emails` field is the only change for `akamai_appsec_activations` resource

* BOTMAN
  * Added resource:
    * `akamai_botman_custom_bot_category_item_sequence` - read and update

* Cloudlets
  * Added import for `akamai_cloudlets_application_load_balancer_activation` resource

* GTM
  * Added data sources:
    * `akamai_gtm_geomap` - reads information for a specific GTM Geographic map
    * `akamai_gtm_geomaps` - reads information for GTM Geographic maps under a given domain

* IAM
 * Fixed issue of generating an incorrect large difference in `granted_roles` update ([I#525](https://github.com/akamai/terraform-provider-akamai/issues/525))

* Network Lists
  * Suppressed activations diff when `notification_emails` field is the only change for `akamai_networklist_activations` resource

* PAPI
  * Added retry logic for `akamai_property_include_activation`
  * Added import of the `certificate` for `akamai_edge_hostname` resource ([I#338](https://github.com/akamai/terraform-provider-akamai/issues/338))
  * NOTE: Certificate modification is not allowed.

#### BUG FIXES:

* Appsec
  * Resolved a drift issue with the `akamai_appsec_advanced_settings_attack_payload_logging` resource
  * Fixed an issue where resource `akamai_appsec_activations` continues in a loop after API throws an error. ([#I528](https://github.com/akamai/terraform-provider-akamai/issues/528))

* CPS
  * Fixed issue where modifications to SAN list in `akamai_cps_third_party_enrollment` of the `akamai_cps_upload_certificate` resource results in to inconsistency terraform plan error.

* DNS
  * Fixed issue in `akamai_dns_record` that modifying `priority` and/or `priority_increment` for `MX` record type was causing an error
* GTM
  * Fixed issue with order of `liveness_test` in `akamai_gtm_property` ([PR#404](https://github.com/akamai/terraform-provider-akamai/pull/404))

#### DEPRECATIONS:

* CPS
  * Deprecated field `unacknowledged_warnings` of `akamai_cps_upload_certificate` resource.

## 6.1.0 (Apr 23, 2024)

#### FEATURES/ENHANCEMENTS:

* DNS
  * Added second mode to `akamai_dns_record` resource where it is possible to provide individual values for priority, weight and port to **every** `SRV` target.
    In such case it is not allowed to provide values for resource level fields `priority`, `weight` and `port`.
    It is not allowed to mix targets with and without those fields.
    ([I#370](https://github.com/akamai/terraform-provider-akamai/issues/370))

* Image and Video Manager
  * Added support for `SmartCrop` transformation in `akamai_imaging_policy_image` datasource

#### BUG FIXES:

* CPS
  * Fixed issue with terraform producing inconsistent final plan for `akamai_cps_upload_certificate` resource on SAN list modification in `akamai_cps_third_party_enrollment` resource.

## 6.0.0 (Mar 26, 2024)

#### BREAKING CHANGES:

* General
  * Migrated to terraform protocol version 6, hence minimal required terraform version is 1.0

* PAPI
    * Added validation to raise an error if the creation of the `akamai_edge_hostname` resource is attempted with an existing edge hostname.
    * Added validation to raise an error during the update of `akamai_edge_hostname` resource for the immutable fields: 'product_id' and 'certificate'.

#### FEATURES/ENHANCEMENTS:

* Global
  * Requests limit value is configurable via field `request_limit` or environment variable `AKAMAI_REQUEST_LIMIT`
  * Added retryable logic for all GET requests to the API.
    This behavior can be disabled using `retry_disabled` field from `akamai` provider configuration or via environment variable `AKAMAI_RETRY_DISABLED`.
    It can be fine-tuned using following fields or environment variables:
    * `retry_max` or `AKAMAI_RETRY_MAX` - The maximum number retires of API requests, default is 10
    * `retry_wait_min` or `AKAMAI_RETRY_WAIT_MIN` - The minimum wait time in seconds between API requests retries, default is 1 sec
    * `retry_wait_max` or `AKAMAI_RETRY_WAIT_MAX` - The maximum wait time in minutes between API requests retries, default is 30 sec
  * Migrated to go 1.21
  * Bumped various dependencies

* Appsec
  * Added resource:
    * `akamai_appsec_penalty_box_conditions` - read and update
    * `akamai_appsec_eval_penalty_box_conditions` - read and update
  * Added new data source:
    * `akamai_appsec_penalty_box_conditions` - read
    * `akamai_appsec_eval_penalty_box_conditions` - read

* CPS
  * Added fields: `org_id`, `assigned_slots`, `staging_slots` and `production_slots` to `data_akamai_cps_enrollment` and `data_akamai_cps_enrollments` data sources

* Edgeworkers
  * Improved error handling in `akamai_edgeworkers_activation` and `resource_akamai_edgeworker` resources
  * Improved error handling in `akamai_edgeworker_activation` datasource

* GTM
  * Added fields:
    * `precedence` inside `traffic_target` in `akamai_gtm_property` resource and `akamai_gtm_domain` data source
    * `sign_and_serve` and `sign_and_serve_algorithm` in `akamai_gtm_domain` data source and resource
    * `http_method`, `http_request_body`, `alternate_ca_certificates` and `pre_2023_security_posture` inside `liveness_test` in `akamai_gtm_property` resource and `akamai_gtm_domain` data source
  * Added support for `ranked-failover` properties in `akamai_gtm_property` resource
  * Enhanced error handling in `akamai_gtm_asmap`, `akamai_gtm_domain`, `akamai_gtm_geomap`, `akamai_gtm_property` and `akamai_gtm_resource` resources

* IMAGING
  * In the event of an API failure during a policy update, reverting to the previous state ([I#491](https://github.com/akamai/terraform-provider-akamai/issues/491), [I#493](https://github.com/akamai/terraform-provider-akamai/issues/493))
  * When performing the read operation, if `activate_on_production` is true, fetch the policy state from the production network; otherwise, obtain it from the staging environment.

* PAPI
  * Added attributes to akamai_property datasource:
    * `contract_id`, `group_id`, `latest_version`, `note`, `production_version`, `product_id`, `property_id`, `rule_format`, `staging_version`
  * `data_akamai_property_rules_builder` is now supporting `v2024-02-12` [rule format](https://techdocs.akamai.com/terraform/reference/rule-format-changes#v2024-02-12)

#### BUG FIXES:

* Appsec
  * Fixed ukraine_geo_control_action drift issue ([I#484](https://github.com/akamai/terraform-provider-akamai/issues/484))

* Cloudlets
    * Allowed empty value for match rules `json` attribute for data sources:
        * `akamai_cloudlets_api_prioritization_match_rule`
        * `akamai_cloudlets_application_load_balancer_match_rule`
        * `akamai_cloudlets_audience_segmentation_match_rule`
        * `akamai_cloudlets_edge_redirector_match_rule`
        * `akamai_cloudlets_forward_rewrite_match_rule`
        * `akamai_cloudlets_phased_release_match_rule`
        * `akamai_cloudlets_request_control_match_rule`
        * `akamai_cloudlets_visitor_prioritization_match_rule`

* CPS
  * Changed below fields from required to optional in `akamai_cps_dv_enrollment` and `akamai_cps_third_party_enrollment`
    for `admin_contact` and `tech_contact` attributes:
    * `organization`
    * `address_line_one`
    * `city`
    * `region`
    * `postal_code`
    * `country_code`

* PAPI
  * Fixed case when `origin_certs_to_honor` field from `origin` behavior mandates presence of empty `custom_certificate_authorities` and/or `custom_certificates` options inside `origin` behavior for `akamai_property_rules_builder` datasource ([I#515](https://github.com/akamai/terraform-provider-akamai/issues/515))

#### DEPRECATIONS

* Appsec
  * `akamai_appsec_selected_hostnames` data source and resource are deprecated with a scheduled end-of-life in v7.0.0 of our provider. Use the `akamai_appsec_configuration` instead.

## 5.6.0 (Feb 19, 2024)

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Added retries in `akamai_appsec_activations` and `akamai_networklist_activations` resources ([I#471](https://github.com/akamai/terraform-provider-akamai/issues/471))
  * Added reactivation support for `akamai_appsec_activations` if the config was deactivated manually ([I#441](https://github.com/akamai/terraform-provider-akamai/issues/441) and [I#442](https://github.com/akamai/terraform-provider-akamai/issues/442))

* Cloudlets
  * Added support for Shared Cloudlets Policies. To use it, provide `is_shared` field in `akamai_cloudlets_policy` resource as `true`. ([I#276](https://github.com/akamai/terraform-provider-akamai/issues/276))
  * Added validation to prevent changing immutable `cloudlet_code` field in `akamai_cloudlets_policy` resource
  * Added support for importing policies without any version
  * Added new data source:
    * `akamai_cloudlets_policy_activation` - read
    * `akamai_cloudlets_shared_policy` - read
  * Changes for `akamai_cloudlets_policy_activation` resource
    * Added support for shared (V3) policies
    * Added import for `akamai_cloudlets_policy_activation`
    * Field `associated_properties` was changed to optional but is still required for non-shared policies
    * Added `is_shared` computed field to indicate if processing policy is shared

* DNS
  * Enhanced handling of `akamai_dns_zone` resource when no `group` is provided:
    * When there is only one group present, the processing should continue with a descriptive warning
    * When there are more than one group present, the processing will fail with descriptive error asking to provide group in the configuration

* Edgeworkers
  * Added `note` attribute to `resource_akamai_edgeworkers_activation` resource

* GTM
  * Added data sources:
    * `akamai_gtm_asmap` - reads information for a specific GTM asmap
    * `akamai_gtm_resources` - reads information for a specific GTM resources under given domain
    * `akamai_gtm_resource` - reads information for a specific GTM resource
    * `akamai_gtm_domain` - reads information for a specific GTM domain
    * `akamai_gtm_domains` - reads list of GTM domains under a given contract
    * `akamai_gtm_cidrmap` - reads information for a specific GTM cidrmap

* IVM
  * Extended `akamai_imaging_policy_image` with new fields:
    * `serve_stale_duration` available under `policy`
    * `allow_pristine_on_downsize` and `prefer_modern_formats` available under `policy.output`

* PAPI
  * Added new resource:
    * `akamai_property_bootstrap` - create, read, update and delete property without specifying rules or edgehostnames. To be used with `akamai_property` resource and its new field `property_id` ([I#466](https://github.com/akamai/terraform-provider-akamai/issues/466))
  * Added `version_notes`, `rule_warnings` and `property_id` attributes to `akamai_property` resource ([I#494](https://github.com/akamai/terraform-provider-akamai/issues/494))
  * Added support for new rule format v2024-01-09 in `data_akamai_property_rules_builder`
  * Improved errors for `akamai_contract` and `akamai_group` datasources when there are multiple groups or contracts
  * Added `name` validation for `akamai_property_include` resource

* Updated various dependencies

#### BUG FIXES:

* Appsec
  * Fixed provider plugin crash in `appsec_attack_group` and `appsec_eval_group` after executing terraform plan ([I#480](https://github.com/akamai/terraform-provider-akamai/issues/480))
  * Fixed drift for struct and list reordering in `akamai_appsec_match_target`

* Cloudlets
  * Fixed handling of version drift for cloudlets policies ([I#478](https://github.com/akamai/terraform-provider-akamai/issues/478))

* CPS
  * Changed `organizational_unit` inside `csr` attribute for `akamai_cps_third_party_enrollment` and `akamai_cps_dv_enrollment`
    resources from required to optional. ([PR#513](https://github.com/akamai/terraform-provider-akamai/pull/513))
  * Changed `state` inside `csr` attribute for `akamai_cps_third_party_enrollment` and `akamai_cps_dv_enrollment` resources from required to optional.

* GTM
  * Fixed 'Inconsistent Final Plan' error for `akamai_gtm_property` resource
  * The diff when reordering `traffic_target` in `akamai_gtm_property` resource at the same time as changing any attribute value inside `traffic_target` will be extensive
  * Added `ForceNew` to the `name` attribute for `akamai_gtm_property` resource as it is not possible to rename it using API

## 5.5.0 (Dec 07, 2023)

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Updated resource:
    * `akamai_appsec_ip_geo` - added `asn_network_lists` attribute to support blocking by ASN client lists
  * Updated data source:
    * `akamai_appsec_ip_geo` - added `asn_network_lists` attribute to list ASN client lists

* BOTMAN
  * Added resource:
    * `akamai_botman_custom_code` - read and update
  * Added data source:
    * `akamai_botman_custom_code` - read
  * Cached api calls for `akamai_botman_akamai_bot_category`, `akamai_botman_akamai_defined_bot` and `akamai_botman_bot_detection` data sources to improve performance.

* Cloudlets
  * Added `origin_description` field to `akamai_cloudlets_application_load_balancer` resource

* PAPI
  * Behavior `restrict_object_caching` is public ([I#314](https://github.com/akamai/terraform-provider-akamai/issues/314) and [#277](https://github.com/akamai/terraform-provider-akamai/issues/277))
  * Added version support for `akamai_property_hostnames` data source ([I#413](https://github.com/akamai/terraform-provider-akamai/issues/413))
  * `data_akamai_property_rules_builder` is now supporting `v2023-10-30` rule format
  * Improved error handling and added retries in `akamai_property_activation` resource
  * Relaxed validation used for includes used in `akamai_property_rules_template`. Files cannot be empty but do not necessary have to be valid json files.

#### BUG FIXES:

* DNS
  * Fixed handling of txt records which are longer than 255 bytes ([I#430](https://github.com/akamai/terraform-provider-akamai/issues/430))

* Image and Video Manager
  * Added suppression when providing `ctr_` prefix in `akamai_imaging_policy_set` ([I#491](https://github.com/akamai/terraform-provider-akamai/issues/491))

## 5.4.0 (Oct 31, 2023)

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Suppressed trigger of new activation for `note` field change in `akamai_networklist_activations` and `akamai_appsec_activations` resources.

* Client Lists
  * Added support for state import for `akamai_clientlist_list` and `akamai_clientlist_activation` resources

* Cloudlets
  * Added `matches_alway` field to `akamai_cloudlets_edge_redirector_match_rule` data source
  * Added configurable timeout for following resources as `timeouts.default` field
    * `akamai_cloudlets_application_load_balancer_activation`
    * `akamai_cloudlets_policy_activation`
    * `akamai_cloudlets_policy`

* CPS
  * Added configurable timeout for following resources as `timeouts.default` field ([I#440](https://github.com/akamai/terraform-provider-akamai/issues/440))
    * `akamai_cps_dv_enrollment`
    * `akamai_cps_dv_validation`
    * `akamai_cps_third_party_enrollment`
    * `akamai_cps_upload_certificate`

* Edgeworkers
  * Added configurable timeout for following resources as `timeouts.default` field
    * `akamai_edgekv_group_items`
    * `akamai_edgeworker`
  * Added configurable timeout for `akamai_edgeworkers_activation` resource as `timeouts.default` and `timeouts.delete` fields

* IAM
  * Phone number is no longer required for IAM user in `akamai_iam_user` resource.

* PAPI
  * Added configurable timeout for following resources as `timeouts.default` field ([I#440](https://github.com/akamai/terraform-provider-akamai/issues/440))
    * `akamai_property_activation`
    * `akamai_property_include_activation`
    * `akamai_edge_hostname`
  * Added configurable timeout for `akamai_cp_code` resource as `timeouts.update` field
  * Changed `version` field in `akamai_property_activation` data source to optional. Now when `version` is not provided,
    datasource automatically finds the active one for given network.
  * Allowed empty values for some fields
    in `akamai_property_builder` ([I#481](https://github.com/akamai/terraform-provider-akamai/issues/481))
  * Added support for new rule format `v2023-09-20`

#### BUG FIXES:

* GTM
  * Fixed problem with wrong datacenters updated in `akamai_gtm_property`.

* IAM
  * Fixed Terraform proposing modifications to user settings when using international phone numbers in `akamai_iam_user`
    resource.
    * NOTE:
      * For international phone numbers there might be a diff during plan. Please apply suggested change to store the
        correct number.
      * Invalid phone numbers will block the plan.

* PAPI
  * Made `status_update_email` attribute optional in `akamai_edge_hostname` resource

## 5.3.0 (Sep 26, 2023)

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Added `sync_point` value in `akamai_networklist_network_lists` data source

* CPS
  * Added `pending_changes` computed field to `akamai_cps_enrollment` data source ([#PR468](https://github.com/akamai/terraform-provider-akamai/pull/468))

* Cloud Wrapper
  * Added support for `comments` argument modification in `akamai_cloudwrapper_configuration` resource

#### BUG FIXES:

* Appsec
  * Fixed `akamai_networklist_network_list` import resulting in null `contract_id` and `group_id`

* PAPI
  * Added errors to `data_property_akamai_contract` and `data_property_akamai_group` data sources, when fetching groups returns multiple inconclusive results
  * Fixed drift issue in `akamai_edge_hostname` resource [(#457)](https://github.com/akamai/terraform-provider-akamai/issues/457)
  * Added missing fields to `akamai_property_builder` for `origin` and `siteShield` behaviors ([#465](https://github.com/akamai/terraform-provider-akamai/issues/465))
  * Improved `akamai_property_rules_builder` empty list transformation ([#438](https://github.com/akamai/terraform-provider-akamai/issues/438))

* GTM
  * Added better drift handling in `akamai_gtm_property` - when property is removed without terraform knowledge, resource doesn't just error on refresh but suggests recreation

## 5.2.0 (Aug 29, 2023)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Cloud Wrapper
  * Added resources:
    * `akamai_cloudwrapper_activation` - activate cloud wrapper configuration, import cloud wrapper configuration activation
    * `akamai_cloudwrapper_configuration` - create, read and update cloud wrapper configuration
  * Added data sources:
    * `akamai_cloudwrapper_capacities` - reads capacities available for the provided contract IDs
    * `akamai_cloudwrapper_configuration` - reads configuration associated with config ID
    * `akamai_cloudwrapper_configurations` - reads all the configurations
    * `akamai_cloudwrapper_location` - reads location for given location name and traffic type
    * `akamai_cloudwrapper_locations` - reads all locations
    * `akamai_cloudwrapper_properties` - reads properties associated with contract IDs with Cloud Wrapper entitlement

* [IMPORTANT] Client Lists
  * Added resources:
    * `akamai_clientlist_list` - create, update and delete Client Lists
    * `akamai_clientlist_activation` - activate a client list
  * Added data source:
    * `akamai_clientlist_lists` - reads Client Lists
      *  Support filter by `name` and/or `types`

* BOTMAN
  * Added resource:
    * `akamai_botman_custom_client_sequence` - read and update custom client sequence
  * Added data source:
    * `akamai_botman_custom_client_sequence` - reads custom client sequence

* PAPI
  * `logStreamName` field from `datastream` behavior has changed from string to array of strings for rule format `v2023-05-30`

## 5.1.0 (Aug 01, 2023)

#### BUG FIXES:

* PAPI
  * Dropped too strict early snippet validation ([#436](https://github.com/akamai/terraform-provider-akamai/issues/436))
  * Fixed issue that `akamai_property` or `akamai_property_include` would sometimes show strange `null -> null` diff
    in `rules` (or dropping `null` in newer Terraform versions) even if no update actually is needed. If there is
    anything else changing in the rule tree, the `null -> null` will be also visible in the diff. That may be fixed in
    later time.
  * Fixed issue that `akamai_property_rules_builder` data source did not support PM variables for fields with validation
    based on regular expressions

#### FEATURES/ENHANCEMENTS:

* Appsec
  * Added resource:
    * `akamai_appsec_security_policy_default_protections`

* BOTMAN
  * Added resource:
    * `akamai_botman_challenge_injection_rules` - read and update
  * Added data sources:
    * `akamai_botman_challenge_injection_rules` - read

* PAPI
  * Added verification to ensure that `akamai_property_rules_builder` data source
    has consistent frozen rule format between parent and it's child.
    Additionally `akamai_property_rules_builder.json` is returning artificial field `_ruleFormat_`.
  * Suppressed trigger of new activation for `note` field change in `akamai_property_activation` and `akamai_property_include_activation` resources.

#### DEPRECATIONS

* Appsec
  * deprecated following resources; use `akamai_appsec_security_policy_default_protections` resource instead:
    * `akamai_appsec_api_constraint_protection`
    * `akamai_appsec_ip_geo_protection`
    * `akamai_appsec_malware_protection`
    * `akamai_appsec_rate_protection`
    * `akamai_appsec_reputation_protection`
    * `akamai_appsec_slowpost_protection`

* BOTMAN
  * deprecated `akamai_botman_challenge_interception_rules` data source and resource; use `akamai_botman_challenge_injection_rules` instead.

## 5.0.1 (Jul 12, 2023)

#### BUG FIXES:

* Reinstated support for configuring provider with environmental variables ([#407](https://github.com/akamai/terraform-provider-akamai/issues/407), [#444](https://github.com/akamai/terraform-provider-akamai/issues/444))
* Fixed `signature does not match` error when using `config` block for authentication ([#444](https://github.com/akamai/terraform-provider-akamai/issues/444), [#446](https://github.com/akamai/terraform-provider-akamai/issues/446))

## 5.0.0 (Jul 5, 2023)

#### BREAKING CHANGES:

* DataStream
  * Changed the following data sources in DataStream 2 V2 API:
    * `akamai_datastream_activation_history` - changed schema and corresponding implementations.
    * `akamai_datastream_dataset_fields` - changed parameter, schema and corresponding implementations.
    * `akamai_datastreams` - changed parameter, schema and corresponding implementations.
  * Changed the following resources in DataStream 2 V2 API:
    * `akamai_datastreams` - changed in schema payload, response attributes and corresponding implementations.
  * Updated attribute names in `datastream.connectors`.
  * Updated methods in `datastream.stream` for the above changes.

* PAPI
  * Changed default value of `auto_acknowledge_rule_warnings` to `false` in `akamai_property_activation` resource

* Removed undocumented support for configuring provider with environment variables (`AKAMAI_ACCESS_TOKEN`, `AKAMAI_CLIENT_TOKEN`, `AKAMAI_HOST`, `AKAMAI_CLIENT_SECRET`, `AKAMAI_MAX_BODY`, and their `AKAMAI_{section}_xxx` equivalents).
  As an alternative users should now use provider's [config](https://techdocs.akamai.com/terraform/docs/gs-authentication#use-inline-credentials) block with [TF_VAR_](https://developer.hashicorp.com/terraform/language/values/variables#environment-variables) envs when wanting to provide configuration through enviroment variables.

##### Removed deprecated schema fields

* Appsec
  * `notes` and `activate` fields in `akamai_appsec_activations` resource
  * `appsec_section` and `appsec` fields in provider schema

* CPS
  * `enable_multi_stacked_certificates` field in `akamai_cps_dv_enrollment` resource

* DNS
  * `dns_section` and `dns` fields in provider schema

* GTM
  * `gtm_section` and `gtm` fields in provider schema

* IAM
  * `is_locked` field in `akamai_iam_user` resource

* Network Lists
  * `activate` field in `akamai_networklist_activations` resource
  * `networklist_section` and `network` fields in provider schema

* PAPI
  * `contract` and `group` fields in `akamai_cp_code` data source
  * `group` field in `akamai_contract` data source
  * `name` and `contract` fields in `akamai_group` data source
  * `contract`, `group` and `product` fields in `akamai_cp_code` resource
  * `contract`, `group` and `product` fields in `akamai_edge_hostname` resource
  * `property` and `rule_warnings` fields in `akamai_property_activation` resource
  * `contract`, `group` and `product` fields in `akamai_property` resource
  * `papi_section`, `property_section` and `property` fields in provider schema

##### Removed deprecated resource

* PAPI
  * `akamai_property_variables`

#### FEATURES/ENHANCEMENTS:

* Provider tested and now supports Terraform 1.4.6
* Migrated `akamai_property_include` data source from SDKv2 to Framework.

* PAPI
  * Added import to `akamai_property_activation` resource
  * Extended `akamai_property_rules_builder` data source: added support for rules frozen format `v2023-01-05` and `v2023-05-30`

* Appsec
  * Updated Geo control to include Action for Ukraine.
  * Added `akamai_appsec_advanced_settings_pii_learning` data source and resource for managing the PII learning advanced setting.

#### DEPRECATIONS

* Deprecated `active` field in `akamai_dns_record` resource

#### BUG FIXES:

* CPS
  * Fixed bug in `akamai_cps_dv_enrollment` resource when MTLS settings are provided ([#339](https://github.com/akamai/terraform-provider-akamai/issues/339))
  * Fixed `sans` field causing perpetual in-place update in `akamai_cps_third_party_enrollment` ([#415](https://github.com/akamai/terraform-provider-akamai/issues/415))

* GTM
  * Made `test_object` inside `liveness_test` required only for `test_object_protocol` values: `HTTP`, `HTTPS` or `FTP` ([I#408](https://github.com/akamai/terraform-provider-akamai/issues/408))

* Cloudlets
  * Added wait for propagation of policy activation deletions, before removing the policy in `akamai_cloudlets_policy` ([I#420](https://github.com/akamai/terraform-provider-akamai/issues/420))

* PAPI
  * Removed hostname validation on `akamai_property` resource ([I#422](https://github.com/akamai/terraform-provider-akamai/issues/422))


## 4.1.0 (Jun 1, 2023)

#### FEATURES/ENHANCEMENTS:

* GTM
  * New data sources:
    * `akamai_gtm_datacenter` - get datacenter information
    * `akamai_gtm_datacenters` - get datacenters information

## 4.0.0 (May 30, 2023)

#### BREAKING CHANGES:

* Appsec
  * Update malware policy `ContentTypes` to include `EncodedContentAttributes`.
  * Malware policy's `ContentTypes` is reported as part of an individual policy but is no longer included in the bulk report of all policies.

* PAPI
  * Remove `cpc_` prefix in `akamai_cp_code` resource and data source IDs

#### FEATURES/ENHANCEMENTS:

* Migrate to Terraform 1.3.7 version

* Akamai
  * Reword returned error when reading edgerc configuration encounters problems ([I#411](https://github.com/akamai/terraform-provider-akamai/issues/411))

* EdgeWorkers
  * Deactivate EdgeWorker versions upon EdgeWorker deletion([I#331](https://github.com/akamai/terraform-provider-akamai/issues/331))

* PAPI
  * Remove enforce `property-snippets` directory check ([I#378](https://github.com/akamai/terraform-provider-akamai/issues/378))
  * Improved variable evaluation logic in `akamai_property_rules_template` data source ([I#324](https://github.com/akamai/terraform-provider-akamai/issues/324), [I#385](https://github.com/akamai/terraform-provider-akamai/issues/385), [I#386](https://github.com/akamai/terraform-provider-akamai/issues/386))
    * Include path can now be provided using data source `variables`
    * `variables` can now reference each other and be used to build other `variables` e.g. `${env.abc} = "${env.prefix} cba"`
    * Variables existence is now verified early across all snippets inside the snippets directory - if variable is used in a snippet which is not included in final template and the variable is not defined, the processing will fail (previously variables were verified only when the snippet was loaded into final result)
  * (Internal usage only) Improved `compliance_record` attribute's syntax for `akamai_property_activation` and `akamai_property_include_activation`

#### BUG FIXES:

* Appsec
  * Fixed issue that in some cases allowed `terraform plan` to create a new config version as a side-effect of reading the current config.

* DNS
  * Fixed TXT record characters escaping issue in akamai_dns_record resource ([I#137](https://github.com/akamai/terraform-provider-akamai/issues/137))
  * Fixed issue when `target` in `akamai_dns_record` resource was not known during plan, the plan failed ([I#410](https://github.com/akamai/terraform-provider-akamai/issues/410))

* Cloudlets
  * Fixed bug related to regex validation for handling property delay in `akamai_cloudlets_policy_activation`
  * Fixed sporadic issue with `akamai_cloudlets_policy_activation` due to network delay

* PAPI
  * Fixed reading float values in `akamai_property_rules_builder`
  * Add validation for hostnames `cname_from` field in `akamai_property` resource
  * Assign only active property activation version in `akamai_property_activation` resource on read

## 3.6.0 (April 27, 2023)

#### FEATURES/ENHANCEMENTS:

* EdgeKV
  * Added resource:
    * `akamai_edgekv_group_items` - create, read, update, delete and import
  * Added data sources:
    * `akamai_edgekv_group_items` - reads group items associated with namespace and network
    * `akamai_edgekv_groups` - reads groups associated with namespace and network
  * Deprecated field `initial_data` under `akamai_edgekv` resource

#### BUG FIXES:

* Cloudlets
  * In some cases `akamai_cloudlets_application_load_balancer_activation` or `akamai_cloudlets_policy_activation` were not activating due to verification delay with property resource.

* CPS
  * Get CSR from long history ([I#403](https://github.com/akamai/terraform-provider-akamai/issues/403))

* GTM
  * Deprecated field `name` of `traffic_target` under `akamai_gtm_property` resource ([I#374](https://github.com/akamai/terraform-provider-akamai/issues/374))

* Image and Video Manager:
  * Fixed diff in `akamai_imaging_policy_image` resource for image policy attributes ([I#383](https://github.com/akamai/terraform-provider-akamai/issues/383)):
    * `Breakpoints.Widths`
    * `Hosts`
    * `Output.AllowedFormats`
    * `Output.ForcedFormats`
    * `Variables`
  * Fixed diff in `akamai_imaging_policy_video` resource for video policy attributes:
    * `Breakpoints.Widths`
    * `Hosts`
    * `Variables`
  * Fixed diff seen in exported imaging policy set - removed default values

* PAPI
  * `is_secure` and `variable` fields can only be used with `default` rule in `akamai_property_rules_builder` data source
  * Delete `product_id` from import of `akamai_edge_hostname`

## 3.5.0 (March 30, 2023)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Advanced Options Settings - New settings added for Request Size Inspection Limit
    * Add data source `akamai_appsec_advanced_settings_request_body`
    * Add resource `akamai_appsec_advanced_settings_request_body`

* BOTMAN
  * Cache OpenAPI calls to improve performance

* Image and Video Manager
  * Add `forced_formats` and `allowed_formats` fields to `output` field

* PAPI
  * Add data source
    * `akamai_property_rules_builder` - create property rule trees directly from HCL (Beta).
  * Add `compliance_record` for `akamai_property_activation` resource

#### BUG FIXES:

* PAPI
  * Fix issue when `akamai_property` imported an older version and during update it didn't create a new version from it

* APPSEC
  * Fix issue updating rule action for ASE AUTO policy

## 3.4.0 (March 2, 2023)

#### FEATURES/ENHANCEMENTS:

* Various dependencies updated
* Updated Akamai Terraform Provider examples to be compliant with current Akamai Terraform Provider version and `TFLint`

* APPSEC
  * Advanced Options Settings - New settings added for Attack Payload Logging
    * Added data source
      [akamai_appsec_advanced_settings_attack_payload_logging](docs/data-sources/appsec_advanced_settings_attack_payload_logging.md)
    * Added resource
      [akamai_appsec_advanced_settings_attack_payload_logging](docs/resources/appsec_advanced_settings_attack_payload_logging.md)

#### BUG FIXES:

* APPSEC
  * Fix drift on `logFilename` element of `malware_policy`
  * Prevent changes to `rate_policy` field of existing `akamai_appsec_rate_policy_action` resource
  * Fix issue that disabled users from using all values allowed by the API in `akamai_appsec_rate_policy_action resource` resource
* PAPI
  * Fix issue when `akamai_property_include_activation` broke during creation, Terraform could not recover
  * Fixed issue that `property_rules_template` data source failed with multiple includes in array ([#387](https://github.com/akamai/terraform-provider-akamai/pull/387))

## 3.3.0 (February 2, 2023)

#### FEATURES/ENHANCEMENTS:

* Support for Go 1.18

* PAPI - Added data source for property activation
  * [akamai_property_activation](docs/data-sources/property_activation.md) - get activation by network

* CPS
  * Add `preferred_trust_chain` to `csr` attribute for `akamai_cps_dv_enrollment` resource

#### BUG FIXES:

* GTM
  * Fixed diff in resources:
    * `resource_akamai_gtm_asmap` for field `assignment.as_numbers`
    * `resource_akamai_gtm_cidrmap` for field `assignment.blocks`
    * `resource_akamai_gtm_geomap` for field `assignment.countries`
    * `resource_akamai_gtm_domain` for field `email_notification_list`
    * `resource_akamai_gtm_resource` for field `resource_instance.load_servers`

* CPS
  * Fixed terraform always showing diff for fields that use unicode characters ([#368](https://github.com/akamai/terraform-provider-akamai/issues/368))

## 3.2.1 (December 16, 2022)

#### BUG FIXES:

* PAPI
  * Fix `rule_format` in `akamai_property` to accept `latest`

## 3.2.0 (December 15 2022)

#### FEATURES/ENHANCEMENTS:

* PAPI - Add support for Property Includes
  * Added resources:
    * [akamai_property_include](docs/resources/property_include.md) - create, read, update, delete and import
    * [akamai_property_include_activation](docs/resources/property_include_activation.md) - create, read, update, delete and import
  * Added data sources:
    * [akamai_property_include_activation](docs/data-sources/property_include_activation.md) - get latest include activation by network
    * [akamai_property_include_parents](docs/data-sources/property_include_parents.md) - get property include parents information
    * [akamai_property_include_rules](docs/data-sources/property_include_rules.md) - get property include version rules information
    * [akamai_property_include](docs/data-sources/property_include.md) - get property include version information
    * [akamai_property_includes](docs/data-sources/property_includes.md) - list property includes information

* APPSEC
  * Add `json` attribute to `akamai_appsec_security_policy` data source to allow obtaining policy name given its ID.

#### BUG FIXES:

* APPSEC
  * Fixed bug that prevented `akamai_appsec_ip_geo` resource from sending correct network lists in `block` mode.
  * Fixed bug that prevented `akamai_appsec_configuration` data source from reporting error correctly when a nonexistent configuration is specified.

## 3.1.0 (December 1, 2022)

#### FEATURES/ENHANCEMENTS:

* CPS
  * New data sources:
    * [akamai_cps_csr](docs/data-sources/cps_csr.md) - returns latest Certificate Signing Request for given enrollment
    * [akamai_cps_deployments](docs/data-sources/cps_deployments.md) - returns deployed certificates for given enrollment
    * [akamai_cps_warnings](docs/data-sources/cps_warnings.md) - returns a map of all possible CPS warnings (ID to warning message). The IDs can be later used to approve warnings (auto_approve_warnings field)
  * Added resources allowing management of third-party enrollments:
    * [akamai_cps_third_party_enrollment](docs/resources/cps_third_party_enrollment.md) - create, read, update, delete and import third-party enrollments
    * [akamai_cps_upload_certificate](docs/resources/cps_upload_certificate.md) - create, read, update and delete
  * Resource cps_dv_enrollment
    * Deprecate `enable_multi_stacked_certificates` field. Now its value is always `false`.

#### BUG FIXES:

* PAPI
  * Fixed update of ip_behavior in `akamai_edge_hostname` resource ([#354](https://github.com/akamai/terraform-provider-akamai/issues/354))

## 3.0.0 (October 27, 2022)

#### BREAKING CHANGES:

* APPSEC
  * Require network list sync point for network list activation ([#326](https://github.com/akamai/terraform-provider-akamai/pull/326))

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Automatically activate network list when contents are modified
  * Increase timeout for security configuration activation to 90 minutes ([#348](https://github.com/akamai/terraform-provider-akamai/issues/348))

* Datastream
  * Added `akamai_datastreams` data source ([#327](https://github.com/akamai/terraform-provider-akamai/issues/327))
  * Added new features to `akamai_datastream` resource
    * new connectors: Elasticsearch, NewRelic and Loggly
    * Splunk and Custom HTTPS connectors were extended with ability to provide mTLS certificates configuration
    * SumoLogic, Splunk and Custom HTTPS connectors were extended with ability to specify custom HTTP headers

#### BUG FIXES:

* APPSEC
  * Fix incorrect payload sent by `akamai_appsec_ip_geo` resource in allow mode

* Datastream
  * Fixed problem with updating the configuration of the following connectors: Splunk, SumoLogic, Custom HTTPS, Datadog

* GTM
  * Fixed unreadable diff when single attribute is changed in traffic target
    > If upgrading to provider version 3.0.0 from an older version and using Akamai GTM Properties, you might see a switch of targets during the first apply. This is needed to get your terraform state in sync with our API, the following terraform plan/apply will not be affected

## 2.4.2 (October 4, 2022)

#### FEATURES/ENHANCEMENTS:

* IAM
  * Update docs for following resources and datasources as they are no longer in Beta
    * `akamai_iam_grantable_roles`
    * `akamai_iam_blocked_user_properties`
    * `akamai_iam_group`
    * `akamai_iam_role`

* Image and Video Manager
  * Update docs for the sub-provider as it's no longer in Beta

#### BUG FIXES:

* Botman
  * Fix page header for the Botman Getting Started Guide.

## 2.4.1 (September 29, 2022)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added Bot Management API Support
  * Added resources allowing management of:
    * `akamai_bot_category_action` - read, update and import
    * `bot_analytics_cookie` - read, update and import
    * `bot_category_exception` - read, update and import
    * `bot_detection_action` - read, update and import
    * `bot_management_settings` - read, update and import
    * `challenge_action` - create, read, update, delete and import
    * `challenge_interception_rules` - read, update and import
    * `client_side_security` - read, update and import
    * `conditional_action` - create, read, update, delete and import
    * `custom_bot_cateogry` - create, read, update, delete and import
    * `custom_bot_category_action` - read, update and import
    * `custom_bot_category_sequence` - read, update and import
    * `custom_client` - create, read, update, delete and import
    * `custom_defined_bot` - create, read, update, delete and import
    * `custom_deny_action` - create, read, update, delete and import
    * `javascript_injection` - read, update and import
    * `recategorized_akamai_defined_bot` - create, read, update, delete and import
    * `serve_alternate_action` - create, read, update, delete and import
    * `transactional_endpoint_protection` - read, update and import
    * `transactional_endpoint` - create, read, update, delete and import
  * Added data sources:
    * `akamai_bot_category` - list akamai bot categories
    * `akamai_bot_category_action` - list akamai bot category actions
    * `akamai_defined_bot` - list akamai defined bots
    * `bot_analytics_cookie` - get bot analytics cookie
    * `bot_analytics_cookie_values` - list bot analytics cookie values
    * `bot_category_exception` - list bot category exceptions
    * `bot_detection` - list bot detections
    * `bot_detection_action` - list bot detection actions
    * `bot_endpoint_coverage_report` - get bot endpoint coverage report
    * `bot_management_settings` - list bot management settings
    * `challenge_action` - list challenge actions
    * `challenge_interception_rules` - list challenge interception rules
    * `client_side_security` - get client side security
    * `conditional_action` - list conditional actions
    * `custom_bot_cateogry` - list custom bot categories
    * `custom_bot_category_action` - list custom bot category actions
    * `custom_bot_category_sequence` - get custom bot category sequence
    * `custom_client` - list custom clients
    * `custom_defined_bot` - list custom defined bots
    * `custom_deny_action` - list custom deny actions
    * `javascript_injection` - get javascript injection
    * `recategorized_akamai_defined_bot` - list recategorized akamai defined bots
    * `response_action` - list response actions
    * `serve_alternate_action` - list serve alternate actions
    * `transactional_endpoint` - list transactional endpoints
    * `transactional_endpoint_protection` - read, update and import

* APPSEC
  * New data sources:
    * `akamai_appsec_malware_content_types` - list available content types for malware protection
    * `akamai_appsec_malware_policies` - list malware policies
    * `akamai_appsec_malware_policy_actions` - list malware policy actions
  * New resources:
    * `akamai_appsec_malware_policy` - create, modify, or delete malware policies
    * `akamai_appsec_malware_policy_action` - create, modify, or delete the actions associated with a malware policy
    * `akamai_appsec_malware_policy_actions` - create, modify, or delete the actions associated with one or more policies within a given security policy
    * `akamai_appsec_malware_protection` - enable or disable malware protection for a security policy

* EdgeWorkers
  * New data sources ([#331](https://github.com/akamai/terraform-provider-akamai/issues/331)):
    * [akamai_edgeworker](docs/data-sources/edgeworker.md) - returns data for specific edgeworker, corresponding version and bundle information
    * [akamai_edgeworker_activation](docs/data-sources/edgeworkers_activation.md) - returns the latest activation in provided network
  * Resources:
    * `akamai_edgeworker_activation` - import

#### BUG FIXES:

* GTM
  * Fix diff for traffic_targets servers in `akamai_gtm_property` resource

## 2.3.0 (August 25, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add notification_emails to activations resource
  * Deprecate existing import functionality; use `cli-terraform export-appsec` instead.

* CPS
  * Extend `akamai_cps_dv_enrollment` with `allow_duplicate_common_name` field
  * New data sources:
    * [akamai_cps_enrollment](docs/data-sources/cps_enrollment.md) - returns data for specific enrollment
    * [akamai_cps_enrollments](docs/data-sources/cps_enrollments.md) - returns data for all of a specific contract's enrollments

#### BUG FIXES:

* Cloudlets
  * Add missing cloudlet codes in Cloudlets documentation ([#323](https://github.com/akamai/terraform-provider-akamai/issues/323))

* EdgeWorker
  * Fix EdgeWorker bundle hash calculation ([#321](https://github.com/akamai/terraform-provider-akamai/issues/321))

* GTM
  * Fix diff for traffic_targets in `akamai_gtm_property` resource
  * Fix `akamai_gtm_domain` shows diff after import
  * Fix `akamai_gtm_resource` shows diff after import
  * Fix terraform import of `akamai_gtm_asmap` does not import assignments

* PAPI
  * Fix error when using uppercase for `edge_hostname` ([#330](https://github.com/akamai/terraform-provider-akamai/issues/330))
  * Fix problematic state file when attempting to change `akamai_cp_code` `group_id`([#322](https://github.com/akamai/terraform-provider-akamai/issues/322))
  * Fix panic on `akamai_property_rules_template` on empty property_snippets file ([#332](https://github.com/akamai/terraform-provider-akamai/issues/332))

## 2.2.0 (June 30, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Added penalty box support for security policy in evaluation mode

* IAM
  * Extended `akamai_iam_user`:
    * `is_locked` field has been deprecated in favor of `lock`
  * Added resources allowing management of:
    * `akamai_iam_blocked_user_properties` - create, read, update and import
    * `akamai_iam_group` - create, read, update, delete and import
    * `akamai_iam_role` - create, read, update, delete and import
  * Added data sources:
    * `akamai_iam_grantable_roles` - list grantable roles
    * `akamai_iam_timezones` - list supported timezones

#### BUG FIXES:

* APPSEC
  * Fix drift in `EffectiveTimePeriod`, `SamplingRate`, `LoggingOptions`, and `Operation` fields of custom rule resource.
  * Fix crash when eval rule API returns an error.
  * Fix incorrect error report when activation API returns an error.

## 2.1.1 (Jun 9, 2022)

#### BUG FIXES:

* Fix vulnerability for HashiCorp go-getter

## 2.1.0 (Jun 2, 2022)

#### FEATURES/ENHANCEMENTS:

* Support for Darwin ARM64 architecture ([GH#236](https://github.com/akamai/terraform-provider-akamai/issues/236))

* Image and Video Manager
  * New data sources:
    * `akamai_imaging_policy_image` - generate JSON for image policy
    * `akamai_imaging_policy_video` - generate JSON for video policy
  * Add `ImQuery` transformation
  * Add `Composite` transformation to `PostBreakpointTransformations`

#### BUG FIXES:

* PAPI
  * Update documentation for `akamai_property_rules_template`
  * Track remote changes in property rules ([#305](https://github.com/akamai/terraform-provider-akamai/issues/305))

* IAM
  * `akamai_iam_user`: remove phone number validation, to allow international phone number format

## 2.0.0 (Apr 28, 2022)

#### BREAKING CHANGES:

* APPSEC
  * Require version number for security configuration activation

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Add tuning recommendations for eval rulesets
  * Require security policy ID for bypass network list data source & resource

#### BUG FIXES:

* PAPI
  * Resource `akamai_property`: handle secure by default API errors

## 1.12.1 (Apr 6, 2022)

#### BUG FIXES:

* Added Image and Video Manager Documentation

* Include `terraform-provider-manifest`

## 1.12.0 (Mar 31, 2022)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added Image and Video Manager API support
  * Added resources allowing management of:
    * `akamai_imaging_policy_image` - create, read, update, delete and import
    * `akamai_imaging_policy_set` - create, read, update, delete and import
    * `akamai_imaging_policy_video` - create, read, update, delete and import

* CLOUDLETS
  * Support for RC cloudlet type (Request Control)

* PAPI
  * Added support for update `akamai_cp_code` resource
  * Added data source:
    * `akamai_properties_search` - list properties matching a specific hostname, edge hostname or property name

* Support for Go 1.17

#### BUG FIXES:

* DATASTREAM
  * Fix ordering sensitivity for JSON based configuration ([#287](https://github.com/akamai/terraform-provider-akamai/issues/287))

* PAPI
  * Fix CP code name forces replacement by adding update functionality in `akamai_cp_code` resource ([#262](https://github.com/akamai/terraform-provider-akamai/issues/262))

* Add metadata required by terraform registry

## 1.11.0 (Mar 3, 2022)

#### FEATURES/ENHANCEMENTS:

* [IMPORTANT] Added EdgeWorkers and EdgeKV API support
  * Added resources allowing management of:
    * EdgeWorker and EdgeWorker activations:
      * `akamai_edgeworker` - create, read, update, delete and import EdgeWorker
      * `akamai_edgeworkers_activation` - create, read, update and delete EdgeWorker activations
    * EdgeKV:
      * `akamai_edgekv` - create, read, update, delete and import an EdgeKV namespace
  * Added data sources for EdgeWorkers:
    * `akamai_edgeworkers_resource_tier` - lists information about resource tiers
    * `akamai_edgeworkers_property_rules` - generates property rule and behavior to associate an EdgeWorker to a property

* CLOUDLETS
  * Support for AS cloudlet type (Audience Segmentation)

#### BUG FIXES:

* APPSEC
  * Prevent 409 Conflict error caused by simultaneous network activation requests
  * Allow updating network list activation without destroying and recreating
  * Update unit tests to remove "NonEmptyPlanExpected" attribute

* CPS
  * Apply on resource `akamai_cps_dv_enrollment` is not idempotent if SANs contain common_name

## 1.10.1 (Feb 10, 2022)

#### FEATURES/ENHANCEMENTS:

* APPSEC
  * Cache OpenAPI calls for config & WAFMode information
  * Allow separate resources for individual protection settings

* CLOUDLETS
  * ALB cloudlet activation: allow modification of the `network` field without destroying the existing activation
  * Policy activation: allow modification of the `network` field without destroying the existing activation

#### BUG FIXES:

* CLOUDLETS
  * Changed schema for `akamai_cloudlets_application_load_balancer` resource, to fix struct validation error during update phase
  * Fixed client side validation to allow a datacenter percentage of 0% in `akamai_cloudlets_application_load_balancer` resource

* PAPI
  * Fix error in `akamai_property_activation` resource, which was blocking rolling back to any previous property activation ([#272](https://github.com/akamai/terraform-provider-akamai/issues/272))

## 1.10.0 (Jan 27, 2022)

#### FEATURES/ENHANCEMENTS:

* CLOUDLETS
  * Support for VP cloudlet type (Visitor Prioritization)
  * Support for CD cloudlet type (Continuous Deployment / Phased Release)
  * Support for FR cloudlet type (Forward Rewrite)
  * Support for AP cloudlet type (API Prioritization)

* APPSEC
  * Remove WAP-only datasource and resources
  * Add support for Evasive Path Match feature

* NETWORK LISTS
  * Include contract_id & group_id in akamai_networklist_network_lists datasource

* PAPI
  * Add support for array type variables in akamai_property_rules_template ([#257](https://github.com/akamai/terraform-provider-akamai/issues/257))

## 1.9.1 (Dec 16, 2021)

#### BUG FIXES:
* DNS
  * Refactored MX Bind processing and target suppress to fix failing import

## 1.9.0 (Dec 6, 2021)

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] Added Cloudlets API support
  * Added resources allowing management of policy and policy activations:
    * `akamai_cloudlets_policy` - create, read, update, delete and import policy
    * `akamai_cloudlets_policy_activation` - create, read, update and delete policy activations
  * Added resources allowing management of application load balancer configuration and application load balancer activations:
    * `akamai_cloudlets_application_load_balancer` - create, read, update, delete and import application load balancer configuration
    * `akamai_cloudlets_application_load_balancer_activation` - create, read, update and delete application load balancer activations
  * Added data sources:
    * `akamai_cloudlets_policy` - lists information about policy
    * `akamai_cloudlets_application_load_balancer` - lists information about application load balancer configuration
    * `akamai_cloudlets_application_load_balancer_match_rule` - lists information about application load balancer match rules
    * `akamai_cloudlets_edge_redirector_match_rule` - lists information about edge redirector match rules
* APPSEC
  * Add group/contract ID support to network list resource ([#243](https://github.com/akamai/terraform-provider-akamai/issues/243))
  * Add tuning recommendations data source
  * Add support for advanced exceptions in ASE rules
  * Update WAP bypass network lists for multi-policy WAP
  * Deprecate WAP-only datasource & resources
* PAPI
  * Updated documentation for data source akamai_property_rules
  * Allowed user to select a rule format in `resource akamai_property`
  * Added optional `use_cases` attribute for `akamai_edge_hostname` resource
#### BUG FIXES:
* Fixed example usage for provider import ([#212](https://github.com/akamai/terraform-provider-akamai/pull/212))
* PAPI
  * Removing default value for `ip_behavior` in `akamai_edge_hostname` resource ([#213](https://github.com/akamai/terraform-provider-akamai/pull/213))
  * Returned an error if `edge_hostname` attribute in `akamai_edge_hostname` resource does not exist  ([#258](https://github.com/akamai/terraform-provider-akamai/issues/258))
* CPS
  * Attribute `dns_challenges` should not be empty on initial apply for `akamai_cps_dv_enrollment` resource ([#253](https://github.com/akamai/terraform-provider-akamai/issues/253))
* DATASTREAM
  * Attribute `dataset_fields_ids` should not be sorted numerically in `akamai_datastream` resource ([#263](https://github.com/akamai/terraform-provider-akamai/issues/263))
* GTM
  * Attribute `datacenter_id` should be required in `akamai_gtm_geomap` resource ([#259](https://github.com/akamai/terraform-provider-akamai/issues/259))

## 1.8.0 (Oct 25, 2021)

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] DATASTREAM - Added DataStream configuration support
  * New [DataStream module](docs/guides/get_started_datastream.md). This module provides scalable, low latency streaming of property data in raw form
  * New resource:
    * [akamai_datastream](docs/resources/datastream.md) - create, read and update log streams
  * New data sources:
    * [akamai_datastream_activation_history](docs/data-sources/datastream_activation_history.md) - list detailed information about the activation status changes for all versions of a stream
    * [akamai_datastream_dataset_fields](docs/data-sources/datastream_dataset_fields.md) - list groups of data set fields available in the template
* PAPI
  * New [akamai_property_rules_template](docs/data-sources/property_rules_template.md) data source, which lets you use JSON template files to configure a rule tree

## 1.7.1 (Sept 30, 2021)

#### FEATURES/ENHANCEMENTS:
* PAPI
  * Handling `note` field during property deactivation
* APPSEC
  * Major documentation updates and clean up

#### BUG FIXES:
* PAPI
  * GRPC limit increased to 64MB ([#220](https://github.com/akamai/terraform-provider-akamai/issues/220))

## 1.7.0 (Aug 19, 2021)

#### FEATURES/ENHANCEMENTS:
* Terraform Plugin SDK updated to v2.7.0
* Provider tested and now supports Terraform 1.0.4

* APPSEC
  * Add wap_selected_hostnames data source and resource
  * Remove import templates for deprecated features
  * Display policy IDs for siem settings in separate table
  * Get an evaluation attack group's or risk score group's action
* NETWORK LISTS
  * Support contract_id and group_id for network list create/update
* PAPI
  * Possibility to set `note` field in property_activation resource
  * Additional checks and validations in `terraform plan` ([#245](https://github.com/akamai/terraform-provider-akamai/issues/245))

#### BUG FIXES:
* APPSEC
  * Configuration drift on reputation_profile create/apply
  * Fix incorrect comments/URL references in inline documentation
  * Data source akamai_appsec_security_policy returning incorrect policy ID
* DNS
  * Trim contract (ctr_) and group (grp_) prefixes when comparing configuration and TF state values ([#242](https://github.com/akamai/terraform-provider-akamai/issues/242))
* GTM
  * Trim contract (ctr_) and group (grp_) prefixes when comparing configuration and TF state values

## 1.6.1 (Jul 21, 2021)

#### BUG FIXES:
* DNS
  * Fixed contract id not being set in zone import and made group optional ([#242](https://github.com/akamai/terraform-provider-akamai/issues/242))
* GTM
  * Fixed documentation mismatch with optional/required fields on nested objects for `akamai_gmt_property` resource ([#240](https://github.com/akamai/terraform-provider-akamai/issues/240))
* PAPI
  * Fixed issue with property hostnames list changing order in diff ([#230](https://github.com/akamai/terraform-provider-akamai/issues/230))
  * Fixed idempotency issue on `akamai_property` resource ([#226](https://github.com/akamai/terraform-provider-akamai/issues/226))
  * Fixed issue with terraform showing misleading diff on `rules` field in `akamai_property` ([#234](https://github.com/akamai/terraform-provider-akamai/issues/234))
* CPS
  * Added `sans` field on `akamai_cps_dv_validation` to enable resending acknowledgement on after SANS are updated

#### FEATURES/ENHANCEMENTS:
* CPS
  * `akamai_cps_dv_enrollment` now accepts `contract_id` with `ctr_` prefix

## 1.6.0 (June 21, 2021)

#### BREAKING CHANGES:
* APPSEC
  * Configuration version numbers are no longer supported for most data sources and resources, as described below.
  * The following data sources are no longer supported:
    * akamai_appsec_attack_group_actions
    * akamai_appsec_attack_group_condition_exception
    * akamai_appsec_eval_rule_actions
    * akamai_appsec_eval_rule_condition_exception
    * akamai_appsec_rule_actions
    * akamai_appsec_rule_condition_exception
  * The following resources are no longer supported:
    * akamai_appsec_attack_group_action
    * akamai_appsec_attack_group_condition_exception
    * akamai_appsec_configuration_clone
    * akamai_appsec_configuration_version_clone
    * akamai_appsec_eval_rule_action
    * akamai_appsec_eval_rule_condition_exception
    * akamai_appsec_rule_action
    * akamai_appsec_rule_condition_exception
    * akamai_appsec_security_policy_clone
    * akamai_appsec_security_policy_protections

#### BUG FIXES:
* PAPI
   * Fixed issue causing edgehostnames not being set properly in state intermittently

#### FEATURES/ENHANCEMENTS:
* [IMPORTANT] CPS - Added Certificate Provisioning API support
  * Added resources allowing management of DV enrollments:
    * akamai_cps_dv_enrollment - create, read, update and delete DV enrollments
    * akamai_cps_dv_validation - inform CPS of finished validation, track change status

* APPSEC
  * The provider now determines automatically the version number to use for data source and resource operations.
    The most recent version of the specified configuration will be used if it is not currently active in either
    staging or production. If the most recent version is currently active, that version will be cloned and the
    newly cloned version will be used. The version attribute has been removed from all resource and data definitions,
    with the exception of the following data sources:
    * akamai_appsec_configuration_version
    * akamai_appsec_export_configuration
  * The export output templates supported by the akamai_appsec_export_configuration data source have been updated
    to remove version attributes.
  * The functionality for cloning and renaming configurations and security policies has been integrated into
    the respective resources. The separate resources for cloning and renaming have been removed. The affected
    elements are listed in the `BREAKING CHANGES` section above.
  * The action and condition_exception functionality for rule, eval-rule and attack-group resources have been
    consolidated into the respective data sources. The individual data sources and resources have been removed,
    and the remaining ones have been renamed. The affected elements are listed in the `BREAKING CHANGES` section above.
  * The akamai_appsec_activation resource's ForceNew attribute is no longer supported.
  * Resource updates that include modifications to the config_id or security_policy_id attributes are forbidden.
  * The akamai_appsec_siem_setting resource's output_text attribute is no longer supported.
  * The tabular output from the export_configuration data source has been improved.
  * The sample configuration file in the source repository has been updated to standardize names and remove
    version attributes.
  * Policy protections are now set individually. The separate resources for setting individual policy_protections
    resources has been removed.
  * The Getting Started guide for Appsec has been updated to include more information on importing resources, including
    a list of the supported output templates.
  * The following data sources have been added:
    * akamai_appsec_advanced_settings_pragma_header
    * akamai_appsec_attack_groups
    * akamai_appsec_eval_rules
    * akamai_appsec_rules
  * The following resources have been added:
    * akamai_appsec_advanced_settings_pragma_header
    * akamai_appsec_api_constraints_protection
    * akamai_appsec_attack_group
    * akamai_appsec_eval_rule
    * akamai_appsec_ip_geo_protection
    * akamai_appsec_rule

* PAPI
  * New optional parameter, which allows to import a specific property version.
    Additional information in [Property resource](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/property#import)

## 1.5.1 (Apr 21, 2021)

#### BUG FIXES:

* APPSEC
  * Suppress 'null' text on output of empty/false values
  * Prevent configuration drift when reapplying configuration after importing or creating resources
  * Update configuration version in local state file when modified in config.tf
  * Use uppercase when managing GEO network list elements
  * Display both API & website match targets in text_output
  * Remove unused output_text from code and documentation
  * Set network_list_id on network list import
  * Add comments to simplify importing resources using "terraform import"

* PAPI
   * Fixed issue causing inconsistent state when activation has rule errors ([#219](https://github.com/akamai/terraform-provider-akamai/issues/219))
   * Fixed issue with `resource_akamai_property` not setting product_id during import ([#224](https://github.com/akamai/terraform-provider-akamai/issues/224))
   * Rule warnings are not set in state anymore in `resource_akamai_property` and `resource_akamai_property_activation` to address size concerns of state file. Users will still be able to see them in logs as warnings

* DNS - Fix panic when zone already exists on create
* GTM - Deprecate and ignore Property field static_ttl. Add warning if present in property resource config

#### FEATURES/ENHANCEMENTS:
* PAPI - `resource_akamai_property_activation` now allows new optional argument `auto_acknowledge_rule_warnings`. Refer to [Property Activation Resource](docs/resources/property_activation.md)

## 1.5.0 (Mar 30, 2021) PAPI - Secure by default integration

#### BREAKING CHANGES:
* PAPI - `resource_akamai_property:` Changed hostnames field to a block type syntax to support additional user inputs. Refer to [Property Resource](docs/resources/property.md) for new syntax.

**Important Note**
Existing terraform users with hostnames defined in older syntax need to manually fix their hostnames configuration and existing state if needed. Additional info in [Property Resource](docs/resources/property.md)

#### BUG FIXES:

* PAPI
   * Fixed issue with version attributes not being set properly ([#208](https://github.com/akamai/terraform-provider-akamai/issues/208))
   * Fixed issue with `data_akamai_property_rules_template` not interpolating `#include` files properly
   * Fixed issue with `data_akamai_property_rules_template` not merging nested files properly

#### FEATURES/ENHANCEMENTS:
* PAPI
   * New [Hostnames Datasource](docs/data-sources/property_hostnames.md) to query hostnames and poll certificate status
   * Improved error handling and error messages in `property` and `property_activation` resources

## 1.4.0 (Mar 17, 2021) Network Lists

These are the operations supported in the Network Lists API v2:

* Create a network list
* Update an existing network list
* Get the existing network lists, optionally filtering by name or type
* Subscribe to a network list
* Activate a network list

## 1.3.1 (Mar 4, 2021)
* PAPI - Fixed issue with rules causing advanced locked behaviors to fail

## 1.3.0 (Feb 25, 2021) APPSEC - Extended list of supported list endpoints from APPSEC API

#### BREAKING CHANGES:
* PAPI
    * `data_akamai_property_rules_template:` snippets files should now be placed under `property-snippets` directory and should have `.json` extension

#### FEATURES/ENHANCEMENTS:
* APPSEC
    * Custom Deny
    * SIEM Setting
    * Advanced Options Settings
    * API Match Target
    * API Request Constraint
    * Create/Delete/Rename Security Policy
    * Host Coverage / Edit Version Notes
    * All WAP Features / WAP Hostname Evaluation
    * Create Security Configuration
    * Rename Security Configuration Version
    * Delete Security Configuration Version
    * Clone Security Configuration
    * Import tool for adding existing resources to Terraform state ([#207](https://github.com/akamai/terraform-provider-akamai/issues/207))
* DNS
    * Create SOA and NS Records on zone read if don't exist.
    * Add HTTPS, SVCB record support
* GTM
    * Add validation for property type and traffic targets combination

#### BUG FIXES:
* PAPI
    * Fixed issue causing hostnames to be appended instead of being replaced
    * Fixed issue causing version and rule comments being dropped ([#55](https://github.com/akamai/terraform-provider-akamai/issues/55))
    * Fixed client side validation to allow certain PAPI errors to passthrough
    * Fixed issue causing incorrect property version being stored in state for certain scenarios
* DNS
    * Suppress NS Record target diff if old and new equal without trailing 'period' ([#189](https://github.com/akamai/terraform-provider-akamai/issues/189))
    * Fail on attempted Zone deletion. Not supported.

## 1.2.1 (Feb 4, 2021)

#### BUG FIXES:
* PAPI -- Fixed crash caused by passing computed cpCode as a variable in rules to akamai_property
* PAPI -- Deprecated "product" attribute in akamai_cp_code resource and changed it "product_id"

## 1.2.0 (Jan 14, 2021) Identity and Access Management support

These are the operations supported in the Identity Management: User Administration API v2:

* Create a new user
* Update a user's profile
* Update a user's role assignments
* Delete a user

## 1.1.1 (Jan 8, 2021)
* APPSEC - Documentation formatting fixes

## 1.1.0 (Jan 7, 2021) APPSEC - Extended list of supported endpoints from APPSEC API:
  * DDoS Protection -- Rate Policy & Action
  * DDoS Protection -- Slowpost setting & Action
  * Application Layer Protection -- Rule Action, Exceptions & Conditions
  * Application Layer Protection -- Rule Evaluation Action, Exceptions & Conditions
  * Application Layer Protection -- Attack Group Action, Exceptions & Conditions
  * Application Layer Protection -- Rule Upgrade & Change Mode for Rule Eval
  * Reputation Profile & Action
  * Network Layer Control -- IP & GEO setting

## 1.0.0 (Dec 9, 2020) Provisioning redesign

#### BREAKING CHANGES:
* provider: configuring via an inline provider block (`property`, `dns`, or `gtm`) has been replaced with a more general `config` block that works the same way.
* There are several breaking changes in the 1.0 release.  You should consult the [Migration guide](docs/guides/1.0_migration.md) for details.
  * resources/akamai_property_activation no longer supports the following fields : activate.  version has gone from being optional to being a required field.
  * data-sources/akamai_property_rules removed in favor of using template JSON object to better work with other Akamai tools and documentation that is all JSON based.
  * resources/akamai_property_variables removed in favor of directly managing the variable segment as part of ruletree object.
  * resources/akamai_cp_code no longer auto-imports on create. If a conflict is detected will error out and to ignore simply import the resource.
  * resources/akamai_edge_hostname no longer supports the following fields : ipv4, ipv6. The revised resource allows setting ip_behavior directly.
  * resources/akamai_property no longer supports the following fields : cp_code, origin, variables, is_secure, contact. The revised resource simplifies the object structure and removes the ability to set the same value more than one way.
#### NOTES:
* provider/papi: changed attribute names in Provisioning to distinguish objects and names from id attributes.  In prior releases, "group" could represent a name, an id, or sometimes both. This release distinguishes them with distinct attribute names "group_name", "group_id" instead of "group"."
#### KNOWN BUGS:
* resources/akamai_property removing hostnames attribute can result in repeated noop update calls because in this case removal means the hostname relationships are un-managed leaving the attribute as empty is a better way to express this change.
#### FEATURES:
* data-sources/akamai_properties added to list properties accessible to the user.
* data-sources/akamai_property_contracts added to list contracts accessible to the user.
* data-sources/akamai_property_groups added to list groups accessible to the user.
* data-sources/akamai_property_products added to list products associated with a given contract.
* data-sources/akamai_property_rule_formats added to list rule_formats.
* data-sources/akamai_property_rules changed to output the structure of a particular rule version on the server. NOTE: this is NOT the same as the deprecated datasource used for rule formatting.
* data-sources/akamai_rules_template added to handle file based JSON templating for rules tree data management
#### ENHANCEMENTS:
* resources/akamai_property_activation aliased property to property_id. Returns these additional attributes : target_version, warnings, errors, activation_id, and status
#### BUG FIXES:
* provider: provider configuration validation requires an edgerc file configured and present even when environment variable-based configuration was used.
* provider: provider inline configuration support was re-introduced as a new config field.
* resources/akamai_property_activation activating and destroying activation for the same property multiple times in a row would fail on second destroy attempt and subsequent destroy attempts with "resource not found error" message.
* resources/akamai_property_activation wrong activation id read for property versions that had been activated and deactivated multiple time.
#### MINOR CHANGES:
* resources/akamai_property aliased property to property_id. contract to contract_id, and product to product_id and account to account_id.  Renamed version to latest_version.
* data-sources/akamai_contract aliased group to group_id and/or group_name.
* data-sources/akamai_cp_code aliased group to group_id and contract to contract_id.
* data-sources/akamai_group aliased name to group_name and contract to contract_id.

## 0.11.0 (Nov 19,2020)

#### NOTES:
* provider: Added support for application security API
#### BUG FIXES:
* provider: Updated edgegrid library to version 2.0.2. This should include the following fixes:
    * Re-enabled global account switch key support in edgerc files for reseller accounts.
    * PAPI - edgehostname updated returns - The System could not find cnameTo value
    * PAPI - property update return error - You provided an Etag that does not represent the last edit. Another edit has occurred, so check your request again before retrying.

## 0.10.2 (Oct 22,2020)
#### NOTES:
* Documentation formatting
#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when environment variable-based configuration was used.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.

## 0.10.1 (Not released)

## 0.10.0 (Oct 20,2020)

#### NOTES:
* provider: The backing edgegrid library was entirely rewritten.  Provider behavior should be preserved but there is chance of incidental changes due to the project size.
* resources/akamai_edge_hostname: edge_hostname field should be provided with an ending of edgesuite.net, edgekey.net, or akamaized.net.  If a required suffix is not provided then edgesuite.net is appended as default.
#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when one should not be needed.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.
#### ENHANCEMENTS:
* provider: improved error handling and improved message consistency
* provider: release notes categorize updates according to Terraform best practices guide.
* resources/akamai_cp: support ids with and without prefixes
* resources/akamai_edge_hostnames: support ids with and without prefixes
* resources/akamai_property: support ids with and without prefixes
* resources/akamai_property_activation: support ids with and without prefixes
#### BUG FIXES:
* resources/akamai_property: [AT-42] Fix criteria_match values handling
* provider: fixed documentation to properly present guides and categories on Hashicorp Terraform registry site
* resources/edge_hostname: added error when neither IPV4 nor IPV6 is selected
* resources/akamai_property: comparisons in rule tree now properly ignore equivalent values with attribute order differences.
* data-sources/akamai_property_rules: comparisons in rule tree now properly ignore equivalent values with attribute order differences.
* provider: updated all error messages to better identify issues and actions required by user
* provider: fixed crash due to unexpected data types from unexpected API responses
* provider: fixed crash due to unexpected data types in Terraform files
* provider: errors now get reported using Terraform diagnostics allowing much more detail to be passed to user when an error occurs.

## 0.9.1 (Sept 02, 2020)
#### BREAKING CHANGES:
* [IMPORTANT] Dropped support for TF clients <= 0.11. Provider now built using Terraform sdk v2 library. Terraform dropped 0.11 client support as part of this update.  This change will make many new enhancements possible. ([See: Terraform v2 sdk](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html))
* resources/akamai_group: contract field (previously optional) now required to ensure contract and group agreement.

#### NOTES:
* [CHANGE] Individual edgerc file sections for different Akamai APIs (i.e., `property_section`, `dns_section`) has been deprecated in favor a common `config_section` used in conjunction with provider aliases ([See: Multiple Provider Configurations](https://www.terraform.io/docs/configuration/providers.html#alias-multiple-provider-configurations))

#### KNOWN BUGS:
* provider: provider configuration validation requires an edgerc file configured and present even when one should not be needed.
* provider: support for configuring the provider via an inline provider block (`property`, `dns`, or `gtm`) no longer works.  Users should use edgerc file or Terraform environment args to configure instead.

#### BUG FIXES:
* [FIX] datasource akamai_group will no longer panic when contract not provided
* [ADD] Project re-organized to prepare for additional APIs to be included
* Fixed build job to compile sub-modules. Code is identical to 0.9.0 release

## 0.9.0 (August 26, 2020)
* [IMPORTANT] This build did not compile all modules properly so use 0.9.1 above instead.

## 0.8.2 (August 13, 2020)
* Initial release via the Terraform Registry. Otherwise identical to 0.8.1 release

## 0.8.1 (July 30, 2020)
* [FIX] Activation is executed, even without changes #139 (`akamai-property-activation`) ([#139](https://github.com/akamai/terraform-provider-template/issues/139))
* [FIX] Cannot find group when there are groups with the same name under multiple contract. #168 (`akamai-property-group`) ([#168](https://github.com/akamai/terraform-provider-template/issues/168))

## 0.8.0 (July 13, 2020)
* [FIX] Corrected Error 401 [Signature does not match] during  new primary zone creation (`akamai-dns`) ([#163](https://github.com/terraform-providers/terraform-provider-template/issues/163))
* [ADD] Updated Getting Started Primary Zone creation description. Added FAQ for Primary zone (`akamai-dns`)
* [FIX] SRV record priority value of 0 not allowed (`akamai-dns`) ([#165](https://github.com/terraform-providers/terraform-provider-template/issues/165))
* [ADD] Initial support for correlation ID in logging (`akamai-property`)

## 0.7.2 (June 11, 2020)
* [FIX] Corrected AAAA record handling of short and long IPv6 notation (`akamai-dns`)
## 0.7.1 (June 01, 2020)
* [FIX] Error after upgrading to 0.7.0 regarding MX records (`akamai-dns`) ([#154](https://github.com/terraform-providers/terraform-provider-template/issues/154))
* [FIX]Error 422 on SOA Record Apply After Creating a Primary Zone (`akamai-dns`) ([#155](https://github.com/terraform-providers/terraform-provider-template/issues/155))
## 0.7.0 (May 21, 2020)
* [ADD] User Agent support for Terraform version and provider version and SDK update
* [FIX] Bugs in Zone Create and Exists (`akamai_dns`) ([#151](https://github.com/terraform-providers/terraform-provider-template/issues/151))
## 0.6.0 (May 18, 2020)
* [ADD] Support the creation of DNS records of type AKAMAICDN (`akamai_dns`) ([#53](https://github.com/terraform-providers/terraform-provider-template/issues/53))
* [ADD] Support akamai_dns_record Import (`akamai_dns`) ([#69](https://github.com/terraform-providers/terraform-provider-template/issues/69))
* [FIX] Cannot remove a backup_cname from GTM property (`akamai_gtm`) ([#124](https://github.com/terraform-providers/terraform-provider-template/issues/124))
* [ADD] DNS Alias Zone Support (`akamai_dns`) ([#125](https://github.com/terraform-providers/terraform-provider-template/issues/125))
* [ADD] DNS TSIG Key support (`akamai_dns`) ([#126](https://github.com/terraform-providers/terraform-provider-template/issues/126))
* [ADD] DNS SOA, AKAMAITLC Record Support (`akamai_dns`) ([#127](https://github.com/terraform-providers/terraform-provider-template/issues/127))
* [FIX] Inverted Parameters - DNS Record Type NAPTR (`akamai_dns`) ([#130](https://github.com/terraform-providers/terraform-provider-template/issues/130))
* [FIX] Inverted Parameters - DNS Record Type NSEC3 (`akamai_dns`) ([#131](https://github.com/terraform-providers/terraform-provider-template/issues/131))
* [FIX] Inverted Parameters - DNS Record Type NSEC3PARAM (`akamai_dns`) ([#132](https://github.com/terraform-providers/terraform-provider-template/issues/132))
* [FIX] Inverted Parameters - DNS Record Type RRSIG (`akamai_dns`) ([#133](https://github.com/terraform-providers/terraform-provider-template/issues/133))
* [FIX] Inverted Parameters - DNS Record Type DS (`akamai_dns`) ([#134](https://github.com/terraform-providers/terraform-provider-template/issues/134))
* [ADD] DNS CAA, TLSA, CERT Record Support (`akamai_dns`) ([#148](https://github.com/terraform-providers/terraform-provider-template/issues/148))

## 0.5.0 (March 06, 2020)
* [FIX] Release edgehostnames and products caching edge library v0.9.10 (`akamai_property`)

## 0.4.0 (March 03, 2020)
* [FIX] Release contract group and cpcode caching edge library v0.9.9 (`akamai_property`)

## 0.3.0 (March 02, 2020)
* [FIX] Provider produced inconsistent final plan #88 add contract group and cpcode caching edge library v0.9.9 (`akamai_property`) ([#88](https://github.com/terraform-providers/terraform-provider-template/issues/88))

## 0.2.0 (February 28, 2020)
* [FIX] Bug - Origin values customhostheader #93 (`akamai_property`) ([#93](https://github.com/terraform-providers/terraform-provider-template/issues/93))
* [FIX] akamai 0.1.5 - err: rpc error: code = Unavailable desc = transport is closing #87 (`akamai_property`) ([#87](https://github.com/terraform-providers/terraform-provider-template/issues/87))
* [FIX] Errors in documentation: akamai_contract and akamai_cp_code #52 (`akamai_property`) ([#52](https://github.com/terraform-providers/terraform-provider-template/issues/52))
* [FIX] Provider produced inconsistent final plan #88 (`akamai_property`) ([#88](https://github.com/terraform-providers/terraform-provider-template/issues/88))
* [FIX] akamai_property_activation creation crashing with Error: rpc error: code = Unavailable desc = transport is closing #102 (`akamai_property`) ([#102](https://github.com/terraform-providers/terraform-provider-template/issues/102))
* [ADD] Add Support for GTM domains and contained elements (domain, datacenter, property, resource, cidrmap, geographicmap, asmap)

## 0.1.5 (January 06, 2020)

* [FIX] Criteria is always end up using must satisfy "all" (`akamai_property`) ([#81](https://github.com/terraform-providers/terraform-provider-template/issues/81))
* [FIX] Provider produced inconsistent final plan (`akamai_property_variables`) ([#82](https://github.com/terraform-providers/terraform-provider-template/issues/82))
* [FIX] Cannot create multiple types of records with the same name (`akamai_dns_record`) ([#11](https://github.com/terraform-providers/terraform-provider-template/issues/11))
* [FIX] akamai_property_activation resource - changing network field causes deactivation of version in staging (`akamai_property_activation`) ([#51](https://github.com/terraform-providers/terraform-provider-template/issues/51))
* [FIX] Multiple MX records creation issue (`akamai_dns_record`) ([#57](https://github.com/terraform-providers/terraform-provider-template/issues/57))

## 0.1.4 (December 06, 2019)
* [FIX] Add support for update of rules state (`akamai_property`) ([#66](https://github.com/terraform-providers/terraform-provider-template/issues/66))
* [FIX] Add support for masters being optional (`akamai_dns_zone`) ([#61](https://github.com/terraform-providers/terraform-provider-template/issues/61))
* [FIX] Create edge hostname 400 error Bad Request Request parameter Slot Number (`akamai_property`) ([#56](https://github.com/terraform-providers/terraform-provider-template/issues/56))
* [FIX] TXT record - State update failure due to sha verification issue (`akamai_dns_zone`) ([#58](https://github.com/terraform-providers/terraform-provider-template/issues/58))
## 0.1.3 (August 12, 2019)

* [FIX] Correct ordering of values for `SRV` records (`akamai_dns_record`) ([#17](https://github.com/terraform-providers/terraform-provider-template/issues/17))
* [FIX] IPV4-only hostnames no longer fail (`akamai_edge_hostname`) ([#21](https://github.com/terraform-providers/terraform-provider-template/issues/21))
* [FIX] Don't try to deactive any version but the current one (`akamai_property_activation`) ([#21](https://github.com/terraform-providers/terraform-provider-template/issues/21))
* [FIX] Fix crash in DNS record validation ([#27](https://github.com/terraform-providers/terraform-provider-template/issues/27))
* [FIX] SiteShield behavior translated correctly to JSON ([#10](https://github.com/terraform-providers/terraform-provider-template/issues/10)] [[#40](https://github.com/terraform-providers/terraform-provider-template/issues/40))
* [FIX] Property rules correctly update (all rules now removed correctly) ([#30](https://github.com/terraform-providers/terraform-provider-template/issues/30))
* [FIX] Property Hostnames correctly update (all hostnames are now removed correctly) ([#44](https://github.com/terraform-providers/terraform-provider-template/issues/44))
* [FIX] Property activation was using the activation ID to fetch the property ([#35](https://github.com/terraform-providers/terraform-provider-template/issues/35))
* [FIX] Ensure property supports `is_secure` for Enhanced TLS ([#42](https://github.com/terraform-providers/terraform-provider-template/issues/42))
* [FIX] Multiple fixes to provider configuration for auth configuration. ([#46](https://github.com/terraform-providers/terraform-provider-template/issues/46))
* [FIX] Ensure the latest version is activated when no `akamai_property_activation.version` is set ([#45](https://github.com/terraform-providers/terraform-provider-template/issues/45))
* [FIX] Multiple records (e.g. using `count`) should now be created correctly ([#11](https://github.com/terraform-providers/terraform-provider-template/issues/11))
* [CHANGE] `akamai_property_rules` has been changed to a data source to ensure dependant resources update correctly, the existing resource now emits an error in all operations ([#47](https://github.com/terraform-providers/terraform-provider-template/issues/47))
* [ADD] Make zone type (primary or secondary) case-insensitive ([#29](https://github.com/terraform-providers/terraform-provider-template/issues/29))

## 0.1.2 (July 26, 2019)

* [FIX] Fixed handling of CPCode behavior in rules.json
* [FIX] Fixed hostname complexity, now a simple `{"public.host" = "edge.host"}` map
* [FIX] Fixed accidental deactivations
* [ADD] Added explicit property and dns credential blocks to provider config
* [ADD] Added better validation to `akamai_dns_record`

## 0.1.1 (July 09, 2019)

* [FIX] Bug fixes

## 0.1.0 (June 19, 2019)

* Initial release
